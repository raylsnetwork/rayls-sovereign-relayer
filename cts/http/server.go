package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/domain"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/service"
)

const (
	publicRelayerParam  = "public_relayer"
	privateRelayerParam = "private_relayer"
)

// KeysService is the subset of service.KeysService needed by the HTTP
// handler — public address lookups only.
type KeysService interface {
	GetPublicRelayerRaylsSignPublicAddresses(ctx context.Context) (domain.PublicRelayerRaylsSignPublicAddresses, error)
	GetPrivateRelayerRaylsSignPublicAddresses(ctx context.Context) (domain.PrivateRelayerRaylsSignPublicAddresses, error)
}

type Server struct {
	httpServer *http.Server
	listener   net.Listener
	ready      atomic.Bool
}

func NewServer(addr string, keysService KeysService) (*Server, error) {
	srv := &Server{}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("GET /ready", func(w http.ResponseWriter, r *http.Request) {
		if srv.ready.Load() {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})
	mux.HandleFunc("GET /public/addresses", handleGetAddresses(keysService))

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	srv.httpServer = &http.Server{Handler: mux}
	srv.listener = lis

	return srv, nil
}

// MarkReady signals that the CTS is fully operational: keys are
// authorized and signing queues are populated. The /ready endpoint
// will start returning 200 after this call.
func (s *Server) MarkReady() {
	s.ready.Store(true)
}

func (s *Server) Serve() error {
	return s.httpServer.Serve(s.listener)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func handleGetAddresses(svc KeysService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		serviceParam := r.URL.Query().Get("service")
		if serviceParam == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{"missing service query parameter"})
			return
		}

		switch serviceParam {
		case publicRelayerParam:
			addresses, err := svc.GetPublicRelayerRaylsSignPublicAddresses(r.Context())
			if err != nil {
				handleKeysError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, publicRelayerAddressResponse{
				PublicChainAddresses:  marshalAddressList(addresses.PublicChainAddresses),
				PrivateChainAddresses: marshalAddressList(addresses.PrivateChainAddresses),
			})

		case privateRelayerParam:
			addresses, err := svc.GetPrivateRelayerRaylsSignPublicAddresses(r.Context())
			if err != nil {
				handleKeysError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, privateRelayerAddressResponse{
				PrivateHubAddresses:            marshalAddressList(addresses.PrivateHubAddresses),
				PrivateChainAddresses:          marshalAddressList(addresses.PrivateChainAddresses),
				PrivateHubDvpOperatorAddresses: marshalAddressList(addresses.PrivateHubDvpOperatorAddresses),
			})

		default:
			writeJSON(w, http.StatusNotFound, errorResponse{"unknown service"})
		}
	}
}

func handleKeysError(w http.ResponseWriter, err error) {
	slog.Warn("Failed to get rayls sign addresses", slog.Any("error", err))

	if errors.Is(err, service.ErrNoRaylsSignKeys) {
		writeJSON(w, http.StatusNotFound, errorResponse{err.Error()})
		return
	}
	writeJSON(w, http.StatusInternalServerError, errorResponse{err.Error()})
}

func marshalAddressList(addresses domain.AddressList) []string {
	out := make([]string, 0, len(addresses))
	for _, addr := range addresses {
		out = append(out, addr.Hex())
	}
	return out
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

type publicRelayerAddressResponse struct {
	PublicChainAddresses  []string `json:"public_chain_addresses"`
	PrivateChainAddresses []string `json:"private_chain_addresses"`
}

type privateRelayerAddressResponse struct {
	PrivateHubAddresses            []string `json:"private_hub_addresses"`
	PrivateChainAddresses          []string `json:"private_chain_addresses"`
	PrivateHubDvpOperatorAddresses []string `json:"private_hub_dvp_operator_addresses"`
}

type errorResponse struct {
	Error string `json:"error"`
}
