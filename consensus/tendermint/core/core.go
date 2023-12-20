package core

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
)

// New creates a Tendermint consensus Core
func New(backend interfaces.Backend, services *interfaces.Services) *Core {
	messagesMap := message.NewMap()
	roundMessage := messagesMap.GetOrCreate(0)
	c := &Core{
		blockPeriod:            1, // todo: retrieve it from contract
		address:                backend.Address(),
		logger:                 backend.Logger(),
		backend:                backend,
		backlogs:               make(map[common.Address][]message.Msg),
		backlogUntrusted:       make(map[uint64][]message.Msg),
		pendingCandidateBlocks: make(map[uint64]*types.Block),
		stopped:                make(chan struct{}, 4),
		committee:              nil,
		futureRoundChange:      make(map[int64]map[common.Address]*big.Int),
		messages:               messagesMap,
		lockedRound:            -1,
		validRound:             -1,
		curRoundMessages:       roundMessage,
		proposeTimeout:         NewTimeout(Propose, backend.Logger()),
		prevoteTimeout:         NewTimeout(Prevote, backend.Logger()),
		precommitTimeout:       NewTimeout(Precommit, backend.Logger()),
		newHeight:              time.Now(),
		newRound:               time.Now(),
		stepChange:             time.Now(),
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
	candidateBlockSub   *event.TypeMuxSubscription
	committedSub        *event.TypeMuxSubscription
	timeoutEventSub     *event.TypeMuxSubscription
	syncEventSub        *event.TypeMuxSubscription
	futureProposalTimer *time.Timer
	stopped             chan struct{}

	backlogs             map[common.Address][]message.Msg
	backlogUntrusted     map[uint64][]message.Msg
	backlogUntrustedSize int
	// map[Height]UnminedBlock
	pendingCandidateBlocks map[uint64]*types.Block

	//
	// Tendermint FSM state fields
	//

	stateMu    sync.RWMutex
	height     *big.Int
	round      int64
	committee  interfaces.Committee
	lastHeader *types.Header
	// height, round, committeeSet and lastHeader are the ONLY guarded fields.
	// everything else MUST be accessed only by the main thread.
	step                  Step
	stepChange            time.Time
	curRoundMessages      *message.RoundMessages
	messages              *message.Map
	sentProposal          bool
	sentPrevote           bool
	sentPrecommit         bool
	setValidRoundAndValue bool

	lockedRound int64
	validRound  int64
	lockedValue *types.Block
	validValue  *types.Block

	proposeTimeout   *Timeout
	prevoteTimeout   *Timeout
	precommitTimeout *Timeout

	futureRoundChange map[int64]map[common.Address]*big.Int

	protocolContracts *autonity.ProtocolContracts

	// tendermint behaviour interfaces, can be used in customizing the behaviours
	// during malicious testing
	broadcaster interfaces.Broadcaster
	prevoter    interfaces.Prevoter
	precommiter interfaces.Precommiter
	proposer    interfaces.Proposer

	// these timestamps are used to compute metrics for tendermint
	newHeight time.Time
	newRound  time.Time
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

func (c *Core) Committee() interfaces.Committee {
	return c.committee
}

func (c *Core) SetCommittee(committee interfaces.Committee) {
	c.committee = committee
}

func (c *Core) Step() Step {
	return c.step
}

func (c *Core) CurRoundMessages() *message.RoundMessages {
	return c.curRoundMessages
}

func (c *Core) Messages() *message.Map {
	return c.messages
}

func (c *Core) SentProposal() bool {
	return c.sentProposal
}

func (c *Core) SetSentProposal(sentProposal bool) {
	c.sentProposal = sentProposal
}

func (c *Core) SentPrevote() bool {
	return c.sentPrevote
}

func (c *Core) SetSentPrevote(sentPrevote bool) {
	c.sentPrevote = sentPrevote
}

func (c *Core) SentPrecommit() bool {
	return c.sentPrecommit
}

func (c *Core) SetSentPrecommit(sentPrecommit bool) {
	c.sentPrecommit = sentPrecommit
}

func (c *Core) SetValidRoundAndValue() bool {
	return c.setValidRoundAndValue
}

func (c *Core) SetSetValidRoundAndValue(setValidRoundAndValue bool) {
	c.setValidRoundAndValue = setValidRoundAndValue
}

func (c *Core) LockedRound() int64 {
	return c.lockedRound
}

func (c *Core) SetLockedRound(lockedRound int64) {
	c.lockedRound = lockedRound
}

func (c *Core) ValidRound() int64 {
	return c.validRound
}

func (c *Core) SetValidRound(validRound int64) {
	c.validRound = validRound
}

func (c *Core) LockedValue() *types.Block {
	return c.lockedValue
}

func (c *Core) SetLockedValue(lockedValue *types.Block) {
	c.lockedValue = lockedValue
}

func (c *Core) ValidValue() *types.Block {
	return c.validValue
}

func (c *Core) SetValidValue(validValue *types.Block) {
	c.validValue = validValue
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

func (c *Core) FutureRoundChange() map[int64]map[common.Address]*big.Int {
	return c.futureRoundChange
}

func (c *Core) SetFutureRoundChange(futureRoundChange map[int64]map[common.Address]*big.Int) {
	c.futureRoundChange = futureRoundChange
}

func (c *Core) Broadcaster() interfaces.Broadcaster {
	return c.broadcaster
}

func (c *Core) Commit(round int64, messages *message.RoundMessages) {
	c.SetStep(PrecommitDone)
	// for metrics
	start := time.Now()
	proposal := messages.Proposal()
	if proposal == nil {
		// Should never happen really.
		c.logger.Error("Core commit called with empty proposal")
		return
	}
	proposalHash := proposal.Block().Header().Hash()
	c.logger.Debug("Committing a block", "hash", proposalHash)

	committedSeals := make([][]byte, 0)
	for _, v := range messages.PrecommitsFor(proposalHash) {
		committedSeals = append(committedSeals, v.Signature())
	}

	if err := c.backend.Commit(proposal.Block(), round, committedSeals); err != nil {
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

// StartRound starts a new round. if round equals to 0, it means to starts a new height
func (c *Core) StartRound(ctx context.Context, round int64) {
	if round > constants.MaxRound {
		c.logger.Crit("⚠️ CONSENSUS FAILED ⚠️")
	}

	c.measureHeightRoundMetrics(round)
	// Set initial FSM state
	c.setInitialState(round)
	// c.setStep(propose) will process the pending unmined blocks sent by the backed.Seal() and set c.lastestPendingRequest
	c.SetStep(Propose)
	c.logger.Debug("Starting new Round", "Height", c.Height(), "Round", round)

	// If the node is the proposer for this round then it would propose validValue or a new block, otherwise,
	// proposeTimeout is started, where the node waits for a proposal from the proposer of the current round.
	if c.IsProposer() {
		// validValue and validRound represent a block they received a quorum of prevote and the round quorum was
		// received, respectively. If the block is not committed in that round then the round is changed.
		// The new proposer will chose the validValue, if present, which was set in one of the previous rounds otherwise
		// they propose a new block.
		if c.validValue != nil {
			c.proposer.SendProposal(ctx, c.validValue)
			return
		}
		// send proposal when there is available candidate rather than blocking the Core event loop, the
		// handleNewCandidateBlockMsg in the Core event loop will send proposal when the available one comes if we
		// don't have it sent here.
		newValue, ok := c.pendingCandidateBlocks[c.Height().Uint64()]
		if ok {
			c.proposer.SendProposal(ctx, newValue)
		}
	} else {
		timeoutDuration := c.timeoutPropose(round)
		c.proposeTimeout.ScheduleTimeout(timeoutDuration, round, c.Height(), c.onTimeoutPropose)
		c.logger.Debug("Scheduled Propose Timeout", "Timeout Duration", timeoutDuration)
	}
}

func (c *Core) setInitialState(r int64) {
	// Start of new height where round is 0
	if r == 0 {
		lastBlockMined := c.backend.HeadBlock()
		c.setHeight(new(big.Int).Add(lastBlockMined.Number(), common.Big1))
		lastHeader := lastBlockMined.Header()
		c.committee.SetLastHeader(lastHeader)
		c.setLastHeader(lastHeader)

		// on epoch rotation, update committee.
		if lastBlockMined.LastEpochBlock().Cmp(lastBlockMined.Number()) == 0 {
			log.Debug("on epoch rotation, update committee!", "number", lastBlockMined.Number())
			c.committee.SetCommittee(lastBlockMined.Header().Committee)
		}
		/*
			// TODO(lorenzo) deal better with stop start failure
			if c.EpochHead() == nil {
				epochHeader := c.backend.BlockChain().GetHeaderByNumber(lastBlockMined.LastEpochBlock().Uint64())
				c.setEpochHead(epochHeader)
			}*/

		c.lockedRound = -1
		c.lockedValue = nil
		c.validRound = -1
		c.validValue = nil
		c.messages.Reset()
		c.futureRoundChange = make(map[int64]map[common.Address]*big.Int)
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
	c.curRoundMessages = c.messages.GetOrCreate(r)
	c.sentProposal = false
	c.sentPrevote = false
	c.sentPrecommit = false
	c.setValidRoundAndValue = false
	c.setRound(r)

	// update round duration timer
	if metrics.Enabled {
		now := time.Now()
		RoundTimer.Update(now.Sub(c.newRound))
		RoundBg.Add(now.Sub(c.newRound).Nanoseconds())
		c.newRound = now
	}
}

/*
	func (c *Core) AcceptVote(roundMsgs *message.RoundMessages, step Step, hash common.Hash, msg message.Message) {
		switch step {
		case Prevote:
			roundMsgs.AddPrevote(hash, msg)
		case Precommit:
			roundMsgs.AddPrecommit(hash, msg)
		}
	}
*/
func (c *Core) SetStep(step Step) {
	now := time.Now()
	if metrics.Enabled {
		switch {
		// "standard" tendermint transitions
		case c.step == PrecommitDone && step == Propose: // precommitdone --> propose
			PrecommitDoneStepTimer.Update(now.Sub(c.stepChange))
			PrecommitDoneStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		case c.step == Propose && step == Prevote: // propose --> prevote
			ProposeStepTimer.Update(now.Sub(c.stepChange))
			ProposeStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		case c.step == Prevote && step == Precommit: // prevote --> precommit
			PrevoteStepTimer.Update(now.Sub(c.stepChange))
			PrevoteStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		case c.step == Precommit && step == PrecommitDone: // precommit --> precommitDone
			PrecommitStepTimer.Update(now.Sub(c.stepChange))
			PrecommitStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		// skipped to a future round
		case c.step == Propose && step == Propose:
			ProposeStepTimer.Update(now.Sub(c.stepChange))
			ProposeStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		case c.step == Prevote && step == Propose:
			PrevoteStepTimer.Update(now.Sub(c.stepChange))
			PrevoteStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		case c.step == Precommit && step == Propose:
			PrecommitStepTimer.Update(now.Sub(c.stepChange))
			PrecommitStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		// committing an old round proposal
		case c.step == Propose && step == PrecommitDone:
			ProposeStepTimer.Update(now.Sub(c.stepChange))
			ProposeStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		case c.step == Prevote && step == PrecommitDone:
			PrevoteStepTimer.Update(now.Sub(c.stepChange))
			PrevoteStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		case c.step == PrecommitDone && step == PrecommitDone:
			//this transition can also happen when we already received 2f+1 precommits but we did not start the new round yet.
			PrecommitDoneStepTimer.Update(now.Sub(c.stepChange))
			PrecommitDoneStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		default:
			// TODO(lorenzo) this ideally should be a .Crit but these transitions do actually happen.
			// see: https://github.com/autonity/autonity/issues/803
			c.logger.Warn("Unexpected tendermint state transition", "c.step", c.step, "step", step)
		}
	}
	c.logger.Debug("moving to step", "step", step.String(), "round", c.Round())
	c.step = step
	c.stepChange = now
	c.processBacklog()
}

func (c *Core) setRound(round int64) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.round = round
}

func (c *Core) setHeight(height *big.Int) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.height = height
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

func (c *Core) Round() int64 {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.round
}

func (c *Core) Height() *big.Int {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.height
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

func (c *Core) CurrentHeightMessages() []message.Msg {
	return c.messages.All()
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
	logger := s.Logger().New("step", s.Step())
	logger.Debug("Broadcasting", "message", msg.String())
	s.BroadcastAll(msg)
}
