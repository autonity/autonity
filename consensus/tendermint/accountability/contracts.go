package accountability

import (
	"errors"
	"fmt"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	engineCore "github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
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
	// error codes of the execution of precompiled contract to verify the input Proof.
	successResult = common.LeftPadBytes([]byte{1}, 32)
	failureResult = make([]byte, 128)

	errNilMessage = errors.New("nil message")
)

const KB = 1024

// LoadPrecompiles init the instances of Fault Detector contracts, and register them into EVM's context
func LoadPrecompiles(chain ChainContext) {
	vm.PrecompiledContractRWMutex.Lock()
	defer vm.PrecompiledContractRWMutex.Unlock()
	pv := InnocenceVerifier{chain: chain}
	cv := MisbehaviourVerifier{chain: chain}
	av := AccusationVerifier{chain: chain}
	setPrecompiles := func(set map[common.Address]vm.PrecompiledContract) {
		set[checkInnocenceAddress] = &pv
		set[checkMisbehaviourAddress] = &cv
		set[checkAccusationAddress] = &av
	}
	setPrecompiles(vm.PrecompiledContractsByzantium)
	setPrecompiles(vm.PrecompiledContractsHomestead)
	setPrecompiles(vm.PrecompiledContractsIstanbul)
	setPrecompiles(vm.PrecompiledContractsBerlin)
	setPrecompiles(vm.PrecompiledContractsBLS)
}

// AccusationVerifier implemented as a native contract to validate if an accusation is valid
type AccusationVerifier struct {
	chain ChainContext
}

// RequiredGas the gas cost to execute AccusationVerifier contract, weighted by input data size.
func (a *AccusationVerifier) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

// Run take the rlp encoded Proof of accusation in byte array, decode it and validate it, if the Proof is valid, then
// the rlp hash of the msg payload and the msg sender is returned.
func (a *AccusationVerifier) Run(input []byte, blockNumber uint64) ([]byte, error) {
	if len(input) <= 32 {
		return failureResult, nil
	}
	// the 1st 32 bytes are length of bytes array in solidity, take RLP bytes after it.
	p, err := decodeRawProof(input[32:])
	if err != nil {
		return failureResult, nil
	}
	evHeight := p.Message.H()
	if evHeight == 0 {
		return failureResult, nil
	}
	// prevent a potential attack: a malicious fault detector can rise an accusation that contain a message
	// corresponding to an old height, while at each Pi, they only buffer specific heights of message in msg store, thus
	// Pi can never provide a valid proof of innocence anymore, making the malicious accusation be valid for slashing.
	if blockNumber > evHeight && (blockNumber-evHeight >= consensus.AccountabilityHeightRange) {
		return failureResult, nil
	}
	if verifyAccusation(a.chain, p) {
		return validReturn(p.Message, p.Rule), nil
	}
	return failureResult, nil
}

// validate the submitted accusation by the contract call.
func verifyAccusation(chain ChainContext, p *Proof) bool {
	// we have only 4 types of rule on accusation.
	// an improvement of this function would be to return an error instead of a bool
	switch p.Rule {
	case autonity.PO:
		if p.Message.Code() != message.LightProposalCode || p.Message.(*message.LightProposal).ValidRound == -1 {
			return false
		}
	case autonity.PVN:
		if p.Message.Code() != message.PrevoteCode {
			return false
		}
	case autonity.PVO:
		if p.Message.Code() != message.PrevoteCode {
			return false
		}
	case autonity.C1:
		if p.Message.Code() != message.PrecommitCode {
			return false
		}
	default:
		return false
	}

	// check if the suspicious msg is from the correct committee of that height.
	h := p.Message.H()
	lastHeader := chain.GetHeaderByNumber(h - 1)
	if lastHeader == nil {
		return false
	}
	if err := p.Message.Validate(lastHeader.CommitteeMember); err != nil {
		return false
	}

	// p case of PVO accusation, we need to check corresponding old proposal of this preVote.
	if p.Rule == autonity.PVO {
		if len(p.Evidences) != 1 {
			return false
		}
		oldProposal := p.Evidences[0]
		// Todo(Youssef): bug possible with Light proposal sig validation // signature may come from accuser not reported
		if err := oldProposal.Validate(lastHeader.CommitteeMember); err != nil {
			return false
		}
		if oldProposal.Code != consensus.MsgLightProposal ||
			oldProposal.R() != p.Message.R() ||
			oldProposal.Value() != p.Message.Value() ||
			oldProposal.ConsensusMsg.(*message.LightProposal).ValidRound == -1 {
			return false
		}
	}

	return true
}

// MisbehaviourVerifier implemented as a native contract to validate if misbehaviour is valid
type MisbehaviourVerifier struct {
	chain ChainContext
}

// RequiredGas the gas cost to execute MisbehaviourVerifier contract, weighted by input data size.
func (c *MisbehaviourVerifier) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

// Run take the rlp encoded Proof of challenge in byte array, decode it and validate it, if the Proof is valid, then
// the rlp hash of the msg payload and the msg sender is returned as the valid identity for Proof management.
func (c *MisbehaviourVerifier) Run(input []byte, _ uint64) ([]byte, error) {
	if len(input) <= 32 {
		return failureResult, nil
	}

	// the 1st 32 bytes are length of bytes array in solidity, take RLP bytes after it.
	p, err := decodeRawProof(input[32:])
	if err != nil {
		if p.Rule == autonity.GarbageMessage && err == errAccountableGarbageMsg && len(p.Evidences) == 0 {
			// since garbage message cannot be decoded for msg height, we return 0 for msg height in such case.
			return validReturn(p.Message, p.Rule), nil
		}
		return failureResult, nil
	}

	evHeight := p.Message.H()
	if evHeight == 0 {
		return failureResult, nil
	}

	return c.validateProof(p), nil
}

// validate the Proof, if the Proof is valid, then the rlp hash of the msg payload and rlp hash of msg sender is
// returned as the valid identity for Proof management.
func (c *MisbehaviourVerifier) validateProof(p *Proof) []byte {

	// check if suspicious message is from correct committee member.
	err := checkMsgSignature(c.chain, p.Message)
	if err != nil {
		return failureResult
	}

	// check if evidence msgs are from committee members of that height.
	h := p.Message.H()
	lastHeader := c.chain.GetHeaderByNumber(h - 1)
	if lastHeader == nil {
		return failureResult
	}

	// check if the number of evidence msgs are exceeded the max to prevent the abuse of the proof msg.
	if len(p.Evidences) > maxEvidenceMessages(lastHeader) {
		return failureResult
	}

	for _, msg := range p.Evidences {
		// the height of msg of the evidences is checked at Validate function.
		if err := msg.Validate(crypto.CheckValidatorSignature, lastHeader); err != nil {
			return failureResult
		}
	}

	if c.validProof(p) {
		return validReturn(p.Message, p.Rule)
	}
	return failureResult
}

// check if the evidence of the misbehaviour is valid or not.
func (c *MisbehaviourVerifier) validProof(p *Proof) bool {
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
		return p.Message.R() > constants.MaxRound
	case autonity.WrongValidRound:
		if lightProposal, ok := p.Message.ConsensusMsg.(*message.LightProposal); ok {
			return lightProposal.ValidRound >= lightProposal.Round
		}
		return false
	case autonity.InvalidProposer:
		return !isProposerValid(c.chain, p.Message)
	case autonity.Equivocation:
		return checkEquivocation(p.Message, p.Evidences) == errEquivocation
	default:
		return false
	}
}

// check if the Proof of challenge of PN is valid,
// node propose a new value when there is a Proof that it precommit at a different value at previous round.
func (c *MisbehaviourVerifier) validMisbehaviourOfPN(p *Proof) bool {
	if len(p.Evidences) != 1 {
		return false
	}

	// should be a new proposal
	proposal := p.Message

	if proposal.Code != consensus.MsgLightProposal || proposal.ConsensusMsg.(*message.LightProposal).ValidRound != -1 {
		return false
	}

	preCommit := p.Evidences[0]
	if preCommit.Sender() == proposal.Sender() &&
		preCommit.Type() == consensus.MsgPrecommit &&
		preCommit.R() < proposal.R() &&
		preCommit.Value() != nilValue {
		return true
	}

	return false
}

// check if the Proof of challenge of PO is valid
func (c *MisbehaviourVerifier) validMisbehaviourOfPO(p *Proof) bool {
	// should be an old proposal
	if p.Message.Type() != consensus.MsgLightProposal {
		return false
	}
	proposal := p.Message.ConsensusMsg.(*message.LightProposal)
	if proposal.ValidRound == -1 {
		return false
	}
	// if the proposal contains an invalid valid round, then it is a valid proof.
	if proposal.ValidRound >= proposal.R() {
		return true
	}

	if len(p.Evidences) == 0 {
		return false
	}
	preCommit := p.Evidences[0]

	if preCommit.Type() == consensus.MsgPrecommit && preCommit.R() == proposal.ValidRound &&
		preCommit.Sender() == p.Message.Sender() && preCommit.Value() != nilValue &&
		preCommit.Value() != proposal.V() {
		return true
	}

	if preCommit.Type() == consensus.MsgPrecommit &&
		preCommit.R() > proposal.ValidRound && preCommit.R() < proposal.R() &&
		preCommit.Sender() == p.Message.Sender() &&
		preCommit.Value() != nilValue {
		return true
	}

	// check if there are quorum prevotes for other value than the proposed value at valid round.
	preVote := p.Evidences[0]
	if preVote.Type() == consensus.MsgPrevote {
		// validate evidences
		for _, pv := range p.Evidences {
			if pv.Type() != consensus.MsgPrevote || pv.R() != proposal.ValidRound || pv.Value() == proposal.V() {
				return false
			}
		}

		if hasEquivocatedVotes(p.Evidences) {
			return false
		}

		// check if preVotes for not V reaches to quorum.
		lastHeader := c.chain.GetHeaderByNumber(p.Message.H() - 1)
		if lastHeader == nil {
			return false
		}
		quorum := bft.Quorum(lastHeader.TotalVotingPower())
		return engineCore.OverQuorumVotes(p.Evidences, quorum) != nil
	}

	return false
}

// check if the Proof of challenge of PVN is valid.
func (c *MisbehaviourVerifier) validMisbehaviourOfPVN(p *Proof) bool {
	if len(p.Evidences) == 0 {
		return false
	}
	prevote := p.Message
	if prevote.Type() != consensus.MsgPrevote || prevote.Value() == nilValue {
		return false
	}

	// check if the corresponding new proposal of preVote for new value is presented.
	correspondingProposal := p.Evidences[0]
	if correspondingProposal.Type() != consensus.MsgLightProposal ||
		correspondingProposal.H() != prevote.H() ||
		correspondingProposal.R() != prevote.R() ||
		correspondingProposal.Value() != prevote.Value() ||
		correspondingProposal.ConsensusMsg.(*message.LightProposal).ValidRound != -1 {
		return false
	}

	// validate preCommits from round R' to R, make sure there is no gap, and the value preCommitted at R'
	// is different from the value preVoted, and the other ones are preCommits of nil.
	preCommits := p.Evidences[1:]

	lastIndex := len(preCommits) - 1

	for i, pc := range preCommits {
		if pc.Type() != consensus.MsgPrecommit || pc.Sender() != prevote.Sender() || pc.R() >= prevote.R() {
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
func (c *MisbehaviourVerifier) validMisbehaviourOfPVO(p *Proof) bool {
	if len(p.Evidences) < 2 {
		return false
	}
	prevote := p.Message
	if prevote.Type() != consensus.MsgPrevote || prevote.Value() == nilValue {
		return false
	}
	// check if the corresponding proposal of preVote is presented.
	correspondingProposal := p.Evidences[0]
	if correspondingProposal.Type() != consensus.MsgLightProposal ||
		correspondingProposal.H() != prevote.H() ||
		correspondingProposal.R() != prevote.R() ||
		correspondingProposal.Value() != prevote.Value() ||
		correspondingProposal.ConsensusMsg.(*message.LightProposal).ValidRound == -1 {
		return false
	}

	validRound := correspondingProposal.ConsensusMsg.(*message.LightProposal).ValidRound
	votedVatVR := p.Evidences[1].Value()

	// check preVotes at evidence.
	for _, pv := range p.Evidences[1:] {
		if pv.Type() != consensus.MsgPrevote || pv.R() != validRound || pv.Value() == nilValue ||
			pv.Value() == correspondingProposal.Value() || pv.Value() != votedVatVR {
			return false
		}
	}

	if hasEquivocatedVotes(p.Evidences[1:]) {
		return false
	}

	// check if quorum prevote for a different value than V at valid round.
	lastHeader := c.chain.GetHeaderByNumber(p.Message.H() - 1)
	if lastHeader == nil {
		return false
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())

	return engineCore.OverQuorumVotes(p.Evidences[1:], quorum) != nil
}

// check if the Proof of challenge of PVO12 is valid.
func (c *MisbehaviourVerifier) validMisbehaviourOfPVO12(p *Proof) bool {
	if len(p.Evidences) < 2 {
		return false
	}
	prevote := p.Message
	if prevote.Type() != consensus.MsgPrevote || prevote.Value() == nilValue {
		return false
	}

	// check if the corresponding proposal of preVote.
	correspondingProposal := p.Evidences[0]
	if correspondingProposal.Type() != consensus.MsgLightProposal ||
		correspondingProposal.H() != prevote.H() ||
		correspondingProposal.R() != prevote.R() ||
		correspondingProposal.Value() != prevote.Value() ||
		correspondingProposal.ConsensusMsg.(*message.LightProposal).ValidRound == -1 {
		return false
	}

	currentRound := correspondingProposal.R()
	validRound := correspondingProposal.ConsensusMsg.(*message.LightProposal).ValidRound
	allPreCommits := p.Evidences[1:]
	// check if there are any msg out of range (validRound, currentRound), and with correct address, height and code.
	// check if all precommits between range (validRound, currentRound) are presented. There should be only one pc per round.
	presentedRounds := make(map[int64]struct{})
	for _, pc := range allPreCommits {
		if pc.R() <= validRound || pc.R() >= currentRound || pc.Type() != consensus.MsgPrecommit || pc.Sender() != prevote.Sender() ||
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
func (c *MisbehaviourVerifier) validMisbehaviourOfPVO3(p *Proof) bool {
	if len(p.Evidences) != 1 {
		return false
	}
	prevote := p.Message
	if prevote.Type() != consensus.MsgPrevote || prevote.Value() == nilValue {
		return false
	}

	// check if the corresponding proposal of preVote is presented, and it contains an invalid validRound.
	correspondingProposal := p.Evidences[0]
	if correspondingProposal.Type() != consensus.MsgLightProposal ||
		correspondingProposal.H() != prevote.H() ||
		correspondingProposal.R() != prevote.R() ||
		correspondingProposal.Value() != prevote.Value() ||
		correspondingProposal.ConsensusMsg.(*message.LightProposal).ValidRound < correspondingProposal.R() {
		return false
	}
	return true
}

// check if the Proof of challenge of C is valid.
func (c *MisbehaviourVerifier) validMisbehaviourOfC(p *Proof) bool {
	if len(p.Evidences) == 0 {
		return false
	}
	preCommit := p.Message
	if preCommit.Type() != consensus.MsgPrecommit || preCommit.Value() == nilValue {
		return false
	}

	// check preVotes for not the same V compares to preCommit.
	for _, m := range p.Evidences {
		if m.Type() != consensus.MsgPrevote || m.Value() == preCommit.Value() || m.R() != preCommit.R() {
			return false
		}
	}

	// check no redundant vote msg in evidence in case of hacking.
	if hasEquivocatedVotes(p.Evidences) {
		return false
	}

	// check if preVotes for not V reaches to quorum.
	lastHeader := c.chain.GetHeaderByNumber(p.Message.H() - 1)
	if lastHeader == nil {
		return false
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())

	return engineCore.OverQuorumVotes(p.Evidences, quorum) != nil
}

// InnocenceVerifier implemented as a native contract to validate an innocence Proof.
type InnocenceVerifier struct {
	chain ChainContext
}

// RequiredGas the gas cost to execute this Proof validator contract, weighted by input data size.
func (c *InnocenceVerifier) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

// Run InnocenceVerifier, take the rlp encoded Proof of innocence, decode it and validate it, if the Proof is valid, then
// return the rlp hash of msg and the rlp hash of msg sender as the valid identity for on-chain management of proofs,
// AC need the check the value returned to match the ID which is on challenge, to remove the challenge from chain.
func (c *InnocenceVerifier) Run(input []byte, blockNumber uint64) ([]byte, error) {
	// take an on-chain innocent Proof, tell the results of the checking
	if len(input) <= 32 || blockNumber == 0 {
		return failureResult, nil
	}

	// the 1st 32 bytes are length of bytes array in solidity, take RLP bytes after it.
	p, err := decodeRawProof(input[32:])
	if err != nil {
		return failureResult, nil
	}
	return c.validateInnocenceProof(p), nil
}

// validate if the innocence Proof is valid, it returns sender address and msg hash in byte array when Proof is valid.
func (c *InnocenceVerifier) validateInnocenceProof(in *Proof) []byte {
	// check if evidence msgs are from committee members of that height.
	h := in.Message.H()

	lastHeader := c.chain.GetHeaderByNumber(h - 1)
	if lastHeader == nil {
		return failureResult
	}

	// validate message.
	if err := in.Message.Validate(crypto.CheckValidatorSignature, lastHeader); err != nil {
		return failureResult
	}

	// to prevent the abuse of the proof message.
	if len(in.Evidences) > maxEvidenceMessages(lastHeader) {
		return failureResult
	}

	for _, m := range in.Evidences {
		// the height of msg of the evidences is checked at Validate function.
		if err := m.Validate(crypto.CheckValidatorSignature, lastHeader); err != nil {
			return failureResult
		}
	}

	if !validInnocenceProof(in, c.chain) {
		return failureResult
	}

	return validReturn(in.Message, in.Rule)
}

func validInnocenceProof(p *Proof, chain ChainContext) bool {
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

// check if the Proof of innocent of PO is valid.
func validInnocenceProofOfPO(p *Proof, chain ChainContext) bool {
	// check if there is quorum number of prevote at the same value on the same valid round
	proposal := p.Message
	if proposal.Type() != consensus.MsgLightProposal || proposal.ConsensusMsg.(*message.LightProposal).ValidRound == -1 {
		return false
	}

	for _, m := range p.Evidences {
		if !(m.Type() == consensus.MsgPrevote &&
			m.Value() == proposal.Value() &&
			m.R() == proposal.ConsensusMsg.(*message.LightProposal).ValidRound) {
			return false
		}
	}

	// check no redundant vote msg in evidence in case of hacking.
	if hasEquivocatedVotes(p.Evidences) {
		return false
	}

	// check quorum prevotes for V at validRound.
	lastHeader := chain.GetHeaderByNumber(proposal.H() - 1)
	if lastHeader == nil {
		return false
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())
	return engineCore.OverQuorumVotes(p.Evidences, quorum) != nil
}

// check if the Proof of innocent of PVN is valid.
func validInnocenceProofOfPVN(p *Proof) bool {
	preVote := p.Message
	if !(preVote.Type() == consensus.MsgPrevote && preVote.Value() != nilValue) {
		return false
	}

	if len(p.Evidences) != 1 {
		return false
	}

	proposal := p.Evidences[0]
	return proposal.Type() == consensus.MsgLightProposal &&
		proposal.H() == preVote.H() &&
		proposal.R() == preVote.R() &&
		proposal.ConsensusMsg.(*message.LightProposal).ValidRound == -1 &&
		proposal.Value() == preVote.Value()
}

// check if the Proof of innocent of PVO is valid.
func validInnocenceProofOfPVO(p *Proof, chain ChainContext) bool {
	// check if there is quorum number of prevote at the same value on the same valid round
	preVote := p.Message
	if !(preVote.Type() == consensus.MsgPrevote && preVote.Value() != nilValue) {
		return false
	}

	if len(p.Evidences) <= 1 {
		return false
	}

	proposal := p.Evidences[0]
	if proposal.Type() != consensus.MsgLightProposal ||
		proposal.Value() != preVote.Value() ||
		proposal.R() != preVote.R() ||
		proposal.ConsensusMsg.(*message.LightProposal).ValidRound == -1 {
		return false
	}

	vr := proposal.ConsensusMsg.(*message.LightProposal).ValidRound
	// check prevotes for V at the valid round.
	for _, m := range p.Evidences[1:] {
		if !(m.Type() == consensus.MsgPrevote && m.Value() == proposal.Value() &&
			m.R() == vr) {
			return false
		}
	}

	// check no redundant vote msg in evidence in case of hacking.
	if hasEquivocatedVotes(p.Evidences[1:]) {
		return false
	}

	// check quorum prevotes at valid round.
	height := preVote.H()
	lastHeader := chain.GetHeaderByNumber(height - 1)
	if lastHeader == nil {
		return false
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())
	return engineCore.OverQuorumVotes(p.Evidences[1:], quorum) != nil
}

// check if the Proof of innocent of C1 is valid.
func validInnocenceProofOfC1(p *Proof, chain ChainContext) bool {
	preCommit := p.Message
	if !(preCommit.Type() == consensus.MsgPrecommit && preCommit.Value() != nilValue) {
		return false
	}

	// check quorum prevotes for V at the same round.
	for _, m := range p.Evidences {
		if !(m.Type() == consensus.MsgPrevote && m.Value() == preCommit.Value() &&
			m.R() == preCommit.R()) {
			return false
		}
	}

	// check no redundant vote msg in evidence in case of hacking.
	if hasEquivocatedVotes(p.Evidences) {
		return false
	}

	height := preCommit.H()
	lastHeader := chain.GetHeaderByNumber(height - 1)
	if lastHeader == nil {
		return false
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())
	return engineCore.OverQuorumVotes(p.Evidences, quorum) != nil
}

func hasEquivocatedVotes(votes []*message.Message) bool {
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

// decodeRawProof decodes an RLP-encoded Proof object.
func decodeRawProof(b []byte) (*Proof, error) {
	p := new(Proof)
	if err := rlp.DecodeBytes(b, p); err != nil {
		return p, err
	}
	if p.Message == nil {
		return p, errNilMessage
	}
	if err := decodeMessage(p.Message); err != nil {
		return p, err
	}
	for _, m := range p.Evidences {
		if m == nil {
			return p, errNilMessage
		}
		if err := decodeMessage(m); err != nil {
			return p, err
		}
	}
	return p, nil
}

// checkMsgSignature checks if the consensus message is from valid member of the committee.
func checkMsgSignature(chain ChainContext, m *message.Message) error {
	lastHeader := chain.GetHeaderByNumber(m.H() - 1)
	if lastHeader == nil {
		return errFutureMsg
	}

	if err := m.Validate(crypto.CheckValidatorSignature, lastHeader); err != nil {
		return errNotCommitteeMsg
	}
	return nil
}

func checkEquivocation(m *message.Message, proof []*message.Message) error {
	if len(proof) == 0 {
		return fmt.Errorf("no proof")
	}
	// check equivocations.
	if m.Hash() != proof[0].Hash() {
		return errEquivocation
	}
	return nil
}

func validReturn(m message.Message, rule autonity.Rule) []byte {
	offender := common.LeftPadBytes(m.Sender().Bytes(), 32)
	ruleID := common.LeftPadBytes([]byte{byte(rule)}, 32)
	block := make([]byte, 32)
	if m.ConsensusMsg != nil {
		block = common.LeftPadBytes(m.ConsensusMsg.H().Bytes(), 32)
	}
	result := make([]byte, 160)
	copy(result[0:32], successResult)
	copy(result[32:64], offender)
	copy(result[64:96], ruleID)
	copy(result[96:128], block)
	copy(result[128:160], m.Hash().Bytes())
	return result
}

func maxEvidenceMessages(header *types.Header) int {
	// todo(youssef): I dont understand that
	committeeSize := len(header.Committee)
	if committeeSize > constants.MaxRound {
		return committeeSize + 1
	}
	return constants.MaxRound + 1
}
