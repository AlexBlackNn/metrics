package main

import (
	"bytes"
	"github.com/AlexBlackNn/metrics/app/server"
	"github.com/AlexBlackNn/metrics/cmd/server/router"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/internal/logger"
	"github.com/AlexBlackNn/metrics/pkg/storage/mockstorage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServerHappyPathMockStorageV2(t *testing.T) {
	// создадим конроллер моков и экземпляр мок-хранилища
	cfg, err := configserver.New()
	if err != nil {
		panic(err)
	}
	log := logger.New(cfg.Env)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMetricStorage := mockstorage.NewMockMetricsStorage(ctrl)
	mockHealthChecker := mockstorage.NewMockHealthChecker(ctrl)

	// определим, какой результат будем получать от «хранилища»
	modelTest := &models.Metric[uint64]{Type: "counter", Name: "test_counter", Value: 10}

	// установим условие: при любом вызове метода ListMessages возвращать массив messages без ошибки
	mockMetricStorage.EXPECT().
		GetMetric(gomock.Any(), gomock.Any()).
		Return(modelTest, nil)

	application, err := server.NewAppInitStorage(mockMetricStorage, mockHealthChecker, cfg, log)
	srv := httptest.NewServer(router.NewChiRouter(application.Cfg, application.Log, application.HandlersV1, application.HandlersV2))
	defer srv.Close()

	client := http.Client{Timeout: 3 * time.Second}

	type Want struct {
		code        int
		response    string
		contentType string
		value       int
	}

	tests := []struct {
		name string
		url  string
		body []byte
		want Want
	}{
		{
			name: "counter with value 10",
			url:  "/value/",
			body: []byte(
				`{
				"id": "test_counter",
				"type": "counter",
				"delta": 10
				}`,
			),
			want: Want{
				code:        http.StatusOK,
				contentType: "application/json",
				response:    `{"id":"test_counter","type":"counter","delta":10}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := srv.URL + tt.url
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(tt.body))
			assert.NoError(t, err, "error creating HTTP request")
			res, err := client.Do(request)
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer res.Body.Close()
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			data, err := io.ReadAll(res.Body)
			assert.NoError(t, err, "error reading response body")
			assert.Equal(t, tt.want.response, string(data))
		})
	}
}
