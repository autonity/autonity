package validator

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
)

func stickyProposer(valSet tendermint.ValidatorSet, proposer common.Address, round uint64) tendermint.Validator {
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

func calcSeed(valSet tendermint.ValidatorSet, proposer common.Address, round uint64) uint64 {
	offset := 0
	if idx, val := valSet.GetByAddress(proposer); val != nil {
		offset = idx
	}
	return uint64(offset) + round
}
