package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

type MerkleTreeRepository struct {
	pool      *pgxpool.Pool
	tableName string
}

func NewMerkleTreeRepository(pool *pgxpool.Pool) *MerkleTreeRepository {
	return &MerkleTreeRepository{
		pool:      pool,
		tableName: MerkleTreeCollectionName,
	}
}

// CreateMerkleTree inserts a new Merkle tree document into PostgreSQL.
func (r *MerkleTreeRepository) CreateMerkleTree(ctx context.Context, tree *types.MerkleTree) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (type, token_address, number, depth, leaves, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	_, err := q.Exec(ctx, query,
		tree.Type, tree.TokenAddress, tree.Number, tree.Depth, tree.Leaves, time.Now(),
	)
	if err != nil {
		return withstack.Wrap(
			fmt.Errorf("error creating merkle tree %v for token_address %s: %w", tree.Number, tree.TokenAddress, err),
		)
	}
	return nil
}

// Get the latest tree for a token address
func (r *MerkleTreeRepository) GetLatestByTokenAddress(ctx context.Context, tokenAddress string) (*types.MerkleTree, error) {
	query := fmt.Sprintf(`
		SELECT type, token_address, number, depth, leaves, created_at
		FROM %s
		WHERE token_address = $1
		ORDER BY number DESC
		LIMIT 1
	`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	rows, err := q.Query(ctx, query, tokenAddress)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("error getting merkle tree for token_address %s: %w", tokenAddress, err))
	}
	defer rows.Close()

	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[MerkleTree])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil //nolint:nilnil // returning nil,nil is intentional for "not found" pattern
		}
		return nil, withstack.Wrap(fmt.Errorf("error getting merkle tree for token_address %s: %w", tokenAddress, err))
	}

	return convertTreeModelToDomain(record)
}

// Get a tree by its number and tokenAddress
func (r *MerkleTreeRepository) GetByNumberAndTokenAddress(
	ctx context.Context,
	treeNumber int,
	tokenAddress string,
) (*types.MerkleTree, error) {
	query := fmt.Sprintf(`
		SELECT type, token_address, number, depth, leaves, created_at
		FROM %s
		WHERE token_address = $1 AND number = $2
	`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	rows, err := q.Query(ctx, query, tokenAddress, treeNumber)
	if err != nil {
		return nil, withstack.Wrap(
			fmt.Errorf("error getting merkle tree for number %v token_address %s: %w", treeNumber, tokenAddress, err),
		)
	}
	defer rows.Close()

	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[MerkleTree])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil //nolint:nilnil // returning nil,nil is intentional for "not found" pattern
		}
		return nil, withstack.Wrap(
			fmt.Errorf("error getting merkle tree for number %v token_address %s: %w", treeNumber, tokenAddress, err),
		)
	}

	return convertTreeModelToDomain(record)
}

// InsertLeaves inserts leaves into a Merkle tree in PostgreSQL.
// Preserves existing leaf order and appends only unique new leaves.
func (r *MerkleTreeRepository) InsertLeaves(ctx context.Context, tokenAddress string, treeNumber int, leaves []string) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET leaves = leaves || (
			SELECT COALESCE(array_agg(new_leaf), ARRAY[]::TEXT[])
			FROM unnest($1::TEXT[]) AS new_leaf
			WHERE new_leaf != ALL(leaves)
		)
		WHERE token_address = $2 AND number = $3
	`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	result, err := q.Exec(ctx, query, leaves, tokenAddress, treeNumber)
	if err != nil {
		return withstack.Wrap(
			fmt.Errorf(
				"error inserting leaves into merkle tree %v token_address %s: %w",
				treeNumber,
				tokenAddress,
				err,
			),
		)
	}

	if result.RowsAffected() == 0 {
		return withstack.Wrap(
			fmt.Errorf("no merkle tree found for tokenAddress=%s, treeNumber=%d", tokenAddress, treeNumber),
		)
	}

	return nil
}

// DeleteMerkleTree removes all Merkle trees (for testing/cleanup)
func (r *MerkleTreeRepository) DeleteMerkleTree(ctx context.Context) error {
	q := GetQuerier(ctx, r.pool)
	_, err := q.Exec(ctx, fmt.Sprintf("DELETE FROM %s", r.tableName))
	if err != nil {
		return withstack.Wrap(fmt.Errorf("error deleting merkle tree: %w", err))
	}
	return nil
}

// convertTreeModelToDomain maps the PostgreSQL MerkleTree model to the domain-level MerkleTree.
func convertTreeModelToDomain(record MerkleTree) (*types.MerkleTree, error) {
	return &types.MerkleTree{
		Type:         record.Type,
		TokenAddress: record.TokenAddress,
		Number:       record.Number,
		Depth:        record.Depth,
		Leaves:       record.Leaves,
		CreatedAt:    record.CreatedAt,
	}, nil
}
