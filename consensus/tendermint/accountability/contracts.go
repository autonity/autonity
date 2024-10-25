package accountability

import (
	"errors"
	"math"
	"math/big"
	"strconv"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/params/generated"

	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/core/types"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/bft"
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
	errProofOffender       = errors.New("accountability proof contains invalid offender")
	errProofMsgCode        = errors.New("accountability proof contains invalid msg code")
	errMaxEvidences        = errors.New("above max evidence threshold")
)

const KB = 1024

// LoadPrecompiles init the instances of Fault Detector contracts, and register them into EVM's context
func LoadPrecompiles() {
	vm.PrecompiledContractRWMutex.Lock()
	defer vm.PrecompiledContractRWMutex.Unlock()
	pv := InnocenceVerifier{checkInnocenceAddress}
	cv := MisbehaviourVerifier{checkMisbehaviourAddress}
	av := AccusationVerifier{checkAccusationAddress}
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
	address common.Address
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
// the rlp hash of the msg payload and the msg signer is returned.
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

	committee, err := committeeByHeight(p.Message.H(), e, a.address)
	if err != nil {
		return failureReturn, err
	}

	if err = verifyProofSignatures(committee, p); err != nil {
		return failureReturn, nil
	}
	if verifyAccusation(p, committee) {
		// the proof carry valid info.
		return validReturn(p.Message, committee.Members[p.OffenderIndex].Address, p.Rule), nil
	}
	return failureReturn, nil
}

// validate the submitted accusation by the contract call.
func verifyAccusation(p *Proof, committee *types.Committee) bool {
	// we have only 4 types of rule on accusation.
	// an improvement of this function would be to return an error instead of a bool
	switch p.Rule {
	case autonity.PO:
		if p.Message.Code() != message.LightProposalCode {
			return false
		}
		lightProposal, ok := p.Message.(*message.LightProposal)
		if !ok {
			return false
		}
		if lightProposal.ValidRound() == -1 || committee.Members[p.OffenderIndex].Address != lightProposal.Signer() {
			return false
		}

	case autonity.PVN:
		if p.Message.Code() != message.PrevoteCode || p.Message.Value() == nilValue {
			return false
		}
		prevote, ok := p.Message.(*message.Prevote)
		if !ok {
			return false
		}
		if !prevote.Signers().Contains(p.OffenderIndex) {
			return false
		}

	case autonity.PVO:
		// theoretically we do not need the non-nil check, since we will check later that prevote.value == proposal.value
		// however added for simplicity of understanding
		if p.Message.Code() != message.PrevoteCode || p.Message.Value() == nilValue {
			return false
		}
		prevote, ok := p.Message.(*message.Prevote)
		if !ok {
			return false
		}
		if !prevote.Signers().Contains(p.OffenderIndex) {
			return false
		}

	case autonity.C1:
		if p.Message.Code() != message.PrecommitCode || p.Message.Value() == nilValue {
			return false
		}

		precommit, ok := p.Message.(*message.Precommit)
		if !ok {
			return false
		}
		if !precommit.Signers().Contains(p.OffenderIndex) {
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

		oldProposal, ok := p.Evidences[0].(*message.LightProposal)
		if !ok {
			return false
		}

		if oldProposal.Code() != message.LightProposalCode ||
			oldProposal.R() != p.Message.R() ||
			oldProposal.Value() != p.Message.Value() ||
			oldProposal.ValidRound() == -1 {
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
	address common.Address
}

// RequiredGas the gas cost to execute MisbehaviourVerifier contract, weighted by input data size.
func (c *MisbehaviourVerifier) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

// Run take the rlp encoded Proof of challenge in byte array, decode it and validate it, if the Proof is valid, then
// the rlp hash of the msg payload and the msg signer is returned as the valid identity for Proof management.
func (c *MisbehaviourVerifier) Run(input []byte, _ uint64, e *vm.EVM, _ common.Address) ([]byte, error) {
	if len(input) <= 32 {
		return failureReturn, nil
	}
	// the 1st 32 bytes are length of bytes array in solidity, take RLP bytes after it.
	p, err := decodeRawProof(input[32:])
	if err != nil {
		return failureReturn, nil
	}

	committee, err := committeeByHeight(p.Message.H(), e, c.address)
	if err != nil {
		return failureReturn, err
	}

	if err = verifyProofSignatures(committee, p); err != nil {
		return failureReturn, nil
	}
	return c.validateFault(p, committee), nil
}

// validate a misbehavior proof, doesn't check the proof signatures.
func (c *MisbehaviourVerifier) validateFault(p *Proof, committee *types.Committee) []byte {
	valid := false
	switch p.Rule {
	case autonity.PN:
		valid = c.validMisbehaviourOfPN(p, committee)
	case autonity.PO:
		valid = c.validMisbehaviourOfPO(p, committee)
	case autonity.PVN:
		valid = c.validMisbehaviourOfPVN(p)
	case autonity.PVO:
		valid = c.validMisbehaviourOfPVO(p, committee)
	case autonity.PVO12:
		valid = c.validMisbehaviourOfPVO12(p)
	case autonity.C:
		valid = c.validMisbehaviourOfC(p, committee)
	case autonity.InvalidProposer:
		if lightProposal, ok := p.Message.(*message.LightProposal); ok {
			proposer := committee.Proposer(lightProposal.H()-1, lightProposal.R())
			valid = (proposer != lightProposal.Signer()) && (committee.Members[p.OffenderIndex].Address == lightProposal.Signer())

		}
	case autonity.Equivocation:
		valid = validMisbehaviourOfEquivocation(p, committee)
	default:
		valid = false
	}

	if valid {
		return validReturn(p.Message, committee.Members[p.OffenderIndex].Address, p.Rule)
	}
	return failureReturn
}

// check if the Proof of challenge of PN is valid,
// node propose a new value when there is a Proof that it precommit at a different value at previous round.
func (c *MisbehaviourVerifier) validMisbehaviourOfPN(p *Proof, committee *types.Committee) bool {
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

	if committee.Members[p.OffenderIndex].Address == proposal.Signer() &&
		preCommit.R() < proposal.R() && preCommit.Value() != nilValue {
		return true
	}
	return false
}

// check if the Proof of challenge of PO is valid
func (c *MisbehaviourVerifier) validMisbehaviourOfPO(p *Proof, committee *types.Committee) bool {
	// should be an old proposal
	proposal, ok := p.Message.(*message.LightProposal)
	if !ok {
		return false
	}
	if proposal.ValidRound() == -1 {
		return false
	}

	if len(p.Evidences) == 0 {
		return false
	}

	switch vote := p.Evidences[0].(type) {
	case *message.Precommit:
		if vote.R() == proposal.ValidRound() &&
			vote.Signers().Contains(p.OffenderIndex) && committee.Members[p.OffenderIndex].Address == proposal.Signer() &&
			vote.Value() != nilValue && vote.Value() != proposal.Value() {
			return true
		}
		if vote.R() > proposal.ValidRound() &&
			vote.R() < proposal.R() &&
			vote.Signers().Contains(p.OffenderIndex) && committee.Members[p.OffenderIndex].Address == proposal.Signer() &&
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

		if hasDifferentVoteOfValues(p.Evidences) || hasDuplicatedVotes(p.Evidences) {
			return false
		}

		// check if preVotes for a not V reaches to quorum.
		quorum := bft.Quorum(committee.TotalVotingPower())
		return message.OverQuorumVotes(p.Evidences, quorum) != nil

	}
	return false
}

// check if the Proof of challenge of PVN is valid.
func (c *MisbehaviourVerifier) validMisbehaviourOfPVN(p *Proof) bool {
	// if there is no evidence or there is expected number of msg evidence field, return false to prevent DoS attack.
	if len(p.Evidences) == 0 || len(p.Evidences) > 2 {
		return false
	}

	prevote, ok := p.Message.(*message.Prevote)
	if !ok {
		return false
	}

	present := prevote.Signers().Contains(p.OffenderIndex)
	if !present || prevote.Code() != message.PrevoteCode || prevote.Value() == nilValue {
		return false
	}

	// check if the corresponding new proposal of preVote for new value is presented.
	correspondingProposal, ok := p.Evidences[0].(*message.LightProposal)
	if !ok {
		return false
	}

	if correspondingProposal.Code() != message.LightProposalCode ||
		correspondingProposal.H() != prevote.H() ||
		correspondingProposal.R() != prevote.R() ||
		correspondingProposal.Value() != prevote.Value() ||
		correspondingProposal.ValidRound() != -1 {
		return false
	}

	// validate preCommits from round R' to R, make sure there is no gap, and the value preCommitted at R'
	// is different from the value preVoted, and the other ones are preCommits of nil.
	preCommits := p.Evidences[1:]

	// Single precommit as the evidence field, process it.
	if len(preCommits) == 1 {
		pc := preCommits[0]
		preC, ok := pc.(*message.Precommit)
		if !ok {
			return false
		}
		if pc.Code() != message.PrecommitCode || !preC.Signers().Contains(p.OffenderIndex) || pc.R() >= prevote.R() {
			return false
		}
		return pc.R()+1 == prevote.R() && pc.Value() != nilValue && pc.Value() != prevote.Value()
	}

	// Otherwise, we have to process aggregated precommits from the aggregated precommits.
	if p.DistinctPrecommits.Len() > 0 {
		preCommits := p.DistinctPrecommits.MsgSigners
		lastIndex := len(preCommits) - 1
		for i, pc := range preCommits {
			// as height and msg code was checked at validate phase, thus we just check round and values at below.
			if !pc.Contains(p.OffenderIndex) || pc.Round >= prevote.R() {
				return false
			}

			// preCommit at R'
			if i == 0 {
				if pc.Value == nilValue || pc.Value == prevote.Value() {
					return false
				}
			} else {
				// preCommits at between R' and R-1, they should be nil.
				if pc.Value != nilValue {
					return false
				}
			}

			// check if there is round gaps between R' and R-1.
			if i < lastIndex && preCommits[i+1].Round-pc.Round > 1 {
				return false
			}

			// check round gap for preCommit at R-1 and R.
			if i == lastIndex {
				return pc.Round+1 == prevote.R()
			}
		}
	}

	return false
}

// check if the proof of challenge of PVO is valid.
func (c *MisbehaviourVerifier) validMisbehaviourOfPVO(p *Proof, committee *types.Committee) bool {
	if len(p.Evidences) < 2 {
		return false
	}
	prevote, ok := p.Message.(*message.Prevote)
	if !ok {
		return false
	}

	present := prevote.Signers().Contains(p.OffenderIndex)
	if !present || prevote.Code() != message.PrevoteCode || prevote.Value() == nilValue {
		return false
	}
	// check if the corresponding proposal of preVote is presented.
	correspondingProposal, ok := p.Evidences[0].(*message.LightProposal)
	if !ok {
		return false
	}

	if correspondingProposal.Code() != message.LightProposalCode ||
		correspondingProposal.H() != prevote.H() ||
		correspondingProposal.R() != prevote.R() ||
		correspondingProposal.Value() != prevote.Value() ||
		correspondingProposal.ValidRound() == -1 {
		return false
	}

	validRound := correspondingProposal.ValidRound()
	votedVatVR := p.Evidences[1].Value()

	// check preVotes at evidence.
	for _, pv := range p.Evidences[1:] {
		if _, ok := pv.(*message.Prevote); !ok {
			return false
		}

		if pv.Code() != message.PrevoteCode || pv.R() != validRound || pv.Value() == nilValue ||
			pv.Value() == correspondingProposal.Value() || pv.Value() != votedVatVR {
			return false
		}
	}

	if hasDifferentVoteOfValues(p.Evidences[1:]) || hasDuplicatedVotes(p.Evidences[1:]) {
		return false
	}

	// check if quorum prevote for a different value than V at valid round.
	quorum := bft.Quorum(committee.TotalVotingPower())
	return message.OverQuorumVotes(p.Evidences[1:], quorum) != nil
}

// check if the Proof of challenge of PVO12 is valid.
func (c *MisbehaviourVerifier) validMisbehaviourOfPVO12(p *Proof) bool {
	// if there is no evidence or there is unexpected number of msg evidence field, return false to prevent DoS attack.
	if len(p.Evidences) == 0 || len(p.Evidences) > 2 {
		return false
	}

	prevote, ok := p.Message.(*message.Prevote)
	if !ok {
		return false
	}

	present := prevote.Signers().Contains(p.OffenderIndex)
	if !present || prevote.Code() != message.PrevoteCode || prevote.Value() == nilValue {
		return false
	}

	// check if the corresponding proposal of preVote.
	correspondingProposal, ok := p.Evidences[0].(*message.LightProposal)
	if !ok {
		return false
	}

	if correspondingProposal.Code() != message.LightProposalCode ||
		correspondingProposal.H() != prevote.H() ||
		correspondingProposal.R() != prevote.R() ||
		correspondingProposal.Value() != prevote.Value() ||
		correspondingProposal.ValidRound() == -1 {
		return false
	}

	currentRound := correspondingProposal.R()
	validRound := correspondingProposal.ValidRound()

	precommits := p.Evidences[1:]
	// Single precommit as the evidence field, process it.
	if len(precommits) == 1 {
		preC, ok := precommits[0].(*message.Precommit)
		if !ok {
			return false
		}
		// check if the msg is in range (validRound, currentRound), and with correct signer, height and code.
		if preC.R() <= validRound || preC.R() >= currentRound || preC.Code() != message.PrecommitCode ||
			!preC.Signers().Contains(p.OffenderIndex) || preC.H() != prevote.H() ||
			int(currentRound-validRound)-1 != len(precommits) {
			return false
		}
		return preC.Value() != prevote.Value() && preC.Value() != nilValue
	}

	// we have distinct aggregated precomits
	if p.DistinctPrecommits.Len() > 0 {

		allPreCommits := p.DistinctPrecommits.MsgSigners
		// check if there are any msg out of range (validRound, currentRound), and with correct address, height and code.
		// check if all precommits between range (validRound, currentRound) are presented.
		// There might have multiple precommits per round due to overlapped aggregation.
		presentedRounds := make(map[int64]struct{})
		for _, pc := range allPreCommits {
			// as the height and code was check at msg validation phase, thus we just check round and signers at below.
			if pc.Round <= validRound || pc.Round >= currentRound || !pc.Contains(p.OffenderIndex) {
				return false
			}
			presentedRounds[pc.Round] = struct{}{}
		}

		if len(presentedRounds) != int(currentRound-validRound)-1 {
			return false
		}

		// If the last precommit for notV is after the last one for V, raise misbehaviour
		// If all precommits are nil, do not raise misbehaviour. It is a valid correct scenario.
		lastRoundForV := int64(-1)
		lastRoundForNotV := int64(-1)
		for _, pc := range allPreCommits {
			if pc.Value == prevote.Value() && pc.Round > lastRoundForV {
				lastRoundForV = pc.Round
			}

			if pc.Value != prevote.Value() && pc.Value != nilValue && pc.Round > lastRoundForNotV {
				lastRoundForNotV = pc.Round
			}
		}

		return lastRoundForNotV > lastRoundForV
	}
	return false
}

// check if the Proof of challenge of C is valid.
func (c *MisbehaviourVerifier) validMisbehaviourOfC(p *Proof, committee *types.Committee) bool {
	if len(p.Evidences) == 0 {
		return false
	}
	preCommit, ok := p.Message.(*message.Precommit)
	if !ok {
		return false
	}
	present := preCommit.Signers().Contains(p.OffenderIndex)
	if !present || preCommit.Code() != message.PrecommitCode || preCommit.Value() == nilValue {
		return false
	}

	// check preVotes for not the same V compares to preCommit.
	for _, m := range p.Evidences {
		if _, ok := m.(*message.Prevote); !ok {
			return false
		}
		if m.Code() != message.PrevoteCode || m.Value() == preCommit.Value() || m.R() != preCommit.R() {
			return false
		}
	}

	// check no redundant vote msg in evidence in case of hacking.
	if hasDifferentVoteOfValues(p.Evidences) || hasDuplicatedVotes(p.Evidences) {
		return false
	}

	// check if preVotes for not V reaches to quorum.
	quorum := bft.Quorum(committee.TotalVotingPower())
	return message.OverQuorumVotes(p.Evidences, quorum) != nil
}

// InnocenceVerifier implemented as a native contract to validate an innocence Proof.
type InnocenceVerifier struct {
	address common.Address
}

// RequiredGas the gas cost to execute this Proof validator contract, weighted by input data size.
func (c *InnocenceVerifier) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

// Run InnocenceVerifier, take the rlp encoded Proof of innocence, decode it and validate it, if the Proof is valid, then
// return the rlp hash of msg and the rlp hash of msg signer as the valid identity for on-chain management of proofs,
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

	committee, err := committeeByHeight(p.Message.H(), e, c.address)
	if err != nil {
		return failureReturn, err
	}

	if err = verifyProofSignatures(committee, p); err != nil {
		return failureReturn, nil
	}

	if !verifyInnocenceProof(p, committee) {
		return failureReturn, nil
	}
	return validReturn(p.Message, committee.Members[p.OffenderIndex].Address, p.Rule), nil
}

func verifyInnocenceProof(p *Proof, committee *types.Committee) bool {
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
func validInnocenceProofOfPO(p *Proof, committee *types.Committee) bool {
	// check if there is quorum number of prevote at the same value on the same valid round
	proposal, ok := p.Message.(*message.LightProposal)
	if !ok {
		return false
	}

	if proposal.Code() != message.LightProposalCode || proposal.ValidRound() == -1 {
		return false
	}

	// check the votes match for the corresponding proposal, and there is no vote for other value in the proof.
	for _, m := range p.Evidences {
		if _, ok := m.(*message.Prevote); !ok {
			return false
		}

		if !(m.Code() == message.PrevoteCode &&
			m.Value() == proposal.Value() &&
			m.R() == proposal.ValidRound()) {
			return false
		}
	}

	// check no redundant vote msg in evidence in case of hacking.
	if hasDuplicatedVotes(p.Evidences) {
		return false
	}

	// check quorum prevotes for V at validRound.
	quorum := bft.Quorum(committee.TotalVotingPower())
	return message.OverQuorumVotes(p.Evidences, quorum) != nil
}

// check if the Proof of innocent of PVN is valid.
func validInnocenceProofOfPVN(p *Proof) bool {
	preVote, ok := p.Message.(*message.Prevote)
	if !ok {
		return false
	}
	if !(preVote.Code() == message.PrevoteCode && preVote.Value() != nilValue) {
		return false
	}

	if len(p.Evidences) != 1 {
		return false
	}

	proposal, ok := p.Evidences[0].(*message.LightProposal)
	if !ok {
		return false
	}
	return proposal.Code() == message.LightProposalCode &&
		proposal.H() == preVote.H() &&
		proposal.R() == preVote.R() &&
		proposal.ValidRound() == -1 &&
		proposal.Value() == preVote.Value()
}

// check if the Proof of innocent of PVO is valid.
func validInnocenceProofOfPVO(p *Proof, committee *types.Committee) bool {
	// check if there is quorum number of prevote at the same value on the same valid round
	preVote, ok := p.Message.(*message.Prevote)
	if !ok {
		return false
	}
	if !(preVote.Code() == message.PrevoteCode && preVote.Value() != nilValue) {
		return false
	}

	if len(p.Evidences) <= 1 {
		return false
	}

	proposal, ok := p.Evidences[0].(*message.LightProposal)
	if !ok {
		return false
	}

	if proposal.Code() != message.LightProposalCode ||
		proposal.Value() != preVote.Value() ||
		proposal.R() != preVote.R() ||
		proposal.ValidRound() == -1 {
		return false
	}

	vr := proposal.ValidRound()
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
	return message.OverQuorumVotes(p.Evidences[1:], quorum) != nil
}

// check if the Proof of innocent of C1 is valid.
func validInnocenceProofOfC1(p *Proof, committee *types.Committee) bool {
	preCommit, ok := p.Message.(*message.Precommit)
	if !ok {
		return false
	}
	if preCommit.Value() == nilValue {
		return false
	}
	// check quorum prevotes for V at the same round, there is no vote for other value.
	for _, m := range p.Evidences {
		if _, ok := m.(*message.Prevote); !ok {
			return false
		}
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
	return message.OverQuorumVotes(p.Evidences, quorum) != nil
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
func hasDifferentVoteOfValues(votes []message.Msg) bool {
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
	if p.Message == nil {
		return p, errors.New("invalid proof")
	}
	return p, nil
}

// verifyProofSignatures checks if the consensus message is from valid member of the committee.
func verifyProofSignatures(committee *types.Committee, p *Proof) error {
	// before signature verification, check if the offender index is valid
	if p.OffenderIndex >= committee.Len() || p.OffenderIndex < 0 {
		return errInvalidOffenderIdx
	}

	// assign power and bls signer key
	if err := p.Message.PreValidate(committee); err != nil {
		return err
	}

	// verify signature
	if err := p.Message.Validate(); err != nil {
		return errNotCommitteeMsg
	}

	// check if the number of evidence msgs are exceeded the max to prevent the abuse of the proof msg.
	if len(p.Evidences) > maxEvidenceMessages(committee.Len()) {
		return errMaxEvidences
	}

	h := p.Message.H()
	for _, msg := range p.Evidences {
		if msg.H() != h {
			return errBadHeight
		}

		if err := msg.PreValidate(committee); err != nil {
			return err
		}

		if err := msg.Validate(); err != nil {
			return err
		}
	}

	// pre-validate and validate the highly aggregated precommits.
	if p.DistinctPrecommits.Len() > 0 {
		if err := p.DistinctPrecommits.PreValidate(committee, h); err != nil {
			return err
		}
		if err := p.DistinctPrecommits.Validate(); err != nil {
			return err
		}
	}

	// check offender idx match with the p.Message.Signer() or Signers().Has(offenderIdx)
	switch m := p.Message.(type) {
	case *message.LightProposal:
		if committee.Members[p.OffenderIndex].Address != m.Signer() {
			return errProofOffender
		}

	case *message.Prevote, *message.Precommit:
		vote1 := p.Message.(message.Vote)
		if !vote1.Signers().Contains(p.OffenderIndex) {
			return errProofOffender
		}
	default:
		return errProofMsgCode
	}
	return nil
}

func validMisbehaviourOfEquivocation(proof *Proof, committee *types.Committee) bool {
	if len(proof.Evidences) != 1 {
		return false
	}

	// as the presents of proof.Message was checked, we can check if the equivocated message have the same msg code.
	if proof.Message.Code() != proof.Evidences[0].Code() {
		return false
	}

	switch msg1 := proof.Message.(type) {
	case *message.LightProposal:
		// check for equivocated proposal with light proposals
		msg2, ok := proof.Evidences[0].(*message.LightProposal)
		if !ok {
			return false
		}

		if msg1.H() == msg2.H() && msg1.R() == msg2.R() && msg1.Signer() == msg2.Signer() &&
			msg1.Signer() == committee.Members[proof.OffenderIndex].Address &&
			(msg1.Value() != msg2.Value() || msg1.ValidRound() != msg2.ValidRound()) {
			return true
		}

	case *message.Prevote, *message.Precommit:
		// check for equivocated proposal with votes
		vote1 := proof.Message.(message.Vote)
		if !vote1.Signers().Contains(proof.OffenderIndex) {
			return false
		}

		vote2, ok := proof.Evidences[0].(message.Vote)
		if !ok {
			return false
		}
		if !vote2.Signers().Contains(proof.OffenderIndex) {
			return false
		}

		if vote1.H() == vote2.H() && vote1.R() == vote2.R() && msg1.Value() != vote2.Value() {
			return true
		}

	default:
		return false
	}

	return false
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

func maxEvidenceMessages(committeeSize int) int {
	if committeeSize > constants.MaxRound {
		return committeeSize + 1
	}
	return constants.MaxRound + 1
}

func committeeByHeight(height uint64, evm *vm.EVM, caller common.Address) (*types.Committee, error) {
	// as the gas of the precompile functions are resolved by its corresponding RequiredGas() interface, thus this gas
	// consumption of reading committee is not counted into the original TXN, thus it does not make sense to resolve
	// the gas cap for this call.
	gas := uint64(math.MaxUint64)
	packedArgs, err := generated.AutonityAbi.Pack("getEpochByHeight", new(big.Int).SetUint64(height))
	if err != nil {
		return nil, err
	}
	// Todo(scott): consider using the existing call on the for callEpochByHeight from the autonity package
	ret, _, err := evm.Call(vm.AccountRef(caller), params.AutonityContractAddress, packedArgs, gas, new(big.Int))
	if err != nil {
		return nil, err
	}

	data, err := generated.AutonityAbi.Unpack("getEpochByHeight", ret)
	if err != nil {
		return nil, err
	}

	info := *abi.ConvertType(data[0], new(autonity.AutonityEpochInfo)).(*autonity.AutonityEpochInfo)

	if len(info.Committee) == 0 {
		panic("get empty committee set for height: " + strconv.FormatUint(height, 10))
	}

	committee := &types.Committee{}
	for _, member := range info.Committee {
		committee.Members = append(committee.Members, types.CommitteeMember{
			Address:           member.Addr,
			VotingPower:       member.VotingPower,
			ConsensusKeyBytes: member.ConsensusKey,
		})
	}
	// As the committee is already sorted by the contract, thus we don't need sort again.
	if err := committee.Enrich(); err != nil {
		return nil, err
	}

	return committee, nil
}
