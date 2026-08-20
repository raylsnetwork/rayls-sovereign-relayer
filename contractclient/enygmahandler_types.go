package contractclient

import (
	"fmt"
)

type EnygmaHandlerClientError struct {
	msg string
	err error
}

func NewEnygmaHandlerClientError(msg string) *EnygmaHandlerClientError {
	return &EnygmaHandlerClientError{
		msg: msg,
	}
}

func WrapInEnygmaHandlerClientError(msg string, err error) *EnygmaHandlerClientError {
	return &EnygmaHandlerClientError{
		msg: msg,
		err: err,
	}
}

func (e *EnygmaHandlerClientError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	} else {
		return e.msg
	}
}

func (e *EnygmaHandlerClientError) Unwrap() error {
	return e.err
}
