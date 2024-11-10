package servergrpc

import (
	"context"
	"log/slog"
	"net"

	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	v1 "github.com/AlexBlackNn/metrics/internal/handlersgrpc/v1"
	"github.com/AlexBlackNn/metrics/internal/logger"
	"github.com/AlexBlackNn/metrics/internal/migrator"
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	"github.com/AlexBlackNn/metrics/pkg/storage/memstorage"
	"github.com/AlexBlackNn/metrics/pkg/storage/postgres"
	"github.com/golang-migrate/migrate/v4"
	"google.golang.org/grpc"
)

type MetricsStorage interface {
	UpdateMetric(
		ctx context.Context,
		metric models.MetricGetter,
	) error
	UpdateSeveralMetrics(
		ctx context.Context,
		metrics map[string]models.MetricGetter,
	) error
	GetMetric(
		ctx context.Context,
		metric models.MetricGetter,
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
	GrpcServer     *grpc.Server
	Cfg            *configserver.Config
	Log            *slog.Logger
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
	if cfg.ServerDataBaseDSN != "" {
		postgresStorage, err := postgres.New(cfg, log)
		if err != nil {
			return nil, err
		}
		log.Info("Starts to apply migrations")
		err = migrator.ApplyMigration(cfg)
		if err != nil {
			if err != migrate.ErrNoChange {
				log.Error("Failed to apply migration", "err", err.Error())
				return nil, err
			}
			log.Info("No migration to apply")
			return RegisterService(postgresStorage, postgresStorage, cfg, log)
		}
		log.Info("Finish to apply migrations")
		return RegisterService(postgresStorage, postgresStorage, cfg, log)
	}
	memStorage, err := memstorage.New(cfg, log)
	if err != nil {
		return nil, err
	}
	return RegisterService(memStorage, memStorage, cfg, log)
}

func RegisterService(ms MetricsStorage, hc HealthChecker, cfg *configserver.Config, log *slog.Logger) (*App, error) {

	prjMetricsService := metricsservice.New(
		log,
		cfg,
		ms,
		hc,
	)

	grpcServer := grpc.NewServer()
	v1.Register(grpcServer, prjMetricsService)

	return &App{
		MetricsService: prjMetricsService,
		GrpcServer:     grpcServer,
		Cfg:            cfg,
		Log:            log,
		DataBase:       ms,
		HealthChecker:  hc,
	}, nil
}

func (a *App) Start() error {
	lis, err := net.Listen("tcp", ":44044")
	if err != nil {
		return err
	}
	a.Log.Info("gRPC server listening on", "addr", ":44044")
	return a.GrpcServer.Serve(lis)
}
