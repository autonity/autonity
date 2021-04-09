package faultdetector

import (
	"fmt"
	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/bft"
	"github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRuleEngine(t *testing.T) {
	height := uint64(100)
	lastHeight := height - 1
	committee, keys := generateCommittee(5)
	proposer := committee[0].Address
	proposerKey := keys[proposer]
	round := int64(3)
	validRound := int64(1)
	totalPower := uint64(len(committee))
	noneNilValue := common.Hash{0x1}

	t.Run("Test de-Equivocated msg", func(t *testing.T) {
		inputMsgs := make([]*core.Message, 2)
		proposal := newProposalMessage(height, round, -1, proposerKey, committee, nil)
		inputMsgs[0] = proposal
		inputMsgs[1] = proposal
		assert.Equal(t, 1, len(deEquivocatedMsgs(inputMsgs)))
	})

	t.Run("getInnocentProof with unprovable rule id", func(t *testing.T) {
		fd := NewFaultDetector(nil, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}))
		var input = Proof{
			Rule: autonity.PVO,
		}

		_, err := fd.getInnocentProof(&input)
		assert.NotNil(t, err)
	})

	t.Run("getInnocentProofOfPO have quorum preVotes", func(t *testing.T) {

		// PO: node propose an old value with an validRound, innocent proof of it should be:
		// there were quorum num of preVote for that value at the validRound.

		fd := NewFaultDetector(nil, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}))
		fd.blocksTotalVotingPower[lastHeight] = totalPower
		// simulate a proposal message with an old value and a valid round.
		proposal := newProposalMessage(height, round, validRound, proposerKey, committee, nil)
		_, err := fd.msgStore.Save(proposal)
		assert.NoError(t, err)

		// simulate at least quorum num of preVotes for a value at a validRound.
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, validRound, msgPrevote, keys[committee[i].Address], proposal.Value(), committee)
			_, err = fd.msgStore.Save(preVote)
			assert.NoError(t, err)
		}

		var accusation = Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.PO,
			Message: proposal,
		}

		proof, err := fd.getInnocentProofOfPO(&accusation)
		assert.NoError(t, err)
		assert.Equal(t, autonity.Innocence, proof.Type)
		assert.Equal(t, proposer, proof.Sender)
		assert.Equal(t, types.RLPHash(proposal), proof.Msghash)
	})

	t.Run("getInnocentProofOfPO no quorum preVotes", func(t *testing.T) {

		// PO: node propose an old value with an validRound, innocent proof of it should be:
		// there were quorum num of preVote for that value at the validRound.

		fd := NewFaultDetector(nil, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}))
		fd.blocksTotalVotingPower[lastHeight] = totalPower
		// simulate a proposal message with an old value and a valid round.
		proposal := newProposalMessage(height, round, validRound, proposerKey, committee, nil)
		_, err := fd.msgStore.Save(proposal)
		assert.NoError(t, err)

		// simulate less than quorum num of preVotes for a value at a validRound.
		preVote := newVoteMsg(height, validRound, msgPrevote, proposerKey, proposal.Value(), committee)
		_, err = fd.msgStore.Save(preVote)
		assert.NoError(t, err)

		var accusation = Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.PO,
			Message: proposal,
		}

		_, err = fd.getInnocentProofOfPO(&accusation)
		assert.Equal(t, errNoEvidenceForPO, err)
	})

	t.Run("getInnocentProofOfPVN have corresponding proposal", func(t *testing.T) {

		// PVN: node prevote for a none nil value, then there must be a corresponding proposal.

		fd := NewFaultDetector(nil, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}))
		// simulate a proposal message with an old value and a valid round.
		proposal := newProposalMessage(height, round, -1, proposerKey, committee, nil)
		_, err := fd.msgStore.Save(proposal)
		assert.NoError(t, err)

		preVote := newVoteMsg(height, round, msgPrevote, proposerKey, proposal.Value(), committee)
		_, err = fd.msgStore.Save(preVote)
		assert.NoError(t, err)

		var accusation = Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.PVN,
			Message: preVote,
		}

		proof, err := fd.getInnocentProofOfPVN(&accusation)
		assert.NoError(t, err)
		assert.Equal(t, autonity.Innocence, proof.Type)
		assert.Equal(t, proposer, proof.Sender)
		assert.Equal(t, types.RLPHash(preVote), proof.Msghash)
	})

	t.Run("getInnocentProofOfPVN have no corresponding proposal", func(t *testing.T) {

		// PVN: node prevote for a none nil value, then there must be a corresponding proposal.
		fd := NewFaultDetector(nil, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}))

		preVote := newVoteMsg(height, round, msgPrevote, proposerKey, noneNilValue, committee)
		_, err := fd.msgStore.Save(preVote)
		assert.NoError(t, err)

		var accusation = Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.PVN,
			Message: preVote,
		}

		_, err = fd.getInnocentProofOfPVN(&accusation)
		assert.Equal(t, errNoEvidenceForPVN, err)
	})

	t.Run("getInnocentProofOfC have corresponding proposal", func(t *testing.T) {

		// C: node preCommit at a none nil value, there must be a corresponding proposal.

		fd := NewFaultDetector(nil, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}))
		// simulate a proposal message with an old value and a valid round.
		proposal := newProposalMessage(height, round, -1, proposerKey, committee, nil)
		_, err := fd.msgStore.Save(proposal)
		assert.NoError(t, err)

		preCommit := newVoteMsg(height, round, msgPrecommit, proposerKey, proposal.Value(), committee)
		_, err = fd.msgStore.Save(preCommit)
		assert.NoError(t, err)

		var accusation = Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.C,
			Message: preCommit,
		}

		proof, err := fd.getInnocentProofOfC(&accusation)
		assert.NoError(t, err)
		assert.Equal(t, autonity.Innocence, proof.Type)
		assert.Equal(t, proposer, proof.Sender)
		assert.Equal(t, types.RLPHash(preCommit), proof.Msghash)
	})

	t.Run("getInnocentProofOfC have no corresponding proposal", func(t *testing.T) {

		// C: node preCommit at a none nil value, there must be a corresponding proposal.

		fd := NewFaultDetector(nil, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}))

		preCommit := newVoteMsg(height, round, msgPrecommit, proposerKey, noneNilValue, committee)
		_, err := fd.msgStore.Save(preCommit)
		assert.NoError(t, err)

		var accusation = Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.C,
			Message: preCommit,
		}

		_, err = fd.getInnocentProofOfC(&accusation)
		assert.Equal(t, errNoEvidenceForC, err)
	})

	t.Run("getInnocentProofOfC1 have quorum preVotes", func(t *testing.T) {

		// C1: node preCommit at a none nil value, there must be quorum corresponding preVotes with same value and round.

		fd := NewFaultDetector(nil, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}))
		fd.blocksTotalVotingPower[lastHeight] = totalPower

		// simulate at least quorum num of preVotes for a value at a validRound.
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, round, msgPrevote, keys[committee[i].Address], noneNilValue, committee)
			_, err := fd.msgStore.Save(preVote)
			assert.NoError(t, err)
		}

		preCommit := newVoteMsg(height, round, msgPrecommit, proposerKey, noneNilValue, committee)
		_, err := fd.msgStore.Save(preCommit)
		assert.NoError(t, err)

		var accusation = Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.C1,
			Message: preCommit,
		}

		proof, err := fd.getInnocentProofOfC1(&accusation)
		assert.NoError(t, err)
		assert.Equal(t, autonity.Innocence, proof.Type)
		assert.Equal(t, proposer, proof.Sender)
		assert.Equal(t, types.RLPHash(preCommit), proof.Msghash)
	})

	t.Run("getInnocentProofOfC1 have no quorum preVotes", func(t *testing.T) {

		// C1: node preCommit at a none nil value, there must be quorum corresponding preVotes with same value and round.

		fd := NewFaultDetector(nil, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}))
		fd.blocksTotalVotingPower[lastHeight] = totalPower

		preCommit := newVoteMsg(height, round, msgPrecommit, proposerKey, noneNilValue, committee)
		_, err := fd.msgStore.Save(preCommit)
		assert.NoError(t, err)

		var accusation = Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.C1,
			Message: preCommit,
		}

		_, err = fd.getInnocentProofOfC1(&accusation)
		assert.Equal(t, errNoEvidenceForC1, err)
	})

	t.Run("Test error to rule mapping", func(t *testing.T) {
		rule, err := errorToRule(errEquivocation)
		assert.NoError(t, err)
		assert.Equal(t, autonity.Equivocation, rule)

		rule, err = errorToRule(errProposer)
		assert.NoError(t, err)
		assert.Equal(t, autonity.InvalidProposer, rule)

		rule, err = errorToRule(errProposal)
		assert.NoError(t, err)
		assert.Equal(t, autonity.InvalidProposal, rule)

		rule, err = errorToRule(errGarbageMsg)
		assert.NoError(t, err)
		assert.Equal(t, autonity.GarbageMessage, rule)

		_, err = errorToRule(fmt.Errorf("unknown err"))
		assert.Error(t, err)
	})

	t.Run("Test calculate power of votes", func(t *testing.T) {
		var preVotes []*core.Message
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, round, msgPrevote, keys[committee[i].Address], noneNilValue, committee)
			preVotes = append(preVotes, preVote)
		}

		// let duplicated msg happens, the counting should skip duplicated ones.
		preVotes = append(preVotes, preVotes...)
		assert.Equal(t, uint64(len(committee)), powerOfVotes(preVotes))
	})

	t.Run("RunRule address the misbehaviour of PN rule", func(t *testing.T) {
		// ------------New Proposal------------
		// PN:  (Mr′<r,P C|pi)∗ <--- (Mr,P|pi)
		// PN1: [nil ∨ ⊥] <--- [V]
		// If one send a maliciousProposal for a new V, then all preCommits for previous rounds from this sender are nil.

		fd := NewFaultDetector(nil, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}))
		quorum := bft.Quorum(totalPower)

		// simulate there was a maliciousProposal at init round 0, and save to msg store.
		initProposal := newProposalMessage(height, 0, -1, keys[committee[1].Address], committee, nil)
		_, err := fd.msgStore.Save(initProposal)
		assert.NoError(t, err)
		// simulate there were quorum preVotes for initProposal at init round 0, and save them.
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, 0, msgPrevote, keys[committee[i].Address], initProposal.Value(), committee)
			_, err = fd.msgStore.Save(preVote)
			assert.NoError(t, err)
		}

		// Node preCommit for init Proposal at init round 0 since there were quorum preVotes for it, and save it.
		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, initProposal.Value(), committee)
		_, err = fd.msgStore.Save(preCommit)
		assert.NoError(t, err)

		// While Node propose a new malicious Proposal at new round with VR as -1 which is malicious, should be addressed by rule PN.
		maliciousProposal := newProposalMessage(height, round, -1, proposerKey, committee, nil)
		_, err = fd.msgStore.Save(maliciousProposal)
		assert.NoError(t, err)

		// Run rule engine over msg store on height.
		onChainProofs := fd.runRulesOverHeight(height, quorum)
		assert.Equal(t, 1, len(onChainProofs))
		assert.Equal(t, autonity.Misbehaviour, onChainProofs[0].Type)
		assert.Equal(t, autonity.PN, onChainProofs[0].Rule)
		assert.Equal(t, maliciousProposal.Signature, onChainProofs[0].Message.Signature)
		assert.Equal(t, preCommit.Signature, onChainProofs[0].Evidence[0].Signature)
	})

	t.Run("RunRule address the misbehaviour of PO rule, the old value proposed is not locked", func(t *testing.T) {
		// ------------Old Proposal------------
		// PO: (Mr′<r,PV) ∧ (Mr′,PC|pi) ∧ (Mr′<r′′<r,P C|pi)∗ <--- (Mr,P|pi)
		// PO1: [#(Mr′,PV|V) ≥ 2f+ 1] ∧ [nil ∨ V ∨ ⊥] ∧ [nil ∨ ⊥] <--- [V]

		// to address below scenario:
		// Is there a precommit for a value other than nil or the proposed value
		// by the current proposer in the valid round? If there is the proposer
		// has proposed a value for which it is not locked on, thus a proof of
		// misbehaviour can be generated.

		fd := NewFaultDetector(nil, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}))
		quorum := bft.Quorum(totalPower)

		// simulate a init proposal at r: 0, with v1.
		initProposal := newProposalMessage(height, 0, -1, keys[committee[1].Address], committee, nil)
		_, err := fd.msgStore.Save(initProposal)
		assert.NoError(t, err)

		// simulate quorum preVotes at r: 0 for v1.
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, 0, msgPrevote, keys[committee[i].Address], initProposal.Value(), committee)
			_, err = fd.msgStore.Save(preVote)
			assert.NoError(t, err)
		}

		// simulate a preCommit at r: 0 for v1 for the node who is going to be addressed as
		// malicious on rule PO for proposing an old value which was not locked at all.
		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, initProposal.Value(), committee)
		_, err = fd.msgStore.Save(preCommit)
		assert.NoError(t, err)

		// simulate malicious proposal at r: 1, vith v2 which was not locked at all.
		// simulate a init proposal at r: 0, with v1.
		maliciousProposal := newProposalMessage(height, 1, 0, proposerKey, committee, nil)
		_, err = fd.msgStore.Save(maliciousProposal)
		assert.NoError(t, err)

		// run rule engine.
		onChainProofs := fd.runRulesOverHeight(height, quorum)
		assert.Equal(t, 2, len(onChainProofs))
		assert.Equal(t, autonity.Misbehaviour, onChainProofs[0].Type)
		assert.Equal(t, autonity.Accusation, onChainProofs[1].Type)
		assert.Equal(t, autonity.PO, onChainProofs[0].Rule)
		assert.Equal(t, autonity.PO, onChainProofs[1].Rule)
		assert.Equal(t, maliciousProposal.Signature, onChainProofs[0].Message.Signature)
		assert.Equal(t, preCommit.Signature, onChainProofs[0].Evidence[0].Signature)
		assert.Equal(t, maliciousProposal.Signature, onChainProofs[1].Message.Signature)
	})

	t.Run("RunRule address the misbehaviour of PO rule, the valid round proposed is not correct", func(t *testing.T) {
		// ------------Old Proposal------------
		// PO: (Mr′<r,PV) ∧ (Mr′,PC|pi) ∧ (Mr′<r′′<r,P C|pi)∗ <--- (Mr,P|pi)
		// PO1: [#(Mr′,PV|V) ≥ 2f+ 1] ∧ [nil ∨ V ∨ ⊥] ∧ [nil ∨ ⊥] <--- [V]

		// To address below scenario:
		// Is there a precommit for anything other than nil from the proposer
		// between the valid round and the round of the proposal? If there is
		// then that implies the proposer saw 2f+1 prevotes in that round and
		// hence it should have set that round as the valid round.

		fd := NewFaultDetector(nil, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}))
		quorum := bft.Quorum(totalPower)
		proposer1 := keys[committee[1].Address]
		maliciousProposer := keys[committee[2].Address]

		header := newBlockHeader(height, committee)
		block := types.NewBlockWithHeader(header)

		// simulate a init proposal at r: 0, with v.
		initProposal := newProposalMessage(height, 0, -1, proposerKey, committee, block)
		_, err := fd.msgStore.Save(initProposal)
		assert.NoError(t, err)

		// simulate quorum preVotes for init proposal
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, 0, msgPrevote, keys[committee[i].Address], initProposal.Value(), committee)
			_, err = fd.msgStore.Save(preVote)
			assert.NoError(t, err)
		}

		// simulate a preCommit for init proposal of proposer1, now valid round == 0.
		preCommit1 := newVoteMsg(height, 0, msgPrecommit, proposer1, initProposal.Value(), committee)
		_, err = fd.msgStore.Save(preCommit1)
		assert.NoError(t, err)

		// assume round changes happens, now proposer1 propose old value with vr: 0.
		// simulate a new proposal at r: 3, with v.
		proposal1 := newProposalMessage(height, 3, 0, proposer1, committee, block)
		_, err = fd.msgStore.Save(proposal1)
		assert.NoError(t, err)

		// now quorum preVotes for proposal1, it makes valid round change to 3.
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, 3, msgPrevote, keys[committee[i].Address], initProposal.Value(), committee)
			_, err = fd.msgStore.Save(preVote)
			assert.NoError(t, err)
		}

		// the malicious proposer did preCommit on this round, make its valid round == 3
		preCommit := newVoteMsg(height, 3, msgPrecommit, maliciousProposer, initProposal.Value(), committee)
		_, err = fd.msgStore.Save(preCommit)
		assert.NoError(t, err)

		// malicious proposer propose at r: 5, with v and a vr: 0 which was not correct anymore.
		maliciousProposal := newProposalMessage(height, 5, 0, maliciousProposer, committee, block)
		_, err = fd.msgStore.Save(maliciousProposal)
		assert.NoError(t, err)

		// run rule engine.
		onChainProofs := fd.runRulesOverHeight(height, quorum)
		assert.Equal(t, 1, len(onChainProofs))
		assert.Equal(t, autonity.Misbehaviour, onChainProofs[0].Type)
		assert.Equal(t, autonity.PO, onChainProofs[0].Rule)
		assert.Equal(t, maliciousProposal.Signature, onChainProofs[0].Message.Signature)
		assert.Equal(t, preCommit.Signature, onChainProofs[0].Evidence[0].Signature)
	})

	t.Run("RunRule address the Accusation of PO rule, no quorum preVotes presented on the valid round", func(t *testing.T) {
		// ------------Old Proposal------------
		// PO: (Mr′<r,PV) ∧ (Mr′,PC|pi) ∧ (Mr′<r′′<r,P C|pi)∗ <--- (Mr,P|pi)
		// PO1: [#(Mr′,PV|V) ≥ 2f+ 1] ∧ [nil ∨ V ∨ ⊥] ∧ [nil ∨ ⊥] <--- [V]

		// To address below accusation scenario:
		// If proposer rise an old proposal, then there must be a quorum preVotes on the valid round.
		// Do we see a quorum of preVotes in the valid round, if not we can raise an accusation, since we cannot be sure
		// that these preVotes don't exist

		fd := NewFaultDetector(nil, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}))
		quorum := bft.Quorum(totalPower)

		header := newBlockHeader(height, committee)
		block := types.NewBlockWithHeader(header)

		// simulate an old proposal at r: 2, with v and vr: 0.
		oldProposal := newProposalMessage(height, 2, 0, proposerKey, committee, block)
		_, err := fd.msgStore.Save(oldProposal)
		assert.NoError(t, err)

		// run rule engine.
		onChainProofs := fd.runRulesOverHeight(height, quorum)
		assert.Equal(t, 1, len(onChainProofs))
		assert.Equal(t, autonity.Accusation, onChainProofs[0].Type)
		assert.Equal(t, autonity.PO, onChainProofs[0].Rule)
		assert.Equal(t, oldProposal.Signature, onChainProofs[0].Message.Signature)
	})

	t.Run("RunRule address the accusation of PVN, no corresponding proposal of preVote", func(t *testing.T) {
		// PVN: (Mr′<r,PC|pi)∧(Mr′<r′′<r,PC|pi)* ∧ (Mr,P|proposer(r)) <--- (Mr,PV|pi)
		// To address below accusation scenario:
		// If there an proVote for a non nil value, then there must be a corresponding proposal at the same round,
		// otherwise an accusation of PVN should be rise.
		fd := NewFaultDetector(nil, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}))
		quorum := bft.Quorum(totalPower)

		// simulate a preVote for v at round, let's the corresponding proposal missing.
		preVote := newVoteMsg(height, round, msgPrevote, keys[committee[1].Address], noneNilValue, committee)
		_, err := fd.msgStore.Save(preVote)
		assert.NoError(t, err)

		// run rule engine.
		onChainProofs := fd.runRulesOverHeight(height, quorum)
		assert.Equal(t, 1, len(onChainProofs))
		assert.Equal(t, autonity.Accusation, onChainProofs[0].Type)
		assert.Equal(t, autonity.PVN, onChainProofs[0].Rule)
		assert.Equal(t, preVote.Signature, onChainProofs[0].Message.Signature)
	})

	t.Run("RunRule address the misbehaviour of PVN, node preVote for value V1 while it preCommitted another value at previous round", func(t *testing.T) {
		// PVN: (Mr′<r,PC|pi)∧(Mr′<r′′<r,PC|pi)* ∧ (Mr,P|proposer(r)) <--- (Mr,PV|pi)
		// PVN2: [nil ∨ ⊥] ∧ [nil ∨ ⊥] ∧ [V:Valid(V)] <--- [V]: r′= 0,∀r′′< r:Mr′′,PC|pi=nil
		// PVN2, If there is a valid proposal V at round r, and pi never
		// ever precommit(locked a value) before, then pi should prevote
		// for V or a nil in case of timeout at this round.

		// To address below misbehaviour scenario:
		// Node preCommitted at v1 at R_x, while it preVote for v2 at R_x + n.

		fd := NewFaultDetector(nil, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}))
		quorum := bft.Quorum(totalPower)
		maliciousNode := keys[committee[1].Address]
		newProposer := keys[committee[2].Address]

		header := newBlockHeader(height, committee)
		block := types.NewBlockWithHeader(header)

		// simulate a init proposal at r: 0, with v.
		initProposal := newProposalMessage(height, 0, -1, proposerKey, committee, block)
		_, err := fd.msgStore.Save(initProposal)
		assert.NoError(t, err)

		// simulate quorum preVotes for init proposal
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, 0, msgPrevote, keys[committee[i].Address], initProposal.Value(), committee)
			_, err = fd.msgStore.Save(preVote)
			assert.NoError(t, err)
		}

		// the malicious node did preCommit for v1 on round 0
		preCommit := newVoteMsg(height, 0, msgPrecommit, maliciousNode, initProposal.Value(), committee)
		_, err = fd.msgStore.Save(preCommit)
		assert.NoError(t, err)

		// assume round changes, some one propose V2 at round 3, and malicious Node now it preVote for this V2.
		newProposal := newProposalMessage(height, 3, -1, newProposer, committee, nil)
		_, err = fd.msgStore.Save(newProposal)
		assert.NoError(t, err)

		// now the malicious node preVote for a new value V2 at higher round 3.
		preVote := newVoteMsg(height, 3, msgPrevote, maliciousNode, newProposal.Value(), committee)
		_, err = fd.msgStore.Save(preVote)
		assert.NoError(t, err)

		onChainProofs := fd.runRulesOverHeight(height, quorum)

		assert.Equal(t, 1, len(onChainProofs))
		assert.Equal(t, autonity.Misbehaviour, onChainProofs[0].Type)
		assert.Equal(t, autonity.PVN, onChainProofs[0].Rule)
		assert.Equal(t, preVote.Signature, onChainProofs[0].Message.Signature)
		assert.Equal(t, preCommit.Signature, onChainProofs[0].Evidence[0].Signature)
	})

	t.Run("RunRule address Accusation of rule C, no corresponding proposal for a preCommit msg", func(t *testing.T) {
		// C: [Mr,P|proposer(r)] ∧ [Mr,PV] <--- [Mr,PC|pi]
		// C1: [V:Valid(V)] ∧ [#(V) ≥ 2f+ 1] <--- [V]

		// To address below accusation scenario:
		// Node preCommit for a V at round R, but we cannot see the corresponding proposal that propose the value at
		// the same round of that preCommit msg.

		fd := NewFaultDetector(nil, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}))
		quorum := bft.Quorum(totalPower)

		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		_, err := fd.msgStore.Save(preCommit)
		assert.NoError(t, err)

		onChainProofs := fd.runRulesOverHeight(height, quorum)
		assert.Equal(t, 1, len(onChainProofs))
		assert.Equal(t, autonity.Accusation, onChainProofs[0].Type)
		assert.Equal(t, autonity.C, onChainProofs[0].Rule)
		assert.Equal(t, preCommit.Signature, onChainProofs[0].Message.Signature)
	})

	t.Run("RunRule address misbehaviour of rule C, node preCommit for V1 with having quorum preVotes of V2", func(t *testing.T) {
		// ------------Precommits------------
		// C: [Mr,P|proposer(r)] ∧ [Mr,PV] <--- [Mr,PC|pi]
		// C1: [V:Valid(V)] ∧ [#(V) ≥ 2f+ 1] <--- [V]

		// To address below misbehaviour scenario:
		// Node preCommit for a value V1, but there was more than quorum preVotes for not V1 at the same round.
		fd := NewFaultDetector(nil, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}))
		quorum := bft.Quorum(totalPower)
		maliciousNode := keys[committee[1].Address]

		header := newBlockHeader(height, committee)
		block := types.NewBlockWithHeader(header)

		// simulate a init proposal at r: 0, with v.
		initProposal := newProposalMessage(height, 0, -1, proposerKey, committee, block)
		_, err := fd.msgStore.Save(initProposal)
		assert.NoError(t, err)

		// simulate more than quorum preVotes for not v.
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, 0, msgPrevote, keys[committee[i].Address], nilValue, committee)
			_, err = fd.msgStore.Save(preVote)
			assert.NoError(t, err)
		}

		// malicious node preCommit to v even through there was no quorum preVotes for v.
		preCommit := newVoteMsg(height, 0, msgPrecommit, maliciousNode, initProposal.Value(), committee)
		_, err = fd.msgStore.Save(preCommit)
		assert.NoError(t, err)

		onChainProofs := fd.runRulesOverHeight(height, quorum)
		assert.Equal(t, 1, len(onChainProofs))
		assert.Equal(t, autonity.Misbehaviour, onChainProofs[0].Type)
		assert.Equal(t, autonity.C, onChainProofs[0].Rule)
		assert.Equal(t, preCommit.Signature, onChainProofs[0].Message.Signature)

		// validate if there is enough preVotes for not v.
		assert.GreaterOrEqual(t, uint64(len(onChainProofs[0].Evidence)), quorum)
		for _, m := range onChainProofs[0].Evidence {
			assert.Equal(t, height, m.H())
			assert.Equal(t, int64(0), m.R())
			assert.Equal(t, msgPrevote, m.Code)
			assert.Equal(t, nilValue, m.Value())
		}
	})

	t.Run("RunRule address accusation of rule C1, no present of quorum preVotes of V to justify the preCommit of V", func(t *testing.T) {
		// ------------Precommits------------
		// C: [Mr,P|proposer(r)] ∧ [Mr,PV] <--- [Mr,PC|pi]
		// C1: [V:Valid(V)] ∧ [#(V) ≥ 2f+ 1] <--- [V]

		// To address below accusation scenario:
		// Node preCommit for a value V, but observer haven't seen quorum preVotes for V at the round, an accusation shall
		// be rise.
		fd := NewFaultDetector(nil, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}))
		quorum := bft.Quorum(totalPower)
		maliciousNode := keys[committee[1].Address]

		header := newBlockHeader(height, committee)
		block := types.NewBlockWithHeader(header)

		// simulate a init proposal at r: 0, with v.
		initProposal := newProposalMessage(height, 0, -1, proposerKey, committee, block)
		_, err := fd.msgStore.Save(initProposal)
		assert.NoError(t, err)

		// malicious node preCommit to v even through there was no quorum preVotes for v.
		preCommit := newVoteMsg(height, 0, msgPrecommit, maliciousNode, initProposal.Value(), committee)
		_, err = fd.msgStore.Save(preCommit)
		assert.NoError(t, err)

		onChainProofs := fd.runRulesOverHeight(height, quorum)
		assert.Equal(t, 1, len(onChainProofs))
		assert.Equal(t, autonity.Accusation, onChainProofs[0].Type)
		assert.Equal(t, autonity.C1, onChainProofs[0].Rule)
		assert.Equal(t, preCommit.Signature, onChainProofs[0].Message.Signature)
	})
}
