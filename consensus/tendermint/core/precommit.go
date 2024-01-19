package core

import (
	"context"
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
)

type Precommiter struct {
	*Core
}

func (c *Precommiter) SendPrecommit(ctx context.Context, isNil bool) {
	value := common.Hash{}
	if !isNil {
		proposal := c.curRoundMessages.Proposal()
		if proposal == nil {
			c.logger.Error("sendPrevote Proposal is empty! It should not be empty!")
			return
		}
		value = proposal.Block().Hash()
		c.logger.Info("Precommiting on proposal", "proposal", proposal.Block().Hash(), "round", c.Round(), "height", c.Height().Uint64())
	} else {
		c.logger.Info("Precommiting on nil", "round", c.Round(), "height", c.Height().Uint64())
	}

	precommit := message.NewPrecommit(c.Round(), c.Height().Uint64(), value, c.backend.Sign)
	c.LogPrecommitMessageEvent("Precommit sent", precommit, c.address.String(), "broadcast")
	c.sentPrecommit = true
	c.Broadcaster().Broadcast(precommit)
}

// HandlePrecommit process the incoming precommit message.
func (c *Precommiter) HandlePrecommit(ctx context.Context, precommit *message.Precommit) error {
	if precommit.R() > c.Round() {
		return constants.ErrFutureRoundMessage
	}
	if precommit.R() < c.Round() {
		// We are receiving a precommit for an old round. We need to check if we have now a quorum
		// in this old round.
		roundMessages := c.messages.GetOrCreate(precommit.R())
		roundMessages.AddPrecommit(precommit)

		oldRoundProposal := roundMessages.Proposal()
		if oldRoundProposal == nil {
			return constants.ErrOldRoundMessage
		}

		// Line 49 in Algorithm 1 of The latest gossip on BFT consensus
		_ = c.quorumPrecommitsCheck(ctx, oldRoundProposal, roundMessages.IsProposalVerified())
		return constants.ErrOldRoundMessage
	}

	// Precommit if for current round from here
	// We don't care about which step we are in to accept a precommit, since it has the highest importance
	c.curRoundMessages.AddPrecommit(precommit)
	c.LogPrecommitMessageEvent("MessageEvent(Precommit): Received", precommit, precommit.Sender().String(), c.address.String())

	c.currentPrecommitChecks(ctx)
	return nil
}

func (c *Precommiter) HandleCommit(ctx context.Context) {
	c.logger.Debug("Received a final committed proposal", "step", c.step)
	lastBlock := c.backend.HeadBlock()
	height := new(big.Int).Add(lastBlock.Number(), common.Big1)
	if height.Cmp(c.Height()) == 0 {
		c.logger.Debug("Discarding event as Core is at the same height", "height", c.Height())
	} else {
		c.logger.Debug("New chain head ahead of consensus Core height", "height", c.Height(), "block_height", height)
		c.StartRound(ctx, 0)
	}
}

func (c *Precommiter) LogPrecommitMessageEvent(message string, precommit *message.Precommit, from, to string) {
	currentProposalHash := c.curRoundMessages.ProposalHash()
	c.logger.Debug(message,
		"from", from,
		"to", to,
		"currentHeight", c.Height(),
		"msgHeight", precommit.H(),
		"currentRound", c.Round(),
		"msgRound", precommit.R(),
		"currentStep", c.step,
		"isProposer", c.IsProposer(),
		"currentProposer", c.CommitteeSet().GetProposer(c.Round()),
		"isNilMsg", precommit.Value() == common.Hash{},
		"hash", precommit.Value(),
		"type", "Precommit",
		"totalVotes", c.curRoundMessages.PrecommitsTotalPower(),
		"totalNilVotes", c.curRoundMessages.PrecommitsPower(common.Hash{}),
		"proposedBlockVote", c.curRoundMessages.PrecommitsPower(currentProposalHash),
	)
}
