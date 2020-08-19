package core

import (
	"fmt"

	common "github.com/clearmatics/autonity/common"
	"github.com/davecgh/go-spew/spew"
)

// Message cache caches messages

type messageCache struct {
	proposals   map[uint64]map[int64]map[common.Address]*Proposal
	prevotes    map[uint64]map[int64]map[common.Address]*Vote
	precommits  map[uint64]map[int64]map[common.Address]*Vote
	rawMessages map[common.Hash]Message

	// TODO use this to impose upper and lower limits on what can go in the
	// message cache.
	currentHeight uint64
}

func (m *messageCache) addMessage(msg Message) error {
	switch msg.Code {
	case msgProposal:
		var p Proposal
		err := msg.Decode(&p)
		if err != nil {
			return err
		}
		roundMap, ok := m.proposals[p.Height.Uint64()]
		if !ok {
			roundMap = make(map[int64]map[common.Address]*Proposal)
		}
		msgMap, ok := roundMap[p.Round]
		if !ok {
			msgMap = make(map[common.Address]*Proposal)
		}
		msgMap[msg.Address] = &p
	case msgPrevote:
		var preV Vote
		err := msg.Decode(&preV)
		if err != nil {
			return err
		}
		addVote(m.prevotes, &preV, msg.Address)
	case msgPrecommit:
		var preC Vote
		err := msg.Decode(&preC)
		if err != nil {
			return err
		}
		addVote(m.prevotes, &preC, msg.Address)
	default:
		panic(fmt.Sprintf("Unrecognised message code %d for message: %s", msg.Code, spew.Sdump(msg)))
	}
	return nil
}

func (m *messageCache) getproposals(height uint64, round int64) map[common.Address]*Proposal {
	roundMap, ok := m.proposals[height]
	if !ok {
		return nil
	}
	return roundMap[round]
}

func (m *messageCache) getprevotes(height uint64, round int64) map[common.Address]*Vote {
	return getVotes(height, round, m.prevotes)
}

func (m *messageCache) getprecommits(height uint64, round int64) map[common.Address]*Vote {
	return getVotes(height, round, m.precommits)
}

func getVotes(height uint64, round int64, votes map[uint64]map[int64]map[common.Address]*Vote) map[common.Address]*Vote {
	roundMap, ok := votes[height]
	if !ok {
		return nil
	}
	return roundMap[round]
}

func addVote(votes map[uint64]map[int64]map[common.Address]*Vote, v *Vote, address common.Address) {
	roundMap, ok := votes[v.Height.Uint64()]
	if !ok {
		roundMap = make(map[int64]map[common.Address]*Vote)
	}
	msgMap, ok := roundMap[v.Round]
	if !ok {
		msgMap = make(map[common.Address]*Vote)
	}
	msgMap[address] = v
}
