package core

import (
	"context"
	"github.com/clearmatics/autonity/common"
	"math/big"
)

// Line 22 in Algorithm 1 of the latest gossip on BFT consensus
func (c *core) checkForNewProposal(ctx context.Context, round int64) error {
	proposalMS := c.getProposalSet(round)
	if proposalMS == nil {
		// Have not received proposal
		return nil
	}
	proposal := proposalMS.proposal()
	proposalMsg := proposalMS.proposalMsg()

	h := proposal.ProposalBlock.Hash()

	if c.isProposerForR(round, proposalMsg.Address) && c.getStep() == propose {
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
			c.sendPrevote(ctx, big.NewInt(c.getHeight().Int64()), big.NewInt(round), h)
		} else {
			c.sendPrevote(ctx, big.NewInt(c.getHeight().Int64()), big.NewInt(round), common.Hash{})
		}
		if err := c.setStep(ctx, prevote); err != nil {
			return err
		}
	}
	return nil
}

// Line 28 in Algorithm 1 of the latest gossip on BFT consensus
func (c *core) checkForOldProposal(ctx context.Context, round int64) error {
	proposalMS := c.getProposalSet(round)
	if proposalMS == nil {
		// Have not received proposal
		return nil
	}
	proposal := proposalMS.proposal()
	proposalMsg := proposalMS.proposalMsg()

	vr := proposal.ValidRound.Int64()

	validRoundPrevotes := c.getPrevotesSet(vr)
	if validRoundPrevotes == nil {
		// Have not received any prevotes for the valid round
		return nil
	}

	h := proposal.ProposalBlock.Hash()

	if c.isProposerForR(round, proposalMsg.Address) && c.quorum(validRoundPrevotes.VotesSize(h)) &&
		c.getStep() == propose && vr >= 0 && vr < round {
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
			c.sendPrevote(ctx, big.NewInt(c.getHeight().Int64()), big.NewInt(round), h)
		} else {
			c.sendPrevote(ctx, big.NewInt(c.getHeight().Int64()), big.NewInt(round), common.Hash{})
		}
		if err := c.setStep(ctx, prevote); err != nil {
			return err
		}
	}
	return nil
}

// Line 34 in Algorithm 1 of the latest gossip on BFT consensus
func (c *core) checkForPrevoteTimeout(round int64, height int64) {
	prevotes := c.getPrevotesSet(round)
	if prevotes == nil {
		// Do not have any prevotes for the round
		return
	}
	if c.getStep() == prevote && !c.prevoteTimeout.timerStarted() && !c.sentPrecommit && c.quorum(prevotes.TotalSize()) {
		timeoutDuration := timeoutPrevote(round)
		c.prevoteTimeout.scheduleTimeout(timeoutDuration, round, height, c.onTimeoutPrevote)
		c.logger.Debug("Scheduled Prevote Timeout", "Timeout Duration", timeoutDuration)
	}
}

// Line 36 in Algorithm 1 of the latest gossip on BFT consensus
func (c *core) checkForQuorumPrevotes(ctx context.Context, round int64) error {
	proposalMS := c.getProposalSet(round)
	if proposalMS == nil {
		// Have not received proposal
		return nil
	}
	proposal := proposalMS.proposal()
	proposalMsg := proposalMS.proposalMsg()

	prevotes := c.getPrevotesSet(round)
	if prevotes == nil {
		// Have not received any prevotes for round
		return nil
	}

	h := proposal.ProposalBlock.Hash()

	// this piece of code should only run once per round
	if c.isProposerForR(round, proposalMsg.Address) && c.quorum(prevotes.VotesSize(h)) &&
		c.getStep() >= prevote && !c.setValidRoundAndValue {
		if valid, err := c.isValid(proposal.ProposalBlock); !valid {
			return err
		}

		if c.prevoteTimeout.timerStarted() {
			if err := c.prevoteTimeout.stopTimer(); err != nil {
				return err
			}
			c.logger.Debug("Stopped Scheduled Prevote Timeout")
		}

		if c.getStep() == prevote {
			c.lockedValue = proposal.ProposalBlock
			c.lockedRound = big.NewInt(round)
			c.sendPrecommit(ctx, big.NewInt(c.getHeight().Int64()), big.NewInt(round), h)
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
	prevotes := c.getPrevotesSet(round)
	if prevotes == nil {
		// Have not received any prevotes for round
		return nil
	}

	if c.getStep() == prevote && c.quorum(prevotes.NilVotesSize()) {
		if c.prevoteTimeout.timerStarted() {
			if err := c.prevoteTimeout.stopTimer(); err != nil {
				return err
			}
			c.logger.Debug("Stopped Scheduled Prevote Timeout")
		}

		c.sendPrecommit(ctx, big.NewInt(c.getHeight().Int64()), big.NewInt(round), common.Hash{})
		_ = c.setStep(ctx, precommit)
	}
	return nil
}

// Line 47 in Algorithm 1 of the latest gossip on BFT consensus
func (c *core) checkForPrecommitTimeout(round int64, height int64) {
	precommits := c.getPrecommitsSet(round)
	if precommits == nil {
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
	proposalMS := c.getProposalSet(round)
	if proposalMS == nil {
		// Have not received proposal
		return nil
	}
	proposal := proposalMS.proposal()
	proposalMsg := proposalMS.proposalMsg()

	precommits := c.getPrecommitsSet(round)
	if precommits == nil {
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

// Line 55 in Algorithm 1 of the latest gossip on BFT consensus
func (c *core) checkForFutureRoundChange(ctx context.Context, round int64) {
	var messages []*Message
	proposalMS, prevotes, precommits := c.getProposalSet(round), c.getPrevotesSet(round), c.getPrecommitsSet(round)

	if proposalMS != nil {
		messages = append(messages, proposalMS.pMsg)
	}

	if prevotes != nil {
		messages = append(messages, prevotes.GetMessages()...)
	}

	if precommits != nil {
		messages = append(messages, precommits.GetMessages()...)
	}

	if len(messages) <= c.valSet.F() {
		// Not enough message to move to future round
		return
	}

	// check for distinct messages
	addrMap := make(map[common.Address]struct{})

	for _, msg := range messages {
		if _, ok := addrMap[msg.Address]; ok {
			// If the message address is already in the map (i.e there are prevote, precommit and/or proposal from the
			// same sender, therefore, continue to next message)
			continue
		}
		addrMap[msg.Address] = struct{}{}
	}

	if len(addrMap) > c.valSet.F() {
		c.logger.Info("Received ceil(N/3) - 1 messages for higher round", "New round", round)
		go c.sendEvent(StartRoundEvent{round: round})
	}
}
