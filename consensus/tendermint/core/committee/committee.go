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
	allProposers      map[int64]*types.CommitteeMember // cached computed values
	roundRobinOffset  int64
	mu                sync.RWMutex // members doesn't need to be protected as it is read-only
}

func NewRoundRobinSet(committee *types.Committee, lastBlockProposer common.Address) (*RoundRobinCommittee, error) {
	// Ensure non empty set
	if committee == nil || len(committee.Members) == 0 {
		return nil, ErrEmptyCommitteeSet
	}

	// Sort committee member
	committee.Sort()

	//Create new roundRobinSet
	set := &RoundRobinCommittee{
		committee:         committee,
		lastBlockProposer: lastBlockProposer,
		allProposers:      make(map[int64]*types.CommitteeMember),
	}

	// calculate offset for round-robin selection of next proposer
	var lastProposerIndex int
	for i, m := range committee.Members {
		if m.Address == lastBlockProposer {
			lastProposerIndex = i
			break
		}
	}

	set.roundRobinOffset = int64(lastProposerIndex)
	if len(committee.Members) > 1 {
		set.roundRobinOffset++
	}
	set.allProposers[0] = set.getNextProposer(0)

	return set, nil
}

func (set *RoundRobinCommittee) CommitteeMember(address common.Address) *types.CommitteeMember {
	return set.committee.CommitteeMember(address)
}

func (set *RoundRobinCommittee) SetCommittee(committee *types.Committee) {
	set.mu.Lock()
	defer set.mu.Unlock()
	set.committee = committee
}

func (set *RoundRobinCommittee) Committee() *types.Committee {
	set.mu.RLock()
	defer set.mu.RUnlock()
	return set.committee
}

func (set *RoundRobinCommittee) GetByIndex(i int) (*types.CommitteeMember, error) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	if i < 0 || i >= len(set.committee.Members) {
		return nil, consensus.ErrCommitteeMemberNotFound
	}
	return set.committee.Members[i], nil
}

func (set *RoundRobinCommittee) GetByAddress(addr common.Address) (*types.CommitteeMember, error) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	member := set.committee.CommitteeMember(addr)
	if member != nil {
		return member, nil
	}
	return nil, consensus.ErrCommitteeMemberNotFound
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
	return bft.Quorum(set.committee.TotalVotingPower())
}

func (set *RoundRobinCommittee) F() *big.Int {
	return bft.F(set.committee.TotalVotingPower())
}

func (set *RoundRobinCommittee) getNextProposer(round int64) *types.CommitteeMember {
	return set.committee.Members[nextProposerIndex(set.roundRobinOffset, round, int64(len(set.committee.Members)))]
}

func nextProposerIndex(offset, round, committeeSize int64) int64 {
	// Round-Robin
	return (offset + round) % committeeSize
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

func (w *WeightedRandomSamplingCommittee) CommitteeMember(address common.Address) *types.CommitteeMember {
	return w.committee.CommitteeMember(address)
}

func (w *WeightedRandomSamplingCommittee) SetCommittee(committee *types.Committee) {
	w.committee = committee
}

func (w *WeightedRandomSamplingCommittee) Committee() *types.Committee {
	return w.committee
}

func (w *WeightedRandomSamplingCommittee) SetLastHeader(header *types.Header) {
	w.previousHeader = header
}

// GetByIndex Get validator by index
func (w *WeightedRandomSamplingCommittee) GetByIndex(i int) (*types.CommitteeMember, error) {
	if i < 0 || i >= len(w.committee.Members) {
		return nil, consensus.ErrCommitteeMemberNotFound
	}
	return w.committee.Members[i], nil
}

// GetByAddress Get validator by given address
func (w *WeightedRandomSamplingCommittee) GetByAddress(addr common.Address) (*types.CommitteeMember, error) {
	m := w.committee.CommitteeMember(addr)
	if m == nil {
		return nil, consensus.ErrCommitteeMemberNotFound
	}

	return m, nil
}

// GetProposer Get the round proposer
func (w *WeightedRandomSamplingCommittee) GetProposer(round int64) *types.CommitteeMember {
	proposer := w.autonityContract.Proposer(w.committee, w.previousHeader.Number.Uint64(), round)
	member := w.committee.CommitteeMember(proposer)
	if member == nil {
		log.Crit("Cannot find elected proposer from current committee")
	}
	return member
}

func (w *WeightedRandomSamplingCommittee) Quorum() *big.Int {
	return bft.Quorum(w.committee.TotalVotingPower())
}

func (w *WeightedRandomSamplingCommittee) F() *big.Int {
	return bft.F(w.committee.TotalVotingPower())
}
