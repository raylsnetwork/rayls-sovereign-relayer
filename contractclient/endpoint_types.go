package contractclient

import "fmt"

type EndpointClientError struct {
	msg string
	err error
}

func NewEndpointClientError(msg string) *EndpointClientError {
	return &EndpointClientError{
		msg: msg,
	}
}

func WrapInEndpointClientError(msg string, err error) *EndpointClientError {
	return &EndpointClientError{
		msg: msg,
		err: err,
	}
}

func (e *EndpointClientError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	} else {
		return e.msg
	}
}

func (e *EndpointClientError) Unwrap() error {
	return e.err
}
