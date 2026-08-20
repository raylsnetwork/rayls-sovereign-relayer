package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

type TransactionManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{pool: pool}
}

func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("failed to begin transaction: %w", err))
	}
	defer func() { _ = tx.Rollback(ctx) }() // Rollback is no-op after Commit

	// Store tx in context so repositories can access it via GetQuerier
	txCtx := ContextWithTx(ctx, tx)

	err = fn(txCtx)
	if err != nil {
		return fmt.Errorf("executing transactional callback: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return withstack.Wrap(fmt.Errorf("failed to commit transaction: %w", err))
	}
	return nil
}
