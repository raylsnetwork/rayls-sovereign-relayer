package contractclient

import (
	"fmt"
)

type EnygmaClientError struct {
	msg string
	err error
}

func NewEnygmaClientError(msg string) *EnygmaClientError {
	return &EnygmaClientError{
		msg: msg,
	}
}

func WrapInEnygmaClientError(msg string, err error) *EnygmaClientError {
	return &EnygmaClientError{
		msg: msg,
		err: err,
	}
}

func (e *EnygmaClientError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	} else {
		return e.msg
	}
}

func (e *EnygmaClientError) Unwrap() error {
	return e.err
}
