package message

import (
	"github.com/autonity/autonity/common"
	"math/big"
	"sync"
)

func NewSet[T Message]() *Set[T] {
	return &Set[T]{
		Votes:    make(map[common.Hash]map[common.Address]T),
		messages: make(map[common.Address]T),
	}
}

type Set[T Message] struct {
	// In some conditions we might receive prevotes or precommit before
	// receiving a proposal, so we must save received message with differents proposed block hash.
	Votes    map[common.Hash]map[common.Address]T // map[proposedBlockHash]map[validatorAddress]vote
	messages map[common.Address]T
	lock     sync.RWMutex
}

func (s *Set[T]) AddVote(blockHash common.Hash, vote T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	sender := vote.Sender()
	// Check first if we already received a message from this pal.
	if _, ok := s.messages[sender]; ok {
		// TODO : double signing fault ! Accountability
		return
	}

	var addressesMap map[common.Address]T

	if _, ok := s.Votes[blockHash]; !ok {
		s.Votes[blockHash] = make(map[common.Address]T)
	}

	addressesMap = s.Votes[blockHash]
	addressesMap[sender] = vote
	s.messages[sender] = vote
}

func (s *Set[T]) Messages() []Message {
	s.lock.RLock()
	defer s.lock.RUnlock()
	result := make([]Message, len(s.messages))
	k := 0
	for _, v := range s.messages {
		result[k] = v
		k++
	}
	return result
}

func (s *Set[T]) VotePower(h common.Hash) *big.Int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if votes, ok := s.Votes[h]; ok {
		power := new(big.Int)
		for _, v := range votes {
			power.Add(power, v.Power())
		}
		return power
	}
	return new(big.Int)
}

func (s *Set[T]) TotalVotePower() *big.Int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	power := new(big.Int)
	for _, msg := range s.messages {
		power.Add(power, msg.Power())
	}
	return power
}

func (s *Set[T]) Values(blockHash common.Hash) []Message {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if _, ok := s.Votes[blockHash]; !ok {
		return nil
	}

	var messages = make([]Message, 0)
	for _, v := range s.Votes[blockHash] {
		messages = append(messages, v)
	}
	return messages
}
