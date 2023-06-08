package misbehaviourdetector

import (
	"fmt"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	proto "github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	engineCore "github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	mUtils "github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"github.com/autonity/autonity/consensus/tendermint/crypto"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/rlp"
)

// precompiled contracts to be call by autonity contract to verify on-chain proofs of accountability events, they are
// a part of consensus.

var (
	checkAccusationAddress   = common.BytesToAddress([]byte{252})
	checkInnocenceAddress    = common.BytesToAddress([]byte{253})
	checkMisbehaviourAddress = common.BytesToAddress([]byte{254})
	// error codes of the execution of precompiled contract to verify the input AccountabilityProof.
	validByte      = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	failure128Byte = make([]byte, 128)
)

const KB = 1024

// LoadAccountabilityPreCompiledContracts init the instances of Fault Detector contracts, and register them into EVM's context
func LoadAccountabilityPreCompiledContracts(chain BlockChainContext) {

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

	vm.PrecompiledContractsBerlin[checkInnocenceAddress] = &pv
	vm.PrecompiledContractsBerlin[checkMisbehaviourAddress] = &cv
	vm.PrecompiledContractsBerlin[checkAccusationAddress] = &av

	vm.PrecompiledContractsBLS[checkInnocenceAddress] = &pv
	vm.PrecompiledContractsBLS[checkMisbehaviourAddress] = &cv
	vm.PrecompiledContractsBLS[checkAccusationAddress] = &av
}

// unregister Fault Detector contracts from EVM's context.
func unRegisterFaultDetectorContracts() {
	vm.PrecompiledContractRWMutex.Lock()
	defer vm.PrecompiledContractRWMutex.Unlock()

	delete(vm.PrecompiledContractsByzantium, checkInnocenceAddress)
	delete(vm.PrecompiledContractsByzantium, checkMisbehaviourAddress)
	delete(vm.PrecompiledContractsByzantium, checkAccusationAddress)

	delete(vm.PrecompiledContractsBerlin, checkInnocenceAddress)
	delete(vm.PrecompiledContractsBerlin, checkMisbehaviourAddress)
	delete(vm.PrecompiledContractsBerlin, checkAccusationAddress)

	delete(vm.PrecompiledContractsBLS, checkInnocenceAddress)
	delete(vm.PrecompiledContractsBLS, checkMisbehaviourAddress)
	delete(vm.PrecompiledContractsBLS, checkAccusationAddress)

	delete(vm.PrecompiledContractsIstanbul, checkInnocenceAddress)
	delete(vm.PrecompiledContractsIstanbul, checkMisbehaviourAddress)
	delete(vm.PrecompiledContractsIstanbul, checkAccusationAddress)

	delete(vm.PrecompiledContractsHomestead, checkInnocenceAddress)
	delete(vm.PrecompiledContractsHomestead, checkMisbehaviourAddress)
	delete(vm.PrecompiledContractsHomestead, checkAccusationAddress)
}

// AccusationVerifier implemented as a native contract to validate if an accusation is valid
type AccusationVerifier struct {
	chain BlockChainContext
}

// RequiredGas the gas cost to execute AccusationVerifier contract, weighted by input data size.
func (a *AccusationVerifier) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

// Run take the rlp encoded AccountabilityProof of accusation in byte array, decode it and validate it, if the AccountabilityProof is valid, then
// the rlp hash of the msg payload and the msg sender is returned.
func (a *AccusationVerifier) Run(input []byte, blockNumber uint64) ([]byte, error) {
	if len(input) <= 32 {
		return failure128Byte, nil
	}

	// the 1st 32 bytes are length of bytes array in solidity, take RLP bytes after it.
	p, err := decodeRawProof(input[32:])
	if err != nil {
		return failure128Byte, nil
	}

	evHeight, err := p.Message.Height()
	if err != nil || evHeight.Uint64() == 0 {
		return failure128Byte, nil
	}

	// prevent a potential attack: a malicious fault detector can rise an accusation that contain a message
	// corresponding to an old height, while at each Pi, they only buffer specific heights of message in msg store, thus
	// Pi can never provide a valid proof of innocence anymore, making the malicious accusation be valid for slashing.
	if blockNumber > evHeight.Uint64() && (blockNumber-evHeight.Uint64() >= proto.AccountabilityHeightRange) {
		return failure128Byte, nil
	}

	return a.validateAccusation(p), nil
}

// validate the submitted accusation by the contract call.
func validAccusation(chain BlockChainContext, in *AccountabilityProof) bool {
	// we have only 4 types of rule on accusation.
	switch in.Rule {
	case autonity.PO:
		if in.Message.Code != proto.MsgLiteProposal || in.Message.ValidRound() == -1 {
			return false
		}
	case autonity.PVN:
		if in.Message.Code != proto.MsgPrevote {
			return false
		}
	case autonity.PVO:
		if in.Message.Code != proto.MsgPrevote {
			return false
		}
	case autonity.C1:
		if in.Message.Code != proto.MsgPrecommit {
			return false
		}
	default:
		return false
	}

	// check if the suspicious msg is from the correct committee of that height.
	h, err := in.Message.Height()
	if err != nil {
		return false
	}

	lastHeader := chain.GetHeaderByNumber(h.Uint64() - 1)
	if lastHeader == nil {
		return false
	}
	if _, err = in.Message.Validate(crypto.CheckValidatorSignature, lastHeader); err != nil {
		return false
	}

	// in case of PVO accusation, we need to check corresponding old proposal of this preVote.
	if in.Rule == autonity.PVO {
		if len(in.Evidence) != 1 {
			return false
		}

		oldProposal := in.Evidence[0]
		if _, err = oldProposal.Validate(crypto.CheckValidatorSignature, lastHeader); err != nil {
			return false
		}

		if oldProposal.Code != proto.MsgLiteProposal || oldProposal.R() != in.Message.R() ||
			oldProposal.Value() != in.Message.Value() || oldProposal.ValidRound() == -1 {
			return false
		}
	}

	return true
}

// validate if the accusation is valid and return the output bytes of the static call from evm.
func (a *AccusationVerifier) validateAccusation(in *AccountabilityProof) []byte {
	if validAccusation(a.chain, in) {
		return resultBytes(in.Message.MsgHash(), in.Message.Address, in.Rule)
	}
	return failure128Byte
}

// MisbehaviourVerifier implemented as a native contract to validate if misbehaviour is valid
type MisbehaviourVerifier struct {
	chain BlockChainContext
}

// RequiredGas the gas cost to execute MisbehaviourVerifier contract, weighted by input data size.
func (c *MisbehaviourVerifier) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

// Run take the rlp encoded AccountabilityProof of challenge in byte array, decode it and validate it, if the AccountabilityProof is valid, then
// the rlp hash of the msg payload and the msg sender is returned as the valid identity for AccountabilityProof management.
func (c *MisbehaviourVerifier) Run(input []byte, blockNumber uint64) ([]byte, error) {
	if len(input) <= 32 {
		return failure128Byte, nil
	}

	// the 1st 32 bytes are length of bytes array in solidity, take RLP bytes after it.
	p, err := decodeRawProof(input[32:])
	if err != nil {
		if p.Rule == autonity.AccountableGarbageMessage && err == errAccountableGarbageMsg && len(p.Evidence) == 0 {
			// since garbage message cannot be decoded for msg height, we return 0 for msg height in such case.
			return resultBytes(p.Message.MsgHash(), p.Message.Address, p.Rule), nil
		}
		return failure128Byte, nil
	}

	evHeight, err := p.Message.Height()
	if err != nil || evHeight.Uint64() == 0 {
		return failure128Byte, nil
	}

	return c.validateProof(p), nil
}

// validate the AccountabilityProof, if the AccountabilityProof is valid, then the rlp hash of the msg payload and rlp hash of msg sender is
// returned as the valid identity for AccountabilityProof management.
func (c *MisbehaviourVerifier) validateProof(p *AccountabilityProof) []byte {

	// check if suspicious message is from correct committee member.
	err := checkMsgSignature(c.chain, p.Message)
	if err != nil {
		return failure128Byte
	}

	// check if evidence msgs are from committee members of that height.
	h, err := p.Message.Height()
	if err != nil {
		return failure128Byte
	}
	lastHeader := c.chain.GetHeaderByNumber(h.Uint64() - 1)
	if lastHeader == nil {
		return failure128Byte
	}

	// check if the number of evidence msgs are exceeded the max to prevent the abuse of the proof msg.
	if len(p.Evidence) > maxEvidenceMessages(lastHeader) {
		return failure128Byte
	}

	for _, msg := range p.Evidence {
		// the height of msg of the evidences is checked at Validate function.
		if _, err := msg.Validate(crypto.CheckValidatorSignature, lastHeader); err != nil {
			return failure128Byte
		}
	}

	if c.validProof(p) {
		return resultBytes(p.Message.MsgHash(), p.Message.Address, p.Rule)
	}
	return failure128Byte
}

// check if the evidence of the misbehaviour is valid or not.
func (c *MisbehaviourVerifier) validProof(p *AccountabilityProof) bool {
	switch p.Rule {
	case autonity.PN:
		return c.validMisbehaviourOfPN(p)
	case autonity.PO:
		return c.validMisbehaviourOfPO(p)
	case autonity.PVN:
		return c.validMisbehaviourOfPVN(p)
	case autonity.PVO:
		return c.validMisbehaviourOfPVO(p)
	case autonity.PVO12:
		return c.validMisbehaviourOfPVO12(p)
	case autonity.PVO3:
		return c.validMisbehaviourOfPVO3(p)
	case autonity.C:
		return c.validMisbehaviourOfC(p)
	case autonity.InvalidRound:
		if p.Message.R() > constants.MaxRound {
			return true
		}
	case autonity.WrongValidRound:
		if p.Message.Type() != proto.MsgLiteProposal {
			return false
		}
		var proposal mUtils.LiteProposal
		err := p.Message.Decode(&proposal)
		if err != nil {
			return false
		}
		return proposal.ValidRound >= proposal.Round
	case autonity.InvalidProposal:
		return validProofOfBadProposal(c.chain, p)
	case autonity.InvalidProposer:
		return !isProposerMsg(c.chain, p.Message)
	case autonity.Equivocation:
		return checkEquivocation(p.Message, p.Evidence) == errEquivocation
	default:
		return false
	}
	return false
}

// check if the AccountabilityProof of challenge of PN is valid,
// node propose a new value when there is a AccountabilityProof that it precommit at a different value at previous round.
func (c *MisbehaviourVerifier) validMisbehaviourOfPN(p *AccountabilityProof) bool {
	if len(p.Evidence) != 1 {
		return false
	}

	// should be a new proposal
	proposal := p.Message

	if proposal.Code != proto.MsgLiteProposal || proposal.ValidRound() != -1 {
		return false
	}

	preCommit := p.Evidence[0]
	if preCommit.Sender() == proposal.Sender() &&
		preCommit.Type() == proto.MsgPrecommit &&
		preCommit.R() < proposal.R() && preCommit.Value() != nilValue {
		return true
	}

	return false
}

// check if the AccountabilityProof of challenge of PO is valid
func (c *MisbehaviourVerifier) validMisbehaviourOfPO(p *AccountabilityProof) bool {
	proposal := p.Message
	// should be an old proposal
	if proposal.Type() != proto.MsgLiteProposal || proposal.ValidRound() == -1 {
		return false
	}

	// if the proposal contains an invalid valid round, then it is a valid proof.
	if proposal.ValidRound() >= proposal.R() {
		return true
	}

	if len(p.Evidence) == 0 {
		return false
	}
	preCommit := p.Evidence[0]

	if preCommit.Type() == proto.MsgPrecommit && preCommit.R() == proposal.ValidRound() &&
		preCommit.Sender() == proposal.Sender() && preCommit.Value() != nilValue &&
		preCommit.Value() != proposal.Value() {
		return true
	}

	if preCommit.Type() == proto.MsgPrecommit &&
		preCommit.R() > proposal.ValidRound() && preCommit.R() < proposal.R() &&
		preCommit.Sender() == proposal.Sender() &&
		preCommit.Value() != nilValue {
		return true
	}

	// check if there are quorum prevotes for other value than the proposed value at valid round.
	preVote := p.Evidence[0]
	if preVote.Type() == proto.MsgPrevote {
		// validate evidences
		for _, pv := range p.Evidence {
			if pv.Type() != proto.MsgPrevote || pv.R() != proposal.ValidRound() || pv.Value() == proposal.Value() {
				return false
			}
		}

		if haveRedundantVotes(p.Evidence) {
			return false
		}

		// check if preVotes for not V reaches to quorum.
		lastHeader := c.chain.GetHeaderByNumber(p.Message.H() - 1)
		if lastHeader == nil {
			return false
		}
		quorum := bft.Quorum(lastHeader.TotalVotingPower())
		return engineCore.OverQuorumVotes(p.Evidence, quorum.Uint64()) != nil
	}

	return false
}

// check if the AccountabilityProof of challenge of PVN is valid.
func (c *MisbehaviourVerifier) validMisbehaviourOfPVN(p *AccountabilityProof) bool {
	if len(p.Evidence) == 0 {
		return false
	}
	prevote := p.Message
	if prevote.Type() != proto.MsgPrevote || prevote.Value() == nilValue {
		return false
	}

	// check if the corresponding new proposal of preVote for new value is presented.
	correspondingProposal := p.Evidence[0]
	if correspondingProposal.Type() != proto.MsgLiteProposal || correspondingProposal.H() != prevote.H() ||
		correspondingProposal.R() != prevote.R() || correspondingProposal.Value() != prevote.Value() ||
		correspondingProposal.ValidRound() != -1 {
		return false
	}

	// validate preCommits from round R' to R, make sure there is no gap, and the value preCommitted at R'
	// is different from the value preVoted, and the other ones are preCommits of nil.
	preCommits := p.Evidence[1:]

	lastIndex := len(preCommits) - 1

	for i, pc := range preCommits {
		if pc.Type() != proto.MsgPrecommit || pc.Sender() != prevote.Sender() || pc.R() >= prevote.R() {
			return false
		}

		// preCommit at R'
		if i == 0 {
			if pc.Value() == nilValue || pc.Value() == prevote.Value() {
				return false
			}
		} else {
			// preCommits at between R' and R-1, they should be nil.
			if pc.Value() != nilValue {
				return false
			}
		}

		// check if there is round gaps between R' and R-1.
		if i < lastIndex && preCommits[i+1].R()-pc.R() > 1 {
			return false
		}

		// check round gap for preCommit at R-1 and R.
		if i == lastIndex {
			return pc.R()+1 == prevote.R()
		}
	}

	return false
}

// check if the proof of challenge of PVO is valid.
func (c *MisbehaviourVerifier) validMisbehaviourOfPVO(p *AccountabilityProof) bool {
	if len(p.Evidence) < 2 {
		return false
	}
	prevote := p.Message
	if prevote.Type() != proto.MsgPrevote || prevote.Value() == nilValue {
		return false
	}
	// check if the corresponding proposal of preVote is presented.
	correspondingProposal := p.Evidence[0]
	if correspondingProposal.Type() != proto.MsgLiteProposal || correspondingProposal.H() != prevote.H() || correspondingProposal.R() != prevote.R() ||
		correspondingProposal.Value() != prevote.Value() || correspondingProposal.ValidRound() == -1 {
		return false
	}

	validRound := correspondingProposal.ValidRound()
	votedVatVR := p.Evidence[1].Value()

	// check preVotes at evidence.
	for _, pv := range p.Evidence[1:] {
		if pv.Type() != proto.MsgPrevote || pv.R() != validRound ||
			pv.Value() == correspondingProposal.Value() || pv.Value() != votedVatVR {
			return false
		}
	}

	if haveRedundantVotes(p.Evidence[1:]) {
		return false
	}

	// check if quorum prevote for a different value than V at valid round.
	lastHeader := c.chain.GetHeaderByNumber(p.Message.H() - 1)
	if lastHeader == nil {
		return false
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())

	return engineCore.OverQuorumVotes(p.Evidence[1:], quorum.Uint64()) != nil
}

// check if the AccountabilityProof of challenge of PVO12 is valid.
func (c *MisbehaviourVerifier) validMisbehaviourOfPVO12(p *AccountabilityProof) bool {
	if len(p.Evidence) < 2 {
		return false
	}
	prevote := p.Message
	if prevote.Type() != proto.MsgPrevote || prevote.Value() == nilValue {
		return false
	}

	// check if the corresponding proposal of preVote.
	correspondingProposal := p.Evidence[0]
	if correspondingProposal.Type() != proto.MsgLiteProposal || correspondingProposal.H() != prevote.H() ||
		correspondingProposal.R() != prevote.R() || correspondingProposal.Value() != prevote.Value() ||
		correspondingProposal.ValidRound() == -1 {
		return false
	}

	currentRound := correspondingProposal.R()
	validRound := correspondingProposal.ValidRound()
	allPreCommits := p.Evidence[1:]
	// check if there are any msg out of range (validRound, currentRound), and with correct address, height and code.
	// check if all precommits between range (validRound, currentRound) are presented. There should be only one pc per round.
	presentedRounds := make(map[int64]struct{})
	for _, pc := range allPreCommits {
		if pc.R() <= validRound || pc.R() >= currentRound || pc.Type() != proto.MsgPrecommit || pc.Sender() != prevote.Sender() ||
			pc.H() != prevote.H() {
			return false
		}
		if _, ok := presentedRounds[pc.R()]; ok {
			return false
		}
		presentedRounds[pc.R()] = struct{}{}
	}

	if len(presentedRounds) != int(currentRound-validRound)-1 {
		return false
	}

	// If the last precommit for notV is after the last one for V, raise misbehaviour
	// If all precommits are nil, do not raise misbehaviour. It is a valid correct scenario.
	lastRoundForV := int64(-1)
	lastRoundForNotV := int64(-1)
	for _, pc := range allPreCommits {
		if pc.Value() == prevote.Value() && pc.R() > lastRoundForV {
			lastRoundForV = pc.R()
		}

		if pc.Value() != prevote.Value() && pc.Value() != nilValue && pc.R() > lastRoundForNotV {
			lastRoundForNotV = pc.R()
		}
	}

	return lastRoundForNotV > lastRoundForV
}

// check if the proof of challenge of PVO3 is valid.
func (c *MisbehaviourVerifier) validMisbehaviourOfPVO3(p *AccountabilityProof) bool {
	if len(p.Evidence) != 1 {
		return false
	}
	prevote := p.Message
	if prevote.Type() != proto.MsgPrevote || prevote.Value() == nilValue {
		return false
	}

	// check if the corresponding proposal of preVote is presented, and it contains an invalid validRound.
	correspondingProposal := p.Evidence[0]
	if correspondingProposal.Type() != proto.MsgLiteProposal || correspondingProposal.H() != prevote.H() ||
		correspondingProposal.R() != prevote.R() || correspondingProposal.Value() != prevote.Value() ||
		correspondingProposal.ValidRound() < correspondingProposal.R() {
		return false
	}
	return true
}

// check if the AccountabilityProof of challenge of C is valid.
func (c *MisbehaviourVerifier) validMisbehaviourOfC(p *AccountabilityProof) bool {
	if len(p.Evidence) == 0 {
		return false
	}
	preCommit := p.Message
	if preCommit.Type() != proto.MsgPrecommit || preCommit.Value() == nilValue {
		return false
	}

	// check preVotes for not the same V compares to preCommit.
	for _, m := range p.Evidence {
		if m.Type() != proto.MsgPrevote || m.Value() == preCommit.Value() || m.R() != preCommit.R() {
			return false
		}
	}

	// check no redundant vote msg in evidence in case of hacking.
	if haveRedundantVotes(p.Evidence) {
		return false
	}

	// check if preVotes for not V reaches to quorum.
	lastHeader := c.chain.GetHeaderByNumber(p.Message.H() - 1)
	if lastHeader == nil {
		return false
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())

	return engineCore.OverQuorumVotes(p.Evidence, quorum.Uint64()) != nil
}

// InnocenceVerifier implemented as a native contract to validate an innocence AccountabilityProof.
type InnocenceVerifier struct {
	chain BlockChainContext
}

// RequiredGas the gas cost to execute this AccountabilityProof validator contract, weighted by input data size.
func (c *InnocenceVerifier) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

// Run InnocenceVerifier, take the rlp encoded AccountabilityProof of innocence, decode it and validate it, if the AccountabilityProof is valid, then
// return the rlp hash of msg and the rlp hash of msg sender as the valid identity for on-chain management of proofs,
// AC need the check the value returned to match the ID which is on challenge, to remove the challenge from chain.
func (c *InnocenceVerifier) Run(input []byte, blockNumber uint64) ([]byte, error) {
	// take an on-chain innocent AccountabilityProof, tell the results of the checking
	if len(input) <= 32 || blockNumber == 0 {
		return failure128Byte, nil
	}

	// the 1st 32 bytes are length of bytes array in solidity, take RLP bytes after it.
	p, err := decodeRawProof(input[32:])
	if err != nil {
		return failure128Byte, nil
	}
	return c.validateInnocenceProof(p), nil
}

// validate if the innocence AccountabilityProof is valid, it returns sender address and msg hash in byte array when AccountabilityProof is valid.
func (c *InnocenceVerifier) validateInnocenceProof(in *AccountabilityProof) []byte {
	// check if evidence msgs are from committee members of that height.
	h, err := in.Message.Height()
	if err != nil {
		return failure128Byte
	}

	lastHeader := c.chain.GetHeaderByNumber(h.Uint64() - 1)
	if lastHeader == nil {
		return failure128Byte
	}

	// validate message.
	if _, err = in.Message.Validate(crypto.CheckValidatorSignature, lastHeader); err != nil {
		return failure128Byte
	}

	// to prevent the abuse of the proof message.
	if len(in.Evidence) > maxEvidenceMessages(lastHeader) {
		return failure128Byte
	}

	for _, m := range in.Evidence {
		// the height of msg of the evidences is checked at Validate function.
		if _, err = m.Validate(crypto.CheckValidatorSignature, lastHeader); err != nil {
			return failure128Byte
		}
	}

	if !validInnocenceProof(in, c.chain) {
		return failure128Byte
	}

	return resultBytes(in.Message.MsgHash(), in.Message.Address, in.Rule)
}

func validInnocenceProof(p *AccountabilityProof, chain BlockChainContext) bool {
	// rule engine only have 4 kind of provable accusation for the time being.
	switch p.Rule {
	case autonity.PO:
		return validInnocenceProofOfPO(p, chain)
	case autonity.PVN:
		return validInnocenceProofOfPVN(p)
	case autonity.PVO:
		return validInnocenceProofOfPVO(p, chain)
	case autonity.C1:
		return validInnocenceProofOfC1(p, chain)
	default:
		return false
	}
}

// check if the AccountabilityProof of innocent of PO is valid.
func validInnocenceProofOfPO(p *AccountabilityProof, chain BlockChainContext) bool {
	// check if there is quorum number of prevote at the same value on the same valid round
	proposal := p.Message
	if proposal.Type() != proto.MsgLiteProposal || proposal.ValidRound() == -1 {
		return false
	}

	for _, m := range p.Evidence {
		if !(m.Type() == proto.MsgPrevote && m.Value() == proposal.Value() &&
			m.R() == proposal.ValidRound()) {
			return false
		}
	}

	// check no redundant vote msg in evidence in case of hacking.
	if haveRedundantVotes(p.Evidence) {
		return false
	}

	// check quorum prevotes for V at validRound.
	lastHeader := chain.GetHeaderByNumber(proposal.H() - 1)
	if lastHeader == nil {
		return false
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())
	return engineCore.OverQuorumVotes(p.Evidence, quorum.Uint64()) != nil
}

// check if the AccountabilityProof of innocent of PVN is valid.
func validInnocenceProofOfPVN(p *AccountabilityProof) bool {
	preVote := p.Message
	if !(preVote.Type() == proto.MsgPrevote && preVote.Value() != nilValue) {
		return false
	}

	if len(p.Evidence) != 1 {
		return false
	}

	proposal := p.Evidence[0]
	return proposal.Type() == proto.MsgLiteProposal && proposal.H() == preVote.H() && proposal.R() == preVote.R() &&
		proposal.ValidRound() == -1 && proposal.Value() == preVote.Value()
}

// check if the AccountabilityProof of innocent of PVO is valid.
func validInnocenceProofOfPVO(p *AccountabilityProof, chain BlockChainContext) bool {
	// check if there is quorum number of prevote at the same value on the same valid round
	preVote := p.Message
	if !(preVote.Type() == proto.MsgPrevote && preVote.Value() != nilValue) {
		return false
	}

	if len(p.Evidence) <= 1 {
		return false
	}

	proposal := p.Evidence[0]
	if proposal.Type() != proto.MsgLiteProposal || proposal.Value() != preVote.Value() ||
		proposal.R() != preVote.R() || proposal.ValidRound() == -1 {
		return false
	}

	vr := proposal.ValidRound()
	// check prevotes for V at the valid round.
	for _, m := range p.Evidence[1:] {
		if !(m.Type() == proto.MsgPrevote && m.Value() == proposal.Value() &&
			m.R() == vr) {
			return false
		}
	}

	// check no redundant vote msg in evidence in case of hacking.
	if haveRedundantVotes(p.Evidence[1:]) {
		return false
	}

	// check quorum prevotes at valid round.
	height := preVote.H()
	lastHeader := chain.GetHeaderByNumber(height - 1)
	if lastHeader == nil {
		return false
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())
	return engineCore.OverQuorumVotes(p.Evidence[1:], quorum.Uint64()) != nil
}

// check if the AccountabilityProof of innocent of C1 is valid.
func validInnocenceProofOfC1(p *AccountabilityProof, chain BlockChainContext) bool {
	preCommit := p.Message
	if !(preCommit.Type() == proto.MsgPrecommit && preCommit.Value() != nilValue) {
		return false
	}

	// check quorum prevotes for V at the same round.
	for _, m := range p.Evidence {
		if !(m.Type() == proto.MsgPrevote && m.Value() == preCommit.Value() &&
			m.R() == preCommit.R()) {
			return false
		}
	}

	// check no redundant vote msg in evidence in case of hacking.
	if haveRedundantVotes(p.Evidence) {
		return false
	}

	height := preCommit.H()
	lastHeader := chain.GetHeaderByNumber(height - 1)
	if lastHeader == nil {
		return false
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())
	return engineCore.OverQuorumVotes(p.Evidence, quorum.Uint64()) != nil
}

func haveRedundantVotes(votes []*mUtils.Message) bool {
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

// decode AccountabilityProof convert AccountabilityProof from rlp encoded bytes into object AccountabilityProof.
func decodeRawProof(b []byte) (*AccountabilityProof, error) {
	p := new(AccountabilityProof)
	err := rlp.DecodeBytes(b, p)
	if err != nil {
		return p, err
	}

	// decode the msg.
	if p.Message == nil {
		return p, fmt.Errorf("nil message")
	}

	err = decodeConsensusMsg(p.Message)
	if err != nil {
		return p, err
	}

	// decode the evidence.
	for _, m := range p.Evidence {
		if m == nil {
			return p, fmt.Errorf("nil message")
		}

		err = decodeConsensusMsg(m)
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

// called by precompiled contract to verify an auto-incriminating proposal.
func validProofOfBadProposal(chain BlockChainContext, proof *AccountabilityProof) bool {
	proposal := proof.Message
	if proposal.Type() != proto.MsgLiteProposal {
		return false
	}
	for _, m := range proof.Evidence {
		if m.Type() != proto.MsgPrevote || m.Value() != nilValue || m.BadProposer() != proposal.Sender() ||
			m.BadValue() != proposal.Value() || m.R() != proposal.R() {
			return false
		}
	}

	height := proposal.H()
	lastHeader := chain.GetHeaderByNumber(height - 1)
	if lastHeader == nil {
		return false
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())
	return engineCore.OverQuorumVotes(proof.Evidence, quorum.Uint64()) != nil
}

// checkMsgSignature, it checks if msg is from valid member of the committee.
func checkMsgSignature(chain BlockChainContext, m *mUtils.Message) error {
	lastHeader := chain.GetHeaderByNumber(m.H() - 1)
	if lastHeader == nil {
		return errFutureMsg
	}

	if _, err := m.Validate(crypto.CheckValidatorSignature, lastHeader); err != nil {
		return errNotCommitteeMsg
	}
	return nil
}

func checkEquivocation(m *mUtils.Message, proof []*mUtils.Message) error {
	if len(proof) == 0 {
		return fmt.Errorf("no proof")
	}
	// check equivocations.
	if !sameConsensusMsg(m, proof[0]) {
		return errEquivocation
	}
	return nil
}

func resultBytes(msgHash common.Hash, sender common.Address, rule autonity.Rule) []byte {
	msgHashBytes := msgHash.Bytes()
	senderBytes := common.LeftPadBytes(sender.Bytes(), 32)
	ruleID := make([]byte, 1)
	ruleID[0] = byte(rule)
	ruleBytes := common.LeftPadBytes(ruleID, 32)
	return append(append(append(senderBytes, msgHashBytes...), validByte...), ruleBytes...)
}

func maxEvidenceMessages(header *types.Header) int {
	committeeSize := len(header.Committee)
	if committeeSize > constants.MaxRound {
		return committeeSize + 1
	}
	return constants.MaxRound + 1
}
