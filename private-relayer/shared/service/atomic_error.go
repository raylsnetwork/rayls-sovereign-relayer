// Decommissioning Teleport (vanilla, atomic).

package service

import "fmt"

var (
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	ErrAlreadyReverted = NewAtomicServiceError("transactions already reverted in private hub")
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	ErrAlreadyExecuted = NewAtomicServiceError("transactions already executed in private hub")
)

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type AtomicServiceError struct {
	msg string
	err error
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func NewAtomicServiceError(msg string) *AtomicServiceError {
	return &AtomicServiceError{
		msg: msg,
	}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func WrapInAtomicServiceError(msg string, err error) *AtomicServiceError {
	return &AtomicServiceError{
		msg: msg,
		err: err,
	}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (e *AtomicServiceError) Error() string {
	if e.err == nil {
		return e.msg
	} else {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (e *AtomicServiceError) Unwrap() error {
	return e.err
}
