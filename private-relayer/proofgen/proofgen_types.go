package proofgen

import "fmt"

type ProofGeneratorError struct {
	msg string
	err error
}

func NewProofGeneratorError(msg string) *ProofGeneratorError {
	return &ProofGeneratorError{
		msg: msg,
	}
}

func WrapInProofGeneratorError(msg string, err error) *ProofGeneratorError {
	return &ProofGeneratorError{
		msg: msg,
		err: err,
	}
}

func (e *ProofGeneratorError) Error() string {
	switch e.err {
	case nil:
		return e.msg
	default:
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	}
}

func (e *ProofGeneratorError) Unwrap() error {
	return e.err
}
