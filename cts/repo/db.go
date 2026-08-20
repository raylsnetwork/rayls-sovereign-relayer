package repo

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

	// Startup ping retry: total budget when the caller does not supply a
	// context. Wraps NewDatabaseConnection() in a 1-minute deadline that
	// absorbs postgres interrupted-shutdown recovery (~65s observed
	// 2026-05-12) without crash-looping the CTS process.
	defaultStartupPingBudget = time.Minute
	pingRetryInitialBackoff  = 250 * time.Millisecond
	pingRetryMaxBackoff      = 4 * time.Second
	pingRetryMultiplier      = 2.0
)

type PostgresDB struct {
	pool         *pgxpool.Pool
	databaseName string
}

// NewDatabaseConnection creates a new PostgreSQL connection pool. The initial
// ping is retried with exponential backoff (default 1-minute budget) to
// absorb transient postgres unavailability during startup or
// interrupted-shutdown recovery. Use NewDatabaseConnectionCtx to control the
// retry budget explicitly.
//
// Migrations must be run separately by calling Migrate().
func NewDatabaseConnection(connectionString string) (*PostgresDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStartupPingBudget)
	defer cancel()
	return NewDatabaseConnectionCtx(ctx, connectionString)
}

// NewDatabaseConnectionCtx is the context-aware variant. Tests inject short
// contexts to exercise the no-DB failure path; production callers can supply
// their own shutdown ctx so a SIGTERM during startup-retry exits cleanly.
func NewDatabaseConnectionCtx(ctx context.Context, connectionString string) (*PostgresDB, error) {
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
// ping succeeds or ctx is cancelled.
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
	db.pool.Close()
}

func (db *PostgresDB) Pool() *pgxpool.Pool {
	return db.pool
}

func (db *PostgresDB) DatabaseName() string {
	return db.databaseName
}

func (db *PostgresDB) NewRaylsViewKeysRepository(tableName string) *RaylsViewKeysRepository {
	return NewRaylsViewKeysRepository(tableName, db.pool)
}

func (db *PostgresDB) NewRaylsSignKeysRepository(tableName string) *RaylsSignKeysRepository {
	return NewRaylsSignKeysRepository(tableName, db.pool)
}

func (db *PostgresDB) NewPaymentSpendKeysRepository(tableName string) *PaymentSpendKeysRepository {
	return NewPaymentSpendKeysRepository(tableName, db.pool)
}

func (db *PostgresDB) NewSharedSecretsRepository(tableName string) *SharedSecretsRepository {
	return NewSharedSecretsRepository(tableName, db.pool)
}

func (db *PostgresDB) NewEnygmaSelfSecretsRepository(tableName string) *EnygmaSelfSecretsRepository {
	return NewEnygmaSelfSecretsRepository(tableName, db.pool)
}

func (db *PostgresDB) NewCTSTransactionRepository() *CTSTransactionRepository {
	return NewCTSTransactionRepository(db.pool)
}
