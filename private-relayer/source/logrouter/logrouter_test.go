package logrouter_test

// import (
// 	"context"
// 	"crypto/ecdsa"
// 	"fmt"
// 	"io"
// 	"log/slog"
// 	"math/big"
// 	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EndpointV1"
// 	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/MessageReceiver"
// 	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/MessageSender"
// 	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/ParticipantStorageReplicaV1"
// 	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/ResourceManager"
// 	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/TokenRegistryReplicaV1"
// 	"github.com/raylsnetwork/rayls-privacy-relayer-api/listener"
// 	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/logparser"
// 	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/service"
// 	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
// 	"testing"
// 	"time"

// 	"github.com/ethereum/go-ethereum/accounts/abi/bind"
// 	"github.com/ethereum/go-ethereum/common"
// 	ethTypes "github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/crypto"
// 	"github.com/ethereum/go-ethereum/ethclient/simulated"
// 	"github.com/stretchr/testify/suite"
// )

// var emptyAddress = common.Address{}

// type LastProcessedBlockRepositoryMock struct {
// 	block *big.Int
// }

// func (l *LastProcessedBlockRepositoryMock) Create(ctx context.Context, service types.LastProcessedBlockDocument, block *big.Int) error {
// 	l.block = block
// 	return nil
// }

// func (l *LastProcessedBlockRepositoryMock) Get(ctx context.Context, service types.LastProcessedBlockDocument) (*big.Int, error) {
// 	if l.block == nil {
// 		return nil, types.ErrLastProcessedBlockNotFound
// 	}
// 	return l.block, nil
// }

// func (l *LastProcessedBlockRepositoryMock) Update(ctx context.Context, service types.LastProcessedBlockDocument, block *big.Int) error {
// 	l.block = block
// 	return nil
// }

// type LogParserTestSuite struct {
// 	suite.Suite
// 	chainID            *big.Int
// 	privateHubChainID      *big.Int
// 	destinationChainID *big.Int

// 	privateKey *ecdsa.PrivateKey
// 	auth       *bind.TransactOpts

// 	genesisAllocations ethTypes.GenesisAlloc
// 	ledger             *simulated.Backend

// 	resourceManager        *ResourceManager.ResourceManager
// 	resourceManagerAddress common.Address

// 	participantStorage        *ParticipantStorageReplicaV1.ParticipantStorageReplicaV1
// 	participantStorageAddress common.Address

// 	tokenRegistry        *TokenRegistryReplicaV1.TokenRegistryReplicaV1
// 	tokenRegistryAddress common.Address

// 	messageReceiver        *MessageReceiver.MessageReceiver
// 	messageReceiverAddress common.Address

// 	messageSender        *MessageSender.MessageSender
// 	messageSenderAddress common.Address

// 	endpointContract *EndpointV1.EndpointV1
// 	endpointAddress  common.Address
// }

// // this function executes before the test suite begins execution
// func (suite *LogParserTestSuite) SetupSuite() {
// 	var err error

// 	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

// 	// Setup private key and genesis for simulated blockchain
// 	suite.chainID = big.NewInt(1337)
// 	suite.privateHubChainID = big.NewInt(1000)
// 	suite.destinationChainID = big.NewInt(1888)

// 	suite.privateKey, err = crypto.GenerateKey()
// 	suite.Require().Nil(err)

// 	suite.auth, err = bind.NewKeyedTransactorWithChainID(suite.privateKey, suite.chainID)
// 	suite.Require().Nil(err)

// 	suite.genesisAllocations = make(ethTypes.GenesisAlloc)
// 	suite.genesisAllocations[suite.auth.From] = ethTypes.Account{Balance: big.NewInt(10000000000000000)}

// 	// Setup simulated blockchain
// 	suite.ledger = simulated.NewBackend(suite.genesisAllocations)
// 	suite.Require().NotNil(suite.ledger)
// }

// // this function executes before each test case
// func (suite *LogParserTestSuite) SetupTest() {
// 	var err error
// 	// ResourceManager (re)deployment
// 	suite.resourceManagerAddress, _, suite.resourceManager, err = ResourceManager.DeployResourceManager(suite.auth, suite.ledger.Client(), emptyAddress)
// 	suite.Require().Nil(err)
// 	suite.ledger.Commit()

// 	// ParticipantStorage (re)deployment
// 	suite.participantStorageAddress, _, suite.participantStorage, err = ParticipantStorageReplicaV1.DeployParticipantStorageReplicaV1(suite.auth, suite.ledger.Client())
// 	suite.Require().Nil(err)
// 	suite.ledger.Commit()

// 	// TokenRegistry (re)deploymen
// 	suite.tokenRegistryAddress, _, suite.tokenRegistry, err = TokenRegistryReplicaV1.DeployTokenRegistryReplicaV1(suite.auth, suite.ledger.Client())
// 	suite.Require().Nil(err)
// 	suite.ledger.Commit()

// 	// MessageSender (re)deployment
// 	suite.messageSenderAddress, _, suite.messageSender, err = MessageSender.DeployMessageSender(suite.auth, suite.ledger.Client(), suite.chainID, suite.privateHubChainID, suite.participantStorageAddress, suite.tokenRegistryAddress)
// 	suite.Require().Nil(err)
// 	suite.ledger.Commit()

// 	// MessageReceiver (re)deployment
// 	suite.messageReceiverAddress, _, suite.messageReceiver, err = MessageReceiver.DeployMessageReceiver(suite.auth, suite.ledger.Client(), emptyAddress, emptyAddress)
// 	suite.Require().Nil(err)
// 	suite.ledger.Commit()

// 	// Endpoint (re)deployment
// 	suite.endpointAddress, _, suite.endpointContract, err = EndpointV1.DeployEndpointV1(suite.auth, suite.ledger.Client())
// 	suite.Require().Nil(err)
// 	suite.ledger.Commit()

// 	// Initialize the endpoint
// 	_, err = suite.endpointContract.Initialize(suite.auth, suite.auth.From, suite.chainID, big.NewInt(1000), big.NewInt(100))
// 	suite.Require().Nil(err)
// 	suite.ledger.Commit()

// 	// Configure the endpoint
// 	_, err = suite.endpointContract.ConfigureEndpoint(suite.auth, emptyAddress, suite.participantStorageAddress, suite.tokenRegistryAddress, suite.resourceManagerAddress, suite.messageSenderAddress, suite.messageSenderAddress, emptyAddress)
// 	suite.Require().Nil(err)
// 	suite.ledger.Commit()

// 	// Initialize the ParticipantStorage
// 	_, err = suite.participantStorage.Initialize(suite.auth, suite.auth.From, suite.endpointAddress)
// 	suite.Require().Nil(err)
// 	suite.ledger.Commit()

// 	// Initialize the TokenRegistry
// 	_, err = suite.tokenRegistry.Initialize(suite.auth, suite.auth.From, suite.endpointAddress)
// 	suite.Require().Nil(err)
// 	suite.ledger.Commit()

// 	// Add authorized address
// 	_, err = suite.endpointContract.AddAuthorizedAddresses(suite.auth, []common.Address{suite.auth.From})
// 	suite.Require().Nil(err)
// 	suite.ledger.Commit()
// }

// // this function executes after each test case
// func (suite *LogParserTestSuite) TearDownTest() {
// 	// Nothing to clean up
// }

// // this function executes after all tests executed
// func (suite *LogParserTestSuite) TearDownSuite() {
// 	// Nothing to clean up
// }

// func (suite *LogParserTestSuite) TestRoutesPrivateHubEvents() {
// 	blockNumber, _ := suite.ledger.Client().BlockNumber(suite.T().Context())

// 	_, err := suite.participantStorage.RequestAllParticipantsDataFromPrivateHub(suite.auth)
// 	suite.Require().Nil(err)
// 	suite.ledger.Commit()

// 	lpbRepo := &LastProcessedBlockRepositoryMock{}

// 	eventBlockNumber, _ := suite.ledger.Client().BlockNumber(suite.T().Context())
// 	messageID := fmt.Sprintf("%d-%d-%d", eventBlockNumber, 0, 0)

// 	wantMessage := service.PrivateHubMessage{
// 		ID: messageID,
// 	}
// 	var gotMessage service.PrivateHubMessage
// 	privateHubMQ := &PrivateHubMQMock{
// 		PushFunc: func(message service.PrivateHubMessage) error {
// 			gotMessage = message
// 			return nil
// 		},
// 	}

// 	listener, err := suite.createPrivateRelayerListener(
// 		new(big.Int).SetUint64(blockNumber), 10, lpbRepo, privateHubMQ,
// 	)

// 	ctx, cancel := context.WithTimeout(suite.T().Context(), time.Second)
// 	defer cancel()

// 	err = listener.Run(ctx)
// 	suite.Require().Nil(err)

// 	suite.Equal(wantMessage, gotMessage)
// }

// func (suite *LogParserTestSuite) TestRoutesCrossChainEvents() {
// 	blockNumber, _ := suite.ledger.Client().BlockNumber(suite.T().Context())

// 	_, err := suite.endpointContract.Send(suite.auth, suite.destinationChainID, emptyAddress, []byte{})
// 	suite.Require().Nil(err)
// 	suite.ledger.Commit()

// 	lpbRepo := &LastProcessedBlockRepositoryMock{}

// 	eventBlockNumber, _ := suite.ledger.Client().BlockNumber(suite.T().Context())
// 	messageID := fmt.Sprintf("%d-%d-%d", eventBlockNumber, 0, 0)

// 	wantMessage := service.PrivateHubMessage{
// 		ID: messageID,
// 	}
// 	var gotMessage service.PrivateHubMessage
// 	privateHubMQ := &PrivateHubMQMock{
// 		PushFunc: func(message service.PrivateHubMessage) error {
// 			gotMessage = message
// 			return nil
// 		},
// 	}

// 	listener, err := suite.createPrivateRelayerListener(
// 		new(big.Int).SetUint64(blockNumber), 10, lpbRepo, privateHubMQ,
// 	)

// 	ctx, cancel := context.WithTimeout(suite.T().Context(), time.Second)
// 	defer cancel()

// 	err = listener.Run(ctx)
// 	suite.Require().Nil(err)

// 	suite.Equal(wantMessage, gotMessage)
// }

// // Helper function to create an endpoint listener using the new LogListener and EndpointWrapper
// func (suite *LogParserTestSuite) createPrivateRelayerListener(
// 	startingBlock *big.Int,
// 	batchSize int,
// 	lpbRepo *LastProcessedBlockRepositoryMock,
// 	privateHubMQ *PrivateHubMQMock,
// ) (*listener.LogListener, error) {

// 	parser, err := logparser.New(privateHubMQ)
// 	suite.Require().Nil(err)
// 	listenerConfig := listener.LogListenerConfig{
// 		Component: types.LastProcessedBlockDocumentPublicChain,

// 		StartingBlock: startingBlock,
// 		BatchSize:     batchSize,
// 		Addresses: []common.Address{
// 			suite.endpointAddress,
// 		},
// 	}
// 	return listener.NewLogListener(
// 		listenerConfig,
// 		parser,
// 		suite.ledger.Client(),
// 		lpbRepo,
// 	)
// }

// func TestLogParserSuite(t *testing.T) {
// 	suite.Run(t, new(LogParserTestSuite))
// }
