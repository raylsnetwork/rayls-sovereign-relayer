package contractclient_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/dvp"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/txsim"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StubDvpIntegrationEncryptor struct {
	encryptedBatches       [][]byte
	encryptedBalanceUpdate []byte
	err                    error

	spyBatches       []*types.EnygmaTransferBatch
	spyBlockNumber   *big.Int
	spyBalanceUpdate types.DvpBalanceUpdated
}

func (s *StubDvpIntegrationEncryptor) EncryptEnygmaTransferBatches(
	_ context.Context,
	batches []*types.EnygmaTransferBatch,
	blockNumber *big.Int,
) ([][]byte, error) {
	s.spyBatches = batches
	s.spyBlockNumber = blockNumber
	return s.encryptedBatches, s.err
}

func (s *StubDvpIntegrationEncryptor) EncryptDvpBalanceUpdated(_ context.Context, message types.DvpBalanceUpdated) ([]byte, error) {
	s.spyBalanceUpdate = message
	return s.encryptedBalanceUpdate, s.err
}

type StubZkdvpIntegrationEthClient struct {
	block *ethTypes.Block
	err   error
}

func (s *StubZkdvpIntegrationEthClient) BlockByNumber(ctx context.Context, number *big.Int) (*ethTypes.Block, error) {
	if s.block == nil && s.err == nil {
		header := &ethTypes.Header{
			Time: uint64(time.Now().Unix()), //nolint:gosec // test data with known values
		}
		return ethTypes.NewBlockWithHeader(header), nil
	}
	return s.block, s.err
}

type StubDvpIntegrationSimulator struct {
	contractError txsim.ContractError
	err           error

	spyHash common.Hash
}

func (s *StubDvpIntegrationSimulator) GetRevertReason(ctx context.Context, hash common.Hash) (txsim.ContractError, error) {
	s.spyHash = hash
	return s.contractError, s.err
}

func (s *StubDvpIntegrationSimulator) DecodeRevertReason(ctx context.Context, err error) (txsim.ContractError, error) {
	return s.contractError, s.err
}

func createTestProof() *types.EnygmaProofResponse {
	return &types.EnygmaProofResponse{
		PiA: [2]*big.Int{big.NewInt(1), big.NewInt(2)},
		PiB: [2][2]*big.Int{
			{big.NewInt(3), big.NewInt(4)},
			{big.NewInt(5), big.NewInt(6)},
		},
		PiC:          [2]*big.Int{big.NewInt(7), big.NewInt(8)},
		PublicSignal: []*big.Int{big.NewInt(9)},
	}
}

func createTestCommitments() []*types.Point {
	return []*types.Point{
		{X: big.NewInt(10), Y: big.NewInt(11)},
		{X: big.NewInt(12), Y: big.NewInt(13)},
	}
}

func createTestBatches() []*types.EnygmaTransferBatch {
	return []*types.EnygmaTransferBatch{
		{
			ToChainID: big.NewInt(100),
		},
		{
			ToChainID: big.NewInt(200),
		},
	}
}

func createTestProofReceipt() *dvp.ProofReceipt {
	return &dvp.ProofReceipt{
		Proof: &dvp.Proof{
			A: [2]*big.Int{big.NewInt(1), big.NewInt(2)},
			B: [2][2]*big.Int{
				{big.NewInt(3), big.NewInt(4)},
				{big.NewInt(5), big.NewInt(6)},
			},
			C: [2]*big.Int{big.NewInt(7), big.NewInt(8)},
		},
		TreeNumbers:      []*big.Int{big.NewInt(1)},
		Message:          big.NewInt(2),
		MerkleRoots:      []*big.Int{big.NewInt(3)},
		Commitments:      []*big.Int{big.NewInt(4)},
		Nullifiers:       []*big.Int{big.NewInt(5)},
		RevertCommitment: big.NewInt(0),
	}
}

func TestDvpIntegrationClient_SignDeposit(t *testing.T) {
	t.Run("successfully signs deposit via executor", func(t *testing.T) {
		batches := createTestBatches()
		proof := createTestProof()
		blockNumber := big.NewInt(100)

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		encryptor := &StubDvpIntegrationEncryptor{
			encryptedBatches: [][]byte{{0xaa}, {0xbb}},
		}

		client := contractclient.NewDvpIntegrationClient(
			executor,
			encryptor,
			&StubZkdvpIntegrationEthClient{},
		)

		err := client.Deposit(context.Background(), "test-event-id", batches, proof, blockNumber, big.NewInt(1337), "test-resource-id", big.NewInt(1000), common.Address{}, common.Hash{}, address)

		require.Nil(t, err)

		assert.Equal(t, batches, encryptor.spyBatches)
		assert.Equal(t, blockNumber, encryptor.spyBlockNumber)
		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotNil(t, executor.spyExecuteCalldata)
	})

	t.Run("wraps encryption errors in EnygmaDvpIntegrationClientError", func(t *testing.T) {
		wantError := errors.New("encryption failed")

		executor := &stubExecutor{}
		encryptor := &StubDvpIntegrationEncryptor{
			err: wantError,
		}

		client := contractclient.NewDvpIntegrationClient(
			executor,
			encryptor,
			&StubZkdvpIntegrationEthClient{},
		)

		err := client.Deposit(context.Background(), "test-event-id", createTestBatches(), createTestProof(), big.NewInt(100), big.NewInt(1337), "test-resource-id", big.NewInt(1000), common.Address{}, common.Hash{}, common.Address{})

		var wrappedErr *contractclient.EnygmaDvpIntegrationClientError
		require.ErrorAs(t, err, &wrappedErr, "error should be wrapped")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})

	t.Run("wraps executor sign errors in EnygmaDvpIntegrationClientError", func(t *testing.T) {
		wantError := errors.New("sign failed")

		executor := &stubExecutor{
			executeErr: wantError,
		}
		encryptor := &StubDvpIntegrationEncryptor{
			encryptedBatches: [][]byte{{0xaa}},
		}

		client := contractclient.NewDvpIntegrationClient(
			executor,
			encryptor,
			&StubZkdvpIntegrationEthClient{},
		)

		err := client.Deposit(context.Background(), "test-event-id", createTestBatches(), createTestProof(), big.NewInt(100), big.NewInt(1337), "test-resource-id", big.NewInt(1000), common.Address{}, common.Hash{}, common.Address{})

		var wrappedErr *contractclient.EnygmaDvpIntegrationClientError
		require.ErrorAs(t, err, &wrappedErr, "error should be wrapped")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})
}

func TestDvpIntegrationClient_SignWithdraw(t *testing.T) {
	t.Run("successfully signs withdraw via executor", func(t *testing.T) {
		batches := createTestBatches()
		proof := createTestProof()
		blockNumber := big.NewInt(100)
		jsProof := createTestProofReceipt()

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		encryptor := &StubDvpIntegrationEncryptor{
			encryptedBatches: [][]byte{{0xaa}, {0xbb}},
		}

		client := contractclient.NewDvpIntegrationClient(
			executor,
			encryptor,
			&StubZkdvpIntegrationEthClient{},
		)

		err := client.Withdraw(context.Background(), "test-event-id", batches, proof, blockNumber, jsProof, big.NewInt(1337), "test-resource-id", big.NewInt(1000), common.Address{}, common.Hash{}, address)

		require.Nil(t, err)

		assert.Equal(t, batches, encryptor.spyBatches)
		assert.Equal(t, blockNumber, encryptor.spyBlockNumber)
		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotNil(t, executor.spyExecuteCalldata)
	})

	t.Run("wraps encryption errors in EnygmaDvpIntegrationClientError", func(t *testing.T) {
		wantError := errors.New("encryption failed")

		executor := &stubExecutor{}
		encryptor := &StubDvpIntegrationEncryptor{
			err: wantError,
		}

		client := contractclient.NewDvpIntegrationClient(
			executor,
			encryptor,
			&StubZkdvpIntegrationEthClient{},
		)

		err := client.Withdraw(context.Background(), "test-event-id", createTestBatches(), createTestProof(), big.NewInt(100), createTestProofReceipt(), big.NewInt(1337), "test-resource-id", big.NewInt(1000), common.Address{}, common.Hash{}, common.Address{})

		var wrappedErr *contractclient.EnygmaDvpIntegrationClientError
		require.ErrorAs(t, err, &wrappedErr, "error should be wrapped")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})

	t.Run("wraps executor sign errors in EnygmaDvpIntegrationClientError", func(t *testing.T) {
		wantError := errors.New("sign failed")

		executor := &stubExecutor{
			executeErr: wantError,
		}
		encryptor := &StubDvpIntegrationEncryptor{
			encryptedBatches: [][]byte{{0xaa}},
		}

		client := contractclient.NewDvpIntegrationClient(
			executor,
			encryptor,
			&StubZkdvpIntegrationEthClient{},
		)

		err := client.Withdraw(context.Background(), "test-event-id", createTestBatches(), createTestProof(), big.NewInt(100), createTestProofReceipt(), big.NewInt(1337), "test-resource-id", big.NewInt(1000), common.Address{}, common.Hash{}, common.Address{})

		var wrappedErr *contractclient.EnygmaDvpIntegrationClientError
		require.ErrorAs(t, err, &wrappedErr, "error should be wrapped")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})
}

func TestDvpIntegrationClient_ConsolidateFunds(t *testing.T) {
	t.Run("successfully consolidates enygma deposit proof", func(t *testing.T) {
		proofReceipt := createTestProofReceipt()

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}

		client := contractclient.NewDvpIntegrationClient(
			executor,
			&StubDvpIntegrationEncryptor{},
			&StubZkdvpIntegrationEthClient{},
		)

		err := client.ConsolidateFunds(context.Background(), "test-event-id", address, proofReceipt)

		require.Nil(t, err)
		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotNil(t, executor.spyExecuteCalldata)
	})

	t.Run("wraps executor errors in EnygmaDvpIntegrationClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{
			executeErr: wantError,
		}

		client := contractclient.NewDvpIntegrationClient(
			executor,
			&StubDvpIntegrationEncryptor{},
			&StubZkdvpIntegrationEthClient{},
		)

		err := client.ConsolidateFunds(context.Background(), "test-event-id", address, createTestProofReceipt())

		var wrappedErr *contractclient.EnygmaDvpIntegrationClientError
		require.ErrorAs(t, err, &wrappedErr, "error should be wrapped")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})
}
