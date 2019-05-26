package core

import (
	"time"

	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint"
	"github.com/clearmatics/autonity/core/types"
)

// TODO: add new message struct for proposal (proposalMessage) and determine how to rlp encode them especially nil
// TODO: add new message for vote (prevote and precommit) and determine how to rlp encode them especially nil
func (c *core) sendProposal(p *types.Block) {
	logger := c.logger.New("step", c.step)

	// If I'm the proposer and I have the same height with the proposal
	if c.currentRoundState.Height().Cmp(p.Number()) == 0 && c.isProposer() && !c.sentProposal {
		r, h, vr := c.currentRoundState.Round(), c.currentRoundState.Height(), c.validRound
		proposal, err := Encode(&tendermint.Proposal{
			Round:         r,
			Height:        h,
			ValidRound:    vr,
			ProposalBlock: *p,
		})
		if err != nil {
			logger.Error("Failed to encode", "Round", r, "Height", h, "ValidRound", vr)
			return
		}
		c.sentProposal = true
		c.backend.SetProposedBlockHash(p.Hash())
		c.broadcast(&message{
			Code: msgProposal,
			Msg:  proposal,
		})
	}
}

func (c *core) handleProposal(msg *message, src tendermint.Validator) error {
	logger := c.logger.New("from", src, "step", c.step)

	var proposal *tendermint.Proposal
	err := msg.Decode(&proposal)
	if err != nil {
		return errFailedDecodeProposal
	}

	// Ensure we have the same view with the Proposal message
	// If it is old message, see if we need to broadcast COMMIT
	if err := c.checkMessage(msgProposal, proposal.Round, proposal.Height); err != nil {
		// We don't care about old proposals so they are ignored
		return err
	}

	// Check if the message comes from currentRoundState proposer
	if !c.valSet.IsProposer(src.Address()) {
		logger.Warn("Ignore proposal messages from non-proposer")
		return errNotFromProposer
	}

	// Verify the proposal we received
	if duration, err := c.backend.Verify(proposal.ProposalBlock); err != nil {
		logger.Warn("Failed to verify proposal", "err", err, "duration", duration)
		// if it's a future block, we will handle it again after the duration
		// TIME FIELD OF HEADER CHECKED HERE - NOT HEIGHT
		if err == consensus.ErrFutureBlock {
			c.stopFutureProposalTimer()
			c.futureProposalTimer = time.AfterFunc(duration, func() {
				c.sendEvent(backlogEvent{
					src: src,
					msg: msg,
				})
			})
		}
		return err
	}

	// Here is about to accept the Proposal
	if c.step == StepAcceptProposal {
		if c.proposeTimeout.started {
			if stopped := c.proposeTimeout.stopTimer(); !stopped {
				return errNilPrevoteSent
			}
		}
		c.acceptProposal(proposal)
		vr := proposal.ValidRound.Int64()
		h := proposal.ProposalBlock.Hash()

		if vr == -1 {
			if c.lockedRound.Int64() == proposal.ValidRound.Int64() || h == c.lockedValue.Hash() {
				c.sendPrevote(false)
			} else {
				c.sendPrevote(true)
			}
			c.setStep(StepProposeDone)
		} else if rs, ok := c.currentHeightRoundsStates[vr]; vr > -1 && vr < c.currentRoundState.round.Int64() && ok && c.quorum(rs.Prevotes.VotesSize(h)) {
			if c.lockedRound.Int64() <= proposal.ValidRound.Int64() || h == c.lockedValue.Hash() {
				c.sendPrevote(false)
			} else {
				c.sendPrevote(true)
			}
			c.setStep(StepProposeDone)
		} else {
			return errNoMajority
		}

	}

	return nil
}

func (c *core) acceptProposal(proposal *tendermint.Proposal) {
	c.currentRoundState.SetProposal(proposal)
}
