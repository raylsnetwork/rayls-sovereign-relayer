// Decommissioning Teleport (vanilla, atomic): atomic members below marked; shared/generic/Enygma/DVP retained.

package service

import (
	"context"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

const (
	// fetchTimeout is the context timeout for fetching cross-chain messages.
	fetchTimeout = 30 * time.Second
	// proofGenerationRoutineCount is the number of goroutines for parallel proof generation.
	proofGenerationRoutineCount = 8
)

//go:generate moq --pkg service_test -out crosschain_mock_test.go . TeleportClient CrossChainConsumer CrossChainEthereumClient CrossChainProofGenerator CrossChainTransactionRepository CrossChainSignatureRepository
type TeleportClient interface {
	StoreEncryptedDataBatch(
		ctx context.Context,
		sharedIDs []string,
		msgs []types.DispatchedMessageToPrivateHub,
		chainID *big.Int,
	) (common.Hash, error)
}

type CrossChainConsumer interface {
	Fetch(ctx context.Context, count int) ([]msgqueue.Message[CrossChainMessage], error)
}

type CrossChainEthereumClient interface {
	BlockByHash(ctx context.Context, hash common.Hash) (*ethTypes.Block, error)
}

type CrossChainProofGenerator interface {
	BatchGenerate(ctx context.Context, txHashes []common.Hash, routineCount int) [][]byte
}

type CrossChainTransactionRepository interface {
	BatchCreateWithStateAndOutcome(ctx context.Context, txs []types.Transaction, state types.TransactionState, outcome types.TransactionOutcome) error
}

type CrossChainSignatureRepository interface {
	BatchCreate(ctx context.Context, signatures []types.CalldataSignature) error
}

type CrossChainService struct {
	ticker    *time.Ticker
	myChainID *big.Int

	consumer CrossChainConsumer

	ethClient CrossChainEthereumClient
	proofGen  CrossChainProofGenerator

	teleportClient TeleportClient
	txRepo         CrossChainTransactionRepository
	// Decommissioning Teleport (vanilla, atomic).
	signatureRepo CrossChainSignatureRepository
}

func NewCrossChainService(
	tickerPeriod time.Duration,
	myChainID *big.Int,

	consumer CrossChainConsumer,

	ethClient CrossChainEthereumClient,
	proofGen CrossChainProofGenerator,

	teleportClient TeleportClient,
	txRepo CrossChainTransactionRepository,
	signatureRepo CrossChainSignatureRepository,
) *CrossChainService {
	return &CrossChainService{
		ticker:    time.NewTicker(tickerPeriod),
		myChainID: myChainID,

		consumer: consumer,

		ethClient: ethClient,
		proofGen:  proofGen,

		teleportClient: teleportClient,
		txRepo:         txRepo,
		signatureRepo:  signatureRepo,
	}
}

func (s *CrossChainService) Run( //nolint:gocognit // cross-chain message processing with multiple steps
	ctx context.Context,
) error {
	// service start log
	slog.Info("CrossChainService started", slog.String("chain_id", s.myChainID.String()))

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

		// tick log
		slog.Debug("CrossChainService tick")

		slog.Debug("Fetching cross-chain messages", slog.Int("batch_size", 100))
		fetchCtx, fetchCancel := context.WithTimeout(ctx, fetchTimeout)
		msgs, err := s.consumer.Fetch(fetchCtx, 100)
		fetchCancel()
		if err != nil {
			slog.Warn("Failed to fetch cross-chain messages", slog.Any("error", err))
			continue
		}
		if len(msgs) == 0 {
			slog.Debug("No cross-chain messages available")
			continue
		}

		slog.Debug("Fetched cross-chain messages", slog.Int("count", len(msgs)))

		var txHashes []common.Hash
		for _, msg := range msgs {
			txHashes = append(txHashes, msg.V.TxHash)
		}

		slog.Debug(
			"Generating cross-chain proofs",
			slog.Int("count", len(txHashes)),
			slog.Int("routine_count", proofGenerationRoutineCount),
		)
		proofs := s.proofGen.BatchGenerate(ctx, txHashes, proofGenerationRoutineCount)

		sharedIDToMessage := make(map[string]msgqueue.Message[CrossChainMessage])

		batchIDs := make(map[uint64]string)
		dispatchBatches := make(map[uint64][]types.DispatchedMessageToPrivateHub)

		// A missing proof must never be dispatched — it strands the token (the
		// destination rejects a nil proof; see #242). Record these teleports as
		// proof_invalid and ack them, instead of discarding them untracked.
		proofFailedTxs := []types.Transaction{}
		proofFailedMsgs := []msgqueue.Message[CrossChainMessage]{}

		// Decommissioning Teleport (vanilla, atomic).
		revertSignatures := []types.CalldataSignature{}
		for i, msg := range msgs {
			if len(proofs[i]) == 0 {
				// batchGenerateRoutine logged the cause; this adds message-level
				// context for the dispatch decision (two ERROR lines is intentional).
				slog.Error(
					"Cross-chain proof generation failed; recording proof_invalid instead of dispatching",
					slog.Any("tx_hash", msg.V.TxHash),
					slog.String("message_id", common.Hash(msg.V.MessageID).Hex()),
					slog.String("id", msg.V.ID),
				)
				proofFailedTxs = append(proofFailedTxs, newProofFailedTransaction(msg.V, s.myChainID))
				proofFailedMsgs = append(proofFailedMsgs, msg)
				continue
			}

			chainIDUint64 := msg.V.ToChainID.Uint64()

			if _, ok := batchIDs[chainIDUint64]; !ok {
				batchIDs[chainIDUint64] = uuid.New().String()
				slog.Debug(
					"Created new dispatch batch",
					slog.Uint64("to_chain_id", chainIDUint64),
					slog.String("batch_id", batchIDs[chainIDUint64]),
				)
			}
			batchID := batchIDs[chainIDUint64]

			var block *ethTypes.Block
			block, err = s.ethClient.BlockByHash(ctx, msg.V.BlockHash)
			if err != nil {
				slog.Warn("Failed to get block by hash", slog.Any("hash", msg.V.BlockHash), slog.Any("error", err))
				continue
			}

			messageToDispatch := newDispatchedMessage(batchID, msg.V, s.myChainID, block, proofs[i])
			dispatchBatches[chainIDUint64] = append(dispatchBatches[chainIDUint64], messageToDispatch)

			sharedIDToMessage[messageToDispatch.SharedId] = msg

			// Decommissioning Teleport (vanilla, atomic).
			if isAtomicMessage(msg.V) {
				revertSignatures = append(revertSignatures, newSignature(messageToDispatch.SharedId, msg.V))
			}
		}

		// Persist the proof_invalid failures before acking. On a persist error,
		// leave them un-acked so the queue redelivers and we retry.
		if len(proofFailedTxs) != 0 {
			if createErr := s.txRepo.BatchCreateWithStateAndOutcome(
				ctx,
				proofFailedTxs,
				types.SourcePublish,
				types.OutcomeFailed,
			); createErr == nil {
				slog.Warn(
					"Recorded proof-invalid cross-chain transactions (proof generation failed)",
					slog.Int("count", len(proofFailedTxs)),
				)
				for _, failedMsg := range proofFailedMsgs {
					_ = failedMsg.Ack(ctx)
				}
			} else {
				slog.Error(
					"Failed to persist proof-invalid cross-chain transactions; leaving un-acked for redelivery",
					slog.Int("count", len(proofFailedTxs)),
					slog.Any("error", createErr),
				)
			}
		}

		// Decommissioning Teleport (vanilla, atomic).
		if len(revertSignatures) != 0 {
			if batchErr := s.signatureRepo.BatchCreate(ctx, revertSignatures); batchErr != nil {
				slog.Error(
					"Failed to persist revert signatures",
					slog.Int("count", len(revertSignatures)),
					slog.Any("error", batchErr),
				)
				continue
			}
		}

		for chainID, batch := range dispatchBatches {
			sharedIDs := make([]string, len(batch))
			for i, msg := range batch {
				sharedIDs[i] = msg.SharedId
			}
			slog.Info(
				"Sending encrypted data batch to Teleport",
				slog.Uint64("to_chain_id", chainID),
				slog.Int("batch_size", len(batch)),
			)
			var batchHash common.Hash
			batchHash, err = s.teleportClient.StoreEncryptedDataBatch(ctx, sharedIDs, batch, new(big.Int).SetUint64(chainID))
			if err == nil {
				slog.Info(
					"Successfully stored encrypted data batch to Teleport",
					slog.Uint64("to_chain_id", chainID),
					slog.Int("batch_size", len(batch)),
				)

				txs := dispatchedMessagesToTransactions(batch, batchHash)
				slog.Debug("Persisting successful cross-chain transaction batch", slog.Int("count", len(txs)))

				if createErr := s.txRepo.BatchCreateWithStateAndOutcome(
					ctx,
					txs,
					types.SourcePublish,
					types.OutcomeSuccess,
				); createErr == nil {
					slog.Info("Persisted successful cross-chain transaction batch", slog.Int("count", len(txs)))

					slog.Debug("Acknowledging successful cross-chain message batch", slog.Int("count", len(batch)))

					for _, dispatchedMsg := range batch {
						// TODO: this could be a critical section - what happens if the service
						// fails and only half the messages of the batch are acked?
						_ = sharedIDToMessage[dispatchedMsg.SharedId].Ack(ctx)
					}
				} else {
					slog.Error(
						"Failed to persist successful cross-chain transactions",
						slog.Int("count", len(txs)),
						slog.Any("error", createErr),
					)
				}
			} else {
				slog.Error(
					"Failed to store encrypted data batch to Teleport",
					slog.Uint64("to_chain_id", chainID),
					slog.Int("batch_size", len(batch)),
					slog.Any("error", err),
				)

				txs := dispatchedMessagesToTransactions(batch, common.Hash{})
				slog.Debug("Persisting failed cross-chain transaction batch", slog.Int("count", len(txs)))

				if createErr := s.txRepo.BatchCreateWithStateAndOutcome(
					ctx,
					txs,
					types.SourcePublish,
					types.OutcomeFailed,
				); createErr == nil {
					slog.Info("Persisted failed cross-chain transaction batch", slog.Int("count", len(txs)))

					slog.Debug("Acknowledging failed cross-chain message batch", slog.Int("count", len(batch)))

					for _, dispatchedMsg := range batch {
						// TODO: this could be a critical section - what happens if the service
						// fails and only half the messages of the batch are acked?
						_ = sharedIDToMessage[dispatchedMsg.SharedId].Ack(ctx)
					}
				} else {
					slog.Error(
						"Failed to persist failed cross-chain transactions",
						slog.Int("count", len(txs)),
						slog.Any("error", createErr),
					)
				}

			}
		}
	}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func isAtomicMessage(msg CrossChainMessage) bool {
	return len(msg.Data.MessageMetadata.RevertPayloadDataSender) != 0 &&
		len(msg.Data.MessageMetadata.RevertPayloadDataReceiver) != 0 &&
		len(msg.Data.MessageMetadata.LockData) != 0
}

func newDispatchedMessage(
	batchID string,
	msg CrossChainMessage,
	fromChainID *big.Int,
	block *ethTypes.Block,
	proof []byte,
) types.DispatchedMessageToPrivateHub {
	return types.DispatchedMessageToPrivateHub{
		SharedId: uuid.New().String(),
		BatchId:  batchID,

		MessageId:   msg.MessageID,
		FromChainId: fromChainID,
		From:        msg.From,
		ToChainId:   msg.ToChainID,
		To:          msg.To,
		Data:        msg.Data,

		TransactionType: types.Transfer,
		IsAtomic:        isAtomicMessage(msg), // maybe move this check else

		BlockNumber: msg.BlockNumber,
		BlockHash:   msg.BlockHash,
		LogIdx:      msg.LogIdx,

		ParentHash: block.ParentHash().Hex(),

		TxHashSource:          msg.TxHash,
		TxHashSourceStatus:    1,
		TxHashSourceTimestamp: block.Time(),

		ResourceId:   msg.Data.MessageMetadata.ResourceId,
		TokenAddress: msg.Data.MessageMetadata.TransferMetadata.TokenAddress,

		Proofs:      proof,
		TxLocation:  msg.TxIdx,
		TxTrieProof: block.TxHash(),
	}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func newSignature(sharedID string, msg CrossChainMessage) types.CalldataSignature {
	return types.CalldataSignature{
		SharedId:      sharedID,
		Signature:     msg.Data.MessageMetadata.RevertPayloadDataSender,
		ResourceId:    msg.Data.MessageMetadata.ResourceId,
		SignatureType: types.RevertOnSenderSide,
	}
}

func dispatchedMessagesToTransactions(
	msgs []types.DispatchedMessageToPrivateHub,
	batchHash common.Hash,
) []types.Transaction {
	txs := make([]types.Transaction, 0, len(msgs))
	for _, msg := range msgs {
		txs = append(txs, dispatchedMessageToTransaction(msg, batchHash))
	}
	return txs
}

func dispatchedMessageToTransaction(msg types.DispatchedMessageToPrivateHub, batchHash common.Hash) types.Transaction {
	return types.Transaction{
		SharedID:            msg.SharedId,
		BatchID:             msg.BatchId,
		BatchPrivateHubHash: batchHash,

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

		UpdatedAt: time.Unix(int64(msg.TxHashSourceTimestamp), 0), //nolint:gosec // timestamps fit in int64

		TransferID:     msg.Data.MessageMetadata.TransferMetadata.Id.String(),
		TransferAmount: msg.Data.MessageMetadata.TransferMetadata.Amount.String(),
	}
}

// newProofFailedTransaction builds a queryable proof_invalid record for a teleport
// whose inclusion proof could not be generated, so a stuck teleport can be found
// (WHERE proof_invalid) instead of lost untracked. State and outcome are set by the
// caller (BatchCreateWithStateAndOutcome); it is built from the source message
// alone, without the block.
func newProofFailedTransaction(msg CrossChainMessage, fromChainID *big.Int) types.Transaction {
	return types.Transaction{
		SharedID:     uuid.New().String(),
		ProofInvalid: true,

		MsgID:       msg.MessageID,
		FromChainID: fromChainID,
		ToChainID:   msg.ToChainID,

		FromContractAddress: msg.Data.MessageMetadata.TransferMetadata.TokenAddress.String(),
		FromUserAddress:     msg.From.String(),
		ResourceID:          common.Hash(msg.Data.MessageMetadata.ResourceId).String(),
		IsAtomic:            isAtomicMessage(msg),

		BlockNumber: msg.BlockNumber,
		TxHash:      msg.TxHash.String(),
		LogIndex:    msg.LogIdx,

		UpdatedAt: time.Now(),

		TransferID:     msg.Data.MessageMetadata.TransferMetadata.Id.String(),
		TransferAmount: msg.Data.MessageMetadata.TransferMetadata.Amount.String(),
	}
}
