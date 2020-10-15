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
	ValidRound int64
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
	PrevoteQThresh(round int64, value *ValueID) bool
	PrecommitQThresh(round int64, value *ValueID) bool
	// FThresh indicates whether we have messages whose voting power exceeds
	// the failure threshold for the given round.
	FThresh(round int64) bool
	Proposer(round int64, nodeID NodeID) bool
}

type Algorithm struct {
	nodeID         NodeID
	height         uint64
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

func New(nodeID NodeID, oracle Oracle) *Algorithm {
	return &Algorithm{
		nodeID:      nodeID,
		lockedRound: -1,
		lockedValue: NilValue,
		validRound:  -1,
		validValue:  NilValue,
		oracle:      oracle,
	}
}

func (a *Algorithm) msg(msgType Step, value ValueID) *ConsensusMessage {
	cm := &ConsensusMessage{
		MsgType: msgType,
		Height:  a.height,
		Round:   a.round,
		Value:   value,
	}
	if a.step == Propose {
		cm.ValidRound = a.validRound
	}
	return cm
}

func (a *Algorithm) timeout(timeoutType Step) *Timeout {
	return &Timeout{
		TimeoutType: timeoutType,
		Height:      a.height,
		Round:       a.round,
		Delay:       1, // TODO
	}
}

// Start round takes the height and round to start as well as a potential value
// to propose. It then clears the first time flags and either returns a Result
// with a proposal to be broadcast if this node is the proposer, or a timeout
// to be scheduled.
func (a *Algorithm) StartRound(height uint64, round int64, value ValueID) *Result {
	println(a.nodeID.String(), height, "isproposer", a.oracle.Proposer(round, a.nodeID))
	// Reset first time flags
	a.line34Executed = false
	a.line36Executed = false
	a.line47Executed = false

	a.height = height
	a.round = round
	a.step = Propose
	if a.oracle.Proposer(round, a.nodeID) {
		if a.validValue != NilValue {
			value = a.validValue
		}
		println(a.nodeID.String(), height, "returning message", value.String())
		return &Result{Broadcast: a.msg(Propose, value)}
	} else { //nolint
		return &Result{Schedule: a.timeout(Propose)}
	}
}

// Result is returned from the methods of Algorithm to indicate the outcome of
// processing and what steps should be taken. Only one of the three fields may
// be set. If StartRound is set it indicates that the caller should call
// StartRound, if Broadcast is set it indicates that the caller should
// broadcast the ConsensusMessage, including sending it to itself and if
// Schedule is set it indicates that the caller should schedule the Timeout.
// Broadcasting and scheduling may be done asynchronously, but starting the next
// round must be done in the calling goroutine so that no other messages are
// processed by Algorithm between the caller receiving the Result and calling
// StartRound.
type Result struct {
	StartRound *RoundChange
	Broadcast  *ConsensusMessage
	Schedule   *Timeout
}

// RoundChange indicates that the caller should initiate a round change by
// calling StartRound with the enclosed Height and Round. If Decision is set
// this indicates that a decision has been reached it will contain the proposal
// that was decided upon, Decision can only be set when Round is 0.
type RoundChange struct {
	Height   uint64
	Round    int64
	Decision *ConsensusMessage
}

// ReceiveMessage processes a consensus message and returns a Result if a state
// change has taken place or nil if no state change has occurred.
func (a *Algorithm) ReceiveMessage(cm *ConsensusMessage) *Result {

	r := a.round
	s := a.step
	o := a.oracle
	t := cm.MsgType

	// look up matching proposal, in the case of a message with msgType
	// proposal the matching proposal is the message.
	p := o.MatchingProposal(cm)

	// Some of the checks in these upon conditions are omitted because they have already been checked.
	//
	// - We do not check height because we only execute this code when the
	// message height matches the current height.
	//
	// - We do not check whether the message comes from a proposer since this
	// is checked before calling this method and we do not process proposals
	// from non proposers.
	//
	// The upon conditions have been re-ordered such that those with outcomes
	// that would supersede the outcome of others come before the others.
	// Specifically the upon conditions for a given step that schedule
	// timeouts, have been moved after the upon conditions for that step that
	// would result in broadcasting a message for a value other than nil or
	// deciding on a value. This ensures that we are able to return when we
	// find a condition that has been met, because we know that the result of
	// this condition will supersede results from other later conditions that
	// may have been met. This approach will hopefully go someway to cutting
	// down unnecessary network traffic between nodes.

	// Line 22
	if t.In(Propose) && cm.Round == r && cm.ValidRound == -1 && s == Propose {
		a.step = Prevote
		if o.Valid(cm.Value) && a.lockedRound == -1 || a.lockedValue == cm.Value {
			println(a.nodeID.String(), a.height, cm.String(), "line 22 val")
			return &Result{Broadcast: a.msg(Prevote, cm.Value)}
		} else { //nolint
			println(a.nodeID.String(), a.height, cm.String(), "line 22 nil")
			return &Result{Broadcast: a.msg(Prevote, NilValue)}
		}
	}

	// Line 28
	if t.In(Propose, Prevote) && p != nil && p.Round == r && o.PrevoteQThresh(p.ValidRound, &p.Value) && s == Propose && (p.ValidRound >= 0 && p.ValidRound < r) {
		a.step = Prevote
		if o.Valid(p.Value) && (a.lockedRound <= p.ValidRound || a.lockedValue == p.Value) {
			println(a.nodeID.String(), a.height, cm.String(), "line 28 val")
			return &Result{Broadcast: a.msg(Prevote, p.Value)}
		} else { //nolint
			println(a.nodeID.String(), a.height, cm.String(), "line 28 nil")
			return &Result{Broadcast: a.msg(Prevote, NilValue)}
		}
	}

	//println(a.nodeId.String(), a.height, t.In(Propose, Prevote), p != nil, p.Round == r, o.PrevoteQThresh(r, &p.Value), o.Valid(p.Value), s >= Prevote, !a.line36Executed)
	// Line 36
	if t.In(Propose, Prevote) && p != nil && p.Round == r && o.PrevoteQThresh(r, &p.Value) && o.Valid(p.Value) && s >= Prevote && !a.line36Executed {
		a.line36Executed = true
		if s == Prevote {
			a.lockedValue = p.Value
			a.lockedRound = r
			a.step = Precommit
		}
		a.validValue = p.Value
		a.validRound = r
		println(a.nodeID.String(), a.height, cm.String(), "line 36 val")
		return &Result{Broadcast: a.msg(Precommit, p.Value)}
	}

	// Line 44
	if t.In(Prevote) && cm.Round == r && o.PrevoteQThresh(r, &NilValue) && s == Prevote {
		a.step = Precommit
		println(a.nodeID.String(), a.height, cm.String(), "line 44 nil")
		return &Result{Broadcast: a.msg(Precommit, NilValue)}
	}

	// Line 34
	if t.In(Prevote) && cm.Round == r && o.PrevoteQThresh(r, nil) && s == Prevote && !a.line34Executed {
		a.line34Executed = true
		println(a.nodeID.String(), a.height, cm.String(), "line 34 timeout")
		return &Result{Schedule: a.timeout(Prevote)}
	}

	// Line 49
	if t.In(Propose, Precommit) && p != nil && o.PrecommitQThresh(p.Round, &p.Value) {
		if o.Valid(p.Value) {
			a.height++
			a.lockedRound = -1
			a.lockedValue = NilValue
			a.validRound = -1
			a.validValue = NilValue
		}
		println(a.nodeID.String(), a.height, cm.String(), "line 49 decide")
		// Return the decided proposal
		return &Result{StartRound: &RoundChange{Height: a.height, Round: 0, Decision: p}}
	}

	// Line 47
	if t.In(Precommit) && cm.Round == r && o.PrecommitQThresh(r, nil) && !a.line47Executed {
		a.line47Executed = true
		println(a.nodeID.String(), a.height, cm.String(), "line 47 timeout")
		return &Result{Schedule: a.timeout(Precommit)}
	}

	// Line 55
	if cm.Round > r && o.FThresh(cm.Round) {
		// TODO account for the fact that many rounds can be skipped here. So
		// what happens to the old round messages? We don't process them, but
		// we can't remove them from the messsage store because they may be
		// used in this round in the condition at line 28. This means that we
		// only should clean the message store when there is a height change,
		// clearing out all messages for the height.
		println(a.nodeID.String(), a.height, cm.String(), "line 55 start round")
		return &Result{StartRound: &RoundChange{Height: a.height, Round: cm.Round}}
	}
	println(a.nodeID.String(), a.height, cm.String(), "no condition match")
	return nil
}

func (a *Algorithm) OnTimeoutPropose(height uint64, round int64) *Result {
	if height == a.height && round == a.round && a.step == Propose {
		a.step = Prevote
		return &Result{Broadcast: a.msg(Prevote, NilValue)}
	}
	return nil
}

func (a *Algorithm) OnTimeoutPrevote(height uint64, round int64) *Result {
	if height == a.height && round == a.round && a.step == Prevote {
		a.step = Precommit
		return &Result{Broadcast: a.msg(Precommit, NilValue)}
	}
	return nil
}

func (a *Algorithm) OnTimeoutPrecommit(height uint64, round int64) *Result {
	if height == a.height && round == a.round {
		return &Result{StartRound: &RoundChange{Height: a.height, Round: a.round + 1}}
	}
	return nil
}
