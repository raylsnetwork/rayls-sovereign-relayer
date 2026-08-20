package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/msgqueue"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/source/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools/fake"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools/spy"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newSingleMessageDvpMQ creates a mock message queue that returns a message only once,
// then blocks until context is canceled on subsequent calls
func newSingleMessageDvpMQ(msg msgqueue.Message[service.DvpSerializedEventBatch]) *DvpBatchMQMock {
	return &DvpBatchMQMock{
		NextFunc: fake.NextMQ(msg),
	}
}

// newDefaultDvpInitiatorMock creates a mock that fails if any handler is unexpectedly called
func newDefaultDvpInitiatorMock(t *testing.T) *DvpInitiatorMock {
	return &DvpInitiatorMock{
		HandleEnygmaSwapERC721Func: func(context.Context, string, *big.Int, common.Address, string, *big.Int, string, string, string, *big.Int, uint64) error {
			assert.Fail(t, "shouldn't call HandleEnygmaSwapERC721")
			return nil
		},
		HandleEnygmaSwapERC1155Func: func(context.Context, string, *big.Int, common.Address, string, *big.Int, string, string, *big.Int, string, *big.Int, uint64) error {
			assert.Fail(t, "shouldn't call HandleEnygmaSwapERC1155")
			return nil
		},
		HandleERC721CreationFunc: func(context.Context, string, string) error {
			assert.Fail(t, "shouldn't call HandleERC721Creation")
			return nil
		},
		HandleERC721MintFunc: func(context.Context, string, string, *big.Int) error {
			assert.Fail(t, "shouldn't call HandleERC721Mint")
			return nil
		},
		HandleERC721BurnFunc: func(context.Context, string, string, *big.Int) error {
			assert.Fail(t, "shouldn't call HandleERC721Burn")
			return nil
		},
		HandleERC721DepositFunc: func(context.Context, string, string, *big.Int, common.Address, string, *big.Int) error {
			assert.Fail(t, "shouldn't call HandleERC721Deposit")
			return nil
		},
		HandleERC721WithdrawFunc: func(context.Context, string, string, *big.Int, common.Address, string, *big.Int) error {
			assert.Fail(t, "shouldn't call HandleERC721Withdraw")
			return nil
		},
		HandleERC721SwapEnygmaFunc: func(context.Context, string, *big.Int, common.Address, string, string, string, *big.Int, string, *big.Int, uint64) error {
			assert.Fail(t, "shouldn't call HandleERC721SwapEnygma")
			return nil
		},
		HandleERC1155CreationFunc: func(context.Context, string, string) error {
			assert.Fail(t, "shouldn't call HandleERC1155Creation")
			return nil
		},
		HandleERC1155MintFunc: func(context.Context, string, string, *big.Int, *big.Int, []byte) error {
			assert.Fail(t, "shouldn't call HandleERC1155Mint")
			return nil
		},
		HandleERC1155BurnFunc: func(context.Context, string, string, *big.Int, *big.Int) error {
			assert.Fail(t, "shouldn't call HandleERC1155Burn")
			return nil
		},
		HandleERC1155DepositFunc: func(context.Context, string, string, *big.Int, *big.Int, []byte, common.Address, string, *big.Int) error {
			assert.Fail(t, "shouldn't call HandleERC1155Deposit")
			return nil
		},
		HandleERC1155WithdrawFunc: func(context.Context, string, string, *big.Int, *big.Int, common.Address, string, *big.Int) error {
			assert.Fail(t, "shouldn't call HandleERC1155Withdraw")
			return nil
		},
		HandleERC1155SwapEnygmaFunc: func(context.Context, string, *big.Int, common.Address, string, *big.Int, string, string, *big.Int, string, *big.Int, uint64) error {
			assert.Fail(t, "shouldn't call HandleERC1155SwapEnygma")
			return nil
		},
	}
}

func TestDvpOrchestrator(t *testing.T) {
	t.Run("calls initiator handler on Enygma swap ERC721 event batch", func(t *testing.T) {
		wantSharedId := "test-shared-id"
		wantDestChainId := big.NewInt(100)
		wantFrom := common.HexToAddress("0xc001babe")
		wantResourceId := "resource-id"
		wantEnygmaAmount := big.NewInt(1000)
		wantNftResourceId := "nft-resource-id"
		wantNftId := "42"

		ackSpy := spy.NewAck()

		events := []service.DvpEnygmaSwapERC721{
			{
				SharedId:      wantSharedId,
				DestChainId:   (wantDestChainId),
				From:          wantFrom,
				ResourceId:    wantResourceId,
				EnygmaAmount:  (wantEnygmaAmount),
				NftResourceId: wantNftResourceId,
				NftId:         wantNftId,
			},
		}
		serializedEvents, _ := json.Marshal(events)

		dvpMQ := newSingleMessageDvpMQ(msgqueue.Message[service.DvpSerializedEventBatch]{
			V: service.DvpSerializedEventBatch{
				BlockNumber:      100,
				Type:             service.DvpEnygmaSwapERC721Event,
				SerializedEvents: serializedEvents,
			},
			Ack: ackSpy.Fn(),
		})

		initiator := newDefaultDvpInitiatorMock(t)
		initiator.HandleEnygmaSwapERC721Func = func(ctx context.Context, sharedId string, toChainId *big.Int, from common.Address, enygmaResourceId string, enygmaAmount *big.Int, nftResourceId string, nftId string, txHash string, txBlockNumber *big.Int, validityTime uint64) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantSharedId, sharedId)
			assert.Equal(t, wantDestChainId, toChainId)
			assert.Equal(t, wantFrom, from)
			assert.Equal(t, wantResourceId, enygmaResourceId)
			assert.Equal(t, wantEnygmaAmount, enygmaAmount)
			assert.Equal(t, wantNftResourceId, nftResourceId)
			assert.Equal(t, wantNftId, nftId)
			return nil
		}

		svc := &service.DvpOrchestrator{}
		// Use reflection or create a proper constructor that accepts dependencies
		// For now, we'll need to expose fields or add a constructor
		svc = service.NewDvpOrchestratorWithDeps(dvpMQ, initiator)

		err := testtools.RunUntilAcked(t, svc.Run, ackSpy, 2*time.Second)
		require.NoError(t, err)

		assert.Equal(t, 1, len(initiator.HandleEnygmaSwapERC721Calls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("calls initiator handler on ERC721 creation event batch", func(t *testing.T) {
		wantResourceId := "test-resource-id"

		ackSpy := spy.NewAck()

		events := []service.Dvp721Creation{
			{
				ResourceId: wantResourceId,
			},
		}
		serializedEvents, _ := json.Marshal(events)

		dvpMQ := newSingleMessageDvpMQ(msgqueue.Message[service.DvpSerializedEventBatch]{
			V: service.DvpSerializedEventBatch{
				BlockNumber:      100,
				Type:             service.Dvp721CreationEvent,
				SerializedEvents: serializedEvents,
			},
			Ack: ackSpy.Fn(),
		})

		initiator := newDefaultDvpInitiatorMock(t)
		initiator.HandleERC721CreationFunc = func(ctx context.Context, _ string, resourceId string) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantResourceId, resourceId)
			return nil
		}

		svc := service.NewDvpOrchestratorWithDeps(dvpMQ, initiator)

		err := testtools.RunUntilAcked(t, svc.Run, ackSpy, 2*time.Second)
		require.NoError(t, err)

		assert.Equal(t, 1, len(initiator.HandleERC721CreationCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("calls initiator handler on ERC721 mint event batch", func(t *testing.T) {
		wantResourceId := "test-resource-id"
		wantNftId := big.NewInt(42)

		ackSpy := spy.NewAck()

		events := []service.Dvp721Mint{
			{
				ResourceId: wantResourceId,
				NftId:      wantNftId,
			},
		}
		serializedEvents, _ := json.Marshal(events)

		dvpMQ := newSingleMessageDvpMQ(msgqueue.Message[service.DvpSerializedEventBatch]{
			V: service.DvpSerializedEventBatch{
				BlockNumber:      100,
				Type:             service.Dvp721MintEvent,
				SerializedEvents: serializedEvents,
			},
			Ack: ackSpy.Fn(),
		})

		initiator := newDefaultDvpInitiatorMock(t)
		initiator.HandleERC721MintFunc = func(ctx context.Context, _ string, resourceId string, nftId *big.Int) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantResourceId, resourceId)
			assert.Equal(t, wantNftId, nftId)
			return nil
		}

		svc := service.NewDvpOrchestratorWithDeps(dvpMQ, initiator)

		err := testtools.RunUntilAcked(t, svc.Run, ackSpy, 2*time.Second)
		require.NoError(t, err)

		assert.Equal(t, 1, len(initiator.HandleERC721MintCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("calls initiator handler on ERC721 deposit event batch with multiple events", func(t *testing.T) {
		wantResourceId := "test-resource-id"
		wantNftId1 := big.NewInt(42)
		wantFrom1 := common.HexToAddress("0xc001babe")
		wantNftId2 := big.NewInt(43)
		wantFrom2 := common.HexToAddress("0xdeadc0de")

		ackSpy := spy.NewAck()

		events := []service.Dvp721DepositIntoDvp{
			{
				ResourceId: wantResourceId,
				NftId:      wantNftId1,
				From:       wantFrom1,
			},
			{
				ResourceId: wantResourceId,
				NftId:      wantNftId2,
				From:       wantFrom2,
			},
		}
		serializedEvents, _ := json.Marshal(events)

		dvpMQ := newSingleMessageDvpMQ(msgqueue.Message[service.DvpSerializedEventBatch]{
			V: service.DvpSerializedEventBatch{
				BlockNumber:      100,
				Type:             service.Dvp721DepositIntoDvpEvent,
				SerializedEvents: serializedEvents,
			},
			Ack: ackSpy.Fn(),
		})

		callCount := 0
		initiator := newDefaultDvpInitiatorMock(t)
		initiator.HandleERC721DepositFunc = func(ctx context.Context, _ string, resourceId string, nftId *big.Int, from common.Address, txHash string, txBlockNumber *big.Int) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling all events")
			assert.Equal(t, wantResourceId, resourceId)

			if callCount == 0 {
				assert.Equal(t, wantNftId1, nftId)
				assert.Equal(t, wantFrom1, from)
			} else if callCount == 1 {
				assert.Equal(t, wantNftId2, nftId)
				assert.Equal(t, wantFrom2, from)
			}
			callCount++
			return nil
		}

		svc := service.NewDvpOrchestratorWithDeps(dvpMQ, initiator)

		err := testtools.RunUntilAcked(t, svc.Run, ackSpy, 2*time.Second)
		require.NoError(t, err)

		assert.Equal(t, 2, len(initiator.HandleERC721DepositCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("calls initiator handler on ERC1155 creation event batch", func(t *testing.T) {
		wantResourceId := "test-resource-id"

		ackSpy := spy.NewAck()

		events := []service.Dvp1155Creation{
			{
				ResourceId: wantResourceId,
			},
		}
		serializedEvents, _ := json.Marshal(events)

		dvpMQ := newSingleMessageDvpMQ(msgqueue.Message[service.DvpSerializedEventBatch]{
			V: service.DvpSerializedEventBatch{
				BlockNumber:      100,
				Type:             service.Dvp1155CreationEvent,
				SerializedEvents: serializedEvents,
			},
			Ack: ackSpy.Fn(),
		})

		initiator := newDefaultDvpInitiatorMock(t)
		initiator.HandleERC1155CreationFunc = func(ctx context.Context, _ string, resourceId string) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantResourceId, resourceId)
			return nil
		}

		svc := service.NewDvpOrchestratorWithDeps(dvpMQ, initiator)

		err := testtools.RunUntilAcked(t, svc.Run, ackSpy, 2*time.Second)
		require.NoError(t, err)

		assert.Equal(t, 1, len(initiator.HandleERC1155CreationCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("calls initiator handler on ERC1155 mint event batch", func(t *testing.T) {
		wantResourceId := "test-resource-id"
		wantTokenId := big.NewInt(100)
		wantValue := big.NewInt(50)
		wantData := []byte{0x01, 0x02}

		ackSpy := spy.NewAck()

		events := []service.Dvp1155Mint{
			{
				ResourceId: wantResourceId,
				TokenId:    wantTokenId,
				Value:      wantValue,
				Data:       wantData,
			},
		}
		serializedEvents, _ := json.Marshal(events)

		dvpMQ := newSingleMessageDvpMQ(msgqueue.Message[service.DvpSerializedEventBatch]{
			V: service.DvpSerializedEventBatch{
				BlockNumber:      100,
				Type:             service.Dvp1155MintEvent,
				SerializedEvents: serializedEvents,
			},
			Ack: ackSpy.Fn(),
		})

		initiator := newDefaultDvpInitiatorMock(t)
		initiator.HandleERC1155MintFunc = func(ctx context.Context, _ string, resourceId string, tokenId *big.Int, tokenAmount *big.Int, tokenData []byte) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantResourceId, resourceId)
			assert.Equal(t, wantTokenId, tokenId)
			assert.Equal(t, wantValue, tokenAmount)
			assert.Equal(t, wantData, tokenData)
			return nil
		}

		svc := service.NewDvpOrchestratorWithDeps(dvpMQ, initiator)

		err := testtools.RunUntilAcked(t, svc.Run, ackSpy, 2*time.Second)
		require.NoError(t, err)

		assert.Equal(t, 1, len(initiator.HandleERC1155MintCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("calls initiator handler on Enygma swap ERC1155 event batch", func(t *testing.T) {
		wantSharedId := "test-shared-id"
		wantDestChainId := big.NewInt(200)
		wantFrom := common.HexToAddress("0xf00dbabe")
		wantResourceId := "resource-id"
		wantEnygmaAmount := big.NewInt(2000)
		wantNftResourceId := "nft-resource-id"
		wantNftId := "100"
		wantNftAmountOrOne := big.NewInt(5)

		ackSpy := spy.NewAck()

		events := []service.DvpEnygmaSwapERC1155{
			{
				SharedId:       wantSharedId,
				DestChainId:    (wantDestChainId),
				From:           wantFrom,
				ResourceId:     wantResourceId,
				EnygmaAmount:   (wantEnygmaAmount),
				NftResourceId:  wantNftResourceId,
				NftId:          wantNftId,
				NftAmountOrOne: (wantNftAmountOrOne),
			},
		}
		serializedEvents, _ := json.Marshal(events)

		dvpMQ := newSingleMessageDvpMQ(msgqueue.Message[service.DvpSerializedEventBatch]{
			V: service.DvpSerializedEventBatch{
				BlockNumber:      100,
				Type:             service.DvpEnygmaSwapERC1155Event,
				SerializedEvents: serializedEvents,
			},
			Ack: ackSpy.Fn(),
		})

		initiator := newDefaultDvpInitiatorMock(t)
		initiator.HandleEnygmaSwapERC1155Func = func(ctx context.Context, sharedId string, toChainId *big.Int, from common.Address, enygmaResourceId string, enygmaAmount *big.Int, nftResourceId string, nftId string, nftAmount *big.Int, txHash string, txBlockNumber *big.Int, validityTime uint64) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantSharedId, sharedId)
			assert.Equal(t, wantDestChainId, toChainId)
			assert.Equal(t, wantFrom, from)
			assert.Equal(t, wantResourceId, enygmaResourceId)
			assert.Equal(t, wantEnygmaAmount, enygmaAmount)
			assert.Equal(t, wantNftResourceId, nftResourceId)
			assert.Equal(t, wantNftId, nftId)
			assert.Equal(t, wantNftAmountOrOne, nftAmount)
			return nil
		}

		svc := service.NewDvpOrchestratorWithDeps(dvpMQ, initiator)

		err := testtools.RunUntilAcked(t, svc.Run, ackSpy, 2*time.Second)
		require.NoError(t, err)

		assert.Equal(t, 1, len(initiator.HandleEnygmaSwapERC1155Calls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("calls initiator handler on ERC721 burn event batch", func(t *testing.T) {
		wantResourceId := "test-resource-id"
		wantNftId := big.NewInt(99)

		ackSpy := spy.NewAck()

		events := []service.Dvp721Burn{
			{
				ResourceId: wantResourceId,
				NftId:      wantNftId,
			},
		}
		serializedEvents, _ := json.Marshal(events)

		dvpMQ := newSingleMessageDvpMQ(msgqueue.Message[service.DvpSerializedEventBatch]{
			V: service.DvpSerializedEventBatch{
				BlockNumber:      100,
				Type:             service.Dvp721BurnEvent,
				SerializedEvents: serializedEvents,
			},
			Ack: ackSpy.Fn(),
		})

		initiator := newDefaultDvpInitiatorMock(t)
		initiator.HandleERC721BurnFunc = func(ctx context.Context, _ string, resourceId string, nftId *big.Int) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantResourceId, resourceId)
			assert.Equal(t, wantNftId, nftId)
			return nil
		}

		svc := service.NewDvpOrchestratorWithDeps(dvpMQ, initiator)

		err := testtools.RunUntilAcked(t, svc.Run, ackSpy, 2*time.Second)
		require.NoError(t, err)

		assert.Equal(t, 1, len(initiator.HandleERC721BurnCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("calls initiator handler on ERC721 withdraw event batch", func(t *testing.T) {
		wantResourceId := "test-resource-id"
		wantNftId := big.NewInt(77)
		wantOwner := common.HexToAddress("0x0wn3r")

		ackSpy := spy.NewAck()

		events := []service.Dvp721WithdrawFromDvp{
			{
				ResourceId: wantResourceId,
				NftId:      wantNftId,
				Owner:      wantOwner,
			},
		}
		serializedEvents, _ := json.Marshal(events)

		dvpMQ := newSingleMessageDvpMQ(msgqueue.Message[service.DvpSerializedEventBatch]{
			V: service.DvpSerializedEventBatch{
				BlockNumber:      100,
				Type:             service.Dvp721WithdrawFromDvpEvent,
				SerializedEvents: serializedEvents,
			},
			Ack: ackSpy.Fn(),
		})

		initiator := newDefaultDvpInitiatorMock(t)
		initiator.HandleERC721WithdrawFunc = func(ctx context.Context, _ string, resourceId string, nftId *big.Int, from common.Address, txHash string, txBlockNumber *big.Int) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantResourceId, resourceId)
			assert.Equal(t, wantNftId, nftId)
			assert.Equal(t, wantOwner, from)
			return nil
		}

		svc := service.NewDvpOrchestratorWithDeps(dvpMQ, initiator)

		err := testtools.RunUntilAcked(t, svc.Run, ackSpy, 2*time.Second)
		require.NoError(t, err)

		assert.Equal(t, 1, len(initiator.HandleERC721WithdrawCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("calls initiator handler on ERC721 swap for Enygma event batch", func(t *testing.T) {
		wantSharedId := "shared-swap-id"
		wantDestChainId := big.NewInt(300)
		wantFrom := common.HexToAddress("0xabc123")
		wantNftResourceId := "nft-res-id"
		wantNftId := "55"
		wantEnygmaResourceId := "enygma-res-id"
		wantEnygmaAmount := big.NewInt(3000)

		ackSpy := spy.NewAck()

		events := []service.Dvp721SwapForEnygma{
			{
				SharedId:         wantSharedId,
				DestChainId:      wantDestChainId,
				From:             wantFrom,
				NftResourceId:    wantNftResourceId,
				NftId:            wantNftId,
				EnygmaResourceId: wantEnygmaResourceId,
				EnygmaAmount:     wantEnygmaAmount,
			},
		}
		serializedEvents, _ := json.Marshal(events)

		dvpMQ := newSingleMessageDvpMQ(msgqueue.Message[service.DvpSerializedEventBatch]{
			V: service.DvpSerializedEventBatch{
				BlockNumber:      100,
				Type:             service.Dvp721SwapForEnygmaEvent,
				SerializedEvents: serializedEvents,
			},
			Ack: ackSpy.Fn(),
		})

		initiator := newDefaultDvpInitiatorMock(t)
		initiator.HandleERC721SwapEnygmaFunc = func(ctx context.Context, sharedId string, toChainId *big.Int, from common.Address, nftResourceId string, nftId string, enygmaResourceId string, enygmaAmount *big.Int, txHash string, txBlockNumber *big.Int, validityTime uint64) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantSharedId, sharedId)
			assert.Equal(t, wantDestChainId, toChainId)
			assert.Equal(t, wantFrom, from)
			assert.Equal(t, wantNftResourceId, nftResourceId)
			assert.Equal(t, wantNftId, nftId)
			assert.Equal(t, wantEnygmaResourceId, enygmaResourceId)
			assert.Equal(t, wantEnygmaAmount, enygmaAmount)
			return nil
		}

		svc := service.NewDvpOrchestratorWithDeps(dvpMQ, initiator)

		err := testtools.RunUntilAcked(t, svc.Run, ackSpy, 2*time.Second)
		require.NoError(t, err)

		assert.Equal(t, 1, len(initiator.HandleERC721SwapEnygmaCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("calls initiator handler on ERC1155 burn event batch", func(t *testing.T) {
		wantResourceId := "test-resource-id"
		wantTokenId := big.NewInt(200)
		wantValue := big.NewInt(25)

		ackSpy := spy.NewAck()

		events := []service.Dvp1155Burn{
			{
				ResourceId: wantResourceId,
				TokenId:    wantTokenId,
				Value:      wantValue,
			},
		}
		serializedEvents, _ := json.Marshal(events)

		dvpMQ := newSingleMessageDvpMQ(msgqueue.Message[service.DvpSerializedEventBatch]{
			V: service.DvpSerializedEventBatch{
				BlockNumber:      100,
				Type:             service.Dvp1155BurnEvent,
				SerializedEvents: serializedEvents,
			},
			Ack: ackSpy.Fn(),
		})

		initiator := newDefaultDvpInitiatorMock(t)
		initiator.HandleERC1155BurnFunc = func(ctx context.Context, _ string, resourceId string, tokenId *big.Int, tokenAmount *big.Int) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantResourceId, resourceId)
			assert.Equal(t, wantTokenId, tokenId)
			assert.Equal(t, wantValue, tokenAmount)
			return nil
		}

		svc := service.NewDvpOrchestratorWithDeps(dvpMQ, initiator)

		err := testtools.RunUntilAcked(t, svc.Run, ackSpy, 2*time.Second)
		require.NoError(t, err)

		assert.Equal(t, 1, len(initiator.HandleERC1155BurnCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("calls initiator handler on ERC1155 deposit event batch", func(t *testing.T) {
		wantResourceId := "test-resource-id"
		wantTokenId := big.NewInt(300)
		wantValue := big.NewInt(75)
		wantData := []byte{0x09, 0x0a}
		wantFrom := common.HexToAddress("0xdeposit0r")

		ackSpy := spy.NewAck()

		events := []service.Dvp1155DepositIntoDvp{
			{
				ResourceId: wantResourceId,
				TokenId:    wantTokenId,
				Value:      wantValue,
				Data:       wantData,
				From:       wantFrom,
			},
		}
		serializedEvents, _ := json.Marshal(events)

		dvpMQ := newSingleMessageDvpMQ(msgqueue.Message[service.DvpSerializedEventBatch]{
			V: service.DvpSerializedEventBatch{
				BlockNumber:      100,
				Type:             service.Dvp1155DepositIntoDvpEvent,
				SerializedEvents: serializedEvents,
			},
			Ack: ackSpy.Fn(),
		})

		initiator := newDefaultDvpInitiatorMock(t)
		initiator.HandleERC1155DepositFunc = func(ctx context.Context, _ string, resourceId string, tokenId *big.Int, tokenAmount *big.Int, tokenData []byte, from common.Address, txHash string, txBlockNumber *big.Int) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantResourceId, resourceId)
			assert.Equal(t, wantTokenId, tokenId)
			assert.Equal(t, wantValue, tokenAmount)
			assert.Equal(t, wantData, tokenData)
			assert.Equal(t, wantFrom, from)
			return nil
		}

		svc := service.NewDvpOrchestratorWithDeps(dvpMQ, initiator)

		err := testtools.RunUntilAcked(t, svc.Run, ackSpy, 2*time.Second)
		require.NoError(t, err)

		assert.Equal(t, 1, len(initiator.HandleERC1155DepositCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("calls initiator handler on ERC1155 withdraw event batch", func(t *testing.T) {
		wantResourceId := "test-resource-id"
		wantTokenId := big.NewInt(400)
		wantValue := big.NewInt(10)
		wantOwner := common.HexToAddress("0xw1thdr4w")
		wantTxHash := "0xabc123"
		wantTxBlockNumber := big.NewInt(99)

		ackSpy := spy.NewAck()

		events := []service.Dvp1155WithdrawFromDvp{
			{
				ResourceId:    wantResourceId,
				TokenId:       wantTokenId,
				Value:         wantValue,
				Owner:         wantOwner,
				TxHash:        wantTxHash,
				TxBlockNumber: wantTxBlockNumber,
			},
		}
		serializedEvents, _ := json.Marshal(events)

		dvpMQ := newSingleMessageDvpMQ(msgqueue.Message[service.DvpSerializedEventBatch]{
			V: service.DvpSerializedEventBatch{
				BlockNumber:      100,
				Type:             service.Dvp1155WithdrawFromDvpEvent,
				SerializedEvents: serializedEvents,
			},
			Ack: ackSpy.Fn(),
		})

		initiator := newDefaultDvpInitiatorMock(t)
		initiator.HandleERC1155WithdrawFunc = func(ctx context.Context, _ string, resourceId string, tokenId *big.Int, tokenAmount *big.Int, from common.Address, txHash string, txBlockNumber *big.Int) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantResourceId, resourceId)
			assert.Equal(t, wantTokenId, tokenId)
			assert.Equal(t, wantValue, tokenAmount)
			assert.Equal(t, wantOwner, from)
			assert.Equal(t, wantTxHash, txHash)
			assert.Equal(t, wantTxBlockNumber, txBlockNumber)
			return nil
		}

		svc := service.NewDvpOrchestratorWithDeps(dvpMQ, initiator)

		err := testtools.RunUntilAcked(t, svc.Run, ackSpy, 2*time.Second)
		require.NoError(t, err)

		assert.Equal(t, 1, len(initiator.HandleERC1155WithdrawCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("calls initiator handler on ERC1155 swap for Enygma event batch", func(t *testing.T) {
		wantSharedId := "shared-swap-1155"
		wantDestChainId := big.NewInt(400)
		wantFrom := common.HexToAddress("0xswapp3r")
		wantTokenResourceId := "token-res-id"
		wantTokenValue := big.NewInt(15)
		wantTokenId := "88"
		wantEnygmaResourceId := "enygma-res-id"
		wantEnygmaAmount := big.NewInt(4000)

		ackSpy := spy.NewAck()

		events := []service.Dvp1155SwapForEnygma{
			{
				SharedId:         wantSharedId,
				DestChainId:      wantDestChainId,
				From:             wantFrom,
				TokenResourceId:  wantTokenResourceId,
				TokenValue:       wantTokenValue,
				TokenId:          wantTokenId,
				EnygmaResourceId: wantEnygmaResourceId,
				EnygmaAmount:     wantEnygmaAmount,
			},
		}
		serializedEvents, _ := json.Marshal(events)

		dvpMQ := newSingleMessageDvpMQ(msgqueue.Message[service.DvpSerializedEventBatch]{
			V: service.DvpSerializedEventBatch{
				BlockNumber:      100,
				Type:             service.Dvp1155SwapForEnygmaEvent,
				SerializedEvents: serializedEvents,
			},
			Ack: ackSpy.Fn(),
		})

		initiator := newDefaultDvpInitiatorMock(t)
		initiator.HandleERC1155SwapEnygmaFunc = func(ctx context.Context, sharedId string, toChainId *big.Int, from common.Address, erc1155ResourceId string, erc1155Amount *big.Int, erc1155Id string, enygmaResourceId string, enygmaAmount *big.Int, txHash string, txBlockNumber *big.Int, validityTime uint64) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantSharedId, sharedId)
			assert.Equal(t, wantDestChainId, toChainId)
			assert.Equal(t, wantFrom, from)
			assert.Equal(t, wantTokenResourceId, erc1155ResourceId)
			assert.Equal(t, wantTokenValue, erc1155Amount)
			assert.Equal(t, wantTokenId, erc1155Id)
			assert.Equal(t, wantEnygmaResourceId, enygmaResourceId)
			assert.Equal(t, wantEnygmaAmount, enygmaAmount)
			return nil
		}

		svc := service.NewDvpOrchestratorWithDeps(dvpMQ, initiator)

		err := testtools.RunUntilAcked(t, svc.Run, ackSpy, 2*time.Second)
		require.NoError(t, err)

		assert.Equal(t, 1, len(initiator.HandleERC1155SwapEnygmaCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("does not ack message when handler fails", func(t *testing.T) {
		ackSpy := spy.NewAck()

		events := []service.Dvp721Creation{
			{
				ResourceId: "test-resource-id",
			},
			{
				ResourceId: "second-resource-id",
			},
		}
		serializedEvents, _ := json.Marshal(events)

		dvpMQ := newSingleMessageDvpMQ(msgqueue.Message[service.DvpSerializedEventBatch]{
			V: service.DvpSerializedEventBatch{
				BlockNumber:      100,
				Type:             service.Dvp721CreationEvent,
				SerializedEvents: serializedEvents,
			},
			Ack: ackSpy.Fn(),
		})

		callCount := 0
		initiator := newDefaultDvpInitiatorMock(t)
		initiator.HandleERC721CreationFunc = func(ctx context.Context, _ string, resourceId string) error {
			callCount++
			if callCount == 1 {
				return errors.New("handler error")
			}
			return nil
		}

		svc := service.NewDvpOrchestratorWithDeps(dvpMQ, initiator)

		// Negative case: the assertion is the *absence* of an ack, which has
		// no synchronisation signal to wait for. Use a generous safety budget
		// to give Run a real chance to attempt processing; 500ms is well
		// above any expected scheduling jitter and keeps the test fast.
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 2, len(initiator.HandleERC721CreationCalls()))
		ackSpy.AssertNotCalled(t, "should not ack message when handler fails, allowing MQ redelivery")
	})

	t.Run("does not ack message when swap cancelled has invalid ERC standard", func(t *testing.T) {
		ackSpy := spy.NewAck()

		events := []service.DvpSwapCancelled{
			{
				SharedId:        "test-shared-id",
				DestChainId:     big.NewInt(100),
				TokenInStandard: 255, // invalid ERC standard
			},
		}
		serializedEvents, _ := json.Marshal(events)

		dvpMQ := newSingleMessageDvpMQ(msgqueue.Message[service.DvpSerializedEventBatch]{
			V: service.DvpSerializedEventBatch{
				BlockNumber:      100,
				Type:             service.DvpSwapCancelledEvent,
				SerializedEvents: serializedEvents,
			},
			Ack: ackSpy.Fn(),
		})

		initiator := newDefaultDvpInitiatorMock(t)

		svc := service.NewDvpOrchestratorWithDeps(dvpMQ, initiator)

		// Negative case: the assertion is the *absence* of an ack. See the
		// adjacent "does not ack message when handler fails" test for the
		// rationale on the 500ms safety budget.
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		err := svc.Run(ctx)
		require.NoError(t, err)

		ackSpy.AssertNotCalled(t, "should not ack message when ERC standard conversion fails, allowing MQ redelivery")
	})
}
