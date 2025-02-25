package economic

import (
	"github.com/ALTree/bigfloat"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/params"
	"math/big"
)

var (
	GoFloatPrecision = uint(100)
)

// inflationEngine mimics the Params struct in the Inflation Controller contract using go's big.Float
type inflationEngine struct {
	rateInitial      *big.Float
	rateTransition   *big.Float
	curveComplexity  *big.Float
	transitionPeriod *big.Float
	decayRate        *big.Float
	genesisTime      *big.Int
}

func newInflationEngine(p *autonity.InflationControllerParams, genesisTime *big.Int) inflationEngine {
	denomination := new(big.Float).SetPrec(GoFloatPrecision).SetInt(params.DecimalFactor)
	return inflationEngine{
		rateInitial:      new(big.Float).Quo(new(big.Float).SetPrec(GoFloatPrecision).SetInt(p.InflationRateInitial), denomination),
		rateTransition:   new(big.Float).Quo(new(big.Float).SetPrec(GoFloatPrecision).SetInt(p.InflationRateTransition), denomination),
		curveComplexity:  new(big.Float).Quo(new(big.Float).SetPrec(GoFloatPrecision).SetInt(p.InflationCurveConvexity), denomination),
		transitionPeriod: new(big.Float).Quo(new(big.Float).SetPrec(GoFloatPrecision).SetInt(p.InflationTransitionPeriod), denomination),
		decayRate:        new(big.Float).Quo(new(big.Float).SetPrec(GoFloatPrecision).SetInt(p.InflationReserveDecayRate), denomination),
		genesisTime:      genesisTime,
	}
}

func (p inflationEngine) calculateSupplyDelta(circulatingSupply, inflationReserve, lastEpochTime, currentTime *big.Int) *big.Int {

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
}

func (p inflationEngine) calculateSupplyDeltaTrans(circulatingSupply, lastEpochTime, currentTime *big.Int) *big.Int {
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

func (p inflationEngine) calculateSupplyDeltaPerm(inflationReserve, lastEpochTime, currentTime *big.Int) *big.Int {
	deltaT := new(big.Float).SetPrec(GoFloatPrecision).SetInt(new(big.Int).Sub(currentTime, lastEpochTime))
	factor := new(big.Float).SetPrec(GoFloatPrecision).SetInt(inflationReserve)
	factor.Mul(factor, deltaT)
	factor.Mul(factor, p.decayRate)
	res, _ := factor.Int(nil)
	return res
}
