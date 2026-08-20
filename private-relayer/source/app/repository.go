package app

import (
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository"
)

type Repositories struct {
	lastProcessedBlockRepository *repository.LastProcessedBlockRepository
	transactionRepository        *repository.TransactionRepository
	signatureRepository          *repository.CalldataSignatureRepository
	atomicStatusRepository       *repository.AtomicStatusRepository

	enygma         *repository.EnygmaRepository
	enygmaHistory  *repository.EnygmaHistoryRepository
	txRecoveryData *repository.TxRecoveryDataRepository
	dvpDeposit     *repository.DvpDepositRepository
	dvpSwap        *repository.DvpSwapRepository
	merkleTree     *repository.MerkleTreeRepository
}

func (p *SourcePrivateRelayer) initializeRepositories() {
	p.repositories = &Repositories{
		lastProcessedBlockRepository: repository.NewLastProcessedBlockRepository(p.pool),
		transactionRepository:        repository.NewTransactionRepository(p.pool),
		signatureRepository:          repository.NewCalldataSignatureRepository(p.pool),
		atomicStatusRepository:       repository.NewAtomicStatusRepository(p.pool),

		enygma:         repository.NewEnygmaRepository(p.pool),
		enygmaHistory:  repository.NewEnygmaHistoryRepository(p.pool),
		txRecoveryData: repository.NewTxRecoveryDataRepository(p.pool),
		dvpDeposit:     repository.NewDvpDepositRepository(p.pool),
		dvpSwap:        repository.NewDvpSwapRepository(p.pool),
		merkleTree:     repository.NewMerkleTreeRepository(p.pool),
	}
}
