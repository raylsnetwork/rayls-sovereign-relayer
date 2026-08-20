// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package PublicRNEndpointV1

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

// RaylsNodeBridgedTransferMetadata is an auto generated low-level Go binding around an user-defined struct.
type RaylsNodeBridgedTransferMetadata struct {
	AssetType    uint8
	Id           *big.Int
	From         common.Address
	To           common.Address
	TokenAddress common.Address
	Amount       *big.Int
}

// RaylsNodeMessage is an auto generated low-level Go binding around an user-defined struct.
type RaylsNodeMessage struct {
	MessageMetadata RaylsNodeMessageMetadata
	Payload         []byte
}

// RaylsNodeMessageMetadata is an auto generated low-level Go binding around an user-defined struct.
type RaylsNodeMessageMetadata struct {
	Nonce               *big.Int
	NewResourceMetadata RaylsNodeNewResourceMetadata
	RevertPayloadData   []byte
	TransferMetadata    RaylsNodeBridgedTransferMetadata
}

// RaylsNodeNewResourceMetadata is an auto generated low-level Go binding around an user-defined struct.
type RaylsNodeNewResourceMetadata struct {
	ResourceDeployType uint8
	Bytecode           []byte
	FactoryTemplate    uint8
	InitializerParams  []byte
}

// PublicRNEndpointV1MetaData contains all meta data concerning the PublicRNEndpointV1 contract.
var PublicRNEndpointV1MetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"configureContracts\",\"inputs\":[{\"name\":\"_messageExecutor\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_messageDispatcher\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"contractVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"currentChainId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getChainId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMessageDispatcherAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNonce\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"authority_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"messageDispatcher\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIMessageDispatcher\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"messageExecutor\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractRNMessageExecutorV1\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nonce\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"receivePayload\",\"inputs\":[{\"name\":\"_srcChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_srcAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_dstAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_raylsMessage\",\"type\":\"tuple\",\"internalType\":\"structRaylsNodeMessage\",\"components\":[{\"name\":\"messageMetadata\",\"type\":\"tuple\",\"internalType\":\"structRaylsNodeMessageMetadata\",\"components\":[{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"newResourceMetadata\",\"type\":\"tuple\",\"internalType\":\"structRaylsNodeNewResourceMetadata\",\"components\":[{\"name\":\"resourceDeployType\",\"type\":\"uint8\",\"internalType\":\"enumRaylsNodeResourceDeployType\"},{\"name\":\"bytecode\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"factoryTemplate\",\"type\":\"uint8\",\"internalType\":\"enumRaylsNodeBridgeableERC\"},{\"name\":\"initializerParams\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"revertPayloadData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"transferMetadata\",\"type\":\"tuple\",\"internalType\":\"structRaylsNodeBridgedTransferMetadata\",\"components\":[{\"name\":\"assetType\",\"type\":\"uint8\",\"internalType\":\"enumRaylsNodeBridgeableERC\"},{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"_messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"send\",\"inputs\":[{\"name\":\"_dstChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_destination\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendToAddress\",\"inputs\":[{\"name\":\"_dstChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_privateChainAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_revertDataPayload\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"transferMetadata\",\"type\":\"tuple\",\"internalType\":\"structRaylsNodeBridgedTransferMetadata\",\"components\":[{\"name\":\"assetType\",\"type\":\"uint8\",\"internalType\":\"enumRaylsNodeBridgeableERC\"},{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PublicRNEndpointV1__InvalidMessageDispatcherAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PublicRNEndpointV1__InvalidMessageExecutorAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PublicRNEndpointV1__SourceAndDestinationChainsSame\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	ID:  "PublicRNEndpointV1",
	Bin: "0x60a06040523060805234801561001457600080fd5b5061001d610022565b6100d4565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff16156100725760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146100d15780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b6080516117546100fd60003960008181610960015281816109890152610ac201526117546000f3fe6080604052600436106100e95760003560e01c80636cbadbfa116100855780636cbadbfa14610220578063a0a8e46014610236578063ad3cb1cc1461024a578063ad552e321461027b578063affed0e01461029b578063bf7e214f146102b1578063d087d288146102c6578063da35a26f146102db578063e0d9d846146102fb57600080fd5b8063150b375f146100ee578063288cbb5d146101215780633408e4701461014357806337ffde051461015857806348628eed146101785780634f1ef286146101a557806352d1902d146101b857806354fd4d50146101cd5780636431087a14610202575b600080fd5b3480156100fa57600080fd5b5061010e610109366004610f1d565b61031b565b6040519081526020015b60405180910390f35b34801561012d57600080fd5b5061014161013c366004610f76565b6103f6565b005b34801561014f57600080fd5b5060005461010e565b34801561016457600080fd5b50610141610173366004611137565b610488565b34801561018457600080fd5b50600254610198906001600160a01b031681565b60405161011891906112e1565b6101416101b33660046112f5565b610513565b3480156101c457600080fd5b5061010e610532565b3480156101d957600080fd5b5060408051808201909152600381526219171b60e91b60208201525b6040516101189190611392565b34801561020e57600080fd5b506002546001600160a01b0316610198565b34801561022c57600080fd5b5061010e60005481565b34801561024257600080fd5b50600161010e565b34801561025657600080fd5b506101f5604051806040016040528060058152602001640352e302e360dc1b81525081565b34801561028757600080fd5b50600154610198906001600160a01b031681565b3480156102a757600080fd5b5061010e60035481565b3480156102bd57600080fd5b5061019861054f565b3480156102d257600080fd5b5060035461010e565b3480156102e757600080fd5b506101416102f63660046113a5565b61055e565b34801561030757600080fd5b5061010e6103163660046113c8565b610666565b6000610333336000356001600160e01b031916610707565b61033b610e90565b6040805160c081018252600080825260208201819052918101829052606081018290526080810182905260a08101919091526103eb6040518060c00160405280898152602001886001600160a01b0316815260200187878080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920182905250938552505050602080830187905260408051918201815291815290820152606001839052610852565b979650505050505050565b61040c336000356001600160e01b031916610707565b6001600160a01b038216610433576040516331dd933f60e21b815260040160405180910390fd5b6001600160a01b03811661045a57604051637a6a228960e11b815260040160405180910390fd5b600180546001600160a01b039384166001600160a01b03199182161790915560028054929093169116179055565b61049e336000356001600160e01b031916610707565b60015460208301516040516390fe4cff60e01b81526001600160a01b03909216916390fe4cff916104da9187919086908b908b9060040161145c565b600060405180830381600087803b1580156104f457600080fd5b505af1158015610508573d6000803e3d6000fd5b505050505050505050565b61051b610955565b610524826109e5565b61052e82826109fe565b5050565b600061053c610ab7565b506000805160206116ff83398151915290565b6000610559610b00565b905090565b6000610568610b19565b805490915060ff600160401b82041615906001600160401b031660008115801561058f5750825b90506000826001600160401b031660011480156105ab5750303b155b9050811580156105b9575080155b156105d75760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561060157845460ff60401b1916600160401b1785555b610609610b44565b600087905561061786610b4c565b831561065d57845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50505050505050565b600061067e336000356001600160e01b031916610707565b610686610e90565b6106fb6040518060c001604052808a8152602001896001600160a01b0316815260200188888080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505050908252506020810184905260408101879052606001859052610852565b98975050505050505050565b6000610711610b8d565b80549091506001600160a01b031680610749576000604051638944034760e01b815260040161074091906112e1565b60405180910390fd5b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa1580156107ad573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107d191906114ad565b9250925092508261065d5780156107fb5760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff8216156108375760405163a426878960e01b81526001600160a01b038816600482015263ffffffff83166024820152604401610740565b86604051632ecd3d0360e21b815260040161074091906112e1565b60008060036000815461086490611511565b91829055506040805160c081018252808201838152606087810151908301526080878101519083015260a08088015190830152815290850151602082015260005485519293509091036108ca57604051635b608f0160e11b815260040160405180910390fd5b6002546000805486516020880151604051630c5a8fc560e21b815293946001600160a01b03169363316a3f1493610909939092339289906004016115a5565b6020604051808303816000875af1158015610928573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061094c91906116b6565b95945050505050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614806109c557507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166109b9610bef565b6001600160a01b031614155b156109e35760405163703e46dd60e11b815260040160405180910390fd5b565b6109fb336000356001600160e01b031916610707565b50565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610a58575060408051601f3d908101601f19168201909252610a55918101906116b6565b60015b610a775781604051634c9c8ce360e01b815260040161074091906112e1565b6000805160206116ff8339815191528114610aa857604051632a87526960e21b815260048101829052602401610740565b610ab28383610c05565b505050565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146109e35760405163703e46dd60e11b815260040160405180910390fd5b6000610b0a610b8d565b546001600160a01b0316919050565b6000807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005b92915050565b6109e3610c5b565b6000610b56610b8d565b80549091506001600160a01b031615610b845781604051638944034760e01b815260040161074091906112e1565b61052e82610c80565b60008060ff19610bbe60017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f356116cf565b604051602001610bd091815260200190565b60408051601f1981840301815291905280516020909101201692915050565b60006000805160206116ff833981519152610b0a565b610c0e82610d10565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a2805115610c5357610ab28282610d6c565b61052e610dd9565b610c63610df8565b6109e357604051631afcd79f60e31b815260040160405180910390fd5b6001600160a01b038116610ca95780604051638944034760e01b815260040161074091906112e1565b6000610cb3610b8d565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b806001600160a01b03163b600003610d3d5780604051634c9c8ce360e01b815260040161074091906112e1565b6000805160206116ff83398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080846001600160a01b031684604051610d8991906116e2565b600060405180830381855af49150503d8060008114610dc4576040519150601f19603f3d011682016040523d82523d6000602084013e610dc9565b606091505b509150915061094c858383610e12565b34156109e35760405163b398979f60e01b815260040160405180910390fd5b6000610e02610b19565b54600160401b900460ff16919050565b606082610e2757610e2282610e68565b610e61565b8151158015610e3e57506001600160a01b0384163b155b15610e5e5783604051639996b31560e01b815260040161074091906112e1565b50805b9392505050565b805115610e7757805160208201fd5b60405163d6bda27560e01b815260040160405180910390fd5b604080516080810190915280600081526060602082015260400160008152602001606081525090565b80356001600160a01b0381168114610ed057600080fd5b919050565b60008083601f840112610ee757600080fd5b5081356001600160401b03811115610efe57600080fd5b602083019150836020828501011115610f1657600080fd5b9250929050565b60008060008060608587031215610f3357600080fd5b84359350610f4360208601610eb9565b925060408501356001600160401b03811115610f5e57600080fd5b610f6a87828801610ed5565b95989497509550505050565b60008060408385031215610f8957600080fd5b610f9283610eb9565b9150610fa060208401610eb9565b90509250929050565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b0381118282101715610fe157610fe1610fa9565b60405290565b604051608081016001600160401b0381118282101715610fe157610fe1610fa9565b600082601f83011261101a57600080fd5b81356001600160401b038082111561103457611034610fa9565b604051601f8301601f19908116603f0116810190828211818310171561105c5761105c610fa9565b8160405283815286602085880101111561107557600080fd5b836020870160208301376000602085830101528094505050505092915050565b803560078110610ed057600080fd5b600060c082840312156110b657600080fd5b60405160c081018181106001600160401b03821117156110d8576110d8610fa9565b6040529050806110e783611095565b8152602083013560208201526110ff60408401610eb9565b604082015261111060608401610eb9565b606082015261112160808401610eb9565b608082015260a083013560a08201525092915050565b600080600080600060a0868803121561114f57600080fd5b8535945061115f60208701610eb9565b935061116d60408701610eb9565b925060608601356001600160401b038082111561118957600080fd5b908701906040828a03121561119d57600080fd5b6111a5610fbf565b81833511156111b357600080fd5b82358301610120818c0312156111c857600080fd5b6111d0610fe7565b813581526020820135848111156111e657600080fd5b82016080818e0312156111f857600080fd5b611200610fe7565b600282351061120e57600080fd5b8135815260208201358681111561122457600080fd5b6112308f828501611009565b60208301525061124260408301611095565b604082015260608201358681111561125957600080fd5b6112658f828501611009565b60608301525060208301525060408201358481111561128357600080fd5b61128f8d828501611009565b6040830152506112a28c606084016110a4565b60608201528252506020830135828111156112bc57600080fd5b6112c88b828601611009565b6020830152509699959850939660800135949350505050565b6001600160a01b0391909116815260200190565b6000806040838503121561130857600080fd5b61131183610eb9565b915060208301356001600160401b0381111561132c57600080fd5b61133885828601611009565b9150509250929050565b60005b8381101561135d578181015183820152602001611345565b50506000910152565b6000815180845261137e816020860160208601611342565b601f01601f19169290920160200192915050565b602081526000610e616020830184611366565b600080604083850312156113b857600080fd5b82359150610fa060208401610eb9565b60008060008060008061014087890312156113e257600080fd5b863595506113f260208801610eb9565b945060408701356001600160401b038082111561140e57600080fd5b61141a8a838b01610ed5565b9096509450606089013591508082111561143357600080fd5b5061144089828a01611009565b92505061145088608089016110a4565b90509295509295509295565b600060018060a01b03808816835260a0602084015261147e60a0840188611366565b6040840196909652606083019490945250911660809091015292915050565b80518015158114610ed057600080fd5b6000806000606084860312156114c257600080fd5b6114cb8461149d565b9250602084015163ffffffff811681146114e457600080fd5b91506114f26040850161149d565b90509250925092565b634e487b7160e01b600052601160045260246000fd5b600060018201611523576115236114fb565b5060010190565b634e487b7160e01b600052602160045260246000fd5b600781106115505761155061152a565b9052565b61155f828251611540565b60208101516020830152604081015160018060a01b038082166040850152806060840151166060850152806080840151166080850152505060a081015160a08301525050565b858152600060018060a01b03808716602084015285604084015280851660608401525060a060808301528251604060a0840152805160e08401526020810151610120806101008601528151600281106116005761160061152a565b61020086015260208201516080610220870152611621610280870182611366565b90506040830151611636610240880182611540565b50606083015192506101ff19868203016102608701526116568184611366565b925050604083015160df1986840301828701526116738382611366565b925050506060820151915061168c610140850183611554565b6020850151848203609f190160c086015291506116a98183611366565b9998505050505050505050565b6000602082840312156116c857600080fd5b5051919050565b81810381811115610b3e57610b3e6114fb565b600082516116f4818460208701611342565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbca2646970667358221220e929cbe1d4ca60c9069200b759055b852344330dc6b742285ccad2bb67a3a91764736f6c63430008180033",
}

// PublicRNEndpointV1 is an auto generated Go binding around an Ethereum contract.
type PublicRNEndpointV1 struct {
	abi abi.ABI
}

// NewPublicRNEndpointV1 creates a new instance of PublicRNEndpointV1.
func NewPublicRNEndpointV1() *PublicRNEndpointV1 {
	parsed, err := PublicRNEndpointV1MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &PublicRNEndpointV1{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *PublicRNEndpointV1) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackUPGRADEINTERFACEVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (publicRNEndpointV1 *PublicRNEndpointV1) PackUPGRADEINTERFACEVERSION() []byte {
	enc, err := publicRNEndpointV1.abi.Pack("UPGRADE_INTERFACE_VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUPGRADEINTERFACEVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackUPGRADEINTERFACEVERSION(data []byte) (string, error) {
	out, err := publicRNEndpointV1.abi.Unpack("UPGRADE_INTERFACE_VERSION", data)
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
func (publicRNEndpointV1 *PublicRNEndpointV1) PackAuthority() []byte {
	enc, err := publicRNEndpointV1.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := publicRNEndpointV1.abi.Unpack("authority", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackConfigureContracts is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x288cbb5d.
//
// Solidity: function configureContracts(address _messageExecutor, address _messageDispatcher) returns()
func (publicRNEndpointV1 *PublicRNEndpointV1) PackConfigureContracts(messageExecutor common.Address, messageDispatcher common.Address) []byte {
	enc, err := publicRNEndpointV1.abi.Pack("configureContracts", messageExecutor, messageDispatcher)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackContractVersion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa0a8e460.
//
// Solidity: function contractVersion() pure returns(uint256)
func (publicRNEndpointV1 *PublicRNEndpointV1) PackContractVersion() []byte {
	enc, err := publicRNEndpointV1.abi.Pack("contractVersion")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackContractVersion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa0a8e460.
//
// Solidity: function contractVersion() pure returns(uint256)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackContractVersion(data []byte) (*big.Int, error) {
	out, err := publicRNEndpointV1.abi.Unpack("contractVersion", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackCurrentChainId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6cbadbfa.
//
// Solidity: function currentChainId() view returns(uint256)
func (publicRNEndpointV1 *PublicRNEndpointV1) PackCurrentChainId() []byte {
	enc, err := publicRNEndpointV1.abi.Pack("currentChainId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackCurrentChainId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6cbadbfa.
//
// Solidity: function currentChainId() view returns(uint256)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackCurrentChainId(data []byte) (*big.Int, error) {
	out, err := publicRNEndpointV1.abi.Unpack("currentChainId", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetChainId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3408e470.
//
// Solidity: function getChainId() view returns(uint256)
func (publicRNEndpointV1 *PublicRNEndpointV1) PackGetChainId() []byte {
	enc, err := publicRNEndpointV1.abi.Pack("getChainId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetChainId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3408e470.
//
// Solidity: function getChainId() view returns(uint256)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackGetChainId(data []byte) (*big.Int, error) {
	out, err := publicRNEndpointV1.abi.Unpack("getChainId", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetMessageDispatcherAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6431087a.
//
// Solidity: function getMessageDispatcherAddress() view returns(address)
func (publicRNEndpointV1 *PublicRNEndpointV1) PackGetMessageDispatcherAddress() []byte {
	enc, err := publicRNEndpointV1.abi.Pack("getMessageDispatcherAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetMessageDispatcherAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6431087a.
//
// Solidity: function getMessageDispatcherAddress() view returns(address)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackGetMessageDispatcherAddress(data []byte) (common.Address, error) {
	out, err := publicRNEndpointV1.abi.Unpack("getMessageDispatcherAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetNonce is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd087d288.
//
// Solidity: function getNonce() view returns(uint256)
func (publicRNEndpointV1 *PublicRNEndpointV1) PackGetNonce() []byte {
	enc, err := publicRNEndpointV1.abi.Pack("getNonce")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetNonce is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd087d288.
//
// Solidity: function getNonce() view returns(uint256)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackGetNonce(data []byte) (*big.Int, error) {
	out, err := publicRNEndpointV1.abi.Unpack("getNonce", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xda35a26f.
//
// Solidity: function initialize(uint256 _chainId, address authority_) returns()
func (publicRNEndpointV1 *PublicRNEndpointV1) PackInitialize(chainId *big.Int, authority common.Address) []byte {
	enc, err := publicRNEndpointV1.abi.Pack("initialize", chainId, authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackMessageDispatcher is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x48628eed.
//
// Solidity: function messageDispatcher() view returns(address)
func (publicRNEndpointV1 *PublicRNEndpointV1) PackMessageDispatcher() []byte {
	enc, err := publicRNEndpointV1.abi.Pack("messageDispatcher")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackMessageDispatcher is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x48628eed.
//
// Solidity: function messageDispatcher() view returns(address)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackMessageDispatcher(data []byte) (common.Address, error) {
	out, err := publicRNEndpointV1.abi.Unpack("messageDispatcher", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackMessageExecutor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad552e32.
//
// Solidity: function messageExecutor() view returns(address)
func (publicRNEndpointV1 *PublicRNEndpointV1) PackMessageExecutor() []byte {
	enc, err := publicRNEndpointV1.abi.Pack("messageExecutor")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackMessageExecutor is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad552e32.
//
// Solidity: function messageExecutor() view returns(address)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackMessageExecutor(data []byte) (common.Address, error) {
	out, err := publicRNEndpointV1.abi.Unpack("messageExecutor", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackNonce is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (publicRNEndpointV1 *PublicRNEndpointV1) PackNonce() []byte {
	enc, err := publicRNEndpointV1.abi.Pack("nonce")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackNonce is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackNonce(data []byte) (*big.Int, error) {
	out, err := publicRNEndpointV1.abi.Unpack("nonce", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackProxiableUUID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (publicRNEndpointV1 *PublicRNEndpointV1) PackProxiableUUID() []byte {
	enc, err := publicRNEndpointV1.abi.Pack("proxiableUUID")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackProxiableUUID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackProxiableUUID(data []byte) ([32]byte, error) {
	out, err := publicRNEndpointV1.abi.Unpack("proxiableUUID", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackReceivePayload is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x37ffde05.
//
// Solidity: function receivePayload(uint256 _srcChainId, address _srcAddress, address _dstAddress, ((uint256,(uint8,bytes,uint8,bytes),bytes,(uint8,uint256,address,address,address,uint256)),bytes) _raylsMessage, bytes32 _messageId) returns()
func (publicRNEndpointV1 *PublicRNEndpointV1) PackReceivePayload(srcChainId *big.Int, srcAddress common.Address, dstAddress common.Address, raylsMessage RaylsNodeMessage, messageId [32]byte) []byte {
	enc, err := publicRNEndpointV1.abi.Pack("receivePayload", srcChainId, srcAddress, dstAddress, raylsMessage, messageId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSend is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x150b375f.
//
// Solidity: function send(uint256 _dstChainId, address _destination, bytes _payload) returns(bytes32 messageId)
func (publicRNEndpointV1 *PublicRNEndpointV1) PackSend(dstChainId *big.Int, destination common.Address, payload []byte) []byte {
	enc, err := publicRNEndpointV1.abi.Pack("send", dstChainId, destination, payload)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSend is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x150b375f.
//
// Solidity: function send(uint256 _dstChainId, address _destination, bytes _payload) returns(bytes32 messageId)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackSend(data []byte) ([32]byte, error) {
	out, err := publicRNEndpointV1.abi.Unpack("send", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackSendToAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe0d9d846.
//
// Solidity: function sendToAddress(uint256 _dstChainId, address _privateChainAddress, bytes _payload, bytes _revertDataPayload, (uint8,uint256,address,address,address,uint256) transferMetadata) returns(bytes32 messageId)
func (publicRNEndpointV1 *PublicRNEndpointV1) PackSendToAddress(dstChainId *big.Int, privateChainAddress common.Address, payload []byte, revertDataPayload []byte, transferMetadata RaylsNodeBridgedTransferMetadata) []byte {
	enc, err := publicRNEndpointV1.abi.Pack("sendToAddress", dstChainId, privateChainAddress, payload, revertDataPayload, transferMetadata)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSendToAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe0d9d846.
//
// Solidity: function sendToAddress(uint256 _dstChainId, address _privateChainAddress, bytes _payload, bytes _revertDataPayload, (uint8,uint256,address,address,address,uint256) transferMetadata) returns(bytes32 messageId)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackSendToAddress(data []byte) ([32]byte, error) {
	out, err := publicRNEndpointV1.abi.Unpack("sendToAddress", data)
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
func (publicRNEndpointV1 *PublicRNEndpointV1) PackUpgradeToAndCall(newImplementation common.Address, data []byte) []byte {
	enc, err := publicRNEndpointV1.abi.Pack("upgradeToAndCall", newImplementation, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackVersion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x54fd4d50.
//
// Solidity: function version() pure returns(string)
func (publicRNEndpointV1 *PublicRNEndpointV1) PackVersion() []byte {
	enc, err := publicRNEndpointV1.abi.Pack("version")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackVersion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x54fd4d50.
//
// Solidity: function version() pure returns(string)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackVersion(data []byte) (string, error) {
	out, err := publicRNEndpointV1.abi.Unpack("version", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PublicRNEndpointV1AuthorityUpdated represents a AuthorityUpdated event raised by the PublicRNEndpointV1 contract.
type PublicRNEndpointV1AuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const PublicRNEndpointV1AuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (PublicRNEndpointV1AuthorityUpdated) ContractEventName() string {
	return PublicRNEndpointV1AuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackAuthorityUpdatedEvent(log *types.Log) (*PublicRNEndpointV1AuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != publicRNEndpointV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PublicRNEndpointV1AuthorityUpdated)
	if len(log.Data) > 0 {
		if err := publicRNEndpointV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range publicRNEndpointV1.abi.Events[event].Inputs {
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

// PublicRNEndpointV1Initialized represents a Initialized event raised by the PublicRNEndpointV1 contract.
type PublicRNEndpointV1Initialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const PublicRNEndpointV1InitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (PublicRNEndpointV1Initialized) ContractEventName() string {
	return PublicRNEndpointV1InitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackInitializedEvent(log *types.Log) (*PublicRNEndpointV1Initialized, error) {
	event := "Initialized"
	if log.Topics[0] != publicRNEndpointV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PublicRNEndpointV1Initialized)
	if len(log.Data) > 0 {
		if err := publicRNEndpointV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range publicRNEndpointV1.abi.Events[event].Inputs {
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

// PublicRNEndpointV1Upgraded represents a Upgraded event raised by the PublicRNEndpointV1 contract.
type PublicRNEndpointV1Upgraded struct {
	Implementation common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const PublicRNEndpointV1UpgradedEventName = "Upgraded"

// ContractEventName returns the user-defined event name.
func (PublicRNEndpointV1Upgraded) ContractEventName() string {
	return PublicRNEndpointV1UpgradedEventName
}

// UnpackUpgradedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Upgraded(address indexed implementation)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackUpgradedEvent(log *types.Log) (*PublicRNEndpointV1Upgraded, error) {
	event := "Upgraded"
	if log.Topics[0] != publicRNEndpointV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PublicRNEndpointV1Upgraded)
	if len(log.Data) > 0 {
		if err := publicRNEndpointV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range publicRNEndpointV1.abi.Events[event].Inputs {
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
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], publicRNEndpointV1.abi.Errors["AddressEmptyCode"].ID.Bytes()[:4]) {
		return publicRNEndpointV1.UnpackAddressEmptyCodeError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicRNEndpointV1.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return publicRNEndpointV1.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicRNEndpointV1.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return publicRNEndpointV1.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicRNEndpointV1.abi.Errors["FailedCall"].ID.Bytes()[:4]) {
		return publicRNEndpointV1.UnpackFailedCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicRNEndpointV1.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return publicRNEndpointV1.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicRNEndpointV1.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return publicRNEndpointV1.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicRNEndpointV1.abi.Errors["PublicRNEndpointV1InvalidMessageDispatcherAddress"].ID.Bytes()[:4]) {
		return publicRNEndpointV1.UnpackPublicRNEndpointV1InvalidMessageDispatcherAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicRNEndpointV1.abi.Errors["PublicRNEndpointV1InvalidMessageExecutorAddress"].ID.Bytes()[:4]) {
		return publicRNEndpointV1.UnpackPublicRNEndpointV1InvalidMessageExecutorAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicRNEndpointV1.abi.Errors["PublicRNEndpointV1SourceAndDestinationChainsSame"].ID.Bytes()[:4]) {
		return publicRNEndpointV1.UnpackPublicRNEndpointV1SourceAndDestinationChainsSameError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicRNEndpointV1.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return publicRNEndpointV1.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicRNEndpointV1.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return publicRNEndpointV1.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicRNEndpointV1.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return publicRNEndpointV1.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicRNEndpointV1.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return publicRNEndpointV1.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicRNEndpointV1.abi.Errors["UUPSUnauthorizedCallContext"].ID.Bytes()[:4]) {
		return publicRNEndpointV1.UnpackUUPSUnauthorizedCallContextError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicRNEndpointV1.abi.Errors["UUPSUnsupportedProxiableUUID"].ID.Bytes()[:4]) {
		return publicRNEndpointV1.UnpackUUPSUnsupportedProxiableUUIDError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// PublicRNEndpointV1AddressEmptyCode represents a AddressEmptyCode error raised by the PublicRNEndpointV1 contract.
type PublicRNEndpointV1AddressEmptyCode struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressEmptyCode(address target)
func PublicRNEndpointV1AddressEmptyCodeErrorID() common.Hash {
	return common.HexToHash("0x9996b315c842ff135b8fc4a08ad5df1c344efbc03d2687aecc0678050d2aac89")
}

// UnpackAddressEmptyCodeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressEmptyCode(address target)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackAddressEmptyCodeError(raw []byte) (*PublicRNEndpointV1AddressEmptyCode, error) {
	out := new(PublicRNEndpointV1AddressEmptyCode)
	if err := publicRNEndpointV1.abi.UnpackIntoInterface(out, "AddressEmptyCode", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicRNEndpointV1ERC1967InvalidImplementation represents a ERC1967InvalidImplementation error raised by the PublicRNEndpointV1 contract.
type PublicRNEndpointV1ERC1967InvalidImplementation struct {
	Implementation common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func PublicRNEndpointV1ERC1967InvalidImplementationErrorID() common.Hash {
	return common.HexToHash("0x4c9c8ce3ceb3130f17f7cdba48d89b5b0129f266a8bac114e6e315a41879b617")
}

// UnpackERC1967InvalidImplementationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackERC1967InvalidImplementationError(raw []byte) (*PublicRNEndpointV1ERC1967InvalidImplementation, error) {
	out := new(PublicRNEndpointV1ERC1967InvalidImplementation)
	if err := publicRNEndpointV1.abi.UnpackIntoInterface(out, "ERC1967InvalidImplementation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicRNEndpointV1ERC1967NonPayable represents a ERC1967NonPayable error raised by the PublicRNEndpointV1 contract.
type PublicRNEndpointV1ERC1967NonPayable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967NonPayable()
func PublicRNEndpointV1ERC1967NonPayableErrorID() common.Hash {
	return common.HexToHash("0xb398979fa84f543c8e222f17890372c487baf85e062276c127fef521eea7224b")
}

// UnpackERC1967NonPayableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967NonPayable()
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackERC1967NonPayableError(raw []byte) (*PublicRNEndpointV1ERC1967NonPayable, error) {
	out := new(PublicRNEndpointV1ERC1967NonPayable)
	if err := publicRNEndpointV1.abi.UnpackIntoInterface(out, "ERC1967NonPayable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicRNEndpointV1FailedCall represents a FailedCall error raised by the PublicRNEndpointV1 contract.
type PublicRNEndpointV1FailedCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedCall()
func PublicRNEndpointV1FailedCallErrorID() common.Hash {
	return common.HexToHash("0xd6bda27508c0fb6d8a39b4b122878dab26f731a7d4e4abe711dd3731899052a4")
}

// UnpackFailedCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedCall()
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackFailedCallError(raw []byte) (*PublicRNEndpointV1FailedCall, error) {
	out := new(PublicRNEndpointV1FailedCall)
	if err := publicRNEndpointV1.abi.UnpackIntoInterface(out, "FailedCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicRNEndpointV1InvalidInitialization represents a InvalidInitialization error raised by the PublicRNEndpointV1 contract.
type PublicRNEndpointV1InvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func PublicRNEndpointV1InvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackInvalidInitializationError(raw []byte) (*PublicRNEndpointV1InvalidInitialization, error) {
	out := new(PublicRNEndpointV1InvalidInitialization)
	if err := publicRNEndpointV1.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicRNEndpointV1NotInitializing represents a NotInitializing error raised by the PublicRNEndpointV1 contract.
type PublicRNEndpointV1NotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func PublicRNEndpointV1NotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackNotInitializingError(raw []byte) (*PublicRNEndpointV1NotInitializing, error) {
	out := new(PublicRNEndpointV1NotInitializing)
	if err := publicRNEndpointV1.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicRNEndpointV1PublicRNEndpointV1InvalidMessageDispatcherAddress represents a PublicRNEndpointV1__InvalidMessageDispatcherAddress error raised by the PublicRNEndpointV1 contract.
type PublicRNEndpointV1PublicRNEndpointV1InvalidMessageDispatcherAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PublicRNEndpointV1__InvalidMessageDispatcherAddress()
func PublicRNEndpointV1PublicRNEndpointV1InvalidMessageDispatcherAddressErrorID() common.Hash {
	return common.HexToHash("0xf4d44512626ac10538e5d7ba676334a028e25c4c02ce2b9df433e9be96d529e6")
}

// UnpackPublicRNEndpointV1InvalidMessageDispatcherAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PublicRNEndpointV1__InvalidMessageDispatcherAddress()
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackPublicRNEndpointV1InvalidMessageDispatcherAddressError(raw []byte) (*PublicRNEndpointV1PublicRNEndpointV1InvalidMessageDispatcherAddress, error) {
	out := new(PublicRNEndpointV1PublicRNEndpointV1InvalidMessageDispatcherAddress)
	if err := publicRNEndpointV1.abi.UnpackIntoInterface(out, "PublicRNEndpointV1InvalidMessageDispatcherAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicRNEndpointV1PublicRNEndpointV1InvalidMessageExecutorAddress represents a PublicRNEndpointV1__InvalidMessageExecutorAddress error raised by the PublicRNEndpointV1 contract.
type PublicRNEndpointV1PublicRNEndpointV1InvalidMessageExecutorAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PublicRNEndpointV1__InvalidMessageExecutorAddress()
func PublicRNEndpointV1PublicRNEndpointV1InvalidMessageExecutorAddressErrorID() common.Hash {
	return common.HexToHash("0xc7764cfc68c4046d6a2027d8bab23ba8d2e109de5738eafabf77f88c7eb6fd45")
}

// UnpackPublicRNEndpointV1InvalidMessageExecutorAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PublicRNEndpointV1__InvalidMessageExecutorAddress()
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackPublicRNEndpointV1InvalidMessageExecutorAddressError(raw []byte) (*PublicRNEndpointV1PublicRNEndpointV1InvalidMessageExecutorAddress, error) {
	out := new(PublicRNEndpointV1PublicRNEndpointV1InvalidMessageExecutorAddress)
	if err := publicRNEndpointV1.abi.UnpackIntoInterface(out, "PublicRNEndpointV1InvalidMessageExecutorAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicRNEndpointV1PublicRNEndpointV1SourceAndDestinationChainsSame represents a PublicRNEndpointV1__SourceAndDestinationChainsSame error raised by the PublicRNEndpointV1 contract.
type PublicRNEndpointV1PublicRNEndpointV1SourceAndDestinationChainsSame struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PublicRNEndpointV1__SourceAndDestinationChainsSame()
func PublicRNEndpointV1PublicRNEndpointV1SourceAndDestinationChainsSameErrorID() common.Hash {
	return common.HexToHash("0xb6c11e02a045c8a579ef61482ee356cb8fe8befa24211a3eef502f5de6c7f075")
}

// UnpackPublicRNEndpointV1SourceAndDestinationChainsSameError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PublicRNEndpointV1__SourceAndDestinationChainsSame()
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackPublicRNEndpointV1SourceAndDestinationChainsSameError(raw []byte) (*PublicRNEndpointV1PublicRNEndpointV1SourceAndDestinationChainsSame, error) {
	out := new(PublicRNEndpointV1PublicRNEndpointV1SourceAndDestinationChainsSame)
	if err := publicRNEndpointV1.abi.UnpackIntoInterface(out, "PublicRNEndpointV1SourceAndDestinationChainsSame", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicRNEndpointV1RaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the PublicRNEndpointV1 contract.
type PublicRNEndpointV1RaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func PublicRNEndpointV1RaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*PublicRNEndpointV1RaylsAccessManagedContractPaused, error) {
	out := new(PublicRNEndpointV1RaylsAccessManagedContractPaused)
	if err := publicRNEndpointV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicRNEndpointV1RaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the PublicRNEndpointV1 contract.
type PublicRNEndpointV1RaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func PublicRNEndpointV1RaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*PublicRNEndpointV1RaylsAccessManagedInvalidAuthority, error) {
	out := new(PublicRNEndpointV1RaylsAccessManagedInvalidAuthority)
	if err := publicRNEndpointV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicRNEndpointV1RaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the PublicRNEndpointV1 contract.
type PublicRNEndpointV1RaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func PublicRNEndpointV1RaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*PublicRNEndpointV1RaylsAccessManagedMustSchedule, error) {
	out := new(PublicRNEndpointV1RaylsAccessManagedMustSchedule)
	if err := publicRNEndpointV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicRNEndpointV1RaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the PublicRNEndpointV1 contract.
type PublicRNEndpointV1RaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func PublicRNEndpointV1RaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*PublicRNEndpointV1RaylsAccessManagedUnauthorized, error) {
	out := new(PublicRNEndpointV1RaylsAccessManagedUnauthorized)
	if err := publicRNEndpointV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicRNEndpointV1UUPSUnauthorizedCallContext represents a UUPSUnauthorizedCallContext error raised by the PublicRNEndpointV1 contract.
type PublicRNEndpointV1UUPSUnauthorizedCallContext struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnauthorizedCallContext()
func PublicRNEndpointV1UUPSUnauthorizedCallContextErrorID() common.Hash {
	return common.HexToHash("0xe07c8dba242a06571ac65fe4bbe20522c9fb111cb33599b799ff8039c1ed18f4")
}

// UnpackUUPSUnauthorizedCallContextError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnauthorizedCallContext()
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackUUPSUnauthorizedCallContextError(raw []byte) (*PublicRNEndpointV1UUPSUnauthorizedCallContext, error) {
	out := new(PublicRNEndpointV1UUPSUnauthorizedCallContext)
	if err := publicRNEndpointV1.abi.UnpackIntoInterface(out, "UUPSUnauthorizedCallContext", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicRNEndpointV1UUPSUnsupportedProxiableUUID represents a UUPSUnsupportedProxiableUUID error raised by the PublicRNEndpointV1 contract.
type PublicRNEndpointV1UUPSUnsupportedProxiableUUID struct {
	Slot [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func PublicRNEndpointV1UUPSUnsupportedProxiableUUIDErrorID() common.Hash {
	return common.HexToHash("0xaa1d49a4c084bfa9aeeee2a0be65267a7f19ba7e1476b114dac513d2c14cb563")
}

// UnpackUUPSUnsupportedProxiableUUIDError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func (publicRNEndpointV1 *PublicRNEndpointV1) UnpackUUPSUnsupportedProxiableUUIDError(raw []byte) (*PublicRNEndpointV1UUPSUnsupportedProxiableUUID, error) {
	out := new(PublicRNEndpointV1UUPSUnsupportedProxiableUUID)
	if err := publicRNEndpointV1.abi.UnpackIntoInterface(out, "UUPSUnsupportedProxiableUUID", raw); err != nil {
		return nil, err
	}
	return out, nil
}
