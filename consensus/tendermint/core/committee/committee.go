package committee

import (
	"errors"
	"github.com/autonity/autonity/autonity"
	"math/big"
	"sort"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	ethcore "github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
)

type RoundRobinCommittee struct {
	members           types.Committee
	lastBlockProposer common.Address
	totalPower        *big.Int
	allProposers      map[int64]types.CommitteeMember // cached computed values
	roundRobinOffset  int64
	mu                sync.RWMutex // members doesn't need to be protected as it is read-only
}

func NewRoundRobinSet(members types.Committee, lastBlockProposer common.Address) (*RoundRobinCommittee, error) {
	// Ensure non empty set
	if len(members) == 0 {
		return nil, ErrEmptyCommitteeSet
	}

	//Create new roundRobinSet
	committee := &RoundRobinCommittee{
		members:           members,
		lastBlockProposer: lastBlockProposer,
		totalPower:        new(big.Int),
		allProposers:      make(map[int64]types.CommitteeMember),
	}

	// sort validator
	sort.Sort(committee.members)

	// calculate total power
	for _, m := range committee.members {
		committee.totalPower.Add(committee.totalPower, m.VotingPower)
	}

	// calculate offset for round robin selection of next proposer
	committee.roundRobinOffset = getMemberIndex(committee.members, lastBlockProposer)
	if len(members) > 1 {
		committee.roundRobinOffset++
	}
	committee.allProposers[0] = committee.getNextProposer(0)

	return committee, nil
}

func (set *RoundRobinCommittee) Committee() types.Committee {
	set.mu.RLock()
	defer set.mu.RUnlock()
	return copyMembers(set.members)
}

func (set *RoundRobinCommittee) GetByIndex(i int) (types.CommitteeMember, error) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	if i < 0 || i >= len(set.members) {
		return types.CommitteeMember{}, consensus.ErrCommitteeMemberNotFound
	}
	return set.members[i], nil
}

func (set *RoundRobinCommittee) GetByAddress(addr common.Address) (int, types.CommitteeMember, error) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	for i, member := range set.members {
		if addr == member.Address {
			return i, member, nil
		}
	}
	return -1, types.CommitteeMember{}, consensus.ErrCommitteeMemberNotFound
}

func (set *RoundRobinCommittee) GetProposer(round int64) types.CommitteeMember {
	set.mu.Lock()
	defer set.mu.Unlock()

	v, ok := set.allProposers[round]
	if !ok {
		v = set.getNextProposer(round)
		set.allProposers[round] = v
	}

	return v
}

func (set *RoundRobinCommittee) SetLastHeader(block *types.Block) {
	return
}

func (set *RoundRobinCommittee) Quorum() *big.Int {
	return bft.Quorum(set.totalPower)
}

func (set *RoundRobinCommittee) F() *big.Int {
	return bft.F(set.totalPower)
}

func (set *RoundRobinCommittee) getNextProposer(round int64) types.CommitteeMember {
	return set.members[nextProposerIndex(set.roundRobinOffset, round, int64(len(set.members)))]
}

func nextProposerIndex(offset, round, committeeSize int64) int64 {
	// Round-Robin
	return (offset + round) % committeeSize
}

func getMemberIndex(members types.Committee, memberAddr common.Address) int64 {
	var index = -1
	for i, member := range members {
		if memberAddr == member.Address {
			index = i
		}
	}
	return int64(index)
}

type WeightedRandomSamplingCommittee struct {
	previousHeader   *types.Header
	bc               *ethcore.BlockChain
	autonityContract *autonity.Contract // proposer buffering are stored at contract wrapper side thus more components can share them.
}

func NewWeightedRandomSamplingCommittee(previousBlock *types.Block, autonityContract *autonity.Contract, bc *ethcore.BlockChain) *WeightedRandomSamplingCommittee {
	return &WeightedRandomSamplingCommittee{
		previousHeader:   previousBlock.Header(),
		bc:               bc,
		autonityContract: autonityContract,
	}
}

// Return the underlying types.Committee
func (w *WeightedRandomSamplingCommittee) Committee() types.Committee {
	return w.previousHeader.Committee
}

func (w *WeightedRandomSamplingCommittee) SetLastHeader(block *types.Block) {
	w.previousHeader = block.Header()
}

// Get validator by index
func (w *WeightedRandomSamplingCommittee) GetByIndex(i int) (types.CommitteeMember, error) {
	if i < 0 || i >= len(w.previousHeader.Committee) {
		return types.CommitteeMember{}, consensus.ErrCommitteeMemberNotFound
	}
	return w.previousHeader.Committee[i], nil
}

// Get validator by given address
func (w *WeightedRandomSamplingCommittee) GetByAddress(addr common.Address) (int, types.CommitteeMember, error) {
	// TODO Promote types.Committee to a struct containing a slice, this will
	// allow for caching of other information like total power ... etc.
	m := w.previousHeader.CommitteeMember(addr)
	if m == nil {
		return -1, types.CommitteeMember{}, consensus.ErrCommitteeMemberNotFound
	}

	return -1, *m, nil
}

// Get the round proposer
func (w *WeightedRandomSamplingCommittee) GetProposer(round int64) types.CommitteeMember {
	// If previous header was the genesis block then we will not yet have
	// deployed the autonity contract so will take the proposer as the first
	// defined validator of the genesis block.
	if w.previousHeader.IsGenesis() {
		sort.Sort(w.previousHeader.Committee)
		return w.previousHeader.Committee[round%int64(len(w.previousHeader.Committee))]
	}
	proposer := w.autonityContract.GetProposer(w.previousHeader, w.previousHeader.Number.Uint64(), round)
	member := w.previousHeader.CommitteeMember(proposer)
	if member == nil {
		//Should not happen in live network, edge case
		return types.CommitteeMember{}
	}
	return *member
}

// Get the optimal quorum size
func (w *WeightedRandomSamplingCommittee) Quorum() *big.Int {
	return bft.Quorum(w.previousHeader.TotalVotingPower())
}

func (w *WeightedRandomSamplingCommittee) F() *big.Int {
	return bft.F(w.previousHeader.TotalVotingPower())
}

var ErrEmptyCommitteeSet = errors.New("committee set can't be empty")

func copyMembers(members types.Committee) types.Committee {
	membersCopy := make(types.Committee, len(members))
	for i, val := range members {
		membersCopy[i] = val
	}
	return membersCopy
}
