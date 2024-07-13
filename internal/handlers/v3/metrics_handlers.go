package v3

import (
	"errors"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
)

type MetricHandlers struct {
	log            *slog.Logger
	metricsService *metricsservice.MetricService
}

func New(log *slog.Logger, metricsService *metricsservice.MetricService) MetricHandlers {
	return MetricHandlers{log: log, metricsService: metricsService}
}

func (m *MetricHandlers) UpdateSeveralMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responseError(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var reqMetrics []Metrics
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

	var metric models.MetricInteraction
	var severalMetrics []models.MetricInteraction
	for _, oneMetric := range reqMetrics {

		if err := validator.New().Struct(oneMetric); err != nil {
			validateErr := err.(validator.ValidationErrors)
			errorText := ValidationError(validateErr)
			responseError(w, r, http.StatusBadRequest, errorText)
			return
		}

		// TODO must be in service layer
		if oneMetric.MType == configserver.MetricTypeCounter {
			metric, err = models.New(
				oneMetric.MType,
				oneMetric.ID,
				fmt.Sprintf("%d", *oneMetric.Delta),
			)
		} else {
			metric, err = models.New(
				oneMetric.MType,
				oneMetric.ID,
				fmt.Sprintf("%g", *oneMetric.Value),
			)
		}
		severalMetrics = append(severalMetrics, metric)
	}
	if err != nil {
		responseError(w, r, http.StatusBadRequest, "metric value conversion error")
		return
	}

	ctx := r.Context()
	err = m.metricsService.UpdateSeveralMetrics(ctx, severalMetrics)
	if err != nil {
		if errors.Is(err, metricsservice.ErrNotValidURL) {
			responseError(w, r, http.StatusNotFound, err.Error())
			return
		}
		responseError(w, r, http.StatusInternalServerError, "internal server error")
		return
	}
	responseOK(w, r, metric)
}
