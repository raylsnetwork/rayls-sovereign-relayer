package service

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Dvp ERC721 Event Types

// Dvp721Creation represents a Dvp ERC721 token creation event
type Dvp721Creation struct {
	ChainEventID string
	ResourceId   string
}

// Dvp721Mint represents a Dvp ERC721 mint event
type Dvp721Mint struct {
	ChainEventID string
	ResourceId   string
	NftId        *big.Int
}

// Dvp721Burn represents a Dvp ERC721 burn event
type Dvp721Burn struct {
	ChainEventID string
	ResourceId   string
	NftId        *big.Int
}

// Dvp721DepositIntoDvp represents a Dvp ERC721 deposit event
type Dvp721DepositIntoDvp struct {
	ChainEventID  string
	ResourceId    string
	NftId         *big.Int
	From          common.Address
	TxHash        string
	TxBlockNumber *big.Int
}

// Dvp721WithdrawFromDvp represents a Dvp ERC721 withdrawal event
type Dvp721WithdrawFromDvp struct {
	ChainEventID  string
	ResourceId    string
	NftId         *big.Int
	Owner         common.Address
	TxHash        string
	TxBlockNumber *big.Int
}

// Dvp721SwapForEnygma represents a Dvp ERC721 swap for Enygma event
type Dvp721SwapForEnygma struct {
	SharedId         string
	DestChainId      *big.Int
	From             common.Address
	NftResourceId    string
	NftId            string
	EnygmaResourceId string
	EnygmaAmount     *big.Int
	TxHash           string
	TxBlockNumber    *big.Int
	ValidityTime     uint64
}

// Dvp ERC1155 Event Types

// Dvp1155Creation represents a Dvp ERC1155 token creation event
type Dvp1155Creation struct {
	ChainEventID string
	ResourceId   string
}

// Dvp1155Mint represents a Dvp ERC1155 mint event
type Dvp1155Mint struct {
	ChainEventID string
	ResourceId   string
	TokenId      *big.Int
	Value        *big.Int
	Data         []byte
}

// Dvp1155Burn represents a Dvp ERC1155 burn event
// Note: The 'To' field from the contract event is not used in business logic
type Dvp1155Burn struct {
	ChainEventID string
	ResourceId   string
	TokenId      *big.Int
	Value        *big.Int
}

// Dvp1155DepositIntoDvp represents a Dvp ERC1155 deposit event
type Dvp1155DepositIntoDvp struct {
	ChainEventID  string
	ResourceId    string
	TokenId       *big.Int
	Value         *big.Int
	Data          []byte
	From          common.Address
	TxHash        string
	TxBlockNumber *big.Int
}

// Dvp1155WithdrawFromDvp represents a Dvp ERC1155 withdrawal event
// Note: The 'Data' field from the contract event is not used in business logic
type Dvp1155WithdrawFromDvp struct {
	ChainEventID  string
	ResourceId    string
	TokenId       *big.Int
	Value         *big.Int
	Owner         common.Address
	TxHash        string
	TxBlockNumber *big.Int
}

// Dvp1155SwapForEnygma represents a Dvp ERC1155 swap for Enygma event
// Note: The 'TokenData' field from the contract event is not used in business logic
type Dvp1155SwapForEnygma struct {
	SharedId         string
	DestChainId      *big.Int
	From             common.Address
	TokenResourceId  string
	TokenValue       *big.Int
	TokenId          string
	EnygmaResourceId string
	EnygmaAmount     *big.Int
	TxHash           string
	TxBlockNumber    *big.Int
	ValidityTime     uint64
}

// ZkDvP Enygma Swap Event Types

// DvpEnygmaSwapERC721 represents an Enygma to ERC721 swap event
type DvpEnygmaSwapERC721 struct {
	SharedId      string
	DestChainId   *big.Int
	From          common.Address
	ResourceId    string
	EnygmaAmount  *big.Int
	NftResourceId string
	NftId         string
	TxHash        string
	TxBlockNumber *big.Int
	ValidityTime  uint64
}

// DvpEnygmaSwapERC1155 represents an Enygma to ERC1155 swap event
type DvpEnygmaSwapERC1155 struct {
	SharedId       string
	DestChainId    *big.Int
	From           common.Address
	ResourceId     string
	EnygmaAmount   *big.Int
	NftResourceId  string
	NftId          string
	NftAmountOrOne *big.Int
	TxHash         string
	TxBlockNumber  *big.Int
	ValidityTime   uint64
}

// DvpSwapCancelled represents a swap cancelled event
type DvpSwapCancelled struct {
	SharedId           string
	DestChainId        *big.Int
	TokenInResourceId  string
	TokenInAmount      *big.Int
	TokenInId          *big.Int
	TokenInStandard    uint8
	TokenOutResourceId string
	TokenOutAmount     *big.Int
	TokenOutId         *big.Int
	TokenOutStandard   uint8
}
