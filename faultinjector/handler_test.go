//go:build faultinjection

package faultinjector

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// newTestServer wires the handlers into an httptest.Server and ensures the
// registry is reset before every test.
func newTestServer(t *testing.T) (*httptest.Server, func()) {
	t.Helper()
	resetGlobal(t)
	Enable()
	mux := http.NewServeMux()
	mux.HandleFunc("/sessions", handleSessions)
	mux.HandleFunc("/sessions/", handleSessionPath)
	srv := httptest.NewServer(mux)
	cleanup := func() {
		srv.Close()
		Disable()
	}
	return srv, cleanup
}

// jsonBody marshals v into a body buffer or fails.
func jsonBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return bytes.NewBuffer(data)
}

func do(t *testing.T, method, urlStr string, body io.Reader) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("http: %v", err)
	}
	return resp
}

func readJSON(t *testing.T, resp *http.Response, out any) {
	t.Helper()
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if err := json.Unmarshal(data, out); err != nil {
		t.Fatalf("unmarshal body %q: %v", data, err)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Session CRUD
// ─────────────────────────────────────────────────────────────────────────────

func TestHTTP_CreateSession(t *testing.T) {
	srv, cleanup := newTestServer(t)
	defer cleanup()

	resp := do(t, http.MethodPost, srv.URL+"/sessions", nil)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var body map[string]string
	readJSON(t, resp, &body)
	if body["id"] == "" {
		t.Fatalf("expected non-empty id in response: %v", body)
	}
}

func TestHTTP_GetSession(t *testing.T) {
	srv, cleanup := newTestServer(t)
	defer cleanup()

	sid := NewSession()
	armOrFail(t, sid, FaultRule{Point: "p", Action: ActionError, Message: "m"})

	resp := do(t, http.MethodGet, srv.URL+"/sessions/"+sid, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var got Session
	readJSON(t, resp, &got)
	if got.ID != sid {
		t.Fatalf("id mismatch: got %s want %s", got.ID, sid)
	}
	if got.Rules["p"].Action != ActionError {
		t.Fatalf("rule round-trip failed: %+v", got.Rules["p"])
	}
}

func TestHTTP_GetSession_NotFound(t *testing.T) {
	srv, cleanup := newTestServer(t)
	defer cleanup()

	resp := do(t, http.MethodGet, srv.URL+"/sessions/does-not-exist", nil)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestHTTP_DeleteSession(t *testing.T) {
	srv, cleanup := newTestServer(t)
	defer cleanup()

	sid := NewSession()
	resp := do(t, http.MethodDelete, srv.URL+"/sessions/"+sid, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if _, err := GetSession(sid); !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("session should have been dropped, err=%v", err)
	}
}

func TestHTTP_ListSessions(t *testing.T) {
	srv, cleanup := newTestServer(t)
	defer cleanup()

	_ = NewSession()
	_ = NewSession()

	resp := do(t, http.MethodGet, srv.URL+"/sessions", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var body struct {
		Enabled  bool             `json:"enabled"`
		Sessions []SessionSummary `json:"sessions"`
	}
	readJSON(t, resp, &body)
	if !body.Enabled {
		t.Fatal("enabled should be true")
	}
	if len(body.Sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(body.Sessions))
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Fault arming and clearing
// ─────────────────────────────────────────────────────────────────────────────

func TestHTTP_ArmFault(t *testing.T) {
	srv, cleanup := newTestServer(t)
	defer cleanup()

	sid := NewSession()
	body := jsonBody(t, FaultRule{Point: "p", Action: ActionError, Message: "m", OneShot: true})
	resp := do(t, http.MethodPost, srv.URL+"/sessions/"+sid+"/faults", body)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	s := mustGetSession(t, sid)
	r, ok := s.Rules["p"]
	if !ok {
		t.Fatal("rule should be armed")
	}
	if !r.OneShot {
		t.Fatal("OneShot flag lost in HTTP round-trip")
	}
}

func TestHTTP_ArmFault_MissingPoint(t *testing.T) {
	srv, cleanup := newTestServer(t)
	defer cleanup()
	sid := NewSession()
	body := jsonBody(t, FaultRule{Action: ActionError, Message: "m"})
	resp := do(t, http.MethodPost, srv.URL+"/sessions/"+sid+"/faults", body)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestHTTP_ArmFault_InvalidAction(t *testing.T) {
	srv, cleanup := newTestServer(t)
	defer cleanup()
	sid := NewSession()
	body := jsonBody(t, FaultRule{Point: "p", Action: "explode"})
	resp := do(t, http.MethodPost, srv.URL+"/sessions/"+sid+"/faults", body)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestHTTP_ArmFault_SessionNotFound(t *testing.T) {
	srv, cleanup := newTestServer(t)
	defer cleanup()
	body := jsonBody(t, FaultRule{Point: "p", Action: ActionError})
	resp := do(t, http.MethodPost, srv.URL+"/sessions/missing/faults", body)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestHTTP_ClearAllRulesInSession_KeepsLog(t *testing.T) {
	srv, cleanup := newTestServer(t)
	defer cleanup()
	sid := NewSession()
	armOrFail(t, sid, FaultRule{Point: "p", Action: ActionError, Message: "m"})
	_ = Check("p")

	resp := do(t, http.MethodDelete, srv.URL+"/sessions/"+sid+"/faults", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	s := mustGetSession(t, sid)
	if len(s.Rules) != 0 {
		t.Fatalf("expected 0 rules, got %d", len(s.Rules))
	}
	if len(s.Log) != 1 {
		t.Fatalf("expected log preserved (1 entry), got %d", len(s.Log))
	}
}

func TestHTTP_ClearRuleByPoint(t *testing.T) {
	srv, cleanup := newTestServer(t)
	defer cleanup()
	sid := NewSession()
	armOrFail(t, sid, FaultRule{Point: "p", Action: ActionError, Message: "m"})
	armOrFail(t, sid, FaultRule{Point: "q", Action: ActionError, Message: "m"})

	resp := do(t, http.MethodDelete, srv.URL+"/sessions/"+sid+"/faults/p", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	s := mustGetSession(t, sid)
	if _, ok := s.Rules["p"]; ok {
		t.Fatal("rule p should be gone")
	}
	if _, ok := s.Rules["q"]; !ok {
		t.Fatal("rule q should remain")
	}
}

func TestHTTP_ClearRuleByPoint_PointWithSlashEncoded(t *testing.T) {
	srv, cleanup := newTestServer(t)
	defer cleanup()
	sid := NewSession()
	const point = "enygma.handler.Receiver.HandleEnygmaCrossTransfer.after_revert_batch"
	armOrFail(t, sid, FaultRule{Point: point, Action: ActionError, Message: "m"})

	encoded := url.PathEscape(point)
	resp := do(t, http.MethodDelete, srv.URL+"/sessions/"+sid+"/faults/"+encoded, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if _, ok := mustGetSession(t, sid).Rules[point]; ok {
		t.Fatalf("rule %q should have been cleared", point)
	}
}

func TestHTTP_ClearLog_KeepsRules(t *testing.T) {
	srv, cleanup := newTestServer(t)
	defer cleanup()
	sid := NewSession()
	armOrFail(t, sid, FaultRule{Point: "p", Action: ActionError, Message: "m"})
	_ = Check("p")

	resp := do(t, http.MethodDelete, srv.URL+"/sessions/"+sid+"/log", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	s := mustGetSession(t, sid)
	if len(s.Log) != 0 {
		t.Fatalf("expected empty log, got %d entries", len(s.Log))
	}
	if _, ok := s.Rules["p"]; !ok {
		t.Fatal("rule should survive clearLog")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Legacy routes are gone — 404 from the default mux
// ─────────────────────────────────────────────────────────────────────────────

func TestHTTP_LegacyRoutes_Return404(t *testing.T) {
	srv, cleanup := newTestServer(t)
	defer cleanup()

	cases := []struct {
		method, path string
	}{
		{http.MethodPost, "/faults"},
		{http.MethodGet, "/faults"},
		{http.MethodDelete, "/faults"},
		{http.MethodDelete, "/faults/some.point"},
	}
	for _, tc := range cases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			var body io.Reader
			if tc.method == http.MethodPost {
				body = strings.NewReader(`{"point":"p","action":"error"}`)
			}
			resp := do(t, tc.method, srv.URL+tc.path, body)
			if resp.StatusCode != http.StatusNotFound {
				t.Fatalf("expected 404 from legacy route %s %s, got %d", tc.method, tc.path, resp.StatusCode)
			}
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Method mismatch
// ─────────────────────────────────────────────────────────────────────────────

func TestHTTP_MethodNotAllowed(t *testing.T) {
	srv, cleanup := newTestServer(t)
	defer cleanup()
	sid := NewSession()

	cases := []struct {
		name, method, path string
	}{
		{"GET on /sessions/{sid}/faults", http.MethodGet, "/sessions/" + sid + "/faults"},
		{"PATCH on /sessions/{sid}", http.MethodPatch, "/sessions/" + sid},
		{"POST on /sessions/{sid}/log", http.MethodPost, "/sessions/" + sid + "/log"},
		{"GET on /sessions/{sid}/log", http.MethodGet, "/sessions/" + sid + "/log"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp := do(t, tc.method, srv.URL+tc.path, nil)
			if resp.StatusCode != http.StatusMethodNotAllowed {
				t.Fatalf("expected 405, got %d", resp.StatusCode)
			}
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Concurrent HTTP arming on the same relayer
// ─────────────────────────────────────────────────────────────────────────────

func TestHTTP_ConcurrentArmingDoesNotCorrupt(t *testing.T) {
	srv, cleanup := newTestServer(t)
	defer cleanup()

	const N = 8
	sids := make([]string, N)
	for i := range sids {
		resp := do(t, http.MethodPost, srv.URL+"/sessions", nil)
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("session creation failed: %d", resp.StatusCode)
		}
		var body map[string]string
		readJSON(t, resp, &body)
		sids[i] = body["id"]
	}

	done := make(chan struct{}, N)
	for i, sid := range sids {
		go func(i int, sid string) {
			defer func() { done <- struct{}{} }()
			for j := 0; j < 25; j++ {
				body := jsonBody(t, FaultRule{
					Point:   "concurrent.http",
					Action:  ActionError,
					Message: sid,
					OneShot: true,
				})
				_ = do(t, http.MethodPost, srv.URL+"/sessions/"+sid+"/faults", body)
			}
		}(i, sid)
	}
	for i := 0; i < N; i++ {
		<-done
	}

	// Every session must end with a well-formed rule (no torn writes).
	for _, sid := range sids {
		s := mustGetSession(t, sid)
		r := s.Rules["concurrent.http"]
		if r == nil {
			continue // the rule may have been consumed or absent — both fine
		}
		if r.Action != ActionError {
			t.Fatalf("session %s rule action corrupted: %+v", sid, r)
		}
		if r.Message != sid {
			t.Fatalf("session %s rule message corrupted: got %q want %q", sid, r.Message, sid)
		}
	}
}
