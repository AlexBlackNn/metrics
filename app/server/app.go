package server

import (
	"context"
	"fmt"
	"github.com/AlexBlackNn/metrics/cmd/server/router"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/internal/handlers/v1"
	v2 "github.com/AlexBlackNn/metrics/internal/handlers/v2"
	"github.com/AlexBlackNn/metrics/internal/logger"
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	"github.com/AlexBlackNn/metrics/pkg/storage/memstorage"
	"github.com/AlexBlackNn/metrics/pkg/storage/postgres"
	"log/slog"
	"net/http"
	"time"
)

type MetricsStorage interface {
	UpdateMetric(
		ctx context.Context,
		metric models.MetricGetter,
	) error
	GetMetric(
		ctx context.Context,
		metricName string,
	) (models.MetricGetter, error)
	GetAllMetrics(
		ctx context.Context,
	) ([]models.MetricGetter, error)
}

type HealthChecker interface {
	HealthCheck(
		ctx context.Context,
	) error
}

// App service consists all entities needed to work.
type App struct {
	MetricsService *metricsservice.MetricService
	HandlersV1     v1.MetricHandlers
	HandlersV2     v2.MetricHandlers
	Cfg            *configserver.Config
	Log            *slog.Logger
	Srv            *http.Server
	DataBase       MetricsStorage
	HealthChecker  HealthChecker
}

// New creates App collecting service layer, config, logger and predefined storage layer.
func New() (*App, error) {

	cfg, err := configserver.New()
	if err != nil {
		return nil, err
	}

	log := logger.New(cfg.Env)

	// Err is now skipped, but when migratings to postgres/sqlite/etc... err will be checked.
	memStorage, _ := memstorage.New(cfg, log)
	postgresStorage, err := postgres.New(cfg, log)
	if err != nil {
		return nil, err
	}
	return NewAppInitStorage(memStorage, postgresStorage, cfg, log)
}

func NewAppInitStorage(ms MetricsStorage, hc HealthChecker, cfg *configserver.Config, log *slog.Logger) (*App, error) {

	// Err is now skipped, but when migratings to postgres/sqlite/etc... err will be checked.
	memStorage := ms
	postgresStorage := hc

	metricsService := metricsservice.New(
		log,
		cfg,
		memStorage,
		postgresStorage,
	)

	projectHandlersV1 := v1.New(log, metricsService)
	projectHandlersV2 := v2.New(log, metricsService)

	srv := &http.Server{
		Addr:         fmt.Sprintf(cfg.ServerAddr),
		Handler:      router.NewChiRouter(cfg, log, projectHandlersV1, projectHandlersV2),
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
		DataBase:       memStorage,
		HealthChecker:  postgresStorage,
	}, nil
}
