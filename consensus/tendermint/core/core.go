package core

import (
	"context"
	"github.com/autonity/autonity/ethdb"
	"math/big"
	"sync"
	"time"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
)

// New creates a Tendermint consensus Core
func New(backend interfaces.Backend, services *interfaces.Services, address common.Address, logger log.Logger,
	noGossip bool, db ethdb.Database) *Core {
	c := &Core{
		blockPeriod:            1, // todo: retrieve it from contract
		address:                address,
		logger:                 logger,
		backend:                backend,
		futureRound:            make(map[int64][]message.Msg),
		futurePower:            make(map[int64]*message.AggregatedPower),
		pendingCandidateBlocks: make(map[uint64]*types.Block),
		stopped:                make(chan struct{}, 4),
		committee:              nil,
		proposeTimeout:         NewTimeout(Propose, logger),
		prevoteTimeout:         NewTimeout(Prevote, logger),
		precommitTimeout:       NewTimeout(Precommit, logger),
		newHeight:              time.Now(),
		newRound:               time.Now(),
		stepChange:             time.Now(),
		roundsState:            newTendermintState(logger, db, backend.BlockChain()),
		noGossip:               noGossip,
	}
	c.SetDefaultHandlers()
	if services != nil {
		c.broadcaster = services.Broadcaster(c)
		c.prevoter = services.Prevoter(c)
		c.precommiter = services.Precommiter(c)
		c.proposer = services.Proposer(c)
	}
	return c
}

func (c *Core) SetDefaultHandlers() {
	c.broadcaster = &Broadcaster{c}
	c.prevoter = &Prevoter{c}
	c.precommiter = &Precommiter{c}
	c.proposer = &Proposer{c}
}

type Core struct {
	blockPeriod uint64
	address     common.Address
	logger      log.Logger

	backend interfaces.Backend
	cancel  context.CancelFunc

	messageSub          *event.TypeMuxSubscription
	candidateBlockCh    chan events.NewCandidateBlockEvent
	committedCh         chan events.CommitEvent
	timeoutEventSub     *event.TypeMuxSubscription
	syncEventSub        *event.TypeMuxSubscription
	futureProposalTimer *time.Timer
	stopped             chan struct{}

	pendingCandidateBlocks map[uint64]*types.Block

	// used to ensure that the aggregator can get the correct power values by calling Power, VotesPower, VotesPowerFor
	roundChangeMu sync.Mutex

	// Tendermint SM state, and its underlying WAL storage.
	roundsState RoundsState
	committee   interfaces.Committee
	lastHeader  *types.Header
	stateMu     sync.RWMutex

	// height, round, committeeSet and lastHeader are the ONLY guarded fields.
	// everything else MUST be accessed only by the main thread.

	stepChange time.Time

	// future round messages are accessed also by the backend (to sync other peers) and the aggregator.
	// they need a lock.
	futureRound     map[int64][]message.Msg
	futurePower     map[int64]*message.AggregatedPower // power cache for future value msgs (per round)
	futureRoundLock sync.RWMutex

	proposeTimeout   *Timeout
	prevoteTimeout   *Timeout
	precommitTimeout *Timeout

	// End of Tendermint FSM fields

	protocolContracts *autonity.ProtocolContracts

	// tendermint behaviour interfaces, can be used in customizing the behaviours
	// during malicious testing
	broadcaster interfaces.Broadcaster
	prevoter    interfaces.Prevoter
	precommiter interfaces.Precommiter
	proposer    interfaces.Proposer

	// these timestamps are used to compute metrics for tendermint
	newHeight          time.Time
	newRound           time.Time
	currBlockTimeStamp time.Time
	noGossip           bool
}

func (c *Core) Prevoter() interfaces.Prevoter {
	return c.prevoter
}

func (c *Core) Precommiter() interfaces.Precommiter {
	return c.precommiter
}

func (c *Core) Proposer() interfaces.Proposer {
	return c.proposer
}

func (c *Core) Address() common.Address {
	return c.address
}

func (c *Core) Post(ev any) {
	switch ev := ev.(type) {
	case events.CommitEvent:
		c.committedCh <- ev
	case events.NewCandidateBlockEvent:
		c.candidateBlockCh <- ev
	}
}

func (c *Core) ProposeTimeout() *Timeout {
	return c.proposeTimeout
}

func (c *Core) PrevoteTimeout() *Timeout {
	return c.prevoteTimeout
}

func (c *Core) PrecommitTimeout() *Timeout {
	return c.precommitTimeout
}

func (c *Core) Broadcaster() interfaces.Broadcaster {
	return c.broadcaster
}

func (c *Core) Commit(ctx context.Context, round int64, messages *message.RoundMessages) {
	c.SetStep(ctx, PrecommitDone)
	// for metrics
	start := time.Now()
	proposal := messages.Proposal()
	if proposal == nil {
		// Should never happen really. Let's panic to catch bugs.
		panic("Core commit called with empty proposal")
		return
	}
	proposalHash := proposal.Block().Header().Hash()
	c.logger.Debug("Committing a block", "hash", proposalHash)

	precommitWithQuorum := messages.PrecommitFor(proposalHash)
	quorumCertificate := types.NewAggregateSignature(precommitWithQuorum.Signature().(*blst.BlsSignature), precommitWithQuorum.Signers())
	// record decision in WAL.
	c.SetDecision(proposal.Block())
	if err := c.backend.Commit(proposal.Block(), round, quorumCertificate); err != nil {
		c.logger.Error("failed to commit a block", "err", err)
		return
	}
	if metrics.Enabled {
		now := time.Now()
		CommitTimer.Update(now.Sub(start))
		CommitBg.Add(now.Sub(start).Nanoseconds())
	}
}

// Metric collecton of round change and height change.
func (c *Core) measureHeightRoundMetrics(round int64) {
	if round == 0 {
		// in case of height change, round changed too, so count it also.
		RoundChangeMeter.Mark(1)
		HeightChangeMeter.Mark(1)
	} else {
		RoundChangeMeter.Mark(1)
	}
}

type backlogMessageEvent struct {
	msg message.Msg
}

// current round == 0 --> height change
func (c *Core) processFuture(previousRound int64, currentRound int64) {
	if currentRound == 0 {
		// if height change, process future height messages
		go c.backend.ProcessFutureMsgs(c.Height().Uint64())
		return
	}

	// round change, process buffered future round messages
	c.futureRoundLock.Lock()
	defer c.futureRoundLock.Unlock()

	for r := previousRound + 1; r <= currentRound; r++ {
		for _, msg := range c.futureRound[r] {
			go c.SendEvent(backlogMessageEvent{
				msg: msg,
			})
		}
		delete(c.futureRound, r)
		delete(c.futurePower, r)
	}
}

// StartRound starts a new round. if round equals to 0, it means to starts a new height
func (c *Core) StartRound(ctx context.Context, round int64) {
	// todo: Jason, shall we keep this? As we read round from WAL, the round could exceed this constants.MaxRound.
	if round > constants.MaxRound {
		c.logger.Crit("⚠️ CONSENSUS FAILED ⚠️")
	}

	previousRound := c.Round()

	c.measureHeightRoundMetrics(round)
	// Set initial FSM state
	c.setInitialState(round)
	c.logger.Debug("Starting new Round", "Height", c.Height(), "Round", round)

	// If the node is the proposer for this round then it would propose validValue or a new block, otherwise,
	// proposeTimeout is started, where the node waits for a proposal from the proposer of the current round.
	if c.IsProposer() {
		// validValue and validRound represent a block they received a quorum of prevote and the round quorum was
		// received, respectively. If the block is not committed in that round then the round is changed.
		// The new proposer will chose the validValue, if present, which was set in one of the previous rounds otherwise
		// they propose a new block.
		if c.ValidValue() != nil {
			c.proposer.SendProposal(ctx, c.ValidValue())
		} else {
			// send proposal when there is available candidate rather than blocking the Core event loop, the
			// handleNewCandidateBlockMsg in the Core event loop will send proposal when the available one comes if we
			// don't have it sent here.
			newValue, ok := c.pendingCandidateBlocks[c.Height().Uint64()]
			if ok {
				c.proposer.SendProposal(ctx, newValue)
			}
		}
	} else {
		timeoutDuration := c.timeoutPropose(round)
		c.proposeTimeout.ScheduleTimeout(timeoutDuration, round, c.Height(), c.onTimeoutPropose)
		c.logger.Debug("Scheduled Propose Timeout", "Timeout Duration", timeoutDuration)
	}
	c.processFuture(previousRound, round)
	c.backend.Post(events.RoundChangeEvent{Height: c.Height().Uint64(), Round: round})
}

func (c *Core) setInitialState(r int64) {
	start := time.Now()
	c.roundChangeMu.Lock()
	RoundChangeMuBg.Add(time.Since(start).Nanoseconds())
	defer c.roundChangeMu.Unlock()

	// Start of new height where round is 0, it will reset all the states in WAL
	if r == 0 {
		lastBlockMined := c.backend.HeadBlock()
		c.SetHeight(new(big.Int).Add(lastBlockMined.Number(), common.Big1))
		lastHeader := lastBlockMined.Header()
		c.committee.SetLastHeader(lastHeader)
		c.setLastHeader(lastHeader)

		c.futureRoundLock.Lock()
		c.futureRound = make(map[int64][]message.Msg)
		c.futurePower = make(map[int64]*message.AggregatedPower)
		c.futureRoundLock.Unlock()
		// update height duration timer
		if metrics.Enabled {
			now := time.Now()
			HeightTimer.Update(now.Sub(c.newHeight))
			HeightBg.Add(now.Sub(c.newHeight).Nanoseconds())
			c.newHeight = now
		}
	}

	c.proposeTimeout.Reset(Propose)
	c.prevoteTimeout.Reset(Prevote)
	c.precommitTimeout.Reset(Precommit)

	c.SetRound(r)

	// update round duration timer
	if metrics.Enabled {
		now := time.Now()
		RoundTimer.Update(now.Sub(c.newRound))
		RoundBg.Add(now.Sub(c.newRound).Nanoseconds())
		c.newRound = now
	}
}

func (c *Core) SetStep(ctx context.Context, step Step) {
	now := time.Now()
	if metrics.Enabled {
		switch {
		// "standard" tendermint transitions
		case c.Step() == PrecommitDone && step == Propose: // precommitdone --> propose
			PrecommitDoneStepTimer.Update(now.Sub(c.stepChange))
			PrecommitDoneStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		case c.Step() == Propose && step == Prevote: // propose --> prevote
			ProposeStepTimer.Update(now.Sub(c.stepChange))
			ProposeStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		case c.Step() == Prevote && step == Precommit: // prevote --> precommit
			PrevoteStepTimer.Update(now.Sub(c.stepChange))
			PrevoteStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		case c.Step() == Precommit && step == PrecommitDone: // precommit --> precommitDone
			PrecommitStepTimer.Update(now.Sub(c.stepChange))
			PrecommitStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		// skipped to a future round
		case c.Step() == Propose && step == Propose:
			ProposeStepTimer.Update(now.Sub(c.stepChange))
			ProposeStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		case c.Step() == Prevote && step == Propose:
			PrevoteStepTimer.Update(now.Sub(c.stepChange))
			PrevoteStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		case c.Step() == Precommit && step == Propose:
			PrecommitStepTimer.Update(now.Sub(c.stepChange))
			PrecommitStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		// committing a proposal (old or current) due to receival of quorum precommits
		case c.Step() == Propose && step == PrecommitDone:
			ProposeStepTimer.Update(now.Sub(c.stepChange))
			ProposeStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		case c.Step() == Prevote && step == PrecommitDone:
			PrevoteStepTimer.Update(now.Sub(c.stepChange))
			PrevoteStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		default:
			// Ideally should be a .Crit, however it does not seem right to me because in the same sceneario the node would:
			// - crash if running the metrics
			// - keep validating without issues if not
			c.logger.Warn("Unexpected tendermint state transition", "c.step", c.Step(), "step", step)
		}
	}
	c.logger.Debug("Step change", "from", c.Step().String(), "to", step.String(), "round", c.Round())
	c.UpdateStep(step)
	c.stepChange = now

	// stop consensus timeouts
	c.stopAllTimeouts()

	// if we are moving from propose to prevote step we need to check again line 34,36 and 44
	// NOTE: this call to stepChangeChecks can cause recursion in the SetStep function.
	// This can happen if the checks cause a transition to Precommit step. It is expected behaviour.
	// If we want to remove this recursion possibility, we could post an Event that signals a step change,
	// which will then be processed in the MainEventLoop
	if c.Step() == Prevote {
		c.stepChangeChecks(ctx)
	}

}

// tries to stop all consensus timeouts
func (c *Core) stopAllTimeouts() {
	if err := c.proposeTimeout.StopTimer(); err != nil {
		c.logger.Debug("Cannot stop propose timer", "c.step", c.Step(), "err", err)
	}
	if err := c.prevoteTimeout.StopTimer(); err != nil {
		c.logger.Debug("Cannot stop prevote timer", "c.step", c.Step(), "err", err)
	}
	if err := c.precommitTimeout.StopTimer(); err != nil {
		c.logger.Debug("Cannot stop precommit timer", "c.step", c.Step(), "err", err)
	}
}

func (c *Core) setCommitteeSet(set interfaces.Committee) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.committee = set
}

func (c *Core) setLastHeader(lastHeader *types.Header) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.lastHeader = lastHeader
}

// state readers:
func (c *Core) Height() *big.Int {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.roundsState.Height()
}

func (c *Core) Round() int64 {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.roundsState.Round()
}

func (c *Core) Step() Step {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.roundsState.Step()
}

func (c *Core) Decision() *types.Block {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.roundsState.Decision()
}

func (c *Core) LockedRound() int64 {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.roundsState.LockedRound()
}

func (c *Core) ValidRound() int64 {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.roundsState.ValidRound()
}

func (c *Core) LockedValue() *types.Block {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.roundsState.LockedValue()
}

func (c *Core) ValidValue() *types.Block {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.roundsState.ValidValue()
}

func (c *Core) SentProposal() bool {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.roundsState.SentProposal()
}

func (c *Core) SentPrevote() bool {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.roundsState.SentPrevote()
}

func (c *Core) SentPrecommit() bool {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.roundsState.SentPrecommit()
}

func (c *Core) ValidRoundAndValueSet() bool {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.roundsState.ValidRoundAndValueSet()
}

func (c *Core) CurRoundMessages() *message.RoundMessages {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.roundsState.CurRoundMessages()
}

// state writers:
func (c *Core) SetHeight(height *big.Int) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.roundsState.StartNewHeight(height)
}

func (c *Core) SetRound(round int64) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.roundsState.StartNewRound(round)
}

func (c *Core) UpdateStep(step Step) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.roundsState.SetStep(step)
}

func (c *Core) SetDecision(block *types.Block) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.roundsState.SetDecision(block)
}

func (c *Core) SetSentProposal() {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.roundsState.SetSentProposal()
}

func (c *Core) SetSentPrevote() {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.roundsState.SetSentPrevote()
}

func (c *Core) SetSentPrecommit() {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.roundsState.SetSentPrecommit()
}

func (c *Core) SetLockedRoundAndValue(lockedRound int64, block *types.Block) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.roundsState.SetLockedRoundAndValue(lockedRound, block)
}

func (c *Core) SetValidRoundAndValue(validRound int64, block *types.Block) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.roundsState.SetValidRoundAndValue(validRound, block)
}

func (c *Core) CommitteeSet() interfaces.Committee {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.committee
}

func (c *Core) LastHeader() *types.Header {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.lastHeader
}

func (c *Core) Power(h uint64, r int64) *message.AggregatedPower {
	start := time.Now()
	c.roundChangeMu.Lock()
	RoundChangeMuBg.Add(time.Since(start).Nanoseconds())
	defer c.roundChangeMu.Unlock()

	if h != c.Height().Uint64() {
		return message.NewAggregatedPower()
	}

	power := message.NewAggregatedPower()
	if r > c.Round() {
		// future round
		c.futureRoundLock.RLock()
		futurePower, ok := c.futurePower[r]
		if ok {
			power = futurePower.Copy()
		}
		c.futureRoundLock.RUnlock()
	} else {
		// old or current round
		power = c.roundsState.GetOrCreate(r).Power()
	}

	return power
}

// NOTE: this assumes that r <= currentRound. If not, the returned power will be 0 even if there might be future round messages in c.futureRound
// This methods should not be used to compute power for future rounds
func (c *Core) VotesPower(h uint64, r int64, code uint8) *message.AggregatedPower {
	start := time.Now()
	c.roundChangeMu.Lock()
	RoundChangeMuBg.Add(time.Since(start).Nanoseconds())
	defer c.roundChangeMu.Unlock()

	if h != c.Height().Uint64() {
		return message.NewAggregatedPower()
	}
	roundMessages := c.roundsState.GetOrCreate(r)
	var power *message.AggregatedPower

	switch code {
	case message.ProposalCode:
		c.logger.Crit("Proposal code passed into VotesPower")
	case message.PrevoteCode:
		power = roundMessages.PrevotesTotalAggregatedPower()
	case message.PrecommitCode:
		power = roundMessages.PrecommitsTotalAggregatedPower()
	default:
		c.logger.Crit("unknown message code", "code", code)
	}
	return power
}

// NOTE: assume r <= currentRound. If not, the returned power will be 0 even if there might be future round messages in c.futureRound
// This methods should not be used to compute power for future rounds
func (c *Core) VotesPowerFor(h uint64, r int64, code uint8, v common.Hash) *message.AggregatedPower {
	start := time.Now()
	c.roundChangeMu.Lock()
	RoundChangeMuBg.Add(time.Since(start).Nanoseconds())
	defer c.roundChangeMu.Unlock()

	if h != c.Height().Uint64() {
		return message.NewAggregatedPower()
	}
	roundMessages := c.roundsState.GetOrCreate(r)
	var power *message.AggregatedPower

	switch code {
	case message.ProposalCode:
		c.logger.Crit("Proposal code passed into VotesPower")
	case message.PrevoteCode:
		power = roundMessages.PrevotesAggregatedPower(v)
	case message.PrecommitCode:
		power = roundMessages.PrecommitsAggregatedPower(v)
	default:
		c.logger.Crit("unknown message code", "code", code)
	}
	return power
}

// TODO: when we sync a peer, should we send him also the future round messages?
func (c *Core) CurrentHeightMessages() []message.Msg {
	return c.roundsState.AllMessages()
}

func (c *Core) Backend() interfaces.Backend {
	return c.backend
}
func (c *Core) Logger() log.Logger {
	return c.logger
}

func (c *Core) IsFromProposer(round int64, address common.Address) bool {
	return c.CommitteeSet().GetProposer(round).Address == address
}

func (c *Core) IsProposer() bool {
	return c.CommitteeSet().GetProposer(c.Round()).Address == c.address
}

func (c *Core) BroadcastAll(msg message.Msg) {
	c.Backend().Broadcast(c.CommitteeSet().Committee(), msg)
}

type Broadcaster struct {
	*Core
}

func (s *Broadcaster) Broadcast(msg message.Msg) {
	s.logger.Debug("Broadcasting", "message", log.Lazy{Fn: msg.String})
	s.BroadcastAll(msg)
}
