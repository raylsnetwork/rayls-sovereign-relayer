// Decommissioning Teleport (vanilla, atomic).

package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	service2 "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/txsim"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Stubs
// ---------------------------------------------------------------------

type StubSignatureSignatureRepository struct {
	spyDestinationUnlockSharedIDs []string
	spyDestinationRevertSharedIDs []string
	spySourceRevertSharedIDs      []string

	signatures []types.CalldataSignature
	err        error
}

func (r *StubSignatureSignatureRepository) GetDestinationUnlocksForSharedIDs(
	ctx context.Context, sharedIDs []string,
) ([]types.CalldataSignature, error) {
	r.spyDestinationUnlockSharedIDs = sharedIDs
	return r.signatures, r.err
}

func (r *StubSignatureSignatureRepository) GetDestinationRevertsForSharedIDs(
	ctx context.Context, sharedIDs []string,
) ([]types.CalldataSignature, error) {
	r.spyDestinationRevertSharedIDs = sharedIDs
	return r.signatures, r.err
}

func (r *StubSignatureSignatureRepository) GetSourceRevertForSharedIDs(
	ctx context.Context, sharedIDs []string,
) ([]types.CalldataSignature, error) {
	r.spySourceRevertSharedIDs = sharedIDs
	return r.signatures, r.err
}

// stateUpdateCall captures one invocation of BatchSetState/BatchSetOutcome/BatchSetStateAndOutcome
// so callback tests can assert each transition independently. Either state or outcome may
// be the zero value when only one of them was set.
type stateUpdateCall struct {
	sharedIDs []string
	state     types.TransactionState
	outcome   types.TransactionOutcome
}

type StubSignatureTransactionRepository struct {
	calls []stateUpdateCall
	err   error
}

func (r *StubSignatureTransactionRepository) BatchSetState(
	ctx context.Context, sharedIDs []string, state types.TransactionState,
) error {
	r.calls = append(r.calls, stateUpdateCall{sharedIDs: sharedIDs, state: state, outcome: types.OutcomePending})
	return r.err
}

func (r *StubSignatureTransactionRepository) BatchSetOutcome(
	ctx context.Context, sharedIDs []string, outcome types.TransactionOutcome,
) error {
	r.calls = append(r.calls, stateUpdateCall{sharedIDs: sharedIDs, outcome: outcome})
	return r.err
}

func (r *StubSignatureTransactionRepository) BatchSetStateAndOutcome(
	ctx context.Context, sharedIDs []string, state types.TransactionState, outcome types.TransactionOutcome,
) error {
	r.calls = append(r.calls, stateUpdateCall{sharedIDs: sharedIDs, state: state, outcome: outcome})
	return r.err
}

// callForState returns the most recent call that transitioned to the given state, or
// nil if no such call happened.
func (r *StubSignatureTransactionRepository) callForState(state types.TransactionState) *stateUpdateCall {
	for i := range r.calls {
		if r.calls[i].state == state {
			return &r.calls[i]
		}
	}
	return nil
}

// callForOutcome returns the most recent call that set the given outcome (without state),
// or nil if no such call happened.
func (r *StubSignatureTransactionRepository) callForOutcome(outcome types.TransactionOutcome) *stateUpdateCall {
	for i := range r.calls {
		if r.calls[i].state == 0 && r.calls[i].outcome == outcome {
			return &r.calls[i]
		}
	}
	return nil
}

type StubSignatureSignatureBatchSender struct {
	spyMessages []types.TxRequest
	err         error
}

func (c *StubSignatureSignatureBatchSender) PushBatch(_ context.Context, msgs []types.TxRequest) error {
	c.spyMessages = msgs
	return c.err
}

type StubSignatureTeleportClient struct {
	spySharedIDs      []string
	spyAdditionalData []types.AtomicTeleportAdditionalData

	called bool
	err    error
}

func (c *StubSignatureTeleportClient) SendAdditionalDataBatch(
	_ context.Context, sharedIDs []string, data []types.AtomicTeleportAdditionalData,
) error {
	c.called = true
	c.spySharedIDs = sharedIDs
	c.spyAdditionalData = data
	return c.err
}

type StubSignatureEthereumClient struct {
	header *ethTypes.Header
	err    error

	spyBlockHash common.Hash
}

func (c *StubSignatureEthereumClient) HeaderByHash(
	ctx context.Context, blockHash common.Hash,
) (*ethTypes.Header, error) {
	c.spyBlockHash = blockHash
	return c.header, c.err
}

type StubSignatureTransactionSimulator struct {
	spyData      []byte
	revertReason txsim.ContractError
	err          error
}

func (s *StubSignatureTransactionSimulator) DecodeRevertBytes(data []byte) (txsim.ContractError, error) {
	s.spyData = data
	return s.revertReason, s.err
}

type StubSignatureGenerator struct {
	spySignatures []types.CalldataSignature
	calldata      []byte
	err           error
}

func (g *StubSignatureGenerator) Generate(sig types.CalldataSignature) ([]byte, error) {
	g.spySignatures = append(g.spySignatures, sig)
	return g.calldata, g.err
}

// ---------------------------------------------------------------------
// Fixture
// ---------------------------------------------------------------------

type fixture struct {
	sigRepo     *StubSignatureSignatureRepository
	txRepo      *StubSignatureTransactionRepository
	sigBatcher  *StubSignatureSignatureBatchSender
	teleportCli *StubSignatureTeleportClient
	ethCli      *StubSignatureEthereumClient
	simulator   *StubSignatureTransactionSimulator
	generator   *StubSignatureGenerator

	destEndpoint common.Address
	svc          *service2.SignatureService
}

func newFixture() *fixture {
	f := &fixture{
		sigRepo:     &StubSignatureSignatureRepository{},
		txRepo:      &StubSignatureTransactionRepository{},
		sigBatcher:  &StubSignatureSignatureBatchSender{},
		teleportCli: &StubSignatureTeleportClient{},
		// Default header avoids a nil-deref in getTimestampForBlockHash;
		// tests that care about the zero-timestamp branch override ethCli.err.
		ethCli:       &StubSignatureEthereumClient{header: &ethTypes.Header{}},
		simulator:    &StubSignatureTransactionSimulator{},
		generator:    &StubSignatureGenerator{calldata: []byte{0xCA, 0xFE}},
		destEndpoint: common.HexToAddress("0x1111111111111111111111111111111111111111"),
	}
	f.svc = service2.NewSignatureService(
		f.sigBatcher,
		f.teleportCli,
		f.ethCli,
		f.sigRepo,
		f.txRepo,
		f.simulator,
		f.generator,
		f.destEndpoint,
	)
	return f
}

// ---------------------------------------------------------------------
// Send path: HandleDestinationExecuted
// ---------------------------------------------------------------------

func TestSignatureService_HandleDestinationExecuted(t *testing.T) {
	t.Run("reads signatures, publishes TxRequests, transitions to DestinationUnlockSigs", func(t *testing.T) {
		sharedIDs := []string{"shared-1", "shared-2"}
		signatures := []types.CalldataSignature{
			{SharedId: "shared-1", SignatureType: types.UnlockOnDestinationSide},
			{SharedId: "shared-2", SignatureType: types.UnlockOnDestinationSide},
		}

		f := newFixture()
		f.sigRepo.signatures = signatures

		err := f.svc.HandleDestinationExecuted(context.Background(), sharedIDs)
		require.NoError(t, err)

		assert.Equal(t, sharedIDs, f.sigRepo.spyDestinationUnlockSharedIDs, "read wrong shared IDs from repo")

		require.Len(t, f.sigBatcher.spyMessages, 2)
		for i, msg := range f.sigBatcher.spyMessages {
			assert.Equal(t, signatures[i].SharedId, msg.CorrelationID)
			assert.Equal(t, f.destEndpoint, msg.Address)
			assert.Equal(t, []byte{0xCA, 0xFE}, msg.Calldata)
			assert.Equal(t, "atomic.destination-unlock", msg.MessageType)
		}

		call := f.txRepo.callForState(types.DestinationUnlockSigs)
		require.NotNil(t, call, "expected transition to DestinationUnlockSigs")
		assert.Equal(t, sharedIDs, call.sharedIDs)
	})

	t.Run("skips signatures the generator cannot handle", func(t *testing.T) {
		f := newFixture()
		f.sigRepo.signatures = []types.CalldataSignature{{SharedId: "x"}}
		f.generator.err = errors.New("boom")

		err := f.svc.HandleDestinationExecuted(context.Background(), []string{"x"})
		require.NoError(t, err)

		assert.Empty(t, f.sigBatcher.spyMessages, "should skip ungeneratable signatures")
	})
}

// ---------------------------------------------------------------------
// Send path: HandleDestinationReverted
// ---------------------------------------------------------------------

func TestSignatureService_HandleDestinationReverted(t *testing.T) {
	t.Run("reads signatures, publishes TxRequests, transitions to DestinationRevertSigs", func(t *testing.T) {
		sharedIDs := []string{"shared-1"}
		signatures := []types.CalldataSignature{
			{SharedId: "shared-1", SignatureType: types.RevertOnDestinationSide},
		}

		f := newFixture()
		f.sigRepo.signatures = signatures

		err := f.svc.HandleDestinationReverted(context.Background(), sharedIDs)
		require.NoError(t, err)

		assert.Equal(t, sharedIDs, f.sigRepo.spyDestinationRevertSharedIDs)
		require.Len(t, f.sigBatcher.spyMessages, 1)
		assert.Equal(t, "atomic.destination-revert", f.sigBatcher.spyMessages[0].MessageType)

		call := f.txRepo.callForState(types.DestinationRevertSigs)
		require.NotNil(t, call)
		assert.Equal(t, sharedIDs, call.sharedIDs)
	})
}

// ---------------------------------------------------------------------
// Send path: HandleSourceReverted
// ---------------------------------------------------------------------

func TestSignatureService_HandleSourceReverted(t *testing.T) {
	t.Run("reads signatures, publishes TxRequests, transitions to SourceRevertSigs", func(t *testing.T) {
		sharedIDs := []string{"shared-1"}
		signatures := []types.CalldataSignature{
			{SharedId: "shared-1", SignatureType: types.RevertOnSenderSide},
		}

		f := newFixture()
		f.sigRepo.signatures = signatures

		err := f.svc.HandleSourceReverted(context.Background(), sharedIDs)
		require.NoError(t, err)

		assert.Equal(t, sharedIDs, f.sigRepo.spySourceRevertSharedIDs)
		require.Len(t, f.sigBatcher.spyMessages, 1)
		assert.Equal(t, "atomic.source-revert", f.sigBatcher.spyMessages[0].MessageType)

		call := f.txRepo.callForState(types.SourceRevertSigs)
		require.NotNil(t, call, "expected transition to SourceRevertSigs")
		assert.Equal(t, sharedIDs, call.sharedIDs)
	})
}

// ---------------------------------------------------------------------
// Send path: HandleSourceExecuted (no migration — just state transition)
// ---------------------------------------------------------------------

func TestSignatureService_HandleSourceExecuted(t *testing.T) {
	t.Run("transitions to SourceFinalized", func(t *testing.T) {
		f := newFixture()

		err := f.svc.HandleSourceExecuted(context.Background(), []string{"shared-1"})
		require.NoError(t, err)

		call := f.txRepo.callForState(types.SourceFinalized)
		require.NotNil(t, call)
		assert.Equal(t, []string{"shared-1"}, call.sharedIDs)
		assert.Equal(t, types.OutcomeSuccess, call.outcome)
	})

	t.Run("wraps repo error in AtomicServiceError", func(t *testing.T) {
		underlying := errors.New("db down")

		f := newFixture()
		f.txRepo.err = underlying

		err := f.svc.HandleSourceExecuted(context.Background(), []string{"shared-1"})
		require.Error(t, err)

		var asErr *service2.AtomicServiceError
		require.ErrorAs(t, err, &asErr)
		require.ErrorIs(t, err, underlying)
	})
}

// ---------------------------------------------------------------------
// Callback: HandleDestinationExecutedCallback
// ---------------------------------------------------------------------

func TestSignatureService_HandleDestinationExecutedCallback(t *testing.T) {
	blockHash := common.HexToHash("0xC0FEBABE")
	txHash := common.HexToHash("0xDEADBEEF")

	t.Run("success results trigger additional data + outcome=success", func(t *testing.T) {
		f := newFixture()
		f.ethCli.header = &ethTypes.Header{Time: 123456}

		results := []types.TxResult{
			{
				CorrelationID: "shared-1",
				Kind:          types.TxResultSuccess,
				TxHash:        txHash,
				Receipt: &ethTypes.Receipt{
					Status:    1,
					TxHash:    txHash,
					BlockHash: blockHash,
				},
			},
		}
		err := f.svc.HandleDestinationExecutedCallback(context.Background(), results)
		require.NoError(t, err)

		// Teleport additional data dispatched using the destination-unlock fields
		require.True(t, f.teleportCli.called)
		assert.Equal(t, []string{"shared-1"}, f.teleportCli.spySharedIDs)
		require.Len(t, f.teleportCli.spyAdditionalData, 1)
		got := f.teleportCli.spyAdditionalData[0]
		assert.Equal(t, "shared-1", got.SharedId)
		assert.Equal(t, txHash, got.TxHashDestination)
		assert.Equal(t, int8(1), got.TxHashDestinationStatus)
		assert.Equal(t, uint64(123456), got.TxHashDestinationTimestamp)

		// Outcome transition
		call := f.txRepo.callForOutcome(types.OutcomeSuccess)
		require.NotNil(t, call)
		assert.Equal(t, []string{"shared-1"}, call.sharedIDs)
	})

	t.Run("revert results set outcome=reverted and decode revert bytes", func(t *testing.T) {
		f := newFixture()
		revertBytes := []byte{0x08, 0xc3, 0x79, 0xa0, 0xaa} // bogus but non-empty

		results := []types.TxResult{
			{CorrelationID: "shared-r", Kind: types.TxResultRevert, RevertData: revertBytes},
		}
		err := f.svc.HandleDestinationExecutedCallback(context.Background(), results)
		require.NoError(t, err)

		assert.Equal(t, revertBytes, f.simulator.spyData, "simulator should decode the revert bytes")
		assert.False(t, f.teleportCli.called, "no teleport call when no successes")

		call := f.txRepo.callForOutcome(types.OutcomeReverted)
		require.NotNil(t, call)
		assert.Equal(t, []string{"shared-r"}, call.sharedIDs)
	})

	t.Run("failed results set outcome=failed", func(t *testing.T) {
		f := newFixture()

		results := []types.TxResult{
			{CorrelationID: "shared-f", Kind: types.TxResultFailed, ErrorReason: "stuck"},
		}
		err := f.svc.HandleDestinationExecutedCallback(context.Background(), results)
		require.NoError(t, err)

		assert.False(t, f.teleportCli.called)

		call := f.txRepo.callForOutcome(types.OutcomeFailed)
		require.NotNil(t, call)
		assert.Equal(t, []string{"shared-f"}, call.sharedIDs)
	})

	t.Run("mixed batch fans out to three independent outcome transitions", func(t *testing.T) {
		f := newFixture()
		f.ethCli.header = &ethTypes.Header{Time: 1}

		results := []types.TxResult{
			{CorrelationID: "ok-1", Kind: types.TxResultSuccess, Receipt: &ethTypes.Receipt{Status: 1}},
			{CorrelationID: "rev-1", Kind: types.TxResultRevert, RevertData: []byte{0x01}},
			{CorrelationID: "fail-1", Kind: types.TxResultFailed, ErrorReason: "dead-lettered"},
			{CorrelationID: "ok-2", Kind: types.TxResultSuccess, Receipt: &ethTypes.Receipt{Status: 1}},
		}
		err := f.svc.HandleDestinationExecutedCallback(context.Background(), results)
		require.NoError(t, err)

		successCall := f.txRepo.callForOutcome(types.OutcomeSuccess)
		require.NotNil(t, successCall)
		assert.Equal(t, []string{"ok-1", "ok-2"}, successCall.sharedIDs, "successes should be in arrival order")

		revertCall := f.txRepo.callForOutcome(types.OutcomeReverted)
		require.NotNil(t, revertCall)
		assert.Equal(t, []string{"rev-1"}, revertCall.sharedIDs)

		failedCall := f.txRepo.callForOutcome(types.OutcomeFailed)
		require.NotNil(t, failedCall)
		assert.Equal(t, []string{"fail-1"}, failedCall.sharedIDs)
	})

	t.Run("empty results is a no-op", func(t *testing.T) {
		f := newFixture()

		err := f.svc.HandleDestinationExecutedCallback(context.Background(), nil)
		require.NoError(t, err)

		assert.False(t, f.teleportCli.called)
		assert.Empty(t, f.txRepo.calls)
	})

	t.Run("zero timestamp when header lookup fails", func(t *testing.T) {
		f := newFixture()
		f.ethCli.err = errors.New("no header")

		results := []types.TxResult{
			{
				CorrelationID: "shared-1",
				Kind:          types.TxResultSuccess,
				Receipt:       &ethTypes.Receipt{Status: 1, BlockHash: blockHash},
			},
		}
		err := f.svc.HandleDestinationExecutedCallback(context.Background(), results)
		require.NoError(t, err)

		require.Len(t, f.teleportCli.spyAdditionalData, 1)
		assert.Equal(t, uint64(0), f.teleportCli.spyAdditionalData[0].TxHashDestinationTimestamp)
		assert.Equal(t, blockHash, f.ethCli.spyBlockHash)
	})

	t.Run("wraps teleport error in AtomicServiceError", func(t *testing.T) {
		underlying := errors.New("teleport down")
		f := newFixture()
		f.teleportCli.err = underlying

		results := []types.TxResult{
			{CorrelationID: "shared-1", Kind: types.TxResultSuccess, Receipt: &ethTypes.Receipt{Status: 1}},
		}
		err := f.svc.HandleDestinationExecutedCallback(context.Background(), results)
		require.Error(t, err)

		var asErr *service2.AtomicServiceError
		require.ErrorAs(t, err, &asErr)
		require.ErrorIs(t, err, underlying)
	})
}

// ---------------------------------------------------------------------
// Callback: HandleDestinationRevertedCallback (happy path — field shape
// differs from Executed, everything else is identical)
// ---------------------------------------------------------------------

func TestSignatureService_HandleDestinationRevertedCallback(t *testing.T) {
	t.Run("success results populate TxHashDestinationRevert* and set outcome=success", func(t *testing.T) {
		blockHash := common.HexToHash("0xAB")
		txHash := common.HexToHash("0xCD")

		f := newFixture()
		f.ethCli.header = &ethTypes.Header{Time: 42}

		results := []types.TxResult{
			{
				CorrelationID: "s-1",
				Kind:          types.TxResultSuccess,
				Receipt: &ethTypes.Receipt{
					Status:    1,
					TxHash:    txHash,
					BlockHash: blockHash,
				},
			},
		}
		err := f.svc.HandleDestinationRevertedCallback(context.Background(), results)
		require.NoError(t, err)

		require.Len(t, f.teleportCli.spyAdditionalData, 1)
		got := f.teleportCli.spyAdditionalData[0]
		assert.Equal(t, "s-1", got.SharedId)
		assert.Equal(t, txHash, got.TxHashDestinationRevert)
		assert.Equal(t, int8(1), got.TxHashDestinationRevertStatus)
		assert.Equal(t, uint64(42), got.TxHashDestinationRevertTimestamp)

		call := f.txRepo.callForOutcome(types.OutcomeSuccess)
		require.NotNil(t, call)
		assert.Equal(t, []string{"s-1"}, call.sharedIDs)
	})
}

// ---------------------------------------------------------------------
// Callback: HandleSourceRevertedCallback (source-side — uses TxHashSourceRevert*
// fields and resolves the SourceRevertSigs phase via outcome=success)
// ---------------------------------------------------------------------

func TestSignatureService_HandleSourceRevertedCallback(t *testing.T) {
	t.Run("success results populate TxHashSourceRevert* and set outcome=success", func(t *testing.T) {
		blockHash := common.HexToHash("0xEF")
		txHash := common.HexToHash("0x01")

		f := newFixture()
		f.ethCli.header = &ethTypes.Header{Time: 7}

		results := []types.TxResult{
			{
				CorrelationID: "s-1",
				Kind:          types.TxResultSuccess,
				Receipt: &ethTypes.Receipt{
					Status:    1,
					TxHash:    txHash,
					BlockHash: blockHash,
				},
			},
		}
		err := f.svc.HandleSourceRevertedCallback(context.Background(), results)
		require.NoError(t, err)

		require.Len(t, f.teleportCli.spyAdditionalData, 1)
		got := f.teleportCli.spyAdditionalData[0]
		assert.Equal(t, "s-1", got.SharedId)
		assert.Equal(t, txHash, got.TxHashSourceRevert)
		assert.Equal(t, int8(1), got.TxHashSourceRevertStatus)
		assert.Equal(t, uint64(7), got.TxHashSourceRevertTimestamp)

		call := f.txRepo.callForOutcome(types.OutcomeSuccess)
		require.NotNil(t, call)
		assert.Equal(t, []string{"s-1"}, call.sharedIDs)
	})
}

// ---------------------------------------------------------------------
// Context propagation
// ---------------------------------------------------------------------

func TestSignatureService_ContextPropagation(t *testing.T) {
	t.Run("propagates cancelled context to signature repository", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		f := newFixture()
		f.sigRepo.err = context.Canceled

		err := f.svc.HandleDestinationExecuted(ctx, []string{"shared-id"})
		require.Error(t, err)
		require.ErrorIs(t, err, context.Canceled)
	})
}

// ---------------------------------------------------------------------
// Error wrapping — table-driven across the three async send paths
// ---------------------------------------------------------------------

func TestSignatureService_SendPathErrorHandling(t *testing.T) {
	sharedIDs := []string{"shared-id"}

	handlers := []struct {
		name string
		fn   func(*service2.SignatureService, context.Context, []string) error
	}{
		{"HandleDestinationExecuted", func(s *service2.SignatureService, ctx context.Context, ids []string) error {
			return s.HandleDestinationExecuted(ctx, ids)
		}},
		{"HandleDestinationReverted", func(s *service2.SignatureService, ctx context.Context, ids []string) error {
			return s.HandleDestinationReverted(ctx, ids)
		}},
		{"HandleSourceReverted", func(s *service2.SignatureService, ctx context.Context, ids []string) error {
			return s.HandleSourceReverted(ctx, ids)
		}},
	}

	t.Run("wraps signature repository errors in AtomicServiceError", func(t *testing.T) {
		for _, h := range handlers {
			t.Run(h.name, func(t *testing.T) {
				underlying := errors.New("sig repo down")
				f := newFixture()
				f.sigRepo.err = underlying

				err := h.fn(f.svc, context.Background(), sharedIDs)
				require.Error(t, err)

				var asErr *service2.AtomicServiceError
				require.ErrorAs(t, err, &asErr)
				require.ErrorIs(t, err, underlying)
			})
		}
	})

	t.Run("wraps batch send errors in AtomicServiceError", func(t *testing.T) {
		for _, h := range handlers {
			t.Run(h.name, func(t *testing.T) {
				underlying := errors.New("publish failed")
				f := newFixture()
				f.sigBatcher.err = underlying

				err := h.fn(f.svc, context.Background(), sharedIDs)
				require.Error(t, err)

				var asErr *service2.AtomicServiceError
				require.ErrorAs(t, err, &asErr)
				require.ErrorIs(t, err, underlying)
			})
		}
	})

	t.Run("wraps transaction repo errors in AtomicServiceError", func(t *testing.T) {
		for _, h := range handlers {
			t.Run(h.name, func(t *testing.T) {
				underlying := errors.New("tx repo down")
				f := newFixture()
				f.txRepo.err = underlying

				err := h.fn(f.svc, context.Background(), sharedIDs)
				require.Error(t, err)

				var asErr *service2.AtomicServiceError
				require.ErrorAs(t, err, &asErr)
				require.ErrorIs(t, err, underlying)
			})
		}
	})
}
