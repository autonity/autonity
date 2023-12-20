package interfaces

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"math/big"
)

type Committee interface {
	// CommitteeMember query a committee member by address.
	CommitteeMember(address common.Address) *types.CommitteeMember

	// SetCommittee update the committee on epoch rotation
	SetCommittee(c *types.Committee)

	// Committee Return the underlying types.Committee
	Committee() *types.Committee

	// GetByIndex Get validator by index
	GetByIndex(i int) (*types.CommitteeMember, error)

	// GetByAddress Get validator by given address
	GetByAddress(addr common.Address) (*types.CommitteeMember, error)

	// GetProposer Get the round proposer
	GetProposer(round int64) *types.CommitteeMember

	// SetLastHeader Update with lastest block header
	SetLastHeader(header *types.Header)

	// Quorum Get the optimal quorum size
	Quorum() *big.Int

	// F Get the maximum number of faulty nodes
	F() *big.Int
}
