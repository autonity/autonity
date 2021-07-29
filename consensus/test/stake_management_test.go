package test

import (
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
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

func TestStakeManagement(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	initHeight := uint64(0)
	startHeight := uint64(1)
	mintStake := new(big.Int).SetUint64(50)
	redeemStake := new(big.Int).SetUint64(100)
	sendStake := new(big.Int).SetUint64(2)
	// prepare chain operator
	operatorKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)

	// genesis hook
	genesisHook := func(g *core.Genesis) *core.Genesis {
		g.Config.AutonityContractConfig.Operator = operatorAddress
		g.Alloc[operatorAddress] = core.GenesisAccount{
			Balance: big.NewInt(100000000000000000),
		}
		return g
	}

	// mint stake hook
	mintStakeHook := func(validator *testNode, address common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock == startHeight {
			return true, nil, interact(validator.rpcPort).tx(operatorKey).mintStake(address, mintStake)
		}
		return false, nil, nil
	}

	stakeChecker := func(t *testing.T, validators map[string]*testNode, stake *big.Int) {
		address := crypto.PubkeyToAddress(validators["VA"].privateKey.PublicKey)
		port := validators["VA"].rpcPort

		initBalance, err := interact(port).call(initHeight).getAccountStake(address)
		require.NoError(t, err)

		newBalance, err := interact(port).call(validators["VA"].lastBlock).getAccountStake(address)
		require.NoError(t, err)

		delta := newBalance.Sub(newBalance, initBalance)
		assert.Equal(t, delta.Int64(), stake.Int64(), "stake balance is not expected")

		initNetworkMetrics, err := interact(validators["VA"].rpcPort).call(initHeight).dumpEconomicsMetricData()
		require.NoError(t, err)

		curNetworkMetrics, err := interact(validators["VA"].rpcPort).call(validators["VA"].lastBlock).dumpEconomicsMetricData()
		require.NoError(t, err)

		totalDelta := curNetworkMetrics.Stakesupply.Sub(curNetworkMetrics.Stakesupply, initNetworkMetrics.Stakesupply)
		assert.Equal(t, totalDelta.Int64(), stake.Int64(), "stake total supply is not expected")
	}

	// mint stake checker hook
	mintStakeCheckerHook := func(t *testing.T, validators map[string]*testNode) {
		stakeChecker(t, validators, mintStake)
	}

	redeemStakeHook := func(validator *testNode, address common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock == startHeight {
			return true, nil, interact(validator.rpcPort).tx(operatorKey).redeemStake(address, redeemStake)
		}
		return false, nil, nil
	}

	redeemStakeCheckerHook := func(t *testing.T, validators map[string]*testNode) {
		stakeChecker(t, validators, redeemStake.Neg(redeemStake))
	}

	pickReceiver := func(validator *testNode) common.Address {
		senderAddress := crypto.PubkeyToAddress(validator.privateKey.PublicKey)
		validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
		toIndex := 0
		if senderAddress == *validatorsList[toIndex].Address {
			toIndex = 1
		}
		return *validatorsList[toIndex].Address
	}

	sendStakeHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock == startHeight {
			receiver := pickReceiver(validator)
			return true, nil, interact(validator.rpcPort).tx(validator.privateKey).sendStake(receiver, sendStake)
		}
		return false, nil, nil
	}

	stakeSendCheckerHook := func(t *testing.T, validators map[string]*testNode) {
		port := validators["VA"].rpcPort

		initNetworkMetrics, err := interact(port).call(initHeight).dumpEconomicsMetricData()
		require.NoError(t, err)

		curNetworkMetrics, err := interact(port).call(validators["VA"].lastBlock).dumpEconomicsMetricData()
		require.NoError(t, err)

		senderAddr := crypto.PubkeyToAddress(validators["VA"].privateKey.PublicKey)
		receiverAddr := pickReceiver(validators["VA"])

		senderInitBalance, err := interact(port).call(initHeight).getAccountStake(senderAddr)
		require.NoError(t, err)
		senderNewBalance, err := interact(port).call(validators["VA"].lastBlock).getAccountStake(senderAddr)
		require.NoError(t, err)
		delta := senderInitBalance.Uint64() - senderNewBalance.Uint64()
		require.Equal(t, delta, sendStake.Uint64())

		receiverInitBalance, err := interact(port).call(initHeight).getAccountStake(receiverAddr)
		require.NoError(t, err)
		receiverNewBalance, err := interact(port).call(validators["VA"].lastBlock).getAccountStake(receiverAddr)
		require.NoError(t, err)
		delta = receiverNewBalance.Uint64() - receiverInitBalance.Uint64()
		require.Equal(t, delta, sendStake.Uint64())

		assert.Equal(t, initNetworkMetrics.Stakesupply.Uint64(), curNetworkMetrics.Stakesupply.Uint64(), "total stake supply is not expected")
	}

	// numBlocks are used to stop the test on current test framework, to let stake management TX to be mined before the test end,
	// bigger numBlocks in below test cases are set.
	testCases := []*testCase{
		{
			name:          "stake management test mint stake",
			numValidators: 6,
			numBlocks:     100,
			txPerPeer:     1,
			sendTransactionHooks: map[string]sendTransactionHook{
				"VA": mintStakeHook,
			},
			genesisHook: genesisHook,
			finalAssert: mintStakeCheckerHook,
		},
		{
			name:          "stake management test redeem stake",
			numValidators: 6,
			numBlocks:     20,
			txPerPeer:     1,
			sendTransactionHooks: map[string]sendTransactionHook{
				"VA": redeemStakeHook,
			},
			genesisHook: genesisHook,
			finalAssert: redeemStakeCheckerHook,
		},
		{
			name:          "stake management test send stake",
			numValidators: 6,
			numBlocks:     20,
			txPerPeer:     1,
			sendTransactionHooks: map[string]sendTransactionHook{
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
