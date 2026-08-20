package handler_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/dvp/handler"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// receiverDeps is a small struct that bundles every collaborator the
// DvpReceiver constructor needs. Callers tweak the specific fields
// relevant to their scenario before invoking newReceiverWithDeps.
type receiverDeps struct {
	swapRepo       *receiverSwapRepositoryMock
	depositRepo    *receiverDepositRepositoryMock
	merkleClient   *receiverMerkleClientMock
	endpointClient *receiverEndpointClientMock
	enygmaHandler  *receiverEnygmaHandlerClientMock
	erc721Handler  *receiverERC721HandlerClientMock
	erc1155Handler *receiverERC1155HandlerClientMock
	depositFinder  *receiverDepositFinderMock
	consolidator   *receiverDepositConsolidatorMock
	commitmentCalc *receiverCommitmentCalculatorMock
	proofGen       *receiverProofGeneratorMock
	retry          *receiverRetryServiceMock
	txManager      *receiverTxManagerMock
}

func newDefaultReceiverDeps() *receiverDeps {
	return &receiverDeps{
		swapRepo: &receiverSwapRepositoryMock{
			CreateSwapFunc: func(ctx context.Context, swap *types.DvpSwap) error { return nil },
			UpdateSwapStatusFunc: func(ctx context.Context, sharedID string, status types.DvpSwapStatus) error {
				return nil
			},
		},
		depositRepo: &receiverDepositRepositoryMock{
			ConfirmDepositFunc: func(ctx context.Context, commitment *big.Int, treeNumber *big.Int) error {
				return nil
			},
			GetDepositByCommitmentAndStatusFunc: func(ctx context.Context, commitment *big.Int, status types.DvpDepositStatus) (*types.DvpDeposit, error) {
				return nil, nil
			},
			UpdateDepositStatusByNullifierFunc: func(ctx context.Context, tokenAddress string, nullifier *big.Int, status types.DvpDepositStatus) error {
				return nil
			},
		},
		merkleClient: &receiverMerkleClientMock{
			PopulateMerkleDbTreeFunc: func(ctx context.Context, tokenAddress string, bigIntTokenType *big.Int, bigIntTreeNumber *big.Int, bigIntLeaves []*big.Int) error {
				return nil
			},
		},
		endpointClient: &receiverEndpointClientMock{
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return common.HexToAddress("0xResource"), nil
			},
		},
		enygmaHandler:  &receiverEnygmaHandlerClientMock{},
		erc721Handler:  &receiverERC721HandlerClientMock{},
		erc1155Handler: &receiverERC1155HandlerClientMock{},
		depositFinder:  &receiverDepositFinderMock{},
		consolidator:   &receiverDepositConsolidatorMock{},
		commitmentCalc: &receiverCommitmentCalculatorMock{},
		proofGen:       &receiverProofGeneratorMock{},
		retry:          &receiverRetryServiceMock{},
		txManager:      &receiverTxManagerMock{},
	}
}

func newReceiverWithDeps(deps *receiverDeps) *handler.DvpReceiver {
	return handler.NewDvpReceiver(
		handler.ReceiverConfig{ChainID: big.NewInt(2)},
		deps.swapRepo,
		deps.depositRepo,
		deps.endpointClient,
		deps.enygmaHandler,
		deps.erc721Handler,
		deps.erc1155Handler,
		deps.depositFinder,
		deps.consolidator,
		deps.commitmentCalc,
		deps.proofGen,
		deps.merkleClient,
		deps.retry,
		deps.txManager,
	)
}

func TestNewDvpReceiver(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("constructs with all v2 dependencies", func(t *testing.T) {
		r := newReceiverWithDeps(newDefaultReceiverDeps())
		require.NotNil(t, r)
	})
}

func TestReceiver_HandleSwapInitiated(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("persists mirrored swap with swapped salt semantics", func(t *testing.T) {
		// salt the initiator used to lock its own self-destination commitment
		// (the token Alice will eventually receive).
		// From Bob's perspective this is the salt for the destination commitment
		// (the token Bob is sending), so Bob stores it as DestSalt.
		initiatorSelfSalt := big.NewInt(0xAAAA)
		// salt the initiator used for Bob's commitment. Bob recovered this from
		// the SwapInitiated event ciphertext via cryptography.RecoverSalt and
		// the log-parser passes it through as InitiatorDestSalt. Bob stores it
		// as SelfSalt (it is the salt Bob will use for its own commitment).
		initiatorDestSalt := big.NewInt(0xBBBB)

		data := &types.DvpSwapInitiatedData{
			InitiatorDestSalt: initiatorDestSalt,
			Message: &types.DvpSwapMessage{
				SharedId: "shared-1",
				To:       "0xBob",
				ChainId:  big.NewInt(7),
				// From Bob's perspective:
				// - msg.TokenIn is the asset Alice deposited (Bob will RECEIVE) -> Bob's TokenOut
				// - msg.TokenOut is what Alice wants (Bob will SEND)             -> Bob's TokenIn
				TokenInAmount:      big.NewInt(100),
				TokenInAddress:     "0xAlicesAsset",
				TokenInResourceID:  "resource-alices-asset",
				TokenInType:        types.DvpEnygma,
				TokenInID:          "",
				TokenOutAmount:     big.NewInt(1),
				TokenOutAddress:    "0xBobsAsset",
				TokenOutResourceID: "resource-bobs-asset",
				TokenOutType:       types.DvpERC721,
				TokenOutID:         "tok-7",
				InitiatorSelfSalt:  initiatorSelfSalt,
			},
		}

		deps := newDefaultReceiverDeps()
		r := newReceiverWithDeps(deps)

		err := r.HandleSwapInitiated(context.Background(), 100, data)
		require.NoError(t, err)

		require.Len(t, deps.swapRepo.CreateSwapCalls(), 1)
		got := deps.swapRepo.CreateSwapCalls()[0].Swap
		assert.Equal(t, "shared-1", got.SharedID)
		assert.Equal(t, "0xBob", got.To)
		assert.Equal(t, types.DvpSwapWaitingConfirmation, got.Status)
		assert.Equal(t, big.NewInt(2), got.SourceChainID, "Bob's SourceChainID is the receiver's local chain")
		assert.Equal(t, big.NewInt(7), got.DestChainID, "DestChainID is the initiator's chain (msg.ChainId)")

		// Salt mirror invariant — the load-bearing rule from receiver.go:441-475.
		assert.Equal(t, initiatorDestSalt, got.SelfSalt,
			"SelfSalt must equal InitiatorDestSalt — Bob uses the salt the initiator created for him")
		assert.Equal(t, initiatorSelfSalt, got.DestSalt,
			"DestSalt must equal Message.InitiatorSelfSalt — the salt the initiator used for its own self-destination commitment")

		// Token-leg mirror — initiator's TokenOut is Bob's TokenIn.
		assert.Equal(t, "0xBobsAsset", got.TokenInAddress)
		assert.Equal(t, types.DvpERC721, got.TokenInType)
		assert.Equal(t, "tok-7", got.TokenInID)
		assert.Equal(t, big.NewInt(1), got.TokenInAmount)
		assert.Equal(t, "0xAlicesAsset", got.TokenOutAddress)
		assert.Equal(t, types.DvpEnygma, got.TokenOutType)
		assert.Equal(t, big.NewInt(100), got.TokenOutAmount)
	})

	t.Run("returns error when CreateSwap fails", func(t *testing.T) {
		deps := newDefaultReceiverDeps()
		want := errors.New("db down")
		deps.swapRepo.CreateSwapFunc = func(ctx context.Context, swap *types.DvpSwap) error {
			return want
		}
		r := newReceiverWithDeps(deps)

		data := &types.DvpSwapInitiatedData{
			InitiatorDestSalt: big.NewInt(1),
			Message: &types.DvpSwapMessage{
				SharedId:          "shared-2",
				ChainId:           big.NewInt(7),
				InitiatorSelfSalt: big.NewInt(2),
			},
		}
		err := r.HandleSwapInitiated(context.Background(), 100, data)
		require.ErrorIs(t, err, want)
	})
}

func TestReceiver_HandleSwapCompleted(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("returns error when swap not found", func(t *testing.T) {
		deps := newDefaultReceiverDeps()
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, sharedID string) (*types.DvpSwap, error) {
			return nil, nil
		}
		r := newReceiverWithDeps(deps)
		err := r.HandleSwapCompleted(context.Background(), "missing")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("Enygma swap: updates status and notifies enygma handler", func(t *testing.T) {
		deps := newDefaultReceiverDeps()
		swap := &types.DvpSwap{
			SharedID:          "shared-eny",
			From:              "0xAlice",
			To:                "0xBob",
			TokenInResourceID: "resource-eny",
			TokenInType:       types.DvpEnygma,
			SourceChainID:     big.NewInt(1),
			DestChainID:       big.NewInt(2),
		}
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		deps.swapRepo.UpdateSwapStatusFunc = func(ctx context.Context, sharedID string, status types.DvpSwapStatus) error {
			return nil
		}
		deps.endpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.Address{}, nil
		}
		deps.enygmaHandler.MarkSwapCompletedFunc = func(ctx context.Context, tokenAddress common.Address, chainID *big.Int, sharedId string) error {
			return nil
		}
		r := newReceiverWithDeps(deps)

		err := r.HandleSwapCompleted(context.Background(), "shared-eny")
		require.NoError(t, err)

		require.Len(t, deps.swapRepo.UpdateSwapStatusCalls(), 1)
		assert.Equal(t, types.DvpSwapCompleted, deps.swapRepo.UpdateSwapStatusCalls()[0].Status)
		require.Len(t, deps.enygmaHandler.MarkSwapCompletedCalls(), 1)
		assert.Equal(t, "shared-eny", deps.enygmaHandler.MarkSwapCompletedCalls()[0].SharedId)
		assert.Equal(t, big.NewInt(2), deps.enygmaHandler.MarkSwapCompletedCalls()[0].ChainID)
	})

	t.Run("ERC721 swap: invokes ERC721 MarkSwapCompleted", func(t *testing.T) {
		deps := newDefaultReceiverDeps()
		swap := &types.DvpSwap{
			SharedID:          "shared-721",
			From:              "0xAlice",
			To:                "0xBob",
			TokenInResourceID: "resource-721",
			TokenInType:       types.DvpERC721,
			TokenInID:         "tok-1",
			SourceChainID:     big.NewInt(1),
			DestChainID:       big.NewInt(2),
		}
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		deps.swapRepo.UpdateSwapStatusFunc = func(ctx context.Context, sharedID string, status types.DvpSwapStatus) error {
			return nil
		}
		deps.endpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.Address{}, nil
		}
		deps.erc721Handler.MarkSwapCompletedFunc = func(ctx context.Context, tokenAddress common.Address, to common.Address, sourceChainId *big.Int, destChainId *big.Int, resourceId string, tokenId string, amount *big.Int, sharedId string) error {
			return nil
		}
		r := newReceiverWithDeps(deps)

		err := r.HandleSwapCompleted(context.Background(), "shared-721")
		require.NoError(t, err)

		require.Len(t, deps.erc721Handler.MarkSwapCompletedCalls(), 1)
		assert.Equal(t, "shared-721", deps.erc721Handler.MarkSwapCompletedCalls()[0].SharedId)
		assert.Equal(t, common.HexToAddress("0xBob"), deps.erc721Handler.MarkSwapCompletedCalls()[0].To)
	})

	t.Run("ERC1155 swap: invokes ERC1155 MarkSwapCompleted", func(t *testing.T) {
		deps := newDefaultReceiverDeps()
		swap := &types.DvpSwap{
			SharedID:          "shared-1155",
			From:              "0xAlice",
			To:                "0xBob",
			TokenInResourceID: "resource-1155",
			TokenInType:       types.DvpERC1155,
			TokenInID:         "tok-7",
			TokenInAmount:     big.NewInt(42),
			SourceChainID:     big.NewInt(1),
			DestChainID:       big.NewInt(2),
		}
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		deps.swapRepo.UpdateSwapStatusFunc = func(ctx context.Context, sharedID string, status types.DvpSwapStatus) error {
			return nil
		}
		deps.endpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.Address{}, nil
		}
		deps.erc1155Handler.MarkSwapCompletedFunc = func(ctx context.Context, tokenAddress common.Address, from common.Address, to common.Address, sourceChainId *big.Int, destChainId *big.Int, resourceId string, tokenId string, amount *big.Int, data []byte, sharedId string) error {
			return nil
		}
		r := newReceiverWithDeps(deps)

		err := r.HandleSwapCompleted(context.Background(), "shared-1155")
		require.NoError(t, err)

		require.Len(t, deps.erc1155Handler.MarkSwapCompletedCalls(), 1)
		assert.Equal(t, "shared-1155", deps.erc1155Handler.MarkSwapCompletedCalls()[0].SharedId)
		assert.Equal(t, common.HexToAddress("0xBob"), deps.erc1155Handler.MarkSwapCompletedCalls()[0].To)
		assert.Equal(t, common.HexToAddress("0xAlice"), deps.erc1155Handler.MarkSwapCompletedCalls()[0].From)
	})

	t.Run("returns error on unsupported token type", func(t *testing.T) {
		deps := newDefaultReceiverDeps()
		swap := &types.DvpSwap{
			SharedID:    "shared-bad",
			TokenInType: types.DvpCustom,
		}
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		r := newReceiverWithDeps(deps)
		err := r.HandleSwapCompleted(context.Background(), "shared-bad")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Unsupported token type")
	})
}

func TestReceiver_HandleSwapRevert(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("rejects unsupported revert status", func(t *testing.T) {
		r := newReceiverWithDeps(newDefaultReceiverDeps())
		err := r.HandleSwapRevert(context.Background(), "x", types.DvpSwapInitiated)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid revert status")
	})

	t.Run("returns error when swap not found", func(t *testing.T) {
		deps := newDefaultReceiverDeps()
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return nil, nil
		}
		r := newReceiverWithDeps(deps)
		err := r.HandleSwapRevert(context.Background(), "x", types.DvpSwapCancelled)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "swap not found")
	})

	t.Run("Cancelled status: cancels swap and notifies cancellation message via enygma handler", func(t *testing.T) {
		deps := newDefaultReceiverDeps()
		swap := &types.DvpSwap{
			SharedID:          "shared-eny",
			TokenInResourceID: "resource-eny",
			TokenInType:       types.DvpEnygma,
		}
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		deps.swapRepo.CancelSwapFunc = func(ctx context.Context, sharedID string, status types.DvpSwapStatus) error {
			return nil
		}
		deps.endpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.Address{}, nil
		}
		deps.enygmaHandler.NotifySenderWithPNCommunicatorFunc = func(ctx context.Context, tokenAddress common.Address, sharedId string, status types.DvpCommunicatorStatus, message string) error {
			return nil
		}
		r := newReceiverWithDeps(deps)

		err := r.HandleSwapRevert(context.Background(), "shared-eny", types.DvpSwapCancelled)
		require.NoError(t, err)

		require.Len(t, deps.swapRepo.CancelSwapCalls(), 1)
		assert.Equal(t, types.DvpSwapCancelled, deps.swapRepo.CancelSwapCalls()[0].Status)
		require.Len(t, deps.enygmaHandler.NotifySenderWithPNCommunicatorCalls(), 1)
		assert.Equal(t, types.SwapCancelled, deps.enygmaHandler.NotifySenderWithPNCommunicatorCalls()[0].Status)
		assert.Equal(t, "swap is cancelled", deps.enygmaHandler.NotifySenderWithPNCommunicatorCalls()[0].Message)
	})

	t.Run("TimedOut status: cancels swap and notifies timeout message via ERC721 handler", func(t *testing.T) {
		deps := newDefaultReceiverDeps()
		swap := &types.DvpSwap{
			SharedID:          "shared-721",
			TokenInResourceID: "resource-721",
			TokenInType:       types.DvpERC721,
		}
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		deps.swapRepo.CancelSwapFunc = func(ctx context.Context, sharedID string, status types.DvpSwapStatus) error {
			return nil
		}
		deps.endpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.Address{}, nil
		}
		deps.erc721Handler.NotifySenderWithPNCommunicatorFunc = func(ctx context.Context, tokenAddress common.Address, sharedId string, status types.DvpCommunicatorStatus, message string) error {
			return nil
		}
		r := newReceiverWithDeps(deps)

		err := r.HandleSwapRevert(context.Background(), "shared-721", types.DvpSwapTimedOut)
		require.NoError(t, err)

		require.Len(t, deps.swapRepo.CancelSwapCalls(), 1)
		assert.Equal(t, types.DvpSwapTimedOut, deps.swapRepo.CancelSwapCalls()[0].Status)
		require.Len(t, deps.erc721Handler.NotifySenderWithPNCommunicatorCalls(), 1)
		assert.Equal(t, types.SwapTimedOut, deps.erc721Handler.NotifySenderWithPNCommunicatorCalls()[0].Status)
		assert.Equal(t, "swap is expired", deps.erc721Handler.NotifySenderWithPNCommunicatorCalls()[0].Message)
	})

	t.Run("TimedOut status: notifies via ERC1155 handler when token type is ERC1155", func(t *testing.T) {
		deps := newDefaultReceiverDeps()
		swap := &types.DvpSwap{
			SharedID:          "shared-1155",
			TokenInResourceID: "resource-1155",
			TokenInType:       types.DvpERC1155,
		}
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		deps.swapRepo.CancelSwapFunc = func(ctx context.Context, sharedID string, status types.DvpSwapStatus) error {
			return nil
		}
		deps.endpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.Address{}, nil
		}
		deps.erc1155Handler.NotifySenderWithPNCommunicatorFunc = func(ctx context.Context, tokenAddress common.Address, sharedId string, status types.DvpCommunicatorStatus, message string) error {
			return nil
		}
		r := newReceiverWithDeps(deps)

		err := r.HandleSwapRevert(context.Background(), "shared-1155", types.DvpSwapTimedOut)
		require.NoError(t, err)
		require.Len(t, deps.erc1155Handler.NotifySenderWithPNCommunicatorCalls(), 1)
		assert.Equal(t, types.SwapTimedOut, deps.erc1155Handler.NotifySenderWithPNCommunicatorCalls()[0].Status)
	})

	t.Run("propagates CancelSwap error", func(t *testing.T) {
		deps := newDefaultReceiverDeps()
		want := errors.New("db down")
		swap := &types.DvpSwap{SharedID: "x", TokenInType: types.DvpEnygma}
		deps.swapRepo.GetSwapBySharedIDFunc = func(ctx context.Context, _ string) (*types.DvpSwap, error) {
			return swap, nil
		}
		deps.swapRepo.CancelSwapFunc = func(ctx context.Context, sharedID string, status types.DvpSwapStatus) error {
			return want
		}
		r := newReceiverWithDeps(deps)
		err := r.HandleSwapRevert(context.Background(), "x", types.DvpSwapCancelled)
		require.ErrorIs(t, err, want)
	})
}

func TestReceiver_HandleCommitments(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("populates merkle tree and confirms each known pending deposit", func(t *testing.T) {
		c1, c2, c3 := big.NewInt(1), big.NewInt(2), big.NewInt(3)
		commits := []*big.Int{c1, c2, c3}

		deps := newDefaultReceiverDeps()
		// Only c1 and c3 correspond to known pending deposits; c2 is unknown
		// and should not trigger ConfirmDeposit.
		deps.depositRepo.GetDepositByCommitmentAndStatusFunc = func(ctx context.Context, c *big.Int, s types.DvpDepositStatus) (*types.DvpDeposit, error) {
			if c.Cmp(c2) == 0 {
				return nil, nil
			}
			return &types.DvpDeposit{Commitment: c}, nil
		}

		r := newReceiverWithDeps(deps)

		data := &types.DvpCommitmentsData{
			TokenAddress: "0xToken",
			TokenType:    big.NewInt(int64(types.DvpERC721)),
			TreeNumber:   big.NewInt(7),
			Commitments:  commits,
		}
		err := r.HandleCommitments(context.Background(), data)
		require.NoError(t, err)

		require.Len(t, deps.merkleClient.PopulateMerkleDbTreeCalls(), 1)
		require.Len(t, deps.depositRepo.ConfirmDepositCalls(), 2,
			"only the two known pending deposits get confirmed; the unknown commitment is skipped")
		confirmedCommitments := []*big.Int{
			deps.depositRepo.ConfirmDepositCalls()[0].Commitment,
			deps.depositRepo.ConfirmDepositCalls()[1].Commitment,
		}
		assert.ElementsMatch(t, []*big.Int{c1, c3}, confirmedCommitments)
	})

	t.Run("propagates merkle population error", func(t *testing.T) {
		want := errors.New("merkle insert failed")
		deps := newDefaultReceiverDeps()
		deps.merkleClient.PopulateMerkleDbTreeFunc = func(ctx context.Context, tokenAddress string, _, _ *big.Int, _ []*big.Int) error {
			return want
		}
		r := newReceiverWithDeps(deps)
		err := r.HandleCommitments(context.Background(), &types.DvpCommitmentsData{
			TokenType:  big.NewInt(0),
			TreeNumber: big.NewInt(0),
		})
		require.ErrorIs(t, err, want)
	})
}

func TestReceiver_HandleNullifiers(t *testing.T) {
	testtools.SilenceLogger()

	const tokenAddress = "0xToken"

	t.Run("marks deposits spent for each nullifier scoped by token", func(t *testing.T) {
		deps := newDefaultReceiverDeps()
		r := newReceiverWithDeps(deps)
		nfs := []*big.Int{big.NewInt(0xCAFE), big.NewInt(0xBEEF)}

		err := r.HandleNullifiers(context.Background(), tokenAddress, nfs)
		require.NoError(t, err)
		require.Len(t, deps.depositRepo.UpdateDepositStatusByNullifierCalls(), 2)
		assert.Equal(t, tokenAddress, deps.depositRepo.UpdateDepositStatusByNullifierCalls()[0].TokenAddress)
		assert.Equal(t, nfs[0], deps.depositRepo.UpdateDepositStatusByNullifierCalls()[0].Nullifier)
		assert.Equal(t, types.DvpDepositSpent, deps.depositRepo.UpdateDepositStatusByNullifierCalls()[0].Status)
		assert.Equal(t, tokenAddress, deps.depositRepo.UpdateDepositStatusByNullifierCalls()[1].TokenAddress)
		assert.Equal(t, nfs[1], deps.depositRepo.UpdateDepositStatusByNullifierCalls()[1].Nullifier)
		assert.Equal(t, types.DvpDepositSpent, deps.depositRepo.UpdateDepositStatusByNullifierCalls()[1].Status)
	})

	t.Run("propagates repository error", func(t *testing.T) {
		want := errors.New("update failed")
		deps := newDefaultReceiverDeps()
		deps.depositRepo.UpdateDepositStatusByNullifierFunc = func(ctx context.Context, ta string, n *big.Int, s types.DvpDepositStatus) error {
			return want
		}
		r := newReceiverWithDeps(deps)
		err := r.HandleNullifiers(context.Background(), tokenAddress, []*big.Int{big.NewInt(1)})
		require.ErrorIs(t, err, want)
	})
}
