package core

import (
	"context"
	"github.com/clearmatics/autonity/consensus/tendermint/wal"
	"math/big"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint"
)

// Start implements core.Engine.Start
func (c *core) Start() error {
	c.subscribeEvents()

	// set currentRoundState before starting go routines
	lastCommittedProposalBlock, _ := c.backend.LastCommittedProposal()
	height := new(big.Int).Add(lastCommittedProposalBlock.Number(), common.Big1)
	c.currentRoundState = NewRoundState(big.NewInt(0), height)

	c.wal = wal.New(c.config.WALDir, height)

	//We need a separate go routine to keep c.latestPendingUnminedBlock up to date
	go c.handleNewUnminedBlockEvent()

	var ctx context.Context
	ctx, c.cancel = context.WithCancel(context.Background())

	//We want to sequentially handle all the event which modify the current consensus state
	go c.handleConsensusEvents(ctx)

	go c.handleConsensusStuck(ctx)

	return nil
}

// Stop implements core.Engine.Stop
func (c *core) Stop() error {
	c.stopFutureProposalTimer()
	c.unsubscribeEvents()

	_ = c.proposeTimeout.stopTimer()
	_ = c.prevoteTimeout.stopTimer()
	_ = c.precommitTimeout.stopTimer()

	// Make sure the handler goroutine exits
	c.cancel()

	c.wal.Close()

	return nil
}

func (c *core) subscribeEvents() {
	c.messageEventSub = c.backend.EventMux().Subscribe(
		// external messages
		tendermint.MessageEvent{},
		// internal messages
		backlogEvent{},
	)
	c.newUnminedBlockEventSub = c.backend.EventMux().Subscribe(
		tendermint.NewUnminedBlockEvent{},
	)
	c.timeoutEventSub = c.backend.EventMux().Subscribe(
		timeoutEvent{},
	)
	c.committedSub = c.backend.EventMux().Subscribe(
		tendermint.CommitEvent{},
	)
}

// Unsubscribe all messageEventSub
func (c *core) unsubscribeEvents() {
	c.messageEventSub.Unsubscribe()
	c.newUnminedBlockEventSub.Unsubscribe()
	c.timeoutEventSub.Unsubscribe()
	c.committedSub.Unsubscribe()
}

// TODO: update all of the TypeMuxSilent to event.Feed and should not use backend.EventMux for core internal messageEventSub: backlogEvent, timeoutEvent

func (c *core) handleNewUnminedBlockEvent() {
	for e := range c.newUnminedBlockEventSub.Chan() {
		c.logger.Debug("Started handling tendermint.NewUnminedBlockEvent")
		newUnminedBlockEvent := e.Data.(tendermint.NewUnminedBlockEvent)

		pb := &newUnminedBlockEvent.NewUnminedBlock

		err := c.handleUnminedBlock(pb)
		switch err {
		case consensus.ErrFutureBlock:
			c.storeUnminedBlockMsg(pb)
		case nil:
			//nothing to do
		default:
			c.logger.Error("core.handleNewUnminedBlockEvent Get message(NewUnminedBlockEvent) failed", "err", err)
		}

		c.logger.Debug("Finished handling tendermint.NewUnminedBlockEvent")
	}
}

func (c *core) handleConsensusEvents(ctx context.Context) {
	// Start a new round from last height + 1
	// Do not want to block listening for events
	c.startRound(ctx, common.Big0)

	for {
		select {
		case ev, ok := <-c.messageEventSub.Chan():
			if !ok {
				return
			}
			// A real ev arrived, process interesting content
			switch e := ev.Data.(type) {
			case tendermint.MessageEvent:
				if len(e.Payload) == 0 {
					c.logger.Error("core.handleConsensusEvents Get message(MessageEvent) empty payload")
				}

				c.logger.Debug("Started handling tendermint.MessageEvent")
				if err := c.handleMsg(ctx, e.Payload); err != nil {
					c.logger.Error("core.handleConsensusEvents Get message(MessageEvent) payload failed", "err", err)
					c.logger.Debug("Finished handling tendermint.MessageEvent with ERROR", "err", err)
					continue
				}
				c.backend.Gossip(ctx, c.valSet.Copy(), e.Payload)
				c.logger.Debug("Finished handling tendermint.MessageEvent")
			case backlogEvent:
				// No need to check signature for internal messages
				c.logger.Debug("Started handling backlogEvent")
				err := c.handleCheckedMsg(ctx, e.msg, e.src)
				if err != nil {
					c.logger.Error("core.handleConsensusEvents handleCheckedMsg message failed", "err", err)
					c.logger.Debug("Finished handling backlogEvent with ERROR", "err", err)
					continue
				}

				p, err := e.msg.Payload()
				if err != nil {
					c.logger.Error("core.handleConsensusEvents Get message payload failed", "err", err)
					c.logger.Debug("Finished handling backlogEvent with ERROR", "err", err)
					continue
				}
				c.backend.Gossip(ctx, c.valSet.Copy(), p)
				c.logger.Debug("Finished handling backlogEvent")
			}
		case ev, ok := <-c.timeoutEventSub.Chan():
			if !ok {
				return
			}
			if timeoutE, ok := ev.Data.(timeoutEvent); ok {
				c.logger.Debug("Started handling timeoutEvent")
				switch timeoutE.step {
				case msgProposal:
					c.handleTimeoutPropose(ctx, timeoutE)
				case msgPrevote:
					c.handleTimeoutPrevote(ctx, timeoutE)
				case msgPrecommit:
					c.handleTimeoutPrecommit(ctx, timeoutE)
				}
				c.logger.Debug("Finished handling timeoutEvent")
			}
		case ev, ok := <-c.committedSub.Chan():
			if !ok {
				return
			}
			switch ev.Data.(type) {
			case tendermint.CommitEvent:
				c.logger.Debug("Started handling CommitEvent")
				c.handleCommit(ctx)
				c.logger.Debug("Finished handling CommitEvent")
			}
		}
	}
}

func (c *core) handleConsensusStuck(ctx context.Context) {
	ticker := time.NewTicker(time.Second*time.Duration(c.config.RequestTimeout/1000))
	defer ticker.Stop()

	currentHeight := new(big.Int).Set(c.currentRoundState.Height())
	currentRound := new(big.Int).Set(c.currentRoundState.Round())

	// once in a while check height/round and if it's dont change - get messages from WAL and send them again
	for {
		select {
		case <-ticker.C:
			height := c.currentRoundState.Height()
			round := c.currentRoundState.Round()

			c.logger.Error("WAL", "height", c.currentRoundState.Height().String(), "round", c.currentRoundState.Round().String(), "currentHeight", currentHeight.String(), "currentRound", currentRound.String())

			if height.Cmp(currentHeight) != 0 {
				currentHeight.Set(height)
				currentRound.Set(round)
				continue
			}

			if round.Cmp(currentRound) != 0 {
				currentRound.Set(round)
				continue
			}

			pastMessages, err := c.wal.Get(height)
			if err != nil {
				c.logger.Error("WAL: cant get messages", "height", height.String(), "round", round.String(), "err", err.Error())
				continue
			}

			c.logger.Warn("WAL: broadcasting", "height", c.currentRoundState.Height().String(), "round", c.currentRoundState.Round().String(), "currentHeight", currentHeight.String(), "currentRound", currentRound.String(), "msg", len(pastMessages))
			for _, msg := range pastMessages {
				c.logger.Debug("WAL: broadcasting message", "height", height.String(), "round", round.String(), "msg", msg)

				if err = c.backend.Broadcast(ctx, c.valSet.Copy(), msg); err != nil {
					c.logger.Error("WAL: failed to broadcast message", "height", height.String(), "round", round.String(), "msg", msg, "err", err)
					continue
				}
			}

		case <-ctx.Done():
			return
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

	if err := msg.FromPayload(payload, c.checkValidatorSignature); err != nil {
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

	return c.handleCheckedMsg(ctx, msg, sender)
}

func (c *core) handleCheckedMsg(ctx context.Context, msg *message, sender tendermint.Validator) error {
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
				var p tendermint.Proposal
				if e := msg.Decode(&p); e != nil {
					return errFailedDecodeProposal
				}
				msgRound = p.Round.Int64()

			} else {
				var v tendermint.Vote
				if e := msg.Decode(&v); e != nil {
					return errFailedDecodeVote
				}
				msgRound = v.Round.Int64()
			}

			c.futureRoundsChange[msgRound] = c.futureRoundsChange[msgRound] + 1
			totalFutureRoundMessages := c.futureRoundsChange[msgRound]

			if totalFutureRoundMessages >= int64(c.valSet.F()) {
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
