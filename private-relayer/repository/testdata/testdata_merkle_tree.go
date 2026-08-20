package testdata

import (
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

var merkleTreeCreatedAt = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

var (
	MerkleTree1 = types.MerkleTree{
		Type:         types.DvpERC20,
		TokenAddress: "0xTokenAddressAAA",
		Number:       1,
		Depth:        16,
		Leaves:       []string{"leaf-a", "leaf-b"},
		CreatedAt:    merkleTreeCreatedAt,
	}
	MerkleTree2 = types.MerkleTree{
		Type:         types.DvpERC20,
		TokenAddress: "0xTokenAddressAAA",
		Number:       2,
		Depth:        16,
		Leaves:       []string{"leaf-c"},
		CreatedAt:    merkleTreeCreatedAt,
	}
	MerkleTree3 = types.MerkleTree{
		Type:         types.DvpERC721,
		TokenAddress: "0xTokenAddressBBB",
		Number:       1,
		Depth:        8,
		Leaves:       []string{"leaf-x"},
		CreatedAt:    merkleTreeCreatedAt,
	}
)

var (
	ModelMerkleTree1 = repository.MerkleTree{
		Type:         MerkleTree1.Type,
		TokenAddress: MerkleTree1.TokenAddress,
		Number:       MerkleTree1.Number,
		Depth:        MerkleTree1.Depth,
		Leaves:       MerkleTree1.Leaves,
		CreatedAt:    MerkleTree1.CreatedAt,
	}
	ModelMerkleTree2 = repository.MerkleTree{
		Type:         MerkleTree2.Type,
		TokenAddress: MerkleTree2.TokenAddress,
		Number:       MerkleTree2.Number,
		Depth:        MerkleTree2.Depth,
		Leaves:       MerkleTree2.Leaves,
		CreatedAt:    MerkleTree2.CreatedAt,
	}
	ModelMerkleTree3 = repository.MerkleTree{
		Type:         MerkleTree3.Type,
		TokenAddress: MerkleTree3.TokenAddress,
		Number:       MerkleTree3.Number,
		Depth:        MerkleTree3.Depth,
		Leaves:       MerkleTree3.Leaves,
		CreatedAt:    MerkleTree3.CreatedAt,
	}
)
