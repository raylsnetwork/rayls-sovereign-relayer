// Decommissioning Teleport (vanilla, atomic).

// Package service implements the legacy public-chain (RN) Teleport bridge relayer.
//
// Deprecated: Decommissioning Teleport (vanilla, atomic).
package service

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/batcher"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/RNMessageDispatcherV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/faultinjector"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

//go:generate moq --pkg service_test -out generator_mock_test.go . MessageConsumer TransactionGenerator Batcher RevertSignatureRepository MessageRecordRepository

const fetchTimeout = 30 * time.Second

type MessageConsumer interface {
	Fetch(ctx context.Context, count int) ([]msgqueue.Message[Message], error)
}

// Batcher is the fire-and-forget publisher the generator hands its generated
// forward / revert messages to. A concrete *batcher.Batcher wired against the
// publicrelayer.forward / publicrelayer.revert subjects satisfies it.
type Batcher interface {
	Send(ctx context.Context, msgs []batcher.Message) error
}

type RevertSignatureRepository interface {
	BatchCreate(ctx context.Context, sigs []RevertSignature) error
	GetByIDs(ctx context.Context, ids []string) ([]RevertSignature, error)
}

type MessageRecordRepository interface {
	BatchCreate(ctx context.Context, records []MessageRecord) error
	UpdateForwardResults(ctx context.Context, updates []ForwardResultUpdate) error
	UpdateRevertResults(ctx context.Context, updates []RevertResultUpdate) error
}

type TransactionGenerator interface {
	Generate(
		fromAddress, toAddress common.Address,
		message RNMessageDispatcherV1.RaylsNodeMessage,
		id common.Hash,
	) ([]byte, error)
}

type GeneratorServiceConfig struct {
	Interval              time.Duration
	EndpointAddress       common.Address
	SourceEndpointAddress common.Address
}

type GeneratorService struct {
	relayTicker           *time.Ticker
	endpointAddress       common.Address
	sourceEndpointAddress common.Address

	consumer  MessageConsumer
	generator TransactionGenerator

	forwardBatcher    Batcher
	revertBatcher     Batcher
	revertRepository  RevertSignatureRepository
	messageRecordRepo MessageRecordRepository
}

func NewGeneratorService(
	config GeneratorServiceConfig,
	consumer MessageConsumer,
	generator TransactionGenerator,
	forwardBatcher Batcher,
	revertBatcher Batcher,
	revertRepository RevertSignatureRepository,
	messageRecordRepo MessageRecordRepository,
) *GeneratorService {
	return &GeneratorService{
		relayTicker:           time.NewTicker(config.Interval),
		endpointAddress:       config.EndpointAddress,
		sourceEndpointAddress: config.SourceEndpointAddress,

		consumer:  consumer,
		generator: generator,

		forwardBatcher:    forwardBatcher,
		revertBatcher:     revertBatcher,
		revertRepository:  revertRepository,
		messageRecordRepo: messageRecordRepo,
	}
}

func (s *GeneratorService) Run(ctx context.Context) error {
	for {
		s.processOnce(ctx)

		select {
		case <-ctx.Done():
			return nil
		case <-s.relayTicker.C:
		}
	}
}

func (s *GeneratorService) processOnce(ctx context.Context) {
	fetchCtx, fetchCancel := context.WithTimeout(ctx, fetchTimeout)
	msgs, err := s.consumer.Fetch(fetchCtx, 100)
	fetchCancel()
	if err != nil {
		slog.Error("Failed to fetch messages from queue", slog.Any("error", err))
		return
	}
	if len(msgs) == 0 {
		return
	}

	forwardBatch, revertSigs, records := s.buildBatch(msgs)

	now := time.Now().UTC()
	for i := range records {
		records[i].CreatedAt = now
		records[i].UpdatedAt = now
	}

	if err := s.messageRecordRepo.BatchCreate(ctx, records); err != nil {
		slog.Error("Failed to persist message records", slog.Any("error", err))
		return
	}

	if err := s.revertRepository.BatchCreate(ctx, revertSigs); err != nil {
		slog.Error("Failed to persist revert signatures", slog.Any("error", err))
		return
	}

	// Fault injection: revert signatures are now durable but the forward batch
	// has not been published yet — exercises recovery from "we have the revert
	// path persisted but never published forward txs". NATS will redeliver the
	// un-acked messages and the BatchCreate ON CONFLICT guard is the
	// idempotency contract under test.
	if faultErr := faultinjector.Check("public_relayer.service.GeneratorService.Run.after_signature_persist"); faultErr != nil {
		slog.Error("Fault injection at after_signature_persist", slog.Any("error", faultErr))
		return
	}

	// Fault injection: simulate failure right before the forward batch is
	// published. On error the messages stay un-acked, NATS redelivers, and
	// generation must be idempotent.
	if faultErr := faultinjector.Check("public_relayer.service.GeneratorService.Run.before_executor_push"); faultErr != nil {
		slog.Error("Fault injection at before_executor_push", slog.Any("error", faultErr))
		return
	}

	if err := s.forwardBatcher.Send(ctx, forwardBatch); err != nil {
		slog.Error("Failed to publish forward batch", slog.Any("error", err))
		return // messages stay un-acked, NATS will redeliver
	}
	slog.Info("Successfully published a batch of forward messages", slog.Int("count", len(forwardBatch)))

	for _, msg := range msgs {
		// Fault injection: classic "message processed but not acked → redelivery"
		// boundary. Fired once per message in the batch.
		if faultErr := faultinjector.Check("public_relayer.service.GeneratorService.Run.before_ack"); faultErr != nil {
			slog.Error("Fault injection at before_ack", slog.Any("error", faultErr))
			continue
		}
		if err := msg.Ack(ctx); err != nil {
			slog.Error("Failed to acknowledge message", slog.Any("error", err))
		}
	}
}

// buildBatch generates forward + revert calldata for every fetched message
// and produces the three parallel slices the Run loop persists and publishes.
// Messages whose calldata generation fails are skipped with a warning so one
// bad message doesn't block the rest of the batch.
func (s *GeneratorService) buildBatch(
	msgs []msgqueue.Message[Message],
) (forward []batcher.Message, reverts []RevertSignature, records []MessageRecord) {
	forward = make([]batcher.Message, 0, len(msgs))
	reverts = make([]RevertSignature, 0, len(msgs))
	records = make([]MessageRecord, 0, len(msgs))

	for _, m := range msgs {
		forwardBytes, revertBytes, err := s.generateForwardAndRevertCalldata(m.V)
		if err != nil {
			slog.Error(
				"Failed to generate calldata for message, skipping",
				slog.String("id", m.V.ID.Hex()),
				slog.Any("error", err),
			)
			continue
		}

		id := m.V.ID.Hex()
		forward = append(forward, batcher.Message{
			ID:       id,
			Address:  s.endpointAddress,
			Calldata: forwardBytes,
		})
		reverts = append(reverts, RevertSignature{
			ID:   id,
			Data: revertBytes,
		})
		records = append(records, MessageRecord{
			ID:     id,
			Status: MessageRecordStatusNew,
		})
	}
	return forward, reverts, records
}

// HandleForwardResults is the async callback for forward tx results. Success
// results transition the record to Succeeded and are terminal. Failure
// results transition the record to Failed, then look up the pre-signed
// revert calldata and publish it via revertBatcher.
func (s *GeneratorService) HandleForwardResults(ctx context.Context, results []types.TxResult) error {
	if len(results) == 0 {
		return nil
	}

	var (
		updates    []ForwardResultUpdate
		successIDs []string
		failedIDs  []string
		failedErrs = map[string]string{}
	)
	for _, res := range results {
		switch res.Kind {
		case types.TxResultSuccess:
			updates = append(updates, ForwardResultUpdate{
				ID:     res.CorrelationID,
				Status: MessageRecordStatusSucceeded,
				Hash:   res.TxHash,
			})
			successIDs = append(successIDs, res.CorrelationID)
		case types.TxResultRevert, types.TxResultFailed:
			updates = append(updates, ForwardResultUpdate{
				ID:     res.CorrelationID,
				Status: MessageRecordStatusFailed,
				Error:  res.ErrorReason,
			})
			failedIDs = append(failedIDs, res.CorrelationID)
			failedErrs[res.CorrelationID] = res.ErrorReason
		}
	}
	slog.Info("Received results for forward messages", slog.Int("successful", len(successIDs)), slog.Int("failed", len(failedIDs)))

	if err := s.messageRecordRepo.UpdateForwardResults(ctx, updates); err != nil {
		slog.Error("Failed to update message records with forward results", slog.Any("error", err))
	}

	if len(failedIDs) == 0 {
		return nil
	}

	sigs, err := s.revertRepository.GetByIDs(ctx, failedIDs)
	if err != nil {
		slog.Error("Failed to load revert signatures for failed forwards",
			slog.Any("error", err),
			slog.Int("failed_count", len(failedIDs)),
		)
		return nil
	}

	revertBatch := make([]batcher.Message, 0, len(sigs))
	for _, sig := range sigs {
		revertBatch = append(revertBatch, batcher.Message{
			ID:       sig.ID,
			Address:  s.sourceEndpointAddress,
			Calldata: sig.Data,
		})
	}

	if err := s.revertBatcher.Send(ctx, revertBatch); err != nil {
		slog.Error("Failed to publish revert batch", slog.Any("error", err))
		return nil
	}

	slog.Info("Published revert batch for failed forwards",
		slog.Int("count", len(revertBatch)),
	)
	return nil
}

// HandleRevertResults is the observability callback for revert tx results.
// Updates the record to RevertSucceeded or RevertFailed — both terminal.
func (s *GeneratorService) HandleRevertResults(ctx context.Context, results []types.TxResult) error {
	if len(results) == 0 {
		return nil
	}

	var (
		successIDs []string
		failedIDs  []string
	)

	updates := make([]RevertResultUpdate, 0, len(results))
	for _, res := range results {
		switch res.Kind {
		case types.TxResultSuccess:
			updates = append(updates, RevertResultUpdate{
				ID:     res.CorrelationID,
				Status: MessageRecordStatusRevertSucceeded,
				Hash:   res.TxHash,
			})
			successIDs = append(successIDs, res.CorrelationID)
		case types.TxResultRevert, types.TxResultFailed:
			updates = append(updates, RevertResultUpdate{
				ID:     res.CorrelationID,
				Status: MessageRecordStatusRevertFailed,
				Error:  res.ErrorReason,
			})
			failedIDs = append(failedIDs, res.CorrelationID)
		}
	}
	slog.Info("Received results for revert messages", slog.Int("successful", len(successIDs)), slog.Int("failed", len(failedIDs)))

	if err := s.messageRecordRepo.UpdateRevertResults(ctx, updates); err != nil {
		slog.Error("Failed to update message records with revert results", slog.Any("error", err))
	}
	return nil
}

func (s *GeneratorService) generateForwardAndRevertCalldata(msg Message) ([]byte, []byte, error) {
	forwardBytes, err := s.generator.Generate(msg.FromAddress, msg.ToAddress, msg.Data, msg.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("generate forward transaction: %w", err)
	}

	revertData := RNMessageDispatcherV1.RaylsNodeMessage{
		MessageMetadata: RNMessageDispatcherV1.RaylsNodeMessageMetadata{
			Nonce:             new(big.Int),
			RevertPayloadData: make([]byte, 0),
			NewResourceMetadata: RNMessageDispatcherV1.RaylsNodeNewResourceMetadata{
				Bytecode:          make([]byte, 0),
				InitializerParams: make([]byte, 0),
			},
			TransferMetadata: RNMessageDispatcherV1.RaylsNodeBridgedTransferMetadata{
				Id:     new(big.Int),
				Amount: new(big.Int),
			},
		},
		Payload: msg.Data.MessageMetadata.RevertPayloadData,
	}

	// Same message ID for the revert tx on the source chain as the forward
	// tx on the destination chain — keeps the two halves of the lifecycle
	// correlated.
	revertBytes, err := s.generator.Generate(msg.ToAddress, msg.FromAddress, revertData, msg.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("generate revert transaction: %w", err)
	}

	return forwardBytes, revertBytes, nil
}
