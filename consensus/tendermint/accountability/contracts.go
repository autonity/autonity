package accountability

import (
	"errors"
	"fmt"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	engineCore "github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
	"github.com/autonity/autonity/rlp"
	"math"
	"math/big"
)

// precompiled contracts to be call by autonity contract to verify on-chain proofs of accountability events, they are
// a part of consensus.

var (
	checkAccusationAddress   = common.BytesToAddress([]byte{0xfc})
	checkInnocenceAddress    = common.BytesToAddress([]byte{0xfd})
	checkMisbehaviourAddress = common.BytesToAddress([]byte{0xfe})
	// error codes of the execution of precompiled contract to verify the input Proof.
	successResult          = common.LeftPadBytes([]byte{1}, 32)
	failureReturn          = make([]byte, 128)
	errBadHeight           = errors.New("height invalid")
	errMaxEvidences        = errors.New("above max evidence threshold")
	errTooRecentAccusation = errors.New("accusation is too recent")
	errTooOldAccusation    = errors.New("accusation is too old")
)

const KB = 1024

// LoadPrecompiles init the instances of Fault Detector contracts, and register them into EVM's context
func LoadPrecompiles() {
	vm.PrecompiledContractRWMutex.Lock()
	defer vm.PrecompiledContractRWMutex.Unlock()
	pv := InnocenceVerifier{}
	cv := MisbehaviourVerifier{}
	av := AccusationVerifier{}
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
}

// RequiredGas the gas cost to execute AccusationVerifier contract, weighted by input data size.
func (a *AccusationVerifier) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

// executes checks that can be done before even verifying signatures
func preVerifyAccusation(m message.Msg, currentHeight uint64) error {
	accusationHeight := m.H()

	// has to be at least DeltaBlocks old
	if currentHeight <= (DeltaBlocks+1) || accusationHeight >= currentHeight-DeltaBlocks {
		return errTooRecentAccusation
	}
	// cannot be too old. Otherwise this could be exploited by a malicious peer to raise an undefendable accusation.
	// additionally we allocate accountabilityHeightRange/4 more blocks as buffer time to avoid race conditions
	// between msgStore garbage collection and innocence proof generation
	if (currentHeight - accusationHeight) > (HeightRange - (HeightRange / 4)) {
		return errTooOldAccusation
	}

	return nil
}

// Run take the rlp encoded Proof of accusation in byte array, decode it and validate it, if the Proof is valid, then
// the rlp hash of the msg payload and the msg sender is returned.
func (a *AccusationVerifier) Run(input []byte, blockNumber uint64, e *vm.EVM, _ common.Address) ([]byte, error) {
	if len(input) <= 32 {
		return failureReturn, nil
	}
	// the 1st 32 bytes are length of bytes array in solidity, take RLP bytes after it.
	p, err := decodeRawProof(input[32:])
	if err != nil {
		return failureReturn, nil
	}

	// Do preliminary checks that do not rely on signature correctness
	// NOTE: We do not have guarantees that: a.chain.CurrentBlock().NumberU64() == blockNumber - 1
	// This is because the chain head can change while we are executing this tx, therefore the blockNumber might become obsolete.
	if err = preVerifyAccusation(p.Message, blockNumber); err != nil {
		return failureReturn, nil
	}

	// if the suspicious message is for a value that got committed in the same height --> reject accusation
	hash := e.Context.GetHash(p.Message.H())
	if hash == p.Message.Value() {
		return failureReturn, nil
	}

	committee, err := committeeOfHeight(p.Message.H(), e)
	if err != nil {
		return failureReturn, nil
	}

	if err := verifyProofSignatures(committee, p); err != nil {
		return failureReturn, nil
	}

	if verifyAccusation(p) {
		return validReturn(p.Message, p.Rule), nil
	}
	return failureReturn, nil
}

// validate the submitted accusation by the contract call.
func verifyAccusation(p *Proof) bool {
	// we have only 4 types of rule on accusation.
	// an improvement of this function would be to return an error instead of a bool
	switch p.Rule {
	case autonity.PO:
		if p.Message.Code() != message.LightProposalCode || p.Message.(*message.LightProposal).ValidRound() == -1 {
			return false
		}
	case autonity.PVN:
		if p.Message.Code() != message.PrevoteCode || p.Message.Value() == nilValue {
			return false
		}
	case autonity.PVO:
		// theoretically we do not need the non-nil check, since we will check later that prevote.value == proposal.value
		// however added for simplicity of understanding
		if p.Message.Code() != message.PrevoteCode || p.Message.Value() == nilValue {
			return false
		}
	case autonity.C1:
		if p.Message.Code() != message.PrecommitCode || p.Message.Value() == nilValue {
			return false
		}
	default:
		return false
	}

	// p case of PVO accusation, we need to check corresponding old proposal of this preVote.
	if p.Rule == autonity.PVO {
		if len(p.Evidences) != 1 {
			return false
		}
		oldProposal := p.Evidences[0]
		if oldProposal.Code() != message.LightProposalCode ||
			oldProposal.R() != p.Message.R() ||
			oldProposal.Value() != p.Message.Value() ||
			oldProposal.(*message.LightProposal).ValidRound() == -1 {
			return false
		}
	} else if len(p.Evidences) > 0 {
		// do not allow useless evidences
		return false
	}

	return true
}

// MisbehaviourVerifier implemented as a native contract to validate if misbehaviour is valid
type MisbehaviourVerifier struct {
}

// RequiredGas the gas cost to execute MisbehaviourVerifier contract, weighted by input data size.
func (c *MisbehaviourVerifier) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

// Run take the rlp encoded Proof of challenge in byte array, decode it and validate it, if the Proof is valid, then
// the rlp hash of the msg payload and the msg sender is returned as the valid identity for Proof management.
func (c *MisbehaviourVerifier) Run(input []byte, _ uint64, e *vm.EVM, _ common.Address) ([]byte, error) {
	if len(input) <= 32 {
		return failureReturn, nil
	}
	// the 1st 32 bytes are length of bytes array in solidity, take RLP bytes after it.
	p, err := decodeRawProof(input[32:])
	if err != nil {
		return failureReturn, nil
	}

	committee, err := committeeOfHeight(p.Message.H(), e)
	if err != nil {
		return failureReturn, nil
	}

	if err = verifyProofSignatures(committee, p); err != nil {
		return failureReturn, nil
	}
	return c.validateFault(committee, p), nil
}

// validate a misbehavior proof, doesn't check the proof signatures.
func (c *MisbehaviourVerifier) validateFault(committee *types.Committee, p *Proof) []byte {
	valid := false
	switch p.Rule {
	case autonity.PN:
		valid = c.validMisbehaviourOfPN(p)
	case autonity.PO:
		valid = c.validMisbehaviourOfPO(committee, p)
	case autonity.PVN:
		valid = c.validMisbehaviourOfPVN(p)
	case autonity.PVO:
		valid = c.validMisbehaviourOfPVO(committee, p)
	case autonity.PVO12:
		valid = c.validMisbehaviourOfPVO12(p)
	case autonity.C:
		valid = c.validMisbehaviourOfC(committee, p)
	case autonity.InvalidProposer:
		if lightProposal, ok := p.Message.(*message.LightProposal); ok {
			// todo: replace this proposer election function once we place election function in contract side.
			proposer := committee.Proposer(lightProposal.H()-1, lightProposal.R())
			valid = proposer != lightProposal.Sender()
		}
	case autonity.Equivocation:
		valid = errors.Is(checkEquivocation(p.Message, p.Evidences), errEquivocation)
	default:
		valid = false
	}

	if valid {
		return validReturn(p.Message, p.Rule)
	}
	return failureReturn
}

// check if the Proof of challenge of PN is valid,
// node propose a new value when there is a Proof that it precommit at a different value at previous round.
func (c *MisbehaviourVerifier) validMisbehaviourOfPN(p *Proof) bool {
	if len(p.Evidences) != 1 {
		return false
	}
	// should be a new proposal
	proposal, ok := p.Message.(*message.LightProposal)
	if !ok {
		return false
	}
	if proposal.ValidRound() != -1 {
		return false
	}
	preCommit, ok := p.Evidences[0].(*message.Precommit)
	if !ok {
		return false
	}
	if preCommit.Sender() == proposal.Sender() &&
		preCommit.R() < proposal.R() &&
		preCommit.Value() != nilValue {
		return true
	}
	return false
}

// check if the Proof of challenge of PO is valid
func (c *MisbehaviourVerifier) validMisbehaviourOfPO(committee *types.Committee, p *Proof) bool {
	// should be an old proposal
	proposal, ok := p.Message.(*message.LightProposal)
	if !ok {
		return false
	}
	if proposal.ValidRound() == -1 {
		return false
	}
	// if the proposal contains an invalid valid round, then it is a valid proof.
	if proposal.ValidRound() >= proposal.R() {
		return true
	}

	if len(p.Evidences) == 0 {
		return false
	}

	switch vote := p.Evidences[0].(type) {
	case *message.Precommit:
		if vote.R() == proposal.ValidRound() &&
			vote.Sender() == p.Message.Sender() &&
			vote.Value() != nilValue &&
			vote.Value() != proposal.Value() {
			return true
		}
		if vote.R() > proposal.ValidRound() &&
			vote.R() < proposal.R() &&
			vote.Sender() == p.Message.Sender() &&
			vote.Value() != nilValue {
			return true
		}
	case *message.Prevote:
		// check if there are quorum prevotes for other value than the proposed value at valid round.
		for _, m := range p.Evidences {
			pv, ok := m.(*message.Prevote)
			if !ok {
				return false
			}
			if pv.R() != proposal.ValidRound() || pv.Value() == proposal.Value() {
				return false
			}
		}
		if hasEquivocatedVotes(p.Evidences) {
			return false
		}

		return overQuorumVotes(committee, p.Evidences)
	}
	return false
}

// check if the Proof of challenge of PVN is valid.
func (c *MisbehaviourVerifier) validMisbehaviourOfPVN(p *Proof) bool {
	if len(p.Evidences) == 0 {
		return false
	}
	prevote := p.Message
	if prevote.Code() != message.PrevoteCode || prevote.Value() == nilValue {
		return false
	}

	// check if the corresponding new proposal of preVote for new value is presented.
	correspondingProposal := p.Evidences[0]
	if correspondingProposal.Code() != message.LightProposalCode ||
		correspondingProposal.H() != prevote.H() ||
		correspondingProposal.R() != prevote.R() ||
		correspondingProposal.Value() != prevote.Value() ||
		correspondingProposal.(*message.LightProposal).ValidRound() != -1 {
		return false
	}

	// validate preCommits from round R' to R, make sure there is no gap, and the value preCommitted at R'
	// is different from the value preVoted, and the other ones are preCommits of nil.
	preCommits := p.Evidences[1:]

	lastIndex := len(preCommits) - 1

	for i, pc := range preCommits {
		if pc.Code() != message.PrecommitCode || pc.Sender() != prevote.Sender() || pc.R() >= prevote.R() {
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
func (c *MisbehaviourVerifier) validMisbehaviourOfPVO(committee *types.Committee, p *Proof) bool {
	if len(p.Evidences) < 2 {
		return false
	}
	prevote := p.Message
	if prevote.Code() != message.PrevoteCode || prevote.Value() == nilValue {
		return false
	}
	// check if the corresponding proposal of preVote is presented.
	correspondingProposal := p.Evidences[0]
	if correspondingProposal.Code() != message.LightProposalCode ||
		correspondingProposal.H() != prevote.H() ||
		correspondingProposal.R() != prevote.R() ||
		correspondingProposal.Value() != prevote.Value() ||
		correspondingProposal.(*message.LightProposal).ValidRound() == -1 {
		return false
	}

	validRound := correspondingProposal.(*message.LightProposal).ValidRound()
	votedVatVR := p.Evidences[1].Value()

	// check preVotes at evidence.
	for _, pv := range p.Evidences[1:] {
		if pv.Code() != message.PrevoteCode || pv.R() != validRound || pv.Value() == nilValue ||
			pv.Value() == correspondingProposal.Value() || pv.Value() != votedVatVR {
			return false
		}
	}

	if hasEquivocatedVotes(p.Evidences[1:]) {
		return false
	}

	// check if quorum prevote for a different value than V at valid round.
	return overQuorumVotes(committee, p.Evidences[1:])
}

// check if the Proof of challenge of PVO12 is valid.
func (c *MisbehaviourVerifier) validMisbehaviourOfPVO12(p *Proof) bool {
	if len(p.Evidences) < 2 {
		return false
	}
	prevote := p.Message
	if prevote.Code() != message.PrevoteCode || prevote.Value() == nilValue {
		return false
	}

	// check if the corresponding proposal of preVote.
	correspondingProposal := p.Evidences[0]
	if correspondingProposal.Code() != message.LightProposalCode ||
		correspondingProposal.H() != prevote.H() ||
		correspondingProposal.R() != prevote.R() ||
		correspondingProposal.Value() != prevote.Value() ||
		correspondingProposal.(*message.LightProposal).ValidRound() == -1 {
		return false
	}

	currentRound := correspondingProposal.R()
	validRound := correspondingProposal.(*message.LightProposal).ValidRound()
	allPreCommits := p.Evidences[1:]
	// check if there are any msg out of range (validRound, currentRound), and with correct address, height and code.
	// check if all precommits between range (validRound, currentRound) are presented. There should be only one pc per round.
	presentedRounds := make(map[int64]struct{})
	for _, pc := range allPreCommits {
		if pc.R() <= validRound || pc.R() >= currentRound || pc.Code() != message.PrecommitCode || pc.Sender() != prevote.Sender() ||
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

// check if the Proof of challenge of C is valid.
func (c *MisbehaviourVerifier) validMisbehaviourOfC(committee *types.Committee, p *Proof) bool {
	if len(p.Evidences) == 0 {
		return false
	}
	preCommit := p.Message
	if preCommit.Code() != message.PrecommitCode || preCommit.Value() == nilValue {
		return false
	}

	// check preVotes for not the same V compares to preCommit.
	for _, m := range p.Evidences {
		if m.Code() != message.PrevoteCode || m.Value() == preCommit.Value() || m.R() != preCommit.R() {
			return false
		}
	}

	// check no redundant vote msg in evidence in case of hacking.
	if hasEquivocatedVotes(p.Evidences) {
		return false
	}

	// check if preVotes for not V reaches to quorum.
	return overQuorumVotes(committee, p.Evidences)
}

// InnocenceVerifier implemented as a native contract to validate an innocence Proof.
type InnocenceVerifier struct {
}

// RequiredGas the gas cost to execute this Proof validator contract, weighted by input data size.
func (c *InnocenceVerifier) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

// Run InnocenceVerifier, take the rlp encoded Proof of innocence, decode it and validate it, if the Proof is valid, then
// return the rlp hash of msg and the rlp hash of msg sender as the valid identity for on-chain management of proofs,
// AC need the check the value returned to match the ID which is on challenge, to remove the challenge from chain.
func (c *InnocenceVerifier) Run(input []byte, blockNumber uint64, e *vm.EVM, _ common.Address) ([]byte, error) {
	if len(input) <= 32 || blockNumber == 0 {
		return failureReturn, nil
	}
	// the 1st 32 bytes are length of bytes array in solidity, take RLP bytes after it.
	p, err := decodeRawProof(input[32:])
	if err != nil {
		return failureReturn, nil
	}

	committee, err := committeeOfHeight(p.Message.H(), e)
	if err != nil {
		return failureReturn, nil
	}

	if err = verifyProofSignatures(committee, p); err != nil {
		return failureReturn, nil
	}
	if !verifyInnocenceProof(committee, p) {
		return failureReturn, nil
	}
	return validReturn(p.Message, p.Rule), nil
}

func verifyInnocenceProof(committee *types.Committee, p *Proof) bool {
	// rule engine only have 4 kind of provable accusation for the time being.
	switch p.Rule {
	case autonity.PO:
		return validInnocenceProofOfPO(committee, p)
	case autonity.PVN:
		return validInnocenceProofOfPVN(p)
	case autonity.PVO:
		return validInnocenceProofOfPVO(committee, p)
	case autonity.C1:
		return validInnocenceProofOfC1(committee, p)
	default:
		return false
	}
}

// check if the Proof of innocent of PO is valid.
func validInnocenceProofOfPO(committee *types.Committee, p *Proof) bool {
	// check if there is quorum number of prevote at the same value on the same valid round
	proposal := p.Message
	if proposal.Code() != message.LightProposalCode || proposal.(*message.LightProposal).ValidRound() == -1 {
		return false
	}

	for _, m := range p.Evidences {
		if !(m.Code() == message.PrevoteCode &&
			m.Value() == proposal.Value() &&
			m.R() == proposal.(*message.LightProposal).ValidRound()) {
			return false
		}
	}

	// check no redundant vote msg in evidence in case of hacking.
	if hasEquivocatedVotes(p.Evidences) {
		return false
	}

	// check quorum prevotes for V at validRound.
	return overQuorumVotes(committee, p.Evidences)
}

// check if the Proof of innocent of PVN is valid.
func validInnocenceProofOfPVN(p *Proof) bool {
	preVote := p.Message
	if !(preVote.Code() == message.PrevoteCode && preVote.Value() != nilValue) {
		return false
	}

	if len(p.Evidences) != 1 {
		return false
	}

	proposal := p.Evidences[0]
	return proposal.Code() == message.LightProposalCode &&
		proposal.H() == preVote.H() &&
		proposal.R() == preVote.R() &&
		proposal.(*message.LightProposal).ValidRound() == -1 &&
		proposal.Value() == preVote.Value()
}

// check if the Proof of innocent of PVO is valid.
func validInnocenceProofOfPVO(committee *types.Committee, p *Proof) bool {
	// check if there is quorum number of prevote at the same value on the same valid round
	preVote := p.Message
	if !(preVote.Code() == message.PrevoteCode && preVote.Value() != nilValue) {
		return false
	}

	if len(p.Evidences) <= 1 {
		return false
	}

	proposal := p.Evidences[0]
	if proposal.Code() != message.LightProposalCode ||
		proposal.Value() != preVote.Value() ||
		proposal.R() != preVote.R() ||
		proposal.(*message.LightProposal).ValidRound() == -1 {
		return false
	}

	vr := proposal.(*message.LightProposal).ValidRound()
	// check prevotes for V at the valid round.
	for _, m := range p.Evidences[1:] {
		if !(m.Code() == message.PrevoteCode && m.Value() == proposal.Value() && m.R() == vr) {
			return false
		}
	}

	// check no redundant vote msg in evidence in case of hacking.
	if hasEquivocatedVotes(p.Evidences[1:]) {
		return false
	}

	// check quorum prevotes at valid round.
	return overQuorumVotes(committee, p.Evidences[1:])
}

// check if the Proof of innocent of C1 is valid.
func validInnocenceProofOfC1(committee *types.Committee, p *Proof) bool {
	preCommit, ok := p.Message.(*message.Precommit)
	if !ok {
		return false
	}
	if preCommit.Value() == nilValue {
		return false
	}
	// check quorum prevotes for V at the same round.
	for _, m := range p.Evidences {
		if !(m.Code() == message.PrevoteCode && m.Value() == preCommit.Value() &&
			m.R() == preCommit.R()) {
			return false
		}
	}
	// check no redundant vote msg in evidence in case of hacking.
	if hasEquivocatedVotes(p.Evidences) {
		return false
	}

	return overQuorumVotes(committee, p.Evidences)
}

func hasEquivocatedVotes(votes []message.Msg) bool {
	voteMap := make(map[common.Address]struct{})
	for _, vote := range votes {
		_, ok := voteMap[vote.Sender()]
		if !ok {
			voteMap[vote.Sender()] = struct{}{}
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
	return p, nil
}

// checkMsgSignature checks if the consensus message is from valid member of the committee.
func verifyProofSignatures(committee *types.Committee, p *Proof) error {
	h := p.Message.H()
	if h == 0 {
		return errBadHeight
	}

	if err := p.Message.Validate(committee.CommitteeMember); err != nil {
		return errNotCommitteeMsg
	}
	// check if the number of evidence msgs are exceeded the max to prevent the abuse of the proof msg.
	if len(p.Evidences) > maxEvidenceMessages(committee) {
		return errMaxEvidences
	}
	for _, msg := range p.Evidences {
		if msg.H() != h {
			return errBadHeight
		}
		if err := msg.Validate(committee.CommitteeMember); err != nil {
			return errNotCommitteeMsg
		}
	}
	return nil
}

func checkEquivocation(m message.Msg, proof []message.Msg) error {
	if len(proof) == 0 {
		return fmt.Errorf("no proof")
	}
	// check equivocations.
	if m.Hash() != proof[0].Hash() {
		return errEquivocation
	}
	return nil
}

func validReturn(m message.Msg, rule autonity.Rule) []byte {
	offender := common.LeftPadBytes(m.Sender().Bytes(), 32)
	ruleID := common.LeftPadBytes([]byte{byte(rule)}, 32)
	block := make([]byte, 32)
	block = common.LeftPadBytes(new(big.Int).SetUint64(m.H()).Bytes(), 32)
	result := make([]byte, 160)
	copy(result[0:32], successResult)
	copy(result[32:64], offender)
	copy(result[64:96], ruleID)
	copy(result[96:128], block)
	copy(result[128:160], m.Hash().Bytes())
	return result
}

func maxEvidenceMessages(committee *types.Committee) int {
	// todo(youssef): I dont understand that
	committeeSize := committee.Len()
	if committeeSize > constants.MaxRound {
		return committeeSize + 1
	}
	return constants.MaxRound + 1
}

func overQuorumVotes(committee *types.Committee, msgs []message.Msg) bool {
	q := bft.Quorum(committee.TotalVotingPower())
	return engineCore.OverQuorumVotes(msgs, q) != nil
}

func committeeOfHeight(height uint64, evm *vm.EVM) (*types.Committee, error) {
	var committeeSet []types.CommitteeMember
	err := acCall(evm, "getCommitteeByHeight", &committeeSet, new(big.Int).SetUint64(height))
	if err != nil {
		return nil, err
	}
	committee := &types.Committee{}
	if len(committeeSet) != 0 {
		committee.Members = make([]*types.CommitteeMember, len(committeeSet))
		for i, m := range committeeSet {
			committee.Members[i] = &types.CommitteeMember{
				Address:      m.Address,
				VotingPower:  new(big.Int).Set(m.VotingPower),
				ConsensusKey: m.ConsensusKey,
			}
		}
		// As the committee is already sorted by the contract, thus we don't need sort again.
	}

	return committee, err
}

func acCall(evm *vm.EVM, function string, result any, args ...any) error {
	packedArgs, err := generated.AutonityAbi.Pack(function, args...)
	if err != nil {
		return err
	}
	gas := uint64(math.MaxUint64)
	ret, _, err := evm.Call(vm.AccountRef(params.DeployerAddress), params.AutonityContractAddress, packedArgs, gas, new(big.Int))
	if err != nil {
		return err
	}
	if err = generated.AutonityAbi.UnpackIntoInterface(result, function, ret); err != nil {
		return err
	}
	return nil
}
