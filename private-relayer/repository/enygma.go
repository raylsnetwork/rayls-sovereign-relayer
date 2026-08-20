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

type EnygmaRepository struct {
	pool      *pgxpool.Pool
	tableName string
}

func NewEnygmaRepository(pool *pgxpool.Pool) *EnygmaRepository {
	return &EnygmaRepository{
		pool:      pool,
		tableName: EnygmaCollectionName,
	}
}

func (repo *EnygmaRepository) CreateEnygma(ctx context.Context, enygma types.Enygma) error {
	model := enygmaTypeToModel(enygma)
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()

	query := fmt.Sprintf(`
		INSERT INTO %s (
			resource_id, finalized_r, finalized_balance, finalized_block_number,
			pending_block_number, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, repo.tableName)

	q := GetQuerier(ctx, repo.pool)
	_, err := q.Exec(ctx, query,
		model.ResourceId, model.FinalizedR, model.FinalizedBalance,
		model.FinalizedBlockNumber, model.PendingBlockNumber,
		model.CreatedAt, model.UpdatedAt,
	)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("PostgreSQL CreateEnygma: %w", err))
	}
	return nil
}

func (repo *EnygmaRepository) UpdateEnygma(
	ctx context.Context,
	resourceId string,
	finalizedBalance *big.Int,
	finalizedR *big.Int,
	finalizedBlockNumber *big.Int,
	pendingBlockNumber *big.Int,
) error {
	query := fmt.Sprintf(`
		UPDATE %s SET
			finalized_balance = $1,
			finalized_r = $2,
			finalized_block_number = $3,
			pending_block_number = $4,
			updated_at = $5
		WHERE resource_id = $6
	`, repo.tableName)

	q := GetQuerier(ctx, repo.pool)
	_, err := q.Exec(ctx, query,
		finalizedBalance.String(),
		finalizedR.String(),
		finalizedBlockNumber.Uint64(),
		pendingBlockNumber.Uint64(),
		time.Now(),
		resourceId,
	)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("PostgreSQL UpdateEnygma: %w", err))
	}
	return nil
}

func (repo *EnygmaRepository) GetEnygmaByResourceId(ctx context.Context, resourceId string) (types.Enygma, error) {
	query := fmt.Sprintf(`
		SELECT resource_id, finalized_r, finalized_balance, finalized_block_number,
		       pending_block_number, created_at, updated_at
		FROM %s
		WHERE resource_id = $1
	`, repo.tableName)

	q := GetQuerier(ctx, repo.pool)
	rows, err := q.Query(ctx, query, resourceId)
	if err != nil {
		return types.Enygma{}, withstack.Wrap(fmt.Errorf("failed to query enygma by resource ID: %w", err))
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Enygma])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.Enygma{}, withstack.Wrap(fmt.Errorf("no enygma found with resourceId: %s", resourceId))
		}
		return types.Enygma{}, withstack.Wrap(fmt.Errorf("failed to decode enygma by resource ID: %w", err))
	}

	return enygmaModelToType(model)
}

func (repo *EnygmaRepository) GetEnygmaByResourceIds(
	ctx context.Context,
	resourceIds []string,
) ([]types.Enygma, error) {
	if len(resourceIds) == 0 {
		return []types.Enygma{}, nil
	}

	query := fmt.Sprintf(`
		SELECT resource_id, finalized_r, finalized_balance, finalized_block_number,
		       pending_block_number, created_at, updated_at
		FROM %s
		WHERE resource_id = ANY($1)
	`, repo.tableName)

	q := GetQuerier(ctx, repo.pool)
	rows, err := q.Query(ctx, query, resourceIds)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("PostgreSQL GetEnygmaByResourceIds: %w", err))
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[Enygma])
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("PostgreSQL GetEnygmaByResourceIds cursor: %w", err))
	}

	var results []types.Enygma
	for _, model := range models {
		enygma, err := enygmaModelToType(model)
		if err != nil {
			return nil, withstack.Wrap(err)
		}

		results = append(results, enygma)
	}

	return results, nil
}

func enygmaTypeToModel(enygma types.Enygma) Enygma {
	return Enygma{
		ResourceId:           enygma.ResourceId,
		FinalizedR:           enygma.FinalizedR.String(),
		FinalizedBalance:     enygma.FinalizedBalance.String(),
		FinalizedBlockNumber: enygma.FinalizedBlockNumber.Uint64(),
		PendingBlockNumber:   enygma.PendingBlockNumber.Uint64(),
	}
}

func enygmaModelToType(model Enygma) (types.Enygma, error) {
	rFactor, ok := new(big.Int).SetString(model.FinalizedR, 10)
	if !ok {
		return types.Enygma{}, withstack.Wrap(
			fmt.Errorf("failed to parse big.Int values from DB for resourceId: %s", model.ResourceId),
		)
	}

	balance, ok := new(big.Int).SetString(model.FinalizedBalance, 10)
	if !ok {
		return types.Enygma{}, withstack.Wrap(
			fmt.Errorf("failed to parse big.Int values from DB for resourceId: %s", model.ResourceId),
		)
	}

	finalizedBlockNumber := new(big.Int).SetUint64(model.FinalizedBlockNumber)
	pendingBlockNumber := new(big.Int).SetUint64(model.PendingBlockNumber)

	return types.Enygma{
		ResourceId:           model.ResourceId,
		FinalizedR:           rFactor,
		FinalizedBalance:     balance,
		FinalizedBlockNumber: finalizedBlockNumber,
		PendingBlockNumber:   pendingBlockNumber,
	}, nil
}
