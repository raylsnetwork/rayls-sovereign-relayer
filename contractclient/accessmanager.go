package contractclient

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/RaylsAccessManagerV1"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

const (
	defaultPollInterval = 5 * time.Second
	// RoleRelayer is the access manager role name checked during authorization polling.
	RoleRelayer          = "RELAYER"
	RoleAuthorizedSender = "AUTHORIZED_SENDER"
)

// AccessManager polls the RaylsAccessManagerV1 contract to verify
// that relayer addresses have been granted the RELAYER.
// The struct name is kept for minimal caller disruption.
type AccessManager struct {
	address      common.Address
	contract     *RaylsAccessManagerV1.RaylsAccessManagerV1
	executor     Executor
	PollInterval time.Duration
}

func NewAccessManager(
	address common.Address,
	executor Executor,
) *AccessManager {
	return &AccessManager{
		address:      address,
		contract:     RaylsAccessManagerV1.NewRaylsAccessManagerV1(),
		executor:     executor,
		PollInterval: defaultPollInterval,
	}
}

// GrantAuthorizedSenderRole grants the AUTHORIZED_SENDER to the given
// token address via the access manager, replacing the old AddAuthorizedSender
// call on the endpoint contract.
func (c *AccessManager) GrantAuthorizedSenderRole(ctx context.Context, tokenAddress common.Address) error {
	// Look up the role ID for AUTHORIZED_SENDER
	calldata := c.contract.PackGetRoleIdByName(RoleAuthorizedSender)
	packedRoleID, err := c.executor.Call(ctx, c.address, calldata)
	if err != nil {
		return fmt.Errorf("failed to get AUTHORIZED_SENDER id: %w", err)
	}

	roleID, err := c.contract.UnpackGetRoleIdByName(packedRoleID)
	if err != nil {
		return fmt.Errorf("failed to unpack role ID: %w", err)
	}

	slog.Info(
		"Submitted GrantRole(AUTHORIZED_SENDER) transaction",
		slog.String("tokenAddress", tokenAddress.Hex()),
	)
	calldata = c.contract.PackGrantRole(roleID, tokenAddress, 0)
	receipt, err := c.executor.Execute(ctx, IDFor("accessmanager.GrantAuthorizedSenderRole", tokenAddress.Hex()), calldata, c.address)
	if err != nil {
		return fmt.Errorf("failed to grant AUTHORIZED_SENDER: %w", err)
	}

	slog.Info(
		"GrantRole(AUTHORIZED_SENDER) transaction confirmed",
		slog.String("txHash", receipt.TxHash.Hex()),
		slog.String("tokenAddress", tokenAddress.Hex()),
	)

	return nil
}

// WaitForAuthorization polls the access manager until all provided private keys
// have the RELAYER. It returns when all addresses derived from the keys
// have the role, or when the context is cancelled.
func (c *AccessManager) WaitForAuthorization(
	ctx context.Context,
	component string,
	keys []string,
) error {
	authLog := slog.With(slog.String("component", component))

	addresses, err := privateKeysToAddresses(keys)
	if err != nil {
		return WrapInAccessManagerError("failed to convert private keys to addresses", err)
	}

WAIT:
	for {
		for _, address := range addresses {
			calldata := c.contract.PackHasRoleByName(RoleRelayer, address)

			data, err := c.executor.Call(ctx, c.address, calldata)
			if err != nil {
				return WrapInAccessManagerError(
					"failed to check RELAYER for address",
					withstack.Wrap(err),
				)
			}

			hasRole, err := c.contract.UnpackHasRoleByName(data)
			if err != nil {
				return WrapInAccessManagerError(
					"failed to unpack calldata",
					withstack.Wrap(err),
				)
			}

			if !hasRole {
				authLog.Warn("Address does not have RELAYER", slog.Any("address", address))
				authLog.Info("Waiting for addresses to be authorized...")
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(c.PollInterval):
				}
				continue WAIT
			}
		}
		return nil
	}
}

func privateKeysToAddresses(privateKeysString []string) ([]common.Address, error) {
	var addresses []common.Address
	for _, marshalledKey := range privateKeysString {
		privateKey, err := crypto.HexToECDSA(marshalledKey)
		if err != nil {
			return nil, withstack.Wrap(fmt.Errorf("failed to unmarshal private key: %w", err))
		}
		addresses = append(addresses, crypto.PubkeyToAddress(privateKey.PublicKey))
	}
	return addresses, nil
}
