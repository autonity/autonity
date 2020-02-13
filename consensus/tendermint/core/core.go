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
	"github.com/clearmatics/autonity/consensus/tendermint/events"
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
	// errFutureHeightMessage is returned when round is earlier than the
	// view of the received message.
	errFutureHeightMessage = errors.New("future height message")
	// errOldHeightMessage is returned when the received message's round is earlier
	// than the validator's round.
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
		allProposals:               make(map[int64]*proposalSet),
		allPrevotes:                make(map[int64]*messageSet),
		allPrecommits:              make(map[int64]*messageSet),
		verifiedProposals:          make(map[common.Hash]bool),
		futureHeightMessages:       make(map[int64][]*Message),
		lockedRound:                big.NewInt(-1),
		validRound:                 big.NewInt(-1),
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
	startRoundEventSub      *event.TypeMuxSubscription
	syncEventSub            *event.TypeMuxSubscription
	futureProposalTimer     *time.Timer
	stopped                 chan struct{}
	isStarted               *uint32
	isStarting              *uint32
	isStopping              *uint32
	isStopped               *uint32

	valSet *validatorSet

	round                *big.Int
	height               *big.Int
	step                 Step
	allProposals         map[int64]*proposalSet
	allPrevotes          map[int64]*messageSet
	allPrecommits        map[int64]*messageSet
	verifiedProposals    map[common.Hash]bool
	futureHeightMessages map[int64][]*Message
	coreMu               sync.RWMutex

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

// StartRoundEvent is posted for Tendermint round change
type StartRoundEvent struct {
	round int64
}

// startRound starts a new round. if round equals to 0, it means to starts a new height
func (c *core) startRound(ctx context.Context, round *big.Int) {
	c.measureHeightRoundMetrics(round)

	proposalBlock, blockProposer := c.backend.LastCommittedProposal()

	c.lastCommittedBlockProposer = blockProposer
	height := new(big.Int).Add(proposalBlock.Number(), common.Big1)
	c.setHeight(height)
	c.setRound(round)
	_ = c.setStep(ctx, propose)

	if round.Int64() == 0 {
		c.lockedRound = big.NewInt(-1)
		c.lockedValue = nil
		c.validRound = big.NewInt(-1)
		c.validValue = nil

		// reset all maps
		c.coreMu.Lock()
		c.allProposals = make(map[int64]*proposalSet)
		c.allPrevotes = make(map[int64]*messageSet)
		c.allPrecommits = make(map[int64]*messageSet)
		c.verifiedProposals = make(map[common.Hash]bool)
		c.coreMu.Unlock()

		// Set validator set for height
		valSet := c.backend.Validators(height.Uint64())
		c.valSet.set(valSet)

		//Send all messages stored in futureHeightMessages to handler
		if ms, ok := c.futureHeightMessages[height.Int64()]; ok {
			for _, m := range ms {
				p, _ := m.Payload()
				go c.sendEvent(events.MessageEvent{Payload: p})
			}
			// Once finished sending messages back to handler delete key value pair
			delete(c.futureHeightMessages, height.Int64())
		}
	}
	// Reset all timeouts
	c.proposeTimeout.reset(propose)
	c.prevoteTimeout.reset(prevote)
	c.precommitTimeout.reset(precommit)

	// Calculate new proposer
	c.valSet.CalcProposer(blockProposer, round.Uint64())

	c.sentProposal = false
	c.sentPrevote = false
	c.sentPrecommit = false
	c.setValidRoundAndValue = false

	c.logger.Debug("Starting new Round", "Height", height, "Round", round)

	// If the node is the proposer for this round then it would propose validValue or a new block, otherwise,
	// proposeTimeout is started, where the node waits for a proposal from the proposer of the current round.
	if c.isProposerForR(round.Int64(), c.address) {
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
		// round, if so then check for propose step condition. It is simpler to check for the conditions here because
		// step is set to propose only once and it is done above. Also, if we are the validator for this round and are
		// honest then we will not have the proposal, therefore checking for the propose step condition in setStep()
		// will require determining whether we are the proposer or not adding to more complexity.
		if proposalMS := c.getProposalSet(round.Int64()); proposalMS != nil {
			proposal := proposalMS.proposal()
			if proposal.ValidRound.Int64() == -1 {
				if err := c.checkForNewProposal(ctx, round.Int64()); err != nil {
					c.logger.Error(err.Error())
				}
			} else if proposal.ValidRound.Int64() >= 0 {
				if err := c.checkForOldProposal(ctx, round.Int64()); err != nil {
					c.logger.Error(err.Error())
				}
			}
		}
	}
}

func (c *core) isValidator(address common.Address) bool {
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
	logger := c.logger.New("step", c.getStep())

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

func (c *core) isProposerForR(r int64, a common.Address) bool {
	return c.valSet.IsProposerForRound(c.lastCommittedBlockProposer, uint64(r), a)
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

func (c *core) stopFutureProposalTimer() {
	if c.futureProposalTimer != nil {
		c.futureProposalTimer.Stop()
	}
}

func (c *core) quorum(i int) bool {
	return float64(i) >= math.Ceil(float64(2)/float64(3)*float64(c.valSet.Size()))
}

func (c *core) addFutureHeighMessage(height int64, msg *Message) {
	c.futureHeightMessages[height] = append(c.futureHeightMessages[height], msg)
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
	if _, ok := c.verifiedProposals[proposal.Hash()]; !ok {
		if _, err := c.backend.VerifyProposal(*proposal); err != nil {
			return false, err
		}
		c.verifiedProposals[proposal.Hash()] = true
	}
	return true, nil
}

func (c *core) setRound(r *big.Int) {
	c.coreMu.Lock()
	defer c.coreMu.Unlock()

	c.round = big.NewInt(r.Int64())
}

func (c *core) getRound() *big.Int {
	c.coreMu.RLock()
	defer c.coreMu.RUnlock()

	return c.round
}

func (c *core) setHeight(height *big.Int) {
	c.coreMu.Lock()
	defer c.coreMu.Unlock()

	c.height = height
}

func (c *core) getHeight() *big.Int {
	c.coreMu.RLock()
	defer c.coreMu.RUnlock()

	return c.height
}

func (c *core) setStep(ctx context.Context, step Step) error {
	c.coreMu.Lock()
	c.step = step
	c.coreMu.Unlock()

	// We need to check for upon conditions which refer to a specific step, so that once a validator moves to that step
	// and no more messages are received, we ensure that if an upon condition is true it is executed. Propose step upon
	// are checked when round is changed in startRound() by resending a future proposal if it was received in an older
	// round to ensure line 22  and line 28 are run. When the validator moves to the prevote step we need to check the
	// prevote step upon conditions here. For precommit step upon condition nothing needs to be done even though line 36
	// predicates on it, this is because line 36 will be run when validator moves to prevote step, prevotes arrival and
	// and proposal arrival in prevote/precommit step. Since quorum prevotes is required to move to precommit step, line
	// 36 would have been executed in prevote step if not because some prevotes were nil, then once a valid prevote
	// arrive then line 36 will be run in precommit step. Otherwise the condition is not satisfied to run line 36.
	// Therefore, nothing needs to be done when a validator moves to the precommit step.

	if step == prevote {
		// Check for line 34, 36 and 44
		curR := c.getRound().Int64()
		curH := c.getHeight().Int64()
		c.checkForPrevoteTimeout(curR, curH)
		if err := c.checkForQuorumPrevotes(ctx, curR); err != nil {
			return err
		}
		if err := c.checkForQuorumPrevotesNil(ctx, curR); err != nil {
			return err
		}
	}

	return nil
}

func (c *core) getStep() Step {
	c.coreMu.RLock()
	defer c.coreMu.RUnlock()
	return c.step
}

func (c *core) currentState() (*big.Int, *big.Int, uint64) {
	return c.getHeight(), c.getRound(), uint64(c.getStep())
}

func (c *core) getAllRoundMessages() []*Message {
	c.coreMu.RLock()
	defer c.coreMu.RUnlock()
	var messages []*Message

	for _, proposalMS := range c.allProposals {
		messages = append(messages, proposalMS.proposalMsg())
	}
	lenProposal := len(messages)
	c.logger.Debug("Collecting messages for sync", "#Proposals", lenProposal)

	for _, prevoteMS := range c.allPrevotes {
		messages = append(messages, prevoteMS.GetMessages()...)
	}
	lenPrevotes := len(messages) - lenProposal
	c.logger.Debug("Collecting messages for sync", "#Prevotes", lenPrevotes)

	for _, precommitMS := range c.allPrecommits {
		messages = append(messages, precommitMS.GetMessages()...)
	}
	lenPrecommits := len(messages) - lenPrevotes - lenProposal
	c.logger.Debug("Collecting messages for sync", "#Precommits", lenPrecommits)

	return messages
}

// Determine if we already have vote from the sender
func (c *core) hasVote(v Vote, m *Message) bool {
	var votes messageSet
	voteRound := v.Round.Int64()
	mCode := m.Code

	if mCode == msgPrevote {
		prevotesSet := c.getPrevotesSet(voteRound)
		if prevotesSet == nil {
			return false
		}
		votes = *prevotesSet
	} else if mCode == msgPrecommit {
		precommitsSet := c.getPrecommitsSet(voteRound)
		if precommitsSet == nil {
			return false
		}
		votes = *precommitsSet
	}
	return votes.hasMessage(*m)
}

func (c *core) getProposalSet(round int64) *proposalSet {
	c.coreMu.RLock()
	defer c.coreMu.RUnlock()

	proposalS, ok := c.allProposals[round]
	if !ok {
		return nil
	}

	return proposalS
}

func (c *core) setProposalSet(round int64, p Proposal, m *Message) {
	c.coreMu.Lock()
	defer c.coreMu.Unlock()
	c.allProposals[round] = newProposalSet(p, m)
}

func (c *core) getPrevotesSet(round int64) *messageSet {
	c.coreMu.RLock()
	defer c.coreMu.RUnlock()

	prevotesS, ok := c.allPrevotes[round]
	if !ok {
		return nil
	}

	return prevotesS
}

func (c *core) setPrevotesSet(round int64) {
	c.coreMu.Lock()
	defer c.coreMu.Unlock()
	c.allPrevotes[round] = newMessageSet()
}

func (c *core) getPrecommitsSet(round int64) *messageSet {
	c.coreMu.RLock()
	defer c.coreMu.RUnlock()

	precommitsS, ok := c.allPrecommits[round]
	if !ok {
		return nil
	}

	return precommitsS
}

func (c *core) setPrecommitsSet(round int64) {
	c.coreMu.Lock()
	defer c.coreMu.Unlock()
	c.allPrecommits[round] = newMessageSet()
}
