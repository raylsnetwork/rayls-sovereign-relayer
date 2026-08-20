package testutil

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const postgresConnectTimeout = 30 * time.Second

type DBConfig struct {
	User     string
	Pass     string
	Database string
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
			joinStrings(columns, ", "),
			joinStrings(placeholders, ", "),
		)

		_, err := pool.Exec(ctx, query, values...)
		if err != nil {
			return fmt.Errorf("failed to insert row: %w", err)
		}
	}

	return nil
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

func SetupLogger() {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
}
