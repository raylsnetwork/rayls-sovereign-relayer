// Decommissioning Teleport (vanilla, atomic).

// Package testdata implements the legacy public-chain (RN) Teleport bridge relayer.
//
// Deprecated: Decommissioning Teleport (vanilla, atomic).
package testdata

import "github.com/raylsnetwork/rayls-sovereign-relayer/public-relayer/service"

var (
	RevertSig1 = service.RevertSignature{ID: "sig-1", Data: []byte{0x01, 0x02, 0x03}}
	RevertSig2 = service.RevertSignature{ID: "sig-2", Data: []byte{0x04, 0x05, 0x06}}
	RevertSig3 = service.RevertSignature{ID: "sig-3", Data: []byte{0x07, 0x08, 0x09}}
)
