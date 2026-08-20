package contractclient

import (
	"fmt"
)

type DvpClientError struct {
	msg string
	err error
}

func NewDvpClientError(msg string) *DvpClientError {
	return &DvpClientError{
		msg: msg,
	}
}

func WrapInDvpClientError(msg string, err error) *DvpClientError {
	return &DvpClientError{
		msg: msg,
		err: err,
	}
}

func (e *DvpClientError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	} else {
		return e.msg
	}
}

func (e *DvpClientError) Unwrap() error {
	return e.err
}

// contract enum
type DvpSwapProofType uint8

const (
	DvpSwapPaymentProofType DvpSwapProofType = iota
	DvpSwapDeliveryProofType
)
