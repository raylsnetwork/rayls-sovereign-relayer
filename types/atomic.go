// Decommissioning Teleport (vanilla, atomic): atomic types/consts below are deprecated; shared members (Transaction, DispatchedMessageToPrivateHub, …) are retained.

package types

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/EndpointV1"
)

type Transaction struct {
	TxHash              string
	TxHashDestination   common.Hash
	LogIndex            uint
	SharedID            string
	State               TransactionState
	Outcome             TransactionOutcome
	ProofInvalid        bool
	FromChainID         *big.Int
	ToChainID           *big.Int
	MsgID               [32]byte
	IsAtomic            bool
	UpdatedAt           time.Time
	CreatedAt           time.Time
	BatchID             string
	BatchPrivateHubHash common.Hash
	ResourceID          string
	FromContractAddress string
	FromUserAddress     string
	TransferID          string
	TransferAmount      string
	BlockNumber         uint64
	ParentHash          string
}

type DispatchedMessageToPrivateHub struct {
	MessageId             [32]byte                `json:"message_id"`
	From                  common.Address          `json:"from"`
	ToChainId             *big.Int                `json:"to_chain_id"`
	To                    common.Address          `json:"to"`
	Data                  EndpointV1.RaylsMessage `json:"data"`
	FromChainId           *big.Int                `json:"from_chain_id"`
	SharedId              string                  `json:"shared_id"`
	TxHashSource          common.Hash             `json:"tx_hash_source"`
	TxHashSourceTimestamp uint64                  `json:"tx_hash_source_timestamp"`
	TxHashSourceStatus    int8                    `json:"tx_hash_source_status"`
	// only used when vanilla; in case of Atomic - the ones in additionalData are used
	TxHashDestination          common.Hash `json:"tx_hash_destination"`
	TxHashDestinationTimestamp uint64      `json:"tx_hash_destination_timestamp"`
	TxHashDestinationStatus    int8        `json:"tx_hash_destination_status"`
	// proofs
	Proofs                     []byte                `json:"proofs"`
	TxTrieProof                common.Hash           `json:"tx_trie_proof"`
	BlockHash                  common.Hash           `json:"block_hash"`
	TxLocation                 uint                  `json:"tx_location"`
	TransactionType            BridgeTransactionType `json:"transaction_type"`
	LogIdx                     uint                  `json:"log_idx"`
	IsAtomic                   bool                  `json:"is_atomic"`
	TxSentToDestinationSuccess bool                  `json:"tx_sent_to_destination"`
	BatchId                    string                `json:"batch_id"`
	BatchPrivateHubHash        common.Hash           `json:"batch_private_hub_hash"`
	TokenAddress               common.Address        `json:"token_address"`
	ResourceId                 common.Hash           `json:"resource_id"`
	ParentHash                 string                `json:"parent_hash"`
	BlockNumber                uint64                `json:"block_number"`
}

func (d DispatchedMessageToPrivateHub) GetID() string {
	return common.Bytes2Hex(d.MessageId[:])
}

type BridgeTransactionType int

func (b BridgeTransactionType) String() string {
	switch b {
	case Transfer:
		return "Transfer"
	case Proof:
		return "Proof"
	default:
		return fmt.Sprintf("%d", int(b))
	}
}

const (
	Transfer BridgeTransactionType = 1
	Proof    BridgeTransactionType = 2
)

// TransactionState identifies the most recent on-chain action attempted on a
// row. Resolution of that action is captured in the orthogonal Outcome column.
//
// Valid transitions ((from, outcome) → next state):
//
//	Source side
//	───────────
//	(row created) → SourcePublish                              (initial)
//	SourcePublish        + success            → SourceFinalized            (via HandleSourceExecuted, when PNH confirms exec)
//	SourcePublish        + success            → SourceTimeoutRevert        (via atomic_expired, after expiry window)
//	SourcePublish        + reverted | failed  → EarlyRevertSigs            (via atomic_earlyrevert)
//	SourceTimeoutRevert  + success            → SourceRevertSigs           (via finalization, when PNH confirms revert)
//
//	Destination side
//	────────────────
//	(row created) → DestinationDispatch                        (initial; vanilla + atomic)
//	DestinationDispatch  + success (atomic)   → HubNotifiedExec            (via atomic_receipt; PNH told "executed")
//	DestinationDispatch  + reverted (atomic)  → HubNotifiedRevert          (via atomic_receipt; PNH told "reverted")
//	DestinationDispatch  + failed  (atomic)   → HubNotifiedRevert          (same)
//	DestinationDispatch  + *       (vanilla)  → (terminal)                 (no further transitions)
//	HubNotifiedExec      + success            → DestinationUnlockSigs      (via finalization + atomic_status)
//	HubNotifiedExec      + !success           → DestinationRevertSigs      (via finalization + atomic_status)
//  HubNotifiedRevert    + success            →                            (terminal on dest; source side handles escrow)

//	Signature flow (atomic; all broadcasts to dest endpoint)
//	────────────────────────────────────────────────────────
//	DestinationUnlockSigs  + success | reverted | failed → (terminal)
//	DestinationRevertSigs  + success | reverted | failed → (terminal)
//	SourceRevertSigs       + success | reverted | failed → (terminal)
//	EarlyRevertSigs        + success | reverted | failed → (terminal)
//
//	Race-condition states (TODO(race-handling) markers in code)
//	───────────────────────────────────────────────────────────
//	HubNotifiedExec      + reverted           → (orphan mint; needs dest-revert recovery)
//	HubNotifiedExec      + failed             → (PNH unreachable; needs retry)
//	HubNotifiedRevert    + failed             → (PNH unreachable; needs retry)
//	SourceTimeoutRevert  + reverted | failed  → (PNH unreachable; needs retry)
//
// SourceFinalized is purely a notification sink — set with outcome=success on entry
// to stop the timeout poller from re-querying the row. Rows that complete via
// SourceRevertSigs reach their own terminal naturally without needing this state.
type TransactionState int16

const (
	// Source side
	SourcePublish TransactionState = 1 // encrypted batch sent to PNH (vanilla + atomic)
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	SourceTimeoutRevert TransactionState = 2 // source notified PNH "revert this batch" after expiry (atomic)
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	EarlyRevertSigs TransactionState = 3 // early-revert sigs sent to dest endpoint (atomic, after publish failure)
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	SourceFinalized TransactionState = 4 // atomic only: PNH confirmed bridge fully executed; nothing left to do

	// Destination side
	DestinationDispatch TransactionState = 10 // dispatch tx broadcast on dest chain (vanilla mint OR atomic exec)
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	HubNotifiedExec TransactionState = 11 // atomic only: dispatch succeeded; dest told PNH "executed"
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	HubNotifiedRevert TransactionState = 12 // atomic only: dispatch failed; dest told PNH "reverted"

	// Atomic signatures
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	DestinationUnlockSigs TransactionState = 20
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	DestinationRevertSigs TransactionState = 21
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	SourceRevertSigs TransactionState = 22
)

type TransactionOutcome string

const (
	OutcomePending  TransactionOutcome = "pending"  // in-flight; phase not yet resolved
	OutcomeSuccess  TransactionOutcome = "success"  // mined OK
	OutcomeReverted TransactionOutcome = "reverted" // mined but reverted on-chain
	OutcomeFailed   TransactionOutcome = "failed"   // never mined / dead-letter
)

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type AtomicStatus uint8

const (
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	AtomicPendingStatus AtomicStatus = 0
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	AtomicExecutedStatus AtomicStatus = 1
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	AtomicRejectedStatus AtomicStatus = 2
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	AtomicRevertedStatus AtomicStatus = 3
)

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type AtomicStatusUpdateMessage struct {
	SharedID string
	Status   AtomicStatus
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type AtomicTeleportMessageStatus string

const (
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	AtomicPending AtomicTeleportMessageStatus = "Pending"
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	AtomicExecuted AtomicTeleportMessageStatus = "Executed"
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	AtomicRejected AtomicTeleportMessageStatus = "Rejected"
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	AtomicReverted AtomicTeleportMessageStatus = "Reverted"
)

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (ms AtomicTeleportMessageStatus) String() string {
	return string(ms)
}

type LastProcessedBlockDocument string

const (
	DocumentIdLastProcessedBlockPrivateHub  LastProcessedBlockDocument = "Private Hub"
	DocumentIdLastProcessedBlockPrivacyNode LastProcessedBlockDocument = "Privacy Node"
	LastProcessedBlockDocumentPublicChain   LastProcessedBlockDocument = "Public Chain"
	LastProcessedBlockDocumentPrivateChain  LastProcessedBlockDocument = "Private Chain"
)

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type AtomicTeleportAdditionalData struct {
	TxHashDestinationRevert          common.Hash `json:"tx_hash_destination_revert,omitempty"`
	TxHashDestinationRevertTimestamp uint64      `json:"tx_hash_destination_revert_timestamp,omitempty"`
	TxHashDestinationRevertStatus    int8        `json:"tx_hash_destination_revert_status,omitempty"`
	TxHashSourceRevert               common.Hash `json:"tx_hash_source_revert,omitempty"`
	TxHashSourceRevertTimestamp      uint64      `json:"tx_hash_source_revert_timestamp,omitempty"`
	TxHashSourceRevertStatus         int8        `json:"tx_hash_source_revert_status,omitempty"`
	TxHashDestination                common.Hash `json:"tx_hash_destination,omitempty"`
	TxHashDestinationTimestamp       uint64      `json:"tx_hash_destination_timestamp,omitempty"`
	TxHashDestinationStatus          int8        `json:"tx_hash_destination_status,omitempty"`
	RevertReason                     string      `json:"revert_reason,omitempty"`
	SharedId                         string      `json:"shared_id,omitempty"`
	BatchPrivateHubHash              common.Hash `json:"batch_private_hub_hash,omitempty"`
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type CalldataSignature struct {
	SharedId      string
	Signature     []byte
	ResourceId    [32]byte
	SignatureType CallDataSignatureType
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (c CalldataSignature) GetID() string {
	return c.SharedId
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type CallDataSignatureType uint8

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (b CallDataSignatureType) String() string {
	switch b {
	case RevertOnSenderSide:
		return "RevertOnSenderSide"
	case RevertOnDestinationSide:
		return "RevertOnDestinationSide"
	case UnlockOnDestinationSide:
		return "UnlockOnDestinationSide"
	default:
		return fmt.Sprintf("%d", b)
	}
}

const (
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	RevertOnSenderSide CallDataSignatureType = 1
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	RevertOnDestinationSide CallDataSignatureType = 2
	// Deprecated: Decommissioning Teleport (vanilla, atomic).
	UnlockOnDestinationSide CallDataSignatureType = 3
)
