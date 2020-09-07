package test

import (
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/keygenerator"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"math/big"
	"testing"
	"time"
)

/*
  In this file, it create couple of test cases to check if the reward distribution works correctly in the autontiy network.
  THe work flow is base on the local e2e test framework's main flow.

  First it setup an autontiy network by according to the genesis hook function, then from the specific chain height, it
  start to issue transaction via the transaction hook function specified for the target node on each height, for example
  in the mintStakeHook, redeemStakeHook, transferStakeHook, it keep issuing transactions to call autonity contract to
  manage stake for members on each height.

  About the reward checking, in the beforeHooks set, we apply rewardChecker function to each member on the network on
  each height during the run time, it first dumps the balance view from parent block, and then it parse and calculate
  each TX's sentAmount, gasUsed, receiveAmount, fee. And also the expected reward fractions are calculated base on
  the stake portions from the parent block's view. Finally it checks if the balance is expected on the new block with a
  simple formula on each account: balanceOnParentBlock + reward + received - sent - gasUsed == balanceOnNewBlock.

  The TX issuing and reward checking work flow are keep running during each new block is mined.
*/

func TestRewardDistribution(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	// prepare chain operator
	operatorKey, err := keygenerator.Next()
	if err != nil {
		t.Fatal(err)
	}
	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)
	// reward checker hook:
	rewardChecker := func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
		type economicMetaData struct {
			Accounts        []common.Address `abi:"accounts"`
			Usertypes       []uint8          `abi:"usertypes"`
			Stakes          []*big.Int       `abi:"stakes"`
			Commissionrates []*big.Int       `abi:"commissionrates"`
			Mingasprice     *big.Int         `abi:"mingasprice"`
			Stakesupply     *big.Int         `abi:"stakesupply"`
		}
		// reward distribution checking start from height-2.
		if block.NumberU64() < 2 {
			return nil
		}
		blockReward := new(big.Int)
		gp := new(core.GasPool).AddGas(block.GasLimit())
		parent := validator.service.BlockChain().GetBlock(block.ParentHash(), block.NumberU64()-1)

		addr := crypto.PubkeyToAddress(validator.privateKey.PublicKey)
		// get parent block view.
		parentState, stateErr := validator.service.BlockChain().StateAt(parent.Root())
		if stateErr != nil {
			return stateErr
		}
		balanceBase := parentState.GetBalance(addr)
		// calculate sent, gasUsed, receive amount and reward base on new block.
		sentAmount := new(big.Int)
		receivedAmount := new(big.Int)
		usedGas := new(big.Int)
		rewardFraction := new(big.Int)
		for i, tx := range block.Transactions() {
			// Apply TX to parent view again to estimate the gasUsed of the TX.
			parentState.Prepare(tx.Hash(), block.Hash(), i)
			receipt, receiptErr := core.ApplyTransaction(validator.service.BlockChain().Config(), validator.service.BlockChain(), nil, gp, parentState, block.Header(), tx, new(uint64), *validator.service.BlockChain().GetVMConfig())
			if receiptErr != nil {
				return receiptErr
			}

			// count block reward.
			txGasConsumed := new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(receipt.GasUsed))
			blockReward.Add(blockReward, txGasConsumed)

			// count received amount on the block.
			if *tx.To() == addr {
				receivedAmount.Add(receivedAmount, tx.Value())
			}

			// count sent and gasUsed amount on the block.
			validator.transactionsMu.Lock()
			if _, ok := validator.transactions[tx.Hash()]; ok {
				sentAmount.Add(sentAmount, tx.Value())
				usedGas.Add(usedGas, txGasConsumed)
			}
			validator.transactionsMu.Unlock()
		}

		// calculate reward fractions.
		parentEcomonimcState := economicMetaData{}
		err := validator.service.BlockChain().GetAutonityContract().AutonityContractCall(parentState, parent.Header(), "dumpEconomicsMetricData", &parentEcomonimcState)
		if err != nil {
			return err
		}

		for i := 0; i < len(parentEcomonimcState.Accounts); i++ {
			if addr == parentEcomonimcState.Accounts[i] {
				rewardFraction = new(big.Int).Mul(blockReward, parentEcomonimcState.Stakes[i])
				rewardFraction = rewardFraction.Div(rewardFraction, parentEcomonimcState.Stakesupply)
				break
			}
		}

		// current block view.
		currentState, stateErr := validator.service.BlockChain().StateAt(block.Root())
		if stateErr != nil {
			return stateErr
		}
		balanceActual := currentState.GetBalance(addr)
		// check if balance is expected: base + reward + received - sent - gasUsed == bActual.
		balanceWant := new(big.Int).Add(balanceBase, rewardFraction)
		balanceWant.Add(balanceWant, receivedAmount)
		balanceWant.Sub(balanceWant, sentAmount)
		balanceWant.Sub(balanceWant, usedGas)
		if balanceWant.Cmp(balanceActual) != 0 {
			return fmt.Errorf("incorrect reward distribution")
		}
		return nil
	}
	// mint stake hook
	mintStakeHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock <= 3 {
			return true, nil, nil
		}

		validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
		index := validator.lastBlock % uint64(len(validatorsList))
		return true, nil, interact(validator.rpcPort).tx(operatorKey).mintStake(*validatorsList[index].Address, new(big.Int).SetUint64(100))
	}
	// send stake hook
	transferStakeHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock <= 3 {
			return true, nil, nil
		}

		validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
		to := validator.lastBlock % uint64(len(validatorsList))
		return true, nil, interact(validator.rpcPort).tx(validator.privateKey).sendStake(*validatorsList[to].Address, new(big.Int).SetUint64(1))
	}
	// redeem stake hook
	redeemStakeHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock <= 3 {
			return true, nil, nil
		}

		validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
		from := validator.lastBlock % uint64(len(validatorsList))
		return true, nil, interact(validator.rpcPort).tx(operatorKey).redeemStake(*validatorsList[from].Address, new(big.Int).SetUint64(1))
	}
	// genesis hook
	genesisHook := func(g *core.Genesis) *core.Genesis {
		g.Config.AutonityContractConfig.Operator = operatorAddress
		g.Alloc[operatorAddress] = core.GenesisAccount{
			Balance: big.NewInt(100000000000000000),
		}
		return g
	}

	testCases := []*testCase{
		{
			name:                 "reward distribution check with fixed staking",
			numValidators:        6,
			numBlocks:            30,
			txPerPeer:            1,
			sendTransactionHooks: make(map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error)),
			// Apply reward checker to all nodes:
			beforeHooks: map[string]hook{
				"VA": rewardChecker,
				"VB": rewardChecker,
				"VC": rewardChecker,
				"VD": rewardChecker,
				"VE": rewardChecker,
				"VF": rewardChecker,
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "reward distribution check with run time stake mint",
			numValidators: 6,
			numBlocks:     30,
			txPerPeer:     1,
			sendTransactionHooks: map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error){
				"VA": mintStakeHook,
			},
			genesisHook: genesisHook,
			// All nodes in the network check its reward distribution.
			beforeHooks: map[string]hook{
				"VA": rewardChecker,
				"VB": rewardChecker,
				"VC": rewardChecker,
				"VD": rewardChecker,
				"VE": rewardChecker,
				"VF": rewardChecker,
			},
		},
		{
			name:          "reward distribution check with run time stake transfer",
			numValidators: 6,
			numBlocks:     30,
			txPerPeer:     1,
			sendTransactionHooks: map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error){
				"VA": transferStakeHook,
			},
			genesisHook: genesisHook,
			// All nodes in the network check its reward distribution.
			beforeHooks: map[string]hook{
				"VA": rewardChecker,
				"VB": rewardChecker,
				"VC": rewardChecker,
				"VD": rewardChecker,
				"VE": rewardChecker,
				"VF": rewardChecker,
			},
		},
		{
			name:          "reward distribution check with run time stake redeem",
			numValidators: 6,
			numBlocks:     30,
			txPerPeer:     1,
			sendTransactionHooks: map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error){
				"VA": redeemStakeHook,
			},
			genesisHook: genesisHook,
			// All nodes in the network check its reward distribution.
			beforeHooks: map[string]hook{
				"VA": rewardChecker,
				"VB": rewardChecker,
				"VC": rewardChecker,
				"VD": rewardChecker,
				"VE": rewardChecker,
				"VF": rewardChecker,
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}
