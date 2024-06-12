package updatemetric

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/appserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/internal/http-server/v1/metrics"
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func PathValidator(r *http.Request) (models.Metric, error) {

	metricType := chi.URLParam(r, "metric_type")

	if metricType != "gauge" && metricType != "counter" {
		return models.Metric{}, metrics.ErrNotValidMetricType
	}

	metricValue := chi.URLParam(r, "metric_value")

	var value interface{} // Store the parsed value here
	var err error

	if metricType == "gauge" {
		value, err = strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return models.Metric{}, metrics.ErrNotValidMetricValue
		}
	} else {
		value, err = strconv.ParseUint(metricValue, 10, 64)
		if err != nil {
			return models.Metric{}, metrics.ErrNotValidMetricValue
		}
	}

	return models.Metric{
		Type:  chi.URLParam(r, "metric_type"),
		Name:  strings.ToLower(chi.URLParam(r, "metric_name")),
		Value: value,
	}, nil
}

func New(log *slog.Logger, application *appserver.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		metric, err := PathValidator(r)

		if errors.Is(err, metrics.ErrNotValidMetricType) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, metrics.ErrNotValidMetricValue) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if application.Cfg.ClientTimeout == 0 {
			application.Cfg.ClientTimeout = 10
		}
		timeout := time.Duration(application.Cfg.ClientTimeout) * time.Second
		fmt.Println("111111111111111111", timeout)
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
