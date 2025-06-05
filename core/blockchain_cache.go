package core

import (
	"fmt"
	"math/big"

	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
)

// NOTE: callers should not modify the returned values and instead treat them as read-only

func (bc *BlockChain) EpochPeriodByHeight(height uint64) (*big.Int, error) {
	config, err := bc.readContractsConfigByHeight(height)
	if err != nil {
		return nil, err
	}
	return config.EpochPeriod, nil
}

func (bc *BlockChain) GasLimitByHeight(height uint64) (*big.Int, error) {
	config, err := bc.readContractsConfigByHeight(height)
	if err != nil {
		return nil, err
	}
	return config.GasLimit, nil
}

func (bc *BlockChain) AccountabilityParamsByHeight(height uint64) (*types.AccountabilityParams, error) {
	config, err := bc.readContractsConfigByHeight(height)
	if err != nil {
		return nil, err
	}
	return &config.Accountability, nil
}

func (bc *BlockChain) Eip1559ParamsByHeight(height uint64) (*types.Eip1559Params, error) {
	config, err := bc.readContractsConfigByHeight(height)
	if err != nil {
		return nil, err
	}
	return &config.Eip1559, nil
}

func (bc *BlockChain) readContractsConfigByHeight(height uint64) (*types.ContractsConfig, error) {
	// unless genesis, the config used for a block is always specified in the previous block contracts config
	if height > 0 {
		height--
	}
	return bc.readContractsConfigAt(height)
}

func (bc *BlockChain) readContractsConfigAt(height uint64) (*types.ContractsConfig, error) {
	// check the in-memory cache first
	if config, isCached := bc.contractsConfigCache.Get(height); isCached {
		return config.(*types.ContractsConfig), nil
	}

	config, _ := rawdb.ReadContractsConfig(bc.db, height)
	// try to restore config from state if not cached
	if config == nil {
		var err error
		log.Debug("Restoring contracts config from state", "height", height)
		config, err = bc.restoreContractsConfigAt(height)
		if err != nil {
			return nil, fmt.Errorf("failed to restore contract config for %d : %w", height, err)
		}
	}
	bc.contractsConfigCache.Add(height, config)
	return config, nil
}

// retrieves the contract config from state and saves it into the cache
func (bc *BlockChain) restoreContractsConfigAt(height uint64) (*types.ContractsConfig, error) {
	header := bc.GetHeaderByNumber(height)
	if header == nil {
		return nil, fmt.Errorf("header %d not found", height)
	}
	state, err := bc.StateAt(header.Root)
	if err != nil {
		return nil, fmt.Errorf("failed to get state for %d: %w", header.Number.Uint64(), err)
	}
	contractsConfig, err := bc.ProtocolContracts().CallGetClientConfig(state, header)
	if err != nil {
		return nil, fmt.Errorf("failed to get contracts config: %w", err)
	}
	rawdb.WriteContractsConfig(bc.db, header.Number.Uint64(), contractsConfig)
	return contractsConfig, nil
}
