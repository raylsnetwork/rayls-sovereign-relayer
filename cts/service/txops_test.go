package service

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Test doubles ---

type stubKeyQueue struct {
	key        *ecdsa.PrivateKey
	dequeueErr error
	enqueueCnt int
	dequeueCnt int
}

func (s *stubKeyQueue) Enqueue(key *ecdsa.PrivateKey) { s.enqueueCnt++ }
func (s *stubKeyQueue) Dequeue(ctx context.Context) (*ecdsa.PrivateKey, error) {
	s.dequeueCnt++
	if s.dequeueErr != nil {
		return nil, s.dequeueErr
	}
	return s.key, nil
}

type stubAuthGen struct {
	err              error
	invalidatedAddrs []common.Address
}

func (s *stubAuthGen) Get(_ context.Context, key *ecdsa.PrivateKey) (*bind.TransactOpts, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &bind.TransactOpts{
		From: crypto.PubkeyToAddress(key.PublicKey),
		Signer: func(addr common.Address, tx *ethtypes.Transaction) (*ethtypes.Transaction, error) {
			return tx, nil
		},
	}, nil
}

func (s *stubAuthGen) InvalidateNonceCache(addr common.Address) {
	s.invalidatedAddrs = append(s.invalidatedAddrs, addr)
}

// stubSyncTxRepo is a minimal txOpsTxSyncRepository double. The zero value
// behaves like an empty ledger: Get reports the id as unseen (ErrTxNotFound),
// so SignAndSend takes the fresh path; Create and Save are no-op successes.
type stubSyncTxRepo struct {
	getTx     *SyncTransaction
	getErr    error
	createErr error
	saveErr   error
}

func (s *stubSyncTxRepo) Create(_ context.Context, _ *SyncTransaction) error { return s.createErr }
func (s *stubSyncTxRepo) Save(_ context.Context, _ *SyncTransaction) error   { return s.saveErr }
func (s *stubSyncTxRepo) Get(_ context.Context, _ string) (*SyncTransaction, error) {
	if s.getTx == nil && s.getErr == nil {
		return nil, ErrTxNotFound
	}
	return s.getTx, s.getErr
}

// fakeDataError implements rpc.DataError with an ABI-encoded revert blob.
type fakeDataError struct{ data string }

func (f fakeDataError) Error() string  { return "execution reverted" }
func (f fakeDataError) ErrorData() any { return f.data }
func (f fakeDataError) ErrorCode() int { return 3 }

// --- Tests ---

func TestSignAndSend_DequeueError(t *testing.T) {
	t.Parallel()

	queue := &stubKeyQueue{dequeueErr: errors.New("queue drained")}
	svc := NewTxOpsService(&stubAuthGen{}, queue, nil, nil, &stubSyncTxRepo{})

	result, err := svc.SignAndSend(context.Background(), "test-id", []byte{0x01}, common.HexToAddress("0xdead"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "dequeue key")
	assert.Nil(t, result.Success)
	assert.Nil(t, result.Revert)
	assert.Equal(t, 0, queue.enqueueCnt, "key should not be re-enqueued if dequeue failed")
}

func TestSignAndSend_AuthGenError_StillRequeuesKey(t *testing.T) {
	t.Parallel()

	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	queue := &stubKeyQueue{key: key}
	gen := &stubAuthGen{err: errors.New("kms offline")}
	svc := NewTxOpsService(gen, queue, nil, nil, &stubSyncTxRepo{})

	result, err := svc.SignAndSend(context.Background(), "test-id", []byte{0x01}, common.HexToAddress("0xdead"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "generate transaction auth")
	assert.Nil(t, result.Success)
	assert.Nil(t, result.Revert)
	assert.Equal(t, 1, queue.enqueueCnt, "key must be returned to the queue on auth failure")
}

func TestExtractRevertData(t *testing.T) {
	t.Parallel()

	raw := "deadbeefcafebabe"
	rawBytes, err := hex.DecodeString(raw)
	require.NoError(t, err)

	tests := []struct {
		name     string
		err      error
		wantData []byte
		wantOk   bool
	}{
		{
			name:     "extracts 0x-prefixed revert blob",
			err:      fakeDataError{data: "0x" + raw},
			wantData: rawBytes,
			wantOk:   true,
		},
		{
			name:     "extracts un-prefixed revert blob",
			err:      fakeDataError{data: raw},
			wantData: rawBytes,
			wantOk:   true,
		},
		{
			name:     "non-DataError returns false",
			err:      errors.New("generic rpc error"),
			wantData: nil,
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotData, gotOk := extractRevertData(tt.err)
			assert.Equal(t, tt.wantOk, gotOk)
			assert.Equal(t, tt.wantData, gotData)
		})
	}
}

// callOnlyEthClient is a minimal txOpsEthClient stub for the Call path. Embeds
// ContractBackend as a nil interface so the other backend methods panic on
// use — tests fail loudly if Call's code path ever starts touching them.
// DeployBackend isn't embedded here because it'd collide with ContractBackend
// on CodeAt (ambiguous selector); the Call path never needs it anyway.
type callOnlyEthClient struct {
	bind.ContractBackend
	callResult []byte
	callErr    error
	callCnt    int
}

// TransactionReceipt satisfies the bind.DeployBackend half of txOpsEthClient
// without embedding DeployBackend (which would collide with ContractBackend
// on CodeAt). Never called from the Call path — panics if touched.
func (c *callOnlyEthClient) TransactionReceipt(_ context.Context, _ common.Hash) (*ethtypes.Receipt, error) {
	panic("callOnlyEthClient.TransactionReceipt should not be called from Call path")
}

// ChainID satisfies the widened txOpsEthClient interface. Never called by the
// Call path — panics if touched.
func (c *callOnlyEthClient) ChainID(_ context.Context) (*big.Int, error) {
	panic("callOnlyEthClient.ChainID should not be called from Call path")
}

// TransactionByHash satisfies the widened txOpsEthClient interface (recovery
// path). Never called from the Call path — panics if touched.
func (c *callOnlyEthClient) TransactionByHash(_ context.Context, _ common.Hash) (*ethtypes.Transaction, bool, error) {
	panic("callOnlyEthClient.TransactionByHash should not be called from Call path")
}

func (c *callOnlyEthClient) CallContract(_ context.Context, _ ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	c.callCnt++
	return c.callResult, c.callErr
}

func TestCall_DequeueError(t *testing.T) {
	t.Parallel()

	queue := &stubKeyQueue{dequeueErr: errors.New("queue drained")}
	svc := NewTxOpsService(&stubAuthGen{}, queue, nil, nil, &stubSyncTxRepo{})

	result, err := svc.Call(context.Background(), []byte{0x01}, common.HexToAddress("0xdead"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "dequeue key")
	assert.Nil(t, result.Value)
	assert.Nil(t, result.Revert)
	assert.Equal(t, 0, queue.enqueueCnt, "key should not be re-enqueued if dequeue failed")
}

func TestCall_Success(t *testing.T) {
	t.Parallel()

	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	queue := &stubKeyQueue{key: key}
	ethClient := &callOnlyEthClient{callResult: []byte{0xde, 0xad, 0xbe, 0xef}}
	svc := NewTxOpsService(&stubAuthGen{}, queue, ethClient, nil, &stubSyncTxRepo{})

	result, err := svc.Call(context.Background(), []byte{0x01}, common.HexToAddress("0xdead"))

	require.NoError(t, err)
	assert.Equal(t, []byte{0xde, 0xad, 0xbe, 0xef}, result.Value)
	assert.Nil(t, result.Revert)
	assert.Equal(t, 1, ethClient.callCnt)
	assert.Equal(t, 1, queue.enqueueCnt, "key must be returned to the queue on success")
}

func TestCall_RevertExtractedFromDataError(t *testing.T) {
	t.Parallel()

	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	revertHex := "08c379a00000000000000000000000000000000000000000000000000000000000000020"
	revertBytes, err := hex.DecodeString(revertHex)
	require.NoError(t, err)

	queue := &stubKeyQueue{key: key}
	ethClient := &callOnlyEthClient{callErr: fakeDataError{data: "0x" + revertHex}}
	svc := NewTxOpsService(&stubAuthGen{}, queue, ethClient, nil, &stubSyncTxRepo{})

	result, err := svc.Call(context.Background(), []byte{0x01}, common.HexToAddress("0xdead"))

	require.NoError(t, err, "revert must be surfaced as data, not as an error")
	assert.Nil(t, result.Value)
	require.NotNil(t, result.Revert)
	assert.Equal(t, revertBytes, result.Revert.RevertData)
	assert.Equal(t, 1, queue.enqueueCnt, "key must be returned to the queue on revert")
}

func TestCall_InfraErrorWrapped(t *testing.T) {
	t.Parallel()

	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	queue := &stubKeyQueue{key: key}
	ethClient := &callOnlyEthClient{callErr: errors.New("connection refused")}
	svc := NewTxOpsService(&stubAuthGen{}, queue, ethClient, nil, &stubSyncTxRepo{})

	result, err := svc.Call(context.Background(), []byte{0x01}, common.HexToAddress("0xdead"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "call contract")
	assert.Contains(t, err.Error(), "connection refused")
	assert.Nil(t, result.Value)
	assert.Nil(t, result.Revert)
	assert.Equal(t, 1, queue.enqueueCnt, "key must be returned to the queue on infra failure")
}

// idempotencyClient is a stub txOpsEthClient whose SendTransaction returns a
// configurable error. Used to assert the broadcastWithIdempotency classifier
// swallows "already known" / "nonce too low" errors without calling
// WaitMined. The rest of the txOpsEthClient surface panics if touched.
type idempotencyClient struct {
	bind.ContractBackend
	sendErr error
	sendCnt int
}

func (c *idempotencyClient) SendTransaction(_ context.Context, _ *ethtypes.Transaction) error {
	c.sendCnt++
	return c.sendErr
}

// TransactionReceipt satisfies bind.DeployBackend without embedding it
// (embedding would collide with ContractBackend on CodeAt). Never called.
func (c *idempotencyClient) TransactionReceipt(_ context.Context, _ common.Hash) (*ethtypes.Receipt, error) {
	panic("idempotencyClient.TransactionReceipt should not be called")
}

// ChainID satisfies the widened txOpsEthClient interface. Never called by the
// broadcastWithIdempotency path.
func (c *idempotencyClient) ChainID(_ context.Context) (*big.Int, error) {
	panic("idempotencyClient.ChainID should not be called")
}

// neverMinesClient always reports the tx as not-yet-mined, so bind.WaitMined
// keeps polling until its ctx is done — the exact "tx never mines" condition the
// WaitMined ceiling guards against.
type neverMinesClient struct {
	bind.ContractBackend
	receiptCnt int
}

func (c *neverMinesClient) TransactionReceipt(_ context.Context, _ common.Hash) (*ethtypes.Receipt, error) {
	c.receiptCnt++
	return nil, ethereum.NotFound
}

func (c *neverMinesClient) ChainID(_ context.Context) (*big.Int, error) {
	panic("neverMinesClient.ChainID should not be called")
}

// TransactionByHash satisfies the widened txOpsEthClient interface (recovery
// path). Never called by the WaitMined-ceiling tests — panics if touched.
func (c *neverMinesClient) TransactionByHash(_ context.Context, _ common.Hash) (*ethtypes.Transaction, bool, error) {
	panic("neverMinesClient.TransactionByHash should not be called")
}

func (c *neverMinesClient) CallContract(_ context.Context, _ ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	panic("neverMinesClient.CallContract should not be called")
}

// TestWaitMined_CeilingFires asserts the service-level WaitMined ceiling caps a
// tx that never mines even when the caller's ctx has no deadline of its own.
func TestWaitMined_CeilingFires(t *testing.T) {
	t.Parallel()

	client := &neverMinesClient{}
	svc := NewTxOpsServiceWithWaitMinedTimeout(&stubAuthGen{}, &stubKeyQueue{}, client, nil, &stubSyncTxRepo{}, 20*time.Millisecond)

	_, err := svc.waitMined(context.Background(), common.HexToHash("0xabc"))

	require.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Positive(t, client.receiptCnt, "WaitMined should have polled at least once before the ceiling fired")
}

// TestWaitMined_TighterCallerDeadlineWins asserts a caller deadline shorter than
// the ceiling still takes precedence — context.WithTimeout keeps the earliest
// deadline, so this is the relayer-handler-bound-wins case.
func TestWaitMined_TighterCallerDeadlineWins(t *testing.T) {
	t.Parallel()

	client := &neverMinesClient{}
	// Ceiling is generous; caller deadline is tight and must be the one that trips.
	svc := NewTxOpsServiceWithWaitMinedTimeout(&stubAuthGen{}, &stubKeyQueue{}, client, nil, &stubSyncTxRepo{}, time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	_, err := svc.waitMined(ctx, common.HexToHash("0xabc"))

	require.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

// TestNewTxOpsService_DefaultWaitMinedTimeout asserts the plain constructor and
// a non-positive override both fall back to the default ceiling.
func TestNewTxOpsService_DefaultWaitMinedTimeout(t *testing.T) {
	t.Parallel()

	assert.Equal(t, defaultWaitMinedTimeout, NewTxOpsService(&stubAuthGen{}, &stubKeyQueue{}, nil, nil, &stubSyncTxRepo{}).waitMinedTimeout)
	assert.Equal(t, defaultWaitMinedTimeout,
		NewTxOpsServiceWithWaitMinedTimeout(&stubAuthGen{}, &stubKeyQueue{}, nil, nil, &stubSyncTxRepo{}, 0).waitMinedTimeout)
}
