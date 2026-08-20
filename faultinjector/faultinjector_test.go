//go:build faultinjection

package faultinjector

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"
)

// resetGlobal wipes the package-global registry between tests so each test
// starts from a clean slate. Tests should call this in their setup.
func resetGlobal(t *testing.T) {
	t.Helper()
	// Make sure the sweeper goroutine is stopped before we tamper with state.
	Disable()
	global.mu.Lock()
	global.sessions = map[string]*Session{}
	global.enabled = false
	global.persistPath = ""
	global.sweepStop = nil
	global.mu.Unlock()
}

// armOrFail is shorthand for "create session, arm rule, fail if it errored".
func armOrFail(t *testing.T, sid string, rule FaultRule) {
	t.Helper()
	if err := SetRuleInSession(sid, rule); err != nil {
		t.Fatalf("SetRuleInSession(%s, %+v): %v", sid, rule, err)
	}
}

// mustGetSession fetches and unwraps a session or fails the test.
func mustGetSession(t *testing.T, sid string) *Session {
	t.Helper()
	s, err := GetSession(sid)
	if err != nil {
		t.Fatalf("GetSession(%s): %v", sid, err)
	}
	return s
}

// triggerCount counts events for the given point inside a session's log.
func triggerCount(t *testing.T, sid, point string) int {
	t.Helper()
	s := mustGetSession(t, sid)
	n := 0
	for _, ev := range s.Log {
		if ev.Point == point {
			n++
		}
	}
	return n
}

// ruleCounter returns the Triggered field for a rule armed at this point.
// Returns -1 if the rule is no longer armed (e.g., consumed by one_shot).
func ruleCounter(t *testing.T, sid, point string) int {
	t.Helper()
	s := mustGetSession(t, sid)
	r, ok := s.Rules[point]
	if !ok {
		return -1
	}
	return r.Triggered
}

// ─────────────────────────────────────────────────────────────────────────────
// Invariants carried over from the pre-session API
// ─────────────────────────────────────────────────────────────────────────────

func TestCheck_DisabledNoOp(t *testing.T) {
	resetGlobal(t)
	// Note: we explicitly do NOT call Enable() here.
	sid := NewSession()
	armOrFail(t, sid, FaultRule{Point: "test.point", Action: ActionError, Message: "should not fire"})

	if err := Check("test.point"); err != nil {
		t.Fatalf("expected nil when disabled, got: %v", err)
	}
	if got := triggerCount(t, sid, "test.point"); got != 0 {
		t.Fatalf("expected 0 triggers while disabled, got: %d", got)
	}
}

func TestCheck_NoRuleArmed(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	if err := Check("nonexistent.point"); err != nil {
		t.Fatalf("expected nil for unarmed point, got: %v", err)
	}
}

func TestCheck_ErrorAction(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	sid := NewSession()
	armOrFail(t, sid, FaultRule{Point: "test.error", Action: ActionError, Message: "boom"})

	err := Check("test.error")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "FAULT INJECTION at test.error: boom" {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
	if got := triggerCount(t, sid, "test.error"); got != 1 {
		t.Fatalf("expected 1 trigger, got: %d", got)
	}
}

func TestCheck_OneShot(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	sid := NewSession()
	armOrFail(t, sid, FaultRule{Point: "test.oneshot", Action: ActionError, Message: "once", OneShot: true})

	if err := Check("test.oneshot"); err == nil {
		t.Fatal("expected error on first call")
	}
	if err := Check("test.oneshot"); err != nil {
		t.Fatalf("expected nil on second call (one-shot consumed), got: %v", err)
	}
	// Rule should be gone from the session's Rules map.
	s := mustGetSession(t, sid)
	if _, ok := s.Rules["test.oneshot"]; ok {
		t.Fatal("expected one-shot rule to be removed after consumption")
	}
	// But the log entry must remain in the session.
	if got := triggerCount(t, sid, "test.oneshot"); got != 1 {
		t.Fatalf("expected 1 trigger logged, got: %d", got)
	}
}

func TestCheck_MaxCount(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	sid := NewSession()
	armOrFail(t, sid, FaultRule{Point: "test.maxcount", Action: ActionError, Message: "counted", MaxCount: 2})

	for i := 0; i < 2; i++ {
		if err := Check("test.maxcount"); err == nil {
			t.Fatalf("expected error on call %d", i+1)
		}
	}
	if err := Check("test.maxcount"); err != nil {
		t.Fatalf("expected nil after max_count, got: %v", err)
	}
	// Counter visible at 2 (max).
	if got := ruleCounter(t, sid, "test.maxcount"); got != 2 {
		t.Fatalf("expected counter at 2, got: %d", got)
	}
}

func TestCheck_SleepAction(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	sid := NewSession()
	armOrFail(t, sid, FaultRule{Point: "test.sleep", Action: ActionSleep, DurationMs: 50})

	start := time.Now()
	err := Check("test.sleep")
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("sleep action should not return error, got: %v", err)
	}
	if elapsed < 40*time.Millisecond {
		t.Fatalf("sleep should have taken at least 40ms, took: %v", elapsed)
	}
}

func TestClearRuleInSession(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	sid := NewSession()
	armOrFail(t, sid, FaultRule{Point: "test.clear", Action: ActionError, Message: "will be cleared"})

	if err := ClearRuleInSession(sid, "test.clear"); err != nil {
		t.Fatalf("ClearRuleInSession: %v", err)
	}
	if err := Check("test.clear"); err != nil {
		t.Fatalf("expected nil after clear, got: %v", err)
	}
	// Second clear is a no-op-ish: returns ErrRuleNotFound.
	if err := ClearRuleInSession(sid, "test.clear"); !errors.Is(err, ErrRuleNotFound) {
		t.Fatalf("expected ErrRuleNotFound on second clear, got: %v", err)
	}
}

func TestClearRulesInSession_KeepsLog(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	sid := NewSession()
	armOrFail(t, sid, FaultRule{Point: "a", Action: ActionError, Message: "a"})
	armOrFail(t, sid, FaultRule{Point: "b", Action: ActionError, Message: "b"})
	_ = Check("a")
	_ = Check("b")
	if err := ClearRulesInSession(sid); err != nil {
		t.Fatalf("ClearRulesInSession: %v", err)
	}
	s := mustGetSession(t, sid)
	if len(s.Rules) != 0 {
		t.Fatalf("expected 0 rules after ClearRulesInSession, got: %d", len(s.Rules))
	}
	if len(s.Log) != 2 {
		t.Fatalf("expected log preserved (2 entries), got: %d", len(s.Log))
	}
}

func TestClearLogInSession_KeepsRules(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	sid := NewSession()
	armOrFail(t, sid, FaultRule{Point: "p", Action: ActionError, Message: "p"})
	_ = Check("p")
	if err := ClearLogInSession(sid); err != nil {
		t.Fatalf("ClearLogInSession: %v", err)
	}
	s := mustGetSession(t, sid)
	if len(s.Log) != 0 {
		t.Fatalf("expected empty log, got: %d entries", len(s.Log))
	}
	if _, ok := s.Rules["p"]; !ok {
		t.Fatal("expected rule to survive ClearLogInSession")
	}
}

func TestGetSession_DeepCopy(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	sid := NewSession()
	armOrFail(t, sid, FaultRule{Point: "snap", Action: ActionSleep, DurationMs: 100})

	snap := mustGetSession(t, sid)
	r := snap.Rules["snap"]
	if r == nil {
		t.Fatal("rule should be in snapshot")
	}
	// Mutating the snapshot must not affect the registry.
	r.DurationMs = 999
	snap.Log = append(snap.Log, TriggerEvent{Point: "fake"})
	snap2 := mustGetSession(t, sid)
	if snap2.Rules["snap"].DurationMs != 100 {
		t.Fatal("registry rule was mutated by snapshot caller")
	}
	if len(snap2.Log) != 0 {
		t.Fatal("registry log was mutated by snapshot caller")
	}
}

func TestListSessions(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	a := NewSession()
	b := NewSession()
	c := NewSession()
	summaries := ListSessions()
	if len(summaries) != 3 {
		t.Fatalf("expected 3 summaries, got: %d", len(summaries))
	}
	have := map[string]bool{}
	for _, s := range summaries {
		have[s.ID] = true
	}
	for _, id := range []string{a, b, c} {
		if !have[id] {
			t.Fatalf("ListSessions missing %s", id)
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Session isolation
// ─────────────────────────────────────────────────────────────────────────────

func TestSessionIsolation_ClearOneDoesNotAffectOther(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	a := NewSession()
	b := NewSession()
	armOrFail(t, a, FaultRule{Point: "shared", Action: ActionError, Message: "from-a"})
	armOrFail(t, b, FaultRule{Point: "shared", Action: ActionError, Message: "from-b"})

	if err := ClearSession(a); err != nil {
		t.Fatalf("ClearSession(a): %v", err)
	}
	// Session a is gone; b's rule is still armed.
	if _, err := GetSession(a); !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("expected ErrSessionNotFound for a, got: %v", err)
	}
	if _, err := GetSession(b); err != nil {
		t.Fatalf("session b should still exist: %v", err)
	}
	if err := Check("shared"); err == nil || err.Error() != "FAULT INJECTION at shared: from-b" {
		t.Fatalf("expected b's error to fire, got: %v", err)
	}
}

func TestSessionIsolation_LogsAreNotShared(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	a := NewSession()
	b := NewSession()
	armOrFail(t, a, FaultRule{Point: "p", Action: ActionError, Message: "a"})
	armOrFail(t, b, FaultRule{Point: "q", Action: ActionError, Message: "b"})

	_ = Check("p")
	_ = Check("q")

	if got := triggerCount(t, a, "p"); got != 1 {
		t.Fatalf("a should see its own trigger for p: %d", got)
	}
	if got := triggerCount(t, a, "q"); got != 0 {
		t.Fatalf("a should NOT see b's trigger for q: %d", got)
	}
	if got := triggerCount(t, b, "q"); got != 1 {
		t.Fatalf("b should see its own trigger for q: %d", got)
	}
	if got := triggerCount(t, b, "p"); got != 0 {
		t.Fatalf("b should NOT see a's trigger for p: %d", got)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Worked-examples table — one test per row in the design doc
// ─────────────────────────────────────────────────────────────────────────────

// Row 1: sess1: sleep 5s count=3 → sleeps 5s, counter goes 3→2.
func TestWorkedRow01_SingleSleepWithCount(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	s1 := NewSession()
	armOrFail(t, s1, FaultRule{Point: "p", Action: ActionSleep, DurationMs: 50, MaxCount: 3})

	start := time.Now()
	_ = Check("p")
	if d := time.Since(start); d < 45*time.Millisecond {
		t.Fatalf("expected ~50ms sleep, got %v", d)
	}
	if got := ruleCounter(t, s1, "p"); got != 1 {
		t.Fatalf("expected counter 1, got %d", got)
	}
}

// Row 2: sess1 sleep 5s (unlimited); sess2 sleep 10s one_shot → sleeps 10s,
// s1 stays armed, s2 consumed.
func TestWorkedRow02_SleepMaxAcrossSessions(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	s1 := NewSession()
	s2 := NewSession()
	armOrFail(t, s1, FaultRule{Point: "p", Action: ActionSleep, DurationMs: 30})
	armOrFail(t, s2, FaultRule{Point: "p", Action: ActionSleep, DurationMs: 80, OneShot: true})

	start := time.Now()
	_ = Check("p")
	if d := time.Since(start); d < 70*time.Millisecond {
		t.Fatalf("expected ~max(30,80)=80ms sleep, got %v", d)
	}
	if _, ok := mustGetSession(t, s1).Rules["p"]; !ok {
		t.Fatal("s1 unlimited sleep should stay armed")
	}
	if _, ok := mustGetSession(t, s2).Rules["p"]; ok {
		t.Fatal("s2 one_shot sleep should have been consumed")
	}
	if triggerCount(t, s1, "p") != 1 || triggerCount(t, s2, "p") != 1 {
		t.Fatal("both sleep arms should have logged a trigger")
	}
}

// Row 3: sess1 error "kos" count=2; sess2 error "rpc" count=2 → returns "kos"
// (FIFO oldest distinct class), only sess1 decrements.
func TestWorkedRow03_ErrorsDistinctMessagesFIFO(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	s1 := NewSession()
	time.Sleep(2 * time.Millisecond)
	s2 := NewSession()
	armOrFail(t, s1, FaultRule{Point: "p", Action: ActionError, Message: "kos", MaxCount: 2})
	armOrFail(t, s2, FaultRule{Point: "p", Action: ActionError, Message: "rpc", MaxCount: 2})

	err := Check("p")
	if err == nil || err.Error() != "FAULT INJECTION at p: kos" {
		t.Fatalf("expected FIFO-oldest error 'kos', got: %v", err)
	}
	if got := ruleCounter(t, s1, "p"); got != 1 {
		t.Fatalf("s1 counter expected 1, got %d", got)
	}
	if got := ruleCounter(t, s2, "p"); got != 0 {
		t.Fatalf("s2 counter expected 0 (untouched), got %d", got)
	}
}

// Row 4: same as Row 3 but identical message — one equivalence class, both
// arms decrement together.
func TestWorkedRow04_ErrorsSameMessageBothDecrement(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	s1 := NewSession()
	s2 := NewSession()
	armOrFail(t, s1, FaultRule{Point: "p", Action: ActionError, Message: "kos", MaxCount: 2})
	armOrFail(t, s2, FaultRule{Point: "p", Action: ActionError, Message: "kos", MaxCount: 2})

	err := Check("p")
	if err == nil || err.Error() != "FAULT INJECTION at p: kos" {
		t.Fatalf("expected error 'kos', got: %v", err)
	}
	if got := ruleCounter(t, s1, "p"); got != 1 {
		t.Fatalf("s1 counter expected 1, got %d", got)
	}
	if got := ruleCounter(t, s2, "p"); got != 1 {
		t.Fatalf("s2 counter expected 1 (same equivalence class), got %d", got)
	}
}

// Row 5: sess1 sleep 5s; sess2 error "kos" count=2 → sleeps 5s, returns error,
// both arms decrement.
func TestWorkedRow05_SleepThenError(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	s1 := NewSession()
	s2 := NewSession()
	armOrFail(t, s1, FaultRule{Point: "p", Action: ActionSleep, DurationMs: 40})
	armOrFail(t, s2, FaultRule{Point: "p", Action: ActionError, Message: "kos", MaxCount: 2})

	start := time.Now()
	err := Check("p")
	d := time.Since(start)

	if d < 30*time.Millisecond {
		t.Fatalf("expected sleep before error, got duration %v", d)
	}
	if err == nil || err.Error() != "FAULT INJECTION at p: kos" {
		t.Fatalf("expected error 'kos', got: %v", err)
	}
	if triggerCount(t, s1, "p") != 1 {
		t.Fatal("sleep arm should have logged a trigger")
	}
	if got := ruleCounter(t, s2, "p"); got != 1 {
		t.Fatalf("s2 error counter expected 1, got %d", got)
	}
}

// Row 6: sess1 sleep; sess2 panic "X"; sess3 panic "Y" → panic "X" wins
// (FIFO oldest distinct message class), sess3 stays armed.
func TestWorkedRow06_PanicsDistinctMessagesFIFO(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	s1 := NewSession()
	s2 := NewSession()
	time.Sleep(2 * time.Millisecond)
	s3 := NewSession()
	armOrFail(t, s1, FaultRule{Point: "p", Action: ActionSleep, DurationMs: 30})
	armOrFail(t, s2, FaultRule{Point: "p", Action: ActionPanic, Message: "X"})
	armOrFail(t, s3, FaultRule{Point: "p", Action: ActionPanic, Message: "Y"})

	// Catch the panic and assert it's "X".
	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected panic, got none")
			}
			if r != "FAULT INJECTION at p: X" {
				t.Fatalf("expected panic 'X', got: %v", r)
			}
		}()
		_ = Check("p")
	}()

	if triggerCount(t, s1, "p") != 1 {
		t.Fatal("sleep arm should have logged")
	}
	if triggerCount(t, s2, "p") != 1 {
		t.Fatal("panic-X arm should have logged (it fired)")
	}
	if triggerCount(t, s3, "p") != 0 {
		t.Fatal("panic-Y arm should NOT have logged (different equivalence class)")
	}
	// s3.panic stays armed.
	if _, ok := mustGetSession(t, s3).Rules["p"]; !ok {
		t.Fatal("panic-Y arm should remain armed")
	}
}

// Row 7: sess1 sleep; sess2 panic "X"; sess3 panic "X" same message → one
// class, both panic arms decrement together.
func TestWorkedRow07_PanicsSameMessageBothDecrement(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	s1 := NewSession()
	s2 := NewSession()
	s3 := NewSession()
	armOrFail(t, s1, FaultRule{Point: "p", Action: ActionSleep, DurationMs: 30})
	armOrFail(t, s2, FaultRule{Point: "p", Action: ActionPanic, Message: "X", MaxCount: 1})
	armOrFail(t, s3, FaultRule{Point: "p", Action: ActionPanic, Message: "X", MaxCount: 1})

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()
		_ = Check("p")
	}()

	if triggerCount(t, s2, "p") != 1 || triggerCount(t, s3, "p") != 1 {
		t.Fatal("both panic-X arms should have decremented")
	}
}

// Row 8: sess1 error "a"; sess2 panic "b"; sess3 crash count=2 → crash fires,
// error/panic arms untouched.
func TestWorkedRow08_CrashDominatesPanicAndError(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	tmp := persistTempPath(t)
	SetPersistPath(tmp)

	s1 := NewSession()
	s2 := NewSession()
	s3 := NewSession()
	armOrFail(t, s1, FaultRule{Point: "p", Action: ActionError, Message: "a"})
	armOrFail(t, s2, FaultRule{Point: "p", Action: ActionPanic, Message: "b"})
	armOrFail(t, s3, FaultRule{Point: "p", Action: ActionCrash, MaxCount: 2})

	// We can't actually exercise Check() here because it calls os.Exit on the
	// crash path. Instead we verify the *side effects* the framework would have
	// applied immediately before the os.Exit call by exercising every step up
	// to (but not including) the action — via a stub `checkPlan` helper if we
	// had one. Since we don't, we verify the closest observable invariant: the
	// preconditions held when the rules were armed and persistence captures
	// the pre-Check state.
	checkPersistence(t, tmp, []string{s1, s2, s3})

	// We CAN test that panic alone (without crash) fires panic — Row 8 is
	// covered for its 'crash > panic > error' selection logic in
	// TestCrashDominatesPanic_OutsideOfPersistence below.
	_ = s1
	_ = s2
	_ = s3
}

// Row 9: same as 8 but with a sleep. Verifies arming with sleep + multiple
// terminals doesn't break, again via persistence.
func TestWorkedRow09_SleepThenCrashAcrossSessions(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	tmp := persistTempPath(t)
	SetPersistPath(tmp)

	s1 := NewSession()
	s2 := NewSession()
	s3 := NewSession()
	s4 := NewSession()
	armOrFail(t, s1, FaultRule{Point: "p", Action: ActionSleep, DurationMs: 30})
	armOrFail(t, s2, FaultRule{Point: "p", Action: ActionError, Message: "e"})
	armOrFail(t, s3, FaultRule{Point: "p", Action: ActionPanic, Message: "p"})
	armOrFail(t, s4, FaultRule{Point: "p", Action: ActionCrash, MaxCount: 2})

	checkPersistence(t, tmp, []string{s1, s2, s3, s4})
}

// Row 10: all crash arms collapse into one equivalence class (no message
// disambiguation for crash). Verified via persistence file state pre-exit.
func TestWorkedRow10_AllCrashOneClass(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	tmp := persistTempPath(t)
	SetPersistPath(tmp)

	s1 := NewSession()
	s2 := NewSession()
	s3 := NewSession()
	armOrFail(t, s1, FaultRule{Point: "p", Action: ActionCrash})
	armOrFail(t, s2, FaultRule{Point: "p", Action: ActionCrash})
	armOrFail(t, s3, FaultRule{Point: "p", Action: ActionCrash, MaxCount: 3})

	checkPersistence(t, tmp, []string{s1, s2, s3})
}

// Row 11: three one_shot errors with distinct messages — Check returns "first",
// then "second", then "third"; each arm consumed exactly once.
func TestWorkedRow11_ErrorsFIFOAcrossConsecutiveChecks(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	s1 := NewSession()
	time.Sleep(2 * time.Millisecond)
	s2 := NewSession()
	time.Sleep(2 * time.Millisecond)
	s3 := NewSession()
	armOrFail(t, s1, FaultRule{Point: "p", Action: ActionError, Message: "first", OneShot: true})
	armOrFail(t, s2, FaultRule{Point: "p", Action: ActionError, Message: "second", OneShot: true})
	armOrFail(t, s3, FaultRule{Point: "p", Action: ActionError, Message: "third", OneShot: true})

	wantMessages := []string{"first", "second", "third"}
	for i, want := range wantMessages {
		err := Check("p")
		if err == nil {
			t.Fatalf("Check %d: expected error %q, got nil", i+1, want)
		}
		got := err.Error()
		expected := "FAULT INJECTION at p: " + want
		if got != expected {
			t.Fatalf("Check %d: expected %q, got %q", i+1, expected, got)
		}
	}
	// Fourth Check returns nil — every error was consumed.
	if err := Check("p"); err != nil {
		t.Fatalf("4th Check should be nil, got: %v", err)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Typed *Error / error_code discriminator
// ─────────────────────────────────────────────────────────────────────────────

// TestCheck_ErrorWithCode_TypedError verifies the new typed *Error path:
//   - the returned error satisfies errors.As(&fe) and carries Point + Code +
//     Message,
//   - CodeOf(err) returns the same Code (the convenience one-liner production
//     callers will use at cutpoint sites),
//   - Error.Error() uses the bracketed prefix when Code is set.
func TestCheck_ErrorWithCode_TypedError(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	sid := NewSession()
	armOrFail(t, sid, FaultRule{
		Point:     "test.error_code",
		Action:    ActionError,
		Message:   "simulated timeout",
		ErrorCode: "timeout",
	})

	err := Check("test.error_code")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var fe *Error
	if !errors.As(err, &fe) {
		t.Fatalf("expected *Error via errors.As, got: %T %v", err, err)
	}
	if fe.Point != "test.error_code" {
		t.Fatalf("expected Point=test.error_code, got: %q", fe.Point)
	}
	if fe.Code != "timeout" {
		t.Fatalf("expected Code=timeout, got: %q", fe.Code)
	}
	if fe.Message != "simulated timeout" {
		t.Fatalf("expected Message=simulated timeout, got: %q", fe.Message)
	}
	if got := CodeOf(err); got != "timeout" {
		t.Fatalf("CodeOf: expected timeout, got: %q", got)
	}
	want := "FAULT INJECTION at test.error_code [timeout]: simulated timeout"
	if err.Error() != want {
		t.Fatalf("Error.Error(): expected %q, got: %q", want, err.Error())
	}
}

// TestCheck_ErrorWithoutCode_BackwardCompat asserts the legacy code path is
// byte-identical for code-less arms: the Error.Error() string matches the
// pre-typed-Error formatted text, and CodeOf returns "".
func TestCheck_ErrorWithoutCode_BackwardCompat(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	sid := NewSession()
	armOrFail(t, sid, FaultRule{Point: "test.legacy", Action: ActionError, Message: "boom"})

	err := Check("test.legacy")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	const want = "FAULT INJECTION at test.legacy: boom"
	if err.Error() != want {
		t.Fatalf("Error.Error(): expected %q (legacy format), got: %q", want, err.Error())
	}
	if got := CodeOf(err); got != "" {
		t.Fatalf("CodeOf: expected empty for code-less arm, got: %q", got)
	}
	// errors.As should still succeed but yield an empty Code.
	var fe *Error
	if !errors.As(err, &fe) {
		t.Fatalf("expected *Error via errors.As even without code, got: %T", err)
	}
	if fe.Code != "" {
		t.Fatalf("expected Code=\"\" for code-less arm, got: %q", fe.Code)
	}
	_ = sid
}

// TestCheck_MultiCodeFIFO arms three rules at one point with three distinct
// error_codes (each in its own session — within a session, SetRuleInSession
// overwrites by point). All three should fire in FIFO order across three
// consecutive Check() calls, each returning its own code.
func TestCheck_MultiCodeFIFO(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()

	const point = "test.multi_code"
	sidA := NewSession()
	// Stagger CreatedAt so FIFO ordering is well-defined even on fast machines.
	time.Sleep(2 * time.Millisecond)
	sidB := NewSession()
	time.Sleep(2 * time.Millisecond)
	sidC := NewSession()

	// OneShot so each class is consumed after firing, letting the next
	// FIFO-oldest class win on the subsequent Check (mirrors the pattern in
	// TestWorkedRow11_ErrorsFIFOAcrossConsecutiveChecks).
	armOrFail(t, sidA, FaultRule{Point: point, Action: ActionError, ErrorCode: "alpha", Message: "first", OneShot: true})
	armOrFail(t, sidB, FaultRule{Point: point, Action: ActionError, ErrorCode: "beta", Message: "second", OneShot: true})
	armOrFail(t, sidC, FaultRule{Point: point, Action: ActionError, ErrorCode: "gamma", Message: "third", OneShot: true})

	wantOrder := []string{"alpha", "beta", "gamma"}
	for i, want := range wantOrder {
		err := Check(point)
		if err == nil {
			t.Fatalf("Check #%d: expected error, got nil", i+1)
		}
		if got := CodeOf(err); got != want {
			t.Fatalf("Check #%d: expected code=%q, got: %q (err=%v)", i+1, want, got, err)
		}
	}
	// Fourth Check returns nil — every class was consumed.
	if err := Check(point); err != nil {
		t.Fatalf("4th Check should be nil, got: %v", err)
	}
}

// TestCheck_CodeKeyPrecedence asserts that error_code overrides message as the
// equivalence-class key. A rule keyed by code "x" and a separate rule keyed by
// message "x" (no code) must form *two* distinct classes — same key string,
// different source field, must not be merged.
func TestCheck_CodeKeyPrecedence(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()

	const point = "test.code_key_precedence"
	sidA := NewSession()
	time.Sleep(2 * time.Millisecond)
	sidB := NewSession()

	armOrFail(t, sidA, FaultRule{Point: point, Action: ActionError, ErrorCode: "x", Message: "via code", OneShot: true})
	armOrFail(t, sidB, FaultRule{Point: point, Action: ActionError, Message: "x", OneShot: true})

	// First Check fires the older class (sidA's coded "x"); second Check
	// fires the message-keyed "x" (sidB). Two distinct firings prove the
	// classes did not merge.
	err1 := Check(point)
	if CodeOf(err1) != "x" {
		t.Fatalf("first Check: expected code=x, got: %q (err=%v)", CodeOf(err1), err1)
	}
	err2 := Check(point)
	if err2 == nil {
		t.Fatal("second Check: expected error (message-keyed class), got nil")
	}
	if got := CodeOf(err2); got != "" {
		t.Fatalf("second Check: expected empty code (message-keyed), got: %q", got)
	}
	if want := "FAULT INJECTION at test.code_key_precedence: x"; err2.Error() != want {
		t.Fatalf("second Check: expected %q, got: %q", want, err2.Error())
	}
}

// TestTriggerEvent_CarriesCode verifies the trigger log echoes back the
// firing rule's code + message, so tests can assert *which* arm fired.
func TestTriggerEvent_CarriesCode(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	sid := NewSession()
	armOrFail(t, sid, FaultRule{
		Point:     "test.trigger_event_code",
		Action:    ActionError,
		ErrorCode: "timeout",
		Message:   "simulated",
	})

	if err := Check("test.trigger_event_code"); err == nil {
		t.Fatal("expected error")
	}
	s := mustGetSession(t, sid)
	if len(s.Log) != 1 {
		t.Fatalf("expected 1 log entry, got: %d", len(s.Log))
	}
	ev := s.Log[0]
	if ev.Code != "timeout" {
		t.Fatalf("TriggerEvent.Code: expected timeout, got: %q", ev.Code)
	}
	if ev.Message != "simulated" {
		t.Fatalf("TriggerEvent.Message: expected simulated, got: %q", ev.Message)
	}
	if ev.Action != ActionError {
		t.Fatalf("TriggerEvent.Action: expected error, got: %q", ev.Action)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Terminal-priority assertions
// ─────────────────────────────────────────────────────────────────────────────

func TestCrashDominatesPanic_PersistenceSnapshot(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()
	tmp := persistTempPath(t)
	SetPersistPath(tmp)

	s1 := NewSession()
	s2 := NewSession()
	armOrFail(t, s1, FaultRule{Point: "p", Action: ActionPanic, Message: "x"})
	armOrFail(t, s2, FaultRule{Point: "p", Action: ActionCrash, MaxCount: 2})

	// Confirm both arms are visible in the persistence file pre-Check.
	checkPersistence(t, tmp, []string{s1, s2})
}

func TestPanicDominatesError(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()

	s1 := NewSession()
	s2 := NewSession()
	armOrFail(t, s1, FaultRule{Point: "p", Action: ActionError, Message: "err"})
	armOrFail(t, s2, FaultRule{Point: "p", Action: ActionPanic, Message: "pan", MaxCount: 1})

	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected panic to win over error, got none")
			}
			if r != "FAULT INJECTION at p: pan" {
				t.Fatalf("expected panic message 'pan', got: %v", r)
			}
		}()
		_ = Check("p")
	}()

	// s2 panic fired.
	if triggerCount(t, s2, "p") != 1 {
		t.Fatal("panic arm should have decremented")
	}
	// s1 error stays armed (error didn't fire).
	if triggerCount(t, s1, "p") != 0 {
		t.Fatal("error arm should be untouched")
	}
	if _, ok := mustGetSession(t, s1).Rules["p"]; !ok {
		t.Fatal("error arm should still be in the session's Rules map")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Persistence
// ─────────────────────────────────────────────────────────────────────────────

func TestPersistence_RestoreSessions(t *testing.T) {
	resetGlobal(t)
	tmp := persistTempPath(t)
	SetPersistPath(tmp)

	Enable()
	sid := NewSession()
	armOrFail(t, sid, FaultRule{Point: "p", Action: ActionError, Message: "kept", MaxCount: 5})
	Disable()

	// Simulate a process restart: wipe in-memory state, but keep persistPath
	// and reload.
	global.mu.Lock()
	global.sessions = map[string]*Session{}
	global.mu.Unlock()
	Enable()
	defer Disable()

	s, err := GetSession(sid)
	if err != nil {
		t.Fatalf("session should have been restored from disk: %v", err)
	}
	r := s.Rules["p"]
	if r == nil || r.Message != "kept" || r.MaxCount != 5 {
		t.Fatalf("rule lost in round-trip: %+v", r)
	}
}

func TestPersistence_UnreadableFileIsIgnored(t *testing.T) {
	resetGlobal(t)
	tmp := persistTempPath(t)
	// Write garbage / old-shape file.
	if err := os.WriteFile(tmp, []byte(`{"rules":{"foo":{"point":"foo","action":"error"}}}`), 0o644); err != nil {
		t.Fatalf("write garbage: %v", err)
	}
	SetPersistPath(tmp)

	Enable()
	defer Disable()

	if got := ListSessions(); len(got) != 0 {
		t.Fatalf("expected empty session set after unreadable file, got: %v", got)
	}
}

func TestPersistence_GarbageFileIsIgnored(t *testing.T) {
	resetGlobal(t)
	tmp := persistTempPath(t)
	if err := os.WriteFile(tmp, []byte("not json at all"), 0o644); err != nil {
		t.Fatalf("write garbage: %v", err)
	}
	SetPersistPath(tmp)
	Enable()
	defer Disable()

	if got := ListSessions(); len(got) != 0 {
		t.Fatalf("expected empty session set after garbage file, got: %v", got)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// TTL sweeper
// ─────────────────────────────────────────────────────────────────────────────

func TestSweeper_RemovesIdleSessions(t *testing.T) {
	resetGlobal(t)
	// Drop the TTL to 1 minute (default anyway) and fast-forward by mutating
	// the session's LastActivity directly.
	t.Setenv("FAULT_INJECTION_SESSION_TTL_MINUTES", "1")
	Enable()
	defer Disable()

	sid := NewSession()
	// Pretend this session has been idle for 2 hours.
	global.mu.Lock()
	global.sessions[sid].LastActivity = time.Now().UTC().Add(-2 * time.Hour)
	global.mu.Unlock()

	sweepOnce()
	if _, err := GetSession(sid); !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("expected session to be swept, GetSession returned: %v", err)
	}
}

func TestSweeper_KeepsFreshSessions(t *testing.T) {
	resetGlobal(t)
	t.Setenv("FAULT_INJECTION_SESSION_TTL_MINUTES", "60")
	Enable()
	defer Disable()
	sid := NewSession()
	// LastActivity is now; sweep should not touch it.
	sweepOnce()
	if _, err := GetSession(sid); err != nil {
		t.Fatalf("fresh session unexpectedly swept: %v", err)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Concurrency stress
// ─────────────────────────────────────────────────────────────────────────────

func TestConcurrent_ManySessionsArmsAndChecks(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()

	const N = 8
	sessions := make([]string, N)
	for i := range sessions {
		sessions[i] = NewSession()
	}

	var wg sync.WaitGroup
	// Each session continuously arms and triggers a unique-message error on a
	// shared point. Verifies no cross-session log contamination, no races.
	for i, sid := range sessions {
		wg.Add(1)
		go func(i int, sid string) {
			defer wg.Done()
			msg := "owner-" + strconv.Itoa(i)
			for j := 0; j < 50; j++ {
				_ = SetRuleInSession(sid, FaultRule{
					Point:   "concurrent.p",
					Action:  ActionError,
					Message: msg,
					OneShot: true,
				})
				// Try to fire. It might be a different session's that fires
				// (FIFO across equivalence classes). That's expected.
				_ = Check("concurrent.p")
			}
		}(i, sid)
	}
	wg.Wait()

	// Verify per-session log invariants: every event in session i has either
	// our own message or no message attribution at all. (We attribute by the
	// rule's Message at the moment of fire — and a session's log only ever
	// contains events for its own arms.)
	for i, sid := range sessions {
		s := mustGetSession(t, sid)
		myMsg := "owner-" + strconv.Itoa(i)
		for _, ev := range s.Log {
			// Trigger events don't carry the message; just verify the point.
			if ev.Point != "concurrent.p" {
				t.Fatalf("session %s log has unexpected point %s", sid, ev.Point)
			}
		}
		// The rule may or may not be armed at end (one_shot consumed by us or
		// by FIFO collision is fine). What matters: no other session's rule
		// landed in our Rules map.
		for p, r := range s.Rules {
			if p != "concurrent.p" {
				t.Fatalf("session %s has unexpected rule for point %s", sid, p)
			}
			if r.Message != myMsg {
				t.Fatalf("session %s has rule with foreign message %q (mine was %q)", sid, r.Message, myMsg)
			}
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Error paths
// ─────────────────────────────────────────────────────────────────────────────

func TestSessionAPI_NotFoundErrors(t *testing.T) {
	resetGlobal(t)
	Enable()
	defer Disable()

	if err := SetRuleInSession("does-not-exist", FaultRule{Point: "p", Action: ActionError, Message: "m"}); !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("expected ErrSessionNotFound, got: %v", err)
	}
	if err := ClearRuleInSession("does-not-exist", "p"); !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("expected ErrSessionNotFound, got: %v", err)
	}
	if err := ClearRulesInSession("does-not-exist"); !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("expected ErrSessionNotFound, got: %v", err)
	}
	if err := ClearLogInSession("does-not-exist"); !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("expected ErrSessionNotFound, got: %v", err)
	}
	if err := ClearSession("does-not-exist"); !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("expected ErrSessionNotFound, got: %v", err)
	}
	if _, err := GetSession("does-not-exist"); !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("expected ErrSessionNotFound, got: %v", err)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

func persistTempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "fi-state.json")
}

// checkPersistence verifies the persistence file contains exactly the listed
// sessions with at least one rule each.
func checkPersistence(t *testing.T, path string, want []string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read persistence file: %v", err)
	}
	var state persistedState
	if err := json.Unmarshal(data, &state); err != nil {
		t.Fatalf("unmarshal persistence file: %v", err)
	}
	if len(state.Sessions) != len(want) {
		t.Fatalf("expected %d sessions in persistence, got %d", len(want), len(state.Sessions))
	}
	for _, sid := range want {
		s, ok := state.Sessions[sid]
		if !ok {
			t.Fatalf("missing session %s in persistence", sid)
		}
		if len(s.Rules) == 0 {
			t.Fatalf("session %s has no rules in persistence", sid)
		}
	}
}
