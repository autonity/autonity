package core

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
)

func TestMsgStore(t *testing.T) {
	height := uint64(100)
	round := int64(0)

	cSize := 5
	proposerIdx := 0
	committee, keys := GenerateCommittee(cSize)
	proposer := committee[proposerIdx].Address
	proposerKey := keys[proposer].consensus

	addrAlice := committee[0].Address
	addrBob := committee[1].Address
	keyBob := keys[addrBob].consensus
	notNilValue := common.Hash{0x1}

	t.Run("query msg store when msg store is empty", func(t *testing.T) {
		ms := NewMsgStore()
		proposals := ms.Get(func(m message.Msg) bool {
			return m.Code() == message.ProposalCode
		}, height, nil)
		assert.Equal(t, 0, len(proposals))
	})

	t.Run("save equivocation msgs in msg store", func(t *testing.T) {
		ms := NewMsgStore()
		preVoteNil := message.NewPrevote(round, height, NilValue, makeSigner(proposerKey), &committee[proposerIdx], cSize)
		ms.Save(preVoteNil, committee)

		preVoteNoneNil := message.NewPrevote(round, height, notNilValue, makeSigner(proposerKey), &committee[proposerIdx], cSize)
		ms.Save(preVoteNoneNil, committee)
		// check equivocated msg is also stored at msg store.
		votes := ms.Get(func(m message.Msg) bool {
			return m.Code() == message.PrevoteCode && m.H() == height && m.R() == round
		}, height, &addrAlice)
		assert.Equal(t, 2, len(votes))
		assert.Equal(t, 1, votes[0].(*message.Prevote).Signers().Len())
	})

	t.Run("Save aggregated votes in msg store", func(t *testing.T) {
		ms := NewMsgStore()
		var prevotes []message.Vote
		for _, c := range committee {
			preVoteNil := message.NewPrevote(round, height, NilValue, makeSigner(keys[c.Address].consensus), &c, cSize)
			prevotes = append(prevotes, preVoteNil)
		}

		aggVote := message.AggregatePrevotes(prevotes)
		ms.Save(aggVote, committee)

		// for every account, they have the prevote saved.
		for i, c := range committee {
			votes := ms.Get(func(m message.Msg) bool {
				return m.Code() == message.PrevoteCode && m.H() == height && m.R() == round && m.Value() == NilValue
			}, height, &c.Address)
			require.Equal(t, 1, len(votes))
			require.Equal(t, true, votes[0].(*message.Prevote).Signers().Contains(i))
		}

		// query for the target aggregated prevote, only 1 prevote is returned.
		votes := ms.Get(func(m message.Msg) bool {
			return m.Code() == message.PrevoteCode && m.H() == height && m.R() == round && m.Value() == NilValue
		}, height, nil)

		require.Equal(t, 1, len(votes))
	})

	t.Run("query a presented preVote from msg store", func(t *testing.T) {
		ms := NewMsgStore()
		preVote := message.NewPrevote(round, height, NilValue, makeSigner(proposerKey), &committee[proposerIdx], cSize)
		ms.Save(preVote, committee)

		votes := ms.Get(func(m message.Msg) bool {
			return m.Code() == message.PrevoteCode && m.H() == height && m.R() == round && m.Value() == NilValue
		}, height, &addrAlice)

		assert.Equal(t, 1, len(votes))
		assert.Equal(t, message.PrevoteCode, votes[0].Code())
		assert.Equal(t, height, votes[0].H())
		assert.Equal(t, round, votes[0].R())
		assert.Equal(t, 1, votes[0].(*message.Prevote).Signers().Len())
		assert.Equal(t, true, votes[0].(message.Vote).Signers().Contains(proposerIdx))
		assert.Equal(t, NilValue, votes[0].Value())
	})

	t.Run("query multiple presented preVote from msg store", func(t *testing.T) {
		ms := NewMsgStore()
		preVoteNil := message.NewPrevote(round, height, NilValue, makeSigner(proposerKey), &committee[proposerIdx], cSize)
		ms.Save(preVoteNil, committee)

		preVoteNoneNil := message.NewPrevote(round, height, notNilValue, makeSigner(keyBob), &committee[1], cSize)
		ms.Save(preVoteNoneNil, committee)

		votes := ms.Get(func(m message.Msg) bool {
			return m.Code() == message.PrevoteCode && m.H() == height && m.R() == round
		}, height, nil)

		assert.Equal(t, 2, len(votes))
		assert.Equal(t, message.PrevoteCode, votes[0].Code())
		assert.Equal(t, message.PrevoteCode, votes[1].Code())
		assert.Equal(t, height, votes[0].H())
		assert.Equal(t, round, votes[0].R())
		assert.Equal(t, height, votes[1].H())
		assert.Equal(t, round, votes[1].R())
		assert.Equal(t, 1, votes[0].(*message.Prevote).Signers().Len())
		assert.Equal(t, 1, votes[1].(*message.Prevote).Signers().Len())
	})

	t.Run("delete msgs at a specific height", func(t *testing.T) {
		ms := NewMsgStore()
		preVoteNil := message.NewPrevote(round, height, NilValue, makeSigner(proposerKey), &committee[proposerIdx], cSize)
		ms.Save(preVoteNil, committee)
		preVoteNoneNil := message.NewPrevote(round, height, notNilValue, makeSigner(keyBob), &committee[1], cSize)
		ms.Save(preVoteNoneNil, committee)
		ms.DeleteOlds(height)
		votes := ms.Get(func(m message.Msg) bool {
			return m.H() == height
		}, height, nil)
		assert.Equal(t, 0, len(votes))
	})

	t.Run("get equivocated votes", func(t *testing.T) {
		ms := NewMsgStore()
		preVoteNil := message.NewPrevote(round, height, NilValue, makeSigner(proposerKey), &committee[proposerIdx], cSize)
		ms.Save(preVoteNil, committee)

		preVoteNoneNil := message.NewPrevote(round, height, notNilValue, makeSigner(proposerKey), &committee[proposerIdx], cSize)
		ms.Save(preVoteNoneNil, committee)

		v := common.Hash{0x23}
		votes := ms.GetEquivocatedVotes(height, round, message.PrevoteCode, proposer, v)
		assert.Equal(t, 2, len(votes))
		assert.Equal(t, height, votes[0].H())
		assert.Equal(t, height, votes[1].H())
		assert.Equal(t, round, votes[0].R())
		assert.Equal(t, round, votes[1].R())
		assert.Equal(t, message.PrevoteCode, votes[0].Code())
		assert.Equal(t, message.PrevoteCode, votes[1].Code())
		assert.NotEqual(t, v, votes[0].Value())
		assert.NotEqual(t, v, votes[1].Value())
		assert.Equal(t, 1, votes[0].(*message.Prevote).Signers().Len())
		assert.Equal(t, 1, votes[1].(*message.Prevote).Signers().Len())
	})
}
