package test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/ethclient"
	"github.com/stretchr/testify/require"
)

/*
  This test checks that at each block rewards are distributed correctly.

  We check four different scenarios.

  * Validators' stake remains constant
  * A validator has stake minted (added) each block.
  * A validator has stake redeemed (removed) each block.
  * A validator sends stake to another validator each block.

  The reward distribution model that we are following is defined here - https://core.clearmatics.net/ASC_core_v3.html#with-no-proposer-bonus-fees.

  It can be summarised as each validator is rewarded with a portion of the total fees for a block that is proportional
  to the amount of stake they own. The total fees for a block is the sum of the fees of the transactions where a transaction
  fee is the gas used multiplied by the gas price.

  For each scenario after each block we check that the balance of each validator is as expected, taking into account its
  stake at that block, any value that has been sent or received at that block and any gas used in that block.
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

	// calculate reward fraction calculates the fraction of a block reward that
	// should be assigned to the node having the given address.
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

	// calculate sentTokens, receivedTokens, gasUsed and blockReward for the given node in the given block.
	calculateRewardMetaPerBlock := func(block *types.Block, node *testNode) (*big.Int, *big.Int, *big.Int, *big.Int, error) {
		client, err := node.node.Attach()
		if err != nil {
			return nil, nil, nil, nil, err
		}
		defer client.Close()
		ec := ethclient.NewClient(client)
		ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
		defer cancel()
		// calculate sent, gasUsed, receive amount and reward on the block.
		sentAmount := new(big.Int)
		receivedAmount := new(big.Int)
		usedGas := new(big.Int)
		blockReward := new(big.Int)
		for _, tx := range block.Transactions() {
			receipt, err := ec.TransactionReceipt(ctx, tx.Hash())
			if err != nil {
				return nil, nil, nil, nil, err
			}

			// count reward on this block.
			txGasConsumed := new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(receipt.GasUsed))
			blockReward.Add(blockReward, txGasConsumed)

			// count received amount on this block.
			if *tx.To() == crypto.PubkeyToAddress(node.privateKey.PublicKey) {
				receivedAmount.Add(receivedAmount, tx.Value())
			}

			// count sent and gasUsed amount on this block.
			node.transactionsMu.Lock()
			if _, ok := node.transactions[tx.Hash()]; ok {
				sentAmount.Add(sentAmount, tx.Value())
				usedGas.Add(usedGas, txGasConsumed)
			}
			node.transactionsMu.Unlock()
		}
		fraction, err := calculateRewardFraction(block.NumberU64(), blockReward, node.rpcPort, crypto.PubkeyToAddress(node.privateKey.PublicKey))
		return sentAmount, receivedAmount, usedGas, fraction, err
	}
	// reward checker hook:
	rewardChecker := func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {

		// get base balance.
		parentNum := block.Number().Sub(block.Number(), big.NewInt(1))
		balanceOnParent, err := getBalance(parentNum, validator)
		if err != nil {
			return err
		}

		// calculate sent, gasUsed, receive amount and reward for the validator in the new block.
		sentAmount, receivedAmount, usedGas, blockReward, err := calculateRewardMetaPerBlock(block, validator)
		if err != nil {
			return err
		}

		balanceNow, err := getBalance(block.Number(), validator)
		if err != nil {
			return err
		}

		// check if balance is expected: balanceOnParent + rewardFraction + received - sent - gasUsed == balanceNow.
		balanceWant := new(big.Int).Add(balanceOnParent, blockReward)
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
		validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidators()
		index := validator.lastBlock % uint64(len(validatorsList))
		return true, nil, interact(validator.rpcPort).tx(operatorKey).mintStake(*validatorsList[index].Address, new(big.Int).SetUint64(100))
	}
	// send stake hook
	transferStakeHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidators()
		to := validator.lastBlock % uint64(len(validatorsList))
		return true, nil, interact(validator.rpcPort).tx(validator.privateKey).sendStake(*validatorsList[to].Address, new(big.Int).SetUint64(1))
	}
	// redeem stake hook
	redeemStakeHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidators()
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
