package service

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type SyncState uint

const (
	StatePending SyncState = iota
	StateMined
	StateReverted
	StateFailed
)

func (s SyncState) String() string {
	switch s {
	case StatePending:
		return "pending"
	case StateMined:
		return "mined"
	case StateReverted:
		return "reverted"
	case StateFailed:
		return "failed"
	default:
		return fmt.Sprintf("SyncState(%d)", uint(s))
	}
}

// IsTerminal reports whether the state is hard-terminal (mined or reverted).
// A retry arriving for a hard-terminal id returns the stored verdict directly
// and never re-enters chain-side recovery.
func (s SyncState) IsTerminal() bool {
	return s == StateMined || s == StateReverted
}

type SyncTransaction struct {
	ID   string
	From common.Address
	Tx   *types.Transaction

	Receipt    *types.Receipt
	RevertData []byte

	State SyncState

	// Version backs optimistic concurrency: it is the value the row carried
	// when loaded by the repository. Save asserts it and bumps it, so a
	// concurrent writer between Get and Save is detected (ErrStaleTransition)
	// rather than silently overwriting a verdict. Zero for a freshly
	// constructed, not-yet-persisted tx.
	Version int64
}

func NewSyncTransaction(id string, from common.Address, tx *types.Transaction) SyncTransaction {
	return SyncTransaction{
		ID:    id,
		From:  from,
		Tx:    tx,
		State: StatePending,
	}
}

func (s *SyncTransaction) ResolveMined(receipt *types.Receipt) error {
	if s.State != StatePending {
		return fmt.Errorf("unsupported state transition")
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("cannont transition to mined - receipt status not successful")
	}

	s.State = StateMined
	s.Receipt = receipt
	return nil
}

func (s *SyncTransaction) ResolveReverted(receipt *types.Receipt, revertData []byte) error {
	if s.State != StatePending {
		return fmt.Errorf("unsupported state transition")
	}

	if receipt.Status != types.ReceiptStatusFailed {
		return fmt.Errorf("cannont transition to reverted - receipt status not reverted ")
	}

	s.State = StateReverted
	s.Receipt = receipt
	s.RevertData = revertData
	return nil

}

func (s *SyncTransaction) Fail() error {
	if s.State != StatePending {
		return fmt.Errorf("unsupported state transition")
	}

	s.State = StateFailed
	return nil
}

func (s *SyncTransaction) Rebind(tx *types.Transaction) error {
	if s.State != StatePending {
		return fmt.Errorf("rebinding is supported only in pending state")
	}

	s.Tx = tx
	return nil
}

func (s *SyncTransaction) Reopen() error {
	if s.State != StateFailed {
		return fmt.Errorf("cannot reopen from a non-failed state")
	}

	s.State = StatePending
	return nil
}

// Function input and output types
type SignAndSendSuccess struct {
	Receipt *types.Receipt
}

type Revert struct {
	RevertData []byte
}

type SignAndSendResult struct {
	Success *SignAndSendSuccess
	Revert  *Revert
}

type DeploySuccess struct {
	Address common.Address
	Receipt *types.Receipt
}

type DeployResult struct {
	Success *DeploySuccess
	Revert  *Revert
}

type CallResult struct {
	Value  []byte
	Revert *Revert
}

type BatchItem struct {
	MsgID   string
	Data    []byte
	Address common.Address
}

type BatchItemResult struct {
	Success *SignAndSendSuccess
	Revert  *Revert
	Err     error
}
