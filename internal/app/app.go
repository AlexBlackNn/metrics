package app

import (
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/services/metrics_service"
	"log/slog"
)

type App struct {
	MetricsService *metrics_service.MetricService
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {

	//init storage

	// init services
	metricsService := metrics_service.New(
		log,
		cfg,
	)
	return &App{MetricsService: metricsService}
}
