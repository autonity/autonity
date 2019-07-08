package core

import (
	"context"
	"math/big"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
)

func (c *core) sendPrevote(ctx context.Context, isNil bool) {
	logger := c.logger.New("step", c.currentRoundState.Step())

	var prevote = &tendermint.Vote{
		Round:  big.NewInt(c.currentRoundState.Round().Int64()),
		Height: big.NewInt(c.currentRoundState.Height().Int64()),
	}

	if isNil {
		prevote.ProposedBlockHash = common.Hash{}
	} else {
		if h := c.currentRoundState.GetCurrentProposalHash(); h == (common.Hash{}) {
			c.logger.Error("sendPrecommit Proposal is empty! It should not be empty!")
			return
		}
		prevote.ProposedBlockHash = c.currentRoundState.GetCurrentProposalHash()
	}

	encodedVote, err := Encode(prevote)
	if err != nil {
		logger.Error("Failed to encode", "subject", prevote)
		return
	}

	c.logPrevoteMessageEvent("MessageEvent(Prevote): Sent", prevote, c.address.String(), "broadcast")

	c.sentPrevote = true
	c.broadcast(ctx, &message{
		Code: msgPrevote,
		Msg:  encodedVote,
	})
}

func (c *core) handlePrevote(ctx context.Context, msg *message) error {
	var preVote *tendermint.Vote
	err := msg.Decode(&preVote)
	if err != nil {
		return errFailedDecodePrevote
	}

	if err = c.checkMessage(preVote.Round, preVote.Height); err != nil {
		// We want to store old round messages for future rounds since it is required for validRound
		if err == errOldRoundMessage {
			// The roundstate must exist as every roundstate is added to c.currentHeightRoundsState at startRound
			// And we only process old rounds while future rounds messages are pushed on to the backlog
			prevoteMS := c.currentHeightRoundsStates[preVote.Round.Int64()].Prevotes
			if preVote.ProposedBlockHash == (common.Hash{}) {
				prevoteMS.AddNilVote(*msg)
			} else {
				prevoteMS.AddVote(preVote.ProposedBlockHash, *msg)
			}
		}
		return err
	}

	// After checking the message we know it is from the same height and round, so we should store it even if
	// c.currentRoundState.Step() < prevote. The propose timeout which is started at the beginning of the round
	// will update the step to at least prevote and when it handle its on preVote(nil), then it will also have
	// votes from other nodes.
	prevoteHash := preVote.ProposedBlockHash
	if prevoteHash == (common.Hash{}) {
		c.currentRoundState.Prevotes.AddNilVote(*msg)
	} else {
		c.currentRoundState.Prevotes.AddVote(prevoteHash, *msg)
	}

	c.logPrevoteMessageEvent("MessageEvent(Prevote): Received", preVote, msg.Address.String(), c.address.String())

	// Now we can add the preVote to our current round state
	if c.currentRoundState.Step() >= prevote {
		curProposalHash := c.currentRoundState.GetCurrentProposalHash()
		curR := c.currentRoundState.Round().Int64()
		curH := c.currentRoundState.Height().Int64()

		// Line 36 in Algorithm 1 of The latest gossip on BFT consensus
		if curProposalHash != (common.Hash{}) && c.quorum(c.currentRoundState.Prevotes.VotesSize(curProposalHash)) && !c.setValidRoundAndValue {
			// this piece of code should only run once
			if err := c.stopPrevoteTimeout(); err != nil {
				return err
			}

			if c.currentRoundState.Step() == prevote {
				c.lockedValue = c.currentRoundState.Proposal().ProposalBlock
				c.lockedRound = big.NewInt(curR)
				c.sendPrecommit(ctx, false)
				c.setStep(precommit)
			}
			c.validValue = c.currentRoundState.Proposal().ProposalBlock
			c.validRound = big.NewInt(curR)
			c.setValidRoundAndValue = true
			// Line 44 in Algorithm 1 of The latest gossip on BFT consensus
		} else if c.currentRoundState.Step() == prevote && c.quorum(c.currentRoundState.Prevotes.NilVotesSize()) {
			if err := c.stopPrevoteTimeout(); err != nil {
				return err
			}
			c.sendPrecommit(ctx, true)
			c.setStep(precommit)

			// Line 34 in Algorithm 1 of The latest gossip on BFT consensus
		} else if c.currentRoundState.Step() == prevote && !c.prevoteTimeout.started && !c.sentPrecommit && c.quorum(c.currentRoundState.Prevotes.TotalSize()) {
			timeoutDuration := timeoutPrevote(curR)
			c.prevoteTimeout.scheduleTimeout(timeoutDuration, curR, curH, c.onTimeoutPrevote)
			c.logger.Debug("Scheduled Prevote Timeout", "Timeout Duration", timeoutDuration)
		}
	}

	return nil
}

func (c *core) stopPrevoteTimeout() error {
	if c.prevoteTimeout.started {
		c.logger.Debug("Stopping Scheduled Prevote Timeout")
		if stopped := c.prevoteTimeout.stopTimer(); !stopped {
			return errNilPrecommitSent
		}
	}
	return nil
}

func (c *core) logPrevoteMessageEvent(message string, prevote *tendermint.Vote, from, to string) {
	currentProposalHash := c.currentRoundState.GetCurrentProposalHash()
	c.logger.Debug(message,
		"from", from,
		"to", to,
		"currentHeight", c.currentRoundState.Height(),
		"msgHeight", prevote.Height,
		"currentRound", c.currentRoundState.Round(),
		"msgRound", prevote.Round,
		"currentStep", c.currentRoundState.Step(),
		"isProposer", c.isProposer(),
		"currentProposer", c.valSet.GetProposer(),
		"isNilMsg", prevote.ProposedBlockHash == common.Hash{},
		"hash", prevote.ProposedBlockHash,
		"type", "Prevote",
		"totalVotes", c.currentRoundState.Prevotes.TotalSize(),
		"totalNilVotes", c.currentRoundState.Prevotes.NilVotesSize(),
		"quorumReject", c.quorum(c.currentRoundState.Prevotes.NilVotesSize()),
		"totalNonNilVotes", c.currentRoundState.Prevotes.VotesSize(currentProposalHash),
		"quorumAccept", c.quorum(c.currentRoundState.Prevotes.VotesSize(currentProposalHash)),
	)
}
