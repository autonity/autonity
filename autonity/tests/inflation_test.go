package tests

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ALTree/bigfloat"
	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/params"
)

var (
	GoFloatPrecision = uint(100)
)

// goParams mimics the Params struct in the Inflation Controller contract using go's big.Float
type goParams struct {
	rateInitial      *big.Float
	rateTransition   *big.Float
	curveComplexity  *big.Float
	transitionPeriod *big.Float
	decayRate        *big.Float
	genesisTime      *big.Int
}

func newGoParams(p *InflationControllerParams, genesisTime *big.Int) goParams {
	denomination := new(big.Float).SetPrec(GoFloatPrecision).SetInt(params.DecimalFactor)
	return goParams{
		rateInitial:      new(big.Float).Quo(new(big.Float).SetPrec(GoFloatPrecision).SetInt(p.InflationRateInitial), denomination),
		rateTransition:   new(big.Float).Quo(new(big.Float).SetPrec(GoFloatPrecision).SetInt(p.InflationRateTransition), denomination),
		curveComplexity:  new(big.Float).Quo(new(big.Float).SetPrec(GoFloatPrecision).SetInt(p.InflationCurveConvexity), denomination),
		transitionPeriod: new(big.Float).Quo(new(big.Float).SetPrec(GoFloatPrecision).SetInt(p.InflationTransitionPeriod), denomination),
		decayRate:        new(big.Float).Quo(new(big.Float).SetPrec(GoFloatPrecision).SetInt(p.InflationReserveDecayRate), denomination),
		genesisTime:      genesisTime,
	}
}

func (p goParams) calculateSupplyDelta(circulatingSupply, inflationReserve, lastEpochTime, currentTime *big.Int) *big.Int {

	t0 := new(big.Int).Sub(lastEpochTime, p.genesisTime)
	t1 := new(big.Int).Sub(currentTime, p.genesisTime)

	if new(big.Float).SetInt(t1).Cmp(p.transitionPeriod) <= 0 {
		return p.calculateSupplyDeltaTrans(circulatingSupply, t0, t1)
	}

	// t1 > p.t from here
	if new(big.Float).SetInt(t0).Cmp(p.transitionPeriod) < 0 {
		pT, _ := p.transitionPeriod.Int(nil)
		untilT := p.calculateSupplyDeltaTrans(circulatingSupply, t0, pT)
		afterT := p.calculateSupplyDeltaPerm(inflationReserve, pT, t1)
		return new(big.Int).Add(untilT, afterT)
	}
	return p.calculateSupplyDeltaPerm(inflationReserve, t0, t1)

	// lExp0 := new(big.Float).SetPrec(GoFloatPrecision).Mul(p.aE, t0)
	// lExp0.Quo(lExp0, p.t)

	// lExp1 := new(big.Float).SetPrec(GoFloatPrecision).Mul(p.aE, t1)
	// lExp1.Quo(lExp1, p.t)

	// deltaT := new(big.Float).SetPrec(GoFloatPrecision).Sub(t1, t0)
	// expTerm1 := new(big.Float).SetPrec(GoFloatPrecision).Mul(p.iInit, deltaT)

	// expTerm2 := new(big.Float).SetPrec(GoFloatPrecision).Sub(p.iInit, p.iTrans)
	// expTerm2.Mul(expTerm2, deltaT)
	// aEExp := bigfloat.Exp(p.aE)
	// temp3 := new(big.Float).SetPrec(GoFloatPrecision).Sub(aEExp, one)
	// expTerm2.Quo(expTerm2, temp3)

	// expTerm1.Add(expTerm1, expTerm2)

	// temp4 := bigfloat.Exp(lExp1)
	// temp5 := bigfloat.Exp(lExp0)
	// temp4.Sub(temp4, temp5)
	// temp4.Mul(temp4, p.t)
	// temp6 := new(big.Float).SetPrec(GoFloatPrecision).Sub(p.iTrans, p.iInit)
	// temp4.Mul(temp4, temp6)
	// temp3.Mul(temp3, p.aE)
	// temp4.Quo(temp4, temp3)

	// expTerm1.Add(expTerm1, temp4)
	// expTerm1 = bigfloat.Exp(expTerm1)

	// circulatingSupplyFloat := new(big.Float).SetPrec(GoFloatPrecision).SetInt(circulatingSupply)
	// expTerm1.Mul(expTerm1, circulatingSupplyFloat)
	// expTerm1.Sub(expTerm1, circulatingSupplyFloat)

	// res, _ := expTerm1.Int(nil)
	// return res
}

func (p goParams) calculateSupplyDeltaTrans(circulatingSupply, lastEpochTime, currentTime *big.Int) *big.Int {
	one := new(big.Float).SetPrec(GoFloatPrecision).SetInt64(1)

	t0 := new(big.Float).SetPrec(GoFloatPrecision).SetInt(lastEpochTime)
	t1 := new(big.Float).SetPrec(GoFloatPrecision).SetInt(currentTime)

	lExp0 := new(big.Float).SetPrec(GoFloatPrecision).Mul(p.curveComplexity, t0)
	lExp0.Quo(lExp0, p.transitionPeriod)

	lExp1 := new(big.Float).SetPrec(GoFloatPrecision).Mul(p.curveComplexity, t1)
	lExp1.Quo(lExp1, p.transitionPeriod)

	deltaT := new(big.Float).SetPrec(GoFloatPrecision).Sub(t1, t0)
	expTerm1 := new(big.Float).SetPrec(GoFloatPrecision).Mul(p.rateInitial, deltaT)

	expTerm2 := new(big.Float).SetPrec(GoFloatPrecision).Sub(p.rateInitial, p.rateTransition)
	expTerm2.Mul(expTerm2, deltaT)
	aEExp := bigfloat.Exp(p.curveComplexity)
	temp3 := new(big.Float).SetPrec(GoFloatPrecision).Sub(aEExp, one)
	expTerm2.Quo(expTerm2, temp3)

	expTerm1.Add(expTerm1, expTerm2)

	temp4 := bigfloat.Exp(lExp1)
	temp5 := bigfloat.Exp(lExp0)
	temp4.Sub(temp4, temp5)
	temp4.Mul(temp4, p.transitionPeriod)
	temp6 := new(big.Float).SetPrec(GoFloatPrecision).Sub(p.rateTransition, p.rateInitial)
	temp4.Mul(temp4, temp6)
	temp3.Mul(temp3, p.curveComplexity)
	temp4.Quo(temp4, temp3)

	expTerm1.Add(expTerm1, temp4)
	expTerm1 = bigfloat.Exp(expTerm1)

	circulatingSupplyFloat := new(big.Float).SetPrec(GoFloatPrecision).SetInt(circulatingSupply)
	expTerm1.Mul(expTerm1, circulatingSupplyFloat)
	expTerm1.Sub(expTerm1, circulatingSupplyFloat)

	res, _ := expTerm1.Int(nil)
	return res
}

func (p goParams) calculateSupplyDeltaPerm(inflationReserve, lastEpochTime, currentTime *big.Int) *big.Int {
	deltaT := new(big.Float).SetPrec(GoFloatPrecision).SetInt(new(big.Int).Sub(currentTime, lastEpochTime))
	factor := new(big.Float).SetPrec(GoFloatPrecision).SetInt(inflationReserve)
	factor.Mul(factor, deltaT)
	factor.Mul(factor, p.decayRate)
	res, _ := factor.Int(nil)
	return res
}

// todo: use approx formula from faical
//func (p goParams) calculateSupplyDeltaApprox() *big.Int {

//}

func TestInflationContract(t *testing.T) {
	r := Setup(t, nil)
	T := big.NewInt(10 * params.SecondsInYear)
	p := &InflationControllerParams{
		InflationRateInitial:      (*big.Int)(params.DefaultInflationControllerGenesis.InflationRateInitial),
		InflationRateTransition:   (*big.Int)(params.DefaultInflationControllerGenesis.InflationRateTransition),
		InflationCurveConvexity:   (*big.Int)(params.DefaultInflationControllerGenesis.InflationCurveConvexity),
		InflationTransitionPeriod: (*big.Int)(params.DefaultInflationControllerGenesis.InflationTransitionPeriod),
		InflationReserveDecayRate: (*big.Int)(params.DefaultInflationControllerGenesis.InflationReserveDecayRate),
	}
	inflationReserve := (*big.Int)(params.TestAutonityContractConfig.InitialInflationReserve)
	genesisTime := r.Evm.Context.Time
	goP := newGoParams(p, genesisTime)
	_, _, inflationControllerContract, err := r.DeployInflationController(nil, *p)
	require.NoError(r.T, err)
	circulatingSupply := new(big.Int).Mul(big.NewInt(60_000_000), params.NTNDecimalFactor) // NTN precision is 18
	epochPeriod := big.NewInt(4 * 60 * 60)
	epochCount := new(big.Int).Div(T, epochPeriod)
	r.T.Log("total epoch", epochCount)
	for i := uint64(0); i < epochCount.Uint64(); i++ {
		lastEpochTime := new(big.Int).Add(genesisTime, new(big.Int).Mul(new(big.Int).SetUint64(i), epochPeriod))
		currentEpochTime := new(big.Int).Add(genesisTime, new(big.Int).Mul(new(big.Int).SetUint64(i+1), epochPeriod))
		currentTime := time.Unix(int64(currentEpochTime.Uint64()), 0)
		days := currentTime.Day()
		years := currentTime.Year()

		delta, gasConsumed, err := inflationControllerContract.CalculateSupplyDelta(nil, circulatingSupply, inflationReserve, lastEpochTime, currentEpochTime)
		require.NoError(r.T, err)
		require.LessOrEqual(r.T, gasConsumed, uint64(30_000))
		goDeltaComputation := goP.calculateSupplyDelta(circulatingSupply, inflationReserve, lastEpochTime, currentEpochTime)
		inflationReserve.Sub(inflationReserve, delta)

		// Compare the go implementation with the solidity one
		diffSolWithGoBasis := new(big.Int).Quo(new(big.Int).Mul(new(big.Int).Sub(goDeltaComputation, delta), big.NewInt(10000)), delta)

		fmt.Println("y:", years, "d:", days, "b:", currentEpochTime, "supply:", circulatingSupply, "delta:", delta, "delta_ntn:", new(big.Int).Div(delta, params.NTNDecimalFactor), "go:", goDeltaComputation, "diffBpts:", diffSolWithGoBasis)
		require.True(r.T, diffSolWithGoBasis.Cmp(common.Big0) == 0, "inflation reward calculation mismatch")

		circulatingSupply.Add(circulatingSupply, delta)
	}
	r.T.Log("final NTN supply", new(big.Int).Div(circulatingSupply, params.NTNDecimalFactor))
}
