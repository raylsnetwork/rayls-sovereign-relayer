// NOTE: shared transaction store (Teleport bridge + generic relay); no atomic-exclusive members, so nothing here is decommissioned.

package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

type TransactionRepository struct {
	pool      *pgxpool.Pool
	tableName string
}

func NewTransactionRepository(pool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{
		pool:      pool,
		tableName: TransactionsCollectionName,
	}
}

const transactionSelectColumns = `
	id, tx_hash, tx_hash_destination, log_index, shared_id, state, outcome, proof_invalid,
	originator_chain_id, destination_chain_id, msg_id, is_atomic,
	updated_at, created_at, batch_id, batch_tx_hash_on_private_hub, resource_id,
	from_contract_address, from_user_address, transfer_metadata_id,
	transfer_metadata_amount, block_number, parent_hash
`

const transactionInsertColumns = `
	tx_hash, tx_hash_destination, log_index, shared_id, state, outcome, proof_invalid,
	originator_chain_id, destination_chain_id, msg_id, is_atomic,
	updated_at, batch_id, batch_tx_hash_on_private_hub, resource_id,
	from_contract_address, from_user_address, transfer_metadata_id,
	transfer_metadata_amount, block_number, parent_hash
`

const transactionInsertPlaceholders = `
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
`

func (r *TransactionRepository) GetByHashes(ctx context.Context, hashes []string) ([]types.Transaction, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE tx_hash = ANY($1)
	`, transactionSelectColumns, r.tableName)

	rows, err := r.pool.Query(ctx, query, hashes)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to find transactions: %w", err))
	}
	defer rows.Close()

	results, err := pgx.CollectRows(rows, pgx.RowToStructByName[Transaction])
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to decode transactions: %w", err))
	}

	txTypeSlice := make([]types.Transaction, 0, len(results))
	for _, tx := range results {
		txType, _ := transactionModelToTransactionType(tx)
		txTypeSlice = append(txTypeSlice, txType)
	}

	return txTypeSlice, nil
}

func (r *TransactionRepository) GetByState(
	ctx context.Context,
	state types.TransactionState,
	opts ...Option,
) ([]types.Transaction, error) {
	queryOptions := GetQueryOptions(opts...)

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE state = $1
	`, transactionSelectColumns, r.tableName)

	args := []any{state}
	if queryOptions.Limit > 0 {
		query += " LIMIT $2"
		args = append(args, queryOptions.Limit)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("error getting transactions by state %v: %w", state, err))
	}
	defer rows.Close()

	txs, err := pgx.CollectRows(rows, pgx.RowToStructByName[Transaction])
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("error decoding transactions by state %v: %w", state, err))
	}

	txTypeSlice := make([]types.Transaction, 0, len(txs))
	for _, record := range txs {
		txType, _ := transactionModelToTransactionType(record)
		txTypeSlice = append(txTypeSlice, txType)
	}

	return txTypeSlice, nil
}

func (r *TransactionRepository) GetByStateAndAtomicity(
	ctx context.Context,
	state types.TransactionState,
	isAtomic bool,
	opts ...Option,
) ([]types.Transaction, error) {
	queryOptions := GetQueryOptions(opts...)

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE state = $1 AND is_atomic = $2
	`, transactionSelectColumns, r.tableName)

	args := []any{state, isAtomic}
	if queryOptions.Limit > 0 {
		query += " LIMIT $3"
		args = append(args, queryOptions.Limit)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, withstack.Wrap(
			fmt.Errorf("error getting transactions by state %v atomicity %v: %w", state, isAtomic, err),
		)
	}
	defer rows.Close()

	txs, err := pgx.CollectRows(rows, pgx.RowToStructByName[Transaction])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, withstack.Wrap(
			fmt.Errorf("error decoding transactions by state %v atomicity %v: %w", state, isAtomic, err),
		)
	}

	txTypeSlice := make([]types.Transaction, 0, len(txs))
	for _, record := range txs {
		txType, _ := transactionModelToTransactionType(record)
		txTypeSlice = append(txTypeSlice, txType)
	}

	return txTypeSlice, nil
}

func (r *TransactionRepository) GetByStateOutcomeAndAtomicity(
	ctx context.Context,
	state types.TransactionState,
	outcome types.TransactionOutcome,
	isAtomic bool,
	opts ...Option,
) ([]types.Transaction, error) {
	queryOptions := GetQueryOptions(opts...)

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE state = $1 AND outcome = $2 AND is_atomic = $3
	`, transactionSelectColumns, r.tableName)

	args := []any{state, outcome, isAtomic}
	if queryOptions.Limit > 0 {
		query += " LIMIT $4"
		args = append(args, queryOptions.Limit)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, withstack.Wrap(
			fmt.Errorf("error getting transactions by state %v outcome %v atomicity %v: %w", state, outcome, isAtomic, err),
		)
	}
	defer rows.Close()

	txs, err := pgx.CollectRows(rows, pgx.RowToStructByName[Transaction])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, withstack.Wrap(
			fmt.Errorf("error decoding transactions by state %v outcome %v atomicity %v: %w", state, outcome, isAtomic, err),
		)
	}

	txTypeSlice := make([]types.Transaction, 0, len(txs))
	for _, record := range txs {
		txType, _ := transactionModelToTransactionType(record)
		txTypeSlice = append(txTypeSlice, txType)
	}

	return txTypeSlice, nil
}

func (r *TransactionRepository) GetByStateOutcomesAndAtomicity(
	ctx context.Context,
	state types.TransactionState,
	outcomes []types.TransactionOutcome,
	isAtomic bool,
	opts ...Option,
) ([]types.Transaction, error) {
	queryOptions := GetQueryOptions(opts...)

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE state = $1 AND outcome = ANY($2) AND is_atomic = $3
	`, transactionSelectColumns, r.tableName)

	args := []any{state, outcomes, isAtomic}
	if queryOptions.Limit > 0 {
		query += " LIMIT $4"
		args = append(args, queryOptions.Limit)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, withstack.Wrap(
			fmt.Errorf("error getting transactions by state %v outcomes %v atomicity %v: %w", state, outcomes, isAtomic, err),
		)
	}
	defer rows.Close()

	txs, err := pgx.CollectRows(rows, pgx.RowToStructByName[Transaction])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, withstack.Wrap(
			fmt.Errorf("error decoding transactions by state %v outcomes %v atomicity %v: %w", state, outcomes, isAtomic, err),
		)
	}

	txTypeSlice := make([]types.Transaction, 0, len(txs))
	for _, record := range txs {
		txType, _ := transactionModelToTransactionType(record)
		txTypeSlice = append(txTypeSlice, txType)
	}

	return txTypeSlice, nil
}

func (r *TransactionRepository) GetByStatesOutcomeAndAtomicity(
	ctx context.Context,
	states []types.TransactionState,
	outcome types.TransactionOutcome,
	isAtomic bool,
	opts ...Option,
) ([]types.Transaction, error) {
	queryOptions := GetQueryOptions(opts...)

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE state = ANY($1) AND outcome = $2 AND is_atomic = $3
	`, transactionSelectColumns, r.tableName)

	args := []any{states, outcome, isAtomic}
	if queryOptions.Limit > 0 {
		query += " LIMIT $4"
		args = append(args, queryOptions.Limit)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, withstack.Wrap(
			fmt.Errorf("error getting transactions by states %v outcome %v atomicity %v: %w", states, outcome, isAtomic, err),
		)
	}
	defer rows.Close()

	txs, err := pgx.CollectRows(rows, pgx.RowToStructByName[Transaction])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, withstack.Wrap(
			fmt.Errorf("error decoding transactions by states %v outcome %v atomicity %v: %w", states, outcome, isAtomic, err),
		)
	}

	txTypeSlice := make([]types.Transaction, 0, len(txs))
	for _, record := range txs {
		txType, _ := transactionModelToTransactionType(record)
		txTypeSlice = append(txTypeSlice, txType)
	}

	return txTypeSlice, nil
}

func (r *TransactionRepository) GetBySharedIDs(ctx context.Context, sharedIDs []string) ([]types.Transaction, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE shared_id = ANY($1)
	`, transactionSelectColumns, r.tableName)

	rows, err := r.pool.Query(ctx, query, sharedIDs)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to query transactions by shared IDs: %w", err))
	}
	defer rows.Close()

	txs, err := pgx.CollectRows(rows, pgx.RowToStructByName[Transaction])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, withstack.Wrap(fmt.Errorf("failed to decode transactions by shared IDs: %w", err))
	}

	txTypeSlice := make([]types.Transaction, 0, len(txs))
	for _, record := range txs {
		txType, _ := transactionModelToTransactionType(record)
		txTypeSlice = append(txTypeSlice, txType)
	}

	return txTypeSlice, nil
}

func (r *TransactionRepository) Create(ctx context.Context, tx types.Transaction) error {
	txModel := transactionTypeToTransactionModel(tx)

	query := fmt.Sprintf(`
		INSERT INTO %s (%s) VALUES (%s)
	`, r.tableName, transactionInsertColumns, transactionInsertPlaceholders)

	_, err := r.pool.Exec(ctx, query,
		txModel.TxHash, txModel.TxHashDestination, txModel.LogIndex, txModel.SharedId, txModel.State, txModel.Outcome, txModel.ProofInvalid,
		txModel.OriginatorChainId, txModel.DestinationChainId, txModel.MsgId[:], txModel.IsAtomic,
		txModel.UpdatedAt, txModel.BatchId, txModel.BatchTxHashOnPrivateHub, txModel.ResourceId,
		txModel.FromContractAddress, txModel.FromUserAddress, txModel.TransferMetadata_Id,
		txModel.TransferMetadata_Amount, txModel.BlockNumber, txModel.ParentHash,
	)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("failed to create transaction: %w", err))
	}
	return nil
}

func (r *TransactionRepository) BatchCreate(ctx context.Context, txs []types.Transaction) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("error creating batch transactions: %w", err))
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := fmt.Sprintf(`
		INSERT INTO %s (%s) VALUES (%s)
	`, r.tableName, transactionInsertColumns, transactionInsertPlaceholders)

	batch := &pgx.Batch{}
	for _, t := range txs {
		txModel := transactionTypeToTransactionModel(t)
		batch.Queue(query,
			txModel.TxHash, txModel.TxHashDestination, txModel.LogIndex, txModel.SharedId, txModel.State, txModel.Outcome, txModel.ProofInvalid,
			txModel.OriginatorChainId, txModel.DestinationChainId, txModel.MsgId[:], txModel.IsAtomic,
			txModel.UpdatedAt, txModel.BatchId, txModel.BatchTxHashOnPrivateHub, txModel.ResourceId,
			txModel.FromContractAddress, txModel.FromUserAddress, txModel.TransferMetadata_Id,
			txModel.TransferMetadata_Amount, txModel.BlockNumber, txModel.ParentHash,
		)
	}

	br := tx.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return withstack.Wrap(fmt.Errorf("failed to insert transactions: %w", err))
	}

	if err := tx.Commit(ctx); err != nil {
		return withstack.Wrap(fmt.Errorf("database transaction failed: %w", err))
	}

	return nil
}

func (r *TransactionRepository) BatchCreateWithStateAndOutcome(
	ctx context.Context,
	txs []types.Transaction,
	state types.TransactionState,
	outcome types.TransactionOutcome,
) error {
	slog.Debug("Batch Creating Transactions...",
		slog.Int("Txs", len(txs)),
		slog.Int("State", int(state)),
		slog.String("Outcome", string(outcome)),
	)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("error creating batch transactions with state: %w", err))
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := fmt.Sprintf(`
		INSERT INTO %s (%s) VALUES (%s)
	`, r.tableName, transactionInsertColumns, transactionInsertPlaceholders)

	batch := &pgx.Batch{}
	for _, t := range txs {
		txModel := transactionTypeToTransactionModel(t)
		txModel.State = state
		txModel.Outcome = outcome
		batch.Queue(query,
			txModel.TxHash, txModel.TxHashDestination, txModel.LogIndex, txModel.SharedId, txModel.State, txModel.Outcome, txModel.ProofInvalid,
			txModel.OriginatorChainId, txModel.DestinationChainId, txModel.MsgId[:], txModel.IsAtomic,
			txModel.UpdatedAt, txModel.BatchId, txModel.BatchTxHashOnPrivateHub, txModel.ResourceId,
			txModel.FromContractAddress, txModel.FromUserAddress, txModel.TransferMetadata_Id,
			txModel.TransferMetadata_Amount, txModel.BlockNumber, txModel.ParentHash,
		)
	}

	br := tx.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return withstack.Wrap(fmt.Errorf("failed to insert transactions: %w", err))
	}

	if err := tx.Commit(ctx); err != nil {
		return withstack.Wrap(fmt.Errorf("database transaction failed: %w", err))
	}

	return nil
}

// BatchSetState transitions the rows to a new state and resets outcome to pending.
func (r *TransactionRepository) BatchSetState(
	ctx context.Context,
	sharedIDs []string,
	state types.TransactionState,
) error {
	query := fmt.Sprintf(`UPDATE %s SET state = $1, outcome = $2 WHERE shared_id = ANY($3)`, r.tableName)

	_, err := r.pool.Exec(ctx, query, state, types.OutcomePending, sharedIDs)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("error updating batch transactions state: %w", err))
	}

	return nil
}

// BatchSetOutcome resolves the current state by setting an outcome; state column is untouched.
func (r *TransactionRepository) BatchSetOutcome(
	ctx context.Context,
	sharedIDs []string,
	outcome types.TransactionOutcome,
) error {
	query := fmt.Sprintf(`UPDATE %s SET outcome = $1 WHERE shared_id = ANY($2)`, r.tableName)

	_, err := r.pool.Exec(ctx, query, outcome, sharedIDs)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("error updating batch transactions outcome: %w", err))
	}

	return nil
}

// BatchSetStateAndOutcome atomically writes both state and outcome.
func (r *TransactionRepository) BatchSetStateAndOutcome(
	ctx context.Context,
	sharedIDs []string,
	state types.TransactionState,
	outcome types.TransactionOutcome,
) error {
	query := fmt.Sprintf(`UPDATE %s SET state = $1, outcome = $2 WHERE shared_id = ANY($3)`, r.tableName)

	_, err := r.pool.Exec(ctx, query, state, outcome, sharedIDs)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("error updating batch transactions state+outcome: %w", err))
	}

	return nil
}

// BatchSetProofInvalid flags rows whose proof verification failed.
func (r *TransactionRepository) BatchSetProofInvalid(
	ctx context.Context,
	sharedIDs []string,
) error {
	query := fmt.Sprintf(`UPDATE %s SET proof_invalid = TRUE WHERE shared_id = ANY($1)`, r.tableName)

	_, err := r.pool.Exec(ctx, query, sharedIDs)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("error flagging batch transactions proof_invalid: %w", err))
	}

	return nil
}

func (r *TransactionRepository) UpdateDestinationHashForSharedID(
	ctx context.Context,
	sharedID string,
	destinationHash common.Hash,
) error {
	query := fmt.Sprintf(`UPDATE %s SET tx_hash_destination = $1 WHERE shared_id = $2`, r.tableName)

	_, err := r.pool.Exec(ctx, query, destinationHash.String(), sharedID)
	if err != nil {
		return withstack.Wrap(err)
	}

	return nil
}

func (r *TransactionRepository) BatchUpdateDestinationHashForSharedIDs(
	ctx context.Context,
	hashBySharedID map[string]common.Hash,
) error {
	if len(hashBySharedID) == 0 {
		return nil
	}

	query := fmt.Sprintf(`UPDATE %s SET tx_hash_destination = $1 WHERE shared_id = $2`, r.tableName)

	batch := &pgx.Batch{}
	for sharedID, hash := range hashBySharedID {
		batch.Queue(query, hash.String(), sharedID)
	}

	br := r.pool.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return withstack.Wrap(fmt.Errorf("failed to batch update destination hashes: %w", err))
	}

	return nil
}

func transactionModelToTransactionType(modelTx Transaction) (types.Transaction, error) {
	fromChainID := new(big.Int)
	if _, ok := fromChainID.SetString(modelTx.OriginatorChainId, 10); !ok {
		return types.Transaction{}, fmt.Errorf("invalid originator_chain_id: %s", modelTx.OriginatorChainId)
	}

	toChainID := new(big.Int)
	if _, ok := toChainID.SetString(modelTx.DestinationChainId, 10); !ok {
		return types.Transaction{}, fmt.Errorf("invalid destination_chain_id: %s", modelTx.DestinationChainId)
	}

	var msgId [32]byte
	copy(msgId[:], modelTx.MsgId)

	return types.Transaction{
		TxHash:              modelTx.TxHash,
		TxHashDestination:   common.HexToHash(modelTx.TxHashDestination),
		LogIndex:            modelTx.LogIndex,
		SharedID:            modelTx.SharedId,
		State:               modelTx.State,
		Outcome:             modelTx.Outcome,
		ProofInvalid:        modelTx.ProofInvalid,
		FromChainID:         fromChainID,
		ToChainID:           toChainID,
		MsgID:               msgId,
		IsAtomic:            modelTx.IsAtomic,
		UpdatedAt:           modelTx.UpdatedAt,
		CreatedAt:           modelTx.CreatedAt,
		BatchID:             modelTx.BatchId,
		BatchPrivateHubHash: common.HexToHash(modelTx.BatchTxHashOnPrivateHub),
		ResourceID:          modelTx.ResourceId,
		FromContractAddress: modelTx.FromContractAddress,
		FromUserAddress:     modelTx.FromUserAddress,
		TransferID:          modelTx.TransferMetadata_Id,
		TransferAmount:      modelTx.TransferMetadata_Amount,
		BlockNumber:         modelTx.BlockNumber,
		ParentHash:          modelTx.ParentHash,
	}, nil
}

func transactionTypeToTransactionModel(typeTx types.Transaction) Transaction {
	outcome := typeTx.Outcome
	if outcome == "" {
		outcome = types.OutcomePending
	}
	return Transaction{
		TxHash:                  typeTx.TxHash,
		TxHashDestination:       typeTx.TxHashDestination.Hex(),
		LogIndex:                typeTx.LogIndex,
		SharedId:                typeTx.SharedID,
		State:                   typeTx.State,
		Outcome:                 outcome,
		ProofInvalid:            typeTx.ProofInvalid,
		OriginatorChainId:       typeTx.FromChainID.String(),
		DestinationChainId:      typeTx.ToChainID.String(),
		MsgId:                   typeTx.MsgID[:],
		IsAtomic:                typeTx.IsAtomic,
		UpdatedAt:               typeTx.UpdatedAt,
		CreatedAt:               typeTx.CreatedAt,
		BatchId:                 typeTx.BatchID,
		BatchTxHashOnPrivateHub: typeTx.BatchPrivateHubHash.Hex(),
		ResourceId:              typeTx.ResourceID,
		FromContractAddress:     typeTx.FromContractAddress,
		FromUserAddress:         typeTx.FromUserAddress,
		TransferMetadata_Id:     typeTx.TransferID,
		TransferMetadata_Amount: typeTx.TransferAmount,
		BlockNumber:             typeTx.BlockNumber,
		ParentHash:              typeTx.ParentHash,
	}
}
