package server

import (
	"fmt"
	"github.com/AlexBlackNn/metrics/cmd/server/router"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/handlers/v1"
	v2 "github.com/AlexBlackNn/metrics/internal/handlers/v2"
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
	HandlersV1     v1.MetricHandlers
	HandlersV2     v2.MetricHandlers
	Cfg            *configserver.Config
	Log            *slog.Logger
	Srv            *http.Server
}

// New creates App collecting service layer, config, logger and predefined storage layer
func New() (*App, error) {

	cfg, err := configserver.New()
	if err != nil {
		return nil, err
	}

	log := logger.New(cfg.Env)

	// err is now skipped, but when migratings to postgres/sqlite/etc... err will be checked
	memStorage, _ := memstorage.New(cfg)

	metricsService := metricsservice.New(
		log,
		cfg,
		memStorage,
	)

	projectHandlersV1 := v1.New(log, metricsService)
	projectHandlersV2 := v2.New(log, metricsService)

	srv := &http.Server{
		Addr:         fmt.Sprintf(cfg.ServerAddr),
		Handler:      router.NewChiRouter(log, projectHandlersV1, projectHandlersV2),
		ReadTimeout:  time.Duration(cfg.ServerReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.ServerWriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.ServerIdleTimeout) * time.Second,
	}

	return &App{
		MetricsService: metricsService,
		HandlersV1:     projectHandlersV1,
		HandlersV2:     projectHandlersV2,
		Srv:            srv,
		Cfg:            cfg,
		Log:            log,
	}, nil
}
