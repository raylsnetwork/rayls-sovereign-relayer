// Decommissioning Teleport (vanilla, atomic).

package contractclient

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	sharedcontractclient "github.com/raylsnetwork/rayls-sovereign-relayer/contractclient"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/PNTokenRegistryV1"
)

// TokenGovernanceClient wraps the PNTokenRegistryV1 facade for the operations
// the public relayer needs: reading a token's registry entry (for deploy
// idempotency and metadata) and writing back the deployed public-chain mirror
// address. It replaces the legacy token-governance client.
type TokenGovernanceClient struct {
	address  common.Address
	contract *PNTokenRegistryV1.PNTokenRegistryV1
	executor sharedcontractclient.Executor
}

func NewTokenGovernanceClient(
	address common.Address,
	executor sharedcontractclient.Executor,
) *TokenGovernanceClient {
	return &TokenGovernanceClient{
		address:  address,
		contract: PNTokenRegistryV1.NewPNTokenRegistryV1(),
		executor: executor,
	}
}

// GetTokenByAddress reads the full registry entry for a token by its private
// (PN) address. Used to fetch deploy metadata (name/symbol/uri/standard) and to
// check the current public-chain status for idempotency.
func (c *TokenGovernanceClient) GetTokenByAddress(
	ctx context.Context,
	tokenAddress common.Address,
) (PNTokenRegistryV1.TokenStructsToken, error) {
	calldata := c.contract.PackGetTokenByAddress(tokenAddress)

	raw, err := c.executor.Call(ctx, c.address, calldata)
	if err != nil {
		return PNTokenRegistryV1.TokenStructsToken{}, fmt.Errorf("failed to get token by address: %w", err)
	}

	token, err := c.contract.UnpackGetTokenByAddress(raw)
	if err != nil {
		return PNTokenRegistryV1.TokenStructsToken{}, fmt.Errorf("failed to unpack token by address: %w", err)
	}

	return token, nil
}

// GetPublicAddressByPrivateAddress returns the deployed public-chain mirror
// address for a token, or the zero address if none is registered yet. Backed by
// GetTokenByAddress on the registry facade; kept for deploy idempotency checks.
func (c *TokenGovernanceClient) GetPublicAddressByPrivateAddress(
	ctx context.Context,
	privateAddress common.Address,
) (common.Address, error) {
	token, err := c.GetTokenByAddress(ctx, privateAddress)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get public address: %w", err)
	}

	return token.PublicTokenAddress, nil
}

func (c *TokenGovernanceClient) UpdatePublicTokenAddress(ctx context.Context, privateAddress, publicAddress common.Address) error {
	calldata := c.contract.PackUpdatePublicTokenAddress(privateAddress, publicAddress)

	receipt, err := c.executor.Execute(ctx, sharedcontractclient.IDFor("tokengovernance.UpdatePublicTokenAddress", privateAddress.Hex(), publicAddress.Hex()), calldata, c.address)
	if err != nil {
		return fmt.Errorf("failed to update public token address: %w", err)
	}
	if receipt.Status != 1 {
		return fmt.Errorf("transaction failed with receipt status: %d", receipt.Status)
	}

	slog.Debug("UpdatePublicTokenAddress", slog.Any("publicAddress", publicAddress))
	return nil
}
