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
	failure64Byte = make([]byte, 64)
)

// init the instances of AFD contracts, and register thems into evm's context
func registerAFDContracts(chain *core.BlockChain) {
	pv := ProofValidator{chain: chain}
	cv := ChallengeValidator{chain: chain}

	vm.PrecompiledContractsByzantium[checkProofAddress] = &pv
	vm.PrecompiledContractsByzantium[checkChallengeAddress] = &cv

	vm.PrecompiledContractsHomestead[checkProofAddress] = &pv
	vm.PrecompiledContractsHomestead[checkChallengeAddress] = &cv

	vm.PrecompiledContractsIstanbul[checkProofAddress] = &pv
	vm.PrecompiledContractsIstanbul[checkChallengeAddress] = &cv

	vm.PrecompiledContractsYoloV1[checkProofAddress] = &pv
	vm.PrecompiledContractsYoloV1[checkChallengeAddress] = &cv
}

// un register AFD contracts from evm's context.
func unRegisterAFDContracts() {
	delete(vm.PrecompiledContractsByzantium, checkProofAddress)
	delete(vm.PrecompiledContractsByzantium, checkChallengeAddress)

	delete(vm.PrecompiledContractsYoloV1, checkProofAddress)
	delete(vm.PrecompiledContractsYoloV1, checkChallengeAddress)

	delete(vm.PrecompiledContractsIstanbul, checkProofAddress)
	delete(vm.PrecompiledContractsIstanbul, checkChallengeAddress)

	delete(vm.PrecompiledContractsHomestead, checkProofAddress)
	delete(vm.PrecompiledContractsHomestead, checkChallengeAddress)
}

// ChallengeValidator implemented as a native contract to validate if challenge is valid
type ChallengeValidator struct{
	chain *core.BlockChain
}

// the gas cost to execute ChallengeValidator contract.
func (c *ChallengeValidator) RequiredGas(_ []byte) uint64 {
	return params.TakeChallengeGas
}

// take the rlp encoded proof of challenge in byte array, decode it and validate it, if the proof is validate, then
// the rlp hash of the msg payload and rlp hash of msg sender is returned as the valid identity for proof management.
func (c *ChallengeValidator) Run(input []byte) ([]byte, error) {
	if len(input) == 0 {
		return failure64Byte, fmt.Errorf("invalid input")
	}

	p, err := decodeProof(input)
	if err != nil {
		return failure64Byte, err
	}

	return c.validateChallenge(p)
}

// validate the proof, if the proof is validate, then the rlp hash of the msg payload and rlp hash of msg sender is
// returned as the valid identity for proof management.
func (c *ChallengeValidator) validateChallenge(p *types.Proof) ([]byte, error) {
	if len(p.Evidence) == 0 {
		return failure64Byte, errNoEvidence
	}

	// check if suspicious message is from correct committee member.
	err := checkMsgSignature(c.chain, &p.Message)
	if err != nil {
		return failure64Byte, err
	}

	// check if evidence msgs are from committee members of that height.
	h, err := p.Message.Height()
	if err != nil {
		return failure64Byte, err
	}
	header := c.chain.GetHeaderByNumber(h.Uint64())

	for i:=0; i < len(p.Evidence); i++ {
		if _, err = p.Evidence[i].Validate(crypto.CheckValidatorSignature, header); err != nil {
			return failure64Byte, err
		}
	}

	if c.validEvidence(p) {
		msgHash := types.RLPHash(p.Message.Payload()).Bytes()
		senderHash := types.RLPHash(p.Message.Address).Bytes()
		return append(msgHash, senderHash...), nil
	}
	return failure64Byte, errInvalidChallenge
}

// check if the evidence of the challenge is valid or not.
func (c *ChallengeValidator) validEvidence(p *types.Proof) bool {
	switch types.Rule(p.Rule) {
	case types.PN:
		//todo Validate evidence of PN rule.
	case types.PO:
		//todo Validate evidence of PO rule.
	case types.PVN:
		//todo Validate evidence of PVN rule.
	case types.PVO:
		//todo Validate evidence of PVO rule.
	case types.C:
		//todo Validate evidence of C rule.
	case types.GarbageMessage:
		return preProcessConsensusMsg(c.chain, &p.Message) == errGarbageMsg
	case types.InvalidProposal:
		return preProcessConsensusMsg(c.chain, &p.Message) == errProposal
	case types.InvalidProposer:
		return preProcessConsensusMsg(c.chain, &p.Message) == errProposer
	case types.Equivocation:
		return checkEquivocation(c.chain, &p.Message, p.Evidence) == errEquivocation
	default:
		return false
	}
	return false
}

// ProofValidator implemented as a native contract to validate an on-chain innocent proof.
type ProofValidator struct{
	chain *core.BlockChain
}

// the gas cost to execute this proof validator contract.
func (c *ProofValidator) RequiredGas(_ []byte) uint64 {
	return params.CheckInnocentGas
}

// ProofValidator, take the rlp encoded proof of innocent, decode it and validate it, if the proof is valid, then
// return the rlp hash of msg and the rlp hash of msg sender as the valid identity for on-chain management of proofs,
// AC need the check the value returned to match the ID which is on challenge, to remove the challenge from chain.
func (c *ProofValidator) Run(input []byte) ([]byte, error) {
	// take an on-chain innocent proof, tell the results of the checking
	if len(input) == 0 {
		return failure64Byte, fmt.Errorf("invalid input")
	}

	p, err := decodeProof(input)
	if err != nil {
		return failure64Byte, err
	}

	return c.validateInnocentProof(p)
}

// validate the innocent proof is valid.
func (c *ProofValidator) validateInnocentProof(in *types.Proof) ([]byte, error) {
	// check if evidence msgs are from committee members of that height.
	h, err := in.Message.Height()
	if err != nil {
		return failure64Byte, err
	}

	header := c.chain.GetHeaderByNumber(h.Uint64())
	// validate message.
	if _, err = in.Message.Validate(crypto.CheckValidatorSignature, header); err != nil {
		return failure64Byte, err
	}

	for i:=0; i < len(in.Evidence); i++ {
		if _, err = in.Evidence[i].Validate(crypto.CheckValidatorSignature, header); err != nil {
			return failure64Byte, err
		}
	}

	// todo: check if the proof is an innocent behavior.

	msgHash := types.RLPHash(in.Message.Payload()).Bytes()
	senderHash := types.RLPHash(in.Message.Address).Bytes()
	return append(msgHash, senderHash...), nil
}

// decode proof convert proof from rlp encoded bytes into object Proof.
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