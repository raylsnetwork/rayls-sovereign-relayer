package service

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	telemetry "github.com/raylsnetwork/rayls-sovereign-relayer/otel"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type finalizationTracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

type finalizationBlockWaiterService interface {
	WaitForBlock(ctx context.Context, targetBlockNumber uint64) (uint64, error)
}

type finalizationEndpointClient interface {
	GetResourceAddress(ctx context.Context, resourceId string) (common.Address, error)
}

type finalizationRetryService interface {
	RetryOperation(
		ctx context.Context,
		operationName string,
		maxRetries int,
		blockNumber uint64,
		executeOperation func(ctx context.Context, nextBlockNumber uint64) error,
	) (uint64, error)
}

type finalizationExecutor interface {
	ExecuteEnygmaCrossTransfer(
		parentCtx context.Context,
		eventID string,
		blockNumber uint64,
		resourceId string,
		txsByChainID map[string][]*types.EnygmaTransferBatchTx,
		enygmaAddress common.Address,
	) error
}

type EnygmaFinalizationService struct {
	tracer             finalizationTracer
	blockWaiterService finalizationBlockWaiterService
	ccEndpointClient   finalizationEndpointClient
	retryService       finalizationRetryService
	executor           finalizationExecutor
}

func NewEnygmaFinalizationService(
	tracer finalizationTracer,
	blockWaiterService finalizationBlockWaiterService,
	ccEndpointClient finalizationEndpointClient,
	retryService finalizationRetryService,
	executor finalizationExecutor,
) *EnygmaFinalizationService {
	return &EnygmaFinalizationService{
		tracer:             tracer,
		blockWaiterService: blockWaiterService,
		ccEndpointClient:   ccEndpointClient,
		retryService:       retryService,
		executor:           executor,
	}
}

func (s *EnygmaFinalizationService) ExecuteFinalization(
	parentCtx context.Context,
	finalizationID string,
	blockNumber uint64,
	resourceId string,
) error {
	ctx, span := s.tracer.Start(parentCtx, telemetry.SPAN_EXECUTE_FINALIZATION)
	defer span.End()

	span.SetAttributes(
		attribute.String(telemetry.ATTR_RESOURCE_ID, resourceId),
		//nolint:gosec // block numbers are within int range
		attribute.Int(telemetry.ATTR_BLOCK_NUMBER, int(blockNumber)),
	)

	// Ensure that the PNH block number is greater than the block number of the last transfer.
	blockNumber, err := s.blockWaiterService.WaitForBlock(ctx, blockNumber)
	if err != nil {
		return fmt.Errorf("wait for block: %w", err)
	}

	enygmaAddress, err := s.ccEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("get resource address: %w", err)
	}

	maxRetries := 10
	_, err = s.retryService.RetryOperation(
		ctx,
		"finalization",
		maxRetries,
		blockNumber,
		func(ctx context.Context, nextBlockNumber uint64) error {
			return s.executor.ExecuteEnygmaCrossTransfer(ctx, finalizationID, nextBlockNumber, resourceId, nil, enygmaAddress)
		},
	)
	if err != nil {
		return fmt.Errorf("retry finalization: %w", err)
	}

	return nil
}
