package contractclient_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StubDvpERC1155Encryptor struct {
	encrData        []byte
	err             error
	spyRawDataBytes []byte
}

func (s *StubDvpERC1155Encryptor) EncryptDvpExtraData(_ context.Context, rawDataBytes []byte) ([]byte, error) {
	s.spyRawDataBytes = rawDataBytes
	if s.encrData != nil {
		return s.encrData, s.err
	}
	// Return valid JSON representing a slice with one struct
	return []byte(`[{"a":"test"}]`), s.err
}

func TestDvpERC1155Client_Burn(t *testing.T) {
	t.Run("successfully burns ERC1155", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		encr := &StubDvpERC1155Encryptor{}

		client := contractclient.NewDvpERC1155Client(executor, encr)

		tokenOwner := common.HexToAddress("0x1234567890123456789012345678901234567890")
		tokenId := big.NewInt(12345)
		tokenAmount := big.NewInt(100)
		err := client.Burn(context.Background(), "test-event-id", address, tokenOwner, tokenId, tokenAmount)

		require.NoError(t, err)
		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotNil(t, executor.spyExecuteCalldata)
	})

	t.Run("wraps executor errors in DvpERC1155ClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{executeErr: wantError}
		encr := &StubDvpERC1155Encryptor{}

		client := contractclient.NewDvpERC1155Client(executor, encr)

		err := client.Burn(context.Background(), "test-event-id", address, common.Address{}, big.NewInt(12345), big.NewInt(100))

		require.Error(t, err)
		var clientErr *contractclient.DvpERC1155ClientError
		require.True(t, errors.As(err, &clientErr))
		assert.Contains(t, err.Error(), "failed to burn ERC1155")
	})
}

func TestDvpERC1155Client_Approve(t *testing.T) {
	t.Run("successfully approves ERC1155 with recipient address", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		encr := &StubDvpERC1155Encryptor{}

		client := contractclient.NewDvpERC1155Client(executor, encr)

		to := common.HexToAddress("0x1234567890123456789012345678901234567890")
		err := client.Approve(context.Background(), "test-event-id", address, to)

		require.NoError(t, err)
		assert.Equal(t, address, executor.spyExecuteAddress)
	})

	t.Run("wraps executor errors in DvpERC1155ClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{executeErr: wantError}
		encr := &StubDvpERC1155Encryptor{}

		client := contractclient.NewDvpERC1155Client(executor, encr)

		err := client.Approve(context.Background(), "test-event-id", address, common.Address{})

		require.Error(t, err)
		var clientErr *contractclient.DvpERC1155ClientError
		require.True(t, errors.As(err, &clientErr))
		assert.Contains(t, err.Error(), "failed to approve ERC1155")
	})
}

func TestDvpERC1155Client_UpdateExtraData(t *testing.T) {
	t.Run("successfully updates extra data with encryption and unmarshal", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		encr := &StubDvpERC1155Encryptor{}

		client := contractclient.NewDvpERC1155Client(executor, encr)

		tokenId := big.NewInt(12345)
		tokenAmount := big.NewInt(100)
		chainId := big.NewInt(1)
		extraDataBytes := []byte("test extra data")
		newOwner := common.HexToAddress("0x1234567890123456789012345678901234567890")

		err := client.UpdateExtraData(context.Background(), "test-event-id", address, tokenId, tokenAmount, chainId, extraDataBytes, newOwner)

		require.NoError(t, err)
		require.Equal(t, extraDataBytes, encr.spyRawDataBytes)
		assert.Equal(t, address, executor.spyExecuteAddress)
	})

	t.Run("wraps encryption errors in DvpERC1155ClientError", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		encr := &StubDvpERC1155Encryptor{
			err: errors.New("encryption failed"),
		}

		client := contractclient.NewDvpERC1155Client(executor, encr)

		err := client.UpdateExtraData(
			context.Background(), "test-event-id",
			address,
			big.NewInt(12345),
			big.NewInt(100),
			big.NewInt(1),
			[]byte("test"),
			common.Address{},
		)

		require.Error(t, err)
		var clientErr *contractclient.DvpERC1155ClientError
		require.True(t, errors.As(err, &clientErr))
		assert.Contains(t, err.Error(), "failed to encrypt dvp ERC1155 extra data")
	})

	t.Run("wraps unmarshal json errors in DvpERC1155ClientError", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		encr := &StubDvpERC1155Encryptor{
			encrData: []byte("invalid json"),
		}

		client := contractclient.NewDvpERC1155Client(executor, encr)

		err := client.UpdateExtraData(
			context.Background(), "test-event-id",
			address,
			big.NewInt(12345),
			big.NewInt(100),
			big.NewInt(1),
			[]byte("test"),
			common.Address{},
		)

		require.Error(t, err)
		var clientErr *contractclient.DvpERC1155ClientError
		require.True(t, errors.As(err, &clientErr))
		assert.Contains(t, err.Error(), "failed to unmarshal extra data")
	})

	t.Run("wraps executor errors in DvpERC1155ClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{executeErr: wantError}
		encr := &StubDvpERC1155Encryptor{}

		client := contractclient.NewDvpERC1155Client(executor, encr)

		err := client.UpdateExtraData(
			context.Background(), "test-event-id",
			address,
			big.NewInt(12345),
			big.NewInt(100),
			big.NewInt(1),
			[]byte("test"),
			common.Address{},
		)

		require.Error(t, err)
		var clientErr *contractclient.DvpERC1155ClientError
		require.True(t, errors.As(err, &clientErr))
		assert.Contains(t, err.Error(), "failed to update ERC1155 extra data")
	})
}
