// Decommissioning Teleport (vanilla, atomic).

package app

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type HealthCheckConfig struct {
	Addr string
	Path string
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

type ErrorResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Error     string `json:"error"`
}

func (n *PublicRelayer) initializeHealthCheckServer(config HealthCheckConfig) {
	mux := &http.ServeMux{}
	mux.Handle(config.Path, NewHealthCheck(n))

	server := &http.Server{
		Addr:              config.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
	n.healthCheckServer = server
}

type HealthCheck struct {
	pubrelayer *PublicRelayer
}

func NewHealthCheck(pubrelayer *PublicRelayer) *HealthCheck {
	return &HealthCheck{
		pubrelayer: pubrelayer,
	}
}

func (h *HealthCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := h.pubrelayer.pool.Ping(r.Context()); err != nil {
		writeServiceUnavailable(w, fmt.Errorf("failed to ping database: %w", err))
		return
	}

	if _, err := h.pubrelayer.privateClient.BlockNumber(r.Context()); err != nil {
		writeServiceUnavailable(w, fmt.Errorf("failed to ping private client: %w", err))
		return
	}

	if _, err := h.pubrelayer.publicClient.BlockNumber(r.Context()); err != nil {
		writeServiceUnavailable(w, fmt.Errorf("failed to ping public client: %w", err))
		return
	}

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to encode health response", slog.Any("error", err))
	}
}

func writeServiceUnavailable(w http.ResponseWriter, err error) {
	response := ErrorResponse{
		Status:    "unhealthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Error:     err.Error(),
	}

	w.WriteHeader(http.StatusServiceUnavailable)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to encode error response", slog.Any("error", err))
	}
}
