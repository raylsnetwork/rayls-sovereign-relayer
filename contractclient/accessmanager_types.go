package contractclient

import (
	"fmt"
)

type AccessManagerError struct {
	msg string
	err error
}

func NewAccessManagerError(msg string) *AccessManagerError {
	return &AccessManagerError{
		msg: msg,
	}
}

func WrapInAccessManagerError(msg string, err error) *AccessManagerError {
	return &AccessManagerError{
		msg: msg,
		err: err,
	}
}

func (e *AccessManagerError) Error() string {
	switch e.err {
	case nil:
		return e.msg
	default:
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	}
}

func (e *AccessManagerError) Unwrap() error {
	return e.err
}
