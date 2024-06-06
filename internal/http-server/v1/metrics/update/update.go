package update

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
	"strconv"
	"time"
)

func PathValidator(r *http.Request) (models.Metric, error) {
	fmt.Println()
	metricType := chi.URLParam(r, "metric_type")

	fmt.Println("0000000000000000000", metricType, chi.URLParam(r, "metric_type"), chi.URLParam(r, "metric_name"), chi.URLParam(r, "metric_value"))
	if metricType != "gauge" && metricType != "counter" {
		fmt.Println("777777777777")
		return models.Metric{}, ErrNotValidMetricType
	}

	metricValue := chi.URLParam(r, "metric_value")
	var value interface{} // Store the parsed value here
	var err error

	fmt.Println("11111111")
	if metricType == "gauge" {
		value, err = strconv.ParseFloat(metricValue, 64)
		fmt.Println("222222")
		if err != nil {
			fmt.Println("333333333")
			return models.Metric{}, ErrNotValidMetricValue
		}
	} else {
		value, err = strconv.ParseUint(metricValue, 10, 64)
		fmt.Println("444444444")
		if err != nil {
			fmt.Println("55555555")
			return models.Metric{}, ErrNotValidMetricValue
		}
	}
	fmt.Println("6666")
	return models.Metric{
		Type:  chi.URLParam(r, "metric_type"),
		Name:  chi.URLParam(r, "metric_name"),
		Value: value,
	}, nil
}

func New(log *slog.Logger, application *appserver.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			fmt.Println("qqqqqqqqqqqqqqqqqq")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		fmt.Println("ppppppppppppppppppppppp", chi.URLParam(r, "metric_type"))

		metric, err := PathValidator(r)

		if errors.Is(err, ErrNotValidMetricType) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, ErrNotValidMetricValue) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = application.MetricsService.UpdateMetricValue(context.Background(), metric)
		if errors.Is(err, metricsservice.ErrNotValidURL) {
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
