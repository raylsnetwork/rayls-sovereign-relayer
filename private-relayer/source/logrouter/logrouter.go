package logrouter

import (
	"context"
	"fmt"
	"log/slog"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

//go:generate moq --pkg logrouter_test -out logrouter_moq_test.go . EndpointMQ EnygmaMQ

type Block struct {
	Number uint64
	Logs   []ethTypes.Log
}

func (b Block) GetID() string {
	return fmt.Sprintf("%d", b.Number)
}

type EndpointMQ interface {
	Push(context.Context, Block) error
}

type EnygmaMQ interface {
	Push(context.Context, Block) error
}

// Split the incoming logs in two message queues - > one for events coming from the
// endpoint contracts, and another for events coming from the enygma contracts.

// EndpointMQ -> Block 1 [10 logs], Block 2 [5 logs]
// receiver fetch(2) -> Blocks 1 [10 logs], Block 2 [5 logs]
type LogRouter struct {
	endpointAddress common.Address
	enygmaAddress   common.Address

	endpointMQ EndpointMQ
	enygmaMQ   EnygmaMQ
}

func New(endpointAddress, enygmaAddress common.Address, endpointMQ EndpointMQ, enygmaMQ EndpointMQ) *LogRouter {
	return &LogRouter{
		endpointAddress: endpointAddress,
		enygmaAddress:   enygmaAddress,

		endpointMQ: endpointMQ,
		enygmaMQ:   enygmaMQ,
	}
}

func (l *LogRouter) Handle(ctx context.Context, logs []ethTypes.Log) error {
	endpointLogs := make(map[uint64][]ethTypes.Log)
	enygmaLogs := make(map[uint64][]ethTypes.Log)

	for _, log := range logs {
		switch log.Address {
		case l.endpointAddress:
			endpointLogs[log.BlockNumber] = append(endpointLogs[log.BlockNumber], log)
		case l.enygmaAddress:
			enygmaLogs[log.BlockNumber] = append(enygmaLogs[log.BlockNumber], log)
		}
	}

	// TODO: Rather than using a map, this logic can be implemented with
	// better time complexity and less additional logic using a splay tree.
	endpointBlockNums := make([]uint64, 0, len(endpointLogs))
	for key := range endpointLogs {
		endpointBlockNums = append(endpointBlockNums, key)
	}
	sort.Slice(endpointBlockNums, func(i, j int) bool {
		return endpointBlockNums[i] < endpointBlockNums[j]
	})

	for _, blockNum := range endpointBlockNums {
		err := l.endpointMQ.Push(ctx, Block{
			Number: blockNum,
			Logs:   endpointLogs[blockNum],
		})
		if err == nil {
			slog.Info("Pushed logs to endpoint queue", slog.Any("count", len(endpointLogs)))
			slog.Debug("Pushed message with block number to endpoint queue", slog.Any("block_number", blockNum))
		} else {
			slog.Error("Failed to push endpoint logs to queue", slog.Any("error", err))
		}
	}

	enygmaBlockNums := make([]uint64, 0, len(enygmaLogs))
	for key := range enygmaLogs {
		enygmaBlockNums = append(enygmaBlockNums, key)
	}
	sort.Slice(enygmaBlockNums, func(i, j int) bool {
		return enygmaBlockNums[i] < enygmaBlockNums[j]
	})

	for _, blockNum := range enygmaBlockNums {
		err := l.enygmaMQ.Push(ctx, Block{
			Number: blockNum,
			Logs:   enygmaLogs[blockNum],
		})
		if err == nil {
			slog.Info("Pushed logs to enygma queue", slog.Any("count", len(enygmaLogs)))
		} else {
			slog.Error("Failed to push enygma logs to queue", slog.Any("error", err))
		}
	}

	return nil
}
