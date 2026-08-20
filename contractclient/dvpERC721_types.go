package contractclient

import (
	"fmt"
)

type DvpERC721ClientError struct {
	msg string
	err error
}

func NewDvpERC721ClientError(msg string) *DvpERC721ClientError {
	return &DvpERC721ClientError{
		msg: msg,
	}
}

func WrapInDvpERC721ClientError(msg string, err error) *DvpERC721ClientError {
	return &DvpERC721ClientError{
		msg: msg,
		err: err,
	}
}

func (e *DvpERC721ClientError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	} else {
		return e.msg
	}
}

func (e *DvpERC721ClientError) Unwrap() error {
	return e.err
}
