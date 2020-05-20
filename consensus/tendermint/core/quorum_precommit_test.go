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
	"testing"
)

func prepareCommittee() (types.Committee, addressKeyMap) {
	// prepare committee.
	minSize := 4
	maxSize := 15
	committeeSize := rand.Intn(maxSize - minSize) + minSize
	committeeSet, keyMap := generateCommittee(committeeSize)
	return committeeSet, keyMap
}

func generateBlock(height *big.Int, nonce types.BlockNonce) *types.Block {
	header := &types.Header{Number: height, Nonce: nonce}
	block := types.NewBlock(header, nil, nil, nil)
	return block
}

// tests for what a node should do when it receives quorum precommits.
// line 49 - 54
func TestTendermintPrecommitTimeout(t *testing.T) {
	t.Run("Line 47 to Line48, schedule for precommit timeout.", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// prepare a random size of committee, and the proposer at last committed block.
		currentCommittee, keyMap := prepareCommittee()
		lastProposer := currentCommittee[len(currentCommittee)-1].Address
		committeeSet, err := committee.NewSet(currentCommittee, lastProposer)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		nextHeight := big.NewInt(11)
		proposalBlock := generateBlock(currentHeight, types.BlockNonce{1, 2, 3, 4, 5, 6, 7, 8})
		nextProposalBlock := generateBlock(nextHeight, types.BlockNonce{1, 2, 3, 4, 5, 6, 7, 9})
		clientAddr := currentCommittee[0].Address

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)

		validRoundProposed := int64(0)
		roundProposed := int64(0)
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.round = roundProposed
		c.lockedRound = roundProposed
		c.lockedValue = proposalBlock
		c.validRound = roundProposed
		c.validValue = proposalBlock
		c.setValidRoundAndValue = true
		c.step = precommit
		c.pendingUnminedBlocks[nextHeight.Uint64()] = nextProposalBlock

		member, _ := committeeSet.GetByIndex(1)
		preCommitMsg, err := preparePrecommitMsg(proposalBlock.Hash(), roundProposed, currentHeight.Int64(), keyMap, member)

		// condition 2f+1 <PRECOMMIT, h_p, round_p, id(v)>, line 49
		var preCommit = Vote{
			Round:             roundProposed,
			Height:            currentHeight,
			ProposedBlockHash: proposalBlock.Hash(),
		}

		encodePreCommit, err := Encode(&preCommit)
		if err != nil {
			t.Error(err)
		}

		receivedPreCommitMsg := Message{
			Code:    		msgPrecommit,
			Msg: 	 		encodePreCommit,
			Address: 		currentCommittee[2].Address,
			power:   		c.CommitteeSet().Quorum() - 1,
			CommittedSeal: 	preCommitMsg.CommittedSeal,
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
		c.curRoundMessages.AddPrecommit(proposalBlock.Hash(), receivedPreCommitMsg)
		c.messages.getOrCreate(validRoundProposed).AddPrecommit(proposalBlock.Hash(), receivedPreCommitMsg)

		// check the value and round commit to the backend for view changing.
		seals := [][]byte{receivedPreCommitMsg.CommittedSeal, preCommitMsg.CommittedSeal}
		backendMock.EXPECT().Commit(proposalBlock, c.round, seals).Return(nil).AnyTimes()
		backendMock.EXPECT().LastCommittedProposal().Return(proposalBlock, clientAddr).AnyTimes()
		backendMock.EXPECT().Committee(gomock.Any()).Return(committeeSet, nil).AnyTimes()
		// behavior to be executed at start round after the commit.
		backendMock.EXPECT().SetProposedBlockHash(nextProposalBlock.Hash())
		backendMock.EXPECT().Sign(gomock.Any()).Return(nil, nil)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(),gomock.Any()).Return(nil)

		err = c.handlePrecommit(context.Background(), preCommitMsg)
		if err != nil {
			t.Error(err)
		}

		// assume that the CommitEvent is sent from backend, and it is handled by handleCommit to start new round.
		c.handleCommit(context.Background())

		// checking tendermint internal states
		assert.Equal(t, c.step, propose)
		assert.Equal(t, c.height, new(big.Int).Add(currentHeight, common.Big1))
		assert.Nil(t, c.lockedValue)
		assert.Nil(t, c.validValue)
		assert.Equal(t, c.lockedRound, int64(-1))
		assert.Nil(t, c.validValue)
		// round msg should be reset to empty set.
		assert.Equal(t, len(c.messages.GetMessages()), 0)
	})
}