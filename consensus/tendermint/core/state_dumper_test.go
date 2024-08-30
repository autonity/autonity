package core

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	tdmcommittee "github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/log"
)

func TestGetLockedValueAndValidValue(t *testing.T) {
	c := &Core{}
	lastHeader := &types.Header{Number: big.NewInt(0)}
	b := generateBlock(new(big.Int).SetUint64(1), lastHeader)
	c.lockedValue = b
	c.validValue = b

	assert.Equal(t, c.lockedValue.Hash(), *getHash(c.lockedValue))
	assert.Equal(t, c.validValue.Hash(), *getHash(c.validValue))
}

func TestGetProposal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodeAddr := common.BytesToAddress([]byte("node"))
	backendMock := interfaces.NewMockBackend(ctrl)
	core := New(backendMock, nil, nodeAddr, log.Root(), false)

	proposal := randomProposal(t)
	core.messages.GetOrCreate(proposal.R()).SetProposal(proposal, true)

	assert.Equal(t, proposal.Block().Hash(), *getProposal(core, proposal.R()))
}

func TestGetRoundState(t *testing.T) {
	sender := common.BytesToAddress([]byte("sender"))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	backendMock := interfaces.NewMockBackend(ctrl)
	c := New(backendMock, nil, sender, log.Root(), false)
	rounds := []int64{0, 1}
	height := big.NewInt(int64(100) + 1)

	// Prepare 2 rounds of messages
	proposals := make([]*message.Propose, 2)
	proposals[0], _ = prepareRoundMsgs(c, rounds[0], height)
	proposals[1], _ = prepareRoundMsgs(c, rounds[1], height)

	// Get the states
	states := getRoundState(c)

	// expect 2 rounds of vote states.
	require.Len(t, states, 2)
	for _, state := range states {
		assert.Contains(t, rounds, state.Round)
		checkRoundState(t, state, rounds[state.Round], proposals[state.Round], true)
	}
}

func TestGetCoreState(t *testing.T) {
	height := big.NewInt(int64(100) + 1)
	prevHeight := height.Sub(height, big.NewInt(1))
	lastHeader := &types.Header{Number: big.NewInt(0)}
	prevBlock := generateBlock(prevHeight, lastHeader)
	knownMsgHash := []common.Hash{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5}, {0, 0, 1, 3, 6}}
	sender := common.BytesToAddress([]byte("sender"))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	backendMock := interfaces.NewMockBackend(ctrl)
	backendMock.EXPECT().KnownMsgHash().Return(knownMsgHash)
	backendMock.EXPECT().FutureMsgs().Return(nil)
	c := New(backendMock, nil, sender, log.Root(), false)

	var rounds = []int64{0, 1}

	// Prepare 2 rounds of messages
	proposals := make([]*message.Propose, 2)
	proposers := make([]common.Address, 2)
	proposals[0], proposers[0] = prepareRoundMsgs(c, rounds[0], height)
	proposals[1], proposers[1] = prepareRoundMsgs(c, rounds[1], height)

	blsKeys := make([]blst.PublicKey, 2)
	blsKey, err := blst.RandKey()
	require.NoError(t, err)
	blsKeys[0] = blsKey.PublicKey()
	blsKey, err = blst.RandKey()
	require.NoError(t, err)
	blsKeys[1] = blsKey.PublicKey()

	one := common.Big1
	committee := new(types.Committee)
	committee.Members = []types.CommitteeMember{
		{Address: proposers[0], VotingPower: one, ConsensusKey: blsKeys[1], ConsensusKeyBytes: blsKeys[1].Marshal(), Index: 0},
		{Address: proposers[1], VotingPower: one, ConsensusKey: blsKeys[0], ConsensusKeyBytes: blsKeys[0].Marshal(), Index: 1},
	}
	committeeSet, err := tdmcommittee.NewRoundRobinSet(committee, proposers[1]) // todo construct set here
	require.NoError(t, err)
	setCoreState(c, height, rounds[1], Propose, proposals[0].Block(), rounds[0], proposals[0].Block(), rounds[0], committeeSet,
		prevBlock.Header())

	var e = StateRequestEvent{
		StateChan: make(chan interfaces.CoreState),
	}
	go c.handleStateDump(e)
	state := <-e.StateChan
	assert.Equal(t, sender, state.Client)
	assert.Equal(t, c.blockPeriod, state.BlockPeriod)
	assert.Len(t, state.CurHeightMessages, 6)
	assert.Equal(t, height, state.Height)
	assert.Equal(t, rounds[1], state.Round)
	assert.Equal(t, uint64(Propose), state.Step)
	assert.Equal(t, proposals[1].Value(), *state.Proposal)
	assert.Equal(t, proposals[0].Value(), *state.LockedValue)
	assert.Equal(t, rounds[0], state.LockedRound)
	assert.Equal(t, proposals[0].Value(), *state.ValidValue)
	assert.Equal(t, rounds[0], state.ValidRound)
	assert.Equal(t, committeeSet.Committee().String(), state.Committee.String())
	assert.Equal(t, committeeSet.GetProposer(rounds[1]).Address, state.Proposer)
	assert.False(t, state.IsProposer)
	assert.Equal(t, committeeSet.Quorum(), state.QuorumVotePower)
	assert.True(t, state.SentProposal)
	assert.True(t, state.SentPrevote)
	assert.True(t, state.SentPrecommit)
	assert.True(t, state.SetValidRoundAndValue)
	assert.False(t, state.ProposeTimerStarted)
	assert.False(t, state.PrevoteTimerStarted)
	assert.False(t, state.PrecommitTimerStarted)
	assert.Equal(t, knownMsgHash, state.KnownMsgHash)

	// expect 2 rounds of vote states.
	require.Len(t, state.RoundStates, 2)
	for _, s := range state.RoundStates {
		assert.Contains(t, rounds, s.Round)
		checkRoundState(t, s, rounds[s.Round], proposals[s.Round], true)
	}
}

func checkRoundState(t *testing.T, s interfaces.RoundState, wantRound int64, wantProposal *message.Propose, wantVerfied bool) {
	require.Equal(t, wantProposal.Block().Hash(), s.Proposal)
	require.Len(t, s.PrevoteState, 1)
	require.Len(t, s.PrecommitState, 1)
	require.Equal(t, wantRound, s.Round)

	require.Equal(t, wantVerfied, s.PrevoteState[0].ProposalVerified)
	require.Equal(t, wantProposal.Block().Hash(), s.PrevoteState[0].Value)

	require.Equal(t, wantVerfied, s.PrecommitState[0].ProposalVerified)
	require.Equal(t, wantProposal.Block().Hash(), s.PrecommitState[0].Value)
}

func prepareRoundMsgs(c *Core, r int64, h *big.Int) (*message.Propose, common.Address) {

	lastHeader := &types.Header{Number: big.NewInt(h.Int64()).Sub(h, common.Big1)}
	proposal := generateBlockProposal(r, h, 0, false, makeSigner(testConsensusKey), testCommitteeMember, lastHeader)
	prevoteMsg := message.NewPrevote(r, h.Uint64(), proposal.Block().Hash(), makeSigner(testConsensusKey), testCommitteeMember, 1)
	precommitMsg := message.NewPrecommit(r, h.Uint64(), proposal.Block().Hash(), makeSigner(testConsensusKey), testCommitteeMember, 1)
	c.messages.GetOrCreate(r).SetProposal(proposal, true)
	c.messages.GetOrCreate(r).AddPrevote(prevoteMsg)
	c.messages.GetOrCreate(r).AddPrecommit(precommitMsg)
	return proposal, proposal.Signer()
}

func setCoreState(c *Core, h *big.Int, r int64, s Step, lv *types.Block, lr int64, vv *types.Block, vr int64, committee interfaces.Committee, header *types.Header) {
	c.setHeight(h)
	c.setRound(r)
	c.SetStep(context.Background(), s)
	c.lockedValue = lv
	c.lockedRound = lr
	c.validValue = vv
	c.validRound = vr
	c.setCommitteeSet(committee)
	c.sentProposal = true
	c.sentPrevote = true
	c.sentPrecommit = true
	c.setValidRoundAndValue = true
}
