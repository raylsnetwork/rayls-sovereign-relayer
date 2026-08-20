package logparser_test

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	encryptpb "github.com/raylsnetwork/rayls-sovereign-relayer/cts/gen/proto/encrypt"
	"github.com/raylsnetwork/rayls-sovereign-relayer/msgqueue"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/logparser"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/logparser/testdata"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/logrouter"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools/fake"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools/spy"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func newDvpSingleBlockMQ(block msgqueue.Message[logrouter.Block]) *DvpTeleportMQMock {
	return &DvpTeleportMQMock{
		NextFunc: fake.NextMQ(block),
	}
}

func newDefaultDvpCTSMock() *DvpCTSClientMock {
	return &DvpCTSClientMock{}
}

func newDefaultDvpSwapRepoMock() *DvpSwapRepositoryMock {
	return &DvpSwapRepositoryMock{
		GetSwapBySharedIDFunc: func(ctx context.Context, sharedID string) (*types.DvpSwap, error) {
			return &types.DvpSwap{SharedID: sharedID}, nil
		},
	}
}

func newDefaultDvpBackoffMock() *DvpBackoffMock {
	return &DvpBackoffMock{
		DoFunc: func(ctx context.Context, maxAttempts int, fn func() error) error {
			return fn()
		},
	}
}

func newDefaultDvpEthClientMock() *DvpEthereumClientMock {
	return &DvpEthereumClientMock{
		BlockByNumberFunc: func(ctx context.Context, number *big.Int) (*ethTypes.Block, error) {
			header := &ethTypes.Header{
				ParentHash: [32]byte{0x01, 0x02, 0x03},
			}
			return ethTypes.NewBlockWithHeader(header), nil
		},
	}
}

func TestDvpTeleportParser(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("supports graceful shutdown", func(t *testing.T) {
		teleportMQ := &DvpTeleportMQMock{
			NextFunc: func(ctx context.Context) (msgqueue.Message[logrouter.Block], error) {
				<-ctx.Done()
				return msgqueue.Message[logrouter.Block]{}, context.Canceled
			},
		}
		orchestratorMQ := &DvpOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.DvpDestMessage) error {
				return nil
			},
		}

		parser := logparser.NewDvpTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			newDefaultDvpCTSMock(),
			newDefaultDvpEthClientMock(),
			newDefaultDvpSwapRepoMock(),
			big.NewInt(12345),
			newDefaultDvpBackoffMock(),
		)

		hasGracefulShutdown := testtools.ShutdownFixture(t, parser.Run, time.Second)
		assert.True(t, hasGracefulShutdown)
	})

	t.Run("acks block after processing logs with no matching events", func(t *testing.T) {
		ackSpy := spy.NewAck()

		teleportMQ := newDvpSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{
				Number: 100,
				Logs:   []ethTypes.Log{},
			},
			Ack: ackSpy.Fn(),
		})
		orchestratorMQ := &DvpOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.DvpDestMessage) error {
				assert.Fail(t, "should not push any messages for empty logs")
				return nil
			},
		}

		parser := logparser.NewDvpTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			newDefaultDvpCTSMock(),
			newDefaultDvpEthClientMock(),
			newDefaultDvpSwapRepoMock(),
			big.NewInt(12345),
			newDefaultDvpBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := parser.Run(ctx)
		require.NoError(t, err)

		ackSpy.AssertCalled(t)
		assert.Equal(t, 0, len(orchestratorMQ.PushCalls()))
	})

	t.Run("constructor creates parser with valid teleport filterer", func(t *testing.T) {
		teleportMQ := &DvpTeleportMQMock{}
		orchestratorMQ := &DvpOrchestratorMQMock{}

		parser := logparser.NewDvpTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			newDefaultDvpCTSMock(),
			newDefaultDvpEthClientMock(),
			newDefaultDvpSwapRepoMock(),
			big.NewInt(12345),
			newDefaultDvpBackoffMock(),
		)

		assert.NotNil(t, parser)
	})

	t.Run("parses Commitments and pushes to orchestrator", func(t *testing.T) {
		wantBlockNumber := uint64(300)
		wantTokenAddress := common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
		wantTokenType := big.NewInt(4)
		wantTreeNumber := big.NewInt(2)
		wantCommitments := []*big.Int{big.NewInt(1000), big.NewInt(2000), big.NewInt(3000)}

		log := testdata.NewCommitmentsLogWith(
			testdata.WithCommitmentsLogBlockNumber(wantBlockNumber),
			testdata.WithCommitmentsTokenAddress(wantTokenAddress),
			testdata.WithCommitmentsTokenType(wantTokenType),
			testdata.WithCommitmentsTreeNumber(wantTreeNumber),
			testdata.WithCommitmentsValues(wantCommitments),
		)

		ackSpy := spy.NewAck()
		teleportMQ := newDvpSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{
				Number: wantBlockNumber,
				Logs:   []ethTypes.Log{log},
			},
			Ack: ackSpy.Fn(),
		})

		var pushedMsg service.DvpDestMessage
		orchestratorMQ := &DvpOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.DvpDestMessage) error {
				pushedMsg = msg
				return nil
			},
		}

		parser := logparser.NewDvpTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			newDefaultDvpCTSMock(),
			newDefaultDvpEthClientMock(),
			newDefaultDvpSwapRepoMock(),
			big.NewInt(12345),
			newDefaultDvpBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := parser.Run(ctx)
		require.NoError(t, err)

		require.Equal(t, 1, len(orchestratorMQ.PushCalls()))
		assert.Equal(t, service.DvpCommitmentsMessage, pushedMsg.Type)
		assert.Equal(t, wantBlockNumber, pushedMsg.BlockNumber)
		assert.Contains(t, pushedMsg.ID, "dvp-commitments-")
		require.NotNil(t, pushedMsg.Commitments)
		assert.Equal(t, wantTokenAddress.Hex(), pushedMsg.Commitments.TokenAddress)
		assert.Equal(t, 0, wantTokenType.Cmp(pushedMsg.Commitments.TokenType))
		assert.Equal(t, 0, wantTreeNumber.Cmp(pushedMsg.Commitments.TreeNumber))
		require.Len(t, pushedMsg.Commitments.Commitments, 3)
		assert.Equal(t, 0, wantCommitments[0].Cmp(pushedMsg.Commitments.Commitments[0]))
		assert.Equal(t, 0, wantCommitments[1].Cmp(pushedMsg.Commitments.Commitments[1]))
		assert.Equal(t, 0, wantCommitments[2].Cmp(pushedMsg.Commitments.Commitments[2]))
		ackSpy.AssertCalled(t)
	})

	t.Run("does not ack block when orchestrator push fails for Commitments", func(t *testing.T) {
		log := testdata.NewCommitmentsLogWith()

		ackSpy := spy.NewAck()
		teleportMQ := newDvpSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{
				Number: 100,
				Logs:   []ethTypes.Log{log},
			},
			Ack: ackSpy.Fn(),
		})

		orchestratorMQ := &DvpOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.DvpDestMessage) error {
				return errors.New("push failed")
			},
		}

		parser := logparser.NewDvpTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			newDefaultDvpCTSMock(),
			newDefaultDvpEthClientMock(),
			newDefaultDvpSwapRepoMock(),
			big.NewInt(12345),
			newDefaultDvpBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := parser.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(orchestratorMQ.PushCalls()))
		ackSpy.AssertNotCalled(t, "should not ack block when push fails, allowing MQ redelivery")
	})

	t.Run("parses Nullifiers and pushes to orchestrator", func(t *testing.T) {
		wantBlockNumber := uint64(400)
		wantNullifiers := []*big.Int{big.NewInt(55555)}

		log := testdata.NewNullifiersLogWith(
			testdata.WithNullifiersLogBlockNumber(wantBlockNumber),
			testdata.WithNullifiersValues(wantNullifiers),
		)

		ackSpy := spy.NewAck()
		teleportMQ := newDvpSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{
				Number: wantBlockNumber,
				Logs:   []ethTypes.Log{log},
			},
			Ack: ackSpy.Fn(),
		})

		var pushedMsg service.DvpDestMessage
		orchestratorMQ := &DvpOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.DvpDestMessage) error {
				pushedMsg = msg
				return nil
			},
		}

		parser := logparser.NewDvpTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			newDefaultDvpCTSMock(),
			newDefaultDvpEthClientMock(),
			newDefaultDvpSwapRepoMock(),
			big.NewInt(12345),
			newDefaultDvpBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := parser.Run(ctx)
		require.NoError(t, err)

		require.Equal(t, 1, len(orchestratorMQ.PushCalls()))
		assert.Equal(t, service.DvpNullifierMessage, pushedMsg.Type)
		assert.Equal(t, wantBlockNumber, pushedMsg.BlockNumber)
		assert.Contains(t, pushedMsg.ID, "dvp-nullifiers-")
		require.NotNil(t, pushedMsg.Nullifiers)
		require.Len(t, pushedMsg.Nullifiers.Nullifiers, 1)
		assert.Equal(t, 0, wantNullifiers[0].Cmp(pushedMsg.Nullifiers.Nullifiers[0]))
		ackSpy.AssertCalled(t)
	})

	t.Run("does not ack block when orchestrator push fails for Nullifiers", func(t *testing.T) {
		log := testdata.NewNullifiersLogWith()

		ackSpy := spy.NewAck()
		teleportMQ := newDvpSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{
				Number: 100,
				Logs:   []ethTypes.Log{log},
			},
			Ack: ackSpy.Fn(),
		})

		orchestratorMQ := &DvpOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.DvpDestMessage) error {
				return errors.New("push failed")
			},
		}

		parser := logparser.NewDvpTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			newDefaultDvpCTSMock(),
			newDefaultDvpEthClientMock(),
			newDefaultDvpSwapRepoMock(),
			big.NewInt(12345),
			newDefaultDvpBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := parser.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(orchestratorMQ.PushCalls()))
		ackSpy.AssertNotCalled(t, "should not ack block when push fails, allowing MQ redelivery")
	})

	t.Run("skips SwapCompleted when swap unknown locally", func(t *testing.T) {
		sharedId := [32]byte{0xaa, 0xbb, 0xcc, 0xdd}
		log := testdata.NewSwapCompletedLogWith(
			testdata.WithSwapCompletedSharedId(sharedId),
		)

		ackSpy := spy.NewAck()
		teleportMQ := newDvpSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{
				Number: 100,
				Logs:   []ethTypes.Log{log},
			},
			Ack: ackSpy.Fn(),
		})

		swapRepo := &DvpSwapRepositoryMock{
			GetSwapBySharedIDFunc: func(ctx context.Context, sharedID string) (*types.DvpSwap, error) {
				return nil, nil
			},
		}
		ctsMock := newDefaultDvpCTSMock()
		orchestratorMQ := &DvpOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.DvpDestMessage) error {
				assert.Fail(t, "should not push when swap is unknown locally")
				return nil
			},
		}

		parser := logparser.NewDvpTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			ctsMock,
			newDefaultDvpEthClientMock(),
			swapRepo,
			big.NewInt(12345),
			newDefaultDvpBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		require.NoError(t, parser.Run(ctx))

		assert.Empty(t, ctsMock.DecryptWithoutFPWithSSCalls())
		assert.Empty(t, swapRepo.UpdateSwapToCalls())
		assert.Empty(t, orchestratorMQ.PushCalls())
		ackSpy.AssertCalled(t)
	})

	t.Run("does not ack SwapCompleted when swap lookup fails", func(t *testing.T) {
		log := testdata.NewSwapCompletedLogWith()

		ackSpy := spy.NewAck()
		teleportMQ := newDvpSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{
				Number: 100,
				Logs:   []ethTypes.Log{log},
			},
			Ack: ackSpy.Fn(),
		})

		swapRepo := &DvpSwapRepositoryMock{
			GetSwapBySharedIDFunc: func(ctx context.Context, sharedID string) (*types.DvpSwap, error) {
				return nil, errors.New("boom")
			},
		}
		ctsMock := newDefaultDvpCTSMock()
		orchestratorMQ := &DvpOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.DvpDestMessage) error {
				assert.Fail(t, "should not push when swap lookup fails")
				return nil
			},
		}

		parser := logparser.NewDvpTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			ctsMock,
			newDefaultDvpEthClientMock(),
			swapRepo,
			big.NewInt(12345),
			newDefaultDvpBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		require.NoError(t, parser.Run(ctx))

		assert.Empty(t, ctsMock.DecryptWithoutFPWithSSCalls())
		assert.Empty(t, swapRepo.UpdateSwapToCalls())
		assert.Empty(t, orchestratorMQ.PushCalls())
		ackSpy.AssertNotCalled(t, "should not ack block when swap lookup fails, allowing MQ redelivery")
	})

	t.Run("pushes SwapCompleted directly for responder side", func(t *testing.T) {
		wantBlockNumber := uint64(500)
		sharedId := [32]byte{0x11, 0x22, 0x33, 0x44}
		txHash := common.HexToHash("0xfeedbeefcafe")
		log := testdata.NewSwapCompletedLogWith(
			testdata.WithSwapCompletedSharedId(sharedId),
			testdata.WithSwapCompletedLogBlockNumber(wantBlockNumber),
			testdata.WithSwapCompletedTxHash(txHash),
		)

		ackSpy := spy.NewAck()
		teleportMQ := newDvpSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{
				Number: wantBlockNumber,
				Logs:   []ethTypes.Log{log},
			},
			Ack: ackSpy.Fn(),
		})

		swapRepo := &DvpSwapRepositoryMock{
			GetSwapBySharedIDFunc: func(ctx context.Context, sharedID string) (*types.DvpSwap, error) {
				return &types.DvpSwap{
					SharedID: sharedID,
					To:       "0xresponder",
					SelfSalt: big.NewInt(1),
				}, nil
			},
		}
		ctsMock := newDefaultDvpCTSMock()

		var pushedMsg service.DvpDestMessage
		orchestratorMQ := &DvpOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.DvpDestMessage) error {
				pushedMsg = msg
				return nil
			},
		}

		parser := logparser.NewDvpTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			ctsMock,
			newDefaultDvpEthClientMock(),
			swapRepo,
			big.NewInt(12345),
			newDefaultDvpBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		require.NoError(t, parser.Run(ctx))

		assert.Empty(t, ctsMock.DecryptWithoutFPWithSSCalls(), "responder side must not decrypt")
		assert.Empty(t, swapRepo.UpdateSwapToCalls(), "responder side must not update swap.To")
		require.Equal(t, 1, len(orchestratorMQ.PushCalls()))
		assert.Equal(t, service.DvpSwapCompletedMessage, pushedMsg.Type)
		assert.Equal(t, wantBlockNumber, pushedMsg.BlockNumber)
		assert.Contains(t, pushedMsg.ID, "dvp-completed-")
		assert.Equal(t, common.Bytes2Hex(sharedId[:]), pushedMsg.SharedID)
		ackSpy.AssertCalled(t)
	})

	t.Run("decrypts SwapCompleted for initiator and updates swap.To", func(t *testing.T) {
		wantBlockNumber := uint64(600)
		sharedId := [32]byte{0x55, 0x66, 0x77, 0x88}
		encryptedData := []byte{0xde, 0xad, 0xbe, 0xef}
		selfSalt := big.NewInt(42)
		wantTo := "0x1234567890AbcdEF1234567890aBcdef12345678"

		log := testdata.NewSwapCompletedLogWith(
			testdata.WithSwapCompletedSharedId(sharedId),
			testdata.WithSwapCompletedEncryptedData(encryptedData),
			testdata.WithSwapCompletedLogBlockNumber(wantBlockNumber),
		)

		ackSpy := spy.NewAck()
		teleportMQ := newDvpSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{
				Number: wantBlockNumber,
				Logs:   []ethTypes.Log{log},
			},
			Ack: ackSpy.Fn(),
		})

		swapRepo := &DvpSwapRepositoryMock{
			GetSwapBySharedIDFunc: func(ctx context.Context, sharedID string) (*types.DvpSwap, error) {
				return &types.DvpSwap{
					SharedID: sharedID,
					To:       "",
					SelfSalt: selfSalt,
				}, nil
			},
			UpdateSwapToFunc: func(ctx context.Context, sharedID string, to string) error {
				return nil
			},
		}

		plaintext, err := json.Marshal(&types.DvpSwapMessage{To: wantTo})
		require.NoError(t, err)
		ctsMock := &DvpCTSClientMock{
			DecryptWithoutFPWithSSFunc: func(ctx context.Context, in *encryptpb.DecryptWithoutFPWithSSRequest, opts ...grpc.CallOption) (*encryptpb.DecryptWithoutFPWithSSResponse, error) {
				return &encryptpb.DecryptWithoutFPWithSSResponse{
					Outcome:   encryptpb.DecryptOutcome_OUTCOME_OK,
					Plaintext: plaintext,
				}, nil
			},
		}

		var pushedMsg service.DvpDestMessage
		orchestratorMQ := &DvpOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.DvpDestMessage) error {
				pushedMsg = msg
				return nil
			},
		}

		parser := logparser.NewDvpTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			ctsMock,
			newDefaultDvpEthClientMock(),
			swapRepo,
			big.NewInt(12345),
			newDefaultDvpBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		require.NoError(t, parser.Run(ctx))

		require.Equal(t, 1, len(ctsMock.DecryptWithoutFPWithSSCalls()))
		decryptCall := ctsMock.DecryptWithoutFPWithSSCalls()[0]
		assert.Equal(t, encryptedData, decryptCall.In.EncryptedData)
		assert.Equal(t, selfSalt.Bytes(), decryptCall.In.Ss)

		require.Equal(t, 1, len(swapRepo.UpdateSwapToCalls()))
		updateCall := swapRepo.UpdateSwapToCalls()[0]
		assert.Equal(t, common.Bytes2Hex(sharedId[:]), updateCall.SharedID)
		assert.Equal(t, wantTo, updateCall.To)

		require.Equal(t, 1, len(orchestratorMQ.PushCalls()))
		assert.Equal(t, service.DvpSwapCompletedMessage, pushedMsg.Type)
		assert.Equal(t, wantBlockNumber, pushedMsg.BlockNumber)
		assert.Equal(t, common.Bytes2Hex(sharedId[:]), pushedMsg.SharedID)
		ackSpy.AssertCalled(t)
	})

	t.Run("skips SwapCompleted when decrypt outcome is NOT_FOR_RECIPIENT", func(t *testing.T) {
		log := testdata.NewSwapCompletedLogWith()

		ackSpy := spy.NewAck()
		teleportMQ := newDvpSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{Number: 100, Logs: []ethTypes.Log{log}},
			Ack: ackSpy.Fn(),
		})

		swapRepo := &DvpSwapRepositoryMock{
			GetSwapBySharedIDFunc: func(ctx context.Context, sharedID string) (*types.DvpSwap, error) {
				return &types.DvpSwap{SharedID: sharedID, To: "", SelfSalt: big.NewInt(1)}, nil
			},
		}
		ctsMock := &DvpCTSClientMock{
			DecryptWithoutFPWithSSFunc: func(ctx context.Context, in *encryptpb.DecryptWithoutFPWithSSRequest, opts ...grpc.CallOption) (*encryptpb.DecryptWithoutFPWithSSResponse, error) {
				return &encryptpb.DecryptWithoutFPWithSSResponse{
					Outcome: encryptpb.DecryptOutcome_OUTCOME_NOT_FOR_RECIPIENT,
				}, nil
			},
		}
		orchestratorMQ := &DvpOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.DvpDestMessage) error {
				assert.Fail(t, "should not push when decrypt outcome is NOT_FOR_RECIPIENT")
				return nil
			},
		}

		parser := logparser.NewDvpTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			ctsMock,
			newDefaultDvpEthClientMock(),
			swapRepo,
			big.NewInt(12345),
			newDefaultDvpBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		require.NoError(t, parser.Run(ctx))

		assert.Equal(t, 1, len(ctsMock.DecryptWithoutFPWithSSCalls()))
		assert.Empty(t, swapRepo.UpdateSwapToCalls())
		assert.Empty(t, orchestratorMQ.PushCalls())
		ackSpy.AssertCalled(t)
	})

	t.Run("skips SwapCompleted when decrypt outcome is TAMPERED", func(t *testing.T) {
		log := testdata.NewSwapCompletedLogWith()

		ackSpy := spy.NewAck()
		teleportMQ := newDvpSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{Number: 100, Logs: []ethTypes.Log{log}},
			Ack: ackSpy.Fn(),
		})

		swapRepo := &DvpSwapRepositoryMock{
			GetSwapBySharedIDFunc: func(ctx context.Context, sharedID string) (*types.DvpSwap, error) {
				return &types.DvpSwap{SharedID: sharedID, To: "", SelfSalt: big.NewInt(1)}, nil
			},
		}
		ctsMock := &DvpCTSClientMock{
			DecryptWithoutFPWithSSFunc: func(ctx context.Context, in *encryptpb.DecryptWithoutFPWithSSRequest, opts ...grpc.CallOption) (*encryptpb.DecryptWithoutFPWithSSResponse, error) {
				return &encryptpb.DecryptWithoutFPWithSSResponse{
					Outcome: encryptpb.DecryptOutcome_OUTCOME_TAMPERED,
				}, nil
			},
		}
		orchestratorMQ := &DvpOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.DvpDestMessage) error {
				assert.Fail(t, "should not push when decrypt outcome is TAMPERED")
				return nil
			},
		}

		parser := logparser.NewDvpTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			ctsMock,
			newDefaultDvpEthClientMock(),
			swapRepo,
			big.NewInt(12345),
			newDefaultDvpBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		require.NoError(t, parser.Run(ctx))

		assert.Equal(t, 1, len(ctsMock.DecryptWithoutFPWithSSCalls()))
		assert.Empty(t, swapRepo.UpdateSwapToCalls())
		assert.Empty(t, orchestratorMQ.PushCalls())
		ackSpy.AssertCalled(t)
	})

	t.Run("does not ack SwapCompleted when decrypt fails", func(t *testing.T) {
		log := testdata.NewSwapCompletedLogWith()

		ackSpy := spy.NewAck()
		teleportMQ := newDvpSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{Number: 100, Logs: []ethTypes.Log{log}},
			Ack: ackSpy.Fn(),
		})

		swapRepo := &DvpSwapRepositoryMock{
			GetSwapBySharedIDFunc: func(ctx context.Context, sharedID string) (*types.DvpSwap, error) {
				return &types.DvpSwap{SharedID: sharedID, To: "", SelfSalt: big.NewInt(1)}, nil
			},
		}
		ctsMock := &DvpCTSClientMock{
			DecryptWithoutFPWithSSFunc: func(ctx context.Context, in *encryptpb.DecryptWithoutFPWithSSRequest, opts ...grpc.CallOption) (*encryptpb.DecryptWithoutFPWithSSResponse, error) {
				return nil, errors.New("cts down")
			},
		}
		orchestratorMQ := &DvpOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.DvpDestMessage) error {
				assert.Fail(t, "should not push when decrypt fails")
				return nil
			},
		}

		parser := logparser.NewDvpTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			ctsMock,
			newDefaultDvpEthClientMock(),
			swapRepo,
			big.NewInt(12345),
			newDefaultDvpBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		require.NoError(t, parser.Run(ctx))

		assert.Equal(t, 1, len(ctsMock.DecryptWithoutFPWithSSCalls()))
		assert.Empty(t, swapRepo.UpdateSwapToCalls())
		assert.Empty(t, orchestratorMQ.PushCalls())
		ackSpy.AssertNotCalled(t, "should not ack block when decrypt fails, allowing MQ redelivery")
	})

	t.Run("does not ack SwapCompleted when decrypted payload is malformed", func(t *testing.T) {
		log := testdata.NewSwapCompletedLogWith()

		ackSpy := spy.NewAck()
		teleportMQ := newDvpSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{Number: 100, Logs: []ethTypes.Log{log}},
			Ack: ackSpy.Fn(),
		})

		swapRepo := &DvpSwapRepositoryMock{
			GetSwapBySharedIDFunc: func(ctx context.Context, sharedID string) (*types.DvpSwap, error) {
				return &types.DvpSwap{SharedID: sharedID, To: "", SelfSalt: big.NewInt(1)}, nil
			},
		}
		ctsMock := &DvpCTSClientMock{
			DecryptWithoutFPWithSSFunc: func(ctx context.Context, in *encryptpb.DecryptWithoutFPWithSSRequest, opts ...grpc.CallOption) (*encryptpb.DecryptWithoutFPWithSSResponse, error) {
				return &encryptpb.DecryptWithoutFPWithSSResponse{
					Outcome:   encryptpb.DecryptOutcome_OUTCOME_OK,
					Plaintext: []byte("not json"),
				}, nil
			},
		}
		orchestratorMQ := &DvpOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.DvpDestMessage) error {
				assert.Fail(t, "should not push when decrypted payload is malformed")
				return nil
			},
		}

		parser := logparser.NewDvpTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			ctsMock,
			newDefaultDvpEthClientMock(),
			swapRepo,
			big.NewInt(12345),
			newDefaultDvpBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		require.NoError(t, parser.Run(ctx))

		assert.Equal(t, 1, len(ctsMock.DecryptWithoutFPWithSSCalls()))
		assert.Empty(t, swapRepo.UpdateSwapToCalls())
		assert.Empty(t, orchestratorMQ.PushCalls())
		ackSpy.AssertNotCalled(t, "should not ack block when payload is malformed, allowing MQ redelivery")
	})

	t.Run("does not ack SwapCompleted when decrypted To address is invalid", func(t *testing.T) {
		cases := []struct {
			name string
			to   string
		}{
			{name: "empty", to: ""},
			{name: "not hex", to: "not-an-address"},
			{name: "wrong length", to: "0xabcdef"},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				log := testdata.NewSwapCompletedLogWith()

				ackSpy := spy.NewAck()
				teleportMQ := newDvpSingleBlockMQ(msgqueue.Message[logrouter.Block]{
					V: logrouter.Block{Number: 100, Logs: []ethTypes.Log{log}},
					Ack: ackSpy.Fn(),
				})

				swapRepo := &DvpSwapRepositoryMock{
					GetSwapBySharedIDFunc: func(ctx context.Context, sharedID string) (*types.DvpSwap, error) {
						return &types.DvpSwap{SharedID: sharedID, To: "", SelfSalt: big.NewInt(1)}, nil
					},
					UpdateSwapToFunc: func(ctx context.Context, sharedID string, to string) error {
						assert.Fail(t, "should not persist when decrypted To address is invalid")
						return nil
					},
				}
				plaintext, err := json.Marshal(&types.DvpSwapMessage{To: tc.to})
				require.NoError(t, err)
				ctsMock := &DvpCTSClientMock{
					DecryptWithoutFPWithSSFunc: func(ctx context.Context, in *encryptpb.DecryptWithoutFPWithSSRequest, opts ...grpc.CallOption) (*encryptpb.DecryptWithoutFPWithSSResponse, error) {
						return &encryptpb.DecryptWithoutFPWithSSResponse{
							Outcome:   encryptpb.DecryptOutcome_OUTCOME_OK,
							Plaintext: plaintext,
						}, nil
					},
				}
				orchestratorMQ := &DvpOrchestratorMQMock{
					PushFunc: func(ctx context.Context, msg service.DvpDestMessage) error {
						assert.Fail(t, "should not push when decrypted To address is invalid")
						return nil
					},
				}

				parser := logparser.NewDvpTeleportParserWithBackoff(
					teleportMQ,
					orchestratorMQ,
					ctsMock,
					newDefaultDvpEthClientMock(),
					swapRepo,
					big.NewInt(12345),
					newDefaultDvpBackoffMock(),
				)

				ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
				defer cancel()

				require.NoError(t, parser.Run(ctx))

				assert.Equal(t, 1, len(ctsMock.DecryptWithoutFPWithSSCalls()))
				assert.Empty(t, swapRepo.UpdateSwapToCalls())
				assert.Empty(t, orchestratorMQ.PushCalls())
				ackSpy.AssertNotCalled(t, "should not ack block when To address is invalid, allowing MQ redelivery")
			})
		}
	})

	t.Run("does not ack SwapCompleted when UpdateSwapTo fails", func(t *testing.T) {
		log := testdata.NewSwapCompletedLogWith()

		ackSpy := spy.NewAck()
		teleportMQ := newDvpSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{Number: 100, Logs: []ethTypes.Log{log}},
			Ack: ackSpy.Fn(),
		})

		swapRepo := &DvpSwapRepositoryMock{
			GetSwapBySharedIDFunc: func(ctx context.Context, sharedID string) (*types.DvpSwap, error) {
				return &types.DvpSwap{SharedID: sharedID, To: "", SelfSalt: big.NewInt(1)}, nil
			},
			UpdateSwapToFunc: func(ctx context.Context, sharedID string, to string) error {
				return errors.New("db err")
			},
		}
		plaintext, err := json.Marshal(&types.DvpSwapMessage{To: "0x1234567890AbcdEF1234567890aBcdef12345678"})
		require.NoError(t, err)
		ctsMock := &DvpCTSClientMock{
			DecryptWithoutFPWithSSFunc: func(ctx context.Context, in *encryptpb.DecryptWithoutFPWithSSRequest, opts ...grpc.CallOption) (*encryptpb.DecryptWithoutFPWithSSResponse, error) {
				return &encryptpb.DecryptWithoutFPWithSSResponse{
					Outcome:   encryptpb.DecryptOutcome_OUTCOME_OK,
					Plaintext: plaintext,
				}, nil
			},
		}
		orchestratorMQ := &DvpOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.DvpDestMessage) error {
				assert.Fail(t, "should not push when UpdateSwapTo fails")
				return nil
			},
		}

		parser := logparser.NewDvpTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			ctsMock,
			newDefaultDvpEthClientMock(),
			swapRepo,
			big.NewInt(12345),
			newDefaultDvpBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		require.NoError(t, parser.Run(ctx))

		assert.Equal(t, 1, len(ctsMock.DecryptWithoutFPWithSSCalls()))
		assert.Equal(t, 1, len(swapRepo.UpdateSwapToCalls()))
		assert.Empty(t, orchestratorMQ.PushCalls())
		ackSpy.AssertNotCalled(t, "should not ack block when UpdateSwapTo fails, allowing MQ redelivery")
	})

	t.Run("does not ack SwapCompleted on unexpected decrypt outcome", func(t *testing.T) {
		log := testdata.NewSwapCompletedLogWith()

		ackSpy := spy.NewAck()
		teleportMQ := newDvpSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{Number: 100, Logs: []ethTypes.Log{log}},
			Ack: ackSpy.Fn(),
		})

		swapRepo := &DvpSwapRepositoryMock{
			GetSwapBySharedIDFunc: func(ctx context.Context, sharedID string) (*types.DvpSwap, error) {
				return &types.DvpSwap{SharedID: sharedID, To: "", SelfSalt: big.NewInt(1)}, nil
			},
		}
		ctsMock := &DvpCTSClientMock{
			DecryptWithoutFPWithSSFunc: func(ctx context.Context, in *encryptpb.DecryptWithoutFPWithSSRequest, opts ...grpc.CallOption) (*encryptpb.DecryptWithoutFPWithSSResponse, error) {
				return &encryptpb.DecryptWithoutFPWithSSResponse{}, nil
			},
		}
		orchestratorMQ := &DvpOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.DvpDestMessage) error {
				assert.Fail(t, "should not push on unexpected decrypt outcome")
				return nil
			},
		}

		parser := logparser.NewDvpTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			ctsMock,
			newDefaultDvpEthClientMock(),
			swapRepo,
			big.NewInt(12345),
			newDefaultDvpBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		require.NoError(t, parser.Run(ctx))

		assert.Equal(t, 1, len(ctsMock.DecryptWithoutFPWithSSCalls()))
		assert.Empty(t, swapRepo.UpdateSwapToCalls())
		assert.Empty(t, orchestratorMQ.PushCalls())
		ackSpy.AssertNotCalled(t, "should not ack block on unexpected decrypt outcome, allowing MQ redelivery")
	})
}
