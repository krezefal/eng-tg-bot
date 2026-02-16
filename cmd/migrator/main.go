package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"

	migfs "github.com/krezefal/eng-tg-bot/migrations"
	"github.com/krezefal/eng-tg-bot/pkg/log"
)

const serviceName = "migrator"

const (
	flagUpName   = "up"
	flagDownName = "down"
	envDBDSN     = "DB_DSN"
)

var logger = log.For(serviceName)

func main() {
	up := flag.Bool(flagUpName, false, "apply all pending migrations")
	down := flag.Bool(flagDownName, false, "rollback all applied migrations")
	flag.Parse()

	logger.Info().Msg("running migrator")
	if err := run(*up, *down); err != nil {
		logger.Fatal().Err(err).Msg("migrator run error")
	}
}

func run(up, down bool) error {
	if up == down {
		return fmt.Errorf("set exactly one flag: --%s or --%s", flagUpName, flagDownName)
	}

	if err := gotenv.Load(); err != nil {
		return fmt.Errorf("load .env: %w", err)
	}

	dsn := strings.TrimSpace(os.Getenv(envDBDSN))
	if dsn == "" {
		return fmt.Errorf("env var %s is empty", envDBDSN)
	}

	m, err := newMigrator(dsn)
	if err != nil {
		return err
	}
	defer closeMigrator(m)

	if up {
		if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migrate up: %w", err)
		}
		logger.Info().Msg("migrations applied")

		return nil
	}

	if err = m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate down: %w", err)
	}
	logger.Info().Msg("migrations rolled back")

	return nil
}

func newMigrator(dsn string) (*migrate.Migrate, error) {
	src, err := iofs.New(migfs.FS, ".")
	if err != nil {
		return nil, fmt.Errorf("create iofs source: %w", err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("postgres.WithInstance: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate.NewWithInstance: %w", err)
	}

	return m, nil
}

func closeMigrator(m *migrate.Migrate) {
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		logger.Warn().Err(srcErr).Msg("migrator close source error")
	}
	if dbErr != nil {
		logger.Warn().Err(dbErr).Msg("migrator close db error")
	}
}
