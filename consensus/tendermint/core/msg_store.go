package core

import (
	"math/big"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
)

var NilValue = common.Hash{}

type MsgStore struct {
	sync.RWMutex
	// the first height that msg are buffered from after node is start.
	firstHeight uint64
	proposals   map[uint64][]*message.Propose
	prevotes    map[uint64][]*message.Prevote
	precommits  map[uint64][]*message.Precommit

	// in the fault detector we only do power computation on prevotes, therefore cache only prevote power
	prevotesPower map[uint64]map[int64]map[common.Hash]*message.PowerInfo
}

func NewMsgStore() *MsgStore {
	return &MsgStore{
		RWMutex:       sync.RWMutex{},
		firstHeight:   uint64(0),
		proposals:     make(map[uint64][]*message.Propose),
		prevotes:      make(map[uint64][]*message.Prevote),
		precommits:    make(map[uint64][]*message.Precommit),
		prevotesPower: make(map[uint64]map[int64]map[common.Hash]*message.PowerInfo),
	}
}

// Save store msg into msg store, it assumes the msg signature was verified, and there is no duplicated msg in the store.
func (ms *MsgStore) Save(m message.Msg) {
	ms.Lock()
	defer ms.Unlock()

	height := m.H()

	if ms.firstHeight == uint64(0) {
		ms.firstHeight = height
	}

	switch msg := m.(type) {
	case *message.Propose:
		_, ok := ms.proposals[height]
		if !ok {
			ms.proposals[height] = make([]*message.Propose, 0)
		}
		ms.proposals[height] = append(ms.proposals[height], msg)
	case *message.Prevote:
		_, ok := ms.prevotes[height]
		if !ok {
			ms.prevotes[height] = make([]*message.Prevote, 0)
		}
		ms.prevotes[height] = append(ms.prevotes[height], msg)

		// update prevotes power cache
		round := msg.R()
		value := msg.Value()
		_, ok = ms.prevotesPower[height]
		if !ok {
			ms.prevotesPower[height] = make(map[int64]map[common.Hash]*message.PowerInfo)
		}
		_, ok = ms.prevotesPower[height][round]
		if !ok {
			ms.prevotesPower[height][round] = make(map[common.Hash]*message.PowerInfo)
		}
		_, ok = ms.prevotesPower[height][round][value]
		if !ok {
			ms.prevotesPower[height][round][value] = message.NewPowerInfo()
		}
		for index, power := range msg.Signers().Powers() {
			ms.prevotesPower[height][round][value].Set(index, power)
		}
	case *message.Precommit:
		_, ok := ms.precommits[height]
		if !ok {
			ms.precommits[height] = make([]*message.Precommit, 0)
		}
		ms.precommits[height] = append(ms.precommits[height], msg)
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
	for h := range ms.proposals {
		if h <= height {
			delete(ms.proposals, h)
		}
	}
	for h := range ms.prevotes {
		if h <= height {
			delete(ms.prevotes, h)
		}
	}
	for h := range ms.precommits {
		if h <= height {
			delete(ms.precommits, h)
		}
	}
	for h := range ms.prevotesPower {
		if h <= height {
			delete(ms.prevotesPower, h)
		}
	}
}

// RemoveMsg only used for integration tests.
func (ms *MsgStore) RemoveMsg(height uint64, code uint8, hash common.Hash) {
	ms.Lock()
	defer ms.Unlock()

	switch code {
	case message.ProposalCode:
		_, ok := ms.proposals[height]
		if !ok {
			return
		}
		var filteredProposals []*message.Propose
		for _, proposal := range ms.proposals[height] {
			if proposal.Hash() != hash {
				filteredProposals = append(filteredProposals, proposal)
			}
		}
		ms.proposals[height] = filteredProposals
	case message.PrevoteCode:
		_, ok := ms.prevotes[height]
		if !ok {
			return
		}
		var filteredPrevotes []*message.Prevote
		for _, prevote := range ms.prevotes[height] {
			if prevote.Hash() != hash {
				filteredPrevotes = append(filteredPrevotes, prevote)
			}
		}
		ms.prevotes[height] = filteredPrevotes

		// update power cache
		ms.prevotesPower = make(map[uint64]map[int64]map[common.Hash]*message.PowerInfo)
		for _, msg := range ms.prevotes[height] {
			round := msg.R()
			value := msg.Value()
			_, ok = ms.prevotesPower[height]
			if !ok {
				ms.prevotesPower[height] = make(map[int64]map[common.Hash]*message.PowerInfo)
			}
			_, ok = ms.prevotesPower[height][round]
			if !ok {
				ms.prevotesPower[height][round] = make(map[common.Hash]*message.PowerInfo)
			}
			_, ok = ms.prevotesPower[height][round][value]
			if !ok {
				ms.prevotesPower[height][round][value] = message.NewPowerInfo()
			}
			for index, power := range msg.Signers().Powers() {
				ms.prevotesPower[height][round][value].Set(index, power)
			}
		}
	case message.PrecommitCode:
		_, ok := ms.precommits[height]
		if !ok {
			return
		}
		var filteredPrecommits []*message.Precommit
		for _, precommit := range ms.precommits[height] {
			if precommit.Hash() != hash {
				filteredPrecommits = append(filteredPrecommits, precommit)
			}
		}
		ms.precommits[height] = filteredPrecommits
	default:
		panic("non-existent code")
	}
}

func (ms *MsgStore) GetProposals(height uint64, query func(*message.Propose) bool) []*message.Propose {
	ms.RLock()
	defer ms.RUnlock()

	var result []*message.Propose
	_, ok := ms.proposals[height]
	if !ok {
		return result
	}

	for _, proposal := range ms.proposals[height] {
		if query(proposal) {
			result = append(result, proposal)
		}
	}
	return result
}

func (ms *MsgStore) GetPrevotes(height uint64, query func(*message.Prevote) bool) []*message.Prevote {
	ms.RLock()
	defer ms.RUnlock()

	var result []*message.Prevote
	_, ok := ms.prevotes[height]
	if !ok {
		return result
	}

	for _, prevote := range ms.prevotes[height] {
		if query(prevote) {
			result = append(result, prevote)
		}
	}
	return result
}

func (ms *MsgStore) GetPrecommits(height uint64, query func(*message.Precommit) bool) []*message.Precommit {
	ms.RLock()
	defer ms.RUnlock()

	var result []*message.Precommit
	_, ok := ms.precommits[height]
	if !ok {
		return result
	}

	for _, precommit := range ms.precommits[height] {
		if query(precommit) {
			result = append(result, precommit)
		}
	}
	return result
}

func (ms *MsgStore) PrevotesPowerFor(height uint64, round int64, value common.Hash) *big.Int {
	ms.RLock()
	defer ms.RUnlock()

	_, ok := ms.prevotesPower[height]
	if !ok {
		return new(big.Int)
	}
	_, ok = ms.prevotesPower[height][round]
	if !ok {
		return new(big.Int)
	}
	_, ok = ms.prevotesPower[height][round][value]
	if !ok {
		return new(big.Int)
	}
	return new(big.Int).Set(ms.prevotesPower[height][round][value].Pow()) // return a copy to avoid data races
}

// this function checks if we have a quorum for a value in (h,r). It excludes the `excludedValue` from the search.
// it is used by the fault detector to verify if we have quorums of prevotes for values != `excludedValue`.
// returns the slice of messages constituting the quorum
func (ms *MsgStore) SearchQuorum(height uint64, round int64, excludedValue common.Hash, quorum *big.Int) []message.Msg {
	ms.Lock()
	defer ms.Unlock()

	var result []message.Msg

	_, ok := ms.prevotesPower[height]
	if !ok {
		return result
	}
	_, ok = ms.prevotesPower[height][round]
	if !ok {
		return result
	}

	for value, powerInfo := range ms.prevotesPower[height][round] {
		if value == excludedValue {
			continue
		}
		if powerInfo.Pow().Cmp(quorum) >= 0 {
			_, ok := ms.prevotes[height]
			if !ok {
				panic("Have quorum in power cache, but cannot find related messages in msgStore")
			}
			for _, prevote := range ms.prevotes[height] {
				if prevote.R() == round && prevote.Value() == value {
					result = append(result, prevote)
				}
			}
			return result
		}
	}

	return result
}
