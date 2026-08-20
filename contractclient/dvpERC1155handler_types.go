package contractclient

import (
	"fmt"
)

type DvpERC1155HandlerClientError struct {
	msg string
	err error
}

func NewDvpERC1155HandlerClientError(msg string) *DvpERC1155HandlerClientError {
	return &DvpERC1155HandlerClientError{
		msg: msg,
	}
}

func WrapInDvpERC1155HandlerClientError(msg string, err error) *DvpERC1155HandlerClientError {
	return &DvpERC1155HandlerClientError{
		msg: msg,
		err: err,
	}
}

func (e *DvpERC1155HandlerClientError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	} else {
		return e.msg
	}
}

func (e *DvpERC1155HandlerClientError) Unwrap() error {
	return e.err
}
