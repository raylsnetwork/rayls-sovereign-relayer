package healthcheck

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//go:generate moq --pkg healthcheck_test -out healthcheck_mock_test.go . DatabaseClient EthereumClient NATSClient

type DatabaseClient interface {
	Ping(ctx context.Context) error
}

type EthereumClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
}

type NATSClient interface {
	IsConnected() bool
}

type Config struct {
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

const (
	// serverReadTimeout is the maximum duration for reading the entire request.
	serverReadTimeout = 5 * time.Second
	// serverWriteTimeout is the maximum duration before timing out writes of the response.
	serverWriteTimeout = 10 * time.Second
	// serverIdleTimeout is the maximum amount of time to wait for the next request when keep-alives are enabled.
	serverIdleTimeout = 120 * time.Second
	// serverReadHeaderTimeout is the amount of time allowed to read request headers.
	serverReadHeaderTimeout = 5 * time.Second
)

type Healthcheck struct {
	server http.Server

	natsConn        NATSClient
	dbConn          DatabaseClient
	privateHubConn  EthereumClient
	privateNodeConn EthereumClient
}

func New(
	config Config,
	natsConn NATSClient,
	dbConn DatabaseClient,
	privateHubConn EthereumClient,
	privateNodeConn EthereumClient,
) *Healthcheck {
	healthcheck := &Healthcheck{
		natsConn:        natsConn,
		dbConn:          dbConn,
		privateHubConn:  privateHubConn,
		privateNodeConn: privateNodeConn,
	}

	mux := &http.ServeMux{}
	mux.Handle(config.Path, healthcheck)

	healthcheck.server = http.Server{
		Addr:              ":" + config.Addr,
		Handler:           mux,
		ReadTimeout:       serverReadTimeout,
		WriteTimeout:      serverWriteTimeout,
		IdleTimeout:       serverIdleTimeout,
		ReadHeaderTimeout: serverReadHeaderTimeout,
	}

	return healthcheck
}

func (h *Healthcheck) ListenAndServe() error {
	return h.server.ListenAndServe()
}

func (h *Healthcheck) Shutdown(ctx context.Context) error {
	// Shutdown HTTP server first to stop accepting new requests
	return h.server.Shutdown(ctx)
}

func (h *Healthcheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only accept GET and HEAD methods
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := h.dbConn.Ping(r.Context()); err != nil {
		writeServiceUnavailable(w, fmt.Errorf("failed to ping database: %w", err))
		return
	}

	if _, err := h.privateHubConn.BlockNumber(r.Context()); err != nil {
		writeServiceUnavailable(w, fmt.Errorf("failed to ping private hub client: %w", err))
		return
	}

	if _, err := h.privateNodeConn.BlockNumber(r.Context()); err != nil {
		writeServiceUnavailable(w, fmt.Errorf("failed to ping private node client: %w", err))
		return
	}

	// Check NATS connection
	if !h.natsConn.IsConnected() {
		writeServiceUnavailable(w, fmt.Errorf("NATS connection is not active"))
		return
	}

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Can't change status code at this point, but at least don't panic
		return
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
		// Can't change status code at this point, but at least don't panic
		return
	}
}
