// Regression tests for the "stuck tx is dead-lettered without any resend"
// architecture issue.
//
// CONTEXT
// -------
// During the 60-message e2e test, a single dispatch transaction submitted by
// CTS-B was lost by PN B's local axyl network (a consensus stall dropped it
// from the mempool). The CTS reaper claimed the row after StuckThreshold,
// emitted TxResultFailed, and marked the row 'failed'. **No re-broadcast or
// resend was ever attempted.** The relayer then collapsed TxResultFailed into
// OutcomeReverted (terminal), so the message was permanently lost.
//
// The bug investigated here is the lack of *any* resend before dead-letter.
// A robust pipeline would re-broadcast the same signed tx N times before
// giving up — re-broadcasts are cheap (the chain dedups by hash if the tx is
// still in mempool, accepts if dropped).
//
// REGRESSION TESTS
// ----------------
// These tests describe the *correct* behaviour for the reaper:
//
//   - TestReaper_StuckRow_BeforeDeadLetter_AttemptsResend
//       When a row sits in 'sent' past StuckThreshold, the reaper MUST first
//       attempt a re-broadcast (`IncrementSendAttempts` + `Batcher.Send`)
//       before publishing TxResultFailed and marking the row 'failed'.
//
//   - TestReaper_ResendAttemptsRespectMaxBound
//       Resends are bounded: after `send_attempts >= MaxResendAttempts`, the
//       reaper proceeds to dead-letter. This prevents infinite resend loops
//       when the chain is genuinely down.
//
//   - TestReaper_FirstResendObservedThroughIncrementSendAttempts
//       Documents that `IncrementSendAttempts` is the contract signal for "a
//       resend attempt happened" — the existing `reap()` never calls it.
//
// Both tests FAIL today against the current reaper.go (which only marks
// failed and never resends) — see assertion messages for the observed-vs-
// expected behaviour they print.
//
// HOW TO MAKE THEM PASS
// ---------------------
// 1. Extend ReaperConfig with `MaxResendAttempts int` and `Sender Batcher`.
// 2. In reap(), before MarkFailed, if r.SendAttempts < MaxResendAttempts,
//    call sender.Send with the original calldata, IncrementSendAttempts,
//    UpdateSentAt(now), and skip the failed-publish/MarkFailed path for this
//    cycle. Only dead-letter when send_attempts has reached the cap.
//
// Until that fix lands these tests stay red — that's the point.

package batcher_test

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/batcher"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/ethrpc"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/stretchr/testify/assert"
)

// stuckRow returns a row representing a tx that was broadcast but never got
// a receipt (the "axyl dropped it" scenario).
func stuckRow(correlationID string) batcher.CTSTransaction {
	oneHourAgo := time.Now().Add(-time.Hour)
	return batcher.CTSTransaction{
		CorrelationID: correlationID,
		Identity:      "privatenode",
		MessageType:   "crosschain.vanilla",
		Address:       common.HexToAddress("0xf836db0e168b91c568a2fbac5aa79759ea1b0258"),
		Calldata:      []byte{0xfe, 0x12, 0xd9, 0xc3, 0xde, 0xad, 0xbe, 0xef}, // executeMessage(...)
		Status:        batcher.StatusSent,
		TxHash:        common.HexToHash("0xa064c4284c57e4421ef96d2cf61d1537a85555f62ef500e0f0ea3e0fb25d3c4a"),
		SendAttempts:  1,
		SentAt:        &oneHourAgo,
	}
}

// newReaperFixture builds a reaper with mocked repo + publisher and exposes
// a Batcher mock the test can use to express the *expected* resend path.
type reaperFixture struct {
	repo      *RepositoryMock
	publisher *ResultPublisherMock
	sender    *BatcherMock
}

func newReaperFixture() *reaperFixture {
	return &reaperFixture{
		repo:      &RepositoryMock{},
		publisher: &ResultPublisherMock{},
		sender:    &BatcherMock{},
	}
}

// TestReaper_StuckRow_BeforeDeadLetter_AttemptsResend documents Issue #2 from
// When the reaper claims a stuck row, it MUST attempt
// to re-broadcast the original signed tx before marking the row failed.
//
// CONSEQUENCES OF THE BUG (current behaviour):
//   - A single transient mempool drop is unrecoverable.
//   - One dropped tx out of 60 → the e2e poll for count==60 times out.
//   - For a token bridge: the source side burns, the destination never mints,
//     tokens are permanently destroyed.
//   - Under axyl's ~30/30min consensus-vote-cancellation rate this WILL
//     reoccur in every long-running e2e run.
func TestReaper_StuckRow_BeforeDeadLetter_AttemptsResend(t *testing.T) {
	f := newReaperFixture()
	row := stuckRow("test-correlation-id-RESEND")

	// Repo returns one stuck row when ClaimStuck is called.
	f.repo.ClaimStuckFunc = func(ctx context.Context, identity string, olderThan time.Time, limit int) ([]batcher.CTSTransaction, error) {
		return []batcher.CTSTransaction{row}, nil
	}
	// The mocks below are wired to record calls so we can audit the
	// reaper's behaviour against the EXPECTED resend semantics.
	f.repo.IncrementSendAttemptsFunc = func(ctx context.Context, keys []batcher.TxKey) error { return nil }
	f.repo.MarkResentFunc = func(ctx context.Context, key batcher.TxKey) error { return nil }
	f.sender.SendFunc = func(ctx context.Context, inputs []ethrpc.TransactionInput) ([]ethrpc.SendResult, error) {
		return []ethrpc.SendResult{{Hash: row.TxHash}}, nil
	}
	f.publisher.PushFunc = func(ctx context.Context, r types.TxResult) error { return nil }
	f.repo.MarkFailedFunc = func(ctx context.Context, key batcher.TxKey, reason string) error { return nil }

	// Build the reaper as it exists today (no sender wired in — that is the
	// bug). When the fix lands, the constructor will accept a sender and
	// MaxResendAttempts, and these tests will be updated to wire them in.
	r := batcher.NewReaperService(
		batcher.ReaperConfig{
			Identity:          "privatenode",
			BatchSize:         100,
			Interval:          time.Minute,
			StuckThreshold:    time.Minute,
			MaxResendAttempts: 3,
		},
		f.repo, f.sender, f.publisher,
	)

	// Trigger a single reap cycle by running the service with a context
	// that cancels almost immediately.
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_ = r.Run(ctx) // best-effort: Run returns nil on ctx.Done.

	// Pre-conditions: the reaper SAW our stuck row.
	assert.GreaterOrEqualf(t, len(f.repo.ClaimStuckCalls()), 1,
		"reaper should have called ClaimStuck at least once")

	// ───── REGRESSION ASSERTION (Issue #2) ──────────────────────────────
	//
	// The reaper SHOULD attempt to resend the stuck row before giving up.
	// A resend is observable via either:
	//   (a) a call to Batcher.Send with the original calldata, OR
	//   (b) IncrementSendAttempts being called from the reaper path.
	//
	// Neither happens today.
	sendCalls := f.sender.SendCalls()
	incCalls := f.repo.IncrementSendAttemptsCalls()
	failedCalls := f.repo.MarkFailedCalls()

	t.Logf("OBSERVED:")
	t.Logf("  Batcher.Send calls (resend attempts):           %d", len(sendCalls))
	t.Logf("  Repository.IncrementSendAttempts calls (reaper): %d", len(incCalls))
	t.Logf("  Repository.MarkFailed calls (dead-letter):      %d", len(failedCalls))

	t.Logf("EXPECTED:")
	t.Logf("  Batcher.Send calls (resend attempts):           >= 1")
	t.Logf("  Repository.IncrementSendAttempts calls:         >= 1")
	t.Logf("  Repository.MarkFailed calls:                    0  (first attempt should resend, not dead-letter)")

	t.Logf("CONSEQUENCE OF BUG: 1 tx out of 60 dropped by axyl => message permanently lost.")
	t.Logf("  In a token-bridge context, the source PL burned the tokens but the dest")
	t.Logf("  PL never minted them. Funds are destroyed and no operator-triggered retry")
	t.Logf("  path exists (relayer DB row state=10/outcome=reverted is terminal).")

	assert.GreaterOrEqualf(t, len(sendCalls), 1,
		"REGRESSION: reaper must attempt to re-broadcast the stuck tx before dead-lettering. "+
			"Currently it goes straight to MarkFailed with zero resend attempts. "+
			"Fix: in reaper.reap(), call sender.Send for rows with send_attempts < MaxResendAttempts "+
			"before publishing TxResultFailed.")

	assert.GreaterOrEqualf(t, len(incCalls), 1,
		"REGRESSION: reaper must IncrementSendAttempts on resend so the resend cap is enforced. "+
			"Currently the reaper never touches send_attempts (only the sender does, on first broadcast). "+
			"See the proposed code shape.")

	assert.Equalf(t, 0, len(failedCalls),
		"REGRESSION: reaper should NOT dead-letter on the first claim of a stuck row. "+
			"Dead-letter must come ONLY after MaxResendAttempts re-broadcasts have all failed.")
}

// TestReaper_ResendAttemptsRespectMaxBound documents the *bound* on resends.
// A correct fix won't loop forever — once `send_attempts >= MaxResendAttempts`
// the reaper finalises the row as failed and lets the relayer compensate
// (or, in 9b's design, retry via a fresh dispatch).
//
// This test wires a row that has already been resent MaxResendAttempts times
// and asserts the reaper proceeds to MarkFailed without another Send. Even
// against a fix that adds resend, this test guards the cap.
func TestReaper_ResendAttemptsRespectMaxBound(t *testing.T) {
	const maxResendAttempts = 3

	f := newReaperFixture()
	row := stuckRow("test-correlation-id-CAP")
	row.SendAttempts = maxResendAttempts // exhausted

	f.repo.ClaimStuckFunc = func(ctx context.Context, identity string, olderThan time.Time, limit int) ([]batcher.CTSTransaction, error) {
		return []batcher.CTSTransaction{row}, nil
	}
	f.repo.MarkFailedFunc = func(ctx context.Context, key batcher.TxKey, reason string) error { return nil }
	f.publisher.PushFunc = func(ctx context.Context, r types.TxResult) error { return nil }
	f.sender.SendFunc = func(ctx context.Context, inputs []ethrpc.TransactionInput) ([]ethrpc.SendResult, error) {
		return []ethrpc.SendResult{{Hash: row.TxHash}}, nil
	}

	r := batcher.NewReaperService(
		batcher.ReaperConfig{
			Identity:          "privatenode",
			BatchSize:         100,
			Interval:          time.Minute,
			StuckThreshold:    time.Minute,
			MaxResendAttempts: 3,
		},
		f.repo, f.sender, f.publisher,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_ = r.Run(ctx)

	// After cap reached, MarkFailed is correct behaviour.
	t.Logf("OBSERVED MarkFailed=%d, Send=%d (row.SendAttempts=%d, cap=%d)",
		len(f.repo.MarkFailedCalls()), len(f.sender.SendCalls()),
		row.SendAttempts, maxResendAttempts)

	assert.Equalf(t, 1, len(f.repo.MarkFailedCalls()),
		"once send_attempts has reached MaxResendAttempts, the reaper MUST dead-letter "+
			"(this is the well-behaved terminal branch).")

	assert.Equalf(t, 0, len(f.sender.SendCalls()),
		"the reaper MUST NOT resend a row whose send_attempts has already reached the cap "+
			"(prevents infinite resend loops when the chain is genuinely down).")
}

// TestReaper_DeadLetterDoesNotMarkFailed_AsFinishedRevert pins down a subtle
// contract that the relayer's downstream (HandleFailedMined) depends on:
// dead-lettered rows have `status=failed` (terminal) with `error_reason`
// populated by `stuckReason`. They MUST NOT be misclassified as
// MarkFinishedRevert (which the relayer treats as an actual on-chain revert
// and would in turn collapse). This test exists to catch a careless
// "unify status=failed and status=finished-revert" refactor.
func TestReaper_DeadLetterDoesNotMarkFailed_AsFinishedRevert(t *testing.T) {
	f := newReaperFixture()
	row := stuckRow("test-correlation-id-NOT-REVERT")
	// Drive the reaper straight into the dead-letter branch — this test
	// is about the *dead-letter* shape, not the resend path, so we
	// disable resend with MaxResendAttempts=0.

	f.repo.ClaimStuckFunc = func(ctx context.Context, identity string, olderThan time.Time, limit int) ([]batcher.CTSTransaction, error) {
		return []batcher.CTSTransaction{row}, nil
	}
	f.repo.MarkFailedFunc = func(ctx context.Context, key batcher.TxKey, reason string) error { return nil }
	f.repo.MarkFinishedRevertFunc = func(ctx context.Context, key batcher.TxKey, revertData []byte) error { return nil }
	f.publisher.PushFunc = func(ctx context.Context, r types.TxResult) error { return nil }

	r := batcher.NewReaperService(
		batcher.ReaperConfig{
			Identity:          "privatenode",
			BatchSize:         100,
			Interval:          time.Minute,
			StuckThreshold:    time.Minute,
			MaxResendAttempts: 0,
		},
		f.repo, f.sender, f.publisher,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_ = r.Run(ctx)

	assert.Equalf(t, 0, len(f.repo.MarkFinishedRevertCalls()),
		"dead-lettered (stuck) rows MUST NOT be persisted as MarkFinishedRevert — they were "+
			"never mined. Misclassifying them as Revert would make the failure look like an "+
			"on-chain require() and confuse the relayer's compensation logic.")
}
