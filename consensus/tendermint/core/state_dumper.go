package core

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"math/big"
	"time"
)

type coreStateRequestEvent struct {
}

// VoteState save the prevote or precommit voting status for a specific value.
type VoteState struct {
	Value            common.Hash
	ProposalVerified bool
	VotePower        uint64
}

// RoundState save the voting status for a specific round.
type RoundState struct {
	Round          int64
	Proposal       common.Hash
	PrevoteState   []VoteState
	PrecommitState []VoteState
}

// TendermintState save an instant status for the tendermint consensus engine.
type TendermintState struct {
	// return error code, 0 for okay, -1 for timeout.
	Code int64
	// validator address
	Client common.Address

	// core state of tendermint
	Height      big.Int
	Round       int64
	Step        uint64
	Proposal    *common.Hash
	LockedValue *common.Hash
	LockedRound int64
	ValidValue  *common.Hash
	ValidRound  int64

	// committee state
	ParentCommittee types.Committee
	Committee       types.Committee
	Proposer        common.Address
	IsProposer      bool
	QuorumVotePower uint64
	RoundStates     []RoundState
	ProposerPolicy  uint64

	// extra state
	SentProposal          bool
	SentPrevote           bool
	SentPrecommit         bool
	SetValidRoundAndValue bool

	// timer state
	BlockPeriod           uint64
	ProposeTimerStarted   bool
	PrevoteTimerStarted   bool
	PrecommitTimerStarted bool

	// current height messages and known message in case of gossip.
	CurHeightMessages []*Message
	KnownMsgHash      []common.Hash
}

func (c *core) CoreState() TendermintState {
	state := TendermintState{}
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
	state := TendermintState{
		Client:            c.address,
		ProposerPolicy:    uint64(c.proposerPolicy),
		BlockPeriod:       c.blockPeriod,
		CurHeightMessages: c.messages.GetMessages(),
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

func getRoundState(c *core) []RoundState {
	rounds := c.messages.getRounds()
	states := make([]RoundState, 0, len(rounds))

	for _, r := range rounds {
		proposal, prevoteState, preCommitState := c.messages.getVoteState(r)
		state := RoundState{
			Round:          r,
			Proposal:       proposal,
			PrevoteState:   prevoteState,
			PrecommitState: preCommitState,
		}
		states = append(states, state)
	}
	return states
}
