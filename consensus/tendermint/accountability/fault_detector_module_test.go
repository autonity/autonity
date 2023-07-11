package accountability

import (
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	proto "github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	tendermintCore "github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewProposalAccountabilityCheck(t *testing.T) {
	committee, keys := generateCommittee()
	height := uint64(0)
	pi := keys[committee[0].Address]

	newProposal0 := newProposalMessage(height, 3, -1, pi, committee, nil)
	nonNilPrecommit0 := newVoteMsg(height, 1, proto.MsgPrecommit, pi, common.BytesToHash([]byte("test")), committee)
	nilPrecommit0 := newVoteMsg(height, 1, proto.MsgPrecommit, pi, common.Hash{}, committee)

	newProposal1 := newProposalMessage(height, 5, -1, pi, committee, nil)
	nilPrecommit1 := newVoteMsg(height, 3, proto.MsgPrecommit, pi, common.Hash{}, committee)

	newProposal0E := newProposalMessage(height, 3, 1, pi, committee, nil)

	t.Run("misbehaviour when pi has sent a non-nil precommit in a previous round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposal0)
		fd.msgStore.Save(nonNilPrecommit0)

		expectedProof := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PN,
			Evidences: []*message.Message{nonNilPrecommit0},
			Message:   newProposal0.ToLightProposal(),
		}

		proofs := fd.newProposalsAccountabilityCheck(0)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof, actualProof)
	})

	t.Run("no proof is returned when proposal is equivocated", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposal0)
		fd.msgStore.Save(nonNilPrecommit0)
		fd.msgStore.Save(newProposal0E)

		proofs := fd.newProposalsAccountabilityCheck(0)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi proposes a new proposal and no precommit has been sent", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposal0)
		fd.msgStore.Save(newProposal1)

		proofs := fd.newProposalsAccountabilityCheck(0)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi proposes a new proposal and has sent nil precommits in previous rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposal0)
		fd.msgStore.Save(nilPrecommit0)
		fd.msgStore.Save(newProposal1)
		fd.msgStore.Save(nilPrecommit1)

		proofs := fd.newProposalsAccountabilityCheck(0)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("multiple proof of misbehaviours when pi has sent non-nil precommits in previous rounds for multiple proposals", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposal0)
		fd.msgStore.Save(nonNilPrecommit0)
		fd.msgStore.Save(newProposal1)

		expectedProof0 := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PN,
			Evidences: []*message.Message{nonNilPrecommit0},
			Message:   newProposal0,
		}

		expectedProof1 := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PN,
			Evidences: []*message.Message{nonNilPrecommit0},
			Message:   newProposal1,
		}

		proofs := fd.newProposalsAccountabilityCheck(0)
		require.Equal(t, 2, len(proofs))

		// The order of proofs is non know apriori
		for _, p := range proofs {
			if p.Message == expectedProof0.Message {
				require.Equal(t, expectedProof0, p)
			}

			if p.Message == expectedProof1.Message {
				// The Evidences list elements can be returned in any order therefore when we have evidence which includes
				// multiple messages we need to check that each message is present separately
				require.Equal(t, expectedProof1.Type, p.Type)
				require.Equal(t, expectedProof1.Rule, p.Rule)
				require.Equal(t, expectedProof1.Message, p.Message)
				require.Contains(t, p.Evidences, nonNilPrecommit0)
			}
		}
	})
}

func TestOldProposalsAccountabilityCheck(t *testing.T) {
	//t.Skip("skip this test for CI jobs, it works in local, anyway there is similar case under misbehaviour_detector_test.go.")
	committee, keys := generateCommittee()
	quorum := bft.Quorum(committee.TotalVotingPower())
	height := uint64(0)
	pi := keys[committee[0].Address]

	header := newBlockHeader(height, committee)
	block := types.NewBlockWithHeader(header)
	header1 := newBlockHeader(height, committee)
	block1 := types.NewBlockWithHeader(header1)

	oldProposal0 := newProposalMessage(height, 3, 0, pi, committee, block)
	oldProposal5 := newProposalMessage(height, 5, 2, pi, committee, block)
	oldProposal0E := newProposalMessage(height, 3, 2, pi, committee, block1)
	oldProposal0E2 := newProposalMessage(height, 3, 0, pi, committee, block1)

	nonNilPrecommit0V := newVoteMsg(height, 0, proto.MsgPrecommit, pi, block.Hash(), committee)
	nonNilPrecommit0VPrime := newVoteMsg(height, 0, proto.MsgPrecommit, pi, block1.Hash(), committee)
	nonNilPrecommit2VPrime := newVoteMsg(height, 2, proto.MsgPrecommit, pi, block1.Hash(), committee)
	nonNilPrecommit1 := newVoteMsg(height, 1, proto.MsgPrecommit, pi, block.Hash(), committee)

	nilPrecommit0 := newVoteMsg(height, 0, proto.MsgPrecommit, pi, nilValue, committee)

	var quorumPrevotes0VPrime []*message.Message
	for i := int64(0); i < quorum.Int64(); i++ {
		quorumPrevotes0VPrime = append(quorumPrevotes0VPrime, newVoteMsg(height, 0, proto.MsgPrevote, keys[committee[i].Address], block1.Hash(), committee))
	}

	var quorumPrevotes0V []*message.Message
	for i := int64(0); i < quorum.Int64(); i++ {
		quorumPrevotes0V = append(quorumPrevotes0V, newVoteMsg(height, 0, proto.MsgPrevote, keys[committee[i].Address], block.Hash(), committee))
	}

	var precommiteNilAfterVR []*message.Message
	for i := 1; i < 3; i++ {
		precommiteNilAfterVR = append(precommiteNilAfterVR, newVoteMsg(height, int64(i), proto.MsgPrecommit, pi, nilValue, committee))
	}

	t.Run("misbehaviour when pi precommited for a different value in valid round than in the old proposal", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		fd.msgStore.Save(nonNilPrecommit0VPrime)

		expectedProof := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PO,
			Evidences: []*message.Message{nonNilPrecommit0VPrime},
			Message:   oldProposal0.ToLightProposal(),
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof, actualProof)
	})

	t.Run("misbehaviour when pi incorrectly set the valid round with a different value than the proposal", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		fd.msgStore.Save(nonNilPrecommit2VPrime)

		expectedProof := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PO,
			Evidences: []*message.Message{nonNilPrecommit2VPrime},
			Message:   oldProposal0.ToLightProposal(),
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof, actualProof)
	})

	t.Run("misbehaviour when pi incorrectly set the valid round with the same value as the proposal", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		fd.msgStore.Save(nonNilPrecommit1)

		expectedProof := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PO,
			Evidences: []*message.Message{nonNilPrecommit1},
			Message:   oldProposal0.ToLightProposal(),
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof, actualProof)
	})

	t.Run("misbehaviour when in valid round there is a quorum of prevotes for a value different than old proposal", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		for _, m := range quorumPrevotes0VPrime {
			fd.msgStore.Save(m)
		}

		expectedProof := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PO,
			Evidences: quorumPrevotes0VPrime,
			Message:   oldProposal0.ToLightProposal(),
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof.Type, actualProof.Type)
		require.Equal(t, expectedProof.Rule, actualProof.Rule)
		require.Equal(t, expectedProof.Message, actualProof.Message)
		// The order of the evidence is not known apriori
		for _, m := range expectedProof.Evidences {
			require.Contains(t, actualProof.Evidences, m)
		}

	})

	t.Run("accusation when no prevotes for proposal value in valid round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)

		expectedProof := &Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.PO,
			Message: oldProposal0.ToLightProposal(),
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof, actualProof)
	})

	t.Run("accusation when less than quorum prevotes for proposal value in valid round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		lessThanQurorumPrevotes := quorumPrevotes0V[:len(quorumPrevotes0V)-2]
		for _, m := range lessThanQurorumPrevotes {
			fd.msgStore.Save(m)
		}

		expectedProof := &Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.PO,
			Message: oldProposal0.ToLightProposal(),
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof, actualProof)
	})

	t.Run("no proof for equivocated proposal with different valid round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		fd.msgStore.Save(oldProposal0E)

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof for equivocated proposal with same valid round however different block value", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		fd.msgStore.Save(oldProposal0E2)

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when quorum of prevotes for V in vr, precommit for V from pi in vr, and precommit nils from pi from vr+1 to r", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		for _, m := range quorumPrevotes0V {
			fd.msgStore.Save(m)
		}
		fd.msgStore.Save(nonNilPrecommit0V)
		for _, m := range precommiteNilAfterVR {
			fd.msgStore.Save(m)
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when quorum of prevotes for V in vr, precommit for V from pi in vr, and some precommit nils from pi from vr+1 to r", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		for _, m := range quorumPrevotes0V {
			fd.msgStore.Save(m)
		}
		fd.msgStore.Save(nonNilPrecommit0V)
		somePrecommits := precommiteNilAfterVR[:len(precommiteNilAfterVR)-2]
		for _, m := range somePrecommits {
			fd.msgStore.Save(m)
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when quorum of prevotes for V in vr, precommit for V from pi in vr", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		for _, m := range quorumPrevotes0V {
			fd.msgStore.Save(m)
		}
		fd.msgStore.Save(nonNilPrecommit0V)

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when quorum of prevotes for V in vr, precommit nil from pi in vr", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		for _, m := range quorumPrevotes0V {
			fd.msgStore.Save(m)
		}
		fd.msgStore.Save(nilPrecommit0)

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when quorum of prevotes for V in vr", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		for _, m := range quorumPrevotes0V {
			fd.msgStore.Save(m)
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("multiple proofs from different rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		fd.msgStore.Save(nonNilPrecommit0VPrime)

		expectedMisbehaviour := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PO,
			Evidences: []*message.Message{nonNilPrecommit0VPrime},
			Message:   oldProposal0.ToLightProposal(),
		}

		fd.msgStore.Save(oldProposal5)
		expectedAccusation := &Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.PO,
			Message: oldProposal5.ToLightProposal(),
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 2, len(proofs))
		require.Contains(t, proofs, expectedMisbehaviour)
		require.Contains(t, proofs, expectedAccusation)
	})
}

func TestPrevotesAccountabilityCheck(t *testing.T) {
	committee, keys := generateCommittee()
	quorum := bft.Quorum(committee.TotalVotingPower())
	height := uint64(0)
	pi := keys[committee[0].Address]

	header := newBlockHeader(height, committee)
	block := types.NewBlockWithHeader(header)
	header1 := newBlockHeader(height, committee)
	block1 := types.NewBlockWithHeader(header1)

	newProposalForB := newProposalMessage(height, 5, -1, keys[committee[1].Address], committee, block)

	prevoteForB := newVoteMsg(height, 5, proto.MsgPrevote, pi, block.Hash(), committee)
	prevoteForB1 := newVoteMsg(height, 5, proto.MsgPrevote, pi, block1.Hash(), committee)

	precommitForB := newVoteMsg(height, 3, proto.MsgPrecommit, pi, block.Hash(), committee)
	precommitForB1 := newVoteMsg(height, 4, proto.MsgPrecommit, pi, block1.Hash(), committee)
	precommitForB1In0 := newVoteMsg(height, 0, proto.MsgPrecommit, pi, block1.Hash(), committee)
	precommitForB1In1 := newVoteMsg(height, 1, proto.MsgPrecommit, pi, block1.Hash(), committee)
	precommitForBIn0 := newVoteMsg(height, 0, proto.MsgPrecommit, pi, block.Hash(), committee)
	precommitForBIn4 := newVoteMsg(height, 4, proto.MsgPrecommit, pi, block.Hash(), committee)

	oldProposalB10 := newProposalMessage(height, 10, 5, keys[committee[1].Address], committee, block)
	newProposalB1In5 := newProposalMessage(height, 5, -1, keys[committee[1].Address], committee, block1)
	newProposalBIn5 := newProposalMessage(height, 5, -1, keys[committee[1].Address], committee, block)

	prevoteForOldB10 := newVoteMsg(height, 10, proto.MsgPrevote, pi, block.Hash(), committee)

	precommitForB1In8 := newVoteMsg(height, 8, proto.MsgPrecommit, pi, block1.Hash(), committee)
	precommitForBIn7 := newVoteMsg(height, 7, proto.MsgPrecommit, pi, block.Hash(), committee)

	t.Run("accusation when there are no corresponding proposals", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(prevoteForB)

		expectedAccusation := &Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.PVN,
			Message: prevoteForB,
		}
		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		require.Contains(t, proofs, expectedAccusation)
	})

	// Testcases for PVN
	t.Run("misbehaviour when pi precommited for a different value in a previous round than the prevoted value", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(precommitForB1)

		expectedMisbehaviour := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PVN,
			Evidences: []*message.Message{newProposalForB.ToLightProposal(), precommitForB1},
			Message:   prevoteForB,
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		require.Equal(t, expectedMisbehaviour, proofs[0])
	})

	t.Run("misbehaviour when pi precommited for a different value in a previous round than the prevoted value, after a flip flop", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(precommitForB1)
		fd.msgStore.Save(precommitForB)

		expectedMisbehaviour := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PVN,
			Evidences: []*message.Message{newProposalForB.ToLightProposal(), precommitForB1},
			Message:   prevoteForB,
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		require.Equal(t, expectedMisbehaviour, proofs[0])
	})

	t.Run("misbehaviour when pi precommited for a different value in a previous round than the prevoted value while precommit nils in middle rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(precommitForB1In0)

		var precommitNilsAfter0 []*message.Message
		for i := 1; i < 5; i++ {
			precommitNil := newVoteMsg(height, int64(i), proto.MsgPrecommit, pi, nilValue, committee)
			precommitNilsAfter0 = append(precommitNilsAfter0, precommitNil)
			fd.msgStore.Save(precommitNil)
		}

		expectedMisbehaviour := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PVN,
			Evidences: []*message.Message{precommitForB1In0},
			Message:   prevoteForB,
		}
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, precommitNilsAfter0...)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedMisbehaviour.Type, actualProof.Type)
		require.Equal(t, expectedMisbehaviour.Rule, actualProof.Rule)
		require.Equal(t, expectedMisbehaviour.Message, actualProof.Message)
		for _, m := range expectedMisbehaviour.Evidences {
			require.Contains(t, actualProof.Evidences, m)
		}
	})

	t.Run("misbehaviour when pi precommited for a different value in a previous round than the prevoted value, after a flip flop, while precommit nils in middle rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(precommitForBIn0)
		fd.msgStore.Save(precommitForB1In1)

		var precommitNilsAfter1 []*message.Message
		for i := 2; i < 5; i++ {
			precommitNil := newVoteMsg(height, int64(i), proto.MsgPrecommit, pi, nilValue, committee)
			precommitNilsAfter1 = append(precommitNilsAfter1, precommitNil)
			fd.msgStore.Save(precommitNil)
		}

		expectedMisbehaviour := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PVN,
			Evidences: []*message.Message{precommitForB1In1},
			Message:   prevoteForB,
		}
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, precommitNilsAfter1...)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedMisbehaviour.Type, actualProof.Type)
		require.Equal(t, expectedMisbehaviour.Rule, actualProof.Rule)
		require.Equal(t, expectedMisbehaviour.Message, actualProof.Message)
		for _, m := range expectedMisbehaviour.Evidences {
			require.Contains(t, actualProof.Evidences, m)
		}
	})

	t.Run("no proof when pi precommited for the same value as the prevoted value in a previous round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(precommitForBIn4)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi precommited for the same value as the prevoted value in a previous round with missing precommits in middle rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(precommitForBIn0)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi precommited for the same value as the prevoted value in a previous round with some missing precommits and precommit nils in middle rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(precommitForBIn0)
		fd.msgStore.Save(newVoteMsg(height, 3, proto.MsgPrecommit, pi, nilValue, committee))

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi precommited for the same value as the prevoted value in a previous round with no missing precommits in middle rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(precommitForBIn0)
		for i := 1; i < 5; i++ {
			precommitNil := newVoteMsg(height, int64(i), proto.MsgPrecommit, pi, nilValue, committee)
			fd.msgStore.Save(precommitNil)
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	// Testcases for PVO
	t.Run("accusation when there is no quorum for the prevote value in the valid round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(prevoteForOldB10)

		expectedAccusation := &Proof{
			Type:      autonity.Accusation,
			Rule:      autonity.PVO,
			Message:   prevoteForOldB10,
			Evidences: []*message.Message{oldProposalB10.ToLightProposal()},
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		acutalProof := proofs[0]
		require.Equal(t, expectedAccusation, acutalProof)
	})

	t.Run("misbehaviour when pi prevotes for an old proposal while in the valid round there is quorum for different value", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		// Need to add this new proposal in valid round so that unwanted accusation are not returned by the prevotes
		// accountability check method. Since we are adding a quorum of prevotes in round 6 we also need to add a new
		// proposal in round 6 to allow for those prevotes to not return accusations.
		fd.msgStore.Save(newProposalB1In5)
		fd.msgStore.Save(prevoteForOldB10)
		// quorum of prevotes for B1 in vr = 6
		var vr5Prevotes []*message.Message
		for i := uint64(0); i < quorum.Uint64(); i++ {
			vr6Prevote := newVoteMsg(height, 5, proto.MsgPrevote, keys[committee[i].Address], block1.Hash(), committee)
			vr5Prevotes = append(vr5Prevotes, vr6Prevote)
			fd.msgStore.Save(vr6Prevote)
		}

		expectedMisbehaviour := &Proof{
			Type:    autonity.Misbehaviour,
			Rule:    autonity.PVO,
			Message: prevoteForOldB10,
		}
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, oldProposalB10.ToLightProposal())
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, vr5Prevotes...)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedMisbehaviour.Type, actualProof.Type)
		require.Equal(t, expectedMisbehaviour.Rule, actualProof.Rule)
		require.Equal(t, expectedMisbehaviour.Message, actualProof.Message)
		for _, m := range expectedMisbehaviour.Evidences {
			require.Contains(t, actualProof.Evidences, m)
		}

	})

	t.Run("misbehaviour when pi has precommited for V in a previous round however the latest precommit from pi is not for V yet pi still prevoted for V in the current round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(prevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(newVoteMsg(height, 5, proto.MsgPrevote, keys[committee[i].Address], block.Hash(), committee))
		}
		for i := newProposalBIn5.R(); i < precommitForBIn7.R(); i++ {
			fd.msgStore.Save(newVoteMsg(height, i, proto.MsgPrecommit, pi, nilValue, committee))
		}
		var precommitsFromPiAfterLatestPrecommitForB []*message.Message
		fd.msgStore.Save(precommitForBIn7)

		precommitsFromPiAfterLatestPrecommitForB = append(precommitsFromPiAfterLatestPrecommitForB, precommitForBIn7)
		fd.msgStore.Save(precommitForB1In8)
		precommitsFromPiAfterLatestPrecommitForB = append(precommitsFromPiAfterLatestPrecommitForB, precommitForB1In8)
		p := newVoteMsg(height, precommitForB1In8.R()+1, proto.MsgPrecommit, pi, nilValue, committee)
		fd.msgStore.Save(p)
		precommitsFromPiAfterLatestPrecommitForB = append(precommitsFromPiAfterLatestPrecommitForB, p)

		expectedMisbehaviour := &Proof{
			Type:    autonity.Misbehaviour,
			Rule:    autonity.PVO12,
			Message: prevoteForOldB10,
		}
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, oldProposalB10.ToLightProposal())
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, precommitsFromPiAfterLatestPrecommitForB...)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedMisbehaviour.Type, actualProof.Type)
		require.Equal(t, expectedMisbehaviour.Rule, actualProof.Rule)
		require.Equal(t, expectedMisbehaviour.Message, actualProof.Message)
		for _, m := range expectedMisbehaviour.Evidences {
			require.Contains(t, actualProof.Evidences, m)
		}
	})

	t.Run("no proof when pi has precommited for V in a previous round and precommit nils afterwards", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(prevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(newVoteMsg(height, 5, proto.MsgPrevote, keys[committee[i].Address], block.Hash(), committee))
		}
		fd.msgStore.Save(precommitForBIn7)
		for i := precommitForBIn7.R() + 1; i < oldProposalB10.R(); i++ {
			fd.msgStore.Save(newVoteMsg(height, i, proto.MsgPrecommit, pi, nilValue, committee))
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))

	})

	t.Run("no proof when pi has precommited for V in a previous round however the latest precommit from pi is not for V yet pi still prevoted for V in the current round"+
		" but there are missing message after latest precommit for V", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(prevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(newVoteMsg(height, 5, proto.MsgPrevote, keys[committee[i].Address], block.Hash(), committee))
		}
		fd.msgStore.Save(precommitForBIn7)
		fd.msgStore.Save(precommitForB1In8)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))

	})

	t.Run("misbehaviour when pi has never precommited for V in a previous round however pi prevoted for V which is being reproposed", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(prevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(newVoteMsg(height, 5, proto.MsgPrevote, keys[committee[i].Address], block.Hash(), committee))
		}

		var precommitsFromPiAfterVR []*message.Message
		for i := newProposalBIn5.R() + 1; i < precommitForB1In8.R(); i++ {
			p := newVoteMsg(height, i, proto.MsgPrecommit, pi, nilValue, committee)
			fd.msgStore.Save(p)
			precommitsFromPiAfterVR = append(precommitsFromPiAfterVR, p)
		}
		fd.msgStore.Save(precommitForB1In8)
		precommitsFromPiAfterVR = append(precommitsFromPiAfterVR, precommitForB1In8)
		p := newVoteMsg(height, precommitForB1In8.R()+1, proto.MsgPrecommit, pi, nilValue, committee)
		fd.msgStore.Save(p)
		precommitsFromPiAfterVR = append(precommitsFromPiAfterVR, p)

		expectedMisbehaviour := &Proof{
			Type:    autonity.Misbehaviour,
			Rule:    autonity.PVO12,
			Message: prevoteForOldB10,
		}
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, oldProposalB10.ToLightProposal())
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, precommitsFromPiAfterVR...)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedMisbehaviour.Type, actualProof.Type)
		require.Equal(t, expectedMisbehaviour.Rule, actualProof.Rule)
		require.Equal(t, expectedMisbehaviour.Message, actualProof.Message)
		for _, m := range expectedMisbehaviour.Evidences {
			require.Contains(t, actualProof.Evidences, m)
		}
	})

	t.Run("no proof when pi has never precommited for V in a previous round however has precommitted nil after VR", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(prevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(newVoteMsg(height, 5, proto.MsgPrevote, keys[committee[i].Address], block.Hash(), committee))
		}

		for i := newProposalBIn5.R() + 1; i < oldProposalB10.R(); i++ {
			fd.msgStore.Save(newVoteMsg(height, i, proto.MsgPrecommit, pi, nilValue, committee))
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi has never precommited for V in a previous round however pi prevoted for V while it has precommited for V' but there are missing precommit before precommit for V'", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(prevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(newVoteMsg(height, 5, proto.MsgPrevote, keys[committee[i].Address], block.Hash(), committee))
		}

		fd.msgStore.Save(precommitForB1In8)

		p := newVoteMsg(height, precommitForB1In8.R()+1, proto.MsgPrecommit, pi, nilValue, committee)
		fd.msgStore.Save(p)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi has never precommited for V in a previous round however pi prevoted for V while it has precommited for V' but there are missing precommit after precommit for V'", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)

		fd.msgStore.Save(prevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(newVoteMsg(height, 5, proto.MsgPrevote, keys[committee[i].Address], block.Hash(), committee))
		}

		for i := newProposalBIn5.R() + 1; i < precommitForB1In8.R(); i++ {
			fd.msgStore.Save(newVoteMsg(height, i, proto.MsgPrecommit, pi, nilValue, committee))
		}
		fd.msgStore.Save(precommitForB1In8)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("prevotes accountability check can return multiple proofs", func(t *testing.T) {
		fd := testFD()

		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(precommitForB1)
		fd.msgStore.Save(precommitForB)

		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(prevoteForOldB10)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(newVoteMsg(height, 6, proto.MsgPrevote, keys[committee[i].Address], block1.Hash(), committee))
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 2, len(proofs))
	})

	t.Run("no proof when prevote is equivocated with different values", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(prevoteForB1)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})
}

func TestPrecommitsAccountabilityCheck(t *testing.T) {
	committee, keys := generateCommittee()
	quorum := bft.Quorum(committee.TotalVotingPower())
	height := uint64(0)
	pi := keys[committee[0].Address]

	header := newBlockHeader(height, committee)
	block := types.NewBlockWithHeader(header)
	header1 := newBlockHeader(height, committee)
	block1 := types.NewBlockWithHeader(header1)

	newProposalForB := newProposalMessage(height, 2, -1, keys[committee[1].Address], committee, block)

	precommitForB := newVoteMsg(height, 2, proto.MsgPrecommit, pi, block.Hash(), committee)
	precommitForB1 := newVoteMsg(height, 2, proto.MsgPrecommit, pi, block1.Hash(), committee)
	precommitForB1In3 := newVoteMsg(height, 3, proto.MsgPrecommit, pi, block1.Hash(), committee)

	t.Run("accusation when prevotes is less than quorum", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(precommitForB)

		for i := int64(0); i < quorum.Int64()-1; i++ {
			fd.msgStore.Save(newVoteMsg(height, 2, proto.MsgPrevote, keys[committee[i].Address], block.Hash(), committee))
		}

		expectedAccusation := &Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.C1,
			Message: precommitForB,
		}
		proofs := fd.precommitsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		require.Equal(t, expectedAccusation, proofs[0])
	})

	t.Run("misbehaviour when there is a quorum for V' than what pi precommitted for", func(t *testing.T) {
		//t.Skip("not stable in CI, but work in local.")
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(precommitForB)

		var prevotesForB1 []*message.Message
		for i := int64(0); i < quorum.Int64(); i++ {
			p := newVoteMsg(height, 2, proto.MsgPrevote, keys[committee[i].Address], block1.Hash(), committee)
			fd.msgStore.Save(p)
			prevotesForB1 = append(prevotesForB1, p)
		}

		expectedMisbehaviour := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.C,
			Evidences: prevotesForB1,
			Message:   precommitForB,
		}
		proofs := fd.precommitsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedMisbehaviour.Type, actualProof.Type)
		require.Equal(t, expectedMisbehaviour.Rule, actualProof.Rule)
		require.Equal(t, expectedMisbehaviour.Message, actualProof.Message)
		for _, m := range expectedMisbehaviour.Evidences {
			require.Contains(t, actualProof.Evidences, m)
		}
	})

	t.Run("multiple proofs can be returned from precommits accountability check", func(t *testing.T) {
		//t.Skip("not stable in CI, but work in local.")
		fd := testFD()
		fd.msgStore.Save(precommitForB1In3)

		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(precommitForB)

		var prevotesForB1 []*message.Message
		for i := int64(0); i < quorum.Int64(); i++ {
			p := newVoteMsg(height, 2, proto.MsgPrevote, keys[committee[i].Address], block1.Hash(), committee)
			fd.msgStore.Save(p)
			prevotesForB1 = append(prevotesForB1, p)
		}

		expectedProof0 := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.C,
			Evidences: prevotesForB1,
			Message:   precommitForB,
		}

		expectedProof1 := &Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.C1,
			Message: precommitForB1In3,
		}
		proofs := fd.precommitsAccountabilityCheck(height, quorum)
		require.Equal(t, 2, len(proofs))

		for _, p := range proofs {
			if p.Message == expectedProof1.Message {
				require.Equal(t, expectedProof1, p)
			}

			if p.Message == expectedProof0.Message {
				// The Evidences list elements can be returned in any order therefore when we have evidence which includes
				// multiple messages we need to check that each message is present separately
				require.Equal(t, expectedProof0.Type, p.Type)
				require.Equal(t, expectedProof0.Rule, p.Rule)
				require.Equal(t, expectedProof0.Message, p.Message)

				for _, m := range expectedProof0.Evidences {
					require.Contains(t, p.Evidences, m)
				}
			}
		}
	})

	t.Run("no proof when there is enough prevotes to form a quorum", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(precommitForB)

		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(newVoteMsg(height, 2, proto.MsgPrevote, keys[committee[i].Address], block.Hash(), committee))
		}

		proofs := fd.precommitsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when there is more than quorum prevotes ", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(precommitForB)

		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(newVoteMsg(height, 2, proto.MsgPrevote, keys[committee[i].Address], block.Hash(), committee))
		}

		proofs := fd.precommitsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when precommit is equivocated with different values", func(t *testing.T) {
		//t.Skip("not stable in CI, but work in local.")
		fd := testFD()
		fd.msgStore.Save(precommitForB)
		fd.msgStore.Save(precommitForB1)

		proofs := fd.precommitsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})
}

func testFD() *FaultDetector {
	return &FaultDetector{
		msgStore: tendermintCore.NewMsgStore(),
		logger:   log.Root(),
	}
}
