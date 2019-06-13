package core

import (
	"math/big"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
)

func (c *core) sendPrecommit(isNil bool) {
	logger := c.logger.New("step", c.step)

	var precommit = &tendermint.Vote{
		Round:  big.NewInt(c.currentRoundState.round.Int64()),
		Height: big.NewInt(c.currentRoundState.Height().Int64()),
	}

	if isNil {
		precommit.ProposedBlockHash = common.Hash{}
	} else {
		precommit.ProposedBlockHash = c.currentRoundState.Proposal().ProposalBlock.Hash()
	}

	encodedVote, err := Encode(precommit)
	if err != nil {
		logger.Error("Failed to encode", "subject", precommit)
		return
	}

	c.logPrecommitMessageEvent("MessageEvent(Precommit): Sent", precommit, c.address.String(), "broadcast")

	c.broadcast(&message{
		Code: msgPrecommit,
		Msg:  encodedVote,
	})
}

// TODO: ensure to check the size of the committed seals as mentioned by Roberto in Correctness and Analysis of IBFT paper
func (c *core) handlePrecommit(msg *message, sender tendermint.Validator) error {
	var precommit *tendermint.Vote
	err := msg.Decode(&precommit)
	if err != nil {
		return errFailedDecodePrecommit
	}

	if err := c.checkMessage(precommit.Round, precommit.Height); err != nil {
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

	c.logPrecommitMessageEvent("MessageEvent(Precommit): Received", precommit, msg.Address.String(), c.address.String())

	// Line 47 in Algorithm 1 of The latest gossip on BFT consensus
	if !c.precommitTimeout.started && c.quorum(c.currentRoundState.Precommits.TotalSize(curProposaleHash)) {
		if err := c.stopPrecommitTimeout(); err != nil {
			return err
		}

		timeoutDuration := timeoutPrecommit(curR)
		c.precommitTimeout.scheduleTimeout(timeoutDuration, curR, curH, c.onTimeoutPrecommit)

		return nil
	}

	// Line 49 in Algorithm 1 of The latest gossip on BFT consensus
	if c.quorum(c.currentRoundState.Precommits.VotesSize(curProposaleHash)) {
		if err := c.stopPrecommitTimeout(); err != nil {
			return err
		}

		c.commit()
	}

	return errNoMajority
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

func (c *core) logPrecommitMessageEvent(message string, precommit *tendermint.Vote, from, to string) {
	currentProposalHash := c.currentRoundState.Proposal().ProposalBlock.Hash()
	c.logger.Info(message,
		"type", "Precommit",
		"from", from,
		"to", to,
		"currentHeight", c.currentRoundState.height,
		"msgHeight", precommit.Height,
		"currentRound", c.currentRoundState.round,
		"msgRound", precommit.Round,
		"currentStep", c.step,
		"isProposer", c.isProposer(),
		"currentProposer", c.valSet.GetProposer(),
		"isNilMsg", precommit.ProposedBlockHash == common.Hash{},
		"hash", precommit.ProposedBlockHash,
		"totalVotes", c.currentRoundState.Precommits.TotalSize(currentProposalHash),
		"totalNilVotes", c.currentRoundState.Precommits.NilVotesSize(),
		"quorumReject", c.quorum(c.currentRoundState.Precommits.NilVotesSize()),
		"totalNonNilVotes", c.currentRoundState.Precommits.VotesSize(currentProposalHash),
		"quorumAccept", c.quorum(c.currentRoundState.Precommits.VotesSize(currentProposalHash)),
	)
}
