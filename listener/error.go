package listener

import (
	"errors"
	"fmt"
)

var (
	ErrLastProcessedBlockAlreadyExists = errors.New("last processed block for chain ID already exists")
	ErrLastProcessedBlockNotFound      = errors.New("last processed block for chain ID doesn't exist")
)

type ListenerError struct {
	msg string
	err error
}

func (e *ListenerError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err)
}

func (e *ListenerError) Unwrap() error {
	return e.err
}

func WrapInListenerError(msg string, err error) *ListenerError {
	return &ListenerError{
		msg: msg,
		err: err,
	}
}
