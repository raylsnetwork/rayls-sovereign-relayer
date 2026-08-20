// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package PublicChainERC721

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

// PublicChainERC721MetaData contains all meta data concerning the PublicChainERC721 contract.
var PublicChainERC721MetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_baseUri\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_raylsNodeEndpoint\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_privateAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getApproved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentTokenId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPublicRaylsNodeEndpoint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"mintToken\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"mintTokenWithId\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerOf\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"receiveTeleportFromPrivacyNode\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"srcChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revertTeleportToPrivacyNode\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"teleportToPrivacyNode\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tokenURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RaylsPublicErc721TokenCreated\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC721IncorrectOwner\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InsufficientApproval\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721NonexistentToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsPublicERC721Handler__CallerIsNotOwnerNorApproved\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsPublicERC721Handler__DestinationIsZeroAddress\",\"inputs\":[]}]",
	ID:  "PublicChainERC721",
	Bin: "0x60806040523480156200001157600080fd5b50604051620021ef380380620021ef83398101604081905262000034916200065f565b600080546001600160a01b0319166001600160a01b03841690811782556040805163bf7e214f60e01b81529051889388938893889333938993889388938893929163bf7e214f9160048083019260209291908290030181865afa158015620000a0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000c6919062000716565b90506001600160a01b03811615620000e357620000e3816200019d565b5060019050620000f48382620007cc565b506002620001038282620007cc565b50600891506200011690508682620007cc565b506009620001258582620007cc565b506007620001348782620007cc565b50600a80546001600160a01b0319166001600160a01b0383161790556200015b826200023d565b60405130907f887df3fa3545f73c5febbd84971c4f17a29f30afb5e97c2baac67cd9b0ee7e7390600090a250506000600b555062000a07975050505050505050565b6001600160a01b038116620001d457604051638944034760e01b81526001600160a01b038216600482015260240160405180910390fd5b6000620001e062000510565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b600080546001600160a01b031662000257576000620002cf565b60008054906101000a90046001600160a01b03166001600160a01b031663bf7e214f6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620002a9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002cf919062000716565b90506001600160a01b038116620002e4575050565b6040805160028082526060820183526000926020830190803683370190505090506340c10f1960e01b8160008151811062000323576200032362000898565b6001600160e01b0319909216602092830291909101909101528051630852cd8d60e31b90829060019081106200035d576200035d62000898565b6001600160e01b031992909216602092830291909101820152604080516002808252606082018352600093919290918301908036833701905050905063bef97c9e60e01b81600081518110620003b757620003b762000898565b6001600160e01b0319909216602092830291909101909101528051632217f5f160e21b9082906001908110620003f157620003f162000898565b6001600160e01b03199290921660209283029190910190910152604080516001808252818301909252600091816020015b6040805180820190915260608082526020820152815260200190600190039081620004225750506040805160808101825260109181019182526f26a2a9a9a0a3a2afa2ac22a1aaaa27a960811b606082015290815260208101849052815191925090829060009062000498576200049862000898565b60209081029190910101526040516337af400760e11b81526001600160a01b03851690636f5e800e90620004d590889087908690600401620008ae565b600060405180830381600087803b158015620004f057600080fd5b505af115801562000505573d6000803e3d6000fd5b505050505050505050565b60008060ff196200054360017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f35620009df565b6040516020016200055691815260200190565b60408051601f1981840301815291905280516020909101201692915050565b634e487b7160e01b600052604160045260246000fd5b60005b83811015620005a85781810151838201526020016200058e565b50506000910152565b600082601f830112620005c357600080fd5b81516001600160401b0380821115620005e057620005e062000575565b604051601f8301601f19908116603f011681019082821181831017156200060b576200060b62000575565b816040528381528660208588010111156200062557600080fd5b620006388460208301602089016200058b565b9695505050505050565b80516001600160a01b03811681146200065a57600080fd5b919050565b600080600080600060a086880312156200067857600080fd5b85516001600160401b03808211156200069057600080fd5b6200069e89838a01620005b1565b96506020880151915080821115620006b557600080fd5b620006c389838a01620005b1565b95506040880151915080821115620006da57600080fd5b50620006e988828901620005b1565b935050620006fa6060870162000642565b91506200070a6080870162000642565b90509295509295909350565b6000602082840312156200072957600080fd5b620007348262000642565b9392505050565b600181811c908216806200075057607f821691505b6020821081036200077157634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620007c7576000816000526020600020601f850160051c81016020861015620007a25750805b601f850160051c820191505b81811015620007c357828155600101620007ae565b5050505b505050565b81516001600160401b03811115620007e857620007e862000575565b6200080081620007f984546200073b565b8462000777565b602080601f8311600181146200083857600084156200081f5750858301515b600019600386901b1c1916600185901b178555620007c3565b600085815260208120601f198616915b82811015620008695788860151825594840194600190910190840162000848565b5085821015620008885787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b634e487b7160e01b600052603260045260246000fd5b6001600160a01b03841681526060602080830182905284518383018190526000929160809182860190888301865b82811015620009045781516001600160e01b03191684529284019290840190600101620008dc565b505050604086820360408801528188518084528484019150848160051b850101858b0160005b83811015620009cc57601f1980888503018652825180518886528051808a8801526200095c818f89018e85016200058b565b918b0151601f9290920190921685018581038d018b87015281518d8201818152928c019350600092918d01905b80841015620009b55784516001600160e01b0319168252938c019360019390930192908c019062000989565b50978b01979550505091880191506001016200092a565b50909d9c50505050505050505050505050565b8181038181111562000a0157634e487b7160e01b600052601160045260246000fd5b92915050565b6117d88062000a176000396000f3fe608060405234801561001057600080fd5b50600436106101385760003560e01c806356189236116100b357806356189236146101f65780635f2e03551461024a5780636352211e1461025b57806370a082311461026e578063885fd7c41461021157806395d89b4114610281578063a22cb46514610289578063b88d4fde1461029c578063bef97c9e146102af578063bf7e214f146102c2578063c87b56dd146102ca578063e985e9c5146102dd57600080fd5b806301173a741461013d57806301ffc9a71461016357806306fdde0314610186578063081812fc1461019b5780630871b193146101bb578063095ea7b3146101d057806314eb966d146101e357806318160ddd146101f657806323b872dd146101fe57806340c10f191461021157806342842e0e1461022457806342966c6814610237575b600080fd5b61015061014b366004611226565b6102f0565b6040519081526020015b60405180910390f35b610176610171366004611257565b610333565b604051901515815260200161015a565b61018e610385565b60405161015a91906112c4565b6101ae6101a93660046112d7565b610417565b60405161015a91906112f0565b6101ce6101c9366004611304565b61042c565b005b6101ce6101de366004611304565b610468565b6101766101f136600461132e565b610473565b600b54610150565b6101ce61020c366004611361565b61061f565b6101ce61021f366004611304565b6106aa565b6101ce610232366004611361565b6106ca565b6101ce6102453660046112d7565b6106ea565b6000546001600160a01b03166101ae565b6101ae6102693660046112d7565b61070c565b61015061027c366004611226565b610717565b61018e61075f565b6101ce6102973660046113ab565b61076e565b6101ce6102aa3660046113f8565b610779565b6101ce6102bd3660046114d4565b610791565b6101ae6108c4565b61018e6102d83660046112d7565b6108dd565b6101766102eb366004611518565b610945565b6000610308336000356001600160e01b031916610973565b600b546103158382610abe565b600b805490600061032583611561565b90915550909150505b919050565b60006001600160e01b031982166380ac58cd60e01b148061036457506001600160e01b03198216635b5e139f60e01b145b8061037f57506301ffc9a760e01b6001600160e01b03198316145b92915050565b6060600880546103949061157a565b80601f01602080910402602001604051908101604052809291908181526020018280546103c09061157a565b801561040d5780601f106103e25761010080835404028352916020019161040d565b820191906000526020600020905b8154815290600101906020018083116103f057829003601f168201915b5050505050905090565b600061042282610b23565b5061037f82610b5b565b610442336000356001600160e01b031916610973565b61044c8282610abe565b600b548110610464576104608160016115b4565b600b555b5050565b610464828233610b76565b60006001600160a01b03841661049c57604051631d0af75360e31b815260040160405180910390fd5b60006104a78461070c565b90506104b4813386610b83565b6104de57338460405163177e802f60e01b81526004016104d59291906115c7565b60405180910390fd5b6104e784610be9565b6040805160c081018252600281526020810186905233818301526001600160a01b038781166060830152306080830152600160a0830152600054600a54935192939082169263e0d9d8469288921690610546908b908b906024016115c7565b60408051601f198184030181529181526020820180516001600160e01b03166330bb92bf60e01b179052516105819088908c906024016115c7565b60408051601f198184030181529181526020820180516001600160e01b0316632217f5f160e21b1790525160e086901b6001600160e01b03191681526105cf949392919088906004016115e0565b6020604051808303816000875af11580156105ee573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106129190611697565b5060019695505050505050565b6001600160a01b038216610649576000604051633250574960e11b81526004016104d591906112f0565b6000610656838333610c24565b9050836001600160a01b0316816001600160a01b0316146106a4576040516364283d7b60e01b81526001600160a01b03808616600483015260248201849052821660448201526064016104d5565b50505050565b6106c0336000356001600160e01b031916610973565b6104648282610abe565b6106e583838360405180602001604052806000815250610779565b505050565b610700336000356001600160e01b031916610973565b61070981610be9565b50565b600061037f82610b23565b60006001600160a01b0382166107435760006040516322718ad960e21b81526004016104d591906112f0565b506001600160a01b031660009081526004602052604090205490565b6060600980546103949061157a565b610464338383610d18565b61078484848461061f565b6106a43385858585610dae565b6107a7336000356001600160e01b031916610973565b6001600160a01b0382166108ba576040805160c08101825260028152602081018390526001600160a01b038681168284018190526060830152306080830152600160a0830152600054600a54935192939082169263e0d9d8469288921690610815908a9088906024016115c7565b60408051601f19818403018152918152602080830180516001600160e01b0316633a03ff1f60e01b1790528151908101825260008152905160e086901b6001600160e01b0319168152610870949392919088906004016115e0565b6020604051808303816000875af115801561088f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108b39190611697565b50506106a4565b6106a48282610abe565b60006108ce610ec7565b546001600160a01b0316919050565b60606108e882610b23565b5060006108f3610f29565b90506000815111610913576040518060200160405280600081525061093e565b8061091d84610f38565b60405160200161092e9291906116b0565b6040516020818303038152906040525b9392505050565b6001600160a01b03918216600090815260066020908152604080832093909416825291909152205460ff1690565b600061097d610ec7565b80549091506001600160a01b0316806109ac576000604051638944034760e01b81526004016104d591906112f0565b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa158015610a10573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a3491906116df565b92509250925082610ab5578015610a5e5760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff821615610a9a5760405163a426878960e01b81526001600160a01b038816600482015263ffffffff831660248201526044016104d5565b86604051632ecd3d0360e21b81526004016104d591906112f0565b50505050505050565b6001600160a01b038216610ae8576000604051633250574960e11b81526004016104d591906112f0565b6000610af683836000610c24565b90506001600160a01b038116156106e55760006040516339e3563760e11b81526004016104d591906112f0565b600080610b2f83610fcb565b90506001600160a01b03811661037f57604051637e27328960e01b8152600481018490526024016104d5565b6000908152600560205260409020546001600160a01b031690565b6106e58383836001610fe6565b60006001600160a01b03831615801590610be15750826001600160a01b0316846001600160a01b03161480610bbd5750610bbd8484610945565b80610be15750826001600160a01b0316610bd683610b5b565b6001600160a01b0316145b949350505050565b6000610bf86000836000610c24565b90506001600160a01b03811661046457604051637e27328960e01b8152600481018390526024016104d5565b600080610c3084610fcb565b90506001600160a01b03831615610c4c57610c4c8184866110e3565b6001600160a01b03811615610c8a57610c69600085600080610fe6565b6001600160a01b038116600090815260046020526040902080546000190190555b6001600160a01b03851615610cb9576001600160a01b0385166000908152600460205260409020805460010190555b60008481526003602052604080822080546001600160a01b0319166001600160a01b0389811691821790925591518793918516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a4949350505050565b6001600160a01b038216610d415781604051630b61174360e31b81526004016104d591906112f0565b6001600160a01b03838116600081815260066020908152604080832094871680845294825291829020805460ff191686151590811790915591519182527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a3505050565b6001600160a01b0383163b15610ec057604051630a85bd0160e11b81526001600160a01b0384169063150b7a0290610df0908890889087908790600401611735565b6020604051808303816000875af1925050508015610e2b575060408051601f3d908101601f19168201909252610e2891810190611772565b60015b610e8b573d808015610e59576040519150601f19603f3d011682016040523d82523d6000602084013e610e5e565b606091505b508051600003610e835783604051633250574960e11b81526004016104d591906112f0565b805160208201fd5b6001600160e01b03198116630a85bd0160e11b14610ebe5783604051633250574960e11b81526004016104d591906112f0565b505b5050505050565b60008060ff19610ef860017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f3561178f565b604051602001610f0a91815260200190565b60408051601f1981840301815291905280516020909101201692915050565b6060600780546103949061157a565b60606000610f4583611139565b600101905060008167ffffffffffffffff811115610f6557610f656113e2565b6040519080825280601f01601f191660200182016040528015610f8f576020820181803683370190505b5090508181016020015b600019016f181899199a1a9b1b9c1cb0b131b232b360811b600a86061a8153600a8504945084610f9957509392505050565b6000908152600360205260409020546001600160a01b031690565b8080610ffa57506001600160a01b03821615155b156110b357600061100a84610b23565b90506001600160a01b038316158015906110365750826001600160a01b0316816001600160a01b031614155b801561104957506110478184610945565b155b15611069578260405163a9fbf51f60e01b81526004016104d591906112f0565b81156110b15783856001600160a01b0316826001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45b505b5050600090815260056020526040902080546001600160a01b0319166001600160a01b0392909216919091179055565b6110ee838383610b83565b6106e5576001600160a01b03831661111c57604051637e27328960e01b8152600481018290526024016104d5565b818160405163177e802f60e01b81526004016104d59291906115c7565b60008072184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b83106111785772184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b830492506040015b6904ee2d6d415b85acef8160201b83106111a2576904ee2d6d415b85acef8160201b830492506020015b662386f26fc1000083106111c057662386f26fc10000830492506010015b6305f5e10083106111d8576305f5e100830492506008015b61271083106111ec57612710830492506004015b606483106111fe576064830492506002015b600a831061037f5760010192915050565b80356001600160a01b038116811461032e57600080fd5b60006020828403121561123857600080fd5b61093e8261120f565b6001600160e01b03198116811461070957600080fd5b60006020828403121561126957600080fd5b813561093e81611241565b60005b8381101561128f578181015183820152602001611277565b50506000910152565b600081518084526112b0816020860160208601611274565b601f01601f19169290920160200192915050565b60208152600061093e6020830184611298565b6000602082840312156112e957600080fd5b5035919050565b6001600160a01b0391909116815260200190565b6000806040838503121561131757600080fd5b6113208361120f565b946020939093013593505050565b60008060006060848603121561134357600080fd5b61134c8461120f565b95602085013595506040909401359392505050565b60008060006060848603121561137657600080fd5b61137f8461120f565b925061138d6020850161120f565b9150604084013590509250925092565b801515811461070957600080fd5b600080604083850312156113be57600080fd5b6113c78361120f565b915060208301356113d78161139d565b809150509250929050565b634e487b7160e01b600052604160045260246000fd5b6000806000806080858703121561140e57600080fd5b6114178561120f565b93506114256020860161120f565b925060408501359150606085013567ffffffffffffffff8082111561144957600080fd5b818701915087601f83011261145d57600080fd5b81358181111561146f5761146f6113e2565b604051601f8201601f19908116603f01168101908382118183101715611497576114976113e2565b816040528281528a60208487010111156114b057600080fd5b82602086016020830137600060208483010152809550505050505092959194509250565b600080600080608085870312156114ea57600080fd5b6114f38561120f565b9350602085013592506115086040860161120f565b9396929550929360600135925050565b6000806040838503121561152b57600080fd5b6115348361120f565b91506115426020840161120f565b90509250929050565b634e487b7160e01b600052601160045260246000fd5b6000600182016115735761157361154b565b5060010190565b600181811c9082168061158e57607f821691505b6020821081036115ae57634e487b7160e01b600052602260045260246000fd5b50919050565b8082018082111561037f5761037f61154b565b6001600160a01b03929092168252602082015260400190565b8581526001600160a01b0385811660208301526101406040830181905260009161160c84830188611298565b915083820360608501526116208287611298565b9250845191506007821061164457634e487b7160e01b600052602160045260246000fd5b608084810192909252602085015160a08501526040850151811660c085015260608501511660e08401528301516001600160a01b0381166101008401525060a08301516101208301529695505050505050565b6000602082840312156116a957600080fd5b5051919050565b600083516116c2818460208801611274565b8351908301906116d6818360208801611274565b01949350505050565b6000806000606084860312156116f457600080fd5b83516116ff8161139d565b602085015190935063ffffffff8116811461171957600080fd5b604085015190925061172a8161139d565b809150509250925092565b6001600160a01b038581168252841660208201526040810183905260806060820181905260009061176890830184611298565b9695505050505050565b60006020828403121561178457600080fd5b815161093e81611241565b8181038181111561037f5761037f61154b56fea2646970667358221220ff8be20eda93dec741509e2838962b1701ff01d4a24167d9cb971a17c654331264736f6c63430008180033",
}

// PublicChainERC721 is an auto generated Go binding around an Ethereum contract.
type PublicChainERC721 struct {
	abi abi.ABI
}

// NewPublicChainERC721 creates a new instance of PublicChainERC721.
func NewPublicChainERC721() *PublicChainERC721 {
	parsed, err := PublicChainERC721MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &PublicChainERC721{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *PublicChainERC721) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(string _baseUri, string _name, string _symbol, address _raylsNodeEndpoint, address _privateAddress) returns()
func (publicChainERC721 *PublicChainERC721) PackConstructor(_baseUri string, _name string, _symbol string, _raylsNodeEndpoint common.Address, _privateAddress common.Address) []byte {
	enc, err := publicChainERC721.abi.Pack("", _baseUri, _name, _symbol, _raylsNodeEndpoint, _privateAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (publicChainERC721 *PublicChainERC721) PackApprove(to common.Address, tokenId *big.Int) []byte {
	enc, err := publicChainERC721.abi.Pack("approve", to, tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (publicChainERC721 *PublicChainERC721) PackAuthority() []byte {
	enc, err := publicChainERC721.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (publicChainERC721 *PublicChainERC721) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := publicChainERC721.abi.Unpack("authority", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (publicChainERC721 *PublicChainERC721) PackBalanceOf(owner common.Address) []byte {
	enc, err := publicChainERC721.abi.Pack("balanceOf", owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (publicChainERC721 *PublicChainERC721) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := publicChainERC721.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackBurn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (publicChainERC721 *PublicChainERC721) PackBurn(tokenId *big.Int) []byte {
	enc, err := publicChainERC721.abi.Pack("burn", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackGetApproved is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (publicChainERC721 *PublicChainERC721) PackGetApproved(tokenId *big.Int) []byte {
	enc, err := publicChainERC721.abi.Pack("getApproved", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetApproved is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (publicChainERC721 *PublicChainERC721) UnpackGetApproved(data []byte) (common.Address, error) {
	out, err := publicChainERC721.abi.Unpack("getApproved", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetCurrentTokenId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x56189236.
//
// Solidity: function getCurrentTokenId() view returns(uint256)
func (publicChainERC721 *PublicChainERC721) PackGetCurrentTokenId() []byte {
	enc, err := publicChainERC721.abi.Pack("getCurrentTokenId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetCurrentTokenId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x56189236.
//
// Solidity: function getCurrentTokenId() view returns(uint256)
func (publicChainERC721 *PublicChainERC721) UnpackGetCurrentTokenId(data []byte) (*big.Int, error) {
	out, err := publicChainERC721.abi.Unpack("getCurrentTokenId", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetPublicRaylsNodeEndpoint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5f2e0355.
//
// Solidity: function getPublicRaylsNodeEndpoint() view returns(address)
func (publicChainERC721 *PublicChainERC721) PackGetPublicRaylsNodeEndpoint() []byte {
	enc, err := publicChainERC721.abi.Pack("getPublicRaylsNodeEndpoint")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetPublicRaylsNodeEndpoint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5f2e0355.
//
// Solidity: function getPublicRaylsNodeEndpoint() view returns(address)
func (publicChainERC721 *PublicChainERC721) UnpackGetPublicRaylsNodeEndpoint(data []byte) (common.Address, error) {
	out, err := publicChainERC721.abi.Unpack("getPublicRaylsNodeEndpoint", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackIsApprovedForAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (publicChainERC721 *PublicChainERC721) PackIsApprovedForAll(owner common.Address, operator common.Address) []byte {
	enc, err := publicChainERC721.abi.Pack("isApprovedForAll", owner, operator)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsApprovedForAll is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (publicChainERC721 *PublicChainERC721) UnpackIsApprovedForAll(data []byte) (bool, error) {
	out, err := publicChainERC721.abi.Unpack("isApprovedForAll", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x40c10f19.
//
// Solidity: function mint(address to, uint256 tokenId) returns()
func (publicChainERC721 *PublicChainERC721) PackMint(to common.Address, tokenId *big.Int) []byte {
	enc, err := publicChainERC721.abi.Pack("mint", to, tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackMintToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01173a74.
//
// Solidity: function mintToken(address to) returns(uint256)
func (publicChainERC721 *PublicChainERC721) PackMintToken(to common.Address) []byte {
	enc, err := publicChainERC721.abi.Pack("mintToken", to)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackMintToken is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01173a74.
//
// Solidity: function mintToken(address to) returns(uint256)
func (publicChainERC721 *PublicChainERC721) UnpackMintToken(data []byte) (*big.Int, error) {
	out, err := publicChainERC721.abi.Unpack("mintToken", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackMintTokenWithId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0871b193.
//
// Solidity: function mintTokenWithId(address to, uint256 tokenId) returns()
func (publicChainERC721 *PublicChainERC721) PackMintTokenWithId(to common.Address, tokenId *big.Int) []byte {
	enc, err := publicChainERC721.abi.Pack("mintTokenWithId", to, tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (publicChainERC721 *PublicChainERC721) PackName() []byte {
	enc, err := publicChainERC721.abi.Pack("name")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (publicChainERC721 *PublicChainERC721) UnpackName(data []byte) (string, error) {
	out, err := publicChainERC721.abi.Unpack("name", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackOwnerOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (publicChainERC721 *PublicChainERC721) PackOwnerOf(tokenId *big.Int) []byte {
	enc, err := publicChainERC721.abi.Pack("ownerOf", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackOwnerOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (publicChainERC721 *PublicChainERC721) UnpackOwnerOf(data []byte) (common.Address, error) {
	out, err := publicChainERC721.abi.Unpack("ownerOf", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackReceiveTeleportFromPrivacyNode is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbef97c9e.
//
// Solidity: function receiveTeleportFromPrivacyNode(address from, uint256 srcChainId, address to, uint256 tokenId) returns()
func (publicChainERC721 *PublicChainERC721) PackReceiveTeleportFromPrivacyNode(from common.Address, srcChainId *big.Int, to common.Address, tokenId *big.Int) []byte {
	enc, err := publicChainERC721.abi.Pack("receiveTeleportFromPrivacyNode", from, srcChainId, to, tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRevertTeleportToPrivacyNode is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x885fd7c4.
//
// Solidity: function revertTeleportToPrivacyNode(address to, uint256 tokenId) returns()
func (publicChainERC721 *PublicChainERC721) PackRevertTeleportToPrivacyNode(to common.Address, tokenId *big.Int) []byte {
	enc, err := publicChainERC721.abi.Pack("revertTeleportToPrivacyNode", to, tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSafeTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (publicChainERC721 *PublicChainERC721) PackSafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) []byte {
	enc, err := publicChainERC721.abi.Pack("safeTransferFrom", from, to, tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSafeTransferFrom0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (publicChainERC721 *PublicChainERC721) PackSafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) []byte {
	enc, err := publicChainERC721.abi.Pack("safeTransferFrom0", from, to, tokenId, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetApprovalForAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (publicChainERC721 *PublicChainERC721) PackSetApprovalForAll(operator common.Address, approved bool) []byte {
	enc, err := publicChainERC721.abi.Pack("setApprovalForAll", operator, approved)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (publicChainERC721 *PublicChainERC721) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := publicChainERC721.abi.Pack("supportsInterface", interfaceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (publicChainERC721 *PublicChainERC721) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := publicChainERC721.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (publicChainERC721 *PublicChainERC721) PackSymbol() []byte {
	enc, err := publicChainERC721.abi.Pack("symbol")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (publicChainERC721 *PublicChainERC721) UnpackSymbol(data []byte) (string, error) {
	out, err := publicChainERC721.abi.Unpack("symbol", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackTeleportToPrivacyNode is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x14eb966d.
//
// Solidity: function teleportToPrivacyNode(address to, uint256 tokenId, uint256 chainId) returns(bool)
func (publicChainERC721 *PublicChainERC721) PackTeleportToPrivacyNode(to common.Address, tokenId *big.Int, chainId *big.Int) []byte {
	enc, err := publicChainERC721.abi.Pack("teleportToPrivacyNode", to, tokenId, chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTeleportToPrivacyNode is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x14eb966d.
//
// Solidity: function teleportToPrivacyNode(address to, uint256 tokenId, uint256 chainId) returns(bool)
func (publicChainERC721 *PublicChainERC721) UnpackTeleportToPrivacyNode(data []byte) (bool, error) {
	out, err := publicChainERC721.abi.Unpack("teleportToPrivacyNode", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackTokenURI is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (publicChainERC721 *PublicChainERC721) PackTokenURI(tokenId *big.Int) []byte {
	enc, err := publicChainERC721.abi.Pack("tokenURI", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTokenURI is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (publicChainERC721 *PublicChainERC721) UnpackTokenURI(data []byte) (string, error) {
	out, err := publicChainERC721.abi.Unpack("tokenURI", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (publicChainERC721 *PublicChainERC721) PackTotalSupply() []byte {
	enc, err := publicChainERC721.abi.Pack("totalSupply")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (publicChainERC721 *PublicChainERC721) UnpackTotalSupply(data []byte) (*big.Int, error) {
	out, err := publicChainERC721.abi.Unpack("totalSupply", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (publicChainERC721 *PublicChainERC721) PackTransferFrom(from common.Address, to common.Address, tokenId *big.Int) []byte {
	enc, err := publicChainERC721.abi.Pack("transferFrom", from, to, tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PublicChainERC721Approval represents a Approval event raised by the PublicChainERC721 contract.
type PublicChainERC721Approval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const PublicChainERC721ApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (PublicChainERC721Approval) ContractEventName() string {
	return PublicChainERC721ApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (publicChainERC721 *PublicChainERC721) UnpackApprovalEvent(log *types.Log) (*PublicChainERC721Approval, error) {
	event := "Approval"
	if log.Topics[0] != publicChainERC721.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PublicChainERC721Approval)
	if len(log.Data) > 0 {
		if err := publicChainERC721.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range publicChainERC721.abi.Events[event].Inputs {
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

// PublicChainERC721ApprovalForAll represents a ApprovalForAll event raised by the PublicChainERC721 contract.
type PublicChainERC721ApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      *types.Log // Blockchain specific contextual infos
}

const PublicChainERC721ApprovalForAllEventName = "ApprovalForAll"

// ContractEventName returns the user-defined event name.
func (PublicChainERC721ApprovalForAll) ContractEventName() string {
	return PublicChainERC721ApprovalForAllEventName
}

// UnpackApprovalForAllEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (publicChainERC721 *PublicChainERC721) UnpackApprovalForAllEvent(log *types.Log) (*PublicChainERC721ApprovalForAll, error) {
	event := "ApprovalForAll"
	if log.Topics[0] != publicChainERC721.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PublicChainERC721ApprovalForAll)
	if len(log.Data) > 0 {
		if err := publicChainERC721.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range publicChainERC721.abi.Events[event].Inputs {
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

// PublicChainERC721AuthorityUpdated represents a AuthorityUpdated event raised by the PublicChainERC721 contract.
type PublicChainERC721AuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const PublicChainERC721AuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (PublicChainERC721AuthorityUpdated) ContractEventName() string {
	return PublicChainERC721AuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (publicChainERC721 *PublicChainERC721) UnpackAuthorityUpdatedEvent(log *types.Log) (*PublicChainERC721AuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != publicChainERC721.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PublicChainERC721AuthorityUpdated)
	if len(log.Data) > 0 {
		if err := publicChainERC721.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range publicChainERC721.abi.Events[event].Inputs {
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

// PublicChainERC721RaylsPublicErc721TokenCreated represents a RaylsPublicErc721TokenCreated event raised by the PublicChainERC721 contract.
type PublicChainERC721RaylsPublicErc721TokenCreated struct {
	TokenAddress common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const PublicChainERC721RaylsPublicErc721TokenCreatedEventName = "RaylsPublicErc721TokenCreated"

// ContractEventName returns the user-defined event name.
func (PublicChainERC721RaylsPublicErc721TokenCreated) ContractEventName() string {
	return PublicChainERC721RaylsPublicErc721TokenCreatedEventName
}

// UnpackRaylsPublicErc721TokenCreatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RaylsPublicErc721TokenCreated(address indexed tokenAddress)
func (publicChainERC721 *PublicChainERC721) UnpackRaylsPublicErc721TokenCreatedEvent(log *types.Log) (*PublicChainERC721RaylsPublicErc721TokenCreated, error) {
	event := "RaylsPublicErc721TokenCreated"
	if log.Topics[0] != publicChainERC721.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PublicChainERC721RaylsPublicErc721TokenCreated)
	if len(log.Data) > 0 {
		if err := publicChainERC721.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range publicChainERC721.abi.Events[event].Inputs {
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

// PublicChainERC721Transfer represents a Transfer event raised by the PublicChainERC721 contract.
type PublicChainERC721Transfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const PublicChainERC721TransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (PublicChainERC721Transfer) ContractEventName() string {
	return PublicChainERC721TransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (publicChainERC721 *PublicChainERC721) UnpackTransferEvent(log *types.Log) (*PublicChainERC721Transfer, error) {
	event := "Transfer"
	if log.Topics[0] != publicChainERC721.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PublicChainERC721Transfer)
	if len(log.Data) > 0 {
		if err := publicChainERC721.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range publicChainERC721.abi.Events[event].Inputs {
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
func (publicChainERC721 *PublicChainERC721) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], publicChainERC721.abi.Errors["ERC721IncorrectOwner"].ID.Bytes()[:4]) {
		return publicChainERC721.UnpackERC721IncorrectOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC721.abi.Errors["ERC721InsufficientApproval"].ID.Bytes()[:4]) {
		return publicChainERC721.UnpackERC721InsufficientApprovalError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC721.abi.Errors["ERC721InvalidApprover"].ID.Bytes()[:4]) {
		return publicChainERC721.UnpackERC721InvalidApproverError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC721.abi.Errors["ERC721InvalidOperator"].ID.Bytes()[:4]) {
		return publicChainERC721.UnpackERC721InvalidOperatorError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC721.abi.Errors["ERC721InvalidOwner"].ID.Bytes()[:4]) {
		return publicChainERC721.UnpackERC721InvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC721.abi.Errors["ERC721InvalidReceiver"].ID.Bytes()[:4]) {
		return publicChainERC721.UnpackERC721InvalidReceiverError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC721.abi.Errors["ERC721InvalidSender"].ID.Bytes()[:4]) {
		return publicChainERC721.UnpackERC721InvalidSenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC721.abi.Errors["ERC721NonexistentToken"].ID.Bytes()[:4]) {
		return publicChainERC721.UnpackERC721NonexistentTokenError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC721.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return publicChainERC721.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC721.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return publicChainERC721.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC721.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return publicChainERC721.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC721.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return publicChainERC721.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC721.abi.Errors["RaylsPublicERC721HandlerCallerIsNotOwnerNorApproved"].ID.Bytes()[:4]) {
		return publicChainERC721.UnpackRaylsPublicERC721HandlerCallerIsNotOwnerNorApprovedError(raw[4:])
	}
	if bytes.Equal(raw[:4], publicChainERC721.abi.Errors["RaylsPublicERC721HandlerDestinationIsZeroAddress"].ID.Bytes()[:4]) {
		return publicChainERC721.UnpackRaylsPublicERC721HandlerDestinationIsZeroAddressError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// PublicChainERC721ERC721IncorrectOwner represents a ERC721IncorrectOwner error raised by the PublicChainERC721 contract.
type PublicChainERC721ERC721IncorrectOwner struct {
	Sender  common.Address
	TokenId *big.Int
	Owner   common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721IncorrectOwner(address sender, uint256 tokenId, address owner)
func PublicChainERC721ERC721IncorrectOwnerErrorID() common.Hash {
	return common.HexToHash("0x64283d7b313c8117c125f736876fa2b4e90ea3831a4716dfdb87d2f540e26289")
}

// UnpackERC721IncorrectOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721IncorrectOwner(address sender, uint256 tokenId, address owner)
func (publicChainERC721 *PublicChainERC721) UnpackERC721IncorrectOwnerError(raw []byte) (*PublicChainERC721ERC721IncorrectOwner, error) {
	out := new(PublicChainERC721ERC721IncorrectOwner)
	if err := publicChainERC721.abi.UnpackIntoInterface(out, "ERC721IncorrectOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC721ERC721InsufficientApproval represents a ERC721InsufficientApproval error raised by the PublicChainERC721 contract.
type PublicChainERC721ERC721InsufficientApproval struct {
	Operator common.Address
	TokenId  *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InsufficientApproval(address operator, uint256 tokenId)
func PublicChainERC721ERC721InsufficientApprovalErrorID() common.Hash {
	return common.HexToHash("0x177e802f6f313bc89797ecace66d6d29ab4719cbaaacbb87367264048b1eb861")
}

// UnpackERC721InsufficientApprovalError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InsufficientApproval(address operator, uint256 tokenId)
func (publicChainERC721 *PublicChainERC721) UnpackERC721InsufficientApprovalError(raw []byte) (*PublicChainERC721ERC721InsufficientApproval, error) {
	out := new(PublicChainERC721ERC721InsufficientApproval)
	if err := publicChainERC721.abi.UnpackIntoInterface(out, "ERC721InsufficientApproval", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC721ERC721InvalidApprover represents a ERC721InvalidApprover error raised by the PublicChainERC721 contract.
type PublicChainERC721ERC721InvalidApprover struct {
	Approver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InvalidApprover(address approver)
func PublicChainERC721ERC721InvalidApproverErrorID() common.Hash {
	return common.HexToHash("0xa9fbf51f86b8e03595d59dc726bb10c329bb24f62589be276d8dd193ca0b69ea")
}

// UnpackERC721InvalidApproverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InvalidApprover(address approver)
func (publicChainERC721 *PublicChainERC721) UnpackERC721InvalidApproverError(raw []byte) (*PublicChainERC721ERC721InvalidApprover, error) {
	out := new(PublicChainERC721ERC721InvalidApprover)
	if err := publicChainERC721.abi.UnpackIntoInterface(out, "ERC721InvalidApprover", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC721ERC721InvalidOperator represents a ERC721InvalidOperator error raised by the PublicChainERC721 contract.
type PublicChainERC721ERC721InvalidOperator struct {
	Operator common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InvalidOperator(address operator)
func PublicChainERC721ERC721InvalidOperatorErrorID() common.Hash {
	return common.HexToHash("0x5b08ba185e8f577075361f3a3555a6580a227ce22734dcc979c1aeadf894658b")
}

// UnpackERC721InvalidOperatorError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InvalidOperator(address operator)
func (publicChainERC721 *PublicChainERC721) UnpackERC721InvalidOperatorError(raw []byte) (*PublicChainERC721ERC721InvalidOperator, error) {
	out := new(PublicChainERC721ERC721InvalidOperator)
	if err := publicChainERC721.abi.UnpackIntoInterface(out, "ERC721InvalidOperator", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC721ERC721InvalidOwner represents a ERC721InvalidOwner error raised by the PublicChainERC721 contract.
type PublicChainERC721ERC721InvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InvalidOwner(address owner)
func PublicChainERC721ERC721InvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x89c62b6479af2e623826dcc39c5133061d35b66d72de92833401dd2fd6567480")
}

// UnpackERC721InvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InvalidOwner(address owner)
func (publicChainERC721 *PublicChainERC721) UnpackERC721InvalidOwnerError(raw []byte) (*PublicChainERC721ERC721InvalidOwner, error) {
	out := new(PublicChainERC721ERC721InvalidOwner)
	if err := publicChainERC721.abi.UnpackIntoInterface(out, "ERC721InvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC721ERC721InvalidReceiver represents a ERC721InvalidReceiver error raised by the PublicChainERC721 contract.
type PublicChainERC721ERC721InvalidReceiver struct {
	Receiver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InvalidReceiver(address receiver)
func PublicChainERC721ERC721InvalidReceiverErrorID() common.Hash {
	return common.HexToHash("0x64a0ae9278f805eaf991dcd18ca78756d280b7508b764ef1b255c55845c11df9")
}

// UnpackERC721InvalidReceiverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InvalidReceiver(address receiver)
func (publicChainERC721 *PublicChainERC721) UnpackERC721InvalidReceiverError(raw []byte) (*PublicChainERC721ERC721InvalidReceiver, error) {
	out := new(PublicChainERC721ERC721InvalidReceiver)
	if err := publicChainERC721.abi.UnpackIntoInterface(out, "ERC721InvalidReceiver", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC721ERC721InvalidSender represents a ERC721InvalidSender error raised by the PublicChainERC721 contract.
type PublicChainERC721ERC721InvalidSender struct {
	Sender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InvalidSender(address sender)
func PublicChainERC721ERC721InvalidSenderErrorID() common.Hash {
	return common.HexToHash("0x73c6ac6e10798e95d99e1f130d923eb40193ecb8d094ec3dce93292564eb3b17")
}

// UnpackERC721InvalidSenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InvalidSender(address sender)
func (publicChainERC721 *PublicChainERC721) UnpackERC721InvalidSenderError(raw []byte) (*PublicChainERC721ERC721InvalidSender, error) {
	out := new(PublicChainERC721ERC721InvalidSender)
	if err := publicChainERC721.abi.UnpackIntoInterface(out, "ERC721InvalidSender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC721ERC721NonexistentToken represents a ERC721NonexistentToken error raised by the PublicChainERC721 contract.
type PublicChainERC721ERC721NonexistentToken struct {
	TokenId *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721NonexistentToken(uint256 tokenId)
func PublicChainERC721ERC721NonexistentTokenErrorID() common.Hash {
	return common.HexToHash("0x7e273289a3a9ef6670f06df7dca227856fc925e956db96980692764a8bc734d7")
}

// UnpackERC721NonexistentTokenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721NonexistentToken(uint256 tokenId)
func (publicChainERC721 *PublicChainERC721) UnpackERC721NonexistentTokenError(raw []byte) (*PublicChainERC721ERC721NonexistentToken, error) {
	out := new(PublicChainERC721ERC721NonexistentToken)
	if err := publicChainERC721.abi.UnpackIntoInterface(out, "ERC721NonexistentToken", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC721RaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the PublicChainERC721 contract.
type PublicChainERC721RaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func PublicChainERC721RaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (publicChainERC721 *PublicChainERC721) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*PublicChainERC721RaylsAccessManagedContractPaused, error) {
	out := new(PublicChainERC721RaylsAccessManagedContractPaused)
	if err := publicChainERC721.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC721RaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the PublicChainERC721 contract.
type PublicChainERC721RaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func PublicChainERC721RaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (publicChainERC721 *PublicChainERC721) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*PublicChainERC721RaylsAccessManagedInvalidAuthority, error) {
	out := new(PublicChainERC721RaylsAccessManagedInvalidAuthority)
	if err := publicChainERC721.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC721RaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the PublicChainERC721 contract.
type PublicChainERC721RaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func PublicChainERC721RaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (publicChainERC721 *PublicChainERC721) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*PublicChainERC721RaylsAccessManagedMustSchedule, error) {
	out := new(PublicChainERC721RaylsAccessManagedMustSchedule)
	if err := publicChainERC721.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC721RaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the PublicChainERC721 contract.
type PublicChainERC721RaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func PublicChainERC721RaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (publicChainERC721 *PublicChainERC721) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*PublicChainERC721RaylsAccessManagedUnauthorized, error) {
	out := new(PublicChainERC721RaylsAccessManagedUnauthorized)
	if err := publicChainERC721.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC721RaylsPublicERC721HandlerCallerIsNotOwnerNorApproved represents a RaylsPublicERC721Handler__CallerIsNotOwnerNorApproved error raised by the PublicChainERC721 contract.
type PublicChainERC721RaylsPublicERC721HandlerCallerIsNotOwnerNorApproved struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsPublicERC721Handler__CallerIsNotOwnerNorApproved()
func PublicChainERC721RaylsPublicERC721HandlerCallerIsNotOwnerNorApprovedErrorID() common.Hash {
	return common.HexToHash("0x37a6b6b38c17beaec892f30b97e1490f5665de5e1803a1f0163530a50eeb0866")
}

// UnpackRaylsPublicERC721HandlerCallerIsNotOwnerNorApprovedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsPublicERC721Handler__CallerIsNotOwnerNorApproved()
func (publicChainERC721 *PublicChainERC721) UnpackRaylsPublicERC721HandlerCallerIsNotOwnerNorApprovedError(raw []byte) (*PublicChainERC721RaylsPublicERC721HandlerCallerIsNotOwnerNorApproved, error) {
	out := new(PublicChainERC721RaylsPublicERC721HandlerCallerIsNotOwnerNorApproved)
	if err := publicChainERC721.abi.UnpackIntoInterface(out, "RaylsPublicERC721HandlerCallerIsNotOwnerNorApproved", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PublicChainERC721RaylsPublicERC721HandlerDestinationIsZeroAddress represents a RaylsPublicERC721Handler__DestinationIsZeroAddress error raised by the PublicChainERC721 contract.
type PublicChainERC721RaylsPublicERC721HandlerDestinationIsZeroAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsPublicERC721Handler__DestinationIsZeroAddress()
func PublicChainERC721RaylsPublicERC721HandlerDestinationIsZeroAddressErrorID() common.Hash {
	return common.HexToHash("0xe857ba98f9495f45ce97c259b2cc41e849a9810510259c895de387c4f9579350")
}

// UnpackRaylsPublicERC721HandlerDestinationIsZeroAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsPublicERC721Handler__DestinationIsZeroAddress()
func (publicChainERC721 *PublicChainERC721) UnpackRaylsPublicERC721HandlerDestinationIsZeroAddressError(raw []byte) (*PublicChainERC721RaylsPublicERC721HandlerDestinationIsZeroAddress, error) {
	out := new(PublicChainERC721RaylsPublicERC721HandlerDestinationIsZeroAddress)
	if err := publicChainERC721.abi.UnpackIntoInterface(out, "RaylsPublicERC721HandlerDestinationIsZeroAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}
