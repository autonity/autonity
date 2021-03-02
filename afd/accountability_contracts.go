package afd

import (
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/bft"
	core2 "github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/consensus/tendermint/crypto"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	crypto2 "github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/params"
	"github.com/clearmatics/autonity/rlp"
)

var (
	checkAccusationAddress = common.BytesToAddress([]byte{252})
	checkProofAddress      = common.BytesToAddress([]byte{253})
	checkChallengeAddress  = common.BytesToAddress([]byte{254})
	failure64Byte          = make([]byte, 64)
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
func (a *AccusationValidator) validateAccusation(in *Proof) ([]byte, error) {
	// we have only 3 types of rule on accusation.
	switch in.Rule {
	case PO:
		if in.Message.Code != msgProposal {
			return failure64Byte, fmt.Errorf("wrong msg for PO rule")
		}
	case PVN:
		if in.Message.Code != msgPrevote {
			return failure64Byte, fmt.Errorf("wrong msg for PVN rule")
		}
	case C:
		if in.Message.Code != msgPrecommit {
			return failure64Byte, fmt.Errorf("wrong msg for rule C")
		}
	case C1:
		if in.Message.Code != msgPrecommit {
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
type ChallengeValidator struct {
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
func (c *ChallengeValidator) validateChallenge(p *Proof) ([]byte, error) {
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

	for i := 0; i < len(p.Evidence); i++ {
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
func (c *ChallengeValidator) validEvidence(p *Proof) bool {
	switch p.Rule {
	case PN:
		return c.validChallengeOfPN(p)
	case PO:
		return c.validChallengeOfPO(p)
	case PVN:
		return c.validChallengeOfPVN(p)
	case C:
		return c.validChallengeOfC(p)
	case GarbageMessage:
		return checkAutoIncriminatingMsg(c.chain, &p.Message) == errGarbageMsg
	case InvalidProposal:
		return checkAutoIncriminatingMsg(c.chain, &p.Message) == errProposal
	case InvalidProposer:
		return checkAutoIncriminatingMsg(c.chain, &p.Message) == errProposer
	case Equivocation:
		return checkEquivocation(c.chain, &p.Message, p.Evidence) == errEquivocation
	default:
		return false
	}
}

// InnocentValidator implemented as a native contract to validate an innocent proof.
type InnocentValidator struct {
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
func (c *InnocentValidator) validateInnocentProof(in *Proof) ([]byte, error) {
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

	for i := 0; i < len(in.Evidence); i++ {
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

func (c *InnocentValidator) validInnocentProof(p *Proof) bool {
	// rule engine only have 3 kind of provable accusation for the time being.
	switch p.Rule {
	case PO:
		return c.validInnocentProofOfPO(p)
	case PVN:
		return c.validInnocentProofOfPVN(p)
	case C:
		return c.validInnocentProofOfC(p)
	case C1:
		return c.validInnocentProofOfC1(p)
	default:
		return false
	}
}

// check if the proof of innocent of PO is valid.
func (c *InnocentValidator) validInnocentProofOfPO(p *Proof) bool {
	// check if there is quorum number of prevote at the same value on the same valid round
	proposal := p.Message
	if proposal.Type() != msgProposal {
		return false
	}

	height := proposal.H()
	quorum := bft.Quorum(c.chain.GetHeaderByNumber(height - 1).TotalVotingPower())

	// check quorum prevotes for V at validRound.
	for i := 0; i < len(p.Evidence); i++ {
		if !(p.Evidence[i].Type() == msgPrevote && p.Evidence[i].Value() == proposal.Value() &&
			p.Evidence[i].R() == proposal.ValidRound()) {
			return false
		}
	}

	// check no redundant vote msg in evidence in case of hacking.
	if haveRedundantVotes(p.Evidence) {
		return false
	}

	if powerOfVotes(p.Evidence) < quorum {
		return false
	}
	return true
}

// check if the proof of innocent of PVN is valid.
func (c *InnocentValidator) validInnocentProofOfPVN(p *Proof) bool {
	// check if there is quorum number of prevote at the same value on the same valid round
	preVote := p.Message
	if !(preVote.Type() == msgPrevote && preVote.Value() != nilValue) {
		return false
	}

	if len(p.Evidence) == 0 {
		return false
	}

	proposal := p.Evidence[0]
	return proposal.Type() == msgProposal && proposal.Value() == preVote.Value() &&
		proposal.R() == preVote.R()
}

// check if the proof of innocent of C is valid.
func (c *InnocentValidator) validInnocentProofOfC(p *Proof) bool {
	preCommit := p.Message
	if !(preCommit.Type() == msgPrecommit && preCommit.Value() != nilValue) {
		return false
	}

	if len(p.Evidence) == 0 {
		return false
	}

	proposal := p.Evidence[0]
	return proposal.Type() == msgProposal && proposal.Value() == preCommit.Value() &&
		proposal.R() == preCommit.R()
}

// check if the proof of innocent of C is valid.
func (c *InnocentValidator) validInnocentProofOfC1(p *Proof) bool {
	preCommit := p.Message
	if !(preCommit.Type() == msgPrecommit && preCommit.Value() != nilValue) {
		return false
	}

	height := preCommit.H()
	quorum := bft.Quorum(c.chain.GetHeaderByNumber(height - 1).TotalVotingPower())

	// check quorum prevotes for V at the same round.
	for i := 0; i < len(p.Evidence); i++ {
		if !(p.Evidence[i].Type() == msgPrevote && p.Evidence[i].Value() == preCommit.Value() &&
			p.Evidence[i].R() == preCommit.R()) {
			return false
		}
	}

	// check no redundant vote msg in evidence in case of hacking.
	if haveRedundantVotes(p.Evidence) {
		return false
	}

	if powerOfVotes(p.Evidence) < quorum {
		return false
	}
	return true
}

// decode proof convert proof from rlp encoded bytes into object Proof.
func decodeProof(proof []byte) (*Proof, error) {
	p := new(RawProof)
	err := rlp.DecodeBytes(proof, p)
	if err != nil {
		return nil, err
	}

	decodedP := new(Proof)
	decodedP.Rule = p.Rule

	// decode consensus message which is rlp encoded.
	msg := new(core2.Message)
	if err := msg.FromPayload(p.Message); err != nil {
		return nil, err
	}
	decodedP.Message = *msg

	for i := 0; i < len(p.Evidence); i++ {
		m := new(core2.Message)
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
func (c *ChallengeValidator) validChallengeOfPN(p *Proof) bool {
	if len(p.Evidence) == 0 {
		return false
	}

	// should be a new proposal
	proposal := p.Message

	if proposal.Code != msgProposal && proposal.ValidRound() != -1 {
		return false
	}

	preCommit := p.Evidence[0]
	if preCommit.Sender() == proposal.Sender() &&
		preCommit.Type() == msgPrecommit &&
		preCommit.R() < proposal.R() && preCommit.Value() != nilValue {
		return true
	}

	return false
}

// check if the proof of challenge of PO is valid
func (c *ChallengeValidator) validChallengeOfPO(p *Proof) bool {
	if len(p.Evidence) == 0 {
		return false
	}
	proposal := p.Message
	// should be an old proposal
	if proposal.Type() != msgProposal && proposal.ValidRound() == -1 {
		return false
	}
	preCommit := p.Evidence[0]

	if preCommit.Type() == msgPrecommit && preCommit.R() == proposal.ValidRound() &&
		preCommit.Sender() == proposal.Sender() && preCommit.Value() != nilValue &&
		preCommit.Value() != proposal.Value() {
		return true
	}

	if preCommit.Type() == msgPrecommit &&
		preCommit.R() > proposal.ValidRound() && preCommit.R() < proposal.R() &&
		preCommit.Sender() == proposal.Sender() &&
		preCommit.Value() != nilValue {
		return true
	}
	return false
}

// check if the proof of challenge of PVN is valid.
func (c *ChallengeValidator) validChallengeOfPVN(p *Proof) bool {
	if len(p.Evidence) == 0 {
		return false
	}
	prevote := p.Message
	if !(prevote.Type() == msgPrevote && prevote.Value() != nilValue) {
		return false
	}

	// get corresponding proposal from last slot.
	correspondingProposal := p.Evidence[len(p.Evidence)-1]
	if !(correspondingProposal.Type() == msgProposal && correspondingProposal.Value() == prevote.Value() &&
		correspondingProposal.R() == prevote.R() && correspondingProposal.ValidRound() == -1) {
		return false
	}

	// validate precommit.
	preCommit := p.Evidence[0]
	if preCommit.Type() == msgPrecommit && preCommit.Value() != nilValue &&
		preCommit.Value() != prevote.Value() && prevote.Sender() == preCommit.Sender() &&
		preCommit.R() < prevote.R() {
		return true
	}

	return false
}

// check if the proof of challenge of C is valid.
func (c *ChallengeValidator) validChallengeOfC(p *Proof) bool {
	if len(p.Evidence) == 0 {
		return false
	}
	preCommit := p.Message
	if !(preCommit.Type() == msgPrecommit && preCommit.Value() != nilValue) {
		return false
	}

	// check prevotes for not the same V of precommit.
	for i := 0; i < len(p.Evidence); i++ {
		if !(p.Evidence[i].Type() == msgPrevote && p.Evidence[i].Value() != preCommit.Value() &&
			p.Evidence[i].R() == preCommit.R()) {
			return false
		}
	}

	// check no redundant vote msg in evidence in case of hacking.
	if haveRedundantVotes(p.Evidence) {
		return false
	}

	// check if prevotes for not V reaches to quorum.
	quorum := bft.Quorum(c.chain.GetHeaderByNumber(p.Message.H() - 1).TotalVotingPower())
	if powerOfVotes(p.Evidence) >= quorum {
		return true
	}

	return true
}

func haveRedundantVotes(votes []core2.Message) bool {
	voteMap := make(map[common.Hash]struct{})
	for _, vote := range votes {
		hash := common.BytesToHash(crypto2.Keccak256(vote.Payload()))
		_, ok := voteMap[hash]
		if !ok {
			voteMap[hash] = struct{}{}
		} else {
			return true
		}
	}

	return false
}
