package testdata

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EnygmaPNEvents"
)

func TestDvp721CreationLog(t *testing.T) {
	log := NewDvp721CreationLogWith()

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackDvp721CreationEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse Dvp721Creation log: %v", err)
	}

	if event.ResourceId != [32]byte{0x0a, 0x0b, 0x0c} {
		t.Errorf("Expected resourceId [10 11 12], got %v", event.ResourceId)
	}

	if log.BlockNumber != 105 {
		t.Errorf("Expected block number 105, got %d", log.BlockNumber)
	}
}

func TestDvp721MintLog(t *testing.T) {
	log := NewDvp721MintLogWith()

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackDvp721MintEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse Dvp721Mint log: %v", err)
	}

	if event.ResourceId != [32]byte{0x0a, 0x0b, 0x0c} {
		t.Errorf("Expected resourceId [10 11 12], got %v", event.ResourceId)
	}

	if event.NftId.Int64() != 42 {
		t.Errorf("Expected nftId 42, got %d", event.NftId.Int64())
	}
}

func TestDvp721BurnLog(t *testing.T) {
	log := NewDvp721BurnLogWith()

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackDvp721BurnEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse Dvp721Burn log: %v", err)
	}

	if event.ResourceId != [32]byte{0x0a, 0x0b, 0x0c} {
		t.Errorf("Expected resourceId [10 11 12], got %v", event.ResourceId)
	}

	if event.NftId.Int64() != 42 {
		t.Errorf("Expected nftId 42, got %d", event.NftId.Int64())
	}
}
