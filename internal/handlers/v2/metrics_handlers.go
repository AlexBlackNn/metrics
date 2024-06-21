package v2

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/internal/lib/response"
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Metrics struct {
	ID    string   `json:"id"`                                  // metrics name
	MType string   `json:"type" validate:"oneof=gauge counter"` // mType = counter || gauge
	Delta *int64   `json:"delta,omitempty"`                     // exists if mType = counter
	Value *float64 `json:"value,omitempty"`                     // exists if mType = gauge
}

func responseOK(w http.ResponseWriter, r *http.Request, metric models.MetricInteraction) {
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
}

func responseError(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	render.Status(r, statusCode)
	render.JSON(w, r, response.Error(message))
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
		responseError(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var reqMetrics Metrics
	err := render.DecodeJSON(r.Body, &reqMetrics)
	if err != nil {
		if errors.Is(err, io.EOF) {
			// Post with empty body
			responseError(w, r, http.StatusBadRequest, "empty request")
			return
		}
		responseError(w, r, http.StatusBadRequest, "failed to decode request")
		return
	}
	if err := validator.New().Struct(reqMetrics); err != nil {
		validateErr := err.(validator.ValidationErrors)
		errorText := response.ValidationError(validateErr)
		responseError(w, r, http.StatusBadRequest, errorText)
		return
	}
	ctx := context.Background()
	metric, err := m.metricsService.GetOneMetricValue(
		ctx, strings.ToLower(reqMetrics.ID),
	)

	if err != nil {
		if errors.Is(err, metricsservice.ErrMetricNotFound) {
			responseError(w, r, http.StatusNotFound, "metric not found")
			return
		}
		responseError(w, r, http.StatusInternalServerError, "internal server error")
		return
	}
	responseOK(w, r, metric)
}

func (m *MetricHandlers) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responseError(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var reqMetrics Metrics
	err := render.DecodeJSON(r.Body, &reqMetrics)
	if err != nil {
		if errors.Is(err, io.EOF) {
			// Post with empty body
			responseError(w, r, http.StatusBadRequest, "empty request")
			return
		}
		responseError(w, r, http.StatusBadRequest, "failed to decode request")
		return
	}
	if err := validator.New().Struct(reqMetrics); err != nil {
		validateErr := err.(validator.ValidationErrors)
		errorText := response.ValidationError(validateErr)
		responseError(w, r, http.StatusBadRequest, errorText)
		return
	}

	var metric models.MetricInteraction

	if reqMetrics.MType == "counter" {
		metric, err = models.New(
			reqMetrics.MType,
			reqMetrics.ID,
			fmt.Sprintf("%d", *reqMetrics.Delta),
		)
	} else {
		metric, err = models.New(
			reqMetrics.MType,
			reqMetrics.ID,
			fmt.Sprintf("%g", *reqMetrics.Value),
		)
	}
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
