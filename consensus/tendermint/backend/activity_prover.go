package backend

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
)

var (
	// ErrNotEmptyActivityProof is returned when the activity proof should have been empty, but it is not
	ErrNotEmptyActivityProof = errors.New("not empty activity proof")
	// ErrEmptyActivityProof is returned if the activity proof field is empty.
	ErrEmptyActivityProof = errors.New("empty activity proof")
	// ErrInvalidActivityProofSignature is returned if the signature is not valid for the aggregated proof.
	ErrInvalidActivityProofSignature = errors.New("invalid activity proof signature")
	// ErrInsufficientActivityProof is returned if the voting power is less than quorum for activity proof.
	ErrInsufficientActivityProof = errors.New("insufficient power for activity proof")
	// ErrImcompleteActivityProof is returned when the proof has a single field nil and the other is not nil
	ErrIncompleteActivityProof = errors.New("incomplete activity proof")
)

// assembleActivityProof assembles the nodes' activity proof of height `h` with the aggregated precommit
// of height: `h-delta`. The proposer is incentivised to include as many signers as possible.
// If the proposer does not have to OR cannot provide a valid activity proof, it should leave the proof empty (internal pointers set to nil)
func (sb *Backend) assembleActivityProof(h uint64) (types.AggregateSignature, uint64, error) {
	number := new(big.Int).SetUint64(h)

	// TODO(lorenzo) re-review this part about lastEpochBlock after epoch-header is merged
	lastEpochBlock, committee, err := sb.consensusViewOfHeight(number)
	if err != nil {
		return types.AggregateSignature{}, 0, fmt.Errorf("Error while fetching lastEpochBlock for height %d: %w", number.Uint64(), err)
	}

	// for the 1st delta blocks of the epoch, the proposer does not have to provide an activity proof
	if h <= lastEpochBlock.Uint64()+tendermint.DeltaBlocks {
		sb.logger.Debug("Skip to assemble activity proof at the start of epoch", "height", h, "lastEpochBlock", lastEpochBlock)
		return types.AggregateSignature{}, 0, nil
	}

	// after delta blocks, get quorum certificates from height h-delta.
	targetHeight := h - tendermint.DeltaBlocks
	targetHeader := sb.BlockChain().GetHeaderByNumber(targetHeight)
	targetRound := targetHeader.Round

	precommits := sb.MsgStore.GetPrecommits(targetHeight, func(m *message.Precommit) bool {
		return m.R() == int64(targetRound) && m.Value() == targetHeader.Hash()
	})

	// we should have provided an activity proof, but we do not have past messages
	if len(precommits) == 0 {
		sb.logger.Warn("Failed to provide activity valid activity proof as proposer", "height", h, "targetHeight", targetHeight)
		return types.AggregateSignature{}, 0, nil
	}

	votes := make([]message.Vote, len(precommits))
	for i, p := range precommits {
		votes[i] = p
	}

	aggregatePrecommit := message.AggregatePrecommits(votes)

	// if we do not have enough voting power, leave the proof empty
	quorum := bft.Quorum(committee.TotalVotingPower())
	if aggregatePrecommit.Power().Cmp(quorum) < 0 {
		sb.logger.Warn("Failed to provide activity valid activity proof as proposer, not enough voting power", "height", h, "targetHeight", targetHeight, "power", aggregatePrecommit.Power(), "quorum", quorum)
		return types.AggregateSignature{}, 0, nil
	}

	return types.NewAggregateSignature(aggregatePrecommit.Signature().(*blst.BlsSignature), aggregatePrecommit.Signers()), targetRound, nil
}

// validateActivityProof validates the activity proof, and returns:
// 1. a boolean indicating whether the proposer did provide a valid proof or not
// 2. the proposer "effort" (amount of voting power exceeding quorum) included into the proof
// 3. the list of absent validators derived from the proof
// 4. an error that indicates whether we should reject that proposal
// Note: The proposer does not have to provide a proof for the 1st delta blocks of an epoch.
func (sb *Backend) validateActivityProof(proof types.AggregateSignature, h uint64, r uint64) (bool, *big.Int, []common.Address, error) {
	number := new(big.Int).SetUint64(h)

	// TODO(lorenzo) re-review this part about lastEpochBlock after epoch-header is merged
	lastEpochBlock, committee, err := sb.consensusViewOfHeight(number)
	if err != nil {
		return false, new(big.Int), []common.Address{}, fmt.Errorf("Error while fetching consensus view for height %d: %w", number.Uint64(), err)
	}

	// during the first delta blocks of the epoch, the proof should be empty. If not, reject proposal
	if h <= lastEpochBlock.Uint64()+tendermint.DeltaBlocks {
		sb.logger.Debug("Validating activity proof in first delta blocks, should be empty", "height", h, "lastEpochBlock", lastEpochBlock)
		if !proof.Empty() {
			return false, new(big.Int), []common.Address{}, ErrNotEmptyActivityProof
		}
		return false, new(big.Int), []common.Address{}, nil
	}

	// at this point the proof should not be empty and should contain at least quorum voting power, otherwise the proposer is faulty
	if proof.Empty() {
		return true, new(big.Int), []common.Address{}, nil
	}

	// if the proof is malformed, reject the proposal
	if proof.Incomplete() {
		return false, new(big.Int), []common.Address{}, ErrIncompleteActivityProof
	}

	activityProof := proof.Copy() // copy so that we do not modify the header when doing Signers.Validate()
	// if the activity proof is malformed however, reject the proposal. We cannot accept arbitrary data in the proposal
	if err := activityProof.Signers.Validate(len(committee)); err != nil {
		return false, new(big.Int), []common.Address{}, fmt.Errorf("Invalid activity proof signers information: %w", err)
	}

	signers, proposerEffort, err := sb.verifyActivityProof(activityProof, committee, h-tendermint.DeltaBlocks, r)
	if err != nil {
		// verifyActivityProof returns an error if the signature is invalid or if the voting power is not enough
		// in both cases, we reject the proposal.
		sb.logger.Info("Rejecting proposal with invalid activity proof signature", "err", err)
		return false, new(big.Int), []common.Address{}, err
	}

	// we have got the signers, let's compute the absentees and return them
	absentees := computeAbsents(signers, committee)

	return false, proposerEffort, absentees, nil
}

func computeAbsents(signers []common.Address, committee types.Committee) []common.Address {
	// TODO(lorenzo) optimization, maybe using a map?
	var absentees []common.Address
	for _, member := range committee {
		found := false
		for _, signer := range signers {
			if signer == member.Address {
				found = true
			}
		}
		if !found {
			absentees = append(absentees, member.Address)
		}
	}
	return absentees
}

// verifyActivityProof validates that the activity proof  is signed  by committee members and that the voting
// power is >= quorum. It returns the node addresses of the signers and the voting power that exceeds quorum.
// Any error in this function will cause the proposer to be faulty for omission accountability.
func (sb *Backend) verifyActivityProof(proof types.AggregateSignature, committee types.Committee, targetHeight uint64, round uint64) ([]common.Address, *big.Int, error) {
	// The data that was signed over for this block
	targetHeader := sb.BlockChain().GetHeaderByNumber(targetHeight)
	headerSeal := message.PrepareCommittedSeal(targetHeader.Hash(), int64(round), targetHeader.Number)

	// Total assembled voting power for the activity proof
	power := new(big.Int)
	signers := make([]common.Address, proof.Signers.Len())
	for i, index := range proof.Signers.FlattenUniq() {
		power.Add(power, committee[index].VotingPower)
		signers[i] = committee[index].Address
	}

	// verify signature
	var keys []blst.PublicKey //nolint
	for _, index := range proof.Signers.Flatten() {
		keys = append(keys, committee[index].ConsensusKey)
	}
	aggregatedKey, err := blst.AggregatePublicKeys(keys)
	if err != nil {
		sb.logger.Crit("Failed to aggregate keys from committee members", "err", err)
	}
	valid := proof.Signature.Verify(aggregatedKey, headerSeal[:])
	if !valid {
		sb.logger.Info("block had invalid activity proof signature")
		return []common.Address{}, new(big.Int), ErrInvalidActivityProofSignature
	}

	// We need at least a quorum for the activity proof.
	quorum := bft.Quorum(committee.TotalVotingPower())
	if power.Cmp(quorum) < 0 {
		sb.logger.Info("block had insufficient voting power in activity proof")
		return []common.Address{}, new(big.Int), ErrInsufficientActivityProof
	}

	proposerEffort := new(big.Int).Set(power)
	proposerEffort.Sub(proposerEffort, quorum)

	return signers, proposerEffort, nil
}

// consensusViewOfHeight returns the last epoch block and the corresponding committee of a specific height, it removes
// the dependence of blockchain.
func (sb *Backend) consensusViewOfHeight(height *big.Int) (*big.Int, types.Committee, error) {
	header := sb.BlockChain().CurrentBlock().Header()
	state, err := sb.BlockChain().State()
	if err != nil {
		return nil, nil, err
	}
	return sb.BlockChain().ProtocolContracts().AutonityContract.GetConsensusViewOfHeight(header, state, height)
}
