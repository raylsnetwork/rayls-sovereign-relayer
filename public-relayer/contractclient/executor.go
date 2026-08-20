// Decommissioning Teleport (vanilla, atomic).

// Package contractclient implements the legacy public-chain (RN) Teleport bridge relayer.
//
// Deprecated: Decommissioning Teleport (vanilla, atomic).
package contractclient

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Client interface {
	bind.ContractBackend
	bind.DeployBackend
}

type IExecutor interface {
	Execute(ctx context.Context, calldata []byte, address common.Address) (*types.Receipt, error)
	Sign(ctx context.Context, calldata []byte, address common.Address) (*types.Transaction, error)
	Call(ctx context.Context, address common.Address, calldata []byte) ([]byte, error)
}

type Executor struct {
	gen    authGen
	queue  keyQueue
	client Client
}

func NewExecutor(gen authGen, queue keyQueue, client Client) *Executor {
	return &Executor{
		gen:    gen,
		queue:  queue,
		client: client,
	}
}

func (e *Executor) Execute(ctx context.Context, calldata []byte, address common.Address) (*types.Receipt, error) {
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

	receipt, err := bind.WaitMined(ctx, e.client, tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("wait mined: %w", err)
	}

	return receipt, nil
}

func (e *Executor) Deploy(ctx context.Context, bytecode []byte, constructor []byte) (common.Address, *types.Receipt, error) {
	key, err := e.queue.Dequeue(ctx)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to dequeue key: %w", err)
	}
	defer e.queue.Enqueue(key)

	auth, err := e.gen.Get(ctx, key)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to generate transaction auth: %w", err)
	}

	address, tx, err := bind.DeployContract(auth, bytecode, e.client, constructor)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to send contract for deployment: %w", err)
	}

	receipt, err := bind.WaitMined(ctx, e.client, tx.Hash())
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("wait mined: %w", err)
	}

	return address, receipt, nil
}

func (e *Executor) Sign(ctx context.Context, calldata []byte, address common.Address) (*types.Transaction, error) {
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

func (e *Executor) Call(ctx context.Context, address common.Address, calldata []byte) ([]byte, error) {
	msg := ethereum.CallMsg{
		To:   &address,
		Data: calldata,
	}

	out, err := e.client.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("call contract: %w", err)
	}

	return out, nil
}
