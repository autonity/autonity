package faultdetector

import (
	"sort"

	"github.com/clearmatics/autonity/consensus/tendermint/bft"
	"github.com/graph-gophers/graphql-go/internal/validation"

	"github.com/clearmatics/autonity/common"
)

type Rule uint8
const (
	PN1 Rule = iota
	PO1
	PVN2
)

type Proof struct {
	Rule Rule
	Evidence []message
	Message message
}

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

type rule interface {
	// Run this whenever we receive a message
	Run(height uint64, store Store) (proofOfMisBehavior, Accusation)
}

// Rules read right to left (find  the right and look for the left)
//
// Rules should be evealuated such that we check all paossible instances and if
// we can't find a single instance that passes then we consider the rule
// failed.
//
// There are 2 types of provable misbehaviors.
// 1. Conflicting messages from a single participant
// 2. A message that conflicts with a quorum of prevotes.
// (precommit for differing value in same round as the prevotes or proposal for an
// old value where in each prior round we can see a quorum of precommits for a distinct value.)
//
// We can ignore rules that have nil on the right hand side.

type message interface {
	Round() uint
	Height() uint
	Sender() common.Address
	Type() byte
	Value() common.Hash // Block hash for a proposal,
	ValidRound() uint
}


func deferProcessing(msg *message){
	timer.executeAfter(msg, delta, func() {
		
	})
}


func (*interceptor) Process(height uint64) ([]*proofOfMisBehavior, []*Accusation){
	// We should be here at time t = timestamp(h+1) + delta

	var proofs []*proofOfMisBehavior
	var accusations []*Accusation

	// PN1 not defendable
	proposalsNew := i.msgStore(height, func(m *message) bool {
		return m.Type() == proposal && m.ValidRound == -1 
	})
	
	for i, proposal := range proposalsNew {
		//check all precommits for previous rounds from this sender are nil
		precommits := i.msgStore(height, func(m *message) bool {
			return m.sender() == proposal.Sender() && m.Type() == precommit && m.Round < proposal.Round && m.Value != nilValue
		})
		if len(precommits) != 0 {
			// construct proof of bad behaviour
			proof := &Proof{
				Rule:     PN1,
				Evidence: []message{precommits[0]},
				Message:  msg,
			}
			proofs = append(proofs, &proof)
		}
	}
	
	// (Mr′<r,PV) ∧ (Mr′,PC|pi) ∧ (Mr′<r′′<r,PC|pi)∗ ⇒ (Mr,P|pi)
	// [#(V) ≥ 2f+1] ∧ [nil ∨ V v ⊥] ∧ [nil v ⊥] ⇒ [V]

	// pv(V, r') pc(V, r') pc(V', r'') -> p(V, r)

	// pv(V, r'), pv(nil v V, r'')...pv(nil v V, r''') ---> P(V', r)

	// PO1 Checking order:
	// 1. Start from r' = r-1
	// Check for proveable misbehaviour
	// 1. r' = r-1 and pc(V', r') --> p(V, r)
	// Here we need 2f+1 prevotes for all rounds and atleast one of them must be
	// for V', then we can create a proof of misbehaviour for the proposal for V
	// 3. pv(V' v nil, 0 <= r' < r)* --> P(V, r)
	// 4. Any precommit for a Value between the valid round and the proposal round.

	// We can see the valid round in the proposal!

	// Check for accusable misbehaviour
	// 
	//

	// PO1 Old
	// proposalsOld := i.msgStore(height, func(m *message) bool {
	// 	return m.Type() == proposal && m.ValidRound > -1 
	// })
	
	// for i, proposal := range proposalsOld {
		
	// 	// Check from the preceeding round backwards
	// 	for j := proposal.Round()-1; j >= 0; j-- {
	// 		prevotes := i.msgStore(height, func(m *message) bool {
	// 			return m.Type() == prevotes && m.Round == j
	// 		})

	// 		// Term 2
				



	// 		// Term 1
	// 		if len(prevotes) > thresh {
	// 			precommit := i.msgStore(height, func(m *message) bool {
	// 				return m.Type() == Precommit && m.Round == j && m.Sender == proposal.Sender()
	// 			})
	// 			// Term 2
	// 			if precommit[0].Value() == proposal.Value() {
	// 				// accusation or proof
	// 			}
	// 		}
	// 	}

	// 	//check all precommits for previous rounds from this sender are nil
	// 	precommits := i.msgStore(height, func(m *message) bool {
	// 		return m.sender() == proposal.Sender() && m.Type() == precommit && m.Round < proposal.Round && m.Value != nilValue
	// 	})

		

	// 	allPrevotes := i.msgStore(height, func(m *message) bool {
	// 		return m.Type() == prevotes && 
	// 	})
		
	// 	if len(precommits) != 0 {
	// 		// construct proof of bad behaviour
	// 		proof := Proof{
	// 			Rule:     PN1,
	// 			Evidence: []message{precommits[0]},
	// 			Message:  msg,
	// 		}
	// 		proofs = append(proofs, &proof)
	// 	}
	// }


	// (Mr′<r,PV) ∧ (Mr′,PC|pi) ∧ (Mr′<r′′<r,PC|pi)∗ ⇒ (Mr,P|pi)
	// [#(V) ≥ 2f+1] ∧ [nil ∨ V v ⊥] ∧ [nil v ⊥] ⇒ [V]

	// PO1 Take 2
	proposalsOld := i.msgStore(height, func(m *message) bool {
		return m.Type() == proposal && m.ValidRound > -1 
	})
	
	for i, proposal := range proposalsOld {
		// Check that in the valid round we see a quorum of prevotes and that
		// there is no precommit at all or a precommit for v or nil.

		validRound := proposal.ValidRound()

		// Is there a precommit for a value other than nil or the proposed value
		// by the current proposer in the valid round? If there is the proposer
		// has proposed a value for which it is not locked on, thus a proof of
		// misbehaviour can be generated.
		precommit := i.msgStore(height, func(m *message) bool {
			return m.Type() == Precommit && m.Round == validRound &&
			m.Sender == proposal.Sender() && m.Value() != nilValue &&
			m.Value() != proposal.Value()
		})
		if len(precommit) > 0 {
			return &Proof{
				Rule:     PO1,
				Evidence: []message{precommit[0]},
				Message:  msg,
			}
		}

		// Is there a precommit for anything other than nil from the proposer
		// between the valid round and the round of the proposal? If there is
		// then that implies the proposer saw 2f+1 prevotes in that round and
		// hence it should have set that round as the valid round.
		precommits := i.msgStore(height, func(m *message) bool {
			return m.Type() == Precommit &&
			m.Round > validRound && m.Round < proposal.Round() &&
			m.Sender == proposal.Sender() &&
			m.Value() != nilValue
		})
		if len(precommits) > 0 {
			return &Proof{
				Rule:     PO1,
				Evidence: []message{precommits[0]},
				Message:  msg,
			}
		}

		// Do we see a quorum of prevotes in the valid round, if not we can
		// raise an accusation, since we cannot be sure that these prevotes
		// don't exist
		prevotes := i.msgStore(height, func(m *message) bool {
			return m.Type() == prevotes && m.Round == validRound
		})
		if len(prevotes) < threshold {
			return &Accusation{
				Rule:     PO1,
				Message:  msg,
			}
		}
	}
	
	// Find all the prevotes
	prevotes := i.msgStore(height, func(m *message) bool{
		return m.Type() == Prevote && m.Value != nilValue
	})

	// iterate over all the prevotes
	// find the proposal which refer to the prevote
	// if a proposal is not found raise an accusation
	// Determine if the proposal is new or old
	// then apply the PVN or PVO rules

	for prevote := range prevotes {
		correspondingProposals := i.msgStore(height, func(m *message) bool {
			m.Type() == Proposal && m.Value == prevote.Value && m.Round() == prevote.Round()
		})

		if len(correspondingProposals) == 0 {
			// raise an accusation 
			 &Accusation{
				Rule:     PV,
				Message:  msg,
			}
		}


		// We need to ensure that we keep all proposals in the message store,
		// so that we have the maximum chance of finding justification for
		// prevotes.
		//
		// P V -1 <-- real - discarded
		// P V 4 <-- fake there never were 2f+1 PV for V in the past
		
		
		for correspondingProposal :=  range correspondingProposals {
			if correspondingProposal.ValidRound == -1 {
				// New Proposal, apply PVN rules
				
				// PVN: (Mr′<r,PC|pi)∧(Mr′<r′′<r,PC|pi)* ∧ (Mr,P|proposer(r)) <--- (Mr,PV|pi)

				// PVN2: [nil ∨ ⊥] ∧ [nil ∨ ⊥] ∧ [V:Valid(V)] <--- [V]: r′= 0,∀r′′< r:Mr′′,PC|pi=nil
				
				// PVN2, If there is a valid proposal V at round r, and pi never
				// ever precommit(locked a value) before, then pi should prevote
				// for V or a nil in case of timeout at this round.

				// PVN3: [V] ∧ [nil ∨ ⊥] ∧ [V:Valid(V)] <--- [V]:∀r′< r′′<r,Mr′′,PC|pi=nil

				// We can check both PVN2 and PVN3 by simply searching for a
				// precommit for a value other than V or nil. This is a proof of
				// misbehaviour. There is no scope to raise an accusation for
				// these rules since the only message in PVN that is not sent by
				// pi is the proposal and you require the proposal before you
				// can even attempt to apply the rule.
				precommits := i.msgStore(height, func(m *message) bool {
					m.Type() == Precommit && m.Value != nilValue &&
					m.Value != prevote.Value && prevote.Sender() == m.Sender() &&
					m.Round() < prevote.Round()
				})

				if len(precommits) > 0 {
					// Proof of misbehaviour
					break
				}

			} else {
				// Old Proposal, apply PVO rules
				// requires 2f+1 PV for V ( but it doesn't exist)	
				// So raise accusation

				// PVO:   (Mr′<r,PC|pi) ∧ (Mr′≤r′′′<r,PV) ∧ (Mr′<r′′<r,PC|pi)* ∧ (Mr,P|proposer(r)) <--- (Mr,PV|pi)

				// PVO1A: [V] ∧ [∗] ∧ [nil v ⊥] ∧ [V] <--- [V]:∀r′<r′′<r,Mr′′,PC|pi=nil

			


				// PVO1A: [V] ∧ [#(nil) ≥ 2f+1] ∧ [nil v ⊥] ∧ [V] <--- [V]:∀r′<r′′<r,Mr′′,PC|pi=nil
				// PVO1A: [V] ∧ [⊥] ∧ [nil v ⊥] ∧ [V] <--- [V]:∀r′<r′′<r,Mr′′,PC|pi=nil

				
				// PVO1A2: [nil] ∧ [#2f+1] ∧ [nil v ⊥] ∧ [V] <--- [V]:∀r′<r′′<r,Mr′′,PC|pi=nil
				// If pi previously precommitted for V and between this precommit and
				// the proposal precommitted for a different value V', then the prevote
				// is considered invalid.

				// PVO1B: [∗] ∧ [∗] ∧ [V:r′′=r−1] ∧ [V] <--- [V] -- not needed as it is a special case of PVO1A

				// PVO2: [V'] ∧ [#(V) ≥ 2f+1] ∧ [nil v ⊥] ∧ [V:validRound(V)=r′′′] <--- [V]:∀r′<r′′<r,Mr′′,PC|pi=nil ∧ ∃r′′′∈[r′,r−1],#(Mr′′′,PV|V) ≥ 2f+ 1

				// Imagine in round 1 we see 2f+1 prevotes for V then in round 2
				// we see an old proposal for V with valid round 1.
				// And we then see a prevote in round 2 for the old value from pi.





				// PVO2: [V']  ∧ [#(V) ≥ 2f+1] ∧ [nil v ⊥] ∧ [V:validRound(V)=r′′′] <--- [V]


				// PVO2: [V] ∧ [#(V) ≥ 2f+1] ∧ [nil v ⊥] ∧ [V:validRound(V)=r′′′] <--- [V]

				
			


				// PVO2: [*]  ∧ [#(V) ≥ 2f+1] ∧ [nil v ⊥] ∧ [V:validRound(V)=r′′′] <--- [V]

				
				// PVO2: [V']  ∧ [#(V) ≥ 2f+1] ∧ [nil v ⊥] ∧ [V:validRound(V)=r′′′] <--- [V]
				// PVO2: [V]   ∧ [#(V) ≥ 2f+1] ∧ [nil v ⊥] ∧ [V:validRound(V)=r′′′] <--- [V]
				// PVO2: [nil] ∧ [#(V) ≥ 2f+1] ∧ [nil v ⊥] ∧ [V:validRound(V)=r′′′] <--- [V]
				// PVO2: [⊥]   ∧ [#(V) ≥ 2f+1] ∧ [nil v ⊥] ∧ [V:validRound(V)=r′′′] <--- [V]

				// If pi previously precommitted for a value V' distinct from
				// the old proposed value V and has not subsequently
				// precommitted for V, we must see a quorum of prevotes for V
				// for the validRound specified in the proposal which must be
				// more recent than the round in which pi precommitted for V'.


			}
		}


	}

	prevoteMap := make(map[common.Hash][]*message)
	for i := range prevotes {
		prevoteMap[prevotes[i].Value()]
	}

	for _, propNew := range proposalsNew {
	}

	for _, propOld := range proposalsOld {
		newPrevotes := i.msgStore(height, func(m *message) bool{
			return m.Type() == Prevote && m.Round == propNew.Round && m.Value() == propNew.Value()
		})

		>>
	}



	if len(proposal) != 1 {
		return nil, nil
	}

	// Check if pi never ever locked a value before.
	precommits := i.msgStore(msg.Height(), func(m *message) bool {
		return m.Type() == precommit && m.Round < msg.Round() && m.Value() != nilValue && m.Sender() == msg.Sender()
	})

	if len(precommits) != 0 {
		return nil, nil
	}

	// If pi never precommit a value before, then it should prevote for nil or V, otherwise generate the proof of
	// violation of PVN2.
	if !(msg.Value() == proposal.Value() || msg.Value() == nilValue) {
		return &Proof{
			Rule:     PVN2,
			Evidence: []message{proposal, precommits},
			Message:  msg,
		}
	}


}

func (i *interceptor) Intercept(msg *message) proofOfMisBehavior {
	// Prerequisite: msg has a valid signature and comes from validator.

	// Validation steps
	//
	// Auto incriminating
	//  - Need to check for proposals if they are coming from the right proposer.
	//  - Is Type valid ? one of (propose, prevote, precommit)
	//  - Check for correctness of old proposals (VR = -1)
	//  - Check that the valid round in the old proposal (VR > -1) is not equal to or greater than the current round.
	//  - Check the validity of the proposal and if it is invalid it is an auto incrimination message
	//

	// proover.Send(interceptor.Intercept(msg))
	//
	// interceptor.Intercept(msg)

	// Saving messages in the store MUST happen before checking for equivocation.

	i.msgStore.Save(msg)

	// if proof := i.checkEquivocation(msg) ; proof != nil {
	// 	return proof
	// }

	if proof := i.checkImmediateFault(msg); proof != nil {
		return proof
	}
	return 
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
			return &Proof{
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
	if msg.Type == prevote && msg.Value() == nilValue {
		proposal := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == propose && m.Round == msg.Round() && notValid(m.Value)
		})

		if len(proposal) != 0 {
			return &Proof{
				Rule:     PVN1,
				Evidence: []message{proposal},
				Message:  msg,
			}
		}
	}

	// PVN2, If there is a valid proposal V at round r, and pi never ever precommit(locked a value) before, then pi should prevote
	// for V or a nil in case of timeout at this round.
	if msg.Type == prevote {
		// Check if we have valid proposal on the round r.
		proposal := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == propose && m.Round == msg.Round() && Valid(m.Value)
		})

		// Check if pi never ever locked a value before.
		precommits := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == precommit && m.Round < msg.Round() && m.Value() != nilValue && m.Sender() == msg.Sender()
		})

		// If pi never precommit a value before, then it should prevote for nil or V, otherwis generate the proof of
		// violation of PVN2.
		if proposal != nil && len(precommits) == 0 && !(msg.Value() == proposal.Value() || msg.Value() == nilValue) {
			return &Proof{
				Rule:     PVN2,
				Evidence: []message{proposal, precommits},
				Message:  msg,
			}
		}
	}

	// PVN3, if V is a valid proposed value, and pi locked it in the previous round, the pi should prevote for V or Nil in case of timeout.
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

		// pi last locked at a distinct value before.
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
			return m.Type() == precommit && (precommits[0].Round < m.Round() || m.Round() < msg.Round()) && m.Value() != proposal.Value() && m.Sender() == msg.Sender()
		})

		// the prevote of pi should be nil or V, otherwise it break the rule
		if proposal != nil && len(precommits) == 1 && len(otherPrecommits) == 0 && !(msg.Value() == proposal.Value() || msg.Value() == nilValue) {
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
		// pi locked at the same value before round r at previous round
		precommits := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == precommit && m.Round < msg.Round() && m.Value() == proposal.Value() && m.Sender() == msg.Sender()
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

	// PVO4. if V is the proposed value at round r and pi did already precommit on V' at the previous round then pi prevotes for nil.
	if msg.Type == prevote {
		// Valid V proposed at round r.
		proposal := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == propose && m.Round == msg.Round() && Valid(m.Value)
		})

		// todo: assume that no equivocation msg on msg store.
		// pi locked at a distinct value V' before round r
		precommits := i.msgStore(msg.Height(), func(m *message) bool {
			return m.Type() == precommit && m.Round < msg.Round() && m.Value() != proposal.Value() && m.Sender() == msg.Sender()
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
