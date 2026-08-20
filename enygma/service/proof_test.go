package service_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/gen/proto/keys"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cryptography"
	"github.com/raylsnetwork/rayls-sovereign-relayer/enygma"
	"github.com/raylsnetwork/rayls-sovereign-relayer/enygma/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/enygma/testutils"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"google.golang.org/grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test helpers - custom EnygmaClient implementation for testing
type testEnygmaClient struct {
	GetPublicValuesFinalizedFunc func(ctx context.Context, tokenAddress common.Address) (*types.EnygmaPublicValues, error)
	callsMutex                   sync.RWMutex
	calls                        int
}

func (c *testEnygmaClient) GetPublicValuesFinalised(ctx context.Context, tokenAddress common.Address) (*types.EnygmaPublicValues, error) {
	c.callsMutex.Lock()
	c.calls++
	c.callsMutex.Unlock()
	if c.GetPublicValuesFinalizedFunc != nil {
		return c.GetPublicValuesFinalizedFunc(ctx, tokenAddress)
	}
	return nil, errors.New("not implemented")
}

func (c *testEnygmaClient) GetCallCount() int {
	c.callsMutex.RLock()
	defer c.callsMutex.RUnlock()
	return c.calls
}

func bigIntsToBytes(ints []*big.Int) [][]byte {
	out := make([][]byte, len(ints))
	for i, v := range ints {
		out[i] = v.Bytes()
	}
	return out
}

// Helper functions to create test data

func createTestPaymentSpendKey() *keys.PaymentSpendKeyResponse {
	secretKey := big.NewInt(12345) // Fixed secret key for testing
	// Compute the public key as Poseidon(sk, sk) mod JubJubPrimeSubGroup
	publicKey, err := cryptography.GetPoseidonHashModNumber(
		[]*big.Int{secretKey, secretKey},
		cryptography.JubJubPrimeSubGroup,
	)
	if err != nil {
		panic(fmt.Sprintf("failed to compute payment spend public key in test helper: %v", err))
	}
	return &keys.PaymentSpendKeyResponse{
		SecretKey: secretKey.Bytes(),
		PublicKey: publicKey.Bytes(),
	}
}

func createTestBatch(fromChainID, toChainID *big.Int, amount *big.Int) *types.EnygmaTransferBatch {
	txs := []*types.EnygmaTransferBatchTx{}

	// Create tx only for other chains
	if fromChainID.Cmp(toChainID) != 0 {
		txs = append(txs, &types.EnygmaTransferBatchTx{
			MessageId:     "msg-1",
			ReferenceId:   [32]byte{1, 2, 3, 4},
			FromAddress:   common.HexToAddress("0x1111111111111111111111111111111111111111"),
			ToAmount:      amount,
			ToAddress:     common.HexToAddress("0x2222222222222222222222222222222222222222"),
			SendTimestamp: 1234567890,
		})
	}

	return &types.EnygmaTransferBatch{
		ResourceId:            "test-resource-123",
		BlockNumberPrivateHub: big.NewInt(100),
		FromChainID:           fromChainID,
		ToChainID:             toChainID,
		ToRValueToAdd:         big.NewInt(0),
		Transactions:          txs,
		BatchId:               "batch-1",
		Ctx:                   context.Background(),
	}
}

func createTestEnygmaState() types.Enygma {
	return types.Enygma{
		ResourceId:           "test-resource-123",
		FinalizedR:           big.NewInt(5000),
		FinalizedBalance:     big.NewInt(100000),
		FinalizedBlockNumber: big.NewInt(50),
	}
}

func createTestEnygmaPublicValues(chainIDs []*big.Int) *types.EnygmaPublicValues {
	commitments := make(map[string]*types.Point)
	publicKeys := make(map[string]*big.Int)

	// Get the test enygma state to calculate correct commitment for sender chain (chainID 1)
	enygmaState := createTestEnygmaState()

	for _, chainID := range chainIDs {
		chainIDStr := chainID.String()

		var commitmentPoint *babyjub.Point
		if chainID.Int64() == 1 {
			// For the sender chain (chainID 1), use the Pedersen commitment from the enygma state
			commitmentPoint = cryptography.PedersenCommitmentEnygma(
				enygmaState.FinalizedBalance,
				enygmaState.FinalizedR,
			)
		} else {
			// For other chains, generate arbitrary point for commitment
			commitmentScalar := big.NewInt(int64(chainID.Int64())*12345 + 1)
			commitmentPoint = babyjub.NewPoint().Mul(commitmentScalar, cryptography.G)
		}

		// Generate public key as big.Int (now computed as Poseidon(sk, sk))
		publicKeyScalar := big.NewInt(int64(chainID.Int64())*54321 + 2)

		commitments[chainIDStr] = &types.Point{
			X: commitmentPoint.X,
			Y: commitmentPoint.Y,
		}
		publicKeys[chainIDStr] = publicKeyScalar
	}

	return &types.EnygmaPublicValues{
		Commitments: commitments,
		PublicKeys:  publicKeys,
	}
}

func TestGenerateTransferProof(t *testing.T) {
	t.Run("returns error if anonymity index is invalid - unique chainIDs exceed anonymity index", func(t *testing.T) {
		ctx := context.Background()
		chainID1 := big.NewInt(1)
		chainID2 := big.NewInt(2)
		chainID3 := big.NewInt(3)
		chainID4 := big.NewInt(4)
		senderAmount := big.NewInt(1000)

		// Create batches with 3 unique destination chains
		batch1 := createTestBatch(chainID1, chainID2, senderAmount)
		batch2 := createTestBatch(chainID1, chainID3, senderAmount)
		batch3 := createTestBatch(chainID1, chainID4, senderAmount)

		params := enygma.TransferProofParams{
			ResourceId:     "test-resource-123",
			AnonymityIndex: 2, // Only supports 2 chains, but we have 3 unique destinations
			SenderAmount:   senderAmount,
			BlockNumber:    big.NewInt(100),
			Batches:        []*types.EnygmaTransferBatch{batch1, batch2, batch3},
			TokenAddress:   common.HexToAddress("0xDEADBEEF"),
		}

		// Create service with mocks (they won't be called due to early validation)
		kosClient := &MockProofKOSClient{}
		proofClient := &MockProofAPIClient{}
		enygmaRepo := &MockProofEnygmaRepository{}
		enygmaClient := &testEnygmaClient{}
		tracer := &testutils.MockTracer{}

		svc := service.NewEnygmaProofService(
			chainID1, kosClient, proofClient, enygmaRepo, enygmaClient, tracer,
		)

		proof, destCommitments, destRFactors, _, err := svc.GenerateTransferProof(ctx, params)

		require.Error(t, err)
		assert.Nil(t, proof)
		assert.Nil(t, destCommitments)
		assert.Nil(t, destRFactors)
		assert.Contains(t, err.Error(), "Only transactions from 1->1 up to 1->k-1 are supported")
	})

	t.Run("returns error if kosClient.GenerateEnygmaSharedSecrets fails", func(t *testing.T) {
		ctx := context.Background()
		chainID1 := big.NewInt(1)
		chainID2 := big.NewInt(2)
		senderAmount := big.NewInt(1000)

		batch := createTestBatch(chainID1, chainID2, senderAmount)

		params := enygma.TransferProofParams{
			ResourceId:     "test-resource-123",
			AnonymityIndex: 2,
			SenderAmount:   senderAmount,
			BlockNumber:    big.NewInt(100),
			Batches:        []*types.EnygmaTransferBatch{batch},
			TokenAddress:   common.HexToAddress("0xDEADBEEF"),
		}

		kosClient := &MockProofKOSClient{
			GetPaymentSpendKeyFunc: func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
				return createTestPaymentSpendKey(), nil
			},
			GenerateEnygmaSharedSecretsFunc: func(ctx context.Context, in *keys.GenerateEnygmaSharedSecretsRequest, opts ...grpc.CallOption) (*keys.GenerateEnygmaSharedSecretsResponse, error) {
				return nil, errors.New("shared secrets service down")
			},
		}
		proofClient := &MockProofAPIClient{}
		enygmaRepo := &MockProofEnygmaRepository{
			GetEnygmaByResourceIdFunc: func(ctx context.Context, resourceId string) (types.Enygma, error) {
				return createTestEnygmaState(), nil
			},
		}
		enygmaClient := &testEnygmaClient{
			GetPublicValuesFinalizedFunc: func(ctx context.Context, _ common.Address) (*types.EnygmaPublicValues, error) {
				return createTestEnygmaPublicValues([]*big.Int{chainID2}), nil
			},
		}
		tracer := &testutils.MockTracer{}

		svc := service.NewEnygmaProofService(
			chainID1, kosClient, proofClient, enygmaRepo, enygmaClient, tracer,
		)

		proof, destCommitments, destRFactors, _, err := svc.GenerateTransferProof(ctx, params)

		require.Error(t, err)
		assert.Nil(t, proof)
		assert.Nil(t, destCommitments)
		assert.Nil(t, destRFactors)
		assert.Contains(t, err.Error(), "shared secrets service down")
	})

	t.Run("returns error if kosClient.GetPaymentSpendKey fails", func(t *testing.T) {
		ctx := context.Background()
		chainID1 := big.NewInt(1)
		chainID2 := big.NewInt(2)
		senderAmount := big.NewInt(1000)

		batch := createTestBatch(chainID1, chainID2, senderAmount)

		params := enygma.TransferProofParams{
			ResourceId:     "test-resource-123",
			AnonymityIndex: 2,
			SenderAmount:   senderAmount,
			BlockNumber:    big.NewInt(100),
			Batches:        []*types.EnygmaTransferBatch{batch},
			TokenAddress:   common.HexToAddress("0xDEADBEEF"),
		}

		kosClient := &MockProofKOSClient{
			GetPaymentSpendKeyFunc: func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
				return nil, errors.New("enygma key error")
			},
		}
		proofClient := &MockProofAPIClient{}
		enygmaRepo := &MockProofEnygmaRepository{
			GetEnygmaByResourceIdFunc: func(ctx context.Context, resourceId string) (types.Enygma, error) {
				return createTestEnygmaState(), nil
			},
		}
		enygmaClient := &testEnygmaClient{}
		tracer := &testutils.MockTracer{}

		svc := service.NewEnygmaProofService(
			chainID1, kosClient, proofClient, enygmaRepo, enygmaClient, tracer,
		)

		proof, destCommitments, destRFactors, _, err := svc.GenerateTransferProof(ctx, params)

		require.Error(t, err)
		assert.Nil(t, proof)
		assert.Nil(t, destCommitments)
		assert.Nil(t, destRFactors)
		assert.Contains(t, err.Error(), "enygma key error")
	})

	t.Run("returns error if enygmaRepository.GetEnygmaByResourceId fails", func(t *testing.T) {
		ctx := context.Background()
		chainID1 := big.NewInt(1)
		chainID2 := big.NewInt(2)
		senderAmount := big.NewInt(1000)

		batch := createTestBatch(chainID1, chainID2, senderAmount)

		params := enygma.TransferProofParams{
			ResourceId:     "test-resource-123",
			AnonymityIndex: 2,
			SenderAmount:   senderAmount,
			BlockNumber:    big.NewInt(100),
			Batches:        []*types.EnygmaTransferBatch{batch},
			TokenAddress:   common.HexToAddress("0xDEADBEEF"),
		}

		kosClient := &MockProofKOSClient{
			GetPaymentSpendKeyFunc: func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
				return createTestPaymentSpendKey(), nil
			},
		}
		proofClient := &MockProofAPIClient{}
		enygmaRepo := &MockProofEnygmaRepository{
			GetEnygmaByResourceIdFunc: func(ctx context.Context, resourceId string) (types.Enygma, error) {
				return types.Enygma{}, errors.New("resource not found")
			},
		}
		enygmaClient := &testEnygmaClient{}
		tracer := &testutils.MockTracer{}

		svc := service.NewEnygmaProofService(
			chainID1, kosClient, proofClient, enygmaRepo, enygmaClient, tracer,
		)

		proof, destCommitments, destRFactors, _, err := svc.GenerateTransferProof(ctx, params)

		require.Error(t, err)
		assert.Nil(t, proof)
		assert.Nil(t, destCommitments)
		assert.Nil(t, destRFactors)
		assert.Contains(t, err.Error(), "resource not found")
	})

	t.Run("returns error if enygmaClient.GetPublicValuesFinalised fails", func(t *testing.T) {
		ctx := context.Background()
		chainID1 := big.NewInt(1)
		chainID2 := big.NewInt(2)
		senderAmount := big.NewInt(1000)

		batch := createTestBatch(chainID1, chainID2, senderAmount)

		params := enygma.TransferProofParams{
			ResourceId:     "test-resource-123",
			AnonymityIndex: 2,
			SenderAmount:   senderAmount,
			BlockNumber:    big.NewInt(100),
			Batches:        []*types.EnygmaTransferBatch{batch},
			TokenAddress:   common.HexToAddress("0xDEADBEEF"),
		}

		kosClient := &MockProofKOSClient{
			GetPaymentSpendKeyFunc: func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
				return createTestPaymentSpendKey(), nil
			},
		}
		proofClient := &MockProofAPIClient{}
		enygmaRepo := &MockProofEnygmaRepository{
			GetEnygmaByResourceIdFunc: func(ctx context.Context, resourceId string) (types.Enygma, error) {
				return createTestEnygmaState(), nil
			},
		}

		enygmaClient := &testEnygmaClient{
			GetPublicValuesFinalizedFunc: func(ctx context.Context, _ common.Address) (*types.EnygmaPublicValues, error) {
				return nil, errors.New("public values error")
			},
		}

		tracer := &testutils.MockTracer{}

		svc := service.NewEnygmaProofService(
			chainID1, kosClient, proofClient, enygmaRepo, enygmaClient, tracer,
		)

		proof, destCommitments, destRFactors, _, err := svc.GenerateTransferProof(ctx, params)

		require.Error(t, err)
		assert.Nil(t, proof)
		assert.Nil(t, destCommitments)
		assert.Nil(t, destRFactors)
		assert.Contains(t, err.Error(), "public values error")
	})

	t.Run("returns error if proofClient.CreateTransferProof fails", func(t *testing.T) {
		ctx := context.Background()
		chainID1 := big.NewInt(1)
		chainID2 := big.NewInt(2)
		senderAmount := big.NewInt(1000)

		// Include sender's chain in destinations (self-transfer) to ensure senderIndex >= 0
		batch1 := createTestBatch(chainID1, chainID1, big.NewInt(500))
		batch2 := createTestBatch(chainID1, chainID2, big.NewInt(500))

		params := enygma.TransferProofParams{
			ResourceId:     "test-resource-123",
			AnonymityIndex: 2,
			SenderAmount:   senderAmount,
			BlockNumber:    big.NewInt(100),
			Batches:        []*types.EnygmaTransferBatch{batch1, batch2},
			TokenAddress:   common.HexToAddress("0xDEADBEEF"),
		}

		kosClient := &MockProofKOSClient{
			GetPaymentSpendKeyFunc: func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
				return createTestPaymentSpendKey(), nil
			},
			GenerateEnygmaSharedSecretsFunc: func(ctx context.Context, in *keys.GenerateEnygmaSharedSecretsRequest, opts ...grpc.CallOption) (*keys.GenerateEnygmaSharedSecretsResponse, error) {
				n := len(in.ChainIDs)
				secrets := make([][]byte, n)
				hashSecrets := make([][]byte, n)
				messageTags := make([][]byte, n)
				for i := 0; i < n; i++ {
					secrets[i] = big.NewInt(int64(1000 + i)).Bytes()
					hashSecrets[i] = big.NewInt(int64(2000 + i)).Bytes()
					messageTags[i] = big.NewInt(int64(3000 + i)).Bytes()
				}
				return &keys.GenerateEnygmaSharedSecretsResponse{Secrets: secrets, HashSecrets: hashSecrets, MessageTags: messageTags}, nil
			},
		}

		proofClient := &MockProofAPIClient{
			CreateTransferProofFunc: func(anonymityIndex int, request types.TransferProofRequest) (types.EnygmaProofResponse, error) {
				return types.EnygmaProofResponse{}, errors.New("proof generation failed")
			},
		}

		enygmaRepo := &MockProofEnygmaRepository{
			GetEnygmaByResourceIdFunc: func(ctx context.Context, resourceId string) (types.Enygma, error) {
				return createTestEnygmaState(), nil
			},
		}

		enygmaClient := &testEnygmaClient{
			GetPublicValuesFinalizedFunc: func(ctx context.Context, _ common.Address) (*types.EnygmaPublicValues, error) {
				return createTestEnygmaPublicValues([]*big.Int{chainID1, chainID2}), nil
			},
		}


		tracer := &testutils.MockTracer{}

		svc := service.NewEnygmaProofService(
			chainID1, kosClient, proofClient, enygmaRepo, enygmaClient, tracer,
		)

		proof, destCommitments, destRFactors, _, err := svc.GenerateTransferProof(ctx, params)

		require.Error(t, err)
		assert.Nil(t, proof)
		assert.Nil(t, destCommitments)
		assert.Nil(t, destRFactors)
		assert.Contains(t, err.Error(), "proof generation failed")
	})

	t.Run("successfully generates a transfer proof", func(t *testing.T) {
		ctx := context.Background()
		chainID1 := big.NewInt(1)
		chainID2 := big.NewInt(2)
		senderAmount := big.NewInt(1000)
		blockNumber := big.NewInt(100)

		// Create two batches to test with multiple transactions
		batch1 := createTestBatch(chainID1, chainID1, big.NewInt(500))
		batch2 := createTestBatch(chainID1, chainID2, big.NewInt(500))

		params := enygma.TransferProofParams{
			ResourceId:     "test-resource-123",
			AnonymityIndex: 2,
			SenderAmount:   senderAmount,
			BlockNumber:    blockNumber,
			Batches:        []*types.EnygmaTransferBatch{batch1, batch2},
			TokenAddress:   common.HexToAddress("0xDEADBEEF"),
		}

		enygmaKey := createTestPaymentSpendKey()
		enygmaState := createTestEnygmaState()
		publicValues := createTestEnygmaPublicValues([]*big.Int{chainID1, chainID2})

		sharedSecrets := []*big.Int{
			big.NewInt(3000),
			big.NewInt(4000),
		}

		expectedProof := types.EnygmaProofResponse{
			PiA: [2]*big.Int{big.NewInt(111), big.NewInt(222)},
			PiB: [2][2]*big.Int{
				{big.NewInt(333), big.NewInt(444)},
				{big.NewInt(555), big.NewInt(666)},
			},
			PiC:          [2]*big.Int{big.NewInt(777), big.NewInt(888)},
			PublicSignal: []*big.Int{big.NewInt(999)},
		}

		kosClient := &MockProofKOSClient{
			GetPaymentSpendKeyFunc: func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
				return enygmaKey, nil
			},
			GenerateEnygmaSharedSecretsFunc: func(ctx context.Context, in *keys.GenerateEnygmaSharedSecretsRequest, opts ...grpc.CallOption) (*keys.GenerateEnygmaSharedSecretsResponse, error) {
				return &keys.GenerateEnygmaSharedSecretsResponse{Secrets: bigIntsToBytes(sharedSecrets), HashSecrets: bigIntsToBytes(sharedSecrets), MessageTags: bigIntsToBytes(sharedSecrets)}, nil
			},
		}

		var capturedProofRequest types.TransferProofRequest

		proofClient := &MockProofAPIClient{
			CreateTransferProofFunc: func(anonymityIndex int, request types.TransferProofRequest) (types.EnygmaProofResponse, error) {
				capturedProofRequest = request
				return expectedProof, nil
			},
		}

		enygmaRepo := &MockProofEnygmaRepository{
			GetEnygmaByResourceIdFunc: func(ctx context.Context, resourceId string) (types.Enygma, error) {
				return enygmaState, nil
			},
		}

		enygmaClient := &testEnygmaClient{
			GetPublicValuesFinalizedFunc: func(ctx context.Context, _ common.Address) (*types.EnygmaPublicValues, error) {
				return publicValues, nil
			},
		}


		tracer := &testutils.MockTracer{}

		svc := service.NewEnygmaProofService(
			chainID1, kosClient, proofClient, enygmaRepo, enygmaClient, tracer,
		)

		proof, destCommitments, destRFactors, _, err := svc.GenerateTransferProof(ctx, params)

		// Verify no error occurred
		require.NoError(t, err)

		// Verify the returned proof is correct
		require.NotNil(t, proof)
		assert.Equal(t, expectedProof, *proof)

		// Verify destination commitments are returned
		require.NotNil(t, destCommitments)
		assert.Len(t, destCommitments, 2)

		// Verify destination random factors are returned
		require.NotNil(t, destRFactors)
		assert.Len(t, destRFactors, 2)

		// Verify all mocks were called exactly once
		assert.Len(t, kosClient.GetPaymentSpendKeyCalls(), 1)
		assert.Len(t, proofClient.CreateTransferProofCalls(), 1)
		assert.Len(t, enygmaRepo.GetEnygmaByResourceIdCalls(), 1)
		assert.Equal(t, enygmaClient.GetCallCount(), 1)
		assert.Len(t, kosClient.GenerateEnygmaSharedSecretsCalls(), 1)

		// Verify captured request has correct structure
		require.NotNil(t, capturedProofRequest.CommonProofRequest)
		assert.Equal(t, params.ResourceId, capturedProofRequest.ResourceId)
		assert.Equal(t, senderAmount, capturedProofRequest.SenderAmount)
		assert.Equal(t, new(big.Int).SetBytes(enygmaKey.SecretKey), capturedProofRequest.SenderSecretKey)
		assert.Equal(t, enygmaState.FinalizedBalance, capturedProofRequest.SenderBalance)
		assert.Equal(t, chainID1, capturedProofRequest.SenderChainId)
		assert.Equal(t, blockNumber, capturedProofRequest.BlockNumber)

		// Verify commitments are valid Point types
		for i, commitment := range capturedProofRequest.DestinationNewCommits {
			assert.NotNil(t, commitment, fmt.Sprintf("Destination commit %d should not be nil", i))
			assert.NotNil(t, commitment.X, fmt.Sprintf("Destination commit %d X should not be nil", i))
			assert.NotNil(t, commitment.Y, fmt.Sprintf("Destination commit %d Y should not be nil", i))
		}

		// Verify nullifier is present and is a valid big.Int
		// Nullifier = Poseidon(arrayHashSecrets[senderIdx], blockNumber) where arrayHashSecrets[senderIdx]
		// is updated in genSecretsAndUpdateDB to be Poseidon(senderSecret, senderSecret)
		require.NotNil(t, capturedProofRequest.Nullifier,
			"Nullifier should be calculated and not nil")
		assert.True(t, capturedProofRequest.Nullifier.Sign() > 0,
			"Nullifier should be a positive big.Int")
		// Verify nullifier is calculated with arrayHashSecrets[senderIndex] and blockNumber
		// senderIndex is 0 because chainID1 is first in destinationChainIDs
		nullifierInputs := []*big.Int{sharedSecrets[0], blockNumber}
		expectedNullifier, err := cryptography.GetPoseidonHashModNumber(nullifierInputs, cryptography.JubJubPrimeGroup)
		require.NoError(t, err)
		assert.Equal(t, expectedNullifier, capturedProofRequest.Nullifier,
			"Nullifier should be hash of arrayHashSecrets[senderIndex] and block number")

		// Verify destination amounts are correct
		// For transfer, the amount to sender chain is (prime - senderAmount) to show outflow
		expectedDestAmounts := []*big.Int{
			cryptography.GetNegative(senderAmount),
			big.NewInt(500),
		}
		assert.Len(t, capturedProofRequest.DestinationAmounts, 2)
		assert.Equal(t, expectedDestAmounts[0], capturedProofRequest.DestinationAmounts[0])
		assert.Equal(t, expectedDestAmounts[1], capturedProofRequest.DestinationAmounts[1])
	})
}

func TestGenerateDepositProof(t *testing.T) {
	t.Run("returns error if anonymity index is invalid - unique chainIDs exceed anonymity index", func(t *testing.T) {
		ctx := context.Background()
		chainID1 := big.NewInt(1)
		chainID2 := big.NewInt(2)
		chainID3 := big.NewInt(3)
		chainID4 := big.NewInt(4)
		senderAmount := big.NewInt(1000)

		// Create batches with 3 unique destination chains
		batch1 := createTestBatch(chainID1, chainID2, senderAmount)
		batch2 := createTestBatch(chainID1, chainID3, senderAmount)
		batch3 := createTestBatch(chainID1, chainID4, senderAmount)

		params := enygma.DepositProofParams{
			ResourceId:        "test-resource-123",
			AnonymityIndex:    2, // Only supports 2 chains, but we have 3 unique destinations
			SenderAmount:      senderAmount,
			BlockNumber:       big.NewInt(100),
			Batches:           []*types.EnygmaTransferBatch{batch1, batch2, batch3},
			TokenAddress:      common.HexToAddress("0xDEADBEEF"),
			DepositCommitment: big.NewInt(999),
			DepositSalt:  big.NewInt(888),
		}

		// Create service with mocks (they won't be called due to early validation)
		kosClient := &MockProofKOSClient{}
		proofClient := &MockProofAPIClient{}
		enygmaRepo := &MockProofEnygmaRepository{}
		enygmaClient := &testEnygmaClient{}
		tracer := &testutils.MockTracer{}

		svc := service.NewEnygmaProofService(
			chainID1, kosClient, proofClient, enygmaRepo, enygmaClient, tracer,
		)

		proof, destCommitments, destRFactors, _, err := svc.GenerateDepositProof(ctx, params)

		require.Error(t, err)
		assert.Nil(t, proof)
		assert.Nil(t, destCommitments)
		assert.Nil(t, destRFactors)
		assert.Contains(t, err.Error(), "Only transactions from 1->1 up to 1->k-1 are supported")
	})

	t.Run("returns error if kosClient.GenerateEnygmaSharedSecrets fails", func(t *testing.T) {
		ctx := context.Background()
		chainID1 := big.NewInt(1)
		chainID2 := big.NewInt(2)
		senderAmount := big.NewInt(1000)

		batch := createTestBatch(chainID1, chainID2, senderAmount)

		params := enygma.DepositProofParams{
			ResourceId:        "test-resource-123",
			AnonymityIndex:    2,
			SenderAmount:      senderAmount,
			BlockNumber:       big.NewInt(100),
			Batches:           []*types.EnygmaTransferBatch{batch},
			TokenAddress:      common.HexToAddress("0xDEADBEEF"),
			DepositCommitment: big.NewInt(999),
			DepositSalt:  big.NewInt(888),
		}

		enygmaClient := &testEnygmaClient{
			GetPublicValuesFinalizedFunc: func(ctx context.Context, _ common.Address) (*types.EnygmaPublicValues, error) {
				return createTestEnygmaPublicValues([]*big.Int{chainID1, chainID2}), nil
			},
		}

		kosClient := &MockProofKOSClient{
			GetPaymentSpendKeyFunc: func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
				return createTestPaymentSpendKey(), nil
			},
			GenerateEnygmaSharedSecretsFunc: func(ctx context.Context, in *keys.GenerateEnygmaSharedSecretsRequest, opts ...grpc.CallOption) (*keys.GenerateEnygmaSharedSecretsResponse, error) {
				return nil, errors.New("shared secrets service down")
			},
		}
		proofClient := &MockProofAPIClient{}
		enygmaRepo := &MockProofEnygmaRepository{
			GetEnygmaByResourceIdFunc: func(ctx context.Context, resourceId string) (types.Enygma, error) {
				return createTestEnygmaState(), nil
			},
		}
		tracer := &testutils.MockTracer{}

		svc := service.NewEnygmaProofService(
			chainID1, kosClient, proofClient, enygmaRepo, enygmaClient, tracer,
		)

		proof, destCommitments, destRFactors, _, err := svc.GenerateDepositProof(ctx, params)

		require.Error(t, err)
		assert.Nil(t, proof)
		assert.Nil(t, destCommitments)
		assert.Nil(t, destRFactors)
		assert.Contains(t, err.Error(), "shared secrets service down")
	})

	t.Run("returns error if kosClient.GetEnygmaKey fails", func(t *testing.T) {
		ctx := context.Background()
		chainID1 := big.NewInt(1)
		chainID2 := big.NewInt(2)
		senderAmount := big.NewInt(1000)

		batch := createTestBatch(chainID1, chainID2, senderAmount)

		params := enygma.DepositProofParams{
			ResourceId:        "test-resource-123",
			AnonymityIndex:    2,
			SenderAmount:      senderAmount,
			BlockNumber:       big.NewInt(100),
			Batches:           []*types.EnygmaTransferBatch{batch},
			TokenAddress:      common.HexToAddress("0xDEADBEEF"),
			DepositCommitment: big.NewInt(999),
			DepositSalt:  big.NewInt(888),
		}

		kosClient := &MockProofKOSClient{
			GetPaymentSpendKeyFunc: func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
				return nil, errors.New("enygma key error")
			},
		}
		proofClient := &MockProofAPIClient{}
		enygmaRepo := &MockProofEnygmaRepository{
			GetEnygmaByResourceIdFunc: func(ctx context.Context, resourceId string) (types.Enygma, error) {
				return createTestEnygmaState(), nil
			},
		}
		enygmaClient := &testEnygmaClient{}
		tracer := &testutils.MockTracer{}

		svc := service.NewEnygmaProofService(
			chainID1, kosClient, proofClient, enygmaRepo, enygmaClient, tracer,
		)

		proof, destCommitments, destRFactors, _, err := svc.GenerateDepositProof(ctx, params)

		require.Error(t, err)
		assert.Nil(t, proof)
		assert.Nil(t, destCommitments)
		assert.Nil(t, destRFactors)
		assert.Contains(t, err.Error(), "enygma key error")
	})

	t.Run("returns error if enygmaRepository.GetEnygmaByResourceId fails", func(t *testing.T) {
		ctx := context.Background()
		chainID1 := big.NewInt(1)
		chainID2 := big.NewInt(2)
		senderAmount := big.NewInt(1000)

		batch := createTestBatch(chainID1, chainID2, senderAmount)

		params := enygma.DepositProofParams{
			ResourceId:        "test-resource-123",
			AnonymityIndex:    2,
			SenderAmount:      senderAmount,
			BlockNumber:       big.NewInt(100),
			Batches:           []*types.EnygmaTransferBatch{batch},
			TokenAddress:      common.HexToAddress("0xDEADBEEF"),
			DepositCommitment: big.NewInt(999),
			DepositSalt:  big.NewInt(888),
		}

		kosClient := &MockProofKOSClient{
			GetPaymentSpendKeyFunc: func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
				return createTestPaymentSpendKey(), nil
			},
		}
		proofClient := &MockProofAPIClient{}
		enygmaRepo := &MockProofEnygmaRepository{
			GetEnygmaByResourceIdFunc: func(ctx context.Context, resourceId string) (types.Enygma, error) {
				return types.Enygma{}, errors.New("resource not found")
			},
		}
		enygmaClient := &testEnygmaClient{}
		tracer := &testutils.MockTracer{}

		svc := service.NewEnygmaProofService(
			chainID1, kosClient, proofClient, enygmaRepo, enygmaClient, tracer,
		)

		proof, destCommitments, destRFactors, _, err := svc.GenerateDepositProof(ctx, params)

		require.Error(t, err)
		assert.Nil(t, proof)
		assert.Nil(t, destCommitments)
		assert.Nil(t, destRFactors)
		assert.Contains(t, err.Error(), "resource not found")
	})

	t.Run("returns error if enygmaClient.GetPublicValuesFinalised fails", func(t *testing.T) {
		ctx := context.Background()
		chainID1 := big.NewInt(1)
		chainID2 := big.NewInt(2)
		senderAmount := big.NewInt(1000)

		batch := createTestBatch(chainID1, chainID2, senderAmount)

		params := enygma.DepositProofParams{
			ResourceId:        "test-resource-123",
			AnonymityIndex:    2,
			SenderAmount:      senderAmount,
			BlockNumber:       big.NewInt(100),
			Batches:           []*types.EnygmaTransferBatch{batch},
			TokenAddress:      common.HexToAddress("0xDEADBEEF"),
			DepositCommitment: big.NewInt(999),
			DepositSalt:  big.NewInt(888),
		}

		kosClient := &MockProofKOSClient{
			GetPaymentSpendKeyFunc: func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
				return createTestPaymentSpendKey(), nil
			},
		}
		proofClient := &MockProofAPIClient{}
		enygmaRepo := &MockProofEnygmaRepository{
			GetEnygmaByResourceIdFunc: func(ctx context.Context, resourceId string) (types.Enygma, error) {
				return createTestEnygmaState(), nil
			},
		}

		enygmaClient := &testEnygmaClient{
			GetPublicValuesFinalizedFunc: func(ctx context.Context, _ common.Address) (*types.EnygmaPublicValues, error) {
				return nil, errors.New("public values error")
			},
		}

		tracer := &testutils.MockTracer{}

		svc := service.NewEnygmaProofService(
			chainID1, kosClient, proofClient, enygmaRepo, enygmaClient, tracer,
		)

		proof, destCommitments, destRFactors, _, err := svc.GenerateDepositProof(ctx, params)

		require.Error(t, err)
		assert.Nil(t, proof)
		assert.Nil(t, destCommitments)
		assert.Nil(t, destRFactors)
		assert.Contains(t, err.Error(), "public values error")
	})

	t.Run("returns error if proofClient.CreateDepositProof fails", func(t *testing.T) {
		ctx := context.Background()
		chainID1 := big.NewInt(1)
		chainID2 := big.NewInt(2)
		senderAmount := big.NewInt(1000)

		// Include sender's chain in destinations (self-transfer) to ensure senderIndex >= 0
		batch1 := createTestBatch(chainID1, chainID1, big.NewInt(500))
		batch2 := createTestBatch(chainID1, chainID2, big.NewInt(500))

		params := enygma.DepositProofParams{
			ResourceId:        "test-resource-123",
			AnonymityIndex:    2,
			SenderAmount:      senderAmount,
			BlockNumber:       big.NewInt(100),
			Batches:           []*types.EnygmaTransferBatch{batch1, batch2},
			TokenAddress:      common.HexToAddress("0xDEADBEEF"),
			DepositCommitment: big.NewInt(999),
			DepositSalt:  big.NewInt(888),
		}

		sharedSecrets := []*big.Int{
			big.NewInt(3000),
			big.NewInt(4000),
		}

		kosClient := &MockProofKOSClient{
			GetPaymentSpendKeyFunc: func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
				return createTestPaymentSpendKey(), nil
			},
			GenerateEnygmaSharedSecretsFunc: func(ctx context.Context, in *keys.GenerateEnygmaSharedSecretsRequest, opts ...grpc.CallOption) (*keys.GenerateEnygmaSharedSecretsResponse, error) {
				return &keys.GenerateEnygmaSharedSecretsResponse{Secrets: bigIntsToBytes(sharedSecrets), HashSecrets: bigIntsToBytes(sharedSecrets), MessageTags: bigIntsToBytes(sharedSecrets)}, nil
			},
		}

		proofClient := &MockProofAPIClient{
			CreateDepositProofFunc: func(anonymityIndex int, request types.DepositProofRequest) (types.EnygmaProofResponse, error) {
				return types.EnygmaProofResponse{}, errors.New("proof generation failed")
			},
		}

		enygmaRepo := &MockProofEnygmaRepository{
			GetEnygmaByResourceIdFunc: func(ctx context.Context, resourceId string) (types.Enygma, error) {
				return createTestEnygmaState(), nil
			},
		}

		enygmaClient := &testEnygmaClient{
			GetPublicValuesFinalizedFunc: func(ctx context.Context, _ common.Address) (*types.EnygmaPublicValues, error) {
				return createTestEnygmaPublicValues([]*big.Int{chainID1, chainID2}), nil
			},
		}


		tracer := &testutils.MockTracer{}

		svc := service.NewEnygmaProofService(
			chainID1, kosClient, proofClient, enygmaRepo, enygmaClient, tracer,
		)

		proof, destCommitments, destRFactors, _, err := svc.GenerateDepositProof(ctx, params)

		require.Error(t, err)
		assert.Nil(t, proof)
		assert.Nil(t, destCommitments)
		assert.Nil(t, destRFactors)
		assert.Contains(t, err.Error(), "proof generation failed")
	})

	t.Run("successfully generates a deposit proof", func(t *testing.T) {
		ctx := context.Background()
		chainID1 := big.NewInt(1)
		chainID2 := big.NewInt(2)
		senderAmount := big.NewInt(1000)
		blockNumber := big.NewInt(100)
		depositCommitment := big.NewInt(777)
		depositSalt := big.NewInt(666)

		// Create two batches to test with multiple transactions
		batch1 := createTestBatch(chainID1, chainID1, big.NewInt(500))
		batch2 := createTestBatch(chainID1, chainID2, big.NewInt(500))

		params := enygma.DepositProofParams{
			ResourceId:        "test-resource-123",
			AnonymityIndex:    2,
			SenderAmount:      senderAmount,
			BlockNumber:       blockNumber,
			Batches:           []*types.EnygmaTransferBatch{batch1, batch2},
			TokenAddress:      common.HexToAddress("0xDEADBEEF"),
			DepositCommitment: depositCommitment,
			DepositSalt:  depositSalt,
		}

		enygmaKey := createTestPaymentSpendKey()
		enygmaState := createTestEnygmaState()
		publicValues := createTestEnygmaPublicValues([]*big.Int{chainID1, chainID2})

		sharedSecrets := []*big.Int{
			big.NewInt(3000),
			big.NewInt(4000),
		}

		expectedProof := types.EnygmaProofResponse{
			PiA: [2]*big.Int{big.NewInt(111), big.NewInt(222)},
			PiB: [2][2]*big.Int{
				{big.NewInt(333), big.NewInt(444)},
				{big.NewInt(555), big.NewInt(666)},
			},
			PiC:          [2]*big.Int{big.NewInt(777), big.NewInt(888)},
			PublicSignal: []*big.Int{big.NewInt(999)},
		}

		kosClient := &MockProofKOSClient{
			GetPaymentSpendKeyFunc: func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
				return enygmaKey, nil
			},
			GenerateEnygmaSharedSecretsFunc: func(ctx context.Context, in *keys.GenerateEnygmaSharedSecretsRequest, opts ...grpc.CallOption) (*keys.GenerateEnygmaSharedSecretsResponse, error) {
				return &keys.GenerateEnygmaSharedSecretsResponse{Secrets: bigIntsToBytes(sharedSecrets), HashSecrets: bigIntsToBytes(sharedSecrets), MessageTags: bigIntsToBytes(sharedSecrets)}, nil
			},
		}

		var capturedDepositProofRequest types.DepositProofRequest

		proofClient := &MockProofAPIClient{
			CreateDepositProofFunc: func(anonymityIndex int, request types.DepositProofRequest) (types.EnygmaProofResponse, error) {
				capturedDepositProofRequest = request
				return expectedProof, nil
			},
		}

		enygmaRepo := &MockProofEnygmaRepository{
			GetEnygmaByResourceIdFunc: func(ctx context.Context, resourceId string) (types.Enygma, error) {
				return enygmaState, nil
			},
		}

		enygmaClient := &testEnygmaClient{
			GetPublicValuesFinalizedFunc: func(ctx context.Context, _ common.Address) (*types.EnygmaPublicValues, error) {
				return publicValues, nil
			},
		}


		tracer := &testutils.MockTracer{}

		svc := service.NewEnygmaProofService(
			chainID1, kosClient, proofClient, enygmaRepo, enygmaClient, tracer,
		)

		proof, destCommitments, destRFactors, _, err := svc.GenerateDepositProof(ctx, params)

		// Verify no error occurred
		require.NoError(t, err)

		// Verify the returned proof is correct
		require.NotNil(t, proof)
		assert.Equal(t, expectedProof, *proof)

		// Verify destination commitments are returned
		require.NotNil(t, destCommitments)
		assert.Len(t, destCommitments, 2)

		// Verify destination random factors are returned
		require.NotNil(t, destRFactors)
		assert.Len(t, destRFactors, 2)

		// Verify all mocks were called exactly once
		assert.Len(t, kosClient.GetPaymentSpendKeyCalls(), 1)
		assert.Len(t, proofClient.CreateDepositProofCalls(), 1)
		assert.Len(t, enygmaRepo.GetEnygmaByResourceIdCalls(), 1)
		assert.Equal(t, enygmaClient.GetCallCount(), 1)
		assert.Len(t, kosClient.GenerateEnygmaSharedSecretsCalls(), 1)

		// Verify captured request has correct structure
		require.NotNil(t, capturedDepositProofRequest.CommonProofRequest)
		assert.Equal(t, params.ResourceId, capturedDepositProofRequest.ResourceId)
		assert.Equal(t, senderAmount, capturedDepositProofRequest.SenderAmount)
		assert.Equal(t, new(big.Int).SetBytes(enygmaKey.SecretKey), capturedDepositProofRequest.SenderSecretKey)
		assert.Equal(t, enygmaState.FinalizedBalance, capturedDepositProofRequest.SenderBalance)
		assert.Equal(t, chainID1, capturedDepositProofRequest.SenderChainId)
		assert.Equal(t, blockNumber, capturedDepositProofRequest.BlockNumber)

		// Verify deposit-specific fields are included
		assert.Equal(t, depositCommitment, capturedDepositProofRequest.DepositCommitment,
			"DepositCommitment should be passed to the proof request")
		assert.Equal(t, depositSalt, capturedDepositProofRequest.DepositSalt,
			"DepositSalt should be passed to the proof request")
		assert.Equal(t, params.TokenAddress, capturedDepositProofRequest.TokenAddress,
			"TokenAddress should be passed to the proof request")

		// Verify commitments are valid Point types
		for i, commitment := range capturedDepositProofRequest.DestinationNewCommits {
			assert.NotNil(t, commitment, fmt.Sprintf("Destination commit %d should not be nil", i))
			assert.NotNil(t, commitment.X, fmt.Sprintf("Destination commit %d X should not be nil", i))
			assert.NotNil(t, commitment.Y, fmt.Sprintf("Destination commit %d Y should not be nil", i))
		}

		// Verify nullifier is present and is a valid big.Int
		// Nullifier = Poseidon(arrayHashSecrets[senderIdx], blockNumber) where arrayHashSecrets[senderIdx]
		// is updated in genSecretsAndUpdateDB to be Poseidon(senderSecret, senderSecret)
		require.NotNil(t, capturedDepositProofRequest.Nullifier,
			"Nullifier should be calculated and not nil")
		assert.True(t, capturedDepositProofRequest.Nullifier.Sign() > 0,
			"Nullifier should be a positive big.Int")
		// Verify nullifier is calculated with blockNumber and arrayHashSecrets[senderIdx]
		// arrayHashSecrets[senderIdx] = Poseidon(senderSecret, senderSecret)
		nullifierInputs := []*big.Int{sharedSecrets[0], blockNumber}
		expectedNullifier, err := cryptography.GetPoseidonHashModNumber(nullifierInputs, cryptography.JubJubPrimeGroup)
		require.NoError(t, err)
		assert.Equal(t, expectedNullifier, capturedDepositProofRequest.Nullifier,
			"Nullifier should be hash of arrayHashSecrets[senderIdx] and block number")

		// Verify destination amounts are correct
		// For deposit, the amount to sender chain is (prime - senderAmount) to show outflow
		expectedDestAmounts := []*big.Int{
			cryptography.GetNegative(senderAmount),
			big.NewInt(500),
		}
		assert.Len(t, capturedDepositProofRequest.DestinationAmounts, 2)
		assert.Equal(t, expectedDestAmounts[0], capturedDepositProofRequest.DestinationAmounts[0])
		assert.Equal(t, expectedDestAmounts[1], capturedDepositProofRequest.DestinationAmounts[1])
	})
}

func TestGenerateWithdrawProof(t *testing.T) {
	t.Run("returns error if anonymity index is invalid - unique chainIDs exceed anonymity index", func(t *testing.T) {
		// Setup
		anonymityIndex := 1
		senderAmount := big.NewInt(1000)
		blockNumber := big.NewInt(12345)

		// Create batches with 2 unique destination chains (exceeds anonymity index of 1)
		chainID1 := big.NewInt(1)
		chainID2 := big.NewInt(2)
		batches := []*types.EnygmaTransferBatch{
			createTestBatch(chainID1, chainID1, big.NewInt(100)),
			createTestBatch(chainID1, chainID2, big.NewInt(200)),
		}

		// Create minimal params
		params := enygma.WithdrawProofParams{
			ResourceId:         "test-resource",
			AnonymityIndex:     anonymityIndex,
			SenderAmount:       senderAmount,
			BlockNumber:        blockNumber,
			Batches:            batches,
			TokenAddress:       common.Address{},
			DepositCommitments: []*big.Int{big.NewInt(111)},
			DepositSecretKeys:  []*big.Int{big.NewInt(222)},
			DepositAmounts:     []*big.Int{big.NewInt(333)},
		}

		tracer := &testutils.MockTracer{}
		proofService := service.NewEnygmaProofService(
			chainID1,
			&MockProofKOSClient{},
			&MockProofAPIClient{},
			&MockProofEnygmaRepository{},
			&testEnygmaClient{},
			tracer,
		)

		proof, commitments, randomFactors, _, err := proofService.GenerateWithdrawProof(context.Background(), params)

		assert.Error(t, err)
		assert.Nil(t, proof)
		assert.Nil(t, commitments)
		assert.Nil(t, randomFactors)
		assert.ErrorContains(t, err, "Only transactions from 1->1 up to 1->k-1 are supported")
	})

	t.Run("returns error if kosClient.GenerateEnygmaSharedSecrets fails", func(t *testing.T) {
		anonymityIndex := 2
		senderAmount := big.NewInt(1000)
		blockNumber := big.NewInt(12345)

		chainID1 := big.NewInt(1)
		batches := []*types.EnygmaTransferBatch{
			createTestBatch(chainID1, chainID1, big.NewInt(100)),
		}

		params := enygma.WithdrawProofParams{
			ResourceId:         "test-resource",
			AnonymityIndex:     anonymityIndex,
			SenderAmount:       senderAmount,
			BlockNumber:        blockNumber,
			Batches:            batches,
			TokenAddress:       common.Address{},
			DepositCommitments: []*big.Int{big.NewInt(111)},
			DepositSecretKeys:  []*big.Int{big.NewInt(222)},
			DepositAmounts:     []*big.Int{big.NewInt(333)},
		}

		kosClient := &MockProofKOSClient{
			GetPaymentSpendKeyFunc: func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
				return createTestPaymentSpendKey(), nil
			},
			GenerateEnygmaSharedSecretsFunc: func(ctx context.Context, in *keys.GenerateEnygmaSharedSecretsRequest, opts ...grpc.CallOption) (*keys.GenerateEnygmaSharedSecretsResponse, error) {
				return nil, fmt.Errorf("kos shared secrets error")
			},
		}

		enygmaClient := &testEnygmaClient{
			GetPublicValuesFinalizedFunc: func(ctx context.Context, _ common.Address) (*types.EnygmaPublicValues, error) {
				return createTestEnygmaPublicValues([]*big.Int{chainID1}), nil
			},
		}


		enygmaRepo := &MockProofEnygmaRepository{
			GetEnygmaByResourceIdFunc: func(ctx context.Context, resourceId string) (types.Enygma, error) {
				return createTestEnygmaState(), nil
			},
		}

		tracer := &testutils.MockTracer{}
		proofService := service.NewEnygmaProofService(
			chainID1,
			kosClient,
			&MockProofAPIClient{},
			enygmaRepo,
			enygmaClient,
			tracer,
		)

		proof, commitments, randomFactors, _, err := proofService.GenerateWithdrawProof(context.Background(), params)

		assert.Error(t, err)
		assert.Nil(t, proof)
		assert.Nil(t, commitments)
		assert.Nil(t, randomFactors)
		assert.ErrorContains(t, err, "kos shared secrets error")
	})

	t.Run("returns error if kosClient.GetEnygmaKey fails", func(t *testing.T) {
		anonymityIndex := 2
		senderAmount := big.NewInt(1000)
		blockNumber := big.NewInt(12345)

		chainID1 := big.NewInt(1)
		batches := []*types.EnygmaTransferBatch{
			createTestBatch(chainID1, chainID1, big.NewInt(100)),
		}

		params := enygma.WithdrawProofParams{
			ResourceId:         "test-resource",
			AnonymityIndex:     anonymityIndex,
			SenderAmount:       senderAmount,
			BlockNumber:        blockNumber,
			Batches:            batches,
			TokenAddress:       common.Address{},
			DepositCommitments: []*big.Int{big.NewInt(111)},
			DepositSecretKeys:  []*big.Int{big.NewInt(222)},
			DepositAmounts:     []*big.Int{big.NewInt(333)},
		}

		kosClient := &MockProofKOSClient{
			GetPaymentSpendKeyFunc: func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
				return nil, fmt.Errorf("kos client error")
			},
		}

		tracer := &testutils.MockTracer{}
		proofService := service.NewEnygmaProofService(
			chainID1,
			kosClient,
			&MockProofAPIClient{},
			&MockProofEnygmaRepository{},
			&testEnygmaClient{},
			tracer,
		)

		proof, commitments, randomFactors, _, err := proofService.GenerateWithdrawProof(context.Background(), params)

		assert.Error(t, err)
		assert.Nil(t, proof)
		assert.Nil(t, commitments)
		assert.Nil(t, randomFactors)
		assert.ErrorContains(t, err, "kos client error")
	})

	t.Run("returns error if enygmaRepository.GetEnygmaByResourceId fails", func(t *testing.T) {
		anonymityIndex := 2
		senderAmount := big.NewInt(1000)
		blockNumber := big.NewInt(12345)

		chainID1 := big.NewInt(1)
		batches := []*types.EnygmaTransferBatch{
			createTestBatch(chainID1, chainID1, big.NewInt(100)),
		}

		params := enygma.WithdrawProofParams{
			ResourceId:         "test-resource",
			AnonymityIndex:     anonymityIndex,
			SenderAmount:       senderAmount,
			BlockNumber:        blockNumber,
			Batches:            batches,
			TokenAddress:       common.Address{},
			DepositCommitments: []*big.Int{big.NewInt(111)},
			DepositSecretKeys:  []*big.Int{big.NewInt(222)},
			DepositAmounts:     []*big.Int{big.NewInt(333)},
		}

		kosClient := &MockProofKOSClient{
			GetPaymentSpendKeyFunc: func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
				return createTestPaymentSpendKey(), nil
			},
		}

		enygmaRepo := &MockProofEnygmaRepository{
			GetEnygmaByResourceIdFunc: func(ctx context.Context, resourceId string) (types.Enygma, error) {
				return types.Enygma{}, fmt.Errorf("repository error")
			},
		}

		tracer := &testutils.MockTracer{}
		proofService := service.NewEnygmaProofService(
			chainID1,
			kosClient,
			&MockProofAPIClient{},
			enygmaRepo,
			&testEnygmaClient{},
			tracer,
		)

		proof, commitments, randomFactors, _, err := proofService.GenerateWithdrawProof(context.Background(), params)

		assert.Error(t, err)
		assert.Nil(t, proof)
		assert.Nil(t, commitments)
		assert.Nil(t, randomFactors)
		assert.ErrorContains(t, err, "repository error")
	})

	t.Run("returns error if enygmaClient.GetPublicValuesFinalised fails", func(t *testing.T) {
		anonymityIndex := 2
		senderAmount := big.NewInt(1000)
		blockNumber := big.NewInt(12345)

		chainID1 := big.NewInt(1)
		batches := []*types.EnygmaTransferBatch{
			createTestBatch(chainID1, chainID1, big.NewInt(100)),
		}

		params := enygma.WithdrawProofParams{
			ResourceId:         "test-resource",
			AnonymityIndex:     anonymityIndex,
			SenderAmount:       senderAmount,
			BlockNumber:        blockNumber,
			Batches:            batches,
			TokenAddress:       common.Address{},
			DepositCommitments: []*big.Int{big.NewInt(111)},
			DepositSecretKeys:  []*big.Int{big.NewInt(222)},
			DepositAmounts:     []*big.Int{big.NewInt(333)},
		}

		kosClient := &MockProofKOSClient{
			GetPaymentSpendKeyFunc: func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
				return createTestPaymentSpendKey(), nil
			},
		}

		enygmaRepo := &MockProofEnygmaRepository{
			GetEnygmaByResourceIdFunc: func(ctx context.Context, resourceId string) (types.Enygma, error) {
				return createTestEnygmaState(), nil
			},
		}

		enygmaClient := &testEnygmaClient{
			GetPublicValuesFinalizedFunc: func(ctx context.Context, _ common.Address) (*types.EnygmaPublicValues, error) {
				return nil, fmt.Errorf("client error")
			},
		}


		tracer := &testutils.MockTracer{}
		proofService := service.NewEnygmaProofService(
			chainID1,
			kosClient,
			&MockProofAPIClient{},
			enygmaRepo,
			enygmaClient,
			tracer,
		)

		proof, commitments, randomFactors, _, err := proofService.GenerateWithdrawProof(context.Background(), params)

		assert.Error(t, err)
		assert.Nil(t, proof)
		assert.Nil(t, commitments)
		assert.Nil(t, randomFactors)
		assert.ErrorContains(t, err, "client error")
	})

	t.Run("returns error if proofClient.CreateWithdrawProof fails", func(t *testing.T) {
		anonymityIndex := 2
		senderAmount := big.NewInt(1000)
		blockNumber := big.NewInt(12345)

		chainID1 := big.NewInt(1)
		batches := []*types.EnygmaTransferBatch{
			createTestBatch(chainID1, chainID1, big.NewInt(100)),
		}

		params := enygma.WithdrawProofParams{
			ResourceId:         "test-resource",
			AnonymityIndex:     anonymityIndex,
			SenderAmount:       senderAmount,
			BlockNumber:        blockNumber,
			Batches:            batches,
			TokenAddress:       common.Address{},
			DepositCommitments: []*big.Int{big.NewInt(111)},
			DepositSecretKeys:  []*big.Int{big.NewInt(222)},
			DepositAmounts:     []*big.Int{big.NewInt(333)},
		}

		sharedSecrets := []*big.Int{big.NewInt(1000)}

		kosClient := &MockProofKOSClient{
			GetPaymentSpendKeyFunc: func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
				return createTestPaymentSpendKey(), nil
			},
			GenerateEnygmaSharedSecretsFunc: func(ctx context.Context, in *keys.GenerateEnygmaSharedSecretsRequest, opts ...grpc.CallOption) (*keys.GenerateEnygmaSharedSecretsResponse, error) {
				return &keys.GenerateEnygmaSharedSecretsResponse{Secrets: bigIntsToBytes(sharedSecrets), HashSecrets: bigIntsToBytes(sharedSecrets), MessageTags: bigIntsToBytes(sharedSecrets)}, nil
			},
		}

		enygmaRepo := &MockProofEnygmaRepository{
			GetEnygmaByResourceIdFunc: func(ctx context.Context, resourceId string) (types.Enygma, error) {
				return createTestEnygmaState(), nil
			},
		}

		enygmaClient := &testEnygmaClient{
			GetPublicValuesFinalizedFunc: func(ctx context.Context, _ common.Address) (*types.EnygmaPublicValues, error) {
				return createTestEnygmaPublicValues([]*big.Int{chainID1}), nil
			},
		}


		proofClient := &MockProofAPIClient{
			CreateWithdrawProofFunc: func(anonymityIndex int, request types.WithdrawProofRequest) (types.EnygmaProofResponse, error) {
				return types.EnygmaProofResponse{}, fmt.Errorf("proof client error")
			},
		}

		tracer := &testutils.MockTracer{}
		proofService := service.NewEnygmaProofService(
			chainID1,
			kosClient,
			proofClient,
			enygmaRepo,
			enygmaClient,
			tracer,
		)

		proof, commitments, randomFactors, _, err := proofService.GenerateWithdrawProof(context.Background(), params)

		assert.Error(t, err)
		assert.Nil(t, proof)
		assert.Nil(t, commitments)
		assert.Nil(t, randomFactors)
		assert.ErrorContains(t, err, "proof client error")
	})

	t.Run("successfully generates a withdraw proof", func(t *testing.T) {
		chainID1 := big.NewInt(1)
		chainID2 := big.NewInt(2)
		anonymityIndex := 2
		senderAmount := big.NewInt(1000)
		blockNumber := big.NewInt(100)
		tokenAddress := common.HexToAddress("0xDEADBEEF")

		// Create two batches to test with multiple transactions
		batch1 := createTestBatch(chainID1, chainID1, senderAmount)
		batch2 := createTestBatch(chainID1, chainID2, big.NewInt(0))
		batches := []*types.EnygmaTransferBatch{batch1, batch2}

		depositCommitments := []*big.Int{big.NewInt(111), big.NewInt(222)}
		depositSecretKeys := []*big.Int{big.NewInt(333), big.NewInt(444)}
		depositAmounts := []*big.Int{big.NewInt(555), big.NewInt(666)}

		params := enygma.WithdrawProofParams{
			ResourceId:         "test-resource",
			AnonymityIndex:     anonymityIndex,
			SenderAmount:       senderAmount,
			BlockNumber:        blockNumber,
			Batches:            batches,
			TokenAddress:       tokenAddress,
			DepositCommitments: depositCommitments,
			DepositSecretKeys:  depositSecretKeys,
			DepositAmounts:     depositAmounts,
		}

		expectedProof := types.EnygmaProofResponse{
			PiA: [2]*big.Int{big.NewInt(1), big.NewInt(2)},
			PiB: [2][2]*big.Int{
				{big.NewInt(3), big.NewInt(4)},
				{big.NewInt(5), big.NewInt(6)},
			},
			PiC:          [2]*big.Int{big.NewInt(7), big.NewInt(8)},
			PublicSignal: []*big.Int{big.NewInt(9), big.NewInt(10)},
		}

		var capturedWithdrawProofRequest types.WithdrawProofRequest

		sharedSecrets := []*big.Int{
			big.NewInt(3000),
			big.NewInt(4000),
		}

		kosClient := &MockProofKOSClient{
			GetPaymentSpendKeyFunc: func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
				return createTestPaymentSpendKey(), nil
			},
			GenerateEnygmaSharedSecretsFunc: func(ctx context.Context, in *keys.GenerateEnygmaSharedSecretsRequest, opts ...grpc.CallOption) (*keys.GenerateEnygmaSharedSecretsResponse, error) {
				return &keys.GenerateEnygmaSharedSecretsResponse{Secrets: bigIntsToBytes(sharedSecrets), HashSecrets: bigIntsToBytes(sharedSecrets), MessageTags: bigIntsToBytes(sharedSecrets)}, nil
			},
		}

		enygmaRepo := &MockProofEnygmaRepository{
			GetEnygmaByResourceIdFunc: func(ctx context.Context, resourceId string) (types.Enygma, error) {
				return createTestEnygmaState(), nil
			},
		}

		enygmaClient := &testEnygmaClient{
			GetPublicValuesFinalizedFunc: func(ctx context.Context, _ common.Address) (*types.EnygmaPublicValues, error) {
				return createTestEnygmaPublicValues([]*big.Int{chainID1, chainID2}), nil
			},
		}


		proofClient := &MockProofAPIClient{
			CreateWithdrawProofFunc: func(anonymityIndex int, request types.WithdrawProofRequest) (types.EnygmaProofResponse, error) {
				capturedWithdrawProofRequest = request
				return expectedProof, nil
			},
		}

		tracer := &testutils.MockTracer{}
		proofService := service.NewEnygmaProofService(
			chainID1,
			kosClient,
			proofClient,
			enygmaRepo,
			enygmaClient,
			tracer,
		)

		proof, commitments, randomFactors, _, err := proofService.GenerateWithdrawProof(context.Background(), params)

		require.NoError(t, err)
		require.NotNil(t, proof)
		assert.Equal(t, &expectedProof, proof)

		// Verify destination commitments are returned
		require.NotNil(t, commitments)
		assert.Len(t, commitments, 2)

		// Verify destination random factors are returned
		require.NotNil(t, randomFactors)
		assert.Len(t, randomFactors, 2)

		// Verify all mocks were called exactly once
		assert.Len(t, kosClient.GetPaymentSpendKeyCalls(), 1)
		assert.Len(t, proofClient.CreateWithdrawProofCalls(), 1)
		assert.Len(t, enygmaRepo.GetEnygmaByResourceIdCalls(), 1)
		assert.Equal(t, enygmaClient.GetCallCount(), 1)
		assert.Len(t, kosClient.GenerateEnygmaSharedSecretsCalls(), 1)

		// Verify captured request has correct structure
		require.NotNil(t, capturedWithdrawProofRequest.CommonProofRequest)
		assert.Equal(t, params.ResourceId, capturedWithdrawProofRequest.ResourceId)
		assert.Equal(t, senderAmount, capturedWithdrawProofRequest.SenderAmount)
		enygmaKey, _ := kosClient.GetPaymentSpendKeyFunc(context.Background(), &keys.GetPaymentSpendKeyRequest{})
		assert.Equal(t, new(big.Int).SetBytes(enygmaKey.SecretKey), capturedWithdrawProofRequest.SenderSecretKey)
		assert.Equal(
			t,
			createTestEnygmaState().FinalizedBalance,
			capturedWithdrawProofRequest.SenderBalance,
		)
		assert.Equal(t, chainID1, capturedWithdrawProofRequest.SenderChainId)
		assert.Equal(t, blockNumber, capturedWithdrawProofRequest.BlockNumber)

		// Verify withdraw-specific fields are included
		assert.Equal(t, params.TokenAddress, capturedWithdrawProofRequest.TokenAddress,
			"TokenAddress should be passed to the withdraw proof request")
		assert.Equal(t, depositCommitments, capturedWithdrawProofRequest.DepositCommitments,
			"DepositCommitments should be passed to the withdraw proof request")
		assert.Equal(t, depositSecretKeys, capturedWithdrawProofRequest.DepositSecretKeys,
			"DepositSecretKeys should be passed to the withdraw proof request")
		assert.Equal(t, depositAmounts, capturedWithdrawProofRequest.DepositAmounts,
			"DepositAmounts should be passed to the withdraw proof request")

		// Verify commitments are valid Point types
		for i, commitment := range capturedWithdrawProofRequest.DestinationNewCommits {
			assert.NotNil(t, commitment, fmt.Sprintf("Destination commit %d should not be nil", i))
			assert.NotNil(t, commitment.X, fmt.Sprintf("Destination commit %d X should not be nil", i))
			assert.NotNil(t, commitment.Y, fmt.Sprintf("Destination commit %d Y should not be nil", i))
		}

		// Verify nullifier is present and is a valid big.Int
		// Nullifier = Poseidon(arrayHashSecrets[senderIdx], blockNumber) where arrayHashSecrets[senderIdx]
		// is updated in genSecretsAndUpdateDB to be Poseidon(senderSecret, senderSecret)
		require.NotNil(t, capturedWithdrawProofRequest.Nullifier,
			"Nullifier should be calculated and not nil")
		assert.True(t, capturedWithdrawProofRequest.Nullifier.Sign() > 0,
			"Nullifier should be a positive big.Int")
		// Verify nullifier is calculated with blockNumber and arrayHashSecrets[senderIdx]
		// arrayHashSecrets[senderIdx] = Poseidon(senderSecret, senderSecret)
		nullifierInputs := []*big.Int{sharedSecrets[0], blockNumber}
		expectedNullifier, err := cryptography.GetPoseidonHashModNumber(nullifierInputs, cryptography.JubJubPrimeGroup)
		require.NoError(t, err)
		assert.Equal(t, expectedNullifier, capturedWithdrawProofRequest.Nullifier,
			"Nullifier should be hash of arrayHashSecrets[senderIdx] and block number")

		// Verify destination amounts are correct
		// For withdraw, the amount to sender chain is senderAmount (not negated like in transfer/deposit)
		expectedDestAmounts := []*big.Int{
			senderAmount,
			big.NewInt(0),
		}
		assert.Len(t, capturedWithdrawProofRequest.DestinationAmounts, 2)
		assert.Equal(t, expectedDestAmounts[0], capturedWithdrawProofRequest.DestinationAmounts[0])
		assert.Equal(t, expectedDestAmounts[1], capturedWithdrawProofRequest.DestinationAmounts[1])
	})
}
