package testdata

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/EnygmaPNEvents"
)

func TestEnygmaCreationLog(t *testing.T) {
	log := GetEnygmaCreationLog()

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackEnygmaCreationEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse EnygmaCreation log: %v", err)
	}

	if event.ResourceId != [32]byte{0x01, 0x02, 0x03} {
		t.Errorf("Expected resourceId [1 2 3], got %v", event.ResourceId)
	}

	if event.InitialSupply.Int64() != 1000000 {
		t.Errorf("Expected initialSupply 1000000, got %d", event.InitialSupply.Int64())
	}
}

func TestEnygmaMintLog(t *testing.T) {
	log := GetEnygmaMintLog()

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackEnygmaMintEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse EnygmaMint log: %v", err)
	}

	if event.ResourceId != [32]byte{0x01, 0x02, 0x03} {
		t.Errorf("Expected resourceId [1 2 3], got %v", event.ResourceId)
	}

	if event.Amount.Int64() != 500 {
		t.Errorf("Expected amount 500, got %d", event.Amount.Int64())
	}
}

func TestEnygmaBurnLog(t *testing.T) {
	log := GetEnygmaBurnLog()

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackEnygmaBurnEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse EnygmaBurn log: %v", err)
	}

	if event.ResourceId != [32]byte{0x01, 0x02, 0x03} {
		t.Errorf("Expected resourceId [1 2 3], got %v", event.ResourceId)
	}

	if event.Amount.Int64() != 200 {
		t.Errorf("Expected amount 200, got %d", event.Amount.Int64())
	}
}

func TestEnygmaSendTransferPNHLog(t *testing.T) {
	log := GetEnygmaSendTransferPNHLog()

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackEnygmaSendTransferPNHEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse EnygmaSendTransferPNH log: %v", err)
	}

	if event.ResourceId != [32]byte{0x01, 0x02, 0x03} {
		t.Errorf("Expected resourceId [1 2 3], got %v", event.ResourceId)
	}
}

func TestEnygmaSwapWithDvpForERC721Log(t *testing.T) {
	log := GetEnygmaSwapWithDvpForERC721Log()

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackEnygmaSwapWithDvpForERC721Event(&log)
	if err != nil {
		t.Fatalf("Failed to parse EnygmaSwapWithDvpForERC721 log: %v", err)
	}

	if event.NftId.Int64() != 123 {
		t.Errorf("Expected nftId 123, got %d", event.NftId.Int64())
	}
}

func TestEnygmaSwapWithDvpForERC1155Log(t *testing.T) {
	log := GetEnygmaSwapWithDvpForERC1155Log()

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackEnygmaSwapWithDvpForERC1155Event(&log)
	if err != nil {
		t.Fatalf("Failed to parse EnygmaSwapWithDvpForERC1155 log: %v", err)
	}

	if event.NftId.Int64() != 456 {
		t.Errorf("Expected nftId 456, got %d", event.NftId.Int64())
	}
}

func TestDvp721DepositIntoDvpLog(t *testing.T) {
	log := GetDvp721DepositIntoDvpLog()

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackDvp721DepositIntoDvpEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse Dvp721DepositIntoDvp log: %v", err)
	}

	if event.NftId.Int64() != 789 {
		t.Errorf("Expected nftId 789, got %d", event.NftId.Int64())
	}
}

func TestDvp721WithdrawFromDvpLog(t *testing.T) {
	log := GetDvp721WithdrawFromDvpLog()

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackDvp721WithdrawFromDvpEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse Dvp721WithdrawFromDvp log: %v", err)
	}

	if event.NftId.Int64() != 789 {
		t.Errorf("Expected nftId 789, got %d", event.NftId.Int64())
	}
}

func TestDvp721SwapForEnygmaLog(t *testing.T) {
	log := GetDvp721SwapForEnygmaLog()

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackDvp721SwapForEnygmaEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse Dvp721SwapForEnygma log: %v", err)
	}

	if event.NftId.Int64() != 999 {
		t.Errorf("Expected nftId 999, got %d", event.NftId.Int64())
	}
}

func TestDvp1155CreationLog(t *testing.T) {
	log := GetDvp1155CreationLog()

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackDvp1155CreationEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse Dvp1155Creation log: %v", err)
	}

	if event.ResourceId != [32]byte{0x0d, 0x0e, 0x0f} {
		t.Errorf("Expected resourceId [13 14 15], got %v", event.ResourceId)
	}
}

func TestDvp1155MintLog(t *testing.T) {
	log := GetDvp1155MintLog()

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackDvp1155MintEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse Dvp1155Mint log: %v", err)
	}

	if event.TokenId.Int64() != 111 {
		t.Errorf("Expected tokenId 111, got %d", event.TokenId.Int64())
	}
}

func TestDvp1155BurnLog(t *testing.T) {
	log := GetDvp1155BurnLog()

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackDvp1155BurnEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse Dvp1155Burn log: %v", err)
	}

	if event.TokenId.Int64() != 111 {
		t.Errorf("Expected tokenId 111, got %d", event.TokenId.Int64())
	}
}

func TestDvp1155DepositIntoDvpLog(t *testing.T) {
	log := GetDvp1155DepositIntoDvpLog()

	filterer := EnygmaPNEvents.NewEnygmaPNEvents()

	event, err := filterer.UnpackDvp1155DepositIntoDvpEvent(&log)
	if err != nil {
		t.Fatalf("Failed to parse Dvp1155DepositIntoDvp log: %v", err)
	}

	if event.TokenId.Int64() != 222 {
		t.Errorf("Expected tokenId 222, got %d", event.TokenId.Int64())
	}
}
