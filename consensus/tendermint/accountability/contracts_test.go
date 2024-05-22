package accountability

import (
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

func TestDecodeProof(t *testing.T) {
	height := uint64(100)
	header := newBlockHeader(height-1, committee)
	rawProposal := newValidatedProposalMessage(height, 3, 0, signer, committee, nil, proposerIdx)
	err := rawProposal.PreValidate(header)
	require.NoError(t, err)
	err = rawProposal.Validate()
	require.NoError(t, err)
	proposal := rawProposal.ToLight()
	preCommit := message.NewPrecommit(3, height, proposal.Value(), signer, self, cSize)

	t.Run("decode with accusation", func(t *testing.T) {
		var p Proof
		p.Type = autonity.Accusation
		p.Rule = autonity.PO
		p.Message = proposal

		rp, err := rlp.EncodeToBytes(&p)
		assert.NoError(t, err)

		decodeProof, err := decodeRawProof(rp)
		assert.NoError(t, err)
		assert.Equal(t, autonity.PO, decodeProof.Rule)
		assert.Equal(t, proposal.Signature(), decodeProof.Message.Signature())
	})

	t.Run("decode with evidence", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		p.Message = proposal
		p.Evidences = append(p.Evidences, preCommit)

		rp, err := rlp.EncodeToBytes(&p)
		assert.NoError(t, err)

		decodeProof, err := decodeRawProof(rp)
		assert.NoError(t, err)
		assert.Equal(t, autonity.PO, decodeProof.Rule)
		assert.Equal(t, proposal.Signature(), decodeProof.Message.Signature())
		assert.Equal(t, preCommit.Signature(), decodeProof.Evidences[0].Signature())
	})
}

func TestAccusationVerifier(t *testing.T) {
	// Todo(youssef): add integration tests for the precompile Run function
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	height := uint64(100)
	lastHeight := height - 1
	lastHeader := newBlockHeader(lastHeight, committee)
	proposal := newValidatedProposalMessage(height, 3, 0, signer, committee, nil, proposerIdx).ToLight()
	chainMock := NewMockChainContext(ctrl)
	chainMock.EXPECT().GetHeaderByNumber(lastHeight).AnyTimes().Return(lastHeader)

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

	t.Run("Test validate accusation, with wrong rule ID", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.Equivocation + 100
		assert.False(t, verifyAccusation(&p, committee))
	})

	t.Run("Test validate accusation, with wrong accusation msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		preVote := message.NewPrevote(0, height, proposal.Value(), signer, self, cSize)
		p.Message = preVote
		assert.Equal(t, false, verifyAccusation(&p, committee))

		p.Rule = autonity.PVN
		p.Message = proposal
		assert.Equal(t, false, verifyAccusation(&p, committee))

		p.Rule = autonity.C
		p.Message = proposal
		assert.Equal(t, false, verifyAccusation(&p, committee))

		p.Rule = autonity.C1
		p.Message = proposal
		assert.Equal(t, false, verifyAccusation(&p, committee))
	})

	t.Run("Test validate accusation, with invalid Signature() of msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		invalidCommittee, invalKeys, _ := generateCommittee()
		p.Message = newValidatedProposalMessage(height, 1, 0, makeSigner(invalKeys[0]), invalidCommittee, nil, 0).ToLight()
		c, err := verifyProofSignatures(chainMock, &p)
		require.Nil(t, c)
		require.NotNil(t, err)
	})

	t.Run("Test validate accusation, with correct accusation msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		p.OffenderIndex = proposerIdx
		p.Offender = proposer
		liteP := newValidatedLightProposal(t, height, 1, 0, signer, committee, lastHeader, nil, proposerIdx)
		p.Message = liteP
		ret := verifyAccusation(&p, committee)
		assert.True(t, ret)
	})

	t.Run("Test validate accusation, with PVO accusation msgs", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVO

		liteP := newValidatedLightProposal(t, height, 1, 0, signer, committee, lastHeader, nil, proposerIdx)

		p.Message = message.NewPrevote(1, height, liteP.Value(), signer, self, cSize)
		p.Evidences = append(p.Evidences, liteP)
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		ret := verifyAccusation(&p, committee)
		assert.True(t, ret)
	})

	t.Run("Test validate accusation, with invalid PVO accusation proof", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVO
		p.Offender = proposer
		p.OffenderIndex = proposerIdx

		liteP := newValidatedLightProposal(t, height, 1, 0, signer, committee, lastHeader, nil, proposerIdx)

		preVote := message.NewPrevote(2, height, liteP.Value(), signer, self, cSize)
		p.Message = preVote
		p.Evidences = append(p.Evidences, liteP)
		ret := verifyAccusation(&p, committee)
		assert.False(t, ret)
	})
}

func TestMisbehaviourVerifier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	height := uint64(100)
	lastHeight := height - 1
	lastHeader := newBlockHeader(lastHeight, committee)
	noneNilValue := common.Hash{0x1}

	chainMock := NewMockChainContext(ctrl)
	chainMock.EXPECT().GetHeaderByNumber(lastHeight).AnyTimes().Return(lastHeader)

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

	t.Run("Test validate misbehaviour Proof, with invalid Signature() of misbehaved msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		invalidCommittee, iKeys, _ := generateCommittee()
		invalidProposal := newValidatedProposalMessage(height, 1, 0, makeSigner(iKeys[0]), invalidCommittee, nil, 0)
		p.Message = invalidProposal.ToLight()

		c, err := verifyProofSignatures(chainMock, &p)
		require.Nil(t, c)
		require.NotNil(t, err)
	})

	t.Run("Test validate misbehaviour Proof, with invalid Signature() of evidence msgs", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		invalidCommittee, ikeys, _ := generateCommittee()
		proposal := newValidatedProposalMessage(height, 1, 0, signer, committee, nil, proposerIdx)
		p.Message = proposal.ToLight()

		invalidPreCommit := message.NewPrecommit(1, height, proposal.Value(), makeSigner(ikeys[0]), &invalidCommittee[0], len(invalidCommittee))
		p.Evidences = append(p.Evidences, invalidPreCommit)
		rawProof, err := rlp.EncodeToBytes(&p)
		require.NoError(t, err)

		decodedProof, err := decodeRawProof(rawProof)
		require.NoError(t, err)

		c, err := verifyProofSignatures(chainMock, decodedProof)
		require.Nil(t, c)
		require.Equal(t, "bad signature", err.Error())
	})

	t.Run("Test validate misbehaviour Proof of PN rule with correct Proof", func(t *testing.T) {
		// prepare a Proof that node proposes for a new value, but he preCommitted a non nil value
		// at previous rounds, such Proof should be valid.
		var p Proof
		p.Rule = autonity.PN

		liteP := newValidatedLightProposal(t, height, 1, -1, signer, committee, lastHeader, nil, proposerIdx)
		p.Message = liteP
		preCommit := newValidatedPrecommit(t, 0, height, noneNilValue, signer, self, cSize, lastHeader)
		p.OffenderIndex = proposerIdx
		p.Offender = proposer

		p.Evidences = append(p.Evidences, preCommit)
		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, validReturn(p.Message, proposer, p.Rule), ret)
	})

	t.Run("Test validate misbehaviour Proof of PN rule with incorrect proposal of Proof", func(t *testing.T) {
		// prepare a p that node propose for an old value.
		var p Proof
		p.Rule = autonity.PN
		p.Offender = proposer
		p.OffenderIndex = proposerIdx

		liteP := newValidatedLightProposal(t, height, 1, 0, signer, committee, lastHeader, nil, proposerIdx)
		p.Message = liteP
		preCommit := newValidatedPrecommit(t, 0, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Evidences = append(p.Evidences, preCommit)

		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of PN rule with no evidence of Proof", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PN
		p.OffenderIndex = proposerIdx
		p.Offender = proposer
		p.Message = newValidatedLightProposal(t, height, 1, -1, signer, committee, lastHeader, nil, proposerIdx)
		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of PO, propose a v rather than the locked one", func(t *testing.T) {
		// simulate a p of misbehaviour of PO, with the proposer proposed an old value that was not
		// the one he locked at previous round, the validation of this p should return true.
		var p Proof
		p.Rule = autonity.PO
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		preCommit := newValidatedPrecommit(t, 0, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Message = newValidatedLightProposal(t, height, 3, 0, signer, committee, lastHeader, nil, proposerIdx)
		p.Evidences = append(p.Evidences, preCommit)
		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, validReturn(p.Message, proposer, p.Rule), ret)
	})

	t.Run("Test validate misbehaviour Proof of PO, propose a valid round rather than the locked one", func(t *testing.T) {
		// simulate a p of misbehaviour of PO, with the proposer proposed a valid round that was not
		// the one he locked at previous round, the validation of this p should return true.
		var p Proof
		p.Rule = autonity.PO
		p.Offender = proposer
		p.OffenderIndex = proposerIdx

		preCommit := newValidatedPrecommit(t, 1, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Message = newValidatedLightProposal(t, height, 3, 0, signer, committee, lastHeader, nil, proposerIdx)
		p.Evidences = append(p.Evidences, preCommit)
		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, validReturn(p.Message, proposer, p.Rule), ret)
	})

	t.Run("Test validate misbehaviour Proof of PO, propose a different value rather than the one that have quorum "+
		"preVotes at valid round.", func(t *testing.T) {
		// simulate a p of misbehaviour of PO, with the proposer proposed a valid round that was not
		// the one he locked at previous round, the validation of this p should return true.

		var p Proof
		p.Rule = autonity.PO
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		p.Message = newValidatedLightProposal(t, height, 3, 0, signer, committee, lastHeader, nil, proposerIdx)
		for i := range committee {
			prevote := newValidatedPrevote(t, 0, height, noneNilValue, makeSigner(keys[i]), &committee[i], cSize, lastHeader)
			p.Evidences = append(p.Evidences, prevote)
		}

		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, validReturn(p.Message, proposer, p.Rule), ret)
	})

	t.Run("Test validate misbehaviour Proof of PO, with no evidence", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		p.Message = newValidatedLightProposal(t, height, 3, 0, signer, committee, lastHeader, nil, proposerIdx)
		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of PO, with a proposal of new value", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		p.Message = newValidatedLightProposal(t, height, 3, -1, signer, committee, lastHeader, nil, proposerIdx)
		preCommit := newValidatedPrecommit(t, 0, height, noneNilValue, signer, self, cSize, lastHeader)

		p.Evidences = append(p.Evidences, preCommit)
		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of PVN rule, with correct Proof", func(t *testing.T) {
		// simulate a p of misbehaviour of PVN, with the node preVote for V1, but he preCommit
		// at a different value V2 at previous round. The validation of the misbehaviour p should
		// return ture.
		var p Proof
		p.Rule = autonity.PVN
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		// node locked at V1 at round 0.
		preCommit := newValidatedPrecommit(t, 0, height, noneNilValue, signer, self, cSize, lastHeader)
		preCommitR1 := newValidatedPrecommit(t, 1, height, nilValue, signer, self, cSize, lastHeader)
		preCommitR2 := newValidatedPrecommit(t, 2, height, nilValue, signer, self, cSize, lastHeader)

		proposal := newValidatedLightProposal(t, height, 3, -1, signer, committee, lastHeader, nil, proposerIdx)
		// node preVote for V2 at round 3
		p.Message = newValidatedPrevote(t, 3, height, proposal.Value(), signer, self, cSize, lastHeader)
		p.Evidences = append(p.Evidences, proposal, preCommit, preCommitR1, preCommitR2)

		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, validReturn(p.Message, proposer, p.Rule), ret)
	})

	t.Run("Test validate misbehaviour Proof of PVN rule, with gaps in preCommits", func(t *testing.T) {
		// simulate a p of misbehaviour of PVN, with the node preVote for V1, but he preCommit
		// at a different value V2 at previous round. The validation of the misbehaviour p should
		// return ture.
		var p Proof
		p.Rule = autonity.PVN
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		// node locked at V1 at round 0.
		preCommit := newValidatedPrecommit(t, 0, height, noneNilValue, signer, self, cSize, lastHeader)
		preCommitR1 := newValidatedPrecommit(t, 1, height, nilValue, signer, self, cSize, lastHeader)
		proposal := newValidatedLightProposal(t, height, 3, -1, signer, committee, lastHeader, nil, proposerIdx)
		// node preVote for V2 at round 3
		preVote := newValidatedPrevote(t, 3, height, proposal.Value(), signer, self, cSize, lastHeader)
		p.Message = preVote
		p.Evidences = append(p.Evidences, preCommit, preCommitR1)

		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of PVN rule, with no evidence", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVN
		p.Offender = proposer
		p.OffenderIndex = proposerIdx

		// node preVote for V2 at round 3
		p.Message = newValidatedPrevote(t, 3, height, noneNilValue, signer, self, cSize, lastHeader)

		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of PVN rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVN
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		// set a wrong type of msg.
		p.Message = newValidatedLightProposal(t, height, 3, -1, signer, committee, lastHeader, nil, proposerIdx)
		preCommit := newValidatedPrecommit(t, 0, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Evidences = append(p.Evidences, preCommit)
		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of PVN rule, with wrong preVote value", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVN
		// node locked at V1 at round 0.
		preCommit := newValidatedPrecommit(t, 0, height, noneNilValue, signer, self, cSize, lastHeader)
		// node preVote for V2 at round 3, with nil value, not provable.
		preVote := newValidatedPrevote(t, 3, height, nilValue, signer, self, cSize, lastHeader)
		p.Message = preVote
		p.Evidences = append(p.Evidences, preCommit)

		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO1 rule, with correct proof", func(t *testing.T) {
		correspondingProposal := newValidatedLightProposal(t, height, 3, 0, signer, committee, lastHeader, nil, proposerIdx)
		// a precommit at round 1, with value v.
		pcForV := newValidatedPrecommit(t, 1, height, correspondingProposal.Value(), signer, self, cSize, lastHeader)
		// a precommit at round 2, with value not v.
		pcForNotV := newValidatedPrecommit(t, 2, height, noneNilValue, signer, self, cSize, lastHeader)

		// a prevote at round 3, with value v.
		preVote := newValidatedPrevote(t, 3, height, correspondingProposal.Value(), signer, self, cSize, lastHeader)
		var p Proof
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		p.Rule = autonity.PVO12
		p.Type = autonity.Misbehaviour
		p.Message = preVote
		p.Evidences = append(p.Evidences, correspondingProposal, pcForV, pcForNotV)
		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, validReturn(p.Message, p.Offender, p.Rule), ret)
	})

	t.Run("Test validate misbehaviour proof of PVO12 rule, with no evidence", func(t *testing.T) {

		correspondingProposal := newValidatedLightProposal(t, height, 3, 0, signer, committee, lastHeader, nil, proposerIdx)
		// a prevote at round 3, with value v.
		preVote := newValidatedPrevote(t, 3, height, correspondingProposal.Value(), signer, self, cSize, lastHeader)
		var p Proof
		p.Rule = autonity.PVO12
		p.Type = autonity.Misbehaviour
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		p.Message = preVote
		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO12 rule, with wrong msg", func(t *testing.T) {
		correspondingProposal := newValidatedLightProposal(t, height, 3, 0, signer, committee, lastHeader, nil, proposerIdx)
		// a precommit at round 1, with value v.
		pcForV := newValidatedPrecommit(t, 1, height, correspondingProposal.Value(), signer, self, cSize, lastHeader)
		// a precommit at round 2, with value not v.
		pcForNotV := newValidatedPrecommit(t, 2, height, noneNilValue, signer, self, cSize, lastHeader)

		var p Proof
		p.Rule = autonity.PVO12
		p.Type = autonity.Misbehaviour
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		p.Message = correspondingProposal
		p.Evidences = append(p.Evidences, correspondingProposal, pcForV, pcForNotV)
		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO12 rule, with in-corresponding proposal", func(t *testing.T) {
		correspondingProposal := newValidatedLightProposal(t, height, 2, 0, signer, committee, lastHeader, nil, proposerIdx)
		// a precommit at round 1, with value v.
		pcForV := newValidatedPrecommit(t, 1, height, correspondingProposal.Value(), signer, self, cSize, lastHeader)
		// a precommit at round 2, with value not v.
		pcForNotV := newValidatedPrecommit(t, 2, height, noneNilValue, signer, self, cSize, lastHeader)

		// a prevote at round 3, with value v.
		preVote := newValidatedPrevote(t, 3, height, correspondingProposal.Value(), signer, self, cSize, lastHeader)
		var p Proof
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		p.Rule = autonity.PVO12
		p.Type = autonity.Misbehaviour
		p.Message = preVote
		p.Evidences = append(p.Evidences, correspondingProposal, pcForV, pcForNotV)
		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO1 rule, with precommits out of round range", func(t *testing.T) {
		correspondingProposal := newValidatedLightProposal(t, height, 3, 0, signer, committee, lastHeader, nil, proposerIdx)
		// a precommit at round 0, with value v.
		pcValidRound := newValidatedPrecommit(t, 0, height, correspondingProposal.Value(), signer, self, cSize, lastHeader)
		// a precommit at round 1, with value v.
		pcForV := newValidatedPrecommit(t, 1, height, correspondingProposal.Value(), signer, self, cSize, lastHeader)

		// a precommit at round 4, with value not v.
		pcForNotV := newValidatedPrecommit(t, 4, height, noneNilValue, signer, self, cSize, lastHeader)

		// a prevote at round 3, with value v.
		preVote := newValidatedPrevote(t, 3, height, correspondingProposal.Value(), signer, self, cSize, lastHeader)
		var p Proof
		p.Rule = autonity.PVO12
		p.Type = autonity.Misbehaviour
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		p.Message = preVote
		p.Evidences = append(p.Evidences, correspondingProposal, pcValidRound, pcForV, pcForNotV)
		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO rule, with correct proof", func(t *testing.T) {
		correspondingProposal := newValidatedLightProposal(t, height, 3, 0, signer, committee, lastHeader, nil, proposerIdx)
		maliciousPreVote := newValidatedPrevote(t, 3, height, correspondingProposal.Value(), signer, self, cSize, lastHeader)
		var p Proof
		p.OffenderIndex = proposerIdx
		p.Offender = proposer
		p.Rule = autonity.PVO
		p.Message = maliciousPreVote
		p.Evidences = append(p.Evidences, correspondingProposal)
		// simulate quorum prevote for not v at valid round.
		for i := range committee {
			preVote := newValidatedPrevote(t, 0, height, noneNilValue, makeSigner(keys[i]), &committee[i], cSize, lastHeader)
			p.Evidences = append(p.Evidences, preVote)
		}

		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, validReturn(p.Message, p.Offender, p.Rule), ret)
	})

	t.Run("Test validate misbehaviour proof of PVO rule, with less quorum preVote for not v", func(t *testing.T) {
		correspondingProposal := newValidatedLightProposal(t, height, 3, 0, signer, committee, lastHeader, nil, proposerIdx)
		maliciousPreVote := newValidatedPrevote(t, 3, height, correspondingProposal.Value(), signer, self, cSize, lastHeader)
		var p Proof
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		p.Rule = autonity.PVO
		p.Message = maliciousPreVote
		p.Evidences = append(p.Evidences, correspondingProposal)
		// simulate only one prevote for not v at valid round.
		preVote := newValidatedPrevote(t, 0, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Evidences = append(p.Evidences, preVote)

		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO rule, with preVotes at wrong valid round", func(t *testing.T) {
		correspondingProposal := newValidatedLightProposal(t, height, 3, 0, signer, committee, lastHeader, nil, proposerIdx)
		maliciousPreVote := newValidatedPrevote(t, 3, height, correspondingProposal.Value(), signer, self, cSize, lastHeader)
		var p Proof
		p.Rule = autonity.PVO
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		p.Message = maliciousPreVote
		p.Evidences = append(p.Evidences, correspondingProposal)
		// simulate quorum prevote for not v at a round rather than valid round
		for i := range committee {
			preVote := newValidatedPrevote(t, 1, height, noneNilValue, makeSigner(keys[i]), &committee[i], cSize, lastHeader)
			p.Evidences = append(p.Evidences, preVote)
		}

		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO2 rule, with precommits of V", func(t *testing.T) {
		correspondingProposal := newValidatedLightProposal(t, height, 3, 0, signer, committee, lastHeader, nil, proposerIdx)
		// a precommit at round 0, with value not v.
		pcVR := newValidatedPrecommit(t, 0, height, correspondingProposal.Value(), signer, self, cSize, lastHeader)
		// a precommit at round 1, with value not v.
		pcR1 := newValidatedPrecommit(t, 1, height, noneNilValue, signer, self, cSize, lastHeader)
		// a precommit at round 2, with value not v.
		pcR2 := newValidatedPrecommit(t, 2, height, noneNilValue, signer, self, cSize, lastHeader)

		// a prevote at round 3, with value v.
		preVote := newValidatedPrevote(t, 3, height, correspondingProposal.Value(), signer, self, cSize, lastHeader)
		var p Proof
		p.Rule = autonity.PVO12
		p.Type = autonity.Misbehaviour
		p.Message = preVote
		p.OffenderIndex = proposerIdx
		p.Offender = proposer
		p.Evidences = append(p.Evidences, correspondingProposal, pcVR, pcR1, pcR2)
		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with correct Proof", func(t *testing.T) {
		// Node preCommit for a V at round R, but in that round, there were quorum PreVotes for notV at that round.
		var p Proof
		p.Rule = autonity.C
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		// Node preCommit for V at round R.
		preCommit := newValidatedPrecommit(t, 0, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Message = preCommit
		for i := range committee {
			preVote := newValidatedPrevote(t, 0, height, common.Hash{0x2}, makeSigner(keys[i]), &committee[i], cSize, lastHeader)
			p.Evidences = append(p.Evidences, preVote)
		}
		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, validReturn(p.Message, p.Offender, p.Rule), ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with no Evidences", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C
		p.OffenderIndex = proposerIdx
		p.Offender = proposer
		// Node preCommit for V at round R.
		preCommit := newValidatedPrecommit(t, 0, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Message = preCommit

		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with wrong preCommit msg", func(t *testing.T) {
		// Node preCommit for nil at round R, not provable
		var p Proof
		p.Rule = autonity.C
		p.OffenderIndex = proposerIdx
		p.Offender = proposer
		preCommit := newValidatedPrecommit(t, 0, height, nilValue, signer, self, cSize, lastHeader)
		p.Message = preCommit
		for i := range committee {
			preVote := newValidatedPrevote(t, 0, height, noneNilValue, makeSigner(keys[i]), &committee[i], cSize, lastHeader)
			p.Evidences = append(p.Evidences, preVote)
		}

		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C
		p.OffenderIndex = proposerIdx
		p.Offender = proposer
		wrongMsg := newValidatedPrevote(t, 0, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Message = wrongMsg
		for i := range committee {
			preVote := newValidatedPrevote(t, 0, height, nilValue, makeSigner(keys[i]), &committee[i], cSize, lastHeader)
			p.Evidences = append(p.Evidences, preVote)
		}

		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with invalid evidence", func(t *testing.T) {
		// the evidence contains same value of preCommit that node preVoted for.
		var p Proof
		p.Rule = autonity.C
		p.OffenderIndex = proposerIdx
		p.Offender = proposer
		preCommit := newValidatedPrecommit(t, 0, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Message = preCommit
		// quorum preVotes of same value, this shouldn't be a valid evidence.
		for i := range committee {
			preVote := newValidatedPrevote(t, 0, height, noneNilValue, makeSigner(keys[i]), &committee[i], cSize, lastHeader)
			p.Evidences = append(p.Evidences, preVote)
		}

		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with invalid evidence: duplicated msg in evidence", func(t *testing.T) {
		// the evidence contains same value of preCommit that node preVoted for.
		var p Proof
		p.Rule = autonity.C
		p.OffenderIndex = proposerIdx
		p.Offender = proposer
		preCommit := newValidatedPrecommit(t, 0, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Message = preCommit
		// duplicated preVotes msg in evidence, should be addressed.
		for i := 0; i < len(committee); i++ {
			preVote := newValidatedPrevote(t, 0, height, nilValue, signer, self, cSize, lastHeader)
			p.Evidences = append(p.Evidences, preVote)
		}

		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with invalid evidence: no quorum preVotes", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C
		p.OffenderIndex = proposerIdx
		p.Offender = proposer
		preCommit := newValidatedPrecommit(t, 0, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Message = preCommit

		// no quorum preVotes msg in evidence, should be addressed.
		preVote := newValidatedPrevote(t, 0, height, common.Hash{0x2}, signer, self, cSize, lastHeader)
		p.Evidences = append(p.Evidences, preVote)

		mv := MisbehaviourVerifier{}
		ret := mv.validateFault(&p, committee)
		assert.Equal(t, failureReturn, ret)
	})
}

func TestInnocenceVerifier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	height := uint64(100)
	lastHeight := height - 1
	noneNilValue := common.Hash{0x1}
	lastHeader := newBlockHeader(lastHeight, committee)
	chainMock := NewMockChainContext(ctrl)
	chainMock.EXPECT().GetHeaderByNumber(lastHeight).AnyTimes().Return(lastHeader)

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
		lHeader := newBlockHeader(lastHeight, invalidCommittee)
		p := &Proof{
			Rule:          autonity.PO,
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Message:       newValidatedLightProposal(t, height, 1, 0, makeSigner(iKeys[0]), invalidCommittee, lHeader, nil, 0),
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
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		invalidCommittee, iKeys, _ := generateCommittee()
		proposal := newValidatedLightProposal(t, height, 1, 0, signer, committee, lastHeader, nil, proposerIdx)
		p.Message = proposal
		iHeader := newBlockHeader(height, invalidCommittee)
		invalidPreVote := newValidatedPrevote(t, 1, height, proposal.Value(), makeSigner(iKeys[0]),
			&invalidCommittee[0], len(invalidCommittee), iHeader)
		p.Evidences = append(p.Evidences, invalidPreVote)

		iv := InnocenceVerifier{chain: chainMock}
		raw, err := rlp.EncodeToBytes(&p)
		require.NoError(t, err)
		ret, err := iv.Run(append(make([]byte, 32), raw...), height, nil, common.Address{})
		require.NoError(t, err)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate innocence Proof of PO rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		wrongMsg := newValidatedPrevote(t, 1, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Message = wrongMsg

		ret := validInnocenceProofOfPO(&p, committee)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence AccountabilityProof of PO rule, with invalid evidence", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		proposal := newValidatedLightProposal(t, height, 1, 0, signer, committee, lastHeader, nil, proposerIdx)
		p.Message = proposal
		// have preVote at different value than proposal
		invalidPrevote := newValidatedPrevote(t, 0, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Evidences = append(p.Evidences, invalidPrevote)
		ret := validInnocenceProofOfPO(&p, committee)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence AccountabilityProof of PO rule, with redundant vote msg", func(t *testing.T) {
		var p Proof
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		p.Rule = autonity.PO
		proposal := newValidatedLightProposal(t, height, 1, 0, signer, committee, lastHeader, nil, proposerIdx)
		p.Message = proposal

		preVote := newValidatedPrevote(t, 0, height, proposal.Value(), signer, self, cSize, lastHeader)
		p.Evidences = append(p.Evidences, preVote)
		// make redundant msg hack.
		p.Evidences = append(p.Evidences, p.Evidences...)

		ret := validInnocenceProofOfPO(&p, committee)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence AccountabilityProof of PO rule, with not quorum vote msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		proposal := newValidatedLightProposal(t, height, 1, 0, signer, committee, lastHeader, nil, proposerIdx)
		p.Message = proposal

		preVote := newValidatedPrevote(t, 0, height, proposal.Value(), signer, self, cSize, lastHeader)
		p.Evidences = append(p.Evidences, preVote)

		ret := validInnocenceProofOfPO(&p, committee)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of PVN rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVN
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		wrongMsg := newValidatedPrecommit(t, 1, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Message = wrongMsg
		ret := validInnocenceProofOfPVN(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of PVN rule, with a wrong preVote for nil", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVN
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		wrongMsg := newValidatedPrevote(t, 1, height, nilValue, signer, self, cSize, lastHeader)
		p.Message = wrongMsg
		ret := validInnocenceProofOfPVN(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of PVN rule, with no evidence", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVN
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		wrongMsg := newValidatedPrevote(t, 1, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Message = wrongMsg
		ret := validInnocenceProofOfPVN(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of PVN rule, with over quorum prevotes", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVN
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		proposal := newValidatedLightProposal(t, height, 1, -1, signer, committee, lastHeader, nil, proposerIdx)
		p.Evidences = append(p.Evidences, proposal)
		preVote := newValidatedPrevote(t, 1, height, proposal.Value(), signer, self, cSize, lastHeader)
		p.Message = preVote
		ret := validInnocenceProofOfPVN(&p)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate innocence proof of PVO rule, with correct Proof", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVO
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		vr := int64(0)
		proposal := newValidatedLightProposal(t, height, 1, vr, signer, committee, lastHeader, nil, proposerIdx)
		preVote := newValidatedPrevote(t, 1, height, proposal.Value(), signer, self, cSize, lastHeader)
		p.Message = preVote
		p.Evidences = append(p.Evidences, proposal)
		// prepare quorum prevotes at valid round.
		for i := range committee {
			prevote := newValidatedPrevote(t, vr, height, proposal.Value(), makeSigner(keys[i]), &committee[i], cSize, lastHeader)
			p.Evidences = append(p.Evidences, prevote)
		}

		ret := validInnocenceProofOfPVO(&p, committee)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate innocence proof of PVO rule, with incorrect proposal", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVO
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		vr := int64(0)
		// with wrong round in proposal.
		proposal := newValidatedLightProposal(t, height, 2, vr, signer, committee, lastHeader, nil, proposerIdx)
		preVote := newValidatedPrevote(t, 1, height, proposal.Value(), signer, self, cSize, lastHeader)
		p.Message = preVote
		p.Evidences = append(p.Evidences, proposal)
		// prepare quorum prevotes at valid round.
		for i := range committee {
			v := newValidatedPrevote(t, vr, height, proposal.Value(), makeSigner(keys[i]), &committee[i], cSize, lastHeader)
			p.Evidences = append(p.Evidences, v)
		}
		ret := validInnocenceProofOfPVO(&p, committee)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence proof of PVO rule, with incorrect preVotes", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVO
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		vr := int64(0)
		// with wrong round in proposal.
		proposal := newValidatedLightProposal(t, height, 2, vr, signer, committee, lastHeader, nil, proposerIdx)
		preVote := newValidatedPrevote(t, 1, height, proposal.Value(), signer, self, cSize, lastHeader)
		p.Message = preVote
		p.Evidences = append(p.Evidences, proposal)
		// prepare quorum prevotes at wrong round.
		for i := range committee {
			v := newValidatedPrevote(t, 1, height, proposal.Value(), makeSigner(keys[i]), &committee[i], cSize, lastHeader)
			p.Evidences = append(p.Evidences, v)
		}
		ret := validInnocenceProofOfPVO(&p, committee)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence proof of PVO rule, with less than quorum preVotes", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVO
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		vr := int64(0)
		proposal := newValidatedLightProposal(t, height, 2, vr, signer, committee, lastHeader, nil, proposerIdx)
		preVote := newValidatedPrevote(t, 2, height, proposal.Value(), signer, self, cSize, lastHeader)
		p.Message = preVote
		p.Evidences = append(p.Evidences, proposal)
		// prepare only one prevotes at valid round.
		v := newValidatedPrevote(t, vr, height, proposal.Value(), signer, self, cSize, lastHeader)
		p.Evidences = append(p.Evidences, v)

		ret := validInnocenceProofOfPVO(&p, committee)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence proof of PVO rule, with preVote for not V", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVO
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		vr := int64(0)
		proposal := newValidatedLightProposal(t, height, 2, vr, signer, committee, lastHeader, nil, proposerIdx)
		preVote := newValidatedPrevote(t, 2, height, proposal.Value(), signer, self, cSize, lastHeader)
		p.Message = preVote
		p.Evidences = append(p.Evidences, proposal)
		// prepare only one prevotes at valid round.
		for i := range committee {
			v := newValidatedPrevote(t, vr, height, noneNilValue, makeSigner(keys[i]), &committee[i], cSize, lastHeader)
			p.Evidences = append(p.Evidences, v)
		}
		ret := validInnocenceProofOfPVO(&p, committee)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C1
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		wrongMsg := newValidatedPrevote(t, 1, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Message = wrongMsg

		ret := validInnocenceProofOfC1(&p, committee)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with a wrong preCommit for nil", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C1
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		wrongMsg := newValidatedPrecommit(t, 1, height, nilValue, signer, self, cSize, lastHeader)
		p.Message = wrongMsg

		ret := validInnocenceProofOfC1(&p, committee)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with a wrong evidence", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C1
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		preCommit := newValidatedPrecommit(t, 1, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Message = preCommit
		// evidence contains a preVote of a different round
		preVote := newValidatedPrevote(t, 0, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Evidences = append(p.Evidences, preVote)

		ret := validInnocenceProofOfC1(&p, committee)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with redundant msgs in evidence ", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C1
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		preCommit := newValidatedPrecommit(t, 1, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Message = preCommit

		preVote := newValidatedPrevote(t, 1, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Evidences = append(p.Evidences, preVote)
		p.Evidences = append(p.Evidences, p.Evidences...)

		ret := validInnocenceProofOfC1(&p, committee)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with no quorum votes of evidence ", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C1
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		preCommit := newValidatedPrecommit(t, 1, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Message = preCommit

		preVote := newValidatedPrevote(t, 1, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Evidences = append(p.Evidences, preVote)

		ret := validInnocenceProofOfC1(&p, committee)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with correct evidence ", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C1
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		preCommit := newValidatedPrecommit(t, 1, height, noneNilValue, signer, self, cSize, lastHeader)
		p.Message = preCommit
		for i := range committee {
			preVote := newValidatedPrevote(t, 1, height, noneNilValue, makeSigner(keys[i]), &committee[i], cSize, lastHeader)
			p.Evidences = append(p.Evidences, preVote)
		}

		ret := validInnocenceProofOfC1(&p, committee)
		assert.Equal(t, true, ret)
	})
}

// More tests are required here
func TestVerifyProofSignatures(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	height := uint64(100)
	lastHeight := height - 1
	round := int64(0)
	chainMock := NewMockChainContext(ctrl)
	currentHeader := newBlockHeader(lastHeight, committee)
	chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(currentHeader).AnyTimes()

	t.Run("normal case, proposal msg is checked correctly", func(t *testing.T) {
		proposal := newValidatedProposalMessage(height, round, -1, signer, committee, nil, proposerIdx)
		_, err := verifyProofSignatures(chainMock, &Proof{Message: proposal})
		require.Nil(t, err)
	})

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
	})

	t.Run("abnormal case, msg is not signed by committee", func(t *testing.T) {
		wrongCommitte, ks, _ := generateCommittee()
		proposal := newValidatedProposalMessage(height, round, -1, makeSigner(ks[0]), wrongCommitte, nil, proposerIdx)
		_, err := verifyProofSignatures(chainMock, &Proof{Message: proposal})
		require.Equal(t, message.ErrUnauthorizedAddress, err)
	})
}

func TestCheckEquivocation(t *testing.T) {
	height := uint64(100)
	round := int64(0)
	lastHeader := newBlockHeader(height-1, committee)
	t.Run("check equivocation with valid Proof of proposal equivocation", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.Equivocation
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		proposal := newValidatedLightProposal(t, height, round, -1, signer, committee, lastHeader, nil, proposerIdx)
		p.Message = proposal
		p2 := newValidatedLightProposal(t, height, round, 1, signer, committee, lastHeader, nil, proposerIdx)
		p.Evidences = append(p.Evidences, p2)
		require.Equal(t, true, validMisbehaviourOfEquivocation(&p, committee))
	})

	t.Run("check equivocation with valid Proof of prevote equivocation", func(t *testing.T) {
		vote1 := newValidatedPrevote(t, round, height, nilValue, signer, self, cSize, lastHeader)
		vote2 := newValidatedPrevote(t, round, height, common.Hash{0x1}, signer, self, cSize, lastHeader)
		var p Proof
		p.Rule = autonity.Equivocation
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		p.Message = vote1
		p.Evidences = append(p.Evidences, vote2)
		require.Equal(t, true, validMisbehaviourOfEquivocation(&p, committee))
	})

	t.Run("check equivocation with valid Proof of precomit equivocation", func(t *testing.T) {
		vote1 := newValidatedPrecommit(t, round, height, nilValue, signer, self, cSize, lastHeader)
		vote2 := newValidatedPrecommit(t, round, height, common.Hash{0x1}, signer, self, cSize, lastHeader)
		var p Proof
		p.Rule = autonity.Equivocation
		p.Offender = proposer
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

func stubVerifier(consensusKey blst.PublicKey) func(address common.Address) *types.CommitteeMember {
	return func(address common.Address) *types.CommitteeMember {
		return &types.CommitteeMember{
			Address:      address,
			VotingPower:  common.Big1,
			ConsensusKey: consensusKey,
		}
	}
}

func newValidatedLightProposal(t *testing.T, height uint64, r int64, vr int64, signer message.Signer, committee types.Committee,
	lastHeader *types.Header, block *types.Block, idx int) *message.LightProposal {
	rawProposal := newValidatedProposalMessage(height, r, vr, signer, committee, block, idx)
	err := rawProposal.PreValidate(lastHeader)
	require.NoError(t, err)
	err = rawProposal.Validate()
	require.NoError(t, err)
	return rawProposal.ToLight()
}

func newValidatedPrecommit(t *testing.T, r int64, height uint64, v common.Hash, signer message.Signer,
	s *types.CommitteeMember, cSize int, lastHeader *types.Header) *message.Precommit {
	preCommit := message.NewPrecommit(r, height, v, signer, s, cSize)
	err := preCommit.PreValidate(lastHeader)
	require.NoError(t, err)
	err = preCommit.Validate()
	require.NoError(t, err)
	return preCommit
}

func newValidatedPrevote(t *testing.T, r int64, height uint64, v common.Hash, signer message.Signer,
	s *types.CommitteeMember, cSize int, lastHeader *types.Header) *message.Prevote {
	prevote := message.NewPrevote(r, height, v, signer, s, cSize)
	err := prevote.PreValidate(lastHeader)
	require.NoError(t, err)
	err = prevote.Validate()
	require.NoError(t, err)
	return prevote
}
