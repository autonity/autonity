package faultdetector

import (
	"fmt"
	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/core"
)

var nilValue = common.Hash{}

func powerOfVotes(votes []core.Message) uint64 {
	power := uint64(0)
	for i := 0; i < len(votes); i++ {
		if votes[i].Type() != msgPrevote || votes[i].Type() != msgPrecommit {
			continue
		}
		power += votes[i].GetPower()
	}
	return power
}

// run rule engine over latest msg store, if the return proofs is not empty, then rise challenge.
func (fd *FaultDetector) runRuleEngine(height uint64, quorum uint64) {
	proofs := fd.runRules(height, quorum)
	if len(proofs) > 0 {
		var onChainProofs []autonity.OnChainProof
		for i := 0; i < len(proofs); i++ {
			p, err := fd.generateOnChainProof(&proofs[i].Message, proofs[i].Evidence, proofs[i].Rule, proofs[i].Type)
			if err != nil {
				fd.logger.Warn("convert proof to on-chain proof", "faultdetector", err)
				continue
			}
			onChainProofs = append(onChainProofs, p)
		}
		fd.sendProofs(true, onChainProofs)
	}
}

// getInnocentProof called by client who is on challenge to get proof of innocent from msg store.
func (fd *FaultDetector) getInnocentProof(c *Proof) (autonity.OnChainProof, error) {
	var proof autonity.OnChainProof
	// rule engine have below provable accusation for the time being:
	switch c.Rule {
	case PO:
		return fd.GetInnocentProofOfPO(c)
	case PVN:
		return fd.GetInnocentProofOfPVN(c)
	case C:
		return fd.GetInnocentProofOfC(c)
	case C1:
		return fd.GetInnocentProofOfC1(c)
	default:
		return proof, fmt.Errorf("not provable rule")
	}
}

// get proof of innocent of PO from msg store.
func (fd *FaultDetector) GetInnocentProofOfPO(c *Proof) (autonity.OnChainProof, error) {
	// PO: node propose an old value with an validRound, innocent proof of it should be:
	// there are quorum num of prevote for that value at the validRound.
	var proof autonity.OnChainProof
	proposal := c.Message
	height := proposal.H()
	validRound := proposal.ValidRound()
	quorum := fd.quorum(height - 1)

	prevotes := fd.msgStore.Get(height, func(m *core.Message) bool {
		return m.Type() == msgPrevote && m.R() == validRound && m.Value() == proposal.Value()
	})

	if powerOfVotes(prevotes) < quorum {
		// cannot proof its innocent for PO, the on-chain contract will fine it latter once the
		// time window for proof ends.
		return proof, fmt.Errorf("node is malicious")
	}

	p, err := fd.generateOnChainProof(&proposal, prevotes, c.Rule, Innocence)
	if err != nil {
		return p, err
	}

	return p, nil
}

// get proof of innocent of PVN from msg store.
func (fd *FaultDetector) GetInnocentProofOfPVN(c *Proof) (autonity.OnChainProof, error) {
	// get innocent proofs for PVN, for a prevote that vote for a new value,
	// then there must be a proposal for this new value.
	var proof autonity.OnChainProof
	prevote := c.Message
	height := prevote.H()

	correspondingProposals := fd.msgStore.Get(height, func(m *core.Message) bool {
		return m.Type() == msgProposal && m.Value() == prevote.Value() && m.R() == prevote.R()
	})

	if len(correspondingProposals) == 0 {
		// cannot proof its innocent for PVN, the on-chain contract will fine it latter once the
		// time window for proof ends.
		return proof, fmt.Errorf("node is malicious")
	}

	p, err := fd.generateOnChainProof(&prevote, correspondingProposals, c.Rule, Innocence)
	if err != nil {
		return p, nil
	}

	return p, nil
}

// get proof of innocent of C from msg store.
func (fd *FaultDetector) GetInnocentProofOfC(c *Proof) (autonity.OnChainProof, error) {
	var proof autonity.OnChainProof
	preCommit := c.Message
	height := preCommit.H()

	proposals := fd.msgStore.Get(height, func(m *core.Message) bool {
		return m.Type() == msgProposal && m.Value() == preCommit.Value() && m.R() == preCommit.R()
	})

	if len(proposals) == 0 {
		// cannot proof its innocent for PVN, the on-chain contract will fine it latter once the
		// time window for proof ends.
		return proof, fmt.Errorf("node is malicious")
	}
	p, err := fd.generateOnChainProof(&preCommit, proposals, c.Rule, Innocence)
	if err != nil {
		return p, err
	}
	return p, nil
}

// get proof of innocent of C1 from msg store.
func (fd *FaultDetector) GetInnocentProofOfC1(c *Proof) (autonity.OnChainProof, error) {
	var proof autonity.OnChainProof
	preCommit := c.Message
	height := preCommit.H()
	quorum := fd.quorum(height - 1)

	prevotesForV := fd.msgStore.Get(height, func(m *core.Message) bool {
		return m.Type() == msgPrevote && m.Value() == preCommit.Value() && m.R() == preCommit.R()
	})

	if powerOfVotes(prevotesForV) < quorum {
		// cannot proof its innocent for PO for now, the on-chain contract will fine it latter once the
		// time window for proof ends.
		return proof, fmt.Errorf("node might be malicious")
	}

	p, err := fd.generateOnChainProof(&preCommit, prevotesForV, c.Rule, Innocence)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (fd *FaultDetector) runRules(height uint64, quorum uint64) (proofs []Proof) {
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

	// ------------New Proposal------------
	// PN:  (Mr′<r,P C|pi)∗ <--- (Mr,P|pi)
	// PN1: [nil ∨ ⊥] <--- [V]

	proposalsNew := fd.msgStore.Get(height, func(m *core.Message) bool {
		return m.Type() == msgProposal && m.ValidRound() == -1
	})

	for _, proposal := range proposalsNew {
		//check all precommits for previous rounds from this sender are nil
		precommits := fd.msgStore.Get(height, func(m *core.Message) bool {
			return m.Sender() == proposal.Sender() && m.Type() == msgPrecommit && m.R() < proposal.R() && m.Value() != nilValue // nolint: scopelint
		})
		if len(precommits) != 0 {
			proof := Proof{
				Type:     Misbehaviour,
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

	proposalsOld := fd.msgStore.Get(height, func(m *core.Message) bool {
		return m.Type() == msgProposal && m.ValidRound() > -1
	})

	for _, proposal := range proposalsOld {
		// Check that in the valid round we see a quorum of prevotes and that
		// there is no precommit at all or a precommit for v or nil.

		validRound := proposal.ValidRound()

		// Is there a precommit for a value other than nil or the proposed value
		// by the current proposer in the valid round? If there is the proposer
		// has proposed a value for which it is not locked on, thus a proof of
		// misbehaviour can be generated.
		precommits := fd.msgStore.Get(height, func(m *core.Message) bool {
			return m.Type() == msgPrecommit && m.R() == validRound &&
				m.Sender() == proposal.Sender() && m.Value() != nilValue && m.Value() != proposal.Value() // nolint: scopelint
		})
		if len(precommits) > 0 {
			proof := Proof{
				Type:     Misbehaviour,
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
		precommits = fd.msgStore.Get(height, func(m *core.Message) bool {
			return m.Type() == msgPrecommit &&
				m.R() > validRound && m.R() < proposal.R() && m.Sender() == proposal.Sender() && m.Value() != nilValue // nolint: scopelint
		})
		if len(precommits) > 0 {
			proof := Proof{
				Type:     Misbehaviour,
				Rule:     PO,
				Evidence: precommits,
				Message:  proposal,
			}
			proofs = append(proofs, proof)
		}

		// Do we see a quorum of prevotes in the valid round, if not we can
		// raise an accusation, since we cannot be sure that these prevotes
		// don't exist
		prevotes := fd.msgStore.Get(height, func(m *core.Message) bool {
			return m.Type() == msgPrevote && m.R() == validRound
		})

		if powerOfVotes(prevotes) < quorum {
			accusation := Proof{
				Type:    Accusation,
				Rule:    PO,
				Message: proposal,
			}
			proofs = append(proofs, accusation)
		}
	}

	// ------------New and Old Prevotes------------

	prevotes := fd.msgStore.Get(height, func(m *core.Message) bool {
		return m.Type() == msgPrevote && m.Value() != nilValue
	})

	for _, prevote := range prevotes {
		correspondingProposals := fd.msgStore.Get(height, func(m *core.Message) bool {
			return m.Type() == msgProposal && m.Value() == prevote.Value() && m.R() == prevote.R() // nolint: scopelint
		})

		if len(correspondingProposals) == 0 {
			accusation := Proof{
				Type: Accusation,
				Rule: PVN, //This could be PVO as well, however, we can't decide since there are no corresponding
				// proposal
				Message: prevote,
			}
			proofs = append(proofs, accusation)
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
				precommits := fd.msgStore.Get(height, func(m *core.Message) bool {
					return m.Type() == msgPrecommit && m.Value() != nilValue &&
						m.Value() != prevote.Value() && prevote.Sender() == m.Sender() && m.R() < prevote.R() // nolint: scopelint
				})

				if len(precommits) > 0 {
					proof := Proof{
						Type: Misbehaviour,
						Rule: PVN,
						// add corresponding proposal at last slot, as the part of evidence to be validated at precompiled contract.
						Evidence: append(precommits, correspondingProposal),
						Message:  prevote,
					}
					proofs = append(proofs, proof)
					break
				}

			} /*else {
				todo: missing PVO rules from D3
				// PVO:   (Mr′<r,PC|pi) ∧ (Mr′≤r′′′<r,PV) ∧ (Mr′<r′′<r,PC|pi)* ∧ (Mr,P|proposer(r)) <--- (Mr,PV|pi)

				// PVO1A: [V] ∧ [∗] ∧ [nil v ⊥] ∧ [V] <--- [V]:∀r′<r′′<r,Mr′′,PC|pi=nil <-- broken we need to see the prevotes for valid round

				// PVO2: [*] ∧ [#(V) ≥ 2f+1] ∧ [nil v ⊥] ∧ [V:validRound(V)=r′′′] <--- [V]:∀r′<r′′<r,Mr′′,PC|pi=nil ∧ ∃r′′′∈[r′,r−1],#(Mr′′′,PV|V) ≥ 2f+ 1

				// If pi previously precommitted for V and between this precommit and
				// the proposal precommitted for a different value V', then the prevote
				// is considered invalid.

				precommits := fd.msgStore.Get(height, func(m *core.Message) bool {
					return m.Type() == msgPrecommit && prevote.Sender() == m.Sender() &&
						m.R() < prevote.R() && m.Value() != nilValue
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
			} */
		}
	}

	// ------------Precommits------------
	// C: [Mr,P|proposer(r)] ∧ [Mr,PV] <--- [Mr,PC|pi]
	// C1: [V:Valid(V)] ∧ [#(V) ≥ 2f+ 1] <--- [V]

	precommits := fd.msgStore.Get(height, func(m *core.Message) bool {
		return m.Type() == msgPrecommit && m.Value() != nilValue
	})

	for _, precommit := range precommits {
		proposals := fd.msgStore.Get(height, func(m *core.Message) bool {
			return m.Type() == msgProposal && m.Value() == precommit.Value() && m.R() == precommit.R() // nolint: scopelint
		})

		if len(proposals) == 0 {
			accusation := Proof{
				Type:    Accusation,
				Rule:    C,
				Message: precommit,
			}
			proofs = append(proofs, accusation)
			continue
		}

		prevotesForNotV := fd.msgStore.Get(height, func(m *core.Message) bool {
			return m.Type() == msgPrevote && m.Value() != precommit.Value() && m.R() == precommit.R() // nolint: scopelint
		})
		prevotesForV := fd.msgStore.Get(height, func(m *core.Message) bool {
			return m.Type() == msgPrevote && m.Value() == precommit.Value() && m.R() == precommit.R() // nolint: scopelint
		})

		if powerOfVotes(prevotesForNotV) >= quorum {
			// In this case there cannot be enough remaining prevotes
			// to justify a precommit for V.
			proof := Proof{
				Type:     Misbehaviour,
				Rule:     C,
				Evidence: prevotesForNotV,
				Message:  precommit,
			}
			proofs = append(proofs, proof)

		} else if powerOfVotes(prevotesForV) < quorum {
			// In this case we simply don't see enough prevotes to
			// justify the precommit.
			accusation := Proof{
				Type:    Accusation,
				Rule:    C1,
				Message: precommit,
			}
			proofs = append(proofs, accusation)
		}
	}

	return proofs
}

func errorToRule(err error) (Rule, error) {
	rule := UnknownRule
	switch err {
	case errEquivocation:
		rule = Equivocation
	case errProposer:
		rule = InvalidProposer
	case errProposal:
		rule = InvalidProposal
	case errGarbageMsg:
		rule = GarbageMessage
	default:
		return rule, fmt.Errorf("errors of not provable")
	}

	return rule, nil
}
