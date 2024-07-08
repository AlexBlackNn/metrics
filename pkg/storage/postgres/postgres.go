package postgres

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/pkg/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
	"html/template"
	"log/slog"
	"strings"
)

type PostStorage struct {
	db *sql.DB
}

// Helper function to get the type
func GetType(m models.MetricGetter) string {
	return m.GetType()
}

// Helper function to get the name
func GetName(m models.MetricGetter) string {
	return m.GetName()
}

func New(cfg *configserver.Config, log *slog.Logger) (*PostStorage, error) {
	db, err := sql.Open("pgx", cfg.ServerDataBaseDSN)
	if err != nil {
		log.Error("Unable to connect to database", "error", err)
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
	metric models.MetricGetter,
) (models.MetricGetter, error) {

	tpl := template.Must(template.New("sqlQuery").Funcs(template.FuncMap{
		"GetType": GetType,
	}).Parse(`
      WITH LatestCounter AS (
        SELECT
          MAX(created) AS latest_created
        FROM
          {{GetType .}}_part
        WHERE
          {{GetType .}}_part.name = $1
      )
      SELECT
        t.name, c.name, c.value
      FROM
        {{GetType .}}_part as c
      JOIN
        app.types as t ON c.metric_id = t.uuid
      JOIN
        LatestCounter lc ON c.created = lc.latest_created
      WHERE
        c.name = $1
      ORDER BY c.created;
  `))

	var sqlTmp bytes.Buffer
	err := tpl.Execute(&sqlTmp, metric)
	if err != nil {
		return nil, err
	}

	row := s.db.QueryRowContext(
		ctx,
		sqlTmp.String(),
		strings.ToLower(metric.GetName()),
	)

	if metric.GetType() == "counter" {
		var metricCounter models.Metric[uint64]
		err = row.Scan(&metricCounter.Type, &metricCounter.Name, &metricCounter.Value)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf(
					"DATA LAYER: storage.postgres.GetMetric: %w",
					storage.ErrMetricNotFound,
				)
			}
			return nil, fmt.Errorf(
				"DATA LAYER: storage.postgres.GetMetric: %w",
				err,
			)
		}
		return &metricCounter, nil
	}
	var metricGauge models.Metric[float64]
	err = row.Scan(&metricGauge.Type, &metricGauge.Name, &metricGauge.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(
				"DATA LAYER: storage.postgres.GetMetric: %w",
				storage.ErrMetricNotFound,
			)
		}
		return nil, fmt.Errorf(
			"DATA LAYER: storage.postgres.GetMetric: %w",
			err,
		)
	}
	return &metricGauge, nil
}

func (s *PostStorage) GetAllMetrics(
	ctx context.Context,
) ([]models.MetricGetter, error) {
	return nil, nil
}

func (s *PostStorage) HealthCheck(
	ctx context.Context,
) error {
	err := s.db.PingContext(ctx)
	fmt.Println("11111111111111111", err)
	return err
}
