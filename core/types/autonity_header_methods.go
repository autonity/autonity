package types

import "github.com/clearmatics/autonity/common"

func (h *Header) IsGenesis() bool {
	return h.Number.Uint64() == 0
}

// CommitteeMember returns the committee member having the given address or
// nil if there is none.
func (h *Header) CommitteeMember(address common.Address) *CommitteeMember {
	return h.CommitteMemberMap()[address]
}

func (h *Header) CommitteMemberMap() map[common.Address]*CommitteeMember {
	if h.committeeMap == nil {
		h.committeeMap = make(map[common.Address]*CommitteeMember)
		for i := range h.Committee {
			member := h.Committee[i]
			h.committeeMap[member.Address] = &member
		}
	}
	return h.committeeMap
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
