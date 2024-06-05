package main

import (
	"github.com/AlexBlackNn/metrics/internal/app_server"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/http-server/v1/metrics/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestServerHappyPath(t *testing.T) {

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

	cfg := &config.Config{
		Env:            "local",
		ServerAddr:     ":8080",
		PollInterval:   2,
		ReportInterval: 5,
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	application := app_server.New(log, cfg)

	handlerUnderTest := update.New(log, application)
	client := http.Client{Timeout: 3 * time.Second}

	// запускаем тестовый сервер, будет выбран первый свободный порт
	srv := httptest.NewServer(handlerUnderTest)
	// останавливаем сервер после завершения теста
	defer srv.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := srv.URL + tt.url
			request, err := http.NewRequest(http.MethodPost, url, nil)
			require.NoError(t, err)
			res, err := client.Do(request)
			require.NoError(t, err)
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestNegativeCasesMetrics(t *testing.T) {

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

	cfg := &config.Config{
		Env:            "local",
		ServerAddr:     ":8080",
		PollInterval:   2,
		ReportInterval: 5,
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	application := app_server.New(log, cfg)

	handlerUnderTest := update.New(log, application)
	client := http.Client{Timeout: 3 * time.Second}

	// запускаем тестовый сервер, будет выбран первый свободный порт
	srv := httptest.NewServer(handlerUnderTest)
	// останавливаем сервер после завершения теста
	defer srv.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodPost, srv.URL+tt.url, nil)
			require.NoError(t, err)
			res, err := client.Do(request)
			require.NoError(t, err)
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))

		})
	}
}

func TestNegativeCasesRequestMethods(t *testing.T) {

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
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name:   "delete request",
			url:    "/update/counter/10",
			method: http.MethodDelete,
			want: Want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name:   "put request",
			url:    "/update/counter/10",
			method: http.MethodPut,
			want: Want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name:   "patch request",
			url:    "/update/counter/10",
			method: http.MethodPatch,
			want: Want{
				code: http.StatusMethodNotAllowed,
			},
		},
	}

	cfg := &config.Config{
		Env:            "local",
		ServerAddr:     ":8080",
		PollInterval:   2,
		ReportInterval: 5,
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	application := app_server.New(log, cfg)

	handlerUnderTest := update.New(log, application)
	client := http.Client{Timeout: 3 * time.Second}

	// запускаем тестовый сервер, будет выбран первый свободный порт
	srv := httptest.NewServer(handlerUnderTest)
	// останавливаем сервер после завершения теста
	defer srv.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(tt.method, srv.URL+tt.url, nil)
			require.NoError(t, err)
			res, err := client.Do(request)
			require.NoError(t, err)
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
		})
	}
}
