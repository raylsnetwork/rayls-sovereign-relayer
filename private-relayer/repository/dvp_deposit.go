package repository

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

type DvpDepositRepository struct {
	pool      *pgxpool.Pool
	tableName string
}

func NewDvpDepositRepository(pool *pgxpool.Pool) *DvpDepositRepository {
	return &DvpDepositRepository{
		pool:      pool,
		tableName: DvpDepositsCollectionName,
	}
}

func (r *DvpDepositRepository) CreateDeposit(ctx context.Context, deposit *types.DvpDeposit) error {
	nullifierStr := ""
	if deposit.Nullifier != nil {
		nullifierStr = deposit.Nullifier.String()
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (
			user_address, salt, token_amount, token_address,
			token_type, token_id, tree_number, commitment,
			nullifier, status, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	_, err := q.Exec(ctx, query,
		deposit.UserAddress,
		deposit.Salt.String(),
		deposit.TokenAmount.String(),
		deposit.TokenAddress,
		deposit.TokenType,
		deposit.TokenID,
		int(deposit.TreeNumber),
		deposit.Commitment.String(),
		nullifierStr,
		deposit.Status,
		time.Now(),
	)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("failed to insert deposit: %w", err))
	}
	return nil
}

// Confirm the deposit by setting the tree number and status to ready
func (r *DvpDepositRepository) ConfirmDeposit(ctx context.Context, commitment *big.Int, treeNumber *big.Int) error {
	query := fmt.Sprintf(`UPDATE %s SET tree_number = $1, status = $2 WHERE commitment = $3`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	_, err := q.Exec(ctx, query,
		int(treeNumber.Int64()), types.DvpDepositUnspent, commitment.String(),
	)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("failed to confirm deposit: %w", err))
	}
	return nil
}

// Update the deposit status to used
func (r *DvpDepositRepository) UpdateDepositStatus(ctx context.Context, commitment *big.Int, status types.DvpDepositStatus) error {
	query := fmt.Sprintf(`UPDATE %s SET status = $1 WHERE commitment = $2`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	_, err := q.Exec(ctx, query, int(status), commitment.String())
	if err != nil {
		return withstack.Wrap(fmt.Errorf("failed to update deposit status: %w", err))
	}

	return nil
}

func (r *DvpDepositRepository) BatchUpdateStatusForCommitments(
	ctx context.Context,
	commitments []string,
	status types.DvpDepositStatus,
) error {
	query := fmt.Sprintf(`UPDATE %s SET status = $1 WHERE commitment = ANY($2)`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	_, err := q.Exec(ctx, query, status, commitments)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("failed to batch update deposit statuses: %w", err))
	}
	return nil
}

func (r *DvpDepositRepository) GetDepositByCommitment(ctx context.Context, commitment *big.Int) (*types.DvpDeposit, error) {
	query := fmt.Sprintf(`
		SELECT user_address, salt, token_amount, token_address,
		       token_type, token_id, tree_number, commitment,
		       nullifier, status, created_at
		FROM %s
		WHERE commitment = $1
	`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	rows, err := q.Query(ctx, query, commitment.String())
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to query deposit by commitment: %w", err))
	}
	defer rows.Close()

	deposit, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[DvpDeposit])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil //nolint:nilnil // returning nil,nil is intentional for "not found" pattern
		}
		return nil, withstack.Wrap(fmt.Errorf("failed to decode deposit by commitment: %w", err))
	}

	return convertModelToDomain(deposit)
}

func (r *DvpDepositRepository) GetDepositByCommitmentAndStatus(
	ctx context.Context,
	commitment *big.Int,
	status types.DvpDepositStatus,
) (*types.DvpDeposit, error) {
	query := fmt.Sprintf(`
		SELECT user_address, salt, token_amount, token_address,
		       token_type, token_id, tree_number, commitment,
		       nullifier, status, created_at
		FROM %s
		WHERE commitment = $1 AND status = $2
	`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	rows, err := q.Query(ctx, query, commitment.String(), int(status))
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to query deposit by commitment and status: %w", err))
	}
	defer rows.Close()

	deposit, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[DvpDeposit])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil //nolint:nilnil // returning nil,nil is intentional for "not found" pattern
		}
		return nil, withstack.Wrap(fmt.Errorf("failed to decode deposit by commitment and status: %w", err))
	}

	return convertModelToDomain(deposit)
}

func (r *DvpDepositRepository) GetNonFungibleDeposit(
	ctx context.Context,
	tokenId string,
	tokenAddress string,
	userAddress string,
	tokenType types.DvpTokenType,
	status types.DvpDepositStatus,
) (*types.DvpDeposit, error) {
	query := fmt.Sprintf(`
		SELECT user_address, salt, token_amount, token_address,
		       token_type, token_id, tree_number, commitment,
		       nullifier, status, created_at
		FROM %s
		WHERE token_id = $1 AND token_address = $2 AND token_type = $3 AND status = $4
	`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	rows, err := q.Query(ctx, query,
		tokenId, tokenAddress, int(tokenType), int(status),
	)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to query non-fungible deposit: %w", err))
	}
	defer rows.Close()

	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[DvpDeposit])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil //nolint:nilnil // returning nil,nil is intentional for "not found" pattern
		}
		return nil, withstack.Wrap(fmt.Errorf("failed to decode non-fungible deposit: %w", err))
	}

	return convertModelToDomain(record)
}

func (r *DvpDepositRepository) GetFungibleDeposits(
	ctx context.Context,
	tokenAddress string,
	userAddress string,
	tokenType types.DvpTokenType,
	status types.DvpDepositStatus,
) ([]types.DvpDeposit, error) {
	query := fmt.Sprintf(`
		SELECT user_address, salt, token_amount, token_address,
		       token_type, token_id, tree_number, commitment,
		       nullifier, status, created_at
		FROM %s
		WHERE token_address = $1 AND user_address = $2 AND token_type = $3 AND status = $4
		ORDER BY created_at ASC
	`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	rows, err := q.Query(ctx, query,
		tokenAddress, userAddress, int(tokenType), int(status),
	)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to query fungible deposits: %w", err))
	}
	defer rows.Close()

	records, err := pgx.CollectRows(rows, pgx.RowToStructByName[DvpDeposit])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil //nolint:nilnil // returning nil,nil is intentional for "not found" pattern
		}
		return nil, withstack.Wrap(fmt.Errorf("failed to decode fungible deposits: %w", err))
	}

	deposits := make([]types.DvpDeposit, len(records))
	for i, record := range records {
		deposit, _ := convertModelToDomain(record)
		deposits[i] = *deposit
	}
	return deposits, nil
}

func (r *DvpDepositRepository) GetDepositsByToken(
	ctx context.Context,
	tokenAddress string,
	tokenId string,
	tokenType types.DvpTokenType,
	userAddress string,
	status types.DvpDepositStatus,
) ([]types.DvpDeposit, error) {
	query := fmt.Sprintf(`
		SELECT user_address, salt, token_amount, token_address,
		       token_type, token_id, tree_number, commitment,
		       nullifier, status, created_at
		FROM %s
		WHERE user_address = $1 AND token_address = $2 AND token_id = $3 AND token_type = $4 AND status = $5
		ORDER BY created_at ASC
	`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	rows, err := q.Query(ctx, query,
		userAddress, tokenAddress, tokenId, int(tokenType), int(status),
	)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to query deposits by token: %w", err))
	}
	defer rows.Close()

	records, err := pgx.CollectRows(rows, pgx.RowToStructByName[DvpDeposit])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil //nolint:nilnil // returning nil,nil is intentional for "not found" pattern
		}
		return nil, withstack.Wrap(fmt.Errorf("failed to decode deposits by token: %w", err))
	}

	deposits := make([]types.DvpDeposit, len(records))
	for i, record := range records {
		deposit, _ := convertModelToDomain(record)
		deposits[i] = *deposit
	}
	return deposits, nil
}

// Upsert the deposit nullifier field
func (r *DvpDepositRepository) UpsertDepositNullifier(ctx context.Context, commitment *big.Int, nullifier *big.Int) error {
	query := fmt.Sprintf(`UPDATE %s SET nullifier = $1 WHERE commitment = $2`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	_, err := q.Exec(ctx, query, nullifier.String(), commitment.String())
	if err != nil {
		return withstack.Wrap(fmt.Errorf("failed to upsert deposit nullifier: %w", err))
	}
	return nil
}

// BatchUpsertNullifiers updates nullifier fields for multiple deposits in a single bulk operation.
// Uses unnest to perform a single atomic UPDATE statement, avoiding the need for manual transaction
// management and ensuring compatibility with context-based transactions via GetQuerier.
func (r *DvpDepositRepository) BatchUpsertNullifiers(
	ctx context.Context,
	commitmentNullifierMap map[string]string,
) error {
	if len(commitmentNullifierMap) == 0 {
		return nil
	}

	commitments := make([]string, 0, len(commitmentNullifierMap))
	nullifiers := make([]string, 0, len(commitmentNullifierMap))
	for c, n := range commitmentNullifierMap {
		commitments = append(commitments, c)
		nullifiers = append(nullifiers, n)
	}

	query := fmt.Sprintf(`
		UPDATE %s AS t
		SET nullifier = data.nullifier
		FROM unnest($1::text[], $2::text[]) AS data(commitment, nullifier)
		WHERE t.commitment = data.commitment`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	_, err := q.Exec(ctx, query, commitments, nullifiers)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("failed to batch upsert nullifiers: %w", err))
	}
	return nil
}

func (r *DvpDepositRepository) UpdateDepositStatusByNullifier(
	ctx context.Context,
	tokenAddress string,
	nullifier *big.Int,
	status types.DvpDepositStatus,
) error {
	query := fmt.Sprintf(
		`UPDATE %s SET status = $1 WHERE token_address = $2 AND nullifier = $3`,
		r.tableName,
	)

	q := GetQuerier(ctx, r.pool)
	_, err := q.Exec(ctx, query, int(status), tokenAddress, nullifier.String())
	if err != nil {
		return withstack.Wrap(fmt.Errorf("failed to update deposit status by nullifier: %w", err))
	}

	return nil
}

func convertModelToDomain(record DvpDeposit) (*types.DvpDeposit, error) {

	salt, ok := new(big.Int).SetString(record.Salt, 10)
	if !ok {
		return nil, fmt.Errorf("failed to convert salt to big.Int: %s", record.Salt)
	}

	tokenAmount, ok := new(big.Int).SetString(record.TokenAmount, 10)
	if !ok {
		return nil, fmt.Errorf("failed to convert token amount to big.Int: %s", record.TokenAmount)
	}
	commitment, ok := new(big.Int).SetString(record.Commitment, 10)
	if !ok {
		return nil, fmt.Errorf("failed to convert commitment to big.Int: %s", record.Commitment)
	}

	nullifier := big.NewInt(0)
	if record.Nullifier != "" {
		nullifier, ok = new(big.Int).SetString(record.Nullifier, 10)
		if !ok {
			return nil, fmt.Errorf("failed to convert nullifier to big.Int: %s", record.Nullifier)
		}
	}

	return &types.DvpDeposit{
		UserAddress:  record.UserAddress,
		Salt:         salt,
		TokenAmount:  tokenAmount,
		TokenAddress: record.TokenAddress,
		TokenType:    record.TokenType,
		TokenID:      record.TokenID,
		TreeNumber:   record.TreeNumber,
		Commitment:   commitment,
		Nullifier:    nullifier,
		Status:       record.Status,
		CreatedAt:    record.CreatedAt,
	}, nil
}
