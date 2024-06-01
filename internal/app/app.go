package app

import (
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/services/metrics_service"
	"github.com/AlexBlackNn/metrics/storage/mem_storage"
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
	memStorage, _ := mem_storage.New()

	// init services
	metricsService := metrics_service.New(
		log,
		cfg,
		memStorage,
	)
	return &App{MetricsService: metricsService}
}
