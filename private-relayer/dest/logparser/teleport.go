// Decommissioning Teleport (vanilla, atomic): atomic members below marked; shared/generic/Enygma/DVP retained.

package logparser

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/backoff"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/TeleportV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/encrypt"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/wireformat"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/logrouter"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"google.golang.org/grpc"
)

//go:generate moq --pkg logparser_test -out teleport_mock_test.go . TeleportMQ CrossChainMQ CTSClient EthereumClient Backoff SUMService
type TeleportMQ interface {
	Next(ctx context.Context) (msgqueue.Message[logrouter.Block], error)
}

type CrossChainMQ interface {
	PushBatch(context.Context, []types.DispatchedMessageToPrivateHub) error
}

type CTSClient interface {
	Decrypt(
		ctx context.Context,
		in *encrypt.DecryptRequest,
		opts ...grpc.CallOption,
	) (*encrypt.DecryptResponse, error)
}

type EthereumClient interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*ethTypes.Block, error)
}

type Backoff interface {
	Do(ctx context.Context, maxAttempts int, fn func() error) error
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type SUMService interface {
	BatchCreate(ctx context.Context, sums []types.AtomicStatusUpdateMessage) error
}

type TeleportParser struct {
	teleportEventParser *TeleportV1.TeleportV1

	teleportMQ   TeleportMQ
	crossChainMQ CrossChainMQ

	ctsClient CTSClient
	ethClient EthereumClient
	// Decommissioning Teleport (vanilla, atomic).
	sumService SUMService

	localChainID *big.Int

	backoff Backoff
}

func NewTeleportParser(
	teleportMQ TeleportMQ,
	crossChainMQ CrossChainMQ,
	ctsClient CTSClient,
	ethClient EthereumClient,
	localChainID *big.Int,
	sumService SUMService,
) *TeleportParser {
	teleportEventParser := TeleportV1.NewTeleportV1()

	backoff, _ := backoff.NewExponential(time.Second, 2, time.Minute)

	return &TeleportParser{
		teleportEventParser: teleportEventParser,

		teleportMQ:   teleportMQ,
		crossChainMQ: crossChainMQ,

		ctsClient:    ctsClient,
		ethClient:    ethClient,
		localChainID: localChainID,
		sumService:   sumService,

		backoff: backoff,
	}
}

func NewTeleportParserWithCustomBackoff(
	teleportMQ TeleportMQ,
	crossChainMQ CrossChainMQ,

	ctsClient CTSClient,
	ethClient EthereumClient,
	localChainID *big.Int,
	sumService SUMService,

	backoff Backoff,
) *TeleportParser {
	teleportEventParser := TeleportV1.NewTeleportV1()

	return &TeleportParser{
		teleportEventParser: teleportEventParser,

		teleportMQ:   teleportMQ,
		crossChainMQ: crossChainMQ,

		ctsClient:    ctsClient,
		ethClient:    ethClient,
		localChainID: localChainID,
		sumService:   sumService,

		backoff: backoff,
	}
}

func (t *TeleportParser) Run(ctx context.Context) error { //nolint:gocognit // log parser handles multiple event types
	slog.Info("TeleportParser started")
	for {
		slog.Debug("Fetching next teleport block from queue")
		block, err := t.teleportMQ.Next(ctx)
		if err != nil {
			if ctx.Err() != nil {
				slog.Info("TeleportLogParserDecryptor shutting down")
				return nil //nolint:nilerr // context cancellation is a clean shutdown signal
			}
			slog.Error("Failed to get next message from teleport MQ", slog.Any("error", err))
			continue
		}

		slog.Debug(
			"Processing teleport logs",
			slog.Uint64("block_number", block.V.Number),
			slog.Int("log_count", len(block.V.Logs)),
		)

		var hasProcessingError bool
		batchesProcessed := 0
		totalMessages := 0
		for _, log := range block.V.Logs {
			if event, parseErr := t.teleportEventParser.UnpackEncryptedDataBatchStoredEvent(&log); parseErr == nil {
				slog.Debug("Attempting to decrypt message batch", slog.Uint64("block_number", block.V.Number))
				msgSlice, decryptErr := t.decryptMessages(ctx, event, event.BlockNumber)
				if decryptErr != nil {
					slog.Error("Failed to process message batch",
						slog.Uint64("block_number", block.V.Number),
						slog.Any("error", decryptErr))
					hasProcessingError = true
					continue
				}
				if msgSlice == nil {
					continue
				}

				slog.Debug("Successfully decrypted messages", slog.Int("count", len(msgSlice)))

				err = t.backoff.Do(ctx, 10, func() error {
					return t.crossChainMQ.PushBatch(ctx, msgSlice)
				})
				if err != nil {
					slog.Error("Failed to push message batch to Cross Chain MQ",
						slog.Uint64("block_number", block.V.Number),
						slog.Int("message_count", len(msgSlice)),
						slog.Any("error", err))
					hasProcessingError = true
				} else {
					slog.Debug("Successfully pushed message batch to Cross Chain MQ", slog.Int("count", len(msgSlice)))
					batchesProcessed++
					totalMessages += len(msgSlice)
				}
			} else if event, parseErr := t.teleportEventParser.UnpackAtomicMessageStatusChangedBatchEvent(&log); parseErr == nil {
				// Decommissioning Teleport (vanilla, atomic).
				sums := make([]types.AtomicStatusUpdateMessage, 0, len(event.MsgIds))
				for _, sharedID := range event.MsgIds {
					sums = append(sums, types.AtomicStatusUpdateMessage{
						SharedID: sharedID,
						Status:   types.AtomicStatus(event.Status),
					})
				}

				if len(sums) == 0 {
					continue
				}

				slog.Info("Creating Status Update Messages",
					slog.Uint64("block_number", block.V.Number),
					slog.Int("count", len(sums)))

				err = t.backoff.Do(ctx, 10, func() error {
					return t.sumService.BatchCreate(ctx, sums)
				})
				if err != nil {
					slog.Error("Failed to create SUMs",
						slog.Uint64("block_number", block.V.Number),
						slog.Int("count", len(sums)),
						slog.Any("error", err))
					hasProcessingError = true
				}
			} else {
				slog.Debug("Skipping unrecognized log topic",
					slog.Uint64("block_number", block.V.Number))
			}
		}

		if batchesProcessed > 0 {
			slog.Info("Processed teleport batches",
				slog.Uint64("block_number", block.V.Number),
				slog.Int("batch_count", batchesProcessed),
				slog.Int("message_count", totalMessages))
		}

		if hasProcessingError {
			slog.Warn("Skipping ack due to processing errors, block will be redelivered",
				slog.Uint64("block_number", block.V.Number))
			continue
		}

		slog.Debug("Acknowledging teleport block", slog.Uint64("block_number", block.V.Number))
		err = t.backoff.Do(ctx, 10, func() error {
			return block.Ack(ctx)
		})
		if err != nil {
			slog.Error("Failed to acknowledge teleport block",
				slog.Uint64("block_number", block.V.Number),
				slog.Any("error", err))
		}
	}
}

func (t *TeleportParser) decryptMessages(
	ctx context.Context,
	e *TeleportV1.TeleportV1EncryptedDataBatchStored,
	blockNumber *big.Int,
) ([]types.DispatchedMessageToPrivateHub, error) {
	ethBlock, err := t.ethClient.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get block %s: %w", blockNumber, err)
	}

	decryptResp, err := t.ctsClient.Decrypt(ctx, &encrypt.DecryptRequest{
		EncryptedData: e.Data,
		Fingerprint:   e.MessageTag,
		BlockNumber:   blockNumber.Uint64(),
		PrevBlockHash: ethBlock.ParentHash().String(),
	})
	if err != nil {
		return nil, fmt.Errorf("decrypting message batch at block %s: %w", blockNumber, err)
	}
	switch decryptResp.GetOutcome() {
	case encrypt.DecryptOutcome_OUTCOME_NOT_FOR_RECIPIENT:
		slog.Debug("Teleport batch not addressed to us, skipping",
			slog.String("blockNumber", blockNumber.String()))
		//nolint:nilnil // empty result is the intended signal for "not for us"
		return nil, nil
	case encrypt.DecryptOutcome_OUTCOME_TAMPERED:
		// As in the Enygma parser: a tag-matched ciphertext that fails AEAD is either tampered/mis-sealed
		// or a transient secret-propagation race. We must not ack-block the shared stream on a possibly
		// poison batch, so we still skip it (return nil, nil) — but loudly, not silently. Emit a
		// structured, greppable ERROR with the dest block, the seal block, and the tx hash so the drop is
		// alertable and forensically reconcilable instead of an invisible WARN.
		slog.Error("teleport_batch_aead_tampered: Teleport batch fingerprint matched but AEAD verification failed on destination; dropping (NOT silently)",
			slog.String("event", "teleport_batch_aead_tampered"),
			slog.String("blockNumber", blockNumber.String()),
			slog.Uint64("sealBlockNumber", e.BlockNumber.Uint64()),
			slog.String("txHash", e.Raw.TxHash.Hex()),
		)
		//nolint:nilnil // skip the batch but don't ack-block on tampering at this layer
		return nil, nil
	case encrypt.DecryptOutcome_OUTCOME_OK:
		// fall through
	default:
		return nil, fmt.Errorf("decrypting message batch at block %s: unexpected outcome %v", blockNumber, decryptResp.GetOutcome())
	}

	msgs, err := wireformat.UnmarshalPlaintext[[]types.DispatchedMessageToPrivateHub](decryptResp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling message batch at block %s: %w", blockNumber, err)
	}
	return msgs, nil
}
