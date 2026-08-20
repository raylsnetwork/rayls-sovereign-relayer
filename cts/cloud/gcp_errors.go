package cloud

import "fmt"

type GCPKMSClientError struct {
	msg string
	err error
}

func WrapInGCPKMSClientError(msg string, err error) *GCPKMSClientError {
	return &GCPKMSClientError{
		msg: msg,
		err: err,
	}
}

func (e *GCPKMSClientError) Error() string {
	return fmt.Sprintf("%s: %e", e.msg, e.err)
}

func (e *GCPKMSClientError) Unwrap() error {
	return e.err
}
