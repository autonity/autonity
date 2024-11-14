package core

import (
	"context"
	"math/big"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"

	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
)

type Precommiter struct {
	*Core
}

func (c *Precommiter) SendPrecommit(_ context.Context, isNil bool) {
	value := common.Hash{}
	if !isNil {
		proposal := c.curRoundMessages.Proposal()
		if proposal == nil {
			c.logger.Error("sendPrevote Proposal is empty! It should not be empty!")
			return
		}
		value = proposal.Block().Hash()
		c.logger.Info("Precommiting on proposal", "proposal", value, "round", c.Round(), "height", c.Height().Uint64())
	} else {
		c.logger.Info("Precommiting on nil", "round", c.Round(), "height", c.Height().Uint64())
	}
	self, err := c.CommitteeSet().MemberByAddress(c.address)
	if err != nil {
		// it can happen in edge case addressed by docker e2e test, that is a validator resets at the epoch boundary,
		// after which it leaves the committee, we cannot panic it in that case.
		c.logger.Error("Validator is no longer in current committee", "err", err, "validator", c.address.String())
		return
	}
	precommit := message.NewPrecommit(c.Round(), c.Height().Uint64(), value, c.backend.Sign, self, c.CommitteeSet().Committee().Len())
	c.LogPrecommitMessageEvent("Precommit sent", precommit)
	c.sentPrecommit = true
	c.Broadcaster().Broadcast(precommit)
	if metrics.Enabled {
		PrecommitSentBlockTSDeltaBg.Add(time.Since(c.currBlockTimeStamp).Nanoseconds())
	}
}

func (c *Precommiter) HandlePrecommit(ctx context.Context, precommit *message.Precommit) error {
	if !precommit.PreVerified() || !precommit.Verified() {
		panic("Handling NON cryptographically verified precommit")
	}

	if precommit.R() > c.Round() {
		return constants.ErrFutureRoundMessage
	}
	if precommit.R() < c.Round() {
		// We are receiving a precommit for an old round. We need to check if we have now a quorum
		// in this old round.
		roundMessages := c.messages.GetOrCreate(precommit.R())
		roundMessages.AddPrecommit(precommit)
		go c.SendEvent(events.PowerChangeEvent{Height: c.Height().Uint64(), Round: c.Round(), Code: message.PrecommitCode, Value: precommit.Value()})

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
	go c.SendEvent(events.PowerChangeEvent{Height: c.Height().Uint64(), Round: c.Round(), Code: message.PrecommitCode, Value: precommit.Value()})
	c.LogPrecommitMessageEvent("MessageEvent(Precommit): Received", precommit)

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

func (c *Precommiter) LogPrecommitMessageEvent(message string, precommit *message.Precommit) {
	c.logger.Debug(message,
		"type", "Precommit",
		"local address", log.Lazy{Fn: func() string { return c.Address().String() }},
		"currentHeight", log.Lazy{Fn: c.Height},
		"msgHeight", precommit.H(),
		"currentRound", log.Lazy{Fn: c.Round},
		"msgRound", precommit.R(),
		"currentStep", c.step,
		"isProposer", log.Lazy{Fn: c.IsProposer},
		"currentProposer", log.Lazy{Fn: func() *types.CommitteeMember { return c.CommitteeSet().GetProposer(c.Round()) }},
		"isNilMsg", precommit.Value() == common.Hash{},
		"value", precommit.Value(),
		"totalVotes", log.Lazy{Fn: c.curRoundMessages.PrecommitsTotalPower},
		"totalNilVotes", log.Lazy{Fn: func() *big.Int { return c.curRoundMessages.PrecommitsPower(common.Hash{}) }},
		"proposedBlockVote", log.Lazy{Fn: func() *big.Int { return c.curRoundMessages.PrecommitsPower(c.curRoundMessages.ProposalHash()) }},
		"precommit", log.Lazy{Fn: func() string { return precommit.String() }},
	)
}
