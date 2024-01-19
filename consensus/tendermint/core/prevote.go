package core

import (
	"context"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
)

type Prevoter struct {
	*Core
}

func (c *Prevoter) SendPrevote(ctx context.Context, isNil bool) {
	value := common.Hash{}
	if !isNil {
		proposal := c.curRoundMessages.Proposal()
		if proposal == nil {
			c.logger.Error("sendPrevote Proposal is empty! It should not be empty!")
			return
		}
		value = proposal.Block().Hash()
		c.logger.Info("Prevoting on proposal", "proposal", value, "round", c.Round(), "height", c.Height().Uint64())
	} else {
		c.logger.Info("Prevoting on nil", "round", c.Round(), "height", c.Height().Uint64())
	}
	prevote := message.NewPrevote(c.Round(), c.Height().Uint64(), value, c.backend.Sign)
	c.LogPrevoteMessageEvent("MessageEvent(Prevote): Sent", prevote, c.address.String(), "broadcast")
	c.sentPrevote = true
	c.Broadcaster().Broadcast(prevote)
}

func (c *Prevoter) HandlePrevote(ctx context.Context, prevote *message.Prevote) error {
	if prevote.R() > c.Round() {
		return constants.ErrFutureRoundMessage
	}
	if prevote.R() < c.Round() {
		// We only process old rounds while future rounds messages are pushed on to the backlog
		oldRoundMessages := c.messages.GetOrCreate(prevote.R())
		oldRoundMessages.AddPrevote(prevote)
		if c.step != Propose {
			return constants.ErrOldRoundMessage
		}
		// Current step is Propose
		// Line 28 in Algorithm 1 of The latest gossip on BFT consensus.
		// ProposalBlock would be nil if node haven't received the proposal yet.
		proposal := c.curRoundMessages.Proposal()
		if proposal == nil {
			return constants.ErrOldRoundMessage
		}
		vr := proposal.ValidRound()
		h := proposal.Block().Hash()
		rs := c.messages.GetOrCreate(vr)
		if vr >= 0 && vr < c.Round() && rs.PrevotesPower(h).Cmp(c.CommitteeSet().Quorum()) >= 0 {
			c.SendPrevote(ctx, !(c.lockedRound <= vr || h == c.lockedValue.Hash()))
			c.SetStep(ctx, Prevote)
			return nil
		}
		return constants.ErrOldRoundMessage
	}

	// After checking the message we know it is from the same height and round, so we should store it even if
	// c.curRoundMessages.Step() < prevote. The propose Timeout which is started at the beginning of the round
	// will update the step to at least prevote and when it handle its on preVote(nil), then it will also have
	// votes from other nodes.
	c.curRoundMessages.AddPrevote(prevote)
	c.LogPrevoteMessageEvent("MessageEvent(Prevote): Received", prevote, prevote.Sender().String(), c.address.String())
	if c.step == Propose {
		return nil
	}
	// We are at step Prevote or Precommit from here
	// Now we can add the preVote to our current round state
	curProposal := c.curRoundMessages.Proposal()
	// Line 36 in Algorithm 1 of The latest gossip on BFT consensus
	if curProposal != nil && c.curRoundMessages.PrevotesPower(curProposal.Block().Hash()).Cmp(c.CommitteeSet().Quorum()) >= 0 && !c.setValidRoundAndValue {
		if c.step == Prevote {
			c.lockedValue = curProposal.Block()
			c.lockedRound = c.Round()
			c.precommiter.SendPrecommit(ctx, false)
			c.SetStep(ctx, Precommit)
		}
		c.validValue = curProposal.Block()
		c.validRound = c.Round()
		c.setValidRoundAndValue = true
		// Line 44 in Algorithm 1 of The latest gossip on BFT consensus
	} else if c.step == Prevote && c.curRoundMessages.PrevotesPower(common.Hash{}).Cmp(c.CommitteeSet().Quorum()) >= 0 {
		c.precommiter.SendPrecommit(ctx, true)
		c.SetStep(ctx, Precommit)
		// Line 34 in Algorithm 1 of The latest gossip on BFT consensus
	} else if c.step == Prevote && !c.prevoteTimeout.TimerStarted() && !c.sentPrecommit && c.curRoundMessages.PrevotesTotalPower().Cmp(c.CommitteeSet().Quorum()) >= 0 {
		timeoutDuration := c.timeoutPrevote(c.Round())
		c.prevoteTimeout.ScheduleTimeout(timeoutDuration, c.Round(), c.Height(), c.onTimeoutPrevote)
		c.logger.Debug("Scheduled Prevote Timeout", "Timeout Duration", timeoutDuration)
	}
	return nil
}

func (c *Prevoter) LogPrevoteMessageEvent(message string, prevote *message.Prevote, from, to string) {
	currentProposalHash := c.curRoundMessages.ProposalHash()
	c.logger.Debug(message,
		"from", from,
		"to", to,
		"currentHeight", c.Height(),
		"msgHeight", prevote.H(),
		"currentRound", c.Round(),
		"msgRound", prevote.R(),
		"currentStep", c.step,
		"isProposer", c.IsProposer(),
		"currentProposer", c.CommitteeSet().GetProposer(c.Round()),
		"isNilMsg", prevote.Value() == common.Hash{},
		"value", prevote.Value(),
		"type", "Prevote",
		"totalVotes", c.curRoundMessages.PrevotesTotalPower(),
		"totalNilVotes", c.curRoundMessages.PrevotesPower(common.Hash{}),
		"quorum", c.committee.Quorum(),
		"VoteProposedBlock", c.curRoundMessages.PrevotesPower(currentProposalHash),
	)
}
