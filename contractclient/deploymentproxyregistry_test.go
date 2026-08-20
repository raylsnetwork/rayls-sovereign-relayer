package contractclient_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contractclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubDeploymentProxyBackend implements bind.ContractBackend but only
// CallContract is meaningful for the deployment proxy registry tests.
type stubDeploymentProxyBackend struct {
	callResult []byte
	callErr    error
}

func (s *stubDeploymentProxyBackend) CallContract(_ context.Context, _ ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	return s.callResult, s.callErr
}

// The remaining methods satisfy bind.ContractBackend but are unused in tests.
func (s *stubDeploymentProxyBackend) CodeAt(_ context.Context, _ common.Address, _ *big.Int) ([]byte, error) {
	panic("unexpected CodeAt")
}
func (s *stubDeploymentProxyBackend) HeaderByNumber(_ context.Context, _ *big.Int) (*ethTypes.Header, error) {
	panic("unexpected HeaderByNumber")
}
func (s *stubDeploymentProxyBackend) PendingCodeAt(_ context.Context, _ common.Address) ([]byte, error) {
	panic("unexpected PendingCodeAt")
}
func (s *stubDeploymentProxyBackend) PendingNonceAt(_ context.Context, _ common.Address) (uint64, error) {
	panic("unexpected PendingNonceAt")
}
func (s *stubDeploymentProxyBackend) SuggestGasPrice(_ context.Context) (*big.Int, error) {
	panic("unexpected SuggestGasPrice")
}
func (s *stubDeploymentProxyBackend) SuggestGasTipCap(_ context.Context) (*big.Int, error) {
	panic("unexpected SuggestGasTipCap")
}
func (s *stubDeploymentProxyBackend) EstimateGas(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
	panic("unexpected EstimateGas")
}
func (s *stubDeploymentProxyBackend) SendTransaction(_ context.Context, _ *ethTypes.Transaction) error {
	panic("unexpected SendTransaction")
}
func (s *stubDeploymentProxyBackend) FilterLogs(_ context.Context, _ ethereum.FilterQuery) ([]ethTypes.Log, error) {
	panic("unexpected FilterLogs")
}
func (s *stubDeploymentProxyBackend) SubscribeFilterLogs(_ context.Context, _ ethereum.FilterQuery, _ chan<- ethTypes.Log) (ethereum.Subscription, error) {
	panic("unexpected SubscribeFilterLogs")
}

// Ensure the stub satisfies the interface at compile time.
var _ bind.ContractBackend = (*stubDeploymentProxyBackend)(nil)

func TestNewDeploymentProxyRegistryClient(t *testing.T) {
	t.Run("returns error when backend CallContract fails", func(t *testing.T) {
		wantError := assert.AnError

		backend := &stubDeploymentProxyBackend{
			callErr: wantError,
		}

		client, err := contractclient.NewDeploymentProxyRegistryClient(common.Address{}, backend)
		require.Error(t, err)
		require.Nil(t, client)
		require.ErrorIs(t, err, wantError)
	})
}
