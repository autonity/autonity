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
	"github.com/clearmatics/autonity/common"
	"time"

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

		if proposalBlock == nil {
			logger.Error("send nil proposed block",
				"Round", c.currentRoundState.round.String(), "Height",
				c.currentRoundState.height.String(), "ValidRound", c.validRound)

			return
		}

		c.sentProposal = true
		c.backend.SetProposedBlockHash(p.Hash())

		c.logProposalMessageEvent("MessageEvent(Proposal): Sent", *proposalBlock, c.address.String(), "broadcast")

		c.broadcast(ctx, &message{
			Code:          msgProposal,
			Msg:           proposal,
			Address:       c.address,
			CommittedSeal: []byte{},
		})
	}
}

func (c *core) handleProposal(ctx context.Context, msg *message) error {
	var proposal Proposal
	err := msg.Decode(&proposal)
	if err != nil {
		return errFailedDecodeProposal
	}

	oldRoundProposal := false
	valSet := c.valSet
	state := c.currentRoundState
	// Ensure we have the same view with the Proposal message
	if err := c.checkMessage(proposal.Round, proposal.Height); err != nil {
		if err == errOldRoundMessage {
			c.logger.Warn("Old round propose message received","round",proposal.Round.Uint64())
			state = c.getOrCreateOldRoundState(proposal.Round)
			//if we already had a proposal for this round, nothing should happen
			//this should never happen in presence of honest validators
			if state.proposal != nil {
				c.logger.Warn("Processing already received propose message for this round")
				return err
			}
			// check if the proposal was coming from the correct old round proposer
			valSet = new(validatorSet)
			valSet.set(c.valSet.Copy())
			valSet.CalcProposer(c.lastProposer, proposal.Round.Uint64())
			oldRoundProposal = true
		} else {
			return err
		}
	}

	// Check if the message comes from currentRoundState proposer
	if !valSet.IsProposer(msg.Address) {
		c.logger.Warn("Ignore proposal messages from non-proposer")
		return errNotFromProposer
	}

	// Verify the proposal we received
	if duration, err := c.backend.VerifyProposal(*proposal.ProposalBlock); err != nil {
		if oldRoundProposal {
			return err
		}
		if timeoutErr := c.proposeTimeout.stopTimer(); timeoutErr != nil {
			return timeoutErr
		}
		c.logger.Debug("Stopped Scheduled Proposal Timeout")
		c.sendPrevote(ctx, true)
		// do not to accept another proposal in current round
		c.setStep(prevote)

		c.logger.Warn("Failed to verify proposal", "err", err, "duration", duration)
		// if it's a future block, we will handle it again after the duration
		// TIME FIELD OF HEADER CHECKED HERE - NOT HEIGHT
		// TODO: implement wiggle time / median time
		if err == consensus.ErrFutureBlock {
			c.stopFutureProposalTimer()
			c.futureProposalTimer = time.AfterFunc(duration, func() {
				_, sender := valSet.GetByAddress(msg.Address)
				c.sendEvent(backlogEvent{
					src: sender,
					msg: msg,
				})
			})
		}
		return err
	}

	if oldRoundProposal {
		state.SetProposal(&proposal, msg)
		if c.CanDecide(state) {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				c.commit(state)
			}
			return nil
		}
	}
	// Here is about to accept the Proposal
	if c.currentRoundState.Step() == propose {
		if err := c.proposeTimeout.stopTimer(); err != nil {
			return err
		}
		c.logger.Debug("Stopped Scheduled Proposal Timeout")

		// Set the proposal for the current round
		c.currentRoundState.SetProposal(&proposal, msg)

		c.logProposalMessageEvent("MessageEvent(Proposal): Received", proposal, msg.Address.String(), c.address.String())

		vr := proposal.ValidRound.Int64()
		h := proposal.ProposalBlock.Hash()
		curR := c.currentRoundState.Round().Int64()

		if vr == -1 {
			// Line 22 in Algorithm 1 of The latest gossip on BFT consensus
			if c.lockedRound.Int64() == vr || h == c.lockedValue.Hash() {
				c.sendPrevote(ctx, false)
			} else {
				c.sendPrevote(ctx, true)
			}
			c.setStep(prevote)
			// Line 28 in Algorithm 1 of The latest gossip on BFT consensus
		} else if rs, ok := c.currentHeightOldRoundsStates[vr]; vr > -1 && vr < curR && ok && c.Quorum(rs.Prevotes.VotesSize(h)) {
			if c.lockedRound.Int64() <= vr || h == c.lockedValue.Hash() {
				c.sendPrevote(ctx, false)
			} else {
				c.sendPrevote(ctx, true)
			}
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
