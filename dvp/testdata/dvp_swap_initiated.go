package testdata

import (
	"math/big"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

// SwapInitiatedOption mutates a *types.DvpSwapInitiatedData in place.
type SwapInitiatedOption func(*types.DvpSwapInitiatedData)

// NewDvpSwapInitiatedData returns the payload that the dest log parser
// pushes to the orchestrator after a successful SwapInitiated decrypt.
// Defaults: a Message built by NewDvpSwapMessage() and a non-nil
// InitiatorDestSalt; ExpiresAt is left nil (the v2 implementation does
// not consume it from the event).
func NewDvpSwapInitiatedData(opts ...SwapInitiatedOption) *types.DvpSwapInitiatedData {
	data := &types.DvpSwapInitiatedData{
		Message:           NewDvpSwapMessage(),
		InitiatorDestSalt: big.NewInt(0xBBBB),
	}
	for _, opt := range opts {
		opt(data)
	}
	return data
}

// WithMessage replaces the embedded *DvpSwapMessage.
func WithMessage(msg *types.DvpSwapMessage) SwapInitiatedOption {
	return func(d *types.DvpSwapInitiatedData) { d.Message = msg }
}

// WithInitiatorDestSalt sets DvpSwapInitiatedData.InitiatorDestSalt — the
// salt the responder recovered from the event ciphertext via
// cryptography.RecoverSalt.
func WithInitiatorDestSalt(salt *big.Int) SwapInitiatedOption {
	return func(d *types.DvpSwapInitiatedData) { d.InitiatorDestSalt = salt }
}

// WithExpiresAt sets DvpSwapInitiatedData.ExpiresAt (currently unused by the
// receiver but populated by the log parser when it lands).
func WithExpiresAt(expiresAt *big.Int) SwapInitiatedOption {
	return func(d *types.DvpSwapInitiatedData) { d.ExpiresAt = expiresAt }
}
