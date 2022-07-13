package core

import (
	"context"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	tctypes "github.com/autonity/autonity/consensus/tendermint/core/types"
	"math/big"
	"reflect"
	"sync"
	"time"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
)

// New creates an Tendermint consensus Core
func New(backend interfaces.Backend) *Core {
	addr := backend.Address()
	logger := log.New("addr", addr.String())
	messagesMap := messageutils.NewMessagesMap()
	roundMessage := messagesMap.GetOrCreate(0)
	c := &Core{
		blockPeriod:            1, // todo: retrieve it from contract
		address:                addr,
		logger:                 logger,
		backend:                backend,
		backlogs:               make(map[common.Address][]*messageutils.Message),
		backlogUnchecked:       make(map[uint64][]*messageutils.Message),
		pendingCandidateBlocks: make(map[uint64]*types.Block),
		stopped:                make(chan struct{}, 4),
		committee:              nil,
		futureRoundChange:      make(map[int64]map[common.Address]uint64),
		messages:               messagesMap,
		lockedRound:            -1,
		validRound:             -1,
		curRoundMessages:       roundMessage,
		proposeTimeout:         tctypes.NewTimeout(tctypes.Propose, logger),
		prevoteTimeout:         tctypes.NewTimeout(tctypes.Prevote, logger),
		precommitTimeout:       tctypes.NewTimeout(tctypes.Precommit, logger),
	}
	c.SetDefaultHandlers()
	return c
}

func (c *Core) SetDefaultHandlers() {
	c.br = &BroadCastService{c}
	c.prevoter = &PrevoteService{c}
	c.precommiter = &PrecommitService{c}
	c.proposer = &ProposeService{c}
}

func (c *Core) SetBroadcastHandler(svc interfaces.Broadcaster) {
	if svc == nil {
		return
	}
	// this would set the current Core object state in the
	// broadcast service object
	field0 := reflect.ValueOf(svc).Elem().Field(0)
	field0.Set(reflect.ValueOf(c))
	c.br = svc
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

	backlogs            map[common.Address][]*messageutils.Message
	backlogUnchecked    map[uint64][]*messageutils.Message
	backlogUncheckedLen int
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
	curRoundMessages      *messageutils.RoundMessages
	messages              *messageutils.MessagesMap
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

	futureRoundChange map[int64]map[common.Address]uint64

	autonityContract *autonity.Contract

	// tendermint behaviour interfaces, can be used in customizing the behaviours
	// during malicious testing
	br          interfaces.Broadcaster
	prevoter    interfaces.Prevoter
	precommiter interfaces.Precommiter
	proposer    interfaces.Proposer
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

func (c *Core) CurRoundMessages() *messageutils.RoundMessages {
	return c.curRoundMessages
}

func (c *Core) Messages() *messageutils.MessagesMap {
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

func (c *Core) FutureRoundChange() map[int64]map[common.Address]uint64 {
	return c.futureRoundChange
}

func (c *Core) SetFutureRoundChange(futureRoundChange map[int64]map[common.Address]uint64) {
	c.futureRoundChange = futureRoundChange
}

func (c *Core) Br() interfaces.Broadcaster {
	return c.br
}

func (c *Core) SetBr(br interfaces.Broadcaster) {
	c.br = br
}

type BroadCastService struct {
	*Core
}

func (s *BroadCastService) Broadcast(ctx context.Context, msg *messageutils.Message) {
	logger := s.Logger().New("step", s.Step())

	payload, err := s.FinalizeMessage(msg)
	if err != nil {
		logger.Error("Failed to finalize message", "msg", msg, "err", err)
		return
	}

	// Broadcast payload
	logger.Debug("broadcasting", "msg", msg.String())
	if err = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), payload); err != nil {
		logger.Error("Failed to broadcast message", "msg", msg, "err", err)
		return
	}
}

func (c *Core) GetCurrentHeightMessages() []*messageutils.Message {
	return c.messages.GetMessages()
}

func (c *Core) IsMember(address common.Address) bool {
	_, _, err := c.CommitteeSet().GetByAddress(address)
	return err == nil
}

func (c *Core) FinalizeMessage(msg *messageutils.Message) ([]byte, error) {
	var err error

	// Sign message
	data, err := msg.PayloadNoSig()
	if err != nil {
		return nil, err
	}
	msg.Signature, err = c.backend.Sign(data)
	if err != nil {
		return nil, err
	}

	return msg.GetPayload(), nil
}

// check if msg sender is proposer for proposal handling.
func (c *Core) IsProposerMsg(round int64, msgAddress common.Address) bool {
	return c.CommitteeSet().GetProposer(round).Address == msgAddress
}
func (c *Core) IsProposer() bool {
	return c.CommitteeSet().GetProposer(c.Round()).Address == c.address
}

func (c *Core) Commit(round int64, messages *messageutils.RoundMessages) {
	c.SetStep(tctypes.PrecommitDone)

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

	c.logger.Info("commit a block", "hash", proposal.ProposalBlock.Header().Hash())

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
}

// Metric collecton of round change and height change.
func (c *Core) MeasureHeightRoundMetrics(round int64) {
	if round == 0 {
		// in case of height change, round changed too, so count it also.
		tctypes.TendermintRoundChangeMeter.Mark(1)
		tctypes.TendermintHeightChangeMeter.Mark(1)
	} else {
		tctypes.TendermintRoundChangeMeter.Mark(1)
	}
}

// StartRound starts a new round. if round equals to 0, it means to starts a new height
func (c *Core) StartRound(ctx context.Context, round int64) {

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
		lastBlockMined, _ := c.backend.LastCommittedProposal()
		c.setHeight(new(big.Int).Add(lastBlockMined.Number(), common.Big1))
		lastHeader := lastBlockMined.Header()
		c.committee.SetLastBlock(lastBlockMined)
		c.setLastHeader(lastHeader)
		c.lockedRound = -1
		c.lockedValue = nil
		c.validRound = -1
		c.validValue = nil
		c.messages.Reset()
		c.futureRoundChange = make(map[int64]map[common.Address]uint64)
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
}

func (c *Core) AcceptVote(roundMsgs *messageutils.RoundMessages, step tctypes.Step, hash common.Hash, msg messageutils.Message) {
	switch step {
	case tctypes.Prevote:
		roundMsgs.AddPrevote(hash, msg)
	case tctypes.Precommit:
		roundMsgs.AddPrecommit(hash, msg)
	}
}

func (c *Core) SetStep(step tctypes.Step) {
	c.logger.Debug("moving to step", "step", step.String(), "round", c.Round())
	c.step = step
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
