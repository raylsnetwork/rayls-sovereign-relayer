// Decommissioning Teleport (vanilla, atomic).

package service_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	sharedservice "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	os.Exit(m.Run())
}

type stubEarlyRevertTransactionRepo struct {
	calls []earlyRevertTxRepoCall
	err   error
}

type earlyRevertTxRepoCall struct {
	sharedIDs []string
	state     types.TransactionState
	outcome   types.TransactionOutcome
}

func (r *stubEarlyRevertTransactionRepo) BatchSetState(
	_ context.Context,
	sharedIDs []string,
	state types.TransactionState,
) error {
	r.calls = append(r.calls, earlyRevertTxRepoCall{sharedIDs: sharedIDs, state: state, outcome: types.OutcomePending})
	return r.err
}

func (r *stubEarlyRevertTransactionRepo) BatchSetOutcome(
	_ context.Context,
	sharedIDs []string,
	outcome types.TransactionOutcome,
) error {
	r.calls = append(r.calls, earlyRevertTxRepoCall{sharedIDs: sharedIDs, outcome: outcome})
	return r.err
}

type stubEarlyRevertSignatureBatchSender struct {
	spyMessages []types.TxRequest
	err         error
}

func (s *stubEarlyRevertSignatureBatchSender) PushBatch(_ context.Context, msgs []types.TxRequest) error {
	s.spyMessages = msgs
	return s.err
}

type stubEarlyRevertSignatureRepository struct {
	err                      error
	signatures               []types.CalldataSignature
	spySourceRevertSharedIDs []string
}

func (s *stubEarlyRevertSignatureRepository) GetSourceRevertForSharedIDs(
	_ context.Context,
	sharedIDs []string,
) ([]types.CalldataSignature, error) {
	s.spySourceRevertSharedIDs = sharedIDs
	return s.signatures, s.err
}

type stubEarlyRevertGenerator struct {
	calldata []byte
	err      error
}

func (g *stubEarlyRevertGenerator) Generate(_ types.CalldataSignature) ([]byte, error) {
	return g.calldata, g.err
}

func TestEarlyRevertService_HandleEarlyRevert(t *testing.T) {
	destEndpoint := common.HexToAddress("0xdeadbeef")

	t.Run("publishes a TxRequest per signature and marks state EarlyRevertSigs", func(t *testing.T) {
		sharedIDs := []string{"shared-1", "shared-2"}
		signatures := []types.CalldataSignature{
			{SharedId: "shared-1"},
			{SharedId: "shared-2"},
		}
		calldata := []byte{0xCA, 0xFE}

		batcher := &stubEarlyRevertSignatureBatchSender{}
		gen := &stubEarlyRevertGenerator{calldata: calldata}
		sigRepo := &stubEarlyRevertSignatureRepository{signatures: signatures}
		txRepo := &stubEarlyRevertTransactionRepo{}

		svc := service.NewEarlyRevertService(batcher, gen, destEndpoint, sigRepo, txRepo)

		require.NoError(t, svc.HandleEarlyRevert(context.Background(), sharedIDs))

		assert.Equal(t, sharedIDs, sigRepo.spySourceRevertSharedIDs)
		require.Len(t, batcher.spyMessages, 2)
		for i, msg := range batcher.spyMessages {
			assert.Equal(t, signatures[i].SharedId, msg.CorrelationID)
			assert.Equal(t, "atomic.early-revert", msg.MessageType)
			assert.Equal(t, destEndpoint, msg.Address)
			assert.Equal(t, calldata, msg.Calldata)
		}
		require.Len(t, txRepo.calls, 1)
		assert.Equal(t, sharedIDs, txRepo.calls[0].sharedIDs)
		assert.Equal(t, types.EarlyRevertSigs, txRepo.calls[0].state)
	})

	t.Run("skips messages whose calldata generation fails but still publishes the rest", func(t *testing.T) {
		sharedIDs := []string{"shared-only"}
		signatures := []types.CalldataSignature{{SharedId: "shared-only"}}

		batcher := &stubEarlyRevertSignatureBatchSender{}
		gen := &stubEarlyRevertGenerator{err: errors.New("generator boom")}
		sigRepo := &stubEarlyRevertSignatureRepository{signatures: signatures}
		txRepo := &stubEarlyRevertTransactionRepo{}

		svc := service.NewEarlyRevertService(batcher, gen, destEndpoint, sigRepo, txRepo)
		require.NoError(t, svc.HandleEarlyRevert(context.Background(), sharedIDs))

		assert.Empty(t, batcher.spyMessages, "skipped signature should not be published")
	})

	t.Run("wraps errors in AtomicServiceError", func(t *testing.T) {
		boom := errors.New("boom")
		var wantErrType *sharedservice.AtomicServiceError

		cases := []struct {
			name       string
			sigRepoErr error
			batcherErr error
			txRepoErr  error
		}{
			{name: "signature repo error", sigRepoErr: boom},
			{name: "batcher publish error", batcherErr: boom},
			{name: "tx repo error", txRepoErr: boom},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				batcher := &stubEarlyRevertSignatureBatchSender{err: tc.batcherErr}
				gen := &stubEarlyRevertGenerator{calldata: []byte{0x01}}
				sigRepo := &stubEarlyRevertSignatureRepository{
					signatures: []types.CalldataSignature{{SharedId: "x"}},
					err:        tc.sigRepoErr,
				}
				txRepo := &stubEarlyRevertTransactionRepo{err: tc.txRepoErr}

				svc := service.NewEarlyRevertService(batcher, gen, destEndpoint, sigRepo, txRepo)
				err := svc.HandleEarlyRevert(context.Background(), []string{"x"})

				require.Error(t, err)
				require.ErrorAs(t, err, &wantErrType)
				require.ErrorIs(t, err, boom)
			})
		}
	})
}

func TestEarlyRevertService_HandleEarlyRevertCallback(t *testing.T) {
	destEndpoint := common.HexToAddress("0xdeadbeef")

	newSvc := func(txRepo *stubEarlyRevertTransactionRepo) *service.EarlyRevertService {
		return service.NewEarlyRevertService(
			&stubEarlyRevertSignatureBatchSender{},
			&stubEarlyRevertGenerator{},
			destEndpoint,
			&stubEarlyRevertSignatureRepository{},
			txRepo,
		)
	}

	t.Run("partitions results by Kind and updates outcome per group", func(t *testing.T) {
		results := []types.TxResult{
			{CorrelationID: "ok-1", Kind: types.TxResultSuccess},
			{CorrelationID: "ok-2", Kind: types.TxResultSuccess},
			{CorrelationID: "rev", Kind: types.TxResultRevert, RevertData: []byte{0xAB}},
			{CorrelationID: "fail", Kind: types.TxResultFailed, ErrorReason: "stuck"},
		}

		txRepo := &stubEarlyRevertTransactionRepo{}
		svc := newSvc(txRepo)
		require.NoError(t, svc.HandleEarlyRevertCallback(context.Background(), results))

		require.Len(t, txRepo.calls, 3)
		callsByOutcome := map[types.TransactionOutcome][]string{}
		for _, c := range txRepo.calls {
			callsByOutcome[c.outcome] = c.sharedIDs
		}
		assert.ElementsMatch(t, []string{"ok-1", "ok-2"}, callsByOutcome[types.OutcomeSuccess])
		assert.ElementsMatch(t, []string{"rev"}, callsByOutcome[types.OutcomeReverted])
		assert.ElementsMatch(t, []string{"fail"}, callsByOutcome[types.OutcomeFailed])
	})

	t.Run("does nothing on empty input", func(t *testing.T) {
		txRepo := &stubEarlyRevertTransactionRepo{}
		svc := newSvc(txRepo)
		require.NoError(t, svc.HandleEarlyRevertCallback(context.Background(), nil))
		assert.Empty(t, txRepo.calls)
	})

	t.Run("propagates AtomicServiceError on repo failure", func(t *testing.T) {
		txRepo := &stubEarlyRevertTransactionRepo{err: errors.New("db down")}
		svc := newSvc(txRepo)

		err := svc.HandleEarlyRevertCallback(context.Background(), []types.TxResult{
			{CorrelationID: "ok", Kind: types.TxResultSuccess},
		})

		var wantErrType *sharedservice.AtomicServiceError
		require.ErrorAs(t, err, &wantErrType)
	})
}
