package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log/slog"
	"os"
	"text/template"
)

type PostStorage struct {
	db *sql.DB
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


	tpl := template.Must(template.New("sqlQuery").Parse(`
     WITH LatestCounter AS (
       SELECT
         MAX(created) AS latest_created
       FROM
         {{.TableName}}
       WHERE
         {{.TableName}}.name = '{{.Name}}'
     )
     SELECT
       t.name, c.name, c.value
     FROM
       {{.TableName}} as c
     JOIN
       app.types as t ON c.metric_id = t.uuid
     JOIN
       LatestCounter lc ON c.created = lc.latest_created
     WHERE
       c.name = '{{.Name}}'
     ORDER BY c.created;
 `))

	data := struct {
		TableName string
		Name      string
	}{
		TableName: metric.GetType()+"s",
		Name:      metric.GetName(),
	}

	// Execute the template and print the result
	err := tpl.Execute(os.Stdout, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}
}

	var row *sql.Row
	switch metric.GetType() {
	case "counter":
		query := "SELECT id, email, pass_hash, is_admin FROM users WHERE (id = $1);"
		row = s.db.QueryRowContext(ctx, query, sqlParam)
	case "gauge":
		query := "SELECT id, email, pass_hash, is_admin FROM users WHERE (email = $1);"
		row = s.db.QueryRowContext(ctx, query, sqlParam)
	default:
		return nil, errors.New("wrong metric type")
	}

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.PassHash, &user.IsAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf(
				"DATA LAYER: storage.postgres.GetUser: %w",
				storage.ErrUserNotFound,
			)
		}
		return models.User{}, fmt.Errorf(
			"DATA LAYER: storage.postgres.GetUser: %w",
			err,
		)
	}
	return user, nil
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
