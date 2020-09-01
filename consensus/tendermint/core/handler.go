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
	"fmt"
	"math/big"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/crypto"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/contracts/autonity"
	"github.com/clearmatics/autonity/core/types"
	autonitycrypto "github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/rlp"
)

// Start implements core.Tendermint.Start
func (c *core) Start(ctx context.Context, contract *autonity.Contract) {
	// Set the autonity contract
	c.autonityContract = contract
	ctx, c.cancel = context.WithCancel(ctx)

	c.subscribeEvents()

	// core.height needs to be set beforehand for unmined block's logic.
	lastBlockMined, _ := c.backend.LastCommittedProposal()
	c.setHeight(new(big.Int).Add(lastBlockMined.Number(), common.Big1))
	// We need a separate go routine to keep c.latestPendingUnminedBlock up to date
	go c.handleNewUnminedBlockEvent(ctx)

	// Tendermint Finite State Machine discrete event loop
	go c.mainEventLoop(ctx)

	go c.backend.HandleUnhandledMsgs(ctx)
}

// Stop implements core.Engine.Stop
func (c *core) Stop() {
	c.logger.Info("stopping tendermint.core", "addr", c.address.String())

	_ = c.proposeTimeout.stopTimer()
	_ = c.prevoteTimeout.stopTimer()
	_ = c.precommitTimeout.stopTimer()

	c.cancel()

	c.stopFutureProposalTimer()
	c.unsubscribeEvents()

	// Ensure all event handling go routines exit
	<-c.stopped
	<-c.stopped
	<-c.stopped
}

func (c *core) subscribeEvents() {
	s := c.backend.Subscribe(events.MessageEvent{}, backlogEvent{})
	c.messageEventSub = s

	s1 := c.backend.Subscribe(events.NewUnminedBlockEvent{})
	c.newUnminedBlockEventSub = s1

	s2 := c.backend.Subscribe(TimeoutEvent{})
	c.timeoutEventSub = s2

	s3 := c.backend.Subscribe(events.CommitEvent{})
	c.committedSub = s3

	s4 := c.backend.Subscribe(events.SyncEvent{})
	c.syncEventSub = s4
}

// Unsubscribe all messageEventSub
func (c *core) unsubscribeEvents() {
	c.messageEventSub.Unsubscribe()
	c.newUnminedBlockEventSub.Unsubscribe()
	c.timeoutEventSub.Unsubscribe()
	c.committedSub.Unsubscribe()
	c.syncEventSub.Unsubscribe()
}

// TODO: update all of the TypeMuxSilent to event.Feed and should not use backend.EventMux for core internal messageEventSub: backlogEvent, TimeoutEvent

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

func (c *core) mainEventLoop(ctx context.Context) {
	// Start a new round from last height + 1
	c.startRound(ctx, 0)

	go c.syncLoop(ctx)

eventLoop:
	for {
		select {
		case ev, ok := <-c.messageEventSub.Chan():
			if !ok {
				break eventLoop
			}
			// A real ev arrived, process interesting content
			switch e := ev.Data.(type) {
			case events.MessageEvent:
				if len(e.Payload) == 0 {
					c.logger.Error("core.mainEventLoop Get message(MessageEvent) empty payload")
				}

				if err := c.handleMsg(ctx, e.Payload); err != nil {
					c.logger.Debug("core.mainEventLoop Get message(MessageEvent) payload failed", "err", err)
					continue
				}
				c.backend.Gossip(ctx, c.committeeSet().Committee(), e.Payload)
			case proposalEvent:
				err := c.handleProposal(ctx, e.proposal)
				if err != nil {
					c.logger.Debug("core.mainEventLoop handleProposal message failed", "err", err)
				}
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
			switch ev.Data.(type) {
			case events.CommitEvent:
				c.handleCommit(ctx)
			}
		case <-ctx.Done():
			c.logger.Info("mainEventLoop is stopped", "event", ctx.Err())
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

	round := c.Round()
	height := c.Height()

	// Ask for sync when the engine starts
	c.backend.AskSync(c.lastHeader)

eventLoop:
	for {
		select {
		case <-timer.C:
			currentRound := c.Round()
			currentHeight := c.Height()

			// we only ask for sync if the current view stayed the same for the past 10 seconds
			if currentHeight.Cmp(height) == 0 && currentRound == round {
				c.backend.AskSync(c.lastHeader)
			}
			round = currentRound
			height = currentHeight
			timer = time.NewTimer(10 * time.Second)

		case ev, ok := <-c.syncEventSub.Chan():
			if !ok {
				break eventLoop
			}
			event := ev.Data.(events.SyncEvent)
			c.logger.Info("Processing sync message", "from", event.Addr)
			c.backend.SyncPeer(event.Addr)
		case <-ctx.Done():
			c.logger.Info("syncLoop is stopped", "event", ctx.Err())
			break eventLoop
		}
	}

	c.stopped <- struct{}{}
}

// sendEvent sends event to mux
func (c *core) sendEvent(ev interface{}) {
	c.backend.Post(ev)
}

type consensusMessage struct {
	step       uint8
	height     uint64
	round      int64
	value      common.Hash
	validRound int64
}

func (c *core) handleMsg(ctx context.Context, payload []byte) error {

	/*
		Basic validity checks
	*/

	m := new(Message)

	// Decode message
	err := rlp.DecodeBytes(payload, m)
	if err != nil {
		return err
	}

	payload, err = m.PayloadNoSig()
	if err != nil {
		return err
	}

	hash := autonitycrypto.Keccak256(payload)

	// Set the hash on the message so that it can be used for indexing.
	m.Hash = common.BytesToHash(hash)

	// Check we haven't already processed this message
	if c.msgCache.Message(m.Hash) != nil {
		// Message was already processed
		return nil
	}

	var proposal Proposal
	var preVote Vote
	var preCommit Vote
	var cm ConsensusMsg
	switch m.Code {
	case msgProposal:
		err := m.Decode(&proposal)
		if err != nil {
			return errFailedDecodeProposal
		}
		cm = &proposal
		err = c.msgCache.addProposal(&proposal, m)
		if err != nil {
			// could be multipe proposal messages from the same proposer
			return err
		}
	case msgPrevote:
		err := m.Decode(&preVote)
		if err != nil {
			return errFailedDecodePrevote
		}
		cm = &preVote
		err = c.msgCache.addPrevote(&preVote, m)
		if err != nil {
			// could be multipe prevotes from same validator
			return err
		}
	case msgPrecommit:
		err := m.Decode(&preCommit)
		if err != nil {
			return errFailedDecodePrecommit
		}
		// Check the committed seal matches the block hash if its a precommit.
		// If not we ignore the message.
		err := c.verifyCommittedSeal(m.Address, append([]byte(nil), m.CommittedSeal...), cm.ProposedValueHash(), cm.GetRound(), cm.GetHeight())
		if err != nil {
			return err
		}
		cm = &preCommit
		err = c.msgCache.addPrecommit(&preCommit, m)
		if err != nil {
			// could be multipe precommits from same validator
			return err
		}
	default:
		return fmt.Errorf("unrecognised consensus message code %q", m.Code)
	}

	// If this message is for a future height then we cannot validate it
	// because we lack the relevant header, we will process it when we reach
	// that height. If it is for a previous height then we are not intersted in
	// it. But it has been added to the msg cache in case other peers would
	// like to sync it.
	if cm.GetHeight().Uint64() != c.Height().Uint64() {
		// Nothing to do here
		return nil
	}

	var vr int64
	if m.Code == msgProposal {
		vr = proposal.ValidRound
	}
	conMsg := &consensusMessage{
		step:       uint8(m.Code),
		height:     cm.GetHeight().Uint64(),
		round:      cm.GetRound(),
		value:      cm.ProposedValueHash(),
		validRound: vr,
	}

	return handleCurrentHeightMessage(conMsg)

}

func (c *core) handleCurrentHeightMessage(cm consensusMessage) error {
	/*
		Domain specific validity checks, now we know that we are at the same
		height as this message we can rely on lastHeader.
	*/

	// Check that the message came from a committee member, if not we ignore it.
	if c.lastHeader.CommitteeMember(m.Address) == nil {
		// TODO turn this into an error type that can be checked for at a
		// higher level to close the connection to this peer.
		return fmt.Errorf("received message from non committee member: %v", m)
	}

	// Again we ignore messges with invalid signatures, they cannot be trusted.
	// TODO replace crypto.CheckValidatorSignature with something more
	// efficient, so we don't have to hash twice.
	address, err := crypto.CheckValidatorSignature(c.lastHeader, payload, m.Signature)
	if err != nil {
		return err
	}

	if address != m.Address {
		// TODO why is Address even a field of Message when the address can be derived?
		return fmt.Errorf("address in message %q and address derived from signature %q don't match", m.Address, address)
	}

	switch m.Code {
	case msgProposal:
		// We ignore proposals from non proposers
		if !c.isProposerMsg(proposal.Round, m.Address) {
			c.logger.Warn("Ignore proposal messages from non-proposer")
			return errNotFromProposer
		}
	}

	switch m.Code {
	case msgProposal:
		err = c.handleProposal(ctx, &proposal)
	case msgPrevote:
		err = c.handlePrevote(ctx, &preVote, c.lastHeader)
	case msgPrecommit:
		err = c.handlePrecommit(ctx, &preCommit, c.lastHeader)
	default:
		panic("should never happen")
	}
	return err

}

func (c *core) handleFutureRoundMsg(ctx context.Context, msg *Message, sender types.CommitteeMember) {
	// Decoding functions can't fail here
	msgRound, err := msg.Round()
	if err != nil {
		c.logger.Error("handleFutureRoundMsg msgRound", "err", err)
		return
	}
	if _, ok := c.futureRoundChange[msgRound]; !ok {
		c.futureRoundChange[msgRound] = make(map[common.Address]uint64)
	}
	c.futureRoundChange[msgRound][sender.Address] = sender.VotingPower.Uint64()

	var totalFutureRoundMessagesPower uint64
	for _, power := range c.futureRoundChange[msgRound] {
		totalFutureRoundMessagesPower += power
	}

	if totalFutureRoundMessagesPower > c.committeeSet().F() {
		c.logger.Info("Received ceil(N/3) - 1 messages power for higher round", "New round", msgRound)
		c.startRound(ctx, msgRound)
	}
}
