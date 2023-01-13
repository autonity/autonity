package messageutils

import (
	"github.com/autonity/autonity/common"
	"math/big"
	"sync"
)

func NewMessageSet() MessageSet {
	return MessageSet{
		Votes:      map[common.Hash]map[common.Address]Message{},
		messages:   make(map[common.Address]*Message),
		messagesMu: new(sync.RWMutex),
	}
}

type MessageSet struct {
	// In some conditions we might receive prevotes or precommit before
	// receiving a proposal, so we must save received message with differents proposed block hash.
	Votes      map[common.Hash]map[common.Address]Message // map[proposedBlockHash]map[validatorAddress]vote
	messages   map[common.Address]*Message
	messagesMu *sync.RWMutex
}

func (ms *MessageSet) AddVote(blockHash common.Hash, msg Message) {
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()

	// Check first if we already received a message from this pal.
	if _, ok := ms.messages[msg.Address]; ok {
		// TODO : double signing fault ! Accountability
		return
	}

	var addressesMap map[common.Address]Message

	if _, ok := ms.Votes[blockHash]; !ok {
		ms.Votes[blockHash] = make(map[common.Address]Message)
	}

	addressesMap = ms.Votes[blockHash]
	addressesMap[msg.Address] = msg
	ms.messages[msg.Address] = &msg
}

func (ms *MessageSet) GetMessages() []*Message {
	ms.messagesMu.RLock()
	defer ms.messagesMu.RUnlock()
	result := make([]*Message, len(ms.messages))
	k := 0
	for _, v := range ms.messages {
		result[k] = v
		k++
	}
	return result
}

func (ms *MessageSet) VotePower(h common.Hash) *big.Int {
	ms.messagesMu.RLock()
	defer ms.messagesMu.RUnlock()
	if msgMap, ok := ms.Votes[h]; ok {
		power := new(big.Int)
		for _, msg := range msgMap {
			power.Add(power, msg.GetPower())
		}
		return power
	}
	return new(big.Int)
}

func (ms *MessageSet) TotalVotePower() *big.Int {
	ms.messagesMu.RLock()
	defer ms.messagesMu.RUnlock()
	power := new(big.Int)
	for _, msg := range ms.messages {
		power.Add(power, msg.GetPower())
	}
	return power
}

func (ms *MessageSet) Values(blockHash common.Hash) []Message {
	ms.messagesMu.RLock()
	defer ms.messagesMu.RUnlock()
	if _, ok := ms.Votes[blockHash]; !ok {
		return nil
	}

	var messages = make([]Message, 0)
	for _, v := range ms.Votes[blockHash] {
		messages = append(messages, v)
	}
	return messages
}
