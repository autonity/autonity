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
	// ErrEmptyActivityProof is returned if the field of activity is empty.
	ErrEmptyActivityProof = errors.New("empty activity proof")
	// ErrInvalidActivityProofSignature is returned if the signature is not valid for the aggregated proof.
	ErrInvalidActivityProofSignature = errors.New("invalid activity proof signature")
	// ErrInsufficientActivityProof is returned if the voting power is less than quorum for activity proof.
	ErrInsufficientActivityProof = errors.New("insufficient power for activity proof")
)

// assembleActivityProof assembles the nodes' activity proof of height `h` with the aggregated precommit
// of height: `h-delta`. The proposer is incentivised to include as many signers as possible, however due to the
// timing of GST + Delta, assembling proof for the first delta blocks in an epoch is not required.
func (sb *Backend) assembleActivityProof(header *types.Header) types.AggregateSignature {
	var defaultProof types.AggregateSignature
	// for the 1st delta blocks, the proposer does not have to prove.
	if header.IsGenesis() {
		return defaultProof
	}

	// as this block haven't been finalized, thus we query the committee with its parent height from state db.
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

// TODO(lorenzo) update comment, I return the inactive ones now
// validateActivityProof validates the validity of the activity proof, and returns the proposer who provides
// an invalid activity proof as omission faulty node of the height, it also returns the signers of a valid
// activity proof which will be submitted to the omission accountability contract. Note: The proposer is innocence
// to provide no proof for the 1st delta blocks, thus, the 1st delta blocks of an epoch is not accountable.
func (sb *Backend) validateActivityProof(curHeader *types.Header) (bool, common.Address, *big.Int, []common.Address) {
	//TODO(lorenzo) double check return values when error or non standard case
	// for the 1st delta blocks, return nothing.
	if curHeader.IsGenesis() {
		return false, common.Address{}, new(big.Int), []common.Address{}
	}

	// todo: could be refined by on top of the epoch header PR.
	// since current block haven't been finalized at this phase,
	// thus we query the committee with its parent height from state db.
	lastHeight := new(big.Int).Sub(curHeader.Number, common.Big1) // TODO(Lorenzo) I don't think we need the minus one here
	lastEpochBlock, committee, err := sb.consensusViewOfHeight(lastHeight)
	if err != nil {
		panic(err)
	}

	if curHeader.Number.Uint64() <= lastEpochBlock.Uint64()+tendermint.DeltaBlocks {
		sb.logger.Debug("Skip to validate activity proof for the 1st delta blocks of epoch",
			"height", curHeader.Number.Uint64(), "lastEpochBlock", lastEpochBlock)
		return false, common.Address{}, new(big.Int), []common.Address{}
	}

	// at block finalization phase, the coinbase was checked already, thus take it as the proposer of the block.
	proposer := curHeader.Coinbase

	// after the 1st delta blocks, the proposer is accountable for not assembling valid activity proof.
	signers, proposerEffort, err := sb.verifyActivityProof(curHeader, committee)
	if err != nil {
		sb.logger.Info("Faulty activity proof addressed, proposer is omission faulty", "proposer", proposer)
		return true, proposer, proposerEffort, []common.Address{}
	}

	// we have got the signers, let's compute the absentees and return them
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

	return false, proposer, proposerEffort, absentees
}

// verifyActivityProof validates that the activity proof for header come from committee members and that the voting
// power constitute a quorum, it returns the node IDs to be submitted to omission accountability contract. Any error
// in this function will cause the proposer to be faulty for omission accountability.
func (sb *Backend) verifyActivityProof(header *types.Header, committee types.Committee) ([]common.Address, *big.Int, error) {
	// todo: replace by committee.Member() once epoch header PR is merged.
	// todo: replace by committee.TotalVotingPower() once epoch header PR is merged.
	// todo: if we submit the address of nodes to omission accountability contract, since ID could changes on committee
	//  reshuffling.

	// TODO(lorenzo) double check return values in case of error

	// un-finalized proposals will have these fields set to nil
	if header.ActivityProof.Signature == nil || header.ActivityProof.Signers == nil {
		return []common.Address{}, new(big.Int), ErrEmptyActivityProof
	}

	activityProof := header.ActivityProof.Copy() // copy so that we do not modify the header when doing Signers.Validate()
	if err := activityProof.Signers.Validate(len(committee)); err != nil {
		return []common.Address{}, new(big.Int), fmt.Errorf("Invalid activity proof signers information: %w", err)
	}

	targetHeight := header.Number.Uint64() - tendermint.DeltaBlocks
	targetHeader := sb.BlockChain().GetHeaderByNumber(targetHeight)
	// TODO(lorenzo) what if target header nil, is that possible?

	// The data that was signed over for this block
	// TODO(lorenzo) problem here, the commit round might not be the same across nodes
	headerSeal := message.PrepareCommittedSeal(targetHeader.Hash(), int64(targetHeader.Round), targetHeader.Number)

	// Total assembled voting power for the activity proof
	power := new(big.Int)
	signers := make([]common.Address, activityProof.Signers.Len())
	for i, index := range activityProof.Signers.FlattenUniq() {
		power.Add(power, committee[index].VotingPower)
		signers[i] = committee[index].Address
	}

	// verify signature
	var keys []blst.PublicKey //nolint
	for _, index := range activityProof.Signers.Flatten() {
		keys = append(keys, committee[index].ConsensusKey)
	}
	aggregatedKey, err := blst.AggregatePublicKeys(keys)
	if err != nil {
		sb.logger.Crit("Failed to aggregate keys from committee members", "err", err)
	}
	valid := activityProof.Signature.Verify(aggregatedKey, headerSeal[:])
	if !valid {
		sb.logger.Error("block had invalid activity proof signature")
		return []common.Address{}, new(big.Int), ErrInvalidActivityProofSignature
	}

	// We need at least a quorum for the activity proof.
	quorum := bft.Quorum(committee.TotalVotingPower())
	if power.Cmp(quorum) < 0 {
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