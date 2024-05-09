package backend

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"sync"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
)

const (
	aggregationPeriod            = 100 * time.Millisecond
	oldMessagesAggregationPeriod = 2 * time.Second
)

// aggregator metrics
var (
	// metrics for different aggregator flows
	MessageBg     = metrics.NewRegisteredBufferedGauge("aggregator/bg/message", nil)     // time it takes to process a single message as received by backend
	RoundBg       = metrics.NewRegisteredBufferedGauge("aggregator/bg/round", nil)       // time it takes to process a round change event
	PowerBg       = metrics.NewRegisteredBufferedGauge("aggregator/bg/power", nil)       // time it takes to process a power change event
	FuturePowerBg = metrics.NewRegisteredBufferedGauge("aggregator/bg/futurepower", nil) // time it takes to process a future power change event
	ProcessBg     = metrics.NewRegisteredBufferedGauge("aggregator/bg/process", nil)     // time it takes to process batches of messages
	BatchesBg     = metrics.NewRegisteredBufferedGauge("aggregator/bg/batches", nil)     // size of batches (aggregated together with a single fastAggregateVerify)
	InvalidBg     = metrics.NewRegisteredBufferedGauge("aggregator/bg/invalid", nil)     // number of invalid sigs
)

type eventerFn func(msg message.Msg, errCh chan<- error) interface{}

// function to create the event for current height messages (they get picked up by Core and by the FD)
func currentHeightEventer(msg message.Msg, errCh chan<- error) interface{} {
	return events.MessageEvent{
		Message: msg,
		ErrCh:   errCh,
	}
}

// function to create the event for old height messages (they get picked up only by the FD)
func oldHeightEventer(msg message.Msg, errCh chan<- error) interface{} {
	return events.OldMessageEvent{
		Message: msg,
		ErrCh:   errCh,
	}
}

func newAggregator(backend interfaces.Backend, core interfaces.Core, logger log.Logger) *aggregator {
	return &aggregator{
		backend:     backend,
		core:        core,
		oldMessages: make(map[common.Hash][]events.UnverifiedMessageEvent),
		messages:    make(map[uint64]map[int64]map[uint8]map[common.Hash][]events.UnverifiedMessageEvent),
		logger:      logger,
		votesFrom:   make(map[common.Address][]common.Hash),
		toIgnore:    make(map[common.Hash]struct{}),
	}
}

type aggregator struct {
	backend interfaces.Backend
	core    interfaces.Core
	// old height messages
	oldMessages map[common.Hash][]events.UnverifiedMessageEvent
	// TODO(lorenzo) might be worth to re-use the set logic that is used in Core (without locks)
	// map[height][round][code][value][]msgs
	messages  map[uint64]map[int64]map[uint8]map[common.Hash][]events.UnverifiedMessageEvent
	votesFrom map[common.Address][]common.Hash
	toIgnore  map[common.Hash]struct{}
	// TODO(lorenzo) can one sub starve the other?
	sub     *event.TypeMuxSubscription // backend events
	coreSub *event.TypeMuxSubscription // core events
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	logger  log.Logger
}

func (a *aggregator) start(ctx context.Context) {
	a.logger.Info("Starting the aggregator routine")
	a.sub = a.backend.Subscribe(events.UnverifiedMessageEvent{})
	a.coreSub = a.backend.Subscribe(events.RoundChangeEvent{}, events.PowerChangeEvent{}, events.FuturePowerChangeEvent{})
	ctx, a.cancel = context.WithCancel(ctx)
	a.wg.Add(1)
	go a.loop(ctx)
}

func (a *aggregator) handleInvalidMessage(errorCh chan<- error, err error, p2pSender common.Address) {
	tryDisconnect(errorCh, err)
	for _, hash := range a.votesFrom[p2pSender] {
		a.toIgnore[hash] = struct{}{}
	}
}

func tryDisconnect(errorCh chan<- error, err error) {
	select {
	case errorCh <- err:
	default: // do nothing
	}
}

func (a *aggregator) saveMessage(e events.UnverifiedMessageEvent) {
	h := e.Message.H()
	r := e.Message.R()
	c := e.Message.Code()
	v := e.Message.Value()

	if _, ok := a.messages[h]; !ok {
		a.messages[h] = make(map[int64]map[uint8]map[common.Hash][]events.UnverifiedMessageEvent)
	}

	if _, ok := a.messages[h][r]; !ok {
		a.messages[h][r] = make(map[uint8]map[common.Hash][]events.UnverifiedMessageEvent)
	}

	if _, ok := a.messages[h][r][c]; !ok {
		a.messages[h][r][c] = make(map[common.Hash][]events.UnverifiedMessageEvent)
	}

	a.messages[h][r][c][v] = append(a.messages[h][r][c][v], e)
}

func (a *aggregator) power(h uint64, r int64) *big.Int {
	return a._power(h, r)
}

func (a *aggregator) votesPower(h uint64, r int64, code uint8) *big.Int {
	return a._power(h, r, code)
}

func (a *aggregator) votesPowerFor(h uint64, r int64, code uint8, v common.Hash) *big.Int {
	return a._power(h, r, code, v)
}

func parse(codeAndValue []interface{}) (uint8, common.Hash, bool, bool) {
	filterByCode := false
	filterByValue := false
	var c uint8
	var v common.Hash
	switch len(codeAndValue) {
	case 2:
		v = codeAndValue[1].(common.Hash)
		filterByValue = true
		fallthrough // if we filter by value, we always also filter by code
	case 1:
		c = codeAndValue[0].(uint8)
		filterByCode = true
	case 0:
		break
	default:
		panic(fmt.Sprintf("Invalid number of variadic parameters: %d", len(codeAndValue)))
	}
	return c, v, filterByCode, filterByValue
}

// TODO(lorenzo) implement some sort of caching
func (a *aggregator) _power(h uint64, r int64, codeAndValue ...interface{}) *big.Int {
	batches := a.batches(h, r, codeAndValue)

	var messages []message.Msg

	for _, batch := range batches {
		for _, e := range batch {
			messages = append(messages, e.Message)
		}
	}

	return message.Power(messages)
}

func (a *aggregator) batches(h uint64, r int64, codeAndValue []interface{}) [][]events.UnverifiedMessageEvent {
	if _, ok := a.messages[h]; !ok {
		return nil
	}

	if _, ok := a.messages[h][r]; !ok {
		return nil
	}

	// deal with variadic arguments if present
	c, v, filterByCode, filterByValue := parse(codeAndValue)

	batches := make([][]events.UnverifiedMessageEvent, 0)

	switch {
	case !filterByCode && !filterByValue:
		for _, valueMap := range a.messages[h][r] {
			for _, events := range valueMap {
				batches = append(batches, events)
			}
		}
	case filterByCode && !filterByValue:
		if _, ok := a.messages[h][r][c]; !ok {
			return nil
		}
		for _, events := range a.messages[h][r][c] {
			batches = append(batches, events)
		}
	case filterByCode && filterByValue:
		if _, ok := a.messages[h][r][c]; !ok {
			return nil
		}
		if _, ok := a.messages[h][r][c][v]; !ok {
			return nil
		}
		batches = append(batches, a.messages[h][r][c][v])
	case !filterByCode && filterByValue:
		a.logger.Crit("Trying to filter by value without filtering by code")
	}

	return batches
}

func (a *aggregator) process(h uint64, r int64, codeAndValue ...interface{}) {
	start := time.Now()
	// a batch is a set of messages for same (height,round,code,value) ---> can be aggregated using FastAggregateVerify
	batches := a.batches(h, r, codeAndValue)
	a.processBatches(batches, currentHeightEventer)

	// clean up
	c, v, filterByCode, filterByValue := parse(codeAndValue)

	switch {
	case !filterByCode && !filterByValue:
		delete(a.messages[h], r)
	case filterByCode && !filterByValue:
		delete(a.messages[h][r], c)
	case filterByCode && filterByValue:
		delete(a.messages[h][r][c], v)
	case !filterByCode && filterByValue:
		a.logger.Crit("Trying to filter by value without filtering by code")
	}
	ProcessBg.Add(time.Now().Sub(start).Nanoseconds())
}

func (a *aggregator) processBatches(batches [][]events.UnverifiedMessageEvent, eventer eventerFn) {
	if len(batches) == 0 {
		return
	}

	processed := 0 // messages that go in the aggregator
	sent := 0      // messages that out of the aggregator (to Core and FD as valid msgs)
	for _, batch := range batches {
		if len(batch) == 0 {
			continue
		}
		BatchesBg.Add(int64(len(batch)))
		processed += len(batch)

		// if batch of proposals, validate them individually
		if batch[0].Message.Code() == message.ProposalCode {
			for _, proposalEvent := range batch {
				if a.toSkip(proposalEvent.Message) {
					continue
				}
				a.processProposal(proposalEvent, eventer)
				sent++
			}
			continue
		}

		var publicKeys []blst.PublicKey
		var signatures []blst.Signature
		var messages []message.Vote
		var p2pSenders []common.Address
		var errChs []chan<- error

		for _, e := range batch {
			m := e.Message
			// skip messages to be ignored or that are already in core
			if a.toSkip(m) {
				continue
			}

			messages = append(messages, m.(message.Vote))
			publicKeys = append(publicKeys, m.SenderKey())
			signatures = append(signatures, m.Signature())
			p2pSenders = append(p2pSenders, e.P2pSender)
			errChs = append(errChs, e.ErrCh)
		}

		// if all messages in the batch got skipped, move to the next batch
		if len(signatures) == 0 {
			continue
		}

		aggregateSignature := blst.Aggregate(signatures)
		hash := batch[0].Message.SignatureInput()
		valid := aggregateSignature.FastAggregateVerify(publicKeys, hash)

		var validVotes []message.Vote
		var invalids []uint

		if !valid {
			// at least one of the signatures is invalid, find at which index
			invalids = blst.FindInvalid(signatures, publicKeys, hash)

			// remove invalid messages and sent the rest of the batch
			// NOTE: the following loop relies on blst.FindInvalid returning invalid indexes sorted according to ascending order
			j := 0
			for i, msg := range messages {
				if j < len(invalids) && uint(i) == invalids[j] {
					j++
					continue
				}
				validVotes = append(validVotes, msg)
			}
		} else {
			// all messages are valid
			validVotes = messages
		}

		sent += len(validVotes)

		if len(validVotes) > 0 {
			// dispatch messages to core and FD
			// repetitive code but I didn't find a way to declare aggregateVotes so that it works both with prevote and preccomit
			switch validVotes[0].(type) {
			case *message.Prevote:
				aggregateVotes := message.AggregatePrevotesSimple(validVotes)
				for _, aggregateVote := range aggregateVotes {
					go a.backend.Post(eventer(aggregateVote, nil)) //TODO(lorenzo) refinements, do we add an errCh here?
				}
			case *message.Precommit:
				aggregateVotes := message.AggregatePrecommitsSimple(validVotes)
				for _, aggregateVote := range aggregateVotes {
					go a.backend.Post(eventer(aggregateVote, nil)) //TODO(lorenzo) refinements, do we add an errCh here?
				}
			default:
				a.logger.Crit("messages being aggregated are not votes", "type", reflect.TypeOf(validVotes[0]))
			}
		}

		// disconnect validators who sent us invalid votes at p2p layer and ignore the msgs coming from them
		InvalidBg.Add(int64(len(invalids)))
		for _, index := range invalids {
			a.logger.Info("Received invalid bls signature from", "peer", p2pSenders[index])
			a.handleInvalidMessage(errChs[index], message.ErrBadSignature, p2pSenders[index])
		}
	}
	a.logger.Debug("Aggregator processed messages", "processed", processed, "sent", sent)
}

func (a *aggregator) processProposal(proposalEvent events.UnverifiedMessageEvent, eventer eventerFn) {
	proposal := proposalEvent.Message
	if err := proposal.Validate(); err != nil {
		a.handleInvalidMessage(proposalEvent.ErrCh, err, proposalEvent.P2pSender)
		return
	}
	go a.backend.Post(eventer(proposal, proposalEvent.ErrCh))
}

// assumes current or old round vote
func (a *aggregator) processVote(voteEvent events.UnverifiedMessageEvent, quorum *big.Int) {
	vote := voteEvent.Message.(message.Vote)
	errCh := voteEvent.ErrCh
	p2pSender := voteEvent.P2pSender

	height := vote.H()
	round := vote.R()
	code := vote.Code()
	value := vote.Value()

	// complex aggregates always carry quorum (enforced at PreValidate)
	// if we do not already have quorum in Core, process right away
	coreVotesForPower := a.core.VotesPowerFor(height, round, code, value)
	coreVotesPower := a.core.VotesPower(height, round, code)
	if vote.Senders().IsComplex() && (coreVotesForPower.Cmp(quorum) < 0 || coreVotesPower.Cmp(quorum) < 0) {
		if err := vote.Validate(); err != nil {
			a.handleInvalidMessage(errCh, err, p2pSender)
			return
		}
		go a.backend.Post(currentHeightEventer(voteEvent.Message, errCh))
		return
	}

	// we are processing an individual vote or a simple aggregate vote
	a.saveMessage(voteEvent)

	// check if we reached quorum voting power on a specific value
	corePower := a.core.VotesPowerFor(height, round, code, value)
	aggregatorPower := a.votesPowerFor(height, round, code, value)
	if corePower.Cmp(quorum) < 0 && corePower.Add(corePower, aggregatorPower).Cmp(quorum) >= 0 {
		a.process(vote.H(), vote.R(), code, value)
		return
	}

	// check if we reached quorum voting power in general
	corePower = a.core.VotesPower(vote.H(), vote.R(), code)
	aggregatorPower = a.votesPower(vote.H(), vote.R(), code)
	if corePower.Cmp(quorum) < 0 && corePower.Add(corePower, aggregatorPower).Cmp(quorum) >= 0 {
		a.process(vote.H(), vote.R(), code)
	}
}

// checks if a message is already in core
// TODO(lorenzo) performance, do something more efficient without iterating over all messages of Core
// Depending on how implemented, it might be possible to remove core.futureRoundLock
func (a *aggregator) alreadyProcessed(msg message.Msg) bool {
	for _, m := range a.core.CurrentHeightMessages() {
		if msg.Hash() == m.Hash() {
			return true
		}
	}
	return false
}

func (a *aggregator) toSkip(msg message.Msg) bool {
	_, ignore := a.toIgnore[msg.Hash()]
	if ignore || a.alreadyProcessed(msg) {
		return true
	}
	return false
}

//TODO(lorenzo) analyze proposal flow. Can we avoid having equivocated (or at least duplicated) proposals buffered in the aggregator?

func (a *aggregator) loop(ctx context.Context) {
	defer a.wg.Done()

	ticker := time.NewTicker(aggregationPeriod)
	defer ticker.Stop()
	oldMessagesTicker := time.NewTicker(oldMessagesAggregationPeriod)
	defer oldMessagesTicker.Stop()

loop:
	for {
		select {
		case ev, ok := <-a.sub.Chan():
			start := time.Now()
			if !ok {
				break loop
			}
			event := ev.Data.(events.UnverifiedMessageEvent)
			msg := event.Message
			p2pSender := event.P2pSender

			a.votesFrom[p2pSender] = append(a.votesFrom[p2pSender], msg.Hash())

			// NOTE: Aggregator and Core run asynchronously. The code needs to take into account that Core can change state at any point here.
			// This also implies that height checks still needs to be done in Core.
			coreHeight := a.core.Height().Uint64()
			if msg.H() < coreHeight {
				a.logger.Debug("Storing old height message in the aggregator", "msgHeight", msg.H(), "coreHeight", coreHeight)
				signatureInput := msg.SignatureInput()
				a.oldMessages[signatureInput] = append(a.oldMessages[signatureInput], event)
				MessageBg.Add(time.Now().Sub(start).Nanoseconds())
				break
			}
			if msg.H() > coreHeight {
				// future messages are dealt with at backend peer handler level
				a.logger.Crit("future message in aggregator", "msgHeight", msg.H(), "coreHeight", coreHeight)
			}

			// if message already in Core, drop it
			if a.alreadyProcessed(msg) {
				a.logger.Debug("Discarding msg, already processed in Core")
				MessageBg.Add(time.Now().Sub(start).Nanoseconds())
				break
			}

			header := a.backend.BlockChain().GetHeaderByNumber(msg.H() - 1)
			if header == nil {
				a.logger.Crit("cannot fetch header for non-future message", "headerHeight", msg.H()-1, "coreHeight", a.core.Height().Uint64())
			}
			quorum := bft.Quorum(header.TotalVotingPower())

			coreRound := a.core.Round()
			if msg.R() > coreRound {
				// NOTE: here we could be buffering a proposal for future round, or a complex vote aggregate.
				a.saveMessage(event)
				// check if power is enough for a round skip
				aggregatorPower := a.power(msg.H(), msg.R())
				corePower := a.core.Power(msg.H(), msg.R())
				if aggregatorPower.Add(aggregatorPower, corePower).Cmp(bft.F(header.TotalVotingPower())) > 0 {
					a.logger.Debug("Processing future round messages due to possible round skip", "height", msg.H(), "round", msg.R(), "coreRound", coreRound)
					a.process(msg.H(), msg.R())
				}
				MessageBg.Add(time.Now().Sub(start).Nanoseconds())
				break
			}

			// current or old round here
			switch msg.(type) {
			// if proposal, verify right away
			case *message.Propose:
				a.processProposal(event, currentHeightEventer)
			case *message.Prevote, *message.Precommit:
				a.processVote(event, quorum)
			default:
				a.logger.Crit("unknown message type arrived in aggregator")
			}
			MessageBg.Add(time.Now().Sub(start).Nanoseconds())
		case ev, ok := <-a.coreSub.Chan():
			start := time.Now()
			if !ok {
				break loop
			}
			switch e := ev.Data.(type) {
			case events.RoundChangeEvent:
				/* a round change happened in Core
				* messages that we had buffered as future round might now be current round, therefore:
				* 1. process right away proposals and complex aggregates
				* 2. re-do quorum checks on individual votes and simple aggregates
				 */
				//TODO(lorenzo) possibly I can also move messages that became old to the oldMessages map. However not really sure if worth it.
				height := e.Height
				round := e.Round

				if _, ok := a.messages[height]; !ok {
					RoundBg.Add(time.Now().Sub(start).Nanoseconds())
					break
				}

				if _, ok := a.messages[height][round]; !ok {
					RoundBg.Add(time.Now().Sub(start).Nanoseconds())
					break
				}

				// process proposals
				if valueMap, ok := a.messages[height][round][message.ProposalCode]; ok {
					for _, proposals := range valueMap {
						for _, proposal := range proposals {
							a.processProposal(proposal, currentHeightEventer)
						}
					}
					delete(a.messages[height][round], message.ProposalCode)
				}

				header := a.backend.BlockChain().GetHeaderByNumber(height - 1)
				if header == nil {
					a.logger.Crit("cannot fetch header for non-future height", "height", height-1)
				}
				quorum := bft.Quorum(header.TotalVotingPower())

				for _, valueMap := range a.messages[height][round] {
					// only prevotes or precommits here, we previously processed and deleted proposals
					for _, evs := range valueMap {
						for _, e := range evs {
							//TODO(lorenzo) possible problem here, processVote might call `process` which is going to delete stuff in the map
							// generally is fine but double check if it causes issue in this case
							a.processVote(e, quorum)
						}
					}
				}
				RoundBg.Add(time.Now().Sub(start).Nanoseconds())
			case events.PowerChangeEvent:
				// a power change happened in Core: re-do quorum checks on individual votes and simple aggregates
				height := e.Height
				round := e.Round
				code := e.Code
				value := e.Value

				header := a.backend.BlockChain().GetHeaderByNumber(height - 1)
				if header == nil {
					panic("cannot fetch header for non-future message")
				}
				quorum := bft.Quorum(header.TotalVotingPower())

				if _, ok := a.messages[height]; !ok {
					PowerBg.Add(time.Now().Sub(start).Nanoseconds())
					break
				}

				if _, ok := a.messages[height][round]; !ok {
					PowerBg.Add(time.Now().Sub(start).Nanoseconds())
					break
				}

				if _, ok := a.messages[height][round][code]; !ok {
					PowerBg.Add(time.Now().Sub(start).Nanoseconds())
					break
				}

				if evs, ok := a.messages[height][round][code][value]; ok {
					if len(evs) > 1 {
						// processing one vote for the value for which power changed is enough to do all necessary checks
						a.processVote(evs[0], quorum)
					}
				}
				PowerBg.Add(time.Now().Sub(start).Nanoseconds())
			case events.FuturePowerChangeEvent:
				height := e.Height
				round := e.Round

				header := a.backend.BlockChain().GetHeaderByNumber(height - 1)
				if header == nil {
					a.logger.Crit("cannot fetch header for non-future message", "height", height-1)
				}

				// check in future round messages power, check again for round skip
				aggregatorPower := a.power(height, round)
				corePower := a.core.Power(height, round)
				if aggregatorPower.Add(aggregatorPower, corePower).Cmp(bft.F(header.TotalVotingPower())) > 0 {
					a.process(height, round)
				}
				FuturePowerBg.Add(time.Now().Sub(start).Nanoseconds())
			}
		case <-ticker.C:
			// process all messages in the aggregator
			for h, roundMap := range a.messages {
				for r, _ := range roundMap {
					a.process(h, r)
				}
			}
			// cleanup
			a.votesFrom = make(map[common.Address][]common.Hash)
			a.toIgnore = make(map[common.Hash]struct{})
		case <-oldMessagesTicker.C:
			batches := make([][]events.UnverifiedMessageEvent, len(a.oldMessages))
			i := 0
			for _, batch := range a.oldMessages {
				batches[i] = batch
				i++
			}
			a.processBatches(batches, oldHeightEventer)

			a.oldMessages = make(map[common.Hash][]events.UnverifiedMessageEvent)
		case <-ctx.Done():
			break loop
		}
	}
}

func (a *aggregator) stop() {
	a.logger.Info("Stopping the aggregator routine")
	a.cancel()
	a.wg.Wait()
}
