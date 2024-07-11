package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
)

func TestMsgStore(t *testing.T) {
	height := uint64(100)
	round := int64(0)

	cSize := 5
	proposerIdx := 0
	committee, keys := GenerateCommittee(cSize)
	proposer := committee.Members[proposerIdx].Address
	proposerKey := keys[proposer].consensus

	indexBob := 1
	addrBob := committee.Members[indexBob].Address
	keyBob := keys[addrBob].consensus
	notNilValue := common.Hash{0x1}

	t.Run("query msg store when msg store is empty", func(t *testing.T) {
		ms := NewMsgStore()
		proposals := ms.GetProposals(height, func(_ *message.Propose) bool {
			return true
		})
		assert.Equal(t, 0, len(proposals))
	})

	t.Run("save equivocation msgs in msg store", func(t *testing.T) {
		ms := NewMsgStore()
		preVoteNil := message.NewPrevote(round, height, NilValue, makeSigner(proposerKey), &committee.Members[proposerIdx], cSize)
		ms.Save(preVoteNil)

		preVoteNoneNil := message.NewPrevote(round, height, notNilValue, makeSigner(proposerKey), &committee.Members[proposerIdx], cSize)
		ms.Save(preVoteNoneNil)
		// check equivocated msg is also stored at msg store.
		votes := ms.GetPrevotes(height, func(m *message.Prevote) bool {
			return m.R() == round && m.Signers().Contains(proposerIdx)
		})
		assert.Equal(t, 2, len(votes))
		assert.Equal(t, 1, votes[0].Signers().Len())
		require.Equal(t, common.Big1, ms.PrevotesPowerFor(height, round, NilValue))
		require.Equal(t, common.Big1, ms.PrevotesPowerFor(height, round, notNilValue))
	})

	t.Run("Save aggregated votes in msg store", func(t *testing.T) {
		ms := NewMsgStore()
		var prevotes []message.Vote
		for _, member := range committee.Members {
			m := member
			preVoteNil := message.NewPrevote(round, height, NilValue, makeSigner(keys[member.Address].consensus), &m, cSize)
			prevotes = append(prevotes, preVoteNil)
		}

		aggVote := message.AggregatePrevotes(prevotes)
		ms.Save(aggVote)

		// for every account, they have the prevote saved.
		for i, member := range committee.Members {
			m := member
			votes := ms.GetPrevotes(height, func(msg *message.Prevote) bool {
				return msg.R() == round && msg.Value() == NilValue && msg.Signers().Contains(int(m.Index))
			})
			require.Equal(t, 1, len(votes))
			require.Equal(t, true, votes[0].Signers().Contains(i))
		}

		// query for the target aggregated prevote, only 1 prevote is returned.
		votes := ms.GetPrevotes(height, func(m *message.Prevote) bool {
			return m.R() == round && m.Value() == NilValue
		})

		require.Equal(t, 1, len(votes))
		require.Equal(t, common.Big5, ms.PrevotesPowerFor(height, round, NilValue))
	})

	t.Run("query a presented preVote from msg store", func(t *testing.T) {
		ms := NewMsgStore()
		preVote := message.NewPrevote(round, height, NilValue, makeSigner(proposerKey), &committee.Members[proposerIdx], cSize)
		ms.Save(preVote)

		votes := ms.GetPrevotes(height, func(m *message.Prevote) bool {
			return m.R() == round && m.Value() == NilValue && m.Signers().Contains(proposerIdx)
		})

		assert.Equal(t, 1, len(votes))
		assert.Equal(t, message.PrevoteCode, votes[0].Code())
		assert.Equal(t, height, votes[0].H())
		assert.Equal(t, round, votes[0].R())
		assert.Equal(t, 1, votes[0].Signers().Len())
		assert.Equal(t, true, votes[0].Signers().Contains(proposerIdx))
		assert.Equal(t, NilValue, votes[0].Value())
	})

	t.Run("query multiple presented preVote from msg store", func(t *testing.T) {
		ms := NewMsgStore()
		preVoteNil := message.NewPrevote(round, height, NilValue, makeSigner(proposerKey), &committee.Members[proposerIdx], cSize)
		ms.Save(preVoteNil)

		preVoteNoneNil := message.NewPrevote(round, height, notNilValue, makeSigner(keyBob), &committee.Members[1], cSize)
		ms.Save(preVoteNoneNil)

		votes := ms.GetPrevotes(height, func(m *message.Prevote) bool {
			return m.R() == round
		})

		assert.Equal(t, 2, len(votes))
		assert.Equal(t, message.PrevoteCode, votes[0].Code())
		assert.Equal(t, message.PrevoteCode, votes[1].Code())
		assert.Equal(t, height, votes[0].H())
		assert.Equal(t, round, votes[0].R())
		assert.Equal(t, height, votes[1].H())
		assert.Equal(t, round, votes[1].R())
		assert.Equal(t, 1, votes[0].Signers().Len())
		assert.Equal(t, 1, votes[1].Signers().Len())
	})

	t.Run("delete msgs at a specific height", func(t *testing.T) {
		ms := NewMsgStore()
		preVoteNil := message.NewPrevote(round, height, NilValue, makeSigner(proposerKey), &committee.Members[proposerIdx], cSize)
		ms.Save(preVoteNil)
		preVoteNoneNil := message.NewPrevote(round, height, notNilValue, makeSigner(keyBob), &committee.Members[1], cSize)
		ms.Save(preVoteNoneNil)
		ms.DeleteOlds(height)
		prevotes := ms.GetPrevotes(height, func(m *message.Prevote) bool {
			return true
		})
		assert.Equal(t, 0, len(prevotes))
		require.Equal(t, uint64(0), ms.PrevotesPowerFor(height, round, NilValue).Uint64())
	})

	t.Run("get equivocated votes", func(t *testing.T) {
		ms := NewMsgStore()
		preVoteNil := message.NewPrevote(round, height, NilValue, makeSigner(proposerKey), &committee.Members[proposerIdx], cSize)
		ms.Save(preVoteNil)

		preVoteNoneNil := message.NewPrevote(round, height, notNilValue, makeSigner(proposerKey), &committee.Members[proposerIdx], cSize)
		ms.Save(preVoteNoneNil)

		v := common.Hash{0x23}
		votes := ms.GetPrevotes(height, func(m *message.Prevote) bool {
			return m.R() == round && m.Signers().Contains(proposerIdx) && m.Value() != v
		})
		assert.Equal(t, 2, len(votes))
		assert.Equal(t, height, votes[0].H())
		assert.Equal(t, height, votes[1].H())
		assert.Equal(t, round, votes[0].R())
		assert.Equal(t, round, votes[1].R())
		assert.Equal(t, message.PrevoteCode, votes[0].Code())
		assert.Equal(t, message.PrevoteCode, votes[1].Code())
		assert.NotEqual(t, v, votes[0].Value())
		assert.NotEqual(t, v, votes[1].Value())
		assert.Equal(t, 1, votes[0].Signers().Len())
		assert.Equal(t, 1, votes[1].Signers().Len())
	})
	t.Run("SearchQuorum correctly detects quorum of prevotes", func(t *testing.T) {
		ms := NewMsgStore()
		preVoteNil := message.NewPrevote(round, height, NilValue, makeSigner(proposerKey), &committee.Members[proposerIdx], cSize)
		ms.Save(preVoteNil)

		require.Equal(t, 0, len(ms.SearchQuorum(height, round, NilValue, common.Big1)))

		preVoteNotNil := message.NewPrevote(round, height, notNilValue, makeSigner(proposerKey), &committee.Members[proposerIdx], cSize)
		ms.Save(preVoteNotNil)

		require.Equal(t, 1, len(ms.SearchQuorum(height, round, NilValue, common.Big1)))
		require.Equal(t, preVoteNotNil.Hash(), ms.SearchQuorum(height, round, NilValue, common.Big1)[0].Hash())

		preVoteNotNil = message.NewPrevote(round, height, notNilValue, makeSigner(keyBob), &committee.Members[indexBob], cSize)
		ms.Save(preVoteNotNil)

		require.Equal(t, 2, len(ms.SearchQuorum(height, round, NilValue, common.Big1)))
		require.Equal(t, 2, len(ms.SearchQuorum(height, round, NilValue, common.Big2)))
		require.Equal(t, 0, len(ms.SearchQuorum(height, round, NilValue, common.Big3)))

		require.Equal(t, 1, len(ms.SearchQuorum(height, round, notNilValue, common.Big1)))
		require.Equal(t, 0, len(ms.SearchQuorum(height, round, notNilValue, common.Big2)))
	})
}
