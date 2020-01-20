package core

import (
	"context"
	"math/big"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
)

func (c *core) commit(ctx context.Context, round int64) {
	_ = c.setStep(ctx, precommitDone)

	proposal := c.getProposal(round)
	if proposal == nil {
		// Should never happen really.
		c.logger.Error("core commit called with empty proposal ")
		return
	}

	if proposal.ProposalBlock == nil {
		// Again should never happen.
		c.logger.Error("commit a NIL block", "block", proposal.ProposalBlock, "height", c.getHeight().String(), "round", c.getRound().String())
		return
	}

	c.logger.Info("commit a block", "hash", proposal.ProposalBlock.Header().Hash())

	precommits := c.allRoundMessages[round].precommits
	committedSeals := make([][]byte, precommits.VotesSize(proposal.ProposalBlock.Hash()))
	for i, v := range precommits.Values(proposal.ProposalBlock.Hash()) {
		committedSeals[i] = make([]byte, types.BFTExtraSeal)
		copy(committedSeals[i][:], v.CommittedSeal[:])
	}

	if err := c.backend.Commit(proposal.ProposalBlock, c.getRound(), committedSeals); err != nil {
		c.logger.Error("failed to commit a block", "err", err)
		return
	}
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
