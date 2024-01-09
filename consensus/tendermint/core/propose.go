package core

import (
	"context"
	"errors"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/metrics"
)

type Proposer struct {
	*Core
}

func (c *Proposer) SendProposal(_ context.Context, block *types.Block) {
	// If I'm the proposer and I have the same height with the proposal
	if c.Height().Cmp(block.Number()) == 0 && c.IsProposer() && !c.sentProposal {
		proposal := message.NewPropose(c.Round(), c.Height().Uint64(), c.validRound, block, c.backend.Sign)
		c.sentProposal = true
		c.backend.SetProposedBlockHash(block.Hash())
		if metrics.Enabled {
			now := time.Now()
			ProposalSentTimer.Update(now.Sub(c.newRound))
			ProposalSentBg.Add(now.Sub(c.newRound).Nanoseconds())
		}
		c.LogProposalMessageEvent("MessageEvent(Proposal): Sent", proposal, c.address.String(), "broadcast")
		c.Broadcaster().Broadcast(proposal)
	}
}

func (c *Proposer) HandleProposal(ctx context.Context, proposal *message.Propose) error {
	// Ensure we have the same view with the Proposal message
	if err := c.checkMessage(proposal.R(), proposal.H()); err != nil {
		// If it's a future round proposal, the only upon condition
		// that can be triggered is L49, but this requires more than F future round messages
		// meaning that a future roundchange will happen before, as such, pushing the
		// message to the backlog is fine.
		if errors.Is(err, constants.ErrOldRoundMessage) {
			roundMessages := c.messages.GetOrCreate(proposal.R())
			// if we already have a proposal for this old round - ignore
			// the first proposal sent by the sender in a round is always the only one we consider.
			if roundMessages.Proposal() != nil {
				return constants.ErrAlreadyProcessed
			}

			if !c.IsFromProposer(proposal.R(), proposal.Sender()) {
				c.logger.Warn("Ignoring proposal from non-proposer")
				return constants.ErrNotFromProposer
			}
			// Save it, but do not verify the proposal yet unless we have enough precommits for it.
			roundMessages.SetProposal(proposal, false)
			if roundMessages.PrecommitsPower(proposal.Block().Hash()).Cmp(c.CommitteeSet().Quorum()) >= 0 {
				if _, err2 := c.backend.VerifyProposal(proposal.Block()); err2 != nil {
					return err2
				}
				c.logger.Debug("Committing old round proposal")
				c.Commit(ctx, proposal.R(), roundMessages)
				return nil
			}
		}
		return err
	}
	// if we already have processed a proposal in this round we ignore.
	if c.curRoundMessages.Proposal() != nil {
		return constants.ErrAlreadyProcessed
	}
	// At this point the local round matches the message round and the current step
	// could be either Proposal, Prevote, or Precommit.

	// Check if the message comes from curRoundMessages proposer
	if !c.IsFromProposer(c.Round(), proposal.Sender()) {
		c.logger.Warn("Ignore proposal messages from non-proposer")
		return constants.ErrNotFromProposer
	}

	// received a current round proposal
	if metrics.Enabled {
		now := time.Now()
		ProposalReceivedTimer.Update(now.Sub(c.newRound))
		ProposalReceivedBg.Add(now.Sub(c.newRound).Nanoseconds())
	}

	// Verify the proposal we received
	start := time.Now()
	duration, err := c.backend.VerifyProposal(proposal.Block()) // youssef: can we skip the verification for our own proposal?

	if metrics.Enabled {
		now := time.Now()
		ProposalVerifiedTimer.Update(now.Sub(start))
		ProposalVerifiedBg.Add(now.Sub(start).Nanoseconds())
	}

	if err != nil {
		// if it's a future block, we will handle it again after the duration
		// TODO: implement wiggle time / median time
		if errors.Is(err, consensus.ErrFutureTimestampBlock) {
			c.StopFutureProposalTimer()
			c.futureProposalTimer = time.AfterFunc(duration, func() {
				c.SendEvent(backlogMessageEvent{
					msg: proposal,
				})
			})
			return err
		}
		// Proposal is invalid here, we need to prevote nil.
		// However, we may have already sent a prevote nil in the past without having processed the proposal
		// because of a timeout, so we need to check if we are still in the Propose step.
		if c.step == Propose {
			c.prevoter.SendPrevote(ctx, true)
			// do not to accept another proposal in current round
			c.SetStep(ctx, Prevote)
		}

		c.logger.Warn("Failed to verify proposal", "err", err, "duration", duration)

		return err
	}

	// Set the proposal for the current round
	c.curRoundMessages.SetProposal(proposal, true)
	c.LogProposalMessageEvent("MessageEvent(Proposal): Received", proposal, proposal.Sender().String(), c.address.String())

	//l49: Check if we have a quorum of precommits for this proposal
	hash := proposal.Block().Hash()
	if c.curRoundMessages.PrecommitsPower(hash).Cmp(c.CommitteeSet().Quorum()) >= 0 {
		c.Commit(ctx, proposal.R(), c.curRoundMessages)
		return nil
	}

	if c.step == Propose {
		vr := proposal.ValidRound()
		// Line 22 in Algorithm 1 of The latest gossip on BFT consensus
		if vr == -1 {
			// When lockedRound is set to any value other than -1 lockedValue is also
			// set to a non nil value. So we can be sure that we will only try to access
			// lockedValue when it is non nil.
			c.prevoter.SendPrevote(ctx, !(c.lockedRound == -1 || hash == c.lockedValue.Hash()))
			c.SetStep(ctx, Prevote)
			return nil
		}
		rs := c.messages.GetOrCreate(vr)
		// Line 28 in Algorithm 1 of The latest gossip on BFT consensus
		// vr >= 0 here
		if vr < c.Round() && rs.PrevotesPower(hash).Cmp(c.CommitteeSet().Quorum()) >= 0 {
			c.prevoter.SendPrevote(ctx, !(c.lockedRound <= vr || hash == c.lockedValue.Hash()))
			c.SetStep(ctx, Prevote)
		}
	}

	// Line 36 in Algorithm 1 of The latest gossip on BFT consensus
	if c.step >= Prevote && c.curRoundMessages.PrevotesPower(proposal.Block().Hash()).Cmp(c.CommitteeSet().Quorum()) >= 0 && !c.setValidRoundAndValue {
		if c.step == Prevote {
			c.lockedValue = proposal.Block()
			c.lockedRound = c.Round()
			c.precommiter.SendPrecommit(ctx, false)
			c.SetStep(ctx, Precommit)
		}
		c.validValue = proposal.Block()
		c.validRound = c.Round()
		c.setValidRoundAndValue = true
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
	if c.IsProposer() && c.step == Propose && !c.sentProposal && c.Height().Cmp(number) == 0 {
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
		"msgHeight", proposal.H(),
		"currentRound", c.Round(),
		"msgRound", proposal.R(),
		"currentStep", c.step,
		"isProposer", c.IsProposer(),
		"currentProposer", c.CommitteeSet().GetProposer(c.Round()),
		"isNilMsg", proposal.Block().Hash() == common.Hash{},
		"hash", proposal.Block().Hash(),
	)
}
