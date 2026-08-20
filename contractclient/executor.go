package contractclient

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
)

type ErrorWithRevertData struct {
	revertData []byte
}

func NewErrorWithRevertData(data []byte) *ErrorWithRevertData {
	return &ErrorWithRevertData{revertData: data}
}

func (e ErrorWithRevertData) Error() string {
	if len(e.revertData) > 0 {
		// Decode the standard Error(string) revert reason (selector 0x08c379a0) into the message so
		// downstream consumers that match on the human-readable reason see it, not just raw hex. In
		// particular the Enygma retry classifier (RetryService.isRetryableError) lists reasons like
		// "Invalid public signal for balance" and "Nullifier already used in pending transaction" as
		// retryable; with a hex-only Error() those substring matches never fire, so a transient
		// stale-proof race (the on-chain finalised balance moved between proof-gen and submission)
		// is misclassified non-retryable and the cross-transfer is dropped instead of regenerated.
		// The hex is retained for debugging and for custom-error reverts that don't decode to a string.
		if reason, derr := abi.UnpackRevert(e.revertData); derr == nil && reason != "" {
			return fmt.Sprintf("transaction reverted: %s (0x%s)", reason, hex.EncodeToString(e.revertData))
		}
		return fmt.Sprintf("transaction reverted: 0x%s", hex.EncodeToString(e.revertData))
	}
	return "transaction reverted"
}

func (e ErrorWithRevertData) GetRevertData() []byte {
	return e.revertData
}

type Client interface {
	bind.ContractBackend
	bind.DeployBackend
}

type Executor interface {
	Execute(ctx context.Context, id string, calldata []byte, address common.Address) (*types.Receipt, error)
	Call(ctx context.Context, address common.Address, calldata []byte) ([]byte, error)
}

type BatchExecutor interface {
	Executor
	BatchExecute(ctx context.Context, items []BatchInput) (map[string]BatchResult, error)
}

type Signer interface {
	Sign(ctx context.Context, calldata []byte, address common.Address) (*ethtypes.Transaction, error)
}

type LocalExecutor struct {
	gen    authGen
	queue  keyQueue
	client Client
}

func NewLocalExecutor(gen authGen, queue keyQueue, client Client) *LocalExecutor {
	return &LocalExecutor{
		gen:    gen,
		queue:  queue,
		client: client,
	}
}

// The second parameter is an id used for idempotency in the CTS.
// Since the local exeuctor doesn't have this feature it is ommited
func (e *LocalExecutor) Execute(ctx context.Context, _ string, calldata []byte, address common.Address) (*types.Receipt, error) {
	key, err := e.queue.Dequeue(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to dequeue key: %w", err)
	}
	defer e.queue.Enqueue(key)

	auth, err := e.gen.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate transaction auth: %w", err)
	}

	bound := bind.NewBoundContract(
		address,
		abi.ABI{},
		e.client,
		e.client,
		e.client,
	)

	tx, err := bound.RawTransact(auth, calldata)
	if err != nil {
		return nil, fmt.Errorf("send tx: %w", err)
	}

	receipt, err := bind.WaitMined(ctx, e.client, tx)
	if err != nil {
		return nil, fmt.Errorf("wait mined: %w", err)
	}
	if receipt.Status == types.ReceiptStatusFailed {
		data, err := e.tryGetRevertReason(ctx, auth.From, tx, receipt.BlockNumber)
		if err != nil {
			return receipt, fmt.Errorf("tx reverted, failed to get revert reason: %w", err)
		}
		return receipt, &ErrorWithRevertData{revertData: data}
	}

	return receipt, nil
}

func (e *LocalExecutor) Sign(ctx context.Context, calldata []byte, address common.Address) (*types.Transaction, error) {
	key, err := e.queue.Dequeue(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to dequeue key: %w", err)
	}
	defer e.queue.Enqueue(key)

	auth, err := e.gen.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate transaction auth: %w", err)
	}
	auth.NoSend = true

	bound := bind.NewBoundContract(
		address,
		abi.ABI{},
		e.client,
		e.client,
		e.client,
	)

	tx, err := bound.RawTransact(auth, calldata)
	if err != nil {
		return nil, fmt.Errorf("sign tx: %w", err)
	}

	return tx, nil
}

func (e *LocalExecutor) Call(ctx context.Context, address common.Address, calldata []byte) ([]byte, error) {
	key, err := e.queue.Dequeue(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to dequeue key for call: %w", err)
	}
	defer e.queue.Enqueue(key)

	from := crypto.PubkeyToAddress(key.PublicKey)

	msg := ethereum.CallMsg{
		From: from,
		To:   &address,
		Data: calldata,
	}

	out, err := e.client.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("call contract: %w", err)
	}

	return out, nil
}

func (e *LocalExecutor) tryGetRevertReason(ctx context.Context, from common.Address, tx *types.Transaction, blockNumber *big.Int) ([]byte, error) {
	msg := ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}

	_, err := e.client.CallContract(ctx, msg, blockNumber)
	// in case we get no error, this means we didn't revert
	if err == nil {
		// should be impossible to reach this, as it would signal
		// a state drift between mine-time and simulate-time
		return nil, nil
	}

	var de rpc.DataError
	if errors.As(err, &de) {
		// error is caused by EVM revert, so extract error data
		dataStr, ok := de.ErrorData().(string)
		if !ok {
			return nil, fmt.Errorf("unexpected revert data type %T", de.ErrorData())
		}
		data, err := hex.DecodeString(strings.TrimPrefix(dataStr, "0x"))
		if err != nil {
			return nil, fmt.Errorf("decode revert data: %w", err)
		}
		return data, nil
	}
	// error was caused by the rpc itself
	return nil, err
}

type UnpackFunc func(raw []byte) (any, error)

func WithVanillaRevert(data []byte, unpack UnpackFunc) (any, error) {
	revertReason, err := unpack(data)
	if err == nil {
		return revertReason, nil
	}

	return abi.UnpackRevert(data)
}

// IsRevertWithSelector reports whether err carries an *ErrorWithRevertData
// whose first 4 bytes match the 4-byte selector of the given custom error.
//
// errorID is the full Keccak-256 hash returned by the generated
// `<Contract><Error>ErrorID()` constants in the contract bindings; only the
// first 4 bytes are compared. Returns false when err does not unwrap to an
// *ErrorWithRevertData, or when the revert payload is shorter than 4 bytes.
func IsRevertWithSelector(err error, errorID common.Hash) bool {
	var rd *ErrorWithRevertData
	if !errors.As(err, &rd) {
		return false
	}
	data := rd.GetRevertData()
	if len(data) < 4 {
		return false
	}
	return bytes.Equal(data[:4], errorID.Bytes()[:4])
}
