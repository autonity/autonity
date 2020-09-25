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
		state.Code = "time out"
		c.logger.Debug("Waiting for tendermint core state timed out", "elapsed", time.Second)
	}

	return state
}

// State Dump is handled in the main loop triggered by an event rather than using RLOCK mutex.
func (c *core) handleStateDump() {
	state := types.TendermintState{
		Client:            c.address,
		ProposerPolicy:    uint64(c.proposerPolicy),
		BlockPeriod:       c.blockPeriod,
		CurHeightMessages: c.messages.CopyMessages(),
		// tendermint core state:
		Height:      *c.Height(),
		Round:       c.Round(),
		Step:        uint64(c.step),
		Proposal:    c.getProposal(c.Round()),
		LockedValue: c.getLockedValue(),
		LockedRound: c.lockedRound,
		ValidValue:  c.getValidValue(),
		ValidRound:  c.validRound,

		// committee state:
		ParentCommittee: c.getParentCommittee(),
		Committee:       c.committeeSet().Committee(),
		Proposer:        c.committeeSet().GetProposer(c.Round()).Address,
		IsProposer:      c.isProposer(),
		QuorumVotePower: c.committeeSet().Quorum(),
		RoundStates:     c.getRoundState(),
		// extra state
		SentProposal:          c.sentProposal,
		SentPrevote:           c.sentPrevote,
		SentPrecommit:         c.sentPrecommit,
		SetValidRoundAndValue: c.setValidRoundAndValue,
		// timer state
		ProposeTimerStarted:   c.proposeTimeout.timerStarted(),
		PrevoteTimerStarted:   c.prevoteTimeout.timerStarted(),
		PrecommitTimerStarted: c.precommitTimeout.timerStarted(),
		// known msgs in case of gossiping.
		KnownMsgHash: c.backend.KnownMsgHash(),
		Code:         "done",
	}
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
