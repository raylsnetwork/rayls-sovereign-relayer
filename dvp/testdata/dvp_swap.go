package testdata

import (
	"math/big"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

// SharedID is the canonical default shared-id used by NewDvpSwap and
// related builders. Exported so tests can pass it to other call sites
// (e.g. contractclient.SharedID-style helpers).
const SharedID = "shared-test-1"

// SwapOption mutates a *types.DvpSwap in place. Every NewDvpSwap option
// must follow this signature.
type SwapOption func(*types.DvpSwap)

// NewDvpSwap returns a swap row populated with sane defaults; opts mutate
// individual fields. Defaults match the most common scenario: an
// already-initiated Enygma -> ERC721 swap with both salts populated.
func NewDvpSwap(opts ...SwapOption) *types.DvpSwap {
	swap := &types.DvpSwap{
		SharedID:           SharedID,
		From:               "0xAlice",
		To:                 "0xBob",
		SourceChainID:      big.NewInt(1),
		DestChainID:        big.NewInt(2),
		TokenInAmount:      big.NewInt(100),
		TokenInAddress:     "0xTokenIn",
		TokenInResourceID:  "resource-in",
		TokenInType:        types.DvpEnygma,
		TokenInID:          "",
		TokenOutAmount:     big.NewInt(200),
		TokenOutAddress:    "0xTokenOut",
		TokenOutResourceID: "resource-out",
		TokenOutType:       types.DvpERC721,
		TokenOutID:         "tok-out",
		Status:             types.DvpSwapInitiated,
		SelfSalt:           big.NewInt(0xAAAA),
		DestSalt:           big.NewInt(0xBBBB),
	}
	for _, opt := range opts {
		opt(swap)
	}
	return swap
}

// WithSharedID sets DvpSwap.SharedID.
func WithSharedID(id string) SwapOption {
	return func(s *types.DvpSwap) { s.SharedID = id }
}

// WithStatus sets DvpSwap.Status.
func WithStatus(status types.DvpSwapStatus) SwapOption {
	return func(s *types.DvpSwap) { s.Status = status }
}

// WithFromTo sets DvpSwap.From and DvpSwap.To.
func WithFromTo(from, to string) SwapOption {
	return func(s *types.DvpSwap) {
		s.From = from
		s.To = to
	}
}

// WithChains sets the SourceChainID and DestChainID.
func WithChains(source, dest *big.Int) SwapOption {
	return func(s *types.DvpSwap) {
		s.SourceChainID = source
		s.DestChainID = dest
	}
}

// WithSalts sets SelfSalt and DestSalt — the two load-bearing fields for
// the v2 swap-confirmation flow.
func WithSalts(self, dest *big.Int) SwapOption {
	return func(s *types.DvpSwap) {
		s.SelfSalt = self
		s.DestSalt = dest
	}
}

// WithLegs replaces the four token-leg fields in one call.
func WithLegs(
	tokenInType types.DvpTokenType, tokenInResourceID, tokenInAddress, tokenInID string, tokenInAmount *big.Int,
	tokenOutType types.DvpTokenType, tokenOutResourceID, tokenOutAddress, tokenOutID string, tokenOutAmount *big.Int,
) SwapOption {
	return func(s *types.DvpSwap) {
		s.TokenInType = tokenInType
		s.TokenInResourceID = tokenInResourceID
		s.TokenInAddress = tokenInAddress
		s.TokenInID = tokenInID
		s.TokenInAmount = tokenInAmount
		s.TokenOutType = tokenOutType
		s.TokenOutResourceID = tokenOutResourceID
		s.TokenOutAddress = tokenOutAddress
		s.TokenOutID = tokenOutID
		s.TokenOutAmount = tokenOutAmount
	}
}

// WithCreatedAt sets DvpSwap.CreatedAt.
func WithCreatedAt(t time.Time) SwapOption {
	return func(s *types.DvpSwap) { s.CreatedAt = t }
}

// NewMirroredSwap returns the swap row Bob persists when his log parser
// fires the SwapInitiated event. The salt mirror is the load-bearing rule
// from dvp/handler/receiver.go:441-475:
//
//	SelfSalt = data.InitiatorDestSalt        // salt Alice picked for Bob's commitment
//	DestSalt = data.Message.InitiatorSelfSalt // salt Alice picked for her own self-destination
//
// The token legs are also mirrored: msg.TokenIn (what Alice spends) becomes
// Bob's TokenOut (what Bob receives), and vice versa.
//
// Centralising this in one builder means future refactors that change the
// mirror semantics only need to update one place.
func NewMirroredSwap(data *types.DvpSwapInitiatedData, opts ...SwapOption) *types.DvpSwap {
	msg := data.Message
	swap := &types.DvpSwap{
		SharedID:           msg.SharedId,
		To:                 msg.To,
		SourceChainID:      big.NewInt(2), // local chain (Bob)
		DestChainID:        msg.ChainId,
		TokenInAmount:      msg.TokenOutAmount,
		TokenInAddress:     msg.TokenOutAddress,
		TokenInResourceID:  msg.TokenOutResourceID,
		TokenInType:        msg.TokenOutType,
		TokenInID:          msg.TokenOutID,
		TokenOutAmount:     msg.TokenInAmount,
		TokenOutAddress:    msg.TokenInAddress,
		TokenOutResourceID: msg.TokenInResourceID,
		TokenOutType:       msg.TokenInType,
		TokenOutID:         msg.TokenInID,
		Status:             types.DvpSwapWaitingConfirmation,
		SelfSalt:           data.InitiatorDestSalt,
		DestSalt:           msg.InitiatorSelfSalt,
	}
	for _, opt := range opts {
		opt(swap)
	}
	return swap
}
