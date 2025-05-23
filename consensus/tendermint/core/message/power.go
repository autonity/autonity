package message

import "math/big"

// auxiliary data structure to take into account aggregated power of a set of signers
type AggregatedPower struct {
	power   *big.Int
	signers *big.Int // used as bitmap, we do not care about coefficients here, only if a validator is present or not
}

// computes the contribution that a vote/aggregate would bring to Core
func Contribution(aggregatorSigners *big.Int, coreSigners *big.Int) *big.Int {
	notCoreSigners := new(big.Int).Not(coreSigners)
	contribution := notCoreSigners.And(notCoreSigners, aggregatorSigners)
	return contribution
}

func (p *AggregatedPower) Set(index int, power *big.Int) {
	if p.signers.Bit(index) == 1 {
		return
	}

	p.signers.SetBit(p.signers, index, 1)
	p.power.Add(p.power, power)
}

func (p *AggregatedPower) Power() *big.Int {
	return p.power

}
func (p *AggregatedPower) Signers() *big.Int {
	return p.signers
}

func (p *AggregatedPower) Copy() *AggregatedPower {
	return &AggregatedPower{power: new(big.Int).Set(p.power), signers: new(big.Int).Set(p.signers)}
}

func NewAggregatedPower() *AggregatedPower {
	return &AggregatedPower{power: new(big.Int), signers: new(big.Int)}
}
