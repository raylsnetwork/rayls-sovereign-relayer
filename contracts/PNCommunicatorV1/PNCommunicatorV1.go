// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package PNCommunicatorV1

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

// SharedObjectsCommunicatiorDataHistory is an auto generated low-level Go binding around an user-defined struct.
type SharedObjectsCommunicatiorDataHistory struct {
	Status      *big.Int
	BlockNumber *big.Int
	Message     string
}

// PNCommunicatorV1MetaData contains all meta data concerning the PNCommunicatorV1 contract.
var PNCommunicatorV1MetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addSharedInfo\",\"inputs\":[{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_status\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_context\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_message\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"contractVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"endpoint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRaylsEndpoint\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAddressByResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllSharedInfo\",\"inputs\":[{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"context\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structSharedObjects.CommunicatiorDataHistory[]\",\"components\":[{\"name\":\"status\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSharedInfo\",\"inputs\":[{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"status\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"context\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSharedInfoAt\",\"inputs\":[{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"status\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"context\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSharedInfoCount\",\"inputs\":[{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasSharedInfo\",\"inputs\":[{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_endpointAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"authority_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeSharedInfo\",\"inputs\":[{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeSharedInfoAt\",\"inputs\":[{\"name\":\"_sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resourceId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SharedInfoAdded\",\"inputs\":[{\"name\":\"sharedId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"status\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SharedInfoRemoved\",\"inputs\":[{\"name\":\"sharedId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	ID:  "PNCommunicatorV1",
	Bin: "0x60a06040523060805234801561001457600080fd5b50608051611a5c61003e60003960008181610cba01528181610ce30152610e1a0152611a5c6000f3fe6080604052600436106100e95760003560e01c80635e280f11116100855780635e280f11146102585780635f997c5b14610278578063976e37441461028e578063a0a8e460146102bc578063ad3cb1cc146102d0578063bf7e214f1461030e578063c4d66de814610323578063cba06a4814610343578063db76edaa1461036357600080fd5b80630f4ea0c5146100ee57806311f50c85146101275780632262357b14610154578063485cc955146101745780634888c7bb146101965780634f1ef286146101d557806352d1902d146101e857806358f48ebf1461020b5780635c7a852014610238575b600080fd5b3480156100fa57600080fd5b5061010e610109366004611344565b610383565b60405161011e94939291906113b6565b60405180910390f35b34801561013357600080fd5b506101476101423660046113e6565b6104d0565b60405161011e91906113ff565b34801561016057600080fd5b5061010e61016f3660046113e6565b610544565b34801561018057600080fd5b5061019461018f366004611428565b6106b9565b005b3480156101a257600080fd5b506101c56101b13660046113e6565b600090815260026020526040902054151590565b604051901515815260200161011e565b6101946101e33660046114ec565b6107ca565b3480156101f457600080fd5b506101fd6107e9565b60405190815260200161011e565b34801561021757600080fd5b506101fd6102263660046113e6565b60009081526002602052604090205490565b34801561024457600080fd5b506101946102533660046113e6565b610806565b34801561026457600080fd5b50600054610147906001600160a01b031681565b34801561028457600080fd5b506101fd60015481565b34801561029a57600080fd5b506102ae6102a93660046113e6565b610842565b60405161011e92919061154f565b3480156102c857600080fd5b5060016101fd565b3480156102dc57600080fd5b50610301604051806040016040528060058152602001640352e302e360dc1b81525081565b60405161011e91906115d7565b34801561031a57600080fd5b5061014761095e565b34801561032f57600080fd5b5061019461033e3660046115ea565b610977565b34801561034f57600080fd5b5061019461035e366004611344565b6109a1565b34801561036f57600080fd5b5061019461037e366004611607565b610b2c565b6000828152600260205260408120548190606090829085106103c05760405162461bcd60e51b81526004016103b790611674565b60405180910390fd5b60008681526002602052604081208054879081106103e0576103e06116a1565b9060005260206000209060030201604051806060016040529081600082015481526020016001820154815260200160028201805461041d906116b7565b80601f0160208091040260200160405190810160405280929190818152602001828054610449906116b7565b80156104965780601f1061046b57610100808354040283529160200191610496565b820191906000526020600020905b81548152906001019060200180831161047957829003601f168201915b505050919092525050815160208084015160409485015160009c8d52600290925293909a2060010154909a92999850965090945050505050565b600080546040516311f50c8560e01b8152600481018490526001600160a01b03909116906311f50c8590602401602060405180830381865afa15801561051a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061053e91906116f1565b92915050565b6000818152600260205260408120548190606090829061059e5760405162461bcd60e51b81526020600482015260156024820152744e6f2073686172656420696e666f2065786973747360581b60448201526064016103b7565b600085815260026020526040812080546105ba9060019061170e565b815481106105ca576105ca6116a1565b90600052602060002090600302016040518060600160405290816000820154815260200160018201548152602001600282018054610607906116b7565b80601f0160208091040260200160405190810160405280929190818152602001828054610633906116b7565b80156106805780601f1061065557610100808354040283529160200191610680565b820191906000526020600020905b81548152906001019060200180831161066357829003601f168201915b505050919092525050815160208084015160409485015160009b8c52600290925293909920600101549099929897509550909350505050565b60006106c3610c3b565b805490915060ff600160401b82041615906001600160401b03166000811580156106ea5750825b90506000826001600160401b031660011480156107065750303b155b905081158015610714575080155b156107325760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561075c57845460ff60401b1916600160401b1785555b610764610c64565b61076d87610977565b600560015561077b86610c6e565b83156107c157845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50505050505050565b6107d2610caf565b6107db82610d3d565b6107e58282610d56565b5050565b60006107f3610e0f565b50600080516020611a0783398151915290565b61081c336000356001600160e01b031916610e58565b60008181526002602052604081209061083582826112a6565b6001820160009055505050565b600081815260026020908152604080832060018101548154835181860281018601909452808452606094919391839190879084015b8282101561094f578382906000526020600020906003020160405180606001604052908160008201548152602001600182015481526020016002820180546108be906116b7565b80601f01602080910402602001604051908101604052809291908181526020018280546108ea906116b7565b80156109375780601f1061090c57610100808354040283529160200191610937565b820191906000526020600020905b81548152906001019060200180831161091a57829003601f168201915b50505050508152505081526020019060010190610877565b50505050905091509150915091565b6000610968610f9a565b546001600160a01b0316919050565b61097f610ffc565b600080546001600160a01b0319166001600160a01b0392909216919091179055565b6109b7336000356001600160e01b031916610e58565b60008281526002602052604090205481106109e45760405162461bcd60e51b81526004016103b790611674565b6000828152600260205260409020546109ff9060019061170e565b811015610aa15760008281526002602052604090208054610a229060019061170e565b81548110610a3257610a326116a1565b9060005260206000209060030201600260008481526020019081526020016000206000018281548110610a6757610a676116a1565b9060005260206000209060030201600082015481600001556001820154816001015560028201816002019081610a9d9190611794565b5050505b6000828152600260205260409020805480610abe57610abe61186e565b6000828152602081206003600019909301928302018181556001810182905590610aeb60028301826112c7565b5050905560408051838152602081018390527fcdf34371fb3ddbefe693dfd0f972e9beeb51f181fbc68ac5094ce3ae89e80402910160405180910390a15050565b610b42336000356001600160e01b031916610e58565b60408051606081018252848152436020808301919091528183018490526000878152600290915291822054909103610b895760008581526002602052604090206001018390555b60008581526002602081815260408084208054600181810183559186529483902086516003909602019485559185015191840191909155830151839291820190610bd39082611884565b5050506000858152600260205260409020547fe00365e63912d463887c3d0161120235b0a762e8d01193f8691ed6a9b66700e4908690869043908690610c1b9060019061170e565b604051610c2c959493929190611937565b60405180910390a15050505050565b6000807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0061053e565b610c6c610ffc565b565b6000610c78610f9a565b80549091506001600160a01b031615610ca65781604051638944034760e01b81526004016103b791906113ff565b6107e582611021565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161480610d1f57507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316610d136110b1565b6001600160a01b031614155b15610c6c5760405163703e46dd60e11b815260040160405180910390fd5b610d53336000356001600160e01b031916610e58565b50565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610db0575060408051601f3d908101601f19168201909252610dad9181019061196e565b60015b610dcf5781604051634c9c8ce360e01b81526004016103b791906113ff565b600080516020611a078339815191528114610e0057604051632a87526960e21b8152600481018290526024016103b7565b610e0a83836110c7565b505050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614610c6c5760405163703e46dd60e11b815260040160405180910390fd5b6000610e62610f9a565b80549091506001600160a01b031680610e91576000604051638944034760e01b81526004016103b791906113ff565b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa158015610ef5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f19919061199c565b925092509250826107c1578015610f435760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff821615610f7f5760405163a426878960e01b81526001600160a01b038816600482015263ffffffff831660248201526044016103b7565b86604051632ecd3d0360e21b81526004016103b791906113ff565b60008060ff19610fcb60017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f3561170e565b604051602001610fdd91815260200190565b60408051601f1981840301815291905280516020909101201692915050565b61100461111d565b610c6c57604051631afcd79f60e31b815260040160405180910390fd5b6001600160a01b03811661104a5780604051638944034760e01b81526004016103b791906113ff565b6000611054610f9a565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b6000600080516020611a07833981519152610968565b6110d082611137565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a280511561111557610e0a8282611193565b6107e5611209565b6000611127610c3b565b54600160401b900460ff16919050565b806001600160a01b03163b6000036111645780604051634c9c8ce360e01b81526004016103b791906113ff565b600080516020611a0783398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080846001600160a01b0316846040516111b091906119ea565b600060405180830381855af49150503d80600081146111eb576040519150601f19603f3d011682016040523d82523d6000602084013e6111f0565b606091505b5091509150611200858383611228565b95945050505050565b3415610c6c5760405163b398979f60e01b815260040160405180910390fd5b60608261123d576112388261127e565b611277565b815115801561125457506001600160a01b0384163b155b156112745783604051639996b31560e01b81526004016103b791906113ff565b50805b9392505050565b80511561128d57805160208201fd5b60405163d6bda27560e01b815260040160405180910390fd5b5080546000825560030290600052602060002090810190610d539190611301565b5080546112d3906116b7565b6000825580601f106112e3575050565b601f016020900490600052602060002090810190610d53919061132f565b8082111561132b5760008082556001820181905561132260028301826112c7565b50600301611301565b5090565b5b8082111561132b5760008155600101611330565b6000806040838503121561135757600080fd5b50508035926020909101359150565b60005b83811015611381578181015183820152602001611369565b50506000910152565b600081518084526113a2816020860160208601611366565b601f01601f19169290920160200192915050565b8481528360208201526080604082015260006113d5608083018561138a565b905082606083015295945050505050565b6000602082840312156113f857600080fd5b5035919050565b6001600160a01b0391909116815260200190565b6001600160a01b0381168114610d5357600080fd5b6000806040838503121561143b57600080fd5b823561144681611413565b9150602083013561145681611413565b809150509250929050565b634e487b7160e01b600052604160045260246000fd5b60006001600160401b038084111561149157611491611461565b604051601f8501601f19908116603f011681019082821181831017156114b9576114b9611461565b816040528093508581528686860111156114d257600080fd5b858560208301376000602087830101525050509392505050565b600080604083850312156114ff57600080fd5b823561150a81611413565b915060208301356001600160401b0381111561152557600080fd5b8301601f8101851361153657600080fd5b61154585823560208401611477565b9150509250929050565b6000604080830185845260206040818601528186518084526060935060608701915060608160051b88010183890160005b838110156115c757898303605f1901855281518051845286810151878501528801518884018890526115b48885018261138a565b9587019593505090850190600101611580565b50909a9950505050505050505050565b602081526000611277602083018461138a565b6000602082840312156115fc57600080fd5b813561127781611413565b6000806000806080858703121561161d57600080fd5b84359350602085013592506040850135915060608501356001600160401b0381111561164857600080fd5b8501601f8101871361165957600080fd5b61166887823560208401611477565b91505092959194509250565b602080825260139082015272496e646578206f7574206f6620626f756e647360681b604082015260600190565b634e487b7160e01b600052603260045260246000fd5b600181811c908216806116cb57607f821691505b6020821081036116eb57634e487b7160e01b600052602260045260246000fd5b50919050565b60006020828403121561170357600080fd5b815161127781611413565b8181038181111561053e57634e487b7160e01b600052601160045260246000fd5b601f821115610e0a576000816000526020600020601f850160051c810160208610156117585750805b601f850160051c820191505b8181101561177757828155600101611764565b505050505050565b600019600383901b1c191660019190911b1790565b81810361179f575050565b6117a982546116b7565b6001600160401b038111156117c0576117c0611461565b6117d4816117ce84546116b7565b8461172f565b6000601f82116001811461180257600083156117f05750848201545b6117fa848261177f565b855550611867565b600085815260209020601f19841690600086815260209020845b8381101561183c578286015482556001958601959091019060200161181c565b508583101561185a5781850154600019600388901b60f8161c191681555b50505060018360011b0184555b5050505050565b634e487b7160e01b600052603160045260246000fd5b81516001600160401b0381111561189d5761189d611461565b6118ab816117ce84546116b7565b602080601f8311600181146118da57600084156118c85750858301515b6118d2858261177f565b865550611777565b600085815260208120601f198616915b82811015611909578886015182559484019460019091019084016118ea565b50858210156119275787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b85815284602082015283604082015260a06060820152600061195c60a083018561138a565b90508260808301529695505050505050565b60006020828403121561198057600080fd5b5051919050565b8051801515811461199757600080fd5b919050565b6000806000606084860312156119b157600080fd5b6119ba84611987565b9250602084015163ffffffff811681146119d357600080fd5b91506119e160408501611987565b90509250925092565b600082516119fc818460208701611366565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbca26469706673582212205c1d1431e5592776daded29dfedbed21c35e5f1bb0e32adf9e248ed088e9a86364736f6c63430008180033",
}

// PNCommunicatorV1 is an auto generated Go binding around an Ethereum contract.
type PNCommunicatorV1 struct {
	abi abi.ABI
}

// NewPNCommunicatorV1 creates a new instance of PNCommunicatorV1.
func NewPNCommunicatorV1() *PNCommunicatorV1 {
	parsed, err := PNCommunicatorV1MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &PNCommunicatorV1{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *PNCommunicatorV1) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackUPGRADEINTERFACEVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (pNCommunicatorV1 *PNCommunicatorV1) PackUPGRADEINTERFACEVERSION() []byte {
	enc, err := pNCommunicatorV1.abi.Pack("UPGRADE_INTERFACE_VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUPGRADEINTERFACEVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackUPGRADEINTERFACEVERSION(data []byte) (string, error) {
	out, err := pNCommunicatorV1.abi.Unpack("UPGRADE_INTERFACE_VERSION", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackAddSharedInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdb76edaa.
//
// Solidity: function addSharedInfo(bytes32 _sharedId, uint256 _status, uint256 _context, string _message) returns()
func (pNCommunicatorV1 *PNCommunicatorV1) PackAddSharedInfo(sharedId [32]byte, status *big.Int, context *big.Int, message string) []byte {
	enc, err := pNCommunicatorV1.abi.Pack("addSharedInfo", sharedId, status, context, message)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (pNCommunicatorV1 *PNCommunicatorV1) PackAuthority() []byte {
	enc, err := pNCommunicatorV1.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := pNCommunicatorV1.abi.Unpack("authority", data)
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
func (pNCommunicatorV1 *PNCommunicatorV1) PackContractVersion() []byte {
	enc, err := pNCommunicatorV1.abi.Pack("contractVersion")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackContractVersion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa0a8e460.
//
// Solidity: function contractVersion() pure returns(uint256)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackContractVersion(data []byte) (*big.Int, error) {
	out, err := pNCommunicatorV1.abi.Unpack("contractVersion", data)
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
func (pNCommunicatorV1 *PNCommunicatorV1) PackEndpoint() []byte {
	enc, err := pNCommunicatorV1.abi.Pack("endpoint")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackEndpoint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackEndpoint(data []byte) (common.Address, error) {
	out, err := pNCommunicatorV1.abi.Unpack("endpoint", data)
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
func (pNCommunicatorV1 *PNCommunicatorV1) PackGetAddressByResourceId(resourceId [32]byte) []byte {
	enc, err := pNCommunicatorV1.abi.Pack("getAddressByResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAddressByResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackGetAddressByResourceId(data []byte) (common.Address, error) {
	out, err := pNCommunicatorV1.abi.Unpack("getAddressByResourceId", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetAllSharedInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x976e3744.
//
// Solidity: function getAllSharedInfo(bytes32 _sharedId) view returns(uint256 context, (uint256,uint256,string)[])
func (pNCommunicatorV1 *PNCommunicatorV1) PackGetAllSharedInfo(sharedId [32]byte) []byte {
	enc, err := pNCommunicatorV1.abi.Pack("getAllSharedInfo", sharedId)
	if err != nil {
		panic(err)
	}
	return enc
}

// GetAllSharedInfoOutput serves as a container for the return parameters of contract
// method GetAllSharedInfo.
type GetAllSharedInfoOutput struct {
	Context *big.Int
	Arg1    []SharedObjectsCommunicatiorDataHistory
}

// UnpackGetAllSharedInfo is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x976e3744.
//
// Solidity: function getAllSharedInfo(bytes32 _sharedId) view returns(uint256 context, (uint256,uint256,string)[])
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackGetAllSharedInfo(data []byte) (GetAllSharedInfoOutput, error) {
	out, err := pNCommunicatorV1.abi.Unpack("getAllSharedInfo", data)
	outstruct := new(GetAllSharedInfoOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Context = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Arg1 = *abi.ConvertType(out[1], new([]SharedObjectsCommunicatiorDataHistory)).(*[]SharedObjectsCommunicatiorDataHistory)
	return *outstruct, err

}

// PackGetSharedInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2262357b.
//
// Solidity: function getSharedInfo(bytes32 _sharedId) view returns(uint256 status, uint256 blockNumber, string message, uint256 context)
func (pNCommunicatorV1 *PNCommunicatorV1) PackGetSharedInfo(sharedId [32]byte) []byte {
	enc, err := pNCommunicatorV1.abi.Pack("getSharedInfo", sharedId)
	if err != nil {
		panic(err)
	}
	return enc
}

// GetSharedInfoOutput serves as a container for the return parameters of contract
// method GetSharedInfo.
type GetSharedInfoOutput struct {
	Status      *big.Int
	BlockNumber *big.Int
	Message     string
	Context     *big.Int
}

// UnpackGetSharedInfo is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2262357b.
//
// Solidity: function getSharedInfo(bytes32 _sharedId) view returns(uint256 status, uint256 blockNumber, string message, uint256 context)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackGetSharedInfo(data []byte) (GetSharedInfoOutput, error) {
	out, err := pNCommunicatorV1.abi.Unpack("getSharedInfo", data)
	outstruct := new(GetSharedInfoOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Status = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.BlockNumber = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.Message = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.Context = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	return *outstruct, err

}

// PackGetSharedInfoAt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0f4ea0c5.
//
// Solidity: function getSharedInfoAt(bytes32 _sharedId, uint256 _index) view returns(uint256 status, uint256 blockNumber, string message, uint256 context)
func (pNCommunicatorV1 *PNCommunicatorV1) PackGetSharedInfoAt(sharedId [32]byte, index *big.Int) []byte {
	enc, err := pNCommunicatorV1.abi.Pack("getSharedInfoAt", sharedId, index)
	if err != nil {
		panic(err)
	}
	return enc
}

// GetSharedInfoAtOutput serves as a container for the return parameters of contract
// method GetSharedInfoAt.
type GetSharedInfoAtOutput struct {
	Status      *big.Int
	BlockNumber *big.Int
	Message     string
	Context     *big.Int
}

// UnpackGetSharedInfoAt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0f4ea0c5.
//
// Solidity: function getSharedInfoAt(bytes32 _sharedId, uint256 _index) view returns(uint256 status, uint256 blockNumber, string message, uint256 context)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackGetSharedInfoAt(data []byte) (GetSharedInfoAtOutput, error) {
	out, err := pNCommunicatorV1.abi.Unpack("getSharedInfoAt", data)
	outstruct := new(GetSharedInfoAtOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Status = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.BlockNumber = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.Message = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.Context = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	return *outstruct, err

}

// PackGetSharedInfoCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x58f48ebf.
//
// Solidity: function getSharedInfoCount(bytes32 _sharedId) view returns(uint256)
func (pNCommunicatorV1 *PNCommunicatorV1) PackGetSharedInfoCount(sharedId [32]byte) []byte {
	enc, err := pNCommunicatorV1.abi.Pack("getSharedInfoCount", sharedId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetSharedInfoCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x58f48ebf.
//
// Solidity: function getSharedInfoCount(bytes32 _sharedId) view returns(uint256)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackGetSharedInfoCount(data []byte) (*big.Int, error) {
	out, err := pNCommunicatorV1.abi.Unpack("getSharedInfoCount", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackHasSharedInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4888c7bb.
//
// Solidity: function hasSharedInfo(bytes32 _sharedId) view returns(bool)
func (pNCommunicatorV1 *PNCommunicatorV1) PackHasSharedInfo(sharedId [32]byte) []byte {
	enc, err := pNCommunicatorV1.abi.Pack("hasSharedInfo", sharedId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackHasSharedInfo is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4888c7bb.
//
// Solidity: function hasSharedInfo(bytes32 _sharedId) view returns(bool)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackHasSharedInfo(data []byte) (bool, error) {
	out, err := pNCommunicatorV1.abi.Unpack("hasSharedInfo", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x485cc955.
//
// Solidity: function initialize(address _endpointAddress, address authority_) returns()
func (pNCommunicatorV1 *PNCommunicatorV1) PackInitialize(endpointAddress common.Address, authority common.Address) []byte {
	enc, err := pNCommunicatorV1.abi.Pack("initialize", endpointAddress, authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackInitialize0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc4d66de8.
//
// Solidity: function initialize(address _endpoint) returns()
func (pNCommunicatorV1 *PNCommunicatorV1) PackInitialize0(endpoint common.Address) []byte {
	enc, err := pNCommunicatorV1.abi.Pack("initialize0", endpoint)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackProxiableUUID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (pNCommunicatorV1 *PNCommunicatorV1) PackProxiableUUID() []byte {
	enc, err := pNCommunicatorV1.abi.Pack("proxiableUUID")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackProxiableUUID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackProxiableUUID(data []byte) ([32]byte, error) {
	out, err := pNCommunicatorV1.abi.Unpack("proxiableUUID", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackRemoveSharedInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c7a8520.
//
// Solidity: function removeSharedInfo(bytes32 _sharedId) returns()
func (pNCommunicatorV1 *PNCommunicatorV1) PackRemoveSharedInfo(sharedId [32]byte) []byte {
	enc, err := pNCommunicatorV1.abi.Pack("removeSharedInfo", sharedId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRemoveSharedInfoAt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcba06a48.
//
// Solidity: function removeSharedInfoAt(bytes32 _sharedId, uint256 _index) returns()
func (pNCommunicatorV1 *PNCommunicatorV1) PackRemoveSharedInfoAt(sharedId [32]byte, index *big.Int) []byte {
	enc, err := pNCommunicatorV1.abi.Pack("removeSharedInfoAt", sharedId, index)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (pNCommunicatorV1 *PNCommunicatorV1) PackResourceId() []byte {
	enc, err := pNCommunicatorV1.abi.Pack("resourceId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackResourceId(data []byte) ([32]byte, error) {
	out, err := pNCommunicatorV1.abi.Unpack("resourceId", data)
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
func (pNCommunicatorV1 *PNCommunicatorV1) PackUpgradeToAndCall(newImplementation common.Address, data []byte) []byte {
	enc, err := pNCommunicatorV1.abi.Pack("upgradeToAndCall", newImplementation, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// PNCommunicatorV1AuthorityUpdated represents a AuthorityUpdated event raised by the PNCommunicatorV1 contract.
type PNCommunicatorV1AuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const PNCommunicatorV1AuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (PNCommunicatorV1AuthorityUpdated) ContractEventName() string {
	return PNCommunicatorV1AuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackAuthorityUpdatedEvent(log *types.Log) (*PNCommunicatorV1AuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != pNCommunicatorV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNCommunicatorV1AuthorityUpdated)
	if len(log.Data) > 0 {
		if err := pNCommunicatorV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNCommunicatorV1.abi.Events[event].Inputs {
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

// PNCommunicatorV1Initialized represents a Initialized event raised by the PNCommunicatorV1 contract.
type PNCommunicatorV1Initialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const PNCommunicatorV1InitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (PNCommunicatorV1Initialized) ContractEventName() string {
	return PNCommunicatorV1InitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackInitializedEvent(log *types.Log) (*PNCommunicatorV1Initialized, error) {
	event := "Initialized"
	if log.Topics[0] != pNCommunicatorV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNCommunicatorV1Initialized)
	if len(log.Data) > 0 {
		if err := pNCommunicatorV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNCommunicatorV1.abi.Events[event].Inputs {
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

// PNCommunicatorV1SharedInfoAdded represents a SharedInfoAdded event raised by the PNCommunicatorV1 contract.
type PNCommunicatorV1SharedInfoAdded struct {
	SharedId    [32]byte
	Status      *big.Int
	BlockNumber *big.Int
	Message     string
	Index       *big.Int
	Raw         *types.Log // Blockchain specific contextual infos
}

const PNCommunicatorV1SharedInfoAddedEventName = "SharedInfoAdded"

// ContractEventName returns the user-defined event name.
func (PNCommunicatorV1SharedInfoAdded) ContractEventName() string {
	return PNCommunicatorV1SharedInfoAddedEventName
}

// UnpackSharedInfoAddedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SharedInfoAdded(bytes32 sharedId, uint256 status, uint256 blockNumber, string message, uint256 index)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackSharedInfoAddedEvent(log *types.Log) (*PNCommunicatorV1SharedInfoAdded, error) {
	event := "SharedInfoAdded"
	if log.Topics[0] != pNCommunicatorV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNCommunicatorV1SharedInfoAdded)
	if len(log.Data) > 0 {
		if err := pNCommunicatorV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNCommunicatorV1.abi.Events[event].Inputs {
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

// PNCommunicatorV1SharedInfoRemoved represents a SharedInfoRemoved event raised by the PNCommunicatorV1 contract.
type PNCommunicatorV1SharedInfoRemoved struct {
	SharedId [32]byte
	Index    *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const PNCommunicatorV1SharedInfoRemovedEventName = "SharedInfoRemoved"

// ContractEventName returns the user-defined event name.
func (PNCommunicatorV1SharedInfoRemoved) ContractEventName() string {
	return PNCommunicatorV1SharedInfoRemovedEventName
}

// UnpackSharedInfoRemovedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SharedInfoRemoved(bytes32 sharedId, uint256 index)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackSharedInfoRemovedEvent(log *types.Log) (*PNCommunicatorV1SharedInfoRemoved, error) {
	event := "SharedInfoRemoved"
	if log.Topics[0] != pNCommunicatorV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNCommunicatorV1SharedInfoRemoved)
	if len(log.Data) > 0 {
		if err := pNCommunicatorV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNCommunicatorV1.abi.Events[event].Inputs {
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

// PNCommunicatorV1Upgraded represents a Upgraded event raised by the PNCommunicatorV1 contract.
type PNCommunicatorV1Upgraded struct {
	Implementation common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const PNCommunicatorV1UpgradedEventName = "Upgraded"

// ContractEventName returns the user-defined event name.
func (PNCommunicatorV1Upgraded) ContractEventName() string {
	return PNCommunicatorV1UpgradedEventName
}

// UnpackUpgradedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Upgraded(address indexed implementation)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackUpgradedEvent(log *types.Log) (*PNCommunicatorV1Upgraded, error) {
	event := "Upgraded"
	if log.Topics[0] != pNCommunicatorV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PNCommunicatorV1Upgraded)
	if len(log.Data) > 0 {
		if err := pNCommunicatorV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pNCommunicatorV1.abi.Events[event].Inputs {
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
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], pNCommunicatorV1.abi.Errors["AddressEmptyCode"].ID.Bytes()[:4]) {
		return pNCommunicatorV1.UnpackAddressEmptyCodeError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNCommunicatorV1.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return pNCommunicatorV1.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNCommunicatorV1.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return pNCommunicatorV1.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNCommunicatorV1.abi.Errors["FailedCall"].ID.Bytes()[:4]) {
		return pNCommunicatorV1.UnpackFailedCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNCommunicatorV1.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return pNCommunicatorV1.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNCommunicatorV1.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return pNCommunicatorV1.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNCommunicatorV1.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return pNCommunicatorV1.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNCommunicatorV1.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return pNCommunicatorV1.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNCommunicatorV1.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return pNCommunicatorV1.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNCommunicatorV1.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return pNCommunicatorV1.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNCommunicatorV1.abi.Errors["UUPSUnauthorizedCallContext"].ID.Bytes()[:4]) {
		return pNCommunicatorV1.UnpackUUPSUnauthorizedCallContextError(raw[4:])
	}
	if bytes.Equal(raw[:4], pNCommunicatorV1.abi.Errors["UUPSUnsupportedProxiableUUID"].ID.Bytes()[:4]) {
		return pNCommunicatorV1.UnpackUUPSUnsupportedProxiableUUIDError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// PNCommunicatorV1AddressEmptyCode represents a AddressEmptyCode error raised by the PNCommunicatorV1 contract.
type PNCommunicatorV1AddressEmptyCode struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressEmptyCode(address target)
func PNCommunicatorV1AddressEmptyCodeErrorID() common.Hash {
	return common.HexToHash("0x9996b315c842ff135b8fc4a08ad5df1c344efbc03d2687aecc0678050d2aac89")
}

// UnpackAddressEmptyCodeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressEmptyCode(address target)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackAddressEmptyCodeError(raw []byte) (*PNCommunicatorV1AddressEmptyCode, error) {
	out := new(PNCommunicatorV1AddressEmptyCode)
	if err := pNCommunicatorV1.abi.UnpackIntoInterface(out, "AddressEmptyCode", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNCommunicatorV1ERC1967InvalidImplementation represents a ERC1967InvalidImplementation error raised by the PNCommunicatorV1 contract.
type PNCommunicatorV1ERC1967InvalidImplementation struct {
	Implementation common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func PNCommunicatorV1ERC1967InvalidImplementationErrorID() common.Hash {
	return common.HexToHash("0x4c9c8ce3ceb3130f17f7cdba48d89b5b0129f266a8bac114e6e315a41879b617")
}

// UnpackERC1967InvalidImplementationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackERC1967InvalidImplementationError(raw []byte) (*PNCommunicatorV1ERC1967InvalidImplementation, error) {
	out := new(PNCommunicatorV1ERC1967InvalidImplementation)
	if err := pNCommunicatorV1.abi.UnpackIntoInterface(out, "ERC1967InvalidImplementation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNCommunicatorV1ERC1967NonPayable represents a ERC1967NonPayable error raised by the PNCommunicatorV1 contract.
type PNCommunicatorV1ERC1967NonPayable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967NonPayable()
func PNCommunicatorV1ERC1967NonPayableErrorID() common.Hash {
	return common.HexToHash("0xb398979fa84f543c8e222f17890372c487baf85e062276c127fef521eea7224b")
}

// UnpackERC1967NonPayableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967NonPayable()
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackERC1967NonPayableError(raw []byte) (*PNCommunicatorV1ERC1967NonPayable, error) {
	out := new(PNCommunicatorV1ERC1967NonPayable)
	if err := pNCommunicatorV1.abi.UnpackIntoInterface(out, "ERC1967NonPayable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNCommunicatorV1FailedCall represents a FailedCall error raised by the PNCommunicatorV1 contract.
type PNCommunicatorV1FailedCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedCall()
func PNCommunicatorV1FailedCallErrorID() common.Hash {
	return common.HexToHash("0xd6bda27508c0fb6d8a39b4b122878dab26f731a7d4e4abe711dd3731899052a4")
}

// UnpackFailedCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedCall()
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackFailedCallError(raw []byte) (*PNCommunicatorV1FailedCall, error) {
	out := new(PNCommunicatorV1FailedCall)
	if err := pNCommunicatorV1.abi.UnpackIntoInterface(out, "FailedCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNCommunicatorV1InvalidInitialization represents a InvalidInitialization error raised by the PNCommunicatorV1 contract.
type PNCommunicatorV1InvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func PNCommunicatorV1InvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackInvalidInitializationError(raw []byte) (*PNCommunicatorV1InvalidInitialization, error) {
	out := new(PNCommunicatorV1InvalidInitialization)
	if err := pNCommunicatorV1.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNCommunicatorV1NotInitializing represents a NotInitializing error raised by the PNCommunicatorV1 contract.
type PNCommunicatorV1NotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func PNCommunicatorV1NotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackNotInitializingError(raw []byte) (*PNCommunicatorV1NotInitializing, error) {
	out := new(PNCommunicatorV1NotInitializing)
	if err := pNCommunicatorV1.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNCommunicatorV1RaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the PNCommunicatorV1 contract.
type PNCommunicatorV1RaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func PNCommunicatorV1RaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*PNCommunicatorV1RaylsAccessManagedContractPaused, error) {
	out := new(PNCommunicatorV1RaylsAccessManagedContractPaused)
	if err := pNCommunicatorV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNCommunicatorV1RaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the PNCommunicatorV1 contract.
type PNCommunicatorV1RaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func PNCommunicatorV1RaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*PNCommunicatorV1RaylsAccessManagedInvalidAuthority, error) {
	out := new(PNCommunicatorV1RaylsAccessManagedInvalidAuthority)
	if err := pNCommunicatorV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNCommunicatorV1RaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the PNCommunicatorV1 contract.
type PNCommunicatorV1RaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func PNCommunicatorV1RaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*PNCommunicatorV1RaylsAccessManagedMustSchedule, error) {
	out := new(PNCommunicatorV1RaylsAccessManagedMustSchedule)
	if err := pNCommunicatorV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNCommunicatorV1RaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the PNCommunicatorV1 contract.
type PNCommunicatorV1RaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func PNCommunicatorV1RaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*PNCommunicatorV1RaylsAccessManagedUnauthorized, error) {
	out := new(PNCommunicatorV1RaylsAccessManagedUnauthorized)
	if err := pNCommunicatorV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNCommunicatorV1UUPSUnauthorizedCallContext represents a UUPSUnauthorizedCallContext error raised by the PNCommunicatorV1 contract.
type PNCommunicatorV1UUPSUnauthorizedCallContext struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnauthorizedCallContext()
func PNCommunicatorV1UUPSUnauthorizedCallContextErrorID() common.Hash {
	return common.HexToHash("0xe07c8dba242a06571ac65fe4bbe20522c9fb111cb33599b799ff8039c1ed18f4")
}

// UnpackUUPSUnauthorizedCallContextError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnauthorizedCallContext()
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackUUPSUnauthorizedCallContextError(raw []byte) (*PNCommunicatorV1UUPSUnauthorizedCallContext, error) {
	out := new(PNCommunicatorV1UUPSUnauthorizedCallContext)
	if err := pNCommunicatorV1.abi.UnpackIntoInterface(out, "UUPSUnauthorizedCallContext", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PNCommunicatorV1UUPSUnsupportedProxiableUUID represents a UUPSUnsupportedProxiableUUID error raised by the PNCommunicatorV1 contract.
type PNCommunicatorV1UUPSUnsupportedProxiableUUID struct {
	Slot [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func PNCommunicatorV1UUPSUnsupportedProxiableUUIDErrorID() common.Hash {
	return common.HexToHash("0xaa1d49a4c084bfa9aeeee2a0be65267a7f19ba7e1476b114dac513d2c14cb563")
}

// UnpackUUPSUnsupportedProxiableUUIDError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func (pNCommunicatorV1 *PNCommunicatorV1) UnpackUUPSUnsupportedProxiableUUIDError(raw []byte) (*PNCommunicatorV1UUPSUnsupportedProxiableUUID, error) {
	out := new(PNCommunicatorV1UUPSUnsupportedProxiableUUID)
	if err := pNCommunicatorV1.abi.UnpackIntoInterface(out, "UUPSUnsupportedProxiableUUID", raw); err != nil {
		return nil, err
	}
	return out, nil
}
