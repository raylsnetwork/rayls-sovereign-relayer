// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TemplateRegistryReplicaV1

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

// ITemplateRegistryTemplate is an auto generated low-level Go binding around an user-defined struct.
type ITemplateRegistryTemplate struct {
	BytecodeHash    [32]byte
	Signature       string
	Selector        [4]byte
	ApprovedAtBlock uint64
	Approved        bool
}

// TemplateRegistryReplicaV1MetaData contains all meta data concerning the TemplateRegistryReplicaV1 contract.
var TemplateRegistryReplicaV1MetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"check\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkWithParamCount\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"paramCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"contractVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"endpoint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRaylsEndpoint\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAddressByResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getKey\",\"inputs\":[{\"name\":\"bytecodeHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getLastUpdatedAt\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTemplate\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structITemplateRegistry.Template\",\"components\":[{\"name\":\"bytecodeHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"approvedAtBlock\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"authority_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onTemplateApproved\",\"inputs\":[{\"name\":\"bytecodeHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"approvedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onTemplateRevoked\",\"inputs\":[{\"name\":\"bytecodeHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"revokedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"resourceId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TemplateMirrored\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"bytecodeHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"indexed\":false,\"internalType\":\"bytes4\"},{\"name\":\"approvedAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TemplateRevocationMirrored\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"bytecodeHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"indexed\":false,\"internalType\":\"bytes4\"},{\"name\":\"revokedAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SignatureParams__Malformed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TemplateRegistryReplicaV1__EmptyBytecodeHash\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TemplateRegistryReplicaV1__EmptySignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TemplateRegistryReplicaV1__NotFromPrivateHub\",\"inputs\":[{\"name\":\"fromChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"hubId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	ID:  "TemplateRegistryReplicaV1",
	Bin: "0x60a06040523060805234801561001457600080fd5b50608051611b3d61003e60003960008181610c0601528181610c2f0152610d610152611b3d6000f3fe6080604052600436106100de5760003560e01c80637db5f731116100855780637db5f7311461021557806388ee8c01146102455780639757739b14610265578063a0a8e46014610292578063ad3cb1cc146102a6578063af368e84146102e4578063bf7e214f14610304578063c4d66de814610319578063fee8c9ac1461033957600080fd5b806311f50c85146100e3578063397349df146101195780633dd9380214610167578063485cc955146101955780634f1ef286146101b757806352d1902d146101ca5780635e280f11146101df5780635f997c5b146101ff575b600080fd5b3480156100ef57600080fd5b506101036100fe366004611498565b610370565b60405161011091906114b1565b60405180910390f35b34801561012557600080fd5b5061014f610134366004611498565b6000908152600360205260409020546001600160401b031690565b6040516001600160401b039091168152602001610110565b34801561017357600080fd5b506101876101823660046114e2565b6103e4565b604051908152602001610110565b3480156101a157600080fd5b506101b56101b0366004611523565b610417565b005b6101b56101c5366004611572565b610528565b3480156101d657600080fd5b50610187610547565b3480156101eb57600080fd5b50600054610103906001600160a01b031681565b34801561020b57600080fd5b5061018760015481565b34801561022157600080fd5b50610235610230366004611635565b610565565b6040519015158152602001610110565b34801561025157600080fd5b506101b5610260366004611678565b6105bd565b34801561027157600080fd5b50610285610280366004611498565b6106cf565b6040516101109190611704565b34801561029e57600080fd5b506001610187565b3480156102b257600080fd5b506102d7604051806040016040528060058152602001640352e302e360dc1b81525081565b6040516101109190611768565b3480156102f057600080fd5b506101b56102ff36600461177b565b6107fc565b34801561031057600080fd5b50610103610a2a565b34801561032557600080fd5b506101b5610334366004611807565b610a43565b34801561034557600080fd5b50610359610354366004611635565b610a6d565b604080519215158352602083019190915201610110565b600080546040516311f50c8560e01b8152600481018490526001600160a01b03909116906311f50c8590602401602060405180830381865afa1580156103ba573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103de9190611824565b92915050565b600082826040516020016103f9929190611841565b60405160208183030381529060405280519060200120905092915050565b6000610421610b7e565b805490915060ff600160401b82041615906001600160401b03166000811580156104485750825b90506000826001600160401b031660011480156104645750303b155b905081158015610472575080155b156104905760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156104ba57845460ff60401b1916600160401b1785555b6104c2610ba7565b6104cb87610a43565b60066001556104d986610bb1565b831561051f57845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50505050505050565b610530610bfb565b61053982610c89565b6105438282610ca2565b5050565b6000610551610d56565b50600080516020611ae88339815191525b90565b600080836001600160a01b03163f83604051602001610585929190611841565b60408051808303601f1901815291815281516020928301206000908152600292839052200154600160601b900460ff16949350505050565b6105d3336000356001600160e01b031916610d9f565b6105db610ee1565b600083836040516020016105f0929190611841565b60408051601f198184030181529181528151602092830120600081815260039093529120549091506001600160401b03908116908316116106315750505050565b600081815260026020818152604080842090920180546cffffffffffffffffff000000001916600160201b6001600160401b03881690810291909117909155600390915291819020805467ffffffffffffffff19169092179091555181907f8ecefc15bacb9e902f69507d1764943a0a43f3d916b1dc50f5fb7bf98854b1d9906106c090879087908790611859565b60405180910390a2505b505050565b6040805160a08101825260008082526060602083018190529282018190529181018290526080810191909152600260008381526020019081526020016000206040518060a00160405290816000820154815260200160018201805461073390611882565b80601f016020809104026020016040519081016040528092919081815260200182805461075f90611882565b80156107ac5780601f10610781576101008083540402835291602001916107ac565b820191906000526020600020905b81548152906001019060200180831161078f57829003601f168201915b50505091835250506002919091015460e081901b6001600160e01b0319166020830152600160201b81046001600160401b03166040830152600160601b900460ff16151560609091015292915050565b610812336000356001600160e01b031916610d9f565b61081a610ee1565b8361083857604051634230fe6b60e01b815260040160405180910390fd5b600082900361085a57604051634da6cb2f60e01b815260040160405180910390fd5b6000838360405161086c9291906118bc565b604051809103902090506000858260405160200161088b929190611841565b60408051601f198184030181529181528151602092830120600081815260039093529120549091506001600160401b03908116908416116108cd575050610a24565b6040518060a0016040528087815260200186868080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201829052509385525050506001600160e01b031985166020808401919091526001600160401b038716604080850191909152600160609094018490528583526002825290912083518155908301519091820190610968908261191c565b506040828101516002929092018054606085015160809095015160e09490941c6001600160601b031990911617600160201b6001600160401b03958616021760ff60601b1916600160601b931515939093029290921790915560008381526003602052819020805467ffffffffffffffff1916928616929092179091555181907f91cf583ebe8e6e22c296b576580931263a234956d20029d76af4a9e4f324303290610a1990899086908890611859565b60405180910390a250505b50505050565b6000610a34610f91565b546001600160a01b0316919050565b610a4b610ff3565b600080546001600160a01b0319166001600160a01b0392909216919091179055565b6000806000846001600160a01b03163f84604051602001610a8f929190611841565b60408051808303601f19018152918152815160209283012060008181526002938490529190912091820154909250600160601b900460ff16610ad957600080935093505050610b77565b6001610b70826001018054610aed90611882565b80601f0160208091040260200160405190810160405280929190818152602001828054610b1990611882565b8015610b665780601f10610b3b57610100808354040283529160200191610b66565b820191906000526020600020905b815481529060010190602001808311610b4957829003601f168201915b5050505050611018565b9350935050505b9250929050565b6000807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a006103de565b610baf610ff3565b565b6000610bbb610f91565b80549091506001600160a01b031615610bf25781604051638944034760e01b8152600401610be991906114b1565b60405180910390fd5b610543826111ff565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161480610c6b57507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316610c5f61128f565b6001600160a01b031614155b15610baf5760405163703e46dd60e11b815260040160405180910390fd5b610c9f336000356001600160e01b031916610d9f565b50565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610cfc575060408051601f3d908101601f19168201909252610cf9918101906119db565b60015b610d1b5781604051634c9c8ce360e01b8152600401610be991906114b1565b600080516020611ae88339815191528114610d4c57604051632a87526960e21b815260048101829052602401610be9565b6106ca83836112a5565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614610baf5760405163703e46dd60e11b815260040160405180910390fd5b6000610da9610f91565b80549091506001600160a01b031680610dd8576000604051638944034760e01b8152600401610be991906114b1565b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa158015610e3c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e609190611a04565b9250925092508261051f578015610e8a5760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff821615610ec65760405163a426878960e01b81526001600160a01b038816600482015263ffffffff83166024820152604401610be9565b86604051632ecd3d0360e21b8152600401610be991906114b1565b6000610eeb6112fb565b905060008060009054906101000a90046001600160a01b03166001600160a01b0316630b39a9516040518163ffffffff1660e01b8152600401602060405180830381865afa158015610f41573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f6591906119db565b905080821461054357604051632430651760e21b81526004810183905260248101829052604401610be9565b60008060ff19610fc260017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f35611a5f565b604051602001610fd491815260200190565b60408051601f1981840301815291905280516020909101201692915050565b610ffb61130f565b610baf57604051631afcd79f60e31b815260040160405180910390fd5b80516000908290600019835b828110156110675783818151811061103e5761103e611a72565b01602001516001600160f81b031916600560fb1b0361105f57809150611067565b600101611024565b506000198114806110a757508261107f600184611a5f565b8151811061108f5761108f611a72565b6020910101516001600160f81b031916602960f81b14155b156110c557604051631f7481bb60e01b815260040160405180910390fd5b6110d0600183611a5f565b6110db826001611a88565b036110eb57506000949350505050565b600080806110fa846001611a88565b90505b611108600186611a5f565b8110156111c957600086828151811061112357611123611a72565b01602001516001600160f81b0319169050600560fb1b819003611152578361114a81611a9b565b9450506111c0565b6001600160f81b03198116602960f81b03611192578360000361118857604051631f7481bb60e01b815260040160405180910390fd5b8361114a81611ab4565b600b60fa1b6001600160f81b031982161480156111ad575083155b156111c057826111bc81611a9b565b9350505b506001016110fd565b5081156111e957604051631f7481bb60e01b815260040160405180910390fd5b6111f4816001611a88565b979650505050505050565b6001600160a01b0381166112285780604051638944034760e01b8152600401610be991906114b1565b6000611232610f91565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b6000600080516020611ae8833981519152610a34565b6112ae82611329565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a28051156112f3576106ca8282611385565b6105436113fb565b600060343610610562575060331936013590565b6000611319610b7e565b54600160401b900460ff16919050565b806001600160a01b03163b6000036113565780604051634c9c8ce360e01b8152600401610be991906114b1565b600080516020611ae883398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080846001600160a01b0316846040516113a29190611acb565b600060405180830381855af49150503d80600081146113dd576040519150601f19603f3d011682016040523d82523d6000602084013e6113e2565b606091505b50915091506113f285838361141a565b95945050505050565b3415610baf5760405163b398979f60e01b815260040160405180910390fd5b60608261142f5761142a82611470565b611469565b815115801561144657506001600160a01b0384163b155b156114665783604051639996b31560e01b8152600401610be991906114b1565b50805b9392505050565b80511561147f57805160208201fd5b60405163d6bda27560e01b815260040160405180910390fd5b6000602082840312156114aa57600080fd5b5035919050565b6001600160a01b0391909116815260200190565b80356001600160e01b0319811681146114dd57600080fd5b919050565b600080604083850312156114f557600080fd5b82359150611505602084016114c5565b90509250929050565b6001600160a01b0381168114610c9f57600080fd5b6000806040838503121561153657600080fd5b82356115418161150e565b915060208301356115518161150e565b809150509250929050565b634e487b7160e01b600052604160045260246000fd5b6000806040838503121561158557600080fd5b82356115908161150e565b915060208301356001600160401b03808211156115ac57600080fd5b818501915085601f8301126115c057600080fd5b8135818111156115d2576115d261155c565b604051601f8201601f19908116603f011681019083821181831017156115fa576115fa61155c565b8160405282815288602084870101111561161357600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b6000806040838503121561164857600080fd5b82356116538161150e565b9150611505602084016114c5565b80356001600160401b03811681146114dd57600080fd5b60008060006060848603121561168d57600080fd5b8335925061169d602085016114c5565b91506116ab60408501611661565b90509250925092565b60005b838110156116cf5781810151838201526020016116b7565b50506000910152565b600081518084526116f08160208601602086016116b4565b601f01601f19169290920160200192915050565b60208152815160208201526000602083015160a0604084015261172a60c08401826116d8565b905063ffffffff60e01b60408501511660608401526001600160401b0360608501511660808401526080840151151560a08401528091505092915050565b60208152600061146960208301846116d8565b6000806000806060858703121561179157600080fd5b8435935060208501356001600160401b03808211156117af57600080fd5b818701915087601f8301126117c357600080fd5b8135818111156117d257600080fd5b8860208285010111156117e457600080fd5b6020830195508094505050506117fc60408601611661565b905092959194509250565b60006020828403121561181957600080fd5b81356114698161150e565b60006020828403121561183657600080fd5b81516114698161150e565b9182526001600160e01b031916602082015260400190565b9283526001600160e01b03199190911660208301526001600160401b0316604082015260600190565b600181811c9082168061189657607f821691505b6020821081036118b657634e487b7160e01b600052602260045260246000fd5b50919050565b8183823760009101908152919050565b601f8211156106ca576000816000526020600020601f850160051c810160208610156118f55750805b601f850160051c820191505b8181101561191457828155600101611901565b505050505050565b81516001600160401b038111156119355761193561155c565b611949816119438454611882565b846118cc565b602080601f83116001811461197e57600084156119665750858301515b600019600386901b1c1916600185901b178555611914565b600085815260208120601f198616915b828110156119ad5788860151825594840194600190910190840161198e565b50858210156119cb5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6000602082840312156119ed57600080fd5b5051919050565b805180151581146114dd57600080fd5b600080600060608486031215611a1957600080fd5b611a22846119f4565b9250602084015163ffffffff81168114611a3b57600080fd5b91506116ab604085016119f4565b634e487b7160e01b600052601160045260246000fd5b818103818111156103de576103de611a49565b634e487b7160e01b600052603260045260246000fd5b808201808211156103de576103de611a49565b600060018201611aad57611aad611a49565b5060010190565b600081611ac357611ac3611a49565b506000190190565b60008251611add8184602087016116b4565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbca26469706673582212208c2660670dda1db5f342bb6214f9593491f53b70c8f997e45e40a26c382fd37c64736f6c63430008180033",
}

// TemplateRegistryReplicaV1 is an auto generated Go binding around an Ethereum contract.
type TemplateRegistryReplicaV1 struct {
	abi abi.ABI
}

// NewTemplateRegistryReplicaV1 creates a new instance of TemplateRegistryReplicaV1.
func NewTemplateRegistryReplicaV1() *TemplateRegistryReplicaV1 {
	parsed, err := TemplateRegistryReplicaV1MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &TemplateRegistryReplicaV1{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *TemplateRegistryReplicaV1) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackUPGRADEINTERFACEVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) PackUPGRADEINTERFACEVERSION() []byte {
	enc, err := templateRegistryReplicaV1.abi.Pack("UPGRADE_INTERFACE_VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUPGRADEINTERFACEVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackUPGRADEINTERFACEVERSION(data []byte) (string, error) {
	out, err := templateRegistryReplicaV1.abi.Unpack("UPGRADE_INTERFACE_VERSION", data)
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
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) PackAuthority() []byte {
	enc, err := templateRegistryReplicaV1.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := templateRegistryReplicaV1.abi.Unpack("authority", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackCheck is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7db5f731.
//
// Solidity: function check(address target, bytes4 selector) view returns(bool)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) PackCheck(target common.Address, selector [4]byte) []byte {
	enc, err := templateRegistryReplicaV1.abi.Pack("check", target, selector)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackCheck is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7db5f731.
//
// Solidity: function check(address target, bytes4 selector) view returns(bool)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackCheck(data []byte) (bool, error) {
	out, err := templateRegistryReplicaV1.abi.Unpack("check", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackCheckWithParamCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfee8c9ac.
//
// Solidity: function checkWithParamCount(address target, bytes4 selector) view returns(bool approved, uint256 paramCount)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) PackCheckWithParamCount(target common.Address, selector [4]byte) []byte {
	enc, err := templateRegistryReplicaV1.abi.Pack("checkWithParamCount", target, selector)
	if err != nil {
		panic(err)
	}
	return enc
}

// CheckWithParamCountOutput serves as a container for the return parameters of contract
// method CheckWithParamCount.
type CheckWithParamCountOutput struct {
	Approved   bool
	ParamCount *big.Int
}

// UnpackCheckWithParamCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfee8c9ac.
//
// Solidity: function checkWithParamCount(address target, bytes4 selector) view returns(bool approved, uint256 paramCount)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackCheckWithParamCount(data []byte) (CheckWithParamCountOutput, error) {
	out, err := templateRegistryReplicaV1.abi.Unpack("checkWithParamCount", data)
	outstruct := new(CheckWithParamCountOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Approved = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ParamCount = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, err

}

// PackContractVersion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa0a8e460.
//
// Solidity: function contractVersion() pure returns(uint256)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) PackContractVersion() []byte {
	enc, err := templateRegistryReplicaV1.abi.Pack("contractVersion")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackContractVersion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa0a8e460.
//
// Solidity: function contractVersion() pure returns(uint256)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackContractVersion(data []byte) (*big.Int, error) {
	out, err := templateRegistryReplicaV1.abi.Unpack("contractVersion", data)
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
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) PackEndpoint() []byte {
	enc, err := templateRegistryReplicaV1.abi.Pack("endpoint")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackEndpoint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackEndpoint(data []byte) (common.Address, error) {
	out, err := templateRegistryReplicaV1.abi.Unpack("endpoint", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetAddressByResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) PackGetAddressByResourceId(resourceId [32]byte) []byte {
	enc, err := templateRegistryReplicaV1.abi.Pack("getAddressByResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAddressByResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackGetAddressByResourceId(data []byte) (common.Address, error) {
	out, err := templateRegistryReplicaV1.abi.Unpack("getAddressByResourceId", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetKey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3dd93802.
//
// Solidity: function getKey(bytes32 bytecodeHash, bytes4 selector) pure returns(bytes32)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) PackGetKey(bytecodeHash [32]byte, selector [4]byte) []byte {
	enc, err := templateRegistryReplicaV1.abi.Pack("getKey", bytecodeHash, selector)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetKey is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3dd93802.
//
// Solidity: function getKey(bytes32 bytecodeHash, bytes4 selector) pure returns(bytes32)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackGetKey(data []byte) ([32]byte, error) {
	out, err := templateRegistryReplicaV1.abi.Unpack("getKey", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackGetLastUpdatedAt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x397349df.
//
// Solidity: function getLastUpdatedAt(bytes32 key) view returns(uint64)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) PackGetLastUpdatedAt(key [32]byte) []byte {
	enc, err := templateRegistryReplicaV1.abi.Pack("getLastUpdatedAt", key)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetLastUpdatedAt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x397349df.
//
// Solidity: function getLastUpdatedAt(bytes32 key) view returns(uint64)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackGetLastUpdatedAt(data []byte) (uint64, error) {
	out, err := templateRegistryReplicaV1.abi.Unpack("getLastUpdatedAt", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, err
}

// PackGetTemplate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9757739b.
//
// Solidity: function getTemplate(bytes32 key) view returns((bytes32,string,bytes4,uint64,bool))
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) PackGetTemplate(key [32]byte) []byte {
	enc, err := templateRegistryReplicaV1.abi.Pack("getTemplate", key)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTemplate is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9757739b.
//
// Solidity: function getTemplate(bytes32 key) view returns((bytes32,string,bytes4,uint64,bool))
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackGetTemplate(data []byte) (ITemplateRegistryTemplate, error) {
	out, err := templateRegistryReplicaV1.abi.Unpack("getTemplate", data)
	if err != nil {
		return *new(ITemplateRegistryTemplate), err
	}
	out0 := *abi.ConvertType(out[0], new(ITemplateRegistryTemplate)).(*ITemplateRegistryTemplate)
	return out0, err
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x485cc955.
//
// Solidity: function initialize(address _endpoint, address authority_) returns()
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) PackInitialize(endpoint common.Address, authority common.Address) []byte {
	enc, err := templateRegistryReplicaV1.abi.Pack("initialize", endpoint, authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackInitialize0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc4d66de8.
//
// Solidity: function initialize(address _endpoint) returns()
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) PackInitialize0(endpoint common.Address) []byte {
	enc, err := templateRegistryReplicaV1.abi.Pack("initialize0", endpoint)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackOnTemplateApproved is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaf368e84.
//
// Solidity: function onTemplateApproved(bytes32 bytecodeHash, string signature, uint64 approvedAt) returns()
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) PackOnTemplateApproved(bytecodeHash [32]byte, signature string, approvedAt uint64) []byte {
	enc, err := templateRegistryReplicaV1.abi.Pack("onTemplateApproved", bytecodeHash, signature, approvedAt)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackOnTemplateRevoked is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x88ee8c01.
//
// Solidity: function onTemplateRevoked(bytes32 bytecodeHash, bytes4 selector, uint64 revokedAt) returns()
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) PackOnTemplateRevoked(bytecodeHash [32]byte, selector [4]byte, revokedAt uint64) []byte {
	enc, err := templateRegistryReplicaV1.abi.Pack("onTemplateRevoked", bytecodeHash, selector, revokedAt)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackProxiableUUID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) PackProxiableUUID() []byte {
	enc, err := templateRegistryReplicaV1.abi.Pack("proxiableUUID")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackProxiableUUID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackProxiableUUID(data []byte) ([32]byte, error) {
	out, err := templateRegistryReplicaV1.abi.Unpack("proxiableUUID", data)
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
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) PackResourceId() []byte {
	enc, err := templateRegistryReplicaV1.abi.Pack("resourceId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackResourceId(data []byte) ([32]byte, error) {
	out, err := templateRegistryReplicaV1.abi.Unpack("resourceId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackUpgradeToAndCall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) PackUpgradeToAndCall(newImplementation common.Address, data []byte) []byte {
	enc, err := templateRegistryReplicaV1.abi.Pack("upgradeToAndCall", newImplementation, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// TemplateRegistryReplicaV1AuthorityUpdated represents a AuthorityUpdated event raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1AuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const TemplateRegistryReplicaV1AuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (TemplateRegistryReplicaV1AuthorityUpdated) ContractEventName() string {
	return TemplateRegistryReplicaV1AuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackAuthorityUpdatedEvent(log *types.Log) (*TemplateRegistryReplicaV1AuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != templateRegistryReplicaV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TemplateRegistryReplicaV1AuthorityUpdated)
	if len(log.Data) > 0 {
		if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range templateRegistryReplicaV1.abi.Events[event].Inputs {
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

// TemplateRegistryReplicaV1Initialized represents a Initialized event raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1Initialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const TemplateRegistryReplicaV1InitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (TemplateRegistryReplicaV1Initialized) ContractEventName() string {
	return TemplateRegistryReplicaV1InitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackInitializedEvent(log *types.Log) (*TemplateRegistryReplicaV1Initialized, error) {
	event := "Initialized"
	if log.Topics[0] != templateRegistryReplicaV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TemplateRegistryReplicaV1Initialized)
	if len(log.Data) > 0 {
		if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range templateRegistryReplicaV1.abi.Events[event].Inputs {
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

// TemplateRegistryReplicaV1TemplateMirrored represents a TemplateMirrored event raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1TemplateMirrored struct {
	Key          [32]byte
	BytecodeHash [32]byte
	Selector     [4]byte
	ApprovedAt   uint64
	Raw          *types.Log // Blockchain specific contextual infos
}

const TemplateRegistryReplicaV1TemplateMirroredEventName = "TemplateMirrored"

// ContractEventName returns the user-defined event name.
func (TemplateRegistryReplicaV1TemplateMirrored) ContractEventName() string {
	return TemplateRegistryReplicaV1TemplateMirroredEventName
}

// UnpackTemplateMirroredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TemplateMirrored(bytes32 indexed key, bytes32 bytecodeHash, bytes4 selector, uint64 approvedAt)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackTemplateMirroredEvent(log *types.Log) (*TemplateRegistryReplicaV1TemplateMirrored, error) {
	event := "TemplateMirrored"
	if log.Topics[0] != templateRegistryReplicaV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TemplateRegistryReplicaV1TemplateMirrored)
	if len(log.Data) > 0 {
		if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range templateRegistryReplicaV1.abi.Events[event].Inputs {
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

// TemplateRegistryReplicaV1TemplateRevocationMirrored represents a TemplateRevocationMirrored event raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1TemplateRevocationMirrored struct {
	Key          [32]byte
	BytecodeHash [32]byte
	Selector     [4]byte
	RevokedAt    uint64
	Raw          *types.Log // Blockchain specific contextual infos
}

const TemplateRegistryReplicaV1TemplateRevocationMirroredEventName = "TemplateRevocationMirrored"

// ContractEventName returns the user-defined event name.
func (TemplateRegistryReplicaV1TemplateRevocationMirrored) ContractEventName() string {
	return TemplateRegistryReplicaV1TemplateRevocationMirroredEventName
}

// UnpackTemplateRevocationMirroredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TemplateRevocationMirrored(bytes32 indexed key, bytes32 bytecodeHash, bytes4 selector, uint64 revokedAt)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackTemplateRevocationMirroredEvent(log *types.Log) (*TemplateRegistryReplicaV1TemplateRevocationMirrored, error) {
	event := "TemplateRevocationMirrored"
	if log.Topics[0] != templateRegistryReplicaV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TemplateRegistryReplicaV1TemplateRevocationMirrored)
	if len(log.Data) > 0 {
		if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range templateRegistryReplicaV1.abi.Events[event].Inputs {
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

// TemplateRegistryReplicaV1Upgraded represents a Upgraded event raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1Upgraded struct {
	Implementation common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const TemplateRegistryReplicaV1UpgradedEventName = "Upgraded"

// ContractEventName returns the user-defined event name.
func (TemplateRegistryReplicaV1Upgraded) ContractEventName() string {
	return TemplateRegistryReplicaV1UpgradedEventName
}

// UnpackUpgradedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Upgraded(address indexed implementation)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackUpgradedEvent(log *types.Log) (*TemplateRegistryReplicaV1Upgraded, error) {
	event := "Upgraded"
	if log.Topics[0] != templateRegistryReplicaV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TemplateRegistryReplicaV1Upgraded)
	if len(log.Data) > 0 {
		if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range templateRegistryReplicaV1.abi.Events[event].Inputs {
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
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], templateRegistryReplicaV1.abi.Errors["AddressEmptyCode"].ID.Bytes()[:4]) {
		return templateRegistryReplicaV1.UnpackAddressEmptyCodeError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryReplicaV1.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return templateRegistryReplicaV1.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryReplicaV1.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return templateRegistryReplicaV1.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryReplicaV1.abi.Errors["FailedCall"].ID.Bytes()[:4]) {
		return templateRegistryReplicaV1.UnpackFailedCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryReplicaV1.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return templateRegistryReplicaV1.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryReplicaV1.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return templateRegistryReplicaV1.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryReplicaV1.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return templateRegistryReplicaV1.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryReplicaV1.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return templateRegistryReplicaV1.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryReplicaV1.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return templateRegistryReplicaV1.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryReplicaV1.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return templateRegistryReplicaV1.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryReplicaV1.abi.Errors["SignatureParamsMalformed"].ID.Bytes()[:4]) {
		return templateRegistryReplicaV1.UnpackSignatureParamsMalformedError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryReplicaV1.abi.Errors["TemplateRegistryReplicaV1EmptyBytecodeHash"].ID.Bytes()[:4]) {
		return templateRegistryReplicaV1.UnpackTemplateRegistryReplicaV1EmptyBytecodeHashError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryReplicaV1.abi.Errors["TemplateRegistryReplicaV1EmptySignature"].ID.Bytes()[:4]) {
		return templateRegistryReplicaV1.UnpackTemplateRegistryReplicaV1EmptySignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryReplicaV1.abi.Errors["TemplateRegistryReplicaV1NotFromPrivateHub"].ID.Bytes()[:4]) {
		return templateRegistryReplicaV1.UnpackTemplateRegistryReplicaV1NotFromPrivateHubError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryReplicaV1.abi.Errors["UUPSUnauthorizedCallContext"].ID.Bytes()[:4]) {
		return templateRegistryReplicaV1.UnpackUUPSUnauthorizedCallContextError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryReplicaV1.abi.Errors["UUPSUnsupportedProxiableUUID"].ID.Bytes()[:4]) {
		return templateRegistryReplicaV1.UnpackUUPSUnsupportedProxiableUUIDError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// TemplateRegistryReplicaV1AddressEmptyCode represents a AddressEmptyCode error raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1AddressEmptyCode struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressEmptyCode(address target)
func TemplateRegistryReplicaV1AddressEmptyCodeErrorID() common.Hash {
	return common.HexToHash("0x9996b315c842ff135b8fc4a08ad5df1c344efbc03d2687aecc0678050d2aac89")
}

// UnpackAddressEmptyCodeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressEmptyCode(address target)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackAddressEmptyCodeError(raw []byte) (*TemplateRegistryReplicaV1AddressEmptyCode, error) {
	out := new(TemplateRegistryReplicaV1AddressEmptyCode)
	if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, "AddressEmptyCode", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryReplicaV1ERC1967InvalidImplementation represents a ERC1967InvalidImplementation error raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1ERC1967InvalidImplementation struct {
	Implementation common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func TemplateRegistryReplicaV1ERC1967InvalidImplementationErrorID() common.Hash {
	return common.HexToHash("0x4c9c8ce3ceb3130f17f7cdba48d89b5b0129f266a8bac114e6e315a41879b617")
}

// UnpackERC1967InvalidImplementationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackERC1967InvalidImplementationError(raw []byte) (*TemplateRegistryReplicaV1ERC1967InvalidImplementation, error) {
	out := new(TemplateRegistryReplicaV1ERC1967InvalidImplementation)
	if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, "ERC1967InvalidImplementation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryReplicaV1ERC1967NonPayable represents a ERC1967NonPayable error raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1ERC1967NonPayable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967NonPayable()
func TemplateRegistryReplicaV1ERC1967NonPayableErrorID() common.Hash {
	return common.HexToHash("0xb398979fa84f543c8e222f17890372c487baf85e062276c127fef521eea7224b")
}

// UnpackERC1967NonPayableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967NonPayable()
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackERC1967NonPayableError(raw []byte) (*TemplateRegistryReplicaV1ERC1967NonPayable, error) {
	out := new(TemplateRegistryReplicaV1ERC1967NonPayable)
	if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, "ERC1967NonPayable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryReplicaV1FailedCall represents a FailedCall error raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1FailedCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedCall()
func TemplateRegistryReplicaV1FailedCallErrorID() common.Hash {
	return common.HexToHash("0xd6bda27508c0fb6d8a39b4b122878dab26f731a7d4e4abe711dd3731899052a4")
}

// UnpackFailedCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedCall()
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackFailedCallError(raw []byte) (*TemplateRegistryReplicaV1FailedCall, error) {
	out := new(TemplateRegistryReplicaV1FailedCall)
	if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, "FailedCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryReplicaV1InvalidInitialization represents a InvalidInitialization error raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1InvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func TemplateRegistryReplicaV1InvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackInvalidInitializationError(raw []byte) (*TemplateRegistryReplicaV1InvalidInitialization, error) {
	out := new(TemplateRegistryReplicaV1InvalidInitialization)
	if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryReplicaV1NotInitializing represents a NotInitializing error raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1NotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func TemplateRegistryReplicaV1NotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackNotInitializingError(raw []byte) (*TemplateRegistryReplicaV1NotInitializing, error) {
	out := new(TemplateRegistryReplicaV1NotInitializing)
	if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryReplicaV1RaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1RaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func TemplateRegistryReplicaV1RaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*TemplateRegistryReplicaV1RaylsAccessManagedContractPaused, error) {
	out := new(TemplateRegistryReplicaV1RaylsAccessManagedContractPaused)
	if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryReplicaV1RaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1RaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func TemplateRegistryReplicaV1RaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*TemplateRegistryReplicaV1RaylsAccessManagedInvalidAuthority, error) {
	out := new(TemplateRegistryReplicaV1RaylsAccessManagedInvalidAuthority)
	if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryReplicaV1RaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1RaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func TemplateRegistryReplicaV1RaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*TemplateRegistryReplicaV1RaylsAccessManagedMustSchedule, error) {
	out := new(TemplateRegistryReplicaV1RaylsAccessManagedMustSchedule)
	if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryReplicaV1RaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1RaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func TemplateRegistryReplicaV1RaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*TemplateRegistryReplicaV1RaylsAccessManagedUnauthorized, error) {
	out := new(TemplateRegistryReplicaV1RaylsAccessManagedUnauthorized)
	if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryReplicaV1SignatureParamsMalformed represents a SignatureParams__Malformed error raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1SignatureParamsMalformed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SignatureParams__Malformed()
func TemplateRegistryReplicaV1SignatureParamsMalformedErrorID() common.Hash {
	return common.HexToHash("0x1f7481bb5d9911fb47f48bfcb3b8825fe62e8920ad080c238653c61ad4152190")
}

// UnpackSignatureParamsMalformedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SignatureParams__Malformed()
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackSignatureParamsMalformedError(raw []byte) (*TemplateRegistryReplicaV1SignatureParamsMalformed, error) {
	out := new(TemplateRegistryReplicaV1SignatureParamsMalformed)
	if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, "SignatureParamsMalformed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryReplicaV1TemplateRegistryReplicaV1EmptyBytecodeHash represents a TemplateRegistryReplicaV1__EmptyBytecodeHash error raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1TemplateRegistryReplicaV1EmptyBytecodeHash struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TemplateRegistryReplicaV1__EmptyBytecodeHash()
func TemplateRegistryReplicaV1TemplateRegistryReplicaV1EmptyBytecodeHashErrorID() common.Hash {
	return common.HexToHash("0x4230fe6bab016e49aaf65c39b6c3834c5d3f7986fa3750d64052ad611712f3f6")
}

// UnpackTemplateRegistryReplicaV1EmptyBytecodeHashError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TemplateRegistryReplicaV1__EmptyBytecodeHash()
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackTemplateRegistryReplicaV1EmptyBytecodeHashError(raw []byte) (*TemplateRegistryReplicaV1TemplateRegistryReplicaV1EmptyBytecodeHash, error) {
	out := new(TemplateRegistryReplicaV1TemplateRegistryReplicaV1EmptyBytecodeHash)
	if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, "TemplateRegistryReplicaV1EmptyBytecodeHash", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryReplicaV1TemplateRegistryReplicaV1EmptySignature represents a TemplateRegistryReplicaV1__EmptySignature error raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1TemplateRegistryReplicaV1EmptySignature struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TemplateRegistryReplicaV1__EmptySignature()
func TemplateRegistryReplicaV1TemplateRegistryReplicaV1EmptySignatureErrorID() common.Hash {
	return common.HexToHash("0x4da6cb2f391d8d778b7a52e1c11ce6bed5a4969e4e56e0e33cc630248b463391")
}

// UnpackTemplateRegistryReplicaV1EmptySignatureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TemplateRegistryReplicaV1__EmptySignature()
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackTemplateRegistryReplicaV1EmptySignatureError(raw []byte) (*TemplateRegistryReplicaV1TemplateRegistryReplicaV1EmptySignature, error) {
	out := new(TemplateRegistryReplicaV1TemplateRegistryReplicaV1EmptySignature)
	if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, "TemplateRegistryReplicaV1EmptySignature", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryReplicaV1TemplateRegistryReplicaV1NotFromPrivateHub represents a TemplateRegistryReplicaV1__NotFromPrivateHub error raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1TemplateRegistryReplicaV1NotFromPrivateHub struct {
	FromChainId *big.Int
	HubId       *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TemplateRegistryReplicaV1__NotFromPrivateHub(uint256 fromChainId, uint256 hubId)
func TemplateRegistryReplicaV1TemplateRegistryReplicaV1NotFromPrivateHubErrorID() common.Hash {
	return common.HexToHash("0x90c1945c6237105309778cec839fe5382832f770c46e3c63a06caa501896dfa5")
}

// UnpackTemplateRegistryReplicaV1NotFromPrivateHubError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TemplateRegistryReplicaV1__NotFromPrivateHub(uint256 fromChainId, uint256 hubId)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackTemplateRegistryReplicaV1NotFromPrivateHubError(raw []byte) (*TemplateRegistryReplicaV1TemplateRegistryReplicaV1NotFromPrivateHub, error) {
	out := new(TemplateRegistryReplicaV1TemplateRegistryReplicaV1NotFromPrivateHub)
	if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, "TemplateRegistryReplicaV1NotFromPrivateHub", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryReplicaV1UUPSUnauthorizedCallContext represents a UUPSUnauthorizedCallContext error raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1UUPSUnauthorizedCallContext struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnauthorizedCallContext()
func TemplateRegistryReplicaV1UUPSUnauthorizedCallContextErrorID() common.Hash {
	return common.HexToHash("0xe07c8dba242a06571ac65fe4bbe20522c9fb111cb33599b799ff8039c1ed18f4")
}

// UnpackUUPSUnauthorizedCallContextError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnauthorizedCallContext()
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackUUPSUnauthorizedCallContextError(raw []byte) (*TemplateRegistryReplicaV1UUPSUnauthorizedCallContext, error) {
	out := new(TemplateRegistryReplicaV1UUPSUnauthorizedCallContext)
	if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, "UUPSUnauthorizedCallContext", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryReplicaV1UUPSUnsupportedProxiableUUID represents a UUPSUnsupportedProxiableUUID error raised by the TemplateRegistryReplicaV1 contract.
type TemplateRegistryReplicaV1UUPSUnsupportedProxiableUUID struct {
	Slot [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func TemplateRegistryReplicaV1UUPSUnsupportedProxiableUUIDErrorID() common.Hash {
	return common.HexToHash("0xaa1d49a4c084bfa9aeeee2a0be65267a7f19ba7e1476b114dac513d2c14cb563")
}

// UnpackUUPSUnsupportedProxiableUUIDError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func (templateRegistryReplicaV1 *TemplateRegistryReplicaV1) UnpackUUPSUnsupportedProxiableUUIDError(raw []byte) (*TemplateRegistryReplicaV1UUPSUnsupportedProxiableUUID, error) {
	out := new(TemplateRegistryReplicaV1UUPSUnsupportedProxiableUUID)
	if err := templateRegistryReplicaV1.abi.UnpackIntoInterface(out, "UUPSUnsupportedProxiableUUID", raw); err != nil {
		return nil, err
	}
	return out, nil
}
