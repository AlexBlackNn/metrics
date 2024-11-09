package tests

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/AlexBlackNn/metrics/app/server"
	"github.com/AlexBlackNn/metrics/cmd/server/router"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/internal/logger"
	"github.com/AlexBlackNn/metrics/pkg/storage/mockstorage"
	"github.com/golang/mock/gomock"
)

func (ms *MetricsSuite) TestServerHappyPathMockStorageV2() {

	ctrl := gomock.NewController(ms.T())
	defer ctrl.Finish()
	mockMetricStorage := mockstorage.NewMockMetricsStorage(ctrl)
	mockHealthChecker := mockstorage.NewMockHealthChecker(ctrl)

	// определим, какой результат будем получать от «хранилища»
	modelTest := &models.Metric[uint64]{Type: "counter", Name: "test_counter", Value: 10}

	// установим условие: при любом вызове метода ListMessages возвращать массив messages без ошибки
	mockMetricStorage.EXPECT().
		GetMetric(gomock.Any(), gomock.Any()).
		Return(modelTest, nil)

	log := logger.New(ms.cfg.Env)
	application, err := server.NewAppInitStorage(mockMetricStorage, mockHealthChecker, ms.cfg, log)
	ms.NoError(err)
	srv := httptest.NewServer(router.NewChiRouter(
		application.Cfg,
		application.Log,
		application.HandlersV1,
		application.HandlersV2,
		application.HealthHandlersV2,
		application.HandlersV3),
	)
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
		ms.Run(tt.name, func() {
			url := srv.URL + tt.url
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(tt.body))
			ms.NoError(err, "error creating HTTP request")
			res, err := client.Do(request)
			ms.NoError(err, "error making HTTP request")
			ms.Equal(tt.want.code, res.StatusCode)
			defer res.Body.Close()
			ms.Equal(tt.want.contentType, res.Header.Get("Content-Type"))
			data, err := io.ReadAll(res.Body)
			ms.NoError(err, "error reading response body")
			ms.JSONEq(tt.want.response, string(data))
		})
	}
}
