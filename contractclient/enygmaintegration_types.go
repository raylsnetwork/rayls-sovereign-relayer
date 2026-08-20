package contractclient

import (
	"fmt"
)

type EnygmaDvpIntegrationClientError struct {
	msg string
	err error
}

func NewEnygmaDvpIntegrationClientError(msg string) *EnygmaDvpIntegrationClientError {
	return &EnygmaDvpIntegrationClientError{
		msg: msg,
	}
}

func WrapInEnygmaDvpIntegrationClientError(msg string, err error) *EnygmaDvpIntegrationClientError {
	return &EnygmaDvpIntegrationClientError{
		msg: msg,
		err: err,
	}
}

func (e *EnygmaDvpIntegrationClientError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	} else {
		return e.msg
	}
}

func (e *EnygmaDvpIntegrationClientError) Unwrap() error {
	return e.err
}
