package core

import (
	"context"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"github.com/clearmatics/autonity/core/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"math/rand"
	"sort"
	"strconv"
	"testing"
)

// TODO: We should create a utility function which can we used across different test files, it can be related to this
// issue https://github.com/clearmatics/autonity/issues/525
func prepareCommittee() types.Committee {
	// prepare committee.
	minSize := 4
	maxSize := 100
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

func TestTendermintStartRoundVariables(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	backendMock := NewMockBackend(ctrl)

	prevHeight := big.NewInt(int64(rand.Intn(100) + 1))
	prevBlock := generateBlock(prevHeight)
	currentHeight := big.NewInt(prevHeight.Int64() + 1)
	currentBlock := generateBlock(currentHeight)
	currentRound := int64(0)

	// We need to sort committee to have deterministic proposer for round 0. Since we have to say who was the last
	// proposer when calling committee.NewSet(), if the committee set passed to the new set is not sorted then the
	// index of last proposer will change and thus the index of proposer for round 0 will also change. Therefore
	// this test will be different every time is run as the client address can be the proposer for round 0 and a
	// different code path will be executed.
	currentCommittee := prepareCommittee()
	sort.Sort(currentCommittee)

	prevBlockProposerIndex := rand.Intn(len(currentCommittee))
	prevBlockProposer := currentCommittee[prevBlockProposerIndex].Address

	// We don't care who is the next proposer so for simplicity we ensure that clientAddress is not the next
	// proposer by setting clientAddress to be the address before prevBlockProposer in the ordered list. This will
	// ensure that the test will not run the broadcast method from backend, since the client will not be proposing
	// for round 0.
	clientAddress := currentCommittee[prevBlockProposerIndex-1%len(currentCommittee)].Address
	committeeSet, err := committee.NewSet(currentCommittee, prevBlockProposer)
	if err != nil {
		t.Error(err)
	}

	t.Run("ensure round 0 state variables are set correctly", func(t *testing.T) {
		backendMock.EXPECT().Address().Return(clientAddress)
		backendMock.EXPECT().LastCommittedProposal().Return(prevBlock, prevBlockProposer)
		backendMock.EXPECT().Committee(currentHeight.Uint64()).Return(committeeSet, nil)

		core := New(backendMock)
		core.startRound(context.Background(), currentRound)

		// Check the initial consensus state
		assert.Equal(t, core.Height(), currentHeight)
		assert.Equal(t, core.Round(), currentRound)
		assert.Equal(t, core.step, propose)
		assert.Nil(t, core.lockedValue)
		assert.Equal(t, core.lockedRound, int64(-1))
		assert.Nil(t, core.validValue)
		assert.Equal(t, core.validRound, int64(-1))
	})

	t.Run("ensure round x state variables are updated correctly", func(t *testing.T) {
		// In this test we are interested in making sure that that values which change in the current round that may
		// have an impact on the actions performed in the following round (in case of round change) are persisted
		// through to the subsequent round.
		backendMock.EXPECT().Address().Return(clientAddress)
		backendMock.EXPECT().LastCommittedProposal().Return(prevBlock, prevBlockProposer).MaxTimes(2)
		backendMock.EXPECT().Committee(currentHeight.Uint64()).Return(committeeSet, nil).MaxTimes(2)

		core := New(backendMock)
		core.startRound(context.Background(), currentRound)

		// Check the initial consensus state
		assert.Equal(t, core.Height(), currentHeight)
		assert.Equal(t, core.Round(), currentRound)
		assert.Equal(t, core.step, propose)
		assert.Nil(t, core.lockedValue)
		assert.Equal(t, core.lockedRound, int64(-1))
		assert.Nil(t, core.validValue)
		assert.Equal(t, core.validRound, int64(-1))

		// Update locked and valid Value (if locked value changes then valid value also changes, ie quorum(prevotes)
		// delivered in prevote step)
		core.lockedValue = currentBlock
		core.lockedRound = currentRound
		core.validValue = currentBlock
		core.validRound = currentRound

		// Move to next round anc check the expected state
		core.startRound(context.Background(), currentRound+1)

		// check consensus state
		assert.Equal(t, core.Height(), currentHeight)
		assert.Equal(t, core.Round(), currentRound+1)
		assert.Equal(t, core.step, propose)
		assert.Equal(t, core.lockedValue, currentBlock)
		assert.Equal(t, core.lockedRound, currentRound)
		assert.Equal(t, core.validValue, currentBlock)
		assert.Equal(t, core.validRound, currentRound)
	})
}

// It test the page-6, from Line-14 to Line 19, StartRound() function from proposer point of view of tendermint pseudo-code.
// Please refer to the algorithm from here: https://arxiv.org/pdf/1807.04938.pdf
func TestTendermintProposerStartRound(t *testing.T) {
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
	if err != nil {
		t.Error(err)
	}
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
	lastProposer := currentCommittee[len(currentCommittee)-1].Address
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
