package temp

import (
	"errors"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"github.com/clearmatics/autonity/core/types"
)

// GetCommittee retrieves the committee for the given header.
func GetCommittee(header *types.Header, chain consensus.ChainReader) (types.Committee, error) {
	if header.IsGenesis() {
		return header.Committee, nil
	}
	parent := chain.GetHeaderByHash(header.ParentHash)
	if parent == nil {
		return nil, errors.New("unknown block")
	}
	return parent.Committee, nil
}

// GetCommitteeSet returns the committee to be used for validating the block
// associated with header. The parent paramer is optional, if it is not
// provided it will be looked up.
func GetCommitteeSet(header *types.Header, chain consensus.ChainReader) (committee.Set, error) {

	var previousProposer common.Address
	// The genesis block has no parent, so the committee is whatever is defined
	// in the block.
	if header.IsGenesis() {
		return committee.NewRoundRobinSet(header.Committee, previousProposer)
	}
	parent := chain.GetHeaderByHash(header.ParentHash)
	if parent == nil {
		return nil, errors.New("unknown block")
	}
	// The genesis block has no ProposerSeal so there is no address to recover
	// in this case.
	if !parent.IsGenesis() {
		var err error
		previousProposer, err = types.Ecrecover(parent)
		if err != nil {
			return nil, err
		}
	}
	return committee.NewRoundRobinSet(parent.Committee, previousProposer)
}
