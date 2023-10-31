package core

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"math/rand"
	"testing"
)

func newBlockHeader(height uint64, committee types.Committee) *types.Header {
	// use random nonce to create different blocks
	var nonce types.BlockNonce
	for i := 0; i < len(nonce); i++ {
		nonce[0] = byte(rand.Intn(256)) //nolint
	}
	return &types.Header{
		Number:    new(big.Int).SetUint64(height),
		Nonce:     nonce,
		Committee: committee,
	}
}

func newVoteMsg(h uint64, r int64, code uint8, senderKey *ecdsa.PrivateKey, value common.Hash, committee types.Committee) *message.Message { //nolint
	lastHeader := newBlockHeader(h-1, committee)
	var vote = message.Vote{
		Round:             r,
		Height:            new(big.Int).SetUint64(h),
		ProposedBlockHash: value,
	}

	encodedVote, err := rlp.EncodeToBytes(&vote)
	if err != nil {
		return nil
	}

	msg := createMsg(encodedVote, code, senderKey)

	return decodeMsg(msg, lastHeader)
}

func CheckValidatorSignature(previousHeader *types.Header, data []byte, sig []byte) (common.Address, error) {
	// 1. Get signature address
	signer, err := types.GetSignatureAddress(data, sig)
	if err != nil {
		return common.Address{}, err
	}

	// 2. Check validator
	val := previousHeader.CommitteeMember(signer)
	if val == nil {
		return common.Address{}, fmt.Errorf("wrong membership")
	}

	return val.Address, nil
}

func createMsg(rlpBytes []byte, code uint8, senderKey *ecdsa.PrivateKey) *message.Message {
	var msg = message.Message{
		Code:          code,
		Payload:       rlpBytes,
		Address:       crypto.PubkeyToAddress(senderKey.PublicKey),
		CommittedSeal: []byte{},
	}
	data, err := msg.BytesNoSignature()
	if err != nil {
		return nil
	}

	hashData := crypto.Keccak256(data)
	msg.Signature, err = crypto.Sign(hashData, senderKey)
	if err != nil {
		return nil
	}
	return &msg
}

// decode msg do the msg decoding and validation to recover the voting power and decodedMsg fields.
func decodeMsg(msg *message.Message, lastHeader *types.Header) *message.Message {
	m, err := message.FromBytes(msg.GetBytes())
	if err != nil {
		return nil
	}
	// validate msg and get voting power with last header.
	if err = m.Validate(CheckValidatorSignature, lastHeader); err != nil {
		return nil
	}
	return m
}

func TestMsgStore(t *testing.T) {
	height := uint64(100)
	round := int64(0)

	committee, keys := tendermint.GenerateCommittee(5)
	proposer := committee[0].Address
	proposerKey := keys[proposer]

	addrAlice := committee[0].Address
	addrBob := committee[1].Address
	keyBob := keys[addrBob]
	noneNilValue := common.Hash{0x1}

	t.Run("query msg store when msg store is empty", func(t *testing.T) {
		ms := NewMsgStore()
		proposals := ms.Get(height, func(m *message.Message) bool {
			return m.Type() == consensus.MsgProposal
		})
		assert.Equal(t, 0, len(proposals))
	})

	t.Run("save equivocation msgs in msg store", func(t *testing.T) {
		ms := NewMsgStore()
		preVoteNil := newVoteMsg(height, round, consensus.MsgPrevote, proposerKey, NilValue, committee)
		ms.Save(preVoteNil)

		preVoteNoneNil := newVoteMsg(height, round, consensus.MsgPrevote, proposerKey, noneNilValue, committee)
		ms.Save(preVoteNoneNil)
		// check equivocated msg is also stored at msg store.
		votes := ms.Get(height, func(m *message.Message) bool {
			return m.Type() == consensus.MsgPrevote && m.H() == height && m.R() == round && m.Sender() == addrAlice
		})
		assert.Equal(t, 2, len(votes))
	})

	t.Run("query a presented preVote from msg store", func(t *testing.T) {
		ms := NewMsgStore()
		preVote := newVoteMsg(height, round, consensus.MsgPrevote, proposerKey, NilValue, committee)
		ms.Save(preVote)

		votes := ms.Get(height, func(m *message.Message) bool {
			return m.Type() == consensus.MsgPrevote && m.H() == height && m.R() == round && m.Sender() == addrAlice &&
				m.Value() == NilValue
		})

		assert.Equal(t, 1, len(votes))
		assert.Equal(t, consensus.MsgPrevote, votes[0].Type())
		assert.Equal(t, height, votes[0].H())
		assert.Equal(t, round, votes[0].R())
		assert.Equal(t, addrAlice, votes[0].Sender())
		assert.Equal(t, NilValue, votes[0].Value())
	})

	t.Run("query multiple presented preVote from msg store", func(t *testing.T) {
		ms := NewMsgStore()
		preVoteNil := newVoteMsg(height, round, consensus.MsgPrevote, proposerKey, NilValue, committee)
		ms.Save(preVoteNil)

		preVoteNoneNil := newVoteMsg(height, round, consensus.MsgPrevote, keyBob, noneNilValue, committee)
		ms.Save(preVoteNoneNil)

		votes := ms.Get(height, func(m *message.Message) bool {
			return m.Type() == consensus.MsgPrevote && m.H() == height && m.R() == round
		})

		assert.Equal(t, 2, len(votes))
		assert.Equal(t, consensus.MsgPrevote, votes[0].Type())
		assert.Equal(t, consensus.MsgPrevote, votes[1].Type())
		assert.Equal(t, height, votes[0].H())
		assert.Equal(t, round, votes[0].R())
		assert.Equal(t, height, votes[1].H())
		assert.Equal(t, round, votes[1].R())
	})

	t.Run("delete msgs at a specific height", func(t *testing.T) {
		ms := NewMsgStore()
		preVoteNil := newVoteMsg(height, round, consensus.MsgPrevote, proposerKey, NilValue, committee)
		ms.Save(preVoteNil)

		preVoteNoneNil := newVoteMsg(height, round, consensus.MsgPrevote, keyBob, noneNilValue, committee)
		ms.Save(preVoteNoneNil)

		ms.DeleteOlds(height)

		votes := ms.Get(height, func(m *message.Message) bool {
			return m.H() == height
		})

		assert.Equal(t, 0, len(votes))
	})

}
