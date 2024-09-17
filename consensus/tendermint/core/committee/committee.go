package committee

import (
	"errors"
	"math/big"
	"sync"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
)

var ErrEmptyCommitteeSet = errors.New("committee set can't be empty")

type RoundRobinCommittee struct {
	committee         *types.Committee
	lastBlockProposer common.Address
	totalPower        *big.Int
	allProposers      map[int64]*types.CommitteeMember // cached computed values
	roundRobinOffset  int64
	mu                sync.RWMutex // members doesn't need to be protected as it is read-only
}

func NewRoundRobinSet(committee *types.Committee, lastBlockProposer common.Address) (*RoundRobinCommittee, error) {
	// Ensure non empty set
	if committee == nil || len(committee.Members) == 0 {
		return nil, ErrEmptyCommitteeSet
	}

	types.SortCommitteeMembers(committee.Members)
	//Create new roundRobinSet
	set := &RoundRobinCommittee{
		committee:         committee,
		lastBlockProposer: lastBlockProposer,
		totalPower:        committee.TotalVotingPower(),
		allProposers:      make(map[int64]*types.CommitteeMember),
	}

	// calculate offset for round robin selection of next proposer
	set.roundRobinOffset = getMemberIndex(set.committee, lastBlockProposer)
	if committee.Len() > 1 {
		set.roundRobinOffset++
	}
	set.allProposers[0] = set.getNextProposer(0)

	return set, nil
}

func (set *RoundRobinCommittee) Committee() *types.Committee {
	set.mu.RLock()
	defer set.mu.RUnlock()
	return set.committee
}

func (set *RoundRobinCommittee) SetCommittee(committee *types.Committee) {
	set.mu.Lock()
	defer set.mu.Unlock()
	set.committee = committee
}

func (set *RoundRobinCommittee) MemberByIndex(i int) (*types.CommitteeMember, error) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	m := set.committee.MemberByIndex(i)
	if m == nil {
		return nil, consensus.ErrCommitteeMemberNotFound
	}
	return m, nil
}

func (set *RoundRobinCommittee) MemberByAddress(addr common.Address) (*types.CommitteeMember, error) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	m := set.committee.MemberByAddress(addr)
	if m == nil {
		return nil, consensus.ErrCommitteeMemberNotFound
	}

	return m, nil
}

func (set *RoundRobinCommittee) GetProposer(round int64) *types.CommitteeMember {
	set.mu.Lock()
	defer set.mu.Unlock()

	v, ok := set.allProposers[round]
	if !ok {
		v = set.getNextProposer(round)
		set.allProposers[round] = v
	}

	return v
}

func (set *RoundRobinCommittee) SetLastHeader(_ *types.Header) {
	return
}

func (set *RoundRobinCommittee) Quorum() *big.Int {
	return bft.Quorum(set.totalPower)
}

func (set *RoundRobinCommittee) F() *big.Int {
	return bft.F(set.totalPower)
}

func (set *RoundRobinCommittee) getNextProposer(round int64) *types.CommitteeMember {
	return &set.committee.Members[nextProposerIndex(set.roundRobinOffset, round, int64(set.committee.Len()))]
}

func nextProposerIndex(offset, round, committeeSize int64) int64 {
	// Round-Robin
	return (offset + round) % committeeSize
}

func getMemberIndex(committee *types.Committee, memberAddr common.Address) int64 {
	var index = -1
	for i, member := range committee.Members {
		if memberAddr == member.Address {
			index = i
		}
	}
	return int64(index)
}

type WeightedRandomSamplingCommittee struct {
	committee        *types.Committee
	previousHeader   *types.Header
	autonityContract *autonity.ProtocolContracts
}

func NewWeightedRandomSamplingCommittee(previousHeader *types.Header, committee *types.Committee, autonityContract *autonity.ProtocolContracts) *WeightedRandomSamplingCommittee {
	return &WeightedRandomSamplingCommittee{
		committee:        committee,
		previousHeader:   previousHeader,
		autonityContract: autonityContract,
	}
}

func (w *WeightedRandomSamplingCommittee) SetCommittee(committee *types.Committee) {
	w.committee = committee
}

// Return the underlying types.Committee
func (w *WeightedRandomSamplingCommittee) Committee() *types.Committee {
	return w.committee
}

func (w *WeightedRandomSamplingCommittee) SetLastHeader(header *types.Header) {
	w.previousHeader = header
}

// Get validator by index
func (w *WeightedRandomSamplingCommittee) MemberByIndex(i int) (*types.CommitteeMember, error) {
	m := w.committee.MemberByIndex(i)
	if m == nil {
		return nil, consensus.ErrCommitteeMemberNotFound
	}
	return m, nil
}

// MemberByAddress Get validator by given address
func (w *WeightedRandomSamplingCommittee) MemberByAddress(addr common.Address) (*types.CommitteeMember, error) {
	m := w.committee.MemberByAddress(addr)
	if m == nil {
		return nil, consensus.ErrCommitteeMemberNotFound
	}

	return m, nil
}

// Get the round proposer
func (w *WeightedRandomSamplingCommittee) GetProposer(round int64) *types.CommitteeMember {
	proposer := w.autonityContract.Proposer(w.committee, nil, w.previousHeader.Number.Uint64(), round)
	member := w.committee.MemberByAddress(proposer)
	if member == nil {
		log.Crit("Cannot find elected proposer from current committee")
	}
	return member
}

// Get the optimal quorum size
func (w *WeightedRandomSamplingCommittee) Quorum() *big.Int {
	return bft.Quorum(w.committee.TotalVotingPower())
}

func (w *WeightedRandomSamplingCommittee) F() *big.Int {
	return bft.F(w.committee.TotalVotingPower())
}
