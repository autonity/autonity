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

// lines 34-35 & 61-64 of page 6 of The latest gossip on BFT consensus to describe the correct behaviours in
// test function for the current implementation of Tendermint.
func TestTendermintPrevoteTimeout(t *testing.T) {
		t.Run("Line34 to Line35, schedule for prevote timeout.", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		currentCommittee, _ := prepareCommittee()
		committeeSet, err := committee.NewSet(currentCommittee, currentCommittee[3].Address)
		if err != nil {
			t.Error(err)
		}

		currentHeight := big.NewInt(10)
		proposalBlock := generateBlock(currentHeight, types.BlockNonce{1, 2, 3, 4, 5, 6, 7, 8})
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
		c.validRound = roundProposed
		c.validValue = proposalBlock
		c.setValidRoundAndValue = true
		c.step = prevote

		// condition 2f+1 <PREVOTE, h_p, round_p, *>, power of pre-vote. line 34
		receivedPrevoteMsg := Message{
			Code:    msgPrevote,
			Address: currentCommittee[2].Address,
			power:   c.CommitteeSet().Quorum() - 1 ,
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

		backendMock.EXPECT().Post(gomock.Any()).Do(func(ev interface{}) {
			timeoutEvent, ok := ev.(TimeoutEvent)
			if !ok {
				t.Error("convert event failure.")
			}
			assert.Equal(t, timeoutEvent.roundWhenCalled, roundProposed)
			assert.Equal(t, timeoutEvent.heightWhenCalled.Uint64(), currentHeight.Uint64())
			assert.Equal(t, timeoutEvent.step, msgPrevote)
		})

		err = c.handlePrevote(context.Background(), preVoteMsg)
		if err != nil {
			t.Error(err)
		}
		// waif for timeout event.
		time.Sleep(timeoutPrevote(roundProposed)*2)
		// checking internal state of tendermint.
		assert.Equal(t, c.lockedRound, int64(-1))
		assert.Nil(t, c.lockedValue)
		assert.Equal(t, c.validRound, roundProposed)
		assert.Equal(t, c.validValue, proposalBlock)
		assert.Equal(t, c.step, prevote)
	})

	t.Run("Line 61 - 64, precommit for nil on timeout of prevote, the height does not change.", func(t *testing.T) {
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
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)

		roundProposed := int64(1)
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.round = roundProposed
		c.lockedRound = -1
		c.lockedValue = nil
		c.validRound = roundProposed
		c.validValue = proposalBlock
		c.setValidRoundAndValue = true
		c.step = prevote

		// prepare the wanted msg which will be broadcast.
		vote := Vote{
			Round:             roundProposed,
			Height:            currentHeight,
			ProposedBlockHash: common.Hash{},
		}
		preCommitMsg, err := Encode(&vote)
		if err != nil {
			t.Error("err")
		}
		wantedMsg, err := c.finalizeMessage(&Message{
			Code:          msgPrecommit,
			Msg:           preCommitMsg,
			Address:       clientAddr,
			CommittedSeal: []byte{},
		})
		if err != nil {
			t.Error(err)
		}
		// should check if broadcast to wanted committeeSet with wanted prevote msg.
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, wantedMsg).Return(nil)

		// timeout event.
		event := TimeoutEvent{
			roundWhenCalled:  roundProposed,
			heightWhenCalled: currentHeight,
			step:             msgPrevote,
		}
		c.handleTimeoutPrevote(context.Background(), event)

		// checking internal sate of tendermint.
		assert.Equal(t, c.round, roundProposed)
		assert.Equal(t, c.height, currentHeight)
		assert.Equal(t, c.lockedRound, int64(-1))
		assert.Nil(t, c.lockedValue)
		assert.Equal(t, c.validRound, roundProposed)
		assert.Equal(t, c.validValue, proposalBlock)
		// step from prevote to precommit with voting for nil.
		assert.Equal(t, c.step, precommit)
	})
}