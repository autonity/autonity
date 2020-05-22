package bft

import (
	"math"
)

func Quorum(totalVotingPower uint64) uint64 {
	return uint64(math.Ceil((2 * float64(totalVotingPower)) / 3.))
}

func F(totalVotingPower uint64) uint64 {
	return uint64(math.Ceil(float64(totalVotingPower)/3)) - 1
}
