package repository

import (
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

const DvpDepositsCollectionName = "dvp_deposits"

type DvpDeposit struct {
	UserAddress  string                 `db:"user_address"`
	Salt         string                 `db:"salt"`
	TokenAmount  string                 `db:"token_amount"`
	TokenAddress string                 `db:"token_address"`
	TokenType    types.DvpTokenType     `db:"token_type"`
	TokenID      string                 `db:"token_id"`
	TreeNumber   int                    `db:"tree_number"`
	Commitment   string                 `db:"commitment"`
	Nullifier    string                 `db:"nullifier"`
	Status       types.DvpDepositStatus `db:"status"`
	CreatedAt    time.Time              `db:"created_at"`
}
