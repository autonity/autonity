package core

import (
	"context"
	"math/big"
)

// Line 22 in Algorithm 1 of the latest gossip on BFT consensus
func (c *core) checkForNewProposal(ctx context.Context, round int64) error {
	proposalMS, ok := c.allProposals[round]
	if !ok {
		// Have not received proposal
		return nil
	}
	proposal, proposalMsg := proposalMS.proposal, proposalMS.proposalMsg

	h := proposal.ProposalBlock.Hash()

	if c.isProposerForR(round, proposalMsg.Address) && c.step == propose {
		if valid, err := c.isValid(proposal.ProposalBlock); !valid {
			return err
		}

		// stop the timeout since a valid proposal has been received, if it cannot be stopped return
		if c.proposeTimeout.timerStarted() {
			if err := c.proposeTimeout.stopTimer(); err != nil {
				return err
			}
			c.logger.Debug("Stopped Scheduled Propose Timeout")
		}

		// Vote for the proposal if proposal is valid(proposal) && (lockedRound = -1 || lockedValue = proposal).
		if c.lockedRound.Int64() == -1 || (c.lockedRound != nil && h == c.lockedValue.Hash()) {
			c.sendPrevote(ctx, false)
		} else {
			c.sendPrevote(ctx, true)
		}
		if err := c.setStep(ctx, prevote); err != nil {
			return err
		}
	}
	return nil
}

// Line 28 in Algorithm 1 of the latest gossip on BFT consensus
func (c *core) checkForOldProposal(ctx context.Context, round int64) error {
	proposalMS, ok := c.allProposals[round]
	if !ok {
		// Have not received proposal
		return nil
	}
	proposal, proposalMsg := proposalMS.proposal, proposalMS.proposalMsg

	vr := proposal.ValidRound.Int64()

	validRoundPrevotes, ok := c.allPrevotes[vr]
	if !ok {
		// Have not received any prevotes for the valid round
		return nil
	}

	h := proposal.ProposalBlock.Hash()

	if c.isProposerForR(round, proposalMsg.Address) && c.quorum(validRoundPrevotes.VotesSize(h)) &&
		c.step == propose && vr >= 0 && vr < round {
		if valid, err := c.isValid(proposal.ProposalBlock); !valid {
			return err
		}

		// stop the timeout since a valid proposal has been received, if it cannot be stopped return
		if c.proposeTimeout.timerStarted() {
			if err := c.proposeTimeout.stopTimer(); err != nil {
				return err
			}
			c.logger.Debug("Stopped Scheduled Propose Timeout")
		}

		// Vote for proposal if valid(proposal) && ((0 <= lockedRound <= vr < curR) || lockedValue == proposal)
		if c.lockedRound.Int64() <= vr || (c.lockedRound != nil && h == c.lockedValue.Hash()) {
			c.sendPrevote(ctx, false)
		} else {
			c.sendPrevote(ctx, true)
		}
		if err := c.setStep(ctx, prevote); err != nil {
			return err
		}
	}
	return nil
}

// Line 34 in Algorithm 1 of the latest gossip on BFT consensus
func (c *core) checkForPrevoteTimeout(round int64, height int64) {
	prevotes, ok := c.allPrevotes[round]
	if !ok {
		// Do not have any prevotes for the round
		return
	}
	if c.step == prevote && !c.prevoteTimeout.timerStarted() && !c.sentPrecommit && c.quorum(prevotes.TotalSize()) {
		timeoutDuration := timeoutPrevote(round)
		c.prevoteTimeout.scheduleTimeout(timeoutDuration, round, height, c.onTimeoutPrevote)
		c.logger.Debug("Scheduled Prevote Timeout", "Timeout Duration", timeoutDuration)
	}
}

// Line 36 in Algorithm 1 of the latest gossip on BFT consensus
func (c *core) checkForQuorumPrevotes(ctx context.Context, round int64) error {
	proposalMS, ok := c.allProposals[round]
	if !ok {
		// Have not received proposal
		return nil
	}
	proposal, proposalMsg := proposalMS.proposal, proposalMS.proposalMsg

	prevotes, ok := c.allPrevotes[round]
	if !ok {
		// Have not received any prevotes for round
		return nil
	}

	h := proposal.ProposalBlock.Hash()

	// this piece of code should only run once per round
	if c.isProposerForR(round, proposalMsg.Address) && c.quorum(prevotes.VotesSize(h)) &&
		c.step >= prevote && !c.setValidRoundAndValue {
		if valid, err := c.isValid(proposal.ProposalBlock); !valid {
			return err
		}

		if c.prevoteTimeout.timerStarted() {
			if err := c.prevoteTimeout.stopTimer(); err != nil {
				return err
			}
			c.logger.Debug("Stopped Scheduled Prevote Timeout")
		}

		if c.step == prevote {
			c.lockedValue = proposal.ProposalBlock
			c.lockedRound = big.NewInt(round)
			c.sendPrecommit(ctx, false)
			_ = c.setStep(ctx, precommit)
		}
		c.validValue = proposal.ProposalBlock
		c.validRound = big.NewInt(round)
		c.setValidRoundAndValue = true

	}
	return nil
}

// Line 44 in Algorithm 1 of the latest gossip on BFT consensus
func (c *core) checkForQuorumPrevotesNil(ctx context.Context, round int64) error {
	prevotes, ok := c.allPrevotes[round]
	if !ok {
		// Have not received any prevotes for round
		return nil
	}

	if c.step == prevote && c.quorum(prevotes.NilVotesSize()) {
		if c.prevoteTimeout.timerStarted() {
			if err := c.prevoteTimeout.stopTimer(); err != nil {
				return err
			}
			c.logger.Debug("Stopped Scheduled Prevote Timeout")
		}

		c.sendPrecommit(ctx, true)
		_ = c.setStep(ctx, precommit)
	}
	return nil
}

// Line 47 in Algorithm 1 of the latest gossip on BFT consensus
func (c *core) checkForPrecommitTimeout(round int64, height int64) {
	precommits, ok := c.allPrecommits[round]
	if !ok {
		// Do not have any precommits for the round
		return
	}
	if !c.precommitTimeout.timerStarted() && c.quorum(precommits.TotalSize()) {
		timeoutDuration := timeoutPrecommit(round)
		c.precommitTimeout.scheduleTimeout(timeoutDuration, round, height, c.onTimeoutPrecommit)
		c.logger.Debug("Scheduled Precommit Timeout", "Timeout Duration", timeoutDuration)
	}
}

// Line 49 in Algorithm 1 of the latest gossip on BFT consensus
func (c *core) checkForConsensus(ctx context.Context, round int64) error {
	proposalMS, ok := c.allProposals[round]
	if !ok {
		// Have not received proposal
		return nil
	}
	proposal, proposalMsg := proposalMS.proposal, proposalMS.proposalMsg

	precommits, ok := c.allPrecommits[round]
	if !ok {
		// Have not received any precommits for round
		return nil
	}

	h := proposal.ProposalBlock.Hash()

	if c.isProposerForR(round, proposalMsg.Address) && c.quorum(precommits.VotesSize(h)) {
		if valid, err := c.isValid(proposal.ProposalBlock); !valid {
			return err
		}

		if c.precommitTimeout.timerStarted() {
			if err := c.precommitTimeout.stopTimer(); err != nil {
				return err
			}
			c.logger.Debug("Stopped Scheduled Precommit Timeout")
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			c.commit(ctx, round)
		}

	}
	return nil
}
