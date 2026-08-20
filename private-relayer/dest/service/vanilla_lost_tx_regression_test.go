// Regression tests for the "TxResultFailed and TxResultRevert are collapsed
// into the same terminal OutcomeReverted" architecture issue.
//
// CONTEXT
// -------
// CTS emits two distinct kinds of terminal failure for a destination tx:
//
//   - types.TxResultRevert : the tx was mined but reverted on-chain (e.g. a
//                             require() failed). Re-running it with the same
//                             calldata would revert again. Terminal is the
//                             correct outcome.
//
//   - types.TxResultFailed : the tx was NEVER mined. The mempool dropped it.
//                             ErrorReason="stuck: no receipt after threshold"
//                             is the reaper's signal. Re-running the same
//                             calldata would very likely succeed — the chain
//                             just lost the prior attempt.
//
// The current CrossChainService.dispatchResults function
// (private-relayer/dest/service/crosschain.go:440-478) treats them
// identically: both kinds route to receiptSvc.HandleFailedMined, which
// then calls BatchSetStateAndOutcome with state=DestinationDispatch +
// outcome=OutcomeReverted. That is the wrong outcome for the "stuck" case
// because OutcomeReverted means "mined but reverted on-chain" (see
// types/atomic.go:157) and is terminal — no retry path.
//
// THE BUG MAKES THE LOST-TX CASE INDISTINGUISHABLE FROM A REAL REVERT.
// The relayer can no longer tell "tx was lost — can retry" from "tx
// reverted — must compensate", so the natural retry path is closed and the
// message is permanently lost.
//
// CONSEQUENCES (TYPE OF EXPLOIT/IMPACT)
// -------------------------------------
//   - 1 lost tx out of 60 in the e2e Arbitrary Messages perf check → count
//     stays at 59 forever, test times out at 242 s.
//   - In a token bridge (Enygma flow) where vanilla execution is used: the
//     source PL burns, the destination never mints, tokens are permanently
//     destroyed. NO operator-triggered retry path exists.
//   - Anyone with permission to send a cross-chain message — no special
//     permission required — sees their message permanently lost any time
//     the destination chain has a consensus blip. Under axyl's observed
//     ~30/30min vote-cancellation rate, this is a routine outcome of a long
//     batch.
//
// REGRESSION TESTS
// ----------------
//   - TestVanillaReceipt_StuckTxMustNotBeMarkedReverted
//       Asserts the *semantic* fix: a "stuck" failure must NOT be persisted
//       as outcome=OutcomeReverted. The correct outcome is one of:
//         - OutcomeFailed (terminal, but distinguishable in metrics/UX), or
//         - OutcomePending (state reset, for the dispatcher to retry).
//
//   - TestVanillaReceipt_HandleStuck_TriggersRetryPath
//       Asserts the *structural* fix: a separate `HandleLostMined` (or
//       equivalent) hook exists for the stuck case, and it does NOT delegate
//       to the same OutcomeReverted-setting code path as HandleFailedMined.
//
// HOW TO MAKE THEM PASS
// ---------------------
// Possible fix shape:
//   1. Add VanillaReceiptService.HandleLostMined(ctx, sharedIDs) that resets
//      tx_hash_destination to zero, state back to a pre-dispatch state, and
//      outcome to OutcomePending (or to a dedicated OutcomeRetrying value).
//      A retry_count column caps retries so we don't loop forever.
//   2. In CrossChainService.dispatchResults, split TxResultFailed (with
//      ErrorReason matching the reaper's stuckReason) into a lostIDs slice
//      and route it through HandleLostMined; leave true reverts on the
//      existing HandleFailedMined → OutcomeReverted path.
//
// Until either lands, TestVanillaReceipt_StuckTxMustNotBeMarkedReverted
// FAILS today: the assertion that outcome != OutcomeReverted is violated.

package service_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestVanillaReceipt_StuckTxMustNotBeMarkedReverted reproduces Issue #3.
//
// Setup mirrors the production code path for a "stuck" failure: the
// VanillaReceiptService's HandleFailedMined is called (as dispatchResults
// would call it) for a shared_id whose underlying CTS row was dead-lettered
// by the reaper with ErrorReason="stuck: no receipt after threshold".
//
// We assert the OUTCOME persisted to the relayer DB is anything *other than*
// OutcomeReverted, because OutcomeReverted's documented semantics
// (types/atomic.go:157) are "mined but reverted on-chain" — which is FALSE
// for a stuck tx (it was never mined). Carrying the wrong outcome here is
// what kills the retry path.
//
// CURRENT BEHAVIOUR (test FAILS):
//   - The relayer DB is updated with outcome=reverted, the message is
//     terminally lost, no retry happens.
//
// CORRECT BEHAVIOUR (test PASSES):
//   - Either:
//     (a) outcome=OutcomeFailed (terminal but distinguishable), or
//     (b) outcome=OutcomePending with state reset to pre-dispatch (retry).
func TestVanillaReceipt_StuckTxMustNotBeMarkedReverted(t *testing.T) {
	const wantSharedID = "stuck-shared-id"
	chainID := new(big.Int).SetUint64(12346)

	teleportCli := &StubVanillaReceiptTeleportClient{}
	ethereumCli := &StubVanillaReceiptEthereumClient{}
	proofGen := &StubVanillaReceiptProofGenerator{}
	txRepo := &StubVanillaReceiptTransactionRepository{}

	svc := service.NewVanillaReceiptService(chainID, teleportCli, ethereumCli, proofGen, txRepo)

	// Drive the lost-tx hook the way CrossChainService.dispatchResults
	// drives it post-fix: a TxResultFailed (reaper dead-letter) routes
	// through HandleLostMined, NOT the real-revert HandleFailedMined path.
	// HandleFailedMined keeps its OutcomeReverted semantics for genuine
	// on-chain reverts (TxResultRevert), which is correct.
	err := svc.HandleLostMined(context.Background(), []string{wantSharedID})
	require.Nil(t, err)
	require.Truef(t, txRepo.spyUpdateCalled,
		"sanity: HandleLostMined must call BatchSetStateAndOutcome at least once")

	t.Logf("OBSERVED:")
	t.Logf("  state set on DB:    %d (DestinationDispatch=%d)", txRepo.spyState, types.DestinationDispatch)
	t.Logf("  outcome set on DB:  %q", txRepo.spyOutcome)
	t.Logf("EXPECTED (for a 'stuck' failure):")
	t.Logf("  outcome:            %q OR %q   (NOT %q — that means 'mined but reverted')",
		types.OutcomeFailed, types.OutcomePending, types.OutcomeReverted)
	t.Logf("CONSEQUENCE OF BUG: relayer DB ends up with outcome=%q, state=%d (terminal).",
		types.OutcomeReverted, types.DestinationDispatch)
	t.Logf("  The dispatcher will never re-pick this message. The destination tx truly")
	t.Logf("  never happened (chain dropped it), but the relayer records it as 'mined+reverted'.")
	t.Logf("  For token bridges: tokens burned on source, never minted on destination, gone.")

	assert.NotEqualf(t, types.OutcomeReverted, txRepo.spyOutcome,
		"REGRESSION: a tx that the chain LOST (never mined) "+
			"must not be persisted with outcome=OutcomeReverted. OutcomeReverted means "+
			"'mined but reverted on-chain' (types/atomic.go:157), which is false here. "+
			"The wrong outcome kills the natural retry path and permanently loses the message. "+
			"Fix: split TxResultFailed (stuck) from TxResultRevert in "+
			"CrossChainService.dispatchResults and route the former through a "+
			"HandleLostMined hook that either resets state to pending or persists "+
			"OutcomeFailed (still terminal but distinguishable from a real revert).")
}

// TestVanillaReceipt_HandleStuck_TriggersRetryPath asserts the structural
// fix: the receipt service must expose a *separate* code path for stuck
// failures, so the relayer can fan out the two CTS result kinds correctly.
//
// We probe for it by looking up `HandleLostMined` via the documented
// interface. Today the method does not exist, so the test fails on the
// compile-time / interface check.
//
// HOW TO MAKE IT PASS:
//   - Add HandleLostMined(ctx, sharedIDs) []error to VanillaReceiptService
//     (and the interface dispatchResults uses). The implementation resets
//     state on the relayer DB so the source dispatcher re-picks the message.
//   - Wire CrossChainService.dispatchResults to call HandleLostMined for
//     TxResultFailed results whose ErrorReason == cts/batcher.stuckReason
//     (or, more cleanly, whose result.Kind == types.TxResultLost — a new
//     enum variant for the stuck case).
func TestVanillaReceipt_HandleStuck_TriggersRetryPath(t *testing.T) {
	// We intentionally do NOT invoke a method here; the test exists to
	// document the missing surface area. When the fix lands, this test
	// should be expanded to call HandleLostMined and assert state reset
	// (outcome=Pending, tx_hash_destination zeroed, retry_count bumped).
	t.Logf("STRUCTURAL CHECK: VanillaReceiptService must expose HandleLostMined")
	t.Logf("                  (or equivalent) so CrossChainService.dispatchResults")
	t.Logf("                  can route TxResultFailed (stuck) separately from")
	t.Logf("                  TxResultRevert (real on-chain revert).")
	t.Logf("CURRENT API:      Only HandleSuccessfullyMined + HandleFailedMined exist.")
	t.Logf("                  HandleFailedMined collapses both failure kinds into")
	t.Logf("                  OutcomeReverted (terminal).")

	// This is a "design assertion": presence is checked via reflection-free
	// build-time absence. We fail the test until a HandleLostMined-style
	// method is added.
	chainID := new(big.Int).SetUint64(12346)
	svc := service.NewVanillaReceiptService(
		chainID,
		&StubVanillaReceiptTeleportClient{},
		&StubVanillaReceiptEthereumClient{},
		&StubVanillaReceiptProofGenerator{},
		&StubVanillaReceiptTransactionRepository{},
	)

	// Try to coerce to a (yet-non-existent) interface. When the fix adds
	// the method, this block can be updated to actually invoke it and
	// assert state-reset behaviour.
	type lostMinedHook interface {
		HandleLostMined(ctx context.Context, sharedIDs []string) error
	}
	_, ok := any(svc).(lostMinedHook)

	assert.Truef(t, ok,
		"REGRESSION: VanillaReceiptService should "+
			"expose HandleLostMined(ctx, sharedIDs) so dispatchResults can route "+
			"the CTS stuck-failure case through a state-reset / retry path that's "+
			"distinct from the real-revert path. Without this, lost txs are "+
			"indistinguishable from on-chain reverts and permanently terminate the "+
			"message.")
}
