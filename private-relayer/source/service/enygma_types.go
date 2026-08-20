package service

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

// EnygmaCreation represents an Enygma token creation event
type EnygmaCreation struct {
	InitialSupply *big.Int
}

// EnygmaTransferTx represents a single transfer transaction split from EnygmaSendTransferCC.
// Each destination in EnygmaSendTransferCC becomes one EnygmaTransferTx
type EnygmaTransferTx struct {
	DestIdx     int
	MessageId   string
	ReferenceId [32]byte
	FromAddress common.Address
	ToChainId   *big.Int
	ToAmount    *big.Int
	ToAddress   common.Address
	// ProgramData is this recipient's array of programmability steps, carried verbatim from the
	// source event's parallel ProgramData[][] field. Converted to the executor's identically-shaped
	// struct type at the executeProgramData call site (abigen emits a distinct type per binding pkg).
	ProgramData []types.EnygmaProgramData
}

type EnygmaSupplyUpdate struct {
	TxHash common.Hash
	To     common.Address
	Amount *big.Int
}

// EnygmaMint represents an Enygma mint event
type EnygmaMint struct {
	To     common.Address
	Amount *big.Int
	TxHash common.Hash
}

// EnygmaBurn represents an Enygma burn event
type EnygmaBurn struct {
	From   common.Address
	Amount *big.Int
	TxHash common.Hash
}

// EnygmaDepositToDvp represents an Enygma deposit to Dvp event
type EnygmaDepositToDvp struct {
	Amount        *big.Int
	From          common.Address
	ReferenceId   [32]byte
	TxHash        common.Hash
	TxBlockNumber *big.Int
}

// EnygmaWithdrawFromDvp represents an Enygma withdrawal from Dvp event
type EnygmaWithdrawFromDvp struct {
	Amount        *big.Int
	To            common.Address
	ReferenceId   [32]byte
	TxHash        common.Hash
	TxBlockNumber *big.Int
}

// pendingFinalization represents a finalization that failed and needs to be retried.
type pendingFinalization struct {
	ID          string
	BlockNumber uint64
	ResourceID  string
	RetryCount  int
}

// EnygmaSwapWithDvpForERC721 represents an Enygma swap for ERC721 event
type EnygmaSwapWithDvpForERC721 struct {
	SharedId      string
	DestChainId   *big.Int
	From          common.Address
	EnygmaAmount  *big.Int
	NftResourceId string
	NftId         string
}

// EnygmaSwapWithDvpForERC1155 represents an Enygma swap for ERC1155 event
type EnygmaSwapWithDvpForERC1155 struct {
	SharedId       string
	DestChainId    *big.Int
	From           common.Address
	EnygmaAmount   *big.Int
	NftResourceId  string
	NftId          string
	NftAmountOrOne *big.Int
}
