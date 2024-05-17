package accountability

import (
	"errors"
	"fmt"
	"github.com/autonity/autonity/core/types"
	"math/big"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	engineCore "github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/rlp"
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
	errTooRecentAccusation = errors.New("accusation is too recent")
	errTooOldAccusation    = errors.New("accusation is too old")
	errValueCommitted      = errors.New("accusation is for a committed value")
)

const KB = 1024

/* TODO: This function subtly breaks the accusation, misbehavior and innocence e2e test.
* This is because since the Precompiled maps are global variables, all nodes in the e2e test
* end up using the same precompiled contracts, which contain the same chain reference.
* I.E. all nodes will use the chain of the last started node when executing precompiled contracts.
* We decided not to fix this issue since it does not affect a standalone client in a production test.
* The real fix here is to remove the chain dependency from precompiled contracts.
 */
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

// executes checks that can be done before even verifying signatures
func preVerifyAccusation(chain ChainContext, m message.Msg, currentHeight uint64) error {
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

	// if the suspicious message is for a value that got committed in the same height --> reject accusation
	if chain.GetBlock(m.Value(), m.H()) != nil {
		return errValueCommitted
	}
	return nil
}

// Run take the rlp encoded Proof of accusation in byte array, decode it and validate it, if the Proof is valid, then
// the rlp hash of the msg payload and the msg signer is returned.
func (a *AccusationVerifier) Run(input []byte, blockNumber uint64, _ *vm.EVM, _ common.Address) ([]byte, error) {
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
	if err := preVerifyAccusation(a.chain, p.Message, blockNumber); err != nil {
		return failureReturn, nil
	}

	committee, err := verifyProofSignatures(a.chain, p)
	if err != nil {
		return failureReturn, nil
	}

	if verifyAccusation(p, committee) {
		// the proof carry valid info.
		return validReturn(p.Message, p.Offender, p.Rule), nil
	}
	return failureReturn, nil
}

// validate the submitted accusation by the contract call.
func verifyAccusation(p *Proof, committee types.Committee) bool {
	// we have only 4 types of rule on accusation.
	// an improvement of this function would be to return an error instead of a bool
	switch p.Rule {
	case autonity.PO:
		if p.Message.Code() != message.LightProposalCode {
			return false
		}
		lightProposal := p.Message.(*message.LightProposal)
		if lightProposal.ValidRound() == -1 || p.Offender != lightProposal.Signer() ||
			p.OffenderIndex != lightProposal.SignerIndex() {
			return false
		}

	case autonity.PVN:
		if p.Message.Code() != message.PrevoteCode || p.Message.Value() == nilValue {
			return false
		}
		prevote := p.Message.(*message.Prevote)
		present := prevote.Signers().Contains(p.OffenderIndex)
		if !present || committee[p.OffenderIndex].Address != p.Offender {
			return false
		}

	case autonity.PVO:
		// theoretically we do not need the non-nil check, since we will check later that prevote.value == proposal.value
		// however added for simplicity of understanding
		if p.Message.Code() != message.PrevoteCode || p.Message.Value() == nilValue {
			return false
		}
		prevote := p.Message.(*message.Prevote)
		present := prevote.Signers().Contains(p.OffenderIndex)
		if !present || committee[p.OffenderIndex].Address != p.Offender {
			return false
		}

	case autonity.C1:
		if p.Message.Code() != message.PrecommitCode || p.Message.Value() == nilValue {
			return false
		}

		precommit := p.Message.(*message.Precommit)
		present := precommit.Signers().Contains(p.OffenderIndex)
		if !present || committee[p.OffenderIndex].Address != p.Offender {
			return false
		}

	default:
		return false
	}

	// in case of PVO accusation, we need to check corresponding old proposal of this preVote.
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
	chain ChainContext
}

// RequiredGas the gas cost to execute MisbehaviourVerifier contract, weighted by input data size.
func (c *MisbehaviourVerifier) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

// Run take the rlp encoded Proof of challenge in byte array, decode it and validate it, if the Proof is valid, then
// the rlp hash of the msg payload and the msg signer is returned as the valid identity for Proof management.
func (c *MisbehaviourVerifier) Run(input []byte, _ uint64, _ *vm.EVM, _ common.Address) ([]byte, error) {
	if len(input) <= 32 {
		return failureReturn, nil
	}
	// the 1st 32 bytes are length of bytes array in solidity, take RLP bytes after it.
	p, err := decodeRawProof(input[32:])
	if err != nil {
		return failureReturn, nil
	}

	committee, err := verifyProofSignatures(c.chain, p)
	if err != nil {
		return failureReturn, nil
	}
	return c.validateFault(p, committee), nil
}

// validate a misbehavior proof, doesn't check the proof signatures.
func (c *MisbehaviourVerifier) validateFault(p *Proof, committee types.Committee) []byte {
	valid := false
	switch p.Rule {
	case autonity.PN:
		valid = c.validMisbehaviourOfPN(p, committee)
	case autonity.PO:
		valid = c.validMisbehaviourOfPO(p, committee)
	case autonity.PVN:
		valid = c.validMisbehaviourOfPVN(p, committee)
	case autonity.PVO:
		valid = c.validMisbehaviourOfPVO(p, committee)
	case autonity.PVO12:
		valid = c.validMisbehaviourOfPVO12(p, committee)
	case autonity.C:
		valid = c.validMisbehaviourOfC(p, committee)
	case autonity.InvalidProposer:
		if lightProposal, ok := p.Message.(*message.LightProposal); ok {
			valid = !isProposerValid(c.chain, lightProposal)
		}
	case autonity.Equivocation:
		valid = errors.Is(checkEquivocation(p.Message, p.Evidences), errEquivocation)
	default:
		valid = false
	}

	if valid {
		return validReturn(p.Message, p.Offender, p.Rule)
	}
	return failureReturn
}

// check if the Proof of challenge of PN is valid,
// node propose a new value when there is a Proof that it precommit at a different value at previous round.
func (c *MisbehaviourVerifier) validMisbehaviourOfPN(p *Proof, committee types.Committee) bool {
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

	if !preCommit.Signers().Contains(p.OffenderIndex) {
		return false
	}

	if committee[p.OffenderIndex].Address != p.Offender && p.Offender == proposal.Signer() &&
		preCommit.R() < proposal.R() &&
		preCommit.Value() != nilValue {
		return true
	}
	return false
}

// check if the Proof of challenge of PO is valid
func (c *MisbehaviourVerifier) validMisbehaviourOfPO(p *Proof, committee types.Committee) bool {
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
			vote.Signers().Contains(p.OffenderIndex) && committee[p.OffenderIndex].Address == proposal.Signer() &&
			p.Offender == proposal.Signer() && vote.Value() != nilValue && vote.Value() != proposal.Value() {
			return true
		}
		if vote.R() > proposal.ValidRound() &&
			vote.R() < proposal.R() &&
			vote.Signers().Contains(p.OffenderIndex) && committee[p.OffenderIndex].Address == proposal.Signer() &&
			p.Offender == proposal.Signer() &&
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

		if hasEquivocatedVotes(p.Evidences) || hasDuplicatedVotes(p.Evidences) {
			return false
		}

		// check if preVotes for a not V reaches to quorum.
		quorum := bft.Quorum(committee.TotalVotingPower())
		return engineCore.OverQuorumVotes(p.Evidences, quorum) != nil

	}
	return false
}

// check if the Proof of challenge of PVN is valid.
func (c *MisbehaviourVerifier) validMisbehaviourOfPVN(p *Proof, committee types.Committee) bool {
	if len(p.Evidences) == 0 {
		return false
	}
	prevote := p.Message.(*message.Prevote)
	present := prevote.Signers().Contains(p.OffenderIndex)
	if !present || prevote.Code() != message.PrevoteCode || prevote.Value() == nilValue ||
		committee[p.OffenderIndex].Address != p.Offender {
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
		preC := pc.(*message.Precommit)
		if pc.Code() != message.PrecommitCode || !preC.Signers().Contains(p.OffenderIndex) || pc.R() >= prevote.R() {
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
func (c *MisbehaviourVerifier) validMisbehaviourOfPVO(p *Proof, committee types.Committee) bool {
	if len(p.Evidences) < 2 {
		return false
	}
	prevote := p.Message.(*message.Prevote)
	present := prevote.Signers().Contains(p.OffenderIndex)
	if !present || prevote.Code() != message.PrevoteCode || prevote.Value() == nilValue ||
		committee[p.OffenderIndex].Address != p.Offender {
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

	if hasEquivocatedVotes(p.Evidences[1:]) || hasDuplicatedVotes(p.Evidences[1:]) {
		return false
	}

	// check if quorum prevote for a different value than V at valid round.
	quorum := bft.Quorum(committee.TotalVotingPower())
	return engineCore.OverQuorumVotes(p.Evidences[1:], quorum) != nil
}

// check if the Proof of challenge of PVO12 is valid.
func (c *MisbehaviourVerifier) validMisbehaviourOfPVO12(p *Proof, committee types.Committee) bool {
	if len(p.Evidences) < 2 {
		return false
	}
	prevote := p.Message.(*message.Prevote)
	present := prevote.Signers().Contains(p.OffenderIndex)
	if !present || prevote.Code() != message.PrevoteCode || prevote.Value() == nilValue ||
		committee[p.OffenderIndex].Address != p.Offender {
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
	// check if all precommits between range (validRound, currentRound) are presented.
	// There might have multiple precommits per round due to overlapped aggregation.
	presentedRounds := make(map[int64]struct{})
	for _, pc := range allPreCommits {
		preC := pc.(*message.Precommit)
		if pc.R() <= validRound || pc.R() >= currentRound || pc.Code() != message.PrecommitCode ||
			!preC.Signers().Contains(p.OffenderIndex) || pc.H() != prevote.H() {
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
func (c *MisbehaviourVerifier) validMisbehaviourOfC(p *Proof, committee types.Committee) bool {
	if len(p.Evidences) == 0 {
		return false
	}
	preCommit := p.Message.(*message.Precommit)
	present := preCommit.Signers().Contains(p.OffenderIndex)
	if !present || preCommit.Code() != message.PrecommitCode || preCommit.Value() == nilValue ||
		committee[p.OffenderIndex].Address != p.Offender {
		return false
	}

	// check preVotes for not the same V compares to preCommit.
	for _, m := range p.Evidences {
		if m.Code() != message.PrevoteCode || m.Value() == preCommit.Value() || m.R() != preCommit.R() {
			return false
		}
	}

	// check no redundant vote msg in evidence in case of hacking.
	if hasEquivocatedVotes(p.Evidences) || hasDuplicatedVotes(p.Evidences) {
		return false
	}

	// check if preVotes for not V reaches to quorum.
	quorum := bft.Quorum(committee.TotalVotingPower())
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
// return the rlp hash of msg and the rlp hash of msg signer as the valid identity for on-chain management of proofs,
// AC need the check the value returned to match the ID which is on challenge, to remove the challenge from chain.
func (c *InnocenceVerifier) Run(input []byte, blockNumber uint64, _ *vm.EVM, _ common.Address) ([]byte, error) {
	if len(input) <= 32 || blockNumber == 0 {
		return failureReturn, nil
	}
	// the 1st 32 bytes are length of bytes array in solidity, take RLP bytes after it.
	p, err := decodeRawProof(input[32:])
	if err != nil {
		return failureReturn, nil
	}
	committee, err := verifyProofSignatures(c.chain, p)
	if err != nil {
		return failureReturn, nil
	}

	if !verifyInnocenceProof(p, committee) {
		return failureReturn, nil
	}
	return validReturn(p.Message, p.Offender, p.Rule), nil
}

func verifyInnocenceProof(p *Proof, committee types.Committee) bool {
	// rule engine only have 4 kind of provable accusation for the time being.
	switch p.Rule {
	case autonity.PO:
		return validInnocenceProofOfPO(p, committee)
	case autonity.PVN:
		return validInnocenceProofOfPVN(p)
	case autonity.PVO:
		return validInnocenceProofOfPVO(p, committee)
	case autonity.C1:
		return validInnocenceProofOfC1(p, committee)
	default:
		return false
	}
}

// check if the Proof of innocent of PO is valid.
func validInnocenceProofOfPO(p *Proof, committee types.Committee) bool {
	// check if there is quorum number of prevote at the same value on the same valid round
	proposal := p.Message
	if proposal.Code() != message.LightProposalCode || proposal.(*message.LightProposal).ValidRound() == -1 {
		return false
	}

	// check the votes match for the corresponding proposal, and there is no vote for other value in the proof.
	for _, m := range p.Evidences {
		if !(m.Code() == message.PrevoteCode &&
			m.Value() == proposal.Value() &&
			m.R() == proposal.(*message.LightProposal).ValidRound()) {
			return false
		}
	}

	// check no redundant vote msg in evidence in case of hacking.
	if hasDuplicatedVotes(p.Evidences) {
		return false
	}

	// check quorum prevotes for V at validRound.
	quorum := bft.Quorum(committee.TotalVotingPower())
	return engineCore.OverQuorumVotes(p.Evidences, quorum) != nil
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
func validInnocenceProofOfPVO(p *Proof, committee types.Committee) bool {
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
	// check prevotes for V at the valid round, no vote for other value.
	for _, m := range p.Evidences[1:] {
		if !(m.Code() == message.PrevoteCode && m.Value() == proposal.Value() && m.R() == vr) {
			return false
		}
	}

	// check no redundant vote msg in evidence in case of hacking.
	if hasDuplicatedVotes(p.Evidences[1:]) {
		return false
	}

	// check quorum prevotes at valid round.
	quorum := bft.Quorum(committee.TotalVotingPower())
	return engineCore.OverQuorumVotes(p.Evidences[1:], quorum) != nil
}

// check if the Proof of innocent of C1 is valid.
func validInnocenceProofOfC1(p *Proof, committee types.Committee) bool {
	preCommit, ok := p.Message.(*message.Precommit)
	if !ok {
		return false
	}
	if preCommit.Value() == nilValue {
		return false
	}
	// check quorum prevotes for V at the same round, there is no vote for other value.
	for _, m := range p.Evidences {
		if !(m.Code() == message.PrevoteCode && m.Value() == preCommit.Value() &&
			m.R() == preCommit.R()) {
			return false
		}
	}
	// check no redundant vote msg in evidence in case of hacking.
	if hasDuplicatedVotes(p.Evidences) {
		return false
	}
	quorum := bft.Quorum(committee.TotalVotingPower())
	return engineCore.OverQuorumVotes(p.Evidences, quorum) != nil
}

// check if there is duplicated vote messages in the set.
func hasDuplicatedVotes(votes []message.Msg) bool {
	hashMap := make(map[common.Hash]struct{})
	for _, vote := range votes {
		_, ok := hashMap[vote.Hash()]
		if !ok {
			hashMap[vote.Hash()] = struct{}{}
		} else {
			return true
		}
	}
	return false
}

// check if there are votes for different values in the set
func hasEquivocatedVotes(votes []message.Msg) bool {
	if len(votes) <= 1 {
		return false
	}
	value := votes[0].Value()
	for _, vote := range votes[1:] {
		if value != vote.Value() {
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

// verifyProofSignatures checks if the consensus message is from valid member of the committee.
func verifyProofSignatures(chain ChainContext, p *Proof) (types.Committee, error) {
	h := p.Message.H()
	lastHeader := chain.GetHeaderByNumber(h - 1)
	if lastHeader == nil {
		return nil, errFutureMsg
	}

	// before signature verification, check if the offender index is valid
	if p.OffenderIndex >= len(lastHeader.Committee) || p.OffenderIndex < 0 {
		return nil, errInvalidOffenderIdx
	}

	// assign power and bls signer key
	if err := p.Message.PreValidate(lastHeader); err != nil {
		return nil, err
	}

	// verify signature
	if err := p.Message.Validate(); err != nil {
		return nil, errNotCommitteeMsg
	}

	for _, msg := range p.Evidences {
		if msg.H() != h {
			return nil, errBadHeight
		}

		if err := msg.PreValidate(lastHeader); err != nil {
			return nil, err
		}

		if err := msg.Validate(); err != nil {
			return nil, err
		}
	}
	return lastHeader.Committee, nil
}

func checkEquivocation(m message.Msg, proof []message.Msg) error {
	if len(proof) == 0 {
		return fmt.Errorf("no proof")
	}
	// todo: fix this bug!!!
	// check equivocations.
	if m.Hash() != proof[0].Hash() {
		return errEquivocation
	}
	return nil
}

func validReturn(m message.Msg, signer common.Address, rule autonity.Rule) []byte {
	offender := common.LeftPadBytes(signer.Bytes(), 32)
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
