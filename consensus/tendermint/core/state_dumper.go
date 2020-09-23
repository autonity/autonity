package core

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"time"
)

type coreStateRequestEvent struct {
}

func (c *core) CoreState() types.TendermintState {
	state := types.TendermintState{}
	// send state dump request.
	var e = coreStateRequestEvent{}
	go c.sendEvent(e)
	// wait for response with timeout.
	timeout := time.After(time.Second)
	select {
	case s := <-c.coreStateCh:
		state = s
	case <-timeout:
		c.logger.Debug("Waiting for tendermint core state timed out", "elapsed", time.Second)
	}

	return state
}

// State Dump is handled in the main loop triggered by an event rather than using RLOCK mutex.
func (c *core) handleStateDump() {
	state := types.TendermintState{
		Client:         c.address,
		ProposerPolicy: uint64(c.proposerPolicy),
		BlockPeriod:    c.blockPeriod,
	}

	state.CurHeightMessages = c.messages.CopyMessages()

	// tendermint core state
	state.Height = *c.Height()
	state.Round = c.Round()
	state.Step = uint64(c.step)
	state.Proposal = c.getProposal(state.Round)
	state.LockedValue = c.getLockedValue()
	state.LockedRound = c.lockedRound
	state.ValidValue = c.getValidValue()
	state.ValidRound = c.validRound

	// committee state
	state.ParentCommittee = c.getParentCommittee()
	state.Committee = c.committeeSet().Committee()
	state.Proposer = c.committeeSet().GetProposer(state.Round).Address
	state.IsProposer = c.isProposer()
	state.QuorumVotePower = c.committeeSet().Quorum()
	state.RoundStates = c.getRoundState()

	// extra state
	state.SentProposal = c.sentProposal
	state.SentPrevote = c.sentPrevote
	state.SentPrecommit = c.sentPrecommit
	state.SetValidRoundAndValue = c.setValidRoundAndValue

	// timer state
	state.ProposeTimerStarted = c.proposeTimeout.timerStarted()
	state.PrevoteTimerStarted = c.prevoteTimeout.timerStarted()
	state.PrecommitTimerStarted = c.precommitTimeout.timerStarted()

	// known msgs in case of gossiping.
	state.KnownMsgHash = c.backend.KnownMsgHash()
	c.coreStateCh <- state
}

func (c *core) getProposal(round int64) common.Hash {
	v := common.Hash{}
	if c.messages.getOrCreate(round).proposal != nil && c.messages.getOrCreate(round).proposal.ProposalBlock != nil {
		v = c.messages.getOrCreate(round).proposal.ProposalBlock.Hash()
	}
	return v
}

func (c *core) getLockedValue() common.Hash {
	v := common.Hash{}
	if c.lockedValue != nil {
		v = c.lockedValue.Hash()
	}
	return v
}

func (c *core) getValidValue() common.Hash {
	v := common.Hash{}
	if c.validValue != nil {
		v = c.validValue.Hash()
	}
	return v
}

func (c *core) getParentCommittee() types.Committee {
	v := types.Committee{}
	if c.lastHeader != nil {
		v = c.lastHeader.Committee
	}
	return v
}

func (c *core) getRoundState() []types.RoundState {
	rounds := c.messages.getRounds()
	states := make([]types.RoundState, 0, len(rounds))

	for _, r := range rounds {
		proposal, prevoteState, preCommitState := c.messages.getVoteState(r)
		state := types.RoundState{
			Round:          r,
			Proposal:       proposal,
			PrevoteState:   prevoteState,
			PrecommitState: preCommitState,
		}
		states = append(states, state)
	}
	return states
}
