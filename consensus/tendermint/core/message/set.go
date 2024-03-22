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
		//indexes: make([]uint64, 100), //TODO(lorenzo) fix size
	}
}

type Set[T isVote] struct {
	// In some conditions we might receive prevotes or precommit before
	// receiving a proposal, so we must save received message with different proposed block hash.
	votes map[common.Hash][]T // map[proposedBlockHash][]vote
	//indexes []uint64
	sync.RWMutex
}

func (s *Set[T]) Add(vote T) {
	s.Lock()
	defer s.Unlock()

	/*
		// Check first if we already received a message from this pal.
		sender := vote.Sender()
		if _, ok := s.messages[sender]; ok {
			// TODO : double signing fault ! Accountability
			return
		}*/

	value := vote.Value()

	if _, ok := s.votes[value]; !ok {
		s.votes[value] = make([]T, 0) //TODO(lorenzo) fix size
	}
	s.votes[value] = append(s.votes[value], vote)
	//s.messages[sender] = vote
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

	if votes, ok := s.votes[h]; ok {
		power := new(big.Int)
		accountedFor := make(map[common.Address]struct{})
		for _, v := range votes {
			senders := v.Senders()
			powers := v.Powers()
			//TODO(lorenzo) twisted logic but should work
			for i, _ := range v.SendersCoeff().FlattenUniq() {
				_, accounted := accountedFor[senders[i]]
				if accounted {
					continue
				}
				power.Add(power, powers[i])
				accountedFor[senders[i]] = struct{}{}
			}
		}
		return power
	}
	return new(big.Int)
}

func (s *Set[T]) TotalPower() *big.Int {
	s.RLock()
	defer s.RUnlock()
	power := new(big.Int)
	accountedFor := make(map[common.Address]struct{})
	for _, vts := range s.votes {
		for _, v := range vts {
			senders := v.Senders()
			powers := v.Powers()
			//TODO(lorenzo) twisted logic but should work
			for i, _ := range v.SendersCoeff().FlattenUniq() {
				_, accounted := accountedFor[senders[i]]
				if accounted {
					continue
				}
				power.Add(power, powers[i])
				accountedFor[senders[i]] = struct{}{}
			}
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
	//TODO(lorenzo) might not need copy here
	messages := make([]T, 0, len(s.votes[blockHash]))
	for _, v := range s.votes[blockHash] {
		messages = append(messages, v)
	}
	return messages
}
