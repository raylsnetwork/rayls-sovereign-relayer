package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Querier is implemented by both *pgxpool.Pool and pgx.Tx, allowing repository
// methods to work transparently with or without an active transaction.
type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type txKey struct{}

// ContextWithTx stores a transaction in context for use by repository methods.
func ContextWithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// TxFromContext retrieves a transaction from context, returns nil if not present.
func TxFromContext(ctx context.Context) pgx.Tx {
	tx, _ := ctx.Value(txKey{}).(pgx.Tx)
	return tx
}

// GetQuerier returns the transaction from context if present, otherwise returns the pool.
// This allows repository methods to participate in transactions transparently.
func GetQuerier(ctx context.Context, pool *pgxpool.Pool) Querier {
	if tx := TxFromContext(ctx); tx != nil {
		return tx
	}
	return pool
}
