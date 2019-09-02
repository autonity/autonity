package validator

import (
	"github.com/clearmatics/autonity/common"
)

func stickyProposer(valSet Set, proposer common.Address, round uint64) Validator {
	size := valSet.Size()
	if size == 0 {
		return nil
	}

	seed := round
	if proposer != (common.Address{}) {
		seed = calcSeed(valSet, proposer, round)
	}

	pick := seed % uint64(size)
	return valSet.GetByIndex(pick)
}

func calcSeed(valSet Set, proposer common.Address, round uint64) uint64 {
	offset := 0
	if idx, val := valSet.GetByAddress(proposer); val != nil {
		offset = idx
	}
	return uint64(offset) + round
}
