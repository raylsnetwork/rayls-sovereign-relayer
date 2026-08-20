// Decommissioning Teleport (vanilla, atomic).

package service

import (
	"context"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	sharedservice "github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

// messageTypeEarlyRevert is the wire MessageType for early-revert
// signatures. Both producer (HandleEarlyRevert) and the result router
// must use this string.
const messageTypeEarlyRevert = "atomic.early-revert"

type earlyRevertSignatureRepository interface {
	GetSourceRevertForSharedIDs(context.Context, []string) ([]types.CalldataSignature, error)
}

// earlyRevertSignatureBatchSender publishes signed-tx requests to CTS
// via the privatehub identity subject. Concrete implementation is
// *msgqueue.Publisher[types.TxRequest].
type earlyRevertSignatureBatchSender interface {
	PushBatch(ctx context.Context, msgs []types.TxRequest) error
}

type earlyRevertGenerator interface {
	Generate(signature types.CalldataSignature) ([]byte, error)
}

type earlyRevertTransactionRepository interface {
	BatchSetState(ctx context.Context, sharedIDs []string, state types.TransactionState) error
	BatchSetOutcome(ctx context.Context, sharedIDs []string, outcome types.TransactionOutcome) error
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type EarlyRevertService struct {
	signatureBatcher earlyRevertSignatureBatchSender
	generator        earlyRevertGenerator
	destEndpoint     common.Address

	signatureRepo earlyRevertSignatureRepository
	txRepo        earlyRevertTransactionRepository
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func NewEarlyRevertService(
	signatureBatcher earlyRevertSignatureBatchSender,
	generator earlyRevertGenerator,
	destEndpoint common.Address,
	signatureRepo earlyRevertSignatureRepository,
	txRepo earlyRevertTransactionRepository,
) *EarlyRevertService {
	return &EarlyRevertService{
		signatureBatcher: signatureBatcher,
		generator:        generator,
		destEndpoint:     destEndpoint,
		signatureRepo:    signatureRepo,
		txRepo:           txRepo,
	}
}

// HandleEarlyRevert is the publish side: load source-revert signatures,
// pack their calldata, publish a TxRequest batch on cts.send.privatehub,
// and transition the affected transactions to EarlyRevertSigs so the
// poller doesn't keep republishing them. The terminal outcome is set
// later by HandleEarlyRevertCallback when the TxResult arrives.
//
// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (s *EarlyRevertService) HandleEarlyRevert(ctx context.Context, sharedIDs []string) error {
	signatures, err := s.signatureRepo.GetSourceRevertForSharedIDs(ctx, sharedIDs)
	if err != nil {
		return sharedservice.WrapInAtomicServiceError("failed to get signatures from repository", err)
	}

	requests := s.generateMessages(signatures)
	if err := s.signatureBatcher.PushBatch(ctx, requests); err != nil {
		return sharedservice.WrapInAtomicServiceError("failed to publish early-revert signatures", err)
	}

	if err := s.txRepo.BatchSetState(
		ctx, sharedIDs, types.EarlyRevertSigs,
	); err != nil {
		return sharedservice.WrapInAtomicServiceError(
			"failed to update transaction state after publishing early-revert signatures", err,
		)
	}
	return nil
}

// HandleEarlyRevertCallback is the callback side: the result router
// invokes this with batches of TxResults that share MessageType
// `atomic.early-revert`. Per-result outcomes:
//
//   - TxResultSuccess → outcome=success on EarlyRevertSigs (terminal).
//   - TxResultRevert  → outcome=reverted on EarlyRevertSigs (terminal).
//   - TxResultFailed  → outcome=failed on EarlyRevertSigs (terminal).
//
// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (s *EarlyRevertService) HandleEarlyRevertCallback(
	ctx context.Context,
	results []types.TxResult,
) error {
	var (
		successIDs  []string
		revertedIDs []string
		failedIDs   []string
	)

	for _, res := range results {
		switch res.Kind {
		case types.TxResultSuccess:
			successIDs = append(successIDs, res.CorrelationID)
		case types.TxResultRevert:
			revertedIDs = append(revertedIDs, res.CorrelationID)
			slog.Error("Early-revert signature reverted",
				slog.String("shared_id", res.CorrelationID),
			)
		case types.TxResultFailed:
			failedIDs = append(failedIDs, res.CorrelationID)
			slog.Error("Early-revert signature failed",
				slog.String("shared_id", res.CorrelationID),
				slog.String("error", res.ErrorReason),
			)
		}
	}

	if len(successIDs) > 0 {
		if err := s.txRepo.BatchSetOutcome(
			ctx, successIDs, types.OutcomeSuccess,
		); err != nil {
			return sharedservice.WrapInAtomicServiceError("updating early-revert success outcome", err)
		}
	}
	if len(revertedIDs) > 0 {
		if err := s.txRepo.BatchSetOutcome(
			ctx, revertedIDs, types.OutcomeReverted,
		); err != nil {
			return sharedservice.WrapInAtomicServiceError("updating early-revert reverted outcome", err)
		}
	}
	if len(failedIDs) > 0 {
		if err := s.txRepo.BatchSetOutcome(
			ctx, failedIDs, types.OutcomeFailed,
		); err != nil {
			return sharedservice.WrapInAtomicServiceError("updating early-revert failed outcome", err)
		}
	}
	return nil
}

func (s *EarlyRevertService) generateMessages(
	signatures []types.CalldataSignature,
) []types.TxRequest {
	msgs := make([]types.TxRequest, 0, len(signatures))
	for _, sig := range signatures {
		calldata, err := s.generator.Generate(sig)
		if err != nil {
			slog.Warn("Failed to generate early-revert calldata, skipping",
				slog.String("shared_id", sig.SharedId),
				slog.Any("err", err),
			)
			continue
		}
		msgs = append(msgs, types.TxRequest{
			CorrelationID: sig.SharedId,
			MessageType:   messageTypeEarlyRevert,
			Address:       s.destEndpoint,
			Calldata:      calldata,
		})
	}
	return msgs
}
