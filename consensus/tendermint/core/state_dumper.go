package core

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	tctypes "github.com/autonity/autonity/consensus/tendermint/core/types"
	"github.com/autonity/autonity/core/types"
	"math/big"
)

func (c *Core) CoreState() tctypes.TendermintState {
	// send state dump request.
	var e = tctypes.CoreStateRequestEvent{
		StateChan: make(chan tctypes.TendermintState),
	}
	go c.SendEvent(e)
	return <-e.StateChan
}

// State Dump is handled in the main loop triggered by an event rather than using RLOCK mutex.
func (c *Core) handleStateDump(e tctypes.CoreStateRequestEvent) {
	state := tctypes.TendermintState{
		Client:            c.address,
		BlockPeriod:       c.blockPeriod,
		CurHeightMessages: msgForDump(c.GetCurrentHeightMessages()),
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

func getBacklogUncheckedMsgs(c *Core) []*tctypes.MsgForDump {
	result := make([]*tctypes.MsgForDump, 0)
	for _, ms := range c.backlogUnchecked {
		result = append(result, msgForDump(ms)...)
	}

	return result
}

// getBacklogUncheckedMsgs and getBacklogMsgs are kind of redundant code,
// don't know how to write it via golang like template in C++, since the only
// difference is the type of the data it operate on.
func getBacklogMsgs(c *Core) []*tctypes.MsgForDump {
	result := make([]*tctypes.MsgForDump, 0)
	for _, ms := range c.backlogs {
		result = append(result, msgForDump(ms)...)
	}

	return result
}

func msgForDump(msgs []*messageutils.Message) []*tctypes.MsgForDump {
	result := make([]*tctypes.MsgForDump, 0, len(msgs))
	for _, m := range msgs {
		msg := new(tctypes.MsgForDump)
		msg.Message = *m
		msg.Power = m.GetPower()
		msg.Hash = types.RLPHash(m.Payload)

		// in case of haven't decode msg yet, set round and height as -1.
		msg.Round = -1
		msg.Height = big.NewInt(-1)
		msg.Round, _ = m.Round()
		msg.Height, _ = m.Height()
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

func getRoundState(c *Core) []tctypes.RoundState {
	rounds := c.messages.GetRounds()
	states := make([]tctypes.RoundState, 0, len(rounds))

	for _, r := range rounds {
		proposal, prevoteState, preCommitState := getVoteState(c.messages, r)
		state := tctypes.RoundState{
			Round:          r,
			Proposal:       proposal,
			PrevoteState:   prevoteState,
			PrecommitState: preCommitState,
		}
		states = append(states, state)
	}
	return states
}

func blockHashes(messages map[common.Hash]map[common.Address]messageutils.Message) []common.Hash {
	blockHashes := make([]common.Hash, 0, len(messages))
	for key := range messages {
		blockHashes = append(blockHashes, key)
	}
	return blockHashes
}

func getVoteState(s *messageutils.MessagesMap, round int64) (common.Hash, []tctypes.VoteState, []tctypes.VoteState) {
	p := common.Hash{}

	if s.GetOrCreate(round).Proposal() != nil && s.GetOrCreate(round).Proposal().ProposalBlock != nil {
		p = s.GetOrCreate(round).Proposal().ProposalBlock.Hash()
	}

	preVoteValues := blockHashes(s.GetOrCreate(round).Prevotes.Votes)
	preCommitValues := blockHashes(s.GetOrCreate(round).Precommits.Votes)
	prevoteState := make([]tctypes.VoteState, 0, len(preVoteValues))
	precommitState := make([]tctypes.VoteState, 0, len(preCommitValues))

	for _, v := range preVoteValues {
		var s = tctypes.VoteState{
			Value:            v,
			ProposalVerified: s.GetOrCreate(round).IsProposalVerified(),
			VotePower:        s.GetOrCreate(round).PrevotesPower(v),
		}
		prevoteState = append(prevoteState, s)
	}

	for _, v := range preCommitValues {
		var s = tctypes.VoteState{
			Value:            v,
			ProposalVerified: s.GetOrCreate(round).IsProposalVerified(),
			VotePower:        s.GetOrCreate(round).PrecommitsPower(v),
		}
		precommitState = append(precommitState, s)
	}

	return p, prevoteState, precommitState
}
