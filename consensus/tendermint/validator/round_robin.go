package validator

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
)

func roundRobinProposer(valSet tendermint.ValidatorSet, proposer common.Address, round uint64) tendermint.Validator {
	if valSet.Size() == 0 {
		return nil
	}

	seed := round
	if proposer != (common.Address{}) {
		seed = calcSeed(valSet, proposer, round) + 1
	}

	pick := seed % uint64(valSet.Size())
	return valSet.GetByIndex(pick)
}