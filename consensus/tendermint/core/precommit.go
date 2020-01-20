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
	currentRound := c.getRound().Int64()

	var precommit = Vote{
		Round:  big.NewInt(c.getRound().Int64()),
		Height: big.NewInt(c.getHeight().Int64()),
	}

	proposalBlockHash := c.getProposal(currentRound).ProposalBlock.Hash()
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
	seal := PrepareCommittedSeal(precommit.ProposedBlockHash, c.getRound(), c.getHeight())
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

	// Check for nil values
	if preCommit.Height == nil || preCommit.Round == nil {
		return errInvalidMessage
	}

	// Ensure prevote is for current height
	if err = c.checkMessage(preCommit.Round, preCommit.Height, prevote); err != nil {
		return err
	}

	// If we already have the prevote do nothing
	if c.hasVote(preCommit, msg) {
		return nil
	}

	// Don't want to decode twice, hence sending preCommit with message
	if err := c.verifyPrecommitCommittedSeal(msg.Address, append([]byte(nil), msg.CommittedSeal...), preCommit.ProposedBlockHash, preCommit.Round, preCommit.Height); err != nil {
		return err
	}
	c.logPrecommitMessageEvent("MessageEvent(Precommit): Received", preCommit, msg.Address.String(), c.address.String())

	// We don't care about which step we are in to accept a preCommit, since it has the highest importance
	curR := c.getRound().Int64()
	curH := c.getHeight().Int64()
	precommitHash := preCommit.ProposedBlockHash

	// The precommit doesn't exists in our current round state, so add it, thus it will add the precommit to the round
	// of the precommit
	precommits := c.allRoundMessages[preCommit.Round.Int64()].precommits
	precommits.Add(precommitHash, *msg)

	roundCmp := preCommit.Round.Cmp(c.getRound())
	if roundCmp == 0 {
		//Check for timeout only if preCommit.Round == curR
		c.checkForPrecommitTimeout(curR, curH)
	}

	// Check for consensus regardless of the precommit round
	if err := c.checkForConsensus(ctx, preCommit.Round.Int64()); err != nil {
		return err
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
	c.logger.Debug("Received a final committed proposal", "step", c.step)
	lastBlock, _ := c.backend.LastCommittedProposal()
	height := new(big.Int).Add(lastBlock.Number(), common.Big1).Uint64()
	if height == c.getHeight().Uint64() {
		c.logger.Debug("Discarding event as core is at the same height", "state_height", c.getHeight().Uint64())
	} else {
		c.logger.Debug("Received proposal is ahead", "state_height", c.getHeight().Uint64(), "block_height", height)
		c.startRound(ctx, common.Big0)
	}
}

func (c *core) logPrecommitMessageEvent(message string, precommit Vote, from, to string) {
	currentRound := c.getRound().Int64()
	currentProposalHash := c.getProposal(currentRound).ProposalBlock.Hash()
	precommits := c.allRoundMessages[currentRound].precommits
	c.logger.Debug(message,
		"from", from,
		"to", to,
		"currentHeight", c.getHeight(),
		"msgHeight", precommit.Height,
		"currentRound", c.getRound(),
		"msgRound", precommit.Round,
		"currentStep", c.step,
		"isProposer", c.isProposerForR(c.getRound().Int64(), c.address),
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
