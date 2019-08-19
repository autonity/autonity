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
)

func newMessageSet() messageSet {
	return messageSet{
		votes:    map[common.Hash]map[common.Address]message{},
		nilvotes: map[common.Address]message{},
	}
}

type messageSet struct {
	votes    map[common.Hash]map[common.Address]message // map[proposedBlockHash]map[validatorAddress]vote
	nilvotes map[common.Address]message                 // map[validatorAddress]vote
}

func (ms *messageSet) AddVote(blockHash common.Hash, msg message) {
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
	if _, ok := ms.nilvotes[msg.Address]; !ok {
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

func (ms *messageSet) TotalSize() int {
	total := ms.NilVotesSize()

	for _, v := range ms.votes {
		total = total + len(v)
	}

	return total
}

func (ms *messageSet) Values(blockHash common.Hash) []message {
	if _, ok := ms.votes[blockHash]; !ok {
		return nil
	}

	var messages = make([]message, 0)
	for _, v := range ms.votes[blockHash] {
		messages = append(messages, v)
	}
	return messages
}
