package message

import (
	"math/big"
	"sync"

	"github.com/autonity/autonity/common"
)

type isVote interface {
	*Prevote | *Precommit
	Msg
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
}

func (s *Set[T]) Add(vote T) {
	s.Lock()
	defer s.Unlock()
	sender := vote.Sender()
	value := vote.Value()
	// Check first if we already received a message from this pal.
	if _, ok := s.messages[sender]; ok {
		// TODO : double signing fault ! Accountability
		return
	}

	if _, ok := s.votes[value]; !ok {
		s.votes[value] = make(map[common.Address]T)
	}
	s.votes[value][sender] = vote
	s.messages[sender] = vote
}

func (s *Set[T]) Messages() []Msg {
	s.RLock()
	defer s.RUnlock()
	result := make([]Msg, len(s.messages))
	k := 0
	for _, v := range s.messages {
		result[k] = v
		k++
	}
	return result
}

func (s *Set[T]) PowerFor(h common.Hash) *big.Int {
	s.RLock()
	defer s.RUnlock()
	if votes, ok := s.votes[h]; ok {
		power := new(big.Int)
		for _, v := range votes {
			power.Add(power, v.Power())
		}
		return power
	}
	return new(big.Int)
}

func (s *Set[T]) TotalPower() *big.Int {
	s.RLock()
	defer s.RUnlock()
	power := new(big.Int)
	for _, msg := range s.messages {
		power.Add(power, msg.Power())
	}
	return power
}

func (s *Set[T]) VotesFor(blockHash common.Hash) []T {
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.votes[blockHash]; !ok {
		return nil
	}
	messages := make([]T, 0, len(s.votes[blockHash]))
	for _, v := range s.votes[blockHash] {
		messages = append(messages, v)
	}
	return messages
}
