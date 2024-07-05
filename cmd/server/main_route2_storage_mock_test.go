package main

import (
	"bytes"
	"context"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/pkg/storage/mockstorage"
	"github.com/golang/mock/gomock"
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
	}

	tests := []struct {
		name string
		url  string
		body []byte
		want Want
	}{
		{
			name: "counter with value 10",
			url:  "/update/",
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
		})
	}
}
