// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package RaylsEnygmaHandler

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

// SharedObjectsEnygmaProgramData is an auto generated low-level Go binding around an user-defined struct.
type SharedObjectsEnygmaProgramData struct {
	ResourceId      [32]byte
	ContractAddress common.Address
	Selector        [4]byte
	Args            []byte
}

// RaylsEnygmaHandlerMetaData contains all meta data concerning the RaylsEnygmaHandler contract.
var RaylsEnygmaHandlerMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"GetERCStandard\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.ErcStandard\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"callWithdrawFromDvp\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelERC1155Swap\",\"inputs\":[{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_toChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_nftId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_nftAmountOrOne\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_nftResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_enygmaAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelERC721Swap\",\"inputs\":[{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_toChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_nftId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_nftResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_enygmaAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"crossBurn\",\"inputs\":[{\"name\":\"_from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"crossMint\",\"inputs\":[{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"crossMintStandard\",\"inputs\":[{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_referenceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"crossRevertMint\",\"inputs\":[{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_reason\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_referenceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"crossTransfer\",\"inputs\":[{\"name\":\"_to\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_value\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"_toChainId\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"_userProgramData\",\"type\":\"tuple[][]\",\"internalType\":\"structSharedObjects.EnygmaProgramData[][]\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"args\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"crossTransferCheck\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"crossTransferFrom\",\"inputs\":[{\"name\":\"_from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_to\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_value\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"_toChainId\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"_userProgramData\",\"type\":\"tuple[][]\",\"internalType\":\"structSharedObjects.EnygmaProgramData[][]\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"args\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"crossTransferRevertBatch\",\"inputs\":[{\"name\":\"_from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_toChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_referenceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"depositToDvp\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvpSwapCompleted\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAddressByResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEndpointAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPNCommunicatorAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"userArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"trusted\",\"type\":\"tuple\",\"internalType\":\"structRaylsTrustedInit\",\"components\":[{\"name\":\"endpoint\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"raylsNodeEndpoint\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"userGovernance\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"linearCrossTransfer\",\"inputs\":[{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_toChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_userProgramData\",\"type\":\"tuple[]\",\"internalType\":\"structSharedObjects.EnygmaProgramData[]\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"args\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"linearCrossTransferFrom\",\"inputs\":[{\"name\":\"_from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_toChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_userProgramData\",\"type\":\"tuple[]\",\"internalType\":\"structSharedObjects.EnygmaProgramData[]\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"args\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"notifySenderAndReceiverWithPNCommunicator\",\"inputs\":[{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_destChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_statusToSender\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.DvpCommunicatiorStatus\"},{\"name\":\"_statusToReceiver\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.DvpCommunicatiorStatus\"},{\"name\":\"_messageToSender\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_messageToReceiver\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"notifySenderWithPNCommunicator\",\"inputs\":[{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_status\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.DvpCommunicatiorStatus\"},{\"name\":\"_message\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"raylsNodeUserGovernance\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIUserGovernance\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"receiveWithdrawFromDvp\",\"inputs\":[{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_referenceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"referenceIdStatus\",\"inputs\":[{\"name\":\"_referenceID\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumRaylsEnygmaHandler.ReferenceIdStatus\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"referenceIdStatusUint\",\"inputs\":[{\"name\":\"_referenceID\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"resourceId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSwapValidityTime\",\"inputs\":[{\"name\":\"_validityTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supplyUpdateRevert\",\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_isMint\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"swapValidityTime\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"swapWithDvpForERC1155\",\"inputs\":[{\"name\":\"_nftId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_nftAmountOrOne\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_nftResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_enygmaAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_destChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_validityTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"swapWithDvpForERC721\",\"inputs\":[{\"name\":\"_nftId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_nftResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_enygmaAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_destChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_validityTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RaylsEnygmaErc20TokenCreated\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenRegistrationSubmitted\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"crossTransferReferenceId\",\"inputs\":[{\"name\":\"_referenceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"transactionReferenceId\",\"inputs\":[{\"name\":\"_referenceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__HubNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"hubStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__PrivacyNodeFrozen\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__PrivacyNodeNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__PublicChainNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"publicChainStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__ResourceNotApproved\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsApp__TokenRegistryNotConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsApp__UnauthorizedTokenRegistry\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenRegistry\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__UserNotRegistered\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsEnygmaHandler__ArrayLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsEnygmaHandler__AuthorityNotSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsEnygmaHandler__CallableExecutionFailed\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsEnygmaHandler__CallableTargetNotContract\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsEnygmaHandler__EmptyArray\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsEnygmaHandler__InvalidDecimals\",\"inputs\":[{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsEnygmaHandler__NotRelayer\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsEnygmaHandler__NotTokenOwnerScoped\",\"inputs\":[{\"name\":\"originSender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsEnygmaHandler__SwapValidityOutOfRange\",\"inputs\":[{\"name\":\"provided\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"min\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"max\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"RaylsEnygmaHandler__TooManyUniqueChainIds\",\"inputs\":[{\"name\":\"count\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"RaylsEnygmaHandler__WrongAddress\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsEnygmaHandler__WrongFunctionForSameChainId\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"RaylsEnygmaHandler__ZeroAddress\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsEnygmaHandler__ZeroAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsEnygmaHandler__ZeroValueArg\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
	ID:  "RaylsEnygmaHandler",
}

// RaylsEnygmaHandler is an auto generated Go binding around an Ethereum contract.
type RaylsEnygmaHandler struct {
	abi abi.ABI
}

// NewRaylsEnygmaHandler creates a new instance of RaylsEnygmaHandler.
func NewRaylsEnygmaHandler() *RaylsEnygmaHandler {
	parsed, err := RaylsEnygmaHandlerMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &RaylsEnygmaHandler{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *RaylsEnygmaHandler) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackGetERCStandard is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x35abee1a.
//
// Solidity: function GetERCStandard() pure returns(uint8)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackGetERCStandard() []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("GetERCStandard")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetERCStandard is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x35abee1a.
//
// Solidity: function GetERCStandard() pure returns(uint8)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackGetERCStandard(data []byte) (uint8, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("GetERCStandard", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, err
}

// PackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackAllowance(owner common.Address, spender common.Address) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("allowance", owner, spender)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAllowance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackAllowance(data []byte) (*big.Int, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("allowance", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackApprove(spender common.Address, value *big.Int) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("approve", spender, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackApprove is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackApprove(data []byte) (bool, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("approve", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackAuthority() []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("authority", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackBalanceOf(account common.Address) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("balanceOf", account)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackBurn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9dc29fac.
//
// Solidity: function burn(address from, uint256 value) returns()
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackBurn(from common.Address, value *big.Int) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("burn", from, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCallWithdrawFromDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3779083a.
//
// Solidity: function callWithdrawFromDvp(uint256 amount) returns(bytes32)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackCallWithdrawFromDvp(amount *big.Int) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("callWithdrawFromDvp", amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackCallWithdrawFromDvp is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3779083a.
//
// Solidity: function callWithdrawFromDvp(uint256 amount) returns(bytes32)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackCallWithdrawFromDvp(data []byte) ([32]byte, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("callWithdrawFromDvp", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackCancelERC1155Swap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x729a1f41.
//
// Solidity: function cancelERC1155Swap(bytes32 _sharedId, uint256 _toChainId, uint256 _nftId, uint256 _nftAmountOrOne, bytes32 _nftResourceId, uint256 _enygmaAmount) returns()
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackCancelERC1155Swap(sharedId [32]byte, toChainId *big.Int, nftId *big.Int, nftAmountOrOne *big.Int, nftResourceId [32]byte, enygmaAmount *big.Int) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("cancelERC1155Swap", sharedId, toChainId, nftId, nftAmountOrOne, nftResourceId, enygmaAmount)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCancelERC721Swap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc397ef31.
//
// Solidity: function cancelERC721Swap(bytes32 _sharedId, uint256 _toChainId, uint256 _nftId, bytes32 _nftResourceId, uint256 _enygmaAmount) returns()
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackCancelERC721Swap(sharedId [32]byte, toChainId *big.Int, nftId *big.Int, nftResourceId [32]byte, enygmaAmount *big.Int) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("cancelERC721Swap", sharedId, toChainId, nftId, nftResourceId, enygmaAmount)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCrossBurn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x589a9e6e.
//
// Solidity: function crossBurn(address _from, uint256 _value) returns()
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackCrossBurn(from common.Address, value *big.Int) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("crossBurn", from, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCrossMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x306c68a6.
//
// Solidity: function crossMint(address _to, uint256 _value) returns()
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackCrossMint(to common.Address, value *big.Int) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("crossMint", to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCrossMintStandard is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x089bfdb0.
//
// Solidity: function crossMintStandard(address _to, uint256 _value, bytes32 _referenceId) returns()
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackCrossMintStandard(to common.Address, value *big.Int, referenceId [32]byte) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("crossMintStandard", to, value, referenceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCrossRevertMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2c69886f.
//
// Solidity: function crossRevertMint(address _to, uint256 _value, string _reason, bytes32 _referenceId) returns()
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackCrossRevertMint(to common.Address, value *big.Int, reason string, referenceId [32]byte) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("crossRevertMint", to, value, reason, referenceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCrossTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8a3c70f0.
//
// Solidity: function crossTransfer(address[] _to, uint256[] _value, uint256[] _toChainId, (bytes32,address,bytes4,bytes)[][] _userProgramData) returns(bytes32)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackCrossTransfer(to []common.Address, value []*big.Int, toChainId []*big.Int, userProgramData [][]SharedObjectsEnygmaProgramData) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("crossTransfer", to, value, toChainId, userProgramData)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackCrossTransfer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8a3c70f0.
//
// Solidity: function crossTransfer(address[] _to, uint256[] _value, uint256[] _toChainId, (bytes32,address,bytes4,bytes)[][] _userProgramData) returns(bytes32)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackCrossTransfer(data []byte) ([32]byte, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("crossTransfer", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackCrossTransferCheck is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x195efa4e.
//
// Solidity: function crossTransferCheck() returns()
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackCrossTransferCheck() []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("crossTransferCheck")
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCrossTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52c0e498.
//
// Solidity: function crossTransferFrom(address _from, address[] _to, uint256[] _value, uint256[] _toChainId, (bytes32,address,bytes4,bytes)[][] _userProgramData) returns(bytes32)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackCrossTransferFrom(from common.Address, to []common.Address, value []*big.Int, toChainId []*big.Int, userProgramData [][]SharedObjectsEnygmaProgramData) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("crossTransferFrom", from, to, value, toChainId, userProgramData)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackCrossTransferFrom is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52c0e498.
//
// Solidity: function crossTransferFrom(address _from, address[] _to, uint256[] _value, uint256[] _toChainId, (bytes32,address,bytes4,bytes)[][] _userProgramData) returns(bytes32)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackCrossTransferFrom(data []byte) ([32]byte, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("crossTransferFrom", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackCrossTransferRevertBatch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe11a96f4.
//
// Solidity: function crossTransferRevertBatch(address _from, address _to, uint256 _value, uint256 _toChainId, bytes32 _referenceId) returns()
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackCrossTransferRevertBatch(from common.Address, to common.Address, value *big.Int, toChainId *big.Int, referenceId [32]byte) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("crossTransferRevertBatch", from, to, value, toChainId, referenceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackDecimals() []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("decimals")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDecimals is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackDecimals(data []byte) (uint8, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("decimals", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, err
}

// PackDepositToDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb778f933.
//
// Solidity: function depositToDvp(uint256 amount) returns(bytes32)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackDepositToDvp(amount *big.Int) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("depositToDvp", amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDepositToDvp is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb778f933.
//
// Solidity: function depositToDvp(uint256 amount) returns(bytes32)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackDepositToDvp(data []byte) ([32]byte, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("depositToDvp", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackDvpSwapCompleted is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x14fe9eb2.
//
// Solidity: function dvpSwapCompleted(uint256 , bytes32 _sharedId) returns()
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackDvpSwapCompleted(arg0 *big.Int, sharedId [32]byte) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("dvpSwapCompleted", arg0, sharedId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackGetAddressByResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackGetAddressByResourceId(resourceId [32]byte) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("getAddressByResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAddressByResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackGetAddressByResourceId(data []byte) (common.Address, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("getAddressByResourceId", data)
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
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackGetEndpointAddress() []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("getEndpointAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetEndpointAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xce884eb5.
//
// Solidity: function getEndpointAddress() view returns(address)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackGetEndpointAddress(data []byte) (common.Address, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("getEndpointAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetPNCommunicatorAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb417c79b.
//
// Solidity: function getPNCommunicatorAddress() view returns(address)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackGetPNCommunicatorAddress() []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("getPNCommunicatorAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetPNCommunicatorAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb417c79b.
//
// Solidity: function getPNCommunicatorAddress() view returns(address)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackGetPNCommunicatorAddress(data []byte) (common.Address, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("getPNCommunicatorAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1d3f35b0.
//
// Solidity: function initialize(bytes userArgs, (address,address,address,address,bytes32,address) trusted) returns()
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackInitialize(userArgs []byte, trusted RaylsTrustedInit) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("initialize", userArgs, trusted)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackLinearCrossTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8be9e327.
//
// Solidity: function linearCrossTransfer(address _to, uint256 _value, uint256 _toChainId, (bytes32,address,bytes4,bytes)[] _userProgramData) returns(bytes32)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackLinearCrossTransfer(to common.Address, value *big.Int, toChainId *big.Int, userProgramData []SharedObjectsEnygmaProgramData) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("linearCrossTransfer", to, value, toChainId, userProgramData)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackLinearCrossTransfer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8be9e327.
//
// Solidity: function linearCrossTransfer(address _to, uint256 _value, uint256 _toChainId, (bytes32,address,bytes4,bytes)[] _userProgramData) returns(bytes32)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackLinearCrossTransfer(data []byte) ([32]byte, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("linearCrossTransfer", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackLinearCrossTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3ab53206.
//
// Solidity: function linearCrossTransferFrom(address _from, address _to, uint256 _value, uint256 _toChainId, (bytes32,address,bytes4,bytes)[] _userProgramData) returns(bytes32)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackLinearCrossTransferFrom(from common.Address, to common.Address, value *big.Int, toChainId *big.Int, userProgramData []SharedObjectsEnygmaProgramData) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("linearCrossTransferFrom", from, to, value, toChainId, userProgramData)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackLinearCrossTransferFrom is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3ab53206.
//
// Solidity: function linearCrossTransferFrom(address _from, address _to, uint256 _value, uint256 _toChainId, (bytes32,address,bytes4,bytes)[] _userProgramData) returns(bytes32)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackLinearCrossTransferFrom(data []byte) ([32]byte, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("linearCrossTransferFrom", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x40c10f19.
//
// Solidity: function mint(address _to, uint256 _value) returns()
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackMint(to common.Address, value *big.Int) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("mint", to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackName() []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("name")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackName(data []byte) (string, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("name", data)
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
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackNotifySenderAndReceiverWithPNCommunicator(sharedId [32]byte, destChainId *big.Int, statusToSender uint8, statusToReceiver uint8, messageToSender string, messageToReceiver string) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("notifySenderAndReceiverWithPNCommunicator", sharedId, destChainId, statusToSender, statusToReceiver, messageToSender, messageToReceiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackNotifySenderWithPNCommunicator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcc90e43b.
//
// Solidity: function notifySenderWithPNCommunicator(bytes32 _sharedId, uint8 _status, string _message) returns()
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackNotifySenderWithPNCommunicator(sharedId [32]byte, status uint8, message string) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("notifySenderWithPNCommunicator", sharedId, status, message)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRaylsNodeUserGovernance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0274c133.
//
// Solidity: function raylsNodeUserGovernance() view returns(address)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackRaylsNodeUserGovernance() []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("raylsNodeUserGovernance")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackRaylsNodeUserGovernance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0274c133.
//
// Solidity: function raylsNodeUserGovernance() view returns(address)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsNodeUserGovernance(data []byte) (common.Address, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("raylsNodeUserGovernance", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackReceiveWithdrawFromDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe4770363.
//
// Solidity: function receiveWithdrawFromDvp(address _to, uint256 _value, bytes32 _referenceId) returns()
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackReceiveWithdrawFromDvp(to common.Address, value *big.Int, referenceId [32]byte) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("receiveWithdrawFromDvp", to, value, referenceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackReferenceIdStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8d742261.
//
// Solidity: function referenceIdStatus(bytes32 _referenceID) view returns(uint8)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackReferenceIdStatus(referenceID [32]byte) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("referenceIdStatus", referenceID)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackReferenceIdStatus is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8d742261.
//
// Solidity: function referenceIdStatus(bytes32 _referenceID) view returns(uint8)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackReferenceIdStatus(data []byte) (uint8, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("referenceIdStatus", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, err
}

// PackReferenceIdStatusUint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd97a5df1.
//
// Solidity: function referenceIdStatusUint(bytes32 _referenceID) view returns(uint256)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackReferenceIdStatusUint(referenceID [32]byte) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("referenceIdStatusUint", referenceID)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackReferenceIdStatusUint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd97a5df1.
//
// Solidity: function referenceIdStatusUint(bytes32 _referenceID) view returns(uint256)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackReferenceIdStatusUint(data []byte) (*big.Int, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("referenceIdStatusUint", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackResourceId() []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("resourceId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackResourceId(data []byte) ([32]byte, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("resourceId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackSetResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa01afbfb.
//
// Solidity: function setResourceId(bytes32 _resourceId) returns()
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackSetResourceId(resourceId [32]byte) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("setResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetSwapValidityTime is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5ddbe8f0.
//
// Solidity: function setSwapValidityTime(uint64 _validityTime) returns()
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackSetSwapValidityTime(validityTime uint64) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("setSwapValidityTime", validityTime)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSupplyUpdateRevert is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x86331214.
//
// Solidity: function supplyUpdateRevert(uint256 _amount, address _recipient, bool _isMint) returns()
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackSupplyUpdateRevert(amount *big.Int, recipient common.Address, isMint bool) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("supplyUpdateRevert", amount, recipient, isMint)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSwapValidityTime is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x521d9498.
//
// Solidity: function swapValidityTime() view returns(uint64)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackSwapValidityTime() []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("swapValidityTime")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSwapValidityTime is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x521d9498.
//
// Solidity: function swapValidityTime() view returns(uint64)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackSwapValidityTime(data []byte) (uint64, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("swapValidityTime", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, err
}

// PackSwapWithDvpForERC1155 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18db13aa.
//
// Solidity: function swapWithDvpForERC1155(uint256 _nftId, uint256 _nftAmountOrOne, bytes32 _nftResourceId, uint256 _enygmaAmount, uint256 _destChainId, bytes32 _sharedId, uint64 _validityTime) returns()
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackSwapWithDvpForERC1155(nftId *big.Int, nftAmountOrOne *big.Int, nftResourceId [32]byte, enygmaAmount *big.Int, destChainId *big.Int, sharedId [32]byte, validityTime uint64) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("swapWithDvpForERC1155", nftId, nftAmountOrOne, nftResourceId, enygmaAmount, destChainId, sharedId, validityTime)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSwapWithDvpForERC721 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdfc507a9.
//
// Solidity: function swapWithDvpForERC721(uint256 _nftId, bytes32 _nftResourceId, uint256 _enygmaAmount, uint256 _destChainId, bytes32 _sharedId, uint64 _validityTime) returns()
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackSwapWithDvpForERC721(nftId *big.Int, nftResourceId [32]byte, enygmaAmount *big.Int, destChainId *big.Int, sharedId [32]byte, validityTime uint64) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("swapWithDvpForERC721", nftId, nftResourceId, enygmaAmount, destChainId, sharedId, validityTime)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackSymbol() []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("symbol")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackSymbol(data []byte) (string, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("symbol", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackTotalSupply() []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("totalSupply")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackTotalSupply(data []byte) (*big.Int, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("totalSupply", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackTransfer(to common.Address, value *big.Int) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("transfer", to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTransfer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackTransfer(data []byte) (bool, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("transfer", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (raylsEnygmaHandler *RaylsEnygmaHandler) PackTransferFrom(from common.Address, to common.Address, value *big.Int) []byte {
	enc, err := raylsEnygmaHandler.abi.Pack("transferFrom", from, to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTransferFrom is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackTransferFrom(data []byte) (bool, error) {
	out, err := raylsEnygmaHandler.abi.Unpack("transferFrom", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// RaylsEnygmaHandlerApproval represents a Approval event raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const RaylsEnygmaHandlerApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (RaylsEnygmaHandlerApproval) ContractEventName() string {
	return RaylsEnygmaHandlerApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackApprovalEvent(log *types.Log) (*RaylsEnygmaHandlerApproval, error) {
	event := "Approval"
	if log.Topics[0] != raylsEnygmaHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsEnygmaHandlerApproval)
	if len(log.Data) > 0 {
		if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsEnygmaHandler.abi.Events[event].Inputs {
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

// RaylsEnygmaHandlerAuthorityUpdated represents a AuthorityUpdated event raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerAuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const RaylsEnygmaHandlerAuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (RaylsEnygmaHandlerAuthorityUpdated) ContractEventName() string {
	return RaylsEnygmaHandlerAuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackAuthorityUpdatedEvent(log *types.Log) (*RaylsEnygmaHandlerAuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != raylsEnygmaHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsEnygmaHandlerAuthorityUpdated)
	if len(log.Data) > 0 {
		if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsEnygmaHandler.abi.Events[event].Inputs {
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

// RaylsEnygmaHandlerInitialized represents a Initialized event raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerInitialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const RaylsEnygmaHandlerInitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (RaylsEnygmaHandlerInitialized) ContractEventName() string {
	return RaylsEnygmaHandlerInitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackInitializedEvent(log *types.Log) (*RaylsEnygmaHandlerInitialized, error) {
	event := "Initialized"
	if log.Topics[0] != raylsEnygmaHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsEnygmaHandlerInitialized)
	if len(log.Data) > 0 {
		if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsEnygmaHandler.abi.Events[event].Inputs {
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

// RaylsEnygmaHandlerRaylsEnygmaErc20TokenCreated represents a RaylsEnygmaErc20TokenCreated event raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsEnygmaErc20TokenCreated struct {
	TokenAddress common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const RaylsEnygmaHandlerRaylsEnygmaErc20TokenCreatedEventName = "RaylsEnygmaErc20TokenCreated"

// ContractEventName returns the user-defined event name.
func (RaylsEnygmaHandlerRaylsEnygmaErc20TokenCreated) ContractEventName() string {
	return RaylsEnygmaHandlerRaylsEnygmaErc20TokenCreatedEventName
}

// UnpackRaylsEnygmaErc20TokenCreatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RaylsEnygmaErc20TokenCreated(address indexed tokenAddress)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsEnygmaErc20TokenCreatedEvent(log *types.Log) (*RaylsEnygmaHandlerRaylsEnygmaErc20TokenCreated, error) {
	event := "RaylsEnygmaErc20TokenCreated"
	if log.Topics[0] != raylsEnygmaHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsEnygmaHandlerRaylsEnygmaErc20TokenCreated)
	if len(log.Data) > 0 {
		if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsEnygmaHandler.abi.Events[event].Inputs {
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

// RaylsEnygmaHandlerTokenRegistrationSubmitted represents a TokenRegistrationSubmitted event raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerTokenRegistrationSubmitted struct {
	TokenAddress common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const RaylsEnygmaHandlerTokenRegistrationSubmittedEventName = "TokenRegistrationSubmitted"

// ContractEventName returns the user-defined event name.
func (RaylsEnygmaHandlerTokenRegistrationSubmitted) ContractEventName() string {
	return RaylsEnygmaHandlerTokenRegistrationSubmittedEventName
}

// UnpackTokenRegistrationSubmittedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenRegistrationSubmitted(address indexed tokenAddress)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackTokenRegistrationSubmittedEvent(log *types.Log) (*RaylsEnygmaHandlerTokenRegistrationSubmitted, error) {
	event := "TokenRegistrationSubmitted"
	if log.Topics[0] != raylsEnygmaHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsEnygmaHandlerTokenRegistrationSubmitted)
	if len(log.Data) > 0 {
		if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsEnygmaHandler.abi.Events[event].Inputs {
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

// RaylsEnygmaHandlerTransfer represents a Transfer event raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   *types.Log // Blockchain specific contextual infos
}

const RaylsEnygmaHandlerTransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (RaylsEnygmaHandlerTransfer) ContractEventName() string {
	return RaylsEnygmaHandlerTransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackTransferEvent(log *types.Log) (*RaylsEnygmaHandlerTransfer, error) {
	event := "Transfer"
	if log.Topics[0] != raylsEnygmaHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsEnygmaHandlerTransfer)
	if len(log.Data) > 0 {
		if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsEnygmaHandler.abi.Events[event].Inputs {
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

// RaylsEnygmaHandlerCrossTransferReferenceId represents a crossTransferReferenceId event raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerCrossTransferReferenceId struct {
	ReferenceId [32]byte
	Raw         *types.Log // Blockchain specific contextual infos
}

const RaylsEnygmaHandlerCrossTransferReferenceIdEventName = "crossTransferReferenceId"

// ContractEventName returns the user-defined event name.
func (RaylsEnygmaHandlerCrossTransferReferenceId) ContractEventName() string {
	return RaylsEnygmaHandlerCrossTransferReferenceIdEventName
}

// UnpackCrossTransferReferenceIdEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event crossTransferReferenceId(bytes32 _referenceId)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackCrossTransferReferenceIdEvent(log *types.Log) (*RaylsEnygmaHandlerCrossTransferReferenceId, error) {
	event := "crossTransferReferenceId"
	if log.Topics[0] != raylsEnygmaHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsEnygmaHandlerCrossTransferReferenceId)
	if len(log.Data) > 0 {
		if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsEnygmaHandler.abi.Events[event].Inputs {
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

// RaylsEnygmaHandlerTransactionReferenceId represents a transactionReferenceId event raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerTransactionReferenceId struct {
	ReferenceId [32]byte
	Raw         *types.Log // Blockchain specific contextual infos
}

const RaylsEnygmaHandlerTransactionReferenceIdEventName = "transactionReferenceId"

// ContractEventName returns the user-defined event name.
func (RaylsEnygmaHandlerTransactionReferenceId) ContractEventName() string {
	return RaylsEnygmaHandlerTransactionReferenceIdEventName
}

// UnpackTransactionReferenceIdEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event transactionReferenceId(bytes32 _referenceId)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackTransactionReferenceIdEvent(log *types.Log) (*RaylsEnygmaHandlerTransactionReferenceId, error) {
	event := "transactionReferenceId"
	if log.Topics[0] != raylsEnygmaHandler.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RaylsEnygmaHandlerTransactionReferenceId)
	if len(log.Data) > 0 {
		if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range raylsEnygmaHandler.abi.Events[event].Inputs {
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
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["ERC20InsufficientAllowance"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackERC20InsufficientAllowanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["ERC20InsufficientBalance"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackERC20InsufficientBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["ERC20InvalidApprover"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackERC20InvalidApproverError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["ERC20InvalidReceiver"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackERC20InvalidReceiverError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["ERC20InvalidSender"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackERC20InvalidSenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["ERC20InvalidSpender"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackERC20InvalidSpenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsAppHubNotActive"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsAppHubNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsAppPrivacyNodeFrozen"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsAppPrivacyNodeFrozenError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsAppPrivacyNodeNotActive"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsAppPrivacyNodeNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsAppPublicChainNotActive"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsAppPublicChainNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsAppResourceNotApproved"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsAppResourceNotApprovedError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsAppTokenRegistryNotConfigured"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsAppTokenRegistryNotConfiguredError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsAppUnauthorizedTokenRegistry"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsAppUnauthorizedTokenRegistryError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsAppUserNotRegistered"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsAppUserNotRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsEnygmaHandlerArrayLengthMismatch"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsEnygmaHandlerArrayLengthMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsEnygmaHandlerAuthorityNotSet"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsEnygmaHandlerAuthorityNotSetError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsEnygmaHandlerCallableExecutionFailed"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsEnygmaHandlerCallableExecutionFailedError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsEnygmaHandlerCallableTargetNotContract"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsEnygmaHandlerCallableTargetNotContractError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsEnygmaHandlerEmptyArray"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsEnygmaHandlerEmptyArrayError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsEnygmaHandlerInvalidDecimals"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsEnygmaHandlerInvalidDecimalsError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsEnygmaHandlerNotRelayer"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsEnygmaHandlerNotRelayerError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsEnygmaHandlerNotTokenOwnerScoped"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsEnygmaHandlerNotTokenOwnerScopedError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsEnygmaHandlerSwapValidityOutOfRange"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsEnygmaHandlerSwapValidityOutOfRangeError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsEnygmaHandlerTooManyUniqueChainIds"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsEnygmaHandlerTooManyUniqueChainIdsError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsEnygmaHandlerWrongAddress"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsEnygmaHandlerWrongAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsEnygmaHandlerWrongFunctionForSameChainId"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsEnygmaHandlerWrongFunctionForSameChainIdError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsEnygmaHandlerZeroAddress"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsEnygmaHandlerZeroAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsEnygmaHandlerZeroAmount"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsEnygmaHandlerZeroAmountError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["RaylsEnygmaHandlerZeroValueArg"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackRaylsEnygmaHandlerZeroValueArgError(raw[4:])
	}
	if bytes.Equal(raw[:4], raylsEnygmaHandler.abi.Errors["ReentrancyGuardReentrantCall"].ID.Bytes()[:4]) {
		return raylsEnygmaHandler.UnpackReentrancyGuardReentrantCallError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// RaylsEnygmaHandlerERC20InsufficientAllowance represents a ERC20InsufficientAllowance error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerERC20InsufficientAllowance struct {
	Spender   common.Address
	Allowance *big.Int
	Needed    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func RaylsEnygmaHandlerERC20InsufficientAllowanceErrorID() common.Hash {
	return common.HexToHash("0xfb8f41b23e99d2101d86da76cdfa87dd51c82ed07d3cb62cbc473e469dbc75c3")
}

// UnpackERC20InsufficientAllowanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackERC20InsufficientAllowanceError(raw []byte) (*RaylsEnygmaHandlerERC20InsufficientAllowance, error) {
	out := new(RaylsEnygmaHandlerERC20InsufficientAllowance)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "ERC20InsufficientAllowance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerERC20InsufficientBalance represents a ERC20InsufficientBalance error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerERC20InsufficientBalance struct {
	Sender  common.Address
	Balance *big.Int
	Needed  *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func RaylsEnygmaHandlerERC20InsufficientBalanceErrorID() common.Hash {
	return common.HexToHash("0xe450d38cd8d9f7d95077d567d60ed49c7254716e6ad08fc9872816c97e0ffec6")
}

// UnpackERC20InsufficientBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackERC20InsufficientBalanceError(raw []byte) (*RaylsEnygmaHandlerERC20InsufficientBalance, error) {
	out := new(RaylsEnygmaHandlerERC20InsufficientBalance)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "ERC20InsufficientBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerERC20InvalidApprover represents a ERC20InvalidApprover error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerERC20InvalidApprover struct {
	Approver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidApprover(address approver)
func RaylsEnygmaHandlerERC20InvalidApproverErrorID() common.Hash {
	return common.HexToHash("0xe602df05cc75712490294c6c104ab7c17f4030363910a7a2626411c6d3118847")
}

// UnpackERC20InvalidApproverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidApprover(address approver)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackERC20InvalidApproverError(raw []byte) (*RaylsEnygmaHandlerERC20InvalidApprover, error) {
	out := new(RaylsEnygmaHandlerERC20InvalidApprover)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "ERC20InvalidApprover", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerERC20InvalidReceiver represents a ERC20InvalidReceiver error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerERC20InvalidReceiver struct {
	Receiver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func RaylsEnygmaHandlerERC20InvalidReceiverErrorID() common.Hash {
	return common.HexToHash("0xec442f055133b72f3b2f9f0bb351c406b178527de2040a7d1feb4e058771f613")
}

// UnpackERC20InvalidReceiverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackERC20InvalidReceiverError(raw []byte) (*RaylsEnygmaHandlerERC20InvalidReceiver, error) {
	out := new(RaylsEnygmaHandlerERC20InvalidReceiver)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "ERC20InvalidReceiver", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerERC20InvalidSender represents a ERC20InvalidSender error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerERC20InvalidSender struct {
	Sender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSender(address sender)
func RaylsEnygmaHandlerERC20InvalidSenderErrorID() common.Hash {
	return common.HexToHash("0x96c6fd1edd0cd6ef7ff0ecc0facdf53148dc0048b57fe58af65755250a7a96bd")
}

// UnpackERC20InvalidSenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSender(address sender)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackERC20InvalidSenderError(raw []byte) (*RaylsEnygmaHandlerERC20InvalidSender, error) {
	out := new(RaylsEnygmaHandlerERC20InvalidSender)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "ERC20InvalidSender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerERC20InvalidSpender represents a ERC20InvalidSpender error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerERC20InvalidSpender struct {
	Spender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSpender(address spender)
func RaylsEnygmaHandlerERC20InvalidSpenderErrorID() common.Hash {
	return common.HexToHash("0x94280d62c347d8d9f4d59a76ea321452406db88df38e0c9da304f58b57b373a2")
}

// UnpackERC20InvalidSpenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSpender(address spender)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackERC20InvalidSpenderError(raw []byte) (*RaylsEnygmaHandlerERC20InvalidSpender, error) {
	out := new(RaylsEnygmaHandlerERC20InvalidSpender)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "ERC20InvalidSpender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerInvalidInitialization represents a InvalidInitialization error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerInvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func RaylsEnygmaHandlerInvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackInvalidInitializationError(raw []byte) (*RaylsEnygmaHandlerInvalidInitialization, error) {
	out := new(RaylsEnygmaHandlerInvalidInitialization)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerNotInitializing represents a NotInitializing error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerNotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func RaylsEnygmaHandlerNotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackNotInitializingError(raw []byte) (*RaylsEnygmaHandlerNotInitializing, error) {
	out := new(RaylsEnygmaHandlerNotInitializing)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func RaylsEnygmaHandlerRaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*RaylsEnygmaHandlerRaylsAccessManagedContractPaused, error) {
	out := new(RaylsEnygmaHandlerRaylsAccessManagedContractPaused)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func RaylsEnygmaHandlerRaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*RaylsEnygmaHandlerRaylsAccessManagedInvalidAuthority, error) {
	out := new(RaylsEnygmaHandlerRaylsAccessManagedInvalidAuthority)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func RaylsEnygmaHandlerRaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*RaylsEnygmaHandlerRaylsAccessManagedMustSchedule, error) {
	out := new(RaylsEnygmaHandlerRaylsAccessManagedMustSchedule)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func RaylsEnygmaHandlerRaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*RaylsEnygmaHandlerRaylsAccessManagedUnauthorized, error) {
	out := new(RaylsEnygmaHandlerRaylsAccessManagedUnauthorized)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsAppHubNotActive represents a RaylsApp__HubNotActive error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsAppHubNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
	HubStatus         uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__HubNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 hubStatus)
func RaylsEnygmaHandlerRaylsAppHubNotActiveErrorID() common.Hash {
	return common.HexToHash("0xdc2ffb0fada912f0dd1b700d4ea9a9ce47e3ecdd1b7b155d2066b9a022a637c2")
}

// UnpackRaylsAppHubNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__HubNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 hubStatus)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsAppHubNotActiveError(raw []byte) (*RaylsEnygmaHandlerRaylsAppHubNotActive, error) {
	out := new(RaylsEnygmaHandlerRaylsAppHubNotActive)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsAppHubNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsAppPrivacyNodeFrozen represents a RaylsApp__PrivacyNodeFrozen error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsAppPrivacyNodeFrozen struct {
	TokenAddress common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__PrivacyNodeFrozen(address tokenAddress)
func RaylsEnygmaHandlerRaylsAppPrivacyNodeFrozenErrorID() common.Hash {
	return common.HexToHash("0xcecb8d3ce0d1417038942c9d252e856b5585275082aa5cdbca675fa64d7bfc24")
}

// UnpackRaylsAppPrivacyNodeFrozenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__PrivacyNodeFrozen(address tokenAddress)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsAppPrivacyNodeFrozenError(raw []byte) (*RaylsEnygmaHandlerRaylsAppPrivacyNodeFrozen, error) {
	out := new(RaylsEnygmaHandlerRaylsAppPrivacyNodeFrozen)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsAppPrivacyNodeFrozen", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsAppPrivacyNodeNotActive represents a RaylsApp__PrivacyNodeNotActive error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsAppPrivacyNodeNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__PrivacyNodeNotActive(address tokenAddress, uint8 privacyNodeStatus)
func RaylsEnygmaHandlerRaylsAppPrivacyNodeNotActiveErrorID() common.Hash {
	return common.HexToHash("0x44c58c43ed8f726e3330349bec7aa7300f000be36837ee0c2cf507d04511e1e8")
}

// UnpackRaylsAppPrivacyNodeNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__PrivacyNodeNotActive(address tokenAddress, uint8 privacyNodeStatus)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsAppPrivacyNodeNotActiveError(raw []byte) (*RaylsEnygmaHandlerRaylsAppPrivacyNodeNotActive, error) {
	out := new(RaylsEnygmaHandlerRaylsAppPrivacyNodeNotActive)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsAppPrivacyNodeNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsAppPublicChainNotActive represents a RaylsApp__PublicChainNotActive error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsAppPublicChainNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
	PublicChainStatus uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__PublicChainNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 publicChainStatus)
func RaylsEnygmaHandlerRaylsAppPublicChainNotActiveErrorID() common.Hash {
	return common.HexToHash("0xd6e23bd403a5000c9afe5c2ed5202b3ff8e25d8c3644c1f51892016fb18e5ab9")
}

// UnpackRaylsAppPublicChainNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__PublicChainNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 publicChainStatus)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsAppPublicChainNotActiveError(raw []byte) (*RaylsEnygmaHandlerRaylsAppPublicChainNotActive, error) {
	out := new(RaylsEnygmaHandlerRaylsAppPublicChainNotActive)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsAppPublicChainNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsAppResourceNotApproved represents a RaylsApp__ResourceNotApproved error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsAppResourceNotApproved struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__ResourceNotApproved()
func RaylsEnygmaHandlerRaylsAppResourceNotApprovedErrorID() common.Hash {
	return common.HexToHash("0x970ad4f73c2c200faa068d3d920e2ef40fca6a5338655abcfb5212557edeed6b")
}

// UnpackRaylsAppResourceNotApprovedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__ResourceNotApproved()
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsAppResourceNotApprovedError(raw []byte) (*RaylsEnygmaHandlerRaylsAppResourceNotApproved, error) {
	out := new(RaylsEnygmaHandlerRaylsAppResourceNotApproved)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsAppResourceNotApproved", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsAppTokenRegistryNotConfigured represents a RaylsApp__TokenRegistryNotConfigured error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsAppTokenRegistryNotConfigured struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__TokenRegistryNotConfigured()
func RaylsEnygmaHandlerRaylsAppTokenRegistryNotConfiguredErrorID() common.Hash {
	return common.HexToHash("0x36a41bd1f6f11cd28b716e935a926fb04f66e11a393b38a49bb660640f3b6dbf")
}

// UnpackRaylsAppTokenRegistryNotConfiguredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__TokenRegistryNotConfigured()
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsAppTokenRegistryNotConfiguredError(raw []byte) (*RaylsEnygmaHandlerRaylsAppTokenRegistryNotConfigured, error) {
	out := new(RaylsEnygmaHandlerRaylsAppTokenRegistryNotConfigured)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsAppTokenRegistryNotConfigured", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsAppUnauthorizedTokenRegistry represents a RaylsApp__UnauthorizedTokenRegistry error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsAppUnauthorizedTokenRegistry struct {
	Caller        common.Address
	TokenRegistry common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__UnauthorizedTokenRegistry(address caller, address tokenRegistry)
func RaylsEnygmaHandlerRaylsAppUnauthorizedTokenRegistryErrorID() common.Hash {
	return common.HexToHash("0x061526480acdfaa09331b795496a6c50aaed25a45d9fca4c9d55fad56af8e09c")
}

// UnpackRaylsAppUnauthorizedTokenRegistryError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__UnauthorizedTokenRegistry(address caller, address tokenRegistry)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsAppUnauthorizedTokenRegistryError(raw []byte) (*RaylsEnygmaHandlerRaylsAppUnauthorizedTokenRegistry, error) {
	out := new(RaylsEnygmaHandlerRaylsAppUnauthorizedTokenRegistry)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsAppUnauthorizedTokenRegistry", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsAppUserNotRegistered represents a RaylsApp__UserNotRegistered error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsAppUserNotRegistered struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__UserNotRegistered(address caller)
func RaylsEnygmaHandlerRaylsAppUserNotRegisteredErrorID() common.Hash {
	return common.HexToHash("0x4c1db902cce08bec31bedc484362fba54949899ac3c0bf0416f3c44af3284baa")
}

// UnpackRaylsAppUserNotRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__UserNotRegistered(address caller)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsAppUserNotRegisteredError(raw []byte) (*RaylsEnygmaHandlerRaylsAppUserNotRegistered, error) {
	out := new(RaylsEnygmaHandlerRaylsAppUserNotRegistered)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsAppUserNotRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsEnygmaHandlerArrayLengthMismatch represents a RaylsEnygmaHandler__ArrayLengthMismatch error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsEnygmaHandlerArrayLengthMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsEnygmaHandler__ArrayLengthMismatch()
func RaylsEnygmaHandlerRaylsEnygmaHandlerArrayLengthMismatchErrorID() common.Hash {
	return common.HexToHash("0xf87365f1c42254da5c61230d0acdcf6f5c1481f4af3d48a5a6734dc99fbb95c0")
}

// UnpackRaylsEnygmaHandlerArrayLengthMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsEnygmaHandler__ArrayLengthMismatch()
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsEnygmaHandlerArrayLengthMismatchError(raw []byte) (*RaylsEnygmaHandlerRaylsEnygmaHandlerArrayLengthMismatch, error) {
	out := new(RaylsEnygmaHandlerRaylsEnygmaHandlerArrayLengthMismatch)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsEnygmaHandlerArrayLengthMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsEnygmaHandlerAuthorityNotSet represents a RaylsEnygmaHandler__AuthorityNotSet error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsEnygmaHandlerAuthorityNotSet struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsEnygmaHandler__AuthorityNotSet()
func RaylsEnygmaHandlerRaylsEnygmaHandlerAuthorityNotSetErrorID() common.Hash {
	return common.HexToHash("0x197f2287fb63db61cf4ca39fe54a3c7550de7e2909c867f85037c76a84a080d4")
}

// UnpackRaylsEnygmaHandlerAuthorityNotSetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsEnygmaHandler__AuthorityNotSet()
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsEnygmaHandlerAuthorityNotSetError(raw []byte) (*RaylsEnygmaHandlerRaylsEnygmaHandlerAuthorityNotSet, error) {
	out := new(RaylsEnygmaHandlerRaylsEnygmaHandlerAuthorityNotSet)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsEnygmaHandlerAuthorityNotSet", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsEnygmaHandlerCallableExecutionFailed represents a RaylsEnygmaHandler__CallableExecutionFailed error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsEnygmaHandlerCallableExecutionFailed struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsEnygmaHandler__CallableExecutionFailed(address target)
func RaylsEnygmaHandlerRaylsEnygmaHandlerCallableExecutionFailedErrorID() common.Hash {
	return common.HexToHash("0x0aae3f2f3c37c778a0d81aca63827f0ea8f5171a8bd2e407b8b9f77f942cbf9e")
}

// UnpackRaylsEnygmaHandlerCallableExecutionFailedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsEnygmaHandler__CallableExecutionFailed(address target)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsEnygmaHandlerCallableExecutionFailedError(raw []byte) (*RaylsEnygmaHandlerRaylsEnygmaHandlerCallableExecutionFailed, error) {
	out := new(RaylsEnygmaHandlerRaylsEnygmaHandlerCallableExecutionFailed)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsEnygmaHandlerCallableExecutionFailed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsEnygmaHandlerCallableTargetNotContract represents a RaylsEnygmaHandler__CallableTargetNotContract error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsEnygmaHandlerCallableTargetNotContract struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsEnygmaHandler__CallableTargetNotContract(address target)
func RaylsEnygmaHandlerRaylsEnygmaHandlerCallableTargetNotContractErrorID() common.Hash {
	return common.HexToHash("0xc33c35d42d6def7a65c2fbb5174006188c82bdbbba09b467b0669bed8879efe5")
}

// UnpackRaylsEnygmaHandlerCallableTargetNotContractError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsEnygmaHandler__CallableTargetNotContract(address target)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsEnygmaHandlerCallableTargetNotContractError(raw []byte) (*RaylsEnygmaHandlerRaylsEnygmaHandlerCallableTargetNotContract, error) {
	out := new(RaylsEnygmaHandlerRaylsEnygmaHandlerCallableTargetNotContract)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsEnygmaHandlerCallableTargetNotContract", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsEnygmaHandlerEmptyArray represents a RaylsEnygmaHandler__EmptyArray error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsEnygmaHandlerEmptyArray struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsEnygmaHandler__EmptyArray()
func RaylsEnygmaHandlerRaylsEnygmaHandlerEmptyArrayErrorID() common.Hash {
	return common.HexToHash("0x80450d04a90ec96ac78759abca2a12a16be6639a47c1d310ebc4f49f2cab3191")
}

// UnpackRaylsEnygmaHandlerEmptyArrayError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsEnygmaHandler__EmptyArray()
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsEnygmaHandlerEmptyArrayError(raw []byte) (*RaylsEnygmaHandlerRaylsEnygmaHandlerEmptyArray, error) {
	out := new(RaylsEnygmaHandlerRaylsEnygmaHandlerEmptyArray)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsEnygmaHandlerEmptyArray", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsEnygmaHandlerInvalidDecimals represents a RaylsEnygmaHandler__InvalidDecimals error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsEnygmaHandlerInvalidDecimals struct {
	Decimals uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsEnygmaHandler__InvalidDecimals(uint8 decimals)
func RaylsEnygmaHandlerRaylsEnygmaHandlerInvalidDecimalsErrorID() common.Hash {
	return common.HexToHash("0x1ee885564370ab8351e4bc7087f2085d16f7eba19fed3a2bb8b47c090f97194c")
}

// UnpackRaylsEnygmaHandlerInvalidDecimalsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsEnygmaHandler__InvalidDecimals(uint8 decimals)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsEnygmaHandlerInvalidDecimalsError(raw []byte) (*RaylsEnygmaHandlerRaylsEnygmaHandlerInvalidDecimals, error) {
	out := new(RaylsEnygmaHandlerRaylsEnygmaHandlerInvalidDecimals)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsEnygmaHandlerInvalidDecimals", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsEnygmaHandlerNotRelayer represents a RaylsEnygmaHandler__NotRelayer error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsEnygmaHandlerNotRelayer struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsEnygmaHandler__NotRelayer(address caller)
func RaylsEnygmaHandlerRaylsEnygmaHandlerNotRelayerErrorID() common.Hash {
	return common.HexToHash("0x1e3ff3a00242d49848256421b5f6522a8d6f56422e9c01571e0e5df35f1fecc6")
}

// UnpackRaylsEnygmaHandlerNotRelayerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsEnygmaHandler__NotRelayer(address caller)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsEnygmaHandlerNotRelayerError(raw []byte) (*RaylsEnygmaHandlerRaylsEnygmaHandlerNotRelayer, error) {
	out := new(RaylsEnygmaHandlerRaylsEnygmaHandlerNotRelayer)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsEnygmaHandlerNotRelayer", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsEnygmaHandlerNotTokenOwnerScoped represents a RaylsEnygmaHandler__NotTokenOwnerScoped error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsEnygmaHandlerNotTokenOwnerScoped struct {
	OriginSender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsEnygmaHandler__NotTokenOwnerScoped(address originSender)
func RaylsEnygmaHandlerRaylsEnygmaHandlerNotTokenOwnerScopedErrorID() common.Hash {
	return common.HexToHash("0xa7ffd9412f95f7678d6d3fb8444e3100af43272043db80770ba8a7f6afe8abf0")
}

// UnpackRaylsEnygmaHandlerNotTokenOwnerScopedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsEnygmaHandler__NotTokenOwnerScoped(address originSender)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsEnygmaHandlerNotTokenOwnerScopedError(raw []byte) (*RaylsEnygmaHandlerRaylsEnygmaHandlerNotTokenOwnerScoped, error) {
	out := new(RaylsEnygmaHandlerRaylsEnygmaHandlerNotTokenOwnerScoped)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsEnygmaHandlerNotTokenOwnerScoped", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsEnygmaHandlerSwapValidityOutOfRange represents a RaylsEnygmaHandler__SwapValidityOutOfRange error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsEnygmaHandlerSwapValidityOutOfRange struct {
	Provided uint64
	Min      uint64
	Max      uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsEnygmaHandler__SwapValidityOutOfRange(uint64 provided, uint64 min, uint64 max)
func RaylsEnygmaHandlerRaylsEnygmaHandlerSwapValidityOutOfRangeErrorID() common.Hash {
	return common.HexToHash("0xc91017e8ad765a3ef66b23f7a4c528dfb73dbd3281289f2a606fc394f8a7f928")
}

// UnpackRaylsEnygmaHandlerSwapValidityOutOfRangeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsEnygmaHandler__SwapValidityOutOfRange(uint64 provided, uint64 min, uint64 max)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsEnygmaHandlerSwapValidityOutOfRangeError(raw []byte) (*RaylsEnygmaHandlerRaylsEnygmaHandlerSwapValidityOutOfRange, error) {
	out := new(RaylsEnygmaHandlerRaylsEnygmaHandlerSwapValidityOutOfRange)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsEnygmaHandlerSwapValidityOutOfRange", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsEnygmaHandlerTooManyUniqueChainIds represents a RaylsEnygmaHandler__TooManyUniqueChainIds error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsEnygmaHandlerTooManyUniqueChainIds struct {
	Count *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsEnygmaHandler__TooManyUniqueChainIds(uint256 count)
func RaylsEnygmaHandlerRaylsEnygmaHandlerTooManyUniqueChainIdsErrorID() common.Hash {
	return common.HexToHash("0x7305724e97f373d09c260312d3bf3146e46d33f61eee274b72afdbe4343de604")
}

// UnpackRaylsEnygmaHandlerTooManyUniqueChainIdsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsEnygmaHandler__TooManyUniqueChainIds(uint256 count)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsEnygmaHandlerTooManyUniqueChainIdsError(raw []byte) (*RaylsEnygmaHandlerRaylsEnygmaHandlerTooManyUniqueChainIds, error) {
	out := new(RaylsEnygmaHandlerRaylsEnygmaHandlerTooManyUniqueChainIds)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsEnygmaHandlerTooManyUniqueChainIds", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsEnygmaHandlerWrongAddress represents a RaylsEnygmaHandler__WrongAddress error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsEnygmaHandlerWrongAddress struct {
	From common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsEnygmaHandler__WrongAddress(address from)
func RaylsEnygmaHandlerRaylsEnygmaHandlerWrongAddressErrorID() common.Hash {
	return common.HexToHash("0xf23e3251d2554893cf0c10bc2f25eaade1e9778de498bfb544b9cfde526bfb52")
}

// UnpackRaylsEnygmaHandlerWrongAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsEnygmaHandler__WrongAddress(address from)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsEnygmaHandlerWrongAddressError(raw []byte) (*RaylsEnygmaHandlerRaylsEnygmaHandlerWrongAddress, error) {
	out := new(RaylsEnygmaHandlerRaylsEnygmaHandlerWrongAddress)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsEnygmaHandlerWrongAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsEnygmaHandlerWrongFunctionForSameChainId represents a RaylsEnygmaHandler__WrongFunctionForSameChainId error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsEnygmaHandlerWrongFunctionForSameChainId struct {
	ChainId *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsEnygmaHandler__WrongFunctionForSameChainId(uint256 chainId)
func RaylsEnygmaHandlerRaylsEnygmaHandlerWrongFunctionForSameChainIdErrorID() common.Hash {
	return common.HexToHash("0xd6a72faa9ed43b664aa7a5e1e987fc6db8259e790c72c27cae363d517db8caa2")
}

// UnpackRaylsEnygmaHandlerWrongFunctionForSameChainIdError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsEnygmaHandler__WrongFunctionForSameChainId(uint256 chainId)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsEnygmaHandlerWrongFunctionForSameChainIdError(raw []byte) (*RaylsEnygmaHandlerRaylsEnygmaHandlerWrongFunctionForSameChainId, error) {
	out := new(RaylsEnygmaHandlerRaylsEnygmaHandlerWrongFunctionForSameChainId)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsEnygmaHandlerWrongFunctionForSameChainId", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsEnygmaHandlerZeroAddress represents a RaylsEnygmaHandler__ZeroAddress error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsEnygmaHandlerZeroAddress struct {
	Addr common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsEnygmaHandler__ZeroAddress(address addr)
func RaylsEnygmaHandlerRaylsEnygmaHandlerZeroAddressErrorID() common.Hash {
	return common.HexToHash("0x2aee3ba916b38eda95440a7643c66bbe2dd5e14787ca0e16854dc105f8f81d27")
}

// UnpackRaylsEnygmaHandlerZeroAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsEnygmaHandler__ZeroAddress(address addr)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsEnygmaHandlerZeroAddressError(raw []byte) (*RaylsEnygmaHandlerRaylsEnygmaHandlerZeroAddress, error) {
	out := new(RaylsEnygmaHandlerRaylsEnygmaHandlerZeroAddress)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsEnygmaHandlerZeroAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsEnygmaHandlerZeroAmount represents a RaylsEnygmaHandler__ZeroAmount error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsEnygmaHandlerZeroAmount struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsEnygmaHandler__ZeroAmount()
func RaylsEnygmaHandlerRaylsEnygmaHandlerZeroAmountErrorID() common.Hash {
	return common.HexToHash("0x043e4a9425394b0cb7e2067e3a4120bb35b1d193bf5651f5c95e03057c67f4aa")
}

// UnpackRaylsEnygmaHandlerZeroAmountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsEnygmaHandler__ZeroAmount()
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsEnygmaHandlerZeroAmountError(raw []byte) (*RaylsEnygmaHandlerRaylsEnygmaHandlerZeroAmount, error) {
	out := new(RaylsEnygmaHandlerRaylsEnygmaHandlerZeroAmount)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsEnygmaHandlerZeroAmount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerRaylsEnygmaHandlerZeroValueArg represents a RaylsEnygmaHandler__ZeroValueArg error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerRaylsEnygmaHandlerZeroValueArg struct {
	Receiver    common.Address
	Value       *big.Int
	DestChainId *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsEnygmaHandler__ZeroValueArg(address receiver, uint256 value, uint256 destChainId)
func RaylsEnygmaHandlerRaylsEnygmaHandlerZeroValueArgErrorID() common.Hash {
	return common.HexToHash("0x90649856b2fb4e8de79a4017209695b06ab8ee9282b996f3e68cbd7782ef821d")
}

// UnpackRaylsEnygmaHandlerZeroValueArgError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsEnygmaHandler__ZeroValueArg(address receiver, uint256 value, uint256 destChainId)
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackRaylsEnygmaHandlerZeroValueArgError(raw []byte) (*RaylsEnygmaHandlerRaylsEnygmaHandlerZeroValueArg, error) {
	out := new(RaylsEnygmaHandlerRaylsEnygmaHandlerZeroValueArg)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "RaylsEnygmaHandlerZeroValueArg", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RaylsEnygmaHandlerReentrancyGuardReentrantCall represents a ReentrancyGuardReentrantCall error raised by the RaylsEnygmaHandler contract.
type RaylsEnygmaHandlerReentrancyGuardReentrantCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ReentrancyGuardReentrantCall()
func RaylsEnygmaHandlerReentrancyGuardReentrantCallErrorID() common.Hash {
	return common.HexToHash("0x3ee5aeb571de7fc460830b4d0017439a1ca56fb0bc39062227ade4fe4a24c1ca")
}

// UnpackReentrancyGuardReentrantCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ReentrancyGuardReentrantCall()
func (raylsEnygmaHandler *RaylsEnygmaHandler) UnpackReentrancyGuardReentrantCallError(raw []byte) (*RaylsEnygmaHandlerReentrancyGuardReentrantCall, error) {
	out := new(RaylsEnygmaHandlerReentrancyGuardReentrantCall)
	if err := raylsEnygmaHandler.abi.UnpackIntoInterface(out, "ReentrancyGuardReentrantCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}
