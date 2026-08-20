package service_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	keyspb "github.com/raylsnetwork/rayls-sovereign-relayer/cts/gen/proto/keys"
	"github.com/raylsnetwork/rayls-sovereign-relayer/dvp"
	"github.com/raylsnetwork/rayls-sovereign-relayer/dvp/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/dvp/testdata"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// consolidationFixture bundles the mocks and config needed by every
// consolidation test. Defaults are the happy-path ERC1155 scenario;
// individual tests override only the fields they care about.
type consolidationFixture struct {
	config                  service.ConsolidationConfig
	repo                    *consolidationDepositRepositoryMock
	spendKey                *consolidationSpendKeyProviderMock
	calc                    *consolidationCommitmentCalculatorMock
	proofService            *consolidationProofServiceMock
	dvpClient               *consolidationDvpClientMock
	enygmaClient            *ConsolidationEnygmaClientMock
	enygmaIntegrationClient *ConsolidationDvpIntegrationClientMock
	waiter                  *consolidationDepositWaiterMock
	txManager               *consolidationTxManagerMock
	viewPK                  []byte
}

func newConsolidationFixture() *consolidationFixture {
	return &consolidationFixture{
		config: service.ConsolidationConfig{
			ChainID:               big.NewInt(1),
			MaxNumberOfJSDeposits: 5,
		},
		repo: &consolidationDepositRepositoryMock{
			BatchUpdateStatusForCommitmentsFunc: func(ctx context.Context, commitments []string, status types.DvpDepositStatus) error {
				return nil
			},
			BatchUpsertNullifiersFunc: func(ctx context.Context, commitmentNullifierMap map[string]string) error {
				return nil
			},
			GetDepositByCommitmentFunc: func(ctx context.Context, commitment *big.Int) (*types.DvpDeposit, error) {
				return &types.DvpDeposit{Commitment: commitment, Status: types.DvpDepositUnspent}, nil
			},
			CreateDepositFunc: func(ctx context.Context, deposit *types.DvpDeposit) error {
				return nil
			},
		},
		spendKey: &consolidationSpendKeyProviderMock{
			GetPaymentSpendKeyFunc: func(ctx context.Context, in *keyspb.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keyspb.PaymentSpendKeyResponse, error) {
				return &keyspb.PaymentSpendKeyResponse{
					SecretKey: big.NewInt(998).Bytes(),
					PublicKey: big.NewInt(999).Bytes(),
				}, nil
			},
		},
		calc: &consolidationCommitmentCalculatorMock{
			CalculateERC1155CommitmentFunc: func(spendPK, salt *big.Int, tokenAddress string, tokenID string, tokenAmount *big.Int) (*big.Int, error) {
				return big.NewInt(888), nil
			},
			CalculatePaymentCommitmentFunc: func(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error) {
				return big.NewInt(888), nil
			},
		},
		proofService: &consolidationProofServiceMock{
			GenerateERC1155JSProofFunc: func(ctx context.Context, sourceViewPublicKey []byte, paymentCommitment *big.Int, destinationPaymentPublicKey *big.Int, destinationSalt *big.Int, paymentAmount *big.Int, userAddress string, tokenAddress string, tokenID string, deposits []*types.DvpDeposit) (*dvp.ProofReceipt, error) {
				return testdata.NewProofReceipt(
					testdata.WithNullifiers(big.NewInt(1), big.NewInt(2)),
					testdata.WithCommitments(big.NewInt(888)),
				), nil
			},
			GenerateEnygmaJSProofFunc: func(ctx context.Context, sourceViewPublicKey []byte, nftCommitment *big.Int, destinationPaymentPublicKey *big.Int, destinationSalt *big.Int, paymentAmount *big.Int, tokenAddress string, deposits []*types.DvpDeposit) (*dvp.ProofReceipt, error) {
				return testdata.NewProofReceipt(
					testdata.WithNullifiers(big.NewInt(1)),
					testdata.WithCommitments(big.NewInt(888)),
				), nil
			},
		},
		dvpClient: &consolidationDvpClientMock{
			MixFundsERC1155Func: func(ctx context.Context, _ string, contractAddress common.Address, proofReceipt *dvp.ProofReceipt) error {
				return nil
			},
		},
		enygmaClient:            &ConsolidationEnygmaClientMock{},
		enygmaIntegrationClient: &ConsolidationDvpIntegrationClientMock{},
		waiter: &consolidationDepositWaiterMock{
			WaitUntilDepositIsConfirmedFunc: func(ctx context.Context, deposit *types.DvpDeposit) error {
				return nil
			},
		},
		txManager: &consolidationTxManagerMock{
			WithTransactionFunc: func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			},
		},
		viewPK: testdata.MlkemEncapsulationKey(),
	}
}

func (f *consolidationFixture) service() *service.ConsolidationService {
	return service.NewConsolidationService(
		f.config, f.repo, f.spendKey, f.calc, f.proofService,
		f.dvpClient, f.enygmaClient, f.enygmaIntegrationClient, f.waiter, f.txManager,
	)
}

// withEnygmaHappyPath wires up the enygma + integration client mocks so
// ConsolidateEnygmaDeposits can run end-to-end.
func (f *consolidationFixture) withEnygmaHappyPath() *consolidationFixture {
	f.enygmaClient.GetDvpIntegrationContractAddressFunc = func(_ context.Context, tokenAddress common.Address) (common.Address, error) {
		return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
	}
	f.enygmaIntegrationClient.ConsolidateFundsFunc = func(ctx context.Context, _ string, enygmaIntegrationAddress common.Address, proofReceipt *dvp.ProofReceipt) error {
		return nil
	}
	return f
}

// --- helpers ---

func erc1155Deposit(amount int64) *types.DvpDeposit {
	return testdata.NewFungibleDeposit(amount, types.DvpERC1155,
		testdata.WithDepositTokenAddress("0x1234"),
		testdata.WithDepositTokenID("42"),
	)
}

func enygmaDeposit(amount int64) *types.DvpDeposit {
	return testdata.NewFungibleDeposit(amount, types.DvpEnygma,
		testdata.WithDepositTokenAddress("0x5678"),
	)
}

// --- tests ---

func TestNewConsolidationService(t *testing.T) {
	t.Run("initializes with all dependencies", func(t *testing.T) {
		f := newConsolidationFixture()
		require.NotNil(t, f.service())
	})
}

func TestConsolidateERC1155Deposits(t *testing.T) {
	t.Run("successfully consolidates multiple deposits with no change", func(t *testing.T) {
		f := newConsolidationFixture()
		deposits := []*types.DvpDeposit{erc1155Deposit(50), erc1155Deposit(50)}

		result, err := f.service().ConsolidateERC1155Deposits(context.Background(), "test-event-id", f.viewPK, deposits, big.NewInt(100))

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result, 1)

		statusCalls := f.repo.BatchUpdateStatusForCommitmentsCalls()
		assert.Len(t, statusCalls, 1)
		assert.Equal(t, types.DvpDepositLocked, statusCalls[0].Status)
	})

	t.Run("returns error when deposits list is empty", func(t *testing.T) {
		f := newConsolidationFixture()

		result, err := f.service().ConsolidateERC1155Deposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{}, big.NewInt(100))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no deposits")
	})

	t.Run("returns error when deposits are from different users", func(t *testing.T) {
		f := newConsolidationFixture()
		d1 := erc1155Deposit(50)
		d2 := erc1155Deposit(50)
		d2.UserAddress = "0xBob"

		result, err := f.service().ConsolidateERC1155Deposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{d1, d2}, big.NewInt(100))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not for the same user")
	})

	t.Run("returns error when deposits are for different tokens", func(t *testing.T) {
		f := newConsolidationFixture()
		d1 := erc1155Deposit(50)
		d2 := erc1155Deposit(50)
		d2.TokenAddress = "0x5678"

		result, err := f.service().ConsolidateERC1155Deposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{d1, d2}, big.NewInt(100))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not for the same token")
	})

	t.Run("returns error when deposits are for different token IDs", func(t *testing.T) {
		f := newConsolidationFixture()
		d1 := erc1155Deposit(50)
		d2 := erc1155Deposit(50)
		d2.TokenID = "43"

		result, err := f.service().ConsolidateERC1155Deposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{d1, d2}, big.NewInt(100))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not for the same token")
	})

	t.Run("returns error when consolidation amount exceeds total amount", func(t *testing.T) {
		f := newConsolidationFixture()
		deposits := []*types.DvpDeposit{erc1155Deposit(50), erc1155Deposit(50)}

		result, err := f.service().ConsolidateERC1155Deposits(context.Background(), "test-event-id", f.viewPK, deposits, big.NewInt(200))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "consolidation amount is greater")
	})

	t.Run("returns error when KOSClient fails", func(t *testing.T) {
		f := newConsolidationFixture()
		f.spendKey.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keyspb.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keyspb.PaymentSpendKeyResponse, error) {
			return nil, errors.New("keypair creation failed")
		}

		result, err := f.service().ConsolidateERC1155Deposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{erc1155Deposit(50)}, big.NewInt(50))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "keypair creation failed")
	})

	t.Run("returns error when proof generation fails", func(t *testing.T) {
		f := newConsolidationFixture()
		f.proofService.GenerateERC1155JSProofFunc = func(ctx context.Context, sourceViewPublicKey []byte, paymentCommitment *big.Int, destinationPaymentPublicKey *big.Int, destinationSalt *big.Int, paymentAmount *big.Int, userAddress string, tokenAddress string, tokenID string, deposits []*types.DvpDeposit) (*dvp.ProofReceipt, error) {
			return nil, errors.New("proof generation failed")
		}

		result, err := f.service().ConsolidateERC1155Deposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{erc1155Deposit(50)}, big.NewInt(50))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "proof generation failed")

		statusCalls := f.repo.BatchUpdateStatusForCommitmentsCalls()
		require.Len(t, statusCalls, 0)
	})

	t.Run("returns error when MixFundsERC1155 fails", func(t *testing.T) {
		f := newConsolidationFixture()
		f.dvpClient.MixFundsERC1155Func = func(ctx context.Context, _ string, contractAddress common.Address, proofReceipt *dvp.ProofReceipt) error {
			return errors.New("mixing funds failed")
		}

		result, err := f.service().ConsolidateERC1155Deposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{erc1155Deposit(50)}, big.NewInt(50))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to consolidate ERC1155 funds: mixing funds failed")

		statusCalls := f.repo.BatchUpdateStatusForCommitmentsCalls()
		require.Len(t, statusCalls, 2)
		assert.Equal(t, types.DvpDepositLocked, statusCalls[0].Status)
		assert.Equal(t, types.DvpDepositUnspent, statusCalls[1].Status)
	})

	t.Run("passes nil commitment when calculator returns nil", func(t *testing.T) {
		f := newConsolidationFixture()
		f.calc.CalculateERC1155CommitmentFunc = func(spendPK, salt *big.Int, tokenAddress string, tokenID string, tokenAmount *big.Int) (*big.Int, error) {
			return nil, nil
		}
		f.proofService.GenerateERC1155JSProofFunc = func(ctx context.Context, sourceViewPublicKey []byte, paymentCommitment *big.Int, destinationPaymentPublicKey *big.Int, destinationSalt *big.Int, paymentAmount *big.Int, userAddress string, tokenAddress string, tokenID string, deposits []*types.DvpDeposit) (*dvp.ProofReceipt, error) {
			return nil, errors.New("proof generation failed with nil commitment")
		}

		result, err := f.service().ConsolidateERC1155Deposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{erc1155Deposit(50)}, big.NewInt(50))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "proof generation failed with nil commitment")
	})

	t.Run("returns error when BatchUpsertNullifiers fails", func(t *testing.T) {
		f := newConsolidationFixture()
		f.repo.BatchUpsertNullifiersFunc = func(ctx context.Context, commitmentNullifierMap map[string]string) error {
			return errors.New("nullifier upsert failed")
		}

		result, err := f.service().ConsolidateERC1155Deposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{erc1155Deposit(50)}, big.NewInt(50))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "nullifier upsert failed")
	})

	t.Run("returns error when CreateDeposit fails", func(t *testing.T) {
		f := newConsolidationFixture()
		f.repo.CreateDepositFunc = func(ctx context.Context, deposit *types.DvpDeposit) error {
			return errors.New("deposit creation failed")
		}

		result, err := f.service().ConsolidateERC1155Deposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{erc1155Deposit(50)}, big.NewInt(50))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "deposit creation failed")
	})

	t.Run("returns error when depositWaiter fails", func(t *testing.T) {
		f := newConsolidationFixture()
		f.waiter.WaitUntilDepositIsConfirmedFunc = func(ctx context.Context, deposit *types.DvpDeposit) error {
			return errors.New("deposit confirmation timeout")
		}
		deposits := []*types.DvpDeposit{erc1155Deposit(50), erc1155Deposit(50)}

		result, err := f.service().ConsolidateERC1155Deposits(context.Background(), "test-event-id", f.viewPK, deposits, big.NewInt(100))

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "deposit confirmation timeout")
	})
}

func TestConsolidateEnygmaDeposits(t *testing.T) {
	t.Run("returns error when deposits list is empty", func(t *testing.T) {
		f := newConsolidationFixture()

		result, err := f.service().ConsolidateEnygmaDeposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{}, big.NewInt(100))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no deposits")
	})

	t.Run("returns error when consolidation amount exceeds total amount", func(t *testing.T) {
		f := newConsolidationFixture()

		result, err := f.service().ConsolidateEnygmaDeposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{enygmaDeposit(50)}, big.NewInt(100))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "consolidation amount is greater")
	})

	t.Run("returns error when proof generation fails", func(t *testing.T) {
		f := newConsolidationFixture()
		f.proofService.GenerateEnygmaJSProofFunc = func(ctx context.Context, sourceViewPublicKey []byte, nftCommitment *big.Int, destinationPaymentPublicKey *big.Int, destinationSalt *big.Int, paymentAmount *big.Int, tokenAddress string, deposits []*types.DvpDeposit) (*dvp.ProofReceipt, error) {
			return nil, errors.New("proof generation failed")
		}

		result, err := f.service().ConsolidateEnygmaDeposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{enygmaDeposit(50)}, big.NewInt(50))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "proof generation failed")

		statusCalls := f.repo.BatchUpdateStatusForCommitmentsCalls()
		require.Len(t, statusCalls, 0)
	})

	t.Run("returns error when BatchUpsertNullifiers fails", func(t *testing.T) {
		f := newConsolidationFixture().withEnygmaHappyPath()
		f.repo.BatchUpsertNullifiersFunc = func(ctx context.Context, commitmentNullifierMap map[string]string) error {
			return errors.New("nullifier upsert failed")
		}

		result, err := f.service().ConsolidateEnygmaDeposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{enygmaDeposit(50)}, big.NewInt(50))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "nullifier upsert failed")
	})

	t.Run("stores nil commitment when calculator returns nil", func(t *testing.T) {
		f := newConsolidationFixture().withEnygmaHappyPath()
		var createdDeposit *types.DvpDeposit
		f.calc.CalculatePaymentCommitmentFunc = func(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error) {
			return nil, nil
		}
		f.repo.CreateDepositFunc = func(ctx context.Context, deposit *types.DvpDeposit) error {
			createdDeposit = deposit
			return nil
		}
		f.repo.GetDepositByCommitmentFunc = func(ctx context.Context, commitment *big.Int) (*types.DvpDeposit, error) {
			return createdDeposit, nil
		}

		result, err := f.service().ConsolidateEnygmaDeposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{enygmaDeposit(50)}, big.NewInt(50))

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Nil(t, createdDeposit.Commitment)
	})

	t.Run("returns error when CreateDeposit fails", func(t *testing.T) {
		f := newConsolidationFixture().withEnygmaHappyPath()
		f.repo.CreateDepositFunc = func(ctx context.Context, deposit *types.DvpDeposit) error {
			return errors.New("deposit creation failed")
		}

		result, err := f.service().ConsolidateEnygmaDeposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{enygmaDeposit(50)}, big.NewInt(50))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "deposit creation failed")
	})

	t.Run("returns error when KOSClient fails to create key pair", func(t *testing.T) {
		f := newConsolidationFixture()
		f.spendKey.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keyspb.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keyspb.PaymentSpendKeyResponse, error) {
			return nil, errors.New("key pair creation failed")
		}

		result, err := f.service().ConsolidateEnygmaDeposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{enygmaDeposit(50)}, big.NewInt(50))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "key pair creation failed")
	})

	t.Run("returns error when GetDvpIntegrationContractAddress fails", func(t *testing.T) {
		f := newConsolidationFixture()
		f.enygmaClient.GetDvpIntegrationContractAddressFunc = func(_ context.Context, tokenAddress common.Address) (common.Address, error) {
			return common.Address{}, errors.New("get integration address failed")
		}

		result, err := f.service().ConsolidateEnygmaDeposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{enygmaDeposit(50)}, big.NewInt(50))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "get integration address failed")
	})

	t.Run("returns error when ConsolidateFunds fails", func(t *testing.T) {
		f := newConsolidationFixture().withEnygmaHappyPath()
		f.enygmaIntegrationClient.ConsolidateFundsFunc = func(ctx context.Context, _ string, enygmaIntegrationAddress common.Address, proofReceipt *dvp.ProofReceipt) error {
			return errors.New("consolidate funds failed")
		}

		result, err := f.service().ConsolidateEnygmaDeposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{enygmaDeposit(50)}, big.NewInt(50))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to consolidate Enygma funds")

		statusCalls := f.repo.BatchUpdateStatusForCommitmentsCalls()
		require.Len(t, statusCalls, 2)
		assert.Equal(t, types.DvpDepositLocked, statusCalls[0].Status)
		assert.Equal(t, types.DvpDepositUnspent, statusCalls[1].Status)
	})

	t.Run("returns error when depositWaiter fails", func(t *testing.T) {
		f := newConsolidationFixture().withEnygmaHappyPath()
		f.waiter.WaitUntilDepositIsConfirmedFunc = func(ctx context.Context, deposit *types.DvpDeposit) error {
			return errors.New("deposit waiter failed")
		}

		result, err := f.service().ConsolidateEnygmaDeposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{enygmaDeposit(50)}, big.NewInt(50))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "deposit waiter failed")
	})

	t.Run("successfully consolidates enygma deposits", func(t *testing.T) {
		f := newConsolidationFixture().withEnygmaHappyPath()
		consolidated := enygmaDeposit(50)
		consolidated.Commitment = big.NewInt(888)
		f.repo.GetDepositByCommitmentFunc = func(ctx context.Context, commitment *big.Int) (*types.DvpDeposit, error) {
			return consolidated, nil
		}

		result, err := f.service().ConsolidateEnygmaDeposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{enygmaDeposit(50)}, big.NewInt(50))

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Len(t, result, 1)
		assert.Equal(t, consolidated, result[0])

		createDepositCalls := f.repo.CreateDepositCalls()
		require.Len(t, createDepositCalls, 1)
		assert.Equal(t, types.DvpEnygma, createDepositCalls[0].Deposit.TokenType)
		assert.Equal(t, types.DvpDepositPending, createDepositCalls[0].Deposit.Status)
	})
}

func TestPrepareDepositsForJSProof(t *testing.T) {
	t.Run("returns deposits as-is when count <= MaxNumberOfJSDeposits", func(t *testing.T) {
		f := newConsolidationFixture()
		deposits := []*types.DvpDeposit{erc1155Deposit(50), erc1155Deposit(50)}

		result, err := f.service().PrepareDepositsForJSProof(context.Background(), "test-event-id", f.viewPK, deposits)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, deposits, result)
	})

	t.Run("returns deposits as-is when empty list", func(t *testing.T) {
		f := newConsolidationFixture()

		result, err := f.service().PrepareDepositsForJSProof(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result, 0)
	})

	t.Run("consolidates deposits when count exceeds MaxNumberOfJSDeposits for ERC1155", func(t *testing.T) {
		f := newConsolidationFixture()
		f.config.MaxNumberOfJSDeposits = 2
		consolidated := erc1155Deposit(100)

		f.repo.GetDepositByCommitmentFunc = func(ctx context.Context, commitment *big.Int) (*types.DvpDeposit, error) {
			return consolidated, nil
		}

		deposits := []*types.DvpDeposit{erc1155Deposit(50), erc1155Deposit(50), erc1155Deposit(50)}

		result, err := f.service().PrepareDepositsForJSProof(context.Background(), "test-event-id", f.viewPK, deposits)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result, 2) // consolidated + remaining
	})

	t.Run("returns error when consolidating fails", func(t *testing.T) {
		f := newConsolidationFixture()
		f.config.MaxNumberOfJSDeposits = 2
		f.spendKey.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keyspb.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keyspb.PaymentSpendKeyResponse, error) {
			return nil, errors.New("key creation failed")
		}

		deposits := []*types.DvpDeposit{erc1155Deposit(50), erc1155Deposit(50), erc1155Deposit(50)}

		result, err := f.service().PrepareDepositsForJSProof(context.Background(), "test-event-id", f.viewPK, deposits)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "key creation failed")
	})

	t.Run("returns error when unsupported token type", func(t *testing.T) {
		f := newConsolidationFixture()
		f.config.MaxNumberOfJSDeposits = 2
		unsupported := testdata.NewDvpDeposit(
			testdata.WithDepositTokenType(types.DvpTokenType(999)),
			testdata.WithDepositTokenAmount(big.NewInt(50)),
		)

		deposits := []*types.DvpDeposit{unsupported, unsupported, unsupported}

		result, err := f.service().PrepareDepositsForJSProof(context.Background(), "test-event-id", f.viewPK, deposits)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "unsupported token")
	})
}

func TestGetConsolidatedDeposits(t *testing.T) {
	t.Run("returns empty list when no commitments provided", func(t *testing.T) {
		f := newConsolidationFixture()
		f.proofService.GenerateERC1155JSProofFunc = func(ctx context.Context, sourceViewPublicKey []byte, paymentCommitment *big.Int, destinationPaymentPublicKey *big.Int, destinationSalt *big.Int, paymentAmount *big.Int, userAddress string, tokenAddress string, tokenID string, deposits []*types.DvpDeposit) (*dvp.ProofReceipt, error) {
			return testdata.NewProofReceipt(
				testdata.WithNullifiers(big.NewInt(1)),
				testdata.WithCommitments(), // empty
			), nil
		}

		result, err := f.service().ConsolidateERC1155Deposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{erc1155Deposit(50)}, big.NewInt(50))

		require.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("skips nil deposits (dummy/change=0)", func(t *testing.T) {
		f := newConsolidationFixture()
		d2 := erc1155Deposit(75)
		callCount := 0
		f.repo.GetDepositByCommitmentFunc = func(ctx context.Context, commitment *big.Int) (*types.DvpDeposit, error) {
			callCount++
			if callCount == 1 {
				return nil, nil //nolint:nilnil // intentional nil return in test mock
			}
			return d2, nil
		}
		f.proofService.GenerateERC1155JSProofFunc = func(ctx context.Context, sourceViewPublicKey []byte, paymentCommitment *big.Int, destinationPaymentPublicKey *big.Int, destinationSalt *big.Int, paymentAmount *big.Int, userAddress string, tokenAddress string, tokenID string, deposits []*types.DvpDeposit) (*dvp.ProofReceipt, error) {
			return testdata.NewProofReceipt(
				testdata.WithNullifiers(big.NewInt(1)),
				testdata.WithCommitments(big.NewInt(777), big.NewInt(888)),
			), nil
		}

		result, err := f.service().ConsolidateERC1155Deposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{erc1155Deposit(50)}, big.NewInt(50))

		require.NoError(t, err)
		assert.Len(t, result, 1)
	})

	t.Run("returns error when GetDepositByCommitment fails", func(t *testing.T) {
		f := newConsolidationFixture()
		f.repo.GetDepositByCommitmentFunc = func(ctx context.Context, commitment *big.Int) (*types.DvpDeposit, error) {
			return nil, errors.New("database error")
		}

		result, err := f.service().ConsolidateERC1155Deposits(context.Background(), "test-event-id", f.viewPK, []*types.DvpDeposit{erc1155Deposit(50)}, big.NewInt(50))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to fetch deposit")
	})
}

func TestCalculateTotalAmountOfDeposits(t *testing.T) {
	t.Run("returns zero for empty list", func(t *testing.T) {
		total := service.CalculateTotalAmountOfDeposits([]*types.DvpDeposit{})
		assert.Equal(t, big.NewInt(0), total)
	})

	t.Run("returns amount for single deposit", func(t *testing.T) {
		total := service.CalculateTotalAmountOfDeposits([]*types.DvpDeposit{erc1155Deposit(100)})
		assert.Equal(t, big.NewInt(100), total)
	})

	t.Run("sums multiple deposits", func(t *testing.T) {
		deposits := []*types.DvpDeposit{erc1155Deposit(100), erc1155Deposit(200), erc1155Deposit(300)}
		total := service.CalculateTotalAmountOfDeposits(deposits)
		assert.Equal(t, big.NewInt(600), total)
	})

	t.Run("handles large amounts", func(t *testing.T) {
		largeAmount := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
		deposit := testdata.NewDvpDeposit(testdata.WithDepositTokenAmount(largeAmount))
		total := service.CalculateTotalAmountOfDeposits([]*types.DvpDeposit{deposit})
		assert.Equal(t, largeAmount, total)
	})
}
