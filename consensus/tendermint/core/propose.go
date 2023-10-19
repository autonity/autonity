package core

import (
	"context"
	"errors"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	tctypes "github.com/autonity/autonity/consensus/tendermint/core/types"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/metrics"
)

type Proposer struct {
	*Core
}

func (c *Proposer) SendProposal(ctx context.Context, p *types.Block) {
	// If I'm the proposer and I have the same height with the proposal
	if c.Height().Cmp(p.Number()) == 0 && c.IsProposer() && !c.sentProposal {
		proposal := message.NewPropose(c.Round(), c.Height().Uint64(), c.validRound, p, c.backend.Sign)

		c.sentProposal = true
		c.backend.SetProposedBlockHash(p.Hash())

		if metrics.Enabled {
			now := time.Now()
			ProposalSentTimer.Update(now.Sub(c.newRound))
			ProposalSentBg.Add(now.Sub(c.newRound).Nanoseconds())
		}

		c.LogProposalMessageEvent("MessageEvent(Proposal): Sent", proposal, c.address.String(), "broadcast")

		c.Br().Broadcast(ctx, proposal)
	}
}

func (c *Proposer) HandleProposal(ctx context.Context, proposal *message.Propose) error {
	// Ensure we have the same view with the Proposal message
	if err := c.CheckMessage(proposal.R(), proposal.H(), Propose); err != nil {
		// If it's a future round proposal2, the only upon condition
		// that can be triggered is L49, but this requires more than F future round messages
		// meaning that a future roundchange will happen before, as such, pushing the
		// message to the backlog is fine.
		if errors.Is(err, constants.ErrOldRoundMessage) {

			roundMessages := c.messages.GetOrCreate(proposal.Round)

			// if we already have a proposal then it must be different than the current one
			// it can't happen unless someone's byzantine.
			if roundMessages.Proposal() != nil {
				return err // do not gossip, TODO: accountability
			}

			if !c.IsFromProposer(proposal.Round, proposal.Sender()) {
				c.logger.Warn("Ignoring proposal from non-proposer")
				return constants.ErrNotFromProposer
			}
			// We do not verify the proposal2 in this case.
			roundMessages.SetProposal(proposal, false)

			if roundMessages.PrecommitsPower(roundMessages.ProposalHash()).Cmp(c.CommitteeSet().Quorum()) >= 0 {
				if _, err2 := c.backend.VerifyProposal(proposal.ProposalBlock); err2 != nil {
					return err2
				}
				c.logger.Debug("Committing old round proposal2")
				c.Commit(proposal.Round, roundMessages)
				return nil
			}
		}
		return err
	}

	// Check if the message comes from curRoundMessages proposer
	if !c.IsFromProposer(c.Round(), proposal.Address) {
		c.logger.Warn("Ignore proposal2 messages from non-proposer")
		return constants.ErrNotFromProposer
	}

	// received a current round proposal2
	if metrics.Enabled {
		now := time.Now()
		ProposalReceivedTimer.Update(now.Sub(c.newRound))
		ProposalReceivedBg.Add(now.Sub(c.newRound).Nanoseconds())
	}

	// Verify the proposal2 we received
	start := time.Now()
	duration, err := c.backend.VerifyProposal(proposal2.ProposalBlock)

	if metrics.Enabled {
		now := time.Now()
		ProposalVerifiedTimer.Update(now.Sub(start))
		ProposalVerifiedBg.Add(now.Sub(start).Nanoseconds())
	}

	if err != nil {
		if timeoutErr := c.proposeTimeout.StopTimer(); timeoutErr != nil {
			return timeoutErr
		}
		// if it's a future block, we will handle it again after the duration
		// TODO: implement wiggle time / median time
		if err == consensus.ErrFutureTimestampBlock {
			c.StopFutureProposalTimer()
			c.futureProposalTimer = time.AfterFunc(duration, func() {
				c.SendEvent(backlogMessageEvent{
					msg: proposal,
				})
			})
			return err
		}
		c.prevoter.SendPrevote(ctx, true)
		// do not to accept another proposal2 in current round
		c.SetStep(tctypes.Prevote)

		c.logger.Warn("Failed to verify proposal2", "err", err, "duration", duration)

		return err
	}

	// Set the proposal2 for the current round
	c.curRoundMessages.SetProposal(proposal2, proposal, true)

	c.LogProposalMessageEvent("MessageEvent(Proposal): Received", proposal2, proposal.Address.String(), c.address.String())

	//l49: Check if we have a quorum of precommits for this proposal2
	curProposalHash := c.curRoundMessages.ProposalHash()
	if c.curRoundMessages.PrecommitsPower(curProposalHash).Cmp(c.CommitteeSet().Quorum()) >= 0 {
		c.Commit(proposal2.Round, c.curRoundMessages)
		return nil
	}

	if c.step == Propose {
		if err := c.proposeTimeout.StopTimer(); err != nil {
			return err
		}

		vr := proposal2.ValidRound
		h := proposal2.ProposalBlock.Hash()

		// Line 22 in Algorithm 1 of The latest gossip on BFT consensus
		if vr == -1 {
			// When lockedRound is set to any value other than -1 lockedValue is also
			// set to a non nil value. So we can be sure that we will only try to access
			// lockedValue when it is non nil.
			c.prevoter.SendPrevote(ctx, !(c.lockedRound == -1 || h == c.lockedValue.Hash()))
			c.SetStep(Prevote)
			return nil
		}

		rs := c.messages.GetOrCreate(vr)

		// Line 28 in Algorithm 1 of The latest gossip on BFT consensus
		// vr >= 0 here
		if vr < c.Round() && rs.PrevotesPower(h).Cmp(c.CommitteeSet().Quorum()) >= 0 {
			c.prevoter.SendPrevote(ctx, !(c.lockedRound <= vr || h == c.lockedValue.Hash()))
			c.SetStep(Prevote)
		}
	}

	return nil
}

func (c *Proposer) HandleNewCandidateBlockMsg(ctx context.Context, candidateBlock *types.Block) {
	if candidateBlock == nil {
		return
	}

	number := candidateBlock.Number()
	if currentIsHigher := c.Height().Cmp(number); currentIsHigher > 0 {
		c.logger.Info("NewCandidateBlockEvent: discarding old height candidateBlock", "number", number.Uint64())
		return
	}

	c.pendingCandidateBlocks[candidateBlock.NumberU64()] = candidateBlock

	// if current node is the proposer of current height and current round at step PROPOSE without available candidate
	// block sent before, if the incoming candidate block is the one it missed, send it now.
	if c.IsProposer() && c.step == tctypes.Propose && !c.sentProposal && c.Height().Cmp(number) == 0 {
		c.logger.Debug("NewCandidateBlockEvent: Sending proposal that was missed before", "number", number.Uint64())
		c.proposer.SendProposal(ctx, candidateBlock)
	}

	// release buffered candidate blocks before the height of current state machine.
	for height := range c.pendingCandidateBlocks {
		if height < c.Height().Uint64() {
			delete(c.pendingCandidateBlocks, height)
		}
	}
}

func (c *Proposer) StopFutureProposalTimer() {
	if c.futureProposalTimer != nil {
		c.futureProposalTimer.Stop()
	}
}

func (c *Proposer) LogProposalMessageEvent(message string, proposal *message.Propose, from, to string) {
	c.logger.Debug(message,
		"type", "Proposal",
		"from", from,
		"to", to,
		"currentHeight", c.Height(),
		"msgHeight", proposal.Height,
		"currentRound", c.Round(),
		"msgRound", proposal.Round,
		"currentStep", c.step,
		"isProposer", c.IsProposer(),
		"currentProposer", c.CommitteeSet().GetProposer(c.Round()),
		"isNilMsg", proposal.ProposalBlock.Hash() == common.Hash{},
		"hash", proposal.ProposalBlock.Hash(),
	)
}
