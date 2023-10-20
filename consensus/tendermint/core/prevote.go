package core

import (
	"context"

	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/core/types"
	"github.com/autonity/autonity/rlp"

	"github.com/autonity/autonity/common"
)

type Prevoter struct {
	*Core
}

func (c *Prevoter) SendPrevote(ctx context.Context, isNil bool) {
	logger := c.logger.New("step", c.step)

	var prevote = &message.Vote{
		Round:  c.Round(),
		Height: c.Height(),
	}

	if isNil {
		prevote.ProposedBlockHash = common.Hash{}
	} else {
		if h := c.curRoundMessages.GetProposalHash(); h == (common.Hash{}) {
			c.logger.Error("sendPrevote Proposal is empty! It should not be empty!")
			return
		}
		prevote.ProposedBlockHash = c.curRoundMessages.GetProposalHash()
	}

	encodedVote, err := rlp.EncodeToBytes(&prevote)
	if err != nil {
		logger.Error("Failed to encode", "subject", prevote)
		return
	}

	c.LogPrevoteMessageEvent("MessageEvent(Prevote): Sent", prevote, c.address.String(), "broadcast")

	c.sentPrevote = true
	c.Br().SignAndBroadcast(ctx, &message.Message{
		Code:          message.MsgPrevote,
		Payload:       encodedVote,
		Address:       c.address,
		CommittedSeal: []byte{},
	})
}

func (c *Prevoter) HandlePrevote(ctx context.Context, msg *message.Message) error {
	preVote := msg.ConsensusMsg.(*message.Vote)
	if err := c.CheckMessage(int64(preVote.Round), preVote.Height.Uint64(), types.Prevote); err != nil {
		// Store old round prevote messages for future rounds since it is required for validRound
		if err == constants.ErrOldRoundMessage {
			// We only process old rounds while future rounds messages are pushed on to the backlog
			oldRoundMessages := c.messages.GetOrCreate(int64(preVote.Round))
			c.AcceptVote(oldRoundMessages, types.Prevote, preVote.ProposedBlockHash, *msg)

			// Line 28 in Algorithm 1 of The latest gossip on BFT consensus.
			if c.step == types.Propose {
				// ProposalBlock would be nil if node haven't receive proposal yet.
				if c.curRoundMessages.ProposalDetails.ProposalBlock != nil {
					vr := c.curRoundMessages.ProposalDetails.ValidRound
					h := c.curRoundMessages.ProposalDetails.ProposalBlock.Hash()
					rs := c.messages.GetOrCreate(vr)

					if vr >= 0 && vr < c.Round() && rs.PrevotesPower(h).Cmp(c.CommitteeSet().Quorum()) >= 0 {
						c.SendPrevote(ctx, !(c.lockedRound <= vr || h == c.lockedValue.Hash()))
						c.SetStep(types.Prevote)
						return nil
					}
				}
			}
		}
		return err
	}

	// After checking the message we know it is from the same height and round, so we should store it even if
	// c.curRoundMessages.Step() < prevote. The propose Timeout which is started at the beginning of the round
	// will update the step to at least prevote and when it handle its on preVote(nil), then it will also have
	// votes from other nodes.
	prevoteHash := preVote.ProposedBlockHash
	c.AcceptVote(c.curRoundMessages, types.Prevote, prevoteHash, *msg)

	c.LogPrevoteMessageEvent("MessageEvent(Prevote): Received", preVote, msg.Address.String(), c.address.String())

	// Now we can add the preVote to our current round state
	if c.step >= types.Prevote {
		curProposalHash := c.curRoundMessages.GetProposalHash()

		// Line 36 in Algorithm 1 of The latest gossip on BFT consensus
		if curProposalHash != (common.Hash{}) && c.curRoundMessages.PrevotesPower(curProposalHash).Cmp(c.CommitteeSet().Quorum()) >= 0 && !c.setValidRoundAndValue {
			// this piece of code should only run once
			if err := c.prevoteTimeout.StopTimer(); err != nil {
				return err
			}
			c.logger.Debug("Stopped Scheduled Prevote Timeout")

			if c.step == types.Prevote {
				c.lockedValue = c.curRoundMessages.Proposal().ProposalBlock
				c.lockedRound = c.Round()
				c.precommiter.SendPrecommit(ctx, false)
				c.SetStep(types.Precommit)
			}
			c.validValue = c.curRoundMessages.Proposal().ProposalBlock
			c.validRound = c.Round()
			c.setValidRoundAndValue = true
			// Line 44 in Algorithm 1 of The latest gossip on BFT consensus
		} else if c.step == types.Prevote && c.curRoundMessages.PrevotesPower(common.Hash{}).Cmp(c.CommitteeSet().Quorum()) >= 0 {
			if err := c.prevoteTimeout.StopTimer(); err != nil {
				return err
			}
			c.logger.Debug("Stopped Scheduled Prevote Timeout")

			c.precommiter.SendPrecommit(ctx, true)
			c.SetStep(types.Precommit)

			// Line 34 in Algorithm 1 of The latest gossip on BFT consensus
		} else if c.step == types.Prevote && !c.prevoteTimeout.TimerStarted() && !c.sentPrecommit && c.curRoundMessages.PrevotesTotalPower().Cmp(c.CommitteeSet().Quorum()) >= 0 {
			timeoutDuration := c.timeoutPrevote(c.Round())
			c.prevoteTimeout.ScheduleTimeout(timeoutDuration, c.Round(), c.Height(), c.onTimeoutPrevote)
			c.logger.Debug("Scheduled Prevote Timeout", "Timeout Duration", timeoutDuration)
		}
	}

	return nil
}

func (c *Prevoter) LogPrevoteMessageEvent(message string, prevote *message.Vote, from, to string) {
	currentProposalHash := c.curRoundMessages.GetProposalHash()
	c.logger.Debug(message,
		"from", from,
		"to", to,
		"currentHeight", c.Height(),
		"msgHeight", prevote.Height,
		"currentRound", c.Round(),
		"msgRound", prevote.Round,
		"currentStep", c.step,
		"isProposer", c.IsProposer(),
		"currentProposer", c.CommitteeSet().GetProposer(c.Round()),
		"isNilMsg", prevote.ProposedBlockHash == common.Hash{},
		"hash", prevote.ProposedBlockHash,
		"type", "Prevote",
		"totalVotes", c.curRoundMessages.PrevotesTotalPower(),
		"totalNilVotes", c.curRoundMessages.PrevotesPower(common.Hash{}),
		"quorum", c.committee.Quorum(),
		"VoteProposedBlock", c.curRoundMessages.PrevotesPower(currentProposalHash),
	)
}
