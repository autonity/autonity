package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	com "github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/metrics"
)

// todo: resolve proper tendermint state synchronization timeout from block period.
const syncTimeOut = 30 * time.Second

// Start implements core.Tendermint.Start
func (c *Core) Start(ctx context.Context, contract *autonity.ProtocolContracts) {
	chainHead := c.backend.HeadBlock().Header()
	epoch, err := c.Backend().EpochByHeight(chainHead.Number.Uint64() + 1)
	if err != nil {
		panic(fmt.Sprintf("failed to fetch epoch information for height: %d, err: %s", chainHead.Number.Uint64()+1, err.Error()))
	}

	c.epoch = epoch
	c.protocolContracts = contract
	committeeSet := com.NewWeightedRandomSamplingCommittee(chainHead, epoch.Committee, c.protocolContracts)
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
	c.candidateBlockCh = make(chan events.NewCandidateBlockEvent, 1)
	c.committedCh = make(chan events.CommitEvent, 1)
	c.timeoutEventSub = c.backend.Subscribe(TimeoutEvent{})
	c.syncEventSub = c.backend.Subscribe(events.SyncEvent{})
}

// Unsubscribe all
func (c *Core) unsubscribeEvents() {
	c.messageSub.Unsubscribe()
	c.timeoutEventSub.Unsubscribe()
	c.syncEventSub.Unsubscribe()
}

func shouldDisconnectSender(err error) bool {
	switch {
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

func recordMessageProcessingTime(code uint8, start time.Time) {
	if !metrics.Enabled {
		return
	}
	switch code {
	case message.ProposalCode:
		MsgProposalBg.Add(time.Since(start).Nanoseconds())
		MsgProposalPackets.Mark(1)
	case message.PrevoteCode:
		MsgPrevoteBg.Add(time.Since(start).Nanoseconds())
		MsgPrevotePackets.Mark(1)
	case message.PrecommitCode:
		MsgPrecommitBg.Add(time.Since(start).Nanoseconds())
		MsgPrecommitPackets.Mark(1)
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

func (c *Core) GossipComplexAggregate(code uint8, round int64, value common.Hash) {
	// We re-add the complex aggregate to the prevote set. If we would substitute the entire set with the complex aggregate, there is a possibility of message loss (if we had multiple un-mergeable complex aggregates in the `messages`). This loss would not harm consensus (we would still have quorum voting power), however it is better to keep all messages in case we have to sync another peer.
	// We can consider changing it only if it considerably harms performance.
	switch code {
	case message.PrevoteCode:
		aggregatePrevote := c.messages.GetOrCreate(round).PrevoteFor(value)
		c.messages.GetOrCreate(round).AddPrevote(aggregatePrevote)
		go c.backend.Gossip(c.CommitteeSet().Committee(), aggregatePrevote)
	case message.PrecommitCode:
		aggregatePrecommit := c.messages.GetOrCreate(round).PrecommitFor(value)
		c.messages.GetOrCreate(round).AddPrecommit(aggregatePrecommit)
		go c.backend.Gossip(c.CommitteeSet().Committee(), aggregatePrecommit)
	}
}

func (c *Core) mainEventLoop(ctx context.Context) {
	go c.syncLoop(ctx)

eventLoop:
	for {
		select {
		case ev, ok := <-c.candidateBlockCh:
			if !ok {
				break eventLoop
			}
			newCandidateBlockEvent := ev
			pb := &newCandidateBlockEvent.NewCandidateBlock
			c.proposer.HandleNewCandidateBlockMsg(ctx, pb)
			if metrics.Enabled && c.IsProposer() {
				CandidateBlockDelayBg.Add(time.Since(newCandidateBlockEvent.CreatedAt).Nanoseconds())
			}
		case ev, ok := <-c.messageSub.Chan():
			if !ok {
				break eventLoop
			}
			start := time.Now()
			// An event arrived, process content
			switch e := ev.Data.(type) {
			case events.MessageEvent:
				if metrics.Enabled {
					AggregatorCoreTransitBg.Add(time.Since(e.Posted).Nanoseconds())
				}
				msg := e.Message

				var hadQuorum bool
				if !c.noGossip {
					// check if we have quorum for message type for this round
					hadQuorum = c.quorumFor(msg.Code(), msg.R(), msg.Value())
				}

				if err := c.handleMsg(ctx, msg); err != nil {
					c.logger.Debug("MessageEvent payload failed", "err", err)
					// filter errors which needs remote peer disconnection
					if shouldDisconnectSender(err) {
						tryDisconnect(e.ErrCh, err)
					}
					// we still want to gossip old round messages
					if !errors.Is(err, constants.ErrOldRoundMessage) {
						break
					}
				}

				if !c.noGossip {
					if !hadQuorum {
						// if we did not have quorum and we reached it now
						// gossip the (complex) aggregate with quorum to everyone instead of the current message
						hasQuorum := c.quorumFor(msg.Code(), msg.R(), msg.Value())
						if hasQuorum {
							c.GossipComplexAggregate(msg.Code(), msg.R(), msg.Value())
							recordMessageProcessingTime(msg.Code(), start)
							break // do not gossip single message, only complex aggregate
						}
					}

					// gossip message. We should arrive here only if we did not already gossip a complex aggregate
					go c.backend.Gossip(c.CommitteeSet().Committee(), msg)
					recordMessageProcessingTime(msg.Code(), start)
				}
			case backlogMessageEvent:
				// TODO: should we check for disconnection also here for future round msgs?
				// need probably to store the errCh? verify if possible.

				msg := e.msg

				var hadQuorum bool
				if !c.noGossip {
					// check if we have quorum for message type for this round
					hadQuorum = c.quorumFor(msg.Code(), msg.R(), msg.Value())
				}

				c.logger.Debug("Handling consensus backlog event")
				if err := c.handleMsg(ctx, msg); err != nil {
					c.logger.Debug("BacklogEvent message handling failed", "err", err)
					continue
				}

				if !c.noGossip {
					if !hadQuorum {
						// if we did not have quorum and we reached it now
						// gossip the (complex) aggregate with quorum to everyone instead of the current message
						hasQuorum := c.quorumFor(msg.Code(), msg.R(), msg.Value())
						if hasQuorum {
							c.GossipComplexAggregate(msg.Code(), msg.R(), msg.Value())
							recordMessageProcessingTime(msg.Code(), start)
							break // do not gossip single message, only complex aggregate
						}
					}

					// gossip message. We should arrive here only if we did not already gossip a complex aggregate
					go c.backend.Gossip(c.CommitteeSet().Committee(), msg)
					recordMessageProcessingTime(msg.Code(), start)
				}
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
		case _, ok := <-c.committedCh:
			if !ok {
				break eventLoop
			}
			c.precommiter.HandleCommit(ctx)
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
	c.backend.AskSync(c.committee.Committee())

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
				c.backend.AskSync(c.committee.Committee())
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
		// TODO: currently old height messages are send directly to the FD, but this check is still needed due to potential TOCTOU race conditions
		// Moreover, I am still wondering if it would be useful to gossip old height messages, as they could be useful for accountability
		c.logger.Debug("ignoring stale consensus message", "msg", msg.String(), "height", c.Height().Uint64())
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

		// update future power
		_, ok := c.futurePower[r]
		if !ok {
			c.futurePower[r] = message.NewAggregatedPower()
		}
		switch m := msg.(type) {
		case *message.Propose:
			c.futurePower[r].Set(m.SignerIndex(), m.Power())
		case *message.Prevote, *message.Precommit:
			for index, power := range m.(message.Vote).Signers().Powers() {
				c.futurePower[r].Set(index, power)
			}
		}
		c.futureRoundLock.Unlock()

		go c.SendEvent(events.FuturePowerChangeEvent{Height: c.Height().Uint64(), Round: r})

		c.roundSkipCheck(ctx, r)
	}

	return err
}

func tryDisconnect(errorCh chan<- error, err error) {
	// errorCh can be nil in case the message is:
	// 1. an aggregated vote (non-complex)
	// 2. a locally created message
	// In both cases no error that causes disconnection can be raised anyways.
	if errorCh == nil {
		return
	}

	select {
	case errorCh <- err:
	default: // do nothing
	}
}
