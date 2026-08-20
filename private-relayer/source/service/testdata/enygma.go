package testdata

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/source/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

// NewEnygmaTransferEventsFixture returns a slice of EnygmaTransferTx transactions
// (pre-split, one per destination) with dummy data suitable for testing transfer event handling.
func NewEnygmaTransferEventsFixture() []service.EnygmaTransferTx {
	return []service.EnygmaTransferTx{
		{
			MessageId:   "msg-1337-1",
			ReferenceId: common.HexToHash("0xaaaa0001"),
			FromAddress: common.HexToAddress("0xdeadbeef"),
			ToChainId:   big.NewInt(1337),
			ToAddress:   common.HexToAddress("0xc001babe"),
			ToAmount:    big.NewInt(10),
		},
		{
			MessageId:   "msg-1000-1",
			ReferenceId: common.HexToHash("0xaaaa0002"),
			FromAddress: common.HexToAddress("0xdeadbeef"),
			ToChainId:   big.NewInt(1000),
			ToAddress:   common.HexToAddress("0xdeadc0de"),
			ToAmount:    big.NewInt(20),
		},
		{
			MessageId:   "msg-1-1",
			ReferenceId: common.HexToHash("0xaaaa0003"),
			FromAddress: common.HexToAddress("0xc001aidd"),
			ToChainId:   big.NewInt(1),
			ToAddress:   common.HexToAddress("0xc0fedead"),
			ToAmount:    big.NewInt(50),
		},
		{
			MessageId:   "msg-2-1",
			ReferenceId: common.HexToHash("0xaaaa0004"),
			FromAddress: common.HexToAddress("0xc001aidd"),
			ToChainId:   big.NewInt(2),
			ToAddress:   common.HexToAddress("0xc0010fff"),
			ToAmount:    big.NewInt(200),
		},
	}
}

// NewTransfersByChainIDFixture returns a map of chain IDs to EnygmaTransferBatchTx slices
// with dummy data suitable for testing transfer grouping by chain ID.
func NewTransfersByChainIDFixture() map[string][]*types.EnygmaTransferBatchTx {
	return map[string][]*types.EnygmaTransferBatchTx{
		"1337": {
			{
				MessageId:     "msg-1337-1",
				ReferenceId:   common.HexToHash("0xaaaa0001"),
				FromAddress:   common.HexToAddress("0xdeadbeef"),
				ToAmount:      big.NewInt(10),
				ToAddress:     common.HexToAddress("0xc001babe"),
				SendTimestamp: 1234567890,
			},
		},
		"1000": {
			{
				MessageId:     "msg-1000-1",
				ReferenceId:   common.HexToHash("0xaaaa0002"),
				FromAddress:   common.HexToAddress("0xdeadbeef"),
				ToAmount:      big.NewInt(20),
				ToAddress:     common.HexToAddress("0xdeadc0de"),
				SendTimestamp: 1234567891,
			},
		},
		"1": {
			{
				MessageId:     "msg-1-1",
				ReferenceId:   common.HexToHash("0xaaaa0003"),
				FromAddress:   common.HexToAddress("0xc001aidd"),
				ToAmount:      big.NewInt(50),
				ToAddress:     common.HexToAddress("0xc0fedead"),
				SendTimestamp: 1234567892,
			},
		},
		"2": {
			{
				MessageId:     "msg-2-1",
				ReferenceId:   common.HexToHash("0xaaaa0004"),
				FromAddress:   common.HexToAddress("0xc001aidd"),
				ToAmount:      big.NewInt(200),
				ToAddress:     common.HexToAddress("0xc0010fff"),
				SendTimestamp: 1234567893,
			},
		},
	}
}

// NewBatchedTransfersFixture returns a slice of batched transfers (maps of chain IDs to EnygmaTransferBatchTx slices)
// with dummy data suitable for testing transfer batching logic.
func NewBatchedTransfersFixture() []map[string][]*types.EnygmaTransferBatchTx {
	return []map[string][]*types.EnygmaTransferBatchTx{
		{
			"1337": {
				{
					MessageId:     "msg-1337-1",
					ReferenceId:   common.HexToHash("0xaaaa0001"),
					FromAddress:   common.HexToAddress("0xdeadbeef"),
					ToAmount:      big.NewInt(10),
					ToAddress:     common.HexToAddress("0xc001babe"),
					SendTimestamp: 1234567890,
				},
			},
			"1000": {
				{
					MessageId:     "msg-1000-1",
					ReferenceId:   common.HexToHash("0xaaaa0002"),
					FromAddress:   common.HexToAddress("0xdeadbeef"),
					ToAmount:      big.NewInt(20),
					ToAddress:     common.HexToAddress("0xdeadc0de"),
					SendTimestamp: 1234567891,
				},
			},
		},
		{
			"1": {
				{
					MessageId:     "msg-1-1",
					ReferenceId:   common.HexToHash("0xaaaa0003"),
					FromAddress:   common.HexToAddress("0xc001aidd"),
					ToAmount:      big.NewInt(50),
					ToAddress:     common.HexToAddress("0xc0fedead"),
					SendTimestamp: 1234567892,
				},
			},
			"2": {
				{
					MessageId:     "msg-2-1",
					ReferenceId:   common.HexToHash("0xaaaa0004"),
					FromAddress:   common.HexToAddress("0xc001aidd"),
					ToAmount:      big.NewInt(200),
					ToAddress:     common.HexToAddress("0xc0010fff"),
					SendTimestamp: 1234567893,
				},
			},
		},
	}
}
