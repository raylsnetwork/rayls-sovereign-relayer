package testdata

import (
	"math/big"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

// MessageOption mutates a *types.DvpSwapMessage in place.
type MessageOption func(*types.DvpSwapMessage)

// NewDvpSwapMessage returns a message populated with sane defaults
// matching the most common scenario: an Enygma -> ERC721 swap with
// InitiatorSelfSalt set. Override fields with the With* options.
func NewDvpSwapMessage(opts ...MessageOption) *types.DvpSwapMessage {
	msg := &types.DvpSwapMessage{
		SharedId:           SharedID,
		To:                 "0xBob",
		ChainId:            big.NewInt(7),
		PNTxHash:           "0xPNTx",
		PNTxTimestamp:      time.Date(2026, 4, 16, 12, 0, 0, 0, time.UTC),
		TokenInAmount:      big.NewInt(100),
		TokenInAddress:     "0xTokenIn",
		TokenInResourceID:  "resource-in",
		TokenInType:        types.DvpEnygma,
		TokenInID:          "",
		TokenOutAmount:     big.NewInt(1),
		TokenOutAddress:    "0xTokenOut",
		TokenOutResourceID: "resource-out",
		TokenOutType:       types.DvpERC721,
		TokenOutID:         "tok-out",
		InitiatorSelfSalt:  big.NewInt(0xAAAA),
	}
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}

// WithMessageSharedID sets DvpSwapMessage.SharedId.
func WithMessageSharedID(id string) MessageOption {
	return func(m *types.DvpSwapMessage) { m.SharedId = id }
}

// WithMessageTo sets DvpSwapMessage.To.
func WithMessageTo(to string) MessageOption {
	return func(m *types.DvpSwapMessage) { m.To = to }
}

// WithMessageChainID sets DvpSwapMessage.ChainId.
func WithMessageChainID(chainID *big.Int) MessageOption {
	return func(m *types.DvpSwapMessage) { m.ChainId = chainID }
}

// WithMessageLegs replaces the eight token-leg fields in one call (mirrors WithLegs on DvpSwap).
func WithMessageLegs(
	tokenInType types.DvpTokenType, tokenInResourceID, tokenInAddress, tokenInID string, tokenInAmount *big.Int,
	tokenOutType types.DvpTokenType, tokenOutResourceID, tokenOutAddress, tokenOutID string, tokenOutAmount *big.Int,
) MessageOption {
	return func(m *types.DvpSwapMessage) {
		m.TokenInType = tokenInType
		m.TokenInResourceID = tokenInResourceID
		m.TokenInAddress = tokenInAddress
		m.TokenInID = tokenInID
		m.TokenInAmount = tokenInAmount
		m.TokenOutType = tokenOutType
		m.TokenOutResourceID = tokenOutResourceID
		m.TokenOutAddress = tokenOutAddress
		m.TokenOutID = tokenOutID
		m.TokenOutAmount = tokenOutAmount
	}
}

// WithInitiatorSelfSalt sets DvpSwapMessage.InitiatorSelfSalt — the salt
// the initiator used to lock its own self-destination commitment.
func WithInitiatorSelfSalt(salt *big.Int) MessageOption {
	return func(m *types.DvpSwapMessage) { m.InitiatorSelfSalt = salt }
}
