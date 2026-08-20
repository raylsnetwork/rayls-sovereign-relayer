package contractclient_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/Dvp"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/dvp"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createDvpTestProofReceipt() *dvp.ProofReceipt {
	return &dvp.ProofReceipt{
		Proof: &dvp.Proof{
			A: [2]*big.Int{big.NewInt(1), big.NewInt(2)},
			B: [2][2]*big.Int{
				{big.NewInt(3), big.NewInt(4)},
				{big.NewInt(5), big.NewInt(6)},
			},
			C: [2]*big.Int{big.NewInt(7), big.NewInt(8)},
		},
		Message:          big.NewInt(12345),
		TreeNumbers:      []*big.Int{big.NewInt(0)},
		MerkleRoots:      []*big.Int{big.NewInt(100)},
		Nullifiers:       []*big.Int{big.NewInt(200)},
		Commitments:      []*big.Int{big.NewInt(300)},
		RevertCommitment: big.NewInt(0),
	}
}

// stubDvpEncryptor satisfies the dvpEncryptor interface in contractclient/dvp.go.
type stubDvpEncryptor struct {
	out       []byte
	err       error
	spySalt   *big.Int
	spyMsg    *types.DvpSwapMessage
	callCount int
}

func (s *stubDvpEncryptor) EncryptDvpSwapMessage(_ context.Context, salt *big.Int, msg *types.DvpSwapMessage) ([]byte, error) {
	s.spySalt = salt
	s.spyMsg = msg
	s.callCount++
	if s.err != nil {
		return nil, s.err
	}
	if s.out == nil {
		return []byte("encrypted-blob"), nil
	}
	return s.out, nil
}

// validSharedID returns a hex string of exactly 32 bytes (64 hex chars) so
// conv.StringToBytes32 succeeds in tests.
func validSharedID() string {
	return "abcd1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab"
}

func swapMsg() *types.DvpSwapMessage {
	return &types.DvpSwapMessage{
		SharedId:       validSharedID(),
		TokenInType:    types.DvpEnygma,
		TokenInAddress: "0x1122334455667788990011223344556677889900",
	}
}

func TestDvpClient_DepositERC721(t *testing.T) {
	t.Run("successfully deposits ERC721 via operator executor", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		operatorExecutor := &stubExecutor{}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, &stubDvpEncryptor{})

		wantContractAddress := common.HexToAddress("0x1111111111111111111111111111111111111111")
		err := client.DepositERC721(context.Background(), "test-event-id", wantContractAddress, big.NewInt(42), big.NewInt(12345), big.NewInt(0xCAFE), []byte("encrypted"))

		require.NoError(t, err)
		assert.Equal(t, address, operatorExecutor.spyExecuteAddress)
		assert.NotEmpty(t, operatorExecutor.spyExecuteCalldata)
		// Non-operator executor should not be invoked.
		assert.Empty(t, executor.spyExecuteCalldata)
	})

	t.Run("wraps operator executor errors in DvpClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		operatorExecutor := &stubExecutor{executeErr: wantError}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, &stubDvpEncryptor{})

		err := client.DepositERC721(context.Background(), "test-event-id", common.Address{}, big.NewInt(42), big.NewInt(12345), big.NewInt(0xCAFE), []byte("encrypted"))

		require.Error(t, err)
		var clientErr *contractclient.DvpClientError
		require.ErrorAs(t, err, &clientErr, "error should be wrapped in DvpClientError")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})
}

func TestDvpClient_WithdrawERC721(t *testing.T) {
	t.Run("successfully withdraws ERC721 via operator executor", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		operatorExecutor := &stubExecutor{}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, &stubDvpEncryptor{})

		proofReceipt := createDvpTestProofReceipt()
		err := client.WithdrawERC721(context.Background(), "test-event-id", common.Address{}, big.NewInt(42), common.Address{}, big.NewInt(0xCAFE), proofReceipt, []byte("encrypted"))

		require.NoError(t, err)
		assert.Equal(t, address, operatorExecutor.spyExecuteAddress)
		assert.NotEmpty(t, operatorExecutor.spyExecuteCalldata)
	})

	t.Run("wraps operator executor errors in DvpClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		operatorExecutor := &stubExecutor{executeErr: wantError}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, &stubDvpEncryptor{})

		proofReceipt := createDvpTestProofReceipt()
		err := client.WithdrawERC721(context.Background(), "test-event-id", common.Address{}, big.NewInt(42), common.Address{}, big.NewInt(0xCAFE), proofReceipt, []byte("encrypted"))

		require.Error(t, err)
		var clientErr *contractclient.DvpClientError
		require.ErrorAs(t, err, &clientErr)
		require.ErrorIs(t, err, wantError)
	})
}

func TestDvpClient_DepositERC1155(t *testing.T) {
	t.Run("successfully deposits ERC1155 via operator executor", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		operatorExecutor := &stubExecutor{}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, &stubDvpEncryptor{})

		err := client.DepositERC1155(context.Background(), "test-event-id", common.Address{}, big.NewInt(1), big.NewInt(100), []byte("data"), big.NewInt(12345), big.NewInt(0xCAFE), []byte("encrypted"))

		require.NoError(t, err)
		assert.Equal(t, address, operatorExecutor.spyExecuteAddress)
		assert.NotEmpty(t, operatorExecutor.spyExecuteCalldata)
	})

	t.Run("wraps operator executor errors in DvpClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		operatorExecutor := &stubExecutor{executeErr: wantError}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, &stubDvpEncryptor{})

		err := client.DepositERC1155(context.Background(), "test-event-id", common.Address{}, big.NewInt(1), big.NewInt(100), []byte("data"), big.NewInt(12345), big.NewInt(0xCAFE), []byte("encrypted"))

		require.Error(t, err)
		var clientErr *contractclient.DvpClientError
		require.ErrorAs(t, err, &clientErr)
		require.ErrorIs(t, err, wantError)
	})
}

func TestDvpClient_WithdrawERC1155(t *testing.T) {
	t.Run("successfully withdraws ERC1155 via operator executor", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		operatorExecutor := &stubExecutor{}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, &stubDvpEncryptor{})

		proofReceipt := createDvpTestProofReceipt()
		err := client.WithdrawERC1155(context.Background(), "test-event-id", common.Address{}, big.NewInt(1), big.NewInt(100), common.Address{}, big.NewInt(0xCAFE), proofReceipt, []byte("encrypted"))

		require.NoError(t, err)
		assert.Equal(t, address, operatorExecutor.spyExecuteAddress)
		assert.NotEmpty(t, operatorExecutor.spyExecuteCalldata)
	})

	t.Run("wraps operator executor errors in DvpClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		operatorExecutor := &stubExecutor{executeErr: wantError}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, &stubDvpEncryptor{})

		proofReceipt := createDvpTestProofReceipt()
		err := client.WithdrawERC1155(context.Background(), "test-event-id", common.Address{}, big.NewInt(1), big.NewInt(100), common.Address{}, big.NewInt(0xCAFE), proofReceipt, []byte("encrypted"))

		require.Error(t, err)
		var clientErr *contractclient.DvpClientError
		require.ErrorAs(t, err, &clientErr)
		require.ErrorIs(t, err, wantError)
	})
}

func TestDvpClient_MixFundsERC1155(t *testing.T) {
	t.Run("successfully mixes funds via executor", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		operatorExecutor := &stubExecutor{}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, &stubDvpEncryptor{})

		proofReceipt := createDvpTestProofReceipt()
		err := client.MixFundsERC1155(context.Background(), "test-event-id", common.Address{}, proofReceipt)

		require.NoError(t, err)
		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotEmpty(t, executor.spyExecuteCalldata)
	})

	t.Run("wraps executor errors in DvpClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{executeErr: wantError}
		operatorExecutor := &stubExecutor{}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, &stubDvpEncryptor{})

		proofReceipt := createDvpTestProofReceipt()
		err := client.MixFundsERC1155(context.Background(), "test-event-id", common.Address{}, proofReceipt)

		require.Error(t, err)
		var clientErr *contractclient.DvpClientError
		require.ErrorAs(t, err, &clientErr)
		require.ErrorIs(t, err, wantError)
	})
}

func TestDvpClient_InitiateSwap(t *testing.T) {
	t.Run("happy path: invokes encryptor with salt+msg, contract receives encrypted blob + ctxt + validityTime", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		operatorExecutor := &stubExecutor{}
		encryptor := &stubDvpEncryptor{out: []byte("ENC")}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, encryptor)

		salt := big.NewInt(0xBBBB)
		ctxt := []byte("ML-KEM-CTXT")
		msg := swapMsg()
		proof := createDvpTestProofReceipt()

		err := client.InitiateSwap(context.Background(), salt, ctxt, msg, proof, 3600, big.NewInt(0xDEAD))
		require.NoError(t, err)

		// Encryption flow
		assert.Equal(t, 1, encryptor.callCount)
		assert.Equal(t, salt, encryptor.spySalt, "encryptor must receive the destination salt")
		assert.Equal(t, msg, encryptor.spyMsg)

		// Execute through standard executor (not operator).
		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotEmpty(t, executor.spyExecuteCalldata)
	})

	t.Run("propagates encryptor error and does NOT call the executor", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		operatorExecutor := &stubExecutor{}
		encryptor := &stubDvpEncryptor{err: errors.New("cts down")}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, encryptor)

		err := client.InitiateSwap(context.Background(), big.NewInt(1), nil, swapMsg(), createDvpTestProofReceipt(), 3600, nil)
		require.Error(t, err)

		var clientErr *contractclient.DvpClientError
		require.ErrorAs(t, err, &clientErr)
		assert.Empty(t, executor.spyExecuteCalldata, "executor must not be invoked when encryption fails")
	})

	t.Run("rejects malformed sharedId before encryption or executor call", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		operatorExecutor := &stubExecutor{}
		encryptor := &stubDvpEncryptor{}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, encryptor)

		msg := swapMsg()
		msg.SharedId = "not-a-32-byte-hex-string"

		err := client.InitiateSwap(context.Background(), big.NewInt(1), nil, msg, createDvpTestProofReceipt(), 3600, nil)
		require.Error(t, err)

		var clientErr *contractclient.DvpClientError
		require.ErrorAs(t, err, &clientErr)
		assert.Empty(t, executor.spyExecuteCalldata)
	})

	t.Run("wraps executor error in DvpClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{executeErr: wantError}
		operatorExecutor := &stubExecutor{}
		encryptor := &stubDvpEncryptor{out: []byte("ENC")}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, encryptor)

		err := client.InitiateSwap(context.Background(), big.NewInt(1), nil, swapMsg(), createDvpTestProofReceipt(), 3600, big.NewInt(0))
		require.Error(t, err)

		var clientErr *contractclient.DvpClientError
		require.ErrorAs(t, err, &clientErr)
		require.ErrorIs(t, err, wantError)
		assert.NotErrorIs(t, err, dvp.ErrSwapAlreadyInitiated)
	})

	t.Run("returns ErrSwapAlreadyInitiated when revert data matches DvpSwapAlreadyExists selector", func(t *testing.T) {
		selector := Dvp.DvpDvpSwapAlreadyExistsErrorID().Bytes()[:4]
		executor := &stubExecutor{executeErr: contractclient.NewErrorWithRevertData(selector)}
		operatorExecutor := &stubExecutor{}
		encryptor := &stubDvpEncryptor{out: []byte("ENC")}

		client := contractclient.NewDvpClient(common.HexToAddress("0x1"), executor, operatorExecutor, encryptor)

		err := client.InitiateSwap(context.Background(), big.NewInt(1), nil, swapMsg(), createDvpTestProofReceipt(), 3600, big.NewInt(0))
		require.Error(t, err)
		require.ErrorIs(t, err, dvp.ErrSwapAlreadyInitiated)
	})

	t.Run("does not return ErrSwapAlreadyInitiated for a different custom error selector", func(t *testing.T) {
		// Use a different no-arg Dvp error to prove selector compare isn't a false positive.
		otherSelector := Dvp.DvpDvpSwapNotExpiredErrorID().Bytes()[:4]
		executor := &stubExecutor{executeErr: contractclient.NewErrorWithRevertData(otherSelector)}
		operatorExecutor := &stubExecutor{}
		encryptor := &stubDvpEncryptor{out: []byte("ENC")}

		client := contractclient.NewDvpClient(common.HexToAddress("0x1"), executor, operatorExecutor, encryptor)

		err := client.InitiateSwap(context.Background(), big.NewInt(1), nil, swapMsg(), createDvpTestProofReceipt(), 3600, big.NewInt(0))
		require.Error(t, err)
		assert.NotErrorIs(t, err, dvp.ErrSwapAlreadyInitiated)

		var clientErr *contractclient.DvpClientError
		require.ErrorAs(t, err, &clientErr)
	})

	t.Run("does not panic on truncated revert data shorter than 4 bytes", func(t *testing.T) {
		executor := &stubExecutor{executeErr: contractclient.NewErrorWithRevertData([]byte{0x01, 0x02})}
		operatorExecutor := &stubExecutor{}
		encryptor := &stubDvpEncryptor{out: []byte("ENC")}

		client := contractclient.NewDvpClient(common.HexToAddress("0x1"), executor, operatorExecutor, encryptor)

		err := client.InitiateSwap(context.Background(), big.NewInt(1), nil, swapMsg(), createDvpTestProofReceipt(), 3600, big.NewInt(0))
		require.Error(t, err)
		assert.NotErrorIs(t, err, dvp.ErrSwapAlreadyInitiated)

		var clientErr *contractclient.DvpClientError
		require.ErrorAs(t, err, &clientErr)
	})
}

func TestDvpClient_CompleteSwap(t *testing.T) {
	t.Run("happy path: encryptor invoked with destSalt, executor receives calldata", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		operatorExecutor := &stubExecutor{}
		encryptor := &stubDvpEncryptor{out: []byte("COMPLETE-ENC")}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, encryptor)

		destSalt := big.NewInt(0xBBBB)
		msg := swapMsg()
		proof := createDvpTestProofReceipt()

		err := client.CompleteSwap(context.Background(), destSalt, msg, proof)
		require.NoError(t, err)

		assert.Equal(t, destSalt, encryptor.spySalt,
			"completion encryption must reuse the original destination salt — same key Alice used")
		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotEmpty(t, executor.spyExecuteCalldata)
	})

	t.Run("propagates encryptor error", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		operatorExecutor := &stubExecutor{}
		encryptor := &stubDvpEncryptor{err: errors.New("cts down")}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, encryptor)

		err := client.CompleteSwap(context.Background(), big.NewInt(1), swapMsg(), createDvpTestProofReceipt())
		require.Error(t, err)

		var clientErr *contractclient.DvpClientError
		require.ErrorAs(t, err, &clientErr)
		assert.Empty(t, executor.spyExecuteCalldata)
	})
}

func TestDvpClient_CancelSwap(t *testing.T) {
	t.Run("calls executor with packed CancelSwap calldata", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		operatorExecutor := &stubExecutor{}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, &stubDvpEncryptor{})

		err := client.CancelSwap(context.Background(), validSharedID(), big.NewInt(0))
		require.NoError(t, err)

		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotEmpty(t, executor.spyExecuteCalldata)
	})

	t.Run("rejects malformed sharedId", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		operatorExecutor := &stubExecutor{}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, &stubDvpEncryptor{})

		err := client.CancelSwap(context.Background(), "not-a-32-byte-hex-string", nil)
		require.Error(t, err)

		var clientErr *contractclient.DvpClientError
		require.ErrorAs(t, err, &clientErr)
		assert.Empty(t, executor.spyExecuteCalldata)
	})
}

func TestDvpClient_ExpireSwap(t *testing.T) {
	t.Run("calls executor with packed ExpireSwap calldata", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		operatorExecutor := &stubExecutor{}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, &stubDvpEncryptor{})

		err := client.ExpireSwap(context.Background(), validSharedID())
		require.NoError(t, err)

		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotEmpty(t, executor.spyExecuteCalldata)
	})

	t.Run("rejects malformed sharedId", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		operatorExecutor := &stubExecutor{}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, &stubDvpEncryptor{})

		err := client.ExpireSwap(context.Background(), "not-a-32-byte-hex-string")
		require.Error(t, err)

		var clientErr *contractclient.DvpClientError
		require.ErrorAs(t, err, &clientErr)
		assert.Empty(t, executor.spyExecuteCalldata)
	})
}

func TestDvpClient_IsSwapExpired(t *testing.T) {
	t.Run("propagates Call error wrapped in DvpClientError", func(t *testing.T) {
		wantError := errors.New("rpc down")
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{callErr: wantError}
		operatorExecutor := &stubExecutor{}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, &stubDvpEncryptor{})

		_, err := client.IsSwapExpired(context.Background(), validSharedID())
		require.Error(t, err)

		var clientErr *contractclient.DvpClientError
		require.ErrorAs(t, err, &clientErr)
		require.ErrorIs(t, err, wantError)
	})

	t.Run("rejects malformed sharedId", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		operatorExecutor := &stubExecutor{}

		client := contractclient.NewDvpClient(address, executor, operatorExecutor, &stubDvpEncryptor{})

		_, err := client.IsSwapExpired(context.Background(), "not-a-32-byte-hex-string")
		require.Error(t, err)

		var clientErr *contractclient.DvpClientError
		require.ErrorAs(t, err, &clientErr)
		assert.Empty(t, executor.spyCallCalldata)
	})
}
