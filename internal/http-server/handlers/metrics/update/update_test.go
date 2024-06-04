package update

import (
	"github.com/AlexBlackNn/metrics/internal/app_server"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNew(t *testing.T) {

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
			name: "positive test,  gauge with value 10.3",
			url:  "/update/gauge/Lookups/10.3",
			want: Want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "positive test, counter with value 10",
			url:  "/update/counter/PoolCount/10",
			want: Want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test, metric name absent",
			url:  "/update/counter/10",
			want: Want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test, metric wrong type value",
			url:  "/update/counter/PoolCount/test",
			want: Want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test, metric wrong metric type",
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

	handlerUnderTest := New(log, application)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.url, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			handlerUnderTest(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()

			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
