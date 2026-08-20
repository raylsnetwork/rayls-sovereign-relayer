// Decommissioning Teleport (vanilla, atomic): atomic members below marked; shared/generic/Enygma/DVP retained.

package service

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/raylsnetwork/rayls-sovereign-relayer/batcher"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/EndpointV1"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cryptography/proofs"
	"github.com/raylsnetwork/rayls-sovereign-relayer/msgqueue"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

//go:generate moq --pkg service_test -out crosschain_mock_test.go . CrossChainMQ TransactionRepository SignatureRepository TransactionGenerator Batcher ReceiptHandler EndpointClient DeployerClient

type CrossChainMQ interface {
	Fetch(ctx context.Context, count int) ([]msgqueue.Message[types.DispatchedMessageToPrivateHub], error)
}

type SignatureRepository interface {
	BatchCreate(ctx context.Context, signatures []types.CalldataSignature) error
}

type TransactionRepository interface {
	BatchCreate(ctx context.Context, txs []types.Transaction) error
	BatchUpdateDestinationHashForSharedIDs(ctx context.Context, hashBySharedID map[string]common.Hash) error
}

type TransactionGenerator interface {
	Generate(
		fromChainID *big.Int,
		fromAddress, toAddress common.Address,
		message EndpointV1.RaylsMessage,
		id common.Hash,
	) ([]byte, error)
}

// Batcher is the fire-and-forget publisher the service hands its generated
// messages to. A concrete *batcher.Batcher wired against the crosschain
// atomic/vanilla subjects satisfies it.
type Batcher interface {
	Send(ctx context.Context, msgs []batcher.Message) error
}

// ReceiptHandler is the subset of *ReceiptService / *VanillaReceiptService
// that the crosschain callbacks delegate to. CrossChainService holds two
// separate instances — one for atomic, one for vanilla — picked per result
// type.
//
// Optional extension: an implementation MAY also satisfy the unexported
// lostMinedHandler interface (see dispatchResults, where the lost-vs-revert
// rationale lives) to handle lost-on-chain outcomes separately. Kept off this
// base interface so ReceiptHandler stays narrow. The vanilla path implements
// it; the atomic path deliberately does not (its compensation machinery
// already covers chain-loss).
type ReceiptHandler interface {
	HandleSuccessfullyMined(ctx context.Context, sharedIDs []string) error
	HandleFailedMined(ctx context.Context, sharedIDs []string) error
}

type EndpointClient interface {
	GetResourceAddress(ctx context.Context, resourceId string) (common.Address, error)
}

type DeployerClient interface {
	DeployResourceAndExecute(
		ctx context.Context,
		resourceId [32]byte,
		message *types.DispatchedMessageToPrivateHub,
	) (common.Hash, error)
}

type VerifierFunc func(mashaledProof []byte, rootTxHash common.Hash, txIdx uint) error

type CrossChainService struct {
	plEndpointAddress common.Address

	consumer       CrossChainMQ
	txRepo         TransactionRepository
	signatureRepo  SignatureRepository
	txGen          TransactionGenerator
	endpointClient EndpointClient
	deployerClient DeployerClient

	// Decommissioning Teleport (vanilla, atomic): atomicBatcher only; vanillaBatcher is retained (generic non-atomic).
	atomicBatcher  Batcher
	vanillaBatcher Batcher

	// Decommissioning Teleport (vanilla, atomic): atomicReceiptSvc only; vanillaReceiptSvc is retained (generic non-atomic).
	atomicReceiptSvc  ReceiptHandler
	vanillaReceiptSvc ReceiptHandler

	verify VerifierFunc

	batchSize int
	ticker    *time.Ticker
}

func NewCrossChainService(
	plEndpointAddress common.Address,
	batchSize int,

	consumer CrossChainMQ,
	txRepo TransactionRepository,
	signatureRepo SignatureRepository,
	txGen TransactionGenerator,
	atomicBatcher Batcher,
	vanillaBatcher Batcher,
	atomicReceiptSvc ReceiptHandler,
	vanillaReceiptSvc ReceiptHandler,
	endpointClient EndpointClient,
	deployerClient DeployerClient,
) *CrossChainService {
	return NewCrossChainServiceWith(
		plEndpointAddress,
		batchSize,
		consumer,
		txRepo,
		signatureRepo,
		txGen,
		atomicBatcher,
		vanillaBatcher,
		atomicReceiptSvc,
		vanillaReceiptSvc,
		endpointClient,
		deployerClient,
		verifyProofs,
	)
}

func NewCrossChainServiceWith(
	plEndpointAddress common.Address,
	batchSize int,

	consumer CrossChainMQ,
	txRepo TransactionRepository,
	signatureRepo SignatureRepository,
	txGen TransactionGenerator,
	atomicBatcher Batcher,
	vanillaBatcher Batcher,
	atomicReceiptSvc ReceiptHandler,
	vanillaReceiptSvc ReceiptHandler,
	endpointClient EndpointClient,
	deployerClient DeployerClient,

	verify VerifierFunc,
) *CrossChainService {
	return &CrossChainService{
		plEndpointAddress: plEndpointAddress,

		consumer:       consumer,
		txRepo:         txRepo,
		signatureRepo:  signatureRepo,
		txGen:          txGen,
		endpointClient: endpointClient,
		deployerClient: deployerClient,

		atomicBatcher:  atomicBatcher,
		vanillaBatcher: vanillaBatcher,

		atomicReceiptSvc:  atomicReceiptSvc,
		vanillaReceiptSvc: vanillaReceiptSvc,

		verify:    verify,
		batchSize: batchSize,
		ticker:    time.NewTicker(time.Second),
	}
}

// messageOutcome is what processMessage accumulates for one fetched
// message. Any combination of the optional fields may be unset:
//
//   - batchMsg is nil when the message took the deploy path (deploy IS
//     the tx, no async broadcast) or when calldata generation failed.
//   - signatures is empty for non-atomic messages.
//   - deployResult is non-nil only when the deploy path was taken; its
//     Kind discriminates success vs failure.
type messageOutcome struct {
	tx           types.Transaction
	atomic       bool
	batchMsg     *batcher.Message
	signatures   []types.CalldataSignature
	deployResult *types.TxResult
}

func (s *CrossChainService) Run(ctx context.Context) error {
	slog.Info("CrossChainService started")
	initialRun := make(chan struct{}, 1)
	initialRun <- struct{}{}

	for {
		select {
		case <-ctx.Done():
			slog.Info("CrossChainService shutting down")
			return nil
		case <-s.ticker.C:
		case <-initialRun:
		}

		if err := s.tick(ctx); err != nil {
			// tick already logs; the returned error is a shutdown signal
			// (context cancelled). Any other per-iteration failure is
			// swallowed and the next tick retries.
			if ctx.Err() != nil {
				slog.Info("CrossChainService shutting down")
				return nil //nolint:nilerr // context cancellation is a clean shutdown signal
			}
		}
	}
}

// tick processes one batch: fetch, per-message handling, publish, persist,
// dispatch synchronous deploy results, ack. Returns a non-nil error only
// when the outer loop should exit (ctx cancellation). All other failures
// are logged and swallowed.
func (s *CrossChainService) tick(ctx context.Context) error {
	slog.Debug("Fetching cross-chain messages", slog.Int("batch_size", s.batchSize))
	msgSlice, err := s.consumer.Fetch(ctx, s.batchSize)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		slog.Error("Failed to get next message from CrossChain MQ", slog.Any("error", err))
		return nil
	}
	if len(msgSlice) == 0 {
		return nil
	}

	slog.Debug("Fetched cross-chain messages", slog.Int("count", len(msgSlice)))

	deployedResourceIDs := map[string]bool{}
	invalidResourceIDs := map[string]bool{}
	s.populateResourceIDMaps(ctx, msgSlice, deployedResourceIDs, invalidResourceIDs)
	slog.Debug("Resource deployment check complete",
		slog.Int("deployed_count", len(deployedResourceIDs)),
		slog.Int("invalid_count", len(invalidResourceIDs)))

	var (
		txSlice              []types.Transaction
		sigSlice             []types.CalldataSignature
		atomicBatchMsgs      []batcher.Message
		vanillaBatchMsgs     []batcher.Message
		atomicDeployResults  []types.TxResult
		vanillaDeployResults []types.TxResult
	)

	for _, msg := range msgSlice {
		outcome := s.processMessage(ctx, msg.V, deployedResourceIDs, invalidResourceIDs)
		if outcome == nil {
			continue
		}

		txSlice = append(txSlice, outcome.tx)
		sigSlice = append(sigSlice, outcome.signatures...)

		if outcome.batchMsg != nil {
			if outcome.atomic {
				atomicBatchMsgs = append(atomicBatchMsgs, *outcome.batchMsg)
			} else {
				vanillaBatchMsgs = append(vanillaBatchMsgs, *outcome.batchMsg)
			}
		}

		if outcome.deployResult != nil {
			if outcome.atomic {
				atomicDeployResults = append(atomicDeployResults, *outcome.deployResult)
			} else {
				vanillaDeployResults = append(vanillaDeployResults, *outcome.deployResult)
			}
		}
	}

	if err := s.publishBatches(ctx, atomicBatchMsgs, vanillaBatchMsgs); err != nil {
		return nil // already logged; skip ack, retry next tick
	}

	if len(sigSlice) != 0 {
		if err := s.signatureRepo.BatchCreate(ctx, sigSlice); err != nil {
			slog.Error("Failed to persist signatures", slog.Any("error", err))
			return nil
		}
		slog.Debug("Successfully persisted signatures", slog.Int("count", len(sigSlice)))
	}

	if len(txSlice) != 0 {
		if err := s.txRepo.BatchCreate(ctx, txSlice); err != nil {
			slog.Error("Failed to persist transactions", slog.Any("error", err))
			return nil
		}
		slog.Debug("Successfully persisted transactions", slog.Int("count", len(txSlice)))
	}

	// Synthesised deploy-path results flow through the exact same
	// callbacks as async NATS results — uniform handling.
	if len(atomicDeployResults) > 0 {
		if err := s.HandleAtomicResults(ctx, atomicDeployResults); err != nil {
			slog.Error("Failed to dispatch atomic deploy results", slog.Any("error", err))
		}
	}
	if len(vanillaDeployResults) > 0 {
		if err := s.HandleVanillaResults(ctx, vanillaDeployResults); err != nil {
			slog.Error("Failed to dispatch vanilla deploy results", slog.Any("error", err))
		}
	}

	for _, msg := range msgSlice {
		_ = msg.Ack(ctx)
	}
	return nil
}

// processMessage runs the full per-message pipeline: proof verify, atomic
// signature emission, resource-deploy check and calldata generation.
// Returns nil when the message is skipped (verify failed).
func (s *CrossChainService) processMessage(
	ctx context.Context,
	msg types.DispatchedMessageToPrivateHub,
	deployedResourceIDs, invalidResourceIDs map[string]bool,
) *messageOutcome {
	slog.Debug("Processing cross-chain message", slog.String("shared_id", msg.SharedId))
	if verifyErr := s.verify(msg.Proofs, msg.TxTrieProof, msg.TxLocation); verifyErr != nil {
		slog.Warn("Invalid message proof",
			slog.String("shared_id", msg.SharedId),
			slog.Any("error", verifyErr),
		)
		return nil
	}

	outcome := &messageOutcome{
		tx:     dispatchedMessageToTransaction(msg),
		atomic: msg.IsAtomic,
	}

	// Decommissioning Teleport (vanilla, atomic).
	if msg.IsAtomic {
		outcome.signatures = []types.CalldataSignature{
			signatureFromMessage(msg, types.UnlockOnDestinationSide),
			signatureFromMessage(msg, types.RevertOnDestinationSide),
		}
	}

	// Resource deployment branch — only for messages that carry a
	// resource ID (arbitrary messages pass through directly).
	if msg.ResourceId != (common.Hash{}) {
		if _, invalid := invalidResourceIDs[msg.ResourceId.String()]; invalid {
			slog.Debug("Skipping message with invalid resource ID",
				slog.String("shared_id", msg.SharedId),
				slog.String("resource_id", msg.ResourceId.String()))
			outcome.tx.State = types.DestinationDispatch
			outcome.tx.Outcome = types.OutcomeFailed
			return outcome
		}

		if _, deployed := deployedResourceIDs[msg.ResourceId.String()]; !deployed {
			slog.Info("Resource not deployed, attempting deployment",
				slog.String("resource_id", msg.ResourceId.String()),
				slog.String("shared_id", msg.SharedId))

			txHash, deployErr := s.deployContractAndExecuteTransaction(
				ctx, msg, deployedResourceIDs, invalidResourceIDs,
			)
			outcome.deployResult = synthesizeDeployResult(msg.SharedId, txHash, deployErr)
			if deployErr != nil {
				outcome.tx.State = types.DestinationDispatch
				outcome.tx.Outcome = types.OutcomeFailed
			} else {
				outcome.tx.TxHashDestination = txHash
				outcome.tx.State = types.DestinationDispatch
				outcome.tx.Outcome = types.OutcomePending
			}
			return outcome
		}
	}

	calldata, err := s.txGen.Generate(msg.FromChainId, msg.From, msg.To, msg.Data, msg.MessageId)
	if err != nil {
		slog.Warn("Failed to generate transaction calldata, skipping message",
			slog.Any("error", err),
			slog.String("shared_id", msg.SharedId),
		)
		outcome.tx.State = types.DestinationDispatch
		outcome.tx.Outcome = types.OutcomeFailed
		return outcome
	}

	outcome.batchMsg = &batcher.Message{
		ID:       msg.SharedId,
		Address:  s.plEndpointAddress,
		Calldata: calldata,
	}
	outcome.tx.State = types.DestinationDispatch
	outcome.tx.Outcome = types.OutcomePending
	return outcome
}

// publishBatches sends the atomic and vanilla batches. Errors are logged
// per-batch; returns a non-nil error only if the caller should skip the
// rest of the tick (publish failure).
func (s *CrossChainService) publishBatches(
	ctx context.Context,
	atomicBatchMsgs, vanillaBatchMsgs []batcher.Message,
) error {
	if len(atomicBatchMsgs) != 0 {
		slog.Debug("Publishing atomic batch", slog.Int("count", len(atomicBatchMsgs)))
		if err := s.atomicBatcher.Send(ctx, atomicBatchMsgs); err != nil {
			slog.Error("Failed to publish atomic batch", slog.Any("error", err))
			return err
		}
		slog.Info("Published atomic batch", slog.Int("count", len(atomicBatchMsgs)))
	}
	if len(vanillaBatchMsgs) != 0 {
		slog.Debug("Publishing vanilla batch", slog.Int("count", len(vanillaBatchMsgs)))
		if err := s.vanillaBatcher.Send(ctx, vanillaBatchMsgs); err != nil {
			slog.Error("Failed to publish vanilla batch", slog.Any("error", err))
			return err
		}
		slog.Info("Published vanilla batch", slog.Int("count", len(vanillaBatchMsgs)))
	}
	return nil
}

// HandleAtomicResults consumes terminal TxResults for atomic crosschain
// messages (including synthesised deploy-path results). Splits by Kind,
// updates the destination hash for successes, and delegates to the
// atomic receipt service.
//
// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (s *CrossChainService) HandleAtomicResults(ctx context.Context, results []types.TxResult) error {
	return s.dispatchResults(ctx, results, s.atomicReceiptSvc)
}

// HandleVanillaResults mirrors HandleAtomicResults for the non-atomic flow.
func (s *CrossChainService) HandleVanillaResults(ctx context.Context, results []types.TxResult) error {
	return s.dispatchResults(ctx, results, s.vanillaReceiptSvc)
}

// dispatchResults is the shared body of the two callbacks.
//
// CTS emits three terminal kinds: TxResultSuccess, TxResultRevert (mined
// and reverted on-chain), and TxResultFailed (never mined — reaper
// dead-letter after MaxResendAttempts, or a synthesised deploy failure).
// The two failure kinds are NOT semantically equivalent: a real revert
// is terminal and must be persisted as OutcomeReverted, whereas a lost
// tx never executed and must NOT be persisted as OutcomeReverted (that
// collapses the lost-tx case into the real-revert state and closes the
// retry path). Routing splits them via the optional HandleLostMined
// hook; handlers that don't expose it (e.g. atomic) fall back to the
// legacy collapsed behaviour without an interface break.
func (s *CrossChainService) dispatchResults(
	ctx context.Context,
	results []types.TxResult,
	receiptSvc interface {
		HandleSuccessfullyMined(ctx context.Context, sharedIDs []string) error
		HandleFailedMined(ctx context.Context, sharedIDs []string) error
	},
) error {
	if len(results) == 0 {
		return nil
	}

	type lostMinedHandler interface {
		HandleLostMined(ctx context.Context, sharedIDs []string) error
	}

	var (
		successIDs  []string
		revertedIDs []string
		lostIDs     []string
	)
	successHashes := map[string]common.Hash{}

	for _, res := range results {
		switch res.Kind {
		case types.TxResultSuccess:
			successIDs = append(successIDs, res.CorrelationID)
			successHashes[res.CorrelationID] = res.TxHash
		case types.TxResultRevert:
			revertedIDs = append(revertedIDs, res.CorrelationID)
		case types.TxResultFailed:
			lostIDs = append(lostIDs, res.CorrelationID)
		}
	}

	if len(successHashes) > 0 {
		if err := s.txRepo.BatchUpdateDestinationHashForSharedIDs(ctx, successHashes); err != nil {
			slog.Error("Failed to update destination hashes", slog.Any("error", err))
		}
	}

	if len(successIDs) > 0 {
		if err := receiptSvc.HandleSuccessfullyMined(ctx, successIDs); err != nil {
			slog.Error("receipt service HandleSuccessfullyMined failed", slog.Any("error", err))
		}
	}

	if h, ok := any(receiptSvc).(lostMinedHandler); ok {
		// Split path: handler exposes a dedicated lost-mined hook.
		if len(revertedIDs) > 0 {
			if err := receiptSvc.HandleFailedMined(ctx, revertedIDs); err != nil {
				slog.Error("receipt service HandleFailedMined failed", slog.Any("error", err))
			}
		}
		if len(lostIDs) > 0 {
			if err := h.HandleLostMined(ctx, lostIDs); err != nil {
				slog.Error("receipt service HandleLostMined failed", slog.Any("error", err))
			}
		}
	} else if combined := append(revertedIDs, lostIDs...); len(combined) > 0 {
		// Fallback path: handlers without a separate lost-mined hook
		// (e.g. the atomic flow, which folds chain-loss into its own
		// compensation machinery) keep the legacy collapsed behaviour
		// — one HandleFailedMined call carrying both kinds.
		if err := receiptSvc.HandleFailedMined(ctx, combined); err != nil {
			slog.Error("receipt service HandleFailedMined failed", slog.Any("error", err))
		}
	}
	return nil
}

func (s *CrossChainService) populateResourceIDMaps(
	ctx context.Context,
	msgSlice []msgqueue.Message[types.DispatchedMessageToPrivateHub],
	deployedResourceIDs, invalidResourceIDs map[string]bool,
) {
	uniqueResourceIDs := map[string]bool{}
	for _, msg := range msgSlice {
		if msg.V.ResourceId == (common.Hash{}) {
			continue
		}
		uniqueResourceIDs[msg.V.ResourceId.String()] = true
	}

	for resourceID := range uniqueResourceIDs {
		resourceAddr, err := s.endpointClient.GetResourceAddress(ctx, resourceID[2:])
		if err != nil {
			slog.Warn("Failed to check resource deployment status",
				slog.Any("error", err),
				slog.String("resource_id", resourceID),
			)
			invalidResourceIDs[resourceID] = true
			continue
		}

		if resourceAddr != (common.Address{}) {
			deployedResourceIDs[resourceID] = true
		}
	}
}

// deployContractAndExecuteTransaction runs a deploy-and-execute on the
// deployer client and updates the caller's deployed/invalid maps. The
// deploy itself is the transaction — the returned hash is the on-chain
// tx hash. Errors are returned so the caller can synthesise a
// TxResultFailed to route through the callback.
func (s *CrossChainService) deployContractAndExecuteTransaction(
	ctx context.Context,
	msg types.DispatchedMessageToPrivateHub,
	deployedResourceIDs, invalidResourceIDs map[string]bool,
) (common.Hash, error) {
	txHash, err := s.deployerClient.DeployResourceAndExecute(ctx, msg.Data.MessageMetadata.ResourceId, &msg)
	if err != nil {
		slog.Error("Failed to deploy resource",
			slog.Any("error", err),
			slog.String("shared_id", msg.SharedId),
			slog.String("resource_id", msg.ResourceId.String()))
		invalidResourceIDs[msg.ResourceId.String()] = true
		return common.Hash{}, err
	}

	slog.Info("Successfully deployed resource",
		slog.String("resource_id", msg.ResourceId.String()),
		slog.String("shared_id", msg.SharedId))
	deployedResourceIDs[msg.ResourceId.String()] = true
	return txHash, nil
}

func synthesizeDeployResult(sharedID string, txHash common.Hash, err error) *types.TxResult {
	if err != nil {
		return &types.TxResult{
			CorrelationID: sharedID,
			Kind:          types.TxResultFailed,
			ErrorReason:   err.Error(),
		}
	}
	return &types.TxResult{
		CorrelationID: sharedID,
		Kind:          types.TxResultSuccess,
		TxHash:        txHash,
	}
}

func dispatchedMessageToTransaction(msg types.DispatchedMessageToPrivateHub) types.Transaction {
	return types.Transaction{
		SharedID:            msg.SharedId,
		BatchID:             msg.BatchId,
		BatchPrivateHubHash: msg.BatchPrivateHubHash,

		MsgID:               msg.MessageId,
		FromChainID:         msg.FromChainId,
		FromContractAddress: msg.TokenAddress.String(),
		FromUserAddress:     msg.From.String(),
		ToChainID:           msg.ToChainId,

		ResourceID: msg.ResourceId.String(),
		IsAtomic:   msg.IsAtomic,

		ParentHash:  msg.ParentHash,
		BlockNumber: msg.BlockNumber,
		TxHash:      msg.TxHashSource.String(),
		LogIndex:    msg.LogIdx,

		//nolint:gosec // block timestamps are within int64 range
		UpdatedAt: time.Unix(int64(msg.TxHashSourceTimestamp), 0),

		TransferID:     msg.Data.MessageMetadata.TransferMetadata.Id.String(),
		TransferAmount: msg.Data.MessageMetadata.TransferMetadata.Amount.String(),
	}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func signatureFromMessage(
	msg types.DispatchedMessageToPrivateHub,
	sigType types.CallDataSignatureType,
) types.CalldataSignature {
	var data []byte
	switch sigType {
	case types.UnlockOnDestinationSide:
		data = msg.Data.MessageMetadata.LockData
	case types.RevertOnDestinationSide:
		data = msg.Data.MessageMetadata.RevertPayloadDataReceiver
	default:
		return types.CalldataSignature{}
	}
	return types.CalldataSignature{
		SignatureType: sigType,
		Signature:     data,
		SharedId:      msg.SharedId,
		ResourceId:    msg.Data.MessageMetadata.ResourceId,
	}
}

func verifyProofs(marshaledProof []byte, rootTxHash common.Hash, txIdx uint) error {
	if len(marshaledProof) == 0 {
		return fmt.Errorf("proof is nil or empty: source proof generation failed")
	}

	key, err := rlp.EncodeToBytes(txIdx)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("failed to encode transaction index in rlp: %w", err))
	}

	proofDB, err := proofs.Import(marshaledProof)
	if err != nil {
		return fmt.Errorf("failed to deserialize proof (len=%d): %w", len(marshaledProof), err)
	}

	_, err = proofs.VerifyProof(rootTxHash, key, proofDB)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("failed to verify proofs: %w", err))
	}
	return nil
}
