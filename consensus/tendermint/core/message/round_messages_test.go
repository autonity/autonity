package message

import (
	"bytes"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"gotest.tools/assert"
	"testing"
)

func TestMessagesMap_newMessageMap(t *testing.T) {
	messagesMap := NewMap()
	assert.Equal(t, 0, len(messagesMap.internal))
}

func TestMessagesMap_reset(t *testing.T) {
	messagesMap := NewMap()
	messagesMap.GetOrCreate(0)
	messagesMap.GetOrCreate(1)
	messagesMap.Reset()
	assert.Equal(t, 0, len(messagesMap.internal))
}

func TestMessagesMap_getOrCreate(t *testing.T) {
	messagesMap := NewMap()
	rm0 := messagesMap.GetOrCreate(0)
	rm1 := messagesMap.GetOrCreate(1)

	assert.Equal(t, rm0, messagesMap.GetOrCreate(0))
	assert.Equal(t, rm1, messagesMap.GetOrCreate(1))
	assert.Equal(t, 2, len(messagesMap.internal))
}

func TestMessagesMap_GetMessages(t *testing.T) {
	messagesMap := NewMap()

	rm0 := messagesMap.GetOrCreate(0)
	rm1 := messagesMap.GetOrCreate(1)
	// let round jump happens.
	rm2 := messagesMap.GetOrCreate(4)

	assert.Equal(t, 3, len(messagesMap.internal))
	assert.Equal(t, 0, len(messagesMap.All()))

	prevoteHash := common.HexToHash("prevoteHash")
	precommitHash := common.HexToHash("precommitHash")

	proposalMsg := &Message{
		Code:          consensus.MsgProposal,
		Payload:       []byte("proposal"),
		Address:       common.HexToAddress("val1"),
		CommittedSeal: []byte{},
	}

	prevoteMsg := &Message{
		Code:          consensus.MsgPrevote,
		Payload:       []byte("prevote"),
		Address:       common.HexToAddress("val1"),
		CommittedSeal: []byte{},
	}

	precommitMsg := &Message{
		Code:          consensus.MsgPrecommit,
		Payload:       []byte("precommit"),
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

	allMessages := messagesMap.All()
	assert.Equal(t, 9, len(allMessages))

	for _, m := range allMessages {
		switch m.Code {
		case consensus.MsgProposal:
			assert.Equal(t, proposalMsg.Code, m.Code)

			r := bytes.Compare(proposalMsg.Payload, m.Payload)
			assert.Equal(t, 0, r)

			r = bytes.Compare(proposalMsg.Address[:], m.Address[:])
			assert.Equal(t, 0, r)

			r = bytes.Compare(proposalMsg.CommittedSeal, m.CommittedSeal)
			assert.Equal(t, 0, r)
		case consensus.MsgPrevote:
			assert.Equal(t, prevoteMsg.Code, m.Code)

			r := bytes.Compare(prevoteMsg.Payload, m.Payload)
			assert.Equal(t, 0, r)

			r = bytes.Compare(prevoteMsg.Address[:], m.Address[:])
			assert.Equal(t, 0, r)

			r = bytes.Compare(prevoteMsg.CommittedSeal, m.CommittedSeal)
			assert.Equal(t, 0, r)
		case consensus.MsgPrecommit:
			assert.Equal(t, precommitMsg.Code, m.Code)

			r := bytes.Compare(precommitMsg.Payload, m.Payload)
			assert.Equal(t, 0, r)

			r = bytes.Compare(precommitMsg.Address[:], m.Address[:])
			assert.Equal(t, 0, r)

			r = bytes.Compare(precommitMsg.CommittedSeal, m.CommittedSeal)
			assert.Equal(t, 0, r)
		}
	}
}
