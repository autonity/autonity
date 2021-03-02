package faultdetector

import (
	"github.com/clearmatics/autonity/common"
)

type Rule uint8

const (
	PN Rule = iota
	PO
	PVN
	PVO
	C
)

type Proof struct {
	parentHash common.Hash // use by precompiled contract to get committee from chain db.
	Rule       Rule
	Message    message
	Evidence   []message
}

type Accusation struct {
	Rule    Rule
	Message message
}

type Store interface {
	Save(message message)
	Get(height uint64, query func(m message) bool) []message
}

type interceptor struct {
	msgStore Store
}

const (
	p byte = iota
	pv
	pc
)

var nilValue = common.Hash{}

type message interface {
	Round() uint
	Height() uint
	Sender() common.Address
	Type() byte
	Value() common.Hash // Block hash for a proposal,
	ValidRound() int
	Payload() []byte // raw bytes of message
}

func (i *interceptor) Intercept(msg message) *Proof {
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
	// participants precommiting or prevoting for an invalid proposal or
	// prevotes for a proposal that is not from the current proposer.

	// Saving messages in the store MUST happen before checking for equivocation.

	i.msgStore.Save(msg)

	// if proof := i.checkEquivocation(msg) ; proof != nil {
	// 	return proof
	// }

	if proof := i.checkImmediateFault(msg); proof != nil {
		return proof
	}

	return nil
}

func (i *interceptor) checkImmediateFault(m message) *Proof {
	// Check for auto-incriminating message and equivocation.
	// The Proof struct as it is defined now may not be sufficient to represent proofs for auto-incriminating and
	// equivocation messages.
	return nil
}

func threshold(height uint64) uint {
	// Determine the quorum threshold for the current height
	return 0
}

func (i *interceptor) Process(height uint64) ([]*Proof, []*Accusation) {
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

	// We should be here at time t = timestamp(h+1) + delta

	var proofs []*Proof
	var accusations []*Accusation

	// ------------New Proposal------------
	// PN:  (Mr′<r,P C|pi)∗ <--- (Mr,P|pi)
	// PN1: [nil ∨ ⊥] <--- [V]

	proposalsNew := i.msgStore.Get(height, func(m message) bool {
		return m.Type() == p && m.ValidRound() == -1
	})

	for _, proposal := range proposalsNew {
		//check all precommits for previous rounds from this sender are nil
		precommits := i.msgStore.Get(height, func(m message) bool {
			return m.Sender() == proposal.Sender() && m.Type() == pc && m.Round() < proposal.Round() && m.Value() != nilValue
		})
		if len(precommits) != 0 {
			proof := &Proof{
				Rule:     PN,
				Evidence: precommits,
				Message:  proposal,
			}
			proofs = append(proofs, proof)
		}
	}

	// ------------Old Proposal------------
	// PO: (Mr′<r,PV) ∧ (Mr′,PC|pi) ∧ (Mr′<r′′<r,P C|pi)∗ <--- (Mr,P|pi)
	// PO1: [#(Mr′,PV|V) ≥ 2f+ 1] ∧ [nil ∨ V ∨ ⊥] ∧ [nil ∨ ⊥] <--- [V]

	proposalsOld := i.msgStore.Get(height, func(m message) bool {
		return m.Type() == p && m.ValidRound() > -1
	})

	for _, proposal := range proposalsOld {
		// Check that in the valid round we see a quorum of prevotes and that
		// there is no precommit at all or a precommit for v or nil.

		validRound := uint(proposal.ValidRound())

		// Is there a precommit for a value other than nil or the proposed value
		// by the current proposer in the valid round? If there is the proposer
		// has proposed a value for which it is not locked on, thus a proof of
		// misbehaviour can be generated.
		precommits := i.msgStore.Get(height, func(m message) bool {
			return m.Type() == pc && m.Round() == validRound &&
				m.Sender() == proposal.Sender() && m.Value() != nilValue &&
				m.Value() != proposal.Value()
		})
		if len(precommits) > 0 {
			proof := &Proof{
				Rule:     PO,
				Evidence: precommits,
				Message:  proposal,
			}
			proofs = append(proofs, proof)
		}

		// Is there a precommit for anything other than nil from the proposer
		// between the valid round and the round of the proposal? If there is
		// then that implies the proposer saw 2f+1 prevotes in that round and
		// hence it should have set that round as the valid round.
		precommits = i.msgStore.Get(height, func(m message) bool {
			return m.Type() == pc &&
				m.Round() > validRound && m.Round() < proposal.Round() &&
				m.Sender() == proposal.Sender() &&
				m.Value() != nilValue
		})
		if len(precommits) > 0 {
			proof := &Proof{
				Rule:     PO,
				Evidence: precommits,
				Message:  proposal,
			}
			proofs = append(proofs, proof)
		}

		// Do we see a quorum of prevotes in the valid round, if not we can
		// raise an accusation, since we cannot be sure that these prevotes
		// don't exist
		prevotes := i.msgStore.Get(height, func(m message) bool {
			return m.Type() == pv && m.Round() == validRound
		})
		if len(prevotes) < int(threshold(height)) {
			accusation := &Accusation{
				Rule:    PO,
				Message: proposal,
			}
			accusations = append(accusations, accusation)
		}
	}

	// ------------New and Old Prevotes------------

	prevotes := i.msgStore.Get(height, func(m message) bool {
		return m.Type() == pv && m.Value() != nilValue
	})

	for _, prevote := range prevotes {
		correspondingProposals := i.msgStore.Get(height, func(m message) bool {
			return m.Type() == p && m.Value() == prevote.Value() && m.Round() == prevote.Round()
		})

		if len(correspondingProposals) == 0 {
			accusation := &Accusation{
				Rule: PVN, //This could be PVO as well, however, we can't decide since there are no corresponding
				// proposal
				Message: prevote,
			}
			accusations = append(accusations, accusation)
		}

		// We need to ensure that we keep all proposals in the message store, so that we have the maximum chance of
		// finding justification for prevotes. This is to account for equivocation where the proposer send 2 proposals
		// with the same value but different valid rounds to different nodes. We can't penalise the sender of prevote
		// since we can't tell which proposal they received. We just want to find a set of message which fit the rule.

		for _, correspondingProposal := range correspondingProposals {
			if correspondingProposal.ValidRound() == -1 {
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
				precommits := i.msgStore.Get(height, func(m message) bool {
					return m.Type() == pc && m.Value() != nilValue &&
						m.Value() != prevote.Value() && prevote.Sender() == m.Sender() &&
						m.Round() < prevote.Round()
				})

				if len(precommits) > 0 {
					proof := &Proof{
						Rule:     PVN,
						Evidence: precommits,
						Message:  prevote,
					}
					proofs = append(proofs, proof)
					break
				}

			} else {
				// PVO:   (Mr′<r,PC|pi) ∧ (Mr′≤r′′′<r,PV) ∧ (Mr′<r′′<r,PC|pi)* ∧ (Mr,P|proposer(r)) <--- (Mr,PV|pi)

				// PVO1A: [V] ∧ [∗] ∧ [nil v ⊥] ∧ [V] <--- [V]:∀r′<r′′<r,Mr′′,PC|pi=nil <-- broken we need to see the prevotes for valid round

				// PVO2: [*] ∧ [#(V) ≥ 2f+1] ∧ [nil v ⊥] ∧ [V:validRound(V)=r′′′] <--- [V]:∀r′<r′′<r,Mr′′,PC|pi=nil ∧ ∃r′′′∈[r′,r−1],#(Mr′′′,PV|V) ≥ 2f+ 1

				// If pi previously precommitted for V and between this precommit and
				// the proposal precommitted for a different value V', then the prevote
				// is considered invalid.

				precommits := i.msgStore.Get(height, func(m message) bool {
					return m.Type() == pc && prevote.Sender() == m.Sender() &&
						m.Round() < prevote.Round() && m.Value() != nilValue
				})
				//check most recent precommit if == V -> pass else --> fail

				// 2f+1 PV(V) round 2

				// round 4 p_i receiveds 2f+1 PV(V') Sends PC(V') and it sets its locked value and locked round=4

				// round 5 proposer proposes P(V, VR=2), so this would mean that p_i prevote nil even though there are 2f+1 prevotes for V in round 2

				// Aneeque's initials thoughts on PVO
				if len(precommits) > 0 {
					// PVO1a

					// sort according to round
					//sort.Sort(precommits)

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

			// ------------Precommits------------
			// C: [Mr,P|proposer(r)] ∧ [Mr,PV] <--- [Mr,PC|pi]
			// C1: [V:Valid(V)] ∧ [#(V) ≥ 2f+ 1] <--- [V]

			precommits := i.msgStore.Get(height, func(m message) bool {
				return m.Type() == pc && m.Value() != nilValue
			})

			for _, precommit := range precommits {
				proposals := i.msgStore.Get(height, func(m message) bool {
					return m.Type() == p && m.Value() == precommit.Value() && m.Round() == precommit.Round()
				})

				if len(proposals) == 0 {
					accusation := &Accusation{
						Rule:    C,
						Message: precommit,
					}
					accusations = append(accusations, accusation)
					continue
				}

				prevotesForNotV := i.msgStore.Get(height, func(m message) bool {
					return m.Type() == pv && m.Value() != precommit.Value() && m.Round() == precommit.Round()
				})
				prevotesForV := i.msgStore.Get(height, func(m message) bool {
					return m.Type() == pv && m.Value() == precommit.Value() && m.Round() == precommit.Round()
				})

				if len(prevotesForNotV) >= int(threshold(height)) {
					// In this case there cannot be enough remaining prevotes
					// to justify a precommit for V.
					proof := &Proof{
						Rule:     C,
						Evidence: prevotesForNotV,
						Message:  precommit,
					}
					proofs = append(proofs, proof)

				} else if len(prevotesForV) < int(threshold(height)) {
					// In this case we simply don't see enough prevotes to
					// justify the precommit.
					accusation := &Accusation{
						Rule:    C,
						Message: precommit,
					}
					accusations = append(accusations, accusation)
				}
			}
		}
	}
	return proofs, accusations
}
