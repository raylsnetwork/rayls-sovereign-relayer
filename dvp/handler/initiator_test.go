package handler_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	ps "github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/ParticipantStorageV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/dvp"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/dvp/handler"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/dvp/testdata"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	keyspb "github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/keys"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

var (
	defaultEncapsulationKeyHex = testdata.MlkemEncapsulationKeyHex()
	defaultDecapsulationKeyHex = testdata.MlkemDecapsulationKeyHex()
)

// initDeps bundles every collaborator NewDvpInitiator needs.
type initDeps struct {
	swapRepo             *initiatorSwapRepositoryMock
	depositRepo          *initiatorDepositRepositoryMock
	psClient             *initiatorParticipantStorageClientMock
	kosClient            *initiatorKeysClientMock
	ccEndpoint           *initiatorEndpointClientMock
	plEndpoint           *initiatorEndpointClientMock
	erc721Client         *InitiatorERC721ClientMock
	erc721HandlerClient  *InitiatorERC721HandlerClientMock
	erc1155Client        *InitiatorERC1155ClientMock
	erc1155HandlerClient *InitiatorERC1155HandlerClientMock
	dvpClient            *initiatorDvpClientMock
	encryptor            *initiatorEncryptorMock
	depositFinder        *initiatorDepositFinderMock
	consolidator         *initiatorDepositConsolidatorMock
	commitmentCalc       *initiatorCommitmentCalculatorMock
	proofGen             *initiatorProofGeneratorMock
	swapWaiter           *initiatorSwapWaiterMock
	swapAgreement        *initiatorSwapAgreementMock
	nodeClient           *initiatorEthClientMock
	hubClient            *initiatorEthClientMock
	txManager            *initiatorTxManagerMock
}

// newDefaultDeps wires every mock with no-op default behavior so individual
// tests need only override the methods relevant to their assertion.
func newDefaultDeps() *initDeps {
	successfulBlock := func(ctx context.Context, number *big.Int) (*ethtypes.Block, error) {
		blockNumber := number
		if blockNumber == nil {
			blockNumber = big.NewInt(1)
		}
		return ethtypes.NewBlockWithHeader(&ethtypes.Header{
			Number: blockNumber,
			Time:   0,
		}), nil
	}

	deps := &initDeps{
		swapRepo: &initiatorSwapRepositoryMock{
			GetSwapBySharedIDFunc: func(ctx context.Context, sharedId string) (*types.DvpSwap, error) {
				return nil, nil
			},
			CreateSwapFunc: func(ctx context.Context, swap *types.DvpSwap) error {
				return nil
			},
			UpdateSwapStatusFunc: func(ctx context.Context, sharedID string, status types.DvpSwapStatus) error {
				return nil
			},
			UpdateSwapFromFunc: func(ctx context.Context, sharedID, from string) error {
				return nil
			},
		},
		depositRepo: &initiatorDepositRepositoryMock{
			CreateDepositFunc: func(ctx context.Context, deposit *types.DvpDeposit) error {
				return nil
			},
			UpdateDepositStatusFunc: func(ctx context.Context, commitment *big.Int, status types.DvpDepositStatus) error {
				return nil
			},
			BatchUpdateStatusForCommitmentsFunc: func(ctx context.Context, commitments []string, status types.DvpDepositStatus) error {
				return nil
			},
			BatchUpsertNullifiersFunc: func(ctx context.Context, m map[string]string) error {
				return nil
			},
			UpsertDepositNullifierFunc: func(ctx context.Context, commitment *big.Int, nullifier *big.Int) error {
				return nil
			},
			GetNonFungibleDepositFunc: func(ctx context.Context, tokenId, tokenAddress, userAddress string, tokenType types.DvpTokenType, status types.DvpDepositStatus) (*types.DvpDeposit, error) {
				return nil, nil
			},
		},
		psClient: &initiatorParticipantStorageClientMock{
			GetPaymentSpendPublicKeyFunc: func(_ context.Context, chainID *big.Int) (*big.Int, error) {
				return big.NewInt(123), nil
			},
			GetChainViewDataFunc: func(_ context.Context, chainID, blockNumber *big.Int) (ps.ParticipantStructsPrivacyNodeViewData, error) {
				return ps.ParticipantStructsPrivacyNodeViewData{
					ChainId:            chainID,
					RaylsViewPublicKey: defaultEncapsulationKeyHex,
				}, nil
			},
		},
		kosClient: &initiatorKeysClientMock{
			GetPaymentSpendKeyFunc: func(ctx context.Context, in *keyspb.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keyspb.PaymentSpendKeyResponse, error) {
				return &keyspb.PaymentSpendKeyResponse{
					PublicKey: big.NewInt(123).Bytes(),
				}, nil
			},
			GetViewPublicKeyFunc: func(ctx context.Context, in *keyspb.GetViewPublicKeyRequest, opts ...grpc.CallOption) (*keyspb.GetViewPublicKeyResponse, error) {
				return &keyspb.GetViewPublicKeyResponse{
					PublicKey: defaultEncapsulationKeyHex,
				}, nil
			},
		},
		ccEndpoint: &initiatorEndpointClientMock{
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return common.HexToAddress("0xCCResource"), nil
			},
		},
		plEndpoint: &initiatorEndpointClientMock{
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return common.HexToAddress("0xPLResource"), nil
			},
		},
		erc721Client: &InitiatorERC721ClientMock{
			ApproveFunc: func(ctx context.Context, _ string, tokenAddress, to common.Address, nftId *big.Int) error {
				return nil
			},
			BurnFunc: func(ctx context.Context, _ string, tokenAddress common.Address, nftId *big.Int) error {
				return nil
			},
			UpdateExtraDataFunc: func(ctx context.Context, _ string, tokenAddress common.Address, nftId, chainId *big.Int, extraDataBytes []byte, newOwner common.Address) error {
				return nil
			},
			MintBatchFunc: func(ctx context.Context, mintDatas []*dvp.DvpERC721MintData) (map[string]contractclient.BatchResult, error) {
				results := make(map[string]contractclient.BatchResult, len(mintDatas))
				for _, md := range mintDatas {
					results[md.GetID()] = contractclient.BatchResult{
						Receipt: &ethtypes.Receipt{Status: ethtypes.ReceiptStatusSuccessful},
					}
				}
				return results, nil
			},
		},
		erc721HandlerClient: &InitiatorERC721HandlerClientMock{
			GetTotalSupplyFunc: func(ctx context.Context, tokenAddress common.Address) ([]*big.Int, error) {
				return []*big.Int{}, nil
			},
			GetExtraDataFunc: func(ctx context.Context, tokenAddress common.Address, nftId *big.Int) ([]byte, error) {
				return []byte{}, nil
			},
			UnlockFunc: func(ctx context.Context, _ string, tokenAddress common.Address, nftId *big.Int) error {
				return nil
			},
		},
		erc1155Client: &InitiatorERC1155ClientMock{
			ApproveFunc: func(ctx context.Context, _ string, tokenAddress, to common.Address) error {
				return nil
			},
			BurnFunc: func(ctx context.Context, _ string, tokenAddress, tokenOwner common.Address, tokenId, tokenAmount *big.Int) error {
				return nil
			},
			UpdateExtraDataFunc: func(ctx context.Context, _ string, tokenAddress common.Address, tokenId, tokenAmount, chainId *big.Int, extraDataBytes []byte, newOwner common.Address) error {
				return nil
			},
			MintBatchFunc: func(ctx context.Context, mintDatas []*dvp.DvpERC1155MintData) (map[string]contractclient.BatchResult, error) {
				results := make(map[string]contractclient.BatchResult, len(mintDatas))
				for _, md := range mintDatas {
					results[md.GetID()] = contractclient.BatchResult{
						Receipt: &ethtypes.Receipt{Status: ethtypes.ReceiptStatusSuccessful},
					}
				}
				return results, nil
			},
		},
		erc1155HandlerClient: &InitiatorERC1155HandlerClientMock{
			GetAllTokenIdsWithSupplyFunc: func(ctx context.Context, tokenAddress common.Address) ([]dvp.DvpERC1155Supply, error) {
				return []dvp.DvpERC1155Supply{}, nil
			},
			GetTokenExtraDataFunc: func(ctx context.Context, tokenAddress common.Address, tokenId *big.Int) ([]byte, error) {
				return []byte{}, nil
			},
			UnlockFunc: func(ctx context.Context, _ string, tokenAddress common.Address, tokenId, tokenAmount *big.Int, to common.Address) error {
				return nil
			},
		},
		dvpClient: &initiatorDvpClientMock{
			DepositERC721Func: func(ctx context.Context, _ string, contractAddress common.Address, nftId, publicKey, salt *big.Int, encryptedBalanceUpdate []byte) error {
				return nil
			},
			WithdrawERC721Func: func(ctx context.Context, _ string, contractAddress common.Address, nftId *big.Int, to common.Address, salt *big.Int, proof *dvp.ProofReceipt, encryptedBalanceUpdate []byte) error {
				return nil
			},
			DepositERC1155Func: func(ctx context.Context, _ string, contractAddress common.Address, tokenId, tokenAmount *big.Int, tokenData []byte, publicKey, salt *big.Int, encryptedBalanceUpdate []byte) error {
				return nil
			},
			WithdrawERC1155Func: func(ctx context.Context, _ string, contractAddress common.Address, tokenId, tokenAmount *big.Int, to common.Address, salt *big.Int, proof *dvp.ProofReceipt, encryptedBalanceUpdate []byte) error {
				return nil
			},
			InitiateSwapFunc: func(ctx context.Context, salt *big.Int, ciphertext []byte, msg *types.DvpSwapMessage, proof *dvp.ProofReceipt, validityTime uint64, passphrase *big.Int) error {
				return nil
			},
			CompleteSwapFunc: func(ctx context.Context, salt *big.Int, swap *types.DvpSwapMessage, proof *dvp.ProofReceipt) error {
				return nil
			},
			CancelSwapFunc: func(ctx context.Context, sharedId string, preimage *big.Int) error {
				return nil
			},
		},
		encryptor: &initiatorEncryptorMock{
			EncryptDvpBalanceUpdatedFunc: func(_ context.Context, message types.DvpBalanceUpdated) ([]byte, error) {
				return []byte{}, nil
			},
		},
		depositFinder: &initiatorDepositFinderMock{
			FindEnygmaDepositsFunc: func(ctx context.Context, from, tokenInAddress string, tokenInAmount *big.Int) ([]*types.DvpDeposit, error) {
				return []*types.DvpDeposit{}, nil
			},
			FindERC1155DepositsForJSProofFunc: func(ctx context.Context, userAddress, tokenAddress, tokenId string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
				return []*types.DvpDeposit{}, nil
			},
			FindERC721DepositFunc: func(ctx context.Context, userAddress, tokenAddress, tokenId string) (*types.DvpDeposit, error) {
				return &types.DvpDeposit{
					Commitment: big.NewInt(1),
				}, nil
			},
		},
		consolidator: &initiatorDepositConsolidatorMock{
			PrepareDepositsForJSProofFunc: func(ctx context.Context, _ string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit) ([]*types.DvpDeposit, error) {
				return deposits, nil
			},
			ConsolidateERC1155DepositsFunc: func(ctx context.Context, _ string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit, consolidationAmount *big.Int) ([]*types.DvpDeposit, error) {
				return deposits, nil
			},
		},
		commitmentCalc: &initiatorCommitmentCalculatorMock{
			CalculateNFTCommitmentFunc: func(spendPK, salt *big.Int, nftID, nftAddress string) (*big.Int, error) {
				return big.NewInt(1), nil
			},
			CalculateERC1155CommitmentFunc: func(spendPK, salt *big.Int, tokenAddress, tokenID string, tokenAmount *big.Int) (*big.Int, error) {
				return big.NewInt(1), nil
			},
			CalculatePaymentCommitmentFunc: func(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error) {
				return big.NewInt(1), nil
			},
			CalculateNullifierFunc: func(spendSK, leafIndex *big.Int) (*big.Int, error) {
				return big.NewInt(1), nil
			},
		},
		proofGen: &initiatorProofGeneratorMock{
			GenerateEnygmaToERC721SwapProofFunc: func(ctx context.Context, swap *types.DvpSwap, deposits []*types.DvpDeposit, sourceViewPublicKey []byte, selfSalt, destSalt, destSpendPubKey *big.Int) (*dvp.ProofReceipt, error) {
				return &dvp.ProofReceipt{Nullifiers: []*big.Int{big.NewInt(1)}}, nil
			},
			GenerateEnygmaToERC1155SwapProofFunc: func(ctx context.Context, swap *types.DvpSwap, deposits []*types.DvpDeposit, sourceViewPublicKey []byte, selfSalt, destSalt, destSpendPubKey *big.Int) (*dvp.ProofReceipt, error) {
				return &dvp.ProofReceipt{Nullifiers: []*big.Int{big.NewInt(1)}}, nil
			},
			GenerateERC721ToEnygmaSwapProofFunc: func(ctx context.Context, swap *types.DvpSwap, deposit *types.DvpDeposit, sourceViewPublicKey []byte, selfSalt, destSalt, destSpendPubKey *big.Int) (*dvp.ProofReceipt, error) {
				return &dvp.ProofReceipt{Nullifiers: []*big.Int{big.NewInt(1)}}, nil
			},
			GenerateERC1155ToEnygmaSwapProofFunc: func(ctx context.Context, swap *types.DvpSwap, deposits []*types.DvpDeposit, sourceViewPublicKey []byte, selfSalt, destSalt, destSpendPubKey *big.Int) (*dvp.ProofReceipt, error) {
				return &dvp.ProofReceipt{Nullifiers: []*big.Int{big.NewInt(1)}}, nil
			},
			GenerateERC721WithdrawProofFunc: func(ctx context.Context, sourceViewPublicKey []byte, destSalt, operatorPublicKey *big.Int, deposit *types.DvpDeposit) (*dvp.ProofReceipt, error) {
				return &dvp.ProofReceipt{Nullifiers: []*big.Int{big.NewInt(1)}}, nil
			},
			GenerateERC1155WithdrawProofFunc: func(ctx context.Context, sourceViewPublicKey []byte, destSalt, operatorPublicKey *big.Int, userAddress, tokenAddress, tokenID string, tokenAmount *big.Int, deposits []*types.DvpDeposit) (*dvp.ProofReceipt, error) {
				return &dvp.ProofReceipt{Nullifiers: []*big.Int{big.NewInt(1)}}, nil
			},
		},
		swapWaiter: &initiatorSwapWaiterMock{
			WaitForSwapInitiationFunc: func(ctx context.Context, sharedId string) (*types.DvpSwap, error) {
				return nil, nil
			},
		},
		swapAgreement: &initiatorSwapAgreementMock{
			VerifyFunc: func(ctx context.Context, swap *types.DvpSwap, destChainId *big.Int, tokenInResourceId string, tokenInAmount *big.Int, tokenInId string, tokenInType types.DvpTokenType, tokenOutResourceId string, tokenOutAmount *big.Int, tokenOutId string, tokenOutType types.DvpTokenType) (string, error) {
				return "", nil
			},
			HandleSwapDisagreementFunc: func(ctx context.Context, sharedId, tokenInResourceId string, tokenInType types.DvpTokenType, message string) error {
				return nil
			},
		},
		nodeClient: &initiatorEthClientMock{
			BlockByNumberFunc: successfulBlock,
		},
		hubClient: &initiatorEthClientMock{
			BlockByNumberFunc: successfulBlock,
		},
		txManager: &initiatorTxManagerMock{
			WithTransactionFunc: func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			},
		},
	}
	return deps
}

func newInitiator(deps *initDeps) *handler.DvpInitiator {
	return handler.NewDvpInitiator(
		handler.InitiatorConfig{ChainID: big.NewInt(1)},
		deps.swapRepo,
		deps.depositRepo,
		deps.psClient,
		deps.kosClient,
		deps.ccEndpoint,
		deps.plEndpoint,
		deps.erc721Client,
		deps.erc721HandlerClient,
		deps.erc1155Client,
		deps.erc1155HandlerClient,
		deps.dvpClient,
		deps.encryptor,
		deps.depositFinder,
		deps.consolidator,
		deps.commitmentCalc,
		deps.proofGen,
		deps.swapWaiter,
		deps.swapAgreement,
		deps.nodeClient,
		deps.hubClient,
		deps.txManager,
	)
}

// preInitiatedSwap returns a swap row in WaitingConfirmation state with both
// salts populated — used by confirmation-branch tests.
func preInitiatedSwap(sharedID string, tokenIn, tokenOut types.DvpTokenType, tokenInResource, tokenOutResource string) *types.DvpSwap {
	swap := testdata.NewDvpSwap(
		testdata.WithSharedID(sharedID),
		testdata.WithStatus(types.DvpSwapWaitingConfirmation),
		testdata.WithLegs(
			tokenIn, tokenInResource, "0xTokenIn", "tok-in", big.NewInt(100),
			tokenOut, tokenOutResource, "0xTokenOut", "tok-out", big.NewInt(200),
		),
	)
	swap.CancelPreimage = big.NewInt(0xCAFE)
	return swap
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestNewDvpInitiator(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("initializes with all dependencies", func(t *testing.T) {
		deps := newDefaultDeps()
		initiator := newInitiator(deps)
		require.NotNil(t, initiator)
	})
}

// ---------------------------------------------------------------------------
// ERC1155 Creation
// ---------------------------------------------------------------------------

func TestHandleERC1155Creation(t *testing.T) {
	ctx := context.Background()
	resourceId := "0xERC1155TokenHandlerAddress"
	tokenHandlerAddress := common.HexToAddress("0xhandler")
	tokenAddress := common.HexToAddress("0xtoken")

	t.Run("successfully creates ERC1155 with initial supply", func(t *testing.T) {
		deps := newDefaultDeps()

		tokenId1 := big.NewInt(1)
		tokenAmount1 := big.NewInt(100)
		tokenId2 := big.NewInt(2)
		tokenAmount2 := big.NewInt(200)

		tokenSupply := []dvp.DvpERC1155Supply{
			{Id: tokenId1, Amount: tokenAmount1},
			{Id: tokenId2, Amount: tokenAmount2},
		}

		extraData1 := []byte("extra_data_1")
		extraData2 := []byte("extra_data_2")

		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenAddress, nil
		}

		deps.erc1155HandlerClient.GetAllTokenIdsWithSupplyFunc = func(ctx context.Context, tokenAddress common.Address) ([]dvp.DvpERC1155Supply, error) {
			return tokenSupply, nil
		}
		deps.erc1155HandlerClient.GetTokenExtraDataFunc = func(ctx context.Context, tokenAddress common.Address, tokenId *big.Int) ([]byte, error) {
			if tokenId.Cmp(tokenId1) == 0 {
				return extraData1, nil
			}
			return extraData2, nil
		}

		var batcherCalled bool
		deps.erc1155Client.MintBatchFunc = func(ctx context.Context, data []*dvp.DvpERC1155MintData) (map[string]contractclient.BatchResult, error) {
			batcherCalled = true
			require.Len(t, data, 2)
			require.Equal(t, tokenId1, data[0].TokenID)
			require.Equal(t, tokenAmount1, data[0].TokenAmount)
			require.Equal(t, extraData1, data[0].ExtraData)
			require.Equal(t, tokenAddress, data[0].TokenAddress)
			require.Equal(t, tokenId2, data[1].TokenID)
			require.Equal(t, tokenAmount2, data[1].TokenAmount)
			require.Equal(t, extraData2, data[1].ExtraData)
			require.Equal(t, tokenAddress, data[1].TokenAddress)
			return map[string]contractclient.BatchResult{
				"1": {Receipt: &ethtypes.Receipt{Status: ethtypes.ReceiptStatusSuccessful}},
				"2": {Receipt: &ethtypes.Receipt{Status: ethtypes.ReceiptStatusSuccessful}},
			}, nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Creation(ctx, "test-event-id", resourceId)
		require.NoError(t, err)
		assert.True(t, batcherCalled)
	})

	t.Run("skips if ERC1155 has no initial supply", func(t *testing.T) {
		deps := newDefaultDeps()

		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenHandlerAddress, nil
		}
		deps.erc1155HandlerClient.GetAllTokenIdsWithSupplyFunc = func(ctx context.Context, tokenAddress common.Address) ([]dvp.DvpERC1155Supply, error) {
			return []dvp.DvpERC1155Supply{}, nil
		}

		var batcherCalled bool
		deps.erc1155Client.MintBatchFunc = func(ctx context.Context, data []*dvp.DvpERC1155MintData) (map[string]contractclient.BatchResult, error) {
			batcherCalled = true
			t.Fatal("MintBatch should not be called when there is no supply")
			return nil, nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Creation(ctx, "test-event-id", resourceId)
		require.NoError(t, err)
		assert.False(t, batcherCalled)
	})

	t.Run("returns error if plEndpointClient.GetResourceAddress fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.Address{}, errors.New("endpoint error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Creation(ctx, "test-event-id", resourceId)
		require.Error(t, err)
		assert.ErrorContains(t, err, "endpoint error")
	})

	t.Run("returns error if GetAllTokenIdsWithSupply fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenHandlerAddress, nil
		}
		deps.erc1155HandlerClient.GetAllTokenIdsWithSupplyFunc = func(ctx context.Context, tokenAddress common.Address) ([]dvp.DvpERC1155Supply, error) {
			return nil, errors.New("supply retrieval error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Creation(ctx, "test-event-id", resourceId)
		require.Error(t, err)
		assert.ErrorContains(t, err, "supply retrieval error")
	})

	t.Run("returns error if ccEndpointClient.GetResourceAddress fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.Address{}, errors.New("cc endpoint error")
		}
		deps.erc1155HandlerClient.GetAllTokenIdsWithSupplyFunc = func(ctx context.Context, tokenAddress common.Address) ([]dvp.DvpERC1155Supply, error) {
			return []dvp.DvpERC1155Supply{{Id: big.NewInt(1), Amount: big.NewInt(100)}}, nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Creation(ctx, "test-event-id", resourceId)
		require.Error(t, err)
		assert.ErrorContains(t, err, "cc endpoint error")
	})

	t.Run("returns error if GetTokenExtraData fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenAddress, nil
		}
		deps.erc1155HandlerClient.GetAllTokenIdsWithSupplyFunc = func(ctx context.Context, tokenAddress common.Address) ([]dvp.DvpERC1155Supply, error) {
			return []dvp.DvpERC1155Supply{{Id: big.NewInt(1), Amount: big.NewInt(100)}}, nil
		}
		deps.erc1155HandlerClient.GetTokenExtraDataFunc = func(ctx context.Context, tokenAddress common.Address, tokenId *big.Int) ([]byte, error) {
			return nil, errors.New("extra data retrieval error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Creation(ctx, "test-event-id", resourceId)
		require.Error(t, err)
		assert.ErrorContains(t, err, "extra data retrieval error")
	})

	t.Run("returns error if MintBatch fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenAddress, nil
		}
		deps.erc1155HandlerClient.GetAllTokenIdsWithSupplyFunc = func(ctx context.Context, tokenAddress common.Address) ([]dvp.DvpERC1155Supply, error) {
			return []dvp.DvpERC1155Supply{{Id: big.NewInt(1), Amount: big.NewInt(100)}}, nil
		}
		deps.erc1155Client.MintBatchFunc = func(ctx context.Context, data []*dvp.DvpERC1155MintData) (map[string]contractclient.BatchResult, error) {
			return nil, errors.New("batcher send error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Creation(ctx, "test-event-id", resourceId)
		require.Error(t, err)
		assert.ErrorContains(t, err, "batcher send error")
	})

	t.Run("handles failed mint results", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenAddress, nil
		}
		deps.erc1155HandlerClient.GetAllTokenIdsWithSupplyFunc = func(ctx context.Context, tokenAddress common.Address) ([]dvp.DvpERC1155Supply, error) {
			return []dvp.DvpERC1155Supply{
				{Id: big.NewInt(1), Amount: big.NewInt(100)},
				{Id: big.NewInt(2), Amount: big.NewInt(200)},
			}, nil
		}
		failedReceipt := &ethtypes.Receipt{Status: ethtypes.ReceiptStatusFailed}
		deps.erc1155Client.MintBatchFunc = func(ctx context.Context, data []*dvp.DvpERC1155MintData) (map[string]contractclient.BatchResult, error) {
			return map[string]contractclient.BatchResult{
				"1": {Receipt: failedReceipt, Err: nil},
				"2": {Receipt: nil, Err: errors.New("mint failed")},
			}, nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Creation(ctx, "test-event-id", resourceId)
		require.NoError(t, err)
	})
}

// ---------------------------------------------------------------------------
// ERC721 Creation
// ---------------------------------------------------------------------------

func TestHandleERC721Creation(t *testing.T) {
	ctx := context.Background()
	resourceId := "0xERC721TokenHandlerAddress"
	nftHandlerAddress := common.HexToAddress("0xhandler")
	nftAddress := common.HexToAddress("0xtoken")

	t.Run("successfully creates ERC721 with initial supply", func(t *testing.T) {
		deps := newDefaultDeps()

		nftId1 := big.NewInt(1)
		nftId2 := big.NewInt(2)
		nftIds := []*big.Int{nftId1, nftId2}
		extraData1 := []byte("extra_data_nft_1")
		extraData2 := []byte("extra_data_nft_2")

		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftAddress, nil
		}
		deps.erc721HandlerClient.GetTotalSupplyFunc = func(ctx context.Context, tokenAddress common.Address) ([]*big.Int, error) {
			return nftIds, nil
		}
		deps.erc721HandlerClient.GetExtraDataFunc = func(ctx context.Context, tokenAddress common.Address, nftId *big.Int) ([]byte, error) {
			if nftId.Cmp(nftId1) == 0 {
				return extraData1, nil
			}
			return extraData2, nil
		}

		var batcherCalled bool
		deps.erc721Client.MintBatchFunc = func(ctx context.Context, data []*dvp.DvpERC721MintData) (map[string]contractclient.BatchResult, error) {
			batcherCalled = true
			require.Len(t, data, 2)
			require.Equal(t, nftId1, data[0].TokenID)
			require.Equal(t, nftAddress, data[0].TokenAddress)
			require.Equal(t, extraData1, data[0].ExtraData)
			require.Equal(t, nftId2, data[1].TokenID)
			require.Equal(t, nftAddress, data[1].TokenAddress)
			require.Equal(t, extraData2, data[1].ExtraData)
			return map[string]contractclient.BatchResult{
				"1": {Receipt: &ethtypes.Receipt{Status: ethtypes.ReceiptStatusSuccessful}},
				"2": {Receipt: &ethtypes.Receipt{Status: ethtypes.ReceiptStatusSuccessful}},
			}, nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Creation(ctx, "test-event-id", resourceId)
		require.NoError(t, err)
		assert.True(t, batcherCalled)
	})

	t.Run("returns error if plEndpointClient GetResourceAddress fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.Address{}, errors.New("endpoint error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Creation(ctx, "test-event-id", resourceId)
		require.Error(t, err)
		assert.ErrorContains(t, err, "endpoint error")
	})

	t.Run("returns error if GetTotalSupply fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftHandlerAddress, nil
		}
		deps.erc721HandlerClient.GetTotalSupplyFunc = func(ctx context.Context, tokenAddress common.Address) ([]*big.Int, error) {
			return nil, errors.New("supply error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Creation(ctx, "test-event-id", resourceId)
		require.Error(t, err)
		assert.ErrorContains(t, err, "supply error")
	})

	t.Run("returns nil if no initial supply", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftHandlerAddress, nil
		}
		deps.erc721HandlerClient.GetTotalSupplyFunc = func(ctx context.Context, tokenAddress common.Address) ([]*big.Int, error) {
			return []*big.Int{}, nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Creation(ctx, "test-event-id", resourceId)
		require.NoError(t, err)
	})

	t.Run("returns error if ccEndpointClient GetResourceAddress fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.Address{}, errors.New("cc endpoint error")
		}
		deps.erc721HandlerClient.GetTotalSupplyFunc = func(ctx context.Context, tokenAddress common.Address) ([]*big.Int, error) {
			return []*big.Int{big.NewInt(1)}, nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Creation(ctx, "test-event-id", resourceId)
		require.Error(t, err)
		assert.ErrorContains(t, err, "cc endpoint error")
	})

	t.Run("returns error if MintBatch fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftAddress, nil
		}
		deps.erc721HandlerClient.GetTotalSupplyFunc = func(ctx context.Context, tokenAddress common.Address) ([]*big.Int, error) {
			return []*big.Int{big.NewInt(1)}, nil
		}
		deps.erc721HandlerClient.GetExtraDataFunc = func(ctx context.Context, tokenAddress common.Address, nftId *big.Int) ([]byte, error) {
			return []byte("extra"), nil
		}
		deps.erc721Client.MintBatchFunc = func(ctx context.Context, data []*dvp.DvpERC721MintData) (map[string]contractclient.BatchResult, error) {
			return nil, errors.New("batcher error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Creation(ctx, "test-event-id", resourceId)
		require.Error(t, err)
		assert.ErrorContains(t, err, "batcher error")
	})

	t.Run("returns error if GetExtraData fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftAddress, nil
		}
		deps.erc721HandlerClient.GetTotalSupplyFunc = func(ctx context.Context, tokenAddress common.Address) ([]*big.Int, error) {
			return []*big.Int{big.NewInt(1)}, nil
		}
		deps.erc721HandlerClient.GetExtraDataFunc = func(ctx context.Context, tokenAddress common.Address, nftId *big.Int) ([]byte, error) {
			return nil, errors.New("extra data fetch failed")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Creation(ctx, "test-event-id", resourceId)
		require.Error(t, err)
		assert.ErrorContains(t, err, "extra data fetch failed")
	})

	t.Run("handles failed mint results", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftAddress, nil
		}
		deps.erc721HandlerClient.GetTotalSupplyFunc = func(ctx context.Context, tokenAddress common.Address) ([]*big.Int, error) {
			return []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}, nil
		}
		deps.erc721HandlerClient.GetExtraDataFunc = func(ctx context.Context, tokenAddress common.Address, nftId *big.Int) ([]byte, error) {
			return []byte("extra_data"), nil
		}
		successReceipt := &ethtypes.Receipt{Status: ethtypes.ReceiptStatusSuccessful}
		failedReceipt := &ethtypes.Receipt{Status: ethtypes.ReceiptStatusFailed}
		deps.erc721Client.MintBatchFunc = func(ctx context.Context, data []*dvp.DvpERC721MintData) (map[string]contractclient.BatchResult, error) {
			return map[string]contractclient.BatchResult{
				"1": {Receipt: successReceipt, Err: nil},
				"2": {Receipt: failedReceipt, Err: nil},
				"3": {Receipt: nil, Err: errors.New("mint failed")},
			}, nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Creation(ctx, "test-event-id", resourceId)
		require.NoError(t, err)
	})
}

// ---------------------------------------------------------------------------
// ERC721 Mint
// ---------------------------------------------------------------------------

func TestHandleERC721Mint(t *testing.T) {
	ctx := context.Background()
	resourceId := "0xERC721TokenHandlerAddress"
	nftHandlerAddress := common.HexToAddress("0xhandler")
	nftAddress := common.HexToAddress("0xnft")
	nftId := big.NewInt(1)
	extraData := []byte("extra_data")

	t.Run("successfully mints ERC721 token", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftAddress, nil
		}
		deps.erc721HandlerClient.GetExtraDataFunc = func(ctx context.Context, tokenAddress common.Address, id *big.Int) ([]byte, error) {
			return extraData, nil
		}

		var batcherCalled bool
		deps.erc721Client.MintBatchFunc = func(ctx context.Context, data []*dvp.DvpERC721MintData) (map[string]contractclient.BatchResult, error) {
			batcherCalled = true
			require.Len(t, data, 1)
			require.Equal(t, nftId, data[0].TokenID)
			require.Equal(t, nftAddress, data[0].TokenAddress)
			require.Equal(t, extraData, data[0].ExtraData)
			return map[string]contractclient.BatchResult{
				"1": {Receipt: &ethtypes.Receipt{Status: ethtypes.ReceiptStatusSuccessful}},
			}, nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Mint(ctx, "test-event-id", resourceId, nftId)
		require.NoError(t, err)
		assert.True(t, batcherCalled)
	})

	t.Run("returns error if plEndpointClient GetResourceAddress fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.Address{}, errors.New("endpoint error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Mint(ctx, "test-event-id", resourceId, nftId)
		require.Error(t, err)
		assert.ErrorContains(t, err, "endpoint error")
	})

	t.Run("returns error if ccEndpointClient GetResourceAddress fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.Address{}, errors.New("cc endpoint error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Mint(ctx, "test-event-id", resourceId, nftId)
		require.Error(t, err)
		assert.ErrorContains(t, err, "cc endpoint error")
	})

	t.Run("returns error if GetExtraData fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftAddress, nil
		}
		deps.erc721HandlerClient.GetExtraDataFunc = func(ctx context.Context, tokenAddress common.Address, id *big.Int) ([]byte, error) {
			return nil, errors.New("extra data error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Mint(ctx, "test-event-id", resourceId, nftId)
		require.Error(t, err)
		assert.ErrorContains(t, err, "extra data error")
	})

	t.Run("returns error if MintBatch fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftAddress, nil
		}
		deps.erc721HandlerClient.GetExtraDataFunc = func(ctx context.Context, tokenAddress common.Address, id *big.Int) ([]byte, error) {
			return extraData, nil
		}
		deps.erc721Client.MintBatchFunc = func(ctx context.Context, data []*dvp.DvpERC721MintData) (map[string]contractclient.BatchResult, error) {
			return nil, errors.New("batcher error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Mint(ctx, "test-event-id", resourceId, nftId)
		require.Error(t, err)
		assert.ErrorContains(t, err, "batcher error")
	})

	t.Run("returns error when mintResult contains Err", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftAddress, nil
		}
		deps.erc721HandlerClient.GetExtraDataFunc = func(ctx context.Context, tokenAddress common.Address, id *big.Int) ([]byte, error) {
			return extraData, nil
		}
		deps.erc721Client.MintBatchFunc = func(ctx context.Context, data []*dvp.DvpERC721MintData) (map[string]contractclient.BatchResult, error) {
			return map[string]contractclient.BatchResult{"1": {Err: errors.New("mint failed")}}, nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Mint(ctx, "test-event-id", resourceId, nftId)
		require.Error(t, err)
		assert.ErrorContains(t, err, "mint failed")
	})
}

// ---------------------------------------------------------------------------
// ERC1155 Mint
// ---------------------------------------------------------------------------

func TestHandleERC1155Mint(t *testing.T) {
	ctx := context.Background()
	resourceId := "0xERC1155TokenHandlerAddress"
	tokenHandlerAddress := common.HexToAddress("0xhandler")
	tokenAddress := common.HexToAddress("0xtoken")
	tokenId := big.NewInt(1)
	tokenAmount := big.NewInt(100)
	tokenData := []byte("token_data")
	extraData := []byte("extra_data")

	t.Run("successfully mints ERC1155 token", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenAddress, nil
		}
		deps.erc1155HandlerClient.GetTokenExtraDataFunc = func(ctx context.Context, tokenAddress common.Address, id *big.Int) ([]byte, error) {
			return extraData, nil
		}

		var batcherCalled bool
		deps.erc1155Client.MintBatchFunc = func(ctx context.Context, data []*dvp.DvpERC1155MintData) (map[string]contractclient.BatchResult, error) {
			batcherCalled = true
			require.Len(t, data, 1)
			require.Equal(t, tokenId, data[0].TokenID)
			require.Equal(t, tokenAmount, data[0].TokenAmount)
			require.Equal(t, tokenData, data[0].Data)
			require.Equal(t, extraData, data[0].ExtraData)
			require.Equal(t, tokenAddress, data[0].TokenAddress)
			return map[string]contractclient.BatchResult{
				"1": {Receipt: &ethtypes.Receipt{Status: ethtypes.ReceiptStatusSuccessful}},
			}, nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Mint(ctx, "test-event-id", resourceId, tokenId, tokenAmount, tokenData)
		require.NoError(t, err)
		assert.True(t, batcherCalled)
	})

	t.Run("returns error if plEndpointClient.GetResourceAddress fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.Address{}, errors.New("endpoint error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Mint(ctx, "test-event-id", resourceId, tokenId, tokenAmount, tokenData)
		require.Error(t, err)
		assert.ErrorContains(t, err, "endpoint error")
	})

	t.Run("returns error if ccEndpointClient.GetResourceAddress fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.Address{}, errors.New("cc endpoint error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Mint(ctx, "test-event-id", resourceId, tokenId, tokenAmount, tokenData)
		require.Error(t, err)
		assert.ErrorContains(t, err, "cc endpoint error")
	})

	t.Run("returns error if GetTokenExtraData fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenAddress, nil
		}
		deps.erc1155HandlerClient.GetTokenExtraDataFunc = func(ctx context.Context, tokenAddress common.Address, id *big.Int) ([]byte, error) {
			return nil, errors.New("extra data error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Mint(ctx, "test-event-id", resourceId, tokenId, tokenAmount, tokenData)
		require.Error(t, err)
		assert.ErrorContains(t, err, "extra data error")
	})

	t.Run("returns error if MintBatch fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenAddress, nil
		}
		deps.erc1155HandlerClient.GetTokenExtraDataFunc = func(ctx context.Context, tokenAddress common.Address, id *big.Int) ([]byte, error) {
			return extraData, nil
		}
		deps.erc1155Client.MintBatchFunc = func(ctx context.Context, data []*dvp.DvpERC1155MintData) (map[string]contractclient.BatchResult, error) {
			return nil, errors.New("batcher send error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Mint(ctx, "test-event-id", resourceId, tokenId, tokenAmount, tokenData)
		require.Error(t, err)
		assert.ErrorContains(t, err, "batcher send error")
	})

	t.Run("handles failed batcher receipt", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenAddress, nil
		}
		deps.erc1155HandlerClient.GetTokenExtraDataFunc = func(ctx context.Context, tokenAddress common.Address, id *big.Int) ([]byte, error) {
			return extraData, nil
		}
		failedReceipt := &ethtypes.Receipt{Status: ethtypes.ReceiptStatusFailed}
		deps.erc1155Client.MintBatchFunc = func(ctx context.Context, data []*dvp.DvpERC1155MintData) (map[string]contractclient.BatchResult, error) {
			return map[string]contractclient.BatchResult{"1": {Receipt: failedReceipt}}, nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Mint(ctx, "test-event-id", resourceId, tokenId, tokenAmount, tokenData)
		require.NoError(t, err)
	})

	t.Run("handles mintResult.Err when batcher returns error in result", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenHandlerAddress, nil
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenAddress, nil
		}
		deps.erc1155HandlerClient.GetTokenExtraDataFunc = func(ctx context.Context, tokenAddress common.Address, id *big.Int) ([]byte, error) {
			return extraData, nil
		}
		deps.erc1155Client.MintBatchFunc = func(ctx context.Context, data []*dvp.DvpERC1155MintData) (map[string]contractclient.BatchResult, error) {
			return map[string]contractclient.BatchResult{"1": {Err: errors.New("mint transaction failed")}}, nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Mint(ctx, "test-event-id", resourceId, tokenId, tokenAmount, tokenData)
		require.NoError(t, err)
	})
}

// ---------------------------------------------------------------------------
// ERC1155 Burn
// ---------------------------------------------------------------------------

func TestHandleERC1155Burn(t *testing.T) {
	ctx := context.Background()
	resourceId := "0xERC1155TokenAddress"
	tokenAddress := common.HexToAddress("0xtoken")
	tokenId := big.NewInt(1)
	tokenAmount := big.NewInt(100)

	t.Run("successfully burns ERC1155 token", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenAddress, nil
		}

		var burnCalled bool
		deps.erc1155Client.BurnFunc = func(ctx context.Context, _ string, tokenAddress, tokenOwner common.Address, tid, tamt *big.Int) error {
			burnCalled = true
			return nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Burn(ctx, "test-event-id", resourceId, tokenId, tokenAmount)
		require.NoError(t, err)
		assert.True(t, burnCalled)
	})

	t.Run("returns error if ccEndpointClient.GetResourceAddress fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.Address{}, errors.New("endpoint error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Burn(ctx, "test-event-id", resourceId, tokenId, tokenAmount)
		require.Error(t, err)
		assert.ErrorContains(t, err, "endpoint error")
	})

	t.Run("returns error if tokenClient.Burn fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenAddress, nil
		}
		deps.erc1155Client.BurnFunc = func(ctx context.Context, _ string, tokenAddress, tokenOwner common.Address, tid, tamt *big.Int) error {
			return errors.New("burn error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Burn(ctx, "test-event-id", resourceId, tokenId, tokenAmount)
		require.Error(t, err)
		assert.ErrorContains(t, err, "burn error")
	})
}

// ---------------------------------------------------------------------------
// ERC721 Burn
// ---------------------------------------------------------------------------

func TestHandleERC721Burn(t *testing.T) {
	ctx := context.Background()
	resourceId := "0xERC721TokenAddress"
	nftAddress := common.HexToAddress("0xnft")
	nftId := big.NewInt(1)

	t.Run("successfully burns ERC721 token", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftAddress, nil
		}

		var burnCalled bool
		deps.erc721Client.BurnFunc = func(ctx context.Context, _ string, tokenAddress common.Address, id *big.Int) error {
			burnCalled = true
			return nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Burn(ctx, "test-event-id", resourceId, nftId)
		require.NoError(t, err)
		assert.True(t, burnCalled)
	})

	t.Run("returns error if ccEndpointClient.GetResourceAddress fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.Address{}, errors.New("endpoint error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Burn(ctx, "test-event-id", resourceId, nftId)
		require.Error(t, err)
		assert.ErrorContains(t, err, "endpoint error")
	})

	t.Run("returns error if tokenClient.Burn fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftAddress, nil
		}
		deps.erc721Client.BurnFunc = func(ctx context.Context, _ string, tokenAddress common.Address, id *big.Int) error {
			return errors.New("burn error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Burn(ctx, "test-event-id", resourceId, nftId)
		require.Error(t, err)
		assert.ErrorContains(t, err, "burn error")
	})
}

// ---------------------------------------------------------------------------
// ERC721 Deposit
// ---------------------------------------------------------------------------

func TestHandleERC721Deposit(t *testing.T) {
	ctx := context.Background()
	resourceId := "0xERC721TokenAddress"
	fromAddress := common.HexToAddress("0xfrom")
	nftAddress := common.HexToAddress("0xnft")
	nftId := big.NewInt(1)
	txBlockNumber := big.NewInt(1)

	t.Run("successfully deposits ERC721 token", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftAddress, nil
		}

		var depositCalled bool
		deps.depositRepo.CreateDepositFunc = func(ctx context.Context, deposit *types.DvpDeposit) error {
			depositCalled = true
			assert.Equal(t, fromAddress.Hex(), deposit.UserAddress)
			assert.Equal(t, nftAddress.Hex(), deposit.TokenAddress)
			assert.Equal(t, nftId.String(), deposit.TokenID)
			assert.Equal(t, types.DvpERC721, deposit.TokenType)
			assert.Equal(t, types.DvpDepositPending, deposit.Status)
			return nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Deposit(ctx, "test-event-id", resourceId, nftId, fromAddress, "0xhash", txBlockNumber)
		require.NoError(t, err)
		assert.True(t, depositCalled)
	})

	t.Run("returns error if ccEndpointClient.GetResourceAddress fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.Address{}, errors.New("endpoint error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Deposit(ctx, "test-event-id", resourceId, nftId, fromAddress, "0xhash", txBlockNumber)
		require.Error(t, err)
		assert.ErrorContains(t, err, "endpoint error")
	})

	t.Run("returns error if nftClient.Approve fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftAddress, nil
		}
		deps.erc721Client.ApproveFunc = func(ctx context.Context, _ string, tokenAddress, to common.Address, id *big.Int) error {
			return errors.New("approve error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Deposit(ctx, "test-event-id", resourceId, nftId, fromAddress, "0xhash", txBlockNumber)
		require.Error(t, err)
		assert.ErrorContains(t, err, "approve error")
	})

	t.Run("returns error if kosClient.GetPaymentSpendKey fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftAddress, nil
		}
		deps.kosClient.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keyspb.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keyspb.PaymentSpendKeyResponse, error) {
			return nil, errors.New("key pair error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Deposit(ctx, "test-event-id", resourceId, nftId, fromAddress, "0xhash", txBlockNumber)
		require.Error(t, err)
		assert.ErrorContains(t, err, "key pair error")
	})

	t.Run("returns error if depositRepository.CreateDeposit fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftAddress, nil
		}
		deps.depositRepo.CreateDepositFunc = func(ctx context.Context, deposit *types.DvpDeposit) error {
			return errors.New("deposit error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Deposit(ctx, "test-event-id", resourceId, nftId, fromAddress, "0xhash", txBlockNumber)
		require.Error(t, err)
		assert.ErrorContains(t, err, "deposit error")
	})
}

// ---------------------------------------------------------------------------
// ERC1155 Deposit
// ---------------------------------------------------------------------------

func TestHandleERC1155Deposit(t *testing.T) {
	ctx := context.Background()
	resourceId := "0xERC1155TokenAddress"
	fromAddress := common.HexToAddress("0xfrom")
	tokenAddress := common.HexToAddress("0xtoken")
	tokenId := big.NewInt(1)
	tokenAmount := big.NewInt(100)
	tokenData := []byte("deposit_data")
	txBlockNumber := big.NewInt(1)

	t.Run("successfully deposits ERC1155 token", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenAddress, nil
		}

		var depositCalled bool
		deps.depositRepo.CreateDepositFunc = func(ctx context.Context, deposit *types.DvpDeposit) error {
			depositCalled = true
			assert.Equal(t, fromAddress.Hex(), deposit.UserAddress)
			assert.Equal(t, tokenAmount.String(), deposit.TokenAmount.String())
			assert.Equal(t, tokenAddress.Hex(), deposit.TokenAddress)
			assert.Equal(t, tokenId.String(), deposit.TokenID)
			assert.Equal(t, types.DvpERC1155, deposit.TokenType)
			assert.Equal(t, types.DvpDepositPending, deposit.Status)
			return nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Deposit(ctx, "test-event-id", resourceId, tokenId, tokenAmount, tokenData, fromAddress, "0xhash", txBlockNumber)
		require.NoError(t, err)
		assert.True(t, depositCalled)
	})

	t.Run("returns error if ccEndpointClient.GetResourceAddress fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.Address{}, errors.New("endpoint error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Deposit(ctx, "test-event-id", resourceId, tokenId, tokenAmount, tokenData, fromAddress, "0xhash", txBlockNumber)
		require.Error(t, err)
		assert.ErrorContains(t, err, "endpoint error")
	})

	t.Run("returns error if tokenClient.Approve fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenAddress, nil
		}
		deps.erc1155Client.ApproveFunc = func(ctx context.Context, _ string, tokenAddress, to common.Address) error {
			return errors.New("approve error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Deposit(ctx, "test-event-id", resourceId, tokenId, tokenAmount, tokenData, fromAddress, "0xhash", txBlockNumber)
		require.Error(t, err)
		assert.ErrorContains(t, err, "approve error")
	})

	t.Run("returns error if kosClient.GetPaymentSpendKey fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenAddress, nil
		}
		deps.kosClient.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keyspb.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keyspb.PaymentSpendKeyResponse, error) {
			return nil, errors.New("key pair error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Deposit(ctx, "test-event-id", resourceId, tokenId, tokenAmount, tokenData, fromAddress, "0xhash", txBlockNumber)
		require.Error(t, err)
		assert.ErrorContains(t, err, "key pair error")
	})

	t.Run("returns error if depositRepository.CreateDeposit fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenAddress, nil
		}
		deps.depositRepo.CreateDepositFunc = func(ctx context.Context, deposit *types.DvpDeposit) error {
			return errors.New("deposit error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Deposit(ctx, "test-event-id", resourceId, tokenId, tokenAmount, tokenData, fromAddress, "0xhash", txBlockNumber)
		require.Error(t, err)
		assert.ErrorContains(t, err, "deposit error")
	})
}

// ---------------------------------------------------------------------------
// ERC721 Withdraw
// ---------------------------------------------------------------------------

func TestHandleERC721Withdraw(t *testing.T) {
	ctx := context.Background()
	resourceId := "0xERC721TokenAddress"
	fromAddress := common.HexToAddress("0xfrom")
	nftAddress := common.HexToAddress("0xnft")
	nftHandlerAddress := common.HexToAddress("0xhandler")
	nftId := big.NewInt(1)
	extraData := []byte("extra_data")
	txBlockNumber := big.NewInt(1)

	deposit := &types.DvpDeposit{
		UserAddress:  fromAddress.Hex(),
		TokenAddress: nftAddress.Hex(),
		TokenID:      nftId.String(),
		TokenAmount:  big.NewInt(1),
		TokenType:    types.DvpERC721,
		Status:       types.DvpDepositUnspent,
		Commitment:   big.NewInt(12345),
	}

	t.Run("successfully withdraws ERC721 NFT", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftAddress, nil
		}
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftHandlerAddress, nil
		}
		deps.depositRepo.GetNonFungibleDepositFunc = func(ctx context.Context, tokenId, tokenAddr, userAddr string, tokenType types.DvpTokenType, status types.DvpDepositStatus) (*types.DvpDeposit, error) {
			return deposit, nil
		}
		deps.erc721HandlerClient.GetExtraDataFunc = func(ctx context.Context, tokenAddress common.Address, id *big.Int) ([]byte, error) {
			return extraData, nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Withdraw(ctx, "test-event-id", resourceId, nftId, fromAddress, "0xhash", txBlockNumber)
		require.NoError(t, err)
	})

	t.Run("returns error if ccEndpointClient.GetResourceAddress fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.Address{}, errors.New("endpoint error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Withdraw(ctx, "test-event-id", resourceId, nftId, fromAddress, "0xhash", txBlockNumber)
		require.Error(t, err)
		assert.ErrorContains(t, err, "endpoint error")
	})

	t.Run("returns error if GetNonFungibleDeposit fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftAddress, nil
		}
		deps.depositRepo.GetNonFungibleDepositFunc = func(ctx context.Context, tokenId, tokenAddr, userAddr string, tokenType types.DvpTokenType, status types.DvpDepositStatus) (*types.DvpDeposit, error) {
			return nil, errors.New("deposit lookup error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Withdraw(ctx, "test-event-id", resourceId, nftId, fromAddress, "0xhash", txBlockNumber)
		require.Error(t, err)
		assert.ErrorContains(t, err, "deposit lookup error")
	})

	t.Run("returns error if deposit is not found", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftAddress, nil
		}
		deps.depositRepo.GetNonFungibleDepositFunc = func(ctx context.Context, tokenId, tokenAddr, userAddr string, tokenType types.DvpTokenType, status types.DvpDepositStatus) (*types.DvpDeposit, error) {
			return nil, nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Withdraw(ctx, "test-event-id", resourceId, nftId, fromAddress, "0xhash", txBlockNumber)
		require.Error(t, err)
		assert.ErrorContains(t, err, "NFT deposit not found")
	})

	t.Run("returns error if GenerateERC721WithdrawProof fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftAddress, nil
		}
		deps.depositRepo.GetNonFungibleDepositFunc = func(ctx context.Context, tokenId, tokenAddr, userAddr string, tokenType types.DvpTokenType, status types.DvpDepositStatus) (*types.DvpDeposit, error) {
			return deposit, nil
		}
		deps.proofGen.GenerateERC721WithdrawProofFunc = func(ctx context.Context, sourceViewPublicKey []byte, destSalt, operatorPublicKey *big.Int, d *types.DvpDeposit) (*dvp.ProofReceipt, error) {
			return nil, errors.New("proof generation error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Withdraw(ctx, "test-event-id", resourceId, nftId, fromAddress, "0xhash", txBlockNumber)
		require.Error(t, err)
		assert.ErrorContains(t, err, "proof generation error")
	})

	t.Run("returns error if plEndpointClient.GetResourceAddress fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return nftAddress, nil
		}
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.Address{}, errors.New("handler endpoint error")
		}
		deps.depositRepo.GetNonFungibleDepositFunc = func(ctx context.Context, tokenId, tokenAddr, userAddr string, tokenType types.DvpTokenType, status types.DvpDepositStatus) (*types.DvpDeposit, error) {
			return deposit, nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC721Withdraw(ctx, "test-event-id", resourceId, nftId, fromAddress, "0xhash", txBlockNumber)
		require.Error(t, err)
		assert.ErrorContains(t, err, "handler endpoint error")
	})
}

// ---------------------------------------------------------------------------
// ERC1155 Withdraw
// ---------------------------------------------------------------------------

func TestHandleERC1155Withdraw(t *testing.T) {
	ctx := context.Background()
	resourceId := "0xERC1155TokenAddress"
	fromAddress := common.HexToAddress("0xfrom")
	tokenAddress := common.HexToAddress("0xtoken")
	tokenHandlerAddress := common.HexToAddress("0xhandler")
	tokenId := big.NewInt(1)
	tokenAmount := big.NewInt(100)
	extraData := []byte("extra_data")
	txBlockNumber := big.NewInt(1)
	commitment1 := big.NewInt(111)

	t.Run("successfully withdraws ERC1155 token with exact amount", func(t *testing.T) {
		deps := newDefaultDeps()

		deposits := []*types.DvpDeposit{
			{
				UserAddress:  fromAddress.Hex(),
				TokenAddress: tokenAddress.Hex(),
				TokenID:      tokenId.String(),
				TokenAmount:  big.NewInt(100),
				Commitment:   commitment1,
				Status:       types.DvpDepositPending,
			},
		}

		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenAddress, nil
		}
		deps.plEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenHandlerAddress, nil
		}
		deps.depositFinder.FindERC1155DepositsForJSProofFunc = func(ctx context.Context, userAddr, tokenAddr, tokenIdStr string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
			return deposits, nil
		}
		deps.erc1155HandlerClient.GetTokenExtraDataFunc = func(ctx context.Context, tokenAddress common.Address, id *big.Int) ([]byte, error) {
			return extraData, nil
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Withdraw(ctx, "test-event-id", resourceId, tokenId, tokenAmount, fromAddress, "0xabc", txBlockNumber)
		require.NoError(t, err)
	})

	t.Run("returns error if ccEndpointClient.GetResourceAddress fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return common.Address{}, errors.New("endpoint error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Withdraw(ctx, "test-event-id", resourceId, tokenId, tokenAmount, fromAddress, "0xabc", txBlockNumber)
		require.Error(t, err)
		assert.ErrorContains(t, err, "endpoint error")
	})

	t.Run("returns error if FindERC1155DepositsForJSProof fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenAddress, nil
		}
		deps.depositFinder.FindERC1155DepositsForJSProofFunc = func(ctx context.Context, userAddr, tokenAddr, tokenIdStr string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
			return nil, errors.New("deposit finder error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Withdraw(ctx, "test-event-id", resourceId, tokenId, tokenAmount, fromAddress, "0xabc", txBlockNumber)
		require.Error(t, err)
		assert.ErrorContains(t, err, "deposit finder error")
	})

	t.Run("returns error if BatchUpdateStatusForCommitments fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deposits := []*types.DvpDeposit{
			{
				UserAddress:  fromAddress.Hex(),
				TokenAddress: tokenAddress.Hex(),
				TokenID:      tokenId.String(),
				TokenAmount:  big.NewInt(100),
				Commitment:   commitment1,
			},
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenAddress, nil
		}
		deps.depositFinder.FindERC1155DepositsForJSProofFunc = func(ctx context.Context, userAddr, tokenAddr, tokenIdStr string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
			return deposits, nil
		}
		deps.depositRepo.BatchUpdateStatusForCommitmentsFunc = func(ctx context.Context, commitments []string, status types.DvpDepositStatus) error {
			return errors.New("batch update error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Withdraw(ctx, "test-event-id", resourceId, tokenId, tokenAmount, fromAddress, "0xabc", txBlockNumber)
		require.Error(t, err)
		assert.ErrorContains(t, err, "batch update error")
	})

	t.Run("returns error if GenerateERC1155WithdrawProof fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deposits := []*types.DvpDeposit{
			{
				UserAddress:  fromAddress.Hex(),
				TokenAddress: tokenAddress.Hex(),
				TokenID:      tokenId.String(),
				TokenAmount:  big.NewInt(100),
				Commitment:   commitment1,
			},
		}
		deps.ccEndpoint.GetResourceAddressFunc = func(context.Context, string) (common.Address, error) {
			return tokenAddress, nil
		}
		deps.depositFinder.FindERC1155DepositsForJSProofFunc = func(ctx context.Context, userAddr, tokenAddr, tokenIdStr string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
			return deposits, nil
		}
		deps.proofGen.GenerateERC1155WithdrawProofFunc = func(ctx context.Context, sourceViewPublicKey []byte, destSalt, operatorPublicKey *big.Int, userAddress, tokenAddr, tokenID string, tamt *big.Int, d []*types.DvpDeposit) (*dvp.ProofReceipt, error) {
			return nil, errors.New("proof generation error")
		}

		initiator := newInitiator(deps)
		err := initiator.HandleERC1155Withdraw(ctx, "test-event-id", resourceId, tokenId, tokenAmount, fromAddress, "0xabc", txBlockNumber)
		require.Error(t, err)
		assert.ErrorContains(t, err, "proof generation error")
	})
}

// --- Swap handler tests ---

func TestHandleEnygmaSwapERC721(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("creates and initiates a new swap", func(t *testing.T) {
		deps := newDefaultDeps()
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC721(
			context.Background(),
			"shared-1",
			big.NewInt(2),
			common.HexToAddress("0xAlice"),
			"resource-eny", big.NewInt(100),
			"resource-721", "tok-7",
			"0xTxHash", big.NewInt(50), 3600,
		)
		require.NoError(t, err)

		require.Len(t, deps.proofGen.GenerateEnygmaToERC721SwapProofCalls(), 1)
		generatedSelfSalt := deps.proofGen.GenerateEnygmaToERC721SwapProofCalls()[0].SelfSalt
		generatedDestSalt := deps.proofGen.GenerateEnygmaToERC721SwapProofCalls()[0].DestSalt
		require.NotNil(t, generatedSelfSalt)
		require.NotNil(t, generatedDestSalt)
		assert.NotEqual(t, generatedSelfSalt, generatedDestSalt)

		require.Len(t, deps.dvpClient.InitiateSwapCalls(), 1)
		call := deps.dvpClient.InitiateSwapCalls()[0]
		assert.Equal(t, generatedDestSalt, call.Salt)
		assert.Equal(t, uint64(3600), call.ValidityTime)
		assert.Equal(t, "shared-1", call.Msg.SharedId)
		assert.Equal(t, generatedSelfSalt, call.Msg.InitiatorSelfSalt)

		require.Len(t, deps.swapRepo.CreateSwapCalls(), 1)
		row := deps.swapRepo.CreateSwapCalls()[0].Swap
		assert.Equal(t, types.DvpSwapInitiated, row.Status)
		assert.Equal(t, generatedSelfSalt, row.SelfSalt)
		assert.Equal(t, generatedDestSalt, row.DestSalt)
		assert.Equal(t, types.DvpEnygma, row.TokenInType)
		assert.Equal(t, types.DvpERC721, row.TokenOutType)
	})

	t.Run("does not persist swap row when contract reports SwapAlreadyExists", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.dvpClient.InitiateSwapFunc = func(ctx context.Context, salt *big.Int, ciphertext []byte, msg *types.DvpSwapMessage, proof *dvp.ProofReceipt, validityTime uint64, passphrase *big.Int) error {
			return dvp.ErrSwapAlreadyInitiated
		}
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC721(
			context.Background(),
			"race-1", big.NewInt(2), common.HexToAddress("0xAlice"),
			"resource-eny", big.NewInt(100),
			"resource-721", "tok-7",
			"0xTxHash", big.NewInt(50), 3600,
		)
		require.Error(t, err)
		require.ErrorIs(t, err, dvp.ErrSwapAlreadyInitiated)
		assert.Len(t, deps.swapRepo.CreateSwapCalls(), 0)
	})

	t.Run("confirms a pre-existing WaitingConfirmation swap with persisted salts", func(t *testing.T) {
		deps := newDefaultDeps()
		swap := preInitiatedSwap("shared-2", types.DvpEnygma, types.DvpERC721, "resource-eny", "resource-721")
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC721(
			context.Background(),
			"shared-2", big.NewInt(2), common.HexToAddress("0xBob"),
			"resource-eny", big.NewInt(100),
			"resource-721", "tok-7",
			"0xTxHash", big.NewInt(50), 3600,
		)
		require.NoError(t, err)

		require.Len(t, deps.swapRepo.UpdateSwapFromCalls(), 1)
		assert.Equal(t, common.HexToAddress("0xBob").Hex(), deps.swapRepo.UpdateSwapFromCalls()[0].From)
		assert.Len(t, deps.swapAgreement.VerifyCalls(), 1)
		require.Len(t, deps.proofGen.GenerateEnygmaToERC721SwapProofCalls(), 1)
		assert.Equal(t, swap.SelfSalt, deps.proofGen.GenerateEnygmaToERC721SwapProofCalls()[0].SelfSalt)
		assert.Equal(t, swap.DestSalt, deps.proofGen.GenerateEnygmaToERC721SwapProofCalls()[0].DestSalt)
		require.Len(t, deps.dvpClient.CompleteSwapCalls(), 1)
		assert.Equal(t, swap.DestSalt, deps.dvpClient.CompleteSwapCalls()[0].Salt)
		assert.Empty(t, deps.dvpClient.InitiateSwapCalls())
		assert.Empty(t, deps.swapRepo.CreateSwapCalls())
	})

	t.Run("propagates dvpClient.InitiateSwap error other than SwapAlreadyExists", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.dvpClient.InitiateSwapFunc = func(ctx context.Context, salt *big.Int, ciphertext []byte, msg *types.DvpSwapMessage, proof *dvp.ProofReceipt, validityTime uint64, passphrase *big.Int) error {
			return errors.New("rpc down")
		}
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC721(
			context.Background(),
			"shared-x", big.NewInt(2), common.HexToAddress("0xAlice"),
			"resource-eny", big.NewInt(100),
			"resource-721", "tok-7",
			"0xTxHash", big.NewInt(50), 3600,
		)
		require.Error(t, err)
		assert.NotErrorIs(t, err, dvp.ErrSwapAlreadyInitiated)
		assert.Empty(t, deps.swapRepo.CreateSwapCalls())
	})

	t.Run("idempotent on already-Initiated swap", func(t *testing.T) {
		deps := newDefaultDeps()
		swap := preInitiatedSwap("shared-i", types.DvpEnygma, types.DvpERC721, "resource-eny", "resource-721")
		swap.Status = types.DvpSwapInitiated
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC721(
			context.Background(),
			"shared-i", big.NewInt(2), common.HexToAddress("0xAlice"),
			"resource-eny", big.NewInt(100),
			"resource-721", "tok-7",
			"0xTxHash", big.NewInt(50), 3600,
		)
		require.NoError(t, err)
		assert.Empty(t, deps.dvpClient.InitiateSwapCalls())
		assert.Empty(t, deps.dvpClient.CompleteSwapCalls())
	})
}

func TestHandleEnygmaSwapERC1155(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("creates and initiates a new swap", func(t *testing.T) {
		deps := newDefaultDeps()
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC1155(
			context.Background(),
			"shared-1155-init", big.NewInt(2), common.HexToAddress("0xAlice"),
			"resource-eny", big.NewInt(100),
			"resource-1155", "tok-3", big.NewInt(7),
			"0xTx", big.NewInt(50), 3600,
		)
		require.NoError(t, err)
		require.Len(t, deps.proofGen.GenerateEnygmaToERC1155SwapProofCalls(), 1)
		require.Len(t, deps.dvpClient.InitiateSwapCalls(), 1)
		require.Len(t, deps.swapRepo.CreateSwapCalls(), 1)
		row := deps.swapRepo.CreateSwapCalls()[0].Swap
		assert.Equal(t, types.DvpSwapInitiated, row.Status)
		assert.Equal(t, types.DvpEnygma, row.TokenInType)
		assert.Equal(t, types.DvpERC1155, row.TokenOutType)
	})

	t.Run("confirms WaitingConfirmation swap with CompleteSwap", func(t *testing.T) {
		deps := newDefaultDeps()
		swap := preInitiatedSwap("shared-1155-conf", types.DvpEnygma, types.DvpERC1155, "resource-eny", "resource-1155")
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC1155(
			context.Background(),
			"shared-1155-conf", big.NewInt(2), common.HexToAddress("0xBob"),
			"resource-eny", big.NewInt(100),
			"resource-1155", "tok-3", big.NewInt(7),
			"0xTx", big.NewInt(50), 3600,
		)
		require.NoError(t, err)
		assert.Len(t, deps.swapAgreement.VerifyCalls(), 1)
		require.Len(t, deps.dvpClient.CompleteSwapCalls(), 1)
		assert.Empty(t, deps.dvpClient.InitiateSwapCalls())
	})

	t.Run("returns error when GetSwapBySharedID fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return nil, errors.New("db down")
		}
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC1155(
			context.Background(),
			"s", big.NewInt(2), common.HexToAddress("0xA"),
			"r-eny", big.NewInt(1), "r-1155", "t", big.NewInt(1),
			"0x", big.NewInt(1), 0,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "getting swap by shared ID")
	})

	t.Run("returns error when nodeClient.BlockByNumber fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.nodeClient.BlockByNumberFunc = func(ctx context.Context, number *big.Int) (*ethtypes.Block, error) {
			return nil, errors.New("rpc error")
		}
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC1155(
			context.Background(),
			"s", big.NewInt(2), common.HexToAddress("0xA"),
			"r-eny", big.NewInt(1), "r-1155", "t", big.NewInt(1),
			"0x", big.NewInt(1), 0,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "getting block by number")
	})

	t.Run("returns error when createSwap fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.Address{}, errors.New("endpoint down")
		}
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC1155(
			context.Background(),
			"s", big.NewInt(2), common.HexToAddress("0xA"),
			"r-eny", big.NewInt(1), "r-1155", "t", big.NewInt(1),
			"0x", big.NewInt(1), 0,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "creating Enygma -> ERC1155 swap")
	})

	t.Run("returns error when prepareSwapProofParams fails via hubClient", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.hubClient.BlockByNumberFunc = func(ctx context.Context, number *big.Int) (*ethtypes.Block, error) {
			return nil, errors.New("hub down")
		}
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC1155(
			context.Background(),
			"s", big.NewInt(2), common.HexToAddress("0xA"),
			"r-eny", big.NewInt(1), "r-1155", "t", big.NewInt(1),
			"0x", big.NewInt(1), 0,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "preparing Enygma -> ERC1155 proof params")
	})

	t.Run("returns error when FindEnygmaDeposits fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.depositFinder.FindEnygmaDepositsFunc = func(ctx context.Context, from, tokenInAddress string, tokenInAmount *big.Int) ([]*types.DvpDeposit, error) {
			return nil, errors.New("finder error")
		}
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC1155(
			context.Background(),
			"s", big.NewInt(2), common.HexToAddress("0xA"),
			"r-eny", big.NewInt(1), "r-1155", "t", big.NewInt(1),
			"0x", big.NewInt(1), 0,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "finding enygma deposits")
	})

	t.Run("returns error when PrepareDepositsForJSProof fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.consolidator.PrepareDepositsForJSProofFunc = func(ctx context.Context, _ string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit) ([]*types.DvpDeposit, error) {
			return nil, errors.New("consolidation error")
		}
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC1155(
			context.Background(),
			"s", big.NewInt(2), common.HexToAddress("0xA"),
			"r-eny", big.NewInt(1), "r-1155", "t", big.NewInt(1),
			"0x", big.NewInt(1), 0,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "preparing deposits for JS proof")
	})

	t.Run("returns error when GenerateEnygmaToERC1155SwapProof fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.proofGen.GenerateEnygmaToERC1155SwapProofFunc = func(ctx context.Context, swap *types.DvpSwap, deposits []*types.DvpDeposit, sourceViewPublicKey []byte, selfSalt, destSalt, destSpendPubKey *big.Int) (*dvp.ProofReceipt, error) {
			return nil, errors.New("proof gen error")
		}
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC1155(
			context.Background(),
			"s", big.NewInt(2), common.HexToAddress("0xA"),
			"r-eny", big.NewInt(1), "r-1155", "t", big.NewInt(1),
			"0x", big.NewInt(1), 0,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "generating Enygma -> ERC1155 proof")
	})

	t.Run("confirmation - returns error when UpdateSwapFrom fails", func(t *testing.T) {
		deps := newDefaultDeps()
		swap := preInitiatedSwap("shared-usf", types.DvpEnygma, types.DvpERC1155, "resource-eny", "resource-1155")
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		deps.swapRepo.UpdateSwapFromFunc = func(ctx context.Context, sharedID string, from string) error {
			return errors.New("update from error")
		}
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC1155(
			context.Background(),
			"shared-usf", big.NewInt(2), common.HexToAddress("0xBob"),
			"resource-eny", big.NewInt(100),
			"resource-1155", "tok-3", big.NewInt(7),
			"0xTx", big.NewInt(50), 3600,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "updating swap from_address")
	})

	t.Run("confirmation - Verify failure triggers HandleSwapDisagreement", func(t *testing.T) {
		deps := newDefaultDeps()
		swap := preInitiatedSwap("shared-disagree", types.DvpEnygma, types.DvpERC1155, "resource-eny", "resource-1155")
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		deps.swapAgreement.VerifyFunc = func(ctx context.Context, _ *types.DvpSwap, _ *big.Int, _ string, _ *big.Int, _ string, _ types.DvpTokenType, _ string, _ *big.Int, _ string, _ types.DvpTokenType) (string, error) {
			return "mismatch", errors.New("amounts disagree")
		}
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC1155(
			context.Background(),
			"shared-disagree", big.NewInt(2), common.HexToAddress("0xBob"),
			"resource-eny", big.NewInt(100),
			"resource-1155", "tok-3", big.NewInt(7),
			"0xTx", big.NewInt(50), 3600,
		)
		require.NoError(t, err)
		assert.Len(t, deps.swapAgreement.HandleSwapDisagreementCalls(), 1)
	})

	t.Run("confirmation - returns error when CompleteSwap fails", func(t *testing.T) {
		deps := newDefaultDeps()
		swap := preInitiatedSwap("shared-cs", types.DvpEnygma, types.DvpERC1155, "resource-eny", "resource-1155")
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		deps.dvpClient.CompleteSwapFunc = func(ctx context.Context, salt *big.Int, swap *types.DvpSwapMessage, proof *dvp.ProofReceipt) error {
			return errors.New("complete swap error")
		}
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC1155(
			context.Background(),
			"shared-cs", big.NewInt(2), common.HexToAddress("0xBob"),
			"resource-eny", big.NewInt(100),
			"resource-1155", "tok-3", big.NewInt(7),
			"0xTx", big.NewInt(50), 3600,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "confirming Enygma -> ERC1155 swap")
	})

	// prepareSwapProofParams indirect error tests
	t.Run("returns error when GetPaymentSpendPublicKey fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.psClient.GetPaymentSpendPublicKeyFunc = func(_ context.Context, chainID *big.Int) (*big.Int, error) {
			return nil, errors.New("ps error")
		}
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC1155(
			context.Background(),
			"s", big.NewInt(2), common.HexToAddress("0xA"),
			"r-eny", big.NewInt(1), "r-1155", "t", big.NewInt(1),
			"0x", big.NewInt(1), 0,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "preparing Enygma -> ERC1155 proof params")
	})

	t.Run("returns error when GetChainViewData fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.psClient.GetChainViewDataFunc = func(_ context.Context, chainID, blockNumber *big.Int) (ps.ParticipantStructsPrivacyNodeViewData, error) {
			return ps.ParticipantStructsPrivacyNodeViewData{}, errors.New("view data error")
		}
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC1155(
			context.Background(),
			"s", big.NewInt(2), common.HexToAddress("0xA"),
			"r-eny", big.NewInt(1), "r-1155", "t", big.NewInt(1),
			"0x", big.NewInt(1), 0,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "preparing Enygma -> ERC1155 proof params")
	})

	t.Run("returns error when dest view key hex decode fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.psClient.GetChainViewDataFunc = func(_ context.Context, chainID, blockNumber *big.Int) (ps.ParticipantStructsPrivacyNodeViewData, error) {
			return ps.ParticipantStructsPrivacyNodeViewData{
				ChainId:            chainID,
				RaylsViewPublicKey: "ZZZZ",
			}, nil
		}
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC1155(
			context.Background(),
			"s", big.NewInt(2), common.HexToAddress("0xA"),
			"r-eny", big.NewInt(1), "r-1155", "t", big.NewInt(1),
			"0x", big.NewInt(1), 0,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "preparing Enygma -> ERC1155 proof params")
	})

	t.Run("returns error when GetViewPublicKey fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.kosClient.GetViewPublicKeyFunc = func(ctx context.Context, in *keyspb.GetViewPublicKeyRequest, opts ...grpc.CallOption) (*keyspb.GetViewPublicKeyResponse, error) {
			return nil, errors.New("kos error")
		}
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC1155(
			context.Background(),
			"s", big.NewInt(2), common.HexToAddress("0xA"),
			"r-eny", big.NewInt(1), "r-1155", "t", big.NewInt(1),
			"0x", big.NewInt(1), 0,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "preparing Enygma -> ERC1155 proof params")
	})

	t.Run("returns error when self view key hex decode fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.kosClient.GetViewPublicKeyFunc = func(ctx context.Context, in *keyspb.GetViewPublicKeyRequest, opts ...grpc.CallOption) (*keyspb.GetViewPublicKeyResponse, error) {
			return &keyspb.GetViewPublicKeyResponse{
				PublicKey: "ZZZZ",
			}, nil
		}
		init := newInitiator(deps)

		err := init.HandleEnygmaSwapERC1155(
			context.Background(),
			"s", big.NewInt(2), common.HexToAddress("0xA"),
			"r-eny", big.NewInt(1), "r-1155", "t", big.NewInt(1),
			"0x", big.NewInt(1), 0,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "preparing Enygma -> ERC1155 proof params")
	})
}

func TestHandleERC721SwapEnygma(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("creates and initiates a new swap", func(t *testing.T) {
		deps := newDefaultDeps()
		init := newInitiator(deps)

		err := init.HandleERC721SwapEnygma(
			context.Background(),
			"shared-721-init", big.NewInt(2), common.HexToAddress("0xAlice"),
			"resource-721", "tok-1",
			"resource-eny", big.NewInt(100),
			"0xTx", big.NewInt(50), 3600,
		)
		require.NoError(t, err)
		require.Len(t, deps.proofGen.GenerateERC721ToEnygmaSwapProofCalls(), 1)
		require.Len(t, deps.dvpClient.InitiateSwapCalls(), 1)
		require.Len(t, deps.swapRepo.CreateSwapCalls(), 1)
		row := deps.swapRepo.CreateSwapCalls()[0].Swap
		assert.Equal(t, types.DvpSwapInitiated, row.Status)
		assert.Equal(t, types.DvpERC721, row.TokenInType)
		assert.Equal(t, types.DvpEnygma, row.TokenOutType)
	})

	t.Run("confirms WaitingConfirmation swap with CompleteSwap", func(t *testing.T) {
		deps := newDefaultDeps()
		swap := preInitiatedSwap("shared-721-conf", types.DvpERC721, types.DvpEnygma, "resource-721", "resource-eny")
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		init := newInitiator(deps)

		err := init.HandleERC721SwapEnygma(
			context.Background(),
			"shared-721-conf", big.NewInt(2), common.HexToAddress("0xBob"),
			"resource-721", "tok-1",
			"resource-eny", big.NewInt(100),
			"0xTx", big.NewInt(50), 3600,
		)
		require.NoError(t, err)
		require.Len(t, deps.dvpClient.CompleteSwapCalls(), 1)
		assert.Empty(t, deps.dvpClient.InitiateSwapCalls())
	})
}

func TestHandleERC1155SwapEnygma(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("creates and initiates a new swap", func(t *testing.T) {
		deps := newDefaultDeps()
		init := newInitiator(deps)

		err := init.HandleERC1155SwapEnygma(
			context.Background(),
			"shared-init", big.NewInt(2), common.HexToAddress("0xAlice"),
			"resource-1155", big.NewInt(7), "tok-3",
			"resource-eny", big.NewInt(100),
			"0xTx", big.NewInt(50), 3600,
		)
		require.NoError(t, err)
		require.Len(t, deps.proofGen.GenerateERC1155ToEnygmaSwapProofCalls(), 1)
		require.Len(t, deps.dvpClient.InitiateSwapCalls(), 1)
		require.Len(t, deps.swapRepo.CreateSwapCalls(), 1)
		row := deps.swapRepo.CreateSwapCalls()[0].Swap
		assert.Equal(t, types.DvpSwapInitiated, row.Status)
		assert.Equal(t, types.DvpERC1155, row.TokenInType)
		assert.Equal(t, types.DvpEnygma, row.TokenOutType)
	})

	t.Run("confirms WaitingConfirmation swap with CompleteSwap", func(t *testing.T) {
		deps := newDefaultDeps()
		swap := preInitiatedSwap("shared-1155-eny-conf", types.DvpERC1155, types.DvpEnygma, "resource-1155", "resource-eny")
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		init := newInitiator(deps)

		err := init.HandleERC1155SwapEnygma(
			context.Background(),
			"shared-1155-eny-conf", big.NewInt(2), common.HexToAddress("0xBob"),
			"resource-1155", big.NewInt(7), "tok-3",
			"resource-eny", big.NewInt(100),
			"0xTx", big.NewInt(50), 3600,
		)
		require.NoError(t, err)
		require.Len(t, deps.dvpClient.CompleteSwapCalls(), 1)
		assert.Empty(t, deps.dvpClient.InitiateSwapCalls())
	})

	t.Run("returns error when GetSwapBySharedID fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return nil, errors.New("db down")
		}
		init := newInitiator(deps)

		err := init.HandleERC1155SwapEnygma(
			context.Background(),
			"s", big.NewInt(2), common.HexToAddress("0xA"),
			"r-1155", big.NewInt(1), "t",
			"r-eny", big.NewInt(1),
			"0x", big.NewInt(1), 0,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "getting swap by shared ID")
	})

	t.Run("returns error when nodeClient.BlockByNumber fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.nodeClient.BlockByNumberFunc = func(ctx context.Context, number *big.Int) (*ethtypes.Block, error) {
			return nil, errors.New("rpc error")
		}
		init := newInitiator(deps)

		err := init.HandleERC1155SwapEnygma(
			context.Background(),
			"s", big.NewInt(2), common.HexToAddress("0xA"),
			"r-1155", big.NewInt(1), "t",
			"r-eny", big.NewInt(1),
			"0x", big.NewInt(1), 0,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "getting block by number")
	})

	t.Run("returns error when createSwap fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.ccEndpoint.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.Address{}, errors.New("endpoint down")
		}
		init := newInitiator(deps)

		err := init.HandleERC1155SwapEnygma(
			context.Background(),
			"s", big.NewInt(2), common.HexToAddress("0xA"),
			"r-1155", big.NewInt(1), "t",
			"r-eny", big.NewInt(1),
			"0x", big.NewInt(1), 0,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "creating ERC1155 -> Enygma swap")
	})

	t.Run("returns error when prepareSwapProofParams fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.hubClient.BlockByNumberFunc = func(ctx context.Context, number *big.Int) (*ethtypes.Block, error) {
			return nil, errors.New("hub down")
		}
		init := newInitiator(deps)

		err := init.HandleERC1155SwapEnygma(
			context.Background(),
			"s", big.NewInt(2), common.HexToAddress("0xA"),
			"r-1155", big.NewInt(1), "t",
			"r-eny", big.NewInt(1),
			"0x", big.NewInt(1), 0,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "preparing ERC1155 -> Enygma proof params")
	})

	t.Run("returns error when FindERC1155DepositsForJSProof fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.depositFinder.FindERC1155DepositsForJSProofFunc = func(ctx context.Context, userAddress, tokenAddress, tokenId string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
			return nil, errors.New("finder error")
		}
		init := newInitiator(deps)

		err := init.HandleERC1155SwapEnygma(
			context.Background(),
			"s", big.NewInt(2), common.HexToAddress("0xA"),
			"r-1155", big.NewInt(1), "t",
			"r-eny", big.NewInt(1),
			"0x", big.NewInt(1), 0,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "finding ERC1155 deposits")
	})

	t.Run("returns error when PrepareDepositsForJSProof fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.consolidator.PrepareDepositsForJSProofFunc = func(ctx context.Context, _ string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit) ([]*types.DvpDeposit, error) {
			return nil, errors.New("consolidation error")
		}
		init := newInitiator(deps)

		err := init.HandleERC1155SwapEnygma(
			context.Background(),
			"s", big.NewInt(2), common.HexToAddress("0xA"),
			"r-1155", big.NewInt(1), "t",
			"r-eny", big.NewInt(1),
			"0x", big.NewInt(1), 0,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "preparing deposits for JS proof")
	})

	t.Run("returns error when GenerateERC1155ToEnygmaSwapProof fails", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.proofGen.GenerateERC1155ToEnygmaSwapProofFunc = func(ctx context.Context, swap *types.DvpSwap, deposits []*types.DvpDeposit, sourceViewPublicKey []byte, selfSalt, destSalt, destSpendPubKey *big.Int) (*dvp.ProofReceipt, error) {
			return nil, errors.New("proof gen error")
		}
		init := newInitiator(deps)

		err := init.HandleERC1155SwapEnygma(
			context.Background(),
			"s", big.NewInt(2), common.HexToAddress("0xA"),
			"r-1155", big.NewInt(1), "t",
			"r-eny", big.NewInt(1),
			"0x", big.NewInt(1), 0,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "generating ERC1155 -> Enygma proof")
	})

	t.Run("confirmation - returns error when UpdateSwapFrom fails", func(t *testing.T) {
		deps := newDefaultDeps()
		swap := preInitiatedSwap("shared-usf", types.DvpERC1155, types.DvpEnygma, "resource-1155", "resource-eny")
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		deps.swapRepo.UpdateSwapFromFunc = func(ctx context.Context, sharedID string, from string) error {
			return errors.New("update from error")
		}
		init := newInitiator(deps)

		err := init.HandleERC1155SwapEnygma(
			context.Background(),
			"shared-usf", big.NewInt(2), common.HexToAddress("0xBob"),
			"resource-1155", big.NewInt(7), "tok-3",
			"resource-eny", big.NewInt(100),
			"0xTx", big.NewInt(50), 3600,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "updating swap from_address")
	})

	t.Run("confirmation - Verify failure triggers HandleSwapDisagreement", func(t *testing.T) {
		deps := newDefaultDeps()
		swap := preInitiatedSwap("shared-disagree", types.DvpERC1155, types.DvpEnygma, "resource-1155", "resource-eny")
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		deps.swapAgreement.VerifyFunc = func(ctx context.Context, _ *types.DvpSwap, _ *big.Int, _ string, _ *big.Int, _ string, _ types.DvpTokenType, _ string, _ *big.Int, _ string, _ types.DvpTokenType) (string, error) {
			return "mismatch", errors.New("amounts disagree")
		}
		init := newInitiator(deps)

		err := init.HandleERC1155SwapEnygma(
			context.Background(),
			"shared-disagree", big.NewInt(2), common.HexToAddress("0xBob"),
			"resource-1155", big.NewInt(7), "tok-3",
			"resource-eny", big.NewInt(100),
			"0xTx", big.NewInt(50), 3600,
		)
		require.NoError(t, err)
		assert.Len(t, deps.swapAgreement.HandleSwapDisagreementCalls(), 1)
	})

	t.Run("confirmation - returns error when CompleteSwap fails", func(t *testing.T) {
		deps := newDefaultDeps()
		swap := preInitiatedSwap("shared-cs", types.DvpERC1155, types.DvpEnygma, "resource-1155", "resource-eny")
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		deps.dvpClient.CompleteSwapFunc = func(ctx context.Context, salt *big.Int, swap *types.DvpSwapMessage, proof *dvp.ProofReceipt) error {
			return errors.New("complete swap error")
		}
		init := newInitiator(deps)

		err := init.HandleERC1155SwapEnygma(
			context.Background(),
			"shared-cs", big.NewInt(2), common.HexToAddress("0xBob"),
			"resource-1155", big.NewInt(7), "tok-3",
			"resource-eny", big.NewInt(100),
			"0xTx", big.NewInt(50), 3600,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "confirming ERC1155 -> Enygma swap")
	})
}

func TestHandleSwapCancellation(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("verifies agreement then calls dvpClient.CancelSwap", func(t *testing.T) {
		deps := newDefaultDeps()
		swap := preInitiatedSwap("shared-cancel", types.DvpEnygma, types.DvpERC721, "resource-eny", "resource-721")
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		init := newInitiator(deps)

		err := init.HandleSwapCancellation(
			context.Background(),
			"shared-cancel", big.NewInt(2),
			"resource-eny", big.NewInt(100), "tok-in", types.DvpEnygma,
			"resource-721", big.NewInt(1), "tok-out", types.DvpERC721,
		)
		require.NoError(t, err)
		assert.Len(t, deps.swapAgreement.VerifyCalls(), 1)
		require.Len(t, deps.dvpClient.CancelSwapCalls(), 1)
		assert.Equal(t, "shared-cancel", deps.dvpClient.CancelSwapCalls()[0].SharedId)
		assert.Equal(t, big.NewInt(0xCAFE), deps.dvpClient.CancelSwapCalls()[0].Preimage)
	})

	t.Run("returns error when swap has no cancel preimage", func(t *testing.T) {
		deps := newDefaultDeps()
		swap := preInitiatedSwap("shared-no-preimage", types.DvpEnygma, types.DvpERC721, "resource-eny", "resource-721")
		swap.CancelPreimage = nil
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		init := newInitiator(deps)

		err := init.HandleSwapCancellation(
			context.Background(),
			"shared-no-preimage", big.NewInt(2),
			"resource-eny", big.NewInt(100), "tok-in", types.DvpEnygma,
			"resource-721", big.NewInt(1), "tok-out", types.DvpERC721,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "has no cancel preimage")
		assert.Empty(t, deps.dvpClient.CancelSwapCalls())
	})

	t.Run("returns error when swap not found", func(t *testing.T) {
		deps := newDefaultDeps()
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return nil, nil
		}
		init := newInitiator(deps)

		err := init.HandleSwapCancellation(
			context.Background(),
			"missing", big.NewInt(2),
			"resource-eny", big.NewInt(100), "", types.DvpEnygma,
			"resource-721", big.NewInt(1), "tok-out", types.DvpERC721,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "swap not found")
		assert.Empty(t, deps.dvpClient.CancelSwapCalls())
	})

	t.Run("returns wrapped error when agreement verification fails", func(t *testing.T) {
		deps := newDefaultDeps()
		swap := preInitiatedSwap("shared-mismatch", types.DvpEnygma, types.DvpERC721, "resource-eny", "resource-721")
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		want := errors.New("amounts disagree")
		deps.swapAgreement.VerifyFunc = func(ctx context.Context, _ *types.DvpSwap, _ *big.Int, _ string, _ *big.Int, _ string, _ types.DvpTokenType, _ string, _ *big.Int, _ string, _ types.DvpTokenType) (string, error) {
			return "amount-mismatch", want
		}
		init := newInitiator(deps)

		err := init.HandleSwapCancellation(
			context.Background(),
			"shared-mismatch", big.NewInt(2),
			"resource-eny", big.NewInt(101), "", types.DvpEnygma,
			"resource-721", big.NewInt(1), "tok-out", types.DvpERC721,
		)
		require.Error(t, err)
		require.ErrorIs(t, err, want)
		assert.Empty(t, deps.dvpClient.CancelSwapCalls())
	})

	t.Run("propagates dvpClient.CancelSwap error", func(t *testing.T) {
		deps := newDefaultDeps()
		swap := preInitiatedSwap("shared-r", types.DvpEnygma, types.DvpERC721, "resource-eny", "resource-721")
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		want := errors.New("revert tx failed")
		deps.dvpClient.CancelSwapFunc = func(ctx context.Context, sharedId string, preimage *big.Int) error {
			return want
		}
		init := newInitiator(deps)

		err := init.HandleSwapCancellation(
			context.Background(),
			"shared-r", big.NewInt(2),
			"resource-eny", big.NewInt(100), "", types.DvpEnygma,
			"resource-721", big.NewInt(1), "tok-out", types.DvpERC721,
		)
		require.Error(t, err)
		require.ErrorIs(t, err, want)
	})
}
