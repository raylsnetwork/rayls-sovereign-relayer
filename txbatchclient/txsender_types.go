package txbatchclient

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type TxSenderError struct {
	errs []error
}

func NewTxSenderError(err ...error) *TxSenderError {
	return &TxSenderError{
		errs: err,
	}
}

func (e *TxSenderError) Error() string {
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

func (e *TxSenderError) Unwrap() []error {
	return e.errs
}

type SendResult struct {
	Hash    string
	Receipt *types.Receipt
	Error   error
}

type SendAsyncResult struct {
	Hash  string
	Error error
}
