package appagent

import (
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/services/agentmetricsservice"
	"log/slog"
)

// AppHTTP service consists all service layers
type AppHTTP struct {
	MetricsService *agentmetricsservice.MetricsHTTPService
}

// NewAppHTTP create App
func NewAppHTTP(
	log *slog.Logger,
	cfg *config.Config,
) *AppHTTP {

	// init services
	metricsService := agentmetricsservice.NewMetricsHTTPService(
		log,
		cfg,
	)
	return &AppHTTP{MetricsService: metricsService}
}
