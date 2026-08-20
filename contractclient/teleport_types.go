package contractclient

import (
	"fmt"
	"math/big"
)

// Deprecated: Decommissioning Teleport (vanilla, atomic).
var ErrUnknownStatus = NewTeleportClientError("got and unknown status for shared ID")

type TeleportClientError struct {
	msg string
	err error
}

func NewTeleportClientError(msg string) *TeleportClientError {
	return &TeleportClientError{
		msg: msg,
	}
}

func WrapInTeleportClientError(msg string, err error) *TeleportClientError {
	return &TeleportClientError{
		msg: msg,
		err: err,
	}
}

func (e *TeleportClientError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	} else {
		return e.msg
	}
}

func (e *TeleportClientError) Unwrap() error {
	return e.err
}

type EncryptedMessages struct {
	BlockNumber *big.Int
	MessageTag  string
	Data        []byte
}
