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
	"github.com/clearmatics/autonity/core/types"
	"math/big"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
)

func (c *core) sendPrecommit(ctx context.Context, isNil bool) {
	logger := c.logger.New("step", c.currentRoundState.Step())

	var precommit = tendermint.Vote{
		Round:  big.NewInt(c.currentRoundState.Round().Int64()),
		Height: big.NewInt(c.currentRoundState.Height().Int64()),
	}

	if isNil {
		precommit.ProposedBlockHash = common.Hash{}
	} else {
		if h := c.currentRoundState.GetCurrentProposalHash(); h == (common.Hash{}) {
			c.logger.Error("core.sendPrecommit Proposal is empty! It should not be empty!")
			return
		}
		precommit.ProposedBlockHash = c.currentRoundState.GetCurrentProposalHash()
	}

	encodedVote, err := Encode(&precommit)
	if err != nil {
		logger.Error("Failed to encode", "subject", precommit)
		return
	}

	c.logPrecommitMessageEvent("MessageEvent(Precommit): Sent", precommit, c.address.String(), "broadcast")

	c.sentPrecommit = true
	c.broadcast(ctx, &message{
		Code: msgPrecommit,
		Msg:  encodedVote,
	})
}

func (c *core) handlePrecommit(ctx context.Context, msg *message) error {
	var preCommit tendermint.Vote
	err := msg.Decode(&preCommit)
	if err != nil {
		return errFailedDecodePrecommit
	}

	if err := c.checkMessage(preCommit.Round, preCommit.Height); err != nil {
		// We don't care about old round preCommit messages, otherwise we would not be in a new round rather a new height
		return err
	}

	curProposalHash := c.currentRoundState.GetCurrentProposalHash()

	// Don't want to decode twice, hence sending preCommit with message
	if err := c.verifyPrecommitCommittedSeal(msg, preCommit); err != nil {
		return err
	}

	// We don't care about which step we are in to accept a preCommit, since it has the highest importance
	precommitHash := preCommit.ProposedBlockHash
	curR := c.currentRoundState.Round().Int64()
	curH := c.currentRoundState.Height().Int64()

	c.acceptPrecommit(precommitHash, *msg)

	c.logPrecommitMessageEvent("MessageEvent(Precommit): Received", preCommit, msg.Address.String(), c.address.String())

	// Line 49 in Algorithm 1 of The latest gossip on BFT consensus
	if curProposalHash != (common.Hash{}) && c.quorum(c.currentRoundState.Precommits.VotesSize(curProposalHash)) {
		if err := c.precommitTimeout.stopTimer(); err != nil {
			return err
		}
		c.logger.Debug("Stopped Scheduled Precommit Timeout")

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			c.commit()
		}

		// Line 47 in Algorithm 1 of The latest gossip on BFT consensus
	} else if !c.precommitTimeout.timerStarted() && c.quorum(c.currentRoundState.Precommits.TotalSize()) {
		timeoutDuration := timeoutPrecommit(curR)
		c.precommitTimeout.scheduleTimeout(timeoutDuration, curR, curH, c.onTimeoutPrecommit)
		c.logger.Debug("Scheduled Precommit Timeout", "Timeout Duration", timeoutDuration)
	}

	return nil
}

func (c *core) verifyPrecommitCommittedSeal(m *message, precommit tendermint.Vote) error {
	addressOfSignerOfCommittedSeal, err := types.GetSignatureAddress(PrepareCommittedSeal(precommit.ProposedBlockHash), m.CommittedSeal)

	if err != nil {
		c.logger.Error("Failed to get signer address", "err", err)
		return err
	}

	// ensure sender signed the committed seal
	if !bytes.Equal(addressOfSignerOfCommittedSeal.Bytes(), m.Address.Bytes()) {
		return errInvalidSenderOfCommittedSeal
	}

	return nil
}

func (c *core) acceptPrecommit(precommitHash common.Hash, msg message) {
	if precommitHash == (common.Hash{}) {
		c.currentRoundState.Precommits.AddNilVote(msg)
	} else {
		c.currentRoundState.Precommits.AddVote(precommitHash, msg)
	}
}

func (c *core) handleCommit(ctx context.Context) {
	c.logger.Debug("Received a final committed proposal", "step", c.currentRoundState.Step())
	c.startRound(ctx, common.Big0)
}

func (c *core) logPrecommitMessageEvent(message string, precommit tendermint.Vote, from, to string) {
	currentProposalHash := c.currentRoundState.GetCurrentProposalHash()
	c.logger.Debug(message,
		"from", from,
		"to", to,
		"currentHeight", c.currentRoundState.Height(),
		"msgHeight", precommit.Height,
		"currentRound", c.currentRoundState.Round(),
		"msgRound", precommit.Round,
		"currentStep", c.currentRoundState.Step(),
		"isProposer", c.isProposer(),
		"currentProposer", c.valSet.GetProposer(),
		"isNilMsg", precommit.ProposedBlockHash == common.Hash{},
		"hash", precommit.ProposedBlockHash,
		"type", "Precommit",
		"totalVotes", c.currentRoundState.Precommits.TotalSize(),
		"totalNilVotes", c.currentRoundState.Precommits.NilVotesSize(),
		"quorumReject", c.quorum(c.currentRoundState.Precommits.NilVotesSize()),
		"totalNonNilVotes", c.currentRoundState.Precommits.VotesSize(currentProposalHash),
		"quorumAccept", c.quorum(c.currentRoundState.Precommits.VotesSize(currentProposalHash)),
	)
}
