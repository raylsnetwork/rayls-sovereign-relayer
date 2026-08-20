package logrouter

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/backoff"
)

const backoffMultiplier = 2.0

type Block struct {
	Number uint64
	Logs   []ethTypes.Log
}

func (b Block) GetID() string {
	return fmt.Sprintf("%d", b.Number)
}

//go:generate moq --pkg logrouter_test -out logrouter_mock_test.go . RouterMQ
type RouterMQ interface {
	Push(context.Context, Block) error
}

type LogRouterConfig struct {
	EndpointAddress     common.Address
	TeleportAddress     common.Address
	EnygmaAddress       common.Address
	DvpAddress          common.Address
	AuditManagerAddress common.Address
}

type LogRouter struct {
	endpontAddress common.Address
	endpointMQ     RouterMQ

	teleportAddress common.Address
	teleportMQ      RouterMQ

	enygmaAddress common.Address
	enygmaMQ      RouterMQ

	dvpAddress common.Address
	dvpMQ      RouterMQ

	auditManagerAddress common.Address
	auditManagerMQ      RouterMQ

	backoff backoff.Strategy
}

func New(
	config LogRouterConfig,
	endpointMQ RouterMQ,
	teleportMQ RouterMQ,
	enygmaMQ RouterMQ,
	dvpMQ RouterMQ,
	auditManagerMQ RouterMQ,
) *LogRouter {
	exp, _ := backoff.NewExponential(time.Second, backoffMultiplier, time.Minute)
	return NewWithCustomBackoff(config, endpointMQ, teleportMQ, enygmaMQ, dvpMQ, auditManagerMQ, exp)
}

func NewWithCustomBackoff(
	config LogRouterConfig,
	endpointMQ RouterMQ,
	teleportMQ RouterMQ,
	enygmaMQ RouterMQ,
	dvpMQ RouterMQ,
	auditManagerMQ RouterMQ,
	backoff backoff.Strategy,
) *LogRouter {
	return &LogRouter{
		endpontAddress: config.EndpointAddress,
		endpointMQ:     endpointMQ,

		teleportAddress: config.TeleportAddress,
		teleportMQ:      teleportMQ,

		enygmaAddress: config.EnygmaAddress,
		enygmaMQ:      enygmaMQ,

		dvpAddress: config.DvpAddress,
		dvpMQ:      dvpMQ,

		auditManagerAddress: config.AuditManagerAddress,
		auditManagerMQ:      auditManagerMQ,

		backoff: backoff,
	}
}

func (l *LogRouter) Handle(ctx context.Context, logs []ethTypes.Log) error {
	slog.Debug("LogRouter handling logs", slog.Int("total_log_count", len(logs)))

	endpointBlocks := map[uint64][]ethTypes.Log{}
	teleportBlocks := map[uint64][]ethTypes.Log{}
	enygmaBlocks := map[uint64][]ethTypes.Log{}
	dvpBlocks := map[uint64][]ethTypes.Log{}
	auditManagerBlocks := map[uint64][]ethTypes.Log{}

	for _, log := range logs {
		switch log.Address {
		case l.endpontAddress:
			endpointBlocks[log.BlockNumber] = append(endpointBlocks[log.BlockNumber], log)
		case l.teleportAddress:
			teleportBlocks[log.BlockNumber] = append(teleportBlocks[log.BlockNumber], log)
		case l.enygmaAddress:
			enygmaBlocks[log.BlockNumber] = append(enygmaBlocks[log.BlockNumber], log)
		case l.dvpAddress:
			dvpBlocks[log.BlockNumber] = append(dvpBlocks[log.BlockNumber], log)
		case l.auditManagerAddress:
			auditManagerBlocks[log.BlockNumber] = append(auditManagerBlocks[log.BlockNumber], log)
		}
	}

	slog.Info("Log routing summary",
		slog.Int("endpoint_blocks", len(endpointBlocks)),
		slog.Int("teleport_blocks", len(teleportBlocks)),
		slog.Int("enygma_blocks", len(enygmaBlocks)),
		slog.Int("dvp_blocks", len(dvpBlocks)),
		slog.Int("audit_manager_blocks", len(auditManagerBlocks)))

	l.pushBlocksToMQ(ctx, l.endpointMQ, endpointBlocks, "endpoint")
	l.pushBlocksToMQ(ctx, l.teleportMQ, teleportBlocks, "teleport")
	l.pushBlocksToMQ(ctx, l.enygmaMQ, enygmaBlocks, "enygma")
	l.pushBlocksToMQ(ctx, l.dvpMQ, dvpBlocks, "dvp")
	l.pushBlocksToMQ(ctx, l.auditManagerMQ, auditManagerBlocks, "audit_manager")

	return nil
}

func (l *LogRouter) pushBlocksToMQ(
	ctx context.Context,
	mq RouterMQ,
	blocks map[uint64][]ethTypes.Log,
	mqName string,
) {
	if len(blocks) == 0 {
		return
	}

	blockNumbers := slices.Collect(maps.Keys(blocks))
	slices.Sort(blockNumbers)

	slog.Debug("Pushing blocks to MQ",
		slog.String("mq", mqName),
		slog.Int("block_count", len(blockNumbers)))

	successCount := 0
	for _, blockNumber := range blockNumbers {
		block := Block{
			Number: blockNumber,
			Logs:   blocks[blockNumber],
		}

		err := l.backoff.Do(ctx, 100, func() error {
			return mq.Push(ctx, block)
		})

		if err != nil {
			slog.Error("Failed to push block to MQ after retries",
				slog.String("mq", mqName),
				slog.Uint64("blockNumber", blockNumber),
				slog.Int("log_count", len(block.Logs)),
				slog.Any("error", err))
			// Continue to next block even if this one failed
		} else {
			successCount++
		}
	}

	if successCount > 0 {
		slog.Debug("Successfully pushed blocks to MQ",
			slog.String("mq", mqName),
			slog.Int("count", successCount))
	}
}
