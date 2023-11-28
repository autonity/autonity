package bft

import (
	"math/big"
)

func Quorum(totalVotingPower *big.Int) *big.Int {
	q := new(big.Int).Mul(big.NewInt(2), totalVotingPower)
	mod := new(big.Int)
	q.DivMod(q, big.NewInt(3), mod)
	if mod.Cmp(big.NewInt(0)) > 0 {
		return q.Add(q, big.NewInt(1))
	}
	return q
}

func F(totalVotingPower *big.Int) *big.Int {
	f, mod := new(big.Int).DivMod(totalVotingPower, big.NewInt(3), new(big.Int))
	if mod.Cmp(big.NewInt(0)) > 0 {
		f.Add(f, big.NewInt(1))
	}
	return f.Sub(f, big.NewInt(1))
}
