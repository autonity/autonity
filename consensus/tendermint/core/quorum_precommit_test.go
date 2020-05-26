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

// tests for what a node should do when it receives quorum precommits.
// line 49 - 54
func TestTendermintQuorumPrecommit(t *testing.T) {
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
		proposalBlock := generateBlock(currentHeight)
		nextProposalBlock := generateBlock(nextHeight)
		clientAddr := currentCommittee[0].Address
		validRoundProposed := int64(0)
		roundProposed := int64(0)

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)

		core := New(backendMock)
		initState(core, committeeSet, currentHeight, roundProposed, proposalBlock, roundProposed, proposalBlock, roundProposed, precommit, true)
		core.pendingUnminedBlocks[nextHeight.Uint64()] = nextProposalBlock

		// prepare round state with messages.
		// condition 2f+1 <PRECOMMIT, h_p, round_p, id(v)>, line 49
		memberCommitted, _ := committeeSet.GetByIndex(2)
		preCommitMsgReceived, err := preparePrecommitMsg(proposalBlock.Hash(), roundProposed, currentHeight.Int64(), keyMap, memberCommitted)
		if err != nil {
			t.Error(err)
		}
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
			Code:          msgPrecommit,
			Msg:           encodePreCommit,
			Address:       currentCommittee[2].Address,
			power:         core.CommitteeSet().Quorum() - 1,
			CommittedSeal: preCommitMsgReceived.CommittedSeal,
		}

		// init the proposal round states.
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

		// init the round states.
		core.curRoundMessages.SetProposal(proposal, proposalMsg, true)
		core.curRoundMessages.AddPrecommit(proposalBlock.Hash(), receivedPreCommitMsg)
		core.messages.getOrCreate(validRoundProposed).AddPrecommit(proposalBlock.Hash(), receivedPreCommitMsg)

		// last preCommitMsg to be handled.
		member, _ := committeeSet.GetByIndex(1)
		preCommitMsg, err := preparePrecommitMsg(proposalBlock.Hash(), roundProposed, currentHeight.Int64(), keyMap, member)
		if err != nil {
			t.Error(err)
		}

		// check the value and round to be commit to backend for view changing.
		seals := [][]byte{receivedPreCommitMsg.CommittedSeal, preCommitMsg.CommittedSeal}
		backendMock.EXPECT().Commit(proposalBlock, core.round, seals).Return(nil).AnyTimes()
		backendMock.EXPECT().LastCommittedProposal().Return(proposalBlock, clientAddr).AnyTimes()
		backendMock.EXPECT().Committee(gomock.Any()).Return(committeeSet, nil).AnyTimes()

		// check the behavior to be executed at start round after the commit.
		backendMock.EXPECT().SetProposedBlockHash(nextProposalBlock.Hash())
		backendMock.EXPECT().Sign(gomock.Any()).Return(nil, nil)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(),gomock.Any()).Return(nil)

		err = core.handlePrecommit(context.Background(), preCommitMsg)
		if err != nil {
			t.Error(err)
		}

		// It is hard to control tendermint's state machine if we construct the full backend since it overwrites the
		// state we simulated on this test context again and again. So we assume the CommitEvent is sent from miner/worker
		// thread via backend's interface, and it is handled to start new round here:
		core.handleCommit(context.Background())

		// checking tendermint internal states
		newHeight := new(big.Int).Add(currentHeight, common.Big1)
		checkState(t, core, newHeight, 0, nil, int64(-1), nil, int64(-1), propose)
		// round msg should be reset to empty set.
		assert.Equal(t, len(core.messages.GetMessages()), 0)
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