package afd

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func newVoteMsg(h uint64, r int64, code uint64, addr common.Address, v common.Hash) *types.ConsensusMessage {
	var vote = types.Vote{
		Round:             r,
		Height:            new(big.Int).SetUint64(h),
		ProposedBlockHash: v,
	}

	encodedVote, err := types.Encode(&vote)
	if err != nil {
		return nil
	}

	var msg = types.ConsensusMessage{
		Code:          code,
		Msg:           encodedVote,
		Address:       addr,
		CommittedSeal: []byte{},
		Power:         0,
	}

	payload, err := msg.PayloadNoSig()
	if err != nil {
		return nil
	}

	m := new(types.ConsensusMessage)
	if err := msg.FromPayload(payload); err != nil {
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
		proposals := ms.Get(height, func(m *types.ConsensusMessage) bool {
			return m.Type() == types.MsgProposal
		})
		assert.Equal(t, 0, len(proposals))
	})

	t.Run("query preVote for nil from msg store", func(t *testing.T) {
		ms := newMsgStore()
		preVote := newVoteMsg(height, round, types.MsgPrevote, nodeAddr, nilValue)
		_, err := ms.Save(preVote)
		if err != nil {
			assert.Error(t, err)
		}

		votes := ms.Get(height, func(m *types.ConsensusMessage) bool {
			return m.Type() == types.MsgPrevote && m.H() == height && m.R() == round && m.Sender() == nodeAddr &&
				m.Value() == nilValue
		})

		assert.Equal(t, 1, len(votes))
	})

}

func TestMsgStore_Save(t *testing.T) {

}

func TestMsgStore_DeleteMsgsAtHeight(t *testing.T) {

}