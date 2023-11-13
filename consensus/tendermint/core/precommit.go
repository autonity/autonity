package core

import (
	"bytes"
	"context"
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/helpers"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	tctypes "github.com/autonity/autonity/consensus/tendermint/core/types"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/rlp"
)

type Precommiter struct {
	*Core
}

func (c *Precommiter) SendPrecommit(ctx context.Context, isNil bool) {
	logger := c.logger.New("step", c.step)

	var precommit = &message.Vote{
		Round:  c.Round(),
		Height: c.Height(),
	}

	if isNil {
		precommit.ProposedBlockHash = common.Hash{}
	} else {
		if h := c.curRoundMessages.GetProposalHash(); h == (common.Hash{}) {
			c.logger.Error("Core.sendPrecommit Proposal is empty! It should not be empty!")
			return
		}
		precommit.ProposedBlockHash = c.curRoundMessages.GetProposalHash()
	}

	encodedVote, err := rlp.EncodeToBytes(&precommit)
	if err != nil {
		logger.Error("Failed to encode", "subject", precommit)
		return
	}

	c.LogPrecommitMessageEvent("MessageEvent(Precommit): Sent", precommit, c.address.String(), "broadcast")

	msg := &message.Message{
		Code:          consensus.MsgPrecommit,
		Payload:       encodedVote,
		Address:       c.address,
		CommittedSeal: []byte{},
	}

	// Create committed seal
	seal := helpers.PrepareCommittedSeal(precommit.ProposedBlockHash, c.Round(), c.Height())
	msg.CommittedSeal, err = c.backend.Sign(seal)
	if err != nil {
		c.logger.Error("Core.sendPrecommit error while signing committed seal", "err", err)
	}

	c.sentPrecommit = true
	c.Broadcaster().SignAndBroadcast(ctx, msg)
}

func (c *Precommiter) HandlePrecommit(ctx context.Context, msg *message.Message) error {
	preCommit := msg.ConsensusMsg.(*message.Vote)
	precommitHash := preCommit.ProposedBlockHash
	if err := c.CheckMessage(preCommit.Round, preCommit.Height.Uint64(), tctypes.Precommit); err != nil {
		if err == constants.ErrOldRoundMessage {
			roundMsgs := c.messages.GetOrCreate(preCommit.Round)
			if err2 := c.VerifyCommittedSeal(msg.Address, append([]byte(nil), msg.CommittedSeal...), preCommit.ProposedBlockHash, preCommit.Round, preCommit.Height); err2 != nil {
				return err2
			}
			c.AcceptVote(roundMsgs, tctypes.Precommit, precommitHash, *msg)
			oldRoundProposalHash := roundMsgs.GetProposalHash()
			if oldRoundProposalHash != (common.Hash{}) && roundMsgs.PrecommitsPower(oldRoundProposalHash).Cmp(c.CommitteeSet().Quorum()) >= 0 {
				c.logger.Info("Quorum on a old round proposal", "round", preCommit.Round)
				if !roundMsgs.IsProposalVerified() {
					if _, err2 := c.backend.VerifyProposal(roundMsgs.Proposal().ProposalBlock); err2 != nil {
						return err2
					}
				}
				c.Commit(preCommit.Round, c.curRoundMessages)
				return nil
			}
		}

		return err
	}

	// Don't want to decode twice, hence sending preCommit with message
	if err := c.VerifyCommittedSeal(msg.Address, append([]byte(nil), msg.CommittedSeal...), preCommit.ProposedBlockHash, preCommit.Round, preCommit.Height); err != nil {
		return err
	}
	// Line 49 in Algorithm 1 of The latest gossip on BFT consensus
	curProposalHash := c.curRoundMessages.GetProposalHash()
	// We don't care about which step we are in to accept a preCommit, since it has the highest importance

	c.AcceptVote(c.curRoundMessages, tctypes.Precommit, precommitHash, *msg)
	c.LogPrecommitMessageEvent("MessageEvent(Precommit): Received", preCommit, msg.Address.String(), c.address.String())
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

func (c *Precommiter) VerifyCommittedSeal(addressMsg common.Address, committedSealMsg []byte, proposedBlockHash common.Hash, round int64, height *big.Int) error {
	committedSeal := helpers.PrepareCommittedSeal(proposedBlockHash, round, height)

	sealerAddress, err := types.GetSignatureAddress(committedSeal, committedSealMsg)
	if err != nil {
		c.logger.Error("Failed to get signer address", "err", err)
		return err
	}

	// ensure sender signed the committed seal
	if !bytes.Equal(sealerAddress.Bytes(), addressMsg.Bytes()) {
		c.logger.Error("verify precommit seal error", "got", addressMsg.String(), "expected", sealerAddress.String())

		return constants.ErrInvalidSenderOfCommittedSeal
	}

	return nil
}

func (c *Precommiter) HandleCommit(ctx context.Context) {
	c.logger.Debug("Received a final committed proposal", "step", c.step)
	lastBlock, _ := c.backend.HeadBlock()
	height := new(big.Int).Add(lastBlock.Number(), common.Big1)
	if height.Cmp(c.Height()) == 0 {
		c.logger.Debug("Discarding event as Core is at the same height", "height", c.Height())
	} else {
		c.logger.Debug("New chain head ahead of consensus Core height", "height", c.Height(), "block_height", height)
		c.StartRound(ctx, 0)
	}
}

func (c *Precommiter) LogPrecommitMessageEvent(message string, precommit *message.Vote, from, to string) {
	currentProposalHash := c.curRoundMessages.GetProposalHash()
	c.logger.Debug(message,
		"from", from,
		"to", to,
		"currentHeight", c.Height(),
		"msgHeight", precommit.Height,
		"currentRound", c.Round(),
		"msgRound", precommit.Round,
		"currentStep", c.step,
		"isProposer", c.IsProposer(),
		"currentProposer", c.CommitteeSet().GetProposer(c.Round()),
		"isNilMsg", precommit.ProposedBlockHash == common.Hash{},
		"hash", precommit.ProposedBlockHash,
		"type", "Precommit",
		"totalVotes", c.curRoundMessages.PrecommitsTotalPower(),
		"totalNilVotes", c.curRoundMessages.PrecommitsPower(common.Hash{}),
		"proposedBlockVote", c.curRoundMessages.PrecommitsPower(currentProposalHash),
	)
}
