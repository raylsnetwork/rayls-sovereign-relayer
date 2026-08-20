package batcher

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

type Message struct {
	ID       string
	Address  common.Address
	Calldata []byte
}

type Publisher interface {
	PushBatch(ctx context.Context, msgs []types.TxRequest) error
}

type Batcher struct {
	messageType string
	publisher   Publisher
}

func NewBatcher(messageType string, pub Publisher) *Batcher {
	return &Batcher{
		messageType: messageType,
		publisher:   pub,
	}

}

func (b *Batcher) Send(ctx context.Context, msgSlice []Message) error {
	requests := make([]types.TxRequest, 0, len(msgSlice))
	for _, msg := range msgSlice {
		requests = append(requests, types.TxRequest{
			CorrelationID: msg.ID,
			MessageType:   b.messageType,
			Address:       msg.Address,
			Calldata:      msg.Calldata,
		})
	}

	return b.publisher.PushBatch(ctx, requests)
}
