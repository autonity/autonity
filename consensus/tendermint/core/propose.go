package core

import (
	"context"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core/types"
)

func (c *core) sendProposal(ctx context.Context, p *types.Block) {
	logger := c.logger.New("step", c.step)

	// If I'm the proposer and I have the same height with the proposal
	if c.Height().Cmp(p.Number()) == 0 && c.isProposer() && !c.sentProposal {
		proposalBlock := NewProposal(c.Round(), c.Height(), c.validRound, p)
		proposal, err := Encode(proposalBlock)
		if err != nil {
			logger.Error("Failed to encode", "Round", proposalBlock.Round, "Height", proposalBlock.Height, "ValidRound", c.validRound)
			return
		}

		c.sentProposal = true
		c.backend.SetProposedBlockHash(p.Hash())

		c.logProposalMessageEvent("MessageEvent(Proposal): Sent", *proposalBlock, c.address.String(), "broadcast")

		c.broadcast(ctx, &Message{
			Code:          msgProposal,
			Msg:           proposal,
			Address:       c.address,
			CommittedSeal: []byte{},
		})
	}
}

func (c *core) handleProposal(ctx context.Context, msg *Message) error {
	var proposal Proposal
	err := msg.Decode(&proposal)
	if err != nil {
		return errFailedDecodeProposal
	}

	// Ensure we have the same view with the Proposal message
	if err := c.checkMessage(proposal.Round, proposal.Height, propose); err != nil {
		// If it's a future round proposal, the only upon condition
		// that can be triggered is L49, but this requires more than F future round messages
		// meaning that a future roundchange will happen before, as such, pushing the
		// message to the backlog is fine.
		if err == errOldRoundMessage {

			roundMsgs := c.messages.getOrCreate(proposal.Round)

			// if we already have a proposal then it must be different than the current one
			// it can't happen unless someone's byzantine.
			if roundMsgs.proposal != nil {
				return err // do not gossip, TODO: accountability
			}

			if !c.isProposerMsg(proposal.Round, msg.Address) {
				c.logger.Warn("Ignore proposal messages from non-proposer")
				return errNotFromProposer
			}
			// We do not verify the proposal in this case.
			roundMsgs.SetProposal(&proposal, msg, false)

			if roundMsgs.PrecommitsPower(roundMsgs.GetProposalHash()) >= c.committeeSet().Quorum() {
				if _, error := c.backend.VerifyProposal(*proposal.ProposalBlock); error != nil {
					return error
				}
				c.logger.Debug("Committing old round proposal")
				c.commit(proposal.Round, roundMsgs)
				return nil
			}
		}
		return err
	}

	// Check if the message comes from curRoundMessages proposer
	if !c.isProposerMsg(c.Round(), msg.Address) {
		c.logger.Warn("Ignore proposal messages from non-proposer")
		return errNotFromProposer
	}

	// Verify the proposal we received
	if duration, err := c.backend.VerifyProposal(*proposal.ProposalBlock); err != nil {

		if timeoutErr := c.proposeTimeout.stopTimer(); timeoutErr != nil {
			return timeoutErr
		}
		// if it's a future block, we will handle it again after the duration
		// TODO: implement wiggle time / median time
		if err == consensus.ErrFutureBlock {
			c.stopFutureProposalTimer()
			c.futureProposalTimer = time.AfterFunc(duration, func() {
				c.sendEvent(backlogEvent{
					msg: msg,
				})
			})
			return err
		}
		c.sendPrevote(ctx, true)
		// do not to accept another proposal in current round
		c.setStep(prevote)

		c.logger.Warn("Failed to verify proposal", "err", err, "duration", duration)

		return err
	}

	// Set the proposal for the current round
	c.curRoundMessages.SetProposal(&proposal, msg, true)

	c.logProposalMessageEvent("MessageEvent(Proposal): Received", proposal, msg.Address.String(), c.address.String())

	//l49: Check if we have a quorum of precommits for this proposal
	curProposalHash := c.curRoundMessages.GetProposalHash()
	if c.curRoundMessages.PrecommitsPower(curProposalHash) >= c.committeeSet().Quorum() {
		c.commit(proposal.Round, c.curRoundMessages)
		return nil
	}

	if c.step == propose {
		if err := c.proposeTimeout.stopTimer(); err != nil {
			return err
		}

		vr := proposal.ValidRound
		h := proposal.ProposalBlock.Hash()

		// Line 22 in Algorithm 1 of The latest gossip on BFT consensus
		if vr == -1 {
			// When lockedRound is set to any value other than -1 lockedValue is also
			// set to a non nil value. So we can be sure that we will only try to access
			// lockedValue when it is non nil.
			c.sendPrevote(ctx, !(c.lockedRound == -1 || h == c.lockedValue.Hash()))
			c.setStep(prevote)
			return nil
		}

		rs := c.messages.getOrCreate(vr)

		// Line 28 in Algorithm 1 of The latest gossip on BFT consensus
		// vr >= 0 here
		if vr < c.Round() && rs.PrevotesPower(h) >= c.committeeSet().Quorum() {
			c.sendPrevote(ctx, !(c.lockedRound <= vr || h == c.lockedValue.Hash()))
			c.setStep(prevote)
		}
	}

	return nil
}

func (c *core) stopFutureProposalTimer() {
	if c.futureProposalTimer != nil {
		c.futureProposalTimer.Stop()
	}
}

func (c *core) logProposalMessageEvent(message string, proposal Proposal, from, to string) {
	c.logger.Debug(message,
		"type", "Proposal",
		"from", from,
		"to", to,
		"currentHeight", c.Height(),
		"msgHeight", proposal.Height,
		"currentRound", c.Round(),
		"msgRound", proposal.Round,
		"currentStep", c.step,
		"isProposer", c.isProposer(),
		"currentProposer", c.committeeSet().GetProposer(c.Round()),
		"isNilMsg", proposal.ProposalBlock.Hash() == common.Hash{},
		"hash", proposal.ProposalBlock.Hash(),
	)
}
