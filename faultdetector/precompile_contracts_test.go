package faultdetector

import (
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
	t.Run("Test accusation verifier required gas", func(t *testing.T) {
		av := AccusationVerifier{chain: nil}
		assert.Equal(t, uint64(params.AutonityPrecompiledContractGas), av.RequiredGas(nil))
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

	})
}

func TestMisbehaviourVerifier(t *testing.T) {

}

func TestInnocenceVerifier(t *testing.T) {

}
