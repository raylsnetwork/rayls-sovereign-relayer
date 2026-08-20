package types

import (
	"fmt"
	"time"
)

type ErcStandard uint8

const (
	CUSTOM ErcStandard = iota
	ERC20
	ERC721
	ERC1155
	ENYGMA
	DVPERC721
	DVPERC1155
	// Test-only variants. Ordinals MUST stay aligned with the on-chain SharedObjects.ErcStandard
	// and RaylsBridgeableERC enums — append only, never reorder. The relayer passes the standard
	// byte straight through to meta.FactoryTemplate, so the receiver factory resolves the *_TEST_KEY
	// seeded bytecode (see ResourceManager._keyForTemplate). No relayer logic branches on these.
	ERC20TEST
	ERC721TEST
	ERC1155TEST
	ENYGMATEST
	DVPERC721TEST
	DVPERC1155TEST

	// ercStandardEnd is a sentinel that MUST stay the last value: it marks the exclusive upper
	// bound of the valid enum range so Uint8ToErcStandard stays correct when new variants are
	// appended above. Append real variants before it, never after.
	ercStandardEnd
)

// ResourceDeployType mirrors the on-chain enum ResourceDeployType (RaylsMessage.sol):
// BYTECODE deploys issuer-supplied runtime bytecode; FACTORY deploys from the seeded
// template registry keyed by the standard's factoryTemplate.
type ResourceDeployType uint8

const (
	ResourceDeployTypeBytecode ResourceDeployType = iota
	ResourceDeployTypeFactory
)

// Convert uint8 to ErcStandard
func Uint8ToErcStandard(value uint8) (ErcStandard, error) {
	ercStandard := ErcStandard(value)
	// CUSTOM is 0, so the lower bound holds for any uint8; ercStandardEnd is the exclusive upper
	// bound, keeping this self-maintaining as variants are appended.
	if ercStandard < ercStandardEnd {
		return ercStandard, nil
	}
	return 0, fmt.Errorf("invalid ErcStandard value: %d", value)
}

const (
	// proof api service routes
	PROOF_API_ENYGMA_2          = "/generateProofTransfer-2"
	PROOF_API_ENYGMA_3          = "/generateProofTransfer-3"
	PROOF_API_ENYGMA_4          = "/generateProofTransfer-4"
	PROOF_API_ENYGMA_5          = "/generateProofTransfer-5"
	PROOF_API_ENYGMA_6          = "/generateProofTransfer-6"
	PROOF_API_ENYGMA_WITHDRAW_2 = "/generateProofWithdraw-2"
	PROOF_API_ENYGMA_WITHDRAW_3 = "/generateProofWithdraw-3"
	PROOF_API_ENYGMA_WITHDRAW_4 = "/generateProofWithdraw-4"
	PROOF_API_ENYGMA_WITHDRAW_5 = "/generateProofWithdraw-5"
	PROOF_API_ENYGMA_WITHDRAW_6 = "/generateProofWithdraw-6"
	PROOF_API_ENYGMA_DEPOSIT_2  = "/generateProofDeposit-2"
	PROOF_API_ENYGMA_DEPOSIT_3  = "/generateProofDeposit-3"
	PROOF_API_ENYGMA_DEPOSIT_4  = "/generateProofDeposit-4"
	PROOF_API_ENYGMA_DEPOSIT_5  = "/generateProofDeposit-5"
	PROOF_API_ENYGMA_DEPOSIT_6  = "/generateProofDeposit-6"
	PROOF_API_JOIN_SPLIT_ENYGMA = "/join-split-enygma"
	PROOF_API_JOIN_SPLIT_1155   = "/join-split-1155"
	PROOF_API_OWNERSHIP_721     = "/ownership-721"
	PROOF_API_OWNERSHIP_1155    = "/ownership-1155"
)

type DvpCommunicatorStatus uint8

const (
	NOSTATUS DvpCommunicatorStatus = iota
	Swap721ForEnygmaSent
	Swap721ForEnygmaReceived
	SwapEnygmaFor721Sent
	SwapEnygmaFor721Received
	Swap721ForEnygmaProcessing
	SwapDoneReadyForWithdraw
	Swap1155WaitingEnygmaSwapSent
	Swap1155WaitingEnygmaSwapReceived
	SwapEnygmaWaiting1155SwapSent
	SwapEnygmaWaiting1155SwapReceived
	Swap1155ForEnygmaProcessing
	SwapError
	SwapCancellationInitiated
	SwapCancelled
	SwapTimedOut
)

type CommunicatorContexts int

const (
	Dvp CommunicatorContexts = iota
)

type ResourceLock struct {
	ResourceId string    `bson:"_id"`
	ExpiresAt  time.Time `bson:"expires_at"`
}
