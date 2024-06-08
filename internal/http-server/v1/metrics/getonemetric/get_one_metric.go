package getonemetric

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/appserver"
	"github.com/AlexBlackNn/metrics/internal/http-server/v1/metrics"
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func PathValidator(r *http.Request) (string, error) {
	fmt.Println()
	metricType := chi.URLParam(r, "metric_type")

	if metricType != "gauge" && metricType != "counter" {
		return "", metrics.ErrNotValidMetricType
	}

	return strings.ToLower(chi.URLParam(r, "metric_name")), nil
}

func New(log *slog.Logger, application *appserver.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		key, err := PathValidator(r)

		if errors.Is(err, metrics.ErrNotValidMetricType) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		metric, err := application.MetricsService.GetOneMetricValue(context.Background(), key)

		if errors.Is(err, metricsservice.ErrMetricNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%v", metric.Value)))
	}
}
