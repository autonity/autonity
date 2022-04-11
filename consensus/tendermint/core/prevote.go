package core

import (
	"context"

	"github.com/autonity/autonity/common"
)

func (c *core) sendPrevote(ctx context.Context, isNil bool) {
	logger := c.logger.New("step", c.step)

	var prevote = Vote{
		Round:  c.Round(),
		Height: c.Height(),
	}

	if isNil {
		prevote.ProposedBlockHash = common.Hash{}
	} else {
		if h := c.curRoundMessages.GetProposalHash(); h == (common.Hash{}) {
			c.logger.Error("sendPrecommit Proposal is empty! It should not be empty!")
			return
		}
		prevote.ProposedBlockHash = c.curRoundMessages.GetProposalHash()
	}

	encodedVote, err := Encode(&prevote)
	if err != nil {
		logger.Error("Failed to encode", "subject", prevote)
		return
	}

	c.logPrevoteMessageEvent("MessageEvent(Prevote): Sent", prevote, c.address.String(), "broadcast")

	c.sentPrevote = true
	c.broadcast(ctx, &Message{
		Code:          msgPrevote,
		Msg:           encodedVote,
		Address:       c.address,
		CommittedSeal: []byte{},
	})
}

func (c *core) handlePrevote(ctx context.Context, msg *Message) error {
	var preVote Vote
	err := msg.Decode(&preVote)
	if err != nil {
		return errFailedDecodePrevote
	}

	if err = c.checkMessage(preVote.Round, preVote.Height, prevote); err != nil {
		// Store old round prevote messages for future rounds since it is required for validRound
		if err == errOldRoundMessage {
			// We only process old rounds while future rounds messages are pushed on to the backlog
			oldRoundMessages := c.messages.getOrCreate(preVote.Round)
			c.acceptVote(oldRoundMessages, prevote, preVote.ProposedBlockHash, *msg)

			// Line 28 in Algorithm 1 of The latest gossip on BFT consensus.
			if c.step == propose {
				// ProposalBlock would be nil if node haven't receive proposal yet.
				if c.curRoundMessages.proposal.ProposalBlock != nil {
					vr := c.curRoundMessages.proposal.ValidRound
					h := c.curRoundMessages.proposal.ProposalBlock.Hash()
					rs := c.messages.getOrCreate(vr)

					if vr >= 0 && vr < c.Round() && rs.PrevotesPower(h) >= c.committeeSet().Quorum() {
						c.sendPrevote(ctx, !(c.lockedRound <= vr || h == c.lockedValue.Hash()))
						c.setStep(prevote)
						return nil
					}
				}
			}
		}
		return err
	}

	// After checking the message we know it is from the same height and round, so we should store it even if
	// c.curRoundMessages.Step() < prevote. The propose timeout which is started at the beginning of the round
	// will update the step to at least prevote and when it handle its on preVote(nil), then it will also have
	// votes from other nodes.
	prevoteHash := preVote.ProposedBlockHash
	c.acceptVote(c.curRoundMessages, prevote, prevoteHash, *msg)

	c.logPrevoteMessageEvent("MessageEvent(Prevote): Received", preVote, msg.Address.String(), c.address.String())

	// Now we can add the preVote to our current round state
	if c.step >= prevote {
		curProposalHash := c.curRoundMessages.GetProposalHash()

		// Line 36 in Algorithm 1 of The latest gossip on BFT consensus
		if curProposalHash != (common.Hash{}) && c.curRoundMessages.PrevotesPower(curProposalHash) >= c.committeeSet().Quorum() && !c.setValidRoundAndValue {
			// this piece of code should only run once
			if err := c.prevoteTimeout.stopTimer(); err != nil {
				return err
			}
			c.logger.Debug("Stopped Scheduled Prevote Timeout")

			if c.step == prevote {
				c.lockedValue = c.curRoundMessages.Proposal().ProposalBlock
				c.lockedRound = c.Round()
				c.sendPrecommit(ctx, false)
				c.setStep(precommit)
			}
			c.validValue = c.curRoundMessages.Proposal().ProposalBlock
			c.validRound = c.Round()
			c.setValidRoundAndValue = true
			// Line 44 in Algorithm 1 of The latest gossip on BFT consensus
		} else if c.step == prevote && c.curRoundMessages.PrevotesPower(common.Hash{}) >= c.committeeSet().Quorum() {
			if err := c.prevoteTimeout.stopTimer(); err != nil {
				return err
			}
			c.logger.Debug("Stopped Scheduled Prevote Timeout")

			c.sendPrecommit(ctx, true)
			c.setStep(precommit)

			// Line 34 in Algorithm 1 of The latest gossip on BFT consensus
		} else if c.step == prevote && !c.prevoteTimeout.timerStarted() && !c.sentPrecommit && c.curRoundMessages.PrevotesTotalPower() >= c.committeeSet().Quorum() {
			timeoutDuration := c.timeoutPrevote(c.Round())
			c.prevoteTimeout.scheduleTimeout(timeoutDuration, c.Round(), c.Height(), c.onTimeoutPrevote)
			c.logger.Debug("Scheduled Prevote Timeout", "Timeout Duration", timeoutDuration)
		}
	}

	return nil
}

func (c *core) logPrevoteMessageEvent(message string, prevote Vote, from, to string) {
	currentProposalHash := c.curRoundMessages.GetProposalHash()
	c.logger.Debug(message,
		"from", from,
		"to", to,
		"currentHeight", c.Height(),
		"msgHeight", prevote.Height,
		"currentRound", c.Round(),
		"msgRound", prevote.Round,
		"currentStep", c.step,
		"isProposer", c.isProposer(),
		"currentProposer", c.committeeSet().GetProposer(c.Round()),
		"isNilMsg", prevote.ProposedBlockHash == common.Hash{},
		"hash", prevote.ProposedBlockHash,
		"type", "Prevote",
		"totalVotes", c.curRoundMessages.PrevotesTotalPower(),
		"totalNilVotes", c.curRoundMessages.PrevotesPower(common.Hash{}),
		"quorum", c.committee.Quorum(),
		"VoteProposedBlock", c.curRoundMessages.PrevotesPower(currentProposalHash),
	)
}
