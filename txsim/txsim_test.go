package txsim_test

import (
	"context"
	"crypto/ecdsa"
	"io"
	"log/slog"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/raylsnetwork/rayls-sovereign-relayer/txsim"
	"github.com/raylsnetwork/rayls-sovereign-relayer/txsim/testdata/Dummy"
	"github.com/stretchr/testify/suite"
)

//go:generate solc --optimize --bin --abi --overwrite --output-dir ./testdata ./testdata/Dummy.sol
//go:generate abigen --bin ./testdata/DummyContract.bin --abi ./testdata/DummyContract.abi --pkg Dummy --out ./testdata/Dummy/Dummy.go

type TxSimTestSuite struct {
	suite.Suite
	chainID            *big.Int
	privateKey         *ecdsa.PrivateKey
	auth               *bind.TransactOpts
	genesisAllocations ethTypes.GenesisAlloc
	ledger             *simulated.Backend

	dummyABI             abi.ABI
	dummyContractAddress common.Address
}

// this function executes before the test suite begins execution
func (suite *TxSimTestSuite) SetupSuite() {
	var err error

	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	// Setup private key and genesis for simulated blockchain
	suite.chainID = big.NewInt(1337)

	suite.privateKey, err = crypto.GenerateKey()
	suite.Require().Nil(err)

	suite.auth, err = bind.NewKeyedTransactorWithChainID(suite.privateKey, suite.chainID)
	suite.Require().Nil(err)

	suite.genesisAllocations = make(ethTypes.GenesisAlloc)
	suite.genesisAllocations[suite.auth.From] = ethTypes.Account{Balance: big.NewInt(10000000000000000)}

	// Setup simulated blockchain
	suite.ledger = simulated.NewBackend(suite.genesisAllocations)
	suite.Require().NotNil(suite.ledger)
}

// this function executes before each test case
func (suite *TxSimTestSuite) SetupTest() {
	var err error

	// Parse Dummy contract ABI
	suite.dummyABI, err = abi.JSON(strings.NewReader(Dummy.DummyABI))
	suite.Require().Nil(err)

	// Dummy contract (re)deployment
	suite.dummyContractAddress, _, _, err = Dummy.DeployDummy(suite.auth, suite.ledger.Client())
	suite.Require().Nil(err)
	suite.ledger.Commit()
}

func (suite *TxSimTestSuite) TestTxSim_GenericError() {
	// Message should be synchronized with the one returned from the Dummy contract
	wantContractError := txsim.ContractError{
		Sig:  "Error(string)",
		Args: []interface{}{"Hit generic error revert!"},
	}

	err := txsim.PopulateErrorMap("./testdata/")
	suite.Require().Nil(err)

	simulator := txsim.NewTransactionSimulator(suite.ledger.Client())
	suite.Require().Nil(err)

	tx := suite.generateHitGenericErrorTx()

	err = suite.ledger.Client().SendTransaction(context.TODO(), tx)
	suite.Require().Nil(err)
	suite.ledger.Commit()

	receipt, err := bind.WaitMined(context.TODO(), suite.ledger.Client(), tx)
	suite.Require().Nil(err)

	gotContractError, err := simulator.GetRevertReason(suite.T().Context(), receipt.TxHash)
	suite.Require().Nil(err)

	suite.Equal(wantContractError, gotContractError)
}

func (suite *TxSimTestSuite) TestTxSim_CustomError() {
	// Message should be synchronized with the one returned from the Dummy contract
	wantContractError := txsim.ContractError{
		Sig:  "CustomError(string)",
		Args: []interface{}{"Hit custom error revert!"},
	}

	err := txsim.PopulateErrorMap("./testdata/")
	suite.Require().Nil(err)

	simulator := txsim.NewTransactionSimulator(suite.ledger.Client())
	suite.Require().Nil(err)

	tx := suite.generateHitCustomErrorTx()

	err = suite.ledger.Client().SendTransaction(context.TODO(), tx)
	suite.Require().Nil(err)
	suite.ledger.Commit()

	receipt, err := bind.WaitMined(context.TODO(), suite.ledger.Client(), tx)
	suite.Require().Nil(err)

	gotContractError, err := simulator.GetRevertReason(suite.T().Context(), receipt.TxHash)
	suite.Require().Nil(err)

	suite.Equal(wantContractError, gotContractError)
}

func (suite *TxSimTestSuite) generateHitGenericErrorTx() *ethTypes.Transaction {
	calldata, err := suite.dummyABI.Pack("hitGenericError")
	suite.Require().NoError(err)

	return suite.newTransactionWithCalldata(calldata)
}

func (suite *TxSimTestSuite) generateHitCustomErrorTx() *ethTypes.Transaction {
	calldata, err := suite.dummyABI.Pack("hitCustomError")
	suite.Require().NoError(err)

	return suite.newTransactionWithCalldata(calldata)
}

func (suite *TxSimTestSuite) newTransactionWithCalldata(calldata []byte) *ethTypes.Transaction {
	nonce, err := suite.ledger.Client().PendingNonceAt(context.TODO(), suite.auth.From)
	suite.Require().Nil(err)

	gasPrice, err := suite.ledger.Client().SuggestGasPrice(context.TODO())
	suite.Require().Nil(err)

	tx := ethTypes.NewTransaction(
		nonce,
		suite.dummyContractAddress,
		big.NewInt(0),
		300000,
		gasPrice,
		calldata,
	)

	signedTx, err := ethTypes.SignTx(tx, ethTypes.NewEIP155Signer(suite.chainID), suite.privateKey)
	suite.Require().Nil(err)

	return signedTx
}

// this function executes after each test case
func (suite *TxSimTestSuite) TearDownTest() {
	// Nothing to clean up
}

// this function executes after all tests executed
func (suite *TxSimTestSuite) TearDownSuite() {
	// Nothing to clean up
}

func TestTxSimTestSuite(t *testing.T) {
	suite.Run(t, new(TxSimTestSuite))
}
