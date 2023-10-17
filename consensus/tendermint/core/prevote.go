package core

import (
	"context"
	"errors"

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
	}
	prevote := message.NewPrevote(c.Round(), c.Height().Uint64(), value, c.backend.Sign)
	c.LogPrevoteMessageEvent("MessageEvent(Prevote): Sent", prevote, c.address.String(), "broadcast")
	c.sentPrevote = true
	c.Broadcaster().Broadcast(prevote)
}

func (c *Prevoter) HandlePrevote(ctx context.Context, prevote *message.Prevote) error {
	if err := c.checkMessageStep(prevote.R(), prevote.H(), Prevote); err != nil {
		// Store old round prevote messages for future rounds since it is required for validRound
		if errors.Is(err, constants.ErrOldRoundMessage) {
			// We only process old rounds while future rounds messages are pushed on to the backlog
			oldRoundMessages := c.messages.GetOrCreate(prevote.R())
			oldRoundMessages.AddPrevote(prevote)

			// Line 28 in Algorithm 1 of The latest gossip on BFT consensus.
			if c.step == Propose {
				// ProposalBlock would be nil if node haven't received the proposal yet.
				if proposal := c.curRoundMessages.Proposal(); proposal != nil {
					vr := proposal.ValidRound()
					h := proposal.Block().Hash()
					rs := c.messages.GetOrCreate(vr)

					if vr >= 0 && vr < c.Round() && rs.PrevotesPower(h).Cmp(c.CommitteeSet().Quorum()) >= 0 {
						c.SendPrevote(ctx, !(c.lockedRound <= vr || h == c.lockedValue.Hash()))
						c.SetStep(Prevote)
						return nil
					}
				}
			}
		}
		return err
	}

	// After checking the message we know it is from the same height and round, so we should store it even if
	// c.curRoundMessages.Step() < prevote. The propose Timeout which is started at the beginning of the round
	// will update the step to at least prevote and when it handle its on preVote(nil), then it will also have
	// votes from other nodes.
	c.curRoundMessages.AddPrevote(prevote)

	c.LogPrevoteMessageEvent("MessageEvent(Prevote): Received", prevote, prevote.Sender().String(), c.address.String())

	// Now we can add the preVote to our current round state
	if c.step >= Prevote {
		curProposal := c.curRoundMessages.Proposal()
		// Line 36 in Algorithm 1 of The latest gossip on BFT consensus
		if curProposal != nil && c.curRoundMessages.PrevotesPower(curProposal.Block().Hash()).Cmp(c.CommitteeSet().Quorum()) >= 0 && !c.setValidRoundAndValue {
			// this piece of code should only run once
			if err := c.prevoteTimeout.StopTimer(); err != nil {
				return err
			}
			c.logger.Debug("Stopped Scheduled Prevote Timeout")

			if c.step == Prevote {
				c.lockedValue = curProposal.Block()
				c.lockedRound = c.Round()
				c.precommiter.SendPrecommit(ctx, false)
				c.SetStep(Precommit)
			}
			c.validValue = curProposal.Block()
			c.validRound = c.Round()
			c.setValidRoundAndValue = true
			// Line 44 in Algorithm 1 of The latest gossip on BFT consensus
		} else if c.step == Prevote && c.curRoundMessages.PrevotesPower(common.Hash{}).Cmp(c.CommitteeSet().Quorum()) >= 0 {
			if err := c.prevoteTimeout.StopTimer(); err != nil {
				return err
			}
			c.logger.Debug("Stopped Scheduled Prevote Timeout")
			c.precommiter.SendPrecommit(ctx, true)
			c.SetStep(Precommit)
			// Line 34 in Algorithm 1 of The latest gossip on BFT consensus
		} else if c.step == Prevote && !c.prevoteTimeout.TimerStarted() && !c.sentPrecommit && c.curRoundMessages.PrevotesTotalPower().Cmp(c.CommitteeSet().Quorum()) >= 0 {
			timeoutDuration := c.timeoutPrevote(c.Round())
			c.prevoteTimeout.ScheduleTimeout(timeoutDuration, c.Round(), c.Height(), c.onTimeoutPrevote)
			c.logger.Debug("Scheduled Prevote Timeout", "Timeout Duration", timeoutDuration)
		}
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
