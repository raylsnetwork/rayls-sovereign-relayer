// Decommissioning Teleport (vanilla, atomic).

package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/txsim"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

const revertReasonTimeout = 30 * time.Second

type signatureSignatureRepository interface {
	GetDestinationUnlocksForSharedIDs(context.Context, []string) ([]types.CalldataSignature, error)
	GetDestinationRevertsForSharedIDs(context.Context, []string) ([]types.CalldataSignature, error)
	GetSourceRevertForSharedIDs(context.Context, []string) ([]types.CalldataSignature, error)
}

type signatureSignatureBatchSender interface {
	PushBatch(ctx context.Context, msgs []types.TxRequest) error
}

type signatureTeleportClient interface {
	SendAdditionalDataBatch(context.Context, []string, []types.AtomicTeleportAdditionalData) error
}

type signatureEthereumClient interface {
	HeaderByHash(context.Context, common.Hash) (*ethTypes.Header, error)
}

type signatureTransactionRepository interface {
	BatchSetState(ctx context.Context, sharedIDs []string, state types.TransactionState) error
	BatchSetOutcome(ctx context.Context, sharedIDs []string, outcome types.TransactionOutcome) error
	BatchSetStateAndOutcome(ctx context.Context, sharedIDs []string, state types.TransactionState, outcome types.TransactionOutcome) error
}

type signatureTransactionSimulator interface {
	DecodeRevertBytes(data []byte) (txsim.ContractError, error)
}

type signatureGenerator interface {
	Generate(signature types.CalldataSignature) ([]byte, error)
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type SignatureService struct {
	signatureBatcher        signatureSignatureBatchSender
	signatureTeleportClient signatureTeleportClient
	ethClient               signatureEthereumClient
	simulator               signatureTransactionSimulator
	generator               signatureGenerator

	destEndpoint common.Address

	signatureRepo   signatureSignatureRepository
	transactionRepo signatureTransactionRepository
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func NewSignatureService(
	signatureBatcher signatureSignatureBatchSender,
	signatureTeleportClient signatureTeleportClient,
	ethClient signatureEthereumClient,
	signatureRepo signatureSignatureRepository,
	transactionRepo signatureTransactionRepository,
	simulator signatureTransactionSimulator,
	generator signatureGenerator,
	destEndpoint common.Address,
) *SignatureService {
	return &SignatureService{
		signatureBatcher:        signatureBatcher,
		signatureTeleportClient: signatureTeleportClient,
		ethClient:               ethClient,
		simulator:               simulator,
		generator:               generator,

		destEndpoint: destEndpoint,

		signatureRepo:   signatureRepo,
		transactionRepo: transactionRepo,
	}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (s *SignatureService) HandleDestinationExecuted(ctx context.Context, sharedIDs []string) error {
	signatures, err := s.signatureRepo.GetDestinationUnlocksForSharedIDs(ctx, sharedIDs)
	if err != nil {
		return WrapInAtomicServiceError("got error while reading unlock signatures from repository", err)
	}

	err = s.signatureBatcher.PushBatch(ctx, s.generateMessages(signatures))
	if err != nil {
		return WrapInAtomicServiceError("failed to send signatures", err)
	}

	err = s.transactionRepo.BatchSetState(
		ctx, sharedIDs, types.DestinationUnlockSigs,
	)
	if err != nil {
		return WrapInAtomicServiceError("failed to update transaction state after sending signatures", err)
	}
	return nil
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (s *SignatureService) HandleDestinationExecutedCallback(ctx context.Context, results []types.TxResult) error {
	var (
		executedSharedIDs []string
		revertedSharedIDs []string
		failedSharedIDs   []string
	)
	executedReceipts := map[string]*ethTypes.Receipt{}

	for _, res := range results {
		switch res.Kind {
		case types.TxResultSuccess:
			executedSharedIDs = append(executedSharedIDs, res.CorrelationID)
			executedReceipts[res.CorrelationID] = res.Receipt
		case types.TxResultRevert:
			revertedSharedIDs = append(revertedSharedIDs, res.CorrelationID)
			reason, decodeErr := s.simulator.DecodeRevertBytes(res.RevertData)
			if decodeErr == nil {
				slog.Error("Destination unlock signature reverted",
					slog.String("shared_id", res.CorrelationID),
					slog.String("reason", reason.String()),
				)
			} else {
				slog.Error("Destination unlock signature reverted",
					slog.String("shared_id", res.CorrelationID),
					slog.String("reason", "unknown"),
					slog.Any("decode_err", decodeErr),
				)
			}
		case types.TxResultFailed:
			failedSharedIDs = append(failedSharedIDs, res.CorrelationID)
			slog.Error("Destination unlock signature failed",
				slog.String("shared_id", res.CorrelationID),
				slog.String("error", res.ErrorReason),
			)
		}
	}

	// Build additional-data batch from the same ordered iteration of
	// executedSharedIDs so the slice pairs positionally (in addition to
	// each struct carrying its own SharedId).
	additionalDataBatch := make([]types.AtomicTeleportAdditionalData, 0, len(executedSharedIDs))
	for _, sharedID := range executedSharedIDs {
		receipt := executedReceipts[sharedID]
		additionalDataBatch = append(additionalDataBatch, types.AtomicTeleportAdditionalData{
			SharedId:          sharedID,
			TxHashDestination: receipt.TxHash,
			//nolint:gosec // Ethereum receipt status is 0 or 1
			TxHashDestinationStatus:    int8(receipt.Status),
			TxHashDestinationTimestamp: s.getTimestampForBlockHash(ctx, receipt.BlockHash),
		})
	}

	if len(executedSharedIDs) > 0 {
		if err := s.signatureTeleportClient.SendAdditionalDataBatch(ctx, executedSharedIDs, additionalDataBatch); err != nil {
			return WrapInAtomicServiceError("failed to send additional data", err)
		}
		if err := s.transactionRepo.BatchSetOutcome(
			ctx, executedSharedIDs, types.OutcomeSuccess,
		); err != nil {
			return WrapInAtomicServiceError("updating executed signature outcome", err)
		}
	}
	if len(revertedSharedIDs) > 0 {
		if err := s.transactionRepo.BatchSetOutcome(
			ctx, revertedSharedIDs, types.OutcomeReverted,
		); err != nil {
			return WrapInAtomicServiceError("updating reverted signature outcome", err)
		}
	}
	if len(failedSharedIDs) > 0 {
		if err := s.transactionRepo.BatchSetOutcome(
			ctx, failedSharedIDs, types.OutcomeFailed,
		); err != nil {
			return WrapInAtomicServiceError("updating failed signature outcome", err)
		}
	}

	return nil
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (s *SignatureService) HandleDestinationReverted(ctx context.Context, sharedIDs []string) error {
	signatures, err := s.signatureRepo.GetDestinationRevertsForSharedIDs(ctx, sharedIDs)
	if err != nil {
		return WrapInAtomicServiceError("got error while reading destination revert signatures from repository", err)
	}

	err = s.signatureBatcher.PushBatch(ctx, s.generateMessages(signatures))
	if err != nil {
		return WrapInAtomicServiceError("failed to send signatures", err)
	}

	err = s.transactionRepo.BatchSetState(
		ctx, sharedIDs, types.DestinationRevertSigs,
	)
	if err != nil {
		return WrapInAtomicServiceError("failed to update transaction state after sending signatures", err)
	}
	return nil
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (s *SignatureService) HandleDestinationRevertedCallback(ctx context.Context, results []types.TxResult) error {
	var (
		revertedSharedIDs []string
		signatureReverted []string
		signatureFailed   []string
	)
	revertedReceipts := map[string]*ethTypes.Receipt{}

	for _, res := range results {
		switch res.Kind {
		case types.TxResultSuccess:
			revertedSharedIDs = append(revertedSharedIDs, res.CorrelationID)
			revertedReceipts[res.CorrelationID] = res.Receipt
		case types.TxResultRevert:
			signatureReverted = append(signatureReverted, res.CorrelationID)
			reason, decodeErr := s.simulator.DecodeRevertBytes(res.RevertData)
			if decodeErr == nil {
				slog.Error("Destination revert signature reverted",
					slog.String("shared_id", res.CorrelationID),
					slog.String("reason", reason.String()),
				)
			} else {
				slog.Error("Destination revert signature reverted",
					slog.String("shared_id", res.CorrelationID),
					slog.String("reason", "unknown"),
					slog.Any("decode_err", decodeErr),
				)
			}
		case types.TxResultFailed:
			signatureFailed = append(signatureFailed, res.CorrelationID)
			slog.Error("Destination revert signature failed",
				slog.String("shared_id", res.CorrelationID),
				slog.String("error", res.ErrorReason),
			)
		}
	}

	additionalDataBatch := make([]types.AtomicTeleportAdditionalData, 0, len(revertedSharedIDs))
	for _, sharedID := range revertedSharedIDs {
		receipt := revertedReceipts[sharedID]
		additionalDataBatch = append(additionalDataBatch, types.AtomicTeleportAdditionalData{
			SharedId:                sharedID,
			TxHashDestinationRevert: receipt.TxHash,
			//nolint:gosec // Ethereum receipt status is 0 or 1
			TxHashDestinationRevertStatus:    int8(receipt.Status),
			TxHashDestinationRevertTimestamp: s.getTimestampForBlockHash(ctx, receipt.BlockHash),
		})
	}

	if len(revertedSharedIDs) > 0 {
		if err := s.signatureTeleportClient.SendAdditionalDataBatch(ctx, revertedSharedIDs, additionalDataBatch); err != nil {
			return WrapInAtomicServiceError("failed to send additional data", err)
		}
		if err := s.transactionRepo.BatchSetOutcome(
			ctx, revertedSharedIDs, types.OutcomeSuccess,
		); err != nil {
			return WrapInAtomicServiceError("updating reverted signature outcome", err)
		}
	}
	if len(signatureReverted) > 0 {
		if err := s.transactionRepo.BatchSetOutcome(
			ctx, signatureReverted, types.OutcomeReverted,
		); err != nil {
			return WrapInAtomicServiceError("updating signature-revert outcome", err)
		}
	}
	if len(signatureFailed) > 0 {
		if err := s.transactionRepo.BatchSetOutcome(
			ctx, signatureFailed, types.OutcomeFailed,
		); err != nil {
			return WrapInAtomicServiceError("updating signature-failed outcome", err)
		}
	}

	return nil
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (s *SignatureService) HandleSourceReverted(ctx context.Context, sharedIDs []string) error {
	signatures, err := s.signatureRepo.GetSourceRevertForSharedIDs(ctx, sharedIDs)
	if err != nil {
		return WrapInAtomicServiceError("got error while reading source revert signatures from repository", err)
	}

	err = s.signatureBatcher.PushBatch(ctx, s.generateMessages(signatures))
	if err != nil {
		return WrapInAtomicServiceError("failed to send signatures", err)
	}

	err = s.transactionRepo.BatchSetState(
		ctx, sharedIDs, types.SourceRevertSigs,
	)
	if err != nil {
		return WrapInAtomicServiceError("failed to update transaction state after sending signatures", err)
	}
	return nil
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (s *SignatureService) HandleSourceRevertedCallback(ctx context.Context, results []types.TxResult) error {
	var (
		revertedSharedIDs []string
		signatureReverted []string
		signatureFailed   []string
	)
	revertedReceipts := map[string]*ethTypes.Receipt{}

	for _, res := range results {
		switch res.Kind {
		case types.TxResultSuccess:
			revertedSharedIDs = append(revertedSharedIDs, res.CorrelationID)
			revertedReceipts[res.CorrelationID] = res.Receipt
		case types.TxResultRevert:
			signatureReverted = append(signatureReverted, res.CorrelationID)
			reason, decodeErr := s.simulator.DecodeRevertBytes(res.RevertData)
			if decodeErr == nil {
				slog.Error("Source revert signature reverted",
					slog.String("shared_id", res.CorrelationID),
					slog.String("reason", reason.String()),
				)
			} else {
				slog.Error("Source revert signature reverted",
					slog.String("shared_id", res.CorrelationID),
					slog.String("reason", "unknown"),
					slog.Any("decode_err", decodeErr),
				)
			}
		case types.TxResultFailed:
			signatureFailed = append(signatureFailed, res.CorrelationID)
			slog.Error("Source revert signature failed",
				slog.String("shared_id", res.CorrelationID),
				slog.String("error", res.ErrorReason),
			)
		}
	}

	additionalDataBatch := make([]types.AtomicTeleportAdditionalData, 0, len(revertedSharedIDs))
	for _, sharedID := range revertedSharedIDs {
		receipt := revertedReceipts[sharedID]
		additionalDataBatch = append(additionalDataBatch, types.AtomicTeleportAdditionalData{
			SharedId:           sharedID,
			TxHashSourceRevert: receipt.TxHash,
			//nolint:gosec // Ethereum receipt status is 0 or 1
			TxHashSourceRevertStatus:    int8(receipt.Status),
			TxHashSourceRevertTimestamp: s.getTimestampForBlockHash(ctx, receipt.BlockHash),
		})
	}

	if len(revertedSharedIDs) > 0 {
		if err := s.signatureTeleportClient.SendAdditionalDataBatch(ctx, revertedSharedIDs, additionalDataBatch); err != nil {
			return WrapInAtomicServiceError("failed to send additional data", err)
		}
		if err := s.transactionRepo.BatchSetOutcome(
			ctx, revertedSharedIDs, types.OutcomeSuccess,
		); err != nil {
			return WrapInAtomicServiceError("updating sender-reverted outcome", err)
		}
	}
	if len(signatureReverted) > 0 {
		if err := s.transactionRepo.BatchSetOutcome(
			ctx, signatureReverted, types.OutcomeReverted,
		); err != nil {
			return WrapInAtomicServiceError("updating signature-revert outcome", err)
		}
	}
	if len(signatureFailed) > 0 {
		if err := s.transactionRepo.BatchSetOutcome(
			ctx, signatureFailed, types.OutcomeFailed,
		); err != nil {
			return WrapInAtomicServiceError("updating signature-failed outcome", err)
		}
	}

	return nil
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (s *SignatureService) HandleSourceExecuted(ctx context.Context, sharedIDs []string) error {
	err := s.transactionRepo.BatchSetStateAndOutcome(ctx, sharedIDs, types.SourceFinalized, types.OutcomeSuccess)
	if err != nil {
		return WrapInAtomicServiceError("got error while trying to update transactions' state", err)
	}
	return nil
}

func (s *SignatureService) getTimestampForBlockHash(ctx context.Context, hash common.Hash) uint64 {
	ctxHeader, cancelHeader := context.WithTimeout(ctx, 10*time.Second)
	header, err := s.ethClient.HeaderByHash(ctxHeader, hash)
	cancelHeader()
	if err == nil {
		return header.Time
	} else {
		return 0
	}
}

func (s *SignatureService) generateMessages(
	signatures []types.CalldataSignature,
) []types.TxRequest {
	msgs := make([]types.TxRequest, 0, len(signatures))
	for _, sig := range signatures {
		calldata, err := s.generator.Generate(sig)
		if err != nil {
			slog.Warn("Failed to generate signature calldata, skipping",
				slog.String("shared_id", sig.SharedId),
				slog.String("signature_type", sig.SignatureType.String()),
				slog.Any("err", err),
			)
			continue
		}
		msgs = append(msgs, types.TxRequest{
			CorrelationID: sig.SharedId,
			MessageType:   messageTypeForSignature(sig.SignatureType),
			Address:       s.destEndpoint,
			Calldata:      calldata,
		})
	}
	return msgs
}

func messageTypeForSignature(st types.CallDataSignatureType) string {
	switch st {
	case types.UnlockOnDestinationSide:
		return "atomic.destination-unlock"
	case types.RevertOnDestinationSide:
		return "atomic.destination-revert"
	case types.RevertOnSenderSide:
		return "atomic.source-revert"
	default:
		return "atomic.unknown"
	}
}
