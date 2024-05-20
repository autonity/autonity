package accountability

import (
	"errors"
	"fmt"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
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
	errMaxEvidences        = errors.New("above max evidence threshold")
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
// the rlp hash of the msg payload and the msg sender is returned.
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

	if err := verifyProofSignatures(a.chain, p); err != nil {
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
	chain ChainContext
}

// RequiredGas the gas cost to execute MisbehaviourVerifier contract, weighted by input data size.
func (c *MisbehaviourVerifier) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

// Run take the rlp encoded Proof of challenge in byte array, decode it and validate it, if the Proof is valid, then
// the rlp hash of the msg payload and the msg sender is returned as the valid identity for Proof management.
func (c *MisbehaviourVerifier) Run(input []byte, _ uint64, _ *vm.EVM, _ common.Address) ([]byte, error) {
	if len(input) <= 32 {
		return failureReturn, nil
	}
	// the 1st 32 bytes are length of bytes array in solidity, take RLP bytes after it.
	p, err := decodeRawProof(input[32:])
	if err != nil {
		return failureReturn, nil
	}
	if err := verifyProofSignatures(c.chain, p); err != nil {
		return failureReturn, nil
	}
	return c.validateFault(p), nil
}

// validate a misbehavior proof, doesn't check the proof signatures.
func (c *MisbehaviourVerifier) validateFault(p *Proof) []byte {
	valid := false
	switch p.Rule {
	case autonity.PN:
		valid = c.validMisbehaviourOfPN(p)
	case autonity.PO:
		valid = c.validMisbehaviourOfPO(p)
	case autonity.PVN:
		valid = c.validMisbehaviourOfPVN(p)
	case autonity.PVO:
		valid = c.validMisbehaviourOfPVO(p)
	case autonity.PVO12:
		valid = c.validMisbehaviourOfPVO12(p)
	case autonity.C:
		valid = c.validMisbehaviourOfC(p)
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
	/* TODO(lorenzo) refinements2, fix
	preCommit, ok := p.Evidences[0].(*message.Precommit)
	if !ok {
		return false
	}
	if preCommit.Signer() == proposal.Signer() &&
		preCommit.R() < proposal.R() &&
		preCommit.Value() != nilValue {
		return true
	}
	*/
	return false
}

// check if the Proof of challenge of PO is valid
func (c *MisbehaviourVerifier) validMisbehaviourOfPO(p *Proof) bool {
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
			//	vote.Signer() == p.Message.Signer() && //TODO(lorenzo) fix this
			vote.Value() != nilValue &&
			vote.Value() != proposal.Value() {
			return true
		}
		if vote.R() > proposal.ValidRound() &&
			vote.R() < proposal.R() &&
			//vote.Signer() == p.Message.Signer() && //TODO(lorenzo) fix this
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
		// check if preVotes for not V reaches to quorum.
		lastHeader := c.chain.GetHeaderByNumber(p.Message.H() - 1)
		if lastHeader == nil {
			return false
		}
		quorum := bft.Quorum(lastHeader.TotalVotingPower())
		return message.OverQuorumVotes(p.Evidences, quorum) != nil

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
		/* //TODO(lorenzo) fix
		if pc.Code() != message.PrecommitCode || pc.Signer() != prevote.Signer() || pc.R() >= prevote.R() {
			return false
		}*/

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
	lastHeader := c.chain.GetHeaderByNumber(p.Message.H() - 1)
	if lastHeader == nil {
		return false
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())

	return message.OverQuorumVotes(p.Evidences[1:], quorum) != nil
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
		/* //TODO(lorenzo) fix this
		if pc.R() <= validRound || pc.R() >= currentRound || pc.Code() != message.PrecommitCode || pc.Signer() != prevote.Signer() ||
			pc.H() != prevote.H() {
			return false
		}*/
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
func (c *MisbehaviourVerifier) validMisbehaviourOfC(p *Proof) bool {
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
	lastHeader := c.chain.GetHeaderByNumber(p.Message.H() - 1)
	if lastHeader == nil {
		return false
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())

	return message.OverQuorumVotes(p.Evidences, quorum) != nil
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
func (c *InnocenceVerifier) Run(input []byte, blockNumber uint64, _ *vm.EVM, _ common.Address) ([]byte, error) {
	if len(input) <= 32 || blockNumber == 0 {
		return failureReturn, nil
	}
	// the 1st 32 bytes are length of bytes array in solidity, take RLP bytes after it.
	p, err := decodeRawProof(input[32:])
	if err != nil {
		return failureReturn, nil
	}
	if err := verifyProofSignatures(c.chain, p); err != nil {
		return failureReturn, nil
	}
	if !verifyInnocenceProof(p, c.chain) {
		return failureReturn, nil
	}
	return validReturn(p.Message, p.Rule), nil
}

func verifyInnocenceProof(p *Proof, chain ChainContext) bool {
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
	lastHeader := chain.GetHeaderByNumber(proposal.H() - 1)
	if lastHeader == nil {
		return false
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())
	return message.OverQuorumVotes(p.Evidences, quorum) != nil
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
func validInnocenceProofOfPVO(p *Proof, chain ChainContext) bool {
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
	height := preVote.H()
	lastHeader := chain.GetHeaderByNumber(height - 1)
	if lastHeader == nil {
		return false
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())
	return message.OverQuorumVotes(p.Evidences[1:], quorum) != nil
}

// check if the Proof of innocent of C1 is valid.
func validInnocenceProofOfC1(p *Proof, chain ChainContext) bool {
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
	height := preCommit.H()
	lastHeader := chain.GetHeaderByNumber(height - 1)
	if lastHeader == nil {
		return false
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())
	return message.OverQuorumVotes(p.Evidences, quorum) != nil
}

func hasEquivocatedVotes(votes []message.Msg) bool {
	/* //TODO(lorenzo) fix
	voteMap := make(map[common.Address]struct{})
	for _, vote := range votes {
		_, ok := voteMap[vote.Signer()]
		if !ok {
			voteMap[vote.Signer()] = struct{}{}
		} else {
			return true
		}
	}

	*/
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
func verifyProofSignatures(chain ChainContext, p *Proof) error {
	/*  //TODO(lorenzo) fix, might need to verify single message signatures here --> DoS?
	h := p.Message.H()
	lastHeader := chain.GetHeaderByNumber(h - 1)
	if lastHeader == nil {
		return errFutureMsg
	}
	if err := p.Message.Validate(lastHeader.CommitteeMember); err != nil {
		return errNotCommitteeMsg
	}
	// check if the number of evidence msgs are exceeded the max to prevent the abuse of the proof msg.
	if len(p.Evidences) > maxEvidenceMessages(lastHeader) {
		return errMaxEvidences
	}
	for _, msg := range p.Evidences {
		if msg.H() != h {
			return errBadHeight
		}
		if err := msg.Validate(lastHeader.CommitteeMember); err != nil {
			return errNotCommitteeMsg
		}
	}
	*/
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
	/* //TODO(lorenzo) fix
	offender := common.LeftPadBytes(m.Signer().Bytes(), 32)
	ruleID := common.LeftPadBytes([]byte{byte(rule)}, 32)
	block := make([]byte, 32)
	block = common.LeftPadBytes(new(big.Int).SetUint64(m.H()).Bytes(), 32)
	*/
	result := make([]byte, 160)
	//copy(result[0:32], successResult)
	//copy(result[32:64], offender)
	//copy(result[64:96], ruleID)
	//copy(result[96:128], block)
	//copy(result[128:160], m.Hash().Bytes())
	return result
}

// TODO(lorenzo) this might not be true anymore due to bls agg
// proofs that include the maximum amount of messages are the ones that require:
// 1. a proposal + a quorum of prevotes (worst-case number of evidence = MaxCommitteeSize + 1)
// 2. a proposal + a list of precommits (worst-case number of evidence = MaxRound + 1).
// thus the maximum number of evidence possible is the max between those two values
func maxEvidenceMessages(header *types.Header) int {
	committeeSize := len(header.Committee)
	if committeeSize > constants.MaxRound {
		return committeeSize + 1
	}
	return constants.MaxRound + 1
}
