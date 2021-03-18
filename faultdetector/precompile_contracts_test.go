package faultdetector

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/params"
	"github.com/clearmatics/autonity/rlp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContractsManagement(t *testing.T) {
	// register contracts into evm package.
	registerAFDContracts(nil)
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
	unRegisterAFDContracts()
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
		var rawProof RawProof
		rawProof.Rule = PO
		rawProof.Message = proposal.Payload()

		rp, err := rlp.EncodeToBytes(&rawProof)
		assert.NoError(t, err)

		decodeProof, err := decodeProof(rp)
		assert.NoError(t, err)
		assert.Equal(t, PO, decodeProof.Rule)
		assert.Equal(t, proposal.Signature, decodeProof.Message.Signature)
	})

	t.Run("decode with evidence", func(t *testing.T) {
		var rawProof RawProof
		rawProof.Rule = PO
		rawProof.Message = proposal.Payload()
		rawProof.Evidence = append(rawProof.Evidence, preCommit.Payload())

		rp, err := rlp.EncodeToBytes(&rawProof)
		assert.NoError(t, err)

		decodeProof, err := decodeProof(rp)
		assert.NoError(t, err)
		assert.Equal(t, PO, decodeProof.Rule)
		assert.Equal(t, proposal.Signature, decodeProof.Message.Signature)
		assert.Equal(t, preCommit.Signature, decodeProof.Evidence[0].Signature)
	})
}

func TestAccusationVerifier(t *testing.T) {
	height := uint64(100)
	committee, keys := generateCommittee(5)
	proposer := committee[0].Address
	proposerKey := keys[proposer]
	proposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)

	t.Run("Test accusation verifier required gas", func(t *testing.T) {
		av := AccusationVerifier{chain: nil}
		assert.Equal(t, params.AutonityPrecompiledContractGas, av.RequiredGas(nil))
	})

	t.Run("Test accusation verifier run with nil bytes", func(t *testing.T) {
		av := AccusationVerifier{chain: nil}
		ret, err := av.Run(nil)
		assert.Equal(t, failure96Byte, ret)
		assert.Nil(t, err)
	})

	t.Run("Test accusation verifier run with invalid rlp bytes", func(t *testing.T) {
		wrongBytes := failure96Byte
		av := AccusationVerifier{chain: nil}
		ret, err := av.Run(wrongBytes)
		assert.Equal(t, failure96Byte, ret)
		assert.Nil(t, err)
	})

	t.Run("Test validate accusation, with wrong rule ID", func(t *testing.T) {
		var proof Proof
		proof.Rule = UnknownRule
		av := AccusationVerifier{chain: nil}
		assert.Equal(t, failure96Byte, av.validateAccusation(&proof, getHeader))
	})

	t.Run("Test validate accusation, with wrong accusation msg", func(t *testing.T) {
		var proof Proof
		av := AccusationVerifier{chain: nil}
		proof.Rule = PO
		preVote := newVoteMsg(height, 0, msgPrevote, proposerKey, proposal.Value(), committee)
		proof.Message = *preVote
		assert.Equal(t, failure96Byte, av.validateAccusation(&proof, getHeader))

		proof.Rule = PVN
		proof.Message = *proposal
		assert.Equal(t, failure96Byte, av.validateAccusation(&proof, getHeader))

		proof.Rule = C
		proof.Message = *proposal
		assert.Equal(t, failure96Byte, av.validateAccusation(&proof, getHeader))

		proof.Rule = C1
		proof.Message = *proposal
		assert.Equal(t, failure96Byte, av.validateAccusation(&proof, getHeader))
	})

	t.Run("Test validate accusation, with invalid signature of msg", func(t *testing.T) {
		var proof Proof
		proof.Rule = PO
		invalidCommittee, keys := generateCommittee(5)
		newProposal := newProposalMessage(height, 1, 0, keys[invalidCommittee[0].Address], invalidCommittee, nil)
		proof.Message = *newProposal

		lastHeader := newBlockHeader(height-1, committee)
		mockGetHeader := func(_ *core.BlockChain, _ uint64) *types.Header {
			return lastHeader
		}

		av := AccusationVerifier{chain: nil}
		ret := av.validateAccusation(&proof, mockGetHeader)
		assert.Equal(t, failure96Byte, ret)
	})

	t.Run("Test validate accusation, with correct accusation msg", func(t *testing.T) {
		av := AccusationVerifier{chain: nil}
		var proof Proof
		proof.Rule = PO
		newProposal := newProposalMessage(height, 1, 0, proposerKey, committee, nil)
		proof.Message = *newProposal
		lastHeader := newBlockHeader(height-1, committee)

		mockGetHeader := func(_ *core.BlockChain, _ uint64) *types.Header {
			return lastHeader
		}

		ret := av.validateAccusation(&proof, mockGetHeader)
		assert.NotEqual(t, failure96Byte, ret)
		assert.Equal(t, common.LeftPadBytes(proposer.Bytes(), 32), ret[0:32])
		assert.Equal(t, types.RLPHash(newProposal.Payload()).Bytes(), ret[32:64])
		assert.Equal(t, validProofByte, ret[64:96])
	})
}

func TestMisbehaviourVerifier(t *testing.T) {
	height := uint64(100)
	committee, keys := generateCommittee(5)
	proposer := committee[0].Address
	proposerKey := keys[proposer]
	noneNilValue := common.Hash{0x1}

	t.Run("Test misbehaviour verifier required gas", func(t *testing.T) {
		mv := MisbehaviourVerifier{chain: nil}
		assert.Equal(t, params.AutonityPrecompiledContractGas, mv.RequiredGas(nil))
	})

	t.Run("Test misbehaviour verifier run with nil bytes", func(t *testing.T) {
		mv := MisbehaviourVerifier{chain: nil}
		ret, err := mv.Run(nil)
		assert.Equal(t, failure96Byte, ret)
		assert.Nil(t, err)
	})

	t.Run("Test misbehaviour verifier run with invalid rlp bytes", func(t *testing.T) {
		wrongBytes := failure96Byte
		mv := MisbehaviourVerifier{chain: nil}
		ret, err := mv.Run(wrongBytes)
		assert.Equal(t, failure96Byte, ret)
		assert.Nil(t, err)
	})

	t.Run("Test validate misbehaviour proof, with invalid signature of misbehaved msg", func(t *testing.T) {
		var proof Proof
		proof.Rule = PO
		invalidCommittee, iKeys := generateCommittee(5)
		invalidProposal := newProposalMessage(height, 1, 0, iKeys[invalidCommittee[0].Address], invalidCommittee, nil)
		proof.Message = *invalidProposal

		lastHeader := newBlockHeader(height-1, committee)
		mockGetHeader := func(_ *core.BlockChain, _ uint64) *types.Header {
			return lastHeader
		}
		currentHeader := newBlockHeader(height, committee)
		mockCurrentHeader := func(_ *core.BlockChain) *types.Header {
			return currentHeader
		}

		mv := MisbehaviourVerifier{chain: nil}
		ret := mv.validateProof(&proof, mockGetHeader, mockCurrentHeader)
		assert.Equal(t, failure96Byte, ret)
	})

	t.Run("Test validate misbehaviour proof, with invalid signature of evidence msgs", func(t *testing.T) {
		var proof Proof
		proof.Rule = PO
		invalidCommittee, ikeys := generateCommittee(5)
		proposal := newProposalMessage(height, 1, 0, proposerKey, committee, nil)
		proof.Message = *proposal
		invalidPreCommit := newVoteMsg(height, 1, msgPrecommit, ikeys[invalidCommittee[0].Address], proposal.Value(), invalidCommittee)
		proof.Evidence = append(proof.Evidence, *invalidPreCommit)

		lastHeader := newBlockHeader(height-1, committee)
		mockGetHeader := func(_ *core.BlockChain, _ uint64) *types.Header {
			return lastHeader
		}
		currentHeader := newBlockHeader(height, committee)
		mockCurrentHeader := func(_ *core.BlockChain) *types.Header {
			return currentHeader
		}

		mv := MisbehaviourVerifier{chain: nil}
		ret := mv.validateProof(&proof, mockGetHeader, mockCurrentHeader)
		assert.Equal(t, failure96Byte, ret)
	})

	t.Run("Test validate misbehaviour proof of PN rule with correct proof", func(t *testing.T) {
		// prepare a proof that node propose for a new value, but he preCommitted a non nil value
		// at previous rounds, such proof should be valid.
		var proof Proof
		proof.Rule = PN
		proposal := newProposalMessage(height, 1, -1, proposerKey, committee, nil)
		proof.Message = *proposal

		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		proof.Evidence = append(proof.Evidence, *preCommit)
		mv := MisbehaviourVerifier{chain: nil}

		ret := mv.validProof(&proof)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate misbehaviour proof of PN rule with incorrect proposal of proof", func(t *testing.T) {
		// prepare a proof that node propose for an old value.
		var proof Proof
		proof.Rule = PN
		proposal := newProposalMessage(height, 1, 0, proposerKey, committee, nil)
		proof.Message = *proposal

		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		proof.Evidence = append(proof.Evidence, *preCommit)
		mv := MisbehaviourVerifier{chain: nil}

		ret := mv.validProof(&proof)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour proof of PN rule with no evidence of proof", func(t *testing.T) {
		var proof Proof
		proof.Rule = PN
		proposal := newProposalMessage(height, 1, -1, proposerKey, committee, nil)
		proof.Message = *proposal

		mv := MisbehaviourVerifier{chain: nil}

		ret := mv.validProof(&proof)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour proof of PO, propose a v rather than the locked one", func(t *testing.T) {
		// simulate a proof of misbehaviour of PO, with the proposer proposed a old value that was not
		// the one he locked at previous round, the validation of this proof should return true.
		var proof Proof
		proof.Rule = PO
		proposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)
		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		proof.Message = *proposal
		proof.Evidence = append(proof.Evidence, *preCommit)
		mv := MisbehaviourVerifier{chain: nil}
		ret := mv.validProof(&proof)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate misbehaviour proof of PO, propose a valid round rather than the locked one", func(t *testing.T) {
		// simulate a proof of misbehaviour of PO, with the proposer proposed a valid round that was not
		// the one he locked at previous round, the validation of this proof should return true.
		var proof Proof
		proof.Rule = PO
		proposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)
		preCommit := newVoteMsg(height, 1, msgPrecommit, proposerKey, noneNilValue, committee)
		proof.Message = *proposal
		proof.Evidence = append(proof.Evidence, *preCommit)
		mv := MisbehaviourVerifier{chain: nil}
		ret := mv.validProof(&proof)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate misbehaviour proof of PO, with no evidence", func(t *testing.T) {
		var proof Proof
		proof.Rule = PO
		proposal := newProposalMessage(height, 3, 0, proposerKey, committee, nil)
		proof.Message = *proposal
		mv := MisbehaviourVerifier{chain: nil}
		ret := mv.validProof(&proof)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour proof of PO, with a proposal of new value", func(t *testing.T) {
		var proof Proof
		proof.Rule = PO
		proposal := newProposalMessage(height, 3, -1, proposerKey, committee, nil)
		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		proof.Message = *proposal
		proof.Evidence = append(proof.Evidence, *preCommit)
		mv := MisbehaviourVerifier{chain: nil}
		ret := mv.validProof(&proof)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour proof of PVN rule, with correct proof", func(t *testing.T) {
		// simulate a proof of misbehaviour of PVN, with the node preVote for V1, but he preCommit
		// at a different value V2 at previous round. The validation of the misbehaviour proof should
		// return ture.
		var proof Proof
		proof.Rule = PVN
		// node locked at V1 at round 0.
		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		proposal := newProposalMessage(height, 3, -1, proposerKey, committee, nil)
		// node preVote for V2 at round 3
		preVote := newVoteMsg(height, 3, msgPrevote, proposerKey, proposal.Value(), committee)
		proof.Message = *preVote
		proof.Evidence = append(proof.Evidence, *preCommit)

		mv := MisbehaviourVerifier{chain: nil}
		ret := mv.validProof(&proof)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate misbehaviour proof of PVN rule, with no evidence", func(t *testing.T) {
		var proof Proof
		proof.Rule = PVN
		proposal := newProposalMessage(height, 3, -1, proposerKey, committee, nil)
		// node preVote for V2 at round 3
		preVote := newVoteMsg(height, 3, msgPrevote, proposerKey, proposal.Value(), committee)
		proof.Message = *preVote

		mv := MisbehaviourVerifier{chain: nil}
		ret := mv.validProof(&proof)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour proof of PVN rule, with wrong msg", func(t *testing.T) {
		var proof Proof
		proof.Rule = PVN
		proposal := newProposalMessage(height, 3, -1, proposerKey, committee, nil)
		// set a wrong type of msg.
		proof.Message = *proposal
		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		proof.Evidence = append(proof.Evidence, *preCommit)
		mv := MisbehaviourVerifier{chain: nil}
		ret := mv.validProof(&proof)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour proof of PVN rule, with wrong preVote value", func(t *testing.T) {
		var proof Proof
		proof.Rule = PVN
		// node locked at V1 at round 0.
		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		// node preVote for V2 at round 3, with nil value, not provable.
		preVote := newVoteMsg(height, 3, msgPrevote, proposerKey, nilValue, committee)
		proof.Message = *preVote
		proof.Evidence = append(proof.Evidence, *preCommit)

		mv := MisbehaviourVerifier{chain: nil}
		ret := mv.validProof(&proof)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour proof of C rule, with correct proof", func(t *testing.T) {
		// Node preCommit for a V at round R, but in that round, there were quorum PreVotes for notV at that round.
		var proof Proof
		proof.Rule = C
		// Node preCommit for V at round R.
		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		proof.Message = *preCommit
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, 0, msgPrevote, keys[committee[i].Address], nilValue, committee)
			proof.Evidence = append(proof.Evidence, *preVote)
		}
		lastHeader := newBlockHeader(height-1, committee)
		mockGetHeader := func(_ *core.BlockChain, _ uint64) *types.Header {
			return lastHeader
		}

		mv := MisbehaviourVerifier{chain: nil}
		ret := mv.validMisbehaviourOfC(&proof, mockGetHeader)
		assert.Equal(t, true, ret)
	})

	t.Run("Test validate misbehaviour proof of C rule, with no Evidence", func(t *testing.T) {
		var proof Proof
		proof.Rule = C
		// Node preCommit for V at round R.
		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		proof.Message = *preCommit
		lastHeader := newBlockHeader(height-1, committee)
		mockGetHeader := func(_ *core.BlockChain, _ uint64) *types.Header {
			return lastHeader
		}

		mv := MisbehaviourVerifier{chain: nil}
		ret := mv.validMisbehaviourOfC(&proof, mockGetHeader)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour proof of C rule, with wrong preCommit msg", func(t *testing.T) {
		var proof Proof
		proof.Rule = C
		// Node preCommit for nil at round R, not provable
		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, nilValue, committee)
		proof.Message = *preCommit
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, 0, msgPrevote, keys[committee[i].Address], noneNilValue, committee)
			proof.Evidence = append(proof.Evidence, *preVote)
		}
		lastHeader := newBlockHeader(height-1, committee)
		mockGetHeader := func(_ *core.BlockChain, _ uint64) *types.Header {
			return lastHeader
		}

		mv := MisbehaviourVerifier{chain: nil}
		ret := mv.validMisbehaviourOfC(&proof, mockGetHeader)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour proof of C rule, with wrong msg", func(t *testing.T) {
		var proof Proof
		proof.Rule = C

		wrongMsg := newVoteMsg(height, 0, msgPrevote, proposerKey, noneNilValue, committee)
		proof.Message = *wrongMsg
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, 0, msgPrevote, keys[committee[i].Address], nilValue, committee)
			proof.Evidence = append(proof.Evidence, *preVote)
		}
		lastHeader := newBlockHeader(height-1, committee)
		mockGetHeader := func(_ *core.BlockChain, _ uint64) *types.Header {
			return lastHeader
		}

		mv := MisbehaviourVerifier{chain: nil}
		ret := mv.validMisbehaviourOfC(&proof, mockGetHeader)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour proof of C rule, with invalid evidence", func(t *testing.T) {
		// the evidence contains same value of preCommit that node preVoted for.
		var proof Proof
		proof.Rule = C

		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		proof.Message = *preCommit
		// quorum preVotes of same value, this shouldn't be a valid evidence.
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, 0, msgPrevote, keys[committee[i].Address], noneNilValue, committee)
			proof.Evidence = append(proof.Evidence, *preVote)
		}
		lastHeader := newBlockHeader(height-1, committee)
		mockGetHeader := func(_ *core.BlockChain, _ uint64) *types.Header {
			return lastHeader
		}

		mv := MisbehaviourVerifier{chain: nil}
		ret := mv.validMisbehaviourOfC(&proof, mockGetHeader)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour proof of C rule, with invalid evidence: duplicated msg in evidence", func(t *testing.T) {
		// the evidence contains same value of preCommit that node preVoted for.
		var proof Proof
		proof.Rule = C

		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		proof.Message = *preCommit
		// duplicated preVotes msg in evidence, should be addressed.
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, 0, msgPrevote, proposerKey, nilValue, committee)
			proof.Evidence = append(proof.Evidence, *preVote)
		}
		lastHeader := newBlockHeader(height-1, committee)
		mockGetHeader := func(_ *core.BlockChain, _ uint64) *types.Header {
			return lastHeader
		}

		mv := MisbehaviourVerifier{chain: nil}
		ret := mv.validMisbehaviourOfC(&proof, mockGetHeader)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour proof of C rule, with invalid evidence: no quorum preVotes", func(t *testing.T) {
		var proof Proof
		proof.Rule = C

		preCommit := newVoteMsg(height, 0, msgPrecommit, proposerKey, noneNilValue, committee)
		proof.Message = *preCommit

		// no quorum preVotes msg in evidence, should be addressed.
		preVote := newVoteMsg(height, 0, msgPrevote, proposerKey, nilValue, committee)
		proof.Evidence = append(proof.Evidence, *preVote)

		lastHeader := newBlockHeader(height-1, committee)
		mockGetHeader := func(_ *core.BlockChain, _ uint64) *types.Header {
			return lastHeader
		}

		mv := MisbehaviourVerifier{chain: nil}
		ret := mv.validMisbehaviourOfC(&proof, mockGetHeader)
		assert.Equal(t, false, ret)
	})

	t.Run("Test validate misbehaviour proof of Garbage Message", func(t *testing.T) {

	})

	t.Run("Test validate misbehaviour proof of Invalid Proposal", func(t *testing.T) {

	})

	t.Run("Test validate misbehaviour proof of Invalid Proposer", func(t *testing.T) {

	})

	t.Run("Test validate misbehaviour proof of equivocation", func(t *testing.T) {

	})
}

func TestInnocenceVerifier(t *testing.T) {

}
