package core

import (
	"context"
	"math/big"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
)

func (c *core) commit(ctx context.Context, round int64) {
	_ = c.setStep(ctx, precommitDone)

	proposalMS := c.getProposalSet(round)
	if proposalMS == nil {
		// Should never happen really.
		c.logger.Error("core commit called with empty proposal ")
		return
	}

	proposal := proposalMS.proposal()
	if proposal.ProposalBlock == nil {
		// Again should never happen.
		c.logger.Error("commit a NIL block", "block", proposal.ProposalBlock, "height", c.getHeight().String(), "round", c.getRound().String())
		return
	}

	c.logger.Info("commit a block", "hash", proposal.ProposalBlock.Header().Hash())

	precommits := c.getPrecommitsSet(round)
	if precommits == nil {
		// This should never be the case
		c.logger.Error("Precommits empty in commit()")
		return
	}
	committedSeals := make([][]byte, precommits.VotesSize(proposal.ProposalBlock.Hash()))
	for i, v := range precommits.AllBlockHashMessages(proposal.ProposalBlock.Hash()) {
		committedSeals[i] = make([]byte, types.BFTExtraSeal)
		copy(committedSeals[i][:], v.CommittedSeal[:])
	}

	if err := c.backend.Commit(proposal.ProposalBlock, c.getRound(), committedSeals); err != nil {
		c.logger.Error("failed to commit a block", "err", err)
		return
	}
}

func (c *core) handleCommit(ctx context.Context) {
	c.logger.Debug("Received a final committed proposal", "step", c.getStep())
	lastBlock, _ := c.backend.LastCommittedProposal()
	height := new(big.Int).Add(lastBlock.Number(), common.Big1).Uint64()
	if height == c.getHeight().Uint64() {
		c.logger.Debug("Discarding event as core is at the same height", "state_height", c.getHeight().Uint64())
	} else {
		c.logger.Debug("Received proposal is ahead", "state_height", c.getHeight().Uint64(), "block_height", height)
		c.startRound(ctx, common.Big0)
	}
}
