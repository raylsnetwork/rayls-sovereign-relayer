// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ParticipantStorageV1

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

// ParticipantStructsAuditInfoData is an auto generated low-level Go binding around an user-defined struct.
type ParticipantStructsAuditInfoData struct {
	ChainId                      *big.Int
	RaylsViewPublicKey           string
	EncryptedRaylsViewPrivateKey []byte
	Mac                          []byte
	BlockNumber                  *big.Int
}

// ParticipantStructsKeyAgreementData is an auto generated low-level Go binding around an user-defined struct.
type ParticipantStructsKeyAgreementData struct {
	ChainId     *big.Int
	Ciphertext  []byte
	Digest      []byte
	BlockNumber *big.Int
}

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

// ParticipantStructsParticipantData is an auto generated low-level Go binding around an user-defined struct.
type ParticipantStructsParticipantData struct {
	ChainId            *big.Int
	Role               uint8
	OwnerId            string
	Name               string
	AllowedToBroadcast bool
}

// ParticipantStructsPrivacyNodeSpendDataSafeReturn is an auto generated low-level Go binding around an user-defined struct.
type ParticipantStructsPrivacyNodeSpendDataSafeReturn struct {
	PaymentSpendPublicKey *big.Int
	PnAddresses           []common.Address
	ChainId               *big.Int
}

// ParticipantStructsPrivacyNodeViewData is an auto generated low-level Go binding around an user-defined struct.
type ParticipantStructsPrivacyNodeViewData struct {
	ChainId            *big.Int
	RaylsViewPublicKey string
	BlockNumber        *big.Int
}

// ParticipantStorageV1MetaData contains all meta data concerning the ParticipantStorageV1 contract.
var ParticipantStorageV1MetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addParticipant\",\"inputs\":[{\"name\":\"_participant\",\"type\":\"tuple\",\"internalType\":\"structParticipantStructs.ParticipantData\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"role\",\"type\":\"uint8\",\"internalType\":\"enumParticipantStructs.Role\"},{\"name\":\"ownerId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"allowedToBroadcast\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addParticipants\",\"inputs\":[{\"name\":\"_participants\",\"type\":\"tuple[]\",\"internalType\":\"structParticipantStructs.ParticipantData[]\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"role\",\"type\":\"uint8\",\"internalType\":\"enumParticipantStructs.Role\"},{\"name\":\"ownerId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"allowedToBroadcast\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"broadcastCurrentParticipants\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkEnygmaAccountAllowed\",\"inputs\":[{\"name\":\"_address\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkEnygmaIssuerAccountAllowed\",\"inputs\":[{\"name\":\"_address\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"configureModules\",\"inputs\":[{\"name\":\"_participantCore\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_auditManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_enygmaManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"contractVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"endpoint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRaylsEndpoint\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAddressByResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllParticipants\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structParticipantStructs.Participant[]\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"role\",\"type\":\"uint8\",\"internalType\":\"enumParticipantStructs.Role\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumParticipantStructs.Status\"},{\"name\":\"ownerId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"allowedToBroadcast\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllParticipantsChainIds\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllPaymentSpendPublicKeys\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structParticipantStructs.PrivacyNodeSpendDataSafeReturn[]\",\"components\":[{\"name\":\"paymentSpendPublicKey\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pnAddresses\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllPrivacyNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAuditInfo\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"data\",\"type\":\"tuple[]\",\"internalType\":\"structParticipantStructs.AuditInfoData[]\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"raylsViewPublicKey\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encryptedRaylsViewPrivateKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"mac\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAuditManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIAuditManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getChainViewData\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"data\",\"type\":\"tuple[]\",\"internalType\":\"structParticipantStructs.PrivacyNodeViewData[]\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"raylsViewPublicKey\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEnygmaAllParticipantsChainIds\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEnygmaManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEnygmaManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getKeyAgreements\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structParticipantStructs.KeyAgreementData[]\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"ciphertext\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"digest\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getParticipant\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structParticipantStructs.Participant\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"role\",\"type\":\"uint8\",\"internalType\":\"enumParticipantStructs.Role\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumParticipantStructs.Status\"},{\"name\":\"ownerId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"allowedToBroadcast\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getParticipantCore\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIParticipantCore\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getParticipantDataBatch\",\"inputs\":[],\"outputs\":[{\"name\":\"pnViewData\",\"type\":\"tuple[]\",\"internalType\":\"structParticipantStructs.PrivacyNodeViewData[]\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"raylsViewPublicKey\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"auditInfo\",\"type\":\"tuple[]\",\"internalType\":\"structParticipantStructs.AuditInfoData[]\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"raylsViewPublicKey\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encryptedRaylsViewPrivateKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"mac\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pnChainIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"auditChainIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPaymentSpendPublicKeyByChainId\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"authority_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initiateKeyAgreement\",\"inputs\":[{\"name\":\"initiatorChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"responderChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"ciphertext\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"digest\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeParticipant\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resourceId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setAuditInfo\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"raylsViewPublicKey\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encryptedRaylsViewPrivateKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"mac\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAuditManager\",\"inputs\":[{\"name\":\"_auditManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setChainViewData\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"raylsViewPublicKey\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setEnygmaManager\",\"inputs\":[{\"name\":\"_enygmaManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setEnygmaPnEventsAddress\",\"inputs\":[{\"name\":\"_pnEnygmaEvents\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setParticipantCore\",\"inputs\":[{\"name\":\"_participantCore\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPaymentSpendPublicKey\",\"inputs\":[{\"name\":\"_chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_paymentSpendPublicKey\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_pnAddresses\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateBroadcastMessagesPermission\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"allowed\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateRole\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"role\",\"type\":\"uint8\",\"internalType\":\"enumParticipantStructs.Role\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateStatus\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumParticipantStructs.Status\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"validateMessageParticipants\",\"inputs\":[{\"name\":\"originChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validateParticipantStatus\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyParticipant\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ModulesConfigured\",\"inputs\":[{\"name\":\"participantCore\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"auditManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"enygmaManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParticipantStorageV1__UnauthorizedCaller\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	ID:  "ParticipantStorageV1",
	Bin: "0x60a06040523060805234801561001457600080fd5b50608051613c4061003e60003960008181611df801528181611e210152611f580152613c406000f3fe6080604052600436106102325760003560e01c80635e280f111161012f578063ad3cb1cc116100b1578063ad3cb1cc1461066a578063b3816fdb146106a8578063bf7e214f146106c8578063c4d66de8146106dd578063c6885dc7146106fd578063c68bab0f1461071f578063c9557f561461073f578063d5c3614f1461075f578063d6b81d5e1461077f578063df5cc3351461079d578063e2446536146107ca578063ede591fd146107ea57600080fd5b80635e280f111461051d5780635f997c5b1461053d5780636279aa9f14610553578063628de2651461057357806365ffa17e14610593578063683f7f27146105b357806381a40de6146105d35780638c8e0c2a146105f3578063987345cb14610613578063a0842e3314610631578063a0a8e4601461065657600080fd5b806337e0d45d116101b857806337e0d45d146103cf5780633a1b3d31146103e45780634017734d14610404578063485cc955146104245780634b9bd78c146104445780634cc168c9146104745780634ec53b2e146104895780634f1ef286146104a957806352bf1c8d146104bc57806352d1902d146104dc578063594c5c76146104ff57600080fd5b8063035c87131461023757806306ec249d14610259578063079c2a491461026e5780630d30d0921461028e5780630fd44407146102c457806311f50c85146102f1578063195ec9ee1461031e5780631b9db2ef146103405780632d3358751461036d57806331b1125f1461038f57806333ded4a8146103af575b600080fd5b34801561024357600080fd5b506102576102523660046122b4565b61080a565b005b34801561026557600080fd5b50610257610871565b34801561027a57600080fd5b506102576102893660046122d1565b61091c565b34801561029a57600080fd5b506102ae6102a9366004612353565b6109c7565b6040516102bb91906123bc565b60405180910390f35b3480156102d057600080fd5b506102e46102df366004612353565b610a69565b6040516102bb91906124d3565b3480156102fd57600080fd5b5061031161030c366004612353565b610b05565b6040516102bb91906124e6565b34801561032a57600080fd5b50610333610b73565b6040516102bb91906125c1565b34801561034c57600080fd5b5061036061035b366004612353565b610c1e565b6040516102bb9190612625565b34801561037957600080fd5b50610382610d02565b6040516102bb9190612674565b34801561039b57600080fd5b506102576103aa366004612687565b610d9c565b3480156103bb57600080fd5b506102576103ca3660046126e0565b610e92565b3480156103db57600080fd5b50610382610f3b565b3480156103f057600080fd5b506102576103ff36600461271d565b610fb9565b34801561041057600080fd5b5061025761041f36600461278a565b611029565b34801561043057600080fd5b5061025761043f3660046127dc565b61109d565b34801561045057600080fd5b5061046461045f3660046122b4565b6111ad565b60405190151581526020016102bb565b34801561048057600080fd5b50610382611248565b34801561049557600080fd5b506102576104a436600461280a565b6112c6565b6102576104b73660046129e7565b61137b565b3480156104c857600080fd5b506102576104d73660046122b4565b61139a565b3480156104e857600080fd5b506104f1611407565b6040519081526020016102bb565b34801561050b57600080fd5b506002546001600160a01b0316610311565b34801561052957600080fd5b50600054610311906001600160a01b031681565b34801561054957600080fd5b506104f160015481565b34801561055f57600080fd5b5061025761056e366004612a36565b611425565b34801561057f57600080fd5b5061025761058e366004612ba3565b6114d7565b34801561059f57600080fd5b506102576105ae366004612c53565b611545565b3480156105bf57600080fd5b506102576105ce366004612353565b6115b5565b3480156105df57600080fd5b506104f16105ee366004612353565b611624565b3480156105ff57600080fd5b5061046461060e366004612c78565b6116ba565b34801561061f57600080fd5b506003546001600160a01b0316610311565b34801561063d57600080fd5b50610646611763565b6040516102bb9493929190612d46565b34801561066257600080fd5b5060016104f1565b34801561067657600080fd5b5061069b604051806040016040528060058152602001640352e302e360dc1b81525081565b6040516102bb9190612d9e565b3480156106b457600080fd5b506104646106c3366004612353565b61181d565b3480156106d457600080fd5b50610311611879565b3480156106e957600080fd5b506102576106f83660046122b4565b611892565b34801561070957600080fd5b506107126118bc565b6040516102bb9190612db1565b34801561072b57600080fd5b5061025761073a3660046122b4565b611956565b34801561074b57600080fd5b5061025761075a366004612353565b6119b4565b34801561076b57600080fd5b5061025761077a366004612e67565b611a34565b34801561078b57600080fd5b506004546001600160a01b0316610311565b3480156107a957600080fd5b506107bd6107b8366004612353565b611abb565b6040516102bb9190612e89565b3480156107d657600080fd5b506102576107e5366004612e9c565b611b57565b3480156107f657600080fd5b506102576108053660046122b4565b611bc5565b610820336000356001600160e01b031916611c23565b6001600160a01b03811661084f5760405162461bcd60e51b815260040161084690612ed8565b60405180910390fd5b600280546001600160a01b0319166001600160a01b0392909216919091179055565b610887336000356001600160e01b031916611c23565b6002546001600160a01b03166108af5760405162461bcd60e51b815260040161084690612f35565b60006108b9611d65565b60025460405163e91638eb60e01b8152600481018390529192506001600160a01b03169063e91638eb906024015b600060405180830381600087803b15801561090157600080fd5b505af1158015610915573d6000803e3d6000fd5b5050505050565b610932336000356001600160e01b031916611c23565b6004546001600160a01b031661095a5760405162461bcd60e51b815260040161084690612f89565b6004805460405163079c2a4960e01b81526001600160a01b039091169163079c2a499161098f91889188918891889101612fdb565b600060405180830381600087803b1580156109a957600080fd5b505af11580156109bd573d6000803e3d6000fd5b5050505050505050565b6003546060906001600160a01b03166109f25760405162461bcd60e51b81526004016108469061302d565b600354604051630698684960e11b8152600481018490526001600160a01b0390911690630d30d09290602401600060405180830381865afa158015610a3b573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610a6391908101906130c3565b92915050565b6003546060906001600160a01b0316610a945760405162461bcd60e51b81526004016108469061302d565b600354604051630fd4440760e01b8152600481018490526001600160a01b0390911690630fd4440790602401600060405180830381865afa158015610add573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610a6391908101906132af565b600080546040516311f50c8560e01b8152600481018490526001600160a01b03909116906311f50c8590602401602060405180830381865afa158015610b4f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a6391906132e3565b6002546060906001600160a01b0316610b9e5760405162461bcd60e51b815260040161084690612f35565b600260009054906101000a90046001600160a01b03166001600160a01b031663195ec9ee6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610bf1573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610c1991908101906133e6565b905090565b610c6960408051610100810190915260008082526020820190815260200160008152602001606081526020016060815260200160008152602001600081526020016000151581525090565b6002546001600160a01b0316610c915760405162461bcd60e51b815260040161084690612f35565b600254604051631b9db2ef60e01b8152600481018490526001600160a01b0390911690631b9db2ef90602401600060405180830381865afa158015610cda573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610a639190810190613489565b6004546060906001600160a01b0316610d2d5760405162461bcd60e51b815260040161084690612f89565b6004805460408051632d33587560e01b815290516001600160a01b0390921692632d3358759282820192600092908290030181865afa158015610d74573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610c19919081019061351c565b610db2336000356001600160e01b031916611c23565b6001600160a01b038316610dd85760405162461bcd60e51b815260040161084690612ed8565b6001600160a01b038216610dfe5760405162461bcd60e51b815260040161084690613550565b6001600160a01b038116610e245760405162461bcd60e51b8152600401610846906135a9565b600280546001600160a01b03199081166001600160a01b03868116918217909355600380548316868516908117909155600480549093169385169384179092556040517ff6c27e15b7e995e3edd056f3e9b7b01098dfe3f91cccf2af78ff33215fc1829d90600090a4505050565b610ea8336000356001600160e01b031916611c23565b6002546001600160a01b0316610ed05760405162461bcd60e51b815260040161084690612f35565b60025460405163067bda9560e31b81526004810184905282151560248201526001600160a01b03909116906333ded4a8906044015b600060405180830381600087803b158015610f1f57600080fd5b505af1158015610f33573d6000803e3d6000fd5b505050505050565b6002546060906001600160a01b0316610f665760405162461bcd60e51b815260040161084690612f35565b600260009054906101000a90046001600160a01b03166001600160a01b03166337e0d45d6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610d74573d6000803e3d6000fd5b610fcf336000356001600160e01b031916611c23565b6002546001600160a01b0316610ff75760405162461bcd60e51b815260040161084690612f35565b600254604051633a1b3d3160e01b81526001600160a01b0390911690633a1b3d3190610f059085908590600401613603565b61103f336000356001600160e01b031916611c23565b6003546001600160a01b03166110675760405162461bcd60e51b81526004016108469061302d565b600354604051634017734d60e01b81526001600160a01b0390911690634017734d9061098f908790879087908790600401613640565b60006110a7611d79565b805490915060ff600160401b82041615906001600160401b03166000811580156110ce5750825b90506000826001600160401b031660011480156110ea5750303b155b9050811580156110f8575080155b156111165760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561114057845460ff60401b1916600160401b1785555b611148611da2565b61115187611892565b6001805561115e86611dac565b83156111a457845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50505050505050565b6004546000906001600160a01b03166111d85760405162461bcd60e51b815260040161084690612f89565b600480546040516312e6f5e360e21b81526001600160a01b0390911691634b9bd78c91611207918691016124e6565b602060405180830381865afa158015611224573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a63919061366b565b6003546060906001600160a01b03166112735760405162461bcd60e51b81526004016108469061302d565b600360009054906101000a90046001600160a01b03166001600160a01b0316634cc168c96040518163ffffffff1660e01b8152600401600060405180830381865afa158015610d74573d6000803e3d6000fd5b6112dc336000356001600160e01b031916611c23565b6003546001600160a01b03166113045760405162461bcd60e51b81526004016108469061302d565b6003546040516327629d9760e11b81526001600160a01b0390911690634ec53b2e90611340908a908a908a908a908a908a908a90600401613688565b600060405180830381600087803b15801561135a57600080fd5b505af115801561136e573d6000803e3d6000fd5b5050505050505050505050565b611383611ded565b61138c82611e7b565b6113968282611e94565b5050565b6113b0336000356001600160e01b031916611c23565b6004546001600160a01b03166113d85760405162461bcd60e51b815260040161084690612f89565b600480546040516352bf1c8d60e01b81526001600160a01b03909116916352bf1c8d916108e7918591016124e6565b6000611411611f4d565b50600080516020613beb8339815191525b90565b61143b336000356001600160e01b031916611c23565b6003546001600160a01b03166114635760405162461bcd60e51b81526004016108469061302d565b600354604051636279aa9f60e01b81526001600160a01b0390911690636279aa9f9061149d908990899089908990899089906004016136d0565b600060405180830381600087803b1580156114b757600080fd5b505af11580156114cb573d6000803e3d6000fd5b50505050505050505050565b6114ed336000356001600160e01b031916611c23565b6002546001600160a01b03166115155760405162461bcd60e51b815260040161084690612f35565b60025460405163628de26560e01b81526001600160a01b039091169063628de265906108e7908490600401613788565b61155b336000356001600160e01b031916611c23565b6002546001600160a01b03166115835760405162461bcd60e51b815260040161084690612f35565b6002546040516332ffd0bf60e11b81526001600160a01b03909116906365ffa17e90610f0590859085906004016137df565b6115cb336000356001600160e01b031916611c23565b6002546001600160a01b03166115f35760405162461bcd60e51b815260040161084690612f35565b60025460405163683f7f2760e01b8152600481018390526001600160a01b039091169063683f7f27906024016108e7565b6004546000906001600160a01b031661164f5760405162461bcd60e51b815260040161084690612f89565b600480546040516340d206f360e11b81529182018490526001600160a01b0316906381a40de690602401602060405180830381865afa158015611696573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a6391906137f3565b6004546000906001600160a01b03166116e55760405162461bcd60e51b815260040161084690612f89565b60048054604051634647061560e11b81526001600160a01b038681169382019390935260248101859052911690638c8e0c2a90604401602060405180830381865afa158015611738573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061175c919061366b565b9392505050565b6003546060908190819081906001600160a01b03166117945760405162461bcd60e51b81526004016108469061302d565b600360009054906101000a90046001600160a01b03166001600160a01b031663a0842e336040518163ffffffff1660e01b8152600401600060405180830381865afa1580156117e7573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261180f9190810190613929565b935093509350935090919293565b6002546000906001600160a01b03166118485760405162461bcd60e51b815260040161084690612f35565b60025460405163b3816fdb60e01b8152600481018490526001600160a01b039091169063b3816fdb90602401611207565b6000611883611f96565b546001600160a01b0316919050565b61189a611ff8565b600080546001600160a01b0319166001600160a01b0392909216919091179055565b6004546060906001600160a01b03166118e75760405162461bcd60e51b815260040161084690612f89565b600480546040805163c6885dc760e01b815290516001600160a01b039092169263c6885dc79282820192600092908290030181865afa15801561192e573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610c1991908101906139d5565b61196c336000356001600160e01b031916611c23565b6001600160a01b0381166119925760405162461bcd60e51b815260040161084690613550565b600380546001600160a01b0319166001600160a01b0392909216919091179055565b6002546001600160a01b03166119dc5760405162461bcd60e51b815260040161084690612f35565b6002546040516364aabfab60e11b8152600481018390526001600160a01b039091169063c9557f569060240160006040518083038186803b158015611a2057600080fd5b505afa158015610915573d6000803e3d6000fd5b6002546001600160a01b0316611a5c5760405162461bcd60e51b815260040161084690612f35565b60025460405163d5c3614f60e01b815260048101849052602481018390526001600160a01b039091169063d5c3614f9060440160006040518083038186803b158015611aa757600080fd5b505afa158015610f33573d6000803e3d6000fd5b6003546060906001600160a01b0316611ae65760405162461bcd60e51b81526004016108469061302d565b60035460405163df5cc33560e01b8152600481018490526001600160a01b039091169063df5cc33590602401600060405180830381865afa158015611b2f573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610a639190810190613b1b565b611b6d336000356001600160e01b031916611c23565b6002546001600160a01b0316611b955760405162461bcd60e51b815260040161084690612f35565b600254604051637122329b60e11b81526001600160a01b039091169063e2446536906108e7908490600401613b4f565b611bdb336000356001600160e01b031916611c23565b6001600160a01b038116611c015760405162461bcd60e51b8152600401610846906135a9565b600480546001600160a01b0319166001600160a01b0392909216919091179055565b6000611c2d611f96565b80549091506001600160a01b031680611c5c576000604051638944034760e01b815260040161084691906124e6565b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa158015611cc0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611ce49190613b62565b925092509250826111a4578015611d0e5760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff821615611d4a5760405163a426878960e01b81526001600160a01b038816600482015263ffffffff83166024820152604401610846565b86604051632ecd3d0360e21b815260040161084691906124e6565b600060343610611422575060331936013590565b6000807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00610a63565b611daa611ff8565b565b6000611db6611f96565b80549091506001600160a01b031615611de45781604051638944034760e01b815260040161084691906124e6565b6113968261201d565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161480611e5d57507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316611e516120ad565b6001600160a01b031614155b15611daa5760405163703e46dd60e11b815260040160405180910390fd5b611e91336000356001600160e01b031916611c23565b50565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015611eee575060408051601f3d908101601f19168201909252611eeb918101906137f3565b60015b611f0d5781604051634c9c8ce360e01b815260040161084691906124e6565b600080516020613beb8339815191528114611f3e57604051632a87526960e21b815260048101829052602401610846565b611f4883836120c3565b505050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614611daa5760405163703e46dd60e11b815260040160405180910390fd5b60008060ff19611fc760017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f35613bad565b604051602001611fd991815260200190565b60408051601f1981840301815291905280516020909101201692915050565b612000612119565b611daa57604051631afcd79f60e31b815260040160405180910390fd5b6001600160a01b0381166120465780604051638944034760e01b815260040161084691906124e6565b6000612050611f96565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b6000600080516020613beb833981519152611883565b6120cc82612133565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a280511561211157611f48828261218f565b611396612205565b6000612123611d79565b54600160401b900460ff16919050565b806001600160a01b03163b6000036121605780604051634c9c8ce360e01b815260040161084691906124e6565b600080516020613beb83398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080846001600160a01b0316846040516121ac9190613bce565b600060405180830381855af49150503d80600081146121e7576040519150601f19603f3d011682016040523d82523d6000602084013e6121ec565b606091505b50915091506121fc858383612224565b95945050505050565b3415611daa5760405163b398979f60e01b815260040160405180910390fd5b6060826122395761223482612277565b61175c565b815115801561225057506001600160a01b0384163b155b156122705783604051639996b31560e01b815260040161084691906124e6565b5092915050565b80511561228657805160208201fd5b60405163d6bda27560e01b815260040160405180910390fd5b6001600160a01b0381168114611e9157600080fd5b6000602082840312156122c657600080fd5b813561175c8161229f565b600080600080606085870312156122e757600080fd5b843593506020850135925060408501356001600160401b038082111561230c57600080fd5b818701915087601f83011261232057600080fd5b81358181111561232f57600080fd5b8860208260051b850101111561234457600080fd5b95989497505060200194505050565b60006020828403121561236557600080fd5b5035919050565b60005b8381101561238757818101518382015260200161236f565b50506000910152565b600081518084526123a881602086016020860161236c565b601f01601f19169290920160200192915050565b600060208083018184528085518083526040925060408601915060408160051b87010184880160005b8381101561244c57603f19898403018552815160808151855288820151818a87015261241382870182612390565b915050878201518582038987015261242b8282612390565b606093840151969093019590955250948701949250908601906001016123e5565b509098975050505050505050565b600082825180855260208086019550808260051b84010181860160005b848110156124c657601f1986840301895281516060815185528582015181878701526124a582870182612390565b60409384015196909301959095525098840198925090830190600101612477565b5090979650505050505050565b60208152600061175c602083018461245a565b6001600160a01b0391909116815260200190565b634e487b7160e01b600052602160045260246000fd5b60038110612520576125206124fa565b9052565b60048110612520576125206124fa565b600061010082518452602083015161254f6020860182612510565b5060408301516125626040860182612524565b50606083015181606086015261257a82860182612390565b915050608083015184820360808601526125948282612390565b91505060a083015160a085015260c083015160c085015260e0830151151560e08501528091505092915050565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b8281101561261857603f19888603018452612606858351612534565b945092850192908501906001016125ea565b5092979650505050505050565b60208152600061175c6020830184612534565b60008151808452602080850194506020840160005b838110156126695781518752958201959082019060010161264d565b509495945050505050565b60208152600061175c6020830184612638565b60008060006060848603121561269c57600080fd5b83356126a78161229f565b925060208401356126b78161229f565b915060408401356126c78161229f565b809150509250925092565b8015158114611e9157600080fd5b600080604083850312156126f357600080fd5b823591506020830135612705816126d2565b809150509250929050565b60048110611e9157600080fd5b6000806040838503121561273057600080fd5b82359150602083013561270581612710565b60008083601f84011261275457600080fd5b5081356001600160401b0381111561276b57600080fd5b60208301915083602082850101111561278357600080fd5b9250929050565b600080600080606085870312156127a057600080fd5b8435935060208501356001600160401b038111156127bd57600080fd5b6127c987828801612742565b9598909750949560400135949350505050565b600080604083850312156127ef57600080fd5b82356127fa8161229f565b915060208301356127058161229f565b600080600080600080600060a0888a03121561282557600080fd5b873596506020880135955060408801356001600160401b038082111561284a57600080fd5b6128568b838c01612742565b909750955060608a013591508082111561286f57600080fd5b5061287c8a828b01612742565b989b979a50959894979596608090950135949350505050565b634e487b7160e01b600052604160045260246000fd5b60405160a081016001600160401b03811182821017156128cd576128cd612895565b60405290565b604051608081016001600160401b03811182821017156128cd576128cd612895565b604051606081016001600160401b03811182821017156128cd576128cd612895565b60405161010081016001600160401b03811182821017156128cd576128cd612895565b604051601f8201601f191681016001600160401b038111828210171561296257612962612895565b604052919050565b60006001600160401b0382111561298357612983612895565b50601f01601f191660200190565b600082601f8301126129a257600080fd5b81356129b56129b08261296a565b61293a565b8181528460208386010111156129ca57600080fd5b816020850160208301376000918101602001919091529392505050565b600080604083850312156129fa57600080fd5b8235612a058161229f565b915060208301356001600160401b03811115612a2057600080fd5b612a2c85828601612991565b9150509250929050565b60008060008060008060a08789031215612a4f57600080fd5b8635955060208701356001600160401b0380821115612a6d57600080fd5b612a798a838b01612742565b90975095506040890135915080821115612a9257600080fd5b612a9e8a838b01612991565b94506060890135915080821115612ab457600080fd5b50612ac189828a01612991565b925050608087013590509295509295509295565b60006001600160401b03821115612aee57612aee612895565b5060051b60200190565b60038110611e9157600080fd5b600060a08284031215612b1757600080fd5b612b1f6128ab565b9050813581526020820135612b3381612af8565b602082015260408201356001600160401b0380821115612b5257600080fd5b612b5e85838601612991565b60408401526060840135915080821115612b7757600080fd5b50612b8484828501612991565b6060830152506080820135612b98816126d2565b608082015292915050565b60006020808385031215612bb657600080fd5b82356001600160401b0380821115612bcd57600080fd5b818501915085601f830112612be157600080fd5b8135612bef6129b082612ad5565b81815260059190911b83018401908481019088831115612c0e57600080fd5b8585015b83811015612c4657803585811115612c2a5760008081fd5b612c388b89838a0101612b05565b845250918601918601612c12565b5098975050505050505050565b60008060408385031215612c6657600080fd5b82359150602083013561270581612af8565b60008060408385031215612c8b57600080fd5b8235612c968161229f565b946020939093013593505050565b600082825180855260208086019550808260051b84010181860160005b848110156124c657601f19868403018952815160a081518552858201518187870152612cef82870182612390565b91505060408083015186830382880152612d098382612390565b9250505060608083015186830382880152612d248382612390565b6080948501519790940196909652505098840198925090830190600101612cc1565b608081526000612d59608083018761245a565b8281036020840152612d6b8187612ca4565b90508281036040840152612d7f8186612638565b90508281036060840152612d938185612638565b979650505050505050565b60208152600061175c6020830184612390565b600060208083018184528085518083526040925060408601915060408160051b8701018488016000805b84811015612e5857898403603f190186528251805185528881015160608a8701819052815190870181905260808701918b019085905b80821015612e3a5782516001600160a01b03168452928c0192918c019160019190910190612e11565b50505090880151948801949094529487019491870191600101612ddb565b50919998505050505050505050565b60008060408385031215612e7a57600080fd5b50508035926020909101359150565b60208152600061175c6020830184612ca4565b600060208284031215612eae57600080fd5b81356001600160401b03811115612ec457600080fd5b612ed084828501612b05565b949350505050565b6020808252603c908201527f5061727469636970616e7453746f7261676556313a205061727469636970616e60408201527f74436f726520616464726573732063616e6e6f74206265207a65726f00000000606082015260800190565b60208082526034908201527f5061727469636970616e7453746f7261676556313a205061727469636970616e6040820152731d10dbdc99481b5bd91d5b19481b9bdd081cd95d60621b606082015260800190565b60208082526032908201527f5061727469636970616e7453746f7261676556313a20456e79676d614d616e6160408201527119d95c881b5bd91d5b19481b9bdd081cd95d60721b606082015260800190565b84815260208082018590526060604083018190528201839052600090849060808401835b86811015612c465783356130128161229f565b6001600160a01b031682529282019290820190600101612fff565b60208082526031908201527f5061727469636970616e7453746f7261676556313a2041756469744d616e6167604082015270195c881b5bd91d5b19481b9bdd081cd95d607a1b606082015260800190565b600082601f83011261308f57600080fd5b815161309d6129b08261296a565b8181528460208386010111156130b257600080fd5b612ed082602083016020870161236c565b600060208083850312156130d657600080fd5b82516001600160401b03808211156130ed57600080fd5b818501915085601f83011261310157600080fd5b815161310f6129b082612ad5565b81815260059190911b8301840190848101908883111561312e57600080fd5b8585015b83811015612c465780518581111561314957600080fd5b86016080818c03601f190112156131605760008081fd5b6131686128d3565b888201518152604080830151888111156131825760008081fd5b6131908e8c8387010161307e565b8b84015250606080840151898111156131a95760008081fd5b6131b78f8d8388010161307e565b928401929092526080939093015192820192909252845250918601918601613132565b600082601f8301126131eb57600080fd5b815160206131fb6129b083612ad5565b82815260059290921b8401810191818101908684111561321a57600080fd5b8286015b848110156132a45780516001600160401b038082111561323e5760008081fd5b908801906060828b03601f19018113156132585760008081fd5b6132606128f5565b8784015181526040808501518481111561327a5760008081fd5b6132888e8b8389010161307e565b8a8401525091909301519083015250835291830191830161321e565b509695505050505050565b6000602082840312156132c157600080fd5b81516001600160401b038111156132d757600080fd5b612ed0848285016131da565b6000602082840312156132f557600080fd5b815161175c8161229f565b805161330b81612af8565b919050565b805161330b81612710565b805161330b816126d2565b6000610100828403121561333957600080fd5b613341612917565b90508151815261335360208301613300565b602082015261336460408301613310565b604082015260608201516001600160401b038082111561338357600080fd5b61338f8583860161307e565b606084015260808401519150808211156133a857600080fd5b506133b58482850161307e565b60808301525060a082015160a082015260c082015160c08201526133db60e0830161331b565b60e082015292915050565b600060208083850312156133f957600080fd5b82516001600160401b038082111561341057600080fd5b818501915085601f83011261342457600080fd5b81516134326129b082612ad5565b81815260059190911b8301840190848101908883111561345157600080fd5b8585015b83811015612c465780518581111561346d5760008081fd5b61347b8b89838a0101613326565b845250918601918601613455565b60006020828403121561349b57600080fd5b81516001600160401b038111156134b157600080fd5b612ed084828501613326565b600082601f8301126134ce57600080fd5b815160206134de6129b083612ad5565b8083825260208201915060208460051b87010193508684111561350057600080fd5b602086015b848110156132a45780518352918301918301613505565b60006020828403121561352e57600080fd5b81516001600160401b0381111561354457600080fd5b612ed0848285016134bd565b60208082526039908201527f5061727469636970616e7453746f7261676556313a2041756469744d616e6167604082015278657220616464726573732063616e6e6f74206265207a65726f60381b606082015260800190565b6020808252603a908201527f5061727469636970616e7453746f7261676556313a20456e79676d614d616e6160408201527967657220616464726573732063616e6e6f74206265207a65726f60301b606082015260800190565b8281526040810161175c6020830184612524565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b84815260606020820152600061365a606083018587613617565b905082604083015295945050505050565b60006020828403121561367d57600080fd5b815161175c816126d2565b87815286602082015260a0604082015260006136a860a083018789613617565b82810360608401526136bb818688613617565b91505082608083015298975050505050505050565b86815260a0602082015260006136ea60a083018789613617565b82810360408401526136fc8187612390565b905082810360608401526137108186612390565b915050826080830152979650505050505050565b805182526000602082015161373c6020850182612510565b50604082015160a0604085015261375660a0850182612390565b90506060830151848203606086015261376f8282612390565b9150506080830151151560808501528091505092915050565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b8281101561261857603f198886030184526137cd858351613724565b945092850192908501906001016137b1565b8281526040810161175c6020830184612510565b60006020828403121561380557600080fd5b5051919050565b600082601f83011261381d57600080fd5b8151602061382d6129b083612ad5565b82815260059290921b8401810191818101908684111561384c57600080fd5b8286015b848110156132a45780516001600160401b03808211156138705760008081fd5b9088019060a0828b03601f190181131561388a5760008081fd5b6138926128ab565b878401518152604080850151848111156138ac5760008081fd5b6138ba8e8b8389010161307e565b8a84015250606080860151858111156138d35760008081fd5b6138e18f8c838a010161307e565b83850152506080915081860151858111156138fc5760008081fd5b61390a8f8c838a010161307e565b9184019190915250919093015190830152508352918301918301613850565b6000806000806080858703121561393f57600080fd5b84516001600160401b038082111561395657600080fd5b613962888389016131da565b9550602087015191508082111561397857600080fd5b6139848883890161380c565b9450604087015191508082111561399a57600080fd5b6139a6888389016134bd565b935060608701519150808211156139bc57600080fd5b506139c9878288016134bd565b91505092959194509250565b600060208083850312156139e857600080fd5b82516001600160401b03808211156139ff57600080fd5b818501915085601f830112613a1357600080fd5b8151613a216129b082612ad5565b81815260059190911b83018401908481019088831115613a4057600080fd5b8585015b83811015612c4657805185811115613a5b57600080fd5b86016060818c03601f19011215613a7157600080fd5b613a796128f5565b888201518152604082015187811115613a9157600080fd5b8201603f81018d13613aa257600080fd5b89810151613ab26129b082612ad5565b81815260059190911b8201604001908b8101908f831115613ad257600080fd5b6040840193505b82841015613afb578351613aec8161229f565b8252928c0192908c0190613ad9565b848d01525050506060919091015160408201528352918601918601613a44565b600060208284031215613b2d57600080fd5b81516001600160401b03811115613b4357600080fd5b612ed08482850161380c565b60208152600061175c6020830184613724565b600080600060608486031215613b7757600080fd5b8351613b82816126d2565b602085015190935063ffffffff81168114613b9c57600080fd5b60408501519092506126c7816126d2565b81810381811115610a6357634e487b7160e01b600052601160045260246000fd5b60008251613be081846020870161236c565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbca2646970667358221220c9902f0d88aa1d57fad71099142e89df4f892febf0178464310cf5d28528872264736f6c63430008180033",
}

// ParticipantStorageV1 is an auto generated Go binding around an Ethereum contract.
type ParticipantStorageV1 struct {
	abi abi.ABI
}

// NewParticipantStorageV1 creates a new instance of ParticipantStorageV1.
func NewParticipantStorageV1() *ParticipantStorageV1 {
	parsed, err := ParticipantStorageV1MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ParticipantStorageV1{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ParticipantStorageV1) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackUPGRADEINTERFACEVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (participantStorageV1 *ParticipantStorageV1) PackUPGRADEINTERFACEVERSION() []byte {
	enc, err := participantStorageV1.abi.Pack("UPGRADE_INTERFACE_VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUPGRADEINTERFACEVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (participantStorageV1 *ParticipantStorageV1) UnpackUPGRADEINTERFACEVERSION(data []byte) (string, error) {
	out, err := participantStorageV1.abi.Unpack("UPGRADE_INTERFACE_VERSION", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackAddParticipant is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe2446536.
//
// Solidity: function addParticipant((uint256,uint8,string,string,bool) _participant) returns()
func (participantStorageV1 *ParticipantStorageV1) PackAddParticipant(participant ParticipantStructsParticipantData) []byte {
	enc, err := participantStorageV1.abi.Pack("addParticipant", participant)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAddParticipants is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x628de265.
//
// Solidity: function addParticipants((uint256,uint8,string,string,bool)[] _participants) returns()
func (participantStorageV1 *ParticipantStorageV1) PackAddParticipants(participants []ParticipantStructsParticipantData) []byte {
	enc, err := participantStorageV1.abi.Pack("addParticipants", participants)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (participantStorageV1 *ParticipantStorageV1) PackAuthority() []byte {
	enc, err := participantStorageV1.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (participantStorageV1 *ParticipantStorageV1) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := participantStorageV1.abi.Unpack("authority", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackBroadcastCurrentParticipants is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06ec249d.
//
// Solidity: function broadcastCurrentParticipants() returns()
func (participantStorageV1 *ParticipantStorageV1) PackBroadcastCurrentParticipants() []byte {
	enc, err := participantStorageV1.abi.Pack("broadcastCurrentParticipants")
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCheckEnygmaAccountAllowed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4b9bd78c.
//
// Solidity: function checkEnygmaAccountAllowed(address _address) view returns(bool)
func (participantStorageV1 *ParticipantStorageV1) PackCheckEnygmaAccountAllowed(address common.Address) []byte {
	enc, err := participantStorageV1.abi.Pack("checkEnygmaAccountAllowed", address)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackCheckEnygmaAccountAllowed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4b9bd78c.
//
// Solidity: function checkEnygmaAccountAllowed(address _address) view returns(bool)
func (participantStorageV1 *ParticipantStorageV1) UnpackCheckEnygmaAccountAllowed(data []byte) (bool, error) {
	out, err := participantStorageV1.abi.Unpack("checkEnygmaAccountAllowed", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackCheckEnygmaIssuerAccountAllowed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8c8e0c2a.
//
// Solidity: function checkEnygmaIssuerAccountAllowed(address _address, uint256 _chainId) view returns(bool)
func (participantStorageV1 *ParticipantStorageV1) PackCheckEnygmaIssuerAccountAllowed(address common.Address, chainId *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("checkEnygmaIssuerAccountAllowed", address, chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackCheckEnygmaIssuerAccountAllowed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8c8e0c2a.
//
// Solidity: function checkEnygmaIssuerAccountAllowed(address _address, uint256 _chainId) view returns(bool)
func (participantStorageV1 *ParticipantStorageV1) UnpackCheckEnygmaIssuerAccountAllowed(data []byte) (bool, error) {
	out, err := participantStorageV1.abi.Unpack("checkEnygmaIssuerAccountAllowed", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackConfigureModules is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x31b1125f.
//
// Solidity: function configureModules(address _participantCore, address _auditManager, address _enygmaManager) returns()
func (participantStorageV1 *ParticipantStorageV1) PackConfigureModules(participantCore common.Address, auditManager common.Address, enygmaManager common.Address) []byte {
	enc, err := participantStorageV1.abi.Pack("configureModules", participantCore, auditManager, enygmaManager)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackContractVersion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa0a8e460.
//
// Solidity: function contractVersion() pure returns(uint256)
func (participantStorageV1 *ParticipantStorageV1) PackContractVersion() []byte {
	enc, err := participantStorageV1.abi.Pack("contractVersion")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackContractVersion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa0a8e460.
//
// Solidity: function contractVersion() pure returns(uint256)
func (participantStorageV1 *ParticipantStorageV1) UnpackContractVersion(data []byte) (*big.Int, error) {
	out, err := participantStorageV1.abi.Unpack("contractVersion", data)
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
func (participantStorageV1 *ParticipantStorageV1) PackEndpoint() []byte {
	enc, err := participantStorageV1.abi.Pack("endpoint")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackEndpoint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (participantStorageV1 *ParticipantStorageV1) UnpackEndpoint(data []byte) (common.Address, error) {
	out, err := participantStorageV1.abi.Unpack("endpoint", data)
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
func (participantStorageV1 *ParticipantStorageV1) PackGetAddressByResourceId(resourceId [32]byte) []byte {
	enc, err := participantStorageV1.abi.Pack("getAddressByResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAddressByResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (participantStorageV1 *ParticipantStorageV1) UnpackGetAddressByResourceId(data []byte) (common.Address, error) {
	out, err := participantStorageV1.abi.Unpack("getAddressByResourceId", data)
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
func (participantStorageV1 *ParticipantStorageV1) PackGetAllParticipants() []byte {
	enc, err := participantStorageV1.abi.Pack("getAllParticipants")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAllParticipants is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x195ec9ee.
//
// Solidity: function getAllParticipants() view returns((uint256,uint8,uint8,string,string,uint256,uint256,bool)[])
func (participantStorageV1 *ParticipantStorageV1) UnpackGetAllParticipants(data []byte) ([]ParticipantStructsParticipant, error) {
	out, err := participantStorageV1.abi.Unpack("getAllParticipants", data)
	if err != nil {
		return *new([]ParticipantStructsParticipant), err
	}
	out0 := *abi.ConvertType(out[0], new([]ParticipantStructsParticipant)).(*[]ParticipantStructsParticipant)
	return out0, err
}

// PackGetAllParticipantsChainIds is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x37e0d45d.
//
// Solidity: function getAllParticipantsChainIds() view returns(uint256[])
func (participantStorageV1 *ParticipantStorageV1) PackGetAllParticipantsChainIds() []byte {
	enc, err := participantStorageV1.abi.Pack("getAllParticipantsChainIds")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAllParticipantsChainIds is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x37e0d45d.
//
// Solidity: function getAllParticipantsChainIds() view returns(uint256[])
func (participantStorageV1 *ParticipantStorageV1) UnpackGetAllParticipantsChainIds(data []byte) ([]*big.Int, error) {
	out, err := participantStorageV1.abi.Unpack("getAllParticipantsChainIds", data)
	if err != nil {
		return *new([]*big.Int), err
	}
	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	return out0, err
}

// PackGetAllPaymentSpendPublicKeys is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc6885dc7.
//
// Solidity: function getAllPaymentSpendPublicKeys() view returns((uint256,address[],uint256)[])
func (participantStorageV1 *ParticipantStorageV1) PackGetAllPaymentSpendPublicKeys() []byte {
	enc, err := participantStorageV1.abi.Pack("getAllPaymentSpendPublicKeys")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAllPaymentSpendPublicKeys is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc6885dc7.
//
// Solidity: function getAllPaymentSpendPublicKeys() view returns((uint256,address[],uint256)[])
func (participantStorageV1 *ParticipantStorageV1) UnpackGetAllPaymentSpendPublicKeys(data []byte) ([]ParticipantStructsPrivacyNodeSpendDataSafeReturn, error) {
	out, err := participantStorageV1.abi.Unpack("getAllPaymentSpendPublicKeys", data)
	if err != nil {
		return *new([]ParticipantStructsPrivacyNodeSpendDataSafeReturn), err
	}
	out0 := *abi.ConvertType(out[0], new([]ParticipantStructsPrivacyNodeSpendDataSafeReturn)).(*[]ParticipantStructsPrivacyNodeSpendDataSafeReturn)
	return out0, err
}

// PackGetAllPrivacyNodes is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4cc168c9.
//
// Solidity: function getAllPrivacyNodes() view returns(uint256[])
func (participantStorageV1 *ParticipantStorageV1) PackGetAllPrivacyNodes() []byte {
	enc, err := participantStorageV1.abi.Pack("getAllPrivacyNodes")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAllPrivacyNodes is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4cc168c9.
//
// Solidity: function getAllPrivacyNodes() view returns(uint256[])
func (participantStorageV1 *ParticipantStorageV1) UnpackGetAllPrivacyNodes(data []byte) ([]*big.Int, error) {
	out, err := participantStorageV1.abi.Unpack("getAllPrivacyNodes", data)
	if err != nil {
		return *new([]*big.Int), err
	}
	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	return out0, err
}

// PackGetAuditInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdf5cc335.
//
// Solidity: function getAuditInfo(uint256 chainId) view returns((uint256,string,bytes,bytes,uint256)[] data)
func (participantStorageV1 *ParticipantStorageV1) PackGetAuditInfo(chainId *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("getAuditInfo", chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAuditInfo is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdf5cc335.
//
// Solidity: function getAuditInfo(uint256 chainId) view returns((uint256,string,bytes,bytes,uint256)[] data)
func (participantStorageV1 *ParticipantStorageV1) UnpackGetAuditInfo(data []byte) ([]ParticipantStructsAuditInfoData, error) {
	out, err := participantStorageV1.abi.Unpack("getAuditInfo", data)
	if err != nil {
		return *new([]ParticipantStructsAuditInfoData), err
	}
	out0 := *abi.ConvertType(out[0], new([]ParticipantStructsAuditInfoData)).(*[]ParticipantStructsAuditInfoData)
	return out0, err
}

// PackGetAuditManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x987345cb.
//
// Solidity: function getAuditManager() view returns(address)
func (participantStorageV1 *ParticipantStorageV1) PackGetAuditManager() []byte {
	enc, err := participantStorageV1.abi.Pack("getAuditManager")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAuditManager is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x987345cb.
//
// Solidity: function getAuditManager() view returns(address)
func (participantStorageV1 *ParticipantStorageV1) UnpackGetAuditManager(data []byte) (common.Address, error) {
	out, err := participantStorageV1.abi.Unpack("getAuditManager", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetChainViewData is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0fd44407.
//
// Solidity: function getChainViewData(uint256 chainId) view returns((uint256,string,uint256)[] data)
func (participantStorageV1 *ParticipantStorageV1) PackGetChainViewData(chainId *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("getChainViewData", chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetChainViewData is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0fd44407.
//
// Solidity: function getChainViewData(uint256 chainId) view returns((uint256,string,uint256)[] data)
func (participantStorageV1 *ParticipantStorageV1) UnpackGetChainViewData(data []byte) ([]ParticipantStructsPrivacyNodeViewData, error) {
	out, err := participantStorageV1.abi.Unpack("getChainViewData", data)
	if err != nil {
		return *new([]ParticipantStructsPrivacyNodeViewData), err
	}
	out0 := *abi.ConvertType(out[0], new([]ParticipantStructsPrivacyNodeViewData)).(*[]ParticipantStructsPrivacyNodeViewData)
	return out0, err
}

// PackGetEnygmaAllParticipantsChainIds is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2d335875.
//
// Solidity: function getEnygmaAllParticipantsChainIds() view returns(uint256[])
func (participantStorageV1 *ParticipantStorageV1) PackGetEnygmaAllParticipantsChainIds() []byte {
	enc, err := participantStorageV1.abi.Pack("getEnygmaAllParticipantsChainIds")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetEnygmaAllParticipantsChainIds is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2d335875.
//
// Solidity: function getEnygmaAllParticipantsChainIds() view returns(uint256[])
func (participantStorageV1 *ParticipantStorageV1) UnpackGetEnygmaAllParticipantsChainIds(data []byte) ([]*big.Int, error) {
	out, err := participantStorageV1.abi.Unpack("getEnygmaAllParticipantsChainIds", data)
	if err != nil {
		return *new([]*big.Int), err
	}
	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	return out0, err
}

// PackGetEnygmaManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd6b81d5e.
//
// Solidity: function getEnygmaManager() view returns(address)
func (participantStorageV1 *ParticipantStorageV1) PackGetEnygmaManager() []byte {
	enc, err := participantStorageV1.abi.Pack("getEnygmaManager")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetEnygmaManager is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd6b81d5e.
//
// Solidity: function getEnygmaManager() view returns(address)
func (participantStorageV1 *ParticipantStorageV1) UnpackGetEnygmaManager(data []byte) (common.Address, error) {
	out, err := participantStorageV1.abi.Unpack("getEnygmaManager", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetKeyAgreements is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0d30d092.
//
// Solidity: function getKeyAgreements(uint256 chainId) view returns((uint256,bytes,bytes,uint256)[])
func (participantStorageV1 *ParticipantStorageV1) PackGetKeyAgreements(chainId *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("getKeyAgreements", chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetKeyAgreements is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0d30d092.
//
// Solidity: function getKeyAgreements(uint256 chainId) view returns((uint256,bytes,bytes,uint256)[])
func (participantStorageV1 *ParticipantStorageV1) UnpackGetKeyAgreements(data []byte) ([]ParticipantStructsKeyAgreementData, error) {
	out, err := participantStorageV1.abi.Unpack("getKeyAgreements", data)
	if err != nil {
		return *new([]ParticipantStructsKeyAgreementData), err
	}
	out0 := *abi.ConvertType(out[0], new([]ParticipantStructsKeyAgreementData)).(*[]ParticipantStructsKeyAgreementData)
	return out0, err
}

// PackGetParticipant is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1b9db2ef.
//
// Solidity: function getParticipant(uint256 chainId) view returns((uint256,uint8,uint8,string,string,uint256,uint256,bool))
func (participantStorageV1 *ParticipantStorageV1) PackGetParticipant(chainId *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("getParticipant", chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetParticipant is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1b9db2ef.
//
// Solidity: function getParticipant(uint256 chainId) view returns((uint256,uint8,uint8,string,string,uint256,uint256,bool))
func (participantStorageV1 *ParticipantStorageV1) UnpackGetParticipant(data []byte) (ParticipantStructsParticipant, error) {
	out, err := participantStorageV1.abi.Unpack("getParticipant", data)
	if err != nil {
		return *new(ParticipantStructsParticipant), err
	}
	out0 := *abi.ConvertType(out[0], new(ParticipantStructsParticipant)).(*ParticipantStructsParticipant)
	return out0, err
}

// PackGetParticipantCore is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x594c5c76.
//
// Solidity: function getParticipantCore() view returns(address)
func (participantStorageV1 *ParticipantStorageV1) PackGetParticipantCore() []byte {
	enc, err := participantStorageV1.abi.Pack("getParticipantCore")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetParticipantCore is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x594c5c76.
//
// Solidity: function getParticipantCore() view returns(address)
func (participantStorageV1 *ParticipantStorageV1) UnpackGetParticipantCore(data []byte) (common.Address, error) {
	out, err := participantStorageV1.abi.Unpack("getParticipantCore", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetParticipantDataBatch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa0842e33.
//
// Solidity: function getParticipantDataBatch() view returns((uint256,string,uint256)[] pnViewData, (uint256,string,bytes,bytes,uint256)[] auditInfo, uint256[] pnChainIds, uint256[] auditChainIds)
func (participantStorageV1 *ParticipantStorageV1) PackGetParticipantDataBatch() []byte {
	enc, err := participantStorageV1.abi.Pack("getParticipantDataBatch")
	if err != nil {
		panic(err)
	}
	return enc
}

// GetParticipantDataBatchOutput serves as a container for the return parameters of contract
// method GetParticipantDataBatch.
type GetParticipantDataBatchOutput struct {
	PnViewData    []ParticipantStructsPrivacyNodeViewData
	AuditInfo     []ParticipantStructsAuditInfoData
	PnChainIds    []*big.Int
	AuditChainIds []*big.Int
}

// UnpackGetParticipantDataBatch is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa0842e33.
//
// Solidity: function getParticipantDataBatch() view returns((uint256,string,uint256)[] pnViewData, (uint256,string,bytes,bytes,uint256)[] auditInfo, uint256[] pnChainIds, uint256[] auditChainIds)
func (participantStorageV1 *ParticipantStorageV1) UnpackGetParticipantDataBatch(data []byte) (GetParticipantDataBatchOutput, error) {
	out, err := participantStorageV1.abi.Unpack("getParticipantDataBatch", data)
	outstruct := new(GetParticipantDataBatchOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.PnViewData = *abi.ConvertType(out[0], new([]ParticipantStructsPrivacyNodeViewData)).(*[]ParticipantStructsPrivacyNodeViewData)
	outstruct.AuditInfo = *abi.ConvertType(out[1], new([]ParticipantStructsAuditInfoData)).(*[]ParticipantStructsAuditInfoData)
	outstruct.PnChainIds = *abi.ConvertType(out[2], new([]*big.Int)).(*[]*big.Int)
	outstruct.AuditChainIds = *abi.ConvertType(out[3], new([]*big.Int)).(*[]*big.Int)
	return *outstruct, err

}

// PackGetPaymentSpendPublicKeyByChainId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x81a40de6.
//
// Solidity: function getPaymentSpendPublicKeyByChainId(uint256 chainId) view returns(uint256)
func (participantStorageV1 *ParticipantStorageV1) PackGetPaymentSpendPublicKeyByChainId(chainId *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("getPaymentSpendPublicKeyByChainId", chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetPaymentSpendPublicKeyByChainId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x81a40de6.
//
// Solidity: function getPaymentSpendPublicKeyByChainId(uint256 chainId) view returns(uint256)
func (participantStorageV1 *ParticipantStorageV1) UnpackGetPaymentSpendPublicKeyByChainId(data []byte) (*big.Int, error) {
	out, err := participantStorageV1.abi.Unpack("getPaymentSpendPublicKeyByChainId", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x485cc955.
//
// Solidity: function initialize(address _endpoint, address authority_) returns()
func (participantStorageV1 *ParticipantStorageV1) PackInitialize(endpoint common.Address, authority common.Address) []byte {
	enc, err := participantStorageV1.abi.Pack("initialize", endpoint, authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackInitialize0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc4d66de8.
//
// Solidity: function initialize(address _endpoint) returns()
func (participantStorageV1 *ParticipantStorageV1) PackInitialize0(endpoint common.Address) []byte {
	enc, err := participantStorageV1.abi.Pack("initialize0", endpoint)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackInitiateKeyAgreement is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4ec53b2e.
//
// Solidity: function initiateKeyAgreement(uint256 initiatorChainId, uint256 responderChainId, bytes ciphertext, bytes digest, uint256 blockNumber) returns()
func (participantStorageV1 *ParticipantStorageV1) PackInitiateKeyAgreement(initiatorChainId *big.Int, responderChainId *big.Int, ciphertext []byte, digest []byte, blockNumber *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("initiateKeyAgreement", initiatorChainId, responderChainId, ciphertext, digest, blockNumber)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackProxiableUUID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (participantStorageV1 *ParticipantStorageV1) PackProxiableUUID() []byte {
	enc, err := participantStorageV1.abi.Pack("proxiableUUID")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackProxiableUUID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (participantStorageV1 *ParticipantStorageV1) UnpackProxiableUUID(data []byte) ([32]byte, error) {
	out, err := participantStorageV1.abi.Unpack("proxiableUUID", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackRemoveParticipant is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x683f7f27.
//
// Solidity: function removeParticipant(uint256 chainId) returns()
func (participantStorageV1 *ParticipantStorageV1) PackRemoveParticipant(chainId *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("removeParticipant", chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (participantStorageV1 *ParticipantStorageV1) PackResourceId() []byte {
	enc, err := participantStorageV1.abi.Pack("resourceId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (participantStorageV1 *ParticipantStorageV1) UnpackResourceId(data []byte) ([32]byte, error) {
	out, err := participantStorageV1.abi.Unpack("resourceId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackSetAuditInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6279aa9f.
//
// Solidity: function setAuditInfo(uint256 chainId, string raylsViewPublicKey, bytes encryptedRaylsViewPrivateKey, bytes mac, uint256 blockNumber) returns()
func (participantStorageV1 *ParticipantStorageV1) PackSetAuditInfo(chainId *big.Int, raylsViewPublicKey string, encryptedRaylsViewPrivateKey []byte, mac []byte, blockNumber *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("setAuditInfo", chainId, raylsViewPublicKey, encryptedRaylsViewPrivateKey, mac, blockNumber)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetAuditManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc68bab0f.
//
// Solidity: function setAuditManager(address _auditManager) returns()
func (participantStorageV1 *ParticipantStorageV1) PackSetAuditManager(auditManager common.Address) []byte {
	enc, err := participantStorageV1.abi.Pack("setAuditManager", auditManager)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetChainViewData is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4017734d.
//
// Solidity: function setChainViewData(uint256 chainId, string raylsViewPublicKey, uint256 blockNumber) returns()
func (participantStorageV1 *ParticipantStorageV1) PackSetChainViewData(chainId *big.Int, raylsViewPublicKey string, blockNumber *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("setChainViewData", chainId, raylsViewPublicKey, blockNumber)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetEnygmaManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xede591fd.
//
// Solidity: function setEnygmaManager(address _enygmaManager) returns()
func (participantStorageV1 *ParticipantStorageV1) PackSetEnygmaManager(enygmaManager common.Address) []byte {
	enc, err := participantStorageV1.abi.Pack("setEnygmaManager", enygmaManager)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetEnygmaPnEventsAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52bf1c8d.
//
// Solidity: function setEnygmaPnEventsAddress(address _pnEnygmaEvents) returns()
func (participantStorageV1 *ParticipantStorageV1) PackSetEnygmaPnEventsAddress(pnEnygmaEvents common.Address) []byte {
	enc, err := participantStorageV1.abi.Pack("setEnygmaPnEventsAddress", pnEnygmaEvents)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetParticipantCore is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x035c8713.
//
// Solidity: function setParticipantCore(address _participantCore) returns()
func (participantStorageV1 *ParticipantStorageV1) PackSetParticipantCore(participantCore common.Address) []byte {
	enc, err := participantStorageV1.abi.Pack("setParticipantCore", participantCore)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetPaymentSpendPublicKey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x079c2a49.
//
// Solidity: function setPaymentSpendPublicKey(uint256 _chainId, uint256 _paymentSpendPublicKey, address[] _pnAddresses) returns()
func (participantStorageV1 *ParticipantStorageV1) PackSetPaymentSpendPublicKey(chainId *big.Int, paymentSpendPublicKey *big.Int, pnAddresses []common.Address) []byte {
	enc, err := participantStorageV1.abi.Pack("setPaymentSpendPublicKey", chainId, paymentSpendPublicKey, pnAddresses)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpdateBroadcastMessagesPermission is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x33ded4a8.
//
// Solidity: function updateBroadcastMessagesPermission(uint256 chainId, bool allowed) returns()
func (participantStorageV1 *ParticipantStorageV1) PackUpdateBroadcastMessagesPermission(chainId *big.Int, allowed bool) []byte {
	enc, err := participantStorageV1.abi.Pack("updateBroadcastMessagesPermission", chainId, allowed)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpdateRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x65ffa17e.
//
// Solidity: function updateRole(uint256 chainId, uint8 role) returns()
func (participantStorageV1 *ParticipantStorageV1) PackUpdateRole(chainId *big.Int, role uint8) []byte {
	enc, err := participantStorageV1.abi.Pack("updateRole", chainId, role)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpdateStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3a1b3d31.
//
// Solidity: function updateStatus(uint256 chainId, uint8 status) returns()
func (participantStorageV1 *ParticipantStorageV1) PackUpdateStatus(chainId *big.Int, status uint8) []byte {
	enc, err := participantStorageV1.abi.Pack("updateStatus", chainId, status)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpgradeToAndCall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (participantStorageV1 *ParticipantStorageV1) PackUpgradeToAndCall(newImplementation common.Address, data []byte) []byte {
	enc, err := participantStorageV1.abi.Pack("upgradeToAndCall", newImplementation, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackValidateMessageParticipants is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd5c3614f.
//
// Solidity: function validateMessageParticipants(uint256 originChainId, uint256 destinationChainId) view returns()
func (participantStorageV1 *ParticipantStorageV1) PackValidateMessageParticipants(originChainId *big.Int, destinationChainId *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("validateMessageParticipants", originChainId, destinationChainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackValidateParticipantStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc9557f56.
//
// Solidity: function validateParticipantStatus(uint256 chainId) view returns()
func (participantStorageV1 *ParticipantStorageV1) PackValidateParticipantStatus(chainId *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("validateParticipantStatus", chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackVerifyParticipant is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb3816fdb.
//
// Solidity: function verifyParticipant(uint256 chainId) view returns(bool)
func (participantStorageV1 *ParticipantStorageV1) PackVerifyParticipant(chainId *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("verifyParticipant", chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackVerifyParticipant is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb3816fdb.
//
// Solidity: function verifyParticipant(uint256 chainId) view returns(bool)
func (participantStorageV1 *ParticipantStorageV1) UnpackVerifyParticipant(data []byte) (bool, error) {
	out, err := participantStorageV1.abi.Unpack("verifyParticipant", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// ParticipantStorageV1AuthorityUpdated represents a AuthorityUpdated event raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1AuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const ParticipantStorageV1AuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (ParticipantStorageV1AuthorityUpdated) ContractEventName() string {
	return ParticipantStorageV1AuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (participantStorageV1 *ParticipantStorageV1) UnpackAuthorityUpdatedEvent(log *types.Log) (*ParticipantStorageV1AuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != participantStorageV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ParticipantStorageV1AuthorityUpdated)
	if len(log.Data) > 0 {
		if err := participantStorageV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range participantStorageV1.abi.Events[event].Inputs {
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

// ParticipantStorageV1Initialized represents a Initialized event raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1Initialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const ParticipantStorageV1InitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (ParticipantStorageV1Initialized) ContractEventName() string {
	return ParticipantStorageV1InitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (participantStorageV1 *ParticipantStorageV1) UnpackInitializedEvent(log *types.Log) (*ParticipantStorageV1Initialized, error) {
	event := "Initialized"
	if log.Topics[0] != participantStorageV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ParticipantStorageV1Initialized)
	if len(log.Data) > 0 {
		if err := participantStorageV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range participantStorageV1.abi.Events[event].Inputs {
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

// ParticipantStorageV1ModulesConfigured represents a ModulesConfigured event raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1ModulesConfigured struct {
	ParticipantCore common.Address
	AuditManager    common.Address
	EnygmaManager   common.Address
	Raw             *types.Log // Blockchain specific contextual infos
}

const ParticipantStorageV1ModulesConfiguredEventName = "ModulesConfigured"

// ContractEventName returns the user-defined event name.
func (ParticipantStorageV1ModulesConfigured) ContractEventName() string {
	return ParticipantStorageV1ModulesConfiguredEventName
}

// UnpackModulesConfiguredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ModulesConfigured(address indexed participantCore, address indexed auditManager, address indexed enygmaManager)
func (participantStorageV1 *ParticipantStorageV1) UnpackModulesConfiguredEvent(log *types.Log) (*ParticipantStorageV1ModulesConfigured, error) {
	event := "ModulesConfigured"
	if log.Topics[0] != participantStorageV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ParticipantStorageV1ModulesConfigured)
	if len(log.Data) > 0 {
		if err := participantStorageV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range participantStorageV1.abi.Events[event].Inputs {
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

// ParticipantStorageV1Upgraded represents a Upgraded event raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1Upgraded struct {
	Implementation common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const ParticipantStorageV1UpgradedEventName = "Upgraded"

// ContractEventName returns the user-defined event name.
func (ParticipantStorageV1Upgraded) ContractEventName() string {
	return ParticipantStorageV1UpgradedEventName
}

// UnpackUpgradedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Upgraded(address indexed implementation)
func (participantStorageV1 *ParticipantStorageV1) UnpackUpgradedEvent(log *types.Log) (*ParticipantStorageV1Upgraded, error) {
	event := "Upgraded"
	if log.Topics[0] != participantStorageV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ParticipantStorageV1Upgraded)
	if len(log.Data) > 0 {
		if err := participantStorageV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range participantStorageV1.abi.Events[event].Inputs {
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
func (participantStorageV1 *ParticipantStorageV1) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["AddressEmptyCode"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackAddressEmptyCodeError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["FailedCall"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackFailedCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["ParticipantStorageV1UnauthorizedCaller"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackParticipantStorageV1UnauthorizedCallerError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["UUPSUnauthorizedCallContext"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackUUPSUnauthorizedCallContextError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["UUPSUnsupportedProxiableUUID"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackUUPSUnsupportedProxiableUUIDError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ParticipantStorageV1AddressEmptyCode represents a AddressEmptyCode error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1AddressEmptyCode struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressEmptyCode(address target)
func ParticipantStorageV1AddressEmptyCodeErrorID() common.Hash {
	return common.HexToHash("0x9996b315c842ff135b8fc4a08ad5df1c344efbc03d2687aecc0678050d2aac89")
}

// UnpackAddressEmptyCodeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressEmptyCode(address target)
func (participantStorageV1 *ParticipantStorageV1) UnpackAddressEmptyCodeError(raw []byte) (*ParticipantStorageV1AddressEmptyCode, error) {
	out := new(ParticipantStorageV1AddressEmptyCode)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "AddressEmptyCode", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1ERC1967InvalidImplementation represents a ERC1967InvalidImplementation error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1ERC1967InvalidImplementation struct {
	Implementation common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func ParticipantStorageV1ERC1967InvalidImplementationErrorID() common.Hash {
	return common.HexToHash("0x4c9c8ce3ceb3130f17f7cdba48d89b5b0129f266a8bac114e6e315a41879b617")
}

// UnpackERC1967InvalidImplementationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func (participantStorageV1 *ParticipantStorageV1) UnpackERC1967InvalidImplementationError(raw []byte) (*ParticipantStorageV1ERC1967InvalidImplementation, error) {
	out := new(ParticipantStorageV1ERC1967InvalidImplementation)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "ERC1967InvalidImplementation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1ERC1967NonPayable represents a ERC1967NonPayable error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1ERC1967NonPayable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967NonPayable()
func ParticipantStorageV1ERC1967NonPayableErrorID() common.Hash {
	return common.HexToHash("0xb398979fa84f543c8e222f17890372c487baf85e062276c127fef521eea7224b")
}

// UnpackERC1967NonPayableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967NonPayable()
func (participantStorageV1 *ParticipantStorageV1) UnpackERC1967NonPayableError(raw []byte) (*ParticipantStorageV1ERC1967NonPayable, error) {
	out := new(ParticipantStorageV1ERC1967NonPayable)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "ERC1967NonPayable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1FailedCall represents a FailedCall error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1FailedCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedCall()
func ParticipantStorageV1FailedCallErrorID() common.Hash {
	return common.HexToHash("0xd6bda27508c0fb6d8a39b4b122878dab26f731a7d4e4abe711dd3731899052a4")
}

// UnpackFailedCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedCall()
func (participantStorageV1 *ParticipantStorageV1) UnpackFailedCallError(raw []byte) (*ParticipantStorageV1FailedCall, error) {
	out := new(ParticipantStorageV1FailedCall)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "FailedCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1InvalidInitialization represents a InvalidInitialization error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1InvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func ParticipantStorageV1InvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (participantStorageV1 *ParticipantStorageV1) UnpackInvalidInitializationError(raw []byte) (*ParticipantStorageV1InvalidInitialization, error) {
	out := new(ParticipantStorageV1InvalidInitialization)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1NotInitializing represents a NotInitializing error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1NotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func ParticipantStorageV1NotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (participantStorageV1 *ParticipantStorageV1) UnpackNotInitializingError(raw []byte) (*ParticipantStorageV1NotInitializing, error) {
	out := new(ParticipantStorageV1NotInitializing)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1ParticipantStorageV1UnauthorizedCaller represents a ParticipantStorageV1__UnauthorizedCaller error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1ParticipantStorageV1UnauthorizedCaller struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ParticipantStorageV1__UnauthorizedCaller(address caller)
func ParticipantStorageV1ParticipantStorageV1UnauthorizedCallerErrorID() common.Hash {
	return common.HexToHash("0x462e5c5b72c0fc8d71c0efe37e91aba3df63f77db44f129a892a044de14a6df8")
}

// UnpackParticipantStorageV1UnauthorizedCallerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ParticipantStorageV1__UnauthorizedCaller(address caller)
func (participantStorageV1 *ParticipantStorageV1) UnpackParticipantStorageV1UnauthorizedCallerError(raw []byte) (*ParticipantStorageV1ParticipantStorageV1UnauthorizedCaller, error) {
	out := new(ParticipantStorageV1ParticipantStorageV1UnauthorizedCaller)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "ParticipantStorageV1UnauthorizedCaller", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1RaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1RaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func ParticipantStorageV1RaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (participantStorageV1 *ParticipantStorageV1) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*ParticipantStorageV1RaylsAccessManagedContractPaused, error) {
	out := new(ParticipantStorageV1RaylsAccessManagedContractPaused)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1RaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1RaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func ParticipantStorageV1RaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (participantStorageV1 *ParticipantStorageV1) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*ParticipantStorageV1RaylsAccessManagedInvalidAuthority, error) {
	out := new(ParticipantStorageV1RaylsAccessManagedInvalidAuthority)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1RaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1RaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func ParticipantStorageV1RaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (participantStorageV1 *ParticipantStorageV1) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*ParticipantStorageV1RaylsAccessManagedMustSchedule, error) {
	out := new(ParticipantStorageV1RaylsAccessManagedMustSchedule)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1RaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1RaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func ParticipantStorageV1RaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (participantStorageV1 *ParticipantStorageV1) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*ParticipantStorageV1RaylsAccessManagedUnauthorized, error) {
	out := new(ParticipantStorageV1RaylsAccessManagedUnauthorized)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1UUPSUnauthorizedCallContext represents a UUPSUnauthorizedCallContext error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1UUPSUnauthorizedCallContext struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnauthorizedCallContext()
func ParticipantStorageV1UUPSUnauthorizedCallContextErrorID() common.Hash {
	return common.HexToHash("0xe07c8dba242a06571ac65fe4bbe20522c9fb111cb33599b799ff8039c1ed18f4")
}

// UnpackUUPSUnauthorizedCallContextError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnauthorizedCallContext()
func (participantStorageV1 *ParticipantStorageV1) UnpackUUPSUnauthorizedCallContextError(raw []byte) (*ParticipantStorageV1UUPSUnauthorizedCallContext, error) {
	out := new(ParticipantStorageV1UUPSUnauthorizedCallContext)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "UUPSUnauthorizedCallContext", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1UUPSUnsupportedProxiableUUID represents a UUPSUnsupportedProxiableUUID error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1UUPSUnsupportedProxiableUUID struct {
	Slot [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func ParticipantStorageV1UUPSUnsupportedProxiableUUIDErrorID() common.Hash {
	return common.HexToHash("0xaa1d49a4c084bfa9aeeee2a0be65267a7f19ba7e1476b114dac513d2c14cb563")
}

// UnpackUUPSUnsupportedProxiableUUIDError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func (participantStorageV1 *ParticipantStorageV1) UnpackUUPSUnsupportedProxiableUUIDError(raw []byte) (*ParticipantStorageV1UUPSUnsupportedProxiableUUID, error) {
	out := new(ParticipantStorageV1UUPSUnsupportedProxiableUUID)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "UUPSUnsupportedProxiableUUID", raw); err != nil {
		return nil, err
	}
	return out, nil
}
