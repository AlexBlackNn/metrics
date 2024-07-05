package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/pkg/storage/mockstorage"
	"github.com/golang/mock/gomock"
	"io"
	"net/http"
)

func (ms *MetricsSuite) TestServerHappyPathMockStorageV2() {
	ctrl := gomock.NewController(ms.T())
	defer ctrl.Finish()

	mockStorage := mockstorage.NewMockMetricsStorage(ctrl)

	ctx := context.Background()
	modelTest := &models.Metric[uint64]{Type: "counter", Name: "test_counter", Value: 22}
	mockStorage.EXPECT().GetMetric(ctx, "test_counter").Return(modelTest, nil)
	result, err := mockStorage.GetMetric(ctx, "test_counter")
	ms.Suite.NoError(err)
	ms.Suite.Equal(modelTest, result)

	ms.application.DataBase = mockStorage

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
				"type": "counter"
				}`,
			),
			want: Want{
				code:        http.StatusOK,
				contentType: "application/json",
				response:    `{"id":"test_counter","type":"counter","delta":10}`,
			},
		},
	}
	// stop server when tests finished
	defer ms.srv.Close()

	for _, tt := range tests {
		ms.Run(tt.name, func() {
			url := ms.srv.URL + tt.url
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(tt.body))
			ms.NoError(err)
			res, err := ms.client.Do(request)
			ms.NoError(err)
			ms.Equal(tt.want.code, res.StatusCode)
			defer res.Body.Close()
			ms.Equal(tt.want.contentType, res.Header.Get("Content-Type"))
			data, err := io.ReadAll(res.Body)
			fmt.Println("1111111111", string(data))
			ms.NoError(err)
			ms.Equal(tt.want.response, string(data))
		})
	}
}
