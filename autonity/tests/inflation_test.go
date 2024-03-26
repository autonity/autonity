package tests

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ALTree/bigfloat"
	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/params"
)

var (
	DecimalPrecision = int64(18)
	SecondsInYear    = int64(365 * 24 * 60 * 60)
	DecimalFactor    = new(big.Int).Exp(big.NewInt(10), big.NewInt(DecimalPrecision), nil)
	NTNDecimalFactor = new(big.Int).SetUint64(params.Ether)
	GoFloatPrecision = uint(100)
)

// goParams mimics the Params struct in the Inflation Controller contract using go's big.Float
type goParams struct {
	iInit  *big.Float
	iTrans *big.Float
	aE     *big.Float
	t      *big.Float
}

func newGoParams(p *InflationControllerParams) goParams {
	denomination := new(big.Float).SetPrec(GoFloatPrecision).SetInt(DecimalFactor)
	return goParams{
		iInit:  new(big.Float).Quo(new(big.Float).SetPrec(GoFloatPrecision).SetInt(p.IInit), denomination),
		iTrans: new(big.Float).Quo(new(big.Float).SetPrec(GoFloatPrecision).SetInt(p.ITrans), denomination),
		aE:     new(big.Float).Quo(new(big.Float).SetPrec(GoFloatPrecision).SetInt(p.AE), denomination),
		t:      new(big.Float).Quo(new(big.Float).SetPrec(GoFloatPrecision).SetInt(p.T), denomination),
	}
}

func (p goParams) calculateSupplyDelta(currentSupply, lastEpochBlock, currentBlock *big.Int) *big.Int {
	one := new(big.Float).SetPrec(GoFloatPrecision).SetInt64(1)

	t0 := new(big.Float).SetPrec(GoFloatPrecision).SetInt(lastEpochBlock)
	t1 := new(big.Float).SetPrec(GoFloatPrecision).SetInt(currentBlock)

	lExp0 := new(big.Float).SetPrec(GoFloatPrecision).Mul(p.aE, t0)
	lExp0.Quo(lExp0, p.t)

	lExp1 := new(big.Float).SetPrec(GoFloatPrecision).Mul(p.aE, t1)
	lExp1.Quo(lExp1, p.t)

	deltaT := new(big.Float).SetPrec(GoFloatPrecision).Sub(t1, t0)
	expTerm1 := new(big.Float).SetPrec(GoFloatPrecision).Mul(p.iInit, deltaT)

	expTerm2 := new(big.Float).SetPrec(GoFloatPrecision).Sub(p.iInit, p.iTrans)
	expTerm2.Mul(expTerm2, deltaT)
	aEExp := bigfloat.Exp(p.aE)
	temp3 := new(big.Float).SetPrec(GoFloatPrecision).Sub(aEExp, one)
	expTerm2.Quo(expTerm2, temp3)

	expTerm1.Add(expTerm1, expTerm2)

	temp4 := bigfloat.Exp(lExp1)
	temp5 := bigfloat.Exp(lExp0)
	temp4.Sub(temp4, temp5)
	temp4.Mul(temp4, p.t)
	temp6 := new(big.Float).SetPrec(GoFloatPrecision).Sub(p.iTrans, p.iInit)
	temp4.Mul(temp4, temp6)
	temp3.Mul(temp3, p.aE)
	temp4.Quo(temp4, temp3)

	expTerm1.Add(expTerm1, temp4)
	expTerm1 = bigfloat.Exp(expTerm1)

	currentSupplyFloat := new(big.Float).SetPrec(GoFloatPrecision).SetInt(currentSupply)
	expTerm1.Mul(expTerm1, currentSupplyFloat)
	expTerm1.Sub(expTerm1, currentSupplyFloat)

	res, _ := expTerm1.Int(nil)
	return res
}

// todo: use approx formula from faical
//func (p goParams) calculateSupplyDeltaApprox() *big.Int {

//}

func TestInflationContract(t *testing.T) {
	r := setup(t, nil)
	T := big.NewInt(10 * SecondsInYear)
	p := &InflationControllerParams{
		IInit:  new(big.Int).Div(new(big.Int).Mul(big.NewInt(10), DecimalFactor), big.NewInt(100*SecondsInYear)),
		ITrans: new(big.Int).Div(new(big.Int).Mul(big.NewInt(6), DecimalFactor), big.NewInt(100*SecondsInYear)),
		AE:     new(big.Int).Div(new(big.Int).Mul(big.NewInt(-124_267), DecimalFactor), big.NewInt(100_000)),
		T:      new(big.Int).Mul(big.NewInt(3*SecondsInYear), DecimalFactor),
	}
	goP := newGoParams(p)
	_, _, inflationControllerContract, err := r.deployInflationController(nil, *p)
	require.NoError(r.t, err)
	currentSupply := new(big.Int).Mul(big.NewInt(60_000_000), NTNDecimalFactor) // NTN precision is 18
	epochPeriod := big.NewInt(4 * 60 * 60)
	epochCount := new(big.Int).Div(T, epochPeriod)
	r.t.Log("total epoch", epochCount)
	for i := uint64(0); i < epochCount.Uint64(); i++ {
		lastEpochBlock := new(big.Int).Mul(new(big.Int).SetUint64(i), epochPeriod)
		currentBlock := new(big.Int).Mul(new(big.Int).SetUint64(i+1), epochPeriod)
		currentTime := time.Duration(currentBlock.Uint64() * 1e9)
		days := (currentTime / (time.Hour * 24)) % 365
		years := currentTime / (time.Hour * 24 * 365)

		delta, gasConsumed, err := inflationControllerContract.CalculateSupplyDelta(nil, currentSupply, lastEpochBlock, currentBlock)
		require.LessOrEqual(r.t, gasConsumed, uint64(30_000))
		goDeltaComputation := goP.calculateSupplyDelta(currentSupply, lastEpochBlock, currentBlock)

		// Compare the go implementation with the solidity one
		diffSolWithGoBasis := new(big.Int).Quo(new(big.Int).Mul(new(big.Int).Sub(goDeltaComputation, delta), big.NewInt(10000)), delta)

		fmt.Println("y:", int(years), "d:", int(days), "b:", currentBlock, "supply:", currentSupply, "delta:", delta, "delta_ntn:", new(big.Int).Div(delta, NTNDecimalFactor), "go:", goDeltaComputation, "diffBpts:", diffSolWithGoBasis)
		require.NoError(r.t, err)
		currentSupply.Add(currentSupply, delta)
	}
	r.t.Log("final NTN supply", new(big.Int).Div(currentSupply, NTNDecimalFactor))
}
