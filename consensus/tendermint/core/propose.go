package core

import (
	"github.com/clearmatics/autonity/common"
	"time"

	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint"
	"github.com/clearmatics/autonity/core/types"
)

func (c *core) sendProposal(p *types.Block) {
	logger := c.logger.New("step", c.step)

	// If I'm the proposer and I have the same height with the proposal
	if c.currentRoundState.Height().Cmp(p.Number()) == 0 && c.isProposer() && !c.sentProposal {
		proposalBlock := tendermint.NewProposal(c.currentRoundState.Round(), c.currentRoundState.Height(), c.validRound, p)
		proposal, err := Encode(proposalBlock)
		if err != nil {
			logger.Error("Failed to encode", "Round", proposalBlock.Round, "Height", proposalBlock.Height, "ValidRound", c.validRound)
			return
		}
		c.sentProposal = true
		c.backend.SetProposedBlockHash(p.Hash())

		c.logProposalMessageEvent("MessageEvent(Proposal): Sent", proposalBlock, c.address.String())

		c.broadcast(&message{
			Code: msgProposal,
			Msg:  proposal,
		})
	}
}

func (c *core) handleProposal(msg *message, sender tendermint.Validator) error {
	logger := c.logger.New("from", sender, "step", c.step)

	var proposal *tendermint.Proposal
	err := msg.Decode(&proposal)
	if err != nil {
		return errFailedDecodeProposal
	}

	c.logProposalMessageEvent("MessageEvent(Proposal): Received", proposal, msg.Address.String())

	// Ensure we have the same view with the Proposal message
	if err := c.checkMessage(proposal.Round, proposal.Height); err != nil {
		// We don't care about old proposals so they are ignored
		return err
	}

	// Check if the message comes from currentRoundState proposer
	if !c.valSet.IsProposer(sender.Address()) {
		logger.Warn("Ignore proposal messages from non-proposer")
		return errNotFromProposer
	}

	// Verify the proposal we received
	if duration, err := c.backend.Verify(*proposal.ProposalBlock); err != nil {
		logger.Warn("Failed to verify proposal", "err", err, "duration", duration)
		// if it's a future block, we will handle it again after the duration
		// TIME FIELD OF HEADER CHECKED HERE - NOT HEIGHT
		// TODO: implement wiggle time / median time
		if err == consensus.ErrFutureBlock {
			c.stopFutureProposalTimer()
			c.futureProposalTimer = time.AfterFunc(duration, func() {
				c.sendEvent(backlogEvent{
					src: sender,
					msg: msg,
				})
			})
		}
		return err
	}

	// Here is about to accept the Proposal
	if c.step == StepAcceptProposal {
		if err := c.stopProposeTimeout(); err == nil {
			// Set the proposal for the current round
			c.currentRoundState.SetProposal(proposal)

			logger.Info("Accepted Proposal", "height", proposal.Height, "round", proposal.Round, "Hash", proposal.ProposalBlock.Hash())

			vr := proposal.ValidRound.Int64()
			h := proposal.ProposalBlock.Hash()
			curR := c.currentRoundState.round.Int64()

			if vr == -1 {
				// Line 22 in Algorithm 1 of The latest gossip on BFT consensus
				if c.lockedRound.Int64() == vr || h == c.lockedValue.Hash() {
					c.sendPrevote(false)
				} else {
					c.sendPrevote(true)
				}
				c.setStep(StepProposeDone)
				// Line 28 in Algorithm 1 of The latest gossip on BFT consensus
			} else if rs, ok := c.currentHeightRoundsStates[vr]; vr > -1 &&
				vr < curR &&
				ok &&
				c.quorum(rs.Prevotes.VotesSize(h)) {
				if c.lockedRound.Int64() <= vr || h == c.lockedValue.Hash() {
					c.sendPrevote(false)
				} else {
					c.sendPrevote(true)
				}
				c.setStep(StepProposeDone)
			} else {
				return errNoMajority
			}
		} else {
			return err
		}
	}

	return nil
}

func (c *core) stopProposeTimeout() error {
	if c.proposeTimeout.started {
		if stopped := c.proposeTimeout.stopTimer(); !stopped {
			return errNilPrevoteSent
		}
	}
	return nil
}

func (c *core) logProposalMessageEvent(message string, proposal *tendermint.Proposal, from string) {
	c.logger.Info(message,
		"from", from,
		"type", "Proposal",
		"currentHeight", c.currentRoundState.height,
		"msgHeight", proposal.Height,
		"currentRound", c.currentRoundState.round,
		"msgRound", proposal.Round,
		"currentSteo", c.step,
		"msgStep", c.step,
		"currentProposer", c.valSet.GetProposer(),
		"isNilMsg", proposal.ProposalBlock.Hash() == common.Hash{},
		"hash", proposal.ProposalBlock.Hash(),
	)
}
