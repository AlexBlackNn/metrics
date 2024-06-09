package main

import (
	"github.com/AlexBlackNn/metrics/internal/appserver"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/stretchr/testify/suite"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

type MetricsSuite struct {
	suite.Suite
	cfg         *config.Config
	log         *slog.Logger
	application *appserver.App
	client      http.Client
}

func (ms *MetricsSuite) SetupTest() {
	ms.cfg = &config.Config{
		Env:            "local",
		PollInterval:   2,
		ReportInterval: 5,
	}

	ms.log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ms.application = appserver.New(ms.log, ms.cfg)
	ms.client = http.Client{Timeout: 3 * time.Second}

}

func (ms *MetricsSuite) TestServerHappyPath() {

	type Want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name string
		url  string
		want Want
	}{
		{
			name: "gauge with value 10.3",
			url:  "/update/gauge/Lookups/10.3",
			want: Want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "counter with value 10",
			url:  "/update/counter/PoolCount/10",
			want: Want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	// запускаем тестовый сервер, будет выбран первый свободный порт
	srv := httptest.NewServer(NewChiRouter(ms.log, ms.application))
	// останавливаем сервер после завершения теста
	defer srv.Close()

	for _, tt := range tests {
		ms.Run(tt.name, func() {
			url := srv.URL + tt.url
			request, err := http.NewRequest(http.MethodPost, url, nil)
			ms.NoError(err)
			res, err := ms.client.Do(request)
			ms.NoError(err)
			// проверяем код ответа
			ms.Equal(tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			ms.Equal(tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func (ms *MetricsSuite) TestNegativeCasesMetrics() {

	type Want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name string
		url  string
		want Want
	}{
		{
			name: "metric name absent",
			url:  "/update/counter/10",
			want: Want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "metric wrong type value",
			url:  "/update/counter/PoolCount/test",
			want: Want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "metric wrong metric type",
			url:  "/update/histogram/PoolCount/10",
			want: Want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	// запускаем тестовый сервер, будет выбран первый свободный порт
	srv := httptest.NewServer(NewChiRouter(ms.log, ms.application))
	// останавливаем сервер после завершения теста
	defer srv.Close()

	for _, tt := range tests {
		ms.Run(tt.name, func() {
			request, err := http.NewRequest(http.MethodPost, srv.URL+tt.url, nil)
			ms.NoError(err)
			res, err := ms.client.Do(request)
			ms.NoError(err)
			// проверяем код ответа
			ms.Equal(tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			ms.Equal(tt.want.contentType, res.Header.Get("Content-Type"))

		})
	}
}

func (ms *MetricsSuite) TestNegativeCasesRequestMethods() {

	type Want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name   string
		url    string
		method string
		want   Want
	}{
		{
			name:   "get request",
			url:    "/update/counter/10",
			method: http.MethodGet,
			want: Want{
				code: http.StatusNotFound,
			},
		},
		{
			name:   "delete request",
			url:    "/update/counter/10",
			method: http.MethodDelete,
			want: Want{
				code: http.StatusNotFound,
			},
		},
		{
			name:   "put request",
			url:    "/update/counter/10",
			method: http.MethodPut,
			want: Want{
				code: http.StatusNotFound,
			},
		},
		{
			name:   "patch request",
			url:    "/update/counter/10",
			method: http.MethodPatch,
			want: Want{
				code: http.StatusNotFound,
			},
		},
	}

	// запускаем тестовый сервер, будет выбран первый свободный порт
	srv := httptest.NewServer(NewChiRouter(ms.log, ms.application))

	for _, tt := range tests {
		ms.Run(tt.name, func() {
			request, err := http.NewRequest(tt.method, srv.URL+tt.url, nil)
			ms.NoError(err)
			res, err := ms.client.Do(request)
			ms.NoError(err)
			defer res.Body.Close()
			// проверяем код ответа
			ms.Equal(tt.want.code, res.StatusCode)
		})
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(MetricsSuite))
}
