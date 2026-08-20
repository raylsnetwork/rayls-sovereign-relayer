package app

import (
	"context"
	"log/slog"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
)

type Repositories struct {
	lastProcessedBlockRepository *repository.LastProcessedBlockRepository
	transactionRepository        *repository.TransactionRepository
	signatureRepository          *repository.CalldataSignatureRepository

	enygma           *repository.EnygmaRepository
	enygmaHistory    *repository.EnygmaHistoryRepository
	txRecoveryData   *repository.TxRecoveryDataRepository
	enygmaCheckpoint *repository.EnygmaCheckpointRepository
	resourceLock     *repository.ResourceLockRepository

	dvpDeposit *repository.DvpDepositRepository
	dvpSwap    *repository.DvpSwapRepository
	merkleTree *repository.MerkleTreeRepository

	atomicStatusRepository *repository.AtomicStatusRepository
}

func (r *DestPrivateRelayer) initializeRepositories(ctx context.Context) {
	resourceLock := repository.NewResourceLockRepository(r.pool)
	if lockErr := resourceLock.RemoveAllLocks(ctx); lockErr != nil {
		slog.Warn("failed to remove all resource locks during initialization", slog.Any("error", lockErr))
	}
	// Start cleanup routine for expired locks
	resourceLock.StartCleanupRoutine(ctx)

	r.repositories = &Repositories{
		lastProcessedBlockRepository: repository.NewLastProcessedBlockRepository(r.pool),
		transactionRepository:        repository.NewTransactionRepository(r.pool),
		signatureRepository:          repository.NewCalldataSignatureRepository(r.pool),

		enygma:           repository.NewEnygmaRepository(r.pool),
		enygmaHistory:    repository.NewEnygmaHistoryRepository(r.pool),
		txRecoveryData:   repository.NewTxRecoveryDataRepository(r.pool),
		enygmaCheckpoint: repository.NewEnygmaCheckpointRepository(r.pool),
		resourceLock:     resourceLock,

		dvpDeposit: repository.NewDvpDepositRepository(r.pool),
		dvpSwap:    repository.NewDvpSwapRepository(r.pool),
		merkleTree: repository.NewMerkleTreeRepository(r.pool),

		atomicStatusRepository: repository.NewAtomicStatusRepository(r.pool),
	}
}
