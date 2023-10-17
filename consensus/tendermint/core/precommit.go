package core

import (
	"context"
	"errors"
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
	}
	precommit := message.NewPrecommit(c.Round(), c.Height().Uint64(), value, c.backend.Sign)
	c.LogPrecommitMessageEvent("Precommit sent", precommit, c.address.String(), "broadcast")

	c.sentPrecommit = true
	c.Br().Broadcast(ctx, precommit)
}

func (c *Precommiter) HandlePrecommit(ctx context.Context, precommit *message.Precommit) error {
	if err := c.checkMessageStep(precommit.R(), precommit.H(), Precommit); err != nil {
		if errors.Is(err, constants.ErrOldRoundMessage) {
			// We are receiving a precommit for an old round. We need to check if we have now a quorum
			// in this old round.
			roundMessages := c.messages.GetOrCreate(precommit.R())
			roundMessages.AddPrecommit(precommit)
			oldRoundProposal := roundMessages.Proposal()
			if oldRoundProposal != nil && roundMessages.PrecommitsPower(oldRoundProposal.Block().Hash()).Cmp(c.CommitteeSet().Quorum()) >= 0 {
				c.logger.Info("Quorum on a old round proposal", "round", precommit.R())
				if !roundMessages.IsProposalVerified() {
					if _, err2 := c.backend.VerifyProposal(roundMessages.Proposal().Block()); err2 != nil {
						// Impossible with the BFT assumptions of 1/3rd honest.
						return err2
					}
				}
				c.Commit(precommit.R(), c.curRoundMessages)
				return nil
			}
		}

		return err
	}

	// Line 49 in Algorithm 1 of The latest gossip on BFT consensus
	curProposalHash := c.curRoundMessages.ProposalHash()
	// We don't care about which step we are in to accept a precommit, since it has the highest importance
	c.curRoundMessages.AddPrecommit(precommit)
	c.LogPrecommitMessageEvent("MessageEvent(Precommit): Received", precommit, precommit.Sender().String(), c.address.String())
	if curProposalHash != (common.Hash{}) && c.curRoundMessages.PrecommitsPower(curProposalHash).Cmp(c.CommitteeSet().Quorum()) >= 0 {
		if err := c.precommitTimeout.StopTimer(); err != nil {
			return err
		}
		c.logger.Debug("Stopped Scheduled Precommit Timeout")

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			c.Commit(c.Round(), c.curRoundMessages)
		}

		// Line 47 in Algorithm 1 of The latest gossip on BFT consensus
	} else if !c.precommitTimeout.TimerStarted() && c.curRoundMessages.PrecommitsTotalPower().Cmp(c.CommitteeSet().Quorum()) >= 0 {
		timeoutDuration := c.timeoutPrecommit(c.Round())
		c.precommitTimeout.ScheduleTimeout(timeoutDuration, c.Round(), c.Height(), c.onTimeoutPrecommit)
		c.logger.Debug("Scheduled Precommit Timeout", "Timeout Duration", timeoutDuration)
	}

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
