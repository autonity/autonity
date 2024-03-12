package core

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
)

func TestMsgStore(t *testing.T) {
	height := uint64(100)
	round := int64(0)

	committee, keys := GenerateCommittee(5)
	proposer := committee[0].Address
	proposerKey := keys[proposer].consensus

	addrAlice := committee[0].Address
	addrBob := committee[1].Address
	keyBob := keys[addrBob].consensus
	notNilValue := common.Hash{0x1}

	t.Run("query msg store when msg store is empty", func(t *testing.T) {
		ms := NewMsgStore()
		proposals := ms.Get(height, func(m message.Msg) bool {
			return m.Code() == message.ProposalCode
		})
		assert.Equal(t, 0, len(proposals))
	})

	t.Run("save equivocation msgs in msg store", func(t *testing.T) {
		ms := NewMsgStore()
		preVoteNil := message.NewPrevote(round, height, NilValue, makeSigner(proposerKey, proposer)).MustVerify(stubVerifier(proposerKey.PublicKey()))
		ms.Save(preVoteNil)

		preVoteNoneNil := message.NewPrevote(round, height, notNilValue, makeSigner(proposerKey, proposer)).MustVerify(stubVerifier(proposerKey.PublicKey()))
		ms.Save(preVoteNoneNil)
		// check equivocated msg is also stored at msg store.
		votes := ms.Get(height, func(m message.Msg) bool {
			return m.Code() == message.PrevoteCode && m.H() == height && m.R() == round && m.Sender() == addrAlice
		})
		assert.Equal(t, 2, len(votes))
	})

	t.Run("query a presented preVote from msg store", func(t *testing.T) {
		ms := NewMsgStore()
		preVote := message.NewPrevote(round, height, NilValue, makeSigner(proposerKey, proposer)).MustVerify(stubVerifier(proposerKey.PublicKey()))
		ms.Save(preVote)

		votes := ms.Get(height, func(m message.Msg) bool {
			return m.Code() == message.PrevoteCode && m.H() == height && m.R() == round && m.Sender() == addrAlice &&
				m.Value() == NilValue
		})

		assert.Equal(t, 1, len(votes))
		assert.Equal(t, message.PrevoteCode, votes[0].Code())
		assert.Equal(t, height, votes[0].H())
		assert.Equal(t, round, votes[0].R())
		assert.Equal(t, addrAlice, votes[0].Sender())
		assert.Equal(t, NilValue, votes[0].Value())
	})

	t.Run("query multiple presented preVote from msg store", func(t *testing.T) {
		ms := NewMsgStore()
		preVoteNil := message.NewPrevote(round, height, NilValue, makeSigner(proposerKey, proposer)).MustVerify(stubVerifier(proposerKey.PublicKey()))
		ms.Save(preVoteNil)

		preVoteNoneNil := message.NewPrevote(round, height, notNilValue, makeSigner(keyBob, addrBob)).MustVerify(stubVerifier(keyBob.PublicKey()))
		ms.Save(preVoteNoneNil)

		votes := ms.Get(height, func(m message.Msg) bool {
			return m.Code() == message.PrevoteCode && m.H() == height && m.R() == round
		})

		assert.Equal(t, 2, len(votes))
		assert.Equal(t, message.PrevoteCode, votes[0].Code())
		assert.Equal(t, message.PrevoteCode, votes[1].Code())
		assert.Equal(t, height, votes[0].H())
		assert.Equal(t, round, votes[0].R())
		assert.Equal(t, height, votes[1].H())
		assert.Equal(t, round, votes[1].R())
	})

	t.Run("delete msgs at a specific height", func(t *testing.T) {
		ms := NewMsgStore()
		preVoteNil := message.NewPrevote(round, height, NilValue, makeSigner(proposerKey, proposer)).MustVerify(stubVerifier(proposerKey.PublicKey()))
		ms.Save(preVoteNil)
		preVoteNoneNil := message.NewPrevote(round, height, notNilValue, makeSigner(keyBob, addrBob)).MustVerify(stubVerifier(keyBob.PublicKey()))
		ms.Save(preVoteNoneNil)
		ms.DeleteOlds(height)
		votes := ms.Get(height, func(m message.Msg) bool {
			return m.H() == height
		})
		assert.Equal(t, 0, len(votes))
	})

}

func stubVerifier(consensusKey blst.PublicKey) func(address common.Address) *types.CommitteeMember {
	return func(address common.Address) *types.CommitteeMember {
		return &types.CommitteeMember{
			Address:      address,
			VotingPower:  common.Big1,
			ConsensusKey: consensusKey,
		}
	}
}

func stubVerifierWithPower(consensusKey blst.PublicKey, power int64) func(address common.Address) *types.CommitteeMember {
	return func(address common.Address) *types.CommitteeMember {
		return &types.CommitteeMember{
			Address:      address,
			VotingPower:  big.NewInt(power),
			ConsensusKey: consensusKey,
		}
	}
}
