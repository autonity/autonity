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
	"errors"
	"math/big"
	"reflect"
	"sync"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
)

const (
	initialProposeTimeout = 5 * time.Second
	//initialPrevoteTimeout   = 5 * time.Second
	//initialPrecommitTimeout = 5 * time.Second
)

var (
	// errInconsistentSubject is returned when received subject is different from
	// currentRoundState subject.
	errInconsistentSubject = errors.New("inconsistent subjects")
	// errNotFromProposer is returned when received message is supposed to be from
	// proposer.
	errNotFromProposer = errors.New("message does not come from proposer")
	// errIgnored is returned when a message was ignored.
	//errIgnored = errors.New("message is ignored")
	// errFutureMessage is returned when currentRoundState view is earlier than the
	// view of the received message.
	errFutureMessage = errors.New("future message")
	// errOldMessage is returned when the received message's view is earlier
	// than currentRoundState view.
	errOldMessage = errors.New("old message")
	// errInvalidMessage is returned when the message is malformed.
	errInvalidMessage = errors.New("invalid message")
	// errFailedDecodeProposal is returned when the PRE-PREPARE message is malformed.
	errFailedDecodeProposal = errors.New("failed to decode PRE-PREPARE")
	// errFailedDecodePrevote is returned when the PREPARE message is malformed.
	errFailedDecodePrevote = errors.New("failed to decode PREPARE")
	// errFailedDecodePrecommit is returned when the COMMIT message is malformed.
	errFailedDecodePrecommit = errors.New("failed to decode COMMIT")
)

var (
	// msgPriority is defined for calculating processing priority to speedup consensus
	// msgProposal > msgPrecommit > msgPrevote
	msgPriority = map[uint64]int{
		msgProposal:  1,
		msgPrecommit: 2,
		msgPrevote:   3,
	}
)

// New creates an Istanbul consensus core
func New(backend tendermint.Backend, config *tendermint.Config) Engine {
	c := &core{
		config:                 config,
		address:                backend.Address(),
		step:                   StepAcceptProposal,
		handlerStopCh:          make(chan struct{}),
		logger:                 log.New("address", backend.Address()),
		backend:                backend,
		backlogs:               make(map[tendermint.Validator]*prque.Prque),
		backlogsMu:             new(sync.Mutex),
		pendingUnminedBlocks:   prque.New(),
		pendingUnminedBlocksMu: new(sync.Mutex),
		proposeTimeout:         new(timeout),
		prevoteTimeout:         new(timeout),
		precommitTimeout:       new(timeout),
	}
	c.validateFn = c.checkValidatorSignature
	return c
}

// ----------------------------------------------------------------------------

type core struct {
	config  *tendermint.Config
	address common.Address
	step    Step
	logger  log.Logger

	backend             tendermint.Backend
	events              *event.TypeMuxSubscription
	finalCommittedSub   *event.TypeMuxSubscription
	timeoutSub          *event.TypeMuxSubscription
	futureProposalTimer *time.Timer

	valSet     tendermint.ValidatorSet
	validateFn func([]byte, []byte) (common.Address, error)

	backlogs   map[tendermint.Validator]*prque.Prque
	backlogsMu *sync.Mutex

	currentRoundState *roundState
	handlerStopCh     chan struct{}

	pendingUnminedBlocks   *prque.Prque
	pendingUnminedBlocksMu *sync.Mutex

	sentProposal bool

	lockedRound *big.Int
	validRound  *big.Int
	lockedValue *types.Block
	validValue  *types.Block

	currentHeightRoundsStates []roundState

	// TODO: may require a mutex
	latestPendingUnminedBlock *types.Block

	proposeTimeout   *timeout
	prevoteTimeout   *timeout
	precommitTimeout *timeout
}

func (c *core) finalizeMessage(msg *message) ([]byte, error) {
	var err error
	// Add sender address
	msg.Address = c.Address()

	// Add proof of consensus
	msg.CommittedSeal = []byte{}
	// Assign the CommittedSeal if it's a COMMIT message and proposal is not nil
	if msg.Code == msgPrecommit && c.currentRoundState.Proposal() != nil {
		seal := PrepareCommittedSeal(c.currentRoundState.Proposal().ProposalBlock.Hash())
		msg.CommittedSeal, err = c.backend.Sign(seal)
		if err != nil {
			return nil, err
		}
	}

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

func (c *core) broadcast(msg *message) {
	logger := c.logger.New("step", c.step)

	payload, err := c.finalizeMessage(msg)
	if err != nil {
		logger.Error("Failed to finalize message", "msg", msg, "err", err)
		return
	}

	// Broadcast payload
	if err = c.backend.Broadcast(c.valSet, payload); err != nil {
		logger.Error("Failed to broadcast message", "msg", msg, "err", err)
		return
	}
}

func (c *core) isProposer() bool {
	v := c.valSet
	if v == nil {
		return false
	}
	return v.IsProposer(c.backend.Address())
}

func (c *core) commit() {
	c.setStep(StepPrecommitDone)

	proposal := c.currentRoundState.Proposal()
	if proposal != nil {
		committedSeals := make([][]byte, c.currentRoundState.Precommits.Size())
		for i, v := range c.currentRoundState.Precommits.Values() {
			committedSeals[i] = make([]byte, types.PoSExtraSeal)
			copy(committedSeals[i][:], v.CommittedSeal[:])
		}

		if err := c.backend.Precommit(proposal.ProposalBlock, committedSeals); err != nil {
			c.currentRoundState.UnlockHash() //Unlock block when insertion fails
			// TODO: go to next height
			return
		}
	}
}

// startRound starts a new round. if round equals to 0, it means to starts a new height
func (c *core) startRound(round *big.Int) {
	//TODO: update the name of lastProposalBlock and LastBlockProposal()
	lastProposalBlock, lastProposalBlockProposer := c.backend.LastProposal()
	height := new(big.Int).Add(lastProposalBlock.Number(), common.Big1)

	// Start of new height where round is 0
	if round.Uint64() == 0 {
		// Set the shared round values to initial values
		c.lockedRound = big.NewInt(-1)
		c.lockedValue = new(types.Block)
		c.validRound = big.NewInt(-1)
		c.validValue = new(types.Block)

		c.valSet = c.backend.Validators(height.Uint64())

		// TODO: Assuming that round == 0 only when the node moves to a new height, need to confirm where exactly the node moves to a new height
		c.currentHeightRoundsStates = nil

	} else {
		// Assuming the above values will be set for round > 0
		// Add the currentRoundState round step to the core previous round states
		c.currentHeightRoundsStates = append(c.currentHeightRoundsStates, *c.currentRoundState)
	}

	c.currentRoundState = newRoundState(round, height, c.valSet, common.Hash{}, nil, c.backend.HasBadProposal)
	c.valSet.CalcProposer(lastProposalBlockProposer, round.Uint64())
	c.sentProposal = false
	// c.setStep(StepAcceptProposal) will process the pending unmined blocks sent by the backed.Seal() and set c.lastestPendingRequest
	c.setStep(StepAcceptProposal)

	var p *types.Block
	if c.isProposer() {
		if c.validValue != nil {
			p = c.validValue
		} else {
			p = c.latestPendingUnminedBlock
		}
		c.sendProposal(p)
	} else {
		c.proposeTimeout.scheduleTimeout(timeoutPropose(round.Int64()), c.onTimeoutPropose)
	}
}

func (c *core) setStep(step Step) {
	// TODO: remove the if
	if c.step != step {
		c.step = step
	}
	if step == StepAcceptProposal {
		c.processPendingRequests()
	}
	c.processBacklog()
}

func (c *core) Address() common.Address {
	return c.address
}

func (c *core) stopFutureProposalTimer() {
	if c.futureProposalTimer != nil {
		c.futureProposalTimer.Stop()
	}
}

func (c *core) stopTimer() {
	c.stopFutureProposalTimer()
}

func (c *core) checkValidatorSignature(data []byte, sig []byte) (common.Address, error) {
	return tendermint.CheckValidatorSignature(c.valSet, data, sig)
}

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal(hash common.Hash) []byte {
	var buf bytes.Buffer
	buf.Write(hash.Bytes())
	buf.Write([]byte{byte(msgPrecommit)})
	return buf.Bytes()
}

func (c *core) onTimeoutPropose() {
}

//func (c *core) onTimeoutPrevote() {
//}
//
//func (c *core) onTimeoutPrecommit() {
//
//}

//---------------------------------------Timeout---------------------------------------

type timeoutEvent struct{}

type timeout struct {
	timer *time.Timer
	sync.RWMutex
}

// runAfterTimeout() will be run in a separate go routine, so values used inside the function needs to be managed separately
func (t *timeout) scheduleTimeout(stepTimeout time.Duration, runAfterTimeout func()) *time.Timer {
	t.Lock()
	defer t.Unlock()
	t.timer = time.AfterFunc(stepTimeout, runAfterTimeout)
	return t.timer
}

//func (t *timeout) stopTimer() bool {
//	t.RLock()
//	defer t.RUnlock()
//	return t.timer.Stop()
//}

// The timeout may need to be changed depending on the Step
func timeoutPropose(round int64) time.Duration {
	return initialProposeTimeout + time.Duration(round)*time.Second
}

//func timeoutPrevote(round int64) time.Duration {
//	return initialProposeTimeout + time.Duration(round)*time.Second
//}
//
//func timeoutPrecommit(round int64) time.Duration {
//	return initialProposeTimeout + time.Duration(round)*time.Second
//}

//---------------------------------------Handler---------------------------------------

// Start implements core.Engine.Start
func (c *core) Start() error {
	// Start a new round from last height + 1
	c.startRound(common.Big0)

	// Tests will handle events itself, so we have to make subscribeEvents()
	// be able to call in test.
	c.subscribeEvents()
	go c.handleEvents()

	return nil
}

// Stop implements core.Engine.Stop
func (c *core) Stop() error {
	c.stopTimer()
	c.unsubscribeEvents()

	// Make sure the handler goroutine exits
	c.handlerStopCh <- struct{}{}
	return nil
}

// TODO: update all of the TypeMuxSilent to event.Feed

// Subscribe both internal and external events
func (c *core) subscribeEvents() {
	c.events = c.backend.EventMux().Subscribe(
		// external events
		tendermint.NewUnminedBlockEvent{},
		tendermint.MessageEvent{},
		// internal events
		backlogEvent{},
	)
	// TODO: not sure why a backend EventMux is being used for core internal events, lazy coding
	c.timeoutSub = c.backend.EventMux().Subscribe(
		timeoutEvent{},
	)
	// TODO: not sure why a backend EventMux is being used for core internal events, lazy coding
	c.finalCommittedSub = c.backend.EventMux().Subscribe(
		tendermint.CommitEvent{},
	)
}

// Unsubscribe all events
func (c *core) unsubscribeEvents() {
	c.events.Unsubscribe()
	c.timeoutSub.Unsubscribe()
	c.finalCommittedSub.Unsubscribe()
}

func (c *core) handleEvents() {
	// Clear step
	defer func() {
		c.currentRoundState = nil
		<-c.handlerStopCh
	}()

	for {
		select {
		case ev, ok := <-c.events.Chan():
			if !ok {
				return
			}
			// A real ev arrived, process interesting content
			switch e := ev.Data.(type) {
			case tendermint.NewUnminedBlockEvent:
				pb := &e.NewUnminedBlock
				err := c.handleUnminedBlock(pb)
				if err == errFutureMessage {
					c.storeUnminedBlockMsg(pb)
				}
			case tendermint.MessageEvent:
				if err := c.handleMsg(e.Payload); err == nil {
					c.backend.Gossip(c.valSet, e.Payload)
				}
			case backlogEvent:
				// No need to check signature for internal messages
				if err := c.handleCheckedMsg(e.msg, e.src); err == nil {
					p, err := e.msg.Payload()
					if err != nil {
						c.logger.Warn("Get message payload failed", "err", err)
						continue
					}
					c.backend.Gossip(c.valSet, p)
				}
			}
		case _, ok := <-c.timeoutSub.Chan():
			if !ok {
				return
			}
			c.handleTimeoutMsg()
		case ev, ok := <-c.finalCommittedSub.Chan():
			if !ok {
				return
			}
			switch ev.Data.(type) {
			case tendermint.CommitEvent:
				c.handleCommit()
			}
		}
	}
}

// sendEvent sends events to mux
func (c *core) sendEvent(ev interface{}) {
	c.backend.EventMux().Post(ev)
}

func (c *core) handleMsg(payload []byte) error {
	logger := c.logger.New()

	// Decode message and check its signature
	msg := new(message)
	if err := msg.FromPayload(payload, c.validateFn); err != nil {
		logger.Error("Failed to decode message from payload", "err", err)
		return err
	}

	// Only accept message if the address is valid
	// TODO: the check is already made in c.validateFn
	_, src := c.valSet.GetByAddress(msg.Address)
	if src == nil {
		logger.Error("Invalid address in message", "msg", msg)
		return tendermint.ErrUnauthorizedAddress
	}

	return c.handleCheckedMsg(msg, src)
}

func (c *core) handleCheckedMsg(msg *message, src tendermint.Validator) error {
	logger := c.logger.New("address", c.address, "from", src)

	// Store the message if it's a future message
	testBacklog := func(err error) error {
		if err == errFutureMessage {
			c.storeBacklog(msg, src)
		}

		return err
	}

	switch msg.Code {
	case msgProposal:
		return testBacklog(c.handleProposal(msg, src))
	case msgPrevote:
		return testBacklog(c.handlePrevote(msg, src))
	case msgPrecommit:
		return testBacklog(c.handlePrecommit(msg, src))
	default:
		logger.Error("Invalid message", "msg", msg)
	}

	return errInvalidMessage
}

// TODO: re-implement to incorporate all three timeouts
func (c *core) handleTimeoutMsg() {
}

//---------------------------------------Backlog---------------------------------------

type backlogEvent struct {
	src tendermint.Validator
	msg *message
}

// checkMessage checks the message step
// return errInvalidMessage if the message is invalid
// return errFutureMessage if the message view is larger than currentRoundState view
// return errOldMessage if the message view is smaller than currentRoundState view
func (c *core) checkMessage(msgCode uint64, round *big.Int, height *big.Int) error {
	if height == nil || round == nil {
		return errInvalidMessage
	}

	// TODO: add current round messages to currentHeightRoundStates
	if height.Cmp(c.currentRoundState.Height()) > 0 {
		return errFutureMessage
	} else if round.Cmp(c.currentRoundState.Round()) > 0 {
		return errFutureMessage
	} else if height.Cmp(c.currentRoundState.Height()) < 0 {
		return errOldMessage
	} else if round.Cmp(c.currentRoundState.Round()) < 0 {
		return errOldMessage
	}

	// StepAcceptProposal only accepts msgProposal
	// other messages are future messages
	if c.step == StepAcceptProposal {
		if msgCode > msgProposal {
			return errFutureMessage
		}
		return nil
	}

	// For steps(StepProposeDone, StepPrevoteDone, StepPrecommitDone),
	// can accept all message types if processing with same view
	return nil
}

func (c *core) storeBacklog(msg *message, src tendermint.Validator) {
	logger := c.logger.New("from", src, "step", c.step)

	if src.Address() == c.Address() {
		logger.Warn("Backlog from self")
		return
	}

	logger.Trace("Store future message")

	c.backlogsMu.Lock()
	defer c.backlogsMu.Unlock()

	backlog := c.backlogs[src]
	if backlog == nil {
		backlog = prque.New()
	}
	switch msg.Code {
	case msgProposal:
		var p *tendermint.Proposal
		err := msg.Decode(&p)
		if err == nil {
			backlog.Push(msg, toPriority(msg.Code, p.Round, p.Height))
		}
		// for msgRoundChange, msgPrevote and msgPrecommit cases
	default:
		var p *tendermint.Subject
		err := msg.Decode(&p)
		if err == nil {
			backlog.Push(msg, toPriority(msg.Code, p.Round, p.Height))
		}
	}
	c.backlogs[src] = backlog
}

func (c *core) processBacklog() {
	c.backlogsMu.Lock()
	defer c.backlogsMu.Unlock()

	for src, backlog := range c.backlogs {
		if backlog == nil {
			continue
		}

		logger := c.logger.New("from", src, "step", c.step)
		var isFuture bool

		// We stop processing if
		//   1. backlog is empty
		//   2. The first message in queue is a future message
		for !(backlog.Empty() || isFuture) {
			m, prio := backlog.Pop()
			msg := m.(*message)
			var round, height *big.Int
			switch msg.Code {
			case msgProposal:
				var m *tendermint.Proposal
				err := msg.Decode(&m)
				if err == nil {
					round, height = m.Round, m.Height
				}
				// for msgRoundChange, msgPrevote and msgPrecommit cases
			default:
				var sub *tendermint.Subject
				err := msg.Decode(&sub)
				if err == nil {
					round, height = sub.Round, sub.Height
				}
			}
			if round == nil || height == nil {
				logger.Debug("Nil view", "msg", msg)
				continue
			}
			// Push back if it's a future message
			err := c.checkMessage(msg.Code, round, height)
			if err != nil {
				if err == errFutureMessage {
					logger.Trace("Stop processing backlog", "msg", msg)
					backlog.Push(msg, prio)
					isFuture = true
					break
				}
				logger.Trace("Skip the backlog event", "msg", msg, "err", err)
				continue
			}
			logger.Trace("Post backlog event", "msg", msg)

			go c.sendEvent(backlogEvent{
				src: src,
				msg: msg,
			})
		}
	}
}

func toPriority(msgCode uint64, r *big.Int, h *big.Int) float32 {
	// FIXME: round will be reset as 0 while new height
	// 10 * Round limits the range of message code is from 0 to 9
	// 1000 * Height limits the range of round is from 0 to 99
	return -float32(h.Uint64()*1000 + r.Uint64()*10 + uint64(msgPriority[msgCode]))
}

//---------------------------------------NewUnminedBlock---------------------------------------

func (c *core) handleUnminedBlock(unminedBlock *types.Block) error {
	logger := c.logger.New("step", c.step, "height", c.currentRoundState.height)

	if err := c.checkUnminedBlockMsg(unminedBlock); err != nil {
		if err == errInvalidMessage {
			logger.Warn("invalid unminedBlock")
			return err
		}
		logger.Warn("unexpected unminedBlock", "err", err, "number", unminedBlock.Number(), "hash", unminedBlock.Hash())
		return err
	}

	logger.Trace("handleUnminedBlock", "number", unminedBlock.Number(), "hash", unminedBlock.Hash())

	c.latestPendingUnminedBlock = unminedBlock
	// TODO: remove, we should not be sending a proposal from handleUnminedBlock
	if c.step == StepAcceptProposal {
		c.sendProposal(unminedBlock)
	}
	return nil
}

// check request step
// return errInvalidMessage if the message is invalid
// return errFutureMessage if the height of proposal is larger than currentRoundState height
// return errOldMessage if the height of proposal is smaller than currentRoundState height
func (c *core) checkUnminedBlockMsg(unminedBlock *types.Block) error {
	if unminedBlock == nil {
		return errInvalidMessage
	}

	if c := c.currentRoundState.height.Cmp(unminedBlock.Number()); c > 0 {
		return errOldMessage
	} else if c < 0 {
		return errFutureMessage
	} else {
		return nil
	}
}

func (c *core) storeUnminedBlockMsg(unminedBlock *types.Block) {
	logger := c.logger.New("step", c.step)

	logger.Trace("Store future unminedBlock", "number", unminedBlock.Number(), "hash", unminedBlock.Hash())

	c.pendingUnminedBlocksMu.Lock()
	defer c.pendingUnminedBlocksMu.Unlock()

	c.pendingUnminedBlocks.Push(unminedBlock, float32(-unminedBlock.Number().Int64()))
}

func (c *core) processPendingRequests() {
	c.pendingUnminedBlocksMu.Lock()
	defer c.pendingUnminedBlocksMu.Unlock()

	for !(c.pendingUnminedBlocks.Empty()) {
		m, prio := c.pendingUnminedBlocks.Pop()
		ub, ok := m.(*types.Block)
		if !ok {
			c.logger.Warn("Malformed request, skip", "msg", m)
			continue
		}
		// Push back if it's a future message
		err := c.checkUnminedBlockMsg(ub)
		if err != nil {
			if err == errFutureMessage {
				c.logger.Trace("Stop processing request", "number", ub.Number(), "hash", ub.Hash())
				c.pendingUnminedBlocks.Push(m, prio)
				break
			}
			c.logger.Trace("Skip the pending request", "number", ub.Number(), "hash", ub.Hash(), "err", err)
			continue
		}
		c.logger.Trace("Post pending request", "number", ub.Number(), "hash", ub.Hash())

		go c.sendEvent(tendermint.NewUnminedBlockEvent{
			NewUnminedBlock: *ub,
		})
	}
}

//---------------------------------------Propose---------------------------------------

// TODO: add new message struct for proposal (proposalMessage) and determine how to rlp encode them especially nil
// TODO: add new message for vote (prevote and precommit) and determine how to rlp encode them especially nil
func (c *core) sendProposal(p *types.Block) {
	logger := c.logger.New("step", c.step)

	// If I'm the proposer and I have the same height with the proposal
	if c.currentRoundState.Height().Cmp(p.Number()) == 0 && c.isProposer() && !c.sentProposal {
		r, h, vr := c.currentRoundState.Round(), c.currentRoundState.Height(), c.validRound
		proposal, err := Encode(&tendermint.Proposal{
			Round:         r,
			Height:        h,
			ValidRound:    vr,
			ProposalBlock: *p,
		})
		if err != nil {
			logger.Error("Failed to encode", "Round", r, "Height", h, "ValidRound", vr)
			return
		}
		c.sentProposal = true
		c.backend.SetProposedBlockHash(p.Hash())
		c.broadcast(&message{
			Code: msgProposal,
			Msg:  proposal,
		})
	}
}

func (c *core) handleProposal(msg *message, src tendermint.Validator) error {
	logger := c.logger.New("from", src, "step", c.step)

	var proposal *tendermint.Proposal
	err := msg.Decode(&proposal)
	if err != nil {
		return errFailedDecodeProposal
	}

	// Ensure we have the same view with the Proposal message
	// If it is old message, see if we need to broadcast COMMIT
	//TODO: fixup
	if err := c.checkMessage(msgProposal, proposal.Round, proposal.Height); err != nil {
		if err == errOldMessage {
			// TODO: keeping it for the time being but rebroadcasting based on old messages should not be required due to partial synchrony assumption
			// TODO: also need to add previous round messages to currentHeightRoundStates and only rebroadcast if older height
			valSet := c.backend.Validators(proposal.ProposalBlock.Number().Uint64()).Copy()
			previousProposer := c.backend.GetProposer(proposal.ProposalBlock.Number().Uint64() - 1)
			valSet.CalcProposer(previousProposer, proposal.Round.Uint64())
			// Broadcast COMMIT if it is an existing block
			// 1. The proposer needs to be a proposer matches the given (Height + Round)
			// 2. The given block must exist
			if valSet.IsProposer(src.Address()) && c.backend.HasPropsal(proposal.ProposalBlock.Hash(), proposal.ProposalBlock.Number()) {
				c.sendPrecommitForOldBlock(proposal.Round, proposal.Height, proposal.ProposalBlock.Hash())
				return nil
			}
		}
		return err
	}

	// Check if the message comes from currentRoundState proposer
	if !c.valSet.IsProposer(src.Address()) {
		logger.Warn("Ignore proposal messages from non-proposer")
		return errNotFromProposer
	}

	// Verify the proposal we received
	if duration, err := c.backend.Verify(proposal.ProposalBlock); err != nil {
		logger.Warn("Failed to verify proposal", "err", err, "duration", duration)
		// if it's a future block, we will handle it again after the duration
		// TIME FIELD OF HEADER CHECKED HERE - NOT HEIGHT
		if err == consensus.ErrFutureBlock {
			c.stopFutureProposalTimer()
			c.futureProposalTimer = time.AfterFunc(duration, func() {
				c.sendEvent(backlogEvent{
					src: src,
					msg: msg,
				})
			})
		}
		return err
	}

	// Here is about to accept the Proposal
	if c.step == StepAcceptProposal {
		c.acceptProposal(proposal)

		if proposal.ValidRound.Int64() == -1 {
			if c.lockedRound == proposal.ValidRound || proposal.ProposalBlock.Hash() == c.lockedValue.Hash() {
				c.sendPrevote(false)
			} else {
				c.sendPrevote(true)
			}
		} else if proposal.ValidRound.Int64() > -1 {

		}
		c.setStep(StepProposeDone)
	}

	return nil
}

func (c *core) acceptProposal(proposal *tendermint.Proposal) {
	c.currentRoundState.SetProposal(proposal)
}

//---------------------------------------Prevote---------------------------------------

func (c *core) sendPrevote(isNil bool) {
	logger := c.logger.New("step", c.step)

	var sub = &tendermint.Subject{
		Round:  big.NewInt(c.currentRoundState.round.Int64()),
		Height: big.NewInt(c.currentRoundState.Height().Int64()),
	}

	if isNil {
		sub.Digest = common.Hash{}
	} else {
		sub.Digest = c.currentRoundState.Proposal().ProposalBlock.Hash()
	}

	encodedSubject, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", "subject", sub)
		return
	}
	c.broadcast(&message{
		Code: msgPrevote,
		Msg:  encodedSubject,
	})
}

func (c *core) handlePrevote(msg *message, src tendermint.Validator) error {
	// Decode PREPARE message
	var prepare *tendermint.Subject
	err := msg.Decode(&prepare)
	if err != nil {
		return errFailedDecodePrevote
	}

	if err = c.checkMessage(msgPrevote, prepare.Round, prepare.Height); err != nil {
		return err
	}

	// If it is locked, it can only process on the locked block.
	// Passing verifyPrevote and checkMessage implies it is processing on the locked block since it was verified in the Proposald step.
	if err = c.verifyPrevote(prepare, src); err != nil {
		return err
	}

	err = c.acceptPrevote(msg)
	if err != nil {
		c.logger.Error("Failed to add PREPARE message to round step",
			"from", src, "step", c.step, "msg", msg, "err", err)
	}

	// Change to Prevoted step if we've received enough PREPARE messages or it is locked
	// and we are in earlier step before Prevoted step.
	if ((c.currentRoundState.IsHashLocked() && prepare.Digest == c.currentRoundState.GetLockedHash()) || c.currentRoundState.GetPrevoteOrPrecommitSize() > 2*c.valSet.F()) &&
		c.step.Cmp(StepPrevoteDone) < 0 {
		c.currentRoundState.LockHash()
		c.setStep(StepPrevoteDone)
		c.sendPrecommit()
	}

	return nil
}

// verifyPrevote verifies if the received PREPARE message is equivalent to our subject
func (c *core) verifyPrevote(prepare *tendermint.Subject, src tendermint.Validator) error {
	logger := c.logger.New("from", src, "step", c.step)

	sub := c.currentRoundState.Subject()
	if !reflect.DeepEqual(prepare, sub) {
		logger.Warn("Inconsistent subjects between PREPARE and proposal", "expected", sub, "got", prepare)
		return errInconsistentSubject
	}

	return nil
}

func (c *core) acceptPrevote(msg *message) error {
	// Add the PREPARE message to currentRoundState round step
	if err := c.currentRoundState.Prevotes.Add(msg); err != nil {
		return err
	}

	return nil
}

//---------------------------------------Precommit---------------------------------------
func (c *core) sendPrecommit() {
	sub := c.currentRoundState.Subject()
	c.broadcastPrecommit(sub)
}

func (c *core) sendPrecommitForOldBlock(r *big.Int, h *big.Int, digest common.Hash) {
	sub := &tendermint.Subject{
		Round:  r,
		Height: h,
		Digest: digest,
	}
	c.broadcastPrecommit(sub)
}

func (c *core) broadcastPrecommit(sub *tendermint.Subject) {
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
	var commit *tendermint.Subject
	err := msg.Decode(&commit)
	if err != nil {
		return errFailedDecodePrecommit
	}

	if err := c.checkMessage(msgPrecommit, commit.Round, commit.Height); err != nil {
		return err
	}

	if err := c.verifyPrecommit(commit, src); err != nil {
		return err
	}

	if err := c.acceptPrecommit(msg); err != nil {
		c.logger.Error("Failed to record commit message", "from", src, "step", c.step, "msg", msg, "err", err)
	}

	// Precommit the proposal once we have enough COMMIT messages and we are not in the Committed step.
	//
	// If we already have a proposal, we may have chance to speed up the consensus process
	// by committing the proposal without PREPARE messages.
	if c.currentRoundState.Precommits.Size() > 2*c.valSet.F() && c.step.Cmp(StepPrecommitDone) < 0 {
		// Still need to call LockHash here since step can skip Prevoted step and jump directly to the Committed step.
		c.currentRoundState.LockHash()
		c.commit()
	}

	return nil
}

// verifyPrecommit verifies if the received COMMIT message is equivalent to our subject
func (c *core) verifyPrecommit(commit *tendermint.Subject, src tendermint.Validator) error {
	logger := c.logger.New("from", src, "step", c.step)

	sub := c.currentRoundState.Subject()
	if !reflect.DeepEqual(commit, sub) {
		logger.Warn("Inconsistent subjects between commit and proposal", "expected", sub, "got", commit)
		return errInconsistentSubject
	}

	return nil
}

// acceptPrecommit adds the COMMIT message to currentRoundState round step
func (c *core) acceptPrecommit(msg *message) error {
	return c.currentRoundState.Precommits.Add(msg)
}

//---------------------------------------Commit---------------------------------------
func (c *core) handleCommit() {
	c.logger.Trace("Received a final committed proposal", "step", c.step)
	c.startRound(common.Big0)
}
