package getonemetric

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/appserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func New(log *slog.Logger, application *appserver.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		err := models.CheckModelType(chi.URLParam(r, "metric_type"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		metric, err := application.MetricsService.GetOneMetricValue(
			context.Background(), strings.ToLower(chi.URLParam(r, "metric_name")),
		)

		if errors.Is(err, metricsservice.ErrMetricNotFound) {
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
		w.Write([]byte(fmt.Sprintf("%v", metric.Value)))
	}
}
