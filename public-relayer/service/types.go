// Decommissioning Teleport (vanilla, atomic).

package service

import (
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/RNMessageDispatcherV1"
)

var ErrDuplicateMessageID = errors.New("record with message id already exists")

type Message struct {
	ID          common.Hash
	FromAddress common.Address
	ToChainID   *big.Int
	ToAddress   common.Address
	Data        RNMessageDispatcherV1.RaylsNodeMessage
}

func (m Message) GetID() string {
	return m.ID.Hex()
}

type MessageRecordStatus int

const (
	// MessageRecordStatusNew means the row was created when the generator
	// published the forward batch; result callbacks have not fired yet.
	MessageRecordStatusNew MessageRecordStatus = iota
	// MessageRecordStatusSucceeded is a terminal success: the forward tx
	// was mined OK.
	MessageRecordStatusSucceeded
	// MessageRecordStatusFailed means the forward tx failed or reverted
	// and the revert has been published; we're awaiting its result.
	MessageRecordStatusFailed
	// MessageRecordStatusRevertSucceeded is a terminal state: the revert
	// tx mined OK.
	MessageRecordStatusRevertSucceeded
	// MessageRecordStatusRevertFailed is a terminal error state: the
	// revert itself failed/reverted. Needs ops attention.
	MessageRecordStatusRevertFailed
)

// MessageRecord is a per-message lifecycle row persisted by the
// GeneratorService. It exists for observability — the forward/revert
// callbacks update it through the legal state transitions:
//
//	New ──success──▶ Succeeded                             (terminal)
//	New ──failure──▶ Failed ──success──▶ RevertSucceeded   (terminal)
//	                        └─failure──▶ RevertFailed      (terminal)
type MessageRecord struct {
	ID          string
	Status      MessageRecordStatus
	ForwardHash common.Hash
	RevertHash  common.Hash
	Error       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ForwardResultUpdate carries a single row's post-forward-callback state
// transition. The repository applies these in a batched UPDATE.
type ForwardResultUpdate struct {
	ID     string
	Status MessageRecordStatus
	Hash   common.Hash
	Error  string
}

// RevertResultUpdate mirrors ForwardResultUpdate for the revert callback.
type RevertResultUpdate struct {
	ID     string
	Status MessageRecordStatus
	Hash   common.Hash
	Error  string
}
