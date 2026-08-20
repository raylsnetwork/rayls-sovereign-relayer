// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package RaylsErc721DvpHandler

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

// SharedObjectsDvp721ExtraData is an auto generated low-level Go binding around an user-defined struct.
type SharedObjectsDvp721ExtraData struct {
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

// RaylsErc721DvpHandlerMetaData contains all meta data concerning the RaylsErc721DvpHandler contract.
var RaylsErc721DvpHandlerMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"GetERCStandard\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.ErcStandard\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"MintFromSwapDvp\",\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_destinationOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_extraDatas\",\"type\":\"tuple[]\",\"internalType\":\"structSharedObjects.Dvp721ExtraData[]\",\"components\":[{\"name\":\"key\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isPublic\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"_tokenName\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"_tokenSymbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"baseUri\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"_id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelSwap\",\"inputs\":[{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_toChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_nftId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_enygmaAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_enygmaResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositIntoDvp\",\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvpSwapCompleted\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structSharedObjects.DvpSwapCompletedParams\",\"components\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAddressByResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApproved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEndpointAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEnygmaEventsAdress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNftExtradaData\",\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structSharedObjects.Dvp721ExtraData[]\",\"components\":[{\"name\":\"key\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isPublic\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPNCommunicatorAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalSupplyAtPN\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"userArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"trusted\",\"type\":\"tuple\",\"internalType\":\"structRaylsTrustedInit\",\"components\":[{\"name\":\"endpoint\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"raylsNodeEndpoint\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"userGovernance\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isTokenLocked\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lockedForDvp\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_extraDatas\",\"type\":\"tuple[]\",\"internalType\":\"structSharedObjects.Dvp721ExtraData[]\",\"components\":[{\"name\":\"key\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isPublic\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"notifySenderAndReceiverWithPNCommunicator\",\"inputs\":[{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_destChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_statusToSender\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.DvpCommunicatiorStatus\"},{\"name\":\"_statusToReceiver\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.DvpCommunicatiorStatus\"},{\"name\":\"_messageToSender\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_messageToReceiver\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"notifySenderWithPNCommunicator\",\"inputs\":[{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_status\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.DvpCommunicatiorStatus\"},{\"name\":\"_message\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ownerOf\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"raylsNodeUserGovernance\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIUserGovernance\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"resourceId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSwapValidityTime\",\"inputs\":[{\"name\":\"_validityTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitTokenUpdate\",\"inputs\":[{\"name\":\"updateType\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.BalanceUpdateType\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"swapValidityTime\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"swapWithDvpForEnygma\",\"inputs\":[{\"name\":\"_nftId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_enygmaAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_enygmaResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_destChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_validityTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unlock\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unlockFromDvp\",\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawFromDvp\",\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RaylsErc721DvpTokenCreated\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenLocked\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenRegistrationSubmitted\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenUnlocked\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC721IncorrectOwner\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InsufficientApproval\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721NonexistentToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__HubNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"hubStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__PrivacyNodeFrozen\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__PrivacyNodeNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__PublicChainNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"publicChainStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__ResourceNotApproved\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsApp__TokenRegistryNotConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsApp__UnauthorizedTokenRegistry\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenRegistry\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__UserNotRegistered\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsErc721DvpHandler__SwapValidityOutOfRange\",\"inputs\":[{\"name\":\"provided\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"min\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"max\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"RaylsErc721DvpHandler__TokenDoesNotExist\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"RaylsErc721DvpHandler__TokenNotApproved\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
	ID:  "RaylsErc721DvpHandler",
}

// RaylsErc721DvpHandler is an auto generated Go binding around an Ethereum contract.
type RaylsErc721DvpHandler struct {
	abi abi.ABI
}

// NewRaylsErc721DvpHandler creates a new instance of RaylsErc721DvpHandler.
func NewRaylsErc721DvpHandler() *RaylsErc721DvpHandler {
	parsed, err := RaylsErc721DvpHandlerMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &RaylsErc721DvpHandler{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *RaylsErc721DvpHandler) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackGetERCStandard is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x35abee1a.
//
// Solidity: function GetERCStandard() pure returns(uint8)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackGetERCStandard() []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("GetERCStandard")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetERCStandard is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x35abee1a.
//
// Solidity: function GetERCStandard() pure returns(uint8)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackGetERCStandard(data []byte) (uint8, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("GetERCStandard", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, err
}

// PackMintFromSwapDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdb9df168.
//
// Solidity: function MintFromSwapDvp(uint256 _tokenId, address _destinationOwner, (string,string,bool)[] _extraDatas) returns()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackMintFromSwapDvp(tokenId *big.Int, destinationOwner common.Address, extraDatas []SharedObjectsDvp721ExtraData) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("MintFromSwapDvp", tokenId, destinationOwner, extraDatas)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackTokenName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9f0c8125.
//
// Solidity: function _tokenName() view returns(string)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackTokenName() []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("_tokenName")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTokenName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9f0c8125.
//
// Solidity: function _tokenName() view returns(string)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackTokenName(data []byte) (string, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("_tokenName", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackTokenSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xecf35cca.
//
// Solidity: function _tokenSymbol() view returns(string)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackTokenSymbol() []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("_tokenSymbol")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTokenSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xecf35cca.
//
// Solidity: function _tokenSymbol() view returns(string)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackTokenSymbol(data []byte) (string, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("_tokenSymbol", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackApprove(to common.Address, tokenId *big.Int) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("approve", to, tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackAuthority() []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("authority", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackBalanceOf(owner common.Address) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("balanceOf", owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackBaseUri is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9abc8320.
//
// Solidity: function baseUri() view returns(string)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackBaseUri() []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("baseUri")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackBaseUri is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9abc8320.
//
// Solidity: function baseUri() view returns(string)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackBaseUri(data []byte) (string, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("baseUri", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackBurn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42966c68.
//
// Solidity: function burn(uint256 _id) returns()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackBurn(id *big.Int) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("burn", id)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCancelSwap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdc8c47e9.
//
// Solidity: function cancelSwap(bytes32 _sharedId, uint256 _toChainId, uint256 _nftId, uint256 _enygmaAmount, bytes32 _enygmaResourceId) returns()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackCancelSwap(sharedId [32]byte, toChainId *big.Int, nftId *big.Int, enygmaAmount *big.Int, enygmaResourceId [32]byte) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("cancelSwap", sharedId, toChainId, nftId, enygmaAmount, enygmaResourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDepositIntoDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x48df7d1f.
//
// Solidity: function depositIntoDvp(uint256 _tokenId) returns()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackDepositIntoDvp(tokenId *big.Int) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("depositIntoDvp", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDvpSwapCompleted is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5a2c90b8.
//
// Solidity: function dvpSwapCompleted((uint256,uint256,address,bytes32) params) returns()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackDvpSwapCompleted(params SharedObjectsDvpSwapCompletedParams) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("dvpSwapCompleted", params)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackGetAddressByResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackGetAddressByResourceId(resourceId [32]byte) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("getAddressByResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAddressByResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackGetAddressByResourceId(data []byte) (common.Address, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("getAddressByResourceId", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetApproved is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackGetApproved(tokenId *big.Int) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("getApproved", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetApproved is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackGetApproved(data []byte) (common.Address, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("getApproved", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetEndpointAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xce884eb5.
//
// Solidity: function getEndpointAddress() view returns(address)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackGetEndpointAddress() []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("getEndpointAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetEndpointAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xce884eb5.
//
// Solidity: function getEndpointAddress() view returns(address)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackGetEndpointAddress(data []byte) (common.Address, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("getEndpointAddress", data)
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
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackGetEnygmaEventsAdress() []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("getEnygmaEventsAdress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetEnygmaEventsAdress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd0acffc1.
//
// Solidity: function getEnygmaEventsAdress() view returns(address)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackGetEnygmaEventsAdress(data []byte) (common.Address, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("getEnygmaEventsAdress", data)
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
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackGetNftExtradaData(tokenId *big.Int) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("getNftExtradaData", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetNftExtradaData is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0e9a7f73.
//
// Solidity: function getNftExtradaData(uint256 _tokenId) view returns((string,string,bool)[])
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackGetNftExtradaData(data []byte) ([]SharedObjectsDvp721ExtraData, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("getNftExtradaData", data)
	if err != nil {
		return *new([]SharedObjectsDvp721ExtraData), err
	}
	out0 := *abi.ConvertType(out[0], new([]SharedObjectsDvp721ExtraData)).(*[]SharedObjectsDvp721ExtraData)
	return out0, err
}

// PackGetPNCommunicatorAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb417c79b.
//
// Solidity: function getPNCommunicatorAddress() view returns(address)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackGetPNCommunicatorAddress() []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("getPNCommunicatorAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetPNCommunicatorAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb417c79b.
//
// Solidity: function getPNCommunicatorAddress() view returns(address)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackGetPNCommunicatorAddress(data []byte) (common.Address, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("getPNCommunicatorAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc4e41b22.
//
// Solidity: function getTotalSupply() view returns(uint256[])
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackGetTotalSupply() []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("getTotalSupply")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc4e41b22.
//
// Solidity: function getTotalSupply() view returns(uint256[])
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackGetTotalSupply(data []byte) ([]*big.Int, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("getTotalSupply", data)
	if err != nil {
		return *new([]*big.Int), err
	}
	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	return out0, err
}

// PackGetTotalSupplyAtPN is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf83a1bcd.
//
// Solidity: function getTotalSupplyAtPN() view returns(uint256[])
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackGetTotalSupplyAtPN() []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("getTotalSupplyAtPN")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTotalSupplyAtPN is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf83a1bcd.
//
// Solidity: function getTotalSupplyAtPN() view returns(uint256[])
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackGetTotalSupplyAtPN(data []byte) ([]*big.Int, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("getTotalSupplyAtPN", data)
	if err != nil {
		return *new([]*big.Int), err
	}
	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	return out0, err
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1d3f35b0.
//
// Solidity: function initialize(bytes userArgs, (address,address,address,address,bytes32,address) trusted) returns()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackInitialize(userArgs []byte, trusted RaylsTrustedInit) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("initialize", userArgs, trusted)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackIsApprovedForAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackIsApprovedForAll(owner common.Address, operator common.Address) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("isApprovedForAll", owner, operator)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsApprovedForAll is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackIsApprovedForAll(data []byte) (bool, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("isApprovedForAll", data)
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
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackIsTokenLocked(account common.Address, id *big.Int) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("isTokenLocked", account, id)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsTokenLocked is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7ae1ccf9.
//
// Solidity: function isTokenLocked(address account, uint256 id) view returns(bool)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackIsTokenLocked(data []byte) (bool, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("isTokenLocked", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackLockedForDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf7ee1bcb.
//
// Solidity: function lockedForDvp(uint256 ) view returns(bool)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackLockedForDvp(arg0 *big.Int) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("lockedForDvp", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackLockedForDvp is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf7ee1bcb.
//
// Solidity: function lockedForDvp(uint256 ) view returns(bool)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackLockedForDvp(data []byte) (bool, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("lockedForDvp", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x87c5d197.
//
// Solidity: function mint(address _to, uint256 _id, (string,string,bool)[] _extraDatas) returns()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackMint(to common.Address, id *big.Int, extraDatas []SharedObjectsDvp721ExtraData) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("mint", to, id, extraDatas)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackName() []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("name")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackName(data []byte) (string, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("name", data)
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
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackNotifySenderAndReceiverWithPNCommunicator(sharedId [32]byte, destChainId *big.Int, statusToSender uint8, statusToReceiver uint8, messageToSender string, messageToReceiver string) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("notifySenderAndReceiverWithPNCommunicator", sharedId, destChainId, statusToSender, statusToReceiver, messageToSender, messageToReceiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackNotifySenderWithPNCommunicator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcc90e43b.
//
// Solidity: function notifySenderWithPNCommunicator(bytes32 _sharedId, uint8 _status, string _message) returns()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackNotifySenderWithPNCommunicator(sharedId [32]byte, status uint8, message string) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("notifySenderWithPNCommunicator", sharedId, status, message)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackOwnerOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackOwnerOf(tokenId *big.Int) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("ownerOf", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackOwnerOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackOwnerOf(data []byte) (common.Address, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("ownerOf", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackRaylsNodeUserGovernance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0274c133.
//
// Solidity: function raylsNodeUserGovernance() view returns(address)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackRaylsNodeUserGovernance() []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("raylsNodeUserGovernance")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackRaylsNodeUserGovernance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0274c133.
//
// Solidity: function raylsNodeUserGovernance() view returns(address)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackRaylsNodeUserGovernance(data []byte) (common.Address, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("raylsNodeUserGovernance", data)
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
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackResourceId() []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("resourceId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackResourceId(data []byte) ([32]byte, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("resourceId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackSafeTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackSafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("safeTransferFrom", from, to, tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSafeTransferFrom0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackSafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("safeTransferFrom0", from, to, tokenId, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetApprovalForAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackSetApprovalForAll(operator common.Address, approved bool) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("setApprovalForAll", operator, approved)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa01afbfb.
//
// Solidity: function setResourceId(bytes32 _resourceId) returns()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackSetResourceId(resourceId [32]byte) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("setResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetSwapValidityTime is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5ddbe8f0.
//
// Solidity: function setSwapValidityTime(uint64 _validityTime) returns()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackSetSwapValidityTime(validityTime uint64) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("setSwapValidityTime", validityTime)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSubmitTokenUpdate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0e5ee780.
//
// Solidity: function submitTokenUpdate(uint8 updateType, uint256 tokenId) returns()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackSubmitTokenUpdate(updateType uint8, tokenId *big.Int) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("submitTokenUpdate", updateType, tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("supportsInterface", interfaceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("supportsInterface", data)
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
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackSwapValidityTime() []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("swapValidityTime")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSwapValidityTime is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x521d9498.
//
// Solidity: function swapValidityTime() view returns(uint64)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackSwapValidityTime(data []byte) (uint64, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("swapValidityTime", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, err
}

// PackSwapWithDvpForEnygma is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x172ddec4.
//
// Solidity: function swapWithDvpForEnygma(uint256 _nftId, uint256 _enygmaAmount, bytes32 _enygmaResourceId, uint256 _destChainId, bytes32 _sharedId, uint64 _validityTime) returns()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackSwapWithDvpForEnygma(nftId *big.Int, enygmaAmount *big.Int, enygmaResourceId [32]byte, destChainId *big.Int, sharedId [32]byte, validityTime uint64) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("swapWithDvpForEnygma", nftId, enygmaAmount, enygmaResourceId, destChainId, sharedId, validityTime)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackSymbol() []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("symbol")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackSymbol(data []byte) (string, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("symbol", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackTokenURI is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackTokenURI(tokenId *big.Int) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("tokenURI", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTokenURI is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackTokenURI(data []byte) (string, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("tokenURI", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackTransferFrom(from common.Address, to common.Address, tokenId *big.Int) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("transferFrom", from, to, tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUnlock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7eee288d.
//
// Solidity: function unlock(address to, uint256 id) returns(bool)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackUnlock(to common.Address, id *big.Int) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("unlock", to, id)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUnlock is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7eee288d.
//
// Solidity: function unlock(address to, uint256 id) returns(bool)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackUnlock(data []byte) (bool, error) {
	out, err := raylsErc721DvpHandler.abi.Unpack("unlock", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackUnlockFromDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x590e613d.
//
// Solidity: function unlockFromDvp(uint256 _tokenId) returns()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackUnlockFromDvp(tokenId *big.Int) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("unlockFromDvp", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackWithdrawFromDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x30163c12.
//
// Solidity: function withdrawFromDvp(uint256 _tokenId) returns()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) PackWithdrawFromDvp(tokenId *big.Int) []byte {
	enc, err := raylsErc721DvpHandler.abi.Pack("withdrawFromDvp", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// RaylsErc721DvpHandlerApproval represents a Approval event raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const RaylsErc721DvpHandlerApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (RaylsErc721DvpHandlerApproval) ContractEventName() string {
	return RaylsErc721DvpHandlerApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackApprovalEvent(log *types.Log) (*RaylsErc721DvpHandlerApproval, error) {
	event := "Approval"
	if log.Topics[0] != raylsErc721DvpHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsErc721DvpHandlerApproval)
	if len(log.Data) > 0 {
		if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsErc721DvpHandler.abi.Events[event].Inputs {
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

// RaylsErc721DvpHandlerApprovalForAll represents a ApprovalForAll event raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      *types.Log // Blockchain specific contextual infos
}

const RaylsErc721DvpHandlerApprovalForAllEventName = "ApprovalForAll"

// ContractEventName returns the user-defined event name.
func (RaylsErc721DvpHandlerApprovalForAll) ContractEventName() string {
	return RaylsErc721DvpHandlerApprovalForAllEventName
}

// UnpackApprovalForAllEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackApprovalForAllEvent(log *types.Log) (*RaylsErc721DvpHandlerApprovalForAll, error) {
	event := "ApprovalForAll"
	if log.Topics[0] != raylsErc721DvpHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsErc721DvpHandlerApprovalForAll)
	if len(log.Data) > 0 {
		if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsErc721DvpHandler.abi.Events[event].Inputs {
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

// RaylsErc721DvpHandlerAuthorityUpdated represents a AuthorityUpdated event raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerAuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const RaylsErc721DvpHandlerAuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (RaylsErc721DvpHandlerAuthorityUpdated) ContractEventName() string {
	return RaylsErc721DvpHandlerAuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackAuthorityUpdatedEvent(log *types.Log) (*RaylsErc721DvpHandlerAuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != raylsErc721DvpHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsErc721DvpHandlerAuthorityUpdated)
	if len(log.Data) > 0 {
		if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsErc721DvpHandler.abi.Events[event].Inputs {
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

// RaylsErc721DvpHandlerInitialized represents a Initialized event raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerInitialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const RaylsErc721DvpHandlerInitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (RaylsErc721DvpHandlerInitialized) ContractEventName() string {
	return RaylsErc721DvpHandlerInitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackInitializedEvent(log *types.Log) (*RaylsErc721DvpHandlerInitialized, error) {
	event := "Initialized"
	if log.Topics[0] != raylsErc721DvpHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsErc721DvpHandlerInitialized)
	if len(log.Data) > 0 {
		if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsErc721DvpHandler.abi.Events[event].Inputs {
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

// RaylsErc721DvpHandlerRaylsErc721DvpTokenCreated represents a RaylsErc721DvpTokenCreated event raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerRaylsErc721DvpTokenCreated struct {
	TokenAddress common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const RaylsErc721DvpHandlerRaylsErc721DvpTokenCreatedEventName = "RaylsErc721DvpTokenCreated"

// ContractEventName returns the user-defined event name.
func (RaylsErc721DvpHandlerRaylsErc721DvpTokenCreated) ContractEventName() string {
	return RaylsErc721DvpHandlerRaylsErc721DvpTokenCreatedEventName
}

// UnpackRaylsErc721DvpTokenCreatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RaylsErc721DvpTokenCreated(address indexed tokenAddress)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackRaylsErc721DvpTokenCreatedEvent(log *types.Log) (*RaylsErc721DvpHandlerRaylsErc721DvpTokenCreated, error) {
	event := "RaylsErc721DvpTokenCreated"
	if log.Topics[0] != raylsErc721DvpHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsErc721DvpHandlerRaylsErc721DvpTokenCreated)
	if len(log.Data) > 0 {
		if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsErc721DvpHandler.abi.Events[event].Inputs {
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

// RaylsErc721DvpHandlerTokenLocked represents a TokenLocked event raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerTokenLocked struct {
	Account common.Address
	TokenId *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const RaylsErc721DvpHandlerTokenLockedEventName = "TokenLocked"

// ContractEventName returns the user-defined event name.
func (RaylsErc721DvpHandlerTokenLocked) ContractEventName() string {
	return RaylsErc721DvpHandlerTokenLockedEventName
}

// UnpackTokenLockedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenLocked(address indexed account, uint256 indexed tokenId)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackTokenLockedEvent(log *types.Log) (*RaylsErc721DvpHandlerTokenLocked, error) {
	event := "TokenLocked"
	if log.Topics[0] != raylsErc721DvpHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsErc721DvpHandlerTokenLocked)
	if len(log.Data) > 0 {
		if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsErc721DvpHandler.abi.Events[event].Inputs {
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

// RaylsErc721DvpHandlerTokenRegistrationSubmitted represents a TokenRegistrationSubmitted event raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerTokenRegistrationSubmitted struct {
	TokenAddress common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const RaylsErc721DvpHandlerTokenRegistrationSubmittedEventName = "TokenRegistrationSubmitted"

// ContractEventName returns the user-defined event name.
func (RaylsErc721DvpHandlerTokenRegistrationSubmitted) ContractEventName() string {
	return RaylsErc721DvpHandlerTokenRegistrationSubmittedEventName
}

// UnpackTokenRegistrationSubmittedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenRegistrationSubmitted(address indexed tokenAddress)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackTokenRegistrationSubmittedEvent(log *types.Log) (*RaylsErc721DvpHandlerTokenRegistrationSubmitted, error) {
	event := "TokenRegistrationSubmitted"
	if log.Topics[0] != raylsErc721DvpHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsErc721DvpHandlerTokenRegistrationSubmitted)
	if len(log.Data) > 0 {
		if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsErc721DvpHandler.abi.Events[event].Inputs {
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

// RaylsErc721DvpHandlerTokenUnlocked represents a TokenUnlocked event raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerTokenUnlocked struct {
	Account common.Address
	TokenId *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const RaylsErc721DvpHandlerTokenUnlockedEventName = "TokenUnlocked"

// ContractEventName returns the user-defined event name.
func (RaylsErc721DvpHandlerTokenUnlocked) ContractEventName() string {
	return RaylsErc721DvpHandlerTokenUnlockedEventName
}

// UnpackTokenUnlockedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenUnlocked(address indexed account, uint256 indexed tokenId)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackTokenUnlockedEvent(log *types.Log) (*RaylsErc721DvpHandlerTokenUnlocked, error) {
	event := "TokenUnlocked"
	if log.Topics[0] != raylsErc721DvpHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsErc721DvpHandlerTokenUnlocked)
	if len(log.Data) > 0 {
		if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsErc721DvpHandler.abi.Events[event].Inputs {
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

// RaylsErc721DvpHandlerTransfer represents a Transfer event raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const RaylsErc721DvpHandlerTransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (RaylsErc721DvpHandlerTransfer) ContractEventName() string {
	return RaylsErc721DvpHandlerTransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackTransferEvent(log *types.Log) (*RaylsErc721DvpHandlerTransfer, error) {
	event := "Transfer"
	if log.Topics[0] != raylsErc721DvpHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsErc721DvpHandlerTransfer)
	if len(log.Data) > 0 {
		if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsErc721DvpHandler.abi.Events[event].Inputs {
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
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["ERC721IncorrectOwner"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackERC721IncorrectOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["ERC721InsufficientApproval"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackERC721InsufficientApprovalError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["ERC721InvalidApprover"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackERC721InvalidApproverError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["ERC721InvalidOperator"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackERC721InvalidOperatorError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["ERC721InvalidOwner"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackERC721InvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["ERC721InvalidReceiver"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackERC721InvalidReceiverError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["ERC721InvalidSender"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackERC721InvalidSenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["ERC721NonexistentToken"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackERC721NonexistentTokenError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["RaylsAppHubNotActive"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackRaylsAppHubNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["RaylsAppPrivacyNodeFrozen"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackRaylsAppPrivacyNodeFrozenError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["RaylsAppPrivacyNodeNotActive"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackRaylsAppPrivacyNodeNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["RaylsAppPublicChainNotActive"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackRaylsAppPublicChainNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["RaylsAppResourceNotApproved"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackRaylsAppResourceNotApprovedError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["RaylsAppTokenRegistryNotConfigured"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackRaylsAppTokenRegistryNotConfiguredError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["RaylsAppUnauthorizedTokenRegistry"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackRaylsAppUnauthorizedTokenRegistryError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["RaylsAppUserNotRegistered"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackRaylsAppUserNotRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["RaylsErc721DvpHandlerSwapValidityOutOfRange"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackRaylsErc721DvpHandlerSwapValidityOutOfRangeError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["RaylsErc721DvpHandlerTokenDoesNotExist"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackRaylsErc721DvpHandlerTokenDoesNotExistError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["RaylsErc721DvpHandlerTokenNotApproved"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackRaylsErc721DvpHandlerTokenNotApprovedError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsErc721DvpHandler.abi.Errors["ReentrancyGuardReentrantCall"].ID.Bytes()[:4]) {
		return raylsErc721DvpHandler.UnpackReentrancyGuardReentrantCallError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// RaylsErc721DvpHandlerERC721IncorrectOwner represents a ERC721IncorrectOwner error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerERC721IncorrectOwner struct {
	Sender  common.Address
	TokenId *big.Int
	Owner   common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721IncorrectOwner(address sender, uint256 tokenId, address owner)
func RaylsErc721DvpHandlerERC721IncorrectOwnerErrorID() common.Hash {
	return common.HexToHash("0x64283d7b313c8117c125f736876fa2b4e90ea3831a4716dfdb87d2f540e26289")
}

// UnpackERC721IncorrectOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721IncorrectOwner(address sender, uint256 tokenId, address owner)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackERC721IncorrectOwnerError(raw []byte) (*RaylsErc721DvpHandlerERC721IncorrectOwner, error) {
	out := new(RaylsErc721DvpHandlerERC721IncorrectOwner)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "ERC721IncorrectOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerERC721InsufficientApproval represents a ERC721InsufficientApproval error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerERC721InsufficientApproval struct {
	Operator common.Address
	TokenId  *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InsufficientApproval(address operator, uint256 tokenId)
func RaylsErc721DvpHandlerERC721InsufficientApprovalErrorID() common.Hash {
	return common.HexToHash("0x177e802f6f313bc89797ecace66d6d29ab4719cbaaacbb87367264048b1eb861")
}

// UnpackERC721InsufficientApprovalError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InsufficientApproval(address operator, uint256 tokenId)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackERC721InsufficientApprovalError(raw []byte) (*RaylsErc721DvpHandlerERC721InsufficientApproval, error) {
	out := new(RaylsErc721DvpHandlerERC721InsufficientApproval)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "ERC721InsufficientApproval", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerERC721InvalidApprover represents a ERC721InvalidApprover error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerERC721InvalidApprover struct {
	Approver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InvalidApprover(address approver)
func RaylsErc721DvpHandlerERC721InvalidApproverErrorID() common.Hash {
	return common.HexToHash("0xa9fbf51f86b8e03595d59dc726bb10c329bb24f62589be276d8dd193ca0b69ea")
}

// UnpackERC721InvalidApproverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InvalidApprover(address approver)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackERC721InvalidApproverError(raw []byte) (*RaylsErc721DvpHandlerERC721InvalidApprover, error) {
	out := new(RaylsErc721DvpHandlerERC721InvalidApprover)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "ERC721InvalidApprover", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerERC721InvalidOperator represents a ERC721InvalidOperator error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerERC721InvalidOperator struct {
	Operator common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InvalidOperator(address operator)
func RaylsErc721DvpHandlerERC721InvalidOperatorErrorID() common.Hash {
	return common.HexToHash("0x5b08ba185e8f577075361f3a3555a6580a227ce22734dcc979c1aeadf894658b")
}

// UnpackERC721InvalidOperatorError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InvalidOperator(address operator)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackERC721InvalidOperatorError(raw []byte) (*RaylsErc721DvpHandlerERC721InvalidOperator, error) {
	out := new(RaylsErc721DvpHandlerERC721InvalidOperator)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "ERC721InvalidOperator", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerERC721InvalidOwner represents a ERC721InvalidOwner error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerERC721InvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InvalidOwner(address owner)
func RaylsErc721DvpHandlerERC721InvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x89c62b6479af2e623826dcc39c5133061d35b66d72de92833401dd2fd6567480")
}

// UnpackERC721InvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InvalidOwner(address owner)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackERC721InvalidOwnerError(raw []byte) (*RaylsErc721DvpHandlerERC721InvalidOwner, error) {
	out := new(RaylsErc721DvpHandlerERC721InvalidOwner)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "ERC721InvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerERC721InvalidReceiver represents a ERC721InvalidReceiver error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerERC721InvalidReceiver struct {
	Receiver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InvalidReceiver(address receiver)
func RaylsErc721DvpHandlerERC721InvalidReceiverErrorID() common.Hash {
	return common.HexToHash("0x64a0ae9278f805eaf991dcd18ca78756d280b7508b764ef1b255c55845c11df9")
}

// UnpackERC721InvalidReceiverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InvalidReceiver(address receiver)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackERC721InvalidReceiverError(raw []byte) (*RaylsErc721DvpHandlerERC721InvalidReceiver, error) {
	out := new(RaylsErc721DvpHandlerERC721InvalidReceiver)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "ERC721InvalidReceiver", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerERC721InvalidSender represents a ERC721InvalidSender error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerERC721InvalidSender struct {
	Sender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InvalidSender(address sender)
func RaylsErc721DvpHandlerERC721InvalidSenderErrorID() common.Hash {
	return common.HexToHash("0x73c6ac6e10798e95d99e1f130d923eb40193ecb8d094ec3dce93292564eb3b17")
}

// UnpackERC721InvalidSenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InvalidSender(address sender)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackERC721InvalidSenderError(raw []byte) (*RaylsErc721DvpHandlerERC721InvalidSender, error) {
	out := new(RaylsErc721DvpHandlerERC721InvalidSender)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "ERC721InvalidSender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerERC721NonexistentToken represents a ERC721NonexistentToken error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerERC721NonexistentToken struct {
	TokenId *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721NonexistentToken(uint256 tokenId)
func RaylsErc721DvpHandlerERC721NonexistentTokenErrorID() common.Hash {
	return common.HexToHash("0x7e273289a3a9ef6670f06df7dca227856fc925e956db96980692764a8bc734d7")
}

// UnpackERC721NonexistentTokenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721NonexistentToken(uint256 tokenId)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackERC721NonexistentTokenError(raw []byte) (*RaylsErc721DvpHandlerERC721NonexistentToken, error) {
	out := new(RaylsErc721DvpHandlerERC721NonexistentToken)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "ERC721NonexistentToken", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerInvalidInitialization represents a InvalidInitialization error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerInvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func RaylsErc721DvpHandlerInvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackInvalidInitializationError(raw []byte) (*RaylsErc721DvpHandlerInvalidInitialization, error) {
	out := new(RaylsErc721DvpHandlerInvalidInitialization)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerNotInitializing represents a NotInitializing error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerNotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func RaylsErc721DvpHandlerNotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackNotInitializingError(raw []byte) (*RaylsErc721DvpHandlerNotInitializing, error) {
	out := new(RaylsErc721DvpHandlerNotInitializing)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerRaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerRaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func RaylsErc721DvpHandlerRaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*RaylsErc721DvpHandlerRaylsAccessManagedContractPaused, error) {
	out := new(RaylsErc721DvpHandlerRaylsAccessManagedContractPaused)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerRaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerRaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func RaylsErc721DvpHandlerRaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*RaylsErc721DvpHandlerRaylsAccessManagedInvalidAuthority, error) {
	out := new(RaylsErc721DvpHandlerRaylsAccessManagedInvalidAuthority)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerRaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerRaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func RaylsErc721DvpHandlerRaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*RaylsErc721DvpHandlerRaylsAccessManagedMustSchedule, error) {
	out := new(RaylsErc721DvpHandlerRaylsAccessManagedMustSchedule)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerRaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerRaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func RaylsErc721DvpHandlerRaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*RaylsErc721DvpHandlerRaylsAccessManagedUnauthorized, error) {
	out := new(RaylsErc721DvpHandlerRaylsAccessManagedUnauthorized)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerRaylsAppHubNotActive represents a RaylsApp__HubNotActive error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerRaylsAppHubNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
	HubStatus         uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__HubNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 hubStatus)
func RaylsErc721DvpHandlerRaylsAppHubNotActiveErrorID() common.Hash {
	return common.HexToHash("0xdc2ffb0fada912f0dd1b700d4ea9a9ce47e3ecdd1b7b155d2066b9a022a637c2")
}

// UnpackRaylsAppHubNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__HubNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 hubStatus)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackRaylsAppHubNotActiveError(raw []byte) (*RaylsErc721DvpHandlerRaylsAppHubNotActive, error) {
	out := new(RaylsErc721DvpHandlerRaylsAppHubNotActive)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "RaylsAppHubNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerRaylsAppPrivacyNodeFrozen represents a RaylsApp__PrivacyNodeFrozen error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerRaylsAppPrivacyNodeFrozen struct {
	TokenAddress common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__PrivacyNodeFrozen(address tokenAddress)
func RaylsErc721DvpHandlerRaylsAppPrivacyNodeFrozenErrorID() common.Hash {
	return common.HexToHash("0xcecb8d3ce0d1417038942c9d252e856b5585275082aa5cdbca675fa64d7bfc24")
}

// UnpackRaylsAppPrivacyNodeFrozenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__PrivacyNodeFrozen(address tokenAddress)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackRaylsAppPrivacyNodeFrozenError(raw []byte) (*RaylsErc721DvpHandlerRaylsAppPrivacyNodeFrozen, error) {
	out := new(RaylsErc721DvpHandlerRaylsAppPrivacyNodeFrozen)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "RaylsAppPrivacyNodeFrozen", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerRaylsAppPrivacyNodeNotActive represents a RaylsApp__PrivacyNodeNotActive error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerRaylsAppPrivacyNodeNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__PrivacyNodeNotActive(address tokenAddress, uint8 privacyNodeStatus)
func RaylsErc721DvpHandlerRaylsAppPrivacyNodeNotActiveErrorID() common.Hash {
	return common.HexToHash("0x44c58c43ed8f726e3330349bec7aa7300f000be36837ee0c2cf507d04511e1e8")
}

// UnpackRaylsAppPrivacyNodeNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__PrivacyNodeNotActive(address tokenAddress, uint8 privacyNodeStatus)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackRaylsAppPrivacyNodeNotActiveError(raw []byte) (*RaylsErc721DvpHandlerRaylsAppPrivacyNodeNotActive, error) {
	out := new(RaylsErc721DvpHandlerRaylsAppPrivacyNodeNotActive)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "RaylsAppPrivacyNodeNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerRaylsAppPublicChainNotActive represents a RaylsApp__PublicChainNotActive error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerRaylsAppPublicChainNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
	PublicChainStatus uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__PublicChainNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 publicChainStatus)
func RaylsErc721DvpHandlerRaylsAppPublicChainNotActiveErrorID() common.Hash {
	return common.HexToHash("0xd6e23bd403a5000c9afe5c2ed5202b3ff8e25d8c3644c1f51892016fb18e5ab9")
}

// UnpackRaylsAppPublicChainNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__PublicChainNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 publicChainStatus)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackRaylsAppPublicChainNotActiveError(raw []byte) (*RaylsErc721DvpHandlerRaylsAppPublicChainNotActive, error) {
	out := new(RaylsErc721DvpHandlerRaylsAppPublicChainNotActive)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "RaylsAppPublicChainNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerRaylsAppResourceNotApproved represents a RaylsApp__ResourceNotApproved error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerRaylsAppResourceNotApproved struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__ResourceNotApproved()
func RaylsErc721DvpHandlerRaylsAppResourceNotApprovedErrorID() common.Hash {
	return common.HexToHash("0x970ad4f73c2c200faa068d3d920e2ef40fca6a5338655abcfb5212557edeed6b")
}

// UnpackRaylsAppResourceNotApprovedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__ResourceNotApproved()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackRaylsAppResourceNotApprovedError(raw []byte) (*RaylsErc721DvpHandlerRaylsAppResourceNotApproved, error) {
	out := new(RaylsErc721DvpHandlerRaylsAppResourceNotApproved)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "RaylsAppResourceNotApproved", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerRaylsAppTokenRegistryNotConfigured represents a RaylsApp__TokenRegistryNotConfigured error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerRaylsAppTokenRegistryNotConfigured struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__TokenRegistryNotConfigured()
func RaylsErc721DvpHandlerRaylsAppTokenRegistryNotConfiguredErrorID() common.Hash {
	return common.HexToHash("0x36a41bd1f6f11cd28b716e935a926fb04f66e11a393b38a49bb660640f3b6dbf")
}

// UnpackRaylsAppTokenRegistryNotConfiguredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__TokenRegistryNotConfigured()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackRaylsAppTokenRegistryNotConfiguredError(raw []byte) (*RaylsErc721DvpHandlerRaylsAppTokenRegistryNotConfigured, error) {
	out := new(RaylsErc721DvpHandlerRaylsAppTokenRegistryNotConfigured)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "RaylsAppTokenRegistryNotConfigured", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerRaylsAppUnauthorizedTokenRegistry represents a RaylsApp__UnauthorizedTokenRegistry error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerRaylsAppUnauthorizedTokenRegistry struct {
	Caller        common.Address
	TokenRegistry common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__UnauthorizedTokenRegistry(address caller, address tokenRegistry)
func RaylsErc721DvpHandlerRaylsAppUnauthorizedTokenRegistryErrorID() common.Hash {
	return common.HexToHash("0x061526480acdfaa09331b795496a6c50aaed25a45d9fca4c9d55fad56af8e09c")
}

// UnpackRaylsAppUnauthorizedTokenRegistryError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__UnauthorizedTokenRegistry(address caller, address tokenRegistry)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackRaylsAppUnauthorizedTokenRegistryError(raw []byte) (*RaylsErc721DvpHandlerRaylsAppUnauthorizedTokenRegistry, error) {
	out := new(RaylsErc721DvpHandlerRaylsAppUnauthorizedTokenRegistry)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "RaylsAppUnauthorizedTokenRegistry", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerRaylsAppUserNotRegistered represents a RaylsApp__UserNotRegistered error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerRaylsAppUserNotRegistered struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__UserNotRegistered(address caller)
func RaylsErc721DvpHandlerRaylsAppUserNotRegisteredErrorID() common.Hash {
	return common.HexToHash("0x4c1db902cce08bec31bedc484362fba54949899ac3c0bf0416f3c44af3284baa")
}

// UnpackRaylsAppUserNotRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__UserNotRegistered(address caller)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackRaylsAppUserNotRegisteredError(raw []byte) (*RaylsErc721DvpHandlerRaylsAppUserNotRegistered, error) {
	out := new(RaylsErc721DvpHandlerRaylsAppUserNotRegistered)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "RaylsAppUserNotRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerRaylsErc721DvpHandlerSwapValidityOutOfRange represents a RaylsErc721DvpHandler__SwapValidityOutOfRange error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerRaylsErc721DvpHandlerSwapValidityOutOfRange struct {
	Provided uint64
	Min      uint64
	Max      uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsErc721DvpHandler__SwapValidityOutOfRange(uint64 provided, uint64 min, uint64 max)
func RaylsErc721DvpHandlerRaylsErc721DvpHandlerSwapValidityOutOfRangeErrorID() common.Hash {
	return common.HexToHash("0xb76bd26fcbd207fc30d2a0feb5a6174afca4fe735df51fe6181ef96ac2128b81")
}

// UnpackRaylsErc721DvpHandlerSwapValidityOutOfRangeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsErc721DvpHandler__SwapValidityOutOfRange(uint64 provided, uint64 min, uint64 max)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackRaylsErc721DvpHandlerSwapValidityOutOfRangeError(raw []byte) (*RaylsErc721DvpHandlerRaylsErc721DvpHandlerSwapValidityOutOfRange, error) {
	out := new(RaylsErc721DvpHandlerRaylsErc721DvpHandlerSwapValidityOutOfRange)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "RaylsErc721DvpHandlerSwapValidityOutOfRange", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerRaylsErc721DvpHandlerTokenDoesNotExist represents a RaylsErc721DvpHandler__TokenDoesNotExist error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerRaylsErc721DvpHandlerTokenDoesNotExist struct {
	TokenId *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsErc721DvpHandler__TokenDoesNotExist(uint256 tokenId)
func RaylsErc721DvpHandlerRaylsErc721DvpHandlerTokenDoesNotExistErrorID() common.Hash {
	return common.HexToHash("0x42a96875537d20bad2b70366b679d2ec9733dc3b4d289b0d56afcb6add36653e")
}

// UnpackRaylsErc721DvpHandlerTokenDoesNotExistError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsErc721DvpHandler__TokenDoesNotExist(uint256 tokenId)
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackRaylsErc721DvpHandlerTokenDoesNotExistError(raw []byte) (*RaylsErc721DvpHandlerRaylsErc721DvpHandlerTokenDoesNotExist, error) {
	out := new(RaylsErc721DvpHandlerRaylsErc721DvpHandlerTokenDoesNotExist)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "RaylsErc721DvpHandlerTokenDoesNotExist", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerRaylsErc721DvpHandlerTokenNotApproved represents a RaylsErc721DvpHandler__TokenNotApproved error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerRaylsErc721DvpHandlerTokenNotApproved struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsErc721DvpHandler__TokenNotApproved()
func RaylsErc721DvpHandlerRaylsErc721DvpHandlerTokenNotApprovedErrorID() common.Hash {
	return common.HexToHash("0x569b2b3cdbe13e3ef25dad58137719c7ab2e6ad54385e34adb610feed699941b")
}

// UnpackRaylsErc721DvpHandlerTokenNotApprovedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsErc721DvpHandler__TokenNotApproved()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackRaylsErc721DvpHandlerTokenNotApprovedError(raw []byte) (*RaylsErc721DvpHandlerRaylsErc721DvpHandlerTokenNotApproved, error) {
	out := new(RaylsErc721DvpHandlerRaylsErc721DvpHandlerTokenNotApproved)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "RaylsErc721DvpHandlerTokenNotApproved", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsErc721DvpHandlerReentrancyGuardReentrantCall represents a ReentrancyGuardReentrantCall error raised by the RaylsErc721DvpHandler contract.
type RaylsErc721DvpHandlerReentrancyGuardReentrantCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ReentrancyGuardReentrantCall()
func RaylsErc721DvpHandlerReentrancyGuardReentrantCallErrorID() common.Hash {
	return common.HexToHash("0x3ee5aeb571de7fc460830b4d0017439a1ca56fb0bc39062227ade4fe4a24c1ca")
}

// UnpackReentrancyGuardReentrantCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ReentrancyGuardReentrantCall()
func (raylsErc721DvpHandler *RaylsErc721DvpHandler) UnpackReentrancyGuardReentrantCallError(raw []byte) (*RaylsErc721DvpHandlerReentrancyGuardReentrantCall, error) {
	out := new(RaylsErc721DvpHandlerReentrancyGuardReentrantCall)
	if err := raylsErc721DvpHandler.abi.UnpackIntoInterface(out, "ReentrancyGuardReentrantCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}
