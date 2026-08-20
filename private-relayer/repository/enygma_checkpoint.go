package repository

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

type EnygmaCheckpointRepository struct {
	pool      *pgxpool.Pool
	tableName string
}

func NewEnygmaCheckpointRepository(pool *pgxpool.Pool) *EnygmaCheckpointRepository {
	return &EnygmaCheckpointRepository{
		pool:      pool,
		tableName: EnygmaCheckpointCollectionName,
	}
}

func (repo *EnygmaCheckpointRepository) CreateEnygmaCheckpoint(
	ctx context.Context,
	checkpoint types.EnygmaCheckpoint,
) error {
	model := enygmaCheckpointTypeToModel(checkpoint)
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()

	query := fmt.Sprintf(`
		INSERT INTO %s (
			resource_id, finalized_public_balance_x, finalized_public_balance_y,
			finalized_block_number, pending_block_number, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, repo.tableName)

	q := GetQuerier(ctx, repo.pool)
	_, err := q.Exec(ctx, query,
		model.ResourceId, model.FinalizedPublicBalanceX, model.FinalizedPublicBalanceY,
		model.FinalizedBlockNumber, model.PendingBlockNumber, model.Status,
		model.CreatedAt, model.UpdatedAt,
	)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("PostgreSQL CreateEnygmaCheckpoint: %w", err))
	}
	return nil
}

// GetLatestCheckpointByFilters returns the latest checkpoint matching the given filters.
// When finalizedBlockNumber is set without pendingBlockNumber, it is used as an upper bound (<).
// When both are set, both are matched exactly.
func (repo *EnygmaCheckpointRepository) GetLatestCheckpointByFilters(
	ctx context.Context,
	resourceId string,
	status *types.EnygmaCheckpointStatus,
	finalizedBlockNumber *big.Int,
	pendingBlockNumber *big.Int,
) (*types.EnygmaCheckpoint, error) {
	query := fmt.Sprintf(`
		SELECT id, resource_id, finalized_public_balance_x, finalized_public_balance_y,
		       finalized_block_number, pending_block_number, status, created_at, updated_at
		FROM %s
		WHERE 1=1
	`, repo.tableName)
	args := []any{}
	argIndex := 1

	if resourceId != "" {
		query += fmt.Sprintf(" AND resource_id = $%d", argIndex)
		args = append(args, resourceId)
		argIndex++
	}

	if status != nil {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, uint8(*status))
		argIndex++
	}

	if finalizedBlockNumber != nil {
		if pendingBlockNumber != nil {
			query += fmt.Sprintf(" AND finalized_block_number = $%d", argIndex)
			args = append(args, finalizedBlockNumber.Uint64())
			argIndex++
			query += fmt.Sprintf(" AND pending_block_number = $%d", argIndex)
			args = append(args, pendingBlockNumber.Uint64())
			argIndex++
		} else {
			query += fmt.Sprintf(" AND finalized_block_number < $%d", argIndex)
			args = append(args, finalizedBlockNumber.Uint64())
		}
	}

	query += " ORDER BY finalized_block_number DESC LIMIT 1"

	q := GetQuerier(ctx, repo.pool)
	rows, err := q.Query(ctx, query, args...)
	if err != nil {
		return nil, withstack.Wrap(err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[EnygmaCheckpoint])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil //nolint:nilnil // returning nil,nil is intentional for "not found" pattern
		}
		return nil, withstack.Wrap(err)
	}

	checkpoint, err := enygmaCheckpointModelToType(model)
	if err != nil {
		return nil, withstack.Wrap(err)
	}

	return &checkpoint, nil
}

func (repo *EnygmaCheckpointRepository) GetValidationCandidates(ctx context.Context) ([]types.EnygmaCheckpoint, error) {
	// Use DISTINCT ON to get the first (lowest finalized_block_number) tentative checkpoint per resource_id
	query := fmt.Sprintf(`
		SELECT DISTINCT ON (resource_id)
		       id, resource_id, finalized_public_balance_x, finalized_public_balance_y,
		       finalized_block_number, pending_block_number, status, created_at, updated_at
		FROM %s
		WHERE status = $1
		ORDER BY resource_id, finalized_block_number ASC
	`, repo.tableName)

	q := GetQuerier(ctx, repo.pool)
	rows, err := q.Query(ctx, query, uint8(types.EnygmaCheckpointStatusTentative))
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("PostgreSQL GetValidationCandidates: %w", err))
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[EnygmaCheckpoint])
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("PostgreSQL GetValidationCandidates decode: %w", err))
	}

	checkpoints := make([]types.EnygmaCheckpoint, len(models))
	for i, model := range models {
		checkpoint, err := enygmaCheckpointModelToType(model)
		if err != nil {
			return nil, withstack.Wrap(fmt.Errorf("PostgreSQL GetValidationCandidates convert: %w", err))
		}
		checkpoints[i] = checkpoint
	}

	return checkpoints, nil
}

func (repo *EnygmaCheckpointRepository) MarkAsFinalized(
	ctx context.Context,
	resourceId string,
	blockNumber *big.Int,
) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET status = $1, updated_at = $2
		WHERE resource_id = $3 AND finalized_block_number = $4
	`, repo.tableName)

	q := GetQuerier(ctx, repo.pool)
	result, err := q.Exec(ctx, query,
		uint8(types.EnygmaCheckpointStatusFinal), time.Now(),
		resourceId, blockNumber.Uint64(),
	)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("PostgreSQL MarkAsFinalized: %w", err))
	}
	if result.RowsAffected() == 0 {
		return withstack.Wrap(fmt.Errorf("PostgreSQL MarkAsFinalized: no checkpoint found"))
	}

	return nil
}

func enygmaCheckpointTypeToModel(checkpoint types.EnygmaCheckpoint) EnygmaCheckpoint {
	return EnygmaCheckpoint{
		ResourceId:              checkpoint.ResourceId,
		FinalizedPublicBalanceX: checkpoint.FinalizedPublicBalanceX.String(),
		FinalizedPublicBalanceY: checkpoint.FinalizedPublicBalanceY.String(),
		FinalizedBlockNumber:    checkpoint.FinalizedBlockNumber.Uint64(),
		PendingBlockNumber:      checkpoint.PendingBlockNumber.Uint64(),
		Status:                  uint8(checkpoint.Status),
	}
}

func enygmaCheckpointModelToType(model EnygmaCheckpoint) (types.EnygmaCheckpoint, error) {
	finalizedPublicBalanceX, ok := new(big.Int).SetString(model.FinalizedPublicBalanceX, 10)
	if !ok {
		return types.EnygmaCheckpoint{}, withstack.Wrap(
			fmt.Errorf("failed to parse FinalizedBalance for resourceId: %s", model.ResourceId),
		)
	}

	finalizedPublicBalanceY, ok := new(big.Int).SetString(model.FinalizedPublicBalanceY, 10)
	if !ok {
		return types.EnygmaCheckpoint{}, withstack.Wrap(
			fmt.Errorf("failed to parse FinalizedBalance for resourceId: %s", model.ResourceId),
		)
	}

	finalizedBlockNumber := new(big.Int).SetUint64(model.FinalizedBlockNumber)
	pendingBlockNumber := new(big.Int).SetUint64(model.PendingBlockNumber)

	return types.EnygmaCheckpoint{
		ID:                      model.ID,
		ResourceId:              model.ResourceId,
		FinalizedPublicBalanceX: finalizedPublicBalanceX,
		FinalizedPublicBalanceY: finalizedPublicBalanceY,
		FinalizedBlockNumber:    finalizedBlockNumber,
		PendingBlockNumber:      pendingBlockNumber,
		Status:                  types.EnygmaCheckpointStatus(model.Status),
	}, nil
}
