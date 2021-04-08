package faultdetector

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/core"
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
func newProposalMessage(h uint64, r int64, vr int64, senderKey *ecdsa.PrivateKey, committee types.Committee, withValue *types.Block) *core.Message {
	header := newBlockHeader(h, committee)
	lastHeader := newBlockHeader(h-1, committee)
	block := withValue
	if withValue == nil {
		block = types.NewBlockWithHeader(header)
	}

	proposal := &core.Proposal{
		Round:         r,
		Height:        new(big.Int).SetUint64(h),
		ValidRound:    vr,
		ProposalBlock: block,
	}
	encodeProposal, err := core.Encode(proposal)
	if err != nil {
		return nil
	}

	msg := createMsg(encodeProposal, msgProposal, senderKey)

	return decodeMsg(msg, lastHeader)
}

func newVoteMsg(h uint64, r int64, code uint8, senderKey *ecdsa.PrivateKey, value common.Hash, committee types.Committee) *core.Message {
	lastHeader := newBlockHeader(h-1, committee)
	var vote = core.Vote{
		Round:             r,
		Height:            new(big.Int).SetUint64(h),
		ProposedBlockHash: value,
	}

	encodedVote, err := core.Encode(&vote)
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

func createMsg(rlpBytes []byte, code uint8, senderKey *ecdsa.PrivateKey) *core.Message {
	var msg = core.Message{
		Code:          uint64(code),
		Msg:           rlpBytes,
		Address:       crypto.PubkeyToAddress(senderKey.PublicKey),
		CommittedSeal: []byte{},
	}
	data, err := msg.PayloadNoSig()
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
func decodeMsg(msg *core.Message, lastHeader *types.Header) *core.Message {
	m := new(core.Message)
	err := m.FromPayload(msg.Payload())
	if err != nil {
		return nil
	}

	// validate msg and get voting power with last header.
	if _, err = m.Validate(CheckValidatorSignature, lastHeader); err != nil {
		return nil
	}
	return m
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
		proposals := ms.Get(height, func(m *core.Message) bool {
			return uint8(m.Type()) == msgProposal
		})
		assert.Equal(t, 0, len(proposals))
	})

	t.Run("save equivocation msgs in msg store", func(t *testing.T) {
		ms := newMsgStore()
		preVoteNil := newVoteMsg(height, round, msgPrevote, proposerKey, nilValue, committee)
		_, err := ms.Save(preVoteNil)
		if err != nil {
			assert.Error(t, err)
		}

		preVoteNoneNil := newVoteMsg(height, round, msgPrevote, proposerKey, noneNilValue, committee)
		equivocatedMsgs, err := ms.Save(preVoteNoneNil)
		assert.NotNil(t, equivocatedMsgs)
		assert.Equal(t, err, errEquivocation)
		assert.Equal(t, nilValue, equivocatedMsgs[0].Value())
		assert.Equal(t, addrAlice, equivocatedMsgs[0].Sender())
		assert.Equal(t, height, equivocatedMsgs[0].H())
		assert.Equal(t, round, equivocatedMsgs[0].R())
		assert.Equal(t, msgPrevote, uint8(equivocatedMsgs[0].Type()))
		// check equivocated msg is also stored at msg store.
		votes := ms.Get(height, func(m *core.Message) bool {
			return uint8(m.Type()) == msgPrevote && m.H() == height && m.R() == round && m.Sender() == addrAlice
		})
		assert.Equal(t, 2, len(votes))
	})

	t.Run("query a presented preVote from msg store", func(t *testing.T) {
		ms := newMsgStore()
		preVote := newVoteMsg(height, round, msgPrevote, proposerKey, nilValue, committee)
		_, err := ms.Save(preVote)
		if err != nil {
			assert.Error(t, err)
		}

		votes := ms.Get(height, func(m *core.Message) bool {
			return uint8(m.Type()) == msgPrevote && m.H() == height && m.R() == round && m.Sender() == addrAlice &&
				m.Value() == nilValue
		})

		assert.Equal(t, 1, len(votes))
		assert.Equal(t, msgPrevote, uint8(votes[0].Type()))
		assert.Equal(t, height, votes[0].H())
		assert.Equal(t, round, votes[0].R())
		assert.Equal(t, addrAlice, votes[0].Sender())
		assert.Equal(t, nilValue, votes[0].Value())
	})

	t.Run("query multiple presented preVote from msg store", func(t *testing.T) {
		ms := newMsgStore()
		preVoteNil := newVoteMsg(height, round, msgPrevote, proposerKey, nilValue, committee)
		_, err := ms.Save(preVoteNil)
		if err != nil {
			assert.Error(t, err)
		}

		preVoteNoneNil := newVoteMsg(height, round, msgPrevote, keyBob, noneNilValue, committee)
		_, err = ms.Save(preVoteNoneNil)
		if err != nil {
			assert.Error(t, err)
		}

		votes := ms.Get(height, func(m *core.Message) bool {
			return uint8(m.Type()) == msgPrevote && m.H() == height && m.R() == round
		})

		assert.Equal(t, 2, len(votes))
		assert.Equal(t, msgPrevote, uint8(votes[0].Type()))
		assert.Equal(t, msgPrevote, uint8(votes[1].Type()))
		assert.Equal(t, height, votes[0].H())
		assert.Equal(t, round, votes[0].R())
		assert.Equal(t, height, votes[1].H())
		assert.Equal(t, round, votes[1].R())
	})

	t.Run("delete msgs at a specific height", func(t *testing.T) {
		ms := newMsgStore()
		preVoteNil := newVoteMsg(height, round, msgPrevote, proposerKey, nilValue, committee)
		_, err := ms.Save(preVoteNil)
		if err != nil {
			assert.Error(t, err)
		}

		preVoteNoneNil := newVoteMsg(height, round, msgPrevote, keyBob, noneNilValue, committee)
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
