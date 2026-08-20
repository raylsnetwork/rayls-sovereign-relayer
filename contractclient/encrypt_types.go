package contractclient

import "fmt"

type EncryptorError struct {
	msg string
	err error
}

func WrapInEncryptorError(msg string, err error) *EncryptorError {
	return &EncryptorError{
		msg: msg,
		err: err,
	}
}

func (e *EncryptorError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

func (e *EncryptorError) Unwrap() error {
	return e.err
}
