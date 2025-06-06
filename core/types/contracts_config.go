package types

import "math/big"

type AccountabilityParams struct {
	Range       *big.Int
	Delta       *big.Int
	GracePeriod *big.Int
}

func (ap AccountabilityParams) Equal(other AccountabilityParams) bool {
	if ap.Range.Cmp(other.Range) != 0 {
		return false
	}
	if ap.Delta.Cmp(other.Delta) != 0 {
		return false
	}
	if ap.GracePeriod.Cmp(other.GracePeriod) != 0 {
		return false
	}
	return true
}

type ContractsConfig struct {
	EpochPeriod    *big.Int
	BlockPeriod    *big.Int
	GasLimit       *big.Int
	Accountability AccountabilityParams
	Eip1559        Eip1559Params
}

func (cc *ContractsConfig) Copy() *ContractsConfig {
	if cc == nil {
		return nil
	}
	return &ContractsConfig{
		EpochPeriod: new(big.Int).Set(cc.EpochPeriod),
		BlockPeriod: new(big.Int).Set(cc.BlockPeriod),
		GasLimit:    new(big.Int).Set(cc.GasLimit),
		Accountability: AccountabilityParams{
			Range:       new(big.Int).Set(cc.Accountability.Range),
			Delta:       new(big.Int).Set(cc.Accountability.Delta),
			GracePeriod: new(big.Int).Set(cc.Accountability.GracePeriod),
		},
		Eip1559: Eip1559Params{
			MinBaseFee:               new(big.Int).Set(cc.Eip1559.MinBaseFee),
			BaseFeeChangeDenominator: new(big.Int).Set(cc.Eip1559.BaseFeeChangeDenominator),
			ElasticityMultiplier:     new(big.Int).Set(cc.Eip1559.ElasticityMultiplier),
			GasLimitBoundDivisor:     new(big.Int).Set(cc.Eip1559.GasLimitBoundDivisor),
		},
	}
}

func (cc *ContractsConfig) Equal(cc2 *ContractsConfig) bool {
	if cc.EpochPeriod.Cmp(cc2.EpochPeriod) != 0 {
		return false
	}
	if cc.BlockPeriod.Cmp(cc2.BlockPeriod) != 0 {
		return false
	}
	if cc.GasLimit.Cmp(cc2.GasLimit) != 0 {
		return false
	}

	// accountability params
	if !cc.Accountability.Equal(cc2.Accountability) {
		return false
	}

	// eip1559 params
	if !cc.Eip1559.Equal(cc2.Eip1559) {
		return false
	}
	return true
}
