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

	logger.Info("MESSAGE: sent external message",
		"type", "precommit",
		"currentHeight", c.currentRoundState.height,
		"currentRound", c.currentRoundState.round,
		"currentStep", c.step,
		"from", c.address.String(),
		"currentProposer", c.isProposer(),
		"msgHeight", vote.Height,
		"msgRound", vote.Round,
		"isNilMsg", vote.ProposedBlockHash == common.Hash{},
		"message", vote,
	)

	c.broadcast(&message{
		Code: msgPrecommit,
		Msg:  encodedVote,
	})
}

// TODO: ensure to check the size of the committed seals as mentioned by Roberto in Correctness and Analysis of IBFT paper
func (c *core) handlePrecommit(msg *message, sender tendermint.Validator) error {
	logger := c.logger.New("from", sender, "step", c.step)

	var precommit *tendermint.Vote
	err := msg.Decode(&precommit)
	if err != nil {
		return errFailedDecodePrecommit
	}

	logger.Info("MESSAGE: got backlog message",
		"type", "precommit",
		"currentHeight", c.currentRoundState.height,
		"currentRound", c.currentRoundState.round,
		"currentStep", c.step,
		"from", msg.Address.String(),
		"sender", sender.Address().String(),
		"to", c.address.String(),
		"currentProposer", c.isProposer(),
		"msgHeight", precommit.Height,
		"msgRound", precommit.Round,
		"isNilMsg", precommit.ProposedBlockHash == common.Hash{},
		"message", precommit,
	)

	if err := c.checkMessage(precommit.Round, precommit.Height); err != nil {
		logger.Info("MESSAGE: backlog message ingored",
			"type", "precommit",
			"currentHeight", c.currentRoundState.height,
			"currentRound", c.currentRoundState.round,
			"currentStep", c.step,
			"from", msg.Address.String(),
			"sender", sender.Address().String(),
			"to", c.address.String(),
			"currentProposer", c.isProposer(),
			"msgHeight", precommit.Height,
			"msgRound", precommit.Round,
			"isNilMsg", precommit.ProposedBlockHash == common.Hash{},
			"message", precommit,
		)
		// We don't care about old round precommit messages, otherwise we would not be in a new round rather a new height
		return err
	}

	// We don't care about which step we are in to accept a precommit, since it has the highest importance
	precommitHash := precommit.ProposedBlockHash
	curProposaleHash := c.currentRoundState.Proposal().ProposalBlock.Hash()
	curR := c.currentRoundState.Round().Int64()
	curH := c.currentRoundState.Height().Int64()

	if precommitHash == (common.Hash{}) {
		c.currentRoundState.Precommits.AddNilVote(*msg)
	} else {
		c.currentRoundState.Precommits.AddVote(precommitHash, *msg)
	}

	logger.Info("Accepted PreCommit", "height", precommit.Height, "round", precommit.Round, "Hash", precommitHash, "quorumReject", c.quorum(c.currentRoundState.Precommits.NilVotesSize()), "totalNilVotes", c.currentRoundState.Precommits.NilVotesSize(), "quorumAccept", c.quorum(c.currentRoundState.Precommits.VotesSize(curProposaleHash)), "totalNonNilVotes", c.currentRoundState.Precommits.VotesSize(curProposaleHash))

	// Line 47 in Algorithm 1 of The latest gossip on BFT consensus
	if !c.precommitTimeout.started && c.quorum(c.currentRoundState.Precommits.NilVotesSize()) {
		timeoutDuration := timeoutPrecommit(curR)
		c.precommitTimeout.scheduleTimeout(timeoutDuration, curR, curH, c.onTimeoutPrecommit)
		// Line 49 in Algorithm 1 of The latest gossip on BFT consensus

		return nil
	}

	if !c.quorum(c.currentRoundState.Precommits.VotesSize(curProposaleHash)) {
		return errNoMajority
	}

	if err := c.stopPrecommitTimeout(); err != nil {
		return err
	}

	c.commit()

	return nil
}

func (c *core) handleCommit() {
	c.logger.Trace("Received a final committed proposal", "step", c.step)
	c.startRound(common.Big0)
}

func (c *core) stopPrecommitTimeout() error {
	if c.prevoteTimeout.started {
		if stopped := c.prevoteTimeout.stopTimer(); !stopped {
			return errMovedToNewRound
		}
	}
	return nil
}
