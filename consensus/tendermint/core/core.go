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

//---------------------------------------Backlog---------------------------------------

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
