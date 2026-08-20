package service

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contractclient"
	"github.com/raylsnetwork/rayls-sovereign-relayer/faultinjector"
	privatehubservice "github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/source/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type revertTracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

type revertEndpointClient interface {
	GetResourceAddress(ctx context.Context, resourceId string) (common.Address, error)
}

type revertEnygmaHandlerClient interface {
	RevertSrcTransferBatch(ctx context.Context, revertTxs []*types.EnygmaTransferFailed) (map[string]contractclient.BatchResult, error)
	RevertSrcSupplyBatch(ctx context.Context, revertTxs []*types.EnygmaSupplyUpdateFailed) (map[string]contractclient.BatchResult, error)
}

type EnygmaRevertService struct {
	tracer              revertTracer
	plEndpointClient    revertEndpointClient
	enygmaHandlerClient revertEnygmaHandlerClient
}

func NewEnygmaRevertService(
	tracer revertTracer,
	plEndpointClient revertEndpointClient,
	enygmaHandlerClient revertEnygmaHandlerClient,
) *EnygmaRevertService {
	return &EnygmaRevertService{
		tracer:              tracer,
		plEndpointClient:    plEndpointClient,
		enygmaHandlerClient: enygmaHandlerClient,
	}
}

// RevertEnygmaSupplyUpdate reverts Enygma supply updates (mints and burns) back to the user
func (s *EnygmaRevertService) RevertEnygmaSupplyUpdate(
	ctx context.Context,
	resourceId string,
	supplyUpdateEvents []privatehubservice.EnygmaSupplyUpdate,
) error {
	ctx, span := s.tracer.Start(ctx, "revert_enygma_supply_update")
	defer span.End()

	span.SetAttributes(
		attribute.String("resource_id", resourceId),
		attribute.Int("supply updates", len(supplyUpdateEvents)),
	)

	enygmaAddress, err := s.plEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("get resource address for supply update revert: %w", err)
	}

	revertTxs := make([]*types.EnygmaSupplyUpdateFailed, 0, len(supplyUpdateEvents))

	for _, event := range supplyUpdateEvents {
		var (
			evType types.EnygmaEventType
			amount *big.Int
		)
		if event.Amount.Cmp(big.NewInt(0)) >= 0 {
			evType = types.EnygmaMint
			amount = event.Amount
		} else {
			evType = types.EnygmaBurn
			amount = big.NewInt(0).Abs(event.Amount)
		}

		revertTx := &types.EnygmaSupplyUpdateFailed{
			TxHash:        event.TxHash.String(),
			EnygmaAddress: enygmaAddress,
			Amount:        amount,
			To:            event.To,
			Type:          evType,
		}
		revertTxs = append(revertTxs, revertTx)
	}

	sendResults, err := s.enygmaHandlerClient.RevertSrcSupplyBatch(ctx, revertTxs)
	if err != nil {
		return fmt.Errorf("send supply update revert batch: %w", err)
	}

	failedTxsCount := 0
	successTxsCount := 0

	for _, sendRes := range sendResults {
		if sendRes.Err != nil {
			failedTxsCount++
			slog.Error("Failed to send enygma revert tx", slog.Any("error", sendRes.Err))
		} else {
			successTxsCount++
		}
	}

	slog.Info(
		"Reverted Enygma supply updates in PL",
		slog.Int("successTxsCount", successTxsCount),
		slog.Int("failedTxsCount", failedTxsCount),
	)

	return nil
}

// RevertEnygmaTransfer reverts failed Enygma cross-chain transfers back to the senders
func (s *EnygmaRevertService) RevertEnygmaTransfer(
	parentCtx context.Context,
	resourceId string,
	txsByChainID map[string][]*types.EnygmaTransferBatchTx,
) error {
	ctx, span := s.tracer.Start(parentCtx, "revert_enygma_transfer")
	defer span.End()

	totalTxs := 0
	for _, txs := range txsByChainID {
		totalTxs += len(txs)
	}
	transfers := make([]*types.EnygmaTransferBatchTx, 0, totalTxs)
	for _, txs := range txsByChainID {
		transfers = append(transfers, txs...)
	}

	span.SetAttributes(
		attribute.String("resource_id", resourceId),
		attribute.Int("transfers_length", len(transfers)),
	)

	senderToFailedTransfer := map[string]*types.EnygmaTransferFailed{}

	for _, transfer := range transfers {
		sender := transfer.FromAddress.Hex()

		if _, exists := senderToFailedTransfer[sender]; !exists {
			senderToFailedTransfer[sender] = &types.EnygmaTransferFailed{
				ReferenceID: transfer.ReferenceId,
				Sender:      transfer.FromAddress,
				Amount:      big.NewInt(0),
				Reason:      "failed to execute transfer on CC",
			}
		}
		senderToFailedTransfer[sender].Amount.Add(senderToFailedTransfer[sender].Amount, transfer.ToAmount)
	}

	revertTxs := make([]*types.EnygmaTransferFailed, 0, len(senderToFailedTransfer))
	for _, tx := range senderToFailedTransfer {
		revertTxs = append(revertTxs, tx)
	}

	err := s.revertTransfersInPL(ctx, resourceId, revertTxs)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed_to_revert_transfer_batch")
		return fmt.Errorf("revert transfers in PL: %w", err)
	}

	slog.Debug("Enygma transfer reverted", slog.String("resourceId", resourceId))

	span.SetStatus(codes.Ok, "success_transfers_reverted")
	return nil
}

// RevertEnygmaDeposit reverts a failed Enygma deposit operation
func (s *EnygmaRevertService) RevertEnygmaDeposit(
	ctx context.Context,
	resourceId string,
	referenceId [32]byte,
	from common.Address,
	amount *big.Int,
) error {
	err := s.revertTransfersInPL(ctx, resourceId, []*types.EnygmaTransferFailed{
		{
			ReferenceID: referenceId,
			Sender:      from,
			Amount:      amount,
			Reason:      "failed to execute deposit on CC",
		},
	})
	if err != nil {
		return fmt.Errorf("revert deposit in PL: %w", err)
	}

	slog.Debug("Enygma deposit reverted", slog.String("resourceId", resourceId))

	return nil
}

// revertTransfersInPL is a private helper that sends revert transactions to the Private Ledger
func (s *EnygmaRevertService) revertTransfersInPL(
	parentCtx context.Context,
	resourceId string,
	revertTxs []*types.EnygmaTransferFailed,
) error {
	_, span := s.tracer.Start(parentCtx, "revert_transfers_in_pl")
	defer span.End()

	span.SetAttributes(
		attribute.String("resource_id", resourceId),
		attribute.Int("revert_txs_count", len(revertTxs)),
	)

	enygmaAddressPL, err := s.plEndpointClient.GetResourceAddress(parentCtx, resourceId)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed_to_get_resource_address")
		return fmt.Errorf("get resource address for revert: %w", err)
	}

	span.SetAttributes(
		attribute.String("enygma_address", enygmaAddressPL.Hex()),
	)

	for _, tx := range revertTxs {
		tx.EnygmaAddress = enygmaAddressPL
	}

	sendResults, err := s.enygmaHandlerClient.RevertSrcTransferBatch(parentCtx, revertTxs)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed_to_send_revert_batch")
		return fmt.Errorf("send transfer revert batch: %w", err)
	}

	if err := faultinjector.Check("enygma.service.EnygmaRevertService.revertTransfersInPL.after_send"); err != nil {
		return fmt.Errorf("fault injection at after_send: %w", err)
	}

	failedTxsCount := 0
	successTxsCount := 0

	for _, sendRes := range sendResults {
		if sendRes.Err != nil {
			failedTxsCount++
			span.AddEvent("revert_tx_failed", trace.WithAttributes(
				attribute.String("error_message", sendRes.Err.Error()),
			))
			slog.Error("Failed to send enygma revert tx", slog.Any("error", sendRes.Err))
		} else {
			successTxsCount++
		}
	}

	span.SetAttributes(
		attribute.Int("success_txs_count", successTxsCount),
		attribute.Int("failed_txs_count", failedTxsCount),
	)

	slog.Info(
		"Reverted Enygma transactions in PL",
		slog.Int("successTxsCount", successTxsCount),
		slog.Int("failedTxsCount", failedTxsCount),
	)

	span.SetStatus(codes.Ok, "success_reverts_sent_to_pl")
	return nil
}
