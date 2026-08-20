// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package EnygmaPNEvents

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

// SharedObjectsEnygmaProgramData is an auto generated low-level Go binding around an user-defined struct.
type SharedObjectsEnygmaProgramData struct {
	ResourceId      [32]byte
	ContractAddress common.Address
	Selector        [4]byte
	Args            []byte
}

// SharedObjectsPNHTransfer is an auto generated low-level Go binding around an user-defined struct.
type SharedObjectsPNHTransfer struct {
	ResourceId  [32]byte
	Value       []*big.Int
	ToChainId   []*big.Int
	To          []common.Address
	From        common.Address
	ReferenceId [32]byte
	ProgramData [][]SharedObjectsEnygmaProgramData
}

// EnygmaPNEventsMetaData contains all meta data concerning the EnygmaPNEvents contract.
var EnygmaPNEventsMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_endpointAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_participantValidator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_tokenValidator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_authority\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelSwap\",\"inputs\":[{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_toChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_tokenInResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_tokenInAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_tokenInId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_tokenInStandard\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.ErcStandard\"},{\"name\":\"_tokenOutResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_tokenOutAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_tokenOutId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_tokenOutStandard\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.ErcStandard\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"creation\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"initialSupply\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositToDvp\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_referenceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvp1155Burn\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvp1155Creation\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvp1155DepositIntoDvp\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvp1155Mint\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvp1155SwapCompleted\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_destinationChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_destinationOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvp1155SwapForEnygma\",\"inputs\":[{\"name\":\"_tokenResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_tokenValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_tokenData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_enygmaAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_enygmaResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_destChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_validityTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvp1155WithdrawFromDvp\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvp721Burn\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_nftId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvp721Creation\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvp721DepositIntoDvp\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_nftId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvp721Mint\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_nftId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvp721SwapCompleted\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_nftId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_destinationChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_destinationOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvp721SwapForEnygma\",\"inputs\":[{\"name\":\"_nftResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_nftId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_enygmaAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_enygmaResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_destChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_validityTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvp721WithdrawFromDvp\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_nftId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAddressByResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEndpointAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"participantValidator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIParticipantValidator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"raylsNodeUserGovernance\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIUserGovernance\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"resourceId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"revertMint\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_reason\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendTransferPNH\",\"inputs\":[{\"name\":\"_pnhTransfer\",\"type\":\"tuple\",\"internalType\":\"structSharedObjects.PNHTransfer\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"value\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"toChainId\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"to\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"referenceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"programData\",\"type\":\"tuple[][]\",\"internalType\":\"structSharedObjects.EnygmaProgramData[][]\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"args\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setParticipantValidator\",\"inputs\":[{\"name\":\"_participantValidator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setTokenValidator\",\"inputs\":[{\"name\":\"_tokenValidator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"swapWithDvpForERC1155\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_nftId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_nftResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_nftAmountOrOne\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_enygmaAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_destChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_validityTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"swapWithDvpForERC721\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_nftId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_nftResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_enygmaAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_destChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_validityTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tokenValidator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractITokenRegistryValidator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawFromDvp\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_referenceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Dvp1155Burn\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_to\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Dvp1155Creation\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Dvp1155DepositIntoDvp\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_tokenId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Dvp1155Mint\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_tokenId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Dvp1155SwapCompleted\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_tokenId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_destinationChainId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_destinationOwner\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Dvp1155SwapForEnygma\",\"inputs\":[{\"name\":\"_tokenResourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_tokenId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_tokenValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_tokenData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"_enygmaAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_enygmaResourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"_destChainId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_validityTime\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Dvp1155WithdrawFromDvp\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_tokenId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Dvp721Burn\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_nftId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Dvp721Creation\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Dvp721DepositIntoDvp\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_nftId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Dvp721Mint\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_nftId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Dvp721SwapCompleted\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_nftId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_destinationChainId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_destinationOwner\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Dvp721SwapForEnygma\",\"inputs\":[{\"name\":\"_nftResourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_nftId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_enygmaAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_enygmaResourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_from\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"_destChainId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_validityTime\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Dvp721WithdrawFromDvp\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_nftId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DvpSwapCancelled\",\"inputs\":[{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_destChainId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_tokenInResourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_tokenInAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_tokenInId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_tokenInStandard\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumSharedObjects.ErcStandard\"},{\"name\":\"_tokenOutResourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_tokenOutAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_tokenOutId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_tokenOutStandard\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumSharedObjects.ErcStandard\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EnygmaBurn\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_from\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EnygmaCreation\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_initialSupply\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EnygmaDepositToDvp\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_from\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"_referenceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EnygmaMint\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_to\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EnygmaRevertMint\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_to\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"_reason\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EnygmaSendTransferPNH\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_value\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"_toChainId\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"_to\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"_from\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"_referenceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_programData\",\"type\":\"tuple[][]\",\"indexed\":false,\"internalType\":\"structSharedObjects.EnygmaProgramData[][]\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"args\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EnygmaSwapWithDvpForERC1155\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_nftId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_nftResourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_nftAmountOrOne\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_enygmaAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_from\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"_destChainId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_validityTime\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EnygmaSwapWithDvpForERC721\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_nftId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_nftResourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_enygmaAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_from\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"_destChainId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"_validityTime\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EnygmaWithdrawFromDvp\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_to\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"_referenceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenRegistrationSubmitted\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EnygmaPNEvents__ValidatorsNotInitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__HubNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"hubStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__PrivacyNodeFrozen\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__PrivacyNodeNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__PublicChainNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"publicChainStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__ResourceNotApproved\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsApp__TokenRegistryNotConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsApp__UnauthorizedTokenRegistry\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenRegistry\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__UserNotRegistered\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	ID:  "EnygmaPNEvents",
	Bin: "0x60806040523480156200001157600080fd5b506040516200416b3803806200416b833981016040819052620000349162000247565b600080546001600160a01b0319166001600160a01b038616178155849050506001600160a01b038316156200007f57600480546001600160a01b0319166001600160a01b0385161790555b6001600160a01b03821615620000ab57600580546001600160a01b0319166001600160a01b0384161790555b6001600160a01b03811615620000c657620000c681620000d0565b50505050620002cc565b6000620000dc62000128565b80549091506001600160a01b0316156200011957604051638944034760e01b81526001600160a01b03831660048201526024015b60405180910390fd5b62000124826200018d565b5050565b60008060ff196200015b60017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f35620002a4565b6040516020016200016e91815260200190565b60408051601f1981840301815291905280516020909101201692915050565b6001600160a01b038116620001c157604051638944034760e01b81526001600160a01b038216600482015260240162000110565b6000620001cd62000128565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b80516001600160a01b03811681146200024257600080fd5b919050565b600080600080608085870312156200025e57600080fd5b62000269856200022a565b935062000279602086016200022a565b925062000289604086016200022a565b915062000299606086016200022a565b905092959194509250565b81810381811115620002c657634e487b7160e01b600052601160045260246000fd5b92915050565b613e8f80620002dc6000396000f3fe608060405234801561001057600080fd5b50600436106101d15760003560e01c806344dc1ce11161010557806380691fba1161009d57806380691fba146103e05780638129fc1c146103f3578063a01afbfb146103fb578063bb9c53ae1461040e578063bf7e214f14610421578063ce884eb514610429578063d1f2ada11461043a578063ea50f46b1461044d578063f205b8701461046057600080fd5b806344dc1ce11461033157806356249c36146103445780635f997c5b1461035757806368f797c91461036e5780636b41b261146103815780636f5be509146103945780637441a427146103a7578063774dcd89146103ba5780637ed9db59146103cd57600080fd5b806322647cd91161017857806322647cd914610286578063268f2e41146102995780632cbc8600146102ac5780633405432c146102bf578063346a9074146102d25780633a048283146102e55780633a06a3f4146102f85780633caf04761461030b5780633d75fb4c1461031e57600080fd5b806302136b01146101d65780630274c133146101eb57806305a3b856146102145780630bcca0d2146102275780630fbfc6611461023a57806311f50c851461024d5780631a8997f7146102605780631e8fcfef14610273575b600080fd5b6101e96101e4366004612eda565b610473565b005b6002546101fe906001600160a01b031681565b60405161020b9190612f13565b60405180910390f35b6101e9610222366004612f27565b610781565b6101e9610235366004612f40565b6107ce565b6101e961024836600461308d565b610adf565b6101fe61025b366004612f27565b610df3565b6101e961026e3660046130fb565b610e67565b6101e961028136600461313c565b610ec0565b6101e96102943660046131b7565b6111c2565b6101e96102a7366004613239565b61122d565b6101e96102ba366004613256565b611294565b6101e96102cd366004613278565b6112e8565b6101e96102e03660046132d1565b611333565b6101e96102f3366004613586565b611389565b6101e9610306366004612eda565b6116a0565b6101e9610319366004613683565b61199e565b6101e961032c366004613704565b6119e9565b6101e961033f3660046130fb565b611d0c565b6101e9610352366004613256565b611d57565b61036060035481565b60405190815260200161020b565b6101e961037c366004612f40565b611da3565b6101e961038f366004613256565b6120a3565b6004546101fe906001600160a01b031681565b6101e96103b5366004613798565b6120ef565b6101e96103c8366004612f27565b612152565b6101e96103db3660046132d1565b612198565b6005546101fe906001600160a01b031681565b6101e96121e1565b6101e9610409366004612f27565b61222e565b6101e961041c3660046137d5565b612279565b6101fe6125cf565b6000546001600160a01b03166101fe565b6101e9610448366004613239565b6125e8565b6101e961045b36600461384e565b612646565b6101e961046e36600461384e565b612963565b610489336000356001600160e01b031916612c6b565b6040805160008152602081019091526004548491906001600160a01b031615806104bc57506005546001600160a01b0316155b156104da57604051631385522360e11b815260040160405180910390fd5b60008060009054906101000a90046001600160a01b03166001600160a01b0316633408e4706040518163ffffffff1660e01b8152600401602060405180830381865afa15801561052e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061055291906138bd565b600480546040516364aabfab60e11b81529293506001600160a01b03169163c9557f56916105869185910190815260200190565b60006040518083038186803b15801561059e57600080fd5b505afa1580156105b2573d6000803e3d6000fd5b505060055460405163030c0d3b60e61b815260048101879052602481018590526001600160a01b03909116925063c3034ec09150604401600060405180830381600087803b15801561060357600080fd5b505af1158015610617573d6000803e3d6000fd5b5050505060005b825181101561073d5760045483516001600160a01b039091169063c9557f5690859084908110610650576106506138d6565b60200260200101516040518263ffffffff1660e01b815260040161067691815260200190565b60006040518083038186803b15801561068e57600080fd5b505afa1580156106a2573d6000803e3d6000fd5b505060055485516001600160a01b03909116925063c3034ec0915086908690859081106106d1576106d16138d6565b60200260200101516040518363ffffffff1660e01b81526004016106ff929190918252602082015260400190565b600060405180830381600087803b15801561071957600080fd5b505af115801561072d573d6000803e3d6000fd5b50506001909201915061061e9050565b507fa16d9a57158a130490ee3d64183abaab5c62c18e70bbc7828c7092dfc7de80c3868686604051610771939291906138ec565b60405180910390a1505050505050565b610797336000356001600160e01b031916612c6b565b6040518181527f81ecc951bab0fae8817006d04acd6a09aeac9e29a57a9bacdb2629bf462c8ec3906020015b60405180910390a150565b6107e4336000356001600160e01b031916612c6b565b6040805160008152602081019091526004548591906001600160a01b0316158061081757506005546001600160a01b0316155b1561083557604051631385522360e11b815260040160405180910390fd5b60008060009054906101000a90046001600160a01b03166001600160a01b0316633408e4706040518163ffffffff1660e01b8152600401602060405180830381865afa158015610889573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108ad91906138bd565b600480546040516364aabfab60e11b81529293506001600160a01b03169163c9557f56916108e19185910190815260200190565b60006040518083038186803b1580156108f957600080fd5b505afa15801561090d573d6000803e3d6000fd5b505060055460405163030c0d3b60e61b815260048101879052602481018590526001600160a01b03909116925063c3034ec09150604401600060405180830381600087803b15801561095e57600080fd5b505af1158015610972573d6000803e3d6000fd5b5050505060005b8251811015610a985760045483516001600160a01b039091169063c9557f56908590849081106109ab576109ab6138d6565b60200260200101516040518263ffffffff1660e01b81526004016109d191815260200190565b60006040518083038186803b1580156109e957600080fd5b505afa1580156109fd573d6000803e3d6000fd5b505060055485516001600160a01b03909116925063c3034ec091508690869085908110610a2c57610a2c6138d6565b60200260200101516040518363ffffffff1660e01b8152600401610a5a929190918252602082015260400190565b600060405180830381600087803b158015610a7457600080fd5b505af1158015610a88573d6000803e3d6000fd5b5050600190920191506109799050565b507fcbe113947dfc01fcaf958ae1475972c6c1a1a83505231bb5c58e8b56d19c4ddb87878787604051610ace949392919061390b565b60405180910390a150505050505050565b610af5336000356001600160e01b031916612c6b565b6040805160008152602081019091526004548691906001600160a01b03161580610b2857506005546001600160a01b0316155b15610b4657604051631385522360e11b815260040160405180910390fd5b60008060009054906101000a90046001600160a01b03166001600160a01b0316633408e4706040518163ffffffff1660e01b8152600401602060405180830381865afa158015610b9a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bbe91906138bd565b600480546040516364aabfab60e11b81529293506001600160a01b03169163c9557f5691610bf29185910190815260200190565b60006040518083038186803b158015610c0a57600080fd5b505afa158015610c1e573d6000803e3d6000fd5b505060055460405163030c0d3b60e61b815260048101879052602481018590526001600160a01b03909116925063c3034ec09150604401600060405180830381600087803b158015610c6f57600080fd5b505af1158015610c83573d6000803e3d6000fd5b5050505060005b8251811015610da95760045483516001600160a01b039091169063c9557f5690859084908110610cbc57610cbc6138d6565b60200260200101516040518263ffffffff1660e01b8152600401610ce291815260200190565b60006040518083038186803b158015610cfa57600080fd5b505afa158015610d0e573d6000803e3d6000fd5b505060055485516001600160a01b03909116925063c3034ec091508690869085908110610d3d57610d3d6138d6565b60200260200101516040518363ffffffff1660e01b8152600401610d6b929190918252602082015260400190565b600060405180830381600087803b158015610d8557600080fd5b505af1158015610d99573d6000803e3d6000fd5b505060019092019150610c8a9050565b507feeac5cbcc1774d60218c32fe342cabbe48bff8a102c772fa005e08ddcd7664f58888888888604051610de1959493929190613975565b60405180910390a15050505050505050565b600080546040516311f50c8560e01b8152600481018490526001600160a01b03909116906311f50c8590602401602060405180830381865afa158015610e3d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e6191906139b4565b92915050565b610e7d336000356001600160e01b031916612c6b565b7f877077a253ed088b196d0f3f5c140386a6805456a27eaf6c39f8aee3da349c5a84848484604051610eb294939291906139d1565b60405180910390a150505050565b610ed6336000356001600160e01b031916612c6b565b6040805160008152602081019091526004548691906001600160a01b03161580610f0957506005546001600160a01b0316155b15610f2757604051631385522360e11b815260040160405180910390fd5b60008060009054906101000a90046001600160a01b03166001600160a01b0316633408e4706040518163ffffffff1660e01b8152600401602060405180830381865afa158015610f7b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f9f91906138bd565b600480546040516364aabfab60e11b81529293506001600160a01b03169163c9557f5691610fd39185910190815260200190565b60006040518083038186803b158015610feb57600080fd5b505afa158015610fff573d6000803e3d6000fd5b505060055460405163030c0d3b60e61b815260048101879052602481018590526001600160a01b03909116925063c3034ec09150604401600060405180830381600087803b15801561105057600080fd5b505af1158015611064573d6000803e3d6000fd5b5050505060005b825181101561118a5760045483516001600160a01b039091169063c9557f569085908490811061109d5761109d6138d6565b60200260200101516040518263ffffffff1660e01b81526004016110c391815260200190565b60006040518083038186803b1580156110db57600080fd5b505afa1580156110ef573d6000803e3d6000fd5b505060055485516001600160a01b03909116925063c3034ec09150869086908590811061111e5761111e6138d6565b60200260200101516040518363ffffffff1660e01b815260040161114c929190918252602082015260400190565b600060405180830381600087803b15801561116657600080fd5b505af115801561117a573d6000803e3d6000fd5b50506001909201915061106b9050565b507ffc3ff5eda49d75e9d98769b0e6f3ca948135aa1bbdcb6087b0fb902c125c782f8888888888604051610de19594939291906139f5565b6111d8336000356001600160e01b031916612c6b565b7f6809767585d1896260b21e2c556a054cc296802ac810f7282c01e382c1191c2e8a8a8a8a8a8a8a8a8a8a6040516112199a99989796959493929190613a55565b60405180910390a150505050505050505050565b611243336000356001600160e01b031916612c6b565b6001600160a01b0381166112725760405162461bcd60e51b815260040161126990613ab6565b60405180910390fd5b600580546001600160a01b0319166001600160a01b0392909216919091179055565b6112aa336000356001600160e01b031916612c6b565b60408051838152602081018390527fbb2b024fb4246e81f90b1f3c31ca27f498ecc010a37356250149c8752041b4b391015b60405180910390a15050565b6112fe336000356001600160e01b031916612c6b565b7fbf7f8235b34fcdd6aa654fdace8d312932fbcdec1c0431d3d99fbc3c93fa130084848484604051610eb29493929190613aed565b611349336000356001600160e01b031916612c6b565b7fa689af8eaab5e7860e64b80da0812e409b2c7193b5d81c84a0561ce837c868d583838360405161137c93929190613b1c565b60405180910390a1505050565b61139f336000356001600160e01b031916612c6b565b805160408201516004546001600160a01b031615806113c757506005546001600160a01b0316155b156113e557604051631385522360e11b815260040160405180910390fd5b60008060009054906101000a90046001600160a01b03166001600160a01b0316633408e4706040518163ffffffff1660e01b8152600401602060405180830381865afa158015611439573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061145d91906138bd565b600480546040516364aabfab60e11b81529293506001600160a01b03169163c9557f56916114919185910190815260200190565b60006040518083038186803b1580156114a957600080fd5b505afa1580156114bd573d6000803e3d6000fd5b505060055460405163030c0d3b60e61b815260048101879052602481018590526001600160a01b03909116925063c3034ec09150604401600060405180830381600087803b15801561150e57600080fd5b505af1158015611522573d6000803e3d6000fd5b5050505060005b82518110156116485760045483516001600160a01b039091169063c9557f569085908490811061155b5761155b6138d6565b60200260200101516040518263ffffffff1660e01b815260040161158191815260200190565b60006040518083038186803b15801561159957600080fd5b505afa1580156115ad573d6000803e3d6000fd5b505060055485516001600160a01b03909116925063c3034ec0915086908690859081106115dc576115dc6138d6565b60200260200101516040518363ffffffff1660e01b815260040161160a929190918252602082015260400190565b600060405180830381600087803b15801561162457600080fd5b505af1158015611638573d6000803e3d6000fd5b5050600190920191506115299050565b507f2ec9ad4f20033b3f239ff38f84f08388760bdc8bb2613ab2d8289c8775b82ea3846000015185602001518660400151876060015188608001518960a001518a60c00151604051610eb29796959493929190613c50565b6116b6336000356001600160e01b031916612c6b565b6040805160008152602081019091526004548491906001600160a01b031615806116e957506005546001600160a01b0316155b1561170757604051631385522360e11b815260040160405180910390fd5b60008060009054906101000a90046001600160a01b03166001600160a01b0316633408e4706040518163ffffffff1660e01b8152600401602060405180830381865afa15801561175b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061177f91906138bd565b600480546040516364aabfab60e11b81529293506001600160a01b03169163c9557f56916117b39185910190815260200190565b60006040518083038186803b1580156117cb57600080fd5b505afa1580156117df573d6000803e3d6000fd5b505060055460405163030c0d3b60e61b815260048101879052602481018590526001600160a01b03909116925063c3034ec09150604401600060405180830381600087803b15801561183057600080fd5b505af1158015611844573d6000803e3d6000fd5b5050505060005b825181101561196a5760045483516001600160a01b039091169063c9557f569085908490811061187d5761187d6138d6565b60200260200101516040518263ffffffff1660e01b81526004016118a391815260200190565b60006040518083038186803b1580156118bb57600080fd5b505afa1580156118cf573d6000803e3d6000fd5b505060055485516001600160a01b03909116925063c3034ec0915086908690859081106118fe576118fe6138d6565b60200260200101516040518363ffffffff1660e01b815260040161192c929190918252602082015260400190565b600060405180830381600087803b15801561194657600080fd5b505af115801561195a573d6000803e3d6000fd5b50506001909201915061184b9050565b507f292c457c9524097b422e1e59c971d0f15c4c395321ca324c634b3ecb2a31c324868686604051610771939291906138ec565b6119b4336000356001600160e01b031916612c6b565b7f489efbdeed278783d61c02fad16bbe7839e6416f21f2f4b07b4312fde2209cfc84848484604051610eb29493929190613cf6565b6119ff336000356001600160e01b031916612c6b565b6040805160008152602081019091526004548b91906001600160a01b03161580611a3257506005546001600160a01b0316155b15611a5057604051631385522360e11b815260040160405180910390fd5b60008060009054906101000a90046001600160a01b03166001600160a01b0316633408e4706040518163ffffffff1660e01b8152600401602060405180830381865afa158015611aa4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611ac891906138bd565b600480546040516364aabfab60e11b81529293506001600160a01b03169163c9557f5691611afc9185910190815260200190565b60006040518083038186803b158015611b1457600080fd5b505afa158015611b28573d6000803e3d6000fd5b505060055460405163030c0d3b60e61b815260048101879052602481018590526001600160a01b03909116925063c3034ec09150604401600060405180830381600087803b158015611b7957600080fd5b505af1158015611b8d573d6000803e3d6000fd5b5050505060005b8251811015611cb35760045483516001600160a01b039091169063c9557f5690859084908110611bc657611bc66138d6565b60200260200101516040518263ffffffff1660e01b8152600401611bec91815260200190565b60006040518083038186803b158015611c0457600080fd5b505afa158015611c18573d6000803e3d6000fd5b505060055485516001600160a01b03909116925063c3034ec091508690869085908110611c4757611c476138d6565b60200260200101516040518363ffffffff1660e01b8152600401611c75929190918252602082015260400190565b600060405180830381600087803b158015611c8f57600080fd5b505af1158015611ca3573d6000803e3d6000fd5b505060019092019150611b949050565b507f788550383917da23eb543efdf0ce3a785143129ddce3c3d33059258087972dfa8d8d8d8d8d8d8d8d8d8d604051611cf59a99989796959493929190613d23565b60405180910390a150505050505050505050505050565b611d22336000356001600160e01b031916612c6b565b7fc8f7f34bcf19399e27dd93e8aaca23900b097730e30bc05970a1d4a5cf8c771284848484604051610eb294939291906139d1565b611d6d336000356001600160e01b031916612c6b565b60408051838152602081018390527fe795ea58eacd609a5f453d52205d2ee157c5acc8053020df74afb5b4be1f071691016112dc565b611db9336000356001600160e01b031916612c6b565b6040805160008152602081019091526004548591906001600160a01b03161580611dec57506005546001600160a01b0316155b15611e0a57604051631385522360e11b815260040160405180910390fd5b60008060009054906101000a90046001600160a01b03166001600160a01b0316633408e4706040518163ffffffff1660e01b8152600401602060405180830381865afa158015611e5e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e8291906138bd565b600480546040516364aabfab60e11b81529293506001600160a01b03169163c9557f5691611eb69185910190815260200190565b60006040518083038186803b158015611ece57600080fd5b505afa158015611ee2573d6000803e3d6000fd5b505060055460405163030c0d3b60e61b815260048101879052602481018590526001600160a01b03909116925063c3034ec09150604401600060405180830381600087803b158015611f3357600080fd5b505af1158015611f47573d6000803e3d6000fd5b5050505060005b825181101561206d5760045483516001600160a01b039091169063c9557f5690859084908110611f8057611f806138d6565b60200260200101516040518263ffffffff1660e01b8152600401611fa691815260200190565b60006040518083038186803b158015611fbe57600080fd5b505afa158015611fd2573d6000803e3d6000fd5b505060055485516001600160a01b03909116925063c3034ec091508690869085908110612001576120016138d6565b60200260200101516040518363ffffffff1660e01b815260040161202f929190918252602082015260400190565b600060405180830381600087803b15801561204957600080fd5b505af115801561205d573d6000803e3d6000fd5b505060019092019150611f4e9050565b507fc5eabd9271ba22cb2104f54a35dfb521e8f629a9e5426eb05e6c682b2aa0f5ad87878787604051610ace949392919061390b565b6120b9336000356001600160e01b031916612c6b565b60408051838152602081018390527ffb80b8c897ce0582b5ddd9a40295f9dc7434f24d5923b18fe60fd15c745ad97f91016112dc565b612105336000356001600160e01b031916612c6b565b604080518581526001600160a01b0385166020820152908101839052606081018290527f5f82d5a5f3e0e200fcbe799a0799c5e07f65e8cb7253daf5d8f0edfc6079e1aa90608001610eb2565b612168336000356001600160e01b031916612c6b565b6040518181527fd93bc166fc42679794d5330004ffe6aeeb2b394844a6a94109c4f9344e05e6ab906020016107c3565b6121ae336000356001600160e01b031916612c6b565b7fd757181b1dd1e2639c416dbb54f82e264b70edf0903b198df93fab075771da4083838360405161137c93929190613b1c565b600354156122275760405162461bcd60e51b8152602060048201526013602482015272105b1c9958591e481a5b9a5d1a585b1a5e9959606a1b6044820152606401611269565b6002600355565b6000612238612db6565b9050336001600160a01b038216146122735760405162c2a4c960e31b81523360048201526001600160a01b0382166024820152604401611269565b50600355565b61228f336000356001600160e01b031916612c6b565b6040805160008152602081019091526004548a91906001600160a01b031615806122c257506005546001600160a01b0316155b156122e057604051631385522360e11b815260040160405180910390fd5b60008060009054906101000a90046001600160a01b03166001600160a01b0316633408e4706040518163ffffffff1660e01b8152600401602060405180830381865afa158015612334573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061235891906138bd565b600480546040516364aabfab60e11b81529293506001600160a01b03169163c9557f569161238c9185910190815260200190565b60006040518083038186803b1580156123a457600080fd5b505afa1580156123b8573d6000803e3d6000fd5b505060055460405163030c0d3b60e61b815260048101879052602481018590526001600160a01b03909116925063c3034ec09150604401600060405180830381600087803b15801561240957600080fd5b505af115801561241d573d6000803e3d6000fd5b5050505060005b82518110156125435760045483516001600160a01b039091169063c9557f5690859084908110612456576124566138d6565b60200260200101516040518263ffffffff1660e01b815260040161247c91815260200190565b60006040518083038186803b15801561249457600080fd5b505afa1580156124a8573d6000803e3d6000fd5b505060055485516001600160a01b03909116925063c3034ec0915086908690859081106124d7576124d76138d6565b60200260200101516040518363ffffffff1660e01b8152600401612505929190918252602082015260400190565b600060405180830381600087803b15801561251f57600080fd5b505af1158015612533573d6000803e3d6000fd5b5050600190920191506124249050565b50604080518d8152602081018d90529081018b9052606081018a9052608081018990526001600160a01b03881660a082015260c0810187905260e081018690526001600160401b0385166101008201527f2f1f0f0acdefebb35ee4e21c054f904602dd3fa4e54bbae0bddca3f2d64030bc906101200160405180910390a1505050505050505050505050565b60006125d9612e50565b546001600160a01b0316919050565b6125fe336000356001600160e01b031916612c6b565b6001600160a01b0381166126245760405162461bcd60e51b815260040161126990613ab6565b600480546001600160a01b0319166001600160a01b0392909216919091179055565b61265c336000356001600160e01b031916612c6b565b6040805160008152602081019091526004548991906001600160a01b0316158061268f57506005546001600160a01b0316155b156126ad57604051631385522360e11b815260040160405180910390fd5b60008060009054906101000a90046001600160a01b03166001600160a01b0316633408e4706040518163ffffffff1660e01b8152600401602060405180830381865afa158015612701573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061272591906138bd565b600480546040516364aabfab60e11b81529293506001600160a01b03169163c9557f56916127599185910190815260200190565b60006040518083038186803b15801561277157600080fd5b505afa158015612785573d6000803e3d6000fd5b505060055460405163030c0d3b60e61b815260048101879052602481018590526001600160a01b03909116925063c3034ec09150604401600060405180830381600087803b1580156127d657600080fd5b505af11580156127ea573d6000803e3d6000fd5b5050505060005b82518110156129105760045483516001600160a01b039091169063c9557f5690859084908110612823576128236138d6565b60200260200101516040518263ffffffff1660e01b815260040161284991815260200190565b60006040518083038186803b15801561286157600080fd5b505afa158015612875573d6000803e3d6000fd5b505060055485516001600160a01b03909116925063c3034ec0915086908690859081106128a4576128a46138d6565b60200260200101516040518363ffffffff1660e01b81526004016128d2929190918252602082015260400190565b600060405180830381600087803b1580156128ec57600080fd5b505af1158015612900573d6000803e3d6000fd5b5050600190920191506127f19050565b507f28558ffad93cabce71b4a8c94795ba0063c4d6def318ba00f1369d8f46784d4a8b8b8b8b8b8b8b8b60405161294e989796959493929190613d92565b60405180910390a15050505050505050505050565b612979336000356001600160e01b031916612c6b565b6040805160008152602081019091526004548991906001600160a01b031615806129ac57506005546001600160a01b0316155b156129ca57604051631385522360e11b815260040160405180910390fd5b60008060009054906101000a90046001600160a01b03166001600160a01b0316633408e4706040518163ffffffff1660e01b8152600401602060405180830381865afa158015612a1e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612a4291906138bd565b600480546040516364aabfab60e11b81529293506001600160a01b03169163c9557f5691612a769185910190815260200190565b60006040518083038186803b158015612a8e57600080fd5b505afa158015612aa2573d6000803e3d6000fd5b505060055460405163030c0d3b60e61b815260048101879052602481018590526001600160a01b03909116925063c3034ec09150604401600060405180830381600087803b158015612af357600080fd5b505af1158015612b07573d6000803e3d6000fd5b5050505060005b8251811015612c2d5760045483516001600160a01b039091169063c9557f5690859084908110612b4057612b406138d6565b60200260200101516040518263ffffffff1660e01b8152600401612b6691815260200190565b60006040518083038186803b158015612b7e57600080fd5b505afa158015612b92573d6000803e3d6000fd5b505060055485516001600160a01b03909116925063c3034ec091508690869085908110612bc157612bc16138d6565b60200260200101516040518363ffffffff1660e01b8152600401612bef929190918252602082015260400190565b600060405180830381600087803b158015612c0957600080fd5b505af1158015612c1d573d6000803e3d6000fd5b505060019092019150612b0e9050565b507f0ced5127db9101ac3fe0785ca9a8ee2168320a7e05a85077dc9a959f45f6dbca8b8b8b8b8b8b8b8b60405161294e989796959493929190613d92565b6000612c75612e50565b80549091506001600160a01b031680612ca4576000604051638944034760e01b81526004016112699190612f13565b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa158015612d08573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d2c9190613dea565b92509250925082612dad578015612d565760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff821615612d925760405163a426878960e01b81526001600160a01b038816600482015263ffffffff83166024820152604401611269565b86604051632ecd3d0360e21b81526004016112699190612f13565b50505050505050565b600080546040516311f50c8560e01b8152600360048201526001600160a01b03909116906311f50c8590602401602060405180830381865afa158015612e00573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612e2491906139b4565b90506001600160a01b038116612e4d576040516336a41bd160e01b815260040160405180910390fd5b90565b60008060ff19612e8160017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f35613e38565b604051602001612e9391815260200190565b60408051601f1981840301815291905280516020909101201692915050565b6001600160a01b0381168114612ec757600080fd5b50565b8035612ed581612eb2565b919050565b600080600060608486031215612eef57600080fd5b83359250602084013591506040840135612f0881612eb2565b809150509250925092565b6001600160a01b0391909116815260200190565b600060208284031215612f3957600080fd5b5035919050565b60008060008060808587031215612f5657600080fd5b84359350602085013592506040850135612f6f81612eb2565b9396929550929360600135925050565b634e487b7160e01b600052604160045260246000fd5b604051608081016001600160401b0381118282101715612fb757612fb7612f7f565b60405290565b60405160e081016001600160401b0381118282101715612fb757612fb7612f7f565b604051601f8201601f191681016001600160401b038111828210171561300757613007612f7f565b604052919050565b60006001600160401b0383111561302857613028612f7f565b61303b601f8401601f1916602001612fdf565b905082815283838301111561304f57600080fd5b828260208301376000602084830101529392505050565b600082601f83011261307757600080fd5b6130868383356020850161300f565b9392505050565b600080600080600060a086880312156130a557600080fd5b85359450602086013593506040860135925060608601356001600160401b038111156130d057600080fd5b6130dc88828901613066565b92505060808601356130ed81612eb2565b809150509295509295909350565b6000806000806080858703121561311157600080fd5b843593506020850135925060408501359150606085013561313181612eb2565b939692955090935050565b600080600080600060a0868803121561315457600080fd5b8535945060208601359350604086013561316d81612eb2565b92506060860135915060808601356001600160401b0381111561318f57600080fd5b61319b88828901613066565b9150509295509295909350565b8035600d8110612ed557600080fd5b6000806000806000806000806000806101408b8d0312156131d757600080fd5b8a35995060208b0135985060408b0135975060608b0135965060808b0135955061320360a08c016131a8565b945060c08b0135935060e08b013592506101008b013591506132286101208c016131a8565b90509295989b9194979a5092959850565b60006020828403121561324b57600080fd5b813561308681612eb2565b6000806040838503121561326957600080fd5b50508035926020909101359150565b6000806000806080858703121561328e57600080fd5b84359350602085013592506040850135915060608501356001600160401b038111156132b957600080fd5b6132c587828801613066565b91505092959194509250565b6000806000606084860312156132e657600080fd5b8335925060208401356132f881612eb2565b929592945050506040919091013590565b60006001600160401b0382111561332257613322612f7f565b5060051b60200190565b600082601f83011261333d57600080fd5b8135602061335261334d83613309565b612fdf565b8083825260208201915060208460051b87010193508684111561337457600080fd5b602086015b848110156133905780358352918301918301613379565b509695505050505050565b600082601f8301126133ac57600080fd5b813560206133bc61334d83613309565b8083825260208201915060208460051b8701019350868411156133de57600080fd5b602086015b848110156133905780356133f681612eb2565b83529183019183016133e3565b600082601f83011261341457600080fd5b61342161334d8335613309565b82358082526020808301929160051b85010185101561343f57600080fd5b602084015b6020853560051b86010181101561357d576001600160401b03808235111561346b57600080fd5b8135860187603f82011261347e57600080fd5b61348e61334d6020830135613309565b602082810135808352908201919060051b83016040018a8111156134b157600080fd5b604084015b818110156135665785813511156134cc57600080fd5b803585016080818e03603f190112156134e457600080fd5b6134ec612f95565b604082013581526135006060830135612eb2565b6060820135602082015260808201356001600160e01b031981161461352457600080fd5b608082013560408201528760a0830135111561353f57600080fd5b6135528e604060a0850135850101613066565b6060820152855250602093840193016134b6565b505086525050602093840193919091019050613444565b50949350505050565b60006020828403121561359857600080fd5b81356001600160401b03808211156135af57600080fd5b9083019060e082860312156135c357600080fd5b6135cb612fbd565b823581526020830135828111156135e157600080fd5b6135ed8782860161332c565b60208301525060408301358281111561360557600080fd5b6136118782860161332c565b60408301525060608301358281111561362957600080fd5b6136358782860161339b565b60608301525061364760808401612eca565b608082015260a083013560a082015260c08301358281111561366857600080fd5b61367487828601613403565b60c08301525095945050505050565b6000806000806080858703121561369957600080fd5b843593506020850135925060408501356136b281612eb2565b915060608501356001600160401b038111156136cd57600080fd5b8501601f810187136136de57600080fd5b6132c58782356020840161300f565b80356001600160401b0381168114612ed557600080fd5b6000806000806000806000806000806101408b8d03121561372457600080fd5b8a35995060208b0135985060408b0135975060608b01356001600160401b0381111561374f57600080fd5b61375b8d828e01613066565b97505060808b0135955060a08b0135945060c08b013561377a81612eb2565b935060e08b013592506101008b013591506132286101208c016136ed565b600080600080608085870312156137ae57600080fd5b8435935060208501356137c081612eb2565b93969395505050506040820135916060013590565b60008060008060008060008060006101208a8c0312156137f457600080fd5b8935985060208a0135975060408a0135965060608a0135955060808a0135945060a08a013561382281612eb2565b935060c08a0135925060e08a0135915061383f6101008b016136ed565b90509295985092959850929598565b600080600080600080600080610100898b03121561386b57600080fd5b88359750602089013596506040890135955060608901359450608089013561389281612eb2565b935060a0890135925060c089013591506138ae60e08a016136ed565b90509295985092959890939650565b6000602082840312156138cf57600080fd5b5051919050565b634e487b7160e01b600052603260045260246000fd5b92835260208301919091526001600160a01b0316604082015260600190565b93845260208401929092526001600160a01b03166040830152606082015260800190565b6000815180845260005b8181101561395557602081850181015186830182015201613939565b506000602082860101526020601f19601f83011685010191505092915050565b85815284602082015283604082015260a06060820152600061399a60a083018561392f565b905060018060a01b03831660808301529695505050505050565b6000602082840312156139c657600080fd5b815161308681612eb2565b938452602084019290925260408301526001600160a01b0316606082015260800190565b85815284602082015260018060a01b038416604082015282606082015260a060808201526000613a2860a083018461392f565b979650505050505050565b600d8110613a5157634e487b7160e01b600052602160045260246000fd5b9052565b6000610140820190508b82528a6020830152896040830152886060830152876080830152613a8660a0830188613a33565b8560c08301528460e083015283610100830152613aa7610120830184613a33565b9b9a5050505050505050505050565b6020808252601c908201527f456e79676d61504e4576656e74733a207a65726f206164647265737300000000604082015260600190565b848152836020820152826040820152608060608201526000613b12608083018461392f565b9695505050505050565b9283526001600160a01b03919091166020830152604082015260600190565b60008151808452602080850194506020840160005b83811015613b6c57815187529582019590820190600101613b50565b509495945050505050565b6000828251808552602080860195506005818360051b8501018287016000805b86811015613c4157601f1988850381018c5283518051808752908801908887019080891b88018a01865b82811015613c2a57898203860184528451805183528c8101516001600160a01b03168d8401526040808201516001600160e01b03191690840152606090810151608091840182905290613c168185018361392f565b968e0196958e019593505050600101613bc1565b509e8a019e97505050938701935050600101613b97565b50919998505050505050505050565b8781526000602060e06020840152613c6b60e084018a613b3b565b8381036040850152613c7d818a613b3b565b848103606086015288518082526020808b0193509091019060005b81811015613cbd5783516001600160a01b031683529284019291840191600101613c98565b50506001600160a01b03881660808601528660a086015284810360c0860152613ce68187613b77565b9c9b505050505050505050505050565b84815283602082015260018060a01b0383166040820152608060608201526000613b12608083018461392f565b60006101408c83528b60208401528a6040840152806060840152613d498184018b61392f565b6080840199909952505060a08101959095526001600160a01b039390931660c085015260e08401919091526101008301526001600160401b031661012090910152949350505050565b9788526020880196909652604087019490945260608601929092526001600160a01b0316608085015260a084015260c08301526001600160401b031660e08201526101000190565b80518015158114612ed557600080fd5b600080600060608486031215613dff57600080fd5b613e0884613dda565b9250602084015163ffffffff81168114613e2157600080fd5b9150613e2f60408501613dda565b90509250925092565b81810381811115610e6157634e487b7160e01b600052601160045260246000fdfea264697066735822122014de77ed38b7f912ad9b75a9b70dc6980ac880ba7f926945d4959bd142878c5764736f6c63430008180033",
}

// EnygmaPNEvents is an auto generated Go binding around an Ethereum contract.
type EnygmaPNEvents struct {
	abi abi.ABI
}

// NewEnygmaPNEvents creates a new instance of EnygmaPNEvents.
func NewEnygmaPNEvents() *EnygmaPNEvents {
	parsed, err := EnygmaPNEventsMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &EnygmaPNEvents{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *EnygmaPNEvents) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address _endpointAddress, address _participantValidator, address _tokenValidator, address _authority) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackConstructor(_endpointAddress common.Address, _participantValidator common.Address, _tokenValidator common.Address, _authority common.Address) []byte {
	enc, err := enygmaPNEvents.abi.Pack("", _endpointAddress, _participantValidator, _tokenValidator, _authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (enygmaPNEvents *EnygmaPNEvents) PackAuthority() []byte {
	enc, err := enygmaPNEvents.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (enygmaPNEvents *EnygmaPNEvents) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := enygmaPNEvents.abi.Unpack("authority", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackBurn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x346a9074.
//
// Solidity: function burn(bytes32 _resourceId, address _from, uint256 _amount) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackBurn(resourceId [32]byte, from common.Address, amount *big.Int) []byte {
	enc, err := enygmaPNEvents.abi.Pack("burn", resourceId, from, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCancelSwap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x22647cd9.
//
// Solidity: function cancelSwap(bytes32 _sharedId, uint256 _toChainId, bytes32 _tokenInResourceId, uint256 _tokenInAmount, uint256 _tokenInId, uint8 _tokenInStandard, bytes32 _tokenOutResourceId, uint256 _tokenOutAmount, uint256 _tokenOutId, uint8 _tokenOutStandard) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackCancelSwap(sharedId [32]byte, toChainId *big.Int, tokenInResourceId [32]byte, tokenInAmount *big.Int, tokenInId *big.Int, tokenInStandard uint8, tokenOutResourceId [32]byte, tokenOutAmount *big.Int, tokenOutId *big.Int, tokenOutStandard uint8) []byte {
	enc, err := enygmaPNEvents.abi.Pack("cancelSwap", sharedId, toChainId, tokenInResourceId, tokenInAmount, tokenInId, tokenInStandard, tokenOutResourceId, tokenOutAmount, tokenOutId, tokenOutStandard)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCreation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x56249c36.
//
// Solidity: function creation(bytes32 _resourceId, uint256 initialSupply) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackCreation(resourceId [32]byte, initialSupply *big.Int) []byte {
	enc, err := enygmaPNEvents.abi.Pack("creation", resourceId, initialSupply)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDepositToDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x68f797c9.
//
// Solidity: function depositToDvp(bytes32 _resourceId, uint256 amount, address _from, bytes32 _referenceId) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackDepositToDvp(resourceId [32]byte, amount *big.Int, from common.Address, referenceId [32]byte) []byte {
	enc, err := enygmaPNEvents.abi.Pack("depositToDvp", resourceId, amount, from, referenceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDvp1155Burn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7441a427.
//
// Solidity: function dvp1155Burn(bytes32 _resourceId, address _to, uint256 _tokenId, uint256 _value) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackDvp1155Burn(resourceId [32]byte, to common.Address, tokenId *big.Int, value *big.Int) []byte {
	enc, err := enygmaPNEvents.abi.Pack("dvp1155Burn", resourceId, to, tokenId, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDvp1155Creation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x05a3b856.
//
// Solidity: function dvp1155Creation(bytes32 _resourceId) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackDvp1155Creation(resourceId [32]byte) []byte {
	enc, err := enygmaPNEvents.abi.Pack("dvp1155Creation", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDvp1155DepositIntoDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1e8fcfef.
//
// Solidity: function dvp1155DepositIntoDvp(bytes32 _resourceId, uint256 _tokenId, address from, uint256 _value, bytes data) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackDvp1155DepositIntoDvp(resourceId [32]byte, tokenId *big.Int, from common.Address, value *big.Int, data []byte) []byte {
	enc, err := enygmaPNEvents.abi.Pack("dvp1155DepositIntoDvp", resourceId, tokenId, from, value, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDvp1155Mint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3405432c.
//
// Solidity: function dvp1155Mint(bytes32 _resourceId, uint256 _tokenId, uint256 _value, bytes data) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackDvp1155Mint(resourceId [32]byte, tokenId *big.Int, value *big.Int, data []byte) []byte {
	enc, err := enygmaPNEvents.abi.Pack("dvp1155Mint", resourceId, tokenId, value, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDvp1155SwapCompleted is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1a8997f7.
//
// Solidity: function dvp1155SwapCompleted(bytes32 _resourceId, uint256 _tokenId, uint256 _destinationChainId, address _destinationOwner) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackDvp1155SwapCompleted(resourceId [32]byte, tokenId *big.Int, destinationChainId *big.Int, destinationOwner common.Address) []byte {
	enc, err := enygmaPNEvents.abi.Pack("dvp1155SwapCompleted", resourceId, tokenId, destinationChainId, destinationOwner)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDvp1155SwapForEnygma is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3d75fb4c.
//
// Solidity: function dvp1155SwapForEnygma(bytes32 _tokenResourceId, uint256 _tokenId, uint256 _tokenValue, bytes _tokenData, uint256 _enygmaAmount, bytes32 _enygmaResourceId, address from, uint256 _destChainId, bytes32 _sharedId, uint64 _validityTime) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackDvp1155SwapForEnygma(tokenResourceId [32]byte, tokenId *big.Int, tokenValue *big.Int, tokenData []byte, enygmaAmount *big.Int, enygmaResourceId [32]byte, from common.Address, destChainId *big.Int, sharedId [32]byte, validityTime uint64) []byte {
	enc, err := enygmaPNEvents.abi.Pack("dvp1155SwapForEnygma", tokenResourceId, tokenId, tokenValue, tokenData, enygmaAmount, enygmaResourceId, from, destChainId, sharedId, validityTime)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDvp1155WithdrawFromDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0fbfc661.
//
// Solidity: function dvp1155WithdrawFromDvp(bytes32 _resourceId, uint256 _tokenId, uint256 _value, bytes data, address owner) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackDvp1155WithdrawFromDvp(resourceId [32]byte, tokenId *big.Int, value *big.Int, data []byte, owner common.Address) []byte {
	enc, err := enygmaPNEvents.abi.Pack("dvp1155WithdrawFromDvp", resourceId, tokenId, value, data, owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDvp721Burn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2cbc8600.
//
// Solidity: function dvp721Burn(bytes32 _resourceId, uint256 _nftId) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackDvp721Burn(resourceId [32]byte, nftId *big.Int) []byte {
	enc, err := enygmaPNEvents.abi.Pack("dvp721Burn", resourceId, nftId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDvp721Creation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x774dcd89.
//
// Solidity: function dvp721Creation(bytes32 _resourceId) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackDvp721Creation(resourceId [32]byte) []byte {
	enc, err := enygmaPNEvents.abi.Pack("dvp721Creation", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDvp721DepositIntoDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x02136b01.
//
// Solidity: function dvp721DepositIntoDvp(bytes32 _resourceId, uint256 _nftId, address from) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackDvp721DepositIntoDvp(resourceId [32]byte, nftId *big.Int, from common.Address) []byte {
	enc, err := enygmaPNEvents.abi.Pack("dvp721DepositIntoDvp", resourceId, nftId, from)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDvp721Mint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6b41b261.
//
// Solidity: function dvp721Mint(bytes32 _resourceId, uint256 _nftId) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackDvp721Mint(resourceId [32]byte, nftId *big.Int) []byte {
	enc, err := enygmaPNEvents.abi.Pack("dvp721Mint", resourceId, nftId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDvp721SwapCompleted is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x44dc1ce1.
//
// Solidity: function dvp721SwapCompleted(bytes32 _resourceId, uint256 _nftId, uint256 _destinationChainId, address _destinationOwner) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackDvp721SwapCompleted(resourceId [32]byte, nftId *big.Int, destinationChainId *big.Int, destinationOwner common.Address) []byte {
	enc, err := enygmaPNEvents.abi.Pack("dvp721SwapCompleted", resourceId, nftId, destinationChainId, destinationOwner)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDvp721SwapForEnygma is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf205b870.
//
// Solidity: function dvp721SwapForEnygma(bytes32 _nftResourceId, uint256 _nftId, uint256 _enygmaAmount, bytes32 _enygmaResourceId, address _from, uint256 _destChainId, bytes32 _sharedId, uint64 _validityTime) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackDvp721SwapForEnygma(nftResourceId [32]byte, nftId *big.Int, enygmaAmount *big.Int, enygmaResourceId [32]byte, from common.Address, destChainId *big.Int, sharedId [32]byte, validityTime uint64) []byte {
	enc, err := enygmaPNEvents.abi.Pack("dvp721SwapForEnygma", nftResourceId, nftId, enygmaAmount, enygmaResourceId, from, destChainId, sharedId, validityTime)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDvp721WithdrawFromDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3a06a3f4.
//
// Solidity: function dvp721WithdrawFromDvp(bytes32 _resourceId, uint256 _nftId, address owner) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackDvp721WithdrawFromDvp(resourceId [32]byte, nftId *big.Int, owner common.Address) []byte {
	enc, err := enygmaPNEvents.abi.Pack("dvp721WithdrawFromDvp", resourceId, nftId, owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackGetAddressByResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (enygmaPNEvents *EnygmaPNEvents) PackGetAddressByResourceId(resourceId [32]byte) []byte {
	enc, err := enygmaPNEvents.abi.Pack("getAddressByResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAddressByResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (enygmaPNEvents *EnygmaPNEvents) UnpackGetAddressByResourceId(data []byte) (common.Address, error) {
	out, err := enygmaPNEvents.abi.Unpack("getAddressByResourceId", data)
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
func (enygmaPNEvents *EnygmaPNEvents) PackGetEndpointAddress() []byte {
	enc, err := enygmaPNEvents.abi.Pack("getEndpointAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetEndpointAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xce884eb5.
//
// Solidity: function getEndpointAddress() view returns(address)
func (enygmaPNEvents *EnygmaPNEvents) UnpackGetEndpointAddress(data []byte) (common.Address, error) {
	out, err := enygmaPNEvents.abi.Unpack("getEndpointAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (enygmaPNEvents *EnygmaPNEvents) PackInitialize() []byte {
	enc, err := enygmaPNEvents.abi.Pack("initialize")
	if err != nil {
		panic(err)
	}
	return enc
}

// PackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7ed9db59.
//
// Solidity: function mint(bytes32 _resourceId, address _to, uint256 _amount) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackMint(resourceId [32]byte, to common.Address, amount *big.Int) []byte {
	enc, err := enygmaPNEvents.abi.Pack("mint", resourceId, to, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackParticipantValidator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6f5be509.
//
// Solidity: function participantValidator() view returns(address)
func (enygmaPNEvents *EnygmaPNEvents) PackParticipantValidator() []byte {
	enc, err := enygmaPNEvents.abi.Pack("participantValidator")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackParticipantValidator is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6f5be509.
//
// Solidity: function participantValidator() view returns(address)
func (enygmaPNEvents *EnygmaPNEvents) UnpackParticipantValidator(data []byte) (common.Address, error) {
	out, err := enygmaPNEvents.abi.Unpack("participantValidator", data)
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
func (enygmaPNEvents *EnygmaPNEvents) PackRaylsNodeUserGovernance() []byte {
	enc, err := enygmaPNEvents.abi.Pack("raylsNodeUserGovernance")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackRaylsNodeUserGovernance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0274c133.
//
// Solidity: function raylsNodeUserGovernance() view returns(address)
func (enygmaPNEvents *EnygmaPNEvents) UnpackRaylsNodeUserGovernance(data []byte) (common.Address, error) {
	out, err := enygmaPNEvents.abi.Unpack("raylsNodeUserGovernance", data)
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
func (enygmaPNEvents *EnygmaPNEvents) PackResourceId() []byte {
	enc, err := enygmaPNEvents.abi.Pack("resourceId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (enygmaPNEvents *EnygmaPNEvents) UnpackResourceId(data []byte) ([32]byte, error) {
	out, err := enygmaPNEvents.abi.Unpack("resourceId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackRevertMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3caf0476.
//
// Solidity: function revertMint(bytes32 _resourceId, uint256 _amount, address _to, string _reason) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackRevertMint(resourceId [32]byte, amount *big.Int, to common.Address, reason string) []byte {
	enc, err := enygmaPNEvents.abi.Pack("revertMint", resourceId, amount, to, reason)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSendTransferPNH is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3a048283.
//
// Solidity: function sendTransferPNH((bytes32,uint256[],uint256[],address[],address,bytes32,(bytes32,address,bytes4,bytes)[][]) _pnhTransfer) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackSendTransferPNH(pnhTransfer SharedObjectsPNHTransfer) []byte {
	enc, err := enygmaPNEvents.abi.Pack("sendTransferPNH", pnhTransfer)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetParticipantValidator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd1f2ada1.
//
// Solidity: function setParticipantValidator(address _participantValidator) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackSetParticipantValidator(participantValidator common.Address) []byte {
	enc, err := enygmaPNEvents.abi.Pack("setParticipantValidator", participantValidator)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa01afbfb.
//
// Solidity: function setResourceId(bytes32 _resourceId) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackSetResourceId(resourceId [32]byte) []byte {
	enc, err := enygmaPNEvents.abi.Pack("setResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetTokenValidator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x268f2e41.
//
// Solidity: function setTokenValidator(address _tokenValidator) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackSetTokenValidator(tokenValidator common.Address) []byte {
	enc, err := enygmaPNEvents.abi.Pack("setTokenValidator", tokenValidator)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSwapWithDvpForERC1155 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbb9c53ae.
//
// Solidity: function swapWithDvpForERC1155(bytes32 _resourceId, uint256 _nftId, bytes32 _nftResourceId, uint256 _nftAmountOrOne, uint256 _enygmaAmount, address _from, uint256 _destChainId, bytes32 _sharedId, uint64 _validityTime) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackSwapWithDvpForERC1155(resourceId [32]byte, nftId *big.Int, nftResourceId [32]byte, nftAmountOrOne *big.Int, enygmaAmount *big.Int, from common.Address, destChainId *big.Int, sharedId [32]byte, validityTime uint64) []byte {
	enc, err := enygmaPNEvents.abi.Pack("swapWithDvpForERC1155", resourceId, nftId, nftResourceId, nftAmountOrOne, enygmaAmount, from, destChainId, sharedId, validityTime)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSwapWithDvpForERC721 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xea50f46b.
//
// Solidity: function swapWithDvpForERC721(bytes32 _resourceId, uint256 _nftId, bytes32 _nftResourceId, uint256 _enygmaAmount, address _from, uint256 _destChainId, bytes32 _sharedId, uint64 _validityTime) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackSwapWithDvpForERC721(resourceId [32]byte, nftId *big.Int, nftResourceId [32]byte, enygmaAmount *big.Int, from common.Address, destChainId *big.Int, sharedId [32]byte, validityTime uint64) []byte {
	enc, err := enygmaPNEvents.abi.Pack("swapWithDvpForERC721", resourceId, nftId, nftResourceId, enygmaAmount, from, destChainId, sharedId, validityTime)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackTokenValidator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x80691fba.
//
// Solidity: function tokenValidator() view returns(address)
func (enygmaPNEvents *EnygmaPNEvents) PackTokenValidator() []byte {
	enc, err := enygmaPNEvents.abi.Pack("tokenValidator")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTokenValidator is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x80691fba.
//
// Solidity: function tokenValidator() view returns(address)
func (enygmaPNEvents *EnygmaPNEvents) UnpackTokenValidator(data []byte) (common.Address, error) {
	out, err := enygmaPNEvents.abi.Unpack("tokenValidator", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackWithdrawFromDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0bcca0d2.
//
// Solidity: function withdrawFromDvp(bytes32 _resourceId, uint256 amount, address _to, bytes32 _referenceId) returns()
func (enygmaPNEvents *EnygmaPNEvents) PackWithdrawFromDvp(resourceId [32]byte, amount *big.Int, to common.Address, referenceId [32]byte) []byte {
	enc, err := enygmaPNEvents.abi.Pack("withdrawFromDvp", resourceId, amount, to, referenceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// EnygmaPNEventsAuthorityUpdated represents a AuthorityUpdated event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsAuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsAuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsAuthorityUpdated) ContractEventName() string {
	return EnygmaPNEventsAuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (enygmaPNEvents *EnygmaPNEvents) UnpackAuthorityUpdatedEvent(log *types.Log) (*EnygmaPNEventsAuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsAuthorityUpdated)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsDvp1155Burn represents a Dvp1155Burn event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsDvp1155Burn struct {
	ResourceId [32]byte
	To         common.Address
	TokenId    *big.Int
	Value      *big.Int
	Raw        *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsDvp1155BurnEventName = "Dvp1155Burn"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsDvp1155Burn) ContractEventName() string {
	return EnygmaPNEventsDvp1155BurnEventName
}

// UnpackDvp1155BurnEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Dvp1155Burn(bytes32 _resourceId, address _to, uint256 _tokenId, uint256 _value)
func (enygmaPNEvents *EnygmaPNEvents) UnpackDvp1155BurnEvent(log *types.Log) (*EnygmaPNEventsDvp1155Burn, error) {
	event := "Dvp1155Burn"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsDvp1155Burn)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsDvp1155Creation represents a Dvp1155Creation event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsDvp1155Creation struct {
	ResourceId [32]byte
	Raw        *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsDvp1155CreationEventName = "Dvp1155Creation"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsDvp1155Creation) ContractEventName() string {
	return EnygmaPNEventsDvp1155CreationEventName
}

// UnpackDvp1155CreationEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Dvp1155Creation(bytes32 _resourceId)
func (enygmaPNEvents *EnygmaPNEvents) UnpackDvp1155CreationEvent(log *types.Log) (*EnygmaPNEventsDvp1155Creation, error) {
	event := "Dvp1155Creation"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsDvp1155Creation)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsDvp1155DepositIntoDvp represents a Dvp1155DepositIntoDvp event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsDvp1155DepositIntoDvp struct {
	ResourceId [32]byte
	TokenId    *big.Int
	From       common.Address
	Value      *big.Int
	Data       []byte
	Raw        *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsDvp1155DepositIntoDvpEventName = "Dvp1155DepositIntoDvp"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsDvp1155DepositIntoDvp) ContractEventName() string {
	return EnygmaPNEventsDvp1155DepositIntoDvpEventName
}

// UnpackDvp1155DepositIntoDvpEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Dvp1155DepositIntoDvp(bytes32 _resourceId, uint256 _tokenId, address from, uint256 _value, bytes data)
func (enygmaPNEvents *EnygmaPNEvents) UnpackDvp1155DepositIntoDvpEvent(log *types.Log) (*EnygmaPNEventsDvp1155DepositIntoDvp, error) {
	event := "Dvp1155DepositIntoDvp"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsDvp1155DepositIntoDvp)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsDvp1155Mint represents a Dvp1155Mint event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsDvp1155Mint struct {
	ResourceId [32]byte
	TokenId    *big.Int
	Value      *big.Int
	Data       []byte
	Raw        *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsDvp1155MintEventName = "Dvp1155Mint"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsDvp1155Mint) ContractEventName() string {
	return EnygmaPNEventsDvp1155MintEventName
}

// UnpackDvp1155MintEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Dvp1155Mint(bytes32 _resourceId, uint256 _tokenId, uint256 _value, bytes data)
func (enygmaPNEvents *EnygmaPNEvents) UnpackDvp1155MintEvent(log *types.Log) (*EnygmaPNEventsDvp1155Mint, error) {
	event := "Dvp1155Mint"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsDvp1155Mint)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsDvp1155SwapCompleted represents a Dvp1155SwapCompleted event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsDvp1155SwapCompleted struct {
	ResourceId         [32]byte
	TokenId            *big.Int
	DestinationChainId *big.Int
	DestinationOwner   common.Address
	Raw                *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsDvp1155SwapCompletedEventName = "Dvp1155SwapCompleted"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsDvp1155SwapCompleted) ContractEventName() string {
	return EnygmaPNEventsDvp1155SwapCompletedEventName
}

// UnpackDvp1155SwapCompletedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Dvp1155SwapCompleted(bytes32 _resourceId, uint256 _tokenId, uint256 _destinationChainId, address _destinationOwner)
func (enygmaPNEvents *EnygmaPNEvents) UnpackDvp1155SwapCompletedEvent(log *types.Log) (*EnygmaPNEventsDvp1155SwapCompleted, error) {
	event := "Dvp1155SwapCompleted"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsDvp1155SwapCompleted)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsDvp1155SwapForEnygma represents a Dvp1155SwapForEnygma event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsDvp1155SwapForEnygma struct {
	TokenResourceId  [32]byte
	TokenId          *big.Int
	TokenValue       *big.Int
	TokenData        []byte
	EnygmaAmount     *big.Int
	EnygmaResourceId [32]byte
	From             common.Address
	DestChainId      *big.Int
	SharedId         [32]byte
	ValidityTime     uint64
	Raw              *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsDvp1155SwapForEnygmaEventName = "Dvp1155SwapForEnygma"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsDvp1155SwapForEnygma) ContractEventName() string {
	return EnygmaPNEventsDvp1155SwapForEnygmaEventName
}

// UnpackDvp1155SwapForEnygmaEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Dvp1155SwapForEnygma(bytes32 _tokenResourceId, uint256 _tokenId, uint256 _tokenValue, bytes _tokenData, uint256 _enygmaAmount, bytes32 _enygmaResourceId, address from, uint256 _destChainId, bytes32 _sharedId, uint64 _validityTime)
func (enygmaPNEvents *EnygmaPNEvents) UnpackDvp1155SwapForEnygmaEvent(log *types.Log) (*EnygmaPNEventsDvp1155SwapForEnygma, error) {
	event := "Dvp1155SwapForEnygma"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsDvp1155SwapForEnygma)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsDvp1155WithdrawFromDvp represents a Dvp1155WithdrawFromDvp event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsDvp1155WithdrawFromDvp struct {
	ResourceId [32]byte
	TokenId    *big.Int
	Value      *big.Int
	Data       []byte
	Owner      common.Address
	Raw        *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsDvp1155WithdrawFromDvpEventName = "Dvp1155WithdrawFromDvp"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsDvp1155WithdrawFromDvp) ContractEventName() string {
	return EnygmaPNEventsDvp1155WithdrawFromDvpEventName
}

// UnpackDvp1155WithdrawFromDvpEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Dvp1155WithdrawFromDvp(bytes32 _resourceId, uint256 _tokenId, uint256 _value, bytes data, address owner)
func (enygmaPNEvents *EnygmaPNEvents) UnpackDvp1155WithdrawFromDvpEvent(log *types.Log) (*EnygmaPNEventsDvp1155WithdrawFromDvp, error) {
	event := "Dvp1155WithdrawFromDvp"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsDvp1155WithdrawFromDvp)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsDvp721Burn represents a Dvp721Burn event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsDvp721Burn struct {
	ResourceId [32]byte
	NftId      *big.Int
	Raw        *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsDvp721BurnEventName = "Dvp721Burn"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsDvp721Burn) ContractEventName() string {
	return EnygmaPNEventsDvp721BurnEventName
}

// UnpackDvp721BurnEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Dvp721Burn(bytes32 _resourceId, uint256 _nftId)
func (enygmaPNEvents *EnygmaPNEvents) UnpackDvp721BurnEvent(log *types.Log) (*EnygmaPNEventsDvp721Burn, error) {
	event := "Dvp721Burn"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsDvp721Burn)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsDvp721Creation represents a Dvp721Creation event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsDvp721Creation struct {
	ResourceId [32]byte
	Raw        *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsDvp721CreationEventName = "Dvp721Creation"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsDvp721Creation) ContractEventName() string {
	return EnygmaPNEventsDvp721CreationEventName
}

// UnpackDvp721CreationEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Dvp721Creation(bytes32 _resourceId)
func (enygmaPNEvents *EnygmaPNEvents) UnpackDvp721CreationEvent(log *types.Log) (*EnygmaPNEventsDvp721Creation, error) {
	event := "Dvp721Creation"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsDvp721Creation)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsDvp721DepositIntoDvp represents a Dvp721DepositIntoDvp event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsDvp721DepositIntoDvp struct {
	ResourceId [32]byte
	NftId      *big.Int
	From       common.Address
	Raw        *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsDvp721DepositIntoDvpEventName = "Dvp721DepositIntoDvp"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsDvp721DepositIntoDvp) ContractEventName() string {
	return EnygmaPNEventsDvp721DepositIntoDvpEventName
}

// UnpackDvp721DepositIntoDvpEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Dvp721DepositIntoDvp(bytes32 _resourceId, uint256 _nftId, address from)
func (enygmaPNEvents *EnygmaPNEvents) UnpackDvp721DepositIntoDvpEvent(log *types.Log) (*EnygmaPNEventsDvp721DepositIntoDvp, error) {
	event := "Dvp721DepositIntoDvp"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsDvp721DepositIntoDvp)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsDvp721Mint represents a Dvp721Mint event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsDvp721Mint struct {
	ResourceId [32]byte
	NftId      *big.Int
	Raw        *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsDvp721MintEventName = "Dvp721Mint"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsDvp721Mint) ContractEventName() string {
	return EnygmaPNEventsDvp721MintEventName
}

// UnpackDvp721MintEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Dvp721Mint(bytes32 _resourceId, uint256 _nftId)
func (enygmaPNEvents *EnygmaPNEvents) UnpackDvp721MintEvent(log *types.Log) (*EnygmaPNEventsDvp721Mint, error) {
	event := "Dvp721Mint"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsDvp721Mint)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsDvp721SwapCompleted represents a Dvp721SwapCompleted event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsDvp721SwapCompleted struct {
	ResourceId         [32]byte
	NftId              *big.Int
	DestinationChainId *big.Int
	DestinationOwner   common.Address
	Raw                *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsDvp721SwapCompletedEventName = "Dvp721SwapCompleted"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsDvp721SwapCompleted) ContractEventName() string {
	return EnygmaPNEventsDvp721SwapCompletedEventName
}

// UnpackDvp721SwapCompletedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Dvp721SwapCompleted(bytes32 _resourceId, uint256 _nftId, uint256 _destinationChainId, address _destinationOwner)
func (enygmaPNEvents *EnygmaPNEvents) UnpackDvp721SwapCompletedEvent(log *types.Log) (*EnygmaPNEventsDvp721SwapCompleted, error) {
	event := "Dvp721SwapCompleted"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsDvp721SwapCompleted)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsDvp721SwapForEnygma represents a Dvp721SwapForEnygma event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsDvp721SwapForEnygma struct {
	NftResourceId    [32]byte
	NftId            *big.Int
	EnygmaAmount     *big.Int
	EnygmaResourceId [32]byte
	From             common.Address
	DestChainId      *big.Int
	SharedId         [32]byte
	ValidityTime     uint64
	Raw              *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsDvp721SwapForEnygmaEventName = "Dvp721SwapForEnygma"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsDvp721SwapForEnygma) ContractEventName() string {
	return EnygmaPNEventsDvp721SwapForEnygmaEventName
}

// UnpackDvp721SwapForEnygmaEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Dvp721SwapForEnygma(bytes32 _nftResourceId, uint256 _nftId, uint256 _enygmaAmount, bytes32 _enygmaResourceId, address _from, uint256 _destChainId, bytes32 _sharedId, uint64 _validityTime)
func (enygmaPNEvents *EnygmaPNEvents) UnpackDvp721SwapForEnygmaEvent(log *types.Log) (*EnygmaPNEventsDvp721SwapForEnygma, error) {
	event := "Dvp721SwapForEnygma"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsDvp721SwapForEnygma)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsDvp721WithdrawFromDvp represents a Dvp721WithdrawFromDvp event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsDvp721WithdrawFromDvp struct {
	ResourceId [32]byte
	NftId      *big.Int
	Owner      common.Address
	Raw        *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsDvp721WithdrawFromDvpEventName = "Dvp721WithdrawFromDvp"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsDvp721WithdrawFromDvp) ContractEventName() string {
	return EnygmaPNEventsDvp721WithdrawFromDvpEventName
}

// UnpackDvp721WithdrawFromDvpEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Dvp721WithdrawFromDvp(bytes32 _resourceId, uint256 _nftId, address owner)
func (enygmaPNEvents *EnygmaPNEvents) UnpackDvp721WithdrawFromDvpEvent(log *types.Log) (*EnygmaPNEventsDvp721WithdrawFromDvp, error) {
	event := "Dvp721WithdrawFromDvp"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsDvp721WithdrawFromDvp)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsDvpSwapCancelled represents a DvpSwapCancelled event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsDvpSwapCancelled struct {
	SharedId           [32]byte
	DestChainId        *big.Int
	TokenInResourceId  [32]byte
	TokenInAmount      *big.Int
	TokenInId          *big.Int
	TokenInStandard    uint8
	TokenOutResourceId [32]byte
	TokenOutAmount     *big.Int
	TokenOutId         *big.Int
	TokenOutStandard   uint8
	Raw                *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsDvpSwapCancelledEventName = "DvpSwapCancelled"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsDvpSwapCancelled) ContractEventName() string {
	return EnygmaPNEventsDvpSwapCancelledEventName
}

// UnpackDvpSwapCancelledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DvpSwapCancelled(bytes32 _sharedId, uint256 _destChainId, bytes32 _tokenInResourceId, uint256 _tokenInAmount, uint256 _tokenInId, uint8 _tokenInStandard, bytes32 _tokenOutResourceId, uint256 _tokenOutAmount, uint256 _tokenOutId, uint8 _tokenOutStandard)
func (enygmaPNEvents *EnygmaPNEvents) UnpackDvpSwapCancelledEvent(log *types.Log) (*EnygmaPNEventsDvpSwapCancelled, error) {
	event := "DvpSwapCancelled"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsDvpSwapCancelled)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsEnygmaBurn represents a EnygmaBurn event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsEnygmaBurn struct {
	ResourceId [32]byte
	From       common.Address
	Amount     *big.Int
	Raw        *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsEnygmaBurnEventName = "EnygmaBurn"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsEnygmaBurn) ContractEventName() string {
	return EnygmaPNEventsEnygmaBurnEventName
}

// UnpackEnygmaBurnEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EnygmaBurn(bytes32 _resourceId, address _from, uint256 _amount)
func (enygmaPNEvents *EnygmaPNEvents) UnpackEnygmaBurnEvent(log *types.Log) (*EnygmaPNEventsEnygmaBurn, error) {
	event := "EnygmaBurn"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsEnygmaBurn)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsEnygmaCreation represents a EnygmaCreation event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsEnygmaCreation struct {
	ResourceId    [32]byte
	InitialSupply *big.Int
	Raw           *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsEnygmaCreationEventName = "EnygmaCreation"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsEnygmaCreation) ContractEventName() string {
	return EnygmaPNEventsEnygmaCreationEventName
}

// UnpackEnygmaCreationEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EnygmaCreation(bytes32 _resourceId, uint256 _initialSupply)
func (enygmaPNEvents *EnygmaPNEvents) UnpackEnygmaCreationEvent(log *types.Log) (*EnygmaPNEventsEnygmaCreation, error) {
	event := "EnygmaCreation"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsEnygmaCreation)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsEnygmaDepositToDvp represents a EnygmaDepositToDvp event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsEnygmaDepositToDvp struct {
	ResourceId  [32]byte
	Amount      *big.Int
	From        common.Address
	ReferenceId [32]byte
	Raw         *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsEnygmaDepositToDvpEventName = "EnygmaDepositToDvp"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsEnygmaDepositToDvp) ContractEventName() string {
	return EnygmaPNEventsEnygmaDepositToDvpEventName
}

// UnpackEnygmaDepositToDvpEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EnygmaDepositToDvp(bytes32 _resourceId, uint256 amount, address _from, bytes32 _referenceId)
func (enygmaPNEvents *EnygmaPNEvents) UnpackEnygmaDepositToDvpEvent(log *types.Log) (*EnygmaPNEventsEnygmaDepositToDvp, error) {
	event := "EnygmaDepositToDvp"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsEnygmaDepositToDvp)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsEnygmaMint represents a EnygmaMint event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsEnygmaMint struct {
	ResourceId [32]byte
	To         common.Address
	Amount     *big.Int
	Raw        *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsEnygmaMintEventName = "EnygmaMint"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsEnygmaMint) ContractEventName() string {
	return EnygmaPNEventsEnygmaMintEventName
}

// UnpackEnygmaMintEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EnygmaMint(bytes32 _resourceId, address _to, uint256 _amount)
func (enygmaPNEvents *EnygmaPNEvents) UnpackEnygmaMintEvent(log *types.Log) (*EnygmaPNEventsEnygmaMint, error) {
	event := "EnygmaMint"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsEnygmaMint)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsEnygmaRevertMint represents a EnygmaRevertMint event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsEnygmaRevertMint struct {
	ResourceId [32]byte
	Amount     *big.Int
	To         common.Address
	Reason     string
	Raw        *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsEnygmaRevertMintEventName = "EnygmaRevertMint"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsEnygmaRevertMint) ContractEventName() string {
	return EnygmaPNEventsEnygmaRevertMintEventName
}

// UnpackEnygmaRevertMintEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EnygmaRevertMint(bytes32 _resourceId, uint256 _amount, address _to, string _reason)
func (enygmaPNEvents *EnygmaPNEvents) UnpackEnygmaRevertMintEvent(log *types.Log) (*EnygmaPNEventsEnygmaRevertMint, error) {
	event := "EnygmaRevertMint"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsEnygmaRevertMint)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsEnygmaSendTransferPNH represents a EnygmaSendTransferPNH event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsEnygmaSendTransferPNH struct {
	ResourceId  [32]byte
	Value       []*big.Int
	ToChainId   []*big.Int
	To          []common.Address
	From        common.Address
	ReferenceId [32]byte
	ProgramData [][]SharedObjectsEnygmaProgramData
	Raw         *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsEnygmaSendTransferPNHEventName = "EnygmaSendTransferPNH"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsEnygmaSendTransferPNH) ContractEventName() string {
	return EnygmaPNEventsEnygmaSendTransferPNHEventName
}

// UnpackEnygmaSendTransferPNHEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EnygmaSendTransferPNH(bytes32 _resourceId, uint256[] _value, uint256[] _toChainId, address[] _to, address _from, bytes32 _referenceId, (bytes32,address,bytes4,bytes)[][] _programData)
func (enygmaPNEvents *EnygmaPNEvents) UnpackEnygmaSendTransferPNHEvent(log *types.Log) (*EnygmaPNEventsEnygmaSendTransferPNH, error) {
	event := "EnygmaSendTransferPNH"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsEnygmaSendTransferPNH)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsEnygmaSwapWithDvpForERC1155 represents a EnygmaSwapWithDvpForERC1155 event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsEnygmaSwapWithDvpForERC1155 struct {
	ResourceId     [32]byte
	NftId          *big.Int
	NftResourceId  [32]byte
	NftAmountOrOne *big.Int
	EnygmaAmount   *big.Int
	From           common.Address
	DestChainId    *big.Int
	SharedId       [32]byte
	ValidityTime   uint64
	Raw            *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsEnygmaSwapWithDvpForERC1155EventName = "EnygmaSwapWithDvpForERC1155"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsEnygmaSwapWithDvpForERC1155) ContractEventName() string {
	return EnygmaPNEventsEnygmaSwapWithDvpForERC1155EventName
}

// UnpackEnygmaSwapWithDvpForERC1155Event is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EnygmaSwapWithDvpForERC1155(bytes32 _resourceId, uint256 _nftId, bytes32 _nftResourceId, uint256 _nftAmountOrOne, uint256 _enygmaAmount, address _from, uint256 _destChainId, bytes32 _sharedId, uint64 _validityTime)
func (enygmaPNEvents *EnygmaPNEvents) UnpackEnygmaSwapWithDvpForERC1155Event(log *types.Log) (*EnygmaPNEventsEnygmaSwapWithDvpForERC1155, error) {
	event := "EnygmaSwapWithDvpForERC1155"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsEnygmaSwapWithDvpForERC1155)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsEnygmaSwapWithDvpForERC721 represents a EnygmaSwapWithDvpForERC721 event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsEnygmaSwapWithDvpForERC721 struct {
	ResourceId    [32]byte
	NftId         *big.Int
	NftResourceId [32]byte
	EnygmaAmount  *big.Int
	From          common.Address
	DestChainId   *big.Int
	SharedId      [32]byte
	ValidityTime  uint64
	Raw           *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsEnygmaSwapWithDvpForERC721EventName = "EnygmaSwapWithDvpForERC721"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsEnygmaSwapWithDvpForERC721) ContractEventName() string {
	return EnygmaPNEventsEnygmaSwapWithDvpForERC721EventName
}

// UnpackEnygmaSwapWithDvpForERC721Event is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EnygmaSwapWithDvpForERC721(bytes32 _resourceId, uint256 _nftId, bytes32 _nftResourceId, uint256 _enygmaAmount, address _from, uint256 _destChainId, bytes32 _sharedId, uint64 _validityTime)
func (enygmaPNEvents *EnygmaPNEvents) UnpackEnygmaSwapWithDvpForERC721Event(log *types.Log) (*EnygmaPNEventsEnygmaSwapWithDvpForERC721, error) {
	event := "EnygmaSwapWithDvpForERC721"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsEnygmaSwapWithDvpForERC721)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsEnygmaWithdrawFromDvp represents a EnygmaWithdrawFromDvp event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsEnygmaWithdrawFromDvp struct {
	ResourceId  [32]byte
	Amount      *big.Int
	To          common.Address
	ReferenceId [32]byte
	Raw         *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsEnygmaWithdrawFromDvpEventName = "EnygmaWithdrawFromDvp"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsEnygmaWithdrawFromDvp) ContractEventName() string {
	return EnygmaPNEventsEnygmaWithdrawFromDvpEventName
}

// UnpackEnygmaWithdrawFromDvpEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EnygmaWithdrawFromDvp(bytes32 _resourceId, uint256 amount, address _to, bytes32 _referenceId)
func (enygmaPNEvents *EnygmaPNEvents) UnpackEnygmaWithdrawFromDvpEvent(log *types.Log) (*EnygmaPNEventsEnygmaWithdrawFromDvp, error) {
	event := "EnygmaWithdrawFromDvp"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsEnygmaWithdrawFromDvp)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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

// EnygmaPNEventsTokenRegistrationSubmitted represents a TokenRegistrationSubmitted event raised by the EnygmaPNEvents contract.
type EnygmaPNEventsTokenRegistrationSubmitted struct {
	TokenAddress common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const EnygmaPNEventsTokenRegistrationSubmittedEventName = "TokenRegistrationSubmitted"

// ContractEventName returns the user-defined event name.
func (EnygmaPNEventsTokenRegistrationSubmitted) ContractEventName() string {
	return EnygmaPNEventsTokenRegistrationSubmittedEventName
}

// UnpackTokenRegistrationSubmittedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenRegistrationSubmitted(address indexed tokenAddress)
func (enygmaPNEvents *EnygmaPNEvents) UnpackTokenRegistrationSubmittedEvent(log *types.Log) (*EnygmaPNEventsTokenRegistrationSubmitted, error) {
	event := "TokenRegistrationSubmitted"
	if log.Topics[0] != enygmaPNEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNEventsTokenRegistrationSubmitted)
	if len(log.Data) > 0 {
		if err := enygmaPNEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNEvents.abi.Events[event].Inputs {
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
func (enygmaPNEvents *EnygmaPNEvents) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], enygmaPNEvents.abi.Errors["EnygmaPNEventsValidatorsNotInitialized"].ID.Bytes()[:4]) {
		return enygmaPNEvents.UnpackEnygmaPNEventsValidatorsNotInitializedError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaPNEvents.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return enygmaPNEvents.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaPNEvents.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return enygmaPNEvents.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaPNEvents.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return enygmaPNEvents.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaPNEvents.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return enygmaPNEvents.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaPNEvents.abi.Errors["RaylsAppHubNotActive"].ID.Bytes()[:4]) {
		return enygmaPNEvents.UnpackRaylsAppHubNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaPNEvents.abi.Errors["RaylsAppPrivacyNodeFrozen"].ID.Bytes()[:4]) {
		return enygmaPNEvents.UnpackRaylsAppPrivacyNodeFrozenError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaPNEvents.abi.Errors["RaylsAppPrivacyNodeNotActive"].ID.Bytes()[:4]) {
		return enygmaPNEvents.UnpackRaylsAppPrivacyNodeNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaPNEvents.abi.Errors["RaylsAppPublicChainNotActive"].ID.Bytes()[:4]) {
		return enygmaPNEvents.UnpackRaylsAppPublicChainNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaPNEvents.abi.Errors["RaylsAppResourceNotApproved"].ID.Bytes()[:4]) {
		return enygmaPNEvents.UnpackRaylsAppResourceNotApprovedError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaPNEvents.abi.Errors["RaylsAppTokenRegistryNotConfigured"].ID.Bytes()[:4]) {
		return enygmaPNEvents.UnpackRaylsAppTokenRegistryNotConfiguredError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaPNEvents.abi.Errors["RaylsAppUnauthorizedTokenRegistry"].ID.Bytes()[:4]) {
		return enygmaPNEvents.UnpackRaylsAppUnauthorizedTokenRegistryError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaPNEvents.abi.Errors["RaylsAppUserNotRegistered"].ID.Bytes()[:4]) {
		return enygmaPNEvents.UnpackRaylsAppUserNotRegisteredError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// EnygmaPNEventsEnygmaPNEventsValidatorsNotInitialized represents a EnygmaPNEvents__ValidatorsNotInitialized error raised by the EnygmaPNEvents contract.
type EnygmaPNEventsEnygmaPNEventsValidatorsNotInitialized struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EnygmaPNEvents__ValidatorsNotInitialized()
func EnygmaPNEventsEnygmaPNEventsValidatorsNotInitializedErrorID() common.Hash {
	return common.HexToHash("0x270aa4460f75f4616fce44fd95afd22f0d316c9175c658b162dce55c71e49681")
}

// UnpackEnygmaPNEventsValidatorsNotInitializedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EnygmaPNEvents__ValidatorsNotInitialized()
func (enygmaPNEvents *EnygmaPNEvents) UnpackEnygmaPNEventsValidatorsNotInitializedError(raw []byte) (*EnygmaPNEventsEnygmaPNEventsValidatorsNotInitialized, error) {
	out := new(EnygmaPNEventsEnygmaPNEventsValidatorsNotInitialized)
	if err := enygmaPNEvents.abi.UnpackIntoInterface(out, "EnygmaPNEventsValidatorsNotInitialized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaPNEventsRaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the EnygmaPNEvents contract.
type EnygmaPNEventsRaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func EnygmaPNEventsRaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (enygmaPNEvents *EnygmaPNEvents) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*EnygmaPNEventsRaylsAccessManagedContractPaused, error) {
	out := new(EnygmaPNEventsRaylsAccessManagedContractPaused)
	if err := enygmaPNEvents.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaPNEventsRaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the EnygmaPNEvents contract.
type EnygmaPNEventsRaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func EnygmaPNEventsRaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (enygmaPNEvents *EnygmaPNEvents) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*EnygmaPNEventsRaylsAccessManagedInvalidAuthority, error) {
	out := new(EnygmaPNEventsRaylsAccessManagedInvalidAuthority)
	if err := enygmaPNEvents.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaPNEventsRaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the EnygmaPNEvents contract.
type EnygmaPNEventsRaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func EnygmaPNEventsRaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (enygmaPNEvents *EnygmaPNEvents) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*EnygmaPNEventsRaylsAccessManagedMustSchedule, error) {
	out := new(EnygmaPNEventsRaylsAccessManagedMustSchedule)
	if err := enygmaPNEvents.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaPNEventsRaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the EnygmaPNEvents contract.
type EnygmaPNEventsRaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func EnygmaPNEventsRaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (enygmaPNEvents *EnygmaPNEvents) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*EnygmaPNEventsRaylsAccessManagedUnauthorized, error) {
	out := new(EnygmaPNEventsRaylsAccessManagedUnauthorized)
	if err := enygmaPNEvents.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaPNEventsRaylsAppHubNotActive represents a RaylsApp__HubNotActive error raised by the EnygmaPNEvents contract.
type EnygmaPNEventsRaylsAppHubNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
	HubStatus         uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__HubNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 hubStatus)
func EnygmaPNEventsRaylsAppHubNotActiveErrorID() common.Hash {
	return common.HexToHash("0xdc2ffb0fada912f0dd1b700d4ea9a9ce47e3ecdd1b7b155d2066b9a022a637c2")
}

// UnpackRaylsAppHubNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__HubNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 hubStatus)
func (enygmaPNEvents *EnygmaPNEvents) UnpackRaylsAppHubNotActiveError(raw []byte) (*EnygmaPNEventsRaylsAppHubNotActive, error) {
	out := new(EnygmaPNEventsRaylsAppHubNotActive)
	if err := enygmaPNEvents.abi.UnpackIntoInterface(out, "RaylsAppHubNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaPNEventsRaylsAppPrivacyNodeFrozen represents a RaylsApp__PrivacyNodeFrozen error raised by the EnygmaPNEvents contract.
type EnygmaPNEventsRaylsAppPrivacyNodeFrozen struct {
	TokenAddress common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__PrivacyNodeFrozen(address tokenAddress)
func EnygmaPNEventsRaylsAppPrivacyNodeFrozenErrorID() common.Hash {
	return common.HexToHash("0xcecb8d3ce0d1417038942c9d252e856b5585275082aa5cdbca675fa64d7bfc24")
}

// UnpackRaylsAppPrivacyNodeFrozenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__PrivacyNodeFrozen(address tokenAddress)
func (enygmaPNEvents *EnygmaPNEvents) UnpackRaylsAppPrivacyNodeFrozenError(raw []byte) (*EnygmaPNEventsRaylsAppPrivacyNodeFrozen, error) {
	out := new(EnygmaPNEventsRaylsAppPrivacyNodeFrozen)
	if err := enygmaPNEvents.abi.UnpackIntoInterface(out, "RaylsAppPrivacyNodeFrozen", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaPNEventsRaylsAppPrivacyNodeNotActive represents a RaylsApp__PrivacyNodeNotActive error raised by the EnygmaPNEvents contract.
type EnygmaPNEventsRaylsAppPrivacyNodeNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__PrivacyNodeNotActive(address tokenAddress, uint8 privacyNodeStatus)
func EnygmaPNEventsRaylsAppPrivacyNodeNotActiveErrorID() common.Hash {
	return common.HexToHash("0x44c58c43ed8f726e3330349bec7aa7300f000be36837ee0c2cf507d04511e1e8")
}

// UnpackRaylsAppPrivacyNodeNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__PrivacyNodeNotActive(address tokenAddress, uint8 privacyNodeStatus)
func (enygmaPNEvents *EnygmaPNEvents) UnpackRaylsAppPrivacyNodeNotActiveError(raw []byte) (*EnygmaPNEventsRaylsAppPrivacyNodeNotActive, error) {
	out := new(EnygmaPNEventsRaylsAppPrivacyNodeNotActive)
	if err := enygmaPNEvents.abi.UnpackIntoInterface(out, "RaylsAppPrivacyNodeNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaPNEventsRaylsAppPublicChainNotActive represents a RaylsApp__PublicChainNotActive error raised by the EnygmaPNEvents contract.
type EnygmaPNEventsRaylsAppPublicChainNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
	PublicChainStatus uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__PublicChainNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 publicChainStatus)
func EnygmaPNEventsRaylsAppPublicChainNotActiveErrorID() common.Hash {
	return common.HexToHash("0xd6e23bd403a5000c9afe5c2ed5202b3ff8e25d8c3644c1f51892016fb18e5ab9")
}

// UnpackRaylsAppPublicChainNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__PublicChainNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 publicChainStatus)
func (enygmaPNEvents *EnygmaPNEvents) UnpackRaylsAppPublicChainNotActiveError(raw []byte) (*EnygmaPNEventsRaylsAppPublicChainNotActive, error) {
	out := new(EnygmaPNEventsRaylsAppPublicChainNotActive)
	if err := enygmaPNEvents.abi.UnpackIntoInterface(out, "RaylsAppPublicChainNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaPNEventsRaylsAppResourceNotApproved represents a RaylsApp__ResourceNotApproved error raised by the EnygmaPNEvents contract.
type EnygmaPNEventsRaylsAppResourceNotApproved struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__ResourceNotApproved()
func EnygmaPNEventsRaylsAppResourceNotApprovedErrorID() common.Hash {
	return common.HexToHash("0x970ad4f73c2c200faa068d3d920e2ef40fca6a5338655abcfb5212557edeed6b")
}

// UnpackRaylsAppResourceNotApprovedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__ResourceNotApproved()
func (enygmaPNEvents *EnygmaPNEvents) UnpackRaylsAppResourceNotApprovedError(raw []byte) (*EnygmaPNEventsRaylsAppResourceNotApproved, error) {
	out := new(EnygmaPNEventsRaylsAppResourceNotApproved)
	if err := enygmaPNEvents.abi.UnpackIntoInterface(out, "RaylsAppResourceNotApproved", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaPNEventsRaylsAppTokenRegistryNotConfigured represents a RaylsApp__TokenRegistryNotConfigured error raised by the EnygmaPNEvents contract.
type EnygmaPNEventsRaylsAppTokenRegistryNotConfigured struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__TokenRegistryNotConfigured()
func EnygmaPNEventsRaylsAppTokenRegistryNotConfiguredErrorID() common.Hash {
	return common.HexToHash("0x36a41bd1f6f11cd28b716e935a926fb04f66e11a393b38a49bb660640f3b6dbf")
}

// UnpackRaylsAppTokenRegistryNotConfiguredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__TokenRegistryNotConfigured()
func (enygmaPNEvents *EnygmaPNEvents) UnpackRaylsAppTokenRegistryNotConfiguredError(raw []byte) (*EnygmaPNEventsRaylsAppTokenRegistryNotConfigured, error) {
	out := new(EnygmaPNEventsRaylsAppTokenRegistryNotConfigured)
	if err := enygmaPNEvents.abi.UnpackIntoInterface(out, "RaylsAppTokenRegistryNotConfigured", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaPNEventsRaylsAppUnauthorizedTokenRegistry represents a RaylsApp__UnauthorizedTokenRegistry error raised by the EnygmaPNEvents contract.
type EnygmaPNEventsRaylsAppUnauthorizedTokenRegistry struct {
	Caller        common.Address
	TokenRegistry common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__UnauthorizedTokenRegistry(address caller, address tokenRegistry)
func EnygmaPNEventsRaylsAppUnauthorizedTokenRegistryErrorID() common.Hash {
	return common.HexToHash("0x061526480acdfaa09331b795496a6c50aaed25a45d9fca4c9d55fad56af8e09c")
}

// UnpackRaylsAppUnauthorizedTokenRegistryError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__UnauthorizedTokenRegistry(address caller, address tokenRegistry)
func (enygmaPNEvents *EnygmaPNEvents) UnpackRaylsAppUnauthorizedTokenRegistryError(raw []byte) (*EnygmaPNEventsRaylsAppUnauthorizedTokenRegistry, error) {
	out := new(EnygmaPNEventsRaylsAppUnauthorizedTokenRegistry)
	if err := enygmaPNEvents.abi.UnpackIntoInterface(out, "RaylsAppUnauthorizedTokenRegistry", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaPNEventsRaylsAppUserNotRegistered represents a RaylsApp__UserNotRegistered error raised by the EnygmaPNEvents contract.
type EnygmaPNEventsRaylsAppUserNotRegistered struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__UserNotRegistered(address caller)
func EnygmaPNEventsRaylsAppUserNotRegisteredErrorID() common.Hash {
	return common.HexToHash("0x4c1db902cce08bec31bedc484362fba54949899ac3c0bf0416f3c44af3284baa")
}

// UnpackRaylsAppUserNotRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__UserNotRegistered(address caller)
func (enygmaPNEvents *EnygmaPNEvents) UnpackRaylsAppUserNotRegisteredError(raw []byte) (*EnygmaPNEventsRaylsAppUserNotRegistered, error) {
	out := new(EnygmaPNEventsRaylsAppUserNotRegistered)
	if err := enygmaPNEvents.abi.UnpackIntoInterface(out, "RaylsAppUserNotRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}
