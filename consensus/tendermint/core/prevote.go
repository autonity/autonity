package core

import (
	"context"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
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
	self, err := c.CommitteeSet().MemberByAddress(c.address)
	if err != nil {
		// it can happen in edge case addressed by docker e2e test, that is a validator resets at the epoch boundary,
		// after which it leaves the committee, we cannot panic it in that case.
		c.logger.Error("Validator is no longer in current committee", "err", err, "validator", c.address.String())
		return
	}
	prevote := message.NewPrevote(c.Round(), c.Height().Uint64(), value, c.backend.Sign, self, c.CommitteeSet().Committee().Len())
	c.LogPrevoteMessageEvent("MessageEvent(Prevote): Sent", prevote)
	c.sentPrevote = true
	c.Broadcaster().Broadcast(prevote)
	if metrics.Enabled() {
		PrevoteSentBlockTSDeltaBg.Add(time.Since(c.currBlockTimeStamp).Nanoseconds())
	}
}

func (c *Prevoter) HandlePrevote(ctx context.Context, prevote *message.Prevote) error {
	if prevote.R() > c.Round() {
		return constants.ErrFutureRoundMessage
	}
	if prevote.R() < c.Round() {
		// We only process old rounds while future rounds messages are pushed on to the backlog
		oldRoundMessages := c.messages.GetOrCreate(prevote.R())
		oldRoundMessages.AddPrevote(prevote)
		go c.SendEvent(events.PowerChangeEvent{Height: c.Height().Uint64(), Round: c.Round(), Code: message.PrevoteCode, Value: prevote.Value()})

		// Proposal would be nil if node haven't received the proposal yet.
		proposal := c.curRoundMessages.Proposal()
		if proposal == nil {
			return constants.ErrOldRoundMessage
		}

		// Line 28 in Algorithm 1 of The latest gossip on BFT consensus.
		// check if we have quorum prevotes on vr
		c.oldProposalCheck(ctx, proposal)
		return constants.ErrOldRoundMessage
	}

	// After checking the message we know it is from the same height and round, so we should store it even if
	// c.curRoundMessages.Step() < prevote. The propose Timeout which is started at the beginning of the round
	// will update the step to at least prevote and when it handle its on preVote(nil), then it will also have
	// votes from other nodes.
	c.curRoundMessages.AddPrevote(prevote)
	go c.SendEvent(events.PowerChangeEvent{Height: c.Height().Uint64(), Round: c.Round(), Code: message.PrevoteCode, Value: prevote.Value()})

	c.LogPrevoteMessageEvent("MessageEvent(Prevote): Received", prevote)
	// check upon conditions for current round proposal
	c.currentPrevoteChecks(ctx)
	return nil
}

func (c *Prevoter) LogPrevoteMessageEvent(message string, prevote *message.Prevote) {
	if c.logger.Enabled(context.Background(), log.LevelDebug) {
		c.logger.Debug(message,
			"type", "Prevote",
			"local address", c.Address().String(),
			"currentHeight", c.Height(),
			"msgHeight", prevote.H(),
			"currentRound", c.Round(),
			"msgRound", prevote.R(),
			"currentStep", c.step,
			"isProposer", c.IsProposer(),
			"currentProposer", c.CommitteeSet().GetProposer(c.Round()),
			"isNilMsg", prevote.Value() == common.Hash{},
			"value", prevote.Value(),
			"totalVotes", c.curRoundMessages.PrevotesTotalPower(),
			"totalNilVotes", c.curRoundMessages.PrevotesPower(common.Hash{}),
			"quorum", c.committee.Quorum(),
			"VoteProposedBlock", c.curRoundMessages.PrevotesPower(c.curRoundMessages.ProposalHash()),
			"prevote", prevote.String(),
		)
	}
}
