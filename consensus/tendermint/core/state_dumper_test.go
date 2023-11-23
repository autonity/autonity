package core

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/autonity/autonity/consensus"
	tdmcommittee "github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	tctypes "github.com/autonity/autonity/consensus/tendermint/core/types"
	"github.com/autonity/autonity/log"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLockedValueAndValidValue(t *testing.T) {
	c := &Core{}
	b := generateBlock(new(big.Int).SetUint64(1))
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
	backendMock.EXPECT().Address().Return(nodeAddr)
	backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())
	core := New(backendMock, nil)

	proposalMsg, proposal := randomProposal(t)
	core.messages.GetOrCreate(proposal.Round).SetProposal(proposal, proposalMsg, true)

	assert.Equal(t, proposal.ProposalBlock.Hash(), *getProposal(core, proposal.Round))
}

func TestGetRoundState(t *testing.T) {
	sender := common.BytesToAddress([]byte("sender"))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	backendMock := interfaces.NewMockBackend(ctrl)
	backendMock.EXPECT().Address().Return(sender)
	backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

	c := New(backendMock, nil)

	rounds := []int64{0, 1}
	height := big.NewInt(int64(100) + 1)

	// Prepare 2 rounds of messages
	proposals := make([]*message.Proposal, 2)
	proposals[0], _ = prepareRoundMsgs(t, c, rounds[0], height, sender)
	proposals[1], _ = prepareRoundMsgs(t, c, rounds[1], height, sender)

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
	prevBlock := generateBlock(prevHeight)
	knownMsgHash := []common.Hash{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5}, {0, 0, 1, 3, 6}}
	sender := common.BytesToAddress([]byte("sender"))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	backendMock := interfaces.NewMockBackend(ctrl)
	backendMock.EXPECT().Address().Return(sender)
	backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())
	backendMock.EXPECT().KnownMsgHash().Return(knownMsgHash)

	c := New(backendMock, nil)

	var rounds = []int64{0, 1}

	// Prepare 2 rounds of messages
	proposals := make([]*message.Proposal, 2)
	proposers := make([]common.Address, 2)
	proposals[0], proposers[0] = prepareRoundMsgs(t, c, rounds[0], height, sender)
	proposals[1], proposers[1] = prepareRoundMsgs(t, c, rounds[1], height, sender)

	one := common.Big1
	committee := new(types.Committee)
	committee.Members = []*types.CommitteeMember{{Address: proposers[0], VotingPower: one}, {Address: proposers[1], VotingPower: one}}
	committeeSet, err := tdmcommittee.NewRoundRobinSet(committee, proposers[1])
	require.NoError(t, err)
	setCoreState(c, height, rounds[1], tctypes.Propose, proposals[0].ProposalBlock, rounds[0], proposals[0].ProposalBlock, rounds[0], committeeSet,
		prevBlock.Header())

	var e = tctypes.CoreStateRequestEvent{
		StateChan: make(chan tctypes.TendermintState),
	}
	go c.handleStateDump(e)
	state := <-e.StateChan
	assert.Equal(t, sender, state.Client)
	assert.Equal(t, c.blockPeriod, state.BlockPeriod)
	assert.Len(t, state.CurHeightMessages, 6)
	assert.Equal(t, height, state.Height)
	assert.Equal(t, rounds[1], state.Round)
	assert.Equal(t, uint64(tctypes.Propose), state.Step)
	assert.Equal(t, proposals[1].ProposalBlock.Hash(), *state.Proposal)
	assert.Equal(t, proposals[0].ProposalBlock.Hash(), *state.LockedValue)
	assert.Equal(t, rounds[0], state.LockedRound)
	assert.Equal(t, proposals[0].ProposalBlock.Hash(), *state.ValidValue)
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

func randomProposal(t *testing.T) (*message.Message, *message.Proposal) {
	currentHeight := big.NewInt(int64(rand.Intn(100) + 1))
	currentRound := int64(rand.Intn(100) + 1)

	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	proposer := crypto.PubkeyToAddress(key.PublicKey)

	return generateBlockProposal(t, currentRound, currentHeight, currentRound-1, proposer, false, key)
}

func checkRoundState(t *testing.T, s tctypes.RoundState, wantRound int64, wantProposal *message.Proposal, wantVerfied bool) {
	require.Equal(t, wantProposal.ProposalBlock.Hash(), s.Proposal)
	require.Len(t, s.PrevoteState, 1)
	require.Len(t, s.PrecommitState, 1)
	require.Equal(t, wantRound, s.Round)

	require.Equal(t, wantVerfied, s.PrevoteState[0].ProposalVerified)
	require.Equal(t, wantProposal.ProposalBlock.Hash(), s.PrevoteState[0].Value)

	require.Equal(t, wantVerfied, s.PrecommitState[0].ProposalVerified)
	require.Equal(t, wantProposal.ProposalBlock.Hash(), s.PrecommitState[0].Value)
}

func prepareRoundMsgs(t *testing.T, c *Core, r int64, h *big.Int, sender common.Address) (proposal *message.Proposal, proposer common.Address) {
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	proposalMsg, proposal := generateBlockProposal(t, r, h, 0, crypto.PubkeyToAddress(privKey.PublicKey), false, privKey)
	prevoteMsg, _, _ := prepareVote(t, consensus.MsgPrevote, r, h, proposal.ProposalBlock.Hash(), sender, privKey)
	precommitMsg, _, _ := prepareVote(t, consensus.MsgPrecommit, r, h, proposal.ProposalBlock.Hash(), sender, privKey)
	c.messages.GetOrCreate(r).SetProposal(proposal, proposalMsg, true)
	c.messages.GetOrCreate(r).AddPrevote(proposal.ProposalBlock.Hash(), *prevoteMsg)
	c.messages.GetOrCreate(r).AddPrecommit(proposal.ProposalBlock.Hash(), *precommitMsg)
	return proposal, proposalMsg.Address
}

func setCoreState(c *Core, h *big.Int, r int64, s tctypes.Step, lv *types.Block, lr int64, vv *types.Block, vr int64, committee interfaces.Committee, header *types.Header) {
	c.setHeight(h)
	c.setRound(r)
	c.SetStep(s)
	c.lockedValue = lv
	c.lockedRound = lr
	c.validValue = vv
	c.validRound = vr
	c.setCommittee(committee)
	c.setLastHeader(header)
	c.sentProposal = true
	c.sentPrevote = true
	c.sentPrecommit = true
	c.setValidRoundAndValue = true
}
