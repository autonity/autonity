package core

import (
	"github.com/autonity/autonity/core/types"

	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
)

var NilValue = common.Hash{}

type MsgStore struct {
	sync.RWMutex
	// the first height that msg are buffered from after node is start.
	firstHeight uint64
	// To keep a more flexible query interface and a better query performance for msg store,
	// we'd need to keep the legacy data schema for msg store, thus the save msg function
	// need to save duplicated pointers of aggregated votes.
	// map[Height]map[Round]map[Step]map[Signer][]*Message
	messages map[uint64]map[int64]map[uint8]map[common.Address][]message.Msg
}

func NewMsgStore() *MsgStore {
	return &MsgStore{
		RWMutex:     sync.RWMutex{},
		firstHeight: uint64(0),
		messages:    make(map[uint64]map[int64]map[uint8]map[common.Address][]message.Msg)}
}

// Save store msg into msg store, it assumes the msg signature was verified, and there is no duplicated msg in the store.
func (ms *MsgStore) Save(m message.Msg, committee types.Committee) {
	ms.Lock()
	defer ms.Unlock()
	if ms.firstHeight == uint64(0) {
		ms.firstHeight = m.H()
	}

	height := m.H()
	roundMap, ok := ms.messages[height]

	if !ok {
		roundMap = make(map[int64]map[uint8]map[common.Address][]message.Msg)
		ms.messages[height] = roundMap
	}

	round := m.R()
	msgTypeMap, ok := roundMap[round]

	if !ok {
		msgTypeMap = make(map[uint8]map[common.Address][]message.Msg)
		roundMap[round] = msgTypeMap
	}

	addressMap, ok := msgTypeMap[m.Code()]
	if !ok {
		addressMap = make(map[common.Address][]message.Msg)
		msgTypeMap[m.Code()] = addressMap
	}

	// as proposal is not aggregatable, save it and return
	if m.Code() == message.ProposalCode {
		signer := m.(*message.Propose).Signer()
		msgs, ok := addressMap[signer]
		if !ok {
			var msgList []message.Msg
			addressMap[signer] = append(msgList, m)
			return
		}
		addressMap[signer] = append(msgs, m)
		return
	}

	// for vote, save vote for each signer.
	signers := m.(message.Vote).Signers()
	for _, valIndex := range signers.FlattenUniq() {
		signer := committee[valIndex].Address
		msgs, ok := addressMap[signer]
		if !ok {
			var msgList []message.Msg
			addressMap[signer] = append(msgList, m)
			return
		}
		addressMap[signer] = append(msgs, m)
	}
}

func (ms *MsgStore) FirstHeightBuffered() uint64 {
	ms.Lock()
	defer ms.Unlock()
	return ms.firstHeight
}

func (ms *MsgStore) DeleteOlds(height uint64) {
	ms.Lock()
	defer ms.Unlock()
	for h := range ms.messages {
		if h <= height {
			// Delete map entry for this height
			delete(ms.messages, h)
		}
	}
}

// Get take height and query conditions to query those msgs from msg store, it returns those msgs satisfied the condition.
func (ms *MsgStore) Get(query func(message.Msg) bool, height uint64, signers ...common.Address) []message.Msg {
	ms.RLock()
	defer ms.RUnlock()

	var result []message.Msg
	roundMap, ok := ms.messages[height]
	if !ok {
		return result
	}

	// querying without the signer address nominated, as it iterates all signers, we'd need to filter duplicated ones.
	if len(signers) == 0 {
		msgHashMap := make(map[common.Hash]struct{})
		for _, msgTypeMap := range roundMap {
			for _, addressMap := range msgTypeMap {
				for _, msgs := range addressMap {
					for _, msg := range msgs {
						if query(msg) {
							if _, ok := msgHashMap[msg.Hash()]; !ok {
								result = append(result, msg)
								msgHashMap[msg.Hash()] = struct{}{}
							}
						}
					}
				}
			}
		}
		return result
	}

	// querying with signer address
	signer := signers[0]
	for _, msgTypeMap := range roundMap {
		for _, addressMap := range msgTypeMap {
			messages, ok := addressMap[signer]
			if !ok {
				continue
			}

			for _, msg := range messages {
				if query(msg) {
					result = append(result, msg)
				}
			}
		}
	}
	return result
}

func (ms *MsgStore) GetEquivocatedVotes(height uint64, round int64, step uint8, signer common.Address, value common.Hash) []message.Msg {
	ms.RLock()
	defer ms.RUnlock()
	var result []message.Msg
	roundMap, ok := ms.messages[height]
	if !ok {
		return result
	}
	stepMap, ok := roundMap[round]
	if !ok {
		return result
	}

	signerMap, ok := stepMap[step]
	if !ok {
		return result
	}

	messages, ok := signerMap[signer]
	if !ok {
		return result
	}

	for _, m := range messages {
		if m.Value() != value {
			result = append(result, m)
		}
	}

	return result
}

func GetStore[T any, PT interface {
	*T
	message.Msg
}](ms *MsgStore, height uint64, query func(*T) bool) []*T {
	ms.RLock()
	defer ms.RUnlock()
	var result []*T
	roundMap, ok := ms.messages[height]
	if !ok {
		return result
	}
	code := PT(new(T)).Code()
	for _, msgTypeMap := range roundMap {
		for _, msgs := range msgTypeMap[code] {
			for _, msg := range msgs {
				if m, ok := msg.(PT); ok {
					if query((*T)(m)) {
						result = append(result, (*T)(m))
					}
				}
			}
		}
	}
	return result
}
