package service_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/batcher"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/EndpointV1"
	"github.com/raylsnetwork/rayls-sovereign-relayer/msgqueue"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools/fake"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools/spy"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newCrossChainSingleMessageMQ(msgSlice []msgqueue.Message[types.DispatchedMessageToPrivateHub]) *CrossChainMQMock {
	return &CrossChainMQMock{
		FetchFunc: fake.FetchMQ(msgSlice),
	}
}

func dummyVerifyFunc(mashaledProof []byte, rootTxHash common.Hash, txIdx uint) error {
	return nil
}

// ccFixture bundles the common mocks + the service under test so individual
// tests can poke at just the dependencies they care about.
type ccFixture struct {
	consumer       *CrossChainMQMock
	txRepo         *TransactionRepositoryMock
	sigRepo        *SignatureRepositoryMock
	txGen          *TransactionGeneratorMock
	atomicBatcher  *BatcherMock
	vanillaBatcher *BatcherMock
	atomicReceipt  *ReceiptHandlerMock
	vanillaReceipt *ReceiptHandlerMock
	endpointClient *EndpointClientMock
	deployerClient *DeployerClientMock

	plEndpointAddr common.Address
}

func newCCFixture() *ccFixture {
	return &ccFixture{
		consumer:       &CrossChainMQMock{},
		txRepo:         &TransactionRepositoryMock{},
		sigRepo:        &SignatureRepositoryMock{},
		txGen:          &TransactionGeneratorMock{},
		atomicBatcher:  &BatcherMock{},
		vanillaBatcher: &BatcherMock{},
		atomicReceipt:  &ReceiptHandlerMock{},
		vanillaReceipt: &ReceiptHandlerMock{},
		endpointClient: &EndpointClientMock{},
		deployerClient: &DeployerClientMock{},
		plEndpointAddr: common.HexToAddress("0x9999999999999999999999999999999999999999"),
	}
}

func (f *ccFixture) build(t *testing.T, verify service.VerifierFunc) *service.CrossChainService {
	t.Helper()
	return service.NewCrossChainServiceWith(
		f.plEndpointAddr,
		100,
		f.consumer,
		f.txRepo,
		f.sigRepo,
		f.txGen,
		f.atomicBatcher,
		f.vanillaBatcher,
		f.atomicReceipt,
		f.vanillaReceipt,
		f.endpointClient,
		f.deployerClient,
		verify,
	)
}

func noopBatchCreateTx() func(ctx context.Context, txs []types.Transaction) error {
	return func(ctx context.Context, txs []types.Transaction) error { return nil }
}

func noopBatchCreateSig() func(ctx context.Context, signatures []types.CalldataSignature) error {
	return func(ctx context.Context, signatures []types.CalldataSignature) error { return nil }
}

func noopBatcherSend() func(ctx context.Context, msgs []batcher.Message) error {
	return func(ctx context.Context, msgs []batcher.Message) error { return nil }
}

// ---------------------------------------------------------------------
// Run path tests
// ---------------------------------------------------------------------

func TestCrossChainService_Run(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("supports graceful shutdown", func(t *testing.T) {
		f := newCCFixture()
		f.consumer.FetchFunc = func(ctx context.Context, count int) ([]msgqueue.Message[types.DispatchedMessageToPrivateHub], error) {
			return nil, nil
		}
		f.txRepo.BatchCreateFunc = noopBatchCreateTx()
		f.sigRepo.BatchCreateFunc = noopBatchCreateSig()
		f.atomicBatcher.SendFunc = noopBatcherSend()
		f.vanillaBatcher.SendFunc = noopBatcherSend()

		svc := f.build(t, dummyVerifyFunc)
		assert.True(t, testtools.ShutdownFixture(t, svc.Run, time.Millisecond))
	})

	t.Run("respects context cancellation on consumer errors", func(t *testing.T) {
		f := newCCFixture()
		f.consumer.FetchFunc = func(ctx context.Context, count int) ([]msgqueue.Message[types.DispatchedMessageToPrivateHub], error) {
			return nil, errors.New("example error")
		}
		f.txRepo.BatchCreateFunc = noopBatchCreateTx()
		f.sigRepo.BatchCreateFunc = noopBatchCreateSig()
		f.atomicBatcher.SendFunc = noopBatcherSend()
		f.vanillaBatcher.SendFunc = noopBatcherSend()

		svc := f.build(t, dummyVerifyFunc)
		assert.True(t, testtools.ShutdownFixture(t, svc.Run, time.Millisecond))
	})

	t.Run("skips message on invalid proof", func(t *testing.T) {
		msg := types.DispatchedMessageToPrivateHub{
			IsAtomic:    true,
			TxLocation:  42,
			TxTrieProof: common.HexToHash("0xc0febabe"),
			Proofs:      common.Hex2Bytes("deadc0de"),
		}
		spyAck := spy.NewAck()
		mqMsg := msgqueue.Message[types.DispatchedMessageToPrivateHub]{V: msg, Ack: spyAck.Fn()}

		f := newCCFixture()
		f.consumer = newCrossChainSingleMessageMQ([]msgqueue.Message[types.DispatchedMessageToPrivateHub]{mqMsg})

		calledVerify := false
		verify := func(proof []byte, rootTxHash common.Hash, txIdx uint) error {
			assert.Equal(t, msg.Proofs, proof)
			assert.Equal(t, msg.TxTrieProof, rootTxHash)
			assert.Equal(t, msg.TxLocation, txIdx)
			calledVerify = true
			return errors.New("invalid proof")
		}

		f.txRepo.BatchCreateFunc = func(ctx context.Context, txs []types.Transaction) error {
			assert.Fail(t, "should not persist transactions on invalid proof")
			return nil
		}
		f.sigRepo.BatchCreateFunc = func(ctx context.Context, sigs []types.CalldataSignature) error {
			assert.Fail(t, "should not persist signatures on invalid proof")
			return nil
		}
		f.txGen.GenerateFunc = func(from *big.Int, _, _ common.Address, _ EndpointV1.RaylsMessage, _ common.Hash) ([]byte, error) {
			assert.Fail(t, "should not generate calldata on invalid proof")
			return nil, nil
		}
		f.atomicBatcher.SendFunc = func(ctx context.Context, msgs []batcher.Message) error {
			assert.Fail(t, "should not publish on invalid proof")
			return nil
		}
		f.vanillaBatcher.SendFunc = func(ctx context.Context, msgs []batcher.Message) error {
			assert.Fail(t, "should not publish on invalid proof")
			return nil
		}
		f.endpointClient.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			assert.Fail(t, "should not query endpoint on invalid proof")
			return common.Address{}, nil
		}
		f.deployerClient.DeployResourceAndExecuteFunc = func(context.Context, [32]byte, *types.DispatchedMessageToPrivateHub) (common.Hash, error) {
			assert.Fail(t, "should not deploy on invalid proof")
			return common.Hash{}, nil
		}

		svc := f.build(t, verify)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		require.NoError(t, svc.Run(ctx))

		assert.True(t, calledVerify)
		spyAck.AssertCalled(t)
	})

	t.Run("atomic message publishes via atomicBatcher and persists signatures + tx", func(t *testing.T) {
		wantLockData := common.Hex2Bytes("c0febabe")
		wantRevertData := common.Hex2Bytes("deadc0de")
		wantSharedID := "shared-atomic-123"
		wantCalldata := common.Hex2Bytes("deadbeef")

		msg := types.DispatchedMessageToPrivateHub{
			IsAtomic: true,
			SharedId: wantSharedID,
			From:     common.HexToAddress("0x1111"),
			To:       common.HexToAddress("0x2222"),
			Data: EndpointV1.RaylsMessage{
				MessageMetadata: EndpointV1.RaylsMessageMetadata{
					LockData:                  wantLockData,
					RevertPayloadDataReceiver: wantRevertData,
				},
			},
		}

		spyAck := spy.NewAck()
		mqMsg := msgqueue.Message[types.DispatchedMessageToPrivateHub]{V: msg, Ack: spyAck.Fn()}

		f := newCCFixture()
		f.consumer = newCrossChainSingleMessageMQ([]msgqueue.Message[types.DispatchedMessageToPrivateHub]{mqMsg})

		sigVerifier := newSignatureVerifier(t, wantLockData, wantRevertData)
		f.sigRepo.BatchCreateFunc = func(ctx context.Context, sigs []types.CalldataSignature) error {
			spyAck.AssertNotCalled(t, "should not ack before persisting signatures")
			return sigVerifier.verify(sigs)
		}
		f.txRepo.BatchCreateFunc = func(ctx context.Context, txs []types.Transaction) error {
			spyAck.AssertNotCalled(t, "should not ack before persisting transactions")
			require.Len(t, txs, 1)
			assert.Equal(t, wantSharedID, txs[0].SharedID)
			assert.Equal(t, types.DestinationDispatch, txs[0].State)
			assert.Equal(t, types.OutcomePending, txs[0].Outcome)
			return nil
		}
		f.txGen.GenerateFunc = func(_ *big.Int, _, _ common.Address, _ EndpointV1.RaylsMessage, _ common.Hash) ([]byte, error) {
			return wantCalldata, nil
		}
		f.atomicBatcher.SendFunc = func(ctx context.Context, msgs []batcher.Message) error {
			require.Len(t, msgs, 1)
			assert.Equal(t, wantSharedID, msgs[0].ID)
			assert.Equal(t, f.plEndpointAddr, msgs[0].Address)
			assert.Equal(t, wantCalldata, msgs[0].Calldata)
			return nil
		}
		f.vanillaBatcher.SendFunc = func(ctx context.Context, msgs []batcher.Message) error {
			assert.Fail(t, "atomic message should not go to vanilla batcher")
			return nil
		}
		f.endpointClient.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.HexToAddress("0x1234"), nil
		}
		f.deployerClient.DeployResourceAndExecuteFunc = func(context.Context, [32]byte, *types.DispatchedMessageToPrivateHub) (common.Hash, error) {
			assert.Fail(t, "should not deploy when resource already deployed")
			return common.Hash{}, nil
		}

		svc := f.build(t, dummyVerifyFunc)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		require.NoError(t, svc.Run(ctx))

		sigVerifier.assertAllSignaturesSeen()
		assert.Equal(t, 1, len(f.atomicBatcher.SendCalls()))
		assert.Equal(t, 0, len(f.vanillaBatcher.SendCalls()))
		assert.Equal(t, 1, len(f.txRepo.BatchCreateCalls()))
		spyAck.AssertCalled(t)
	})

	t.Run("vanilla message publishes via vanillaBatcher (no signatures)", func(t *testing.T) {
		wantSharedID := "shared-vanilla-456"
		wantCalldata := common.Hex2Bytes("cafebabe")

		msg := types.DispatchedMessageToPrivateHub{
			IsAtomic: false,
			SharedId: wantSharedID,
			From:     common.HexToAddress("0x3333"),
			To:       common.HexToAddress("0x4444"),
		}
		spyAck := spy.NewAck()
		mqMsg := msgqueue.Message[types.DispatchedMessageToPrivateHub]{V: msg, Ack: spyAck.Fn()}

		f := newCCFixture()
		f.consumer = newCrossChainSingleMessageMQ([]msgqueue.Message[types.DispatchedMessageToPrivateHub]{mqMsg})

		f.sigRepo.BatchCreateFunc = func(ctx context.Context, sigs []types.CalldataSignature) error {
			assert.Fail(t, "vanilla message should not produce signatures")
			return nil
		}
		f.txRepo.BatchCreateFunc = noopBatchCreateTx()
		f.txGen.GenerateFunc = func(_ *big.Int, _, _ common.Address, _ EndpointV1.RaylsMessage, _ common.Hash) ([]byte, error) {
			return wantCalldata, nil
		}
		f.atomicBatcher.SendFunc = func(ctx context.Context, msgs []batcher.Message) error {
			assert.Fail(t, "vanilla message should not go to atomic batcher")
			return nil
		}
		f.vanillaBatcher.SendFunc = func(ctx context.Context, msgs []batcher.Message) error {
			require.Len(t, msgs, 1)
			assert.Equal(t, wantSharedID, msgs[0].ID)
			return nil
		}
		f.endpointClient.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.HexToAddress("0x1234"), nil
		}

		svc := f.build(t, dummyVerifyFunc)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		require.NoError(t, svc.Run(ctx))

		assert.Equal(t, 0, len(f.atomicBatcher.SendCalls()))
		assert.Equal(t, 1, len(f.vanillaBatcher.SendCalls()))
		assert.Equal(t, 0, len(f.sigRepo.BatchCreateCalls()))
		spyAck.AssertCalled(t)
	})

	t.Run("deploy-success routes to atomicReceiptSvc.HandleSuccessfullyMined", func(t *testing.T) {
		wantSharedID := "deploy-ok-123"
		wantResourceID := common.HexToHash("0xaabbccdd")
		wantDeploymentTxHash := common.HexToHash("0xdeadbeef")

		msg := types.DispatchedMessageToPrivateHub{
			SharedId:   wantSharedID,
			ResourceId: wantResourceID,
			IsAtomic:   true,
			From:       common.HexToAddress("0x1111"),
			To:         common.HexToAddress("0x2222"),
			Data: EndpointV1.RaylsMessage{
				MessageMetadata: EndpointV1.RaylsMessageMetadata{
					ResourceId: func() [32]byte { var a [32]byte; copy(a[:], wantResourceID[:]); return a }(),
				},
			},
		}
		spyAck := spy.NewAck()
		mqMsg := msgqueue.Message[types.DispatchedMessageToPrivateHub]{V: msg, Ack: spyAck.Fn()}

		f := newCCFixture()
		f.consumer = newCrossChainSingleMessageMQ([]msgqueue.Message[types.DispatchedMessageToPrivateHub]{mqMsg})

		f.endpointClient.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.Address{}, nil // zero → not deployed
		}
		deployCalled := false
		f.deployerClient.DeployResourceAndExecuteFunc = func(ctx context.Context, resourceId [32]byte, m *types.DispatchedMessageToPrivateHub) (common.Hash, error) {
			deployCalled = true
			assert.Equal(t, wantResourceID[:], resourceId[:])
			return wantDeploymentTxHash, nil
		}
		f.txGen.GenerateFunc = func(_ *big.Int, _, _ common.Address, _ EndpointV1.RaylsMessage, _ common.Hash) ([]byte, error) {
			assert.Fail(t, "deploy path should not call generator")
			return nil, nil
		}
		f.atomicBatcher.SendFunc = func(ctx context.Context, msgs []batcher.Message) error {
			assert.Empty(t, msgs, "deploy path should not publish")
			return nil
		}
		f.vanillaBatcher.SendFunc = noopBatcherSend()
		f.sigRepo.BatchCreateFunc = noopBatchCreateSig()
		f.txRepo.BatchCreateFunc = func(ctx context.Context, txs []types.Transaction) error {
			require.Len(t, txs, 1)
			assert.Equal(t, wantDeploymentTxHash, txs[0].TxHashDestination)
			assert.Equal(t, types.DestinationDispatch, txs[0].State)
			assert.Equal(t, types.OutcomePending, txs[0].Outcome)
			return nil
		}
		f.txRepo.BatchUpdateDestinationHashForSharedIDsFunc = func(ctx context.Context, hashBySharedID map[string]common.Hash) error {
			assert.Equal(t, map[string]common.Hash{wantSharedID: wantDeploymentTxHash}, hashBySharedID)
			return nil
		}
		f.atomicReceipt.HandleSuccessfullyMinedFunc = func(ctx context.Context, sharedIDs []string) error {
			assert.Equal(t, []string{wantSharedID}, sharedIDs)
			return nil
		}
		f.atomicReceipt.HandleFailedMinedFunc = func(ctx context.Context, sharedIDs []string) error {
			assert.Fail(t, "successful deploy should not hit HandleFailedMined")
			return nil
		}

		svc := f.build(t, dummyVerifyFunc)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		require.NoError(t, svc.Run(ctx))

		assert.True(t, deployCalled)
		assert.Equal(t, 1, len(f.atomicReceipt.HandleSuccessfullyMinedCalls()))
		assert.Equal(t, 0, len(f.vanillaReceipt.HandleSuccessfullyMinedCalls()))
		spyAck.AssertCalled(t)
	})

	t.Run("deploy-failure routes to atomicReceiptSvc.HandleFailedMined and marks tx failed", func(t *testing.T) {
		wantSharedID := "deploy-fail-456"
		wantResourceID := common.HexToHash("0xfa110000")

		msg := types.DispatchedMessageToPrivateHub{
			SharedId:   wantSharedID,
			ResourceId: wantResourceID,
			IsAtomic:   true,
			From:       common.HexToAddress("0x1111"),
			To:         common.HexToAddress("0x2222"),
		}
		spyAck := spy.NewAck()
		mqMsg := msgqueue.Message[types.DispatchedMessageToPrivateHub]{V: msg, Ack: spyAck.Fn()}

		f := newCCFixture()
		f.consumer = newCrossChainSingleMessageMQ([]msgqueue.Message[types.DispatchedMessageToPrivateHub]{mqMsg})

		f.endpointClient.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.Address{}, nil
		}
		f.deployerClient.DeployResourceAndExecuteFunc = func(context.Context, [32]byte, *types.DispatchedMessageToPrivateHub) (common.Hash, error) {
			return common.Hash{}, errors.New("out of gas")
		}
		f.atomicBatcher.SendFunc = noopBatcherSend()
		f.vanillaBatcher.SendFunc = noopBatcherSend()
		f.sigRepo.BatchCreateFunc = noopBatchCreateSig()
		f.txRepo.BatchCreateFunc = func(ctx context.Context, txs []types.Transaction) error {
			require.Len(t, txs, 1)
			assert.Equal(t, types.DestinationDispatch, txs[0].State)
			assert.Equal(t, types.OutcomeFailed, txs[0].Outcome)
			assert.Equal(t, common.Hash{}, txs[0].TxHashDestination)
			return nil
		}
		f.atomicReceipt.HandleFailedMinedFunc = func(ctx context.Context, sharedIDs []string) error {
			assert.Equal(t, []string{wantSharedID}, sharedIDs)
			return nil
		}
		f.atomicReceipt.HandleSuccessfullyMinedFunc = func(ctx context.Context, sharedIDs []string) error {
			assert.Fail(t, "failed deploy should not hit HandleSuccessfullyMined")
			return nil
		}

		svc := f.build(t, dummyVerifyFunc)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		require.NoError(t, svc.Run(ctx))

		assert.Equal(t, 1, len(f.atomicReceipt.HandleFailedMinedCalls()))
		// Failed deploys have no hash to update on the tx row, so no hash-update call.
		assert.Equal(t, 0, len(f.txRepo.BatchUpdateDestinationHashForSharedIDsCalls()))
		spyAck.AssertCalled(t)
	})

	t.Run("empty resource ID skips deployment check", func(t *testing.T) {
		msg := types.DispatchedMessageToPrivateHub{
			SharedId:   "arbitrary-789",
			ResourceId: common.Hash{},
			From:       common.HexToAddress("0x1111"),
			To:         common.HexToAddress("0x2222"),
		}
		spyAck := spy.NewAck()
		mqMsg := msgqueue.Message[types.DispatchedMessageToPrivateHub]{V: msg, Ack: spyAck.Fn()}

		f := newCCFixture()
		f.consumer = newCrossChainSingleMessageMQ([]msgqueue.Message[types.DispatchedMessageToPrivateHub]{mqMsg})
		f.endpointClient.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			assert.Fail(t, "should not check endpoint for empty resource ID")
			return common.Address{}, nil
		}
		f.deployerClient.DeployResourceAndExecuteFunc = func(context.Context, [32]byte, *types.DispatchedMessageToPrivateHub) (common.Hash, error) {
			assert.Fail(t, "should not deploy for empty resource ID")
			return common.Hash{}, nil
		}
		f.txGen.GenerateFunc = func(_ *big.Int, _, _ common.Address, _ EndpointV1.RaylsMessage, _ common.Hash) ([]byte, error) {
			return common.Hex2Bytes("ab"), nil
		}
		f.vanillaBatcher.SendFunc = func(ctx context.Context, msgs []batcher.Message) error {
			require.Len(t, msgs, 1)
			return nil
		}
		f.atomicBatcher.SendFunc = noopBatcherSend()
		f.sigRepo.BatchCreateFunc = noopBatchCreateSig()
		f.txRepo.BatchCreateFunc = noopBatchCreateTx()

		svc := f.build(t, dummyVerifyFunc)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		require.NoError(t, svc.Run(ctx))

		assert.Equal(t, 1, len(f.vanillaBatcher.SendCalls()))
		spyAck.AssertCalled(t)
	})

	t.Run("calldata generation failure marks tx failed and skips batch", func(t *testing.T) {
		msg := types.DispatchedMessageToPrivateHub{
			SharedId:   "gen-fail-101",
			ResourceId: common.Hash{},
			IsAtomic:   false,
		}
		spyAck := spy.NewAck()
		mqMsg := msgqueue.Message[types.DispatchedMessageToPrivateHub]{V: msg, Ack: spyAck.Fn()}

		f := newCCFixture()
		f.consumer = newCrossChainSingleMessageMQ([]msgqueue.Message[types.DispatchedMessageToPrivateHub]{mqMsg})
		f.endpointClient.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.HexToAddress("0x1234"), nil
		}
		f.txGen.GenerateFunc = func(_ *big.Int, _, _ common.Address, _ EndpointV1.RaylsMessage, _ common.Hash) ([]byte, error) {
			return nil, errors.New("bad format")
		}
		f.atomicBatcher.SendFunc = func(ctx context.Context, msgs []batcher.Message) error {
			assert.Fail(t, "should not publish on generation failure")
			return nil
		}
		f.vanillaBatcher.SendFunc = func(ctx context.Context, msgs []batcher.Message) error {
			assert.Fail(t, "should not publish on generation failure")
			return nil
		}
		f.sigRepo.BatchCreateFunc = noopBatchCreateSig()
		f.txRepo.BatchCreateFunc = func(ctx context.Context, txs []types.Transaction) error {
			require.Len(t, txs, 1)
			assert.Equal(t, types.DestinationDispatch, txs[0].State)
			assert.Equal(t, types.OutcomeFailed, txs[0].Outcome)
			return nil
		}

		svc := f.build(t, dummyVerifyFunc)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		require.NoError(t, svc.Run(ctx))

		assert.Equal(t, 1, len(f.txRepo.BatchCreateCalls()))
		spyAck.AssertCalled(t)
	})

	t.Run("does not ack when tx persistence fails", func(t *testing.T) {
		msg := types.DispatchedMessageToPrivateHub{
			SharedId:   "persist-fail-202",
			ResourceId: common.Hash{},
		}
		spyAck := spy.NewAck()
		mqMsg := msgqueue.Message[types.DispatchedMessageToPrivateHub]{V: msg, Ack: spyAck.Fn()}

		f := newCCFixture()
		f.consumer = newCrossChainSingleMessageMQ([]msgqueue.Message[types.DispatchedMessageToPrivateHub]{mqMsg})
		f.endpointClient.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.HexToAddress("0x1234"), nil
		}
		f.txGen.GenerateFunc = func(_ *big.Int, _, _ common.Address, _ EndpointV1.RaylsMessage, _ common.Hash) ([]byte, error) {
			return common.Hex2Bytes("ff"), nil
		}
		f.atomicBatcher.SendFunc = noopBatcherSend()
		f.vanillaBatcher.SendFunc = noopBatcherSend()
		f.sigRepo.BatchCreateFunc = noopBatchCreateSig()
		f.txRepo.BatchCreateFunc = func(ctx context.Context, txs []types.Transaction) error {
			return errors.New("db down")
		}

		svc := f.build(t, dummyVerifyFunc)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		require.NoError(t, svc.Run(ctx))

		spyAck.AssertNotCalled(t, "should not ack on persistence failure")
	})

	t.Run("does not ack when atomic batch publish fails", func(t *testing.T) {
		msg := types.DispatchedMessageToPrivateHub{
			SharedId:   "publish-fail-303",
			IsAtomic:   true,
			ResourceId: common.Hash{},
		}
		spyAck := spy.NewAck()
		mqMsg := msgqueue.Message[types.DispatchedMessageToPrivateHub]{V: msg, Ack: spyAck.Fn()}

		f := newCCFixture()
		f.consumer = newCrossChainSingleMessageMQ([]msgqueue.Message[types.DispatchedMessageToPrivateHub]{mqMsg})
		f.endpointClient.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.HexToAddress("0x1234"), nil
		}
		f.txGen.GenerateFunc = func(_ *big.Int, _, _ common.Address, _ EndpointV1.RaylsMessage, _ common.Hash) ([]byte, error) {
			return common.Hex2Bytes("ff"), nil
		}
		f.atomicBatcher.SendFunc = func(ctx context.Context, msgs []batcher.Message) error {
			return errors.New("publish down")
		}
		f.vanillaBatcher.SendFunc = noopBatcherSend()
		f.sigRepo.BatchCreateFunc = noopBatchCreateSig()
		f.txRepo.BatchCreateFunc = func(ctx context.Context, txs []types.Transaction) error {
			assert.Fail(t, "should not persist when publish fails")
			return nil
		}

		svc := f.build(t, dummyVerifyFunc)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		require.NoError(t, svc.Run(ctx))

		spyAck.AssertNotCalled(t, "should not ack on publish failure")
	})
}

// ---------------------------------------------------------------------
// Callback path tests
// ---------------------------------------------------------------------

func TestCrossChainService_HandleAtomicResults(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("success results update destination hash and call HandleSuccessfullyMined", func(t *testing.T) {
		wantHash := common.HexToHash("0xaa")
		results := []types.TxResult{
			{CorrelationID: "s-1", Kind: types.TxResultSuccess, TxHash: wantHash},
		}

		f := newCCFixture()
		f.txRepo.BatchUpdateDestinationHashForSharedIDsFunc = func(ctx context.Context, hashMap map[string]common.Hash) error {
			assert.Equal(t, map[string]common.Hash{"s-1": wantHash}, hashMap)
			return nil
		}
		f.atomicReceipt.HandleSuccessfullyMinedFunc = func(ctx context.Context, sharedIDs []string) error {
			assert.Equal(t, []string{"s-1"}, sharedIDs)
			return nil
		}

		svc := f.build(t, dummyVerifyFunc)
		require.NoError(t, svc.HandleAtomicResults(context.Background(), results))

		assert.Equal(t, 1, len(f.txRepo.BatchUpdateDestinationHashForSharedIDsCalls()))
		assert.Equal(t, 1, len(f.atomicReceipt.HandleSuccessfullyMinedCalls()))
		assert.Equal(t, 0, len(f.atomicReceipt.HandleFailedMinedCalls()))
	})

	t.Run("revert and failed results both route to HandleFailedMined", func(t *testing.T) {
		results := []types.TxResult{
			{CorrelationID: "r-1", Kind: types.TxResultRevert, RevertData: []byte{0x01}},
			{CorrelationID: "f-1", Kind: types.TxResultFailed, ErrorReason: "stuck"},
		}

		f := newCCFixture()
		f.atomicReceipt.HandleFailedMinedFunc = func(ctx context.Context, sharedIDs []string) error {
			assert.ElementsMatch(t, []string{"r-1", "f-1"}, sharedIDs)
			return nil
		}

		svc := f.build(t, dummyVerifyFunc)
		require.NoError(t, svc.HandleAtomicResults(context.Background(), results))

		assert.Equal(t, 0, len(f.atomicReceipt.HandleSuccessfullyMinedCalls()))
		assert.Equal(t, 0, len(f.txRepo.BatchUpdateDestinationHashForSharedIDsCalls()))
		assert.Equal(t, 1, len(f.atomicReceipt.HandleFailedMinedCalls()))
	})

	t.Run("mixed batch splits correctly", func(t *testing.T) {
		results := []types.TxResult{
			{CorrelationID: "ok-1", Kind: types.TxResultSuccess, TxHash: common.HexToHash("0xab")},
			{CorrelationID: "rev-1", Kind: types.TxResultRevert},
			{CorrelationID: "ok-2", Kind: types.TxResultSuccess, TxHash: common.HexToHash("0xcd")},
			{CorrelationID: "fail-1", Kind: types.TxResultFailed},
		}

		f := newCCFixture()
		f.txRepo.BatchUpdateDestinationHashForSharedIDsFunc = func(ctx context.Context, hashMap map[string]common.Hash) error {
			assert.Len(t, hashMap, 2)
			return nil
		}
		f.atomicReceipt.HandleSuccessfullyMinedFunc = func(ctx context.Context, sharedIDs []string) error {
			assert.Equal(t, []string{"ok-1", "ok-2"}, sharedIDs)
			return nil
		}
		f.atomicReceipt.HandleFailedMinedFunc = func(ctx context.Context, sharedIDs []string) error {
			assert.Equal(t, []string{"rev-1", "fail-1"}, sharedIDs)
			return nil
		}

		svc := f.build(t, dummyVerifyFunc)
		require.NoError(t, svc.HandleAtomicResults(context.Background(), results))

		assert.Equal(t, 1, len(f.atomicReceipt.HandleSuccessfullyMinedCalls()))
		assert.Equal(t, 1, len(f.atomicReceipt.HandleFailedMinedCalls()))
	})

	t.Run("empty results is a no-op", func(t *testing.T) {
		f := newCCFixture()

		svc := f.build(t, dummyVerifyFunc)
		require.NoError(t, svc.HandleAtomicResults(context.Background(), nil))

		assert.Equal(t, 0, len(f.atomicReceipt.HandleSuccessfullyMinedCalls()))
		assert.Equal(t, 0, len(f.atomicReceipt.HandleFailedMinedCalls()))
	})

	t.Run("does not touch vanilla receipt service", func(t *testing.T) {
		results := []types.TxResult{
			{CorrelationID: "s-1", Kind: types.TxResultSuccess, TxHash: common.HexToHash("0xaa")},
		}

		f := newCCFixture()
		f.txRepo.BatchUpdateDestinationHashForSharedIDsFunc = func(context.Context, map[string]common.Hash) error { return nil }
		f.atomicReceipt.HandleSuccessfullyMinedFunc = func(context.Context, []string) error { return nil }

		svc := f.build(t, dummyVerifyFunc)
		require.NoError(t, svc.HandleAtomicResults(context.Background(), results))

		assert.Equal(t, 0, len(f.vanillaReceipt.HandleSuccessfullyMinedCalls()))
		assert.Equal(t, 0, len(f.vanillaReceipt.HandleFailedMinedCalls()))
	})
}

func TestCrossChainService_HandleVanillaResults(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("success results route to vanillaReceiptSvc", func(t *testing.T) {
		wantHash := common.HexToHash("0xbb")
		results := []types.TxResult{
			{CorrelationID: "v-1", Kind: types.TxResultSuccess, TxHash: wantHash},
		}

		f := newCCFixture()
		f.txRepo.BatchUpdateDestinationHashForSharedIDsFunc = func(ctx context.Context, hashMap map[string]common.Hash) error {
			assert.Equal(t, map[string]common.Hash{"v-1": wantHash}, hashMap)
			return nil
		}
		f.vanillaReceipt.HandleSuccessfullyMinedFunc = func(ctx context.Context, sharedIDs []string) error {
			assert.Equal(t, []string{"v-1"}, sharedIDs)
			return nil
		}

		svc := f.build(t, dummyVerifyFunc)
		require.NoError(t, svc.HandleVanillaResults(context.Background(), results))

		assert.Equal(t, 1, len(f.vanillaReceipt.HandleSuccessfullyMinedCalls()))
		assert.Equal(t, 0, len(f.atomicReceipt.HandleSuccessfullyMinedCalls()))
	})

	t.Run("failed results route to vanillaReceiptSvc.HandleFailedMined", func(t *testing.T) {
		results := []types.TxResult{
			{CorrelationID: "v-f", Kind: types.TxResultFailed, ErrorReason: "stuck"},
		}

		f := newCCFixture()
		f.vanillaReceipt.HandleFailedMinedFunc = func(ctx context.Context, sharedIDs []string) error {
			assert.Equal(t, []string{"v-f"}, sharedIDs)
			return nil
		}

		svc := f.build(t, dummyVerifyFunc)
		require.NoError(t, svc.HandleVanillaResults(context.Background(), results))

		assert.Equal(t, 1, len(f.vanillaReceipt.HandleFailedMinedCalls()))
	})
}
