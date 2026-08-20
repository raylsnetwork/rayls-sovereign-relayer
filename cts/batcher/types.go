// Package batcher owns the CTS-side sign/send/receipt pipeline for a single
// signing identity (e.g. privatehub, privatenode, dvpoperator). A NATS-fed
// IngestionService writes pending rows to `cts_transaction`; a self-polling
// BatcherService signs+broadcasts them; a self-polling ReceipterService
// fetches receipts and publishes terminal results; a self-polling
// ReaperService dead-letters transactions stuck without a receipt.
package batcher

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/ethrpc"
	"github.com/raylsnetwork/rayls-sovereign-relayer/msgqueue"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

// Status is the state-machine stage of a single transaction row.
type Status string

const (
	StatusPending  Status = "pending"
	StatusSent     Status = "sent"
	StatusFinished Status = "finished"
	StatusFailed   Status = "failed"
)

// CTSTransaction is the internal domain record materialised from a DB
// row. Not exposed to NATS — the wire format is IngestRequest / Result.
type CTSTransaction struct {
	CorrelationID   string
	Identity        string
	MessageType     string
	Address         common.Address
	Calldata        []byte
	Status          Status
	TxHash          common.Hash
	ReceiptStatus   uint64
	RevertData      []byte
	ErrorReason     string
	SendAttempts    int
	ReceiptAttempts int
	SentAt          *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

//go:generate moq --pkg batcher_test -out mocks_test.go . MessageConsumer ResultPublisher Batcher Repository
type MessageConsumer interface {
	Fetch(ctx context.Context, count int) ([]msgqueue.Message[types.TxRequest], error)
}

type ResultPublisher interface {
	Push(ctx context.Context, r types.TxResult) error
}

type Batcher interface {
	Send(ctx context.Context, inputs []ethrpc.TransactionInput) ([]ethrpc.SendResult, error)
	GetReceipts(ctx context.Context, hashes []common.Hash) ([]ethrpc.ReceiptResult, error)
}

// TxKey identifies a single cts_transaction row. The composite primary
// key is (correlation_id, message_type), so a single shared id may
// carry multiple message types (crosschain.atomic + destination-unlock
// + destination-revert, etc.) — each its own row.
type TxKey struct {
	CorrelationID string
	MessageType   string
}

// KeyOf extracts the row key from a CTSTransaction.
func KeyOf(t CTSTransaction) TxKey {
	return TxKey{CorrelationID: t.CorrelationID, MessageType: t.MessageType}
}

type Repository interface {
	// Insert is idempotent via ON CONFLICT (correlation_id, message_type)
	// DO NOTHING.
	Insert(ctx context.Context, req types.TxRequest) error

	// ClaimPending returns up to `limit` pending rows for `identity`,
	// using FOR UPDATE SKIP LOCKED as defense in depth. Rows are NOT
	// transitioned by this call — per-row MarkX happens after the batch
	// Send/GetReceipts call returns.
	ClaimPending(ctx context.Context, identity string, limit int) ([]CTSTransaction, error)
	ClaimSent(ctx context.Context, identity string, limit int) ([]CTSTransaction, error)
	ClaimStuck(ctx context.Context, identity string, olderThan time.Time, limit int) ([]CTSTransaction, error)

	MarkSent(ctx context.Context, key TxKey, hash common.Hash) error
	// MarkResent refreshes sent_at on an already-sent row so the reaper's
	// stuck-claim cutoff respects the re-broadcast cycle. Does not touch
	// tx_hash — the receipter keeps polling the original hash so a delayed
	// receipt of the first broadcast (the common case under transient
	// receipter blips) still resolves the row normally.
	MarkResent(ctx context.Context, key TxKey) error
	MarkFinishedSuccess(ctx context.Context, key TxKey, receiptStatus uint64) error
	MarkFinishedRevert(ctx context.Context, key TxKey, revertData []byte) error
	MarkFailed(ctx context.Context, key TxKey, reason string) error

	IncrementSendAttempts(ctx context.Context, keys []TxKey) error
	IncrementReceiptAttempts(ctx context.Context, keys []TxKey) error
}
