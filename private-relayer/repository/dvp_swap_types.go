package repository

import (
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

const DvpSwapCollectionName = "dvp_swaps"

type DvpSwap struct {
	SharedID      string `db:"shared_id"`
	From          string `db:"from_address"`
	To            string `db:"to_address"`
	SourceChainID string `db:"source_chain_id"`
	DestChainID   string `db:"dest_chain_id"`
	TokenInAmount          string              `db:"token_in_amount"`
	TokenInAddress         string              `db:"token_in_address"`
	TokenInResourceID      string              `db:"token_in_resource_id"`
	TokenInType            types.DvpTokenType  `db:"token_in_type"`
	TokenInID              string              `db:"token_in_id"`
	TokenOutAmount         string              `db:"token_out_amount"`
	TokenOutAddress        string              `db:"token_out_address"`
	TokenOutResourceID     string              `db:"token_out_resource_id"`
	TokenOutType           types.DvpTokenType  `db:"token_out_type"`
	TokenOutID             string              `db:"token_out_id"`
	Status                 types.DvpSwapStatus `db:"status"`
	CreatedAt              time.Time           `db:"created_at"`
	CancelledAt            *time.Time          `db:"cancelled_at"`
	SelfSalt               *string             `db:"self_salt"`
	DestSalt               *string             `db:"dest_salt"`
	CancelPreimage         *string             `db:"cancel_preimage"`
}
