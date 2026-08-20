package testdata

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/EndpointV1"
	"github.com/raylsnetwork/rayls-sovereign-relayer/msgqueue"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/source/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

// CrossChainMessage helpers

// CrossChainMessageOption represents a modification to a CrossChainMessage fixture.
type CrossChainMessageOption func(*service.CrossChainMessage)

// NewCrossChainMessage creates a CrossChainMessage with sensible defaults.
func NewCrossChainMessage() service.CrossChainMessage {
	return service.CrossChainMessage{
		ID: "1000-42-7",

		MessageID: common.HexToHash("0xc0febabe"),
		From:      common.HexToAddress("0xdeadc0de"),
		ToChainID: big.NewInt(1337),
		To:        common.HexToAddress("0xdeadbeef"),
		Data: EndpointV1.RaylsMessage{
			MessageMetadata: EndpointV1.RaylsMessageMetadata{
				RevertPayloadDataSender:   common.Hex2Bytes("cafefeed"),
				RevertPayloadDataReceiver: common.Hex2Bytes("deadbabe"),
				LockData:                  common.Hex2Bytes("c0010fff"),
			},
		},

		BlockHash: common.HexToHash("0xdead10cc"),
		TxHash:    common.HexToHash("0xc0cac01a"),

		BlockNumber: 1000,
		TxIdx:       42,
		LogIdx:      7,
	}
}

// WithMessageIDOpt sets the MessageID field on a CrossChainMessage when using options.
func WithRevertPayloadSednerOpts(data []byte) CrossChainMessageOption {
	return func(m *service.CrossChainMessage) {
		m.Data.MessageMetadata.RevertPayloadDataSender = data
	}
}

// WithMessageIDOpt sets the MessageID field on a CrossChainMessage when using options.
func WithMessageIDOpt(id common.Hash) CrossChainMessageOption {
	return func(m *service.CrossChainMessage) {
		m.MessageID = id
	}
}

// WithToChainIDOpt sets the ToChainID field on a CrossChainMessage when using options.
func WithToChainIDOpt(id *big.Int) CrossChainMessageOption {
	return func(m *service.CrossChainMessage) {
		m.ToChainID = id
	}
}

// NewCrossChainMessageWith creates a CrossChainMessage using the default
// fixture and then applies any provided options.
func NewCrossChainMessageWith(opts ...CrossChainMessageOption) service.CrossChainMessage {
	msg := NewCrossChainMessage()
	for _, opt := range opts {
		opt(&msg)
	}
	return msg
}

// QueueMessageBuilder builds test msgqueue.Message[service.CrossChainMessage] instances
type QueueMessageBuilder struct {
	msg     service.CrossChainMessage
	ackFunc func(context.Context) error
}

// NewQueueMessage creates a QueueMessageBuilder with sensible defaults
func NewQueueMessage() *QueueMessageBuilder {
	return &QueueMessageBuilder{
		msg: NewCrossChainMessage(),
		ackFunc: func(context.Context) error {
			return nil
		},
	}
}

func (b *QueueMessageBuilder) WithCrossChainMessage(msg service.CrossChainMessage) *QueueMessageBuilder {
	b.msg = msg
	return b
}

func (b *QueueMessageBuilder) WithAckFunc(ackFunc func(context.Context) error) *QueueMessageBuilder {
	b.ackFunc = ackFunc
	return b
}

// Convenience methods that delegate to the underlying CrossChainMessage
func (b *QueueMessageBuilder) WithToChainID(toChainID *big.Int) *QueueMessageBuilder {
	b.msg.ToChainID = toChainID
	return b
}

func (b *QueueMessageBuilder) WithID(id string) *QueueMessageBuilder {
	b.msg.ID = id
	return b
}

func (b *QueueMessageBuilder) WithMessageID(messageID common.Hash) *QueueMessageBuilder {
	b.msg.MessageID = messageID
	return b
}

func (b *QueueMessageBuilder) WithTxHash(txHash common.Hash) *QueueMessageBuilder {
	b.msg.TxHash = txHash
	return b
}

func (b *QueueMessageBuilder) WithBlockHash(blockHash common.Hash) *QueueMessageBuilder {
	b.msg.BlockHash = blockHash
	return b
}

func (b *QueueMessageBuilder) Build() msgqueue.Message[service.CrossChainMessage] {
	return msgqueue.Message[service.CrossChainMessage]{
		V:   b.msg,
		Ack: b.ackFunc,
	}
}

// NewBlock creates a test Ethereum block with sensible defaults.
func NewBlock() *ethTypes.Block {
	header := &ethTypes.Header{
		Time:       uint64(time.Now().Unix()),
		ParentHash: common.HexToHash("0xfeedface"),
		TxHash:     common.HexToHash("0xd00d2bad"),
	}
	return ethTypes.NewBlock(header, nil, nil, nil)
}

// DispatchedMessageOption represents a modification to a DispatchedMessageToPrivateHub fixture.
type DispatchedMessageOption func(*types.DispatchedMessageToPrivateHub)

// WithAtomicOpt sets the IsAtomic flag on a DispatchedMessageToPrivateHub when using options.
func WithAtomicOpt(isAtomic bool) DispatchedMessageOption {
	return func(m *types.DispatchedMessageToPrivateHub) {
		m.IsAtomic = isAtomic
	}
}

// NewDispatchedMessage creates a DispatchedMessageToPrivateHub from a CrossChainMessage
// and Ethereum block, applying any provided options.
func NewDispatchedMessage(ccMsg service.CrossChainMessage, ourChainID *big.Int, block *ethTypes.Block, opts ...DispatchedMessageOption) types.DispatchedMessageToPrivateHub {
	msg := types.DispatchedMessageToPrivateHub{
		MessageId:   ccMsg.MessageID,
		FromChainId: ourChainID,
		From:        ccMsg.From,
		ToChainId:   ccMsg.ToChainID,
		To:          ccMsg.To,
		Data:        ccMsg.Data,

		TransactionType: types.Transfer,
		IsAtomic:        true,

		BlockNumber: ccMsg.BlockNumber,
		BlockHash:   ccMsg.BlockHash,
		LogIdx:      ccMsg.LogIdx,

		ParentHash: block.ParentHash().Hex(),

		TxHashSource:          ccMsg.TxHash,
		TxHashSourceStatus:    1,
		TxHashSourceTimestamp: block.Time(),

		ResourceId:   ccMsg.Data.MessageMetadata.ResourceId,
		TokenAddress: ccMsg.Data.MessageMetadata.TransferMetadata.TokenAddress,

		Proofs:      common.Hex2Bytes("fee1dead"),
		TxLocation:  ccMsg.TxIdx,
		TxTrieProof: block.TxHash(),
	}

	for _, opt := range opts {
		opt(&msg)
	}

	return msg
}

// NewTransaction creates a Transaction with sensible defaults from a DispatchedMessageToPrivateHub.
func NewTransaction(msg types.DispatchedMessageToPrivateHub, batchHash common.Hash) types.Transaction {
	return types.Transaction{
		SharedID:            msg.SharedId,
		BatchID:             msg.BatchId,
		BatchPrivateHubHash: batchHash,

		MsgID:               msg.MessageId,
		FromChainID:         msg.FromChainId,
		FromContractAddress: msg.TokenAddress.String(),
		FromUserAddress:     msg.From.String(),
		ToChainID:           msg.ToChainId,

		ResourceID: msg.ResourceId.String(),
		IsAtomic:   msg.IsAtomic,

		ParentHash:  msg.ParentHash,
		BlockNumber: msg.BlockNumber,
		TxHash:      msg.TxHashSource.String(),
		LogIndex:    msg.LogIdx,

		UpdatedAt: time.Unix(int64(msg.TxHashSourceTimestamp), 0),

		TransferID:     msg.Data.MessageMetadata.TransferMetadata.Id.String(),
		TransferAmount: msg.Data.MessageMetadata.TransferMetadata.Amount.String(),
	}
}
