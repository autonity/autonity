package test

import (
	"context"
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
			conn, err := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(validator.rpcPort))
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()

			nonce, err := conn.PendingNonceAt(context.Background(), operatorAddress)
			if err != nil {
				t.Fatal(err)
			}

			gasPrice, err := conn.SuggestGasPrice(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			auth := bind.NewKeyedTransactor(operatorKey)
			auth.From = operatorAddress
			auth.Nonce = big.NewInt(int64(nonce))
			auth.GasLimit = uint64(300000) // in units
			auth.GasPrice = gasPrice

			contractAddress := autonity.ContractAddress
			instance, err := NewAutonity(contractAddress, conn)
			if err != nil {
				t.Fatal(err)
			}
			validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
			_, err = instance.MintStake(auth, *validatorsList[0].Address, stakeDelta)
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
			conn, err := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(validator.rpcPort))
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()

			nonce, err := conn.PendingNonceAt(context.Background(), operatorAddress)
			if err != nil {
				t.Fatal(err)
			}

			gasPrice, err := conn.SuggestGasPrice(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			auth := bind.NewKeyedTransactor(operatorKey)
			auth.From = operatorAddress
			auth.Nonce = big.NewInt(int64(nonce))
			auth.GasLimit = uint64(300000) // in units
			auth.GasPrice = gasPrice

			contractAddress := autonity.ContractAddress
			instance, err := NewAutonity(contractAddress, conn)
			if err != nil {
				t.Fatal(err)
			}
			validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
			_, err = instance.RedeemStake(auth, *validatorsList[0].Address, stakeDelta)
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
			conn, err := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(validator.rpcPort))
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()

			senderAddress := crypto.PubkeyToAddress(validator.privateKey.PublicKey)
			nonce, err := conn.PendingNonceAt(context.Background(), senderAddress)
			if err != nil {
				t.Fatal(err)
			}

			gasPrice, err := conn.SuggestGasPrice(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			auth := bind.NewKeyedTransactor(validator.privateKey)
			auth.From = senderAddress
			auth.Nonce = big.NewInt(int64(nonce))
			auth.GasLimit = uint64(30000000) // in units
			auth.GasPrice = gasPrice

			contractAddress := autonity.ContractAddress
			instance, err := NewAutonity(contractAddress, conn)
			if err != nil {
				t.Fatal(err)
			}
			validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
			toIndex := 0
			if senderAddress == *validatorsList[toIndex].Address {
				toIndex = 1
			}
			_, err = instance.Send(auth, *validatorsList[toIndex].Address, stakeDelta)
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
		conn, err := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(validators["VA"].rpcPort))
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		contractAddress := autonity.ContractAddress
		instance, err := NewAutonity(contractAddress, conn)
		if err != nil {
			t.Fatal(err)
		}

		auth := bind.CallOpts{
			Pending:     false,
			From:        common.Address{},
			BlockNumber: new(big.Int).SetUint64(3),
			Context:     context.Background(),
		}

		initNetworkMetrics, err := instance.DumpEconomicsMetricData(&auth)
		if err != nil {
			t.Fatal(err)
		}

		auth.BlockNumber.SetUint64(validators["VA"].lastBlock)
		curNetworkMetrics, err := instance.DumpEconomicsMetricData(&auth)
		if err != nil {
			t.Fatal(err)
		}
		validatorsList := validators["VA"].service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
		// check acount stake balance.
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
		conn, err := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(validators["VA"].rpcPort))
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		contractAddress := autonity.ContractAddress
		instance, err := NewAutonity(contractAddress, conn)
		if err != nil {
			t.Fatal(err)
		}

		auth := bind.CallOpts{
			Pending:     false,
			From:        common.Address{},
			BlockNumber: new(big.Int).SetUint64(3),
			Context:     context.Background(),
		}

		initNetworkMetrics, err := instance.DumpEconomicsMetricData(&auth)
		if err != nil {
			t.Fatal(err)
		}

		auth.BlockNumber.SetUint64(validators["VA"].lastBlock)
		curNetworkMetrics, err := instance.DumpEconomicsMetricData(&auth)
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
			numBlocks:     10,
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
			numBlocks:     10,
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
			numBlocks:     10,
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
