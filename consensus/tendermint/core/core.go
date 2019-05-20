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
	"github.com/clearmatics/autonity/metrics"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
)

// New creates an Istanbul consensus core
func New(backend tendermint.Backend, config *tendermint.Config) Engine {
	c := &core{
		config:             config,
		address:            backend.Address(),
		state:              StateAcceptRequest,
		handlerStopCh:      make(chan struct{}),
		logger:             log.New("address", backend.Address()),
		backend:            backend,
		backlogs:           make(map[tendermint.Validator]*prque.Prque),
		backlogsMu:         new(sync.Mutex),
		pendingRequests:    prque.New(),
		pendingRequestsMu:  new(sync.Mutex),
		consensusTimestamp: time.Time{},
		roundMeter:         metrics.NewRegisteredMeter("consensus/tendermint/core/round", nil),
		sequenceMeter:      metrics.NewRegisteredMeter("consensus/tendermint/core/sequence", nil),
		consensusTimer:     metrics.NewRegisteredTimer("consensus/tendermint/core/consensus", nil),
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

	valSet                tendermint.ValidatorSet
	waitingForRoundChange bool
	validateFn            func([]byte, []byte) (common.Address, error)

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

	// TODO: remove these ethstats from the code temporarily since it probably needs to be implemented again
	consensusTimestamp time.Time
	// the meter to record the round change rate
	roundMeter metrics.Meter
	// the meter to record the sequence update rate
	sequenceMeter metrics.Meter
	// the timer to record consensus duration (from accepting a proposal to final committed stage)
	consensusTimer metrics.Timer

	sentProposal bool

	lockedRound *big.Int
	validRound  *big.Int
	lockedValue tendermint.ProposalBlock
	validValue  tendermint.ProposalBlock

	currentHeightRoundsStates []roundState

	// TODO: may require a mutex
	latestPendingRequest *tendermint.Request
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
	logger := c.logger.New("Previous Round", round.Uint64()-1, "Previous Height", height.Uint64()-1)

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
		c.scheduleTimeoutPropose()
	}

	// TODO remove the following
	if c.current == nil {
		logger = c.logger.New("old_round", -1, "old_seq", 0)
	} else {
		logger = c.logger.New("old_round", c.current.Round(), "old_seq", c.current.Sequence())
	}

	roundChange := false
	// Try to get last proposal
	lastProposal, lastProposer := c.backend.LastProposal()
	if c.current == nil {
		logger.Trace("Start to the initial round")
	} else if lastProposal.Number().Cmp(c.current.Sequence()) >= 0 {
		diff := new(big.Int).Sub(lastProposal.Number(), c.current.Sequence())
		c.sequenceMeter.Mark(new(big.Int).Add(diff, common.Big1).Int64())

		if !c.consensusTimestamp.IsZero() {
			c.consensusTimer.UpdateSince(c.consensusTimestamp)
			c.consensusTimestamp = time.Time{}
		}
		logger.Trace("Catch up latest proposal", "number", lastProposal.Number().Uint64(), "hash", lastProposal.Hash())
	} else if lastProposal.Number().Cmp(big.NewInt(c.current.Sequence().Int64()-1)) == 0 {
		if round.Cmp(common.Big0) == 0 { //TODO : Bug ?
			// same seq and round, don't need to start new round
			return
		} else if round.Cmp(c.current.Round()) < 0 {
			logger.Warn("New round should not be smaller than current round", "seq", lastProposal.Number().Int64(), "new_round", round, "old_round", c.current.Round())
			return
		}
		roundChange = true
	} else {
		logger.Warn("New sequence should be larger than current sequence", "new_seq", lastProposal.Number().Int64())
		return
	}

	var newView *tendermint.View
	if roundChange {
		newView = &tendermint.View{
			Sequence: new(big.Int).Set(c.current.Sequence()),
			Round:    new(big.Int).Set(round),
		}
	} else {
		newView = &tendermint.View{
			Sequence: new(big.Int).Add(lastProposal.Number(), common.Big1),
			Round:    new(big.Int),
		}
		c.valSet = c.backend.Validators(lastProposal.Number().Uint64() + 1)
	}

	// Update logger
	logger = logger.New("old_proposer", c.valSet.GetProposer())
	// Clear invalid ROUND CHANGE messages
	c.roundChangeSet = newRoundChangeSet(c.valSet)
	// New snapshot for new round
	c.updateRoundState(newView, c.valSet, roundChange)
	// Calculate new proposer
	c.valSet.CalcProposer(lastProposer, newView.Round.Uint64())
	c.waitingForRoundChange = false
	c.sentProposal = false
	c.setState(StateAcceptRequest)
	if roundChange && c.isProposer() && c.current != nil {
		// If it is locked, propose the old proposal
		// If we have pending request, propose pending request
		if c.current.IsHashLocked() {
			r := &tendermint.Request{
				ProposalBlock: c.current.Proposal().ProposalBlock, //c.current.ProposalBlock would be the locked proposal by previous proposer, see updateRoundState
			}
			c.sendProposal(r)
		} else if c.current.pendingRequest != nil {
			c.sendProposal(c.current.pendingRequest)
		}
	}
	c.newRoundChangeTimer()

	logger.Debug("New round", "new_round", newView.Round, "new_seq", newView.Sequence, "new_proposer", c.valSet.GetProposer(), "valSet", c.valSet.List(), "size", c.valSet.Size(), "isProposer", c.isProposer())
}

func (c *core) catchUpRound(view *tendermint.View) {
	logger := c.logger.New("old_round", c.current.Round(), "old_seq", c.current.Sequence(), "old_proposer", c.valSet.GetProposer())

	if view.Round.Cmp(c.current.Round()) > 0 {
		c.roundMeter.Mark(new(big.Int).Sub(view.Round, c.current.Round()).Int64())
	}
	c.waitingForRoundChange = true

	// Need to keep block locked for round catching up
	c.updateRoundState(view, c.valSet, true)
	c.roundChangeSet.Clear(view.Round)
	c.newRoundChangeTimer()

	logger.Trace("Catch up round", "new_round", view.Round, "new_seq", view.Sequence, "new_proposer", c.valSet)
}

// updateRoundState updates round state by checking if locking block is necessary
func (c *core) updateRoundState(view *tendermint.View, validatorSet tendermint.ValidatorSet, roundChange bool) {
	// Lock only if both roundChange is true and it is locked
	if roundChange && c.current != nil {
		if c.current.IsHashLocked() {
			c.current = newRoundState(view, validatorSet, c.current.GetLockedHash(), c.current.Proposal(), c.current.pendingRequest, c.backend.HasBadProposal)
		} else {
			c.current = newRoundState(view, validatorSet, common.Hash{}, nil, c.current.pendingRequest, c.backend.HasBadProposal)
		}
	} else {
		c.current = newRoundState(view, validatorSet, common.Hash{}, nil, nil, c.backend.HasBadProposal)
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

func (c *core) scheduleTimeoutPropose() {
	// TODO: add a timeoutField to the core struct and us the time lib to start the timeout
	// TODO: implement OnTimeoutPropose
}

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal(hash common.Hash) []byte {
	var buf bytes.Buffer
	buf.Write(hash.Bytes())
	buf.Write([]byte{byte(msgPrecommit)})
	return buf.Bytes()
}
