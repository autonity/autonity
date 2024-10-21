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
	"github.com/autonity/autonity/log"
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

	// We start preparing block as soon as proposal is verified, but there are situation
	// that verified proposal is not finalized in the particular round hence this safety
	// check to ensure that the block parent hash is same as last hash in core
	if c.Backend().HeadBlock().Hash() != block.ParentHash() {
		log.Info("verified proposal was not finalized in the last round", "aborting send proposal", "last header hash", c.Backend().HeadBlock().Hash(), "block parent hash", block.ParentHash())
		return
	}

	if c.sentProposal {
		return
	}
	self, err := c.CommitteeSet().MemberByAddress(c.address)
	if err != nil {
		// it can happen in edge case addressed by docker e2e test, that is a validator resets at the epoch boundary,
		// after which it leaves the committee, we cannot panic it in that case.
		c.logger.Error("Validator is no longer in current committee", "err", err, "validator", c.address.String())
		return
	}

	proposal := message.NewPropose(c.Round(), c.Height().Uint64(), c.validRound, block, c.backend.Sign, self)
	c.sentProposal = true
	c.backend.SetProposedBlockHash(block.Hash())
	c.LogProposalMessageEvent("MessageEvent(Proposal): Sent", proposal)
	c.Broadcaster().Broadcast(proposal)
	if metrics.Enabled {
		now := time.Now()
		ProposalSentTimer.Update(now.Sub(c.newRound))
		c.currBlockTimeStamp = time.Unix(int64(proposal.Block().Header().Time), 0)
		ProposalSentBlockTSDeltaBg.Add(time.Since(c.currBlockTimeStamp).Nanoseconds())
	}
}

func (c *Proposer) HandleProposal(ctx context.Context, proposal *message.Propose) error {
	if !proposal.PreVerified() || !proposal.Verified() {
		panic("Handling NON cryptographically verified proposal")
	}

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
	if !c.IsFromProposer(proposal.R(), proposal.Signer()) {
		c.logger.Warn("Ignoring proposal from non-proposer", "signer", proposal.Signer())
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
		c.currBlockTimeStamp = time.Unix(int64(proposal.Block().Header().Time), 0)
		ProposalReceivedBlockTSDeltaBg.Add(time.Since(c.currBlockTimeStamp).Nanoseconds())
	}

	var (
		duration time.Duration
		err      error
		start    = time.Now()
	)

	// skip verification for our own proposal
	// skip if our own OR cached
	// verify if not in our own and not state cached
	if c.backend.ProposedBlockHash() != proposal.Block().Hash() && !c.backend.IsProposalStateCached(proposal.Block().Hash()) {
		duration, err = c.backend.VerifyProposal(proposal.Block())
	}

	if metrics.Enabled {
		now := time.Now()
		ProposalVerifiedTimer.Update(now.Sub(start))
		ProposalVerifiedBg.Add(now.Sub(start).Nanoseconds())
	}

	if err != nil {
		// if it's a future block, we will handle it again after the duration
		// TODO: implement wiggle time / median time
		if errors.Is(err, consensus.ErrFutureTimestampBlock) {
			c.logger.Debug("delaying processing of proposal due to future timestamp", "delay", duration)
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

	// notify miner
	go c.Backend().ProposalVerified(proposal.Block())
	// Set the proposal for the current round
	c.curRoundMessages.SetProposal(proposal, true)
	c.LogProposalMessageEvent("MessageEvent(Proposal): Received", proposal)

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

func (c *Proposer) LogProposalMessageEvent(message string, proposal *message.Propose) {
	c.logger.Debug(message,
		"type", "Proposal",
		"local address", log.Lazy{Fn: func() string { return c.Address().String() }},
		"currentHeight", log.Lazy{Fn: c.Height},
		"msgHeight", proposal.H(),
		"currentRound", log.Lazy{Fn: c.Round},
		"msgRound", proposal.R(),
		"currentStep", c.step,
		"isProposer", log.Lazy{Fn: c.IsProposer},
		"currentProposer", log.Lazy{Fn: func() *types.CommitteeMember { return c.CommitteeSet().GetProposer(c.Round()) }},
		"isNilMsg", log.Lazy{Fn: func() bool { return proposal.Block().Hash() == common.Hash{} }},
		"value", log.Lazy{Fn: func() common.Hash { return proposal.Block().Hash() }},
		"proposal", log.Lazy{Fn: func() string { return proposal.String() }},
	)
}
