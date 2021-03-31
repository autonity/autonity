package faultdetector

import (
	"crypto/ecdsa"
	"github.com/clearmatics/autonity/common"
	tdm "github.com/clearmatics/autonity/consensus/tendermint"
	algo "github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"math/rand"
	"sort"
	"testing"
)

type addressKeyMap map[common.Address]*ecdsa.PrivateKey

func generateCommittee(n int) (types.Committee, addressKeyMap) { // nolint: unparam
	vals := make(types.Committee, 0)
	keymap := make(addressKeyMap)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		committeeMember := types.CommitteeMember{
			Address:     crypto.PubkeyToAddress(privateKey.PublicKey),
			VotingPower: new(big.Int).SetUint64(1),
		}
		vals = append(vals, committeeMember)
		keymap[committeeMember.Address] = privateKey
	}
	sort.Sort(vals)
	return vals, keymap
}

func newBlockHeader(height uint64, committee types.Committee) *types.Header {
	// use random nonce to create different blocks
	var nonce types.BlockNonce
	for i := 0; i < len(nonce); i++ {
		nonce[0] = byte(rand.Intn(256))
	}
	return &types.Header{
		Number:    new(big.Int).SetUint64(height),
		Nonce:     nonce,
		Committee: committee,
	}
}

// new proposal with meta data, if the withValue is not nil, it will use the value as proposal, otherwise an
// random block will be used as the value for proposal.
func newProposalMessage(h uint64, r int64, vr int64, senderKey *ecdsa.PrivateKey, committee types.Committee, withValue *types.Block) *tdm.Message {
	header := newBlockHeader(h, committee)
	block := withValue
	if withValue == nil {
		block = types.NewBlockWithHeader(header)
	}

	p := &algo.ConsensusMessage{
		MsgType:    algo.Propose,
		Height:     h,
		Round:      r,
		Value:      algo.ValueID(block.Hash()),
		ValidRound: vr,
	}

	msgBytes, err := tdm.EncodeSignedMessage(p, senderKey, block)
	if err != nil {
		return nil
	}

	msg, err := tdm.DecodeSignedMessage(msgBytes)
	if err != nil {
		return nil
	}
	return msg
}

func newVoteMsg(h uint64, r int64, msgType algo.Step, senderKey *ecdsa.PrivateKey, value common.Hash) *tdm.Message {
	vote := &algo.ConsensusMessage{
		MsgType: msgType,
		Height:  h,
		Round:   r,
		Value:   algo.ValueID(value),
	}
	msgBytes, err := tdm.EncodeSignedMessage(vote, senderKey, nil)
	if err != nil {
		return nil
	}

	msg, err := tdm.DecodeSignedMessage(msgBytes)
	if err != nil {
		return nil
	}
	return msg
}

func TestMsgStore(t *testing.T) {
	height := uint64(100)
	round := int64(0)

	committee, keys := generateCommittee(5)
	proposer := committee[0].Address
	proposerKey := keys[proposer]

	addrAlice := committee[0].Address
	addrBob := committee[1].Address
	keyBob := keys[addrBob]
	noneNilValue := common.Hash{0x1}

	t.Run("query msg store when msg store is empty", func(t *testing.T) {
		ms := newMsgStore()
		proposals := ms.Get(height, func(m *tdm.Message) bool {
			return m.Type() == algo.Propose
		})
		assert.Equal(t, 0, len(proposals))
	})

	t.Run("save equivocation msgs in msg store", func(t *testing.T) {
		ms := newMsgStore()
		preVoteNil := newVoteMsg(height, round, algo.Prevote, proposerKey, nilValue)
		_, err := ms.Save(preVoteNil)
		if err != nil {
			assert.Error(t, err)
		}

		preVoteNoneNil := newVoteMsg(height, round, algo.Prevote, proposerKey, noneNilValue)
		equivocatedMsgs, err := ms.Save(preVoteNoneNil)
		assert.NotNil(t, equivocatedMsgs)
		assert.Equal(t, err, errEquivocation)
		assert.Equal(t, nilValue, equivocatedMsgs[0].V())
		assert.Equal(t, addrAlice, equivocatedMsgs[0].Sender())
		assert.Equal(t, height, equivocatedMsgs[0].H())
		assert.Equal(t, round, equivocatedMsgs[0].R())
		assert.Equal(t, algo.Prevote, equivocatedMsgs[0].Type())
		// check equivocated msg is also stored at msg store.
		votes := ms.Get(height, func(m *tdm.Message) bool {
			return m.Type() == algo.Prevote && m.H() == height && m.R() == round && m.Sender() == addrAlice
		})
		assert.Equal(t, 2, len(votes))
	})

	t.Run("query a presented preVote from msg store", func(t *testing.T) {
		ms := newMsgStore()
		preVote := newVoteMsg(height, round, algo.Prevote, proposerKey, nilValue)
		_, err := ms.Save(preVote)
		if err != nil {
			assert.Error(t, err)
		}

		votes := ms.Get(height, func(m *tdm.Message) bool {
			return m.Type() == algo.Prevote && m.H() == height && m.R() == round && m.Sender() == addrAlice &&
				m.V() == nilValue
		})

		assert.Equal(t, 1, len(votes))
		assert.Equal(t, algo.Prevote, votes[0].Type())
		assert.Equal(t, height, votes[0].H())
		assert.Equal(t, round, votes[0].R())
		assert.Equal(t, addrAlice, votes[0].Sender())
		assert.Equal(t, nilValue, votes[0].V())
	})

	t.Run("query multiple presented preVote from msg store", func(t *testing.T) {
		ms := newMsgStore()
		preVoteNil := newVoteMsg(height, round, algo.Prevote, proposerKey, nilValue)
		_, err := ms.Save(preVoteNil)
		if err != nil {
			assert.Error(t, err)
		}

		preVoteNoneNil := newVoteMsg(height, round, algo.Prevote, keyBob, noneNilValue)
		_, err = ms.Save(preVoteNoneNil)
		if err != nil {
			assert.Error(t, err)
		}

		votes := ms.Get(height, func(m *tdm.Message) bool {
			return m.Type() == algo.Prevote && m.H() == height && m.R() == round
		})

		assert.Equal(t, 2, len(votes))
		assert.Equal(t, algo.Prevote, votes[0].Type())
		assert.Equal(t, algo.Prevote, votes[1].Type())
		assert.Equal(t, height, votes[0].H())
		assert.Equal(t, round, votes[0].R())
		assert.Equal(t, height, votes[1].H())
		assert.Equal(t, round, votes[1].R())
	})

	t.Run("delete msgs at a specific height", func(t *testing.T) {
		ms := newMsgStore()
		preVoteNil := newVoteMsg(height, round, algo.Prevote, proposerKey, nilValue)
		_, err := ms.Save(preVoteNil)
		if err != nil {
			assert.Error(t, err)
		}

		preVoteNoneNil := newVoteMsg(height, round, algo.Prevote, keyBob, noneNilValue)
		_, err = ms.Save(preVoteNoneNil)
		if err != nil {
			assert.Error(t, err)
		}

		ms.DeleteMsgsAtHeight(height)

		votes := ms.Get(height, func(m *tdm.Message) bool {
			return m.H() == height
		})

		assert.Equal(t, 0, len(votes))
	})

}
