// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package RaylsErc1155DvpHandler

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

// RaylsTrustedInit is an auto generated low-level Go binding around an user-defined struct.
type RaylsTrustedInit struct {
	Endpoint          common.Address
	RaylsNodeEndpoint common.Address
	UserGovernance    common.Address
	Owner             common.Address
	ResourceId        [32]byte
	Caller            common.Address
}

// SharedObjectsDvp1155ExtraData is an auto generated low-level Go binding around an user-defined struct.
type SharedObjectsDvp1155ExtraData struct {
	Key      string
	Value    string
	IsPublic bool
}

// SharedObjectsDvpSwapCompletedParams is an auto generated low-level Go binding around an user-defined struct.
type SharedObjectsDvpSwapCompletedParams struct {
	TokenId            *big.Int
	DestinationChainId *big.Int
	DestinationOwner   common.Address
	SharedId           [32]byte
}

// SharedObjectsERC1155Supply is an auto generated low-level Go binding around an user-defined struct.
type SharedObjectsERC1155Supply struct {
	Id     *big.Int
	Amount *big.Int
}

// RaylsErc1155DvpHandlerMetaData contains all meta data concerning the RaylsErc1155DvpHandler contract.
var RaylsErc1155DvpHandlerMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"GetERCStandard\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.ErcStandard\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"MintFromSwapDvp\",\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_destinationOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_extraDatas\",\"type\":\"tuple[]\",\"internalType\":\"structSharedObjects.Dvp1155ExtraData[]\",\"components\":[{\"name\":\"key\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isPublic\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOfBatch\",\"inputs\":[{\"name\":\"accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"ids\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"_from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelSwap\",\"inputs\":[{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_toChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_tokenValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_enygmaResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_enygmaAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositIntoDvp\",\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvpSwapCompleted\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structSharedObjects.DvpSwapCompletedParams\",\"components\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"exists\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAddressByResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllTokenIdsWithSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"supply\",\"type\":\"tuple[]\",\"internalType\":\"structSharedObjects.ERC1155Supply[]\",\"components\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBatchedAllTokenIdsWithSupply\",\"inputs\":[{\"name\":\"startIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"howMany\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"supplies\",\"type\":\"tuple[]\",\"internalType\":\"structSharedObjects.ERC1155Supply[]\",\"components\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEndpointAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEnygmaEventsAdress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNftExtradaData\",\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structSharedObjects.Dvp1155ExtraData[]\",\"components\":[{\"name\":\"key\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isPublic\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNumberAllTokens\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPNCommunicatorAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenExtraData\",\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structSharedObjects.Dvp1155ExtraData[]\",\"components\":[{\"name\":\"key\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isPublic\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"userArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"trusted\",\"type\":\"tuple\",\"internalType\":\"structRaylsTrustedInit\",\"components\":[{\"name\":\"endpoint\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"raylsNodeEndpoint\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"userGovernance\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isTokenLocked\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lockedForDvp\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_extraDatas\",\"type\":\"tuple[]\",\"internalType\":\"structSharedObjects.Dvp1155ExtraData[]\",\"components\":[{\"name\":\"key\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isPublic\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"notifySenderAndReceiverWithPNCommunicator\",\"inputs\":[{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_destChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_statusToSender\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.DvpCommunicatiorStatus\"},{\"name\":\"_statusToReceiver\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.DvpCommunicatiorStatus\"},{\"name\":\"_messageToSender\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_messageToReceiver\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"notifySenderWithPNCommunicator\",\"inputs\":[{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_status\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.DvpCommunicatiorStatus\"},{\"name\":\"_message\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"raylsNodeUserGovernance\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIUserGovernance\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"resourceId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"safeBatchTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"ids\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"values\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSwapValidityTime\",\"inputs\":[{\"name\":\"_validityTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitTokenUpdate\",\"inputs\":[{\"name\":\"updateType\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.BalanceUpdateType\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"swapValidityTime\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"swapWithDvpForEnygma\",\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_tokenValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenDataParam\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_enygmaAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_enygmaResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_destChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_validityTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unlock\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unlockFromDvp\",\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"uri\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawFromDvp\",\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RaylsErc1155DvpTokenCreated\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenLocked\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenRegistrationSubmitted\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenUnlocked\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TransferBatch\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"ids\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"values\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TransferSingle\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"URI\",\"inputs\":[{\"name\":\"value\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC1155InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC1155InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1155InvalidArrayLength\",\"inputs\":[{\"name\":\"idsLength\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"valuesLength\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC1155InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1155InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1155InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1155MissingApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__HubNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"hubStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__PrivacyNodeFrozen\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__PrivacyNodeNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__PublicChainNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"publicChainStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__ResourceNotApproved\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsApp__TokenRegistryNotConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsApp__UnauthorizedTokenRegistry\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenRegistry\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__UserNotRegistered\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsErc1155DvpHandler__SwapValidityOutOfRange\",\"inputs\":[{\"name\":\"provided\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"min\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"max\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
	ID:  "RaylsErc1155DvpHandler",
}

// RaylsErc1155DvpHandler is an auto generated Go binding around an Ethereum contract.
type RaylsErc1155DvpHandler struct {
	abi abi.ABI
}

// NewRaylsErc1155DvpHandler creates a new instance of RaylsErc1155DvpHandler.
func NewRaylsErc1155DvpHandler() *RaylsErc1155DvpHandler {
	parsed, err := RaylsErc1155DvpHandlerMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &RaylsErc1155DvpHandler{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *RaylsErc1155DvpHandler) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackGetERCStandard is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x35abee1a.
//
// Solidity: function GetERCStandard() pure returns(uint8)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackGetERCStandard() []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("GetERCStandard")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetERCStandard is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x35abee1a.
//
// Solidity: function GetERCStandard() pure returns(uint8)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackGetERCStandard(data []byte) (uint8, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("GetERCStandard", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, err
}

// PackMintFromSwapDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x19f4df79.
//
// Solidity: function MintFromSwapDvp(uint256 _tokenId, address _destinationOwner, uint256 _value, bytes data, (string,string,bool)[] _extraDatas) returns()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackMintFromSwapDvp(tokenId *big.Int, destinationOwner common.Address, value *big.Int, data []byte, extraDatas []SharedObjectsDvp1155ExtraData) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("MintFromSwapDvp", tokenId, destinationOwner, value, data, extraDatas)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackAuthority() []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("authority", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackBalanceOf(account common.Address, id *big.Int) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("balanceOf", account, id)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackBalanceOfBatch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackBalanceOfBatch(accounts []common.Address, ids []*big.Int) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("balanceOfBatch", accounts, ids)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackBalanceOfBatch is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackBalanceOfBatch(data []byte) ([]*big.Int, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("balanceOfBatch", data)
	if err != nil {
		return *new([]*big.Int), err
	}
	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	return out0, err
}

// PackBurn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf5298aca.
//
// Solidity: function burn(address _from, uint256 _id, uint256 _value) returns()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackBurn(from common.Address, id *big.Int, value *big.Int) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("burn", from, id, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCancelSwap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb2677ae4.
//
// Solidity: function cancelSwap(bytes32 _sharedId, uint256 _toChainId, uint256 _tokenId, uint256 _tokenValue, bytes32 _enygmaResourceId, uint256 _enygmaAmount) returns()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackCancelSwap(sharedId [32]byte, toChainId *big.Int, tokenId *big.Int, tokenValue *big.Int, enygmaResourceId [32]byte, enygmaAmount *big.Int) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("cancelSwap", sharedId, toChainId, tokenId, tokenValue, enygmaResourceId, enygmaAmount)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDepositIntoDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9f32f69f.
//
// Solidity: function depositIntoDvp(uint256 _tokenId, uint256 _value, bytes _data) returns()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackDepositIntoDvp(tokenId *big.Int, value *big.Int, data []byte) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("depositIntoDvp", tokenId, value, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDvpSwapCompleted is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3f431e73.
//
// Solidity: function dvpSwapCompleted((uint256,uint256,address,bytes32) params, address from, uint256 _value, bytes data) returns()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackDvpSwapCompleted(params SharedObjectsDvpSwapCompletedParams, from common.Address, value *big.Int, data []byte) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("dvpSwapCompleted", params, from, value, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackExists is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f558e79.
//
// Solidity: function exists(uint256 id) view returns(bool)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackExists(id *big.Int) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("exists", id)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackExists is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4f558e79.
//
// Solidity: function exists(uint256 id) view returns(bool)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackExists(data []byte) (bool, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("exists", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackGetAddressByResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackGetAddressByResourceId(resourceId [32]byte) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("getAddressByResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAddressByResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackGetAddressByResourceId(data []byte) (common.Address, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("getAddressByResourceId", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetAllTokenIdsWithSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xef156674.
//
// Solidity: function getAllTokenIdsWithSupply() view returns((uint256,uint256)[] supply)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackGetAllTokenIdsWithSupply() []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("getAllTokenIdsWithSupply")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAllTokenIdsWithSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xef156674.
//
// Solidity: function getAllTokenIdsWithSupply() view returns((uint256,uint256)[] supply)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackGetAllTokenIdsWithSupply(data []byte) ([]SharedObjectsERC1155Supply, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("getAllTokenIdsWithSupply", data)
	if err != nil {
		return *new([]SharedObjectsERC1155Supply), err
	}
	out0 := *abi.ConvertType(out[0], new([]SharedObjectsERC1155Supply)).(*[]SharedObjectsERC1155Supply)
	return out0, err
}

// PackGetBatchedAllTokenIdsWithSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7a87be0f.
//
// Solidity: function getBatchedAllTokenIdsWithSupply(uint256 startIndex, uint256 howMany) view returns((uint256,uint256)[] supplies)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackGetBatchedAllTokenIdsWithSupply(startIndex *big.Int, howMany *big.Int) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("getBatchedAllTokenIdsWithSupply", startIndex, howMany)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetBatchedAllTokenIdsWithSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7a87be0f.
//
// Solidity: function getBatchedAllTokenIdsWithSupply(uint256 startIndex, uint256 howMany) view returns((uint256,uint256)[] supplies)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackGetBatchedAllTokenIdsWithSupply(data []byte) ([]SharedObjectsERC1155Supply, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("getBatchedAllTokenIdsWithSupply", data)
	if err != nil {
		return *new([]SharedObjectsERC1155Supply), err
	}
	out0 := *abi.ConvertType(out[0], new([]SharedObjectsERC1155Supply)).(*[]SharedObjectsERC1155Supply)
	return out0, err
}

// PackGetEndpointAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xce884eb5.
//
// Solidity: function getEndpointAddress() view returns(address)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackGetEndpointAddress() []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("getEndpointAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetEndpointAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xce884eb5.
//
// Solidity: function getEndpointAddress() view returns(address)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackGetEndpointAddress(data []byte) (common.Address, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("getEndpointAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetEnygmaEventsAdress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd0acffc1.
//
// Solidity: function getEnygmaEventsAdress() view returns(address)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackGetEnygmaEventsAdress() []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("getEnygmaEventsAdress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetEnygmaEventsAdress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd0acffc1.
//
// Solidity: function getEnygmaEventsAdress() view returns(address)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackGetEnygmaEventsAdress(data []byte) (common.Address, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("getEnygmaEventsAdress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetNftExtradaData is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0e9a7f73.
//
// Solidity: function getNftExtradaData(uint256 _tokenId) view returns((string,string,bool)[])
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackGetNftExtradaData(tokenId *big.Int) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("getNftExtradaData", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetNftExtradaData is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0e9a7f73.
//
// Solidity: function getNftExtradaData(uint256 _tokenId) view returns((string,string,bool)[])
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackGetNftExtradaData(data []byte) ([]SharedObjectsDvp1155ExtraData, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("getNftExtradaData", data)
	if err != nil {
		return *new([]SharedObjectsDvp1155ExtraData), err
	}
	out0 := *abi.ConvertType(out[0], new([]SharedObjectsDvp1155ExtraData)).(*[]SharedObjectsDvp1155ExtraData)
	return out0, err
}

// PackGetNumberAllTokens is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8e368ab0.
//
// Solidity: function getNumberAllTokens() view returns(uint256)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackGetNumberAllTokens() []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("getNumberAllTokens")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetNumberAllTokens is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8e368ab0.
//
// Solidity: function getNumberAllTokens() view returns(uint256)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackGetNumberAllTokens(data []byte) (*big.Int, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("getNumberAllTokens", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetPNCommunicatorAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb417c79b.
//
// Solidity: function getPNCommunicatorAddress() view returns(address)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackGetPNCommunicatorAddress() []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("getPNCommunicatorAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetPNCommunicatorAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb417c79b.
//
// Solidity: function getPNCommunicatorAddress() view returns(address)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackGetPNCommunicatorAddress(data []byte) (common.Address, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("getPNCommunicatorAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetTokenExtraData is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3f567b31.
//
// Solidity: function getTokenExtraData(uint256 _tokenId) view returns((string,string,bool)[])
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackGetTokenExtraData(tokenId *big.Int) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("getTokenExtraData", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTokenExtraData is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3f567b31.
//
// Solidity: function getTokenExtraData(uint256 _tokenId) view returns((string,string,bool)[])
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackGetTokenExtraData(data []byte) ([]SharedObjectsDvp1155ExtraData, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("getTokenExtraData", data)
	if err != nil {
		return *new([]SharedObjectsDvp1155ExtraData), err
	}
	out0 := *abi.ConvertType(out[0], new([]SharedObjectsDvp1155ExtraData)).(*[]SharedObjectsDvp1155ExtraData)
	return out0, err
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1d3f35b0.
//
// Solidity: function initialize(bytes userArgs, (address,address,address,address,bytes32,address) trusted) returns()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackInitialize(userArgs []byte, trusted RaylsTrustedInit) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("initialize", userArgs, trusted)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackIsApprovedForAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackIsApprovedForAll(account common.Address, operator common.Address) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("isApprovedForAll", account, operator)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsApprovedForAll is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackIsApprovedForAll(data []byte) (bool, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("isApprovedForAll", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackIsTokenLocked is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7ae1ccf9.
//
// Solidity: function isTokenLocked(address account, uint256 id) view returns(bool)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackIsTokenLocked(account common.Address, id *big.Int) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("isTokenLocked", account, id)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsTokenLocked is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7ae1ccf9.
//
// Solidity: function isTokenLocked(address account, uint256 id) view returns(bool)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackIsTokenLocked(data []byte) (bool, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("isTokenLocked", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackLockedForDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf7eb1445.
//
// Solidity: function lockedForDvp(address , uint256 ) view returns(uint256)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackLockedForDvp(arg0 common.Address, arg1 *big.Int) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("lockedForDvp", arg0, arg1)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackLockedForDvp is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf7eb1445.
//
// Solidity: function lockedForDvp(address , uint256 ) view returns(uint256)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackLockedForDvp(data []byte) (*big.Int, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("lockedForDvp", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6c6a848b.
//
// Solidity: function mint(address _to, uint256 _id, uint256 value, bytes data, (string,string,bool)[] _extraDatas) returns()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackMint(to common.Address, id *big.Int, value *big.Int, data []byte, extraDatas []SharedObjectsDvp1155ExtraData) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("mint", to, id, value, data, extraDatas)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackName() []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("name")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackName(data []byte) (string, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("name", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackNotifySenderAndReceiverWithPNCommunicator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbe4f2cac.
//
// Solidity: function notifySenderAndReceiverWithPNCommunicator(bytes32 _sharedId, uint256 _destChainId, uint8 _statusToSender, uint8 _statusToReceiver, string _messageToSender, string _messageToReceiver) returns()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackNotifySenderAndReceiverWithPNCommunicator(sharedId [32]byte, destChainId *big.Int, statusToSender uint8, statusToReceiver uint8, messageToSender string, messageToReceiver string) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("notifySenderAndReceiverWithPNCommunicator", sharedId, destChainId, statusToSender, statusToReceiver, messageToSender, messageToReceiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackNotifySenderWithPNCommunicator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcc90e43b.
//
// Solidity: function notifySenderWithPNCommunicator(bytes32 _sharedId, uint8 _status, string _message) returns()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackNotifySenderWithPNCommunicator(sharedId [32]byte, status uint8, message string) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("notifySenderWithPNCommunicator", sharedId, status, message)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRaylsNodeUserGovernance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0274c133.
//
// Solidity: function raylsNodeUserGovernance() view returns(address)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackRaylsNodeUserGovernance() []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("raylsNodeUserGovernance")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackRaylsNodeUserGovernance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0274c133.
//
// Solidity: function raylsNodeUserGovernance() view returns(address)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackRaylsNodeUserGovernance(data []byte) (common.Address, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("raylsNodeUserGovernance", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackResourceId() []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("resourceId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackResourceId(data []byte) ([32]byte, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("resourceId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackSafeBatchTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] values, bytes data) returns()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackSafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, values []*big.Int, data []byte) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("safeBatchTransferFrom", from, to, ids, values, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSafeTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 value, bytes data) returns()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackSafeTransferFrom(from common.Address, to common.Address, id *big.Int, value *big.Int, data []byte) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("safeTransferFrom", from, to, id, value, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetApprovalForAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackSetApprovalForAll(operator common.Address, approved bool) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("setApprovalForAll", operator, approved)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa01afbfb.
//
// Solidity: function setResourceId(bytes32 _resourceId) returns()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackSetResourceId(resourceId [32]byte) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("setResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetSwapValidityTime is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5ddbe8f0.
//
// Solidity: function setSwapValidityTime(uint64 _validityTime) returns()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackSetSwapValidityTime(validityTime uint64) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("setSwapValidityTime", validityTime)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSubmitTokenUpdate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfb3ea8f3.
//
// Solidity: function submitTokenUpdate(uint8 updateType, uint256 tokenId, uint256 value) returns()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackSubmitTokenUpdate(updateType uint8, tokenId *big.Int, value *big.Int) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("submitTokenUpdate", updateType, tokenId, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("supportsInterface", interfaceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackSwapValidityTime is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x521d9498.
//
// Solidity: function swapValidityTime() view returns(uint64)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackSwapValidityTime() []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("swapValidityTime")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSwapValidityTime is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x521d9498.
//
// Solidity: function swapValidityTime() view returns(uint64)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackSwapValidityTime(data []byte) (uint64, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("swapValidityTime", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, err
}

// PackSwapWithDvpForEnygma is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd061a68.
//
// Solidity: function swapWithDvpForEnygma(uint256 _tokenId, uint256 _tokenValue, bytes tokenDataParam, uint256 _enygmaAmount, bytes32 _enygmaResourceId, uint256 _destChainId, bytes32 _sharedId, uint64 _validityTime) returns()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackSwapWithDvpForEnygma(tokenId *big.Int, tokenValue *big.Int, tokenDataParam []byte, enygmaAmount *big.Int, enygmaResourceId [32]byte, destChainId *big.Int, sharedId [32]byte, validityTime uint64) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("swapWithDvpForEnygma", tokenId, tokenValue, tokenDataParam, enygmaAmount, enygmaResourceId, destChainId, sharedId, validityTime)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackTotalSupply() []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("totalSupply")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackTotalSupply(data []byte) (*big.Int, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("totalSupply", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackTotalSupply0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbd85b039.
//
// Solidity: function totalSupply(uint256 id) view returns(uint256)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackTotalSupply0(id *big.Int) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("totalSupply0", id)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTotalSupply0 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbd85b039.
//
// Solidity: function totalSupply(uint256 id) view returns(uint256)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackTotalSupply0(data []byte) (*big.Int, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("totalSupply0", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackUnlock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe70cc6aa.
//
// Solidity: function unlock(address to, uint256 id, uint256 value, bytes data) returns(bool)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackUnlock(to common.Address, id *big.Int, value *big.Int, data []byte) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("unlock", to, id, value, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUnlock is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe70cc6aa.
//
// Solidity: function unlock(address to, uint256 id, uint256 value, bytes data) returns(bool)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackUnlock(data []byte) (bool, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("unlock", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackUnlockFromDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xef210bd3.
//
// Solidity: function unlockFromDvp(uint256 _tokenId, uint256 _value, address _to) returns()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackUnlockFromDvp(tokenId *big.Int, value *big.Int, to common.Address) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("unlockFromDvp", tokenId, value, to)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUri is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0e89341c.
//
// Solidity: function uri(uint256 ) view returns(string)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackUri(arg0 *big.Int) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("uri", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUri is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0e89341c.
//
// Solidity: function uri(uint256 ) view returns(string)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackUri(data []byte) (string, error) {
	out, err := raylsErc1155DvpHandler.abi.Unpack("uri", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackWithdrawFromDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9b02b9e0.
//
// Solidity: function withdrawFromDvp(uint256 _tokenId, uint256 _value, bytes data) returns()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) PackWithdrawFromDvp(tokenId *big.Int, value *big.Int, data []byte) []byte {
	enc, err := raylsErc1155DvpHandler.abi.Pack("withdrawFromDvp", tokenId, value, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// RaylsErc1155DvpHandlerApprovalForAll represents a ApprovalForAll event raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerApprovalForAll struct {
	Account  common.Address
	Operator common.Address
	Approved bool
	Raw      *types.Log // Blockchain specific contextual infos
}

const RaylsErc1155DvpHandlerApprovalForAllEventName = "ApprovalForAll"

// ContractEventName returns the user-defined event name.
func (RaylsErc1155DvpHandlerApprovalForAll) ContractEventName() string {
	return RaylsErc1155DvpHandlerApprovalForAllEventName
}

// UnpackApprovalForAllEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackApprovalForAllEvent(log *types.Log) (*RaylsErc1155DvpHandlerApprovalForAll, error) {
	event := "ApprovalForAll"
	if log.Topics[0] != raylsErc1155DvpHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsErc1155DvpHandlerApprovalForAll)
	if len(log.Data) > 0 {
		if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsErc1155DvpHandler.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// RaylsErc1155DvpHandlerAuthorityUpdated represents a AuthorityUpdated event raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerAuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const RaylsErc1155DvpHandlerAuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (RaylsErc1155DvpHandlerAuthorityUpdated) ContractEventName() string {
	return RaylsErc1155DvpHandlerAuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackAuthorityUpdatedEvent(log *types.Log) (*RaylsErc1155DvpHandlerAuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != raylsErc1155DvpHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsErc1155DvpHandlerAuthorityUpdated)
	if len(log.Data) > 0 {
		if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsErc1155DvpHandler.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// RaylsErc1155DvpHandlerInitialized represents a Initialized event raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerInitialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const RaylsErc1155DvpHandlerInitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (RaylsErc1155DvpHandlerInitialized) ContractEventName() string {
	return RaylsErc1155DvpHandlerInitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackInitializedEvent(log *types.Log) (*RaylsErc1155DvpHandlerInitialized, error) {
	event := "Initialized"
	if log.Topics[0] != raylsErc1155DvpHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsErc1155DvpHandlerInitialized)
	if len(log.Data) > 0 {
		if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsErc1155DvpHandler.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// RaylsErc1155DvpHandlerRaylsErc1155DvpTokenCreated represents a RaylsErc1155DvpTokenCreated event raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerRaylsErc1155DvpTokenCreated struct {
	TokenAddress common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const RaylsErc1155DvpHandlerRaylsErc1155DvpTokenCreatedEventName = "RaylsErc1155DvpTokenCreated"

// ContractEventName returns the user-defined event name.
func (RaylsErc1155DvpHandlerRaylsErc1155DvpTokenCreated) ContractEventName() string {
	return RaylsErc1155DvpHandlerRaylsErc1155DvpTokenCreatedEventName
}

// UnpackRaylsErc1155DvpTokenCreatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RaylsErc1155DvpTokenCreated(address indexed tokenAddress)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackRaylsErc1155DvpTokenCreatedEvent(log *types.Log) (*RaylsErc1155DvpHandlerRaylsErc1155DvpTokenCreated, error) {
	event := "RaylsErc1155DvpTokenCreated"
	if log.Topics[0] != raylsErc1155DvpHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsErc1155DvpHandlerRaylsErc1155DvpTokenCreated)
	if len(log.Data) > 0 {
		if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsErc1155DvpHandler.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// RaylsErc1155DvpHandlerTokenLocked represents a TokenLocked event raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerTokenLocked struct {
	Account common.Address
	TokenId *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const RaylsErc1155DvpHandlerTokenLockedEventName = "TokenLocked"

// ContractEventName returns the user-defined event name.
func (RaylsErc1155DvpHandlerTokenLocked) ContractEventName() string {
	return RaylsErc1155DvpHandlerTokenLockedEventName
}

// UnpackTokenLockedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenLocked(address indexed account, uint256 indexed tokenId)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackTokenLockedEvent(log *types.Log) (*RaylsErc1155DvpHandlerTokenLocked, error) {
	event := "TokenLocked"
	if log.Topics[0] != raylsErc1155DvpHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsErc1155DvpHandlerTokenLocked)
	if len(log.Data) > 0 {
		if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsErc1155DvpHandler.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// RaylsErc1155DvpHandlerTokenRegistrationSubmitted represents a TokenRegistrationSubmitted event raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerTokenRegistrationSubmitted struct {
	TokenAddress common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const RaylsErc1155DvpHandlerTokenRegistrationSubmittedEventName = "TokenRegistrationSubmitted"

// ContractEventName returns the user-defined event name.
func (RaylsErc1155DvpHandlerTokenRegistrationSubmitted) ContractEventName() string {
	return RaylsErc1155DvpHandlerTokenRegistrationSubmittedEventName
}

// UnpackTokenRegistrationSubmittedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenRegistrationSubmitted(address indexed tokenAddress)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackTokenRegistrationSubmittedEvent(log *types.Log) (*RaylsErc1155DvpHandlerTokenRegistrationSubmitted, error) {
	event := "TokenRegistrationSubmitted"
	if log.Topics[0] != raylsErc1155DvpHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsErc1155DvpHandlerTokenRegistrationSubmitted)
	if len(log.Data) > 0 {
		if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsErc1155DvpHandler.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// RaylsErc1155DvpHandlerTokenUnlocked represents a TokenUnlocked event raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerTokenUnlocked struct {
	Account common.Address
	TokenId *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const RaylsErc1155DvpHandlerTokenUnlockedEventName = "TokenUnlocked"

// ContractEventName returns the user-defined event name.
func (RaylsErc1155DvpHandlerTokenUnlocked) ContractEventName() string {
	return RaylsErc1155DvpHandlerTokenUnlockedEventName
}

// UnpackTokenUnlockedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenUnlocked(address indexed account, uint256 indexed tokenId)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackTokenUnlockedEvent(log *types.Log) (*RaylsErc1155DvpHandlerTokenUnlocked, error) {
	event := "TokenUnlocked"
	if log.Topics[0] != raylsErc1155DvpHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsErc1155DvpHandlerTokenUnlocked)
	if len(log.Data) > 0 {
		if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsErc1155DvpHandler.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// RaylsErc1155DvpHandlerTransferBatch represents a TransferBatch event raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerTransferBatch struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Ids      []*big.Int
	Values   []*big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const RaylsErc1155DvpHandlerTransferBatchEventName = "TransferBatch"

// ContractEventName returns the user-defined event name.
func (RaylsErc1155DvpHandlerTransferBatch) ContractEventName() string {
	return RaylsErc1155DvpHandlerTransferBatchEventName
}

// UnpackTransferBatchEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackTransferBatchEvent(log *types.Log) (*RaylsErc1155DvpHandlerTransferBatch, error) {
	event := "TransferBatch"
	if log.Topics[0] != raylsErc1155DvpHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsErc1155DvpHandlerTransferBatch)
	if len(log.Data) > 0 {
		if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsErc1155DvpHandler.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// RaylsErc1155DvpHandlerTransferSingle represents a TransferSingle event raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerTransferSingle struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Id       *big.Int
	Value    *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const RaylsErc1155DvpHandlerTransferSingleEventName = "TransferSingle"

// ContractEventName returns the user-defined event name.
func (RaylsErc1155DvpHandlerTransferSingle) ContractEventName() string {
	return RaylsErc1155DvpHandlerTransferSingleEventName
}

// UnpackTransferSingleEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackTransferSingleEvent(log *types.Log) (*RaylsErc1155DvpHandlerTransferSingle, error) {
	event := "TransferSingle"
	if log.Topics[0] != raylsErc1155DvpHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsErc1155DvpHandlerTransferSingle)
	if len(log.Data) > 0 {
		if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsErc1155DvpHandler.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// RaylsErc1155DvpHandlerURI represents a URI event raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerURI struct {
	Value string
	Id    *big.Int
	Raw   *types.Log // Blockchain specific contextual infos
}

const RaylsErc1155DvpHandlerURIEventName = "URI"

// ContractEventName returns the user-defined event name.
func (RaylsErc1155DvpHandlerURI) ContractEventName() string {
	return RaylsErc1155DvpHandlerURIEventName
}

// UnpackURIEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event URI(string value, uint256 indexed id)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackURIEvent(log *types.Log) (*RaylsErc1155DvpHandlerURI, error) {
	event := "URI"
	if log.Topics[0] != raylsErc1155DvpHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsErc1155DvpHandlerURI)
	if len(log.Data) > 0 {
		if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsErc1155DvpHandler.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["ERC1155InsufficientBalance"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackERC1155InsufficientBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["ERC1155InvalidApprover"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackERC1155InvalidApproverError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["ERC1155InvalidArrayLength"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackERC1155InvalidArrayLengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["ERC1155InvalidOperator"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackERC1155InvalidOperatorError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["ERC1155InvalidReceiver"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackERC1155InvalidReceiverError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["ERC1155InvalidSender"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackERC1155InvalidSenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["ERC1155MissingApprovalForAll"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackERC1155MissingApprovalForAllError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["RaylsAppHubNotActive"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackRaylsAppHubNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["RaylsAppPrivacyNodeFrozen"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackRaylsAppPrivacyNodeFrozenError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["RaylsAppPrivacyNodeNotActive"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackRaylsAppPrivacyNodeNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["RaylsAppPublicChainNotActive"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackRaylsAppPublicChainNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["RaylsAppResourceNotApproved"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackRaylsAppResourceNotApprovedError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["RaylsAppTokenRegistryNotConfigured"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackRaylsAppTokenRegistryNotConfiguredError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["RaylsAppUnauthorizedTokenRegistry"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackRaylsAppUnauthorizedTokenRegistryError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["RaylsAppUserNotRegistered"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackRaylsAppUserNotRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["RaylsErc1155DvpHandlerSwapValidityOutOfRange"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackRaylsErc1155DvpHandlerSwapValidityOutOfRangeError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc1155DvpHandler.abi.Errors["ReentrancyGuardReentrantCall"].ID.Bytes()[:4]) {
		return raylsErc1155DvpHandler.UnpackReentrancyGuardReentrantCallError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// RaylsErc1155DvpHandlerERC1155InsufficientBalance represents a ERC1155InsufficientBalance error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerERC1155InsufficientBalance struct {
	Sender  common.Address
	Balance *big.Int
	Needed  *big.Int
	TokenId *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1155InsufficientBalance(address sender, uint256 balance, uint256 needed, uint256 tokenId)
func RaylsErc1155DvpHandlerERC1155InsufficientBalanceErrorID() common.Hash {
	return common.HexToHash("0x03dee4c573c982787b5f3537d6323ffaca9d864448aa6bd828ada9e5d0837036")
}

// UnpackERC1155InsufficientBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1155InsufficientBalance(address sender, uint256 balance, uint256 needed, uint256 tokenId)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackERC1155InsufficientBalanceError(raw []byte) (*RaylsErc1155DvpHandlerERC1155InsufficientBalance, error) {
	out := new(RaylsErc1155DvpHandlerERC1155InsufficientBalance)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "ERC1155InsufficientBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerERC1155InvalidApprover represents a ERC1155InvalidApprover error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerERC1155InvalidApprover struct {
	Approver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1155InvalidApprover(address approver)
func RaylsErc1155DvpHandlerERC1155InvalidApproverErrorID() common.Hash {
	return common.HexToHash("0x3e31884e33c33ce0039d1905e3c252950ae3b95240f36d4fff81f5ff6752ef99")
}

// UnpackERC1155InvalidApproverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1155InvalidApprover(address approver)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackERC1155InvalidApproverError(raw []byte) (*RaylsErc1155DvpHandlerERC1155InvalidApprover, error) {
	out := new(RaylsErc1155DvpHandlerERC1155InvalidApprover)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "ERC1155InvalidApprover", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerERC1155InvalidArrayLength represents a ERC1155InvalidArrayLength error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerERC1155InvalidArrayLength struct {
	IdsLength    *big.Int
	ValuesLength *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1155InvalidArrayLength(uint256 idsLength, uint256 valuesLength)
func RaylsErc1155DvpHandlerERC1155InvalidArrayLengthErrorID() common.Hash {
	return common.HexToHash("0x5b0599913619cfa5633692652638ed25cafcd079c9beae8c251b12c23dcc83f2")
}

// UnpackERC1155InvalidArrayLengthError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1155InvalidArrayLength(uint256 idsLength, uint256 valuesLength)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackERC1155InvalidArrayLengthError(raw []byte) (*RaylsErc1155DvpHandlerERC1155InvalidArrayLength, error) {
	out := new(RaylsErc1155DvpHandlerERC1155InvalidArrayLength)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "ERC1155InvalidArrayLength", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerERC1155InvalidOperator represents a ERC1155InvalidOperator error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerERC1155InvalidOperator struct {
	Operator common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1155InvalidOperator(address operator)
func RaylsErc1155DvpHandlerERC1155InvalidOperatorErrorID() common.Hash {
	return common.HexToHash("0xced3e10010b9d2aa24827119d0db4a8feec73aea48b4b3e470d8a9f3ff723569")
}

// UnpackERC1155InvalidOperatorError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1155InvalidOperator(address operator)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackERC1155InvalidOperatorError(raw []byte) (*RaylsErc1155DvpHandlerERC1155InvalidOperator, error) {
	out := new(RaylsErc1155DvpHandlerERC1155InvalidOperator)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "ERC1155InvalidOperator", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerERC1155InvalidReceiver represents a ERC1155InvalidReceiver error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerERC1155InvalidReceiver struct {
	Receiver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1155InvalidReceiver(address receiver)
func RaylsErc1155DvpHandlerERC1155InvalidReceiverErrorID() common.Hash {
	return common.HexToHash("0x57f447ceed621d9e134e26de5772c88799abb7322ce2a87f95dce247d47105c6")
}

// UnpackERC1155InvalidReceiverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1155InvalidReceiver(address receiver)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackERC1155InvalidReceiverError(raw []byte) (*RaylsErc1155DvpHandlerERC1155InvalidReceiver, error) {
	out := new(RaylsErc1155DvpHandlerERC1155InvalidReceiver)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "ERC1155InvalidReceiver", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerERC1155InvalidSender represents a ERC1155InvalidSender error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerERC1155InvalidSender struct {
	Sender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1155InvalidSender(address sender)
func RaylsErc1155DvpHandlerERC1155InvalidSenderErrorID() common.Hash {
	return common.HexToHash("0x01a83514e94b34009110b75cac6742ba33bd7c62f18a3616bafea52855d3b175")
}

// UnpackERC1155InvalidSenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1155InvalidSender(address sender)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackERC1155InvalidSenderError(raw []byte) (*RaylsErc1155DvpHandlerERC1155InvalidSender, error) {
	out := new(RaylsErc1155DvpHandlerERC1155InvalidSender)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "ERC1155InvalidSender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerERC1155MissingApprovalForAll represents a ERC1155MissingApprovalForAll error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerERC1155MissingApprovalForAll struct {
	Operator common.Address
	Owner    common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1155MissingApprovalForAll(address operator, address owner)
func RaylsErc1155DvpHandlerERC1155MissingApprovalForAllErrorID() common.Hash {
	return common.HexToHash("0xe237d922be9fac42efeaaaffb42cc6b57e0ff95d94a1b74daeff69adc7657754")
}

// UnpackERC1155MissingApprovalForAllError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1155MissingApprovalForAll(address operator, address owner)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackERC1155MissingApprovalForAllError(raw []byte) (*RaylsErc1155DvpHandlerERC1155MissingApprovalForAll, error) {
	out := new(RaylsErc1155DvpHandlerERC1155MissingApprovalForAll)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "ERC1155MissingApprovalForAll", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerInvalidInitialization represents a InvalidInitialization error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerInvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func RaylsErc1155DvpHandlerInvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackInvalidInitializationError(raw []byte) (*RaylsErc1155DvpHandlerInvalidInitialization, error) {
	out := new(RaylsErc1155DvpHandlerInvalidInitialization)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerNotInitializing represents a NotInitializing error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerNotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func RaylsErc1155DvpHandlerNotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackNotInitializingError(raw []byte) (*RaylsErc1155DvpHandlerNotInitializing, error) {
	out := new(RaylsErc1155DvpHandlerNotInitializing)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerRaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerRaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func RaylsErc1155DvpHandlerRaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*RaylsErc1155DvpHandlerRaylsAccessManagedContractPaused, error) {
	out := new(RaylsErc1155DvpHandlerRaylsAccessManagedContractPaused)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerRaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerRaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func RaylsErc1155DvpHandlerRaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*RaylsErc1155DvpHandlerRaylsAccessManagedInvalidAuthority, error) {
	out := new(RaylsErc1155DvpHandlerRaylsAccessManagedInvalidAuthority)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerRaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerRaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func RaylsErc1155DvpHandlerRaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*RaylsErc1155DvpHandlerRaylsAccessManagedMustSchedule, error) {
	out := new(RaylsErc1155DvpHandlerRaylsAccessManagedMustSchedule)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerRaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerRaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func RaylsErc1155DvpHandlerRaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*RaylsErc1155DvpHandlerRaylsAccessManagedUnauthorized, error) {
	out := new(RaylsErc1155DvpHandlerRaylsAccessManagedUnauthorized)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerRaylsAppHubNotActive represents a RaylsApp__HubNotActive error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerRaylsAppHubNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
	HubStatus         uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__HubNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 hubStatus)
func RaylsErc1155DvpHandlerRaylsAppHubNotActiveErrorID() common.Hash {
	return common.HexToHash("0xdc2ffb0fada912f0dd1b700d4ea9a9ce47e3ecdd1b7b155d2066b9a022a637c2")
}

// UnpackRaylsAppHubNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__HubNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 hubStatus)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackRaylsAppHubNotActiveError(raw []byte) (*RaylsErc1155DvpHandlerRaylsAppHubNotActive, error) {
	out := new(RaylsErc1155DvpHandlerRaylsAppHubNotActive)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "RaylsAppHubNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerRaylsAppPrivacyNodeFrozen represents a RaylsApp__PrivacyNodeFrozen error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerRaylsAppPrivacyNodeFrozen struct {
	TokenAddress common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__PrivacyNodeFrozen(address tokenAddress)
func RaylsErc1155DvpHandlerRaylsAppPrivacyNodeFrozenErrorID() common.Hash {
	return common.HexToHash("0xcecb8d3ce0d1417038942c9d252e856b5585275082aa5cdbca675fa64d7bfc24")
}

// UnpackRaylsAppPrivacyNodeFrozenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__PrivacyNodeFrozen(address tokenAddress)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackRaylsAppPrivacyNodeFrozenError(raw []byte) (*RaylsErc1155DvpHandlerRaylsAppPrivacyNodeFrozen, error) {
	out := new(RaylsErc1155DvpHandlerRaylsAppPrivacyNodeFrozen)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "RaylsAppPrivacyNodeFrozen", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerRaylsAppPrivacyNodeNotActive represents a RaylsApp__PrivacyNodeNotActive error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerRaylsAppPrivacyNodeNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__PrivacyNodeNotActive(address tokenAddress, uint8 privacyNodeStatus)
func RaylsErc1155DvpHandlerRaylsAppPrivacyNodeNotActiveErrorID() common.Hash {
	return common.HexToHash("0x44c58c43ed8f726e3330349bec7aa7300f000be36837ee0c2cf507d04511e1e8")
}

// UnpackRaylsAppPrivacyNodeNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__PrivacyNodeNotActive(address tokenAddress, uint8 privacyNodeStatus)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackRaylsAppPrivacyNodeNotActiveError(raw []byte) (*RaylsErc1155DvpHandlerRaylsAppPrivacyNodeNotActive, error) {
	out := new(RaylsErc1155DvpHandlerRaylsAppPrivacyNodeNotActive)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "RaylsAppPrivacyNodeNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerRaylsAppPublicChainNotActive represents a RaylsApp__PublicChainNotActive error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerRaylsAppPublicChainNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
	PublicChainStatus uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__PublicChainNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 publicChainStatus)
func RaylsErc1155DvpHandlerRaylsAppPublicChainNotActiveErrorID() common.Hash {
	return common.HexToHash("0xd6e23bd403a5000c9afe5c2ed5202b3ff8e25d8c3644c1f51892016fb18e5ab9")
}

// UnpackRaylsAppPublicChainNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__PublicChainNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 publicChainStatus)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackRaylsAppPublicChainNotActiveError(raw []byte) (*RaylsErc1155DvpHandlerRaylsAppPublicChainNotActive, error) {
	out := new(RaylsErc1155DvpHandlerRaylsAppPublicChainNotActive)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "RaylsAppPublicChainNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerRaylsAppResourceNotApproved represents a RaylsApp__ResourceNotApproved error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerRaylsAppResourceNotApproved struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__ResourceNotApproved()
func RaylsErc1155DvpHandlerRaylsAppResourceNotApprovedErrorID() common.Hash {
	return common.HexToHash("0x970ad4f73c2c200faa068d3d920e2ef40fca6a5338655abcfb5212557edeed6b")
}

// UnpackRaylsAppResourceNotApprovedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__ResourceNotApproved()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackRaylsAppResourceNotApprovedError(raw []byte) (*RaylsErc1155DvpHandlerRaylsAppResourceNotApproved, error) {
	out := new(RaylsErc1155DvpHandlerRaylsAppResourceNotApproved)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "RaylsAppResourceNotApproved", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerRaylsAppTokenRegistryNotConfigured represents a RaylsApp__TokenRegistryNotConfigured error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerRaylsAppTokenRegistryNotConfigured struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__TokenRegistryNotConfigured()
func RaylsErc1155DvpHandlerRaylsAppTokenRegistryNotConfiguredErrorID() common.Hash {
	return common.HexToHash("0x36a41bd1f6f11cd28b716e935a926fb04f66e11a393b38a49bb660640f3b6dbf")
}

// UnpackRaylsAppTokenRegistryNotConfiguredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__TokenRegistryNotConfigured()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackRaylsAppTokenRegistryNotConfiguredError(raw []byte) (*RaylsErc1155DvpHandlerRaylsAppTokenRegistryNotConfigured, error) {
	out := new(RaylsErc1155DvpHandlerRaylsAppTokenRegistryNotConfigured)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "RaylsAppTokenRegistryNotConfigured", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerRaylsAppUnauthorizedTokenRegistry represents a RaylsApp__UnauthorizedTokenRegistry error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerRaylsAppUnauthorizedTokenRegistry struct {
	Caller        common.Address
	TokenRegistry common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__UnauthorizedTokenRegistry(address caller, address tokenRegistry)
func RaylsErc1155DvpHandlerRaylsAppUnauthorizedTokenRegistryErrorID() common.Hash {
	return common.HexToHash("0x061526480acdfaa09331b795496a6c50aaed25a45d9fca4c9d55fad56af8e09c")
}

// UnpackRaylsAppUnauthorizedTokenRegistryError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__UnauthorizedTokenRegistry(address caller, address tokenRegistry)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackRaylsAppUnauthorizedTokenRegistryError(raw []byte) (*RaylsErc1155DvpHandlerRaylsAppUnauthorizedTokenRegistry, error) {
	out := new(RaylsErc1155DvpHandlerRaylsAppUnauthorizedTokenRegistry)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "RaylsAppUnauthorizedTokenRegistry", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerRaylsAppUserNotRegistered represents a RaylsApp__UserNotRegistered error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerRaylsAppUserNotRegistered struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__UserNotRegistered(address caller)
func RaylsErc1155DvpHandlerRaylsAppUserNotRegisteredErrorID() common.Hash {
	return common.HexToHash("0x4c1db902cce08bec31bedc484362fba54949899ac3c0bf0416f3c44af3284baa")
}

// UnpackRaylsAppUserNotRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__UserNotRegistered(address caller)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackRaylsAppUserNotRegisteredError(raw []byte) (*RaylsErc1155DvpHandlerRaylsAppUserNotRegistered, error) {
	out := new(RaylsErc1155DvpHandlerRaylsAppUserNotRegistered)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "RaylsAppUserNotRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerRaylsErc1155DvpHandlerSwapValidityOutOfRange represents a RaylsErc1155DvpHandler__SwapValidityOutOfRange error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerRaylsErc1155DvpHandlerSwapValidityOutOfRange struct {
	Provided uint64
	Min      uint64
	Max      uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsErc1155DvpHandler__SwapValidityOutOfRange(uint64 provided, uint64 min, uint64 max)
func RaylsErc1155DvpHandlerRaylsErc1155DvpHandlerSwapValidityOutOfRangeErrorID() common.Hash {
	return common.HexToHash("0x5cd2b39d2741ff468a769e235e3b365f7e68ec4a510bcee4a5e15af4922eaca1")
}

// UnpackRaylsErc1155DvpHandlerSwapValidityOutOfRangeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsErc1155DvpHandler__SwapValidityOutOfRange(uint64 provided, uint64 min, uint64 max)
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackRaylsErc1155DvpHandlerSwapValidityOutOfRangeError(raw []byte) (*RaylsErc1155DvpHandlerRaylsErc1155DvpHandlerSwapValidityOutOfRange, error) {
	out := new(RaylsErc1155DvpHandlerRaylsErc1155DvpHandlerSwapValidityOutOfRange)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "RaylsErc1155DvpHandlerSwapValidityOutOfRange", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc1155DvpHandlerReentrancyGuardReentrantCall represents a ReentrancyGuardReentrantCall error raised by the RaylsErc1155DvpHandler contract.
type RaylsErc1155DvpHandlerReentrancyGuardReentrantCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ReentrancyGuardReentrantCall()
func RaylsErc1155DvpHandlerReentrancyGuardReentrantCallErrorID() common.Hash {
	return common.HexToHash("0x3ee5aeb571de7fc460830b4d0017439a1ca56fb0bc39062227ade4fe4a24c1ca")
}

// UnpackReentrancyGuardReentrantCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ReentrancyGuardReentrantCall()
func (raylsErc1155DvpHandler *RaylsErc1155DvpHandler) UnpackReentrancyGuardReentrantCallError(raw []byte) (*RaylsErc1155DvpHandlerReentrancyGuardReentrantCall, error) {
	out := new(RaylsErc1155DvpHandlerReentrancyGuardReentrantCall)
	if err := raylsErc1155DvpHandler.abi.UnpackIntoInterface(out, "ReentrancyGuardReentrantCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}
