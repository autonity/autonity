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
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/contracts/autonity"
	"github.com/clearmatics/autonity/core/types"
	"github.com/davecgh/go-spew/spew"
)

var errStopped error = errors.New("stopped")

// Start implements core.Tendermint.Start
func (c *core) Start(ctx context.Context, contract *autonity.Contract) {
	println("starting")
	// Set the autonity contract
	c.autonityContract = contract
	ctx, c.cancel = context.WithCancel(ctx)

	// Subscribe
	c.eventsSub = c.backend.Subscribe(events.MessageEvent{}, &algorithm.Timeout{}, events.CommitEvent{})
	c.syncEventSub = c.backend.Subscribe(events.SyncEvent{})
	c.newUnminedBlockEventSub = c.backend.Subscribe(events.NewUnminedBlockEvent{})

	c.wg = &sync.WaitGroup{}
	// We need a separate go routine to keep c.latestPendingUnminedBlock up to date
	c.wg.Add(1)
	go c.handleNewUnminedBlockEvent(ctx)

	// Tendermint Finite State Machine discrete event loop
	c.wg.Add(1)
	go c.mainEventLoop(ctx)

	go c.backend.HandleUnhandledMsgs(ctx)
}

// stop implements core.Engine.stop
func (c *core) Stop() {
	println(addr(c.address), c.height, "stopping")

	c.logger.Info("stopping tendermint.core", "addr", addr(c.address))

	c.cancel()

	// Signal to wake up await value if it is waiting.
	c.valueSet.L.Lock()
	c.valueSet.Signal()
	c.valueSet.L.Unlock()

	// Unsubscribe
	c.eventsSub.Unsubscribe()
	c.syncEventSub.Unsubscribe()

	println(addr(c.address), c.height, "almost stopped")
	// Ensure all event handling go routines exit
	c.wg.Wait()
}

func (c *core) handleNewUnminedBlockEvent(ctx context.Context) {
	defer c.wg.Done()
eventLoop:
	for {
		select {
		case e, ok := <-c.newUnminedBlockEventSub.Chan():
			if !ok {
				break eventLoop
			}
			block := e.Data.(events.NewUnminedBlockEvent).NewUnminedBlock
			c.SetValue(&block)
		case <-ctx.Done():
			c.logger.Info("handleNewUnminedBlockEvent is stopped", "event", ctx.Err())
			break eventLoop
		}
	}
}

func (c *core) newHeight(ctx context.Context, height uint64) error {
	c.syncTimer = time.NewTimer(20 * time.Second)
	newHeight := new(big.Int).SetUint64(height)
	// set the new height
	c.height = newHeight
	var err error
	c.currentBlock, err = c.AwaitValue(ctx, newHeight)
	if err != nil {
		return err
	}
	prevBlock, _ := c.backend.LastCommittedProposal()

	c.lastHeader = prevBlock.Header()
	committeeSet := c.createCommittee(prevBlock)
	c.committee = committeeSet

	// Update internals of oracle
	c.ora.lastHeader = c.lastHeader
	c.ora.committeeSet = committeeSet

	// Handle messages for the new height
	r := c.algo.StartRound(newHeight.Uint64(), 0, algorithm.ValueID(c.currentBlock.Hash()))

	// If we are making a proposal, we need to ensure that we add the proposal
	// block to the msg store, so that it can be picked up in buildMessage.
	if r.Broadcast != nil {
		println(addr(c.address), "adding value", height, c.currentBlock.Hash().String())
		c.msgCache.addValue(c.currentBlock.Hash(), c.currentBlock)
	}

	// Note that we don't risk enterning an infinite loop here since
	// start round can only return results with brodcasts or schedules.
	// TODO actually don't return result from Start round.
	err = c.handleResult(ctx, r)
	if err != nil {
		return err
	}
	for _, msg := range c.msgCache.heightMessages(newHeight.Uint64()) {
		go func(m *message) {
			err := c.handleCurrentHeightMessage(ctx, m)
			c.logger.Error("failed to handle current height message", "message", m.String, "err", err)
		}(msg)
	}
	return nil
}

func (c *core) handleResult(ctx context.Context, r *algorithm.Result) error {

	switch {
	case r == nil:
		return nil
	case r.StartRound != nil:
		sr := r.StartRound
		if sr.Round == 0 && sr.Decision == nil {
			panic("round changes of 0 must be accompanied with a decision")
		}
		if sr.Decision != nil {
			// A decision has been reached
			println(addr(c.address), "decided on block", sr.Decision.Height,
				common.Hash(sr.Decision.Value).String())

			// This will ultimately lead to a commit event, which we will pick
			// up on but we will ignore it because instead we will wait here to
			// select the next value that matches this height.
			_, err := c.Commit(sr.Decision)
			if err != nil {
				panic(fmt.Sprintf("%s Failed to commit sr.Decision: %s err: %v", algorithm.NodeID(c.address).String(), spew.Sdump(sr.Decision), err))
			}
			err = c.newHeight(ctx, sr.Height)
			if err != nil {
				return err
			}

		} else {
			// sanity check
			currBlockNum := c.currentBlock.Number().Uint64()
			if currBlockNum != sr.Height {
				panic(fmt.Sprintf("current block number %d out of sync with  height %d", currBlockNum, sr.Height))
			}

			rr := c.algo.StartRound(sr.Height, sr.Round, algorithm.ValueID(c.currentBlock.Hash()))
			// Note that we don't risk enterning an infinite loop here since
			// start round can only return results with brodcasts or schedules.
			// TODO actually don't return result from Start round.
			err := c.handleResult(ctx, rr)
			if err != nil {
				return err
			}
		}
	case r.Broadcast != nil:
		println(addr(c.address), c.height.String(), r.Broadcast.String(), "sending")
		// Broadcasting ends with the message reaching us eventually

		// We must build message here since buildMessage relies on accessing
		// the msg store, and since the message store is not syncronised we
		// need to do it from the handler routine.
		msg, err := encodeSignedMessage(r.Broadcast, c.key, c.msgCache)
		if err != nil {
			panic(fmt.Sprintf(
				"%s We were unable to build a message, this indicates a programming error: %v",
				addr(c.address),
				err,
			))
		}

		// Broadcast in a new goroutine
		go func(committee types.Committee) {
			err := c.backend.Broadcast(ctx, committee, msg)
			if err != nil {
				c.logger.Error("Failed to broadcast message", "msg", msg, "err", err)
			}
		}(c.lastHeader.Committee)

	case r.Schedule != nil:
		time.AfterFunc(time.Duration(r.Schedule.Delay)*time.Second, func() {
			c.backend.Post(r.Schedule)
		})

	}
	return nil
}

func (c *core) mainEventLoop(ctx context.Context) {
	defer c.wg.Done()
	// Start a new round from last height + 1
	c.algo = algorithm.New(algorithm.NodeID(c.address), c.ora)
	lastBlockMined, _ := c.backend.LastCommittedProposal()
	err := c.newHeight(ctx, lastBlockMined.NumberU64()+1)
	if err != nil {
		println(addr(c.address), c.height.Uint64(), "exiting main event loop", "err", err)
		return
	}

	// Ask for sync when the engine starts
	c.backend.AskSync(c.lastHeader)

eventLoop:
	for {
		select {
		case <-c.syncTimer.C:
			c.backend.AskSync(c.lastHeader)
			c.syncTimer = time.NewTimer(20 * time.Second)

		case ev, ok := <-c.syncEventSub.Chan():
			if !ok {
				break eventLoop
			}
			event := ev.Data.(events.SyncEvent)
			c.logger.Info("Processing sync message", "from", event.Addr)
			c.backend.SyncPeer(event.Addr, c.msgCache.rawHeightMessages(c.height.Uint64()))
		case ev, ok := <-c.eventsSub.Chan():
			if !ok {
				break eventLoop
			}
			// A real ev arrived, process interesting content
			switch e := ev.Data.(type) {
			case events.MessageEvent:
				err := c.handleMsg(ctx, e.Payload)
				if err == errStopped {
					return
				}
				if err != nil {
					c.logger.Debug("core.mainEventLoop problem processing message", "err", err)
					continue
				}
				c.backend.Gossip(ctx, c.lastHeader.Committee, e.Payload)
			case *algorithm.ConsensusMessage:
				println(addr(c.address), e.String(), "message from self")
				// This is a message we sent ourselves we do not need to broadcast it
				if c.height.Uint64() == e.Height {
					r := c.algo.ReceiveMessage(e)
					err := c.handleResult(ctx, r)
					if err != nil {
						println(addr(c.address), c.height.Uint64(), "exiting main event loop", "err", err)
						return
					}
				}
			case *algorithm.Timeout:
				var r *algorithm.Result
				switch e.TimeoutType {
				case algorithm.Propose:
					println(addr(c.address), "on timeout propose", e.Height, "round", e.Round)
					r = c.algo.OnTimeoutPropose(e.Height, e.Round)
				case algorithm.Prevote:
					println(addr(c.address), "on timeout prevote", e.Height, "round", e.Round)
					r = c.algo.OnTimeoutPrevote(e.Height, e.Round)
				case algorithm.Precommit:
					println(addr(c.address), "on timeout precommit", e.Height, "round", e.Round)
					r = c.algo.OnTimeoutPrecommit(e.Height, e.Round)
				}
				if r != nil && r.Broadcast != nil {
					println("nonnil timeout")
				}
				err := c.handleResult(ctx, r)
				if err != nil {
					println(addr(c.address), c.height.Uint64(), "exiting main event loop", "err", err)
					return
				}
			case events.CommitEvent:
				println(addr(c.address), "commit event")
				c.logger.Debug("Received a final committed proposal")
				lastBlock, _ := c.backend.LastCommittedProposal()
				height := new(big.Int).Add(lastBlock.Number(), common.Big1)
				if height.Cmp(c.height) == 0 {
					println(addr(c.address), "Discarding event as core is at the same height", "height", c.height)
					c.logger.Debug("Discarding event as core is at the same height", "height", c.height)
				} else {
					println(addr(c.address), "Received proposal is ahead", "height", c.height, "block_height", height.String())
					c.logger.Debug("Received proposal is ahead", "height", c.height, "block_height", height)
					err := c.newHeight(ctx, height.Uint64())
					if err != nil {
						println(addr(c.address), c.height.Uint64(), "exiting main event loop", "err", err)
						return
					}
				}
			}
		case <-ctx.Done():
			c.logger.Info("mainEventLoop is stopped", "event", ctx.Err())
			break eventLoop
		}
	}

}

func (c *core) handleMsg(ctx context.Context, msgBytes []byte) error {

	println("got a message")
	/*
		Basic validity checks
	*/

	m, err := decodeSignedMessage(msgBytes)
	if err != nil {
		fmt.Printf("some error: %v\n", err)
		return err
	}
	// Check we haven't already processed this message
	if c.msgCache.Message(m.hash) != nil {
		// Message was already processed
		return nil
	}
	err = c.msgCache.addMessage(m, msgBytes)
	if err != nil {
		// could be multiple proposal messages from the same proposer
		return err
	}
	if m.consensusMessage.MsgType == algorithm.Propose {
		c.msgCache.addValue(m.value.Hash(), m.value)
	}

	// If this message is for a future height then we cannot validate it
	// because we lack the relevant header, we will process it when we reach
	// that height. If it is for a previous height then we are not intersted in
	// it. But it has been added to the msg cache in case other peers would
	// like to sync it.
	if m.consensusMessage.Height != c.height.Uint64() {
		// Nothing to do here
		return nil
	}

	return c.handleCurrentHeightMessage(ctx, m)

}

func (c *core) handleCurrentHeightMessage(ctx context.Context, m *message) error {
	println(addr(c.address), c.height.String(), m.String(), "received")
	cm := m.consensusMessage
	/*
		Domain specific validity checks, now we know that we are at the same
		height as this message we can rely on lastHeader.
	*/

	// Check that the message came from a committee member, if not we ignore it.
	if c.lastHeader.CommitteeMember(m.address) == nil {
		// TODO turn this into an error type that can be checked for at a
		// higher level to close the connection to this peer.
		return fmt.Errorf("received message from non committee member: %v", m)
	}

	switch cm.MsgType {
	case algorithm.Propose:
		// We ignore proposals from non proposers
		if c.committee.GetProposer(cm.Round).Address != m.address {
			c.logger.Warn("Ignore proposal messages from non-proposer")
			return errNotFromProposer

			// TODO verify proposal here.
			//
			// If we are introducing time into the mix then what we are saying
			// is that we don't expect different participants' clocks to drift
			// out of sync more than some delta. And if they do then we don't
			// expect consensus to work.
			//
			// So in the case that clocks drift too far out of sync and say a
			// node considers a proposal invalid that 2f+1 other nodes
			// precommit for that node becomes stuck and can only continue in
			// consensus by re-syncing the blocks.
			//
			// So in verifying the proposal wrt time we should verify once
			// within reasonable clock sync bounds and then set the validity
			// based on that and never re-process the message again.

		}
		// Proposals values are allowed to be invalid.
		if _, err := c.backend.VerifyProposal(*c.msgCache.value(common.Hash(cm.Value))); err == nil {
			println(addr(c.address), "valid", cm.Value.String())
			c.msgCache.setValid(common.Hash(cm.Value))
		}
	default:
		// All other messages that have reached this point are valid, but we are not marking the vlaue valid here, we are marking the message valid.
		c.msgCache.setValid(m.hash)
	}

	r := c.algo.ReceiveMessage(cm)
	err := c.handleResult(ctx, r)
	if err != nil {
		return err
	}
	return nil
}
