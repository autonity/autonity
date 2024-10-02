package core

import (
	"context"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/metrics"
	"time"
)

type Committer struct {
	*Core
}

func (c *Committer) Commit(ctx context.Context, round int64, messages *message.RoundMessages) {
	c.SetStep(ctx, PrecommitDone)
	// for metrics
	start := time.Now()
	proposal := messages.Proposal()
	if proposal == nil {
		// Should never happen really. Let's panic to catch bugs.
		panic("Core commit called with empty proposal")
	}
	proposalHash := proposal.Block().Header().Hash()
	c.logger.Debug("Committing a block", "hash", proposalHash)

	precommitWithQuorum := messages.PrecommitFor(proposalHash)
	quorumCertificate := types.NewAggregateSignature(precommitWithQuorum.Signature().(*blst.BlsSignature), precommitWithQuorum.Signers())

	// todo: Jason, since commit() checks extra conditions, shall we need to move those condition here before we flush
	//  the decision? However this movement is quit heavy, since we have consensus conditions and chain context conditions
	//  to be checked.
	// record decision in WAL before the submission.
	c.SetDecision(proposal.Block(), proposal.R())

	if err := c.backend.Commit(proposal.Block(), round, quorumCertificate); err != nil {
		c.logger.Error("failed to commit a block", "err", err)
		return
	}

	if metrics.Enabled {
		now := time.Now()
		CommitTimer.Update(now.Sub(start))
		CommitBg.Add(now.Sub(start).Nanoseconds())
	}
}
