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
	committeeSize := rand.Intn(maxSize - minSize) + minSize
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

// It test the page-6, from Line-14 to Line 19, StartRound() function from proposer point of view of tendermint pseudo-code.
// Please refer to the algorithm from here: https://arxiv.org/pdf/1807.04938.pdf
func TestTendermintProposerStartRound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// prepare a random size of committee, and the proposer at last committed block.
	currentCommittee := prepareCommittee()
	lastProposer := currentCommittee[len(currentCommittee) - 1].Address
	committeeSet, err := committee.NewSet(currentCommittee, lastProposer)
	if err != nil {
		t.Error(err)
	}

	currentHeight := big.NewInt(10)
	currentBlock := generateBlock(currentHeight)
	proposalHeight := big.NewInt(11)
	proposalBlock := generateBlock(proposalHeight)
	clientAddr := currentCommittee[0].Address

	// create consensus core.
	backendMock := NewMockBackend(ctrl)
	backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
	c := New(backendMock)
	// init core's context data.
	c.pendingUnminedBlocks[proposalHeight.Uint64()] = proposalBlock
	c.committeeSet = committeeSet
	c.sentProposal = false
	c.height = currentHeight
	round := int64(0)
	// since the default value of step and round is are both 0 which is to be one of the expected state, so we set them
	// into different value.
	c.step = precommitDone
	c.round = -1
	backendMock.EXPECT().LastCommittedProposal().Return(currentBlock, lastProposer)
	backendMock.EXPECT().Committee(proposalHeight.Uint64()).Return(committeeSet, nil)
	backendMock.EXPECT().SetProposedBlockHash(proposalBlock.Hash())
	backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)

	// prepare the wanted msg which will be broadcast.
	proposal := NewProposal(round, proposalHeight, int64(-1), proposalBlock)
	proposalMsg, err := Encode(proposal)
	if err != nil {
		t.Error("err")
	}
	wantedMsg, err := c.finalizeMessage(&Message{
		Code:          msgProposal,
		Msg:           proposalMsg,
		Address:       clientAddr,
		CommittedSeal: []byte{},
	})
	// should check if broadcast to wanted committeeSet with wanted MSG.
	backendMock.EXPECT().Broadcast(context.Background(), committeeSet, wantedMsg).Return(nil)
	c.startRound(context.Background(), round)

	// checking consensus state machine states
	assert.True(t, c.sentProposal)
	assert.Equal(t, c.step, propose)
	assert.Equal(t, c.Round(), round)
	assert.Nil(t, c.lockedValue)
	assert.Equal(t, c.lockedRound, int64(-1))
	assert.Nil(t, c.validValue)
	assert.Equal(t, c.validRound, int64(-1))
}

// It test the page-6, line-21, StartRound() function from follower point of view of tendermint pseudo-code.
// Please refer to the algorithm from here: https://arxiv.org/pdf/1807.04938.pdf
func TestTendermintFollowerStartRound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currentCommittee := prepareCommittee()
	lastProposer := currentCommittee[len(currentCommittee) - 1].Address
	committeeSet, err := committee.NewSet(currentCommittee, lastProposer)
	if err != nil {
		t.Error(err)
	}

	currentHeight := big.NewInt(10)
	currentBlock := generateBlock(currentHeight)
	clientAddr := currentCommittee[0].Address

	backendMock := NewMockBackend(ctrl)
	backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
	backendMock.EXPECT().LastCommittedProposal().AnyTimes().Return(currentBlock, lastProposer)

	// create consensus core.
	c := New(backendMock)
	c.committeeSet = committeeSet
	round := int64(1)
	// since the default value of step and round is are both 0 which is to be one of the expected state, so we set them
	// into different value.
	c.step = precommitDone
	c.round = -1
	c.startRound(context.Background(), round)

	// checking consensus state machine states
	assert.True(t, c.proposeTimeout.started)
	assert.Equal(t, c.step, propose)
	assert.Equal(t, c.Round(), round)
	assert.Nil(t, c.lockedValue)
	assert.Equal(t, c.lockedRound, int64(-1))
	assert.Nil(t, c.validValue)
	assert.Equal(t, c.validRound, int64(-1))
	// clean up timer otherwise it would panic due to the core object is released.
	err = c.proposeTimeout.stopTimer()
	if err != nil {
		t.Error(err)
	}
}

// It test the page-6, upon proposal logic blocks from Line22 to Line33 from tendermint pseudo-code.
// Please refer to the algorithm from here: https://arxiv.org/pdf/1807.04938.pdf
func TestTendermintUponProposal(t *testing.T) {
	// Below 4 test cases cover line 22 to line 27 of tendermint pseudo-code.
	// It test line 24 was executed and step was forwarded on line 27.
	t.Run("on proposal with validRound as (-1) proposed and local lockedRound as (-1)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		currentCommittee := prepareCommittee()
		lastProposer := currentCommittee[len(currentCommittee) - 1].Address
		committeeSet, err := committee.NewSet(currentCommittee, lastProposer)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, gomock.Any()).AnyTimes().Return(nil)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), nil)

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

		// create consensus core.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.lockedRound = -1

		err = c.handleProposal(context.Background(), msg)
		if err != nil {
			t.Error(err)
		}

		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
	})

	// It test line 24 was executed and step was forwarded on line 27 but with below different condition.
	t.Run("on proposal with validRound as (-1) proposed and local lockedRound as (1) and lockedValue as a valid value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		currentCommittee := prepareCommittee()
		committeeSet, err := committee.NewSet(currentCommittee, currentCommittee[3].Address)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, gomock.Any()).AnyTimes().Return(nil)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), nil)

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

		// create consensus core.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.lockedRound = 1 // set lockedRound as 1.
		c.lockedValue = proposalBlock

		err = c.handleProposal(context.Background(), msg)
		if err != nil {
			t.Error(err)
		}

		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
	})

	// It test line 26 was executed and step was forwarded on line 27 but with below different condition.
	t.Run("on proposal with validRound as (-1) proposed and local lockedRound as (1) and lockedValue as a nil value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		currentCommittee := prepareCommittee()
		committeeSet, err := committee.NewSet(currentCommittee, currentCommittee[3].Address)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, gomock.Any()).AnyTimes().Return(nil)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), nil)

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

		// create consensus core.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.lockedRound = 1
		c.lockedValue = nil

		err = c.handleProposal(context.Background(), msg)
		if err != nil {
			t.Error(err)
		}

		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
	})

	// It test line 26 was executed and step was forwarded on line 27 but with below different condition.
	t.Run("on proposal with invalid block, follower should step forward with voting for nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		currentCommittee := prepareCommittee()
		committeeSet, err := committee.NewSet(currentCommittee, currentCommittee[3].Address)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, gomock.Any()).AnyTimes().Return(nil)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), errors.New(
			"invalid block"))

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

		// create consensus core.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.lockedRound = -1
		c.lockedValue = proposalBlock

		err = c.handleProposal(context.Background(), msg)
		if err != nil {
			assert.Equal(t, err.Error(), "invalid block")
		}

		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
	})

	// Below 4 test cases cover line 28 to line 33 of tendermint pseudo-code.
	// It test line 30 was executed and step was forwarded on line 33.
	t.Run("on proposal with pre-vote power satisfy the quorum and on the same vr view", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		currentCommittee := prepareCommittee()
		committeeSet, err := committee.NewSet(currentCommittee, currentCommittee[3].Address)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, gomock.Any()).AnyTimes().Return(nil)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), nil)

		// condition vr >= 0 && vr < round_p, line 28.
		validRoundProposed := int64(0)
		roundProposed := int64(1)

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

		// create consensus core.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.round = roundProposed

		// condition lockedRound_p <= vr, line 29.
		c.lockedRound = -1

		// condition step_p = propose, line 28.
		c.step = propose

		// condition 2f+1 <PREVOTE, h_p, vr, id(v)>, power of pre-vote on the same valid round meets quorum, line 28.
		prevoteMsg := Message{
			Code:    msgPrevote,
			Address: currentCommittee[2].Address,
			power:   3,
		}
		c.messages.getOrCreate(validRoundProposed).AddPrevote(proposalBlock.Hash(), prevoteMsg)

		err = c.handleProposal(context.Background(), proposalMsg)
		if err != nil {
			t.Error(err)
		}

		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
	})

	// It test line 30 was executed and step was forwarded on line 33.
	t.Run("on proposal with pre-vote power satisfy the quorum and on the same vr view, but lockedRound > vr", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		currentCommittee := prepareCommittee()
		committeeSet, err := committee.NewSet(currentCommittee, currentCommittee[3].Address)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, gomock.Any()).AnyTimes().Return(nil)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), nil)

		// condition vr >= 0 && vr < round_p, line 28.
		validRoundProposed := int64(0)
		roundProposed := int64(1)

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

		// create consensus core.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.round = roundProposed

		// condition (lockedRound_p <= vr || lockedValue_p = v, line 29.
		c.lockedRound = 1
		c.lockedValue = proposalBlock

		// condition step_p = propose, line 28.
		c.step = propose

		// condition 2f+1 <PREVOTE, h_p, vr, id(v)>, power of pre-vote on the same valid round meets quorum, line 28.
		prevoteMsg := Message{
			Code:    msgPrevote,
			Address: currentCommittee[2].Address,
			power:   3,
		}
		c.messages.getOrCreate(validRoundProposed).AddPrevote(proposalBlock.Hash(), prevoteMsg)

		err = c.handleProposal(context.Background(), proposalMsg)
		if err != nil {
			t.Error(err)
		}

		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
	})

	// It test line 32 was executed and step was forwarded on line 33.
	t.Run("on proposal with pre-vote power satisfy the quorum and on the same vr view, but with un-expected locked round and locked value, engine should pre-vote for nil and step to pre-vote", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		currentCommittee := prepareCommittee()
		committeeSet, err := committee.NewSet(currentCommittee, currentCommittee[3].Address)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, gomock.Any()).AnyTimes().Return(nil)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), nil)

		// condition vr >= 0 && vr < round_p, line 28.
		validRoundProposed := int64(0)
		roundProposed := int64(1)

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

		// create consensus core.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.round = roundProposed

		// condition (lockedRound_p <= vr || lockedValue_p = v, line 29.
		c.lockedRound = 1
		c.lockedValue = nil

		// condition step_p = propose, line 28.
		c.step = propose

		// condition 2f+1 <PREVOTE, h_p, vr, id(v)>, power of pre-vote on the same valid round meets quorum, line 28.
		prevoteMsg := Message{
			Code:    msgPrevote,
			Address: currentCommittee[2].Address,
			power:   3,
		}
		c.messages.getOrCreate(validRoundProposed).AddPrevote(proposalBlock.Hash(), prevoteMsg)

		err = c.handleProposal(context.Background(), proposalMsg)
		if err != nil {
			t.Error(err)
		}

		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
	})

	// It test line 32 was executed and step was forwarded on line 33.
	t.Run("on proposal with all condition satisfied but with invalid value(block)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		currentCommittee := prepareCommittee()
		committeeSet, err := committee.NewSet(currentCommittee, currentCommittee[3].Address)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, gomock.Any()).AnyTimes().Return(nil)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), errors.New("invalid block"))

		// condition vr >= 0 && vr < round_p, line 28.
		validRoundProposed := int64(0)
		roundProposed := int64(1)

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

		// condition 2f+1 <PREVOTE, h_p, vr, id(v)>, power of pre-vote on the same valid round meets quorum, line 28.
		prevoteMsg := Message{
			Code:    msgPrevote,
			Address: currentCommittee[2].Address,
			power:   3,
		}
		c.messages.getOrCreate(validRoundProposed).AddPrevote(proposalBlock.Hash(), prevoteMsg)

		err = c.handleProposal(context.Background(), proposalMsg)
		if err != nil {
			assert.Equal(t, err.Error(), "invalid block")
		}

		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
	})
}

// It test the page-6, logic block from Line36 to Line43, Line 34 to Line 35 from tendermint pseudo-code.
func TestTendermintUponPrevote(t *testing.T) {
	t.Run("Line36 to Line43, on prevote for the first time.", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		currentCommittee := prepareCommittee()
		committeeSet, err := committee.NewSet(currentCommittee, currentCommittee[3].Address)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, gomock.Any()).AnyTimes().Return(nil)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)

		validRoundProposed := int64(0)
		roundProposed := int64(0)

		preVoteMsg := createPrevote(t, proposalBlock.Hash(), roundProposed, currentHeight, committeeSet.Committee()[0])
		// create consensus core and conditions.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.round = roundProposed
		c.lockedRound = -1
		c.lockedValue = nil
		c.validRound = -1
		c.validValue = nil
		c.setValidRoundAndValue = false
		c.step = prevote

		// condition 2f+1 <PREVOTE, h_p, round_p, id(v)>, power of pre-vote. line 36.
		receivedPrevoteMsg := Message{
			Code:    msgPrevote,
			Address: currentCommittee[2].Address,
			power:   3,
		}

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

		c.curRoundMessages.SetProposal(proposal, proposalMsg, true)
		c.curRoundMessages.AddPrevote(proposalBlock.Hash(), receivedPrevoteMsg)
		c.messages.getOrCreate(validRoundProposed).AddPrevote(proposalBlock.Hash(), receivedPrevoteMsg)

		err = c.handlePrevote(context.Background(), preVoteMsg)
		if err != nil {
			t.Error(err)
		}

		assert.True(t, c.sentPrecommit)
		assert.Equal(t, c.step, precommit)
		assert.Equal(t, c.lockedRound, roundProposed)
		assert.Equal(t, c.lockedValue, proposalBlock)
		assert.Equal(t, c.validRound, roundProposed)
		assert.Equal(t, c.validValue, proposalBlock)
	})

	t.Run("Line36 to Line41, on prevote for the first time, with step > prevote.", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		currentCommittee := prepareCommittee()
		committeeSet, err := committee.NewSet(currentCommittee, currentCommittee[3].Address)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, gomock.Any()).AnyTimes().Return(nil)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)

		validRoundProposed := int64(0)
		roundProposed := int64(0)

		preVoteMsg := createPrevote(t, proposalBlock.Hash(), roundProposed, currentHeight, committeeSet.Committee()[0])
		// create consensus core and conditions.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.round = roundProposed
		c.lockedRound = -1
		c.lockedValue = nil
		c.validRound = -1
		c.validValue = nil
		c.setValidRoundAndValue = false
		c.step = precommitDone

		// condition 2f+1 <PREVOTE, h_p, round_p, id(v)>, power of pre-vote. line 36.
		receivedPrevoteMsg := Message{
			Code:    msgPrevote,
			Address: currentCommittee[2].Address,
			power:   3,
		}
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

		c.curRoundMessages.SetProposal(proposal, proposalMsg, true)
		c.curRoundMessages.AddPrevote(proposalBlock.Hash(), receivedPrevoteMsg)
		c.messages.getOrCreate(validRoundProposed).AddPrevote(proposalBlock.Hash(), receivedPrevoteMsg)

		err = c.handlePrevote(context.Background(), preVoteMsg)
		if err != nil {
			t.Error(err)
		}

		assert.False(t, c.sentPrecommit)
		assert.Equal(t, c.step, precommitDone)
		assert.Equal(t, c.lockedRound, int64(-1))
		assert.Nil(t, c.lockedValue)
		assert.Equal(t, c.validRound, roundProposed)
		assert.Equal(t, c.validValue, proposalBlock)
	})

	t.Run("Line34 to Line35, schedule for prevote timeout.", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		currentCommittee := prepareCommittee()
		committeeSet, err := committee.NewSet(currentCommittee, currentCommittee[3].Address)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, gomock.Any()).AnyTimes().Return(nil)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)

		validRoundProposed := int64(0)
		roundProposed := int64(0)

		preVoteMsg := createPrevote(t, proposalBlock.Hash(), roundProposed, currentHeight, committeeSet.Committee()[0])
		// create consensus core and conditions.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.round = roundProposed
		c.lockedRound = roundProposed
		c.lockedValue = proposalBlock
		c.validRound = roundProposed
		c.validValue = proposalBlock
		c.setValidRoundAndValue = true
		c.step = prevote

		// condition 2f+1 <PREVOTE, h_p, round_p, id(v)>, power of pre-vote. line 36.
		receivedPrevoteMsg := Message{
			Code:    msgPrevote,
			Address: currentCommittee[2].Address,
			power:   3,
		}
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

		c.curRoundMessages.SetProposal(proposal, proposalMsg, true)
		c.curRoundMessages.AddPrevote(proposalBlock.Hash(), receivedPrevoteMsg)
		c.messages.getOrCreate(validRoundProposed).AddPrevote(proposalBlock.Hash(), receivedPrevoteMsg)

		err = c.handlePrevote(context.Background(), preVoteMsg)
		if err != nil {
			t.Error(err)
		}
		assert.True(t, c.prevoteTimeout.started)
		// clean up timer otherwise it would panic due to the core object is released.
		err = c.prevoteTimeout.stopTimer()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Line44 to Line46, step from prevote to precommit by voting for nil.", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		currentCommittee := prepareCommittee()
		committeeSet, err := committee.NewSet(currentCommittee, currentCommittee[3].Address)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, gomock.Any()).AnyTimes().Return(nil)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)

		validRoundProposed := int64(0)
		roundProposed := int64(0)

		preVoteMsg := createPrevote(t, proposalBlock.Hash(), roundProposed, currentHeight, committeeSet.Committee()[0])
		// create consensus core and conditions.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.round = roundProposed
		c.lockedRound = roundProposed
		c.lockedValue = proposalBlock
		c.validRound = roundProposed
		c.validValue = proposalBlock
		c.setValidRoundAndValue = true
		c.step = prevote

		// condition 2f+1 <PREVOTE, h_p, round_p, id(v)>, power of pre-vote. line 36.
		receivedPrevoteMsg := Message{
			Code:    msgPrevote,
			Address: currentCommittee[2].Address,
			power:   3,
		}
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

		c.curRoundMessages.SetProposal(proposal, proposalMsg, true)
		c.curRoundMessages.AddPrevote(common.Hash{}, receivedPrevoteMsg)
		c.messages.getOrCreate(validRoundProposed).AddPrevote(proposalBlock.Hash(), receivedPrevoteMsg)

		err = c.handlePrevote(context.Background(), preVoteMsg)
		if err != nil {
			t.Error(err)
		}
		assert.True(t, c.sentPrecommit)
		assert.Equal(t, c.step, precommit)
	})
}

// It test the page-6, logic block from Line 47 to Line 56 from tendermint pseudo-code.
func TestTendermintUponPrecommit(t *testing.T) {
	t.Run("Line 47 to Line48, schedule for precommit timeout.", func(t *testing.T) {
		// todo: test line 47 - 48, schedule precommit timeout.
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		currentCommittee := prepareCommittee()
		committeeSet, err := committee.NewSet(currentCommittee, currentCommittee[3].Address)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, gomock.Any()).AnyTimes().Return(nil)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)

		validRoundProposed := int64(0)
		roundProposed := int64(0)

		preVoteMsg := createPrevote(t, proposalBlock.Hash(), roundProposed, currentHeight, committeeSet.Committee()[0])
		// create consensus core and conditions.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.round = roundProposed
		c.lockedRound = roundProposed
		c.lockedValue = proposalBlock
		c.validRound = roundProposed
		c.validValue = proposalBlock
		c.setValidRoundAndValue = true
		c.step = prevote

		// condition 2f+1 <PREVOTE, h_p, round_p, id(v)>, power of pre-vote. line 36.
		receivedPrevoteMsg := Message{
			Code:    msgPrevote,
			Address: currentCommittee[2].Address,
			power:   3,
		}
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

		c.curRoundMessages.SetProposal(proposal, proposalMsg, true)
		c.curRoundMessages.AddPrevote(proposalBlock.Hash(), receivedPrevoteMsg)
		c.messages.getOrCreate(validRoundProposed).AddPrevote(proposalBlock.Hash(), receivedPrevoteMsg)

		err = c.handlePrevote(context.Background(), preVoteMsg)
		if err != nil {
			t.Error(err)
		}
		assert.True(t, c.prevoteTimeout.started)
		// clean up timer otherwise it would panic due to the core object is released.
		err = c.prevoteTimeout.stopTimer()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Line 49 - Line 54, start round with new height.", func(t *testing.T) {
		// todo: start round with new height.
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		currentCommittee := prepareCommittee()
		committeeSet, err := committee.NewSet(currentCommittee, currentCommittee[3].Address)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, gomock.Any()).AnyTimes().Return(nil)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)

		validRoundProposed := int64(0)
		roundProposed := int64(0)

		preVoteMsg := createPrevote(t, proposalBlock.Hash(), roundProposed, currentHeight, committeeSet.Committee()[0])
		// create consensus core and conditions.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.round = roundProposed
		c.lockedRound = roundProposed
		c.lockedValue = proposalBlock
		c.validRound = roundProposed
		c.validValue = proposalBlock
		c.setValidRoundAndValue = true
		c.step = prevote

		// condition 2f+1 <PREVOTE, h_p, round_p, id(v)>, power of pre-vote. line 36.
		receivedPrevoteMsg := Message{
			Code:    msgPrevote,
			Address: currentCommittee[2].Address,
			power:   3,
		}
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

		c.curRoundMessages.SetProposal(proposal, proposalMsg, true)
		c.curRoundMessages.AddPrevote(proposalBlock.Hash(), receivedPrevoteMsg)
		c.messages.getOrCreate(validRoundProposed).AddPrevote(proposalBlock.Hash(), receivedPrevoteMsg)

		err = c.handlePrevote(context.Background(), preVoteMsg)
		if err != nil {
			t.Error(err)
		}
		assert.True(t, c.prevoteTimeout.started)
		// clean up timer otherwise it would panic due to the core object is released.
		err = c.prevoteTimeout.stopTimer()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Line55 to Line56, start round with higher round from committee members.", func(t *testing.T) {
		// todo: start round with higher round from committee members.
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		currentCommittee := prepareCommittee()
		committeeSet, err := committee.NewSet(currentCommittee, currentCommittee[3].Address)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, gomock.Any()).AnyTimes().Return(nil)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)

		validRoundProposed := int64(0)
		roundProposed := int64(0)

		preVoteMsg := createPrevote(t, proposalBlock.Hash(), roundProposed, currentHeight, committeeSet.Committee()[0])
		// create consensus core and conditions.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.round = roundProposed
		c.lockedRound = roundProposed
		c.lockedValue = proposalBlock
		c.validRound = roundProposed
		c.validValue = proposalBlock
		c.setValidRoundAndValue = true
		c.step = prevote

		// condition 2f+1 <PREVOTE, h_p, round_p, id(v)>, power of pre-vote. line 36.
		receivedPrevoteMsg := Message{
			Code:    msgPrevote,
			Address: currentCommittee[2].Address,
			power:   3,
		}
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

		c.curRoundMessages.SetProposal(proposal, proposalMsg, true)
		c.curRoundMessages.AddPrevote(proposalBlock.Hash(), receivedPrevoteMsg)
		c.messages.getOrCreate(validRoundProposed).AddPrevote(proposalBlock.Hash(), receivedPrevoteMsg)

		err = c.handlePrevote(context.Background(), preVoteMsg)
		if err != nil {
			t.Error(err)
		}
		assert.True(t, c.prevoteTimeout.started)
		// clean up timer otherwise it would panic due to the core object is released.
		err = c.prevoteTimeout.stopTimer()
		if err != nil {
			t.Error(err)
		}
	})
}
