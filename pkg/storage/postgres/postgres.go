// Package postgres is a data layer which communicate with postgresql database.
package postgres

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/pkg/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// PostStorage client to communicate with postgresql
type PostStorage struct {
	tmpl Tmpl
	db   *sql.DB
}

// GetType is a helper function to get the type.
func GetType(m models.MetricGetter) string {
	return m.GetType()
}

// New creates client for concurrent use by multiple goroutines and maintains its own pool of idle connections.
func New(cfg *configserver.Config, log *slog.Logger) (*PostStorage, error) {
	db, err := sql.Open("pgx", cfg.ServerDataBaseDSN)
	if err != nil {
		log.Error("Unable to connect to database", "error", err)
		return nil, fmt.Errorf(
			"DATA LAYER: storage.postgres.GetMetric: %w - %v",
			storage.ErrConnectionFailed, err,
		)
	}
	return &PostStorage{db: db, tmpl: NewTemplate()}, nil
}

// Stop closes connection to DB.
func (s *PostStorage) Stop() error {
	return s.db.Close()
}

// UpdateMetric updates metric value.
func (s *PostStorage) UpdateMetric(
	ctx context.Context,
	metric models.MetricGetter,
) error {
	var sqlTmp bytes.Buffer
	err := s.tmpl["UpdateMetric"].Execute(&sqlTmp, metric)
	if err != nil {
		return fmt.Errorf(
			"DATA LAYER: storage.postgres.UpdateMetric: couldn't create template: %w - %v",
			storage.ErrSQLExec, err,
		)
	}

	_, err = s.db.ExecContext(
		ctx, sqlTmp.String(), metric.GetType(), metric.GetName(), metric.GetValue(),
	)
	if err != nil {
		return fmt.Errorf(
			"DATA LAYER: storage.postgres.UpdateMetric: couldn't save metric: %w - %v",
			storage.ErrSQLExec, err,
		)
	}
	return nil
}

// UpdateSeveralMetrics updates several metric value.
func (s *PostStorage) UpdateSeveralMetrics(
	ctx context.Context,
	metrics map[string]models.MetricGetter,
) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf(
			"DATA LAYER: storage.postgres.UpdateSeveralMetrics: couldn't open transaction: %w - %v",
			storage.ErrSQLExec, err,
		)
	}
	defer func(tx *sql.Tx) {
		tx.Rollback()
	}(tx)

	sqlTmpStms := make(map[string]string)
	sqlTmpStms[configserver.MetricTypeGauge] = "INSERT INTO app.gauge_part (metric_id, name, value) VALUES ((SELECT uuid FROM app.types WHERE name = $1), $2, $3)"
	sqlTmpStms[configserver.MetricTypeCounter] = "INSERT INTO app.counter_part (metric_id, name, value) VALUES ((SELECT uuid FROM app.types WHERE name = $1), $2, $3)"

	// This prepared statements seem to be unnecessary, because Exec creates Prepare statement under the hood.
	preparedStmt := make(map[string]*sql.Stmt)
	for name, onesqlTmpStms := range sqlTmpStms {
		stmt, err := tx.PrepareContext(ctx, onesqlTmpStms)
		if err != nil {
			return fmt.Errorf(
				"DATA LAYER: storage.postgres.UpdateSeveralMetrics: couldn't prepare context: %w - %v",
				storage.ErrSQLExec, err,
			)
		}
		preparedStmt[name] = stmt
	}
	// The statements prepared for a transaction by calling the transaction's Tx.Prepare or Tx.Stmt methods
	//are closed by the call to Tx.Commit or Tx.Rollback. https://pkg.go.dev/database/sql#Tx
	// BUT it must be proved in pgx docs or code
	//TODO: prove or refute
	defer func() {
		for _, onesqlTmpStms := range preparedStmt {
			onesqlTmpStms.Close()
		}
	}()

	for _, oneMetric := range metrics {
		_, err = preparedStmt[oneMetric.GetType()].ExecContext(
			ctx, oneMetric.GetType(), oneMetric.GetName(), oneMetric.GetValue(),
		)
		if err != nil {
			return fmt.Errorf(
				"DATA LAYER: storage.postgres.UpdateSeveralMetrics: couldn't save metric: %w - %v",
				storage.ErrSQLExec, err,
			)
		}
	}
	return tx.Commit()
}

// GetMetric updates gets metric value.
func (s *PostStorage) GetMetric(
	ctx context.Context,
	metric models.MetricGetter,
) (models.MetricGetter, error) {

	var sqlTmp bytes.Buffer
	err := s.tmpl["GetMetric"].Execute(&sqlTmp, metric)
	if err != nil {
		return nil, err
	}

	row := s.db.QueryRowContext(
		ctx,
		sqlTmp.String(),
		metric.GetName(),
	)

	var tmpMetric TempMetric
	err = row.Scan(&tmpMetric.Type, &tmpMetric.Name, &tmpMetric.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(
				"DATA LAYER: storage.postgres.GetMetric: %w - %v",
				storage.ErrMetricNotFound, err,
			)
		}
		return nil, fmt.Errorf(
			"DATA LAYER: storage.postgres.GetMetric: %w - %v",
			storage.ErrSQLExec, err,
		)
	}
	metricDB, err := models.New(
		tmpMetric.GetType(),
		tmpMetric.GetName(),
		tmpMetric.GetStringValue(),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"DATA LAYER: storage.postgres.GetMetric: models.New %w - %v",
			storage.ErrUnexpectedBehavior, err,
		)
	}
	return metricDB, nil
}

// GetAllMetrics gets all metric values.
func (s *PostStorage) GetAllMetrics(
	ctx context.Context,
) ([]models.MetricGetter, error) {

	var metrics []models.MetricGetter

	var sqlTmp bytes.Buffer
	err := s.tmpl["GetAllMetrics"].Execute(&sqlTmp, nil)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryContext(
		ctx,
		sqlTmp.String(),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"DATA LAYER: storage.postgres.GetAllMetrics: %w - %v",
			storage.ErrSQLExec, err,
		)
	}
	defer func(rows *sql.Rows) {
		err = errors.Join(err, rows.Close())
	}(rows)

	for rows.Next() {
		var tmpMetric TempMetric
		err = rows.Scan(&tmpMetric.Type, &tmpMetric.Name, &tmpMetric.Value)
		if err != nil {
			return nil, fmt.Errorf(
				"DATA LAYER: storage.postgres.GetAllMetrics: rows.Scan %w - %v",
				storage.ErrUnexpectedBehavior, err,
			)
		}

		metricDB, err := models.New(
			tmpMetric.GetType(),
			tmpMetric.GetName(),
			tmpMetric.GetStringValue(),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"DATA LAYER: storage.postgres.GetAllMetrics: models.New %w - %v",
				storage.ErrUnexpectedBehavior, err,
			)
		}
		metrics = append(metrics, metricDB)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf(
			"DATA LAYER: storage.postgres.GetAllMetrics: rows.Err %w - %v",
			storage.ErrUnexpectedBehavior, err,
		)
	}
	return metrics, nil
}

// HealthCheck provides service health check.
func (s *PostStorage) HealthCheck(
	ctx context.Context,
) error {
	err := s.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf(
			"DATA LAYER: storage.postgres.HealthCheck: %w - %v",
			storage.ErrConnectionFailed, err,
		)
	}
	return nil
}
