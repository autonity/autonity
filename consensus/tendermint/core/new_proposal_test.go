package core

import (
	"context"
	"errors"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"math/rand"
	"sort"
	"testing"
	"time"
)

// The following tests aim to test lines 22 - 27 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestTendermintNewProposal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backendMock := NewMockBackend(ctrl)

	minSize, maxSize := 4, 100
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address

	backendMock.EXPECT().Address().Return(clientAddr)
	c := New(backendMock)
	c.setCommitteeSet(committeeSet)

	t.Run("receive invalid proposal for current round", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))

		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(propose)

		var invalidProposal Proposal
		// members[currentRound] means that the sender is the proposer for the current round
		// assume that the message is from a member of committee set and the signature is signing the contents, however,
		// the proposal block inside the message is invalid
		invalidMsg := generateBlockProposal(t, currentRound, currentHeight, -1, members[currentRound].Address, true)
		err := invalidMsg.Decode(&invalidProposal)
		assert.Nil(t, err)

		// prepare prevote nil
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := preparePrevote(t, currentRound, currentHeight, common.Hash{}, clientAddr)

		backendMock.EXPECT().VerifyProposal(*invalidProposal.ProposalBlock).Return(time.Duration(1), errors.New("invalid proposal"))
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, prevoteMsgRLPWithSig).Return(nil)

		err = c.handleCheckedMsg(context.Background(), invalidMsg, members[currentRound])
		assert.Error(t, err, "expected an error for invalid proposal")
		assert.Equal(t, prevote, c.step)
	})
	t.Run("receive proposal with validRound = -1 and client's lockedRound = -1", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		clientLockedRound := int64(-1)

		var proposal Proposal
		proposalMsg := generateBlockProposal(t, currentRound, currentHeight, clientLockedRound, members[currentRound].Address, false)
		err := proposalMsg.Decode(&proposal) // we have to do this because encoding and decoding changes some default values
		assert.Nil(t, err)

		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := preparePrevote(t, currentRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr)

		// if lockedRround = - 1 then lockedValue = nil
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(propose)
		c.lockedRound = -1
		c.lockedValue = nil

		backendMock.EXPECT().VerifyProposal(*proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, prevoteMsgRLPWithSig).Return(nil)

		err = c.handleCheckedMsg(context.Background(), proposalMsg, members[currentRound])
		assert.Nil(t, err)
		assert.Equal(t, prevote, c.step)
		assert.Nil(t, c.lockedValue)
		assert.Equal(t, clientLockedRound, c.lockedRound)
	})
	t.Run("receive proposal with validRound = -1 and client's lockedValue is same as proposal block", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		clientLockedRound := int64(0)

		var proposal Proposal
		proposalMsg := generateBlockProposal(t, currentRound, currentHeight, -1, members[currentRound].Address, false)
		// we have to do this because encoding and decoding changes some default values and thus same blocks are no longer equal
		err := proposalMsg.Decode(&proposal)
		assert.Nil(t, err)

		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := preparePrevote(t, currentRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr)

		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(propose)
		c.lockedRound = clientLockedRound
		c.lockedValue = proposal.ProposalBlock
		c.validRound = clientLockedRound
		c.validValue = proposal.ProposalBlock

		backendMock.EXPECT().VerifyProposal(*proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, prevoteMsgRLPWithSig).Return(nil)

		err = c.handleCheckedMsg(context.Background(), proposalMsg, members[currentRound])
		assert.Nil(t, err)
		assert.Equal(t, prevote, c.step)
		assert.Equal(t, proposal.ProposalBlock, c.lockedValue)
		assert.Equal(t, clientLockedRound, c.lockedRound)
		assert.Equal(t, proposal.ProposalBlock, c.validValue)
		assert.Equal(t, clientLockedRound, c.validRound)
	})
	t.Run("receive proposal with validRound = -1 and client's lockedValue is different from proposal block", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		clientLockedRound := int64(0)
		clientLockedValue := generateBlock(currentHeight)

		var proposal Proposal
		proposalMsg := generateBlockProposal(t, currentRound, currentHeight, -1, members[currentRound].Address, false)
		// we have to do this because encoding and decoding changes some default values and thus same blocks are no longer equal
		err := proposalMsg.Decode(&proposal)
		assert.Nil(t, err)

		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := preparePrevote(t, currentRound, currentHeight, common.Hash{}, clientAddr)

		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(propose)
		c.lockedRound = clientLockedRound
		c.lockedValue = clientLockedValue
		c.validRound = clientLockedRound
		c.validValue = clientLockedValue

		backendMock.EXPECT().VerifyProposal(*proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, prevoteMsgRLPWithSig).Return(nil)

		err = c.handleCheckedMsg(context.Background(), proposalMsg, members[currentRound])
		assert.Nil(t, err)
		assert.Equal(t, prevote, c.step)
		assert.Equal(t, clientLockedValue, c.lockedValue)
		assert.Equal(t, clientLockedRound, c.lockedRound)
		assert.Equal(t, clientLockedValue, c.validValue)
		assert.Equal(t, clientLockedRound, c.validRound)
	})
}

// The following tests are not specific to proposal messages but rather apply to all messages
func TestHandleMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key1, err := crypto.GenerateKey()
	assert.Nil(t, err)
	key2, err := crypto.GenerateKey()
	assert.Nil(t, err)

	key1PubAddr := crypto.PubkeyToAddress(key1.PublicKey)
	key2PubAddr := crypto.PubkeyToAddress(key2.PublicKey)

	committeeSet, err := committee.NewSet(types.Committee{types.CommitteeMember{
		Address:     key1PubAddr,
		VotingPower: big.NewInt(1),
	}}, key1PubAddr)
	assert.Nil(t, err)

	backendMock := NewMockBackend(ctrl)
	backendMock.EXPECT().Address().Return(key1PubAddr)
	core := New(backendMock)

	t.Run("message sender is not in the committee set", func(t *testing.T) {
		// Prepare message
		msg := &Message{Address: key2PubAddr, Code: uint64(rand.Intn(3)), Msg: []byte("random message1")}

		msgRlpNoSig, err := msg.PayloadNoSig()
		assert.Nil(t, err)

		msg.Signature, err = crypto.Sign(crypto.Keccak256(msgRlpNoSig), key2)
		assert.Nil(t, err)

		msgRlpWithSig, err := msg.Payload()
		assert.Nil(t, err)

		core.setCommitteeSet(committeeSet)
		err = core.handleMsg(context.Background(), msgRlpWithSig)

		assert.Error(t, err, "unauthorised sender, sender is not is committees set")
	})

	t.Run("message sender is not the message siger", func(t *testing.T) {
		msg := &Message{Address: key1PubAddr, Code: uint64(rand.Intn(3)), Msg: []byte("random message2")}

		msgRlpNoSig, err := msg.PayloadNoSig()
		assert.Nil(t, err)

		msg.Signature, err = crypto.Sign(crypto.Keccak256(msgRlpNoSig), key1)
		assert.Nil(t, err)

		msgRlpWithSig, err := msg.Payload()
		assert.Nil(t, err)

		core.setCommitteeSet(committeeSet)
		err = core.handleMsg(context.Background(), msgRlpWithSig)

		assert.Error(t, err, "unauthorised sender, sender is not the signer of the message")
	})

	t.Run("malicious sender sends incorrect signature", func(t *testing.T) {
		sig, err := crypto.Sign(crypto.Keccak256([]byte("random bytes")), key1)
		assert.Nil(t, err)

		msg := &Message{Address: key1PubAddr, Code: uint64(rand.Intn(3)), Msg: []byte("random message2"), Signature: sig}
		msgRlpWithSig, err := msg.Payload()

		core.setCommitteeSet(committeeSet)
		err = core.handleMsg(context.Background(), msgRlpWithSig)

		assert.Error(t, err, "malicious sender sends different signature to signature of message")
	})
}

func prepareCommittee(t *testing.T, cSize int) *committee.Set {
	c := types.Committee{}
	for i := 1; i <= cSize; i++ {
		key, err := crypto.GenerateKey()
		assert.Nil(t, err)
		member := types.CommitteeMember{
			Address:     crypto.PubkeyToAddress(key.PublicKey),
			VotingPower: new(big.Int).SetInt64(1),
		}
		c = append(c, member)
	}

	sort.Sort(c)
	committeeSet, err := committee.NewSet(c, c[len(c)-1].Address)
	assert.Nil(t, err)
	return committeeSet
}

func generateBlock(height *big.Int) *types.Block {
	// use random nonce to create different blocks
	var nonce types.BlockNonce
	for i := 0; i < len(nonce); i++ {
		nonce[0] = byte(rand.Intn(256))
	}
	header := &types.Header{Number: height, Nonce: nonce}
	block := types.NewBlock(header, nil, nil, nil)
	return block
}

func generateBlockProposal(t *testing.T, r int64, h *big.Int, vr int64, src common.Address, invalid bool) *Message {
	var block *types.Block
	if invalid {
		header := &types.Header{Number: h}
		header.Difficulty = nil
		block = types.NewBlock(header, nil, nil, nil)
	} else {
		block = generateBlock(h)
	}
	proposal := NewProposal(r, h, vr, block)
	proposalRlp, err := Encode(proposal)
	assert.Nil(t, err)
	return &Message{
		Code:    msgProposal,
		Msg:     proposalRlp,
		Address: src,
	}
}

func preparePrevote(t *testing.T, round int64, height *big.Int, blockHash common.Hash, clientAddr common.Address) (*Message, []byte, []byte) {
	// prepare the proposal message
	voteRLP, err := Encode(&Vote{Round: round, Height: height, ProposedBlockHash: blockHash})
	assert.Nil(t, err)
	prevoteMsg := &Message{Code: msgPrevote, Msg: voteRLP, Address: clientAddr, Signature: []byte("prevote signature")}
	prevoteMsgRLPNoSig, err := prevoteMsg.PayloadNoSig()
	assert.Nil(t, err)
	prevoteMsgRLPWithSig, err := prevoteMsg.Payload()
	assert.Nil(t, err)
	return prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig
}
