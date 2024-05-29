package backend

import (
	"context"
	"math/big"
	"reflect"
	"sync"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
)

const (
	aggregationPeriod            = 150 * time.Millisecond
	oldMessagesAggregationPeriod = 2 * time.Second
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

	RoundBg                    = metrics.NewRegisteredBufferedGauge("aggregator/round", nil, nil)                                   // time it takes to process a round change event
	PowerBg                    = metrics.NewRegisteredBufferedGauge("aggregator/power", nil, nil)                                   // time it takes to process a power change event
	FuturePowerBg              = metrics.NewRegisteredBufferedGauge("aggregator/futurepower", nil, nil)                             // time it takes to process a future power change event
	BatchesBg                  = metrics.NewRegisteredBufferedGauge("aggregator/batches", nil, metrics.GetIntPointer(100))          // size of batches (aggregated together with a single fastAggregateVerify)
	InvalidBg                  = metrics.NewRegisteredBufferedGauge("aggregator/invalid", nil, metrics.GetIntPointer(100))          // number of invalid sigs
	BackendAggregatorTransitBg = metrics.NewRegisteredBufferedGauge("aggregator/backend/transit", nil, metrics.GetIntPointer(1000)) // measures time for message passing from backend to aggregator
)

func recordMessageProcessingTime(code uint8, start time.Time) {
	if !metrics.Enabled {
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

type eventBuilder func(msg message.Msg, errCh chan<- error) interface{}

// function to create the event for current height messages (they get picked up by Core and by the FD)
func currentHeightEventBuilder(msg message.Msg, errCh chan<- error) interface{} {
	return events.MessageEvent{
		Message: msg,
		ErrCh:   errCh,
		Posted:  time.Now(),
	}
}

// function to create the event for old height messages (they get picked up only by the FD)
func oldHeightEventBuilder(msg message.Msg, errCh chan<- error) interface{} {
	return events.OldMessageEvent{
		Message: msg,
		ErrCh:   errCh,
	}
}

// computes how much new voting power will the messages in the aggregator apport to core
func powerContribution(aggregatorSigners *big.Int, coreSigners *big.Int, committee types.Committee) *big.Int {
	contribution := message.Contribution(aggregatorSigners, coreSigners)
	if contribution.Cmp(common.Big0) == 0 {
		return new(big.Int) // no power contribution
	}
	// there is a contribution, compute how much
	contributionPower := new(big.Int)
	for i, member := range committee {
		if contribution.Bit(i) == 1 {
			contributionPower.Add(contributionPower, member.VotingPower)
		}
	}
	return contributionPower
}

func newAggregator(backend interfaces.Backend, core interfaces.Core, logger log.Logger) *aggregator {
	return &aggregator{
		backend:       backend,
		core:          core,
		staleMessages: make(map[common.Hash][]events.UnverifiedMessageEvent),
		messages:      make(map[uint64]map[int64]*RoundInfo),
		logger:        logger,
		messagesFrom:  make(map[common.Address][]common.Hash),
		toIgnore:      make(map[common.Hash]struct{}),
	}
}

type RoundInfo struct {
	proposals []events.UnverifiedMessageEvent

	prevotes         map[common.Hash][]events.UnverifiedMessageEvent
	prevotesPower    *message.AggregatedPower
	prevotesPowerFor map[common.Hash]*message.AggregatedPower

	precommits         map[common.Hash][]events.UnverifiedMessageEvent
	precommitsPower    *message.AggregatedPower
	precommitsPowerFor map[common.Hash]*message.AggregatedPower

	power *message.AggregatedPower // entire round power
}

func NewRoundInfo() *RoundInfo {
	return &RoundInfo{
		proposals:          make([]events.UnverifiedMessageEvent, 0),
		prevotes:           make(map[common.Hash][]events.UnverifiedMessageEvent),
		prevotesPower:      message.NewAggregatedPower(),
		prevotesPowerFor:   make(map[common.Hash]*message.AggregatedPower),
		precommits:         make(map[common.Hash][]events.UnverifiedMessageEvent),
		precommitsPower:    message.NewAggregatedPower(),
		precommitsPowerFor: make(map[common.Hash]*message.AggregatedPower),
		power:              message.NewAggregatedPower(),
	}
}

type aggregator struct {
	backend interfaces.Backend
	core    interfaces.Core

	staleMessages map[common.Hash][]events.UnverifiedMessageEvent
	messages      map[uint64]map[int64]*RoundInfo

	messagesFrom map[common.Address][]common.Hash
	toIgnore     map[common.Hash]struct{}

	coreSub *event.TypeMuxSubscription // core events
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	logger  log.Logger
}

func (a *aggregator) start(ctx context.Context) {
	a.logger.Info("Starting the aggregator routine")
	a.coreSub = a.backend.Subscribe(events.RoundChangeEvent{}, events.PowerChangeEvent{}, events.FuturePowerChangeEvent{})
	ctx, a.cancel = context.WithCancel(ctx)
	a.wg.Add(1)
	go a.loop(ctx)
}

func (a *aggregator) handleInvalidMessage(errorCh chan<- error, err error, sender common.Address) {
	tryDisconnect(errorCh, err)
	for _, hash := range a.messagesFrom[sender] {
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
	case message.ProposalCode:
		roundInfo.proposals = append(roundInfo.proposals, e)

		// update round power cache
		proposal := e.Message.(*message.Propose)
		roundInfo.power.Set(proposal.SignerIndex(), proposal.Power())
	case message.PrevoteCode:
		roundInfo.prevotes[v] = append(roundInfo.prevotes[v], e)

		_, ok := roundInfo.prevotesPowerFor[v]
		if !ok {
			roundInfo.prevotesPowerFor[v] = message.NewAggregatedPower()
		}

		// update power caches
		vote := e.Message.(message.Vote)
		for index, power := range vote.Signers().Powers() {
			roundInfo.power.Set(index, power)
			roundInfo.prevotesPower.Set(index, power)
			roundInfo.prevotesPowerFor[v].Set(index, power)
		}
	case message.PrecommitCode:
		roundInfo.precommits[v] = append(roundInfo.precommits[v], e)

		_, ok := roundInfo.precommitsPowerFor[v]
		if !ok {
			roundInfo.precommitsPowerFor[v] = message.NewAggregatedPower()
		}

		// update power caches
		vote := e.Message.(message.Vote)
		for index, power := range vote.Signers().Powers() {
			roundInfo.power.Set(index, power)
			roundInfo.precommitsPower.Set(index, power)
			roundInfo.precommitsPowerFor[v].Set(index, power)
		}
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

func (a *aggregator) power(h uint64, r int64) *message.AggregatedPower {
	if a.empty(h, r) {
		return message.NewAggregatedPower()
	}
	return a.messages[h][r].power.Copy() // return a copy as the aggregator is going to modify this value
}

func (a *aggregator) votesPower(h uint64, r int64, c uint8) *message.AggregatedPower {
	if a.empty(h, r) {
		return message.NewAggregatedPower()
	}

	roundInfo := a.messages[h][r]
	var power *message.AggregatedPower
	switch c {
	case message.PrevoteCode:
		power = roundInfo.prevotesPower.Copy()
	case message.PrecommitCode:
		power = roundInfo.precommitsPower.Copy()
	default:
		a.logger.Crit("Unexpected code", "c", c)
	}
	return power // return a copy as the aggregator is going to modify this value
}

func (a *aggregator) votesPowerFor(h uint64, r int64, c uint8, v common.Hash) *message.AggregatedPower {
	if a.empty(h, r) {
		return message.NewAggregatedPower()
	}

	roundInfo := a.messages[h][r]
	var power *message.AggregatedPower
	var ok bool // necessary to not override power declaration inside the switch
	switch c {
	case message.PrevoteCode:
		power, ok = roundInfo.prevotesPowerFor[v]
	case message.PrecommitCode:
		power, ok = roundInfo.precommitsPowerFor[v]
	default:
		a.logger.Crit("Unexpected code", "c", c)
	}

	if !ok {
		return message.NewAggregatedPower()
	}
	return power.Copy() // return a copy as the aggregator is going to modify this value
}

func (a *aggregator) processRound(h uint64, r int64) {
	if a.empty(h, r) {
		return
	}

	roundInfo := a.messages[h][r]

	// TODO(lorenzo) can processing proposal before votes cause performance issues? (in case of a lot of proposals)
	for _, proposalEvent := range roundInfo.proposals {
		if a.toSkip(proposalEvent.Message) {
			continue
		}
		a.processProposal(proposalEvent, currentHeightEventBuilder)
	}

	nBatches := len(roundInfo.prevotes) + len(roundInfo.precommits)
	batches := make([][]events.UnverifiedMessageEvent, nBatches)
	i := 0

	// batch prevotes
	for _, events := range roundInfo.prevotes {
		batches[i] = events
		i++
	}

	// batch precommits
	for _, events := range roundInfo.precommits {
		batches[i] = events
		i++
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

		// clean up
		roundInfo.prevotes = make(map[common.Hash][]events.UnverifiedMessageEvent)
		roundInfo.prevotesPower = message.NewAggregatedPower()
		roundInfo.prevotesPowerFor = make(map[common.Hash]*message.AggregatedPower)

		// recompute total power for the round (precommits power + proposals)
		roundInfo.power = roundInfo.precommitsPower.Copy()
		for _, proposalEvent := range roundInfo.proposals {
			proposal := proposalEvent.Message.(*message.Propose)
			roundInfo.power.Set(proposal.SignerIndex(), proposal.Power())
		}
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
		roundInfo.precommits = make(map[common.Hash][]events.UnverifiedMessageEvent)
		roundInfo.precommitsPower = message.NewAggregatedPower()
		roundInfo.precommitsPowerFor = make(map[common.Hash]*message.AggregatedPower)

		// recompute total power for the round (prevotes power + proposals)
		roundInfo.power = roundInfo.prevotesPower.Copy()
		for _, proposalEvent := range roundInfo.proposals {
			proposal := proposalEvent.Message.(*message.Propose)
			roundInfo.power.Set(proposal.SignerIndex(), proposal.Power())
		}
	default:
		a.logger.Crit("Unexpected code", "c", c)
	}

}

func (a *aggregator) processVotesFor(h uint64, r int64, c uint8, v common.Hash) {
	if a.empty(h, r) {
		return
	}

	roundInfo := a.messages[h][r]

	// fetch batch
	switch c {
	case message.PrevoteCode:
		batch, ok := roundInfo.prevotes[v]
		if !ok {
			return
		}

		a.processBatches([][]events.UnverifiedMessageEvent{batch}, currentHeightEventBuilder)

		// clean up
		delete(roundInfo.prevotes, v)
		delete(roundInfo.prevotesPowerFor, v)

		//TODO(lorenzo) this might be too computation heavy

		// re-compute round power and prevote power
		roundInfo.prevotesPower = message.NewAggregatedPower()
		roundInfo.power = message.NewAggregatedPower()

		// prevotes
		for _, prevotesEvent := range roundInfo.prevotes {
			for _, e := range prevotesEvent {
				vote := e.Message.(message.Vote)
				for index, power := range vote.Signers().Powers() {
					roundInfo.power.Set(index, power)
					roundInfo.prevotesPower.Set(index, power)
				}
			}
		}

		// precommits
		for _, precommitsEvent := range roundInfo.precommits {
			for _, e := range precommitsEvent {
				vote := e.Message.(message.Vote)
				for index, power := range vote.Signers().Powers() {
					roundInfo.power.Set(index, power)
				}
			}
		}

		// proposals
		for _, proposalEvent := range roundInfo.proposals {
			proposal := proposalEvent.Message.(*message.Propose)
			roundInfo.power.Set(proposal.SignerIndex(), proposal.Power())
		}
	case message.PrecommitCode:
		batch, ok := roundInfo.precommits[v]
		if !ok {
			return
		}

		a.processBatches([][]events.UnverifiedMessageEvent{batch}, currentHeightEventBuilder)

		// clean up
		delete(roundInfo.precommits, v)
		delete(roundInfo.precommitsPowerFor, v)

		//TODO(lorenzo) this might be too computation heavy

		// re-compute round power and prevote power
		roundInfo.precommitsPower = message.NewAggregatedPower()
		roundInfo.power = message.NewAggregatedPower()

		// prevotes
		for _, prevotesEvent := range roundInfo.prevotes {
			for _, e := range prevotesEvent {
				vote := e.Message.(message.Vote)
				for index, power := range vote.Signers().Powers() {
					roundInfo.power.Set(index, power)
				}
			}
		}

		// precommits
		for _, precommitsEvent := range roundInfo.precommits {
			for _, e := range precommitsEvent {
				vote := e.Message.(message.Vote)
				for index, power := range vote.Signers().Powers() {
					roundInfo.power.Set(index, power)
					roundInfo.precommitsPower.Set(index, power)
				}
			}
		}

		// proposals
		for _, proposalEvent := range roundInfo.proposals {
			proposal := proposalEvent.Message.(*message.Propose)
			roundInfo.power.Set(proposal.SignerIndex(), proposal.Power())
		}
	default:
		a.logger.Crit("Unexpected code", "c", c)
	}

}

// a batch is a set of messages for same (height,round,code,value) ---> can be aggregated using FastAggregateVerify
func (a *aggregator) processBatches(batches [][]events.UnverifiedMessageEvent, eventer eventBuilder) {
	if len(batches) == 0 {
		return
	}

	processed := 0 // messages that go in the aggregator
	sent := 0      // messages that out of the aggregator (to Core and FD as valid msgs)
	for _, batch := range batches {
		if len(batch) == 0 {
			continue
		}
		if metrics.Enabled {
			BatchesBg.Add(int64(len(batch)))
		}
		processed += len(batch)

		var publicKeys []blst.PublicKey
		var signatures []blst.Signature
		var messages []message.Vote
		var senders []common.Address
		var errChs []chan<- error

		for _, e := range batch {
			m := e.Message
			// skip messages to be ignored or that are already in core
			if a.toSkip(m) {
				continue
			}

			messages = append(messages, m.(message.Vote))
			publicKeys = append(publicKeys, m.SignerKey())
			signatures = append(signatures, m.Signature())
			senders = append(senders, e.Sender)
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
		if metrics.Enabled {
			InvalidBg.Add(int64(len(invalids)))
		}
		for _, index := range invalids {
			a.logger.Info("Received invalid bls signature from", "peer", senders[index])
			a.handleInvalidMessage(errChs[index], message.ErrBadSignature, senders[index])
		}
	}
	a.logger.Debug("Aggregator processed messages", "processed", processed, "sent", sent)
}

func (a *aggregator) processProposal(proposalEvent events.UnverifiedMessageEvent, eventer eventBuilder) {
	proposal := proposalEvent.Message
	if err := proposal.Validate(); err != nil {
		a.handleInvalidMessage(proposalEvent.ErrCh, err, proposalEvent.Sender)
		return
	}
	go a.backend.Post(eventer(proposal, proposalEvent.ErrCh))
}

// assumes current or old round vote
// if add == true, the msg is saved in the aggregator.
// if add == false, the msg is not saved and only the power checks are done.
func (a *aggregator) handleVote(voteEvent events.UnverifiedMessageEvent, committee types.Committee, quorum *big.Int, add bool) {
	vote := voteEvent.Message.(message.Vote)
	errCh := voteEvent.ErrCh
	sender := voteEvent.Sender

	height := vote.H()
	round := vote.R()
	code := vote.Code()
	value := vote.Value()

	// complex aggregates always carry quorum (enforced at PreValidate)
	// if we do not already have quorum in Core, process right away
	coreVotesForPower := a.core.VotesPowerFor(height, round, code, value)
	coreVotesPower := a.core.VotesPower(height, round, code)
	if vote.Signers().IsComplex() && (coreVotesForPower.Power().Cmp(quorum) < 0 || coreVotesPower.Power().Cmp(quorum) < 0) {
		if err := vote.Validate(); err != nil {
			a.handleInvalidMessage(errCh, err, sender)
			return
		}
		go a.backend.Post(currentHeightEventBuilder(voteEvent.Message, errCh))
		return
	}

	// we are processing an individual vote or a simple aggregate vote
	if add {
		a.saveMessage(voteEvent)
	}

	// check if we reached quorum voting power on a specific value
	corePower := a.core.VotesPowerFor(height, round, code, value)
	aggregatorPower := a.votesPowerFor(height, round, code, value)
	contribution := powerContribution(aggregatorPower.Signers(), corePower.Signers(), committee)
	if corePower.Power().Cmp(quorum) < 0 && contribution.Add(contribution, corePower.Power()).Cmp(quorum) >= 0 {
		a.processVotesFor(height, round, code, value)
		return
	}

	// check if we reached quorum voting power in general
	corePower = a.core.VotesPower(height, round, code)
	aggregatorPower = a.votesPower(height, round, code)
	contribution = powerContribution(aggregatorPower.Signers(), corePower.Signers(), committee)
	if corePower.Power().Cmp(quorum) < 0 && contribution.Add(contribution, corePower.Power()).Cmp(quorum) >= 0 {
		a.processVotes(height, round, code)
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

	// channel where the aggregator will receive msgs from the backend handlers
	messageCh := a.backend.MessageCh()

loop:
	for {
		select {
		case event, ok := <-messageCh:
			start := time.Now()
			if !ok {
				break loop
			}
			if metrics.Enabled {
				BackendAggregatorTransitBg.Add(time.Since(event.Posted).Nanoseconds())
			}
			msg := event.Message
			sender := event.Sender

			a.messagesFrom[sender] = append(a.messagesFrom[sender], msg.Hash())

			// NOTE: Aggregator and Core run asynchronously. The code needs to take into account that Core can change state at any point here.
			// This also implies that height checks still needs to be done in Core.
			coreHeight := a.core.Height().Uint64()
			if msg.H() < coreHeight {
				a.logger.Debug("Storing old height message in the aggregator", "msgHeight", msg.H(), "coreHeight", coreHeight)
				signatureInput := msg.SignatureInput()
				a.staleMessages[signatureInput] = append(a.staleMessages[signatureInput], event)
				break
			}
			if msg.H() > coreHeight {
				// future messages are dealt with at backend peer handler level
				a.logger.Crit("future message in aggregator", "msgHeight", msg.H(), "coreHeight", coreHeight)
			}

			// if message already in Core, drop it
			if a.alreadyProcessed(msg) {
				a.logger.Debug("Discarding msg, already processed in Core")
				break
			}

			header := a.backend.BlockChain().GetHeaderByNumber(msg.H() - 1)
			if header == nil { //nolint
				a.logger.Crit("cannot fetch header for non-future message", "headerHeight", msg.H()-1, "coreHeight", coreHeight)
			}
			committee := header.Committee //nolint
			quorum := bft.Quorum(header.TotalVotingPower())

			coreRound := a.core.Round()
			if msg.R() > coreRound {
				// NOTE: here we could be buffering a proposal for future round, or a complex vote aggregate.
				a.saveMessage(event)
				// check if power is enough for a round skip
				aggregatorPower := a.power(msg.H(), msg.R())
				corePower := a.core.Power(msg.H(), msg.R())
				contribution := powerContribution(aggregatorPower.Signers(), corePower.Signers(), header.Committee)
				if contribution.Add(contribution, corePower.Power()).Cmp(bft.F(header.TotalVotingPower())) > 0 {
					a.logger.Debug("Processing future round messages due to possible round skip", "height", msg.H(), "round", msg.R(), "coreRound", coreRound)
					a.processRound(msg.H(), msg.R())
				}
				recordMessageProcessingTime(msg.Code(), start)
				break
			}

			// current or old round here
			switch msg.(type) {
			// if proposal, verify right away
			case *message.Propose:
				a.processProposal(event, currentHeightEventBuilder)
			case *message.Prevote, *message.Precommit:
				a.handleVote(event, committee, quorum, true)
			default:
				a.logger.Crit("unknown message type arrived in aggregator")
			}
			recordMessageProcessingTime(msg.Code(), start)
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

				if a.empty(height, round) {
					break
				}

				roundInfo := a.messages[height][round]

				// process proposals
				for _, proposalEvent := range roundInfo.proposals {
					if a.toSkip(proposalEvent.Message) {
						continue
					}
					a.processProposal(proposalEvent, currentHeightEventBuilder)
				}
				//clean up
				roundInfo.proposals = make([]events.UnverifiedMessageEvent, 0)

				header := a.backend.BlockChain().GetHeaderByNumber(height - 1)
				if header == nil { //nolint
					a.logger.Crit("cannot fetch header for non-future height", "height", height-1)
				}
				quorum := bft.Quorum(header.TotalVotingPower())
				committee := header.Committee //nolint

				for _, evs := range roundInfo.precommits {
					for _, e := range evs {
						//TODO(lorenzo) possible problem here, handleVote might call `process*` which is going to delete stuff in the map
						// generally is fine but double check if it causes issue in this case
						// also, can I just process one event for value to re-do the checks?
						a.handleVote(e, committee, quorum, false)
					}
				}

				for _, evs := range roundInfo.prevotes {
					for _, e := range evs {
						//TODO(lorenzo) possible problem here, processVote might call `process` which is going to delete stuff in the map
						// generally is fine but double check if it causes issue in this case
						// also, can I just process one event for value to re-do the checks?
						a.handleVote(e, committee, quorum, false)
					}
				}
				if metrics.Enabled {
					RoundBg.Add(time.Since(start).Nanoseconds())
				}
			case events.PowerChangeEvent:
				// a power change happened in Core: re-do quorum checks on individual votes and simple aggregates
				height := e.Height
				round := e.Round
				code := e.Code
				value := e.Value
				if a.empty(height, round) {
					break
				}

				roundInfo := a.messages[height][round]
				header := a.backend.BlockChain().GetHeaderByNumber(height - 1)
				if header == nil {
					panic("cannot fetch header for non-future message")
				}
				quorum := bft.Quorum(header.TotalVotingPower())
				committee := header.Committee

				var votesEvent []events.UnverifiedMessageEvent
				var ok bool
				switch code {
				case message.PrevoteCode:
					votesEvent, ok = roundInfo.prevotes[value]
				case message.PrecommitCode:
					votesEvent, ok = roundInfo.precommits[value]
				default:
					a.logger.Crit("Unexpected code", "code", code)
				}

				if !ok || len(votesEvent) == 0 {
					break
				}

				// processing one vote for the value for which power changed is enough to do all necessary checks
				a.handleVote(votesEvent[0], committee, quorum, false)
				if metrics.Enabled {
					PowerBg.Add(time.Since(start).Nanoseconds())
				}
			case events.FuturePowerChangeEvent:
				height := e.Height
				round := e.Round

				header := a.backend.BlockChain().GetHeaderByNumber(height - 1)
				if header == nil { //nolint
					a.logger.Crit("cannot fetch header for non-future message", "height", height-1)
				}
				committee := header.Committee //nolint

				// check in future round messages power, check again for round skip
				aggregatorPower := a.power(height, round)
				corePower := a.core.Power(height, round)
				contribution := powerContribution(aggregatorPower.Signers(), corePower.Signers(), committee)
				if contribution.Add(contribution, corePower.Power()).Cmp(bft.F(header.TotalVotingPower())) > 0 {
					a.processRound(height, round)
				}
				if metrics.Enabled {
					FuturePowerBg.Add(time.Since(start).Nanoseconds())
				}
			}
		case <-ticker.C:
			coreHeight := a.core.Height().Uint64()

			// process all messages in the aggregator
			for h, roundMap := range a.messages {
				for r, roundInfo := range roundMap {
					// if old height messages, move them to the staleMessages data structure. They are not useful for Core anymore.
					if h < coreHeight {
						//proposals
						for _, proposal := range roundInfo.proposals {
							signatureInput := proposal.Message.SignatureInput()
							a.staleMessages[signatureInput] = append(a.staleMessages[signatureInput], proposal)
						}
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
			a.messagesFrom = make(map[common.Address][]common.Hash)
			a.toIgnore = make(map[common.Hash]struct{})
		case <-oldMessagesTicker.C:
			var batches [][]events.UnverifiedMessageEvent
			for _, batch := range a.staleMessages {
				// if batch of proposals, validate them individually
				if batch[0].Message.Code() == message.ProposalCode {
					for _, proposalEvent := range batch {
						if a.toSkip(proposalEvent.Message) {
							continue
						}
						a.processProposal(proposalEvent, oldHeightEventBuilder)
					}
					continue
				}
				batches = append(batches, batch)
			}
			a.processBatches(batches, oldHeightEventBuilder)

			a.staleMessages = make(map[common.Hash][]events.UnverifiedMessageEvent)
		case <-ctx.Done():
			break loop
		}
	}
}

func (a *aggregator) stop() {
	a.logger.Info("Stopping the aggregator routine")
	a.cancel()
	a.coreSub.Unsubscribe()
	a.wg.Wait()
}
