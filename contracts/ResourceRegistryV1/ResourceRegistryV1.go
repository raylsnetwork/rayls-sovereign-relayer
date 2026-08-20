// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ResourceRegistryV1

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

// ResourceRegistryV1Resource is an auto generated low-level Go binding around an user-defined struct.
type ResourceRegistryV1Resource struct {
	ResourceId        [32]byte
	Standard          uint8
	Bytecode          []byte
	InitializerParams []byte
}

// ResourceRegistryV1MetaData contains all meta data concerning the ResourceRegistryV1 contract.
var ResourceRegistryV1MetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"contractVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getResourceById\",\"inputs\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structResourceRegistryV1.Resource\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"standard\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.ErcStandard\"},{\"name\":\"bytecode\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"initializerParams\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"authority_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerResource\",\"inputs\":[{\"name\":\"standard\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.ErcStandard\"},{\"name\":\"bytecode\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"initializerParams\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setTokenRegistry\",\"inputs\":[{\"name\":\"tokenRegistryAt\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ResourceRegistryV1__TokenRegistryNotSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ResourceRegistryV1__UnauthorizedCaller\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	ID:  "ResourceRegistryV1",
	Bin: "0x60a06040523060805234801561001457600080fd5b5060805161138961003e60003960008181610886015281816108af01526109e801526113896000f3fe6080604052600436106100765760003560e01c806335a5af921461007b578063417004881461009d5780634f1ef286146100d057806352d1902d146100e3578063a0a8e460146100f8578063ad3cb1cc1461010c578063bc8a6f201461014a578063bf7e214f14610177578063c4d66de814610199575b600080fd5b34801561008757600080fd5b5061009b610096366004610e77565b6101b9565b005b3480156100a957600080fd5b506100bd6100b8366004610f34565b6101fb565b6040519081526020015b60405180910390f35b61009b6100de366004610fad565b6103c8565b3480156100ef57600080fd5b506100bd6103e7565b34801561010457600080fd5b5060016100bd565b34801561011857600080fd5b5061013d604051806040016040528060058152602001640352e302e360dc1b81525081565b6040516100c7919061104a565b34801561015657600080fd5b5061016a61016536600461105d565b610404565b6040516100c7919061108c565b34801561018357600080fd5b5061018c6105a0565b6040516100c791906110f9565b3480156101a557600080fd5b5061009b6101b4366004610e77565b6105b9565b6101cf336000356001600160e01b0319166106bb565b60006101d9610806565b80546001600160a01b0319166001600160a01b03939093169290921790915550565b600080610206610806565b80549091506001600160a01b03166102315760405163046d162160e11b815260040160405180910390fd5b80546001600160a01b031633146102665733604051636543a02360e01b815260040161025d91906110f9565b60405180910390fd5b6000610270610806565b9050600061027c61082a565b6000818152600284016020526040902054909150156102f55760405162461bcd60e51b815260206004820152602f60248201527f5265736f75726365526567697374727956313a205265736f7572636520616c7260448201526e1958591e481c9959da5cdd195c9959608a1b606482015260840161025d565b81600101604051806080016040528083815260200189600c81111561031c5761031c611076565b815260208082018a905260409091018890528254600181810185556000948552938290208351600490920201908155908201518184018054939492939192909160ff19169083600c81111561037357610373611076565b02179055506040820151600282019061038c908261118f565b50606082015160038201906103a1908261118f565b50505060018201546000828152600290930160205260409092209190915595945050505050565b6103d061087b565b6103d98261090b565b6103e38282610924565b5050565b60006103f16109dd565b5060008051602061133483398151915290565b604080516080810182526000808252602082015260609181018290528181019190915261043082610a26565b6040805160808101909152815481526001820154909190602083019060ff16600c81111561046057610460611076565b600c81111561047157610471611076565b81526020016002820180546104859061110d565b80601f01602080910402602001604051908101604052809291908181526020018280546104b19061110d565b80156104fe5780601f106104d3576101008083540402835291602001916104fe565b820191906000526020600020905b8154815290600101906020018083116104e157829003601f168201915b505050505081526020016003820180546105179061110d565b80601f01602080910402602001604051908101604052809291908181526020018280546105439061110d565b80156105905780601f1061056557610100808354040283529160200191610590565b820191906000526020600020905b81548152906001019060200180831161057357829003601f168201915b5050505050815250509050919050565b60006105aa610adb565b546001600160a01b0316919050565b60006105c3610b3d565b805490915060ff600160401b82041615906001600160401b03166000811580156105ea5750825b90506000826001600160401b031660011480156106065750303b155b905081158015610614575080155b156106325760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561065c57845460ff60401b1916600160401b1785555b610664610b68565b61066d86610b70565b83156106b357845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b505050505050565b60006106c5610adb565b80549091506001600160a01b0316806106f4576000604051638944034760e01b815260040161025d91906110f9565b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa158015610758573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061077c919061125e565b925092509250826107fd5780156107a65760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff8216156107e25760405163a426878960e01b81526001600160a01b038816600482015263ffffffff8316602482015260440161025d565b86604051632ecd3d0360e21b815260040161025d91906110f9565b50505050505050565b7f5f092f6026c66587fbb8d90f812020d9443212d551a1ca8fbb1a34b93d86cc0090565b600080610835610806565b600381018054919250600190600061084d83856112c2565b9091555050604080516020810183905201604051602081830303815290604052805190602001209250505090565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614806108eb57507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166108df610bb1565b6001600160a01b031614155b156109095760405163703e46dd60e11b815260040160405180910390fd5b565b610921336000356001600160e01b0319166106bb565b50565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa92505050801561097e575060408051601f3d908101601f1916820190925261097b918101906112d5565b60015b61099d5781604051634c9c8ce360e01b815260040161025d91906110f9565b60008051602061133483398151915281146109ce57604051632a87526960e21b81526004810182905260240161025d565b6109d88383610bc7565b505050565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146109095760405163703e46dd60e11b815260040160405180910390fd5b600080610a31610806565b6000848152600282016020526040812054919250819003610aa35760405162461bcd60e51b815260206004820152602660248201527f5265736f75726365526567697374727956313a205265736f75726365206e6f7460448201526508199bdd5b9960d21b606482015260840161025d565b81600101600182610ab491906112ee565b81548110610ac457610ac4611301565b906000526020600020906004020192505050919050565b60008060ff19610b0c60017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f356112ee565b604051602001610b1e91815260200190565b60408051601f1981840301815291905280516020909101201692915050565b6000807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005b92915050565b610909610c1d565b6000610b7a610adb565b80549091506001600160a01b031615610ba85781604051638944034760e01b815260040161025d91906110f9565b6103e382610c42565b60006000805160206113348339815191526105aa565b610bd082610cd2565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a2805115610c15576109d88282610d2e565b6103e3610da4565b610c25610dc3565b61090957604051631afcd79f60e31b815260040160405180910390fd5b6001600160a01b038116610c6b5780604051638944034760e01b815260040161025d91906110f9565b6000610c75610adb565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b806001600160a01b03163b600003610cff5780604051634c9c8ce360e01b815260040161025d91906110f9565b60008051602061133483398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080846001600160a01b031684604051610d4b9190611317565b600060405180830381855af49150503d8060008114610d86576040519150601f19603f3d011682016040523d82523d6000602084013e610d8b565b606091505b5091509150610d9b858383610ddd565b95945050505050565b34156109095760405163b398979f60e01b815260040160405180910390fd5b6000610dcd610b3d565b54600160401b900460ff16919050565b606082610df257610ded82610e33565b610e2c565b8151158015610e0957506001600160a01b0384163b155b15610e295783604051639996b31560e01b815260040161025d91906110f9565b50805b9392505050565b805115610e4257805160208201fd5b60405163d6bda27560e01b815260040160405180910390fd5b80356001600160a01b0381168114610e7257600080fd5b919050565b600060208284031215610e8957600080fd5b610e2c82610e5b565b634e487b7160e01b600052604160045260246000fd5b600082601f830112610eb957600080fd5b81356001600160401b0380821115610ed357610ed3610e92565b604051601f8301601f19908116603f01168101908282118183101715610efb57610efb610e92565b81604052838152866020858801011115610f1457600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080600060608486031215610f4957600080fd5b8335600d8110610f5857600080fd5b925060208401356001600160401b0380821115610f7457600080fd5b610f8087838801610ea8565b93506040860135915080821115610f9657600080fd5b50610fa386828701610ea8565b9150509250925092565b60008060408385031215610fc057600080fd5b610fc983610e5b565b915060208301356001600160401b03811115610fe457600080fd5b610ff085828601610ea8565b9150509250929050565b60005b83811015611015578181015183820152602001610ffd565b50506000910152565b60008151808452611036816020860160208601610ffa565b601f01601f19169290920160200192915050565b602081526000610e2c602083018461101e565b60006020828403121561106f57600080fd5b5035919050565b634e487b7160e01b600052602160045260246000fd5b602081528151602082015260006020830151600d81106110bc57634e487b7160e01b600052602160045260246000fd5b806040840152506040830151608060608401526110dc60a084018261101e565b90506060840151601f19848303016080850152610d9b828261101e565b6001600160a01b0391909116815260200190565b600181811c9082168061112157607f821691505b60208210810361114157634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156109d8576000816000526020600020601f850160051c810160208610156111705750805b601f850160051c820191505b818110156106b35782815560010161117c565b81516001600160401b038111156111a8576111a8610e92565b6111bc816111b6845461110d565b84611147565b602080601f8311600181146111f157600084156111d95750858301515b600019600386901b1c1916600185901b1785556106b3565b600085815260208120601f198616915b8281101561122057888601518255948401946001909101908401611201565b508582101561123e5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b80518015158114610e7257600080fd5b60008060006060848603121561127357600080fd5b61127c8461124e565b9250602084015163ffffffff8116811461129557600080fd5b91506112a36040850161124e565b90509250925092565b634e487b7160e01b600052601160045260246000fd5b80820180821115610b6257610b626112ac565b6000602082840312156112e757600080fd5b5051919050565b81810381811115610b6257610b626112ac565b634e487b7160e01b600052603260045260246000fd5b60008251611329818460208701610ffa565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbca2646970667358221220ba01128f9a33a63f783121689869599a2fd6f00eeebd213b9c759a083b3d567a64736f6c63430008180033",
}

// ResourceRegistryV1 is an auto generated Go binding around an Ethereum contract.
type ResourceRegistryV1 struct {
	abi abi.ABI
}

// NewResourceRegistryV1 creates a new instance of ResourceRegistryV1.
func NewResourceRegistryV1() *ResourceRegistryV1 {
	parsed, err := ResourceRegistryV1MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ResourceRegistryV1{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ResourceRegistryV1) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackUPGRADEINTERFACEVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (resourceRegistryV1 *ResourceRegistryV1) PackUPGRADEINTERFACEVERSION() []byte {
	enc, err := resourceRegistryV1.abi.Pack("UPGRADE_INTERFACE_VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUPGRADEINTERFACEVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (resourceRegistryV1 *ResourceRegistryV1) UnpackUPGRADEINTERFACEVERSION(data []byte) (string, error) {
	out, err := resourceRegistryV1.abi.Unpack("UPGRADE_INTERFACE_VERSION", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (resourceRegistryV1 *ResourceRegistryV1) PackAuthority() []byte {
	enc, err := resourceRegistryV1.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (resourceRegistryV1 *ResourceRegistryV1) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := resourceRegistryV1.abi.Unpack("authority", data)
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
func (resourceRegistryV1 *ResourceRegistryV1) PackContractVersion() []byte {
	enc, err := resourceRegistryV1.abi.Pack("contractVersion")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackContractVersion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa0a8e460.
//
// Solidity: function contractVersion() pure returns(uint256)
func (resourceRegistryV1 *ResourceRegistryV1) UnpackContractVersion(data []byte) (*big.Int, error) {
	out, err := resourceRegistryV1.abi.Unpack("contractVersion", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetResourceById is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbc8a6f20.
//
// Solidity: function getResourceById(bytes32 resourceId) view returns((bytes32,uint8,bytes,bytes))
func (resourceRegistryV1 *ResourceRegistryV1) PackGetResourceById(resourceId [32]byte) []byte {
	enc, err := resourceRegistryV1.abi.Pack("getResourceById", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetResourceById is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbc8a6f20.
//
// Solidity: function getResourceById(bytes32 resourceId) view returns((bytes32,uint8,bytes,bytes))
func (resourceRegistryV1 *ResourceRegistryV1) UnpackGetResourceById(data []byte) (ResourceRegistryV1Resource, error) {
	out, err := resourceRegistryV1.abi.Unpack("getResourceById", data)
	if err != nil {
		return *new(ResourceRegistryV1Resource), err
	}
	out0 := *abi.ConvertType(out[0], new(ResourceRegistryV1Resource)).(*ResourceRegistryV1Resource)
	return out0, err
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc4d66de8.
//
// Solidity: function initialize(address authority_) returns()
func (resourceRegistryV1 *ResourceRegistryV1) PackInitialize(authority common.Address) []byte {
	enc, err := resourceRegistryV1.abi.Pack("initialize", authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackProxiableUUID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (resourceRegistryV1 *ResourceRegistryV1) PackProxiableUUID() []byte {
	enc, err := resourceRegistryV1.abi.Pack("proxiableUUID")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackProxiableUUID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (resourceRegistryV1 *ResourceRegistryV1) UnpackProxiableUUID(data []byte) ([32]byte, error) {
	out, err := resourceRegistryV1.abi.Unpack("proxiableUUID", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackRegisterResource is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x41700488.
//
// Solidity: function registerResource(uint8 standard, bytes bytecode, bytes initializerParams) returns(bytes32)
func (resourceRegistryV1 *ResourceRegistryV1) PackRegisterResource(standard uint8, bytecode []byte, initializerParams []byte) []byte {
	enc, err := resourceRegistryV1.abi.Pack("registerResource", standard, bytecode, initializerParams)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackRegisterResource is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x41700488.
//
// Solidity: function registerResource(uint8 standard, bytes bytecode, bytes initializerParams) returns(bytes32)
func (resourceRegistryV1 *ResourceRegistryV1) UnpackRegisterResource(data []byte) ([32]byte, error) {
	out, err := resourceRegistryV1.abi.Unpack("registerResource", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackSetTokenRegistry is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x35a5af92.
//
// Solidity: function setTokenRegistry(address tokenRegistryAt) returns()
func (resourceRegistryV1 *ResourceRegistryV1) PackSetTokenRegistry(tokenRegistryAt common.Address) []byte {
	enc, err := resourceRegistryV1.abi.Pack("setTokenRegistry", tokenRegistryAt)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpgradeToAndCall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (resourceRegistryV1 *ResourceRegistryV1) PackUpgradeToAndCall(newImplementation common.Address, data []byte) []byte {
	enc, err := resourceRegistryV1.abi.Pack("upgradeToAndCall", newImplementation, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// ResourceRegistryV1AuthorityUpdated represents a AuthorityUpdated event raised by the ResourceRegistryV1 contract.
type ResourceRegistryV1AuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const ResourceRegistryV1AuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (ResourceRegistryV1AuthorityUpdated) ContractEventName() string {
	return ResourceRegistryV1AuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (resourceRegistryV1 *ResourceRegistryV1) UnpackAuthorityUpdatedEvent(log *types.Log) (*ResourceRegistryV1AuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != resourceRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ResourceRegistryV1AuthorityUpdated)
	if len(log.Data) > 0 {
		if err := resourceRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range resourceRegistryV1.abi.Events[event].Inputs {
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

// ResourceRegistryV1Initialized represents a Initialized event raised by the ResourceRegistryV1 contract.
type ResourceRegistryV1Initialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const ResourceRegistryV1InitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (ResourceRegistryV1Initialized) ContractEventName() string {
	return ResourceRegistryV1InitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (resourceRegistryV1 *ResourceRegistryV1) UnpackInitializedEvent(log *types.Log) (*ResourceRegistryV1Initialized, error) {
	event := "Initialized"
	if log.Topics[0] != resourceRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ResourceRegistryV1Initialized)
	if len(log.Data) > 0 {
		if err := resourceRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range resourceRegistryV1.abi.Events[event].Inputs {
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

// ResourceRegistryV1Upgraded represents a Upgraded event raised by the ResourceRegistryV1 contract.
type ResourceRegistryV1Upgraded struct {
	Implementation common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const ResourceRegistryV1UpgradedEventName = "Upgraded"

// ContractEventName returns the user-defined event name.
func (ResourceRegistryV1Upgraded) ContractEventName() string {
	return ResourceRegistryV1UpgradedEventName
}

// UnpackUpgradedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Upgraded(address indexed implementation)
func (resourceRegistryV1 *ResourceRegistryV1) UnpackUpgradedEvent(log *types.Log) (*ResourceRegistryV1Upgraded, error) {
	event := "Upgraded"
	if log.Topics[0] != resourceRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ResourceRegistryV1Upgraded)
	if len(log.Data) > 0 {
		if err := resourceRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range resourceRegistryV1.abi.Events[event].Inputs {
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
func (resourceRegistryV1 *ResourceRegistryV1) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], resourceRegistryV1.abi.Errors["AddressEmptyCode"].ID.Bytes()[:4]) {
		return resourceRegistryV1.UnpackAddressEmptyCodeError(raw[4:])
	}
	if bytes.Equal(raw[:4], resourceRegistryV1.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return resourceRegistryV1.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], resourceRegistryV1.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return resourceRegistryV1.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], resourceRegistryV1.abi.Errors["FailedCall"].ID.Bytes()[:4]) {
		return resourceRegistryV1.UnpackFailedCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], resourceRegistryV1.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return resourceRegistryV1.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], resourceRegistryV1.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return resourceRegistryV1.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], resourceRegistryV1.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return resourceRegistryV1.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], resourceRegistryV1.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return resourceRegistryV1.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], resourceRegistryV1.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return resourceRegistryV1.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], resourceRegistryV1.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return resourceRegistryV1.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], resourceRegistryV1.abi.Errors["ResourceRegistryV1TokenRegistryNotSet"].ID.Bytes()[:4]) {
		return resourceRegistryV1.UnpackResourceRegistryV1TokenRegistryNotSetError(raw[4:])
	}
	if bytes.Equal(raw[:4], resourceRegistryV1.abi.Errors["ResourceRegistryV1UnauthorizedCaller"].ID.Bytes()[:4]) {
		return resourceRegistryV1.UnpackResourceRegistryV1UnauthorizedCallerError(raw[4:])
	}
	if bytes.Equal(raw[:4], resourceRegistryV1.abi.Errors["UUPSUnauthorizedCallContext"].ID.Bytes()[:4]) {
		return resourceRegistryV1.UnpackUUPSUnauthorizedCallContextError(raw[4:])
	}
	if bytes.Equal(raw[:4], resourceRegistryV1.abi.Errors["UUPSUnsupportedProxiableUUID"].ID.Bytes()[:4]) {
		return resourceRegistryV1.UnpackUUPSUnsupportedProxiableUUIDError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ResourceRegistryV1AddressEmptyCode represents a AddressEmptyCode error raised by the ResourceRegistryV1 contract.
type ResourceRegistryV1AddressEmptyCode struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressEmptyCode(address target)
func ResourceRegistryV1AddressEmptyCodeErrorID() common.Hash {
	return common.HexToHash("0x9996b315c842ff135b8fc4a08ad5df1c344efbc03d2687aecc0678050d2aac89")
}

// UnpackAddressEmptyCodeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressEmptyCode(address target)
func (resourceRegistryV1 *ResourceRegistryV1) UnpackAddressEmptyCodeError(raw []byte) (*ResourceRegistryV1AddressEmptyCode, error) {
	out := new(ResourceRegistryV1AddressEmptyCode)
	if err := resourceRegistryV1.abi.UnpackIntoInterface(out, "AddressEmptyCode", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ResourceRegistryV1ERC1967InvalidImplementation represents a ERC1967InvalidImplementation error raised by the ResourceRegistryV1 contract.
type ResourceRegistryV1ERC1967InvalidImplementation struct {
	Implementation common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func ResourceRegistryV1ERC1967InvalidImplementationErrorID() common.Hash {
	return common.HexToHash("0x4c9c8ce3ceb3130f17f7cdba48d89b5b0129f266a8bac114e6e315a41879b617")
}

// UnpackERC1967InvalidImplementationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func (resourceRegistryV1 *ResourceRegistryV1) UnpackERC1967InvalidImplementationError(raw []byte) (*ResourceRegistryV1ERC1967InvalidImplementation, error) {
	out := new(ResourceRegistryV1ERC1967InvalidImplementation)
	if err := resourceRegistryV1.abi.UnpackIntoInterface(out, "ERC1967InvalidImplementation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ResourceRegistryV1ERC1967NonPayable represents a ERC1967NonPayable error raised by the ResourceRegistryV1 contract.
type ResourceRegistryV1ERC1967NonPayable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967NonPayable()
func ResourceRegistryV1ERC1967NonPayableErrorID() common.Hash {
	return common.HexToHash("0xb398979fa84f543c8e222f17890372c487baf85e062276c127fef521eea7224b")
}

// UnpackERC1967NonPayableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967NonPayable()
func (resourceRegistryV1 *ResourceRegistryV1) UnpackERC1967NonPayableError(raw []byte) (*ResourceRegistryV1ERC1967NonPayable, error) {
	out := new(ResourceRegistryV1ERC1967NonPayable)
	if err := resourceRegistryV1.abi.UnpackIntoInterface(out, "ERC1967NonPayable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ResourceRegistryV1FailedCall represents a FailedCall error raised by the ResourceRegistryV1 contract.
type ResourceRegistryV1FailedCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedCall()
func ResourceRegistryV1FailedCallErrorID() common.Hash {
	return common.HexToHash("0xd6bda27508c0fb6d8a39b4b122878dab26f731a7d4e4abe711dd3731899052a4")
}

// UnpackFailedCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedCall()
func (resourceRegistryV1 *ResourceRegistryV1) UnpackFailedCallError(raw []byte) (*ResourceRegistryV1FailedCall, error) {
	out := new(ResourceRegistryV1FailedCall)
	if err := resourceRegistryV1.abi.UnpackIntoInterface(out, "FailedCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ResourceRegistryV1InvalidInitialization represents a InvalidInitialization error raised by the ResourceRegistryV1 contract.
type ResourceRegistryV1InvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func ResourceRegistryV1InvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (resourceRegistryV1 *ResourceRegistryV1) UnpackInvalidInitializationError(raw []byte) (*ResourceRegistryV1InvalidInitialization, error) {
	out := new(ResourceRegistryV1InvalidInitialization)
	if err := resourceRegistryV1.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ResourceRegistryV1NotInitializing represents a NotInitializing error raised by the ResourceRegistryV1 contract.
type ResourceRegistryV1NotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func ResourceRegistryV1NotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (resourceRegistryV1 *ResourceRegistryV1) UnpackNotInitializingError(raw []byte) (*ResourceRegistryV1NotInitializing, error) {
	out := new(ResourceRegistryV1NotInitializing)
	if err := resourceRegistryV1.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ResourceRegistryV1RaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the ResourceRegistryV1 contract.
type ResourceRegistryV1RaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func ResourceRegistryV1RaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (resourceRegistryV1 *ResourceRegistryV1) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*ResourceRegistryV1RaylsAccessManagedContractPaused, error) {
	out := new(ResourceRegistryV1RaylsAccessManagedContractPaused)
	if err := resourceRegistryV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ResourceRegistryV1RaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the ResourceRegistryV1 contract.
type ResourceRegistryV1RaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func ResourceRegistryV1RaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (resourceRegistryV1 *ResourceRegistryV1) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*ResourceRegistryV1RaylsAccessManagedInvalidAuthority, error) {
	out := new(ResourceRegistryV1RaylsAccessManagedInvalidAuthority)
	if err := resourceRegistryV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ResourceRegistryV1RaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the ResourceRegistryV1 contract.
type ResourceRegistryV1RaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func ResourceRegistryV1RaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (resourceRegistryV1 *ResourceRegistryV1) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*ResourceRegistryV1RaylsAccessManagedMustSchedule, error) {
	out := new(ResourceRegistryV1RaylsAccessManagedMustSchedule)
	if err := resourceRegistryV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ResourceRegistryV1RaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the ResourceRegistryV1 contract.
type ResourceRegistryV1RaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func ResourceRegistryV1RaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (resourceRegistryV1 *ResourceRegistryV1) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*ResourceRegistryV1RaylsAccessManagedUnauthorized, error) {
	out := new(ResourceRegistryV1RaylsAccessManagedUnauthorized)
	if err := resourceRegistryV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ResourceRegistryV1ResourceRegistryV1TokenRegistryNotSet represents a ResourceRegistryV1__TokenRegistryNotSet error raised by the ResourceRegistryV1 contract.
type ResourceRegistryV1ResourceRegistryV1TokenRegistryNotSet struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ResourceRegistryV1__TokenRegistryNotSet()
func ResourceRegistryV1ResourceRegistryV1TokenRegistryNotSetErrorID() common.Hash {
	return common.HexToHash("0x08da2c427b3b2adf19489632a03ff3e6b8ae3707b1e24abd71c8c28dea061b2d")
}

// UnpackResourceRegistryV1TokenRegistryNotSetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ResourceRegistryV1__TokenRegistryNotSet()
func (resourceRegistryV1 *ResourceRegistryV1) UnpackResourceRegistryV1TokenRegistryNotSetError(raw []byte) (*ResourceRegistryV1ResourceRegistryV1TokenRegistryNotSet, error) {
	out := new(ResourceRegistryV1ResourceRegistryV1TokenRegistryNotSet)
	if err := resourceRegistryV1.abi.UnpackIntoInterface(out, "ResourceRegistryV1TokenRegistryNotSet", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ResourceRegistryV1ResourceRegistryV1UnauthorizedCaller represents a ResourceRegistryV1__UnauthorizedCaller error raised by the ResourceRegistryV1 contract.
type ResourceRegistryV1ResourceRegistryV1UnauthorizedCaller struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ResourceRegistryV1__UnauthorizedCaller(address caller)
func ResourceRegistryV1ResourceRegistryV1UnauthorizedCallerErrorID() common.Hash {
	return common.HexToHash("0x6543a023dac2837882104d3e6a1606af6a589938c1a575f8d97c3f385625f533")
}

// UnpackResourceRegistryV1UnauthorizedCallerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ResourceRegistryV1__UnauthorizedCaller(address caller)
func (resourceRegistryV1 *ResourceRegistryV1) UnpackResourceRegistryV1UnauthorizedCallerError(raw []byte) (*ResourceRegistryV1ResourceRegistryV1UnauthorizedCaller, error) {
	out := new(ResourceRegistryV1ResourceRegistryV1UnauthorizedCaller)
	if err := resourceRegistryV1.abi.UnpackIntoInterface(out, "ResourceRegistryV1UnauthorizedCaller", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ResourceRegistryV1UUPSUnauthorizedCallContext represents a UUPSUnauthorizedCallContext error raised by the ResourceRegistryV1 contract.
type ResourceRegistryV1UUPSUnauthorizedCallContext struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnauthorizedCallContext()
func ResourceRegistryV1UUPSUnauthorizedCallContextErrorID() common.Hash {
	return common.HexToHash("0xe07c8dba242a06571ac65fe4bbe20522c9fb111cb33599b799ff8039c1ed18f4")
}

// UnpackUUPSUnauthorizedCallContextError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnauthorizedCallContext()
func (resourceRegistryV1 *ResourceRegistryV1) UnpackUUPSUnauthorizedCallContextError(raw []byte) (*ResourceRegistryV1UUPSUnauthorizedCallContext, error) {
	out := new(ResourceRegistryV1UUPSUnauthorizedCallContext)
	if err := resourceRegistryV1.abi.UnpackIntoInterface(out, "UUPSUnauthorizedCallContext", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ResourceRegistryV1UUPSUnsupportedProxiableUUID represents a UUPSUnsupportedProxiableUUID error raised by the ResourceRegistryV1 contract.
type ResourceRegistryV1UUPSUnsupportedProxiableUUID struct {
	Slot [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func ResourceRegistryV1UUPSUnsupportedProxiableUUIDErrorID() common.Hash {
	return common.HexToHash("0xaa1d49a4c084bfa9aeeee2a0be65267a7f19ba7e1476b114dac513d2c14cb563")
}

// UnpackUUPSUnsupportedProxiableUUIDError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func (resourceRegistryV1 *ResourceRegistryV1) UnpackUUPSUnsupportedProxiableUUIDError(raw []byte) (*ResourceRegistryV1UUPSUnsupportedProxiableUUID, error) {
	out := new(ResourceRegistryV1UUPSUnsupportedProxiableUUID)
	if err := resourceRegistryV1.abi.UnpackIntoInterface(out, "UUPSUnsupportedProxiableUUID", raw); err != nil {
		return nil, err
	}
	return out, nil
}
