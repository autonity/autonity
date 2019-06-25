package core

import (
	"math/big"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
)

func (c *core) sendPrecommit(isNil bool) {
	logger := c.logger.New("step", c.step)

	var precommit = &tendermint.Vote{
		Round:  big.NewInt(c.currentRoundState.Round().Int64()),
		Height: big.NewInt(c.currentRoundState.Height().Int64()),
	}

	if isNil {
		precommit.ProposedBlockHash = common.Hash{}
	} else {
		if h := c.currentRoundState.GetCurrentProposalHash(); h == (common.Hash{}) {
			c.logger.Error("core.sendPrecommit Proposal is empty! It should not be empty!")
			return
		}
		precommit.ProposedBlockHash = c.currentRoundState.GetCurrentProposalHash()
	}

	encodedVote, err := Encode(precommit)
	if err != nil {
		logger.Error("Failed to encode", "subject", precommit)
		return
	}

	c.logPrecommitMessageEvent("MessageEvent(Precommit): Sent", precommit, c.address.String(), "broadcast")

	c.sentPrecommit = true
	c.broadcast(&message{
		Code: msgPrecommit,
		Msg:  encodedVote,
	})
}

// TODO: ensure to check the size of the committed seals as mentioned by Roberto in Correctness and Analysis of IBFT paper
func (c *core) handlePrecommit(msg *message) error {
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
	curProposalHash := c.currentRoundState.GetCurrentProposalHash()
	curR := c.currentRoundState.Round().Int64()
	curH := c.currentRoundState.Height().Int64()

	if precommitHash == (common.Hash{}) {
		c.currentRoundState.Precommits.AddNilVote(*msg)
	} else {
		c.currentRoundState.Precommits.AddVote(precommitHash, *msg)
	}

	c.logPrecommitMessageEvent("MessageEvent(Precommit): Received", precommit, msg.Address.String(), c.address.String())

	// Line 49 in Algorithm 1 of The latest gossip on BFT consensus
	if curProposalHash != (common.Hash{}) && c.quorum(c.currentRoundState.Precommits.VotesSize(curProposalHash)) {
		if err := c.stopPrecommitTimeout(); err != nil {
			return err
		}

		c.commit()

		// Line 47 in Algorithm 1 of The latest gossip on BFT consensus
	} else if !c.precommitTimeout.started && c.quorum(c.currentRoundState.Precommits.TotalSize()) {
		timeoutDuration := timeoutPrecommit(curR)
		c.precommitTimeout.scheduleTimeout(timeoutDuration, curR, curH, c.onTimeoutPrecommit)
		c.logger.Debug("Scheduled Precommit Timeout", "Timeout Duration", timeoutDuration)
	}

	return nil
}

func (c *core) handleCommit() {
	c.logger.Trace("Received a final committed proposal", "step", c.step)
	c.startRound(common.Big0)
}

func (c *core) stopPrecommitTimeout() error {
	if c.precommitTimeout.started {
		c.logger.Debug("Stopping Scheduled Precommit Timeout")
		if stopped := c.precommitTimeout.stopTimer(); !stopped {
			return errMovedToNewRound
		}
	}
	return nil
}

func (c *core) logPrecommitMessageEvent(message string, precommit *tendermint.Vote, from, to string) {
	currentProposalHash := c.currentRoundState.GetCurrentProposalHash()
	c.logger.Debug(message,
		"from", from,
		"to", to,
		"currentHeight", c.currentRoundState.Height(),
		"msgHeight", precommit.Height,
		"currentRound", c.currentRoundState.Round(),
		"msgRound", precommit.Round,
		"currentStep", c.step,
		"isProposer", c.isProposer(),
		"currentProposer", c.valSet.GetProposer(),
		"isNilMsg", precommit.ProposedBlockHash == common.Hash{},
		"hash", precommit.ProposedBlockHash,
		"type", "Precommit",
		"totalVotes", c.currentRoundState.Precommits.TotalSize(),
		"totalNilVotes", c.currentRoundState.Precommits.NilVotesSize(),
		"quorumReject", c.quorum(c.currentRoundState.Precommits.NilVotesSize()),
		"totalNonNilVotes", c.currentRoundState.Precommits.VotesSize(currentProposalHash),
		"quorumAccept", c.quorum(c.currentRoundState.Precommits.VotesSize(currentProposalHash)),
	)
}
