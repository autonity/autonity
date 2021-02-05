package afd

import (
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/crypto"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/params"
	"github.com/clearmatics/autonity/rlp"
)

var (
	checkProofAddress = common.BytesToAddress([]byte{253})
	checkChallengeAddress = common.BytesToAddress([]byte{254})
	// true32Byte is returned if the bn256 pairing check succeeds.
	true32Byte = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	// false32Byte is returned if the bn256 pairing check fails.
	false32Byte = make([]byte, 32)
)

func initAccountabilityContracts(chain *core.BlockChain) {

	proofChecker := checkProof{blockchain: chain}
	challengeChecker := checkChallenge{blockchain: chain}

	vm.PrecompiledContractsByzantium[checkProofAddress] = &proofChecker
	vm.PrecompiledContractsByzantium[checkChallengeAddress] = &challengeChecker

	vm.PrecompiledContractsHomestead[checkProofAddress] = &proofChecker
	vm.PrecompiledContractsHomestead[checkChallengeAddress] = &challengeChecker

	vm.PrecompiledContractsIstanbul[checkProofAddress] = &proofChecker
	vm.PrecompiledContractsIstanbul[checkChallengeAddress] = &challengeChecker

	vm.PrecompiledContractsYoloV1[checkProofAddress] = &proofChecker
	vm.PrecompiledContractsYoloV1[checkChallengeAddress] = &challengeChecker
}

func cleanContracts() {
	delete(vm.PrecompiledContractsByzantium, checkProofAddress)
	delete(vm.PrecompiledContractsByzantium, checkChallengeAddress)

	delete(vm.PrecompiledContractsYoloV1, checkProofAddress)
	delete(vm.PrecompiledContractsYoloV1, checkChallengeAddress)

	delete(vm.PrecompiledContractsIstanbul, checkProofAddress)
	delete(vm.PrecompiledContractsIstanbul, checkChallengeAddress)

	delete(vm.PrecompiledContractsHomestead, checkProofAddress)
	delete(vm.PrecompiledContractsHomestead, checkChallengeAddress)
}

// checkChallenge implemented as a native contract to take an on-chain challenge.
type checkChallenge struct{
	blockchain *core.BlockChain
}

func (c *checkChallenge) RequiredGas(_ []byte) uint64 {
	return params.TakeChallengeGas
}

// checkChallenge, take challenge from AC by copy the packed byte array, decode and
// validate it, the on challenge client should send the proof of innocent via a transaction.
func (c *checkChallenge) Run(input []byte) ([]byte, error) {
	if len(input) == 0 {
		panic(fmt.Errorf("invalid proof of innocent - empty"))
	}

	hash, err := c.CheckChallenge(input)
	if err != nil {
		return common.Hash{}.Bytes(), err
	}

	return hash, nil
}

// validate the proof is a valid challenge.
func (c *checkChallenge) validateChallenge(p *types.Proof) ([]byte, error) {
	// check if evidence msgs are from committee members of that height.
	h, err := p.Message.Height()
	if err != nil {
		return nil, err
	}

	header := c.blockchain.GetHeaderByNumber(h.Uint64())
	// validate message.
	if _, err = p.Message.Validate(crypto.CheckValidatorSignature, header); err != nil {
		return nil, err
	}

	for i:=0; i < len(p.Evidence); i++ {
		if _, err = p.Evidence[i].Validate(crypto.CheckValidatorSignature, header); err != nil {
			return nil, err
		}
	}

	// todo: check if the proof is a valid suspicion.
	// if valid, then return msgHash, nil
	return types.RLPHash(p.Message.Payload()).Bytes(), nil
}

// validate challenge, call from EVM package.
func (c *checkChallenge) CheckChallenge(packedProof []byte) ([]byte, error) {
	p, err := decodeProof(packedProof)
	if err != nil {
		return nil, err
	}

	return c.validateChallenge(p)
}

// checkProof implemented as a native contract to validate an on-chain innocent proof.
type checkProof struct{
	blockchain *core.BlockChain
}

func (c *checkProof) RequiredGas(_ []byte) uint64 {
	return params.CheckInnocentGas
}

// checkProof, take proof from AC by copy the packed byte array, decode and validate it.
func (c *checkProof) Run(input []byte) ([]byte, error) {
	// take an on-chain innocent proof, tell the results of the checking
	if len(input) == 0 {
		panic(fmt.Errorf("invalid proof of innocent - empty"))
	}

	hash, err := c.CheckProof(input)
	if err != nil {
		return common.Hash{}.Bytes(), err
	}

	return hash, nil
}

// Check the proof of innocent, it is called from precompiled contracts of EVM package.
func (c *checkProof) CheckProof(packedProof []byte) ([]byte, error) {

	p, err := decodeProof(packedProof)
	if err != nil {
		return nil, err
	}

	return c.validateInnocentProof(p)
}

// validate the innocent proof is valid.
func (c *checkProof) validateInnocentProof(in *types.Proof) ([]byte, error) {
	// check if evidence msgs are from committee members of that height.
	h, err := in.Message.Height()
	if err != nil {
		return nil, err
	}

	header := c.blockchain.GetHeaderByNumber(h.Uint64())
	// validate message.
	if _, err = in.Message.Validate(crypto.CheckValidatorSignature, header); err != nil {
		return nil, err
	}

	for i:=0; i < len(in.Evidence); i++ {
		if _, err = in.Evidence[i].Validate(crypto.CheckValidatorSignature, header); err != nil {
			return nil, err
		}
	}

	// todo: check if the proof is an innocent behavior.
	// return msg hash and nil when its proved as valid.
	return types.RLPHash(in.Message.Payload()).Bytes(), nil
}

func decodeProof(proof []byte) (*types.Proof, error) {
	p := new(types.RawProof)
	err := rlp.DecodeBytes(proof, p)
	if err != nil {
		return nil, err
	}

	decodedP := new(types.Proof)
	decodedP.Rule = p.Rule

	// decode consensus message which is rlp encoded.
	msg := new(types.ConsensusMessage)
	if err := msg.FromPayload(p.Message); err != nil {
		return nil, err
	}
	decodedP.Message = *msg

	for i:= 0; i < len(p.Evidence); i++ {
		m := new(types.ConsensusMessage)
		if err := m.FromPayload(p.Evidence[i]); err != nil {
			return nil, fmt.Errorf("msg cannot be decoded")
		}
		decodedP.Evidence = append(decodedP.Evidence, *m)
	}
	return decodedP, nil
}