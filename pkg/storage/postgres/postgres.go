package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log/slog"
)

type PostStorage struct {
	db *sql.DB
}

func New(cfg *configserver.Config, log *slog.Logger) (*PostStorage, error) {
	db, err := sql.Open("pgx", cfg.ServerDataBaseDSN)
	if err != nil {
		return nil, fmt.Errorf(
			"DATA LAYER: storage.postgres.New: couldn't open a database: %w",
			err,
		)
	}
	return &PostStorage{db: db}, nil
}

func (s *PostStorage) Stop() error {
	return s.db.Close()
}

func (s *PostStorage) UpdateMetric(
	ctx context.Context,
	metric models.MetricGetter,
) error {
	return nil
}

func (s *PostStorage) GetMetric(
	ctx context.Context,
	metricName string,
) (models.MetricGetter, error) {
	return nil, nil
}

func (s *PostStorage) GetAllMetrics(
	ctx context.Context,
) ([]models.MetricGetter, error) {
	return nil, nil
}

func (s *PostStorage) HealthCheck(
	ctx context.Context,
) error {
	return s.db.PingContext(ctx)
}
