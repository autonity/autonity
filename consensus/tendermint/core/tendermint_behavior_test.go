package core

import (
	"context"
	"errors"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"github.com/clearmatics/autonity/core/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func prepareCommittee() types.Committee {
	// prepare committee.
	minSize := 4
	maxSize := 15
	committeeSize := rand.Intn(maxSize-minSize) + minSize
	committeeSet := types.Committee{}
	for i := 1; i <= committeeSize; i++ {
		hexString := "0x01234567890" + strconv.Itoa(i)
		member := types.CommitteeMember{
			Address:     common.HexToAddress(hexString),
			VotingPower: new(big.Int).SetInt64(1),
		}
		committeeSet = append(committeeSet, member)
	}
	return committeeSet
}

func generateBlock(height *big.Int) *types.Block {
	header := &types.Header{Number: height}
	block := types.NewBlock(header, nil, nil, nil)
	return block
}


// It test the page-6, line 28 to line 33, on proposal logic of tendermint pseudo-code.
// Please refer to the algorithm from here: https://arxiv.org/pdf/1807.04938.pdf
func TestTendermintOnOldProposal(t *testing.T) {
	// It test line 30 was executed and step was forwarded on line 33.
	t.Run("on proposal with pre-vote power satisfy the quorum and on the same vr view", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// prepare a random size of committee, and the proposer at last committed block.
		currentCommittee := prepareCommittee()
		lastProposer := currentCommittee[len(currentCommittee)-1].Address
		committeeSet, err := committee.NewSet(currentCommittee, lastProposer)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
		// condition vr >= 0 && vr < round_p, line 28.
		validRoundProposed := int64(0)
		roundProposed := int64(1)

		// create consensus core, and prepare context.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.round = roundProposed
		// condition lockedRound_p <= vr, line 29.
		c.lockedRound = -1
		// condition step_p = propose, line 28.
		c.step = propose
		// condition 2f+1 <PREVOTE, h_p, vr, id(v)>, power of pre-vote on the same valid round meets quorum, line 28.
		receivedPrevoteMsg := Message{
			Code:    msgPrevote,
			Address: currentCommittee[2].Address,
			power:   c.CommitteeSet().Quorum(),
		}
		c.messages.getOrCreate(validRoundProposed).AddPrevote(proposalBlock.Hash(), receivedPrevoteMsg)

		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), nil)

		// prepare input msg
		proposal := NewProposal(roundProposed, currentHeight, validRoundProposed, proposalBlock)
		encodedProposal, err := Encode(proposal)
		if err != nil {
			t.Error(err)
		}
		proposalMsg := &Message{
			Code:    msgProposal,
			Msg:     encodedProposal,
			Address: currentCommittee[1].Address,
		}
		// prepare the wanted msg which will be broadcast.
		vote := Vote{
			Round:             roundProposed,
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

		err = c.handleProposal(context.Background(), proposalMsg)
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

	// It test line 30 was executed and step was forwarded on line 33.
	// FAILED TOO.
	t.Run("on proposal with pre-vote power satisfy the quorum and on the same vr view, but lockedRound > vr", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// prepare a random size of committee, and the proposer at last committed block.
		currentCommittee := prepareCommittee()
		lastProposer := currentCommittee[len(currentCommittee)-1].Address
		committeeSet, err := committee.NewSet(currentCommittee, lastProposer)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address
		// condition vr >= 0 && vr < round_p, line 28.
		validRoundProposed := int64(0)
		roundProposed := int64(1)

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
		// create consensus core.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.round = roundProposed
		// condition (lockedRound_p <= vr || lockedValue_p = v, line 29.
		c.lockedRound = 1
		c.lockedValue = proposalBlock
		c.validRound = 1
		c.validValue = proposalBlock
		// condition step_p = propose, line 28.
		c.step = propose
		// condition 2f+1 <PREVOTE, h_p, vr, id(v)>, power of pre-vote on the same valid round meets quorum, line 28.
		receivedPrevoteMsg := Message{
			Code:    msgPrevote,
			Address: currentCommittee[2].Address,
			power:   c.committeeSet.Quorum(),
		}
		c.messages.getOrCreate(validRoundProposed).AddPrevote(proposalBlock.Hash(), receivedPrevoteMsg)

		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), nil)

		// prepare input msg.
		proposal := NewProposal(roundProposed, currentHeight, validRoundProposed, proposalBlock)
		encodedProposal, err := Encode(proposal)
		if err != nil {
			t.Error(err)
		}
		proposalMsg := &Message{
			Code:    msgProposal,
			Msg:     encodedProposal,
			Address: currentCommittee[1].Address,
		}

		// prepare the wanted msg which will be broadcast.
		vote := Vote{
			Round:             roundProposed,
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

		err = c.handleProposal(context.Background(), proposalMsg)
		if err != nil {
			t.Error(err)
		}

		// checking internal state of tendermint.
		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
		assert.Equal(t, c.lockedValue, proposalBlock)
		assert.Equal(t, c.lockedRound, int64(1))
		assert.Equal(t, c.validValue, proposalBlock)
		assert.Equal(t, c.validRound, int64(1))
	})

	// It test line 32 was executed and step was forwarded on line 33.
	t.Run("on proposal with pre-vote power satisfy the quorum and on the same vr view, but with un-expected locked round and locked value, engine should pre-vote for nil and step to pre-vote", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// prepare a random size of committee, and the proposer at last committed block.
		currentCommittee := prepareCommittee()
		lastProposer := currentCommittee[len(currentCommittee)-1].Address
		committeeSet, err := committee.NewSet(currentCommittee, lastProposer)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address
		// condition vr >= 0 && vr < round_p, line 28.
		validRoundProposed := int64(0)
		roundProposed := int64(1)

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
		// create consensus core.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.round = roundProposed
		// condition (lockedRound_p <= vr || lockedValue_p = v, line 29.
		c.lockedRound = 1
		lockedValue := generateBlock(currentHeight)
		c.lockedValue = lockedValue
		// condition step_p = propose, line 28.
		c.step = propose

		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), nil)

		// condition 2f+1 <PREVOTE, h_p, vr, id(v)>, power of pre-vote on the same valid round meets quorum, line 28.
		receivedPrevoteMsg := Message{
			Code:    msgPrevote,
			Address: currentCommittee[2].Address,
			power:   c.CommitteeSet().Quorum(),
		}
		c.messages.getOrCreate(validRoundProposed).AddPrevote(proposalBlock.Hash(), receivedPrevoteMsg)

		// prepare input message.
		proposal := NewProposal(roundProposed, currentHeight, validRoundProposed, proposalBlock)
		encodedProposal, err := Encode(proposal)
		if err != nil {
			t.Error(err)
		}
		proposalMsg := &Message{
			Code:    msgProposal,
			Msg:     encodedProposal,
			Address: currentCommittee[1].Address,
		}

		// prepare the wanted msg which will be broadcast.
		vote := Vote{
			Round:             roundProposed,
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

		err = c.handleProposal(context.Background(), proposalMsg)
		if err != nil {
			t.Error(err)
		}

		// checking internal state of tendermint.
		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
		assert.Equal(t, c.lockedValue, lockedValue)
		assert.Equal(t, c.lockedRound, int64(1))
		assert.Nil(t, c.validValue)
		assert.Equal(t, c.validRound, int64(-1))
	})

	// It test line 32 was executed and step was forwarded on line 33.
	t.Run("on proposal with all condition satisfied but with invalid value(block)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// prepare a random size of committee, and the proposer at last committed block.
		currentCommittee := prepareCommittee()
		lastProposer := currentCommittee[len(currentCommittee)-1].Address
		committeeSet, err := committee.NewSet(currentCommittee, lastProposer)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address
		// condition vr >= 0 && vr < round_p, line 28.
		validRoundProposed := int64(0)
		roundProposed := int64(1)
		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
		// create consensus core.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.round = roundProposed
		// condition (lockedRound_p <= vr || lockedValue_p = v, line 29.
		c.lockedRound = 0
		c.lockedValue = proposalBlock
		// condition step_p = propose, line 28.
		c.step = propose

		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), errors.New("invalid block"))

		// prepare input message.
		proposal := NewProposal(roundProposed, currentHeight, validRoundProposed, proposalBlock)
		encodedProposal, err := Encode(proposal)
		if err != nil {
			t.Error(err)
		}

		proposalMsg := &Message{
			Code:    msgProposal,
			Msg:     encodedProposal,
			Address: currentCommittee[1].Address,
		}

		// condition 2f+1 <PREVOTE, h_p, vr, id(v)>, power of pre-vote on the same valid round meets quorum, line 28.
		receivedPrevoteMsg := Message{
			Code:    msgPrevote,
			Address: currentCommittee[2].Address,
			power:   c.CommitteeSet().Quorum(),
		}
		c.messages.getOrCreate(validRoundProposed).AddPrevote(proposalBlock.Hash(), receivedPrevoteMsg)

		// prepare the wanted msg which will be broadcast.
		vote := Vote{
			Round:             roundProposed,
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

		err = c.handleProposal(context.Background(), proposalMsg)
		if err != nil {
			assert.Equal(t, err.Error(), "invalid block")
		}

		// checking internal state of tendermint.
		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
		assert.Equal(t, c.lockedValue, proposalBlock)
		assert.Equal(t, c.lockedRound, int64(0))
		assert.Nil(t, c.validValue)
		assert.Equal(t, c.validRound, int64(-1))
	})
}
