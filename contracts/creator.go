package contracts

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

type ContractConstructor[T any] func(common.Address, bind.ContractBackend) (*T, error)

func CreateContract[T any](address common.Address, client bind.ContractBackend, constructor ContractConstructor[T]) (*T, error) {
	ctxCode, cancelCode := context.WithTimeout(context.Background(), 30*time.Second)
	code, err := client.CodeAt(ctxCode, address, nil)
	cancelCode()
	if err != nil {
		return nil, err
	}
	if len(code) == 0 {
		return nil, fmt.Errorf("no contract code at address %s", address.Hex())
	}

	instance, err := constructor(address, client)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("cannot create instance at %s: %w", address.String(), err))
	}

	return instance, nil
}
