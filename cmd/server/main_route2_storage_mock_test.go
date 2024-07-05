package main

//
//import (
//	"bytes"
//	"context"
//	"github.com/AlexBlackNn/metrics/app/server"
//	"github.com/AlexBlackNn/metrics/cmd/server/router"
//	"github.com/AlexBlackNn/metrics/internal/domain/models"
//	"github.com/AlexBlackNn/metrics/pkg/storage/mockstorage"
//	"github.com/golang/mock/gomock"
//	"github.com/stretchr/testify/suite"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//	"time"
//)
//
//type MetricsSuiteV2 struct {
//	suite.Suite
//	application *server.App
//	client      http.Client
//	srv         *httptest.Server
//}
//
//func (ms *MetricsSuiteV2) SetupSuite() {
//	var err error
//	ms.application, err = server.New()
//
//	ctrl := gomock.NewController(ms.T())
//	defer ctrl.Finish()
//
//	mockStorage := mockstorage.NewMockMetricsStorage(ctrl)
//	ctx := context.Background()
//	modelTest := &models.Metric[uint64]{Type: "counter", Name: "test_counter", Value: 22}
//	mockStorage.EXPECT().GetMetric(ctx, "test_counter").Return(modelTest, nil)
//	ms.application.DataBase = mockStorage
//
//	if err != nil {
//		ms.T().Fatal(err)
//	}
//	ms.client = http.Client{Timeout: 3 * time.Second}
//}
//
//func (ms *MetricsSuiteV2) BeforeTest(suiteName, testName string) {
//	// Starts server with first random port.
//	ms.srv = httptest.NewServer(router.NewChiRouter(ms.application.Cfg, ms.application.Log, ms.application.HandlersV1, ms.application.HandlersV2))
//}
//
//func (ms *MetricsSuiteV2) AfterTest(suiteName, testName string) {
//	ms.srv = nil
//}
//
//func (ms *MetricsSuiteV2) TestServerHappyPathMockStorageV2() {
//	type Want struct {
//		code        int
//		response    string
//		contentType string
//	}
//
//	tests := []struct {
//		name string
//		url  string
//		body []byte
//		want Want
//	}{
//		{
//			name: "counter with value 10",
//			url:  "/update/",
//			body: []byte(
//				`{
//				"id": "test_counter",
//				"type": "counter",
//				"delta": 10
//				}`,
//			),
//			want: Want{
//				code:        http.StatusOK,
//				contentType: "application/json",
//			},
//		},
//	}
//	// stop server when tests finished
//	defer ms.srv.Close()
//
//	for _, tt := range tests {
//		ms.Run(tt.name, func() {
//			url := ms.srv.URL + tt.url
//			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(tt.body))
//			ms.NoError(err)
//			res, err := ms.client.Do(request)
//			ms.NoError(err)
//			ms.Equal(tt.want.code, res.StatusCode)
//			defer res.Body.Close()
//			ms.Equal(tt.want.contentType, res.Header.Get("Content-Type"))
//		})
//	}
//}
//
//func TestSuiteV2(t *testing.T) {
//	suite.Run(t, new(MetricsSuiteV2))
//}
