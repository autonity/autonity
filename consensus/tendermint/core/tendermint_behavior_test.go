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

// It test the page-6, line 22 to line 27, on new proposal logic of tendermint pseudo-code.
// Please refer to the algorithm from here: https://arxiv.org/pdf/1807.04938.pdf
func TestTendermintNewProposal(t *testing.T) {
	// Below 4 test cases cover line 22 to line 27 of tendermint pseudo-code.
	// It test line 24 was executed and step was forwarded on line 27.

	t.Run("on proposal with validRound as (-1) proposed and local lockedRound as (-1)", func(t *testing.T) {
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
		// create consensus core.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.lockedRound = -1
		c.step = propose

		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), nil)

		// prepare input msg
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

		// prepare the wanted msg which will be broadcast.
		vote := Vote{
			Round:             0,
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
		// checking consensus state machine states
		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
		assert.Nil(t, c.lockedValue)
		assert.Equal(t, c.lockedRound, int64(-1))
		assert.Nil(t, c.validValue)
		assert.Equal(t, c.validRound, int64(-1))
	})

	// It test line 24 was executed and step was forwarded on line 27 but with below different condition.
	t.Run("on proposal with validRound as (-1) proposed and local lockedRound as (1) and lockedValue as the same value proposed", func(t *testing.T) {
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
		// create consensus core.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
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
		// checking consensus state machine states
		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
		assert.Equal(t, c.lockedValue, proposalBlock)
		assert.Equal(t, c.lockedRound, int64(1))
		assert.Equal(t, c.validValue, proposalBlock)
		assert.Equal(t, c.validRound, int64(1))
	})

	// It test line 26 was executed and step was forwarded on line 27 but with below different condition.
	t.Run("on proposal with validRound as (-1) proposed and local lockedRound as (1) and locked at different value, vote for nil", func(t *testing.T) {
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
		// create consensus core.
		lockedValue := generateBlock(big.NewInt(10))
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
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

		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
		assert.Equal(t, c.lockedValue, lockedValue)
		assert.Equal(t, c.lockedRound, int64(1))
		assert.Equal(t, c.validValue, lockedValue)
		assert.Equal(t, c.validRound, int64(1))
	})

	// It test line 26 was executed and step was forwarded on line 27 but with invalid value proposed.
	t.Run("on proposal with invalid block, follower should step forward with voting for nil", func(t *testing.T) {
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
		// create consensus core.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
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

		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
		assert.Nil(t, c.lockedValue)
		assert.Equal(t, c.lockedRound, int64(-1))
		assert.Nil(t, c.validValue)
		assert.Equal(t, c.validRound, int64(-1))
	})
}