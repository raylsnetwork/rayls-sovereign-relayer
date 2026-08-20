// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package PNTokenRegistryV1

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

// TokenStructsFrozenToken is an auto generated low-level Go binding around an user-defined struct.
type TokenStructsFrozenToken struct {
	ResourceId         [32]byte
	FrozenParticipants []*big.Int
}

// TokenStructsToken is an auto generated low-level Go binding around an user-defined struct.
type TokenStructsToken struct {
	ResourceId         [32]byte
	Name               string
	Symbol             string
	Uri                string
	TokenAddress       common.Address
	PublicTokenAddress common.Address
	IssuerChainId      *big.Int
	ErcStandard        uint8
	PrivacyNodeStatus  uint8
	HubStatus          uint8
	PublicChainStatus  uint8
	CreatedAt          *big.Int
	UpdatedAt          *big.Int
}

// PNTokenRegistryV1MetaData contains all meta data concerning the PNTokenRegistryV1 contract.
var PNTokenRegistryV1MetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"activateToken\",\"inputs\":[{\"name\":\"tokenResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"ercStandard\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"contractVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"deprecateOnPublicChain\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"endpoint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRaylsEndpoint\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"freezeOnPrivacyNode\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"freezeOnPublicChain\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getActiveTokenCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAddressByResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllTokens\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structTokenStructs.Token[]\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"uri\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"publicTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"issuerChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"ercStandard\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.ErcStandard\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\"},{\"name\":\"hubStatus\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.HubStatus\"},{\"name\":\"publicChainStatus\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.PublicChainStatus\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFrozenTokenForParticipant\",\"inputs\":[{\"name\":\"tokenResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getHubStatus\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.HubStatus\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPrivacyNodeStatus\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPublicChainStatus\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.PublicChainStatus\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenByAddress\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structTokenStructs.Token\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"uri\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"publicTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"issuerChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"ercStandard\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.ErcStandard\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\"},{\"name\":\"hubStatus\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.HubStatus\"},{\"name\":\"publicChainStatus\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.PublicChainStatus\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenByResourceId\",\"inputs\":[{\"name\":\"tokenResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structTokenStructs.Token\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"uri\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"publicTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"issuerChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"ercStandard\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.ErcStandard\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\"},{\"name\":\"hubStatus\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.HubStatus\"},{\"name\":\"publicChainStatus\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.PublicChainStatus\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenBySymbol\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structTokenStructs.Token\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"uri\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"publicTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"issuerChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"ercStandard\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.ErcStandard\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\"},{\"name\":\"hubStatus\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.HubStatus\"},{\"name\":\"publicChainStatus\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.PublicChainStatus\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenCore\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractITokenCore\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenFreezeManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractITokenFreezeManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokensByHubStatus\",\"inputs\":[{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.HubStatus\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structTokenStructs.Token[]\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"uri\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"publicTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"issuerChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"ercStandard\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.ErcStandard\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\"},{\"name\":\"hubStatus\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.HubStatus\"},{\"name\":\"publicChainStatus\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.PublicChainStatus\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokensByPrivacyNodeStatus\",\"inputs\":[{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structTokenStructs.Token[]\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"uri\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"publicTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"issuerChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"ercStandard\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.ErcStandard\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\"},{\"name\":\"hubStatus\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.HubStatus\"},{\"name\":\"publicChainStatus\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.PublicChainStatus\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"endpointAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"authority_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isTokenActiveForHub\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isTokenActiveForPublicChain\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isTokenFullyOperational\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerHubToken\",\"inputs\":[{\"name\":\"tokenResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerToken\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rejectToken\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeFrozenToken\",\"inputs\":[{\"name\":\"unfrozenToken\",\"type\":\"tuple\",\"internalType\":\"structTokenStructs.FrozenToken\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"frozenParticipants\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestAllFrozenTokensDataFromPrivateHub\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resourceId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setTokenCore\",\"inputs\":[{\"name\":\"tokenCoreAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setTokenFreezeManager\",\"inputs\":[{\"name\":\"freezeManagerAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitToHub\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitToPublicChain\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"syncFrozenTokens\",\"inputs\":[{\"name\":\"frozenTokens\",\"type\":\"tuple[]\",\"internalType\":\"structTokenStructs.FrozenToken[]\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"frozenParticipants\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tokenExists\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unfreezeOnPrivacyNode\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unfreezeOnPublicChain\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateFrozenToken\",\"inputs\":[{\"name\":\"frozenToken\",\"type\":\"tuple\",\"internalType\":\"structTokenStructs.FrozenToken\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"frozenParticipants\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updatePrivacyNodeStatus\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updatePublicTokenAddress\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"publicTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"validateTokenForParticipant\",\"inputs\":[{\"name\":\"tokenResourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenCoreSet\",\"inputs\":[{\"name\":\"tokenCore\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenFreezeManagerSet\",\"inputs\":[{\"name\":\"tokenFreezeManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenRegistryModulesConfigured\",\"inputs\":[{\"name\":\"tokenCore\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenFreezeManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenRegistryV1__InvalidTokenCoreAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TokenRegistryV1__InvalidTokenFreezeManagerAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TokenRegistryV1__TokenCoreNotConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TokenRegistryV1__TokenFreezeManagerNotConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	ID:  "PNTokenRegistryV1",
	Bin: "0x60a0604052306080523480156200001557600080fd5b506200002062000026565b620000da565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff1615620000775760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b0390811614620000d75780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b608051612a2862000104600039600081816119ec01528181611a150152611b4c0152612a286000f3fe6080604052600436106102525760003560e01c80636f0656f01161013a578063b33f78ca116100b1578063b33f78ca146106ff578063bf7e214f1461071f578063c3034ec014610734578063c4d66de814610754578063d8dc051014610774578063de0c362514610794578063e3ade0ec146107a9578063efa74f1f146107c9578063f05ab9de146107e9578063f0a599bb14610809578063f356cc0314610827578063f5aeb4341461084557600080fd5b80636f0656f01461057657806378a895671461059657806378dd4e6f146105ab5780637ae93604146105cb5780637fbc532d146105eb5780638a58b79f1461060b57806391ded8fa1461063857806393d85727146106585780639eb7b4eb1461066d578063a0a8e4601461068d578063a0b9b9b1146106a1578063ad3cb1cc146106c157600080fd5b80632a5c792a116101ce5780632a5c792a146104085780632f73e6e41461042a5780633b9373491461044a578063457415671461046a57806347ecc9391461048a578063485cc955146104aa5780634da8e6c3146104ca5780634f1ef286146104ea57806352d1902d146104fd5780635e280f11146105205780635f997c5b1461054057806364c0673a1461055657600080fd5b806214bc461461025757806301162d431461028c578063060ed143146102b957806309824a80146102db5780630a75cf4e146102fb57806310dcd5201461031b57806311f50c851461033b57806313aa1a1d1461036857806314ae646b146103885780631fe56159146103a8578063222e1b16146103c857806322e594ff146103e8575b600080fd5b34801561026357600080fd5b50610277610272366004611f36565b610865565b60405190151581526020015b60405180910390f35b34801561029857600080fd5b506102ac6102a7366004611f36565b6108e0565b6040516102839190611f86565b3480156102c557600080fd5b506102d96102d4366004611f99565b610956565b005b3480156102e757600080fd5b506102d96102f6366004611f36565b610a45565b34801561030757600080fd5b506102d9610316366004611f36565b610ac3565b34801561032757600080fd5b506102d9610336366004611f36565b610b0c565b34801561034757600080fd5b5061035b610356366004611fe1565b610b55565b6040516102839190611ffa565b34801561037457600080fd5b506102d9610383366004611f36565b610bc3565b34801561039457600080fd5b506102d96103a3366004611f36565b610c89565b3480156103b457600080fd5b506102d96103c3366004611f36565b610d4f565b3480156103d457600080fd5b506102d96103e336600461200e565b610d98565b3480156103f457600080fd5b506102ac610403366004611f36565b610e19565b34801561041457600080fd5b5061041d610e4e565b60405161028391906121e9565b34801561043657600080fd5b506102d9610445366004611f36565b610ec2565b34801561045657600080fd5b5061041d61046536600461225a565b610f0b565b34801561047657600080fd5b50610277610485366004611f36565b610f85565b34801561049657600080fd5b506102d96104a5366004612277565b610fba565b3480156104b657600080fd5b506102d96104c53660046122a7565b61100f565b3480156104d657600080fd5b506102d96104e5366004611f36565b611120565b6102d96104f83660046123a9565b611169565b34801561050957600080fd5b50610512611188565b604051908152602001610283565b34801561052c57600080fd5b5060005461035b906001600160a01b031681565b34801561054c57600080fd5b5061051260015481565b34801561056257600080fd5b50610277610571366004611f36565b6111a6565b34801561058257600080fd5b5061041d61059136600461225a565b6111db565b3480156105a257600080fd5b50610512611210565b3480156105b757600080fd5b506102776105c636600461240c565b61127b565b3480156105d757600080fd5b506102d96105e6366004611f36565b6112fe565b3480156105f757600080fd5b506102d9610606366004611f36565b611347565b34801561061757600080fd5b5061062b610626366004611fe1565b611390565b604051610283919061242e565b34801561064457600080fd5b5061062b610653366004611f36565b611412565b34801561066457600080fd5b5061051261144d565b34801561067957600080fd5b506102d9610688366004612441565b611494565b34801561069957600080fd5b506001610512565b3480156106ad57600080fd5b506102d96106bc3660046122a7565b6114dd565b3480156106cd57600080fd5b506106f2604051806040016040528060058152602001640352e302e360dc1b81525081565b604051610283919061247b565b34801561070b57600080fd5b5061027761071a366004611f36565b611533565b34801561072b57600080fd5b5061035b611568565b34801561074057600080fd5b506102d961074f36600461240c565b611581565b34801561076057600080fd5b506102d961076f366004611f36565b6115e6565b34801561078057600080fd5b506102d961078f36600461248e565b611610565b3480156107a057600080fd5b506102d961165b565b3480156107b557600080fd5b506102ac6107c4366004611f36565b6116cd565b3480156107d557600080fd5b5061062b6107e43660046124bc565b611702565b3480156107f557600080fd5b506102d9610804366004611f36565b61173d565b34801561081557600080fd5b506002546001600160a01b031661035b565b34801561083357600080fd5b506003546001600160a01b031661035b565b34801561085157600080fd5b506102d9610860366004612441565b611786565b600061086f6117cf565b6001600160a01b03166214bc46836040518263ffffffff1660e01b81526004016108999190611ffa565b602060405180830381865afa1580156108b6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108da9190612521565b92915050565b60006108ea6117cf565b6001600160a01b03166301162d43836040518263ffffffff1660e01b81526004016109159190611ffa565b602060405180830381865afa158015610932573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108da9190612547565b61096c336000356001600160e01b0319166117f8565b60405163a01afbfb60e01b8152600481018490526001600160a01b0383169063a01afbfb90602401600060405180830381600087803b1580156109ae57600080fd5b505af11580156109c2573d6000803e3d6000fd5b5050505060006109d06117cf565b60405163060ed14360e01b8152600481018690526001600160a01b03858116602483015260ff851660448301529192509082169063060ed14390606401600060405180830381600087803b158015610a2757600080fd5b505af1158015610a3b573d6000803e3d6000fd5b5050505050505050565b610a5b336000356001600160e01b0319166117f8565b610a636117cf565b6001600160a01b03166309824a80826040518263ffffffff1660e01b8152600401610a8e9190611ffa565b600060405180830381600087803b158015610aa857600080fd5b505af1158015610abc573d6000803e3d6000fd5b5050505050565b610ad9336000356001600160e01b0319166117f8565b610ae16117cf565b6001600160a01b0316630a75cf4e826040518263ffffffff1660e01b8152600401610a8e9190611ffa565b610b22336000356001600160e01b0319166117f8565b610b2a611943565b6001600160a01b03166310dcd520826040518263ffffffff1660e01b8152600401610a8e9190611ffa565b600080546040516311f50c8560e01b8152600481018490526001600160a01b03909116906311f50c8590602401602060405180830381865afa158015610b9f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108da919061256f565b610bd9336000356001600160e01b0319166117f8565b6001600160a01b038116610c00576040516365dce40760e11b815260040160405180910390fd5b600380546001600160a01b0319166001600160a01b0383169081179091556040517f2dd3698b1731a899f810b168007b9bd902aec761bfe35f98fb17348bb174a92190600090a26003546002546040516001600160a01b0392831692909116907f70e091790811d89fdfce2f92748c357fdb8847d6e97f65b45fecf8d412d32deb90600090a350565b610c9f336000356001600160e01b0319166117f8565b6001600160a01b038116610cc65760405163743840cd60e01b815260040160405180910390fd5b600280546001600160a01b0319166001600160a01b0383169081179091556040517fcd3287a073a6c154081e3c3168a5f5966880588e313e9dd388e935d177bec14b90600090a26003546002546040516001600160a01b0392831692909116907f70e091790811d89fdfce2f92748c357fdb8847d6e97f65b45fecf8d412d32deb90600090a350565b610d65336000356001600160e01b0319166117f8565b610d6d6117cf565b6001600160a01b0316631fe56159826040518263ffffffff1660e01b8152600401610a8e9190611ffa565b610dae336000356001600160e01b0319166117f8565b610db6611943565b6001600160a01b031663222e1b1683836040518363ffffffff1660e01b8152600401610de3929190612610565b600060405180830381600087803b158015610dfd57600080fd5b505af1158015610e11573d6000803e3d6000fd5b505050505050565b6000610e236117cf565b6001600160a01b03166322e594ff836040518263ffffffff1660e01b81526004016109159190611ffa565b6060610e586117cf565b6001600160a01b0316632a5c792a6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610e95573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610ebd91908101906127ff565b905090565b610ed8336000356001600160e01b0319166117f8565b610ee0611943565b6001600160a01b0316632f73e6e4826040518263ffffffff1660e01b8152600401610a8e9190611ffa565b6060610f156117cf565b6001600160a01b0316633b937349836040518263ffffffff1660e01b8152600401610f409190611f86565b600060405180830381865afa158015610f5d573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526108da91908101906127ff565b6000610f8f6117cf565b6001600160a01b03166345741567836040518263ffffffff1660e01b81526004016108999190611ffa565b610fd0336000356001600160e01b0319166117f8565b610fd86117cf565b6040516347ecc93960e01b8152600481018490526001600160a01b03838116602483015291909116906347ecc93990604401610de3565b600061101961196d565b805490915060ff600160401b82041615906001600160401b03166000811580156110405750825b90506000826001600160401b0316600114801561105c5750303b155b90508115801561106a575080155b156110885760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156110b257845460ff60401b1916600160401b1785555b6110ba611996565b6110c3876115e6565b60036001556110d1866119a0565b831561111757845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50505050505050565b611136336000356001600160e01b0319166117f8565b61113e611943565b6001600160a01b0316634da8e6c3826040518263ffffffff1660e01b8152600401610a8e9190611ffa565b6111716119e1565b61117a82611a6f565b6111848282611a88565b5050565b6000611192611b41565b506000805160206129d38339815191525b90565b60006111b06117cf565b6001600160a01b03166364c0673a836040518263ffffffff1660e01b81526004016108999190611ffa565b60606111e56117cf565b6001600160a01b0316636f0656f0836040518263ffffffff1660e01b8152600401610f409190611f86565b600061121a6117cf565b6001600160a01b03166378a895676040518163ffffffff1660e01b8152600401602060405180830381865afa158015611257573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ebd91906128c1565b6000611285611943565b6040516378dd4e6f60e01b815260048101859052602481018490526001600160a01b0391909116906378dd4e6f90604401602060405180830381865afa1580156112d3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112f79190612521565b9392505050565b611314336000356001600160e01b0319166117f8565b61131c611943565b6001600160a01b0316637ae93604826040518263ffffffff1660e01b8152600401610a8e9190611ffa565b61135d336000356001600160e01b0319166117f8565b6113656117cf565b6001600160a01b0316637fbc532d826040518263ffffffff1660e01b8152600401610a8e9190611ffa565b611398611e93565b6113a06117cf565b6001600160a01b0316638a58b79f836040518263ffffffff1660e01b81526004016113cd91815260200190565b600060405180830381865afa1580156113ea573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526108da91908101906128da565b61141a611e93565b6114226117cf565b6001600160a01b03166391ded8fa836040518263ffffffff1660e01b81526004016113cd9190611ffa565b60006114576117cf565b6001600160a01b03166393d857276040518163ffffffff1660e01b8152600401602060405180830381865afa158015611257573d6000803e3d6000fd5b6114aa336000356001600160e01b0319166117f8565b6114b2611943565b6001600160a01b0316639eb7b4eb826040518263ffffffff1660e01b8152600401610a8e919061290e565b6114f3336000356001600160e01b0319166117f8565b6114fb6117cf565b60405163a0b9b9b160e01b81526001600160a01b0384811660048301528381166024830152919091169063a0b9b9b190604401610de3565b600061153d6117cf565b6001600160a01b031663b33f78ca836040518263ffffffff1660e01b81526004016108999190611ffa565b6000611572611b8a565b546001600160a01b0316919050565b611589611943565b60405163030c0d3b60e61b815260048101849052602481018390526001600160a01b03919091169063c3034ec09060440160006040518083038186803b1580156115d257600080fd5b505afa158015610e11573d6000803e3d6000fd5b6115ee611bec565b600080546001600160a01b0319166001600160a01b0392909216919091179055565b611626336000356001600160e01b0319166117f8565b61162e6117cf565b6001600160a01b031663d8dc051083836040518363ffffffff1660e01b8152600401610de3929190612921565b611671336000356001600160e01b0319166117f8565b611679611943565b6001600160a01b031663de0c36256040518163ffffffff1660e01b8152600401600060405180830381600087803b1580156116b357600080fd5b505af11580156116c7573d6000803e3d6000fd5b50505050565b60006116d76117cf565b6001600160a01b031663e3ade0ec836040518263ffffffff1660e01b81526004016109159190611ffa565b61170a611e93565b6117126117cf565b6001600160a01b031663efa74f1f836040518263ffffffff1660e01b81526004016113cd919061247b565b611753336000356001600160e01b0319166117f8565b61175b6117cf565b6001600160a01b031663f05ab9de826040518263ffffffff1660e01b8152600401610a8e9190611ffa565b61179c336000356001600160e01b0319166117f8565b6117a4611943565b6001600160a01b031663f5aeb434826040518263ffffffff1660e01b8152600401610a8e919061290e565b6002546001600160a01b0316806111a35760405162dee40560e41b815260040160405180910390fd5b6000611802611b8a565b80549091506001600160a01b03168061183a576000604051638944034760e01b81526004016118319190611ffa565b60405180910390fd5b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa15801561189e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906118c29190612947565b925092509250826111175780156118ec5760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff8216156119285760405163a426878960e01b81526001600160a01b038816600482015263ffffffff83166024820152604401611831565b86604051632ecd3d0360e21b81526004016118319190611ffa565b6003546001600160a01b0316806111a357604051635e11777d60e01b815260040160405180910390fd5b6000807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a006108da565b61199e611bec565b565b60006119aa611b8a565b80549091506001600160a01b0316156119d85781604051638944034760e01b81526004016118319190611ffa565b61118482611c11565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161480611a5157507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316611a45611ca1565b6001600160a01b031614155b1561199e5760405163703e46dd60e11b815260040160405180910390fd5b611a85336000356001600160e01b0319166117f8565b50565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015611ae2575060408051601f3d908101601f19168201909252611adf918101906128c1565b60015b611b015781604051634c9c8ce360e01b81526004016118319190611ffa565b6000805160206129d38339815191528114611b3257604051632a87526960e21b815260048101829052602401611831565b611b3c8383611cb7565b505050565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161461199e5760405163703e46dd60e11b815260040160405180910390fd5b60008060ff19611bbb60017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f35612995565b604051602001611bcd91815260200190565b60408051601f1981840301815291905280516020909101201692915050565b611bf4611d0d565b61199e57604051631afcd79f60e31b815260040160405180910390fd5b6001600160a01b038116611c3a5780604051638944034760e01b81526004016118319190611ffa565b6000611c44611b8a565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b60006000805160206129d3833981519152611572565b611cc082611d27565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a2805115611d0557611b3c8282611d83565b611184611df9565b6000611d1761196d565b54600160401b900460ff16919050565b806001600160a01b03163b600003611d545780604051634c9c8ce360e01b81526004016118319190611ffa565b6000805160206129d383398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080846001600160a01b031684604051611da091906129b6565b600060405180830381855af49150503d8060008114611ddb576040519150601f19603f3d011682016040523d82523d6000602084013e611de0565b606091505b5091509150611df0858383611e18565b95945050505050565b341561199e5760405163b398979f60e01b815260040160405180910390fd5b606082611e2d57611e2882611e6b565b6112f7565b8151158015611e4457506001600160a01b0384163b155b15611e645783604051639996b31560e01b81526004016118319190611ffa565b5092915050565b805115611e7a57805160208201fd5b60405163d6bda27560e01b815260040160405180910390fd5b604051806101a001604052806000801916815260200160608152602001606081526020016060815260200160006001600160a01b0316815260200160006001600160a01b03168152602001600081526020016000600c811115611ef857611ef8611f53565b815260200160008152602001600081526020016000815260200160008152602001600081525090565b6001600160a01b0381168114611a8557600080fd5b600060208284031215611f4857600080fd5b81356112f781611f21565b634e487b7160e01b600052602160045260246000fd5b60058110611a8557611a85611f53565b611f8281611f69565b9052565b60208101611f9383611f69565b91905290565b600080600060608486031215611fae57600080fd5b833592506020840135611fc081611f21565b9150604084013560ff81168114611fd657600080fd5b809150509250925092565b600060208284031215611ff357600080fd5b5035919050565b6001600160a01b0391909116815260200190565b6000806020838503121561202157600080fd5b82356001600160401b038082111561203857600080fd5b818501915085601f83011261204c57600080fd5b81358181111561205b57600080fd5b8660208260051b850101111561207057600080fd5b60209290920196919550909350505050565b60005b8381101561209d578181015183820152602001612085565b50506000910152565b600081518084526120be816020860160208601612082565b601f01601f19169290920160200192915050565b600d8110611f8257611f82611f53565b60006101a0825184526020830151816020860152612102828601826120a6565b9150506040830151848203604086015261211c82826120a6565b9150506060830151848203606086015261213682826120a6565b915050608083015161215360808601826001600160a01b03169052565b5060a083015161216e60a08601826001600160a01b03169052565b5060c083015160c085015260e083015161218b60e08601826120d2565b506101008084015161219f82870182611f79565b5050610120808401516121b482870182611f79565b5050610140808401516121c982870182611f79565b505061016083810151908501526101809283015192909301919091525090565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b8281101561224057603f1988860301845261222e8583516120e2565b94509285019290850190600101612212565b5092979650505050505050565b60058110611a8557600080fd5b60006020828403121561226c57600080fd5b81356112f78161224d565b6000806040838503121561228a57600080fd5b82359150602083013561229c81611f21565b809150509250929050565b600080604083850312156122ba57600080fd5b82356122c581611f21565b9150602083013561229c81611f21565b634e487b7160e01b600052604160045260246000fd5b6040516101a081016001600160401b038111828210171561230e5761230e6122d5565b60405290565b604051601f8201601f191681016001600160401b038111828210171561233c5761233c6122d5565b604052919050565b60006001600160401b0382111561235d5761235d6122d5565b50601f01601f191660200190565b600061237e61237984612344565b612314565b905082815283838301111561239257600080fd5b828260208301376000602084830101529392505050565b600080604083850312156123bc57600080fd5b82356123c781611f21565b915060208301356001600160401b038111156123e257600080fd5b8301601f810185136123f357600080fd5b6124028582356020840161236b565b9150509250929050565b6000806040838503121561241f57600080fd5b50508035926020909101359150565b6020815260006112f760208301846120e2565b60006020828403121561245357600080fd5b81356001600160401b0381111561246957600080fd5b8201604081850312156112f757600080fd5b6020815260006112f760208301846120a6565b600080604083850312156124a157600080fd5b82356124ac81611f21565b9150602083013561229c8161224d565b6000602082840312156124ce57600080fd5b81356001600160401b038111156124e457600080fd5b8201601f810184136124f557600080fd5b6125048482356020840161236b565b949350505050565b8051801515811461251c57600080fd5b919050565b60006020828403121561253357600080fd5b6112f78261250c565b805161251c8161224d565b60006020828403121561255957600080fd5b81516112f78161224d565b805161251c81611f21565b60006020828403121561258157600080fd5b81516112f781611f21565b8035825260006020820135601e198336030181126125a957600080fd5b82016020810190356001600160401b038111156125c557600080fd5b8060051b8036038313156125d857600080fd5b60406020870181905286018290526001600160fb1b038211156125fa57600080fd5b8083606088013794909401606001949350505050565b60208082528181018390526000906040600585901b840181019084018684805b8881101561267457878503603f190184528235368b9003603e19018112612655578283fd5b612661868c830161258c565b9550509285019291850191600101612630565b509298975050505050505050565b600082601f83011261269357600080fd5b81516126a161237982612344565b8181528460208386010111156126b657600080fd5b612504826020830160208701612082565b8051600d811061251c57600080fd5b60006101a082840312156126e957600080fd5b6126f16122eb565b90508151815260208201516001600160401b038082111561271157600080fd5b61271d85838601612682565b6020840152604084015191508082111561273657600080fd5b61274285838601612682565b6040840152606084015191508082111561275b57600080fd5b5061276884828501612682565b60608301525061277a60808301612564565b608082015261278b60a08301612564565b60a082015260c082015160c08201526127a660e083016126c7565b60e08201526101006127b981840161253c565b908201526101206127cb83820161253c565b908201526101406127dd83820161253c565b9082015261016082810151908201526101809182015191810191909152919050565b6000602080838503121561281257600080fd5b82516001600160401b038082111561282957600080fd5b818501915085601f83011261283d57600080fd5b81518181111561284f5761284f6122d5565b8060051b61285e858201612314565b918252838101850191858101908984111561287857600080fd5b86860192505b838310156128b4578251858111156128965760008081fd5b6128a48b89838a01016126d6565b835250918601919086019061287e565b9998505050505050505050565b6000602082840312156128d357600080fd5b5051919050565b6000602082840312156128ec57600080fd5b81516001600160401b0381111561290257600080fd5b612504848285016126d6565b6020815260006112f7602083018461258c565b6001600160a01b03831681526040810161293a83611f69565b8260208301529392505050565b60008060006060848603121561295c57600080fd5b6129658461250c565b9250602084015163ffffffff8116811461297e57600080fd5b915061298c6040850161250c565b90509250925092565b818103818111156108da57634e487b7160e01b600052601160045260246000fd5b600082516129c8818460208701612082565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbca2646970667358221220603d2507f6a43ae911866f6ab2a081d91d0b1ebfb2712a79ec5b143fcf565dd064736f6c63430008180033",
}

// PNTokenRegistryV1 is an auto generated Go binding around an Ethereum contract.
type PNTokenRegistryV1 struct {
	abi abi.ABI
}

// NewPNTokenRegistryV1 creates a new instance of PNTokenRegistryV1.
func NewPNTokenRegistryV1() *PNTokenRegistryV1 {
	parsed, err := PNTokenRegistryV1MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &PNTokenRegistryV1{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *PNTokenRegistryV1) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackUPGRADEINTERFACEVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackUPGRADEINTERFACEVERSION() []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("UPGRADE_INTERFACE_VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUPGRADEINTERFACEVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackUPGRADEINTERFACEVERSION(data []byte) (string, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("UPGRADE_INTERFACE_VERSION", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackActivateToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x060ed143.
//
// Solidity: function activateToken(bytes32 tokenResourceId, address tokenAddress, uint8 ercStandard) returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackActivateToken(tokenResourceId [32]byte, tokenAddress common.Address, ercStandard uint8) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("activateToken", tokenResourceId, tokenAddress, ercStandard)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackAuthority() []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("authority", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackContractVersion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa0a8e460.
//
// Solidity: function contractVersion() pure returns(uint256)
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackContractVersion() []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("contractVersion")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackContractVersion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa0a8e460.
//
// Solidity: function contractVersion() pure returns(uint256)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackContractVersion(data []byte) (*big.Int, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("contractVersion", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackDeprecateOnPublicChain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1fe56159.
//
// Solidity: function deprecateOnPublicChain(address tokenAddress) returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackDeprecateOnPublicChain(tokenAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("deprecateOnPublicChain", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackEndpoint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackEndpoint() []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("endpoint")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackEndpoint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackEndpoint(data []byte) (common.Address, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("endpoint", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackFreezeOnPrivacyNode is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f73e6e4.
//
// Solidity: function freezeOnPrivacyNode(address tokenAddress) returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackFreezeOnPrivacyNode(tokenAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("freezeOnPrivacyNode", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackFreezeOnPublicChain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4da8e6c3.
//
// Solidity: function freezeOnPublicChain(address tokenAddress) returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackFreezeOnPublicChain(tokenAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("freezeOnPublicChain", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackGetActiveTokenCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x93d85727.
//
// Solidity: function getActiveTokenCount() view returns(uint256)
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackGetActiveTokenCount() []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("getActiveTokenCount")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetActiveTokenCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x93d85727.
//
// Solidity: function getActiveTokenCount() view returns(uint256)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackGetActiveTokenCount(data []byte) (*big.Int, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("getActiveTokenCount", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetAddressByResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackGetAddressByResourceId(resourceId [32]byte) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("getAddressByResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAddressByResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackGetAddressByResourceId(data []byte) (common.Address, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("getAddressByResourceId", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetAllTokens is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2a5c792a.
//
// Solidity: function getAllTokens() view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256)[])
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackGetAllTokens() []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("getAllTokens")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAllTokens is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2a5c792a.
//
// Solidity: function getAllTokens() view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256)[])
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackGetAllTokens(data []byte) ([]TokenStructsToken, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("getAllTokens", data)
	if err != nil {
		return *new([]TokenStructsToken), err
	}
	out0 := *abi.ConvertType(out[0], new([]TokenStructsToken)).(*[]TokenStructsToken)
	return out0, err
}

// PackGetFrozenTokenForParticipant is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x78dd4e6f.
//
// Solidity: function getFrozenTokenForParticipant(bytes32 tokenResourceId, uint256 chainId) view returns(bool)
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackGetFrozenTokenForParticipant(tokenResourceId [32]byte, chainId *big.Int) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("getFrozenTokenForParticipant", tokenResourceId, chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetFrozenTokenForParticipant is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x78dd4e6f.
//
// Solidity: function getFrozenTokenForParticipant(bytes32 tokenResourceId, uint256 chainId) view returns(bool)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackGetFrozenTokenForParticipant(data []byte) (bool, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("getFrozenTokenForParticipant", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackGetHubStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe3ade0ec.
//
// Solidity: function getHubStatus(address tokenAddress) view returns(uint8)
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackGetHubStatus(tokenAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("getHubStatus", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetHubStatus is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe3ade0ec.
//
// Solidity: function getHubStatus(address tokenAddress) view returns(uint8)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackGetHubStatus(data []byte) (uint8, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("getHubStatus", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, err
}

// PackGetPrivacyNodeStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x22e594ff.
//
// Solidity: function getPrivacyNodeStatus(address tokenAddress) view returns(uint8)
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackGetPrivacyNodeStatus(tokenAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("getPrivacyNodeStatus", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetPrivacyNodeStatus is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x22e594ff.
//
// Solidity: function getPrivacyNodeStatus(address tokenAddress) view returns(uint8)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackGetPrivacyNodeStatus(data []byte) (uint8, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("getPrivacyNodeStatus", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, err
}

// PackGetPublicChainStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01162d43.
//
// Solidity: function getPublicChainStatus(address tokenAddress) view returns(uint8)
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackGetPublicChainStatus(tokenAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("getPublicChainStatus", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetPublicChainStatus is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01162d43.
//
// Solidity: function getPublicChainStatus(address tokenAddress) view returns(uint8)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackGetPublicChainStatus(data []byte) (uint8, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("getPublicChainStatus", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, err
}

// PackGetTokenByAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91ded8fa.
//
// Solidity: function getTokenByAddress(address tokenAddress) view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256))
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackGetTokenByAddress(tokenAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("getTokenByAddress", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTokenByAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x91ded8fa.
//
// Solidity: function getTokenByAddress(address tokenAddress) view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256))
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackGetTokenByAddress(data []byte) (TokenStructsToken, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("getTokenByAddress", data)
	if err != nil {
		return *new(TokenStructsToken), err
	}
	out0 := *abi.ConvertType(out[0], new(TokenStructsToken)).(*TokenStructsToken)
	return out0, err
}

// PackGetTokenByResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8a58b79f.
//
// Solidity: function getTokenByResourceId(bytes32 tokenResourceId) view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256))
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackGetTokenByResourceId(tokenResourceId [32]byte) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("getTokenByResourceId", tokenResourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTokenByResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8a58b79f.
//
// Solidity: function getTokenByResourceId(bytes32 tokenResourceId) view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256))
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackGetTokenByResourceId(data []byte) (TokenStructsToken, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("getTokenByResourceId", data)
	if err != nil {
		return *new(TokenStructsToken), err
	}
	out0 := *abi.ConvertType(out[0], new(TokenStructsToken)).(*TokenStructsToken)
	return out0, err
}

// PackGetTokenBySymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xefa74f1f.
//
// Solidity: function getTokenBySymbol(string symbol) view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256))
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackGetTokenBySymbol(symbol string) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("getTokenBySymbol", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTokenBySymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xefa74f1f.
//
// Solidity: function getTokenBySymbol(string symbol) view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256))
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackGetTokenBySymbol(data []byte) (TokenStructsToken, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("getTokenBySymbol", data)
	if err != nil {
		return *new(TokenStructsToken), err
	}
	out0 := *abi.ConvertType(out[0], new(TokenStructsToken)).(*TokenStructsToken)
	return out0, err
}

// PackGetTokenCore is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf0a599bb.
//
// Solidity: function getTokenCore() view returns(address)
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackGetTokenCore() []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("getTokenCore")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTokenCore is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf0a599bb.
//
// Solidity: function getTokenCore() view returns(address)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackGetTokenCore(data []byte) (common.Address, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("getTokenCore", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetTokenCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x78a89567.
//
// Solidity: function getTokenCount() view returns(uint256)
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackGetTokenCount() []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("getTokenCount")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTokenCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x78a89567.
//
// Solidity: function getTokenCount() view returns(uint256)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackGetTokenCount(data []byte) (*big.Int, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("getTokenCount", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetTokenFreezeManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf356cc03.
//
// Solidity: function getTokenFreezeManager() view returns(address)
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackGetTokenFreezeManager() []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("getTokenFreezeManager")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTokenFreezeManager is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf356cc03.
//
// Solidity: function getTokenFreezeManager() view returns(address)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackGetTokenFreezeManager(data []byte) (common.Address, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("getTokenFreezeManager", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetTokensByHubStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6f0656f0.
//
// Solidity: function getTokensByHubStatus(uint8 status) view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256)[])
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackGetTokensByHubStatus(status uint8) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("getTokensByHubStatus", status)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTokensByHubStatus is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6f0656f0.
//
// Solidity: function getTokensByHubStatus(uint8 status) view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256)[])
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackGetTokensByHubStatus(data []byte) ([]TokenStructsToken, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("getTokensByHubStatus", data)
	if err != nil {
		return *new([]TokenStructsToken), err
	}
	out0 := *abi.ConvertType(out[0], new([]TokenStructsToken)).(*[]TokenStructsToken)
	return out0, err
}

// PackGetTokensByPrivacyNodeStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3b937349.
//
// Solidity: function getTokensByPrivacyNodeStatus(uint8 status) view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256)[])
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackGetTokensByPrivacyNodeStatus(status uint8) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("getTokensByPrivacyNodeStatus", status)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTokensByPrivacyNodeStatus is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3b937349.
//
// Solidity: function getTokensByPrivacyNodeStatus(uint8 status) view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256)[])
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackGetTokensByPrivacyNodeStatus(data []byte) ([]TokenStructsToken, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("getTokensByPrivacyNodeStatus", data)
	if err != nil {
		return *new([]TokenStructsToken), err
	}
	out0 := *abi.ConvertType(out[0], new([]TokenStructsToken)).(*[]TokenStructsToken)
	return out0, err
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x485cc955.
//
// Solidity: function initialize(address endpointAddress, address authority_) returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackInitialize(endpointAddress common.Address, authority common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("initialize", endpointAddress, authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackInitialize0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc4d66de8.
//
// Solidity: function initialize(address _endpoint) returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackInitialize0(endpoint common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("initialize0", endpoint)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackIsTokenActiveForHub is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0014bc46.
//
// Solidity: function isTokenActiveForHub(address tokenAddress) view returns(bool)
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackIsTokenActiveForHub(tokenAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("isTokenActiveForHub", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsTokenActiveForHub is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0014bc46.
//
// Solidity: function isTokenActiveForHub(address tokenAddress) view returns(bool)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackIsTokenActiveForHub(data []byte) (bool, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("isTokenActiveForHub", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackIsTokenActiveForPublicChain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x64c0673a.
//
// Solidity: function isTokenActiveForPublicChain(address tokenAddress) view returns(bool)
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackIsTokenActiveForPublicChain(tokenAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("isTokenActiveForPublicChain", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsTokenActiveForPublicChain is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x64c0673a.
//
// Solidity: function isTokenActiveForPublicChain(address tokenAddress) view returns(bool)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackIsTokenActiveForPublicChain(data []byte) (bool, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("isTokenActiveForPublicChain", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackIsTokenFullyOperational is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x45741567.
//
// Solidity: function isTokenFullyOperational(address tokenAddress) view returns(bool)
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackIsTokenFullyOperational(tokenAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("isTokenFullyOperational", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsTokenFullyOperational is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x45741567.
//
// Solidity: function isTokenFullyOperational(address tokenAddress) view returns(bool)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackIsTokenFullyOperational(data []byte) (bool, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("isTokenFullyOperational", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackProxiableUUID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackProxiableUUID() []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("proxiableUUID")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackProxiableUUID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackProxiableUUID(data []byte) ([32]byte, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("proxiableUUID", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackRegisterHubToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x47ecc939.
//
// Solidity: function registerHubToken(bytes32 tokenResourceId, address tokenAddress) returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackRegisterHubToken(tokenResourceId [32]byte, tokenAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("registerHubToken", tokenResourceId, tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRegisterToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x09824a80.
//
// Solidity: function registerToken(address tokenAddress) returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackRegisterToken(tokenAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("registerToken", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRejectToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf05ab9de.
//
// Solidity: function rejectToken(address tokenAddress) returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackRejectToken(tokenAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("rejectToken", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRemoveFrozenToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf5aeb434.
//
// Solidity: function removeFrozenToken((bytes32,uint256[]) unfrozenToken) returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackRemoveFrozenToken(unfrozenToken TokenStructsFrozenToken) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("removeFrozenToken", unfrozenToken)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRequestAllFrozenTokensDataFromPrivateHub is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xde0c3625.
//
// Solidity: function requestAllFrozenTokensDataFromPrivateHub() returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackRequestAllFrozenTokensDataFromPrivateHub() []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("requestAllFrozenTokensDataFromPrivateHub")
	if err != nil {
		panic(err)
	}
	return enc
}

// PackResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackResourceId() []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("resourceId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackResourceId(data []byte) ([32]byte, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("resourceId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackSetTokenCore is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x14ae646b.
//
// Solidity: function setTokenCore(address tokenCoreAddress) returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackSetTokenCore(tokenCoreAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("setTokenCore", tokenCoreAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetTokenFreezeManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x13aa1a1d.
//
// Solidity: function setTokenFreezeManager(address freezeManagerAddress) returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackSetTokenFreezeManager(freezeManagerAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("setTokenFreezeManager", freezeManagerAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSubmitToHub is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7fbc532d.
//
// Solidity: function submitToHub(address tokenAddress) returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackSubmitToHub(tokenAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("submitToHub", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSubmitToPublicChain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0a75cf4e.
//
// Solidity: function submitToPublicChain(address tokenAddress) returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackSubmitToPublicChain(tokenAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("submitToPublicChain", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSyncFrozenTokens is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x222e1b16.
//
// Solidity: function syncFrozenTokens((bytes32,uint256[])[] frozenTokens) returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackSyncFrozenTokens(frozenTokens []TokenStructsFrozenToken) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("syncFrozenTokens", frozenTokens)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackTokenExists is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb33f78ca.
//
// Solidity: function tokenExists(address tokenAddress) view returns(bool)
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackTokenExists(tokenAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("tokenExists", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTokenExists is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb33f78ca.
//
// Solidity: function tokenExists(address tokenAddress) view returns(bool)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackTokenExists(data []byte) (bool, error) {
	out, err := pNTokenRegistryV1.abi.Unpack("tokenExists", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackUnfreezeOnPrivacyNode is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x10dcd520.
//
// Solidity: function unfreezeOnPrivacyNode(address tokenAddress) returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackUnfreezeOnPrivacyNode(tokenAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("unfreezeOnPrivacyNode", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUnfreezeOnPublicChain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7ae93604.
//
// Solidity: function unfreezeOnPublicChain(address tokenAddress) returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackUnfreezeOnPublicChain(tokenAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("unfreezeOnPublicChain", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpdateFrozenToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9eb7b4eb.
//
// Solidity: function updateFrozenToken((bytes32,uint256[]) frozenToken) returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackUpdateFrozenToken(frozenToken TokenStructsFrozenToken) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("updateFrozenToken", frozenToken)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpdatePrivacyNodeStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd8dc0510.
//
// Solidity: function updatePrivacyNodeStatus(address tokenAddress, uint8 status) returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackUpdatePrivacyNodeStatus(tokenAddress common.Address, status uint8) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("updatePrivacyNodeStatus", tokenAddress, status)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpdatePublicTokenAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa0b9b9b1.
//
// Solidity: function updatePublicTokenAddress(address tokenAddress, address publicTokenAddress) returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackUpdatePublicTokenAddress(tokenAddress common.Address, publicTokenAddress common.Address) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("updatePublicTokenAddress", tokenAddress, publicTokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpgradeToAndCall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackUpgradeToAndCall(newImplementation common.Address, data []byte) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("upgradeToAndCall", newImplementation, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackValidateTokenForParticipant is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc3034ec0.
//
// Solidity: function validateTokenForParticipant(bytes32 tokenResourceId, uint256 chainId) view returns()
func (pNTokenRegistryV1 *PNTokenRegistryV1) PackValidateTokenForParticipant(tokenResourceId [32]byte, chainId *big.Int) []byte {
	enc, err := pNTokenRegistryV1.abi.Pack("validateTokenForParticipant", tokenResourceId, chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PNTokenRegistryV1AuthorityUpdated represents a AuthorityUpdated event raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1AuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const PNTokenRegistryV1AuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (PNTokenRegistryV1AuthorityUpdated) ContractEventName() string {
	return PNTokenRegistryV1AuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackAuthorityUpdatedEvent(log *types.Log) (*PNTokenRegistryV1AuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != pNTokenRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNTokenRegistryV1AuthorityUpdated)
	if len(log.Data) > 0 {
		if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNTokenRegistryV1.abi.Events[event].Inputs {
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

// PNTokenRegistryV1Initialized represents a Initialized event raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1Initialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const PNTokenRegistryV1InitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (PNTokenRegistryV1Initialized) ContractEventName() string {
	return PNTokenRegistryV1InitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackInitializedEvent(log *types.Log) (*PNTokenRegistryV1Initialized, error) {
	event := "Initialized"
	if log.Topics[0] != pNTokenRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNTokenRegistryV1Initialized)
	if len(log.Data) > 0 {
		if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNTokenRegistryV1.abi.Events[event].Inputs {
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

// PNTokenRegistryV1TokenCoreSet represents a TokenCoreSet event raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1TokenCoreSet struct {
	TokenCore common.Address
	Raw       *types.Log // Blockchain specific contextual infos
}

const PNTokenRegistryV1TokenCoreSetEventName = "TokenCoreSet"

// ContractEventName returns the user-defined event name.
func (PNTokenRegistryV1TokenCoreSet) ContractEventName() string {
	return PNTokenRegistryV1TokenCoreSetEventName
}

// UnpackTokenCoreSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenCoreSet(address indexed tokenCore)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackTokenCoreSetEvent(log *types.Log) (*PNTokenRegistryV1TokenCoreSet, error) {
	event := "TokenCoreSet"
	if log.Topics[0] != pNTokenRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNTokenRegistryV1TokenCoreSet)
	if len(log.Data) > 0 {
		if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNTokenRegistryV1.abi.Events[event].Inputs {
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

// PNTokenRegistryV1TokenFreezeManagerSet represents a TokenFreezeManagerSet event raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1TokenFreezeManagerSet struct {
	TokenFreezeManager common.Address
	Raw                *types.Log // Blockchain specific contextual infos
}

const PNTokenRegistryV1TokenFreezeManagerSetEventName = "TokenFreezeManagerSet"

// ContractEventName returns the user-defined event name.
func (PNTokenRegistryV1TokenFreezeManagerSet) ContractEventName() string {
	return PNTokenRegistryV1TokenFreezeManagerSetEventName
}

// UnpackTokenFreezeManagerSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenFreezeManagerSet(address indexed tokenFreezeManager)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackTokenFreezeManagerSetEvent(log *types.Log) (*PNTokenRegistryV1TokenFreezeManagerSet, error) {
	event := "TokenFreezeManagerSet"
	if log.Topics[0] != pNTokenRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNTokenRegistryV1TokenFreezeManagerSet)
	if len(log.Data) > 0 {
		if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNTokenRegistryV1.abi.Events[event].Inputs {
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

// PNTokenRegistryV1TokenRegistryModulesConfigured represents a TokenRegistryModulesConfigured event raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1TokenRegistryModulesConfigured struct {
	TokenCore          common.Address
	TokenFreezeManager common.Address
	Raw                *types.Log // Blockchain specific contextual infos
}

const PNTokenRegistryV1TokenRegistryModulesConfiguredEventName = "TokenRegistryModulesConfigured"

// ContractEventName returns the user-defined event name.
func (PNTokenRegistryV1TokenRegistryModulesConfigured) ContractEventName() string {
	return PNTokenRegistryV1TokenRegistryModulesConfiguredEventName
}

// UnpackTokenRegistryModulesConfiguredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenRegistryModulesConfigured(address indexed tokenCore, address indexed tokenFreezeManager)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackTokenRegistryModulesConfiguredEvent(log *types.Log) (*PNTokenRegistryV1TokenRegistryModulesConfigured, error) {
	event := "TokenRegistryModulesConfigured"
	if log.Topics[0] != pNTokenRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNTokenRegistryV1TokenRegistryModulesConfigured)
	if len(log.Data) > 0 {
		if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNTokenRegistryV1.abi.Events[event].Inputs {
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

// PNTokenRegistryV1Upgraded represents a Upgraded event raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1Upgraded struct {
	Implementation common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const PNTokenRegistryV1UpgradedEventName = "Upgraded"

// ContractEventName returns the user-defined event name.
func (PNTokenRegistryV1Upgraded) ContractEventName() string {
	return PNTokenRegistryV1UpgradedEventName
}

// UnpackUpgradedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Upgraded(address indexed implementation)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackUpgradedEvent(log *types.Log) (*PNTokenRegistryV1Upgraded, error) {
	event := "Upgraded"
	if log.Topics[0] != pNTokenRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNTokenRegistryV1Upgraded)
	if len(log.Data) > 0 {
		if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNTokenRegistryV1.abi.Events[event].Inputs {
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
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], pNTokenRegistryV1.abi.Errors["AddressEmptyCode"].ID.Bytes()[:4]) {
		return pNTokenRegistryV1.UnpackAddressEmptyCodeError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenRegistryV1.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return pNTokenRegistryV1.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenRegistryV1.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return pNTokenRegistryV1.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenRegistryV1.abi.Errors["FailedCall"].ID.Bytes()[:4]) {
		return pNTokenRegistryV1.UnpackFailedCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenRegistryV1.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return pNTokenRegistryV1.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenRegistryV1.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return pNTokenRegistryV1.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenRegistryV1.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return pNTokenRegistryV1.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenRegistryV1.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return pNTokenRegistryV1.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenRegistryV1.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return pNTokenRegistryV1.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenRegistryV1.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return pNTokenRegistryV1.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenRegistryV1.abi.Errors["TokenRegistryV1InvalidTokenCoreAddress"].ID.Bytes()[:4]) {
		return pNTokenRegistryV1.UnpackTokenRegistryV1InvalidTokenCoreAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenRegistryV1.abi.Errors["TokenRegistryV1InvalidTokenFreezeManagerAddress"].ID.Bytes()[:4]) {
		return pNTokenRegistryV1.UnpackTokenRegistryV1InvalidTokenFreezeManagerAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenRegistryV1.abi.Errors["TokenRegistryV1TokenCoreNotConfigured"].ID.Bytes()[:4]) {
		return pNTokenRegistryV1.UnpackTokenRegistryV1TokenCoreNotConfiguredError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenRegistryV1.abi.Errors["TokenRegistryV1TokenFreezeManagerNotConfigured"].ID.Bytes()[:4]) {
		return pNTokenRegistryV1.UnpackTokenRegistryV1TokenFreezeManagerNotConfiguredError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenRegistryV1.abi.Errors["UUPSUnauthorizedCallContext"].ID.Bytes()[:4]) {
		return pNTokenRegistryV1.UnpackUUPSUnauthorizedCallContextError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenRegistryV1.abi.Errors["UUPSUnsupportedProxiableUUID"].ID.Bytes()[:4]) {
		return pNTokenRegistryV1.UnpackUUPSUnsupportedProxiableUUIDError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// PNTokenRegistryV1AddressEmptyCode represents a AddressEmptyCode error raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1AddressEmptyCode struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressEmptyCode(address target)
func PNTokenRegistryV1AddressEmptyCodeErrorID() common.Hash {
	return common.HexToHash("0x9996b315c842ff135b8fc4a08ad5df1c344efbc03d2687aecc0678050d2aac89")
}

// UnpackAddressEmptyCodeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressEmptyCode(address target)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackAddressEmptyCodeError(raw []byte) (*PNTokenRegistryV1AddressEmptyCode, error) {
	out := new(PNTokenRegistryV1AddressEmptyCode)
	if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, "AddressEmptyCode", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenRegistryV1ERC1967InvalidImplementation represents a ERC1967InvalidImplementation error raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1ERC1967InvalidImplementation struct {
	Implementation common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func PNTokenRegistryV1ERC1967InvalidImplementationErrorID() common.Hash {
	return common.HexToHash("0x4c9c8ce3ceb3130f17f7cdba48d89b5b0129f266a8bac114e6e315a41879b617")
}

// UnpackERC1967InvalidImplementationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackERC1967InvalidImplementationError(raw []byte) (*PNTokenRegistryV1ERC1967InvalidImplementation, error) {
	out := new(PNTokenRegistryV1ERC1967InvalidImplementation)
	if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, "ERC1967InvalidImplementation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenRegistryV1ERC1967NonPayable represents a ERC1967NonPayable error raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1ERC1967NonPayable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967NonPayable()
func PNTokenRegistryV1ERC1967NonPayableErrorID() common.Hash {
	return common.HexToHash("0xb398979fa84f543c8e222f17890372c487baf85e062276c127fef521eea7224b")
}

// UnpackERC1967NonPayableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967NonPayable()
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackERC1967NonPayableError(raw []byte) (*PNTokenRegistryV1ERC1967NonPayable, error) {
	out := new(PNTokenRegistryV1ERC1967NonPayable)
	if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, "ERC1967NonPayable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenRegistryV1FailedCall represents a FailedCall error raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1FailedCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedCall()
func PNTokenRegistryV1FailedCallErrorID() common.Hash {
	return common.HexToHash("0xd6bda27508c0fb6d8a39b4b122878dab26f731a7d4e4abe711dd3731899052a4")
}

// UnpackFailedCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedCall()
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackFailedCallError(raw []byte) (*PNTokenRegistryV1FailedCall, error) {
	out := new(PNTokenRegistryV1FailedCall)
	if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, "FailedCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenRegistryV1InvalidInitialization represents a InvalidInitialization error raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1InvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func PNTokenRegistryV1InvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackInvalidInitializationError(raw []byte) (*PNTokenRegistryV1InvalidInitialization, error) {
	out := new(PNTokenRegistryV1InvalidInitialization)
	if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenRegistryV1NotInitializing represents a NotInitializing error raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1NotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func PNTokenRegistryV1NotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackNotInitializingError(raw []byte) (*PNTokenRegistryV1NotInitializing, error) {
	out := new(PNTokenRegistryV1NotInitializing)
	if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenRegistryV1RaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1RaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func PNTokenRegistryV1RaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*PNTokenRegistryV1RaylsAccessManagedContractPaused, error) {
	out := new(PNTokenRegistryV1RaylsAccessManagedContractPaused)
	if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenRegistryV1RaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1RaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func PNTokenRegistryV1RaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*PNTokenRegistryV1RaylsAccessManagedInvalidAuthority, error) {
	out := new(PNTokenRegistryV1RaylsAccessManagedInvalidAuthority)
	if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenRegistryV1RaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1RaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func PNTokenRegistryV1RaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*PNTokenRegistryV1RaylsAccessManagedMustSchedule, error) {
	out := new(PNTokenRegistryV1RaylsAccessManagedMustSchedule)
	if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenRegistryV1RaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1RaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func PNTokenRegistryV1RaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*PNTokenRegistryV1RaylsAccessManagedUnauthorized, error) {
	out := new(PNTokenRegistryV1RaylsAccessManagedUnauthorized)
	if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenRegistryV1TokenRegistryV1InvalidTokenCoreAddress represents a TokenRegistryV1__InvalidTokenCoreAddress error raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1TokenRegistryV1InvalidTokenCoreAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenRegistryV1__InvalidTokenCoreAddress()
func PNTokenRegistryV1TokenRegistryV1InvalidTokenCoreAddressErrorID() common.Hash {
	return common.HexToHash("0x743840cd612f7e36e5814125efeb3c972964e780bff028a10d4e7ba4b7bd47f3")
}

// UnpackTokenRegistryV1InvalidTokenCoreAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenRegistryV1__InvalidTokenCoreAddress()
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackTokenRegistryV1InvalidTokenCoreAddressError(raw []byte) (*PNTokenRegistryV1TokenRegistryV1InvalidTokenCoreAddress, error) {
	out := new(PNTokenRegistryV1TokenRegistryV1InvalidTokenCoreAddress)
	if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, "TokenRegistryV1InvalidTokenCoreAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenRegistryV1TokenRegistryV1InvalidTokenFreezeManagerAddress represents a TokenRegistryV1__InvalidTokenFreezeManagerAddress error raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1TokenRegistryV1InvalidTokenFreezeManagerAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenRegistryV1__InvalidTokenFreezeManagerAddress()
func PNTokenRegistryV1TokenRegistryV1InvalidTokenFreezeManagerAddressErrorID() common.Hash {
	return common.HexToHash("0xcbb9c80e5b7e9993353677692525f58d399465ade8d5206d324f23f7cca6eaf9")
}

// UnpackTokenRegistryV1InvalidTokenFreezeManagerAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenRegistryV1__InvalidTokenFreezeManagerAddress()
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackTokenRegistryV1InvalidTokenFreezeManagerAddressError(raw []byte) (*PNTokenRegistryV1TokenRegistryV1InvalidTokenFreezeManagerAddress, error) {
	out := new(PNTokenRegistryV1TokenRegistryV1InvalidTokenFreezeManagerAddress)
	if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, "TokenRegistryV1InvalidTokenFreezeManagerAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenRegistryV1TokenRegistryV1TokenCoreNotConfigured represents a TokenRegistryV1__TokenCoreNotConfigured error raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1TokenRegistryV1TokenCoreNotConfigured struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenRegistryV1__TokenCoreNotConfigured()
func PNTokenRegistryV1TokenRegistryV1TokenCoreNotConfiguredErrorID() common.Hash {
	return common.HexToHash("0x0dee4050ae8d6c6f1f485a566a6e78c577c2f1134e26cef638d9ecc511270fa2")
}

// UnpackTokenRegistryV1TokenCoreNotConfiguredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenRegistryV1__TokenCoreNotConfigured()
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackTokenRegistryV1TokenCoreNotConfiguredError(raw []byte) (*PNTokenRegistryV1TokenRegistryV1TokenCoreNotConfigured, error) {
	out := new(PNTokenRegistryV1TokenRegistryV1TokenCoreNotConfigured)
	if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, "TokenRegistryV1TokenCoreNotConfigured", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenRegistryV1TokenRegistryV1TokenFreezeManagerNotConfigured represents a TokenRegistryV1__TokenFreezeManagerNotConfigured error raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1TokenRegistryV1TokenFreezeManagerNotConfigured struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenRegistryV1__TokenFreezeManagerNotConfigured()
func PNTokenRegistryV1TokenRegistryV1TokenFreezeManagerNotConfiguredErrorID() common.Hash {
	return common.HexToHash("0x5e11777d4fffe4b1cbb3c3552b2d82e0fb7c9d31e45257f293bfae903a2645b0")
}

// UnpackTokenRegistryV1TokenFreezeManagerNotConfiguredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenRegistryV1__TokenFreezeManagerNotConfigured()
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackTokenRegistryV1TokenFreezeManagerNotConfiguredError(raw []byte) (*PNTokenRegistryV1TokenRegistryV1TokenFreezeManagerNotConfigured, error) {
	out := new(PNTokenRegistryV1TokenRegistryV1TokenFreezeManagerNotConfigured)
	if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, "TokenRegistryV1TokenFreezeManagerNotConfigured", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenRegistryV1UUPSUnauthorizedCallContext represents a UUPSUnauthorizedCallContext error raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1UUPSUnauthorizedCallContext struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnauthorizedCallContext()
func PNTokenRegistryV1UUPSUnauthorizedCallContextErrorID() common.Hash {
	return common.HexToHash("0xe07c8dba242a06571ac65fe4bbe20522c9fb111cb33599b799ff8039c1ed18f4")
}

// UnpackUUPSUnauthorizedCallContextError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnauthorizedCallContext()
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackUUPSUnauthorizedCallContextError(raw []byte) (*PNTokenRegistryV1UUPSUnauthorizedCallContext, error) {
	out := new(PNTokenRegistryV1UUPSUnauthorizedCallContext)
	if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, "UUPSUnauthorizedCallContext", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenRegistryV1UUPSUnsupportedProxiableUUID represents a UUPSUnsupportedProxiableUUID error raised by the PNTokenRegistryV1 contract.
type PNTokenRegistryV1UUPSUnsupportedProxiableUUID struct {
	Slot [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func PNTokenRegistryV1UUPSUnsupportedProxiableUUIDErrorID() common.Hash {
	return common.HexToHash("0xaa1d49a4c084bfa9aeeee2a0be65267a7f19ba7e1476b114dac513d2c14cb563")
}

// UnpackUUPSUnsupportedProxiableUUIDError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func (pNTokenRegistryV1 *PNTokenRegistryV1) UnpackUUPSUnsupportedProxiableUUIDError(raw []byte) (*PNTokenRegistryV1UUPSUnsupportedProxiableUUID, error) {
	out := new(PNTokenRegistryV1UUPSUnsupportedProxiableUUID)
	if err := pNTokenRegistryV1.abi.UnpackIntoInterface(out, "UUPSUnsupportedProxiableUUID", raw); err != nil {
		return nil, err
	}
	return out, nil
}
