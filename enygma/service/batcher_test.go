package service_test

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/enygma/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/enygma/testutils"
	privatehubservice "github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/source/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockBatchingParticipantClient struct {
	GetEnygmaParticipantsFunc func() ([]*big.Int, error)

	GetEnygmaParticipantsCalls struct {
		sync.RWMutex
		Calls []struct{}
	}
}

func (m *MockBatchingParticipantClient) GetEnygmaParticipants(_ context.Context) ([]*big.Int, error) {
	m.GetEnygmaParticipantsCalls.Lock()
	m.GetEnygmaParticipantsCalls.Calls = append(m.GetEnygmaParticipantsCalls.Calls, struct{}{})
	m.GetEnygmaParticipantsCalls.Unlock()

	if m.GetEnygmaParticipantsFunc != nil {
		return m.GetEnygmaParticipantsFunc()
	}
	return nil, nil //nolint:nilnil // intentional nil return in test mock
}

func createResourceIdArray(id string) [32]byte {
	var arr [32]byte
	copy(arr[:], id)
	return arr
}

func createEnygmaTransferTx(messageID string, referenceId [32]byte, from common.Address, toChainId *big.Int, toAddress common.Address, toAmount *big.Int) privatehubservice.EnygmaTransferTx {
	return privatehubservice.EnygmaTransferTx{
		MessageId:   messageID,
		ReferenceId: referenceId,
		FromAddress: from,
		ToChainId:   toChainId,
		ToAddress:   toAddress,
		ToAmount:    toAmount,
	}
}

func TestEnygmaBatcher_BatchEnygmaSupplyUpdateEvents(t *testing.T) {
	t.Run("batches mints correctly", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 100,
			},
			mockClient,
			mockTracer,
		)

		mints := []*big.Int{big.NewInt(100), big.NewInt(200)}

		result := batcher.BatchEnygmaSupplyUpdateEvents(mints, nil)

		assert.Equal(t, types.EnygmaMint, result.Type)
		assert.Equal(t, big.NewInt(300), result.Amount)
	})

	t.Run("batches burns correctly", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 100,
			},
			mockClient,
			mockTracer,
		)

		burns := []*big.Int{big.NewInt(50), big.NewInt(75)}

		result := batcher.BatchEnygmaSupplyUpdateEvents(nil, burns)

		assert.Equal(t, types.EnygmaBurn, result.Type)
		assert.Equal(t, big.NewInt(125), result.Amount)
	})

	t.Run("net mints when mints exceed burns", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 100,
			},
			mockClient,
			mockTracer,
		)

		mints := []*big.Int{big.NewInt(500)}
		burns := []*big.Int{big.NewInt(200)}

		result := batcher.BatchEnygmaSupplyUpdateEvents(mints, burns)

		assert.Equal(t, types.EnygmaMint, result.Type)
		assert.Equal(t, big.NewInt(300), result.Amount)
	})

	t.Run("net burns when burns exceed mints", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 100,
			},
			mockClient,
			mockTracer,
		)

		mints := []*big.Int{big.NewInt(100)}
		burns := []*big.Int{big.NewInt(400)}

		result := batcher.BatchEnygmaSupplyUpdateEvents(mints, burns)

		assert.Equal(t, types.EnygmaBurn, result.Type)
		assert.Equal(t, big.NewInt(300), result.Amount)
	})

	t.Run("handles empty events", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 100,
			},
			mockClient,
			mockTracer,
		)

		result := batcher.BatchEnygmaSupplyUpdateEvents(nil, nil)

		assert.Equal(t, big.NewInt(0), result.Amount)
	})

	t.Run("aggregates multiple mints", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 100,
			},
			mockClient,
			mockTracer,
		)

		mints := []*big.Int{big.NewInt(100), big.NewInt(200), big.NewInt(150)}

		result := batcher.BatchEnygmaSupplyUpdateEvents(mints, nil)

		assert.Equal(t, types.EnygmaMint, result.Type)
		assert.Equal(t, big.NewInt(450), result.Amount)
	})

	t.Run("aggregates multiple burns", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 100,
			},
			mockClient,
			mockTracer,
		)

		burns := []*big.Int{big.NewInt(50), big.NewInt(75), big.NewInt(100)}

		result := batcher.BatchEnygmaSupplyUpdateEvents(nil, burns)

		assert.Equal(t, types.EnygmaBurn, result.Type)
		assert.Equal(t, big.NewInt(225), result.Amount)
	})
}

func TestEnygmaBatcher_GetTransferBatchLimits(t *testing.T) {
	t.Run("returns correct limits for 2 participants", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{
			GetEnygmaParticipantsFunc: func() ([]*big.Int, error) {
				return []*big.Int{big.NewInt(1), big.NewInt(2)}, nil
			},
		}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 100,
			},
			mockClient,
			mockTracer,
		)

		maxTxs, maxChainIDs, err := batcher.GetTransferBatchLimits(context.Background())

		require.NoError(t, err)
		assert.Equal(t, 100, maxTxs)
		assert.Equal(t, 1, maxChainIDs) // anonymity index 2 - 1
	})

	t.Run("returns correct limits for 5 participants", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{
			GetEnygmaParticipantsFunc: func() ([]*big.Int, error) {
				return []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4), big.NewInt(5)}, nil
			},
		}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 50,
			},
			mockClient,
			mockTracer,
		)

		maxTxs, maxChainIDs, err := batcher.GetTransferBatchLimits(context.Background())

		require.NoError(t, err)
		assert.Equal(t, 50, maxTxs)
		assert.Equal(t, 4, maxChainIDs) // anonymity index 5 - 1
	})

	t.Run("returns capped limits for more than 6 participants", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{
			GetEnygmaParticipantsFunc: func() ([]*big.Int, error) {
				return []*big.Int{
					big.NewInt(1),
					big.NewInt(2),
					big.NewInt(3),
					big.NewInt(4),
					big.NewInt(5),
					big.NewInt(6),
					big.NewInt(7),
				}, nil
			},
		}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 200,
			},
			mockClient,
			mockTracer,
		)

		maxTxs, maxChainIDs, err := batcher.GetTransferBatchLimits(context.Background())

		require.NoError(t, err)
		assert.Equal(t, 200, maxTxs)
		assert.Equal(t, 5, maxChainIDs) // anonymity index capped at 6 - 1
	})

	t.Run("returns error for insufficient participants", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{
			GetEnygmaParticipantsFunc: func() ([]*big.Int, error) {
				return []*big.Int{big.NewInt(1)}, nil
			},
		}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 100,
			},
			mockClient,
			mockTracer,
		)

		_, _, err := batcher.GetTransferBatchLimits(context.Background())

		require.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient banks")
	})

	t.Run("returns error when client fails", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{
			GetEnygmaParticipantsFunc: func() ([]*big.Int, error) {
				return nil, errors.New("client error")
			},
		}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 100,
			},
			mockClient,
			mockTracer,
		)

		_, _, err := batcher.GetTransferBatchLimits(context.Background())

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get enygma participants")
	})
}

func TestEnygmaBatcher_GroupTransfersByChainID(t *testing.T) {
	t.Run("groups transfers by chain ID and maps fields correctly", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 100,
			},
			mockClient,
			mockTracer,
		)

		refId1 := createResourceIdArray("ref1")
		refId2 := createResourceIdArray("ref2")
		refId3 := createResourceIdArray("ref3")

		txs := []privatehubservice.EnygmaTransferTx{
			createEnygmaTransferTx("msg-1", refId1, common.Address{1}, big.NewInt(2), common.Address{2}, big.NewInt(100)),
			createEnygmaTransferTx("msg-2", refId2, common.Address{1}, big.NewInt(2), common.Address{3}, big.NewInt(200)),
			createEnygmaTransferTx("msg-3", refId3, common.Address{1}, big.NewInt(3), common.Address{4}, big.NewInt(300)),
		}

		result := batcher.GroupTransfersByChainID(txs)

		// Verify grouping by chain ID
		assert.Equal(t, 2, len(result))
		assert.Equal(t, 2, len(result["2"]))
		assert.Equal(t, 1, len(result["3"]))

		// Verify field mapping for chain ID "2" transactions
		chain2Txs := result["2"]
		assert.Equal(t, "msg-1", chain2Txs[0].MessageId)
		assert.Equal(t, refId1, chain2Txs[0].ReferenceId)
		assert.Equal(t, common.Address{2}, chain2Txs[0].ToAddress)
		assert.Equal(t, big.NewInt(100), chain2Txs[0].ToAmount)

		assert.Equal(t, "msg-2", chain2Txs[1].MessageId)
		assert.Equal(t, refId2, chain2Txs[1].ReferenceId)
		assert.Equal(t, common.Address{3}, chain2Txs[1].ToAddress)
		assert.Equal(t, big.NewInt(200), chain2Txs[1].ToAmount)

		// Verify field mapping for chain ID "3" transaction
		chain3Txs := result["3"]
		assert.Equal(t, "msg-3", chain3Txs[0].MessageId)
		assert.Equal(t, refId3, chain3Txs[0].ReferenceId)
		assert.Equal(t, common.Address{4}, chain3Txs[0].ToAddress)
		assert.Equal(t, big.NewInt(300), chain3Txs[0].ToAmount)
	})

	t.Run("handles empty events", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 100,
			},
			mockClient,
			mockTracer,
		)

		result := batcher.GroupTransfersByChainID([]privatehubservice.EnygmaTransferTx{})

		assert.Equal(t, 0, len(result))
	})
}

func TestEnygmaBatcher_BatchTransfers(t *testing.T) {
	t.Run("creates single batch when under limits", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{
			GetEnygmaParticipantsFunc: func() ([]*big.Int, error) {
				// Return 6 participants to get maxChainIDs=5 (high enough for 2 chain IDs)
				return []*big.Int{
					big.NewInt(1),
					big.NewInt(2),
					big.NewInt(3),
					big.NewInt(4),
					big.NewInt(5),
					big.NewInt(6),
				}, nil
			},
		}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 100,
			},
			mockClient,
			mockTracer,
		)

		txsByChainID := map[string][]*types.EnygmaTransferBatchTx{
			"2": {
				{ToAddress: common.Address{1}, ToAmount: big.NewInt(100)},
				{ToAddress: common.Address{2}, ToAmount: big.NewInt(200)},
			},
			"3": {
				{ToAddress: common.Address{3}, ToAmount: big.NewInt(300)},
			},
		}

		batches, err := batcher.BatchTransfers(context.Background(), txsByChainID)
		require.NoError(t, err)
		assert.Equal(t, 1, len(batches))
		assert.Equal(t, 2, len(batches[0]))
		assert.Equal(t, 2, len(batches[0]["2"]))
		assert.Equal(t, 1, len(batches[0]["3"]))
	})

	t.Run("splits into multiple batches when tx limit exceeded", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{
			GetEnygmaParticipantsFunc: func() ([]*big.Int, error) {
				// Return 6 participants to get maxChainIDs=5 (high enough for 1 chain ID)
				return []*big.Int{
					big.NewInt(1),
					big.NewInt(2),
					big.NewInt(3),
					big.NewInt(4),
					big.NewInt(5),
					big.NewInt(6),
				}, nil
			},
		}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 2, // Set to 2 to force splitting 3 txs into 2 batches
			},
			mockClient,
			mockTracer,
		)

		txsByChainID := map[string][]*types.EnygmaTransferBatchTx{
			"2": {
				{ToAddress: common.Address{1}, ToAmount: big.NewInt(100)},
				{ToAddress: common.Address{2}, ToAmount: big.NewInt(200)},
				{ToAddress: common.Address{3}, ToAmount: big.NewInt(300)},
			},
		}

		batches, err := batcher.BatchTransfers(context.Background(), txsByChainID)
		require.NoError(t, err)

		require.NoError(t, err)
		assert.Equal(t, 2, len(batches))
		// Total transaction count across all batches should be 3
		totalTxs := 0
		for _, batch := range batches {
			for _, txs := range batch {
				totalTxs += len(txs)
			}
		}
		assert.Equal(t, 3, totalTxs)
	})

	t.Run("splits into multiple batches when chain ID limit exceeded", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{
			GetEnygmaParticipantsFunc: func() ([]*big.Int, error) {
				// Return 2 participants to get maxChainIDs=1 (forces splitting across 3 chain IDs)
				return []*big.Int{big.NewInt(1), big.NewInt(2)}, nil
			},
		}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 100,
			},
			mockClient,
			mockTracer,
		)

		txsByChainID := map[string][]*types.EnygmaTransferBatchTx{
			"2": {{ToAddress: common.Address{1}, ToAmount: big.NewInt(100)}},
			"3": {{ToAddress: common.Address{2}, ToAmount: big.NewInt(200)}},
			"4": {{ToAddress: common.Address{3}, ToAmount: big.NewInt(300)}},
		}

		batches, err := batcher.BatchTransfers(context.Background(), txsByChainID)
		require.NoError(t, err)

		require.NoError(t, err)
		assert.Equal(t, 3, len(batches))
	})

	t.Run("handles empty transfers", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{
			GetEnygmaParticipantsFunc: func() ([]*big.Int, error) {
				// Return 6 participants to get maxChainIDs=5
				return []*big.Int{
					big.NewInt(1),
					big.NewInt(2),
					big.NewInt(3),
					big.NewInt(4),
					big.NewInt(5),
					big.NewInt(6),
				}, nil
			},
		}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 100,
			},
			mockClient,
			mockTracer,
		)

		batches, err := batcher.BatchTransfers(context.Background(), map[string][]*types.EnygmaTransferBatchTx{})
		require.NoError(t, err)

		require.NoError(t, err)
		assert.Equal(t, 0, len(batches))
	})
}

func TestEnygmaBatcher_CreateBatchesWithAnonimity(t *testing.T) {
	t.Run("fulfills anonymity set when transfers to only 1 destination", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{
			GetEnygmaParticipantsFunc: func() ([]*big.Int, error) {
				return []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}, nil
			},
		}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 100,
			},
			mockClient,
			mockTracer,
		)

		txsByChainID := map[string][]*types.EnygmaTransferBatchTx{
			"2": {{MessageId: "msg1", ToAddress: common.Address{1}, ToAmount: big.NewInt(100)}},
		}

		ctx := context.Background()
		batches, err := batcher.CreateBatchesWithAnonimity(ctx, "resourceId", big.NewInt(100), txsByChainID)

		require.NoError(t, err)
		assert.Equal(t, 3, len(batches))
		assert.Equal(t, big.NewInt(1), batches[0].ToChainID)
		assert.Equal(t, big.NewInt(2), batches[1].ToChainID)
		assert.Equal(t, big.NewInt(3), batches[2].ToChainID)
		assert.Equal(t, 0, len(batches[0].Transactions))
		assert.Equal(t, 1, len(batches[1].Transactions))
		assert.Equal(t, 0, len(batches[2].Transactions))
	})

	t.Run("fulfills anonymity set when transfers to only 1 destination (sorted)", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{
			GetEnygmaParticipantsFunc: func() ([]*big.Int, error) {
				return []*big.Int{big.NewInt(3), big.NewInt(2), big.NewInt(1)}, nil
			},
		}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 100,
			},
			mockClient,
			mockTracer,
		)

		txsByChainID := map[string][]*types.EnygmaTransferBatchTx{
			"2": {{MessageId: "msg1", ToAddress: common.Address{1}, ToAmount: big.NewInt(100)}},
		}

		ctx := context.Background()
		batches, err := batcher.CreateBatchesWithAnonimity(ctx, "resourceId", big.NewInt(100), txsByChainID)

		require.NoError(t, err)
		assert.Equal(t, 3, len(batches))
		assert.Equal(t, big.NewInt(1), batches[0].ToChainID)
		assert.Equal(t, big.NewInt(2), batches[1].ToChainID)
		assert.Equal(t, big.NewInt(3), batches[2].ToChainID)
		assert.Equal(t, 0, len(batches[0].Transactions))
		assert.Equal(t, 1, len(batches[1].Transactions))
		assert.Equal(t, 0, len(batches[2].Transactions))
	})

	t.Run("returns error when getting participants fails", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{
			GetEnygmaParticipantsFunc: func() ([]*big.Int, error) {
				return nil, errors.New("db error")
			},
		}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 100,
			},
			mockClient,
			mockTracer,
		)

		ctx := context.Background()
		_, err := batcher.CreateBatchesWithAnonimity(
			ctx,
			"resourceId",
			big.NewInt(100),
			map[string][]*types.EnygmaTransferBatchTx{},
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get enygma participants")
	})

	t.Run("creates batch for each destination chain", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{
			GetEnygmaParticipantsFunc: func() ([]*big.Int, error) {
				return []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}, nil
			},
		}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 100,
			},
			mockClient,
			mockTracer,
		)

		txsByChainID := map[string][]*types.EnygmaTransferBatchTx{
			"2": {{MessageId: "msg1", ToAddress: common.Address{1}, ToAmount: big.NewInt(100)}},
			"3": {{MessageId: "msg2", ToAddress: common.Address{2}, ToAmount: big.NewInt(200)}},
		}

		ctx := context.Background()
		batches, err := batcher.CreateBatchesWithAnonimity(ctx, "resourceId", big.NewInt(100), txsByChainID)

		require.NoError(t, err)
		// Should have a batch for each destination chain + padding
		assert.Equal(t, 3, len(batches))
		assert.Equal(t, 0, len(batches[0].Transactions))
		assert.Equal(t, 1, len(batches[1].Transactions))
		assert.Equal(t, 1, len(batches[2].Transactions))
		assert.Equal(t, big.NewInt(1), batches[0].ToChainID)
		assert.Equal(t, big.NewInt(2), batches[1].ToChainID)
		assert.Equal(t, big.NewInt(3), batches[2].ToChainID)
	})

	t.Run("fulfills anonymity set when no transfers", func(t *testing.T) {
		mockClient := &MockBatchingParticipantClient{
			GetEnygmaParticipantsFunc: func() ([]*big.Int, error) {
				return []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}, nil
			},
		}
		mockTracer := &testutils.MockTracer{}

		batcher := service.NewEnygmaBatcher(
			&service.EnygmaBatcherConfig{
				ChainID:        big.NewInt(1),
				MaxTxsPerBatch: 100,
			},
			mockClient,
			mockTracer,
		)

		ctx := context.Background()
		batches, err := batcher.CreateBatchesWithAnonimity(
			ctx,
			"resourceId",
			big.NewInt(100),
			map[string][]*types.EnygmaTransferBatchTx{},
		)

		require.NoError(t, err)
		assert.Equal(t, 3, len(batches))
		assert.Equal(t, 0, len(batches[0].Transactions))
		assert.Equal(t, 0, len(batches[1].Transactions))
		assert.Equal(t, 0, len(batches[2].Transactions))
	})
}
