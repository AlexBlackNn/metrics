package app_agent

import (
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/services/agentmetricsservice"
	"log/slog"
)

// App service consists all service layers
type AppHttp struct {
	MetricsService *agentmetricsservice.MetricsHttpService
}

// New create App
func NewAppHttp(
	log *slog.Logger,
	cfg *config.Config,
) *AppHttp {

	// init services
	metricsService := agentmetricsservice.NewMetricsHttpService(
		log,
		cfg,
	)
	return &AppHttp{MetricsService: metricsService}
}
