package cloud

import (
	"errors"
	"fmt"
)

var ErrNoKey = errors.New("key does not exist")

type AWSKMSClientError struct {
	msg string
	err error
}

func WrapInAWSKMSClientError(msg string, err error) error {
	return &AWSKMSClientError{
		msg: msg,
		err: err,
	}
}

func (e *AWSKMSClientError) Error() string {
	return fmt.Sprintf("%s: %e", e.msg, e.err)
}

func (e *AWSKMSClientError) Unwrap() error {
	return e.err
}
