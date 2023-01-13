package interfaces

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"math/big"
)

type Committee interface {
	// Return the underlying types.Committee
	Committee() types.Committee
	// Get validator by index
	GetByIndex(i int) (types.CommitteeMember, error)
	// Get validator by given address
	GetByAddress(addr common.Address) (int, types.CommitteeMember, error)
	// Get the round proposer
	GetProposer(round int64) types.CommitteeMember
	// Update with lastest block
	SetLastBlock(block *types.Block)
	// Get the optimal quorum size
	Quorum() *big.Int
	// Get the maximum number of faulty nodes
	F() *big.Int
}
