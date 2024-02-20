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
	c.LogPrevoteMessageEvent("MessageEvent(Prevote): Received", prevote, prevote.Sender().String(), c.address.String())

	// check upon conditions for current round proposal
	c.currentPrevoteChecks(ctx)
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
