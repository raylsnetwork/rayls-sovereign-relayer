package client

import "fmt"

type ProofAPIError struct {
	msg string
	err error
}

func NewProofAPIError(msg string) *ProofAPIError {
	return &ProofAPIError{
		msg: msg,
	}
}

func (p *ProofAPIError) Error() string {
	if p.err != nil {
		return fmt.Sprintf("%s: %s", p.msg, p.err.Error())
	}
	return p.msg
}

func (p *ProofAPIError) Unwrap() error {
	return p.err
}
