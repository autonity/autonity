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
		proposalBlock := generateBlock(currentHeight)
		validRoundProposed := int64(0)
		roundProposed := int64(0)
		clientAddr := currentCommittee[0].Address

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)

		// create the triggering msg.
		preVoteMsg := createPrevote(t, proposalBlock.Hash(), roundProposed, currentHeight, committeeSet.Committee()[0])

		// init tendermint internal states.
		core := New(backendMock)
		initState(core, committeeSet, currentHeight, roundProposed, nil, int64(-1), proposalBlock, roundProposed, prevote, true)

		// condition 2f+1 <PREVOTE, h_p, round_p, *>, quorum power of pre-vote for *. line 34
		// prepare the received prevote msg, and init the round state.
		receivedPrevoteMsg := Message{
			Code:    msgPrevote,
			Address: currentCommittee[2].Address,
			// the full quorum will be triggered by input msg by counting 1 to below power.
			power: core.CommitteeSet().Quorum() - 1,
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

		core.curRoundMessages.SetProposal(proposal, proposalMsg, true)
		core.curRoundMessages.AddPrevote(proposalBlock.Hash(), receivedPrevoteMsg)
		core.messages.getOrCreate(validRoundProposed).AddPrevote(proposalBlock.Hash(), receivedPrevoteMsg)

		// subscribe to check the prevote timeout event is sent to backend at the specific view.
		backendMock.EXPECT().Post(gomock.Any()).Do(func(ev interface{}) {
			timeoutEvent, ok := ev.(TimeoutEvent)
			if !ok {
				t.Error("convert event failure.")
			}
			assert.Equal(t, roundProposed, timeoutEvent.roundWhenCalled)
			assert.Equal(t, currentHeight.Uint64(), timeoutEvent.heightWhenCalled.Uint64())
			assert.Equal(t, msgPrevote, timeoutEvent.step)
		})

		_, msgSender, _ := core.committeeSet.GetByAddress(clientAddr)
		err = core.handleCheckedMsg(context.Background(), preVoteMsg, msgSender)
		if err != nil {
			t.Error(err)
		}
		// wait for a duration and check the state changes was triggered handlePrevote.
		duration := timeoutPrevote(roundProposed) * 2
		time.Sleep(duration)
		// checking internal state of tendermint.
		checkState(t, core, currentHeight, roundProposed, nil, int64(-1), proposalBlock, roundProposed, prevote)
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
		roundProposed := int64(1)
		proposalBlock := generateBlock(currentHeight)
		lockedRound := int64(-1)
		step := prevote
		clientAddr := currentCommittee[0].Address

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)

		core := New(backendMock)
		initState(core, committeeSet, currentHeight, roundProposed, nil, lockedRound, proposalBlock, roundProposed, step, true)

		// subscribe for checking the broadcast msg: precommit for nil on the specific view.
		checkBroadcastMsg(t, backendMock, core, msgPrecommit, currentHeight, roundProposed, common.Hash{})

		// prepare timeout event.
		event := TimeoutEvent{
			roundWhenCalled:  roundProposed,
			heightWhenCalled: currentHeight,
			step:             msgPrevote,
		}
		core.handleTimeoutPrevote(context.Background(), event)

		// checking internal state of tendermint.
		checkState(t, core, currentHeight, roundProposed, nil, int64(-1), proposalBlock, roundProposed, precommit)
	})
}

func checkBroadcastMsg(t *testing.T, backendMock *MockBackend, core *core, msgCode uint64, height *big.Int, round int64, hash common.Hash) {
	// prepare the wanted msg which will be broadcast.
	vote := Vote{
		Round:             round,
		Height:            height,
		ProposedBlockHash: hash,
	}
	msg, err := Encode(&vote)
	if err != nil {
		t.Error(err)
	}
	wantedMsg, err := core.finalizeMessage(&Message{
		Code:          msgCode,
		Msg:           msg,
		Address:       core.address,
		CommittedSeal: []byte{},
	})
	if err != nil {
		t.Error(err)
	}
	// should check if broadcast to wanted committeeSet with wanted msg.
	backendMock.EXPECT().Broadcast(context.Background(), core.committeeSet, wantedMsg).Return(nil)
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
