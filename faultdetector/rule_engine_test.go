package faultdetector

import (
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
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

	t.Run("getInnocentProof with unprovable rule id", func(t *testing.T) {
		fd := NewFaultDetector(nil, proposer)
		var input = Proof{
			Rule: PVO,
		}

		_, err := fd.getInnocentProof(&input)
		assert.NotNil(t, err)
	})

	t.Run("GetInnocentProofOfPO have quorum preVotes", func(t *testing.T) {

		// PO: node propose an old value with an validRound, innocent proof of it should be:
		// there were quorum num of preVote for that value at the validRound.

		fd := NewFaultDetector(nil, proposer)
		fd.savePower(lastHeight, totalPower)
		// simulate a proposal message with an old value and a valid round.
		proposal := newProposalMessage(height, round, validRound, proposerKey, committee)
		_, err := fd.msgStore.Save(proposal)
		assert.NoError(t, err)

		// simulate at least quorum num of preVotes for a value at a validRound.
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, validRound, msgPrevote, keys[committee[i].Address], proposal.Value(), committee)
			_, err = fd.msgStore.Save(preVote)
			assert.NoError(t, err)
		}

		var accusation = Proof{
			Type:    Accusation,
			Rule:    PO,
			Message: *proposal,
		}

		proof, err := fd.GetInnocentProofOfPO(&accusation)
		assert.NoError(t, err)
		assert.Equal(t, uint64(Innocence), proof.Type.Uint64())
		assert.Equal(t, proposer, proof.Sender)
		assert.Equal(t, types.RLPHash(proposal.Payload()), proof.Msghash)
	})

	t.Run("GetInnocentProofOfPO no quorum preVotes", func(t *testing.T) {

		// PO: node propose an old value with an validRound, innocent proof of it should be:
		// there were quorum num of preVote for that value at the validRound.

		fd := NewFaultDetector(nil, proposer)
		fd.savePower(lastHeight, totalPower)
		// simulate a proposal message with an old value and a valid round.
		proposal := newProposalMessage(height, round, validRound, proposerKey, committee)
		_, err := fd.msgStore.Save(proposal)
		assert.NoError(t, err)

		// simulate less than quorum num of preVotes for a value at a validRound.
		preVote := newVoteMsg(height, validRound, msgPrevote, proposerKey, proposal.Value(), committee)
		_, err = fd.msgStore.Save(preVote)
		assert.NoError(t, err)

		var accusation = Proof{
			Type:    Accusation,
			Rule:    PO,
			Message: *proposal,
		}

		_, err = fd.GetInnocentProofOfPO(&accusation)
		assert.Equal(t, errNoEvidenceForPO, err)
	})

	t.Run("GetInnocentProofOfPVN have corresponding proposal", func(t *testing.T) {

		// PVN: node prevote for a none nil value, then there must be a corresponding proposal.

		fd := NewFaultDetector(nil, proposer)
		fd.savePower(lastHeight, totalPower)
		// simulate a proposal message with an old value and a valid round.
		proposal := newProposalMessage(height, round, -1, proposerKey, committee)
		_, err := fd.msgStore.Save(proposal)
		assert.NoError(t, err)

		preVote := newVoteMsg(height, round, msgPrevote, proposerKey, proposal.Value(), committee)
		_, err = fd.msgStore.Save(preVote)
		assert.NoError(t, err)

		var accusation = Proof{
			Type:    Accusation,
			Rule:    PVN,
			Message: *preVote,
		}

		proof, err := fd.GetInnocentProofOfPVN(&accusation)
		assert.NoError(t, err)
		assert.Equal(t, uint64(Innocence), proof.Type.Uint64())
		assert.Equal(t, proposer, proof.Sender)
		assert.Equal(t, types.RLPHash(preVote.Payload()), proof.Msghash)
	})

	t.Run("GetInnocentProofOfPVN have no corresponding proposal", func(t *testing.T) {

		// PVN: node prevote for a none nil value, then there must be a corresponding proposal.
		fd := NewFaultDetector(nil, proposer)
		fd.savePower(lastHeight, totalPower)

		preVote := newVoteMsg(height, round, msgPrevote, proposerKey, noneNilValue, committee)
		_, err := fd.msgStore.Save(preVote)
		assert.NoError(t, err)

		var accusation = Proof{
			Type:    Accusation,
			Rule:    PVN,
			Message: *preVote,
		}

		_, err = fd.GetInnocentProofOfPVN(&accusation)
		assert.Equal(t, errNoEvidenceForPVN, err)
	})

	t.Run("GetInnocentProofOfC have corresponding proposal", func(t *testing.T) {

		// C: node preCommit at a none nil value, there must be a corresponding proposal.

		fd := NewFaultDetector(nil, proposer)
		fd.savePower(lastHeight, totalPower)
		// simulate a proposal message with an old value and a valid round.
		proposal := newProposalMessage(height, round, -1, proposerKey, committee)
		_, err := fd.msgStore.Save(proposal)
		assert.NoError(t, err)

		preCommit := newVoteMsg(height, round, msgPrecommit, proposerKey, proposal.Value(), committee)
		_, err = fd.msgStore.Save(preCommit)
		assert.NoError(t, err)

		var accusation = Proof{
			Type:    Accusation,
			Rule:    C,
			Message: *preCommit,
		}

		proof, err := fd.GetInnocentProofOfC(&accusation)
		assert.NoError(t, err)
		assert.Equal(t, uint64(Innocence), proof.Type.Uint64())
		assert.Equal(t, proposer, proof.Sender)
		assert.Equal(t, types.RLPHash(preCommit.Payload()), proof.Msghash)
	})

	t.Run("GetInnocentProofOfC have no corresponding proposal", func(t *testing.T) {

		// C: node preCommit at a none nil value, there must be a corresponding proposal.

		fd := NewFaultDetector(nil, proposer)
		fd.savePower(lastHeight, totalPower)

		preCommit := newVoteMsg(height, round, msgPrecommit, proposerKey, noneNilValue, committee)
		_, err := fd.msgStore.Save(preCommit)
		assert.NoError(t, err)

		var accusation = Proof{
			Type:    Accusation,
			Rule:    C,
			Message: *preCommit,
		}

		_, err = fd.GetInnocentProofOfC(&accusation)
		assert.Equal(t, errNoEvidenceForC, err)
	})

	t.Run("GetInnocentProofOfC1 have quorum preVotes", func(t *testing.T) {

		// C1: node preCommit at a none nil value, there must be quorum corresponding preVotes with same value and round.

		fd := NewFaultDetector(nil, proposer)
		fd.savePower(lastHeight, totalPower)

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
			Type:    Accusation,
			Rule:    C1,
			Message: *preCommit,
		}

		proof, err := fd.GetInnocentProofOfC1(&accusation)
		assert.NoError(t, err)
		assert.Equal(t, uint64(Innocence), proof.Type.Uint64())
		assert.Equal(t, proposer, proof.Sender)
		assert.Equal(t, types.RLPHash(preCommit.Payload()), proof.Msghash)
	})

	t.Run("GetInnocentProofOfC1 have no quorum preVotes", func(t *testing.T) {

		// C1: node preCommit at a none nil value, there must be quorum corresponding preVotes with same value and round.

		fd := NewFaultDetector(nil, proposer)
		fd.savePower(lastHeight, totalPower)

		preCommit := newVoteMsg(height, round, msgPrecommit, proposerKey, noneNilValue, committee)
		_, err := fd.msgStore.Save(preCommit)
		assert.NoError(t, err)

		var accusation = Proof{
			Type:    Accusation,
			Rule:    C1,
			Message: *preCommit,
		}

		_, err = fd.GetInnocentProofOfC1(&accusation)
		assert.Equal(t, errNoEvidenceForC1, err)
	})

	t.Run("Test error to rule mapping", func(t *testing.T) {
		rule, err := errorToRule(errEquivocation)
		assert.NoError(t, err)
		assert.Equal(t, Equivocation, rule)

		rule, err = errorToRule(errProposer)
		assert.NoError(t, err)
		assert.Equal(t, InvalidProposer, rule)

		rule, err = errorToRule(errProposal)
		assert.NoError(t, err)
		assert.Equal(t, InvalidProposal, rule)

		rule, err = errorToRule(errGarbageMsg)
		assert.NoError(t, err)
		assert.Equal(t, GarbageMessage, rule)

		rule, err = errorToRule(fmt.Errorf("unknown err"))
		assert.Error(t, err)
	})
}
