package repository

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-migrate/migrate/v4"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/raylsnetwork/rayls-sovereign-relayer/backoff"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

//go:embed migrations/*.sql
var migrations embed.FS

const (
	defaultMaxConns        = 25
	defaultMinConns        = 2
	defaultMaxConnLifetime = 1 * time.Hour
	defaultMaxConnIdleTime = 5 * time.Minute
	defaultHealthCheckTime = 30 * time.Second
	defaultConnectTimeout  = 10 * time.Second

	// Startup ping retry parameters. The loop runs until ctx cancels
	// (typically a 1-minute deadline at the run.go callsite). Tuned to
	// absorb the postgres interrupted-shutdown recovery window observed
	// on 2026-05-12 (~65s).
	pingRetryInitialBackoff = 250 * time.Millisecond
	pingRetryMaxBackoff     = 4 * time.Second
	pingRetryMultiplier     = 2.0
)

type PostgresDB struct {
	pool         *pgxpool.Pool
	databaseName string
}

// Connect creates a new PostgreSQL connection pool. The initial ping is
// retried with exponential backoff until the context is cancelled or the
// server responds — this absorbs transient unavailability during postgres
// startup, interrupted-shutdown recovery, and brief network blips so the
// relayer does not crash-loop while waiting for its database.
//
// Migrations must be run separately by calling Migrate().
func Connect(ctx context.Context, connectionString string) (*PostgresDB, error) {
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("parse connection string: %w", err))
	}

	databaseName := config.ConnConfig.Database

	config.MaxConns = defaultMaxConns
	config.MinConns = defaultMinConns
	config.MaxConnLifetime = defaultMaxConnLifetime
	config.MaxConnIdleTime = defaultMaxConnIdleTime
	config.HealthCheckPeriod = defaultHealthCheckTime
	config.ConnConfig.ConnectTimeout = defaultConnectTimeout

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("create connection pool: %w", err))
	}

	if err := pingWithBackoff(ctx, pool, databaseName); err != nil {
		pool.Close()
		return nil, err
	}

	db := &PostgresDB{
		pool:         pool,
		databaseName: databaseName,
	}

	return db, nil
}

// pingWithBackoff retries pool.Ping with exponential backoff until either the
// ping succeeds or ctx is cancelled. Delays come from backoff.Exponential.Next()
// (the shared utility used throughout this repo); the loop itself is open-coded
// because backoff.Do() requires a positive maxAttempts and we want to retry
// until ctx terminates the wait — there is no good upper bound for how long
// postgres recovery can legitimately take.
//
// The returned error wraps the last ping failure with the elapsed retry budget
// so an operator can distinguish "DB never came up" from "DB rejected us".
func pingWithBackoff(ctx context.Context, pool *pgxpool.Pool, databaseName string) error {
	strategy, err := backoff.NewExponential(pingRetryInitialBackoff, pingRetryMultiplier, pingRetryMaxBackoff)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("invalid postgres ping backoff config: %w", err))
	}

	attempt := 0
	start := time.Now()
	for {
		attempt++
		pingErr := pool.Ping(ctx)
		if pingErr == nil {
			if attempt > 1 {
				slog.Info(
					"Postgres reachable after retries",
					slog.String("database", databaseName),
					slog.Int("attempts", attempt),
					slog.Duration("elapsed", time.Since(start)),
				)
			}
			return nil
		}

		// Context cancelled by caller? Stop and report budget context.
		if ctxErr := ctx.Err(); ctxErr != nil {
			return withstack.Wrap(fmt.Errorf(
				"ping database %q failed after %d attempts (%s) [ctx: %s]: %w",
				databaseName, attempt, time.Since(start).Truncate(time.Millisecond), ctxErr, pingErr,
			))
		}

		delay := strategy.Next(attempt)
		slog.Warn(
			"Postgres ping failed; retrying",
			slog.String("database", databaseName),
			slog.Int("attempt", attempt),
			slog.Duration("backoff", delay),
			slog.Any("error", pingErr),
		)

		// Sleep, but wake immediately if the context is cancelled mid-backoff.
		t := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			t.Stop()
			return withstack.Wrap(fmt.Errorf(
				"ping database %q cancelled after %d attempts (%s): %w",
				databaseName, attempt, time.Since(start).Truncate(time.Millisecond), ctx.Err(),
			))
		case <-t.C:
		}
	}
}

// Migrate runs embedded database migrations.
func (db *PostgresDB) Migrate() error {
	// Create source driver from embedded filesystem
	sourceDriver, err := iofs.New(migrations, "migrations")
	if err != nil {
		return withstack.Wrap(fmt.Errorf("create migration source: %w", err))
	}

	// Convert pool to *sql.DB (pgx migrate driver requires *sql.DB)
	sqlDB := stdlib.OpenDBFromPool(db.pool)

	dbDriver, err := pgxmigrate.WithInstance(sqlDB, &pgxmigrate.Config{})
	if err != nil {
		return withstack.Wrap(fmt.Errorf("create migration driver: %w", err))
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, db.databaseName, dbDriver)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("create migration instance: %w", err))
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return withstack.Wrap(fmt.Errorf("run migrations: %w", err))
	}

	return nil
}

func (db *PostgresDB) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

func (db *PostgresDB) Pool() *pgxpool.Pool {
	return db.pool
}

func (db *PostgresDB) DatabaseName() string {
	return db.databaseName
}
