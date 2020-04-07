package committee

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"math/big"
)

// weightedRoundRobinProposer, it distribute the height+round into a 256bits value space by Keccak256, then schedule
// the proposer weighted by voting power.
func weightedRoundRobinProposer(valSet Set, proposer common.Address, round int64, height *big.Int) types.CommitteeMember {

	// if total voting power turn to 0, fall back to round robin.
	totalPower := uint64(0)
	for _, val := range valSet.Committee() {
		totalPower += val.VotingPower.Uint64()
	}

	if totalPower == 0 {
		return roundRobinProposer(valSet, proposer, round, height)
	}

	// distributed the round into 256bit value space by keccak256.
	inputSeed := new(big.Int).Add(height, new(big.Int).SetInt64(round))
	seed := new(big.Int).SetBytes(crypto.Keccak256([]byte(inputSeed.String())))

	// weighted by voting power.
	offset := seed.Uint64() % totalPower
	selectedIndex := 0

	counter := uint64(0)
	for i:= 0; i < valSet.Committee().Len(); i ++ {
		if valSet.Committee()[i].VotingPower.Uint64() == 0 {
			continue
		}
		counter += valSet.Committee()[i].VotingPower.Uint64()
		if offset <= (counter - 1) {
			selectedIndex = i
			break
		}
	}
	selectedProposer, _ := valSet.GetByIndex(selectedIndex)
	return selectedProposer
}
