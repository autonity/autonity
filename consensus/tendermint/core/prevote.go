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

// TODO: possibly need to add sentPrevote and sentPrecommit
func (c *core) handlePrevote(msg *message, src tendermint.Validator) error {
	// Decode PREPARE message
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
	// TODO: manage prevote timer
	if c.step >= StepProposeDone {
		prevoteHash := prevote.ProposedBlockHash
		curProposaleHash := c.currentRoundState.Proposal().ProposalBlock.Hash()

		if prevoteHash == (common.Hash{}) {
			c.currentRoundState.Prevotes.AddNilVote(*msg)
		} else {
			c.currentRoundState.Prevotes.AddVote(prevoteHash, *msg)
		}

		if c.step == StepProposeDone && !c.prevoteTimeout.started && c.quorum(c.currentRoundState.Prevotes.TotalSize(curProposaleHash)) {
			timeoutDuration := timeoutPrevote(c.currentRoundState.Round().Int64())
			c.prevoteTimeout.scheduleTimeout(timeoutDuration, c.currentRoundState.Round().Int64(), c.currentRoundState.Height().Int64(), c.onTimeoutPrevote)
		} else if c.step == StepProposeDone && c.quorum(c.currentRoundState.Prevotes.NilVotesSize()) {
			// TODO: probably need to stop timer, same in the other if branches need to fix this
			c.sendPrecommit(true)
			c.setStep(StepPrevoteDone)
		} else if c.quorum(c.currentRoundState.Prevotes.VotesSize(curProposaleHash)) && !c.setValidRoundAndValue {
			// this piece of code should only run once
			if c.step == StepProposeDone {
				c.lockedValue = &c.currentRoundState.Proposal().ProposalBlock
				c.lockedRound = big.NewInt(c.currentRoundState.Round().Int64())
				c.sendPrecommit(false)
				c.setStep(StepPrevoteDone)
			}
			c.validValue = &c.currentRoundState.Proposal().ProposalBlock
			c.validRound = big.NewInt(c.currentRoundState.Round().Int64())
			c.setValidRoundAndValue = true
		}
	}

	return nil
}
