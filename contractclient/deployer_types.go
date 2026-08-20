package contractclient

import "fmt"

type DeployerClientError struct {
	msg string
	err error
}

func NewDeployerClientError(msg string) *DeployerClientError {
	return &DeployerClientError{
		msg: msg,
	}
}

func WrapInDeployerClientError(msg string, err error) *DeployerClientError {
	return &DeployerClientError{
		msg: msg,
		err: err,
	}
}

func (e *DeployerClientError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	} else {
		return e.msg
	}
}

func (e *DeployerClientError) Unwrap() error {
	return e.err
}
