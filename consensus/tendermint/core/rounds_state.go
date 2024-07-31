package core

import (
	"github.com/autonity/autonity/common"
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
	SetProposal(proposal *message.Propose, verified bool)

	CurRoundMessages() *message.RoundMessages
	CurRoundProposal() *message.Propose
	CurRoundPrevotesTotalPower() *big.Int
	CurRoundPrecommitsTotalPower() *big.Int
	CurRoundPrevotesPower(hash common.Hash) *big.Int
	CurRoundPrecommitsPower(hash common.Hash) *big.Int

	AddPrevote(vote *message.Prevote)
	AddPrecommit(vote *message.Precommit)

	RoundProposal(r int64) *message.Propose
	RoundPower(r int64) *big.Int
	RoundPrevotesTotalPower(r int64) *big.Int
	RoundPrecommitsTotalPower(r int64) *big.Int
	RoundPrevotesPower(r int64, hash common.Hash) *big.Int
	RoundPrecommitsPower(r int64, hash common.Hash) *big.Int

	GetOrCreate(r int64) *message.RoundMessages // todo: shall it be a internal one?
	GetRounds() []int64

	AggregatedPrevoteFor(round int64, hash common.Hash) *message.Prevote
	AggregatedPrecommitFor(round int64, hash common.Hash) *message.Precommit

	All() []message.Msg
}

type TendermintStateInterface interface {
	// round state writer functions
	StartNewHeight(h *big.Int)
	StartNewRound(r int64)
	SetStep(s Step)
	SetDecision(block *types.Block)
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

// TendermintState has the raw state of tendermint state machine, it also contains some extra state base on Autonity
// implementation of TBFT.
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

// RoundsStateImpl stores the tendermint state of a consensus instance into file system. For every successfully applied
// consensus messages, they are flushed to WAL, and for every update on the tendermint state, the state are flushed to
// WAL as well, thus the validator can recover the tendermint state from a disaster by loading the state from WAL. Note
// that the view flushed in WAL might become stale if the network grows the chain head into a higher view, in this case,
// the consensus engine should start a new height to overwrite the view of WAL on start up.
type RoundsStateImpl struct {
	TendermintState                         // state that to be flushed on update.
	curRoundMessages *message.RoundMessages // current round messages of Autonity tendermint.
	messages         *message.Map           // round messages of current height.
	db               *RoundsStateDB         // storage layer of tendermint state and round messages.
	logger           log.Logger
}

// newRoundsState, load rounds state from underlying database if there was state stored, otherwise return default state.
func newRoundsState(logger log.Logger, db ethdb.Database) RoundsState {
	// todo: jason, shall we keep the limit of max round to 99? As with WAL, validator can restore the round, thus it
	//  wouldn't be reset to 0 on start up, this limit of round number isn't a standard implementation of BFT, it might
	//  introduce live-ness issue for the consensus engine.

	// load tendermint state and rounds messages from database.
	walDB := newRoundStateDB(db)
	roundMsgs := walDB.RoundMsgsFromDB()
	lastState := walDB.GetLastTendermintState()
	state := TendermintState{
		height:                lastState.height,
		round:                 lastState.round,
		step:                  lastState.step,
		decision:              lastState.decision,
		lockedRound:           lastState.lockedRound,
		validRound:            lastState.validRound,
		sentProposal:          lastState.sentProposal,
		sentPrevote:           lastState.sentPrevote,
		sentPrecommit:         lastState.sentPrecommit,
		setValidRoundAndValue: lastState.setValidRoundAndValue,
	}

	if lastState.lockedRound != -1 {
		state.lockedValue = roundMsgs.GetOrCreate(lastState.lockedRound).Proposal().Block()
	}
	if lastState.validRound != -1 {
		state.validValue = roundMsgs.GetOrCreate(lastState.validRound).Proposal().Block()
	}

	return &RoundsStateImpl{
		TendermintState:  state,
		messages:         roundMsgs,
		curRoundMessages: roundMsgs.GetOrCreate(state.round),
		db:               walDB,
		logger:           logger,
	}
}

// round state writer functions

func (rs *RoundsStateImpl) StartNewHeight(h *big.Int) {
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

	if err := rs.db.UpdateLastRoundState(rs.TendermintState, true); err != nil {
		rs.logger.Error("failed to flush round state in WAL", "error", err)
	}

	rs.messages.Reset()
	rs.curRoundMessages = rs.messages.GetOrCreate(0)
}

func (rs *RoundsStateImpl) StartNewRound(r int64) {
	rs.curRoundMessages = rs.messages.GetOrCreate(r)

	rs.round = r
	rs.step = Propose
	rs.sentProposal = false
	rs.sentPrevote = false
	rs.sentPrecommit = false
	rs.setValidRoundAndValue = false
	if err := rs.db.UpdateLastRoundState(rs.TendermintState, false); err != nil {
		rs.logger.Error("failed to flush round state in WAL", "error", err)
	}
}

func (rs *RoundsStateImpl) SetStep(s Step) {
	rs.step = s
	if err := rs.db.UpdateLastRoundState(rs.TendermintState, false); err != nil {
		rs.logger.Error("failed to flush round state in WAL", "error", err)
	}
}

func (rs *RoundsStateImpl) SetDecision(block *types.Block) {
	rs.decision = block
	if err := rs.db.UpdateLastRoundState(rs.TendermintState, false); err != nil {
		rs.logger.Error("failed to flush round state in WAL", "error", err)
	}
}

func (rs *RoundsStateImpl) SetLockedRoundAndValue(r int64, block *types.Block) {
	rs.lockedRound = r
	rs.lockedValue = block
	if err := rs.db.UpdateLastRoundState(rs.TendermintState, false); err != nil {
		rs.logger.Error("failed to flush round state in WAL", "error", err)
	}
}

func (rs *RoundsStateImpl) SetValidRoundAndValue(r int64, block *types.Block) {
	rs.validRound = r
	rs.validValue = block
	rs.setValidRoundAndValue = true
	if err := rs.db.UpdateLastRoundState(rs.TendermintState, false); err != nil {
		rs.logger.Error("failed to flush round state in WAL", "error", err)
	}
}

func (rs *RoundsStateImpl) SetSentProposal() {
	rs.sentProposal = true
	if err := rs.db.UpdateLastRoundState(rs.TendermintState, false); err != nil {
		rs.logger.Error("failed to flush round state in WAL", "error", err)
	}
}

func (rs *RoundsStateImpl) SetSentPrevote() {
	rs.sentPrevote = true
	if err := rs.db.UpdateLastRoundState(rs.TendermintState, false); err != nil {
		rs.logger.Error("failed to flush round state in WAL", "error", err)
	}
}

func (rs *RoundsStateImpl) SetSentPrecommit() {
	rs.sentPrecommit = true
	if err := rs.db.UpdateLastRoundState(rs.TendermintState, false); err != nil {
		rs.logger.Error("failed to flush round state in WAL", "error", err)
	}
}

// message writer functions:

func (rs *RoundsStateImpl) SetProposal(proposal *message.Propose, verified bool) {
	if proposal.R() == rs.round {
		rs.curRoundMessages.SetProposal(proposal, verified)
		rs.db.AddMsg(proposal, verified)
		return
	}
	if proposal.R() < rs.round {
		rs.messages.GetOrCreate(proposal.R()).SetProposal(proposal, verified)
		rs.db.AddMsg(proposal, verified)
		return
	}
	panic("future round proposal should be in backlog")
}
func (rs *RoundsStateImpl) AddPrevote(vote *message.Prevote) {
	if vote.R() == rs.round {
		rs.curRoundMessages.AddPrevote(vote)
		rs.db.AddMsg(vote, true)
		return
	}
	if vote.R() < rs.round {
		rs.messages.GetOrCreate(vote.R()).AddPrevote(vote)
		rs.db.AddMsg(vote, true)
		return
	}
	panic("future round prevote should be in backlog")
}

func (rs *RoundsStateImpl) AddPrecommit(vote *message.Precommit) {
	if vote.R() == rs.round {
		rs.curRoundMessages.AddPrecommit(vote)
		rs.db.AddMsg(vote, true)
		return
	}
	if vote.R() < rs.round {
		rs.messages.GetOrCreate(vote.R()).AddPrecommit(vote)
		rs.db.AddMsg(vote, true)
		return
	}
	panic("future round precommit should be in backlog")
}

// state reader functions:

func (rs *RoundsStateImpl) Height() *big.Int {
	return rs.height
}

func (rs *RoundsStateImpl) Round() int64 {
	return rs.round
}

func (rs *RoundsStateImpl) Step() Step {
	return rs.step
}

func (rs *RoundsStateImpl) Decision() *types.Block {
	return rs.decision
}

func (rs *RoundsStateImpl) LockedRound() int64 {
	return rs.lockedRound
}

func (rs *RoundsStateImpl) ValidRound() int64 {
	return rs.validRound
}

func (rs *RoundsStateImpl) LockedValue() *types.Block {
	return rs.lockedValue
}

func (rs *RoundsStateImpl) ValidValue() *types.Block {
	return rs.validValue
}

func (rs *RoundsStateImpl) SentProposal() bool {
	return rs.sentProposal
}

func (rs *RoundsStateImpl) SentPrevote() bool {
	return rs.sentPrevote
}

func (rs *RoundsStateImpl) SentPrecommit() bool {
	return rs.sentPrecommit
}

func (rs *RoundsStateImpl) ValidRoundAndValueSet() bool {
	return rs.setValidRoundAndValue
}

// round messages reader functions:

func (rs *RoundsStateImpl) CurRoundMessages() *message.RoundMessages {
	return rs.curRoundMessages
}

func (rs *RoundsStateImpl) CurRoundProposal() *message.Propose {
	return rs.curRoundMessages.Proposal()
}

func (rs *RoundsStateImpl) CurRoundPrevotesTotalPower() *big.Int {
	return rs.curRoundMessages.PrevotesTotalPower()
}

func (rs *RoundsStateImpl) CurRoundPrecommitsTotalPower() *big.Int {
	return rs.curRoundMessages.PrecommitsTotalPower()
}

func (rs *RoundsStateImpl) CurRoundPrevotesPower(hash common.Hash) *big.Int {
	return rs.curRoundMessages.PrevotesPower(hash)
}

func (rs *RoundsStateImpl) CurRoundPrecommitsPower(hash common.Hash) *big.Int {
	return rs.curRoundMessages.PrecommitsPower(hash)
}

func (rs *RoundsStateImpl) RoundProposal(r int64) *message.Propose {
	if r == rs.round {
		return rs.curRoundMessages.Proposal()
	}
	if r < rs.round {
		return rs.messages.GetOrCreate(r).Proposal()
	}
	panic("future round proposal should be in backlog")
}

func (rs *RoundsStateImpl) RoundPower(r int64) *big.Int {
	if r == rs.round {
		return rs.curRoundMessages.Power().Power()
	}
	if r < rs.round {
		return rs.messages.GetOrCreate(r).Power().Power()
	}
	panic("round power cannot be calculated for future round messages")
}

func (rs *RoundsStateImpl) RoundPrevotesTotalPower(r int64) *big.Int {
	if r == rs.round {
		return rs.curRoundMessages.PrevotesTotalPower()
	}
	if r < rs.round {
		return rs.messages.GetOrCreate(r).PrevotesTotalPower()
	}
	panic("prevote power cannot be calculated for future round messages")
}

func (rs *RoundsStateImpl) RoundPrecommitsTotalPower(r int64) *big.Int {
	if r == rs.round {
		return rs.curRoundMessages.PrecommitsTotalPower()
	}
	if r < rs.round {
		return rs.messages.GetOrCreate(r).PrecommitsTotalPower()
	}
	panic("precommit power cannot be calculated for future round messages")
}

func (rs *RoundsStateImpl) RoundPrevotesPower(r int64, hash common.Hash) *big.Int {
	if r == rs.round {
		return rs.curRoundMessages.PrevotesPower(hash)
	}
	if r < rs.round {
		return rs.messages.GetOrCreate(r).PrevotesPower(hash)
	}
	panic("prevote power cannot be calculated for future round messages")
}

func (rs *RoundsStateImpl) RoundPrecommitsPower(r int64, hash common.Hash) *big.Int {
	if r == rs.round {
		return rs.curRoundMessages.PrecommitsPower(hash)
	}
	if r < rs.round {
		return rs.messages.GetOrCreate(r).PrecommitsPower(hash)
	}
	panic("precommit power cannot be calculated for future round messages")
}

func (rs *RoundsStateImpl) GetOrCreate(r int64) *message.RoundMessages {
	return rs.messages.GetOrCreate(r)
}

func (rs *RoundsStateImpl) GetRounds() []int64 {
	return rs.messages.GetRounds()
}

func (rs *RoundsStateImpl) AggregatedPrevoteFor(round int64, hash common.Hash) *message.Prevote {
	if round == rs.round {
		return rs.curRoundMessages.PrevoteFor(hash)
	}
	if round < rs.round {
		return rs.messages.GetOrCreate(round).PrevoteFor(hash)
	}
	panic("cannot aggregate future round prevotes")
}

func (rs *RoundsStateImpl) AggregatedPrecommitFor(round int64, hash common.Hash) *message.Precommit {
	if round == rs.round {
		return rs.curRoundMessages.PrecommitFor(hash)
	}
	if round < rs.round {
		return rs.messages.GetOrCreate(round).PrecommitFor(hash)
	}
	panic("cannot aggregate future round precommits")
}

func (rs *RoundsStateImpl) All() []message.Msg {
	return rs.messages.All()
}
