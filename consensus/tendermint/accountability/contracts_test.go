package accountability

import (
	"crypto/ecdsa"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/crypto"
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
	proposal := newProposalMessage(height, 3, 0, signer, committee, nil).MustVerify(stubVerifier).ToLight()
	preCommit := message.NewPrecommit(3, height, proposal.Value(), signer)

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
	height := uint64(100)
	proposal := newProposalMessage(height, 3, 0, signer, committee, nil)

	t.Run("Test accusation verifier required gas", func(t *testing.T) {
		av := AccusationVerifier{}
		assert.Equal(t, params.AutonityAFDContractGasPerKB, av.RequiredGas(nil))
	})

	t.Run("Test accusation verifier run with nil bytes", func(t *testing.T) {
		av := AccusationVerifier{}
		ret, err := av.Run(nil, height)
		assert.Equal(t, failureReturn, ret)
		assert.Nil(t, err)
	})

	t.Run("Test accusation verifier run with invalid rlp bytes", func(t *testing.T) {
		wrongBytes := failureReturn
		av := AccusationVerifier{}
		ret, err := av.Run(wrongBytes, height)
		assert.Equal(t, failureReturn, ret)
		assert.Nil(t, err)
	})

	t.Run("Test validate accusation, with wrong rule ID", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.Equivocation + 100
		assert.False(t, verifyAccusation(&p))
	})

	t.Run("Test validate accusation, with wrong accusation msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		preVote := message.NewPrevote(0, height, proposal.Value(), signer)
		p.Message = preVote
		assert.Equal(t, false, verifyAccusation(&p))

		p.Rule = autonity.PVN
		p.Message = proposal
		assert.Equal(t, false, verifyAccusation(&p))

		p.Rule = autonity.C
		p.Message = proposal
		assert.Equal(t, false, verifyAccusation(&p))

		p.Rule = autonity.C1
		p.Message = proposal
		assert.Equal(t, false, verifyAccusation(&p))
	})

	t.Run("Test validate accusation, with invalid Signature() of msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		invalidCommittee, invalKeys := generateCommittee()
		p.Message = newProposalMessage(height, 1, 0, makeSigner(invalKeys[0], invalidCommittee[0]), invalidCommittee, nil)

		ret := verifyAccusation(&p)
		assert.False(t, ret)
	})

	t.Run("Test validate accusation, with correct accusation msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		newProposal := newProposalMessage(height, 1, 0, signer, committee, nil).MustVerify(stubVerifier)
		p.Message = newProposal.ToLight()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ret := verifyAccusation(&p)
		assert.True(t, ret)
		/*
			assert.Equal(t, common.LeftPadBytes(proposer.Bytes(), 32), ret[0:32])
			assert.Equal(t, newProposal.Hash().Bytes(), ret[32:64])
			assert.Equal(t, successResult, ret[64:96])
		*/
	})

	t.Run("Test validate accusation, with PVO accusation msgs", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVO
		oldProposal := newProposalMessage(height, 1, 0, signer, committee, nil).MustVerify(stubVerifier)
		p.Message = message.NewPrevote(1, height, oldProposal.Value(), signer)
		p.Evidences = append(p.Evidences, oldProposal.ToLight())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ret := verifyAccusation(&p)
		assert.True(t, ret)
		/*
			assert.NotEqual(t, failureResult, ret)
			assert.Equal(t, common.LeftPadBytes(proposer.Bytes(), 32), ret[0:32])
			assert.Equal(t, preVote.Hash().Bytes(), ret[32:64])
			assert.Equal(t, successResult, ret[64:96])
		*/
	})

	t.Run("Test validate accusation, with invalid PVO accusation proof", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVO
		oldProposal := newProposalMessage(height, 1, 0, signer, committee, nil)
		preVote := message.NewPrevote(2, height, oldProposal.Value(), signer)
		p.Message = preVote
		p.Evidences = append(p.Evidences, oldProposal)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ret := verifyAccusation(&p)
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
		mv := MisbehaviourVerifier{chain: chainMock}
		assert.Equal(t, params.AutonityAFDContractGasPerKB, mv.RequiredGas(nil))
	})

	t.Run("Test misbehaviour verifier run with nil bytes", func(t *testing.T) {
		mv := MisbehaviourVerifier{chain: chainMock}
		ret, err := mv.Run(nil, height)
		assert.Equal(t, failureReturn, ret)
		assert.Nil(t, err)
	})

	t.Run("Test misbehaviour verifier run with invalid rlp bytes", func(t *testing.T) {
		wrongBytes := failureReturn
		mv := MisbehaviourVerifier{chain: chainMock}
		ret, err := mv.Run(wrongBytes, height)
		assert.Equal(t, failureReturn, ret)
		assert.Nil(t, err)
	})

	t.Run("Test validate misbehaviour Proof, with invalid Signature() of misbehaved msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		invalidCommittee, iKeys := generateCommittee()
		invalidProposal := newProposalMessage(height, 1, 0, makeSigner(iKeys[0], invalidCommittee[0]), invalidCommittee, nil)
		p.Message = invalidProposal
		mv := MisbehaviourVerifier{chain: chainMock}

		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof, with invalid Signature() of evidence msgs", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO

		invalidCommittee, ikeys := generateCommittee()
		proposal := newProposalMessage(height, 1, 0, signer, committee, nil)
		p.Message = proposal

		invalidPreCommit := message.NewPrecommit(1, height, proposal.Value(), makeSigner(ikeys[0], invalidCommittee[0]))
		p.Evidences = append(p.Evidences, invalidPreCommit)

		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of PN rule with correct Proof", func(t *testing.T) {
		// prepare a Proof that node proposes for a new value, but he preCommitted a non nil value
		// at previous rounds, such Proof should be valid.
		var p Proof
		p.Rule = autonity.PN
		proposal := newProposalMessage(height, 1, -1, signer, committee, nil).MustVerify(stubVerifier)
		p.Message = proposal.ToLight()

		preCommit := message.NewPrecommit(0, height, noneNilValue, signer).MustVerify(stubVerifier)

		p.Evidences = append(p.Evidences, preCommit)
		mv := MisbehaviourVerifier{chain: chainMock}

		ret := mv.validateFault(&p)
		assert.Equal(t, validReturn(p.Message, p.Rule), ret)
	})

	t.Run("Test validate misbehaviour Proof of PN rule with incorrect proposal of Proof", func(t *testing.T) {
		// prepare a p that node propose for an old value.
		var p Proof
		p.Rule = autonity.PN
		p.Message = newProposalMessage(height, 1, 0, signer, committee, nil)

		preCommit := message.NewPrecommit(0, height, noneNilValue, signer)
		p.Evidences = append(p.Evidences, preCommit)

		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of PN rule with no evidence of Proof", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PN
		p.Message = newProposalMessage(height, 1, -1, signer, committee, nil)

		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of PO, propose a v rather than the locked one", func(t *testing.T) {
		// simulate a p of misbehaviour of PO, with the proposer proposed an old value that was not
		// the one he locked at previous round, the validation of this p should return true.
		var p Proof
		p.Rule = autonity.PO
		proposal := newProposalMessage(height, 3, 0, signer, committee, nil).MustVerify(stubVerifier)
		preCommit := message.NewPrecommit(0, height, noneNilValue, signer).MustVerify(stubVerifier)
		p.Message = proposal.ToLight()
		p.Evidences = append(p.Evidences, preCommit)
		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, validReturn(p.Message, p.Rule), ret)
	})

	t.Run("Test validate misbehaviour Proof of PO, propose a valid round rather than the locked one", func(t *testing.T) {
		// simulate a p of misbehaviour of PO, with the proposer proposed a valid round that was not
		// the one he locked at previous round, the validation of this p should return true.
		var p Proof
		p.Rule = autonity.PO
		proposal := newProposalMessage(height, 3, 0, signer, committee, nil).MustVerify(stubVerifier)
		preCommit := message.NewPrecommit(1, height, noneNilValue, signer)
		p.Message = proposal.ToLight()
		p.Evidences = append(p.Evidences, preCommit)
		mv := MisbehaviourVerifier{chain: chainMock}
		verifyProofSignatures(chainMock, &p)
		ret := mv.validateFault(&p)
		assert.Equal(t, validReturn(p.Message, p.Rule), ret)
	})

	t.Run("Test validate misbehaviour Proof of PO, propose a different value rather than the one that have quorum "+
		"preVotes at valid round.", func(t *testing.T) {
		// simulate a p of misbehaviour of PO, with the proposer proposed a valid round that was not
		// the one he locked at previous round, the validation of this p should return true.

		var p Proof
		p.Rule = autonity.PO
		p.Message = newProposalMessage(height, 3, 0, signer, committee, nil).MustVerify(stubVerifier).ToLight()
		for i := range committee {
			prevote := message.NewPrevote(0, height, noneNilValue, makeSigner(keys[i], committee[i]))
			p.Evidences = append(p.Evidences, prevote)
		}

		mv := MisbehaviourVerifier{chain: chainMock}
		verifyProofSignatures(chainMock, &p)
		ret := mv.validateFault(&p)
		assert.Equal(t, validReturn(p.Message, p.Rule), ret)
	})

	t.Run("Test validate misbehaviour Proof of PO, with no evidence", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		p.Message = newProposalMessage(height, 3, 0, signer, committee, nil).MustVerify(stubVerifier).ToLight()
		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of PO, with a proposal of new value", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		p.Message = newProposalMessage(height, 3, -1, signer, committee, nil).MustVerify(stubVerifier).ToLight()
		preCommit := message.NewPrecommit(0, height, noneNilValue, signer).MustVerify(stubVerifier)

		p.Evidences = append(p.Evidences, preCommit)
		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of PVN rule, with correct Proof", func(t *testing.T) {
		// simulate a p of misbehaviour of PVN, with the node preVote for V1, but he preCommit
		// at a different value V2 at previous round. The validation of the misbehaviour p should
		// return ture.
		var p Proof
		p.Rule = autonity.PVN
		// node locked at V1 at round 0.
		preCommit := message.NewPrecommit(0, height, noneNilValue, signer)
		preCommitR1 := message.NewPrecommit(1, height, nilValue, signer)
		preCommitR2 := message.NewPrecommit(2, height, nilValue, signer)
		proposal := message.NewLightProposal(newProposalMessage(height, 3, -1, signer, committee, nil).MustVerify(lastHeader.CommitteeMember))
		// node preVote for V2 at round 3
		p.Message = message.NewPrevote(3, height, proposal.Value(), signer)
		p.Evidences = append(p.Evidences, proposal, preCommit, preCommitR1, preCommitR2)

		mv := MisbehaviourVerifier{chain: chainMock}
		verifyProofSignatures(chainMock, &p)
		ret := mv.validateFault(&p)
		assert.Equal(t, validReturn(p.Message, p.Rule), ret)
	})

	t.Run("Test validate misbehaviour Proof of PVN rule, with gaps in preCommits", func(t *testing.T) {
		// simulate a p of misbehaviour of PVN, with the node preVote for V1, but he preCommit
		// at a different value V2 at previous round. The validation of the misbehaviour p should
		// return ture.
		var p Proof
		p.Rule = autonity.PVN
		// node locked at V1 at round 0.
		preCommit := message.NewPrecommit(0, height, noneNilValue, signer)
		preCommitR1 := message.NewPrecommit(1, height, nilValue, signer)
		proposal := newProposalMessage(height, 3, -1, signer, committee, nil).MustVerify(stubVerifier)
		// node preVote for V2 at round 3
		preVote := message.NewPrevote(3, height, proposal.Value(), signer)
		p.Message = preVote
		p.Evidences = append(p.Evidences, preCommit, preCommitR1)

		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of PVN rule, with no evidence", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVN
		proposal := newProposalMessage(height, 3, -1, signer, committee, nil).MustVerify(stubVerifier)
		// node preVote for V2 at round 3
		preVote := message.NewPrevote(3, height, proposal.Value(), signer)
		p.Message = preVote

		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of PVN rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVN
		proposal := newProposalMessage(height, 3, -1, signer, committee, nil).MustVerify(stubVerifier)
		// set a wrong type of msg.
		p.Message = proposal
		preCommit := message.NewPrecommit(0, height, noneNilValue, signer)
		p.Evidences = append(p.Evidences, preCommit)
		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of PVN rule, with wrong preVote value", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVN
		// node locked at V1 at round 0.
		preCommit := message.NewPrecommit(0, height, noneNilValue, signer)
		// node preVote for V2 at round 3, with nil value, not provable.
		preVote := message.NewPrevote(3, height, nilValue, signer)
		p.Message = preVote
		p.Evidences = append(p.Evidences, preCommit)

		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO1 rule, with correct proof", func(t *testing.T) {
		correspondingProposal := newProposalMessage(height, 3, 0, signer, committee, nil).MustVerify(stubVerifier)
		// a precommit at round 1, with value v.
		pcForV := message.NewPrecommit(1, height, correspondingProposal.Value(), signer)
		// a precommit at round 2, with value not v.
		pcForNotV := message.NewPrecommit(2, height, noneNilValue, signer)

		// a prevote at round 3, with value v.
		preVote := message.NewPrevote(3, height, correspondingProposal.Value(), signer)
		var p Proof
		p.Rule = autonity.PVO12
		p.Type = autonity.Misbehaviour
		p.Message = preVote
		p.Evidences = append(p.Evidences, message.NewLightProposal(correspondingProposal), pcForV, pcForNotV)
		mv := MisbehaviourVerifier{chain: chainMock}
		verifyProofSignatures(chainMock, &p)
		ret := mv.validateFault(&p)
		assert.Equal(t, validReturn(p.Message, p.Rule), ret)
	})

	t.Run("Test validate misbehaviour proof of PVO12 rule, with no evidence", func(t *testing.T) {
		correspondingProposal := newProposalMessage(height, 3, 0, signer, committee, nil).MustVerify(stubVerifier)
		// a prevote at round 3, with value v.
		preVote := message.NewPrevote(3, height, correspondingProposal.Value(), signer)
		var p Proof
		p.Rule = autonity.PVO12
		p.Type = autonity.Misbehaviour
		p.Message = preVote
		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO12 rule, with wrong msg", func(t *testing.T) {
		correspondingProposal := newProposalMessage(height, 3, 0, signer, committee, nil).MustVerify(stubVerifier)
		// a precommit at round 1, with value v.
		pcForV := message.NewPrecommit(1, height, correspondingProposal.Value(), signer)
		// a precommit at round 2, with value not v.
		pcForNotV := message.NewPrecommit(2, height, noneNilValue, signer)

		var p Proof
		p.Rule = autonity.PVO12
		p.Type = autonity.Misbehaviour
		p.Message = correspondingProposal
		p.Evidences = append(p.Evidences, correspondingProposal, pcForV, pcForNotV)
		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO12 rule, with in-corresponding proposal", func(t *testing.T) {
		correspondingProposal := newProposalMessage(height, 2, 0, signer, committee, nil).MustVerify(stubVerifier)
		// a precommit at round 1, with value v.
		pcForV := message.NewPrecommit(1, height, correspondingProposal.Value(), signer)
		// a precommit at round 2, with value not v.
		pcForNotV := message.NewPrecommit(2, height, noneNilValue, signer)

		// a prevote at round 3, with value v.
		preVote := message.NewPrevote(3, height, correspondingProposal.Value(), signer)
		var p Proof
		p.Rule = autonity.PVO12
		p.Type = autonity.Misbehaviour
		p.Message = preVote
		p.Evidences = append(p.Evidences, correspondingProposal, pcForV, pcForNotV)
		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO1 rule, with precommits out of round range", func(t *testing.T) {
		correspondingProposal := newProposalMessage(height, 3, 0, signer, committee, nil).MustVerify(stubVerifier)
		// a precommit at round 0, with value v.
		pcValidRound := message.NewPrecommit(0, height, correspondingProposal.Value(), signer)
		// a precommit at round 1, with value v.
		pcForV := message.NewPrecommit(1, height, correspondingProposal.Value(), signer)

		// a precommit at round 4, with value not v.
		pcForNotV := message.NewPrecommit(4, height, noneNilValue, signer)

		// a prevote at round 3, with value v.
		preVote := message.NewPrevote(3, height, correspondingProposal.Value(), signer)
		var p Proof
		p.Rule = autonity.PVO12
		p.Type = autonity.Misbehaviour
		p.Message = preVote
		p.Evidences = append(p.Evidences, correspondingProposal, pcValidRound, pcForV, pcForNotV)
		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO rule, with correct proof", func(t *testing.T) {
		correspondingProposal := newProposalMessage(height, 3, 0, signer, committee, nil).MustVerify(stubVerifier)
		maliciousPreVote := message.NewPrevote(3, height, correspondingProposal.Value(), signer)
		var p Proof
		p.Rule = autonity.PVO
		p.Message = maliciousPreVote
		p.Evidences = append(p.Evidences, message.NewLightProposal(correspondingProposal))
		// simulate quorum prevote for not v at valid round.
		for i := range committee {
			preVote := message.NewPrevote(0, height, noneNilValue, makeSigner(keys[i], committee[i]))
			p.Evidences = append(p.Evidences, preVote)
		}

		mv := MisbehaviourVerifier{chain: chainMock}
		verifyProofSignatures(chainMock, &p)
		ret := mv.validateFault(&p)
		assert.Equal(t, validReturn(p.Message, p.Rule), ret)
	})

	t.Run("Test validate misbehaviour proof of PVO rule, with less quorum preVote for not v", func(t *testing.T) {
		correspondingProposal := newProposalMessage(height, 3, 0, signer, committee, nil).MustVerify(stubVerifier)
		maliciousPreVote := message.NewPrevote(3, height, correspondingProposal.Value(), signer)
		var p Proof
		p.Rule = autonity.PVO
		p.Message = maliciousPreVote
		p.Evidences = append(p.Evidences, correspondingProposal)
		// simulate only one prevote for not v at valid round.
		preVote := message.NewPrevote(0, height, noneNilValue, signer)
		p.Evidences = append(p.Evidences, preVote)

		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO rule, with preVotes at wrong valid round", func(t *testing.T) {
		correspondingProposal := newProposalMessage(height, 3, 0, signer, committee, nil).MustVerify(stubVerifier)
		maliciousPreVote := message.NewPrevote(3, height, correspondingProposal.Value(), signer)
		var p Proof
		p.Rule = autonity.PVO
		p.Message = maliciousPreVote
		p.Evidences = append(p.Evidences, correspondingProposal)
		// simulate quorum prevote for not v at a round rather than valid round
		for i := range committee {
			preVote := message.NewPrevote(1, height, noneNilValue, makeSigner(keys[i], committee[i]))
			p.Evidences = append(p.Evidences, preVote)
		}

		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO2 rule, with precommits of V", func(t *testing.T) {
		correspondingProposal := newProposalMessage(height, 3, 0, signer, committee, nil).MustVerify(stubVerifier)
		// a precommit at round 0, with value not v.
		pcVR := message.NewPrecommit(0, height, correspondingProposal.Value(), signer)
		// a precommit at round 1, with value not v.
		pcR1 := message.NewPrecommit(1, height, noneNilValue, signer)
		// a precommit at round 2, with value not v.
		pcR2 := message.NewPrecommit(2, height, noneNilValue, signer)

		// a prevote at round 3, with value v.
		preVote := message.NewPrevote(3, height, correspondingProposal.Value(), signer)
		var p Proof
		p.Rule = autonity.PVO12
		p.Type = autonity.Misbehaviour
		p.Message = preVote
		p.Evidences = append(p.Evidences, correspondingProposal, pcVR, pcR1, pcR2)
		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with correct Proof", func(t *testing.T) {
		// Node preCommit for a V at round R, but in that round, there were quorum PreVotes for notV at that round.
		var p Proof
		p.Rule = autonity.C
		// Node preCommit for V at round R.
		preCommit := message.NewPrecommit(0, height, noneNilValue, signer)
		p.Message = preCommit
		for i := range committee {
			preVote := message.NewPrevote(0, height, common.Hash{0x2}, makeSigner(keys[i], committee[i]))
			p.Evidences = append(p.Evidences, preVote)
		}
		mv := MisbehaviourVerifier{chain: chainMock}
		verifyProofSignatures(chainMock, &p)
		ret := mv.validateFault(&p)
		assert.Equal(t, validReturn(p.Message, p.Rule), ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with no Evidences", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C
		// Node preCommit for V at round R.
		preCommit := message.NewPrecommit(0, height, noneNilValue, signer)
		p.Message = preCommit

		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with wrong preCommit msg", func(t *testing.T) {
		// Node preCommit for nil at round R, not provable
		var p Proof
		p.Rule = autonity.C
		preCommit := message.NewPrecommit(0, height, nilValue, signer)
		p.Message = preCommit
		for i := range committee {
			preVote := message.NewPrevote(0, height, noneNilValue, makeSigner(keys[i], committee[i]))
			p.Evidences = append(p.Evidences, preVote)
		}

		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C

		wrongMsg := message.NewPrevote(0, height, noneNilValue, signer)
		p.Message = wrongMsg
		for i := range committee {
			preVote := message.NewPrevote(0, height, nilValue, makeSigner(keys[i], committee[i]))
			p.Evidences = append(p.Evidences, preVote)
		}

		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with invalid evidence", func(t *testing.T) {
		// the evidence contains same value of preCommit that node preVoted for.
		var p Proof
		p.Rule = autonity.C

		preCommit := message.NewPrecommit(0, height, noneNilValue, signer)
		p.Message = preCommit
		// quorum preVotes of same value, this shouldn't be a valid evidence.
		for i := range committee {
			preVote := message.NewPrevote(0, height, noneNilValue, makeSigner(keys[i], committee[i]))
			p.Evidences = append(p.Evidences, preVote)
		}

		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with invalid evidence: duplicated msg in evidence", func(t *testing.T) {
		// the evidence contains same value of preCommit that node preVoted for.
		var p Proof
		p.Rule = autonity.C

		preCommit := message.NewPrecommit(0, height, noneNilValue, signer)
		p.Message = preCommit
		// duplicated preVotes msg in evidence, should be addressed.
		for i := 0; i < len(committee); i++ {
			preVote := message.NewPrevote(0, height, nilValue, signer)
			p.Evidences = append(p.Evidences, preVote)
		}

		mv := MisbehaviourVerifier{chain: chainMock}
		verifyProofSignatures(chainMock, &p)
		ret := mv.validateFault(&p)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with invalid evidence: no quorum preVotes", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C

		preCommit := message.NewPrecommit(0, height, noneNilValue, signer)
		p.Message = preCommit

		// no quorum preVotes msg in evidence, should be addressed.
		preVote := message.NewPrevote(0, height, common.Hash{0x2}, signer)
		p.Evidences = append(p.Evidences, preVote)

		mv := MisbehaviourVerifier{chain: chainMock}
		verifyProofSignatures(chainMock, &p)
		ret := mv.validateFault(&p)
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
		ret, err := iv.Run(nil, height)
		assert.Equal(t, failureReturn, ret)
		assert.Nil(t, err)
	})

	t.Run("Test validate innocence Proof with invalid Signature() of message", func(t *testing.T) {
		invalidCommittee, iKeys := generateCommittee()
		p := &Proof{
			Rule:    autonity.PO,
			Message: message.NewLightProposal(newProposalMessage(height, 1, 0, makeSigner(iKeys[0], invalidCommittee[0]), invalidCommittee, nil).MustVerify(stubVerifier)),
		}
		iv := InnocenceVerifier{chain: chainMock}
		raw, err := rlp.EncodeToBytes(&p)
		require.NoError(t, err)
		ret, err := iv.Run(append(make([]byte, 32), raw...), height)
		require.NoError(t, err)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate innocence Proof, with invalid Signature() of evidence msgs", func(t *testing.T) {

		var p Proof
		p.Rule = autonity.PO
		invalidCommittee, iKeys := generateCommittee()
		proposal := newProposalMessage(height, 1, 0, signer, committee, nil).MustVerify(stubVerifier)
		p.Message = message.NewLightProposal(proposal)
		invalidPreVote := message.NewPrevote(1, height, proposal.Value(), makeSigner(iKeys[0], invalidCommittee[0]))
		p.Evidences = append(p.Evidences, invalidPreVote)

		iv := InnocenceVerifier{chain: chainMock}
		raw, err := rlp.EncodeToBytes(&p)
		require.NoError(t, err)
		ret, err := iv.Run(append(make([]byte, 32), raw...), height)
		require.NoError(t, err)
		assert.Equal(t, failureReturn, ret)
	})

	t.Run("Test validate innocence Proof of PO rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		wrongMsg := message.NewPrevote(1, height, noneNilValue, signer)
		p.Message = wrongMsg

		ret := validInnocenceProofOfPO(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence AccountabilityProof of PO rule, with invalid evidence", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		proposal := newProposalMessage(height, 1, 0, signer, committee, nil).MustVerify(stubVerifier)
		p.Message = proposal
		// have preVote at different value than proposal
		invalidPreVote := message.NewPrevote(0, height, noneNilValue, signer)
		p.Evidences = append(p.Evidences, invalidPreVote)
		ret := validInnocenceProofOfPO(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence AccountabilityProof of PO rule, with redundant vote msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		proposal := newProposalMessage(height, 1, 0, signer, committee, nil).MustVerify(stubVerifier)
		p.Message = proposal

		preVote := message.NewPrevote(0, height, proposal.Value(), signer)
		p.Evidences = append(p.Evidences, preVote)
		// make redundant msg hack.
		p.Evidences = append(p.Evidences, p.Evidences...)

		ret := validInnocenceProofOfPO(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence AccountabilityProof of PO rule, with not quorum vote msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		proposal := newProposalMessage(height, 1, 0, signer, committee, nil).MustVerify(stubVerifier)
		p.Message = proposal

		preVote := message.NewPrevote(0, height, proposal.Value(), signer)
		p.Evidences = append(p.Evidences, preVote)

		ret := validInnocenceProofOfPO(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of PVN rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVN
		wrongMsg := message.NewPrecommit(1, height, noneNilValue, signer)
		p.Message = wrongMsg
		ret := validInnocenceProofOfPVN(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of PVN rule, with a wrong preVote for nil", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVN
		wrongMsg := message.NewPrevote(1, height, nilValue, signer)
		p.Message = wrongMsg
		ret := validInnocenceProofOfPVN(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of PVN rule, with no evidence", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVN
		wrongMsg := message.NewPrevote(1, height, noneNilValue, signer)
		p.Message = wrongMsg
		ret := validInnocenceProofOfPVN(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of PVN rule, with over quorum prevotes", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVN
		proposal := newProposalMessage(height, 1, -1, signer, committee, nil).MustVerify(stubVerifier)
		p.Evidences = append(p.Evidences, message.NewLightProposal(proposal))
		preVote := message.NewPrevote(1, height, proposal.Value(), signer)
		p.Message = preVote
		ret := validInnocenceProofOfPVN(&p)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate innocence proof of PVO rule, with correct Proof", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVO
		vr := int64(0)
		proposal := newProposalMessage(height, 1, vr, signer, committee, nil).MustVerify(stubVerifier)
		preVote := message.NewPrevote(1, height, proposal.Value(), signer)
		p.Message = preVote
		p.Evidences = append(p.Evidences, message.NewLightProposal(proposal))
		// prepare quorum prevotes at valid round.
		for i := range committee {
			prevote := message.NewPrevote(vr, height, proposal.Value(), makeSigner(keys[i], committee[i]))
			p.Evidences = append(p.Evidences, prevote)
		}

		validateProof(&p, lastHeader)
		ret := validInnocenceProofOfPVO(&p, chainMock)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate innocence proof of PVO rule, with incorrect proposal", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVO
		vr := int64(0)
		// with wrong round in proposal.
		proposal := newProposalMessage(height, 2, vr, signer, committee, nil).MustVerify(stubVerifier)
		preVote := message.NewPrevote(1, height, proposal.Value(), signer)
		p.Message = preVote
		p.Evidences = append(p.Evidences, proposal)
		// prepare quorum prevotes at valid round.
		for i := range committee {
			preVote := message.NewPrevote(vr, height, proposal.Value(), makeSigner(keys[i], committee[i]))
			p.Evidences = append(p.Evidences, preVote)
		}
		validateProof(&p, lastHeader)
		ret := validInnocenceProofOfPVO(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence proof of PVO rule, with incorrect preVotes", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVO
		vr := int64(0)
		// with wrong round in proposal.
		proposal := newProposalMessage(height, 2, vr, signer, committee, nil).MustVerify(stubVerifier)
		preVote := message.NewPrevote(1, height, proposal.Value(), signer)
		p.Message = preVote
		p.Evidences = append(p.Evidences, proposal)
		// prepare quorum prevotes at wrong round.
		for i := range committee {
			preVote := message.NewPrevote(1, height, proposal.Value(), makeSigner(keys[i], committee[i]))
			p.Evidences = append(p.Evidences, preVote)
		}
		validateProof(&p, lastHeader)
		ret := validInnocenceProofOfPVO(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence proof of PVO rule, with less than quorum preVotes", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVO
		vr := int64(0)
		proposal := newProposalMessage(height, 2, vr, signer, committee, nil).MustVerify(stubVerifier)
		preVote := message.NewPrevote(2, height, proposal.Value(), signer)
		p.Message = preVote
		p.Evidences = append(p.Evidences, proposal)
		// prepare only one prevotes at valid round.
		v := message.NewPrevote(vr, height, proposal.Value(), signer)
		p.Evidences = append(p.Evidences, v)

		validateProof(&p, lastHeader)
		ret := validInnocenceProofOfPVO(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence proof of PVO rule, with preVote for not V", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVO
		vr := int64(0)
		proposal := newProposalMessage(height, 2, vr, signer, committee, nil).MustVerify(stubVerifier)
		preVote := message.NewPrevote(2, height, proposal.Value(), signer)
		p.Message = preVote
		p.Evidences = append(p.Evidences, proposal)
		// prepare only one prevotes at valid round.
		for i := range committee {
			preVote := message.NewPrevote(vr, height, noneNilValue, makeSigner(keys[i], committee[i]))
			p.Evidences = append(p.Evidences, preVote)
		}
		validateProof(&p, lastHeader)
		ret := validInnocenceProofOfPVO(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C1
		wrongMsg := message.NewPrevote(1, height, noneNilValue, signer)
		p.Message = wrongMsg

		validateProof(&p, lastHeader)
		ret := validInnocenceProofOfC1(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with a wrong preCommit for nil", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C1
		wrongMsg := message.NewPrecommit(1, height, nilValue, signer)
		p.Message = wrongMsg

		validateProof(&p, lastHeader)
		ret := validInnocenceProofOfC1(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with a wrong evidence", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C1
		preCommit := message.NewPrecommit(1, height, noneNilValue, signer)
		p.Message = preCommit
		// evidence contains a preVote of a different round
		preVote := message.NewPrevote(0, height, noneNilValue, signer)
		p.Evidences = append(p.Evidences, preVote)

		validateProof(&p, lastHeader)
		ret := validInnocenceProofOfC1(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with redundant msgs in evidence ", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C1
		preCommit := message.NewPrecommit(1, height, noneNilValue, signer)
		p.Message = preCommit

		preVote := message.NewPrevote(1, height, noneNilValue, signer)
		p.Evidences = append(p.Evidences, preVote)
		p.Evidences = append(p.Evidences, p.Evidences...)

		validateProof(&p, lastHeader)
		ret := validInnocenceProofOfC1(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with no quorum votes of evidence ", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C1
		preCommit := message.NewPrecommit(1, height, noneNilValue, signer)
		p.Message = preCommit

		preVote := message.NewPrevote(1, height, noneNilValue, signer)
		p.Evidences = append(p.Evidences, preVote)

		validateProof(&p, lastHeader)
		ret := validInnocenceProofOfC1(&p, chainMock)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with correct evidence ", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C1
		preCommit := message.NewPrecommit(1, height, noneNilValue, signer)
		p.Message = preCommit
		for i := range committee {
			preVote := message.NewPrevote(1, height, noneNilValue, makeSigner(keys[i], committee[i]))
			p.Evidences = append(p.Evidences, preVote)
		}

		validateProof(&p, lastHeader)
		ret := validInnocenceProofOfC1(&p, chainMock)
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
		proposal := newProposalMessage(height, round, -1, signer, committee, nil)
		require.Nil(t, verifyProofSignatures(chainMock, &Proof{Message: proposal}))
	})

	t.Run("a future msg is received, expect an error of errFutureMsg", func(t *testing.T) {
		futureHeight := height + 1
		proposal := newProposalMessage(futureHeight, round, -1, signer, committee, nil)
		chainMock.EXPECT().GetHeaderByNumber(height).Return(nil)
		require.Equal(t, errFutureMsg, verifyProofSignatures(chainMock, &Proof{Message: proposal}))
	})

	t.Run("chain cannot provide the last header of the height that msg votes on, expect an error of errFutureMsg", func(t *testing.T) {
		proposal := newProposalMessage(height-5, round, -1, signer, committee, nil)
		chainMock.EXPECT().GetHeaderByNumber(height - 6).Return(nil)
		require.Equal(t, errFutureMsg, verifyProofSignatures(chainMock, &Proof{Message: proposal}))
	})

	t.Run("abnormal case, msg is not signed by committee", func(t *testing.T) {
		wrongCommitte, ks := generateCommittee()
		proposal := newProposalMessage(height, round, -1, makeSigner(ks[0], wrongCommitte[0]), wrongCommitte, nil)
		require.Equal(t, errNotCommitteeMsg, verifyProofSignatures(chainMock, &Proof{Message: proposal}))
	})
}

func TestCheckEquivocation(t *testing.T) {
	height := uint64(100)
	round := int64(0)
	t.Run("check equivocation with valid Proof of equivocation", func(t *testing.T) {
		proposal := newProposalMessage(height, round, -1, signer, committee, nil).MustVerify(stubVerifier)
		vote1 := message.NewPrevote(round, height, proposal.Value(), signer)
		vote2 := message.NewPrevote(round, height, nilValue, signer)
		var proofs []message.Msg
		proofs = append(proofs, vote2)
		require.Equal(t, errEquivocation, checkEquivocation(vote1, proofs))
	})

	t.Run("check equivocation with invalid Proof of equivocation", func(t *testing.T) {
		proposal := newProposalMessage(height, round, -1, signer, committee, nil).MustVerify(stubVerifier)
		vote1 := message.NewPrevote(round, height, proposal.Value(), signer)
		var proofs []message.Msg
		proofs = append(proofs, vote1)
		require.Nil(t, checkEquivocation(vote1, proofs))
	})
}

func validateProof(p *Proof, header *types.Header) {
	p.Message.Validate(header.CommitteeMember)
	for _, m := range p.Evidences {
		m.Validate(header.CommitteeMember)
	}
}

func makeSigner(key *ecdsa.PrivateKey, val types.CommitteeMember) message.Signer {
	return func(hash common.Hash) ([]byte, common.Address) {
		out, _ := crypto.Sign(hash[:], key)
		return out, val.Address
	}
}

func stubVerifier(address common.Address) *types.CommitteeMember {
	return &types.CommitteeMember{
		Address:     address,
		VotingPower: common.Big1,
	}
}
