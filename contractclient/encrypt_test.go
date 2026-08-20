package contractclient_test

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	ps "github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/ParticipantStorageV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/encrypt"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type StubEncryptEthereumClient struct {
	block *ethTypes.Block
	err   error

	spyBlockNumber *big.Int
}

func (s *StubEncryptEthereumClient) BlockByNumber(ctx context.Context, blockNum *big.Int) (*ethTypes.Block, error) {
	s.spyBlockNumber = blockNum
	return s.block, s.err
}

type StubParticipantStorageClient struct {
	chainInfo  ps.ParticipantStructsPrivacyNodeViewData
	chainInfos []ps.ParticipantStructsPrivacyNodeViewData
	err        error

	spyBlockNumber *big.Int
}

func (c *StubParticipantStorageClient) GetVenOperatorChainInfo(
	_ context.Context,
	blockNumber *big.Int,
) (ps.ParticipantStructsPrivacyNodeViewData, error) {
	return c.chainInfo, c.err
}

func (c *StubParticipantStorageClient) GetMyChainInfo(
	_ context.Context,
	blockNumber *big.Int,
) (ps.ParticipantStructsPrivacyNodeViewData, error) {
	c.spyBlockNumber = blockNumber
	return c.chainInfo, c.err
}

func (c *StubParticipantStorageClient) GetChainViewData(
	_ context.Context,
	chainID *big.Int,
	blockNumber *big.Int,
) (ps.ParticipantStructsPrivacyNodeViewData, error) {
	c.spyBlockNumber = blockNumber
	return c.chainInfo, c.err
}

func (c *StubParticipantStorageClient) GetChainViewDataBatch(
	_ context.Context,
	chainIDs []*big.Int,
) ([]ps.ParticipantStructsPrivacyNodeViewData, error) {
	return c.chainInfos, c.err
}

// StubKOSClient implements the ctsClient gRPC interface
// (Encrypt + EncryptWithoutFP + EncryptWithoutFPWithSS).
type StubKOSClient struct {
	lastEncryptReq                *encrypt.EncryptRequest
	lastEncryptWithoutFPReq       *encrypt.EncryptWithoutFPRequest
	lastEncryptWithoutFPWithSSReq *encrypt.EncryptWithoutFPWithSSRequest

	encryptResp                *encrypt.EncryptResponse
	encryptWithoutFPResp       *encrypt.EncryptWithoutFPResponse
	encryptWithoutFPWithSSResp *encrypt.EncryptWithoutFPWithSSResponse
	err                        error
}

func (c *StubKOSClient) Encrypt(
	ctx context.Context,
	req *encrypt.EncryptRequest,
	opts ...grpc.CallOption,
) (*encrypt.EncryptResponse, error) {
	c.lastEncryptReq = req
	return c.encryptResp, c.err
}

func (c *StubKOSClient) EncryptWithoutFP(
	ctx context.Context,
	in *encrypt.EncryptWithoutFPRequest,
	opts ...grpc.CallOption,
) (*encrypt.EncryptWithoutFPResponse, error) {
	c.lastEncryptWithoutFPReq = in
	return c.encryptWithoutFPResp, c.err
}

func (c *StubKOSClient) EncryptWithoutFPWithSS(
	ctx context.Context,
	in *encrypt.EncryptWithoutFPWithSSRequest,
	opts ...grpc.CallOption,
) (*encrypt.EncryptWithoutFPWithSSResponse, error) {
	c.lastEncryptWithoutFPWithSSReq = in
	return c.encryptWithoutFPWithSSResp, c.err
}

func TestAdditionalDataEncryptor(t *testing.T) {
	t.Run("encrypts additional data through the KOS", func(t *testing.T) {
		wantEncryptedData := []byte("encrypted data bytes")

		wantDestChainInfo := ps.ParticipantStructsPrivacyNodeViewData{
			ChainId: new(big.Int).SetUint64(1337),
		}
		wantBlockNumber := new(big.Int).SetUint64(100)

		var wantLatestBlockNumber *big.Int = nil

		header := &ethTypes.Header{
			Number:     wantBlockNumber,
			ParentHash: common.Hash([32]byte{0xDE, 0xAD, 0xBE, 0xEF}),
		}
		block := ethTypes.NewBlockWithHeader(header)

		ethClient := &StubEncryptEthereumClient{
			block: block,
		}
		psClient := &StubParticipantStorageClient{
			chainInfo: wantDestChainInfo,
		}
		kos := &StubKOSClient{
			encryptWithoutFPResp: &encrypt.EncryptWithoutFPResponse{
				EncryptedData: wantEncryptedData,
			},
		}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		gotEncrData, err := encr.EncryptAdditionalData(context.Background(), []types.AtomicTeleportAdditionalData{
			{SharedId: "example-shared-id"},
		})
		require.Nil(t, err)

		assert.Equal(t, wantLatestBlockNumber, ethClient.spyBlockNumber)
		assert.Equal(t, wantDestChainInfo.ChainId.Uint64(), kos.lastEncryptWithoutFPReq.ChainID)
		assert.Equal(t, wantBlockNumber.Uint64(), kos.lastEncryptWithoutFPReq.BlockNumber)
		assert.Equal(t, common.Bytes2Hex(wantEncryptedData), gotEncrData)
	})

	t.Run("wraps errors from ethereum client in EncryptorError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.EncryptorError{}

		header := &ethTypes.Header{Number: new(big.Int).SetUint64(100)}
		block := ethTypes.NewBlockWithHeader(header)

		ethClient := &StubEncryptEthereumClient{block: block, err: wantError}
		psClient := &StubParticipantStorageClient{}
		kos := &StubKOSClient{}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		_, gotErr := encr.EncryptAdditionalData(context.Background(), []types.AtomicTeleportAdditionalData{{SharedId: "x"}})

		require.ErrorAs(t, gotErr, &wantErrorType)
		require.ErrorIs(t, gotErr, wantError)
	})

	t.Run("wraps errors from participant storage client in EncryptorError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.EncryptorError{}

		header := &ethTypes.Header{Number: new(big.Int).SetUint64(100)}
		block := ethTypes.NewBlockWithHeader(header)

		ethClient := &StubEncryptEthereumClient{block: block}
		psClient := &StubParticipantStorageClient{err: wantError}
		kos := &StubKOSClient{}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		_, gotErr := encr.EncryptAdditionalData(context.Background(), []types.AtomicTeleportAdditionalData{{SharedId: "x"}})

		require.ErrorAs(t, gotErr, &wantErrorType)
		require.ErrorIs(t, gotErr, wantError)
	})

	t.Run("wraps errors from kos client in EncryptorError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.EncryptorError{}

		header := &ethTypes.Header{Number: new(big.Int).SetUint64(100)}
		block := ethTypes.NewBlockWithHeader(header)

		ethClient := &StubEncryptEthereumClient{block: block}
		psClient := &StubParticipantStorageClient{
			chainInfo: ps.ParticipantStructsPrivacyNodeViewData{ChainId: big.NewInt(1)},
		}
		kos := &StubKOSClient{err: wantError}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		_, gotErr := encr.EncryptAdditionalData(context.Background(), []types.AtomicTeleportAdditionalData{{SharedId: "x"}})

		require.ErrorAs(t, gotErr, &wantErrorType)
		require.ErrorIs(t, gotErr, wantError)
	})
}

func TestEncryptor_EncryptMessage(t *testing.T) {
	t.Run("encrypts message through KOS", func(t *testing.T) {
		wantBlockNumber := new(big.Int).SetUint64(1337)
		wantFingerprint := "example-message-tag"
		wantEncryptedData := []byte("example-data")
		wantPrevBlockHash := common.Hash([32]byte{0xDE, 0xAD, 0xBE, 0xEF})

		wantDestChainInfo := ps.ParticipantStructsPrivacyNodeViewData{
			ChainId: new(big.Int).SetUint64(1337),
		}

		header := &ethTypes.Header{
			Number:     wantBlockNumber,
			ParentHash: wantPrevBlockHash,
		}
		block := ethTypes.NewBlockWithHeader(header)

		ethClient := &StubEncryptEthereumClient{block: block}
		psClient := &StubParticipantStorageClient{chainInfo: wantDestChainInfo}
		kos := &StubKOSClient{
			encryptResp: &encrypt.EncryptResponse{
				Fingerprint:   wantFingerprint,
				EncryptedData: wantEncryptedData,
			},
		}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		got, err := encr.EncryptMessages(context.Background(),
			[]types.DispatchedMessageToPrivateHub{{SharedId: "example-shared-id"}},
			wantDestChainInfo.ChainId,
		)
		require.Nil(t, err)

		assert.Equal(t, wantBlockNumber, psClient.spyBlockNumber)
		assert.Equal(t, wantBlockNumber.Uint64(), kos.lastEncryptReq.BlockNumber)
		assert.Equal(t, wantPrevBlockHash.String(), kos.lastEncryptReq.PrevBlockHash)
		assert.Equal(t, wantDestChainInfo.ChainId.Uint64(), kos.lastEncryptReq.ChainID)

		assert.Equal(t, wantBlockNumber, got.BlockNumber)
		assert.Equal(t, wantFingerprint, got.MessageTag)
		assert.Equal(t, wantEncryptedData, got.Data)
	})

	t.Run("wraps errors from ethereum client in EncryptorError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.EncryptorError{}

		header := &ethTypes.Header{Number: new(big.Int).SetUint64(100)}
		block := ethTypes.NewBlockWithHeader(header)

		ethClient := &StubEncryptEthereumClient{block: block, err: wantError}
		psClient := &StubParticipantStorageClient{}
		kos := &StubKOSClient{}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		_, gotErr := encr.EncryptMessages(context.Background(),
			[]types.DispatchedMessageToPrivateHub{{SharedId: "x"}},
			new(big.Int).SetUint64(1337),
		)

		require.ErrorAs(t, gotErr, &wantErrorType)
		require.ErrorIs(t, gotErr, wantError)
	})

	t.Run("wraps errors from participant storage client in EncryptorError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.EncryptorError{}

		header := &ethTypes.Header{Number: new(big.Int).SetUint64(100)}
		block := ethTypes.NewBlockWithHeader(header)

		ethClient := &StubEncryptEthereumClient{block: block}
		psClient := &StubParticipantStorageClient{err: wantError}
		kos := &StubKOSClient{}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		_, gotErr := encr.EncryptMessages(context.Background(),
			[]types.DispatchedMessageToPrivateHub{{SharedId: "x"}},
			new(big.Int).SetUint64(1337),
		)

		require.ErrorAs(t, gotErr, &wantErrorType)
		require.ErrorIs(t, gotErr, wantError)
	})

	t.Run("wraps errors from kos client in EncryptorError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.EncryptorError{}

		header := &ethTypes.Header{Number: new(big.Int).SetUint64(100)}
		block := ethTypes.NewBlockWithHeader(header)

		ethClient := &StubEncryptEthereumClient{block: block}
		psClient := &StubParticipantStorageClient{
			chainInfo: ps.ParticipantStructsPrivacyNodeViewData{
				ChainId: new(big.Int).SetUint64(1),
			},
		}
		kos := &StubKOSClient{err: wantError}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		_, gotErr := encr.EncryptMessages(context.Background(),
			[]types.DispatchedMessageToPrivateHub{{SharedId: "x"}},
			new(big.Int).SetUint64(1337),
		)

		require.ErrorAs(t, gotErr, &wantErrorType)
		require.ErrorIs(t, gotErr, wantError)
	})
}

func TestEncryptor_EncryptEnygmaTransferBatches(t *testing.T) {
	t.Run("encrypts multiple batches for different chains", func(t *testing.T) {
		wantEncryptedData := []byte("encrypted-batch-data")
		blockNumber := new(big.Int).SetInt64(100)

		batches := []*types.EnygmaTransferBatch{
			{ToChainID: new(big.Int).SetInt64(200), BatchId: "batch-1"},
			{ToChainID: new(big.Int).SetInt64(300), BatchId: "batch-2"},
		}

		chainInfos := []ps.ParticipantStructsPrivacyNodeViewData{
			{ChainId: new(big.Int).SetInt64(200), RaylsViewPublicKey: "pubkey-200"},
			{ChainId: new(big.Int).SetInt64(300), RaylsViewPublicKey: "pubkey-300"},
		}

		ethClient := &StubEncryptEthereumClient{}
		psClient := &StubParticipantStorageClient{chainInfos: chainInfos}
		kos := &StubKOSClient{
			encryptWithoutFPResp: &encrypt.EncryptWithoutFPResponse{
				EncryptedData: wantEncryptedData,
			},
		}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		gotEncryptedBatches, err := encr.EncryptEnygmaTransferBatches(context.Background(), batches, blockNumber)
		require.Nil(t, err)

		require.Len(t, gotEncryptedBatches, 2)
		assert.Equal(t, wantEncryptedData, gotEncryptedBatches[0])
		assert.Equal(t, wantEncryptedData, gotEncryptedBatches[1])
		assert.Equal(t, blockNumber.Uint64(), kos.lastEncryptWithoutFPReq.BlockNumber)
	})

	t.Run("wraps errors from participant storage client in EncryptorError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.EncryptorError{}

		batches := []*types.EnygmaTransferBatch{
			{ToChainID: new(big.Int).SetInt64(200), BatchId: "batch-1"},
		}
		blockNumber := new(big.Int).SetInt64(100)

		ethClient := &StubEncryptEthereumClient{}
		psClient := &StubParticipantStorageClient{err: wantError}
		kos := &StubKOSClient{}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		_, gotErr := encr.EncryptEnygmaTransferBatches(context.Background(), batches, blockNumber)

		require.ErrorAs(t, gotErr, &wantErrorType)
		require.ErrorIs(t, gotErr, wantError)
	})

	t.Run("wraps errors from kos client in EncryptorError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.EncryptorError{}

		batches := []*types.EnygmaTransferBatch{
			{ToChainID: new(big.Int).SetInt64(200), BatchId: "batch-1"},
		}
		blockNumber := new(big.Int).SetInt64(100)

		chainInfos := []ps.ParticipantStructsPrivacyNodeViewData{
			{ChainId: new(big.Int).SetInt64(200), RaylsViewPublicKey: "pubkey-200"},
		}

		ethClient := &StubEncryptEthereumClient{}
		psClient := &StubParticipantStorageClient{chainInfos: chainInfos}
		kos := &StubKOSClient{err: wantError}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		_, gotErr := encr.EncryptEnygmaTransferBatches(context.Background(), batches, blockNumber)

		require.ErrorAs(t, gotErr, &wantErrorType)
		require.ErrorIs(t, gotErr, wantError)
	})
}

func TestEncryptor_EncryptEnygmaTransferBatchCompleted(t *testing.T) {
	t.Run("encrypts completed transfer batch messages", func(t *testing.T) {
		wantEncryptedData := []byte("encrypted-completed-data")
		wantBlockNumber := new(big.Int).SetUint64(100)

		messages := []types.EnygmaTransferCompleted{
			{MessageId: "msg-1", TransactionHash: "tx-hash-1"},
			{MessageId: "msg-2", TransactionHash: "tx-hash-2"},
		}

		chainInfo := ps.ParticipantStructsPrivacyNodeViewData{
			ChainId:            new(big.Int).SetInt64(42),
			RaylsViewPublicKey: "my-pubkey",
		}

		header := &ethTypes.Header{Number: wantBlockNumber}
		block := ethTypes.NewBlockWithHeader(header)

		ethClient := &StubEncryptEthereumClient{block: block}
		psClient := &StubParticipantStorageClient{chainInfo: chainInfo}
		kos := &StubKOSClient{
			encryptWithoutFPResp: &encrypt.EncryptWithoutFPResponse{
				EncryptedData: wantEncryptedData,
			},
		}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		gotEncryptedData, err := encr.EncryptEnygmaTransferBatchCompleted(context.Background(), messages)
		require.Nil(t, err)

		assert.Equal(t, wantEncryptedData, gotEncryptedData)
		assert.Equal(t, wantBlockNumber, psClient.spyBlockNumber)
		assert.Equal(t, wantBlockNumber.Uint64(), kos.lastEncryptWithoutFPReq.BlockNumber)
		assert.Equal(t, chainInfo.ChainId.Uint64(), kos.lastEncryptWithoutFPReq.ChainID)
	})

	t.Run("wraps errors from ethereum client in EncryptorError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.EncryptorError{}

		ethClient := &StubEncryptEthereumClient{err: wantError}
		psClient := &StubParticipantStorageClient{}
		kos := &StubKOSClient{}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		_, gotErr := encr.EncryptEnygmaTransferBatchCompleted(context.Background(), []types.EnygmaTransferCompleted{{MessageId: "msg-1"}})

		require.ErrorAs(t, gotErr, &wantErrorType)
		require.ErrorIs(t, gotErr, wantError)
	})

	t.Run("wraps errors from participant storage client in EncryptorError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.EncryptorError{}

		header := &ethTypes.Header{Number: new(big.Int).SetUint64(100)}
		block := ethTypes.NewBlockWithHeader(header)

		ethClient := &StubEncryptEthereumClient{block: block}
		psClient := &StubParticipantStorageClient{err: wantError}
		kos := &StubKOSClient{}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		_, gotErr := encr.EncryptEnygmaTransferBatchCompleted(context.Background(), []types.EnygmaTransferCompleted{{MessageId: "msg-1"}})

		require.ErrorAs(t, gotErr, &wantErrorType)
		require.ErrorIs(t, gotErr, wantError)
	})

	t.Run("wraps errors from kos client in EncryptorError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.EncryptorError{}

		header := &ethTypes.Header{Number: new(big.Int).SetUint64(100)}
		block := ethTypes.NewBlockWithHeader(header)

		chainInfo := ps.ParticipantStructsPrivacyNodeViewData{ChainId: new(big.Int).SetInt64(42)}

		ethClient := &StubEncryptEthereumClient{block: block}
		psClient := &StubParticipantStorageClient{chainInfo: chainInfo}
		kos := &StubKOSClient{err: wantError}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		_, gotErr := encr.EncryptEnygmaTransferBatchCompleted(context.Background(), []types.EnygmaTransferCompleted{{MessageId: "msg-1"}})

		require.ErrorAs(t, gotErr, &wantErrorType)
		require.ErrorIs(t, gotErr, wantError)
	})
}

func TestEncryptor_EncryptDvpSwapMessage(t *testing.T) {
	t.Run("encrypts dvp swap message through KOS", func(t *testing.T) {
		wantEncryptedData := []byte("encrypted-swap-data")
		wantSalt := new(big.Int).SetInt64(999)
		wantMessage := types.DvpSwapMessage{
			SharedId: "msg-123",
			To:       "0x123",
		}

		kos := &StubKOSClient{
			encryptWithoutFPWithSSResp: &encrypt.EncryptWithoutFPWithSSResponse{
				EncryptedData: wantEncryptedData,
			},
		}
		ethClient := &StubEncryptEthereumClient{}
		psClient := &StubParticipantStorageClient{}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		gotData, err := encr.EncryptDvpSwapMessage(context.Background(), wantSalt, &wantMessage)
		require.Nil(t, err)

		assert.Equal(t, wantSalt.Bytes(), kos.lastEncryptWithoutFPWithSSReq.Ss, "didn't provide correct salt to KOS")
		assert.Equal(t, wantEncryptedData, gotData, "didn't return encrypted data from KOS")
	})

	t.Run("wraps errors from kos client in EncryptorError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.EncryptorError{}

		message := types.DvpSwapMessage{
			SharedId: "msg-123",
		}
		salt := new(big.Int).SetInt64(999)

		kos := &StubKOSClient{
			err: wantError,
		}
		ethClient := &StubEncryptEthereumClient{}
		psClient := &StubParticipantStorageClient{}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		_, gotErr := encr.EncryptDvpSwapMessage(context.Background(), salt, &message)

		require.ErrorAs(t, gotErr, &wantErrorType)
		require.ErrorIs(t, gotErr, wantError)
	})
}

func TestEncryptor_EncryptDvpExtraData(t *testing.T) {
	t.Run("encrypts private fields and leaves public fields unchanged", func(t *testing.T) {
		wantBlockNumber := new(big.Int).SetInt64(100)
		wantEncryptedValue := []byte("encrypted-secret-value")

		rawExtraData := []byte(`[
			{"Key": "name", "Value": "public-name", "IsPublic": true},
			{"Key": "secret", "Value": "secret-value", "IsPublic": false}
		]`)

		chainInfo := ps.ParticipantStructsPrivacyNodeViewData{
			ChainId:            new(big.Int).SetInt64(42),
			RaylsViewPublicKey: "my-pubkey",
		}

		header := &ethTypes.Header{Number: wantBlockNumber}
		block := ethTypes.NewBlockWithHeader(header)

		ethClient := &StubEncryptEthereumClient{block: block}
		psClient := &StubParticipantStorageClient{chainInfo: chainInfo}
		kos := &StubKOSClient{
			encryptWithoutFPResp: &encrypt.EncryptWithoutFPResponse{
				EncryptedData: wantEncryptedValue,
			},
		}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		gotExtraDataBytes, err := encr.EncryptDvpExtraData(context.Background(), rawExtraData)
		require.Nil(t, err)

		var result []map[string]interface{}
		err = json.Unmarshal(gotExtraDataBytes, &result)
		require.Nil(t, err)

		require.Len(t, result, 2)
		assert.Equal(t, "public-name", result[0]["Value"])
		assert.Equal(t, true, result[0]["IsPublic"])
		expectedBase64 := "ZW5jcnlwdGVkLXNlY3JldC12YWx1ZQ=="
		assert.Equal(t, expectedBase64, result[1]["Value"])
		assert.Equal(t, false, result[1]["IsPublic"])
		assert.Equal(t, chainInfo.ChainId.Uint64(), kos.lastEncryptWithoutFPReq.ChainID)
		assert.Equal(t, wantBlockNumber.Uint64(), kos.lastEncryptWithoutFPReq.BlockNumber)
	})

	t.Run("wraps JSON parsing errors in EncryptorError", func(t *testing.T) {
		header := &ethTypes.Header{Number: new(big.Int).SetInt64(100)}
		block := ethTypes.NewBlockWithHeader(header)

		ethClient := &StubEncryptEthereumClient{block: block}
		psClient := &StubParticipantStorageClient{}
		kos := &StubKOSClient{}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		_, gotErr := encr.EncryptDvpExtraData(context.Background(), []byte("not valid json"))

		var wantErrorType *contractclient.EncryptorError
		require.ErrorAs(t, gotErr, &wantErrorType)
	})

	t.Run("wraps errors from ethereum client in EncryptorError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.EncryptorError{}

		ethClient := &StubEncryptEthereumClient{err: wantError}
		psClient := &StubParticipantStorageClient{}
		kos := &StubKOSClient{}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		_, gotErr := encr.EncryptDvpExtraData(context.Background(), []byte(`[{"Key": "name", "Value": "v", "IsPublic": true}]`))

		require.ErrorAs(t, gotErr, &wantErrorType)
		require.ErrorIs(t, gotErr, wantError)
	})

	t.Run("wraps errors from participant storage client in EncryptorError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.EncryptorError{}

		header := &ethTypes.Header{Number: new(big.Int).SetInt64(100)}
		block := ethTypes.NewBlockWithHeader(header)

		ethClient := &StubEncryptEthereumClient{block: block}
		psClient := &StubParticipantStorageClient{err: wantError}
		kos := &StubKOSClient{}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		_, gotErr := encr.EncryptDvpExtraData(context.Background(), []byte(`[{"Key": "name", "Value": "v", "IsPublic": true}]`))

		require.ErrorAs(t, gotErr, &wantErrorType)
		require.ErrorIs(t, gotErr, wantError)
	})

	t.Run("wraps errors from kos client in EncryptorError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &contractclient.EncryptorError{}

		header := &ethTypes.Header{Number: new(big.Int).SetInt64(100)}
		block := ethTypes.NewBlockWithHeader(header)

		chainInfo := ps.ParticipantStructsPrivacyNodeViewData{
			ChainId:            new(big.Int).SetInt64(42),
			RaylsViewPublicKey: "my-pubkey",
		}

		ethClient := &StubEncryptEthereumClient{
			block: block,
		}
		psClient := &StubParticipantStorageClient{
			chainInfo: chainInfo,
		}
		kos := &StubKOSClient{
			err: wantError,
		}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		_, gotErr := encr.EncryptDvpExtraData(context.Background(), []byte(`[{"Key": "secret", "Value": "v", "IsPublic": false}]`))

		require.ErrorAs(t, gotErr, &wantErrorType, "didn't wrap error in EncryptorError")
		require.ErrorIs(t, gotErr, wantError, "didn't wrap underlying error")
	})

	t.Run("successfully marshals result after encryption", func(t *testing.T) {
		// Create valid JSON that will unmarshal but has the right structure
		rawExtraData := []byte(`[{"Key": "name", "Value": "public-name", "IsPublic": true}]`)

		header := &ethTypes.Header{
			Number: new(big.Int).SetInt64(100),
		}
		block := ethTypes.NewBlockWithHeader(header)

		ethClient := &StubEncryptEthereumClient{
			block: block,
		}
		psClient := &StubParticipantStorageClient{}
		kos := &StubKOSClient{}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		gotExtraDataBytes, err := encr.EncryptDvpExtraData(context.Background(), rawExtraData)
		require.Nil(t, err)

		// Verify result is still valid JSON
		var result []any
		err = json.Unmarshal(gotExtraDataBytes, &result)
		require.Nil(t, err)
	})

	t.Run("encrypts multiple private fields correctly", func(t *testing.T) {
		wantBlockNumber := new(big.Int).SetInt64(100)
		wantEncryptedValue := []byte("encrypted-value")

		rawExtraData := []byte(`[
			{"Key": "public1", "Value": "public-value1", "IsPublic": true},
			{"Key": "private1", "Value": "private-value1", "IsPublic": false},
			{"Key": "public2", "Value": "public-value2", "IsPublic": true},
			{"Key": "private2", "Value": "private-value2", "IsPublic": false}
		]`)

		chainInfo := ps.ParticipantStructsPrivacyNodeViewData{
			ChainId:            new(big.Int).SetInt64(42),
			RaylsViewPublicKey: "my-pubkey",
		}

		header := &ethTypes.Header{Number: wantBlockNumber}
		block := ethTypes.NewBlockWithHeader(header)

		ethClient := &StubEncryptEthereumClient{block: block}
		psClient := &StubParticipantStorageClient{chainInfo: chainInfo}
		kos := &StubKOSClient{
			encryptWithoutFPResp: &encrypt.EncryptWithoutFPResponse{
				EncryptedData: wantEncryptedValue,
			},
		}
		encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)

		gotExtraDataBytes, err := encr.EncryptDvpExtraData(context.Background(), rawExtraData)
		require.Nil(t, err)

		var result []map[string]interface{}
		err = json.Unmarshal(gotExtraDataBytes, &result)
		require.Nil(t, err)

		require.Len(t, result, 4)
		assert.Equal(t, "public-value1", result[0]["Value"])
		assert.Equal(t, true, result[0]["IsPublic"])
		assert.Equal(t, "public-value2", result[2]["Value"])
		assert.Equal(t, true, result[2]["IsPublic"])

		expectedBase64 := "ZW5jcnlwdGVkLXZhbHVl"
		assert.Equal(t, expectedBase64, result[1]["Value"])
		assert.Equal(t, expectedBase64, result[3]["Value"])
	})
}

func TestEncryptor_EncryptDvpBalanceUpdated(t *testing.T) {
	message := types.DvpBalanceUpdated{
		ErcId:              new(big.Int).SetInt64(1),
		TokenType:          1,
		ResourceId:         "resource-test",
		SourceChainId:      new(big.Int).SetInt64(9),
		DestinationChainId: new(big.Int).SetInt64(9),
		Amount:             new(big.Int).SetInt64(500),
		UpdateType:         types.Mint,
	}

	tests := []struct {
		name          string
		setup         func() (*contractclient.Encryptor, *StubEncryptEthereumClient, *StubParticipantStorageClient, *StubKOSClient, error)
		assertSuccess func(t *testing.T, data []byte, eth *StubEncryptEthereumClient, ps *StubParticipantStorageClient, kos *StubKOSClient)
	}{
		{
			name: "encrypts successfully dvp balance update using KOS",
			setup: func() (*contractclient.Encryptor, *StubEncryptEthereumClient, *StubParticipantStorageClient, *StubKOSClient, error) {
				wantBlockNumber := new(big.Int).SetInt64(200)
				header := &ethTypes.Header{Number: wantBlockNumber}
				block := ethTypes.NewBlockWithHeader(header)

				ethClient := &StubEncryptEthereumClient{block: block}
				operatorInfo := ps.ParticipantStructsPrivacyNodeViewData{
					ChainId:            new(big.Int).SetInt64(777),
					RaylsViewPublicKey: "operator-pubkey",
				}
				psClient := &StubParticipantStorageClient{chainInfo: operatorInfo}
				kos := &StubKOSClient{
					encryptWithoutFPResp: &encrypt.EncryptWithoutFPResponse{
						EncryptedData: []byte("encrypted-dvp-balance-data"),
					},
				}

				encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)
				return encr, ethClient, psClient, kos, nil
			},
			assertSuccess: func(t *testing.T, data []byte, eth *StubEncryptEthereumClient, psc *StubParticipantStorageClient, kos *StubKOSClient) {
				assert.Equal(t, []byte("encrypted-dvp-balance-data"), data)
				assert.Nil(t, eth.spyBlockNumber, "should request latest block number")
				assert.Equal(t, eth.block.NumberU64(), kos.lastEncryptWithoutFPReq.BlockNumber)
				assert.Equal(t, psc.chainInfo.ChainId.Uint64(), kos.lastEncryptWithoutFPReq.ChainID)
			},
		},
		{
			name: "wraps ethereum client errors",
			setup: func() (*contractclient.Encryptor, *StubEncryptEthereumClient, *StubParticipantStorageClient, *StubKOSClient, error) {
				wantErr := errors.New("example error")
				ethClient := &StubEncryptEthereumClient{err: wantErr}
				psClient := &StubParticipantStorageClient{}
				kos := &StubKOSClient{}
				encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)
				return encr, ethClient, psClient, kos, wantErr
			},
		},
		{
			name: "wraps participant storage client errors",
			setup: func() (*contractclient.Encryptor, *StubEncryptEthereumClient, *StubParticipantStorageClient, *StubKOSClient, error) {
				wantErr := errors.New("example error")
				header := &ethTypes.Header{Number: new(big.Int).SetInt64(100)}
				block := ethTypes.NewBlockWithHeader(header)
				ethClient := &StubEncryptEthereumClient{block: block}
				psClient := &StubParticipantStorageClient{err: wantErr}
				kos := &StubKOSClient{}
				encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)
				return encr, ethClient, psClient, kos, wantErr
			},
		},
		{
			name: "wraps kos client errors",
			setup: func() (*contractclient.Encryptor, *StubEncryptEthereumClient, *StubParticipantStorageClient, *StubKOSClient, error) {
				wantErr := errors.New("example error")
				header := &ethTypes.Header{Number: new(big.Int).SetInt64(100)}
				block := ethTypes.NewBlockWithHeader(header)
				ethClient := &StubEncryptEthereumClient{block: block}
				operatorInfo := ps.ParticipantStructsPrivacyNodeViewData{ChainId: new(big.Int).SetInt64(55)}
				psClient := &StubParticipantStorageClient{chainInfo: operatorInfo}
				kos := &StubKOSClient{err: wantErr}
				encr := contractclient.NewEncryptor(kos, psClient, ethClient, nil)
				return encr, ethClient, psClient, kos, wantErr
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encr, ethClient, psClient, kosClient, wantErr := tt.setup()
			gotData, err := encr.EncryptDvpBalanceUpdated(context.Background(), message)

			if wantErr != nil {
				require.Error(t, err)
				var encryptErr *contractclient.EncryptorError
				require.ErrorAs(t, err, &encryptErr)
				require.ErrorIs(t, err, wantErr)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, tt.assertSuccess)
			tt.assertSuccess(t, gotData, ethClient, psClient, kosClient)
		})
	}
}
