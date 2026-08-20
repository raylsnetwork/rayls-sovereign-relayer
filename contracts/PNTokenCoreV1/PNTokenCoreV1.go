// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package PNTokenCoreV1

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

// PNTokenCoreLibMetaData contains all meta data concerning the PNTokenCoreLib contract.
var PNTokenCoreLibMetaData = bind.MetaData{
	ABI: "[]",
	ID:  "7c2d9d4b3468ae44947ccbf0e482c64ec5",
}

// PNTokenCoreLib is an auto generated Go binding around an Ethereum contract.
type PNTokenCoreLib struct {
	abi abi.ABI
}

// NewPNTokenCoreLib creates a new instance of PNTokenCoreLib.
func NewPNTokenCoreLib() *PNTokenCoreLib {
	parsed, err := PNTokenCoreLibMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &PNTokenCoreLib{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *PNTokenCoreLib) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PNTokenCoreV1MetaData contains all meta data concerning the PNTokenCoreV1 contract.
var PNTokenCoreV1MetaData = bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"UPGRADE_INTERFACE_VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"resourceId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"ercStandard\",\"type\":\"uint8\"}],\"name\":\"activateToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"authority\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"deprecateOnPublicChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"endpoint\",\"outputs\":[{\"internalType\":\"contractIRaylsEndpoint\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enygmaPNEvents\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getActiveTokenCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllTokens\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"resourceId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"publicTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"issuerChainId\",\"type\":\"uint256\"},{\"internalType\":\"enumSharedObjects.ErcStandard\",\"name\":\"ercStandard\",\"type\":\"uint8\"},{\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\",\"name\":\"privacyNodeStatus\",\"type\":\"uint8\"},{\"internalType\":\"enumTokenStructs.HubStatus\",\"name\":\"hubStatus\",\"type\":\"uint8\"},{\"internalType\":\"enumTokenStructs.PublicChainStatus\",\"name\":\"publicChainStatus\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"createdAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"internalType\":\"structTokenStructs.Token[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"getHubStatus\",\"outputs\":[{\"internalType\":\"enumTokenStructs.HubStatus\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"getPrivacyNodeStatus\",\"outputs\":[{\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"getPublicChainStatus\",\"outputs\":[{\"internalType\":\"enumTokenStructs.PublicChainStatus\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"getTokenByAddress\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"resourceId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"publicTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"issuerChainId\",\"type\":\"uint256\"},{\"internalType\":\"enumSharedObjects.ErcStandard\",\"name\":\"ercStandard\",\"type\":\"uint8\"},{\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\",\"name\":\"privacyNodeStatus\",\"type\":\"uint8\"},{\"internalType\":\"enumTokenStructs.HubStatus\",\"name\":\"hubStatus\",\"type\":\"uint8\"},{\"internalType\":\"enumTokenStructs.PublicChainStatus\",\"name\":\"publicChainStatus\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"createdAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"internalType\":\"structTokenStructs.Token\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"resourceId\",\"type\":\"bytes32\"}],\"name\":\"getTokenByResourceId\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"resourceId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"publicTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"issuerChainId\",\"type\":\"uint256\"},{\"internalType\":\"enumSharedObjects.ErcStandard\",\"name\":\"ercStandard\",\"type\":\"uint8\"},{\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\",\"name\":\"privacyNodeStatus\",\"type\":\"uint8\"},{\"internalType\":\"enumTokenStructs.HubStatus\",\"name\":\"hubStatus\",\"type\":\"uint8\"},{\"internalType\":\"enumTokenStructs.PublicChainStatus\",\"name\":\"publicChainStatus\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"createdAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"internalType\":\"structTokenStructs.Token\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"getTokenBySymbol\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"resourceId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"publicTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"issuerChainId\",\"type\":\"uint256\"},{\"internalType\":\"enumSharedObjects.ErcStandard\",\"name\":\"ercStandard\",\"type\":\"uint8\"},{\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\",\"name\":\"privacyNodeStatus\",\"type\":\"uint8\"},{\"internalType\":\"enumTokenStructs.HubStatus\",\"name\":\"hubStatus\",\"type\":\"uint8\"},{\"internalType\":\"enumTokenStructs.PublicChainStatus\",\"name\":\"publicChainStatus\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"createdAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"internalType\":\"structTokenStructs.Token\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTokenCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumTokenStructs.HubStatus\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"getTokensByHubStatus\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"resourceId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"publicTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"issuerChainId\",\"type\":\"uint256\"},{\"internalType\":\"enumSharedObjects.ErcStandard\",\"name\":\"ercStandard\",\"type\":\"uint8\"},{\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\",\"name\":\"privacyNodeStatus\",\"type\":\"uint8\"},{\"internalType\":\"enumTokenStructs.HubStatus\",\"name\":\"hubStatus\",\"type\":\"uint8\"},{\"internalType\":\"enumTokenStructs.PublicChainStatus\",\"name\":\"publicChainStatus\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"createdAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"internalType\":\"structTokenStructs.Token[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"getTokensByPrivacyNodeStatus\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"resourceId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"publicTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"issuerChainId\",\"type\":\"uint256\"},{\"internalType\":\"enumSharedObjects.ErcStandard\",\"name\":\"ercStandard\",\"type\":\"uint8\"},{\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\",\"name\":\"privacyNodeStatus\",\"type\":\"uint8\"},{\"internalType\":\"enumTokenStructs.HubStatus\",\"name\":\"hubStatus\",\"type\":\"uint8\"},{\"internalType\":\"enumTokenStructs.PublicChainStatus\",\"name\":\"publicChainStatus\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"createdAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"internalType\":\"structTokenStructs.Token[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"authority_\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"isTokenActiveForHub\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"isTokenActiveForPublicChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"isTokenFullyOperational\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"resourceId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"registerHubToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"registerToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"rejectToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"relayAuthorizationRegistry\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"endpointAddress\",\"type\":\"address\"}],\"name\":\"setEndpoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"enygmaPNEventsAddress\",\"type\":\"address\"}],\"name\":\"setEnygmaPNEvents\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"enumTokenStructs.FreezeLayer\",\"name\":\"layer\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"frozen\",\"type\":\"bool\"}],\"name\":\"setFreezeStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"relayAuthRegistry\",\"type\":\"address\"}],\"name\":\"setRelayAuthorizationRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"freezeManager\",\"type\":\"address\"}],\"name\":\"setTokenFreezeManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenRegistry\",\"type\":\"address\"}],\"name\":\"setTokenRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"submitToHub\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"submitToPublicChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"tokenExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"resourceId\",\"type\":\"bytes32\"}],\"name\":\"tokenExistsByResourceId\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenFreezeManager\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenRegistryAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"updatePrivacyNodeStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"publicTokenAddress\",\"type\":\"address\"}],\"name\":\"updatePublicTokenAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldAuthority\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAuthority\",\"type\":\"address\"}],\"name\":\"AuthorityUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumTokenStructs.HubStatus\",\"name\":\"previousStatus\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumTokenStructs.HubStatus\",\"name\":\"newStatus\",\"type\":\"uint8\"}],\"name\":\"HubStatusUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\",\"name\":\"previousStatus\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\",\"name\":\"newStatus\",\"type\":\"uint8\"}],\"name\":\"PrivacyNodeStatusUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumTokenStructs.PublicChainStatus\",\"name\":\"previousStatus\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumTokenStructs.PublicChainStatus\",\"name\":\"newStatus\",\"type\":\"uint8\"}],\"name\":\"PublicChainStatusUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"publicTokenAddress\",\"type\":\"address\"}],\"name\":\"PublicTokenAddressUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"resourceId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumSharedObjects.ErcStandard\",\"name\":\"ercStandard\",\"type\":\"uint8\"}],\"name\":\"TokenActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"tokenRegistry\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"tokenFreezeManager\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"endpoint\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"enygmaPNEvents\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"relayAuthorizationRegistry\",\"type\":\"address\"}],\"name\":\"TokenCoreModulesConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"enumSharedObjects.ErcStandard\",\"name\":\"ercStandard\",\"type\":\"uint8\"}],\"name\":\"TokenRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RaylsAccessManaged__ContractPaused\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"delay\",\"type\":\"uint32\"}],\"name\":\"RaylsAccessManaged__MustSchedule\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"RaylsAccessManaged__Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenCoreV1__EndpointNotConfigured\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"enumTokenStructs.HubStatus\",\"name\":\"currentStatus\",\"type\":\"uint8\"}],\"name\":\"TokenCoreV1__HubApprovalNotPending\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenCoreV1__InvalidEndpointAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenCoreV1__InvalidEnygmaPNEventsAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"enumTokenStructs.PrivacyNodeStatus\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"TokenCoreV1__InvalidPrivacyNodeStatus\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenCoreV1__InvalidTokenAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenCoreV1__InvalidTokenFreezeManagerAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenCoreV1__InvalidTokenRegistryAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenCoreV1__PrivacyNodeAuthorizationRequired\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"resourceId\",\"type\":\"bytes32\"}],\"name\":\"TokenCoreV1__ResourceIdAlreadyAssigned\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenCoreV1__ResourceIdAlreadySet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenCoreV1__StatusAlreadySet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenCoreV1__TokenAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenCoreV1__TokenDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenCoreV1__TokenRegistryNotConfigured\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenCoreV1__TokenSymbolAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"TokenCoreV1__UnauthorizedCaller\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"TokenCoreV1__UnauthorizedFreezeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UUPSUnauthorizedCallContext\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"UUPSUnsupportedProxiableUUID\",\"type\":\"error\"}]",
	ID:  "3d6ed29fc7e4bc9b7edc149018a527513f",
	Bin: "0x60a0604052306080523480156200001557600080fd5b506200002062000026565b620000da565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff1615620000775760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b0390811614620000d75780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b60805161538b6200010460003960008181613d9001528181613db90152613eea015261538b6000f3fe6080604052600436106102055760003560e01c80636f0656f011610119578063b33f78ca116100a6578063b33f78ca14610617578063b58bdb4914610637578063bf7e214f14610657578063c4d66de81461066c578063cd4b4a511461068c578063d8dc0510146106ac578063dbbb4155146106cc578063dff11d4c146106ec578063e3ade0ec1461070c578063efa74f1f1461072c578063f05ab9de1461074c57600080fd5b80636f0656f0146104c257806378a89567146104e25780637fbc532d146104f75780638212f073146105175780638a58b79f1461053757806391ded8fa1461056457806393d85727146105845780639d53ea4c14610599578063a0b9b9b1146105b9578063ad3cb1cc146105d957600080fd5b806335a5af921161019757806335a5af921461037f5780633b9373491461039f5780633cd17302146103bf57806345741567146103df57806347ecc939146103ff5780634f1ef2861461041f57806352d1902d146104325780635be2aca0146104555780635e280f111461048257806364c0673a146104a257600080fd5b806214bc461461020a57806301162d431461023f578063060ed1431461026c57806309824a801461028e5780630a75cf4e146102ae5780631047100b146102ce57806313aa1a1d146102fd5780631fe561591461031d57806322e594ff1461033d5780632a5c792a1461035d575b600080fd5b34801561021657600080fd5b5061022a610225366004614828565b61076c565b60405190151581526020015b60405180910390f35b34801561024b57600080fd5b5061025f61025a366004614828565b6107ee565b6040516102369190614878565b34801561027857600080fd5b5061028c61028736600461488b565b61082e565b005b34801561029a57600080fd5b5061028c6102a9366004614828565b610c78565b3480156102ba57600080fd5b5061028c6102c9366004614828565b6111ae565b3480156102da57600080fd5b5061022a6102e93660046148d3565b600090815260026020526040902054151590565b34801561030957600080fd5b5061028c610318366004614828565b6112f3565b34801561032957600080fd5b5061028c610338366004614828565b61138f565b34801561034957600080fd5b5061025f610358366004614828565b6113d2565b34801561036957600080fd5b50610372611410565b6040516102369190614a53565b34801561038b57600080fd5b5061028c61039a366004614828565b6117be565b3480156103ab57600080fd5b506103726103ba366004614ac4565b611850565b3480156103cb57600080fd5b5061028c6103da366004614828565b611cd8565b3480156103eb57600080fd5b5061022a6103fa366004614828565b611d6b565b34801561040b57600080fd5b5061028c61041a366004614ae1565b611e0f565b61028c61042d366004614bbc565b61207a565b34801561043e57600080fd5b50610447612095565b604051908152602001610236565b34801561046157600080fd5b50600454610475906001600160a01b031681565b6040516102369190614c1f565b34801561048e57600080fd5b50600854610475906001600160a01b031681565b3480156104ae57600080fd5b5061022a6104bd366004614828565b6120b2565b3480156104ce57600080fd5b506103726104dd366004614ac4565b612102565b3480156104ee57600080fd5b50610447612580565b34801561050357600080fd5b5061028c610512366004614828565b6125a4565b34801561052357600080fd5b5061028c610532366004614828565b61294d565b34801561054357600080fd5b506105576105523660046148d3565b6129ba565b6040516102369190614c33565b34801561057057600080fd5b5061055761057f366004614828565b612d02565b34801561059057600080fd5b50610447613021565b3480156105a557600080fd5b50600554610475906001600160a01b031681565b3480156105c557600080fd5b5061028c6105d4366004614c46565b613097565b3480156105e557600080fd5b5061060a604051806040016040528060058152602001640352e302e360dc1b81525081565b6040516102369190614c74565b34801561062357600080fd5b5061022a610632366004614828565b613156565b34801561064357600080fd5b5061028c610652366004614c95565b613173565b34801561066357600080fd5b50610475613202565b34801561067857600080fd5b5061028c610687366004614828565b61321b565b34801561069857600080fd5b50600754610475906001600160a01b031681565b3480156106b857600080fd5b5061028c6106c7366004614cd9565b613328565b3480156106d857600080fd5b5061028c6106e7366004614828565b61339b565b3480156106f857600080fd5b50600654610475906001600160a01b031681565b34801561071857600080fd5b5061025f610727366004614828565b613430565b34801561073857600080fd5b50610557610747366004614d07565b61346f565b34801561075857600080fd5b5061028c610767366004614828565b6134bb565b600080600061077a846134fb565b8154811061078a5761078a614d57565b60009182526020909120600a90910201905060026007820154610100900460ff1660048111156107bc576107bc614845565b1480156107e757506002600782015462010000900460ff1660048111156107e5576107e5614845565b145b9392505050565b6000806107fa836134fb565b8154811061080a5761080a614d57565b60009182526020909120600a90910201600701546301000000900460ff1692915050565b6004546001600160a01b03163314610864573360405163394557ff60e21b815260040161085b9190614c1f565b60405180910390fd5b6008546001600160a01b031661088d5760405163f117d70960e01b815260040160405180910390fd5b6000610898836134fb565b905060008082815481106108ae576108ae614d57565b60009182526020909120600a909102018054909150158015906108d2575080548514155b156108f2578360405163f117323b60e01b815260040161085b9190614c1f565b60008581526002602052604090205480158015906109105750828114155b156109315760405163c3bf60f960e01b81526004810187905260240161085b565b6001600783015462010000900460ff16600481111561095257610952614845565b1461098557848260070160029054906101000a900460ff16604051630fa5bf8d60e21b815260040161085b929190614d6d565b85825560008681526002602052604090819020849055426009840155600854905163d4f951c760e01b81526001600160a01b039091169063d4f951c7906109d29089908990600401614d93565b600060405180830381600087803b1580156109ec57600080fd5b505af1158015610a00573d6000803e3d6000fd5b5050505073__$7c2d9d4b3468ae44947ccbf0e482c64ec5$__63febc04a1600860009054906101000a90046001600160a01b03166001600160a01b031663bf7e214f6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610a71573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a959190614daa565b876040518363ffffffff1660e01b8152600401610ab3929190614dc7565b60006040518083038186803b158015610acb57600080fd5b505af4158015610adf573d6000803e3d6000fd5b50505050600073__$7c2d9d4b3468ae44947ccbf0e482c64ec5$__63fecd5415878760ff16600c811115610b1557610b15614845565b6040518363ffffffff1660e01b8152600401610b32929190614de1565b602060405180830381865af4158015610b4f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b739190614dfe565b60065490915073__$7c2d9d4b3468ae44947ccbf0e482c64ec5$__9063c8f35196906001600160a01b03168960ff8916600c811115610bb457610bb4614845565b856040518563ffffffff1660e01b8152600401610bd49493929190614e17565b60006040518083038186803b158015610bec57600080fd5b505af4158015610c00573d6000803e3d6000fd5b50505050610c0f84600261353a565b856001600160a01b0316877f4e10bd9e15ecb70ffc648529318b2aa9c5fe9ca728e75df29cb427f4c49b22b58760ff16600c811115610c5057610c50614845565b604051610c5d9190614e4a565b60405180910390a3610c6f8785613632565b50505050505050565b6004546001600160a01b03163314610ca5573360405163394557ff60e21b815260040161085b9190614c1f565b6001600160a01b038116610ccc57604051636e8a975960e01b815260040160405180910390fd5b610cd581613156565b15610cf3576040516360c91b5960e11b815260040160405180910390fd5b6008546001600160a01b0316610d1c5760405163f117d70960e01b815260040160405180910390fd5b604051634676e7b960e01b815260009073__$7c2d9d4b3468ae44947ccbf0e482c64ec5$__90634676e7b990610d56908590600401614c1f565b602060405180830381865af4158015610d73573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d979190614e58565b9050600080600073__$7c2d9d4b3468ae44947ccbf0e482c64ec5$__634e24dbc286866040518363ffffffff1660e01b8152600401610dd7929190614de1565b600060405180830381865af4158015610df4573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610e1c9190810190614ec9565b925092509250600382604051610e329190614f50565b908152602001604051809103902054600014610e615760405163fae70dbd60e01b815260040160405180910390fd5b6000604051806101a001604052806000801b8152602001858152602001848152602001838152602001876001600160a01b0316815260200160006001600160a01b03168152602001600860009054906101000a90046001600160a01b03166001600160a01b0316633408e4706040518163ffffffff1660e01b8152600401602060405180830381865afa158015610efc573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f209190614dfe565b815260200186600c811115610f3757610f37614845565b81526020016001815260200160008152602001600081524260208083018290526040909201528254600181810185556000948552938290208351600a90920201908155908201519192909190820190610f909082614fe8565b5060408201516002820190610fa59082614fe8565b5060608201516003820190610fba9082614fe8565b5060808201516004820180546001600160a01b039283166001600160a01b03199182161790915560a084015160058401805491909316911617905560c0820151600682015560e082015160078201805460ff1916600183600c81111561102257611022614845565b0217905550610100828101516007830180549192909161ff0019169083600481111561105057611050614845565b021790555061012082015160078201805462ff000019166201000083600481111561107d5761107d614845565b021790555061014082015160078201805463ff000000191663010000008360048111156110ac576110ac614845565b0217905550610160820151600882015561018090910151600990910155600080546110d9906001906150bd565b6001600160a01b0387166000908152600160205260409081902082905551909150819060039061110a908690614f50565b908152602001604051809103902081905550856001600160a01b03167f063765b272ccfe8c5b747b73e618a2d710b3796d89b51946c128404ca6c85cc8858588604051611159939291906150d0565b60405180910390a2856001600160a01b03167f962ec3691afeb0804244549c8e163d9e4a89f507a6ed4a8a87c6fbca721699d26000600160405161119e929190615105565b60405180910390a2505050505050565b6004546001600160a01b031633146111db573360405163394557ff60e21b815260040161085b9190614c1f565b60006111e6826134fb565b905060008082815481106111fc576111fc614d57565b60009182526020909120600a90910201905060026007820154610100900460ff16600481111561122e5761122e614845565b1461124e5782604051631f9afd4360e11b815260040161085b9190614c1f565b600160078201546301000000900460ff16600481111561127057611270614845565b148061129b5750600260078201546301000000900460ff16600481111561129957611299614845565b145b806112c55750600360078201546301000000900460ff1660048111156112c3576112c3614845565b145b156112e357604051637d113f6d60e11b815260040160405180910390fd5b6112ee82600161374c565b505050565b611309336000356001600160e01b031916613839565b6001600160a01b038116611330576040516392e680fb60e01b815260040160405180910390fd5b600580546001600160a01b0319166001600160a01b03838116918217909255600854600454600654600754604051938616959283169360008051602061533683398151915293611384938116921690614dc7565b60405180910390a450565b6004546001600160a01b031633146113bc573360405163394557ff60e21b815260040161085b9190614c1f565b6113cf6113c8826134fb565b600461374c565b50565b6000806113de836134fb565b815481106113ee576113ee614d57565b60009182526020909120600a9091020160070154610100900460ff1692915050565b6060600061141c612580565b90506000816001600160401b0381111561143857611438614b11565b60405190808252806020026020018201604052801561147157816020015b61145e614785565b8152602001906001900390816114565790505b50905060005b828110156117b757600061148c82600161511e565b8154811061149c5761149c614d57565b90600052602060002090600a0201604051806101a0016040529081600082015481526020016001820180546114d090614f6c565b80601f01602080910402602001604051908101604052809291908181526020018280546114fc90614f6c565b80156115495780601f1061151e57610100808354040283529160200191611549565b820191906000526020600020905b81548152906001019060200180831161152c57829003601f168201915b5050505050815260200160028201805461156290614f6c565b80601f016020809104026020016040519081016040528092919081815260200182805461158e90614f6c565b80156115db5780601f106115b0576101008083540402835291602001916115db565b820191906000526020600020905b8154815290600101906020018083116115be57829003601f168201915b505050505081526020016003820180546115f490614f6c565b80601f016020809104026020016040519081016040528092919081815260200182805461162090614f6c565b801561166d5780601f106116425761010080835404028352916020019161166d565b820191906000526020600020905b81548152906001019060200180831161165057829003601f168201915b505050918352505060048201546001600160a01b039081166020830152600583015416604082015260068201546060820152600782015460809091019060ff16600c8111156116be576116be614845565b600c8111156116cf576116cf614845565b81526020016007820160019054906101000a900460ff1660048111156116f7576116f7614845565b600481111561170857611708614845565b81526020016007820160029054906101000a900460ff16600481111561173057611730614845565b600481111561174157611741614845565b81526020016007820160039054906101000a900460ff16600481111561176957611769614845565b600481111561177a5761177a614845565b8152602001600882015481526020016009820154815250508282815181106117a4576117a4614d57565b6020908102919091010152600101611477565b5092915050565b6117d4336000356001600160e01b031916613839565b6001600160a01b0381166117fb5760405163fce06f6360e01b815260040160405180910390fd5b600480546001600160a01b0319166001600160a01b0383811691821790925560085460055460065460075460405193861695928316949360008051602061533683398151915293611384938116921690614dc7565b6060600060015b6000548110156118d15783600481111561187357611873614845565b6000828154811061188657611886614d57565b90600052602060002090600a020160070160019054906101000a900460ff1660048111156118b6576118b6614845565b036118c957816118c581615131565b9250505b600101611857565b506000816001600160401b038111156118ec576118ec614b11565b60405190808252806020026020018201604052801561192557816020015b611912614785565b81526020019060019003908161190a5790505b509050600060015b600054811015611cce5785600481111561194957611949614845565b6000828154811061195c5761195c614d57565b90600052602060002090600a020160070160019054906101000a900460ff16600481111561198c5761198c614845565b03611cc657600081815481106119a4576119a4614d57565b90600052602060002090600a0201604051806101a0016040529081600082015481526020016001820180546119d890614f6c565b80601f0160208091040260200160405190810160405280929190818152602001828054611a0490614f6c565b8015611a515780601f10611a2657610100808354040283529160200191611a51565b820191906000526020600020905b815481529060010190602001808311611a3457829003601f168201915b50505050508152602001600282018054611a6a90614f6c565b80601f0160208091040260200160405190810160405280929190818152602001828054611a9690614f6c565b8015611ae35780601f10611ab857610100808354040283529160200191611ae3565b820191906000526020600020905b815481529060010190602001808311611ac657829003601f168201915b50505050508152602001600382018054611afc90614f6c565b80601f0160208091040260200160405190810160405280929190818152602001828054611b2890614f6c565b8015611b755780601f10611b4a57610100808354040283529160200191611b75565b820191906000526020600020905b815481529060010190602001808311611b5857829003601f168201915b505050918352505060048201546001600160a01b039081166020830152600583015416604082015260068201546060820152600782015460809091019060ff16600c811115611bc657611bc6614845565b600c811115611bd757611bd7614845565b81526020016007820160019054906101000a900460ff166004811115611bff57611bff614845565b6004811115611c1057611c10614845565b81526020016007820160029054906101000a900460ff166004811115611c3857611c38614845565b6004811115611c4957611c49614845565b81526020016007820160039054906101000a900460ff166004811115611c7157611c71614845565b6004811115611c8257611c82614845565b815260200160088201548152602001600982015481525050838381518110611cac57611cac614d57565b60200260200101819052508180611cc290615131565b9250505b60010161192d565b5090949350505050565b611cee336000356001600160e01b031916613839565b6001600160a01b038116611d1557604051632bac709560e01b815260040160405180910390fd5b600680546001600160a01b0319166001600160a01b038381169182179092556008546005546004546007546040519386169592831694918316936000805160206153368339815191529361138493921690614dc7565b6000806000611d79846134fb565b81548110611d8957611d89614d57565b60009182526020909120600a90910201905060026007820154610100900460ff166004811115611dbb57611dbb614845565b148015611de657506002600782015462010000900460ff166004811115611de457611de4614845565b145b80156107e75750600260078201546301000000900460ff1660048111156107e5576107e5614845565b6004546001600160a01b03163314611e3c573360405163394557ff60e21b815260040161085b9190614c1f565b6008546001600160a01b0316611e655760405163f117d70960e01b815260040160405180910390fd5b6001600160a01b038116611e8c57604051636e8a975960e01b815260040160405180910390fd5b611e9581613156565b6120765760008281526002602052604090205461207657604051634676e7b960e01b815260009073__$7c2d9d4b3468ae44947ccbf0e482c64ec5$__90634676e7b990611ee6908590600401614c1f565b602060405180830381865af4158015611f03573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f279190614e58565b60085460405163d4f951c760e01b81529192506001600160a01b03169063d4f951c790611f5a9086908690600401614d93565b600060405180830381600087803b158015611f7457600080fd5b505af1158015611f88573d6000803e3d6000fd5b5050505073__$7c2d9d4b3468ae44947ccbf0e482c64ec5$__63febc04a1600860009054906101000a90046001600160a01b03166001600160a01b031663bf7e214f6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611ff9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061201d9190614daa565b846040518363ffffffff1660e01b815260040161203b929190614dc7565b60006040518083038186803b15801561205357600080fd5b505af4158015612067573d6000803e3d6000fd5b505050506112ee83838361397b565b5050565b612082613d85565b61208b82613e15565b6120768282613e2b565b600061209f613edf565b5060008051602061531683398151915290565b60008060006120c0846134fb565b815481106120d0576120d0614d57565b60009182526020909120600a90910201905060026007820154610100900460ff166004811115611de457611de4614845565b6060600060015b6000548110156121835783600481111561212557612125614845565b6000828154811061213857612138614d57565b90600052602060002090600a020160070160029054906101000a900460ff16600481111561216857612168614845565b0361217b578161217781615131565b9250505b600101612109565b506000816001600160401b0381111561219e5761219e614b11565b6040519080825280602002602001820160405280156121d757816020015b6121c4614785565b8152602001906001900390816121bc5790505b509050600060015b600054811015611cce578560048111156121fb576121fb614845565b6000828154811061220e5761220e614d57565b90600052602060002090600a020160070160029054906101000a900460ff16600481111561223e5761223e614845565b03612578576000818154811061225657612256614d57565b90600052602060002090600a0201604051806101a00160405290816000820154815260200160018201805461228a90614f6c565b80601f01602080910402602001604051908101604052809291908181526020018280546122b690614f6c565b80156123035780601f106122d857610100808354040283529160200191612303565b820191906000526020600020905b8154815290600101906020018083116122e657829003601f168201915b5050505050815260200160028201805461231c90614f6c565b80601f016020809104026020016040519081016040528092919081815260200182805461234890614f6c565b80156123955780601f1061236a57610100808354040283529160200191612395565b820191906000526020600020905b81548152906001019060200180831161237857829003601f168201915b505050505081526020016003820180546123ae90614f6c565b80601f01602080910402602001604051908101604052809291908181526020018280546123da90614f6c565b80156124275780601f106123fc57610100808354040283529160200191612427565b820191906000526020600020905b81548152906001019060200180831161240a57829003601f168201915b505050918352505060048201546001600160a01b039081166020830152600583015416604082015260068201546060820152600782015460809091019060ff16600c81111561247857612478614845565b600c81111561248957612489614845565b81526020016007820160019054906101000a900460ff1660048111156124b1576124b1614845565b60048111156124c2576124c2614845565b81526020016007820160029054906101000a900460ff1660048111156124ea576124ea614845565b60048111156124fb576124fb614845565b81526020016007820160039054906101000a900460ff16600481111561252357612523614845565b600481111561253457612534614845565b81526020016008820154815260200160098201548152505083838151811061255e5761255e614d57565b6020026020010181905250818061257490615131565b9250505b6001016121df565b6000805481036125905750600090565b60005461259f906001906150bd565b905090565b6004546001600160a01b031633146125d1573360405163394557ff60e21b815260040161085b9190614c1f565b6008546001600160a01b03166125fa5760405163f117d70960e01b815260040160405180910390fd5b6004546001600160a01b03166126235760405163292914f360e11b815260040160405180910390fd5b600061262e826134fb565b9050600080828154811061264457612644614d57565b60009182526020909120600a90910201905060026007820154610100900460ff16600481111561267657612676614845565b146126965782604051631f9afd4360e11b815260040161085b9190614c1f565b6001600782015462010000900460ff1660048111156126b7576126b7614845565b14806126e157506002600782015462010000900460ff1660048111156126df576126df614845565b145b8061270a57506004600782015462010000900460ff16600481111561270857612708614845565b145b1561272857604051637d113f6d60e11b815260040160405180910390fd5b600781015460048054600684015460405163725e6e6760e01b815260009473__$7c2d9d4b3468ae44947ccbf0e482c64ec5$__9463725e6e6794612791948b9460ff909416936001600160a01b0390921692909160018b019160028c019160038d0191016151c7565b600060405180830381865af41580156127ae573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526127d6919081019061523b565b60085460408051630b39a95160e01b815290519293506001600160a01b039091169163150b375f918391630b39a951916004808201926020929091908290030181865afa15801561282b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061284f9190614dfe565b600854604051631a24be5560e21b815260206004820152600d60248201526c546f6b656e526567697374727960981b60448201526001600160a01b0390911690636892f95490606401602060405180830381865afa1580156128b5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906128d99190614daa565b846040518463ffffffff1660e01b81526004016128f893929190615283565b6020604051808303816000875af1158015612917573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061293b9190614dfe565b5061294783600161353a565b50505050565b612963336000356001600160e01b031916613839565b600780546001600160a01b0319166001600160a01b03838116918217909255600854600554600454600654604051938616959283169491831693600080516020615336833981519152936113849392169190614dc7565b6129c2614785565b600082815260026020526040812054908190036129f25760405163b2be0e7f60e01b815260040160405180910390fd5b60008181548110612a0557612a05614d57565b90600052602060002090600a0201604051806101a001604052908160008201548152602001600182018054612a3990614f6c565b80601f0160208091040260200160405190810160405280929190818152602001828054612a6590614f6c565b8015612ab25780601f10612a8757610100808354040283529160200191612ab2565b820191906000526020600020905b815481529060010190602001808311612a9557829003601f168201915b50505050508152602001600282018054612acb90614f6c565b80601f0160208091040260200160405190810160405280929190818152602001828054612af790614f6c565b8015612b445780601f10612b1957610100808354040283529160200191612b44565b820191906000526020600020905b815481529060010190602001808311612b2757829003601f168201915b50505050508152602001600382018054612b5d90614f6c565b80601f0160208091040260200160405190810160405280929190818152602001828054612b8990614f6c565b8015612bd65780601f10612bab57610100808354040283529160200191612bd6565b820191906000526020600020905b815481529060010190602001808311612bb957829003601f168201915b505050918352505060048201546001600160a01b039081166020830152600583015416604082015260068201546060820152600782015460809091019060ff16600c811115612c2757612c27614845565b600c811115612c3857612c38614845565b81526020016007820160019054906101000a900460ff166004811115612c6057612c60614845565b6004811115612c7157612c71614845565b81526020016007820160029054906101000a900460ff166004811115612c9957612c99614845565b6004811115612caa57612caa614845565b81526020016007820160039054906101000a900460ff166004811115612cd257612cd2614845565b6004811115612ce357612ce3614845565b8152600882015460208201526009909101546040909101529392505050565b612d0a614785565b6000612d15836134fb565b81548110612d2557612d25614d57565b90600052602060002090600a0201604051806101a001604052908160008201548152602001600182018054612d5990614f6c565b80601f0160208091040260200160405190810160405280929190818152602001828054612d8590614f6c565b8015612dd25780601f10612da757610100808354040283529160200191612dd2565b820191906000526020600020905b815481529060010190602001808311612db557829003601f168201915b50505050508152602001600282018054612deb90614f6c565b80601f0160208091040260200160405190810160405280929190818152602001828054612e1790614f6c565b8015612e645780601f10612e3957610100808354040283529160200191612e64565b820191906000526020600020905b815481529060010190602001808311612e4757829003601f168201915b50505050508152602001600382018054612e7d90614f6c565b80601f0160208091040260200160405190810160405280929190818152602001828054612ea990614f6c565b8015612ef65780601f10612ecb57610100808354040283529160200191612ef6565b820191906000526020600020905b815481529060010190602001808311612ed957829003601f168201915b505050918352505060048201546001600160a01b039081166020830152600583015416604082015260068201546060820152600782015460809091019060ff16600c811115612f4757612f47614845565b600c811115612f5857612f58614845565b81526020016007820160019054906101000a900460ff166004811115612f8057612f80614845565b6004811115612f9157612f91614845565b81526020016007820160029054906101000a900460ff166004811115612fb957612fb9614845565b6004811115612fca57612fca614845565b81526020016007820160039054906101000a900460ff166004811115612ff257612ff2614845565b600481111561300357613003614845565b81526008820154602082015260099091015460409091015292915050565b60008060015b6000548110156130915760026000828154811061304657613046614d57565b90600052602060002090600a020160070160019054906101000a900460ff16600481111561307657613076614845565b03613089578161308581615131565b9250505b600101613027565b50919050565b6004546001600160a01b031633146130c4573360405163394557ff60e21b815260040161085b9190614c1f565b60006130cf836134fb565b905060008082815481106130e5576130e5614d57565b600091825260208220600a91909102016005810180546001600160a01b0319166001600160a01b0387811691821790925542600984015560405192945092908716917fbe21ca43a190bc4f13c11b845e4461a995832bbdc47ce0a8dd4cba565fb238c49190a361294782600261374c565b6001600160a01b0316600090815260016020526040902054151590565b6005546001600160a01b031633146131a05733604051631ca0693360e11b815260040161085b9190614c1f565b60006131ab846134fb565b905060008360028111156131c1576131c1614845565b036131d5576131d08183613f28565b612947565b60018360028111156131e9576131e9614845565b036131f8576131d0818361407b565b61294781836141d1565b600061320c61432a565b546001600160a01b0316919050565b600061322561438c565b805490915060ff600160401b82041615906001600160401b031660008115801561324c5750825b90506000826001600160401b031660011480156132685750303b155b905081158015613276575080155b156132945760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156132be57845460ff60401b1916600160401b1785555b6132c66143b5565b6132cf866143bd565b6000805460010181558052831561332057845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b505050505050565b6004546001600160a01b03163314613355573360405163394557ff60e21b815260040161085b9190614c1f565b600481600481111561336957613369614845565b036133895780604051633d80b7c560e11b815260040161085b9190614878565b612076613395836134fb565b826143fe565b6133b1336000356001600160e01b031916613839565b6001600160a01b0381166133d857604051631c03d1ef60e11b815260040160405180910390fd5b600880546001600160a01b0319166001600160a01b0383811691821790925560055460045460065460075460405194959384169492841693600080516020615336833981519152936113849382169290911690614dc7565b60008061343c836134fb565b8154811061344c5761344c614d57565b60009182526020909120600a909102016007015462010000900460ff1692915050565b613477614785565b60006003836040516134899190614f50565b9081526020016040518091039020549050806000036129f25760405163b2be0e7f60e01b815260040160405180910390fd5b6004546001600160a01b031633146134e8573360405163394557ff60e21b815260040161085b9190614c1f565b6113cf6134f4826134fb565b600361353a565b6001600160a01b0381166000908152600160205260408120548082036135345760405163b2be0e7f60e01b815260040160405180910390fd5b92915050565b600080838154811061354e5761354e614d57565b60009182526020909120600a90910201600781015490915062010000900460ff1682600481111561358157613581614845565b81600481111561359357613593614845565b036135b157604051637d113f6d60e11b815260040160405180910390fd5b60078201805484919062ff00001916620100008360048111156135d6576135d6614845565b021790555042600983015560048201546040516001600160a01b03909116907f4f5d530726e69549f324811acc2350417dba052ef6a08d05f84ac97d4924bc66906136249084908790615105565b60405180910390a250505050565b6005546001600160a01b0316158061365457506005546001600160a01b03163b155b1561365d575050565b600554600854604080516303408e4760e41b815290516001600160a01b03938416936378dd4e6f938793911691633408e470916004808201926020929091908290030181865afa1580156136b5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906136d99190614dfe565b6040516001600160e01b031960e085901b16815260048101929092526024820152604401602060405180830381865afa925050508015613736575060408051601f3d908101601f19168201909252613733918101906152ad565b60015b156120765780156112ee576112ee82600161407b565b600080838154811061376057613760614d57565b60009182526020909120600a9091020160078101549091506301000000900460ff1682600481111561379457613794614845565b8160048111156137a6576137a6614845565b036137c457604051637d113f6d60e11b815260040160405180910390fd5b60078201805484919063ff000000191663010000008360048111156137eb576137eb614845565b021790555042600983015560048201546040516001600160a01b03909116907fd2ed21564a0773a7e31d3467b6565cd7085b6cb3985ddd93d4fc8bda98bdc68a906136249084908790615105565b600061384361432a565b80549091506001600160a01b031680613872576000604051638944034760e01b815260040161085b9190614c1f565b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa1580156138d6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138fa91906152ca565b92509250925082610c6f5780156139245760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff8216156139605760405163a426878960e01b81526001600160a01b038816600482015263ffffffff8316602482015260440161085b565b86604051632ecd3d0360e21b815260040161085b9190614c1f565b600080600073__$7c2d9d4b3468ae44947ccbf0e482c64ec5$__634e24dbc286866040518363ffffffff1660e01b81526004016139b9929190614de1565b600060405180830381865af41580156139d6573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526139fe9190810190614ec9565b9250925092506000600383604051613a169190614f50565b9081526020016040518091039020546000141590506000604051806101a00160405280898152602001868152602001858152602001848152602001886001600160a01b0316815260200160006001600160a01b03168152602001600860009054906101000a90046001600160a01b03166001600160a01b0316633408e4706040518163ffffffff1660e01b8152600401602060405180830381865afa158015613ac3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613ae79190614dfe565b815260200187600c811115613afe57613afe614845565b81526020016002815260200160028152602001600081524260208083018290526040909201528254600181810185556000948552938290208351600a90920201908155908201519192909190820190613b579082614fe8565b5060408201516002820190613b6c9082614fe8565b5060608201516003820190613b819082614fe8565b5060808201516004820180546001600160a01b039283166001600160a01b03199182161790915560a084015160058401805491909316911617905560c0820151600682015560e082015160078201805460ff1916600183600c811115613be957613be9614845565b0217905550610100828101516007830180549192909161ff00191690836004811115613c1757613c17614845565b021790555061012082015160078201805462ff0000191662010000836004811115613c4457613c44614845565b021790555061014082015160078201805463ff00000019166301000000836004811115613c7357613c73614845565b021790555061016082015160088201556101809091015160099091015560008054613ca0906001906150bd565b6001600160a01b03881660009081526001602090815260408083208490558b835260029091529020819055905081613cf45780600385604051613ce39190614f50565b908152604051908190036020019020555b866001600160a01b03167f063765b272ccfe8c5b747b73e618a2d710b3796d89b51946c128404ca6c85cc8868689604051613d31939291906150d0565b60405180910390a2866001600160a01b0316887f4e10bd9e15ecb70ffc648529318b2aa9c5fe9ca728e75df29cb427f4c49b22b588604051613d739190614e4a565b60405180910390a35050505050505050565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161480613df557507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316613de96144e5565b6001600160a01b031614155b15613e135760405163703e46dd60e11b815260040160405180910390fd5b565b6113cf336000356001600160e01b031916613839565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015613e85575060408051601f3d908101601f19168201909252613e8291810190614dfe565b60015b613ea45781604051634c9c8ce360e01b815260040161085b9190614c1f565b6000805160206153168339815191528114613ed557604051632a87526960e21b81526004810182905260240161085b565b6112ee83836144fb565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614613e135760405163703e46dd60e11b815260040160405180910390fd5b6000808381548110613f3c57613f3c614d57565b90600052602060002090600a020190508115613feb5760046007820154610100900460ff166004811115613f7257613f72614845565b03613f9057604051637d113f6d60e11b815260040160405180910390fd5b60078101546004808301546001600160a01b03166000908152600960205260409020805461010090930460ff1692909160ff19909116906001908490811115613fdb57613fdb614845565b02179055506112ee8360046143fe565b60046007820154610100900460ff16600481111561400b5761400b614845565b1461402957604051637d113f6d60e11b815260040160405180910390fd5b60048101546001600160a01b031660009081526009602052604090205461405490849060ff166143fe565b600401546001600160a01b03166000908152600960205260409020805460ff191690555050565b600080838154811061408f5761408f614d57565b90600052602060002090600a020190508115614140576004600782015462010000900460ff1660048111156140c6576140c6614845565b036140e457604051637d113f6d60e11b815260040160405180910390fd5b60078101546004808301546001600160a01b03166000908152600a6020526040902080546201000090930460ff1692909160ff1990911690600190849081111561413057614130614845565b02179055506112ee83600461353a565b6004600782015462010000900460ff16600481111561416157614161614845565b1461417f57604051637d113f6d60e11b815260040160405180910390fd5b60048101546001600160a01b03166000908152600a60205260409020546141aa90849060ff1661353a565b600401546001600160a01b03166000908152600a60205260409020805460ff191690555050565b60008083815481106141e5576141e5614d57565b90600052602060002090600a02019050811561429857600360078201546301000000900460ff16600481111561421d5761421d614845565b0361423b57604051637d113f6d60e11b815260040160405180910390fd5b60078101546004808301546001600160a01b03166000908152600b602052604090208054630100000090930460ff1692909160ff1990911690600190849081111561428857614288614845565b02179055506112ee83600361374c565b600360078201546301000000900460ff1660048111156142ba576142ba614845565b146142d857604051637d113f6d60e11b815260040160405180910390fd5b60048101546001600160a01b03166000908152600b602052604090205461430390849060ff1661374c565b600401546001600160a01b03166000908152600b60205260409020805460ff191690555050565b60008060ff1961435b60017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f356150bd565b60405160200161436d91815260200190565b60408051601f1981840301815291905280516020909101201692915050565b6000807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00613534565b613e13614551565b60006143c761432a565b80549091506001600160a01b0316156143f55781604051638944034760e01b815260040161085b9190614c1f565b61207682614576565b600080838154811061441257614412614d57565b60009182526020909120600a909102016007810154909150610100900460ff1682600481111561444457614444614845565b81600481111561445657614456614845565b0361447457604051637d113f6d60e11b815260040160405180910390fd5b60078201805484919061ff00191661010083600481111561449757614497614845565b021790555042600983015560048201546040516001600160a01b03909116907f962ec3691afeb0804244549c8e163d9e4a89f507a6ed4a8a87c6fbca721699d2906136249084908790615105565b600060008051602061531683398151915261320c565b61450482614606565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a2805115614549576112ee8282614662565b6120766146d8565b6145596146f7565b613e1357604051631afcd79f60e31b815260040160405180910390fd5b6001600160a01b03811661459f5780604051638944034760e01b815260040161085b9190614c1f565b60006145a961432a565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b806001600160a01b03163b6000036146335780604051634c9c8ce360e01b815260040161085b9190614c1f565b60008051602061531683398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080846001600160a01b03168460405161467f9190614f50565b600060405180830381855af49150503d80600081146146ba576040519150601f19603f3d011682016040523d82523d6000602084013e6146bf565b606091505b50915091506146cf858383614711565b95945050505050565b3415613e135760405163b398979f60e01b815260040160405180910390fd5b600061470161438c565b54600160401b900460ff16919050565b606082614726576147218261475d565b6107e7565b815115801561473d57506001600160a01b0384163b155b156117b75783604051639996b31560e01b815260040161085b9190614c1f565b80511561476c57805160208201fd5b60405163d6bda27560e01b815260040160405180910390fd5b604051806101a001604052806000801916815260200160608152602001606081526020016060815260200160006001600160a01b0316815260200160006001600160a01b03168152602001600081526020016000600c8111156147ea576147ea614845565b815260200160008152602001600081526020016000815260200160008152602001600081525090565b6001600160a01b03811681146113cf57600080fd5b60006020828403121561483a57600080fd5b81356107e781614813565b634e487b7160e01b600052602160045260246000fd5b600581106113cf576113cf614845565b6148748161485b565b9052565b602081016148858361485b565b91905290565b6000806000606084860312156148a057600080fd5b8335925060208401356148b281614813565b9150604084013560ff811681146148c857600080fd5b809150509250925092565b6000602082840312156148e557600080fd5b5035919050565b60005b838110156149075781810151838201526020016148ef565b50506000910152565b600081518084526149288160208601602086016148ec565b601f01601f19169290920160200192915050565b600d811061487457614874614845565b60006101a082518452602083015181602086015261496c82860182614910565b915050604083015184820360408601526149868282614910565b915050606083015184820360608601526149a08282614910565b91505060808301516149bd60808601826001600160a01b03169052565b5060a08301516149d860a08601826001600160a01b03169052565b5060c083015160c085015260e08301516149f560e086018261493c565b5061010080840151614a098287018261486b565b505061012080840151614a1e8287018261486b565b505061014080840151614a338287018261486b565b505061016083810151908501526101809283015192909301919091525090565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b82811015614aaa57603f19888603018452614a9885835161494c565b94509285019290850190600101614a7c565b5092979650505050505050565b600581106113cf57600080fd5b600060208284031215614ad657600080fd5b81356107e781614ab7565b60008060408385031215614af457600080fd5b823591506020830135614b0681614813565b809150509250929050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b0381118282101715614b4f57614b4f614b11565b604052919050565b60006001600160401b03821115614b7057614b70614b11565b50601f01601f191660200190565b6000614b91614b8c84614b57565b614b27565b9050828152838383011115614ba557600080fd5b828260208301376000602084830101529392505050565b60008060408385031215614bcf57600080fd5b8235614bda81614813565b915060208301356001600160401b03811115614bf557600080fd5b8301601f81018513614c0657600080fd5b614c1585823560208401614b7e565b9150509250929050565b6001600160a01b0391909116815260200190565b6020815260006107e7602083018461494c565b60008060408385031215614c5957600080fd5b8235614c6481614813565b91506020830135614b0681614813565b6020815260006107e76020830184614910565b80151581146113cf57600080fd5b600080600060608486031215614caa57600080fd5b8335614cb581614813565b9250602084013560038110614cc957600080fd5b915060408401356148c881614c87565b60008060408385031215614cec57600080fd5b8235614cf781614813565b91506020830135614b0681614ab7565b600060208284031215614d1957600080fd5b81356001600160401b03811115614d2f57600080fd5b8201601f81018413614d4057600080fd5b614d4f84823560208401614b7e565b949350505050565b634e487b7160e01b600052603260045260246000fd5b6001600160a01b038316815260408101614d868361485b565b8260208301529392505050565b9182526001600160a01b0316602082015260400190565b600060208284031215614dbc57600080fd5b81516107e781614813565b6001600160a01b0392831681529116602082015260400190565b6001600160a01b0383168152604081016107e7602083018461493c565b600060208284031215614e1057600080fd5b5051919050565b6001600160a01b03851681526020810184905260808101614e3b604083018561493c565b82606083015295945050505050565b60208101613534828461493c565b600060208284031215614e6a57600080fd5b8151600d81106107e757600080fd5b6000614e87614b8c84614b57565b9050828152838383011115614e9b57600080fd5b6107e78360208301846148ec565b600082601f830112614eba57600080fd5b6107e783835160208501614e79565b600080600060608486031215614ede57600080fd5b83516001600160401b0380821115614ef557600080fd5b614f0187838801614ea9565b94506020860151915080821115614f1757600080fd5b614f2387838801614ea9565b93506040860151915080821115614f3957600080fd5b50614f4686828701614ea9565b9150509250925092565b60008251614f628184602087016148ec565b9190910192915050565b600181811c90821680614f8057607f821691505b60208210810361309157634e487b7160e01b600052602260045260246000fd5b601f8211156112ee576000816000526020600020601f850160051c81016020861015614fc95750805b601f850160051c820191505b8181101561332057828155600101614fd5565b81516001600160401b0381111561500157615001614b11565b6150158161500f8454614f6c565b84614fa0565b602080601f83116001811461504a57600084156150325750858301515b600019600386901b1c1916600185901b178555613320565b600085815260208120601f198616915b828110156150795788860151825594840194600190910190840161505a565b50858210156150975787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b634e487b7160e01b600052601160045260246000fd5b81810381811115613534576135346150a7565b6060815260006150e36060830186614910565b82810360208401526150f58186614910565b915050614d4f604083018461493c565b604081016151128461485b565b838252614d868361485b565b80820180821115613534576135346150a7565b600060018201615143576151436150a7565b5060010190565b6000815461515781614f6c565b808552602060018381168015615174576001811461518e576151bc565b60ff1985168884015283151560051b8801830195506151bc565b866000528260002060005b858110156151b45781548a8201860152908301908401615199565b890184019650505b505050505092915050565b6001600160a01b0388811682526000906151e4602084018a61493c565b80881660408401525085606083015260e0608083015261520760e083018661514a565b82810360a0840152615219818661514a565b905082810360c084015261522d818561514a565b9a9950505050505050505050565b60006020828403121561524d57600080fd5b81516001600160401b0381111561526357600080fd5b8201601f8101841361527457600080fd5b614d4f84825160208401614e79565b8381526001600160a01b03831660208201526060604082018190526000906146cf90830184614910565b6000602082840312156152bf57600080fd5b81516107e781614c87565b6000806000606084860312156152df57600080fd5b83516152ea81614c87565b602085015190935063ffffffff8116811461530457600080fd5b60408501519092506148c881614c8756fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc187af05701e638e7e5f432ce06f1b41424d001704ef286fd435865056ec8542fa264697066735822122042814595737acec10a092212244c90975bb80295731b9672978348863e69b74f64736f6c63430008180033",
	Deps: []*bind.MetaData{
		&PNTokenCoreLibMetaData,
	},
}

// PNTokenCoreV1 is an auto generated Go binding around an Ethereum contract.
type PNTokenCoreV1 struct {
	abi abi.ABI
}

// NewPNTokenCoreV1 creates a new instance of PNTokenCoreV1.
func NewPNTokenCoreV1() *PNTokenCoreV1 {
	parsed, err := PNTokenCoreV1MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &PNTokenCoreV1{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *PNTokenCoreV1) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackUPGRADEINTERFACEVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (pNTokenCoreV1 *PNTokenCoreV1) PackUPGRADEINTERFACEVERSION() []byte {
	enc, err := pNTokenCoreV1.abi.Pack("UPGRADE_INTERFACE_VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUPGRADEINTERFACEVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackUPGRADEINTERFACEVERSION(data []byte) (string, error) {
	out, err := pNTokenCoreV1.abi.Unpack("UPGRADE_INTERFACE_VERSION", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackActivateToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x060ed143.
//
// Solidity: function activateToken(bytes32 resourceId, address tokenAddress, uint8 ercStandard) returns()
func (pNTokenCoreV1 *PNTokenCoreV1) PackActivateToken(resourceId [32]byte, tokenAddress common.Address, ercStandard uint8) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("activateToken", resourceId, tokenAddress, ercStandard)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (pNTokenCoreV1 *PNTokenCoreV1) PackAuthority() []byte {
	enc, err := pNTokenCoreV1.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := pNTokenCoreV1.abi.Unpack("authority", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDeprecateOnPublicChain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1fe56159.
//
// Solidity: function deprecateOnPublicChain(address tokenAddress) returns()
func (pNTokenCoreV1 *PNTokenCoreV1) PackDeprecateOnPublicChain(tokenAddress common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("deprecateOnPublicChain", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackEndpoint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (pNTokenCoreV1 *PNTokenCoreV1) PackEndpoint() []byte {
	enc, err := pNTokenCoreV1.abi.Pack("endpoint")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackEndpoint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackEndpoint(data []byte) (common.Address, error) {
	out, err := pNTokenCoreV1.abi.Unpack("endpoint", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackEnygmaPNEvents is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdff11d4c.
//
// Solidity: function enygmaPNEvents() view returns(address)
func (pNTokenCoreV1 *PNTokenCoreV1) PackEnygmaPNEvents() []byte {
	enc, err := pNTokenCoreV1.abi.Pack("enygmaPNEvents")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackEnygmaPNEvents is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdff11d4c.
//
// Solidity: function enygmaPNEvents() view returns(address)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackEnygmaPNEvents(data []byte) (common.Address, error) {
	out, err := pNTokenCoreV1.abi.Unpack("enygmaPNEvents", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetActiveTokenCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x93d85727.
//
// Solidity: function getActiveTokenCount() view returns(uint256)
func (pNTokenCoreV1 *PNTokenCoreV1) PackGetActiveTokenCount() []byte {
	enc, err := pNTokenCoreV1.abi.Pack("getActiveTokenCount")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetActiveTokenCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x93d85727.
//
// Solidity: function getActiveTokenCount() view returns(uint256)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackGetActiveTokenCount(data []byte) (*big.Int, error) {
	out, err := pNTokenCoreV1.abi.Unpack("getActiveTokenCount", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetAllTokens is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2a5c792a.
//
// Solidity: function getAllTokens() view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256)[])
func (pNTokenCoreV1 *PNTokenCoreV1) PackGetAllTokens() []byte {
	enc, err := pNTokenCoreV1.abi.Pack("getAllTokens")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAllTokens is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2a5c792a.
//
// Solidity: function getAllTokens() view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256)[])
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackGetAllTokens(data []byte) ([]TokenStructsToken, error) {
	out, err := pNTokenCoreV1.abi.Unpack("getAllTokens", data)
	if err != nil {
		return *new([]TokenStructsToken), err
	}
	out0 := *abi.ConvertType(out[0], new([]TokenStructsToken)).(*[]TokenStructsToken)
	return out0, err
}

// PackGetHubStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe3ade0ec.
//
// Solidity: function getHubStatus(address tokenAddress) view returns(uint8)
func (pNTokenCoreV1 *PNTokenCoreV1) PackGetHubStatus(tokenAddress common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("getHubStatus", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetHubStatus is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe3ade0ec.
//
// Solidity: function getHubStatus(address tokenAddress) view returns(uint8)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackGetHubStatus(data []byte) (uint8, error) {
	out, err := pNTokenCoreV1.abi.Unpack("getHubStatus", data)
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
func (pNTokenCoreV1 *PNTokenCoreV1) PackGetPrivacyNodeStatus(tokenAddress common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("getPrivacyNodeStatus", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetPrivacyNodeStatus is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x22e594ff.
//
// Solidity: function getPrivacyNodeStatus(address tokenAddress) view returns(uint8)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackGetPrivacyNodeStatus(data []byte) (uint8, error) {
	out, err := pNTokenCoreV1.abi.Unpack("getPrivacyNodeStatus", data)
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
func (pNTokenCoreV1 *PNTokenCoreV1) PackGetPublicChainStatus(tokenAddress common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("getPublicChainStatus", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetPublicChainStatus is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01162d43.
//
// Solidity: function getPublicChainStatus(address tokenAddress) view returns(uint8)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackGetPublicChainStatus(data []byte) (uint8, error) {
	out, err := pNTokenCoreV1.abi.Unpack("getPublicChainStatus", data)
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
func (pNTokenCoreV1 *PNTokenCoreV1) PackGetTokenByAddress(tokenAddress common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("getTokenByAddress", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTokenByAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x91ded8fa.
//
// Solidity: function getTokenByAddress(address tokenAddress) view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256))
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackGetTokenByAddress(data []byte) (TokenStructsToken, error) {
	out, err := pNTokenCoreV1.abi.Unpack("getTokenByAddress", data)
	if err != nil {
		return *new(TokenStructsToken), err
	}
	out0 := *abi.ConvertType(out[0], new(TokenStructsToken)).(*TokenStructsToken)
	return out0, err
}

// PackGetTokenByResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8a58b79f.
//
// Solidity: function getTokenByResourceId(bytes32 resourceId) view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256))
func (pNTokenCoreV1 *PNTokenCoreV1) PackGetTokenByResourceId(resourceId [32]byte) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("getTokenByResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTokenByResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8a58b79f.
//
// Solidity: function getTokenByResourceId(bytes32 resourceId) view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256))
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackGetTokenByResourceId(data []byte) (TokenStructsToken, error) {
	out, err := pNTokenCoreV1.abi.Unpack("getTokenByResourceId", data)
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
func (pNTokenCoreV1 *PNTokenCoreV1) PackGetTokenBySymbol(symbol string) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("getTokenBySymbol", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTokenBySymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xefa74f1f.
//
// Solidity: function getTokenBySymbol(string symbol) view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256))
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackGetTokenBySymbol(data []byte) (TokenStructsToken, error) {
	out, err := pNTokenCoreV1.abi.Unpack("getTokenBySymbol", data)
	if err != nil {
		return *new(TokenStructsToken), err
	}
	out0 := *abi.ConvertType(out[0], new(TokenStructsToken)).(*TokenStructsToken)
	return out0, err
}

// PackGetTokenCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x78a89567.
//
// Solidity: function getTokenCount() view returns(uint256)
func (pNTokenCoreV1 *PNTokenCoreV1) PackGetTokenCount() []byte {
	enc, err := pNTokenCoreV1.abi.Pack("getTokenCount")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTokenCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x78a89567.
//
// Solidity: function getTokenCount() view returns(uint256)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackGetTokenCount(data []byte) (*big.Int, error) {
	out, err := pNTokenCoreV1.abi.Unpack("getTokenCount", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetTokensByHubStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6f0656f0.
//
// Solidity: function getTokensByHubStatus(uint8 status) view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256)[])
func (pNTokenCoreV1 *PNTokenCoreV1) PackGetTokensByHubStatus(status uint8) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("getTokensByHubStatus", status)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTokensByHubStatus is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6f0656f0.
//
// Solidity: function getTokensByHubStatus(uint8 status) view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256)[])
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackGetTokensByHubStatus(data []byte) ([]TokenStructsToken, error) {
	out, err := pNTokenCoreV1.abi.Unpack("getTokensByHubStatus", data)
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
func (pNTokenCoreV1 *PNTokenCoreV1) PackGetTokensByPrivacyNodeStatus(status uint8) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("getTokensByPrivacyNodeStatus", status)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTokensByPrivacyNodeStatus is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3b937349.
//
// Solidity: function getTokensByPrivacyNodeStatus(uint8 status) view returns((bytes32,string,string,string,address,address,uint256,uint8,uint8,uint8,uint8,uint256,uint256)[])
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackGetTokensByPrivacyNodeStatus(data []byte) ([]TokenStructsToken, error) {
	out, err := pNTokenCoreV1.abi.Unpack("getTokensByPrivacyNodeStatus", data)
	if err != nil {
		return *new([]TokenStructsToken), err
	}
	out0 := *abi.ConvertType(out[0], new([]TokenStructsToken)).(*[]TokenStructsToken)
	return out0, err
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc4d66de8.
//
// Solidity: function initialize(address authority_) returns()
func (pNTokenCoreV1 *PNTokenCoreV1) PackInitialize(authority common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("initialize", authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackIsTokenActiveForHub is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0014bc46.
//
// Solidity: function isTokenActiveForHub(address tokenAddress) view returns(bool)
func (pNTokenCoreV1 *PNTokenCoreV1) PackIsTokenActiveForHub(tokenAddress common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("isTokenActiveForHub", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsTokenActiveForHub is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0014bc46.
//
// Solidity: function isTokenActiveForHub(address tokenAddress) view returns(bool)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackIsTokenActiveForHub(data []byte) (bool, error) {
	out, err := pNTokenCoreV1.abi.Unpack("isTokenActiveForHub", data)
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
func (pNTokenCoreV1 *PNTokenCoreV1) PackIsTokenActiveForPublicChain(tokenAddress common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("isTokenActiveForPublicChain", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsTokenActiveForPublicChain is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x64c0673a.
//
// Solidity: function isTokenActiveForPublicChain(address tokenAddress) view returns(bool)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackIsTokenActiveForPublicChain(data []byte) (bool, error) {
	out, err := pNTokenCoreV1.abi.Unpack("isTokenActiveForPublicChain", data)
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
func (pNTokenCoreV1 *PNTokenCoreV1) PackIsTokenFullyOperational(tokenAddress common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("isTokenFullyOperational", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsTokenFullyOperational is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x45741567.
//
// Solidity: function isTokenFullyOperational(address tokenAddress) view returns(bool)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackIsTokenFullyOperational(data []byte) (bool, error) {
	out, err := pNTokenCoreV1.abi.Unpack("isTokenFullyOperational", data)
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
func (pNTokenCoreV1 *PNTokenCoreV1) PackProxiableUUID() []byte {
	enc, err := pNTokenCoreV1.abi.Pack("proxiableUUID")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackProxiableUUID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackProxiableUUID(data []byte) ([32]byte, error) {
	out, err := pNTokenCoreV1.abi.Unpack("proxiableUUID", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackRegisterHubToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x47ecc939.
//
// Solidity: function registerHubToken(bytes32 resourceId, address tokenAddress) returns()
func (pNTokenCoreV1 *PNTokenCoreV1) PackRegisterHubToken(resourceId [32]byte, tokenAddress common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("registerHubToken", resourceId, tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRegisterToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x09824a80.
//
// Solidity: function registerToken(address tokenAddress) returns()
func (pNTokenCoreV1 *PNTokenCoreV1) PackRegisterToken(tokenAddress common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("registerToken", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRejectToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf05ab9de.
//
// Solidity: function rejectToken(address tokenAddress) returns()
func (pNTokenCoreV1 *PNTokenCoreV1) PackRejectToken(tokenAddress common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("rejectToken", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRelayAuthorizationRegistry is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcd4b4a51.
//
// Solidity: function relayAuthorizationRegistry() view returns(address)
func (pNTokenCoreV1 *PNTokenCoreV1) PackRelayAuthorizationRegistry() []byte {
	enc, err := pNTokenCoreV1.abi.Pack("relayAuthorizationRegistry")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackRelayAuthorizationRegistry is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcd4b4a51.
//
// Solidity: function relayAuthorizationRegistry() view returns(address)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackRelayAuthorizationRegistry(data []byte) (common.Address, error) {
	out, err := pNTokenCoreV1.abi.Unpack("relayAuthorizationRegistry", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackSetEndpoint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdbbb4155.
//
// Solidity: function setEndpoint(address endpointAddress) returns()
func (pNTokenCoreV1 *PNTokenCoreV1) PackSetEndpoint(endpointAddress common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("setEndpoint", endpointAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetEnygmaPNEvents is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3cd17302.
//
// Solidity: function setEnygmaPNEvents(address enygmaPNEventsAddress) returns()
func (pNTokenCoreV1 *PNTokenCoreV1) PackSetEnygmaPNEvents(enygmaPNEventsAddress common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("setEnygmaPNEvents", enygmaPNEventsAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetFreezeStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb58bdb49.
//
// Solidity: function setFreezeStatus(address tokenAddress, uint8 layer, bool frozen) returns()
func (pNTokenCoreV1 *PNTokenCoreV1) PackSetFreezeStatus(tokenAddress common.Address, layer uint8, frozen bool) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("setFreezeStatus", tokenAddress, layer, frozen)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetRelayAuthorizationRegistry is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8212f073.
//
// Solidity: function setRelayAuthorizationRegistry(address relayAuthRegistry) returns()
func (pNTokenCoreV1 *PNTokenCoreV1) PackSetRelayAuthorizationRegistry(relayAuthRegistry common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("setRelayAuthorizationRegistry", relayAuthRegistry)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetTokenFreezeManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x13aa1a1d.
//
// Solidity: function setTokenFreezeManager(address freezeManager) returns()
func (pNTokenCoreV1 *PNTokenCoreV1) PackSetTokenFreezeManager(freezeManager common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("setTokenFreezeManager", freezeManager)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetTokenRegistry is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x35a5af92.
//
// Solidity: function setTokenRegistry(address tokenRegistry) returns()
func (pNTokenCoreV1 *PNTokenCoreV1) PackSetTokenRegistry(tokenRegistry common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("setTokenRegistry", tokenRegistry)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSubmitToHub is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7fbc532d.
//
// Solidity: function submitToHub(address tokenAddress) returns()
func (pNTokenCoreV1 *PNTokenCoreV1) PackSubmitToHub(tokenAddress common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("submitToHub", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSubmitToPublicChain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0a75cf4e.
//
// Solidity: function submitToPublicChain(address tokenAddress) returns()
func (pNTokenCoreV1 *PNTokenCoreV1) PackSubmitToPublicChain(tokenAddress common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("submitToPublicChain", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackTokenExists is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb33f78ca.
//
// Solidity: function tokenExists(address tokenAddress) view returns(bool)
func (pNTokenCoreV1 *PNTokenCoreV1) PackTokenExists(tokenAddress common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("tokenExists", tokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTokenExists is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb33f78ca.
//
// Solidity: function tokenExists(address tokenAddress) view returns(bool)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenExists(data []byte) (bool, error) {
	out, err := pNTokenCoreV1.abi.Unpack("tokenExists", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackTokenExistsByResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1047100b.
//
// Solidity: function tokenExistsByResourceId(bytes32 resourceId) view returns(bool)
func (pNTokenCoreV1 *PNTokenCoreV1) PackTokenExistsByResourceId(resourceId [32]byte) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("tokenExistsByResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTokenExistsByResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1047100b.
//
// Solidity: function tokenExistsByResourceId(bytes32 resourceId) view returns(bool)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenExistsByResourceId(data []byte) (bool, error) {
	out, err := pNTokenCoreV1.abi.Unpack("tokenExistsByResourceId", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackTokenFreezeManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9d53ea4c.
//
// Solidity: function tokenFreezeManager() view returns(address)
func (pNTokenCoreV1 *PNTokenCoreV1) PackTokenFreezeManager() []byte {
	enc, err := pNTokenCoreV1.abi.Pack("tokenFreezeManager")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTokenFreezeManager is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9d53ea4c.
//
// Solidity: function tokenFreezeManager() view returns(address)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenFreezeManager(data []byte) (common.Address, error) {
	out, err := pNTokenCoreV1.abi.Unpack("tokenFreezeManager", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackTokenRegistryAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5be2aca0.
//
// Solidity: function tokenRegistryAddress() view returns(address)
func (pNTokenCoreV1 *PNTokenCoreV1) PackTokenRegistryAddress() []byte {
	enc, err := pNTokenCoreV1.abi.Pack("tokenRegistryAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTokenRegistryAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5be2aca0.
//
// Solidity: function tokenRegistryAddress() view returns(address)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenRegistryAddress(data []byte) (common.Address, error) {
	out, err := pNTokenCoreV1.abi.Unpack("tokenRegistryAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackUpdatePrivacyNodeStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd8dc0510.
//
// Solidity: function updatePrivacyNodeStatus(address tokenAddress, uint8 status) returns()
func (pNTokenCoreV1 *PNTokenCoreV1) PackUpdatePrivacyNodeStatus(tokenAddress common.Address, status uint8) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("updatePrivacyNodeStatus", tokenAddress, status)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpdatePublicTokenAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa0b9b9b1.
//
// Solidity: function updatePublicTokenAddress(address tokenAddress, address publicTokenAddress) returns()
func (pNTokenCoreV1 *PNTokenCoreV1) PackUpdatePublicTokenAddress(tokenAddress common.Address, publicTokenAddress common.Address) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("updatePublicTokenAddress", tokenAddress, publicTokenAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpgradeToAndCall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (pNTokenCoreV1 *PNTokenCoreV1) PackUpgradeToAndCall(newImplementation common.Address, data []byte) []byte {
	enc, err := pNTokenCoreV1.abi.Pack("upgradeToAndCall", newImplementation, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// PNTokenCoreV1AuthorityUpdated represents a AuthorityUpdated event raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1AuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const PNTokenCoreV1AuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (PNTokenCoreV1AuthorityUpdated) ContractEventName() string {
	return PNTokenCoreV1AuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackAuthorityUpdatedEvent(log *types.Log) (*PNTokenCoreV1AuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != pNTokenCoreV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNTokenCoreV1AuthorityUpdated)
	if len(log.Data) > 0 {
		if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNTokenCoreV1.abi.Events[event].Inputs {
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

// PNTokenCoreV1HubStatusUpdated represents a HubStatusUpdated event raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1HubStatusUpdated struct {
	TokenAddress   common.Address
	PreviousStatus uint8
	NewStatus      uint8
	Raw            *types.Log // Blockchain specific contextual infos
}

const PNTokenCoreV1HubStatusUpdatedEventName = "HubStatusUpdated"

// ContractEventName returns the user-defined event name.
func (PNTokenCoreV1HubStatusUpdated) ContractEventName() string {
	return PNTokenCoreV1HubStatusUpdatedEventName
}

// UnpackHubStatusUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event HubStatusUpdated(address indexed tokenAddress, uint8 previousStatus, uint8 newStatus)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackHubStatusUpdatedEvent(log *types.Log) (*PNTokenCoreV1HubStatusUpdated, error) {
	event := "HubStatusUpdated"
	if log.Topics[0] != pNTokenCoreV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNTokenCoreV1HubStatusUpdated)
	if len(log.Data) > 0 {
		if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNTokenCoreV1.abi.Events[event].Inputs {
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

// PNTokenCoreV1Initialized represents a Initialized event raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1Initialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const PNTokenCoreV1InitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (PNTokenCoreV1Initialized) ContractEventName() string {
	return PNTokenCoreV1InitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackInitializedEvent(log *types.Log) (*PNTokenCoreV1Initialized, error) {
	event := "Initialized"
	if log.Topics[0] != pNTokenCoreV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNTokenCoreV1Initialized)
	if len(log.Data) > 0 {
		if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNTokenCoreV1.abi.Events[event].Inputs {
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

// PNTokenCoreV1PrivacyNodeStatusUpdated represents a PrivacyNodeStatusUpdated event raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1PrivacyNodeStatusUpdated struct {
	TokenAddress   common.Address
	PreviousStatus uint8
	NewStatus      uint8
	Raw            *types.Log // Blockchain specific contextual infos
}

const PNTokenCoreV1PrivacyNodeStatusUpdatedEventName = "PrivacyNodeStatusUpdated"

// ContractEventName returns the user-defined event name.
func (PNTokenCoreV1PrivacyNodeStatusUpdated) ContractEventName() string {
	return PNTokenCoreV1PrivacyNodeStatusUpdatedEventName
}

// UnpackPrivacyNodeStatusUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PrivacyNodeStatusUpdated(address indexed tokenAddress, uint8 previousStatus, uint8 newStatus)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackPrivacyNodeStatusUpdatedEvent(log *types.Log) (*PNTokenCoreV1PrivacyNodeStatusUpdated, error) {
	event := "PrivacyNodeStatusUpdated"
	if log.Topics[0] != pNTokenCoreV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNTokenCoreV1PrivacyNodeStatusUpdated)
	if len(log.Data) > 0 {
		if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNTokenCoreV1.abi.Events[event].Inputs {
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

// PNTokenCoreV1PublicChainStatusUpdated represents a PublicChainStatusUpdated event raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1PublicChainStatusUpdated struct {
	TokenAddress   common.Address
	PreviousStatus uint8
	NewStatus      uint8
	Raw            *types.Log // Blockchain specific contextual infos
}

const PNTokenCoreV1PublicChainStatusUpdatedEventName = "PublicChainStatusUpdated"

// ContractEventName returns the user-defined event name.
func (PNTokenCoreV1PublicChainStatusUpdated) ContractEventName() string {
	return PNTokenCoreV1PublicChainStatusUpdatedEventName
}

// UnpackPublicChainStatusUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PublicChainStatusUpdated(address indexed tokenAddress, uint8 previousStatus, uint8 newStatus)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackPublicChainStatusUpdatedEvent(log *types.Log) (*PNTokenCoreV1PublicChainStatusUpdated, error) {
	event := "PublicChainStatusUpdated"
	if log.Topics[0] != pNTokenCoreV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNTokenCoreV1PublicChainStatusUpdated)
	if len(log.Data) > 0 {
		if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNTokenCoreV1.abi.Events[event].Inputs {
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

// PNTokenCoreV1PublicTokenAddressUpdated represents a PublicTokenAddressUpdated event raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1PublicTokenAddressUpdated struct {
	TokenAddress       common.Address
	PublicTokenAddress common.Address
	Raw                *types.Log // Blockchain specific contextual infos
}

const PNTokenCoreV1PublicTokenAddressUpdatedEventName = "PublicTokenAddressUpdated"

// ContractEventName returns the user-defined event name.
func (PNTokenCoreV1PublicTokenAddressUpdated) ContractEventName() string {
	return PNTokenCoreV1PublicTokenAddressUpdatedEventName
}

// UnpackPublicTokenAddressUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PublicTokenAddressUpdated(address indexed tokenAddress, address indexed publicTokenAddress)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackPublicTokenAddressUpdatedEvent(log *types.Log) (*PNTokenCoreV1PublicTokenAddressUpdated, error) {
	event := "PublicTokenAddressUpdated"
	if log.Topics[0] != pNTokenCoreV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNTokenCoreV1PublicTokenAddressUpdated)
	if len(log.Data) > 0 {
		if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNTokenCoreV1.abi.Events[event].Inputs {
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

// PNTokenCoreV1TokenActivated represents a TokenActivated event raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenActivated struct {
	ResourceId   [32]byte
	TokenAddress common.Address
	ErcStandard  uint8
	Raw          *types.Log // Blockchain specific contextual infos
}

const PNTokenCoreV1TokenActivatedEventName = "TokenActivated"

// ContractEventName returns the user-defined event name.
func (PNTokenCoreV1TokenActivated) ContractEventName() string {
	return PNTokenCoreV1TokenActivatedEventName
}

// UnpackTokenActivatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenActivated(bytes32 indexed resourceId, address indexed tokenAddress, uint8 ercStandard)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenActivatedEvent(log *types.Log) (*PNTokenCoreV1TokenActivated, error) {
	event := "TokenActivated"
	if log.Topics[0] != pNTokenCoreV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNTokenCoreV1TokenActivated)
	if len(log.Data) > 0 {
		if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNTokenCoreV1.abi.Events[event].Inputs {
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

// PNTokenCoreV1TokenCoreModulesConfigured represents a TokenCoreModulesConfigured event raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenCoreModulesConfigured struct {
	TokenRegistry              common.Address
	TokenFreezeManager         common.Address
	Endpoint                   common.Address
	EnygmaPNEvents             common.Address
	RelayAuthorizationRegistry common.Address
	Raw                        *types.Log // Blockchain specific contextual infos
}

const PNTokenCoreV1TokenCoreModulesConfiguredEventName = "TokenCoreModulesConfigured"

// ContractEventName returns the user-defined event name.
func (PNTokenCoreV1TokenCoreModulesConfigured) ContractEventName() string {
	return PNTokenCoreV1TokenCoreModulesConfiguredEventName
}

// UnpackTokenCoreModulesConfiguredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenCoreModulesConfigured(address indexed tokenRegistry, address indexed tokenFreezeManager, address indexed endpoint, address enygmaPNEvents, address relayAuthorizationRegistry)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenCoreModulesConfiguredEvent(log *types.Log) (*PNTokenCoreV1TokenCoreModulesConfigured, error) {
	event := "TokenCoreModulesConfigured"
	if log.Topics[0] != pNTokenCoreV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNTokenCoreV1TokenCoreModulesConfigured)
	if len(log.Data) > 0 {
		if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNTokenCoreV1.abi.Events[event].Inputs {
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

// PNTokenCoreV1TokenRegistered represents a TokenRegistered event raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenRegistered struct {
	TokenAddress common.Address
	Name         string
	Symbol       string
	ErcStandard  uint8
	Raw          *types.Log // Blockchain specific contextual infos
}

const PNTokenCoreV1TokenRegisteredEventName = "TokenRegistered"

// ContractEventName returns the user-defined event name.
func (PNTokenCoreV1TokenRegistered) ContractEventName() string {
	return PNTokenCoreV1TokenRegisteredEventName
}

// UnpackTokenRegisteredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenRegistered(address indexed tokenAddress, string name, string symbol, uint8 ercStandard)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenRegisteredEvent(log *types.Log) (*PNTokenCoreV1TokenRegistered, error) {
	event := "TokenRegistered"
	if log.Topics[0] != pNTokenCoreV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNTokenCoreV1TokenRegistered)
	if len(log.Data) > 0 {
		if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNTokenCoreV1.abi.Events[event].Inputs {
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

// PNTokenCoreV1Upgraded represents a Upgraded event raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1Upgraded struct {
	Implementation common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const PNTokenCoreV1UpgradedEventName = "Upgraded"

// ContractEventName returns the user-defined event name.
func (PNTokenCoreV1Upgraded) ContractEventName() string {
	return PNTokenCoreV1UpgradedEventName
}

// UnpackUpgradedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Upgraded(address indexed implementation)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackUpgradedEvent(log *types.Log) (*PNTokenCoreV1Upgraded, error) {
	event := "Upgraded"
	if log.Topics[0] != pNTokenCoreV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNTokenCoreV1Upgraded)
	if len(log.Data) > 0 {
		if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNTokenCoreV1.abi.Events[event].Inputs {
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
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["AddressEmptyCode"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackAddressEmptyCodeError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["FailedCall"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackFailedCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["TokenCoreV1EndpointNotConfigured"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackTokenCoreV1EndpointNotConfiguredError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["TokenCoreV1HubApprovalNotPending"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackTokenCoreV1HubApprovalNotPendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["TokenCoreV1InvalidEndpointAddress"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackTokenCoreV1InvalidEndpointAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["TokenCoreV1InvalidEnygmaPNEventsAddress"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackTokenCoreV1InvalidEnygmaPNEventsAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["TokenCoreV1InvalidPrivacyNodeStatus"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackTokenCoreV1InvalidPrivacyNodeStatusError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["TokenCoreV1InvalidTokenAddress"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackTokenCoreV1InvalidTokenAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["TokenCoreV1InvalidTokenFreezeManagerAddress"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackTokenCoreV1InvalidTokenFreezeManagerAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["TokenCoreV1InvalidTokenRegistryAddress"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackTokenCoreV1InvalidTokenRegistryAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["TokenCoreV1PrivacyNodeAuthorizationRequired"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackTokenCoreV1PrivacyNodeAuthorizationRequiredError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["TokenCoreV1ResourceIdAlreadyAssigned"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackTokenCoreV1ResourceIdAlreadyAssignedError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["TokenCoreV1ResourceIdAlreadySet"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackTokenCoreV1ResourceIdAlreadySetError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["TokenCoreV1StatusAlreadySet"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackTokenCoreV1StatusAlreadySetError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["TokenCoreV1TokenAlreadyExists"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackTokenCoreV1TokenAlreadyExistsError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["TokenCoreV1TokenDoesNotExist"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackTokenCoreV1TokenDoesNotExistError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["TokenCoreV1TokenRegistryNotConfigured"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackTokenCoreV1TokenRegistryNotConfiguredError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["TokenCoreV1TokenSymbolAlreadyExists"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackTokenCoreV1TokenSymbolAlreadyExistsError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["TokenCoreV1UnauthorizedCaller"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackTokenCoreV1UnauthorizedCallerError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["TokenCoreV1UnauthorizedFreezeManager"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackTokenCoreV1UnauthorizedFreezeManagerError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["UUPSUnauthorizedCallContext"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackUUPSUnauthorizedCallContextError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNTokenCoreV1.abi.Errors["UUPSUnsupportedProxiableUUID"].ID.Bytes()[:4]) {
		return pNTokenCoreV1.UnpackUUPSUnsupportedProxiableUUIDError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// PNTokenCoreV1AddressEmptyCode represents a AddressEmptyCode error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1AddressEmptyCode struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressEmptyCode(address target)
func PNTokenCoreV1AddressEmptyCodeErrorID() common.Hash {
	return common.HexToHash("0x9996b315c842ff135b8fc4a08ad5df1c344efbc03d2687aecc0678050d2aac89")
}

// UnpackAddressEmptyCodeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressEmptyCode(address target)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackAddressEmptyCodeError(raw []byte) (*PNTokenCoreV1AddressEmptyCode, error) {
	out := new(PNTokenCoreV1AddressEmptyCode)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "AddressEmptyCode", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1ERC1967InvalidImplementation represents a ERC1967InvalidImplementation error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1ERC1967InvalidImplementation struct {
	Implementation common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func PNTokenCoreV1ERC1967InvalidImplementationErrorID() common.Hash {
	return common.HexToHash("0x4c9c8ce3ceb3130f17f7cdba48d89b5b0129f266a8bac114e6e315a41879b617")
}

// UnpackERC1967InvalidImplementationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackERC1967InvalidImplementationError(raw []byte) (*PNTokenCoreV1ERC1967InvalidImplementation, error) {
	out := new(PNTokenCoreV1ERC1967InvalidImplementation)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "ERC1967InvalidImplementation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1ERC1967NonPayable represents a ERC1967NonPayable error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1ERC1967NonPayable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967NonPayable()
func PNTokenCoreV1ERC1967NonPayableErrorID() common.Hash {
	return common.HexToHash("0xb398979fa84f543c8e222f17890372c487baf85e062276c127fef521eea7224b")
}

// UnpackERC1967NonPayableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967NonPayable()
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackERC1967NonPayableError(raw []byte) (*PNTokenCoreV1ERC1967NonPayable, error) {
	out := new(PNTokenCoreV1ERC1967NonPayable)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "ERC1967NonPayable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1FailedCall represents a FailedCall error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1FailedCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedCall()
func PNTokenCoreV1FailedCallErrorID() common.Hash {
	return common.HexToHash("0xd6bda27508c0fb6d8a39b4b122878dab26f731a7d4e4abe711dd3731899052a4")
}

// UnpackFailedCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedCall()
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackFailedCallError(raw []byte) (*PNTokenCoreV1FailedCall, error) {
	out := new(PNTokenCoreV1FailedCall)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "FailedCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1InvalidInitialization represents a InvalidInitialization error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1InvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func PNTokenCoreV1InvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackInvalidInitializationError(raw []byte) (*PNTokenCoreV1InvalidInitialization, error) {
	out := new(PNTokenCoreV1InvalidInitialization)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1NotInitializing represents a NotInitializing error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1NotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func PNTokenCoreV1NotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackNotInitializingError(raw []byte) (*PNTokenCoreV1NotInitializing, error) {
	out := new(PNTokenCoreV1NotInitializing)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1RaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1RaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func PNTokenCoreV1RaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*PNTokenCoreV1RaylsAccessManagedContractPaused, error) {
	out := new(PNTokenCoreV1RaylsAccessManagedContractPaused)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1RaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1RaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func PNTokenCoreV1RaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*PNTokenCoreV1RaylsAccessManagedInvalidAuthority, error) {
	out := new(PNTokenCoreV1RaylsAccessManagedInvalidAuthority)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1RaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1RaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func PNTokenCoreV1RaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*PNTokenCoreV1RaylsAccessManagedMustSchedule, error) {
	out := new(PNTokenCoreV1RaylsAccessManagedMustSchedule)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1RaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1RaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func PNTokenCoreV1RaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*PNTokenCoreV1RaylsAccessManagedUnauthorized, error) {
	out := new(PNTokenCoreV1RaylsAccessManagedUnauthorized)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1TokenCoreV1EndpointNotConfigured represents a TokenCoreV1__EndpointNotConfigured error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenCoreV1EndpointNotConfigured struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenCoreV1__EndpointNotConfigured()
func PNTokenCoreV1TokenCoreV1EndpointNotConfiguredErrorID() common.Hash {
	return common.HexToHash("0xf117d70932d46523a51c5a42c7e85b830e1389b0ae4da9c462131a4fe82d9c37")
}

// UnpackTokenCoreV1EndpointNotConfiguredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenCoreV1__EndpointNotConfigured()
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenCoreV1EndpointNotConfiguredError(raw []byte) (*PNTokenCoreV1TokenCoreV1EndpointNotConfigured, error) {
	out := new(PNTokenCoreV1TokenCoreV1EndpointNotConfigured)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "TokenCoreV1EndpointNotConfigured", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1TokenCoreV1HubApprovalNotPending represents a TokenCoreV1__HubApprovalNotPending error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenCoreV1HubApprovalNotPending struct {
	TokenAddress  common.Address
	CurrentStatus uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenCoreV1__HubApprovalNotPending(address tokenAddress, uint8 currentStatus)
func PNTokenCoreV1TokenCoreV1HubApprovalNotPendingErrorID() common.Hash {
	return common.HexToHash("0x3e96fe34725ab61d01b36f0e53a49abcd466b4e8ce2e634d0f52377753e4b551")
}

// UnpackTokenCoreV1HubApprovalNotPendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenCoreV1__HubApprovalNotPending(address tokenAddress, uint8 currentStatus)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenCoreV1HubApprovalNotPendingError(raw []byte) (*PNTokenCoreV1TokenCoreV1HubApprovalNotPending, error) {
	out := new(PNTokenCoreV1TokenCoreV1HubApprovalNotPending)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "TokenCoreV1HubApprovalNotPending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1TokenCoreV1InvalidEndpointAddress represents a TokenCoreV1__InvalidEndpointAddress error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenCoreV1InvalidEndpointAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenCoreV1__InvalidEndpointAddress()
func PNTokenCoreV1TokenCoreV1InvalidEndpointAddressErrorID() common.Hash {
	return common.HexToHash("0x3807a3de4859afa92e6bb3b2a05f05ad15c2cb966f65587fd3f772b1d89f5f68")
}

// UnpackTokenCoreV1InvalidEndpointAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenCoreV1__InvalidEndpointAddress()
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenCoreV1InvalidEndpointAddressError(raw []byte) (*PNTokenCoreV1TokenCoreV1InvalidEndpointAddress, error) {
	out := new(PNTokenCoreV1TokenCoreV1InvalidEndpointAddress)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "TokenCoreV1InvalidEndpointAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1TokenCoreV1InvalidEnygmaPNEventsAddress represents a TokenCoreV1__InvalidEnygmaPNEventsAddress error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenCoreV1InvalidEnygmaPNEventsAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenCoreV1__InvalidEnygmaPNEventsAddress()
func PNTokenCoreV1TokenCoreV1InvalidEnygmaPNEventsAddressErrorID() common.Hash {
	return common.HexToHash("0x2bac7095e008575e204716710300aa2ef4cd2ff837617a0fd97a80b0fa980940")
}

// UnpackTokenCoreV1InvalidEnygmaPNEventsAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenCoreV1__InvalidEnygmaPNEventsAddress()
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenCoreV1InvalidEnygmaPNEventsAddressError(raw []byte) (*PNTokenCoreV1TokenCoreV1InvalidEnygmaPNEventsAddress, error) {
	out := new(PNTokenCoreV1TokenCoreV1InvalidEnygmaPNEventsAddress)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "TokenCoreV1InvalidEnygmaPNEventsAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1TokenCoreV1InvalidPrivacyNodeStatus represents a TokenCoreV1__InvalidPrivacyNodeStatus error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenCoreV1InvalidPrivacyNodeStatus struct {
	Status uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenCoreV1__InvalidPrivacyNodeStatus(uint8 status)
func PNTokenCoreV1TokenCoreV1InvalidPrivacyNodeStatusErrorID() common.Hash {
	return common.HexToHash("0x7b016f8a02c3b5a00ed0cd3466b95963106a8d4484c3cb2c9c2e831097c61430")
}

// UnpackTokenCoreV1InvalidPrivacyNodeStatusError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenCoreV1__InvalidPrivacyNodeStatus(uint8 status)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenCoreV1InvalidPrivacyNodeStatusError(raw []byte) (*PNTokenCoreV1TokenCoreV1InvalidPrivacyNodeStatus, error) {
	out := new(PNTokenCoreV1TokenCoreV1InvalidPrivacyNodeStatus)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "TokenCoreV1InvalidPrivacyNodeStatus", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1TokenCoreV1InvalidTokenAddress represents a TokenCoreV1__InvalidTokenAddress error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenCoreV1InvalidTokenAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenCoreV1__InvalidTokenAddress()
func PNTokenCoreV1TokenCoreV1InvalidTokenAddressErrorID() common.Hash {
	return common.HexToHash("0x6e8a97591d19ebcdb71cde0b86fbaee4ccd39ab0dc19e4b0caccece5a6d5c04a")
}

// UnpackTokenCoreV1InvalidTokenAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenCoreV1__InvalidTokenAddress()
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenCoreV1InvalidTokenAddressError(raw []byte) (*PNTokenCoreV1TokenCoreV1InvalidTokenAddress, error) {
	out := new(PNTokenCoreV1TokenCoreV1InvalidTokenAddress)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "TokenCoreV1InvalidTokenAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1TokenCoreV1InvalidTokenFreezeManagerAddress represents a TokenCoreV1__InvalidTokenFreezeManagerAddress error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenCoreV1InvalidTokenFreezeManagerAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenCoreV1__InvalidTokenFreezeManagerAddress()
func PNTokenCoreV1TokenCoreV1InvalidTokenFreezeManagerAddressErrorID() common.Hash {
	return common.HexToHash("0x92e680fb638a1d432747dbcaf536ebbfd6a768b5ca875b79111b8775e51d992d")
}

// UnpackTokenCoreV1InvalidTokenFreezeManagerAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenCoreV1__InvalidTokenFreezeManagerAddress()
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenCoreV1InvalidTokenFreezeManagerAddressError(raw []byte) (*PNTokenCoreV1TokenCoreV1InvalidTokenFreezeManagerAddress, error) {
	out := new(PNTokenCoreV1TokenCoreV1InvalidTokenFreezeManagerAddress)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "TokenCoreV1InvalidTokenFreezeManagerAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1TokenCoreV1InvalidTokenRegistryAddress represents a TokenCoreV1__InvalidTokenRegistryAddress error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenCoreV1InvalidTokenRegistryAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenCoreV1__InvalidTokenRegistryAddress()
func PNTokenCoreV1TokenCoreV1InvalidTokenRegistryAddressErrorID() common.Hash {
	return common.HexToHash("0xfce06f63b45e7104fae4d06cfdc285dd57bea95075c54d753bc3e8614e35796c")
}

// UnpackTokenCoreV1InvalidTokenRegistryAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenCoreV1__InvalidTokenRegistryAddress()
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenCoreV1InvalidTokenRegistryAddressError(raw []byte) (*PNTokenCoreV1TokenCoreV1InvalidTokenRegistryAddress, error) {
	out := new(PNTokenCoreV1TokenCoreV1InvalidTokenRegistryAddress)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "TokenCoreV1InvalidTokenRegistryAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1TokenCoreV1PrivacyNodeAuthorizationRequired represents a TokenCoreV1__PrivacyNodeAuthorizationRequired error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenCoreV1PrivacyNodeAuthorizationRequired struct {
	TokenAddress common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenCoreV1__PrivacyNodeAuthorizationRequired(address tokenAddress)
func PNTokenCoreV1TokenCoreV1PrivacyNodeAuthorizationRequiredErrorID() common.Hash {
	return common.HexToHash("0x3f35fa868f58890cba4168967530b14b475141fb70565bb132aada2eb192f8bc")
}

// UnpackTokenCoreV1PrivacyNodeAuthorizationRequiredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenCoreV1__PrivacyNodeAuthorizationRequired(address tokenAddress)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenCoreV1PrivacyNodeAuthorizationRequiredError(raw []byte) (*PNTokenCoreV1TokenCoreV1PrivacyNodeAuthorizationRequired, error) {
	out := new(PNTokenCoreV1TokenCoreV1PrivacyNodeAuthorizationRequired)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "TokenCoreV1PrivacyNodeAuthorizationRequired", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1TokenCoreV1ResourceIdAlreadyAssigned represents a TokenCoreV1__ResourceIdAlreadyAssigned error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenCoreV1ResourceIdAlreadyAssigned struct {
	ResourceId [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenCoreV1__ResourceIdAlreadyAssigned(bytes32 resourceId)
func PNTokenCoreV1TokenCoreV1ResourceIdAlreadyAssignedErrorID() common.Hash {
	return common.HexToHash("0xc3bf60f9a03596b7836a92099fb183b3178012c8752850af4f7c4995fdbc8331")
}

// UnpackTokenCoreV1ResourceIdAlreadyAssignedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenCoreV1__ResourceIdAlreadyAssigned(bytes32 resourceId)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenCoreV1ResourceIdAlreadyAssignedError(raw []byte) (*PNTokenCoreV1TokenCoreV1ResourceIdAlreadyAssigned, error) {
	out := new(PNTokenCoreV1TokenCoreV1ResourceIdAlreadyAssigned)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "TokenCoreV1ResourceIdAlreadyAssigned", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1TokenCoreV1ResourceIdAlreadySet represents a TokenCoreV1__ResourceIdAlreadySet error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenCoreV1ResourceIdAlreadySet struct {
	TokenAddress common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenCoreV1__ResourceIdAlreadySet(address tokenAddress)
func PNTokenCoreV1TokenCoreV1ResourceIdAlreadySetErrorID() common.Hash {
	return common.HexToHash("0xf117323b7559a5f4bdf0f0a0270628e2fb053103547ec1a0a0b88a4c32f7b019")
}

// UnpackTokenCoreV1ResourceIdAlreadySetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenCoreV1__ResourceIdAlreadySet(address tokenAddress)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenCoreV1ResourceIdAlreadySetError(raw []byte) (*PNTokenCoreV1TokenCoreV1ResourceIdAlreadySet, error) {
	out := new(PNTokenCoreV1TokenCoreV1ResourceIdAlreadySet)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "TokenCoreV1ResourceIdAlreadySet", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1TokenCoreV1StatusAlreadySet represents a TokenCoreV1__StatusAlreadySet error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenCoreV1StatusAlreadySet struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenCoreV1__StatusAlreadySet()
func PNTokenCoreV1TokenCoreV1StatusAlreadySetErrorID() common.Hash {
	return common.HexToHash("0xfa227edaca97825570dcda4ab681fa445912bc818f9341c2d0bd781dd69c9857")
}

// UnpackTokenCoreV1StatusAlreadySetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenCoreV1__StatusAlreadySet()
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenCoreV1StatusAlreadySetError(raw []byte) (*PNTokenCoreV1TokenCoreV1StatusAlreadySet, error) {
	out := new(PNTokenCoreV1TokenCoreV1StatusAlreadySet)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "TokenCoreV1StatusAlreadySet", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1TokenCoreV1TokenAlreadyExists represents a TokenCoreV1__TokenAlreadyExists error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenCoreV1TokenAlreadyExists struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenCoreV1__TokenAlreadyExists()
func PNTokenCoreV1TokenCoreV1TokenAlreadyExistsErrorID() common.Hash {
	return common.HexToHash("0xc19236b2fc948de9cfe1351f0934b073e27457bdfe21e1f525b89077d089cf41")
}

// UnpackTokenCoreV1TokenAlreadyExistsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenCoreV1__TokenAlreadyExists()
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenCoreV1TokenAlreadyExistsError(raw []byte) (*PNTokenCoreV1TokenCoreV1TokenAlreadyExists, error) {
	out := new(PNTokenCoreV1TokenCoreV1TokenAlreadyExists)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "TokenCoreV1TokenAlreadyExists", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1TokenCoreV1TokenDoesNotExist represents a TokenCoreV1__TokenDoesNotExist error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenCoreV1TokenDoesNotExist struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenCoreV1__TokenDoesNotExist()
func PNTokenCoreV1TokenCoreV1TokenDoesNotExistErrorID() common.Hash {
	return common.HexToHash("0xb2be0e7fd66a9214d291149ed062d8a8e19bdafa5efa9b7611afc35478702a7c")
}

// UnpackTokenCoreV1TokenDoesNotExistError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenCoreV1__TokenDoesNotExist()
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenCoreV1TokenDoesNotExistError(raw []byte) (*PNTokenCoreV1TokenCoreV1TokenDoesNotExist, error) {
	out := new(PNTokenCoreV1TokenCoreV1TokenDoesNotExist)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "TokenCoreV1TokenDoesNotExist", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1TokenCoreV1TokenRegistryNotConfigured represents a TokenCoreV1__TokenRegistryNotConfigured error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenCoreV1TokenRegistryNotConfigured struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenCoreV1__TokenRegistryNotConfigured()
func PNTokenCoreV1TokenCoreV1TokenRegistryNotConfiguredErrorID() common.Hash {
	return common.HexToHash("0x525229e6c70d169cd7fd2b6bc6c94ce34d4241af11be7000c06887310ede58e2")
}

// UnpackTokenCoreV1TokenRegistryNotConfiguredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenCoreV1__TokenRegistryNotConfigured()
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenCoreV1TokenRegistryNotConfiguredError(raw []byte) (*PNTokenCoreV1TokenCoreV1TokenRegistryNotConfigured, error) {
	out := new(PNTokenCoreV1TokenCoreV1TokenRegistryNotConfigured)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "TokenCoreV1TokenRegistryNotConfigured", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1TokenCoreV1TokenSymbolAlreadyExists represents a TokenCoreV1__TokenSymbolAlreadyExists error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenCoreV1TokenSymbolAlreadyExists struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenCoreV1__TokenSymbolAlreadyExists()
func PNTokenCoreV1TokenCoreV1TokenSymbolAlreadyExistsErrorID() common.Hash {
	return common.HexToHash("0xfae70dbd6cac0c6d322de00b8e13f044fe24d790efb1ac43bd0bdc3a32209d62")
}

// UnpackTokenCoreV1TokenSymbolAlreadyExistsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenCoreV1__TokenSymbolAlreadyExists()
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenCoreV1TokenSymbolAlreadyExistsError(raw []byte) (*PNTokenCoreV1TokenCoreV1TokenSymbolAlreadyExists, error) {
	out := new(PNTokenCoreV1TokenCoreV1TokenSymbolAlreadyExists)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "TokenCoreV1TokenSymbolAlreadyExists", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1TokenCoreV1UnauthorizedCaller represents a TokenCoreV1__UnauthorizedCaller error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenCoreV1UnauthorizedCaller struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenCoreV1__UnauthorizedCaller(address caller)
func PNTokenCoreV1TokenCoreV1UnauthorizedCallerErrorID() common.Hash {
	return common.HexToHash("0xe5155ffc8ff18d9b9064a3a00a769f93101218c4e7681850e00ca86c640281d7")
}

// UnpackTokenCoreV1UnauthorizedCallerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenCoreV1__UnauthorizedCaller(address caller)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenCoreV1UnauthorizedCallerError(raw []byte) (*PNTokenCoreV1TokenCoreV1UnauthorizedCaller, error) {
	out := new(PNTokenCoreV1TokenCoreV1UnauthorizedCaller)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "TokenCoreV1UnauthorizedCaller", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1TokenCoreV1UnauthorizedFreezeManager represents a TokenCoreV1__UnauthorizedFreezeManager error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1TokenCoreV1UnauthorizedFreezeManager struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenCoreV1__UnauthorizedFreezeManager(address caller)
func PNTokenCoreV1TokenCoreV1UnauthorizedFreezeManagerErrorID() common.Hash {
	return common.HexToHash("0x3940d26625c1a66a8d70a4832ab82970b6cdea7260d418045cae0654df21396b")
}

// UnpackTokenCoreV1UnauthorizedFreezeManagerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenCoreV1__UnauthorizedFreezeManager(address caller)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackTokenCoreV1UnauthorizedFreezeManagerError(raw []byte) (*PNTokenCoreV1TokenCoreV1UnauthorizedFreezeManager, error) {
	out := new(PNTokenCoreV1TokenCoreV1UnauthorizedFreezeManager)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "TokenCoreV1UnauthorizedFreezeManager", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1UUPSUnauthorizedCallContext represents a UUPSUnauthorizedCallContext error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1UUPSUnauthorizedCallContext struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnauthorizedCallContext()
func PNTokenCoreV1UUPSUnauthorizedCallContextErrorID() common.Hash {
	return common.HexToHash("0xe07c8dba242a06571ac65fe4bbe20522c9fb111cb33599b799ff8039c1ed18f4")
}

// UnpackUUPSUnauthorizedCallContextError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnauthorizedCallContext()
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackUUPSUnauthorizedCallContextError(raw []byte) (*PNTokenCoreV1UUPSUnauthorizedCallContext, error) {
	out := new(PNTokenCoreV1UUPSUnauthorizedCallContext)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "UUPSUnauthorizedCallContext", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNTokenCoreV1UUPSUnsupportedProxiableUUID represents a UUPSUnsupportedProxiableUUID error raised by the PNTokenCoreV1 contract.
type PNTokenCoreV1UUPSUnsupportedProxiableUUID struct {
	Slot [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func PNTokenCoreV1UUPSUnsupportedProxiableUUIDErrorID() common.Hash {
	return common.HexToHash("0xaa1d49a4c084bfa9aeeee2a0be65267a7f19ba7e1476b114dac513d2c14cb563")
}

// UnpackUUPSUnsupportedProxiableUUIDError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func (pNTokenCoreV1 *PNTokenCoreV1) UnpackUUPSUnsupportedProxiableUUIDError(raw []byte) (*PNTokenCoreV1UUPSUnsupportedProxiableUUID, error) {
	out := new(PNTokenCoreV1UUPSUnsupportedProxiableUUID)
	if err := pNTokenCoreV1.abi.UnpackIntoInterface(out, "UUPSUnsupportedProxiableUUID", raw); err != nil {
		return nil, err
	}
	return out, nil
}
