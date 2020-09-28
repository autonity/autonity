package main

type ValueID [32]byte

var nilValue ValueID

type NodeID [20]byte

type Step uint8

const (
	Propose Step = iota
	Prevote
	Precommit
)

func (s Step) in(steps ...Step) bool {
	for _, step := range steps {
		if s == step {
			return true
		}
	}
	return false
}

type Timeout struct {
	timeoutType Step
	delay       uint
	height      uint64
	round       int64
}

type ConsensusMessage struct {
	MsgType    Step
	Height     uint64
	Round      int64
	Value      ValueID
	ValidRound int64
}

// Oracle is used to answer questions the algorithm may have about its
// state, such as 'Am I the proposer' or 'Have i reached prevote quorum
// threshold for value with id v?'
type Oracle interface {
	Valid(ValueID) bool
	MatchingProposal(*ConsensusMessage) *ConsensusMessage
	PrevoteQThresh(round int64, value *ValueID) bool
	PrecommitQThresh(round int64, value *ValueID) bool
	FThresh(round int64, value *ValueID) bool
	Proposer(NodeID) bool
	Value() ValueID
}

type Algorithm struct {
	nodeId         NodeID
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

func (a *Algorithm) msg(msgType Step, value ValueID) *ConsensusMessage {
	cm := &ConsensusMessage{
		MsgType: a.step,
		Height:  a.height,
		Round:   a.round,
		Value:   value,
	}
	if a.step == Propose {
		cm.ValidRound = a.validRound
	}
	return cm
}

func (a *Algorithm) timeout(msgType Step) *Timeout {
	return &Timeout{
		timeoutType: Propose,
		height:      a.height,
		round:       a.round,
		delay:       1, // todo
	}
}

// Start round takes the round to start, clears the first time flags and then
// either broadcasts a proposal if this node is the proposer, or schedules a
// proposal timeout.
func (a *Algorithm) StartRound(round int64) (*ConsensusMessage, *Timeout) {
	// Reset first time flags
	a.line34Executed = false
	a.line36Executed = false
	a.line47Executed = false

	a.round = round
	a.step = Propose
	if a.oracle.Proposer(a.nodeId) {
		var v ValueID
		if a.validValue != nilValue {
			v = a.validValue
		} else {
			v = a.oracle.Value()
		}
		return a.msg(Propose, v), nil
	} else {
		return nil, a.timeout(Propose)
	}
}

// ReceiveMessage processes a consensus message and returns a ConsensusMessage
// or Timeout, indicating that either the message should be broadcast or that
// the timeout should be scheduled. At most one of the return values can be non
// nil, but it is possible for both to be nil in the case that the processed
// messge does not result in a state change.
func (a *Algorithm) ReceiveMessage(cm *ConsensusMessage) (*ConsensusMessage, *Timeout) {

	r := a.round
	s := a.step
	o := a.oracle
	t := cm.MsgType

	// look up matching proposal, in the case of a message with msgType
	// proposal the matching proposal is the message.
	p := o.MatchingProposal(cm)

	// Some of the checks in these upon conditions are omitted because they have alrady been checked.
	//
	// - We do not check height because we only execute this code when the
	// message height matches the current height.
	//
	// - We do not check whether the message comes from a proposer since this
	// is checkded before calling this method and we do not process proposals
	// from non proposers.
	//
	// The upon conditions have been re-ordered such that those with outcomes
	// that would supercede the oucome of others come before the others.
	// Specifically the upon conditions for a given step that schedule
	// timeouts, have been moved after the upon conditions for that step that
	// would result in broadcasting a message for a value other than nil or
	// deciding on a value. This ensures that we are able to return when we
	// find a condition that has been met, becuase we know that the result of
	// this condition will supercede results from other later conditions that
	// may have been met. This approach will hopefully go someway to cutting
	// down unneccesary network traffic between nodes.

	// Line 22
	if t.in(Propose) && cm.Round == r && cm.ValidRound == -1 && s == Propose {
		a.step = Prevote
		if o.Valid(cm.Value) && a.lockedRound == -1 || a.lockedValue == cm.Value {
			return a.msg(Prevote, cm.Value), nil
		} else {
			return a.msg(Prevote, nilValue), nil
		}
	}

	// Line 28
	if t.in(Propose, Prevote) && p != nil && p.Round == r && o.PrevoteQThresh(p.ValidRound, &p.Value) && s == Propose && (p.ValidRound >= 0 && p.ValidRound < r) {
		a.step = Prevote
		if o.Valid(p.Value) && (a.lockedRound <= p.ValidRound || a.lockedValue == p.Value) {
			return a.msg(Prevote, p.Value), nil
		} else {
			return a.msg(Prevote, nilValue), nil
		}
	}

	// Line 36
	if t.in(Propose, Prevote) && p != nil && p.Round == r && o.PrevoteQThresh(r, &p.Value) && o.Valid(p.Value) && s >= Prevote && !a.line36Executed {
		a.line36Executed = true
		if s == Prevote {
			a.lockedValue = p.Value
			a.lockedRound = r
			a.step = Precommit
		}
		a.validValue = p.Value
		a.validRound = r
		return a.msg(Precommit, p.Value), nil
	}

	// Line 44
	if t.in(Prevote) && cm.Round == r && o.PrevoteQThresh(r, &nilValue) && s == Prevote {
		a.step = Precommit
		return a.msg(Precommit, nilValue), nil
	}

	// Line 34
	if t.in(Prevote) && cm.Round == r && o.PrevoteQThresh(r, nil) && s == Prevote && !a.line34Executed {
		a.line34Executed = true
		return nil, a.timeout(Prevote)
	}

	// Line 49
	if t.in(Propose, Precommit) && p != nil && o.PrecommitQThresh(p.Round, &p.Value) {
		if o.Valid(p.Value) {
			// TODO commit here commit(p.Value)
			a.height++
			a.lockedRound = -1
			a.lockedValue = nilValue
			a.validRound = -1
			a.validValue = nilValue
		}
		return a.StartRound(0)
	}

	// Line 47
	if t.in(Precommit) && cm.Round == r && o.PrecommitQThresh(r, nil) && !a.line47Executed {
		a.line47Executed = true
		return nil, a.timeout(Precommit)
	}

	// Line 55
	if cm.Round > r && o.FThresh(cm.Round, nil) {
		// TODO account for the fact that many rounds can be skipped here.  so
		// what happens to the old round messages? We don't process them, but
		// we can't remove them from the cache because they may be used in this
		// round. in the conditon at line 28. This means that we only should
		// clean the message cache when there is a height change, clearing out
		// all messages for the height.
		return a.StartRound(cm.Round)
	}
	return nil, nil
}

func (a *Algorithm) onTimeoutPropose(height uint64, round int64) *ConsensusMessage {
	if height == a.height && round == a.round && a.step == Propose {
		a.step = Prevote
		return a.msg(Prevote, nilValue)
	}
	return nil
}
func (a *Algorithm) onTimeoutPrevote(height uint64, round int64) *ConsensusMessage {
	if height == a.height && round == a.round && a.step == Prevote {
		a.step = Precommit
		return a.msg(Precommit, nilValue)
	}
	return nil
}
func (a *Algorithm) onTimeoutPrecommit(height uint64, round int64) (*ConsensusMessage, *Timeout) {
	if height == a.height && round == a.round {
		return a.StartRound(a.round + 1)
	}
	return nil, nil
}
