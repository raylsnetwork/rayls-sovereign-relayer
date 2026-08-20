package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

const lockExpirationTime = 5 * time.Minute

// ErrLockAlreadyExists is returned by InsertNewLock when a non-expired lock for the
// given resource already exists. Callers should use errors.Is to check for this condition.
var ErrLockAlreadyExists = errors.New("resource lock already exists and has not expired")

type ResourceLockRepository struct {
	pool      *pgxpool.Pool
	tableName string
}

func NewResourceLockRepository(pool *pgxpool.Pool) *ResourceLockRepository {
	return &ResourceLockRepository{
		pool:      pool,
		tableName: ResourceLockCollectionName,
	}
}

// StartCleanupRoutine starts a goroutine that periodically removes expired locks.
// This provides automatic expiration of stale locks.
func (repo *ResourceLockRepository) StartCleanupRoutine(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := repo.cleanupExpiredLocks(ctx); err != nil {
					slog.Warn("failed to cleanup expired locks", "error", err)
				}
			}
		}
	}()
}

func (repo *ResourceLockRepository) cleanupExpiredLocks(ctx context.Context) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE expires_at < $1`, repo.tableName)
	_, err := repo.pool.Exec(ctx, query, time.Now().UTC())
	if err != nil {
		return withstack.Wrap(fmt.Errorf("PostgreSQL cleanupExpiredLocks: %w", err))
	}
	return nil
}

func (repo *ResourceLockRepository) InsertNewLock(ctx context.Context, resourceLock string) error {
	now := time.Now().UTC()
	expiresAt := now.Add(lockExpirationTime).Truncate(time.Millisecond)

	// Try to insert or update only if lock doesn't exist or is expired
	query := fmt.Sprintf(`
		INSERT INTO %s (resource_id, expires_at)
		VALUES ($1, $2)
		ON CONFLICT (resource_id) DO UPDATE
		SET expires_at = EXCLUDED.expires_at
		WHERE %s.expires_at < $3
	`, repo.tableName, repo.tableName)

	result, err := repo.pool.Exec(ctx, query, resourceLock, expiresAt, now)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return withstack.Wrap(ErrLockAlreadyExists)
		}
		return withstack.Wrap(fmt.Errorf("PostgreSQL InsertNewLock: %w", err))
	}

	// If no rows were affected, the lock exists and hasn't expired
	if result.RowsAffected() == 0 {
		return withstack.Wrap(ErrLockAlreadyExists)
	}

	return nil
}

func (repo *ResourceLockRepository) RemoveLock(ctx context.Context, resourceId string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE resource_id = $1`, repo.tableName)

	_, err := repo.pool.Exec(ctx, query, resourceId)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("PostgreSQL RemoveLock: %w", err))
	}
	return nil
}

func (repo *ResourceLockRepository) GetLock(ctx context.Context, resourceId string) (types.ResourceLock, error) {
	query := fmt.Sprintf(`SELECT resource_id, expires_at FROM %s WHERE resource_id = $1`, repo.tableName)

	rows, err := repo.pool.Query(ctx, query, resourceId)
	if err != nil {
		return types.ResourceLock{}, withstack.Wrap(fmt.Errorf("PostgreSQL GetLock: %w", err))
	}
	defer rows.Close()

	lock, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[ResourceLock])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.ResourceLock{}, nil
		}
		return types.ResourceLock{}, withstack.Wrap(fmt.Errorf("PostgreSQL GetLock: %w", err))
	}

	return types.ResourceLock{
		ResourceId: lock.ResourceId,
		ExpiresAt:  lock.ExpiresAt,
	}, nil
}

func (repo *ResourceLockRepository) RemoveAllLocks(ctx context.Context) error {
	_, err := repo.pool.Exec(ctx, fmt.Sprintf("DELETE FROM %s", repo.tableName))
	if err != nil {
		return withstack.Wrap(fmt.Errorf("PostgreSQL RemoveAllLocks: %w", err))
	}
	return nil
}
