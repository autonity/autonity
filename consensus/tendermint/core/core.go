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
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/contracts/autonity"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"github.com/davecgh/go-spew/spew"
)

var (
	// errNotFromProposer is returned when received message is supposed to be from
	// proposer.
	errNotFromProposer = errors.New("message does not come from proposer")
	// errFutureHeightMessage is returned when curRoundMessages view is earlier than the
	// view of the received message.
	errFutureHeightMessage = errors.New("future height message")
	// errOldHeightMessage is returned when the received message's view is earlier
	// than curRoundMessages view.
	errOldHeightMessage = errors.New("old height message")
	// errOldRoundMessage message is returned when message is of the same Height but form a smaller round
	errOldRoundMessage = errors.New("same height but old round message")
	// errFutureRoundMessage message is returned when message is of the same Height but form a newer round
	errFutureRoundMessage = errors.New("same height but future round message")
	// errFutureStepMessage message is returned when it's a prevote or precommit message of the same Height same round
	// while the current step is propose.
	errFutureStepMessage = errors.New("same round but future step message")
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

const (
	MaxRound = 99 // consequence of backlog priority
)

// New creates an Tendermint consensus core
func New(backend Backend, config *config.Config) *core {
	addr := backend.Address()
	logger := log.New("addr", addr.String())
	return &core{
		proposerPolicy:        config.ProposerPolicy,
		blockPeriod:           config.BlockPeriod,
		address:               addr,
		logger:                logger,
		backend:               backend,
		backlogs:              make(map[types.CommitteeMember][]*Message),
		pendingUnminedBlocks:  make(map[uint64]*types.Block),
		pendingUnminedBlockCh: make(chan *types.Block),
		stopped:               make(chan struct{}, 4),
		committee:             nil,
		futureRoundChange:     make(map[int64]map[common.Address]uint64),
		lockedRound:           -1,
		validRound:            -1,
		proposeTimeout:        newTimeout(propose, logger),
		prevoteTimeout:        newTimeout(prevote, logger),
		precommitTimeout:      newTimeout(precommit, logger),
	}
}

type core struct {
	proposerPolicy config.ProposerPolicy
	blockPeriod    uint64
	address        common.Address
	logger         log.Logger

	backend Backend
	cancel  context.CancelFunc

	messageEventSub         *event.TypeMuxSubscription
	newUnminedBlockEventSub *event.TypeMuxSubscription
	committedSub            *event.TypeMuxSubscription
	timeoutEventSub         *event.TypeMuxSubscription
	consensusMessageSub     *event.TypeMuxSubscription
	syncEventSub            *event.TypeMuxSubscription
	futureProposalTimer     *time.Timer
	stopped                 chan struct{}

	msgCache   *messageCache
	backlogs   map[types.CommitteeMember][]*Message
	backlogsMu sync.Mutex
	// map[Height]UnminedBlock
	pendingUnminedBlocks     map[uint64]*types.Block
	pendingUnminedBlocksMu   sync.Mutex
	pendingUnminedBlockCh    chan *types.Block
	isWaitingForUnminedBlock bool

	//
	// Tendermint FSM state fields
	//

	height     *big.Int
	round      int64
	committee  committee
	lastHeader *types.Header
	// height, round and committeeSet are the ONLY guarded fields.
	// everything else MUST be accessed only by the main thread.
	stateMu               sync.RWMutex
	step                  Step
	sentProposal          bool
	sentPrevote           bool
	sentPrecommit         bool
	setValidRoundAndValue bool

	lockedRound int64
	validRound  int64
	lockedValue *types.Block
	validValue  *types.Block

	proposeTimeout   *timeout
	prevoteTimeout   *timeout
	precommitTimeout *timeout

	futureRoundChange map[int64]map[common.Address]uint64

	autonityContract *autonity.Contract

	line34Executed bool
	line36Executed bool
	line47Executed bool
}

func (c *core) GetCurrentHeightMessages() []*Message {
	return c.msgCache.heightMessages(c.Height().Uint64())
}

func (c *core) IsMember(address common.Address) bool {
	_, _, err := c.committeeSet().GetByAddress(address)
	return err == nil
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
	logger := c.logger.New("step", c.step)

	payload, err := c.finalizeMessage(msg)
	if err != nil {
		logger.Error("Failed to finalize message", "msg", msg, "err", err)
		return
	}

	// Broadcast payload
	logger.Debug("broadcasting", "msg", msg.String())
	if err = c.backend.Broadcast(ctx, c.committeeSet().Committee(), payload); err != nil {
		logger.Error("Failed to broadcast message", "msg", msg, "err", err)
		return
	}
}

// check if msg sender is proposer for proposal handling.
func (c *core) isProposerMsg(round int64, msgAddress common.Address) bool {
	return c.committeeSet().GetProposer(round).Address == msgAddress
}
func (c *core) isProposer() bool {
	return c.committeeSet().GetProposer(c.Round()).Address == c.address
}

func (c *core) commit(block *types.Block, round int64) {
	c.setStep(precommitDone)

	// Sanity check
	if block == nil {
		panic(fmt.Sprintf("Attempted to commit nil block: %s", spew.Sdump(block)))
	}

	c.logger.Info("commit a block", "hash", block.Hash())

	committedSeals := c.msgCache.signatures(block.Hash(), round, block.NumberU64())
	if err := c.backend.Commit(block, round, committedSeals); err != nil {
		c.logger.Error("failed to commit a block", "err", err)
	}
}

// Metric collecton of round change and height change.
func (c *core) measureHeightRoundMetrics(round int64) {
	if round == 0 {
		tendermintHeightChangeMeter.Mark(1)
	}
	tendermintRoundChangeMeter.Mark(1)
}

// startRound starts a new round. if round equals to 0, it means to starts a new height
func (c *core) startRound(ctx context.Context, round int64) {

	c.measureHeightRoundMetrics(round)
	// Set initial FSM state
	c.setInitialState(round)
	// c.setStep(propose) will process the pending unmined blocks sent by the backed.Seal() and set c.lastestPendingRequest
	c.setStep(propose)
	c.logger.Debug("Starting new Round", "Height", c.Height(), "Round", round)

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
		timeoutDuration := c.timeoutPropose(round)
		c.proposeTimeout.scheduleTimeout(timeoutDuration, round, c.Height(), c.onTimeoutPropose)
		c.logger.Debug("Scheduled Propose Timeout", "Timeout Duration", timeoutDuration)
	}

}

func (c *core) reprocessConsensusMessage(cm *consensusMessage) error {
	go c.sendEvent(cm) // Could we use less go routines?
	return nil
}

func (c *core) setInitialState(r int64) {
	// Start of new height where round is 0
	if r == 0 {
		lastBlockMined, _ := c.backend.LastCommittedProposal()
		c.setHeight(new(big.Int).Add(lastBlockMined.Number(), common.Big1))

		lastHeader := lastBlockMined.Header()
		var committeeSet committee
		var err error
		var lastProposer common.Address
		switch c.proposerPolicy {
		case config.RoundRobin:
			if !lastHeader.IsGenesis() {
				var err error
				lastProposer, err = types.Ecrecover(lastHeader)
				if err != nil {
					panic(fmt.Sprintf("unable to recover proposer address from header %q: %v", lastHeader, err))
				}
			}
			committeeSet, err = newRoundRobinSet(lastHeader.Committee, lastProposer)
			if err != nil {
				panic(fmt.Sprintf("failed to construct committee %v", err))
			}
		case config.WeightedRandomSampling:
			committeeSet = newWeightedRandomSamplingCommittee(lastBlockMined, c.autonityContract, c.backend.BlockChain())
		default:
			panic(fmt.Sprintf("unrecognised proposer policy %q", c.proposerPolicy))
		}

		c.lastHeader = lastHeader
		c.setCommitteeSet(committeeSet)
		c.lockedRound = -1
		c.lockedValue = nil
		c.validRound = -1
		c.validValue = nil
		c.futureRoundChange = make(map[int64]map[common.Address]uint64)
	}

	c.proposeTimeout.reset(propose)
	c.prevoteTimeout.reset(prevote)
	c.precommitTimeout.reset(precommit)
	c.sentProposal = false
	c.sentPrevote = false
	c.sentPrecommit = false
	c.setValidRoundAndValue = false
	c.setRound(r)
}

func (c *core) setStep(step Step) {
	// Hmm thinking about reprocessing now.
	//
	// If we have changed step to propose then the only 2 upon checks which
	// have now, become acceissible are both rely on a proposal message, so by
	// reprocessing the proposal message we can ensure that we hit all relevant checks.
	//
	//
	// Consider what checks needs to be reconsidered.
	//
	// We assume we process all messages for a height when we reach that
	// height, we don't consider this reprocessing because those messages have
	// never been fully processed before.
	//
	// First off checks that never need reprocessing are
	//
	// 49 Its not specific to any round or step, and so would be triggered as
	// soon either the right proposals or precommits are received.
	//
	// 55 Again its not specific to any round or step, and so would be
	// triggered as soon as the right threshold is met.
	//
	// ---
	//
	// For a round change, aka step change from any step to propose. We need to
	// check every check that is specific to round. I.E has a roundp value
	// somewhere in the check.
	//
	// Those checks are. (all checks not including the 49 and 55 which we know
	// do not need to be rechecked)
	//
	// Line 22
	// Line 28
	// Line 34
	// Line 36
	// Line 44
	// Line 47
	//
	// What messages shall we re-process to hit all those checks? Re-processing
	// just the proposals will cover 22, 28, 36. That leaves 34, 44 and 47.
	//
	// 34 can be hit by sending a prevote with any value, so we could just pick
	// the first prevote for that round and process it.
	//
	// 44 can be hit by sending a prevote for nil, so we would have to find a
	// prevote for nil in the message cache and send that. We can't make one up
	// because just one node may constitute quorum threshold. Actually I think
	// we can, because we would not see any stake attached to the message if
	// noone voted.
	//
	// 47 can be hit by a precommit with any value
	//
	// So we can hit both 34 and 44 with a prevote for nil and 47 with a
	// precommit for any value. So we can hit all checks by reprocessing just
	// the proposals and a fake precommit and fake prevote.
	//
	// ---
	//
	// For step change propose to prevote, we will need to check all checks
	// that were not accessible with step propose and are now accessible with
	// step prevote.
	//
	// Those checks are.
	//
	// Line 34
	// Line 36
	// Line 44
	//
	// What messages shall we reprocess to hit all those checks. Reprocessing
	// the proposals will cover 36. That leaves 34 and 44. As we can see from
	// our reasoning above a fake prevote for nil will hit those checks.
	//
	//
	// For step change prevote to precommit, we will need to check all checks
	// that were not accessible with step prevote or propose and are now
	// accessible with step precommit. There are no checks that require a step
	// of precommit, any relevant checks would have already triggered as soon
	// as messages were recived in the propose or prevote steps, so nothing
	// needs to be re-processed.
	c.logger.Debug("moving to step", "step", step.String(), "round", c.Round())
	c.step = step
	switch step {
	case propose:
		c.msgCache.roundMessages(c.height.Uint64(), c.round, c.reprocessConsensusMessage)
	case prevote:
	}
	c.processBacklog()
}

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal(hash common.Hash, round int64, height *big.Int) []byte {
	var buf bytes.Buffer
	roundBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(roundBytes, uint64(round))
	buf.Write(roundBytes)
	buf.Write(height.Bytes())
	buf.Write(hash.Bytes())
	return buf.Bytes()
}

func (c *core) setRound(round int64) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.round = round
}

func (c *core) setHeight(height *big.Int) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.height = height
}
func (c *core) setCommitteeSet(set committee) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.committee = set
}

func (c *core) Round() int64 {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.round
}

func (c *core) Height() *big.Int {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.height
}
func (c *core) committeeSet() committee {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.committee
}
