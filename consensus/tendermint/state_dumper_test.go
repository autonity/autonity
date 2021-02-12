package tendermint

// import (
// 	"math/big"
// 	"math/rand"
// 	"testing"

// 	"github.com/clearmatics/autonity/common"
// 	"github.com/clearmatics/autonity/consensus/tendermint/config"
// 	"github.com/clearmatics/autonity/core/types"
// 	"github.com/clearmatics/autonity/crypto"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func TestGetLockedValueAndValidValue(t *testing.T) {
// 	c := &core{}
// 	b := generateBlock(new(big.Int).SetUint64(1))
// 	c.lockedValue = b
// 	c.validValue = b

// 	assert.Equal(t, c.lockedValue.Hash(), *getHash(c.lockedValue))
// 	assert.Equal(t, c.validValue.Hash(), *getHash(c.validValue))
// }

// func TestGetProposal(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	nodeAddr := common.BytesToAddress([]byte("node"))
// 	backendMock := NewMockBackend(ctrl)
// 	backendMock.EXPECT().Address().Return(nodeAddr)
// 	core := New(backendMock, config.RoundRobinConfig())

// 	proposalMsg, proposal := randomProposal(t)
// 	core.messages.getOrCreate(proposal.Round).SetProposal(&proposal, proposalMsg, true)

// 	assert.Equal(t, proposal.ProposalBlock.Hash(), *getProposal(core, proposal.Round))
// }

// func TestGetRoundState(t *testing.T) {
// 	sender := common.BytesToAddress([]byte("sender"))

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	backendMock := NewMockBackend(ctrl)
// 	backendMock.EXPECT().Address().Return(sender)

// 	c := New(backendMock, config.DefaultConfig())

// 	var rounds []int64 = []int64{0, 1}
// 	height := big.NewInt(int64(100) + 1)

// 	// Prepare 2 rounds of messages
// 	proposals := make([]Proposal, 2)
// 	proposals[0], _ = prepareRoundMsgs(t, c, rounds[0], height, sender)
// 	proposals[1], _ = prepareRoundMsgs(t, c, rounds[1], height, sender)

// 	// Get the states
// 	states := getRoundState(c)

// 	// expect 2 rounds of vote states.
// 	require.Len(t, states, 2)
// 	for _, state := range states {
// 		assert.Contains(t, rounds, state.Round)
// 		checkRoundState(t, state, rounds[state.Round], &proposals[state.Round], true)
// 	}
// }

// func TestGetCoreState(t *testing.T) {
// 	height := big.NewInt(int64(100) + 1)
// 	prevHeight := height.Sub(height, big.NewInt(1))
// 	prevBlock := generateBlock(prevHeight)
// 	knownMsgHash := []common.Hash{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5}, {0, 0, 1, 3, 6}}
// 	sender := common.BytesToAddress([]byte("sender"))

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	backendMock := NewMockBackend(ctrl)
// 	backendMock.EXPECT().Address().Return(sender)
// 	backendMock.EXPECT().KnownMsgHash().Return(knownMsgHash)

// 	c := New(backendMock, config.DefaultConfig())

// 	var rounds []int64 = []int64{0, 1}

// 	// Prepare 2 rounds of messages
// 	proposals := make([]Proposal, 2)
// 	proposers := make([]common.Address, 2)
// 	proposals[0], proposers[0] = prepareRoundMsgs(t, c, rounds[0], height, sender)
// 	proposals[1], proposers[1] = prepareRoundMsgs(t, c, rounds[1], height, sender)

// 	one := common.Big1
// 	members := []types.CommitteeMember{{Address: proposers[0], VotingPower: one}, {Address: proposers[1], VotingPower: one}}
// 	committeeSet, err := newRoundRobinSet(members, proposers[1]) // todo construct set here
// 	require.NoError(t, err)
// 	setCoreState(c, height, rounds[1], propose, proposals[0].ProposalBlock, rounds[0], proposals[0].ProposalBlock, rounds[0], committeeSet,
// 		prevBlock.Header())

// 	var e = coreStateRequestEvent{
// 		stateChan: make(chan TendermintState),
// 	}
// 	go c.handleStateDump(e)
// 	state := <-e.stateChan

// 	assert.Equal(t, sender, state.Client)
// 	assert.Equal(t, uint64(c.proposerPolicy), state.ProposerPolicy)
// 	assert.Equal(t, c.blockPeriod, state.BlockPeriod)
// 	assert.Len(t, state.CurHeightMessages, 6)
// 	assert.Equal(t, height, state.Height)
// 	assert.Equal(t, rounds[1], state.Round)
// 	assert.Equal(t, uint64(propose), state.Step)
// 	assert.Equal(t, proposals[1].ProposalBlock.Hash(), *state.Proposal)
// 	assert.Equal(t, proposals[0].ProposalBlock.Hash(), *state.LockedValue)
// 	assert.Equal(t, rounds[0], state.LockedRound)
// 	assert.Equal(t, proposals[0].ProposalBlock.Hash(), *state.ValidValue)
// 	assert.Equal(t, rounds[0], state.ValidRound)
// 	assert.Equal(t, committeeSet.Committee().String(), state.Committee.String())
// 	assert.Equal(t, committeeSet.GetProposer(rounds[1]).Address, state.Proposer)
// 	assert.False(t, state.IsProposer)
// 	assert.Equal(t, committeeSet.Quorum(), state.QuorumVotePower)
// 	assert.True(t, state.SentProposal)
// 	assert.True(t, state.SentPrevote)
// 	assert.True(t, state.SentPrecommit)
// 	assert.True(t, state.SetValidRoundAndValue)
// 	assert.False(t, state.ProposeTimerStarted)
// 	assert.False(t, state.PrevoteTimerStarted)
// 	assert.False(t, state.PrecommitTimerStarted)
// 	assert.Equal(t, knownMsgHash, state.KnownMsgHash)

// 	// expect 2 rounds of vote states.
// 	require.Len(t, state.RoundStates, 2)
// 	for _, s := range state.RoundStates {
// 		assert.Contains(t, rounds, s.Round)
// 		checkRoundState(t, s, rounds[s.Round], &proposals[s.Round], true)
// 	}
// }

// func randomProposal(t *testing.T) (*Message, Proposal) {
// 	currentHeight := big.NewInt(int64(rand.Intn(100) + 1))
// 	currentRound := int64(rand.Intn(100) + 1)

// 	proposer := common.BytesToAddress([]byte("proposer"))
// 	return generateBlockProposal(t, currentRound, currentHeight, currentRound-1, proposer, false)
// }

// func checkRoundState(t *testing.T, s RoundState, wantRound int64, wantProposal *Proposal, wantVerfied bool) {
// 	require.Equal(t, wantProposal.ProposalBlock.Hash(), s.Proposal)
// 	require.Len(t, s.PrevoteState, 1)
// 	require.Len(t, s.PrecommitState, 1)
// 	require.Equal(t, wantRound, s.Round)

// 	require.Equal(t, wantVerfied, s.PrevoteState[0].ProposalVerified)
// 	require.Equal(t, wantProposal.ProposalBlock.Hash(), s.PrevoteState[0].Value)

// 	require.Equal(t, wantVerfied, s.PrecommitState[0].ProposalVerified)
// 	require.Equal(t, wantProposal.ProposalBlock.Hash(), s.PrecommitState[0].Value)
// }

// func prepareRoundMsgs(t *testing.T, c *core, r int64, h *big.Int, sender common.Address) (proposal Proposal, proposer common.Address) {
// 	privKey, err := crypto.GenerateKey()
// 	require.NoError(t, err)
// 	proposalMsg, proposal := generateBlockProposal(t, r, h, 0, crypto.PubkeyToAddress(privKey.PublicKey), false)
// 	prevoteMsg, _, _ := prepareVote(t, msgPrevote, r, h, proposal.ProposalBlock.Hash(), sender, privKey)
// 	precommitMsg, _, _ := prepareVote(t, msgPrecommit, r, h, proposal.ProposalBlock.Hash(), sender, privKey)
// 	c.messages.getOrCreate(r).SetProposal(&proposal, proposalMsg, true)
// 	c.messages.getOrCreate(r).AddPrevote(proposal.ProposalBlock.Hash(), *prevoteMsg)
// 	c.messages.getOrCreate(r).AddPrecommit(proposal.ProposalBlock.Hash(), *precommitMsg)
// 	return proposal, proposalMsg.Address
// }

// func setCoreState(c *core, h *big.Int, r int64, s Step, lv *types.Block, lr int64, vv *types.Block, vr int64, committee committee, header *types.Header) {
// 	c.setHeight(h)
// 	c.setRound(r)
// 	c.setStep(s)
// 	c.lockedValue = lv
// 	c.lockedRound = lr
// 	c.validValue = vv
// 	c.validRound = vr
// 	c.setCommitteeSet(committee)
// 	c.lastHeader = header
// 	c.sentProposal = true
// 	c.sentPrevote = true
// 	c.sentPrecommit = true
// 	c.setValidRoundAndValue = true
// }
