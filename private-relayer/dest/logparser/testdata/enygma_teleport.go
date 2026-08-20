package testdata

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EnygmaTeleport"
)

// EnygmaTransferLogOption represents a modification to an EnygmaTransfer log fixture.
type EnygmaTransferLogOption func(*types.Log, *enygmaTransferData)

type enygmaTransferData struct {
	resourceID       [32]byte
	toChainId        *big.Int
	encryptedMessage []byte
	messageTag       *big.Int
	blockNumber      *big.Int
	anonymitySet     []*big.Int
	arrayHashSecrets []*big.Int
}

// WithTransferResourceID sets the ResourceID field.
func WithTransferResourceID(resourceID [32]byte) EnygmaTransferLogOption {
	return func(log *types.Log, data *enygmaTransferData) {
		data.resourceID = resourceID
	}
}

// WithTransferToChainId sets the ToChainId field.
func WithTransferToChainId(chainId *big.Int) EnygmaTransferLogOption {
	return func(log *types.Log, data *enygmaTransferData) {
		data.toChainId = chainId
	}
}

// WithTransferEncryptedMessage sets the EncryptedMessage field.
func WithTransferEncryptedMessage(msg []byte) EnygmaTransferLogOption {
	return func(log *types.Log, data *enygmaTransferData) {
		data.encryptedMessage = msg
	}
}

// WithTransferBlockNumber sets the BlockNumber field.
func WithTransferBlockNumber(blockNumber uint64) EnygmaTransferLogOption {
	return func(log *types.Log, data *enygmaTransferData) {
		log.BlockNumber = blockNumber
	}
}

// WithTransferTxHash sets the TxHash field.
func WithTransferTxHash(hash common.Hash) EnygmaTransferLogOption {
	return func(log *types.Log, data *enygmaTransferData) {
		log.TxHash = hash
	}
}

// WithTransferMessageTag sets the MessageTag field.
func WithTransferMessageTag(messageTag *big.Int) EnygmaTransferLogOption {
	return func(log *types.Log, data *enygmaTransferData) {
		data.messageTag = messageTag
	}
}

// WithTransferEventBlockNumber sets the blockNumber field in the event data.
func WithTransferEventBlockNumber(blockNumber *big.Int) EnygmaTransferLogOption {
	return func(log *types.Log, data *enygmaTransferData) {
		data.blockNumber = blockNumber
	}
}

// WithTransferAnonymitySet sets the AnonymitySet field.
func WithTransferAnonymitySet(anonymitySet []*big.Int) EnygmaTransferLogOption {
	return func(log *types.Log, data *enygmaTransferData) {
		data.anonymitySet = anonymitySet
	}
}

// WithTransferArrayHashSecrets sets the ArrayHashSecrets field.
func WithTransferArrayHashSecrets(arrayHashSecrets []*big.Int) EnygmaTransferLogOption {
	return func(log *types.Log, data *enygmaTransferData) {
		data.arrayHashSecrets = arrayHashSecrets
	}
}

// NewEnygmaTransferLogWith creates an EnygmaTransfer log with default values and applies options.
func NewEnygmaTransferLogWith(opts ...EnygmaTransferLogOption) types.Log {
	data := enygmaTransferData{
		resourceID:       [32]byte{0x01, 0x02, 0x03},
		toChainId:        big.NewInt(12345),
		encryptedMessage: []byte{0xde, 0xad, 0xbe, 0xef},
		messageTag:       big.NewInt(1),
		blockNumber:      big.NewInt(100),
		anonymitySet:     []*big.Int{big.NewInt(111), big.NewInt(222)},
		arrayHashSecrets: []*big.Int{big.NewInt(333), big.NewInt(444)},
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 100,
		TxHash:      common.HexToHash("0xabcdef1234567890"),
		TxIndex:     1,
		BlockHash:   common.HexToHash("0x123456"),
		Index:       0,
	}

	for _, opt := range opts {
		opt(&log, &data)
	}

	abi, _ := EnygmaTeleport.EnygmaTeleportMetaData.ParseABI()
	// Non-indexed parameters: encryptedMessage, messageTag, blockNumber, anonymitySet, arrayHashSecrets, toChainId
	packedData, _ := abi.Events["EnygmaTransfer"].Inputs.NonIndexed().Pack(
		data.encryptedMessage,
		data.messageTag,
		data.blockNumber,
		data.anonymitySet,
		data.arrayHashSecrets,
		data.toChainId,
	)

	// Event: EnygmaTransfer(bytes32 indexed resourceId, bytes encryptedMessage, uint256 messageTag, uint256 blockNumber, uint256[] anonymitySet, uint256[] arrayHashSecrets, uint256 toChainId)
	// Event signature hash: 0x2b925a0de145bd51c571f8da5994d78f6cce29dc9f967e20a425139699375683
	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("EnygmaTransfer(bytes32,bytes,uint256,uint256,uint256[],uint256[],uint256)")),
		common.BytesToHash(data.resourceID[:]),
	}
	log.Data = packedData

	return log
}

// BalancesFinalizedLogOption represents a modification to a BalancesFinalized log fixture.
type BalancesFinalizedLogOption func(*types.Log, *balancesFinalizedData)

type balancesFinalizedData struct {
	resourceID           [32]byte
	finalizedBlockNumber *big.Int
	pendingBlockNumber   *big.Int
	balances             []EnygmaTeleport.IEnygmaV1EnygmaPointWithChainId
}

// WithFinalizedResourceID sets the ResourceID field.
func WithFinalizedResourceID(resourceID [32]byte) BalancesFinalizedLogOption {
	return func(log *types.Log, data *balancesFinalizedData) {
		data.resourceID = resourceID
	}
}

// WithFinalizedBlockNumbers sets both block number fields.
func WithFinalizedBlockNumbers(finalized, pending *big.Int) BalancesFinalizedLogOption {
	return func(log *types.Log, data *balancesFinalizedData) {
		data.finalizedBlockNumber = finalized
		data.pendingBlockNumber = pending
	}
}

// WithFinalizedBalances sets the balances array.
func WithFinalizedBalances(balances []EnygmaTeleport.IEnygmaV1EnygmaPointWithChainId) BalancesFinalizedLogOption {
	return func(log *types.Log, data *balancesFinalizedData) {
		data.balances = balances
	}
}

// WithFinalizedLogBlockNumber sets the log's BlockNumber.
func WithFinalizedLogBlockNumber(blockNumber uint64) BalancesFinalizedLogOption {
	return func(log *types.Log, data *balancesFinalizedData) {
		log.BlockNumber = blockNumber
	}
}

// NewBalancesFinalizedLogWith creates a BalancesFinalized log with default values and applies options.
func NewBalancesFinalizedLogWith(opts ...BalancesFinalizedLogOption) types.Log {
	data := balancesFinalizedData{
		resourceID:           [32]byte{0x01, 0x02, 0x03},
		finalizedBlockNumber: big.NewInt(100),
		pendingBlockNumber:   big.NewInt(110),
		balances: []EnygmaTeleport.IEnygmaV1EnygmaPointWithChainId{
			{
				ChainId: big.NewInt(12345),
				C1:      big.NewInt(1000),
				C2:      big.NewInt(2000),
			},
		},
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 100,
		TxHash:      common.HexToHash("0xabcdef1234567890"),
		TxIndex:     1,
		BlockHash:   common.HexToHash("0x123456"),
		Index:       0,
	}

	for _, opt := range opts {
		opt(&log, &data)
	}

	abi, _ := EnygmaTeleport.EnygmaTeleportMetaData.ParseABI()
	// Non-indexed parameters: pendingBlockNumber, balances
	packedData, _ := abi.Events["BalancesFinalized"].Inputs.NonIndexed().Pack(
		data.pendingBlockNumber,
		data.balances,
	)

	// Event: BalancesFinalized(bytes32 indexed resourceId, uint256 indexed finalizedBlockNumber, uint256 pendingBlockNumber, (uint256,uint256,uint256)[] balances)
	// Event signature hash: 0xcab9daef74ea3374e6a930e337db120037cc08e628a458f0cd2c756dc74946b9
	log.Topics = []common.Hash{
		common.HexToHash("0xcab9daef74ea3374e6a930e337db120037cc08e628a458f0cd2c756dc74946b9"),
		common.BytesToHash(data.resourceID[:]),
		common.BigToHash(data.finalizedBlockNumber),
	}
	log.Data = packedData

	return log
}
