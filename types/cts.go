package types

import (
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

// TxRequest is the wire format the relayer publishes on
// `cts.send.<identity>`. CTS ingests it, signs, broadcasts, and
// eventually publishes a matching TxResult on `cts.result.<identity>`.
// The caller-supplied CorrelationID is the DB primary key, the NATS
// dedup key, and the handle the relayer uses to correlate the eventual
// TxResult back to its own domain state.
type TxRequest struct {
	CorrelationID string         `json:"correlation_id"`
	Identity      string         `json:"identity"`
	MessageType   string         `json:"message_type"`
	Address       common.Address `json:"address"`
	Calldata      []byte         `json:"calldata"`
}

// GetID is the NATS dedup key. Multiple message types share the same
// `cts.send.<identity>` subject (forward, unlock, destination-revert,
// source-revert, early-revert, privatehub.execute), and they all key
// off the same shared-id-shaped CorrelationID. CorrelationID alone
// would collide and silently drop every message after the first one
// per shared id. Combining MessageType + CorrelationID keeps retries
// dedupping (same type + same id) while letting genuinely-different
// follow-up messages pass through.
func (r TxRequest) GetID() string { return r.MessageType + ":" + r.CorrelationID }

// TxResultKind discriminates the terminal outcome of a TxRequest:
//   - TxResultSuccess: tx mined successfully; Receipt is populated
//   - TxResultRevert:  mined but reverted, or caught as a pre-flight
//                      revert by gas estimation; RevertData carries
//                      the ABI-encoded revert bytes
//   - TxResultFailed:  terminal non-mined failure (dead-lettered after
//                      retries, stuck past threshold, etc.); ErrorReason
//                      is populated
type TxResultKind string

const (
	TxResultSuccess TxResultKind = "success"
	TxResultRevert  TxResultKind = "revert"
	TxResultFailed  TxResultKind = "failed"
)

// TxResult is the wire format published on `cts.result.<identity>` once
// a TxRequest reaches a terminal state. MessageType echoes the request
// so the relayer can dispatch to the right callback.
//
// Receipt is fully decoded by the time a callback sees it — go-ethereum's
// *Receipt implements JSON marshal/unmarshal so it survives the NATS
// hop transparently. Callbacks read Receipt.BlockHash, Receipt.Logs,
// etc. directly without a separate decode step.
type TxResult struct {
	CorrelationID string            `json:"correlation_id"`
	MessageType   string            `json:"message_type"`
	Identity      string            `json:"identity"`
	Kind          TxResultKind      `json:"kind"`
	TxHash        common.Hash       `json:"tx_hash,omitempty"`
	Receipt       *ethTypes.Receipt `json:"receipt,omitempty"`
	RevertData    []byte            `json:"revert_data,omitempty"`
	ErrorReason   string            `json:"error_reason,omitempty"`
}

// GetID — same dedup logic as TxRequest. Multiple result types share
// the `cts.result.<identity>` subject for the same shared id; without
// MessageType in the key, only the first result per shared id would
// reach the router.
func (r TxResult) GetID() string { return r.MessageType + ":" + r.CorrelationID }
