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
	c := &Core{roundsState: newTendermintState(log.New(), nil, nil)}
	b := generateBlock(new(big.Int).SetUint64(1))
	c.SetValidRoundAndValue(0, b)
	c.SetLockedRoundAndValue(0, b)

	assert.Equal(t, c.LockedValue().Hash(), *getHash(c.LockedValue()))
	assert.Equal(t, c.ValidValue().Hash(), *getHash(c.ValidValue()))
}

func TestGetProposal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodeAddr := common.BytesToAddress([]byte("node"))
	backendMock := interfaces.NewMockBackend(ctrl)
	core := New(backendMock, nil, nodeAddr, log.Root(), false, nil)
	core.roundsState = newTendermintState(log.New(), nil, nil)
	proposal := randomProposal(t)
	core.roundsState.GetOrCreate(proposal.R()).SetProposal(proposal, true)

	assert.Equal(t, proposal.Block().Hash(), *getProposal(core, proposal.R()))
}

func TestGetRoundState(t *testing.T) {
	sender := common.BytesToAddress([]byte("sender"))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	backendMock := interfaces.NewMockBackend(ctrl)
	c := New(backendMock, nil, sender, log.Root(), false, nil)
	c.roundsState = newTendermintState(log.New(), nil, nil)
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
	knownMsgHash := []common.Hash{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5}, {0, 0, 1, 3, 6}}
	sender := common.BytesToAddress([]byte("sender"))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	backendMock := interfaces.NewMockBackend(ctrl)
	backendMock.EXPECT().KnownMsgHash().Return(knownMsgHash)
	backendMock.EXPECT().FutureMsgs().Return(nil)
	c := New(backendMock, nil, sender, log.Root(), false, nil)
	c.roundsState = newTendermintState(log.New(), nil, nil)
	c.SetHeight(height)
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
	setCoreState(c, rounds[1], Propose, proposals[0].Block(), rounds[0], proposals[0].Block(), rounds[0], committeeSet)

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
	proposal := generateBlockProposal(r, h, 0, false, makeSigner(testConsensusKey), testCommitteeMember)
	prevoteMsg := message.NewPrevote(r, h.Uint64(), proposal.Block().Hash(), makeSigner(testConsensusKey), testCommitteeMember, 1)
	precommitMsg := message.NewPrecommit(r, h.Uint64(), proposal.Block().Hash(), makeSigner(testConsensusKey), testCommitteeMember, 1)
	c.roundsState.GetOrCreate(r).SetProposal(proposal, true)
	c.roundsState.GetOrCreate(r).AddPrevote(prevoteMsg)
	c.roundsState.GetOrCreate(r).AddPrecommit(precommitMsg)
	return proposal, proposal.Signer()
}

func setCoreState(c *Core, r int64, s Step, lv *types.Block, lr int64, vv *types.Block, vr int64, committee interfaces.Committee) {
	c.SetRound(r)
	c.SetStep(context.Background(), s)
	c.SetLockedRoundAndValue(lr, lv)
	c.SetValidRoundAndValue(vr, vv)

	c.setCommitteeSet(committee)
	c.SetSentProposal()
	c.SetSentPrevote()
	c.SetSentPrecommit()
	c.ValidRoundAndValueSet()
}
