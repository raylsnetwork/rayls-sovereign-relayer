package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/batcher"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

const CTSTransactionTable = "cts_transaction"

// CTSTransactionRepository persists cts_transaction rows for the
// cts/service/batcher pipeline. Implements batcher.Repository.
type CTSTransactionRepository struct {
	pool *pgxpool.Pool
}

func NewCTSTransactionRepository(pool *pgxpool.Pool) *CTSTransactionRepository {
	return &CTSTransactionRepository{pool: pool}
}

// CTSTransactionModel is the wire-to-DB shape. pgx decodes into this via
// RowToStructByName matching the column names.
type CTSTransactionModel struct {
	CorrelationID   string     `db:"correlation_id"`
	Identity        string     `db:"identity"`
	MessageType     string     `db:"message_type"`
	Address         []byte     `db:"address"`
	Calldata        []byte     `db:"calldata"`
	Status          string     `db:"status"`
	TxHash          []byte     `db:"tx_hash"`
	ReceiptStatus   *int16     `db:"receipt_status"`
	RevertData      []byte     `db:"revert_data"`
	ErrorReason     *string    `db:"error_reason"`
	SendAttempts    int        `db:"send_attempts"`
	ReceiptAttempts int        `db:"receipt_attempts"`
	SentAt          *time.Time `db:"sent_at"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
}

func toCTSTransaction(m CTSTransactionModel) batcher.CTSTransaction {
	t := batcher.CTSTransaction{
		CorrelationID:   m.CorrelationID,
		Identity:        m.Identity,
		MessageType:     m.MessageType,
		Address:         common.BytesToAddress(m.Address),
		Calldata:        m.Calldata,
		Status:          batcher.Status(m.Status),
		SendAttempts:    m.SendAttempts,
		ReceiptAttempts: m.ReceiptAttempts,
		SentAt:          m.SentAt,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		RevertData:      m.RevertData,
	}
	if len(m.TxHash) > 0 {
		t.TxHash = common.BytesToHash(m.TxHash)
	}
	if m.ReceiptStatus != nil {
		t.ReceiptStatus = uint64(*m.ReceiptStatus) //nolint:gosec // status is 0 or 1
	}
	if m.ErrorReason != nil {
		t.ErrorReason = *m.ErrorReason
	}
	return t
}

func (r *CTSTransactionRepository) Insert(ctx context.Context, req types.TxRequest) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (correlation_id, identity, message_type, address, calldata, status)
		VALUES ($1, $2, $3, $4, $5, 'pending')
		ON CONFLICT (correlation_id, message_type) DO NOTHING
	`, CTSTransactionTable)

	_, err := r.pool.Exec(ctx, query,
		req.CorrelationID, req.Identity, req.MessageType, req.Address.Bytes(), req.Calldata,
	)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("inserting cts_transaction: %w", err))
	}
	return nil
}

func (r *CTSTransactionRepository) ClaimPending(
	ctx context.Context, identity string, limit int,
) ([]batcher.CTSTransaction, error) {
	query := fmt.Sprintf(`
		SELECT * FROM %s
		WHERE status = 'pending' AND identity = $1
		ORDER BY created_at
		LIMIT $2
		FOR UPDATE SKIP LOCKED
	`, CTSTransactionTable)
	return r.queryClaim(ctx, query, identity, limit)
}

func (r *CTSTransactionRepository) ClaimSent(
	ctx context.Context, identity string, limit int,
) ([]batcher.CTSTransaction, error) {
	query := fmt.Sprintf(`
		SELECT * FROM %s
		WHERE status = 'sent' AND identity = $1
		ORDER BY sent_at
		LIMIT $2
		FOR UPDATE SKIP LOCKED
	`, CTSTransactionTable)
	return r.queryClaim(ctx, query, identity, limit)
}

func (r *CTSTransactionRepository) ClaimStuck(
	ctx context.Context, identity string, olderThan time.Time, limit int,
) ([]batcher.CTSTransaction, error) {
	query := fmt.Sprintf(`
		SELECT * FROM %s
		WHERE status = 'sent' AND identity = $1 AND sent_at < $2
		ORDER BY sent_at
		LIMIT $3
		FOR UPDATE SKIP LOCKED
	`, CTSTransactionTable)
	return r.queryClaim(ctx, query, identity, olderThan, limit)
}

// queryClaim runs a FOR UPDATE SKIP LOCKED SELECT inside a short-lived
// transaction. The lock exists only to guard the SELECT against a
// hypothetical second writer — it is not held across the subsequent RPC
// batch call.
func (r *CTSTransactionRepository) queryClaim(
	ctx context.Context, query string, args ...any,
) ([]batcher.CTSTransaction, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("starting claim tx: %w", err))
	}
	defer func() { _ = tx.Rollback(ctx) }()

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("querying cts_transaction: %w", err))
	}
	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[CTSTransactionModel])
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("collecting cts_transaction rows: %w", err))
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, withstack.Wrap(fmt.Errorf("committing claim tx: %w", err))
	}

	out := make([]batcher.CTSTransaction, len(models))
	for i, m := range models {
		out[i] = toCTSTransaction(m)
	}
	return out, nil
}

func (r *CTSTransactionRepository) MarkSent(
	ctx context.Context, key batcher.TxKey, hash common.Hash,
) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET status = 'sent', tx_hash = $3, sent_at = NOW(), updated_at = NOW()
		WHERE correlation_id = $1 AND message_type = $2 AND status = 'pending'
	`, CTSTransactionTable)

	_, err := r.pool.Exec(ctx, query, key.CorrelationID, key.MessageType, hash.Bytes())
	if err != nil {
		return withstack.Wrap(fmt.Errorf("marking sent: %w", err))
	}
	return nil
}

func (r *CTSTransactionRepository) MarkResent(
	ctx context.Context, key batcher.TxKey,
) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET sent_at = NOW(), updated_at = NOW()
		WHERE correlation_id = $1 AND message_type = $2 AND status = 'sent'
	`, CTSTransactionTable)

	_, err := r.pool.Exec(ctx, query, key.CorrelationID, key.MessageType)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("marking resent: %w", err))
	}
	return nil
}

func (r *CTSTransactionRepository) MarkFinishedSuccess(
	ctx context.Context, key batcher.TxKey, receiptStatus uint64,
) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET status = 'finished', receipt_status = $3, updated_at = NOW()
		WHERE correlation_id = $1 AND message_type = $2
	`, CTSTransactionTable)

	_, err := r.pool.Exec(ctx, query, key.CorrelationID, key.MessageType, int16(receiptStatus)) //nolint:gosec // 0 or 1
	if err != nil {
		return withstack.Wrap(fmt.Errorf("marking finished success: %w", err))
	}
	return nil
}

func (r *CTSTransactionRepository) MarkFinishedRevert(
	ctx context.Context, key batcher.TxKey, revertData []byte,
) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET status = 'finished', receipt_status = 0, revert_data = $3, updated_at = NOW()
		WHERE correlation_id = $1 AND message_type = $2
	`, CTSTransactionTable)

	_, err := r.pool.Exec(ctx, query, key.CorrelationID, key.MessageType, revertData)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("marking finished revert: %w", err))
	}
	return nil
}

func (r *CTSTransactionRepository) MarkFailed(
	ctx context.Context, key batcher.TxKey, reason string,
) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET status = 'failed', error_reason = $3, updated_at = NOW()
		WHERE correlation_id = $1 AND message_type = $2
	`, CTSTransactionTable)

	_, err := r.pool.Exec(ctx, query, key.CorrelationID, key.MessageType, reason)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("marking failed: %w", err))
	}
	return nil
}

// incrementAttempts is the common body for IncrementSendAttempts /
// IncrementReceiptAttempts. It uses unnest with two parallel arrays so
// the (correlation_id, message_type) pairs are matched positionally,
// not as a cross product.
func (r *CTSTransactionRepository) incrementAttempts(
	ctx context.Context, column string, keys []batcher.TxKey,
) error {
	if len(keys) == 0 {
		return nil
	}
	correlationIDs := make([]string, len(keys))
	messageTypes := make([]string, len(keys))
	for i, k := range keys {
		correlationIDs[i] = k.CorrelationID
		messageTypes[i] = k.MessageType
	}
	query := fmt.Sprintf(`
		UPDATE %s
		SET %s = %s + 1, updated_at = NOW()
		WHERE (correlation_id, message_type) IN (
			SELECT * FROM unnest($1::text[], $2::text[])
		)
	`, CTSTransactionTable, column, column)

	_, err := r.pool.Exec(ctx, query, correlationIDs, messageTypes)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("incrementing %s: %w", column, err))
	}
	return nil
}

func (r *CTSTransactionRepository) IncrementSendAttempts(
	ctx context.Context, keys []batcher.TxKey,
) error {
	return r.incrementAttempts(ctx, "send_attempts", keys)
}

func (r *CTSTransactionRepository) IncrementReceiptAttempts(
	ctx context.Context, keys []batcher.TxKey,
) error {
	return r.incrementAttempts(ctx, "receipt_attempts", keys)
}
