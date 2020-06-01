package core

import (
	"context"
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

// The following tests aim to test lines 28 - 33 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestOldProposal(t *testing.T) {
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

	t.Run("receive proposal with vr >= 0 and client's lockedRound <= vr", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))

		// vr >= 0 && vr < round_p
		proposalValidRound := int64(rand.Intn(int(currentRound)))
		// c.lockedRound <= vr
		choice := rand.Intn(2)
		clientLockedRound := int64(-1)
		if choice != 0 {
			clientLockedRound = int64(rand.Intn(int(proposalValidRound + 1)))
		}

		var proposal Proposal
		proposalMsg := generateBlockProposal(t, currentRound, currentHeight, proposalValidRound, members[currentRound].Address, false)
		// we have to do this because encoding and decoding changes some default values and thus same blocks are no longer equal
		err := proposalMsg.Decode(&proposal)
		assert.Nil(t, err)

		// expected message to be broadcast
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := preparePrevote(t, currentRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr)

		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(propose)
		c.lockedRound = clientLockedRound
		c.validRound = clientLockedRound
		// Although the following is not possible it is required to ensure that c.lockRound <= proposalValidRound is
		// responsible for sending the prevote for the incoming proposal
		c.lockedValue = nil
		c.validValue = nil
		c.messages.getOrCreate(proposalValidRound).AddPrevote(proposal.ProposalBlock.Hash(), Message{Code: msgPrevote, power: c.CommitteeSet().Quorum()})

		backendMock.EXPECT().VerifyProposal(*proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, prevoteMsgRLPWithSig).Return(nil)

		err = c.handleCheckedMsg(context.Background(), proposalMsg, members[currentRound])
		assert.Nil(t, err)
		assert.Equal(t, prevote, c.step)
		assert.Nil(t, c.lockedValue)
		assert.Equal(t, clientLockedRound, c.lockedRound)
		assert.Nil(t, c.validValue)
		assert.Equal(t, clientLockedRound, c.validRound)
	})
	t.Run("receive proposal with vr >= 0 and client's lockedValue is same as proposal block", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))

		// vr >= 0 && vr < round_p
		proposalValidRound := int64(rand.Intn(int(currentRound)))
		var proposal Proposal
		proposalMsg := generateBlockProposal(t, currentRound, currentHeight, proposalValidRound, members[currentRound].Address, false)
		// we have to do this because encoding and decoding changes some default values and thus same blocks are no longer equal
		err := proposalMsg.Decode(&proposal)
		assert.Nil(t, err)

		// expected message to be broadcast
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := preparePrevote(t, currentRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr)

		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(propose)
		// Although the following is not possible it is required to ensure that c.lockedValue = proposal is responsible
		// for sending the prevote for the incoming proposal
		c.lockedRound = proposalValidRound + 1
		c.validRound = proposalValidRound + 1
		c.lockedValue = proposal.ProposalBlock
		c.validValue = proposal.ProposalBlock
		c.messages.getOrCreate(proposalValidRound).AddPrevote(proposal.ProposalBlock.Hash(), Message{Code: msgPrevote, power: c.CommitteeSet().Quorum()})

		backendMock.EXPECT().VerifyProposal(*proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, prevoteMsgRLPWithSig).Return(nil)

		err = c.handleCheckedMsg(context.Background(), proposalMsg, members[currentRound])
		assert.Nil(t, err)
		assert.Equal(t, prevote, c.step)
		assert.Equal(t, proposalValidRound+1, c.lockedRound)
		assert.Equal(t, proposalValidRound+1, c.validRound)
		assert.Equal(t, proposal.ProposalBlock, c.lockedValue)
		assert.Equal(t, proposal.ProposalBlock, c.validValue)
	})
	t.Run("receive proposal with vr >= 0 and clients is lockedRound > vr with a different value", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		clientLockedValue := generateBlock(currentHeight)

		// vr >= 0 && vr < round_p
		proposalValidRound := int64(rand.Intn(int(currentRound)))
		var proposal Proposal
		proposalMsg := generateBlockProposal(t, currentRound, currentHeight, proposalValidRound, members[currentRound].Address, false)
		// we have to do this because encoding and decoding changes some default values and thus same blocks are no longer equal
		err := proposalMsg.Decode(&proposal)
		assert.Nil(t, err)

		// expected message to be broadcast
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := preparePrevote(t, currentRound, currentHeight, common.Hash{}, clientAddr)

		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(propose)
		// Although the following is not possible it is required to ensure that c.lockedValue = proposal is responsible
		// for sending the prevote for the incoming proposal
		c.lockedRound = proposalValidRound + 1
		c.validRound = proposalValidRound + 1
		c.lockedValue = clientLockedValue
		c.validValue = clientLockedValue
		c.messages.getOrCreate(proposalValidRound).AddPrevote(proposal.ProposalBlock.Hash(), Message{Code: msgPrevote, power: c.CommitteeSet().Quorum()})

		backendMock.EXPECT().VerifyProposal(*proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, prevoteMsgRLPWithSig).Return(nil)

		err = c.handleCheckedMsg(context.Background(), proposalMsg, members[currentRound])
		assert.Nil(t, err)
		assert.Equal(t, prevote, c.step)
		assert.Equal(t, proposalValidRound+1, c.lockedRound)
		assert.Equal(t, proposalValidRound+1, c.validRound)
		assert.Equal(t, clientLockedValue, c.lockedValue)
		assert.Equal(t, clientLockedValue, c.validValue)
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
