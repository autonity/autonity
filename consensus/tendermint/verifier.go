package tendermint

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/hexutil"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	"github.com/clearmatics/autonity/consensus/tendermint/bft"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/log"
)

var (
	// errUnknownBlock is returned when the list of committee is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")
	// errUnauthorized is returned if a header is signed by a non authorized entity.
	errUnauthorized = errors.New("unauthorized")
	// errInvalidCoindbase is returned if the signer is not the coinbase address,
	errInvalidCoinbase = errors.New("invalid coinbase")
	// errInvalidDifficulty is returned if the difficulty of a block is not 1
	errInvalidDifficulty = errors.New("invalid difficulty")
	// errInvalidMixDigest is returned if a block's mix digest is not BFT digest.
	errInvalidMixDigest = errors.New("invalid BFT mix digest")
	// errInvalidNonce is returned if a block's nonce is invalid
	errInvalidNonce = errors.New("invalid nonce")
	// errInvalidUncleHash is returned if a block contains an non-empty uncle list.
	errInvalidUncleHash = errors.New("non empty uncle hash")
	// errInvalidTimestamp is returned if the timestamp of a block is lower than the previous block's timestamp + the minimum block period.
	errInvalidTimestamp = errors.New("invalid timestamp")
)

var (
	defaultDifficulty = big.NewInt(1)
	nilUncleHash      = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.
	emptyNonce        = types.BlockNonce{}

	nonceAuthVote = hexutil.MustDecode("0xffffffffffffffff") // Magic nonce number to vote on adding a new validator
	nonceDropVote = hexutil.MustDecode("0x0000000000000000") // Magic nonce number to vote on removing a validator.
)

type Verifier struct {
	logger      log.Logger
	vmConfig    *vm.Config
	finalizer   Finalizer
	blockPeriod uint64
}

func NewVerifier(c *vm.Config, finalizer Finalizer, blockPeriod uint64) *Verifier {
	return &Verifier{
		logger:    log.New(),
		vmConfig:  c,
		finalizer: finalizer,
	}
}

// VerifyProposal verifies the proposal. If a consensus.ErrFutureBlock error is returned,
// the time difference of the proposal and current time is also returned.
func (v *Verifier) VerifyProposal(proposal types.Block, blockchain *core.BlockChain, address string) (time.Duration, error) {
	// Check if the proposal is a valid block
	// TODO: fix always false statement and check for non nil
	// TODO: use interface instead of type
	block := &proposal
	//if block == nil {
	//	sb.logger.Error("Invalid proposal, %v", proposal)
	//	return 0, errInvalidProposal
	//}

	// check bad block
	if blockchain.HasBadBlock(block.Hash()) {
		return 0, core.ErrBlacklistedHash
	}

	// verify the header of proposed block
	err := v.VerifyHeader(blockchain, block.Header(), false)
	// ignore errEmptyCommittedSeals error because we don't have the committed seals yet
	if err == nil || err == types.ErrEmptyCommittedSeals {
		var (
			receipts types.Receipts

			usedGas        = new(uint64)
			gp             = new(core.GasPool).AddGas(block.GasLimit())
			header         = block.Header()
			proposalNumber = header.Number.Uint64()
			parent         = blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1)
		)

		// We need to process all of the transaction to get the latest state to get the latest committee
		state, stateErr := blockchain.StateAt(parent.Root())
		if stateErr != nil {
			return 0, stateErr
		}

		// Validate the body of the proposal
		if err = blockchain.Validator().ValidateBody(block); err != nil {
			return 0, err
		}

		// sb.blockchain.Processor().Process() was not called because it calls back Finalize() and would have modified the proposal
		// Instead only the transactions are applied to the copied state
		for i, tx := range block.Transactions() {
			state.Prepare(tx.Hash(), block.Hash(), i)
			// Might be vulnerable to DoS Attack depending on gaslimit
			// Todo : Double check
			receipt, receiptErr := core.ApplyTransaction(blockchain.Config(), blockchain, nil, gp, state, header, tx, usedGas, *v.vmConfig)
			if receiptErr != nil {
				return 0, receiptErr
			}
			receipts = append(receipts, receipt)
		}

		state.Prepare(common.ACHash(block.Number()), block.Hash(), len(block.Transactions()))
		committeeSet, receipt, err := v.finalizer.Finalize(blockchain, header, state, block.Transactions(), nil, receipts)
		if err != nil {
			return 0, err
		}
		receipts = append(receipts, receipt)
		//Validate the state of the proposal
		if err = blockchain.Validator().ValidateState(block, state, receipts, *usedGas); err != nil {
			return 0, err
		}

		//Perform the actual comparison
		if len(header.Committee) != len(committeeSet) {
			v.logger.Error("wrong committee set",
				"proposalNumber", proposalNumber,
				"extraLen", len(header.Committee),
				"currentLen", len(committeeSet),
				"committee", header.Committee,
				"current", committeeSet,
			)
			return 0, consensus.ErrInconsistentCommitteeSet
		}

		for i := range committeeSet {
			if header.Committee[i].Address != committeeSet[i].Address ||
				header.Committee[i].VotingPower.Cmp(committeeSet[i].VotingPower) != 0 {
				v.logger.Error("wrong committee member in the set",
					"index", i,
					"currentVerifier", address,
					"proposalNumber", proposalNumber,
					"headerCommittee", header.Committee[i],
					"computedCommittee", committeeSet[i],
					"fullHeader", header.Committee,
					"fullComputed", committeeSet,
				)
				return 0, consensus.ErrInconsistentCommitteeSet
			}
		}
		// At this stage committee field is consistent with the validator list returned by Soma-contract

		return 0, nil
	} else if err == consensus.ErrFutureBlock {
		return time.Until(time.Unix(int64(block.Header().Time), 0)), consensus.ErrFutureBlock
	}
	return 0, err
}

// VerifyHeader checks whether a header conforms to the consensus rules of a
// given engine. Verifying the seal may be done optionally here, or explicitly
// via the VerifySeal method.
func (v *Verifier) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, checkSeals bool) error {
	return v.verifyHeader(header, chain.GetHeaderByHash(header.ParentHash), checkSeals)
}

// verifyHeader checks whether a header conforms to the consensus rules. It
// expects the parent header to be provided unless header is the genesis
// header.
func (v *Verifier) verifyHeader(header, parent *types.Header, checkSeals bool) error {
	if header.Number == nil {
		return errUnknownBlock
	}
	// Don't waste time checking blocks from the future
	if big.NewInt(int64(header.Time)).Cmp(big.NewInt(time.Now().Unix())) > 0 {
		return consensus.ErrFutureBlock // This looks wrong, what if my clock has slipped slightly?
	}

	// Ensure that the coinbase is valid
	if header.Nonce != (emptyNonce) && !bytes.Equal(header.Nonce[:], nonceAuthVote) && !bytes.Equal(header.Nonce[:], nonceDropVote) {
		return errInvalidNonce
	}
	// Ensure that the mix digest is zero as we don't have fork protection currently
	if header.MixDigest != types.BFTDigest {
		return errInvalidMixDigest
	}
	// Ensure that the block doesn't contain any uncles which are meaningless in BFT
	if header.UncleHash != nilUncleHash {
		return errInvalidUncleHash
	}
	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	if header.Difficulty == nil || header.Difficulty.Cmp(defaultDifficulty) != 0 {
		return errInvalidDifficulty
	}

	// If this is the genesis block there is no further verification to be
	// done.
	if header.IsGenesis() {
		return nil
	}
	// We expect the parent to be non nil when header is not the genesis header.
	if parent == nil {
		return errUnknownBlock
	}
	return v.verifyHeaderAgainstParent(header, parent, checkSeals)
}

// verifyHeaderAgainstParent verifies that the given header is valid with respect to its parent.
func (v *Verifier) verifyHeaderAgainstParent(header, parent *types.Header, checkSeals bool) error {
	if parent.Number.Uint64() != header.Number.Uint64()-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}
	// Ensure that the block's timestamp isn't too close to it's parent
	if parent.Time+v.blockPeriod > header.Time {
		return errInvalidTimestamp

	}
	if !checkSeals {
		return nil
	}
	if err := v.verifySigner(header, parent); err != nil {
		return err
	}

	return v.verifyCommittedSeals(header, parent)
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications (the order is that of
// the input slice).
func (v *Verifier) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{}, 1)
	results := make(chan error, len(headers))
	go func() {
		for i, header := range headers {
			var parent *types.Header
			switch {
			case i > 0:
				parent = headers[i-1]
			case i == 0:
				parent = chain.GetHeaderByHash(header.ParentHash)
			}
			err := v.verifyHeader(header, parent, true)
			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

// VerifyUncles verifies that the given block's uncles conform to the consensus
// rules of a given engine.
func (v *Verifier) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	if len(block.Uncles()) > 0 {
		return errInvalidUncleHash
	}
	return nil
}

// verifySigner checks that the signer is part of the committee.
func (v *Verifier) verifySigner(header, parent *types.Header) error {
	// resolve the authorization key and check against signers
	signer, err := types.Ecrecover(header)
	if err != nil {
		return err
	}

	if header.Coinbase != signer {
		return errInvalidCoinbase
	}

	// Signer should be in the validator set of previous block's extraData.
	if parent.CommitteeMember(signer) != nil {
		return nil
	}

	return errUnauthorized
}

// verifyCommittedSeals validates that the committed seals for header come from
// committee members and that the voting power of the committed seals constitutes
// a quorum.
func (v *Verifier) verifyCommittedSeals(header, parent *types.Header) error {
	// The length of Committed seals should be larger than 0
	if len(header.CommittedSeals) == 0 {
		return types.ErrEmptyCommittedSeals
	}

	// Setup map to track votes made by committee members
	votes := make(map[common.Address]int, len(parent.Committee))

	// Calculate total voting power
	var committeeVotingPower uint64
	for _, member := range parent.Committee {
		committeeVotingPower += member.VotingPower.Uint64()
	}

	// Total Voting power for this block
	var power uint64
	// The data that was signed over for this block
	proposerSeal := header.ProposerSeal
	commitment, err := BuildCommitment(proposerSeal, header.Number.Uint64(), int64(header.Round), algorithm.ValueID(header.Hash()))
	if err != nil {
		return err
	}

	// 1. Get committed seals from current header
	for _, signedSeal := range header.CommittedSeals {
		// 2. Get the address from signature
		addr, err := types.GetSignatureAddressHash(commitment, signedSeal)
		if err != nil {
			v.logger.Error("not a valid address", "err", err)
			return types.ErrInvalidSignature
		}

		member := parent.CommitteeMember(addr)
		if member == nil {
			v.logger.Error(fmt.Sprintf("block %d had seal from non committee member %q", header.Number.Uint64(), addr.String()))
			return types.ErrInvalidCommittedSeals
		}

		votes[member.Address]++
		if votes[member.Address] > 1 {
			v.logger.Error(fmt.Sprintf("committee member %q had multiple seals on block %d", addr.String(), header.Number.Uint64()))
			return types.ErrInvalidCommittedSeals
		}
		power += member.VotingPower.Uint64()
	}

	// We need at least a quorum for the block to be considered valid
	if power < bft.Quorum(committeeVotingPower) {
		return types.ErrInvalidCommittedSeals
	}

	return nil
}

// VerifySeal checks whether the crypto seal on a header is valid according to
// the consensus rules of the given engine.
func (v *Verifier) VerifySeal(chain consensus.ChainHeaderReader, header *types.Header) error {
	// Ensure the signer is part of the committee

	// The genesis block is not signed.
	if header.IsGenesis() {
		return errUnknownBlock
	}

	// ensure that the difficulty equals to defaultDifficulty
	if header.Difficulty.Cmp(defaultDifficulty) != 0 {
		return errInvalidDifficulty
	}

	parent := chain.GetHeaderByHash(header.ParentHash)
	if parent == nil {
		// TODO make this ErrUnknownAncestor
		return errUnknownBlock
	}
	return v.verifySigner(header, parent)
}
