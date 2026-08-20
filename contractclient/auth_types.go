package contractclient

import "fmt"

type AuthGenError struct {
	msg string
	err error
}

func WrapInAuthGenError(msg string, err error) *AuthGenError {
	return &AuthGenError{
		msg: msg,
		err: err,
	}
}

func (e *AuthGenError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

func (e *AuthGenError) Unwrap() error {
	return e.err
}
