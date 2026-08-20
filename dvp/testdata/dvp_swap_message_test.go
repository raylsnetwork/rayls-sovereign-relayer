package testdata

import (
	"math/big"
	"testing"

	"github.com/raylsnetwork/rayls-sovereign-relayer/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDvpSwapMessage_Defaults(t *testing.T) {
	msg := NewDvpSwapMessage()
	require.NotNil(t, msg)
	assert.Equal(t, SharedID, msg.SharedId)
	assert.Equal(t, "0xBob", msg.To)
	assert.Equal(t, big.NewInt(7), msg.ChainId)
	assert.Equal(t, types.DvpEnygma, msg.TokenInType)
	assert.Equal(t, types.DvpERC721, msg.TokenOutType)
	assert.Equal(t, big.NewInt(0xAAAA), msg.InitiatorSelfSalt)
}

func TestNewDvpSwapMessage_Options(t *testing.T) {
	t.Run("WithMessageSharedID", func(t *testing.T) {
		assert.Equal(t, "id-1", NewDvpSwapMessage(WithMessageSharedID("id-1")).SharedId)
	})
	t.Run("WithMessageTo", func(t *testing.T) {
		assert.Equal(t, "0xCarol", NewDvpSwapMessage(WithMessageTo("0xCarol")).To)
	})
	t.Run("WithMessageChainID", func(t *testing.T) {
		assert.Equal(t, big.NewInt(99), NewDvpSwapMessage(WithMessageChainID(big.NewInt(99))).ChainId)
	})
	t.Run("WithInitiatorSelfSalt", func(t *testing.T) {
		assert.Equal(t, big.NewInt(0xABCD), NewDvpSwapMessage(WithInitiatorSelfSalt(big.NewInt(0xABCD))).InitiatorSelfSalt)
	})

	t.Run("WithMessageLegs swaps token-leg fields", func(t *testing.T) {
		msg := NewDvpSwapMessage(WithMessageLegs(
			types.DvpERC1155, "rin", "0xIn", "5", big.NewInt(50),
			types.DvpEnygma, "rout", "0xOut", "", big.NewInt(75),
		))
		assert.Equal(t, types.DvpERC1155, msg.TokenInType)
		assert.Equal(t, "5", msg.TokenInID)
		assert.Equal(t, big.NewInt(50), msg.TokenInAmount)
		assert.Equal(t, types.DvpEnygma, msg.TokenOutType)
		assert.Equal(t, big.NewInt(75), msg.TokenOutAmount)
	})
}
