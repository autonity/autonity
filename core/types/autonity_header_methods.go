package types

import (
	"math"

	"github.com/clearmatics/autonity/common"
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
func (h *Header) TotalVotingPower() uint64 {
	return h.Committee.TotalVotingPower()
}

// TotalVotingPower returns the total voting power contained in the committee.
func (c Committee) TotalVotingPower() uint64 {
	var total uint64
	for _, m := range c {
		total += m.VotingPower.Uint64()
	}
	return total
}

// Qurum returns the voting power that constitutes a quorum.
func (c Committee) Quorum() uint64 {
	return uint64(math.Ceil((2 * float64(c.TotalVotingPower())) / 3.))
}

// F returns the voting power that constitutes the maximum tolerable failures.
func (c Committee) F() uint64 {
	return uint64(math.Ceil(float64(c.TotalVotingPower())/3)) - 1
}
