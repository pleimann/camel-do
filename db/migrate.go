package db

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFs embed.FS

func Migrate(database *sql.DB) error {
	config := &sqlite3.Config{}
	driver, driverErr := sqlite3.WithInstance(database, config)
	if driverErr != nil {
		return fmt.Errorf("cannot open database migration driver: %w", driverErr)
	}
	source, err := iofs.New(migrationsFs, "migrations")
	if err != nil {
		return fmt.Errorf("error creating migration filesystem: %w", err)
	}

	defer source.Close()

	migrator, err := migrate.NewWithInstance("iofs", source, "", driver)
	if err != nil {
		return fmt.Errorf("error creating migration instance: %w", err)
	}

	if err := migrator.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("cannot perform database migration: %w", err)
		}
	}

	return nil
}
