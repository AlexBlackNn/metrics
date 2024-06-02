package update

import (
	"context"
	"errors"
	"github.com/AlexBlackNn/metrics/internal/app"
	"github.com/AlexBlackNn/metrics/internal/services/metrics_service"
	"log/slog"
	"net/http"
	"time"
)

// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
// @Summary J,y
// @Description Создает новое выражение на сервере
// @Tags Calculations
// @Accept json
// @Produce json
// @Param body body Request true "Запрос на создание выражения"
// @Success 201 {object} Response
// @Router /expression [post]
// @Security BearerAuth
func New(log *slog.Logger, application *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		err := application.MetricsService.UpdateMetricValue(context.Background(), r.URL.Path)
		if errors.Is(err, metrics_service.ErrNotValidUrl) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}
}
