package txsim

import "fmt"

var (
	ErrTranasctionNotReverted = NewTransactionSimulatorError("transaction is not reverted - cannot get revert reason")
	ErrUnknownInterface       = NewTransactionSimulatorError(
		"failed to cast error to rpc.DataError - cannot extract error data",
	)
	ErrUnkownDataType = NewTransactionSimulatorError(
		"failed to cast error data to string - cannot decode revert reason",
	)
	ErrUnknownError = NewTransactionSimulatorError("revert reason not defined in local contract ABIs")
	ErrInvalidData  = NewTransactionSimulatorError("invalid data for unpacking")
)

type TransactionSimulatorError struct {
	msg string
	err error
}

func NewTransactionSimulatorError(msg string) *TransactionSimulatorError {
	return &TransactionSimulatorError{
		msg: msg,
	}
}

func WrapInTransactionSimulatorError(msg string, err error) *TransactionSimulatorError {
	return &TransactionSimulatorError{
		msg: msg,
		err: err,
	}
}

func (e *TransactionSimulatorError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	} else {
		return e.msg
	}
}

func (e *TransactionSimulatorError) Unwrap() error {
	return e.err
}
