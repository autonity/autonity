// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"github.com/clearmatics/autonity/common"
	"sync"
)

func newMessageSet() *messageSet {
	return &messageSet{
		votes:      map[common.Hash]map[common.Address]Message{},
		messagesMu: new(sync.RWMutex),
	}
}

// empty block hash (common.Hash{}) is for nil votes
type messageSet struct {
	votes      map[common.Hash]map[common.Address]Message // map[proposedBlockHash]map[validatorAddress]vote
	messagesMu *sync.RWMutex
}

func newProposalSet(p Proposal, m *Message) *proposalSet {
	return &proposalSet{
		p:    p,
		pMsg: m,
		mu:   new(sync.RWMutex),
	}
}

type proposalSet struct {
	p    Proposal
	pMsg *Message
	mu   *sync.RWMutex
}

func (ms *messageSet) Add(hash common.Hash, msg Message) {
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()

	var addressesMap map[common.Address]Message
	var ok bool

	if _, ok = ms.votes[hash]; !ok {
		ms.votes[hash] = make(map[common.Address]Message)
	}

	addressesMap = ms.votes[hash]

	if _, ok := addressesMap[msg.Address]; ok {
		return
	}

	addressesMap[msg.Address] = msg

}

func (ms *messageSet) GetMessages() []*Message {
	ms.messagesMu.RLock()
	defer ms.messagesMu.RUnlock()

	var result []*Message

	for _, addressMap := range ms.votes {
		for _, msg := range addressMap {
			result = append(result, &msg)
		}
	}
	return result
}

func (ms *messageSet) VotesSize(h common.Hash) int {
	ms.messagesMu.RLock()
	defer ms.messagesMu.RUnlock()

	if m, ok := ms.votes[h]; ok {
		return len(m)
	}
	return 0
}

func (ms *messageSet) NilVotesSize() int {
	ms.messagesMu.RLock()
	defer ms.messagesMu.RUnlock()

	return len(ms.votes[common.Hash{}])
}

func (ms *messageSet) TotalSize() int {
	ms.messagesMu.RLock()
	defer ms.messagesMu.RUnlock()

	var total int
	for _, v := range ms.votes {
		total = total + len(v)
	}
	return total
}

func (ms *messageSet) AllBlockHashMessages(blockHash common.Hash) []Message {
	ms.messagesMu.RLock()
	defer ms.messagesMu.RUnlock()

	if _, ok := ms.votes[blockHash]; !ok {
		return nil
	}

	var messages = make([]Message, 0)
	for _, v := range ms.votes[blockHash] {
		messages = append(messages, v)
	}

	var result = make([]Message, len(messages))
	copy(result, messages)
	return result
}

func (ms *messageSet) hasMessage(h common.Hash, m Message) bool {
	ms.messagesMu.RLock()
	defer ms.messagesMu.RUnlock()

	var addressesMap map[common.Address]Message
	var ok bool

	if addressesMap, ok = ms.votes[h]; !ok {
		return false
	}

	if _, ok = addressesMap[m.Address]; !ok {
		return false
	}

	return true
}

func (ps *proposalSet) proposal() Proposal {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.p
}

func (ps *proposalSet) proposalMsg() *Message {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.pMsg
}
