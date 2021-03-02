package afd

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func newVoteMsg(h uint64, r int64, code uint64, addr common.Address, value common.Hash) *core.Message {
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

func TestMsgStore_Get(t *testing.T) {
	height := uint64(100)
	round := int64(0)
	nodeAddr := common.Address{0x1}

	t.Run("msg store is empty", func(t *testing.T) {
		ms := newMsgStore()
		proposals := ms.Get(height, func(m *core.Message) bool {
			return m.Type() == msgProposal
		})
		assert.Equal(t, 0, len(proposals))
	})

	t.Run("query preVote for nil from msg store", func(t *testing.T) {
		ms := newMsgStore()
		preVote := newVoteMsg(height, round, msgPrevote, nodeAddr, nilValue)
		_, err := ms.Save(preVote)
		if err != nil {
			assert.Error(t, err)
		}

		votes := ms.Get(height, func(m *core.Message) bool {
			return m.Type() == msgPrevote && m.H() == height && m.R() == round && m.Sender() == nodeAddr &&
				m.Value() == nilValue
		})

		assert.Equal(t, 1, len(votes))
	})

}

func TestMsgStore_Save(t *testing.T) {

}

func TestMsgStore_DeleteMsgsAtHeight(t *testing.T) {

}
