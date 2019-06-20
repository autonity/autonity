package core

import (
	"sync"

	"github.com/clearmatics/autonity/common"
)

func newMessageSet() messageSet {
	return messageSet{
		votes:    map[common.Hash]map[common.Address]message{},
		nilvotes: map[common.Address]message{},
		mu:       new(sync.RWMutex),
	}
}

type messageSet struct {
	votes    map[common.Hash]map[common.Address]message // map[proposedBlockHash]map[validatorAddress]vote
	nilvotes map[common.Address]message                 // map[validatorAddress]vote
	mu       *sync.RWMutex
}

func (ms *messageSet) AddVote(blockHash common.Hash, msg message) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var addressesMap map[common.Address]message
	var ok bool

	if _, ok = ms.votes[blockHash]; !ok {
		ms.votes[blockHash] = make(map[common.Address]message)
	}

	addressesMap = ms.votes[blockHash]

	if _, ok := addressesMap[msg.Address]; ok {
		return
	}

	addressesMap[msg.Address] = msg
}

func (ms *messageSet) AddNilVote(msg message) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if _, ok := ms.nilvotes[msg.Address]; !ok {
		ms.nilvotes[msg.Address] = msg
	}
}

func (ms *messageSet) VotesSize(h common.Hash) int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	if m, ok := ms.votes[h]; ok {
		return len(m)
	}
	return 0
}

func (ms *messageSet) NilVotesSize() int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return len(ms.nilvotes)
}

func (ms *messageSet) TotalSize() int {
	total := ms.NilVotesSize()

	ms.mu.RLock()
	defer ms.mu.RUnlock()
	for _, v := range ms.votes {
		total = total + len(v)
	}

	return total
}

func (ms *messageSet) Values(blockHash common.Hash) []message {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	if _, ok := ms.votes[blockHash]; !ok {
		return nil
	}

	var messages = make([]message, 0)
	for _, v := range ms.votes[blockHash] {
		messages = append(messages, v)
	}
	return messages
}
