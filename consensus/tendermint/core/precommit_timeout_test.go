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
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address
		validRoundProposed := int64(0)
		roundProposed := int64(0)

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)

		// init tendermint internal states.
		core := New(backendMock)
		initState(core, committeeSet, currentHeight, roundProposed, proposalBlock, roundProposed, proposalBlock, roundProposed, precommit, true)

		// condition 2f+1 <PRECOMMIT, h_p, round_p, *>, line 47. For the round state preparation.
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
			Msg:     encodePreCommit,
			Address: currentCommittee[2].Address,
			power:   core.CommitteeSet().Quorum() - 1,
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

		// init round states.
		core.curRoundMessages.SetProposal(proposal, proposalMsg, true)
		core.curRoundMessages.AddPrecommit(proposalBlock.Hash(), receivedPreCommitMsg)
		core.messages.getOrCreate(validRoundProposed).AddPrecommit(proposalBlock.Hash(), receivedPreCommitMsg)

		// new precommit msg which vote for different value.
		member, _ := committeeSet.GetByIndex(1)
		preCommitMsg, err := preparePrecommitMsg(generateBlock(currentHeight).Hash(), roundProposed, currentHeight.Int64(), keyMap, member)
		checkTimeOutEvent(t, backendMock, msgPrecommit, currentHeight, roundProposed)
		err = core.handleCheckedMsg(context.Background(), preCommitMsg, member)
		if err != nil {
			t.Error(err)
		}

		// waif for timeout event.
		time.Sleep(timeoutPrecommit(roundProposed) * 2)
		// checking internal state of tendermint.
		checkState(t, core, currentHeight, roundProposed, proposalBlock, roundProposed, proposalBlock, roundProposed, precommit)
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
		proposalBlock := generateBlock(currentHeight)
		clientAddr := currentCommittee[0].Address
		roundProposed := int64(1)

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)

		//init tendermint internal states.
		core := New(backendMock)
		initState(core, committeeSet, currentHeight, roundProposed, proposalBlock, roundProposed, proposalBlock, roundProposed, precommit, true)

		event := TimeoutEvent{
			roundWhenCalled:  roundProposed,
			heightWhenCalled: currentHeight,
			step:             msgPrecommit,
		}
		core.handleTimeoutPrecommit(context.Background(), event)
		checkState(t, core, currentHeight, roundProposed + 1, proposalBlock, roundProposed, proposalBlock, roundProposed, propose)
	})
}

func checkTimeOutEvent(t *testing.T, backendMock *MockBackend, msgCode uint64, height *big.Int, round int64) {
	backendMock.EXPECT().Post(gomock.Any()).Do(func(ev interface{}) {
		timeoutEvent, ok := ev.(TimeoutEvent)
		if !ok {
			t.Error("convert event failure.")
		}
		assert.Equal(t, round, timeoutEvent.roundWhenCalled)
		assert.Equal(t, height.Uint64(), timeoutEvent.heightWhenCalled.Uint64())
		assert.Equal(t, msgCode, timeoutEvent.step)
	})
}

func initState(core *core, committee *committee.Set, height *big.Int, round int64, lockedValue *types.Block, lockedRound int64, validValue *types.Block, validRound int64, step Step, setValid bool) {
	core.committeeSet = committee
	core.height = height
	core.round = round
	core.lockedRound = lockedRound
	core.lockedValue = lockedValue
	core.validRound = validRound
	core.validValue = validValue
	core.setValidRoundAndValue = setValid
	core.step = step
}

// checking internal state of tendermint.
func checkState(t *testing.T, core *core, height *big.Int, round int64, lockedValue *types.Block, lockedRound int64, validValue *types.Block, validRound int64, step Step) {
	assert.Equal(t, height, core.Height())
	assert.Equal(t, round, core.Round())
	assert.Equal(t, validValue, core.validValue)
	assert.Equal(t, validRound, core.validRound)
	assert.Equal(t, lockedValue, core.lockedValue)
	assert.Equal(t, lockedRound, core.lockedRound)
	assert.Equal(t, step, core.step)
}

func prepareCommittee() (types.Committee, addressKeyMap) {
	// prepare committee.
	minSize := 4
	maxSize := 15
	committeeSize := rand.Intn(maxSize-minSize) + minSize
	committeeSet, keyMap := generateCommittee(committeeSize)
	return committeeSet, keyMap
}

func generateBlock(height *big.Int) *types.Block {
	// use random nonce to create different blocks
	var nonce types.BlockNonce
	for i := 0; i < len(nonce); i++ {
		nonce[0] = byte(rand.Intn(256))
	}
	header := &types.Header{Number: height, Nonce: nonce}
	block := types.NewBlock(header, nil, nil, nil)
	return block
}