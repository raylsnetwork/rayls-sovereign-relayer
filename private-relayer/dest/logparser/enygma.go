package logparser

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/backoff"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EnygmaTeleport"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/conv"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/encrypt"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/wireformat"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/logrouter"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"google.golang.org/grpc"
)

//go:generate moq --pkg logparser_test -out enygma_mock_test.go . EnygmaTeleportMQ EnygmaOrchestratorMQ KOSDecryptor EnygmaBackoff

type EnygmaBackoff interface {
	Do(ctx context.Context, maxAttempts int, fn func() error) error
}

type EnygmaTeleportMQ interface {
	Next(ctx context.Context) (msgqueue.Message[logrouter.Block], error)
}

type EnygmaOrchestratorMQ interface {
	Push(ctx context.Context, msg service.EnygmaDestMessage) error
}

type KOSDecryptor interface {
	DecryptEnygmaTransferBatch(
		ctx context.Context,
		in *encrypt.DecryptEnygmaTransferBatchRequest,
		opts ...grpc.CallOption,
	) (*encrypt.DecryptEnygmaTransferBatchResponse, error)
}

type EnygmaTeleportParser struct {
	teleportMQ     EnygmaTeleportMQ
	orchestratorMQ EnygmaOrchestratorMQ

	kosClient KOSDecryptor

	localChainId *big.Int

	teleportFilterer *EnygmaTeleport.EnygmaTeleport
	backoff          EnygmaBackoff
}

func NewEnygmaTeleportParser(
	teleportMQ EnygmaTeleportMQ,
	orchestratorMQ EnygmaOrchestratorMQ,
	kosClient KOSDecryptor,
	localChainId *big.Int,
) *EnygmaTeleportParser {
	backoff, _ := backoff.NewExponential(time.Second, 2, time.Minute)
	return NewEnygmaTeleportParserWithBackoff(teleportMQ, orchestratorMQ, kosClient, localChainId, backoff)
}

func NewEnygmaTeleportParserWithBackoff(
	teleportMQ EnygmaTeleportMQ,
	orchestratorMQ EnygmaOrchestratorMQ,
	kosClient KOSDecryptor,
	localChainId *big.Int,
	backoff EnygmaBackoff,
) *EnygmaTeleportParser {
	teleportFilterer := EnygmaTeleport.NewEnygmaTeleport()

	return &EnygmaTeleportParser{
		teleportMQ:       teleportMQ,
		orchestratorMQ:   orchestratorMQ,
		kosClient:        kosClient,
		localChainId:     localChainId,
		teleportFilterer: teleportFilterer,
		backoff:          backoff,
	}
}

func (p *EnygmaTeleportParser) Run(ctx context.Context) error {
	slog.Info("EnygmaTeleportParser started")
	for {
		slog.Debug("Fetching next enygma teleport block from queue")
		block, err := p.teleportMQ.Next(ctx)
		if err != nil {
			if ctx.Err() != nil {
				slog.Info("EnygmaTeleportParser shutting down")
				return nil //nolint:nilerr // context cancellation is a clean shutdown signal
			}
			slog.Error("Failed to get next message from enygma teleport MQ", slog.Any("error", err))
			continue
		}

		slog.Debug(
			"Processing enygma teleport block",
			slog.Uint64("blockNumber", block.V.Number),
			slog.Int("log_count", len(block.V.Logs)),
		)

		var hasProcessingError bool
		for _, log := range block.V.Logs {
			if transferEvent, err := p.teleportFilterer.UnpackEnygmaTransferEvent(&log); err == nil {
				if err := p.handleEnygmaTransfer(ctx, transferEvent, block.V.Number); err != nil {
					slog.Error("Failed to handle EnygmaTransfer event",
						slog.Any("error", err),
						slog.Uint64("blockNumber", block.V.Number),
					)
					hasProcessingError = true
				}
			} else if balancesEvent, err := p.teleportFilterer.UnpackBalancesFinalizedEvent(&log); err == nil {
				if err := p.handleBalancesFinalized(ctx, balancesEvent, block.V.Number); err != nil {
					slog.Error("Failed to handle BalancesFinalized event",
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

		slog.Debug("Acknowledging enygma teleport block", slog.Uint64("blockNumber", block.V.Number))
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

func (p *EnygmaTeleportParser) handleEnygmaTransfer(
	ctx context.Context,
	event *EnygmaTeleport.EnygmaTeleportEnygmaTransfer,
	blockNumber uint64, // BlockNumber where event was mined
) error {
	slog.Info("handleEnygmaTransfer: Processing transfer event",
		slog.Uint64("blockNumber", blockNumber),
		slog.String("toChainId", event.ToChainId.String()),
		slog.String("localChainId", p.localChainId.String()),
		slog.String("resourceId", hex.EncodeToString(event.ResourceId[:])),
	)

	if event.ToChainId.Cmp(p.localChainId) != 0 {
		slog.Info("Skipping EnygmaTransfer not destined for local chain",
			slog.String("toChainId", event.ToChainId.String()),
			slog.String("localChainId", p.localChainId.String()),
		)
		return nil
	}

	decryptResp, err := p.kosClient.DecryptEnygmaTransferBatch(
		ctx,
		&encrypt.DecryptEnygmaTransferBatchRequest{
			EncryptedData: event.EncryptedMessage,
			MessageTag:    event.MessageTag.Bytes(),
			AnonymitySet:  conv.BigIntsToUint64s(event.AnonymitySet),
			BlockNumber:   event.BlockNumber.Uint64(),
			ChainID:       p.localChainId.Uint64(),
			ResourceID:    event.ResourceId[:],
		},
	)
	if err != nil {
		return fmt.Errorf("decrypting EnygmaTransfer batch at block %d: %w", blockNumber, err)
	}
	switch decryptResp.GetOutcome() {
	case encrypt.DecryptOutcome_OUTCOME_NOT_FOR_RECIPIENT:
		slog.Debug("EnygmaTransfer batch not addressed to us, skipping event",
			slog.Uint64("blockNumber", blockNumber),
			slog.String("resourceId", hex.EncodeToString(event.ResourceId[:])),
		)
		return nil
	case encrypt.DecryptOutcome_OUTCOME_TAMPERED:
		// A tag-matched ciphertext that fails AEAD verification is either a genuinely tampered/mis-sealed
		// batch or a transient secret-propagation race. We deliberately do NOT ack-block the shared dest
		// stream on it (a poison block must not head-of-line-stall every other token's transfers), so this
		// still returns nil and the block is acked. But the drop must NOT be silent: emit a structured,
		// greppable ERROR (severity, not WARN) carrying enough context to alert on and forensically
		// reconcile/replay the source emission — resourceId, the dest block, the seal block the source
		// used, the PNH tx hash, and the chain IDs. Alerting rules should fire on
		// event=enygma_transfer_aead_tampered.
		slog.Error("enygma_transfer_aead_tampered: EnygmaTransfer batch failed AEAD verification on destination; dropping (NOT silently)",
			slog.String("event", "enygma_transfer_aead_tampered"),
			slog.Uint64("blockNumber", blockNumber),
			slog.Uint64("sealBlockNumber", event.BlockNumber.Uint64()),
			slog.String("resourceId", hex.EncodeToString(event.ResourceId[:])),
			slog.String("privateHubTxHash", event.Raw.TxHash.Hex()),
			slog.String("toChainId", event.ToChainId.String()),
			slog.String("localChainId", p.localChainId.String()),
		)
		return nil
	case encrypt.DecryptOutcome_OUTCOME_OK:
		// fall through
	default:
		return fmt.Errorf("decrypting EnygmaTransfer batch at block %d: unexpected outcome %v", blockNumber, decryptResp.GetOutcome())
	}
	decryptedBatch, err := wireformat.UnmarshalPlaintext[types.EnygmaTransferBatch](decryptResp)
	if err != nil {
		return fmt.Errorf("unmarshaling EnygmaTransfer batch at block %d: %w", blockNumber, err)
	}

	decryptedBatch.PrivateHubTxHash = event.Raw.TxHash.Hex()

	messageID := fmt.Sprintf("enygma-transfer-%d-%s", blockNumber, hex.EncodeToString(event.Raw.TxHash[:8]))

	msg := service.EnygmaDestMessage{
		ID:            messageID,
		Type:          service.EnygmaTransferBatchMessage,
		BlockNumber:   blockNumber,
		TransferBatch: &decryptedBatch,
	}

	if err := p.backoff.Do(ctx, 10, func() error {
		return p.orchestratorMQ.Push(ctx, msg)
	}); err != nil {
		return fmt.Errorf("failed to push EnygmaTransferBatch to orchestrator: %w", err)
	}

	slog.Info("Pushed EnygmaTransferBatch to orchestrator",
		slog.String("messageId", messageID),
		slog.String("resourceId", decryptedBatch.ResourceId),
	)

	return nil
}

func (p *EnygmaTeleportParser) handleBalancesFinalized(
	ctx context.Context,
	event *EnygmaTeleport.EnygmaTeleportBalancesFinalized,
	blockNumber uint64,
) error {
	slog.Debug("Processing BalancesFinalized event",
		slog.Uint64("blockNumber", blockNumber),
		slog.Int("balance_count", len(event.Balances)),
		slog.String("resourceId", hex.EncodeToString(event.ResourceId[:8])))

	for _, balance := range event.Balances {
		if balance.ChainId.Cmp(p.localChainId) != 0 {
			continue
		}

		// The message ID must identify the FINALIZATION by its settlement window
		// (resourceId, chainId, finalizedBlockNumber, pendingBlockNumber), NOT the commit-chain block
		// it was emitted in. The emit block is the wrong identity axis: BalancesFinalized carries its
		// own window in the event, and that window — not the block the tx happened to land in — is what
		// the service-level dedup in ProcessEnygmaFinalizedBalances keys on (GetLatestCheckpointByFilters
		// on finalizedBlockNumber + pendingBlockNumber). Keying the MQ ID on the emit block instead
		// desynchronises the two dedup layers: on a resync/re-fetch the same finalization can be re-read
		// under a different block context (false miss → duplicate checkpoint attempt), and two distinct
		// finalizations for one resource emitted in the same block collapse onto one ID (false collision →
		// one silently dropped). Either way the orchestrator MQ dedup (jetstream.WithMsgID + 1h Duplicates
		// window) makes the wrong call, the affected checkpoint is never created, the DB falls a window
		// behind the chain, and every later transfer proof fails the SC-vs-DB commitment guard (deadlock).
		// Keying on the settlement window aligns the MQ ID with the service dedup and still dedups genuine
		// redeliveries/resync re-fetches of the same finalization.
		messageID := fmt.Sprintf("enygma-finalized-%s-%s-%s-%s",
			event.FinalizedBlockNumber.String(),
			event.PendingBlockNumber.String(),
			hex.EncodeToString(event.ResourceId[:8]),
			balance.ChainId.String(),
		)

		finalizedBalance := &types.EnygmaFinalizedBalance{
			ResourceId:           hex.EncodeToString(event.ResourceId[:]),
			FinalizedBlockNumber: event.FinalizedBlockNumber,
			PendingBlockNumber:   event.PendingBlockNumber,
			BalanceX:             balance.C1,
			BalanceY:             balance.C2,
		}

		msg := service.EnygmaDestMessage{
			ID:               messageID,
			Type:             service.EnygmaFinalizedBalanceMessage,
			BlockNumber:      blockNumber,
			FinalizedBalance: finalizedBalance,
		}

		if err := p.backoff.Do(ctx, 10, func() error {
			return p.orchestratorMQ.Push(ctx, msg)
		}); err != nil {
			return fmt.Errorf("failed to push EnygmaFinalizedBalance to orchestrator: %w", err)
		}

		slog.Info("Pushed EnygmaFinalizedBalance to orchestrator",
			slog.String("messageId", messageID),
			slog.String("resourceId", finalizedBalance.ResourceId),
		)
	}

	return nil
}
