package core

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"math/big"
)

func (c *Core) CoreState() TendermintState {
	// send state dump request.
	var e = CoreStateRequestEvent{
		StateChan: make(chan TendermintState),
	}
	go c.SendEvent(e)
	return <-e.StateChan
}

// State Dump is handled in the main loop triggered by an event rather than using RLOCK mutex.
func (c *Core) handleStateDump(e CoreStateRequestEvent) {
	state := TendermintState{
		Client:            c.address,
		BlockPeriod:       c.blockPeriod,
		CurHeightMessages: msgForDump(c.CurrentHeightMessages()),
		BacklogMessages:   getBacklogMsgs(c),
		UncheckedMsgs:     getBacklogUncheckedMsgs(c),
		// tendermint Core state:
		Height:      c.Height(),
		Round:       c.Round(),
		Step:        uint64(c.step),
		Proposal:    getProposal(c, c.Round()),
		LockedValue: getHash(c.lockedValue),
		LockedRound: c.lockedRound,
		ValidValue:  getHash(c.validValue),
		ValidRound:  c.validRound,

		// committee state
		Committee:       c.CommitteeSet().Committee(),
		Proposer:        c.CommitteeSet().GetProposer(c.Round()).Address,
		IsProposer:      c.IsProposer(),
		QuorumVotePower: c.CommitteeSet().Quorum(),
		RoundStates:     getRoundState(c),
		// extra state
		SentProposal:          c.sentProposal,
		SentPrevote:           c.sentPrevote,
		SentPrecommit:         c.sentPrecommit,
		SetValidRoundAndValue: c.setValidRoundAndValue,
		// timer state
		ProposeTimerStarted:   c.proposeTimeout.TimerStarted(),
		PrevoteTimerStarted:   c.prevoteTimeout.TimerStarted(),
		PrecommitTimerStarted: c.precommitTimeout.TimerStarted(),
		// known msgs in case of gossiping.
		KnownMsgHash: c.backend.KnownMsgHash(),
	}

	// for none blocking send state.
	c.logger.Debug("sending Core state msg")
	e.StateChan <- state
	// let sender to close channel.
	close(e.StateChan)
}

func getBacklogUncheckedMsgs(c *Core) []*MsgForDump {
	result := make([]*MsgForDump, 0)
	for _, ms := range c.backlogUntrusted {
		result = append(result, msgForDump(ms)...)
	}

	return result
}

// getBacklogUncheckedMsgs and getBacklogMsgs are kind of redundant code,
// don't know how to write it via golang like template in C++, since the only
// difference is the type of the data it operate on.
func getBacklogMsgs(c *Core) []*MsgForDump {
	result := make([]*MsgForDump, 0)
	for _, ms := range c.backlogs {
		result = append(result, msgForDump(ms)...)
	}

	return result
}

func msgForDump(messages []message.Message) []*MsgForDump {
	result := make([]*MsgForDump, 0, len(messages))
	for _, m := range messages {
		msg := new(MsgForDump)
		msg.Message = m
		msg.Power = m.Power()
		msg.Hash = types.RLPHash(m.Bytes) // ]

		// in case of haven't decode msg yet, set round and height as -1.
		msg.Round = -1
		msg.Height = big.NewInt(-1)
		msg.Round = m.R()
		msg.Height = big.NewInt(int64(m.H()))
		result = append(result, msg)
	}
	return result
}

func getProposal(c *Core, round int64) *common.Hash {
	if c.messages.GetOrCreate(round).ProposalDetails != nil && c.messages.GetOrCreate(round).ProposalDetails.ProposalBlock != nil {
		v := c.messages.GetOrCreate(round).ProposalDetails.ProposalBlock.Hash()
		return &v
	}
	return nil
}

func getHash(b *types.Block) *common.Hash {
	if b != nil {
		v := b.Hash()
		return &v
	}
	return nil
}

func getRoundState(c *Core) []RoundState {
	rounds := c.messages.GetRounds()
	states := make([]RoundState, 0, len(rounds))

	for _, r := range rounds {
		proposal, prevoteState, preCommitState := getVoteState(c.messages, r)
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

func blockHashes(messages map[common.Hash]map[common.Address]message.Message) []common.Hash {
	blockHashes := make([]common.Hash, 0, len(messages))
	for key := range messages {
		blockHashes = append(blockHashes, key)
	}
	return blockHashes
}

func getVoteState(s *message.Map, round int64) (common.Hash, []VoteState, []VoteState) {
	p := common.Hash{}

	if s.GetOrCreate(round).Proposal() != nil && s.GetOrCreate(round).Proposal().ProposalBlock != nil {
		p = s.GetOrCreate(round).Proposal().ProposalBlock.Hash()
	}

	preVoteValues := blockHashes(s.GetOrCreate(round).Prevotes.Votes)
	preCommitValues := blockHashes(s.GetOrCreate(round).Precommits.Votes)
	prevoteState := make([]VoteState, 0, len(preVoteValues))
	precommitState := make([]VoteState, 0, len(preCommitValues))

	for _, v := range preVoteValues {
		var s = VoteState{
			Value:            v,
			ProposalVerified: s.GetOrCreate(round).IsProposalVerified(),
			VotePower:        s.GetOrCreate(round).PrevotesPower(v),
		}
		prevoteState = append(prevoteState, s)
	}

	for _, v := range preCommitValues {
		var s = VoteState{
			Value:            v,
			ProposalVerified: s.GetOrCreate(round).IsProposalVerified(),
			VotePower:        s.GetOrCreate(round).PrecommitsPower(v),
		}
		precommitState = append(precommitState, s)
	}

	return p, prevoteState, precommitState
}
