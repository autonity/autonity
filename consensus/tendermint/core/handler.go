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
	c.stateEventSub = c.backend.Subscribe(StateRequestEvent{})
	c.candidateBlockCh = make(chan events.NewCandidateBlockEvent, 1)
	c.committedCh = make(chan events.CommitEvent, 1)
	c.messageEventCh = make(chan events.MessageEventer, 1000) // todo(review) : channel size
	c.timeoutEventSub = c.backend.Subscribe(TimeoutEvent{})
}

// Unsubscribe all
func (c *Core) unsubscribeEvents() {
	c.stateEventSub.Unsubscribe()
	c.timeoutEventSub.Unsubscribe()
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
	case errors.Is(err, consensus.ErrFutureTimestampBlock):
		fallthrough
	case errors.Is(err, consensus.ErrPrunedAncestor):
		fallthrough
	case errors.Is(err, constants.ErrRedundantVote):
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
		quorum = c.messages.GetOrCreate(round).PrevotesPower(value).Cmp(c.CommitteeSet().Quorum()) >= 0
	case message.PrecommitCode:
		quorum = c.messages.GetOrCreate(round).PrecommitsPower(value).Cmp(c.CommitteeSet().Quorum()) >= 0
	}
	return quorum
}

func (c *Core) GossipComplexAggregate(code uint8, round int64, value common.Hash) {
	// We re-add the complex aggregate to the prevote set. If we would substitute the entire set with the complex aggregate,
	// there is a possibility of message loss (if we had multiple un-mergeable complex aggregates in the `messages`).
	// This loss would not harm consensus (we would still have quorum voting power), however it is better to keep all messages
	// in case we have to sync another peer. We can consider changing it only if it considerably harms performance.
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
	go c.livenessTrackerLoop(ctx)

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
		case ev, ok := <-c.messageEventCh:
			if !ok {
				break eventLoop
			}
			start := time.Now()
			// An event arrived, process content
			switch e := ev.(type) {
			case events.MessageEvent:
				if metrics.Enabled {
					AggregatorCoreTransitBg.Add(time.Since(e.Posted()).Nanoseconds())
				}
				msg := e.Message()

				if c.Height().Uint64() > msg.H() {
					// TODO: currently old height messages are send directly to the FD, but this check is still needed due to potential TOCTOU race conditions
					c.logger.Debug("mainEventLoop: ignoring stale consensus message", "msg type", msg.Code(), "core height", c.Height().Uint64(), "msgHeight", msg.H(), "msgRound", msg.R())
					break
				}
				var hadQuorum, hasQuorum bool
				if !c.noGossip {
					// check if we have quorum for message type for this round
					hadQuorum = c.quorumFor(msg.Code(), msg.R(), msg.Value())
				}
				needGossip := true
				var err error
				if e.Sender() == c.backend.Address() && !hadQuorum {
					go c.backend.Gossip(c.CommitteeSet().Committee(), msg)
					needGossip = false
				}

				if err = c.handleMsg(ctx, msg); err != nil {
					c.logger.Debug("MessageEvent payload failed", "err", err, "current Height", c.Height().Uint64(), "msg Height", msg.H(), "msg Round", msg.R())
					// filter errors which needs remote peer disconnection
					if shouldDisconnectSender(err) {
						tryDisconnect(e.ErrCh(), err)
						break
					}
					if errors.Is(err, constants.ErrFutureRoundMessage) && msg.Code() != message.ProposalCode && e.Sender() != c.backend.Address() {
						// immediately gossip future round votes
						go c.backend.Router().Forward(c.CommitteeSet().Committee(), msg, e.Sender())
						recordMessageProcessingTime(msg.Code(), start)
						break
					}
					// we still want to gossip old round messages and redundant votes
					if !errors.Is(err, constants.ErrOldRoundMessage) && !errors.Is(err, constants.ErrRedundantVote) {
						break
					}
				}

				// proposals are already gossiped in backend
				if msg.Code() == message.ProposalCode {
					recordMessageProcessingTime(msg.Code(), start)
					break
				}

				// valid message, mark liveness time unless it was redundant
				//if !errors.Is(err, constants.ErrRedundantVote) {
				c.syncState.setLastLivenessTime(time.Now())
				//}
				if !c.noGossip {
					if !hadQuorum {
						// if we did not have quorum and we reached it now
						// gossip the (complex) aggregate with quorum to everyone instead of the current message
						hasQuorum = c.quorumFor(msg.Code(), msg.R(), msg.Value())
						if hasQuorum {
							c.GossipComplexAggregate(msg.Code(), msg.R(), msg.Value())
							recordMessageProcessingTime(msg.Code(), start)
							break // do not gossip single message, only complex aggregate
						}
					}
					if !needGossip {
						// already gossiped
						break
					}

					if err != nil && errors.Is(err, constants.ErrOldRoundMessage) {
						go c.backend.SlowGossip(c.CommitteeSet().Committee(), msg)
					} else {
						// todo: refactor
						go c.backend.Router().Forward(c.CommitteeSet().Committee(), msg, e.Sender())
					}
				}
				recordMessageProcessingTime(msg.Code(), start)
			case backlogMessageEvent:
				// TODO: should we check for disconnection also here for future round msgs?
				// need probably to store the errCh? verify if possible.

				msg := e.Message()

				if c.Height().Uint64() > msg.H() {
					// TODO: currently old height messages are send directly to the FD, but this check is still needed due to potential TOCTOU race conditions
					c.logger.Debug("mainEventLoop: ignoring stale consensus message", "msg type", msg.Code(), "core height", c.Height().Uint64(), "msgHeight", msg.H(), "msgRound", msg.R())
					break
				}

				var hadQuorum, hasQuorum bool
				if !c.noGossip {
					// check if we have quorum for message type for this round
					hadQuorum = c.quorumFor(msg.Code(), msg.R(), msg.Value())
				}

				c.logger.Debug("Handling consensus backlog event")
				var err error
				if err = c.handleMsg(ctx, msg); err != nil {
					c.logger.Debug("BacklogEvent message handling failed", "err", err)
					// we still want to gossip old round messages and redundant votes
					if !errors.Is(err, constants.ErrOldRoundMessage) && !errors.Is(err, constants.ErrRedundantVote) {
						continue
					}
				}

				// valid message, mark liveness time unless it was redundant
				if !errors.Is(err, constants.ErrRedundantVote) {
					c.syncState.setLastLivenessTime(time.Now())
				}

				//// proposals are already gossiped in backend
				//if msg.Code() == message.ProposalCode {
				//	recordMessageProcessingTime(msg.Code(), start)
				//	continue
				//}
				if !c.noGossip {
					if !hadQuorum {
						// if we did not have quorum and we reached it now
						// gossip the (complex) aggregate with quorum to everyone instead of the current message
						hasQuorum = c.quorumFor(msg.Code(), msg.R(), msg.Value())
						if hasQuorum {
							c.GossipComplexAggregate(msg.Code(), msg.R(), msg.Value())
							recordMessageProcessingTime(msg.Code(), start)
							break // do not gossip single message, only complex aggregate
						}
					}
				}
				//if err != nil && errors.Is(err, constants.ErrOldRoundMessage) {
				//	go c.backend.SlowGossip(c.CommitteeSet().Committee(), msg)
				//} else {
				//	go c.backend.Router().Forward(c.CommitteeSet().Committee(), msg, e.Sender)
				//}
				recordMessageProcessingTime(msg.Code(), start)
			}
		case ev, _ := <-c.stateEventSub.Chan():
			c.handleStateDump(ev.Data.(StateRequestEvent))
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

// this method is responsible for:
// - tracking whether we are out of consensus sync
// - asking the network to send us the current consensus state if so
func (c *Core) livenessTrackerLoop(ctx context.Context) {
	defer func() {
		c.stopped <- struct{}{}
	}()

	// Ask for sync when the engine starts. Retry until sync succeeds or we are stopped
	for {
		err := c.backend.AskSync(c.committee.Committee(), c.createSyncMsg())
		if err == nil {
			break
		}
		select {
		case <-ctx.Done():
			c.logger.Debug("livenessTrackerLoop has been stopped before initial sync", "event", ctx.Err())
			return
		default:
			c.logger.Warn("Failed to ask initial consensus sync, retrying...", "err", err)
			time.Sleep(100 * time.Millisecond)
		}
	}

	ticker := time.NewTicker(constants.AskSyncInterval)
	defer ticker.Stop()

eventLoop:
	for {
		select {
		case <-ticker.C:

			elapsedTime := time.Since(c.syncState.getLastLivenessTime())
			currentSyncTimeout := c.syncState.getSyncTimeout()
			if elapsedTime < currentSyncTimeout {
				c.logger.Debug("Sync timeout not reached yet", "elapsed time", elapsedTime, "current timeout", currentSyncTimeout)
				continue
			}

			// no liveness for more than currentSyncTimeout --> askSync to the other nodes
			c.logger.Warn("⚠️ Consensus liveliness lost", "node", c.Address(), "height", c.Height(), "round", c.Round(), "step", c.Step())
			err := c.backend.AskSync(c.committee.Committee(), c.createSyncMsg())
			if err != nil {
				c.logger.Warn("Failed to ask consensus sync", "err", err)
				// will automatically retry at next iteration
			}
		case <-ctx.Done():
			c.logger.Debug("livenessTrackerLoop is stopped", "event", ctx.Err())
			break eventLoop
		}
	}
}

// SendEvent sends event to mux
func (c *Core) SendEvent(ev any) {
	switch ev := ev.(type) {
	case events.CoreEvent:
		c.eventCh <- ev
	default:
		c.backend.Post(ev)
	}
}

func (c *Core) handleMsg(ctx context.Context, msg message.Msg) error {

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

		c.SendEvent(events.NewFuturePowerChangeEvent(c.Height().Uint64(), r))

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

func redundancyError(didContribute bool) error {
	if didContribute {
		return nil
	}
	return constants.ErrRedundantVote
}
