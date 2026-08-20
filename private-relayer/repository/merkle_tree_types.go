package repository

import (
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

const MerkleTreeCollectionName = "merkle_trees"

type MerkleTree struct {
	Type         types.DvpTokenType `db:"type"`
	TokenAddress string             `db:"token_address"`
	Number       int                `db:"number"`
	Depth        int                `db:"depth"`
	Leaves       []string           `db:"leaves"`
	CreatedAt    time.Time          `db:"created_at"`
}
