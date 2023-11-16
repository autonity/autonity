package accountability

import (
	"crypto/ecdsa"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
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
	committee, keys := generateCommittee()
	proposer := committee[0].Address
	proposerKey := keys[proposer]
	proposal := message.NewLightProposal(newProposalMessage(height, 3, 0, proposerKey, committee, nil))
	preCommit := message.NewPrecommit(3, height, proposal.Value(), makeSigner(proposerKey))

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
	lastHeight := height - 1
	committee, keys := generateCommittee()
	proposer := committee[0].Address
	proposerKey := keys[proposer]
	proposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)

	t.Run("Test accusation verifier required gas", func(t *testing.T) {
		av := AccusationVerifier{}
		assert.Equal(t, params.AutonityAFDContractGasPerKB, av.RequiredGas(nil))
	})

	t.Run("Test accusation verifier run with nil bytes", func(t *testing.T) {
		av := AccusationVerifier{}
		ret, err := av.Run(nil, height)
		assert.Equal(t, failureResult, ret)
		assert.Nil(t, err)
	})

	t.Run("Test accusation verifier run with invalid rlp bytes", func(t *testing.T) {
		wrongBytes := failureResult
		av := AccusationVerifier{}
		ret, err := av.Run(wrongBytes, height)
		assert.Equal(t, failureResult, ret)
		assert.Nil(t, err)
	})

	t.Run("Test validate accusation, with wrong rule ID", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.InvalidRound + 100
		assert.False(t, verifyAccusation(nil, &p))
	})

	t.Run("Test validate accusation, with wrong accusation msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		preVote := message.NewPrevote(0, height, proposal.Value(), makeSigner(proposerKey))
		p.Message = preVote
		assert.Equal(t, false, verifyAccusation(nil, &p))

		p.Rule = autonity.PVN
		p.Message = proposal
		assert.Equal(t, false, verifyAccusation(nil, &p))

		p.Rule = autonity.C
		p.Message = proposal
		assert.Equal(t, false, verifyAccusation(nil, &p))

		p.Rule = autonity.C1
		p.Message = proposal
		assert.Equal(t, false, verifyAccusation(nil, &p))
	})

	t.Run("Test validate accusation, with invalid Signature() of msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		invalidCommittee, keys := generateCommittee()
		newProposal := newProposalMessage(height, 1, 0, keys[invalidCommittee[0].Address], invalidCommittee, nil)
		p.Message = newProposal
		ret := verifyAccusation(nil, &p)
		assert.False(t, ret)
	})

	t.Run("Test validate accusation, with correct accusation msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		newProposal := newProposalMessage(height, 1, 0, proposerKey, committee, nil)
		p.Message = message.NewLightProposal(newProposal)
		lastHeader := newBlockHeader(lastHeight, committee)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)

		ret := verifyAccusation(chainMock, &p)
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
		oldProposal := newProposalMessage(height, 1, 0, proposerKey, committee, nil)
		preVote := message.NewPrevote(1, height, oldProposal.Value(), makeSigner(proposerKey))
		p.Message = preVote
		p.Evidences = append(p.Evidences, message.NewLightProposal(oldProposal))
		lastHeader := newBlockHeader(lastHeight, committee)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		ret := verifyAccusation(chainMock, &p)
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
		oldProposal := newProposalMessage(height, 1, 0, proposerKey, committee, nil)
		preVote := message.NewPrevote(2, height, oldProposal.Value(), makeSigner(proposerKey))
		p.Message = preVote
		p.Evidences = append(p.Evidences, oldProposal)
		lastHeader := newBlockHeader(lastHeight, committee)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)

		ret := verifyAccusation(chainMock, &p)
		assert.False(t, ret)
	})
}

func TestMisbehaviourVerifier(t *testing.T) {
	height := uint64(100)
	lastHeight := height - 1
	committee, keys := generateCommittee()
	proposer := committee[0].Address
	proposerKey := keys[proposer]
	noneNilValue := common.Hash{0x1}

	t.Run("Test misbehaviour verifier required gas", func(t *testing.T) {
		mv := MisbehaviourVerifier{}
		assert.Equal(t, params.AutonityAFDContractGasPerKB, mv.RequiredGas(nil))
	})

	t.Run("Test misbehaviour verifier run with nil bytes", func(t *testing.T) {
		mv := MisbehaviourVerifier{}
		ret, err := mv.Run(nil, height)
		assert.Equal(t, failureResult, ret)
		assert.Nil(t, err)
	})

	t.Run("Test misbehaviour verifier run with invalid rlp bytes", func(t *testing.T) {
		wrongBytes := failureResult
		mv := MisbehaviourVerifier{}
		ret, err := mv.Run(wrongBytes, height)
		assert.Equal(t, failureResult, ret)
		assert.Nil(t, err)
	})

	t.Run("Test validate misbehaviour Proof, with invalid Signature() of misbehaved msg", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = autonity.PO
		invalidCommittee, iKeys := generateCommittee()
		invalidProposal := newProposalMessage(height, 1, 0, iKeys[invalidCommittee[0].Address], invalidCommittee, nil)
		p.Message = invalidProposal

		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		mv := MisbehaviourVerifier{chain: chainMock}

		ret := mv.validateProof(&p)
		assert.Equal(t, failureResult, ret)
	})

	t.Run("Test validate misbehaviour Proof, with invalid Signature() of evidence msgs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = autonity.PO
		invalidCommittee, ikeys := generateCommittee()
		proposal := newProposalMessage(height, 1, 0, proposerKey, committee, nil)
		p.Message = proposal
		invalidPreCommit := message.NewPrecommit(1, height, proposal.Value(), makeSigner(ikeys[invalidCommittee[0].Address]))
		p.Evidences = append(p.Evidences, invalidPreCommit)

		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).AnyTimes().Return(lastHeader)
		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateProof(&p)
		assert.Equal(t, failureResult, ret)
	})

	t.Run("Test validate misbehaviour Proof of PN rule with correct Proof", func(t *testing.T) {
		// prepare a Proof that node proposes for a new value, but he preCommitted a non nil value
		// at previous rounds, such Proof should be valid.
		var p Proof
		p.Rule = autonity.PN
		proposal := newProposalMessage(height, 1, -1, proposerKey, committee, nil)
		p.Message = message.NewLightProposal(proposal)

		preCommit := message.NewPrecommit(0, height, noneNilValue, makeSigner(proposerKey))
		p.Evidences = append(p.Evidences, preCommit)
		mv := MisbehaviourVerifier{}

		ret := mv.validProof(&p)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate misbehaviour Proof of PN rule with incorrect proposal of Proof", func(t *testing.T) {
		// prepare a p that node propose for an old value.
		var p Proof
		p.Rule = autonity.PN
		proposal := newProposalMessage(height, 1, 0, proposerKey, committee, nil)
		p.Message = proposal

		preCommit := message.NewPrecommit(0, height, noneNilValue, makeSigner(proposerKey))
		p.Evidences = append(p.Evidences, preCommit)
		mv := MisbehaviourVerifier{}

		ret := mv.validProof(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of PN rule with no evidence of Proof", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PN
		proposal := newProposalMessage(height, 1, -1, proposerKey, committee, nil)
		p.Message = proposal

		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of PO, propose a v rather than the locked one", func(t *testing.T) {
		// simulate a p of misbehaviour of PO, with the proposer proposed an old value that was not
		// the one he locked at previous round, the validation of this p should return true.
		var p Proof
		p.Rule = autonity.PO
		proposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)
		preCommit := message.NewPrecommit(0, height, noneNilValue, makeSigner(proposerKey))
		p.Message = message.NewLightProposal(proposal)
		p.Evidences = append(p.Evidences, preCommit)
		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate misbehaviour Proof of PO, propose a valid round rather than the locked one", func(t *testing.T) {
		// simulate a p of misbehaviour of PO, with the proposer proposed a valid round that was not
		// the one he locked at previous round, the validation of this p should return true.
		var p Proof
		p.Rule = autonity.PO
		proposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)
		preCommit := message.NewPrecommit(1, height, noneNilValue, makeSigner(proposerKey))
		p.Message = message.NewLightProposal(proposal)
		p.Evidences = append(p.Evidences, preCommit)
		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate misbehaviour Proof of PO, propose a different value rather than the one that have quorum "+
		"preVotes at valid round.", func(t *testing.T) {
		// simulate a p of misbehaviour of PO, with the proposer proposed a valid round that was not
		// the one he locked at previous round, the validation of this p should return true.
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = autonity.PO
		proposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)
		p.Message = message.NewLightProposal(proposal)
		for _, c := range committee {
			preVotes := message.NewPrevote(0, height, noneNilValue, makeSigner(keys[c.Address]))
			p.Evidences = append(p.Evidences, preVotes)
		}
		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validProof(&p)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate misbehaviour Proof of PO, with no evidence", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		proposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)
		p.Message = proposal
		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of PO, with a proposal of new value", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		proposal := newProposalMessage(height, 3, -1, proposerKey, committee, nil)
		preCommit := message.NewPrecommit(0, height, noneNilValue, makeSigner(proposerKey))
		p.Message = proposal
		p.Evidences = append(p.Evidences, preCommit)
		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of PVN rule, with correct Proof", func(t *testing.T) {
		// simulate a p of misbehaviour of PVN, with the node preVote for V1, but he preCommit
		// at a different value V2 at previous round. The validation of the misbehaviour p should
		// return ture.
		var p Proof
		p.Rule = autonity.PVN
		// node locked at V1 at round 0.
		preCommit := message.NewPrecommit(0, height, noneNilValue, makeSigner(proposerKey))
		preCommitR1 := message.NewPrecommit(1, height, nilValue, makeSigner(proposerKey))
		preCommitR2 := message.NewPrecommit(2, height, nilValue, makeSigner(proposerKey))
		proposal := newProposalMessage(height, 3, -1, proposerKey, committee, nil)
		// node preVote for V2 at round 3
		preVote := message.NewPrevote(3, height, proposal.Value(), makeSigner(proposerKey))
		p.Message = preVote
		p.Evidences = append(p.Evidences, message.NewLightProposal(proposal), preCommit, preCommitR1, preCommitR2)

		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate misbehaviour Proof of PVN rule, with gaps in preCommits", func(t *testing.T) {
		// simulate a p of misbehaviour of PVN, with the node preVote for V1, but he preCommit
		// at a different value V2 at previous round. The validation of the misbehaviour p should
		// return ture.
		var p Proof
		p.Rule = autonity.PVN
		// node locked at V1 at round 0.
		preCommit := message.NewPrecommit(0, height, noneNilValue, makeSigner(proposerKey))
		preCommitR1 := message.NewPrecommit(1, height, nilValue, makeSigner(proposerKey))
		proposal := newProposalMessage(height, 3, -1, proposerKey, committee, nil)
		// node preVote for V2 at round 3
		preVote := message.NewPrevote(3, height, proposal.Value(), makeSigner(proposerKey))
		p.Message = preVote
		p.Evidences = append(p.Evidences, preCommit, preCommitR1)

		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of PVN rule, with no evidence", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVN
		proposal := newProposalMessage(height, 3, -1, proposerKey, committee, nil)
		// node preVote for V2 at round 3
		preVote := message.NewPrevote(3, height, proposal.Value(), makeSigner(proposerKey))
		p.Message = preVote

		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of PVN rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVN
		proposal := newProposalMessage(height, 3, -1, proposerKey, committee, nil)
		// set a wrong type of msg.
		p.Message = proposal
		preCommit := message.NewPrecommit(0, height, noneNilValue, makeSigner(proposerKey))
		p.Evidences = append(p.Evidences, preCommit)
		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of PVN rule, with wrong preVote value", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVN
		// node locked at V1 at round 0.
		preCommit := message.NewPrecommit(0, height, noneNilValue, makeSigner(proposerKey))
		// node preVote for V2 at round 3, with nil value, not provable.
		preVote := message.NewPrevote(3, height, nilValue, makeSigner(proposerKey))
		p.Message = preVote
		p.Evidences = append(p.Evidences, preCommit)

		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO1 rule, with correct proof", func(t *testing.T) {
		correspondingProposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)
		// a precommit at round 1, with value v.
		pcForV := message.NewPrecommit(1, height, correspondingProposal.Value(), makeSigner(proposerKey))
		// a precommit at round 2, with value not v.
		pcForNotV := message.NewPrecommit(2, height, noneNilValue, makeSigner(proposerKey))

		// a prevote at round 3, with value v.
		preVote := message.NewPrevote(3, height, correspondingProposal.Value(), makeSigner(proposerKey))
		var p Proof
		p.Rule = autonity.PVO12
		p.Type = autonity.Misbehaviour
		p.Message = preVote
		p.Evidences = append(p.Evidences, message.NewLightProposal(correspondingProposal), pcForV, pcForNotV)
		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO12 rule, with no evidence", func(t *testing.T) {
		correspondingProposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)
		// a prevote at round 3, with value v.
		preVote := message.NewPrevote(3, height, correspondingProposal.Value(), makeSigner(proposerKey))
		var p Proof
		p.Rule = autonity.PVO12
		p.Type = autonity.Misbehaviour
		p.Message = preVote
		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO12 rule, with wrong msg", func(t *testing.T) {
		correspondingProposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)
		// a precommit at round 1, with value v.
		pcForV := message.NewPrecommit(1, height, correspondingProposal.Value(), makeSigner(proposerKey))
		// a precommit at round 2, with value not v.
		pcForNotV := message.NewPrecommit(2, height, noneNilValue, makeSigner(proposerKey))

		var p Proof
		p.Rule = autonity.PVO12
		p.Type = autonity.Misbehaviour
		p.Message = correspondingProposal
		p.Evidences = append(p.Evidences, correspondingProposal, pcForV, pcForNotV)
		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO12 rule, with in-corresponding proposal", func(t *testing.T) {
		correspondingProposal := newProposalMessage(height, 2, 0, proposerKey, committee, nil)
		// a precommit at round 1, with value v.
		pcForV := message.NewPrecommit(1, height, correspondingProposal.Value(), makeSigner(proposerKey))
		// a precommit at round 2, with value not v.
		pcForNotV := message.NewPrecommit(2, height, noneNilValue, makeSigner(proposerKey))

		// a prevote at round 3, with value v.
		preVote := message.NewPrevote(3, height, correspondingProposal.Value(), makeSigner(proposerKey))
		var p Proof
		p.Rule = autonity.PVO12
		p.Type = autonity.Misbehaviour
		p.Message = preVote
		p.Evidences = append(p.Evidences, correspondingProposal, pcForV, pcForNotV)
		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO1 rule, with precommits out of round range", func(t *testing.T) {
		correspondingProposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)
		// a precommit at round 0, with value v.
		pcValidRound := message.NewPrecommit(0, height, correspondingProposal.Value(), makeSigner(proposerKey))
		// a precommit at round 1, with value v.
		pcForV := message.NewPrecommit(1, height, correspondingProposal.Value(), makeSigner(proposerKey))

		// a precommit at round 4, with value not v.
		pcForNotV := message.NewPrecommit(4, height, noneNilValue, makeSigner(proposerKey))

		// a prevote at round 3, with value v.
		preVote := message.NewPrevote(3, height, correspondingProposal.Value(), makeSigner(proposerKey))
		var p Proof
		p.Rule = autonity.PVO12
		p.Type = autonity.Misbehaviour
		p.Message = preVote
		p.Evidences = append(p.Evidences, correspondingProposal, pcValidRound, pcForV, pcForNotV)
		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO rule, with correct proof", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		correspondingProposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)
		maliciousPreVote := message.NewPrevote(3, height, correspondingProposal.Value(), makeSigner(proposerKey))
		var p Proof
		p.Rule = autonity.PVO
		p.Message = maliciousPreVote
		p.Evidences = append(p.Evidences, message.NewLightProposal(correspondingProposal))
		// simulate quorum prevote for not v at valid round.
		for _, c := range committee {
			preVote := message.NewPrevote(0, height, noneNilValue, makeSigner(keys[c.Address]))
			p.Evidences = append(p.Evidences, preVote)
		}
		lastHeader := newBlockHeader(lastHeight, committee)

		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validMisbehaviourOfPVO(&p)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO rule, with less quorum preVote for not v", func(t *testing.T) {
		correspondingProposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)
		maliciousPreVote := message.NewPrevote(3, height, correspondingProposal.Value(), makeSigner(proposerKey))
		var p Proof
		p.Rule = autonity.PVO
		p.Message = maliciousPreVote
		p.Evidences = append(p.Evidences, correspondingProposal)
		// simulate only one prevote for not v at valid round.
		preVote := message.NewPrevote(0, height, noneNilValue, makeSigner(proposerKey))
		p.Evidences = append(p.Evidences, preVote)

		mv := MisbehaviourVerifier{chain: nil}
		ret := mv.validMisbehaviourOfPVO(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO rule, with preVotes at wrong valid round", func(t *testing.T) {
		correspondingProposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)
		maliciousPreVote := message.NewPrevote(3, height, correspondingProposal.Value(), makeSigner(proposerKey))
		var p Proof
		p.Rule = autonity.PVO
		p.Message = maliciousPreVote
		p.Evidences = append(p.Evidences, correspondingProposal)
		// simulate quorum prevote for not v at a round rather than valid round
		for _, c := range committee {
			preVote := message.NewPrevote(1, height, noneNilValue, makeSigner(keys[c.Address]))
			p.Evidences = append(p.Evidences, preVote)
		}

		mv := MisbehaviourVerifier{chain: nil}
		ret := mv.validMisbehaviourOfPVO(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour proof of PVO2 rule, with precommits of V", func(t *testing.T) {
		correspondingProposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)
		// a precommit at round 0, with value not v.
		pcVR := message.NewPrecommit(0, height, correspondingProposal.Value(), makeSigner(proposerKey))
		// a precommit at round 1, with value not v.
		pcR1 := message.NewPrecommit(1, height, noneNilValue, makeSigner(proposerKey))
		// a precommit at round 2, with value not v.
		pcR2 := message.NewPrecommit(2, height, noneNilValue, makeSigner(proposerKey))

		// a prevote at round 3, with value v.
		preVote := message.NewPrevote(3, height, correspondingProposal.Value(), makeSigner(proposerKey))
		var p Proof
		p.Rule = autonity.PVO12
		p.Type = autonity.Misbehaviour
		p.Message = preVote
		p.Evidences = append(p.Evidences, correspondingProposal, pcVR, pcR1, pcR2)
		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with correct Proof", func(t *testing.T) {
		// Node preCommit for a V at round R, but in that round, there were quorum PreVotes for notV at that round.
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = autonity.C
		// Node preCommit for V at round R.
		preCommit := message.NewPrecommit(0, height, noneNilValue, makeSigner(proposerKey))
		p.Message = preCommit
		for _, c := range committee {
			preVote := message.NewPrevote(0, height, common.Hash{0x2}, makeSigner(keys[c.Address]))
			p.Evidences = append(p.Evidences, preVote)
		}
		lastHeader := newBlockHeader(lastHeight, committee)

		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validMisbehaviourOfC(&p)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with no Evidences", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C
		// Node preCommit for V at round R.
		preCommit := message.NewPrecommit(0, height, noneNilValue, makeSigner(proposerKey))
		p.Message = preCommit

		mv := MisbehaviourVerifier{}
		ret := mv.validMisbehaviourOfC(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with wrong preCommit msg", func(t *testing.T) {
		// Node preCommit for nil at round R, not provable
		var p Proof
		p.Rule = autonity.C
		preCommit := message.NewPrecommit(0, height, nilValue, makeSigner(proposerKey))
		p.Message = preCommit
		for _, c := range committee {
			preVote := message.NewPrevote(0, height, noneNilValue, makeSigner(keys[c.Address]))
			p.Evidences = append(p.Evidences, preVote)
		}

		mv := MisbehaviourVerifier{}
		ret := mv.validMisbehaviourOfC(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C

		wrongMsg := message.NewPrevote(0, height, noneNilValue, makeSigner(proposerKey))
		p.Message = wrongMsg
		for _, c := range committee {
			preVote := message.NewPrevote(0, height, nilValue, makeSigner(keys[c.Address]))
			p.Evidences = append(p.Evidences, preVote)
		}

		mv := MisbehaviourVerifier{}
		ret := mv.validMisbehaviourOfC(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with invalid evidence", func(t *testing.T) {
		// the evidence contains same value of preCommit that node preVoted for.
		var p Proof
		p.Rule = autonity.C

		preCommit := message.NewPrecommit(0, height, noneNilValue, makeSigner(proposerKey))
		p.Message = preCommit
		// quorum preVotes of same value, this shouldn't be a valid evidence.
		for _, c := range committee {
			preVote := message.NewPrevote(0, height, noneNilValue, makeSigner(keys[c.Address]))
			p.Evidences = append(p.Evidences, preVote)
		}

		mv := MisbehaviourVerifier{}
		ret := mv.validMisbehaviourOfC(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with invalid evidence: duplicated msg in evidence", func(t *testing.T) {
		// the evidence contains same value of preCommit that node preVoted for.
		var p Proof
		p.Rule = autonity.C

		preCommit := message.NewPrecommit(0, height, noneNilValue, makeSigner(proposerKey))
		p.Message = preCommit
		// duplicated preVotes msg in evidence, should be addressed.
		for i := 0; i < len(committee); i++ {
			preVote := message.NewPrevote(0, height, nilValue, makeSigner(proposerKey))
			p.Evidences = append(p.Evidences, preVote)
		}

		mv := MisbehaviourVerifier{}
		ret := mv.validMisbehaviourOfC(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with invalid evidence: no quorum preVotes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = autonity.C

		preCommit := message.NewPrecommit(0, height, noneNilValue, makeSigner(proposerKey))
		p.Message = preCommit

		// no quorum preVotes msg in evidence, should be addressed.
		preVote := message.NewPrevote(0, height, common.Hash{0x2}, makeSigner(proposerKey))
		p.Evidences = append(p.Evidences, preVote)
		lastHeader := newBlockHeader(lastHeight, committee)

		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)

		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validMisbehaviourOfC(&p)
		assert.Equal(t, false, ret)
	})
}

func TestInnocenceVerifier(t *testing.T) {
	height := uint64(100)
	lastHeight := height - 1
	committee, keys := generateCommittee()
	proposer := committee[0].Address
	proposerKey := keys[proposer]
	noneNilValue := common.Hash{0x1}
	t.Run("Test innocence verifier required gas", func(t *testing.T) {
		iv := InnocenceVerifier{chain: nil}
		assert.Equal(t, params.AutonityAFDContractGasPerKB, iv.RequiredGas(nil))
	})

	t.Run("Test innocence verifier run with nil bytes", func(t *testing.T) {
		iv := InnocenceVerifier{chain: nil}
		ret, err := iv.Run(nil, height)
		assert.Equal(t, failureResult, ret)
		assert.Nil(t, err)
	})

	t.Run("Test validate innocence Proof with invalid Signature() of message", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = autonity.PO
		invalidCommittee, iKeys := generateCommittee()
		invalidProposal := newProposalMessage(height, 1, 0, iKeys[invalidCommittee[0].Address], invalidCommittee, nil)
		p.Message = invalidProposal

		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		iv := InnocenceVerifier{chain: chainMock}
		ret := iv.validateInnocenceProof(&p)
		assert.Equal(t, failureResult, ret)
	})

	t.Run("Test validate innocence Proof, with invalid Signature() of evidence msgs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = autonity.PO
		invalidCommittee, iKeys := generateCommittee()
		proposal := newProposalMessage(height, 1, 0, proposerKey, committee, nil)
		p.Message = proposal
		invalidPreVote := message.NewPrevote(1, height, proposal.Value(), makeSigner(iKeys[invalidCommittee[0].Address]))
		p.Evidences = append(p.Evidences, invalidPreVote)

		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		iv := InnocenceVerifier{chain: chainMock}
		ret := iv.validateInnocenceProof(&p)
		assert.Equal(t, failureResult, ret)
	})

	t.Run("Test validate innocence Proof of PO rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		wrongMsg := message.NewPrevote(1, height, noneNilValue, makeSigner(proposerKey))
		p.Message = wrongMsg

		ret := validInnocenceProofOfPO(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence AccountabilityProof of PO rule, with invalid evidence", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		proposal := newProposalMessage(height, 1, 0, proposerKey, committee, nil)
		p.Message = proposal
		// have preVote at different value than proposal
		invalidPreVote := message.NewPrevote(0, height, noneNilValue, makeSigner(proposerKey))
		p.Evidences = append(p.Evidences, invalidPreVote)
		ret := validInnocenceProofOfPO(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence AccountabilityProof of PO rule, with redundant vote msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		proposal := newProposalMessage(height, 1, 0, proposerKey, committee, nil)
		p.Message = proposal

		preVote := message.NewPrevote(0, height, proposal.Value(), makeSigner(proposerKey))
		p.Evidences = append(p.Evidences, preVote)
		// make redundant msg hack.
		p.Evidences = append(p.Evidences, p.Evidences...)

		ret := validInnocenceProofOfPO(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence AccountabilityProof of PO rule, with not quorum vote msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PO
		proposal := newProposalMessage(height, 1, 0, proposerKey, committee, nil)
		p.Message = proposal

		preVote := message.NewPrevote(0, height, proposal.Value(), makeSigner(proposerKey))
		p.Evidences = append(p.Evidences, preVote)

		ret := validInnocenceProofOfPO(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of PVN rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVN
		wrongMsg := message.NewPrecommit(1, height, noneNilValue, makeSigner(proposerKey))
		p.Message = wrongMsg
		ret := validInnocenceProofOfPVN(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of PVN rule, with a wrong preVote for nil", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVN
		wrongMsg := message.NewPrevote(1, height, nilValue, makeSigner(proposerKey))
		p.Message = wrongMsg
		ret := validInnocenceProofOfPVN(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of PVN rule, with no evidence", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVN
		wrongMsg := message.NewPrevote(1, height, noneNilValue, makeSigner(proposerKey))
		p.Message = wrongMsg
		ret := validInnocenceProofOfPVN(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of PVN rule, with over quorum prevotes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var p Proof
		p.Rule = autonity.PVN
		proposal := newProposalMessage(height, 1, -1, proposerKey, committee, nil)
		p.Evidences = append(p.Evidences, message.NewLightProposal(proposal))
		preVote := message.NewPrevote(1, height, proposal.Value(), makeSigner(proposerKey))
		p.Message = preVote
		ret := validInnocenceProofOfPVN(&p)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate innocence proof of PVO rule, with correct Proof", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = autonity.PVO
		vr := int64(0)
		proposal := newProposalMessage(height, 1, vr, proposerKey, committee, nil)
		preVote := message.NewPrevote(1, height, proposal.Value(), makeSigner(proposerKey))
		p.Message = preVote
		p.Evidences = append(p.Evidences, message.NewLightProposal(proposal))
		// prepare quorum prevotes at valid round.
		for _, c := range committee {
			preVote := message.NewPrevote(vr, height, proposal.Value(), makeSigner(keys[c.Address]))
			p.Evidences = append(p.Evidences, preVote)
		}
		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		ret := validInnocenceProofOfPVO(&p, chainMock)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate innocence proof of PVO rule, with incorrect proposal", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVO
		vr := int64(0)
		// with wrong round in proposal.
		proposal := newProposalMessage(height, 2, vr, proposerKey, committee, nil)
		preVote := message.NewPrevote(1, height, proposal.Value(), makeSigner(proposerKey))
		p.Message = preVote
		p.Evidences = append(p.Evidences, proposal)
		// prepare quorum prevotes at valid round.
		for _, c := range committee {
			preVote := message.NewPrevote(vr, height, proposal.Value(), makeSigner(keys[c.Address]))
			p.Evidences = append(p.Evidences, preVote)
		}

		ret := validInnocenceProofOfPVO(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence proof of PVO rule, with incorrect preVotes", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVO
		vr := int64(0)
		// with wrong round in proposal.
		proposal := newProposalMessage(height, 2, vr, proposerKey, committee, nil)
		preVote := message.NewPrevote(1, height, proposal.Value(), makeSigner(proposerKey))
		p.Message = preVote
		p.Evidences = append(p.Evidences, proposal)
		// prepare quorum prevotes at wrong round.
		for _, c := range committee {
			preVote := message.NewPrevote(1, height, proposal.Value(), makeSigner(keys[c.Address]))
			p.Evidences = append(p.Evidences, preVote)
		}
		ret := validInnocenceProofOfPVO(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence proof of PVO rule, with less than quorum preVotes", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVO
		vr := int64(0)
		proposal := newProposalMessage(height, 2, vr, proposerKey, committee, nil)
		preVote := message.NewPrevote(2, height, proposal.Value(), makeSigner(proposerKey))
		p.Message = preVote
		p.Evidences = append(p.Evidences, proposal)
		// prepare only one prevotes at valid round.
		v := message.NewPrevote(vr, height, proposal.Value(), makeSigner(proposerKey))
		p.Evidences = append(p.Evidences, v)

		ret := validInnocenceProofOfPVO(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence proof of PVO rule, with preVote for not V", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.PVO
		vr := int64(0)
		proposal := newProposalMessage(height, 2, vr, proposerKey, committee, nil)
		preVote := message.NewPrevote(2, height, proposal.Value(), makeSigner(proposerKey))
		p.Message = preVote
		p.Evidences = append(p.Evidences, proposal)
		// prepare only one prevotes at valid round.
		for _, c := range committee {
			preVote := message.NewPrevote(vr, height, noneNilValue, makeSigner(keys[c.Address]))
			p.Evidences = append(p.Evidences, preVote)
		}

		ret := validInnocenceProofOfPVO(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C1
		wrongMsg := message.NewPrevote(1, height, noneNilValue, makeSigner(proposerKey))
		p.Message = wrongMsg
		ret := validInnocenceProofOfC1(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with a wrong preCommit for nil", func(t *testing.T) {
		var p Proof
		p.Rule = autonity.C1
		wrongMsg := message.NewPrecommit(1, height, nilValue, makeSigner(proposerKey))
		p.Message = wrongMsg
		ret := validInnocenceProofOfC1(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with a wrong evidence", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = autonity.C1
		preCommit := message.NewPrecommit(1, height, noneNilValue, makeSigner(proposerKey))
		p.Message = preCommit
		// evidence contains a preVote of a different round
		preVote := message.NewPrevote(0, height, noneNilValue, makeSigner(proposerKey))
		p.Evidences = append(p.Evidences, preVote)
		ret := validInnocenceProofOfC1(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with redundant msgs in evidence ", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = autonity.C1
		preCommit := message.NewPrecommit(1, height, noneNilValue, makeSigner(proposerKey))
		p.Message = preCommit

		preVote := message.NewPrevote(1, height, noneNilValue, makeSigner(proposerKey))
		p.Evidences = append(p.Evidences, preVote)
		p.Evidences = append(p.Evidences, p.Evidences...)
		ret := validInnocenceProofOfC1(&p, nil)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with no quorum votes of evidence ", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = autonity.C1
		preCommit := message.NewPrecommit(1, height, noneNilValue, makeSigner(proposerKey))
		p.Message = preCommit

		preVote := message.NewPrevote(1, height, noneNilValue, makeSigner(proposerKey))
		p.Evidences = append(p.Evidences, preVote)
		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		ret := validInnocenceProofOfC1(&p, chainMock)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with correct evidence ", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = autonity.C1
		preCommit := message.NewPrecommit(1, height, noneNilValue, makeSigner(proposerKey))
		p.Message = preCommit
		for _, c := range committee {
			preVote := message.NewPrevote(1, height, noneNilValue, makeSigner(keys[c.Address]))
			p.Evidences = append(p.Evidences, preVote)
		}

		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		ret := validInnocenceProofOfC1(&p, chainMock)
		assert.Equal(t, true, ret)
	})
}

func TestCheckMsgSignature(t *testing.T) {
	height := uint64(100)
	lastHeight := height - 1
	round := int64(0)
	committee, keys := generateCommittee()

	t.Run("normal case, proposal msg is checked correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		currentHeader := newBlockHeader(lastHeight, committee)
		proposal := newProposalMessage(height, round, -1, keys[committee[0].Address], committee, nil)
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(currentHeader)
		require.Nil(t, checkMsgSignature(chainMock, proposal))
	})

	t.Run("a future msg is received, expect an error of errFutureMsg", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		futureHeight := height + 1
		proposal := newProposalMessage(futureHeight, round, -1, keys[committee[0].Address], committee, nil)
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(height).Return(nil)
		require.Equal(t, errFutureMsg, checkMsgSignature(chainMock, proposal))
	})

	t.Run("chain cannot provide the last header of the height that msg votes on, expect an error of errFutureMsg", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		proposal := newProposalMessage(height, round, -1, keys[committee[0].Address], committee, nil)
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(nil)
		require.Equal(t, errFutureMsg, checkMsgSignature(chainMock, proposal))
	})

	t.Run("abnormal case, msg is not signed by committee", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		wrongCommitte, ks := generateCommittee()
		currentHeader := newBlockHeader(lastHeight, committee)
		proposal := newProposalMessage(height, round, -1, ks[wrongCommitte[0].Address], wrongCommitte, nil)
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(currentHeader)
		require.Equal(t, errNotCommitteeMsg, checkMsgSignature(chainMock, proposal))
	})
}

func TestCheckEquivocation(t *testing.T) {
	height := uint64(100)
	round := int64(0)
	committee, keys := generateCommittee()

	t.Run("check equivocation with valid Proof of equivocation", func(t *testing.T) {
		proposal := newProposalMessage(height, round, -1, keys[committee[0].Address], committee, nil)
		vote1 := message.NewPrevote(round, height, proposal.Value(), makeSigner(keys[committee[0].Address]))
		vote2 := message.NewPrevote(round, height, nilValue, makeSigner(keys[committee[0].Address]))
		var proofs []message.Msg
		proofs = append(proofs, vote2)
		require.Equal(t, errEquivocation, checkEquivocation(vote1, proofs))
	})

	t.Run("check equivocation with invalid Proof of equivocation", func(t *testing.T) {
		proposal := newProposalMessage(height, round, -1, keys[committee[0].Address], committee, nil)
		vote1 := message.NewPrevote(round, height, proposal.Value(), makeSigner(keys[committee[0].Address]))
		var proofs []message.Msg
		proofs = append(proofs, vote1)
		require.Nil(t, checkEquivocation(vote1, proofs))
	})
}

func makeSigner(key *ecdsa.PrivateKey) message.Signer {
	return func(hash common.Hash) (output []byte, err error) {
		return crypto.Sign(hash[:], key)
	}
}
