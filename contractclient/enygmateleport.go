package contractclient

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EnygmaTeleport"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

type enygmaTeleportEncryptor interface {
	EncryptEnygmaTransferBatchCompleted(ctx context.Context, messages []types.EnygmaTransferCompleted) ([]byte, error)
}

type EnygmaTeleportClient struct {
	address  common.Address
	contract *EnygmaTeleport.EnygmaTeleport
	executor Executor
	encr     enygmaTeleportEncryptor
}

func NewEnygmaTeleportClient(
	address common.Address,
	executor Executor,
	encr enygmaTeleportEncryptor,
) *EnygmaTeleportClient {
	return &EnygmaTeleportClient{
		address:  address,
		contract: EnygmaTeleport.NewEnygmaTeleport(),
		executor: executor,
		encr:     encr,
	}
}

func (c *EnygmaTeleportClient) SendTransferCompleted(
	ctx context.Context,
	messages []types.EnygmaTransferCompleted,
) error {
	encryptedData, err := c.encr.EncryptEnygmaTransferBatchCompleted(ctx, messages)
	if err != nil {
		return WrapInEnygmaTeleportClientError("failed to encrypt enygma transfer batch completed", err)
	}

	msgIDs := make([]string, len(messages))
	for i, m := range messages {
		msgIDs[i] = m.MessageId
	}

	calldata := c.contract.PackEnygmaTransferCompleted(encryptedData)

	// Best-effort idempotency key: opportunistic batch membership is not stable
	// across restarts, so this only suppresses an exact-same-batch resend. The
	// real double-effect guard is on-chain (per-message dedup at the destination).
	_, err = c.executor.Execute(ctx, IDFor("enygmateleport.SendTransferCompleted", HashIDs(msgIDs)), calldata, c.address)
	if err != nil {
		return WrapInEnygmaTeleportClientError("failed to send enygma transfer completed", withstack.Wrap(err))
	}

	return nil
}
