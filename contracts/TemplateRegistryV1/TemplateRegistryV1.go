// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TemplateRegistryV1

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

// TemplateRegistryV1MetaData contains all meta data concerning the TemplateRegistryV1 contract.
var TemplateRegistryV1MetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"contractVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"endpoint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRaylsEndpoint\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAddressByResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getKey\",\"inputs\":[{\"name\":\"bytecodeHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getTemplate\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structITemplateRegistry.Template\",\"components\":[{\"name\":\"bytecodeHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"approvedAtBlock\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"authority_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"propose\",\"inputs\":[{\"name\":\"bytecodeHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reseedStandardTemplate\",\"inputs\":[{\"name\":\"bytecodeHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resourceId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"revoke\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"seedStandardTemplate\",\"inputs\":[{\"name\":\"bytecodeHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StandardTemplateReseeded\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"bytecodeHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"indexed\":false,\"internalType\":\"bytes4\"},{\"name\":\"approvedAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TemplateApproved\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"bytecodeHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"indexed\":false,\"internalType\":\"bytes4\"},{\"name\":\"approvedAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TemplateProposed\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"bytecodeHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"indexed\":false,\"internalType\":\"bytes4\"},{\"name\":\"proposer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TemplateRevoked\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"bytecodeHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"indexed\":false,\"internalType\":\"bytes4\"},{\"name\":\"revokedAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TemplateRegistryV1__AlreadyApproved\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"TemplateRegistryV1__AlreadyRegistered\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"TemplateRegistryV1__EmptyBytecodeHash\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TemplateRegistryV1__EmptySignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TemplateRegistryV1__NotApproved\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"TemplateRegistryV1__NotProposed\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"TemplateRegistryV1__RevokedTemplate\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	ID:  "TemplateRegistryV1",
	Bin: "0x60a06040523060805234801561001457600080fd5b50608051611fa761003e600039600081816110060152818161102f01526111610152611fa76000f3fe6080604052600436106100de5760003560e01c80639757739b116100855780639757739b146101e3578063a042893b14610210578063a0a8e46014610230578063a53a1adf14610244578063ad3cb1cc14610264578063b75c7dc6146102a2578063bf7e214f146102c2578063c4d66de8146102d7578063e844d93c146102f757600080fd5b8063042e9a43146100e357806311f50c85146101165780633dd9380214610143578063485cc955146101635780634f1ef2861461018557806352d1902d146101985780635e280f11146101ad5780635f997c5b146101cd575b600080fd5b3480156100ef57600080fd5b506101036100fe36600461175b565b610317565b6040519081526020015b60405180910390f35b34801561012257600080fd5b506101366101313660046117d6565b6104d1565b60405161010d91906117ef565b34801561014f57600080fd5b5061010361015e366004611803565b610545565b34801561016f57600080fd5b5061018361017e366004611855565b610578565b005b610183610193366004611899565b610689565b3480156101a457600080fd5b506101036106a8565b3480156101b957600080fd5b50600054610136906001600160a01b031681565b3480156101d957600080fd5b5061010360015481565b3480156101ef57600080fd5b506102036101fe3660046117d6565b6106c5565b60405161010d91906119ac565b34801561021c57600080fd5b5061018361022b36600461175b565b6107f2565b34801561023c57600080fd5b506002610103565b34801561025057600080fd5b5061018361025f3660046117d6565b6109fd565b34801561027057600080fd5b50610295604051806040016040528060058152602001640352e302e360dc1b81525081565b60405161010d9190611a10565b3480156102ae57600080fd5b506101836102bd3660046117d6565b610b7f565b3480156102ce57600080fd5b50610136610c5f565b3480156102e357600080fd5b506101836102f2366004611a23565b610c78565b34801561030357600080fd5b5061018361031236600461175b565b610ca2565b60008361033757604051631b64251960e01b815260040160405180910390fd5b600082900361035957604051634248016960e11b815260040160405180910390fd5b6000610366858585610f39565b600081815260026020526040902054909350909150156103a15760405163b60533b960e01b8152600481018390526024015b60405180910390fd5b6040518060a0016040528086815260200185858080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201829052509385525050506001600160e01b0319841660208084019190915260408084018390526060909301829052858252600281529190208251815590820151600182019061042f9082611ac2565b506040828101516002909201805460608501516080909501511515600160601b0260ff60601b196001600160401b03909616600160201b026001600160601b031990921660e09590951c949094171793909316919091179091555182907f4f5dc5cc6e287f6e7b6dba0daf890f0a47ac3f5cd64975f3aa2e0165ae1a67fb906104c19088908890889087903390611baa565b60405180910390a2509392505050565b600080546040516311f50c8560e01b8152600481018490526001600160a01b03909116906311f50c8590602401602060405180830381865afa15801561051b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061053f9190611bf0565b92915050565b6000828260405160200161055a929190611c0d565b60405160208183030381529060405280519060200120905092915050565b6000610582610f87565b805490915060ff600160401b82041615906001600160401b03166000811580156105a95750825b90506000826001600160401b031660011480156105c55750303b155b9050811580156105d3575080155b156105f15760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561061b57845460ff60401b1916600160401b1785555b610623610fb0565b61062c87610c78565b600660015561063a86610fba565b831561068057845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50505050505050565b610691610ffb565b61069a82611089565b6106a482826110a2565b5050565b60006106b2611156565b50600080516020611f5283398151915290565b6040805160a08101825260008082526060602083018190529282018190529181018290526080810191909152600260008381526020019081526020016000206040518060a00160405290816000820154815260200160018201805461072990611a40565b80601f016020809104026020016040519081016040528092919081815260200182805461075590611a40565b80156107a25780601f10610777576101008083540402835291602001916107a2565b820191906000526020600020905b81548152906001019060200180831161078557829003601f168201915b50505091835250506002919091015460e081901b6001600160e01b0319166020830152600160201b81046001600160401b03166040830152600160601b900460ff16151560609091015292915050565b610808336000356001600160e01b03191661119f565b8261082657604051631b64251960e01b815260040160405180910390fd5b600081900361084857604051634248016960e11b815260040160405180910390fd5b600080610856858585610f39565b60008181526002602052604090205491935091501561088b5760405163b60533b960e01b815260048101829052602401610398565b60004390506040518060a0016040528087815260200186868080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201829052509385525050506001600160e01b031986166020808401919091526001600160401b03851660408085019190915260016060909401849052868352600282529091208351815590830151909182019061092b9082611ac2565b506040828101516002909201805460608501516080909501511515600160601b0260ff60601b196001600160401b03909616600160201b026001600160601b031990921660e09590951c94909417179390931691909117909155518290600080516020611f32833981519152906109ab9089908990899089908890611c25565b60405180910390a26109f58686868080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508692506112e1915050565b505050505050565b610a13336000356001600160e01b03191661119f565b60008181526002602052604090208054610a4357604051631503e89360e21b815260048101839052602401610398565b6002810154600160601b900460ff1615610a73576040516345acca9960e11b815260048101839052602401610398565b6002810180546001600160401b0343908116600160201b02600160201b600160681b031990921691909117600160601b1791829055825460405191928592600080516020611f3283398151915292610ad7929091600188019160e01b908790611c6b565b60405180910390a2610b7a8260000154836001018054610af690611a40565b80601f0160208091040260200160405190810160405280929190818152602001828054610b2290611a40565b8015610b6f5780601f10610b4457610100808354040283529160200191610b6f565b820191906000526020600020905b815481529060010190602001808311610b5257829003601f168201915b5050505050836112e1565b505050565b610b95336000356001600160e01b03191661119f565b600081815260026020819052604090912090810154600160601b900460ff16610bd45760405163116a690760e21b815260048101839052602401610398565b6002810180546001600160401b0343908116600160201b02600160201b600160681b031990921691909117918290558254604051919285927f475624dece05a7ddfc13e4e9ff675a83aa03ab9b12aa1c7c45bead51cd2ed5d392610c4292909160e09190911b908690611d2a565b60405180910390a281546002830154610b7a919060e01b836113d9565b6000610c69611418565b546001600160a01b0316919050565b610c8061147a565b600080546001600160a01b0319166001600160a01b0392909216919091179055565b610cb8336000356001600160e01b03191661119f565b82610cd657604051631b64251960e01b815260040160405180910390fd5b6000819003610cf857604051634248016960e11b815260040160405180910390fd5b600080610d06858585610f39565b600081815260026020526040902080549294509092509015801590610d3757506002810154600160601b900460ff16155b8015610d5657506002810154600160201b90046001600160401b031615155b15610d7757604051637151679b60e01b815260048101839052602401610398565b6002810154600160601b900460ff1615610d9357505050505050565b60004390506040518060a0016040528088815260200187878080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201829052509385525050506001600160e01b031987166020808401919091526001600160401b038516604080850191909152600160609094018490528783526002825290912083518155908301519091820190610e339082611ac2565b506040828101516002909201805460608501516080909501511515600160601b0260ff60601b196001600160401b03909616600160201b026001600160601b031990921660e09590951c94909417179390931691909117909155518390600080516020611f3283398151915290610eb3908a908a908a908a908890611c25565b60405180910390a2827f68cecfa18fa362272c2bb2f520dca293d275c46a533ac95357cb4f2bcc83078b888684604051610eef93929190611d2a565b60405180910390a26106808787878080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508692506112e1915050565b6000808383604051610f4c929190611d53565b6040519081900381209250610f679086908490602001611c0d565b604051602081830303815290604052805190602001209050935093915050565b6000807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0061053f565b610fb861147a565b565b6000610fc4611418565b80549091506001600160a01b031615610ff25781604051638944034760e01b815260040161039891906117ef565b6106a48261149f565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016148061106b57507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031661105f61152f565b6001600160a01b031614155b15610fb85760405163703e46dd60e11b815260040160405180910390fd5b61109f336000356001600160e01b03191661119f565b50565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa9250505080156110fc575060408051601f3d908101601f191682019092526110f991810190611d63565b60015b61111b5781604051634c9c8ce360e01b815260040161039891906117ef565b600080516020611f52833981519152811461114c57604051632a87526960e21b815260048101829052602401610398565b610b7a8383611545565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614610fb85760405163703e46dd60e11b815260040160405180910390fd5b60006111a9611418565b80549091506001600160a01b0316806111d8576000604051638944034760e01b815260040161039891906117ef565b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa15801561123c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112609190611d91565b9250925092508261068057801561128a5760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff8216156112c65760405163a426878960e01b81526001600160a01b038816600482015263ffffffff83166024820152604401610398565b86604051632ecd3d0360e21b815260040161039891906117ef565b6112e9611724565b600080546040516001600160a01b03909116916341d7174491600690632bcda3a160e21b90611320908a908a908a90602401611ddf565b60408051601f19818403018152918152602080830180516001600160e01b03166001600160e01b03199586161790528151808201835260008082528351808401855281815284519384018552908352925160e089901b909516855261138f969594909291908a90600401611e11565b6020604051808303816000875af11580156113ae573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113d29190611d63565b5050505050565b6113e1611724565b600080546040516001600160a01b03909116916341d71744916006906388ee8c0160e01b90611320908a908a908a90602401611d2a565b60008060ff1961144960017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f35611ef4565b60405160200161145b91815260200190565b60408051601f1981840301815291905280516020909101201692915050565b61148261159b565b610fb857604051631afcd79f60e31b815260040160405180910390fd5b6001600160a01b0381166114c85780604051638944034760e01b815260040161039891906117ef565b60006114d2611418565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b6000600080516020611f52833981519152610c69565b61154e826115b5565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a280511561159357610b7a8282611611565b6106a4611687565b60006115a5610f87565b54600160401b900460ff16919050565b806001600160a01b03163b6000036115e25780604051634c9c8ce360e01b815260040161039891906117ef565b600080516020611f5283398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080846001600160a01b03168460405161162e9190611f15565b600060405180830381855af49150503d8060008114611669576040519150601f19603f3d011682016040523d82523d6000602084013e61166e565b606091505b509150915061167e8583836116a6565b95945050505050565b3415610fb85760405163b398979f60e01b815260040160405180910390fd5b6060826116bb576116b6826116fc565b6116f5565b81511580156116d257506001600160a01b0384163b155b156116f25783604051639996b31560e01b815260040161039891906117ef565b50805b9392505050565b80511561170b57805160208201fd5b60405163d6bda27560e01b815260040160405180910390fd5b6040805160c08101909152806000815260006020820181905260408201819052606082018190526080820181905260a09091015290565b60008060006040848603121561177057600080fd5b8335925060208401356001600160401b038082111561178e57600080fd5b818601915086601f8301126117a257600080fd5b8135818111156117b157600080fd5b8760208285010111156117c357600080fd5b6020830194508093505050509250925092565b6000602082840312156117e857600080fd5b5035919050565b6001600160a01b0391909116815260200190565b6000806040838503121561181657600080fd5b8235915060208301356001600160e01b03198116811461183557600080fd5b809150509250929050565b6001600160a01b038116811461109f57600080fd5b6000806040838503121561186857600080fd5b823561187381611840565b9150602083013561183581611840565b634e487b7160e01b600052604160045260246000fd5b600080604083850312156118ac57600080fd5b82356118b781611840565b915060208301356001600160401b03808211156118d357600080fd5b818501915085601f8301126118e757600080fd5b8135818111156118f9576118f9611883565b604051601f8201601f19908116603f0116810190838211818310171561192157611921611883565b8160405282815288602084870101111561193a57600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b60005b8381101561197757818101518382015260200161195f565b50506000910152565b6000815180845261199881602086016020860161195c565b601f01601f19169290920160200192915050565b60208152815160208201526000602083015160a060408401526119d260c0840182611980565b905063ffffffff60e01b60408501511660608401526001600160401b0360608501511660808401526080840151151560a08401528091505092915050565b6020815260006116f56020830184611980565b600060208284031215611a3557600080fd5b81356116f581611840565b600181811c90821680611a5457607f821691505b602082108103611a7457634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115610b7a576000816000526020600020601f850160051c81016020861015611aa35750805b601f850160051c820191505b818110156109f557828155600101611aaf565b81516001600160401b03811115611adb57611adb611883565b611aef81611ae98454611a40565b84611a7a565b602080601f831160018114611b245760008415611b0c5750858301515b600019600386901b1c1916600185901b1785556109f5565b600085815260208120601f198616915b82811015611b5357888601518255948401946001909101908401611b34565b5085821015611b715787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b858152608060208201526000611bc4608083018688611b81565b6001600160e01b0319949094166040830152506001600160a01b03919091166060909101529392505050565b600060208284031215611c0257600080fd5b81516116f581611840565b9182526001600160e01b031916602082015260400190565b858152608060208201526000611c3f608083018688611b81565b6001600160e01b0319949094166040830152506001600160401b03919091166060909101529392505050565b848152600060206080602084015260008654611c8681611a40565b80608087015260a0600180841660008114611ca85760018114611cc457611cf4565b60ff19851660a08a015260a084151560051b8a01019550611cf4565b8b600052602060002060005b85811015611ceb5781548b8201860152908301908801611cd0565b8a0160a0019650505b505050506001600160e01b031987166040860152509150611d129050565b6001600160401b038316606083015295945050505050565b9283526001600160e01b03199190911660208301526001600160401b0316604082015260600190565b8183823760009101908152919050565b600060208284031215611d7557600080fd5b5051919050565b80518015158114611d8c57600080fd5b919050565b600080600060608486031215611da657600080fd5b611daf84611d7c565b9250602084015163ffffffff81168114611dc857600080fd5b9150611dd660408501611d7c565b90509250925092565b838152606060208201526000611df86060830185611980565b90506001600160401b0383166040830152949350505050565b6000610180898352886020840152806040840152611e3181840189611980565b90508281036060840152611e458188611980565b90508281036080840152611e598187611980565b905082810360a0840152611e6d8186611980565b9150508251600d8110611e9057634e487b7160e01b600052602160045260246000fd5b60c0830152602083015160e083015260408301516001600160a01b039081166101008401526060840151166101208301526080830151611edc6101408401826001600160a01b03169052565b5060a083015161016083015298975050505050505050565b8181038181111561053f57634e487b7160e01b600052601160045260246000fd5b60008251611f2781846020870161195c565b919091019291505056fe17b3b38bbe4b9f31d33cd2b5bc659cce23eb821304267fa93cb462d4cb5ed6f5360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbca2646970667358221220aac588a68d515be9ebe2b83358a76a60d0116c70563330e9e970ff7511f9e93264736f6c63430008180033",
}

// TemplateRegistryV1 is an auto generated Go binding around an Ethereum contract.
type TemplateRegistryV1 struct {
	abi abi.ABI
}

// NewTemplateRegistryV1 creates a new instance of TemplateRegistryV1.
func NewTemplateRegistryV1() *TemplateRegistryV1 {
	parsed, err := TemplateRegistryV1MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &TemplateRegistryV1{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *TemplateRegistryV1) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackUPGRADEINTERFACEVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (templateRegistryV1 *TemplateRegistryV1) PackUPGRADEINTERFACEVERSION() []byte {
	enc, err := templateRegistryV1.abi.Pack("UPGRADE_INTERFACE_VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUPGRADEINTERFACEVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (templateRegistryV1 *TemplateRegistryV1) UnpackUPGRADEINTERFACEVERSION(data []byte) (string, error) {
	out, err := templateRegistryV1.abi.Unpack("UPGRADE_INTERFACE_VERSION", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa53a1adf.
//
// Solidity: function approve(bytes32 key) returns()
func (templateRegistryV1 *TemplateRegistryV1) PackApprove(key [32]byte) []byte {
	enc, err := templateRegistryV1.abi.Pack("approve", key)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (templateRegistryV1 *TemplateRegistryV1) PackAuthority() []byte {
	enc, err := templateRegistryV1.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (templateRegistryV1 *TemplateRegistryV1) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := templateRegistryV1.abi.Unpack("authority", data)
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
func (templateRegistryV1 *TemplateRegistryV1) PackContractVersion() []byte {
	enc, err := templateRegistryV1.abi.Pack("contractVersion")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackContractVersion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa0a8e460.
//
// Solidity: function contractVersion() pure returns(uint256)
func (templateRegistryV1 *TemplateRegistryV1) UnpackContractVersion(data []byte) (*big.Int, error) {
	out, err := templateRegistryV1.abi.Unpack("contractVersion", data)
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
func (templateRegistryV1 *TemplateRegistryV1) PackEndpoint() []byte {
	enc, err := templateRegistryV1.abi.Pack("endpoint")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackEndpoint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (templateRegistryV1 *TemplateRegistryV1) UnpackEndpoint(data []byte) (common.Address, error) {
	out, err := templateRegistryV1.abi.Unpack("endpoint", data)
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
func (templateRegistryV1 *TemplateRegistryV1) PackGetAddressByResourceId(resourceId [32]byte) []byte {
	enc, err := templateRegistryV1.abi.Pack("getAddressByResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAddressByResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (templateRegistryV1 *TemplateRegistryV1) UnpackGetAddressByResourceId(data []byte) (common.Address, error) {
	out, err := templateRegistryV1.abi.Unpack("getAddressByResourceId", data)
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
func (templateRegistryV1 *TemplateRegistryV1) PackGetKey(bytecodeHash [32]byte, selector [4]byte) []byte {
	enc, err := templateRegistryV1.abi.Pack("getKey", bytecodeHash, selector)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetKey is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3dd93802.
//
// Solidity: function getKey(bytes32 bytecodeHash, bytes4 selector) pure returns(bytes32)
func (templateRegistryV1 *TemplateRegistryV1) UnpackGetKey(data []byte) ([32]byte, error) {
	out, err := templateRegistryV1.abi.Unpack("getKey", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackGetTemplate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9757739b.
//
// Solidity: function getTemplate(bytes32 key) view returns((bytes32,string,bytes4,uint64,bool))
func (templateRegistryV1 *TemplateRegistryV1) PackGetTemplate(key [32]byte) []byte {
	enc, err := templateRegistryV1.abi.Pack("getTemplate", key)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTemplate is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9757739b.
//
// Solidity: function getTemplate(bytes32 key) view returns((bytes32,string,bytes4,uint64,bool))
func (templateRegistryV1 *TemplateRegistryV1) UnpackGetTemplate(data []byte) (ITemplateRegistryTemplate, error) {
	out, err := templateRegistryV1.abi.Unpack("getTemplate", data)
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
func (templateRegistryV1 *TemplateRegistryV1) PackInitialize(endpoint common.Address, authority common.Address) []byte {
	enc, err := templateRegistryV1.abi.Pack("initialize", endpoint, authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackInitialize0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc4d66de8.
//
// Solidity: function initialize(address _endpoint) returns()
func (templateRegistryV1 *TemplateRegistryV1) PackInitialize0(endpoint common.Address) []byte {
	enc, err := templateRegistryV1.abi.Pack("initialize0", endpoint)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackPropose is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x042e9a43.
//
// Solidity: function propose(bytes32 bytecodeHash, string signature) returns(bytes32 key)
func (templateRegistryV1 *TemplateRegistryV1) PackPropose(bytecodeHash [32]byte, signature string) []byte {
	enc, err := templateRegistryV1.abi.Pack("propose", bytecodeHash, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackPropose is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x042e9a43.
//
// Solidity: function propose(bytes32 bytecodeHash, string signature) returns(bytes32 key)
func (templateRegistryV1 *TemplateRegistryV1) UnpackPropose(data []byte) ([32]byte, error) {
	out, err := templateRegistryV1.abi.Unpack("propose", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackProxiableUUID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (templateRegistryV1 *TemplateRegistryV1) PackProxiableUUID() []byte {
	enc, err := templateRegistryV1.abi.Pack("proxiableUUID")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackProxiableUUID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (templateRegistryV1 *TemplateRegistryV1) UnpackProxiableUUID(data []byte) ([32]byte, error) {
	out, err := templateRegistryV1.abi.Unpack("proxiableUUID", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackReseedStandardTemplate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe844d93c.
//
// Solidity: function reseedStandardTemplate(bytes32 bytecodeHash, string signature) returns()
func (templateRegistryV1 *TemplateRegistryV1) PackReseedStandardTemplate(bytecodeHash [32]byte, signature string) []byte {
	enc, err := templateRegistryV1.abi.Pack("reseedStandardTemplate", bytecodeHash, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (templateRegistryV1 *TemplateRegistryV1) PackResourceId() []byte {
	enc, err := templateRegistryV1.abi.Pack("resourceId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (templateRegistryV1 *TemplateRegistryV1) UnpackResourceId(data []byte) ([32]byte, error) {
	out, err := templateRegistryV1.abi.Unpack("resourceId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackRevoke is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb75c7dc6.
//
// Solidity: function revoke(bytes32 key) returns()
func (templateRegistryV1 *TemplateRegistryV1) PackRevoke(key [32]byte) []byte {
	enc, err := templateRegistryV1.abi.Pack("revoke", key)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSeedStandardTemplate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa042893b.
//
// Solidity: function seedStandardTemplate(bytes32 bytecodeHash, string signature) returns()
func (templateRegistryV1 *TemplateRegistryV1) PackSeedStandardTemplate(bytecodeHash [32]byte, signature string) []byte {
	enc, err := templateRegistryV1.abi.Pack("seedStandardTemplate", bytecodeHash, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpgradeToAndCall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (templateRegistryV1 *TemplateRegistryV1) PackUpgradeToAndCall(newImplementation common.Address, data []byte) []byte {
	enc, err := templateRegistryV1.abi.Pack("upgradeToAndCall", newImplementation, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// TemplateRegistryV1AuthorityUpdated represents a AuthorityUpdated event raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1AuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const TemplateRegistryV1AuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (TemplateRegistryV1AuthorityUpdated) ContractEventName() string {
	return TemplateRegistryV1AuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (templateRegistryV1 *TemplateRegistryV1) UnpackAuthorityUpdatedEvent(log *types.Log) (*TemplateRegistryV1AuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != templateRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TemplateRegistryV1AuthorityUpdated)
	if len(log.Data) > 0 {
		if err := templateRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range templateRegistryV1.abi.Events[event].Inputs {
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

// TemplateRegistryV1Initialized represents a Initialized event raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1Initialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const TemplateRegistryV1InitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (TemplateRegistryV1Initialized) ContractEventName() string {
	return TemplateRegistryV1InitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (templateRegistryV1 *TemplateRegistryV1) UnpackInitializedEvent(log *types.Log) (*TemplateRegistryV1Initialized, error) {
	event := "Initialized"
	if log.Topics[0] != templateRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TemplateRegistryV1Initialized)
	if len(log.Data) > 0 {
		if err := templateRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range templateRegistryV1.abi.Events[event].Inputs {
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

// TemplateRegistryV1StandardTemplateReseeded represents a StandardTemplateReseeded event raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1StandardTemplateReseeded struct {
	Key          [32]byte
	BytecodeHash [32]byte
	Selector     [4]byte
	ApprovedAt   uint64
	Raw          *types.Log // Blockchain specific contextual infos
}

const TemplateRegistryV1StandardTemplateReseededEventName = "StandardTemplateReseeded"

// ContractEventName returns the user-defined event name.
func (TemplateRegistryV1StandardTemplateReseeded) ContractEventName() string {
	return TemplateRegistryV1StandardTemplateReseededEventName
}

// UnpackStandardTemplateReseededEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event StandardTemplateReseeded(bytes32 indexed key, bytes32 bytecodeHash, bytes4 selector, uint64 approvedAt)
func (templateRegistryV1 *TemplateRegistryV1) UnpackStandardTemplateReseededEvent(log *types.Log) (*TemplateRegistryV1StandardTemplateReseeded, error) {
	event := "StandardTemplateReseeded"
	if log.Topics[0] != templateRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TemplateRegistryV1StandardTemplateReseeded)
	if len(log.Data) > 0 {
		if err := templateRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range templateRegistryV1.abi.Events[event].Inputs {
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

// TemplateRegistryV1TemplateApproved represents a TemplateApproved event raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1TemplateApproved struct {
	Key          [32]byte
	BytecodeHash [32]byte
	Signature    string
	Selector     [4]byte
	ApprovedAt   uint64
	Raw          *types.Log // Blockchain specific contextual infos
}

const TemplateRegistryV1TemplateApprovedEventName = "TemplateApproved"

// ContractEventName returns the user-defined event name.
func (TemplateRegistryV1TemplateApproved) ContractEventName() string {
	return TemplateRegistryV1TemplateApprovedEventName
}

// UnpackTemplateApprovedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TemplateApproved(bytes32 indexed key, bytes32 bytecodeHash, string signature, bytes4 selector, uint64 approvedAt)
func (templateRegistryV1 *TemplateRegistryV1) UnpackTemplateApprovedEvent(log *types.Log) (*TemplateRegistryV1TemplateApproved, error) {
	event := "TemplateApproved"
	if log.Topics[0] != templateRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TemplateRegistryV1TemplateApproved)
	if len(log.Data) > 0 {
		if err := templateRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range templateRegistryV1.abi.Events[event].Inputs {
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

// TemplateRegistryV1TemplateProposed represents a TemplateProposed event raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1TemplateProposed struct {
	Key          [32]byte
	BytecodeHash [32]byte
	Signature    string
	Selector     [4]byte
	Proposer     common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const TemplateRegistryV1TemplateProposedEventName = "TemplateProposed"

// ContractEventName returns the user-defined event name.
func (TemplateRegistryV1TemplateProposed) ContractEventName() string {
	return TemplateRegistryV1TemplateProposedEventName
}

// UnpackTemplateProposedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TemplateProposed(bytes32 indexed key, bytes32 bytecodeHash, string signature, bytes4 selector, address proposer)
func (templateRegistryV1 *TemplateRegistryV1) UnpackTemplateProposedEvent(log *types.Log) (*TemplateRegistryV1TemplateProposed, error) {
	event := "TemplateProposed"
	if log.Topics[0] != templateRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TemplateRegistryV1TemplateProposed)
	if len(log.Data) > 0 {
		if err := templateRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range templateRegistryV1.abi.Events[event].Inputs {
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

// TemplateRegistryV1TemplateRevoked represents a TemplateRevoked event raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1TemplateRevoked struct {
	Key          [32]byte
	BytecodeHash [32]byte
	Selector     [4]byte
	RevokedAt    uint64
	Raw          *types.Log // Blockchain specific contextual infos
}

const TemplateRegistryV1TemplateRevokedEventName = "TemplateRevoked"

// ContractEventName returns the user-defined event name.
func (TemplateRegistryV1TemplateRevoked) ContractEventName() string {
	return TemplateRegistryV1TemplateRevokedEventName
}

// UnpackTemplateRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TemplateRevoked(bytes32 indexed key, bytes32 bytecodeHash, bytes4 selector, uint64 revokedAt)
func (templateRegistryV1 *TemplateRegistryV1) UnpackTemplateRevokedEvent(log *types.Log) (*TemplateRegistryV1TemplateRevoked, error) {
	event := "TemplateRevoked"
	if log.Topics[0] != templateRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TemplateRegistryV1TemplateRevoked)
	if len(log.Data) > 0 {
		if err := templateRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range templateRegistryV1.abi.Events[event].Inputs {
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

// TemplateRegistryV1Upgraded represents a Upgraded event raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1Upgraded struct {
	Implementation common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const TemplateRegistryV1UpgradedEventName = "Upgraded"

// ContractEventName returns the user-defined event name.
func (TemplateRegistryV1Upgraded) ContractEventName() string {
	return TemplateRegistryV1UpgradedEventName
}

// UnpackUpgradedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Upgraded(address indexed implementation)
func (templateRegistryV1 *TemplateRegistryV1) UnpackUpgradedEvent(log *types.Log) (*TemplateRegistryV1Upgraded, error) {
	event := "Upgraded"
	if log.Topics[0] != templateRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TemplateRegistryV1Upgraded)
	if len(log.Data) > 0 {
		if err := templateRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range templateRegistryV1.abi.Events[event].Inputs {
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
func (templateRegistryV1 *TemplateRegistryV1) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], templateRegistryV1.abi.Errors["AddressEmptyCode"].ID.Bytes()[:4]) {
		return templateRegistryV1.UnpackAddressEmptyCodeError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryV1.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return templateRegistryV1.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryV1.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return templateRegistryV1.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryV1.abi.Errors["FailedCall"].ID.Bytes()[:4]) {
		return templateRegistryV1.UnpackFailedCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryV1.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return templateRegistryV1.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryV1.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return templateRegistryV1.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryV1.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return templateRegistryV1.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryV1.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return templateRegistryV1.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryV1.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return templateRegistryV1.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryV1.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return templateRegistryV1.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryV1.abi.Errors["TemplateRegistryV1AlreadyApproved"].ID.Bytes()[:4]) {
		return templateRegistryV1.UnpackTemplateRegistryV1AlreadyApprovedError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryV1.abi.Errors["TemplateRegistryV1AlreadyRegistered"].ID.Bytes()[:4]) {
		return templateRegistryV1.UnpackTemplateRegistryV1AlreadyRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryV1.abi.Errors["TemplateRegistryV1EmptyBytecodeHash"].ID.Bytes()[:4]) {
		return templateRegistryV1.UnpackTemplateRegistryV1EmptyBytecodeHashError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryV1.abi.Errors["TemplateRegistryV1EmptySignature"].ID.Bytes()[:4]) {
		return templateRegistryV1.UnpackTemplateRegistryV1EmptySignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryV1.abi.Errors["TemplateRegistryV1NotApproved"].ID.Bytes()[:4]) {
		return templateRegistryV1.UnpackTemplateRegistryV1NotApprovedError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryV1.abi.Errors["TemplateRegistryV1NotProposed"].ID.Bytes()[:4]) {
		return templateRegistryV1.UnpackTemplateRegistryV1NotProposedError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryV1.abi.Errors["TemplateRegistryV1RevokedTemplate"].ID.Bytes()[:4]) {
		return templateRegistryV1.UnpackTemplateRegistryV1RevokedTemplateError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryV1.abi.Errors["UUPSUnauthorizedCallContext"].ID.Bytes()[:4]) {
		return templateRegistryV1.UnpackUUPSUnauthorizedCallContextError(raw[4:])
	}
	if bytes.Equal(raw[:4], templateRegistryV1.abi.Errors["UUPSUnsupportedProxiableUUID"].ID.Bytes()[:4]) {
		return templateRegistryV1.UnpackUUPSUnsupportedProxiableUUIDError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// TemplateRegistryV1AddressEmptyCode represents a AddressEmptyCode error raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1AddressEmptyCode struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressEmptyCode(address target)
func TemplateRegistryV1AddressEmptyCodeErrorID() common.Hash {
	return common.HexToHash("0x9996b315c842ff135b8fc4a08ad5df1c344efbc03d2687aecc0678050d2aac89")
}

// UnpackAddressEmptyCodeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressEmptyCode(address target)
func (templateRegistryV1 *TemplateRegistryV1) UnpackAddressEmptyCodeError(raw []byte) (*TemplateRegistryV1AddressEmptyCode, error) {
	out := new(TemplateRegistryV1AddressEmptyCode)
	if err := templateRegistryV1.abi.UnpackIntoInterface(out, "AddressEmptyCode", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryV1ERC1967InvalidImplementation represents a ERC1967InvalidImplementation error raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1ERC1967InvalidImplementation struct {
	Implementation common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func TemplateRegistryV1ERC1967InvalidImplementationErrorID() common.Hash {
	return common.HexToHash("0x4c9c8ce3ceb3130f17f7cdba48d89b5b0129f266a8bac114e6e315a41879b617")
}

// UnpackERC1967InvalidImplementationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func (templateRegistryV1 *TemplateRegistryV1) UnpackERC1967InvalidImplementationError(raw []byte) (*TemplateRegistryV1ERC1967InvalidImplementation, error) {
	out := new(TemplateRegistryV1ERC1967InvalidImplementation)
	if err := templateRegistryV1.abi.UnpackIntoInterface(out, "ERC1967InvalidImplementation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryV1ERC1967NonPayable represents a ERC1967NonPayable error raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1ERC1967NonPayable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967NonPayable()
func TemplateRegistryV1ERC1967NonPayableErrorID() common.Hash {
	return common.HexToHash("0xb398979fa84f543c8e222f17890372c487baf85e062276c127fef521eea7224b")
}

// UnpackERC1967NonPayableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967NonPayable()
func (templateRegistryV1 *TemplateRegistryV1) UnpackERC1967NonPayableError(raw []byte) (*TemplateRegistryV1ERC1967NonPayable, error) {
	out := new(TemplateRegistryV1ERC1967NonPayable)
	if err := templateRegistryV1.abi.UnpackIntoInterface(out, "ERC1967NonPayable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryV1FailedCall represents a FailedCall error raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1FailedCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedCall()
func TemplateRegistryV1FailedCallErrorID() common.Hash {
	return common.HexToHash("0xd6bda27508c0fb6d8a39b4b122878dab26f731a7d4e4abe711dd3731899052a4")
}

// UnpackFailedCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedCall()
func (templateRegistryV1 *TemplateRegistryV1) UnpackFailedCallError(raw []byte) (*TemplateRegistryV1FailedCall, error) {
	out := new(TemplateRegistryV1FailedCall)
	if err := templateRegistryV1.abi.UnpackIntoInterface(out, "FailedCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryV1InvalidInitialization represents a InvalidInitialization error raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1InvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func TemplateRegistryV1InvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (templateRegistryV1 *TemplateRegistryV1) UnpackInvalidInitializationError(raw []byte) (*TemplateRegistryV1InvalidInitialization, error) {
	out := new(TemplateRegistryV1InvalidInitialization)
	if err := templateRegistryV1.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryV1NotInitializing represents a NotInitializing error raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1NotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func TemplateRegistryV1NotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (templateRegistryV1 *TemplateRegistryV1) UnpackNotInitializingError(raw []byte) (*TemplateRegistryV1NotInitializing, error) {
	out := new(TemplateRegistryV1NotInitializing)
	if err := templateRegistryV1.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryV1RaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1RaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func TemplateRegistryV1RaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (templateRegistryV1 *TemplateRegistryV1) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*TemplateRegistryV1RaylsAccessManagedContractPaused, error) {
	out := new(TemplateRegistryV1RaylsAccessManagedContractPaused)
	if err := templateRegistryV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryV1RaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1RaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func TemplateRegistryV1RaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (templateRegistryV1 *TemplateRegistryV1) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*TemplateRegistryV1RaylsAccessManagedInvalidAuthority, error) {
	out := new(TemplateRegistryV1RaylsAccessManagedInvalidAuthority)
	if err := templateRegistryV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryV1RaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1RaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func TemplateRegistryV1RaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (templateRegistryV1 *TemplateRegistryV1) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*TemplateRegistryV1RaylsAccessManagedMustSchedule, error) {
	out := new(TemplateRegistryV1RaylsAccessManagedMustSchedule)
	if err := templateRegistryV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryV1RaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1RaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func TemplateRegistryV1RaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (templateRegistryV1 *TemplateRegistryV1) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*TemplateRegistryV1RaylsAccessManagedUnauthorized, error) {
	out := new(TemplateRegistryV1RaylsAccessManagedUnauthorized)
	if err := templateRegistryV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryV1TemplateRegistryV1AlreadyApproved represents a TemplateRegistryV1__AlreadyApproved error raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1TemplateRegistryV1AlreadyApproved struct {
	Key [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TemplateRegistryV1__AlreadyApproved(bytes32 key)
func TemplateRegistryV1TemplateRegistryV1AlreadyApprovedErrorID() common.Hash {
	return common.HexToHash("0x8b5995322cda1cd8981e7c3159c49dc18dcbcb36140f69047217fbbbefabaeb7")
}

// UnpackTemplateRegistryV1AlreadyApprovedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TemplateRegistryV1__AlreadyApproved(bytes32 key)
func (templateRegistryV1 *TemplateRegistryV1) UnpackTemplateRegistryV1AlreadyApprovedError(raw []byte) (*TemplateRegistryV1TemplateRegistryV1AlreadyApproved, error) {
	out := new(TemplateRegistryV1TemplateRegistryV1AlreadyApproved)
	if err := templateRegistryV1.abi.UnpackIntoInterface(out, "TemplateRegistryV1AlreadyApproved", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryV1TemplateRegistryV1AlreadyRegistered represents a TemplateRegistryV1__AlreadyRegistered error raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1TemplateRegistryV1AlreadyRegistered struct {
	Key [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TemplateRegistryV1__AlreadyRegistered(bytes32 key)
func TemplateRegistryV1TemplateRegistryV1AlreadyRegisteredErrorID() common.Hash {
	return common.HexToHash("0xb60533b993c6aa535ec473363d089400c5dc78f190addf85f278b1c89f297219")
}

// UnpackTemplateRegistryV1AlreadyRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TemplateRegistryV1__AlreadyRegistered(bytes32 key)
func (templateRegistryV1 *TemplateRegistryV1) UnpackTemplateRegistryV1AlreadyRegisteredError(raw []byte) (*TemplateRegistryV1TemplateRegistryV1AlreadyRegistered, error) {
	out := new(TemplateRegistryV1TemplateRegistryV1AlreadyRegistered)
	if err := templateRegistryV1.abi.UnpackIntoInterface(out, "TemplateRegistryV1AlreadyRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryV1TemplateRegistryV1EmptyBytecodeHash represents a TemplateRegistryV1__EmptyBytecodeHash error raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1TemplateRegistryV1EmptyBytecodeHash struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TemplateRegistryV1__EmptyBytecodeHash()
func TemplateRegistryV1TemplateRegistryV1EmptyBytecodeHashErrorID() common.Hash {
	return common.HexToHash("0x1b642519aaf02e985483690f45bcb07a900bf77a2e18743cc72bbbdb3ca177a5")
}

// UnpackTemplateRegistryV1EmptyBytecodeHashError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TemplateRegistryV1__EmptyBytecodeHash()
func (templateRegistryV1 *TemplateRegistryV1) UnpackTemplateRegistryV1EmptyBytecodeHashError(raw []byte) (*TemplateRegistryV1TemplateRegistryV1EmptyBytecodeHash, error) {
	out := new(TemplateRegistryV1TemplateRegistryV1EmptyBytecodeHash)
	if err := templateRegistryV1.abi.UnpackIntoInterface(out, "TemplateRegistryV1EmptyBytecodeHash", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryV1TemplateRegistryV1EmptySignature represents a TemplateRegistryV1__EmptySignature error raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1TemplateRegistryV1EmptySignature struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TemplateRegistryV1__EmptySignature()
func TemplateRegistryV1TemplateRegistryV1EmptySignatureErrorID() common.Hash {
	return common.HexToHash("0x849002d2d4f0f160c0d78a1d3adb3c22638449c3d514525a0642c97ba6e43f3d")
}

// UnpackTemplateRegistryV1EmptySignatureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TemplateRegistryV1__EmptySignature()
func (templateRegistryV1 *TemplateRegistryV1) UnpackTemplateRegistryV1EmptySignatureError(raw []byte) (*TemplateRegistryV1TemplateRegistryV1EmptySignature, error) {
	out := new(TemplateRegistryV1TemplateRegistryV1EmptySignature)
	if err := templateRegistryV1.abi.UnpackIntoInterface(out, "TemplateRegistryV1EmptySignature", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryV1TemplateRegistryV1NotApproved represents a TemplateRegistryV1__NotApproved error raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1TemplateRegistryV1NotApproved struct {
	Key [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TemplateRegistryV1__NotApproved(bytes32 key)
func TemplateRegistryV1TemplateRegistryV1NotApprovedErrorID() common.Hash {
	return common.HexToHash("0x45a9a41cdc1a0623a51248048f0d8fd65dbeab7af7e2669679a37486788e5157")
}

// UnpackTemplateRegistryV1NotApprovedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TemplateRegistryV1__NotApproved(bytes32 key)
func (templateRegistryV1 *TemplateRegistryV1) UnpackTemplateRegistryV1NotApprovedError(raw []byte) (*TemplateRegistryV1TemplateRegistryV1NotApproved, error) {
	out := new(TemplateRegistryV1TemplateRegistryV1NotApproved)
	if err := templateRegistryV1.abi.UnpackIntoInterface(out, "TemplateRegistryV1NotApproved", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryV1TemplateRegistryV1NotProposed represents a TemplateRegistryV1__NotProposed error raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1TemplateRegistryV1NotProposed struct {
	Key [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TemplateRegistryV1__NotProposed(bytes32 key)
func TemplateRegistryV1TemplateRegistryV1NotProposedErrorID() common.Hash {
	return common.HexToHash("0x540fa24c7f72a9ed122bc8cd7ee9f2bd026d9eb98271fa4b224453c1b92494b2")
}

// UnpackTemplateRegistryV1NotProposedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TemplateRegistryV1__NotProposed(bytes32 key)
func (templateRegistryV1 *TemplateRegistryV1) UnpackTemplateRegistryV1NotProposedError(raw []byte) (*TemplateRegistryV1TemplateRegistryV1NotProposed, error) {
	out := new(TemplateRegistryV1TemplateRegistryV1NotProposed)
	if err := templateRegistryV1.abi.UnpackIntoInterface(out, "TemplateRegistryV1NotProposed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryV1TemplateRegistryV1RevokedTemplate represents a TemplateRegistryV1__RevokedTemplate error raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1TemplateRegistryV1RevokedTemplate struct {
	Key [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TemplateRegistryV1__RevokedTemplate(bytes32 key)
func TemplateRegistryV1TemplateRegistryV1RevokedTemplateErrorID() common.Hash {
	return common.HexToHash("0x7151679b2b0e4b55508f437bdfcab632e3fe1e417f1a42bd6696baecbf4f9a9e")
}

// UnpackTemplateRegistryV1RevokedTemplateError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TemplateRegistryV1__RevokedTemplate(bytes32 key)
func (templateRegistryV1 *TemplateRegistryV1) UnpackTemplateRegistryV1RevokedTemplateError(raw []byte) (*TemplateRegistryV1TemplateRegistryV1RevokedTemplate, error) {
	out := new(TemplateRegistryV1TemplateRegistryV1RevokedTemplate)
	if err := templateRegistryV1.abi.UnpackIntoInterface(out, "TemplateRegistryV1RevokedTemplate", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryV1UUPSUnauthorizedCallContext represents a UUPSUnauthorizedCallContext error raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1UUPSUnauthorizedCallContext struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnauthorizedCallContext()
func TemplateRegistryV1UUPSUnauthorizedCallContextErrorID() common.Hash {
	return common.HexToHash("0xe07c8dba242a06571ac65fe4bbe20522c9fb111cb33599b799ff8039c1ed18f4")
}

// UnpackUUPSUnauthorizedCallContextError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnauthorizedCallContext()
func (templateRegistryV1 *TemplateRegistryV1) UnpackUUPSUnauthorizedCallContextError(raw []byte) (*TemplateRegistryV1UUPSUnauthorizedCallContext, error) {
	out := new(TemplateRegistryV1UUPSUnauthorizedCallContext)
	if err := templateRegistryV1.abi.UnpackIntoInterface(out, "UUPSUnauthorizedCallContext", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateRegistryV1UUPSUnsupportedProxiableUUID represents a UUPSUnsupportedProxiableUUID error raised by the TemplateRegistryV1 contract.
type TemplateRegistryV1UUPSUnsupportedProxiableUUID struct {
	Slot [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func TemplateRegistryV1UUPSUnsupportedProxiableUUIDErrorID() common.Hash {
	return common.HexToHash("0xaa1d49a4c084bfa9aeeee2a0be65267a7f19ba7e1476b114dac513d2c14cb563")
}

// UnpackUUPSUnsupportedProxiableUUIDError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func (templateRegistryV1 *TemplateRegistryV1) UnpackUUPSUnsupportedProxiableUUIDError(raw []byte) (*TemplateRegistryV1UUPSUnsupportedProxiableUUID, error) {
	out := new(TemplateRegistryV1UUPSUnsupportedProxiableUUID)
	if err := templateRegistryV1.abi.UnpackIntoInterface(out, "UUPSUnsupportedProxiableUUID", raw); err != nil {
		return nil, err
	}
	return out, nil
}
