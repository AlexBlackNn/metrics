package app_agent

import (
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/services/agentmetrics"
	"log/slog"
)

// App service consists all service layers
type App struct {
	MetricsService *agentmetrics.MetricsService
}

// New create App
func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {

	// init services
	metricsService := agentmetrics.New(
		log,
		cfg,
	)
	return &App{MetricsService: metricsService}
}
