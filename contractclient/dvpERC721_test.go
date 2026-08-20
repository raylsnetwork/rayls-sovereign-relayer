package contractclient_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StubDvpERC721Encryptor struct {
	encrData        []byte
	err             error
	spyRawDataBytes []byte
}

func (s *StubDvpERC721Encryptor) EncryptDvpExtraData(_ context.Context, rawDataBytes []byte) ([]byte, error) {
	s.spyRawDataBytes = rawDataBytes
	if s.encrData != nil {
		return s.encrData, s.err
	}
	// Return valid JSON representing a slice with one struct
	return []byte(`[{"a":"test"}]`), s.err
}

func createMockTransaction() *ethTypes.Transaction {
	return ethTypes.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil)
}

func createMockReceipt() *ethTypes.Receipt {
	return &ethTypes.Receipt{
		Status: 1, // Success
		TxHash: common.Hash{},
	}
}

func TestDvpERC721Client_Burn(t *testing.T) {
	t.Run("successfully burns ERC721", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		encr := &StubDvpERC721Encryptor{}

		client := contractclient.NewDvpERC721Client(executor, encr)

		nftId := big.NewInt(12345)
		err := client.Burn(context.Background(), "test-event-id", address, nftId)

		require.NoError(t, err)
		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotNil(t, executor.spyExecuteCalldata)
	})

	t.Run("wraps executor errors in DvpERC721ClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{executeErr: wantError}
		encr := &StubDvpERC721Encryptor{}

		client := contractclient.NewDvpERC721Client(executor, encr)

		err := client.Burn(context.Background(), "test-event-id", address, big.NewInt(12345))

		require.Error(t, err)
		var clientErr *contractclient.DvpERC721ClientError
		require.True(t, errors.As(err, &clientErr))
		assert.Contains(t, err.Error(), "failed to burn ERC721")
	})
}

func TestDvpERC721Client_Approve(t *testing.T) {
	t.Run("successfully approves ERC721 with recipient address", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		encr := &StubDvpERC721Encryptor{}

		client := contractclient.NewDvpERC721Client(executor, encr)

		to := common.HexToAddress("0x1234567890123456789012345678901234567890")
		nftId := big.NewInt(12345)
		err := client.Approve(context.Background(), "test-event-id", address, to, nftId)

		require.NoError(t, err)
		assert.Equal(t, address, executor.spyExecuteAddress)
	})

	t.Run("wraps revert with unrecognised revert data in DvpERC721ClientError", func(t *testing.T) {
		revertErr := contractclient.NewErrorWithRevertData([]byte{0x01, 0x02, 0x03, 0x04})

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{executeErr: revertErr}
		encr := &StubDvpERC721Encryptor{}

		client := contractclient.NewDvpERC721Client(executor, encr)

		err := client.Approve(context.Background(), "test-event-id", address, common.Address{}, big.NewInt(12345))

		require.Error(t, err)
		var clientErr *contractclient.DvpERC721ClientError
		require.True(t, errors.As(err, &clientErr))
		require.ErrorIs(t, err, revertErr)
		assert.Contains(t, err.Error(), "failed to approve ERC721")
	})

	t.Run("wraps executor execute errors in DvpERC721ClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{executeErr: wantError}
		encr := &StubDvpERC721Encryptor{}

		client := contractclient.NewDvpERC721Client(executor, encr)

		err := client.Approve(context.Background(), "test-event-id", address, common.Address{}, big.NewInt(12345))

		require.Error(t, err)
		var clientErr *contractclient.DvpERC721ClientError
		require.True(t, errors.As(err, &clientErr))
		assert.Contains(t, err.Error(), "failed to approve ERC721")
	})
}

func TestDvpERC721Client_UpdateExtraData(t *testing.T) {
	t.Run("successfully updates extra data with encryption and unmarshal", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		encr := &StubDvpERC721Encryptor{}

		client := contractclient.NewDvpERC721Client(executor, encr)

		nftId := big.NewInt(12345)
		chainId := big.NewInt(1)
		extraDataBytes := []byte("test extra data")
		newOwner := common.HexToAddress("0x1234567890123456789012345678901234567890")

		err := client.UpdateExtraData(context.Background(), "test-event-id", address, nftId, chainId, extraDataBytes, newOwner)

		require.NoError(t, err)
		require.Equal(t, extraDataBytes, encr.spyRawDataBytes)
		assert.Equal(t, address, executor.spyExecuteAddress)
	})

	t.Run("wraps encryption errors in DvpERC721ClientError", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		encr := &StubDvpERC721Encryptor{
			err: errors.New("encryption failed"),
		}

		client := contractclient.NewDvpERC721Client(executor, encr)

		err := client.UpdateExtraData(
			context.Background(), "test-event-id",
			address,
			big.NewInt(12345),
			big.NewInt(1),
			[]byte("test"),
			common.Address{},
		)

		require.Error(t, err)
		var clientErr *contractclient.DvpERC721ClientError
		require.True(t, errors.As(err, &clientErr))
		assert.Contains(t, err.Error(), "failed to encrypt ERC721 extra data")
	})

	t.Run("wraps unmarshal json errors in DvpERC721ClientError", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		encr := &StubDvpERC721Encryptor{
			encrData: []byte("invalid json"),
		}

		client := contractclient.NewDvpERC721Client(executor, encr)

		err := client.UpdateExtraData(
			context.Background(), "test-event-id",
			address,
			big.NewInt(12345),
			big.NewInt(1),
			[]byte("test"),
			common.Address{},
		)

		require.Error(t, err)
		var clientErr *contractclient.DvpERC721ClientError
		require.True(t, errors.As(err, &clientErr))
		assert.Contains(t, err.Error(), "failed to unmarshal extra data")
	})

	t.Run("wraps executor errors in DvpERC721ClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{executeErr: wantError}
		encr := &StubDvpERC721Encryptor{}

		client := contractclient.NewDvpERC721Client(executor, encr)

		err := client.UpdateExtraData(
			context.Background(), "test-event-id",
			address,
			big.NewInt(12345),
			big.NewInt(1),
			[]byte("test"),
			common.Address{},
		)

		require.Error(t, err)
		var clientErr *contractclient.DvpERC721ClientError
		require.True(t, errors.As(err, &clientErr))
		assert.Contains(t, err.Error(), "failed to update ERC721 infos")
	})
}
