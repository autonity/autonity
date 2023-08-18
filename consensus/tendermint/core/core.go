package core

import (
	"context"
	"math/big"
	"reflect"
	"sync"
	"time"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	tctypes "github.com/autonity/autonity/consensus/tendermint/core/types"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
)

// New creates an Tendermint consensus Core
func New(backend interfaces.Backend) *Core {
	addr := backend.Address()
	messagesMap := message.NewMessagesMap()
	roundMessage := messagesMap.GetOrCreate(0)
	c := &Core{
		blockPeriod:            1, // todo: retrieve it from contract
		address:                addr,
		logger:                 backend.Logger(),
		backend:                backend,
		backlogs:               make(map[common.Address][]*message.Message),
		backlogUntrusted:       make(map[uint64][]*message.Message),
		pendingCandidateBlocks: make(map[uint64]*types.Block),
		stopped:                make(chan struct{}, 4),
		committee:              nil,
		futureRoundChange:      make(map[int64]map[common.Address]*big.Int),
		messages:               messagesMap,
		lockedRound:            -1,
		validRound:             -1,
		curRoundMessages:       roundMessage,
		proposeTimeout:         tctypes.NewTimeout(tctypes.Propose, backend.Logger()),
		prevoteTimeout:         tctypes.NewTimeout(tctypes.Prevote, backend.Logger()),
		precommitTimeout:       tctypes.NewTimeout(tctypes.Precommit, backend.Logger()),
		newHeight:              time.Now(),
		newRound:               time.Now(),
		stepChange:             time.Now(),
	}
	c.SetDefaultHandlers()
	return c
}

func (c *Core) SetDefaultHandlers() {
	c.broadcaster = &Broadcaster{c}
	c.prevoter = &Prevoter{c}
	c.precommiter = &Precommiter{c}
	c.proposer = &Proposer{c}
}

func (c *Core) SetBroadcaster(svc interfaces.Broadcaster) {
	if svc == nil {
		return
	}
	// this would set the current Core object state in the
	// broadcast service object
	field0 := reflect.ValueOf(svc).Elem().Field(0)
	field0.Set(reflect.ValueOf(c))
	c.broadcaster = svc
}
func (c *Core) SetPrevoter(svc interfaces.Prevoter) {
	if svc == nil {
		return
	}
	fields := reflect.ValueOf(svc).Elem()
	// Set up default Core
	field0 := fields.Field(0)
	field0.Set(reflect.ValueOf(c))
	// Set up default prevote service
	if fields.NumField() > 1 {
		field1 := fields.Field(1)
		field1.Set(reflect.ValueOf(c.prevoter))
	}
	c.prevoter = svc
}

func (c *Core) SetPrecommitter(svc interfaces.Precommiter) {
	if svc == nil {
		return
	}
	fields := reflect.ValueOf(svc).Elem()
	// Set up default Core
	field0 := fields.Field(0)
	field0.Set(reflect.ValueOf(c))
	// Set up default precommit service
	if fields.NumField() > 1 {
		field1 := fields.Field(1)
		field1.Set(reflect.ValueOf(c.precommiter))
	}

	c.precommiter = svc
}

func (c *Core) SetProposer(svc interfaces.Proposer) {
	if svc == nil {
		return
	}
	fields := reflect.ValueOf(svc).Elem()
	// Set up default Core
	field0 := fields.Field(0)
	field0.Set(reflect.ValueOf(c))
	// Set up default propose service
	if fields.NumField() > 1 {
		field1 := fields.Field(1)
		field1.Set(reflect.ValueOf(c.proposer))
	}
	c.proposer = svc
}

type Core struct {
	blockPeriod uint64
	address     common.Address
	logger      log.Logger

	backend interfaces.Backend
	cancel  context.CancelFunc

	messageEventSub        *event.TypeMuxSubscription
	candidateBlockEventSub *event.TypeMuxSubscription
	committedSub           *event.TypeMuxSubscription
	timeoutEventSub        *event.TypeMuxSubscription
	syncEventSub           *event.TypeMuxSubscription
	futureProposalTimer    *time.Timer
	stopped                chan struct{}

	backlogs             map[common.Address][]*message.Message
	backlogUntrusted     map[uint64][]*message.Message
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
	step                  tctypes.Step
	stepChange            time.Time
	curRoundMessages      *message.RoundMessages
	messages              *message.MessagesMap
	sentProposal          bool
	sentPrevote           bool
	sentPrecommit         bool
	setValidRoundAndValue bool

	lockedRound int64
	validRound  int64
	lockedValue *types.Block
	validValue  *types.Block

	proposeTimeout   *tctypes.Timeout
	prevoteTimeout   *tctypes.Timeout
	precommitTimeout *tctypes.Timeout

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

func (c *Core) GetPrevoter() interfaces.Prevoter {
	return c.prevoter
}

func (c *Core) GetPrecommiter() interfaces.Precommiter {
	return c.precommiter
}

func (c *Core) GetProposer() interfaces.Proposer {
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

func (c *Core) Step() tctypes.Step {
	return c.step
}

func (c *Core) CurRoundMessages() *message.RoundMessages {
	return c.curRoundMessages
}

func (c *Core) Messages() *message.MessagesMap {
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

func (c *Core) ProposeTimeout() *tctypes.Timeout {
	return c.proposeTimeout
}

func (c *Core) PrevoteTimeout() *tctypes.Timeout {
	return c.prevoteTimeout
}

func (c *Core) PrecommitTimeout() *tctypes.Timeout {
	return c.precommitTimeout
}

func (c *Core) FutureRoundChange() map[int64]map[common.Address]*big.Int {
	return c.futureRoundChange
}

func (c *Core) SetFutureRoundChange(futureRoundChange map[int64]map[common.Address]*big.Int) {
	c.futureRoundChange = futureRoundChange
}

func (c *Core) Br() interfaces.Broadcaster {
	return c.broadcaster
}

func (c *Core) SetBr(br interfaces.Broadcaster) {
	c.broadcaster = br
}

func (c *Core) CurrentHeightMessages() []*message.Message {
	return c.messages.Messages()
}

func (c *Core) SignMessage(msg *message.Message) ([]byte, error) {
	data, err := msg.BytesNoSignature()
	if err != nil {
		return nil, err
	}
	if msg.Signature, err = c.backend.Sign(data); err != nil {
		return nil, err
	}
	return msg.GetBytes(), nil
}

func (c *Core) Commit(round int64, messages *message.RoundMessages) {
	c.SetStep(tctypes.PrecommitDone)

	// for metrics
	start := time.Now()

	proposal := messages.Proposal()
	if proposal == nil {
		// Should never happen really.
		c.logger.Error("Core commit called with empty proposal ")
		return
	}

	if proposal.ProposalBlock == nil {
		// Again should never happen.
		c.logger.Error("commit a NIL block",
			"block", proposal.ProposalBlock,
			"height", c.Height(),
			"round", round)
		return
	}

	c.logger.Debug("commit a block", "hash", proposal.ProposalBlock.Header().Hash())

	committedSeals := make([][]byte, 0)
	for _, v := range messages.CommitedSeals(proposal.ProposalBlock.Hash()) {
		seal := make([]byte, types.BFTExtraSeal)
		copy(seal[:], v.CommittedSeal[:])
		committedSeals = append(committedSeals, seal)
	}

	if err := c.backend.Commit(proposal.ProposalBlock, round, committedSeals); err != nil {
		c.logger.Error("failed to commit a block", "err", err)
		return
	}

	if metrics.Enabled {
		now := time.Now()
		tctypes.CommitTimer.Update(now.Sub(start))
		tctypes.CommitBg.Add(now.Sub(start).Nanoseconds())
	}
}

// Metric collecton of round change and height change.
func (c *Core) MeasureHeightRoundMetrics(round int64) {
	if round == 0 {
		// in case of height change, round changed too, so count it also.
		tctypes.RoundChangeMeter.Mark(1)
		tctypes.HeightChangeMeter.Mark(1)
	} else {
		tctypes.RoundChangeMeter.Mark(1)
	}
}

// StartRound starts a new round. if round equals to 0, it means to starts a new height
func (c *Core) StartRound(ctx context.Context, round int64) {
	if round > constants.MaxRound {
		c.logger.Crit("⚠️ CONSENSUS FAILED ⚠️")
	}

	c.MeasureHeightRoundMetrics(round)
	// Set initial FSM state
	c.SetInitialState(round)
	// c.setStep(propose) will process the pending unmined blocks sent by the backed.Seal() and set c.lastestPendingRequest
	c.SetStep(tctypes.Propose)
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

func (c *Core) SetInitialState(r int64) {
	// Start of new height where round is 0
	if r == 0 {
		lastBlockMined, _ := c.backend.HeadBlock()
		c.setHeight(new(big.Int).Add(lastBlockMined.Number(), common.Big1))
		lastHeader := lastBlockMined.Header()
		c.committee.SetLastHeader(lastHeader)
		c.setLastHeader(lastHeader)
		c.lockedRound = -1
		c.lockedValue = nil
		c.validRound = -1
		c.validValue = nil
		c.messages.Reset()
		c.futureRoundChange = make(map[int64]map[common.Address]*big.Int)
		// update height duration timer
		if metrics.Enabled {
			now := time.Now()
			tctypes.HeightTimer.Update(now.Sub(c.newHeight))
			tctypes.HeightBg.Add(now.Sub(c.newHeight).Nanoseconds())
			c.newHeight = now
		}
	}

	c.proposeTimeout.Reset(tctypes.Propose)
	c.prevoteTimeout.Reset(tctypes.Prevote)
	c.precommitTimeout.Reset(tctypes.Precommit)
	c.curRoundMessages = c.messages.GetOrCreate(r)
	c.sentProposal = false
	c.sentPrevote = false
	c.sentPrecommit = false
	c.setValidRoundAndValue = false
	c.setRound(r)

	// update round duration timer
	if metrics.Enabled {
		now := time.Now()
		tctypes.RoundTimer.Update(now.Sub(c.newRound))
		tctypes.RoundBg.Add(now.Sub(c.newRound).Nanoseconds())
		c.newRound = now
	}
}

func (c *Core) AcceptVote(roundMsgs *message.RoundMessages, step tctypes.Step, hash common.Hash, msg message.Message) {
	switch step {
	case tctypes.Prevote:
		roundMsgs.AddPrevote(hash, msg)
	case tctypes.Precommit:
		roundMsgs.AddPrecommit(hash, msg)
	}
}

func (c *Core) SetStep(step tctypes.Step) {
	now := time.Now()
	if metrics.Enabled {
		switch {
		// "standard" tendermint transitions
		case c.step == tctypes.PrecommitDone && step == tctypes.Propose: // precommitdone --> propose
			tctypes.PrecommitDoneStepTimer.Update(now.Sub(c.stepChange))
			tctypes.PrecommitDoneStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		case c.step == tctypes.Propose && step == tctypes.Prevote: // propose --> prevote
			tctypes.ProposeStepTimer.Update(now.Sub(c.stepChange))
			tctypes.ProposeStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		case c.step == tctypes.Prevote && step == tctypes.Precommit: // prevote --> precommit
			tctypes.PrevoteStepTimer.Update(now.Sub(c.stepChange))
			tctypes.PrevoteStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		case c.step == tctypes.Precommit && step == tctypes.PrecommitDone: // precommit --> precommitDone
			tctypes.PrecommitStepTimer.Update(now.Sub(c.stepChange))
			tctypes.PrecommitStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		// skipped to a future round
		case c.step == tctypes.Propose && step == tctypes.Propose:
			tctypes.ProposeStepTimer.Update(now.Sub(c.stepChange))
			tctypes.ProposeStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		case c.step == tctypes.Prevote && step == tctypes.Propose:
			tctypes.PrevoteStepTimer.Update(now.Sub(c.stepChange))
			tctypes.PrevoteStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		case c.step == tctypes.Precommit && step == tctypes.Propose:
			tctypes.PrecommitStepTimer.Update(now.Sub(c.stepChange))
			tctypes.PrecommitStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		// committing an old round proposal
		case c.step == tctypes.Propose && step == tctypes.PrecommitDone:
			tctypes.ProposeStepTimer.Update(now.Sub(c.stepChange))
			tctypes.ProposeStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		case c.step == tctypes.Prevote && step == tctypes.PrecommitDone:
			tctypes.PrevoteStepTimer.Update(now.Sub(c.stepChange))
			tctypes.PrevoteStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
		case c.step == tctypes.PrecommitDone && step == tctypes.PrecommitDone:
			//this transition can also happen when we already received 2f+1 precommits but we did not start the new round yet.
			tctypes.PrecommitDoneStepTimer.Update(now.Sub(c.stepChange))
			tctypes.PrecommitDoneStepBg.Add(now.Sub(c.stepChange).Nanoseconds())
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

type Broadcaster struct {
	*Core
}

func (s *Broadcaster) SignAndBroadcast(ctx context.Context, msg *message.Message) {
	logger := s.Logger().New("step", s.Step())
	payload, err := s.SignMessage(msg)
	if err != nil {
		// This should not fail ..
		logger.Error("Failed to finalize message", "message", msg, "err", err)
		return
	}
	// SignAndBroadcast payload
	logger.Debug("Broadcasting", "message", msg.String())
	if err := s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), payload); err != nil {
		logger.Error("Failed to broadcast message", "msg", msg, "err", err)
		return
	}
}
