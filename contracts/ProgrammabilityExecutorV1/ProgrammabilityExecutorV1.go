// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ProgrammabilityExecutorV1

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

// ProgrammabilityExecutorV1MetaData contains all meta data concerning the ProgrammabilityExecutorV1 contract.
var ProgrammabilityExecutorV1MetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"contractVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"endpoint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRaylsEndpoint\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"executeProgramData\",\"inputs\":[{\"name\":\"steps\",\"type\":\"tuple[]\",\"internalType\":\"structSharedObjects.EnygmaProgramData[]\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"args\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"expectedMintTotal\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"originSender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAddressByResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_templateRegistryReplica\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"authority_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"resourceId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"templateRegistryReplica\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractITemplateRegistryReplica\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProgramDataExecuted\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"resourceId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"target\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"indexed\":false,\"internalType\":\"bytes4\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ProgramData__BothTargetsProvided\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProgramData__MintTotalMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProgramData__NoTargetProvided\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProgramData__Reverted\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"ret\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ProgramData__TooManyBlobs\",\"inputs\":[{\"name\":\"count\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProgramData__UnapprovedTemplate\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"bytecodeHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}]},{\"type\":\"error\",\"name\":\"ProgramData__UnknownResourceId\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ProgrammabilityExecutorV1__ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	ID:  "ProgrammabilityExecutorV1",
	Bin: "0x60a06040523060805234801561001457600080fd5b5060805161162361003e60003960008181610891015281816108ba01526109f301526116236000f3fe6080604052600436106100975760003560e01c806311f50c851461009c5780634f1ef286146100d257806352d1902d146100e75780635e280f111461010a5780635f997c5b1461012a5780636a0ed73d14610140578063a0a8e46014610160578063ad3cb1cc14610174578063bf7e214f146101b2578063c0c53b8b146101c7578063c4d66de8146101e7578063e6877d7314610207575b600080fd5b3480156100a857600080fd5b506100bc6100b73660046110a2565b610227565b6040516100c991906110bb565b60405180910390f35b6100e56100e03660046110fa565b61029b565b005b3480156100f357600080fd5b506100fc6102ba565b6040519081526020016100c9565b34801561011657600080fd5b506000546100bc906001600160a01b031681565b34801561013657600080fd5b506100fc60015481565b34801561014c57600080fd5b506002546100bc906001600160a01b031681565b34801561016c57600080fd5b5060016100fc565b34801561018057600080fd5b506101a5604051806040016040528060058152602001640352e302e360dc1b81525081565b6040516100c9919061120d565b3480156101be57600080fd5b506100bc6102d7565b3480156101d357600080fd5b506100e56101e2366004611220565b6102f0565b3480156101f357600080fd5b506100e561020236600461126b565b61045b565b34801561021357600080fd5b506100e5610222366004611288565b610485565b600080546040516311f50c8560e01b8152600481018490526001600160a01b03909116906311f50c8590602401602060405180830381865afa158015610271573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102959190611315565b92915050565b6102a3610886565b6102ac82610916565b6102b6828261092f565b5050565b60006102c46109e8565b506000805160206115ce83398151915290565b60006102e1610a31565b546001600160a01b0316919050565b60006102fa610a93565b805490915060ff600160401b82041615906001600160401b03166000811580156103215750825b90506000826001600160401b0316600114801561033d5750303b155b90508115801561034b575080155b156103695760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561039357845460ff60401b1916600160401b1785555b6001600160a01b03881615806103b057506001600160a01b038716155b156103ce57604051634e4990e360e11b815260040160405180910390fd5b6103d6610abc565b6103de610ac4565b6103e78861045b565b600280546001600160a01b0319166001600160a01b03891617905561040b86610ad4565b831561045157845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b5050505050505050565b610463610b15565b600080546001600160a01b0319166001600160a01b0392909216919091179055565b61049b336000356001600160e01b031916610b3a565b6104a3610c85565b6101008311156104ce57604051631fea353160e01b8152600481018490526024015b60405180910390fd5b6000805b8481101561084c5760008686838181106104ee576104ee611332565b90506020028101906105009190611348565b610511906060810190604001611368565b9050600087878481811061052757610527611332565b90506020028101906105399190611348565b610547906060810190611392565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201829052509394506105e392508691508b90508a8281811061059657610596611332565b90506020028101906105a89190611348565b358b8b888181106105bb576105bb611332565b90506020028101906105cd9190611348565b6105de90604081019060200161126b565b610cbb565b600254604051637db5f73160e01b81529192506001600160a01b031690637db5f7319061061690849087906004016113df565b602060405180830381865afa158015610633573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106579190611417565b6106975760405163ccdb0e9f60e01b8152600481018590526001600160a01b0382163f60248201526001600160e01b0319841660448201526064016104c5565b6060630f76402560e41b6001600160e01b0319851601610702576000838060200190518101906106c79190611432565b5091506106d69050818861147f565b965084846040516020016106eb929190611492565b60405160208183030381529060405291505061072c565b83838860601b60405160200161071a939291906114c3565b60405160208183030381529060405290505b600080836001600160a01b0316836040516107479190611509565b6000604051808303816000865af19150503d8060008114610784576040519150601f19603f3d011682016040523d82523d6000602084013e610789565b606091505b5091509150816107d857868c8c898181106107a6576107a6611332565b90506020028101906107b89190611348565b604051634f3853f960e11b81526104c5929135908990859060040161151b565b8b8b888181106107ea576107ea611332565b90506020028101906107fc9190611348565b60000135877f8689e39bf5f1db26858b4e0ae6745d8a87bb8b336e7aebca5d38c28a0ddb3f3986896040516108329291906113df565b60405180910390a35050600190940193506104d292505050565b508281146108775760405163781895a560e11b815260048101849052602481018290526044016104c5565b50610880610de3565b50505050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614806108f657507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166108ea610df4565b6001600160a01b031614155b156109145760405163703e46dd60e11b815260040160405180910390fd5b565b61092c336000356001600160e01b031916610b3a565b50565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610989575060408051601f3d908101601f1916820190925261098691810190611553565b60015b6109a85781604051634c9c8ce360e01b81526004016104c591906110bb565b6000805160206115ce83398151915281146109d957604051632a87526960e21b8152600481018290526024016104c5565b6109e38383610e0a565b505050565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146109145760405163703e46dd60e11b815260040160405180910390fd5b60008060ff19610a6260017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f3561156c565b604051602001610a7491815260200190565b60408051601f1981840301815291905280516020909101201692915050565b6000807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00610295565b610914610b15565b610acc610b15565b610914610e60565b6000610ade610a31565b80549091506001600160a01b031615610b0c5781604051638944034760e01b81526004016104c591906110bb565b6102b682610e68565b610b1d610ef8565b61091457604051631afcd79f60e31b815260040160405180910390fd5b6000610b44610a31565b80549091506001600160a01b031680610b73576000604051638944034760e01b81526004016104c591906110bb565b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa158015610bd7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bfb919061157f565b92509250925082610c7c578015610c255760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff821615610c615760405163a426878960e01b81526001600160a01b038816600482015263ffffffff831660248201526044016104c5565b86604051632ecd3d0360e21b81526004016104c591906110bb565b50505050505050565b6000610c8f610f12565b805490915060011901610cb557604051633ee5aeb560e01b815260040160405180910390fd5b60029055565b600082158015906001600160a01b0384161515908290610cd85750805b15610cfc5760405160016201398560e41b03198152600481018790526024016104c5565b81158015610d08575080155b15610d295760405163046a1f2960e31b8152600481018790526024016104c5565b8115610dd6576000546040516311f50c8560e01b8152600481018790526001600160a01b03909116906311f50c8590602401602060405180830381865afa158015610d78573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d9c9190611315565b92506001600160a01b038316610dcf5760405163cdbe10bf60e01b815260048101879052602481018690526044016104c5565b5050610ddc565b83925050505b9392505050565b6000610ded610f12565b6001905550565b60006000805160206115ce8339815191526102e1565b610e1382610f36565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a2805115610e58576109e38282610f92565b6102b6611008565b610de3610b15565b6001600160a01b038116610e915780604051638944034760e01b81526004016104c591906110bb565b6000610e9b610a31565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b6000610f02610a93565b54600160401b900460ff16919050565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0090565b806001600160a01b03163b600003610f635780604051634c9c8ce360e01b81526004016104c591906110bb565b6000805160206115ce83398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080846001600160a01b031684604051610faf9190611509565b600060405180830381855af49150503d8060008114610fea576040519150601f19603f3d011682016040523d82523d6000602084013e610fef565b606091505b5091509150610fff858383611027565b95945050505050565b34156109145760405163b398979f60e01b815260040160405180910390fd5b60608261103c576110378261107a565b610ddc565b815115801561105357506001600160a01b0384163b155b156110735783604051639996b31560e01b81526004016104c591906110bb565b5080610ddc565b80511561108957805160208201fd5b60405163d6bda27560e01b815260040160405180910390fd5b6000602082840312156110b457600080fd5b5035919050565b6001600160a01b0391909116815260200190565b6001600160a01b038116811461092c57600080fd5b634e487b7160e01b600052604160045260246000fd5b6000806040838503121561110d57600080fd5b8235611118816110cf565b915060208301356001600160401b038082111561113457600080fd5b818501915085601f83011261114857600080fd5b81358181111561115a5761115a6110e4565b604051601f8201601f19908116603f01168101908382118183101715611182576111826110e4565b8160405282815288602084870101111561119b57600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b60005b838110156111d85781810151838201526020016111c0565b50506000910152565b600081518084526111f98160208601602086016111bd565b601f01601f19169290920160200192915050565b602081526000610ddc60208301846111e1565b60008060006060848603121561123557600080fd5b8335611240816110cf565b92506020840135611250816110cf565b91506040840135611260816110cf565b809150509250925092565b60006020828403121561127d57600080fd5b8135610ddc816110cf565b6000806000806060858703121561129e57600080fd5b84356001600160401b03808211156112b557600080fd5b818701915087601f8301126112c957600080fd5b8135818111156112d857600080fd5b8860208260051b85010111156112ed57600080fd5b602092830196509450508501359150604085013561130a816110cf565b939692955090935050565b60006020828403121561132757600080fd5b8151610ddc816110cf565b634e487b7160e01b600052603260045260246000fd5b60008235607e1983360301811261135e57600080fd5b9190910192915050565b60006020828403121561137a57600080fd5b81356001600160e01b031981168114610ddc57600080fd5b6000808335601e198436030181126113a957600080fd5b8301803591506001600160401b038211156113c357600080fd5b6020019150368190038213156113d857600080fd5b9250929050565b6001600160a01b039290921682526001600160e01b031916602082015260400190565b8051801515811461141257600080fd5b919050565b60006020828403121561142957600080fd5b610ddc82611402565b60008060006060848603121561144757600080fd5b8351611452816110cf565b602085015160409095015190969495509392505050565b634e487b7160e01b600052601160045260246000fd5b8082018082111561029557610295611469565b6001600160e01b03198316815281516000906114b58160048501602087016111bd565b919091016004019392505050565b6001600160e01b03198416815282516000906114e68160048501602088016111bd565b6001600160601b0319939093169190920160048101919091526018019392505050565b6000825161135e8184602087016111bd565b84815283602082015263ffffffff60e01b8316604082015260806060820152600061154960808301846111e1565b9695505050505050565b60006020828403121561156557600080fd5b5051919050565b8181038181111561029557610295611469565b60008060006060848603121561159457600080fd5b61159d84611402565b9250602084015163ffffffff811681146115b657600080fd5b91506115c460408501611402565b9050925092509256fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbca264697066735822122089dce4cfc8f90f2e9bedcfc242454b2973a02a46b8c51d733e8b369ab1ae6eb964736f6c63430008180033",
}

// ProgrammabilityExecutorV1 is an auto generated Go binding around an Ethereum contract.
type ProgrammabilityExecutorV1 struct {
	abi abi.ABI
}

// NewProgrammabilityExecutorV1 creates a new instance of ProgrammabilityExecutorV1.
func NewProgrammabilityExecutorV1() *ProgrammabilityExecutorV1 {
	parsed, err := ProgrammabilityExecutorV1MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ProgrammabilityExecutorV1{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ProgrammabilityExecutorV1) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackUPGRADEINTERFACEVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) PackUPGRADEINTERFACEVERSION() []byte {
	enc, err := programmabilityExecutorV1.abi.Pack("UPGRADE_INTERFACE_VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUPGRADEINTERFACEVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackUPGRADEINTERFACEVERSION(data []byte) (string, error) {
	out, err := programmabilityExecutorV1.abi.Unpack("UPGRADE_INTERFACE_VERSION", data)
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
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) PackAuthority() []byte {
	enc, err := programmabilityExecutorV1.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := programmabilityExecutorV1.abi.Unpack("authority", data)
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
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) PackContractVersion() []byte {
	enc, err := programmabilityExecutorV1.abi.Pack("contractVersion")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackContractVersion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa0a8e460.
//
// Solidity: function contractVersion() pure returns(uint256)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackContractVersion(data []byte) (*big.Int, error) {
	out, err := programmabilityExecutorV1.abi.Unpack("contractVersion", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackEndpoint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) PackEndpoint() []byte {
	enc, err := programmabilityExecutorV1.abi.Pack("endpoint")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackEndpoint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackEndpoint(data []byte) (common.Address, error) {
	out, err := programmabilityExecutorV1.abi.Unpack("endpoint", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackExecuteProgramData is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe6877d73.
//
// Solidity: function executeProgramData((bytes32,address,bytes4,bytes)[] steps, uint256 expectedMintTotal, address originSender) returns()
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) PackExecuteProgramData(steps []SharedObjectsEnygmaProgramData, expectedMintTotal *big.Int, originSender common.Address) []byte {
	enc, err := programmabilityExecutorV1.abi.Pack("executeProgramData", steps, expectedMintTotal, originSender)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackGetAddressByResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) PackGetAddressByResourceId(resourceId [32]byte) []byte {
	enc, err := programmabilityExecutorV1.abi.Pack("getAddressByResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAddressByResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackGetAddressByResourceId(data []byte) (common.Address, error) {
	out, err := programmabilityExecutorV1.abi.Unpack("getAddressByResourceId", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc0c53b8b.
//
// Solidity: function initialize(address _endpoint, address _templateRegistryReplica, address authority_) returns()
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) PackInitialize(endpoint common.Address, templateRegistryReplica common.Address, authority common.Address) []byte {
	enc, err := programmabilityExecutorV1.abi.Pack("initialize", endpoint, templateRegistryReplica, authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackInitialize0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc4d66de8.
//
// Solidity: function initialize(address _endpoint) returns()
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) PackInitialize0(endpoint common.Address) []byte {
	enc, err := programmabilityExecutorV1.abi.Pack("initialize0", endpoint)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackProxiableUUID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) PackProxiableUUID() []byte {
	enc, err := programmabilityExecutorV1.abi.Pack("proxiableUUID")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackProxiableUUID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackProxiableUUID(data []byte) ([32]byte, error) {
	out, err := programmabilityExecutorV1.abi.Unpack("proxiableUUID", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) PackResourceId() []byte {
	enc, err := programmabilityExecutorV1.abi.Pack("resourceId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackResourceId(data []byte) ([32]byte, error) {
	out, err := programmabilityExecutorV1.abi.Unpack("resourceId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackTemplateRegistryReplica is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6a0ed73d.
//
// Solidity: function templateRegistryReplica() view returns(address)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) PackTemplateRegistryReplica() []byte {
	enc, err := programmabilityExecutorV1.abi.Pack("templateRegistryReplica")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTemplateRegistryReplica is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6a0ed73d.
//
// Solidity: function templateRegistryReplica() view returns(address)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackTemplateRegistryReplica(data []byte) (common.Address, error) {
	out, err := programmabilityExecutorV1.abi.Unpack("templateRegistryReplica", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackUpgradeToAndCall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) PackUpgradeToAndCall(newImplementation common.Address, data []byte) []byte {
	enc, err := programmabilityExecutorV1.abi.Pack("upgradeToAndCall", newImplementation, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// ProgrammabilityExecutorV1AuthorityUpdated represents a AuthorityUpdated event raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1AuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const ProgrammabilityExecutorV1AuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (ProgrammabilityExecutorV1AuthorityUpdated) ContractEventName() string {
	return ProgrammabilityExecutorV1AuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackAuthorityUpdatedEvent(log *types.Log) (*ProgrammabilityExecutorV1AuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != programmabilityExecutorV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProgrammabilityExecutorV1AuthorityUpdated)
	if len(log.Data) > 0 {
		if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range programmabilityExecutorV1.abi.Events[event].Inputs {
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

// ProgrammabilityExecutorV1Initialized represents a Initialized event raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1Initialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const ProgrammabilityExecutorV1InitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (ProgrammabilityExecutorV1Initialized) ContractEventName() string {
	return ProgrammabilityExecutorV1InitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackInitializedEvent(log *types.Log) (*ProgrammabilityExecutorV1Initialized, error) {
	event := "Initialized"
	if log.Topics[0] != programmabilityExecutorV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProgrammabilityExecutorV1Initialized)
	if len(log.Data) > 0 {
		if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range programmabilityExecutorV1.abi.Events[event].Inputs {
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

// ProgrammabilityExecutorV1ProgramDataExecuted represents a ProgramDataExecuted event raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1ProgramDataExecuted struct {
	Index      *big.Int
	ResourceId [32]byte
	Target     common.Address
	Selector   [4]byte
	Raw        *types.Log // Blockchain specific contextual infos
}

const ProgrammabilityExecutorV1ProgramDataExecutedEventName = "ProgramDataExecuted"

// ContractEventName returns the user-defined event name.
func (ProgrammabilityExecutorV1ProgramDataExecuted) ContractEventName() string {
	return ProgrammabilityExecutorV1ProgramDataExecutedEventName
}

// UnpackProgramDataExecutedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ProgramDataExecuted(uint256 indexed index, bytes32 indexed resourceId, address target, bytes4 selector)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackProgramDataExecutedEvent(log *types.Log) (*ProgrammabilityExecutorV1ProgramDataExecuted, error) {
	event := "ProgramDataExecuted"
	if log.Topics[0] != programmabilityExecutorV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProgrammabilityExecutorV1ProgramDataExecuted)
	if len(log.Data) > 0 {
		if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range programmabilityExecutorV1.abi.Events[event].Inputs {
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

// ProgrammabilityExecutorV1Upgraded represents a Upgraded event raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1Upgraded struct {
	Implementation common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const ProgrammabilityExecutorV1UpgradedEventName = "Upgraded"

// ContractEventName returns the user-defined event name.
func (ProgrammabilityExecutorV1Upgraded) ContractEventName() string {
	return ProgrammabilityExecutorV1UpgradedEventName
}

// UnpackUpgradedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Upgraded(address indexed implementation)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackUpgradedEvent(log *types.Log) (*ProgrammabilityExecutorV1Upgraded, error) {
	event := "Upgraded"
	if log.Topics[0] != programmabilityExecutorV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProgrammabilityExecutorV1Upgraded)
	if len(log.Data) > 0 {
		if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range programmabilityExecutorV1.abi.Events[event].Inputs {
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
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["AddressEmptyCode"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackAddressEmptyCodeError(raw[4:])
	}
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["FailedCall"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackFailedCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["ProgramDataBothTargetsProvided"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackProgramDataBothTargetsProvidedError(raw[4:])
	}
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["ProgramDataMintTotalMismatch"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackProgramDataMintTotalMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["ProgramDataNoTargetProvided"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackProgramDataNoTargetProvidedError(raw[4:])
	}
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["ProgramDataReverted"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackProgramDataRevertedError(raw[4:])
	}
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["ProgramDataTooManyBlobs"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackProgramDataTooManyBlobsError(raw[4:])
	}
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["ProgramDataUnapprovedTemplate"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackProgramDataUnapprovedTemplateError(raw[4:])
	}
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["ProgramDataUnknownResourceId"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackProgramDataUnknownResourceIdError(raw[4:])
	}
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["ProgrammabilityExecutorV1ZeroAddress"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackProgrammabilityExecutorV1ZeroAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["ReentrancyGuardReentrantCall"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackReentrancyGuardReentrantCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["UUPSUnauthorizedCallContext"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackUUPSUnauthorizedCallContextError(raw[4:])
	}
	if bytes.Equal(raw[:4], programmabilityExecutorV1.abi.Errors["UUPSUnsupportedProxiableUUID"].ID.Bytes()[:4]) {
		return programmabilityExecutorV1.UnpackUUPSUnsupportedProxiableUUIDError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ProgrammabilityExecutorV1AddressEmptyCode represents a AddressEmptyCode error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1AddressEmptyCode struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressEmptyCode(address target)
func ProgrammabilityExecutorV1AddressEmptyCodeErrorID() common.Hash {
	return common.HexToHash("0x9996b315c842ff135b8fc4a08ad5df1c344efbc03d2687aecc0678050d2aac89")
}

// UnpackAddressEmptyCodeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressEmptyCode(address target)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackAddressEmptyCodeError(raw []byte) (*ProgrammabilityExecutorV1AddressEmptyCode, error) {
	out := new(ProgrammabilityExecutorV1AddressEmptyCode)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "AddressEmptyCode", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProgrammabilityExecutorV1ERC1967InvalidImplementation represents a ERC1967InvalidImplementation error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1ERC1967InvalidImplementation struct {
	Implementation common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func ProgrammabilityExecutorV1ERC1967InvalidImplementationErrorID() common.Hash {
	return common.HexToHash("0x4c9c8ce3ceb3130f17f7cdba48d89b5b0129f266a8bac114e6e315a41879b617")
}

// UnpackERC1967InvalidImplementationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackERC1967InvalidImplementationError(raw []byte) (*ProgrammabilityExecutorV1ERC1967InvalidImplementation, error) {
	out := new(ProgrammabilityExecutorV1ERC1967InvalidImplementation)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "ERC1967InvalidImplementation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProgrammabilityExecutorV1ERC1967NonPayable represents a ERC1967NonPayable error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1ERC1967NonPayable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967NonPayable()
func ProgrammabilityExecutorV1ERC1967NonPayableErrorID() common.Hash {
	return common.HexToHash("0xb398979fa84f543c8e222f17890372c487baf85e062276c127fef521eea7224b")
}

// UnpackERC1967NonPayableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967NonPayable()
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackERC1967NonPayableError(raw []byte) (*ProgrammabilityExecutorV1ERC1967NonPayable, error) {
	out := new(ProgrammabilityExecutorV1ERC1967NonPayable)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "ERC1967NonPayable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProgrammabilityExecutorV1FailedCall represents a FailedCall error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1FailedCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedCall()
func ProgrammabilityExecutorV1FailedCallErrorID() common.Hash {
	return common.HexToHash("0xd6bda27508c0fb6d8a39b4b122878dab26f731a7d4e4abe711dd3731899052a4")
}

// UnpackFailedCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedCall()
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackFailedCallError(raw []byte) (*ProgrammabilityExecutorV1FailedCall, error) {
	out := new(ProgrammabilityExecutorV1FailedCall)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "FailedCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProgrammabilityExecutorV1InvalidInitialization represents a InvalidInitialization error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1InvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func ProgrammabilityExecutorV1InvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackInvalidInitializationError(raw []byte) (*ProgrammabilityExecutorV1InvalidInitialization, error) {
	out := new(ProgrammabilityExecutorV1InvalidInitialization)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProgrammabilityExecutorV1NotInitializing represents a NotInitializing error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1NotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func ProgrammabilityExecutorV1NotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackNotInitializingError(raw []byte) (*ProgrammabilityExecutorV1NotInitializing, error) {
	out := new(ProgrammabilityExecutorV1NotInitializing)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProgrammabilityExecutorV1ProgramDataBothTargetsProvided represents a ProgramData__BothTargetsProvided error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1ProgramDataBothTargetsProvided struct {
	Index *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ProgramData__BothTargetsProvided(uint256 index)
func ProgrammabilityExecutorV1ProgramDataBothTargetsProvidedErrorID() common.Hash {
	return common.HexToHash("0xffec67b06e167fe03e4bb1e8d8776c3cdfb49a388a1bc65f148bb945f0153774")
}

// UnpackProgramDataBothTargetsProvidedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ProgramData__BothTargetsProvided(uint256 index)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackProgramDataBothTargetsProvidedError(raw []byte) (*ProgrammabilityExecutorV1ProgramDataBothTargetsProvided, error) {
	out := new(ProgrammabilityExecutorV1ProgramDataBothTargetsProvided)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "ProgramDataBothTargetsProvided", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProgrammabilityExecutorV1ProgramDataMintTotalMismatch represents a ProgramData__MintTotalMismatch error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1ProgramDataMintTotalMismatch struct {
	Expected *big.Int
	Actual   *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ProgramData__MintTotalMismatch(uint256 expected, uint256 actual)
func ProgrammabilityExecutorV1ProgramDataMintTotalMismatchErrorID() common.Hash {
	return common.HexToHash("0xf0312b4a7b2aebcb8f28b1648782a5492688b7e5f3920994f3c53a4f900562e8")
}

// UnpackProgramDataMintTotalMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ProgramData__MintTotalMismatch(uint256 expected, uint256 actual)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackProgramDataMintTotalMismatchError(raw []byte) (*ProgrammabilityExecutorV1ProgramDataMintTotalMismatch, error) {
	out := new(ProgrammabilityExecutorV1ProgramDataMintTotalMismatch)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "ProgramDataMintTotalMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProgrammabilityExecutorV1ProgramDataNoTargetProvided represents a ProgramData__NoTargetProvided error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1ProgramDataNoTargetProvided struct {
	Index *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ProgramData__NoTargetProvided(uint256 index)
func ProgrammabilityExecutorV1ProgramDataNoTargetProvidedErrorID() common.Hash {
	return common.HexToHash("0x2350f948358ff2bcafe4b9efcb56d7c8244da1c6c1d8415c660338a228120558")
}

// UnpackProgramDataNoTargetProvidedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ProgramData__NoTargetProvided(uint256 index)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackProgramDataNoTargetProvidedError(raw []byte) (*ProgrammabilityExecutorV1ProgramDataNoTargetProvided, error) {
	out := new(ProgrammabilityExecutorV1ProgramDataNoTargetProvided)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "ProgramDataNoTargetProvided", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProgrammabilityExecutorV1ProgramDataReverted represents a ProgramData__Reverted error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1ProgramDataReverted struct {
	Index      *big.Int
	ResourceId [32]byte
	Selector   [4]byte
	Ret        []byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ProgramData__Reverted(uint256 index, bytes32 resourceId, bytes4 selector, bytes ret)
func ProgrammabilityExecutorV1ProgramDataRevertedErrorID() common.Hash {
	return common.HexToHash("0x9e70a7f2d97749b0f96c842adb00e195f10ad941f92af8a67d22e3b32063b5ea")
}

// UnpackProgramDataRevertedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ProgramData__Reverted(uint256 index, bytes32 resourceId, bytes4 selector, bytes ret)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackProgramDataRevertedError(raw []byte) (*ProgrammabilityExecutorV1ProgramDataReverted, error) {
	out := new(ProgrammabilityExecutorV1ProgramDataReverted)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "ProgramDataReverted", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProgrammabilityExecutorV1ProgramDataTooManyBlobs represents a ProgramData__TooManyBlobs error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1ProgramDataTooManyBlobs struct {
	Count *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ProgramData__TooManyBlobs(uint256 count)
func ProgrammabilityExecutorV1ProgramDataTooManyBlobsErrorID() common.Hash {
	return common.HexToHash("0x1fea3531678aee8e5fe9206ff8fb3e76a4ae52edca545d57d4170a5dd9171f0b")
}

// UnpackProgramDataTooManyBlobsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ProgramData__TooManyBlobs(uint256 count)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackProgramDataTooManyBlobsError(raw []byte) (*ProgrammabilityExecutorV1ProgramDataTooManyBlobs, error) {
	out := new(ProgrammabilityExecutorV1ProgramDataTooManyBlobs)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "ProgramDataTooManyBlobs", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProgrammabilityExecutorV1ProgramDataUnapprovedTemplate represents a ProgramData__UnapprovedTemplate error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1ProgramDataUnapprovedTemplate struct {
	Index        *big.Int
	BytecodeHash [32]byte
	Selector     [4]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ProgramData__UnapprovedTemplate(uint256 index, bytes32 bytecodeHash, bytes4 selector)
func ProgrammabilityExecutorV1ProgramDataUnapprovedTemplateErrorID() common.Hash {
	return common.HexToHash("0xccdb0e9f6d4741f701b4dd8ea7a0aaac9f0b9c8784cf63d54e45e5154c2bec6d")
}

// UnpackProgramDataUnapprovedTemplateError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ProgramData__UnapprovedTemplate(uint256 index, bytes32 bytecodeHash, bytes4 selector)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackProgramDataUnapprovedTemplateError(raw []byte) (*ProgrammabilityExecutorV1ProgramDataUnapprovedTemplate, error) {
	out := new(ProgrammabilityExecutorV1ProgramDataUnapprovedTemplate)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "ProgramDataUnapprovedTemplate", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProgrammabilityExecutorV1ProgramDataUnknownResourceId represents a ProgramData__UnknownResourceId error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1ProgramDataUnknownResourceId struct {
	Index      *big.Int
	ResourceId [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ProgramData__UnknownResourceId(uint256 index, bytes32 resourceId)
func ProgrammabilityExecutorV1ProgramDataUnknownResourceIdErrorID() common.Hash {
	return common.HexToHash("0xcdbe10bfdee94179945f075946e6cc677edeba9832aec4ee48829cc45a5e9134")
}

// UnpackProgramDataUnknownResourceIdError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ProgramData__UnknownResourceId(uint256 index, bytes32 resourceId)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackProgramDataUnknownResourceIdError(raw []byte) (*ProgrammabilityExecutorV1ProgramDataUnknownResourceId, error) {
	out := new(ProgrammabilityExecutorV1ProgramDataUnknownResourceId)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "ProgramDataUnknownResourceId", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProgrammabilityExecutorV1ProgrammabilityExecutorV1ZeroAddress represents a ProgrammabilityExecutorV1__ZeroAddress error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1ProgrammabilityExecutorV1ZeroAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ProgrammabilityExecutorV1__ZeroAddress()
func ProgrammabilityExecutorV1ProgrammabilityExecutorV1ZeroAddressErrorID() common.Hash {
	return common.HexToHash("0x9c9321c6fd25c09b920629252863a62ee14161d2decaf77147e2923cfcff4fac")
}

// UnpackProgrammabilityExecutorV1ZeroAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ProgrammabilityExecutorV1__ZeroAddress()
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackProgrammabilityExecutorV1ZeroAddressError(raw []byte) (*ProgrammabilityExecutorV1ProgrammabilityExecutorV1ZeroAddress, error) {
	out := new(ProgrammabilityExecutorV1ProgrammabilityExecutorV1ZeroAddress)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "ProgrammabilityExecutorV1ZeroAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProgrammabilityExecutorV1RaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1RaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func ProgrammabilityExecutorV1RaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*ProgrammabilityExecutorV1RaylsAccessManagedContractPaused, error) {
	out := new(ProgrammabilityExecutorV1RaylsAccessManagedContractPaused)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProgrammabilityExecutorV1RaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1RaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func ProgrammabilityExecutorV1RaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*ProgrammabilityExecutorV1RaylsAccessManagedInvalidAuthority, error) {
	out := new(ProgrammabilityExecutorV1RaylsAccessManagedInvalidAuthority)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProgrammabilityExecutorV1RaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1RaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func ProgrammabilityExecutorV1RaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*ProgrammabilityExecutorV1RaylsAccessManagedMustSchedule, error) {
	out := new(ProgrammabilityExecutorV1RaylsAccessManagedMustSchedule)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProgrammabilityExecutorV1RaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1RaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func ProgrammabilityExecutorV1RaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*ProgrammabilityExecutorV1RaylsAccessManagedUnauthorized, error) {
	out := new(ProgrammabilityExecutorV1RaylsAccessManagedUnauthorized)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProgrammabilityExecutorV1ReentrancyGuardReentrantCall represents a ReentrancyGuardReentrantCall error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1ReentrancyGuardReentrantCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ReentrancyGuardReentrantCall()
func ProgrammabilityExecutorV1ReentrancyGuardReentrantCallErrorID() common.Hash {
	return common.HexToHash("0x3ee5aeb571de7fc460830b4d0017439a1ca56fb0bc39062227ade4fe4a24c1ca")
}

// UnpackReentrancyGuardReentrantCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ReentrancyGuardReentrantCall()
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackReentrancyGuardReentrantCallError(raw []byte) (*ProgrammabilityExecutorV1ReentrancyGuardReentrantCall, error) {
	out := new(ProgrammabilityExecutorV1ReentrancyGuardReentrantCall)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "ReentrancyGuardReentrantCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProgrammabilityExecutorV1UUPSUnauthorizedCallContext represents a UUPSUnauthorizedCallContext error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1UUPSUnauthorizedCallContext struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnauthorizedCallContext()
func ProgrammabilityExecutorV1UUPSUnauthorizedCallContextErrorID() common.Hash {
	return common.HexToHash("0xe07c8dba242a06571ac65fe4bbe20522c9fb111cb33599b799ff8039c1ed18f4")
}

// UnpackUUPSUnauthorizedCallContextError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnauthorizedCallContext()
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackUUPSUnauthorizedCallContextError(raw []byte) (*ProgrammabilityExecutorV1UUPSUnauthorizedCallContext, error) {
	out := new(ProgrammabilityExecutorV1UUPSUnauthorizedCallContext)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "UUPSUnauthorizedCallContext", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProgrammabilityExecutorV1UUPSUnsupportedProxiableUUID represents a UUPSUnsupportedProxiableUUID error raised by the ProgrammabilityExecutorV1 contract.
type ProgrammabilityExecutorV1UUPSUnsupportedProxiableUUID struct {
	Slot [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func ProgrammabilityExecutorV1UUPSUnsupportedProxiableUUIDErrorID() common.Hash {
	return common.HexToHash("0xaa1d49a4c084bfa9aeeee2a0be65267a7f19ba7e1476b114dac513d2c14cb563")
}

// UnpackUUPSUnsupportedProxiableUUIDError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func (programmabilityExecutorV1 *ProgrammabilityExecutorV1) UnpackUUPSUnsupportedProxiableUUIDError(raw []byte) (*ProgrammabilityExecutorV1UUPSUnsupportedProxiableUUID, error) {
	out := new(ProgrammabilityExecutorV1UUPSUnsupportedProxiableUUID)
	if err := programmabilityExecutorV1.abi.UnpackIntoInterface(out, "UUPSUnsupportedProxiableUUID", raw); err != nil {
		return nil, err
	}
	return out, nil
}
