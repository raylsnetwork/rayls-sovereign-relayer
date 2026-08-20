package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
	privatehubservice "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type EnygmaBatcherConfig struct {
	ChainID        *big.Int
	MaxTxsPerBatch int
}

type battingParticipantClient interface {
	GetEnygmaParticipants(ctx context.Context) ([]*big.Int, error)
}

type batcherTracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

type EnygmaBatcher struct {
	conf     *EnygmaBatcherConfig
	psClient battingParticipantClient
	tracer   batcherTracer
}

func NewEnygmaBatcher(
	conf *EnygmaBatcherConfig,
	psClient battingParticipantClient,
	tracer batcherTracer,
) *EnygmaBatcher {
	return &EnygmaBatcher{
		conf:     conf,
		psClient: psClient,
		tracer:   tracer,
	}
}

func (s *EnygmaBatcher) BatchEnygmaSupplyUpdateEvents(
	mints []*big.Int,
	burns []*big.Int,
) types.EnygmaSupplyUpdate {
	totalMintAmount := big.NewInt(0)
	totalBurnAmount := big.NewInt(0)

	// Sum all mint amounts
	for _, amount := range mints {
		totalMintAmount.Add(totalMintAmount, amount)
	}

	// Sum all burn amounts
	for _, amount := range burns {
		totalBurnAmount.Add(totalBurnAmount, amount)
	}

	// Return net supply update
	supplyUpdate := types.EnygmaSupplyUpdate{}
	if totalMintAmount.Cmp(totalBurnAmount) > 0 {
		supplyUpdate.Type = types.EnygmaMint
		supplyUpdate.Amount = totalMintAmount.Sub(totalMintAmount, totalBurnAmount)
	} else {
		supplyUpdate.Type = types.EnygmaBurn
		supplyUpdate.Amount = totalBurnAmount.Sub(totalBurnAmount, totalMintAmount)
	}

	return supplyUpdate
}

// GetTransferBatchLimits calculates the maximum number of transactions per batch
// and maximum number of unique chain IDs per batch based on the anonymity index
func (s *EnygmaBatcher) GetTransferBatchLimits(ctx context.Context) (int, int, error) {
	enygmaParticipants, err := s.psClient.GetEnygmaParticipants(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get enygma participants: %w", err)
	}

	anonymityIndex, err := getAnonymityIndex(len(enygmaParticipants))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get anonymity index: %w", err)
	}

	// Our chain is included in the anonymity index, so we must reserve one slot for it in the batch.
	// Later, we will send an empty batch to ourselves which might exceed the anonymity index(PLs > k).
	maxChainIDsPerBatch := anonymityIndex - 1

	return s.conf.MaxTxsPerBatch, maxChainIDsPerBatch, nil
}

// GroupTransfersByChainID groups transfer transactions by their destination chain IDs.
// Input transactions are already split (one per destination) with unique MessageIDs.
func (s *EnygmaBatcher) GroupTransfersByChainID(
	events []privatehubservice.EnygmaTransferTx,
) map[string][]*types.EnygmaTransferBatchTx {
	txsByChainID := make(map[string][]*types.EnygmaTransferBatchTx)

	for _, event := range events {
		chainID := event.ToChainId.String()

		if _, exists := txsByChainID[chainID]; !exists {
			txsByChainID[chainID] = []*types.EnygmaTransferBatchTx{}
		}

		txsByChainID[chainID] = append(txsByChainID[chainID], &types.EnygmaTransferBatchTx{
			MessageId:     event.MessageId,
			ReferenceId:   event.ReferenceId,
			FromAddress:   event.FromAddress,
			ToAmount:      event.ToAmount,
			ToAddress:     event.ToAddress,
			ProgramData:   event.ProgramData,
			SendTimestamp: time.Now().UnixMilli(),
		})
	}

	return txsByChainID
}

// BatchTransfers creates batches from transfers while respecting constraints:
// 1. Total transactions across all chains in a batch <= maxTxsPerBatch
// 2. Number of unique chain IDs in a batch <= maxChainIDsPerBatch
func (s *EnygmaBatcher) BatchTransfers(
	ctx context.Context,
	txsByChainID map[string][]*types.EnygmaTransferBatchTx,
) ([]map[string][]*types.EnygmaTransferBatchTx, error) {
	maxTxsPerBatch, maxChainIDsPerBatch, err := s.GetTransferBatchLimits(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transfer batch limits: %w", err)
	}

	batches := make([]map[string][]*types.EnygmaTransferBatchTx, 0)
	currentBatch := make(map[string][]*types.EnygmaTransferBatchTx)
	currentBatchTxCount := 0

	for chainID, txs := range txsByChainID {
		for _, tx := range txs {
			// Check if adding this transaction would exceed limits
			txLimitReached := currentBatchTxCount >= maxTxsPerBatch
			chainLimitReached := len(currentBatch) >= maxChainIDsPerBatch && currentBatch[chainID] == nil

			if txLimitReached || chainLimitReached {
				if currentBatchTxCount > 0 {
					batches = append(batches, currentBatch)
				}
				currentBatch = make(map[string][]*types.EnygmaTransferBatchTx)
				currentBatchTxCount = 0
			}

			// Add transaction to the current batch
			if currentBatch[chainID] == nil {
				currentBatch[chainID] = make([]*types.EnygmaTransferBatchTx, 0)
			}
			currentBatch[chainID] = append(currentBatch[chainID], tx)
			currentBatchTxCount++
		}
	}

	// Add final batch if it has transactions
	if currentBatchTxCount > 0 {
		batches = append(batches, currentBatch)
	}

	return batches, nil
}

// CreateBatchesWithAnonimity creates transfer batches with anonymity set fulfillment
// Ensures that the anonymity set is filled to the required size k
func (s *EnygmaBatcher) CreateBatchesWithAnonimity(
	ctx context.Context,
	resourceId string,
	blockNumber *big.Int,
	txsByChainID map[string][]*types.EnygmaTransferBatchTx,
) ([]*types.EnygmaTransferBatch, error) {
	_, span := s.tracer.Start(ctx, "create_batches_with_anonymity")
	defer span.End()

	// Construct batches
	enygmaParticipants, err := s.psClient.GetEnygmaParticipants(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get enygma participants: %w", err)
	}

	anonymityIndex, err := getAnonymityIndex(len(enygmaParticipants))
	if err != nil {
		return nil, fmt.Errorf("failed to get anonymity index: %w", err)
	}

	destChainIDs := []*big.Int{}

	// Extract all destination chain IDs.
	for chainID := range txsByChainID {
		var chainIDInt int64
		chainIDInt, err = strconv.ParseInt(chainID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse chain ID %s: %w", chainID, err)
		}
		destChainIDs = append(destChainIDs, big.NewInt(chainIDInt))
	}
	// Include the sender chain as a destination.
	destChainIDs = append([]*big.Int{s.conf.ChainID}, destChainIDs...)
	destChainIDs, err = fulfillAnonimitySet(anonymityIndex, destChainIDs, enygmaParticipants)
	if err != nil {
		return nil, fmt.Errorf("failed to fulfill anonymity set: %w", err)
	}

	// The proof requires the chain IDs to be sorted in an ascending order.
	// If they're not sorted - PL B -> PL A fails in the proofApi API.
	destChainIDs = sortAscendingOrderBigIntArray(destChainIDs)

	batches := make([]*types.EnygmaTransferBatch, len(destChainIDs))
	// Initializing the batches.
	for i, toChainID := range destChainIDs {
		batches[i] = &types.EnygmaTransferBatch{
			BatchId:                    uuid.New().String(),
			ResourceId:                 resourceId,
			BlockNumberPrivateHub:      blockNumber,
			FromChainID:                s.conf.ChainID,
			ToChainID:                  toChainID,
			SoftFinalityStartTimestamp: time.Now().UnixMilli(),
			// Will be populated once the proof is generated.
			ToRValueToAdd: nil,
			Transactions:  []*types.EnygmaTransferBatchTx{},
		}
	}

	for _, batch := range batches {
		if txs, ok := txsByChainID[batch.ToChainID.String()]; ok {
			batch.Transactions = txs
		}
	}

	span.SetStatus(codes.Ok, "batches created successfully")
	return batches, nil
}

// maxAnonymityBanks is the maximum number of banks used for the anonymity set.
const maxAnonymityBanks = 6

// getAnonymityIndex determines the anonymity index based on the number of banks/participants
// Returns a value between 2 and 6, or an error if there are fewer than 2 participants
func getAnonymityIndex(qtyBanks int) (int, error) {
	if qtyBanks < 2 {
		return 0, fmt.Errorf("insufficient banks: need at least 2, got %d", qtyBanks)
	}
	if qtyBanks > maxAnonymityBanks {
		return maxAnonymityBanks, nil
	}
	return qtyBanks, nil
}

// fulfillAnonimitySet fills the anonymity set to the required size k by adding padding chain IDs.
// Returns exactly k unique chain IDs sorted in ascending order.
func fulfillAnonimitySet(k int, currentChainIDs []*big.Int, allChainIDs []*big.Int) ([]*big.Int, error) {
	// Build deduplicated sets
	currentSet := make(map[string]*big.Int)
	for _, id := range currentChainIDs {
		currentSet[id.String()] = id
	}

	allSet := make(map[string]*big.Int)
	for _, id := range allChainIDs {
		allSet[id.String()] = id
	}

	// Validate: currentChainIDs must be a subset of allChainIDs
	if len(currentSet) > len(allSet) {
		return nil, fmt.Errorf(
			"currentChainIDs size (%d) cannot be greater than allChainIDs size (%d)",
			len(currentSet),
			len(allSet),
		)
	}

	// Validate: allChainIDs must have at least k elements
	if len(allSet) < k {
		return nil, fmt.Errorf("allChainIDs size (%d) is less than required anonymity set size k (%d)", len(allSet), k)
	}

	// If allChainIDs size == k, return all of them sorted
	if len(allSet) == k {
		result := make([]*big.Int, 0, k)
		for _, id := range allSet {
			result = append(result, id)
		}
		return sortAscendingOrderBigIntArray(result), nil
	}

	// Collect IDs not in currentSet
	remaining := make([]*big.Int, 0, len(allSet)-len(currentSet))
	for key, id := range allSet {
		if _, exists := currentSet[key]; !exists {
			remaining = append(remaining, id)
		}
	}

	// Start result with current IDs
	result := make([]*big.Int, 0, k)
	for _, id := range currentSet {
		result = append(result, id)
	}

	// Fill up to k by randomly selecting from remaining
	needed := k - len(result)
	if needed > 0 {
		// Fisher-Yates shuffle using crypto/rand
		for i := len(remaining) - 1; i > 0; i-- {
			jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
			if err != nil {
				return nil, fmt.Errorf("generating random index for shuffle: %w", err)
			}
			j := int(jBig.Int64())
			remaining[i], remaining[j] = remaining[j], remaining[i]
		}
		result = append(result, remaining[:needed]...)
	}

	return sortAscendingOrderBigIntArray(result), nil
}

// sortAscendingOrderBigIntArray sorts a slice of big.Int values in ascending order
// Returns a new sorted slice without modifying the input
func sortAscendingOrderBigIntArray(arr []*big.Int) []*big.Int {
	// Create a copy of the input array
	sortedArr := make([]*big.Int, len(arr))
	for i, v := range arr {
		sortedArr[i] = new(big.Int).Set(v) // Create a deep copy of each big.Int
	}

	// Sort the copy in ascending order
	sort.Slice(sortedArr, func(i, j int) bool {
		return sortedArr[i].Cmp(sortedArr[j]) < 0
	})

	return sortedArr
}
