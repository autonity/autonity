package faultdetector

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/params"
	"github.com/clearmatics/autonity/rlp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContractsManagement(t *testing.T) {
	// register contracts into evm package.
	registerFaultDetectorContracts(nil)
	assert.NotNil(t, vm.PrecompiledContractsByzantium[checkInnocenceAddress])
	assert.NotNil(t, vm.PrecompiledContractsByzantium[checkMisbehaviourAddress])
	assert.NotNil(t, vm.PrecompiledContractsByzantium[checkAccusationAddress])

	assert.NotNil(t, vm.PrecompiledContractsHomestead[checkInnocenceAddress])
	assert.NotNil(t, vm.PrecompiledContractsHomestead[checkMisbehaviourAddress])
	assert.NotNil(t, vm.PrecompiledContractsHomestead[checkAccusationAddress])

	assert.NotNil(t, vm.PrecompiledContractsIstanbul[checkInnocenceAddress])
	assert.NotNil(t, vm.PrecompiledContractsIstanbul[checkMisbehaviourAddress])
	assert.NotNil(t, vm.PrecompiledContractsIstanbul[checkAccusationAddress])

	assert.NotNil(t, vm.PrecompiledContractsYoloV1[checkInnocenceAddress])
	assert.NotNil(t, vm.PrecompiledContractsYoloV1[checkAccusationAddress])
	assert.NotNil(t, vm.PrecompiledContractsYoloV1[checkMisbehaviourAddress])
	// un-register them from evm package.
	unRegisterFaultDetectorContracts()
	assert.Nil(t, vm.PrecompiledContractsByzantium[checkInnocenceAddress])
	assert.Nil(t, vm.PrecompiledContractsByzantium[checkMisbehaviourAddress])
	assert.Nil(t, vm.PrecompiledContractsByzantium[checkAccusationAddress])

	assert.Nil(t, vm.PrecompiledContractsHomestead[checkInnocenceAddress])
	assert.Nil(t, vm.PrecompiledContractsHomestead[checkMisbehaviourAddress])
	assert.Nil(t, vm.PrecompiledContractsHomestead[checkAccusationAddress])

	assert.Nil(t, vm.PrecompiledContractsIstanbul[checkInnocenceAddress])
	assert.Nil(t, vm.PrecompiledContractsIstanbul[checkMisbehaviourAddress])
	assert.Nil(t, vm.PrecompiledContractsIstanbul[checkAccusationAddress])

	assert.Nil(t, vm.PrecompiledContractsYoloV1[checkInnocenceAddress])
	assert.Nil(t, vm.PrecompiledContractsYoloV1[checkAccusationAddress])
	assert.Nil(t, vm.PrecompiledContractsYoloV1[checkMisbehaviourAddress])
}

func TestDecodeProof(t *testing.T) {
	height := uint64(100)
	committee, keys := generateCommittee(5)
	proposer := committee[0].Address
	proposerKey := keys[proposer]
	proposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)
	preCommit := newVoteMsg(height, 3, msgPrecommit, proposerKey, proposal.Value(), committee)

	t.Run("decode with accusation", func(t *testing.T) {
		var p Proof
		p.Rule = PO
		p.Message = proposal

		rp, err := rlp.EncodeToBytes(&p)
		assert.NoError(t, err)

		decodeProof, err := decodeRawProof(rp)
		assert.NoError(t, err)
		assert.Equal(t, PO, decodeProof.Rule)
		assert.Equal(t, proposal.Signature, decodeProof.Message.Signature)
	})

	t.Run("decode with evidence", func(t *testing.T) {
		var p Proof
		p.Rule = PO
		p.Message = proposal
		p.Evidence = append(p.Evidence, preCommit)

		rp, err := rlp.EncodeToBytes(&p)
		assert.NoError(t, err)

		decodeProof, err := decodeRawProof(rp)
		assert.NoError(t, err)
		assert.Equal(t, PO, decodeProof.Rule)
		assert.Equal(t, proposal.Signature, decodeProof.Message.Signature)
		assert.Equal(t, preCommit.Signature, decodeProof.Evidence[0].Signature)
	})
}

func TestAccusationVerifier(t *testing.T) {
	height := uint64(100)
	lastHeight := height - 1
	committee, keys := generateCommittee(5)
	proposer := committee[0].Address
	proposerKey := keys[proposer]
	proposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)

	t.Run("Test accusation verifier required gas", func(t *testing.T) {
		av := AccusationVerifier{}
		assert.Equal(t, params.AutonityPrecompiledContractGas, av.RequiredGas(nil))
	})

	t.Run("Test accusation verifier run with nil bytes", func(t *testing.T) {
		av := AccusationVerifier{}
		ret, err := av.Run(nil)
		assert.Equal(t, failure96Byte, ret)
		assert.Nil(t, err)
	})

	t.Run("Test accusation verifier run with invalid rlp bytes", func(t *testing.T) {
		wrongBytes := failure96Byte
		av := AccusationVerifier{}
		ret, err := av.Run(wrongBytes)
		assert.Equal(t, failure96Byte, ret)
		assert.Nil(t, err)
	})

	t.Run("Test validate accusation, with wrong rule ID", func(t *testing.T) {
		var p Proof
		p.Rule = UnknownRule
		av := AccusationVerifier{}
		assert.Equal(t, failure96Byte, av.validateAccusation(&p))
	})

	t.Run("Test validate accusation, with wrong accusation msg", func(t *testing.T) {
		var p Proof
		av := AccusationVerifier{}
		p.Rule = PO
		preVote := newVoteMsg(height, 0, msgPrevote, proposerKey, proposal.Value(), committee)
		p.Message = preVote
		assert.Equal(t, failure96Byte, av.validateAccusation(&p))

		p.Rule = PVN
		p.Message = proposal
		assert.Equal(t, failure96Byte, av.validateAccusation(&p))

		p.Rule = C
		p.Message = proposal
		assert.Equal(t, failure96Byte, av.validateAccusation(&p))

		p.Rule = C1
		p.Message = proposal
		assert.Equal(t, failure96Byte, av.validateAccusation(&p))
	})

	t.Run("Test validate accusation, with invalid signature of msg", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var p Proof
		p.Rule = PO
		invalidCommittee, keys := generateCommittee(5)
		newProposal := newProposalMessage(height, 1, 0, keys[invalidCommittee[0].Address], invalidCommittee, nil)
		p.Message = newProposal

		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		av := AccusationVerifier{chain: chainMock}

		ret := av.validateAccusation(&p)
		assert.Equal(t, failure96Byte, ret)
	})

	t.Run("Test validate accusation, with correct accusation msg", func(t *testing.T) {
		var p Proof
		p.Rule = PO
		newProposal := newProposalMessage(height, 1, 0, proposerKey, committee, nil)
		p.Message = newProposal
		lastHeader := newBlockHeader(lastHeight, committee)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)

		av := AccusationVerifier{chain: chainMock}
		ret := av.validateAccusation(&p)
		assert.NotEqual(t, failure96Byte, ret)
		assert.Equal(t, common.LeftPadBytes(proposer.Bytes(), 32), ret[0:32])
		assert.Equal(t, types.RLPHash(newProposal).Bytes(), ret[32:64])
		assert.Equal(t, validProofByte, ret[64:96])
	})
}

func TestMisbehaviourVerifier(t *testing.T) {
	height := uint64(100)
	lastHeight := height - 1
	committee, keys := generateCommittee(5)
	proposer := committee[0].Address
	proposerKey := keys[proposer]
	noneNilValue := common.Hash{0x1}

	t.Run("Test misbehaviour verifier required gas", func(t *testing.T) {
		mv := MisbehaviourVerifier{}
		assert.Equal(t, params.AutonityPrecompiledContractGas, mv.RequiredGas(nil))
	})

	t.Run("Test misbehaviour verifier run with nil bytes", func(t *testing.T) {
		mv := MisbehaviourVerifier{}
		ret, err := mv.Run(nil)
		assert.Equal(t, failure96Byte, ret)
		assert.Nil(t, err)
	})

	t.Run("Test misbehaviour verifier run with invalid rlp bytes", func(t *testing.T) {
		wrongBytes := failure96Byte
		mv := MisbehaviourVerifier{}
		ret, err := mv.Run(wrongBytes)
		assert.Equal(t, failure96Byte, ret)
		assert.Nil(t, err)
	})

	t.Run("Test validate misbehaviour Proof, with invalid signature of misbehaved msg", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = PO
		invalidCommittee, iKeys := generateCommittee(5)
		invalidProposal := newProposalMessage(height, 1, 0, iKeys[invalidCommittee[0].Address], invalidCommittee, nil)
		p.Message = invalidProposal

		lastHeader := newBlockHeader(lastHeight, committee)
		currentHeader := newBlockHeader(height, committee)
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		chainMock.EXPECT().CurrentHeader().Return(currentHeader)
		mv := MisbehaviourVerifier{chain: chainMock}

		ret := mv.validateProof(&p)
		assert.Equal(t, failure96Byte, ret)
	})

	t.Run("Test validate misbehaviour Proof, with invalid signature of evidence msgs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = PO
		invalidCommittee, ikeys := generateCommittee(5)
		proposal := newProposalMessage(height, 1, 0, proposerKey, committee, nil)
		p.Message = proposal
		invalidPreCommit := newVoteMsg(height, 1, msgPrecommit, ikeys[invalidCommittee[0].Address], proposal.Value(), invalidCommittee)
		p.Evidence = append(p.Evidence, invalidPreCommit)

		lastHeader := newBlockHeader(lastHeight, committee)
		currentHeader := newBlockHeader(height, committee)
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).AnyTimes().Return(lastHeader)
		chainMock.EXPECT().CurrentHeader().Return(currentHeader)
		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validateProof(&p)
		assert.Equal(t, failure96Byte, ret)
	})

	t.Run("Test validate misbehaviour Proof of PN rule with correct Proof", func(t *testing.T) {
		// prepare a Proof that node propose for a new value, but he preCommitted a non nil value
		// at previous rounds, such Proof should be valid.
		var p Proof
		p.Rule = PN
		proposal := newProposalMessage(height, 1, -1, proposerKey, committee, nil)
		p.Message = proposal

		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		p.Evidence = append(p.Evidence, preCommit)
		mv := MisbehaviourVerifier{}

		ret := mv.validProof(&p)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate misbehaviour Proof of PN rule with incorrect proposal of Proof", func(t *testing.T) {
		// prepare a p that node propose for an old value.
		var p Proof
		p.Rule = PN
		proposal := newProposalMessage(height, 1, 0, proposerKey, committee, nil)
		p.Message = proposal

		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		p.Evidence = append(p.Evidence, preCommit)
		mv := MisbehaviourVerifier{}

		ret := mv.validProof(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of PN rule with no evidence of Proof", func(t *testing.T) {
		var p Proof
		p.Rule = PN
		proposal := newProposalMessage(height, 1, -1, proposerKey, committee, nil)
		p.Message = proposal

		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of PO, propose a v rather than the locked one", func(t *testing.T) {
		// simulate a p of misbehaviour of PO, with the proposer proposed a old value that was not
		// the one he locked at previous round, the validation of this p should return true.
		var p Proof
		p.Rule = PO
		proposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)
		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		p.Message = proposal
		p.Evidence = append(p.Evidence, preCommit)
		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate misbehaviour Proof of PO, propose a valid round rather than the locked one", func(t *testing.T) {
		// simulate a p of misbehaviour of PO, with the proposer proposed a valid round that was not
		// the one he locked at previous round, the validation of this p should return true.
		var p Proof
		p.Rule = PO
		proposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)
		preCommit := newVoteMsg(height, 1, msgPrecommit, proposerKey, noneNilValue, committee)
		p.Message = proposal
		p.Evidence = append(p.Evidence, preCommit)
		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate misbehaviour Proof of PO, with no evidence", func(t *testing.T) {
		var p Proof
		p.Rule = PO
		proposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)
		p.Message = proposal
		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of PO, with a proposal of new value", func(t *testing.T) {
		var p Proof
		p.Rule = PO
		proposal := newProposalMessage(height, 3, -1, proposerKey, committee, nil)
		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		p.Message = proposal
		p.Evidence = append(p.Evidence, preCommit)
		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of PVN rule, with correct Proof", func(t *testing.T) {
		// simulate a p of misbehaviour of PVN, with the node preVote for V1, but he preCommit
		// at a different value V2 at previous round. The validation of the misbehaviour p should
		// return ture.
		var p Proof
		p.Rule = PVN
		// node locked at V1 at round 0.
		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		proposal := newProposalMessage(height, 3, -1, proposerKey, committee, nil)
		// node preVote for V2 at round 3
		preVote := newVoteMsg(height, 3, msgPrevote, proposerKey, proposal.Value(), committee)
		p.Message = preVote
		p.Evidence = append(p.Evidence, preCommit)

		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate misbehaviour Proof of PVN rule, with no evidence", func(t *testing.T) {
		var p Proof
		p.Rule = PVN
		proposal := newProposalMessage(height, 3, -1, proposerKey, committee, nil)
		// node preVote for V2 at round 3
		preVote := newVoteMsg(height, 3, msgPrevote, proposerKey, proposal.Value(), committee)
		p.Message = preVote

		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of PVN rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = PVN
		proposal := newProposalMessage(height, 3, -1, proposerKey, committee, nil)
		// set a wrong type of msg.
		p.Message = proposal
		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		p.Evidence = append(p.Evidence, preCommit)
		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of PVN rule, with wrong preVote value", func(t *testing.T) {
		var p Proof
		p.Rule = PVN
		// node locked at V1 at round 0.
		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		// node preVote for V2 at round 3, with nil value, not provable.
		preVote := newVoteMsg(height, 3, msgPrevote, proposerKey, nilValue, committee)
		p.Message = preVote
		p.Evidence = append(p.Evidence, preCommit)

		mv := MisbehaviourVerifier{}
		ret := mv.validProof(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with correct Proof", func(t *testing.T) {
		// Node preCommit for a V at round R, but in that round, there were quorum PreVotes for notV at that round.
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = C
		// Node preCommit for V at round R.
		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		p.Message = preCommit
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, 0, msgPrevote, keys[committee[i].Address], nilValue, committee)
			p.Evidence = append(p.Evidence, preVote)
		}
		lastHeader := newBlockHeader(lastHeight, committee)

		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validMisbehaviourOfC(&p)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with no Evidence", func(t *testing.T) {
		var p Proof
		p.Rule = C
		// Node preCommit for V at round R.
		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		p.Message = preCommit

		mv := MisbehaviourVerifier{}
		ret := mv.validMisbehaviourOfC(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with wrong preCommit msg", func(t *testing.T) {
		// Node preCommit for nil at round R, not provable
		var p Proof
		p.Rule = C
		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, nilValue, committee)
		p.Message = preCommit
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, 0, msgPrevote, keys[committee[i].Address], noneNilValue, committee)
			p.Evidence = append(p.Evidence, preVote)
		}

		mv := MisbehaviourVerifier{}
		ret := mv.validMisbehaviourOfC(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = C

		wrongMsg := newVoteMsg(height, 0, msgPrevote, proposerKey, noneNilValue, committee)
		p.Message = wrongMsg
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, 0, msgPrevote, keys[committee[i].Address], nilValue, committee)
			p.Evidence = append(p.Evidence, preVote)
		}

		mv := MisbehaviourVerifier{}
		ret := mv.validMisbehaviourOfC(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with invalid evidence", func(t *testing.T) {
		// the evidence contains same value of preCommit that node preVoted for.
		var p Proof
		p.Rule = C

		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		p.Message = preCommit
		// quorum preVotes of same value, this shouldn't be a valid evidence.
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, 0, msgPrevote, keys[committee[i].Address], noneNilValue, committee)
			p.Evidence = append(p.Evidence, preVote)
		}

		mv := MisbehaviourVerifier{}
		ret := mv.validMisbehaviourOfC(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with invalid evidence: duplicated msg in evidence", func(t *testing.T) {
		// the evidence contains same value of preCommit that node preVoted for.
		var p Proof
		p.Rule = C

		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		p.Message = preCommit
		// duplicated preVotes msg in evidence, should be addressed.
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, 0, msgPrevote, proposerKey, nilValue, committee)
			p.Evidence = append(p.Evidence, preVote)
		}

		mv := MisbehaviourVerifier{}
		ret := mv.validMisbehaviourOfC(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour Proof of C rule, with invalid evidence: no quorum preVotes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = C

		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		p.Message = preCommit

		// no quorum preVotes msg in evidence, should be addressed.
		preVote := newVoteMsg(height, 0, msgPrevote, proposerKey, nilValue, committee)
		p.Evidence = append(p.Evidence, preVote)
		lastHeader := newBlockHeader(lastHeight, committee)

		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)

		mv := MisbehaviourVerifier{chain: chainMock}
		ret := mv.validMisbehaviourOfC(&p)
		assert.Equal(t, false, ret)
	})
}

func TestInnocenceVerifier(t *testing.T) {
	height := uint64(100)
	lastHeight := height - 1
	committee, keys := generateCommittee(5)
	proposer := committee[0].Address
	proposerKey := keys[proposer]
	noneNilValue := common.Hash{0x1}
	t.Run("Test innocence verifier required gas", func(t *testing.T) {
		iv := InnocenceVerifier{chain: nil}
		assert.Equal(t, params.AutonityPrecompiledContractGas, iv.RequiredGas(nil))
	})

	t.Run("Test innocence verifier run with nil bytes", func(t *testing.T) {
		iv := InnocenceVerifier{chain: nil}
		ret, err := iv.Run(nil)
		assert.Equal(t, failure96Byte, ret)
		assert.Nil(t, err)
	})

	t.Run("Test validate innocence Proof with invalid signature of message", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = PO
		invalidCommittee, iKeys := generateCommittee(5)
		invalidProposal := newProposalMessage(height, 1, 0, iKeys[invalidCommittee[0].Address], invalidCommittee, nil)
		p.Message = invalidProposal

		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		iv := InnocenceVerifier{chain: chainMock}
		ret := iv.validateInnocenceProof(&p)
		assert.Equal(t, failure96Byte, ret)
	})

	t.Run("Test validate innocence Proof, with invalid signature of evidence msgs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = PO
		invalidCommittee, iKeys := generateCommittee(5)
		proposal := newProposalMessage(height, 1, 0, proposerKey, committee, nil)
		p.Message = proposal
		invalidPreVote := newVoteMsg(height, 1, msgPrevote, iKeys[invalidCommittee[0].Address], proposal.Value(), invalidCommittee)
		p.Evidence = append(p.Evidence, invalidPreVote)

		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		iv := InnocenceVerifier{chain: chainMock}
		ret := iv.validateInnocenceProof(&p)
		assert.Equal(t, failure96Byte, ret)
	})

	t.Run("Test validate innocence Proof of PO rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = PO
		wrongMsg := newVoteMsg(height, 1, msgPrevote, proposerKey, noneNilValue, committee)
		p.Message = wrongMsg
		iv := InnocenceVerifier{chain: nil}
		ret := iv.validInnocenceProofOfPO(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of PO rule, with invalid evidence", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = PO
		proposal := newProposalMessage(height, 1, 0, proposerKey, committee, nil)
		p.Message = proposal
		// have preVote at different value than proposal
		invalidPreVote := newVoteMsg(height, 0, msgPrevote, proposerKey, noneNilValue, committee)
		p.Evidence = append(p.Evidence, invalidPreVote)
		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)

		iv := InnocenceVerifier{chain: chainMock}
		ret := iv.validInnocenceProofOfPO(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of PO rule, with redundant vote msg", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = PO
		proposal := newProposalMessage(height, 1, 0, proposerKey, committee, nil)
		p.Message = proposal

		preVote := newVoteMsg(height, 0, msgPrevote, proposerKey, proposal.Value(), committee)
		p.Evidence = append(p.Evidence, preVote)
		// make redundant msg hack.
		p.Evidence = append(p.Evidence, p.Evidence...)

		lastHeader := newBlockHeader(lastHeight, committee)

		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		iv := InnocenceVerifier{chain: chainMock}
		ret := iv.validInnocenceProofOfPO(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of PO rule, with not quorum vote msg", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = PO
		proposal := newProposalMessage(height, 1, 0, proposerKey, committee, nil)
		p.Message = proposal

		preVote := newVoteMsg(height, 0, msgPrevote, proposerKey, proposal.Value(), committee)
		p.Evidence = append(p.Evidence, preVote)

		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)

		iv := InnocenceVerifier{chain: chainMock}
		ret := iv.validInnocenceProofOfPO(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of PVN rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = PVN
		wrongMsg := newVoteMsg(height, 1, msgPrecommit, proposerKey, noneNilValue, committee)
		p.Message = wrongMsg
		iv := InnocenceVerifier{chain: nil}
		ret := iv.validInnocenceProofOfPVN(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of PVN rule, with a wrong preVote for nil", func(t *testing.T) {
		var p Proof
		p.Rule = PVN
		wrongMsg := newVoteMsg(height, 1, msgPrevote, proposerKey, nilValue, committee)
		p.Message = wrongMsg
		iv := InnocenceVerifier{chain: nil}
		ret := iv.validInnocenceProofOfPVN(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of PVN rule, with no evidence", func(t *testing.T) {
		var p Proof
		p.Rule = PVN
		wrongMsg := newVoteMsg(height, 1, msgPrevote, proposerKey, noneNilValue, committee)
		p.Message = wrongMsg
		iv := InnocenceVerifier{chain: nil}
		ret := iv.validInnocenceProofOfPVN(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of PVN rule, with correct Proof", func(t *testing.T) {
		var p Proof
		p.Rule = PVN
		proposal := newProposalMessage(height, 1, -1, proposerKey, committee, nil)
		preVote := newVoteMsg(height, 1, msgPrevote, proposerKey, proposal.Value(), committee)
		p.Message = preVote
		p.Evidence = append(p.Evidence, proposal)
		iv := InnocenceVerifier{chain: nil}
		ret := iv.validInnocenceProofOfPVN(&p)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate innocence Proof of C rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = C
		wrongMsg := newVoteMsg(height, 1, msgPrevote, proposerKey, noneNilValue, committee)
		p.Message = wrongMsg
		iv := InnocenceVerifier{chain: nil}
		ret := iv.validInnocenceProofOfC(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C rule, with a wrong preCommit for nil", func(t *testing.T) {
		var p Proof
		p.Rule = C
		wrongMsg := newVoteMsg(height, 1, msgPrecommit, proposerKey, nilValue, committee)
		p.Message = wrongMsg
		iv := InnocenceVerifier{chain: nil}
		ret := iv.validInnocenceProofOfC(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C rule, with no evidence", func(t *testing.T) {
		var p Proof
		p.Rule = C
		preCommit := newVoteMsg(height, 1, msgPrecommit, proposerKey, noneNilValue, committee)
		p.Message = preCommit
		iv := InnocenceVerifier{chain: nil}
		ret := iv.validInnocenceProofOfC(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C rule, with correct Proof", func(t *testing.T) {
		var p Proof
		p.Rule = C
		proposal := newProposalMessage(height, 1, -1, proposerKey, committee, nil)
		preCommit := newVoteMsg(height, 1, msgPrecommit, proposerKey, proposal.Value(), committee)
		p.Message = preCommit
		p.Evidence = append(p.Evidence, proposal)
		iv := InnocenceVerifier{chain: nil}
		ret := iv.validInnocenceProofOfC(&p)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with wrong msg", func(t *testing.T) {
		var p Proof
		p.Rule = C1
		wrongMsg := newVoteMsg(height, 1, msgPrevote, proposerKey, noneNilValue, committee)
		p.Message = wrongMsg
		iv := InnocenceVerifier{chain: nil}
		ret := iv.validInnocenceProofOfC1(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with a wrong preCommit for nil", func(t *testing.T) {
		var p Proof
		p.Rule = C1
		wrongMsg := newVoteMsg(height, 1, msgPrecommit, proposerKey, nilValue, committee)
		p.Message = wrongMsg
		iv := InnocenceVerifier{chain: nil}
		ret := iv.validInnocenceProofOfC1(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with a wrong evidence", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = C1
		preCommit := newVoteMsg(height, 1, msgPrecommit, proposerKey, noneNilValue, committee)
		p.Message = preCommit
		// evidence contains a preVote of a different round
		preVote := newVoteMsg(height, 0, msgPrevote, proposerKey, noneNilValue, committee)
		p.Evidence = append(p.Evidence, preVote)
		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		iv := InnocenceVerifier{chain: chainMock}
		ret := iv.validInnocenceProofOfC1(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with redundant msgs in evidence ", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = C1
		preCommit := newVoteMsg(height, 1, msgPrecommit, proposerKey, noneNilValue, committee)
		p.Message = preCommit

		preVote := newVoteMsg(height, 1, msgPrevote, proposerKey, noneNilValue, committee)
		p.Evidence = append(p.Evidence, preVote)
		p.Evidence = append(p.Evidence, p.Evidence...)
		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		iv := InnocenceVerifier{chain: chainMock}
		ret := iv.validInnocenceProofOfC1(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with no quorum votes of evidence ", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = C1
		preCommit := newVoteMsg(height, 1, msgPrecommit, proposerKey, noneNilValue, committee)
		p.Message = preCommit

		preVote := newVoteMsg(height, 1, msgPrevote, proposerKey, noneNilValue, committee)
		p.Evidence = append(p.Evidence, preVote)
		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		iv := InnocenceVerifier{chain: chainMock}
		ret := iv.validInnocenceProofOfC1(&p)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate innocence Proof of C1 rule, with correct evidence ", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var p Proof
		p.Rule = C1
		preCommit := newVoteMsg(height, 1, msgPrecommit, proposerKey, noneNilValue, committee)
		p.Message = preCommit
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, 1, msgPrevote, keys[committee[i].Address], noneNilValue, committee)
			p.Evidence = append(p.Evidence, preVote)
		}

		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		iv := InnocenceVerifier{chain: chainMock}
		ret := iv.validInnocenceProofOfC1(&p)
		assert.Equal(t, true, ret)
	})
}
