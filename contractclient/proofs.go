package contractclient

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/Proofs"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

// SubmitResult contains the result of submitting header proofs
type SubmitResult struct {
	StartBlock         *big.Int
	EndBlock           *big.Int
	NextExpectedBlock  *big.Int
	IncorrectHashEvent *Proofs.ProofsIncorrectParentHashEvent
}

// ProofsClient wraps the Proofs contract with transaction management
type ProofsClient struct {
	address  common.Address
	contract *Proofs.Proofs
	executor Executor
}

func NewProofsClient(
	address common.Address,
	executor Executor,
) *ProofsClient {
	return &ProofsClient{
		address:  address,
		contract: Proofs.NewProofs(),
		executor: executor,
	}
}

func (c *ProofsClient) GetNextBlockNumber(ctx context.Context, chainID *big.Int) (*big.Int, error) {
	calldata := c.contract.PackGetNextHeaderBlockNumber(chainID)

	raw, err := c.executor.Call(ctx, c.address, calldata)
	if err != nil {
		return nil, WrapInProofsClientError("failed to get next header block number", withstack.Wrap(err))
	}

	blockNumber, err := c.contract.UnpackGetNextHeaderBlockNumber(raw)
	if err != nil {
		return nil, WrapInProofsClientError("failed to unpack next header block number", withstack.Wrap(err))
	}
	return blockNumber, nil
}

func (c *ProofsClient) SubmitBatchHeaders(
	ctx context.Context,
	chainID *big.Int,
	headers []Proofs.ProofsHeader,
) (*SubmitResult, error) {
	if len(headers) == 0 {
		return nil, NewProofsClientError("no headers to submit")
	}

	calldata := c.contract.PackAddBatchHeaders(chainID, headers)

	// The proofs contract contains deduplication logic,
	// so this should be alright as a first version.
	// TODO: A more stable way to derive an id
	receipt, err := c.executor.Execute(ctx, IDFor("proofs.SubmitBatchHeaders", uuid.New().String()), calldata, c.address)
	if err != nil {
		return nil, WrapInProofsClientError("failed to submit batch headers", withstack.Wrap(err))
	}

	// Build result
	startBlock := headers[0].Number
	endBlock := headers[len(headers)-1].Number
	result := &SubmitResult{
		StartBlock:        startBlock,
		EndBlock:          endBlock,
		NextExpectedBlock: new(big.Int).Add(endBlock, big.NewInt(1)),
	}

	// Check for IncorrectParentHashEvent
	for _, log := range receipt.Logs {
		if event, err := c.contract.UnpackIncorrectParentHashEventEvent(log); err == nil {
			result.IncorrectHashEvent = event
			result.NextExpectedBlock = event.BlockNumber
			break
		}
	}

	return result, nil
}

// ProofsClientError is the error type for ProofsClient operations
type ProofsClientError struct {
	msg string
	err error
}

func NewProofsClientError(msg string) *ProofsClientError {
	return &ProofsClientError{msg: msg}
}

func WrapInProofsClientError(msg string, err error) *ProofsClientError {
	return &ProofsClientError{msg: msg, err: err}
}

func (e *ProofsClientError) Error() string {
	if e.err == nil {
		return e.msg
	}
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

func (e *ProofsClientError) Unwrap() error {
	return e.err
}
