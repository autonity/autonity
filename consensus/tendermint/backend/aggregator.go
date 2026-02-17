package backend

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
)

const (
	aggregationPeriod            = 150 * time.Millisecond
	oldMessagesAggregationPeriod = 2 * time.Second
	oldMessagesStatsPeriod       = 2 * time.Second
	numComputeWorkers            = 8
)

// aggregator metrics
var (
	ProposalBg  = metrics.NewRegisteredBufferedGauge("aggregator/proposal", nil, nil)                         // time it takes to process a proposal as received by backend
	PrevoteBg   = metrics.NewRegisteredBufferedGauge("aggregator/prevote", nil, metrics.GetIntPointer(200))   // time it takes to process a prevote as received by backend
	PrecommitBg = metrics.NewRegisteredBufferedGauge("aggregator/precommit", nil, metrics.GetIntPointer(200)) // time it takes to process a precommit as received by backend

	// packet meters
	ProposalPackets  = metrics.NewRegisteredMeter("aggregator/proposal/packets", nil)  //nolint:goconst
	PrevotePackets   = metrics.NewRegisteredMeter("aggregator/prevote/packets", nil)   //nolint:goconst
	PrecommitPackets = metrics.NewRegisteredMeter("aggregator/precommit/packets", nil) //nolint:goconst

	BatchesBg                  = metrics.NewRegisteredBufferedGauge("aggregator/batches", nil, metrics.GetIntPointer(100))          // size of batches (aggregated together with a single fastAggregateVerify)
	InvalidBg                  = metrics.NewRegisteredBufferedGauge("aggregator/invalid", nil, metrics.GetIntPointer(100))          // number of invalid sigs
	BackendAggregatorTransitBg = metrics.NewRegisteredBufferedGauge("aggregator/backend/transit", nil, metrics.GetIntPointer(1000)) // measures time for message passing from backend to aggregator
)

func recordMessageProcessingTime(code uint8, start time.Time) {
	if !metrics.Enabled() {
		return
	}
	switch code {
	case message.ProposalCode:
		ProposalBg.Add(time.Since(start).Nanoseconds())
		ProposalPackets.Mark(1)
	case message.PrevoteCode:
		PrevoteBg.Add(time.Since(start).Nanoseconds())
		PrevotePackets.Mark(1)
	case message.PrecommitCode:
		PrecommitBg.Add(time.Since(start).Nanoseconds())
		PrecommitPackets.Mark(1)
	}
}

type eventBuilder func(msg message.Msg, errCh chan<- error, sender common.Address, disseminated bool) interface{}

// function to create the event for current height messages (they get picked up by Core and by the FD)
func currentHeightEventBuilder(msg message.Msg, errCh chan<- error, sender common.Address, disseminated bool) interface{} {
	return events.NewMessageEvent(msg, errCh, sender, time.Now(), disseminated)
}

// function to create the event for old height messages (they get picked up only by the FD)
func oldHeightEventBuilder(msg message.Msg, errCh chan<- error, sender common.Address, _ bool) interface{} {
	return events.NewOldMessageEvent(msg, errCh, sender, time.Now())
}

func newAggregator(backend interfaces.Backend, core interfaces.Core, logger log.Logger) *aggregator {
	return &aggregator{
		backend:           backend,
		core:              core,
		staleMessages:     make(map[common.Hash][]events.UnverifiedMessageEvent),
		messages:          make(map[uint64]map[int64]*RoundInfo),
		logger:            logger,
		messagesFrom:      make(map[common.Address][]common.Hash),
		toIgnore:          make(map[common.Hash]struct{}),
		internalCoreCh:    make(chan events.MessageEventer, 1),
		internalFdCh:      make(chan events.MessageEventer, 1),
		internalBacklogCh: make(chan events.UnverifiedMessageEvent, 1000),
		computeWorkersCh:  make(chan events.UnverifiedMessageEvent, aggregatorMessageQueue),
		signerSetCache:    newAggregatorCache(),
	}
}

type RoundInfo struct {
	prevotes   map[common.Hash][]events.UnverifiedMessageEvent
	precommits map[common.Hash][]events.UnverifiedMessageEvent
}

func NewRoundInfo() *RoundInfo {
	return &RoundInfo{
		prevotes:   make(map[common.Hash][]events.UnverifiedMessageEvent),
		precommits: make(map[common.Hash][]events.UnverifiedMessageEvent),
	}
}

type aggregator struct {
	backend interfaces.Backend
	core    interfaces.Core

	staleMessages map[common.Hash][]events.UnverifiedMessageEvent

	messages     map[uint64]map[int64]*RoundInfo
	messagesFrom map[common.Address][]common.Hash

	toIgnore map[common.Hash]struct{}

	signerSetCache *aggregatorCache
	afdDispatchCh  chan<- events.MessageEventer

	internalCoreCh    chan events.MessageEventer
	internalFdCh      chan events.MessageEventer
	internalBacklogCh chan events.UnverifiedMessageEvent // backlog for messages that were filtered by cache, but reinjected

	computeWorkersCh chan events.UnverifiedMessageEvent

	cancel context.CancelFunc
	wg     sync.WaitGroup
	logger log.Logger
}

func (a *aggregator) start(ctx context.Context) {
	a.logger.Info("Starting the aggregator routine")
	ctx, a.cancel = context.WithCancel(ctx)
	// if the aggregator was previously stopped, these will be nil, we have to recreate them
	if a.internalCoreCh == nil {
		a.internalCoreCh = make(chan events.MessageEventer, 1) // buffered to avoid deadlock
	}
	if a.internalFdCh == nil {
		a.internalFdCh = make(chan events.MessageEventer, 1) // buffered to avoid deadlock
	}
	if a.internalBacklogCh == nil {
		a.internalBacklogCh = make(chan events.UnverifiedMessageEvent, 1000) // buffered to avoid deadlock
	}
	if a.computeWorkersCh == nil {
		a.computeWorkersCh = make(chan events.UnverifiedMessageEvent, aggregatorMessageQueue)
	}
	a.wg.Add(numComputeWorkers)
	for i := 0; i < numComputeWorkers; i++ {
		go func() {
			defer a.wg.Done()
			a.computeWork(ctx)
		}()
	}
	a.wg.Add(1)
	go a.loop(ctx)
}

func (a *aggregator) computeWork(ctx context.Context) {
	for {
		select {
		case ev := <-a.computeWorkersCh:
			// trigger the lazy computations
			ev.Message.SignerKey()
			_, _ = ev.Message.Signature() // error will be captured later during validation
		case <-ctx.Done():
			return
		}
	}
}

func (a *aggregator) handleInvalidMessage(event events.UnverifiedMessageEvent, err error) {
	tryDisconnect(event.ErrCh, err)
	a.signerSetCache.invalidateVotes(event.Message.(message.Vote))
	for _, hash := range a.messagesFrom[event.Sender] {
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
		a.messages[h] = make(map[int64]*RoundInfo)
	}

	if _, ok := a.messages[h][r]; !ok {
		a.messages[h][r] = NewRoundInfo()
	}

	roundInfo := a.messages[h][r]

	switch c {
	case message.PrevoteCode:
		roundInfo.prevotes[v] = append(roundInfo.prevotes[v], e)
	case message.PrecommitCode:
		roundInfo.precommits[v] = append(roundInfo.precommits[v], e)
	}
}

func (a *aggregator) empty(h uint64, r int64) bool {
	if _, ok := a.messages[h]; !ok {
		return true
	}

	if _, ok := a.messages[h][r]; !ok {
		return true
	}
	return false
}

func (a *aggregator) processRound(h uint64, r int64) {
	if a.empty(h, r) {
		return
	}

	roundInfo := a.messages[h][r]

	nBatches := len(roundInfo.prevotes) + len(roundInfo.precommits)
	batches := make([][]events.UnverifiedMessageEvent, 0, nBatches)

	// batch prevotes
	for _, evs := range roundInfo.prevotes {
		batches = append(batches, evs)
	}

	// batch precommits
	for _, evs := range roundInfo.precommits {
		batches = append(batches, evs)
	}

	a.processBatches(batches, currentHeightEventBuilder)

	//clean up
	delete(a.messages[h], r)
}

func (a *aggregator) processVotes(h uint64, r int64, c uint8) {
	if a.empty(h, r) {
		return
	}

	roundInfo := a.messages[h][r]

	// fill up batches matrix
	switch c {
	case message.PrevoteCode:
		nBatches := len(roundInfo.prevotes)
		batches := make([][]events.UnverifiedMessageEvent, nBatches)

		i := 0
		for _, events := range roundInfo.prevotes {
			batches[i] = events
			i++
		}

		a.processBatches(batches, currentHeightEventBuilder)

		// clean up. Use `clear` instead of `make` so that we stop iterating over the map if we end up here from RoundChangeEvent
		clear(roundInfo.prevotes)
	case message.PrecommitCode:
		nBatches := len(roundInfo.precommits)
		batches := make([][]events.UnverifiedMessageEvent, nBatches)

		i := 0
		for _, events := range roundInfo.precommits {
			batches[i] = events
			i++
		}

		a.processBatches(batches, currentHeightEventBuilder)

		// clean up
		clear(roundInfo.precommits)
	default:
		a.logger.Crit("Unexpected code", "c", c)
	}

}

func (a *aggregator) processVotesFor(h uint64, r int64, c uint8, v common.Hash) {
	if a.empty(h, r) {
		return
	}

	roundInfo := a.messages[h][r]

	switch c {
	case message.PrevoteCode:
		batch, ok := roundInfo.prevotes[v]
		if !ok {
			return
		}
		a.processBatches([][]events.UnverifiedMessageEvent{batch}, currentHeightEventBuilder)

		// clean up
		delete(roundInfo.prevotes, v)
	case message.PrecommitCode:
		batch, ok := roundInfo.precommits[v]
		if !ok {
			return
		}
		a.processBatches([][]events.UnverifiedMessageEvent{batch}, currentHeightEventBuilder)
		// clean up
		delete(roundInfo.precommits, v)
	default:
		a.logger.Crit("Unexpected code", "c", c)
	}
}

func (a *aggregator) DispatchCoreEvents(_ context.Context) {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		for {
			select {
			case event, ok := <-a.internalCoreCh:
				if !ok {
					a.logger.Warn("Aggregator internal core channel closed, stopping dispatching to core")

					return
				}
				// This is the only place that blocks for the Core
				if ev, ok := event.(events.MessageEvent); ok { // only new message events are send to core
					a.backend.DispatchToCore(ev)
				}
			}
		}
	}()
}

func (a *aggregator) DispatchFaultDetectorEvents(_ context.Context) {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		for {
			select {
			case event, ok := <-a.internalFdCh:
				if !ok {
					a.logger.Warn("Aggregator internal fault detector channel closed, stopping dispatching to FD")
					return
				}
				// This is the only place that blocks the fault detector
				a.backend.DispatchToFD(event)
			}
		}
	}()
}

func (a *aggregator) filterSkipped(evs []events.UnverifiedMessageEvent) []events.UnverifiedMessageEvent {
	filtered := make([]events.UnverifiedMessageEvent, 0, len(evs))
	for _, e := range evs {
		m := e.Message
		// skip messages to be ignored
		if a.toSkip(m) {
			continue
		}
		filtered = append(filtered, e)
	}
	return filtered
}

// a batch is a set of messages for same (height,round,code,value) ---> can be aggregated using FastAggregateVerify
func (a *aggregator) processBatches(batches [][]events.UnverifiedMessageEvent, eventer eventBuilder) {
	if len(batches) == 0 {
		return
	}
	for _, batch := range batches {
		batch = a.filterSkipped(batch) // filter out messages that are to be skipped
		a.processAndValidateBatch(batch, eventer)
	}
}

func (a *aggregator) processAndValidateBatch(batch []events.UnverifiedMessageEvent, eventer eventBuilder) {
	if len(batch) == 0 {
		return
	}
	if metrics.Enabled() {
		BatchesBg.Add(int64(len(batch)))
	}
	candidates := make([]events.UnverifiedMessageEvent, 0, len(batch))
	publicKeys := make([]blst.PublicKey, 0, len(batch))
	signatures := make([]blst.Signature, 0, len(batch))
	reInjectFiltered := false

	defer func() {
		if reInjectFiltered {
			// we have a backup of filtered messages, re-inject them in the backlog channel
			a.reInjectFilteredMessages(batch[0].Message.(message.Vote))
		}
	}()
	for _, event := range batch {
		m := event.Message.(message.Vote)
		publicKey := m.SignerKey()
		signature, err := m.Signature()
		if err != nil {
			a.logger.Debug("Signature decoding failed for message", "peer", event.Sender, "err", err)
			a.handleInvalidMessage(event, err)
			reInjectFiltered = true
			continue
		}
		candidates = append(candidates, event)
		publicKeys = append(publicKeys, publicKey)
		signatures = append(signatures, signature)
	}

	if len(candidates) == 0 {
		return
	}

	hash := batch[0].Message.SignatureInput()
	if blst.FastAggregateVerifyBatch(signatures, publicKeys, hash) {
		// all signatures are valid
		validVotes := make([]message.Vote, len(candidates))
		for i, event := range candidates {
			validVotes[i] = event.Message.(message.Vote)
		}
		a.aggregateAndDispatch(validVotes, eventer)
		return
	}

	// at least one of the signatures is invalid, find at which index
	reInjectFiltered = true
	invalids := blst.FindInvalid(signatures, publicKeys, hash)
	if metrics.Enabled() {
		InvalidBg.Add(int64(len(invalids)))
	}
	invalidSet := make(map[int]struct{}, len(invalids))
	for _, idx := range invalids {
		invalidSet[int(idx)] = struct{}{} // #nosec
	}
	validVotes := make([]message.Vote, 0, len(candidates)-len(invalids))
	for i, event := range candidates {
		if _, isInvalid := invalidSet[i]; !isInvalid {
			validVotes = append(validVotes, event.Message.(message.Vote))
			continue
		}
		a.logger.Info("Received invalid bls signature from", "peer", event.Sender)
		a.handleInvalidMessage(event, message.ErrBadSignature)
	}
	a.aggregateAndDispatch(validVotes, eventer)
}

func (a *aggregator) reInjectFilteredMessages(vote message.Vote) {
	filtered := a.signerSetCache.emptyFiltered(vote.H(), vote.R(), vote.Code(), vote.Value())
	warned := false
	for i, e := range filtered {
		select {
		case a.internalBacklogCh <- e:
		default:
			if !warned {
				log.Warn("Aggregator backlog channel is full, dropping messages", "number of messages dropped", len(filtered)-i,
					"height", vote.H(), "round", vote.R(), "code", vote.Code(), "value", vote.Value())
				warned = true
			}
		}
	}
}

func (a *aggregator) aggregateAndDispatch(validVotes []message.Vote, eventer eventBuilder) {
	if len(validVotes) == 0 {
		return
	}
	switch validVotes[0].(type) {
	case *message.Prevote:
		aggregateVotes := message.AggregatePrevotes(validVotes)
		for _, aggregateVote := range aggregateVotes {
			a.signerSetCache.addVote(aggregateVote, stepDispatched)
			a.internalCoreCh <- eventer(aggregateVote, nil, a.backend.Address(), false).(events.MessageEventer)
			a.internalFdCh <- eventer(aggregateVote, nil, a.backend.Address(), false).(events.MessageEventer)
		}
	case *message.Precommit:
		aggregateVotes := message.AggregatePrecommits(validVotes)
		for _, aggregateVote := range aggregateVotes {
			a.signerSetCache.addVote(aggregateVote, stepDispatched)
			a.internalCoreCh <- eventer(aggregateVote, nil, a.backend.Address(), false).(events.MessageEventer)
			a.internalFdCh <- eventer(aggregateVote, nil, a.backend.Address(), false).(events.MessageEventer)
		}
	default:
		a.logger.Crit("messages being aggregated are not votes", "type", reflect.TypeOf(validVotes[0]))
	}
}

func (a *aggregator) processProposal(proposalEvent events.UnverifiedMessageEvent, eventer eventBuilder) {
	proposal := proposalEvent.Message
	a.signerSetCache.addEvent(proposalEvent, stepDispatched)
	a.internalCoreCh <- eventer(proposal, proposalEvent.ErrCh, proposalEvent.Sender, proposalEvent.Disseminated).(events.MessageEventer) // send to core
	a.internalFdCh <- eventer(proposal, proposalEvent.ErrCh, proposalEvent.Sender, proposalEvent.Disseminated).(events.MessageEventer)   // send to fault detector
}

func (a *aggregator) areAllSignersJailed(msg message.Msg, committee *types.Committee) bool {
	vote, ok := msg.(message.Vote)
	if !ok { // proposals are already checked for jailing in backend handler
		return false
	}
	// check if all signers of the message are jailed
	if a.backend.JailedCount() >= vote.Signers().Len() {
		allJailed := true
		vote.Signers().ForEachDistinctSigner(func(signerIndex int) {
			signer := committee.Members[signerIndex].Address
			if !a.backend.IsJailed(signer) {
				allJailed = false
			}
		})
		// unless all signers are jailed, we still process aggregates
		if allJailed {
			a.logger.Debug("Vote message contains only signatures from jailed validators, ignoring message", "signers", vote.Signers().String())
			return true
		}
	}
	return false
}

// assumes current or old round vote
// if add == true, the msg is saved in the aggregator.
// if add == false, the msg is not saved and only the power checks are done.
func (a *aggregator) handleVote(voteEvent events.UnverifiedMessageEvent, committee *types.Committee) {
	vote := voteEvent.Message.(message.Vote)
	height := vote.H()
	round := vote.R()
	code := vote.Code()
	value := vote.Value()

	a.saveMessage(voteEvent)

	if round > a.core.Round() {
		a.processFutureRound(committee, voteEvent.Message)
		return
	}

	quorum := bft.Quorum(committee.TotalVotingPower())
	//// check if we reached quorum voting power on a specific value
	votingPowerReceived := a.signerSetCache.presentPowerForValue(height, round, value, code, stepReceived)
	if votingPowerReceived.Cmp(quorum) >= 0 {
		votingPowerDispatched := a.signerSetCache.presentPowerForValue(height, round, value, code, stepDispatched)
		if votingPowerDispatched.Cmp(quorum) < 0 {
			// we have enough votes, but not enough dispatched, so process the votes
			a.processVotesFor(height, round, code, value)
			return
		}
	}

	//// check if we have reached a quorum of votes for a specific round (regardless of value)
	totalPowerReceived := a.signerSetCache.totalPowerForCode(height, round, code, stepReceived)
	if totalPowerReceived.Cmp(quorum) >= 0 {
		totalPowerDispatched := a.signerSetCache.totalPowerForCode(height, round, code, stepDispatched)
		if totalPowerDispatched.Cmp(quorum) < 0 {
			a.processVotes(height, round, code)
			return
		}
	}
	// blocking call to trigger lazy computations of signatures and pubkeys
	a.computeWorkersCh <- voteEvent
}

func (a *aggregator) toSkip(msg message.Msg) bool {
	_, ignore := a.toIgnore[msg.Hash()]
	return ignore
}

func (a *aggregator) processFutureRound(committee *types.Committee, msg message.Msg) {
	// check power of ALL messages for that round, including proposals
	// if > 1/3 (ie. bft.F) then process that round
	receivedPowerForRound := a.signerSetCache.totalPowerForRound(msg.H(), msg.R(), stepReceived)
	if receivedPowerForRound.Cmp(bft.F(committee.TotalVotingPower())) > 0 {
		dispatchedPowerForRound := a.signerSetCache.totalPowerForRound(msg.H(), msg.R(), stepDispatched)
		if dispatchedPowerForRound.Cmp(bft.F(committee.TotalVotingPower())) <= 0 {
			a.processRound(msg.H(), msg.R())
		}
	}
}

func (a *aggregator) handleEvent(event events.UnverifiedMessageEvent) {
	start := time.Now()
	msg := event.Message
	sender := event.Sender

	a.messagesFrom[sender] = append(a.messagesFrom[sender], msg.Hash())

	// NOTE: Aggregator and Core run asynchronously. The code needs to take into account that Core can change state at any point here.
	// This also implies that height checks still needs to be done in Core.
	coreHeight := a.core.Height().Uint64()
	committee, err := a.backend.CommitteeByHeight(msg.H())
	if err != nil {
		panic(fmt.Sprintf("cannot get committee of height: %d", msg.H()))
	}
	if a.signerSetCache.filter(event) {
		return // already processed a message with more signers
	}
	if a.areAllSignersJailed(msg, committee) {
		return
	}
	// mark committee size for the height to avoid any more calls to CommitteeByHeight
	a.signerSetCache.markCommittee(msg.H(), committee)
	a.signerSetCache.addEvent(event, stepReceived)

	// old height proposals are processed right away, while votes are placed in the stale messages map
	if msg.H() < coreHeight {
		if msg.Code() == message.ProposalCode {
			a.processProposal(event, oldHeightEventBuilder)
		} else {
			signatureInput := msg.SignatureInput()
			a.staleMessages[signatureInput] = append(a.staleMessages[signatureInput], event)
		}
		return // note: processing ends here
	}
	if msg.H() > coreHeight {
		// future messages are dealt with at backend peer handler level
		a.logger.Crit("future message in aggregator", "msgHeight", msg.H(), "coreHeight", coreHeight)
	}

	switch msg.(type) {
	case *message.Propose:
		a.processProposal(event, currentHeightEventBuilder)
		// if it was a future round proposal, check for round skip
		if msg.R() > a.core.Round() {
			a.processFutureRound(committee, msg)
		}
	case *message.Prevote, *message.Precommit:
		a.handleVote(event, committee)
	default:
		panic("unknown message type: " + reflect.TypeOf(msg).String())
	}
	recordMessageProcessingTime(msg.Code(), start)
}

func (a *aggregator) oldHeightStats() {
	if !a.logger.Enabled(context.Background(), log.LevelDebug) {
		return
	}

	stats := make(map[uint64]map[int64][3]int)
	for _, batch := range a.staleMessages {
		for _, event := range batch {
			height := event.Message.H()
			round := event.Message.R()
			code := event.Message.Code()

			if stats[height] == nil {
				stats[height] = make(map[int64][3]int)
			}

			counts := stats[height][round]
			counts[code]++
			stats[height][round] = counts
		}
	}

	sb := strings.Builder{}
	sb.Grow(len(stats) * 100)

	sb.WriteString("Stale message Statistics by Height and Round\n")
	for height, rounds := range stats {
		for round, counts := range rounds {
			fmt.Fprintf(&sb, "H: %d R: %d | ", height, round)
			if counts[message.PrevoteCode] > 0 {
				fmt.Fprintf(&sb, "prevotes: %d ", counts[message.PrevoteCode])
			}
			if counts[message.PrecommitCode] > 0 {
				fmt.Fprintf(&sb, "precommits: %d ", counts[message.PrecommitCode])
			}
			sb.WriteByte('\n')
		}
		sb.WriteString("-------------------------------------")
	}
	a.logger.Debug("Stale message statistics", "stats", sb.String())
}

func (a *aggregator) loop(ctx context.Context) {
	defer a.wg.Done()
	a.DispatchCoreEvents(ctx)
	a.DispatchFaultDetectorEvents(ctx)

	aggTicker := time.NewTicker(aggregationPeriod)
	defer aggTicker.Stop()
	oldMessagesTicker := time.NewTicker(oldMessagesAggregationPeriod)
	defer oldMessagesTicker.Stop()

	oldMessagesStatsTicker := time.NewTicker(oldMessagesStatsPeriod)
	defer oldMessagesStatsTicker.Stop()

	// channel where the aggregator will receive msgs from the backend handlers
	messageCh := a.backend.MessageCh()
loop:
	for {
		select {
		case event, ok := <-messageCh:
			if !ok {
				break loop
			}
			if metrics.Enabled() {
				BackendAggregatorTransitBg.Add(time.Since(event.Posted).Nanoseconds())
			}
			a.handleEvent(event)
		case event, ok := <-a.internalBacklogCh:
			// handle backlog messages that were filtered by cache, but reinjected
			if !ok {
				break loop
			}
			a.handleEvent(event)
		case <-aggTicker.C:
			coreHeight := a.core.Height().Uint64()

			// process all messages in the aggregator
			for h, roundMap := range a.messages {
				for r, roundInfo := range roundMap {
					// if old height messages, move them to the staleMessages data structure. They are not useful for Core anymore.
					if h < coreHeight {
						// prevotes
						for _, sameValueVotes := range roundInfo.prevotes {
							if len(sameValueVotes) == 0 {
								continue
							}
							signatureInput := sameValueVotes[0].Message.SignatureInput() // all votes have same (h,r,c,v)
							a.staleMessages[signatureInput] = append(a.staleMessages[signatureInput], sameValueVotes...)
						}
						// precommits
						for _, sameValueVotes := range roundInfo.precommits {
							if len(sameValueVotes) == 0 {
								continue
							}
							signatureInput := sameValueVotes[0].Message.SignatureInput() // all votes have same (h,r,c,v)
							a.staleMessages[signatureInput] = append(a.staleMessages[signatureInput], sameValueVotes...)
						}
						delete(a.messages[h], r)
						continue
					}
					// if current height, process them
					a.processRound(h, r)
				}
			}
			// cleanup
			clear(a.messages)
			clear(a.messagesFrom)
			clear(a.toIgnore)
			a.cleanUp(coreHeight)
		case <-oldMessagesTicker.C:
			a.logger.Trace("Processing stale messages in the aggregator")
			var batches [][]events.UnverifiedMessageEvent
			for _, batch := range a.staleMessages {
				batches = append(batches, batch)
			}
			a.processBatches(batches, oldHeightEventBuilder)
			clear(a.staleMessages)
		case <-oldMessagesStatsTicker.C:
			a.oldHeightStats()
		case <-ctx.Done():
			break loop
		}
	}
	close(a.internalCoreCh)
	close(a.internalFdCh)
	close(a.internalBacklogCh)
}

func (a *aggregator) cleanUp(coreHeight uint64) {
	minHeight := coreHeight

	// clean up stale messages that are older than core height
	for _, batch := range a.staleMessages {
		if batch[0].Message.H() < minHeight {
			minHeight = batch[0].Message.H()
		}
	}

	// clean up cache
	a.signerSetCache.pruneToHeight(minHeight)
}

func (a *aggregator) stop() {
	a.logger.Info("Stopping the aggregator routine")
	a.cancel()
	a.wg.Wait()
	a.computeWorkersCh = nil
	a.internalCoreCh = nil
	a.internalFdCh = nil
	a.internalBacklogCh = nil
}
