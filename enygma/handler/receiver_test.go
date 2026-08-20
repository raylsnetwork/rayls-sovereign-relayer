package handler_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/txsim"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
)

// createBatch creates a batch for testing HandleEnygmaCrossTransfer
func createBatch(fromChainID, toChainID *big.Int, numTxs int) *types.EnygmaTransferBatch {
	txs := make([]*types.EnygmaTransferBatchTx, 0)

	if fromChainID.Cmp(toChainID) != 0 {
		for i := 0; i < numTxs; i++ {
			txs = append(txs, &types.EnygmaTransferBatchTx{
				MessageId:   uuid.New().String(), // random message id
				ReferenceId: [32]byte{byte(i + 1)},
				FromAddress: common.HexToAddress("0x1111111111111111111111111111111111111111"),
				ToAmount:    big.NewInt(100),
				ToAddress:   common.HexToAddress("0x2222222222222222222222222222222222222222"),
				ProgramData: []types.EnygmaProgramData{
					{ResourceId: [32]byte{byte(i + 1)}, Selector: [4]byte{0xAA}, Args: []byte{byte(i + 1), 0xAA}},
				}, // single mint step per recipient
				SendTimestamp: time.Now().UnixMilli(),
			})
		}
	}

	// random batch id
	return &types.EnygmaTransferBatch{
		ResourceId:            "0000000000000000000000000000000000000000000000000000000000000001", // fixed resource id - no 0x prefix
		BlockNumberPrivateHub: big.NewInt(100),
		FromChainID:           fromChainID,
		ToChainID:             toChainID,
		ToRValueToAdd:         big.NewInt(0),
		Transactions:          txs,
		BatchId:               uuid.New().String(), // random batch id
		Ctx:                   context.Background(),
	}
}

// makeSuccessfulMint returns a ReceiveDestTransferBatch func where every transfer reports a successful receipt.
func makeSuccessfulMint() func(ctx context.Context, transfers []*types.EnygmaCrossTransferData) (map[string]contractclient.BatchResult, error) {
	return func(ctx context.Context, transfers []*types.EnygmaCrossTransferData) (map[string]contractclient.BatchResult, error) {
		results := make(map[string]contractclient.BatchResult, len(transfers))
		for _, t := range transfers {
			results[t.MessageId] = contractclient.BatchResult{
				Receipt: &ethTypes.Receipt{Status: 1, TxHash: common.HexToHash("0xabc")},
			}
		}
		return results, nil
	}
}

// makeAlternatingMint returns success/failure receipts alternately for the supplied transfers.
func makeAlternatingMint() func(ctx context.Context, transfers []*types.EnygmaCrossTransferData) (map[string]contractclient.BatchResult, error) {
	return func(ctx context.Context, transfers []*types.EnygmaCrossTransferData) (map[string]contractclient.BatchResult, error) {
		results := make(map[string]contractclient.BatchResult, len(transfers))
		for i, t := range transfers {
			results[t.MessageId] = contractclient.BatchResult{
				Receipt: &ethTypes.Receipt{Status: 1 - uint64(i%2), TxHash: common.HexToHash("0xdef")},
			}
		}
		return results, nil
	}
}

// makeFailedMint returns BatchResults with Err set for every transfer.
func makeFailedMint() func(ctx context.Context, transfers []*types.EnygmaCrossTransferData) (map[string]contractclient.BatchResult, error) {
	return func(ctx context.Context, transfers []*types.EnygmaCrossTransferData) (map[string]contractclient.BatchResult, error) {
		results := make(map[string]contractclient.BatchResult, len(transfers))
		for _, t := range transfers {
			results[t.MessageId] = contractclient.BatchResult{
				Err: errors.New("send failed"),
			}
		}
		return results, nil
	}
}

// makeRevertedMint returns BatchResults with Receipt.Status=0 for every transfer (on-chain revert).
func makeRevertedMint() func(ctx context.Context, transfers []*types.EnygmaCrossTransferData) (map[string]contractclient.BatchResult, error) {
	return func(ctx context.Context, transfers []*types.EnygmaCrossTransferData) (map[string]contractclient.BatchResult, error) {
		results := make(map[string]contractclient.BatchResult, len(transfers))
		for _, t := range transfers {
			results[t.MessageId] = contractclient.BatchResult{
				Receipt: &ethTypes.Receipt{Status: 0, TxHash: common.HexToHash("0xfee")},
			}
		}
		return results, nil
	}
}

// makeSuccessfulRevert returns successful revert results.
func makeSuccessfulRevert() func(ctx context.Context, tokenAddress common.Address, reverts []*types.EnygmaCrossTransferData) (map[string]contractclient.BatchResult, error) {
	return func(ctx context.Context, tokenAddress common.Address, reverts []*types.EnygmaCrossTransferData) (map[string]contractclient.BatchResult, error) {
		results := make(map[string]contractclient.BatchResult, len(reverts))
		for _, r := range reverts {
			results[r.MessageId] = contractclient.BatchResult{
				Receipt: &ethTypes.Receipt{Status: 1, TxHash: common.HexToHash("0x1234")},
			}
		}
		return results, nil
	}
}

func TestHandleEnygmaCrossTransfer(t *testing.T) { //nolint:gocognit // comprehensive test with many assertions
	t.Run("Early exit cases", func(t *testing.T) {
		t.Run("Sender is receiver - should skip processing", func(t *testing.T) {
			deps := setupReceiverDeps()

			batch := createBatch(big.NewInt(1), big.NewInt(1), 1)
			rcvr := setupReceiver(deps)

			err := rcvr.HandleEnygmaCrossTransfer(context.Background(), batch)

			assert.NoError(t, err)
			assert.Equal(t, 1, len(deps.Tracer.StartCalls.Calls))
			// Verify unused mocks were never called (early exit case)
			assert.Equal(t, 0, len(deps.EndpointClient.GetResourceAddressCalls()))
			assert.Equal(t, 0, len(deps.Deployer.DeployCalls()))
			assert.Equal(t, 0, len(deps.CreationServiceMock.CreateEnygmaCalls()))
			assert.Equal(t, 0, len(deps.HistoryRepository.GetEnygmaHistoryByUniqueKeyCalls()))
			assert.Equal(t, 0, len(deps.HandlerClient.ReceiveDestTransferBatchCalls()))
			assert.Equal(t, 0, len(deps.HandlerClient.RevertDestTransferBatchCalls()))
			assert.Equal(t, 0, len(deps.TeleportClient.SendTransferCompletedCalls()))
			assert.Equal(t, 0, len(deps.Simulator.GetRevertReasonCalls()))
		})

		t.Run("Invalid resource ID format - StringToBytes32 fails", func(t *testing.T) {
			deps := setupReceiverDeps()

			batch := createBatch(big.NewInt(2), big.NewInt(1), 1)
			batch.ResourceId = "invalid_resource_id"
			rcvr := setupReceiver(deps)

			err := rcvr.HandleEnygmaCrossTransfer(context.Background(), batch)

			assert.Error(t, err)
			// Verify unused mocks were never called (invalid resource ID - early exit)
			assert.Equal(t, 0, len(deps.EndpointClient.GetResourceAddressCalls()))
			assert.Equal(t, 0, len(deps.Deployer.DeployCalls()))
			assert.Equal(t, 0, len(deps.CreationServiceMock.CreateEnygmaCalls()))
			assert.Equal(t, 0, len(deps.HistoryRepository.GetEnygmaHistoryByUniqueKeyCalls()))
			assert.Equal(t, 0, len(deps.HandlerClient.ReceiveDestTransferBatchCalls()))
			assert.Equal(t, 0, len(deps.HandlerClient.RevertDestTransferBatchCalls()))
			assert.Equal(t, 0, len(deps.TeleportClient.SendTransferCompletedCalls()))
			assert.Equal(t, 0, len(deps.Simulator.GetRevertReasonCalls()))
		})

		t.Run("Batch already processed - tx_recovery is Confirmed, skip mint but re-insert history for resync", func(t *testing.T) {
			batch := createBatch(big.NewInt(2), big.NewInt(1), 1)

			deps := setupReceiverDeps()
			deps.EndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
				return common.HexToAddress("0x3333333333333333333333333333333333333333"), nil
			}
			deps.TxRecoveryRepository.GetByPrivateHubTxHashFunc = func(ctx context.Context, privateHubTxHash string) (*types.TxRecoveryData, error) {
				return &types.TxRecoveryData{
					PrivateHubTxHash: privateHubTxHash,
					Status:           types.HistoryStatusConfirmed,
				}, nil
			}
			deps.HistoryRepository.InsertEnygmaHistoryFunc = func(ctx context.Context, history types.EnygmaHistory) error {
				return nil
			}

			rcvr := setupReceiver(deps)

			err := rcvr.HandleEnygmaCrossTransfer(context.Background(), batch)

			assert.NoError(t, err)
			assert.Equal(t, 1, len(deps.EndpointClient.GetResourceAddressCalls()))
			// History IS re-inserted when recovery is confirmed — the resync flow
			// deletes history rows after rolling back checkpoints and expects the
			// receiver to recreate them.
			assert.Equal(t, 1, len(deps.HistoryRepository.InsertEnygmaHistoryCalls()))
			// Mint must NOT be redispatched (the original mint already landed on chain).
			assert.Equal(t, 0, len(deps.Deployer.DeployCalls()))
			assert.Equal(t, 0, len(deps.CreationServiceMock.CreateEnygmaCalls()))
			assert.Equal(t, 0, len(deps.HandlerClient.ReceiveDestTransferBatchCalls()))
			assert.Equal(t, 0, len(deps.HandlerClient.RevertDestTransferBatchCalls()))
			assert.Equal(t, 0, len(deps.TeleportClient.SendTransferCompletedCalls()))
			assert.Equal(t, 0, len(deps.Simulator.GetRevertReasonCalls()))
		})
	})

	t.Run("Token deployment", func(t *testing.T) {
		t.Run("First time deployment succeeds", func(t *testing.T) {
			deps := setupReceiverDeps()
			deps.EndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
				return common.HexToAddress("0x"), nil
			}
			deps.Deployer.DeployFunc = func(ctx context.Context, resourceId [32]byte, initiatorChainId *big.Int) (common.Address, error) {
				return common.HexToAddress("0x4444444444444444444444444444444444444444"), nil
			}
			deps.HistoryRepository.InsertEnygmaHistoryFunc = func(ctx context.Context, history types.EnygmaHistory) error {
				return nil
			}
			deps.HandlerClient.ReceiveDestTransferBatchFunc = makeSuccessfulMint()
			deps.TeleportClient.SendTransferCompletedFunc = func(ctx context.Context, messages []types.EnygmaTransferCompleted) error {
				return nil
			}

			batch := createBatch(big.NewInt(2), big.NewInt(1), 1)
			rcvr := setupReceiver(deps)

			err := rcvr.HandleEnygmaCrossTransfer(context.Background(), batch)

			assert.NoError(t, err)
			assert.Equal(t, 1, len(deps.Deployer.DeployCalls()))
			assert.Equal(t, 1, len(deps.CreationServiceMock.CreateEnygmaCalls()))
			// Verify unused mocks were never called
			assert.Equal(t, 0, len(deps.HandlerClient.RevertDestTransferBatchCalls()))
			assert.Equal(t, 0, len(deps.Simulator.GetRevertReasonCalls()))
		})

		t.Run("First time deployment fails", func(t *testing.T) {
			deps := setupReceiverDeps()
			deps.EndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
				return common.HexToAddress("0x"), nil
			}
			deps.Deployer.DeployFunc = func(ctx context.Context, resourceId [32]byte, initiatorChainId *big.Int) (common.Address, error) {
				return common.Address{}, errors.New("deployment failed")
			}

			batch := createBatch(big.NewInt(2), big.NewInt(1), 1)
			rcvr := setupReceiver(deps)

			err := rcvr.HandleEnygmaCrossTransfer(context.Background(), batch)

			assert.Error(t, err)
			assert.Equal(t, 1, len(deps.Deployer.DeployCalls()))
			// Verify unused mocks were never called
			assert.Equal(t, 0, len(deps.CreationServiceMock.CreateEnygmaCalls()))
			assert.Equal(t, 0, len(deps.HistoryRepository.GetEnygmaHistoryByUniqueKeyCalls()))
			assert.Equal(t, 0, len(deps.HistoryRepository.InsertEnygmaHistoryCalls()))
			assert.Equal(t, 0, len(deps.HandlerClient.ReceiveDestTransferBatchCalls()))
			assert.Equal(t, 0, len(deps.HandlerClient.RevertDestTransferBatchCalls()))
			assert.Equal(t, 0, len(deps.TeleportClient.SendTransferCompletedCalls()))
			assert.Equal(t, 0, len(deps.Simulator.GetRevertReasonCalls()))
		})

		t.Run("Token already exists", func(t *testing.T) {
			deps := setupReceiverDeps()
			deps.EndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
				return common.HexToAddress("0x5555555555555555555555555555555555555555"), nil
			}
			deps.HistoryRepository.InsertEnygmaHistoryFunc = func(ctx context.Context, history types.EnygmaHistory) error {
				return nil
			}
			deps.HandlerClient.ReceiveDestTransferBatchFunc = makeSuccessfulMint()
			deps.TeleportClient.SendTransferCompletedFunc = func(ctx context.Context, messages []types.EnygmaTransferCompleted) error {
				return nil
			}

			batch := createBatch(big.NewInt(2), big.NewInt(1), 1)
			rcvr := setupReceiver(deps)

			err := rcvr.HandleEnygmaCrossTransfer(context.Background(), batch)

			assert.NoError(t, err)
			// Verify unused mocks were never called
			assert.Equal(t, 0, len(deps.Deployer.DeployCalls()))
			assert.Equal(t, 0, len(deps.CreationServiceMock.CreateEnygmaCalls()))
			assert.Equal(t, 0, len(deps.HandlerClient.RevertDestTransferBatchCalls()))
			assert.Equal(t, 0, len(deps.Simulator.GetRevertReasonCalls()))
		})
	})

	t.Run("Minting operations", func(t *testing.T) {
		t.Run("All mints succeed", func(t *testing.T) {
			deps := setupReceiverDeps()
			deps.EndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
				return common.HexToAddress("0x6666666666666666666666666666666666666666"), nil
			}
			deps.HistoryRepository.InsertEnygmaHistoryFunc = func(ctx context.Context, history types.EnygmaHistory) error {
				return nil
			}
			deps.HandlerClient.ReceiveDestTransferBatchFunc = makeSuccessfulMint()
			deps.TeleportClient.SendTransferCompletedFunc = func(ctx context.Context, messages []types.EnygmaTransferCompleted) error {
				assert.Equal(t, 3, len(messages))
				return nil
			}

			batch := createBatch(big.NewInt(2), big.NewInt(1), 3)
			batch.PrivateHubTxHash = "0xdeadbeef"
			rcvr := setupReceiver(deps)

			err := rcvr.HandleEnygmaCrossTransfer(context.Background(), batch)

			assert.NoError(t, err)
			receiveCalls := deps.HandlerClient.ReceiveDestTransferBatchCalls()
			assert.Equal(t, 1, len(receiveCalls))
			// Every recipient's ProgramData blobs must reach the handler client verbatim —
			// they are the argument to the on-chain executeProgramData dispatch.
			assert.Equal(t, 3, len(receiveCalls[0].Transfers))
			for _, transfer := range receiveCalls[0].Transfers {
				assert.NotEmpty(t, transfer.ProgramData, "ProgramData must be forwarded to the executor")
			}
			assert.Equal(t, 1, len(deps.TeleportClient.SendTransferCompletedCalls()))
			// Verify unused mocks were never called
			assert.Equal(t, 0, len(deps.Deployer.DeployCalls()))
			assert.Equal(t, 0, len(deps.CreationServiceMock.CreateEnygmaCalls()))
			assert.Equal(t, 0, len(deps.HandlerClient.RevertDestTransferBatchCalls()))
			assert.Equal(t, 0, len(deps.Simulator.GetRevertReasonCalls()))
		})

		t.Run("Mixed success and failure in mints", func(t *testing.T) {
			deps := setupReceiverDeps()
			deps.EndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
				return common.HexToAddress("0x7777777777777777777777777777777777777777"), nil
			}
			deps.HistoryRepository.InsertEnygmaHistoryFunc = func(ctx context.Context, history types.EnygmaHistory) error {
				return nil
			}
			deps.HandlerClient.ReceiveDestTransferBatchFunc = makeAlternatingMint()
			deps.HandlerClient.RevertDestTransferBatchFunc = makeSuccessfulRevert()
			deps.TeleportClient.SendTransferCompletedFunc = func(ctx context.Context, messages []types.EnygmaTransferCompleted) error {
				return nil
			}

			batch := createBatch(big.NewInt(2), big.NewInt(1), 2)
			batch.PrivateHubTxHash = "0xdeadbeef"
			rcvr := setupReceiver(deps)

			err := rcvr.HandleEnygmaCrossTransfer(context.Background(), batch)

			assert.NoError(t, err)
			assert.Equal(t, 1, len(deps.HandlerClient.ReceiveDestTransferBatchCalls()))
			assert.Equal(t, 1, len(deps.HandlerClient.RevertDestTransferBatchCalls()))
			// Verify unused mocks were never called
			assert.Equal(t, 0, len(deps.Deployer.DeployCalls()))
			assert.Equal(t, 0, len(deps.CreationServiceMock.CreateEnygmaCalls()))
		})

		t.Run("All mints fail", func(t *testing.T) {
			deps := setupReceiverDeps()
			deps.EndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
				return common.HexToAddress("0x8888888888888888888888888888888888888888"), nil
			}
			deps.HistoryRepository.InsertEnygmaHistoryFunc = func(ctx context.Context, history types.EnygmaHistory) error {
				return nil
			}
			deps.HandlerClient.ReceiveDestTransferBatchFunc = makeFailedMint()
			deps.HandlerClient.RevertDestTransferBatchFunc = func(ctx context.Context, tokenAddress common.Address, reverts []*types.EnygmaCrossTransferData) (map[string]contractclient.BatchResult, error) {
				results := make(map[string]contractclient.BatchResult, len(reverts))
				for _, r := range reverts {
					results[r.MessageId] = contractclient.BatchResult{
						Err: errors.New("send failed"),
					}
				}
				return results, nil
			}

			batch := createBatch(big.NewInt(2), big.NewInt(1), 2)
			batch.PrivateHubTxHash = "0xdeadbeef"
			rcvr := setupReceiver(deps)

			err := rcvr.HandleEnygmaCrossTransfer(context.Background(), batch)

			assert.NoError(t, err)
			assert.Equal(t, 1, len(deps.HandlerClient.ReceiveDestTransferBatchCalls()))
			assert.Equal(t, 1, len(deps.HandlerClient.RevertDestTransferBatchCalls()))
			// All mints failed: no governance notification.
			assert.Equal(t, 0, len(deps.TeleportClient.SendTransferCompletedCalls()))
			// Verify unused mocks were never called
			assert.Equal(t, 0, len(deps.Deployer.DeployCalls()))
			assert.Equal(t, 0, len(deps.CreationServiceMock.CreateEnygmaCalls()))
		})

		t.Run("ReceiveDestTransferBatch returns error", func(t *testing.T) {
			deps := setupReceiverDeps()
			deps.EndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
				return common.HexToAddress("0x9999999999999999999999999999999999999999"), nil
			}
			deps.HistoryRepository.InsertEnygmaHistoryFunc = func(ctx context.Context, history types.EnygmaHistory) error {
				return nil
			}
			deps.HandlerClient.ReceiveDestTransferBatchFunc = func(ctx context.Context, transfers []*types.EnygmaCrossTransferData) (map[string]contractclient.BatchResult, error) {
				return nil, errors.New("handler client failure")
			}

			batch := createBatch(big.NewInt(2), big.NewInt(1), 1)
			rcvr := setupReceiver(deps)

			err := rcvr.HandleEnygmaCrossTransfer(context.Background(), batch)

			assert.Error(t, err)
			// Verify unused mocks were never called
			assert.Equal(t, 0, len(deps.Deployer.DeployCalls()))
			assert.Equal(t, 0, len(deps.CreationServiceMock.CreateEnygmaCalls()))
			assert.Equal(t, 0, len(deps.HandlerClient.RevertDestTransferBatchCalls()))
			assert.Equal(t, 0, len(deps.TeleportClient.SendTransferCompletedCalls()))
			assert.Equal(t, 0, len(deps.Simulator.GetRevertReasonCalls()))
		})
	})

	t.Run("Governance notification", func(t *testing.T) {
		t.Run("Sends transfer completed successfully", func(t *testing.T) {
			deps := setupReceiverDeps()
			deps.EndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
				return common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), nil
			}
			deps.HistoryRepository.InsertEnygmaHistoryFunc = func(ctx context.Context, history types.EnygmaHistory) error {
				return nil
			}
			deps.HandlerClient.ReceiveDestTransferBatchFunc = makeSuccessfulMint()
			deps.TeleportClient.SendTransferCompletedFunc = func(ctx context.Context, messages []types.EnygmaTransferCompleted) error {
				assert.Equal(t, 1, len(messages))
				// Message ID is dynamically generated - just verify we got a message
				assert.NotEmpty(t, messages[0].MessageId)
				return nil
			}

			batch := createBatch(big.NewInt(2), big.NewInt(1), 1)
			rcvr := setupReceiver(deps)

			err := rcvr.HandleEnygmaCrossTransfer(context.Background(), batch)

			assert.NoError(t, err)
			// Verify unused mocks were never called
			assert.Equal(t, 0, len(deps.Deployer.DeployCalls()))
			assert.Equal(t, 0, len(deps.CreationServiceMock.CreateEnygmaCalls()))
			assert.Equal(t, 0, len(deps.HandlerClient.RevertDestTransferBatchCalls()))
			assert.Equal(t, 0, len(deps.Simulator.GetRevertReasonCalls()))
		})

		t.Run("Fails to send transfer completed", func(t *testing.T) {
			deps := setupReceiverDeps()
			deps.EndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
				return common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"), nil
			}
			deps.HistoryRepository.InsertEnygmaHistoryFunc = func(ctx context.Context, history types.EnygmaHistory) error {
				return nil
			}
			deps.HandlerClient.ReceiveDestTransferBatchFunc = makeSuccessfulMint()
			deps.TeleportClient.SendTransferCompletedFunc = func(ctx context.Context, messages []types.EnygmaTransferCompleted) error {
				return errors.New("teleport failed")
			}

			batch := createBatch(big.NewInt(2), big.NewInt(1), 1)
			rcvr := setupReceiver(deps)

			err := rcvr.HandleEnygmaCrossTransfer(context.Background(), batch)

			assert.Error(t, err)
			// Verify unused mocks were never called
			assert.Equal(t, 0, len(deps.Deployer.DeployCalls()))
			assert.Equal(t, 0, len(deps.CreationServiceMock.CreateEnygmaCalls()))
			assert.Equal(t, 0, len(deps.HandlerClient.RevertDestTransferBatchCalls()))
			assert.Equal(t, 0, len(deps.Simulator.GetRevertReasonCalls()))
		})
	})

	t.Run("Revert operations", func(t *testing.T) {
		t.Run("Successful reverts after failed mints", func(t *testing.T) {
			deps := setupReceiverDeps()
			deps.EndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
				return common.HexToAddress("0xcccccccccccccccccccccccccccccccccccccccc"), nil
			}
			deps.HistoryRepository.InsertEnygmaHistoryFunc = func(ctx context.Context, history types.EnygmaHistory) error {
				return nil
			}
			deps.HandlerClient.ReceiveDestTransferBatchFunc = makeFailedMint()
			deps.HandlerClient.RevertDestTransferBatchFunc = makeSuccessfulRevert()

			batch := createBatch(big.NewInt(2), big.NewInt(1), 1)
			rcvr := setupReceiver(deps)

			err := rcvr.HandleEnygmaCrossTransfer(context.Background(), batch)

			assert.NoError(t, err)
			assert.Equal(t, 1, len(deps.HandlerClient.RevertDestTransferBatchCalls()))
			// All mints failed -> no governance notification.
			assert.Equal(t, 0, len(deps.TeleportClient.SendTransferCompletedCalls()))
			// Verify unused mocks were never called
			assert.Equal(t, 0, len(deps.Deployer.DeployCalls()))
			assert.Equal(t, 0, len(deps.CreationServiceMock.CreateEnygmaCalls()))
			assert.Equal(t, 0, len(deps.Simulator.GetRevertReasonCalls()))
		})

		t.Run("RevertDestTransferBatch returns error", func(t *testing.T) {
			deps := setupReceiverDeps()
			deps.EndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
				return common.HexToAddress("0xdddddddddddddddddddddddddddddddddddddddd"), nil
			}
			deps.HistoryRepository.InsertEnygmaHistoryFunc = func(ctx context.Context, history types.EnygmaHistory) error {
				return nil
			}
			deps.HandlerClient.ReceiveDestTransferBatchFunc = makeRevertedMint()
			deps.HandlerClient.RevertDestTransferBatchFunc = func(ctx context.Context, tokenAddress common.Address, reverts []*types.EnygmaCrossTransferData) (map[string]contractclient.BatchResult, error) {
				return nil, errors.New("revert batcher error")
			}

			batch := createBatch(big.NewInt(2), big.NewInt(1), 1)
			rcvr := setupReceiver(deps)

			err := rcvr.HandleEnygmaCrossTransfer(context.Background(), batch)

			assert.Error(t, err)
			// Verify unused mocks were never called
			assert.Equal(t, 0, len(deps.Deployer.DeployCalls()))
			assert.Equal(t, 0, len(deps.CreationServiceMock.CreateEnygmaCalls()))
			assert.Equal(t, 0, len(deps.TeleportClient.SendTransferCompletedCalls()))
		})

		t.Run("Reverts fail on chain", func(t *testing.T) {
			deps := setupReceiverDeps()
			deps.EndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
				return common.HexToAddress("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"), nil
			}
			deps.HistoryRepository.InsertEnygmaHistoryFunc = func(ctx context.Context, history types.EnygmaHistory) error {
				return nil
			}
			deps.HandlerClient.ReceiveDestTransferBatchFunc = makeRevertedMint()
			deps.HandlerClient.RevertDestTransferBatchFunc = func(ctx context.Context, tokenAddress common.Address, reverts []*types.EnygmaCrossTransferData) (map[string]contractclient.BatchResult, error) {
				results := make(map[string]contractclient.BatchResult, len(reverts))
				for _, r := range reverts {
					results[r.MessageId] = contractclient.BatchResult{
						Receipt: &ethTypes.Receipt{Status: 0, TxHash: common.HexToHash("0x5678")},
					}
				}
				return results, nil
			}

			batch := createBatch(big.NewInt(2), big.NewInt(1), 1)
			rcvr := setupReceiver(deps)

			err := rcvr.HandleEnygmaCrossTransfer(context.Background(), batch)

			assert.NoError(t, err)
			// Verify unused mocks were never called
			assert.Equal(t, 0, len(deps.Deployer.DeployCalls()))
			assert.Equal(t, 0, len(deps.CreationServiceMock.CreateEnygmaCalls()))
			assert.Equal(t, 0, len(deps.TeleportClient.SendTransferCompletedCalls()))
		})
	})

	t.Run("Error handling", func(t *testing.T) {
		t.Run("Cannot retrieve resource address", func(t *testing.T) {
			deps := setupReceiverDeps()
			deps.EndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
				return common.Address{}, errors.New("endpoint error")
			}

			batch := createBatch(big.NewInt(2), big.NewInt(1), 1)
			rcvr := setupReceiver(deps)

			err := rcvr.HandleEnygmaCrossTransfer(context.Background(), batch)

			assert.Error(t, err)
			// Verify unused mocks were never called
			assert.Equal(t, 0, len(deps.Deployer.DeployCalls()))
			assert.Equal(t, 0, len(deps.CreationServiceMock.CreateEnygmaCalls()))
			assert.Equal(t, 0, len(deps.HistoryRepository.GetEnygmaHistoryByUniqueKeyCalls()))
			assert.Equal(t, 0, len(deps.HistoryRepository.InsertEnygmaHistoryCalls()))
			assert.Equal(t, 0, len(deps.HandlerClient.ReceiveDestTransferBatchCalls()))
			assert.Equal(t, 0, len(deps.HandlerClient.RevertDestTransferBatchCalls()))
			assert.Equal(t, 0, len(deps.TeleportClient.SendTransferCompletedCalls()))
			assert.Equal(t, 0, len(deps.Simulator.GetRevertReasonCalls()))
		})

		t.Run("Database fails to insert history", func(t *testing.T) {
			deps := setupReceiverDeps()
			deps.EndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
				return common.HexToAddress("0xffffffffffffffffffffffffffffffffffffffff"), nil
			}
			deps.HistoryRepository.InsertEnygmaHistoryFunc = func(ctx context.Context, history types.EnygmaHistory) error {
				return errors.New("insert failed")
			}

			batch := createBatch(big.NewInt(2), big.NewInt(1), 1)
			rcvr := setupReceiver(deps)

			err := rcvr.HandleEnygmaCrossTransfer(context.Background(), batch)

			assert.Error(t, err)
			assert.Equal(t, 1, len(deps.HistoryRepository.InsertEnygmaHistoryCalls()))
			// Verify unused mocks were never called - mint not attempted on history insert failure.
			assert.Equal(t, 0, len(deps.HandlerClient.ReceiveDestTransferBatchCalls()))
			assert.Equal(t, 0, len(deps.Deployer.DeployCalls()))
			assert.Equal(t, 0, len(deps.CreationServiceMock.CreateEnygmaCalls()))
			assert.Equal(t, 0, len(deps.HandlerClient.RevertDestTransferBatchCalls()))
			assert.Equal(t, 0, len(deps.TeleportClient.SendTransferCompletedCalls()))
			assert.Equal(t, 0, len(deps.Simulator.GetRevertReasonCalls()))
		})

		t.Run("Revert reason simulator unused on success path", func(t *testing.T) {
			deps := setupReceiverDeps()
			deps.EndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
				return common.HexToAddress("0x1010101010101010101010101010101010101010"), nil
			}
			deps.HistoryRepository.InsertEnygmaHistoryFunc = func(ctx context.Context, history types.EnygmaHistory) error {
				return nil
			}
			deps.HandlerClient.ReceiveDestTransferBatchFunc = makeRevertedMint()
			deps.Simulator.GetRevertReasonFunc = func(ctx context.Context, txHash common.Hash) (txsim.ContractError, error) {
				return txsim.ContractError{}, errors.New("simulation failed")
			}
			deps.HandlerClient.RevertDestTransferBatchFunc = makeSuccessfulRevert()

			batch := createBatch(big.NewInt(2), big.NewInt(1), 1)
			rcvr := setupReceiver(deps)

			err := rcvr.HandleEnygmaCrossTransfer(context.Background(), batch)

			assert.NoError(t, err)
			// Verify unused mocks were never called
			assert.Equal(t, 0, len(deps.Deployer.DeployCalls()))
			assert.Equal(t, 0, len(deps.CreationServiceMock.CreateEnygmaCalls()))
			assert.Equal(t, 0, len(deps.TeleportClient.SendTransferCompletedCalls()))
		})
	})
}
