// Decommissioning Teleport (vanilla, atomic).

package service

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	sharedservice "github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/txsim"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

const revertReasonTimeout = 30 * time.Second

type receiptTeleportClient interface {
	ExecuteAtomicMessageBatch(context.Context, []string, []types.AtomicTeleportAdditionalData) error
	RevertAtomicMessageBatch(context.Context, []string, []types.AtomicTeleportAdditionalData) error
}

type receiptTransactionRepository interface {
	GetBySharedIDs(context.Context, []string) ([]types.Transaction, error)
	BatchSetStateAndOutcome(ctx context.Context, sharedIDs []string, state types.TransactionState, outcome types.TransactionOutcome) error
}

type receiptEthereumClient interface {
	TransactionReceipt(context.Context, common.Hash) (*ethTypes.Receipt, error)
	HeaderByHash(context.Context, common.Hash) (*ethTypes.Header, error)
}

type receiptTransactionSimulator interface {
	GetRevertReason(context.Context, common.Hash) (txsim.ContractError, error)
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type ReceiptService struct {
	teleportCli    receiptTeleportClient
	ethereumCli    receiptEthereumClient
	transactionSim receiptTransactionSimulator

	transactionRepo receiptTransactionRepository
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func NewReceiptService(
	teleportCli receiptTeleportClient,
	ethereumCli receiptEthereumClient,
	transactionRepo receiptTransactionRepository,
	transactionSim receiptTransactionSimulator,
) *ReceiptService {
	return &ReceiptService{
		teleportCli:    teleportCli,
		ethereumCli:    ethereumCli,
		transactionSim: transactionSim,

		transactionRepo: transactionRepo,
	}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (s *ReceiptService) HandleSuccessfullyMined(ctx context.Context, sharedIDs []string) error {
	var additionalDataSlice []types.AtomicTeleportAdditionalData

	txs, err := s.transactionRepo.GetBySharedIDs(ctx, sharedIDs)
	if err != nil {
		return sharedservice.WrapInAtomicServiceError("got error while querying transactions by shared IDs", err)
	}

	for _, tx := range txs {
		additionalData := types.AtomicTeleportAdditionalData{
			TxHashDestination:          tx.TxHashDestination,
			TxHashDestinationStatus:    1,
			TxHashDestinationTimestamp: s.getTimestampForTxHash(ctx, tx.TxHashDestination),
			SharedId:                   tx.SharedID,
			BatchPrivateHubHash:        tx.BatchPrivateHubHash,
		}
		additionalDataSlice = append(additionalDataSlice, additionalData)
	}

	err = s.teleportCli.ExecuteAtomicMessageBatch(ctx, sharedIDs, additionalDataSlice)
	outcome := types.OutcomeSuccess
	if err != nil {
		if !errors.Is(err, sharedservice.ErrAlreadyReverted) {
			return sharedservice.WrapInAtomicServiceError("got error while marking message batch as executed", err)
		}
		// Race lost: we minted on dest but PNH already considers this batch reverted
		// (source's timeout-revert beat us). Recording outcome=reverted lets the
		// finalization poller route this row through dest-revert sigs to undo the
		// orphan mint, without needing to wait for atomic_status to confirm.
		outcome = types.OutcomeReverted
	}

	err = s.transactionRepo.BatchSetStateAndOutcome(ctx, sharedIDs, types.HubNotifiedExec, outcome)
	if err != nil {
		return sharedservice.WrapInAtomicServiceError("got error while updating transaction state for shared IDs", err)
	}
	return nil
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (s *ReceiptService) HandleFailedMined(ctx context.Context, sharedIDs []string) error {
	var additionalDataSlice []types.AtomicTeleportAdditionalData

	txs, err := s.transactionRepo.GetBySharedIDs(ctx, sharedIDs)
	if err != nil {
		return sharedservice.WrapInAtomicServiceError("got error while querying transactions by shared IDs", err)
	}

	for _, tx := range txs {
		var revertReason txsim.ContractError
		if tx.TxHashDestination != (common.Hash{}) {
			ctxReason, cancelReason := context.WithTimeout(ctx, revertReasonTimeout)
			revertReason, _ = s.transactionSim.GetRevertReason(ctxReason, tx.TxHashDestination)
			cancelReason()
		}
		additionalData := types.AtomicTeleportAdditionalData{
			TxHashDestination:          tx.TxHashDestination,
			TxHashDestinationStatus:    1,
			TxHashDestinationTimestamp: s.getTimestampForTxHash(ctx, tx.TxHashDestination),
			RevertReason:               revertReason.String(),
			SharedId:                   tx.SharedID,
			BatchPrivateHubHash:        tx.BatchPrivateHubHash,
		}
		additionalDataSlice = append(additionalDataSlice, additionalData)
	}

	err = s.teleportCli.RevertAtomicMessageBatch(ctx, sharedIDs, additionalDataSlice)
	if err != nil {
		return sharedservice.WrapInAtomicServiceError("got error while marking message batch as reverted", err)
	}

	err = s.transactionRepo.BatchSetStateAndOutcome(ctx, sharedIDs, types.HubNotifiedRevert, types.OutcomeSuccess)
	if err != nil {
		return sharedservice.WrapInAtomicServiceError("got error while updating transaction state for shared IDs", err)
	}
	return nil
}

func (s *ReceiptService) getTimestampForTxHash(ctx context.Context, txHash common.Hash) uint64 {
	ctxReceipt, cancelReceipt := context.WithTimeout(ctx, 10*time.Second)
	receipt, err := s.ethereumCli.TransactionReceipt(ctxReceipt, txHash)
	cancelReceipt()
	if err != nil {
		return 0
	}

	ctxHeader, cancelHeader := context.WithTimeout(ctx, 10*time.Second)
	header, err := s.ethereumCli.HeaderByHash(ctxHeader, receipt.BlockHash)
	cancelHeader()
	if err == nil {
		return header.Time
	} else {
		return 0
	}
}
