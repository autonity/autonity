package core

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
)

// todo: resolve proper tendermint state synchronization timeout from block period.
const syncTimeOut = 30 * time.Second

var ErrValidatorJailed = errors.New("jailed validator")

// Start implements core.Tendermint.Start
func (c *Core) Start(ctx context.Context, contract *autonity.ProtocolContracts) {
	c.protocolContracts = contract
	committeeSet := committee.NewWeightedRandomSamplingCommittee(c.backend.BlockChain().CurrentBlock(),
		c.protocolContracts,
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

	c.stopAllTimeouts()

	c.cancel()

	c.proposer.StopFutureProposalTimer()
	c.unsubscribeEvents()

	// Ensure all event handling go routines exit
	<-c.stopped
	<-c.stopped
}

func (c *Core) subscribeEvents() {
	c.messageSub = c.backend.Subscribe(
		events.MessageEvent{},
		backlogMessageEvent{},
		backlogUntrustedMessageEvent{},
		StateRequestEvent{})
	c.candidateBlockSub = c.backend.Subscribe(events.NewCandidateBlockEvent{})
	c.timeoutEventSub = c.backend.Subscribe(TimeoutEvent{})
	c.committedSub = c.backend.Subscribe(events.CommitEvent{})
	c.syncEventSub = c.backend.Subscribe(events.SyncEvent{})
}

// Unsubscribe all messageSub
func (c *Core) unsubscribeEvents() {
	c.messageSub.Unsubscribe()
	c.candidateBlockSub.Unsubscribe()
	c.timeoutEventSub.Unsubscribe()
	c.committedSub.Unsubscribe()
	c.syncEventSub.Unsubscribe()
}

func shouldDisconnectSender(err error) bool {
	switch {
	case errors.Is(err, constants.ErrFutureHeightMessage):
		fallthrough
	case errors.Is(err, constants.ErrOldHeightMessage):
		fallthrough
	case errors.Is(err, constants.ErrOldRoundMessage):
		fallthrough
	case errors.Is(err, constants.ErrFutureRoundMessage):
		fallthrough
	case errors.Is(err, constants.ErrNilPrevoteSent):
		fallthrough
	case errors.Is(err, constants.ErrNilPrecommitSent):
		fallthrough
	case errors.Is(err, constants.ErrMovedToNewRound):
		fallthrough
	case errors.Is(err, constants.ErrHeightClosed):
		fallthrough
	case errors.Is(err, constants.ErrAlreadyHaveBlock):
		fallthrough
	case errors.Is(err, consensus.ErrPrunedAncestor):
		fallthrough
	case errors.Is(err, constants.ErrAlreadyHaveProposal):
		return false
	case errors.Is(err, ErrValidatorJailed):
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
		case ev, ok := <-c.messageSub.Chan():
			if !ok {
				break eventLoop
			}
			// A real ev arrived, process interesting content
			switch e := ev.Data.(type) {
			case events.MessageEvent:

				// At this stage, a message is parsed and all the internal fields must be accessible
				if err := c.handleMsg(ctx, e.Message); err != nil {
					c.logger.Debug("MessageEvent payload failed", "err", err)
					// filter errors which needs remote peer disconnection
					if shouldDisconnectSender(err) {
						tryDisconnect(e.ErrCh, err)
					}
					continue
				}
				c.backend.Gossip(c.CommitteeSet().Committee(), e.Message)
			case backlogMessageEvent:
				// No need to check signature for internal messages
				c.logger.Debug("Started handling consensus backlog event")
				if err := c.handleValidMsg(ctx, e.msg); err != nil {
					c.logger.Debug("BacklogEvent message handling failed", "err", err)
					continue
				}
				c.backend.Gossip(c.CommitteeSet().Committee(), e.msg)

			case backlogUntrustedMessageEvent:
				c.logger.Debug("Started handling backlog unchecked event")
				// messages in the untrusted buffer were successfully decoded
				if err := c.handleMsg(ctx, e.msg); err != nil {
					c.logger.Debug("BacklogUntrustedMessageEvent message failed", "err", err)
					continue
				}
				c.backend.Gossip(c.CommitteeSet().Committee(), e.msg)
			case StateRequestEvent:
				// Process Tendermint state dump request.
				c.handleStateDump(e)
			}
		case ev, ok := <-c.timeoutEventSub.Chan():
			if !ok {
				break eventLoop
			}
			if timeoutE, ok := ev.Data.(TimeoutEvent); ok {
				// if we already decided on this height block, ignore the timeout. It is useless by now.
				if c.step == PrecommitDone {
					c.logTimeoutEvent("Timer expired while at PrecommitDone step, ignoring", "", timeoutE)
					continue
				}
				switch timeoutE.Step {
				case Propose:
					c.handleTimeoutPropose(ctx, timeoutE)
				case Prevote:
					c.handleTimeoutPrevote(ctx, timeoutE)
				case Precommit:
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
		case ev, ok := <-c.candidateBlockSub.Chan():
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
				c.logger.Warn("⚠️ Consensus liveliness lost")
				c.logger.Warn("Broadcasting sync request..")
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
func (c *Core) handleMsg(ctx context.Context, msg message.Msg) error {
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
	// current height message

	// if we already decided on this height block, discard the message. It is useless by now.
	if c.step == PrecommitDone {
		return constants.ErrHeightClosed
	}

	if err := msg.Validate(c.LastHeader().CommitteeMember); err != nil {
		c.logger.Error("Failed to validate message", "err", err)
		c.logger.Error(msg.String())
		return err
	}
	if c.backend.IsJailed(msg.Sender()) {
		c.logger.Debug("Jailed validator, ignoring message", "address", msg.Sender())
		return ErrValidatorJailed
	}
	return c.handleValidMsg(ctx, msg)
}

func (c *Core) handleValidMsg(ctx context.Context, msg message.Msg) error {
	logger := c.logger.New("from", msg.Sender())

	// These checks need to be repeated here due to backlogged messages being re-injected

	if c.Height().Uint64() > msg.H() {
		return constants.ErrOldHeightMessage
	}

	if c.Height().Uint64() < msg.H() {
		panic("Processing future message")
	}

	// if we already decided on this height block, discard the message. It is useless by now.
	if c.step == PrecommitDone {
		return constants.ErrHeightClosed
	}

	// Store the message if it's a future message
	testBacklog := func(err error) error {
		// We want to store only future messages in backlog
		switch {
		case errors.Is(err, constants.ErrFutureRoundMessage):
			logger.Debug("Storing future round message in backlog")
			c.storeBacklog(msg, msg.Sender())
			// decoding must have been successful to return
			c.roundSkipCheck(ctx, msg, msg.Sender())
		}
		return err
	}

	switch m := msg.(type) {
	case *message.Propose:
		logger.Debug("Handling Proposal")
		return testBacklog(c.proposer.HandleProposal(ctx, m))
	case *message.Prevote:
		logger.Debug("Handling Prevote")
		return testBacklog(c.prevoter.HandlePrevote(ctx, m))
	case *message.Precommit:
		logger.Debug("Handling Precommit")
		return testBacklog(c.precommiter.HandlePrecommit(ctx, m))
	default:
		logger.Error("Invalid message", "msg", msg)
	}
	// this should never happen, decoding only returns us propose, prevote or precommit
	panic("handled message that is not propose, prevote or precommit")
}

func tryDisconnect(errorCh chan<- error, err error) {
	select {
	case errorCh <- err:
	default: // do nothing
	}
}
