package v2

import (
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	"log/slog"
	"net/http"
)

type HealthHandlers struct {
	log            *slog.Logger
	metricsService *metricsservice.MetricService
}

func NewHealth(log *slog.Logger, metricsService *metricsservice.MetricService) HealthHandlers {
	return HealthHandlers{log: log, metricsService: metricsService}
}

func (m *HealthHandlers) ReadinessProbe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responseError(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	ctx := r.Context()
	err := m.metricsService.HealthCheck(ctx)

	if err != nil {
		responseError(w, r, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (m *HealthHandlers) LivenessProbe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responseError(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	w.WriteHeader(http.StatusOK)
}
