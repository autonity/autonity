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
	// errInconsistentSubject is returned when received subject is different from
	// currentRoundState subject.
	//errInconsistentSubject = errors.New("inconsistent subjects")
	// errNotFromProposer is returned when received message is supposed to be from
	// proposer.
	errNotFromProposer = errors.New("message does not come from proposer")
	// errIgnored is returned when a message was ignored.
	//errIgnored = errors.New("message is ignored")
	// errFutureHeightMessage is returned when currentRoundState view is earlier than the
	// view of the received message.
	errFutureHeightMessage = errors.New("future height message")
	// errOldHeightMessage is returned when the received message's view is earlier
	// than currentRoundState view.
	errOldHeightMessage = errors.New("old height message")
	// errOldRoundMessage message is returned when message is of the same Height but form a smaller round
	errOldRoundMessage = errors.New("same height but old round message")
	// errFutureRoundMessage message is returned when message is of the same Height but form a newer round
	errFutureRoundMessage = errors.New("same height but future round message")
	// errInvalidMessage is returned when the message is malformed.
	errInvalidMessage = errors.New("invalid message")
	// errFailedDecodeProposal is returned when the PRE-PREPARE message is malformed.
	errFailedDecodeProposal = errors.New("failed to decode PRE-PREPARE")
	// errFailedDecodePrevote is returned when the PREPARE message is malformed.
	errFailedDecodePrevote = errors.New("failed to decode PREPARE")
	// errFailedDecodePrecommit is returned when the COMMIT message is malformed.
	errFailedDecodePrecommit = errors.New("failed to decode COMMIT")
	// errNoMajority is returned when test value is less than math.Ceil(2/3*N)
	errNoMajority = errors.New("no majority")
	// errNilPrevoteSent is returned when timer could be stopped in time
	errNilPrevoteSent = errors.New("timer expired and nil prevote sent")
)

type Engine interface {
	Start() error
	Stop() error
}

// TODO: add better logging
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

	sentProposal          bool
	setValidRoundAndValue bool

	lockedRound *big.Int
	validRound  *big.Int
	lockedValue *types.Block
	validValue  *types.Block

	currentHeightRoundsStates map[int64]roundState

	// TODO: may require a mutex
	latestPendingUnminedBlock *types.Block

	proposeTimeout   *timeout
	prevoteTimeout   *timeout
	precommitTimeout *timeout

	//TODO: Add futureRounds to keep track of when to move to next round after receiving f + 1 messages from that round
	//map[futureRoundNumber]NumberOfMessagesReceivedForTheRound
	futureRoundsChange map[int64]int64
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
		committedSeals := make([][]byte, c.currentRoundState.Precommits.VotesSize(proposal.ProposalBlock.Hash()))
		for i, v := range c.currentRoundState.Precommits.Values(proposal.ProposalBlock.Hash()) {
			committedSeals[i] = make([]byte, types.PoSExtraSeal)
			copy(committedSeals[i][:], v.CommittedSeal[:])
		}

		if err := c.backend.Precommit(proposal.ProposalBlock, committedSeals); err != nil {
			return
		}
	}
}

// startRound starts a new round. if round equals to 0, it means to starts a new height
func (c *core) startRound(round *big.Int) {
	lastCommittedProposalBlock, lastCommittedProposalBlockProposer := c.backend.LastCommittedProposal()
	height := new(big.Int).Add(lastCommittedProposalBlock.Number(), common.Big1)

	// Start of new height where round is 0
	if round.Uint64() == 0 {
		// Set the shared round values to initial values
		c.lockedRound = big.NewInt(-1)
		c.lockedValue = new(types.Block)
		c.validRound = big.NewInt(-1)
		c.validValue = new(types.Block)

		c.valSet = c.backend.Validators(height.Uint64())

		// Assuming that round == 0 only when the node moves to a new height
		c.currentHeightRoundsStates = make(map[int64]roundState)
	} else {
		// Assuming the above values will be set for round > 0
		// Add the currentRoundState round step to the core previous round states
		// TODO: fix up with references so that I don't overwrite messages when receiving old messages
		c.currentHeightRoundsStates[c.currentRoundState.Round().Int64()] = *c.currentRoundState
	}

	c.proposeTimeout = nil
	c.prevoteTimeout = nil
	c.precommitTimeout = nil

	var rounds []int64
	for k := range c.futureRoundsChange {
		rounds = append(rounds, k)
	}

	for _, r := range rounds {
		if r <= round.Int64() {
			delete(c.futureRoundsChange, r)
		}
	}

	c.currentRoundState = newRoundState(round, height, c.backend.HasBadProposal)
	c.valSet.CalcProposer(lastCommittedProposalBlockProposer, round.Uint64())
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
		c.proposeTimeout.scheduleTimeout(timeoutPropose(round.Int64()), c.currentRoundState.Height().Int64(), c.currentRoundState.Round().Int64(), c.onTimeoutPropose)
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

func (c *core) quorum(i int) bool {
	return float64(i) >= math.Ceil(float64(2)/float64(3)*float64(c.valSet.Size()))
}

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal(hash common.Hash) []byte {
	var buf bytes.Buffer
	buf.Write(hash.Bytes())
	buf.Write([]byte{byte(msgPrecommit)})
	return buf.Bytes()
}
