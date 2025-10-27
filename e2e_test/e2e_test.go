package e2e

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"text/tabwriter"
	"time"

	bindings1 "github.com/autonity/autonity/autonity/bindings/1"
	generated1 "github.com/autonity/autonity/params/upgrades/generated/1"
	generatedtests "github.com/autonity/autonity/params/upgrades/generated/tests"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/autonity/bindings"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/backend"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	ccore "github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/rlp"
)

// This test checks that we can process transactions that transfer value from
// one participant to another.
func TestSendingValue(t *testing.T) {
	network, err := NewNetwork(t, 7, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown(t)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	err = network[0].SendAUTtracked(ctx, network[1].Address, 10)
	require.NoError(t, err)
	_ = network.WaitToMineNBlocks(100, 100, false)
}

func TestProtocolContractsDeployment(t *testing.T) {
	network, err := NewNetwork(t, 2, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown(t)
	// Autonity Contract
	autonityContract, _ := bindings.NewAutonity(params.AutonityContractAddress, network[0].WsClient)
	autonityConfig, err := autonityContract.GetConfig(nil)
	require.NoError(t, err)

	validators, err := autonityContract.GetValidators(nil)
	require.NoError(t, err)

	// gengen sets the operator to the address of the first validator
	require.Equal(t, validators[0], autonityConfig.Protocol.OperatorAccount)
	require.Equal(t, params.TestAutonityContractConfig.BlockPeriod, autonityConfig.Protocol.BlockPeriod.Uint64())
	require.Equal(t, params.TestAutonityContractConfig.EpochPeriod, autonityConfig.Protocol.EpochPeriod.Uint64())
	require.Equal(t, params.TestAutonityContractConfig.GasLimit, autonityConfig.Protocol.GasLimit.Uint64())
	require.Equal(t, params.TestAutonityContractConfig.GasLimitBoundDivisor, autonityConfig.Protocol.GasLimitBoundDivisor.Uint64())
	require.Equal(t, params.TestAutonityContractConfig.BaseFeeChangeDenominator, autonityConfig.Policy.BaseFeeChangeDenominator.Uint64())
	require.Equal(t, params.TestAutonityContractConfig.ElasticityMultiplier, autonityConfig.Policy.ElasticityMultiplier.Uint64())
	require.Equal(t, params.TestAutonityContractConfig.MaxCommitteeSize, autonityConfig.Protocol.CommitteeSize.Uint64())
	require.Equal(t, params.TestAutonityContractConfig.DelegationRate, autonityConfig.Policy.DelegationRate.Uint64())
	require.Equal(t, params.TestAutonityContractConfig.MinBaseFee, autonityConfig.Policy.MinBaseFee.Uint64())
	require.Equal(t, params.TestAutonityContractConfig.Treasury, autonityConfig.Policy.TreasuryAccount)
	require.Equal(t, params.TestAutonityContractConfig.WithheldRewardsPool, autonityConfig.Policy.WithheldRewardsPool)
	require.Equal(t, params.TestAutonityContractConfig.TreasuryFee, autonityConfig.Policy.TreasuryFee.Uint64())
	require.Equal(t, params.TestAutonityContractConfig.UnbondingPeriod, autonityConfig.Policy.UnbondingPeriod.Uint64())
	require.Equal(t, params.UpgradeManagerContractAddress, autonityConfig.Contracts.UpgradeManagerContract)
	require.Equal(t, params.AccountabilityContractAddress, autonityConfig.Contracts.AccountabilityContract)
	require.Equal(t, params.OmissionAccountabilityContractAddress, autonityConfig.Contracts.OmissionAccountabilityContract)
	require.Equal(t, params.ACUContractAddress, autonityConfig.Contracts.AcuContract)
	require.Equal(t, params.StabilizationContractAddress, autonityConfig.Contracts.StabilizationContract)
	require.Equal(t, params.OracleContractAddress, autonityConfig.Contracts.OracleContract)
	require.Equal(t, params.SupplyControlContractAddress, autonityConfig.Contracts.SupplyControlContract)
	// Accountability Contract
	accountabilityContract, _ := bindings.NewAccountability(params.AccountabilityContractAddress, network[0].WsClient)
	accountabilityConfig, err := accountabilityContract.GetConfig(nil)
	require.NoError(t, err)
	require.Equal(t, params.TestAccountabilityConfig.HistoryFactor, accountabilityConfig.Factors.History.Uint64())
	require.Equal(t, params.TestAccountabilityConfig.CollusionFactor, accountabilityConfig.Factors.Collusion.Uint64())
	require.Equal(t, params.TestAccountabilityConfig.JailFactor, accountabilityConfig.Factors.Jail.Uint64())
	require.Equal(t, params.TestAccountabilityConfig.BaseSlashingRateLow, accountabilityConfig.BaseSlashingRates.Low.Uint64())
	require.Equal(t, params.TestAccountabilityConfig.BaseSlashingRateMid, accountabilityConfig.BaseSlashingRates.Mid.Uint64())
	require.Equal(t, params.TestAccountabilityConfig.BaseSlashingRateHigh, accountabilityConfig.BaseSlashingRates.High.Uint64())
	require.Equal(t, params.TestAccountabilityConfig.InnocenceProofSubmissionWindow, accountabilityConfig.InnocenceProofSubmissionWindow.Uint64())
	// Oracle Contract -- todo
	// ACU Contract -- todo
	// Supply Control Contract -- todo
	// Stabilization Contract -- todo
	// Omission Contract -- todo
	// Upgrade Manager Contract
	upgradeManagerContract, _ := bindings.NewUpgradeManager(params.UpgradeManagerContractAddress, network[0].WsClient)
	upgradeManagerAutonityAddress, err := upgradeManagerContract.GetAutonity(nil)
	require.NoError(t, err)
	require.Equal(t, params.AutonityContractAddress, upgradeManagerAutonityAddress)
	err = network.WaitToMineNBlocks(2, 15, false)
	require.NoError(t, err)
}

// tests if upgrading all ASM contracts at once is feasible in terms of gas and tx size
func TestAsmAtomicUpgrade(t *testing.T) {
	network, err := NewNetwork(t, 2, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown(t)

	// verify that upgrade manager version 1.1.0 has been deployed at genesis
	upgradeManagerCode, err := network[0].WsClient.CodeAt(context.Background(), params.UpgradeManagerContractAddress, nil)
	require.NoError(t, err)
	require.True(t, bytes.Equal(upgradeManagerCode, generated1.UpgradeManagerRuntimeBytecode))

	// build ASM contracts upgrade tx
	upgradeManager, err := bindings1.NewUpgradeManager1(params.UpgradeManagerContractAddress, network[0].WsClient)
	require.NoError(t, err)

	operatorKey := network[0].Key
	transactOpts, err := bind.NewKeyedTransactorWithChainID(operatorKey, params.TestChainConfig.ChainID)
	require.NoError(t, err)
	tx, err := upgradeManager.UpgradeMultiple(transactOpts,
		[]common.Address{
			params.ACUContractAddress,
			params.SupplyControlContractAddress,
			params.StabilizationContractAddress,
			params.InflationControllerContractAddress,
			params.AuctioneerContractAddress,
		},
		[]string{
			string(generatedtests.ACUTestUpgradeBytecode),
			string(generatedtests.SupplyControlTestUpgradeBytecode),
			string(generatedtests.StabilizationTestUpgradeBytecode),
			string(generatedtests.InflationControllerTestUpgradeBytecode),
			string(generatedtests.AuctioneerTestUpgradeBytecode),
		},
	)
	require.NoError(t, err)
	t.Logf("upgrade transaction size: %s", tx.Size().String())
	require.True(t, tx.Size() < ccore.TxMaxSize)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = network.AwaitTransactions(ctx, tx)
	require.NoError(t, err)

	receipt, err := network[0].WsClient.TransactionReceipt(ctx, tx.Hash())
	require.NoError(t, err)
	require.True(t, receipt.Status == types.ReceiptStatusSuccessful)
	t.Logf("gas used %d, gas limit on mainnet %d", receipt.GasUsed, params.AutMainnetChainConfig.AutonityContractConfig.GasLimit)
	require.True(t, receipt.GasUsed < params.AutMainnetChainConfig.AutonityContractConfig.GasLimit)
}

func fetchMinimumBaseFee(t *testing.T, chain *ccore.BlockChain, number *uint64) *big.Int {
	// if number is nil, take the latest core height as number
	if number == nil {
		number = new(uint64)
		*number = chain.CurrentBlock().NumberU64() + 1
	}
	eip1559Params, err := chain.Eip1559ParamsByHeight(*number)
	require.NoError(t, err)
	return eip1559Params.MinBaseFee
}

func fetchGasLimit(t *testing.T, chain *ccore.BlockChain, number *uint64) *big.Int {
	// if number is nil, take the latest core height as number
	if number == nil {
		number = new(uint64)
		*number = chain.CurrentBlock().NumberU64() + 1
	}
	gasLimit, err := chain.GasLimitByHeight(*number)
	require.NoError(t, err)
	return gasLimit
}

func fetchElasticityMultiplier(t *testing.T, chain *ccore.BlockChain, number *uint64) *big.Int {
	// if number is nil, take the latest core height as number
	if number == nil {
		number = new(uint64)
		*number = chain.CurrentBlock().NumberU64() + 1
	}
	eip1559Params, err := chain.Eip1559ParamsByHeight(*number)
	require.NoError(t, err)
	return eip1559Params.ElasticityMultiplier
}

func fetchGasLimitBoundDivisor(t *testing.T, chain *ccore.BlockChain, number *uint64) *big.Int {
	// if number is nil, take the latest core height as number
	if number == nil {
		number = new(uint64)
		*number = chain.CurrentBlock().NumberU64() + 1
	}
	eip1559Params, err := chain.Eip1559ParamsByHeight(*number)
	require.NoError(t, err)
	return eip1559Params.GasLimitBoundDivisor
}

// wrapper function that facilitates setting one or more eip1559 params while leaving the other ones unchanged
func setEip1559Params(t *testing.T, autonity *bindings.Autonity, transactOpts *bind.TransactOpts, customizeFn func(eip1559 *bindings.IAutonityEip1559)) (*types.Transaction, error) {
	config, err := autonity.GetConfig(nil)
	require.NoError(t, err)

	eip1559Params := &bindings.IAutonityEip1559{
		MinBaseFee:               config.Policy.MinBaseFee,
		BaseFeeChangeDenominator: config.Policy.BaseFeeChangeDenominator,
		ElasticityMultiplier:     config.Policy.ElasticityMultiplier,
		GasLimitBoundDivisor:     config.Protocol.GasLimitBoundDivisor,
	}
	customizeFn(eip1559Params)

	return autonity.SetEip1559Params(transactOpts, *eip1559Params)
}

func TestCachedProtocolParameterChange(t *testing.T) {
	t.Run("If minimum base fee is updated, at epoch end cached value is updated as well", func(t *testing.T) {
		network, err := NewNetwork(t, 2, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
		require.NoError(t, err)
		defer network.Shutdown(t)
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		initialMinBaseFee := new(big.Int).SetUint64(params.TestMinBaseFee)
		require.Equal(t, initialMinBaseFee.String(), fetchMinimumBaseFee(t, network[0].Eth.BlockChain(), nil).String())
		require.Equal(t, initialMinBaseFee.String(), fetchMinimumBaseFee(t, network[1].Eth.BlockChain(), nil).String())

		// update min base fee
		updatedMinBaseFee, _ := new(big.Int).SetString("30000000000", 10)
		autonityContract, _ := bindings.NewAutonity(params.AutonityContractAddress, network[0].WsClient)
		transactOpts, _ := bind.NewKeyedTransactorWithChainID(network[0].Key, params.TestChainConfig.ChainID)
		tx, err := setEip1559Params(t, autonityContract, transactOpts, func(params *bindings.IAutonityEip1559) {
			params.MinBaseFee = updatedMinBaseFee
		})
		require.NoError(t, err)
		err = network.AwaitTransactions(ctx, tx)
		require.NoError(t, err)

		// close epoch
		epochPeriod := params.TestAutonityContractConfig.EpochPeriod
		err = network.WaitForHeight(epochPeriod, int(epochPeriod))
		require.NoError(t, err)

		// contract should be updated
		minBaseFee, err := autonityContract.GetMinimumBaseFee(new(bind.CallOpts))
		require.NoError(t, err)
		require.Equal(t, updatedMinBaseFee.String(), minBaseFee.String())

		// caches should be updated from the new epoch start
		require.Equal(t, initialMinBaseFee.String(), fetchMinimumBaseFee(t, network[0].Eth.BlockChain(), &epochPeriod).String())
		require.Equal(t, initialMinBaseFee.String(), fetchMinimumBaseFee(t, network[1].Eth.BlockChain(), &epochPeriod).String())
		firstBlockOfNewEpoch := epochPeriod + 1
		require.Equal(t, updatedMinBaseFee.String(), fetchMinimumBaseFee(t, network[0].Eth.BlockChain(), &firstBlockOfNewEpoch).String())
		require.Equal(t, updatedMinBaseFee.String(), fetchMinimumBaseFee(t, network[1].Eth.BlockChain(), &firstBlockOfNewEpoch).String())
	})
	t.Run("gas limit parameters update", func(t *testing.T) {
		validators, err := Validators(t, 2, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
		require.NoError(t, err)
		// set genesis gas limit to the target to ease future computations
		network, err := NewNetworkFromValidators(t, validators, true, func(genesis *ccore.Genesis) {
			genesis.GasLimit = params.TestChainConfig.AutonityContractConfig.GasLimit
		})
		require.NoError(t, err)
		defer network.Shutdown(t)

		initialGasLimit := new(big.Int).SetUint64(params.TestChainConfig.AutonityContractConfig.GasLimit)
		t.Logf("initial gas limit %s", initialGasLimit.String())
		require.Equal(t, initialGasLimit.String(), fetchGasLimit(t, network[0].Eth.BlockChain(), nil).String())
		require.Equal(t, initialGasLimit.String(), fetchGasLimit(t, network[1].Eth.BlockChain(), nil).String())

		// mine a couple blocks
		require.NoError(t, network.WaitToMineNBlocks(5, 20, false))

		startBlock := network[0].Eth.BlockChain().CurrentBlock()
		startGasLimit := startBlock.GasLimit()
		startNumber := startBlock.NumberU64()
		t.Logf("at block %d, gasLimit: %d", startNumber, startGasLimit)

		// update gas limit
		updatedGasLimit := new(big.Int).SetUint64(initialGasLimit.Uint64() * 10)
		t.Logf("updated gas limit %s", updatedGasLimit.String())
		autonityContract, _ := bindings.NewAutonity(params.AutonityContractAddress, network[0].WsClient)
		transactOpts, _ := bind.NewKeyedTransactorWithChainID(network[0].Key, params.TestChainConfig.ChainID)
		tx, err := autonityContract.SetGasLimit(transactOpts, updatedGasLimit)
		require.NoError(t, err)
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		err = network.AwaitTransactions(ctx, tx)
		require.NoError(t, err)

		// check in which block the transaction was mined
		receipt, err := network[0].WsClient.TransactionReceipt(ctx, tx.Hash())
		require.NoError(t, err)
		changeBlockNumber := receipt.BlockNumber.Uint64()
		t.Logf("gas limit update tx mined at block: %d", changeBlockNumber)

		// contract should be updated
		config, err := autonityContract.GetConfig(new(bind.CallOpts))
		require.NoError(t, err)
		require.Equal(t, updatedGasLimit.String(), config.Protocol.GasLimit.String())

		// caches should be updated only from the block after the change tx was mined
		require.Equal(t, initialGasLimit.String(), fetchGasLimit(t, network[0].Eth.BlockChain(), &changeBlockNumber).String())
		require.Equal(t, initialGasLimit.String(), fetchGasLimit(t, network[1].Eth.BlockChain(), &changeBlockNumber).String())
		firstBlockAfterChange := changeBlockNumber + 1
		require.Equal(t, updatedGasLimit.String(), fetchGasLimit(t, network[0].Eth.BlockChain(), &firstBlockAfterChange).String())
		require.Equal(t, updatedGasLimit.String(), fetchGasLimit(t, network[1].Eth.BlockChain(), &firstBlockAfterChange).String())

		// mined block gas limit should start to increase towards the new limit
		err = network.WaitToMineNBlocks(20, 30, false)
		require.NoError(t, err)

		endBlock := network[0].Eth.BlockChain().CurrentBlock()
		endGasLimit := endBlock.GasLimit()
		endNumber := endBlock.NumberU64()

		t.Logf("at block %d, gasLimit: %d", endNumber, endGasLimit)

		diff := endNumber - startNumber
		t.Logf("diff: %d, gasLimit now: %d, gasLimit before: %d", diff, endGasLimit, startGasLimit)
		require.True(t, endGasLimit > startGasLimit)

		// verify that the increase is correct
		numBlocks := endNumber - changeBlockNumber

		expectedIncrease := uint64(0)
		for i := 0; i < int(numBlocks); i++ {
			expectedIncrease += (startGasLimit+expectedIncrease)/params.DefaultGasLimitBoundDivisor - 1
		}
		t.Logf("expected increase: %d, actual increase: %d", expectedIncrease, endGasLimit-startGasLimit)
		require.Equal(t, expectedIncrease, endGasLimit-startGasLimit)

		// now change the gas limit bound divisor, the gas limit should start increasing faster
		updatedGasLimitBoundDivisor := new(big.Int).SetUint64(params.DefaultGasLimitBoundDivisor / 10)
		t.Logf("updated gas limit bound divisor %s (from %d)", updatedGasLimitBoundDivisor.String(), params.DefaultGasLimitBoundDivisor)
		tx, err = setEip1559Params(t, autonityContract, transactOpts, func(params *bindings.IAutonityEip1559) {
			params.GasLimitBoundDivisor = updatedGasLimitBoundDivisor
		})
		require.NoError(t, err)
		ctx, cancel2 := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel2()
		err = network.AwaitTransactions(ctx, tx)
		require.NoError(t, err)

		ctx, cancel3 := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel3()
		receipt, err = network[0].WsClient.TransactionReceipt(ctx, tx.Hash())
		t.Logf("gas limit bound divisor change tx mined at block %d", receipt.BlockNumber.Uint64())

		require.NoError(t, network.WaitToMineNBlocks(2, 10, false))

		// cache should be on the old value still
		require.Equal(t, params.DefaultGasLimitBoundDivisor, fetchGasLimitBoundDivisor(t, network[0].Eth.BlockChain(), nil).Uint64())

		// close the epoch so change is applied
		epochPeriod := params.TestChainConfig.AutonityContractConfig.EpochPeriod
		require.True(t, network[0].Eth.BlockChain().CurrentBlock().NumberU64() < epochPeriod)
		require.NoError(t, network.WaitForHeight(epochPeriod+5, int(epochPeriod*2)))

		require.Equal(t, params.DefaultGasLimitBoundDivisor, fetchGasLimitBoundDivisor(t, network[0].Eth.BlockChain(), &epochPeriod).Uint64())
		firstBlock := epochPeriod + 1
		require.Equal(t, params.DefaultGasLimitBoundDivisor/10, fetchGasLimitBoundDivisor(t, network[0].Eth.BlockChain(), &firstBlock).Uint64())

		startBlock = network[0].Eth.BlockChain().CurrentBlock()
		startNumber = startBlock.NumberU64()
		startGasLimit = startBlock.GasLimit()

		require.NoError(t, network.WaitToMineNBlocks(20, 30, false))

		endBlock = network[0].Eth.BlockChain().CurrentBlock()
		endNumber = endBlock.NumberU64()
		endGasLimit = endBlock.GasLimit()

		t.Logf("startNumber: %d, gasLimit: %d", startNumber, startGasLimit)
		t.Logf("endNumber: %d, gasLimit: %d", endNumber, endGasLimit)
		t.Logf("number diff: %d, gas diff: %d", endNumber-startNumber, endGasLimit-startGasLimit)

		expectedIncrease = uint64(0)
		for i := 0; i < int(endNumber-startNumber); i++ {
			expectedIncrease += (startGasLimit+expectedIncrease)/updatedGasLimitBoundDivisor.Uint64() - 1
		}
		t.Logf("expected increase: %d, actual increase: %d", expectedIncrease, endGasLimit-startGasLimit)
		require.Equal(t, expectedIncrease, endGasLimit-startGasLimit)

		// now increase the gasLimit bound divisor to a very high amount > gasLimit
		previousGasLimitBoundDivisor := new(big.Int).Set(updatedGasLimitBoundDivisor)
		updatedGasLimitBoundDivisor = new(big.Int).SetUint64(updatedGasLimit.Uint64() * 10)
		t.Logf("updated gas limit bound divisor %s (from %s)", updatedGasLimitBoundDivisor.String(), previousGasLimitBoundDivisor.String())
		tx, err = setEip1559Params(t, autonityContract, transactOpts, func(params *bindings.IAutonityEip1559) {
			params.GasLimitBoundDivisor = updatedGasLimitBoundDivisor
		})
		require.NoError(t, err)
		ctx, cancel4 := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel4()
		err = network.AwaitTransactions(ctx, tx)
		require.NoError(t, err)

		ctx, cancel5 := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel5()
		receipt, err = network[0].WsClient.TransactionReceipt(ctx, tx.Hash())
		t.Logf("gas limit bound divisor change tx mined at block %d", receipt.BlockNumber.Uint64())

		// close the epoch so change is applied
		require.True(t, network[0].Eth.BlockChain().CurrentBlock().NumberU64() < epochPeriod*2)
		require.NoError(t, network.WaitForHeight(epochPeriod*2+5, int(epochPeriod*2)))

		// mine some blocks with the new params
		require.NoError(t, network.WaitToMineNBlocks(10, 20, false))
	})
	t.Run("gas limit == gas limit bound divisor doesn't cause issues", func(t *testing.T) {
		validators, err := Validators(t, 2, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
		require.NoError(t, err)
		network, err := NewNetworkFromValidators(t, validators, true, func(genesis *ccore.Genesis) {
			genesis.GasLimit = params.MinGasLimit
			genesis.Config.AutonityContractConfig.GasLimit = params.MinGasLimit
			genesis.Config.AutonityContractConfig.GasLimitBoundDivisor = params.MinGasLimit
		})
		require.NoError(t, err)
		defer network.Shutdown(t)

		require.NoError(t, network.WaitToMineNBlocks(20, 30, false))
	})
	t.Run("gas limit < min gas limit doesn't cause issues", func(t *testing.T) {
		validators, err := Validators(t, 2, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
		require.NoError(t, err)
		network, err := NewNetworkFromValidators(t, validators, true, func(genesis *ccore.Genesis) {
			genesis.Config.AutonityContractConfig.GasLimit = 1
		})
		require.NoError(t, err)
		defer network.Shutdown(t)

		require.NoError(t, network.WaitToMineNBlocks(20, 30, false))
	})
	t.Run("minbasefee == 0 does not cause issues", func(t *testing.T) {
		validators, err := Validators(t, 2, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
		require.NoError(t, err)
		network, err := NewNetworkFromValidators(t, validators, true, func(genesis *ccore.Genesis) {
			genesis.Config.AutonityContractConfig.MinBaseFee = 0
			// we want the base fee to stay to 0
			genesis.BaseFee = new(big.Int)
		})
		require.NoError(t, err)
		defer network.Shutdown(t)

		require.NoError(t, network.WaitToMineNBlocks(10, 20, false))

		// send tx with gasprice 0, should be mined
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		rawTx := types.NewTransaction(0, network[1].Address, new(big.Int).SetUint64(10), 30_000, new(big.Int).SetUint64(0), nil)
		signedTx, err := types.SignTx(rawTx, types.LatestSigner(network[0].EthConfig.Genesis.Config), network[0].Key)
		require.NoError(t, err)
		require.Equal(t, common.Big0.String(), signedTx.GasPrice().String())

		err = network[0].WsClient.SendTransaction(ctx, signedTx)
		require.NoError(t, err)

		err = network.AwaitTransactions(ctx, signedTx)
		require.NoError(t, err)

		// check that it was mined
		receipt, err := network[0].WsClient.TransactionReceipt(ctx, signedTx.Hash())
		require.NoError(t, err)

		t.Logf("tx mined at block %d", receipt.BlockNumber.Uint64())
		require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

		// same with dynamic tx
		ctx, cancel2 := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel2()
		rawTx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   params.TestChainConfig.ChainID,
			Nonce:     uint64(1),
			GasTipCap: new(big.Int),
			GasFeeCap: new(big.Int),
			Gas:       30_000,
			To:        &network[1].Address,
			Value:     new(big.Int).SetUint64(10),
			Data:      nil,
		})
		signedTx, err = types.SignTx(rawTx, types.LatestSigner(network[0].EthConfig.Genesis.Config), network[0].Key)
		require.NoError(t, err)
		require.Equal(t, common.Big0.String(), signedTx.GasPrice().String())

		err = network[0].WsClient.SendTransaction(ctx, signedTx)
		require.NoError(t, err)

		err = network.AwaitTransactions(ctx, signedTx)
		require.NoError(t, err)

		// check that it was mined
		receipt, err = network[0].WsClient.TransactionReceipt(ctx, signedTx.Hash())
		require.NoError(t, err)

		t.Logf("dynamic tx mined at block %d", receipt.BlockNumber.Uint64())
		require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
	})
	t.Run("testing changes of baseFeeChangeDenominator", func(t *testing.T) {
		validators, err := Validators(t, 2, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
		require.NoError(t, err)
		network, err := NewNetworkFromValidators(t, validators, true, func(genesis *ccore.Genesis) {
			// set a high gas limit, base fee should decrease
			// faster with lower denominator
			// slower with higher
			genesis.BaseFee = new(big.Int).SetUint64(100_000_000_000_00)
			genesis.GasLimit = 100_000_000
			genesis.Config.AutonityContractConfig.GasLimit = genesis.GasLimit * 10
			genesis.Config.AutonityContractConfig.BaseFeeChangeDenominator = 32
		})
		require.NoError(t, err)
		defer network.Shutdown(t)

		startBlock := network[0].Eth.BlockChain().CurrentBlock()
		t.Logf("start %s: baseFee %s, gasUsed %d", startBlock.Number().String(), startBlock.BaseFee().String(), startBlock.GasUsed())

		require.NoError(t, network.WaitToMineNBlocks(10, 20, false))

		endBlock := network[0].Eth.BlockChain().GetBlockByNumber(startBlock.Number().Uint64() + 5)
		t.Logf("end %s: baseFee %s, gasUsed %d", endBlock.Number().String(), endBlock.BaseFee().String(), endBlock.GasUsed())

		require.True(t, startBlock.BaseFee().Uint64() > endBlock.BaseFee().Uint64())

		diffBaseFee := new(big.Int).Sub(startBlock.BaseFee(), endBlock.BaseFee())
		diffBaseFeeFloat, _ := diffBaseFee.Float64()
		startBaseFeeFloat, _ := startBlock.BaseFee().Float64()
		diffBaseFeePerc := diffBaseFeeFloat * 100 / startBaseFeeFloat
		t.Logf("diff baseFee %s (%.2f %%)", diffBaseFee.String(), diffBaseFeePerc)

		// decrease the value of denominator, base fee should drop faster
		updatedBaseFeeChangeDenominator := new(big.Int).SetUint64(16)
		t.Logf("updated base fee change denominator %s (from 32)", updatedBaseFeeChangeDenominator.String())
		autonityContract, _ := bindings.NewAutonity(params.AutonityContractAddress, network[0].WsClient)
		transactOpts, _ := bind.NewKeyedTransactorWithChainID(network[0].Key, params.TestChainConfig.ChainID)
		tx, err := setEip1559Params(t, autonityContract, transactOpts, func(params *bindings.IAutonityEip1559) {
			params.BaseFeeChangeDenominator = updatedBaseFeeChangeDenominator
		})
		require.NoError(t, err)
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		err = network.AwaitTransactions(ctx, tx)
		require.NoError(t, err)

		receipt, err := network[0].WsClient.TransactionReceipt(ctx, tx.Hash())
		t.Logf("basefee change denominator change tx mined at block %d", receipt.BlockNumber.Uint64())

		// close the epoch so change is applied
		epochPeriod := params.TestChainConfig.AutonityContractConfig.EpochPeriod
		require.True(t, network[0].Eth.BlockChain().CurrentBlock().NumberU64() < epochPeriod)
		require.NoError(t, network.WaitForHeight(epochPeriod+5, int(epochPeriod*2)))

		startBlock = network[0].Eth.BlockChain().CurrentBlock()
		t.Logf("start %s: baseFee %s, gasUsed %d", startBlock.Number().String(), startBlock.BaseFee().String(), startBlock.GasUsed())

		require.NoError(t, network.WaitToMineNBlocks(10, 20, false))

		endBlock = network[0].Eth.BlockChain().GetBlockByNumber(startBlock.Number().Uint64() + 5)
		t.Logf("end %s: baseFee %s, gasUsed %d", endBlock.Number().String(), endBlock.BaseFee().String(), endBlock.GasUsed())

		require.True(t, startBlock.BaseFee().Uint64() > endBlock.BaseFee().Uint64())

		diffBaseFee2 := new(big.Int).Sub(startBlock.BaseFee(), endBlock.BaseFee())
		diffBaseFee2Float, _ := diffBaseFee2.Float64()
		startBaseFeeFloat, _ = startBlock.BaseFee().Float64()
		diffBaseFee2Perc := diffBaseFee2Float * 100 / startBaseFeeFloat
		t.Logf("diff baseFee %s (%.2f %%)", diffBaseFee2.String(), diffBaseFee2Perc)

		// basefee should have gone down faster with a lower denominator
		t.Logf("phase 1 decrease: %.2f, phase 2 decrease: %.2f", diffBaseFeePerc, diffBaseFee2Perc)
		require.True(t, diffBaseFee2Perc > diffBaseFeePerc)

		// increase the value of denominator, base fee should drop slower
		updatedBaseFeeChangeDenominator = new(big.Int).SetUint64(64)
		t.Logf("updated base fee change denominator %s (from 16)", updatedBaseFeeChangeDenominator.String())
		tx, err = setEip1559Params(t, autonityContract, transactOpts, func(params *bindings.IAutonityEip1559) {
			params.BaseFeeChangeDenominator = updatedBaseFeeChangeDenominator
		})
		require.NoError(t, err)
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		err = network.AwaitTransactions(ctx, tx)
		require.NoError(t, err)

		receipt, err = network[0].WsClient.TransactionReceipt(ctx, tx.Hash())
		t.Logf("base fee change denominator change tx mined at block %d", receipt.BlockNumber.Uint64())

		// close the epoch so change is applied
		epochPeriod = params.TestChainConfig.AutonityContractConfig.EpochPeriod
		require.True(t, network[0].Eth.BlockChain().CurrentBlock().NumberU64() < epochPeriod*2)
		require.NoError(t, network.WaitForHeight(epochPeriod*2+5, int(epochPeriod*2)))

		startBlock = network[0].Eth.BlockChain().CurrentBlock()
		t.Logf("start %s: baseFee %s, gasUsed %d", startBlock.Number().String(), startBlock.BaseFee().String(), startBlock.GasUsed())

		require.NoError(t, network.WaitToMineNBlocks(10, 20, false))

		endBlock = network[0].Eth.BlockChain().GetBlockByNumber(startBlock.Number().Uint64() + 5)
		t.Logf("end %s: baseFee %s, gasUsed %d", endBlock.Number().String(), endBlock.BaseFee().String(), endBlock.GasUsed())

		require.True(t, startBlock.BaseFee().Uint64() > endBlock.BaseFee().Uint64())

		diffBaseFee3 := new(big.Int).Sub(startBlock.BaseFee(), endBlock.BaseFee())
		diffBaseFee3Float, _ := diffBaseFee3.Float64()
		startBaseFeeFloat, _ = startBlock.BaseFee().Float64()
		diffBaseFee3Perc := diffBaseFee3Float * 100 / startBaseFeeFloat
		t.Logf("diff baseFee %s (%.2f %%)", diffBaseFee3.String(), diffBaseFee3Perc)

		// basefee should have gone down slower with a higher denominator
		t.Logf("phase 1 decrease: %.2f, phase 2 decrease: %.2f, phase 3 decrease: %.2f", diffBaseFeePerc, diffBaseFee2Perc, diffBaseFee3Perc)
		require.True(t, diffBaseFee3Perc < diffBaseFeePerc)
	})
	t.Run("test very big base fee change denominator", func(t *testing.T) {
		validators, err := Validators(t, 2, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
		require.NoError(t, err)
		network, err := NewNetworkFromValidators(t, validators, true, func(genesis *ccore.Genesis) {
			genesis.BaseFee = new(big.Int).SetUint64(100_000_000_000_00)
			genesis.GasLimit = 100_000_000
			genesis.Config.AutonityContractConfig.GasLimit = genesis.GasLimit * 10
			genesis.Config.AutonityContractConfig.BaseFeeChangeDenominator = math.MaxUint64
		})
		require.NoError(t, err)
		defer network.Shutdown(t)

		startBlock := network[0].Eth.BlockChain().CurrentBlock()
		t.Logf("start %s: baseFee %s, gasUsed %d", startBlock.Number().String(), startBlock.BaseFee().String(), startBlock.GasUsed())

		require.NoError(t, network.WaitToMineNBlocks(10, 20, false))

		endBlock := network[0].Eth.BlockChain().GetBlockByNumber(startBlock.Number().Uint64() + 5)
		t.Logf("end %s: baseFee %s, gasUsed %d", endBlock.Number().String(), endBlock.BaseFee().String(), endBlock.GasUsed())

		require.True(t, startBlock.BaseFee().Uint64() > endBlock.BaseFee().Uint64())

		diffBaseFee := new(big.Int).Sub(startBlock.BaseFee(), endBlock.BaseFee())
		diffBaseFeeFloat, _ := diffBaseFee.Float64()
		startBaseFeeFloat, _ := startBlock.BaseFee().Float64()
		diffBaseFeePerc := diffBaseFeeFloat * 100 / startBaseFeeFloat
		t.Logf("diff baseFee %s (%.2f %%)", diffBaseFee.String(), diffBaseFeePerc)
	})
	t.Run("test very big elasticity multiplier", func(t *testing.T) {
		validators, err := Validators(t, 2, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
		require.NoError(t, err)
		network, err := NewNetworkFromValidators(t, validators, true, func(genesis *ccore.Genesis) {
			genesis.BaseFee = new(big.Int).SetUint64(100_000_000_000_00)
			genesis.GasLimit = 100_000_000
			genesis.Config.AutonityContractConfig.GasLimit = genesis.GasLimit * 10
			genesis.Config.AutonityContractConfig.ElasticityMultiplier = math.MaxUint64
			genesis.GasUsed = 100
		})
		require.NoError(t, err)
		defer network.Shutdown(t)

		startBlock := network[0].Eth.BlockChain().CurrentBlock()
		t.Logf("start %s: baseFee %s, gasUsed %d", startBlock.Number().String(), startBlock.BaseFee().String(), startBlock.GasUsed())

		require.NoError(t, network.WaitToMineNBlocks(10, 20, false))

		endBlock := network[0].Eth.BlockChain().GetBlockByNumber(startBlock.Number().Uint64() + 5)
		t.Logf("end %s: baseFee %s, gasUsed %d", endBlock.Number().String(), endBlock.BaseFee().String(), endBlock.GasUsed())

		require.True(t, startBlock.BaseFee().Uint64() < endBlock.BaseFee().Uint64())

		diffBaseFee := new(big.Int).Sub(endBlock.BaseFee(), startBlock.BaseFee())
		diffBaseFeeFloat, _ := diffBaseFee.Float64()
		startBaseFeeFloat, _ := startBlock.BaseFee().Float64()
		diffBaseFeePerc := diffBaseFeeFloat * 100 / startBaseFeeFloat
		t.Logf("diff baseFee %s (%.2f %%)", diffBaseFee.String(), diffBaseFeePerc)
	})
	t.Run("testing changes of elasticityMultiplier", func(t *testing.T) {
		validators, err := Validators(t, 2, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
		require.NoError(t, err)
		network, err := NewNetworkFromValidators(t, validators, true, func(genesis *ccore.Genesis) {
			genesis.BaseFee = new(big.Int).SetUint64(100_000_000_000_00)
			genesis.GasLimit = 100_000_000
			genesis.Config.AutonityContractConfig.GasLimit = genesis.GasLimit * 10
			genesis.Config.AutonityContractConfig.ElasticityMultiplier = 4
			genesis.Config.AutonityContractConfig.BaseFeeChangeDenominator = 64
		})
		require.NoError(t, err)
		defer network.Shutdown(t)

		startBlock := network[0].Eth.BlockChain().CurrentBlock()
		t.Logf("start %s: baseFee %s, gasUsed %d", startBlock.Number().String(), startBlock.BaseFee().String(), startBlock.GasUsed())

		require.NoError(t, network.WaitToMineNBlocks(10, 20, false))

		endBlock := network[0].Eth.BlockChain().GetBlockByNumber(startBlock.Number().Uint64() + 5)
		t.Logf("end %s: baseFee %s, gasUsed %d", endBlock.Number().String(), endBlock.BaseFee().String(), endBlock.GasUsed())

		require.True(t, startBlock.BaseFee().Uint64() > endBlock.BaseFee().Uint64())

		diffBaseFee := new(big.Int).Sub(startBlock.BaseFee(), endBlock.BaseFee())
		diffBaseFeeFloat, _ := diffBaseFee.Float64()
		startBaseFeeFloat, _ := startBlock.BaseFee().Float64()
		diffBaseFeePerc := diffBaseFeeFloat * 100 / startBaseFeeFloat
		t.Logf("diff baseFee %s (%.2f %%)", diffBaseFee.String(), diffBaseFeePerc)

		updatedElasticityMultiplier := new(big.Int).SetUint64(8)
		t.Logf("updated elasticity multiplier %s (from 4)", updatedElasticityMultiplier.String())
		autonityContract, _ := bindings.NewAutonity(params.AutonityContractAddress, network[0].WsClient)
		transactOpts, _ := bind.NewKeyedTransactorWithChainID(network[0].Key, params.TestChainConfig.ChainID)
		tx, err := setEip1559Params(t, autonityContract, transactOpts, func(params *bindings.IAutonityEip1559) {
			params.ElasticityMultiplier = updatedElasticityMultiplier
		})
		require.NoError(t, err)
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		err = network.AwaitTransactions(ctx, tx)
		require.NoError(t, err)

		receipt, err := network[0].WsClient.TransactionReceipt(ctx, tx.Hash())
		t.Logf("elasticity multiplier change tx mined at block %d", receipt.BlockNumber.Uint64())

		// close the epoch so change is applied
		epochPeriod := params.TestChainConfig.AutonityContractConfig.EpochPeriod
		require.True(t, network[0].Eth.BlockChain().CurrentBlock().NumberU64() < epochPeriod)
		require.NoError(t, network.WaitForHeight(epochPeriod+5, int(epochPeriod*2)))

		// check if cache has changed
		require.Equal(t, uint64(4), fetchElasticityMultiplier(t, network[0].Eth.BlockChain(), &epochPeriod).Uint64())
		require.Equal(t, uint64(4), fetchElasticityMultiplier(t, network[1].Eth.BlockChain(), &epochPeriod).Uint64())
		firstBlock := epochPeriod + 1
		require.Equal(t, uint64(8), fetchElasticityMultiplier(t, network[0].Eth.BlockChain(), &firstBlock).Uint64())
		require.Equal(t, uint64(8), fetchElasticityMultiplier(t, network[1].Eth.BlockChain(), &firstBlock).Uint64())

		startBlock = network[0].Eth.BlockChain().CurrentBlock()
		t.Logf("start %s: baseFee %s, gasUsed %d", startBlock.Number().String(), startBlock.BaseFee().String(), startBlock.GasUsed())

		require.NoError(t, network.WaitToMineNBlocks(10, 20, false))

		endBlock = network[0].Eth.BlockChain().GetBlockByNumber(startBlock.Number().Uint64() + 5)
		t.Logf("end %s: baseFee %s, gasUsed %d", endBlock.Number().String(), endBlock.BaseFee().String(), endBlock.GasUsed())

		require.True(t, startBlock.BaseFee().Uint64() > endBlock.BaseFee().Uint64())

		diffBaseFee2 := new(big.Int).Sub(startBlock.BaseFee(), endBlock.BaseFee())
		diffBaseFee2Float, _ := diffBaseFee2.Float64()
		startBaseFeeFloat, _ = startBlock.BaseFee().Float64()
		diffBaseFee2Perc := diffBaseFee2Float * 100 / startBaseFeeFloat
		t.Logf("diff baseFee %s (%.2f %%)", diffBaseFee2.String(), diffBaseFee2Perc)

		t.Logf("phase 1 decrease: %.2f, phase 2 decrease: %.2f", diffBaseFeePerc, diffBaseFee2Perc)

		updatedElasticityMultiplier = new(big.Int).SetUint64(2)
		t.Logf("updated elasticity multiplier %s (from 8)", updatedElasticityMultiplier.String())
		tx, err = setEip1559Params(t, autonityContract, transactOpts, func(params *bindings.IAutonityEip1559) {
			params.ElasticityMultiplier = updatedElasticityMultiplier
		})
		require.NoError(t, err)
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		err = network.AwaitTransactions(ctx, tx)
		require.NoError(t, err)

		receipt, err = network[0].WsClient.TransactionReceipt(ctx, tx.Hash())
		t.Logf("gas limit bound divisor change tx mined at block %d", receipt.BlockNumber.Uint64())

		// close the epoch so change is applied
		epochPeriod = params.TestChainConfig.AutonityContractConfig.EpochPeriod
		require.True(t, network[0].Eth.BlockChain().CurrentBlock().NumberU64() < epochPeriod*2)
		require.NoError(t, network.WaitForHeight(epochPeriod*2+5, int(epochPeriod*2)))

		// check if cache has changed
		lastBlockOfEpoch := epochPeriod * 2
		require.Equal(t, uint64(8), fetchElasticityMultiplier(t, network[0].Eth.BlockChain(), &lastBlockOfEpoch).Uint64())
		require.Equal(t, uint64(8), fetchElasticityMultiplier(t, network[1].Eth.BlockChain(), &lastBlockOfEpoch).Uint64())
		firstBlock = lastBlockOfEpoch + 1
		require.Equal(t, uint64(2), fetchElasticityMultiplier(t, network[0].Eth.BlockChain(), &firstBlock).Uint64())
		require.Equal(t, uint64(2), fetchElasticityMultiplier(t, network[1].Eth.BlockChain(), &firstBlock).Uint64())

		startBlock = network[0].Eth.BlockChain().CurrentBlock()
		t.Logf("start %s: baseFee %s, gasUsed %d", startBlock.Number().String(), startBlock.BaseFee().String(), startBlock.GasUsed())

		require.NoError(t, network.WaitToMineNBlocks(10, 20, false))

		endBlock = network[0].Eth.BlockChain().GetBlockByNumber(startBlock.Number().Uint64() + 5)
		t.Logf("end %s: baseFee %s, gasUsed %d", endBlock.Number().String(), endBlock.BaseFee().String(), endBlock.GasUsed())

		require.True(t, startBlock.BaseFee().Uint64() > endBlock.BaseFee().Uint64())

		diffBaseFee3 := new(big.Int).Sub(startBlock.BaseFee(), endBlock.BaseFee())
		diffBaseFee3Float, _ := diffBaseFee3.Float64()
		startBaseFeeFloat, _ = startBlock.BaseFee().Float64()
		diffBaseFee3Perc := diffBaseFee3Float * 100 / startBaseFeeFloat
		t.Logf("diff baseFee %s (%.2f %%)", diffBaseFee3.String(), diffBaseFee3Perc)

		t.Logf("phase 1 decrease: %.2f, phase 2 decrease: %.2f, phase 3 decrease: %.2f", diffBaseFeePerc, diffBaseFee2Perc, diffBaseFee3Perc)
	})
}

func assembleClientConfig(t *testing.T, node *Node, target uint64) *types.ContractsConfig {
	eip1559Params, err := node.Eth.BlockChain().Eip1559ParamsByHeight(target)
	require.NoError(t, err)
	accountabilityParams, err := node.Eth.BlockChain().AccountabilityParamsByHeight(target)
	require.NoError(t, err)
	epochPeriod, err := node.Eth.BlockChain().EpochPeriodByHeight(target)
	require.NoError(t, err)
	gasLimit, err := node.Eth.BlockChain().GasLimitByHeight(target)
	require.NoError(t, err)
	clusteringThreshold, err := node.Eth.BlockChain().ClusteringThresholdByHeight(target)
	require.NoError(t, err)
	return &types.ContractsConfig{
		EpochPeriod:         epochPeriod,
		BlockPeriod:         new(big.Int).SetUint64(1), // for now block period is hardcoded to 1s always
		GasLimit:            gasLimit,
		ClusteringThreshold: clusteringThreshold,
		Accountability:      *accountabilityParams,
		Eip1559:             *eip1559Params,
	}
}

func fetchClientConfig(t *testing.T, node *Node, target uint64) *types.ContractsConfig {
	header := node.Eth.BlockChain().GetHeaderByNumber(target)
	state, err := node.Eth.BlockChain().StateAt(header.Root)
	require.NoError(t, err)
	config, err := node.Eth.BlockChain().ProtocolContracts().CallGetClientConfig(state, header)
	require.NoError(t, err)
	return config
}

func rlpEncodeConfig(config *types.ContractsConfig) []byte {
	data, err := rlp.EncodeToBytes(config)
	if err != nil {
		panic("Failed to RLP encode contracts config: " + err.Error())
	}
	return data
}

func TestCacheCoherentWithState(t *testing.T) {
	users, err := Validators(t, 2, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	network, err := NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	err = network.WaitToMineNBlocks(5, 20, false)
	require.NoError(t, err)

	// cached config and state config should always be coherent
	require.Equal(t,
		rlpEncodeConfig(fetchClientConfig(t, network[0], 0)),
		rlpEncodeConfig(assembleClientConfig(t, network[0], 0)),
	)

	require.Equal(t,
		rlpEncodeConfig(fetchClientConfig(t, network[0], 3)),
		rlpEncodeConfig(assembleClientConfig(t, network[0], 3)),
	)

}

// a change in the delta parameter should be applied at epoch end
func TestOmissionDeltaUpdate(t *testing.T) {
	network, err := NewNetwork(t, 2, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown(t)

	omissionContract, err := bindings.NewOmissionAccountability(params.OmissionAccountabilityContractAddress, network[0].WsClient)
	require.NoError(t, err)
	autonityContract, err := bindings.NewAutonity(params.AutonityContractAddress, network[0].WsClient)
	require.NoError(t, err)

	sendAndWait := func(tx *types.Transaction) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		err = network.AwaitTransactions(ctx, tx)
		require.NoError(t, err)
	}

	// set lookback window to 1, to have more leeway in changing delta
	transactOpts, err := bind.NewKeyedTransactorWithChainID(network[0].Key, params.TestChainConfig.ChainID)
	require.NoError(t, err)
	tx, err := omissionContract.SetLookbackWindow(transactOpts, common.Big1)
	require.NoError(t, err)
	sendAndWait(tx)

	deltaRespected := func(node *Node, delta uint64, epochPeriod uint64) bool {
		currentBlock := node.Eth.BlockChain().CurrentBlock().Number().Uint64()
		if currentBlock%epochPeriod <= delta {
			t.Fatal("Not enough blocks in the epoch") // cannot really check if delta is correctly respected if we don't have at least delta + 1 blocks in the epoch
		}
		epochBlock := currentBlock - (currentBlock % epochPeriod)
		t.Logf("epochBlock: %d, epochPeriod: %d, currentBlock: %d, delta: %d", epochBlock, epochPeriod, currentBlock, delta)
		for h := epochBlock + 1; h <= currentBlock; h++ {
			block := node.Eth.BlockChain().GetBlockByNumber(h)
			t.Logf("Checking block num: %d", h)
			if h <= epochBlock+delta {
				// needs to be empty
				if block.Header().ActivityProof != nil {
					t.Log("Block should be empty but is not")
					return false
				}
			} else {
				// needs to be not empty
				if block.Header().ActivityProof == nil {
					t.Log("Block should not be empty but it is")
					return false
				}
			}

		}
		return true
	}

	initialDelta := params.DefaultOmissionAccountabilityConfig.Delta
	require.NotEqual(t, uint64(1), initialDelta)
	epoch0, err := network[0].Eth.BlockChain().LatestEpoch()
	require.NoError(t, err)
	delta0 := epoch0.OmissionDelta.Uint64()
	epoch1, err := network[1].Eth.BlockChain().LatestEpoch()
	require.NoError(t, err)
	delta1 := epoch1.OmissionDelta.Uint64()
	require.Equal(t, initialDelta, delta0)
	require.Equal(t, initialDelta, delta1)
	t.Log(initialDelta)

	epochPeriodBig, err := autonityContract.GetEpochPeriod(nil)
	require.NoError(t, err)
	epochPeriod := epochPeriodBig.Uint64()
	t.Log(epochPeriod)

	// reduce delta to 2
	newDelta := new(big.Int).Set(common.Big2)
	tx, err = omissionContract.SetDelta(transactOpts, newDelta)
	require.NoError(t, err)
	sendAndWait(tx)

	// we should still be in epoch 0, getDelta should already return the new value
	epochID, err := autonityContract.GetEpochID(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(0), epochID.Uint64())
	newDeltaFetched, err := omissionContract.GetDelta(nil)
	require.NoError(t, err)
	require.Equal(t, newDelta.String(), newDeltaFetched.String())

	// wait for some blocks than check that delta is respected
	err = network.WaitForHeight(initialDelta*2, int(initialDelta*3))
	require.NoError(t, err)
	require.True(t, deltaRespected(network[1], initialDelta, epochPeriod))

	// wait for the epoch to end
	err = network.WaitForHeight(epochPeriod, int(epochPeriod))
	require.NoError(t, err)

	// should be in epoch 1
	epochID, err = autonityContract.GetEpochID(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(1), epochID.Uint64())

	// wait for some more blocks and check that the activity proofs have been set accordingly to the new delta
	err = network.WaitToMineNBlocks(10, 15, false)
	require.NoError(t, err)

	require.True(t, deltaRespected(network[1], newDelta.Uint64(), epochPeriod))

	// update delta again, this time increasing it
	newDelta = new(big.Int).SetUint64(20)
	tx, err = omissionContract.SetDelta(transactOpts, newDelta)
	require.NoError(t, err)
	sendAndWait(tx)

	// wait for the epoch to end
	err = network.WaitForHeight(epochPeriod*2, int(epochPeriod*2))
	require.NoError(t, err)

	// should be in epoch 2
	epochID, err = autonityContract.GetEpochID(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(2), epochID.Uint64())

	// wait for some more blocks and check that the activity proofs have been set accordingly to the new delta
	err = network.WaitToMineNBlocks(35, 50, false)
	require.NoError(t, err)

	require.True(t, deltaRespected(network[1], newDelta.Uint64(), epochPeriod))

}

// This test checks that when a transaction is processed the fees are divided
// between validators and stakeholders.
func TestFeeRedistributionValidatorsAndDelegators(t *testing.T) {
	t.Skip("Is broken with Penalty Absorbing Stake")
	//todo: fix. Genesis validators are no longer issued Liquid Newton. Need to introduce 3rd party delegators.
	vals, err := Validators(t, 3, "10e18,v,10000,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	vals[2].Stake = 25000

	network, err := NewNetworkFromValidators(t, vals, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	n := network[0]

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	// retrieve current balance
	// send liquid newton to some random address
	// check balance - shouldnt have increased
	// wait for epoch
	// check claimable fees
	// redeem fees

	// Setup Bindings
	autonityContract, _ := bindings.NewAutonity(params.AutonityContractAddress, n.WsClient)
	valAddrs, _ := autonityContract.GetValidators(nil)
	liquidStateContracts := make([]*bindings.ILiquid, len(valAddrs))
	validators := make([]bindings.IAutonityValidator, len(valAddrs))
	for i, valAddr := range valAddrs {
		validators[i], _ = autonityContract.GetValidator(nil, valAddr)
		liquidStateContracts[i], _ = bindings.NewILiquid(validators[i].LiquidStateContract, n.WsClient)
	}
	transactor, _ := bind.NewKeyedTransactorWithChainID(vals[0].TreasuryKey, big.NewInt(1234))
	tx, err := liquidStateContracts[0].Transfer(
		transactor,
		common.Address{66, 66}, big.NewInt(1337),
	)

	require.NoError(t, err)
	_ = network.WaitToMineNBlocks(2, 20, false)
	tx2, err := n.SendAUT(ctx, network[1].Address, 10)
	require.NoError(t, err)
	err = network.AwaitTransactions(ctx, tx, tx2)
	require.NoError(t, err)
	// claimable fees should be 0 before epoch
	for i := range liquidStateContracts {
		unclaimed, _ := liquidStateContracts[i].UnclaimedRewards(&bind.CallOpts{}, validators[i].Treasury)
		require.Equal(t, big.NewInt(0).Bytes(), unclaimed.Bytes())
	}

	// wait for epoch

	// calculate reward pool
	r1, _ := n.WsClient.TransactionReceipt(context.Background(), tx.Hash())
	r2, _ := n.WsClient.TransactionReceipt(context.Background(), tx2.Hash())

	b1, _ := n.WsClient.BlockByNumber(context.Background(), r1.BlockNumber)
	b2, _ := n.WsClient.BlockByNumber(context.Background(), r2.BlockNumber)

	rewardT1 := new(big.Int).Mul(new(big.Int).SetUint64(r1.CumulativeGasUsed), b1.BaseFee())
	rewardT2 := new(big.Int).Mul(new(big.Int).SetUint64(r2.CumulativeGasUsed), b2.BaseFee())

	totalFees := new(big.Int).Add(rewardT1, rewardT2)
	treasuryRewards := new(big.Int).Div(new(big.Int).Mul(totalFees, new(big.Int).SetUint64(15)), big.NewInt(10000))
	totalRewards := new(big.Int).Sub(totalFees, treasuryRewards)

	balanceBeforeEpoch, _ := n.WsClient.BalanceAt(context.Background(), params.AutonityContractAddress, nil)
	require.Equal(t, totalFees, balanceBeforeEpoch)

	err = network.WaitToMineNBlocks(30, 90, false)
	require.NoError(t, err)

	fmt.Println("total rewards", totalRewards)
	balanceGlobalTreasury, _ := n.WsClient.BalanceAt(context.Background(), common.Address{120}, nil)
	cfg, _ := autonityContract.GetConfig(nil)
	fmt.Println(cfg)
	require.Equal(t, treasuryRewards, balanceGlobalTreasury)

	stake := []int64{10000 - 1337, 10000, 25000}
	epochStake := []int64{10000, 10000, 25000}
	totalStake := int64(45000)
	for i := range liquidStateContracts {
		unclaimed, _ := liquidStateContracts[i].UnclaimedRewards(&bind.CallOpts{}, validators[i].Treasury)
		totalValRewards := new(big.Int).Div(new(big.Int).Mul(totalRewards, big.NewInt(epochStake[i])), big.NewInt(totalStake))
		valCommission := new(big.Int).Div(new(big.Int).Mul(totalValRewards, big.NewInt(12)), big.NewInt(100))
		stakerReward := new(big.Int).Sub(totalValRewards, valCommission)
		require.Equal(t, new(big.Int).Div(new(big.Int).Mul(stakerReward, big.NewInt(stake[i])), big.NewInt(epochStake[i])), unclaimed)
	}
}

// a node is verifying a proposal, but while he is verifying the finalized block is injected from p2p layer
func TestNodeAlreadyHasProposedBlock(t *testing.T) {
	vals, err := Validators(t, 2, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	// trick to get a handle to Proposer of node 0
	var proposer interfaces.Proposer
	var node0Core *core.Core
	vals[0].TendermintServices = &interfaces.Services{Proposer: func(c interfaces.Core) interfaces.Proposer {
		proposer = &core.Proposer{Core: c.(*core.Core)}
		node0Core = c.(*core.Core)
		return proposer
	}}

	// start the network
	network, err := NewNetworkFromValidators(t, vals, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// mine a couple blocks
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err)

	// stop consensus engine of node 0, this will halt the network
	node := network[0]
	err = node.Eth.Engine().Close()
	require.NoError(t, err)
	// reset current round message to force verify proposal
	node0Core.Messages().Reset()

	// get latest inserted block and generate proposal out of it
	block := node.Eth.BlockChain().CurrentBlock()
	proposal := message.NewPropose(0, block.NumberU64(), -1, block, func(hash common.Hash) blst.Signature {
		return node.ConsensusKey.Sign(hash.Bytes())
	}, &types.CommitteeMember{
		Address:           node.Address,
		VotingPower:       common.Big1,
		ConsensusKeyBytes: node.ConsensusKey.PublicKey().Marshal(),
		ConsensusKey:      node.ConsensusKey.PublicKey(),
	})
	//reset cache to force verify proposal
	ethDb := rawdb.NewMemoryDatabase()
	db := state.NewDatabase(ethDb)
	stateDB, _ := state.New(common.Hash{}, db, nil)
	node0Core.Backend().BlockChain().CacheProposalState(common.Hash{}, nil, 0, stateDB, &types.ContractsConfig{})

	// handle the proposal
	err = proposer.HandleProposal(context.TODO(), proposal)
	require.ErrorIs(t, err, constants.ErrAlreadyHaveBlock)
}

func TestStartingAndStoppingNodes(t *testing.T) {
	network, err := NewNetwork(t, 5, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown(t)
	n := network[0]
	// Send a TX to see that the network is working
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = n.SendAUTtracked(ctx, network[1].Address, 10)
	require.NoError(t, err)
	// Stop a node
	err = network[1].Close(true)
	network[1].Wait()
	require.NoError(t, err)
	// Send a TX to see that the network is working
	ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = n.SendAUTtracked(ctx, network[1].Address, 10)
	require.NoError(t, err)

	// Stop a node
	err = network[2].Close(true)
	network[2].Wait()
	require.NoError(t, err)
	// We have now stopped more than F nodes, so we expect TX processing to time out.
	// Well wait 5 times the avgTransactionDuration before we assume the TX is not being processed.
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = n.SendAUTtracked(ctx, network[1].Address, 10)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expecting %q, instead got: %v ", context.DeadlineExceeded.Error(), err)
	}

	// We start a node again and expect the previously unprocessed transaction to be processed
	err = network[2].Start()
	require.NoError(t, err)

	// Ensure that the previously sent transaction is now processed
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	err = n.AwaitSentTransactions(ctx)
	require.NoError(t, err)
	// Send a TX to see that the network is still working
	ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = n.SendAUTtracked(ctx, network[1].Address, 10)
	require.NoError(t, err)

	// Start the last stopped node
	err = network[1].Start()
	require.NoError(t, err)
	// Send a TX to see that the network is still working
	err = n.SendAUTtracked(context.Background(), network[1].Address, 10)
	require.NoError(t, err)
}

func NewBroadcasterWithCheck(currentHeight uint64, stopped bool) func(c interfaces.Core) interfaces.Broadcaster {
	return func(c interfaces.Core) interfaces.Broadcaster {
		return &broadcasterWithCheck{c.(*core.Core), currentHeight, stopped}
	}
}

type broadcasterWithCheck struct {
	*core.Core
	currentHeight uint64
	stopped       bool
}

func (s *broadcasterWithCheck) Broadcast(msg message.Msg) {
	logger := s.Logger().New("step", s.Step())
	logger.Debug("Broadcasting", "message", msg.String())

	if s.stopped && msg.H() <= s.currentHeight {
		panic("stopped node sent old consensus message once he came back up")
	}

	s.BroadcastAll(msg)
}

// Tests that a stopped node, once restarted, will sync up to chain head before sending consensus messages
func TestWaitForChainSyncAfterStop(t *testing.T) {
	network, err := NewNetwork(t, 5, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown(t)

	err = network.WaitToMineNBlocks(10, 60, false)
	require.NoError(t, err)

	// Stop node 0, he will need to sync up once restarted
	err = network[0].Close(false)
	require.NoError(t, err)
	network[0].Wait()

	err = network.WaitToMineNBlocks(10, 60, false)
	require.NoError(t, err)

	// Stop node 1. This will make the network lose liveness and halt
	err = network[1].Close(false)
	require.NoError(t, err)
	network[1].Wait()

	// network should be stalled now
	err = network.WaitToMineNBlocks(1, 5, false)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expecting %q, instead got: %v ", context.DeadlineExceeded.Error(), err)
	}

	// restart node 1. He will rightfully send a consensus message for the consensus instance he stopped at (chainHeight+1)
	chainHeight := network[2].GetChainHeight()
	network[1].CustHandler = &interfaces.Services{Broadcaster: NewBroadcasterWithCheck(chainHeight, true)}
	err = network[1].Start()
	require.NoError(t, err)

	// give time for node 1 to come back up
	err = network.WaitToMineNBlocks(15, 60, false)
	require.NoError(t, err)

	// restart node 0. He should sync up and not send old consensus messages
	chainHeight = network[2].GetChainHeight()
	network[0].CustHandler = &interfaces.Services{Broadcaster: NewBroadcasterWithCheck(chainHeight, true)}
	err = network[0].Start()
	require.NoError(t, err)

	// give time for node 0 to come back up
	err = network.WaitToMineNBlocks(15, 60, false)
	require.NoError(t, err)

}

// Test details
// a.setup 7 validators with 100 voting power on each, and keep 7 committee seats as well.
// b.start the network with 1st 3 nodes only, the network should be on-hold since the online voting power is less than 2/3 of 7
// c.after the on-holding for a while, start the 4th node, then the network should start to produce blocks without any on-holding.
func TestTendermintQuorum(t *testing.T) {
	users, err := Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	network, err := NewNetworkFromValidators(t, users, false)
	require.NoError(t, err)
	defer network.Shutdown(t)
	for i, n := range network {
		// start 3 nodes
		if i < 3 {
			err = n.Start()
			require.NoError(t, err)
		}
	}
	// check if network on hold
	err = network.WaitToMineNBlocks(3, 60, false)
	require.Error(t, err, "Network is not supposed to be mining blocks at this point")
	// start 4th node
	err = network[3].Start()
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

// a. setup 6 validators with 100 voting power on each, and keep 6 committee seats as well.
// b. start all the 6 validators for a while, as the network is producing blocks for a while.
// c.stop 3 nodes one by one, then the network should on-hold when the online voting power is less than 2/3 of 6.
// d.after the on-holding for a while, recover the stopped nodes, then the network should start to produce blocks without any on-holding.
func TestTendermintQuorum2(t *testing.T) {
	users, err := Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	// creates a network of 6 users and starts all the nodes in it
	network, err := NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err)
	defer network.Shutdown(t)
	// stop 3 nodes and verify if transaction processing is halted
	for i, n := range network {
		// stop last 3 nodes
		if i > 2 {
			err = n.Close(true)
			n.Wait()
			require.NoError(t, err)
		}
	}
	// check if network on hold
	err = network.WaitToMineNBlocks(3, 60, false)
	require.Error(t, err, "Network is not supposed to be mining blocks at this point")

	// start the nodes back again
	for i, n := range network {
		// start back 4,5 & 6th node
		if i > 2 {
			err = n.Start()
			require.NoError(t, err)

		}
	}
	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

// a. setup 7 validators (A, B, C, D, E, F, G) with 100 voting power on each, and keep 7 committee seats as well.
// b. start the network for a while(the network is producing blocks for a while).
// c. Shut down nodes A and B, now the network keeps liveness with C, D, E, F online, TXs should be mined.
// d. After a while, shut down Node C and D, with only E and F online, the network should on-hold, no transaction should be mined.
// e. Then start up node A and B, they should be synchronized with node E and F, and network liveness should recovered, new TXs should be mined after the recover.
// f. Then start up node C and D, both C and D should get synchronized to the latest chain height.
// g. Then shut down node E and F, the network should still keep liveness, TXs are mined.
// h. Recover E and F, they should get synchronized finally.
func TestTendermintQuorum4(t *testing.T) {
	users, err := Validators(t, 7, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	// creates a network of 7 users and starts all the nodes in it
	network, err := NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown(t)
	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
	i := 0
	for i < 2 {
		// stop 1st and 2nd node
		err = network[i].Close(true)
		require.NoError(t, err)
		network[i].Wait()
		i++
	}
	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
	// shutting down 3rd and 4th node
	for i < 4 {
		// stop 1st and 2nd node
		err = network[i].Close(true)
		require.NoError(t, err)
		network[i].Wait()
		i++
	}

	// network should be on hold
	err = network.WaitToMineNBlocks(3, 60, false)
	require.Error(t, err, "Network is not supposed to be mining blocks at this point")

	// start back 1st and 2nd node
	i = 0
	for i < 2 {
		err = network[i].Start()
		require.NoError(t, err)
		i++
	}

	err = network.WaitForSyncComplete()
	require.NoError(t, err)

	// we are restoring liveliness we would wait for
	// network should be back up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")

	// start back 3rd and 4th node
	for i < 4 {
		err = network[i].Start()
		require.NoError(t, err)
		i++
	}

	// wait for sync completion
	err = network.WaitForSyncComplete()
	require.NoError(t, err)

	// network should be back up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")

	// bring down 5th and 6th node
	for i < 6 {
		// stop 1st and 2nd node
		err = network[i].Close(true)
		require.NoError(t, err)
		network[i].Wait()
		i++
	}

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
	i = 4
	// start back 5th and 6th node
	for i < 4 {
		err = network[i].Start()
		require.NoError(t, err)
		i++
	}

	// wait for sync completion
	err = network.WaitForSyncComplete()
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

// setup up a network of 12 nodes
// ensure the newtwork is running and blocks are getting mined
// start/stop nodes in parallel
// wait for network to go on hold
// restart all nodes and ensures that network resumes mining new blocks
func TestStartStopAllNodesInParallel(t *testing.T) {
	const nodeCount = 12
	users, err := Validators(t, nodeCount, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	// creates a network of 6 users and starts all the nodes in it
	network, err := NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown(t)
	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")

	type nodeStatus struct {
		lock   sync.Mutex
		status bool
	}
	var wg sync.WaitGroup
	mlock := sync.RWMutex{}
	m := make(map[int]*nodeStatus, nodeCount) // map to maintain start status of all nodes
	for i := 0; i < nodeCount; i++ {
		m[i] = &nodeStatus{sync.Mutex{}, true}
	}

	// randomly start stop nodes
	for j := 0; j < 40; j++ {
		i := rand.Intn(nodeCount) //nolint:gosec
		switch i % 4 {
		case 0, 1, 2:
			wg.Add(1)
			go func() {
				defer wg.Done()
				mlock.RLock()
				nodeStatus := m[i]
				mlock.RUnlock()
				nodeStatus.lock.Lock()
				defer nodeStatus.lock.Unlock()
				if !nodeStatus.status {
					return
				}
				e := network[i].Close(true)
				require.NoError(t, e)
				network[i].Wait()
				nodeStatus.status = false
				mlock.Lock()
				m[i] = nodeStatus
				mlock.Unlock()
			}()
		case 3:
			wg.Add(1)
			go func() {
				defer wg.Done()
				mlock.RLock()
				nodeStatus := m[i]
				mlock.RUnlock()
				nodeStatus.lock.Lock()
				defer nodeStatus.lock.Unlock()
				if nodeStatus.status {
					return
				}
				e := network[i].Start()
				require.NoError(t, e)
				nodeStatus.status = true
				mlock.Lock()
				m[i] = nodeStatus
				mlock.Unlock()
			}()
		}
	}
	// waiting for all go routines to be over
	wg.Wait()
	// start back all nodes
	for _, n := range network {
		err = n.Start()
		require.NoError(t, err)
	}

	// Verify network is not on hold anymore and producing blocks
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

func TestValidatorMigration(t *testing.T) {
	vals, err := Validators(t, 5, "10e18,v,10000,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	network, err := NewNetworkFromValidators(t, vals, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// wait for validator to connect to each other
	// pause a validator
	// turn off validator
	// update ip of a single validator
	// re-activate validator with new ip:
	// activate validator again
	// ensure all validators are connected in full mesh with new ip

	// Setup Bindings
	autonityContract, _ := bindings.NewAutonity(params.AutonityContractAddress, network[1].WsClient)

	transactor, err := bind.NewKeyedTransactorWithChainID(vals[0].TreasuryKey, params.TestChainConfig.ChainID)
	require.NoError(t, err)

	validator, _ := autonityContract.GetValidator(nil, network[0].Address)
	oldEnode, err := enode.ParseV4(validator.Enode)
	require.NoError(t, err)
	// use new port in the newer enode
	newPort := freeport.GetOne(t)
	newEnode := enode.AppendConsensusEndpoint(oldEnode.Host(), strconv.Itoa(newPort), oldEnode.String())

	_, err = autonityContract.PauseValidator(transactor, network[0].Address)
	require.NoError(t, err)

	inCommitee := func(addr common.Address) bool {
		members, err := autonityContract.GetCommittee(nil)
		require.NoError(t, err)
		found := false
		for _, m := range members {
			if m.Addr == addr {
				found = true
				break
			}
		}
		return found
	}

	// wait for few blocks to get the validator out of committee
loop:
	for {
		select {
		case <-time.After(1 * time.Second):
			if !inCommitee(network[0].Address) {
				break loop
			}
		case <-time.After(11 * time.Second): // more than 2 epochs
			require.False(t, inCommitee(network[0].Address), "validator is still in committee after pausing")
		}
	}

	_, err = autonityContract.UpdateEnode(transactor, network[0].Address, newEnode)
	require.NoError(t, err)
	// turn off validator to restart it with new port
	network[0].Close(true)

	err = network.WaitToMineNBlocks(1, 10, false)
	require.NoError(t, err)
	// verify other client only has 3 connection now
	require.Equal(t, 3, network[1].ConsensusServer().PeerCount(), "connection with paused validator is not dropped after it got evicted from committee")

	// start the client again
	network[0].Config.ConsensusP2P.ListenAddr = oldEnode.Host() + ":" + strconv.Itoa(newPort)
	network[0].Start()

	// re-activate validator
	_, err = autonityContract.ActivateValidator(transactor, network[0].Address)
	require.NoError(t, err)

	// wait for few blocks to get the validator in committee
loop1:
	for {
		select {
		case <-time.After(1 * time.Second):
			if inCommitee(network[0].Address) {
				break loop1
			}
		case <-time.After(11 * time.Second): // more than 2 epochs
			require.True(t, inCommitee(network[0].Address), "validator is not in committee after activating")
		}
	}
	// we are in committee again wait for approximately 30 second to connect everyone
	err = network.WaitToMineNBlocks(30, 60, false)
	require.NoError(t, err)

	fullyConnected := false
	tries := 30
loop2:
	for {
		time.Sleep(1 * time.Second)
		if network[0].ConsensusServer().PeerCount() == 4 {
			fullyConnected = true
			break loop2
		}
		tries--
		if tries == 0 {
			break loop2
		}
	}
	// verify we are connected fully again
	require.True(t, fullyConnected)
}

/*
// UNSUPPORTED ON NON UNIX DEV ENV
func updateRlimit() {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Getting Rlimit ", err)
	}
	fmt.Println(rLimit)
	rLimit.Max = 65536
	rLimit.Cur = 65536
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Setting Rlimit ", err)
	}
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Getting Rlimit ", err)
	}
	fmt.Println("Rlimit Final", rLimit)
}
*/

func TestCommitteeSizeChangeMidEpoch(t *testing.T) {
	params.TestAutonityContractConfig.EpochPeriod = 10
	params.DefaultOmissionAccountabilityConfig.Delta = 4
	params.DefaultOmissionAccountabilityConfig.LookbackWindow = 4
	params.TestOracleConfig.VotePeriod = 4

	log.DefaultVerbosity = log.LvlError
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlError, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	vals, err := Validators(t, 20, "10e18,v,10000,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	network, err := NewNetworkFromValidators(t, vals, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	err = network.WaitToMineNBlocks(2, 10, false)
	require.NoError(t, err)

	// wait for validator to connect to each other
	// update the committee Size to increase proposer rewards

	// Setup Bindings
	autonityContract, _ := bindings.NewAutonity(params.AutonityContractAddress, network[0].WsClient)

	transactor, err := bind.NewKeyedTransactorWithChainID(vals[0].NodeKey, params.TestChainConfig.ChainID)
	require.NoError(t, err)

	_, err = autonityContract.SetCommitteeSize(transactor, big.NewInt(2))
	require.NoError(t, err)

	err = network.WaitToMineNBlocks(20, 60, false)
	require.NoError(t, err)
}

// TestLargeNetwork test internally +100 nodes setups. This will be moved later
// in its own cmd package.
func TestLargeNetwork(t *testing.T) {
	t.Skip("only on demand")
	//
	//------ Config section -------
	//

	// DefaultVerbosity will set the log levels for the main components: consensus, eth, blockchain..
	log.DefaultVerbosity = log.LvlError
	//Set the root logger level for everything else.
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlError, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	// Fast epoch to see changes in committee reflected fast
	// too short epoch can cause validators be omission faulty, 20 is relatively safe
	params.TestAutonityContractConfig.EpochPeriod = 20
	params.DefaultOmissionAccountabilityConfig.Delta = 5
	params.DefaultOmissionAccountabilityConfig.LookbackWindow = 10
	params.TestOracleConfig.VotePeriod = 10
	// total peers to be deployed
	const peerCount = 40
	// peers above max committee are participants
	const maxCommittee = 38
	var newMaxCommittee = maxCommittee
	var epoch = params.TestAutonityContractConfig.EpochPeriod

	//
	//----End config -------------
	//

	validators, _ := Validators(t, peerCount, "10e18,v,1000,127.0.0.1:%s,%s,%s,%s")
	network, err := NewInMemoryNetwork(t, validators, true, func(genesis *ccore.Genesis) {
		genesis.Config.AutonityContractConfig.MaxCommitteeSize = maxCommittee
	})
	require.NoError(t, err)
	defer network.Shutdown(t)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	autonityContract, err := bindings.NewAutonity(params.AutonityContractAddress, network[0].WsClient)
	require.NoError(t, err)
	transactOpts, err := bind.NewKeyedTransactorWithChainID(network[0].Key, params.TestChainConfig.ChainID)
	require.NoError(t, err)

	getNetworkState := func() (height uint64, committee *types.Committee) {
		for _, n := range network {
			nodeHeight := n.Eth.BlockChain().CurrentHeader().Number.Uint64()
			if nodeHeight > height {
				height = nodeHeight
				epoch, _ := n.Eth.BlockChain().LatestEpoch()
				committee = epoch.Committee
			}
		}
		return
	}

	inCommittee := func(address common.Address, committee *types.Committee) bool {
		return committee.MemberByAddress(address) != nil
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 2, 1, ' ', 0)
	tickCount := 0
	done := make(chan struct{})
	go func() {
		for range ticker.C {
			tickCount++
			fmt.Println("-------------", tickCount, "---------------")

			select {
			case <-done:
				return
			default:
			}

			if tickCount%int(epoch) == int(epoch-1) {
				newMaxCommittee = 1
				fmt.Println("setting new committee size:", newMaxCommittee)
				autonityContract.SetCommitteeSize(transactOpts, big.NewInt(int64(newMaxCommittee)))
			}
			maxHeight, committee := getNetworkState()
			if committee == nil {
				continue
			}
			fmt.Println("max height:", maxHeight)
			fmt.Println("go routines:", runtime.NumGoroutine())
			var peersBuf strings.Builder
			peersBuf.WriteString("n\t")
			var committeeBuf strings.Builder
			committeeBuf.WriteString("c\t")
			var execCountBuf strings.Builder
			execCountBuf.WriteString("en\t")
			var consCountBuf strings.Builder
			consCountBuf.WriteString("cn\t")
			var stateBuf strings.Builder
			stateBuf.WriteString("b\t")
			for i, node := range network {
				peersBuf.WriteString(strconv.Itoa(i) + "\t")
				execCountBuf.WriteString(strconv.Itoa(node.ExecutionServer().PeerCount()) + "\t")
				consCountBuf.WriteString(strconv.Itoa(node.ConsensusServer().PeerCount()) + "\t")
				stateBuf.WriteString(strconv.Itoa(int(node.Eth.BlockChain().CurrentHeader().Number.Uint64())) + "\t")
				if inCommittee(node.Address, committee) {
					committeeBuf.WriteString("X" + "\t")
				} else {
					committeeBuf.WriteString(" " + "\t")
				}
			}
			fmt.Fprintln(w, peersBuf.String())
			fmt.Fprintln(w, committeeBuf.String())
			fmt.Fprintln(w, execCountBuf.String())
			fmt.Fprintln(w, consCountBuf.String())
			fmt.Fprintln(w, stateBuf.String())
			w.Flush()
		}
	}()

	err = network.WaitToMineNBlocks(120, 180, false)
	close(done)
	require.NoError(t, err)
	_, committee := getNetworkState()
	for _, n := range network {
		if inCommittee(n.Address, committee) {
			// nolint:typecheck // ConsensusServer is defined, linter false positive
			acnCount := n.ConsensusServer().PeerCount()
			require.Equal(t, newMaxCommittee-1, uint64(acnCount))
		}
	}
}

func TestLoad(t *testing.T) {
	t.Skip("only on demand")
	peerCount := 7
	targetTPS := 1000
	validators, _ := Validators(t, peerCount, "10e18,v,1000,127.0.0.1:%s,%s,%s,%s")
	network, err := NewInMemoryNetwork(t, validators, true)
	require.NoError(t, err)
	gasFeeCap := new(big.Int).SetUint64(network[0].EthConfig.Genesis.Config.AutonityContractConfig.MinBaseFee)
	//err = network[0].Eth.StartMining(1)
	require.NoError(t, err)
	fmt.Println(network[0].Eth.BlockChain().CurrentHeader().Number)
	closeCh := make(chan struct{})
	for j := 0; j < peerCount; j++ {
		go func(id int) {
			timer := time.NewTicker(time.Second)
			defer timer.Stop()
			nonce := 0
			for {
				select {
				case <-timer.C:
					for i := 0; i < targetTPS/peerCount; i++ {
						rawTx := types.NewTx(&types.DynamicFeeTx{
							Nonce:     uint64(nonce),
							GasTipCap: common.Big1,
							GasFeeCap: new(big.Int).Mul(gasFeeCap, common.Big2),
							Gas:       21000,
							To:        &common.Address{},
							Value:     common.Big1,
							Data:      nil,
						})
						signed, err := types.SignTx(rawTx, types.LatestSigner(network[0].EthConfig.Genesis.Config), network[id].Key)
						require.NoError(t, err)
						network[id].Eth.TxPool().AddLocal(signed)
						nonce++
					}
				case <-closeCh:
					return
				}
			}
		}(j)
	}
	err = network.WaitToMineNBlocks(1800, 2300, false)
	require.NoError(t, err)
	close(closeCh)
}

func newPrevoteEquivocator(c interfaces.Core) interfaces.Broadcaster {
	return &PrevoteEquivocator{hasEquivocated: false, equivocateAfter: 0, Core: c.(*core.Core)}
}

type PrevoteEquivocator struct {
	hasEquivocated  bool   // will only equivocate once
	equivocateAfter uint64 // will only equivocate after this height
	*core.Core
}

// will send only a single equivocated prevote
func (s *PrevoteEquivocator) Broadcast(msg message.Msg) {
	s.BroadcastAll(msg)
	if s.hasEquivocated || msg.H() <= s.equivocateAfter {
		return
	}
	if _, isPrevote := msg.(*message.Prevote); !isPrevote {
		return
	}
	committee, err := s.Core.Backend().BlockChain().CommitteeByHeight(msg.H())
	if err != nil {
		panic(err)
	}

	self := committee.MemberByAddress(s.Core.Address())
	msgEq := message.NewPrevote(msg.R(), msg.H(), NonNilValue, s.Backend().Sign, self, committee)
	s.Logger().Info("Equivocation simulation", "height", msg.H())
	s.BroadcastAll(msgEq)
	s.hasEquivocated = true
}

func extractBackend(engine consensus.Engine) *backend.Backend {
	back, ok := engine.(*backend.Backend)
	if !ok {
		panic("cannot cast engine to consensus backend")
	}
	return back
}

func TestJailingPersistence(t *testing.T) {

	vals, err := Validators(t, 12, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	// set node that is doing equivocation
	equivocatorIndex := 11
	vals[equivocatorIndex].TendermintServices = &interfaces.Services{Broadcaster: newPrevoteEquivocator}

	customEpochPeriod := uint64(100)
	network, err := NewNetworkFromValidators(t, vals, true, func(genesis *ccore.Genesis) {
		genesis.Config.AutonityContractConfig.EpochPeriod = customEpochPeriod
		genesis.Config.OmissionAccountabilityConfig.InactivityThreshold = 10000 // do not jail for inactivity in this test
	})
	require.NoError(t, err)
	defer network.Shutdown(t)

	// mine some blocks, needs to be greater than the reporting period (currently 20) of the fault detector
	// otherwise the misbehaviour will not be posted on chain
	err = network.WaitToMineNBlocks(40, 60, false)
	require.NoError(t, err)

	// val 11 should have been jailed
	require.True(t, extractBackend(network[0].Eth.Engine()).IsJailed(network[equivocatorIndex].Address))

	// shutdown val 0 and bring it back up, val 11 should still be in the jailed mapping
	err = network[0].Restart()
	require.NoError(t, err)

	time.Sleep(1 * time.Second) // give time to faulty validator watcher routine to run first time

	require.True(t, extractBackend(network[0].Eth.Engine()).IsJailed(network[equivocatorIndex].Address))

	// shutdown node 1 and restart it after the epoch end, the jailed mapping should be cleared
	err = network[1].Close(false)
	require.NoError(t, err)
	network[1].Wait()

	// close the epoch, the jailed mapping should be cleared
	err = network.WaitForHeight(customEpochPeriod, int(customEpochPeriod))
	require.NoError(t, err)

	require.False(t, extractBackend(network[0].Eth.Engine()).IsJailed(network[equivocatorIndex].Address))

	// when restarting node 1 should clear his jailed database according to the epoch change
	err = network[1].Start()
	require.NoError(t, err)
	time.Sleep(1 * time.Second)
	require.False(t, extractBackend(network[1].Eth.Engine()).IsJailed(network[equivocatorIndex].Address))
}

func firstEquivocator(c interfaces.Core) interfaces.Broadcaster {
	return &PrevoteEquivocator{hasEquivocated: false, equivocateAfter: 0, Core: c.(*core.Core)}
}

func secondEquivocator(c interfaces.Core) interfaces.Broadcaster {
	return &PrevoteEquivocator{hasEquivocated: false, equivocateAfter: 50, Core: c.(*core.Core)}
}

func thirdEquivocator(c interfaces.Core) interfaces.Broadcaster {
	return &PrevoteEquivocator{hasEquivocated: false, equivocateAfter: 100, Core: c.(*core.Core)}
}

// test that once we p2p jail 1/3 of voting power, we unjail the oldest jailed validators to avoid losing liveness
func TestJailingRotation(t *testing.T) {
	t.Skip("TODO")
	// 6 nodes --> f = 2
	vals, err := Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	// set nodes that are doing equivocation
	firstEquivocatorIndex := 1
	secondEquivocatorIndex := 2
	thirdEquivocatorIndex := 3
	vals[firstEquivocatorIndex].TendermintServices = &interfaces.Services{Broadcaster: firstEquivocator}
	vals[secondEquivocatorIndex].TendermintServices = &interfaces.Services{Broadcaster: secondEquivocator}
	vals[thirdEquivocatorIndex].TendermintServices = &interfaces.Services{Broadcaster: thirdEquivocator}

	customEpochPeriod := uint64(200)
	network, err := NewNetworkFromValidators(t, vals, true, func(genesis *ccore.Genesis) {
		genesis.Config.AutonityContractConfig.EpochPeriod = customEpochPeriod
		genesis.Config.OmissionAccountabilityConfig.InactivityThreshold = 10000 // do not jail for inactivity in this test
	})
	require.NoError(t, err)
	defer network.Shutdown(t)

	// mine some blocks, first equivocator should get jailed
	err = network.WaitToMineNBlocks(50, int(customEpochPeriod), false)
	require.NoError(t, err)

	require.True(t, extractBackend(network[0].Eth.Engine()).IsJailed(network[firstEquivocatorIndex].Address))

	// mine some more, also the second equivocator should get jailed. the first one should remain jailed
	err = network.WaitToMineNBlocks(50, int(customEpochPeriod), false)
	require.NoError(t, err)

	require.True(t, extractBackend(network[0].Eth.Engine()).IsJailed(network[firstEquivocatorIndex].Address))
	require.True(t, extractBackend(network[0].Eth.Engine()).IsJailed(network[secondEquivocatorIndex].Address))

	// mine some more, third equivocator should get jailed. Second one should remain jailed, while first one should get unjailed
	err = network.WaitToMineNBlocks(50, int(customEpochPeriod), false)
	require.NoError(t, err)

	require.False(t, extractBackend(network[0].Eth.Engine()).IsJailed(network[firstEquivocatorIndex].Address))
	require.True(t, extractBackend(network[0].Eth.Engine()).IsJailed(network[secondEquivocatorIndex].Address))
	require.True(t, extractBackend(network[0].Eth.Engine()).IsJailed(network[thirdEquivocatorIndex].Address))

	// close the epoch, everyone free
	err = network.WaitForHeight(customEpochPeriod, int(customEpochPeriod))
	require.NoError(t, err)

	require.False(t, extractBackend(network[0].Eth.Engine()).IsJailed(network[firstEquivocatorIndex].Address))
	require.False(t, extractBackend(network[0].Eth.Engine()).IsJailed(network[secondEquivocatorIndex].Address))
	require.False(t, extractBackend(network[0].Eth.Engine()).IsJailed(network[thirdEquivocatorIndex].Address))
}
