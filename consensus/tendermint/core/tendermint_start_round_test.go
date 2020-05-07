package core

import (
	"context"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"github.com/clearmatics/autonity/core/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
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
