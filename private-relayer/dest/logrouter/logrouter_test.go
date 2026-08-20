package logrouter_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/google/go-cmp/cmp"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/logrouter"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/logrouter/testdata"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogRouter(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("routes Endpoint logs to EndpointParser MQ", func(t *testing.T) {
		endpointAddress := common.HexToAddress("0xdeadc0de")

		blockNumber := uint64(1337)
		endpointLog := testdata.NewEndpointV1MessageDispatchedLogWith(
			testdata.WithMessageDispatchedBlockNumber(blockNumber),
			testdata.WithMessageDispatchedAddress(endpointAddress),
		)
		logs := []ethTypes.Log{endpointLog}

		stub := newLogRouterStub(t,
			withEndpointMQ(func(ctx context.Context, block logrouter.Block) error {
				assert.Equal(t, blockNumber, block.Number)
				assertLogsEqual(t, logs, block.Logs)
				return nil
			}),
			withEndpointAddress(endpointAddress),
		)

		router := stub.newLogRouter()

		err := router.Handle(context.TODO(), logs)
		require.Nil(t, err)
	})

	t.Run("splits logs by block number", func(t *testing.T) {
		endpointAddress := common.HexToAddress("0xdeadc0de")

		blockNumberA := uint64(1000)
		blockNumberB := uint64(1337)
		endpointLogA := testdata.NewEndpointV1MessageDispatchedLogWith(
			testdata.WithMessageDispatchedBlockNumber(blockNumberA),
			testdata.WithMessageDispatchedAddress(endpointAddress),
		)
		endpointLogB := testdata.NewEndpointV1MessageDispatchedLogWith(
			testdata.WithMessageDispatchedBlockNumber(blockNumberB),
			testdata.WithMessageDispatchedAddress(endpointAddress),
		)
		pushFuncCalls := 0
		logs := []ethTypes.Log{endpointLogB, endpointLogA}

		stub := newLogRouterStub(t,
			withEndpointMQ(func(ctx context.Context, block logrouter.Block) error {
				switch pushFuncCalls {
				case 0:
					assert.Equal(t, blockNumberA, block.Number)
					assertLogsEqual(t, []ethTypes.Log{endpointLogA}, block.Logs)
				case 1:
					assert.Equal(t, blockNumberB, block.Number)
					assertLogsEqual(t, []ethTypes.Log{endpointLogB}, block.Logs)
				}
				pushFuncCalls++
				return nil
			}),
			withEndpointAddress(endpointAddress),
		)

		router := stub.newLogRouter()

		err := router.Handle(context.TODO(), logs)
		require.Nil(t, err)
	})

	t.Run("routes Teleport logs to TeleportParser MQ", func(t *testing.T) {
		teleportAddress := common.HexToAddress("0xdeadc0de")

		blockNumber := uint64(1337)
		teleportLog := testdata.NewTeleportV1EncryptedDataBatchStoredLogWith(
			testdata.WithEncryptedDataBatchStoredBlockNumber(blockNumber),
			testdata.WithEncryptedDataBatchStoredAddress(teleportAddress),
		)
		logs := []ethTypes.Log{teleportLog}

		stub := newLogRouterStub(t,
			withTeleportMQ(func(ctx context.Context, block logrouter.Block) error {
				assert.Equal(t, blockNumber, block.Number)
				assertLogsEqual(t, logs, block.Logs)
				return nil
			}),
			withTeleportAddress(teleportAddress),
		)

		router := stub.newLogRouter()

		err := router.Handle(context.TODO(), logs)
		require.Nil(t, err)
	})

	t.Run("routes Enygma logs to EnygmaParser MQ", func(t *testing.T) {
		enygmaAddress := common.HexToAddress("0xdeadc0de")

		blockNumber := uint64(1337)
		enygmaLog := testdata.NewEnygmaTeleportEnygmaTransferLogWith(
			testdata.WithEnygmaTransferBlockNumber(blockNumber),
			testdata.WithEnygmaTransferAddress(enygmaAddress),
		)
		logs := []ethTypes.Log{enygmaLog}

		stub := newLogRouterStub(t,
			withEnygmaMQ(func(ctx context.Context, block logrouter.Block) error {
				assert.Equal(t, blockNumber, block.Number)
				assertLogsEqual(t, logs, block.Logs)
				return nil
			}),
			withEnygmaAddress(enygmaAddress),
		)

		router := stub.newLogRouter()

		err := router.Handle(context.TODO(), logs)
		require.Nil(t, err)
	})

	t.Run("routes Dvp logs to DvpParser MQ", func(t *testing.T) {
		dvpAddress := common.HexToAddress("0xdeadc0de")

		blockNumber := uint64(1337)
		dvpLog := testdata.NewDvpTeleportTransferEncryptedDataLogWith(
			testdata.WithDvpTransferEncryptedDataBlockNumber(blockNumber),
			testdata.WithDvpTransferEncryptedDataAddress(dvpAddress),
		)
		logs := []ethTypes.Log{dvpLog}

		stub := newLogRouterStub(t,
			withDvpMQ(func(ctx context.Context, block logrouter.Block) error {
				assert.Equal(t, blockNumber, block.Number)
				assertLogsEqual(t, logs, block.Logs)
				return nil
			}),
			withDvpAddress(dvpAddress),
		)

		router := stub.newLogRouter()

		err := router.Handle(context.TODO(), logs)
		require.Nil(t, err)
	})
}

func assertLogsEqual(t *testing.T, expected, actual []ethTypes.Log) {
	t.Helper()

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("Logs mismatch (-want +got):\n%s", diff)
	}
}

func TestLogRouter_RetriesOnFailure(t *testing.T) {
	t.Run("succeeds after retry", func(t *testing.T) {
		endpointAddress := common.HexToAddress("0xdeadc0de")

		blockNumber := uint64(1337)
		endpointLog := testdata.NewEndpointV1MessageDispatchedLogWith(
			testdata.WithMessageDispatchedBlockNumber(blockNumber),
			testdata.WithMessageDispatchedAddress(endpointAddress),
		)
		logs := []ethTypes.Log{endpointLog}

		attempts := 0
		stub := newLogRouterStub(t,
			withEndpointMQ(func(ctx context.Context, block logrouter.Block) error {
				attempts++
				if attempts < 3 {
					return assert.AnError // Fail first 2 attempts
				}
				return nil // Succeed on 3rd attempt
			}),
			withEndpointAddress(endpointAddress),
		)

		router := stub.newLogRouter()

		err := router.Handle(context.TODO(), logs)
		require.Nil(t, err)
		assert.Equal(t, 3, attempts)
	})

	t.Run("continues to next block after max retries", func(t *testing.T) {
		endpointAddress := common.HexToAddress("0xdeadc0de")

		blockNumber1 := uint64(1000)
		blockNumber2 := uint64(1001)
		log1 := testdata.NewEndpointV1MessageDispatchedLogWith(
			testdata.WithMessageDispatchedBlockNumber(blockNumber1),
			testdata.WithMessageDispatchedAddress(endpointAddress),
		)
		log2 := testdata.NewEndpointV1MessageDispatchedLogWith(
			testdata.WithMessageDispatchedBlockNumber(blockNumber2),
			testdata.WithMessageDispatchedAddress(endpointAddress),
		)
		logs := []ethTypes.Log{log1, log2}

		seenBlocks := []uint64{}
		stub := newLogRouterStub(t,
			withEndpointMQ(func(ctx context.Context, block logrouter.Block) error {
				seenBlocks = append(seenBlocks, block.Number)
				if block.Number == blockNumber1 {
					return assert.AnError // First block always fails
				}
				return nil // Second block succeeds
			}),
			withEndpointAddress(endpointAddress),
		)

		router := stub.newLogRouter()

		err := router.Handle(context.TODO(), logs)
		require.Nil(t, err)

		// Should have tried block 1 many times (100+), then succeeded with block 2 once
		assert.Greater(t, len(seenBlocks), 100)
		assert.Equal(t, blockNumber2, seenBlocks[len(seenBlocks)-1])
	})
}
