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

var maxHeightOrRound = 100

// The following tests aim to test lines 22 - 27 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestTendermintNewProposal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backendMock := NewMockBackend(ctrl)

	committeeSet := prepareCommittee(t)
	members := committeeSet.Committee()
	currentHeight := big.NewInt(int64(rand.Intn(maxHeightOrRound) + 1))
	currentRound := int64(rand.Intn(maxHeightOrRound))
	clientAddr := members[0].Address

	proposalBlock := generateBlock(currentHeight) //Probably don't need it as it is only required for a few cases

	backendMock.EXPECT().Address().Return(clientAddr)
	c := New(backendMock)
	c.setHeight(currentHeight)
	c.setRound(currentRound)
	c.setCommitteeSet(committeeSet)

	preparePrevote := func(t *testing.T, round int64, height *big.Int, blockHash common.Hash, clientAddr common.Address) (*Message, []byte, []byte) {
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

	t.Run("receive invalid proposal for current round", func(t *testing.T) {
		c.setStep(propose)

		var invalidProposal Proposal
		// members[currentRound] means that the sender is the proposer for the current round
		invalidMsg := generateInvalidBlockProposal(t, currentHeight, currentRound, members[currentRound].Address)
		err := invalidMsg.Decode(&invalidProposal)
		assert.Nil(t, err)

		// prepare prevote nil
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := preparePrevote(t, currentRound, currentHeight, common.Hash{}, clientAddr)

		backendMock.EXPECT().VerifyProposal(*invalidProposal.ProposalBlock).Return(time.Duration(1), errors.New("invalid proposal"))
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, prevoteMsgRLPWithSig).Return(nil)

		err = c.handleCheckedMsg(context.Background(), invalidMsg, members[currentRound])
		assert.Error(t, err, "expected an error for invalid proposal")
	})
	t.Run("receive proposal with validRound = -1 and client's lockedRound = -1", func(t *testing.T) {
		c.lockedRound = -1
		c.step = propose

		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), nil)

		// prepare input msg
		validRoundProposed := int64(-1)
		proposal := NewProposal(currentRound, currentHeight, validRoundProposed, proposalBlock)
		encodedProposal, err := Encode(proposal)
		if err != nil {
			t.Error(err)
		}
		msg := &Message{
			Code:    msgProposal,
			Msg:     encodedProposal,
			Address: clientAddr,
		}

		// prepare the wanted msg which will be broadcast.
		vote := Vote{
			Round:             currentRound,
			Height:            currentHeight,
			ProposedBlockHash: proposalBlock.Hash(),
		}
		prevoteMsg, err := Encode(&vote)
		if err != nil {
			t.Error("err")
		}
		wantedMsg, err := c.finalizeMessage(&Message{
			Code:          msgPrevote,
			Msg:           prevoteMsg,
			Address:       clientAddr,
			CommittedSeal: []byte{},
		})
		if err != nil {
			t.Error(err)
		}

		// should check if broadcast to wanted committeeSet with wanted prevote msg.
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, wantedMsg).Return(nil)

		err = c.handleProposal(context.Background(), msg)
		if err != nil {
			t.Error(err)
		}
		// checking internal state of tendermint.
		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
		assert.Nil(t, c.lockedValue)
		assert.Equal(t, c.lockedRound, int64(-1))
		assert.Nil(t, c.validValue)
		assert.Equal(t, c.validRound, int64(-1))
	})

	// It test line 24 was executed and step was forwarded on line 27 but with below different condition.
	t.Run("on proposal with validRound as (-1) proposed and local lockedRound as (1), but locked at the same value as proposed already.", func(t *testing.T) {
		c.lockedRound = 1 // set lockedRound as 1.
		c.lockedValue = proposalBlock
		c.validRound = 1
		c.validValue = proposalBlock
		c.step = propose

		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), nil)

		validRoundProposed := int64(-1)
		proposal := NewProposal(0, currentHeight, validRoundProposed, proposalBlock)
		encodedProposal, err := Encode(proposal)
		if err != nil {
			t.Error(err)
		}

		// prepare input msg
		msg := &Message{
			Code:    msgProposal,
			Msg:     encodedProposal,
			Address: clientAddr,
		}

		// prepare the wanted msg which will be broadcast.
		vote := Vote{
			Round:             0,
			Height:            currentHeight,
			ProposedBlockHash: proposalBlock.Hash(),
			//ProposedBlockHash: common.Hash{},
		}
		prevoteMsg, err := Encode(&vote)
		if err != nil {
			t.Error("err")
		}
		wantedMsg, err := c.finalizeMessage(&Message{
			Code:          msgPrevote,
			Msg:           prevoteMsg,
			Address:       clientAddr,
			CommittedSeal: []byte{},
		})
		if err != nil {
			t.Error(err)
		}

		// should check if broadcast to wanted committeeSet with wanted prevote msg.
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, wantedMsg).Return(nil)
		err = c.handleProposal(context.Background(), msg)
		if err != nil {
			t.Error(err)
		}
		// checking internal state of tendermint
		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
		assert.Equal(t, c.lockedValue, proposalBlock)
		assert.Equal(t, c.lockedRound, int64(1))
		assert.Equal(t, c.validValue, proposalBlock)
		assert.Equal(t, c.validRound, int64(1))
	})

	// It test line 26 was executed and step was forwarded on line 27 but with below different condition.
	t.Run("on proposal with validRound as (-1) proposed and local lockedRound as (1) and locked at different value, vote for nil", func(t *testing.T) {
		lockedValue := generateBlock(big.NewInt(11))
		c.lockedRound = 1
		c.lockedValue = lockedValue
		c.validRound = 1
		c.validValue = lockedValue

		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), nil)

		// prepare input proposal msg.
		validRoundProposed := int64(-1)
		proposal := NewProposal(0, currentHeight, validRoundProposed, proposalBlock)
		encodedProposal, err := Encode(proposal)
		if err != nil {
			t.Error(err)
		}

		msg := &Message{
			Code:    msgProposal,
			Msg:     encodedProposal,
			Address: clientAddr,
		}

		// prepare the wanted vote for nil msg which will be broadcast.
		vote := Vote{
			Round:             0,
			Height:            currentHeight,
			ProposedBlockHash: common.Hash{},
		}
		prevoteMsg, err := Encode(&vote)
		if err != nil {
			t.Error("err")
		}
		wantedMsg, err := c.finalizeMessage(&Message{
			Code:          msgPrevote,
			Msg:           prevoteMsg,
			Address:       clientAddr,
			CommittedSeal: []byte{},
		})
		if err != nil {
			t.Error(err)
		}

		// should check if broadcast to wanted committeeSet with wanted prevote msg.
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, wantedMsg).Return(nil)
		err = c.handleProposal(context.Background(), msg)
		if err != nil {
			t.Error(err)
		}

		// checking internal state of tendermint.
		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
		assert.Equal(t, c.lockedValue, lockedValue)
		assert.Equal(t, c.lockedRound, int64(1))
		assert.Equal(t, c.validValue, lockedValue)
		assert.Equal(t, c.validRound, int64(1))
	})

	// It test line 26 was executed and step was forwarded on line 27 but with invalid value proposed.
	t.Run("on proposal with invalid block, follower should step forward with voting for nil", func(t *testing.T) {
		c.lockedRound = -1
		c.lockedValue = nil
		c.validRound = -1
		c.validValue = nil

		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), errors.New("invalid block"))

		// prepare input proposal msg.
		validRoundProposed := int64(-1)
		proposal := NewProposal(0, currentHeight, validRoundProposed, proposalBlock)
		encodedProposal, err := Encode(proposal)
		if err != nil {
			t.Error(err)
		}

		msg := &Message{
			Code:    msgProposal,
			Msg:     encodedProposal,
			Address: clientAddr,
		}

		// prepare the wanted vote for nil msg which will be broadcast.
		vote := Vote{
			Round:             0,
			Height:            currentHeight,
			ProposedBlockHash: common.Hash{},
		}
		prevoteMsg, err := Encode(&vote)
		if err != nil {
			t.Error("err")
		}
		wantedMsg, err := c.finalizeMessage(&Message{
			Code:          msgPrevote,
			Msg:           prevoteMsg,
			Address:       clientAddr,
			CommittedSeal: []byte{},
		})
		if err != nil {
			t.Error(err)
		}

		// should check if broadcast to wanted committeeSet with wanted prevote msg.
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, wantedMsg).Return(nil)
		err = c.handleProposal(context.Background(), msg)
		if err != nil {
			assert.Equal(t, err.Error(), "invalid block")
		}

		// checking inernal state of tendermint.
		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
		assert.Nil(t, c.lockedValue)
		assert.Equal(t, c.lockedRound, int64(-1))
		assert.Nil(t, c.validValue)
		assert.Equal(t, c.validRound, int64(-1))
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

func prepareCommittee(t *testing.T) *committee.Set {
	minSize := 4
	maxSize := 100
	cSize := rand.Intn(maxSize-minSize) + minSize
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

func generateInvalidBlockProposal(t *testing.T, h *big.Int, r int64, src common.Address) *Message {
	header := &types.Header{Number: big.NewInt(int64(rand.Intn(maxHeightOrRound)))}
	header.Difficulty = nil
	block := types.NewBlock(header, nil, nil, nil)
	proposal := NewProposal(r, h, -1, block)
	proposalRlp, err := Encode(proposal)
	assert.Nil(t, err)
	return &Message{
		Code:    msgProposal,
		Msg:     proposalRlp,
		Address: src,
	}
}
