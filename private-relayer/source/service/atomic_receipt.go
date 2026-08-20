// Decommissioning Teleport (vanilla, atomic).

package service

import (
	"context"

	sharedservice "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

type receiptTransactionRepository interface {
	BatchSetStateAndOutcome(ctx context.Context, sharedIDs []string, state types.TransactionState, outcome types.TransactionOutcome) error
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type ReceiptService struct {
	txRepo receiptTransactionRepository
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func NewReceiptService(txRepo receiptTransactionRepository) *ReceiptService {
	return &ReceiptService{
		txRepo: txRepo,
	}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (s *ReceiptService) HandleSuccessfullyMined(ctx context.Context, sharedIDs []string) error {
	err := s.txRepo.BatchSetStateAndOutcome(ctx, sharedIDs, types.SourcePublish, types.OutcomeSuccess)
	if err != nil {
		return sharedservice.WrapInAtomicServiceError("got error while trying to update state for shared IDs", err)
	}
	return nil
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (s *ReceiptService) HandleFailedMined(ctx context.Context, sharedIDs []string) error {
	err := s.txRepo.BatchSetStateAndOutcome(ctx, sharedIDs, types.SourcePublish, types.OutcomeFailed)
	if err != nil {
		return sharedservice.WrapInAtomicServiceError("got error while trying to update state for shared IDs", err)
	}
	return nil
}
