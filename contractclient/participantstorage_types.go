package contractclient

import (
	"fmt"
)

var (
	ErrNoChainInfo             = NewParticipantStorageClientError("no chain info for chain ID")
	ErrNoAuditInfo             = NewParticipantStorageClientError("no audit info for chain ID")
	ErrOutdatedKeyAgreement    = NewParticipantStorageClientError("key agreement already exists for newer block number")
	ErrNoPaymentSpendPublicKey = NewParticipantStorageClientError("no payment spend public key for chain ID")
)

// outdatedKeyAgreementRevertReason is the substring emitted by the
// ParticipantStorage contract's require() when InitiateKeyAgreement is called
// with a block number older than the latest stored agreement. It is matched
// both on the synchronous call-error path and on the revert-reason decoded
// from a failed receipt. Kept as a constant so a single edit covers both call
// sites if the on-chain message ever changes.
const outdatedKeyAgreementRevertReason = "Block number is lower than the latest key agreement"

type ParticipantStorageClientError struct {
	msg string
	err error
}

func NewParticipantStorageClientError(msg string) *ParticipantStorageClientError {
	return &ParticipantStorageClientError{
		msg: msg,
	}
}

func WrapInParticipantStorageClientError(msg string, err error) *ParticipantStorageClientError {
	return &ParticipantStorageClientError{
		msg: msg,
		err: err,
	}
}

func (e *ParticipantStorageClientError) Error() string {
	switch e.err {
	case nil:
		return e.msg
	default:
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	}
}

func (e *ParticipantStorageClientError) Unwrap() error {
	return e.err
}
