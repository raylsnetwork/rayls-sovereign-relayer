//go:build faultinjection

// Package faultinjector provides a runtime fault-injection framework for E2E
// testing. External test harnesses (e.g. the TypeScript e2e suite) create
// isolated sessions, arm rules at named code points within those sessions, and
// observe how the relayer behaves when those points fire.
//
// Fault points use path-based naming, e.g.
//
//	"enygma.handler.Receiver.HandleEnygmaCrossTransfer.after_revert_batch"
//
// Sessions isolate test state from neighbours: each test file's `before` hook
// creates a session, its `after` hook drops it. Rules, trigger logs, and
// counters live inside their session and don't bleed into other tests on the
// same relayer.
//
// Multi-arm semantics — when multiple sessions arm rules at the same point:
//
//   - sleep arms are *equivalent* regardless of duration; on a single Check the
//     framework sleeps for max(durations) and every sleep arm decrements;
//   - error/panic arms are keyed by `error_code` when it is set, falling back
//     to `message` when it is not. Arms sharing the same key form one
//     equivalence class and decrement together; arms with different keys form
//     distinct classes, and the FIFO-oldest class wins this Check;
//   - crash arms are always equivalent (no message field); all crash arms
//     decrement together when the crash fires;
//   - terminal precedence is crash > panic > error: the most destructive
//     armed action wins per Check. The arms of action types that didn't win
//     remain armed for subsequent Checks.
//
// The `error` action returns a typed *Error value carrying both the rule's
// `error_code` (machine-readable, optional) and `message` (human-readable).
// Production callers can errors.As / CodeOf the result to switch on the code
// without parsing the .Error() string.
//
// Persistence: when FAULT_INJECTION_PERSIST_PATH is set, the whole session
// table is fsynced to disk on every mutation and restored on Enable(). This
// makes one-shot crash rules safe across process restart — the rule is
// persistently consumed before os.Exit.
//
// This package is gated behind a build tag (`-tags faultinjection`) and the
// FAULT_INJECTION_ENABLED runtime flag. Never ship a binary built with the
// tag to production.
package faultinjector

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
)

// FaultAction defines what happens when a fault point is triggered.
type FaultAction string

const (
	ActionCrash FaultAction = "crash" // os.Exit(1)
	ActionSleep FaultAction = "sleep" // time.Sleep(duration_ms)
	ActionPanic FaultAction = "panic" // panic(message)
	ActionError FaultAction = "error" // Check returns error(message)
)

// FaultRule is a single armed rule scoped to one session.
type FaultRule struct {
	Point      string      `json:"point"`
	Action     FaultAction `json:"action"`
	DurationMs int         `json:"duration_ms,omitempty"`
	Message    string      `json:"message,omitempty"`
	// ErrorCode is a machine-readable discriminator the production cutpoint
	// can switch on via CodeOf(err) when this rule fires an ActionError or
	// ActionPanic. When ErrorCode is set it also becomes the equivalence-class
	// key for multi-arm grouping (overriding Message for grouping purposes,
	// while Message remains the human-readable text).
	ErrorCode string `json:"error_code,omitempty"`
	OneShot   bool   `json:"one_shot"`
	MaxCount  int    `json:"max_count,omitempty"`
	Triggered int    `json:"triggered,omitempty"`
}

// TriggerEvent records a single firing of an armed rule. Code and Message
// echo the firing rule's `error_code` / `message`, so tests can assert *which*
// arm fired in multi-class scenarios — not just *that* a fire happened.
type TriggerEvent struct {
	Point     string      `json:"point"`
	Action    FaultAction `json:"action"`
	Code      string      `json:"code,omitempty"`
	Message   string      `json:"message,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// Error is the typed error returned by Check when an ActionError rule fires.
// Production callers can errors.As(err, &fe) — or the CodeOf(err) one-liner —
// to read the test-supplied error_code without parsing the .Error() string.
//
// The .Error() text is byte-identical to the legacy formatted string when the
// rule's error_code is empty, so existing callers / log lines / tests that
// match on the literal `"FAULT INJECTION at <point>: <message>"` keep working.
type Error struct {
	Point   string
	Code    string // empty when the arming rule didn't set error_code
	Message string // human-readable; survives existing log lines
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Code != "" {
		return fmt.Sprintf("FAULT INJECTION at %s [%s]: %s", e.Point, e.Code, e.Message)
	}
	return fmt.Sprintf("FAULT INJECTION at %s: %s", e.Point, e.Message)
}

// CodeOf returns the test-supplied error_code carried by a fault-injected
// error, or "" if err isn't a fault-injection error (or carries no code).
// Idiomatic use at a cutpoint call site:
//
//	if err := faultinjector.Check(point); err != nil {
//	    switch faultinjector.CodeOf(err) {
//	    case "timeout":   // ... retry path
//	    case "db_locked": // ... revert path
//	    default:          // unknown code (or no FI): bubble up
//	    }
//	    return err
//	}
func CodeOf(err error) string {
	var fe *Error
	if errors.As(err, &fe) && fe != nil {
		return fe.Code
	}
	return ""
}

// Session is an isolated arm/log scope owned by one test (or operator).
type Session struct {
	ID           string                `json:"id"`
	CreatedAt    time.Time             `json:"created_at"`
	LastActivity time.Time             `json:"last_activity"`
	Rules        map[string]*FaultRule `json:"rules"` // keyed by point
	Log          []TriggerEvent        `json:"log"`
}

// SessionSummary is the lightweight shape returned by ListSessions.
type SessionSummary struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	LastActivity time.Time `json:"last_activity"`
	RuleCount    int       `json:"rule_count"`
	LogCount     int       `json:"log_count"`
}

// Registry holds the framework's global state.
type Registry struct {
	mu          sync.Mutex
	sessions    map[string]*Session
	enabled     bool
	persistPath string

	// sweeper goroutine: present iff a sweeper is running.
	// Disable() closes sweepStop to terminate it.
	sweepStop chan struct{}
}

var (
	global = &Registry{sessions: make(map[string]*Session)}

	ErrSessionNotFound = errors.New("fault-injection session not found")
	ErrRuleNotFound    = errors.New("fault-injection rule not found")
)

const (
	defaultSessionTTLMinutes = 60
	sweepInterval            = time.Minute
)

// SetPersistPath configures file-based persistence. When set, the session
// table is saved to this path on every mutation and restored on Enable().
func SetPersistPath(path string) {
	global.mu.Lock()
	defer global.mu.Unlock()
	global.persistPath = path
	slog.Info("Fault-injection persistence configured", slog.String("path", path))
}

// Enable activates fault-injection checking. Idempotent: a second call while
// already enabled is a no-op.
func Enable() {
	global.mu.Lock()
	if global.enabled {
		global.mu.Unlock()
		return
	}
	global.enabled = true
	if global.persistPath != "" {
		restoreLocked()
	}
	var startSweeper bool
	if global.sweepStop == nil {
		global.sweepStop = make(chan struct{})
		startSweeper = true
	}
	stopCh := global.sweepStop
	global.mu.Unlock()

	if startSweeper {
		go runSweeper(stopCh)
	}
	slog.Warn("FAULT INJECTION ENABLED — this must only be used in test environments")
}

// Disable deactivates fault-injection checking and stops the sweeper.
func Disable() {
	global.mu.Lock()
	if !global.enabled {
		global.mu.Unlock()
		return
	}
	global.enabled = false
	stopCh := global.sweepStop
	global.sweepStop = nil
	global.mu.Unlock()

	if stopCh != nil {
		close(stopCh)
	}
}

// IsEnabled returns whether fault-injection checking is currently active.
func IsEnabled() bool {
	global.mu.Lock()
	defer global.mu.Unlock()
	return global.enabled
}

// ─────────────────────────────────────────────────────────────────────────────
// Session API
// ─────────────────────────────────────────────────────────────────────────────

// NewSession creates a fresh isolated session and returns its server-assigned
// UUID. Sessions can be created regardless of whether the framework is enabled
// — Check() simply returns nil while disabled.
func NewSession() string {
	global.mu.Lock()
	defer global.mu.Unlock()
	sid := uuid.New().String()
	now := time.Now().UTC()
	global.sessions[sid] = &Session{
		ID:           sid,
		CreatedAt:    now,
		LastActivity: now,
		Rules:        map[string]*FaultRule{},
		Log:          []TriggerEvent{},
	}
	persistLocked()
	slog.Info("Fault-injection session created", slog.String("session_id", sid))
	return sid
}

// SetRuleInSession arms a rule on a point inside the given session. Any
// existing rule on the same point in this session is overwritten (within a
// single session there is no multi-arm — only across sessions).
func SetRuleInSession(sid string, rule FaultRule) error {
	global.mu.Lock()
	defer global.mu.Unlock()
	s, ok := global.sessions[sid]
	if !ok {
		return ErrSessionNotFound
	}
	r := rule
	r.Triggered = 0
	s.Rules[rule.Point] = &r
	s.LastActivity = time.Now().UTC()
	persistLocked()
	slog.Warn("Fault rule armed",
		slog.String("session_id", sid),
		slog.String("point", rule.Point),
		slog.String("action", string(rule.Action)),
	)
	return nil
}

// ClearRuleInSession disarms a single rule. Returns ErrRuleNotFound if the
// point wasn't armed in this session.
func ClearRuleInSession(sid, point string) error {
	global.mu.Lock()
	defer global.mu.Unlock()
	s, ok := global.sessions[sid]
	if !ok {
		return ErrSessionNotFound
	}
	if _, ok := s.Rules[point]; !ok {
		return ErrRuleNotFound
	}
	delete(s.Rules, point)
	s.LastActivity = time.Now().UTC()
	persistLocked()
	return nil
}

// ClearRulesInSession drops every rule in the session. The session and its
// trigger log are kept intact — use ClearSession to also drop the log.
func ClearRulesInSession(sid string) error {
	global.mu.Lock()
	defer global.mu.Unlock()
	s, ok := global.sessions[sid]
	if !ok {
		return ErrSessionNotFound
	}
	s.Rules = map[string]*FaultRule{}
	s.LastActivity = time.Now().UTC()
	persistLocked()
	return nil
}

// ClearLogInSession empties the trigger log without dropping any rules. Use
// case: arm rules, run scenario A, ClearLog, run scenario B, assert only
// scenario B's triggers are present.
func ClearLogInSession(sid string) error {
	global.mu.Lock()
	defer global.mu.Unlock()
	s, ok := global.sessions[sid]
	if !ok {
		return ErrSessionNotFound
	}
	s.Log = []TriggerEvent{}
	s.LastActivity = time.Now().UTC()
	persistLocked()
	return nil
}

// ClearSession removes the session entirely: rules, log, metadata.
func ClearSession(sid string) error {
	global.mu.Lock()
	defer global.mu.Unlock()
	if _, ok := global.sessions[sid]; !ok {
		return ErrSessionNotFound
	}
	delete(global.sessions, sid)
	persistLocked()
	slog.Info("Fault-injection session dropped", slog.String("session_id", sid))
	return nil
}

// GetSession returns a deep-copy snapshot of the session.
func GetSession(sid string) (*Session, error) {
	global.mu.Lock()
	defer global.mu.Unlock()
	s, ok := global.sessions[sid]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return cloneSession(s), nil
}

// ListSessions returns lightweight summaries of every session.
func ListSessions() []SessionSummary {
	global.mu.Lock()
	defer global.mu.Unlock()
	out := make([]SessionSummary, 0, len(global.sessions))
	for _, s := range global.sessions {
		out = append(out, SessionSummary{
			ID:           s.ID,
			CreatedAt:    s.CreatedAt,
			LastActivity: s.LastActivity,
			RuleCount:    len(s.Rules),
			LogCount:     len(s.Log),
		})
	}
	return out
}

// cloneSession returns a deep copy. Must be called with global.mu held.
func cloneSession(s *Session) *Session {
	out := &Session{
		ID:           s.ID,
		CreatedAt:    s.CreatedAt,
		LastActivity: s.LastActivity,
		Rules:        make(map[string]*FaultRule, len(s.Rules)),
		Log:          make([]TriggerEvent, len(s.Log)),
	}
	for k, v := range s.Rules {
		r := *v
		out.Rules[k] = &r
	}
	copy(out.Log, s.Log)
	return out
}

// ─────────────────────────────────────────────────────────────────────────────
// Check — the multi-arm execution engine
// ─────────────────────────────────────────────────────────────────────────────

// armedRule pairs a session with the rule it owns at the current point.
type armedRule struct {
	session *Session
	rule    *FaultRule
}

// Check is called at instrumented points in the relayer code. See the package
// doc for the multi-arm execution semantics.
func Check(point string) error {
	global.mu.Lock()

	if !global.enabled {
		global.mu.Unlock()
		return nil
	}

	var sleeps, errs, panics, crashes []armedRule
	for _, s := range global.sessions {
		r, ok := s.Rules[point]
		if !ok {
			continue
		}
		if r.MaxCount > 0 && r.Triggered >= r.MaxCount {
			continue
		}
		switch r.Action {
		case ActionSleep:
			sleeps = append(sleeps, armedRule{s, r})
		case ActionError:
			errs = append(errs, armedRule{s, r})
		case ActionPanic:
			panics = append(panics, armedRule{s, r})
		case ActionCrash:
			crashes = append(crashes, armedRule{s, r})
		}
	}

	if len(sleeps)+len(errs)+len(panics)+len(crashes) == 0 {
		global.mu.Unlock()
		return nil
	}

	// Sleep equivalence class: any sleep arm matches. Effective duration is max.
	var maxSleep time.Duration
	for _, a := range sleeps {
		d := time.Duration(a.rule.DurationMs) * time.Millisecond
		if d > maxSleep {
			maxSleep = d
		}
	}
	for _, a := range sleeps {
		markFiredLocked(a, point, "", "")
	}

	// Terminal selection — crash > panic > error.
	var (
		winningAction  FaultAction
		winningCode    string
		winningMessage string
		fired          []armedRule
	)
	switch {
	case len(crashes) > 0:
		winningAction = ActionCrash
		fired = crashes // single equivalence class — no code/message disambiguation
	case len(panics) > 0:
		winningAction = ActionPanic
		fired, winningCode, winningMessage = pickFIFOClassByErrorKey(panics)
	case len(errs) > 0:
		winningAction = ActionError
		fired, winningCode, winningMessage = pickFIFOClassByErrorKey(errs)
	}
	for _, a := range fired {
		markFiredLocked(a, point, winningCode, winningMessage)
	}

	persistLocked()
	global.mu.Unlock()

	slog.Warn("FAULT INJECTION TRIGGERED",
		slog.String("point", point),
		slog.Int("sleep_arms", len(sleeps)),
		slog.String("winning_action", string(winningAction)),
		slog.String("winning_code", winningCode),
		slog.Int("terminal_arms", len(fired)),
	)

	if maxSleep > 0 {
		slog.Warn("FAULT INJECTION: SLEEPING",
			slog.String("point", point),
			slog.Duration("duration", maxSleep),
		)
		time.Sleep(maxSleep)
	}

	switch winningAction {
	case ActionCrash:
		slog.Error("FAULT INJECTION: CRASHING PROCESS", slog.String("point", point))
		os.Exit(1)
	case ActionPanic:
		// Panic message preserves the legacy text exactly when no error_code
		// is set; with a code, the bracketed prefix mirrors Error.Error().
		if winningCode != "" {
			panic(fmt.Sprintf("FAULT INJECTION at %s [%s]: %s", point, winningCode, winningMessage))
		}
		panic(fmt.Sprintf("FAULT INJECTION at %s: %s", point, winningMessage))
	case ActionError:
		return &Error{Point: point, Code: winningCode, Message: winningMessage}
	}
	return nil
}

// pickFIFOClassByErrorKey groups by the rule's discriminator key — its
// ErrorCode when set, else its Message — and returns every armedRule in the
// class whose oldest Session.CreatedAt is the smallest. Returns the winning
// (code, message) tuple drawn from the first armedRule of the winning class
// (within a class all rules share the same key, but Message may vary when
// ErrorCode is what tied them together; first-wins is deterministic and
// matches the FIFO selection).
func pickFIFOClassByErrorKey(candidates []armedRule) ([]armedRule, string, string) {
	if len(candidates) == 0 {
		return nil, "", ""
	}
	classes := map[string][]armedRule{}
	oldest := map[string]time.Time{}
	for _, a := range candidates {
		key := classKey(a.rule)
		classes[key] = append(classes[key], a)
		if t, ok := oldest[key]; !ok || a.session.CreatedAt.Before(t) {
			oldest[key] = a.session.CreatedAt
		}
	}
	var (
		winnerKey  string
		winnerTime time.Time
		first      = true
	)
	for key, t := range oldest {
		if first || t.Before(winnerTime) {
			winnerKey = key
			winnerTime = t
			first = false
		}
	}
	winners := classes[winnerKey]
	// Code/message come from any winner (all share the class key; Message is
	// stable when the key was Message, free when the key was ErrorCode — pick
	// the first arm's values for determinism).
	winCode := winners[0].rule.ErrorCode
	winMsg := winners[0].rule.Message
	return winners, winCode, winMsg
}

// classKey returns the equivalence-class key for an error/panic rule. The
// key is namespaced by source field so that arms keyed on `error_code` and
// arms keyed on `message` never collide on the same literal string (e.g. a
// rule with ErrorCode="x" and another with Message="x" must form two
// distinct classes). See the package doc for semantics.
func classKey(r *FaultRule) string {
	if r.ErrorCode != "" {
		return "C:" + r.ErrorCode
	}
	return "M:" + r.Message
}

// markFiredLocked records that an arm fired: bump counter, append log event,
// touch LastActivity, consume one-shot. `code` and `message` are the values
// the firing surfaced (empty for sleep). Must hold global.mu.
func markFiredLocked(a armedRule, point, code, message string) {
	a.rule.Triggered++
	now := time.Now().UTC()
	a.session.Log = append(a.session.Log, TriggerEvent{
		Point:     point,
		Action:    a.rule.Action,
		Code:      code,
		Message:   message,
		Timestamp: now.Format(time.RFC3339Nano),
	})
	a.session.LastActivity = now
	if a.rule.OneShot {
		delete(a.session.Rules, point)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Persistence
// ─────────────────────────────────────────────────────────────────────────────

type persistedState struct {
	Sessions map[string]*Session `json:"sessions"`
}

// persistLocked writes the current session table to disk. Must hold global.mu.
func persistLocked() {
	if global.persistPath == "" {
		return
	}
	state := persistedState{Sessions: global.sessions}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		slog.Error("Failed to marshal fault-injection sessions",
			slog.Any("error", err))
		return
	}
	if err := os.WriteFile(global.persistPath, data, 0o600); err != nil {
		slog.Error("Failed to persist fault-injection sessions",
			slog.Any("error", err),
			slog.String("path", global.persistPath))
	}
}

// restoreLocked attempts to read the persistence file. On any parse error or
// missing/unrecognised shape, logs a warning and leaves the session map empty.
// Must hold global.mu.
func restoreLocked() {
	data, err := os.ReadFile(global.persistPath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		slog.Warn("Could not read fault-injection persistence file; starting with empty session set",
			slog.Any("error", err),
			slog.String("path", global.persistPath),
		)
		return
	}
	var state persistedState
	if err := json.Unmarshal(data, &state); err != nil {
		slog.Warn("Ignoring unreadable fault-injection persistence file; starting with empty session set",
			slog.Any("error", err),
			slog.String("path", global.persistPath),
		)
		return
	}
	if state.Sessions == nil {
		slog.Warn("Ignoring unreadable fault-injection persistence file (no sessions key); starting with empty session set",
			slog.String("path", global.persistPath),
		)
		return
	}
	global.sessions = state.Sessions
	for _, s := range global.sessions {
		if s.Rules == nil {
			s.Rules = map[string]*FaultRule{}
		}
		if s.Log == nil {
			s.Log = []TriggerEvent{}
		}
	}
	slog.Info("Restored persisted fault-injection sessions",
		slog.Int("count", len(state.Sessions)),
	)
}

// ─────────────────────────────────────────────────────────────────────────────
// Sweeper — drops sessions inactive longer than the configured TTL
// ─────────────────────────────────────────────────────────────────────────────

func sessionTTL() time.Duration {
	if v := os.Getenv("FAULT_INJECTION_SESSION_TTL_MINUTES"); v != "" {
		if minutes, err := strconv.Atoi(v); err == nil && minutes > 0 {
			return time.Duration(minutes) * time.Minute
		}
	}
	return defaultSessionTTLMinutes * time.Minute
}

func runSweeper(stop chan struct{}) {
	ticker := time.NewTicker(sweepInterval)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			sweepOnce()
		}
	}
}

// sweepOnce drops any session whose LastActivity is older than the TTL.
// Exposed for unit tests via a separate test build-tag file; production code
// only triggers it via the ticker.
func sweepOnce() {
	cutoff := time.Now().UTC().Add(-sessionTTL())
	global.mu.Lock()
	defer global.mu.Unlock()
	var dropped []string
	for sid, s := range global.sessions {
		if s.LastActivity.Before(cutoff) {
			delete(global.sessions, sid)
			dropped = append(dropped, sid)
		}
	}
	if len(dropped) > 0 {
		persistLocked()
		for _, sid := range dropped {
			slog.Info("Fault-injection session swept (TTL exceeded)",
				slog.String("session_id", sid),
			)
		}
	}
}
