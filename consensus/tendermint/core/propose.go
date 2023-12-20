package core

import (
	"context"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	tctypes "github.com/autonity/autonity/consensus/tendermint/core/types"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/metrics"
	"github.com/autonity/autonity/rlp"
)

type Proposer struct {
	*Core
}

func (c *Proposer) SendProposal(ctx context.Context, p *types.Block) {
	logger := c.logger.New("step", c.step)

	// If I'm the proposer and I have the same height with the proposal
	if c.Height().Cmp(p.Number()) == 0 && c.IsProposer() && !c.sentProposal {
		proposalBlock := message.NewProposal(c.Round(), c.Height(), c.validRound, p, c.backend.Sign)
		proposal, err := rlp.EncodeToBytes(proposalBlock)
		if err != nil {
			logger.Error("Failed to encode", "Round", proposalBlock.Round, "Height", proposalBlock.Height, "ValidRound", c.validRound)
			return
		}

		c.sentProposal = true
		c.backend.SetProposedBlockHash(p.Hash())

		if metrics.Enabled {
			now := time.Now()
			tctypes.ProposalSentTimer.Update(now.Sub(c.newRound))
			tctypes.ProposalSentBg.Add(now.Sub(c.newRound).Nanoseconds())
		}
		c.LogProposalMessageEvent("MessageEvent(Proposal): Sent", proposalBlock, c.address.String(), "broadcast")

		c.Br().SignAndBroadcast(ctx, &message.Message{
			Code:          consensus.MsgProposal,
			Payload:       proposal,
			Address:       c.address,
			CommittedSeal: []byte{},
		})
	}
}

func (c *Proposer) HandleProposal(ctx context.Context, msg *message.Message) error {
	proposal := msg.ConsensusMsg.(*message.Proposal)
	// Ensure we have the same view with the Proposal message
	if err := c.CheckMessage(proposal.Round, proposal.Height.Uint64(), tctypes.Propose); err != nil {
		// If it's a future round proposal, the only upon condition
		// that can be triggered is L49, but this requires more than F future round messages
		// meaning that a future roundchange will happen before, as such, pushing the
		// message to the backlog is fine.
		if err == constants.ErrOldRoundMessage {

			roundMsgs := c.messages.GetOrCreate(proposal.Round)

			// if we already have a proposal for this old round - ignore
			// the first proposal sent by the sender in a round is always the only one we consider.
			if roundMsgs.ProposalMsg != nil {
				return constants.ErrAlreadyProcessed
			}

			if !c.IsFromProposer(proposal.Round, msg.Address) {
				c.logger.Warn("Ignoring proposal from non-proposer")
				return constants.ErrNotFromProposer
			}

			// Save it, but do not verify the proposal yet unless we have enough precommits for it.
			roundMsgs.SetProposal(proposal, msg, false)

			if roundMsgs.PrecommitsPower(roundMsgs.GetProposalHash()).Cmp(c.CommitteeSet().Quorum()) >= 0 {
				if _, err2 := c.backend.VerifyProposal(proposal.ProposalBlock); err2 != nil {
					return err2
				}
				c.logger.Debug("Committing old round proposal")
				c.Commit(proposal.Round, roundMsgs)
				return nil
			}
		}
		return err
	}
	// if we already have processed a proposal in this round we ignore.
	if c.curRoundMessages.ProposalMsg != nil {
		return constants.ErrAlreadyProcessed
	}
	// At this point the local round matches the message round and the current step
	// could be either Proposal, Prevote, or Precommit.

	// Check if the message comes from curRoundMessages proposer
	if !c.IsFromProposer(c.Round(), msg.Address) {
		c.logger.Warn("Ignore proposal messages from non-proposer")
		return constants.ErrNotFromProposer
	}

	// received a current round proposal
	if metrics.Enabled {
		now := time.Now()
		tctypes.ProposalReceivedTimer.Update(now.Sub(c.newRound))
		tctypes.ProposalReceivedBg.Add(now.Sub(c.newRound).Nanoseconds())
	}

	// Verify the proposal we received
	start := time.Now()
	duration, err := c.backend.VerifyProposal(proposal.ProposalBlock)

	if metrics.Enabled {
		now := time.Now()
		tctypes.ProposalVerifiedTimer.Update(now.Sub(start))
		tctypes.ProposalVerifiedBg.Add(now.Sub(start).Nanoseconds())
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
					msg: msg,
				})
			})
			return err
		}
		// Proposal is invalid here, we need to prevote nil.
		// However, we may have already sent a prevote nil in the past without having processed the proposal
		// because of a timeout, so we need to check if we are still in the Propose step.
		if c.step == tctypes.Propose {
			c.prevoter.SendPrevote(ctx, true)
			// do not to accept another proposal in current round
			c.SetStep(tctypes.Prevote)
		}
		c.logger.Warn("Failed to verify proposal", "err", err, "duration", duration)

		return err
	}

	// Set the proposal for the current round
	c.curRoundMessages.SetProposal(proposal, msg, true)

	c.LogProposalMessageEvent("MessageEvent(Proposal): Received", proposal, msg.Address.String(), c.address.String())

	//l49: Check if we have a quorum of precommits for this proposal
	curProposalHash := c.curRoundMessages.GetProposalHash()
	if c.curRoundMessages.PrecommitsPower(curProposalHash).Cmp(c.CommitteeSet().Quorum()) >= 0 {
		c.Commit(proposal.Round, c.curRoundMessages)
		return nil
	}

	if c.step == tctypes.Propose {
		if err := c.proposeTimeout.StopTimer(); err != nil {
			return err
		}

		vr := proposal.ValidRound
		h := proposal.ProposalBlock.Hash()

		// Line 22 in Algorithm 1 of The latest gossip on BFT consensus
		if vr == -1 {
			// When lockedRound is set to any value other than -1 lockedValue is also
			// set to a non nil value. So we can be sure that we will only try to access
			// lockedValue when it is non nil.
			c.prevoter.SendPrevote(ctx, !(c.lockedRound == -1 || h == c.lockedValue.Hash()))
			c.SetStep(tctypes.Prevote)
			return nil
		}

		rs := c.messages.GetOrCreate(vr)

		// Line 28 in Algorithm 1 of The latest gossip on BFT consensus
		// vr >= 0 here
		if vr < c.Round() && rs.PrevotesPower(h).Cmp(c.CommitteeSet().Quorum()) >= 0 {
			c.prevoter.SendPrevote(ctx, !(c.lockedRound <= vr || h == c.lockedValue.Hash()))
			c.SetStep(tctypes.Prevote)
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

func (c *Proposer) LogProposalMessageEvent(message string, proposal *message.Proposal, from, to string) {
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
