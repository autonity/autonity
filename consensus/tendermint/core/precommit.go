package core

import (
	"math/big"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
)

func (c *core) sendPrecommit(isNil bool) {
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
		Code: msgPrecommit,
		Msg:  encodedVote,
	})
}

func (c *core) sendPrecommitForOldBlock(r *big.Int, h *big.Int, digest common.Hash) {
	sub := &tendermint.Vote{
		Round:             r,
		Height:            h,
		ProposedBlockHash: digest,
	}
	c.broadcastPrecommit(sub)
}

func (c *core) broadcastPrecommit(sub *tendermint.Vote) {
	logger := c.logger.New("step", c.step)

	encodedSubject, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", "subject", sub)
		return
	}
	c.broadcast(&message{
		Code: msgPrecommit,
		Msg:  encodedSubject,
	})
}

func (c *core) handlePrecommit(msg *message, src tendermint.Validator) error {
	// Decode COMMIT message
	var precommit *tendermint.Vote
	err := msg.Decode(&precommit)
	if err != nil {
		return errFailedDecodePrecommit
	}

	if err := c.checkMessage(msgPrecommit, precommit.Round, precommit.Height); err != nil {
		// We don't care about old round messages, so if the message is in the correct height and round it should be fine
		return err
	}
	// TODO: manage precommit timer
	// We don't care about which step we are in to accept a precommit, since it has the highest importance
	precommitHash := precommit.ProposedBlockHash
	curProposaleHash := c.currentRoundState.Proposal().ProposalBlock.Hash()

	if precommitHash == (common.Hash{}) {
		c.currentRoundState.Precommits.AddNilVote(*msg)
	} else {
		c.currentRoundState.Precommits.AddVote(precommitHash, *msg)
	}

	if !c.precommitTimeout.started && c.quorum(c.currentRoundState.Precommits.NilVotesSize()) {
		timeoutDuration := timeoutPrecommit(c.currentRoundState.Round().Int64())
		c.precommitTimeout.scheduleTimeout(timeoutDuration, c.currentRoundState.Round().Int64(), c.currentRoundState.Height().Int64(), c.onTimeoutPrecommit)
	} else if c.quorum(c.currentRoundState.Precommits.VotesSize(curProposaleHash)) {
		c.commit()
	}

	return nil
}

func (c *core) handleCommit() {
	c.logger.Trace("Received a final committed proposal", "step", c.step)
	c.startRound(common.Big0)
}
