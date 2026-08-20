// Decommissioning Teleport (vanilla, atomic).

package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

var (
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	ErrStartingSessionForCreatingBatchAtomicStatuses = errors.New(
		"error starting session for creating batch atomic statuses",
	)
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	ErrInsertingBatchAtomicStatuses = errors.New("error inserting batch atomic statuses")
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	ErrCreatingBatchAtomicStatuses = errors.New("error creating batch atomic statuses")
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	ErrFindingUnprocessedBatch = errors.New("error finding unprocessed batch atomic statuses")
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	ErrGettingUnprocessedBatch = errors.New("error getting unprocessed batch atomic statuses")
)

// Deprecated: Decommissioning Teleport (vanilla, atomic).
const AtomicStatusCollectionName = "atomic_status"

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type AtomicStatusRepository struct {
	pool      *pgxpool.Pool
	tableName string
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func NewAtomicStatusRepository(pool *pgxpool.Pool) *AtomicStatusRepository {
	return &AtomicStatusRepository{
		pool:      pool,
		tableName: AtomicStatusCollectionName,
	}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (a *AtomicStatusRepository) BatchCreate(ctx context.Context, sums []types.AtomicStatusUpdateMessage) error {
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("%w: %w", ErrCreatingBatchAtomicStatuses, err))
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// ON CONFLICT DO NOTHING makes SUM creation idempotent. The logparser projects Hub
	// AtomicMessageStatusChangedBatch events into this table under at-least-once delivery
	// (blocks get redelivered on any ack failure), and a message reaches exactly one terminal
	// status (Executed XOR Reverted), so first-writer-wins is correct. Without this, a
	// redelivered block — or a duplicate status event — hits atomic_status_pkey and wedges the
	// whole block in an endless "skip ack -> redeliver" loop, stalling all downstream delivery.
	query := fmt.Sprintf(`INSERT INTO %s (shared_id, status, is_processed) VALUES ($1, $2, $3) ON CONFLICT (shared_id) DO NOTHING`, a.tableName)

	batch := &pgx.Batch{}
	for _, sum := range sums {
		sumModel := sumTypeToSumModel(sum, false)
		batch.Queue(query, sumModel.SharedId, sumModel.Status, sumModel.IsProcessed)
	}

	br := tx.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return withstack.Wrap(fmt.Errorf("%w: %w", ErrCreatingBatchAtomicStatuses, err))
	}

	if err := tx.Commit(ctx); err != nil {
		return withstack.Wrap(fmt.Errorf("%w: %w", ErrCreatingBatchAtomicStatuses, err))
	}

	return nil
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (a *AtomicStatusRepository) GetUnprocessedBySharedIDs(
	ctx context.Context,
	sharedIDs []string,
	opts ...Option,
) ([]types.AtomicStatusUpdateMessage, error) {
	queryOptions := GetQueryOptions(opts...)

	query := fmt.Sprintf(`
		SELECT shared_id, status, is_processed
		FROM %s
		WHERE shared_id = ANY($1) AND is_processed = false
	`, a.tableName)

	args := []any{sharedIDs}
	if queryOptions.Limit > 0 {
		query += " LIMIT $2"
		args = append(args, queryOptions.Limit)
	}

	rows, err := a.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("%w: %w", ErrFindingUnprocessedBatch, err))
	}
	defer rows.Close()

	sumsModel, err := pgx.CollectRows(rows, pgx.RowToStructByName[AtomicSUM])
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("%w: %w", ErrGettingUnprocessedBatch, err))
	}

	sums := make([]types.AtomicStatusUpdateMessage, len(sumsModel))
	for i, sumModel := range sumsModel {
		sums[i] = sumModelToSumType(sumModel)
	}

	return sums, nil
}

func sumModelToSumType(sumModel AtomicSUM) (sumType types.AtomicStatusUpdateMessage) {
	sumType = types.AtomicStatusUpdateMessage{
		SharedID: sumModel.SharedId,
		Status:   types.AtomicStatus(sumModel.Status),
	}
	return
}

func sumTypeToSumModel(sumType types.AtomicStatusUpdateMessage, isProcessed bool) (sumModel AtomicSUM) {
	sumModel = AtomicSUM{
		SharedId:    sumType.SharedID,
		Status:      uint8(sumType.Status),
		IsProcessed: isProcessed,
	}
	return
}
