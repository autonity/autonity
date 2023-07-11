package types

import (
	"math/big"

	"github.com/autonity/autonity/common"
)

func (h *Header) IsGenesis() bool {
	return h.Number.Uint64() == 0
}

// CommitteeMember returns the committee member having the given address or
// nil if there is none.
func (h *Header) CommitteeMember(address common.Address) *CommitteeMember {
	h.once.Do(func() {
		h.committeeMap = make(map[common.Address]*CommitteeMember)
		for i := range h.Committee {
			member := h.Committee[i]
			h.committeeMap[member.Address] = &member
		}
	})
	return h.committeeMap[address]
}

// TotalVotingPower returns the total voting power contained in the committee
// for the block associated with this header.
func (h *Header) TotalVotingPower() *big.Int {
	return h.Committee.TotalVotingPower()
}

// TotalVotingPower returns the total voting power contained in the committee.
func (c Committee) TotalVotingPower() *big.Int {
	total := new(big.Int)
	for _, m := range c {
		total.Add(total, m.VotingPower)
	}
	return total
}
