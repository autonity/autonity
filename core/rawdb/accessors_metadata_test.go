package rawdb

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
)

func TestReadWriteContractsConfig(t *testing.T) {
	log.Root().SetHandler(log.StderrHandler)

	db := NewMemoryDatabase()

	config := &types.ContractsConfig{
		BlockPeriod:         common.Big1,
		EpochPeriod:         common.Big256,
		GasLimit:            new(big.Int).SetUint64(20_000_000),
		ClusteringThreshold: new(big.Int).SetUint64(64),
		Accountability: types.AccountabilityParams{
			Range:       big.NewInt(100),
			Delta:       big.NewInt(10),
			GracePeriod: common.Big0,
		},
		Eip1559: types.Eip1559Params{
			MinBaseFee:               common.Big5,
			BaseFeeChangeDenominator: new(big.Int).SetUint64(8),
			ElasticityMultiplier:     new(big.Int).SetUint64(2),
			GasLimitBoundDivisor:     new(big.Int).SetUint64(1024),
		},
	}
	targetNumber := uint64(0)

	// reading should return nil
	fetchedConfig, blockNumber := ReadContractsConfig(db, targetNumber)
	require.Nil(t, fetchedConfig)
	require.Equal(t, uint64(0), blockNumber)

	// write a config and fetch it
	WriteContractsConfig(db, targetNumber, config)
	fetchedConfig, blockNumber = ReadContractsConfig(db, targetNumber)
	require.Equal(t, targetNumber, blockNumber)
	require.True(t, fetchedConfig.Equal(config))

	// rewrite config, should be fine
	WriteContractsConfig(db, targetNumber, config)
	fetchedConfig, blockNumber = ReadContractsConfig(db, targetNumber)
	require.Equal(t, targetNumber, blockNumber)
	require.True(t, fetchedConfig.Equal(config))

	// write two times same config, deduplication should happen
	WriteContractsConfig(db, targetNumber+1, config)
	WriteContractsConfig(db, targetNumber+2, config)
	fetchedConfig, blockNumber = ReadContractsConfig(db, targetNumber+1)
	require.Equal(t, targetNumber, blockNumber)
	require.True(t, fetchedConfig.Equal(config))
	fetchedConfig, blockNumber = ReadContractsConfig(db, targetNumber+2)
	require.Equal(t, targetNumber, blockNumber)
	require.True(t, fetchedConfig.Equal(config))

	// change config and write. No deduplication
	config.BlockPeriod = common.Big5
	WriteContractsConfig(db, targetNumber+3, config)
	fetchedConfig, blockNumber = ReadContractsConfig(db, targetNumber+3)
	require.Equal(t, targetNumber+3, blockNumber)
	require.Equal(t, common.Big5.Uint64(), fetchedConfig.BlockPeriod.Uint64())

	// write to a far away block, should still succeed
	WriteContractsConfig(db, targetNumber+100, config)
	fetchedConfig, blockNumber = ReadContractsConfig(db, targetNumber+100)
	require.Equal(t, targetNumber+100, blockNumber)
	require.True(t, config.Equal(fetchedConfig))

	// add same config after (will be deduplicated)
	WriteContractsConfig(db, targetNumber+101, config)
	fetchedConfig, blockNumber = ReadContractsConfig(db, targetNumber+101)
	require.Equal(t, targetNumber+100, blockNumber)
	require.True(t, config.Equal(fetchedConfig))

	// simulate db corruption by overwriting config at block 100 with empty bytearray
	require.NoError(t, db.Put(contractsConfigKey(targetNumber+100), append(contractsConfigDataPrefix, []byte{}...)))

	// reading that config should fail now
	fetchedConfig, blockNumber = ReadContractsConfig(db, targetNumber+100)
	require.Equal(t, uint64(0), blockNumber)
	require.Nil(t, fetchedConfig)
	fetchedConfig, blockNumber = ReadContractsConfig(db, targetNumber+101)
	require.Equal(t, uint64(0), blockNumber)
	require.Nil(t, fetchedConfig)

	// mess it up even more by changing the prefix
	require.NoError(t, db.Put(contractsConfigKey(targetNumber+100), append([]byte("z"), []byte{}...)))

	fetchedConfig, blockNumber = ReadContractsConfig(db, targetNumber+100)
	require.Equal(t, uint64(0), blockNumber)
	require.Nil(t, fetchedConfig)
	fetchedConfig, blockNumber = ReadContractsConfig(db, targetNumber+101)
	require.Equal(t, uint64(0), blockNumber)
	require.Nil(t, fetchedConfig)

	// write another config after that, should be fine
	WriteContractsConfig(db, targetNumber+102, config)
	fetchedConfig, blockNumber = ReadContractsConfig(db, targetNumber+102)
	require.Equal(t, targetNumber+102, blockNumber)
	require.True(t, config.Equal(fetchedConfig))

	WriteContractsConfig(db, targetNumber+103, config)
	fetchedConfig, blockNumber = ReadContractsConfig(db, targetNumber+103)
	require.Equal(t, targetNumber+102, blockNumber)
	require.True(t, config.Equal(fetchedConfig))

	// corrupt pointer to 102
	require.NoError(t, db.Put(contractsConfigKey(targetNumber+103), append(blockNumberDataPrefix, []byte{}...)))
	fetchedConfig, blockNumber = ReadContractsConfig(db, targetNumber+103)
	require.Equal(t, uint64(0), blockNumber)
	require.Nil(t, fetchedConfig)
}
