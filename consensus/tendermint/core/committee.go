package core

import (
	"errors"
	"sort"
	"sync"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	ethcore "github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
)

type committee interface {
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
	Quorum() uint64
	// Get the maximum number of faulty nodes
	F() uint64
}

type roundRobinCommittee struct {
	members           types.Committee
	lastBlockProposer common.Address
	totalPower        uint64
	allProposers      map[int64]types.CommitteeMember // cached computed values
	roundRobinOffset  int64
	mu                sync.RWMutex // members doesn't need to be protected as it is read-only
}

func newRoundRobinSet(members types.Committee, lastBlockProposer common.Address) (*roundRobinCommittee, error) {
	// Ensure non empty set
	if len(members) == 0 {
		return nil, ErrEmptyCommitteeSet
	}

	//Create new roundRobinSet
	committee := &roundRobinCommittee{
		members:           members,
		lastBlockProposer: lastBlockProposer,
		allProposers:      make(map[int64]types.CommitteeMember),
	}

	// sort validator
	sort.Sort(committee.members)

	// calculate total power
	for _, m := range committee.members {
		committee.totalPower += m.VotingPower.Uint64()
	}

	// calculate offset for round robin selection of next proposer
	committee.roundRobinOffset = getMemberIndex(committee.members, lastBlockProposer)
	if len(members) > 1 {
		committee.roundRobinOffset++
	}
	committee.allProposers[0] = committee.getNextProposer(0)

	return committee, nil
}

func (set *roundRobinCommittee) Committee() types.Committee {
	set.mu.RLock()
	defer set.mu.RUnlock()
	return copyMembers(set.members)
}

func (set *roundRobinCommittee) GetByIndex(i int) (types.CommitteeMember, error) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	if i < 0 || i >= len(set.members) {
		return types.CommitteeMember{}, consensus.ErrCommitteeMemberNotFound
	}
	return set.members[i], nil
}

func (set *roundRobinCommittee) GetByAddress(addr common.Address) (int, types.CommitteeMember, error) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	for i, member := range set.members {
		if addr == member.Address {
			return i, member, nil
		}
	}
	return -1, types.CommitteeMember{}, consensus.ErrCommitteeMemberNotFound
}

func (set *roundRobinCommittee) GetProposer(round int64) types.CommitteeMember {
	set.mu.Lock()
	defer set.mu.Unlock()

	v, ok := set.allProposers[round]
	if !ok {
		v = set.getNextProposer(round)
		set.allProposers[round] = v
	}

	return v
}

func (set *roundRobinCommittee) SetLastBlock(block *types.Block) {
	return
}

func (set *roundRobinCommittee) Quorum() uint64 {
	return bft.Quorum(set.totalPower)
}

func (set *roundRobinCommittee) F() uint64 {
	return bft.F(set.totalPower)
}

func (set *roundRobinCommittee) getNextProposer(round int64) types.CommitteeMember {
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

type weightedRandomSamplingCommittee struct {
	previousHeader         *types.Header
	bc                     *ethcore.BlockChain // Todo : remove this dependency
	autonityContract       *autonity.Contract
	previousBlockStateRoot common.Hash
	cachedProposer         map[int64]types.CommitteeMember
}

func newWeightedRandomSamplingCommittee(previousBlock *types.Block, autonityContract *autonity.Contract, bc *ethcore.BlockChain) *weightedRandomSamplingCommittee {
	return &weightedRandomSamplingCommittee{
		previousHeader:         previousBlock.Header(),
		bc:                     bc,
		autonityContract:       autonityContract,
		previousBlockStateRoot: previousBlock.Root(),
		cachedProposer:         make(map[int64]types.CommitteeMember),
	}
}

// Return the underlying types.Committee
func (w *weightedRandomSamplingCommittee) Committee() types.Committee {
	return w.previousHeader.Committee
}

func (w *weightedRandomSamplingCommittee) SetLastBlock(block *types.Block) {
	w.previousHeader = block.Header()
	w.previousBlockStateRoot = block.Root()
	w.cachedProposer = make(map[int64]types.CommitteeMember)
}

// Get validator by index
func (w *weightedRandomSamplingCommittee) GetByIndex(i int) (types.CommitteeMember, error) {
	if i < 0 || i >= len(w.previousHeader.Committee) {
		return types.CommitteeMember{}, consensus.ErrCommitteeMemberNotFound
	}
	return w.previousHeader.Committee[i], nil
}

// Get validator by given address
func (w *weightedRandomSamplingCommittee) GetByAddress(addr common.Address) (int, types.CommitteeMember, error) {
	// TODO Promote types.Committee to a struct containing a slice, this will
	// allow for caching of other information like total power ... etc.
	m := w.previousHeader.CommitteeMember(addr)
	if m == nil {
		return -1, types.CommitteeMember{}, consensus.ErrCommitteeMemberNotFound
	}

	return -1, *m, nil
}

// Get the round proposer
func (w *weightedRandomSamplingCommittee) GetProposer(round int64) types.CommitteeMember {
	if res, ok := w.cachedProposer[round]; ok {
		return res
	}
	// If previous header was the genesis block then we will not yet have
	// deployed the autonity contract so will take the proposer as the first
	// defined validator of the genesis block.
	if w.previousHeader.IsGenesis() {
		sort.Sort(w.previousHeader.Committee)
		return w.previousHeader.Committee[round%int64(len(w.previousHeader.Committee))]
	}
	// state.New has started taking a snapshot.Tree but it seems to be only for
	// performance, see - https://github.com/autonity/autonity/pull/20152
	statedb, err := state.New(w.previousBlockStateRoot, w.bc.StateCache(), nil)
	if err != nil {
		log.Error("cannot load state from block chain.")
		return types.CommitteeMember{}
	}
	proposer := w.autonityContract.GetProposerFromAC(w.previousHeader, statedb, w.previousHeader.Number.Uint64(), round)
	member := w.previousHeader.CommitteeMember(proposer)

	w.cachedProposer[round] = *member
	return *member
	// TODO make this return an error
}

// Get the optimal quorum size
func (w *weightedRandomSamplingCommittee) Quorum() uint64 {
	return bft.Quorum(w.previousHeader.TotalVotingPower())
}

func (w *weightedRandomSamplingCommittee) F() uint64 {
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
