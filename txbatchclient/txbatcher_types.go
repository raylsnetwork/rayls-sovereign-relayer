package txbatchclient

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
)

type CallError struct {
	msg string
	err error
}

func NewCallError(err error) *CallError {
	return &CallError{
		msg: "error while trying to send batch",
		err: err,
	}
}

func (e *CallError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

func (e *CallError) Unwrap() error {
	return e.err
}

type RPCError struct {
	msg string
	err rpc.Error
}

func NewRPCError(err rpc.Error) *RPCError {
	return &RPCError{
		msg: "error while processing batch in PL",
		err: err,
	}
}

func (e *RPCError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

func (e *RPCError) Unwrap() error {
	return e.err
}
