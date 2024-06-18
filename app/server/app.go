package server

import (
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/logger"
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	"github.com/AlexBlackNn/metrics/pkg/storage/memstorage"
	"log/slog"
)

// App service consists all entities needed to work
type App struct {
	MetricsService *metricsservice.Monitor
	Cfg            *config.Config
	Log            *slog.Logger
}

// New creates App collecting service layer, config, logger and predefined storage layer
func New() (*App, error) {

	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	log := logger.New(cfg.Env)
	log.Info("starting application", slog.String("cfg", cfg.String()))

	// err is now skipped, but when migratings to postgres/sqlite/etc... err will be checked
	memStorage, _ := memstorage.New()

	metricsService := metricsservice.New(
		log,
		cfg,
		memStorage,
	)

	return &App{MetricsService: metricsService, Cfg: cfg, Log: log}, nil
}
