package temp

import (
	"errors"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"github.com/clearmatics/autonity/core/types"
)

// retrieve list of committee for the block at height passed as parameter
func SavedCommittee(number uint64, chain consensus.ChainReader) (committee.Set, error) {
	var lastProposer common.Address
	var err error
	if number == 0 {
		number = 1
	}
	parentHeader := chain.GetHeaderByNumber(number - 1)
	if parentHeader == nil {
		return nil, errors.New("unknown block")
	}
	// For the genesis block, lastProposer is no one (empty).
	if number > 1 {
		lastProposer, err = types.Ecrecover(parentHeader)
		if err != nil {
			return nil, err
		}
	}
	return committee.NewRoundRobinSet(parentHeader.Committee, lastProposer)
}
