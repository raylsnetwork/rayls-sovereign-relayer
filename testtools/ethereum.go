package testtools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// FakeEthServerConfig configures the fake Ethereum JSON-RPC server.
type FakeEthServerConfig struct {
	ChainID  *big.Int
	GasPrice *big.Int

	// DeploymentRegistries maps a deployment proxy registry address to the
	// contract names and addresses that GetAllContracts() should return.
	DeploymentRegistries map[common.Address]DeploymentRegistry
}

// DeploymentRegistry holds the names and addresses returned by GetAllContracts().
type DeploymentRegistry struct {
	Names     []string
	Addresses []common.Address
}

// FakeEthServer is a fake Ethereum JSON-RPC server for testing.
type FakeEthServer struct {
	Server *httptest.Server
	URL    string

	config                  FakeEthServerConfig
	getAllContractsSelector []byte
	encodedRegistries       map[common.Address][]byte
}

type jsonRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *jsonRPCError   `json:"error,omitempty"`
}

type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewFakeEthServer creates a fake Ethereum JSON-RPC server.
// The server handles eth_chainId, eth_gasPrice, eth_getCode, eth_call, and eth_blockNumber.
func NewFakeEthServer(t testing.TB, config FakeEthServerConfig) *FakeEthServer {
	t.Helper()

	s := &FakeEthServer{
		config:                  config,
		getAllContractsSelector: crypto.Keccak256([]byte("getAllContracts()"))[:4],
		encodedRegistries:       make(map[common.Address][]byte),
	}

	// Pre-encode all registry responses
	for addr, reg := range config.DeploymentRegistries {
		encoded, err := encodeGetAllContracts(reg)
		if err != nil {
			t.Fatalf("failed to ABI-encode GetAllContracts for %s: %s", addr.Hex(), err)
		}
		s.encodedRegistries[addr] = encoded
	}

	s.Server = httptest.NewServer(http.HandlerFunc(s.serveHTTP))
	s.URL = s.Server.URL
	t.Cleanup(func() { s.Server.Close() })

	return s
}

func (s *FakeEthServer) serveHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Detect batch vs single request
	trimmed := strings.TrimSpace(string(body))
	if strings.HasPrefix(trimmed, "[") {
		var reqs []jsonRPCRequest
		if err := json.Unmarshal(body, &reqs); err != nil {
			http.Error(w, "bad JSON", http.StatusBadRequest)
			return
		}
		resps := make([]jsonRPCResponse, len(reqs))
		for i, req := range reqs {
			resps[i] = s.handleRPC(req)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resps)
	} else {
		var req jsonRPCRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "bad JSON", http.StatusBadRequest)
			return
		}
		resp := s.handleRPC(req)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func (s *FakeEthServer) handleRPC(req jsonRPCRequest) jsonRPCResponse {
	switch req.Method {
	case "eth_chainId":
		return s.respond(req.ID, hexutil.EncodeBig(s.config.ChainID))

	case "eth_gasPrice":
		return s.respond(req.ID, hexutil.EncodeBig(s.config.GasPrice))

	case "eth_getCode":
		// Return non-empty bytecode so contracts.CreateContract succeeds
		return s.respond(req.ID, "0x01")

	case "eth_blockNumber":
		// Return block 0 so the listener stays idle
		return s.respond(req.ID, "0x0")

	case "eth_getLogs":
		// Return empty logs so the listener finds nothing to process
		return s.respond(req.ID, []any{})

	case "eth_call":
		return s.handleEthCall(req)

	default:
		return s.respondError(req.ID, -32601, fmt.Sprintf("method %q not supported by fake server", req.Method))
	}
}

func (s *FakeEthServer) handleEthCall(req jsonRPCRequest) jsonRPCResponse {
	// Parse params: [{to, data, ...}, "latest"]
	var params []json.RawMessage
	if err := json.Unmarshal(req.Params, &params); err != nil || len(params) < 1 {
		return s.respondError(req.ID, -32602, "invalid params for eth_call")
	}

	var callObj struct {
		To    string `json:"to"`
		Data  string `json:"data"`
		Input string `json:"input"`
	}
	if err := json.Unmarshal(params[0], &callObj); err != nil {
		return s.respondError(req.ID, -32602, "invalid call object")
	}

	toAddr := common.HexToAddress(callObj.To)
	// go-ethereum sends call data as "input", not "data"
	callData := callObj.Input
	if callData == "" {
		callData = callObj.Data
	}
	data := common.FromHex(callData)

	// Check if this is a GetAllContracts() call by matching the 4-byte selector
	const ethMethodSelectorSize = 4
	if len(data) >= ethMethodSelectorSize {
		selector := data[:ethMethodSelectorSize]
		if bytes.Equal(selector, s.getAllContractsSelector) {
			if encoded, ok := s.encodedRegistries[toAddr]; ok {
				return s.respond(req.ID, hexutil.Encode(encoded))
			}
		}
	}

	// For any other eth_call (contract constructor checks, etc), return empty
	return s.respond(req.ID, "0x")
}

func (s *FakeEthServer) respond(id json.RawMessage, result interface{}) jsonRPCResponse {
	return jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

func (s *FakeEthServer) respondError(id json.RawMessage, code int, message string) jsonRPCResponse {
	return jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &jsonRPCError{Code: code, Message: message},
	}
}

// ErrorStringRevertData builds the ABI payload of a standard Solidity `Error(string)` revert
// (selector 0x08c379a0) carrying the given reason — the exact on-chain shape that surfaces as a
// *contractclient.ErrorWithRevertData. Shared so the contractclient and enygma retry tests don't each
// hand-roll the same encoding.
func ErrorStringRevertData(t testing.TB, reason string) []byte {
	t.Helper()
	strType, err := abi.NewType("string", "", nil)
	if err != nil {
		t.Fatalf("ErrorStringRevertData: string type: %s", err)
	}
	encoded, err := abi.Arguments{{Type: strType}}.Pack(reason)
	if err != nil {
		t.Fatalf("ErrorStringRevertData: pack reason %q: %s", reason, err)
	}
	return append([]byte{0x08, 0xc3, 0x79, 0xa0}, encoded...)
}

func encodeGetAllContracts(reg DeploymentRegistry) ([]byte, error) {
	strArrayType, err := abi.NewType("string[]", "", nil)
	if err != nil {
		return nil, fmt.Errorf("string[] type: %w", err)
	}
	addrArrayType, err := abi.NewType("address[]", "", nil)
	if err != nil {
		return nil, fmt.Errorf("address[] type: %w", err)
	}

	args := abi.Arguments{
		{Type: strArrayType},
		{Type: addrArrayType},
	}

	return args.Pack(reg.Names, reg.Addresses)
}
