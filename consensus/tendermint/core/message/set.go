package message

import (
	"sync"

	"github.com/autonity/autonity/common"
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

func (s *Set) Add(vote Vote) {
	s.Lock()
	defer s.Unlock()

	value := vote.Value()
	previousVotes, ok := s.votes[value]
	if !ok {
		s.votes[value] = make([]Vote, 1)
		s.powers[value] = NewAggregatedPower()
	}

	// update total power and power for value
	for index, power := range vote.Signers().Powers() {
		s.totalPower.Set(index, power)
		s.powers[value].Set(index, power)
	}

	// check if we are adding the first vote
	if len(previousVotes) == 0 {
		s.votes[value][0] = vote
		return
	}

	// if not first vote, aggregate previous votes and new vote
	switch vote.(type) {
	case *Prevote:
		aggregatedVotes := AggregatePrevotesSimple(append(previousVotes, vote))
		s.votes[value] = make([]Vote, len(aggregatedVotes))
		for i, aggregatedVote := range aggregatedVotes {
			s.votes[value][i] = aggregatedVote
		}
	case *Precommit:
		aggregatedVotes := AggregatePrecommitsSimple(append(previousVotes, vote))
		s.votes[value] = make([]Vote, len(aggregatedVotes))
		for i, aggregatedVote := range aggregatedVotes {
			s.votes[value][i] = aggregatedVote
		}
	default:
		panic("Trying to add a vote that is not Prevote nor Precommit")
	}
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
