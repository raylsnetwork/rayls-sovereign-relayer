// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package PublicChainERC20

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

// PublicChainERC20MetaData contains all meta data concerning the PublicChainERC20 contract.
var PublicChainERC20MetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_raylsNodeEndpoint\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_initialSupply\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"privateAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPublicRaylsNodeEndpoint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"privateChainAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"receiveTeleportFromPrivacyNode\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"srcChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revertTeleportToPrivacyNode\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"teleportToPrivacyNode\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RaylsPublicErc20TokenCreated\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsPublicERC20Handler__AmountMustBeGreaterThanZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsPublicERC20Handler__DestinationIsZeroAddress\",\"inputs\":[]}]",
	ID:  "PublicChainERC20",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001a2e38038062001a2e833981016040819052620000349162000760565b600080546001600160a01b0319166001600160a01b03851690811782556040805163bf7e214f60e01b815290518893889388933393889388938893889363bf7e214f916004808201926020929091908290030181865afa1580156200009d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000c39190620007f9565b90506001600160a01b03811615620000e057620000e0816200019a565b5060049050620000f18382620008af565b506005620001008282620008af565b50600691506200011390508682620008af565b506007620001228582620008af565b50600880546001600160a01b0319166001600160a01b03831617905562000149826200023b565b60405130907f46a23e472b0a659f24e62ee736c903350af82814eb3501a323d7c9cd1eb366be90600090a2505050505060008211156200018f576200018f338362000485565b505050505062000b0a565b6001600160a01b038116620001d257604051638944034760e01b81526001600160a01b03821660048201526024015b60405180910390fd5b6000620001de620004c3565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b60006200024762000528565b90506001600160a01b0381166200025c575050565b6040805160028082526060820183526000926020830190803683370190505090506340c10f1960e01b816000815181106200029b576200029b6200097b565b6001600160e01b0319909216602092830291909101909101528051632770a7eb60e21b9082906001908110620002d557620002d56200097b565b6001600160e01b03199290921660209283029190910190910152604080516001808252818301909252600091816020015b604080518082019091526060808252602082015281526020019060019003908162000306575050604080516002808252606082018352929350600092909160208301908036833701905050905063bef97c9e60e01b816000815181106200037157620003716200097b565b6001600160e01b0319909216602092830291909101909101528051632217f5f160e21b9082906001908110620003ab57620003ab6200097b565b6001600160e01b03199092166020928302919091018201526040805160808101825260109181019182526f26a2a9a9a0a3a2afa2ac22a1aaaa27a960811b6060820152908152908101829052825183906000906200040d576200040d6200097b565b60209081029190910101526040516337af400760e11b81526001600160a01b03851690636f5e800e906200044a9088908790879060040162000991565b600060405180830381600087803b1580156200046557600080fd5b505af11580156200047a573d6000803e3d6000fd5b505050505050505050565b6001600160a01b038216620004b15760405163ec442f0560e01b815260006004820152602401620001c9565b620004bf6000838362000543565b5050565b60008060ff19620004f660017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f3562000ad8565b6040516020016200050991815260200190565b60408051601f1981840301815291905280516020909101201692915050565b600062000534620004c3565b546001600160a01b0316919050565b6001600160a01b0383166200057257806003600082825462000566919062000af4565b90915550620005e69050565b6001600160a01b03831660009081526001602052604090205481811015620005c75760405163391434e360e21b81526001600160a01b03851660048201526024810182905260448101839052606401620001c9565b6001600160a01b03841660009081526001602052604090209082900390555b6001600160a01b038216620006045760038054829003905562000623565b6001600160a01b03821660009081526001602052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516200066991815260200190565b60405180910390a3505050565b634e487b7160e01b600052604160045260246000fd5b60005b83811015620006a95781810151838201526020016200068f565b50506000910152565b600082601f830112620006c457600080fd5b81516001600160401b0380821115620006e157620006e162000676565b604051601f8301601f19908116603f011681019082821181831017156200070c576200070c62000676565b816040528381528660208588010111156200072657600080fd5b620007398460208301602089016200068c565b9695505050505050565b80516001600160a01b03811681146200075b57600080fd5b919050565b600080600080600060a086880312156200077957600080fd5b85516001600160401b03808211156200079157600080fd5b6200079f89838a01620006b2565b96506020880151915080821115620007b657600080fd5b50620007c588828901620006b2565b945050620007d66040870162000743565b925060608601519150620007ed6080870162000743565b90509295509295909350565b6000602082840312156200080c57600080fd5b620008178262000743565b9392505050565b600181811c908216806200083357607f821691505b6020821081036200085457634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620008aa576000816000526020600020601f850160051c81016020861015620008855750805b601f850160051c820191505b81811015620008a65782815560010162000891565b5050505b505050565b81516001600160401b03811115620008cb57620008cb62000676565b620008e381620008dc84546200081e565b846200085a565b602080601f8311600181146200091b5760008415620009025750858301515b600019600386901b1c1916600185901b178555620008a6565b600085815260208120601f198616915b828110156200094c578886015182559484019460019091019084016200092b565b50858210156200096b5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b634e487b7160e01b600052603260045260246000fd5b6001600160a01b03841681526060602080830182905284518383018190526000929160809182860190888301865b82811015620009e75781516001600160e01b03191684529284019290840190600101620009bf565b505050604086820360408801528188518084528484019150848160051b850101858b0160005b8381101562000aaf57601f1980888503018652825180518886528051808a88015262000a3f818f89018e85016200068c565b918b0151601f9290920190921685018581038d018b87015281518d8201818152928c019350600092918d01905b8084101562000a985784516001600160e01b0319168252938c019360019390930192908c019062000a6c565b50978b019795505050918801915060010162000a0d565b50909d9c50505050505050505050505050565b634e487b7160e01b600052601160045260246000fd5b8181038181111562000aee5762000aee62000ac2565b92915050565b8082018082111562000aee5762000aee62000ac2565b610f148062000b1a6000396000f3fe608060405234801561001057600080fd5b50600436106100eb5760003560e01c806370a082311161009257806370a08231146101a757806380a0a8b3146101d0578063885fd7c41461017857806395d89b41146101e15780639dc29fac146101e9578063a9059cbb146101fc578063bef97c9e1461020f578063bf7e214f14610222578063dd62ed3e1461022a57600080fd5b806306fdde03146100f0578063095ea7b31461010e57806314eb966d1461013157806318160ddd1461014457806323b872dd14610156578063313ce5671461016957806340c10f19146101785780635f2e03551461018d575b600080fd5b6100f861023d565b6040516101059190610b8b565b60405180910390f35b61012161011c366004610bc1565b6102cf565b6040519015158152602001610105565b61012161013f366004610beb565b6102e9565b6003545b604051908152602001610105565b610121610164366004610c1e565b610475565b60405160128152602001610105565b61018b610186366004610bc1565b610499565b005b6000546001600160a01b03165b6040516101059190610c5a565b6101486101b5366004610c6e565b6001600160a01b031660009081526001602052604090205490565b6008546001600160a01b031661019a565b6100f86104bd565b61018b6101f7366004610bc1565b6104cc565b61012161020a366004610bc1565b6104ec565b61018b61021d366004610c89565b6104fa565b61019a610634565b610148610238366004610ccd565b61064d565b60606006805461024c90610d00565b80601f016020809104026020016040519081016040528092919081815260200182805461027890610d00565b80156102c55780601f1061029a576101008083540402835291602001916102c5565b820191906000526020600020905b8154815290600101906020018083116102a857829003601f168201915b5050505050905090565b6000336102dd818585610678565b60019150505b92915050565b60006001600160a01b038416610312576040516378682caf60e11b815260040160405180910390fd5b8260000361033357604051639e89fced60e01b815260040160405180910390fd5b61033d338461068a565b6040805160c0810182526001815260006020820181905233828401526001600160a01b03878116606084015230608084015260a083018790529054600854935192939082169263e0d9d846928792169061039d908a908a90602401610d3a565b60408051601f198184030181529181526020820180516001600160e01b03166330bb92bf60e01b179052516103d89033908b90602401610d3a565b60408051601f198184030181529181526020820180516001600160e01b0316632217f5f160e21b1790525160e086901b6001600160e01b031916815261042694939291908890600401610d53565b6020604051808303816000875af1158015610445573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104699190610e0a565b50600195945050505050565b6000336104838582856106c9565b61048e858585610717565b506001949350505050565b6104af336000356001600160e01b031916610776565b6104b982826108c1565b5050565b60606007805461024c90610d00565b6104e2336000356001600160e01b031916610776565b6104b9828261068a565b6000336102dd818585610717565b610510336000356001600160e01b031916610776565b6001600160a01b038216610624576040805160c081018252600181526000602082018190526001600160a01b03878116838501819052606084015230608084015260a083018590529054600854935192939082169263e0d9d846928892169061057f908a908890602401610d3a565b60408051601f19818403018152918152602080830180516001600160e01b0316633a03ff1f60e01b1790528151908101825260008152905160e086901b6001600160e01b03191681526105da94939291908890600401610d53565b6020604051808303816000875af11580156105f9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061061d9190610e0a565b505061062e565b61062e82826108c1565b50505050565b600061063e6108f7565b546001600160a01b0316919050565b6001600160a01b03918216600090815260026020908152604080832093909416825291909152205490565b6106858383836001610959565b505050565b6001600160a01b0382166106bd576000604051634b637e8f60e11b81526004016106b49190610c5a565b60405180910390fd5b6104b982600083610a2e565b60006106d5848461064d565b905060001981101561062e578181101561070857828183604051637dc7a0d960e11b81526004016106b493929190610e23565b61062e84848484036000610959565b6001600160a01b038316610741576000604051634b637e8f60e11b81526004016106b49190610c5a565b6001600160a01b03821661076b57600060405163ec442f0560e01b81526004016106b49190610c5a565b610685838383610a2e565b60006107806108f7565b80549091506001600160a01b0316806107af576000604051638944034760e01b81526004016106b49190610c5a565b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa158015610813573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108379190610e54565b925092509250826108b85780156108615760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff82161561089d5760405163a426878960e01b81526001600160a01b038816600482015263ffffffff831660248201526044016106b4565b86604051632ecd3d0360e21b81526004016106b49190610c5a565b50505050505050565b6001600160a01b0382166108eb57600060405163ec442f0560e01b81526004016106b49190610c5a565b6104b960008383610a2e565b60008060ff1961092860017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f35610eb8565b60405160200161093a91815260200190565b60408051601f1981840301815291905280516020909101201692915050565b6001600160a01b03841661098357600060405163e602df0560e01b81526004016106b49190610c5a565b6001600160a01b0383166109ad576000604051634a1406b160e11b81526004016106b49190610c5a565b6001600160a01b038085166000908152600260209081526040808320938716835292905220829055801561062e57826001600160a01b0316846001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92584604051610a2091815260200190565b60405180910390a350505050565b6001600160a01b038316610a59578060036000828254610a4e9190610ecb565b90915550610ab89050565b6001600160a01b03831660009081526001602052604090205481811015610a995783818360405163391434e360e21b81526004016106b493929190610e23565b6001600160a01b03841660009081526001602052604090209082900390555b6001600160a01b038216610ad457600380548290039055610af3565b6001600160a01b03821660009081526001602052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051610b3891815260200190565b60405180910390a3505050565b6000815180845260005b81811015610b6b57602081850181015186830182015201610b4f565b506000602082860101526020601f19601f83011685010191505092915050565b602081526000610b9e6020830184610b45565b9392505050565b80356001600160a01b0381168114610bbc57600080fd5b919050565b60008060408385031215610bd457600080fd5b610bdd83610ba5565b946020939093013593505050565b600080600060608486031215610c0057600080fd5b610c0984610ba5565b95602085013595506040909401359392505050565b600080600060608486031215610c3357600080fd5b610c3c84610ba5565b9250610c4a60208501610ba5565b9150604084013590509250925092565b6001600160a01b0391909116815260200190565b600060208284031215610c8057600080fd5b610b9e82610ba5565b60008060008060808587031215610c9f57600080fd5b610ca885610ba5565b935060208501359250610cbd60408601610ba5565b9396929550929360600135925050565b60008060408385031215610ce057600080fd5b610ce983610ba5565b9150610cf760208401610ba5565b90509250929050565b600181811c90821680610d1457607f821691505b602082108103610d3457634e487b7160e01b600052602260045260246000fd5b50919050565b6001600160a01b03929092168252602082015260400190565b8581526001600160a01b03858116602083015261014060408301819052600091610d7f84830188610b45565b91508382036060850152610d938287610b45565b92508451915060078210610db757634e487b7160e01b600052602160045260246000fd5b608084810192909252602085015160a08501526040850151811660c085015260608501511660e08401528301516001600160a01b0381166101008401525060a08301516101208301529695505050505050565b600060208284031215610e1c57600080fd5b5051919050565b6001600160a01b039390931683526020830191909152604082015260600190565b80518015158114610bbc57600080fd5b600080600060608486031215610e6957600080fd5b610e7284610e44565b9250602084015163ffffffff81168114610e8b57600080fd5b9150610e9960408501610e44565b90509250925092565b634e487b7160e01b600052601160045260246000fd5b818103818111156102e3576102e3610ea2565b808201808211156102e3576102e3610ea256fea2646970667358221220b01e7b1965c8878dc6f9ca7fbc6283d4600b6bb469019155559137bda298a28064736f6c63430008180033",
}

// PublicChainERC20 is an auto generated Go binding around an Ethereum contract.
type PublicChainERC20 struct {
	abi abi.ABI
}

// NewPublicChainERC20 creates a new instance of PublicChainERC20.
func NewPublicChainERC20() *PublicChainERC20 {
	parsed, err := PublicChainERC20MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &PublicChainERC20{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *PublicChainERC20) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(string _name, string _symbol, address _raylsNodeEndpoint, uint256 _initialSupply, address privateAddress) returns()
func (publicChainERC20 *PublicChainERC20) PackConstructor(_name string, _symbol string, _raylsNodeEndpoint common.Address, _initialSupply *big.Int, privateAddress common.Address) []byte {
	enc, err := publicChainERC20.abi.Pack("", _name, _symbol, _raylsNodeEndpoint, _initialSupply, privateAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (publicChainERC20 *PublicChainERC20) PackAllowance(owner common.Address, spender common.Address) []byte {
	enc, err := publicChainERC20.abi.Pack("allowance", owner, spender)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAllowance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (publicChainERC20 *PublicChainERC20) UnpackAllowance(data []byte) (*big.Int, error) {
	out, err := publicChainERC20.abi.Unpack("allowance", data)
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
func (publicChainERC20 *PublicChainERC20) PackApprove(spender common.Address, value *big.Int) []byte {
	enc, err := publicChainERC20.abi.Pack("approve", spender, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackApprove is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (publicChainERC20 *PublicChainERC20) UnpackApprove(data []byte) (bool, error) {
	out, err := publicChainERC20.abi.Unpack("approve", data)
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
func (publicChainERC20 *PublicChainERC20) PackAuthority() []byte {
	enc, err := publicChainERC20.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (publicChainERC20 *PublicChainERC20) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := publicChainERC20.abi.Unpack("authority", data)
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
func (publicChainERC20 *PublicChainERC20) PackBalanceOf(account common.Address) []byte {
	enc, err := publicChainERC20.abi.Pack("balanceOf", account)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (publicChainERC20 *PublicChainERC20) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := publicChainERC20.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackBurn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9dc29fac.
//
// Solidity: function burn(address from, uint256 amount) returns()
func (publicChainERC20 *PublicChainERC20) PackBurn(from common.Address, amount *big.Int) []byte {
	enc, err := publicChainERC20.abi.Pack("burn", from, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (publicChainERC20 *PublicChainERC20) PackDecimals() []byte {
	enc, err := publicChainERC20.abi.Pack("decimals")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDecimals is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (publicChainERC20 *PublicChainERC20) UnpackDecimals(data []byte) (uint8, error) {
	out, err := publicChainERC20.abi.Unpack("decimals", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, err
}

// PackGetPublicRaylsNodeEndpoint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5f2e0355.
//
// Solidity: function getPublicRaylsNodeEndpoint() view returns(address)
func (publicChainERC20 *PublicChainERC20) PackGetPublicRaylsNodeEndpoint() []byte {
	enc, err := publicChainERC20.abi.Pack("getPublicRaylsNodeEndpoint")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetPublicRaylsNodeEndpoint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5f2e0355.
//
// Solidity: function getPublicRaylsNodeEndpoint() view returns(address)
func (publicChainERC20 *PublicChainERC20) UnpackGetPublicRaylsNodeEndpoint(data []byte) (common.Address, error) {
	out, err := publicChainERC20.abi.Unpack("getPublicRaylsNodeEndpoint", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (publicChainERC20 *PublicChainERC20) PackMint(to common.Address, amount *big.Int) []byte {
	enc, err := publicChainERC20.abi.Pack("mint", to, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (publicChainERC20 *PublicChainERC20) PackName() []byte {
	enc, err := publicChainERC20.abi.Pack("name")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (publicChainERC20 *PublicChainERC20) UnpackName(data []byte) (string, error) {
	out, err := publicChainERC20.abi.Unpack("name", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackPrivateChainAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x80a0a8b3.
//
// Solidity: function privateChainAddress() view returns(address)
func (publicChainERC20 *PublicChainERC20) PackPrivateChainAddress() []byte {
	enc, err := publicChainERC20.abi.Pack("privateChainAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackPrivateChainAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x80a0a8b3.
//
// Solidity: function privateChainAddress() view returns(address)
func (publicChainERC20 *PublicChainERC20) UnpackPrivateChainAddress(data []byte) (common.Address, error) {
	out, err := publicChainERC20.abi.Unpack("privateChainAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackReceiveTeleportFromPrivacyNode is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbef97c9e.
//
// Solidity: function receiveTeleportFromPrivacyNode(address from, uint256 srcChainId, address to, uint256 amount) returns()
func (publicChainERC20 *PublicChainERC20) PackReceiveTeleportFromPrivacyNode(from common.Address, srcChainId *big.Int, to common.Address, amount *big.Int) []byte {
	enc, err := publicChainERC20.abi.Pack("receiveTeleportFromPrivacyNode", from, srcChainId, to, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRevertTeleportToPrivacyNode is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x885fd7c4.
//
// Solidity: function revertTeleportToPrivacyNode(address to, uint256 amount) returns()
func (publicChainERC20 *PublicChainERC20) PackRevertTeleportToPrivacyNode(to common.Address, amount *big.Int) []byte {
	enc, err := publicChainERC20.abi.Pack("revertTeleportToPrivacyNode", to, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (publicChainERC20 *PublicChainERC20) PackSymbol() []byte {
	enc, err := publicChainERC20.abi.Pack("symbol")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (publicChainERC20 *PublicChainERC20) UnpackSymbol(data []byte) (string, error) {
	out, err := publicChainERC20.abi.Unpack("symbol", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackTeleportToPrivacyNode is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x14eb966d.
//
// Solidity: function teleportToPrivacyNode(address to, uint256 amount, uint256 chainId) returns(bool)
func (publicChainERC20 *PublicChainERC20) PackTeleportToPrivacyNode(to common.Address, amount *big.Int, chainId *big.Int) []byte {
	enc, err := publicChainERC20.abi.Pack("teleportToPrivacyNode", to, amount, chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTeleportToPrivacyNode is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x14eb966d.
//
// Solidity: function teleportToPrivacyNode(address to, uint256 amount, uint256 chainId) returns(bool)
func (publicChainERC20 *PublicChainERC20) UnpackTeleportToPrivacyNode(data []byte) (bool, error) {
	out, err := publicChainERC20.abi.Unpack("teleportToPrivacyNode", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (publicChainERC20 *PublicChainERC20) PackTotalSupply() []byte {
	enc, err := publicChainERC20.abi.Pack("totalSupply")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (publicChainERC20 *PublicChainERC20) UnpackTotalSupply(data []byte) (*big.Int, error) {
	out, err := publicChainERC20.abi.Unpack("totalSupply", data)
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
func (publicChainERC20 *PublicChainERC20) PackTransfer(to common.Address, value *big.Int) []byte {
	enc, err := publicChainERC20.abi.Pack("transfer", to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTransfer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (publicChainERC20 *PublicChainERC20) UnpackTransfer(data []byte) (bool, error) {
	out, err := publicChainERC20.abi.Unpack("transfer", data)
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
func (publicChainERC20 *PublicChainERC20) PackTransferFrom(from common.Address, to common.Address, value *big.Int) []byte {
	enc, err := publicChainERC20.abi.Pack("transferFrom", from, to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTransferFrom is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (publicChainERC20 *PublicChainERC20) UnpackTransferFrom(data []byte) (bool, error) {
	out, err := publicChainERC20.abi.Unpack("transferFrom", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PublicChainERC20Approval represents a Approval event raised by the PublicChainERC20 contract.
type PublicChainERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const PublicChainERC20ApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (PublicChainERC20Approval) ContractEventName() string {
	return PublicChainERC20ApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (publicChainERC20 *PublicChainERC20) UnpackApprovalEvent(log *types.Log) (*PublicChainERC20Approval, error) {
	event := "Approval"
	if log.Topics[0] != publicChainERC20.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PublicChainERC20Approval)
	if len(log.Data) > 0 {
		if err := publicChainERC20.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range publicChainERC20.abi.Events[event].Inputs {
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

// PublicChainERC20AuthorityUpdated represents a AuthorityUpdated event raised by the PublicChainERC20 contract.
type PublicChainERC20AuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const PublicChainERC20AuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (PublicChainERC20AuthorityUpdated) ContractEventName() string {
	return PublicChainERC20AuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (publicChainERC20 *PublicChainERC20) UnpackAuthorityUpdatedEvent(log *types.Log) (*PublicChainERC20AuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != publicChainERC20.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PublicChainERC20AuthorityUpdated)
	if len(log.Data) > 0 {
		if err := publicChainERC20.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range publicChainERC20.abi.Events[event].Inputs {
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

// PublicChainERC20RaylsPublicErc20TokenCreated represents a RaylsPublicErc20TokenCreated event raised by the PublicChainERC20 contract.
type PublicChainERC20RaylsPublicErc20TokenCreated struct {
	TokenAddress common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const PublicChainERC20RaylsPublicErc20TokenCreatedEventName = "RaylsPublicErc20TokenCreated"

// ContractEventName returns the user-defined event name.
func (PublicChainERC20RaylsPublicErc20TokenCreated) ContractEventName() string {
	return PublicChainERC20RaylsPublicErc20TokenCreatedEventName
}

// UnpackRaylsPublicErc20TokenCreatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RaylsPublicErc20TokenCreated(address indexed tokenAddress)
func (publicChainERC20 *PublicChainERC20) UnpackRaylsPublicErc20TokenCreatedEvent(log *types.Log) (*PublicChainERC20RaylsPublicErc20TokenCreated, error) {
	event := "RaylsPublicErc20TokenCreated"
	if log.Topics[0] != publicChainERC20.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PublicChainERC20RaylsPublicErc20TokenCreated)
	if len(log.Data) > 0 {
		if err := publicChainERC20.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range publicChainERC20.abi.Events[event].Inputs {
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

// PublicChainERC20Transfer represents a Transfer event raised by the PublicChainERC20 contract.
type PublicChainERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   *types.Log // Blockchain specific contextual infos
}

const PublicChainERC20TransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (PublicChainERC20Transfer) ContractEventName() string {
	return PublicChainERC20TransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (publicChainERC20 *PublicChainERC20) UnpackTransferEvent(log *types.Log) (*PublicChainERC20Transfer, error) {
	event := "Transfer"
	if log.Topics[0] != publicChainERC20.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PublicChainERC20Transfer)
	if len(log.Data) > 0 {
		if err := publicChainERC20.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range publicChainERC20.abi.Events[event].Inputs {
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
func (publicChainERC20 *PublicChainERC20) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], publicChainERC20.abi.Errors["ERC20InsufficientAllowance"].ID.Bytes()[:4]) {
		return publicChainERC20.UnpackERC20InsufficientAllowanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC20.abi.Errors["ERC20InsufficientBalance"].ID.Bytes()[:4]) {
		return publicChainERC20.UnpackERC20InsufficientBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC20.abi.Errors["ERC20InvalidApprover"].ID.Bytes()[:4]) {
		return publicChainERC20.UnpackERC20InvalidApproverError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC20.abi.Errors["ERC20InvalidReceiver"].ID.Bytes()[:4]) {
		return publicChainERC20.UnpackERC20InvalidReceiverError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC20.abi.Errors["ERC20InvalidSender"].ID.Bytes()[:4]) {
		return publicChainERC20.UnpackERC20InvalidSenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC20.abi.Errors["ERC20InvalidSpender"].ID.Bytes()[:4]) {
		return publicChainERC20.UnpackERC20InvalidSpenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC20.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return publicChainERC20.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC20.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return publicChainERC20.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC20.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return publicChainERC20.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC20.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return publicChainERC20.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC20.abi.Errors["RaylsPublicERC20HandlerAmountMustBeGreaterThanZero"].ID.Bytes()[:4]) {
		return publicChainERC20.UnpackRaylsPublicERC20HandlerAmountMustBeGreaterThanZeroError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC20.abi.Errors["RaylsPublicERC20HandlerDestinationIsZeroAddress"].ID.Bytes()[:4]) {
		return publicChainERC20.UnpackRaylsPublicERC20HandlerDestinationIsZeroAddressError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// PublicChainERC20ERC20InsufficientAllowance represents a ERC20InsufficientAllowance error raised by the PublicChainERC20 contract.
type PublicChainERC20ERC20InsufficientAllowance struct {
	Spender   common.Address
	Allowance *big.Int
	Needed    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func PublicChainERC20ERC20InsufficientAllowanceErrorID() common.Hash {
	return common.HexToHash("0xfb8f41b23e99d2101d86da76cdfa87dd51c82ed07d3cb62cbc473e469dbc75c3")
}

// UnpackERC20InsufficientAllowanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func (publicChainERC20 *PublicChainERC20) UnpackERC20InsufficientAllowanceError(raw []byte) (*PublicChainERC20ERC20InsufficientAllowance, error) {
	out := new(PublicChainERC20ERC20InsufficientAllowance)
	if err := publicChainERC20.abi.UnpackIntoInterface(out, "ERC20InsufficientAllowance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC20ERC20InsufficientBalance represents a ERC20InsufficientBalance error raised by the PublicChainERC20 contract.
type PublicChainERC20ERC20InsufficientBalance struct {
	Sender  common.Address
	Balance *big.Int
	Needed  *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func PublicChainERC20ERC20InsufficientBalanceErrorID() common.Hash {
	return common.HexToHash("0xe450d38cd8d9f7d95077d567d60ed49c7254716e6ad08fc9872816c97e0ffec6")
}

// UnpackERC20InsufficientBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func (publicChainERC20 *PublicChainERC20) UnpackERC20InsufficientBalanceError(raw []byte) (*PublicChainERC20ERC20InsufficientBalance, error) {
	out := new(PublicChainERC20ERC20InsufficientBalance)
	if err := publicChainERC20.abi.UnpackIntoInterface(out, "ERC20InsufficientBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC20ERC20InvalidApprover represents a ERC20InvalidApprover error raised by the PublicChainERC20 contract.
type PublicChainERC20ERC20InvalidApprover struct {
	Approver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidApprover(address approver)
func PublicChainERC20ERC20InvalidApproverErrorID() common.Hash {
	return common.HexToHash("0xe602df05cc75712490294c6c104ab7c17f4030363910a7a2626411c6d3118847")
}

// UnpackERC20InvalidApproverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidApprover(address approver)
func (publicChainERC20 *PublicChainERC20) UnpackERC20InvalidApproverError(raw []byte) (*PublicChainERC20ERC20InvalidApprover, error) {
	out := new(PublicChainERC20ERC20InvalidApprover)
	if err := publicChainERC20.abi.UnpackIntoInterface(out, "ERC20InvalidApprover", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC20ERC20InvalidReceiver represents a ERC20InvalidReceiver error raised by the PublicChainERC20 contract.
type PublicChainERC20ERC20InvalidReceiver struct {
	Receiver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func PublicChainERC20ERC20InvalidReceiverErrorID() common.Hash {
	return common.HexToHash("0xec442f055133b72f3b2f9f0bb351c406b178527de2040a7d1feb4e058771f613")
}

// UnpackERC20InvalidReceiverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func (publicChainERC20 *PublicChainERC20) UnpackERC20InvalidReceiverError(raw []byte) (*PublicChainERC20ERC20InvalidReceiver, error) {
	out := new(PublicChainERC20ERC20InvalidReceiver)
	if err := publicChainERC20.abi.UnpackIntoInterface(out, "ERC20InvalidReceiver", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC20ERC20InvalidSender represents a ERC20InvalidSender error raised by the PublicChainERC20 contract.
type PublicChainERC20ERC20InvalidSender struct {
	Sender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSender(address sender)
func PublicChainERC20ERC20InvalidSenderErrorID() common.Hash {
	return common.HexToHash("0x96c6fd1edd0cd6ef7ff0ecc0facdf53148dc0048b57fe58af65755250a7a96bd")
}

// UnpackERC20InvalidSenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSender(address sender)
func (publicChainERC20 *PublicChainERC20) UnpackERC20InvalidSenderError(raw []byte) (*PublicChainERC20ERC20InvalidSender, error) {
	out := new(PublicChainERC20ERC20InvalidSender)
	if err := publicChainERC20.abi.UnpackIntoInterface(out, "ERC20InvalidSender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC20ERC20InvalidSpender represents a ERC20InvalidSpender error raised by the PublicChainERC20 contract.
type PublicChainERC20ERC20InvalidSpender struct {
	Spender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSpender(address spender)
func PublicChainERC20ERC20InvalidSpenderErrorID() common.Hash {
	return common.HexToHash("0x94280d62c347d8d9f4d59a76ea321452406db88df38e0c9da304f58b57b373a2")
}

// UnpackERC20InvalidSpenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSpender(address spender)
func (publicChainERC20 *PublicChainERC20) UnpackERC20InvalidSpenderError(raw []byte) (*PublicChainERC20ERC20InvalidSpender, error) {
	out := new(PublicChainERC20ERC20InvalidSpender)
	if err := publicChainERC20.abi.UnpackIntoInterface(out, "ERC20InvalidSpender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC20RaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the PublicChainERC20 contract.
type PublicChainERC20RaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func PublicChainERC20RaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (publicChainERC20 *PublicChainERC20) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*PublicChainERC20RaylsAccessManagedContractPaused, error) {
	out := new(PublicChainERC20RaylsAccessManagedContractPaused)
	if err := publicChainERC20.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC20RaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the PublicChainERC20 contract.
type PublicChainERC20RaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func PublicChainERC20RaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (publicChainERC20 *PublicChainERC20) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*PublicChainERC20RaylsAccessManagedInvalidAuthority, error) {
	out := new(PublicChainERC20RaylsAccessManagedInvalidAuthority)
	if err := publicChainERC20.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC20RaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the PublicChainERC20 contract.
type PublicChainERC20RaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func PublicChainERC20RaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (publicChainERC20 *PublicChainERC20) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*PublicChainERC20RaylsAccessManagedMustSchedule, error) {
	out := new(PublicChainERC20RaylsAccessManagedMustSchedule)
	if err := publicChainERC20.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC20RaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the PublicChainERC20 contract.
type PublicChainERC20RaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func PublicChainERC20RaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (publicChainERC20 *PublicChainERC20) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*PublicChainERC20RaylsAccessManagedUnauthorized, error) {
	out := new(PublicChainERC20RaylsAccessManagedUnauthorized)
	if err := publicChainERC20.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC20RaylsPublicERC20HandlerAmountMustBeGreaterThanZero represents a RaylsPublicERC20Handler__AmountMustBeGreaterThanZero error raised by the PublicChainERC20 contract.
type PublicChainERC20RaylsPublicERC20HandlerAmountMustBeGreaterThanZero struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsPublicERC20Handler__AmountMustBeGreaterThanZero()
func PublicChainERC20RaylsPublicERC20HandlerAmountMustBeGreaterThanZeroErrorID() common.Hash {
	return common.HexToHash("0x9e89fcede98fda39f1542bb7ec2c47dd22b165e3513408261b79d77a5f30a35b")
}

// UnpackRaylsPublicERC20HandlerAmountMustBeGreaterThanZeroError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsPublicERC20Handler__AmountMustBeGreaterThanZero()
func (publicChainERC20 *PublicChainERC20) UnpackRaylsPublicERC20HandlerAmountMustBeGreaterThanZeroError(raw []byte) (*PublicChainERC20RaylsPublicERC20HandlerAmountMustBeGreaterThanZero, error) {
	out := new(PublicChainERC20RaylsPublicERC20HandlerAmountMustBeGreaterThanZero)
	if err := publicChainERC20.abi.UnpackIntoInterface(out, "RaylsPublicERC20HandlerAmountMustBeGreaterThanZero", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC20RaylsPublicERC20HandlerDestinationIsZeroAddress represents a RaylsPublicERC20Handler__DestinationIsZeroAddress error raised by the PublicChainERC20 contract.
type PublicChainERC20RaylsPublicERC20HandlerDestinationIsZeroAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsPublicERC20Handler__DestinationIsZeroAddress()
func PublicChainERC20RaylsPublicERC20HandlerDestinationIsZeroAddressErrorID() common.Hash {
	return common.HexToHash("0xf0d0595e68847cc60e4292c338ce4055b075bd5dfc42b6ee225ece7b11b8f15b")
}

// UnpackRaylsPublicERC20HandlerDestinationIsZeroAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsPublicERC20Handler__DestinationIsZeroAddress()
func (publicChainERC20 *PublicChainERC20) UnpackRaylsPublicERC20HandlerDestinationIsZeroAddressError(raw []byte) (*PublicChainERC20RaylsPublicERC20HandlerDestinationIsZeroAddress, error) {
	out := new(PublicChainERC20RaylsPublicERC20HandlerDestinationIsZeroAddress)
	if err := publicChainERC20.abi.UnpackIntoInterface(out, "RaylsPublicERC20HandlerDestinationIsZeroAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}
