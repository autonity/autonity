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
		state.Code = -1
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
		Proposal:    getProposal(c, c.Round()),
		LockedValue: getLockedValue(c),
		LockedRound: c.lockedRound,
		ValidValue:  getValidValue(c),
		ValidRound:  c.validRound,

		// committee state:
		ParentCommittee: getParentCommittee(c),
		Committee:       c.committeeSet().Committee(),
		Proposer:        c.committeeSet().GetProposer(c.Round()).Address,
		IsProposer:      c.isProposer(),
		QuorumVotePower: c.committeeSet().Quorum(),
		RoundStates:     getRoundState(c),
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
		Code:         0,
	}
	c.coreStateCh <- state
}

func getProposal(c *core, round int64) *common.Hash {
	if c.messages.getOrCreate(round).proposal != nil && c.messages.getOrCreate(round).proposal.ProposalBlock != nil {
		v := c.messages.getOrCreate(round).proposal.ProposalBlock.Hash()
		return &v
	}
	return nil
}

func getLockedValue(c *core) *common.Hash {
	if c.lockedValue != nil {
		v := c.lockedValue.Hash()
		return &v
	}
	return nil
}

func getValidValue(c *core) *common.Hash {
	if c.validValue != nil {
		v := c.validValue.Hash()
		return &v
	}
	return nil
}

func getParentCommittee(c *core) types.Committee {
	v := types.Committee{}
	if c.lastHeader != nil {
		v = c.lastHeader.Committee
	}
	return v
}

func getRoundState(c *core) []types.RoundState {
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
