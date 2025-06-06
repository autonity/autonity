// Copyright 2021 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package misc

import (
	"fmt"
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/math"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/params"
)

// VerifyEip1559Header verifies some header attributes which were changed in EIP-1559,
// - gas limit check
// - basefee check
func VerifyEip1559Header(config *params.ChainConfig, eip1559Params *types.Eip1559Params, parent, header *types.Header) error {
	// Verify that the gas limit remains within allowed bounds
	parentGasLimit := parent.GasLimit
	if !config.IsLondon(parent.Number) {
		// take the genesis elasticity multiplier, however this
		// should never happen on Autonity, since EIP1559 is activated from genesis
		parentGasLimit = parent.GasLimit * params.DefaultElasticityMultiplier
	}
	if err := VerifyGaslimit(parentGasLimit, header.GasLimit, eip1559Params.GasLimitBoundDivisor.Uint64()); err != nil {
		return err
	}
	// Verify the header is not malformed
	if header.BaseFee == nil {
		return fmt.Errorf("header is missing baseFee")
	}
	// Verify the baseFee is correct based on the parent header.
	expectedBaseFee := CalcBaseFee(config, parent, eip1559Params)
	if header.BaseFee.Cmp(expectedBaseFee) != 0 {
		return fmt.Errorf("invalid baseFee: have %s, want %s, parentBaseFee %s, parentGasUsed %d",
			header.BaseFee, expectedBaseFee, parent.BaseFee, parent.GasUsed)
	}
	return nil
}

// CalcBaseFee calculates the basefee of the header.
func CalcBaseFee(config *params.ChainConfig, parent *types.Header, eip1559Params *types.Eip1559Params) *big.Int {
	// If the current block is the first EIP-1559 block, return the InitialBaseFee.
	if !config.IsLondon(parent.Number) {
		return new(big.Int).SetUint64(params.InitialBaseFee)
	}

	// compute the correct gas target
	var parentGasTarget uint64
	if parent.GasLimit >= eip1559Params.ElasticityMultiplier.Uint64() {
		// standard case
		parentGasTarget = parent.GasLimit / eip1559Params.ElasticityMultiplier.Uint64()
	} else {
		// extreme edge case, let's handle it gracefully
		parentGasTarget = 1 // avoid targetting 0, which would cause division by 0 later on
	}
	parentGasTargetBig := new(big.Int).SetUint64(parentGasTarget)

	// If the parent gasUsed is the same as the target, the baseFee remains unchanged.
	if parent.GasUsed == parentGasTarget {
		return new(big.Int).Set(parent.BaseFee)
	}
	if parent.GasUsed > parentGasTarget {
		// If the parent block used more gas than its target, the baseFee should increase.
		gasUsedDelta := new(big.Int).SetUint64(parent.GasUsed - parentGasTarget)
		x := new(big.Int).Mul(parent.BaseFee, gasUsedDelta)
		y := x.Div(x, parentGasTargetBig)
		baseFeeDelta := math.BigMax(
			x.Div(y, eip1559Params.BaseFeeChangeDenominator),
			common.Big1,
		)
		return x.Add(parent.BaseFee, baseFeeDelta)
	} else {
		// Otherwise if the parent block used less gas than its target, the baseFee should decrease.
		gasUsedDelta := new(big.Int).SetUint64(parentGasTarget - parent.GasUsed)
		x := new(big.Int).Mul(parent.BaseFee, gasUsedDelta)
		y := x.Div(x, parentGasTargetBig)
		baseFeeDelta := x.Div(y, eip1559Params.BaseFeeChangeDenominator)

		// if the baseFeeDelta is 0  due to a very high baseFeeChangeDenominator (and not due to parent.BaseFee == 0),
		// bump it to 1, as we want the basefee to decrease if we are below target
		if parent.BaseFee.Uint64() > 0 && baseFeeDelta.Uint64() == 0 {
			baseFeeDelta = common.Big1
		}

		return math.BigMax(
			x.Sub(parent.BaseFee, baseFeeDelta),
			new(big.Int).Set(eip1559Params.MinBaseFee),
		)
	}
}
