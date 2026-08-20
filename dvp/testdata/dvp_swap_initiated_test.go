package testdata

import (
	"math/big"
	"testing"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDvpSwapInitiatedData_Defaults(t *testing.T) {
	d := NewDvpSwapInitiatedData()
	require.NotNil(t, d)
	require.NotNil(t, d.Message, "default Message comes from NewDvpSwapMessage")
	assert.Equal(t, SharedID, d.Message.SharedId)
	assert.Equal(t, big.NewInt(0xBBBB), d.InitiatorDestSalt)
	assert.Nil(t, d.ExpiresAt, "ExpiresAt is currently unused by the receiver — leave nil by default")
}

func TestNewDvpSwapInitiatedData_Options(t *testing.T) {
	t.Run("WithMessage replaces the entire embedded message", func(t *testing.T) {
		custom := &types.DvpSwapMessage{SharedId: "custom-1", InitiatorSelfSalt: big.NewInt(0x1234)}
		d := NewDvpSwapInitiatedData(WithMessage(custom))
		assert.Same(t, custom, d.Message)
	})

	t.Run("WithInitiatorDestSalt", func(t *testing.T) {
		d := NewDvpSwapInitiatedData(WithInitiatorDestSalt(big.NewInt(0xCCCC)))
		assert.Equal(t, big.NewInt(0xCCCC), d.InitiatorDestSalt)
	})

	t.Run("WithExpiresAt", func(t *testing.T) {
		d := NewDvpSwapInitiatedData(WithExpiresAt(big.NewInt(1700000000)))
		assert.Equal(t, big.NewInt(1700000000), d.ExpiresAt)
	})
}
