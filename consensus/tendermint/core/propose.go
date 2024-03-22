package core

import (
	"context"
	"errors"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/metrics"
)

type Proposer struct {
	*Core
}

func (c *Proposer) SendProposal(_ context.Context, block *types.Block) {
	// Required to have proposal block current and being current proposer
	// Defensively panic here to catch bugs
	if c.Height().Cmp(block.Number()) != 0 {
		panic("proposal block height incorrect")
	}
	if !c.IsProposer() {
		panic("not proposer")
	}
	if c.sentProposal {
		return
	}
	proposal := message.NewPropose(c.Round(), c.Height().Uint64(), c.validRound, block, c.backend.Sign)
	c.sentProposal = true
	c.backend.SetProposedBlockHash(block.Hash())
	c.logger.Info("Proposing new block", "proposal", proposal.Block().Hash(), "round", c.Round(), "height", c.Height().Uint64())
	c.LogProposalMessageEvent("MessageEvent(Proposal): Sent", proposal, c.address.String(), "broadcast")
	c.Broadcaster().Broadcast(proposal)
	if metrics.Enabled {
		now := time.Now()
		ProposalSentTimer.Update(now.Sub(c.newRound))
		ProposalSentBg.Add(now.Sub(c.newRound).Nanoseconds())
		ProposalSentBlockTSDeltaBg.Add(time.Unix(int64(block.Header().Time), 0).Sub(now).Nanoseconds())
		c.proposalSent = now
	}
}

func (c *Proposer) HandleProposal(ctx context.Context, proposal *message.Propose) error {
	if proposal.R() > c.Round() {
		// If it's a future round proposal, the only upon condition
		// that can be triggered is L49, but this requires more than F future round messages
		// meaning that a future roundchange will happen before, as such, pushing the
		// message to the backlog is fine.
		return constants.ErrFutureRoundMessage
	}

	// proposal is either for current round or old round
	roundMessages := c.messages.GetOrCreate(proposal.R())

	// if we already have a proposal for this round - ignore
	// the first proposal sent by the sender in a round is always the only one we consider.
	if roundMessages.Proposal() != nil {
		return constants.ErrAlreadyHaveProposal
	}

	// check if proposal comes from the correct proposer for pair (h,r)
	if !c.IsFromProposer(proposal.R(), proposal.Sender()) {
		c.logger.Warn("Ignoring proposal from non-proposer", "sender", proposal.Sender())
		return constants.ErrNotFromProposer
	}

	if proposal.R() < c.Round() {
		// old round proposal, check if we have quorum precommits on it
		// Save it, but do not verify the proposal yet unless we have enough precommits for it.
		roundMessages.SetProposal(proposal, false)

		// Line 49 in Algorithm 1 of The latest gossip on BFT consensus
		// check if we have a quorum of precommits for this proposal
		_ = c.quorumPrecommitsCheck(ctx, proposal, false)

		return constants.ErrOldRoundMessage
	}

	// At this point the local round matches the message round (roundMessages == c.curRoundMessages)
	// current step could be either Proposal, Prevote, or Precommit.

	// received a current round proposal
	if metrics.Enabled {
		now := time.Now()
		ProposalReceivedTimer.Update(now.Sub(c.newRound))
		ProposalReceivedBg.Add(now.Sub(c.newRound).Nanoseconds())
		ProposalReceivedBlockTSDeltaBg.Add(time.Since(c.proposalSent).Nanoseconds())
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
		// if the proposal block is already in the chain, no need to prevote for nil
		if errors.Is(err, core.ErrKnownBlock) || errors.Is(err, constants.ErrAlreadyHaveBlock) {
			c.logger.Info("Verified proposal that was already in our local chain", "err", err, "duration", duration)
			c.SetStep(ctx, PrecommitDone) // we do not need to process any more consensus messages for this height
			return constants.ErrAlreadyHaveBlock
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

	// check upon conditions for current round proposal
	c.currentProposalChecks(ctx, proposal)

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
