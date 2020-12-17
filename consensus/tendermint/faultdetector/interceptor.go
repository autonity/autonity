package faultdetector

import (
	"github.com/clearmatics/autonity/common"
)

type Rule uint8

const (
	PN1 Rule = iota
	PO1
	PVN2
)

type Proof struct {
	Rule     Rule
	Evidence []message
	Message  message
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

func deferProcessing(msg *message) {
	timer.executeAfter(msg, delta, func() {

	})
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
	// We currently don't include auto-incriminating messages in the message
	// store for simplicity, but there are some cases in which we would be able
	// to prove bad behaviour by including auto-incriminating messages, such as
	// participants precommiting or prevoting for an invalid proposal.

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

func (*interceptor) Process(height uint64) ([]*proofOfMisBehavior, []*Accusation) {
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
				Rule:    PO1,
				Message: msg,
			}
		}
	}

	// Find all the prevotes
	prevotes := i.msgStore(height, func(m *message) bool {
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
				Rule:    PV,
				Message: msg,
			}
		}

		// We need to ensure that we keep all proposals in the message store,
		// so that we have the maximum chance of finding justification for
		// prevotes.
		//
		// P V -1 <-- real - discarded
		// P V 4 <-- fake there never were 2f+1 PV for V in the past

		for correspondingProposal := range correspondingProposals {
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
				// PVO:   (Mr′<r,PC|pi) ∧ (Mr′≤r′′′<r,PV) ∧ (Mr′<r′′<r,PC|pi)* ∧ (Mr,P|proposer(r)) <--- (Mr,PV|pi)

				// PVO1A: [V] ∧ [∗] ∧ [nil v ⊥] ∧ [V] <--- [V]:∀r′<r′′<r,Mr′′,PC|pi=nil <-- broken we need to see the prevotes for valid round

				// PVO2: [*] ∧ [#(V) ≥ 2f+1] ∧ [nil v ⊥] ∧ [V:validRound(V)=r′′′] <--- [V]:∀r′<r′′<r,Mr′′,PC|pi=nil ∧ ∃r′′′∈[r′,r−1],#(Mr′′′,PV|V) ≥ 2f+ 1

				// If pi previously precommitted for V and between this precommit and
				// the proposal precommitted for a different value V', then the prevote
				// is considered invalid.

				precommits := i.msgStore(height, func(m *message) bool {
					m.Type() == Precommit && prevote.Sender() == m.Sender() &&
						m.Round() < prevote.Round() && m.Value != nilValue
				})
				//check most recent precommit if == V -> pass else --> fail

				// 2f+1 PV(V) round 2

				// round 4 p_i receiveds 2f+1 PV(V') Sends PC(V') and it sets its locked value and locked round=4

				// round 5 proposer proposes P(V, VR=2), so this would mean that p_i prevote nil even though there are 2f+1 prevotes for V in round 2

				// Aneeque's initials thoughts on PVO
				if len(precommits) > 0 {
					// PVO1a

					// sort according to round
					Sort(precommits)

					// Proof of misbehaviour:

					// Get the lastest precommit
					// Check the precommit value
					// if it precommit.Value() != prevote.Value
					// 		check all round from precommit to current round for 2f+1 prevotes
					// 		if even a single round doesn't have 2f+1 prevotes, raise an accusation
					//		else we have proof of misbehaviour if non of the 2f+1 prevotes are for precommit.Value()

					// if it precommit.Value() == prevote.Value
					// 		Check that if we 2f+1 prevotes for all rounds since precommit.Round() till current round,
					//      if yes, than non of them can be for value other than prevote.Value, otherwise we have proof of misbehaviour
					// 		if there are gaps then the condition passes

				} else {
					// PVO2

					// We don't have a precommit from the p_i
					// check that in valid round we have 2f+1 prevotes for V rule passes, otherwise raise an accustion
				}

				// PVO1B: [∗] ∧ [∗] ∧ [V:r′′=r−1] ∧ [V] <--- [V] -- not needed as it is a special case of PVO1A

				// PVO2: [*] ∧ [#(V) ≥ 2f+1] ∧ [nil v ⊥] ∧ [V:validRound(V)=r′′′] <--- [V]:∀r′<r′′<r,Mr′′,PC|pi=nil ∧ ∃r′′′∈[r′,r−1],#(Mr′′′,PV|V) ≥ 2f+ 1
				// If we can see an old proposal for V with valid round vr and
				// 2f+1 prevotes for the V in round vr, then pi could have also
				// seen them and hence be able to prevote for the old proposal.

			}

			// C: [Mr,P|proposer(r)] ∧ [Mr,PV] <--- [Mr,PC|pi]

			// C1: [V:Valid(V)] ∧ [#(V) ≥ 2f+ 1] <--- [V]

			precommits := i.msgStore(height, func(m *message) bool {
				m.Type() == Precommit && m.Value != nilValue
			})

			for _, precommit := range precommits {
				proposals := i.msgStore(height, func(m *message) bool {
					m.Type() == Proposal && m.Value == precommit.Value && m.Round() == precommit.Round()
				})

				if len(proposals) == 0 {
					// raise an accusation
					&Accusation{
						Rule:    C,
						Message: msg,
					}
					continue
				}

				prevotesForNotV := i.msgStore(height, func(m *message) bool {
					m.Type() == Prevote && m.Value != precommit.Value() && m.Round() == precommit.Round()
				})
				prevotesForV := i.msgStore(height, func(m *message) bool {
					m.Type() == Prevote && m.Value == precommit.Value() && m.Round() == precommit.Round()
				})
				if len(prevotesForNotV) >= threshold {
					// proof of misbehaviour
					&Proof{
						Rule:     C,
						Evidence: []message{prevotesForNotV},
						Message:  msg,
					}

				} else if len(prevotesForV) < threshold {
					//raise an accusation
					&Accusation{
						Rule:    C,
						Message: msg,
					}

				}
			}
		}
	}
}
