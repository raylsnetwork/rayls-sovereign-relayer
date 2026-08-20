// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ParticipantStorageReplicaV1

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

// ParticipantStructsParticipant is an auto generated low-level Go binding around an user-defined struct.
type ParticipantStructsParticipant struct {
	ChainId            *big.Int
	Role               uint8
	Status             uint8
	OwnerId            string
	Name               string
	CreatedAt          *big.Int
	UpdatedAt          *big.Int
	AllowedToBroadcast bool
}

// ParticipantStorageReplicaV1MetaData contains all meta data concerning the ParticipantStorageReplicaV1 contract.
var ParticipantStorageReplicaV1MetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addOrUpdateParticipants\",\"inputs\":[{\"name\":\"_participants\",\"type\":\"tuple[]\",\"internalType\":\"structParticipantStructs.Participant[]\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"role\",\"type\":\"uint8\",\"internalType\":\"enumParticipantStructs.Role\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumParticipantStructs.Status\"},{\"name\":\"ownerId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"allowedToBroadcast\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"contractVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"endpoint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRaylsEndpoint\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAddressByResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllParticipants\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structParticipantStructs.Participant[]\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"role\",\"type\":\"uint8\",\"internalType\":\"enumParticipantStructs.Role\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumParticipantStructs.Status\"},{\"name\":\"ownerId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"allowedToBroadcast\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"authority_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestAllParticipantsDataFromPrivateHub\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resourceId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"validateMessageParticipants\",\"inputs\":[{\"name\":\"originChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validateParticipantStatus\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParticipantStorageReplicaV1__NotFromPrivateHub\",\"inputs\":[{\"name\":\"fromChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"hubId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ParticipantStorageReplicaV1__ParticipantNotActive\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ParticipantStorageReplicaV1__ParticipantNotRegistered\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	ID:  "ParticipantStorageReplicaV1",
	Bin: "0x60a06040523060805234801561001457600080fd5b50608051611c6c61003e60003960008181610a5f01528181610a880152610bbc0152611c6c6000f3fe6080604052600436106100b85760003560e01c806311f50c85146100bd578063195ec9ee146100f3578063485cc955146101155780634f1ef2861461013757806352d1902d1461014a5780635e280f111461016d5780635f997c5b1461018d578063899cb662146101a357806394b0fd87146101c3578063a0a8e460146101d8578063ad3cb1cc146101ec578063bf7e214f1461022a578063c4d66de81461023f578063c9557f561461025f578063d5c3614f1461027f575b600080fd5b3480156100c957600080fd5b506100dd6100d83660046113fe565b61029f565b6040516100ea9190611417565b60405180910390f35b3480156100ff57600080fd5b50610108610313565b6040516100ea91906114a5565b34801561012157600080fd5b506101356101303660046115a8565b61053c565b005b6101356101453660046116a7565b61064c565b34801561015657600080fd5b5061015f61066b565b6040519081526020016100ea565b34801561017957600080fd5b506000546100dd906001600160a01b031681565b34801561019957600080fd5b5061015f60015481565b3480156101af57600080fd5b506101356101be366004611766565b610689565b3480156101cf57600080fd5b50610135610790565b3480156101e457600080fd5b50600161015f565b3480156101f857600080fd5b5061021d604051806040016040528060058152602001640352e302e360dc1b81525081565b6040516100ea91906118d6565b34801561023657600080fd5b506100dd610885565b34801561024b57600080fd5b5061013561025a3660046118e9565b61089e565b34801561026b57600080fd5b5061013561027a3660046113fe565b6108c8565b34801561028b57600080fd5b5061013561029a366004611906565b61094f565b600080546040516311f50c8560e01b8152600481018490526001600160a01b03909116906311f50c8590602401602060405180830381865afa1580156102e9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061030d9190611928565b92915050565b60606002805480602002602001604051908101604052809291908181526020016000905b82821015610533578382906000526020600020906007020160405180610100016040529081600082015481526020016001820160009054906101000a900460ff1660028111156103895761038961142b565b600281111561039a5761039a61142b565b81526020016001820160019054906101000a900460ff1660038111156103c2576103c261142b565b60038111156103d3576103d361142b565b81526020016002820180546103e790611945565b80601f016020809104026020016040519081016040528092919081815260200182805461041390611945565b80156104605780601f1061043557610100808354040283529160200191610460565b820191906000526020600020905b81548152906001019060200180831161044357829003601f168201915b5050505050815260200160038201805461047990611945565b80601f01602080910402602001604051908101604052809291908181526020018280546104a590611945565b80156104f25780601f106104c7576101008083540402835291602001916104f2565b820191906000526020600020905b8154815290600101906020018083116104d557829003601f168201915b505050918352505060048201546020808301919091526005830154604083015260069092015460ff1615156060909101529082526001929092019101610337565b50505050905090565b60006105466109e0565b805490915060ff600160401b82041615906001600160401b031660008115801561056d5750825b90506000826001600160401b031660011480156105895750303b155b905081158015610597575080155b156105b55760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156105df57845460ff60401b1916600160401b1785555b6105e7610a09565b6105f08761089e565b600180556105fd86610a13565b831561064357845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50505050505050565b610654610a54565b61065d82610ae2565b6106678282610af8565b5050565b6000610675610bb1565b50600080516020611c178339815191525b90565b61069f336000356001600160e01b031916610bfa565b60006106a9610d3c565b905060008060009054906101000a90046001600160a01b03166001600160a01b0316630b39a9516040518163ffffffff1660e01b8152600401602060405180830381865afa1580156106ff573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610723919061197f565b905080821461075457604051631ce86f9d60e11b815260048101839052602481018290526044015b60405180910390fd5b60005b835181101561078a5761078284828151811061077557610775611998565b6020026020010151610d50565b600101610757565b50505050565b6040805160c08101825260008082526020808301829052828401829052606083018290526080830182905260a0830182905290548351630b39a95160e01b815293519293610882936001600160a01b0390921692630b39a951926004808401938290030181865afa158015610809573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061082d919061197f565b6001546040805160048152602481018252602080820180516001600160e01b03166306ec249d60e01b179052825180820184526000808252845180840186528181528551938401909552825291929087610fb6565b50565b600061088f61103f565b546001600160a01b0316919050565b6108a66110a1565b600080546001600160a01b0319166001600160a01b0392909216919091179055565b60008054906101000a90046001600160a01b03166001600160a01b0316630b39a9516040518163ffffffff1660e01b8152600401602060405180830381865afa158015610919573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061093d919061197f565b81036109465750565b610882816110c6565b60008054906101000a90046001600160a01b03166001600160a01b0316630b39a9516040518163ffffffff1660e01b8152600401602060405180830381865afa1580156109a0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109c4919061197f565b81036109ce575050565b6109d7826110c6565b610667816110c6565b6000807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0061030d565b610a116110a1565b565b6000610a1d61103f565b80549091506001600160a01b031615610a4b5781604051638944034760e01b815260040161074b9190611417565b61066782611179565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161480610ac457507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316610ab8611209565b6001600160a01b031614155b15610a115760405163703e46dd60e11b815260040160405180910390fd5b610882336000356001600160e01b031916610bfa565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610b52575060408051601f3d908101601f19168201909252610b4f9181019061197f565b60015b610b715781604051634c9c8ce360e01b815260040161074b9190611417565b600080516020611c178339815191528114610ba257604051632a87526960e21b81526004810182905260240161074b565b610bac838361121f565b505050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614610a115760405163703e46dd60e11b815260040160405180910390fd5b6000610c0461103f565b80549091506001600160a01b031680610c33576000604051638944034760e01b815260040161074b9190611417565b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa158015610c97573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610cbb91906119ae565b92509250925082610643578015610ce55760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff821615610d215760405163a426878960e01b81526001600160a01b038816600482015263ffffffff8316602482015260440161074b565b86604051632ecd3d0360e21b815260040161074b9190611417565b600060343610610686575060331936013590565b805160008181526004602052604081205490819003610ec45760028054600181810183556000839052855160079092027f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace810192835560208701517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5acf909101805488959293919260ff1990911691908490811115610df057610df061142b565b0217905550604082015160018201805461ff001916610100836003811115610e1a57610e1a61142b565b021790555060608201516002820190610e339082611a54565b5060808201516003820190610e489082611a54565b5060a082015160048281019190915560c0830151600583015560e0909201516006909101805460ff191691151591909117905560038054600181019091557fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b018390556002546000848152602092909252604090912055505050565b6000610ed1600183611b13565b90508360028281548110610ee757610ee7611998565b90600052602060002090600702016000820151816000015560208201518160010160006101000a81548160ff02191690836002811115610f2957610f2961142b565b0217905550604082015160018201805461ff001916610100836003811115610f5357610f5361142b565b021790555060608201516002820190610f6c9082611a54565b5060808201516003820190610f819082611a54565b5060a0820151600482015560c0820151600582015560e0909101516006909101805460ff191691151591909117905550505050565b600054604051631075c5d160e21b81526001600160a01b03909116906341d7174490610ff2908a908a908a908a908a908a908a90600401611b34565b6020604051808303816000875af1158015611011573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611035919061197f565b5050505050505050565b60008060ff1961107060017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f35611b13565b60405160200161108291815260200190565b60408051601f1981840301815291905280516020909101201692915050565b6110a9611275565b610a1157604051631afcd79f60e31b815260040160405180910390fd5b806000036110d15750565b600081815260046020526040812054908190036111045760405163914ee7d560e01b81526004810183905260240161074b565b6000611111600183611b13565b905060016002828154811061112857611128611998565b906000526020600020906007020160010160019054906101000a900460ff1660038111156111585761115861142b565b14610bac5760405163063a8c9560e31b81526004810184905260240161074b565b6001600160a01b0381166111a25780604051638944034760e01b815260040161074b9190611417565b60006111ac61103f565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b6000600080516020611c1783398151915261088f565b6112288261128f565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a280511561126d57610bac82826112eb565b610667611361565b600061127f6109e0565b54600160401b900460ff16919050565b806001600160a01b03163b6000036112bc5780604051634c9c8ce360e01b815260040161074b9190611417565b600080516020611c1783398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080846001600160a01b0316846040516113089190611bfa565b600060405180830381855af49150503d8060008114611343576040519150601f19603f3d011682016040523d82523d6000602084013e611348565b606091505b5091509150611358858383611380565b95945050505050565b3415610a115760405163b398979f60e01b815260040160405180910390fd5b60608261139557611390826113d6565b6113cf565b81511580156113ac57506001600160a01b0384163b155b156113cc5783604051639996b31560e01b815260040161074b9190611417565b50805b9392505050565b8051156113e557805160208201fd5b60405163d6bda27560e01b815260040160405180910390fd5b60006020828403121561141057600080fd5b5035919050565b6001600160a01b0391909116815260200190565b634e487b7160e01b600052602160045260246000fd5b600481106114515761145161142b565b9052565b60005b83811015611470578181015183820152602001611458565b50506000910152565b60008151808452611491816020860160208601611455565b601f01601f19169290920160200192915050565b600060208083018184528085518083526040925060408601915060408160051b87010184880160005b8381101561158557603f1989840301855281516101008151855288820151600381106114fc576114fc61142b565b858a01528188015161151089870182611441565b50606080830151828288015261152883880182611479565b92505050608080830151868303828801526115438382611479565b60a0858101519089015260c0808601519089015260e094850151801515868a01529490935091506115719050565b5095880195935050908601906001016114ce565b509098975050505050505050565b6001600160a01b038116811461088257600080fd5b600080604083850312156115bb57600080fd5b82356115c681611593565b915060208301356115d681611593565b809150509250929050565b634e487b7160e01b600052604160045260246000fd5b60405161010081016001600160401b038111828210171561161a5761161a6115e1565b60405290565b604051601f8201601f191681016001600160401b0381118282101715611648576116486115e1565b604052919050565b60006001600160401b03831115611669576116696115e1565b61167c601f8401601f1916602001611620565b905082815283838301111561169057600080fd5b828260208301376000602084830101529392505050565b600080604083850312156116ba57600080fd5b82356116c581611593565b915060208301356001600160401b038111156116e057600080fd5b8301601f810185136116f157600080fd5b61170085823560208401611650565b9150509250929050565b80356003811061171957600080fd5b919050565b80356004811061171957600080fd5b600082601f83011261173e57600080fd5b6113cf83833560208501611650565b801515811461088257600080fd5b80356117198161174d565b6000602080838503121561177957600080fd5b82356001600160401b038082111561179057600080fd5b818501915085601f8301126117a457600080fd5b8135818111156117b6576117b66115e1565b8060051b6117c5858201611620565b91825283810185019185810190898411156117df57600080fd5b86860192505b838310156118c9578235858111156117fc57600080fd5b8601610100818c03601f190181131561181457600080fd5b61181c6115f7565b89830135815261182e6040840161170a565b8a82015261183e6060840161171e565b60408201526080830135888111156118565760008081fd5b6118648e8c8387010161172d565b60608301525060a0808401358981111561187e5760008081fd5b61188c8f8d8388010161172d565b60808401525060c0808501358284015260e091508185013581840152506118b483850161175b565b908201528452505091860191908601906117e5565b9998505050505050505050565b6020815260006113cf6020830184611479565b6000602082840312156118fb57600080fd5b81356113cf81611593565b6000806040838503121561191957600080fd5b50508035926020909101359150565b60006020828403121561193a57600080fd5b81516113cf81611593565b600181811c9082168061195957607f821691505b60208210810361197957634e487b7160e01b600052602260045260246000fd5b50919050565b60006020828403121561199157600080fd5b5051919050565b634e487b7160e01b600052603260045260246000fd5b6000806000606084860312156119c357600080fd5b83516119ce8161174d565b602085015190935063ffffffff811681146119e857600080fd5b60408501519092506119f98161174d565b809150509250925092565b601f821115610bac576000816000526020600020601f850160051c81016020861015611a2d5750805b601f850160051c820191505b81811015611a4c57828155600101611a39565b505050505050565b81516001600160401b03811115611a6d57611a6d6115e1565b611a8181611a7b8454611945565b84611a04565b602080601f831160018114611ab65760008415611a9e5750858301515b600019600386901b1c1916600185901b178555611a4c565b600085815260208120601f198616915b82811015611ae557888601518255948401946001909101908401611ac6565b5085821015611b035787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b8181038181111561030d57634e487b7160e01b600052601160045260246000fd5b6000610180898352886020840152806040840152611b5481840189611479565b90508281036060840152611b688188611479565b90508281036080840152611b7c8187611479565b905082810360a0840152611b908186611479565b9150508251600d8110611ba557611ba561142b565b60c0830152602083015160e083015260408301516001600160a01b039081166101008401526060840151811661012084015260808401511661014083015260a090920151610160909101529695505050505050565b60008251611c0c818460208701611455565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbca2646970667358221220aee5fbb256314545dd773e3803222ae175d87c2f2c31e503f02c8115db4cd1ae64736f6c63430008180033",
}

// ParticipantStorageReplicaV1 is an auto generated Go binding around an Ethereum contract.
type ParticipantStorageReplicaV1 struct {
	abi abi.ABI
}

// NewParticipantStorageReplicaV1 creates a new instance of ParticipantStorageReplicaV1.
func NewParticipantStorageReplicaV1() *ParticipantStorageReplicaV1 {
	parsed, err := ParticipantStorageReplicaV1MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ParticipantStorageReplicaV1{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ParticipantStorageReplicaV1) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackUPGRADEINTERFACEVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) PackUPGRADEINTERFACEVERSION() []byte {
	enc, err := participantStorageReplicaV1.abi.Pack("UPGRADE_INTERFACE_VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUPGRADEINTERFACEVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackUPGRADEINTERFACEVERSION(data []byte) (string, error) {
	out, err := participantStorageReplicaV1.abi.Unpack("UPGRADE_INTERFACE_VERSION", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackAddOrUpdateParticipants is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x899cb662.
//
// Solidity: function addOrUpdateParticipants((uint256,uint8,uint8,string,string,uint256,uint256,bool)[] _participants) returns()
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) PackAddOrUpdateParticipants(participants []ParticipantStructsParticipant) []byte {
	enc, err := participantStorageReplicaV1.abi.Pack("addOrUpdateParticipants", participants)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) PackAuthority() []byte {
	enc, err := participantStorageReplicaV1.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := participantStorageReplicaV1.abi.Unpack("authority", data)
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
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) PackContractVersion() []byte {
	enc, err := participantStorageReplicaV1.abi.Pack("contractVersion")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackContractVersion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa0a8e460.
//
// Solidity: function contractVersion() pure returns(uint256)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackContractVersion(data []byte) (*big.Int, error) {
	out, err := participantStorageReplicaV1.abi.Unpack("contractVersion", data)
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
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) PackEndpoint() []byte {
	enc, err := participantStorageReplicaV1.abi.Pack("endpoint")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackEndpoint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackEndpoint(data []byte) (common.Address, error) {
	out, err := participantStorageReplicaV1.abi.Unpack("endpoint", data)
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
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) PackGetAddressByResourceId(resourceId [32]byte) []byte {
	enc, err := participantStorageReplicaV1.abi.Pack("getAddressByResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAddressByResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackGetAddressByResourceId(data []byte) (common.Address, error) {
	out, err := participantStorageReplicaV1.abi.Unpack("getAddressByResourceId", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetAllParticipants is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x195ec9ee.
//
// Solidity: function getAllParticipants() view returns((uint256,uint8,uint8,string,string,uint256,uint256,bool)[])
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) PackGetAllParticipants() []byte {
	enc, err := participantStorageReplicaV1.abi.Pack("getAllParticipants")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAllParticipants is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x195ec9ee.
//
// Solidity: function getAllParticipants() view returns((uint256,uint8,uint8,string,string,uint256,uint256,bool)[])
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackGetAllParticipants(data []byte) ([]ParticipantStructsParticipant, error) {
	out, err := participantStorageReplicaV1.abi.Unpack("getAllParticipants", data)
	if err != nil {
		return *new([]ParticipantStructsParticipant), err
	}
	out0 := *abi.ConvertType(out[0], new([]ParticipantStructsParticipant)).(*[]ParticipantStructsParticipant)
	return out0, err
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x485cc955.
//
// Solidity: function initialize(address _endpoint, address authority_) returns()
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) PackInitialize(endpoint common.Address, authority common.Address) []byte {
	enc, err := participantStorageReplicaV1.abi.Pack("initialize", endpoint, authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackInitialize0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc4d66de8.
//
// Solidity: function initialize(address _endpoint) returns()
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) PackInitialize0(endpoint common.Address) []byte {
	enc, err := participantStorageReplicaV1.abi.Pack("initialize0", endpoint)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackProxiableUUID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) PackProxiableUUID() []byte {
	enc, err := participantStorageReplicaV1.abi.Pack("proxiableUUID")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackProxiableUUID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackProxiableUUID(data []byte) ([32]byte, error) {
	out, err := participantStorageReplicaV1.abi.Unpack("proxiableUUID", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackRequestAllParticipantsDataFromPrivateHub is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x94b0fd87.
//
// Solidity: function requestAllParticipantsDataFromPrivateHub() returns()
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) PackRequestAllParticipantsDataFromPrivateHub() []byte {
	enc, err := participantStorageReplicaV1.abi.Pack("requestAllParticipantsDataFromPrivateHub")
	if err != nil {
		panic(err)
	}
	return enc
}

// PackResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) PackResourceId() []byte {
	enc, err := participantStorageReplicaV1.abi.Pack("resourceId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackResourceId(data []byte) ([32]byte, error) {
	out, err := participantStorageReplicaV1.abi.Unpack("resourceId", data)
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
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) PackUpgradeToAndCall(newImplementation common.Address, data []byte) []byte {
	enc, err := participantStorageReplicaV1.abi.Pack("upgradeToAndCall", newImplementation, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackValidateMessageParticipants is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd5c3614f.
//
// Solidity: function validateMessageParticipants(uint256 originChainId, uint256 destinationChainId) view returns()
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) PackValidateMessageParticipants(originChainId *big.Int, destinationChainId *big.Int) []byte {
	enc, err := participantStorageReplicaV1.abi.Pack("validateMessageParticipants", originChainId, destinationChainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackValidateParticipantStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc9557f56.
//
// Solidity: function validateParticipantStatus(uint256 chainId) view returns()
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) PackValidateParticipantStatus(chainId *big.Int) []byte {
	enc, err := participantStorageReplicaV1.abi.Pack("validateParticipantStatus", chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// ParticipantStorageReplicaV1AuthorityUpdated represents a AuthorityUpdated event raised by the ParticipantStorageReplicaV1 contract.
type ParticipantStorageReplicaV1AuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const ParticipantStorageReplicaV1AuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (ParticipantStorageReplicaV1AuthorityUpdated) ContractEventName() string {
	return ParticipantStorageReplicaV1AuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackAuthorityUpdatedEvent(log *types.Log) (*ParticipantStorageReplicaV1AuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != participantStorageReplicaV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ParticipantStorageReplicaV1AuthorityUpdated)
	if len(log.Data) > 0 {
		if err := participantStorageReplicaV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range participantStorageReplicaV1.abi.Events[event].Inputs {
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

// ParticipantStorageReplicaV1Initialized represents a Initialized event raised by the ParticipantStorageReplicaV1 contract.
type ParticipantStorageReplicaV1Initialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const ParticipantStorageReplicaV1InitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (ParticipantStorageReplicaV1Initialized) ContractEventName() string {
	return ParticipantStorageReplicaV1InitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackInitializedEvent(log *types.Log) (*ParticipantStorageReplicaV1Initialized, error) {
	event := "Initialized"
	if log.Topics[0] != participantStorageReplicaV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ParticipantStorageReplicaV1Initialized)
	if len(log.Data) > 0 {
		if err := participantStorageReplicaV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range participantStorageReplicaV1.abi.Events[event].Inputs {
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

// ParticipantStorageReplicaV1Upgraded represents a Upgraded event raised by the ParticipantStorageReplicaV1 contract.
type ParticipantStorageReplicaV1Upgraded struct {
	Implementation common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const ParticipantStorageReplicaV1UpgradedEventName = "Upgraded"

// ContractEventName returns the user-defined event name.
func (ParticipantStorageReplicaV1Upgraded) ContractEventName() string {
	return ParticipantStorageReplicaV1UpgradedEventName
}

// UnpackUpgradedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Upgraded(address indexed implementation)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackUpgradedEvent(log *types.Log) (*ParticipantStorageReplicaV1Upgraded, error) {
	event := "Upgraded"
	if log.Topics[0] != participantStorageReplicaV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ParticipantStorageReplicaV1Upgraded)
	if len(log.Data) > 0 {
		if err := participantStorageReplicaV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range participantStorageReplicaV1.abi.Events[event].Inputs {
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
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], participantStorageReplicaV1.abi.Errors["AddressEmptyCode"].ID.Bytes()[:4]) {
		return participantStorageReplicaV1.UnpackAddressEmptyCodeError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageReplicaV1.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return participantStorageReplicaV1.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageReplicaV1.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return participantStorageReplicaV1.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageReplicaV1.abi.Errors["FailedCall"].ID.Bytes()[:4]) {
		return participantStorageReplicaV1.UnpackFailedCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageReplicaV1.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return participantStorageReplicaV1.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageReplicaV1.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return participantStorageReplicaV1.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageReplicaV1.abi.Errors["ParticipantStorageReplicaV1NotFromPrivateHub"].ID.Bytes()[:4]) {
		return participantStorageReplicaV1.UnpackParticipantStorageReplicaV1NotFromPrivateHubError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageReplicaV1.abi.Errors["ParticipantStorageReplicaV1ParticipantNotActive"].ID.Bytes()[:4]) {
		return participantStorageReplicaV1.UnpackParticipantStorageReplicaV1ParticipantNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageReplicaV1.abi.Errors["ParticipantStorageReplicaV1ParticipantNotRegistered"].ID.Bytes()[:4]) {
		return participantStorageReplicaV1.UnpackParticipantStorageReplicaV1ParticipantNotRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageReplicaV1.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return participantStorageReplicaV1.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageReplicaV1.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return participantStorageReplicaV1.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageReplicaV1.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return participantStorageReplicaV1.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageReplicaV1.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return participantStorageReplicaV1.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageReplicaV1.abi.Errors["UUPSUnauthorizedCallContext"].ID.Bytes()[:4]) {
		return participantStorageReplicaV1.UnpackUUPSUnauthorizedCallContextError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageReplicaV1.abi.Errors["UUPSUnsupportedProxiableUUID"].ID.Bytes()[:4]) {
		return participantStorageReplicaV1.UnpackUUPSUnsupportedProxiableUUIDError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ParticipantStorageReplicaV1AddressEmptyCode represents a AddressEmptyCode error raised by the ParticipantStorageReplicaV1 contract.
type ParticipantStorageReplicaV1AddressEmptyCode struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressEmptyCode(address target)
func ParticipantStorageReplicaV1AddressEmptyCodeErrorID() common.Hash {
	return common.HexToHash("0x9996b315c842ff135b8fc4a08ad5df1c344efbc03d2687aecc0678050d2aac89")
}

// UnpackAddressEmptyCodeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressEmptyCode(address target)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackAddressEmptyCodeError(raw []byte) (*ParticipantStorageReplicaV1AddressEmptyCode, error) {
	out := new(ParticipantStorageReplicaV1AddressEmptyCode)
	if err := participantStorageReplicaV1.abi.UnpackIntoInterface(out, "AddressEmptyCode", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageReplicaV1ERC1967InvalidImplementation represents a ERC1967InvalidImplementation error raised by the ParticipantStorageReplicaV1 contract.
type ParticipantStorageReplicaV1ERC1967InvalidImplementation struct {
	Implementation common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func ParticipantStorageReplicaV1ERC1967InvalidImplementationErrorID() common.Hash {
	return common.HexToHash("0x4c9c8ce3ceb3130f17f7cdba48d89b5b0129f266a8bac114e6e315a41879b617")
}

// UnpackERC1967InvalidImplementationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackERC1967InvalidImplementationError(raw []byte) (*ParticipantStorageReplicaV1ERC1967InvalidImplementation, error) {
	out := new(ParticipantStorageReplicaV1ERC1967InvalidImplementation)
	if err := participantStorageReplicaV1.abi.UnpackIntoInterface(out, "ERC1967InvalidImplementation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageReplicaV1ERC1967NonPayable represents a ERC1967NonPayable error raised by the ParticipantStorageReplicaV1 contract.
type ParticipantStorageReplicaV1ERC1967NonPayable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967NonPayable()
func ParticipantStorageReplicaV1ERC1967NonPayableErrorID() common.Hash {
	return common.HexToHash("0xb398979fa84f543c8e222f17890372c487baf85e062276c127fef521eea7224b")
}

// UnpackERC1967NonPayableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967NonPayable()
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackERC1967NonPayableError(raw []byte) (*ParticipantStorageReplicaV1ERC1967NonPayable, error) {
	out := new(ParticipantStorageReplicaV1ERC1967NonPayable)
	if err := participantStorageReplicaV1.abi.UnpackIntoInterface(out, "ERC1967NonPayable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageReplicaV1FailedCall represents a FailedCall error raised by the ParticipantStorageReplicaV1 contract.
type ParticipantStorageReplicaV1FailedCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedCall()
func ParticipantStorageReplicaV1FailedCallErrorID() common.Hash {
	return common.HexToHash("0xd6bda27508c0fb6d8a39b4b122878dab26f731a7d4e4abe711dd3731899052a4")
}

// UnpackFailedCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedCall()
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackFailedCallError(raw []byte) (*ParticipantStorageReplicaV1FailedCall, error) {
	out := new(ParticipantStorageReplicaV1FailedCall)
	if err := participantStorageReplicaV1.abi.UnpackIntoInterface(out, "FailedCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageReplicaV1InvalidInitialization represents a InvalidInitialization error raised by the ParticipantStorageReplicaV1 contract.
type ParticipantStorageReplicaV1InvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func ParticipantStorageReplicaV1InvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackInvalidInitializationError(raw []byte) (*ParticipantStorageReplicaV1InvalidInitialization, error) {
	out := new(ParticipantStorageReplicaV1InvalidInitialization)
	if err := participantStorageReplicaV1.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageReplicaV1NotInitializing represents a NotInitializing error raised by the ParticipantStorageReplicaV1 contract.
type ParticipantStorageReplicaV1NotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func ParticipantStorageReplicaV1NotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackNotInitializingError(raw []byte) (*ParticipantStorageReplicaV1NotInitializing, error) {
	out := new(ParticipantStorageReplicaV1NotInitializing)
	if err := participantStorageReplicaV1.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageReplicaV1ParticipantStorageReplicaV1NotFromPrivateHub represents a ParticipantStorageReplicaV1__NotFromPrivateHub error raised by the ParticipantStorageReplicaV1 contract.
type ParticipantStorageReplicaV1ParticipantStorageReplicaV1NotFromPrivateHub struct {
	FromChainId *big.Int
	HubId       *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ParticipantStorageReplicaV1__NotFromPrivateHub(uint256 fromChainId, uint256 hubId)
func ParticipantStorageReplicaV1ParticipantStorageReplicaV1NotFromPrivateHubErrorID() common.Hash {
	return common.HexToHash("0x39d0df3a9a65531410c339b65b4e3ecb9c36c1653ebf78c2d64dcf85d846b19d")
}

// UnpackParticipantStorageReplicaV1NotFromPrivateHubError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ParticipantStorageReplicaV1__NotFromPrivateHub(uint256 fromChainId, uint256 hubId)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackParticipantStorageReplicaV1NotFromPrivateHubError(raw []byte) (*ParticipantStorageReplicaV1ParticipantStorageReplicaV1NotFromPrivateHub, error) {
	out := new(ParticipantStorageReplicaV1ParticipantStorageReplicaV1NotFromPrivateHub)
	if err := participantStorageReplicaV1.abi.UnpackIntoInterface(out, "ParticipantStorageReplicaV1NotFromPrivateHub", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageReplicaV1ParticipantStorageReplicaV1ParticipantNotActive represents a ParticipantStorageReplicaV1__ParticipantNotActive error raised by the ParticipantStorageReplicaV1 contract.
type ParticipantStorageReplicaV1ParticipantStorageReplicaV1ParticipantNotActive struct {
	ChainId *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ParticipantStorageReplicaV1__ParticipantNotActive(uint256 chainId)
func ParticipantStorageReplicaV1ParticipantStorageReplicaV1ParticipantNotActiveErrorID() common.Hash {
	return common.HexToHash("0x31d464a875bd0eeda0fe49ee2737fedec9762984a5039a701f0c59c1e94a38dc")
}

// UnpackParticipantStorageReplicaV1ParticipantNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ParticipantStorageReplicaV1__ParticipantNotActive(uint256 chainId)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackParticipantStorageReplicaV1ParticipantNotActiveError(raw []byte) (*ParticipantStorageReplicaV1ParticipantStorageReplicaV1ParticipantNotActive, error) {
	out := new(ParticipantStorageReplicaV1ParticipantStorageReplicaV1ParticipantNotActive)
	if err := participantStorageReplicaV1.abi.UnpackIntoInterface(out, "ParticipantStorageReplicaV1ParticipantNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageReplicaV1ParticipantStorageReplicaV1ParticipantNotRegistered represents a ParticipantStorageReplicaV1__ParticipantNotRegistered error raised by the ParticipantStorageReplicaV1 contract.
type ParticipantStorageReplicaV1ParticipantStorageReplicaV1ParticipantNotRegistered struct {
	ChainId *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ParticipantStorageReplicaV1__ParticipantNotRegistered(uint256 chainId)
func ParticipantStorageReplicaV1ParticipantStorageReplicaV1ParticipantNotRegisteredErrorID() common.Hash {
	return common.HexToHash("0x914ee7d5bb8142ae3c0a33ee6f1daf8d227c786c33fa81319bc2da6bfcd18a2b")
}

// UnpackParticipantStorageReplicaV1ParticipantNotRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ParticipantStorageReplicaV1__ParticipantNotRegistered(uint256 chainId)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackParticipantStorageReplicaV1ParticipantNotRegisteredError(raw []byte) (*ParticipantStorageReplicaV1ParticipantStorageReplicaV1ParticipantNotRegistered, error) {
	out := new(ParticipantStorageReplicaV1ParticipantStorageReplicaV1ParticipantNotRegistered)
	if err := participantStorageReplicaV1.abi.UnpackIntoInterface(out, "ParticipantStorageReplicaV1ParticipantNotRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageReplicaV1RaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the ParticipantStorageReplicaV1 contract.
type ParticipantStorageReplicaV1RaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func ParticipantStorageReplicaV1RaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*ParticipantStorageReplicaV1RaylsAccessManagedContractPaused, error) {
	out := new(ParticipantStorageReplicaV1RaylsAccessManagedContractPaused)
	if err := participantStorageReplicaV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageReplicaV1RaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the ParticipantStorageReplicaV1 contract.
type ParticipantStorageReplicaV1RaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func ParticipantStorageReplicaV1RaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*ParticipantStorageReplicaV1RaylsAccessManagedInvalidAuthority, error) {
	out := new(ParticipantStorageReplicaV1RaylsAccessManagedInvalidAuthority)
	if err := participantStorageReplicaV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageReplicaV1RaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the ParticipantStorageReplicaV1 contract.
type ParticipantStorageReplicaV1RaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func ParticipantStorageReplicaV1RaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*ParticipantStorageReplicaV1RaylsAccessManagedMustSchedule, error) {
	out := new(ParticipantStorageReplicaV1RaylsAccessManagedMustSchedule)
	if err := participantStorageReplicaV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageReplicaV1RaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the ParticipantStorageReplicaV1 contract.
type ParticipantStorageReplicaV1RaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func ParticipantStorageReplicaV1RaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*ParticipantStorageReplicaV1RaylsAccessManagedUnauthorized, error) {
	out := new(ParticipantStorageReplicaV1RaylsAccessManagedUnauthorized)
	if err := participantStorageReplicaV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageReplicaV1UUPSUnauthorizedCallContext represents a UUPSUnauthorizedCallContext error raised by the ParticipantStorageReplicaV1 contract.
type ParticipantStorageReplicaV1UUPSUnauthorizedCallContext struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnauthorizedCallContext()
func ParticipantStorageReplicaV1UUPSUnauthorizedCallContextErrorID() common.Hash {
	return common.HexToHash("0xe07c8dba242a06571ac65fe4bbe20522c9fb111cb33599b799ff8039c1ed18f4")
}

// UnpackUUPSUnauthorizedCallContextError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnauthorizedCallContext()
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackUUPSUnauthorizedCallContextError(raw []byte) (*ParticipantStorageReplicaV1UUPSUnauthorizedCallContext, error) {
	out := new(ParticipantStorageReplicaV1UUPSUnauthorizedCallContext)
	if err := participantStorageReplicaV1.abi.UnpackIntoInterface(out, "UUPSUnauthorizedCallContext", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageReplicaV1UUPSUnsupportedProxiableUUID represents a UUPSUnsupportedProxiableUUID error raised by the ParticipantStorageReplicaV1 contract.
type ParticipantStorageReplicaV1UUPSUnsupportedProxiableUUID struct {
	Slot [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func ParticipantStorageReplicaV1UUPSUnsupportedProxiableUUIDErrorID() common.Hash {
	return common.HexToHash("0xaa1d49a4c084bfa9aeeee2a0be65267a7f19ba7e1476b114dac513d2c14cb563")
}

// UnpackUUPSUnsupportedProxiableUUIDError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func (participantStorageReplicaV1 *ParticipantStorageReplicaV1) UnpackUUPSUnsupportedProxiableUUIDError(raw []byte) (*ParticipantStorageReplicaV1UUPSUnsupportedProxiableUUID, error) {
	out := new(ParticipantStorageReplicaV1UUPSUnsupportedProxiableUUID)
	if err := participantStorageReplicaV1.abi.UnpackIntoInterface(out, "UUPSUnsupportedProxiableUUID", raw); err != nil {
		return nil, err
	}
	return out, nil
}
