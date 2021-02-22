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
	checkAccusationAddress = common.BytesToAddress([]byte{252})
	checkProofAddress = common.BytesToAddress([]byte{253})
	checkChallengeAddress = common.BytesToAddress([]byte{254})
	failure64Byte = make([]byte, 64)
)

// init the instances of AFD contracts, and register thems into evm's context
func registerAFDContracts(chain *core.BlockChain) {
	pv := InnocentValidator{chain: chain}
	cv := ChallengeValidator{chain: chain}
	av := AccusationValidator{chain: chain}

	vm.PrecompiledContractsByzantium[checkProofAddress] = &pv
	vm.PrecompiledContractsByzantium[checkChallengeAddress] = &cv
	vm.PrecompiledContractsByzantium[checkAccusationAddress] = &av

	vm.PrecompiledContractsHomestead[checkProofAddress] = &pv
	vm.PrecompiledContractsHomestead[checkChallengeAddress] = &cv
	vm.PrecompiledContractsHomestead[checkAccusationAddress] = &av

	vm.PrecompiledContractsIstanbul[checkProofAddress] = &pv
	vm.PrecompiledContractsIstanbul[checkChallengeAddress] = &cv
	vm.PrecompiledContractsIstanbul[checkAccusationAddress] = &av

	vm.PrecompiledContractsYoloV1[checkProofAddress] = &pv
	vm.PrecompiledContractsYoloV1[checkChallengeAddress] = &cv
	vm.PrecompiledContractsYoloV1[checkAccusationAddress] = &av
}

// un register AFD contracts from evm's context.
func unRegisterAFDContracts() {
	delete(vm.PrecompiledContractsByzantium, checkProofAddress)
	delete(vm.PrecompiledContractsByzantium, checkChallengeAddress)
	delete(vm.PrecompiledContractsByzantium, checkAccusationAddress)

	delete(vm.PrecompiledContractsYoloV1, checkProofAddress)
	delete(vm.PrecompiledContractsYoloV1, checkChallengeAddress)
	delete(vm.PrecompiledContractsYoloV1, checkAccusationAddress)

	delete(vm.PrecompiledContractsIstanbul, checkProofAddress)
	delete(vm.PrecompiledContractsIstanbul, checkChallengeAddress)
	delete(vm.PrecompiledContractsIstanbul, checkAccusationAddress)

	delete(vm.PrecompiledContractsHomestead, checkProofAddress)
	delete(vm.PrecompiledContractsHomestead, checkChallengeAddress)
	delete(vm.PrecompiledContractsHomestead, checkAccusationAddress)
}

// AccusationValidator implemented as a native contract to validate if a accusation is valid
type AccusationValidator struct {
	chain *core.BlockChain
}

// the gas cost to execute AccusationValidator contract.
func (a *AccusationValidator) RequiredGas(_ []byte) uint64 {
	return params.MinimumGas
}

// take the rlp encoded proof of accusation in byte array, decode it and validate it, if the proof is validate, then
// the rlp hash of the msg payload and the msg sender is returned.
func (a *AccusationValidator) Run(input []byte) ([]byte, error) {
	if len(input) == 0 {
		return failure64Byte, fmt.Errorf("invalid input")
	}

	p, err := decodeProof(input)
	if err != nil {
		return failure64Byte, err
	}

	return a.validateAccusation(p)
}

// validate if the accusation is valid.
func (a *AccusationValidator) validateAccusation(in *types.Proof) ([]byte, error) {
	// we have only 3 types of rule on accusation.
	switch in.Rule {
	case types.PO:
		if in.Message.Code != types.MsgProposal {
			return failure64Byte, fmt.Errorf("wrong msg for PO rule")
		}
	case types.PVN:
		if in.Message.Code != types.MsgPrevote {
			return failure64Byte, fmt.Errorf("wrong msg for PVN rule")
		}
	case types.C:
		if in.Message.Code != types.MsgPrecommit {
			return failure64Byte, fmt.Errorf("wrong msg for rule C")
		}
	default:
		return failure64Byte, fmt.Errorf("not provable accusation rule")
	}

	// check if the suspicious msg is from the correct committee of that height.
	h, err := in.Message.Height()
	if err != nil {
		return failure64Byte, err
	}

	header := a.chain.GetHeaderByNumber(h.Uint64())
	if _, err = in.Message.Validate(crypto.CheckValidatorSignature, header); err != nil {
		return failure64Byte, err
	}

	msgHash := types.RLPHash(in.Message.Payload()).Bytes()
	sender := common.LeftPadBytes(in.Message.Address.Bytes(), 32)
	return append(sender, msgHash...), nil
}

// ChallengeValidator implemented as a native contract to validate if challenge is valid
type ChallengeValidator struct{
	chain *core.BlockChain
}

// the gas cost to execute ChallengeValidator contract.
func (c *ChallengeValidator) RequiredGas(_ []byte) uint64 {
	return params.MinimumGas
}

// take the rlp encoded proof of challenge in byte array, decode it and validate it, if the proof is validate, then
// the rlp hash of the msg payload and the msg sender is returned as the valid identity for proof management.
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
		return validChallengeOfPN(p)
	case types.PO:
		return validChallengeOfPO(p)
	case types.PVN:
		return validChallengeOfPVN(p)
	case types.C:
		return validChallengeOfC(p)
	case types.GarbageMessage:
		return checkAutoIncriminatingMsg(c.chain, &p.Message) == errGarbageMsg
	case types.InvalidProposal:
		return checkAutoIncriminatingMsg(c.chain, &p.Message) == errProposal
	case types.InvalidProposer:
		return checkAutoIncriminatingMsg(c.chain, &p.Message) == errProposer
	case types.Equivocation:
		return checkEquivocation(c.chain, &p.Message, p.Evidence) == errEquivocation
	default:
		return false
	}
}

// InnocentValidator implemented as a native contract to validate an innocent proof.
type InnocentValidator struct{
	chain *core.BlockChain
}

// the gas cost to execute this proof validator contract.
func (c *InnocentValidator) RequiredGas(_ []byte) uint64 {
	return params.MinimumGas
}

// InnocentValidator, take the rlp encoded proof of innocent, decode it and validate it, if the proof is valid, then
// return the rlp hash of msg and the rlp hash of msg sender as the valid identity for on-chain management of proofs,
// AC need the check the value returned to match the ID which is on challenge, to remove the challenge from chain.
func (c *InnocentValidator) Run(input []byte) ([]byte, error) {
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

// validate if the innocent proof is valid, it returns sender address and msg hash in byte array when proof is valid.
func (c *InnocentValidator) validateInnocentProof(in *types.Proof) ([]byte, error) {
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

	if !c.validInnocentProof(in) {
		return failure64Byte, fmt.Errorf("invalid proof of innocent")
	}

	msgHash := types.RLPHash(in.Message.Payload()).Bytes()
	sender := common.LeftPadBytes(in.Message.Address.Bytes(), 32)
	return append(sender, msgHash...), nil
}

func (c *InnocentValidator) validInnocentProof(p *types.Proof) bool {
	// rule engine only have 3 kind of provable accusation for the time being.
	switch types.Rule(p.Rule) {
	case types.PO:
		return validInnocentProofOfPO(p)
	case types.PVN:
		return validInnocentProofOfPVN(p)
	case types.C:
		return validInnocentProofOfC(p)
	default:
		return false
	}
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

///////////////////////////////////////////////////////////////////////
// validate proof of challenge for rules.
// check if the proof of challenge of PN is valid,
// node propose a new value when there is a proof that it precommit at a different value at previous round.
func validChallengeOfPN(c *types.Proof) bool {
	if len(c.Evidence) == 0 {
		return false
	}

	// should be a new proposal
	proposal := c.Message

	if proposal.Code != types.MsgProposal && proposal.ValidRound() != -1 {
		return false
	}

	preCommit := c.Evidence[0]
	if preCommit.Sender() == proposal.Sender() &&
		preCommit.Type() == types.MsgPrecommit &&
		preCommit.R() < proposal.R() && preCommit.Value() != nilValue {
		return true
	}

	return false
}

// check if the proof of challenge of PO is valid
func validChallengeOfPO(c *types.Proof) bool {
	if len(c.Evidence) == 0 {
		return false
	}
	proposal := c.Message
	// should be an old proposal
	if proposal.Type() != types.MsgProposal && proposal.ValidRound() == -1 {
		return false
	}
	preCommit := c.Evidence[0]

	if preCommit.Type() == types.MsgPrecommit && preCommit.R() == proposal.ValidRound() &&
		preCommit.Sender() == proposal.Sender() && preCommit.Value() != nilValue &&
		preCommit.Value() != proposal.Value() {
		return true
	}

	if preCommit.Type() == types.MsgPrecommit &&
		preCommit.R() > proposal.ValidRound() && preCommit.R() < proposal.R() &&
		preCommit.Sender() == proposal.Sender() &&
		preCommit.Value() != nilValue {
		return true
	}
	return false
}

// check if the proof of challenge of PVN is valid.
func validChallengeOfPVN(c *types.Proof) bool {
	if len(c.Evidence) == 0 {
		return false
	}
	prevote := c.Message
	if !(prevote.Type() == types.MsgPrevote && prevote.Value() != nilValue) {
		return false
	}

	// get corresponding proposal from last slot.
	correspondingProposal := c.Evidence[len(c.Evidence)-1]
	if !(correspondingProposal.Type() == types.MsgProposal && correspondingProposal.Value() == prevote.Value() &&
		correspondingProposal.R() == prevote.R() && correspondingProposal.ValidRound() == -1) {
		return false
	}

	// validate precommit.
	preCommit := c.Evidence[0]
	if preCommit.Type() == types.MsgPrecommit && preCommit.Value() != nilValue &&
		preCommit.Value() != prevote.Value() && prevote.Sender() == preCommit.Sender() &&
		preCommit.R() < prevote.R() {
		return true
	}

	return false
}

// check if the proof of challenge of C is valid.
func validChallengeOfC(c *types.Proof) bool {
	// todo: check challenge of C is valid
	return true
}