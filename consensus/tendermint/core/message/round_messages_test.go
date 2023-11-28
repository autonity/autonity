package message

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
)

func TestMapReset(t *testing.T) {
	messages := NewMap()
	messages.GetOrCreate(0).AddPrevote(NewPrevote(1, 2, common.Hash{}, defaultSigner).MustVerify(stubVerifier))
	messages.GetOrCreate(1).AddPrecommit(NewPrecommit(1, 2, common.Hash{}, defaultSigner).MustVerify(stubVerifier))
	messages.Reset()
	assert.Equal(t, 0, len(messages.All()))
}

func TestGetOrCreate(t *testing.T) {
	messages := NewMap()
	rm0 := messages.GetOrCreate(0)
	rm1 := messages.GetOrCreate(1)
	rm1.AddPrevote(NewPrevote(1, 2, common.Hash{}, defaultSigner).MustVerify(stubVerifier))
	assert.Equal(t, rm0, messages.GetOrCreate(0))
	assert.Equal(t, rm1, messages.GetOrCreate(1))
}

func TestGetMessages(t *testing.T) {
	messages := NewMap()

	rm0 := messages.GetOrCreate(0)
	rm1 := messages.GetOrCreate(1)
	// let round jump happens.
	rm2 := messages.GetOrCreate(4)

	assert.Equal(t, 3, len(messages.internal))
	assert.Equal(t, 0, len(messages.All()))

	prevoteHash := common.HexToHash("prevoteHash")
	precommitHash := common.HexToHash("precommitHash")
	block := types.NewBlockWithHeader(&types.Header{Number: common.Big1})

	proposal := NewPropose(1, 2, -1, block, defaultSigner).MustVerify(stubVerifier)
	prevote := NewPrevote(1, 2, prevoteHash, defaultSigner).MustVerify(stubVerifier)
	precommit := NewPrecommit(1, 2, precommitHash, defaultSigner).MustVerify(stubVerifier)

	rm0.SetProposal(proposal, false)
	rm0.AddPrevote(prevote)
	rm0.AddPrecommit(precommit)

	rm1.SetProposal(proposal, false)
	rm1.AddPrevote(prevote)
	rm1.AddPrecommit(precommit)

	rm2.SetProposal(proposal, false)
	rm2.AddPrevote(prevote)
	rm2.AddPrecommit(precommit)

	allMessages := messages.All()
	assert.Equal(t, 9, len(allMessages))

	for _, m := range allMessages {
		switch m.Code() {
		case ProposalCode:
			r := bytes.Compare(proposal.Payload(), m.Payload())
			assert.Equal(t, 0, r)
		case PrevoteCode:
			r := bytes.Compare(prevote.Payload(), m.Payload())
			assert.Equal(t, 0, r)
		case PrecommitCode:
			r := bytes.Compare(precommit.Payload(), m.Payload())
			assert.Equal(t, 0, r)
		}
	}
}
