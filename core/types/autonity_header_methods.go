package types

import "github.com/clearmatics/autonity/common"

func (h *Header) IsGenesis() bool {
	return h.Number.Uint64() == 0
}

// CommitteeMember returns the committee member having the given address or
// nil if there is none.
func (h *Header) CommitteeMember(address common.Address) *CommitteeMember {
	if h.committeeMap == nil {
		h.committeeMap = make(map[common.Address]*CommitteeMember)
		for i := range h.Committee {
			member := h.Committee[i]
			h.committeeMap[member.Address] = &member
		}
	}
	return h.committeeMap[address]
}
