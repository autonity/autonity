package faultdetector

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func newVoteMsg(h uint64, r int64, code uint64, addr common.Address, value common.Hash) *core.Message { // nolint: unparam
	var vote = core.Vote{
		Round:             r,
		Height:            new(big.Int).SetUint64(h),
		ProposedBlockHash: value,
	}

	encodedVote, err := core.Encode(&vote)
	if err != nil {
		return nil
	}

	var msg = core.Message{
		Code:          code,
		Msg:           encodedVote,
		Address:       addr,
		CommittedSeal: []byte{},
	}

	payload, err := msg.PayloadNoSig()
	if err != nil {
		return nil
	}

	m := new(core.Message)
	if err := m.FromPayload(payload); err != nil {
		return nil
	}

	return m
}

func TestMsgStore(t *testing.T) {
	height := uint64(100)
	round := int64(0)
	addrAlice := common.Address{0x1}
	addrBob := common.Address{0x2}
	noneNilValue := common.Hash{0x1}

	t.Run("query msg store when msg store is empty", func(t *testing.T) {
		ms := newMsgStore()
		proposals := ms.Get(height, func(m *core.Message) bool {
			return m.Type() == msgProposal
		})
		assert.Equal(t, 0, len(proposals))
	})

	t.Run("save equivocation msgs in msg store", func(t *testing.T) {
		ms := newMsgStore()
		preVoteNil := newVoteMsg(height, round, msgPrevote, addrAlice, nilValue)
		_, err := ms.Save(preVoteNil)
		if err != nil {
			assert.Error(t, err)
		}

		preVoteNoneNil := newVoteMsg(height, round, msgPrevote, addrAlice, noneNilValue)
		equivocatedMsg, err := ms.Save(preVoteNoneNil)
		assert.NotNil(t, equivocatedMsg)
		assert.Equal(t, err, errEquivocation)
		assert.Equal(t, nilValue, equivocatedMsg.Value())
		assert.Equal(t, addrAlice, equivocatedMsg.Sender())
		assert.Equal(t, height, equivocatedMsg.H())
		assert.Equal(t, round, equivocatedMsg.R())
		assert.Equal(t, msgPrevote, equivocatedMsg.Type())
	})

	t.Run("query a presented preVote from msg store", func(t *testing.T) {
		ms := newMsgStore()
		preVote := newVoteMsg(height, round, msgPrevote, addrAlice, nilValue)
		_, err := ms.Save(preVote)
		if err != nil {
			assert.Error(t, err)
		}

		votes := ms.Get(height, func(m *core.Message) bool {
			return m.Type() == msgPrevote && m.H() == height && m.R() == round && m.Sender() == addrAlice &&
				m.Value() == nilValue
		})

		assert.Equal(t, 1, len(votes))
		assert.Equal(t, msgPrevote, votes[0].Type())
		assert.Equal(t, height, votes[0].H())
		assert.Equal(t, round, votes[0].R())
		assert.Equal(t, addrAlice, votes[0].Sender())
		assert.Equal(t, nilValue, votes[0].Value())
	})

	t.Run("query multiple presented preVote from msg store", func(t *testing.T) {
		ms := newMsgStore()
		preVoteNil := newVoteMsg(height, round, msgPrevote, addrAlice, nilValue)
		_, err := ms.Save(preVoteNil)
		if err != nil {
			assert.Error(t, err)
		}

		preVoteNoneNil := newVoteMsg(height, round, msgPrevote, addrBob, noneNilValue)
		_, err = ms.Save(preVoteNoneNil)
		if err != nil {
			assert.Error(t, err)
		}

		votes := ms.Get(height, func(m *core.Message) bool {
			return m.Type() == msgPrevote && m.H() == height && m.R() == round
		})

		assert.Equal(t, 2, len(votes))
		assert.Equal(t, msgPrevote, votes[0].Type())
		assert.Equal(t, msgPrevote, votes[1].Type())
		assert.Equal(t, height, votes[0].H())
		assert.Equal(t, round, votes[0].R())
		assert.Equal(t, height, votes[1].H())
		assert.Equal(t, round, votes[1].R())
	})

	t.Run("delete msgs at a specific height", func(t *testing.T) {
		ms := newMsgStore()
		preVoteNil := newVoteMsg(height, round, msgPrevote, addrAlice, nilValue)
		_, err := ms.Save(preVoteNil)
		if err != nil {
			assert.Error(t, err)
		}

		preVoteNoneNil := newVoteMsg(height, round, msgPrevote, addrBob, noneNilValue)
		_, err = ms.Save(preVoteNoneNil)
		if err != nil {
			assert.Error(t, err)
		}

		ms.DeleteMsgsAtHeight(height)

		votes := ms.Get(height, func(m *core.Message) bool {
			return m.H() == height
		})

		assert.Equal(t, 0, len(votes))
	})

}
