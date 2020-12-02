package faultdetector

import (
	"github.com/clearmatics/autonity/consensus/tendermint/bft"
	"sort"

	"github.com/clearmatics/autonity/common"
)

type interceptor struct {
	msgStore Store
}

type Verifier interface {
	Verify(m *message, s Store) error
}

type Rule byte

const (
	pv byte = iota
	pc
	p
)

type message interface {
	Round() uint
	Height() uint
	Sender() common.Address
	Type() byte
	Value() common.Hash // Block hash for a proposal,
	ValidRound() uint
}

func (i *interceptor) Intercept(msg *message) {
	// Prerequisite: msg has a valid signature and comes from validator.

	// Validation steps
	//
	// Auto incriminating
	//
	// Is Type valid one of (propose, prevote, precommit)
	//
	//
	// PN1 not defendable
	if msg.Type == Propose && msg.ValidRound() != -1 {
		//check all precommits for previous rounds from this sender are nil
		precommits := i.msgStore(msg.Height(), func(m *message) bool {
			return m.sender() == msg.Sender() && m.Type() == precommit && m.Round < msg.Round && m.Value != nilValue
		})
		if len(precommits) != 0 {
			// construct proof of bad behaviour
			precommits = append(precommits, msg)
			return Proof{
				Rule:     PN1,
				Evidence: []message{precommits[0]},
				Message:  msg,
			}
		}

	}

	// PO1
	//
	if msg.Type == Propose && msg.ValidRound() > -1 {
		//check all precommits for previous rounds from this sender are not nil or V

		// We need to find the latest round where this participant precommitted a non nil value (assuming there are no gaps between the precommits)?  Was that value
		precommits := i.msgStore(msg.Height(), func(m *message) bool {
			return m.sender() == msg.Sender() && m.Type() == precommit && m.Round == msg.Round-1 && m.Value != nilValue && m.Value != msg.Value
		})

		if len(precommits) > 0 {
			// Bad behaviour
		}

		sort.Slice(precommits, func(a, b int) {
			precommits[a].Round > precommits[b].Round()
		})

		latest := precommits[0]

		if latest.Value() != msg.Value() {

		}

		// if !i.hasQuorum(msg.Height(), prevotes){
		// 	// We have bad behaviour
		// 	// do we need apply a state: pending for this case, and it will be re-process latter. - No we cant ever check this case, we may also see the qui
		// }

		if len(precommits) != 0 {
			// construct proof of bad behaviour
			precommits = append(precommits, msg)
			return Proof{
				Rule:     PN1,
				Evidence: []message{precommits[0]},
				Message:  msg,
			}
		}

	}

	// PVT1, timeout case, cannot to proof.


	// PVN, rules for prevote for a new value.
	// PVN1, if there existed An invalid proposal at this round, then node should prevote for nil at the same round.
	// todo: assume the equivocation of proposer is address by other rules.
	if msg.Type == prevote && msg.Value() != nilValue {
		proposal := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == propose && m.Round == msg.Round() && notValid(m.Value)
		})

		if proposal != nil {
			return Proof{
				Rule:     PVN1,
				Evidence: []message{proposal},
				Message:  msg,
			}
		}
	}

	// PVN2, If there is a valid proposal V, and pi never ever precommit(locked the value) before, then pi should prevote
	// for this value or a nil in case of timeout.
	if msg.Type == prevote {
		proposal := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == propose && m.Round == msg.Round() && Valid(m.Value)
		})

		// pi never ever locked a value before.
		precommits := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == precommit && m.Round < msg.Round() && m.Value() != nilValue && m.Sender() == msg.Sender()
		})

		// the prevote of pi should be nil or same as the value proposed.
		if proposal != nil && len(precommits) > 0 && !(msg.Value() != proposal.Value() || msg.Value() != nilValue) {
			return Proof{
				Rule:     PVN2,
				Evidence: []message{proposal, precommits},
				Message:  msg,
			}
		}
	}

	// PVN3, if V is a valid proposed value, and pi locked it in the previous round, the pi should prevote for V or Nil
	if msg.Type == prevote {
		proposal := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == propose && m.Round == msg.Round() && Valid(m.Value)
		})

		// pi locked the same value before.
		precommits := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == precommit && m.Round < msg.Round() && m.Value() == proposal.Value() && m.Sender() == msg.Sender()
		})

		// the prevote of pi should be nil or same as the value proposed.
		if proposal != nil && len(precommits) > 0 && !(msg.Value() == proposal.Value() || msg.Value() == nilValue) {
			return Proof{
				Rule:     PVN3,
				Evidence: []message{proposal, precommits},
				Message:  msg,
			}
		}
	}

	// PVN4, if V is a newly proposed value, and pi last locked on a distinct value, then pi should only prevote for nil.
	if msg.Type == prevote {
		proposal := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == propose && m.Round == msg.Round() && Valid(m.Value)
		})

		// pi locked at a distinct value before.
		precommits := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == precommit && m.Round < msg.Round() && m.Value() != proposal.Value() && m.Sender() == msg.Sender()
		})

		// the prevote of pi should be nil, otherwise it break the rule
		if proposal != nil && len(precommits) > 0 && msg.Value() != nilValue {
			return Proof{
				Rule:     PVN4,
				Evidence: []message{proposal, precommits},
				Message:  msg,
			}
		}
	}

	// PVO, rules for prevote for an old value.
	// PVO1a. if V is valid proposal at round r, and pi did already locked on V at round r' < r, and pi never precommit for other Values
	// in any round between r' and r, then in round r, pi should prevote for V or nil in case of timeout.
	if msg.Type == prevote {
		// Valid V proposed at round r.
		proposal := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == propose && m.Round == msg.Round() && Valid(m.Value)
		})

		// todo: assume that no equivocation msg on msg store.
		// pi locked at the same value before round r at r'
		precommits := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == precommit && m.Round < msg.Round() && m.Value() == proposal.Value() && m.Sender() == msg.Sender()
		})

		// pi never locked for other values between round r' and r.
		otherPrecommits := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == precommit && precommits[0].Round < msg.Round() && m.Value() != proposal.Value() && m.Sender() == msg.Sender()
		})

		// the prevote of pi should be nil or V, otherwise it break the rule
		if proposal != nil && len(precommits) == 1 && len(otherPrecommits) > 0 && !(msg.Value() == proposal.Value() || msg.Value() == nilValue) {
			return Proof{
				Rule:     PVO1a,
				Evidence: []message{proposal, precommits, otherPrecommits},
				Message:  msg,
			}
		}
	}

	// PVO1b. if V is the proposed value at round r and Pi did already precommit on V at the previous round then either Pi prevotes for nil or V .
	if msg.Type == prevote {
		// Valid V proposed at round r.
		proposal := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == propose && m.Round == msg.Round() && Valid(m.Value)
		})

		// todo: assume that no equivocation msg on msg store.
		// pi locked at the same value before round r at r-1
		precommits := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == precommit && m.Round == msg.Round()-1 && m.Value() == proposal.Value() && m.Sender() == msg.Sender()
		})

		// the prevote of pi should be nil or V, otherwise it break the rule
		if proposal != nil && len(precommits) == 1 && !(msg.Value() == proposal.Value() || msg.Value() == nilValue) {
			return Proof{
				Rule:     PVO1b,
				Evidence: []message{proposal, precommits},
				Message:  msg,
			}
		}
	}

	// PVO2, if V is the proposed value at round r and Pi did already precommit on V' in the past, at round r' < r (it locked on it)
	// but there were 2f + 1 prevotes for V for a round r''' between r' and r − 1 then in round r either pi prevotes for V or nil (in case of a timeout)
	if msg.Type == prevote {
		// Valid V proposed at round r.
		proposal := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == propose && m.Round == msg.Round() && Valid(m.Value)
		})

		// todo: assume that no equivocation msg on msg store.
		// pi locked at a distinct value V' before round r at r'
		precommits := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == precommit && m.Round < msg.Round() && m.Value() != proposal.Value() && m.Sender() == msg.Sender()
		})

		// there exist quorum prevote for V for a round r''' between r' and r-1.
		r1 := precommits[0].Round()
		var prevotes = nil

		for round range (r1, msg.Round()-1) {
			results := i.msgStore(msg.Height(), func(m *message) bool {
				return m.Type() == prevote && m.Round() == round  && m.Value() == proposal.Value()
			})
			if len(results) >= bft.Quorum() {
				prevotes := results
			}
		}

		// the prevote of pi at round r should be nil or V, otherwise it break the rule
		if proposal != nil && len(precommits) == 1 && len(prevotes) >= bft.Quorum() && !(msg.Value() == proposal.Value() || msg.Value() == nilValue) {
			return Proof{
				Rule:     PVO2,
				Evidence: []message{proposal, precommits, prevotes},
				Message:  msg,
			}
		}
	}

	// PVO3. if V is the proposed value at round r and pi did already precommit on V' in the past, at round r' < r (it locked on it)
	// and r' is greater than the valid round associated with the proposal then pi prevotes for nil.
	if msg.Type == prevote {
		// Valid V proposed at round r.
		proposal := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == propose && m.Round == msg.Round() && Valid(m.Value)
		})

		// todo: assume that no equivocation msg on msg store.
		// pi locked at a distinct value V' before round r at r'
		precommits := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == precommit && m.Round < msg.Round() && m.Value() != proposal.Value() && m.Sender() == msg.Sender()
		})

		// pi should propose nil at round r if r' > validRound of the proposal
		if proposal != nil && len(precommits) == 1 && precommits[0].Round() > proposal.ValidRound() && msg.Value() != nilValue {
			return Proof{
				Rule:     PVO3,
				Evidence: []message{proposal, precommits},
				Message:  msg,
			}
		}
	}

	// PV04. if V is the proposed value at round r and pi did already precommit on V' at the previous round then pi prevotes for nil.
	if msg.Type == prevote {
		// Valid V proposed at round r.
		proposal := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == propose && m.Round == msg.Round() && Valid(m.Value)
		})

		// todo: assume that no equivocation msg on msg store.
		// pi locked at a distinct value V' before round r at r-1
		precommits := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == precommit && m.Round == msg.Round()-1 && m.Value() != proposal.Value() && m.Sender() == msg.Sender()
		})

		// the prevote of pi should be nil at round r, otherwise it break the rule
		if proposal != nil && len(precommits) == 1 && msg.Value() != nilValue {
			return Proof{
				Rule:     PV04,
				Evidence: []message{proposal, precommits},
				Message:  msg,
			}
		}
	}


	// CT,
	// CT1. Time out case, cannot to proof it.

	// C. Rules for precommit.

	// C1, precommit for a none nil value, it seems to impossible to check precommit for nil in case of timeout.
	if msg.Type == Precommit && msg.Value() != nil {
		// in the round of msg.Round, there must be a a proposal for this value && there must be a quorum number of prevote for this value.
		proposal := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == propose && m.Round == msg.Round() && m.Value == msg.Value()
		})

		prevotes := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == prevote && m.Round == msg.Round() && m.Value == msg.Value()
		})

		// todo: we assume that 2f+1 prevotes msg received before GST.
		if proposal == nil || len(prevotes) < bft.Quorum() {
			// construct proof of misbehavior of C1:
			return Proof{
				Rule:     C1,
				Evidence: []message{proposal, prevotes},
				Message:  msg,
			}
		}
	}

	// C2, If there exist quorum of prevote for nil at round r, then precommit for nil at round is valid.
	if msg.Type == Precommit ** msg.Value() == nil {
		prevotes := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == prevote ** m.Round == msg.Round() && m.Value == msg.Value()
		})

		// todo: we assume that 2f+1 prevotes msg received before GST.
		if len(prevotes) < bft.Quorum() {
			// construct proof of misbehavior of C2:
			return Proof{
				Rule:     C2,
				Evidence: []message{prevotes},
				Message:  msg,
			}
		}
	}

}

// Proposer A proposes V2
// Validator B receives proposal for V2
// Validator B checks the precommits from Proposer A, it finds a precommit for V1
// However, Proposer A locked in round 1 for V1, then unlocked and locked again in round 4 for V2.
// Validator B is yet to receive the precommit for Proposer A for V2
// Now it seems incorrect behaviour from Proposer A, but it is justified in sending the proposal for V2
