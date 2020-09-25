package core

import (
	"crypto/ecdsa"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/core/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func TestStateDumper_GetProposal(t *testing.T) {
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet, privateKeys := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address

	genRoundVoteMessages := func(c *core, round int64, proposal *Proposal, proposalMsg *Message, preVoteMsg *Message, preCommitMsg *Message, verified bool) {
		c.messages.getOrCreate(round).SetProposal(proposal, proposalMsg, verified)
		c.messages.getOrCreate(round).AddPrevote(proposal.ProposalBlock.Hash(), *preVoteMsg)
		c.messages.getOrCreate(round).AddPrecommit(proposal.ProposalBlock.Hash(), *preCommitMsg)
	}

	checkRoundState := func(t *testing.T, s types.RoundState, wantRound int64, wantProposal *Proposal, wantVerfied bool) {
		require.Equal(t, wantProposal.ProposalBlock.Hash(), s.Proposal)
		require.Len(t, s.PrevoteState, 1)
		require.Len(t, s.PrecommitState, 1)
		require.Equal(t, wantRound, s.Round)

		require.Equal(t, wantVerfied, s.PrevoteState[0].ProposalVerified)
		require.Equal(t, wantProposal.ProposalBlock.Hash(), s.PrevoteState[0].Value)

		require.Equal(t, wantVerfied, s.PrecommitState[0].ProposalVerified)
		require.Equal(t, wantProposal.ProposalBlock.Hash(), s.PrecommitState[0].Value)
	}

	prepareRoundMsgs := func(t *testing.T, r int64, h *big.Int, vr int64, proposer common.Address, sender common.Address,
		privKey *ecdsa.PrivateKey) (*Message, Proposal, *Message, *Message) {
		proposalMsg, proposal := generateBlockProposal(t, r, h, vr, proposer, false)
		prevoteMsg, _, _ := prepareVote(t, msgPrevote, r, h, proposal.ProposalBlock.Hash(), sender, privKey)
		precommitMsg, _, _ := prepareVote(t, msgPrecommit, r, h, proposal.ProposalBlock.Hash(), sender, privKey)
		return proposalMsg, proposal, prevoteMsg, precommitMsg
	}

	setCoreState := func(c *core, h *big.Int, r int64, s Step, lv *types.Block, lr int64, vv *types.Block, vr int64, committee committee, header *types.Header) {
		c.setHeight(h)
		c.setRound(r)
		c.setStep(s)
		c.lockedValue = lv
		c.lockedRound = lr
		c.validValue = vv
		c.validRound = vr
		c.setCommitteeSet(committee)
		c.lastHeader = header
		c.sentProposal = true
		c.sentPrevote = true
		c.sentPrecommit = true
		c.setValidRoundAndValue = true
	}

	checkState := func(t *testing.T, c *core, state types.TendermintState, currentHeight *big.Int, currentRound int64, initRound int64,
		initProposal Proposal, newProposal Proposal, prevBlock *types.Block, knownMsgHash []common.Hash) {

		require.Equal(t, int64(0), state.Code)
		require.Equal(t, clientAddr, state.Client)
		require.Equal(t, uint64(c.proposerPolicy), state.ProposerPolicy)
		require.Equal(t, c.blockPeriod, state.BlockPeriod)
		require.Len(t, state.CurHeightMessages, 6)
		require.Equal(t, *currentHeight, state.Height)
		require.Equal(t, currentRound, state.Round)
		require.Equal(t, uint64(propose), state.Step)
		require.Equal(t, newProposal.ProposalBlock.Hash(), state.Proposal)
		require.Equal(t, initProposal.ProposalBlock.Hash(), state.LockedValue)
		require.Equal(t, initRound, state.LockedRound)
		require.Equal(t, initProposal.ProposalBlock.Hash(), state.ValidValue)
		require.Equal(t, initRound, state.ValidRound)
		require.Equal(t, c.getParentCommittee().String(), prevBlock.Header().Committee.String())
		require.Equal(t, committeeSet.Committee().String(), state.Committee.String())
		require.Equal(t, members[currentRound].Address, state.Proposer)
		require.False(t, state.IsProposer)
		require.Equal(t, committeeSet.Quorum(), state.QuorumVotePower)
		require.True(t, state.SentProposal)
		require.True(t, state.SentPrevote)
		require.True(t, state.SentPrecommit)
		require.True(t, state.SetValidRoundAndValue)
		require.False(t, state.ProposeTimerStarted)
		require.False(t, state.PrevoteTimerStarted)
		require.False(t, state.PrecommitTimerStarted)
		require.Equal(t, knownMsgHash, state.KnownMsgHash)

		// expect 2 rounds of vote states.
		require.Len(t, state.RoundStates, 2)
		for _, v := range state.RoundStates {
			require.Contains(t, []int64{initRound, currentRound}, v.Round)
			if v.Round == currentRound {
				checkRoundState(t, v, currentRound, &newProposal, true)
			}
			if v.Round == initRound {
				checkRoundState(t, v, initRound, &initProposal, true)
			}
		}
	}

	t.Run("get proposal, locked value and valid value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		core := New(backendMock, config.RoundRobinConfig())
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, int64(rand.Intn(int(currentRound+1)-1)), members[currentRound].Address, false)
		core.messages.getOrCreate(currentRound).SetProposal(&proposal, proposalMsg, true)
		fact := core.getProposal(currentRound)
		core.lockedValue = proposal.ProposalBlock
		core.validValue = proposal.ProposalBlock
		assert.Equal(t, proposal.ProposalBlock.Hash(), fact)
		assert.Equal(t, core.getLockedValue(), core.lockedValue.Hash())
		assert.Equal(t, core.getValidValue(), core.validValue.Hash())
	})

	t.Run("get parent block committee", func(t *testing.T) {
		prevHeight := big.NewInt(int64(rand.Intn(100) + 1))
		prevBlock := generateBlock(prevHeight)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		core := New(backendMock, config.DefaultConfig())
		core.setCommitteeSet(committeeSet)
		core.lastHeader = prevBlock.Header()

		assert.Equal(t, core.getParentCommittee().String(), prevBlock.Header().Committee.String())
	})

	t.Run("get round state", func(t *testing.T) {

		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		// round 0 messages
		initRound := int64(0)
		initProposalMsg, initProposal, initRoundPrevoteMsg, initRoundPrecommitMsg := prepareRoundMsgs(t, initRound, currentHeight, 0, members[initRound].Address, clientAddr, privateKeys[clientAddr])

		// round 1 messages
		currentRound := int64(1)
		sender := 1
		newProposalMsg, newProposal, curRoundPrevoteMsg, curRoundPrecommitMsg := prepareRoundMsgs(t, currentRound, currentHeight, 0, members[currentRound].Address, members[sender].Address, privateKeys[members[sender].Address])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.DefaultConfig())

		// round 0 messages:
		genRoundVoteMessages(c, initRound, &initProposal, initProposalMsg, initRoundPrevoteMsg, initRoundPrecommitMsg, true)
		// current round messages:
		genRoundVoteMessages(c, currentRound, &newProposal, newProposalMsg, curRoundPrevoteMsg, curRoundPrecommitMsg, true)

		states := c.getRoundState()

		// expect 2 rounds of vote states.
		require.Len(t, states, 2)
		for _, v := range states {
			require.Contains(t, []int64{initRound, currentRound}, v.Round)
			if v.Round == currentRound {
				checkRoundState(t, v, currentRound, &newProposal, true)
			}
			if v.Round == initRound {
				checkRoundState(t, v, initRound, &initProposal, true)
			}
		}
	})

	t.Run("test dump state handler", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		prevHeight := currentHeight.Sub(currentHeight, big.NewInt(1))
		prevBlock := generateBlock(prevHeight)
		knownMsgHash := []common.Hash{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5}, {0, 0, 1, 3, 6}}

		// round 0 messages
		initRound := int64(0)
		initProposalMsg, initProposal, initRoundPrevoteMsg, initRoundPrecommitMsg := prepareRoundMsgs(t, initRound, currentHeight, 0, members[initRound].Address, clientAddr, privateKeys[clientAddr])

		// round 1 messages
		currentRound := int64(1)
		sender := 1
		newProposalMsg, newProposal, curRoundPrevoteMsg, curRoundPrecommitMsg := prepareRoundMsgs(t, currentRound, currentHeight, 0, members[currentRound].Address, members[sender].Address, privateKeys[members[sender].Address])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().KnownMsgHash().Return(knownMsgHash)

		c := New(backendMock, config.DefaultConfig())
		c.address = clientAddr
		// round 0 messages:
		genRoundVoteMessages(c, initRound, &initProposal, initProposalMsg, initRoundPrevoteMsg, initRoundPrecommitMsg, true)
		// current round messages:
		genRoundVoteMessages(c, currentRound, &newProposal, newProposalMsg, curRoundPrevoteMsg, curRoundPrecommitMsg, true)

		setCoreState(c, currentHeight, currentRound, propose, initProposal.ProposalBlock, initRound, initProposal.ProposalBlock, initRound, committeeSet,
			prevBlock.Header())

		go c.handleStateDump()

		state := types.TendermintState{}
		// wait for response with timeout.
		timeout := time.After(time.Second)
		select {
		case s := <-c.coreStateCh:
			state = s
		case <-timeout:
			t.Fatal("fetch tendermint state time out")
		}

		checkState(t, c, state, currentHeight, currentRound, initRound, initProposal, newProposal, prevBlock, knownMsgHash)
	})

	t.Run("test RPC callback to dump state time out", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Post(coreStateRequestEvent{}).Times(1)
		core := New(backendMock, config.DefaultConfig())
		state := core.CoreState()
		require.Equal(t, int64(-1), state.Code)
	})
}
