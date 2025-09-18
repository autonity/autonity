package message

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/log"
)

type Set struct {
	// In some conditions we might receive prevotes or precommit before
	// receiving a proposal, so we must save received message with different proposed block hash.
	votes map[common.Hash][]Vote // map[proposedBlockHash][]vote

	/* we use AggregatedPower because we cannot simply sum the voting power of the votes. This is because we might have:
	* 1. duplicated votes between different overlapping aggregates for the same value
	* 2. equivocated votes from the same validator across different values
	 */
	powers     map[common.Hash]*AggregatedPower // cumulative voting power for each value
	totalPower *AggregatedPower                 // total voting power of votes

	sync.RWMutex
}

func NewSet() *Set {
	return &Set{
		votes:      make(map[common.Hash][]Vote),
		powers:     make(map[common.Hash]*AggregatedPower),
		totalPower: NewAggregatedPower(),
	}
}

// returns whether the new signers increased the power or the vote was redundant
func (s *Set) Add(vote Vote) bool {
	s.Lock()
	defer s.Unlock()

	// will be set to true if the vote brings an increase in voting power in Core.
	// votes that contribute by equivocation (increasing power of a specific value, but
	// not the total power for the votes set) are considered not redundant (as they should be gossiped)
	voteContributed := false

	value := vote.Value()
	previousVotes, ok := s.votes[value]
	if !ok {
		s.votes[value] = make([]Vote, 1)
		s.powers[value] = NewAggregatedPower()
	}

	// update total power and power for value
	it := vote.Signers().NewIterator()
	for it.Next() {
		signerContributed := s.totalPower.Set(it.Index(), vote.Signers().PowerByIndex(it.Index()))
		signerContributedToValue := s.powers[value].Set(it.Index(), vote.Signers().PowerByIndex(it.Index()))
		voteContributed = voteContributed || signerContributed || signerContributedToValue
	}

	// check if we are adding the first vote
	if len(previousVotes) == 0 {
		s.votes[value][0] = vote
		return voteContributed
	}

	// not adding the first vote, aggregate it with previous ones
	// if the vote is redundant, there is no need to add it to the msg store,
	// as it will be discarded during aggregation anyway
	if voteContributed {
		s.add(previousVotes, vote)
	}

	return voteContributed
}

// returns -1 if not aggregatable, index of aggregatable vote otherwise
func isAggregatable(previousVotes []Vote, newVote Vote) int {
	index := -1
	for i, previousVote := range previousVotes {
		if previousVote.Signers().RespectsBoundaries(newVote.Signers()) {
			index = i
			break
		}
	}
	return index
}

func logErrorPrevotes(previousVotes []Vote, newVote Vote, index int, aggregatedVotes []*Prevote) {
	log.Error("aggregated votes length not equal to 1",
		"previousVotes[index]", previousVotes[index].String(),
		"newVote", newVote.String(), "index", index)

	for i, aggregatedVote := range aggregatedVotes {
		log.Error("aggregated votes", "i", i, "aggregatedVote", aggregatedVote.String())
	}
}

func logErrorPrecommits(previousVotes []Vote, newVote Vote, index int, aggregatedVotes []*Precommit) {
	log.Error("aggregated votes length not equal to 1",
		"previousVotes[index]", previousVotes[index].String(),
		"newVote", newVote.String(), "index", index)

	for i, aggregatedVote := range aggregatedVotes {
		log.Error("aggregated votes", "i", i, "aggregatedVote", aggregatedVote.String())
	}
}

// assumes that `previousVotes` are all non-mergeable with each other
func (s *Set) add(previousVotes []Vote, newVote Vote) {
	value := newVote.Value()

	// check if it can be merged with any previous vote
	index := isAggregatable(previousVotes, newVote)
	if index == -1 {
		// non aggregatable, append at the end and return
		previousVotes = append(previousVotes, newVote)
		s.votes[value] = previousVotes
		return
	}

	switch newVote.Code() {
	case PrevoteCode:
		aggregatedVotes := AggregatePrevotes([]Vote{previousVotes[index], newVote})
		if len(aggregatedVotes) != 1 {
			logErrorPrevotes(previousVotes, newVote, index, aggregatedVotes)
			panic("aggregated votes should have length 1")
		}
		previousVotes[index] = aggregatedVotes[0]
	case PrecommitCode:
		aggregatedVotes := AggregatePrecommits([]Vote{previousVotes[index], newVote})
		if len(aggregatedVotes) != 1 {
			logErrorPrecommits(previousVotes, newVote, index, aggregatedVotes)
			panic("aggregated votes should have length 1")
		}
		previousVotes[index] = aggregatedVotes[0]
	default:
		panic(fmt.Sprintf("Trying to add a vote that is not Prevote nor Precommit: %d", newVote.Code()))
	}
	s.votes[value] = previousVotes
}

func (s *Set) Messages() []Msg {
	s.RLock()
	defer s.RUnlock()

	messages := make([]Msg, 0)
	for _, votes := range s.votes {
		for _, vote := range votes {
			messages = append(messages, vote.(Msg))
		}
	}
	return messages
}

func (s *Set) PowerFor(h common.Hash) *AggregatedPower {
	s.RLock()
	defer s.RUnlock()

	_, ok := s.powers[h]
	if ok {
		return s.powers[h].Copy() // return copy to avoid data race
	}
	return NewAggregatedPower()
}

func (s *Set) TotalPower() *AggregatedPower {
	s.RLock()
	defer s.RUnlock()

	// NOTE: in case of equivocated messages, we count power only once
	return s.totalPower.Copy() // return copy to avoid data race
}

func (s *Set) VotesFor(blockHash common.Hash) []Vote {
	s.RLock()
	defer s.RUnlock()

	return s.votes[blockHash]
}

func (s *Set) DumpMsgView() ([]common.Hash, []*big.Int) {
	s.RLock()         //nolint
	defer s.RUnlock() //nolint

	var values []common.Hash //nolint
	var signers []*big.Int   //nolint
	for v, votes := range s.powers {
		values = append(values, v)
		signers = append(signers, new(big.Int).Set(votes.Signers()))
	}

	return values, signers
}
