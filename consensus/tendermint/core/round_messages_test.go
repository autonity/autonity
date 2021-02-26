package core

import (
	"bytes"
	"github.com/clearmatics/autonity/common"
	"gotest.tools/assert"
	"testing"
)

func TestMessagesMap_newMessageMap(t *testing.T) {
	messagesMap := newMessagesMap()
	assert.Equal(t, 0, len(messagesMap.internal))
}

func TestMessagesMap_reset(t *testing.T) {
	messagesMap := newMessagesMap()
	messagesMap.getOrCreate(0)
	messagesMap.getOrCreate(1)
	messagesMap.reset()
	assert.Equal(t, 0, len(messagesMap.internal))
}

func TestMessagesMap_getOrCreate(t *testing.T) {
	messagesMap := newMessagesMap()
	rm0 := messagesMap.getOrCreate(0)
	rm1 := messagesMap.getOrCreate(1)

	assert.Equal(t, rm0, messagesMap.getOrCreate(0))
	assert.Equal(t, rm1, messagesMap.getOrCreate(1))
	assert.Equal(t, 2, len(messagesMap.internal))
}

func TestMessagesMap_GetMessages(t *testing.T) {
	messagesMap := newMessagesMap()

	rm0 := messagesMap.getOrCreate(0)
	rm1 := messagesMap.getOrCreate(1)
	// let round jump happens.
	rm2 := messagesMap.getOrCreate(4)

	assert.Equal(t, 3, len(messagesMap.internal))
	assert.Equal(t, 0, len(messagesMap.GetMessages()))

	prevoteHash := common.HexToHash("prevoteHash")
	precommitHash := common.HexToHash("precommitHash")

	proposalMsg := &ConsensusMessage{
		Code:          msgProposal,
		Msg:           []byte("proposal"),
		Address:       common.HexToAddress("val1"),
		CommittedSeal: []byte{},
	}

	prevoteMsg := &ConsensusMessage{
		Code:          msgPrevote,
		Msg:           []byte("prevote"),
		Address:       common.HexToAddress("val1"),
		CommittedSeal: []byte{},
	}

	precommitMsg := &ConsensusMessage{
		Code:          msgPrecommit,
		Msg:           []byte("precommit"),
		Address:       common.HexToAddress("val1"),
		CommittedSeal: []byte("committed seal"),
	}

	rm0.SetProposal(&Proposal{}, proposalMsg, false)
	rm0.AddPrevote(prevoteHash, *prevoteMsg)
	rm0.AddPrecommit(precommitHash, *precommitMsg)

	rm1.SetProposal(&Proposal{}, proposalMsg, false)
	rm1.AddPrevote(prevoteHash, *prevoteMsg)
	rm1.AddPrecommit(precommitHash, *precommitMsg)

	rm2.SetProposal(&Proposal{}, proposalMsg, false)
	rm2.AddPrevote(prevoteHash, *prevoteMsg)
	rm2.AddPrecommit(precommitHash, *precommitMsg)

	allMessages := messagesMap.GetMessages()
	assert.Equal(t, 9, len(allMessages))

	for _, m := range allMessages {
		switch m.Code {
		case msgProposal:
			assert.Equal(t, proposalMsg.Code, m.Code)

			r := bytes.Compare(proposalMsg.Msg, m.Msg)
			assert.Equal(t, 0, r)

			r = bytes.Compare(proposalMsg.Address[:], m.Address[:])
			assert.Equal(t, 0, r)

			r = bytes.Compare(proposalMsg.CommittedSeal, m.CommittedSeal)
			assert.Equal(t, 0, r)
		case msgPrevote:
			assert.Equal(t, prevoteMsg.Code, m.Code)

			r := bytes.Compare(prevoteMsg.Msg, m.Msg)
			assert.Equal(t, 0, r)

			r = bytes.Compare(prevoteMsg.Address[:], m.Address[:])
			assert.Equal(t, 0, r)

			r = bytes.Compare(prevoteMsg.CommittedSeal, m.CommittedSeal)
			assert.Equal(t, 0, r)
		case msgPrecommit:
			assert.Equal(t, precommitMsg.Code, m.Code)

			r := bytes.Compare(precommitMsg.Msg, m.Msg)
			assert.Equal(t, 0, r)

			r = bytes.Compare(precommitMsg.Address[:], m.Address[:])
			assert.Equal(t, 0, r)

			r = bytes.Compare(precommitMsg.CommittedSeal, m.CommittedSeal)
			assert.Equal(t, 0, r)
		}
	}
}
