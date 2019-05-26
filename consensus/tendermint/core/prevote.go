package core

import (
	"math/big"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
)

func (c *core) sendPrevote(isNil bool) {
	logger := c.logger.New("step", c.step)

	var vote = &tendermint.Vote{
		Round:  big.NewInt(c.currentRoundState.round.Int64()),
		Height: big.NewInt(c.currentRoundState.Height().Int64()),
	}

	if isNil {
		vote.ProposedBlockHash = common.Hash{}
	} else {
		vote.ProposedBlockHash = c.currentRoundState.Proposal().ProposalBlock.Hash()
	}

	encodedVote, err := Encode(vote)
	if err != nil {
		logger.Error("Failed to encode", "subject", vote)
		return
	}
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

	if err = c.checkMessage(msgPrevote, prevote.Round, prevote.Height); err != nil {
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

		logger.Info("Accepted Prevote", prevoteHash)

		// Line 34 in Algorithm 1 of The latest gossip on BFT consensus
		if c.step == StepProposeDone && !c.prevoteTimeout.started && c.quorum(c.currentRoundState.Prevotes.TotalSize(curProposaleHash)) {
			if err := c.stopPrevoteTimeout(); err == nil {
				timeoutDuration := timeoutPrevote(curR)
				c.prevoteTimeout.scheduleTimeout(timeoutDuration, curR, curH, c.onTimeoutPrevote)
			} else {
				return err
			}
			// Line 44 in Algorithm 1 of The latest gossip on BFT consensus
		} else if c.step == StepProposeDone && c.quorum(c.currentRoundState.Prevotes.NilVotesSize()) {
			if err := c.stopPrevoteTimeout(); err == nil {
				c.sendPrecommit(true)
				c.setStep(StepPrevoteDone)
			} else {
				return err
			}
			// Line 36 in Algorithm 1 of The latest gossip on BFT consensus
		} else if c.quorum(c.currentRoundState.Prevotes.VotesSize(curProposaleHash)) && !c.setValidRoundAndValue {
			// this piece of code should only run once
			if err := c.stopPrevoteTimeout(); err == nil {
				if c.step == StepProposeDone {
					c.lockedValue = &c.currentRoundState.Proposal().ProposalBlock
					c.lockedRound = big.NewInt(curR)
					c.sendPrecommit(false)
					c.setStep(StepPrevoteDone)
				}
				c.validValue = &c.currentRoundState.Proposal().ProposalBlock
				c.validRound = big.NewInt(curR)
				c.setValidRoundAndValue = true
			} else {
				return err
			}
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
