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
	GoFloatPrecision = uint(100)
)

// goParams mimics the Params struct in the Inflation Controller contract using go's big.Float
type goParams struct {
	iInit       *big.Float
	iTrans      *big.Float
	aE          *big.Float
	t           *big.Float
	genesisTime *big.Int
}

func newGoParams(p *InflationControllerParams, genesisTime *big.Int) goParams {
	denomination := new(big.Float).SetPrec(GoFloatPrecision).SetInt(params.DecimalFactor)
	return goParams{
		iInit:       new(big.Float).Quo(new(big.Float).SetPrec(GoFloatPrecision).SetInt(p.IInit), denomination),
		iTrans:      new(big.Float).Quo(new(big.Float).SetPrec(GoFloatPrecision).SetInt(p.ITrans), denomination),
		aE:          new(big.Float).Quo(new(big.Float).SetPrec(GoFloatPrecision).SetInt(p.AE), denomination),
		t:           new(big.Float).Quo(new(big.Float).SetPrec(GoFloatPrecision).SetInt(p.T), denomination),
		genesisTime: genesisTime,
	}
}

func (p goParams) calculateSupplyDelta(currentSupply, lastEpochTime, currentTime *big.Int) *big.Int {
	one := new(big.Float).SetPrec(GoFloatPrecision).SetInt64(1)

	t0 := new(big.Float).SetPrec(GoFloatPrecision).SetInt(new(big.Int).Sub(lastEpochTime, p.genesisTime))
	t1 := new(big.Float).SetPrec(GoFloatPrecision).SetInt(new(big.Int).Sub(currentTime, p.genesisTime))

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
	T := big.NewInt(10 * params.SecondsInYear)
	p := &InflationControllerParams{
		IInit:  (*big.Int)(params.DefaultInflationControllerGenesis.IInit),
		ITrans: (*big.Int)(params.DefaultInflationControllerGenesis.ITrans),
		AE:     (*big.Int)(params.DefaultInflationControllerGenesis.Ae),
		T:      (*big.Int)(params.DefaultInflationControllerGenesis.T),
		IPerm:  (*big.Int)(params.DefaultInflationControllerGenesis.IPerm),
	}
	inflationReserve := (*big.Int)(params.TestAutonityContractConfig.InitialInflationReserve)
	genesisTime := r.evm.Context.Time
	goP := newGoParams(p, genesisTime)
	_, _, inflationControllerContract, err := r.deployInflationController(nil, *p)
	require.NoError(r.t, err)
	currentSupply := new(big.Int).Mul(big.NewInt(60_000_000), params.NTNDecimalFactor) // NTN precision is 18
	epochPeriod := big.NewInt(4 * 60 * 60)
	epochCount := new(big.Int).Div(T, epochPeriod)
	r.t.Log("total epoch", epochCount)
	for i := uint64(0); i < epochCount.Uint64(); i++ {
		lastEpochTime := new(big.Int).Add(genesisTime, new(big.Int).Mul(new(big.Int).SetUint64(i), epochPeriod))
		currentEpochTime := new(big.Int).Add(genesisTime, new(big.Int).Mul(new(big.Int).SetUint64(i+1), epochPeriod))
		currentTime := time.Unix(int64(currentEpochTime.Uint64()), 0)
		days := currentTime.Day()
		years := currentTime.Year()

		delta, gasConsumed, err := inflationControllerContract.CalculateSupplyDelta(nil, currentSupply, inflationReserve, lastEpochTime, currentEpochTime)
		require.NoError(r.t, err)
		require.LessOrEqual(r.t, gasConsumed, uint64(30_000))
		inflationReserve.Sub(inflationReserve, delta)
		goDeltaComputation := goP.calculateSupplyDelta(currentSupply, lastEpochTime, currentEpochTime)

		// Compare the go implementation with the solidity one
		diffSolWithGoBasis := new(big.Int).Quo(new(big.Int).Mul(new(big.Int).Sub(goDeltaComputation, delta), big.NewInt(10000)), delta)

		fmt.Println("y:", years, "d:", days, "b:", currentEpochTime, "supply:", currentSupply, "delta:", delta, "delta_ntn:", new(big.Int).Div(delta, params.NTNDecimalFactor), "go:", goDeltaComputation, "diffBpts:", diffSolWithGoBasis)

		currentSupply.Add(currentSupply, delta)
	}
	r.t.Log("final NTN supply", new(big.Int).Div(currentSupply, params.NTNDecimalFactor))
	r.waitNextEpoch()
	r.waitNextEpoch()
	r.waitNextEpoch()
	r.waitNextEpoch()
}
