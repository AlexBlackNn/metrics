package v2

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
)

type HealthHandlers struct {
	log            *slog.Logger
	metricsService *metricsservice.MetricService
}

func NewHealth(log *slog.Logger, metricsService *metricsservice.MetricService) HealthHandlers {
	return HealthHandlers{log: log, metricsService: metricsService}
}

// ReadinessProbe проверка готовности приложения.
// @ID infoHealth
// @Summary Проверка готовности приложения
// @Description Определяет можно ли подавать трафик на сервис
// @Tags Health
// @Produce json
// @Success 200 {object} v2.Response
// @Failure 500 {object} v2.Response
// @Router /ping [get]
func (m *HealthHandlers) ReadinessProbe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responseError(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	ctx, cancel := context.WithTimeoutCause(r.Context(), 300*time.Millisecond, errors.New("readinessProbe timeout"))
	defer cancel()

	err := m.metricsService.HealthCheck(ctx)

	if err != nil {
		responseError(w, r, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	responseHealth(w, r, http.StatusOK, "ready")
}

func (m *HealthHandlers) LivenessProbe(w http.ResponseWriter, r *http.Request) {
	fmt.Println("11111111111")
	if r.Method != http.MethodGet {
		responseError(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	fmt.Println("2222222222")
	responseHealth(w, r, http.StatusOK, "alive")
}
