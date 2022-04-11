package core

import (
	"context"
	"math/big"
	"time"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/crypto"
	"github.com/autonity/autonity/consensus/tendermint/events"
)

// Start implements core.Tendermint.Start
func (c *core) Start(ctx context.Context, contract *autonity.Contract) {
	c.autonityContract = contract
	committeeSet := newWeightedRandomSamplingCommittee(c.backend.BlockChain().CurrentBlock(),
		c.autonityContract,
		c.backend.BlockChain())
	c.setCommitteeSet(committeeSet)

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
	s := c.backend.Subscribe(events.MessageEvent{}, backlogEvent{}, backlogUncheckedEvent{}, coreStateRequestEvent{})
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
			c.storeUnminedBlockMsg(ctx, pb)
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
				msg := new(Message)
				if err := msg.FromPayload(e.Payload); err != nil {
					c.logger.Error("consensus message invalid payload", "err", err)
					continue
				}
				if err := c.handleMsg(ctx, msg); err != nil {
					c.logger.Debug("MessageEvent payload failed", "err", err)
					continue
				}
				c.backend.Gossip(ctx, c.committeeSet().Committee(), e.Payload)
			case backlogEvent:
				// No need to check signature for internal messages
				c.logger.Debug("started handling backlogEvent")
				if err := c.handleCheckedMsg(ctx, e.msg); err != nil {
					c.logger.Debug("backlogEvent message handling failed", "err", err)
					continue
				}
				c.backend.Gossip(ctx, c.committeeSet().Committee(), e.msg.Payload())

			case backlogUncheckedEvent:
				c.logger.Debug("started handling backlogUncheckedEvent")
				if err := c.handleMsg(ctx, e.msg); err != nil {
					c.logger.Debug("backlogUncheckedEvent message failed", "err", err)
					continue
				}
				c.backend.Gossip(ctx, c.committeeSet().Committee(), e.msg.Payload())
			case coreStateRequestEvent:
				// Process Tendermint state dump request.
				c.handleStateDump(e)
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

func (c *core) handleMsg(ctx context.Context, msg *Message) error {

	msgHeight, err := msg.Height()
	if err != nil {
		return err
	}
	if msgHeight.Cmp(c.Height()) > 0 {
		// Future height message. Skip processing and put it in the untrusted backlog buffer.
		c.storeUncheckedBacklog(msg)
		return errFutureHeightMessage // No gossip
	}
	if msgHeight.Cmp(c.Height()) < 0 {
		// Old height messages. Do nothing.
		return errOldHeightMessage // No gossip
	}

	if _, err = msg.Validate(crypto.CheckValidatorSignature, c.lastHeader); err != nil {
		c.logger.Error("Failed to validate message", "err", err)
		return err
	}

	return c.handleCheckedMsg(ctx, msg)
}

func (c *core) handleFutureRoundMsg(ctx context.Context, msg *Message, sender common.Address) {
	// Decoding functions can't fail here
	msgRound, err := msg.Round()
	if err != nil {
		c.logger.Error("handleFutureRoundMsg msgRound", "err", err)
		return
	}
	if _, ok := c.futureRoundChange[msgRound]; !ok {
		c.futureRoundChange[msgRound] = make(map[common.Address]uint64)
	}
	c.futureRoundChange[msgRound][sender] = msg.power

	var totalFutureRoundMessagesPower uint64
	for _, power := range c.futureRoundChange[msgRound] {
		totalFutureRoundMessagesPower += power
	}

	if totalFutureRoundMessagesPower > c.committeeSet().F() {
		c.logger.Info("Received ceil(N/3) - 1 messages power for higher round", "New round", msgRound)
		c.startRound(ctx, msgRound)
	}
}

func (c *core) handleCheckedMsg(ctx context.Context, msg *Message) error {
	logger := c.logger.New("address", c.address, "from", msg.Address)

	// Store the message if it's a future message
	testBacklog := func(err error) error {
		// We want to store only future messages in backlog
		if err == errFutureHeightMessage {
			//Future messages should never be processed and reach this point. Panic.
			panic("Processed future message")
		} else if err == errFutureRoundMessage {
			logger.Debug("Storing future round message in backlog")
			c.storeBacklog(msg, msg.Address)
			// decoding must have been successful to return
			c.handleFutureRoundMsg(ctx, msg, msg.Address)
		} else if err == errFutureStepMessage {
			logger.Debug("Storing future step message in backlog")
			c.storeBacklog(msg, msg.Address)
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
