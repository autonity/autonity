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
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
)

var (
	initialProposeTimeout   = 5 * time.Second
	initialPrevoteTimeout   = 5 * time.Second
	initialPrecommitTimeout = 5 * time.Second
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

	roundChangeSet     *roundChangeSet
	roundChangeTimer   *time.Timer
	roundChangeTimerMu sync.RWMutex

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
			c.sendNextRoundChange()
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

	c.roundChangeTimerMu.RLock()
	defer c.roundChangeTimerMu.RUnlock()
	if c.roundChangeTimer != nil {
		c.roundChangeTimer.Stop()
	}
}

func (c *core) newRoundChangeTimer() {
	c.stopTimer()

	// set timeout based on the round number
	timeout := time.Duration(c.config.RequestTimeout) * time.Millisecond
	round := c.current.Round().Uint64()
	if round > 0 {
		timeout += time.Duration(math.Pow(2, float64(round))) * time.Second
	}

	c.roundChangeTimerMu.Lock()
	defer c.roundChangeTimerMu.Unlock()
	c.roundChangeTimer = time.AfterFunc(timeout, func() {
		c.sendEvent(timeoutEvent{})
	})
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
