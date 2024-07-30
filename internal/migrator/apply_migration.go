package migrator

import (
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	// migration lib
	"github.com/golang-migrate/migrate/v4"
	// driver for migration applying postgres
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// driver to get migrations from files (*.sql in our case)
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func ApplyMigration(cfg *configserver.Config) error {
	m, err := migrate.New(
		"file://"+"./migrations",
		fmt.Sprintf(cfg.ServerDataBaseDSN),
	)
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil {
		return err
	}
	return nil
}
