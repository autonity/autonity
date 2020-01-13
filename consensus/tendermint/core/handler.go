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
	"math/big"
	"sync/atomic"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/crypto"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/clearmatics/autonity/core/types"
)

// Start implements core.Engine.Start
func (c *core) Start(ctx context.Context, chain consensus.ChainReader, currentBlock func() *types.Block, hasBadBlock func(hash common.Hash) bool) error {
	// prevent double start
	if atomic.LoadUint32(c.isStarted) == 1 {
		return nil
	}
	if !atomic.CompareAndSwapUint32(c.isStarting, 0, 1) {
		return nil
	}
	defer func() {
		atomic.StoreUint32(c.isStarting, 0)
		atomic.StoreUint32(c.isStopped, 0)
		atomic.StoreUint32(c.isStarted, 1)
	}()

	ctx, c.cancel = context.WithCancel(ctx)

	err := c.backend.Start(ctx, chain, currentBlock, hasBadBlock)
	if err != nil {
		return err
	}

	c.subscribeEvents()

	//We need a separate go routine to keep c.latestPendingUnminedBlock up to date
	go c.handleNewUnminedBlockEvent(ctx)

	//We want to sequentially handle all the event which modify the current consensus state
	go c.handleConsensusEvents(ctx)

	go c.backend.HandleUnhandledMsgs(ctx)

	return nil
}

// Stop implements core.Engine.Stop
func (c *core) Stop() error {
	// prevent double stop
	if atomic.LoadUint32(c.isStopped) == 1 {
		return nil
	}
	if !atomic.CompareAndSwapUint32(c.isStopping, 0, 1) {
		return nil
	}
	defer func() {
		atomic.StoreUint32(c.isStopping, 0)
		atomic.StoreUint32(c.isStopped, 1)
		atomic.StoreUint32(c.isStarted, 0)
	}()

	c.logger.Info("stopping tendermint.core", "addr", c.address.String())

	_ = c.proposeTimeout.stopTimer()
	_ = c.prevoteTimeout.stopTimer()
	_ = c.precommitTimeout.stopTimer()

	c.cancel()

	c.stopFutureProposalTimer()
	c.unsubscribeEvents()

	<-c.stopped
	<-c.stopped

	err := c.backend.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *core) subscribeEvents() {
	c.messageEventSub = c.backend.Subscribe(events.MessageEvent{})
	c.newUnminedBlockEventSub = c.backend.Subscribe(events.NewUnminedBlockEvent{})
	c.timeoutEventSub = c.backend.Subscribe(TimeoutEvent{})
	c.committedSub = c.backend.Subscribe(events.CommitEvent{})
	c.syncEventSub = c.backend.Subscribe(events.SyncEvent{})
}

// Unsubscribe all messageEventSub
func (c *core) unsubscribeEvents() {
	c.messageEventSub.Unsubscribe()
	c.newUnminedBlockEventSub.Unsubscribe()
	c.timeoutEventSub.Unsubscribe()
	c.committedSub.Unsubscribe()
	c.syncEventSub.Unsubscribe()
}

func (c *core) handleNewUnminedBlockEvent(ctx context.Context) {
eventLoop:
	for {
		select {
		case e, ok := <-c.newUnminedBlockEventSub.Chan():
			if !ok {
				break eventLoop
			}
			newUnminedBlockEvent := e.Data.(events.NewUnminedBlockEvent)
			pb := &newUnminedBlockEvent.NewUnminedBlock
			c.storeUnminedBlockMsg(pb)
		case <-ctx.Done():
			c.logger.Info("handleNewUnminedBlockEvent is stopped", "event", ctx.Err())
			break eventLoop
		}
	}

	c.stopped <- struct{}{}
}

func (c *core) handleConsensusEvents(ctx context.Context) {
	// Start a new round from last height + 1
	c.startRound(ctx, common.Big0)

	go c.syncLoop(ctx)

eventLoop:
	for {
		select {
		case ev, ok := <-c.messageEventSub.Chan():
			if !ok {
				break eventLoop
			}
			// A real ev arrived, process interesting content
			if messageE, ok := ev.Data.(events.MessageEvent); ok {
				if len(messageE.Payload) == 0 {
					c.logger.Error("core.handleConsensusEvents Get message(MessageEvent) empty payload")
				}

				if err := c.handleMsg(ctx, messageE.Payload); err != nil {
					c.logger.Debug("core.handleConsensusEvents Get message(MessageEvent) payload failed", "err", err)
					continue
				}
				c.backend.Gossip(ctx, c.valSet.Copy(), messageE.Payload)
			}
		case ev, ok := <-c.timeoutEventSub.Chan():
			if !ok {
				break eventLoop
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
				break eventLoop
			}
			if _, ok := ev.Data.(events.CommitEvent); ok {
				c.handleCommit(ctx)
			}
		case <-ctx.Done():
			c.logger.Info("handleConsensusEvents is stopped", "event", ctx.Err())
			break eventLoop
		}
	}

	c.stopped <- struct{}{}
}

func (c *core) syncLoop(ctx context.Context) {
	/*
		this method is responsible for asking the network to send us the current consensus state
		and to process sync queries events.
	*/
	timer := time.NewTimer(10 * time.Second)

	round := c.roundState.Round()
	height := c.roundState.Height()

	// Ask for sync when the engine starts
	c.backend.AskSync(c.valSet.Copy())

	for {
		select {
		case <-timer.C:
			currentRound := c.roundState.Round()
			currentHeight := c.roundState.Height()

			// we only ask for sync if the current view stayed the same for the past 10 seconds
			if currentHeight.Cmp(height) == 0 && currentRound.Cmp(round) == 0 {
				c.backend.AskSync(c.valSet.Copy())
			}
			round = currentRound
			height = currentHeight
			timer = time.NewTimer(10 * time.Second)
		case ev, ok := <-c.syncEventSub.Chan():
			if !ok {
				return
			}
			event := ev.Data.(events.SyncEvent)
			c.logger.Info("Processing sync message", "from", event.Addr)
			c.SyncPeer(event.Addr)
		case <-ctx.Done():
			return
		}
	}
}

// sendEvent sends event to mux
func (c *core) sendEvent(ev interface{}) {
	c.backend.Post(ev)
}

func (c *core) handleMsg(ctx context.Context, payload []byte) error {
	logger := c.logger.New()
	// Decode message and check its signature
	msg := new(Message)
	sender, err := msg.FromPayload(payload, c.valSet.Copy(), crypto.CheckValidatorSignature)
	if err != nil {
		logger.Error("Failed to decode message from payload", "err", err)
		return err
	}

	return c.handleCheckedMsg(ctx, msg, *sender)
}

func (c *core) handleCheckedMsg(ctx context.Context, msg *Message, sender validator.Validator) error {
	logger := c.logger.New("address", c.address, "from", sender)

	switch msg.Code {
	case msgProposal:
		logger.Debug("tendermint.MessageEvent: PROPOSAL")
		return c.handleProposal(ctx, msg)
	case msgPrevote:
		logger.Debug("tendermint.MessageEvent: PREVOTE")
		return c.handlePrevote(ctx, msg)
	case msgPrecommit:
		logger.Debug("tendermint.MessageEvent: PRECOMMIT")
		return c.handlePrecommit(ctx, msg)
	default:
		logger.Error("Invalid message", "msg", msg)
	}

	return errInvalidMessage
}

// checkMessage checks the message step
// return errInvalidMessage if the message is invalid
// return errFutureHeightMessage if the message view is larger than roundState view
// return errOldHeightMessage if the message view is smaller than roundState view
// return errFutureStepMessage if we are at the same view but at the propose step and it's a voting message.
func (c *core) checkMessage(round *big.Int, height *big.Int, step Step) error {
	if height == nil || round == nil {
		return errInvalidMessage
	}

	if height.Cmp(c.roundState.Height()) > 0 {
		return errFutureHeightMessage
	} else if height.Cmp(c.roundState.Height()) < 0 {
		return errOldHeightMessage
	}

	return nil
}
