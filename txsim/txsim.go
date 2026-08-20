package txsim

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// minFunctionSelectorLength is the minimum data length required to extract an ABI function selector (4 bytes).
const minFunctionSelectorLength = 4

type ethereumClient interface {
	CallContract(context.Context, ethereum.CallMsg, *big.Int) ([]byte, error)
	TransactionByHash(context.Context, common.Hash) (*types.Transaction, bool, error)
}

type TransactionSimulator struct {
	client ethereumClient
}

func NewTransactionSimulator(client ethereumClient) *TransactionSimulator {
	return &TransactionSimulator{
		client: client,
	}
}

func (s *TransactionSimulator) GetRevertReason(ctx context.Context, hash common.Hash) (ContractError, error) {
	tx, _, err := s.client.TransactionByHash(ctx, hash)
	if err != nil {
		return ContractError{}, WrapInTransactionSimulatorError("failed to get transaction by hash", err)
	}

	signer := types.NewEIP155Signer(tx.ChainId())
	from, err := types.Sender(signer, tx)
	if err != nil {
		return ContractError{}, WrapInTransactionSimulatorError("failed to recover sender", err)
	}

	msg := ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}

	_, txErr := s.client.CallContract(ctx, msg, nil)
	if txErr == nil {
		return ContractError{}, ErrTranasctionNotReverted
	}

	return s.DecodeRevertReason(ctx, txErr)
}

func (s *TransactionSimulator) DecodeRevertBytes(data []byte) (ContractError, error) {
	return DecodeRevertBytes(data)
}

// DecodeRevertBytes decodes raw EVM revert data into a ContractError. It resolves both the
// standard `Error(string)` revert and the custom errors registered in the package-level error map
// (see PopulateErrorMap). It relies only on that map and the ABI decoder — no RPC client — so it is
// safe to call from any context that has revert bytes in hand (e.g. an ErrorWithRevertData unwound
// from a failed cross-transfer batch), not just the simulator.
func DecodeRevertBytes(data []byte) (ContractError, error) {
	if len(data) < minFunctionSelectorLength {
		return ContractError{}, ErrInvalidData
	}

	revertReason, err := abi.UnpackRevert(data)
	if err == nil {
		return ContractError{
			Sig:  "Error(string)",
			Args: []interface{}{revertReason},
		}, nil
	}

	return unpackCustomRevert(data)
}

func (s *TransactionSimulator) DecodeRevertReason(ctx context.Context, txErr error) (ContractError, error) {
	var jsonError rpc.DataError
	if !errors.As(txErr, &jsonError) {
		return ContractError{}, ErrUnknownInterface
	}

	dataStr, ok := jsonError.ErrorData().(string)
	if !ok {
		return ContractError{}, ErrUnkownDataType
	}

	revertReason, err := abi.UnpackRevert(common.FromHex(dataStr))
	if err == nil {
		return ContractError{
			Sig:  "Error(string)",
			Args: []interface{}{revertReason},
		}, nil
	}

	contractError, err := unpackCustomRevert(common.FromHex(dataStr))
	if err != nil {
		return ContractError{}, err
	}

	return contractError, nil
}

func unpackCustomRevert(data []byte) (ContractError, error) {
	if len(data) < minFunctionSelectorLength {
		return ContractError{}, ErrInvalidData
	}

	def, exists := contractErrors[common.Bytes2Hex(data[:minFunctionSelectorLength])]
	if !exists {
		return ContractError{}, ErrUnknownError
	}

	args, err := def.Inputs.Unpack(data[minFunctionSelectorLength:])
	if err != nil {
		return ContractError{}, WrapInTransactionSimulatorError("failed to unpack error inputs", err)
	}

	return ContractError{
		Sig:  def.Sig,
		Args: args,
	}, nil
}
