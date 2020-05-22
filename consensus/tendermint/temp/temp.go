package temp

import (
	"errors"

	"github.com/clearmatics/autonity/consensus"
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
