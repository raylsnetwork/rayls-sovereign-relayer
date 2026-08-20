package testdata

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EnygmaPNEvents"
)

func TestEnygmaCreationLogWithOptions(t *testing.T) {
	// Test with custom resource ID
	customResourceID := [32]byte{0xaa, 0xbb, 0xcc}
	log := NewEnygmaCreationLogWith(
		WithCreateResourceID(customResourceID),
	)

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackEnygmaCreationEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse EnygmaCreation log: %v", err)
	}

	if event.ResourceId != customResourceID {
		t.Errorf("Expected resourceId %v, got %v", customResourceID, event.ResourceId)
	}

	// Should still have default initial supply
	if event.InitialSupply.Int64() != 1000000 {
		t.Errorf("Expected default initialSupply 1000000, got %d", event.InitialSupply.Int64())
	}
}

func TestEnygmaCreationLogWithMultipleOptions(t *testing.T) {
	// Test with multiple custom options
	customResourceID := [32]byte{0x11, 0x22, 0x33}
	customSupply := big.NewInt(5000000)
	customBlockNumber := uint64(999)

	log := NewEnygmaCreationLogWith(
		WithCreateResourceID(customResourceID),
		WithCreateInitialSupply(customSupply),
		WithCreateBlockNumber(customBlockNumber),
	)

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackEnygmaCreationEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse EnygmaCreation log: %v", err)
	}

	if event.ResourceId != customResourceID {
		t.Errorf("Expected resourceId %v, got %v", customResourceID, event.ResourceId)
	}

	if event.InitialSupply.Cmp(customSupply) != 0 {
		t.Errorf("Expected initialSupply %v, got %v", customSupply, event.InitialSupply)
	}

	if log.BlockNumber != customBlockNumber {
		t.Errorf("Expected blockNumber %d, got %d", customBlockNumber, log.BlockNumber)
	}
}

func TestEnygmaDepositToDvpLogWithOptions(t *testing.T) {
	// Test with custom values
	customResourceID := [32]byte{0xdd, 0xee, 0xff}
	customAmount := big.NewInt(999)
	customFrom := common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	customReferenceID := [32]byte{0x11, 0x22, 0x33}
	customBlockNumber := uint64(555)

	log := NewEnygmaDepositToDvpLogWith(
		WithDepositResourceID(customResourceID),
		WithDepositAmount(customAmount),
		WithDepositFrom(customFrom),
		WithDepositReferenceID(customReferenceID),
		WithDepositBlockNumber(customBlockNumber),
	)

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackEnygmaDepositToDvpEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse EnygmaDepositToDvp log: %v", err)
	}

	if event.ResourceId != customResourceID {
		t.Errorf("Expected resourceId %v, got %v", customResourceID, event.ResourceId)
	}

	if event.Amount.Cmp(customAmount) != 0 {
		t.Errorf("Expected amount %v, got %v", customAmount, event.Amount)
	}

	if event.From != customFrom {
		t.Errorf("Expected from %v, got %v", customFrom, event.From)
	}

	if event.ReferenceId != customReferenceID {
		t.Errorf("Expected referenceId %v, got %v", customReferenceID, event.ReferenceId)
	}

	if log.BlockNumber != customBlockNumber {
		t.Errorf("Expected blockNumber %d, got %d", customBlockNumber, log.BlockNumber)
	}
}
