package faultdetector

import (
	"github.com/clearmatics/autonity/common"
	tdm "github.com/clearmatics/autonity/consensus/tendermint"
	algo "github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	"github.com/clearmatics/autonity/core/types"
)

type MsgStore struct {
	// map[Height]map[Round]map[MsgType]map[common.address][]*Message
	messages map[uint64]map[int64]map[algo.Step]map[common.Address][]*tdm.Message
}

func newMsgStore() *MsgStore {
	return &MsgStore{messages: make(map[uint64]map[int64]map[algo.Step]map[common.Address][]*tdm.Message)}
}

// store msg into msg store, it returns msgs that is equivocation than the input msg, and an errEquivocation.
// otherwise it return nil, nil
func (ms *MsgStore) Save(m *tdm.Message) ([]*tdm.Message, error) {
	height := m.H()
	roundMap, ok := ms.messages[height]
	if !ok {
		roundMap = make(map[int64]map[algo.Step]map[common.Address][]*tdm.Message)
		ms.messages[height] = roundMap
	}

	round := m.R()
	msgTypeMap, ok := roundMap[round]
	if !ok {
		msgTypeMap = make(map[algo.Step]map[common.Address][]*tdm.Message)
		roundMap[round] = msgTypeMap
	}

	addressMap, ok := msgTypeMap[m.Type()]
	if !ok {
		addressMap = make(map[common.Address][]*tdm.Message)
		msgTypeMap[m.Type()] = addressMap
	}

	msgs, ok := addressMap[m.Address]
	if !ok {
		var msgList []*tdm.Message
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
func (ms *MsgStore) Get(height uint64, query func(*tdm.Message) bool) []tdm.Message {

	var result []tdm.Message
	roundMap, ok := ms.messages[height]
	if !ok {
		return result
	}

	for _, msgTypeMap := range roundMap {
		for _, addressMap := range msgTypeMap {
			for _, msgs := range addressMap {
				for i := 0; i < len(msgs); i++ {
					if query(msgs[i]) {
						result = append(result, *msgs[i])
					}
				}
			}
		}
	}

	return result
}
