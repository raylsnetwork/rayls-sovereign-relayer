// Decommissioning Teleport (vanilla, atomic).
//
// FIXME(integration): does not compile under `-tags integration` — assigns ResourceId [32]byte into
// the repository model's []byte. Pre-existing drift; left unfixed by request.

package testdata

import (
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

var (
	CalldataSignature1 = types.CalldataSignature{
		SharedId:  "shared-id-1",
		Signature: []byte{0xDE, 0xAD, 0xBE, 0xEF},
		ResourceId: [32]byte{
			0x01, 0x02, 0x03, 0x04,
			0x05, 0x06, 0x07, 0x08,
			0x09, 0x0A, 0x0B, 0x0C,
			0x0D, 0x0E, 0x0F, 0x10,
			0x11, 0x12, 0x13, 0x14,
			0x15, 0x16, 0x17, 0x18,
			0x19, 0x1A, 0x1B, 0x1C,
			0x1D, 0x1E, 0x1F, 0x20,
		},
		SignatureType: types.RevertOnSenderSide,
	}
	CalldataSignature2 = types.CalldataSignature{
		SharedId:  "shared-id-2",
		Signature: []byte{0xAA, 0xBB, 0xCC, 0xDD},
		ResourceId: [32]byte{
			0x20, 0x1F, 0x1E, 0x1D,
			0x1C, 0x1B, 0x1A, 0x19,
			0x18, 0x17, 0x16, 0x15,
			0x14, 0x13, 0x12, 0x11,
			0x10, 0x0F, 0x0E, 0x0D,
			0x0C, 0x0B, 0x0A, 0x09,
			0x08, 0x07, 0x06, 0x05,
			0x04, 0x03, 0x02, 0x01,
		},
		SignatureType: types.RevertOnDestinationSide,
	}
)

var (
	CalldataSignatureModel1 = repository.CalldataSignature{
		SharedId:                CalldataSignature1.SharedId,
		Status:                  0,
		Signature:               CalldataSignature1.Signature,
		ResourceId:              CalldataSignature1.ResourceId,
		SignatureExecuteChainId: "0",
		DestinationChainId:      "0",
		SignatureType:           uint8(CalldataSignature1.SignatureType),
	}

	CalldataSignatureModel2 = repository.CalldataSignature{
		SharedId:                CalldataSignature2.SharedId,
		Status:                  0,
		Signature:               CalldataSignature2.Signature,
		ResourceId:              CalldataSignature2.ResourceId,
		SignatureExecuteChainId: "0",
		DestinationChainId:      "0",
		SignatureType:           uint8(CalldataSignature2.SignatureType),
	}
)
