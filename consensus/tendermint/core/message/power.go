package message

import "math/big"

// auxiliary data structure to take into account aggregated power of a set of signers
type AggregatedPower struct {
	power *big.Int
	//todo:(review) => integration with new signer
	signers *big.Int // used as bitmap, we do not care about coefficients here, only if a validator is present or not
}

// computes the contribution that a vote/aggregate would bring to Core
func Contribution(aggregatorSigners *big.Int, coreSigners *big.Int) *big.Int {
	notCoreSigners := new(big.Int).Not(coreSigners)
	contribution := notCoreSigners.And(notCoreSigners, aggregatorSigners)
	return contribution
}

// returns whether the new signer increased the power or was redundant
func (p *AggregatedPower) Set(index int, power *big.Int) bool {
	if p.signers.Bit(index) == 1 {
		return false // no power increase, the signer was already included
	}

	p.signers.SetBit(p.signers, index, 1)
	p.power.Add(p.power, power)
	return true
}

func (p *AggregatedPower) Subtract(other *AggregatedPower) bool {
	p.power.Sub(p.power, other.power)
	mask := new(big.Int).Not(other.signers)
	p.signers.And(other.signers, mask)
	return true
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
