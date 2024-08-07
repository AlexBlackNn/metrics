package v2

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
	"strings"
)

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
	// TODO: easyjson can be used?
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
		errorText := ValidationError(validateErr)
		responseError(w, r, http.StatusBadRequest, errorText)
		return
	}
	ctx := r.Context()
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
		errorText := ValidationError(validateErr)
		responseError(w, r, http.StatusBadRequest, errorText)
		return
	}

	var metric models.MetricInteraction

	// TODO must be in service layer
	if reqMetrics.MType == configserver.MetricTypeCounter {
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
		responseError(w, r, http.StatusBadRequest, "metric value conversion error")
		return
	}

	ctx := r.Context()
	err = m.metricsService.UpdateMetricValue(ctx, metric)
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
