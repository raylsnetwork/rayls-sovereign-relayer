package resync

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"google.golang.org/grpc"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EnygmaTeleport"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/encrypt"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/wireformat"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/conv"
	enygmaService "github.com/raylsnetwork/rayls-privacy-relayer-api/enygma/service"
	repository "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"go.opentelemetry.io/otel/trace"
)

//go:generate moq --pkg resync -out service_mock_test.go . resourceLockRepository checkpointRepository ethClient kosDecryptor enygmaReceiver tracer

// resourceLockRepository provides distributed locking to prevent concurrent resyncs
type resourceLockRepository interface {
	InsertNewLock(ctx context.Context, resourceId string) error
	RemoveLock(ctx context.Context, resourceId string) error
}

// checkpointRepository provides access to enygma checkpoints
type checkpointRepository interface {
	GetValidationCandidates(ctx context.Context) ([]types.EnygmaCheckpoint, error)
	MarkAsFinalized(ctx context.Context, resourceId string, finalizedBlockNumber *big.Int) error
	GetLatestCheckpointByFilters(
		ctx context.Context,
		resourceId string,
		status *types.EnygmaCheckpointStatus,
		finalizedBlockNumber *big.Int,
		pendingBlockNumber *big.Int,
	) (*types.EnygmaCheckpoint, error)
	CreateEnygmaCheckpoint(ctx context.Context, checkpoint types.EnygmaCheckpoint) error
}

// ethClient provides access to the commit chain for filtering logs
type ethClient interface {
	FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]ethTypes.Log, error)
}

// kosDecryptor decrypts encrypted enygma transfer batches
type kosDecryptor interface {
	DecryptEnygmaTransferBatch(
		ctx context.Context,
		in *encrypt.DecryptEnygmaTransferBatchRequest,
		opts ...grpc.CallOption,
	) (*encrypt.DecryptEnygmaTransferBatchResponse, error)
}

// enygmaReceiver handles the business logic for enygma transfers
type enygmaReceiver interface {
	HandleEnygmaCrossTransfer(ctx context.Context, batch *types.EnygmaTransferBatch) error
}

// tracer provides OpenTelemetry tracing
type tracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

// resyncBlockRange is the number of blocks to fetch when re-syncing enygma events.
const resyncBlockRange = 25

// EnygmaResyncService handles resynchronization of enygma checkpoints
// when validation fails repeatedly. It re-fetches logs from the blockchain
// and re-processes them to recover missing history.
type EnygmaResyncService struct {
	enygmaTeleportAddress common.Address
	plChainId             *big.Int

	resourceLockRepo resourceLockRepository
	checkpointRepo   checkpointRepository
	ethClient        ethClient
	kosClient        kosDecryptor
	receiver         enygmaReceiver
	tracer           tracer

	teleportFilterer *EnygmaTeleport.EnygmaTeleport
}

// NewEnygmaResyncService creates a new resync service with all dependencies
func NewEnygmaResyncService(
	enygmaTeleportAddress common.Address,
	plChainId *big.Int,
	resourceLockRepo resourceLockRepository,
	checkpointRepo checkpointRepository,
	ethClient ethClient,
	kosClient kosDecryptor,
	receiver enygmaReceiver,
	tracer tracer,
) (*EnygmaResyncService, error) {
	teleportFilterer := EnygmaTeleport.NewEnygmaTeleport()

	return &EnygmaResyncService{
		enygmaTeleportAddress: enygmaTeleportAddress,
		plChainId:             plChainId,
		resourceLockRepo:      resourceLockRepo,
		checkpointRepo:        checkpointRepo,
		ethClient:             ethClient,
		kosClient:             kosClient,
		receiver:              receiver,
		tracer:                tracer,
		teleportFilterer:      teleportFilterer,
	}, nil
}

// ResyncEnygma attempts to recover missing enygma history for a resource.
// It acquires a distributed lock, finds the latest checkpoint, re-fetches
// logs from the blockchain, and re-processes them through the receiver.
func (s *EnygmaResyncService) ResyncEnygma(ctx context.Context, resourceId string) error {
	slog.Debug("Checkpoint resync was initialized", "resourceId", resourceId)

	err := s.resourceLockRepo.InsertNewLock(ctx, resourceId)
	if err != nil {
		if errors.Is(err, repository.ErrLockAlreadyExists) {
			slog.DebugContext(
				ctx,
				"Resource lock already exists and has not expired - skipping resync",
				"resourceId",
				resourceId,
			)
			return nil
		}
		return fmt.Errorf("acquiring resource lock for resync %s: %w", resourceId, err)
	}

	// Always release the lock when done, even on error
	defer func() {
		if unlockErr := s.resourceLockRepo.RemoveLock(ctx, resourceId); unlockErr != nil {
			slog.Error(
				"Failed to remove resource lock",
				slog.String("resourceId", resourceId),
				slog.String("err", unlockErr.Error()),
			)
		}
	}()

	return s.processEnygmaResync(ctx, resourceId)
}

func (s *EnygmaResyncService) processEnygmaResync(ctx context.Context, resourceId string) error {
	status := types.EnygmaCheckpointStatusFinal
	checkpoint, err := s.checkpointRepo.GetLatestCheckpointByFilters(ctx, resourceId, &status, nil, nil)
	if err != nil {
		slog.Error(
			"Failed to get latest finalized checkpoint",
			slog.String("resourceId", resourceId),
			slog.String("err", err.Error()),
		)
		return fmt.Errorf("getting latest finalized checkpoint for resource %s: %w", resourceId, err)
	}

	var blockNumberToUse *big.Int

	if checkpoint == nil {
		status = types.EnygmaCheckpointStatusTentative
		checkpoint, err = s.checkpointRepo.GetLatestCheckpointByFilters(ctx, resourceId, &status, nil, nil)
		if err != nil {
			slog.Error(
				"Failed to get latest checkpoint",
				slog.String("resourceId", resourceId),
				slog.String("err", err.Error()),
			)
			return fmt.Errorf("getting latest tentative checkpoint for resource %s: %w", resourceId, err)
		}

		if checkpoint == nil {
			slog.Error("No checkpoint exists for resync", slog.String("resourceId", resourceId))
			return errors.New("no checkpoint found on resync, something is wrong")
		}
		blockNumberToUse = checkpoint.FinalizedBlockNumber
	} else {
		blockNumberToUse = checkpoint.PendingBlockNumber
	}

	return s.refetchEnygmaEvents(ctx, resourceId, blockNumberToUse)
}

func (s *EnygmaResyncService) refetchEnygmaEvents(
	ctx context.Context,
	resourceId string,
	blockNumberToFetchEvents *big.Int,
) error {
	query := ethereum.FilterQuery{
		FromBlock: blockNumberToFetchEvents,
		ToBlock:   new(big.Int).Add(blockNumberToFetchEvents, big.NewInt(resyncBlockRange)),
		Addresses: []common.Address{s.enygmaTeleportAddress},
	}

	logs, err := s.ethClient.FilterLogs(ctx, query)
	if err != nil {
		slog.Error(
			"Failed to filter logs for enygma resync",
			slog.String("resourceId", resourceId),
			slog.String("err", err.Error()),
		)
		return fmt.Errorf("filtering logs for enygma resync of resource %s: %w", resourceId, err)
	}

	return s.parseEnygmaLogs(ctx, logs)
}

// parseEnygmaLogs processes ONLY enygma-related events from logs.
// This is a focused version that doesn't process Endpoint, Teleport, or ZkDVP events.
func (s *EnygmaResyncService) parseEnygmaLogs(ctx context.Context, logs []ethTypes.Log) error {
	var enygmaTransferBatches []*types.EnygmaTransferBatch
	var enygmaFinalizedBalances []*types.EnygmaFinalizedBalance

	for _, log := range logs {
		// Try to parse as EnygmaTransfer
		if transferEvent, err := s.teleportFilterer.UnpackEnygmaTransferEvent(&log); err == nil {
			if transferEvent.ToChainId.Cmp(s.plChainId) != 0 {
				continue
			}

			decryptResp, err := s.kosClient.DecryptEnygmaTransferBatch(
				ctx,
				&encrypt.DecryptEnygmaTransferBatchRequest{
					EncryptedData: transferEvent.EncryptedMessage,
					MessageTag:    transferEvent.MessageTag.Bytes(),
					AnonymitySet:  conv.BigIntsToUint64s(transferEvent.AnonymitySet),
					BlockNumber:   transferEvent.BlockNumber.Uint64(),
					ChainID:       s.plChainId.Uint64(),
					ResourceID:    transferEvent.ResourceId[:],
				},
			)
			if err != nil {
				slog.Warn("Failed to decrypt EnygmaTransfer batch during resync, skipping",
					slog.Any("error", err),
					slog.String("resourceId", hex.EncodeToString(transferEvent.ResourceId[:])),
				)
				continue
			}
			switch decryptResp.GetOutcome() {
			case encrypt.DecryptOutcome_OUTCOME_NOT_FOR_RECIPIENT:
				slog.Debug("EnygmaTransfer batch during resync not addressed to us, skipping",
					slog.String("resourceId", hex.EncodeToString(transferEvent.ResourceId[:])),
				)
				continue
			case encrypt.DecryptOutcome_OUTCOME_TAMPERED:
				slog.Warn("EnygmaTransfer batch during resync AEAD verification failed",
					slog.String("resourceId", hex.EncodeToString(transferEvent.ResourceId[:])),
				)
				continue
			case encrypt.DecryptOutcome_OUTCOME_OK:
				// fall through
			default:
				slog.Warn("EnygmaTransfer batch during resync returned unexpected outcome, skipping",
					slog.String("resourceId", hex.EncodeToString(transferEvent.ResourceId[:])),
					slog.Any("outcome", decryptResp.GetOutcome()),
				)
				continue
			}
			batch, err := wireformat.UnmarshalPlaintext[types.EnygmaTransferBatch](decryptResp)
			if err != nil {
				slog.Warn("Failed to unmarshal EnygmaTransfer batch during resync, skipping",
					slog.Any("error", err),
					slog.String("resourceId", hex.EncodeToString(transferEvent.ResourceId[:])),
				)
				continue
			}

			batch.Ctx = ctx
			batch.PrivateHubTxHash = log.TxHash.Hex()
			enygmaTransferBatches = append(enygmaTransferBatches, &batch)
			continue
		}

		// Try to parse as BalancesFinalized
		if balancesEvent, err := s.teleportFilterer.UnpackBalancesFinalizedEvent(&log); err == nil {
			for _, balance := range balancesEvent.Balances {
				if balance.ChainId.Cmp(s.plChainId) == 0 {
					enygmaFinalizedBalances = append(enygmaFinalizedBalances, &types.EnygmaFinalizedBalance{
						ResourceId:           hex.EncodeToString(balancesEvent.ResourceId[:]),
						FinalizedBlockNumber: balancesEvent.FinalizedBlockNumber,
						PendingBlockNumber:   balancesEvent.PendingBlockNumber,
						BalanceX:             balance.C1,
						BalanceY:             balance.C2,
					})
				}
			}
		}
	}

	// Process transfer batches through receiver (idempotent - will skip existing history)
	for _, batch := range enygmaTransferBatches {
		if err := s.receiver.HandleEnygmaCrossTransfer(ctx, batch); err != nil {
			slog.Error("Error executing Enygma cross transfer batch during resync", slog.String("err", err.Error()))
		}
	}

	// Process finalized balances (idempotent - will skip existing checkpoints)
	if err := enygmaService.ProcessEnygmaFinalizedBalances(
		ctx,
		enygmaFinalizedBalances,
		s.checkpointRepo,
		s.tracer,
	); err != nil {
		slog.Error("Error processing Enygma finalized balances during resync", slog.String("err", err.Error()))
	}

	return nil
}
