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
	"math/big"
	"sync"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
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
	logger := log.New("addr", backend.Address().String())
	messagesMap := newMessagesMap()
	roundMessage := messagesMap.getOrCreate(0)
	return &core{
		config:                config,
		address:               backend.Address(),
		logger:                logger,
		backend:               backend,
		backlogs:              make(map[types.CommitteeMember]*prque.Prque),
		pendingUnminedBlocks:  make(map[uint64]*types.Block),
		pendingUnminedBlockCh: make(chan *types.Block),
		stopped:               make(chan struct{}, 3),
		isStarting:            new(uint32),
		isStarted:             new(uint32),
		isStopping:            new(uint32),
		isStopped:             new(uint32),
		committeeSet:          nil,
		futureRoundChange:     make(map[int64]map[common.Address]struct{}),
		messages:              messagesMap,
		lockedRound:           -1,
		validRound:            -1,
		curRoundMessages:      roundMessage,
		proposeTimeout:        newTimeout(propose, logger),
		prevoteTimeout:        newTimeout(prevote, logger),
		precommitTimeout:      newTimeout(precommit, logger),
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

	backlogs   map[types.CommitteeMember]*prque.Prque
	backlogsMu sync.Mutex
	// map[Height]UnminedBlock
	pendingUnminedBlocks     map[uint64]*types.Block
	pendingUnminedBlocksMu   sync.Mutex
	pendingUnminedBlockCh    chan *types.Block
	isWaitingForUnminedBlock bool

	//
	// Tendermint FSM state fields
	//

	height       *big.Int
	round        int64
	committeeSet committee.Set
	// height, round and committeeSet are the ONLY guarded fields.
	// everything else MUST be accessed only by the main thread.
	stateMu               sync.RWMutex
	step                  Step
	curRoundMessages      *roundMessages
	messages              messagesMap
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

	futureRoundChange map[int64]map[common.Address]struct{}
}

func (c *core) GetCurrentHeightMessages() []*Message {
	return c.messages.GetMessages()
}

func (c *core) IsMember(address common.Address) bool {
	_, _, err := c.CommitteeSet().GetByAddress(address)
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
	if err = c.backend.Broadcast(ctx, c.CommitteeSet(), payload); err != nil {
		logger.Error("Failed to broadcast message", "msg", msg, "err", err)
		return
	}
}

func (c *core) isProposer() bool {
	return c.CommitteeSet().IsProposer(c.Round(), c.address)
}

func (c *core) commit(round int64, messages *roundMessages) {
	c.setStep(precommitDone)

	proposal := messages.Proposal()
	if proposal == nil {
		// Should never happen really.
		c.logger.Error("core commit called with empty proposal ")
		return
	}

	if proposal.ProposalBlock == nil {
		// Again should never happen.
		c.logger.Error("commit a NIL block",
			"block", proposal.ProposalBlock,
			"height", c.Height(),
			"round", round)
		return
	}

	c.logger.Info("commit a block", "hash", proposal.ProposalBlock.Header().Hash())

	committedSeals := make([][]byte, messages.PrecommitsCount(proposal.ProposalBlock.Hash()))
	for i, v := range messages.CommitedSeals(proposal.ProposalBlock.Hash()) {
		committedSeals[i] = make([]byte, types.BFTExtraSeal)
		copy(committedSeals[i][:], v.CommittedSeal[:])
	}

	if err := c.backend.Commit(proposal.ProposalBlock, round, committedSeals); err != nil {
		c.logger.Error("failed to commit a block", "err", err)
		return
	}
}

// Metric collecton of round change and height change.
func (c *core) measureHeightRoundMetrics(round int64) {
	if round == 0 {
		// in case of height change, round changed too, so count it also.
		tendermintRoundChangeMeter.Mark(1)
		tendermintHeightChangeMeter.Mark(1)
	} else {
		tendermintRoundChangeMeter.Mark(1)
	}
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
		timeoutDuration := timeoutPropose(round)
		c.proposeTimeout.scheduleTimeout(timeoutDuration, round, c.Height(), c.onTimeoutPropose)
		c.logger.Debug("Scheduled Propose Timeout", "Timeout Duration", timeoutDuration)
	}
}

func (c *core) setInitialState(r int64) {
	// Start of new height where round is 0
	if r == 0 {
		lastBlockMined, _ := c.backend.LastCommittedProposal()
		c.setHeight(new(big.Int).Add(lastBlockMined.Number(), common.Big1))
		committeeSet, err := c.backend.Committee(c.Height().Uint64())
		if err != nil {
			c.logger.Error("fatal error: could not retrieve saved committee")
			panic(err)
		}
		c.setCommitteeSet(committeeSet)
		c.lockedRound = -1
		c.lockedValue = nil
		c.validRound = -1
		c.validValue = nil
		c.messages.reset()
		c.futureRoundChange = make(map[int64]map[common.Address]struct{})
	}

	c.proposeTimeout.reset(propose)
	c.prevoteTimeout.reset(prevote)
	c.precommitTimeout.reset(precommit)
	c.curRoundMessages = c.messages.getOrCreate(r)
	c.sentProposal = false
	c.sentPrevote = false
	c.sentPrecommit = false
	c.setValidRoundAndValue = false
	c.setRound(r)
}

func (c *core) acceptVote(roundMsgs *roundMessages, step Step, hash common.Hash, msg Message) {
	switch step {
	case prevote:
		roundMsgs.AddPrevote(hash, msg)
	case precommit:
		roundMsgs.AddPrecommit(hash, msg)
	}
}

func (c *core) setStep(step Step) {
	c.logger.Debug("moving to step", "step", step.String(), "round", c.Round())
	c.step = step
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
func (c *core) setCommitteeSet(set committee.Set) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.committeeSet = set
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
func (c *core) CommitteeSet() committee.Set {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.committeeSet
}
