package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

// ErrRecoveryAlreadyExists is returned when an Insert call hits a duplicate-key conflict,
// meaning the recovery record has already been persisted.
var ErrRecoveryAlreadyExists = errors.New("tx recovery data already exists")

type TxRecoveryDataRepository struct {
	pool      *pgxpool.Pool
	tableName string
}

func NewTxRecoveryDataRepository(pool *pgxpool.Pool) *TxRecoveryDataRepository {
	return &TxRecoveryDataRepository{
		pool:      pool,
		tableName: TxRecoveryDataCollectionName,
	}
}

const txRecoveryDataColumns = `private_hub_tx_hash, resource_id, private_hub_block_number, from_chain_id,
	tx_bytes, event_type, tx_nature, status, created_at`

// Insert persists a new recovery record. Returns ErrRecoveryAlreadyExists on duplicate key.
func (repo *TxRecoveryDataRepository) Insert(ctx context.Context, data types.TxRecoveryData) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (%s)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, repo.tableName, txRecoveryDataColumns)

	q := GetQuerier(ctx, repo.pool)
	_, err := q.Exec(ctx, query,
		data.PrivateHubTxHash,
		data.ResourceID,
		data.PrivateHubBlockNumber,
		data.FromChainID,
		data.TxBytes,
		uint8(data.EventType),
		string(data.TxNature),
		int(data.Status),
		time.Now(),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrRecoveryAlreadyExists
		}
		return withstack.Wrap(fmt.Errorf("PostgreSQL Insert tx_recovery_data: %w", err))
	}
	return nil
}

// GetByPrivateHubTxHash returns the recovery record by private_hub_tx_hash regardless of status.
// Returns (nil, nil) if no record exists.
func (repo *TxRecoveryDataRepository) GetByPrivateHubTxHash(ctx context.Context, privateHubTxHash string) (*types.TxRecoveryData, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM %s
		WHERE private_hub_tx_hash = $1
	`, txRecoveryDataColumns, repo.tableName)

	q := GetQuerier(ctx, repo.pool)
	rows, err := q.Query(ctx, query, privateHubTxHash)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("PostgreSQL GetByPrivateHubTxHash: %w", err))
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[TxRecoveryData])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil //nolint:nilnil // returning nil,nil is intentional for "not found" pattern
		}
		return nil, withstack.Wrap(fmt.Errorf("PostgreSQL GetByPrivateHubTxHash: %w", err))
	}

	return txRecoveryModelToType(model), nil
}

// GetPendingByEventKey returns a pending recovery record matching the compound key
// (resource_id, private_hub_block_number, from_chain_id, event_type).
// Returns (nil, nil) if no pending record exists.
func (repo *TxRecoveryDataRepository) GetPendingByEventKey(
	ctx context.Context,
	resourceID string,
	blockNumber uint64,
	fromChainID string,
	eventType types.EnygmaEventType,
) (*types.TxRecoveryData, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM %s
		WHERE resource_id = $1
		  AND private_hub_block_number = $2
		  AND from_chain_id = $3
		  AND event_type = $4
		  AND status = $5
	`, txRecoveryDataColumns, repo.tableName)

	q := GetQuerier(ctx, repo.pool)
	rows, err := q.Query(ctx, query, resourceID, blockNumber, fromChainID, uint8(eventType), int(types.HistoryStatusPending))
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("PostgreSQL GetPendingByEventKey: %w", err))
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[TxRecoveryData])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil //nolint:nilnil // returning nil,nil is intentional for "not found" pattern
		}
		return nil, withstack.Wrap(fmt.Errorf("PostgreSQL GetPendingByEventKey: %w", err))
	}

	return txRecoveryModelToType(model), nil
}

// MarkConfirmed sets status=confirmed and clears tx_bytes for the given private_hub_tx_hash.
func (repo *TxRecoveryDataRepository) MarkConfirmed(ctx context.Context, privateHubTxHash string) error {
	query := fmt.Sprintf(`
		UPDATE %s SET status = $1, tx_bytes = NULL WHERE private_hub_tx_hash = $2
	`, repo.tableName)

	q := GetQuerier(ctx, repo.pool)
	result, err := q.Exec(ctx, query, int(types.HistoryStatusConfirmed), privateHubTxHash)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("PostgreSQL MarkConfirmed: %w", err))
	}
	if result.RowsAffected() == 0 {
		return withstack.Wrap(fmt.Errorf("PostgreSQL MarkConfirmed: no record found for private_hub_tx_hash %s", privateHubTxHash))
	}
	return nil
}

func txRecoveryModelToType(model TxRecoveryData) *types.TxRecoveryData {
	return &types.TxRecoveryData{
		ResourceID:            model.ResourceID,
		PrivateHubBlockNumber: model.PrivateHubBlockNumber,
		FromChainID:           model.FromChainID,
		PrivateHubTxHash:      model.PrivateHubTxHash,
		TxBytes:               model.TxBytes,
		EventType:             types.EnygmaEventType(model.EventType),
		TxNature:              types.TxNature(model.TxNature),
		Status:                types.HistoryStatus(model.Status),
	}
}
