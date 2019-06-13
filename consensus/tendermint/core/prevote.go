package core

import (
	"math/big"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
)

func (c *core) sendPrevote(isNil bool) {
	logger := c.logger.New("step", c.step)

	var prevote = &tendermint.Vote{
		Round:  big.NewInt(c.currentRoundState.round.Int64()),
		Height: big.NewInt(c.currentRoundState.Height().Int64()),
	}

	if isNil {
		prevote.ProposedBlockHash = common.Hash{}
	} else {
		prevote.ProposedBlockHash = c.currentRoundState.Proposal().ProposalBlock.Hash()
	}

	encodedVote, err := Encode(prevote)
	if err != nil {
		logger.Error("Failed to encode", "subject", prevote)
		return
	}

	c.logPrevoteMessageEvent("MessageEvent(Prevote): Sent", prevote, c.address.String())

	c.broadcast(&message{
		Code: msgPrevote,
		Msg:  encodedVote,
	})
}

func (c *core) handlePrevote(msg *message, sender tendermint.Validator) error {
	logger := c.logger.New("from", sender, "step", c.step)

	var prevote *tendermint.Vote
	err := msg.Decode(&prevote)
	if err != nil {
		return errFailedDecodePrevote
	}

	c.logPrevoteMessageEvent("MessageEvent(Prevote): Received", prevote, msg.Address.String())

	if err = c.checkMessage(prevote.Round, prevote.Height); err != nil {
		// We want to store old round messages for future rounds since it is required for validRound
		if err == errOldRoundMessage {
			// The roundstate must exist as every roundstate is added to c.currentHeightRoundsState at startRound
			// And we only process old rounds while future rounds messages are pushed on to the backlog
			prevoteMS := c.currentHeightRoundsStates[prevote.Round.Int64()].Prevotes
			if prevote.ProposedBlockHash == (common.Hash{}) {
				prevoteMS.AddNilVote(*msg)
			} else {
				prevoteMS.AddVote(prevote.ProposedBlockHash, *msg)
			}
		}
		return err
	}

	// Now we can add the prevote to our current round state
	if c.step >= StepProposeDone {
		prevoteHash := prevote.ProposedBlockHash
		curProposaleHash := c.currentRoundState.Proposal().ProposalBlock.Hash()
		curR := c.currentRoundState.Round().Int64()
		curH := c.currentRoundState.Height().Int64()

		if prevoteHash == (common.Hash{}) {
			c.currentRoundState.Prevotes.AddNilVote(*msg)
		} else {
			c.currentRoundState.Prevotes.AddVote(prevoteHash, *msg)
		}

		logger.Info("Accepted Prevote", "height", prevote.Height, "round", prevote.Round, "Hash", prevoteHash, "quorumReject", c.quorum(c.currentRoundState.Prevotes.NilVotesSize()), "totalNilVotes", c.currentRoundState.Prevotes.NilVotesSize(), "quorumAccept", c.quorum(c.currentRoundState.Prevotes.TotalSize(curProposaleHash)), "totalNonNilVotes", c.currentRoundState.Prevotes.TotalSize(curProposaleHash))

		// Line 34 in Algorithm 1 of The latest gossip on BFT consensus
		if c.step == StepProposeDone && !c.prevoteTimeout.started && c.quorum(c.currentRoundState.Prevotes.TotalSize(curProposaleHash)) {
			if err := c.stopPrevoteTimeout(); err != nil {
				return err
			}

			timeoutDuration := timeoutPrevote(curR)
			c.prevoteTimeout.scheduleTimeout(timeoutDuration, curR, curH, c.onTimeoutPrevote)
			// Line 44 in Algorithm 1 of The latest gossip on BFT consensus
		} else if c.step == StepProposeDone && c.quorum(c.currentRoundState.Prevotes.NilVotesSize()) {
			if err := c.stopPrevoteTimeout(); err != nil {
				return err
			}
			c.sendPrecommit(true)
			c.setStep(StepPrevoteDone)
			// Line 36 in Algorithm 1 of The latest gossip on BFT consensus

		} else if c.quorum(c.currentRoundState.Prevotes.VotesSize(curProposaleHash)) && !c.setValidRoundAndValue {
			// this piece of code should only run once
			if err := c.stopPrevoteTimeout(); err != nil {
				return err
			}

			if c.step == StepProposeDone {
				c.lockedValue = c.currentRoundState.Proposal().ProposalBlock
				c.lockedRound = big.NewInt(curR)
				c.sendPrecommit(false)
				c.setStep(StepPrevoteDone)
			}
			c.validValue = c.currentRoundState.Proposal().ProposalBlock
			c.validRound = big.NewInt(curR)
			c.setValidRoundAndValue = true
		} else {
			return errNoMajority
		}
	}

	return nil
}

func (c *core) stopPrevoteTimeout() error {
	if c.prevoteTimeout.started {
		if stopped := c.prevoteTimeout.stopTimer(); !stopped {
			return errNilPrecommitSent
		}
	}
	return nil
}

func (c *core) logPrevoteMessageEvent(message string, prevote *tendermint.Vote, from string) {
	c.logger.Info(message,
		"from", from,
		"type", "Prevote",
		"currentHeight", c.currentRoundState.height,
		"msgHeight", prevote.Height,
		"currentRound", c.currentRoundState.round,
		"msgRound", prevote.Round,
		"currentSteo", c.step,
		"msgStep", c.step,
		"currentProposer", c.valSet.GetProposer(),
		"isNilMsg", prevote.ProposedBlockHash == common.Hash{},
	)
}
