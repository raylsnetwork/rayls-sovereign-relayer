package testdata

import (
	"math/big"
	"testing"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDvpSwap_Defaults(t *testing.T) {
	swap := NewDvpSwap()
	require.NotNil(t, swap)
	assert.Equal(t, SharedID, swap.SharedID)
	assert.Equal(t, types.DvpSwapInitiated, swap.Status)
	assert.Equal(t, types.DvpEnygma, swap.TokenInType)
	assert.Equal(t, types.DvpERC721, swap.TokenOutType)
	assert.Equal(t, big.NewInt(0xAAAA), swap.SelfSalt, "default SelfSalt is the canonical 0xAAAA")
	assert.Equal(t, big.NewInt(0xBBBB), swap.DestSalt, "default DestSalt is the canonical 0xBBBB")
}

func TestNewDvpSwap_Options(t *testing.T) {
	t.Run("WithSharedID", func(t *testing.T) {
		swap := NewDvpSwap(WithSharedID("custom"))
		assert.Equal(t, "custom", swap.SharedID)
	})

	t.Run("WithStatus", func(t *testing.T) {
		swap := NewDvpSwap(WithStatus(types.DvpSwapWaitingConfirmation))
		assert.Equal(t, types.DvpSwapWaitingConfirmation, swap.Status)
	})

	t.Run("WithFromTo", func(t *testing.T) {
		swap := NewDvpSwap(WithFromTo("0xCarol", "0xDave"))
		assert.Equal(t, "0xCarol", swap.From)
		assert.Equal(t, "0xDave", swap.To)
	})

	t.Run("WithChains", func(t *testing.T) {
		swap := NewDvpSwap(WithChains(big.NewInt(7), big.NewInt(11)))
		assert.Equal(t, big.NewInt(7), swap.SourceChainID)
		assert.Equal(t, big.NewInt(11), swap.DestChainID)
	})

	t.Run("WithSalts", func(t *testing.T) {
		swap := NewDvpSwap(WithSalts(big.NewInt(0xCC), big.NewInt(0xDD)))
		assert.Equal(t, big.NewInt(0xCC), swap.SelfSalt)
		assert.Equal(t, big.NewInt(0xDD), swap.DestSalt)
	})

	t.Run("WithLegs", func(t *testing.T) {
		swap := NewDvpSwap(WithLegs(
			types.DvpERC1155, "rin", "0xIn", "5", big.NewInt(50),
			types.DvpEnygma, "rout", "0xOut", "", big.NewInt(75),
		))
		assert.Equal(t, types.DvpERC1155, swap.TokenInType)
		assert.Equal(t, "5", swap.TokenInID)
		assert.Equal(t, big.NewInt(50), swap.TokenInAmount)
		assert.Equal(t, types.DvpEnygma, swap.TokenOutType)
		assert.Equal(t, big.NewInt(75), swap.TokenOutAmount)
	})

	t.Run("WithCreatedAt", func(t *testing.T) {
		ts := time.Date(2026, 4, 16, 10, 0, 0, 0, time.UTC)
		swap := NewDvpSwap(WithCreatedAt(ts))
		assert.Equal(t, ts, swap.CreatedAt)
	})

	t.Run("multiple options compose left-to-right", func(t *testing.T) {
		swap := NewDvpSwap(
			WithSharedID("compose-1"),
			WithStatus(types.DvpSwapCompleted),
			WithChains(big.NewInt(3), big.NewInt(4)),
		)
		assert.Equal(t, "compose-1", swap.SharedID)
		assert.Equal(t, types.DvpSwapCompleted, swap.Status)
		assert.Equal(t, big.NewInt(3), swap.SourceChainID)
		assert.Equal(t, big.NewInt(4), swap.DestChainID)
	})
}

func TestNewMirroredSwap_PreservesSaltMirror(t *testing.T) {
	t.Run("salt mirror invariant from receiver.go:441-475", func(t *testing.T) {
		initiatorSelfSalt := big.NewInt(0x1111)
		initiatorDestSalt := big.NewInt(0x2222)
		data := &types.DvpSwapInitiatedData{
			InitiatorDestSalt: initiatorDestSalt,
			Message: &types.DvpSwapMessage{
				SharedId:          "mirror-1",
				To:                "0xBob",
				ChainId:           big.NewInt(7),
				TokenInAmount:     big.NewInt(100),
				TokenInAddress:    "0xAlicesAsset",
				TokenInType:       types.DvpEnygma,
				TokenOutAmount:    big.NewInt(1),
				TokenOutAddress:   "0xBobsAsset",
				TokenOutType:      types.DvpERC721,
				TokenOutID:        "tok-7",
				InitiatorSelfSalt: initiatorSelfSalt,
			},
		}
		swap := NewMirroredSwap(data)
		require.NotNil(t, swap)
		assert.Equal(t, types.DvpSwapWaitingConfirmation, swap.Status)
		// Salt mirror: SelfSalt == data.InitiatorDestSalt (Alice's salt for Bob's commitment).
		assert.Equal(t, initiatorDestSalt, swap.SelfSalt,
			"SelfSalt must mirror InitiatorDestSalt — Bob uses the salt Alice generated for him")
		assert.Equal(t, initiatorSelfSalt, swap.DestSalt,
			"DestSalt must mirror Message.InitiatorSelfSalt — Alice's own self-destination salt")
		// Leg mirror: msg.TokenOut becomes Bob's TokenIn.
		assert.Equal(t, "0xBobsAsset", swap.TokenInAddress)
		assert.Equal(t, types.DvpERC721, swap.TokenInType)
		assert.Equal(t, big.NewInt(1), swap.TokenInAmount)
		assert.Equal(t, "0xAlicesAsset", swap.TokenOutAddress)
		assert.Equal(t, types.DvpEnygma, swap.TokenOutType)
		assert.Equal(t, big.NewInt(100), swap.TokenOutAmount)
	})

	t.Run("opts apply on top of the mirror", func(t *testing.T) {
		data := &types.DvpSwapInitiatedData{
			InitiatorDestSalt: big.NewInt(1),
			Message: &types.DvpSwapMessage{
				SharedId:          "mirror-2",
				ChainId:           big.NewInt(7),
				InitiatorSelfSalt: big.NewInt(2),
			},
		}
		swap := NewMirroredSwap(data, WithStatus(types.DvpSwapCancelled))
		assert.Equal(t, types.DvpSwapCancelled, swap.Status)
	})
}
