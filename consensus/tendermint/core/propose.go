// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"context"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core/types"
)

func (c *core) sendProposal(ctx context.Context, p *types.Block) {
	logger := c.logger.New("step", c.currentRoundState.Step())

	// If I'm the proposer and I have the same height with the proposal
	if c.currentRoundState.Height().Int64() == p.Number().Int64() && c.isProposer() && !c.sentProposal {
		proposalBlock := NewProposal(c.currentRoundState.Round(), c.currentRoundState.Height(), c.validRound, p)
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
		// If it's a future round proposal, the only upon conditon
		// that can be triggered is L49, but this requires more than F future round messages
		// meaning that a future roundchange will happen before, as such, pushing the
		// message to the backlog is fine.
		if err == errOldRoundMessage {
			// if we already have a proposal then abort.
		}
		return err
	}

	// Check if the message comes from curRoundMessages proposer
	if !c.CommitteeSet().IsProposer(c.Round(), msg.Address) {
		c.logger.Warn("Ignore proposal messages from non-proposer")
		return errNotFromProposer
	}

	// Verify the proposal we received
	if duration, err := c.backend.VerifyProposal(*proposal.ProposalBlock); err != nil {
		if timeoutErr := c.proposeTimeout.stopTimer(); timeoutErr != nil {
			return timeoutErr
		}
		c.sendPrevote(ctx, true)
		// do not to accept another proposal in current round
		c.setStep(prevote)

		c.logger.Warn("Failed to verify proposal", "err", err, "duration", duration)
		// if it's a future block, we will handle it again after the duration
		// TODO: implement wiggle time / median time
		if err == consensus.ErrFutureBlock {
			c.stopFutureProposalTimer()
			c.futureProposalTimer = time.AfterFunc(duration, func() {
				_, sender := c.valSet.GetByAddress(msg.Address)
				c.sendEvent(backlogEvent{
					src: sender,
					msg: msg,
				})
			})
		}
		return err
	}

	// Here is about to accept the Proposal
	if c.currentRoundState.Step() == propose {
		if err := c.proposeTimeout.stopTimer(); err != nil {
			return err
		}

		// Set the proposal for the current round
		c.currentRoundState.SetProposal(&proposal, msg)

		c.logProposalMessageEvent("MessageEvent(Proposal): Received", proposal, msg.Address.String(), c.address.String())

		vr := proposal.ValidRound.Int64()
		h := proposal.ProposalBlock.Hash()
		curR := c.currentRoundState.Round().Int64()

		c.currentHeightOldRoundsStatesMu.RLock()
		defer c.currentHeightOldRoundsStatesMu.RUnlock()

		// Line 22 in Algorithm 1 of The latest gossip on BFT consensus
		if vr == -1 {
			var voteForProposal = false
			if c.lockedValue != nil {
				voteForProposal = c.lockedRound.Int64() == -1 || h == c.lockedValue.Hash()
			}
			c.sendPrevote(ctx, voteForProposal)
			c.setStep(prevote)
			return nil
		}

		rs, ok := c.currentHeightOldRoundsStates[vr]
		if !ok {
			c.logger.Error("handleProposal. unknown old round",
				"proposalHeight", h,
				"proposalRound", vr,
				"currentHeight", c.currentRoundState.height.Uint64(),
				"currentRound", c.currentRoundState.round.Uint64(),
			)
		}

		// Line 28 in Algorithm 1 of The latest gossip on BFT consensus
		if ok && vr < curR && c.Quorum(rs.Prevotes.VotesSize(h)) {
			var voteForProposal = false
			if c.lockedValue != nil {
				voteForProposal = c.lockedRound.Int64() <= vr || h == c.lockedValue.Hash()
			}
			c.sendPrevote(ctx, voteForProposal)
			c.setStep(prevote)
		}
	}

	return nil
}

func (c *core) logProposalMessageEvent(message string, proposal Proposal, from, to string) {
	c.logger.Debug(message,
		"type", "Proposal",
		"from", from,
		"to", to,
		"currentHeight", c.currentRoundState.Height(),
		"msgHeight", proposal.Height,
		"currentRound", c.currentRoundState.Round(),
		"msgRound", proposal.Round,
		"currentStep", c.currentRoundState.Step(),
		"isProposer", c.isProposer(),
		"currentProposer", c.valSet.GetProposer(),
		"isNilMsg", proposal.ProposalBlock.Hash() == common.Hash{},
		"hash", proposal.ProposalBlock.Hash(),
	)
}
