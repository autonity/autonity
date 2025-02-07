package core

import (
	"encoding/binary"
	"math/big"
	"testing"

	"github.com/autonity/autonity/consensus/ethash"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/ethdb"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
	"github.com/stretchr/testify/require"
)

var ethHashConfig = &types.ContractsConfig{
	EpochPeriod: new(big.Int).SetUint64(params.TestChainConfig.AutonityContractConfig.EpochPeriod),
	BlockPeriod: new(big.Int).SetUint64(params.TestChainConfig.AutonityContractConfig.BlockPeriod),
	GasLimit:    new(big.Int).SetUint64(params.TestChainConfig.AutonityContractConfig.GasLimit),
	Accountability: types.AccountabilityParams{
		Range:       new(big.Int).SetUint64(params.TestChainConfig.AccountabilityConfig.Range),
		Delta:       new(big.Int).SetUint64(params.TestChainConfig.AccountabilityConfig.Delta),
		GracePeriod: new(big.Int),
	},
	Eip1559: types.Eip1559Params{
		MinBaseFee:               new(big.Int).SetUint64(params.TestMinBaseFee),
		BaseFeeChangeDenominator: new(big.Int).SetUint64(params.DefaultBaseFeeChangeDenominator),
		ElasticityMultiplier:     new(big.Int).SetUint64(params.DefaultElasticityMultiplier),
		GasLimitBoundDivisor:     new(big.Int).SetUint64(params.DefaultGasLimitBoundDivisor),
	},
}

// helper function to arbitrarly assign slots in cache
func tamperCache(t *testing.T, db ethdb.Database, number uint64, tamperedData []byte) {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)

	contractsConfigPrefix := []byte("C")
	err := db.Put(append(contractsConfigPrefix, enc...), tamperedData)
	require.NoError(t, err)
}

func TestBlockchainCache(t *testing.T) {
	log.Root().SetHandler(log.StderrHandler)

	blockToBeMined := 100
	targetBlock := uint64(blockToBeMined / 2)

	db, chain, err := newCanonical(t, ethash.NewFaker(), blockToBeMined, true)
	require.NoError(t, err)

	// all block should store the fake ethhash config
	eip1559Params, err := chain.Eip1559ParamsByHeight(targetBlock)
	require.NoError(t, err)
	require.Equal(t, ethHashConfig.Eip1559.MinBaseFee.String(), eip1559Params.MinBaseFee.String())

	epochPeriod, err := chain.EpochPeriodByHeight(targetBlock + 1)
	require.NoError(t, err)
	require.Equal(t, ethHashConfig.EpochPeriod.String(), epochPeriod.String())

	accountabilityParams, err := chain.AccountabilityParamsByHeight(targetBlock - 1)
	require.NoError(t, err)
	require.Equal(t, ethHashConfig.Accountability.Delta.String(), accountabilityParams.Delta.String())

	accountabilityParams, err = chain.AccountabilityParamsByHeight(targetBlock - 10)
	require.NoError(t, err)
	require.Equal(t, ethHashConfig.Accountability.Range.String(), accountabilityParams.Range.String())

	accountabilityParams, err = chain.AccountabilityParamsByHeight(targetBlock + 10)
	require.NoError(t, err)
	require.Equal(t, ethHashConfig.Accountability.GracePeriod.String(), accountabilityParams.GracePeriod.String())

	eip1559Params, err = chain.Eip1559ParamsByHeight(targetBlock)
	require.NoError(t, err)
	require.Equal(t, ethHashConfig.Eip1559.BaseFeeChangeDenominator.String(), eip1559Params.BaseFeeChangeDenominator.String())

	eip1559Params, err = chain.Eip1559ParamsByHeight(targetBlock + 20)
	require.NoError(t, err)
	require.Equal(t, ethHashConfig.Eip1559.ElasticityMultiplier.String(), eip1559Params.ElasticityMultiplier.String())

	gasLimit, err := chain.GasLimitByHeight(targetBlock + 30)
	require.NoError(t, err)
	require.Equal(t, ethHashConfig.GasLimit.String(), gasLimit.String())

	eip1559Params, err = chain.Eip1559ParamsByHeight(0)
	require.NoError(t, err)
	require.Equal(t, ethHashConfig.Eip1559.GasLimitBoundDivisor.String(), eip1559Params.GasLimitBoundDivisor.String())

	// check that recover from state works correctly
	t.Logf("tampering %d slot", targetBlock-1)
	tamperCache(t, db, targetBlock-1, []byte{})

	eip1559Params, err = chain.Eip1559ParamsByHeight(targetBlock)
	require.NoError(t, err)
	require.Equal(t, ethHashConfig.Eip1559.GasLimitBoundDivisor.String(), eip1559Params.GasLimitBoundDivisor.String())

	t.Logf("tampering %d slot", 0)
	tamperCache(t, db, 0, []byte{})

	eip1559Params, err = chain.Eip1559ParamsByHeight(targetBlock)
	require.NoError(t, err)
	require.Equal(t, ethHashConfig.Eip1559.GasLimitBoundDivisor.String(), eip1559Params.GasLimitBoundDivisor.String())

	// delete all cached info. Recovery from state should still work
	t.Logf("nuke everything 💀")
	for i := 0; i <= blockToBeMined; i++ {
		tamperCache(t, db, uint64(i), []byte{})
	}

	gasLimit, err = chain.GasLimitByHeight(targetBlock + 30)
	require.NoError(t, err)
	require.Equal(t, ethHashConfig.GasLimit.String(), gasLimit.String())

	t.Logf("fetching full config at %d", targetBlock)
	config, err := chain.readContractsConfigAt(targetBlock)
	require.NoError(t, err)
	require.Equal(t, ethHashConfig.Eip1559.MinBaseFee.String(), config.Eip1559.MinBaseFee.String())
	require.Equal(t, ethHashConfig.EpochPeriod.String(), config.EpochPeriod.String())
	require.Equal(t, ethHashConfig.Accountability.Delta.String(), config.Accountability.Delta.String())
	require.Equal(t, ethHashConfig.Accountability.Range.String(), config.Accountability.Range.String())
	require.Equal(t, ethHashConfig.Accountability.GracePeriod.String(), config.Accountability.GracePeriod.String())
	require.Equal(t, ethHashConfig.Eip1559.BaseFeeChangeDenominator.String(), config.Eip1559.BaseFeeChangeDenominator.String())
	require.Equal(t, ethHashConfig.Eip1559.ElasticityMultiplier.String(), config.Eip1559.ElasticityMultiplier.String())
	require.Equal(t, ethHashConfig.GasLimit.String(), config.GasLimit.String())
	require.Equal(t, ethHashConfig.Eip1559.GasLimitBoundDivisor.String(), config.Eip1559.GasLimitBoundDivisor.String())

	t.Logf("restoring full config at %d", targetBlock)
	config, err = chain.restoreContractsConfigAt(targetBlock)
	require.NoError(t, err)
	require.Equal(t, ethHashConfig.Eip1559.MinBaseFee.String(), config.Eip1559.MinBaseFee.String())
	require.Equal(t, ethHashConfig.EpochPeriod.String(), config.EpochPeriod.String())
	require.Equal(t, ethHashConfig.Accountability.Delta.String(), config.Accountability.Delta.String())
	require.Equal(t, ethHashConfig.Accountability.Range.String(), config.Accountability.Range.String())
	require.Equal(t, ethHashConfig.Accountability.GracePeriod.String(), config.Accountability.GracePeriod.String())
	require.Equal(t, ethHashConfig.Eip1559.BaseFeeChangeDenominator.String(), config.Eip1559.BaseFeeChangeDenominator.String())
	require.Equal(t, ethHashConfig.Eip1559.ElasticityMultiplier.String(), config.Eip1559.ElasticityMultiplier.String())
	require.Equal(t, ethHashConfig.GasLimit.String(), config.GasLimit.String())
	require.Equal(t, ethHashConfig.Eip1559.GasLimitBoundDivisor.String(), config.Eip1559.GasLimitBoundDivisor.String())

	// trying to fetch config in the future fails
	t.Logf("trying to fetch at %d", targetBlock+1000)
	eip1559Params, err = chain.Eip1559ParamsByHeight(targetBlock + 1000)
	require.Error(t, err)
	t.Log(err)
	require.Nil(t, eip1559Params)

	// fetching for pending block is fine
	t.Logf("fetching at %d", blockToBeMined+1)
	gasLimit, err = chain.GasLimitByHeight(uint64(blockToBeMined + 1))
	require.NoError(t, err)
	require.Equal(t, ethHashConfig.GasLimit.String(), gasLimit.String())
}
