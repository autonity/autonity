package core

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
)

type StateRequestEvent struct {
	StateChan chan interfaces.CoreState
}

func (c *Core) CoreState() interfaces.CoreState {
	// send state dump request.
	var e = StateRequestEvent{
		StateChan: make(chan interfaces.CoreState),
	}
	go c.SendEvent(e)
	return <-e.StateChan
}

func msgForDump(msg message.Msg) interfaces.MsgForDump {
	mDump := interfaces.MsgForDump{
		Code:           msg.Code(),
		Hash:           msg.Hash(),
		Payload:        msg.Payload(),
		Height:         msg.H(),
		Round:          msg.R(),
		SignatureInput: msg.SignatureInput(),
		Signature:      msg.Signature().Marshal(),
		Power:          msg.Power(),
	}

	switch msg.Code() {
	case message.ProposalCode:
		propose := msg.(*message.Propose)

		mDump.Block = propose.Block()
		mDump.ValidRound = propose.ValidRound()
		mDump.Signer = propose.Signer()

	case message.PrevoteCode, message.PrecommitCode:
		vote := msg.(message.Vote)
		mDump.Signers = vote.Signers()
		mDump.Value = vote.Value()

	default:
		panic("Unknown message type")

	}
	return mDump
}

func msgsForDump(msgs []message.Msg) []interfaces.MsgForDump {
	mDumps := make([]interfaces.MsgForDump, 0, len(msgs))
	for _, msg := range msgs {
		mDumps = append(mDumps, msgForDump(msg))
	}
	return mDumps
}

// State Dump is handled in the main loop triggered by an event rather than using RLOCK mutex.
func (c *Core) handleStateDump(e StateRequestEvent) {
	state := interfaces.CoreState{
		Client:              c.address,
		BlockPeriod:         c.blockPeriod,
		CurHeightMessages:   msgsForDump(c.messages.All()),
		FutureRoundMessages: msgsForDump(getFutureRoundMsgs(c)),
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
		Committee:       *c.CommitteeSet().Committee().Copy(),
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

func getFutureRoundMsgs(c *Core) []message.Msg {
	c.futureRoundLock.RLock()
	defer c.futureRoundLock.RUnlock()
	result := make([]message.Msg, 0)
	for _, msgs := range c.futureRound {
		result = append(result, msgs...)
	}

	return result
}

func getProposal(c *Core, round int64) *common.Hash {
	if c.messages.GetOrCreate(round).Proposal() != nil && c.messages.GetOrCreate(round).Proposal().Block() != nil {
		v := c.messages.GetOrCreate(round).Proposal().Block().Hash()
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

func getRoundState(c *Core) []interfaces.RoundState {
	rounds := c.messages.GetRounds()
	states := make([]interfaces.RoundState, 0, len(rounds))

	for _, r := range rounds {
		proposal, prevoteState, preCommitState := getVoteState(c.messages, r)
		state := interfaces.RoundState{
			Round:          r,
			Proposal:       proposal,
			PrevoteState:   prevoteState,
			PrecommitState: preCommitState,
		}
		states = append(states, state)
	}
	return states
}

func blockHashes[T interface{ Value() common.Hash }](messages []message.Msg) []common.Hash {
	blockHashes := make([]common.Hash, 0, len(messages))
	for _, m := range messages {
		blockHashes = append(blockHashes, m.(T).Value())
	}
	return blockHashes
}

func getVoteState(s *message.Map, round int64) (common.Hash, []interfaces.VoteState, []interfaces.VoteState) {
	messages := s.GetOrCreate(round)

	p := messages.ProposalHash()
	preVoteValues := blockHashes[*message.Prevote](messages.AllPrevotes())
	preCommitValues := blockHashes[*message.Precommit](messages.AllPrecommits())
	prevoteState := make([]interfaces.VoteState, 0, len(preVoteValues))
	precommitState := make([]interfaces.VoteState, 0, len(preCommitValues))

	for _, v := range preVoteValues {
		var s = interfaces.VoteState{
			Value:            v,
			ProposalVerified: s.GetOrCreate(round).IsProposalVerified(),
			VotePower:        s.GetOrCreate(round).PrevotesPower(v),
		}
		prevoteState = append(prevoteState, s)
	}

	for _, v := range preCommitValues {
		var s = interfaces.VoteState{
			Value:            v,
			ProposalVerified: s.GetOrCreate(round).IsProposalVerified(),
			VotePower:        s.GetOrCreate(round).PrecommitsPower(v),
		}
		precommitState = append(precommitState, s)
	}

	return p, prevoteState, precommitState
}
