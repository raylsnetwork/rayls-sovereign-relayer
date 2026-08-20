// Decommissioning Teleport (vanilla, atomic).

package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raylsnetwork/rayls-sovereign-relayer/public-relayer/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

type MessageRecordRepository struct {
	pool      *pgxpool.Pool
	tableName string
}

func NewMessageRecordRepository(
	tableName string,
	pool *pgxpool.Pool,
) *MessageRecordRepository {
	return &MessageRecordRepository{
		pool:      pool,
		tableName: tableName,
	}
}

// BatchCreate inserts lifecycle rows for a batch of forward messages.
// Uses ON CONFLICT DO NOTHING so NATS redeliveries don't overwrite existing
// rows (the first publish wins; later deliveries are idempotent).
func (r *MessageRecordRepository) BatchCreate(ctx context.Context, records []service.MessageRecord) error {
	if len(records) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("begin transaction: %w", err))
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := fmt.Sprintf(`
		INSERT INTO %s (id, status, forward_hash, revert_hash, error, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO NOTHING
	`, r.tableName)

	batch := &pgx.Batch{}
	for _, m := range records {
		row := messageRecordToModel(m)
		batch.Queue(query,
			row.ID, row.Status,
			row.ForwardHash, row.RevertHash, row.Error,
			row.CreatedAt, row.UpdatedAt,
		)
	}

	br := tx.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return withstack.Wrap(fmt.Errorf("insert message records: %w", err))
	}

	if err := tx.Commit(ctx); err != nil {
		return withstack.Wrap(fmt.Errorf("commit transaction: %w", err))
	}
	return nil
}

// UpdateForwardResults stamps the forward outcome (status, hash, error) for
// each id. Rows are matched by primary key; missing rows are silently
// skipped (UPDATE ... WHERE id = $1 is a no-op on missing rows).
func (r *MessageRecordRepository) UpdateForwardResults(ctx context.Context, updates []service.ForwardResultUpdate) error {
	if len(updates) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("begin transaction: %w", err))
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := fmt.Sprintf(`
		UPDATE %s
		SET status = $2, forward_hash = $3, error = $4, updated_at = $5
		WHERE id = $1
	`, r.tableName)

	now := time.Now().UTC()
	batch := &pgx.Batch{}
	for _, u := range updates {
		batch.Queue(query, u.ID, int(u.Status), hashToString(u.Hash), u.Error, now)
	}

	br := tx.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return withstack.Wrap(fmt.Errorf("update forward results: %w", err))
	}

	if err := tx.Commit(ctx); err != nil {
		return withstack.Wrap(fmt.Errorf("commit transaction: %w", err))
	}
	return nil
}

// UpdateRevertResults mirrors UpdateForwardResults for the revert leg.
func (r *MessageRecordRepository) UpdateRevertResults(ctx context.Context, updates []service.RevertResultUpdate) error {
	if len(updates) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("begin transaction: %w", err))
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := fmt.Sprintf(`
		UPDATE %s
		SET status = $2, revert_hash = $3, error = $4, updated_at = $5
		WHERE id = $1
	`, r.tableName)

	now := time.Now().UTC()
	batch := &pgx.Batch{}
	for _, u := range updates {
		batch.Queue(query, u.ID, int(u.Status), hashToString(u.Hash), u.Error, now)
	}

	br := tx.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return withstack.Wrap(fmt.Errorf("update revert results: %w", err))
	}

	if err := tx.Commit(ctx); err != nil {
		return withstack.Wrap(fmt.Errorf("commit transaction: %w", err))
	}
	return nil
}

func messageRecordToModel(m service.MessageRecord) MessageRecord {
	return MessageRecord{
		ID:          m.ID,
		Status:      int(m.Status),
		ForwardHash: hashToString(m.ForwardHash),
		RevertHash:  hashToString(m.RevertHash),
		Error:       m.Error,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// hashToString returns an empty string for the zero hash so NOT NULL
// columns with a DEFAULT '' accept it cleanly.
func hashToString(h common.Hash) string {
	if h == (common.Hash{}) {
		return ""
	}
	return h.Hex()
}
