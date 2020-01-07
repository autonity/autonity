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
	"math/big"

	"github.com/clearmatics/autonity/common"
)

func (c *core) sendPrevote(ctx context.Context, isNil bool) {
	logger := c.logger.New("step", c.roundState.Step())
	currentRound := c.roundState.Round().Int64()

	var prevote = Vote{
		Round:  big.NewInt(c.roundState.Round().Int64()),
		Height: big.NewInt(c.roundState.Height().Int64()),
	}

	proposalBlockHash := c.roundState.Proposal(currentRound).ProposalBlock.Hash()
	if isNil {
		prevote.ProposedBlockHash = common.Hash{}
	} else {
		if h := proposalBlockHash; h == (common.Hash{}) {
			c.logger.Error("sendPrecommit Proposal is empty! It should not be empty!")
			return
		}
		prevote.ProposedBlockHash = proposalBlockHash
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

func (c *core) handlePrevote(ctx context.Context, msg *Message) error {
	var preVote Vote
	err := msg.Decode(&preVote)
	if err != nil {
		return errFailedDecodePrevote
	}

	if err = c.checkMessage(preVote.Round, preVote.Height, prevote); err != nil {
		// Store old round prevote messages for future rounds since it is required for validRound
		if err == errOldRoundMessage {
			// We only process old rounds while future rounds messages are pushed on to the backlog
			oldRoundMessageSet, ok := c.roundState.allRoundMessages[preVote.Round.Int64()]
			if !ok {
				oldRoundMessageSet = newRoundMessageSet()
				c.roundState.allRoundMessages[preVote.Round.Int64()] = oldRoundMessageSet
			}

			oldRoundPrevotes := oldRoundMessageSet.prevotes
			oldRoundPrevotes.Add(preVote.ProposedBlockHash, *msg)
		}
		return err
	}

	// After checking the message we know it is from the same height and round, so we should store it even if
	// c.roundState.Step() < prevote. The propose timeout which is started at the beginning of the round
	// will update the step to at least prevote and when it handle its on preVote(nil), then it will also have
	// votes from other nodes.
	curR := c.roundState.Round().Int64()
	curH := c.roundState.Height().Int64()
	curProposalHash := c.roundState.Proposal(curR).ProposalBlock.Hash()
	prevoteHash := preVote.ProposedBlockHash
	prevotes := c.roundState.allRoundMessages[curR].prevotes
	prevotes.Add(prevoteHash, *msg)

	c.logPrevoteMessageEvent("MessageEvent(Prevote): Received", preVote, msg.Address.String(), c.address.String())

	// Now we can add the preVote to our current round state
	if c.roundState.Step() >= prevote {
		// Line 36 in Algorithm 1 of The latest gossip on BFT consensus
		if curProposalHash != (common.Hash{}) && c.Quorum(prevotes.VotesSize(curProposalHash)) && !c.setValidRoundAndValue {
			// this piece of code should only run once
			if err := c.prevoteTimeout.stopTimer(); err != nil {
				return err
			}
			c.logger.Debug("Stopped Scheduled Prevote Timeout")

			if c.roundState.Step() == prevote {
				c.lockedValue = c.roundState.Proposal(curR).ProposalBlock
				c.lockedRound = big.NewInt(curR)
				c.sendPrecommit(ctx, false)
				c.setStep(precommit)
			}
			c.validValue = c.roundState.Proposal(curR).ProposalBlock
			c.validRound = big.NewInt(curR)
			c.setValidRoundAndValue = true
			// Line 44 in Algorithm 1 of The latest gossip on BFT consensus
		} else if c.roundState.Step() == prevote && c.Quorum(prevotes.NilVotesSize()) {
			if err := c.prevoteTimeout.stopTimer(); err != nil {
				return err
			}
			c.logger.Debug("Stopped Scheduled Prevote Timeout")

			c.sendPrecommit(ctx, true)
			c.setStep(precommit)

			// Line 34 in Algorithm 1 of The latest gossip on BFT consensus
		} else if c.roundState.Step() == prevote && !c.prevoteTimeout.timerStarted() && !c.sentPrecommit && c.Quorum(prevotes.TotalSize()) {
			timeoutDuration := timeoutPrevote(curR)
			c.prevoteTimeout.scheduleTimeout(timeoutDuration, curR, curH, c.onTimeoutPrevote)
			c.logger.Debug("Scheduled Prevote Timeout", "Timeout Duration", timeoutDuration)
		}
	}

	return nil
}

func (c *core) logPrevoteMessageEvent(message string, prevote Vote, from, to string) {
	currentRound := c.roundState.Round().Int64()
	currentProposalHash := c.roundState.Proposal(currentRound).ProposalBlock.Hash()
	prevotes := c.roundState.allRoundMessages[currentRound].prevotes
	c.logger.Debug(message,
		"from", from,
		"to", to,
		"currentHeight", c.roundState.Height(),
		"msgHeight", prevote.Height,
		"currentRound", c.roundState.Round(),
		"msgRound", prevote.Round,
		"currentStep", c.roundState.Step(),
		"isProposer", c.isProposer(),
		"currentProposer", c.valSet.GetProposer(),
		"isNilMsg", prevote.ProposedBlockHash == common.Hash{},
		"hash", prevote.ProposedBlockHash,
		"type", "Prevote",
		"totalVotes", prevotes.TotalSize(),
		"totalNilVotes", prevotes.NilVotesSize(),
		"quorumReject", c.Quorum(prevotes.NilVotesSize()),
		"totalNonNilVotes", prevotes.VotesSize(currentProposalHash),
		"quorumAccept", c.Quorum(prevotes.VotesSize(currentProposalHash)),
	)
}
