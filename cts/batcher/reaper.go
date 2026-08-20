package batcher

import (
	"context"
	"log/slog"
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/ethrpc"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

const stuckReason = "stuck: no receipt after threshold"

type ReaperConfig struct {
	Identity       string
	BatchSize      int
	Interval       time.Duration
	StuckThreshold time.Duration // how long a row may sit in 'sent' before being re-broadcast / dead-lettered
	// MaxResendAttempts caps how many total broadcasts (initial + resends)
	// the pipeline performs before dead-lettering. When a row's send_attempts
	// reaches this value the reaper publishes TxResultFailed and marks the
	// row 'failed'. A value of 0 disables resend (reaper behaves as before).
	MaxResendAttempts int
}

type ReaperService struct {
	identity          string
	batchSize         int
	interval          time.Duration
	stuckThreshold    time.Duration
	maxResendAttempts int

	repo      Repository
	sender    Batcher
	publisher ResultPublisher
}

// NewReaperService wires a reaper. `sender` is the broadcast surface used
// to re-send rows that have been stuck past StuckThreshold; pass nil to
// disable resend (the reaper reverts to its pre-resend dead-letter-only
// behaviour). `sender` is the same `ethrpc.Batcher` the `BatcherService`
// holds — sharing it keeps the private-key queue (and therefore the
// nonce) under a single owner.
func NewReaperService(conf ReaperConfig, repo Repository, sender Batcher, pub ResultPublisher) *ReaperService {
	return &ReaperService{
		identity:          conf.Identity,
		batchSize:         conf.BatchSize,
		interval:          conf.Interval,
		stuckThreshold:    conf.StuckThreshold,
		maxResendAttempts: conf.MaxResendAttempts,
		repo:              repo,
		sender:            sender,
		publisher:         pub,
	}
}

func (s *ReaperService) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	initial := make(chan struct{}, 1)
	initial <- struct{}{}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		case <-initial:
		}
		s.reap(ctx)
	}
}

func (s *ReaperService) reap(ctx context.Context) {
	cutoff := time.Now().Add(-s.stuckThreshold)
	rows, err := s.repo.ClaimStuck(ctx, s.identity, cutoff, s.batchSize)
	if err != nil {
		slog.Error("reap: claim stuck failed", slog.String("identity", s.identity), slog.Any("err", err))
		return
	}
	if len(rows) == 0 {
		return
	}

	for _, r := range rows {
		// Resend path: re-broadcast the original calldata before giving
		// up. The chain dedups by hash if the original tx is still in
		// any peer's mempool (the common transient case); accepts a
		// fresh broadcast if the original was dropped. Either way the
		// receipter (which keeps polling the row's tx_hash) eventually
		// finalises the row and the reaper backs off.
		if s.sender != nil && r.SendAttempts < s.maxResendAttempts {
			s.attemptResend(ctx, r)
			continue
		}

		// Dead-letter path: cap reached (or resend disabled).
		result := types.TxResult{
			CorrelationID: r.CorrelationID,
			MessageType:   r.MessageType,
			Identity:      s.identity,
			Kind:          types.TxResultFailed,
			TxHash:        r.TxHash,
			ErrorReason:   stuckReason,
		}
		if err := s.publisher.Push(ctx, result); err != nil {
			slog.Error("reap: publish failed",
				slog.String("correlation_id", r.CorrelationID), slog.Any("err", err))
			continue
		}
		if err := s.repo.MarkFailed(ctx, KeyOf(r), stuckReason); err != nil {
			slog.Error("reap: mark failed failed",
				slog.String("correlation_id", r.CorrelationID), slog.Any("err", err))
		}
	}
}

// attemptResend re-broadcasts a stuck row's original calldata. The row
// stays 'sent' on every path except the pre-flight-revert one (where
// the row is finalised as finished-revert). The receipter keeps polling
// r.TxHash regardless, so a delayed receipt of the first broadcast
// (the common transient case) still resolves the row normally.
//
// Known limitation — truly-dropped txs cannot be revived here.
// `s.sender.Send` reserves a NEW nonce from the shared cache and signs
// a NEW tx (different hash). If the original tx was dropped from every
// peer's mempool — i.e. the chain genuinely lost it, not just delayed
// — the receipter will keep polling the original (now-dead) hash and
// the resend's hash is not tracked anywhere on the row. Net effect:
// MaxResendAttempts is fully effective for transient mempool / worker
// blips, but a truly-dropped tx will cycle through all resend attempts
// and eventually dead-letter. Operators watching for "true drops"
// should alert on `MarkFailed` rows whose final send_attempts equals
// MaxResendAttempts and whose `Re-broadcasting stuck tx` log line
// fired the expected number of times. A future fix would either
// invalidate the nonce cache + persist the resend's hash (replacing
// the original tracking, with care around idempotent on-chain
// rejects), or extend the schema to track multiple hashes per row.
//
// Replay safety across reaped message types. The resend re-broadcasts the same
// calldata at a fresh nonce, so a delayed-not-dropped original and the resend can
// both mine; correctness then depends on the destination being idempotent for that
// message. The reaped set is exactly the async cts_transaction rows, and every
// reaped message type (crosschain.vanilla, crosschain.atomic, atomic.*,
// privatehub.execute) dispatches Endpoint.receivePayload into
// RNMessageExecutorV1.executeMessage, which reverts on a duplicate messageId (the
// on-chain executed[messageId] guard). A second delivery of anything the reaper
// resends is therefore rejected on chain. Enygma cross-transfer and DvP mints are
// NOT reaped: they run through the synchronous BatchSignAndSend path (no
// cts_transaction row) and are additionally guarded on chain by referenceIdsStatus.
// The one unguarded pair (receiveWithdrawFromDvp, supplyUpdateRevert) is tracked in
// rayls-privacy-contracts#272 and is likewise not reaper-reachable.
func (s *ReaperService) attemptResend(ctx context.Context, r CTSTransaction) {
	key := KeyOf(r)
	input := ethrpc.TransactionInput{ID: r.CorrelationID, Data: r.Calldata, Address: r.Address}

	// Bump send_attempts BEFORE broadcast so a mid-resend crash leaves
	// the counter in the bumped state (preserves the cap invariant
	// across restarts). Mirrors BatcherService.send.
	if err := s.repo.IncrementSendAttempts(ctx, []TxKey{key}); err != nil {
		slog.Error("reap: bump send_attempts failed",
			slog.String("identity", s.identity),
			slog.String("correlation_id", r.CorrelationID), slog.Any("err", err))
		return
	}

	slog.Info("Re-broadcasting stuck tx",
		slog.String("service", s.identity),
		slog.String("correlation_id", r.CorrelationID),
		slog.Int("send_attempts", r.SendAttempts+1),
		slog.Int("max_resend_attempts", s.maxResendAttempts))

	results, err := s.sender.Send(ctx, []ethrpc.TransactionInput{input})
	if err != nil {
		// Whole-batch retryable failure — leave the row alone; next
		// reap tick will try again under the bumped send_attempts.
		slog.Warn("reap: resend batch send failed",
			slog.String("identity", s.identity),
			slog.String("correlation_id", r.CorrelationID), slog.Any("err", err))
		return
	}
	if len(results) != 1 {
		// Contract violation from s.sender — Send is documented to
		// return exactly one SendResult per input. A mismatch here
		// means a programming error somewhere in the sender pipeline.
		// Log loudly: send_attempts was already incremented above, so
		// silent return would re-fire on every reap tick with no
		// indication of why.
		slog.Error("reap: sender returned unexpected result count (expected 1)",
			slog.String("identity", s.identity),
			slog.String("correlation_id", r.CorrelationID),
			slog.Int("got", len(results)))
		return
	}

	res := results[0]
	switch {
	case res.Err != nil:
		// Per-tx broadcast error (e.g., "nonce too low" because the
		// original eventually mined). Harmless — the receipter is
		// still polling r.TxHash; it'll see the receipt next tick.
		// Refresh sent_at so the reaper paces resends by
		// StuckThreshold instead of busy-looping every Interval.
		slog.Warn("reap: per-tx resend err",
			slog.String("identity", s.identity),
			slog.String("correlation_id", r.CorrelationID),
			slog.Any("err", res.Err))
		if err := s.repo.MarkResent(ctx, key); err != nil {
			slog.Error("reap: mark resent failed (after per-tx err)",
				slog.String("correlation_id", r.CorrelationID), slog.Any("err", err))
		}
	case res.Revert != nil:
		// Pre-flight revert during RESEND ≠ pre-flight revert during the
		// FIRST send. Here it almost always means the ORIGINAL tx mined
		// successfully and changed chain state in a way that gas
		// estimation of the new (same-calldata) tx now reverts — e.g.,
		// an idempotent executeMessage rejecting a second delivery of
		// the same msgId. We CANNOT tell that apart from "the original
		// would have reverted too" by looking only at the resend's
		// estimation, so finalising the row here would (a) misreport a
		// successful delivery as a revert and (b) close the retry path
		// at the relayer. Safer: leave the row 'sent' and let the
		// receipter (post-FI-expiry) poll the original tx_hash and
		// classify by the actual mined receipt. Refresh sent_at so
		// the reaper paces the next cycle.
		slog.Warn("reap: per-tx resend pre-flight revert (treating as transient — receipter will classify by original receipt)",
			slog.String("identity", s.identity),
			slog.String("correlation_id", r.CorrelationID),
			slog.Int("revert_data_len", len(res.Revert.Data)))
		if err := s.repo.MarkResent(ctx, key); err != nil {
			slog.Error("reap: mark resent failed (after pre-flight revert)",
				slog.String("correlation_id", r.CorrelationID), slog.Any("err", err))
		}
	default:
		// Successful re-broadcast. The chain accepted the new tx, but
		// we intentionally do NOT overwrite r.TxHash on the row — the
		// receipter keeps polling the ORIGINAL hash so a delayed receipt
		// of the first broadcast (the common case) still resolves the
		// row, and a real on-chain executor reject of a duplicate (an
		// idempotent executeMessage) won't masquerade as a fresh
		// revert. Refresh sent_at to pace the next reap cycle.
		if err := s.repo.MarkResent(ctx, key); err != nil {
			slog.Error("reap: mark resent failed",
				slog.String("correlation_id", r.CorrelationID), slog.Any("err", err))
		}
	}
}
