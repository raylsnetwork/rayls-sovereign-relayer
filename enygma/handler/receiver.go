package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contractclient"
	"github.com/raylsnetwork/rayls-sovereign-relayer/conv"
	"github.com/raylsnetwork/rayls-sovereign-relayer/faultinjector"
	telemetry "github.com/raylsnetwork/rayls-sovereign-relayer/otel"
	repository "github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/txsim"
	"github.com/raylsnetwork/rayls-sovereign-relayer/txutil"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/maps"
)

const millisPerSecond = txutil.MillisPerSecond

type receiverEndpointClient interface {
	GetResourceAddress(ctx context.Context, resourceId string) (common.Address, error)
}

type receiverDeployer interface {
	Deploy(ctx context.Context, resourceId [32]byte, initiatorChainId *big.Int) (common.Address, error)
}

type receiverEnygmaHistoryRepository interface {
	InsertEnygmaHistory(ctx context.Context, history types.EnygmaHistory) error
	GetEnygmaHistoryByUniqueKey(
		ctx context.Context,
		resourceId string,
		blockNumberPrivateHub *big.Int,
		fromChainId *big.Int,
		eventType types.EnygmaEventType,
	) (*types.EnygmaHistory, error)
}

type receiverTxRecoveryRepository interface {
	Insert(ctx context.Context, data types.TxRecoveryData) error
	GetByPrivateHubTxHash(ctx context.Context, privateHubTxHash string) (*types.TxRecoveryData, error)
	MarkConfirmed(ctx context.Context, privateHubTxHash string) error
}

type receiverTransactionSimulator interface {
	GetRevertReason(context.Context, common.Hash) (txsim.ContractError, error)
	DecodeRevertBytes(data []byte) (txsim.ContractError, error)
}

type receiverTeleportClient interface {
	SendTransferCompleted(ctx context.Context, messages []types.EnygmaTransferCompleted) error
}

type receiverTracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

type receiverEnygmaCreationService interface {
	CreateEnygma(ctx context.Context, resourceId string, fromChainId *big.Int, blockNumberPrivateHub *big.Int) error
}
type enygmaHandlerClient interface {
	ReceiveDestTransferBatch(ctx context.Context, transfers []*types.EnygmaCrossTransferData) (map[string]contractclient.BatchResult, error)
	RevertDestTransferBatch(ctx context.Context, tokenAddress common.Address, reverts []*types.EnygmaCrossTransferData) (map[string]contractclient.BatchResult, error)
}

type Receiver struct {
	plChainId               *big.Int
	plEndpointClient        receiverEndpointClient
	plSimulator             receiverTransactionSimulator
	enygmaHistoryRepository receiverEnygmaHistoryRepository
	txRecoveryRepository    receiverTxRecoveryRepository
	destDeployer            receiverDeployer
	teleportClient          receiverTeleportClient
	enygmaHandlerClient     enygmaHandlerClient
	tracer                  receiverTracer
	creationService         receiverEnygmaCreationService
	mintReceiptTimeout      time.Duration
}

func NewReceiver(
	plChainId *big.Int,
	plEndpointClient receiverEndpointClient,
	plSimulator receiverTransactionSimulator,
	destDeployer receiverDeployer,
	enygmaHistoryRepository receiverEnygmaHistoryRepository,
	txRecoveryRepository receiverTxRecoveryRepository,
	teleportClient receiverTeleportClient,
	enygmaHandlerClient enygmaHandlerClient,
	tracer receiverTracer,
	creationService receiverEnygmaCreationService,
	mintReceiptTimeout time.Duration,
) *Receiver {
	return &Receiver{
		plChainId:               plChainId,
		plEndpointClient:        plEndpointClient,
		plSimulator:             plSimulator,
		destDeployer:            destDeployer,
		enygmaHistoryRepository: enygmaHistoryRepository,
		txRecoveryRepository:    txRecoveryRepository,
		teleportClient:          teleportClient,
		enygmaHandlerClient:     enygmaHandlerClient,
		tracer:                  tracer,
		creationService:         creationService,
		mintReceiptTimeout:      mintReceiptTimeout,
	}
}

// HandleEnygmaCrossTransfer finishes the cross-chain Enygma token transfers on the destination side.
// It deploys tokens if needed, executes mint batches, reverts failed transactions, and notifies governance.
//
//nolint:gocognit,gocyclo // complex cross-chain transfer with proof verification
func (r *Receiver) HandleEnygmaCrossTransfer(
	ctx context.Context,
	batch *types.EnygmaTransferBatch,
) error { //nolint:gocognit,gocyclo // multi-step cross-chain transfer with error handling at each step
	ctx, span := r.tracer.Start(ctx, telemetry.SPAN_FINISH_ENYGMA_CROSS_TRANSFER)
	defer span.End()

	// Validate batch fields before accessing them
	if batch == nil {
		slog.Error("HandleEnygmaCrossTransfer: received nil batch")
		return fmt.Errorf("received nil batch")
	}
	if batch.FromChainID == nil {
		slog.Error("HandleEnygmaCrossTransfer: batch.FromChainID is nil", slog.String("resourceId", batch.ResourceId))
		return fmt.Errorf("batch.FromChainID is nil for resource %s", batch.ResourceId)
	}

	blockNumberStr := "nil"
	if batch.BlockNumberPrivateHub != nil {
		blockNumberStr = batch.BlockNumberPrivateHub.String()
	}

	span.SetAttributes(
		attribute.String(telemetry.ATTR_RESOURCE_ID, batch.ResourceId),
		attribute.String(telemetry.ATTR_BLOCK_NUMBER_PNH, blockNumberStr),
		attribute.String(telemetry.ATTR_SENDER_CHAIN_ID, batch.FromChainID.String()),
	)

	// We are the sender. We do not have to do anything.
	if batch.FromChainID.Cmp(r.plChainId) == 0 {
		slog.Info("Enygma tx sent by us. Skipping...")
		return nil
	}

	slog.Info(
		"Executing cross transfer batch!",
		slog.Any("Total Transactions", len(batch.Transactions)),
		slog.Any("ResourceId", batch.ResourceId),
	)

	resourceId, err := conv.StringToBytes32(batch.ResourceId)
	if err != nil {
		return fmt.Errorf("convert resource ID to bytes32: %w", err)
	}

	enygmaPLAddress, err := r.plEndpointClient.GetResourceAddress(ctx, batch.ResourceId)
	if err != nil {
		return fmt.Errorf("error getting contract address from resource ID %s: %w", batch.ResourceId, err)
	}

	if batch.BlockNumberPrivateHub == nil {
		slog.Error("HandleEnygmaCrossTransfer: batch.BlockNumberPrivateHub is nil", slog.String("resourceId", batch.ResourceId))
		return fmt.Errorf("batch.BlockNumberPrivateHub is nil for resource %s", batch.ResourceId)
	}

	// Enygma token does not exist on the destination PL.
	// It means that we receive this Enygma token for the very first time, so we must deploy it on our PL and initialize it in the DB.
	if enygmaPLAddress == common.HexToAddress("0x") {
		slog.Debug(
			"Received enygma transfer for the very first time. Deploying token to the PL",
			slog.Any("ResourceId", batch.ResourceId),
		)

		enygmaPLAddress, err = r.destDeployer.Deploy(ctx, resourceId, batch.FromChainID)
		if err != nil {
			return fmt.Errorf("deploy enygma token: %w", err)
		}

		err = r.creationService.CreateEnygma(ctx, batch.ResourceId, r.plChainId, batch.BlockNumberPrivateHub)
		if err != nil {
			return fmt.Errorf("create enygma after deploy: %w", err)
		}

		slog.Debug("Successfully deployed enygma token to the PL", slog.Any("ResourceId", batch.ResourceId))
	}

	// Crash-recovery dedup: the recovery record is the source of truth for
	// whether a cross-transfer has been fully processed on this side. The
	// previous version used `enygma_history` as the dedup key, but that
	// row is inserted BEFORE the mint dispatch — a crash between the
	// insert and the mint would leave a history row with no mint, and a
	// NATS-redelivered retry would skip the batch entirely (token loss).
	//
	// Recovery flow:
	//   - Confirmed → fully processed; the mint already landed on chain.
	//     Skip the mint but re-insert the history row idempotently — the
	//     Enygma_resync flow deletes history rows after rolling
	//     enygma_checkpoints back, then expects the receiver to recreate
	//     them by re-processing past events. Recovery rows are not
	//     touched by that flow, so without this arm a confirmed recovery
	//     would prevent resync.
	//   - Pending → first attempt crashed before MarkConfirmed; retry the
	//     mint. crossMint() is idempotent on referenceId (issue #75) so
	//     re-broadcasting a successful mint is a silent on-chain no-op.
	//   - Not present → fresh batch; insert Pending recovery + history,
	//     then proceed to mint.
	existingRecovery, err := r.txRecoveryRepository.GetByPrivateHubTxHash(ctx, batch.PrivateHubTxHash)
	if err != nil {
		return fmt.Errorf("failed to check existing recovery for resource %s: %w", batch.ResourceId, err)
	}

	totalReceivedBalance := big.NewInt(0)
	for _, tx := range batch.Transactions {
		totalReceivedBalance.Add(totalReceivedBalance, tx.ToAmount)
	}

	if existingRecovery != nil && existingRecovery.Status == types.HistoryStatusConfirmed {
		slog.Info("Recovery already confirmed, ensuring history is present",
			slog.Any("Total Transactions", len(batch.Transactions)),
			slog.Any("ResourceId", batch.ResourceId),
		)
		err = r.enygmaHistoryRepository.InsertEnygmaHistory(ctx, types.EnygmaHistory{
			ResourceId:            batch.ResourceId,
			FromChainId:           batch.FromChainID,
			RFactor:               batch.ToRValueToAdd,
			BlockNumberPrivateHub: batch.BlockNumberPrivateHub,
			BalanceChange:         totalReceivedBalance,
			EventType:             types.EnygmaTransfer,
			PrivateHubTxHash:      batch.PrivateHubTxHash,
		})
		if err != nil && !errors.Is(err, repository.ErrAlreadyProcessed) {
			slog.Error("Failed to re-insert history during resync", slog.String("resourceId", batch.ResourceId), slog.String("err", err.Error()))
			return err
		}
		return nil
	}

	if existingRecovery == nil {
		// Fresh batch: persist Pending recovery + history. The recovery
		// goes in first so a crash between recovery insert and history
		// insert still routes the next attempt through the Pending arm.
		err = r.txRecoveryRepository.Insert(ctx, types.TxRecoveryData{
			PrivateHubTxHash:      batch.PrivateHubTxHash,
			ResourceID:            batch.ResourceId,
			PrivateHubBlockNumber: batch.BlockNumberPrivateHub.Uint64(),
			FromChainID:           batch.FromChainID.String(),
			TxBytes:               nil, // CTS owns signed bytes; no relayer-side replay
			EventType:             types.EnygmaTransfer,
			TxNature:              types.TxNatureEnygma,
			Status:                types.HistoryStatusPending,
		})
		if err != nil && !errors.Is(err, repository.ErrRecoveryAlreadyExists) {
			slog.Error("Failed to insert pending recovery", slog.String("resourceId", batch.ResourceId), slog.String("err", err.Error()))
			return err
		}

		err = r.enygmaHistoryRepository.InsertEnygmaHistory(ctx, types.EnygmaHistory{
			ResourceId:            batch.ResourceId,
			FromChainId:           batch.FromChainID,
			RFactor:               batch.ToRValueToAdd,
			BlockNumberPrivateHub: batch.BlockNumberPrivateHub,
			BalanceChange:         totalReceivedBalance,
			EventType:             types.EnygmaTransfer,
			PrivateHubTxHash:      batch.PrivateHubTxHash,
		})
		if err != nil && !errors.Is(err, repository.ErrAlreadyProcessed) {
			slog.Error("Failed to update database", slog.String("resourceId", batch.ResourceId), slog.String("err", err.Error()))
			return err
		}
	} else {
		slog.Info("Resuming pending cross transfer after crash",
			slog.String("resourceId", batch.ResourceId),
			slog.String("privateHubTxHash", batch.PrivateHubTxHash),
		)
	}

	if len(batch.Transactions) == 0 {
		slog.Info("No transactions to process", slog.String("resourceId", batch.ResourceId))
		return r.txRecoveryRepository.MarkConfirmed(ctx, batch.PrivateHubTxHash)
	}

	if err := faultinjector.Check("enygma.handler.Receiver.HandleEnygmaCrossTransfer.after_insert_history"); err != nil {
		return fmt.Errorf("fault injection at after_insert_history: %w", err)
	}

	msgIdToTxMap := make(map[string]*types.EnygmaCrossTransferData)
	failedMsgIDs := make([]string, 0)
	executedMsgIDs := make([]string, 0)

	for _, tx := range batch.Transactions {
		msgIdToTxMap[tx.MessageId] = &types.EnygmaCrossTransferData{
			EnygmaTransferBatchTx: tx,
			EnygmaAddress:         enygmaPLAddress,
			FromChainID:           batch.FromChainID,
		}
	}

	_, mintSpan := r.tracer.Start(ctx, telemetry.SPAN_SEND_CROSS_MINT)
	mintResults, err := r.enygmaHandlerClient.ReceiveDestTransferBatch(ctx, maps.Values(msgIdToTxMap))
	if err != nil {
		mintSpan.RecordError(err)
		mintSpan.SetStatus(codes.Error, telemetry.STATUS_ERROR_FAILED_TO_SEND_MINT)
		return fmt.Errorf("mint transactions: %w", err)
	}
	mintSpan.SetStatus(codes.Ok, telemetry.STATUS_SUCCESS_MINT_SENT)
	mintSpan.End()

	if err := faultinjector.Check("enygma.handler.Receiver.HandleEnygmaCrossTransfer.after_mint_batch"); err != nil {
		return fmt.Errorf("fault injection at after_mint_batch: %w", err)
	}

	for msgID, mintRes := range mintResults {
		// Fault injection (test-only; no-op unless armed): force this message onto the
		// destination-mint-failure path so it routes to RevertDestTransferBatch below
		// (source re-credited, dest mint reverted -> net-zero, governance stays PENDING).
		// Reproduces the dest-failure outcome that issuer revert-traps (e.g. an
		// addressToFail override) provided before canonical factory bytecode began
		// shipping to Enygma destinations (commit 88c798ce).
		if faultErr := faultinjector.Check("enygma.handler.Receiver.HandleEnygmaCrossTransfer.fail_mint"); faultErr != nil {
			failedMsgIDs = append(failedMsgIDs, msgID)
			slog.Warn("Fault injection: forcing enygma destination mint onto the revert path", slog.String("msg id", msgID), slog.String("fault", faultErr.Error()))
			continue
		}
		if mintRes.Err != nil {
			failedMsgIDs = append(failedMsgIDs, msgID)
			logAttrs := []any{slog.String("msg id", msgID), slog.Any("error", mintRes.Err)}
			// Decode the on-chain revert data into a concise reason instead of leaving only the raw
			// revert hex. ContractError.Reason() returns the message string for a Solidity
			// require/Error(string) revert, and falls back to the error NAME (e.g.
			// "ProgramData__Unapproved", argument list stripped) for custom errors.
			var revertErr *contractclient.ErrorWithRevertData
			if errors.As(mintRes.Err, &revertErr) {
				if decoded, decErr := r.plSimulator.DecodeRevertBytes(revertErr.GetRevertData()); decErr == nil && !decoded.IsEmpty() {
					logAttrs = append(logAttrs, slog.String("revert reason", decoded.Reason()))
				}
			}
			slog.Error("Failed to execute enygma payment programmability transaction on destination", logAttrs...)
			continue
		}
		if mintRes.Receipt.Status == 1 {
			executedMsgIDs = append(executedMsgIDs, msgID)
			slog.Debug("Enygma mint succeeded on destination", slog.String("msg id", msgID))
		} else {
			failedMsgIDs = append(failedMsgIDs, msgID)
			slog.Error("Enygma mint transaction failed on destination", slog.String("msg id", msgID))
		}
	}

	// Notify governance about all the successfully executed transactions in the batch
	completedBatchTxs := make([]types.EnygmaTransferCompleted, 0, len(executedMsgIDs))
	for _, msgID := range executedMsgIDs {
		result := mintResults[msgID]
		completedBatchTxs = append(completedBatchTxs, types.EnygmaTransferCompleted{
			MessageId:       msgID,
			TransactionHash: result.Receipt.TxHash.Hex(),
			ChainId:         r.plChainId,
		})
	}

	if err := faultinjector.Check("enygma.handler.Receiver.HandleEnygmaCrossTransfer.before_transfer_completed"); err != nil {
		return fmt.Errorf("fault injection at before_transfer_completed: %w", err)
	}

	if len(completedBatchTxs) > 0 {
		_, transferCompletedSpan := r.tracer.Start(ctx, telemetry.SPAN_TRANSFER_COMPLETED)
		defer transferCompletedSpan.End()
		slog.Debug(
			"Sending Enygma transfer batch completed message to Private Hub",
			slog.Int("Total Messages", len(completedBatchTxs)),
		)
		err = r.teleportClient.SendTransferCompleted(ctx, completedBatchTxs)
		if err != nil {
			transferCompletedSpan.RecordError(err)
			transferCompletedSpan.SetStatus(codes.Error, telemetry.STATUS_ERROR_FAILED_TO_SEND_TRANSFER_COMPLETED)
			return fmt.Errorf("failed to send transfer completed: %w", err)
		}
		transferCompletedSpan.SetStatus(codes.Ok, telemetry.STATUS_SUCCESS_TRANSFER_COMPLETED_SENT)
	}

	/*
	   Revert the transactions that failed
	*/

	slog.Info(
		"Execution of cross transfer batch completed!",
		slog.Int("Total Transactions", len(batch.Transactions)),
		slog.Int("Failed Transactions", len(failedMsgIDs)),
	)

	if len(failedMsgIDs) == 0 {
		return r.txRecoveryRepository.MarkConfirmed(ctx, batch.PrivateHubTxHash)
	}

	slog.InfoContext(
		ctx,
		"Submitting failed transactions to be reverted!",
		slog.Int("Total Transactions", len(failedMsgIDs)),
	)

	revertDataList := make([]*types.EnygmaCrossTransferData, 0, len(failedMsgIDs))
	for _, msgId := range failedMsgIDs {
		tx := msgIdToTxMap[msgId]
		revertDataList = append(revertDataList, &types.EnygmaCrossTransferData{
			EnygmaTransferBatchTx: &types.EnygmaTransferBatchTx{
				MessageId:   tx.MessageId,
				ReferenceId: tx.ReferenceId,
				// Switching the addresses because we are reverting the transaction
				FromAddress: tx.ToAddress,
				ToAddress:   tx.FromAddress,
				ToAmount:    tx.ToAmount,
			},
			EnygmaAddress: enygmaPLAddress,
			FromChainID:   batch.FromChainID,
		})
	}
	_, revertSpan := r.tracer.Start(ctx, telemetry.SPAN_REVERT_ENYGMA_TRANSFER)
	defer revertSpan.End()
	revertResults, err := r.enygmaHandlerClient.RevertDestTransferBatch(ctx, enygmaPLAddress, revertDataList)
	if err != nil {
		return fmt.Errorf("send revert batch: %w", err)
	}

	if err := faultinjector.Check("enygma.handler.Receiver.HandleEnygmaCrossTransfer.after_revert_batch"); err != nil {
		return fmt.Errorf("fault injection at after_revert_batch: %w", err)
	}

	executedReverts := 0
	for _, revertRes := range revertResults {
		if revertRes.Err != nil {
			slog.Error("Failed to send enygma cross transfer revert", slog.Any("error", revertRes.Err))
			continue
		}
		if revertRes.Receipt.Status == 1 {
			executedReverts++
		} else {
			slog.Error("Enygma cross transfer revert failed on-chain")
		}
	}

	slog.Info(
		"Failed transactions were successfully submitted for revert to the sender PL!",
		slog.Int("Total Submitted", executedReverts),
		slog.Int("Total Failed", len(failedMsgIDs)),
	)

	return r.txRecoveryRepository.MarkConfirmed(ctx, batch.PrivateHubTxHash)
}
