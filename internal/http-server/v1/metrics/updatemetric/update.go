package updatemetric

import (
	"context"
	"errors"
	"github.com/AlexBlackNn/metrics/internal/appserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"time"
)

func New(log *slog.Logger, application *appserver.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		metric, err := models.Load(
			chi.URLParam(r, "metric_type"),
			chi.URLParam(r, "metric_name"),
			chi.URLParam(r, "metric_value"),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// TODO: in some tests somehow ClientTimeout gets 0, which creates DEADLINE ERROR
		if application.Cfg.ClientTimeout == 0 {
			application.Cfg.ClientTimeout = 10
		}
		timeout := time.Duration(application.Cfg.ClientTimeout) * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		err = application.MetricsService.UpdateMetricValue(ctx, metric)
		if errors.Is(err, metricsservice.ErrNotValidURL) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

	}
}
