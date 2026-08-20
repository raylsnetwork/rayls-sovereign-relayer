// Decommissioning Teleport (vanilla, atomic).

package service

import (
	"context"
	"errors"

	sharedservice "github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

type expiredTeleportClient interface {
	RevertAtomicMessageBatch(context.Context, []string, []types.AtomicTeleportAdditionalData) error
	GetAtomicMessageStatuses(context.Context, []string) ([]types.AtomicStatusUpdateMessage, error)
}

type expiredTransactionRepository interface {
	BatchSetStateAndOutcome(ctx context.Context, sharedIDs []string, state types.TransactionState, outcome types.TransactionOutcome) error
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type ExpiredService struct {
	teleportCli expiredTeleportClient
	txRepo      expiredTransactionRepository
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func NewExpiredService(teleportCli expiredTeleportClient, txRepo expiredTransactionRepository) *ExpiredService {
	return &ExpiredService{
		teleportCli: teleportCli,
		txRepo:      txRepo,
	}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (s *ExpiredService) Handle(ctx context.Context, sharedIDs []string) error {
	atomicStatuses, err := s.teleportCli.GetAtomicMessageStatuses(ctx, sharedIDs)
	if err != nil {
		return sharedservice.WrapInAtomicServiceError("error while attempting to get atomic message statuses", err)
	}

	var pendingSharedIDs []string
	for _, atomicStatus := range atomicStatuses {
		if atomicStatus.Status == types.AtomicPendingStatus {
			pendingSharedIDs = append(pendingSharedIDs, atomicStatus.SharedID)
		}
	}
	if len(pendingSharedIDs) == 0 {
		return nil
	}

	err = s.teleportCli.RevertAtomicMessageBatch(ctx, pendingSharedIDs, []types.AtomicTeleportAdditionalData{})
	if err != nil {
		// TODO(race-handling): when PNH returns ErrAlreadyExecuted the row stays at
		// SourcePublish+success forever; should transition to SourceFinalized so the
		// timeout poller stops re-querying.
		if errors.Is(err, sharedservice.ErrAlreadyExecuted) {
			return nil
		}
		return sharedservice.WrapInAtomicServiceError("error while attempting revert atomic message batch", err)
	}

	// RevertAtomicMessageBatch is synchronous and only returns nil once the Hub tx is mined,
	// so the timeout-revert is already confirmed on-chain here. Record SourceTimeoutRevert+success
	// (NOT +pending) so the finalization poller — which consumes SourceTimeoutRevert+success —
	// picks the row up and drives the sender-side refund (SourceRevertSigs → revertTeleportMint).
	// Using BatchSetState would force outcome=pending and orphan the row, leaving the sender
	// un-refunded forever.
	err = s.txRepo.BatchSetStateAndOutcome(ctx, pendingSharedIDs, types.SourceTimeoutRevert, types.OutcomeSuccess)
	if err != nil {
		return sharedservice.WrapInAtomicServiceError("failed to update statuses for transactions", err)
	}
	return nil
}
