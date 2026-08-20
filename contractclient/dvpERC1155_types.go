package contractclient

import (
	"fmt"
)

type DvpERC1155ClientError struct {
	msg string
	err error
}

func NewDvpERC1155ClientError(msg string) *DvpERC1155ClientError {
	return &DvpERC1155ClientError{
		msg: msg,
	}
}

func WrapInDvpERC1155ClientError(msg string, err error) *DvpERC1155ClientError {
	return &DvpERC1155ClientError{
		msg: msg,
		err: err,
	}
}

func (e *DvpERC1155ClientError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	} else {
		return e.msg
	}
}

func (e *DvpERC1155ClientError) Unwrap() error {
	return e.err
}
