// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"bytes"
	"context"
	"errors"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
)

var (
	// errNotFromProposer is returned when received message is supposed to be from
	// proposer.
	errNotFromProposer = errors.New("message does not come from proposer")
	// errFutureHeightMessage is returned when roundState view is earlier than the
	// view of the received message.
	errFutureHeightMessage = errors.New("future height message")
	// errOldHeightMessage is returned when the received message's view is earlier
	// than roundState view.
	errOldHeightMessage = errors.New("old height message")
	// errInvalidMessage is returned when the message is malformed.
	errInvalidMessage = errors.New("invalid message")
	// errInvalidSenderOfCommittedSeal is returned when the committed seal is not from the sender of the message.
	errInvalidSenderOfCommittedSeal = errors.New("invalid sender of committed seal")
	// errFailedDecodeProposal is returned when the PROPOSAL message is malformed.
	errFailedDecodeProposal = errors.New("failed to decode PROPOSAL")
	// errFailedDecodePrevote is returned when the PREVOTE message is malformed.
	errFailedDecodePrevote = errors.New("failed to decode PREVOTE")
	// errFailedDecodePrecommit is returned when the PRECOMMIT message is malformed.
	errFailedDecodePrecommit = errors.New("failed to decode PRECOMMIT")
	// errFailedDecodeVote is returned for when PREVOTE or PRECOMMIT is malformed.
	errFailedDecodeVote = errors.New("failed to decode vote")
	// errNilPrevoteSent is returned when timer could be stopped in time
	errNilPrevoteSent = errors.New("timer expired and nil prevote sent")
	// errNilPrecommitSent is returned when timer could be stopped in time
	errNilPrecommitSent = errors.New("timer expired and nil precommit sent")
	// errMovedToNewRound is returned when timer could be stopped in time
	errMovedToNewRound = errors.New("timer expired and new round started")
)

// New creates an Tendermint consensus core
func New(backend Backend, config *config.Config) *core {
	logger := log.New("addr", backend.Address().String())
	return &core{
		config:                     config,
		address:                    backend.Address(),
		logger:                     logger,
		backend:                    backend,
		pendingUnminedBlocks:       make(map[uint64]*types.Block),
		pendingUnminedBlockCh:      make(chan *types.Block),
		stopped:                    make(chan struct{}, 3),
		isStarting:                 new(uint32),
		isStarted:                  new(uint32),
		isStopping:                 new(uint32),
		isStopped:                  new(uint32),
		valSet:                     new(validatorSet),
		lockedRound:                big.NewInt(-1),
		validRound:                 big.NewInt(-1),
		roundState:                 new(roundState),
		proposeTimeout:             newTimeout(propose, logger),
		prevoteTimeout:             newTimeout(prevote, logger),
		precommitTimeout:           newTimeout(precommit, logger),
		lastCommittedBlockProposer: common.Address{},
	}
}

type core struct {
	config  *config.Config
	address common.Address
	logger  log.Logger

	backend Backend
	cancel  context.CancelFunc

	messageEventSub         *event.TypeMuxSubscription
	newUnminedBlockEventSub *event.TypeMuxSubscription
	committedSub            *event.TypeMuxSubscription
	timeoutEventSub         *event.TypeMuxSubscription
	syncEventSub            *event.TypeMuxSubscription
	futureProposalTimer     *time.Timer
	stopped                 chan struct{}
	isStarted               *uint32
	isStarting              *uint32
	isStopping              *uint32
	isStopped               *uint32

	valSet *validatorSet

	roundState *roundState

	// map[Height]UnminedBlock
	pendingUnminedBlocks     map[uint64]*types.Block
	pendingUnminedBlocksMu   sync.Mutex
	pendingUnminedBlockCh    chan *types.Block
	isWaitingForUnminedBlock bool

	sentProposal          bool
	sentPrevote           bool
	sentPrecommit         bool
	setValidRoundAndValue bool

	lockedRound *big.Int
	validRound  *big.Int
	lockedValue *types.Block
	validValue  *types.Block

	proposeTimeout   *timeout
	prevoteTimeout   *timeout
	precommitTimeout *timeout

	lastCommittedBlockProposer common.Address
}

func (c *core) IsValidator(address common.Address) bool {
	_, val := c.valSet.GetByAddress(address)
	return val != nil
}

func (c *core) finalizeMessage(msg *Message) ([]byte, error) {
	var err error

	// Sign message
	data, err := msg.PayloadNoSig()
	if err != nil {
		return nil, err
	}
	msg.Signature, err = c.backend.Sign(data)
	if err != nil {
		return nil, err
	}

	// Convert to payload
	payload, err := msg.Payload()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (c *core) broadcast(ctx context.Context, msg *Message) {
	logger := c.logger.New("step", c.roundState.Step())

	payload, err := c.finalizeMessage(msg)
	if err != nil {
		logger.Error("Failed to finalize message", "msg", msg, "err", err)
		return
	}

	// Broadcast payload
	logger.Debug("broadcasting", "msg", msg.String())
	if err = c.backend.Broadcast(ctx, c.valSet.Copy(), payload); err != nil {
		logger.Error("Failed to broadcast message", "msg", msg, "err", err)
		return
	}
}

func (c *core) isProposer() bool {
	return c.valSet.IsProposer(c.address)
}

func (c *core) isProposerForR(r int64, a common.Address) bool {
	return c.valSet.IsProposerForRound(c.lastCommittedBlockProposer, uint64(r), a)
}

func (c *core) commit(round int64) {
	c.setStep(precommitDone)

	proposal := c.roundState.Proposal(round)
	if proposal == nil {
		// Should never happen really.
		c.logger.Error("core commit called with empty proposal ")
		return
	}

	if proposal.ProposalBlock == nil {
		// Again should never happen.
		c.logger.Error("commit a NIL block",
			"block", proposal.ProposalBlock,
			"height", c.roundState.height.String(),
			"round", c.roundState.round.String())
		return
	}

	c.logger.Info("commit a block", "hash", proposal.ProposalBlock.Header().Hash())

	precommits := c.roundState.allRoundMessages[round].precommits
	committedSeals := make([][]byte, precommits.VotesSize(proposal.ProposalBlock.Hash()))
	for i, v := range precommits.Values(proposal.ProposalBlock.Hash()) {
		committedSeals[i] = make([]byte, types.BFTExtraSeal)
		copy(committedSeals[i][:], v.CommittedSeal[:])
	}

	if err := c.backend.Commit(proposal.ProposalBlock, c.roundState.Round(), committedSeals); err != nil {
		c.logger.Error("failed to commit a block", "err", err)
		return
	}
}

// Metric collecton of round change and height change.
func (c *core) measureHeightRoundMetrics(round *big.Int) {
	if round.Cmp(common.Big0) == 0 {
		// in case of height change, round changed too, so count it also.
		tendermintRoundChangeMeter.Mark(1)
		tendermintHeightChangeMeter.Mark(1)
	} else {
		tendermintRoundChangeMeter.Mark(1)
	}
}

// startRound starts a new round. if round equals to 0, it means to starts a new height
func (c *core) startRound(ctx context.Context, round *big.Int) {

	c.measureHeightRoundMetrics(round)
	lastCommittedProposalBlock, lastCommittedBlockProposer := c.backend.LastCommittedProposal()
	height := new(big.Int).Add(lastCommittedProposalBlock.Number(), common.Big1)

	c.setCore(round, height, lastCommittedBlockProposer)

	// c.setStep(propose) will process the pending unmined blocks sent by the backed.Seal() and set c.lastestPendingRequest
	c.setStep(propose)

	c.logger.Debug("Starting new Round", "Height", height, "Round", round)

	// If the node is the proposer for this round then it would propose validValue or a new block, otherwise,
	// proposeTimeout is started, where the node waits for a proposal from the proposer of the current round.
	if c.isProposer() {
		// validValue and validRound represent a block they received a quorum of prevote and the round quorum was
		// received, respectively. If the block is not committed in that round then the round is changed.
		// The new proposer will chose the validValue, if present, which was set in one of the previous rounds otherwise
		// they propose a new block.
		var p *types.Block
		if c.validValue != nil {
			p = c.validValue
		} else {
			p = c.getUnminedBlock()
			if p == nil {
				select {
				case <-ctx.Done():
					return
				case p = <-c.pendingUnminedBlockCh:
				}
			}
		}
		c.sendProposal(ctx, p)
	} else {
		timeoutDuration := timeoutPropose(round.Int64())
		c.proposeTimeout.scheduleTimeout(timeoutDuration, round.Int64(), height.Int64(), c.onTimeoutPropose)
		c.logger.Debug("Scheduled Propose Timeout", "Timeout Duration", timeoutDuration)

		// Check if we already have the proposal (this will be true if a future proposal was received an a previous
		// round, if so send the proposal message to handler to handle the proposal
		if c.roundState.Proposal(round.Int64()) != nil {
			c.sendEvent(c.roundState.allRoundMessages[round.Int64()].proposalMsg)
		}
	}
}

func (c *core) setCore(r *big.Int, h *big.Int, lastProposer common.Address) {
	// Start of new height where round is 0
	if r.Int64() == 0 {
		// Set the shared round values to initial values
		c.lockedRound = big.NewInt(-1)
		c.lockedValue = nil
		c.validRound = big.NewInt(-1)
		c.validValue = nil

		c.lastCommittedBlockProposer = lastProposer

		// Set validator set for height
		valSet := c.backend.Validators(h.Uint64())
		c.valSet.set(valSet)
	}
	// Reset all timeouts
	c.proposeTimeout.reset(propose)
	c.prevoteTimeout.reset(prevote)
	c.precommitTimeout.reset(precommit)

	// update the round and height
	c.roundState.Update(r, h)

	// Calculate new proposer
	c.valSet.CalcProposer(lastProposer, r.Uint64())
	c.sentProposal = false
	c.sentPrevote = false
	c.sentPrecommit = false
	c.setValidRoundAndValue = false
}

func (c *core) setStep(step Step) {
	c.roundState.SetStep(step)
}

func (c *core) stopFutureProposalTimer() {
	if c.futureProposalTimer != nil {
		c.futureProposalTimer.Stop()
	}
}

func (c *core) Quorum(i int) bool {
	return float64(i) >= math.Ceil(float64(2)/float64(3)*float64(c.valSet.Size()))
}

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal(hash common.Hash, round *big.Int, height *big.Int) []byte {
	var buf bytes.Buffer
	buf.Write(round.Bytes())
	buf.Write(height.Bytes())
	buf.Write(hash.Bytes())
	return buf.Bytes()
}

func (c *core) isValid(proposal *types.Block) (bool, error) {
	if _, ok := c.roundState.verifiedProposals[proposal.Hash()]; !ok {
		if _, err := c.backend.VerifyProposal(*proposal); err != nil {
			return false, err
		}
		c.roundState.verifiedProposals[proposal.Hash()] = true
	}
	return true, nil
}

// Line 49 in Algorithm 1 of the latest gossip on BFT consensus
func (c *core) checkForConsensus(ctx context.Context, round int64) error {
	proposal := c.roundState.Proposal(round)
	proposalMsg := c.roundState.allRoundMessages[round].proposalMsg
	precommits := c.roundState.allRoundMessages[round].precommits
	h := proposal.ProposalBlock.Hash()

	if proposal != nil && c.isProposerForR(round, proposalMsg.Address) && c.Quorum(precommits.VotesSize(h)) {
		if valid, err := c.isValid(proposal.ProposalBlock); !valid {
			return err
		} else {

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				c.commit(round)
			}

		}
	}
	return nil
}

// Line 22 in Algorithm 1 of the latest gossip on BFT consensus
func (c *core) checkForNewProposal(ctx context.Context, round int64) error {
	proposal := c.roundState.Proposal(round)
	proposalMsg := c.roundState.allRoundMessages[round].proposalMsg
	h := proposal.ProposalBlock.Hash()

	if proposal != nil && c.isProposerForR(round, proposalMsg.Address) && c.roundState.Step() == propose {
		valid, err := c.isValid(proposal.ProposalBlock)

		// Vote for the proposal if proposal is valid(proposal) && (lockedRound = -1 || lockedValue = proposal).
		if valid && (c.lockedRound.Int64() == -1 || (c.lockedRound != nil && h == c.lockedValue.Hash())) {
			c.sendPrevote(ctx, true)
			c.setStep(prevote)
			return nil
		}
		c.sendPrevote(ctx, false)
		c.setStep(prevote)
		return err
	}
	return nil
}

// Line 28 in Algorithm 1 of the latest gossip on BFT consensus
func (c *core) checkForOldProposal(ctx context.Context, round int64) error {
	proposal := c.roundState.Proposal(round)
	proposalMsg := c.roundState.allRoundMessages[round].proposalMsg
	vr := proposal.ValidRound.Int64()
	validRoundPrevotes := c.roundState.allRoundMessages[vr].prevotes
	h := proposal.ProposalBlock.Hash()

	if proposal != nil && c.isProposerForR(round, proposalMsg.Address) && c.Quorum(validRoundPrevotes.VotesSize(h)) &&
		c.roundState.Step() == propose && vr >= 0 && vr < round {
		valid, err := c.isValid(proposal.ProposalBlock)

		// Vote for proposal if valid(proposal) && ((0 <= lockedRound <= vr < curR) || lockedValue == proposal)
		if valid && (c.lockedRound.Int64() <= vr || (c.lockedRound != nil && h == c.lockedValue.Hash())) {
			c.sendPrevote(ctx, true)
			c.setStep(prevote)
			return nil
		}
		c.sendPrevote(ctx, false)
		c.setStep(prevote)
		return err
	}
	return nil
}

// Line 36 in Algorithm 1 of the latest gossip on BFT consensus
func (c *core) checkForQuorumPrevotes(ctx context.Context, round int64) error {
	proposal := c.roundState.Proposal(round)
	proposalMsg := c.roundState.allRoundMessages[round].proposalMsg
	prevotes := c.roundState.allRoundMessages[round].prevotes
	h := proposal.ProposalBlock.Hash()

	if proposal != nil && c.isProposerForR(round, proposalMsg.Address) && c.Quorum(prevotes.VotesSize(h)) &&
		c.roundState.Step() >= prevote && !c.setValidRoundAndValue {
		if valid, err := c.isValid(proposal.ProposalBlock); !valid {
			return err
		}

		if c.roundState.Step() == prevote {
			c.lockedValue = proposal.ProposalBlock
			c.lockedRound = big.NewInt(round)
			c.sendPrecommit(ctx, false)
			c.setStep(precommit)
		}
		c.validValue = proposal.ProposalBlock
		c.validRound = big.NewInt(round)
		c.setValidRoundAndValue = true

	}
	return nil
}
