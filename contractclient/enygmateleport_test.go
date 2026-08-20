package contractclient_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contractclient"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StubEnygmaTeleportEncryptor struct {
	encryptedData []byte
	err           error

	spyMessages []types.EnygmaTransferCompleted
}

func (s *StubEnygmaTeleportEncryptor) EncryptEnygmaTransferBatchCompleted(
	_ context.Context,
	messages []types.EnygmaTransferCompleted,
) ([]byte, error) {
	s.spyMessages = messages
	return s.encryptedData, s.err
}

type StubEnygmaTeleportExecutor struct {
	spyCalldata []byte
	spyAddress  common.Address
	receipt     *ethTypes.Receipt
	err         error
}

func (s *StubEnygmaTeleportExecutor) Execute(
	_ context.Context,
	_ string,
	calldata []byte,
	address common.Address,
) (*ethTypes.Receipt, error) {
	s.spyCalldata = calldata
	s.spyAddress = address
	return s.receipt, s.err
}

func (s *StubEnygmaTeleportExecutor) Sign(
	_ context.Context,
	_ []byte,
	_ common.Address,
) (*ethTypes.Transaction, error) {
	panic("unexpected Sign")
}

func (s *StubEnygmaTeleportExecutor) Call(
	_ context.Context,
	_ common.Address,
	_ []byte,
) ([]byte, error) {
	panic("unexpected Call")
}

func TestEnygmaTeleportClient_SendTransferCompleted(t *testing.T) {
	t.Run("successfully sends transfer completed with all dependencies working", func(t *testing.T) {
		wantMessages := []types.EnygmaTransferCompleted{{MessageId: "msg-1"}, {MessageId: "msg-2"}}
		wantEncryptedData := []byte{0xaa, 0xbb, 0xcc, 0xdd}

		address := common.HexToAddress("0x1234567890123456789012345678901234567890")

		encryptor := &StubEnygmaTeleportEncryptor{
			encryptedData: wantEncryptedData,
		}
		executor := &StubEnygmaTeleportExecutor{
			receipt: &ethTypes.Receipt{Status: 1},
		}

		client := contractclient.NewEnygmaTeleportClient(address, executor, encryptor)

		err := client.SendTransferCompleted(context.Background(), wantMessages)

		require.Nil(t, err)

		assert.Equal(t, wantMessages, encryptor.spyMessages)
		assert.NotNil(t, executor.spyCalldata)
		assert.Equal(t, address, executor.spyAddress)
	})

	t.Run("wraps encryption errors in EnygmaTeleportClientError", func(t *testing.T) {
		wantError := errors.New("encryption failed")

		address := common.Address{}
		executor := &StubEnygmaTeleportExecutor{}

		encryptor := &StubEnygmaTeleportEncryptor{
			err: wantError,
		}

		client := contractclient.NewEnygmaTeleportClient(address, executor, encryptor)

		err := client.SendTransferCompleted(context.Background(), []types.EnygmaTransferCompleted{})

		var wrappedErr *contractclient.EnygmaTeleportClientError
		require.ErrorAs(t, err, &wrappedErr, "error should be wrapped in EnygmaTeleportClientError")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})

	t.Run("wraps contract call errors in EnygmaTeleportClientError", func(t *testing.T) {
		wantError := errors.New("contract call failed")

		address := common.Address{}

		encryptor := &StubEnygmaTeleportEncryptor{
			encryptedData: []byte{0x01, 0x02},
		}
		executor := &StubEnygmaTeleportExecutor{
			err: wantError,
		}

		client := contractclient.NewEnygmaTeleportClient(address, executor, encryptor)

		err := client.SendTransferCompleted(context.Background(), []types.EnygmaTransferCompleted{})

		var wrappedErr *contractclient.EnygmaTeleportClientError
		require.ErrorAs(t, err, &wrappedErr, "error should be wrapped in EnygmaTeleportClientError")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})

	t.Run("wraps receipt fetch errors in EnygmaTeleportClientError", func(t *testing.T) {
		wantError := errors.New("receipt fetch failed")

		encryptor := &StubEnygmaTeleportEncryptor{
			encryptedData: []byte{0x01, 0x02},
		}
		executor := &StubEnygmaTeleportExecutor{
			err: wantError,
		}

		client := contractclient.NewEnygmaTeleportClient(common.Address{}, executor, encryptor)

		err := client.SendTransferCompleted(context.Background(), []types.EnygmaTransferCompleted{})

		var wrappedErr *contractclient.EnygmaTeleportClientError
		require.ErrorAs(t, err, &wrappedErr, "error should be wrapped in EnygmaTeleportClientError")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})

}
