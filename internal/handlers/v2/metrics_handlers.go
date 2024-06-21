package v2

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/internal/lib/response"
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Metrics struct {
	ID    string   `json:"id"`                                  // имя метрики
	MType string   `json:"type" validate:"oneof=gauge counter"` // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"`                     // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"`                     // значение метрики в случае передачи gauge
}

type MetricHandlers struct {
	log            *slog.Logger
	metricsService *metricsservice.MetricService
}

func New(log *slog.Logger, metricsService *metricsservice.MetricService) MetricHandlers {
	return MetricHandlers{log: log, metricsService: metricsService}
}

func (m *MetricHandlers) GetOneMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("1111")
	var reqMetrics Metrics
	err := render.DecodeJSON(r.Body, &reqMetrics)
	if err != nil {
		if errors.Is(err, io.EOF) {
			// Post with empty body
			fmt.Println("2222 err")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("empty request"))
			return
		}
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to decode request"))
		fmt.Println("333 err")
		return
	}
	fmt.Println("44444")
	if err := validator.New().Struct(reqMetrics); err != nil {
		fmt.Println("555 err")
		validateErr := err.(validator.ValidationErrors)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ValidationError(validateErr))
		return
	}
	fmt.Println("6666")
	ctx := context.Background()
	metric, err := m.metricsService.GetOneMetricValue(
		ctx, strings.ToLower(reqMetrics.ID),
	)

	if err != nil {
		if errors.Is(err, metricsservice.ErrMetricNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if metric.GetType() == "counter" {
		metricValue := int64(metric.GetValue().(uint64))
		render.JSON(w, r, Metrics{
			ID:    metric.GetName(),
			MType: metric.GetType(),
			Delta: &metricValue,
		})
		return
	}
	metricValue := metric.GetValue().(float64)
	render.JSON(w, r, Metrics{
		ID:    metric.GetName(),
		MType: metric.GetType(),
		Value: &metricValue,
	})
	return
}

func (m *MetricHandlers) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metric, err := models.New(
		chi.URLParam(r, "metric_type"),
		chi.URLParam(r, "metric_name"),
		chi.URLParam(r, "metric_value"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err = m.metricsService.UpdateMetricValue(ctx, metric)
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
