package logparser

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/backoff"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/DvpTeleport"
	encryptpb "github.com/raylsnetwork/rayls-sovereign-relayer/cts/gen/proto/encrypt"
	keyspb "github.com/raylsnetwork/rayls-sovereign-relayer/cts/gen/proto/keys"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/wireformat"
	"github.com/raylsnetwork/rayls-sovereign-relayer/msgqueue"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/logrouter"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"google.golang.org/grpc"
)

//go:generate moq --pkg logparser_test -out dvp_mock_test.go . DvpTeleportMQ DvpOrchestratorMQ DvpCTSClient DvpBackoff DvpEthereumClient DvpSwapRepository

type DvpBackoff interface {
	Do(ctx context.Context, maxAttempts int, fn func() error) error
}

type DvpEthereumClient interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*ethTypes.Block, error)
}

type DvpSwapRepository interface {
	GetSwapBySharedID(ctx context.Context, sharedID string) (*types.DvpSwap, error)
	UpdateSwapTo(ctx context.Context, sharedID string, to string) error
}

type DvpTeleportMQ interface {
	Next(ctx context.Context) (msgqueue.Message[logrouter.Block], error)
}

type DvpOrchestratorMQ interface {
	Push(ctx context.Context, msg service.DvpDestMessage) error
}

type DvpCTSClient interface {
	RecoverViewSalt(ctx context.Context, in *keyspb.RecoverViewSaltRequest, opts ...grpc.CallOption) (*keyspb.RecoverViewSaltResponse, error)
	DecryptWithoutFPWithSS(ctx context.Context, in *encryptpb.DecryptWithoutFPWithSSRequest, opts ...grpc.CallOption) (*encryptpb.DecryptWithoutFPWithSSResponse, error)
}

type DvpTeleportParser struct {
	teleportMQ     DvpTeleportMQ
	orchestratorMQ DvpOrchestratorMQ

	kosClient DvpCTSClient
	ethClient DvpEthereumClient
	swapRepo  DvpSwapRepository

	localChainId *big.Int

	teleportFilterer *DvpTeleport.DvpTeleport
	backoff          DvpBackoff
}

func NewDvpTeleportParser(
	teleportMQ DvpTeleportMQ,
	orchestratorMQ DvpOrchestratorMQ,
	kosClient DvpCTSClient,
	ethClient DvpEthereumClient,
	swapRepo DvpSwapRepository,
	localChainId *big.Int,
) *DvpTeleportParser {
	backoff, _ := backoff.NewExponential(time.Second, 2, time.Minute)
	return NewDvpTeleportParserWithBackoff(
		teleportMQ,
		orchestratorMQ,
		kosClient,
		ethClient,
		swapRepo,
		localChainId,
		backoff,
	)
}

func NewDvpTeleportParserWithBackoff(
	teleportMQ DvpTeleportMQ,
	orchestratorMQ DvpOrchestratorMQ,
	kosClient DvpCTSClient,
	ethClient DvpEthereumClient,
	swapRepo DvpSwapRepository,
	localChainId *big.Int,
	backoff DvpBackoff,
) *DvpTeleportParser {
	teleportFilterer := DvpTeleport.NewDvpTeleport()

	return &DvpTeleportParser{
		teleportMQ:       teleportMQ,
		orchestratorMQ:   orchestratorMQ,
		kosClient:        kosClient,
		ethClient:        ethClient,
		swapRepo:         swapRepo,
		localChainId:     localChainId,
		teleportFilterer: teleportFilterer,
		backoff:          backoff,
	}
}

func (p *DvpTeleportParser) Run( //nolint:gocognit // log parser handles multiple event types
	ctx context.Context,
) error {
	slog.Info("DvpTeleportParser started")
	for {
		slog.Debug("Fetching next dvp teleport block from queue")
		block, err := p.teleportMQ.Next(ctx)
		if err != nil {
			if ctx.Err() != nil {
				slog.Info("DvpTeleportParser shutting down")
				return nil //nolint:nilerr // context cancellation is a clean shutdown signal
			}
			slog.Error("Failed to get next message from dvp teleport MQ", slog.Any("error", err))
			continue
		}

		slog.Debug(
			"Processing dvp teleport block",
			slog.Uint64("blockNumber", block.V.Number),
			slog.Int("log_count", len(block.V.Logs)),
		)

		var hasProcessingError bool
		for _, log := range block.V.Logs {
			if commitmentsEvent, err := p.teleportFilterer.UnpackCommitmentsEvent(&log); err == nil {
				if err := p.handleCommitments(ctx, commitmentsEvent, block.V.Number); err != nil {
					slog.Error("Failed to handle Commitments event",
						slog.Any("error", err),
						slog.Uint64("blockNumber", block.V.Number),
					)
					hasProcessingError = true
				}
			} else if nullifiersEvent, err := p.teleportFilterer.UnpackNullifiersEvent(&log); err == nil {
				if err := p.handleNullifiers(ctx, nullifiersEvent, block.V.Number); err != nil {
					slog.Error("Failed to handle Nullifiers event",
						slog.Any("error", err),
						slog.Uint64("blockNumber", block.V.Number),
					)
					hasProcessingError = true
				}
			} else if swapEvent, err := p.teleportFilterer.UnpackSwapInitiatedEvent(&log); err == nil {
				if err := p.handleSwapInitiated(
					ctx,
					swapEvent,
					swapEvent.Raw.BlockNumber,
				); err != nil {
					slog.Error("Failed to handle SwapInitiated event",
						slog.Any("error", err),
						slog.Uint64("blockNumber", block.V.Number),
					)
					hasProcessingError = true
				}
			} else if completedEvent, err := p.teleportFilterer.UnpackSwapCompletedEvent(&log); err == nil {
				if err := p.handleSwapCompleted(ctx, completedEvent, block.V.Number); err != nil {
					slog.Error("Failed to handle SwapCompleted event",
						slog.Any("error", err),
						slog.Uint64("blockNumber", block.V.Number),
					)
					hasProcessingError = true
				}
			} else if cancelledEvent, err := p.teleportFilterer.UnpackSwapCancelledEvent(&log); err == nil {
				if err := p.handleSwapCancelled(ctx, cancelledEvent, block.V.Number); err != nil {
					slog.Error("Failed to handle SwapCancelled event",
						slog.Any("error", err),
						slog.Uint64("blockNumber", block.V.Number),
					)
					hasProcessingError = true
				}
			} else if timedOutEvent, err := p.teleportFilterer.UnpackSwapTimedOutEvent(&log); err == nil {
				if err := p.handleSwapTimedOut(ctx, timedOutEvent, block.V.Number); err != nil {
					slog.Error("Failed to handle SwapTimedOut event",
						slog.Any("error", err),
						slog.Uint64("blockNumber", block.V.Number),
					)
					hasProcessingError = true
				}
			}
		}

		if hasProcessingError {
			slog.Warn("Skipping ack due to processing errors, block will be redelivered",
				slog.Uint64("blockNumber", block.V.Number))
			continue
		}

		slog.Debug("Acknowledging dvp teleport block", slog.Uint64("blockNumber", block.V.Number))
		if err := p.backoff.Do(ctx, 10, func() error {
			return block.Ack(ctx)
		}); err != nil {
			slog.Error("Failed to acknowledge block after retries",
				slog.Uint64("blockNumber", block.V.Number),
				slog.Any("error", err),
			)
		}
	}
}

func (p *DvpTeleportParser) handleSwapInitiated(
	ctx context.Context,
	event *DvpTeleport.DvpTeleportSwapInitiated,
	blockNumber uint64,
) error {
	sharedId := hex.EncodeToString(event.SharedId[:])

	slog.Info("handleSwapInitiated: Processing swap initiated event",
		slog.Uint64("blockNumber", blockNumber),
		slog.String("sharedId", sharedId),
	)

	saltResp, err := p.kosClient.RecoverViewSalt(ctx,
		&keyspb.RecoverViewSaltRequest{
			BlockNumber: event.Raw.BlockNumber,
			Ctxt:        event.Ctxt,
		})
	if err != nil {
		return fmt.Errorf("recover salt from ciphertext at block %d: %w", blockNumber, err)
	}

	decryptResp, err := p.kosClient.DecryptWithoutFPWithSS(ctx, &encryptpb.DecryptWithoutFPWithSSRequest{
		EncryptedData: event.EncryptedData,
		Ss:            saltResp.GetSalt(),
	})
	if err != nil {
		// Genuine failure (CTS unreachable, malformed event, CTS bug). "Not for me"
		// is now an outcome on the response, not an error.
		return fmt.Errorf("decrypting dvp swap message at block %d: %w", blockNumber, err)
	}
	switch decryptResp.GetOutcome() {
	case encryptpb.DecryptOutcome_OUTCOME_NOT_FOR_RECIPIENT:
		slog.Debug("DvpSwapInitiated not addressed to us, skipping event",
			slog.Uint64("blockNumber", blockNumber),
			slog.String("sharedId", sharedId),
		)
		return nil
	case encryptpb.DecryptOutcome_OUTCOME_TAMPERED:
		slog.Warn("DvpSwapInitiated AEAD verification failed unexpectedly",
			slog.Uint64("blockNumber", blockNumber),
			slog.String("sharedId", sharedId),
		)
		return nil
	case encryptpb.DecryptOutcome_OUTCOME_OK:
		// fall through
	default:
		return fmt.Errorf("decrypting dvp swap message at block %d: unexpected outcome %v", blockNumber, decryptResp.GetOutcome())
	}
	decryptedMsg, err := wireformat.UnmarshalPlaintext[*types.DvpSwapMessage](decryptResp)
	if err != nil {
		return fmt.Errorf("unmarshaling dvp swap message at block %d: %w", blockNumber, err)
	}

	if decryptedMsg != nil {
		messageID := fmt.Sprintf("dvp-initiated-to-us-%d-%s", blockNumber, hex.EncodeToString(event.Raw.TxHash[:8]))
		msg := service.DvpDestMessage{
			ID:          messageID,
			Type:        service.DvpSwapInitiatedMessage,
			BlockNumber: blockNumber,
			SharedID:    sharedId,
			SwapInitiated: &types.DvpSwapInitiatedData{
				Message:           decryptedMsg,
				InitiatorDestSalt: new(big.Int).SetBytes(saltResp.GetSalt()),
				// ExpiresAt:  event.ExpiresAt,
			},
		}
		if err := p.backoff.Do(ctx, 10, func() error {
			return p.orchestratorMQ.Push(ctx, msg)
		}); err != nil {
			return fmt.Errorf("failed to push DvpSwapInitiated to orchestrator: %w", err)
		}
		slog.Info("Pushed DvpSwapInitiated to orchestrator",
			slog.String("messageId", messageID),
			slog.String("sharedId", sharedId),
		)
		return nil
	}

	return nil
}

func (p *DvpTeleportParser) handleCommitments(
	ctx context.Context,
	event *DvpTeleport.DvpTeleportCommitments,
	blockNumber uint64,
) error {
	slog.Info("handleCommitments: Processing commitments event",
		slog.Uint64("blockNumber", blockNumber),
		slog.String("tokenAddress", event.TokenAddress.Hex()),
		slog.String("treeNumber", event.TreeNumber.String()),
		slog.Int("commitmentCount", len(event.Commitments)),
	)

	messageID := fmt.Sprintf(
		"dvp-commitments-%d-%s-%d",
		blockNumber,
		hex.EncodeToString(event.Raw.TxHash[:8]),
		event.Raw.Index,
	)

	msg := service.DvpDestMessage{
		ID:          messageID,
		Type:        service.DvpCommitmentsMessage,
		BlockNumber: blockNumber,
		Commitments: &types.DvpCommitmentsData{
			TokenAddress: event.TokenAddress.Hex(),
			TokenType:    event.TokenType,
			TreeNumber:   event.TreeNumber,
			Commitments:  event.Commitments,
		},
	}

	if err := p.backoff.Do(ctx, 10, func() error {
		return p.orchestratorMQ.Push(ctx, msg)
	}); err != nil {
		return fmt.Errorf("failed to push DvpCommitments to orchestrator: %w", err)
	}

	slog.Info("Pushed DvpCommitments to orchestrator",
		slog.String("messageId", messageID),
		slog.String("tokenAddress", event.TokenAddress.Hex()),
		slog.String("treeNumber", event.TreeNumber.String()),
	)

	return nil
}

func (p *DvpTeleportParser) handleNullifiers(
	ctx context.Context,
	event *DvpTeleport.DvpTeleportNullifiers,
	blockNumber uint64,
) error {
	slog.Info("handleNullifiers: Processing nullifiers event",
		slog.Uint64("blockNumber", blockNumber),
		slog.Int("nullifierCount", len(event.Nullifiers)),
	)

	messageID := fmt.Sprintf(
		"dvp-nullifiers-%d-%s-%d",
		blockNumber,
		hex.EncodeToString(event.Raw.TxHash[:8]),
		event.Raw.Index,
	)

	msg := service.DvpDestMessage{
		ID:          messageID,
		Type:        service.DvpNullifierMessage,
		BlockNumber: blockNumber,
		Nullifiers: &types.DvpNullifierData{
			TokenAddress: event.TokenAddress.Hex(),
			Nullifiers:   event.Nullifiers,
		},
	}

	if err := p.backoff.Do(ctx, 10, func() error {
		return p.orchestratorMQ.Push(ctx, msg)
	}); err != nil {
		return fmt.Errorf("failed to push DvpNullifiers to orchestrator: %w", err)
	}

	slog.Info("Pushed DvpNullifiers to orchestrator",
		slog.String("messageId", messageID),
		slog.Int("nullifierCount", len(event.Nullifiers)),
	)

	return nil
}

func (p *DvpTeleportParser) handleSwapCompleted(
	ctx context.Context,
	event *DvpTeleport.DvpTeleportSwapCompleted,
	blockNumber uint64,
) error {
	sharedId := hex.EncodeToString(event.SharedId[:])

	slog.Info("handleSwapCompleted: Processing swap completed event",
		slog.Uint64("blockNumber", blockNumber),
		slog.String("sharedId", sharedId),
	)

	// Skip events for swaps we don't know about locally.
	swap, err := p.swapRepo.GetSwapBySharedID(ctx, sharedId)
	if err != nil {
		return fmt.Errorf("failed to lookup swap %s for SwapCompleted at block %d: %w", sharedId, blockNumber, err)
	}
	if swap == nil {
		return nil
	}

	// Both sides are receiving this event.
	// Only the initiator needs to decrypt the message to get the to address.
	if swap.To == "" {
		decryptResp, err := p.kosClient.DecryptWithoutFPWithSS(ctx, &encryptpb.DecryptWithoutFPWithSSRequest{
			EncryptedData: event.EncryptedData,
			Ss:            swap.SelfSalt.Bytes(),
		})
		if err != nil {
			return fmt.Errorf("decrypting dvp swap completed message at block %d: %w", blockNumber, err)
		}
		switch decryptResp.GetOutcome() {
		case encryptpb.DecryptOutcome_OUTCOME_NOT_FOR_RECIPIENT:
			slog.Debug("DvpSwapCompleted not addressed to us, skipping event",
				slog.Uint64("blockNumber", blockNumber),
				slog.String("sharedId", sharedId),
			)
			return nil
		case encryptpb.DecryptOutcome_OUTCOME_TAMPERED:
			slog.Warn("DvpSwapCompleted AEAD verification failed unexpectedly",
				slog.Uint64("blockNumber", blockNumber),
				slog.String("sharedId", sharedId),
			)
			return nil
		case encryptpb.DecryptOutcome_OUTCOME_OK:
			// fall through
		default:
			return fmt.Errorf("decrypting dvp swap completed message at block %d: unexpected outcome %v", blockNumber, decryptResp.GetOutcome())
		}
		decryptedMsg, err := wireformat.UnmarshalPlaintext[*types.DvpSwapMessage](decryptResp)
		if err != nil {
			return fmt.Errorf("unmarshaling dvp swap completed message at block %d: %w", blockNumber, err)
		}
		if decryptedMsg == nil {
			return fmt.Errorf("decrypted dvp swap completed message at block %d: nil message", blockNumber)
		}
		if !common.IsHexAddress(decryptedMsg.To) {
			return fmt.Errorf("decrypted dvp swap completed message at block %d has invalid To address %q for swap %s", blockNumber, decryptedMsg.To, sharedId)
		}
		if err := p.swapRepo.UpdateSwapTo(ctx, sharedId, decryptedMsg.To); err != nil {
			return fmt.Errorf("updating swap.To for %s at block %d: %w", sharedId, blockNumber, err)
		}
		slog.Info("Updated swap.To from decrypted SwapCompleted payload", slog.String("sharedId", sharedId))
	}

	messageID := fmt.Sprintf("dvp-completed-%d-%s", blockNumber, hex.EncodeToString(event.Raw.TxHash[:8]))
	msg := service.DvpDestMessage{
		ID:          messageID,
		Type:        service.DvpSwapCompletedMessage,
		BlockNumber: blockNumber,
		SharedID:    sharedId,
	}
	if err := p.backoff.Do(ctx, 10, func() error {
		return p.orchestratorMQ.Push(ctx, msg)
	}); err != nil {
		return fmt.Errorf("failed to push DvpSwapCompleted to orchestrator: %w", err)
	}
	slog.Info("Pushed DvpSwapCompleted to orchestrator",
		slog.String("messageId", messageID),
		slog.String("sharedId", sharedId),
	)
	return nil
}

func (p *DvpTeleportParser) handleSwapCancelled(
	ctx context.Context,
	event *DvpTeleport.DvpTeleportSwapCancelled,
	blockNumber uint64,
) error {
	sharedId := hex.EncodeToString(event.SharedId[:])

	slog.Info("handleSwapCancelled: Processing swap cancelled event",
		slog.Uint64("blockNumber", blockNumber),
		slog.String("sharedId", sharedId),
	)

	forUs, err := p.isEventForUs(ctx, event.SharedId)
	if err != nil {
		return fmt.Errorf("checking if SwapCancelled event is for us: %w", err)
	}
	if !forUs {
		return nil
	}

	messageID := fmt.Sprintf("dvp-cancelled-%d-%s", blockNumber, hex.EncodeToString(event.Raw.TxHash[:8]))
	msg := service.DvpDestMessage{
		ID:          messageID,
		Type:        service.DvpSwapCancelledMessage,
		BlockNumber: blockNumber,
		SharedID:    sharedId,
	}
	if err := p.backoff.Do(ctx, 10, func() error {
		return p.orchestratorMQ.Push(ctx, msg)
	}); err != nil {
		return fmt.Errorf("failed to push DvpSwapCancelled to orchestrator: %w", err)
	}
	slog.Info("Pushed DvpSwapCancelled to orchestrator",
		slog.String("messageId", messageID),
		slog.String("sharedId", sharedId),
	)
	return nil
}

func (p *DvpTeleportParser) handleSwapTimedOut(
	ctx context.Context,
	event *DvpTeleport.DvpTeleportSwapTimedOut,
	blockNumber uint64,
) error {
	sharedId := hex.EncodeToString(event.SharedId[:])

	slog.Info("handleSwapTimedOut: Processing swap timed out event",
		slog.Uint64("blockNumber", blockNumber),
		slog.String("sharedId", sharedId),
	)

	forUs, err := p.isEventForUs(ctx, event.SharedId)
	if err != nil {
		return fmt.Errorf("checking if SwapTimedOut event is for us: %w", err)
	}
	if !forUs {
		return nil
	}

	messageID := fmt.Sprintf("dvp-timedout-%d-%s", blockNumber, hex.EncodeToString(event.Raw.TxHash[:8]))
	msg := service.DvpDestMessage{
		ID:          messageID,
		Type:        service.DvpSwapTimedOutMessage,
		BlockNumber: blockNumber,
		SharedID:    sharedId,
	}
	if err := p.backoff.Do(ctx, 10, func() error {
		return p.orchestratorMQ.Push(ctx, msg)
	}); err != nil {
		return fmt.Errorf("failed to push DvpSwapTimedOut to orchestrator: %w", err)
	}
	slog.Info("Pushed DvpSwapTimedOut to orchestrator",
		slog.String("messageId", messageID),
		slog.String("sharedId", sharedId),
	)
	return nil
}

func (p *DvpTeleportParser) isEventForUs(ctx context.Context, sharedIDBytes [32]byte) (bool, error) {
	sharedID := hex.EncodeToString(sharedIDBytes[:])
	swap, err := p.swapRepo.GetSwapBySharedID(ctx, sharedID)
	if err != nil {
		return false, fmt.Errorf("failed to check swap for shared ID %s: %w", sharedID, err)
	}

	return swap != nil, nil
}
