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
	logger := c.logger.New("step", c.roundState.Step())
	currentRound := c.roundState.Round().Int64()

	var precommit = Vote{
		Round:  big.NewInt(c.roundState.Round().Int64()),
		Height: big.NewInt(c.roundState.Height().Int64()),
	}

	proposalBlockHash := c.roundState.Proposal(currentRound).ProposalBlock.Hash()
	if isNil {
		precommit.ProposedBlockHash = common.Hash{}
	} else {
		if h := proposalBlockHash; h == (common.Hash{}) {
			c.logger.Error("core.sendPrecommit Proposal is empty! It should not be empty!")
			return
		}
		precommit.ProposedBlockHash = proposalBlockHash
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
	seal := PrepareCommittedSeal(precommit.ProposedBlockHash, c.roundState.Round(), c.roundState.Height())
	msg.CommittedSeal, err = c.backend.Sign(seal)
	if err != nil {
		c.logger.Error("core.sendPrecommit error while signing committed seal", "err", err)
	}

	c.sentPrecommit = true
	c.broadcast(ctx, msg)
}

func (c *core) handlePrecommit(ctx context.Context, msg *Message) error {
	var preCommit Vote
	err := msg.Decode(&preCommit)
	if err != nil {
		return errFailedDecodePrecommit
	}

	if err := c.checkMessage(preCommit.Round, preCommit.Height, precommit); err != nil {
		return err
	}

	// Don't want to decode twice, hence sending preCommit with message
	if err := c.verifyPrecommitCommittedSeal(msg.Address, append([]byte(nil), msg.CommittedSeal...), preCommit.ProposedBlockHash, preCommit.Round, preCommit.Height); err != nil {
		return err
	}

	// We don't care about which step we are in to accept a preCommit, since it has the highest importance
	// TODO: use uints instead of ints
	precommitHash := preCommit.ProposedBlockHash
	curR := c.roundState.Round().Int64()
	curH := c.roundState.Height().Int64()
	precommits := c.roundState.allRoundMessages[curR].precommits

	precommits.Add(precommitHash, *msg)

	c.logPrecommitMessageEvent("MessageEvent(Precommit): Received", preCommit, msg.Address.String(), c.address.String())

	// Line 49 in Algorithm 1 of The latest gossip on BFT consensus
	curProposalHash := c.roundState.Proposal(curR).ProposalBlock.Hash()
	if curProposalHash != (common.Hash{}) && c.Quorum(precommits.VotesSize(curProposalHash)) {
		if err := c.precommitTimeout.stopTimer(); err != nil {
			return err
		}
		c.logger.Debug("Stopped Scheduled Precommit Timeout")

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			c.commit(curR)
		}

		// Line 47 in Algorithm 1 of The latest gossip on BFT consensus
	} else if !c.precommitTimeout.timerStarted() && c.Quorum(precommits.TotalSize()) {
		timeoutDuration := timeoutPrecommit(curR)
		c.precommitTimeout.scheduleTimeout(timeoutDuration, curR, curH, c.onTimeoutPrecommit)
		c.logger.Debug("Scheduled Precommit Timeout", "Timeout Duration", timeoutDuration)
	}

	return nil
}

func (c *core) verifyPrecommitCommittedSeal(addressMsg common.Address, committedSealMsg []byte, proposedBlockHash common.Hash, round *big.Int, height *big.Int) error {
	committedSeal := PrepareCommittedSeal(proposedBlockHash, round, height)

	addressOfSignerOfCommittedSeal, err := types.GetSignatureAddress(committedSeal, committedSealMsg)
	if err != nil {
		c.logger.Error("Failed to get signer address", "err", err)
		return err
	}

	// ensure sender signed the committed seal
	if !bytes.Equal(addressOfSignerOfCommittedSeal.Bytes(), addressMsg.Bytes()) {
		c.logger.Error("verify precommit seal error", "got", addressMsg.String(), "expected", addressOfSignerOfCommittedSeal.String())

		return errInvalidSenderOfCommittedSeal
	}

	return nil
}

func (c *core) handleCommit(ctx context.Context) {
	c.logger.Debug("Received a final committed proposal", "step", c.roundState.Step())
	lastBlock, _ := c.backend.LastCommittedProposal()
	height := new(big.Int).Add(lastBlock.Number(), common.Big1).Uint64()
	if height == c.roundState.Height().Uint64() {
		c.logger.Debug("Discarding event as core is at the same height", "state_height", c.roundState.Height().Uint64())
	} else {
		c.logger.Debug("Received proposal is ahead", "state_height", c.roundState.Height().Uint64(), "block_height", height)
		c.startRound(ctx, common.Big0)
	}
}

func (c *core) logPrecommitMessageEvent(message string, precommit Vote, from, to string) {
	currentRound := c.roundState.Round().Int64()
	currentProposalHash := c.roundState.Proposal(currentRound).ProposalBlock.Hash()
	precommits := c.roundState.allRoundMessages[currentRound].precommits
	c.logger.Debug(message,
		"from", from,
		"to", to,
		"currentHeight", c.roundState.Height(),
		"msgHeight", precommit.Height,
		"currentRound", c.roundState.Round(),
		"msgRound", precommit.Round,
		"currentStep", c.roundState.Step(),
		"isProposer", c.isProposer(),
		"currentProposer", c.valSet.GetProposer(),
		"isNilMsg", precommit.ProposedBlockHash == common.Hash{},
		"hash", precommit.ProposedBlockHash,
		"type", "Precommit",
		"totalVotes", precommits.TotalSize(),
		"totalNilVotes", precommits.NilVotesSize(),
		"quorumReject", c.Quorum(precommits.NilVotesSize()),
		"totalNonNilVotes", precommits.VotesSize(currentProposalHash),
		"quorumAccept", c.Quorum(precommits.VotesSize(currentProposalHash)),
	)
}
