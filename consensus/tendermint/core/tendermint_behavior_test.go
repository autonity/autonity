package core

import (
	"context"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"github.com/clearmatics/autonity/core/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"math/rand"
	"testing"
	"time"
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

// tests for when a node should schedule precommit timeout and what should be done when it expires.
// line 47 - 48 & 65-67
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
		proposalBlock := generateBlock(currentHeight, types.BlockNonce{1, 2, 3, 4, 5, 6, 7, 8})
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

		// condition 2f+1 <PRECOMMIT, h_p, round_p, *>, line 47
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
			Code:    msgPrecommit,
			Msg: 	 encodePreCommit,
			Address: currentCommittee[2].Address,
			power:   c.CommitteeSet().Quorum() - 1,
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

		member, _ := committeeSet.GetByIndex(1)
		preCommitMsg, err := preparePrecommitMsg(generateBlock(currentHeight, types.BlockNonce{1, 2, 3, 4, 5, 6, 7, 9}).Hash(), roundProposed, currentHeight.Int64(), keyMap, member)
		backendMock.EXPECT().Post(gomock.Any()).Do(func(ev interface{}) {
			timeoutEvent, ok := ev.(TimeoutEvent)
			if !ok {
				t.Error("convert event failure.")
			}
			assert.Equal(t, timeoutEvent.roundWhenCalled, roundProposed)
			assert.Equal(t, timeoutEvent.heightWhenCalled.Uint64(), currentHeight.Uint64())
			assert.Equal(t, timeoutEvent.step, msgPrecommit)
		})

		err = c.handlePrecommit(context.Background(), preCommitMsg)
		if err != nil {
			t.Error(err)
		}

		// waif for timeout event.
		time.Sleep(timeoutPrecommit(roundProposed)*2)
		// checking internal state of tendermint.
		assert.Equal(t, c.lockedRound, roundProposed)
		assert.Equal(t, c.lockedValue, proposalBlock)
		assert.Equal(t, c.validRound, roundProposed)
		assert.Equal(t, c.validValue, proposalBlock)
	})

	t.Run("Line 65, 66, 67, start round on timeout of precommit, the height does not change.", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// prepare a random size of committee, and the proposer at last committed block.
		currentCommittee, _ := prepareCommittee()
		lastProposer := currentCommittee[len(currentCommittee)-1].Address
		committeeSet, err := committee.NewSet(currentCommittee, lastProposer)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight, types.BlockNonce{1, 2, 3, 4, 5, 6, 7, 8})
		clientAddr := currentCommittee[0].Address

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)

		roundProposed := int64(1)
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

		event := TimeoutEvent{
			roundWhenCalled:  roundProposed,
			heightWhenCalled: currentHeight,
			step:             msgPrecommit,
		}
		c.handleTimeoutPrecommit(context.Background(), event)
		assert.Equal(t, c.round, roundProposed + 1)
		assert.Equal(t, c.height, currentHeight)
		assert.Equal(t, c.lockedRound, roundProposed)
		assert.Equal(t, c.lockedValue, proposalBlock)
		assert.Equal(t, c.validRound, roundProposed)
		assert.Equal(t, c.validValue, proposalBlock)
		assert.Equal(t, c.step, propose)
	})
}