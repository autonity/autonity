package core

import (
	"math/big"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint"
)

// Start implements core.Engine.Start
func (c *core) Start() error {
	// Tests will handle events itself, so we have to make subscribeEvents()
	// be able to call in test.
	c.subscribeEvents()
	go c.handleEvents()

	// Start a new round from last height + 1
	go c.startRound(common.Big0)
	return nil
}

// Stop implements core.Engine.Stop
func (c *core) Stop() error {
	c.stopTimer()
	c.unsubscribeEvents()

	// Make sure the handler goroutine exits
	c.handlerStopCh <- struct{}{}
	return nil
}

// TODO: update all of the TypeMuxSilent to event.Feed and should not use backend.EventMux for core internal events: backlogEvent, timeoutEvent

// Subscribe both internal and external events
func (c *core) subscribeEvents() {
	c.events = c.backend.EventMux().Subscribe(
		// external events
		tendermint.NewUnminedBlockEvent{},
		tendermint.MessageEvent{},
		// internal events
		backlogEvent{},
	)
	c.timeoutSub = c.backend.EventMux().Subscribe(
		timeoutEvent{},
	)
	c.finalCommittedSub = c.backend.EventMux().Subscribe(
		tendermint.CommitEvent{},
	)
}

// Unsubscribe all events
func (c *core) unsubscribeEvents() {
	c.events.Unsubscribe()
	c.timeoutSub.Unsubscribe()
	c.finalCommittedSub.Unsubscribe()
}

func (c *core) handleEvents() {
	// Clear step
	defer func() {
		c.currentRoundState = nil
		<-c.handlerStopCh
	}()

	for {
		select {
		case ev, ok := <-c.events.Chan():
			if !ok {
				return
			}
			// A real ev arrived, process interesting content
			switch e := ev.Data.(type) {
			case tendermint.NewUnminedBlockEvent:
				pb := &e.NewUnminedBlock
				err := c.handleUnminedBlock(pb)
				if err == consensus.ErrFutureBlock {
					c.storeUnminedBlockMsg(pb)
				}
			case tendermint.MessageEvent:
				if err := c.handleMsg(e.Payload); err == nil {
					c.backend.Gossip(c.valSet, e.Payload)
				}
			case backlogEvent:
				// No need to check signature for internal messages
				if err := c.handleCheckedMsg(e.msg, e.src); err == nil {
					p, err := e.msg.Payload()
					if err != nil {
						c.logger.Warn("Get message payload failed", "err", err)
						continue
					}
					c.backend.Gossip(c.valSet, p)
				}
			}
		case ev, ok := <-c.timeoutSub.Chan():
			if !ok {
				return
			}
			if timeoutE, ok := ev.Data.(timeoutEvent); ok {
				switch timeoutE.step {
				case msgProposal:
					c.handleTimeoutPropose(timeoutE)
				case msgPrevote:
					c.handleTimeoutPrevote(timeoutE)
				case msgPrecommit:
					c.handleTimeoutPrecommit(timeoutE)
				}
			}
		case ev, ok := <-c.finalCommittedSub.Chan():
			if !ok {
				return
			}
			switch ev.Data.(type) {
			case tendermint.CommitEvent:
				c.handleCommit()
			}
		}
	}
}

// sendEvent sends events to mux
func (c *core) sendEvent(ev interface{}) {
	c.backend.EventMux().Post(ev)
}

func (c *core) handleMsg(payload []byte) error {
	logger := c.logger.New()

	// Decode message and check its signature
	msg := new(message)
	if err := msg.FromPayload(payload, c.validateFn); err != nil {
		logger.Error("Failed to decode message from payload", "err", err)
		return err
	}

	// Only accept message if the address is valid
	// TODO: the check is already made in c.validateFn
	_, sender := c.valSet.GetByAddress(msg.Address)
	if sender == nil {
		logger.Error("Invalid address in message", "msg", msg)
		return tendermint.ErrUnauthorizedAddress
	}

	return c.handleCheckedMsg(msg, sender)
}

func (c *core) handleCheckedMsg(msg *message, sender tendermint.Validator) error {
	logger := c.logger.New("address", c.address, "from", sender)

	// Store the message if it's a future message
	testBacklog := func(err error) error {
		// We want to store only future messages in backlog
		if err == errFutureRoundMessage {
			//We cannot move to a round in a new height without receiving a new block
			var msgRound int64
			if msg.Code == msgProposal {
				var p tendermint.Proposal
				if e := msg.Decode(p); e != nil {
					return errFailedDecodeProposal
				}
				msgRound = p.Round.Int64()

			} else {
				var v tendermint.Vote
				if e := msg.Decode(v); e != nil {
					// TODO: introduce new error: errFailedDecodeVote
					return errFailedDecodePrecommit
				}
				msgRound = v.Round.Int64()
			}

			c.futureRoundsChange[msgRound] = c.futureRoundsChange[msgRound] + 1
			totalFutureRoundMessages := c.futureRoundsChange[msgRound]

			if totalFutureRoundMessages >= int64(c.valSet.F()) {
				c.startRound(big.NewInt(msgRound))
			}

		}
		if err == errFutureHeightMessage || err == errFutureRoundMessage {
			c.storeBacklog(msg, sender)
		}

		return err
	}

	switch msg.Code {
	case msgProposal:
		return testBacklog(c.handleProposal(msg, sender))
	case msgPrevote:
		return testBacklog(c.handlePrevote(msg, sender))
	case msgPrecommit:
		return testBacklog(c.handlePrecommit(msg, sender))
	default:
		logger.Error("Invalid message", "msg", msg)
	}

	return errInvalidMessage
}
