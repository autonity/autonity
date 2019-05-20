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
	"time"

	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint"
)

// TODO: add new message struct for proposal (proposalMessage) and determine how to rlp encode them especially nil
// TODO: add new message for vote (prevote and precommit) and determine how to rlp encode them especially nil
func (c *core) sendProposal(request *tendermint.Request) {
	logger := c.logger.New("state", c.state)

	// If I'm the proposer and I have the same sequence with the proposal
	if c.current.Sequence().Cmp(request.ProposalBlock.Number()) == 0 && c.isProposer() && !c.sentProposal {
		curView := c.currentView()
		proposal, err := Encode(&tendermint.Proposal{
			View:          curView,
			ProposalBlock: request.ProposalBlock,
		})
		if err != nil {
			logger.Error("Failed to encode", "view", curView)
			return
		}
		c.sentProposal = true
		c.backend.SetProposedBlockHash(request.ProposalBlock.Hash())
		c.broadcast(&message{
			Code: msgProposal,
			Msg:  proposal,
		})
	}
}

func (c *core) handleProposal(msg *message, src tendermint.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Decode PRE-PREPARE
	var proposal *tendermint.Proposal
	err := msg.Decode(&proposal)
	if err != nil {
		return errFailedDecodeProposal
	}

	// Ensure we have the same view with the PRE-PREPARE message
	// If it is old message, see if we need to broadcast COMMIT
	if err := c.checkMessage(msgProposal, proposal.View); err != nil {
		if err == errOldMessage {
			//TODO : EIP Says ignore?
			// Get validator set for the given proposal
			valSet := c.backend.Validators(proposal.ProposalBlock.Number().Uint64()).Copy()
			previousProposer := c.backend.GetProposer(proposal.ProposalBlock.Number().Uint64() - 1)
			valSet.CalcProposer(previousProposer, proposal.View.Round.Uint64())
			// Broadcast COMMIT if it is an existing block
			// 1. The proposer needs to be a proposer matches the given (Sequence + Round)
			// 2. The given block must exist
			if valSet.IsProposer(src.Address()) && c.backend.HasPropsal(proposal.ProposalBlock.Hash(), proposal.ProposalBlock.Number()) {
				c.sendPrecommitForOldBlock(proposal.View, proposal.ProposalBlock.Hash())
				return nil
			}
		}
		return err
	}

	// Check if the message comes from current proposer
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
		} else {
			c.sendNextRoundChange()
		}
		return err
	}

	// Here is about to accept the PRE-PREPARE
	if c.state == StateAcceptRequest {
		// Send ROUND CHANGE if the locked proposal and the received proposal are different
		if c.current.IsHashLocked() {
			if proposal.ProposalBlock.Hash() == c.current.GetLockedHash() {
				// Broadcast COMMIT and enters Prevoted state directly
				c.acceptProposal(proposal)
				c.setState(StatePrevoteDone)
				c.sendPrecommit() // TODO : double check, why not PREPARE?
			} else {
				// Send round change
				c.sendNextRoundChange()
			}
		} else {
			// Either
			//   1. the locked proposal and the received proposal match
			//   2. we have no locked proposal
			c.acceptProposal(proposal)
			c.setState(StateProposeDone)
			c.sendPrevote()
		}
	}

	return nil
}

func (c *core) acceptProposal(proposal *tendermint.Proposal) {
	c.current.SetProposal(proposal)
}
