package core

import (
	"context"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"github.com/autonity/autonity/consensus/tendermint/core/types"
	"github.com/autonity/autonity/consensus/tendermint/crypto"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"time"
)

// todo: resolve proper tendermint state synchronization timeout from block period.
const syncTimeOut = 30 * time.Second

// Start implements core.Tendermint.Start
func (c *Core) Start(ctx context.Context, contract *autonity.Contract) {
	c.autonityContract = contract
	committeeSet := committee.NewWeightedRandomSamplingCommittee(c.backend.BlockChain().CurrentBlock(),
		c.autonityContract,
		c.backend.BlockChain())
	c.setCommitteeSet(committeeSet)

	ctx, c.cancel = context.WithCancel(ctx)
	c.subscribeEvents()
	// Tendermint Finite State Machine discrete event loop
	go c.mainEventLoop(ctx)
	go c.backend.HandleUnhandledMsgs(ctx)
}

// Stop implements Core.Engine.Stop
func (c *Core) Stop() {
	c.logger.Info("stopping tendermint.Core", "addr", c.address.String())

	_ = c.proposeTimeout.StopTimer()
	_ = c.prevoteTimeout.StopTimer()
	_ = c.precommitTimeout.StopTimer()

	c.cancel()

	c.proposer.StopFutureProposalTimer()
	c.unsubscribeEvents()

	// Ensure all event handling go routines exit
	<-c.stopped
	<-c.stopped
}

func (c *Core) subscribeEvents() {
	s := c.backend.Subscribe(events.MessageEvent{}, backlogEvent{}, backlogUncheckedEvent{}, types.CoreStateRequestEvent{})
	c.messageEventSub = s

	s1 := c.backend.Subscribe(events.NewCandidateBlockEvent{})
	c.candidateBlockEventSub = s1

	s2 := c.backend.Subscribe(types.TimeoutEvent{})
	c.timeoutEventSub = s2

	s3 := c.backend.Subscribe(events.CommitEvent{})
	c.committedSub = s3

	s4 := c.backend.Subscribe(events.SyncEvent{})
	c.syncEventSub = s4
}

// Unsubscribe all messageEventSub
func (c *Core) unsubscribeEvents() {
	c.messageEventSub.Unsubscribe()
	c.candidateBlockEventSub.Unsubscribe()
	c.timeoutEventSub.Unsubscribe()
	c.committedSub.Unsubscribe()
	c.syncEventSub.Unsubscribe()
}

func needsPeerDisconnect(err error) bool {
	switch err {
	case constants.ErrFutureHeightMessage:
		fallthrough
	case constants.ErrOldHeightMessage:
		fallthrough
	case constants.ErrOldRoundMessage:
		fallthrough
	case constants.ErrFutureRoundMessage:
		fallthrough
	case constants.ErrFutureStepMessage:
		fallthrough
	case constants.ErrNilPrevoteSent:
		fallthrough
	case constants.ErrNilPrecommitSent:
		fallthrough
	case constants.ErrMovedToNewRound:
		return false
	default:
		return true
	}
}

func (c *Core) mainEventLoop(ctx context.Context) {
	// Start a new round from last height + 1
	c.StartRound(ctx, 0)

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
				msg := new(messageutils.Message)
				if err := msg.FromPayload(e.Payload); err != nil {
					c.logger.Error("consensus message invalid payload", "err", err)
					select {
					case e.ErrCh <- err:
					default: // do nothing
					}
					continue
				}
				if err := c.handleMsg(ctx, msg); err != nil {
					c.logger.Debug("MessageEvent payload failed", "err", err)
					// filter errors which needs remote peer disconnection
					if needsPeerDisconnect(err) {
						select {
						case e.ErrCh <- err:
						default: // do nothing
						}
					}
					continue
				}
				c.backend.Gossip(ctx, c.CommitteeSet().Committee(), e.Payload)
			case backlogEvent:
				// No need to check signature for internal messages
				c.logger.Debug("started handling backlogEvent")
				if err := c.handleCheckedMsg(ctx, e.msg); err != nil {
					c.logger.Debug("backlogEvent message handling failed", "err", err)
					continue
				}
				c.backend.Gossip(ctx, c.CommitteeSet().Committee(), e.msg.GetPayload())

			case backlogUncheckedEvent:
				c.logger.Debug("started handling backlogUncheckedEvent")
				if err := c.handleMsg(ctx, e.msg); err != nil {
					c.logger.Debug("backlogUncheckedEvent message failed", "err", err)
					continue
				}
				c.backend.Gossip(ctx, c.CommitteeSet().Committee(), e.msg.GetPayload())
			case types.CoreStateRequestEvent:
				// Process Tendermint state dump request.
				c.handleStateDump(e)
			}
		case ev, ok := <-c.timeoutEventSub.Chan():
			if !ok {
				break eventLoop
			}
			if timeoutE, ok := ev.Data.(types.TimeoutEvent); ok {
				switch timeoutE.Step {
				case messageutils.MsgProposal:
					c.handleTimeoutPropose(ctx, timeoutE)
				case messageutils.MsgPrevote:
					c.handleTimeoutPrevote(ctx, timeoutE)
				case messageutils.MsgPrecommit:
					c.handleTimeoutPrecommit(ctx, timeoutE)
				}
			}
		case ev, ok := <-c.committedSub.Chan():
			if !ok {
				break eventLoop
			}
			switch ev.Data.(type) {
			case events.CommitEvent:
				c.precommiter.HandleCommit(ctx)
			}
		case ev, ok := <-c.candidateBlockEventSub.Chan():
			if !ok {
				break eventLoop
			}
			newCandidateBlockEvent := ev.Data.(events.NewCandidateBlockEvent)
			pb := &newCandidateBlockEvent.NewCandidateBlock
			c.proposer.HandleNewCandidateBlockMsg(ctx, pb)
		case <-ctx.Done():
			c.logger.Info("mainEventLoop is stopped", "event", ctx.Err())
			break eventLoop
		}
	}

	c.stopped <- struct{}{}
}

func (c *Core) syncLoop(ctx context.Context) {
	/*
		this method is responsible for asking the network to send us the current consensus state
		and to process sync queries events.
	*/
	timer := time.NewTimer(syncTimeOut)

	round := c.Round()
	height := c.Height()

	// Ask for sync when the engine starts
	c.backend.AskSync(c.LastHeader())

eventLoop:
	for {
		select {
		case <-timer.C:
			currentRound := c.Round()
			currentHeight := c.Height()

			// we only ask for sync if the current view stayed the same for the past 10 seconds
			if currentHeight.Cmp(height) == 0 && currentRound == round {
				c.backend.AskSync(c.LastHeader())
			}
			round = currentRound
			height = currentHeight
			timer = time.NewTimer(syncTimeOut)

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

// SendEvent sends event to mux
func (c *Core) SendEvent(ev interface{}) {
	c.backend.Post(ev)
}

func (c *Core) handleMsg(ctx context.Context, msg *messageutils.Message) error {

	msgHeight, err := msg.Height()
	if err != nil {
		return err
	}
	if msgHeight.Cmp(c.Height()) > 0 {
		// Future height message. Skip processing and put it in the untrusted backlog buffer.
		c.storeUncheckedBacklog(msg)
		return constants.ErrFutureHeightMessage // No gossip
	}
	if msgHeight.Cmp(c.Height()) < 0 {
		// Old height messages. Do nothing.
		return constants.ErrOldHeightMessage // No gossip
	}

	if _, err = msg.Validate(crypto.CheckValidatorSignature, c.LastHeader()); err != nil {
		c.logger.Error("Failed to validate message", "err", err)
		return err
	}

	return c.handleCheckedMsg(ctx, msg)
}

func (c *Core) handleFutureRoundMsg(ctx context.Context, msg *messageutils.Message, sender common.Address) {
	// Decoding functions can't fail here
	msgRound, err := msg.Round()
	if err != nil {
		c.logger.Error("handleFutureRoundMsg msgRound", "err", err)
		return
	}
	if _, ok := c.futureRoundChange[msgRound]; !ok {
		c.futureRoundChange[msgRound] = make(map[common.Address]uint64)
	}
	c.futureRoundChange[msgRound][sender] = msg.Power

	var totalFutureRoundMessagesPower uint64
	for _, power := range c.futureRoundChange[msgRound] {
		totalFutureRoundMessagesPower += power
	}

	if totalFutureRoundMessagesPower > c.CommitteeSet().F() {
		c.logger.Info("Received ceil(N/3) - 1 messages power for higher round", "New round", msgRound)
		c.StartRound(ctx, msgRound)
	}
}

func (c *Core) handleCheckedMsg(ctx context.Context, msg *messageutils.Message) error {
	logger := c.logger.New("address", c.address, "from", msg.Address)

	// Store the message if it's a future message
	testBacklog := func(err error) error {
		// We want to store only future messages in backlog
		if err == constants.ErrFutureHeightMessage {
			//Future messages should never be processed and reach this point. Panic.
			panic("Processed future message")
		} else if err == constants.ErrFutureRoundMessage {
			logger.Debug("Storing future round message in backlog")
			c.storeBacklog(msg, msg.Address)
			// decoding must have been successful to return
			c.handleFutureRoundMsg(ctx, msg, msg.Address)
		} else if err == constants.ErrFutureStepMessage {
			logger.Debug("Storing future step message in backlog")
			c.storeBacklog(msg, msg.Address)
		}

		return err
	}

	switch msg.Code {
	case messageutils.MsgProposal:
		logger.Debug("tendermint.MessageEvent: PROPOSAL")
		return testBacklog(c.proposer.HandleProposal(ctx, msg))
	case messageutils.MsgPrevote:
		logger.Debug("tendermint.MessageEvent: PREVOTE")
		return testBacklog(c.prevoter.HandlePrevote(ctx, msg))
	case messageutils.MsgPrecommit:
		logger.Debug("tendermint.MessageEvent: PRECOMMIT")
		return testBacklog(c.precommiter.HandlePrecommit(ctx, msg))
	default:
		logger.Error("Invalid message", "msg", msg)
	}

	return constants.ErrInvalidMessage
}
