package types

import "math/big"

type AccountabilityParams struct {
	Range       *big.Int
	Delta       *big.Int
	GracePeriod *big.Int
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
	if cc.Accountability.Range.Cmp(cc2.Accountability.Range) != 0 {
		return false
	}
	if cc.Accountability.Delta.Cmp(cc2.Accountability.Delta) != 0 {
		return false
	}
	if cc.Accountability.GracePeriod.Cmp(cc2.Accountability.GracePeriod) != 0 {
		return false
	}

	// eip1559 params
	if cc.Eip1559.MinBaseFee.Cmp(cc2.Eip1559.MinBaseFee) != 0 {
		return false
	}
	if cc.Eip1559.BaseFeeChangeDenominator.Cmp(cc2.Eip1559.BaseFeeChangeDenominator) != 0 {
		return false
	}
	if cc.Eip1559.ElasticityMultiplier.Cmp(cc2.Eip1559.ElasticityMultiplier) != 0 {
		return false
	}
	if cc.Eip1559.GasLimitBoundDivisor.Cmp(cc2.Eip1559.GasLimitBoundDivisor) != 0 {
		return false
	}
	return true
}
