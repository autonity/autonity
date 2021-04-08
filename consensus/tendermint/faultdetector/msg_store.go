package faultdetector

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/core/types"
	"sync"
)

type MsgStore struct {
	sync.RWMutex
	// map[Height]map[Round]map[MsgType]map[common.address][]*Message
	messages map[uint64]map[int64]map[uint64]map[common.Address][]*core.Message
}

func newMsgStore() *MsgStore {
	return &MsgStore{
		RWMutex:  sync.RWMutex{},
		messages: make(map[uint64]map[int64]map[uint64]map[common.Address][]*core.Message)}
}

// store msg into msg store, it returns msgs that is equivocation than the input msg, and an errEquivocation.
// otherwise it return nil, nil
func (ms *MsgStore) Save(m *core.Message) ([]*core.Message, error) {
	ms.Lock()
	defer ms.Unlock()

	height, _ := m.Height()
	roundMap, ok := ms.messages[height.Uint64()]
	if !ok {
		roundMap = make(map[int64]map[uint64]map[common.Address][]*core.Message)
		ms.messages[height.Uint64()] = roundMap
	}

	round, _ := m.Round()
	msgTypeMap, ok := roundMap[round]
	if !ok {
		msgTypeMap = make(map[uint64]map[common.Address][]*core.Message)
		roundMap[round] = msgTypeMap
	}

	addressMap, ok := msgTypeMap[m.Code]
	if !ok {
		addressMap = make(map[common.Address][]*core.Message)
		msgTypeMap[m.Code] = addressMap
	}

	msgs, ok := addressMap[m.Address]
	if !ok {
		var msgList []*core.Message
		addressMap[m.Address] = append(msgList, m)
		return nil, nil
	}

	presented := false
	for i := 0; i < len(msgs); i++ {
		if types.RLPHash(msgs[i].Payload()) == types.RLPHash(m.Payload()) {
			presented = true
		}
	}

	if !presented {
		addressMap[m.Address] = append(msgs, m)
		return msgs, errEquivocation
	}

	return nil, nil
}

func (ms *MsgStore) DeleteMsgsAtHeight(height uint64) {
	ms.Lock()
	defer ms.Unlock()

	// Remove all messgages at this height
	for round, roundMsgMap := range ms.messages[height] {
		for code, typeMsgMap := range roundMsgMap {
			for addr, _ := range typeMsgMap { // nolint
				delete(ms.messages[height][round][code], addr)
			}
			delete(ms.messages[height][round], code)
		}
		delete(ms.messages[height], round)
	}
	// Delete map entry for this height
	delete(ms.messages, height)
}

// get take height and query conditions to query those msgs from msg store, it returns those msgs satisfied the condition.
func (ms *MsgStore) Get(height uint64, query func(*core.Message) bool) []*core.Message {
	ms.RLock()
	defer ms.RUnlock()

	var result []*core.Message
	roundMap, ok := ms.messages[height]
	if !ok {
		return result
	}

	for _, msgTypeMap := range roundMap {
		for _, addressMap := range msgTypeMap {
			for _, msgs := range addressMap {
				for i := 0; i < len(msgs); i++ {
					if query(msgs[i]) {
						result = append(result, msgs[i])
					}
				}
			}
		}
	}

	return result
}
