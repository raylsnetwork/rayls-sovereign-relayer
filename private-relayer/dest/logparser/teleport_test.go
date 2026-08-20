package logparser_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/gen/proto/encrypt"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/wireformat"
	"github.com/raylsnetwork/rayls-sovereign-relayer/msgqueue"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/logparser"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/logparser/testdata"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/logrouter"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools/fake"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools/spy"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func newTeleportSingleMessageMQ(msg msgqueue.Message[logrouter.Block]) *TeleportMQMock {
	return &TeleportMQMock{
		NextFunc: fake.NextMQ(msg),
	}
}

func newDefaultCrossChainMQ(t *testing.T) *CrossChainMQMock {
	return &CrossChainMQMock{
		PushBatchFunc: func(ctx context.Context, msgSlice []types.DispatchedMessageToPrivateHub) error {
			assert.Fail(t, "shouldn't call cross chain MQ push batch")
			return nil
		},
	}
}

func newDefaultCTSClient(t *testing.T) *CTSClientMock {
	return &CTSClientMock{
		DecryptFunc: func(ctx context.Context, in *encrypt.DecryptRequest, opts ...grpc.CallOption) (*encrypt.DecryptResponse, error) {
			assert.Fail(t, "shouldn't call CTS client decrypt")
			return nil, nil //nolint:nilnil // intentional nil return in test mock
		},
	}
}

func newDefaultEthereumClient(t *testing.T) *EthereumClientMock {
	return &EthereumClientMock{
		BlockByNumberFunc: func(ctx context.Context, number *big.Int) (*ethTypes.Block, error) {
			assert.Fail(t, "shouldn't call ethereum client block by number")
			return nil, nil //nolint:nilnil // intentional nil return in test mock
		},
	}
}

func newDefaultBackoff(t *testing.T) *BackoffMock {
	return &BackoffMock{
		DoFunc: func(ctx context.Context, maxAttempts int, fn func() error) error {
			assert.Fail(t, "shouldn't call backoff do")
			return nil
		},
	}
}

func newNoopSUMService() *SUMServiceMock {
	return &SUMServiceMock{
		BatchCreateFunc: func(ctx context.Context, sums []types.AtomicStatusUpdateMessage) error {
			return nil
		},
	}
}

func newDefaultSUMService(t *testing.T) *SUMServiceMock {
	return &SUMServiceMock{
		BatchCreateFunc: func(ctx context.Context, sums []types.AtomicStatusUpdateMessage) error {
			assert.Fail(t, "shouldn't call SUM service batch create")
			return nil
		},
	}
}

func TestTeleportParser(t *testing.T) {
	testtools.SilenceLogger()

	localChainID := big.NewInt(12345)
	hasGracefulShutdown := false
	t.Run("supports graceful shitdown", func(t *testing.T) {
		teleportMQ := &TeleportMQMock{
			NextFunc: func(ctx context.Context) (msgqueue.Message[logrouter.Block], error) {
				<-ctx.Done()
				return msgqueue.Message[logrouter.Block]{}, context.Canceled
			},
		}
		ctsClient := newDefaultCTSClient(t)
		ethClient := newDefaultEthereumClient(t)
		crossChainMQ := newDefaultCrossChainMQ(t)
		sumService := newNoopSUMService()

		parser := logparser.NewTeleportParser(
			teleportMQ,
			crossChainMQ,
			ctsClient,
			ethClient,
			localChainID,
			sumService,
		)
		hasGracefulShutdown = testtools.ShutdownFixture(t, parser.Run, time.Millisecond)
		assert.True(t, hasGracefulShutdown)
	})

	require.True(t, hasGracefulShutdown, "doesn't support graceful shutdown - exiting tests early")
	t.Run("parses events and pushes them to the MQ on successful decryption", func(t *testing.T) {
		wantData := common.Hex2Bytes("deadc0de")
		log := testdata.NewTeleportV1EncryptedDataBatchStoredLogWith(
			testdata.WithEncryptedDataBatchStoredData(wantData),
			testdata.WithEncryptedDataBatchStoredBlockNumber(1337),
		)

		wantPrevHash := common.HexToHash("0xc0cac01a")
		ethBlock := ethTypes.NewBlock(&ethTypes.Header{
			ParentHash: wantPrevHash,
		}, nil, nil, nil)

		wantNumber := big.NewInt(1337)
		block := logrouter.Block{
			Number: wantNumber.Uint64(),
			Logs:   []ethTypes.Log{log},
		}

		wantMsgs := []types.DispatchedMessageToPrivateHub{
			{
				SharedId: "example-shared-id-1",
			},
			{
				SharedId: "example-shared-id-2",
			},
		}

		ackSpy := spy.NewAck()
		teleportMQ := newTeleportSingleMessageMQ(msgqueue.Message[logrouter.Block]{
			V:   block,
			Ack: ackSpy.Fn(),
		})
		wantMsgsPlaintext, err := wireformat.MarshalPlaintext(wantMsgs)
		require.NoError(t, err)

		ctsClient := newDefaultCTSClient(t)
		ctsClient.DecryptFunc = func(ctx context.Context, in *encrypt.DecryptRequest, opts ...grpc.CallOption) (*encrypt.DecryptResponse, error) {
			ackSpy.AssertNotCalled(t, "should decrypt before acking")
			assert.Equal(t, wantPrevHash.String(), in.PrevBlockHash)
			assert.Equal(t, wantNumber.Uint64(), in.BlockNumber)
			assert.Equal(t, wantData, in.EncryptedData)
			return &encrypt.DecryptResponse{Plaintext: wantMsgsPlaintext, Outcome: encrypt.DecryptOutcome_OUTCOME_OK}, nil
		}
		ethClient := newDefaultEthereumClient(t)
		ethClient.BlockByNumberFunc = func(ctx context.Context, number *big.Int) (*ethTypes.Block, error) {
			assert.Equal(t, wantNumber, number)
			return ethBlock, nil
		}
		crossChainMQ := newDefaultCrossChainMQ(t)
		crossChainMQ.PushBatchFunc = func(ctx context.Context, msgSlice []types.DispatchedMessageToPrivateHub) error {
			ackSpy.AssertNotCalled(t, "shouldn't ack before pushing to the destination MQ")
			assert.Equal(t, wantMsgs, msgSlice)
			return nil
		}
		sumService := newNoopSUMService()

		parser := logparser.NewTeleportParser(
			teleportMQ,
			crossChainMQ,
			ctsClient,
			ethClient,
			localChainID,
			sumService,
		)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err = parser.Run(ctx)
		require.Nil(t, err)

		assert.Equal(t, 1, len(crossChainMQ.PushBatchCalls()), "should push to destination message queue exactly once")
		ackSpy.AssertCalled(t)
	})

	t.Run("implements exponential backoff on fail to push to destination message queue", func(t *testing.T) {
		crossChainMQError := errors.New("cross-chain-error")

		log := testdata.NewTeleportV1EncryptedDataBatchStoredLogWith()
		ethBlock := ethTypes.NewBlock(&ethTypes.Header{
			ParentHash: common.HexToHash("0xc0cac01a"),
		}, nil, nil, nil)
		block := logrouter.Block{
			Number: 1337,
			Logs:   []ethTypes.Log{log},
		}

		ackSpy := spy.NewAck()
		teleportMQ := newTeleportSingleMessageMQ(msgqueue.Message[logrouter.Block]{
			V:   block,
			Ack: ackSpy.Fn(),
		})
		decryptedPlaintext, err := wireformat.MarshalPlaintext([]types.DispatchedMessageToPrivateHub{{SharedId: "example-shared-id"}})
		require.NoError(t, err)

		ctsClient := newDefaultCTSClient(t)
		ctsClient.DecryptFunc = func(ctx context.Context, in *encrypt.DecryptRequest, opts ...grpc.CallOption) (*encrypt.DecryptResponse, error) {
			return &encrypt.DecryptResponse{Plaintext: decryptedPlaintext, Outcome: encrypt.DecryptOutcome_OUTCOME_OK}, nil
		}
		ethClient := newDefaultEthereumClient(t)
		ethClient.BlockByNumberFunc = func(ctx context.Context, number *big.Int) (*ethTypes.Block, error) {
			return ethBlock, nil
		}
		crossChainMQ := newDefaultCrossChainMQ(t)
		crossChainMQ.PushBatchFunc = func(ctx context.Context, msgSlice []types.DispatchedMessageToPrivateHub) error {
			return crossChainMQError
		}
		sumService := newNoopSUMService()

		backoffInCrossChainMQ := false
		backoff := newDefaultBackoff(t)
		backoff.DoFunc = func(ctx context.Context, maxAttempts int, fn func() error) error {
			ackSpy.AssertNotCalled(t, "shouldn't ack message before successfully pushing to destination MQ")
			err := fn()
			if errors.Is(err, crossChainMQError) {
				backoffInCrossChainMQ = true
			}
			return nil
		}

		parser := logparser.NewTeleportParserWithCustomBackoff(
			teleportMQ,
			crossChainMQ,
			ctsClient,
			ethClient,
			localChainID,
			sumService,
			backoff,
		)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err = parser.Run(ctx)
		require.Nil(t, err)

		assert.True(t, backoffInCrossChainMQ)
		ackSpy.AssertCalled(t)
	})

	t.Run("implements exponential backoff on fail to ack source message", func(t *testing.T) {
		ackError := errors.New("ack-error")
		log := testdata.NewTeleportV1EncryptedDataBatchStoredLogWith()
		ethBlock := ethTypes.NewBlock(&ethTypes.Header{
			ParentHash: common.HexToHash("0xc0cac01a"),
		}, nil, nil, nil)
		block := logrouter.Block{
			Number: 1337,
			Logs:   []ethTypes.Log{log},
		}

		ackSpy := spy.NewAck()
		teleportMQ := newTeleportSingleMessageMQ(msgqueue.Message[logrouter.Block]{
			V:   block,
			Ack: ackSpy.Fn(),
		})
		decryptedPlaintext, err := wireformat.MarshalPlaintext([]types.DispatchedMessageToPrivateHub{{SharedId: "example-shared-id"}})
		require.NoError(t, err)

		ctsClient := newDefaultCTSClient(t)
		ctsClient.DecryptFunc = func(ctx context.Context, in *encrypt.DecryptRequest, opts ...grpc.CallOption) (*encrypt.DecryptResponse, error) {
			return &encrypt.DecryptResponse{Plaintext: decryptedPlaintext, Outcome: encrypt.DecryptOutcome_OUTCOME_OK}, nil
		}
		ethClient := newDefaultEthereumClient(t)
		ethClient.BlockByNumberFunc = func(ctx context.Context, number *big.Int) (*ethTypes.Block, error) {
			return ethBlock, nil
		}
		crossChainMQ := newDefaultCrossChainMQ(t)
		crossChainMQ.PushBatchFunc = func(ctx context.Context, msgSlice []types.DispatchedMessageToPrivateHub) error {
			return ackError
		}
		sumService := newNoopSUMService()

		backoffInAck := false
		backoff := newDefaultBackoff(t)
		backoff.DoFunc = func(ctx context.Context, maxAttempts int, fn func() error) error {
			ackSpy.AssertNotCalled(t, "shouldn't ack message before successfully pushing to destination MQ")
			err := fn()
			if errors.Is(err, ackError) {
				backoffInAck = true
			}
			return nil
		}

		parser := logparser.NewTeleportParserWithCustomBackoff(
			teleportMQ,
			crossChainMQ,
			ctsClient,
			ethClient,
			localChainID,
			sumService,
			backoff,
		)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err = parser.Run(ctx)
		require.Nil(t, err)

		assert.True(t, backoffInAck)
		ackSpy.AssertCalled(t)
	})

	// AC-1: AtomicMessageStatusChangedBatch logs call SUMService.BatchCreate with correct SharedIDs and Status
	t.Run("parses AtomicMessageStatusChangedBatch and calls SUMService.BatchCreate", func(t *testing.T) {
		wantMsgIds := []string{"shared-id-1", "shared-id-2", "shared-id-3"}
		wantStatus := uint8(1) // AtomicExecutedStatus

		log := testdata.NewTeleportV1AtomicMessageStatusChangedBatchLogWith(
			testdata.WithAtomicStatusChangedMsgIds(wantMsgIds),
			testdata.WithAtomicStatusChangedStatus(wantStatus),
		)
		block := logrouter.Block{
			Number: 500,
			Logs:   []ethTypes.Log{log},
		}

		ackSpy := spy.NewAck()
		teleportMQ := newTeleportSingleMessageMQ(msgqueue.Message[logrouter.Block]{
			V:   block,
			Ack: ackSpy.Fn(),
		})
		ctsClient := newDefaultCTSClient(t)
		ethClient := newDefaultEthereumClient(t)
		crossChainMQ := newDefaultCrossChainMQ(t)

		sumService := newDefaultSUMService(t)
		sumService.BatchCreateFunc = func(ctx context.Context, sums []types.AtomicStatusUpdateMessage) error {
			require.Len(t, sums, 3)
			for i, sum := range sums {
				assert.Equal(t, wantMsgIds[i], sum.SharedID)
				assert.Equal(t, types.AtomicStatus(wantStatus), sum.Status)
			}
			return nil
		}

		backoff := newDefaultBackoff(t)
		backoff.DoFunc = func(ctx context.Context, maxAttempts int, fn func() error) error {
			return fn()
		}

		parser := logparser.NewTeleportParserWithCustomBackoff(
			teleportMQ,
			crossChainMQ,
			ctsClient,
			ethClient,
			localChainID,
			sumService,
			backoff,
		)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := parser.Run(ctx)
		require.Nil(t, err)

		assert.Equal(t, 1, len(sumService.BatchCreateCalls()), "should call SUM service batch create exactly once")
		ackSpy.AssertCalled(t)
	})

	// AC-2: Both EncryptedDataBatchStored and AtomicMessageStatusChangedBatch are handled independently
	t.Run(
		"handles both EncryptedDataBatchStored and AtomicMessageStatusChangedBatch independently",
		func(t *testing.T) {
			encryptedLog := testdata.NewTeleportV1EncryptedDataBatchStoredLogWith()
			atomicLog := testdata.NewTeleportV1AtomicMessageStatusChangedBatchLogWith(
				testdata.WithAtomicStatusChangedMsgIds([]string{"atomic-shared-id"}),
				testdata.WithAtomicStatusChangedStatus(uint8(3)), // AtomicRevertedStatus
			)

			ethBlock := ethTypes.NewBlock(&ethTypes.Header{
				ParentHash: common.HexToHash("0xc0cac01a"),
			}, nil, nil, nil)

			block := logrouter.Block{
				Number: 600,
				Logs:   []ethTypes.Log{encryptedLog, atomicLog},
			}

			ackSpy := spy.NewAck()
			teleportMQ := newTeleportSingleMessageMQ(msgqueue.Message[logrouter.Block]{
				V:   block,
				Ack: ackSpy.Fn(),
			})
			decryptedPlaintext, err := wireformat.MarshalPlaintext([]types.DispatchedMessageToPrivateHub{{SharedId: "encrypted-shared-id"}})
			require.NoError(t, err)

			ctsClient := newDefaultCTSClient(t)
			ctsClient.DecryptFunc = func(ctx context.Context, in *encrypt.DecryptRequest, opts ...grpc.CallOption) (*encrypt.DecryptResponse, error) {
				return &encrypt.DecryptResponse{Plaintext: decryptedPlaintext, Outcome: encrypt.DecryptOutcome_OUTCOME_OK}, nil
			}
			ethClient := newDefaultEthereumClient(t)
			ethClient.BlockByNumberFunc = func(ctx context.Context, number *big.Int) (*ethTypes.Block, error) {
				return ethBlock, nil
			}
			crossChainMQ := newDefaultCrossChainMQ(t)
			crossChainMQ.PushBatchFunc = func(ctx context.Context, msgSlice []types.DispatchedMessageToPrivateHub) error {
				return nil
			}
			sumService := newDefaultSUMService(t)
			sumService.BatchCreateFunc = func(ctx context.Context, sums []types.AtomicStatusUpdateMessage) error {
				require.Len(t, sums, 1)
				assert.Equal(t, "atomic-shared-id", sums[0].SharedID)
				assert.Equal(t, types.AtomicRevertedStatus, sums[0].Status)
				return nil
			}

			backoff := newDefaultBackoff(t)
			backoff.DoFunc = func(ctx context.Context, maxAttempts int, fn func() error) error {
				return fn()
			}

			parser := logparser.NewTeleportParserWithCustomBackoff(
				teleportMQ,
				crossChainMQ,
				ctsClient,
				ethClient,
				localChainID,
				sumService,
				backoff,
			)
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
			defer cancel()

			err = parser.Run(ctx)
			require.Nil(t, err)

			assert.Equal(t, 1, len(crossChainMQ.PushBatchCalls()), "should push encrypted messages to cross chain MQ")
			assert.Equal(t, 1, len(sumService.BatchCreateCalls()), "should call SUM service for atomic status")
			ackSpy.AssertCalled(t)
		},
	)

	// AC-3: BatchCreate error retries via backoff.Do, block still acked
	t.Run("retries BatchCreate via backoff on error and still acks block", func(t *testing.T) {
		batchCreateError := errors.New("batch-create-error")

		log := testdata.NewTeleportV1AtomicMessageStatusChangedBatchLogWith(
			testdata.WithAtomicStatusChangedMsgIds([]string{"retry-shared-id"}),
		)
		block := logrouter.Block{
			Number: 700,
			Logs:   []ethTypes.Log{log},
		}

		ackSpy := spy.NewAck()
		teleportMQ := newTeleportSingleMessageMQ(msgqueue.Message[logrouter.Block]{
			V:   block,
			Ack: ackSpy.Fn(),
		})
		ctsClient := newDefaultCTSClient(t)
		ethClient := newDefaultEthereumClient(t)
		crossChainMQ := newDefaultCrossChainMQ(t)

		sumService := newDefaultSUMService(t)
		sumService.BatchCreateFunc = func(ctx context.Context, sums []types.AtomicStatusUpdateMessage) error {
			return batchCreateError
		}

		backoffInBatchCreate := false
		backoff := newDefaultBackoff(t)
		backoff.DoFunc = func(ctx context.Context, maxAttempts int, fn func() error) error {
			err := fn()
			if errors.Is(err, batchCreateError) {
				backoffInBatchCreate = true
			}
			return nil // simulate backoff succeeding (retry exhaustion handled gracefully)
		}

		parser := logparser.NewTeleportParserWithCustomBackoff(
			teleportMQ,
			crossChainMQ,
			ctsClient,
			ethClient,
			localChainID,
			sumService,
			backoff,
		)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := parser.Run(ctx)
		require.Nil(t, err)

		assert.True(t, backoffInBatchCreate, "should use backoff for BatchCreate retry")
		ackSpy.AssertCalled(t)
	})

	// AC-3 (error path): BatchCreate fails after backoff exhaustion, block NOT acked for redelivery
	t.Run("does not ack block when BatchCreate fails after backoff exhaustion", func(t *testing.T) {
		batchCreateError := errors.New("batch-create-error")

		log := testdata.NewTeleportV1AtomicMessageStatusChangedBatchLogWith(
			testdata.WithAtomicStatusChangedMsgIds([]string{"fail-shared-id"}),
		)
		block := logrouter.Block{
			Number: 750,
			Logs:   []ethTypes.Log{log},
		}

		ackSpy := spy.NewAck()
		teleportMQ := newTeleportSingleMessageMQ(msgqueue.Message[logrouter.Block]{
			V:   block,
			Ack: ackSpy.Fn(),
		})
		ctsClient := newDefaultCTSClient(t)
		ethClient := newDefaultEthereumClient(t)
		crossChainMQ := newDefaultCrossChainMQ(t)

		sumService := newDefaultSUMService(t)
		sumService.BatchCreateFunc = func(ctx context.Context, sums []types.AtomicStatusUpdateMessage) error {
			return batchCreateError
		}

		backoff := newDefaultBackoff(t)
		backoff.DoFunc = func(ctx context.Context, maxAttempts int, fn func() error) error {
			return fn()
		}

		parser := logparser.NewTeleportParserWithCustomBackoff(
			teleportMQ,
			crossChainMQ,
			ctsClient,
			ethClient,
			localChainID,
			sumService,
			backoff,
		)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := parser.Run(ctx)
		require.Nil(t, err)

		assert.Equal(t, 1, len(sumService.BatchCreateCalls()), "should attempt BatchCreate")
		ackSpy.AssertNotCalled(t, "should not ack block when BatchCreate fails, allowing MQ redelivery")
	})

	t.Run("does not ack block when BlockByNumber fails in decryptMessages", func(t *testing.T) {
		log := testdata.NewTeleportV1EncryptedDataBatchStoredLogWith(
			testdata.WithEncryptedDataBatchStoredBlockNumber(1337),
		)
		block := logrouter.Block{
			Number: 1337,
			Logs:   []ethTypes.Log{log},
		}

		ackSpy := spy.NewAck()
		teleportMQ := newTeleportSingleMessageMQ(msgqueue.Message[logrouter.Block]{
			V:   block,
			Ack: ackSpy.Fn(),
		})
		ctsClient := newDefaultCTSClient(t)
		ethClient := newDefaultEthereumClient(t)
		ethClient.BlockByNumberFunc = func(ctx context.Context, number *big.Int) (*ethTypes.Block, error) {
			return nil, errors.New("ethereum node unavailable")
		}
		crossChainMQ := newDefaultCrossChainMQ(t)
		sumService := newNoopSUMService()

		parser := logparser.NewTeleportParser(
			teleportMQ,
			crossChainMQ,
			ctsClient,
			ethClient,
			localChainID,
			sumService,
		)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := parser.Run(ctx)
		require.Nil(t, err)

		assert.Equal(
			t,
			0,
			len(ctsClient.DecryptCalls()),
			"should not attempt decryption when block fetch fails",
		)
		assert.Equal(t, 0, len(crossChainMQ.PushBatchCalls()), "should not push to cross chain MQ")
		ackSpy.AssertNotCalled(t, "should not ack block when BlockByNumber fails, allowing MQ redelivery")
	})

	// AC-4: Unrecognized log topic is skipped, no service calls made, block still acked
	t.Run("skips unrecognized log topics and still acks block", func(t *testing.T) {
		unrecognizedLog := ethTypes.Log{
			Address:     common.HexToAddress("0x9999999999999999999999999999999999999999"),
			BlockNumber: 900,
			Topics: []common.Hash{
				common.HexToHash("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
			},
			Data: []byte{0x01, 0x02, 0x03},
		}
		block := logrouter.Block{
			Number: 900,
			Logs:   []ethTypes.Log{unrecognizedLog},
		}

		ackSpy := spy.NewAck()
		teleportMQ := newTeleportSingleMessageMQ(msgqueue.Message[logrouter.Block]{
			V:   block,
			Ack: ackSpy.Fn(),
		})
		ctsClient := newDefaultCTSClient(t)
		ethClient := newDefaultEthereumClient(t)
		crossChainMQ := newDefaultCrossChainMQ(t)
		sumService := newDefaultSUMService(t)

		backoff := newDefaultBackoff(t)
		backoff.DoFunc = func(ctx context.Context, maxAttempts int, fn func() error) error {
			return fn()
		}

		parser := logparser.NewTeleportParserWithCustomBackoff(
			teleportMQ,
			crossChainMQ,
			ctsClient,
			ethClient,
			localChainID,
			sumService,
			backoff,
		)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := parser.Run(ctx)
		require.Nil(t, err)

		assert.Equal(t, 0, len(crossChainMQ.PushBatchCalls()), "should not call cross chain MQ for unrecognized log")
		assert.Equal(t, 0, len(sumService.BatchCreateCalls()), "should not call SUM service for unrecognized log")
		ackSpy.AssertCalled(t)
	})

	// AC-6: Empty AtomicMessageStatusChangedBatch (no msgIds) does not call BatchCreate
	t.Run("skips BatchCreate when AtomicMessageStatusChangedBatch has no msgIds", func(t *testing.T) {
		log := testdata.NewTeleportV1AtomicMessageStatusChangedBatchLogWith(
			testdata.WithAtomicStatusChangedMsgIds([]string{}),
		)
		block := logrouter.Block{
			Number: 800,
			Logs:   []ethTypes.Log{log},
		}

		ackSpy := spy.NewAck()
		teleportMQ := newTeleportSingleMessageMQ(msgqueue.Message[logrouter.Block]{
			V:   block,
			Ack: ackSpy.Fn(),
		})
		ctsClient := newDefaultCTSClient(t)
		ethClient := newDefaultEthereumClient(t)
		crossChainMQ := newDefaultCrossChainMQ(t)
		sumService := newDefaultSUMService(t)

		backoff := newDefaultBackoff(t)
		backoff.DoFunc = func(ctx context.Context, maxAttempts int, fn func() error) error {
			return fn()
		}

		parser := logparser.NewTeleportParserWithCustomBackoff(
			teleportMQ,
			crossChainMQ,
			ctsClient,
			ethClient,
			localChainID,
			sumService,
			backoff,
		)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := parser.Run(ctx)
		require.Nil(t, err)

		assert.Equal(t, 0, len(sumService.BatchCreateCalls()), "should not call BatchCreate for empty msgIds")
		ackSpy.AssertCalled(t)
	})

	t.Run("does not ack block when PushBatch fails after backoff", func(t *testing.T) {
		log := testdata.NewTeleportV1EncryptedDataBatchStoredLogWith(
			testdata.WithEncryptedDataBatchStoredBlockNumber(1337),
		)
		ethBlock := ethTypes.NewBlock(&ethTypes.Header{
			ParentHash: common.HexToHash("0xc0cac01a"),
		}, nil, nil, nil)
		block := logrouter.Block{
			Number: 1337,
			Logs:   []ethTypes.Log{log},
		}

		ackSpy := spy.NewAck()
		teleportMQ := newTeleportSingleMessageMQ(msgqueue.Message[logrouter.Block]{
			V:   block,
			Ack: ackSpy.Fn(),
		})
		decryptedPlaintext, err := wireformat.MarshalPlaintext([]types.DispatchedMessageToPrivateHub{{SharedId: "example-shared-id"}})
		require.NoError(t, err)

		ctsClient := newDefaultCTSClient(t)
		ctsClient.DecryptFunc = func(ctx context.Context, in *encrypt.DecryptRequest, opts ...grpc.CallOption) (*encrypt.DecryptResponse, error) {
			return &encrypt.DecryptResponse{Plaintext: decryptedPlaintext, Outcome: encrypt.DecryptOutcome_OUTCOME_OK}, nil
		}
		ethClient := newDefaultEthereumClient(t)
		ethClient.BlockByNumberFunc = func(ctx context.Context, number *big.Int) (*ethTypes.Block, error) {
			return ethBlock, nil
		}
		crossChainMQ := newDefaultCrossChainMQ(t)
		crossChainMQ.PushBatchFunc = func(ctx context.Context, msgSlice []types.DispatchedMessageToPrivateHub) error {
			return errors.New("cross chain MQ unavailable")
		}
		sumService := newNoopSUMService()

		backoff := newDefaultBackoff(t)
		backoff.DoFunc = func(ctx context.Context, maxAttempts int, fn func() error) error {
			return fn()
		}

		parser := logparser.NewTeleportParserWithCustomBackoff(
			teleportMQ,
			crossChainMQ,
			ctsClient,
			ethClient,
			localChainID,
			sumService,
			backoff,
		)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err = parser.Run(ctx)
		require.Nil(t, err)

		assert.Equal(t, 1, len(crossChainMQ.PushBatchCalls()), "should attempt PushBatch")
		ackSpy.AssertNotCalled(t, "should not ack block when PushBatch fails, allowing MQ redelivery")
	})
}
