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
	initialProposeTimeout   = 5 * time.Second
	initialPrevoteTimeout   = 5 * time.Second
	initialPrecommitTimeout = 5 * time.Second
)

var (
	// errInconsistentSubject is returned when received subject is different from
	// current subject.
	errInconsistentSubject = errors.New("inconsistent subjects")
	// errNotFromProposer is returned when received message is supposed to be from
	// proposer.
	errNotFromProposer = errors.New("message does not come from proposer")
	// errIgnored is returned when a message was ignored.
	errIgnored = errors.New("message is ignored")
	// errFutureMessage is returned when current view is earlier than the
	// view of the received message.
	errFutureMessage = errors.New("future message")
	// errOldMessage is returned when the received message's view is earlier
	// than current view.
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
		config:            config,
		address:           backend.Address(),
		state:             StateAcceptRequest,
		handlerStopCh:     make(chan struct{}),
		logger:            log.New("address", backend.Address()),
		backend:           backend,
		backlogs:          make(map[tendermint.Validator]*prque.Prque),
		backlogsMu:        new(sync.Mutex),
		pendingRequests:   prque.New(),
		pendingRequestsMu: new(sync.Mutex),
		proposeTimeout:    new(timeout),
		prevoteTimeout:    new(timeout),
		precommitTimeout:  new(timeout),
	}
	c.validateFn = c.checkValidatorSignature
	return c
}

// ----------------------------------------------------------------------------

type core struct {
	config  *tendermint.Config
	address common.Address
	// TODO change the name to step Step
	state  State
	logger log.Logger

	backend             tendermint.Backend
	events              *event.TypeMuxSubscription
	finalCommittedSub   *event.TypeMuxSubscription
	timeoutSub          *event.TypeMuxSubscription
	futureProposalTimer *time.Timer

	valSet     tendermint.ValidatorSet
	validateFn func([]byte, []byte) (common.Address, error)

	backlogs   map[tendermint.Validator]*prque.Prque
	backlogsMu *sync.Mutex

	// TODO: update the name to currentRoundState
	current       *roundState
	handlerStopCh chan struct{}

	pendingRequests   *prque.Prque
	pendingRequestsMu *sync.Mutex

	sentProposal bool

	lockedRound *big.Int
	validRound  *big.Int
	lockedValue tendermint.ProposalBlock
	validValue  tendermint.ProposalBlock

	currentHeightRoundsStates []roundState

	// TODO: may require a mutex
	latestPendingRequest *tendermint.Request

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
	if msg.Code == msgPrecommit && c.current.Proposal() != nil {
		seal := PrepareCommittedSeal(c.current.Proposal().ProposalBlock.Hash())
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
	logger := c.logger.New("state", c.state)

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

func (c *core) currentView() *tendermint.View {
	return &tendermint.View{
		Sequence: new(big.Int).Set(c.current.Sequence()),
		Round:    new(big.Int).Set(c.current.Round()),
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
	c.setState(StatePrecommitDone)

	proposal := c.current.Proposal()
	if proposal != nil {
		committedSeals := make([][]byte, c.current.Precommits.Size())
		for i, v := range c.current.Precommits.Values() {
			committedSeals[i] = make([]byte, types.PoSExtraSeal)
			copy(committedSeals[i][:], v.CommittedSeal[:])
		}

		if err := c.backend.Precommit(proposal.ProposalBlock, committedSeals); err != nil {
			c.current.UnlockHash() //Unlock block when insertion fails
			// TODO: go to next height
			return
		}
	}
}

// startNewRound starts a new round. if round equals to 0, it means to starts a new sequence
// TODO: change name to startRound
func (c *core) startNewRound(round *big.Int) {
	//TODO: update the name of lastProposalBlock and LastBlockProposal()
	lastProposalBlock, lastProposalBlockProposer := c.backend.LastProposal()
	height := new(big.Int).Add(lastProposalBlock.Number(), common.Big1)

	// Start of new height where round is 0
	if round.Uint64() == 0 {
		// Set the shared round values to initial values
		c.lockedRound = nil
		c.lockedValue = nil
		c.validRound = nil
		c.validValue = nil

		c.valSet = c.backend.Validators(height.Uint64())

		// TODO: Assuming that round == 0 only when the node moves to a new height, need to confirm where exactly the node moves to a new height
		c.currentHeightRoundsStates = nil

	} else {
		// Assuming the above values will be set for round > 0
		// Add the current round state to the core previous round states
		c.currentHeightRoundsStates = append(c.currentHeightRoundsStates, *c.current)
	}

	// Update the current round state
	curView := tendermint.View{
		Round:    new(big.Int).Set(round),
		Sequence: new(big.Int).Set(height),
	}

	c.current = newRoundState(
		&curView,
		c.valSet,
		common.Hash{},
		nil,
		nil,
		c.backend.HasBadProposal,
	)

	c.valSet.CalcProposer(lastProposalBlockProposer, round.Uint64())
	c.sentProposal = false
	// c.setState(StateAcceptRequest) will process the pending request sent by the backed.Seal() and set c.lastestPendingRequest
	c.setState(StateAcceptRequest)

	var proposalRequest *tendermint.Request
	if c.isProposer() {
		if c.validValue != nil {
			proposalRequest = &tendermint.Request{ProposalBlock: c.validValue}
		} else {
			proposalRequest = c.latestPendingRequest
		}
		c.sendProposal(proposalRequest)
	} else {
		c.proposeTimeout.scheduleTimeout(timeoutPropose(round.Int64()), c.onTimeoutPropose)
	}
}

func (c *core) setState(state State) {
	if c.state != state {
		c.state = state
	}
	if state == StateAcceptRequest {
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

func (c *core) onTimeoutPrevote() {
}

func (c *core) onTimeoutPrecommit() {

}

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

func (t *timeout) stopTimer() bool {
	t.RLock()
	defer t.RUnlock()
	return t.timer.Stop()
}

// The timeout may need to be changed depending on the State
func timeoutPropose(round int64) time.Duration {
	return initialProposeTimeout + time.Duration(round)*time.Second
}

func timeoutPrevote(round int64) time.Duration {
	return initialProposeTimeout + time.Duration(round)*time.Second
}

func timeoutPrecommit(round int64) time.Duration {
	return initialProposeTimeout + time.Duration(round)*time.Second
}

//---------------------------------------Handler---------------------------------------

// Start implements core.Engine.Start
func (c *core) Start() error {
	// Start a new round from last sequence + 1
	c.startNewRound(common.Big0)

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
		tendermint.RequestEvent{},
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
	// Clear state
	defer func() {
		c.current = nil
		<-c.handlerStopCh
	}()

	for {
		select {
		case event, ok := <-c.events.Chan():
			if !ok {
				return
			}
			// A real event arrived, process interesting content
			switch ev := event.Data.(type) {
			case tendermint.RequestEvent:
				r := &tendermint.Request{
					ProposalBlock: ev.ProposalBlock,
				}
				err := c.handleRequest(r)
				if err == errFutureMessage {
					c.storeRequestMsg(r)
				}
			case tendermint.MessageEvent:
				if err := c.handleMsg(ev.Payload); err == nil {
					c.backend.Gossip(c.valSet, ev.Payload)
				}
			case backlogEvent:
				// No need to check signature for internal messages
				if err := c.handleCheckedMsg(ev.msg, ev.src); err == nil {
					p, err := ev.msg.Payload()
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
		case event, ok := <-c.finalCommittedSub.Chan():
			if !ok {
				return
			}
			switch event.Data.(type) {
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
	// If we're not waiting for round change yet, we can try to catch up
	// the max round with F+1 round change message. We only need to catch up
	// if the max round is larger than current round.

	lastProposal, _ := c.backend.LastProposal()
	if lastProposal != nil && lastProposal.Number().Cmp(c.current.Sequence()) >= 0 {
		c.logger.Trace("round change timeout, catch up latest sequence", "number", lastProposal.Number().Uint64())
		c.startNewRound(common.Big0)
	}
}

//---------------------------------------Backlog---------------------------------------

type backlogEvent struct {
	src tendermint.Validator
	msg *message
}

// checkMessage checks the message state
// return errInvalidMessage if the message is invalid
// return errFutureMessage if the message view is larger than current view
// return errOldMessage if the message view is smaller than current view
func (c *core) checkMessage(msgCode uint64, view *tendermint.View) error {
	if view == nil || view.Sequence == nil || view.Round == nil {
		return errInvalidMessage
	}

	if view.Cmp(c.currentView()) > 0 {
		return errFutureMessage
	}

	if view.Cmp(c.currentView()) < 0 {
		return errOldMessage
	}

	// StateAcceptRequest only accepts msgProposal
	// other messages are future messages
	if c.state == StateAcceptRequest {
		if msgCode > msgProposal {
			return errFutureMessage
		}
		return nil
	}

	// For states(StateProposeDone, StatePrevoteDone, StatePrecommitDone),
	// can accept all message types if processing with same view
	return nil
}

func (c *core) storeBacklog(msg *message, src tendermint.Validator) {
	logger := c.logger.New("from", src, "state", c.state)

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
			backlog.Push(msg, toPriority(msg.Code, p.View))
		}
		// for msgRoundChange, msgPrevote and msgPrecommit cases
	default:
		var p *tendermint.Subject
		err := msg.Decode(&p)
		if err == nil {
			backlog.Push(msg, toPriority(msg.Code, p.View))
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

		logger := c.logger.New("from", src, "state", c.state)
		var isFuture bool

		// We stop processing if
		//   1. backlog is empty
		//   2. The first message in queue is a future message
		for !(backlog.Empty() || isFuture) {
			m, prio := backlog.Pop()
			msg := m.(*message)
			var view *tendermint.View
			switch msg.Code {
			case msgProposal:
				var m *tendermint.Proposal
				err := msg.Decode(&m)
				if err == nil {
					view = m.View
				}
				// for msgRoundChange, msgPrevote and msgPrecommit cases
			default:
				var sub *tendermint.Subject
				err := msg.Decode(&sub)
				if err == nil {
					view = sub.View
				}
			}
			if view == nil {
				logger.Debug("Nil view", "msg", msg)
				continue
			}
			// Push back if it's a future message
			err := c.checkMessage(msg.Code, view)
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

func toPriority(msgCode uint64, view *tendermint.View) float32 {
	// FIXME: round will be reset as 0 while new sequence
	// 10 * Round limits the range of message code is from 0 to 9
	// 1000 * Sequence limits the range of round is from 0 to 99
	return -float32(view.Sequence.Uint64()*1000 + view.Round.Uint64()*10 + uint64(msgPriority[msgCode]))
}

//---------------------------------------Request---------------------------------------

func (c *core) handleRequest(request *tendermint.Request) error {
	logger := c.logger.New("state", c.state, "seq", c.current.sequence)

	if err := c.checkRequestMsg(request); err != nil {
		if err == errInvalidMessage {
			logger.Warn("invalid request")
			return err
		}
		logger.Warn("unexpected request", "err", err, "number", request.ProposalBlock.Number(), "hash", request.ProposalBlock.Hash())
		return err
	}

	logger.Trace("handleRequest", "number", request.ProposalBlock.Number(), "hash", request.ProposalBlock.Hash())

	c.latestPendingRequest = request
	// TODO: remove, we should not be sending a proposal from handleRequest
	if c.state == StateAcceptRequest {
		c.sendProposal(request)
	}
	return nil
}

// check request state
// return errInvalidMessage if the message is invalid
// return errFutureMessage if the sequence of proposal is larger than current sequence
// return errOldMessage if the sequence of proposal is smaller than current sequence
func (c *core) checkRequestMsg(request *tendermint.Request) error {
	if request == nil || request.ProposalBlock == nil {
		return errInvalidMessage
	}

	if c := c.current.sequence.Cmp(request.ProposalBlock.Number()); c > 0 {
		return errOldMessage
	} else if c < 0 {
		return errFutureMessage
	} else {
		return nil
	}
}

func (c *core) storeRequestMsg(request *tendermint.Request) {
	logger := c.logger.New("state", c.state)

	logger.Trace("Store future request", "number", request.ProposalBlock.Number(), "hash", request.ProposalBlock.Hash())

	c.pendingRequestsMu.Lock()
	defer c.pendingRequestsMu.Unlock()

	c.pendingRequests.Push(request, float32(-request.ProposalBlock.Number().Int64()))
}

func (c *core) processPendingRequests() {
	c.pendingRequestsMu.Lock()
	defer c.pendingRequestsMu.Unlock()

	for !(c.pendingRequests.Empty()) {
		m, prio := c.pendingRequests.Pop()
		r, ok := m.(*tendermint.Request)
		if !ok {
			c.logger.Warn("Malformed request, skip", "msg", m)
			continue
		}
		// Push back if it's a future message
		err := c.checkRequestMsg(r)
		if err != nil {
			if err == errFutureMessage {
				c.logger.Trace("Stop processing request", "number", r.ProposalBlock.Number(), "hash", r.ProposalBlock.Hash())
				c.pendingRequests.Push(m, prio)
				break
			}
			c.logger.Trace("Skip the pending request", "number", r.ProposalBlock.Number(), "hash", r.ProposalBlock.Hash(), "err", err)
			continue
		}
		c.logger.Trace("Post pending request", "number", r.ProposalBlock.Number(), "hash", r.ProposalBlock.Hash())

		go c.sendEvent(tendermint.RequestEvent{
			ProposalBlock: r.ProposalBlock,
		})
	}
}

//---------------------------------------Propose---------------------------------------

// TODO: add new message struct for proposal (proposalMessage) and determine how to rlp encode them especially nil
// TODO: add new message for vote (prevote and precommit) and determine how to rlp encode them especially nil
func (c *core) sendProposal(request *tendermint.Request) {
	logger := c.logger.New("state", c.state)

	// If I'm the proposer and I have the same sequence with the proposal
	if c.current.Sequence().Cmp(request.ProposalBlock.Number()) == 0 && c.isProposer() && !c.sentProposal {
		curView := c.currentView()
		proposal, err := Encode(&tendermint.Proposal{
			View:          curView,
			ProposalBlock: request.ProposalBlock,
		})
		if err != nil {
			logger.Error("Failed to encode", "view", curView)
			return
		}
		c.sentProposal = true
		c.backend.SetProposedBlockHash(request.ProposalBlock.Hash())
		c.broadcast(&message{
			Code: msgProposal,
			Msg:  proposal,
		})
	}
}

func (c *core) handleProposal(msg *message, src tendermint.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Decode PRE-PREPARE
	var proposal *tendermint.Proposal
	err := msg.Decode(&proposal)
	if err != nil {
		return errFailedDecodeProposal
	}

	// Ensure we have the same view with the PRE-PREPARE message
	// If it is old message, see if we need to broadcast COMMIT
	if err := c.checkMessage(msgProposal, proposal.View); err != nil {
		if err == errOldMessage {
			//TODO : EIP Says ignore?
			// Get validator set for the given proposal
			valSet := c.backend.Validators(proposal.ProposalBlock.Number().Uint64()).Copy()
			previousProposer := c.backend.GetProposer(proposal.ProposalBlock.Number().Uint64() - 1)
			valSet.CalcProposer(previousProposer, proposal.View.Round.Uint64())
			// Broadcast COMMIT if it is an existing block
			// 1. The proposer needs to be a proposer matches the given (Sequence + Round)
			// 2. The given block must exist
			if valSet.IsProposer(src.Address()) && c.backend.HasPropsal(proposal.ProposalBlock.Hash(), proposal.ProposalBlock.Number()) {
				c.sendPrecommitForOldBlock(proposal.View, proposal.ProposalBlock.Hash())
				return nil
			}
		}
		return err
	}

	// Check if the message comes from current proposer
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
		} else {
			// TODO: possibly send propose(nil) (need to update)
		}
		return err
	}

	// Here is about to accept the PRE-PREPARE
	if c.state == StateAcceptRequest {
		// Send ROUND CHANGE if the locked proposal and the received proposal are different
		if c.current.IsHashLocked() {
			if proposal.ProposalBlock.Hash() == c.current.GetLockedHash() {
				// Broadcast COMMIT and enters Prevoted state directly
				c.acceptProposal(proposal)
				c.setState(StatePrevoteDone)
				c.sendPrecommit() // TODO : double check, why not PREPARE?
			} else {
				// TODO: possibly send propose(nil) (need to update)
			}
		} else {
			// Either
			//   1. the locked proposal and the received proposal match
			//   2. we have no locked proposal
			c.acceptProposal(proposal)
			c.setState(StateProposeDone)
			c.sendPrevote()
		}
	}

	return nil
}

func (c *core) acceptProposal(proposal *tendermint.Proposal) {
	c.current.SetProposal(proposal)
}

//---------------------------------------Prevote---------------------------------------

func (c *core) sendPrevote() {
	logger := c.logger.New("state", c.state)

	sub := c.current.Subject()
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

	if err = c.checkMessage(msgPrevote, prepare.View); err != nil {
		return err
	}

	// If it is locked, it can only process on the locked block.
	// Passing verifyPrevote and checkMessage implies it is processing on the locked block since it was verified in the Proposald state.
	if err = c.verifyPrevote(prepare, src); err != nil {
		return err
	}

	err = c.acceptPrevote(msg)
	if err != nil {
		c.logger.Error("Failed to add PREPARE message to round state",
			"from", src, "state", c.state, "msg", msg, "err", err)
	}

	// Change to Prevoted state if we've received enough PREPARE messages or it is locked
	// and we are in earlier state before Prevoted state.
	if ((c.current.IsHashLocked() && prepare.Digest == c.current.GetLockedHash()) || c.current.GetPrevoteOrPrecommitSize() > 2*c.valSet.F()) &&
		c.state.Cmp(StatePrevoteDone) < 0 {
		c.current.LockHash()
		c.setState(StatePrevoteDone)
		c.sendPrecommit()
	}

	return nil
}

// verifyPrevote verifies if the received PREPARE message is equivalent to our subject
func (c *core) verifyPrevote(prepare *tendermint.Subject, src tendermint.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	sub := c.current.Subject()
	if !reflect.DeepEqual(prepare, sub) {
		logger.Warn("Inconsistent subjects between PREPARE and proposal", "expected", sub, "got", prepare)
		return errInconsistentSubject
	}

	return nil
}

func (c *core) acceptPrevote(msg *message) error {
	// Add the PREPARE message to current round state
	if err := c.current.Prevotes.Add(msg); err != nil {
		return err
	}

	return nil
}

//---------------------------------------Precommit---------------------------------------
func (c *core) sendPrecommit() {
	sub := c.current.Subject()
	c.broadcastPrecommit(sub)
}

func (c *core) sendPrecommitForOldBlock(view *tendermint.View, digest common.Hash) {
	sub := &tendermint.Subject{
		View:   view,
		Digest: digest,
	}
	c.broadcastPrecommit(sub)
}

func (c *core) broadcastPrecommit(sub *tendermint.Subject) {
	logger := c.logger.New("state", c.state)

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

	if err := c.checkMessage(msgPrecommit, commit.View); err != nil {
		return err
	}

	if err := c.verifyPrecommit(commit, src); err != nil {
		return err
	}

	if err := c.acceptPrecommit(msg); err != nil {
		c.logger.Error("Failed to record commit message", "from", src, "state", c.state, "msg", msg, "err", err)
	}

	// Precommit the proposal once we have enough COMMIT messages and we are not in the Committed state.
	//
	// If we already have a proposal, we may have chance to speed up the consensus process
	// by committing the proposal without PREPARE messages.
	if c.current.Precommits.Size() > 2*c.valSet.F() && c.state.Cmp(StatePrecommitDone) < 0 {
		// Still need to call LockHash here since state can skip Prevoted state and jump directly to the Committed state.
		c.current.LockHash()
		c.commit()
	}

	return nil
}

// verifyPrecommit verifies if the received COMMIT message is equivalent to our subject
func (c *core) verifyPrecommit(commit *tendermint.Subject, src tendermint.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	sub := c.current.Subject()
	if !reflect.DeepEqual(commit, sub) {
		logger.Warn("Inconsistent subjects between commit and proposal", "expected", sub, "got", commit)
		return errInconsistentSubject
	}

	return nil
}

// acceptPrecommit adds the COMMIT message to current round state
func (c *core) acceptPrecommit(msg *message) error {
	return c.current.Precommits.Add(msg)
}

//---------------------------------------Commit---------------------------------------
func (c *core) handleCommit() {
	c.logger.Trace("Received a final committed proposal", "state", c.state)
	c.startNewRound(common.Big0)
}
