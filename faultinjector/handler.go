//go:build faultinjection

package faultinjector

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// NewHTTPServer creates and returns an HTTP server for the fault-injection
// session API. Routes:
//
//	POST    /sessions                              create a session
//	GET     /sessions                              list all sessions (debug)
//	GET     /sessions/{sid}                        full session snapshot
//	DELETE  /sessions/{sid}                        drop session
//	POST    /sessions/{sid}/faults                 arm a rule in session
//	DELETE  /sessions/{sid}/faults                 clear all rules in session (keep log)
//	DELETE  /sessions/{sid}/faults/{point}         clear one rule in session
//	DELETE  /sessions/{sid}/log                    empty the session's trigger log
//
// The old flat `/faults` routes are intentionally absent — any caller still
// using them gets a 404 from the default mux.
func NewHTTPServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/sessions", handleSessions)
	mux.HandleFunc("/sessions/", handleSessionPath)

	return &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

// handleSessions serves /sessions (no trailing path).
func handleSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodPost:
		sid := NewSession()
		// Read back the session so the response reflects the actual CreatedAt
		// recorded on creation rather than a fresh timestamp.
		s, err := GetSession(sid)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		writeJSON(w, map[string]any{
			"id":         sid,
			"created_at": s.CreatedAt.UTC().Format(time.RFC3339Nano),
		})
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		writeJSON(w, map[string]any{
			"enabled":  IsEnabled(),
			"sessions": ListSessions(),
		})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// handleSessionPath serves /sessions/{sid}[/...]
func handleSessionPath(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rest := strings.TrimPrefix(r.URL.Path, "/sessions/")
	if rest == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "session id required"})
		return
	}
	// Split into up to 3 segments: sid, sub-resource, optional point.
	parts := strings.SplitN(rest, "/", 3)
	sid := parts[0]
	if sid == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "session id required"})
		return
	}

	switch {
	case len(parts) == 1:
		handleSession(w, r, sid)
	case len(parts) == 2 && parts[1] == "faults":
		handleSessionFaults(w, r, sid)
	case len(parts) == 3 && parts[1] == "faults":
		handleSessionFaultByPoint(w, r, sid, parts[2])
	case len(parts) == 2 && parts[1] == "log":
		handleSessionLog(w, r, sid)
	default:
		w.WriteHeader(http.StatusNotFound)
		writeJSON(w, map[string]string{"error": "unknown subresource for session"})
	}
}

func handleSession(w http.ResponseWriter, r *http.Request, sid string) {
	switch r.Method {
	case http.MethodGet:
		s, err := GetSession(sid)
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		writeJSON(w, s)
	case http.MethodDelete:
		if err := ClearSession(sid); err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		writeJSON(w, map[string]string{"status": "cleared", "id": sid})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleSessionFaults(w http.ResponseWriter, r *http.Request, sid string) {
	switch r.Method {
	case http.MethodPost:
		var rule FaultRule
		if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
			writeError(w, http.StatusBadRequest, errors.New("invalid JSON: "+err.Error()))
			return
		}
		if rule.Point == "" {
			writeError(w, http.StatusBadRequest, errors.New("point is required"))
			return
		}
		switch rule.Action {
		case ActionCrash, ActionSleep, ActionPanic, ActionError:
			// valid
		default:
			writeError(w, http.StatusBadRequest, errors.New("invalid action: must be crash, sleep, panic, or error"))
			return
		}
		if err := SetRuleInSession(sid, rule); err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, ErrSessionNotFound) {
				status = http.StatusNotFound
			}
			writeError(w, status, err)
			return
		}
		slog.Info("Fault rule armed via API",
			slog.String("session_id", sid),
			slog.String("point", rule.Point),
			slog.String("action", string(rule.Action)),
			slog.Bool("one_shot", rule.OneShot),
		)
		w.WriteHeader(http.StatusCreated)
		writeJSON(w, map[string]string{
			"status":  "armed",
			"session": sid,
			"point":   rule.Point,
			"action":  string(rule.Action),
		})
	case http.MethodDelete:
		if err := ClearRulesInSession(sid); err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		writeJSON(w, map[string]string{"status": "rules_cleared", "session": sid})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleSessionFaultByPoint(w http.ResponseWriter, r *http.Request, sid, point string) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if point == "" {
		writeError(w, http.StatusBadRequest, errors.New("point name required"))
		return
	}
	if err := ClearRuleInSession(sid, point); err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJSON(w, map[string]string{"status": "cleared", "session": sid, "point": point})
}

func handleSessionLog(w http.ResponseWriter, r *http.Request, sid string) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err := ClearLogInSession(sid); err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJSON(w, map[string]string{"status": "log_cleared", "session": sid})
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	writeJSON(w, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, v any) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("Failed to write JSON response", slog.Any("error", err))
	}
}
