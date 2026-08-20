package testtools

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	postgresConnectTimeout = 30 * time.Second
)

type DBConfig struct {
	User           string
	Pass           string
	Database       string
	MigrationsPath string // Optional: file:// URL to migrations directory
}

func SetupPostgres(t testing.TB, config DBConfig) (*pgxpool.Pool, func()) {
	ctx := context.Background()

	if config.User == "" {
		config.User = "test"
	}
	if config.Pass == "" {
		config.Pass = "test"
	}
	if config.Database == "" {
		config.Database = "testdb"
	}

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(config.Database),
		postgres.WithUsername(config.User),
		postgres.WithPassword(config.Pass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(postgresConnectTimeout),
		),
	)
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		t.Fatalf("failed to parse connection string: %s", err)
	}

	poolConfig.MinConns = 1
	poolConfig.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		t.Fatalf("failed to connect to postgres: %s", err)
	}

	// Run migrations if path is provided
	if config.MigrationsPath != "" {
		if err := runMigrations(pool, config.MigrationsPath); err != nil {
			pool.Close()
			t.Fatalf("failed to run migrations: %s", err)
		}
	}

	return pool, func() {
		pool.Close()
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}
}

func SeedTable(pool *pgxpool.Pool, tableName string, rows []map[string]any) error {
	if len(rows) == 0 {
		return nil
	}

	ctx := context.Background()

	for _, row := range rows {
		if len(row) == 0 {
			continue
		}

		columns := make([]string, 0, len(row))
		placeholders := make([]string, 0, len(row))
		values := make([]any, 0, len(row))

		i := 1
		for col, val := range row {
			columns = append(columns, col)
			placeholders = append(placeholders, fmt.Sprintf("$%d", i))
			values = append(values, val)
			i++
		}

		query := fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES (%s)",
			tableName,
			strings.Join(columns, ", "),
			strings.Join(placeholders, ", "),
		)

		_, err := pool.Exec(ctx, query, values...)
		if err != nil {
			return fmt.Errorf("failed to insert row: %w", err)
		}
	}

	return nil
}

// runMigrations runs database migrations using the pool-based driver.
func runMigrations(pool *pgxpool.Pool, migrationsPath string) error {
	// Convert pool to *sql.DB (pgx migrate driver requires *sql.DB)
	sqlDB := stdlib.OpenDBFromPool(pool)

	driver, err := pgxmigrate.WithInstance(sqlDB, &pgxmigrate.Config{})
	if err != nil {
		return fmt.Errorf("create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("create migration instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
