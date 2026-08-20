package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/faultinjector"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"golang.org/x/sync/errgroup"
)

//go:generate moq --pkg service_test -out enygma_mock_test.go . EnygmaBatchMQ EnygmaInitiator EnygmaReverter EnygmaBlockWaiterService EnygmaFinalizationService EnygmaBatcher
type EnygmaBatchMQ interface {
	Fetch(ctx context.Context, count int) ([]msgqueue.Message[EnygmaSerializedEvent], error)
}

type EnygmaEthereumClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
}

type EnygmaBlockWaiterService interface {
	WaitForBlock(ctx context.Context, targetBlockNumber uint64) (uint64, error)
}

type EnygmaInitiator interface {
	HandleEnygmaCreation(ctx context.Context, chainEventID string, resourceId string, blockNumber uint64, initialSupply *big.Int) (uint64, error)
	HandleEnygmaSupplyUpdates(ctx context.Context, batchID string, resourceId string, blockNumber uint64, batch types.EnygmaSupplyUpdate) (uint64, error)
	HandleEnygmaCrossTransfer(ctx context.Context, batchID string, resourceID string, blockNumber uint64, batch map[string][]*types.EnygmaTransferBatchTx) (uint64, error)
	HandleEnygmaDeposit(ctx context.Context, chainEventID string, blockNumber uint64, resourceID string, referenceID [32]byte, from common.Address, amount *big.Int, txHash common.Hash) (uint64, error)
	HandleEnygmaWithdrawal(ctx context.Context, chainEventID string, blockNumber uint64, resourceID string, referenceID [32]byte, to common.Address, amount *big.Int, txHash common.Hash) (uint64, error)
}

type EnygmaReverter interface {
	RevertEnygmaSupplyUpdate(ctx context.Context, resourceId string, supplyUpdateEvents []EnygmaSupplyUpdate) error
	RevertEnygmaDeposit(
		ctx context.Context,
		resourceId string,
		referenceId [32]byte,
		from common.Address,
		amount *big.Int,
	) error
	RevertEnygmaTransfer(
		parentCtx context.Context,
		resourceId string,
		txsByChainID map[string][]*types.EnygmaTransferBatchTx,
	) error
}

type EnygmaBatcher interface {
	BatchTransfers(ctx context.Context, txsByChainID map[string][]*types.EnygmaTransferBatchTx) ([]map[string][]*types.EnygmaTransferBatchTx, error)
	GroupTransfersByChainID(events []EnygmaTransferTx) map[string][]*types.EnygmaTransferBatchTx
}

type EnygmaFinalizationService interface {
	ExecuteFinalization(ctx context.Context, id string, blockNumber uint64, resourceID string) error
}

type EnygmaOrchestrator struct {
	ticker                   *time.Ticker
	enygmaMQ                 EnygmaBatchMQ
	blockWaiter              EnygmaBlockWaiterService
	enygmaBatcher            EnygmaBatcher
	initiator                EnygmaInitiator
	reverter                 EnygmaReverter
	finalizer                EnygmaFinalizationService
	pendingFinalizations     chan pendingFinalization
	maxConcurrentResourceIDs int
}

func NewEnygmaOrchestrator(tickerPeriod time.Duration, enygmaMQ EnygmaBatchMQ, blockWaiter EnygmaBlockWaiterService, enygmaBatcher EnygmaBatcher, initiator EnygmaInitiator, reverter EnygmaReverter, finalizer EnygmaFinalizationService, maxConcurrentResourceIDs int) *EnygmaOrchestrator {
	if maxConcurrentResourceIDs < 1 {
		maxConcurrentResourceIDs = 1
	}
	return &EnygmaOrchestrator{
		ticker:                   time.NewTicker(tickerPeriod),
		enygmaMQ:                 enygmaMQ,
		blockWaiter:              blockWaiter,
		enygmaBatcher:            enygmaBatcher,
		initiator:                initiator,
		reverter:                 reverter,
		finalizer:                finalizer,
		pendingFinalizations:     make(chan pendingFinalization, 100),
		maxConcurrentResourceIDs: maxConcurrentResourceIDs,
	}
}

const (
	maxFinalizationRetries     = 5
	finalizationRetryBaseDelay = 2 * time.Second
)

func (o *EnygmaOrchestrator) groupMsgsByResourceAndEvent(
	msgs []msgqueue.Message[EnygmaSerializedEvent],
) map[string]map[EnygmaEventType][]msgqueue.Message[EnygmaSerializedEvent] {

	grouped := make(map[string]map[EnygmaEventType][]msgqueue.Message[EnygmaSerializedEvent])

	for _, msg := range msgs {
		resourceID := msg.V.ResourceID
		evType := msg.V.Type

		if grouped[resourceID] == nil {
			grouped[resourceID] = make(map[EnygmaEventType][]msgqueue.Message[EnygmaSerializedEvent])
		}
		grouped[resourceID][evType] = append(grouped[resourceID][evType], msg)
	}

	return grouped
}

func (o *EnygmaOrchestrator) Run(ctx context.Context) error {
	slog.Info("EnygmaOrchestrator started")

	go o.processFinalizationRetryQueue(ctx)

	initialRun := make(chan struct{}, 1)
	initialRun <- struct{}{}

	for {
		select {
		case <-ctx.Done():
			slog.Info("EnygmaOrchestrator shutting down")
			return nil
		case <-o.ticker.C:
		case <-initialRun:
		}

		slog.Debug("Fetching enygma messages", slog.Int("max_size", 1000))
		fetchCtx, fetchCancel := context.WithTimeout(ctx, 30*time.Second)
		msgs, err := o.enygmaMQ.Fetch(fetchCtx, 1000)
		fetchCancel()
		if err != nil {
			slog.Warn("Failed to fetch enygma messages", slog.Any("error", err))
			continue
		}
		if len(msgs) == 0 {
			slog.Debug("No enygma messages available")
			continue
		}

		slog.Info("Fetched enygma messages", slog.Int("count", len(msgs)))

		grouped := o.groupMsgsByResourceAndEvent(msgs)

		eg, egCtx := errgroup.WithContext(ctx)
		eg.SetLimit(o.maxConcurrentResourceIDs)

		for resourceID, typeMap := range grouped {
			eg.Go(func() error {
				slog.Info("Processing enygma for resource", slog.String("ID", resourceID))
				o.processResource(egCtx, resourceID, typeMap)
				return nil
			})
		}

		eg.Wait()
	}
}

func (o *EnygmaOrchestrator) processResource(
	ctx context.Context,
	resourceID string,
	typeMap map[EnygmaEventType][]msgqueue.Message[EnygmaSerializedEvent],
) {
	for evType := range EnygmaEventTypeCount {
		messages, ok := typeMap[evType]
		if !ok {
			continue
		}

		// Sample the latest block per event type rather than once for the whole resource.
		// Processing one event type can take long enough for the chain to advance, so a
		// single shared sample would hand every later event type a stale baseline.
		slog.Debug("Waiting for latest block", slog.String("resource_id", resourceID))
		latestBlock, err := o.blockWaiter.WaitForBlock(ctx, 0)
		if err != nil {
			slog.Warn("Failed to wait for block", slog.Any("error", err), slog.String("resource_id", resourceID))
			return
		}
		slog.Debug("Got latest block", slog.Uint64("block_number", latestBlock), slog.String("resource_id", resourceID))

		switch evType {
		case EnygmaCreationEvent:
			o.processEnygmaCreationMsgs(ctx, resourceID, messages, latestBlock)
		case EnygmaSupplyUpdateEvent:
			o.processEnygmaSupplyUpdateMsgs(ctx, resourceID, messages, latestBlock)
		case EnygmaTransferEvent:
			o.processEnygmaTransferMsgs(ctx, resourceID, messages, latestBlock)
		case EnygmaDepositEvent:
			o.processEnygmaDepositMsgs(ctx, resourceID, messages, latestBlock)
		case EnygmaWithdrawEvent:
			o.processEnygmaWithdrawalMsgs(ctx, resourceID, messages, latestBlock)
		}
	}
}

func (o *EnygmaOrchestrator) processEnygmaCreationMsgs(ctx context.Context, resourceID string, messages []msgqueue.Message[EnygmaSerializedEvent], latestBlock uint64) {
	if len(messages) == 0 {
		return
	}

	resourceIDLog := slog.With(slog.String("resource_id", resourceID))
	resourceIDLog.Info("Processing enygma creation messages", slog.Int("count", len(messages)))

	for _, msg := range messages {
		resourceIDLog.Debug("Deserializing creation message")
		typed, err := deserializeEnygmaMessage[EnygmaCreation](msg.V)
		if err != nil {
			slog.Error("Failed to deserialize enygma creation events", slog.Any("error", err))
			continue
		}

		lastBlock, shouldAck, shouldFinalize := o.handleEnygmaCreation(ctx, msg.V.GetID(), resourceID, typed.Event, latestBlock)
		if shouldAck {
			o.ackMessage(ctx, msg)
		}

		if shouldFinalize {
			resourceIDLog.Debug("Executing finalization for creation", slog.Uint64("block_number", lastBlock))
			if err := o.finalizer.ExecuteFinalization(ctx, msg.V.GetID(), lastBlock, resourceID); err != nil {
				resourceIDLog.Error("Failed to execute finalization, queuing for retry", slog.Any("error", err))
				o.queueFinalizationRetry(msg.V.GetID(), lastBlock, resourceID)
			}
		}
	}
}

func (o *EnygmaOrchestrator) processEnygmaSupplyUpdateMsgs(ctx context.Context, resourceID string, messages []msgqueue.Message[EnygmaSerializedEvent], latestBlock uint64) {
	if len(messages) == 0 {
		return
	}

	resourceIDLog := slog.With(slog.String("resource_id", resourceID))
	resourceIDLog.Info("Processing enygma supply update messages", slog.Int("count", len(messages)))

	var batch []EnygmaSupplyUpdate
	var msgsToAck []msgqueue.Message[EnygmaSerializedEvent]

	for _, msg := range messages {
		typed, err := deserializeEnygmaMessage[EnygmaSupplyUpdate](msg.V)
		if err != nil {
			slog.Error("Failed to deserialize supply update message", slog.Any("error", err))
			continue
		}
		batch = append(batch, typed.Event)
		msgsToAck = append(msgsToAck, msg)
	}

	resourceIDLog.Info("Processing enygma supply update batch")
	// TODO: ??? what message id can I use here ???
	lastBlock, shouldAck, shouldFinalize := o.handleEnygmaSupplyUpdateBatch(ctx, uuid.New().String(), resourceID, batch, latestBlock)
	if shouldAck {
		resourceIDLog.Info("Acknowledging messages")
		for _, msg := range msgsToAck {
			o.ackMessage(ctx, msg)
		}
	}

	if shouldFinalize {
		// TODO: ??? what message id can I use here ???
		finalizationID := uuid.New().String()
		resourceIDLog.Debug("Executing finalization for supply update", slog.Uint64("block_number", lastBlock))
		if err := o.finalizer.ExecuteFinalization(ctx, finalizationID, lastBlock, resourceID); err != nil {
			resourceIDLog.Error("Failed to execute finalization, queuing for retry", slog.Any("error", err))
			o.queueFinalizationRetry(finalizationID, lastBlock, resourceID)
		}
	}
}

func (o *EnygmaOrchestrator) processEnygmaTransferMsgs(ctx context.Context, resourceID string, messages []msgqueue.Message[EnygmaSerializedEvent], latestBlock uint64) {
	if len(messages) == 0 {
		return
	}

	resourceIDLog := slog.With(slog.String("resource_id", resourceID))
	resourceIDLog.Info("Processing enygma cross-transfer messages", slog.Int("count", len(messages)))

	events := make([]EnygmaTransferTx, 0)
	msgIdToMsg := make(map[string]msgqueue.Message[EnygmaSerializedEvent])

	for _, msg := range messages {
		typed, err := deserializeEnygmaMessage[EnygmaTransferTx](msg.V)
		if err != nil {
			slog.Error("Failed to deserialize enygma transfer message", slog.Any("error", err))
			continue
		}
		events = append(events, typed.Event)
		msgIdToMsg[typed.Event.MessageId] = msg
	}

	txsByChainID := o.enygmaBatcher.GroupTransfersByChainID(events)
	batches, err := o.enygmaBatcher.BatchTransfers(ctx, txsByChainID)
	if err != nil {
		slog.Error("Failed to batch transfers", slog.Any("error", err))
		return
	}

	var needFinalization bool

	for i, batch := range batches {
		batchLog := resourceIDLog.With(slog.String("batch_index", formatMsgIndex(i, batches)))
		batchLog.Info("Processing enygma cross-transfer batch")

		// TODO: ??? what message id can I use here ???
		lastTransferBlock, shouldAck, shouldFinalize := o.handleEnygmaTransferBatch(ctx, uuid.New().String(), resourceID, batch, latestBlock)
		if shouldAck {
			batchLog.Info("Acknowledging batch messages")
			for _, txs := range batch {
				for _, tx := range txs {
					if msg, ok := msgIdToMsg[tx.MessageId]; ok {
						o.ackMessage(ctx, msg)
					} else {
						slog.Error("No message found for tx MessageId", slog.String("msg_id", tx.MessageId))
					}
				}
			}
		}
		// We need to execute finalization tx for at least 1 successful transfer
		if shouldFinalize {
			needFinalization = true
		}
		latestBlock = lastTransferBlock

		batchLog.Info("Enygma cross-transfer batch processed")
	}

	if needFinalization {
		// TODO: ??? what message id can I use here ???
		finalizationID := uuid.New().String()
		resourceIDLog.Debug("Executing finalization for cross-transfers", slog.Uint64("block_number", latestBlock))
		if err := o.finalizer.ExecuteFinalization(ctx, finalizationID, latestBlock, resourceID); err != nil {
			resourceIDLog.Error("Failed to execute finalization, queuing for retry", slog.Any("error", err))
			o.queueFinalizationRetry(finalizationID, latestBlock, resourceID)
		}
	}
}

func (o *EnygmaOrchestrator) processEnygmaDepositMsgs(ctx context.Context, resourceID string, messages []msgqueue.Message[EnygmaSerializedEvent], latestBlock uint64) {
	if len(messages) == 0 {
		return
	}

	resourceIDLog := slog.With(slog.String("resource_id", resourceID))
	resourceIDLog.Info("Processing enygma deposit messages", slog.Int("count", len(messages)))

	var needFinalization bool

	for i, msg := range messages {
		msgLog := resourceIDLog.With(slog.String("msg_index", formatMsgIndex(i, messages)))
		msgLog.Debug("Deserializing deposit message")
		typed, err := deserializeEnygmaMessage[EnygmaDepositToDvp](msg.V)
		if err != nil {
			slog.Error("Failed to deserialize enygma deposit event", slog.Any("error", err))
			continue
		}

		msgLog.Info("Processing deposit message")
		lastDepositBlock, shouldAck, shouldFinalize := o.handleEnygmaDeposit(ctx, msg.V.GetID(), resourceID, typed.Event, latestBlock)
		if shouldAck {
			o.ackMessage(ctx, msg)
		}
		// We need to execute finalization tx for at least 1 successful deposit
		if shouldFinalize {
			needFinalization = true
		}
		latestBlock = lastDepositBlock

		msgLog.Info("Deposit message processed")
	}

	if needFinalization {
		// TODO: ??? what message id can I use here ???
		finalizationID := uuid.New().String()
		resourceIDLog.Debug("Executing finalization for deposits", slog.Uint64("block_number", latestBlock))
		if err := o.finalizer.ExecuteFinalization(ctx, finalizationID, latestBlock, resourceID); err != nil {
			resourceIDLog.Error("Failed to execute finalization, queuing for retry", slog.Any("error", err))
			o.queueFinalizationRetry(finalizationID, latestBlock, resourceID)
		}
	}
}

func (o *EnygmaOrchestrator) processEnygmaWithdrawalMsgs(ctx context.Context, resourceID string, messages []msgqueue.Message[EnygmaSerializedEvent], latestBlock uint64) {
	if len(messages) == 0 {
		return
	}

	resourceIDLog := slog.With(slog.String("resource_id", resourceID))
	resourceIDLog.Info("Processing enygma withdrawal messages", slog.Int("count", len(messages)))

	var needFinalization bool

	for i, msg := range messages {
		msgLog := resourceIDLog.With(slog.String("msg_index", formatMsgIndex(i, messages)))
		msgLog.Debug("Deserializing withdrawal message")
		typed, err := deserializeEnygmaMessage[EnygmaWithdrawFromDvp](msg.V)
		if err != nil {
			slog.Error("Failed to deserialize enygma withdraw event", slog.Any("error", err))
			continue
		}

		msgLog.Info("Processing withdrawal message")
		lastWithdrawalBlock, shouldAck, shouldFinalize := o.handleEnygmaWithdrawal(ctx, msg.V.GetID(), resourceID, typed.Event, latestBlock)
		if shouldAck {
			o.ackMessage(ctx, msg)
		}
		// We need to execute finalization tx for at least 1 successful withdrawal
		if shouldFinalize {
			needFinalization = true
		}
		latestBlock = lastWithdrawalBlock

		msgLog.Info("Withdrawal message processed")
	}

	if needFinalization {
		// TODO: ??? what message id can I use here ???
		finalizationID := uuid.New().String()
		resourceIDLog.Debug("Executing finalization for withdrawals", slog.Uint64("block_number", latestBlock))
		if err := o.finalizer.ExecuteFinalization(ctx, finalizationID, latestBlock, resourceID); err != nil {
			resourceIDLog.Error("Failed to execute finalization, queuing for retry", slog.Any("error", err))
			o.queueFinalizationRetry(finalizationID, latestBlock, resourceID)
		}
	}
}

func (o *EnygmaOrchestrator) handleEnygmaCreation(ctx context.Context, chainEventID string, resourceID string, event EnygmaCreation, latestBlock uint64) (uint64, bool, bool) {
	resourceIDLog := slog.With(slog.String("resource_id", resourceID))
	resourceIDLog.Info("Handling enygma creation", slog.String("initial_supply", event.InitialSupply.String()))

	creationBlockNumber, err := o.initiator.HandleEnygmaCreation(ctx, chainEventID, resourceID, latestBlock, event.InitialSupply)
	if err != nil {
		resourceIDLog.Error("Failed to handle enygma creation", slog.Any("error", err))
		return creationBlockNumber, false, false
	}
	resourceIDLog.Info("Successfully handled enygma creation", slog.String("resource_id", resourceID))

	return creationBlockNumber, true, true
}

func (o *EnygmaOrchestrator) handleEnygmaSupplyUpdateBatch(ctx context.Context, batchID string, resourceID string, events []EnygmaSupplyUpdate, latestBlock uint64) (uint64, bool, bool) {
	resourceIDLog := slog.With(slog.String("resource_id", resourceID))
	resourceIDLog.Info("Handling enygma supply update event batch")

	var batch types.EnygmaSupplyUpdate
	batch.Amount = big.NewInt(0)

	for _, event := range events {
		batch.Amount = batch.Amount.Add(batch.Amount, event.Amount)
	}
	if batch.Amount.Cmp(big.NewInt(0)) >= 0 {
		batch.Type = types.EnygmaMint
	} else {
		batch.Type = types.EnygmaBurn
		batch.Amount = batch.Amount.Abs(batch.Amount)
	}

	resourceIDLog.Debug("Handling enygma supply update", slog.Int("type", int(batch.Type)), slog.String("amount", batch.Amount.String()))
	supplyBlockNumber, err := o.initiator.HandleEnygmaSupplyUpdates(ctx, batchID, resourceID, latestBlock, batch)
	if err != nil {
		resourceIDLog.Error("Failed to handle enygma supply updates", slog.Any("error", err))
		resourceIDLog.Info("Reverting enygma supply update")
		err = o.reverter.RevertEnygmaSupplyUpdate(ctx, resourceID, events)
		if err != nil {
			resourceIDLog.Error("Failed to revert enygma supply update", slog.Any("error", err))
			return supplyBlockNumber, false, false
		}
		resourceIDLog.Info("Successfully reverted enygma supply update")
		return supplyBlockNumber, true, false
	}

	resourceIDLog.Info("Successfully handled enygma supply update")

	return supplyBlockNumber, true, true
}

func (o *EnygmaOrchestrator) handleEnygmaTransferBatch(ctx context.Context, batchID string, resourceID string, batch map[string][]*types.EnygmaTransferBatchTx, latestBlock uint64) (uint64, bool, bool) {
	resourceIDLog := slog.With(slog.String("resource_id", resourceID))
	resourceIDLog.Info("Handling enygma cross-transfer event batch")

	latestBlock, err := o.blockWaiter.WaitForBlock(ctx, latestBlock)
	if err != nil {
		resourceIDLog.Error("Failed to wait for block", slog.Any("error", err))
		return latestBlock, false, false
	}

	transferBlockNumber, err := o.initiator.HandleEnygmaCrossTransfer(ctx, batchID, resourceID, latestBlock, batch)
	if err != nil {
		// On-chain reverts are decoded into a concise `revert reason` at the executor level
		// (ExecuteEnygmaCrossTransfer), where the revert originates — so the contractclient/txsim
		// revert types stay out of this top-level orchestration logic.
		resourceIDLog.Error("Failed to handle enygma cross-transfer batch", slog.Any("error", err))
		resourceIDLog.Info("Reverting enygma cross-transfer batch")
		err = o.reverter.RevertEnygmaTransfer(ctx, resourceID, batch)
		if err != nil {
			resourceIDLog.Error("Failed to revert enygma cross-transfer batch", slog.Any("error", err))
			return transferBlockNumber, false, false
		}
		resourceIDLog.Info("Successfully reverted enygma cross-transfer batch")
		if faultErr := faultinjector.Check("private_relayer.source.service.EnygmaOrchestrator.handleEnygmaTransfer.after_revert"); faultErr != nil {
			resourceIDLog.Error("Fault injection triggered after revert", slog.Any("error", faultErr))
		}
		return transferBlockNumber, true, false
	}

	resourceIDLog.Info("Successfully handled enygma cross-transfer batch")
	if faultErr := faultinjector.Check("private_relayer.source.service.EnygmaOrchestrator.handleEnygmaTransfer.after_cross_transfer"); faultErr != nil {
		resourceIDLog.Error("Fault injection triggered after cross transfer", slog.Any("error", faultErr))
	}

	return transferBlockNumber, true, true
}

func (o *EnygmaOrchestrator) handleEnygmaDeposit(ctx context.Context, chainEventID string, resourceID string, event EnygmaDepositToDvp, latestBlock uint64) (uint64, bool, bool) {
	resourceIDLog := slog.With(slog.String("resource_id", resourceID))
	resourceIDLog.Info("Handling enygma deposit event")

	latestBlock, err := o.blockWaiter.WaitForBlock(ctx, latestBlock)
	if err != nil {
		slog.Error("Failed to wait for block after deposit", slog.Any("error", err))
		return latestBlock, false, false
	}

	depositLog := resourceIDLog.With(slog.String("from", event.From.Hex()), slog.String("amount", event.Amount.String()))
	depositLog.Debug("Handling enygma deposit event")

	depositBlockNumber, err := o.initiator.HandleEnygmaDeposit(ctx, chainEventID, latestBlock, resourceID, event.ReferenceId, event.From, event.Amount, event.TxHash)
	if err != nil {
		resourceIDLog.Error("Failed to handle enygma deposit", slog.Any("error", err))

		depositLog.Info("Reverting enygma deposit")
		err = o.reverter.RevertEnygmaDeposit(ctx, resourceID, event.ReferenceId, event.From, event.Amount)
		if err != nil {
			depositLog.Error("Failed to revert enygma deposit", slog.Any("error", err))
			return depositBlockNumber, false, false
		} else {
			depositLog.Info("Successfully reverted enygma deposit")
			return depositBlockNumber, true, false
		}
	}

	resourceIDLog.Info("Successfully handled enygma deposit")
	depositLog.Debug("Successfully handled enygma deposit")

	return depositBlockNumber, true, true
}

func (o *EnygmaOrchestrator) handleEnygmaWithdrawal(ctx context.Context, chainEventID string, resourceID string, event EnygmaWithdrawFromDvp, latestBlock uint64) (uint64, bool, bool) {
	resourceIDLog := slog.With(slog.String("resource_id", resourceID))
	resourceIDLog.Info("Handling enygma withdrawal event")

	latestBlock, err := o.blockWaiter.WaitForBlock(ctx, latestBlock)
	if err != nil {
		resourceIDLog.Error("Failed to wait for block after withdrawal", slog.Any("error", err))
		return latestBlock, false, false
	}

	withdrawalLog := resourceIDLog.With(slog.String("to", event.To.Hex()), slog.String("amount", event.Amount.String()))
	withdrawalLog.Debug("Handling enygma withdrawal event")

	withdrawalBlockNumber, err := o.initiator.HandleEnygmaWithdrawal(ctx, chainEventID, latestBlock, resourceID, event.ReferenceId, event.To, event.Amount, event.TxHash)
	if err != nil {
		withdrawalLog.Error("Failed to handle enygma withdrawal", slog.Any("error", err))
		return withdrawalBlockNumber, false, false
	}

	withdrawalLog.Debug("Successfully handled enygma withdrawal")
	resourceIDLog.Info("Successfully handled enygma withdrawal")

	return withdrawalBlockNumber, true, true
}

func (o *EnygmaOrchestrator) ackMessage(ctx context.Context, msg msgqueue.Message[EnygmaSerializedEvent]) {
	ackLog := slog.With(slog.String("resource_id", msg.V.ResourceID), slog.String("msg_id", msg.V.Id))
	ackLog.Debug("Acknowledging message")
	if err := msg.Ack(ctx); err != nil {
		ackLog.Error("Failed to acknowledge message", slog.Any("error", err))
	} else {
		ackLog.Debug("Message acknowledged successfully")
	}
}

// processFinalizationRetryQueue drains the pendingFinalizations channel, waits for the next block,
// and retries finalization. Items exceeding maxFinalizationRetries are dropped with an error log.
func (o *EnygmaOrchestrator) processFinalizationRetryQueue(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case pf := <-o.pendingFinalizations:
			retryLog := slog.With(
				slog.String("resource_id", pf.ResourceID),
				slog.Uint64("block_number", pf.BlockNumber),
				slog.Int("retry_count", pf.RetryCount),
			)

			retryLog.Info("Retrying finalization")

			// Backoff delay to avoid rapid-fire retries when the target block has
			// already passed and WaitForBlock returns immediately.
			delay := finalizationRetryBaseDelay * time.Duration(pf.RetryCount+1)
			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
			}

			_, err := o.blockWaiter.WaitForBlock(ctx, pf.BlockNumber)
			if err != nil {
				retryLog.Warn("Failed to wait for block during finalization retry", slog.Any("error", err))
				o.requeueFinalization(pf, retryLog)
				continue
			}

			if err := o.finalizer.ExecuteFinalization(ctx, pf.ID, pf.BlockNumber, pf.ResourceID); err != nil {
				retryLog.Error("Finalization retry failed", slog.Any("error", err))
				o.requeueFinalization(pf, retryLog)
			} else {
				retryLog.Info("Finalization retry succeeded")
			}
		}
	}
}

// requeueFinalization re-enqueues a failed finalization with an incremented retry count,
// or drops it if max retries are exceeded.
func (o *EnygmaOrchestrator) requeueFinalization(pf pendingFinalization, log *slog.Logger) {
	pf.RetryCount++
	if pf.RetryCount > maxFinalizationRetries {
		log.Error("Finalization exceeded max retries, dropping",
			slog.Int("max_retries", maxFinalizationRetries),
		)
		return
	}

	select {
	case o.pendingFinalizations <- pf:
	default:
		log.Error("Finalization retry queue is full, dropping finalization")
	}
}

// queueFinalizationRetry enqueues a finalization for retry. Used when initial finalization fails.
func (o *EnygmaOrchestrator) queueFinalizationRetry(id string, blockNumber uint64, resourceID string) {
	pf := pendingFinalization{
		ID:          id,
		BlockNumber: blockNumber,
		ResourceID:  resourceID,
		RetryCount:  0,
	}

	select {
	case o.pendingFinalizations <- pf:
		slog.Info("Queued finalization for retry",
			slog.String("resource_id", resourceID),
			slog.Uint64("block_number", blockNumber),
		)
	default:
		slog.Error("Finalization retry queue is full, dropping finalization",
			slog.String("resource_id", resourceID),
			slog.Uint64("block_number", blockNumber),
		)
	}
}

func deserializeEnygmaMessage[T any](serialized EnygmaSerializedEvent) (EnygmaTypedEvent[T], error) {
	var zero EnygmaTypedEvent[T]

	typed := EnygmaTypedEvent[T]{
		LogIndex:    serialized.LogIndex,
		TxHash:      serialized.TxHash,
		BlockNumber: serialized.BlockNumber,
		ResourceID:  serialized.ResourceID,
		Type:        serialized.Type,
	}
	err := json.Unmarshal(serialized.SerializedEvent, &typed.Event)
	if err != nil {
		return zero, fmt.Errorf("failed to unmarshal events")
	}
	return typed, nil
}

func formatMsgIndex[T any](i int, arr []T) string {
	return fmt.Sprintf("%d/%d", i+1, len(arr))
}
