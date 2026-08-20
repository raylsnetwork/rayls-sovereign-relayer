package types

import (
	"fmt"
	"math/big"
	"time"
)

type MerkleProof struct {
	Indices  *big.Int   `json:"indices"`
	Elements []*big.Int `json:"elements"`
}

type DvpSwapMessage struct {
	SharedId      string
	To            string
	ChainId       *big.Int
	PNTxHash      string
	PNTxTimestamp time.Time

	TokenInAmount      *big.Int
	TokenInAddress     string
	TokenInResourceID  string
	TokenInType        DvpTokenType
	TokenInID          string
	TokenOutAmount     *big.Int
	TokenOutAddress    string
	TokenOutResourceID string
	TokenOutType       DvpTokenType
	TokenOutID         string

	// Salt the initiator used to create his self-destination commitment(the token he'll receive)
	InitiatorSelfSalt *big.Int
}

type DvpSwapStatus uint8

const (
	DvpSwapCreated DvpSwapStatus = iota
	DvpSwapInitiated
	DvpSwapInitiationFailed
	DvpSwapWaitingConfirmation
	DvpSwapCompleted
	DvpSwapFailed
	DvpSwapCancelled
	DvpSwapTimedOut
)

type DvpSwap struct {
	SharedID           string
	From               string
	To                 string
	SourceChainID      *big.Int
	SourceTxHash       string
	SourceBlockNumber  *big.Int
	SourceTxTimestamp  time.Time
	DestChainID        *big.Int
	TokenInAmount      *big.Int
	TokenInAddress     string // TODO delete this
	TokenInResourceID  string
	TokenInType        DvpTokenType
	TokenInID          string
	TokenOutAmount     *big.Int
	TokenOutAddress    string // TODO delete this
	TokenOutResourceID string
	TokenOutType       DvpTokenType
	TokenOutID         string
	Status             DvpSwapStatus
	CreatedAt          time.Time
	CancelledAt        *time.Time
	DestSalt           *big.Int
	SelfSalt           *big.Int
	CancelPreimage     *big.Int
}

type DvpSwapType int

const (
	DvpSwapTypeEnygmaERC721 DvpSwapType = iota
	DvpSwapTypeERC721Enygma
	DvpSwapTypeEnygmaERC1155
	DvpSwapTypeERC1155Enygma
)

func (swap *DvpSwap) Type() (DvpSwapType, error) {
	// Enygma -> ERC721
	if swap.TokenInType == DvpEnygma && swap.TokenOutType == DvpERC721 {
		return DvpSwapTypeEnygmaERC721, nil
	}
	if swap.TokenInType == DvpERC721 && swap.TokenOutType == DvpEnygma {
		return DvpSwapTypeERC721Enygma, nil
	}
	// Enygma -> ERC1155_NON_FUNGIBLE
	if swap.TokenInType == DvpEnygma && swap.TokenOutType == DvpERC1155 {
		return DvpSwapTypeEnygmaERC1155, nil
	}
	if swap.TokenInType == DvpERC1155 && swap.TokenOutType == DvpEnygma {
		return DvpSwapTypeERC1155Enygma, nil
	}

	return 0, fmt.Errorf("unsupported swap type")
}

type DvpTokenType int

const (
	DvpCustom     DvpTokenType = iota // 0
	DvpERC20                          // 1
	DvpERC721                         // 2
	DvpERC1155                        // 3
	DvpEnygma                         // 4
	DvpDVPERC721                      // 5 - DVP-enabled ERC721
	DvpDVPERC1155                     // 6 - DVP-enabled ERC1155
)

type DvpDepositStatus int

const (
	// Created, waiting to be confirmed
	DvpDepositPending DvpDepositStatus = iota
	// Confirmed, ready to be spent
	DvpDepositUnspent
	// In use, being processed by some operation
	DvpDepositLocked
	// Spent, cannot be used anymore
	DvpDepositSpent
)

type DvpDeposit struct {
	UserAddress  string
	Salt         *big.Int
	TokenAmount  *big.Int
	TokenAddress string
	TokenType    DvpTokenType
	TokenID      string
	TreeNumber   int
	Commitment   *big.Int
	Nullifier    *big.Int
	Status       DvpDepositStatus
	CreatedAt    time.Time
}

type MerkleTree struct {
	Type         DvpTokenType
	TokenAddress string
	Number       int
	Depth        int
	Leaves       []string
	CreatedAt    time.Time
}

// DvpBalanceUpdateType defines whether tokens are being burned or minted in a operation.
type DvpBalanceUpdateType uint8

const (
	Burn DvpBalanceUpdateType = 0
	Mint DvpBalanceUpdateType = 1
)

// DvpBalanceUpdated represents an encrypted balance update message sent to the Private Network Hub.
type DvpBalanceUpdated struct {
	ErcId      *big.Int
	TokenType  uint8
	ResourceId string
	SharedId   string

	From string
	To   string

	SourceChainId     *big.Int
	SourceTxHash      string
	SourceTxTimestamp time.Time

	DestinationTxHash      string
	DestinationTxTimestamp time.Time
	DestinationChainId     *big.Int

	Amount     *big.Int
	UpdateType DvpBalanceUpdateType
}

// DvpCommitmentsData carries data from a Commitments event for processing by the receiver.
type DvpCommitmentsData struct {
	TokenAddress string
	TokenType    *big.Int
	TreeNumber   *big.Int
	Commitments  []*big.Int
}

// DvpNullifierData carries data from a Nullifiers event for processing by the receiver.
type DvpNullifierData struct {
	TokenAddress string
	Nullifiers   []*big.Int
}

// DvpSwapInitiatedData carries data from a SwapInitiated event whose
// payload was successfully decrypted by CTS — i.e. this chain is the responder.
type DvpSwapInitiatedData struct {
	Message           *DvpSwapMessage
	InitiatorDestSalt *big.Int
	ExpiresAt         *big.Int
}

