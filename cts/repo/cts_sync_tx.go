package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

const CTSSyncTxTable = "cts_sync_tx"

// uniqueViolationCode is the PostgreSQL SQLSTATE for a unique-constraint
// violation. We map it to service.ErrTxAlreadyExists on insert so the caller
// can branch on the sentinel rather than parsing driver errors.
const uniqueViolationCode = "23505"

// CTSSyncTxRepository persists cts_sync_tx rows for the sync SignAndSend
// idempotency path. State transitions live on service.SyncTransaction; this
// repository is a dumb persistence port — it loads and stores the aggregate,
// with the single exception of Create, whose UNIQUE(id) constraint is the
// load-bearing serialization point for dedup.
type CTSSyncTxRepository struct {
	pool *pgxpool.Pool
}

func NewCTSSyncTxRepository(pool *pgxpool.Pool) *CTSSyncTxRepository {
	return &CTSSyncTxRepository{pool: pool}
}

// ctsSyncTxModel is the wire-to-DB shape. pgx decodes into this via
// RowToStructByName matching the column names.
type ctsSyncTxModel struct {
	ID          string    `db:"id"`
	FromAddress []byte    `db:"from_address"`
	TxHash      []byte    `db:"tx_hash"`
	TxRLP       []byte    `db:"tx_rlp"`
	ResultState string    `db:"result_state"`
	ReceiptJSON []byte    `db:"receipt_json"`
	RevertData  []byte    `db:"revert_data"`
	Version     int64     `db:"version"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// stateToColumn maps the domain state to its DB string. Kept here so the
// persistence vocabulary lives at the persistence boundary.
func stateToColumn(s service.SyncState) (string, error) {
	switch s {
	case service.StatePending:
		return "pending", nil
	case service.StateMined:
		return "mined", nil
	case service.StateReverted:
		return "reverted", nil
	case service.StateFailed:
		return "failed", nil
	default:
		return "", fmt.Errorf("unknown sync state %d", uint(s))
	}
}

func stateFromColumn(s string) (service.SyncState, error) {
	switch s {
	case "pending":
		return service.StatePending, nil
	case "mined":
		return service.StateMined, nil
	case "reverted":
		return service.StateReverted, nil
	case "failed":
		return service.StateFailed, nil
	default:
		return 0, fmt.Errorf("unknown result_state %q", s)
	}
}

// toSyncTransaction rebuilds the aggregate from a row: the tx is decoded from
// its canonical (EIP-2718) binary form, the receipt from its stored JSON.
func toSyncTransaction(m ctsSyncTxModel) (*service.SyncTransaction, error) {
	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(m.TxRLP); err != nil {
		return nil, withstack.Wrap(fmt.Errorf("decoding tx_rlp for %s: %w", m.ID, err))
	}

	state, err := stateFromColumn(m.ResultState)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("decoding state for %s: %w", m.ID, err))
	}

	var receipt *types.Receipt
	if len(m.ReceiptJSON) > 0 {
		receipt = new(types.Receipt)
		if err := json.Unmarshal(m.ReceiptJSON, receipt); err != nil {
			return nil, withstack.Wrap(fmt.Errorf("decoding receipt_json for %s: %w", m.ID, err))
		}
	}

	return &service.SyncTransaction{
		ID:         m.ID,
		From:       common.BytesToAddress(m.FromAddress),
		Tx:         tx,
		Receipt:    receipt,
		RevertData: m.RevertData,
		State:      state,
		Version:    m.Version,
	}, nil
}

// Create inserts a freshly-built (pending) row. The UNIQUE(id) PRIMARY KEY is
// the serialization point: a duplicate id returns service.ErrTxAlreadyExists,
// which the caller treats as the signal to enter the recovery flow. On success
// the in-memory Version is reset to the stored default (0).
func (r *CTSSyncTxRepository) Create(ctx context.Context, t *service.SyncTransaction) error {
	rlp, err := t.Tx.MarshalBinary()
	if err != nil {
		return withstack.Wrap(fmt.Errorf("encoding tx %s: %w", t.Tx.Hash(), err))
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (id, from_address, tx_hash, tx_rlp)
		VALUES ($1, $2, $3, $4)
	`, CTSSyncTxTable)

	_, err = r.pool.Exec(ctx, query, t.ID, t.From.Bytes(), t.Tx.Hash().Bytes(), rlp)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return service.ErrTxAlreadyExists
		}
		return withstack.Wrap(fmt.Errorf("inserting cts_sync_tx: %w", err))
	}
	t.Version = 0
	return nil
}

// Get loads the full aggregate for id, including its current State and (if
// terminal) its stored result. Returns service.ErrTxNotFound if no row exists.
func (r *CTSSyncTxRepository) Get(ctx context.Context, id string) (*service.SyncTransaction, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, CTSSyncTxTable)

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("querying cts_sync_tx: %w", err))
	}
	m, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[ctsSyncTxModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, service.ErrTxNotFound
		}
		return nil, withstack.Wrap(fmt.Errorf("collecting cts_sync_tx row: %w", err))
	}
	return toSyncTransaction(m)
}

// Save persists the current in-memory state of an existing row using
// optimistic concurrency: the UPDATE asserts the version the aggregate was
// loaded at and bumps it. If no row matches (the version advanced under a
// concurrent writer, or the row is gone) it returns service.ErrStaleTransition
// and the caller should re-Get and re-evaluate. On success the in-memory
// Version is advanced so the aggregate can be saved again.
func (r *CTSSyncTxRepository) Save(ctx context.Context, t *service.SyncTransaction) error {
	rlp, err := t.Tx.MarshalBinary()
	if err != nil {
		return withstack.Wrap(fmt.Errorf("encoding tx %s: %w", t.Tx.Hash(), err))
	}

	stateCol, err := stateToColumn(t.State)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("saving cts_sync_tx %s: %w", t.ID, err))
	}

	var receiptJSON []byte
	if t.Receipt != nil {
		receiptJSON, err = json.Marshal(t.Receipt)
		if err != nil {
			return withstack.Wrap(fmt.Errorf("encoding receipt for %s: %w", t.ID, err))
		}
	}

	query := fmt.Sprintf(`
		UPDATE %s
		SET from_address = $2,
		    tx_hash      = $3,
		    tx_rlp       = $4,
		    result_state = $5,
		    receipt_json = $6,
		    revert_data  = $7,
		    version      = version + 1,
		    updated_at   = NOW()
		WHERE id = $1 AND version = $8
	`, CTSSyncTxTable)

	tag, err := r.pool.Exec(ctx, query,
		t.ID, t.From.Bytes(), t.Tx.Hash().Bytes(), rlp,
		stateCol, receiptJSON, t.RevertData, t.Version,
	)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("updating cts_sync_tx %s: %w", t.ID, err))
	}
	if tag.RowsAffected() == 0 {
		return service.ErrStaleTransition
	}
	t.Version++
	return nil
}

// DeleteTerminalOlderThan removes hard-terminal rows (mined, reverted) last
// updated before cutoff. pending and failed rows are never deleted regardless
// of age — pending rows may be mid-recovery and failed rows hold the only
// (hash, RLP) needed to ask the chain whether a lost tx landed. Returns the
// number of rows deleted.
func (r *CTSSyncTxRepository) DeleteTerminalOlderThan(ctx context.Context, cutoff time.Time) (int64, error) {
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE result_state IN ('mined', 'reverted') AND updated_at < $1
	`, CTSSyncTxTable)

	tag, err := r.pool.Exec(ctx, query, cutoff)
	if err != nil {
		return 0, withstack.Wrap(fmt.Errorf("deleting terminal cts_sync_tx rows: %w", err))
	}
	return tag.RowsAffected(), nil
}
