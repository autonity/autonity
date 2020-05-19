package temp

import (
	"errors"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"github.com/clearmatics/autonity/core/types"
)

// retrieve list of getCommittee for the block header passed as parameter
func GetCommittee(header *types.Header, parents []*types.Header, chain consensus.ChainReader) (committee.Set, error) {

	// We can't use savedCommittee if parents are being passed :
	// those blocks are not yet saved in the blockchain.
	// autonity will stop processing the received blockchain from the moment an error appears.
	// See insertChain in blockchain.go
	if len(parents) > 0 {
		parent := parents[len(parents)-1]
		lastMiner, err := types.Ecrecover(parent)
		if err != nil {
			return nil, err
		}
		return committee.NewRoundRobinSet(parent.Committee, lastMiner)
	}

	number := header.Number.Uint64()
	if number == 0 {
		number = 1
	}
	// Check for existence of parent
	parentHeader := chain.GetHeaderByNumber(number - 1)
	if parentHeader == nil {
		return nil, errors.New("unknown block")
	}

	var lastProposer common.Address
	if number > 1 {
		var err error
		lastProposer, err = types.Ecrecover(parentHeader)
		if err != nil {
			return nil, err
		}
	}
	return committee.NewRoundRobinSet(parentHeader.Committee, lastProposer)
}
