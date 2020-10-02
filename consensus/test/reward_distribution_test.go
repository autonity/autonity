package test

import (
	"context"
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/ethclient"
	"github.com/stretchr/testify/require"
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
	operatorKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)

	// get balance on specific block for an node
	getBalance := func(blockNum *big.Int, node *testNode) (*big.Int, error) {
		client, err := node.node.Attach()
		if err != nil {
			return nil, err
		}
		defer client.Close()
		ec := ethclient.NewClient(client)
		ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
		defer cancel()

		addr := crypto.PubkeyToAddress(node.privateKey.PublicKey)
		balance, err := ec.BalanceAt(ctx, addr, blockNum)
		if err != nil {
			return nil, err
		}
		return balance, nil
	}
	// calculate sentTokens, receivedTokens, gasUsed, blockReward on per block for node
	calculateBalanceFactors := func(block *types.Block, node *testNode) (*big.Int, *big.Int, *big.Int, *big.Int, error) {
		client, err := node.node.Attach()
		if err != nil {
			return nil, nil, nil, nil, err
		}
		defer client.Close()
		ec := ethclient.NewClient(client)
		ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
		defer cancel()
		// calculate sent, gasUsed, receive amount and reward base on new block.
		sentAmount := new(big.Int)
		receivedAmount := new(big.Int)
		usedGas := new(big.Int)
		blockReward := new(big.Int)
		for _, tx := range block.Transactions() {
			receipt, err := ec.TransactionReceipt(ctx, tx.Hash())
			if err != nil {
				return nil, nil, nil, nil, err
			}

			// count block reward.
			txGasConsumed := new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(receipt.GasUsed))
			blockReward.Add(blockReward, txGasConsumed)

			// count received amount on the block.
			if *tx.To() == crypto.PubkeyToAddress(node.privateKey.PublicKey) {
				receivedAmount.Add(receivedAmount, tx.Value())
			}

			// count sent and gasUsed amount on the block.
			node.transactionsMu.Lock()
			if _, ok := node.transactions[tx.Hash()]; ok {
				sentAmount.Add(sentAmount, tx.Value())
				usedGas.Add(usedGas, txGasConsumed)
			}
			node.transactionsMu.Unlock()
		}
		return sentAmount, receivedAmount, usedGas, blockReward, nil
	}
	// calculate reward fraction
	calculateRewardFraction := func(blockNum uint64, blockReward *big.Int, port int, address common.Address) (*big.Int, error) {
		rewardFraction := new(big.Int)
		// calculate reward fractions base on latest economic state
		economicState, err := interact(port).call(blockNum).dumpEconomicsMetricData()
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(economicState.Accounts); i++ {
			if address == economicState.Accounts[i] {
				rewardFraction = new(big.Int).Mul(blockReward, economicState.Stakes[i])
				rewardFraction = rewardFraction.Div(rewardFraction, economicState.Stakesupply)
				break
			}
		}
		return rewardFraction, nil
	}

	// reward checker hook:
	rewardChecker := func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {

		// get base balance.
		parentNum := block.Number().Sub(block.Number(), big.NewInt(1))
		balanceOnParent, err := getBalance(parentNum, validator)
		if err != nil {
			return err
		}

		// calculate sent, gasUsed, receive amount and reward base on new block.
		sentAmount, receivedAmount, usedGas, blockReward, err := calculateBalanceFactors(block, validator)
		if err != nil {
			return err
		}

		// calculate reward fraction.
		rewardFraction, err := calculateRewardFraction(block.NumberU64(), blockReward, validator.rpcPort, crypto.PubkeyToAddress(validator.privateKey.PublicKey))
		if err != nil {
			return err
		}

		balanceNow, err := getBalance(block.Number(), validator)
		if err != nil {
			return err
		}

		// check if balance is expected: balanceOnParent + rewardFraction + received - sent - gasUsed == balanceNow.
		balanceWant := new(big.Int).Add(balanceOnParent, rewardFraction)
		balanceWant.Add(balanceWant, receivedAmount)
		balanceWant.Sub(balanceWant, sentAmount)
		balanceWant.Sub(balanceWant, usedGas)
		if balanceWant.Cmp(balanceNow) != 0 {
			return fmt.Errorf("incorrect reward distribution")
		}
		return nil
	}
	// mint stake hook
	mintStakeHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
		index := validator.lastBlock % uint64(len(validatorsList))
		return true, nil, interact(validator.rpcPort).tx(operatorKey).mintStake(*validatorsList[index].Address, new(big.Int).SetUint64(100))
	}
	// send stake hook
	transferStakeHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
		to := validator.lastBlock % uint64(len(validatorsList))
		return true, nil, interact(validator.rpcPort).tx(validator.privateKey).sendStake(*validatorsList[to].Address, new(big.Int).SetUint64(1))
	}
	// redeem stake hook
	redeemStakeHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
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
			sendTransactionHooks: make(map[string]sendTransactionHook),
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
			sendTransactionHooks: map[string]sendTransactionHook{
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
			sendTransactionHooks: map[string]sendTransactionHook{
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
			sendTransactionHooks: map[string]sendTransactionHook{
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
