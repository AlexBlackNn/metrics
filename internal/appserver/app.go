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

// New create App
func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {

	//init storage
	// err is now skipped, but when migrating to postgres/sqlite/etc... err will be checked
	memStorage, _ := memstorage.New()

	// init services
	metricsService := metricsservice.New(
		log,
		cfg,
		memStorage,
	)
	return &App{MetricsService: metricsService}
}
