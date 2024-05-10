package message

import (
	"math/big"
	"sync"

	"github.com/autonity/autonity/common"
)

//TODO(lorenzo) refinements2, analyze more duplicated msgs and equivocation scnearios

func NewSet() *Set {
	return &Set{
		votes: make(map[common.Hash][]Vote),
		//powers: make(map[common.Hash]*big.Int),
	}
}

type Set struct {
	// In some conditions we might receive prevotes or precommit before
	// receiving a proposal, so we must save received message with different proposed block hash.
	votes map[common.Hash][]Vote // map[proposedBlockHash][]vote
	//powers       map[common.Hash]*big.Int //map[proposedBlockHash][]vote
	sync.RWMutex //TODO(lorenzo) refinements, do we need this lock since there is already one is round_messages?
}

func (s *Set) Add(vote Vote) {
	s.Lock()
	defer s.Unlock()

	value := vote.Value()
	previousVotes, ok := s.votes[value]

	if !ok {
		s.votes[value] = []Vote{vote}
		//s.powers[value] = vote.Senders().Power()
		return
	}

	// aggregate previous votes and vote
	//TODO(lorenzo) performance, verify that this doesn't create too much memory
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

func (s *Set) PowerFor(h common.Hash) *big.Int {
	s.RLock()
	defer s.RUnlock()

	votes, ok := s.votes[h]
	if !ok {
		return new(big.Int)
	}

	var messages []Msg
	for _, vote := range votes {
		messages = append(messages, vote.(Msg))
	}

	return Power(messages)
}

func (s *Set) TotalPower() *big.Int {
	s.RLock()
	defer s.RUnlock()

	// NOTE: in case of equivocated messages, we count power only once
	// TODO(lorenzo) refinements, write a test for it

	var messages []Msg

	for _, votes := range s.votes {
		for _, vote := range votes {
			messages = append(messages, vote.(Msg))
		}
	}

	return Power(messages)
}

func (s *Set) VotesFor(blockHash common.Hash) []Vote {
	s.RLock()
	defer s.RUnlock()

	return s.votes[blockHash]
}
