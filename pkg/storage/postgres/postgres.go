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

	var tmpMetric TempMetric
	err = row.Scan(&tmpMetric.Type, &tmpMetric.Name, &tmpMetric.Value)
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
	metricDB, err := models.New(
		tmpMetric.GetType(),
		tmpMetric.GetName(),
		tmpMetric.GetStringValue(),
	)
	if err != nil {
		return nil, err
	}
	return metricDB, nil
}

func (s *PostStorage) GetAllMetrics(
	ctx context.Context,
) ([]models.MetricGetter, error) {

	var metrics []models.MetricGetter

	sqlTmp := `
	WITH LatestCounter AS (
		SELECT
			MAX(created) AS latest_created, name
		FROM
			app.counter_part
		GROUP BY
			name
	), LatestGauge AS (
		SELECT
			MAX(created) AS latest_created, name
		FROM
			app.gauge_part
		GROUP BY
			name
	)
	SELECT
		t.name, c.name, c.value
	FROM
		app.counter_part as c
			JOIN
		app.types as t ON c.metric_id = t.uuid
			JOIN
		LatestCounter lc ON c.created = lc.latest_created AND c.name = lc.name
	UNION
	SELECT
		t.name, g.name, g.value
	FROM
		app.gauge_part as g
			JOIN
		app.types as t ON g.metric_id = t.uuid
			JOIN
		LatestGauge lc ON g.created = lc.latest_created AND g.name = lc.name
`

	rows, err := s.db.QueryContext(
		ctx,
		sqlTmp,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tmpMetric TempMetric
		err = rows.Scan(&tmpMetric.Type, &tmpMetric.Name, &tmpMetric.Value)
		if err != nil {
			return nil, err
		}

		metricDB, err := models.New(
			tmpMetric.GetType(),
			tmpMetric.GetName(),
			tmpMetric.GetStringValue(),
		)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metricDB)
	}
	return metrics, nil
}

func (s *PostStorage) HealthCheck(
	ctx context.Context,
) error {
	return s.db.PingContext(ctx)
}
