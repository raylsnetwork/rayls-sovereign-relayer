package repository

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raylsnetwork/rayls-sovereign-relayer/conv"
	"github.com/raylsnetwork/rayls-sovereign-relayer/listener"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

// LastProcessedBlockRepository is responsible for interacting with the last processed block table in PostgreSQL
type LastProcessedBlockRepository struct {
	pool        *pgxpool.Pool
	tableName   string
	updateMutex sync.RWMutex
}

// NewLastProcessedBlockRepository creates a new instance of LastProcessedBlockRepository
func NewLastProcessedBlockRepository(pool *pgxpool.Pool) *LastProcessedBlockRepository {
	return &LastProcessedBlockRepository{
		pool:        pool,
		tableName:   LastProcessedBlockNumberCollectionName,
		updateMutex: sync.RWMutex{},
	}
}

// Get retrieves the last processed ID from PostgreSQL
func (blockRepo *LastProcessedBlockRepository) Get(
	ctx context.Context,
	documentID types.LastProcessedBlockDocument,
) (*big.Int, error) {
	query := fmt.Sprintf(`SELECT chain, last_block, updated_at FROM %s WHERE chain = $1`, blockRepo.tableName)

	rows, err := blockRepo.pool.Query(ctx, query, string(documentID))
	if err != nil {
		return nil, withstack.Wrap(
			fmt.Errorf("error querying last processed block number from DB for %s: %w", string(documentID), err),
		)
	}
	defer rows.Close()

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[LastProcessedBlockNumber])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, withstack.Wrap(listener.ErrLastProcessedBlockNotFound)
		}
		return nil, withstack.Wrap(
			fmt.Errorf("error querying last processed block number from DB for %s: %w", string(documentID), err),
		)
	}

	lastBlockBigInt, err := conv.StringToBigInt(result.LastBlock)
	if err != nil {
		return nil, withstack.Wrap(err)
	}

	return lastBlockBigInt, nil
}

func (r *LastProcessedBlockRepository) Create(
	ctx context.Context,
	chainID types.LastProcessedBlockDocument,
	blockNum *big.Int,
) error {
	_, err := r.Get(ctx, chainID)
	if !errors.Is(err, listener.ErrLastProcessedBlockNotFound) {
		if err != nil {
			return withstack.Wrap(err)
		} else {
			return withstack.Wrap(listener.ErrLastProcessedBlockAlreadyExists)
		}
	}

	query := fmt.Sprintf(`INSERT INTO %s (chain, last_block, updated_at) VALUES ($1, $2, $3)`, r.tableName)

	_, err = r.pool.Exec(ctx, query, string(chainID), blockNum.String(), time.Now())
	if err != nil {
		return withstack.Wrap(err)
	}

	return nil
}

func (r *LastProcessedBlockRepository) Update(
	ctx context.Context,
	chainID types.LastProcessedBlockDocument,
	blockNum *big.Int,
) error {
	query := fmt.Sprintf(`UPDATE %s SET last_block = $1, updated_at = $2 WHERE chain = $3`, r.tableName)

	result, err := r.pool.Exec(ctx, query, blockNum.String(), time.Now(), string(chainID))
	if err != nil {
		return withstack.Wrap(err)
	}
	if result.RowsAffected() == 0 {
		return withstack.Wrap(listener.ErrLastProcessedBlockNotFound)
	}

	return nil
}

// UpdateWithLock forces an update with exclusive lock (for reset operations)
func (r *LastProcessedBlockRepository) UpdateWithLock(
	ctx context.Context,
	chainID types.LastProcessedBlockDocument,
	blockNum *big.Int,
) error {
	r.updateMutex.Lock()
	defer r.updateMutex.Unlock()

	query := fmt.Sprintf(`UPDATE %s SET last_block = $1, updated_at = $2 WHERE chain = $3`, r.tableName)

	result, err := r.pool.Exec(ctx, query, blockNum.String(), time.Now(), string(chainID))
	if err != nil {
		return withstack.Wrap(err)
	}
	if result.RowsAffected() == 0 {
		return withstack.Wrap(listener.ErrLastProcessedBlockNotFound)
	}

	return nil
}

// ResetLastProcessedBlock forces to update the last processed ID in PostgreSQL
func (blockRepo *LastProcessedBlockRepository) ResetLastProcessedBlock(
	ctx context.Context,
	blockNum *big.Int,
	documentId types.LastProcessedBlockDocument,
) error {
	blockRepo.updateMutex.Lock()
	defer blockRepo.updateMutex.Unlock()

	query := fmt.Sprintf(`
		INSERT INTO %s (chain, last_block, updated_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (chain) DO UPDATE SET last_block = $2, updated_at = $3
	`, blockRepo.tableName)

	_, err := blockRepo.pool.Exec(ctx, query, string(documentId), blockNum.String(), time.Now())
	if err != nil {
		return withstack.Wrap(
			fmt.Errorf("error upserting last processed block number into DB for %s: %w", string(documentId), err),
		)
	}
	return nil
}
