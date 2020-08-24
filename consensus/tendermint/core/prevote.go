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
	types "github.com/clearmatics/autonity/core/types"
)

func (c *core) sendPrevote(ctx context.Context, isNil bool) {
	logger := c.logger.New("step", c.step)

	var prevote = Vote{
		Round:  c.Round(),
		Height: c.Height(),
	}

	if isNil {
		prevote.ProposedBlockHash = common.Hash{}
	} else {
		if h := c.curRoundMessages.GetProposalHash(); h == (common.Hash{}) {
			c.logger.Error("sendPrecommit Proposal is empty! It should not be empty!")
			return
		}
		prevote.ProposedBlockHash = c.curRoundMessages.GetProposalHash()
	}

	encodedVote, err := Encode(&prevote)
	if err != nil {
		logger.Error("Failed to encode", "subject", prevote)
		return
	}

	c.logPrevoteMessageEvent("MessageEvent(Prevote): Sent", prevote, c.address.String(), "broadcast")

	c.sentPrevote = true
	c.broadcast(ctx, &Message{
		Code:          msgPrevote,
		Msg:           encodedVote,
		Address:       c.address,
		CommittedSeal: []byte{},
	})
}

func (c *core) handlePrevote(ctx context.Context, preVote *Vote, header *types.Header) error {
	if preVote.Round > c.Round() || preVote.Round < c.Round() {
		// If it's a future or past round prevote leave it.
		return nil
	}

	proposal := c.msgCache.proposal(preVote.Height.Uint64(), preVote.Round, c.committee.GetProposer(preVote.Round).Address)

	// Line 36 in Algorithm 1 of The latest gossip on BFT consensus
	if proposal != nil && c.step >= prevote && c.msgCache.prevotePower(proposal.ProposedValueHash(), header) >= c.committeeSet().Quorum() && !c.setValidRoundAndValue {
		// this piece of code should only run once
		if err := c.prevoteTimeout.stopTimer(); err != nil {
			return err
		}
		c.logger.Debug("Stopped Scheduled Prevote Timeout")

		if c.step == prevote {
			c.lockedValue = proposal.ProposalBlock
			c.lockedRound = proposal.Round
			c.sendPrecommit(ctx, false)
			c.setStep(precommit)
		}
		c.validValue = proposal.ProposalBlock
		c.validRound = proposal.Round
		c.setValidRoundAndValue = true
		// Line 44 in Algorithm 1 of The latest gossip on BFT consensus
	} else if c.step == prevote && c.msgCache.prevotePower(common.Hash{}, header) >= c.committeeSet().Quorum() {
		if err := c.prevoteTimeout.stopTimer(); err != nil {
			return err
		}
		c.logger.Debug("Stopped Scheduled Prevote Timeout")

		c.sendPrecommit(ctx, true)
		c.setStep(precommit)

		// Line 34 in Algorithm 1 of The latest gossip on BFT consensus
	} else if c.step == prevote && !c.prevoteTimeout.timerStarted() && !c.sentPrecommit && c.msgCache.totalPrevotePower(preVote.Round, header) >= c.committeeSet().Quorum() {
		timeoutDuration := c.timeoutPrevote(c.Round())
		c.prevoteTimeout.scheduleTimeout(timeoutDuration, c.Round(), c.Height(), c.onTimeoutPrevote)
		c.logger.Debug("Scheduled Prevote Timeout", "Timeout Duration", timeoutDuration)
	}

	return nil
}

func (c *core) logPrevoteMessageEvent(message string, prevote Vote, from, to string) {
	currentProposalHash := c.curRoundMessages.GetProposalHash()
	c.logger.Debug(message,
		"from", from,
		"to", to,
		"currentHeight", c.Height(),
		"msgHeight", prevote.Height,
		"currentRound", c.Round(),
		"msgRound", prevote.Round,
		"currentStep", c.step,
		"isProposer", c.isProposer(),
		"currentProposer", c.committeeSet().GetProposer(c.Round()),
		"isNilMsg", prevote.ProposedBlockHash == common.Hash{},
		"hash", prevote.ProposedBlockHash,
		"type", "Prevote",
		"totalVotes", c.curRoundMessages.PrevotesTotalPower(),
		"totalNilVotes", c.curRoundMessages.PrevotesPower(common.Hash{}),
		"VoteProposedBlock", c.curRoundMessages.PrevotesPower(currentProposalHash),
	)
}
