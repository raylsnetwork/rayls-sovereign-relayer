package contractclient

import "fmt"

type ResourceRegistryClientError struct {
	msg string
	err error
}

func NewResourceRegistryClientError(msg string) *ResourceRegistryClientError {
	return &ResourceRegistryClientError{
		msg: msg,
	}
}

func WrapInResourceRegistryClientError(msg string, err error) *ResourceRegistryClientError {
	return &ResourceRegistryClientError{
		msg: msg,
		err: err,
	}
}

func (e *ResourceRegistryClientError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	} else {
		return e.msg
	}
}

func (e *ResourceRegistryClientError) Unwrap() error {
	return e.err
}
