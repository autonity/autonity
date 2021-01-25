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

	var precommit = types.Vote{
		Round:  c.Round(),
		Height: c.Height(),
	}

	if isNil {
		precommit.ProposedBlockHash = common.Hash{}
	} else {
		if h := c.curRoundMessages.GetProposalHash(); h == (common.Hash{}) {
			c.logger.Error("core.sendPrecommit Proposal is empty! It should not be empty!")
			return
		}
		precommit.ProposedBlockHash = c.curRoundMessages.GetProposalHash()
	}

	encodedVote, err := types.Encode(&precommit)
	if err != nil {
		logger.Error("Failed to encode", "subject", precommit)
		return
	}

	c.logPrecommitMessageEvent("MessageEvent(Precommit): Sent", precommit, c.address.String(), "broadcast")

	msg := &types.ConsensusMessage{
		Code:          types.MsgPrecommit,
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

func (c *core) handlePrecommit(ctx context.Context, msg *types.ConsensusMessage) error {
	var preCommit types.Vote
	err := msg.Decode(&preCommit)
	if err != nil {
		return errFailedDecodePrecommit
	}
	precommitHash := preCommit.ProposedBlockHash

	if err := c.checkMessage(preCommit.Round, preCommit.Height, precommit); err != nil {

		if err == errOldRoundMessage {
			roundMsgs := c.messages.getOrCreate(preCommit.Round)
			if error := c.verifyCommittedSeal(msg.Address, append([]byte(nil), msg.CommittedSeal...), preCommit.ProposedBlockHash, preCommit.Round, preCommit.Height); error != nil {
				return error
			}
			c.acceptVote(roundMsgs, precommit, precommitHash, *msg)
			oldRoundProposalHash := roundMsgs.GetProposalHash()
			if oldRoundProposalHash != (common.Hash{}) && roundMsgs.PrecommitsPower(oldRoundProposalHash) >= c.committeeSet().Quorum() {
				c.logger.Info("Quorum on a old round proposal", "round", preCommit.Round)
				if !roundMsgs.isProposalVerified() {
					if _, error := c.backend.VerifyProposal(*roundMsgs.Proposal().ProposalBlock); error != nil {
						return error
					}
				}
				c.commit(preCommit.Round, c.curRoundMessages)
				return nil
			}
		}

		return err
	}

	// Don't want to decode twice, hence sending preCommit with message
	if err := c.verifyCommittedSeal(msg.Address, append([]byte(nil), msg.CommittedSeal...), preCommit.ProposedBlockHash, preCommit.Round, preCommit.Height); err != nil {
		return err
	}
	// Line 49 in Algorithm 1 of The latest gossip on BFT consensus
	curProposalHash := c.curRoundMessages.GetProposalHash()
	// We don't care about which step we are in to accept a preCommit, since it has the highest importance

	c.acceptVote(c.curRoundMessages, precommit, precommitHash, *msg)
	c.logPrecommitMessageEvent("MessageEvent(Precommit): Received", preCommit, msg.Address.String(), c.address.String())
	if curProposalHash != (common.Hash{}) && c.curRoundMessages.PrecommitsPower(curProposalHash) >= c.committeeSet().Quorum() {
		if err := c.precommitTimeout.stopTimer(); err != nil {
			return err
		}
		c.logger.Debug("Stopped Scheduled Precommit Timeout")

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			c.commit(c.Round(), c.curRoundMessages)
		}

		// Line 47 in Algorithm 1 of The latest gossip on BFT consensus
	} else if !c.precommitTimeout.timerStarted() && c.curRoundMessages.PrecommitsTotalPower() >= c.committeeSet().Quorum() {
		timeoutDuration := c.timeoutPrecommit(c.Round())
		c.precommitTimeout.scheduleTimeout(timeoutDuration, c.Round(), c.Height(), c.onTimeoutPrecommit)
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

func (c *core) logPrecommitMessageEvent(message string, precommit types.Vote, from, to string) {
	currentProposalHash := c.curRoundMessages.GetProposalHash()
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
		"totalVotes", c.curRoundMessages.PrecommitsTotalPower(),
		"totalNilVotes", c.curRoundMessages.PrecommitsPower(common.Hash{}),
		"proposedBlockVote", c.curRoundMessages.PrecommitsPower(currentProposalHash),
	)
}
