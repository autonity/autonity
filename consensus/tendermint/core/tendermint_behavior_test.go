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
	"strconv"
	"testing"
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
