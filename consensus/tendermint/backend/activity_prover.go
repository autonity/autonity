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

// assembleActivityProof assembles the nodes' activity proof of height: h with the aggregated precommit
// of height: h-dela. Proposer is incentivised to assemble proof as much as possible, however due to the
// timing of GST + Delta, assembling proof for the first delta blocks in an epoch is not required.
func (sb *Backend) assembleActivityProof(header *types.Header) types.AggregateSignature {
	var defaultProof types.AggregateSignature
	// for the 1st delta blocks, the proposer does not have to prove.
	if header.IsGenesis() {
		return defaultProof
	}

	lastHeight := new(big.Int).Sub(header.Number, common.Big1)
	lastEpochBlock, _, err := sb.consensusViewOfHeight(lastHeight)
	if err != nil {
		panic(err)
	}
	if header.Number.Uint64() <= lastEpochBlock.Uint64()+tendermint.DeltaBlocks {
		sb.logger.Debug("Skip to assemble activity proof at the starting of epoch",
			"height", header.Number.Uint64(), "lastEpochBlock", lastEpochBlock)
		return defaultProof
	}

	// after delta blocks, get quorum certificates from height h-delta.
	targetHeight := header.Number.Uint64() - tendermint.DeltaBlocks
	targetHeader := sb.BlockChain().GetHeaderByNumber(targetHeight)

	// get precommits for the same value of the height h-delta, aggregate the missing ones of the
	precommits := sb.MsgStore.GetPrecommits(targetHeight, func(m *message.Precommit) bool {
		return m.R() == int64(targetHeader.Round) && m.Value() == targetHeader.Hash()
	})

	votes := make([]message.Vote, len(precommits))
	for i, p := range precommits {
		votes[i] = p
	}

	if len(votes) == 0 {
		return defaultProof
	}

	aggregate := message.AggregatePrecommits(votes)
	defaultProof.Signature = aggregate.Signature().(*blst.BlsSignature)
	defaultProof.Signers = aggregate.Signers()
	return defaultProof
}

// validateActivityProof validates the validity of the activity proof, and returns the proposer who provides
// an invalid activity proof as omission faulty node of the height, it also returns the signers of a valid
// activity proof which will be submitted to the omission accountability contract. Note: The proposer is innocence
// to provide no proof for the 1st delta blocks, thus, the 1st delta blocks of an epoch is not accountable.
func (sb *Backend) validateActivityProof(curHeader *types.Header) (bool, []*big.Int) {
	// for the 1st delta blocks, return nothing.
	if curHeader.IsGenesis() {
		return false, []*big.Int{}
	}

	// todo: could be refined by on top of the epoch header PR.
	lastHeight := new(big.Int).Sub(curHeader.Number, common.Big1)
	lastEpochBlock, committee, err := sb.consensusViewOfHeight(lastHeight)
	if err != nil {
		panic(err)
	}

	if curHeader.Number.Uint64() <= lastEpochBlock.Uint64()+tendermint.DeltaBlocks {
		sb.logger.Debug("Skip to validate activity proof for the 1st delta blocks of epoch",
			"height", curHeader.Number.Uint64(), "lastEpochBlock", lastEpochBlock)
		return false, []*big.Int{}
	}

	// at block finalization phase, the coinbase was checked already, thus take it as the proposer of the block.
	proposer := curHeader.Coinbase

	// after the 1st delta blocks, the proposer is accountable for not assembling valid activity proof.
	signers, err := sb.verifyActivityProof(curHeader, committee, proposer)
	if err != nil {
		sb.logger.Info("Faulty activity proof addressed, proposer is omission faulty", "proposer", proposer)
		return true, signers
	}

	return false, signers
}

// verifyActivityProof validates that the activity proof for header come from committee members and that the voting
// power constitute a quorum, it returns the node IDs to be submitted to omission accountability contract. Any error
// in this function will cause the proposer to be faulty for omission accountability.
func (sb *Backend) verifyActivityProof(header *types.Header, committee types.Committee, proposer common.Address) ([]*big.Int, error) {

	// todo: replace by committee.Member() once epoch header PR is merged.
	// todo: replace by committee.TotalVotingPower() once epoch header PR is merged.
	// todo: if we submit the address of nodes to omission accountability contract, since ID could changes on committee
	//  reshuffling.
	// Calculate total voting power of committee
	committeeVotingPower := new(big.Int)
	proposerID := new(big.Int)
	for _, member := range committee {
		if member.Address == proposer {
			proposerID.SetUint64(member.Index)
		}
		committeeVotingPower.Add(committeeVotingPower, member.VotingPower)
	}

	// un-finalized proposals will have these fields set to nil
	if header.ActivityProof.Signature == nil || header.ActivityProof.Signers == nil {
		return []*big.Int{proposerID}, ErrEmptyActivityProof
	}

	activityProof := header.ActivityProof.Copy() // copy so that we do not modify the header when doing Signers.Validate()
	if err := activityProof.Signers.Validate(len(committee)); err != nil {
		return []*big.Int{proposerID}, fmt.Errorf("Invalid activity proof signers information: %w", err)
	}

	targetHeight := header.Number.Uint64() - tendermint.DeltaBlocks
	targetHeader := sb.BlockChain().GetHeaderByNumber(targetHeight)

	// The data that was signed over for this block
	headerSeal := message.PrepareCommittedSeal(targetHeader.Hash(), int64(targetHeader.Round), targetHeader.Number)

	// Total assembled voting power for the activity proof
	power := new(big.Int)
	signers := activityProof.Signers.FlattenUniq()
	IDs := make([]*big.Int, len(signers))
	for i, index := range signers {
		IDs[i] = new(big.Int).SetInt64(int64(index))
		power.Add(power, committee[index].VotingPower)
	}

	// verify signature
	var keys [][]byte //nolint
	for _, index := range activityProof.Signers.Flatten() {
		keys = append(keys, committee[index].ConsensusKeyBytes)
	}
	aggregatedKey, err := blst.AggregatePublicKeys(keys)
	if err != nil {
		sb.logger.Crit("Failed to aggregate keys from committee members", "err", err)
	}
	valid := activityProof.Signature.Verify(aggregatedKey, headerSeal[:])
	if !valid {
		sb.logger.Error("block had invalid activity proof signature")
		return []*big.Int{proposerID}, ErrInvalidActivityProofSignature
	}

	// We need at least a quorum for the activity proof.
	if power.Cmp(bft.Quorum(committeeVotingPower)) < 0 {
		return []*big.Int{proposerID}, ErrInsufficientActivityProof
	}

	return IDs, nil
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
