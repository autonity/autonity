package core

import (
	"context"
	"errors"
	"time"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
)

// todo: resolve proper tendermint state synchronization timeout from block period.
const syncTimeOut = 30 * time.Second

// Start implements core.Tendermint.Start
func (c *Core) Start(ctx context.Context, contract *autonity.ProtocolContracts) {
	c.protocolContracts = contract
	committeeSet := committee.NewWeightedRandomSamplingCommittee(c.backend.BlockChain().CurrentBlock(),
		c.protocolContracts,
		c.backend.BlockChain())
	c.setCommitteeSet(committeeSet)
	ctx, c.cancel = context.WithCancel(ctx)
	c.subscribeEvents()

	// Start a new round from last height + 1
	c.StartRound(ctx, 0)

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
	/* //TODO(lorenzo) refinements2, double check. Also this is kinda broken due to aggregator not sending an ErrCh
	case errors.Is(err, constants.ErrFutureHeightMessage):
		fallthrough
	*/
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
	default:
		return true
	}
}

func (c *Core) quorumFor(code uint8, round int64, value common.Hash) bool {
	quorum := false
	switch code {
	case message.ProposalCode:
		break
	case message.PrevoteCode:
		quorum = (c.messages.GetOrCreate(round).PrevotesPower(value).Cmp(c.CommitteeSet().Quorum()) >= 0)
	case message.PrecommitCode:
		quorum = (c.messages.GetOrCreate(round).PrecommitsPower(value).Cmp(c.CommitteeSet().Quorum()) >= 0)
	}
	return quorum
}

// TODO(lorenzo) maybe I can substitute the existing votes with the aggregate one (instead of just adding it)
func (c *Core) GossipComplexAggregate(code uint8, round int64, value common.Hash) {
	switch code {
	case message.PrevoteCode:
		aggregatePrevote := c.messages.GetOrCreate(round).PrevoteFor(value)
		c.messages.GetOrCreate(round).AddPrevote(aggregatePrevote)
		c.backend.Gossip(c.CommitteeSet().Committee(), aggregatePrevote)
	case message.PrecommitCode:
		aggregatePrecommit := c.messages.GetOrCreate(round).PrecommitFor(value)
		c.messages.GetOrCreate(round).AddPrecommit(aggregatePrecommit)
		c.backend.Gossip(c.CommitteeSet().Committee(), aggregatePrecommit)
	}
}

func (c *Core) mainEventLoop(ctx context.Context) {
	go c.syncLoop(ctx)

eventLoop:
	for {
		select {
		case ev, ok := <-c.messageSub.Chan():
			if !ok {
				break eventLoop
			}
			// An event arrived, process content
			switch e := ev.Data.(type) {
			case events.MessageEvent:
				msg := e.Message

				// check if we have quorum for message type for this round
				hadQuorum := c.quorumFor(msg.Code(), msg.R(), msg.Value())

				if err := c.handleMsg(ctx, msg); err != nil {
					c.logger.Debug("MessageEvent payload failed", "err", err)
					// filter errors which needs remote peer disconnection
					if shouldDisconnectSender(err) {
						tryDisconnect(e.ErrCh, err)
					}
					break
				}

				if !hadQuorum {
					// if we did not have quorum and we reached it now
					// gossip the (complex) aggregate with quorum to everyone instead of the current message
					hasQuorum := c.quorumFor(msg.Code(), msg.R(), msg.Value())
					if hasQuorum {
						c.GossipComplexAggregate(msg.Code(), msg.R(), msg.Value())
						break // do not gossip single message, only complex aggregate
					}
				}

				// gossip message. We should arrive here only if we did not already gossip a complex aggregate
				c.backend.Gossip(c.CommitteeSet().Committee(), msg)
			case backlogMessageEvent:
				// TODO(lorenzo) refinements, should we check for disconnection also here?
				// I am not sure we can get the error ch though

				msg := e.msg

				// check if we have quorum for message type for this round
				hadQuorum := c.quorumFor(msg.Code(), msg.R(), msg.Value())

				c.logger.Debug("Handling consensus backlog event")
				if err := c.handleMsg(ctx, msg); err != nil {
					c.logger.Debug("BacklogEvent message handling failed", "err", err)
					continue
				}

				if !hadQuorum {
					// if we did not have quorum and we reached it now
					// gossip the (complex) aggregate with quorum to everyone instead of the current message
					hasQuorum := c.quorumFor(msg.Code(), msg.R(), msg.Value())
					if hasQuorum {
						c.GossipComplexAggregate(msg.Code(), msg.R(), msg.Value())
						break // do not gossip single message, only complex aggregate
					}
				}

				// gossip message. We should arrive here only if we did not already gossip a complex aggregate
				c.backend.Gossip(c.CommitteeSet().Committee(), msg)
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

func (c *Core) handleMsg(ctx context.Context, msg message.Msg) error {
	// These checks need to be repeated here due to backlogged messages being re-injected
	if c.Height().Uint64() > msg.H() {
		//TODO(lorenzo) refinements, this seems to happen quite a lot. Understand why.
		// also should we gossip old height messages?
		c.logger.Debug("ignoring stale consensus message", "hash", msg.Hash())
		return constants.ErrOldHeightMessage
	}

	if c.Height().Uint64() < msg.H() {
		panic("Processing future height message")
	}

	// if we already decided on this height block, discard the message. It is useless by now.
	if c.step == PrecommitDone {
		return constants.ErrHeightClosed
	}

	var err error
	switch m := msg.(type) {
	case *message.Propose:
		c.logger.Debug("Handling Proposal")
		err = c.proposer.HandleProposal(ctx, m)
	case *message.Prevote:
		c.logger.Debug("Handling Prevote")
		err = c.prevoter.HandlePrevote(ctx, m)
	case *message.Precommit:
		c.logger.Debug("Handling Precommit")
		err = c.precommiter.HandlePrecommit(ctx, m)
	default:
		// this should never happen, decoding only returns us propose, prevote or precommit
		panic("handled message that is not propose, prevote or precommit. Msg: " + msg.String())
	}

	// Store the message if it is a future round message
	if errors.Is(err, constants.ErrFutureRoundMessage) {
		c.logger.Debug("Storing future round message")

		r := msg.R()
		c.futureRoundLock.Lock()
		c.futureRound[r] = append(c.futureRound[r], msg)
		c.futureRoundLock.Unlock()

		c.backend.Post(events.FuturePowerChangeEvent{Height: c.Height().Uint64(), Round: r})

		c.roundSkipCheck(ctx, r)
	}

	return err
}

func tryDisconnect(errorCh chan<- error, err error) {
	//TODO(lorenzo) refinements2, if aggregated vote or local message, we will have no error channel.
	// maybe I can send back the error to the aggregator so the he can do the disconnection and the removal of messages
	if errorCh == nil {
		return
	}

	select {
	case errorCh <- err:
	default: // do nothing
	}
}
