package committee

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
)

func stickyProposer(valSet Set, proposer common.Address, round int64) types.CommitteeMember {
	size := valSet.Size()
	seed := int(round)
	if proposer != (common.Address{}) {
		seed = calcSeed(valSet, proposer, round)
	}

	pick := seed % size
	return valSet.GetByIndex(pick)
}

func calcSeed(valSet Set, proposer common.Address, round int64) int {
	offset := 0
	if idx, _ := valSet.GetByAddress(proposer); idx != -1 {
		offset = idx
	}
	return offset + int(round)
}
