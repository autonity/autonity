package accountability

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/rlp"
)

//TODO(lorenzo) add tests where we sent future height messages to the precompiled Run() methods. Maybe start from the following:
/*
t.Run("a future msg is received, expect an error of errFutureMsg", func(t *testing.T) {
		futureHeight := height + 1
		proposal := newValidatedProposalMessage(futureHeight, round, -1, signer, committee, nil, proposerIdx)
		chainMock.EXPECT().GetHeaderByNumber(height).Return(nil)
		_, err := verifyProofSignatures(chainMock, &Proof{Message: proposal})
		require.Equal(t, errFutureMsg, err)
	})

	t.Run("chain cannot provide the last header of the height that msg votes on, expect an error of errFutureMsg", func(t *testing.T) {
		proposal := newValidatedProposalMessage(height-5, round, -1, signer, committee, nil, proposerIdx)
		chainMock.EXPECT().GetHeaderByNumber(height - 6).Return(nil)
		_, err := verifyProofSignatures(chainMock, &Proof{Message: proposal})
		require.Equal(t, errFutureMsg, err)
	})*/
var (
	height              = uint64(100)
	defRound            = int64(0)
	defValidRound       = int64(-1)
	newRound            = int64(4)
	validRound          = int64(2)
	noneNilValue        = common.Hash{0x1}
	defNewProposal      = newValidatedProposalMessage(height, defRound, defValidRound, signer, committee, nil, proposerIdx)
	defLightNewProposal = defNewProposal.ToLight()
	defOldProposal      = newValidatedProposalMessage(height, newRound, validRound, signer, committee, nil, proposerIdx)
	defLightOldProposal = defOldProposal.ToLight()
	oldProposal2        = newValidatedProposalMessage(height, newRound-1, validRound, signer, committee, nil, proposerIdx)
	oldLightProposal2   = oldProposal2.ToLight()

	prevoteForOldProposal1 = newValidatedPrevote(newRound, height, defOldProposal.Value(), signer, self, cSize)
	prevoteForOldProposal2 = newValidatedPrevote(newRound, height, defOldProposal.Value(), makeSigner(keys[1]), &committee.Members[1], cSize)
	aggPrevoteForOld       = message.AggregatePrevotes([]message.Vote{prevoteForOldProposal1, prevoteForOldProposal2})

	nilPrevote1     = newValidatedPrevote(defRound, height, nilValue, signer, self, cSize)
	nilPrevote2     = newValidatedPrevote(defRound, height, nilValue, makeSigner(keys[1]), &committee.Members[1], cSize)
	aggNilPrevote   = message.AggregatePrevotes([]message.Vote{nilPrevote1, nilPrevote2})
	nilPrecommit1   = newValidatedPrecommit(defRound, height, nilValue, signer, self, cSize)
	nilPrecommit2   = newValidatedPrecommit(defRound, height, nilValue, makeSigner(keys[1]), &committee.Members[1], cSize)
	aggNilPrecommit = message.AggregatePrecommits([]message.Vote{nilPrecommit1, nilPrecommit2})

	prevote1     = newValidatedPrevote(defRound, height, defNewProposal.Value(), signer, self, cSize)
	prevote2     = newValidatedPrevote(defRound, height, defNewProposal.Value(), makeSigner(keys[1]), &committee.Members[1], cSize)
	aggPrevote   = message.AggregatePrevotes([]message.Vote{prevote1, prevote2})
	precommit1   = newValidatedPrecommit(defRound, height, defNewProposal.Value(), signer, self, cSize)
	precommit2   = newValidatedPrecommit(defRound, height, defNewProposal.Value(), makeSigner(keys[1]), &committee.Members[1], cSize)
	aggPrecommit = message.AggregatePrecommits([]message.Vote{precommit1, precommit2})

	committee2, keys2, _ = generateCommittee()
	proposal2            = newValidatedProposalMessage(height, defRound, defValidRound, makeSigner(keys2[0]), committee2, nil, proposerIdx)
	invalidPrecommit     = newValidatedPrecommit(defRound, height, proposal2.Value(), makeSigner(keys2[0]), &committee2.Members[0], committee2.Len())
)

func TestContractsManagement(t *testing.T) {
	// register contracts into evm package.
	LoadPrecompiles(nil)
	assert.NotNil(t, vm.PrecompiledContractsByzantium[checkInnocenceAddress])
	assert.NotNil(t, vm.PrecompiledContractsByzantium[checkMisbehaviourAddress])
	assert.NotNil(t, vm.PrecompiledContractsByzantium[checkAccusationAddress])

	assert.NotNil(t, vm.PrecompiledContractsHomestead[checkInnocenceAddress])
	assert.NotNil(t, vm.PrecompiledContractsHomestead[checkMisbehaviourAddress])
	assert.NotNil(t, vm.PrecompiledContractsHomestead[checkAccusationAddress])

	assert.NotNil(t, vm.PrecompiledContractsIstanbul[checkInnocenceAddress])
	assert.NotNil(t, vm.PrecompiledContractsIstanbul[checkMisbehaviourAddress])
	assert.NotNil(t, vm.PrecompiledContractsIstanbul[checkAccusationAddress])

	assert.NotNil(t, vm.PrecompiledContractsBerlin[checkInnocenceAddress])
	assert.NotNil(t, vm.PrecompiledContractsBerlin[checkAccusationAddress])
	assert.NotNil(t, vm.PrecompiledContractsBerlin[checkMisbehaviourAddress])

	assert.NotNil(t, vm.PrecompiledContractsBLS[checkInnocenceAddress])
	assert.NotNil(t, vm.PrecompiledContractsBLS[checkAccusationAddress])
	assert.NotNil(t, vm.PrecompiledContractsBLS[checkMisbehaviourAddress])

}

func TestDecodeAndVerifyProofs(t *testing.T) {
	type testCase struct {
		Proof
		outCome error
	}
	var validProof Proof
	validProof.Type = autonity.AccountabilityEventType(rand.Intn(3))
	validProof.Rule = autonity.Rule(rand.Intn(10))
	validProof.Message = defNewProposal.ToLight()
	validProof.OffenderIndex = proposerIdx
	validProof.Evidences = append(validProof.Evidences, aggPrevote, aggPrecommit)

	p2 := validProof
	p2.OffenderIndex = committee.Len()

	p3 := validProof
	p3.OffenderIndex = committee.Len() - 1

	p4 := validProof
	p4.Message = proposal2.ToLight()

	var proofWithInvalidSignature Proof
	proofWithInvalidSignature.Rule = autonity.PO
	proofWithInvalidSignature.Message = defLightOldProposal
	proofWithInvalidSignature.Evidences = append(proofWithInvalidSignature.Evidences, invalidPrecommit)

	cases := []testCase{

		{
			validProof,
			nil,
		},
		{
			p2,
			errInvalidOffenderIdx,
		},
		{
			p3,
			errProofOffender,
		},
		{
			p4,
			message.ErrUnauthorizedAddress,
		},
		{
			proofWithInvalidSignature,
			message.ErrBadSignature,
		},
	}

	for i, tc := range cases {
		proof := tc.Proof
		rp, err := rlp.EncodeToBytes(&proof)
		assert.NoError(t, err)
		decodeProof, err := decodeRawProof(rp)
		assert.NoError(t, err)
		err = verifyProofSignatures(committee, decodeProof)
		t.Log("Running TestDecodeAndVerifyProofs case", "case id", i, "actual err", err, "expected", tc.outCome)
		require.Equal(t, tc.outCome, err)
		if tc.outCome == nil {
			assert.Equal(t, tc.Proof.Rule, decodeProof.Rule)
			assert.Equal(t, tc.Proof.Message.Signature(), decodeProof.Message.Signature())
			assert.Equal(t, tc.Proof.Evidences, decodeProof.Evidences)
		}
	}
}

type testCase struct {
	proof   Proof
	outCome bool
}

func TestAccusationVerifier(t *testing.T) {
	// Todo(youssef): add integration tests for the precompile Run function

	t.Run("Test accusation verifier required gas", func(t *testing.T) {
		av := AccusationVerifier{}
		assert.Equal(t, params.AutonityAFDContractGasPerKB, av.RequiredGas(nil))
	})

	t.Run("Test accusation verifier run with nil bytes", func(t *testing.T) {
		av := AccusationVerifier{}
		ret, err := av.Run(nil, height, nil, common.Address{})
		assert.Equal(t, failureReturn, ret)
		assert.Nil(t, err)
	})

	t.Run("Test accusation verifier run with invalid rlp bytes", func(t *testing.T) {
		wrongBytes := failureReturn
		av := AccusationVerifier{}
		ret, err := av.Run(wrongBytes, height, nil, common.Address{})
		assert.Equal(t, failureReturn, ret)
		assert.Nil(t, err)
	})

	cases := []testCase{
		{
			proof: Proof{
				Rule: autonity.Equivocation + 100,
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:    autonity.PO,
				Message: aggPrevote,
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:    autonity.PVN,
				Message: defLightNewProposal,
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:    autonity.PVO,
				Message: defLightNewProposal,
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:    autonity.C1,
				Message: defLightNewProposal,
			},
			outCome: false,
		},
		// PO tests comes here:
		{
			proof: Proof{
				Rule:    autonity.PO,
				Message: defNewProposal,
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:    autonity.PO,
				Message: defLightNewProposal,
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:          autonity.PO,
				Message:       defLightOldProposal,
				OffenderIndex: proposerIdx + 1,
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:          autonity.PO,
				Message:       defLightOldProposal,
				OffenderIndex: proposerIdx,
			},
			outCome: true,
		},
		// PVN tests comes here:
		{
			proof: Proof{
				Rule:    autonity.PVN,
				Message: defLightOldProposal,
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:    autonity.PVN,
				Message: aggNilPrevote,
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:          autonity.PVN,
				Message:       aggPrevote,
				OffenderIndex: 3,
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:          autonity.PVN,
				Message:       aggPrevote,
				OffenderIndex: 0,
			},
			outCome: true,
		},
		// PVO tests comes here:
		{
			proof: Proof{
				Rule:    autonity.PVO,
				Message: defLightOldProposal,
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:    autonity.PVO,
				Message: aggNilPrevote,
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:          autonity.PVO,
				Message:       aggPrevote,
				OffenderIndex: 3,
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:          autonity.PVO,
				Message:       aggPrevote,
				OffenderIndex: 1,
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:          autonity.PVO,
				Message:       aggPrevote,
				OffenderIndex: 1,
				Evidences:     []message.Msg{defLightNewProposal},
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:          autonity.PVO,
				Message:       aggPrevote,
				OffenderIndex: 1,
				Evidences:     []message.Msg{defLightNewProposal},
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:          autonity.PVO,
				Message:       aggPrevoteForOld,
				OffenderIndex: 1,
				Evidences:     []message.Msg{oldLightProposal2},
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:          autonity.PVO,
				Message:       aggPrevoteForOld,
				OffenderIndex: 1,
				Evidences:     []message.Msg{defLightOldProposal},
			},
			outCome: true,
		},

		// C1 tests comes here:
		{
			proof: Proof{
				Rule:    autonity.C1,
				Message: aggPrevote,
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:    autonity.C1,
				Message: aggNilPrecommit,
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:          autonity.C1,
				Message:       aggPrecommit,
				OffenderIndex: 3,
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:          autonity.C1,
				Message:       aggPrecommit,
				OffenderIndex: 1,
			},
			outCome: true,
		},
	}

	for i, tc := range cases {
		proof := tc.proof
		actual := verifyAccusation(&proof, committee)
		t.Log("running TestAccusationVerifier", "case", i, "expected", tc.outCome, "actual", actual)
		assert.Equal(t, tc.outCome, actual)
	}
}

func TestMisbehaviourVerifier(t *testing.T) {
	type testCase struct {
		proof   Proof
		outCome []byte
	}

	liteNewP := newValidatedLightProposal(height, 1, -1, signer, committee, nil, proposerIdx)
	liteOldP := newValidatedLightProposal(height, 3, 0, signer, committee, nil, proposerIdx)

	prevotes := make([]message.Vote, committee.Len())
	for i := range committee.Members {
		prevotes[i] = newValidatedPrevote(0, height, noneNilValue, makeSigner(keys[i]), &committee.Members[i], cSize)
	}
	aggVote := message.AggregatePrevotes(prevotes)
	aggVoteNoQuorum := message.AggregatePrevotes(prevotes[1:2])
	fakedVote1 := newValidatedPrevote(1, height, noneNilValue, signer, self, cSize)
	fakedVote2 := newValidatedPrevote(0, height, liteOldP.Value(), signer, self, cSize)
	fakedVote3 := newValidatedPrevote(0, height, nilValue, signer, self, cSize)
	commit1 := newValidatedPrecommit(0, height, noneNilValue, signer, self, cSize)
	commit2 := newValidatedPrecommit(0, height, noneNilValue, makeSigner(keys[1]), &committee.Members[1], cSize)
	commit3 := newValidatedPrecommit(0, height, noneNilValue, makeSigner(keys[2]), &committee.Members[2], cSize)
	commit4 := newValidatedPrecommit(0, height, nilValue, signer, self, cSize)
	commit5 := newValidatedPrecommit(2, height, noneNilValue, signer, self, cSize)
	aggCommit := message.AggregatePrecommits([]message.Vote{commit1, commit2})

	// node locked at V1 at round 0.
	preCommitPVN := newValidatedPrecommit(0, height, noneNilValue, signer, self, cSize)
	preCommitR1PVN := newValidatedPrecommit(1, height, nilValue, signer, self, cSize)
	preCommitR1PVN2 := newValidatedPrecommit(1, height, nilValue, makeSigner(keys[1]), &committee.Members[1], cSize)
	aggPrecomitR1PVN := message.AggregatePrecommits([]message.Vote{preCommitR1PVN, preCommitR1PVN2})

	preCommitR2PVN := newValidatedPrecommit(2, height, nilValue, signer, self, cSize)
	proposalPVN := newValidatedLightProposal(height, 3, -1, signer, committee, nil, proposerIdx)
	// node preVote for V2 at round 3
	prevotePVN := newValidatedPrevote(3, height, proposalPVN.Value(), signer, self, cSize)

	// PVO settings
	correspondingProposalPVO := newValidatedLightProposal(height, 3, 0, signer, committee, nil, proposerIdx)
	maliciousPreVotePVO := newValidatedPrevote(3, height, correspondingProposalPVO.Value(), signer, self, cSize)
	// simulate quorum prevote for not v at valid round.
	votesPVO := make([]message.Vote, committee.Len())
	for i := range committee.Members {
		votesPVO[i] = newValidatedPrevote(0, height, noneNilValue, makeSigner(keys[i]), &committee.Members[i], cSize)
	}
	aggVotePVO := message.AggregatePrevotes(votesPVO)
	aggVotePVONoQuorum := message.AggregatePrevotes(votesPVO[2:3])

	// PVO12 settings.
	// a precommit at round 1, with value v.
	pcForVPVO12 := newValidatedPrecommit(1, height, correspondingProposalPVO.Value(), signer, self, cSize)
	// a precommit at round 2, with value not v.
	pcForNotVPVO12 := newValidatedPrecommit(2, height, noneNilValue, signer, self, cSize)
	// a prevote at round 3, with value v.
	preVotePVO12 := newValidatedPrevote(3, height, correspondingProposalPVO.Value(), signer, self, cSize)

	// Rule C settings.
	preCommitC := newValidatedPrecommit(0, height, noneNilValue, signer, self, cSize)
	preCommitNilC := newValidatedPrecommit(0, height, nilValue, signer, self, cSize)
	votesC := make([]message.Vote, committee.Len())
	for i := range committee.Members {
		votesC[i] = newValidatedPrevote(0, height, common.Hash{0x2}, makeSigner(keys[i]), &committee.Members[i], cSize)
	}
	aggVoteC := message.AggregatePrevotes(votesC)
	aggVoteCNoQuorum := message.AggregatePrevotes(votesC[2:3])

	t.Run("Test misbehaviour verifier required gas", func(t *testing.T) {
		mv := MisbehaviourVerifier{}
		assert.Equal(t, params.AutonityAFDContractGasPerKB, mv.RequiredGas(nil))
	})

	t.Run("Test misbehaviour verifier run with nil bytes", func(t *testing.T) {
		mv := MisbehaviourVerifier{}
		ret, err := mv.Run(nil, height, nil, common.Address{})
		assert.Equal(t, failureReturn, ret)
		assert.Nil(t, err)
	})

	t.Run("Test misbehaviour verifier run with invalid rlp bytes", func(t *testing.T) {
		wrongBytes := failureReturn
		mv := MisbehaviourVerifier{}
		ret, err := mv.Run(wrongBytes, height, nil, common.Address{})
		assert.Equal(t, failureReturn, ret)
		assert.Nil(t, err)
	})

	tests := []testCase{
		// test proof of misbehaviour of PN handling comes here:
		{
			proof: Proof{
				Rule:          autonity.PN,
				Message:       newValidatedLightProposal(height, 1, 0, signer, committee, nil, proposerIdx),
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{newValidatedPrecommit(0, height, noneNilValue, signer, self, cSize)},
			},
			outCome: failureReturn,
		},
		{
			proof: Proof{
				Rule:          autonity.PN,
				Message:       liteNewP,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{newValidatedPrevote(0, height, noneNilValue, signer, self, cSize)},
			},
			outCome: failureReturn,
		},
		{
			proof: Proof{
				Rule:          autonity.PN,
				Message:       liteNewP,
				OffenderIndex: 3,
				Evidences:     []message.Msg{aggCommit},
			},
			outCome: failureReturn,
		},
		{
			proof: Proof{
				Rule:          autonity.PN,
				Message:       liteNewP,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{commit3},
			},
			outCome: failureReturn,
		},
		{
			proof: Proof{
				Rule:          autonity.PN,
				Message:       liteNewP,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{commit4},
			},
			outCome: failureReturn,
		},
		{
			proof: Proof{
				Rule:          autonity.PN,
				Message:       liteNewP,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{commit5},
			},
			outCome: failureReturn,
		},
		{
			proof: Proof{
				Rule:          autonity.PN,
				Message:       liteNewP,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{aggCommit},
			},
			outCome: validReturn(liteNewP, proposer, autonity.PN),
		},
		// test proof of misbehaviour of PO handling comes here:
		{
			proof: Proof{
				Rule:          autonity.PO,
				Message:       aggCommit,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{aggCommit},
			},
			outCome: failureReturn,
		},
		{
			proof: Proof{
				Rule:          autonity.PO,
				Message:       liteNewP,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{aggCommit},
			},
			outCome: failureReturn,
		},
		{
			proof: Proof{
				Rule:          autonity.PO,
				Message:       liteOldP,
				OffenderIndex: proposerIdx,
			},
			outCome: failureReturn,
		},
		{
			proof: Proof{
				Rule:          autonity.PO,
				Message:       liteOldP,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{aggCommit},
			},
			outCome: validReturn(liteOldP, proposer, autonity.PO),
		},
		{
			proof: Proof{
				Rule:          autonity.PO,
				Message:       liteOldP,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{commit5},
			},
			outCome: validReturn(liteOldP, proposer, autonity.PO),
		},
		{
			proof: Proof{
				Rule:          autonity.PO,
				Message:       liteOldP,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{aggVote, aggVote},
			},
			outCome: failureReturn,
		},
		{
			proof: Proof{
				Rule:          autonity.PO,
				Message:       liteOldP,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{aggVote, fakedVote1},
			},
			outCome: failureReturn,
		},
		{
			proof: Proof{
				Rule:          autonity.PO,
				Message:       liteOldP,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{aggVote, fakedVote2},
			},
			outCome: failureReturn,
		},
		{
			proof: Proof{
				Rule:          autonity.PO,
				Message:       liteOldP,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{aggVoteNoQuorum},
			},
			outCome: failureReturn,
		},
		{
			proof: Proof{
				Rule:          autonity.PO,
				Message:       liteOldP,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{aggVote},
			},
			outCome: validReturn(liteOldP, proposer, autonity.PO),
		},
		// Misbehaviour of PVN tests comes here:
		{
			// no evidence
			proof: Proof{
				Rule:          autonity.PVN,
				OffenderIndex: proposerIdx,
			},
			outCome: failureReturn,
		},
		{
			// wrong message
			proof: Proof{
				Rule:          autonity.PVN,
				OffenderIndex: proposerIdx,
				Message:       liteOldP,
				Evidences:     []message.Msg{aggVote},
			},
			outCome: failureReturn,
		},
		{
			// vote signers does not contain offender idx.
			proof: Proof{
				Rule:          autonity.PVN,
				OffenderIndex: proposerIdx,
				Message:       aggVoteNoQuorum,
				Evidences:     []message.Msg{aggVote},
			},
			outCome: failureReturn,
		},
		{
			// vote for nil value
			proof: Proof{
				Rule:          autonity.PVN,
				OffenderIndex: proposerIdx,
				Message:       fakedVote3,
				Evidences:     []message.Msg{aggVote},
			},
			outCome: failureReturn,
		},
		{
			// vote for nil value
			proof: Proof{
				Rule:          autonity.PVN,
				OffenderIndex: proposerIdx,
				Message:       fakedVote3,
				Evidences:     []message.Msg{aggVote},
			},
			outCome: failureReturn,
		},
		{
			// no corresponding proposal provided in the proof.
			proof: Proof{
				Rule:          autonity.PVN,
				OffenderIndex: proposerIdx,
				Message:       prevotePVN,
				Evidences:     []message.Msg{preCommitPVN, preCommitR1PVN, preCommitR2PVN},
			},
			outCome: failureReturn,
		},
		{
			// invalid proposal provided in the proof.
			proof: Proof{
				Rule:          autonity.PVN,
				OffenderIndex: proposerIdx,
				Message:       prevotePVN,
				Evidences:     []message.Msg{liteOldP, preCommitPVN, preCommitR1PVN, preCommitR2PVN},
			},
			outCome: failureReturn,
		},
		{
			// evidence contains invalid precomit message.
			proof: Proof{
				Rule:          autonity.PVN,
				OffenderIndex: proposerIdx,
				Message:       prevotePVN,
				Evidences:     []message.Msg{proposalPVN, preCommitPVN, aggPrecomitR1PVN, preCommitR2PVN, commit1},
			},
			outCome: failureReturn,
		},
		{
			// proof contains round gaps in precommits
			proof: Proof{
				Rule:          autonity.PVN,
				OffenderIndex: proposerIdx,
				Message:       prevotePVN,
				Evidences:     []message.Msg{proposalPVN, preCommitPVN, aggPrecomitR1PVN},
			},
			outCome: failureReturn,
		},
		{
			proof: Proof{
				Rule:          autonity.PVN,
				OffenderIndex: proposerIdx,
				Message:       prevotePVN,
				Evidences:     []message.Msg{proposalPVN, preCommitPVN, aggPrecomitR1PVN, preCommitR2PVN},
			},
			outCome: validReturn(prevotePVN, proposer, autonity.PVN),
		},
		// PVO tests comes here:
		{
			// no message
			proof: Proof{
				Rule:          autonity.PVO,
				OffenderIndex: proposerIdx,
			},
			outCome: failureReturn,
		},
		{
			// invalid message type.
			proof: Proof{
				Rule:          autonity.PVO,
				Message:       liteOldP,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{proposalPVN, preCommitPVN, aggPrecomitR1PVN, preCommitR2PVN},
			},
			outCome: failureReturn,
		},
		{
			// vote signers does not contain offender idx.
			proof: Proof{
				Rule:          autonity.PVO,
				Message:       fakedVote3,
				OffenderIndex: proposerIdx + 1,
				Evidences:     []message.Msg{proposalPVN, preCommitPVN, aggPrecomitR1PVN, preCommitR2PVN},
			},
			outCome: failureReturn,
		},
		{
			// vote signers does not contain offender idx.
			proof: Proof{
				Rule:          autonity.PVO,
				Message:       fakedVote3,
				OffenderIndex: proposerIdx + 1,
				Evidences:     []message.Msg{proposalPVN, preCommitPVN, aggPrecomitR1PVN, preCommitR2PVN},
			},
			outCome: failureReturn,
		},
		{
			// vote for nil cannot be accountable
			proof: Proof{
				Rule:          autonity.PVO,
				Message:       fakedVote3,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{proposalPVN, preCommitPVN, aggPrecomitR1PVN, preCommitR2PVN},
			},
			outCome: failureReturn,
		},
		{
			// there is no corresponding proposal in the proof.
			proof: Proof{
				Rule:          autonity.PVO,
				Message:       fakedVote2,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{preCommitPVN, aggPrecomitR1PVN, preCommitR2PVN},
			},
			outCome: failureReturn,
		},
		{
			// there is no corresponding proposal in the proof.
			proof: Proof{
				Rule:          autonity.PVO,
				Message:       fakedVote2,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{preCommitPVN, aggPrecomitR1PVN, preCommitR2PVN},
			},
			outCome: failureReturn,
		},
		{
			// there is invalid proposal in the proof.
			proof: Proof{
				Rule:          autonity.PVO,
				Message:       maliciousPreVotePVO,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{oldLightProposal2, aggVotePVO},
			},
			outCome: failureReturn,
		},
		{
			// there is invalid prvotes in the proof
			proof: Proof{
				Rule:          autonity.PVO,
				Message:       maliciousPreVotePVO,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{correspondingProposalPVO, aggPrevote},
			},
			outCome: failureReturn,
		},
		{
			// there is duplicated vote in the proof.
			proof: Proof{
				Rule:          autonity.PVO,
				Message:       maliciousPreVotePVO,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{correspondingProposalPVO, aggVotePVO, aggVotePVO},
			},
			outCome: failureReturn,
		},
		{
			// there is less than quorum prevotes in the proof.
			proof: Proof{
				Rule:          autonity.PVO,
				Message:       maliciousPreVotePVO,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{correspondingProposalPVO, aggVotePVONoQuorum},
			},
			outCome: failureReturn,
		},
		{
			proof: Proof{
				Rule:          autonity.PVO,
				Message:       maliciousPreVotePVO,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{correspondingProposalPVO, aggVotePVO},
			},
			outCome: validReturn(maliciousPreVotePVO, proposer, autonity.PVO),
		},
		// PVO12 tests comes here:
		{
			// no evidence.
			proof: Proof{
				Rule: autonity.PVO12,
			},
			outCome: failureReturn,
		},
		{
			// wrong msg type
			proof: Proof{
				Rule:      autonity.PVO12,
				Message:   liteOldP,
				Evidences: []message.Msg{correspondingProposalPVO, pcForVPVO12, pcForNotVPVO12},
			},
			outCome: failureReturn,
		},
		{
			// wrong offender index
			proof: Proof{
				Rule:          autonity.PVO12,
				Message:       preVotePVO12,
				OffenderIndex: proposerIdx + 1,
				Evidences:     []message.Msg{correspondingProposalPVO, pcForVPVO12, pcForNotVPVO12},
			},
			outCome: failureReturn,
		},
		{
			// prevote for nil is not accountable
			proof: Proof{
				Rule:          autonity.PVO12,
				Message:       fakedVote3,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{correspondingProposalPVO, pcForVPVO12, pcForNotVPVO12},
			},
			outCome: failureReturn,
		},
		{
			// there is no corresponding proposal from the proof
			proof: Proof{
				Rule:          autonity.PVO12,
				Message:       preVotePVO12,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{pcForVPVO12, pcForNotVPVO12},
			},
			outCome: failureReturn,
		},
		{
			// there is invalid proposal from the proof
			proof: Proof{
				Rule:          autonity.PVO12,
				Message:       preVotePVO12,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{liteOldP, pcForVPVO12, pcForNotVPVO12},
			},
			outCome: failureReturn,
		},
		{
			// there is invalid precommits from the proof
			proof: Proof{
				Rule:          autonity.PVO12,
				Message:       preVotePVO12,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{correspondingProposalPVO, aggPrecommit, pcForVPVO12, pcForNotVPVO12},
			},
			outCome: failureReturn,
		},
		{
			// there is invalid msg from the proof
			proof: Proof{
				Rule:          autonity.PVO12,
				Message:       preVotePVO12,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{correspondingProposalPVO, aggPrevote, pcForVPVO12, pcForNotVPVO12},
			},
			outCome: failureReturn,
		},
		{
			// there is missing round of precommits from the proof
			proof: Proof{
				Rule:          autonity.PVO12,
				Message:       preVotePVO12,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{correspondingProposalPVO, pcForVPVO12},
			},
			outCome: failureReturn,
		},
		{
			proof: Proof{
				Rule:          autonity.PVO12,
				Message:       preVotePVO12,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{correspondingProposalPVO, pcForVPVO12, pcForNotVPVO12},
			},
			outCome: validReturn(preVotePVO12, proposer, autonity.PVO12),
		},
		// C1 misbehaviour proof handling starts from here:
		{
			// no evidence
			proof: Proof{
				Rule: autonity.C,
			},
			outCome: failureReturn,
		},
		{
			// wrong message type.
			proof: Proof{
				Rule:      autonity.C,
				Evidences: []message.Msg{aggVoteC},
			},
			outCome: failureReturn,
		},
		{
			// msg signer does not contain offender.
			proof: Proof{
				Rule:          autonity.C,
				Message:       preCommitC,
				OffenderIndex: proposerIdx + 1,
				Evidences:     []message.Msg{aggVoteC},
			},
			outCome: failureReturn,
		},
		{
			// precommit for nil is not accountable.
			proof: Proof{
				Rule:          autonity.C,
				Message:       preCommitNilC,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{aggVoteC},
			},
			outCome: failureReturn,
		},
		{
			// there is invalid msg in the evidence.
			proof: Proof{
				Rule:          autonity.C,
				Message:       preCommitC,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{aggVoteC, liteOldP, preCommitC},
			},
			outCome: failureReturn,
		},
		{
			// there is duplicated msg in the evidence.
			proof: Proof{
				Rule:          autonity.C,
				Message:       preCommitC,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{aggVoteC, aggVoteC},
			},
			outCome: failureReturn,
		},
		{
			// there is less quorum prevotes in the evidence.
			proof: Proof{
				Rule:          autonity.C,
				Message:       preCommitC,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{aggVoteCNoQuorum},
			},
			outCome: failureReturn,
		},
		{
			proof: Proof{
				Rule:          autonity.C,
				Message:       preCommitC,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{aggVoteC},
			},
			outCome: validReturn(preCommitC, proposer, autonity.C),
		},
	}
	mv := MisbehaviourVerifier{}
	for i, tc := range tests {
		proof := tc.proof
		ret := mv.validateFault(&proof, committee)
		if !bytes.Equal(tc.outCome, ret) {
			t.Log("TestMisbehaviourVerifier", "case", i, "config", tc, "expected", tc.outCome, "actual", ret)
		}
		assert.Equal(t, tc.outCome, ret)
	}
}

func TestInnocenceVerifier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	chainMock := NewMockChainContext(ctrl)
	chainMock.EXPECT().CommitteeOfHeight(height).AnyTimes().Return(committee, nil)

	proposalPO := newValidatedLightProposal(height, 1, 0, signer, committee, nil, proposerIdx)
	votesPO := make([]message.Vote, committee.Len())
	for i := range committee.Members {
		votesPO[i] = newValidatedPrevote(0, height, proposalPO.Value(), makeSigner(keys[i]), &committee.Members[i], cSize)
	}
	aggVotesPO := message.AggregatePrevotes(votesPO)
	aggVotesPONoQuorum := message.AggregatePrevotes(votesPO[2:3])
	votesForOtherValue := newValidatedPrevote(0, height, noneNilValue, signer, self, cSize)

	nilPrevote := newValidatedPrevote(1, height, nilValue, signer, self, cSize)

	// PVN settings
	lightProposalPVN := newValidatedLightProposal(height, 1, -1, signer, committee, nil, proposerIdx)

	// PVO settings
	proposalPVO := newValidatedLightProposal(height, 1, 0, signer, committee, nil, proposerIdx)
	preVotePVO := newValidatedPrevote(1, height, proposalPVO.Value(), signer, self, cSize)
	preVoteNilPVO := newValidatedPrevote(1, height, nilValue, signer, self, cSize)
	// prepare quorum prevotes at valid round.
	votesPVO := make([]message.Vote, committee.Len())
	for i := range committee.Members {
		votesPVO[i] = newValidatedPrevote(0, height, proposalPVO.Value(), makeSigner(keys[i]), &committee.Members[i], cSize)
	}
	aggVotePVO := message.AggregatePrevotes(votesPVO)
	aggVotePVONoQuorum := message.AggregatePrevotes(votesPVO[2:3])

	// C1 settings
	preCommitC1 := newValidatedPrecommit(1, height, noneNilValue, signer, self, cSize)
	preCommitC1Nil := newValidatedPrecommit(1, height, nilValue, signer, self, cSize)
	votesC1 := make([]message.Vote, committee.Len())
	for i := range committee.Members {
		votesC1[i] = newValidatedPrevote(1, height, noneNilValue, makeSigner(keys[i]), &committee.Members[i], cSize)
	}
	aggVoteC1 := message.AggregatePrevotes(votesC1)
	preVoteC1ForOtherV := newValidatedPrevote(1, height, proposalPO.Value(), signer, self, cSize)
	aggVoteC1NoQuorum := message.AggregatePrevotes(votesC1[2:3])

	t.Run("Test innocence verifier required gas", func(t *testing.T) {
		iv := InnocenceVerifier{chain: nil}
		assert.Equal(t, params.AutonityAFDContractGasPerKB, iv.RequiredGas(nil))
	})

	t.Run("Test innocence verifier run with nil bytes", func(t *testing.T) {
		iv := InnocenceVerifier{chain: nil}
		ret, err := iv.Run(nil, height, nil, common.Address{})
		assert.Equal(t, failureReturn, ret)
		assert.Nil(t, err)
	})

	t.Run("Test validate innocence Proof with invalid Signature() of message", func(t *testing.T) {
		invalidCommittee, iKeys, _ := generateCommittee()
		p := &Proof{
			Rule:          autonity.PO,
			OffenderIndex: proposerIdx,
			Message:       newValidatedLightProposal(height, 1, 0, makeSigner(iKeys[0]), invalidCommittee, nil, 0),
		}
		iv := InnocenceVerifier{chain: chainMock}
		raw, err := rlp.EncodeToBytes(&p)
		require.NoError(t, err)
		ret, err := iv.Run(append(make([]byte, 32), raw...), height, nil, common.Address{})
		require.NoError(t, err)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate innocence Proof, with invalid Signature() of evidence msgs", func(t *testing.T) {

		var p Proof
		p.Rule = autonity.PO
		p.OffenderIndex = proposerIdx
		invalidCommittee, iKeys, _ := generateCommittee()
		proposal := newValidatedLightProposal(height, 1, 0, signer, committee, nil, proposerIdx)
		p.Message = proposal
		invalidPreVote := newValidatedPrevote(1, height, proposal.Value(), makeSigner(iKeys[0]),
			&invalidCommittee.Members[0], invalidCommittee.Len())
		p.Evidences = append(p.Evidences, invalidPreVote)

		iv := InnocenceVerifier{chain: chainMock}
		raw, err := rlp.EncodeToBytes(&p)
		require.NoError(t, err)
		ret, err := iv.Run(append(make([]byte, 32), raw...), height, nil, common.Address{})
		require.NoError(t, err)
		assert.Equal(t, failureReturn, ret)
	})

	tests := []testCase{
		// Innocence proof of PO test comes here:
		{
			// wrong msg provided in the proof.
			proof: Proof{
				Rule:          autonity.PO,
				Message:       newValidatedPrevote(1, height, noneNilValue, signer, self, cSize),
				OffenderIndex: proposerIdx,
			},
			outCome: false,
		},
		{
			// wrong proposal provided in the proof.
			proof: Proof{
				Rule:          autonity.PO,
				Message:       newValidatedLightProposal(height, 1, -1, signer, committee, nil, proposerIdx),
				OffenderIndex: proposerIdx,
			},
			outCome: false,
		},
		{
			// have preVote of different value than proposal
			proof: Proof{
				Rule:          autonity.PO,
				Message:       newValidatedLightProposal(height, 1, 0, signer, committee, nil, proposerIdx),
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{newValidatedPrevote(0, height, noneNilValue, signer, self, cSize)},
			},
			outCome: false,
		},
		{
			// have prevotes for other value
			proof: Proof{
				Rule:          autonity.PO,
				Message:       proposalPO,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{votesForOtherValue},
			},
			outCome: false,
		},
		{
			// have no quorum prevotes from proof
			proof: Proof{
				Rule:          autonity.PO,
				Message:       proposalPO,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{aggVotesPONoQuorum},
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:          autonity.PO,
				Message:       proposalPO,
				OffenderIndex: proposerIdx,
				Evidences:     []message.Msg{aggVotesPO},
			},
			outCome: true,
		},
		// Innocence proof of PVN test comes here:
		{
			// wrong msg provided in the proof.
			proof: Proof{
				Rule: autonity.PVN,
			},
			outCome: false,
		},
		{
			// prevote for nil is not accountable.
			proof: Proof{
				Rule:          autonity.PVN,
				OffenderIndex: proposerIdx,
				Message:       nilPrevote,
			},
			outCome: false,
		},
		{
			// no evidence
			proof: Proof{
				Rule:          autonity.PVN,
				OffenderIndex: proposerIdx,
				Message:       aggVotesPO,
			},
			outCome: false,
		},
		{
			// wrong msg in evidence.
			proof: Proof{
				Rule:          autonity.PVN,
				OffenderIndex: proposerIdx,
				Message:       aggVotesPO,
				Evidences:     []message.Msg{aggVotesPO},
			},
			outCome: false,
		},
		{
			// wrong msg in evidence.
			proof: Proof{
				Rule:          autonity.PVN,
				OffenderIndex: proposerIdx,
				Message:       newValidatedPrevote(1, height, lightProposalPVN.Value(), signer, self, cSize),
				Evidences:     []message.Msg{lightProposalPVN},
			},
			outCome: true,
		},
		// Innocence proof of PVO test comes here:
		{
			// wrong message
			proof: Proof{
				Rule:          autonity.PVO,
				OffenderIndex: proposerIdx,
				Message:       lightProposalPVN,
			},
			outCome: false,
		},
		{
			// prevote nil is not accountable
			proof: Proof{
				Rule:          autonity.PVO,
				Evidences:     []message.Msg{proposalPVO, aggVotePVO},
				OffenderIndex: proposerIdx,
				Message:       preVoteNilPVO,
			},
			outCome: false,
		},
		{
			// no evidence
			proof: Proof{
				Rule:          autonity.PVO,
				OffenderIndex: proposerIdx,
				Message:       preVotePVO,
			},
			outCome: false,
		},
		{
			// there is no corresponding proposal in the evidence.
			proof: Proof{
				Rule:          autonity.PVO,
				Evidences:     []message.Msg{aggVotePVO},
				OffenderIndex: proposerIdx,
				Message:       preVotePVO,
			},
			outCome: false,
		},
		{
			// there is votes for other value in the evidence.
			proof: Proof{
				Rule:          autonity.PVO,
				Evidences:     []message.Msg{proposalPVO, aggVotePVO, votesForOtherValue},
				OffenderIndex: proposerIdx,
				Message:       preVotePVO,
			},
			outCome: false,
		},
		{
			// there are duplicated votes in evidence.
			proof: Proof{
				Rule:          autonity.PVO,
				Evidences:     []message.Msg{proposalPVO, aggVotePVO, aggVotePVO},
				OffenderIndex: proposerIdx,
				Message:       preVotePVO,
			},
			outCome: false,
		},
		{
			// there are less quorum votes in evidence.
			proof: Proof{
				Rule:          autonity.PVO,
				Evidences:     []message.Msg{proposalPVO, aggVotePVONoQuorum},
				OffenderIndex: proposerIdx,
				Message:       preVotePVO,
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:          autonity.PVO,
				Evidences:     []message.Msg{proposalPVO, aggVotePVO},
				OffenderIndex: proposerIdx,
				Message:       preVotePVO,
			},
			outCome: true,
		},
		// Innocence proof of C1 test comes here:
		{
			// wrong message type.
			proof: Proof{
				Rule:          autonity.C1,
				Evidences:     []message.Msg{aggVoteC1},
				OffenderIndex: proposerIdx,
				Message:       prevote2,
			},
			outCome: false,
		},
		{
			// precomit for nil is not accountable.
			proof: Proof{
				Rule:          autonity.C1,
				Evidences:     []message.Msg{aggVoteC1},
				OffenderIndex: proposerIdx,
				Message:       preCommitC1Nil,
			},
			outCome: false,
		},
		{
			// votes for other value.
			proof: Proof{
				Rule:          autonity.C1,
				Evidences:     []message.Msg{aggVoteC1NoQuorum, preVoteC1ForOtherV},
				OffenderIndex: proposerIdx,
				Message:       preCommitC1,
			},
			outCome: false,
		},
		{
			// duplicated msg in evidence.
			proof: Proof{
				Rule:          autonity.C1,
				Evidences:     []message.Msg{aggVoteC1, aggVoteC1},
				OffenderIndex: proposerIdx,
				Message:       preCommitC1,
			},
			outCome: false,
		},
		{
			proof: Proof{
				Rule:          autonity.C1,
				Evidences:     []message.Msg{aggVoteC1},
				OffenderIndex: proposerIdx,
				Message:       preCommitC1,
			},
			outCome: true,
		},
	}

	for i, tc := range tests {
		proof := tc.proof
		ret := verifyInnocenceProof(&proof, committee)
		if ret != tc.outCome {
			t.Log("TestInnocenceVerifier", "case", i, "config", tc, "expected", tc.outCome, "actual", ret)
		}
		assert.Equal(t, tc.outCome, ret)
	}
}

func TestCheckEquivocation(t *testing.T) {
	round := int64(0)
	t.Run("check equivocation with valid Proof of proposal equivocation", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.Equivocation
		p.OffenderIndex = proposerIdx
		proposal := newValidatedLightProposal(height, round, -1, signer, committee, nil, proposerIdx)
		p.Message = proposal
		p2 := newValidatedLightProposal(height, round, 1, signer, committee, nil, proposerIdx)
		p.Evidences = append(p.Evidences, p2)
		require.Equal(t, true, validMisbehaviourOfEquivocation(&p, committee))
	})

	t.Run("check equivocation with valid Proof of prevote equivocation", func(t *testing.T) {
		vote1 := newValidatedPrevote(round, height, nilValue, signer, self, cSize)
		vote2 := newValidatedPrevote(round, height, common.Hash{0x1}, signer, self, cSize)
		var p Proof
		p.Rule = autonity.Equivocation
		p.OffenderIndex = proposerIdx
		p.Message = vote1
		p.Evidences = append(p.Evidences, vote2)
		require.Equal(t, true, validMisbehaviourOfEquivocation(&p, committee))
	})

	t.Run("check equivocation with valid Proof of precomit equivocation", func(t *testing.T) {
		vote1 := newValidatedPrecommit(round, height, nilValue, signer, self, cSize)
		vote2 := newValidatedPrecommit(round, height, common.Hash{0x1}, signer, self, cSize)
		var p Proof
		p.Rule = autonity.Equivocation
		p.OffenderIndex = proposerIdx
		p.Message = vote1
		p.Evidences = append(p.Evidences, vote2)
		require.Equal(t, true, validMisbehaviourOfEquivocation(&p, committee))
	})
}

func makeSigner(key blst.SecretKey) message.Signer {
	return func(hash common.Hash) blst.Signature {
		signature := key.Sign(hash[:])
		return signature
	}
}

func newValidatedLightProposal(height uint64, r int64, vr int64, signer message.Signer, committee *types.Committee,
	block *types.Block, idx int) *message.LightProposal { //nolint
	rawProposal := newValidatedProposalMessage(height, r, vr, signer, committee, block, idx)
	return rawProposal.ToLight()
}

func newValidatedPrecommit(r int64, height uint64, v common.Hash, signer message.Signer,
	s *types.CommitteeMember, cSize int) *message.Precommit {
	preCommit := message.NewPrecommit(r, height, v, signer, s, cSize)
	return preCommit
}

func newValidatedPrevote(r int64, height uint64, v common.Hash, signer message.Signer,
	s *types.CommitteeMember, cSize int) *message.Prevote {
	prevote := message.NewPrevote(r, height, v, signer, s, cSize)
	return prevote
}
