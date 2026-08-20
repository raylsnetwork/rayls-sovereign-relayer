package txbatchclient

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type TxReceipterError struct {
	errs []error
}

func NewTxReceipterError(err ...error) *TxReceipterError {
	return &TxReceipterError{
		errs: err,
	}
}

func (e *TxReceipterError) Error() string {
	if len(e.errs) == 1 {
		return e.errs[0].Error()
	}

	b := []byte(e.errs[0].Error())
	for _, err := range e.errs[1:] {
		b = append(b, '\n')
		b = append(b, err.Error()...)
	}
	return string(b)
}

func (e *TxReceipterError) Unwrap() []error {
	return e.errs
}

type ReceiptResult struct {
	Receipt *types.Receipt
	Error   error
}
