package faultdetector

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/bft"
	tendermintCore "github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/consensus/tendermint/crypto"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/params"
	"github.com/clearmatics/autonity/rlp"
)

var (
	checkAccusationAddress   = common.BytesToAddress([]byte{252})
	checkInnocenceAddress    = common.BytesToAddress([]byte{253})
	checkMisbehaviourAddress = common.BytesToAddress([]byte{254})
	// error codes of the execution of precompiled contract to verify the input Proof.
	validProofByte = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	failure96Byte  = make([]byte, 96)
)

// init the instances of Fault Detector contracts, and register thems into evm's context
func registerFaultDetectorContracts(chain BlockChainContext) {
	vm.PrecompiledContractRWMutex.Lock()
	defer vm.PrecompiledContractRWMutex.Unlock()
	pv := InnocenceVerifier{chain: chain}
	cv := MisbehaviourVerifier{chain: chain}
	av := AccusationVerifier{chain: chain}

	vm.PrecompiledContractsByzantium[checkInnocenceAddress] = &pv
	vm.PrecompiledContractsByzantium[checkMisbehaviourAddress] = &cv
	vm.PrecompiledContractsByzantium[checkAccusationAddress] = &av

	vm.PrecompiledContractsHomestead[checkInnocenceAddress] = &pv
	vm.PrecompiledContractsHomestead[checkMisbehaviourAddress] = &cv
	vm.PrecompiledContractsHomestead[checkAccusationAddress] = &av

	vm.PrecompiledContractsIstanbul[checkInnocenceAddress] = &pv
	vm.PrecompiledContractsIstanbul[checkMisbehaviourAddress] = &cv
	vm.PrecompiledContractsIstanbul[checkAccusationAddress] = &av

	vm.PrecompiledContractsYoloV1[checkInnocenceAddress] = &pv
	vm.PrecompiledContractsYoloV1[checkMisbehaviourAddress] = &cv
	vm.PrecompiledContractsYoloV1[checkAccusationAddress] = &av
}

// unregister Faul Detector contracts from evm's context.
func unRegisterFaultDetectorContracts() {
	vm.PrecompiledContractRWMutex.Lock()
	defer vm.PrecompiledContractRWMutex.Unlock()

	delete(vm.PrecompiledContractsByzantium, checkInnocenceAddress)
	delete(vm.PrecompiledContractsByzantium, checkMisbehaviourAddress)
	delete(vm.PrecompiledContractsByzantium, checkAccusationAddress)

	delete(vm.PrecompiledContractsYoloV1, checkInnocenceAddress)
	delete(vm.PrecompiledContractsYoloV1, checkMisbehaviourAddress)
	delete(vm.PrecompiledContractsYoloV1, checkAccusationAddress)

	delete(vm.PrecompiledContractsIstanbul, checkInnocenceAddress)
	delete(vm.PrecompiledContractsIstanbul, checkMisbehaviourAddress)
	delete(vm.PrecompiledContractsIstanbul, checkAccusationAddress)

	delete(vm.PrecompiledContractsHomestead, checkInnocenceAddress)
	delete(vm.PrecompiledContractsHomestead, checkMisbehaviourAddress)
	delete(vm.PrecompiledContractsHomestead, checkAccusationAddress)
}

// AccusationVerifier implemented as a native contract to validate if a accusation is valid
type AccusationVerifier struct {
	chain BlockChainContext
}

// the gas cost to execute AccusationVerifier contract.
func (a *AccusationVerifier) RequiredGas(_ []byte) uint64 {
	return params.AutonityPrecompiledContractGas
}

// take the rlp encoded Proof of accusation in byte array, decode it and validate it, if the Proof is validate, then
// the rlp hash of the msg payload and the msg sender is returned.
func (a *AccusationVerifier) Run(input []byte) ([]byte, error) {
	if len(input) <= 32 {
		return failure96Byte, nil
	}

	// the 1st 32 bytes are length of bytes array in solidity, take RLP bytes after it.
	p, err := decodeRawProof(input[32:])
	if err != nil {
		return failure96Byte, nil
	}

	return a.validateAccusation(p), nil
}

// validate if the accusation is valid.
func (a *AccusationVerifier) validateAccusation(in *Proof) []byte {
	// we have only 4 types of rule on accusation.
	switch in.Rule {
	case PO:
		if in.Message.Code != msgProposal {
			return failure96Byte
		}
	case PVN:
		if in.Message.Code != msgPrevote {
			return failure96Byte
		}
	case PVO:
		if in.Message.Code != msgPrevote {
			return failure96Byte
		}
	case C:
		if in.Message.Code != msgPrecommit {
			return failure96Byte
		}
	case C1:
		if in.Message.Code != msgPrecommit {
			return failure96Byte
		}
	default:
		return failure96Byte
	}

	// check if the suspicious msg is from the correct committee of that height.
	h, err := in.Message.Height()
	if err != nil {
		return failure96Byte
	}

	lastHeader := a.chain.GetHeaderByNumber(h.Uint64() - 1)
	if lastHeader == nil {
		return failure96Byte
	}
	if _, err = in.Message.Validate(crypto.CheckValidatorSignature, lastHeader); err != nil {
		return failure96Byte
	}

	// in case of PVO accusation, we need to check corresponding old proposal of this preVote.
	if in.Rule == PVO {
		if len(in.Evidence) == 0 {
			return failure96Byte
		}

		oldProposal := in.Evidence[0]
		if _, err = oldProposal.Validate(crypto.CheckValidatorSignature, lastHeader); err != nil {
			return failure96Byte
		}

		if oldProposal.Code != msgProposal || oldProposal.R() != in.Message.R() ||
			oldProposal.Value() != in.Message.Value() || oldProposal.ValidRound() == -1 {
			return failure96Byte
		}
	}

	msgHash := in.Message.MsgHash().Bytes()
	sender := common.LeftPadBytes(in.Message.Address.Bytes(), 32)
	return append(append(sender, msgHash...), validProofByte...)
}

// MisbehaviourVerifier implemented as a native contract to validate if misbehaviour is valid
type MisbehaviourVerifier struct {
	chain BlockChainContext
}

// the gas cost to execute MisbehaviourVerifier contract.
func (c *MisbehaviourVerifier) RequiredGas(_ []byte) uint64 {
	return params.AutonityPrecompiledContractGas
}

// take the rlp encoded Proof of challenge in byte array, decode it and validate it, if the Proof is validate, then
// the rlp hash of the msg payload and the msg sender is returned as the valid identity for Proof management.
func (c *MisbehaviourVerifier) Run(input []byte) ([]byte, error) {
	if len(input) <= 32 {
		return failure96Byte, nil
	}

	// the 1st 32 bytes are length of bytes array in solidity, take RLP bytes after it.
	p, err := decodeRawProof(input[32:])
	if err != nil {
		return failure96Byte, nil
	}

	return c.validateProof(p), nil
}

// validate the Proof, if the Proof is validate, then the rlp hash of the msg payload and rlp hash of msg sender is
// returned as the valid identity for Proof management.
func (c *MisbehaviourVerifier) validateProof(p *Proof) []byte {

	// check if suspicious message is from correct committee member.
	err := checkMsgSignature(c.chain, p.Message)
	if err != nil {
		return failure96Byte
	}

	// check if evidence msgs are from committee members of that height.
	h, err := p.Message.Height()
	if err != nil {
		return failure96Byte
	}
	lastHeader := c.chain.GetHeaderByNumber(h.Uint64() - 1)
	if lastHeader == nil {
		return failure96Byte
	}

	for i := 0; i < len(p.Evidence); i++ {
		// the height of msg of the evidences is checked at Validate function.
		if _, err = p.Evidence[i].Validate(crypto.CheckValidatorSignature, lastHeader); err != nil {
			return failure96Byte
		}
	}

	if c.validProof(p) {
		msgHash := p.Message.MsgHash().Bytes()
		sender := common.LeftPadBytes(p.Message.Address.Bytes(), 32)
		return append(append(sender, msgHash...), validProofByte...)
	}
	return failure96Byte
}

// check if the evidence of the misbehaviour is valid or not.
func (c *MisbehaviourVerifier) validProof(p *Proof) bool {
	switch p.Rule {
	case PN:
		return c.validMisbehaviourOfPN(p)
	case PO:
		return c.validMisbehaviourOfPO(p)
	case PVN:
		return c.validMisbehaviourOfPVN(p)
	case C:
		return c.validMisbehaviourOfC(p)
	case GarbageMessage:
		return checkAutoIncriminatingMsg(c.chain, p.Message) == errGarbageMsg
	case InvalidProposal:
		return checkAutoIncriminatingMsg(c.chain, p.Message) == errProposal
	case InvalidProposer:
		return checkAutoIncriminatingMsg(c.chain, p.Message) == errProposer
	case Equivocation:
		return checkEquivocation(c.chain, p.Message, p.Evidence) == errEquivocation
	default:
		return false
	}
}

// check if the Proof of challenge of PN is valid,
// node propose a new value when there is a Proof that it precommit at a different value at previous round.
func (c *MisbehaviourVerifier) validMisbehaviourOfPN(p *Proof) bool {
	if len(p.Evidence) == 0 {
		return false
	}

	// should be a new proposal
	proposal := p.Message

	if proposal.Code != msgProposal || proposal.ValidRound() != -1 {
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

// check if the Proof of challenge of PO is valid
func (c *MisbehaviourVerifier) validMisbehaviourOfPO(p *Proof) bool {
	if len(p.Evidence) == 0 {
		return false
	}
	proposal := p.Message
	// should be an old proposal
	if proposal.Type() != msgProposal || proposal.ValidRound() == -1 {
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

// check if the Proof of challenge of PVN is valid.
func (c *MisbehaviourVerifier) validMisbehaviourOfPVN(p *Proof) bool {
	if len(p.Evidence) == 0 {
		return false
	}
	prevote := p.Message
	if prevote.Type() != msgPrevote || prevote.Value() == nilValue {
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

// check if the Proof of challenge of C is valid.
func (c *MisbehaviourVerifier) validMisbehaviourOfC(p *Proof) bool {
	if len(p.Evidence) == 0 {
		return false
	}
	preCommit := p.Message
	if preCommit.Type() != msgPrecommit || preCommit.Value() == nilValue {
		return false
	}

	// check preVotes for not the same V compares to preCommit.
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

	// check if preVotes for not V reaches to quorum.
	quorum := bft.Quorum(c.chain.GetHeaderByNumber(p.Message.H() - 1).TotalVotingPower())
	return powerOfVotes(p.Evidence) >= quorum
}

// InnocenceVerifier implemented as a native contract to validate an innocence Proof.
type InnocenceVerifier struct {
	chain BlockChainContext
}

// the gas cost to execute this Proof validator contract.
func (c *InnocenceVerifier) RequiredGas(_ []byte) uint64 {
	return params.AutonityPrecompiledContractGas
}

// InnocenceVerifier, take the rlp encoded Proof of innocence, decode it and validate it, if the Proof is valid, then
// return the rlp hash of msg and the rlp hash of msg sender as the valid identity for on-chain management of proofs,
// AC need the check the value returned to match the ID which is on challenge, to remove the challenge from chain.
func (c *InnocenceVerifier) Run(input []byte) ([]byte, error) {
	// take an on-chain innocent Proof, tell the results of the checking
	if len(input) <= 32 {
		return failure96Byte, nil
	}

	// the 1st 32 bytes are length of bytes array in solidity, take RLP bytes after it.
	p, err := decodeRawProof(input[32:])
	if err != nil {
		return failure96Byte, nil
	}

	return c.validateInnocenceProof(p), nil
}

// validate if the innocence Proof is valid, it returns sender address and msg hash in byte array when Proof is valid.
func (c *InnocenceVerifier) validateInnocenceProof(in *Proof) []byte {
	// check if evidence msgs are from committee members of that height.
	h, err := in.Message.Height()
	if err != nil {
		return failure96Byte
	}

	lastHeader := c.chain.GetHeaderByNumber(h.Uint64() - 1)
	if lastHeader == nil {
		return failure96Byte
	}

	// validate message.
	if _, err = in.Message.Validate(crypto.CheckValidatorSignature, lastHeader); err != nil {
		return failure96Byte
	}

	for i := 0; i < len(in.Evidence); i++ {
		// the height of msg of the evidences is checked at Validate function.
		if _, err = in.Evidence[i].Validate(crypto.CheckValidatorSignature, lastHeader); err != nil {
			return failure96Byte
		}
	}

	if !c.validInnocenceProof(in) {
		return failure96Byte
	}

	msgHash := in.Message.MsgHash().Bytes()
	sender := common.LeftPadBytes(in.Message.Address.Bytes(), 32)
	return append(append(sender, msgHash...), validProofByte...)
}

func (c *InnocenceVerifier) validInnocenceProof(p *Proof) bool {
	// rule engine only have 3 kind of provable accusation for the time being.
	switch p.Rule {
	case PO:
		return c.validInnocenceProofOfPO(p)
	case PVN:
		return c.validInnocenceProofOfPVN(p)
	case C:
		return c.validInnocenceProofOfC(p)
	case C1:
		return c.validInnocenceProofOfC1(p)
	default:
		return false
	}
}

// check if the Proof of innocent of PO is valid.
func (c *InnocenceVerifier) validInnocenceProofOfPO(p *Proof) bool {
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

// check if the Proof of innocent of PVN is valid.
func (c *InnocenceVerifier) validInnocenceProofOfPVN(p *Proof) bool {
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

// check if the Proof of innocent of C is valid.
func (c *InnocenceVerifier) validInnocenceProofOfC(p *Proof) bool {
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

// check if the Proof of innocent of C is valid.
func (c *InnocenceVerifier) validInnocenceProofOfC1(p *Proof) bool {
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

func haveRedundantVotes(votes []*tendermintCore.Message) bool {
	voteMap := make(map[common.Address]struct{})
	for _, vote := range votes {
		_, ok := voteMap[vote.Address]
		if !ok {
			voteMap[vote.Address] = struct{}{}
		} else {
			return true
		}
	}

	return false
}

// decode Proof convert Proof from rlp encoded bytes into object Proof.
func decodeRawProof(b []byte) (*Proof, error) {
	p := new(Proof)
	err := rlp.DecodeBytes(b, p)
	if err != nil {
		return nil, err
	}

	// decode the msg.
	err = decodeMessage(p.Message)
	if err != nil {
		return nil, err
	}

	// decode the evidence.
	for _, m := range p.Evidence {
		err = decodeMessage(m)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func decodeMessage(m *tendermintCore.Message) error {
	switch m.Code {
	case msgProposal:
		var proposal tendermintCore.Proposal
		return m.Decode(&proposal)
	case msgPrevote, msgPrecommit:
		var vote tendermintCore.Vote
		return m.Decode(&vote)
	default:
		return errGarbageMsg
	}
}
