package core

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/ethdb"
	"github.com/autonity/autonity/log"
	"math/big"
)

type RoundsState interface {
	TendermintStateInterface
	RoundMessageInterface
}

type RoundMessageInterface interface {
	Messages() *message.Map
	AllMessages() []message.Msg
	CurRoundMessages() *message.RoundMessages

	SetProposal(proposal *message.Propose, verified bool)
	AddPrevote(vote *message.Prevote)
	AddPrecommit(vote *message.Precommit)
	GetOrCreate(r int64) *message.RoundMessages
}

type TendermintStateInterface interface {
	// round state writer functions
	StartNewHeight(h *big.Int)
	StartNewRound(r int64)
	SetStep(s Step)
	SetDecision(block *types.Block, round int64)
	SetLockedRoundAndValue(r int64, block *types.Block)
	SetValidRoundAndValue(r int64, block *types.Block)
	SetSentProposal()
	SetSentPrevote()
	SetSentPrecommit()

	// round state reader functions
	Height() *big.Int
	Round() int64
	Step() Step
	Decision() *types.Block
	DecisionRound() int64
	LockedRound() int64
	ValidRound() int64
	LockedValue() *types.Block
	ValidValue() *types.Block
	SentProposal() bool
	SentPrevote() bool
	SentPrecommit() bool
	ValidRoundAndValueSet() bool
}

const (
	garbageCollectionInterval = 60
)

var nilValue = common.Hash{}

// TendermintState has the raw state of tendermint state machine,
// it also contains some extra state base on Autonity implementation of TBFT.
type TendermintState struct {
	// raw states of tendermint, the decision is recorded to recommit the decision if one failed to commit it during a
	// disaster scenario, this protects the safety of the consensus protocol.
	height   *big.Int
	round    int64
	step     Step
	decision *types.Block

	lockedRound int64
	validRound  int64
	lockedValue *types.Block
	validValue  *types.Block

	// extra states of Autonity TBFT implementation.
	sentProposal          bool
	sentPrevote           bool
	sentPrecommit         bool
	setValidRoundAndValue bool
}

// TendermintStateImpl stores the tendermint state of a consensus instance into file system. For every successfully applied
// consensus messages, they are flushed to WAL, and for every update on the tendermint state, the state are flushed to
// WAL as well, thus the validator can recover the tendermint state from a disaster by loading the state from WAL. Note
// that the view flushed in WAL might become stale if the network grows the chain head into a higher view, in this case,
// the consensus engine should start a new height to overwrite the view of WAL on start up.
type TendermintStateImpl struct {
	TendermintState                         // in-memory state that to be flushed on update.
	curRoundMessages *message.RoundMessages // in-memory current round messages of Autonity tendermint.
	messages         *message.Map           // in-memory round messages of current height.
	roundOfDecision  int64                  // in-memory round on which the decision is committed with, enriched from WAL msgs.
	db               *TendermintStateDB     // storage layer of tendermint state and round messages.
	logger           log.Logger
}

// newTendermintState, load rounds state from underlying database if there was state stored, otherwise return default state.
func newTendermintState(logger log.Logger, db ethdb.Database, chain consensus.ChainHeaderReader) RoundsState {

	// if the WAL logging is not enabled, return default state on start up.
	// the default state will be overwritten on consensus engine start phase by respecting chain height.
	if db == nil {
		return newStateWithoutWAL(logger)
	}

	// load tendermint state and rounds messages from database.
	walDB := newTendermintStateDB(logger, db)
	roundMsgs := walDB.RoundMsgsFromDB(chain)
	lastState := walDB.GetLastTendermintState()
	state := TendermintState{
		height:                new(big.Int).SetUint64(lastState.Height),
		round:                 int64(lastState.Round),
		step:                  lastState.Step,
		lockedRound:           int64(lastState.LockedRound),
		validRound:            int64(lastState.ValidRound),
		sentProposal:          lastState.SentProposal,
		sentPrevote:           lastState.SentPrevote,
		sentPrecommit:         lastState.SentPrecommit,
		setValidRoundAndValue: lastState.SetValidRoundAndValue,
	}

	if lastState.IsLockedRoundNil {
		state.lockedRound = -1
	}

	if lastState.IsValidRoundNil {
		state.validRound = -1
	}

	// by according to tendermint paper: https://arxiv.org/pdf/1807.04938, page-6, line-36 to line-43.
	// the locked value and valid value are the in the proposal of the corresponding locked round and valid round.
	if state.lockedRound != -1 && roundMsgs.GetOrCreate(state.lockedRound).Proposal() != nil {
		state.lockedValue = roundMsgs.GetOrCreate(state.lockedRound).Proposal().Block()
	}
	if state.validRound != -1 && roundMsgs.GetOrCreate(state.validRound).Proposal() != nil {
		state.validValue = roundMsgs.GetOrCreate(state.validRound).Proposal().Block()
	}

	roundOfDecision := int64(-1)
	// load the decision from round messages, as the decision is not always be in the last round,
	// we need to iterate round messages' proposal to find it.
	if lastState.Decision != nilValue {
		allRounds := roundMsgs.GetRounds()
		for _, r := range allRounds {
			proposal := roundMsgs.GetOrCreate(r).Proposal()
			if proposal != nil && proposal.Block().Hash() == lastState.Decision {
				state.decision = proposal.Block()
				roundOfDecision = proposal.R()
				break
			}
		}
	}

	return &TendermintStateImpl{
		TendermintState:  state,
		roundOfDecision:  roundOfDecision,
		messages:         roundMsgs,
		curRoundMessages: roundMsgs.GetOrCreate(state.round),
		db:               walDB,
		logger:           logger,
	}
}

// newStateWithoutWAL returns a default state of tendermint.
func newStateWithoutWAL(logger log.Logger) *TendermintStateImpl {
	state := TendermintState{
		height:      common.Big0,
		round:       0,
		step:        Propose,
		decision:    nil,
		lockedRound: -1,
		validRound:  -1,
		lockedValue: nil,
		validValue:  nil,

		// extra states of Autonity TBFT implementation.
		sentProposal:          false,
		sentPrevote:           false,
		sentPrecommit:         false,
		setValidRoundAndValue: false,
	}

	roundMsgs := message.NewMap()
	return &TendermintStateImpl{
		TendermintState:  state,
		roundOfDecision:  -1,
		messages:         roundMsgs,
		curRoundMessages: roundMsgs.GetOrCreate(0),
		db:               nil,
		logger:           logger,
	}
}

// round state writer functions

func (rs *TendermintStateImpl) StartNewHeight(h *big.Int) {
	rs.height = h
	rs.round = 0
	rs.step = Propose
	rs.decision = nil

	rs.lockedRound = -1
	rs.lockedValue = nil
	rs.validRound = -1
	rs.validValue = nil

	rs.sentProposal = false
	rs.sentProposal = false
	rs.sentPrecommit = false
	rs.setValidRoundAndValue = false

	rs.roundOfDecision = -1
	rs.messages.Reset()
	rs.curRoundMessages = rs.messages.GetOrCreate(0)

	rs.tryLogStateTransition(true)
}

func (rs *TendermintStateImpl) StartNewRound(r int64) {
	rs.curRoundMessages = rs.messages.GetOrCreate(r)

	rs.round = r
	rs.step = Propose
	rs.sentProposal = false
	rs.sentPrevote = false
	rs.sentPrecommit = false
	rs.setValidRoundAndValue = false
	rs.tryLogStateTransition(false)
}

func (rs *TendermintStateImpl) SetStep(s Step) {
	rs.step = s
	rs.tryLogStateTransition(false)
}

func (rs *TendermintStateImpl) SetDecision(block *types.Block, r int64) {
	rs.decision = block
	rs.roundOfDecision = r
	rs.tryLogStateTransition(false)
}

func (rs *TendermintStateImpl) SetLockedRoundAndValue(r int64, block *types.Block) {
	rs.lockedRound = r
	rs.lockedValue = block
	rs.tryLogStateTransition(false)
}

func (rs *TendermintStateImpl) SetValidRoundAndValue(r int64, block *types.Block) {
	rs.validRound = r
	rs.validValue = block
	rs.setValidRoundAndValue = true
	rs.tryLogStateTransition(false)
}

func (rs *TendermintStateImpl) SetSentProposal() {
	rs.sentProposal = true
	rs.tryLogStateTransition(false)
}

func (rs *TendermintStateImpl) SetSentPrevote() {
	rs.sentPrevote = true
	rs.tryLogStateTransition(false)
}

func (rs *TendermintStateImpl) SetSentPrecommit() {
	rs.sentPrecommit = true
	rs.tryLogStateTransition(false)
}

// message writer functions:

func (rs *TendermintStateImpl) SetProposal(proposal *message.Propose, verified bool) {
	if proposal.R() == rs.round {
		rs.curRoundMessages.SetProposal(proposal, verified)
		rs.tryLogMsg(proposal, verified)
		return
	}
	if proposal.R() < rs.round {
		rs.messages.GetOrCreate(proposal.R()).SetProposal(proposal, verified)
		rs.tryLogMsg(proposal, verified)
		return
	}
	panic("future round proposal should be in backlog")
}
func (rs *TendermintStateImpl) AddPrevote(vote *message.Prevote) {
	if vote.R() == rs.round {
		rs.curRoundMessages.AddPrevote(vote)
		rs.tryLogMsg(vote, true)
		return
	}
	if vote.R() < rs.round {
		rs.messages.GetOrCreate(vote.R()).AddPrevote(vote)
		rs.tryLogMsg(vote, true)
		return
	}
	panic("future round prevote should be in backlog")
}

func (rs *TendermintStateImpl) AddPrecommit(vote *message.Precommit) {
	if vote.R() == rs.round {
		rs.curRoundMessages.AddPrecommit(vote)
		rs.tryLogMsg(vote, true)
		return
	}
	if vote.R() < rs.round {
		rs.messages.GetOrCreate(vote.R()).AddPrecommit(vote)
		rs.tryLogMsg(vote, true)
		return
	}
	panic("future round precommit should be in backlog")
}

func (rs *TendermintStateImpl) tryLogStateTransition(startNewHeight bool) {
	// log state transition only on WAL is enabled.
	if rs.db != nil {
		if err := rs.db.UpdateLastRoundState(&rs.TendermintState, startNewHeight); err != nil {
			rs.logger.Error("failed to flush state transition in WAL", "error", err)
		}
	}
}

func (rs *TendermintStateImpl) tryLogMsg(msg message.Msg, verified bool) {
	// flush message only on WAL is enabled.
	if rs.db != nil {
		if err := rs.db.AddMsg(msg, verified); err != nil {
			rs.logger.Error("failed to flush msg in WAL", "error", err)
		}
	}
}

// state reader functions:

func (rs *TendermintStateImpl) Height() *big.Int {
	return rs.height
}

func (rs *TendermintStateImpl) Round() int64 {
	return rs.round
}

func (rs *TendermintStateImpl) Step() Step {
	return rs.step
}

func (rs *TendermintStateImpl) DecisionRound() int64 {
	return rs.roundOfDecision
}

func (rs *TendermintStateImpl) Decision() *types.Block {
	return rs.decision
}

func (rs *TendermintStateImpl) LockedRound() int64 {
	return rs.lockedRound
}

func (rs *TendermintStateImpl) ValidRound() int64 {
	return rs.validRound
}

func (rs *TendermintStateImpl) LockedValue() *types.Block {
	return rs.lockedValue
}

func (rs *TendermintStateImpl) ValidValue() *types.Block {
	return rs.validValue
}

func (rs *TendermintStateImpl) SentProposal() bool {
	return rs.sentProposal
}

func (rs *TendermintStateImpl) SentPrevote() bool {
	return rs.sentPrevote
}

func (rs *TendermintStateImpl) SentPrecommit() bool {
	return rs.sentPrecommit
}

func (rs *TendermintStateImpl) ValidRoundAndValueSet() bool {
	return rs.setValidRoundAndValue
}

// round messages reader functions:

func (rs *TendermintStateImpl) Messages() *message.Map {
	return rs.messages
}

func (rs *TendermintStateImpl) CurRoundMessages() *message.RoundMessages {
	return rs.curRoundMessages
}

func (rs *TendermintStateImpl) GetOrCreate(r int64) *message.RoundMessages {
	return rs.messages.GetOrCreate(r)
}

func (rs *TendermintStateImpl) AllMessages() []message.Msg {
	return rs.messages.All()
}
