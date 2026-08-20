package testdata

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

var (
	ModelTransaction1 = repository.Transaction{
		TxHash:                  "0xabc123def456ghi789",
		TxHashDestination:       "0xdef456ghi789abc123",
		LogIndex:                1,
		SharedId:                uuid.New().String(),
		State:                   types.SourcePublish,
		Outcome:                 types.OutcomePending,
		OriginatorChainId:       "1",
		DestinationChainId:      "2",
		MsgId:                   [32]byte{0xaa, 0xbb, 0xcc},
		IsAtomic:                true,
		UpdatedAt:               time.Now().UTC().Truncate(time.Millisecond),
		BatchId:                 uuid.New().String(),
		BatchTxHashOnPrivateHub: "0x123456abcdef",
		ResourceId:              "resource-123",
		FromContractAddress:     "0xcontractaddress1",
		FromUserAddress:         "0xuseraddress1",
		TransferMetadata_Id:     uuid.New().String(),
		TransferMetadata_Amount: "1000000000000000000",
		BlockNumber:             12345678,
		ParentHash:              "0xparenthash1",
	}

	ModelTransaction2 = repository.Transaction{
		TxHash:                  "0x987zyx654wvu321tsr",
		TxHashDestination:       "0x654wvu321tsr987zyx",
		LogIndex:                2,
		SharedId:                uuid.New().String(),
		State:                   types.SourcePublish,
		Outcome:                 types.OutcomeSuccess,
		OriginatorChainId:       "3",
		DestinationChainId:      "4",
		MsgId:                   [32]byte{0xdd, 0xee, 0xff},
		IsAtomic:                false,
		UpdatedAt:               time.Now().UTC().Truncate(time.Millisecond),
		BatchId:                 uuid.New().String(),
		BatchTxHashOnPrivateHub: "0xabcdef123456",
		ResourceId:              "resource-456",
		FromContractAddress:     "0xcontractaddress2",
		FromUserAddress:         "0xuseraddress2",
		TransferMetadata_Id:     uuid.New().String(),
		TransferMetadata_Amount: "500000000000000000",
		BlockNumber:             87654321,
		ParentHash:              "0xparenthash2",
	}
	ModelTransaction3 = repository.Transaction{
		TxHash:                  "0xaaaabbbbccccddddeeee",
		TxHashDestination:       "0xbbbbccccddddeeeeffff",
		LogIndex:                3,
		SharedId:                uuid.New().String(),
		State:                   types.DestinationDispatch,
		Outcome:                 types.OutcomePending,
		OriginatorChainId:       "5",
		DestinationChainId:      "6",
		MsgId:                   [32]byte{0x11, 0x22, 0x33},
		IsAtomic:                false,
		UpdatedAt:               time.Now().UTC(),
		BatchId:                 uuid.New().String(),
		BatchTxHashOnPrivateHub: "0xdeadbeefdeadbeef",
		ResourceId:              "resource-789",
		FromContractAddress:     "0xcontractaddress3",
		FromUserAddress:         "0xuseraddress3",
		TransferMetadata_Id:     uuid.New().String(),
		TransferMetadata_Amount: "2500000000000000000",
		BlockNumber:             13579135,
		ParentHash:              "0xparenthash3",
	}
)

var (
	Transaction1 = types.Transaction{
		TxHash:              ModelTransaction1.TxHash,
		TxHashDestination:   common.HexToHash(ModelTransaction1.TxHashDestination),
		LogIndex:            ModelTransaction1.LogIndex,
		SharedID:            ModelTransaction1.SharedId,
		State:               ModelTransaction1.State,
		Outcome:             ModelTransaction1.Outcome,
		FromChainID:         big.NewInt(1),
		ToChainID:           big.NewInt(2),
		MsgID:               ModelTransaction1.MsgId,
		IsAtomic:            ModelTransaction1.IsAtomic,
		UpdatedAt:           ModelTransaction1.UpdatedAt,
		BatchID:             ModelTransaction1.BatchId,
		BatchPrivateHubHash: common.HexToHash(ModelTransaction1.BatchTxHashOnPrivateHub),
		ResourceID:          ModelTransaction1.ResourceId,
		FromContractAddress: ModelTransaction1.FromContractAddress,
		FromUserAddress:     ModelTransaction1.FromUserAddress,
		TransferID:          ModelTransaction1.TransferMetadata_Id,
		TransferAmount:      ModelTransaction1.TransferMetadata_Amount,
		BlockNumber:         ModelTransaction1.BlockNumber,
		ParentHash:          ModelTransaction1.ParentHash,
	}

	Transaction2 = types.Transaction{
		TxHash:              ModelTransaction2.TxHash,
		TxHashDestination:   common.HexToHash(ModelTransaction2.TxHashDestination),
		LogIndex:            ModelTransaction2.LogIndex,
		SharedID:            ModelTransaction2.SharedId,
		State:               ModelTransaction2.State,
		Outcome:             ModelTransaction2.Outcome,
		FromChainID:         big.NewInt(3),
		ToChainID:           big.NewInt(4),
		MsgID:               ModelTransaction2.MsgId,
		IsAtomic:            ModelTransaction2.IsAtomic,
		UpdatedAt:           ModelTransaction2.UpdatedAt,
		BatchID:             ModelTransaction2.BatchId,
		BatchPrivateHubHash: common.HexToHash(ModelTransaction2.BatchTxHashOnPrivateHub),
		ResourceID:          ModelTransaction2.ResourceId,
		FromContractAddress: ModelTransaction2.FromContractAddress,
		FromUserAddress:     ModelTransaction2.FromUserAddress,
		TransferID:          ModelTransaction2.TransferMetadata_Id,
		TransferAmount:      ModelTransaction2.TransferMetadata_Amount,
		BlockNumber:         ModelTransaction2.BlockNumber,
		ParentHash:          ModelTransaction2.ParentHash,
	}
	Transaction3 = types.Transaction{
		TxHash:              ModelTransaction3.TxHash,
		TxHashDestination:   common.HexToHash(ModelTransaction3.TxHashDestination),
		LogIndex:            ModelTransaction3.LogIndex,
		SharedID:            ModelTransaction3.SharedId,
		State:               ModelTransaction3.State,
		Outcome:             ModelTransaction3.Outcome,
		FromChainID:         big.NewInt(5),
		ToChainID:           big.NewInt(6),
		MsgID:               ModelTransaction3.MsgId,
		IsAtomic:            ModelTransaction3.IsAtomic,
		UpdatedAt:           ModelTransaction3.UpdatedAt,
		BatchID:             ModelTransaction3.BatchId,
		BatchPrivateHubHash: common.HexToHash(ModelTransaction3.BatchTxHashOnPrivateHub),
		ResourceID:          ModelTransaction3.ResourceId,
		FromContractAddress: ModelTransaction3.FromContractAddress,
		FromUserAddress:     ModelTransaction3.FromUserAddress,
		TransferID:          ModelTransaction3.TransferMetadata_Id,
		TransferAmount:      ModelTransaction3.TransferMetadata_Amount,
		BlockNumber:         ModelTransaction3.BlockNumber,
		ParentHash:          ModelTransaction3.ParentHash,
	}
)
