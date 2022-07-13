package core

import (
	"context"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	tctypes "github.com/autonity/autonity/consensus/tendermint/core/types"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core/types"
)

type ProposeService struct {
	*Core
}

func (c *ProposeService) SendProposal(ctx context.Context, p *types.Block) {
	logger := c.logger.New("step", c.step)

	// If I'm the proposer and I have the same height with the proposal
	if c.Height().Cmp(p.Number()) == 0 && c.IsProposer() && !c.sentProposal {
		proposalBlock := messageutils.NewProposal(c.Round(), c.Height(), c.validRound, p)
		proposal, err := messageutils.Encode(proposalBlock)
		if err != nil {
			logger.Error("Failed to encode", "Round", proposalBlock.Round, "Height", proposalBlock.Height, "ValidRound", c.validRound)
			return
		}

		c.sentProposal = true
		c.backend.SetProposedBlockHash(p.Hash())

		c.LogProposalMessageEvent("MessageEvent(Proposal): Sent", *proposalBlock, c.address.String(), "broadcast")

		c.Br().Broadcast(ctx, &messageutils.Message{
			Code:          messageutils.MsgProposal,
			Msg:           proposal,
			Address:       c.address,
			CommittedSeal: []byte{},
		})
	}
}

func (c *ProposeService) HandleProposal(ctx context.Context, msg *messageutils.Message) error {
	var proposal messageutils.Proposal
	err := msg.Decode(&proposal)
	if err != nil {
		return constants.ErrFailedDecodeProposal
	}

	// Ensure we have the same view with the Proposal message
	if err := c.CheckMessage(proposal.Round, proposal.Height, tctypes.Propose); err != nil {
		// If it's a future round proposal, the only upon condition
		// that can be triggered is L49, but this requires more than F future round messages
		// meaning that a future roundchange will happen before, as such, pushing the
		// message to the backlog is fine.
		if err == constants.ErrOldRoundMessage {

			roundMsgs := c.messages.GetOrCreate(proposal.Round)

			// if we already have a proposal then it must be different than the current one
			// it can't happen unless someone's byzantine.
			if roundMsgs.ProposalDetails != nil {
				return err // do not gossip, TODO: accountability
			}

			if !c.IsProposerMsg(proposal.Round, msg.Address) {
				c.logger.Warn("Ignore proposal messages from non-proposer")
				return constants.ErrNotFromProposer
			}
			// We do not verify the proposal in this case.
			roundMsgs.SetProposal(&proposal, msg, false)

			if roundMsgs.PrecommitsPower(roundMsgs.GetProposalHash()) >= c.CommitteeSet().Quorum() {
				if _, error := c.backend.VerifyProposal(*proposal.ProposalBlock); error != nil {
					return error
				}
				c.logger.Debug("Committing old round proposal")
				c.Commit(proposal.Round, roundMsgs)
				return nil
			}
		}
		return err
	}

	// Check if the message comes from curRoundMessages proposer
	if !c.IsProposerMsg(c.Round(), msg.Address) {
		c.logger.Warn("Ignore proposal messages from non-proposer")
		return constants.ErrNotFromProposer
	}

	// Verify the proposal we received
	if duration, err := c.backend.VerifyProposal(*proposal.ProposalBlock); err != nil {

		if timeoutErr := c.proposeTimeout.StopTimer(); timeoutErr != nil {
			return timeoutErr
		}
		// if it's a future block, we will handle it again after the duration
		// TODO: implement wiggle time / median time
		if err == consensus.ErrFutureBlock {
			c.StopFutureProposalTimer()
			c.futureProposalTimer = time.AfterFunc(duration, func() {
				c.SendEvent(backlogEvent{
					msg: msg,
				})
			})
			return err
		}
		c.prevoter.SendPrevote(ctx, true)
		// do not to accept another proposal in current round
		c.SetStep(tctypes.Prevote)

		c.logger.Warn("Failed to verify proposal", "err", err, "duration", duration)

		return err
	}

	// Set the proposal for the current round
	c.curRoundMessages.SetProposal(&proposal, msg, true)

	c.LogProposalMessageEvent("MessageEvent(Proposal): Received", proposal, msg.Address.String(), c.address.String())

	//l49: Check if we have a quorum of precommits for this proposal
	curProposalHash := c.curRoundMessages.GetProposalHash()
	if c.curRoundMessages.PrecommitsPower(curProposalHash) >= c.CommitteeSet().Quorum() {
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
		if vr < c.Round() && rs.PrevotesPower(h) >= c.CommitteeSet().Quorum() {
			c.prevoter.SendPrevote(ctx, !(c.lockedRound <= vr || h == c.lockedValue.Hash()))
			c.SetStep(tctypes.Prevote)
		}
	}

	return nil
}

func (c *ProposeService) HandleNewCandidateBlockMsg(ctx context.Context, candidateBlock *types.Block) {
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

func (c *ProposeService) StopFutureProposalTimer() {
	if c.futureProposalTimer != nil {
		c.futureProposalTimer.Stop()
	}
}

func (c *ProposeService) LogProposalMessageEvent(message string, proposal messageutils.Proposal, from, to string) {
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
