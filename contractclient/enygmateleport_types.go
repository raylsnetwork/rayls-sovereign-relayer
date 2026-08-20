package contractclient

import "fmt"

type EnygmaTeleportClientError struct {
	msg string
	err error
}

func NewEnygmaTeleportClientError(msg string) *EnygmaTeleportClientError {
	return &EnygmaTeleportClientError{
		msg: msg,
	}
}

func WrapInEnygmaTeleportClientError(msg string, err error) *EnygmaTeleportClientError {
	return &EnygmaTeleportClientError{
		msg: msg,
		err: err,
	}
}

func (e *EnygmaTeleportClientError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	} else {
		return e.msg
	}
}

func (e *EnygmaTeleportClientError) Unwrap() error {
	return e.err
}
