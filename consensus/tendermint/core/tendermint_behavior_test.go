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
	"testing"
	"time"
)

func prepareCommittee() types.Committee {
	// prepare committee.
	member1 := types.CommitteeMember{
		Address:     common.HexToAddress("0x01234567890"),
		VotingPower: new(big.Int).SetInt64(1),
	}
	member2 := types.CommitteeMember{
		Address:     common.HexToAddress("0x01234567891"),
		VotingPower: new(big.Int).SetInt64(1),
	}
	member3 := types.CommitteeMember{
		Address:     common.HexToAddress("0x01234567892"),
		VotingPower: new(big.Int).SetInt64(1),
	}
	member4 := types.CommitteeMember{
		Address:     common.HexToAddress("0x01234567892"),
		VotingPower: new(big.Int).SetInt64(1),
	}
	return types.Committee{member1, member2, member3, member4}
}

func generateBlock(height *big.Int) *types.Block {
	header := &types.Header{Number: height}
	block := types.NewBlock(header, nil, nil, nil)
	return block
}

// It test the page-6, line-11, StartRound() function from proposer point of view of tendermint pseudo-code.
// Please refer to the algorithm from here: https://arxiv.org/pdf/1807.04938.pdf
func TestTendermintProposerStartRound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currentCommittee := prepareCommittee()
	committeeSet, err := committee.NewSet(currentCommittee, currentCommittee[3].Address)
	if err != nil {
		t.Error(err)
	}

	currentHeight := big.NewInt(10)
	currentBlock := generateBlock(currentHeight)
	proposalHeight := big.NewInt(11)
	proposalBlock := generateBlock(proposalHeight)
	clientAddr := currentCommittee[0].Address

	backendMock := NewMockBackend(ctrl)
	backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
	backendMock.EXPECT().LastCommittedProposal().AnyTimes().Return(currentBlock, currentCommittee[3].Address)
	backendMock.EXPECT().Committee(proposalHeight.Uint64()).AnyTimes().Return(committeeSet, nil)
	backendMock.EXPECT().SetProposedBlockHash(proposalBlock.Hash()).AnyTimes()
	backendMock.EXPECT().Broadcast(context.Background(), committeeSet, gomock.Any()).AnyTimes().Return(nil)
	backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)

	// create consensus core.
	c := New(backendMock)
	c.pendingUnminedBlocks[proposalHeight.Uint64()] = proposalBlock
	c.committeeSet = committeeSet
	c.sentProposal = false
	c.height = currentHeight
	round := int64(0)
	// since the default value of step and round is are both 0 which is to be one of the expected state, so we set them
	// into different value.
	c.step = precommitDone
	c.round = -1
	c.startRound(context.Background(), round)

	assert.Equal(t, c.sentProposal, true)
	assert.Equal(t, c.step, propose)
	assert.Equal(t, c.Round(), round)
}

// It test the page-6, line-11, StartRound() function from follower point of view of tendermint pseudo-code.
// Please refer to the algorithm from here: https://arxiv.org/pdf/1807.04938.pdf
func TestTendermintFollowerStartRound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currentCommittee := prepareCommittee()
	committeeSet, err := committee.NewSet(currentCommittee, currentCommittee[3].Address)
	if err != nil {
		t.Error(err)
	}

	currentHeight := big.NewInt(10)
	currentBlock := generateBlock(currentHeight)
	clientAddr := currentCommittee[0].Address

	backendMock := NewMockBackend(ctrl)
	backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
	backendMock.EXPECT().LastCommittedProposal().AnyTimes().Return(currentBlock, currentCommittee[3].Address)

	// create consensus core.
	c := New(backendMock)
	c.committeeSet = committeeSet
	round := int64(1)
	// since the default value of step and round is are both 0 which is to be one of the expected state, so we set them
	// into different value.
	c.step = precommitDone
	c.round = -1
	c.startRound(context.Background(), round)

	assert.Equal(t, c.proposeTimeout.started, true)
	assert.Equal(t, c.step, propose)
	assert.Equal(t, c.Round(), round)
}

// It test the upon proposal logic blocks from tendermint pseudo-code.
// Please refer to the algorithm from here: https://arxiv.org/pdf/1807.04938.pdf
func TestTendermintUponProposal(t *testing.T) {
	// Below 4 test cases cover line 22 to line 27 of tendermint pseudo-code.
	// It test line 24 was executed and step was forwarded on line 27.
	t.Run("on proposal with validRound as (-1) proposed and local lockedRound as (-1)", func(t *testing.T) {
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
		c.lockedRound = -1

		err = c.handleProposal(context.Background(), msg)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, c.sentPrevote, true)
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

		assert.Equal(t, c.sentPrevote, true)
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

		assert.Equal(t, c.sentPrevote, true)
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

		assert.Equal(t, c.sentPrevote, true)
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

		// condidtion vr >= 0 && vr < round_p, line 28.
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

		assert.Equal(t, c.sentPrevote, true)
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

		// condidtion vr >= 0 && vr < round_p, line 28.
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

		assert.Equal(t, c.sentPrevote, true)
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

		// condidtion vr >= 0 && vr < round_p, line 28.
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

		assert.Equal(t, c.sentPrevote, true)
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

		// condidtion vr >= 0 && vr < round_p, line 28.
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

		assert.Equal(t, c.sentPrevote, true)
		assert.Equal(t, c.step, prevote)
	})

}
