package app

import (
	"context"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
)

const migrationsPath = "migrations"

//go:embed  migrations/*.sql
var fs embed.FS

func (a *app) startMigrate(ctx context.Context, migratePath string, dbName string, db *sqlx.DB) error {
	err := db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("db connection not alive: %w", err)
	}
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{
		DatabaseName: dbName,
		SchemaName:   "public",
	})
	if err != nil {
		return fmt.Errorf("db migration database driver error: %w", err)
	}
	source, err := iofs.New(fs, migratePath)
	if err != nil {
		return fmt.Errorf("db migration source driver error: %w", err)
	}
	instance, err := migrate.NewWithInstance("fs", source, dbName, driver)
	if err != nil {
		return fmt.Errorf("db migration instance error: %w", err)
	}
	if err := instance.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("db migration up error: %w", err)
	}

	return nil
}
