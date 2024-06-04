package update

import (
	"github.com/AlexBlackNn/metrics/internal/app_server"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/utils"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
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
		want Want
	}{{
		name: "positive test #1",
		want: Want{
			code:        200,
			contentType: "text/plain; charset=utf-8",
		}},
	}

	cfg := &config.Config{
		Env:            "local",
		ServerAddr:     ":8080",
		PollInterval:   2,
		ReportInterval: 5,
	}

	// init logger
	log := utils.SetupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))
	application := app_server.New(log, cfg)
	handlerUnderTest := New(log, application)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/update/gauge/Lookups/10", nil)
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
