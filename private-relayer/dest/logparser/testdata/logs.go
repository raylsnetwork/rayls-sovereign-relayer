package testdata

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/DvpTeleport"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/EndpointV1"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/EnygmaTeleport"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/TeleportV1"
)

// ============================================================================
// EndpointV1MessageDispatched
// ============================================================================

// EndpointV1MessageDispatchedLogOption represents a modification to an EndpointV1MessageDispatched log fixture.
type EndpointV1MessageDispatchedLogOption func(*types.Log, *endpointV1MessageDispatchedData)

type endpointV1MessageDispatchedData struct {
	messageId [32]byte
	from      common.Address
	toChainId *big.Int
	to        common.Address
	data      EndpointV1.RaylsMessage
}

// WithMessageDispatchedMessageId sets the MessageId field on an EndpointV1MessageDispatched log.
func WithMessageDispatchedMessageId(messageId [32]byte) EndpointV1MessageDispatchedLogOption {
	return func(log *types.Log, data *endpointV1MessageDispatchedData) {
		data.messageId = messageId
	}
}

// WithMessageDispatchedMessageID sets the MessageId field on an EndpointV1MessageDispatched log using common.Hash.
func WithMessageDispatchedMessageID(messageID common.Hash) EndpointV1MessageDispatchedLogOption {
	return func(log *types.Log, data *endpointV1MessageDispatchedData) {
		data.messageId = messageID
	}
}

// WithMessageDispatchedFrom sets the From field on an EndpointV1MessageDispatched log.
func WithMessageDispatchedFrom(from common.Address) EndpointV1MessageDispatchedLogOption {
	return func(log *types.Log, data *endpointV1MessageDispatchedData) {
		data.from = from
	}
}

// WithMessageDispatchedToChainId sets the ToChainId field on an EndpointV1MessageDispatched log.
func WithMessageDispatchedToChainId(toChainId *big.Int) EndpointV1MessageDispatchedLogOption {
	return func(log *types.Log, data *endpointV1MessageDispatchedData) {
		data.toChainId = toChainId
	}
}

// WithMessageDispatchedTo sets the To field on an EndpointV1MessageDispatched log.
func WithMessageDispatchedTo(to common.Address) EndpointV1MessageDispatchedLogOption {
	return func(log *types.Log, data *endpointV1MessageDispatchedData) {
		data.to = to
	}
}

// WithMessageDispatchedData sets the Data field on an EndpointV1MessageDispatched log.
func WithMessageDispatchedData(raylsData EndpointV1.RaylsMessage) EndpointV1MessageDispatchedLogOption {
	return func(log *types.Log, data *endpointV1MessageDispatchedData) {
		data.data = raylsData
	}
}

// WithMessageDispatchedBlockNumber sets the BlockNumber field on an EndpointV1MessageDispatched log.
func WithMessageDispatchedBlockNumber(blockNumber uint64) EndpointV1MessageDispatchedLogOption {
	return func(log *types.Log, data *endpointV1MessageDispatchedData) {
		log.BlockNumber = blockNumber
	}
}

// WithMessageDispatchedTxHash sets the TxHash field on an EndpointV1MessageDispatched log.
func WithMessageDispatchedTxHash(txHash common.Hash) EndpointV1MessageDispatchedLogOption {
	return func(log *types.Log, data *endpointV1MessageDispatchedData) {
		log.TxHash = txHash
	}
}

// WithMessageDispatchedAddress sets the contract Address field on an EndpointV1MessageDispatched log.
func WithMessageDispatchedAddress(address common.Address) EndpointV1MessageDispatchedLogOption {
	return func(log *types.Log, data *endpointV1MessageDispatchedData) {
		log.Address = address
	}
}

// GetEndpointV1MessageDispatchedLog returns a sample EndpointV1MessageDispatched event log with default values.
func GetEndpointV1MessageDispatchedLog() types.Log {
	return NewEndpointV1MessageDispatchedLogWith()
}

// NewEndpointV1MessageDispatchedLogWith creates an EndpointV1MessageDispatched log using default values
// and then applies any provided options.
func NewEndpointV1MessageDispatchedLogWith(opts ...EndpointV1MessageDispatchedLogOption) types.Log {
	// Default values
	data := endpointV1MessageDispatchedData{
		messageId: [32]byte{0xaa, 0xbb, 0xcc},
		from:      common.HexToAddress("0x1111111111111111111111111111111111111111"),
		toChainId: big.NewInt(1888),
		to:        common.HexToAddress("0x2222222222222222222222222222222222222222"),
		data: EndpointV1.RaylsMessage{
			MessageMetadata: EndpointV1.RaylsMessageMetadata{
				Valid:        true,
				Nonce:        big.NewInt(1),
				ResourceId:   [32]byte{},
				LockData:     []byte{},
				IgnoresNonce: false,
				NewResourceMetadata: EndpointV1.NewResourceMetadata{
					Valid:              false,
					ResourceDeployType: 0,
					Bytecode:           []byte{},
					FactoryTemplate:    0,
					InitializerParams:  []byte{},
				},
				RevertPayloadDataSender:   []byte{},
				RevertPayloadDataReceiver: []byte{},
				TransferMetadata: EndpointV1.BridgedTransferMetadata{
					AssetType:    0,
					Id:           big.NewInt(0),
					From:         common.Address{},
					To:           common.Address{},
					TokenAddress: common.Address{},
					Amount:       big.NewInt(0),
				},
			},
			Payload: []byte{},
		},
	}

	log := types.Log{
		Address:     common.HexToAddress("0x3333333333333333333333333333333333333333"),
		BlockNumber: 100,
		TxHash:      common.HexToHash("0xabcdef0123456789"),
		TxIndex:     1,
		BlockHash:   common.HexToHash("0x123456789abcdef0"),
		Index:       0,
	}

	// Apply options
	for _, opt := range opts {
		opt(&log, &data)
	}

	// Pack the data using ABI
	abi, _ := EndpointV1.EndpointV1MetaData.ParseABI()
	packedData, _ := abi.Events["MessageDispatched"].Inputs.NonIndexed().Pack(data.to, data.data)

	// Set topics: event signature + indexed parameters (messageId, from, toChainId)
	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("MessageDispatched(bytes32,address,uint256,address,((bool,uint256,(bool,uint8,bytes,uint8,bytes),bytes32,bytes,bytes,bytes,(uint8,uint256,address,address,address,uint256),bool),bytes))")),
		data.messageId,
		common.BytesToHash(data.from.Bytes()),
		common.BigToHash(data.toChainId),
	}
	log.Data = packedData

	return log
}

// ============================================================================
// TeleportV1EncryptedDataBatchStored
// ============================================================================

// TeleportV1EncryptedDataBatchStoredLogOption represents a modification to a TeleportV1EncryptedDataBatchStored log fixture.
type TeleportV1EncryptedDataBatchStoredLogOption func(*types.Log, *teleportV1EncryptedDataBatchStoredData)

type teleportV1EncryptedDataBatchStoredData struct {
	print       string
	data        []byte
	blockNumber *big.Int
}

// WithEncryptedDataBatchStoredPrint sets the Print field on a TeleportV1EncryptedDataBatchStored log.
func WithEncryptedDataBatchStoredPrint(print string) TeleportV1EncryptedDataBatchStoredLogOption {
	return func(log *types.Log, data *teleportV1EncryptedDataBatchStoredData) {
		data.print = print
	}
}

// WithEncryptedDataBatchStoredData sets the Data field on a TeleportV1EncryptedDataBatchStored log.
func WithEncryptedDataBatchStoredData(encryptedData []byte) TeleportV1EncryptedDataBatchStoredLogOption {
	return func(log *types.Log, data *teleportV1EncryptedDataBatchStoredData) {
		data.data = encryptedData
	}
}

// WithEncryptedDataBatchStoredBlockNumber sets the BlockNumber field on a TeleportV1EncryptedDataBatchStored log.
func WithEncryptedDataBatchStoredBlockNumber(blockNumber uint64) TeleportV1EncryptedDataBatchStoredLogOption {
	return func(log *types.Log, data *teleportV1EncryptedDataBatchStoredData) {
		data.blockNumber = new(big.Int).SetUint64(blockNumber)
		log.BlockNumber = blockNumber
	}
}

// WithEncryptedDataBatchStoredTxHash sets the TxHash field on a TeleportV1EncryptedDataBatchStored log.
func WithEncryptedDataBatchStoredTxHash(txHash common.Hash) TeleportV1EncryptedDataBatchStoredLogOption {
	return func(log *types.Log, data *teleportV1EncryptedDataBatchStoredData) {
		log.TxHash = txHash
	}
}

// WithEncryptedDataBatchStoredAddress sets the contract Address field on a TeleportV1EncryptedDataBatchStored log.
func WithEncryptedDataBatchStoredAddress(address common.Address) TeleportV1EncryptedDataBatchStoredLogOption {
	return func(log *types.Log, data *teleportV1EncryptedDataBatchStoredData) {
		log.Address = address
	}
}

// GetTeleportV1EncryptedDataBatchStoredLog returns a sample TeleportV1EncryptedDataBatchStored event log with default values.
func GetTeleportV1EncryptedDataBatchStoredLog() types.Log {
	return NewTeleportV1EncryptedDataBatchStoredLogWith()
}

// NewTeleportV1EncryptedDataBatchStoredLogWith creates a TeleportV1EncryptedDataBatchStored log using default values
// and then applies any provided options.
func NewTeleportV1EncryptedDataBatchStoredLogWith(opts ...TeleportV1EncryptedDataBatchStoredLogOption) types.Log {
	// Default values
	data := teleportV1EncryptedDataBatchStoredData{
		print:       "batch-fingerprint-001",
		data:        []byte{0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe},
		blockNumber: big.NewInt(100),
	}

	log := types.Log{
		Address:     common.HexToAddress("0x4444444444444444444444444444444444444444"),
		BlockNumber: 100,
		TxHash:      common.HexToHash("0xfedcba9876543210"),
		TxIndex:     2,
		BlockHash:   common.HexToHash("0x0fedcba987654321"),
		Index:       1,
	}

	// Apply options
	for _, opt := range opts {
		opt(&log, &data)
	}

	// Pack the data using ABI
	abi, _ := TeleportV1.TeleportV1MetaData.ParseABI()
	packedData, _ := abi.Events["EncryptedDataBatchStored"].Inputs.NonIndexed().Pack(data.print, data.data)

	// Set topics: event signature + indexed parameters (blockNumber)
	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("EncryptedDataBatchStored(string,bytes,uint256)")),
		common.BigToHash(data.blockNumber),
	}
	log.Data = packedData

	return log
}

// ============================================================================
// TeleportV1AtomicMessageStatusChangedBatch
// ============================================================================

// TeleportV1AtomicMessageStatusChangedBatchLogOption represents a modification to a TeleportV1AtomicMessageStatusChangedBatch log fixture.
type TeleportV1AtomicMessageStatusChangedBatchLogOption func(*types.Log, *teleportV1AtomicMessageStatusChangedBatchData)

type teleportV1AtomicMessageStatusChangedBatchData struct {
	msgIds []string
	status uint8
}

// WithAtomicStatusChangedMsgIds sets the MsgIds field on a TeleportV1AtomicMessageStatusChangedBatch log.
func WithAtomicStatusChangedMsgIds(msgIds []string) TeleportV1AtomicMessageStatusChangedBatchLogOption {
	return func(log *types.Log, data *teleportV1AtomicMessageStatusChangedBatchData) {
		data.msgIds = msgIds
	}
}

// WithAtomicStatusChangedStatus sets the Status field on a TeleportV1AtomicMessageStatusChangedBatch log.
func WithAtomicStatusChangedStatus(status uint8) TeleportV1AtomicMessageStatusChangedBatchLogOption {
	return func(log *types.Log, data *teleportV1AtomicMessageStatusChangedBatchData) {
		data.status = status
	}
}

// WithAtomicStatusChangedBlockNumber sets the BlockNumber field on a TeleportV1AtomicMessageStatusChangedBatch log.
func WithAtomicStatusChangedBlockNumber(blockNumber uint64) TeleportV1AtomicMessageStatusChangedBatchLogOption {
	return func(log *types.Log, data *teleportV1AtomicMessageStatusChangedBatchData) {
		log.BlockNumber = blockNumber
	}
}

// WithAtomicStatusChangedTxHash sets the TxHash field on a TeleportV1AtomicMessageStatusChangedBatch log.
func WithAtomicStatusChangedTxHash(txHash common.Hash) TeleportV1AtomicMessageStatusChangedBatchLogOption {
	return func(log *types.Log, data *teleportV1AtomicMessageStatusChangedBatchData) {
		log.TxHash = txHash
	}
}

// WithAtomicStatusChangedAddress sets the contract Address field on a TeleportV1AtomicMessageStatusChangedBatch log.
func WithAtomicStatusChangedAddress(address common.Address) TeleportV1AtomicMessageStatusChangedBatchLogOption {
	return func(log *types.Log, data *teleportV1AtomicMessageStatusChangedBatchData) {
		log.Address = address
	}
}

// NewTeleportV1AtomicMessageStatusChangedBatchLogWith creates a TeleportV1AtomicMessageStatusChangedBatch log
// using default values and then applies any provided options.
func NewTeleportV1AtomicMessageStatusChangedBatchLogWith(opts ...TeleportV1AtomicMessageStatusChangedBatchLogOption) types.Log {
	// Default values
	data := teleportV1AtomicMessageStatusChangedBatchData{
		msgIds: []string{"shared-id-001", "shared-id-002"},
		status: 1, // AtomicExecutedStatus
	}

	log := types.Log{
		Address:     common.HexToAddress("0x4444444444444444444444444444444444444444"),
		BlockNumber: 100,
		TxHash:      common.HexToHash("0xabcdef1234567890"),
		TxIndex:     5,
		BlockHash:   common.HexToHash("0x0987654321fedcba"),
		Index:       4,
	}

	// Apply options
	for _, opt := range opts {
		opt(&log, &data)
	}

	// Pack the data using ABI
	abi, _ := TeleportV1.TeleportV1MetaData.ParseABI()
	packedData, _ := abi.Events["AtomicMessageStatusChangedBatch"].Inputs.NonIndexed().Pack(data.msgIds, data.status)

	// Set topics: event signature only (no indexed parameters)
	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("AtomicMessageStatusChangedBatch(string[],uint8)")),
	}
	log.Data = packedData

	return log
}

// ============================================================================
// DvpTeleportTransferEncryptedData
// ============================================================================

// DvpTeleportTransferEncryptedDataLogOption represents a modification to a DvpTeleportTransferEncryptedData log fixture.
type DvpTeleportTransferEncryptedDataLogOption func(*types.Log, *dvpTeleportTransferEncryptedDataData)

type dvpTeleportTransferEncryptedDataData struct {
	print       string
	data        []byte
	blockNumber *big.Int
}

// WithDvpTransferEncryptedDataPrint sets the Print field on a DvpTeleportTransferEncryptedData log.
func WithDvpTransferEncryptedDataPrint(print string) DvpTeleportTransferEncryptedDataLogOption {
	return func(log *types.Log, data *dvpTeleportTransferEncryptedDataData) {
		data.print = print
	}
}

// WithDvpTransferEncryptedDataData sets the Data field on a DvpTeleportTransferEncryptedData log.
func WithDvpTransferEncryptedDataData(encryptedData []byte) DvpTeleportTransferEncryptedDataLogOption {
	return func(log *types.Log, data *dvpTeleportTransferEncryptedDataData) {
		data.data = encryptedData
	}
}

// WithDvpTransferEncryptedDataBlockNumber sets the BlockNumber field on a DvpTeleportTransferEncryptedData log.
func WithDvpTransferEncryptedDataBlockNumber(blockNumber uint64) DvpTeleportTransferEncryptedDataLogOption {
	return func(log *types.Log, data *dvpTeleportTransferEncryptedDataData) {
		data.blockNumber = new(big.Int).SetUint64(blockNumber)
		log.BlockNumber = blockNumber
	}
}

// WithDvpTransferEncryptedDataTxHash sets the TxHash field on a DvpTeleportTransferEncryptedData log.
func WithDvpTransferEncryptedDataTxHash(txHash common.Hash) DvpTeleportTransferEncryptedDataLogOption {
	return func(log *types.Log, data *dvpTeleportTransferEncryptedDataData) {
		log.TxHash = txHash
	}
}

// WithDvpTransferEncryptedDataAddress sets the contract Address field on a DvpTeleportTransferEncryptedData log.
func WithDvpTransferEncryptedDataAddress(address common.Address) DvpTeleportTransferEncryptedDataLogOption {
	return func(log *types.Log, data *dvpTeleportTransferEncryptedDataData) {
		log.Address = address
	}
}

// GetDvpTeleportTransferEncryptedDataLog returns a sample DvpTeleportTransferEncryptedData event log with default values.
func GetDvpTeleportTransferEncryptedDataLog() types.Log {
	return NewDvpTeleportTransferEncryptedDataLogWith()
}

// NewDvpTeleportTransferEncryptedDataLogWith creates a DvpTeleportTransferEncryptedData log using default values
// and then applies any provided options.
func NewDvpTeleportTransferEncryptedDataLogWith(opts ...DvpTeleportTransferEncryptedDataLogOption) types.Log {
	// Default values
	data := dvpTeleportTransferEncryptedDataData{
		print:       "dvp-transfer-fingerprint-001",
		data:        []byte{0xba, 0xad, 0xf0, 0x0d, 0xc0, 0xde},
		blockNumber: big.NewInt(200),
	}

	log := types.Log{
		Address:     common.HexToAddress("0x5555555555555555555555555555555555555555"),
		BlockNumber: 200,
		TxHash:      common.HexToHash("0x1122334455667788"),
		TxIndex:     3,
		BlockHash:   common.HexToHash("0x8877665544332211"),
		Index:       2,
	}

	// Apply options
	for _, opt := range opts {
		opt(&log, &data)
	}

	// Pack the data using ABI
	abi, _ := DvpTeleport.DvpTeleportMetaData.ParseABI()
	packedData, _ := abi.Events["TransferEncryptedData"].Inputs.NonIndexed().Pack(data.print, data.data)

	// Set topics: event signature + indexed parameters (blockNumber)
	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("TransferEncryptedData(string,bytes,uint256)")),
		common.BigToHash(data.blockNumber),
	}
	log.Data = packedData

	return log
}

// ============================================================================
// EnygmaTeleportEnygmaTransfer
// ============================================================================

// EnygmaTeleportEnygmaTransferLogOption represents a modification to an EnygmaTeleportEnygmaTransfer log fixture.
type EnygmaTeleportEnygmaTransferLogOption func(*types.Log, *enygmaTeleportEnygmaTransferData)

type enygmaTeleportEnygmaTransferData struct {
	resourceId       [32]byte
	toChainId        *big.Int
	encryptedMessage []byte
}

// WithEnygmaTransferResourceId sets the ResourceId field on an EnygmaTeleportEnygmaTransfer log.
func WithEnygmaTransferResourceId(resourceId [32]byte) EnygmaTeleportEnygmaTransferLogOption {
	return func(log *types.Log, data *enygmaTeleportEnygmaTransferData) {
		data.resourceId = resourceId
	}
}

// WithEnygmaTransferToChainId sets the ToChainId field on an EnygmaTeleportEnygmaTransfer log.
func WithEnygmaTransferToChainId(toChainId *big.Int) EnygmaTeleportEnygmaTransferLogOption {
	return func(log *types.Log, data *enygmaTeleportEnygmaTransferData) {
		data.toChainId = toChainId
	}
}

// WithEnygmaTransferEncryptedMessage sets the EncryptedMessage field on an EnygmaTeleportEnygmaTransfer log.
func WithEnygmaTransferEncryptedMessage(encryptedMessage []byte) EnygmaTeleportEnygmaTransferLogOption {
	return func(log *types.Log, data *enygmaTeleportEnygmaTransferData) {
		data.encryptedMessage = encryptedMessage
	}
}

// WithEnygmaTransferBlockNumber sets the BlockNumber field on an EnygmaTeleportEnygmaTransfer log.
func WithEnygmaTransferBlockNumber(blockNumber uint64) EnygmaTeleportEnygmaTransferLogOption {
	return func(log *types.Log, data *enygmaTeleportEnygmaTransferData) {
		log.BlockNumber = blockNumber
	}
}

// WithEnygmaTransferTxHash sets the TxHash field on an EnygmaTeleportEnygmaTransfer log.
func WithEnygmaTransferTxHash(txHash common.Hash) EnygmaTeleportEnygmaTransferLogOption {
	return func(log *types.Log, data *enygmaTeleportEnygmaTransferData) {
		log.TxHash = txHash
	}
}

// WithEnygmaTransferAddress sets the contract Address field on an EnygmaTeleportEnygmaTransfer log.
func WithEnygmaTransferAddress(address common.Address) EnygmaTeleportEnygmaTransferLogOption {
	return func(log *types.Log, data *enygmaTeleportEnygmaTransferData) {
		log.Address = address
	}
}

// GetEnygmaTeleportEnygmaTransferLog returns a sample EnygmaTeleportEnygmaTransfer event log with default values.
func GetEnygmaTeleportEnygmaTransferLog() types.Log {
	return NewEnygmaTeleportEnygmaTransferLogWith()
}

// NewEnygmaTeleportEnygmaTransferLogWith creates an EnygmaTeleportEnygmaTransfer log using default values
// and then applies any provided options.
func NewEnygmaTeleportEnygmaTransferLogWith(opts ...EnygmaTeleportEnygmaTransferLogOption) types.Log {
	// Default values
	data := enygmaTeleportEnygmaTransferData{
		resourceId:       [32]byte{0xc0, 0xff, 0xee, 0xba, 0xbe},
		toChainId:        big.NewInt(1337),
		encryptedMessage: []byte{0xde, 0xca, 0xfb, 0xad, 0xc0, 0xff, 0xee},
	}

	log := types.Log{
		Address:     common.HexToAddress("0x6666666666666666666666666666666666666666"),
		BlockNumber: 300,
		TxHash:      common.HexToHash("0xaabbccddeeff0011"),
		TxIndex:     4,
		BlockHash:   common.HexToHash("0x1100ffeeddccbbaa"),
		Index:       3,
	}

	// Apply options
	for _, opt := range opts {
		opt(&log, &data)
	}

	// Pack the data using ABI
	abi, _ := EnygmaTeleport.EnygmaTeleportMetaData.ParseABI()
	packedData, _ := abi.Events["EnygmaTransfer"].Inputs.NonIndexed().Pack(data.encryptedMessage)

	// Set topics: event signature + indexed parameters (resourceId, toChainId)
	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("EnygmaTransfer(bytes32,uint256,bytes)")),
		data.resourceId,
		common.BigToHash(data.toChainId),
	}
	log.Data = packedData

	return log
}
