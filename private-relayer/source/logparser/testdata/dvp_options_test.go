package testdata

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/EnygmaPNEvents"
)

func TestDvp721CreationLogWithOptions(t *testing.T) {
	// Test with custom resource ID
	customResourceID := [32]byte{0xaa, 0xbb, 0xcc}
	log := NewDvp721CreationLogWith(
		WithDvp721CreationResourceID(customResourceID),
	)

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackDvp721CreationEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse Dvp721Creation log: %v", err)
	}

	if event.ResourceId != customResourceID {
		t.Errorf("Expected resourceId %v, got %v", customResourceID, event.ResourceId)
	}

	// Should still have default block number
	if log.BlockNumber != 105 {
		t.Errorf("Expected default block number 105, got %d", log.BlockNumber)
	}
}

func TestDvp721CreationLogWithMultipleOptions(t *testing.T) {
	// Test with multiple custom options
	customResourceID := [32]byte{0x11, 0x22, 0x33}
	customBlockNumber := uint64(999)

	log := NewDvp721CreationLogWith(
		WithDvp721CreationResourceID(customResourceID),
		WithDvp721CreationBlockNumber(customBlockNumber),
	)

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackDvp721CreationEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse Dvp721Creation log: %v", err)
	}

	if event.ResourceId != customResourceID {
		t.Errorf("Expected resourceId %v, got %v", customResourceID, event.ResourceId)
	}

	if log.BlockNumber != customBlockNumber {
		t.Errorf("Expected blockNumber %d, got %d", customBlockNumber, log.BlockNumber)
	}
}

func TestDvp721MintLogWithOptions(t *testing.T) {
	// Test with custom values
	customResourceID := [32]byte{0xdd, 0xee, 0xff}
	customNftId := big.NewInt(999)
	customBlockNumber := uint64(555)

	log := NewDvp721MintLogWith(
		WithDvp721MintResourceID(customResourceID),
		WithDvp721MintNftId(customNftId),
		WithDvp721MintBlockNumber(customBlockNumber),
	)

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackDvp721MintEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse Dvp721Mint log: %v", err)
	}

	if event.ResourceId != customResourceID {
		t.Errorf("Expected resourceId %v, got %v", customResourceID, event.ResourceId)
	}

	if event.NftId.Cmp(customNftId) != 0 {
		t.Errorf("Expected nftId %v, got %v", customNftId, event.NftId)
	}

	if log.BlockNumber != customBlockNumber {
		t.Errorf("Expected blockNumber %d, got %d", customBlockNumber, log.BlockNumber)
	}
}

func TestDvp721BurnLogWithOptions(t *testing.T) {
	// Test with custom values
	customResourceID := [32]byte{0x44, 0x55, 0x66}
	customNftId := big.NewInt(777)
	customBlockNumber := uint64(888)

	log := NewDvp721BurnLogWith(
		WithDvp721BurnResourceID(customResourceID),
		WithDvp721BurnNftId(customNftId),
		WithDvp721BurnBlockNumber(customBlockNumber),
	)

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackDvp721BurnEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse Dvp721Burn log: %v", err)
	}

	if event.ResourceId != customResourceID {
		t.Errorf("Expected resourceId %v, got %v", customResourceID, event.ResourceId)
	}

	if event.NftId.Cmp(customNftId) != 0 {
		t.Errorf("Expected nftId %v, got %v", customNftId, event.NftId)
	}

	if log.BlockNumber != customBlockNumber {
		t.Errorf("Expected blockNumber %d, got %d", customBlockNumber, log.BlockNumber)
	}
}
