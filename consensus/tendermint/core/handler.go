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
	"context"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/crypto"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/clearmatics/autonity/core/types"
	"math/big"
)

// Start implements core.Engine.Start
func (c *core) Start(chain consensus.ChainReader, currentBlock func() *types.Block, hasBadBlock func(hash common.Hash) bool) error {
	err := c.backend.Start(chain, currentBlock, hasBadBlock)
	if err != nil {
		return err
	}

	c.subscribeEvents()

	// set currentRoundState before starting go routines
	lastCommittedProposalBlock, _ := c.backend.LastCommittedProposal()
	height := new(big.Int).Add(lastCommittedProposalBlock.Number(), common.Big1)
	c.currentRoundState = NewRoundState(big.NewInt(0), height)

	//We need a separate go routine to keep c.latestPendingUnminedBlock up to date
	go c.handleNewUnminedBlockEvent()

	var ctx context.Context
	ctx, c.cancel = context.WithCancel(context.Background())

	//We want to sequentially handle all the event which modify the current consensus state
	go c.handleConsensusEvents(ctx)

	return nil
}

// Stop implements core.Engine.Stop
func (c *core) Stop() error {
	err := c.backend.Close()
	if err != nil {
		return err
	}

	c.stopFutureProposalTimer()
	c.unsubscribeEvents()

	_ = c.proposeTimeout.stopTimer()
	_ = c.prevoteTimeout.stopTimer()
	_ = c.precommitTimeout.stopTimer()

	// Make sure the handler goroutine exits
	c.cancel()
	return nil
}

func (c *core) subscribeEvents() {
	c.messageEventSub = c.backend.EventMux().Subscribe(
		// external messages
		events.MessageEvent{},
		// internal messages
		backlogEvent{},
	)
	c.newUnminedBlockEventSub = c.backend.EventMux().Subscribe(
		events.NewUnminedBlockEvent{},
	)
	c.timeoutEventSub = c.backend.EventMux().Subscribe(
		TimeoutEvent{},
	)
	c.committedSub = c.backend.EventMux().Subscribe(
		events.CommitEvent{},
	)
}

// Unsubscribe all messageEventSub
func (c *core) unsubscribeEvents() {
	c.messageEventSub.Unsubscribe()
	c.newUnminedBlockEventSub.Unsubscribe()
	c.timeoutEventSub.Unsubscribe()
	c.committedSub.Unsubscribe()
}

// TODO: update all of the TypeMuxSilent to event.Feed and should not use backend.EventMux for core internal messageEventSub: backlogEvent, TimeoutEvent

func (c *core) handleNewUnminedBlockEvent() {
	for e := range c.newUnminedBlockEventSub.Chan() {
		newUnminedBlockEvent := e.Data.(events.NewUnminedBlockEvent)
		pb := &newUnminedBlockEvent.NewUnminedBlock
		c.storeUnminedBlockMsg(pb)
	}
}

func (c *core) handleConsensusEvents(ctx context.Context) {
	// Start a new round from last height + 1
	c.startRound(ctx, common.Big0)

	for {
		select {
		case ev, ok := <-c.messageEventSub.Chan():
			if !ok {
				return
			}
			// A real ev arrived, process interesting content
			switch e := ev.Data.(type) {
			case events.MessageEvent:
				if len(e.Payload) == 0 {
					c.logger.Error("core.handleConsensusEvents Get message(MessageEvent) empty payload")
				}

				if err := c.handleMsg(ctx, e.Payload); err != nil {
					c.logger.Debug("core.handleConsensusEvents Get message(MessageEvent) payload failed", "err", err)
					continue
				}
				c.backend.Gossip(ctx, c.valSet.Copy(), e.Payload)
			case backlogEvent:
				// No need to check signature for internal messages
				c.logger.Debug("Started handling backlogEvent")
				err := c.handleCheckedMsg(ctx, e.msg, e.src)
				if err != nil {
					c.logger.Debug("core.handleConsensusEvents handleCheckedMsg message failed", "err", err)
					continue
				}

				p, err := e.msg.Payload()
				if err != nil {
					c.logger.Debug("core.handleConsensusEvents Get message payload failed", "err", err)
					continue
				}

				c.backend.Gossip(ctx, c.valSet.Copy(), p)
			}
		case ev, ok := <-c.timeoutEventSub.Chan():
			if !ok {
				return
			}
			if timeoutE, ok := ev.Data.(TimeoutEvent); ok {
				switch timeoutE.step {
				case msgProposal:
					c.handleTimeoutPropose(ctx, timeoutE)
				case msgPrevote:
					c.handleTimeoutPrevote(ctx, timeoutE)
				case msgPrecommit:
					c.handleTimeoutPrecommit(ctx, timeoutE)
				}
			}
		case ev, ok := <-c.committedSub.Chan():
			if !ok {
				return
			}
			switch ev.Data.(type) {
			case events.CommitEvent:
				c.handleCommit(ctx)
			}
		}
	}
}

// sendEvent sends event to mux
func (c *core) sendEvent(ev interface{}) {
	c.backend.EventMux().Post(ev)
}

func (c *core) handleMsg(ctx context.Context, payload []byte) error {
	logger := c.logger.New()

	// Decode message and check its signature
	msg := new(message)

	sender, err := msg.FromPayload(payload, c.valSet.Copy(), crypto.CheckValidatorSignature)
	if err != nil {
		logger.Error("Failed to decode message from payload", "err", err)
		return err
	}

	return c.handleCheckedMsg(ctx, msg, *sender)
}

func (c *core) handleCheckedMsg(ctx context.Context, msg *message, sender validator.Validator) error {
	logger := c.logger.New("address", c.address, "from", sender)

	// Store the message if it's a future message
	testBacklog := func(err error) error {
		// We want to store only future messages in backlog
		if err == errFutureHeightMessage {
			logger.Debug("Storing future height message in backlog")
			c.storeBacklog(msg, sender)
		} else if err == errFutureRoundMessage {
			logger.Debug("Storing future height message in backlog")
			c.storeBacklog(msg, sender)
			//We cannot move to a round in a new height without receiving a new block
			var msgRound int64
			if msg.Code == msgProposal {
				var p Proposal
				if e := msg.Decode(&p); e != nil {
					return errFailedDecodeProposal
				}
				msgRound = p.Round.Int64()

			} else {
				var v Vote
				if e := msg.Decode(&v); e != nil {
					return errFailedDecodeVote
				}
				msgRound = v.Round.Int64()
			}

			c.futureRoundsChange[msgRound] = c.futureRoundsChange[msgRound] + 1
			totalFutureRoundMessages := c.futureRoundsChange[msgRound]

			if totalFutureRoundMessages > int64(c.valSet.F()) {
				logger.Debug("Received ceil(N/3) - 1 messages for higher round", "New round", msgRound)
				c.startRound(ctx, big.NewInt(msgRound))
			}
		}

		return err
	}

	switch msg.Code {
	case msgProposal:
		logger.Debug("tendermint.MessageEvent: PROPOSAL")
		return testBacklog(c.handleProposal(ctx, msg))
	case msgPrevote:
		logger.Debug("tendermint.MessageEvent: PREVOTE")
		return testBacklog(c.handlePrevote(ctx, msg))
	case msgPrecommit:
		logger.Debug("tendermint.MessageEvent: PRECOMMIT")
		return testBacklog(c.handlePrecommit(ctx, msg))
	default:
		logger.Error("Invalid message", "msg", msg)
	}

	return errInvalidMessage
}
