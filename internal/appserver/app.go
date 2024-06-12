package appserver

import (
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	"github.com/AlexBlackNn/metrics/storage/memstorage"
	"log/slog"
)

// App service consists all services needed to work
type App struct {
	MetricsService *metricsservice.MetricService
}

// New creates App collecting service layer with predefine storage layer
func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {

	// err is now skipped, but when migrating to postgres/sqlite/etc... err will be checked
	memStorage, _ := memstorage.New()

	metricsService := metricsservice.New(
		log,
		cfg,
		memStorage,
	)
	return &App{MetricsService: metricsService}
}
