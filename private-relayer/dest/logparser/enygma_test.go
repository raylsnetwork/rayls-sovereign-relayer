package logparser_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EnygmaTeleport"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/encrypt"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/wireformat"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/logparser"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/logparser/testdata"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/logrouter"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools/fake"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools/spy"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func newSingleBlockMQ(block msgqueue.Message[logrouter.Block]) *EnygmaTeleportMQMock {
	return &EnygmaTeleportMQMock{
		NextFunc: fake.NextMQ(block),
	}
}

func newDefaultKOSMock(t *testing.T) *KOSDecryptorMock {
	t.Helper()
	batch := types.EnygmaTransferBatch{
		ResourceId:  "0102030000000000000000000000000000000000000000000000000000000000",
		FromChainID: big.NewInt(999),
	}
	plaintext, err := wireformat.MarshalPlaintext(batch)
	require.NoError(t, err)

	return &KOSDecryptorMock{
		DecryptEnygmaTransferBatchFunc: func(ctx context.Context, in *encrypt.DecryptEnygmaTransferBatchRequest, opts ...grpc.CallOption) (*encrypt.DecryptEnygmaTransferBatchResponse, error) {
			return &encrypt.DecryptEnygmaTransferBatchResponse{Plaintext: plaintext, Outcome: encrypt.DecryptOutcome_OUTCOME_OK}, nil
		},
	}
}

func newDefaultBackoffMock() *EnygmaBackoffMock {
	return &EnygmaBackoffMock{
		DoFunc: func(ctx context.Context, maxAttempts int, fn func() error) error {
			return fn()
		},
	}
}

func TestEnygmaTeleportParser(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("supports graceful shutdown", func(t *testing.T) {
		teleportMQ := &EnygmaTeleportMQMock{
			NextFunc: func(ctx context.Context) (msgqueue.Message[logrouter.Block], error) {
				<-ctx.Done()
				return msgqueue.Message[logrouter.Block]{}, context.Canceled
			},
		}
		orchestratorMQ := &EnygmaOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.EnygmaDestMessage) error {
				return nil
			},
		}

		parser := logparser.NewEnygmaTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			newDefaultKOSMock(t),
			big.NewInt(12345),
			newDefaultBackoffMock(),
		)

		hasGracefulShutdown := testtools.ShutdownFixture(t, parser.Run, time.Second)
		assert.True(t, hasGracefulShutdown)
	})

	t.Run("acks block after processing logs with no matching events", func(t *testing.T) {
		ackSpy := spy.NewAck()

		teleportMQ := newSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{
				Number: 100,
				Logs:   []ethTypes.Log{},
			},
			Ack: ackSpy.Fn(),
		})
		orchestratorMQ := &EnygmaOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.EnygmaDestMessage) error {
				assert.Fail(t, "should not push any messages for empty logs")
				return nil
			},
		}

		parser := logparser.NewEnygmaTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			newDefaultKOSMock(t),
			big.NewInt(12345),
			newDefaultBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := parser.Run(ctx)
		require.NoError(t, err)

		ackSpy.AssertCalled(t)
		assert.Equal(t, 0, len(orchestratorMQ.PushCalls()))
	})

	t.Run("constructor creates parser with valid teleport filterer", func(t *testing.T) {
		teleportMQ := &EnygmaTeleportMQMock{}
		orchestratorMQ := &EnygmaOrchestratorMQMock{}

		parser := logparser.NewEnygmaTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			newDefaultKOSMock(t),
			big.NewInt(12345),
			newDefaultBackoffMock(),
		)

		assert.NotNil(t, parser)
	})

	t.Run("parses EnygmaTransfer and pushes to orchestrator", func(t *testing.T) {
		localChainId := big.NewInt(12345)
		wantBlockNumber := uint64(100)
		wantResourceID := "0102030000000000000000000000000000000000000000000000000000000000"

		log := testdata.NewEnygmaTransferLogWith(
			testdata.WithTransferBlockNumber(wantBlockNumber),
			testdata.WithTransferToChainId(localChainId),
		)

		ackSpy := spy.NewAck()
		teleportMQ := newSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{
				Number: wantBlockNumber,
				Logs:   []ethTypes.Log{log},
			},
			Ack: ackSpy.Fn(),
		})

		var pushedMsg service.EnygmaDestMessage
		orchestratorMQ := &EnygmaOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.EnygmaDestMessage) error {
				pushedMsg = msg
				return nil
			},
		}

		batch := types.EnygmaTransferBatch{
			ResourceId:  wantResourceID,
			FromChainID: big.NewInt(999),
		}
		batchPlaintext, err := wireformat.MarshalPlaintext(batch)
		require.NoError(t, err)

		kosClient := &KOSDecryptorMock{
			DecryptEnygmaTransferBatchFunc: func(ctx context.Context, in *encrypt.DecryptEnygmaTransferBatchRequest, opts ...grpc.CallOption) (*encrypt.DecryptEnygmaTransferBatchResponse, error) {
				assert.Equal(t, wantBlockNumber, in.BlockNumber)
				return &encrypt.DecryptEnygmaTransferBatchResponse{Plaintext: batchPlaintext, Outcome: encrypt.DecryptOutcome_OUTCOME_OK}, nil
			},
		}

		parser := logparser.NewEnygmaTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			kosClient,
			localChainId,
			newDefaultBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err = parser.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(orchestratorMQ.PushCalls()))
		assert.Equal(t, service.EnygmaTransferBatchMessage, pushedMsg.Type)
		assert.Equal(t, wantBlockNumber, pushedMsg.BlockNumber)
		assert.NotNil(t, pushedMsg.TransferBatch)
		assert.Equal(t, wantResourceID, pushedMsg.TransferBatch.ResourceId)
		ackSpy.AssertCalled(t)
	})

	t.Run("skips EnygmaTransfer not destined for local chain", func(t *testing.T) {
		localChainId := big.NewInt(12345)
		otherChainId := big.NewInt(99999)

		log := testdata.NewEnygmaTransferLogWith(
			testdata.WithTransferToChainId(otherChainId),
		)

		ackSpy := spy.NewAck()
		teleportMQ := newSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{
				Number: 100,
				Logs:   []ethTypes.Log{log},
			},
			Ack: ackSpy.Fn(),
		})

		orchestratorMQ := &EnygmaOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.EnygmaDestMessage) error {
				assert.Fail(t, "should not push message for other chain")
				return nil
			},
		}

		parser := logparser.NewEnygmaTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			newDefaultKOSMock(t),
			localChainId,
			newDefaultBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := parser.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 0, len(orchestratorMQ.PushCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("acks and skips when batch is not addressed to us (OUTCOME_NOT_FOR_RECIPIENT)", func(t *testing.T) {
		localChainId := big.NewInt(12345)

		log := testdata.NewEnygmaTransferLogWith(
			testdata.WithTransferToChainId(localChainId),
		)

		ackSpy := spy.NewAck()
		teleportMQ := newSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{
				Number: 100,
				Logs:   []ethTypes.Log{log},
			},
			Ack: ackSpy.Fn(),
		})

		orchestratorMQ := &EnygmaOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.EnygmaDestMessage) error {
				assert.Fail(t, "should not push message when decryption fails")
				return nil
			},
		}

		kosClient := &KOSDecryptorMock{
			DecryptEnygmaTransferBatchFunc: func(ctx context.Context, in *encrypt.DecryptEnygmaTransferBatchRequest, opts ...grpc.CallOption) (*encrypt.DecryptEnygmaTransferBatchResponse, error) {
				return &encrypt.DecryptEnygmaTransferBatchResponse{Outcome: encrypt.DecryptOutcome_OUTCOME_NOT_FOR_RECIPIENT}, nil
			},
		}

		parser := logparser.NewEnygmaTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			kosClient,
			localChainId,
			newDefaultBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := parser.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 0, len(orchestratorMQ.PushCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("does not ack on generic decryption error so block is retried", func(t *testing.T) {
		localChainId := big.NewInt(12345)

		log := testdata.NewEnygmaTransferLogWith(
			testdata.WithTransferToChainId(localChainId),
		)

		ackSpy := spy.NewAck()
		teleportMQ := newSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{
				Number: 100,
				Logs:   []ethTypes.Log{log},
			},
			Ack: ackSpy.Fn(),
		})

		orchestratorMQ := &EnygmaOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.EnygmaDestMessage) error {
				assert.Fail(t, "should not push message when decryption fails")
				return nil
			},
		}

		kosClient := &KOSDecryptorMock{
			DecryptEnygmaTransferBatchFunc: func(ctx context.Context, in *encrypt.DecryptEnygmaTransferBatchRequest, opts ...grpc.CallOption) (*encrypt.DecryptEnygmaTransferBatchResponse, error) {
				return nil, errors.New("decryption failed")
			},
		}

		parser := logparser.NewEnygmaTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			kosClient,
			localChainId,
			newDefaultBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := parser.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 0, len(orchestratorMQ.PushCalls()))
		ackSpy.AssertNotCalled(t, "expected ack to be skipped so block is redelivered")
	})

	t.Run("parses BalancesFinalized and pushes to orchestrator", func(t *testing.T) {
		localChainId := big.NewInt(12345)
		wantBlockNumber := uint64(100)
		wantBalanceX := big.NewInt(1000)
		wantBalanceY := big.NewInt(2000)

		log := testdata.NewBalancesFinalizedLogWith(
			testdata.WithFinalizedLogBlockNumber(wantBlockNumber),
			testdata.WithFinalizedBalances([]EnygmaTeleport.IEnygmaV1EnygmaPointWithChainId{
				{
					ChainId: localChainId,
					C1:      wantBalanceX,
					C2:      wantBalanceY,
				},
			}),
		)

		ackSpy := spy.NewAck()
		teleportMQ := newSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{
				Number: wantBlockNumber,
				Logs:   []ethTypes.Log{log},
			},
			Ack: ackSpy.Fn(),
		})

		var pushedMsg service.EnygmaDestMessage
		orchestratorMQ := &EnygmaOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.EnygmaDestMessage) error {
				pushedMsg = msg
				return nil
			},
		}

		parser := logparser.NewEnygmaTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			newDefaultKOSMock(t),
			localChainId,
			newDefaultBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := parser.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(orchestratorMQ.PushCalls()))
		assert.Equal(t, service.EnygmaFinalizedBalanceMessage, pushedMsg.Type)
		assert.NotNil(t, pushedMsg.FinalizedBalance)
		assert.Equal(t, wantBalanceX, pushedMsg.FinalizedBalance.BalanceX)
		assert.Equal(t, wantBalanceY, pushedMsg.FinalizedBalance.BalanceY)
		ackSpy.AssertCalled(t)
	})

	// Regression for the Enygma SC<->DB finalization deadlock.
	// When a single commit-chain transaction finalizes more than one settlement window,
	// EnygmaV1 emits multiple BalancesFinalized events in the SAME block. The orchestrator MQ
	// dedups on EnygmaDestMessage.GetID() (jetstream.WithMsgID + 1h Duplicates window). If the
	// message ID is keyed on the emit block alone, those distinct finalizations collide on one ID,
	// the MQ silently drops the second checkpoint message, the relayer DB falls one finalization
	// window behind the chain, and every subsequent transfer proof fails the SC-vs-DB commitment
	// guard forever (proof.go: "balance commitment mismatch between SC and DB"). Each distinct
	// finalization (finalizedBlockNumber, pendingBlockNumber) MUST therefore produce a distinct ID.
	t.Run("distinct finalizations in the same block get distinct message IDs (no MQ dedup collision)", func(t *testing.T) {
		localChainId := big.NewInt(12345)
		emitBlock := uint64(15559)
		var resourceID [32]byte
		copy(resourceID[:], []byte{0x5b, 0xfa, 0x74, 0xc7, 0x43, 0x91, 0x40, 0x28})

		balances := []EnygmaTeleport.IEnygmaV1EnygmaPointWithChainId{
			{ChainId: localChainId, C1: big.NewInt(1000), C2: big.NewInt(2000)},
		}
		// Two finalizations emitted in the SAME commit-chain block for DIFFERENT settlement windows.
		logFinalize15554 := testdata.NewBalancesFinalizedLogWith(
			testdata.WithFinalizedResourceID(resourceID),
			testdata.WithFinalizedLogBlockNumber(emitBlock),
			testdata.WithFinalizedBlockNumbers(big.NewInt(15554), big.NewInt(15556)),
			testdata.WithFinalizedBalances(balances),
		)
		logFinalize15556 := testdata.NewBalancesFinalizedLogWith(
			testdata.WithFinalizedResourceID(resourceID),
			testdata.WithFinalizedLogBlockNumber(emitBlock),
			testdata.WithFinalizedBlockNumbers(big.NewInt(15556), big.NewInt(15558)),
			testdata.WithFinalizedBalances(balances),
		)

		ackSpy := spy.NewAck()
		teleportMQ := newSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{
				Number: emitBlock,
				Logs:   []ethTypes.Log{logFinalize15554, logFinalize15556},
			},
			Ack: ackSpy.Fn(),
		})

		var pushedIDs []string
		orchestratorMQ := &EnygmaOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.EnygmaDestMessage) error {
				pushedIDs = append(pushedIDs, msg.GetID())
				return nil
			},
		}

		parser := logparser.NewEnygmaTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			newDefaultKOSMock(t),
			localChainId,
			newDefaultBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := parser.Run(ctx)
		require.NoError(t, err)

		require.Len(t, pushedIDs, 2, "both finalizations in the block must be pushed to the orchestrator")
		assert.NotEqual(t, pushedIDs[0], pushedIDs[1],
			"two distinct finalizations emitted in the same block must produce distinct MQ message IDs; "+
				"otherwise jetstream.WithMsgID dedup drops the second checkpoint and the relayer DB deadlocks "+
				"one finalization window behind the chain")
	})

	t.Run("skips BalancesFinalized for other chains", func(t *testing.T) {
		localChainId := big.NewInt(12345)
		otherChainId := big.NewInt(99999)

		log := testdata.NewBalancesFinalizedLogWith(
			testdata.WithFinalizedBalances([]EnygmaTeleport.IEnygmaV1EnygmaPointWithChainId{
				{
					ChainId: otherChainId,
					C1:      big.NewInt(1000),
					C2:      big.NewInt(2000),
				},
			}),
		)

		ackSpy := spy.NewAck()
		teleportMQ := newSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{
				Number: 100,
				Logs:   []ethTypes.Log{log},
			},
			Ack: ackSpy.Fn(),
		})

		orchestratorMQ := &EnygmaOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.EnygmaDestMessage) error {
				assert.Fail(t, "should not push message for other chain")
				return nil
			},
		}

		parser := logparser.NewEnygmaTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			newDefaultKOSMock(t),
			localChainId,
			newDefaultBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := parser.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 0, len(orchestratorMQ.PushCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("does not ack block when orchestrator push fails for transfer", func(t *testing.T) {
		localChainId := big.NewInt(12345)

		log := testdata.NewEnygmaTransferLogWith(
			testdata.WithTransferToChainId(localChainId),
		)

		ackSpy := spy.NewAck()
		teleportMQ := newSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{
				Number: 100,
				Logs:   []ethTypes.Log{log},
			},
			Ack: ackSpy.Fn(),
		})

		orchestratorMQ := &EnygmaOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.EnygmaDestMessage) error {
				return errors.New("push failed")
			},
		}

		parser := logparser.NewEnygmaTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			newDefaultKOSMock(t),
			localChainId,
			newDefaultBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := parser.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(orchestratorMQ.PushCalls()))
		ackSpy.AssertNotCalled(t, "should not ack block when push fails, allowing MQ redelivery")
	})

	t.Run("does not ack block when orchestrator push fails for BalancesFinalized", func(t *testing.T) {
		localChainId := big.NewInt(12345)

		log := testdata.NewBalancesFinalizedLogWith(
			testdata.WithFinalizedBalances([]EnygmaTeleport.IEnygmaV1EnygmaPointWithChainId{
				{
					ChainId: localChainId,
					C1:      big.NewInt(1000),
					C2:      big.NewInt(2000),
				},
			}),
		)

		ackSpy := spy.NewAck()
		teleportMQ := newSingleBlockMQ(msgqueue.Message[logrouter.Block]{
			V: logrouter.Block{
				Number: 100,
				Logs:   []ethTypes.Log{log},
			},
			Ack: ackSpy.Fn(),
		})

		orchestratorMQ := &EnygmaOrchestratorMQMock{
			PushFunc: func(ctx context.Context, msg service.EnygmaDestMessage) error {
				return errors.New("push failed")
			},
		}

		parser := logparser.NewEnygmaTeleportParserWithBackoff(
			teleportMQ,
			orchestratorMQ,
			newDefaultKOSMock(t),
			localChainId,
			newDefaultBackoffMock(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := parser.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(orchestratorMQ.PushCalls()))
		ackSpy.AssertNotCalled(t, "should not ack block when BalancesFinalized push fails, allowing MQ redelivery")
	})
}

// captureLogs installs a thread-safe text slog handler as the default for the duration of the test and
// returns a function yielding everything logged so far. It restores the previous default on cleanup.
func captureLogs(t *testing.T) func() string {
	t.Helper()
	var mu sync.Mutex
	buf := &bytes.Buffer{}
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&syncWriter{mu: &mu, w: buf}, &slog.HandlerOptions{Level: slog.LevelDebug})))
	t.Cleanup(func() { slog.SetDefault(prev) })
	return func() string {
		mu.Lock()
		defer mu.Unlock()
		return buf.String()
	}
}

type syncWriter struct {
	mu *sync.Mutex
	w  *bytes.Buffer
}

func (s *syncWriter) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.w.Write(p)
}

// Regression for the AEAD-TAMPERED silent-drop data-loss bug (#229).
//
// A tag-matched EnygmaTransfer batch that fails AEAD verification on the destination used to be dropped
// with a WARN + return nil: no structured event, no severity, no forensic trail — indistinguishable from
// routine noise. The fix keeps the head-of-line-safe behaviour (the batch is NOT pushed and the block is
// still acked so one poison block cannot stall the shared dest stream), but makes the drop LOUD: a
// structured ERROR event=enygma_transfer_aead_tampered carrying the resourceId, the dest/seal block
// numbers and the tx hash, so operators can alert on and reconcile it.
func TestEnygmaTeleportParser_TamperedBatchIsSurfacedNotSilentlyDropped(t *testing.T) {
	logs := captureLogs(t)

	localChainId := big.NewInt(12345)
	wantBlockNumber := uint64(12118)

	var resourceID [32]byte
	copy(resourceID[:], []byte{0x01, 0x02, 0x03})

	log := testdata.NewEnygmaTransferLogWith(
		testdata.WithTransferBlockNumber(wantBlockNumber),
		testdata.WithTransferToChainId(localChainId),
		testdata.WithTransferResourceID(resourceID),
	)

	ackSpy := spy.NewAck()
	teleportMQ := newSingleBlockMQ(msgqueue.Message[logrouter.Block]{
		V: logrouter.Block{
			Number: wantBlockNumber,
			Logs:   []ethTypes.Log{log},
		},
		Ack: ackSpy.Fn(),
	})

	orchestratorMQ := &EnygmaOrchestratorMQMock{
		PushFunc: func(ctx context.Context, msg service.EnygmaDestMessage) error {
			assert.Fail(t, "tampered batch must NOT be pushed to the orchestrator")
			return nil
		},
	}

	// KOS returns OUTCOME_TAMPERED: the message tag matched but AEAD verification failed.
	kosClient := &KOSDecryptorMock{
		DecryptEnygmaTransferBatchFunc: func(ctx context.Context, in *encrypt.DecryptEnygmaTransferBatchRequest, opts ...grpc.CallOption) (*encrypt.DecryptEnygmaTransferBatchResponse, error) {
			return &encrypt.DecryptEnygmaTransferBatchResponse{Outcome: encrypt.DecryptOutcome_OUTCOME_TAMPERED}, nil
		},
	}

	parser := logparser.NewEnygmaTeleportParserWithBackoff(
		teleportMQ,
		orchestratorMQ,
		kosClient,
		localChainId,
		newDefaultBackoffMock(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	require.NoError(t, parser.Run(ctx))

	// Behavioural guarantees: not pushed (the batch is dropped at the parser), but the block IS still
	// acked — head-of-line safety on the shared dest stream is preserved (a poison block must not stall
	// every other token's transfers).
	assert.Empty(t, orchestratorMQ.PushCalls(), "tampered batch must not be pushed to the orchestrator")
	ackSpy.AssertCalled(t)

	// The drop must be SURFACED, not silent: a structured ERROR event with the forensic context.
	out := logs()
	assert.Contains(t, out, "event=enygma_transfer_aead_tampered",
		"the tampered drop must emit the greppable structured event for alerting")
	assert.Contains(t, out, "level=ERROR",
		"the tampered drop must be ERROR severity, not a buried WARN")
	assert.Contains(t, out, strings.ToLower("010203"),
		"the structured event must carry the resourceId for forensic reconciliation")
}
