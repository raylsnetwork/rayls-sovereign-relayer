package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/txutil"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

// txPersisterRecoveryRepository defines the recovery repository methods needed for crash recovery.
type txPersisterRecoveryRepository interface {
	Insert(ctx context.Context, data types.TxRecoveryData) error
	GetPendingByEventKey(ctx context.Context, resourceID string, blockNumber uint64, fromChainID string, eventType types.EnygmaEventType) (*types.TxRecoveryData, error)
	MarkConfirmed(ctx context.Context, privateHubTxHash string) error
}

var _ txPersisterRecoveryRepository = (*repository.TxRecoveryDataRepository)(nil)

// TxPersister implements the persist-then-broadcast pattern for DvP blockchain calls.
type TxPersister struct {
	recoveryRepo txPersisterRecoveryRepository
	broadcaster  txutil.Broadcaster
	receipter    txutil.Receipter
	txSimulator  txutil.TransactionSimulator
}

// NewTxPersister creates a new TxPersister.
func NewTxPersister(
	recoveryRepo txPersisterRecoveryRepository,
	broadcaster txutil.Broadcaster,
	receipter txutil.Receipter,
	txSimulator txutil.TransactionSimulator,
) *TxPersister {
	return &TxPersister{
		recoveryRepo: recoveryRepo,
		broadcaster:  broadcaster,
		receipter:    receipter,
		txSimulator:  txSimulator,
	}
}

// PersistAndBroadcast encodes a signed tx, persists a pending recovery record,
// broadcasts the tx, waits for the receipt, and marks the recovery as confirmed.
func (p *TxPersister) PersistAndBroadcast(
	ctx context.Context,
	tx *ethTypes.Transaction,
	recoveryData types.TxRecoveryData,
) error {
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("failed to RLP-encode signed tx: %w", err))
	}

	recoveryData.PrivateHubTxHash = tx.Hash().Hex()
	recoveryData.TxBytes = txBytes
	recoveryData.Status = types.HistoryStatusPending

	if err := p.recoveryRepo.Insert(ctx, recoveryData); err != nil {
		if errors.Is(err, repository.ErrRecoveryAlreadyExists) {
			slog.Info("DvP persist-and-broadcast: skipping already-processed event",
				slog.String("privateHubTxHash", recoveryData.PrivateHubTxHash),
				slog.String("eventType", recoveryData.EventType.String()),
			)
			return nil
		}
		return withstack.Wrap(err)
	}

	if err := p.broadcastSignedTx(ctx, tx); err != nil {
		return err
	}

	if err := p.waitForReceipt(ctx, tx.Hash().Hex()); err != nil {
		return err
	}

	return p.recoveryRepo.MarkConfirmed(ctx, recoveryData.PrivateHubTxHash)
}

// ResumePendingTx re-broadcasts a stored signed tx from crash recovery and waits for its receipt.
func (p *TxPersister) ResumePendingTx(ctx context.Context, pending *types.TxRecoveryData) error {
	var tx ethTypes.Transaction
	if err := rlp.DecodeBytes(pending.TxBytes, &tx); err != nil {
		return withstack.Wrap(fmt.Errorf("failed to decode pending signed tx: %w", err))
	}

	slog.Info("DvP crash recovery: resuming pending transaction",
		slog.String("privateHubTxHash", pending.PrivateHubTxHash),
		slog.String("eventType", pending.EventType.String()),
	)

	if err := p.broadcastSignedTx(ctx, &tx); err != nil {
		return err
	}

	if err := p.waitForReceipt(ctx, pending.PrivateHubTxHash); err != nil {
		return err
	}

	return p.recoveryRepo.MarkConfirmed(ctx, pending.PrivateHubTxHash)
}

// CheckPendingRecovery checks if there's a pending recovery record for crash recovery.
// Returns the pending record if found, nil if not.
func (p *TxPersister) CheckPendingRecovery(
	ctx context.Context,
	resourceID string,
	blockNumber uint64,
	fromChainID string,
	eventType types.EnygmaEventType,
) (*types.TxRecoveryData, error) {
	return p.recoveryRepo.GetPendingByEventKey(ctx, resourceID, blockNumber, fromChainID, eventType)
}

func (p *TxPersister) broadcastSignedTx(ctx context.Context, tx *ethTypes.Transaction) error {
	return txutil.BroadcastSignedTx(ctx, p.broadcaster, tx)
}

func (p *TxPersister) waitForReceipt(ctx context.Context, txHash string) error {
	return txutil.WaitForReceipt(ctx, p.receipter, p.txSimulator, txHash)
}
