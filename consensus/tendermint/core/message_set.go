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
	"github.com/clearmatics/autonity/core/types"
	"sync"
)

func newMessageSet() messageSet {
	return messageSet{
		votes:      map[common.Hash]map[common.Address]types.ConsensusMessage{},
		messages:   make(map[common.Address]*types.ConsensusMessage),
		messagesMu: new(sync.RWMutex),
	}
}

type messageSet struct {
	// In some conditions we might receive prevotes or precommit before
	// receiving a proposal, so we must save received message with differents proposed block hash.
	votes      map[common.Hash]map[common.Address]types.ConsensusMessage // map[proposedBlockHash]map[validatorAddress]vote
	messages   map[common.Address]*types.ConsensusMessage
	messagesMu *sync.RWMutex
}

func (ms *messageSet) AddVote(blockHash common.Hash, msg types.ConsensusMessage) {
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()

	// Check first if we already received a message from this pal.
	if _, ok := ms.messages[msg.Address]; ok {
		// TODO : double signing fault ! Accountability
		return
	}

	var addressesMap map[common.Address]types.ConsensusMessage

	if _, ok := ms.votes[blockHash]; !ok {
		ms.votes[blockHash] = make(map[common.Address]types.ConsensusMessage)
	}

	addressesMap = ms.votes[blockHash]
	addressesMap[msg.Address] = msg
	ms.messages[msg.Address] = &msg
}

func (ms *messageSet) GetMessages() []*types.ConsensusMessage {
	ms.messagesMu.RLock()
	defer ms.messagesMu.RUnlock()
	result := make([]*types.ConsensusMessage, len(ms.messages))
	k := 0
	for _, v := range ms.messages {
		result[k] = v
		k++
	}
	return result
}

func (ms *messageSet) VotePower(h common.Hash) uint64 {
	ms.messagesMu.RLock()
	defer ms.messagesMu.RUnlock()
	if msgMap, ok := ms.votes[h]; ok {
		var power uint64
		for _, msg := range msgMap {
			power += msg.GetPower()
		}
		return power
	}
	return 0
}

func (ms *messageSet) TotalVotePower() uint64 {
	ms.messagesMu.RLock()
	defer ms.messagesMu.RUnlock()
	var power uint64
	for _, msg := range ms.messages {
		power += msg.GetPower()
	}
	return power
}

func (ms *messageSet) Values(blockHash common.Hash) []types.ConsensusMessage {
	ms.messagesMu.RLock()
	defer ms.messagesMu.RUnlock()
	if _, ok := ms.votes[blockHash]; !ok {
		return nil
	}

	var messages = make([]types.ConsensusMessage, 0)
	for _, v := range ms.votes[blockHash] {
		messages = append(messages, v)
	}
	return messages
}
