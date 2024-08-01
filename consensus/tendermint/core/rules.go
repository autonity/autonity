package core

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/metrics"
)

// Line 22 in Algorithm 1 of The latest gossip on BFT consensus
// checks if we can prevote for a new proposal
// Assumes:
// 1. proposal is for current height
// 2. proposal is for current round
// 3. proposal is from current proposer
// 4. proposal is verified and valid
func (c *Core) newProposalCheck(ctx context.Context, proposal *message.Propose) {
	if c.Step() != Propose || proposal.ValidRound() != -1 {
		return
	}
	// When lockedRound is set to any value other than -1 lockedValue is also
	// set to a non nil value. So we can be sure that we will only try to access
	// lockedValue when it is non nil.
	c.prevoter.SendPrevote(ctx, !(c.LockedRound() == -1 || proposal.Block().Hash() == c.LockedValue().Hash()))
	c.SetStep(ctx, Prevote)
}

// Line 28 in Algorithm 1 of The latest gossip on BFT consensus
// checks if we can prevote for an old proposal
// Assumes:
// 1. proposal is for current height
// 2. proposal is for current round
// 3. proposal is from current proposer
// 4. proposal is verified and valid
func (c *Core) oldProposalCheck(ctx context.Context, proposal *message.Propose) {
	vr := proposal.ValidRound()
	if c.Step() != Propose || vr == -1 || vr >= c.Round() {
		return
	}

	hash := proposal.Block().Hash()
	rm := c.roundsState.GetOrCreate(vr)
	if rm.PrevotesPower(hash).Cmp(c.CommitteeSet().Quorum()) >= 0 {
		c.prevoter.SendPrevote(ctx, !(c.LockedRound() <= vr || hash == c.LockedValue().Hash()))
		c.SetStep(ctx, Prevote)
	}
}

// Line 34 in Algorithm 1 of The latest gossip on BFT consensus
// checks if we have to schedule the prevote timeout
func (c *Core) prevoteTimeoutCheck() {
	if c.Step() != Prevote {
		return
	}
	if !c.prevoteTimeout.TimerStarted() && !c.SentPrecommit() && c.CurRoundMessages().PrevotesTotalPower().Cmp(c.CommitteeSet().Quorum()) >= 0 {
		timeoutDuration := c.timeoutPrevote(c.Round())
		c.prevoteTimeout.ScheduleTimeout(timeoutDuration, c.Round(), c.Height(), c.onTimeoutPrevote)
		c.logger.Debug("Scheduled Prevote Timeout", "Timeout Duration", timeoutDuration)
	}
}

// Line 36 in Algorithm 1 of The latest gossip on BFT consensus
// Assumes:
// 1. proposal is for current height
// 2. proposal is for current round
// 3. proposal is from current proposer
// 4. proposal is verified and valid
func (c *Core) quorumPrevotesCheck(ctx context.Context, proposal *message.Propose) {
	if c.Step() == Propose {
		return
	}
	// we are at prevote or precommit step
	if c.CurRoundMessages().PrevotesPower(proposal.Block().Hash()).Cmp(c.CommitteeSet().Quorum()) >= 0 && !c.ValidRoundAndValueSet() {
		if metrics.Enabled {
			PrevoteQuorumBlockTSDeltaBg.Add(time.Since(c.currBlockTimeStamp).Nanoseconds())
		}
		if c.Step() == Prevote {
			// todo: Jason, double check if the locked round equals to the round of the proposal.
			c.SetLockedRoundAndValue(c.Round(), proposal.Block())
			c.precommiter.SendPrecommit(ctx, false)
			c.SetStep(ctx, Precommit)
		}
		c.SetValidRoundAndValue(c.Round(), proposal.Block())
	}
}

// Line 44 in Algorithm 1 of The latest gossip on BFT consensus
// checks if we have to precommit nil because we received quorum prevotes nil
func (c *Core) quorumPrevotesNilCheck(ctx context.Context) {
	if c.Step() != Prevote {
		return
	}
	if c.CurRoundMessages().PrevotesPower(common.Hash{}).Cmp(c.CommitteeSet().Quorum()) >= 0 {
		if metrics.Enabled {
			PrevoteQuorumBlockTSDeltaBg.Add(time.Since(c.currBlockTimeStamp).Nanoseconds())
		}
		c.precommiter.SendPrecommit(ctx, true)
		c.SetStep(ctx, Precommit)
	}
}

// Line 47 in Algorithm 1 of The latest gossip on BFT consensus
// checks if we have to schedule the precommit timeout
func (c *Core) precommitTimeoutCheck() {
	if !c.precommitTimeout.TimerStarted() && c.CurRoundMessages().PrecommitsTotalPower().Cmp(c.CommitteeSet().Quorum()) >= 0 {
		timeoutDuration := c.timeoutPrecommit(c.Round())
		c.precommitTimeout.ScheduleTimeout(timeoutDuration, c.Round(), c.Height(), c.onTimeoutPrecommit)
		c.logger.Debug("Scheduled Precommit Timeout", "Timeout Duration", timeoutDuration)
	}
}

// Line 49 in Algorithm 1 of The latest gossip on BFT consensus
// Assumes:
// 1. proposal is for current height
// 2. proposal is from correct proposer from (h,r) of the proposal
// returns whether the proposal was committed or not
func (c *Core) quorumPrecommitsCheck(ctx context.Context, proposal *message.Propose, verified bool) bool {
	hash := proposal.Block().Hash()
	rm := c.roundsState.GetOrCreate(proposal.R())

	// if no quorum, return without verifying the proposal
	if rm.PrecommitsPower(hash).Cmp(c.CommitteeSet().Quorum()) < 0 {
		return false
	}
	if metrics.Enabled {
		PrecommitQuorumBlockTSDeltaBg.Add(time.Since(c.currBlockTimeStamp).Nanoseconds())
	}

	// if there is a quorum, verify the proposal if needed
	if !verified {
		if _, err := c.backend.VerifyProposal(proposal.Block()); err != nil {
			// This can happen if while we are processing the proposal,
			// we actually receive the finalized proposed block from p2p block propagation (other peers already reached quorum on it)
			// In this case we can just consider the proposal as committed.
			if errors.Is(err, core.ErrKnownBlock) || errors.Is(err, constants.ErrAlreadyHaveBlock) {
				c.logger.Info("Verified proposal that was already in our local chain", "err", err)
				c.SetStep(ctx, PrecommitDone) // we do not need to process any more consensus messages for this height
				return true
			}
			// Impossible with the BFT assumptions of 1/3rd honest.
			panic("Fatal Safety Error: Quorum on unverifiable proposal. err: " + err.Error())
		}
	}

	// all good, commit
	c.logger.Debug("Committing proposal", "height", c.Height(), "round", c.Round(), "proposal round", proposal.R())
	c.Commit(ctx, proposal.R(), rm)
	return true
}

// Line 55 in Algorithm 1 of The latest gossip on BFT consensus
// check if we need to skip to a new round
func (c *Core) roundSkipCheck(ctx context.Context, r int64) {
	futurePower := new(big.Int)

	c.futureRoundLock.RLock()
	futureAggregatedPower, ok := c.futurePower[r]
	if ok {
		futurePower.Set(futureAggregatedPower.Power())
	}
	c.futureRoundLock.RUnlock()

	if futurePower.Cmp(c.CommitteeSet().F()) > 0 {
		c.logger.Debug("Received messages with F + 1 total power for a higher round", "New round", r)
		c.StartRound(ctx, r, false)
	}
}

// -------------------------------------
// UTILITY FUNCTIONS
// These functions do not directly map to tendermint upon conditions
// but instead group together multiple conditions that need to be checked
// based on the type of msg received (propose, prevote or precommit)
// -------------------------------------

// upon condition rules to check when receiving a current round proposal
func (c *Core) currentProposalChecks(ctx context.Context, proposal *message.Propose) {
	// Line 49 in Algorithm 1 of The latest gossip on BFT consensus
	// check if we have a quorum of precommits for this proposal.
	// If so, no need to check the other rules
	if committed := c.quorumPrecommitsCheck(ctx, proposal, true); committed {
		return
	}

	// Line 22 in Algorithm 1 of The latest gossip on BFT consensus
	// check if to prevote this proposal in case proposal.vr == -1
	c.newProposalCheck(ctx, proposal)

	// Line 28 in Algorithm 1 of The latest gossip on BFT consensus
	// check if to prevote this proposal in case proposal.vr >= 0
	c.oldProposalCheck(ctx, proposal)

	// Line 36 in Algorithm 1 of The latest gossip on BFT consensus
	// check if we have quorum prevotes on the proposal
	c.quorumPrevotesCheck(ctx, proposal)
}

// upon condition rules to check when receiving a current round prevote
func (c *Core) currentPrevoteChecks(ctx context.Context) {
	// fetch current proposal
	curProposal := c.CurRoundMessages().Proposal()

	if curProposal != nil {
		// Line 36 in Algorithm 1 of The latest gossip on BFT consensus
		// check if we have quorum prevotes for the proposal
		c.quorumPrevotesCheck(ctx, curProposal)
	}

	// Line 44 in Algorithm 1 of The latest gossip on BFT consensus
	// check if we have quorum prevotes for nil, if so precommit nil
	c.quorumPrevotesNilCheck(ctx)

	// Line 34 in Algorithm 1 of The latest gossip on BFT consensus
	// check if we have to schedule prevote timeout
	c.prevoteTimeoutCheck()
}

// Rules to check at step change:
// 1. need to be checked only when we change from propose to prevote step
// 2. coincide with the ones to check when receiving a prevote
func (c *Core) stepChangeChecks(ctx context.Context) {
	if c.Step() != Prevote {
		panic("Step change tendermint checks done when transitioning to a step != Prevote")
	}
	c.currentPrevoteChecks(ctx)
}

// upon condition rules to check when receiving a current round precommit
func (c *Core) currentPrecommitChecks(ctx context.Context) {
	curProposal := c.CurRoundMessages().Proposal()

	// Line 49 in Algorithm 1 of The latest gossip on BFT consensus
	// check if we reached quorum precommits for the current proposal
	if curProposal != nil {
		// if we commit, no need to check the other rules
		if committed := c.quorumPrecommitsCheck(ctx, curProposal, true); committed {
			return
		}
	}

	// Line 47 in Algorithm 1 of The latest gossip on BFT consensus
	// check if we need to schedule the precommit timeout
	c.precommitTimeoutCheck()
}
