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
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core/types"
)

func (c *core) sendProposal(ctx context.Context, p *types.Block) {
	logger := c.logger.New("step", c.getStep())

	// If I'm the proposer and I have the same height with the proposal
	if c.getHeight().Int64() == p.Number().Int64() && c.isProposerForR(c.getRound().Int64(), c.address) && !c.sentProposal {
		proposalBlock := NewProposal(c.getRound(), c.getHeight(), c.validRound, p, c.logger)
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

	// Check for nil values
	if proposal.Height == nil || proposal.Round == nil || proposal.ProposalBlock == nil {
		return errInvalidMessage
	}

	// Ensure proposal is for current height
	if err := c.checkMessage(proposal.Round, proposal.Height, propose); err != nil {
		return err
	}

	// Ensure proposer of the Proposal is for round proposal.Round
	if !c.isProposerForR(proposal.Round.Int64(), msg.Address) {
		c.logger.Warn("Ignore proposal messages from non-proposer")
		return errNotFromProposer
	}

	// Check block header (will check body later)
	if duration, err := c.backend.VerifyProposalHeader(*proposal.ProposalBlock); err != nil {
		if err == consensus.ErrFutureBlock {
			c.logger.Warn("Proposal timestamp greater than local time. Setting timer to handle message again.", "err", err, "duration", duration)
			c.stopFutureProposalTimer()
			c.futureProposalTimer = time.AfterFunc(duration, func() {
				p, _ := msg.Payload()
				c.sendEvent(events.MessageEvent{Payload: p})
			})
		}
		return err
	}

	// if the proposal is different from what is stored in the round state, then proposer is byzantine, therefore,
	// ignore the proposal. otherwise, if the proposal block is received more than once through gossip we need to
	// ignore since the state will not change.
	if ps := c.getProposalSet(proposal.Round.Int64()); ps != nil {
		c.logger.Debug("Already have proposal so ignoring...")
		return nil
	}
	// We don't have old, current or future proposal, then add the proposal to the relevant round message set
	// and since the state has changed we need to check for consensus on this proposal block.
	c.setProposalSet(proposal.Round.Int64(), proposal, msg)

	c.logProposalMessageEvent("MessageEvent(Proposal): Received", proposal, msg.Address.String(), c.address.String())

	if err := c.checkForConsensus(ctx, proposal.Round.Int64()); err != nil {
		return err
	}

	roundCmp := proposal.Round.Cmp(c.getRound())
	if roundCmp < 0 {
		// Nothing more to do for old round proposal
		c.logger.Debug("Received old round proposal")
	} else if roundCmp > 0 {
		c.logger.Debug("Received future round proposal")
		c.checkForFutureRoundChange(ctx, proposal.Round.Int64())
	} else {
		// Proposal is for current round, i.e. proposal.Round.Int64() = c.getRound().Int64()
		if c.getStep() == propose {
			if proposal.ValidRound.Int64() == -1 {
				return c.checkForNewProposal(ctx, proposal.Round.Int64())
			} else if proposal.ValidRound.Int64() >= 0 {
				return c.checkForOldProposal(ctx, proposal.Round.Int64())
			}
		} else if c.getStep() >= prevote {
			return c.checkForQuorumPrevotes(ctx, proposal.Round.Int64())
		}
	}

	return nil
}

func (c *core) logProposalMessageEvent(message string, proposal Proposal, from, to string) {
	c.logger.Debug(message,
		"type", "Proposal",
		"from", from,
		"to", to,
		"currentHeight", c.getHeight(),
		"msgHeight", proposal.Height,
		"currentRound", c.getRound(),
		"msgRound", proposal.Round,
		"currentStep", c.getStep(),
		"isProposer", c.isProposerForR(proposal.Round.Int64(), c.address),
		"currentProposer", c.valSet.GetProposer(),
		"isNilMsg", proposal.ProposalBlock.Hash() == common.Hash{},
		"hash", proposal.ProposalBlock.Hash(),
	)
}
