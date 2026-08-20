package types

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Point struct {
	X *big.Int
	Y *big.Int
}
type EnygmaPublicValues struct {
	Commitments map[string]*Point
	PublicKeys  map[string]*big.Int // Changed from Point to big.Int - now computed as Poseidon(sk, sk)
}

type Enygma struct {
	ResourceId           string
	FinalizedR           *big.Int
	FinalizedBalance     *big.Int
	FinalizedBlockNumber *big.Int
	PendingBlockNumber   *big.Int
}

// HistoryStatus represents the lifecycle state of an EnygmaHistory record.
type HistoryStatus int

const (
	HistoryStatusPending   HistoryStatus = 0
	HistoryStatusConfirmed HistoryStatus = 1
)

type EnygmaHistory struct {
	ResourceId            string
	FromChainId           *big.Int
	BlockNumberPrivateHub *big.Int
	RFactor               *big.Int
	BalanceChange         *big.Int // Positive for incoming, negative for outgoing
	EventType             EnygmaEventType
	PrivateHubTxHash      string // Private Hub tx hash of the transaction
	SignedTxBytes         []byte // RLP-encoded signed tx for crash-recovery re-broadcast
}

type EnygmaEventType uint8

// Matching the enum IEnygmaV1.TxType (0-5).
// Values 6+ are relayer-internal and do not map to any contract enum.
const (
	EnygmaCreation EnygmaEventType = iota
	EnygmaMint
	EnygmaBurn
	EnygmaTransfer
	EnygmaDepositToDvp
	EnygmaWithdrawFromDvp

	// Relayer-internal DvP event types for crash-recovery tracking.
	DvpERC721Withdraw       EnygmaEventType = 6
	DvpERC1155Withdraw      EnygmaEventType = 7
	DvpSwapExecution        EnygmaEventType = 8
	EnygmaWithdrawFromDvpPL EnygmaEventType = 9
	DvpSwapCompletion       EnygmaEventType = 10
)

func (e EnygmaEventType) String() string {
	switch e {
	case EnygmaCreation:
		return "creation"
	case EnygmaMint:
		return "mint"
	case EnygmaBurn:
		return "burn"
	case EnygmaTransfer:
		return "transfer"
	case EnygmaDepositToDvp:
		return "deposit"
	case EnygmaWithdrawFromDvp:
		return "withdraw"
	case DvpERC721Withdraw:
		return "dvp_erc721_withdraw"
	case DvpERC1155Withdraw:
		return "dvp_erc1155_withdraw"
	case DvpSwapExecution:
		return "dvp_swap_execution"
	case EnygmaWithdrawFromDvpPL:
		return "withdraw_pl"
	case DvpSwapCompletion:
		return "dvp_swap_completion"
	}

	return "unknown"
}

type EnygmaCheckpoint struct {
	ID                      string
	ResourceId              string
	FinalizedR              *big.Int
	FinalizedPublicBalanceX *big.Int
	FinalizedPublicBalanceY *big.Int
	FinalizedBlockNumber    *big.Int
	PendingBlockNumber      *big.Int
	Status                  EnygmaCheckpointStatus
}

type EnygmaCheckpointStatus uint8

const (
	EnygmaCheckpointStatusTentative EnygmaCheckpointStatus = iota
	EnygmaCheckpointStatusFinal
)

// Common parameters shared by all proof requests
type CommonProofRequest struct {
	ResourceId                 string
	SenderSecretKey            *big.Int
	SenderBalance              *big.Int
	SenderRandomFactor         *big.Int
	SenderChainId              *big.Int
	SenderAmount               *big.Int
	DestinationChainIDs        []*big.Int
	DestinationAmounts         []*big.Int
	DestinationPublicKeys      []*big.Int // Changed from []*Point to []*big.Int - now Poseidon(sk, sk)
	DestinationPreviousCommits []*Point
	DestinationNewCommits      []*Point
	DestinationRandomFactors   []*big.Int
	Nullifier                  *big.Int
	BlockNumber                *big.Int
	DestinationSharedSecrets   []*big.Int // k array of secrets: secrets[senderIdx] = Poseidon(previousR, sk), others from shared secrets
	ArrayHashSecrets           []*big.Int // k array of arrayHashSecret[i] = Poseidon(secrets[i], secrets[i])
	MessageTags                []*big.Int // k array of tag messages
}

type TransferProofRequest struct {
	*CommonProofRequest
}
type DepositProofRequest struct {
	*CommonProofRequest
	TokenAddress      common.Address
	DepositCommitment *big.Int
	DepositSalt       *big.Int
	DepositPublicKey  *big.Int
}

type WithdrawProofRequest struct {
	*CommonProofRequest
	TokenAddress       common.Address
	DepositCommitments []*big.Int
	DepositSecretKeys  []*big.Int
	DepositAmounts     []*big.Int
	DepositSalts       []*big.Int
}
type ResponseEnygmaProofAPI struct {
	Pi_A          []string   `json:"pi_a"`
	Pi_B          [][]string `json:"pi_b"`
	Pi_C          []string   `json:"pi_c"`
	Public_Signal []string   `json:"public_signal"`
}

type EnygmaProofResponse struct {
	PiA          [2]*big.Int    `json:"pi_a"`
	PiB          [2][2]*big.Int `json:"pi_b"`
	PiC          [2]*big.Int    `json:"pi_c"`
	PublicSignal []*big.Int     `json:"public_signal"`
}

type EnygmaMintCompleted struct {
	ResourceId         [32]byte
	Amount             *big.Int
	DestinationChainId *big.Int
}

type EnygmaPlTransferCompleted struct {
	ReferenceId     [32]byte
	TransactionHash [32]byte
	Timestamp       *big.Int
	ChainId         *big.Int
}

type EnygmaTransferCompleted struct {
	MessageId       string
	TransactionHash string
	ChainId         *big.Int
}

// Enygma Batch Transfer types
type EnygmaTransferBatch struct {
	ResourceId            string
	BlockNumberPrivateHub *big.Int
	FromChainID           *big.Int
	ToChainID             *big.Int
	ToRValueToAdd         *big.Int
	Transactions          []*EnygmaTransferBatchTx
	BatchId               string
	PrivateHubTxHash      string // Private Hub tx hash that emitted this batch
	// This context is needed in order to be able to link the traces
	// After we decrypt the batch we want to see the full trace chain
	Ctx context.Context
	// Unix timestamp for measuring soft finality
	// This is the time between the tx processing has stared until it is sent to the Private Network Hub
	SoftFinalityStartTimestamp int64
}

// EnygmaProgramData is the domain representation of a single programmability step that the
// destination ProgrammabilityExecutor resolves and dispatches. It mirrors the on-chain
// SharedObjects.EnygmaProgramData struct but keeps the domain layer free of contract bindings;
// converters translate to/from the abigen-generated per-package structs at the chain boundary.
type EnygmaProgramData struct {
	ResourceId      [32]byte
	ContractAddress common.Address
	Selector        [4]byte
	Args            []byte
}

type EnygmaTransferBatchTx struct {
	MessageId   string
	ReferenceId [32]byte
	FromAddress common.Address
	ToAmount    *big.Int
	ToAddress   common.Address
	// ProgramData is this recipient's ordered array of programmability steps. It becomes the
	// steps argument to one ProgrammabilityExecutor.executeProgramData(EnygmaProgramData[],
	// uint256,address) tx on the destination PN (with ToAmount as expectedMintTotal and
	// FromAddress as originSender). For a plain transfer the sender stamps a single-element
	// [mintStep]; composed transfers carry [mintStep, userStep...]. The split into one batch tx
	// per recipient happens upstream in convertEnygmaSendTransferCC, so this is a single
	// per-recipient array, not the event's per-recipient-of-arrays.
	ProgramData []EnygmaProgramData
	// Unix timestamp for measuring hard finality
	// This is the time between the batch has been created until the time it was mined by the destination node
	SendTimestamp int64
}

// Needed for the tx generator
func (e *EnygmaTransferBatchTx) GetID() string {
	return e.MessageId
}

type EnygmaCrossTransferData struct {
	*EnygmaTransferBatchTx
	EnygmaAddress common.Address
	FromChainID   *big.Int
}

// Needed for the tx generator
func (e *EnygmaCrossTransferData) GetID() string {
	return e.MessageId
}

type EnygmaFinalizedBalance struct {
	ResourceId           string
	FinalizedBlockNumber *big.Int
	PendingBlockNumber   *big.Int
	BalanceX             *big.Int
	BalanceY             *big.Int
}

type EnygmaSupplyUpdate struct {
	Amount *big.Int
	Type   EnygmaEventType
}

type EnygmaTransferFailed struct {
	ReferenceID   [32]byte
	EnygmaAddress common.Address
	Sender        common.Address
	Amount        *big.Int
	Reason        string
}

// required by the enygma revert tx generator
func (e *EnygmaTransferFailed) GetID() string {
	return hex.EncodeToString(e.ReferenceID[:])
}

type EnygmaSupplyUpdateFailed struct {
	TxHash        string
	EnygmaAddress common.Address
	Amount        *big.Int
	To            common.Address
	Type          EnygmaEventType
}

func (e *EnygmaSupplyUpdateFailed) GetID() string {
	return e.TxHash
}

type PaymentSpendKey struct {
	SecretKey *big.Int
	PublicKey *big.Int
}
