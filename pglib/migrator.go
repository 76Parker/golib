package pglib

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	postgresmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func ApplyMigrations(pool *pgxpool.Pool, pathToMigrationFiles, dbName string) error {
	if pool == nil {
		return errors.New("pgx pool is nil")
	}
	if pathToMigrationFiles == "" {
		pathToMigrationFiles = "./migrations"
	}
	if dbName == "" {
		dbName = "postgres"
	}

	sqlDB := stdlib.OpenDBFromPool(pool)

	driver, err := postgresmigrate.WithInstance(sqlDB, &postgresmigrate.Config{
		DatabaseName: dbName,
	})
	if err != nil {
		return fmt.Errorf("create postgres migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+pathToMigrationFiles,
		dbName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}
	defer m.Close()
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}

		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
