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
	"bytes"
	"context"
	"math/big"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
)

func (c *core) sendPrecommit(ctx context.Context, isNil bool) {
	logger := c.logger.New("step", c.step)

	var precommit = Vote{
		Round:  c.Round(),
		Height: c.Height(),
	}

	if isNil {
		precommit.ProposedBlockHash = common.Hash{}
	} else {
		hash := c.msgCache.proposal(c.Height().Uint64(), c.Round(), c.committee.GetProposer(c.Round()).Address).ProposedValueHash()
		if hash == (common.Hash{}) {
			c.logger.Error("core.sendPrecommit Proposal is empty! It should not be empty!")
			return
		}
		precommit.ProposedBlockHash = hash
	}

	encodedVote, err := Encode(&precommit)
	if err != nil {
		logger.Error("Failed to encode", "subject", precommit)
		return
	}

	c.logPrecommitMessageEvent("MessageEvent(Precommit): Sent", precommit, c.address.String(), "broadcast")

	msg := &Message{
		Code:          msgPrecommit,
		Msg:           encodedVote,
		Address:       c.address,
		CommittedSeal: []byte{},
	}

	// Create committed seal
	seal := PrepareCommittedSeal(precommit.ProposedBlockHash, c.Round(), c.Height())
	msg.CommittedSeal, err = c.backend.Sign(seal)
	if err != nil {
		c.logger.Error("core.sendPrecommit error while signing committed seal", "err", err)
	}

	c.sentPrecommit = true
	c.broadcast(ctx, msg)
}

func (c *core) handlePrecommit(ctx context.Context, preCommit *Vote, header *types.Header) error {
	if preCommit.Round > c.Round() {
		// If it's a future round precommit leave it.
		return nil
	}

	proposal := c.msgCache.proposal(preCommit.Height.Uint64(), preCommit.Round, c.committee.GetProposer(preCommit.Round).Address)
	proposalHash := proposal.ProposedValueHash()

	if preCommit.Round < c.Round() {
		if proposalHash != (common.Hash{}) && c.msgCache.precommitPower(proposalHash, proposal.Round, header) >= c.committeeSet().Quorum() {
			c.logger.Info("Quorum on a old round proposal", "round",
				preCommit.Round)
			if c.msgCache.proposalVerified(proposalHash) {
				c.commit(proposal)
			}
			return nil
		}
	}

	// At this point we know we have a precommit that has the same height and
	// round as core. We don't know if it is a vote for the proposal we looked
	// up though.

	// Line 49 in Algorithm 1 of The latest gossip on BFT consensus
	// We don't care about which step we are in to accept a preCommit, since it has the highest importance

	if proposalHash != (common.Hash{}) && c.msgCache.precommitPower(proposalHash, proposal.Round, header) >= c.committeeSet().Quorum() {
		if err := c.precommitTimeout.stopTimer(); err != nil {
			return err
		}
		c.logger.Debug("Stopped Scheduled Precommit Timeout")

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			c.commit(proposal)
		}

		// Line 47 in Algorithm 1 of The latest gossip on BFT consensus
	} else if !c.precommitTimeout.timerStarted() && c.msgCache.totalPrecommitPower(proposal.Round, header) >= c.committeeSet().Quorum() {
		timeoutDuration := c.timeoutPrecommit(proposal.Round)
		c.precommitTimeout.scheduleTimeout(timeoutDuration, proposal.Round, proposal.Height, c.onTimeoutPrecommit)
		c.logger.Debug("Scheduled Precommit Timeout", "Timeout Duration", timeoutDuration)
	}

	return nil
}

func (c *core) verifyCommittedSeal(addressMsg common.Address, committedSealMsg []byte, proposedBlockHash common.Hash, round int64, height *big.Int) error {
	committedSeal := PrepareCommittedSeal(proposedBlockHash, round, height)

	sealerAddress, err := types.GetSignatureAddress(committedSeal, committedSealMsg)
	if err != nil {
		c.logger.Error("Failed to get signer address", "err", err)
		return err
	}

	// ensure sender signed the committed seal
	if !bytes.Equal(sealerAddress.Bytes(), addressMsg.Bytes()) {
		c.logger.Error("verify precommit seal error", "got", addressMsg.String(), "expected", sealerAddress.String())

		return errInvalidSenderOfCommittedSeal
	}

	return nil
}

func (c *core) handleCommit(ctx context.Context) {
	c.logger.Debug("Received a final committed proposal", "step", c.step)
	lastBlock, _ := c.backend.LastCommittedProposal()
	height := new(big.Int).Add(lastBlock.Number(), common.Big1)
	if height.Cmp(c.Height()) == 0 {
		c.logger.Debug("Discarding event as core is at the same height", "height", c.Height())
	} else {
		c.logger.Debug("Received proposal is ahead", "height", c.Height(), "block_height", height)
		c.startRound(ctx, 0)
	}
}

func (c *core) logPrecommitMessageEvent(message string, precommit Vote, from, to string) {
	currentProposalHash := c.msgCache.proposal(precommit.Height.Uint64(), precommit.Round, c.committee.GetProposer(precommit.Round).Address).ProposedValueHash()
	c.logger.Debug(message,
		"from", from,
		"to", to,
		"currentHeight", c.Height(),
		"msgHeight", precommit.Height,
		"currentRound", c.Round(),
		"msgRound", precommit.Round,
		"currentStep", c.step,
		"isProposer", c.isProposer(),
		"currentProposer", c.committeeSet().GetProposer(c.Round()),
		"isNilMsg", precommit.ProposedBlockHash == common.Hash{},
		"hash", precommit.ProposedBlockHash,
		"type", "Precommit",
		"totalVotes", c.msgCache.totalPrecommitPower(precommit.Round, c.lastHeader),
		"totalNilVotes", c.msgCache.precommitPower(common.Hash{}, precommit.Round, c.lastHeader),
		"VoteProposedBlock", c.msgCache.precommitPower(currentProposalHash, precommit.Round, c.lastHeader),
	)
}
