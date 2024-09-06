package interfaces

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"math/big"
)

type Committee interface {
	// Committee Return the underlying types.Committee
	Committee() *types.Committee

	// MemberByIndex Get validator by index
	MemberByIndex(i int) (*types.CommitteeMember, error)

	// MemberByAddress Get validator by given address
	MemberByAddress(addr common.Address) (*types.CommitteeMember, error)

	// GetProposer Get the round proposer
	GetProposer(round int64) *types.CommitteeMember

	// SetLastHeader Update with lastest block header
	SetLastHeader(block *types.Header)

	SetCommittee(committee *types.Committee)

	// Quorum Get the optimal quorum size
	Quorum() *big.Int

	// F Get the maximum number of faulty nodes
	F() *big.Int
}
