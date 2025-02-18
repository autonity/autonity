package rawdb

import (
	"github.com/autonity/autonity/autonity/bindings"
	"github.com/autonity/autonity/common"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestIsEqual(t *testing.T) {
	config := &bindings.AutonityClientAwareConfig{
		BlockPeriod: common.Big1,
		EpochPeriod: common.Big256,
		MinBaseFee:  common.Big5,
	}
	config2 := &bindings.AutonityClientAwareConfig{
		BlockPeriod: new(big.Int).Set(config.BlockPeriod),
		EpochPeriod: new(big.Int).Set(config.EpochPeriod),
		MinBaseFee:  new(big.Int).Set(config.MinBaseFee),
	}

	require.True(t, isEqual(config, config))
	require.True(t, isEqual(config, config2))
	config2.BlockPeriod = common.Big4
	require.False(t, isEqual(config, config2))
}

func TestReadWriteContractsConfig(t *testing.T) {
	db := NewMemoryDatabase()

	config := &bindings.AutonityClientAwareConfig{
		BlockPeriod:         common.Big1,
		EpochPeriod:         common.Big256,
		MinBaseFee:          common.Big5,
		AccountabilityDelta: big.NewInt(10),
		AccountabilityRange: big.NewInt(100),
	}
	targetNumber := uint64(0)

	WriteContractsConfig(db, targetNumber, config, false)
	fetchedConfig, blockNumber := ReadContractsConfig(db, targetNumber)
	require.Equal(t, targetNumber, blockNumber)
	require.True(t, isEqual(fetchedConfig, config))

	// rewrite config, should be fine
	WriteContractsConfig(db, targetNumber, config, false)
	fetchedConfig, blockNumber = ReadContractsConfig(db, targetNumber)
	require.Equal(t, targetNumber, blockNumber)
	require.True(t, isEqual(fetchedConfig, config))

	// write two times same config, deduplication should happen
	WriteContractsConfig(db, targetNumber+1, config, false)
	WriteContractsConfig(db, targetNumber+2, config, false)
	fetchedConfig, blockNumber = ReadContractsConfig(db, targetNumber+1)
	require.Equal(t, targetNumber, blockNumber)
	require.True(t, isEqual(fetchedConfig, config))
	fetchedConfig, blockNumber = ReadContractsConfig(db, targetNumber+2)
	require.Equal(t, targetNumber, blockNumber)
	require.True(t, isEqual(fetchedConfig, config))

	// change config and write. No deduplication
	config.BlockPeriod = common.Big5
	WriteContractsConfig(db, targetNumber+3, config, false)
	fetchedConfig, blockNumber = ReadContractsConfig(db, targetNumber+3)
	require.Equal(t, targetNumber+3, blockNumber)
	require.Equal(t, common.Big5.Uint64(), fetchedConfig.BlockPeriod.Uint64())

}
