package core

import (
	"sync"

	"github.com/clearmatics/autonity/common"
)

func newMessageSet() messageSet {
	return messageSet{
		// map[proposedBlockHash]map[validatorAddress]vote
		votes: map[common.Hash]map[common.Address]message{},
		// map[validatorAddress]nilvote
		nilvotes: map[common.Address]message{},
	}
}

type messageSet struct {
	votes    map[common.Hash]map[common.Address]message
	nilvotes map[common.Address]message
	mu       *sync.RWMutex
}

func (ms *messageSet) AddVote(blockHash common.Hash, msg message) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var addressesMap map[common.Address]message
	var ok bool

	if addressesMap, ok = ms.votes[blockHash]; !ok {
		return
	}

	if _, ok := addressesMap[msg.Address]; ok {
		return
	}

	addressesMap[msg.Address] = msg
}

func (ms *messageSet) AddNilVote(msg message) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if _, ok := ms.nilvotes[msg.Address]; ok {
		ms.nilvotes[msg.Address] = msg
	}
}

func (ms *messageSet) VotesSize(h common.Hash) int {
	if m, ok := ms.votes[h]; ok {
		return len(m)
	}
	return 0
}

func (ms *messageSet) NilVotesSize() int {
	return len(ms.nilvotes)
}

func (ms *messageSet) TotalSize(blockHash common.Hash) int {
	total := len(ms.nilvotes)
	if _, ok := ms.votes[blockHash]; !ok {
		return total
	}

	return total + len(ms.votes[blockHash])
}

func (ms *messageSet) Values(blockHash common.Hash) []message {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	if _, ok := ms.votes[blockHash]; !ok {
		return nil
	}

	var messages []message
	for _, v := range ms.votes[blockHash] {
		messages = append(messages, v)
	}
	return messages
}

//
//func (ms *messageSet) Size() int {
//	ms.mu.RLock()
//	defer ms.mu.RUnlock()
//	return len(ms.messages)
//}
