package contractclient

import (
	"fmt"
)

type DvpERC721HandlerClientError struct {
	msg string
	err error
}

func NewDvpERC721HandlerClientError(msg string) *DvpERC721HandlerClientError {
	return &DvpERC721HandlerClientError{
		msg: msg,
	}
}

func WrapInDvpERC721HandlerClientError(msg string, err error) *DvpERC721HandlerClientError {
	return &DvpERC721HandlerClientError{
		msg: msg,
		err: err,
	}
}

func (e *DvpERC721HandlerClientError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	} else {
		return e.msg
	}
}

func (e *DvpERC721HandlerClientError) Unwrap() error {
	return e.err
}
