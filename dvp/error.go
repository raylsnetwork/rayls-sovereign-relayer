package dvp

import "fmt"

var (
	// ErrSwapAlreadyInitiated maps to the Dvp__SwapAlreadyExists() custom error
	// returned by the Dvp contract's initiateSwap when a swap with the same shared
	// id was already registered (race condition with the other side).
	ErrSwapAlreadyInitiated = NewDvpServiceError("swap already initiated")
	ErrDvpSwapNotFound      = NewDvpServiceError("dvp swap not found")
	ErrDvpSwapNotPending    = NewDvpServiceError("dvp swap not pending")
	ErrDvpSwapExpired       = NewDvpServiceError("dvp swap expired")
	ErrDvpInvalidVaultId    = NewDvpServiceError("dvp invalid vault id")
	ErrDvpInvalidProof      = NewDvpServiceError("dvp invalid proof")
)

type DvpServiceError struct {
	msg string
	err error
}

func NewDvpServiceError(msg string) *DvpServiceError {
	return &DvpServiceError{
		msg: msg,
	}
}

func WrapInDvpServiceError(msg string, err error) *DvpServiceError {
	return &DvpServiceError{
		msg: msg,
		err: err,
	}
}

func (e *DvpServiceError) Error() string {
	if e.err == nil {
		return e.msg
	} else {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	}
}

func (e *DvpServiceError) Unwrap() error {
	return e.err
}
