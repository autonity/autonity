package core

import (
	"context"
	"errors"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/core/types"
	"github.com/autonity/autonity/consensus/tendermint/crypto"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"math/big"
	"time"
)

// todo: resolve proper tendermint state synchronization timeout from block period.
const syncTimeOut = 30 * time.Second

var ErrValidatorJailed = errors.New("jailed validator")

// Start implements core.Tendermint.Start
func (c *Core) Start(ctx context.Context, contract *autonity.Contracts) {
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
	c.logger.Debug("Stopping Tendermint Core", "addr", c.address.String())

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
	s := c.backend.Subscribe(events.MessageEvent{}, backlogMessageEvent{}, backlogUntrustedMessageEvent{}, types.CoreStateRequestEvent{})
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

func shouldDisconnectSender(err error) bool {
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
	case ErrValidatorJailed:
		// this one is tricky. Ideally yes, we want to disconnect the sender but we can't
		// really assume that all the other committee members have the same view on the
		// jailed validator list before gossip, that is risking then to disconnect honest nodes.
		// This needs to verified though. Returning false for the time being.
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
				msg, err := message.FromBytes(e.Payload) // todo(youssef): move that ahead
				if err != nil {
					c.logger.Error("Received consensus message with invalid payload", "err", err)
					tryDisconnect(e.ErrCh, err)
					continue
				}
				// At this stage, a message is parsed and all the internal fields must be accessible
				if err := c.handleMsg(ctx, msg); err != nil {
					c.logger.Debug("MessageEvent payload failed", "err", err)
					// filter errors which needs remote peer disconnection
					if shouldDisconnectSender(err) {
						tryDisconnect(e.ErrCh, err)
					}
					continue
				}
				c.backend.Gossip(ctx, c.CommitteeSet().Committee(), e.Payload)
			case backlogMessageEvent:
				// No need to check signature for internal messages
				c.logger.Debug("Started handling consensus backlog event")
				if err := c.handleValidMsg(ctx, e.msg); err != nil {
					c.logger.Debug("BacklogEvent message handling failed", "err", err)
					continue
				}
				c.backend.Gossip(ctx, c.CommitteeSet().Committee(), e.msg.GetBytes())

			case backlogUntrustedMessageEvent:
				c.logger.Debug("Started handling backlog unchecked event")
				// messages in the untrusted buffer were successfully decoded
				if err := c.handleMsg(ctx, e.msg); err != nil {
					c.logger.Debug("BacklogUntrustedMessageEvent message failed", "err", err)
					continue
				}
				c.backend.Gossip(ctx, c.CommitteeSet().Committee(), e.msg.GetBytes())
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
				case consensus.MsgProposal:
					c.handleTimeoutPropose(ctx, timeoutE)
				case consensus.MsgPrevote:
					c.handleTimeoutPrevote(ctx, timeoutE)
				case consensus.MsgPrecommit:
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
			c.logger.Debug("Tendermint core main loop stopped", "event", ctx.Err())
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
			c.logger.Debug("Processing sync message", "from", event.Addr)
			c.backend.SyncPeer(event.Addr)
		case <-ctx.Done():
			c.logger.Debug("syncLoop is stopped", "event", ctx.Err())
			break eventLoop
		}
	}

	c.stopped <- struct{}{}
}

// SendEvent sends event to mux
func (c *Core) SendEvent(ev any) {
	c.backend.Post(ev)
}

// handleMsg assume msg has already been decoded
func (c *Core) handleMsg(ctx context.Context, msg *message.Message) error {
	msgHeight := new(big.Int).SetUint64(msg.H())
	if msgHeight.Cmp(c.Height()) > 0 {
		// Future height message. Skip processing and put it in the untrusted backlog buffer.
		c.storeFutureMessage(msg)
		return constants.ErrFutureHeightMessage // No gossip
	}
	if msgHeight.Cmp(c.Height()) < 0 {
		// Old height messages. Do nothing.
		return constants.ErrOldHeightMessage // No gossip
	}
	if err := msg.Validate(crypto.CheckValidatorSignature, c.LastHeader()); err != nil {
		c.logger.Error("Failed to validate message", "err", err)
		return err
	}
	if c.backend.IsJailed(msg.Address) {
		c.logger.Debug("Jailed validator, ignoring message", "address", msg.Address)
		return ErrValidatorJailed
	}
	return c.handleValidMsg(ctx, msg)
}

func (c *Core) handleFutureRoundMsg(ctx context.Context, msg *message.Message, sender common.Address) {
	// Decoding functions can't fail here
	msgRound := msg.R()
	if _, ok := c.futureRoundChange[msgRound]; !ok {
		c.futureRoundChange[msgRound] = make(map[common.Address]*big.Int)
	}
	c.futureRoundChange[msgRound][sender] = msg.Power

	totalFutureRoundMessagesPower := new(big.Int)
	for _, power := range c.futureRoundChange[msgRound] {
		totalFutureRoundMessagesPower.Add(totalFutureRoundMessagesPower, power)
	}

	if totalFutureRoundMessagesPower.Cmp(c.CommitteeSet().F()) > 0 {
		c.logger.Debug("Received messages with F + 1 total power for a higher round", "New round", msgRound)
		c.StartRound(ctx, msgRound)
	}
}

func (c *Core) handleValidMsg(ctx context.Context, msg *message.Message) error {
	logger := c.logger.New("address", c.address, "from", msg.Address)

	// Store the message if it's a future message
	testBacklog := func(err error) error {
		// We want to store only future messages in backlog
		if err == constants.ErrFutureHeightMessage {
			//Future messages should never be processed and reach this point. Panic.
			panic("Processed future message as a valid message")
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
	case consensus.MsgProposal:
		logger.Debug("Handling Proposal")
		return testBacklog(c.proposer.HandleProposal(ctx, msg))
	case consensus.MsgPrevote:
		logger.Debug("Handling Prevote")
		return testBacklog(c.prevoter.HandlePrevote(ctx, msg))
	case consensus.MsgPrecommit:
		logger.Debug("Handling Precommit")
		return testBacklog(c.precommiter.HandlePrecommit(ctx, msg))
	default:
		logger.Error("Invalid message", "msg", msg)
	}

	return constants.ErrInvalidMessage
}

func tryDisconnect(errorCh chan<- error, err error) {
	select {
	case errorCh <- err:
	default: // do nothing
	}
}
