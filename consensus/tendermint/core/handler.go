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
	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	"github.com/clearmatics/autonity/consensus/tendermint/crypto"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/contracts/autonity"
	autonitycrypto "github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/rlp"
	"github.com/davecgh/go-spew/spew"
)

// Start implements core.Tendermint.Start
func (c *core) Start(ctx context.Context, contract *autonity.Contract) {
	// Set the autonity contract
	c.autonityContract = contract
	ctx, c.cancel = context.WithCancel(ctx)

	// Subscribe
	c.eventsSub = c.backend.Subscribe(events.MessageEvent{}, events.NewUnminedBlockEvent{}, &algorithm.ConsensusMessage{}, &algorithm.Timeout{}, events.CommitEvent{})
	c.syncEventSub = c.backend.Subscribe(events.SyncEvent{})
	c.newUnminedBlockEventSub = c.backend.Subscribe(events.NewUnminedBlockEvent{})

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

	c.cancel()

	// Unsubscribe
	c.eventsSub.Unsubscribe()
	c.syncEventSub.Unsubscribe()

	// Ensure all event handling go routines exit
	<-c.stopped
	<-c.stopped
	<-c.stopped
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

func (c *core) handleResult(ctx context.Context, m *algorithm.ConsensusMessage, t *algorithm.Timeout, proposal *algorithm.ConsensusMessage) {
	if proposal != nil {
		// A decision has been reached
		err := c.Commit(proposal)
		if err != nil {
			panic(fmt.Sprintf("Failed to commit proposal: %s err: %v", spew.Sdump(proposal), err))
		}

		c.updateLatestBlock()
	}
	switch {
	case m != nil:
		go c.broadcast(ctx, m)
		go c.backend.Post(m)
	case t != nil:
		time.AfterFunc(time.Duration(t.Delay)*time.Second, func() {
			c.backend.Post(t)
		})
	}
}

func (c *core) mainEventLoop(ctx context.Context) {
	// Start a new round from last height + 1
	c.algo = algorithm.New(algorithm.NodeID(c.address), nil)
	m, t := c.algo.StartRound(c.Height().Uint64(), 0)
	c.handleResult(ctx, m, t, nil)
	go c.syncLoop(ctx)

eventLoop:
	for {
		select {
		case ev, ok := <-c.eventsSub.Chan():
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
			case *algorithm.ConsensusMessage:
				// This is a message we sent ourselves we do not need to broadcast it
				if c.Height().Uint64() == e.Height {
					m, t, p := c.algo.ReceiveMessage(e)
					c.handleResult(ctx, m, t, p)
				}
			case *algorithm.Timeout:
				var m *algorithm.ConsensusMessage
				var t *algorithm.Timeout
				switch e.TimeoutType {
				case algorithm.Propose:
					m = c.algo.OnTimeoutPropose(e.Height, e.Round)
				case algorithm.Prevote:
					m = c.algo.OnTimeoutPrevote(e.Height, e.Round)
				case algorithm.Precommit:
					m, t = c.algo.OnTimeoutPrecommit(e.Height, e.Round)
				}
				c.handleResult(ctx, m, t, nil)
			case events.CommitEvent:
				c.logger.Debug("Received a final committed proposal")
				lastBlock, _ := c.backend.LastCommittedProposal()
				height := new(big.Int).Add(lastBlock.Number(), common.Big1)
				if height.Cmp(c.Height()) == 0 {
					c.logger.Debug("Discarding event as core is at the same height", "height", c.Height())
				} else {
					c.logger.Debug("Received proposal is ahead", "height", c.Height(), "block_height", height)
					c.updateLatestBlock()
					m, t := c.algo.StartRound(c.height.Uint64(), 0)
					c.handleResult(ctx, m, t, nil)
				}
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
	timer := time.NewTimer(20 * time.Second)

	height := c.Height()

	// Ask for sync when the engine starts
	c.backend.AskSync(c.lastHeader)

eventLoop:
	for {
		select {
		case <-timer.C:
			currentHeight := c.Height()

			// we only ask for sync if the current view stayed the same for the past 10 seconds
			if currentHeight.Cmp(height) == 0 {
				c.backend.AskSync(c.lastHeader)
			}
			height = currentHeight
			timer = time.NewTimer(20 * time.Second)

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

func (c *core) handleMsg(ctx context.Context, payload []byte) error {

	/*
		Basic validity checks
	*/

	m := new(Message)

	// Set the hash on the message so that it can be used for indexing.
	m.Hash = common.BytesToHash(autonitycrypto.Keccak256(payload))

	// Check we haven't already processed this message
	if c.msgCache.Message(m.Hash) != nil {
		// Message was already processed
		return nil
	}

	// Decode message
	err := rlp.DecodeBytes(payload, m)
	if err != nil {
		return err
	}

	var proposal Proposal
	var preVote Vote
	var preCommit Vote
	var conMsg *algorithm.ConsensusMessage
	switch m.Code {
	case msgProposal:
		err := m.Decode(&proposal)
		if err != nil {
			return errFailedDecodeProposal
		}

		valueHash := proposal.ProposalBlock.Hash()
		conMsg = &algorithm.ConsensusMessage{
			MsgType:    algorithm.Step(m.Code),
			Height:     proposal.Height.Uint64(),
			Round:      proposal.Round,
			Value:      algorithm.ValueID(valueHash),
			ValidRound: proposal.ValidRound,
		}

		err = c.msgCache.addMessage(m, conMsg)
		if err != nil {
			// could be multiple proposal messages from the same proposer
			return err
		}
		c.msgCache.addValue(valueHash, proposal.ProposalBlock)

	case msgPrevote:
		err := m.Decode(&preVote)
		if err != nil {
			return errFailedDecodePrevote
		}
		conMsg = &algorithm.ConsensusMessage{
			MsgType: algorithm.Step(m.Code),
			Height:  preVote.Height.Uint64(),
			Round:   preVote.Round,
			Value:   algorithm.ValueID(preVote.ProposedBlockHash),
		}

		err = c.msgCache.addMessage(m, conMsg)
		if err != nil {
			// could be multiple precommits from same validator
			return err
		}
	case msgPrecommit:
		err := m.Decode(&preCommit)
		if err != nil {
			return errFailedDecodePrecommit
		}
		// Check the committed seal matches the block hash if its a precommit.
		// If not we ignore the message.
		//
		// Note this method does not make use of any blockchain state, so it is
		// safe to call it now. In fact it only uses the logger of c so I think
		// it could easily be detached from c.
		err = c.verifyCommittedSeal(m.Address, append([]byte(nil), m.CommittedSeal...), preCommit.ProposedBlockHash, preCommit.Round, preCommit.Height)
		if err != nil {
			return err
		}
		conMsg = &algorithm.ConsensusMessage{
			MsgType: algorithm.Step(m.Code),
			Height:  preCommit.Height.Uint64(),
			Round:   preCommit.Round,
			Value:   algorithm.ValueID(preVote.ProposedBlockHash),
		}

		err = c.msgCache.addMessage(m, conMsg)
		if err != nil {
			// could be multiple precommits from same validator
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
	if conMsg.Height != c.Height().Uint64() {
		// Nothing to do here
		return nil
	}

	return c.handleCurrentHeightMessage(m, conMsg)

}

func (c *core) handleCurrentHeightMessage(m *Message, cm *algorithm.ConsensusMessage) error {
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

	payload, err := m.PayloadNoSig()
	if err != nil {
		return err
	}

	// Again we ignore messges with invalid signatures, they cannot be trusted.
	// TODO make crypto.CheckValidatorSignature accept Message so that it can
	// handle generating the payload and checking it with the sig and address.
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
		if !c.isProposerMsg(cm.Round, m.Address) {
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
		if _, err := c.backend.VerifyProposal(*c.msgCache.values[common.Hash(cm.Value)]); err == nil {
			c.msgCache.setValid(common.Hash(cm.Value))
		}
	default:
		c.msgCache.setValid(m.Hash)

	}

	cm, t, p := c.algo.ReceiveMessage(cm)
	c.handleResult(context.Background(), cm, t, p)
	return nil
}
