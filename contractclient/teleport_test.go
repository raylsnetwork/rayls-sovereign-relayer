package contractclient_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contractclient"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/TeleportV1"
	sharedservice "github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/shared/service"
	raylstypes "github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// packTeleportStringError builds a full revert payload for custom Solidity
// errors of the shape `error X(string msgId)`: the 4-byte error selector
// followed by the ABI-encoded `string msgId` argument. This mirrors what the
// EVM emits on revert.
func packTeleportStringError(t *testing.T, selector []byte, msgId string) []byte {
	t.Helper()
	stringType, err := abi.NewType("string", "", nil)
	require.NoError(t, err)
	packed, err := abi.Arguments{{Type: stringType}}.Pack(msgId)
	require.NoError(t, err)
	return append(append([]byte{}, selector...), packed...)
}

type StubEncryptor struct {
	encrData           string
	dataWithMessageTag contractclient.EncryptedMessages
	err                error

	spyData    []raylstypes.AtomicTeleportAdditionalData
	spyMessage []raylstypes.DispatchedMessageToPrivateHub
}

func (e *StubEncryptor) EncryptAdditionalData(_ context.Context, data []raylstypes.AtomicTeleportAdditionalData) (string, error) {
	e.spyData = data
	return e.encrData, e.err
}

func (e *StubEncryptor) EncryptMessages(
	_ context.Context,
	data []raylstypes.DispatchedMessageToPrivateHub,
	chainID *big.Int,
) (contractclient.EncryptedMessages, error) {
	e.spyMessage = data
	result := e.dataWithMessageTag
	if result.BlockNumber == nil {
		result.BlockNumber = big.NewInt(0)
	}
	return result, e.err
}

type stubExecutor struct {
	executeErr      error
	callErr         error
	callResult      []byte
	batchExecuteErr error
	batchResults    map[string]contractclient.BatchResult

	spyExecuteID       string
	spyExecuteCalldata []byte
	spyExecuteAddress  common.Address

	spyCallCalldata []byte
	spyCallAddress  common.Address

	spyBatchItems []contractclient.BatchInput
}

func (s *stubExecutor) Execute(ctx context.Context, id string, calldata []byte, address common.Address) (*ethTypes.Receipt, error) {
	s.spyExecuteID = id
	s.spyExecuteCalldata = calldata
	s.spyExecuteAddress = address

	if s.executeErr != nil {
		return nil, s.executeErr
	}

	return &ethTypes.Receipt{
		TxHash: common.HexToHash("0x1"),
		Status: 1,
	}, nil
}

func (s *stubExecutor) Sign(ctx context.Context, calldata []byte, address common.Address) (*ethTypes.Transaction, error) {
	s.spyExecuteCalldata = calldata
	s.spyExecuteAddress = address

	if s.executeErr != nil {
		return nil, s.executeErr
	}

	return ethTypes.NewTransaction(0, address, big.NewInt(0), 0, big.NewInt(0), calldata), nil
}

func (s *stubExecutor) Call(ctx context.Context, address common.Address, calldata []byte) ([]byte, error) {
	s.spyCallCalldata = calldata
	s.spyCallAddress = address

	if s.callErr != nil {
		return nil, s.callErr
	}

	return s.callResult, nil
}

func (s *stubExecutor) BatchExecute(ctx context.Context, items []contractclient.BatchInput) (map[string]contractclient.BatchResult, error) {
	s.spyBatchItems = items
	if s.batchExecuteErr != nil {
		return nil, s.batchExecuteErr
	}
	return s.batchResults, nil
}

func TestTeleportClient_SendAdditionalDataBatch(t *testing.T) {
	t.Run("encrypts additional data and sends it to the teleport contract", func(t *testing.T) {
		wantSharedIDs := []string{"example-shared-id"}
		wantEncrData := "encrypted data"

		addData := []raylstypes.AtomicTeleportAdditionalData{
			{
				SharedId: wantSharedIDs[0],
			},
		}

		address := common.HexToAddress("0x1")

		executor := &stubExecutor{}
		encr := &StubEncryptor{
			encrData: wantEncrData,
		}
		cli := contractclient.NewTeleportClient(address, executor, encr)

		err := cli.SendAdditionalDataBatch(context.Background(), wantSharedIDs, addData)
		require.Nil(t, err)

		// Check encryptor parameters
		assert.Equal(t, addData, encr.spyData)
		// We don't assert on calldata contents here, only that we attempted to execute against the right address.
		assert.Equal(t, address, executor.spyExecuteAddress)
	})

	t.Run("wraps executor errors in TeleportClientError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.TeleportClientError{}

		sharedIDs := []string{"example-shared-id"}
		addData := []raylstypes.AtomicTeleportAdditionalData{
			{
				SharedId: sharedIDs[0],
			},
		}

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{
			executeErr: wantError,
		}
		encr := &StubEncryptor{}
		cli := contractclient.NewTeleportClient(address, executor, encr)

		err := cli.SendAdditionalDataBatch(context.Background(), sharedIDs, addData)

		require.ErrorAs(t, err, &wantErrorType, "didn't wrap error in TeleportClientError")
		require.ErrorIs(t, err, wantError, "didn't wrap undelying error")
	})

	t.Run("wraps teleport contract errors in TeleportClientError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.TeleportClientError{}

		sharedIDs := []string{"example-shared-id"}
		addData := []raylstypes.AtomicTeleportAdditionalData{
			{
				SharedId: sharedIDs[0],
			},
		}

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{
			executeErr: wantError,
		}
		encr := &StubEncryptor{}
		cli := contractclient.NewTeleportClient(address, executor, encr)

		err := cli.SendAdditionalDataBatch(context.Background(), sharedIDs, addData)

		require.ErrorAs(t, err, &wantErrorType, "didn't wrap error in TeleportClientError")
		require.ErrorIs(t, err, wantError, "didn't wrap undelying error")
	})

	t.Run("wraps encryptor errors in TeleportClientError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.TeleportClientError{}

		sharedIDs := []string{"example-shared-id"}
		addData := []raylstypes.AtomicTeleportAdditionalData{
			{
				SharedId: sharedIDs[0],
			},
		}

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		encr := &StubEncryptor{
			err: wantError,
		}
		cli := contractclient.NewTeleportClient(address, executor, encr)

		err := cli.SendAdditionalDataBatch(context.Background(), sharedIDs, addData)

		require.ErrorAs(t, err, &wantErrorType, "didn't wrap error in TeleportClientError")
		require.ErrorIs(t, err, wantError, "didn't wrap undelying error")
	})
}

func TestTeleportClient_ExecuteAtomicMessageBatch(t *testing.T) {
	t.Run("encrypts additional data and sends it to the teleport contract", func(t *testing.T) {
		wantSharedIDs := []string{"example-shared-id"}
		wantEncrData := "encrypted data"

		addData := []raylstypes.AtomicTeleportAdditionalData{
			{
				SharedId: wantSharedIDs[0],
			},
		}

		address := common.HexToAddress("0x1")

		executor := &stubExecutor{}
		encr := &StubEncryptor{
			encrData: wantEncrData,
		}
		cli := contractclient.NewTeleportClient(address, executor, encr)

		err := cli.ExecuteAtomicMessageBatch(context.Background(), wantSharedIDs, addData)
		require.Nil(t, err)

		// Check encryptor parameters
		assert.Equal(t, addData, encr.spyData)
		// Ensure we executed against the correct address and produced calldata.
		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotNil(t, executor.spyExecuteCalldata)
	})

	t.Run(
		"returns sharedservice.ErrAlreadyReverted when contract reverts with TeleportV1__MessageAlreadyReverted",
		func(t *testing.T) {
			sharedIDs := []string{"example-shared-id"}
			addData := []raylstypes.AtomicTeleportAdditionalData{
				{
					SharedId: sharedIDs[0],
				},
			}

			address := common.HexToAddress("0x1")
			selector := TeleportV1.TeleportV1TeleportV1MessageAlreadyRevertedErrorID().Bytes()[:4]
			executor := &stubExecutor{
				executeErr: contractclient.NewErrorWithRevertData(packTeleportStringError(t, selector, "example-msg-id")),
			}
			encr := &StubEncryptor{}

			cli := contractclient.NewTeleportClient(address, executor, encr)

			err := cli.ExecuteAtomicMessageBatch(context.Background(), sharedIDs, addData)

			require.ErrorIs(t, err, sharedservice.ErrAlreadyReverted)
		},
	)

	t.Run("wraps executor execute errors in TeleportClientError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.TeleportClientError{}

		sharedIDs := []string{"example-shared-id"}
		addData := []raylstypes.AtomicTeleportAdditionalData{
			{
				SharedId: sharedIDs[0],
			},
		}

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{
			executeErr: wantError,
		}
		encr := &StubEncryptor{}
		cli := contractclient.NewTeleportClient(address, executor, encr)

		err := cli.ExecuteAtomicMessageBatch(context.Background(), sharedIDs, addData)

		require.ErrorAs(t, err, &wantErrorType, "didn't wrap error in TeleportClientError")
		require.ErrorIs(t, err, wantError, "didn't wrap undelying error")
	})

	t.Run("wraps revert with unrecognised revert data in TeleportClientError", func(t *testing.T) {
		wantErrorType := &contractclient.TeleportClientError{}
		revertErr := contractclient.NewErrorWithRevertData([]byte{0x01, 0x02, 0x03})

		sharedIDs := []string{"example-shared-id"}
		addData := []raylstypes.AtomicTeleportAdditionalData{
			{
				SharedId: sharedIDs[0],
			},
		}

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{
			executeErr: revertErr,
		}
		encr := &StubEncryptor{}

		cli := contractclient.NewTeleportClient(address, executor, encr)

		err := cli.ExecuteAtomicMessageBatch(context.Background(), sharedIDs, addData)

		require.ErrorAs(t, err, &wantErrorType, "didn't wrap error in TeleportClientError")
		require.ErrorIs(t, err, revertErr, "didn't wrap underlying error")
	})

	t.Run("wraps encryptor errors in TeleportClientError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.TeleportClientError{}

		sharedIDs := []string{"example-shared-id"}
		addData := []raylstypes.AtomicTeleportAdditionalData{
			{
				SharedId: sharedIDs[0],
			},
		}

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		encr := &StubEncryptor{
			err: wantError,
		}
		cli := contractclient.NewTeleportClient(address, executor, encr)

		err := cli.ExecuteAtomicMessageBatch(context.Background(), sharedIDs, addData)

		require.ErrorAs(t, err, &wantErrorType, "didn't wrap error in TeleportClientError")
		require.ErrorIs(t, err, wantError, "didn't wrap undelying error")
	})
}

func TestTeleportClient_RevertAtomicMessageBatch(t *testing.T) {
	t.Run("encrypts additional data and sends it to the teleport contract", func(t *testing.T) {
		wantSharedIDs := []string{"example-shared-id"}
		wantEncrData := "encrypted data"

		addData := []raylstypes.AtomicTeleportAdditionalData{
			{
				SharedId: wantSharedIDs[0],
			},
		}

		address := common.HexToAddress("0x1")

		executor := &stubExecutor{}
		encr := &StubEncryptor{
			encrData: wantEncrData,
		}
		cli := contractclient.NewTeleportClient(address, executor, encr)

		err := cli.RevertAtomicMessageBatch(context.Background(), wantSharedIDs, addData)
		require.Nil(t, err)

		// Check encryptor parameters
		assert.Equal(t, addData, encr.spyData)
		// Ensure we executed against the correct address and produced calldata.
		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotNil(t, executor.spyExecuteCalldata)
	})

	t.Run(
		"returns sharedservice.ErrAlreadyExecuted when contract reverts with TeleportV1__MessageAlreadyExecuted",
		func(t *testing.T) {
			sharedIDs := []string{"example-shared-id"}
			addData := []raylstypes.AtomicTeleportAdditionalData{
				{
					SharedId: sharedIDs[0],
				},
			}

			address := common.HexToAddress("0x1")
			selector := TeleportV1.TeleportV1TeleportV1MessageAlreadyExecutedErrorID().Bytes()[:4]
			executor := &stubExecutor{
				executeErr: contractclient.NewErrorWithRevertData(packTeleportStringError(t, selector, "example-msg-id")),
			}
			encr := &StubEncryptor{}

			cli := contractclient.NewTeleportClient(address, executor, encr)

			err := cli.RevertAtomicMessageBatch(context.Background(), sharedIDs, addData)

			require.ErrorIs(t, err, sharedservice.ErrAlreadyExecuted)
		},
	)

	t.Run("wraps executor execute errors in TeleportClientError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.TeleportClientError{}

		sharedIDs := []string{"example-shared-id"}
		addData := []raylstypes.AtomicTeleportAdditionalData{
			{
				SharedId: sharedIDs[0],
			},
		}

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{
			executeErr: wantError,
		}
		encr := &StubEncryptor{}
		cli := contractclient.NewTeleportClient(address, executor, encr)

		err := cli.RevertAtomicMessageBatch(context.Background(), sharedIDs, addData)

		require.ErrorAs(t, err, &wantErrorType, "didn't wrap error in TeleportClientError")
		require.ErrorIs(t, err, wantError, "didn't wrap undelying error")
	})

	t.Run("wraps revert with unrecognised revert data in TeleportClientError", func(t *testing.T) {
		wantErrorType := &contractclient.TeleportClientError{}
		revertErr := contractclient.NewErrorWithRevertData([]byte{0x01, 0x02, 0x03})

		sharedIDs := []string{"example-shared-id"}
		addData := []raylstypes.AtomicTeleportAdditionalData{
			{
				SharedId: sharedIDs[0],
			},
		}

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{
			executeErr: revertErr,
		}
		encr := &StubEncryptor{}

		cli := contractclient.NewTeleportClient(address, executor, encr)

		err := cli.RevertAtomicMessageBatch(context.Background(), sharedIDs, addData)

		require.ErrorAs(t, err, &wantErrorType, "didn't wrap error in TeleportClientError")
		require.ErrorIs(t, err, revertErr, "didn't wrap underlying error")
	})

	t.Run("wraps encryptor errors in TeleportClientError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.TeleportClientError{}

		sharedIDs := []string{"example-shared-id"}
		addData := []raylstypes.AtomicTeleportAdditionalData{
			{
				SharedId: sharedIDs[0],
			},
		}

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		encr := &StubEncryptor{
			err: wantError,
		}
		cli := contractclient.NewTeleportClient(address, executor, encr)

		err := cli.RevertAtomicMessageBatch(context.Background(), sharedIDs, addData)

		require.ErrorAs(t, err, &wantErrorType, "didn't wrap error in TeleportClientError")
		require.ErrorIs(t, err, wantError, "didn't wrap undelying error")
	})
}

func TestTeleportClient_StoreEncryptedDataBatch(t *testing.T) {
	t.Run("encrypts messages for the commit chain and sends them", func(t *testing.T) {
		wantEncrBatch := TeleportV1.TeleportV1dataBatch{
			MessageTag: "example-finger-print",
			Data:       []byte("example-encrypted-data"),
		}
		sharedIDs := []string{"example-shared-id"}
		messagesToPrivateHub := []raylstypes.DispatchedMessageToPrivateHub{
			{
				SharedId: "example-shared-id",
			},
		}
		dataWithMessageTag := contractclient.EncryptedMessages{
			MessageTag:  wantEncrBatch.MessageTag,
			Data:        wantEncrBatch.Data,
			BlockNumber: big.NewInt(100),
		}

		address := common.HexToAddress("0x1")

		executor := &stubExecutor{}
		encr := &StubEncryptor{
			dataWithMessageTag: dataWithMessageTag,
		}
		cli := contractclient.NewTeleportClient(address, executor, encr)

		hash, err := cli.StoreEncryptedDataBatch(context.Background(), sharedIDs, messagesToPrivateHub, new(big.Int).SetUint64(1337))
		require.Nil(t, err)
		require.NotEqual(t, hash, [32]byte{}, "should return a non-zero hash")

		// Check that we encrypted the correct data
		require.Equal(t, messagesToPrivateHub, encr.spyMessage)
		// We can't easily inspect calldata without ABI decoding; just assert we executed.
		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotNil(t, executor.spyExecuteCalldata)
	})

	t.Run("wraps auth generator errors in TeleportClientError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.TeleportClientError{}

		sharedIDs := []string{"example-shared-id"}
		messagesToPrivateHub := []raylstypes.DispatchedMessageToPrivateHub{
			{
				SharedId: "example-shared-id",
			},
		}

		address := common.HexToAddress("0x1")

		executor := &stubExecutor{
			executeErr: wantError,
		}
		encr := &StubEncryptor{}
		cli := contractclient.NewTeleportClient(address, executor, encr)

		_, err := cli.StoreEncryptedDataBatch(context.Background(), sharedIDs, messagesToPrivateHub, new(big.Int).SetUint64(1337))

		require.ErrorAs(t, err, &wantErrorType, "didn't wrap error in TeleportClientError")
		require.ErrorIs(t, err, wantError, "didn't wrap undelying error")
	})

	t.Run("wraps teleport contract errors in TeleportClientError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.TeleportClientError{}

		sharedIDs := []string{"example-shared-id"}
		messagesToPrivateHub := []raylstypes.DispatchedMessageToPrivateHub{
			{
				SharedId: "example-shared-id",
			},
		}

		address := common.HexToAddress("0x1")

		executor := &stubExecutor{
			executeErr: wantError,
		}
		encr := &StubEncryptor{}
		cli := contractclient.NewTeleportClient(address, executor, encr)

		_, err := cli.StoreEncryptedDataBatch(context.Background(), sharedIDs, messagesToPrivateHub, new(big.Int).SetUint64(1337))

		require.ErrorAs(t, err, &wantErrorType, "didn't wrap error in TeleportClientError")
		require.ErrorIs(t, err, wantError, "didn't wrap undelying error")
	})

	t.Run("wraps encryptor errors in TeleportClientError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.TeleportClientError{}

		sharedIDs := []string{"example-shared-id"}
		messagesToPrivateHub := []raylstypes.DispatchedMessageToPrivateHub{
			{
				SharedId: "example-shared-id",
			},
		}

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		encr := &StubEncryptor{
			err: wantError,
		}
		cli := contractclient.NewTeleportClient(address, executor, encr)

		_, err := cli.StoreEncryptedDataBatch(context.Background(), sharedIDs, messagesToPrivateHub, new(big.Int).SetUint64(1337))

		require.ErrorAs(t, err, &wantErrorType, "didn't wrap error in TeleportClientError")
		require.ErrorIs(t, err, wantError, "didn't wrap undelying error")
	})
}

func TestTeleportClient_GetAtomicMessageStatuses(t *testing.T) {
	t.Run("gets atomic message statuses from teleport contract", func(t *testing.T) {
		wantSharedIDs := []string{"pending-shared-id", "executed-shared-id", "rejected-shared-id", "reverted-shared-id"}
		wantSUMs := []raylstypes.AtomicStatusUpdateMessage{
			{
				SharedID: wantSharedIDs[0],
				Status:   raylstypes.AtomicPendingStatus,
			},
			{
				SharedID: wantSharedIDs[1],
				Status:   raylstypes.AtomicExecutedStatus,
			},
			{
				SharedID: wantSharedIDs[2],
				Status:   raylstypes.AtomicRejectedStatus,
			},
			{
				SharedID: wantSharedIDs[3],
				Status:   raylstypes.AtomicRevertedStatus,
			},
		}

		teleportSUMs := []TeleportV1.TeleportV1MessageStatusResult{
			{
				MsgId:  wantSUMs[0].SharedID,
				Status: raylstypes.AtomicPending.String(),
			},
			{
				MsgId:  wantSUMs[1].SharedID,
				Status: raylstypes.AtomicExecuted.String(),
			},
			{
				MsgId:  wantSUMs[2].SharedID,
				Status: raylstypes.AtomicRejected.String(),
			},
			{
				MsgId:  wantSUMs[3].SharedID,
				Status: raylstypes.AtomicReverted.String(),
			},
		}

		address := common.HexToAddress("0x1")

		parsed, err := TeleportV1.TeleportV1MetaData.ParseABI()
		require.NoError(t, err)
		method := parsed.Methods["getAtomicMessageStatuses"]
		callResult, err := method.Outputs.Pack(teleportSUMs)
		require.NoError(t, err)

		executor := &stubExecutor{
			callResult: callResult,
		}
		encr := &StubEncryptor{}
		cli := contractclient.NewTeleportClient(address, executor, encr)

		gotSUMs, err := cli.GetAtomicMessageStatuses(context.Background(), wantSharedIDs)
		require.Nil(t, err)

		assert.ElementsMatch(t, wantSUMs, gotSUMs, "didn't return correct atomic message statuses")
	})

	t.Run("returns ErrUnknownStatus in case of an unknown status from contract", func(t *testing.T) {
		wantErr := contractclient.ErrUnknownStatus

		sharedIDs := []string{"example-shared-id"}
		teleportSUMs := []TeleportV1.TeleportV1MessageStatusResult{
			{
				MsgId:  sharedIDs[0],
				Status: "Unknown Status",
			},
		}

		address := common.HexToAddress("0x1")

		parsed, err := TeleportV1.TeleportV1MetaData.ParseABI()
		require.NoError(t, err)
		method := parsed.Methods["getAtomicMessageStatuses"]
		callResult, err := method.Outputs.Pack(teleportSUMs)
		require.NoError(t, err)

		executor := &stubExecutor{
			callResult: callResult,
		}
		encr := &StubEncryptor{}
		cli := contractclient.NewTeleportClient(address, executor, encr)

		_, gotErr := cli.GetAtomicMessageStatuses(context.Background(), sharedIDs)

		assert.ErrorIs(t, gotErr, wantErr)
	})

	t.Run("wraps teleport contract errors in TeleportClientError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.TeleportClientError{}

		sharedIDs := []string{"example-shared-id"}

		address := common.HexToAddress("0x1")

		executor := &stubExecutor{
			callErr: wantError,
		}
		encr := &StubEncryptor{}
		cli := contractclient.NewTeleportClient(address, executor, encr)

		_, err := cli.GetAtomicMessageStatuses(context.Background(), sharedIDs)

		require.ErrorAs(t, err, &wantErrorType, "didn't wrap error in TeleportClientError")
		require.ErrorIs(t, err, wantError, "didn't wrap undelying error")
	})
}
