package message

import (
	"math/big"
	"sync"

	"github.com/autonity/autonity/common"
)

//TODO(lorenzo) analyze more duplicated msgs and equivocation scnearios

/*
type isVote interface {
	*Prevote | *Precommit
	IndividualMsg
}

func NewSet[T isVote]() *Set[T] {
	return &Set[T]{
		votes:    make(map[common.Hash]map[common.Address]T),
		messages: make(map[common.Address]T),
	}
}

type Set[T isVote] struct {
	// In some conditions we might receive prevotes or precommit before
	// receiving a proposal, so we must save received message with different proposed block hash.
	votes    map[common.Hash]map[common.Address]T // map[proposedBlockHash]map[validatorAddress]vote
	messages map[common.Address]T
	sync.RWMutex
}*/

type isVote interface {
	*AggregatePrevote | *AggregatePrecommit
	AggregateMsg
}

func NewSet[T isVote]() *Set[T] {
	return &Set[T]{
		votes: make(map[common.Hash][]T),
	}
}

type Set[T isVote] struct {
	// In some conditions we might receive prevotes or precommit before
	// receiving a proposal, so we must save received message with different proposed block hash.
	votes        map[common.Hash][]T // map[proposedBlockHash][]vote
	sync.RWMutex                     //TODO(lorenzo) why this mutex (for state dumper?)
}

func (s *Set[T]) Add(vote T) {
	s.Lock()
	defer s.Unlock()

	//TODO(lorenzo) now we can have equivocated messages in core, how does this impact core?

	value := vote.Value()

	if _, ok := s.votes[value]; !ok {
		s.votes[value] = make([]T, 0) //TODO(lorenzo) allocate some more capacity?
	}
	s.votes[value] = append(s.votes[value], vote)
}

func (s *Set[T]) Messages() []Msg {
	s.RLock()
	defer s.RUnlock()
	result := make([]Msg, 0)
	for _, v := range s.votes {
		for _, vote := range v {
			result = append(result, vote)
		}
	}
	return result
}

func (s *Set[T]) PowerFor(h common.Hash) *big.Int {
	s.RLock()
	defer s.RUnlock()

	//TODO(lorenzo) can I just call the s.power?
	if votes, ok := s.votes[h]; ok {
		accountedFor := make(map[common.Address]struct{})
		return s.power(votes, accountedFor)
	}
	return new(big.Int)
}

func (s *Set[T]) TotalPower() *big.Int {
	s.RLock()
	defer s.RUnlock()
	power := new(big.Int)
	// NOTE: in case of equivocated messages, we count power only once --> write a test for it
	accountedFor := make(map[common.Address]struct{})
	for _, votes := range s.votes {
		power.Add(power, s.power(votes, accountedFor))
	}

	return power
}

// TODO(lorenzo) This logic is a bit twisted but it works.
// the key here is that we have to count power only once per sender
// across multiple aggregate votes.
// TODO(lorenzo) write tests for it, and make sure accountedFor works as intended (for totalpower)
func (s *Set[T]) power(votes []T, accountedFor map[common.Address]struct{}) *big.Int {
	power := new(big.Int)

	for _, v := range votes {
		addresses := v.Senders().Addresses()
		powers := v.Senders().Powers()
		for _, index := range v.Senders().FlattenUniq() {
			_, accounted := accountedFor[addresses[index]]
			if accounted {
				continue
			}
			power.Add(power, powers[index])
			accountedFor[addresses[index]] = struct{}{}
		}
	}
	return power
}

func (s *Set[T]) VotesFor(blockHash common.Hash) []T {
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.votes[blockHash]; !ok {
		return nil
	}
	//TODO(lorenzo) might not need copy here, double check
	messages := make([]T, 0, len(s.votes[blockHash]))
	for _, v := range s.votes[blockHash] {
		messages = append(messages, v)
	}
	return messages
}
