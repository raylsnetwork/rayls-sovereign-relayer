// Decommissioning Teleport (vanilla, atomic).

package contractclient_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/PNTokenRegistryV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/public-relayer/contractclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// packTokenByAddressResponse ABI-encodes a TokenStructsToken as the getTokenByAddress
// return value, so the stub executor can return a realistic Call response.
func packTokenByAddressResponse(t *testing.T, token PNTokenRegistryV1.TokenStructsToken) []byte {
	t.Helper()

	contractABI, err := PNTokenRegistryV1.PNTokenRegistryV1MetaData.ParseABI()
	require.NoError(t, err)

	if token.IssuerChainId == nil {
		token.IssuerChainId = big.NewInt(0)
	}
	if token.CreatedAt == nil {
		token.CreatedAt = big.NewInt(0)
	}
	if token.UpdatedAt == nil {
		token.UpdatedAt = big.NewInt(0)
	}

	resp, err := contractABI.Methods["getTokenByAddress"].Outputs.Pack(token)
	require.NoError(t, err)
	return resp
}

type tokenGovernanceStubExecutor struct {
	callAddress    common.Address
	callCalldata   []byte
	callResponse   []byte
	callErr        error
	execAddress    common.Address
	execCalldata   []byte
	execReceipt    *ethTypes.Receipt
	execErr        error
	signAddress    common.Address
	signCalldata   []byte
	signTx         *ethTypes.Transaction
	signErr        error
}

func (s *tokenGovernanceStubExecutor) Execute(_ context.Context, _ string, calldata []byte, address common.Address) (*ethTypes.Receipt, error) {
	s.execAddress = address
	s.execCalldata = calldata
	return s.execReceipt, s.execErr
}

func (s *tokenGovernanceStubExecutor) Sign(_ context.Context, calldata []byte, address common.Address) (*ethTypes.Transaction, error) {
	s.signAddress = address
	s.signCalldata = calldata
	return s.signTx, s.signErr
}

func (s *tokenGovernanceStubExecutor) Call(_ context.Context, address common.Address, calldata []byte) ([]byte, error) {
	s.callAddress = address
	s.callCalldata = calldata
	return s.callResponse, s.callErr
}

func TestTokenGovernanceClient_GetTokenByAddress(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("returns the registry token entry for a given address", func(t *testing.T) {
		wantAddr := common.HexToAddress("0x1111")
		wantPublic := common.HexToAddress("0x2222")

		binding := PNTokenRegistryV1.NewPNTokenRegistryV1()
		calldata := binding.PackGetTokenByAddress(wantAddr)
		resp := packTokenByAddressResponse(t, PNTokenRegistryV1.TokenStructsToken{
			TokenAddress:       wantAddr,
			PublicTokenAddress: wantPublic,
			Symbol:             "TT",
			ErcStandard:        1,
		})

		exec := &tokenGovernanceStubExecutor{callResponse: resp}
		address := common.HexToAddress("0xAAAA")
		client := contractclient.NewTokenGovernanceClient(address, exec)

		token, err := client.GetTokenByAddress(context.Background(), wantAddr)

		require.NoError(t, err)
		assert.Equal(t, wantAddr, token.TokenAddress)
		assert.Equal(t, wantPublic, token.PublicTokenAddress)
		assert.Equal(t, "TT", token.Symbol)
		assert.Equal(t, uint8(1), token.ErcStandard)
		assert.Equal(t, address, exec.callAddress)
		assert.Equal(t, calldata, exec.callCalldata)
	})

	t.Run("wraps contract error", func(t *testing.T) {
		wantErr := errors.New("contract read failed")
		exec := &tokenGovernanceStubExecutor{callErr: wantErr}
		client := contractclient.NewTokenGovernanceClient(common.Address{}, exec)

		_, err := client.GetTokenByAddress(context.Background(), common.Address{})

		require.ErrorIs(t, err, wantErr)
		assert.Contains(t, err.Error(), "failed to get token by address")
	})
}

func TestTokenGovernanceClient_GetPublicAddressByPrivateAddress(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("returns public address for given private address", func(t *testing.T) {
		wantPrivate := common.HexToAddress("0x1111")
		wantPublic := common.HexToAddress("0x2222")

		binding := PNTokenRegistryV1.NewPNTokenRegistryV1()
		calldata := binding.PackGetTokenByAddress(wantPrivate)
		resp := packTokenByAddressResponse(t, PNTokenRegistryV1.TokenStructsToken{
			TokenAddress:       wantPrivate,
			PublicTokenAddress: wantPublic,
		})

		exec := &tokenGovernanceStubExecutor{
			callResponse: resp,
		}
		address := common.HexToAddress("0xAAAA")
		client := contractclient.NewTokenGovernanceClient(address, exec)

		gotPublic, err := client.GetPublicAddressByPrivateAddress(context.Background(), wantPrivate)

		require.NoError(t, err)
		assert.Equal(t, wantPublic, gotPublic)
		assert.Equal(t, address, exec.callAddress)
		assert.Equal(t, calldata, exec.callCalldata)
	})

	t.Run("wraps contract error", func(t *testing.T) {
		wantErr := errors.New("contract read failed")
		exec := &tokenGovernanceStubExecutor{callErr: wantErr}
		client := contractclient.NewTokenGovernanceClient(common.Address{}, exec)

		_, err := client.GetPublicAddressByPrivateAddress(context.Background(), common.Address{})

		require.ErrorIs(t, err, wantErr)
		assert.Contains(t, err.Error(), "failed to get public address")
	})
}

func TestTokenGovernanceClient_UpdatePublicTokenAddress(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("successfully updates public token address", func(t *testing.T) {
		wantPrivate := common.HexToAddress("0xaaaa")
		wantPublic := common.HexToAddress("0xbbbb")
		address := common.HexToAddress("0xAAAA")

		binding := PNTokenRegistryV1.NewPNTokenRegistryV1()
		calldata := binding.PackUpdatePublicTokenAddress(wantPrivate, wantPublic)
		exec := &tokenGovernanceStubExecutor{
			execReceipt: &ethTypes.Receipt{Status: 1},
		}

		client := contractclient.NewTokenGovernanceClient(address, exec)

		err := client.UpdatePublicTokenAddress(context.Background(), wantPrivate, wantPublic)

		require.NoError(t, err)

		assert.Equal(t, address, exec.execAddress)
		assert.Equal(t, calldata, exec.execCalldata)
	})

	t.Run("wraps execute error", func(t *testing.T) {
		wantErr := errors.New("contract write failed")
		exec := &tokenGovernanceStubExecutor{execErr: wantErr}
		client := contractclient.NewTokenGovernanceClient(common.Address{}, exec)

		err := client.UpdatePublicTokenAddress(context.Background(), common.Address{}, common.Address{})

		require.ErrorIs(t, err, wantErr)
		assert.Contains(t, err.Error(), "failed to update public token address")
	})

	t.Run("returns error when receipt status is not 1", func(t *testing.T) {
		exec := &tokenGovernanceStubExecutor{
			execReceipt: &ethTypes.Receipt{Status: 0},
		}
		client := contractclient.NewTokenGovernanceClient(common.Address{}, exec)

		err := client.UpdatePublicTokenAddress(context.Background(), common.Address{}, common.Address{})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "transaction failed with receipt status")
	})
}
