package service_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	keyspb "github.com/raylsnetwork/rayls-sovereign-relayer/cts/gen/proto/keys"
	"github.com/raylsnetwork/rayls-sovereign-relayer/dvp"
	"github.com/raylsnetwork/rayls-sovereign-relayer/dvp/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/dvp/testdata"
	merkle "github.com/raylsnetwork/rayls-sovereign-relayer/merkle-service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// proofFixture bundles the mocks needed by every ProofService test.
// Defaults are the happy-path scenario; individual tests override fields.
type proofFixture struct {
	config service.ProofServiceConfig
	merkle *proofMerkleClientMock
	kos    *proofKeysClientMock
	api    *proofAPIClientMock
	calc   *proofCommitmentCalculatorMock
	repo   *proofDepositRepositoryMock
	txMgr  *proofTxManagerMock
	viewPK []byte
}

// defaultProofSignal is the PublicSignal returned by the default gnark mock.
// proofToReceipt (numIn=1, numOut=2) reads it as:
//
//	[0]=message [1]=merkleRoot [2]=nullifier [3]=treeNumber [4:6]=commitments [6]=revertCommitment
//
// so [2]=0xDEAD (57005) is the input nullifier and [6]=0xBEEF (48879) is the
// revert commitment — both asserted by the swap tests.
var defaultProofSignal = []*big.Int{
	big.NewInt(0xD0), big.NewInt(0xD1), big.NewInt(0xDEAD), big.NewInt(0xD3),
	big.NewInt(0xD4), big.NewInt(0xD5), big.NewInt(0xBEEF), big.NewInt(0xD7),
}

func defaultProof() dvp.Proof {
	return dvp.Proof{
		A:            [2]*big.Int{big.NewInt(1), big.NewInt(2)},
		B:            [2][2]*big.Int{{big.NewInt(3), big.NewInt(4)}, {big.NewInt(5), big.NewInt(6)}},
		C:            [2]*big.Int{big.NewInt(7), big.NewInt(8)},
		PublicSignal: defaultProofSignal,
	}
}

// newProofFixture wires happy-path defaults for every mock method (moq panics on
// an unset Func). Individual tests override a single Func to exercise error paths.
func newProofFixture() *proofFixture {
	return &proofFixture{
		config: service.ProofServiceConfig{
			ChainID:            big.NewInt(1),
			MerkleTreeDepth:    32,
			NumberOfJSParamsIn: 1,
		},
		merkle: &proofMerkleClientMock{
			GenerateMerkleProofFunc: func(context.Context, *big.Int, int, string) (*merkle.MerkleProof, error) {
				return &merkle.MerkleProof{
					Element:  big.NewInt(1),
					Elements: []*big.Int{big.NewInt(2), big.NewInt(3)},
					Indices:  big.NewInt(0),
					Root:     big.NewInt(4),
				}, nil
			},
		},
		kos: &proofKeysClientMock{
			GetPaymentSpendKeyFunc: func(context.Context, *keyspb.GetPaymentSpendKeyRequest, ...grpc.CallOption) (*keyspb.PaymentSpendKeyResponse, error) {
				return &keyspb.PaymentSpendKeyResponse{
					SecretKey: big.NewInt(0x5EC).Bytes(),
					PublicKey: big.NewInt(0x9AB).Bytes(),
				}, nil
			},
		},
		api: &proofAPIClientMock{
			CreateEnygmaJSProofFunc:        func(*dvp.EnygmaJoinSplitProofRequest) (dvp.Proof, error) { return defaultProof(), nil },
			CreateErc1155JSProofFunc:       func(*dvp.ERC1155JoinSplitProofRequest) (dvp.Proof, error) { return defaultProof(), nil },
			CreateErc721OwnershipProofFunc: func(*dvp.ERC721OwnershipProofRequest) (dvp.Proof, error) { return defaultProof(), nil },
		},
		// Distinct, non-0xBEEF commitments so change/self-destination deposits are
		// never mistaken for the revert deposit (whose commitment is 0xBEEF, from
		// the proof signal) in the swap tests' loops.
		calc: &proofCommitmentCalculatorMock{
			CalculateERC1155CommitmentFunc: func(*big.Int, *big.Int, string, string, *big.Int) (*big.Int, error) { return big.NewInt(0xC0DE), nil },
			CalculateNFTCommitmentFunc:     func(*big.Int, *big.Int, string, string) (*big.Int, error) { return big.NewInt(0xFACE), nil },
			CalculatePaymentCommitmentFunc: func(*big.Int, *big.Int, *big.Int, string) (*big.Int, error) { return big.NewInt(0xCAFE), nil },
			GetNFTUniqueIDFunc:             func(string, string) (*big.Int, error) { return big.NewInt(0xC1), nil },
		},
		repo: &proofDepositRepositoryMock{
			CreateDepositFunc:          func(context.Context, *types.DvpDeposit) error { return nil },
			BatchUpsertNullifiersFunc:  func(context.Context, map[string]string) error { return nil },
			UpsertDepositNullifierFunc: func(context.Context, *big.Int, *big.Int) error { return nil },
			GetDepositsByTokenFunc: func(context.Context, string, string, types.DvpTokenType, string, types.DvpDepositStatus) ([]types.DvpDeposit, error) {
				return nil, nil
			},
			GetFungibleDepositsFunc: func(context.Context, string, string, types.DvpTokenType, types.DvpDepositStatus) ([]types.DvpDeposit, error) {
				return nil, nil
			},
		},
		txMgr: &proofTxManagerMock{
			WithTransactionFunc: func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) },
		},
		viewPK: testdata.MlkemEncapsulationKey(),
	}
}

func (f *proofFixture) service() *service.ProofService {
	return service.NewProofService(
		f.config, f.merkle, f.kos, f.api, f.calc, f.repo, f.txMgr,
	)
}

// --- deposit helpers ---

func proofERC1155Deposit(amount int64) *types.DvpDeposit {
	return testdata.NewFungibleDeposit(amount, types.DvpERC1155,
		testdata.WithDepositTokenAddress("0x1234"),
		testdata.WithDepositTokenID("42"),
	)
}

func proofEnygmaDeposit(amount int64) *types.DvpDeposit {
	return testdata.NewFungibleDeposit(amount, types.DvpEnygma,
		testdata.WithDepositTokenAddress("0x5678"),
	)
}

func proofNFTDeposit(nftID string) *types.DvpDeposit {
	return testdata.NewNFTDeposit(nftID,
		testdata.WithDepositTokenAddress("0x1234"),
	)
}

func proofSampleSwap(tokenIn, tokenOut types.DvpTokenType) *types.DvpSwap {
	return testdata.NewDvpSwap(
		testdata.WithSharedID("shared-x"),
		testdata.WithLegs(
			tokenIn, "", "0xTokenIn", "", big.NewInt(100),
			tokenOut, "", "0xTokenOut", "tok-out", big.NewInt(200),
		),
	)
}

// --- constructor test ---

func TestNewProofService(t *testing.T) {
	t.Run("initializes all fields correctly", func(t *testing.T) {
		f := newProofFixture()
		require.NotNil(t, f.service())
		assert.Equal(t, big.NewInt(1), f.config.ChainID)
		assert.Equal(t, 32, f.config.MerkleTreeDepth)
		assert.Equal(t, 1, f.config.NumberOfJSParamsIn)
	})
}

// --- ERC1155 JS proof ---

func TestProofService_GenerateERC1155JSProof(t *testing.T) {
	t.Run("successfully generates proof with single deposit", func(t *testing.T) {
		f := newProofFixture()

		result, err := f.service().GenerateERC1155JSProof(
			context.Background(), f.viewPK,
			big.NewInt(777), big.NewInt(456), big.NewInt(0xA1), big.NewInt(50),
			"0xAlice", "0x1234", "42",
			[]*types.DvpDeposit{proofERC1155Deposit(100)},
		)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Proof)
	})

	t.Run("creates change deposit when change amount > 0", func(t *testing.T) {
		f := newProofFixture()

		result, err := f.service().GenerateERC1155JSProof(
			context.Background(), f.viewPK,
			big.NewInt(777), big.NewInt(456), big.NewInt(0xA1), big.NewInt(50),
			"0xAlice", "0x1234", "42",
			[]*types.DvpDeposit{proofERC1155Deposit(100)},
		)

		require.NoError(t, err)
		require.NotNil(t, result)
		// v2: 1 change deposit + 1 revert deposit
		assert.Len(t, f.repo.CreateDepositCalls(), 2)
	})

	t.Run("skips change deposit when change amount = 0", func(t *testing.T) {
		f := newProofFixture()

		result, err := f.service().GenerateERC1155JSProof(
			context.Background(), f.viewPK,
			big.NewInt(777), big.NewInt(456), big.NewInt(0xA1), big.NewInt(100),
			"0xAlice", "0x1234", "42",
			[]*types.DvpDeposit{proofERC1155Deposit(100)},
		)

		require.NoError(t, err)
		require.NotNil(t, result)
		// v2: revert deposit is always persisted even without change
		assert.Len(t, f.repo.CreateDepositCalls(), 1)
	})

	t.Run("returns error when merkleClient fails", func(t *testing.T) {
		f := newProofFixture()
		f.merkle.GenerateMerkleProofFunc = func(ctx context.Context, commitment *big.Int, treeNumber int, tokenAddress string) (*merkle.MerkleProof, error) {
			return nil, errors.New("merkle proof generation failed")
		}

		result, err := f.service().GenerateERC1155JSProof(
			context.Background(), f.viewPK,
			big.NewInt(777), big.NewInt(456), big.NewInt(0xA1), big.NewInt(50),
			"0xAlice", "0x1234", "42",
			[]*types.DvpDeposit{proofERC1155Deposit(100)},
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "merkle proof generation failed")
	})

	t.Run("returns error when kosClient fails", func(t *testing.T) {
		f := newProofFixture()
		f.kos.GetPaymentSpendKeyFunc = func(context.Context, *keyspb.GetPaymentSpendKeyRequest, ...grpc.CallOption) (*keyspb.PaymentSpendKeyResponse, error) {
			return nil, errors.New("keypair retrieval failed")
		}

		result, err := f.service().GenerateERC1155JSProof(
			context.Background(), f.viewPK,
			big.NewInt(777), big.NewInt(456), big.NewInt(0xA1), big.NewInt(50),
			"0xAlice", "0x1234", "42",
			[]*types.DvpDeposit{proofERC1155Deposit(100)},
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "keypair retrieval failed")
	})

	t.Run("returns error when proofClient fails", func(t *testing.T) {
		f := newProofFixture()
		f.api.CreateErc1155JSProofFunc = func(req *dvp.ERC1155JoinSplitProofRequest) (dvp.Proof, error) {
			return dvp.Proof{}, errors.New("proof generation failed")
		}

		result, err := f.service().GenerateERC1155JSProof(
			context.Background(), f.viewPK,
			big.NewInt(777), big.NewInt(456), big.NewInt(0xA1), big.NewInt(50),
			"0xAlice", "0x1234", "42",
			[]*types.DvpDeposit{proofERC1155Deposit(100)},
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "proof generation failed")
	})

	t.Run("returns error when depositRepository fails", func(t *testing.T) {
		f := newProofFixture()
		f.repo.CreateDepositFunc = func(ctx context.Context, d *types.DvpDeposit) error {
			return errors.New("deposit creation failed")
		}

		result, err := f.service().GenerateERC1155JSProof(
			context.Background(), f.viewPK,
			big.NewInt(777), big.NewInt(456), big.NewInt(0xA1), big.NewInt(50),
			"0xAlice", "0x1234", "42",
			[]*types.DvpDeposit{proofERC1155Deposit(100)},
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "deposit creation failed")
	})

	t.Run("returns error when tokenId parsing fails", func(t *testing.T) {
		f := newProofFixture()

		result, err := f.service().GenerateERC1155JSProof(
			context.Background(), f.viewPK,
			big.NewInt(777), big.NewInt(456), big.NewInt(0xA1), big.NewInt(50),
			"0xAlice", "0x1234", "not-a-number",
			[]*types.DvpDeposit{proofERC1155Deposit(100)},
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to parse erc1155 token id")
	})
}

// --- Enygma JS proof ---

func TestProofService_GenerateEnygmaJSProof(t *testing.T) {
	t.Run("successfully generates proof with single deposit", func(t *testing.T) {
		f := newProofFixture()

		result, err := f.service().GenerateEnygmaJSProof(
			context.Background(), f.viewPK,
			big.NewInt(777), big.NewInt(456), big.NewInt(0xA1), big.NewInt(50),
			"0x5678",
			[]*types.DvpDeposit{proofEnygmaDeposit(100)},
		)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Proof)
	})

	t.Run("creates change deposit when change amount > 0", func(t *testing.T) {
		f := newProofFixture()

		result, err := f.service().GenerateEnygmaJSProof(
			context.Background(), f.viewPK,
			big.NewInt(777), big.NewInt(456), big.NewInt(0xA1), big.NewInt(50),
			"0x5678",
			[]*types.DvpDeposit{proofEnygmaDeposit(100)},
		)

		require.NoError(t, err)
		require.NotNil(t, result)
		// v2: 1 change deposit + 1 revert deposit
		assert.Len(t, f.repo.CreateDepositCalls(), 2)
	})

	t.Run("returns error when merkleClient fails", func(t *testing.T) {
		f := newProofFixture()
		f.merkle.GenerateMerkleProofFunc = func(ctx context.Context, commitment *big.Int, treeNumber int, tokenAddress string) (*merkle.MerkleProof, error) {
			return nil, errors.New("merkle proof generation failed")
		}

		result, err := f.service().GenerateEnygmaJSProof(
			context.Background(), f.viewPK,
			big.NewInt(777), big.NewInt(456), big.NewInt(0xA1), big.NewInt(50),
			"0x5678",
			[]*types.DvpDeposit{proofEnygmaDeposit(100)},
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "merkle proof generation failed")
	})

	t.Run("returns error when kosClient fails", func(t *testing.T) {
		f := newProofFixture()
		f.kos.GetPaymentSpendKeyFunc = func(context.Context, *keyspb.GetPaymentSpendKeyRequest, ...grpc.CallOption) (*keyspb.PaymentSpendKeyResponse, error) {
			return nil, errors.New("keypair retrieval failed")
		}

		result, err := f.service().GenerateEnygmaJSProof(
			context.Background(), f.viewPK,
			big.NewInt(777), big.NewInt(456), big.NewInt(0xA1), big.NewInt(50),
			"0x5678",
			[]*types.DvpDeposit{proofEnygmaDeposit(100)},
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "keypair retrieval failed")
	})

	t.Run("returns error when proofClient fails", func(t *testing.T) {
		f := newProofFixture()
		f.api.CreateEnygmaJSProofFunc = func(req *dvp.EnygmaJoinSplitProofRequest) (dvp.Proof, error) {
			return dvp.Proof{}, errors.New("proof generation failed")
		}

		result, err := f.service().GenerateEnygmaJSProof(
			context.Background(), f.viewPK,
			big.NewInt(777), big.NewInt(456), big.NewInt(0xA1), big.NewInt(50),
			"0x5678",
			[]*types.DvpDeposit{proofEnygmaDeposit(100)},
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "proof generation failed")
	})

	t.Run("returns error when depositRepository fails", func(t *testing.T) {
		f := newProofFixture()
		f.repo.CreateDepositFunc = func(ctx context.Context, d *types.DvpDeposit) error {
			return errors.New("deposit creation failed")
		}

		result, err := f.service().GenerateEnygmaJSProof(
			context.Background(), f.viewPK,
			big.NewInt(777), big.NewInt(456), big.NewInt(0xA1), big.NewInt(50),
			"0x5678",
			[]*types.DvpDeposit{proofEnygmaDeposit(100)},
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "deposit creation failed")
	})
}

// --- Ownership proof ---

func TestProofService_GenerateOwnershipProof(t *testing.T) {
	t.Run("successfully generates ownership proof", func(t *testing.T) {
		f := newProofFixture()

		result, err := f.service().GenerateOwnershipProof(
			context.Background(), f.viewPK,
			big.NewInt(888), big.NewInt(456), big.NewInt(0xA1),
			proofNFTDeposit("42"),
		)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Proof)
	})

	t.Run("returns error when GetNFTUniqueID fails", func(t *testing.T) {
		f := newProofFixture()

		result, err := service.NewProofService(
			f.config, f.merkle, f.kos, f.api,
			&proofCalcWithNFTError{err: errors.New("nft uid calculation failed")},
			f.repo, f.txMgr,
		).GenerateOwnershipProof(
			context.Background(), f.viewPK,
			big.NewInt(888), big.NewInt(456), big.NewInt(0xA1),
			proofNFTDeposit("42"),
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "nft uid calculation failed")
	})

	t.Run("returns error when merkleClient fails", func(t *testing.T) {
		f := newProofFixture()
		f.merkle.GenerateMerkleProofFunc = func(ctx context.Context, commitment *big.Int, treeNumber int, tokenAddress string) (*merkle.MerkleProof, error) {
			return nil, errors.New("merkle proof failed")
		}

		result, err := f.service().GenerateOwnershipProof(
			context.Background(), f.viewPK,
			big.NewInt(888), big.NewInt(456), big.NewInt(0xA1),
			proofNFTDeposit("42"),
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "merkle proof failed")
	})

	t.Run("returns error when kosClient fails", func(t *testing.T) {
		f := newProofFixture()
		f.kos.GetPaymentSpendKeyFunc = func(context.Context, *keyspb.GetPaymentSpendKeyRequest, ...grpc.CallOption) (*keyspb.PaymentSpendKeyResponse, error) {
			return nil, errors.New("keypair retrieval failed")
		}

		result, err := f.service().GenerateOwnershipProof(
			context.Background(), f.viewPK,
			big.NewInt(888), big.NewInt(456), big.NewInt(0xA1),
			proofNFTDeposit("42"),
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "keypair retrieval failed")
	})

	t.Run("returns error when proofClient fails", func(t *testing.T) {
		f := newProofFixture()
		f.api.CreateErc721OwnershipProofFunc = func(req *dvp.ERC721OwnershipProofRequest) (dvp.Proof, error) {
			return dvp.Proof{}, errors.New("proof generation failed")
		}

		result, err := f.service().GenerateOwnershipProof(
			context.Background(), f.viewPK,
			big.NewInt(888), big.NewInt(456), big.NewInt(0xA1),
			proofNFTDeposit("42"),
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "proof generation failed")
	})
}

// --- Swap proofs ---

func TestProofService_GenerateEnygmaToERC721SwapProof(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("calls gnark with revertSalt and persists revert + self-destination + change deposits", func(t *testing.T) {
		f := newProofFixture()

		swap := proofSampleSwap(types.DvpEnygma, types.DvpERC721)
		deposits := []*types.DvpDeposit{proofEnygmaDeposit(150)}
		selfSalt := big.NewInt(0xAAAA)
		destSalt := big.NewInt(0xBBBB)
		destSpendPK := big.NewInt(0xCCCC)

		proof, err := f.service().GenerateEnygmaToERC721SwapProof(context.Background(), swap, deposits, f.viewPK, selfSalt, destSalt, destSpendPK)
		require.NoError(t, err)
		require.NotNil(t, proof)
		assert.Equal(t, big.NewInt(0xBEEF), proof.RevertCommitment)

		require.Len(t, f.api.CreateEnygmaJSProofCalls(), 1)
		req := f.api.CreateEnygmaJSProofCalls()[0].Req
		require.NotNil(t, req.RevertSalt)
		require.Len(t, req.SaltsIn, 1)
		assert.Equal(t, deposits[0].Salt, req.SaltsIn[0])
		require.Len(t, req.SaltsOut, 2)
		assert.Equal(t, destSalt, req.SaltsOut[0])

		require.GreaterOrEqual(t, len(f.repo.CreateDepositCalls()), 2)

		var sawRevertDeposit, sawSelfDestERC721 bool
		for _, c := range f.repo.CreateDepositCalls() {
			d := c.Deposit
			if d.Commitment.Cmp(big.NewInt(0xBEEF)) == 0 {
				sawRevertDeposit = true
				assert.Equal(t, types.DvpDepositPending, d.Status)
				assert.Equal(t, req.RevertSalt, d.Salt)
			}
			if d.TokenType == types.DvpERC721 && d.TokenID == "tok-out" {
				sawSelfDestERC721 = true
				assert.Equal(t, selfSalt, d.Salt)
			}
		}
		assert.True(t, sawRevertDeposit)
		assert.True(t, sawSelfDestERC721)

		require.Len(t, f.repo.BatchUpsertNullifiersCalls(), 1)
		require.Equal(t, "1500", firstKey(f.repo.BatchUpsertNullifiersCalls()[0].CommitmentNullifierMap))
	})

	t.Run("propagates gnark API error", func(t *testing.T) {
		f := newProofFixture()
		want := errors.New("gnark down")
		f.api.CreateEnygmaJSProofFunc = func(req *dvp.EnygmaJoinSplitProofRequest) (dvp.Proof, error) {
			return dvp.Proof{}, want
		}

		_, err := f.service().GenerateEnygmaToERC721SwapProof(
			context.Background(),
			proofSampleSwap(types.DvpEnygma, types.DvpERC721),
			[]*types.DvpDeposit{proofEnygmaDeposit(100)},
			f.viewPK,
			big.NewInt(1), big.NewInt(2), big.NewInt(3),
		)
		require.Error(t, err)
		require.ErrorIs(t, err, want)
	})
}

func TestProofService_GenerateEnygmaToERC1155SwapProof(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("persists self-destination ERC1155 deposit and revert deposit", func(t *testing.T) {
		f := newProofFixture()

		swap := proofSampleSwap(types.DvpEnygma, types.DvpERC1155)
		_, err := f.service().GenerateEnygmaToERC1155SwapProof(
			context.Background(),
			swap,
			[]*types.DvpDeposit{proofEnygmaDeposit(150)},
			f.viewPK,
			big.NewInt(0xAAAA), big.NewInt(0xBBBB), big.NewInt(0xCCCC),
		)
		require.NoError(t, err)

		var sawSelfDestERC1155, sawRevert bool
		for _, c := range f.repo.CreateDepositCalls() {
			d := c.Deposit
			if d.TokenType == types.DvpERC1155 {
				sawSelfDestERC1155 = true
			}
			if d.Commitment.Cmp(big.NewInt(0xBEEF)) == 0 {
				sawRevert = true
			}
		}
		assert.True(t, sawSelfDestERC1155)
		assert.True(t, sawRevert)
	})

	t.Run("returns error when GetPaymentSpendKey fails", func(t *testing.T) {
		f := newProofFixture()
		want := errors.New("kos down")
		f.kos.GetPaymentSpendKeyFunc = func(context.Context, *keyspb.GetPaymentSpendKeyRequest, ...grpc.CallOption) (*keyspb.PaymentSpendKeyResponse, error) {
			return nil, want
		}

		_, err := f.service().GenerateEnygmaToERC1155SwapProof(
			context.Background(),
			proofSampleSwap(types.DvpEnygma, types.DvpERC1155),
			[]*types.DvpDeposit{proofEnygmaDeposit(100)},
			f.viewPK,
			big.NewInt(1), big.NewInt(2), big.NewInt(3),
		)
		require.Error(t, err)
		require.ErrorIs(t, err, want)
		assert.Contains(t, err.Error(), "getting payment spend key")
	})

	t.Run("returns error when CalculateERC1155Commitment fails", func(t *testing.T) {
		f := newProofFixture()
		want := errors.New("calc error")
		f.calc.CalculateERC1155CommitmentFunc = func(spendPK, salt *big.Int, tokenAddress, tokenID string, tokenAmount *big.Int) (*big.Int, error) {
			return nil, want
		}

		_, err := f.service().GenerateEnygmaToERC1155SwapProof(
			context.Background(),
			proofSampleSwap(types.DvpEnygma, types.DvpERC1155),
			[]*types.DvpDeposit{proofEnygmaDeposit(100)},
			f.viewPK,
			big.NewInt(1), big.NewInt(2), big.NewInt(3),
		)
		require.Error(t, err)
		require.ErrorIs(t, err, want)
		assert.Contains(t, err.Error(), "calculating ERC1155 commitment")
	})

	t.Run("returns error when CreateEnygmaJSProof fails", func(t *testing.T) {
		f := newProofFixture()
		want := errors.New("gnark down")
		f.api.CreateEnygmaJSProofFunc = func(req *dvp.EnygmaJoinSplitProofRequest) (dvp.Proof, error) {
			return dvp.Proof{}, want
		}

		_, err := f.service().GenerateEnygmaToERC1155SwapProof(
			context.Background(),
			proofSampleSwap(types.DvpEnygma, types.DvpERC1155),
			[]*types.DvpDeposit{proofEnygmaDeposit(100)},
			f.viewPK,
			big.NewInt(1), big.NewInt(2), big.NewInt(3),
		)
		require.Error(t, err)
		require.ErrorIs(t, err, want)
	})

	t.Run("returns error when CreateDeposit fails for ERC1155", func(t *testing.T) {
		f := newProofFixture()
		want := errors.New("db error")
		callCount := 0
		f.repo.CreateDepositFunc = func(ctx context.Context, deposit *types.DvpDeposit) error {
			callCount++
			// The first CreateDeposit call is for the revert deposit inside GenerateEnygmaJSProof.
			// The second is the ERC1155 self-destination deposit we want to fail.
			if callCount >= 2 {
				return want
			}
			return nil
		}

		_, err := f.service().GenerateEnygmaToERC1155SwapProof(
			context.Background(),
			proofSampleSwap(types.DvpEnygma, types.DvpERC1155),
			[]*types.DvpDeposit{proofEnygmaDeposit(100)},
			f.viewPK,
			big.NewInt(1), big.NewInt(2), big.NewInt(3),
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "creating pending ERC1155 deposit")
	})

	t.Run("returns error when BatchUpsertNullifiers fails", func(t *testing.T) {
		f := newProofFixture()
		want := errors.New("nullifier error")
		f.repo.BatchUpsertNullifiersFunc = func(ctx context.Context, m map[string]string) error {
			return want
		}

		_, err := f.service().GenerateEnygmaToERC1155SwapProof(
			context.Background(),
			proofSampleSwap(types.DvpEnygma, types.DvpERC1155),
			[]*types.DvpDeposit{proofEnygmaDeposit(100)},
			f.viewPK,
			big.NewInt(1), big.NewInt(2), big.NewInt(3),
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "batch upserting nullifiers")
	})
}

func TestProofService_GenerateERC721ToEnygmaSwapProof(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("locks input + persists pending Enygma deposit atomically via txManager", func(t *testing.T) {
		f := newProofFixture()

		swap := proofSampleSwap(types.DvpERC721, types.DvpEnygma)
		input := proofNFTDeposit("tok-7")

		_, err := f.service().GenerateERC721ToEnygmaSwapProof(
			context.Background(), swap, input, f.viewPK,
			big.NewInt(0xAAAA), big.NewInt(0xBBBB), big.NewInt(0xCCCC),
		)
		require.NoError(t, err)

		assert.Len(t, f.txMgr.WithTransactionCalls(), 1)
		require.Len(t, f.repo.UpsertDepositNullifierCalls(), 1)
		assert.Equal(t, input.Commitment, f.repo.UpsertDepositNullifierCalls()[0].Commitment)
	})

	t.Run("does not create deposits when txManager rolls back", func(t *testing.T) {
		f := newProofFixture()
		want := errors.New("rollback triggered")
		f.txMgr.WithTransactionFunc = func(context.Context, func(context.Context) error) error {
			return want
		}

		_, err := f.service().GenerateERC721ToEnygmaSwapProof(
			context.Background(),
			proofSampleSwap(types.DvpERC721, types.DvpEnygma),
			proofNFTDeposit("tok-7"),
			f.viewPK,
			big.NewInt(1), big.NewInt(2), big.NewInt(3),
		)
		require.ErrorIs(t, err, want)
		assert.Empty(t, f.repo.UpsertDepositNullifierCalls())
	})
}

func TestProofService_GenerateERC1155ToEnygmaSwapProof(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("batch-upserts nullifiers for spent ERC1155 inputs atomically", func(t *testing.T) {
		f := newProofFixture()

		swap := proofSampleSwap(types.DvpERC1155, types.DvpEnygma)
		swap.TokenInID = "3"
		input := proofERC1155Deposit(100)
		input.Commitment = big.NewInt(100)

		_, err := f.service().GenerateERC1155ToEnygmaSwapProof(
			context.Background(), swap,
			[]*types.DvpDeposit{input},
			f.viewPK,
			big.NewInt(0xAAAA), big.NewInt(0xBBBB), big.NewInt(0xCCCC),
		)
		require.NoError(t, err)

		assert.Len(t, f.txMgr.WithTransactionCalls(), 1)
		require.GreaterOrEqual(t, len(f.repo.BatchUpsertNullifiersCalls()), 1)

		var sawInput bool
		for _, c := range f.repo.BatchUpsertNullifiersCalls() {
			mp := c.CommitmentNullifierMap
			if _, ok := mp["100"]; ok {
				sawInput = true
				assert.Equal(t, "57005", mp["100"], "nullifier from defaultProofSignal[2] = 0xDEAD = 57005")
				break
			}
		}
		assert.True(t, sawInput)
	})

	t.Run("returns error when GetPaymentSpendKey fails", func(t *testing.T) {
		f := newProofFixture()
		want := errors.New("kos down")
		f.kos.GetPaymentSpendKeyFunc = func(context.Context, *keyspb.GetPaymentSpendKeyRequest, ...grpc.CallOption) (*keyspb.PaymentSpendKeyResponse, error) {
			return nil, want
		}

		swap := proofSampleSwap(types.DvpERC1155, types.DvpEnygma)
		swap.TokenInID = "3"
		_, err := f.service().GenerateERC1155ToEnygmaSwapProof(
			context.Background(), swap,
			[]*types.DvpDeposit{proofERC1155Deposit(100)},
			f.viewPK,
			big.NewInt(1), big.NewInt(2), big.NewInt(3),
		)
		require.Error(t, err)
		require.ErrorIs(t, err, want)
		assert.Contains(t, err.Error(), "getting payment spend key")
	})

	t.Run("returns error when CalculatePaymentCommitment fails", func(t *testing.T) {
		f := newProofFixture()
		want := errors.New("calc error")
		f.calc.CalculatePaymentCommitmentFunc = func(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error) {
			return nil, want
		}

		swap := proofSampleSwap(types.DvpERC1155, types.DvpEnygma)
		swap.TokenInID = "3"
		_, err := f.service().GenerateERC1155ToEnygmaSwapProof(
			context.Background(), swap,
			[]*types.DvpDeposit{proofERC1155Deposit(100)},
			f.viewPK,
			big.NewInt(1), big.NewInt(2), big.NewInt(3),
		)
		require.Error(t, err)
		require.ErrorIs(t, err, want)
		assert.Contains(t, err.Error(), "calculating payment commitment")
	})

	t.Run("returns error when GenerateERC1155JSProof fails", func(t *testing.T) {
		f := newProofFixture()
		want := errors.New("gnark down")
		f.api.CreateErc1155JSProofFunc = func(req *dvp.ERC1155JoinSplitProofRequest) (dvp.Proof, error) {
			return dvp.Proof{}, want
		}

		swap := proofSampleSwap(types.DvpERC1155, types.DvpEnygma)
		swap.TokenInID = "3"
		_, err := f.service().GenerateERC1155ToEnygmaSwapProof(
			context.Background(), swap,
			[]*types.DvpDeposit{proofERC1155Deposit(100)},
			f.viewPK,
			big.NewInt(1), big.NewInt(2), big.NewInt(3),
		)
		require.Error(t, err)
		require.ErrorIs(t, err, want)
	})

	t.Run("does not create deposits when txManager rolls back", func(t *testing.T) {
		f := newProofFixture()
		want := errors.New("rollback triggered")
		f.txMgr.WithTransactionFunc = func(context.Context, func(context.Context) error) error {
			return want
		}

		swap := proofSampleSwap(types.DvpERC1155, types.DvpEnygma)
		swap.TokenInID = "3"
		_, err := f.service().GenerateERC1155ToEnygmaSwapProof(
			context.Background(), swap,
			[]*types.DvpDeposit{proofERC1155Deposit(100)},
			f.viewPK,
			big.NewInt(1), big.NewInt(2), big.NewInt(3),
		)
		require.ErrorIs(t, err, want)
		assert.Empty(t, f.repo.BatchUpsertNullifiersCalls())
	})

	t.Run("returns error when BatchUpsertNullifiers fails inside tx", func(t *testing.T) {
		f := newProofFixture()
		want := errors.New("nullifier error")
		f.repo.BatchUpsertNullifiersFunc = func(ctx context.Context, m map[string]string) error {
			return want
		}

		swap := proofSampleSwap(types.DvpERC1155, types.DvpEnygma)
		swap.TokenInID = "3"
		_, err := f.service().GenerateERC1155ToEnygmaSwapProof(
			context.Background(), swap,
			[]*types.DvpDeposit{proofERC1155Deposit(100)},
			f.viewPK,
			big.NewInt(1), big.NewInt(2), big.NewInt(3),
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "locking deposits and creating pending enygma deposit")
	})
}

// --- Withdraw proofs ---

func TestProofService_GenerateERC721WithdrawProof(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("delegates to GenerateOwnershipProof with operator pubkey and zero payment commitment", func(t *testing.T) {
		f := newProofFixture()

		operatorPK := big.NewInt(0xDEADBEEF)
		_, err := f.service().GenerateERC721WithdrawProof(context.Background(), f.viewPK, big.NewInt(0xBBBB), operatorPK, proofNFTDeposit("tok-with-7"))
		require.NoError(t, err)

		require.Len(t, f.api.CreateErc721OwnershipProofCalls(), 1)
		req := f.api.CreateErc721OwnershipProofCalls()[0].Req
		assert.Equal(t, big.NewInt(0), req.PaymentCommitment)
		assert.Equal(t, operatorPK, req.PubKeyOut)
		assert.Equal(t, big.NewInt(0xBBBB), req.SaltOut)
	})
}

func TestProofService_GenerateERC1155WithdrawProof(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("delegates to GenerateERC1155JSProof with operator pubkey and zero payment commitment", func(t *testing.T) {
		f := newProofFixture()

		operatorPK := big.NewInt(0xDEADBEEF)
		_, err := f.service().GenerateERC1155WithdrawProof(
			context.Background(), f.viewPK,
			big.NewInt(0xBBBB), operatorPK,
			"0xUser", "0xToken", "42", big.NewInt(50),
			[]*types.DvpDeposit{proofERC1155Deposit(50)},
		)
		require.NoError(t, err)

		require.Len(t, f.api.CreateErc1155JSProofCalls(), 1)
		req := f.api.CreateErc1155JSProofCalls()[0].Req
		assert.Equal(t, big.NewInt(0), req.NftCommitment)
		assert.Len(t, req.PubKeysOut, 2)
		assert.Equal(t, operatorPK, req.PubKeysOut[0])
	})
}

// --- helpers ---

// proofCalcWithNFTError implements proofCommitmentCalculator but returns an
// error from GetNFTUniqueID. The zero-value proofCommitmentCalculatorMock always succeeds.
type proofCalcWithNFTError struct {
	err error
}

func (m *proofCalcWithNFTError) GetNFTUniqueID(_, _ string) (*big.Int, error) {
	return nil, m.err
}
func (m *proofCalcWithNFTError) CalculateNFTCommitment(_, _ *big.Int, _, _ string) (*big.Int, error) {
	return big.NewInt(0xC2), nil
}
func (m *proofCalcWithNFTError) CalculateERC1155Commitment(_, _ *big.Int, _, _ string, _ *big.Int) (*big.Int, error) {
	return big.NewInt(0xC3), nil
}
func (m *proofCalcWithNFTError) CalculatePaymentCommitment(_, _, _ *big.Int, _ string) (*big.Int, error) {
	return big.NewInt(0xC4), nil
}

// firstKey returns the first key in a map (iteration order is non-deterministic
// but we use it only when we expect a single-entry map).
func firstKey(m map[string]string) string {
	for k := range m {
		return k
	}
	return ""
}
