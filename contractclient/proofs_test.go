package contractclient_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contractclient"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/Proofs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeProofsHeader(number int64) Proofs.ProofsHeader {
	return Proofs.ProofsHeader{
		Number:        big.NewInt(number),
		Difficulty:    big.NewInt(0),
		GasLimit:      big.NewInt(0),
		GasUsed:       big.NewInt(0),
		Time:          big.NewInt(0),
		BaseFeePerGas: big.NewInt(0),
	}
}

func TestProofsClient_GetNextBlockNumber(t *testing.T) {
	t.Run("returns block number via executor call", func(t *testing.T) {
		wantChainID := big.NewInt(100)

		// We need ABI-encoded result. Use the Proofs contract to pack a response.
		parsed, err := Proofs.ProofsMetaData.ParseABI()
		require.NoError(t, err)
		method := parsed.Methods["getNextHeaderBlockNumber"]
		callResult, err := method.Outputs.Pack(big.NewInt(42))
		require.NoError(t, err)

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{
			callResult: callResult,
		}

		client := contractclient.NewProofsClient(address, executor)

		gotBlockNumber, err := client.GetNextBlockNumber(context.Background(), wantChainID)
		require.NoError(t, err)

		assert.Equal(t, big.NewInt(42), gotBlockNumber)
		assert.Equal(t, address, executor.spyCallAddress)
	})

	t.Run("wraps executor errors in ProofsClientError", func(t *testing.T) {
		wantError := errors.New("call error")
		wantErrorType := &contractclient.ProofsClientError{}

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{
			callErr: wantError,
		}

		client := contractclient.NewProofsClient(address, executor)

		_, err := client.GetNextBlockNumber(context.Background(), big.NewInt(100))

		require.ErrorAs(t, err, &wantErrorType)
		require.ErrorIs(t, err, wantError)
	})
}

func TestProofsClient_SubmitBatchHeaders(t *testing.T) {
	t.Run("returns ProofsClientError for empty headers", func(t *testing.T) {
		wantErrorType := &contractclient.ProofsClientError{}

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}

		client := contractclient.NewProofsClient(address, executor)

		_, err := client.SubmitBatchHeaders(context.Background(), big.NewInt(100), []Proofs.ProofsHeader{})

		require.ErrorAs(t, err, &wantErrorType)
		assert.Contains(t, err.Error(), "no headers to submit")
	})

	t.Run("wraps executor execute errors in ProofsClientError", func(t *testing.T) {
		wantError := errors.New("execute error")
		wantErrorType := &contractclient.ProofsClientError{}

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{
			executeErr: wantError,
		}

		client := contractclient.NewProofsClient(address, executor)

		headers := []Proofs.ProofsHeader{makeProofsHeader(10)}
		_, err := client.SubmitBatchHeaders(context.Background(), big.NewInt(100), headers)

		require.ErrorAs(t, err, &wantErrorType)
		require.ErrorIs(t, err, wantError)
		assert.Contains(t, err.Error(), "failed to submit batch headers")
	})

	t.Run("successfully submits and returns result with start end and next block", func(t *testing.T) {
		headers := []Proofs.ProofsHeader{
			makeProofsHeader(10),
			makeProofsHeader(11),
			makeProofsHeader(12),
		}

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}

		client := contractclient.NewProofsClient(address, executor)

		result, err := client.SubmitBatchHeaders(context.Background(), big.NewInt(100), headers)
		require.NoError(t, err)

		assert.Equal(t, big.NewInt(10), result.StartBlock)
		assert.Equal(t, big.NewInt(12), result.EndBlock)
		assert.Equal(t, big.NewInt(13), result.NextExpectedBlock)
	})
}
