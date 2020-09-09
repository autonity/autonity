package main

type ValueID [32]byte

var nilValue ValueID

type Step uint8

func (s Step) in(steps ...Step) bool {
	for _, step := range steps {
		if s == step {
			return true
		}
	}
	return false
}

const (
	Propose Step = iota
	Prevote
	Precommit
)

type ConsensusMessage struct {
	MsgType    Step
	Height     uint64
	Round      int64
	Value      ValueID
	ValidRound int64
}

type Sender interface {
	Send(cm *ConsensusMessage)
}

type Algorithm struct {
	height         uint64
	round          int64
	step           Step
	lockedRound    int64
	lockedValue    ValueID
	validRound     int64
	validValue     ValueID
	sender         Sender
	line34Executed bool
	line36Executed bool
	line47Executed bool
}

type Oracle interface {
	Valid(ValueID) bool
	MatchingProposal(*ConsensusMessage) *ConsensusMessage
	PrevoteQuorum(round int64, value *ValueID) bool
	PrecommitQuorum(round int64, value *ValueID) bool
	FailThreshold(round int64, value *ValueID) bool
}

func (a *Algorithm) send(msgType Step, value ValueID) {
	cm := &ConsensusMessage{
		MsgType: msgType,
		Height:  a.height,
		Round:   a.round,
		Value:   value,
	}
	if msgType == Propose {
		cm.ValidRound = a.validRound
	}
	a.sender.Send(cm)
}

func (a *Algorithm) ReceiveMessage(cm *ConsensusMessage, o Oracle) {

	h := a.height
	r := a.round
	s := a.step
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

	// Line 22
	if t.in(Propose) && cm.Round == r && cm.ValidRound == -1 && s == Propose {
		if o.Valid(cm.Value) && a.lockedRound == -1 || a.lockedValue == cm.Value {
			a.send(Prevote, cm.Value)
		} else {
			a.send(Prevote, nilValue)
		}
	}

	// Line 28
	if t.in(Propose, Prevote) && p != nil && p.Round == r && o.PrevoteQuorum(p.ValidRound, &p.Value) && s == Propose && (p.ValidRound >= 0 && p.ValidRound < r) {
		if o.Valid(p.Value) && (a.lockedRound <= p.ValidRound || a.lockedValue == p.Value) {
			a.send(Prevote, p.Value)
		} else {
			a.send(Prevote, nilValue)
		}
	}

	// Line 34
	if t.in(Prevote) && cm.Round == r && o.PrevoteQuorum(r, nil) && s == Prevote && !a.line34Executed {
		//c.prevoteTimeout.scheduleTimeout(c.timeoutPrevote(r), r, h, c.onTimeoutPrecommit)
	}

	// Line 36
	if t.in(Propose, Prevote) && p != nil && p.Round == r && o.PrevoteQuorum(r, &p.Value) && o.Valid(p.Value) && s >= Prevote && !a.line36Executed {
		if s == Prevote {
			a.lockedValue = p.Value
			a.lockedRound = r
			a.send(Precommit, p.Value)
			s = Precommit // TODO set steps in all situations where we set the steps
			a.step = Precommit
		}
		a.validValue = p.Value
		a.validRound = r
	}

	// Line 44
	if t.in(Prevote) && cm.Round == r && o.PrevoteQuorum(r, &nilValue) && s == Prevote {
		a.send(Precommit, p.Value)
		s = Precommit
		a.step = Precommit
	}

	// Line 47
	if t.in(Precommit) && cm.Round == r && o.PrecommitQuorum(r, nil) && !a.line47Executed {
		//c.precommitTimeout.scheduleTimeout(c.timeoutPrecommit(r), r, h, c.onTimeoutPrecommit) // TODO handle the timers
	}

	// Line 49
	if t.in(Propose, Precommit) && p != nil && o.PrecommitQuorum(p.Round, &p.Value) {
		if o.Valid(p.Value) {
			// TODO commit here commit(p.Value)
			a.height++
			a.lockedRound = -1
			a.lockedValue = nilValue
			a.validRound = -1
			a.validValue = nilValue
		}

		// Not quite sure how to start the round nicely
		// need to ensure that we don't stack overflow in the case that the
		// next height messages are sufficient for consensus when we
		// process them and so on and so on.  So I need to set the start
		// round states and then queue the messages for processing. And I
		// need to ensure that I get a list of messages to process in an
		// atomic step from the msg cache so that I don't end up trying to
		// process the same message twice.
	}

	// Line 55
	if cm.Round > r && o.FailThreshold(cm.Round, nil) {
		// StartRound(cm.Round) // TODO
	}
}

func (a *Algorithm) SendMessage() {
}
