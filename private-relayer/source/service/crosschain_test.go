package service_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/service/testdata"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCrossChainService(t *testing.T) { //nolint:gocognit // comprehensive test with detailed mock setup
	testtools.SilenceLogger()

	// tickerPeriod is set far longer than any per-subtest ctx timeout, so the
	// only iteration each subtest's Run loop performs is the initial one (via
	// the internal initialRun channel). This keeps exact-count assertions on
	// downstream mocks (Fetch, BatchCreate, StoreEncryptedDataBatch, etc.)
	// deterministic — otherwise a race between ticker.C and ctx.Done can
	// trigger a second iteration when the test budget overlaps the ticker.
	tickerPeriod := 1 * time.Minute
	ourChainID := big.NewInt(777)

	t.Run("supports graceful shutdown", func(t *testing.T) {
		consumer := &CrossChainConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[service.CrossChainMessage], error) {
				return []msgqueue.Message[service.CrossChainMessage]{}, nil
			},
		}
		ethClient := &CrossChainEthereumClientMock{
			BlockByHashFunc: func(ctx context.Context, hash common.Hash) (*ethTypes.Block, error) {
				return nil, nil //nolint:nilnil // returning nil,nil is intentional for test stub
			},
		}
		proofGen := &CrossChainProofGeneratorMock{
			BatchGenerateFunc: func(_ context.Context, txHashes []common.Hash, routineCount int) [][]byte {
				return [][]byte{}
			},
		}
		teleportClient := &TeleportClientMock{
			StoreEncryptedDataBatchFunc: func(_ context.Context, sharedIDs []string, msgs []types.DispatchedMessageToPrivateHub, chainID *big.Int) (common.Hash, error) {
				return common.Hash{}, nil
			},
		}
		transactionRepo := &CrossChainTransactionRepositoryMock{
			BatchCreateWithStateAndOutcomeFunc: func(ctx context.Context, txs []types.Transaction, state types.TransactionState, outcome types.TransactionOutcome) error {
				return nil
			},
		}
		signatureRepo := &CrossChainSignatureRepositoryMock{
			BatchCreateFunc: func(ctx context.Context, signatures []types.CalldataSignature) error {
				return nil
			},
		}

		svc := service.NewCrossChainService(
			tickerPeriod,
			ourChainID,
			consumer,
			ethClient,
			proofGen,
			teleportClient,
			transactionRepo,
			signatureRepo,
		)
		hasGracefulShutdown := testtools.ShutdownFixture(t, svc.Run, time.Millisecond)

		assert.True(t, hasGracefulShutdown, "service should shutdown gracefully when context is cancelled")
	})

	t.Run("respects context cancelation on errors", func(t *testing.T) {
		consumer := &CrossChainConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[service.CrossChainMessage], error) {
				return []msgqueue.Message[service.CrossChainMessage]{}, errors.New("example error")
			},
		}
		ethClient := &CrossChainEthereumClientMock{
			BlockByHashFunc: func(ctx context.Context, hash common.Hash) (*ethTypes.Block, error) {
				return nil, nil //nolint:nilnil // returning nil,nil is intentional for test stub
			},
		}
		proofGen := &CrossChainProofGeneratorMock{
			BatchGenerateFunc: func(_ context.Context, txHashes []common.Hash, routineCount int) [][]byte {
				return [][]byte{}
			},
		}
		teleportClient := &TeleportClientMock{
			StoreEncryptedDataBatchFunc: func(_ context.Context, sharedIDs []string, msgs []types.DispatchedMessageToPrivateHub, chainID *big.Int) (common.Hash, error) {
				return common.Hash{}, nil
			},
		}
		transactionRepo := &CrossChainTransactionRepositoryMock{
			BatchCreateWithStateAndOutcomeFunc: func(ctx context.Context, txs []types.Transaction, state types.TransactionState, outcome types.TransactionOutcome) error {
				return nil
			},
		}
		signatureRepo := &CrossChainSignatureRepositoryMock{
			BatchCreateFunc: func(ctx context.Context, signatures []types.CalldataSignature) error {
				return nil
			},
		}

		svc := service.NewCrossChainService(
			tickerPeriod,
			ourChainID,
			consumer,
			ethClient,
			proofGen,
			teleportClient,
			transactionRepo,
			signatureRepo,
		)
		respectsContextOnError := testtools.ShutdownFixture(t, svc.Run, time.Millisecond)

		assert.True(t, respectsContextOnError, "service should respect context cancellation even on errors")
	})

	t.Run("continues to ticker on no messages to process", func(t *testing.T) {
		consumer := &CrossChainConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[service.CrossChainMessage], error) {
				return []msgqueue.Message[service.CrossChainMessage]{}, nil
			},
		}
		ethClient := &CrossChainEthereumClientMock{
			BlockByHashFunc: func(ctx context.Context, hash common.Hash) (*ethTypes.Block, error) {
				return nil, nil //nolint:nilnil // returning nil,nil is intentional for test stub
			},
		}
		proofGen := &CrossChainProofGeneratorMock{
			BatchGenerateFunc: func(_ context.Context, txHashes []common.Hash, routineCount int) [][]byte {
				return [][]byte{}
			},
		}
		teleportClient := &TeleportClientMock{
			StoreEncryptedDataBatchFunc: func(_ context.Context, sharedIDs []string, msgs []types.DispatchedMessageToPrivateHub, chainID *big.Int) (common.Hash, error) {
				return common.Hash{}, nil
			},
		}
		transactionRepo := &CrossChainTransactionRepositoryMock{
			BatchCreateWithStateAndOutcomeFunc: func(ctx context.Context, txs []types.Transaction, state types.TransactionState, outcome types.TransactionOutcome) error {
				return nil
			},
		}
		signatureRepo := &CrossChainSignatureRepositoryMock{
			BatchCreateFunc: func(ctx context.Context, signatures []types.CalldataSignature) error {
				return nil
			},
		}

		svc := service.NewCrossChainService(
			tickerPeriod,
			ourChainID,
			consumer,
			ethClient,
			proofGen,
			teleportClient,
			transactionRepo,
			signatureRepo,
		)

		// 1s safety budget — well above any scheduling jitter under load.
		// Run processes the prepared messages once, then loops back to a
		// blocking consumer.Fetch which waits for ctx cancel. The downstream
		// assertions (mock-call counts) don't expose an ack-callback to
		// synchronise on, so a generous timeout is the cleanest pattern.
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := svc.Run(ctx)
		require.Nil(t, err, "service should complete without error")

		// verify that we fetch messages
		assert.Equal(t, 1, len(consumer.FetchCalls()), "should call consumer.Fetch exactly once")

		// verify none of the dependencies get called
		assert.Equal(t, 0, len(ethClient.BlockByHashCalls()), "should not call ethClient when no messages to process")
		assert.Equal(t, 0, len(proofGen.BatchGenerateCalls()), "should not generate proofs when no messages to process")
		assert.Equal(
			t,
			0,
			len(teleportClient.StoreEncryptedDataBatchCalls()),
			"should not call teleport when no messages to process",
		)
		assert.Equal(
			t,
			0,
			len(transactionRepo.BatchCreateWithStateAndOutcomeCalls()),
			"should not persist transactions when no messages to process",
		)
	})

	t.Run("generates DispatchedMessageToPrivateHub and sends it to the Teleport contract", func(t *testing.T) {
		proof := common.Hex2Bytes("fee1dead")
		block := testdata.NewBlock()
		batchHash := common.HexToHash("0x80088008")

		msg := testdata.NewCrossChainMessage()
		messageToDispatch := testdata.NewDispatchedMessage(msg, ourChainID, block)
		transaction := testdata.NewTransaction(messageToDispatch, batchHash)

		// spy function used to verify that the message was acked
		calledAck := false
		spyAck := func(context.Context) error {
			calledAck = true
			return nil
		}

		consumer := &CrossChainConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[service.CrossChainMessage], error) {
				return []msgqueue.Message[service.CrossChainMessage]{
					{
						V:   msg,
						Ack: spyAck,
					},
				}, nil
			},
		}
		ethClient := &CrossChainEthereumClientMock{
			BlockByHashFunc: func(ctx context.Context, hash common.Hash) (*ethTypes.Block, error) {
				// verify correct hash as function parameter
				assert.Equal(t, msg.BlockHash, hash, "should request block using message's block hash")

				return block, nil
			},
		}
		proofGen := &CrossChainProofGeneratorMock{
			BatchGenerateFunc: func(_ context.Context, txHashes []common.Hash, routineCount int) [][]byte {
				// verify correct number of hashes
				require.Equal(t, 1, len(txHashes), "should generate proof for exactly one transaction hash")
				// verify correct hash is used for proof generation
				assert.Equal(t, txHashes[0], msg.TxHash, "should use message's transaction hash for proof generation")

				return [][]byte{proof}
			},
		}
		teleportClient := &TeleportClientMock{
			StoreEncryptedDataBatchFunc: func(_ context.Context, sharedIDs []string, msgs []types.DispatchedMessageToPrivateHub, chainID *big.Int) (common.Hash, error) {
				// verify we didn't ack the message before sending it to the Private Hub
				assert.False(t, calledAck, "should not ack message before sending to teleport")
				// verify we received correct number of messages
				require.Equal(t, 1, len(msgs), "should send exactly one message to teleport")
				// check deterministic fields
				assertDeterministicFieldsDMTPH(t, messageToDispatch, msgs[0])
				// check that shared ID is set
				assert.NotEmpty(t, msgs[0].SharedId, "dispatched message should have a shared ID")
				// check that batch ID is set
				assert.NotEmpty(t, msgs[0].BatchId, "dispatched message should have a batch ID")

				return batchHash, nil
			},
		}
		transactionRepo := &CrossChainTransactionRepositoryMock{
			BatchCreateWithStateAndOutcomeFunc: func(ctx context.Context, txs []types.Transaction, state types.TransactionState, outcome types.TransactionOutcome) error {
				// verify we didn't ack the message before persisting its transaction
				assert.False(t, calledAck, "should not ack message before persisting transaction")
				// verify correct number of transactions
				require.Equal(t, 1, len(txs), "should persist exactly one transaction")

				// verify deterministic fields
				assertDeterministicFieldsTransaction(t, transaction, txs[0])
				// verify batch hash from PNH is set
				assert.Equal(
					t,
					batchHash,
					txs[0].BatchPrivateHubHash,
					"transaction should have correct batch hash from teleport",
				)
				// check the correct state and outcome are set
				assert.Equal(
					t,
					types.SourcePublish,
					state,
					"transaction should be created with SourcePublish state",
				)
				assert.Equal(
					t,
					types.OutcomeSuccess,
					outcome,
					"transaction should be created with OutcomeSuccess",
				)

				return nil
			},
		}
		signatureRepo := &CrossChainSignatureRepositoryMock{
			BatchCreateFunc: func(ctx context.Context, signatures []types.CalldataSignature) error {
				// verify we didn't ack the message before persisting its signature
				assert.False(t, calledAck, "should not ack message before persisting signature")

				require.Equal(t, len(signatures), 1, "should persist exactly one signature")

				sig := signatures[0]
				assert.NotEmpty(t, sig.SharedId, "signature should have a shared ID")
				assert.Equal(
					t,
					msg.Data.MessageMetadata.RevertPayloadDataSender,
					sig.Signature,
					"signature should match message's revert payload",
				)
				assert.Equal(
					t,
					msg.Data.MessageMetadata.ResourceId,
					sig.ResourceId,
					"signature should have correct resource ID",
				)
				assert.Equal(
					t,
					types.RevertOnSenderSide,
					sig.SignatureType,
					"signature type should be RevertOnSenderSide",
				)

				return nil
			},
		}

		svc := service.NewCrossChainService(
			tickerPeriod,
			ourChainID,
			consumer,
			ethClient,
			proofGen,
			teleportClient,
			transactionRepo,
			signatureRepo,
		)

		// 1s safety budget — well above any scheduling jitter under load.
		// Run processes the prepared messages once, then loops back to a
		// blocking consumer.Fetch which waits for ctx cancel. The downstream
		// assertions (mock-call counts) don't expose an ack-callback to
		// synchronise on, so a generous timeout is the cleanest pattern.
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := svc.Run(ctx)
		require.Nil(t, err, "service should complete without error")

		// check that fetch was called
		assert.Equal(t, 1, len(consumer.FetchCalls()), "should call consumer.Fetch exactly once")
		// check teleport client was called
		assert.Equal(
			t,
			1,
			len(teleportClient.StoreEncryptedDataBatchCalls()),
			"should call teleport.StoreEncryptedDataBatch exactly once",
		)
		// check proofs were generated
		assert.Equal(t, 1, len(proofGen.BatchGenerateCalls()), "should call proofGen.BatchGenerate exactly once")
		// check we persisted transaction
		assert.Equal(
			t,
			1,
			len(transactionRepo.BatchCreateWithStateAndOutcomeCalls()),
			"should call transactionRepo.BatchCreateWithState exactly once",
		)
		// check we persisted signature
		assert.Equal(t, 1, len(signatureRepo.BatchCreateCalls()), "should call signatureRepo.BatchCreate exactly once")

		// check that message was acked
		assert.True(t, calledAck, "message should be acknowledged after successful processing")
	})

	t.Run("splits messages with different chain IDs in different batches", func(t *testing.T) {
		chainIDA := big.NewInt(1111)
		chainIDB := big.NewInt(1878)

		proof := common.Hex2Bytes("fee1dead")

		block := testdata.NewBlock()
		batchHash := common.HexToHash("0x80088008")

		msgA1 := testdata.NewCrossChainMessageWith(
			testdata.WithToChainIDOpt(chainIDA),
			testdata.WithMessageIDOpt(common.HexToHash("0xa1")),
		)
		msgA2 := testdata.NewCrossChainMessageWith(
			testdata.WithToChainIDOpt(chainIDA),
			testdata.WithMessageIDOpt(common.HexToHash("0xa2")),
		)
		msgB1 := testdata.NewCrossChainMessageWith(
			testdata.WithToChainIDOpt(chainIDB),
			testdata.WithMessageIDOpt(common.HexToHash("0xb1")),
		)
		msgB2 := testdata.NewCrossChainMessageWith(
			testdata.WithToChainIDOpt(chainIDB),
			testdata.WithMessageIDOpt(common.HexToHash("0xb2")),
		)

		dummyAck := func(context.Context) error {
			return nil
		}

		consumer := &CrossChainConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[service.CrossChainMessage], error) {
				return []msgqueue.Message[service.CrossChainMessage]{
					{
						V:   msgA1,
						Ack: dummyAck,
					},
					{
						V:   msgA2,
						Ack: dummyAck,
					},
					{
						V:   msgB1,
						Ack: dummyAck,
					},
					{
						V:   msgB2,
						Ack: dummyAck,
					},
				}, nil
			},
		}
		ethClient := &CrossChainEthereumClientMock{
			BlockByHashFunc: func(ctx context.Context, hash common.Hash) (*ethTypes.Block, error) {
				return block, nil
			},
		}
		proofGen := &CrossChainProofGeneratorMock{
			BatchGenerateFunc: func(_ context.Context, txHashes []common.Hash, routineCount int) [][]byte {
				// verify correct number of hashes
				require.Equal(t, 4, len(txHashes), "should generate proofs for all 4 transaction hashes")

				return [][]byte{proof, proof, proof, proof}
			},
		}

		seenChains := map[string]int{}
		seenMsgIDs := map[common.Hash]int{}

		teleportClient := &TeleportClientMock{
			StoreEncryptedDataBatchFunc: func(_ context.Context, sharedIDs []string, msgs []types.DispatchedMessageToPrivateHub, chainID *big.Int) (common.Hash, error) {
				require.Equal(t, 2, len(msgs), "each batch should contain exactly 2 messages for the same chain")

				// All messages in this batch must share the same chain ID
				toChain := msgs[0].ToChainId
				for _, m := range msgs {
					if m.ToChainId.Cmp(chainID) != 0 {
						assert.Fail(t, "mismatch between chainID parameter and message's ToChainId")
					}
					if m.ToChainId.Cmp(toChain) != 0 {
						assert.Fail(t, "all messages in a batch must have the same destination chain ID")
					}
					seenMsgIDs[m.MessageId]++
				}

				seenChains[toChain.String()]++

				return batchHash, nil
			},
		}
		transactionRepo := &CrossChainTransactionRepositoryMock{
			BatchCreateWithStateAndOutcomeFunc: func(ctx context.Context, txs []types.Transaction, state types.TransactionState, outcome types.TransactionOutcome) error {
				return nil
			},
		}
		signatureRepo := &CrossChainSignatureRepositoryMock{
			BatchCreateFunc: func(ctx context.Context, signatures []types.CalldataSignature) error {
				return nil
			},
		}

		svc := service.NewCrossChainService(
			tickerPeriod,
			ourChainID,
			consumer,
			ethClient,
			proofGen,
			teleportClient,
			transactionRepo,
			signatureRepo,
		)

		// 1s safety budget — well above any scheduling jitter under load.
		// Run processes the prepared messages once, then loops back to a
		// blocking consumer.Fetch which waits for ctx cancel. The downstream
		// assertions (mock-call counts) don't expose an ack-callback to
		// synchronise on, so a generous timeout is the cleanest pattern.
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := svc.Run(ctx)
		require.Nil(t, err, "service should complete without error")

		// check that fetch was called
		assert.Equal(t, 1, len(consumer.FetchCalls()), "should call consumer.Fetch exactly once")
		// check proofs were generated
		assert.Equal(t, 1, len(proofGen.BatchGenerateCalls()), "should call proofGen.BatchGenerate exactly once")
		// check we persisted transactions
		assert.NotEqual(
			t,
			0,
			len(transactionRepo.BatchCreateWithStateAndOutcomeCalls()),
			"should persist transactions at least once",
		)
		// check we persisted revert signatures
		assert.Equal(t, 1, len(signatureRepo.BatchCreateCalls()), "should call signatureRepo.BatchCreate exactly once")

		// check teleport client was called twice (once per destination chain)
		assert.Equal(
			t,
			2,
			len(teleportClient.StoreEncryptedDataBatchCalls()),
			"should call teleport twice, once for each destination chain",
		)

		// we should have exactly one batch per chain
		assert.Equal(t, 1, seenChains[chainIDA.String()], "should create exactly one batch for chain A")
		assert.Equal(t, 1, seenChains[chainIDB.String()], "should create exactly one batch for chain B")

		// and all four messages should appear exactly once
		assert.Equal(t, 1, seenMsgIDs[msgA1.MessageID], "message A1 should appear in exactly one batch")
		assert.Equal(t, 1, seenMsgIDs[msgA2.MessageID], "message A2 should appear in exactly one batch")
		assert.Equal(t, 1, seenMsgIDs[msgB1.MessageID], "message B1 should appear in exactly one batch")
		assert.Equal(t, 1, seenMsgIDs[msgB2.MessageID], "message B2 should appear in exactly one batch")
	})

	t.Run("doesn't persist signatures for non-atomic transactions", func(t *testing.T) {
		proof := common.Hex2Bytes("fee1dead")
		block := testdata.NewBlock()
		batchHash := common.HexToHash("0x80088008")

		msgAtomicA := testdata.NewCrossChainMessageWith(
			testdata.WithMessageIDOpt(common.HexToHash("0xa1")),
		)
		msgVanillaA := testdata.NewCrossChainMessageWith(
			testdata.WithMessageIDOpt(common.HexToHash("0xv1")),
			testdata.WithRevertPayloadSednerOpts(nil),
		)
		msgAtomicB := testdata.NewCrossChainMessageWith(
			testdata.WithMessageIDOpt(common.HexToHash("0xa2")),
		)
		msgVanillaB := testdata.NewCrossChainMessageWith(
			testdata.WithMessageIDOpt(common.HexToHash("0xv2")),
			testdata.WithRevertPayloadSednerOpts(nil),
		)

		dummyAck := func(context.Context) error {
			return nil
		}
		consumer := &CrossChainConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[service.CrossChainMessage], error) {
				return []msgqueue.Message[service.CrossChainMessage]{
					{
						V:   msgAtomicA,
						Ack: dummyAck,
					},
					{
						V:   msgVanillaA,
						Ack: dummyAck,
					},
					{
						V:   msgAtomicB,
						Ack: dummyAck,
					},
					{
						V:   msgVanillaB,
						Ack: dummyAck,
					},
				}, nil
			},
		}
		ethClient := &CrossChainEthereumClientMock{
			BlockByHashFunc: func(ctx context.Context, hash common.Hash) (*ethTypes.Block, error) {
				return block, nil
			},
		}
		proofGen := &CrossChainProofGeneratorMock{
			BatchGenerateFunc: func(_ context.Context, txHashes []common.Hash, routineCount int) [][]byte {
				return [][]byte{proof, proof, proof, proof}
			},
		}
		teleportClient := &TeleportClientMock{
			StoreEncryptedDataBatchFunc: func(_ context.Context, sharedIDs []string, msgs []types.DispatchedMessageToPrivateHub, chainID *big.Int) (common.Hash, error) {
				return batchHash, nil
			},
		}
		transactionRepo := &CrossChainTransactionRepositoryMock{
			BatchCreateWithStateAndOutcomeFunc: func(ctx context.Context, txs []types.Transaction, state types.TransactionState, outcome types.TransactionOutcome) error {
				return nil
			},
		}
		signatureRepo := &CrossChainSignatureRepositoryMock{
			BatchCreateFunc: func(ctx context.Context, signatures []types.CalldataSignature) error {
				require.Equal(t, len(signatures), 2, "should persist exactly 2 signatures (only for atomic messages)")

				assert.NotEmpty(t, signatures[0].Signature, "first signature should not be empty")
				assert.NotEmpty(t, signatures[1].Signature, "second signature should not be empty")

				return nil
			},
		}

		svc := service.NewCrossChainService(
			tickerPeriod,
			ourChainID,
			consumer,
			ethClient,
			proofGen,
			teleportClient,
			transactionRepo,
			signatureRepo,
		)

		// 1s safety budget — well above any scheduling jitter under load.
		// Run processes the prepared messages once, then loops back to a
		// blocking consumer.Fetch which waits for ctx cancel. The downstream
		// assertions (mock-call counts) don't expose an ack-callback to
		// synchronise on, so a generous timeout is the cleanest pattern.
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := svc.Run(ctx)
		require.Nil(t, err, "service should complete without error")

		// check we persisted revert signatures only for atomic messages
		assert.Equal(t, 1, len(signatureRepo.BatchCreateCalls()), "should call signatureRepo.BatchCreate exactly once")
	})

	t.Run("doesn't ack messages on fail to create transactions in repository", func(t *testing.T) {
		proof := common.Hex2Bytes("fee1dead")
		block := testdata.NewBlock()

		msg1 := testdata.NewCrossChainMessageWith(
			testdata.WithMessageIDOpt(common.HexToHash("0xa1")),
		)
		msg2 := testdata.NewCrossChainMessageWith(
			testdata.WithMessageIDOpt(common.HexToHash("0xa2")),
		)

		// Track which messages were acked
		ackedMessages := make(map[common.Hash]bool)
		makeAckSpy := func(msgID common.Hash) func(context.Context) error {
			return func(context.Context) error {
				ackedMessages[msgID] = true
				return nil
			}
		}

		consumer := &CrossChainConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[service.CrossChainMessage], error) {
				return []msgqueue.Message[service.CrossChainMessage]{
					{
						V:   msg1,
						Ack: makeAckSpy(msg1.MessageID),
					},
					{
						V:   msg2,
						Ack: makeAckSpy(msg2.MessageID),
					},
				}, nil
			},
		}
		ethClient := &CrossChainEthereumClientMock{
			BlockByHashFunc: func(ctx context.Context, hash common.Hash) (*ethTypes.Block, error) {
				return block, nil
			},
		}
		proofGen := &CrossChainProofGeneratorMock{
			BatchGenerateFunc: func(_ context.Context, txHashes []common.Hash, routineCount int) [][]byte {
				require.Equal(t, 2, len(txHashes), "should generate proofs for both transaction hashes")
				return [][]byte{proof, proof}
			},
		}
		teleportClient := &TeleportClientMock{
			StoreEncryptedDataBatchFunc: func(_ context.Context, sharedIDs []string, msgs []types.DispatchedMessageToPrivateHub, chainID *big.Int) (common.Hash, error) {
				return common.HexToHash("0x80088008"), nil
			},
		}
		transactionRepo := &CrossChainTransactionRepositoryMock{
			BatchCreateWithStateAndOutcomeFunc: func(ctx context.Context, txs []types.Transaction, state types.TransactionState, outcome types.TransactionOutcome) error {
				// Verify we got both transactions
				require.Equal(t, 2, len(txs), "should attempt to persist both transactions")
				// Return an error to simulate failure
				return assert.AnError
			},
		}
		signatureRepo := &CrossChainSignatureRepositoryMock{
			BatchCreateFunc: func(ctx context.Context, signatures []types.CalldataSignature) error {
				return nil
			},
		}

		svc := service.NewCrossChainService(
			tickerPeriod,
			ourChainID,
			consumer,
			ethClient,
			proofGen,
			teleportClient,
			transactionRepo,
			signatureRepo,
		)

		// 1s safety budget — well above any scheduling jitter under load.
		// Run processes the prepared messages once, then loops back to a
		// blocking consumer.Fetch which waits for ctx cancel. The downstream
		// assertions (mock-call counts) don't expose an ack-callback to
		// synchronise on, so a generous timeout is the cleanest pattern.
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := svc.Run(ctx)
		require.Nil(t, err, "service should complete without error")

		// verify that we attempted to persist transaction
		assert.Equal(
			t,
			1,
			len(transactionRepo.BatchCreateWithStateAndOutcomeCalls()),
			"should attempt to persist transactions once",
		)

		// check that NO messages were acked due to transaction persistence failure
		assert.False(
			t,
			ackedMessages[msg1.MessageID],
			"message should not be acknowledged when transaction repository fails",
		)
		assert.False(
			t,
			ackedMessages[msg2.MessageID],
			"message should not be acknowledged when transaction repository fails",
		)
	})

	t.Run("persists and acks both successful and failed batches with appropriate states", func(t *testing.T) {
		chainSuccess := big.NewInt(1111)
		chainFail := big.NewInt(2222)

		proof := common.Hex2Bytes("fee1dead")
		block := testdata.NewBlock()
		batchHash := common.HexToHash("0x80088008")

		msgSuccess1 := testdata.NewCrossChainMessageWith(
			testdata.WithToChainIDOpt(chainSuccess),
			testdata.WithMessageIDOpt(common.HexToHash("0xa1")),
		)
		msgSuccess2 := testdata.NewCrossChainMessageWith(
			testdata.WithToChainIDOpt(chainSuccess),
			testdata.WithMessageIDOpt(common.HexToHash("0xa2")),
		)
		msgFail1 := testdata.NewCrossChainMessageWith(
			testdata.WithToChainIDOpt(chainFail),
			testdata.WithMessageIDOpt(common.HexToHash("0xb1")),
		)
		msgFail2 := testdata.NewCrossChainMessageWith(
			testdata.WithToChainIDOpt(chainFail),
			testdata.WithMessageIDOpt(common.HexToHash("0xb2")),
		)

		// Track which messages were acked
		ackedMessages := make(map[common.Hash]bool)
		makeAckSpy := func(msgID common.Hash) func(context.Context) error {
			return func(context.Context) error {
				ackedMessages[msgID] = true
				return nil
			}
		}

		consumer := &CrossChainConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[service.CrossChainMessage], error) {
				return []msgqueue.Message[service.CrossChainMessage]{
					{
						V:   msgSuccess1,
						Ack: makeAckSpy(msgSuccess1.MessageID),
					},
					{
						V:   msgSuccess2,
						Ack: makeAckSpy(msgSuccess2.MessageID),
					},
					{
						V:   msgFail1,
						Ack: makeAckSpy(msgFail1.MessageID),
					},
					{
						V:   msgFail2,
						Ack: makeAckSpy(msgFail2.MessageID),
					},
				}, nil
			},
		}
		ethClient := &CrossChainEthereumClientMock{
			BlockByHashFunc: func(ctx context.Context, hash common.Hash) (*ethTypes.Block, error) {
				return block, nil
			},
		}
		proofGen := &CrossChainProofGeneratorMock{
			BatchGenerateFunc: func(_ context.Context, txHashes []common.Hash, routineCount int) [][]byte {
				require.Equal(t, 4, len(txHashes), "should generate proofs for all 4 transaction hashes")
				return [][]byte{proof, proof, proof, proof}
			},
		}
		teleportClient := &TeleportClientMock{
			StoreEncryptedDataBatchFunc: func(_ context.Context, sharedIDs []string, msgs []types.DispatchedMessageToPrivateHub, chainID *big.Int) (common.Hash, error) {
				// Succeed for chainSuccess, fail for chainFail
				if chainID.Cmp(chainSuccess) == 0 {
					return batchHash, nil
				}
				return common.Hash{}, assert.AnError
			},
		}
		transactionRepo := &CrossChainTransactionRepositoryMock{
			BatchCreateWithStateAndOutcomeFunc: func(ctx context.Context, txs []types.Transaction, state types.TransactionState, outcome types.TransactionOutcome) error {
				// Should be called twice - once for successful chain, once for failed chain
				if len(txs) == 2 {
					// Could be either successful or failed batch
					firstTxChainID := txs[0].ToChainID
					if firstTxChainID.Cmp(chainSuccess) == 0 {
						assert.Equal(
							t,
							types.SourcePublish,
							state,
							"successful messages should be persisted with SourcePublish state",
						)
						assert.Equal(
							t,
							types.OutcomeSuccess,
							outcome,
							"successful messages should be persisted with OutcomeSuccess",
						)
						assert.Equal(
							t,
							batchHash,
							txs[0].BatchPrivateHubHash,
							"successful messages should have the batch hash from teleport",
						)
						for _, tx := range txs {
							assert.Equal(
								t,
								chainSuccess,
								tx.ToChainID,
								"all transactions in successful batch should be for successful chain",
							)
						}
					} else if firstTxChainID.Cmp(chainFail) == 0 {
						assert.Equal(
							t,
							types.SourcePublish,
							state,
							"failed messages should be persisted with SourcePublish state",
						)
						assert.Equal(
							t,
							types.OutcomeFailed,
							outcome,
							"failed messages should be persisted with OutcomeFailed",
						)
						assert.Equal(
							t,
							common.Hash{},
							txs[0].BatchPrivateHubHash,
							"failed messages should have empty batch hash",
						)
						for _, tx := range txs {
							assert.Equal(
								t,
								chainFail,
								tx.ToChainID,
								"all transactions in failed batch should be for failed chain",
							)
						}
					}
				}
				return nil
			},
		}
		signatureRepo := &CrossChainSignatureRepositoryMock{
			BatchCreateFunc: func(ctx context.Context, signatures []types.CalldataSignature) error {
				return nil
			},
		}

		svc := service.NewCrossChainService(
			tickerPeriod,
			ourChainID,
			consumer,
			ethClient,
			proofGen,
			teleportClient,
			transactionRepo,
			signatureRepo,
		)

		// 1s safety budget — well above any scheduling jitter under load.
		// Run processes the prepared messages once, then loops back to a
		// blocking consumer.Fetch which waits for ctx cancel. The downstream
		// assertions (mock-call counts) don't expose an ack-callback to
		// synchronise on, so a generous timeout is the cleanest pattern.
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := svc.Run(ctx)
		require.Nil(t, err, "service should complete without error")

		// verify that we attempted to send to teleport for both chains
		assert.Equal(
			t,
			2,
			len(teleportClient.StoreEncryptedDataBatchCalls()),
			"should attempt to send to teleport for both chains",
		)

		// verify transaction repository was called for both successful and failed messages
		assert.Equal(
			t,
			2,
			len(transactionRepo.BatchCreateWithStateAndOutcomeCalls()),
			"should persist transactions for both successful and failed chains",
		)

		// Verify all messages were acked (both successful and failed)
		assert.True(t, ackedMessages[msgSuccess1.MessageID], "message for successful chain should be acked")
		assert.True(t, ackedMessages[msgSuccess2.MessageID], "message for successful chain should be acked")
		assert.True(
			t,
			ackedMessages[msgFail1.MessageID],
			"message for failed chain should be acked after persisting with failed state",
		)
		assert.True(
			t,
			ackedMessages[msgFail2.MessageID],
			"message for failed chain should be acked after persisting with failed state",
		)
	})

	t.Run("doesn't ack messages on fail to create revert signature in repository", func(t *testing.T) {
		proof := common.Hex2Bytes("fee1dead")
		block := testdata.NewBlock()
		batchHash := common.HexToHash("0x80088008")

		msg1 := testdata.NewCrossChainMessageWith(
			testdata.WithMessageIDOpt(common.HexToHash("0xa1")),
		)
		msg2 := testdata.NewCrossChainMessageWith(
			testdata.WithMessageIDOpt(common.HexToHash("0xa2")),
		)

		// Track which messages were acked
		ackedMessages := make(map[common.Hash]bool)
		makeAckSpy := func(msgID common.Hash) func(context.Context) error {
			return func(context.Context) error {
				ackedMessages[msgID] = true
				return nil
			}
		}

		consumer := &CrossChainConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[service.CrossChainMessage], error) {
				return []msgqueue.Message[service.CrossChainMessage]{
					{
						V:   msg1,
						Ack: makeAckSpy(msg1.MessageID),
					},
					{
						V:   msg2,
						Ack: makeAckSpy(msg2.MessageID),
					},
				}, nil
			},
		}
		ethClient := &CrossChainEthereumClientMock{
			BlockByHashFunc: func(ctx context.Context, hash common.Hash) (*ethTypes.Block, error) {
				return block, nil
			},
		}
		proofGen := &CrossChainProofGeneratorMock{
			BatchGenerateFunc: func(_ context.Context, txHashes []common.Hash, routineCount int) [][]byte {
				require.Equal(t, 2, len(txHashes), "should generate proofs for both transaction hashes")
				return [][]byte{proof, proof}
			},
		}
		teleportClient := &TeleportClientMock{
			StoreEncryptedDataBatchFunc: func(_ context.Context, sharedIDs []string, msgs []types.DispatchedMessageToPrivateHub, chainID *big.Int) (common.Hash, error) {
				return batchHash, nil
			},
		}
		transactionRepo := &CrossChainTransactionRepositoryMock{
			BatchCreateWithStateAndOutcomeFunc: func(ctx context.Context, txs []types.Transaction, state types.TransactionState, outcome types.TransactionOutcome) error {
				return nil
			},
		}
		signatureRepo := &CrossChainSignatureRepositoryMock{
			BatchCreateFunc: func(ctx context.Context, signatures []types.CalldataSignature) error {
				// Verify we got both signatures
				require.Equal(t, 2, len(signatures), "should attempt to persist both signatures")
				// Return an error to simulate failure
				return assert.AnError
			},
		}

		svc := service.NewCrossChainService(
			tickerPeriod,
			ourChainID,
			consumer,
			ethClient,
			proofGen,
			teleportClient,
			transactionRepo,
			signatureRepo,
		)

		// 1s safety budget — well above any scheduling jitter under load.
		// Run processes the prepared messages once, then loops back to a
		// blocking consumer.Fetch which waits for ctx cancel. The downstream
		// assertions (mock-call counts) don't expose an ack-callback to
		// synchronise on, so a generous timeout is the cleanest pattern.
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := svc.Run(ctx)
		require.Nil(t, err, "service should complete without error")

		// verify that we attempted to persist revert signature
		assert.Equal(t, 1, len(signatureRepo.BatchCreateCalls()), "should attempt to persist revert signatures once")

		// check that NO messages were acked due to signature persistence failure
		assert.False(
			t,
			ackedMessages[msg1.MessageID],
			"message should not be acknowledged when signature repository fails",
		)
		assert.False(
			t,
			ackedMessages[msg2.MessageID],
			"message should not be acknowledged when signature repository fails",
		)
	})

	t.Run("records proof_invalid failure and acks message when proof generation fails", func(t *testing.T) {
		block := testdata.NewBlock()
		batchHash := common.HexToHash("0x80088008")

		goodTxHash := common.HexToHash("0x600d600d")
		badTxHash := common.HexToHash("0xbadbad")
		goodProof := common.Hex2Bytes("fee1dead")

		goodAcked := false
		badAcked := false

		msgGood := testdata.NewQueueMessage().
			WithMessageID(common.HexToHash("0x01")).
			WithTxHash(goodTxHash).
			WithAckFunc(func(context.Context) error { goodAcked = true; return nil }).
			Build()
		msgBad := testdata.NewQueueMessage().
			WithMessageID(common.HexToHash("0x02")).
			WithTxHash(badTxHash).
			WithAckFunc(func(context.Context) error { badAcked = true; return nil }).
			Build()

		consumer := &CrossChainConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[service.CrossChainMessage], error) {
				return []msgqueue.Message[service.CrossChainMessage]{msgGood, msgBad}, nil
			},
		}
		ethClient := &CrossChainEthereumClientMock{
			BlockByHashFunc: func(ctx context.Context, hash common.Hash) (*ethTypes.Block, error) {
				return block, nil
			},
		}
		// Valid proof for the good tx, nil proof for the bad tx (keyed by hash so
		// the assertion holds regardless of routine ordering).
		proofGen := &CrossChainProofGeneratorMock{
			BatchGenerateFunc: func(_ context.Context, txHashes []common.Hash, routineCount int) [][]byte {
				out := make([][]byte, len(txHashes))
				for i, h := range txHashes {
					if h == goodTxHash {
						out[i] = goodProof
					}
				}
				return out
			},
		}
		teleportClient := &TeleportClientMock{
			StoreEncryptedDataBatchFunc: func(_ context.Context, sharedIDs []string, msgs []types.DispatchedMessageToPrivateHub, chainID *big.Int) (common.Hash, error) {
				// a nil proof must never leave the source: only the valid-proof
				// message may ship.
				require.Len(t, msgs, 1, "only the message with a valid proof should be dispatched")
				assert.Equal(t, goodTxHash, msgs[0].TxHashSource, "dispatched msg must be the valid-proof one")
				assert.Equal(t, goodProof, msgs[0].Proofs, "dispatched message must carry the valid proof")
				return batchHash, nil
			},
		}
		transactionRepo := &CrossChainTransactionRepositoryMock{
			BatchCreateWithStateAndOutcomeFunc: func(ctx context.Context, txs []types.Transaction, state types.TransactionState, outcome types.TransactionOutcome) error {
				return nil
			},
		}
		signatureRepo := &CrossChainSignatureRepositoryMock{
			BatchCreateFunc: func(ctx context.Context, signatures []types.CalldataSignature) error {
				return nil
			},
		}

		svc := service.NewCrossChainService(
			tickerPeriod,
			ourChainID,
			consumer,
			ethClient,
			proofGen,
			teleportClient,
			transactionRepo,
			signatureRepo,
		)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := svc.Run(ctx)
		require.Nil(t, err, "service should complete without error")

		// a nil proof must never be dispatched
		assert.Equal(t, 1, len(teleportClient.StoreEncryptedDataBatchCalls()), "should dispatch exactly one batch")

		// two persistence calls: the proof_invalid failure and the successful dispatch
		var failedTx, successTx *types.Transaction
		for _, call := range transactionRepo.BatchCreateWithStateAndOutcomeCalls() {
			require.Len(t, call.Txs, 1, "each persistence call should carry exactly one transaction")
			tx := call.Txs[0]
			switch call.Outcome {
			case types.OutcomeFailed:
				failedTx = &tx
			case types.OutcomeSuccess:
				successTx = &tx
			}
		}

		require.NotNil(t, failedTx, "the proof-generation failure must be persisted as a failed transaction")
		assert.True(t, failedTx.ProofInvalid, "the failed transaction must be flagged proof_invalid")
		assert.Equal(t, badTxHash.String(), failedTx.TxHash, "the proof_invalid record must be for the failing tx")

		require.NotNil(t, successTx, "the dispatched message must be persisted as a success")
		assert.False(t, successTx.ProofInvalid, "the dispatched transaction must not be flagged proof_invalid")

		// both messages are acked: the valid one after dispatch, the failing one
		// after its proof_invalid failure is durably recorded — neither is left to
		// silently churn or drop from the queue.
		assert.True(t, goodAcked, "message with a valid proof should be acknowledged after dispatch")
		assert.True(t, badAcked, "message with a failed proof should be acked after being recorded proof_invalid")
	})
}

func assertDeterministicFieldsDMTPH(t testing.TB, want, got types.DispatchedMessageToPrivateHub) {
	t.Helper()

	// Message routing fields
	assert.Equal(t, want.MessageId, got.MessageId, "MessageId mismatch")
	assert.Equal(t, want.FromChainId, got.FromChainId, "FromChainId mismatch")
	assert.Equal(t, want.From, got.From, "From mismatch")
	assert.Equal(t, want.ToChainId, got.ToChainId, "ToChainId mismatch")
	assert.Equal(t, want.To, got.To, "To mismatch")
	assert.Equal(t, want.Data, got.Data, "Data mismatch")

	// Transaction type fields
	assert.Equal(t, want.TransactionType, got.TransactionType, "TransactionType mismatch")
	assert.Equal(t, want.IsAtomic, got.IsAtomic, "IsAtomic mismatch")

	// Block and log location fields
	assert.Equal(t, want.BlockNumber, got.BlockNumber, "BlockNumber mismatch")
	assert.Equal(t, want.BlockHash, got.BlockHash, "BlockHash mismatch")
	assert.Equal(t, want.LogIdx, got.LogIdx, "LogIdx mismatch")
	assert.Equal(t, want.ParentHash, got.ParentHash, "ParentHash mismatch")

	// Transaction hash and status fields
	assert.Equal(t, want.TxHashSource, got.TxHashSource, "TxHashSource mismatch")
	assert.Equal(t, want.TxHashSourceStatus, got.TxHashSourceStatus, "TxHashSourceStatus mismatch")
	assert.Equal(t, want.TxHashSourceTimestamp, got.TxHashSourceTimestamp, "TxHashSourceTimestamp mismatch")

	// Resource and token fields
	assert.Equal(t, want.ResourceId, got.ResourceId, "ResourceId mismatch")
	assert.Equal(t, want.TokenAddress, got.TokenAddress, "TokenAddress mismatch")

	// Proof fields
	assert.Equal(t, want.Proofs, got.Proofs, "Proofs mismatch")
	assert.Equal(t, want.TxLocation, got.TxLocation, "TxLocation mismatch")
	assert.Equal(t, want.TxTrieProof, got.TxTrieProof, "TxTrieProof mismatch")
}

func assertDeterministicFieldsTransaction(t testing.TB, want, got types.Transaction) {
	t.Helper()

	// Message ID and routing fields
	assert.Equal(t, want.MsgID, got.MsgID, "MsgID mismatch")
	assert.Equal(t, want.FromChainID, got.FromChainID, "FromChainID mismatch")
	assert.Equal(t, want.FromContractAddress, got.FromContractAddress, "FromContractAddress mismatch")
	assert.Equal(t, want.FromUserAddress, got.FromUserAddress, "FromUserAddress mismatch")
	assert.Equal(t, want.ToChainID, got.ToChainID, "ToChainID mismatch")

	// Resource and atomic fields
	assert.Equal(t, want.ResourceID, got.ResourceID, "ResourceID mismatch")
	assert.Equal(t, want.IsAtomic, got.IsAtomic, "IsAtomic mismatch")

	// Block and transaction location fields
	assert.Equal(t, want.ParentHash, got.ParentHash, "ParentHash mismatch")
	assert.Equal(t, want.BlockNumber, got.BlockNumber, "BlockNumber mismatch")
	assert.Equal(t, want.TxHash, got.TxHash, "TxHash mismatch")
	assert.Equal(t, want.LogIndex, got.LogIndex, "LogIndex mismatch")

	// Timestamp field
	assert.Equal(t, want.UpdatedAt, got.UpdatedAt, "UpdatedAt mismatch")

	// Transfer metadata fields
	assert.Equal(t, want.TransferID, got.TransferID, "TransferID mismatch")
	assert.Equal(t, want.TransferAmount, got.TransferAmount, "TransferAmount mismatch")
}
