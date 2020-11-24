package faultdetector

import (
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

}

// Proposer A proposes V2
// Validator B receives proposal for V2
// Validator B checks the precommits from Proposer A, it finds a precommit for V1
// However, Proposer A locked in round 1 for V1, then unlocked and locked again in round 4 for V2.
// Validator B is yet to receive the precommit for Proposer A for V2
// Now it seems incorrect behaviour from Proposer A, but it is justified in sending the proposal for V2
