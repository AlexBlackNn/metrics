package server

import (
	"fmt"
	"github.com/AlexBlackNn/metrics/cmd/server/router"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/handlers"
	"github.com/AlexBlackNn/metrics/internal/logger"
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	"github.com/AlexBlackNn/metrics/pkg/storage/memstorage"
	"log/slog"
	"net/http"
	"time"
)

// App service consists all entities needed to work
type App struct {
	MetricsService *metricsservice.MetricService
	Handlers       handlers.MetricHandlers
	Cfg            *config.Config
	Log            *slog.Logger
	Srv            *http.Server
}

// New creates App collecting service layer, config, logger and predefined storage layer
func New() (*App, error) {

	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	log := logger.New(cfg.Env)

	// err is now skipped, but when migratings to postgres/sqlite/etc... err will be checked
	memStorage, _ := memstorage.New()

	metricsService := metricsservice.New(
		log,
		cfg,
		memStorage,
	)

	projectHandlers := handlers.New(log, metricsService)

	srv := &http.Server{
		Addr:         fmt.Sprintf(cfg.ServerAddr),
		Handler:      router.NewChiRouter(log, projectHandlers),
		ReadTimeout:  time.Duration(cfg.ServerReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.ServerWriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.ServerIdleTimeout) * time.Second,
	}

	return &App{
		MetricsService: metricsService,
		Handlers:       projectHandlers,
		Srv:            srv,
		Cfg:            cfg,
		Log:            log,
	}, nil
}
