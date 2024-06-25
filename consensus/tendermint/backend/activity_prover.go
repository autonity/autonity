package backend

import (
	"errors"
	"fmt"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
	"math/big"
)

var (
	// ErrEmptyActivityProof is returned if the field of activity is empty.
	ErrEmptyActivityProof = errors.New("empty activity proof")
	// ErrInvalidActivityProofSignature is returned if the signature is not valid for the aggregated proof.
	ErrInvalidActivityProofSignature = errors.New("invalid activity proof signature")
	// ErrInsufficientActivityProof is returned if the voting power is less than quorum for activity proof.
	ErrInsufficientActivityProof = errors.New("insufficient power for activity proof")
)

type ActivityReport struct {
	FaultyProposer common.Address
	Signers        []int
}

// assembleActivityProof assembles the nodes' activity proof of height: h with the aggregated precommit
// of height: h-dela. Proposer is incentivised to assemble proof as much as possible, however due to the
// timing of GST + Delta, assembling proof for the first delta blocks in an epoch is not required.
func (sb *Backend) assembleActivityProof(h uint64) (types.AggregateSignature, error) {
	var defaultProof types.AggregateSignature
	// for the 1st delta blocks, the proposer does not have to prove.
	if h == 0 {
		return defaultProof, nil
	}
	lastEpochBlock, err := sb.lastEpochBlockOfHeight(h)
	if err != nil {
		panic(err)
	}
	if h <= lastEpochBlock+tendermint.DeltaBlocks {
		sb.logger.Debug("Skip to assemble activity proof at the starting of epoch",
			"height", h, "lastEpochBlock", lastEpochBlock)
		return defaultProof, nil
	}

	// after delta blocks, get quorum certificates from height h-delta.
	targetHeight := h - tendermint.DeltaBlocks
	header := sb.BlockChain().GetHeaderByNumber(targetHeight)

	// get precommits for the same value of the height h-delta, aggregate the missing ones of the
	precommits := sb.MsgStore.GetPrecommits(targetHeight, func(m *message.Precommit) bool {
		return m.R() == int64(header.Round) && m.Value() == header.Hash()
	})

	votes := make([]message.Vote, len(precommits))
	for i, p := range precommits {
		votes[i] = p
	}
	aggregate := message.AggregatePrecommits(votes)
	defaultProof.Signature = aggregate.Signature().(*blst.BlsSignature)
	defaultProof.Signers = aggregate.Signers()
	return defaultProof, nil
}

// validateActivityProof validates the validity of the activity proof, and returns the proposer who provides
// an invalid activity proof as omission faulty node of the height, it also returns the signers of a valid
// activity proof which will be submitted to the omission accountability contract. Note: The proposer is innocence
// to provide no proof for the 1st delta blocks, thus, the 1st delta blocks of an epoch is not accountable.
func (sb *Backend) validateActivityProof(curHeader, parent *types.Header) *ActivityReport {
	report := &ActivityReport{}
	// for the 1st delta blocks, return nothing.
	if curHeader.IsGenesis() {
		return report
	}

	h := curHeader.Number.Uint64()
	lastEpochBlock, err := sb.lastEpochBlockOfHeight(h)
	if err != nil {
		panic(err)
	}

	if h <= lastEpochBlock+tendermint.DeltaBlocks {
		sb.logger.Debug("Skip to validate activity proof for the 1st delta blocks of epoch",
			"height", h, "lastEpochBlock", lastEpochBlock)
		return report
	}

	proposer, err := types.ECRecover(curHeader)
	if err != nil {
		panic(err) // as at this phase, the block proposer/signer was verified.
	}

	// after the 1st delta blocks, the proposer is accountable for not assembling valid activity proof.
	if err = sb.verifyActivityProof(curHeader, parent); err != nil {
		sb.logger.Info("Faulty activity proof addressed", "proposer", proposer)
		report.FaultyProposer = proposer
		return report
	}

	// todo: (Jason) double check if the flattenUniq() is deterministic?
	report.Signers = curHeader.ActivityProof.Signers.FlattenUniq()
	return report
}

// verifyActivityProof validates that the activity proof for header come from committee members and
// that the voting power constitute a quorum.
func (sb *Backend) verifyActivityProof(header, parent *types.Header) error {
	// un-finalized proposals will have these fields set to nil
	if header.ActivityProof.Signature == nil || header.ActivityProof.Signers == nil {
		return ErrEmptyActivityProof
	}
	activityProof := header.ActivityProof.Copy() // copy so that we do not modify the header when doing Signers.Validate()
	if err := activityProof.Signers.Validate(len(parent.Committee)); err != nil {
		return fmt.Errorf("Invalid activity proof signers information: %w", err)
	}

	// todo: replace by committee.TotalVotingPower()
	// Calculate total voting power of committee
	committeeVotingPower := new(big.Int)
	for _, member := range parent.Committee {
		committeeVotingPower.Add(committeeVotingPower, member.VotingPower)
	}

	// The data that was signed over for this block
	headerSeal := message.PrepareCommittedSeal(header.Hash(), int64(header.Round), header.Number)

	// Total assembled voting power for the activity proof
	power := new(big.Int)
	for _, index := range activityProof.Signers.FlattenUniq() {
		power.Add(power, parent.Committee[index].VotingPower)
	}

	// verify signature
	var keys [][]byte //nolint
	for _, index := range activityProof.Signers.Flatten() {
		keys = append(keys, parent.Committee[index].ConsensusKeyBytes)
	}
	aggregatedKey, err := blst.AggregatePublicKeys(keys)
	if err != nil {
		sb.logger.Crit("Failed to aggregate keys from committee members", "err", err)
	}
	valid := activityProof.Signature.Verify(aggregatedKey, headerSeal[:])
	if !valid {
		sb.logger.Error("block had invalid committed seal")
		return ErrInvalidActivityProofSignature
	}

	// We need at least a quorum for the activity proof.
	if power.Cmp(bft.Quorum(committeeVotingPower)) < 0 {
		return ErrInsufficientActivityProof
	}

	return nil
}

// lastEpochBlockOfHeight get the last epoch block of height from AC contract.
// todo(Jason) replace this by query from the header chain when epoch-header feature is merged.
func (sb *Backend) lastEpochBlockOfHeight(height uint64) (uint64, error) {
	return uint64(0), nil
}
