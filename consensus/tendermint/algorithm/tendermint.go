package algorithm

import (
	"encoding/hex"
	"fmt"
)

type ValueID [32]byte

func (v ValueID) String() string {
	return hex.EncodeToString(v[:3])
}

var NilValue ValueID

type NodeID [20]byte

func (n NodeID) String() string {
	return hex.EncodeToString(n[:3])
}

type Step uint8

const (
	Propose Step = iota
	Prevote
	Precommit
)

func (s Step) String() string {
	switch s {
	case Propose:
		return "Propose"
	case Prevote:
		return "Prevote"
	case Precommit:
		return "Precommit"
	default:
		panic(fmt.Sprintf("Unrecognised step value %d", s))
	}
}

func (s Step) ShortString() string {
	switch s {
	case Propose:
		return "pp"
	case Prevote:
		return "pv"
	case Precommit:
		return "pc"
	default:
		panic(fmt.Sprintf("Unrecognised step value %d", s))
	}
}

func (s Step) In(steps ...Step) bool {
	for _, step := range steps {
		if s == step {
			return true
		}
	}
	return false
}

type Timeout struct {
	TimeoutType Step
	Delay       uint
	Height      uint64
	Round       int64
}

type ConsensusMessage struct {
	MsgType    Step
	Height     uint64
	Round      int64
	Value      ValueID
	ValidRound int64 // This field only has meaning for propose step. For prevote and precommit this value is ignored.
}

func (cm *ConsensusMessage) String() string {
	return fmt.Sprintf("s:%-3s h:%-3d r:%-3d v:%-6s", cm.MsgType.ShortString(), cm.Height, cm.Round, cm.Value.String())
}

// Oracle is used to answer questions the algorithm may have about its
// state, such as 'Am I the proposer' or 'Have i reached prevote quorum
// threshold for value with id v?'
type Oracle interface {
	Valid(ValueID) bool
	MatchingProposal(*ConsensusMessage) *ConsensusMessage
	// TODO: merge the functions into QThresh since the calculation is always the same for both, instead define private
	// functions for readability.
	PrevoteQThresh(round int64, value *ValueID) bool
	PrecommitQThresh(round int64, value *ValueID) bool
	// FThresh indicates whether we have messages whose voting power exceeds
	// the failure threshold for the given round.
	FThresh(round int64) bool
	Proposer(round int64, nodeID NodeID) bool
	Height() uint64
	Value() (ValueID, error)
}

type OneShotTendermint struct {
	nodeID         NodeID
	round          int64
	step           Step
	lockedRound    int64
	lockedValue    ValueID
	validRound     int64
	validValue     ValueID
	line34Executed bool
	line36Executed bool
	line47Executed bool
	oracle         Oracle
}

func New(nodeID NodeID, oracle Oracle) *OneShotTendermint {
	return &OneShotTendermint{
		nodeID:      nodeID,
		round:       -1,
		lockedRound: -1,
		lockedValue: NilValue,
		validRound:  -1,
		validValue:  NilValue,
		oracle:      oracle,
	}
}

func (ost OneShotTendermint) height() uint64 {
	return ost.oracle.Height()
}

func (ost *OneShotTendermint) msg(msgType Step, value ValueID) *ConsensusMessage {
	cm := &ConsensusMessage{
		MsgType: msgType,
		Height:  ost.height(),
		Round:   ost.round,
		Value:   value,
	}
	if ost.step == Propose {
		cm.ValidRound = ost.validRound
	}
	return cm
}

func (ost *OneShotTendermint) timeout(timeoutType Step) *Timeout {
	return &Timeout{
		TimeoutType: timeoutType,
		Height:      ost.height(),
		Round:       ost.round,
		Delay:       1, // TODO
	}
}

// Start round takes a round to start. It then clears the first time flags and either returns a proposal
// ConsensusMessage to be broadcast, if this node is the proposer or if not, a Timeout to be scheduled.
func (ost *OneShotTendermint) StartRound(round int64) (*ConsensusMessage, *Timeout, error) {
	//println(ost.nodeID.String(), height, "isproposer", ost.oracle.Proposer(round, ost.nodeID))

	// sanity check
	switch {
	case round < 0:
		panic(fmt.Sprintf("New round cannot be less than 0. Previous round: %-3d, new round: %-3d", ost.round, round))
	case round <= ost.round:
		panic(fmt.Sprintf("New round must be more than the current round. Previous round: %-3d, new round: %-3d", ost.round, round))
	}

	// Reset first time flags
	ost.line34Executed = false
	ost.line36Executed = false
	ost.line47Executed = false

	ost.round = round
	ost.step = Propose
	if ost.oracle.Proposer(round, ost.nodeID) {
		var value ValueID
		var err error

		if ost.validValue != NilValue {
			value = ost.validValue
		} else {
			value, err = ost.oracle.Value()
			if err != nil {
				return nil, nil, err
			}
		}
		//println(a.nodeID.String(), a.height(), "returning message", value.String())
		return ost.msg(Propose, value), nil, nil
	} else { //nolint
		return nil, ost.timeout(Propose), nil
	}
}

// RoundChange indicates that the caller should initiate a round change by
// calling StartRound with the enclosed Height and Round. If Decision is set
// this indicates that a decision has been reached it will contain the proposal
// that was decided upon, Decision can only be set when Round is 0.
type RoundChange struct {
	Height   uint64 //TODO: consider removing this since this is an internal message which will not be broadcasted
	Round    int64
	Decision *ConsensusMessage //TODO: consider changing this to ValueID
}

// ReceiveMessage processes a consensus message and returns 3 values of which
// at most one can be non nil, although all can be nil, which indicates no
// state change.
//
// The values that can be returned are as follows:
//
// - *ConsensusMessage - This should be broadcast to the rest of the network,
//   including ourselves. This action can be taken asynchronously.
//
// - *RoundChange - This indicates that we need to progress to the next round,
//   and possibly next height, ultimately leading to calling StartRound with the
//   enclosed Height and Round. The call to StartRound must be executed by the
//   calling goroutine before any other call to ReceiveMessage.
//
// - *Timeout - This should be scheduled based to call the corresponding OnTimeout*
//   method after the Delay with the enclosed Height and Round. This action can be
//   taken asynchronously.
func (ost *OneShotTendermint) ReceiveMessage(cm *ConsensusMessage) (*RoundChange, *ConsensusMessage, *Timeout) {

	r := ost.round
	s := ost.step
	o := ost.oracle
	t := cm.MsgType

	// look up matching proposal, in the case of ost message with msgType
	// proposal the matching proposal is the message.
	p := o.MatchingProposal(cm)

	// Some of the checks in these upon conditions are omitted because they have already been checked.
	//
	// - We do not check height because we only execute this code when the
	// message height matches the current height.
	//
	// - We do not check whether the message comes from ost proposer since this
	// is checked before calling this method and we do not process proposals
	// from non proposers.
	//
	// The upon conditions have been re-ordered such that those with outcomes
	// that would supersede the outcome of others come before the others.
	// Specifically the upon conditions for ost given step that schedule
	// timeouts, have been moved after the upon conditions for that step that
	// would result in broadcasting ost message for ost value other than nil or
	// deciding on ost value. This ensures that we are able to return when we
	// find ost condition that has been met, because we know that the result of
	// this condition will supersede results from other later conditions that
	// may have been met. This approach will hopefully go someway to cutting
	// down unnecessary network traffic between nodes.

	// Line 22
	if t.In(Propose) && cm.Round == r && cm.ValidRound == -1 && s == Propose {
		ost.step = Prevote
		if o.Valid(cm.Value) && ost.lockedRound == -1 || ost.lockedValue == cm.Value {
			//println(ost.nodeID.String(), ost.height(), cm.String(), "line 22 val")
			return nil, ost.msg(Prevote, cm.Value), nil
		} else { //nolint
			//println(ost.nodeID.String(), ost.height(), cm.String(), "line 22 nil")
			return nil, ost.msg(Prevote, NilValue), nil
		}
	}

	// Line 28
	if t.In(Propose, Prevote) && p != nil && p.Round == r && o.PrevoteQThresh(p.ValidRound, &p.Value) && s == Propose && (p.ValidRound >= 0 && p.ValidRound < r) {
		ost.step = Prevote
		if o.Valid(p.Value) && (ost.lockedRound <= p.ValidRound || ost.lockedValue == p.Value) {
			//println(ost.nodeID.String(), ost.height(), cm.String(), "line 28 val")
			return nil, ost.msg(Prevote, p.Value), nil
		} else { //nolint
			//println(ost.nodeID.String(), ost.height(), cm.String(), "line 28 nil")
			return nil, ost.msg(Prevote, NilValue), nil
		}
	}

	////println(ost.nodeId.String(), ost.height(), t.In(Propose, Prevote), p != nil, p.Round == r, o.PrevoteQThresh(r, &p.Value), o.Valid(p.Value), s >= Prevote, !ost.line36Executed)
	// Line 36
	if t.In(Propose, Prevote) && p != nil && p.Round == r && o.PrevoteQThresh(r, &p.Value) && o.Valid(p.Value) && s >= Prevote && !ost.line36Executed {
		ost.line36Executed = true
		if s == Prevote {
			ost.lockedValue = p.Value
			ost.lockedRound = r
			ost.step = Precommit
		}
		ost.validValue = p.Value
		ost.validRound = r
		//println(ost.nodeID.String(), ost.height(), cm.String(), "line 36 val")
		return nil, ost.msg(Precommit, p.Value), nil
	}

	// Line 44
	if t.In(Prevote) && cm.Round == r && o.PrevoteQThresh(r, &NilValue) && s == Prevote {
		ost.step = Precommit
		//println(ost.nodeID.String(), ost.height(), cm.String(), "line 44 nil")
		return nil, ost.msg(Precommit, NilValue), nil
	}

	// Line 34
	if t.In(Prevote) && cm.Round == r && o.PrevoteQThresh(r, nil) && s == Prevote && !ost.line34Executed {
		ost.line34Executed = true
		//println(ost.nodeID.String(), ost.height(), cm.String(), "line 34 timeout")
		return nil, nil, ost.timeout(Prevote)
	}

	// Line 49
	if t.In(Propose, Precommit) && p != nil && o.PrecommitQThresh(p.Round, &p.Value) {
		if o.Valid(p.Value) {
			ost.lockedRound = -1
			ost.lockedValue = NilValue
			ost.validRound = -1
			ost.validValue = NilValue
		}
		//println(ost.nodeID.String(), ost.height(), cm.String(), "line 49 decide")
		// Return the decided proposal
		return &RoundChange{Height: ost.height(), Round: 0, Decision: p}, nil, nil
	}

	// Line 47
	if t.In(Precommit) && cm.Round == r && o.PrecommitQThresh(r, nil) && !ost.line47Executed {
		ost.line47Executed = true
		//println(ost.nodeID.String(), ost.height(), cm.String(), "line 47 timeout")
		return nil, nil, ost.timeout(Precommit)
	}

	// Line 55
	if cm.Round > r && o.FThresh(cm.Round) {
		// TODO account for the fact that many rounds can be skipped here. So
		// what happens to the old round messages? We don't process them, but
		// we can't remove them from the messsage store because they may be
		// used in this round in the condition at line 28. This means that we
		// only should clean the message store when there is ost height change,
		// clearing out all messages for the height.
		//println(ost.nodeID.String(), ost.height(), cm.String(), "line 55 start round")
		return &RoundChange{Height: ost.height(), Round: cm.Round}, nil, nil
	}
	//println(ost.nodeID.String(), ost.height(), cm.String(), "no condition match")
	return nil, nil, nil
}

func (ost *OneShotTendermint) OnTimeoutPropose(height uint64, round int64) *ConsensusMessage {
	if height == ost.height() && round == ost.round && ost.step == Propose {
		ost.step = Prevote
		return ost.msg(Prevote, NilValue)
	}
	return nil
}

func (ost *OneShotTendermint) OnTimeoutPrevote(height uint64, round int64) *ConsensusMessage {
	if height == ost.height() && round == ost.round && ost.step == Prevote {
		ost.step = Precommit
		return ost.msg(Precommit, NilValue)
	}
	return nil
}

func (ost *OneShotTendermint) OnTimeoutPrecommit(height uint64, round int64) *RoundChange {
	if height == ost.height() && round == ost.round {
		return &RoundChange{Height: ost.height(), Round: ost.round + 1}
	}
	return nil
}
