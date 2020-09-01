package test

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/clearmatics/autonity/accounts/abi/bind"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/keygenerator"
	"github.com/clearmatics/autonity/contracts/autonity"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/ethclient"
	"math/big"
	"strconv"
	"sync"
	"testing"
)

/*
  In this file, it create 3 test cases which have similar work flow base on the local e2e test framework's main flow.

  First it setup an autontiy network by according to the genesis hook function, then from the specific chain height, it
  start to issue transactions via the transaction hook function specified for the target node, for example in the
  mintStakeHook, redeemStakeHook, and sendStakeHook, it issues transaction to call autonity contract via operator account
  to manage the stakes on the members.

  Then the test case verify the output from its finalAssert hook function on the specified height of the blockchain, for
  example, it checks the stake balance in different height to compare if the balance is expected.
*/

type testAutonity struct {
	autonity *Autonity
	transactionOpt *bind.TransactOpts
	callOpt *bind.CallOpts
	client *ethclient.Client
}

func (a *testAutonity) Close() {
	if a.client != nil {
		a.client.Close()
	}
}

func autonityInstance(operatorKey *ecdsa.PrivateKey, node *testNode) (*testAutonity, error) {

	contract := new(testAutonity)

	conn, err := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(node.rpcPort))
	contract.client = conn
	if err != nil {
		return contract, err
	}

	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)
	nonce, err := conn.PendingNonceAt(context.Background(), operatorAddress)
	if err != nil {
		return contract, err
	}

	gasPrice, err := conn.SuggestGasPrice(context.Background())
	if err != nil {
		return contract, err
	}

	txOpt := bind.NewKeyedTransactor(operatorKey)
	txOpt.From = operatorAddress
	txOpt.Nonce = big.NewInt(int64(nonce))
	txOpt.GasLimit = uint64(300000000) // in units
	txOpt.GasPrice = gasPrice
	instance, err := NewAutonity(autonity.ContractAddress, conn)
	if err != nil {
		return contract, err
	}

	callOpt := &bind.CallOpts{
		Pending:     false,
		From:        common.Address{},
		BlockNumber: new(big.Int),
		Context:     context.Background(),
	}

	contract.autonity = instance
	contract.client = conn
	contract.transactionOpt = txOpt
	contract.callOpt = callOpt
	return contract, nil
}

func TestStakeManagement(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	onceMint := sync.Once{}
	onceRedeem := sync.Once{}
	onceSend := sync.Once{}
	stakeDelta := new(big.Int).SetUint64(50)
	// prepare chain operator
	operatorKey, err := keygenerator.Next()
	if err != nil {
		t.Fatal(err)
	}
	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)
	// mint stake hook
	mintStakeHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock <= 3 {
			return true, nil, nil
		}
		onceMint.Do(func() {

			contract, err := autonityInstance(operatorKey, validator)
			defer contract.Close()

			if err != nil {
				t.Fatal(err)
			}

			validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
			_, err = contract.autonity.MintStake(contract.transactionOpt, *validatorsList[0].Address, stakeDelta)
			if err != nil {
				t.Fatal(err)
			}
		})
		return false, nil, nil
	}

	redeemStakeHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock <= 3 {
			return true, nil, nil
		}
		onceRedeem.Do(func() {

			contract, err := autonityInstance(operatorKey, validator)
			defer contract.Close()

			if err != nil {
				t.Fatal(err)
			}

			validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
			_, err = contract.autonity.RedeemStake(contract.transactionOpt, *validatorsList[0].Address, stakeDelta)
			if err != nil {
				t.Fatal(err)
			}
		})
		return false, nil, nil
	}

	sendStakeHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock <= 3 {
			return true, nil, nil
		}
		onceSend.Do(func() {

			contract, err := autonityInstance(validator.privateKey, validator)
			defer contract.Close()
			if err != nil {
				t.Fatal(err)
			}

			senderAddress := crypto.PubkeyToAddress(validator.privateKey.PublicKey)
			validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
			toIndex := 0
			if senderAddress == *validatorsList[toIndex].Address {
				toIndex = 1
			}
			_, err = contract.autonity.Send(contract.transactionOpt, *validatorsList[toIndex].Address, stakeDelta)
			if err != nil {
				t.Fatal(err)
			}
		})
		return false, nil, nil
	}

	// genesis hook
	genesisHook := func(g *core.Genesis) *core.Genesis {
		g.Config.AutonityContractConfig.Operator = operatorAddress
		g.Alloc[operatorAddress] = core.GenesisAccount{
			Balance: big.NewInt(100000000000000000),
		}
		return g
	}

	stakeCheckerHook := func(t *testing.T, validators map[string]*testNode) {
		contract, err := autonityInstance(operatorKey, validators["VA"])
		defer contract.Close()

		if err != nil {
			t.Fatal(err)
		}

		contract.callOpt.BlockNumber.SetUint64(3)
		initNetworkMetrics, err := contract.autonity.DumpEconomicsMetricData(contract.callOpt)
		if err != nil {
			t.Fatal(err)
		}

		contract.callOpt.BlockNumber.SetUint64(validators["VA"].lastBlock)
		curNetworkMetrics, err := contract.autonity.DumpEconomicsMetricData(contract.callOpt)
		if err != nil {
			t.Fatal(err)
		}
		validatorsList := validators["VA"].service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
		// check account stake balance.
		founded := false
		for index, v := range initNetworkMetrics.Accounts {
			if v == *validatorsList[0].Address {
				founded = true
				initBalance := initNetworkMetrics.Stakes[index]
				newBalance := curNetworkMetrics.Stakes[index]
				delta := new(big.Int).Abs(initBalance.Sub(initBalance, newBalance))
				if delta.Uint64() != stakeDelta.Uint64() {
					t.Fatal("stake balance is not expected")
				}
				totalDelta := new(big.Int).Abs(initNetworkMetrics.Stakesupply.Sub(initNetworkMetrics.Stakesupply, curNetworkMetrics.Stakesupply))
				if totalDelta.Uint64() != stakeDelta.Uint64() {
					t.Fatal("stake total supply is not expected")
				}
			}
		}
		if !founded {
			t.Fatal("cannot find wanted account from chain DB")
		}
	}

	stakeSendCheckerHook := func(t *testing.T, validators map[string]*testNode) {
		contract, err := autonityInstance(operatorKey, validators["VA"])
		defer contract.Close()

		if err != nil {
			t.Fatal(err)
		}

		contract.callOpt.BlockNumber.SetUint64(3)
		initNetworkMetrics, err := contract.autonity.DumpEconomicsMetricData(contract.callOpt)
		if err != nil {
			t.Fatal(err)
		}

		contract.callOpt.BlockNumber.SetUint64(validators["VA"].lastBlock)
		curNetworkMetrics, err := contract.autonity.DumpEconomicsMetricData(contract.callOpt)
		if err != nil {
			t.Fatal(err)
		}

		validatorsList := validators["VA"].service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
		senderAddr := crypto.PubkeyToAddress(validators["VA"].privateKey.PublicKey)
		toIndex := 0
		if senderAddr == *validatorsList[toIndex].Address {
			toIndex = 1
		}
		receiverAddr := *validatorsList[toIndex].Address

		// check account stake balance.
		senderPassed := false
		receiverPassed := false
		for index, v := range curNetworkMetrics.Accounts {
			if v == senderAddr {
				curBalance := curNetworkMetrics.Stakes[index].Uint64()
				initBalance := initNetworkMetrics.Stakes[index].Uint64()
				delta := initBalance - curBalance
				if delta == stakeDelta.Uint64() {
					senderPassed = true
					continue
				}
			}
			if v == receiverAddr {
				curBalance := curNetworkMetrics.Stakes[index].Uint64()
				initBalance := initNetworkMetrics.Stakes[index].Uint64()
				delta := curBalance - initBalance
				if delta == stakeDelta.Uint64() {
					receiverPassed = true
					continue
				}
			}
		}

		if !senderPassed || !receiverPassed {
			t.Fatal("stake balance checking failed")
		}

		if initNetworkMetrics.Stakesupply.Uint64() != curNetworkMetrics.Stakesupply.Uint64() {
			t.Fatal("total stake supply is not expected")
		}
	}

	testCases := []*testCase{
		{
			name:          "stake management test mint stake",
			numValidators: 6,
			numBlocks:     20,
			txPerPeer:     1,
			sendTransactionHooks: map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error){
				"VA": mintStakeHook,
			},
			genesisHook: genesisHook,
			finalAssert: stakeCheckerHook,
		},
		{
			name:          "stake management test redeem stake",
			numValidators: 6,
			numBlocks:     20,
			txPerPeer:     1,
			sendTransactionHooks: map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error){
				"VA": redeemStakeHook,
			},
			genesisHook: genesisHook,
			finalAssert: stakeCheckerHook,
		},
		{
			name:          "stake management test send stake",
			numValidators: 6,
			numBlocks:     20,
			txPerPeer:     1,
			sendTransactionHooks: map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error){
				"VA": sendStakeHook,
			},
			genesisHook: genesisHook,
			finalAssert: stakeSendCheckerHook,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}
