// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package backend

import (
	"bytes"
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/clearmatics/autonity/consensus/tendermint/committee"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/hexutil"
	"github.com/clearmatics/autonity/consensus"
	tendermintCore "github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/rpc"
)

const (
	inmemorySnapshots = 128 // Number of recent vote snapshots to keep in memory
	inmemoryPeers     = 40
	inmemoryMessages  = 1024
)

// ErrStartedEngine is returned if the engine is already started
var ErrStartedEngine = errors.New("started engine")

var (
	// errInvalidProposal is returned when a prposal is malformed.
	//errInvalidProposal = errors.New("invalid proposal")
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
	// errInvalidRound is returned if the round exceed maximum round number.
	errInvalidRound = errors.New("invalid round")
)
var (
	defaultDifficulty = big.NewInt(1)
	nilUncleHash      = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.
	emptyNonce        = types.BlockNonce{}
	now               = time.Now

	nonceAuthVote = hexutil.MustDecode("0xffffffffffffffff") // Magic nonce number to vote on adding a new validator
	nonceDropVote = hexutil.MustDecode("0x0000000000000000") // Magic nonce number to vote on removing a validator.
)

// Author retrieves the Ethereum address of the account that minted the given
// block, which may be different from the header's coinbase if a consensus
// engine is based on signatures.
func (sb *Backend) Author(header *types.Header) (common.Address, error) {
	return types.Ecrecover(header)
}

// VerifyHeader checks whether a header conforms to the consensus rules of a
// given engine. Verifying the seal may be done optionally here, or explicitly
// via the VerifySeal method.
func (sb *Backend) VerifyHeader(chain consensus.ChainReader, header *types.Header, _ bool) error {
	return sb.verifyHeader(chain, header, nil)
}

// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (sb *Backend) verifyHeader(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	if header.Number == nil {
		return errUnknownBlock
	}
	if header.Round > tendermintCore.MaxRound {
		return errInvalidRound
	}
	// Don't waste time checking blocks from the future
	if big.NewInt(int64(header.Time)).Cmp(big.NewInt(now().Unix())) > 0 {
		return consensus.ErrFutureBlock
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

	return sb.verifyCascadingFields(chain, header, parents)
}

// verifyCascadingFields verifies all the header fields that are not standalone,
// rather depend on a batch of previous headers. The caller may optionally pass
// in a batch of parents (ascending order) to avoid looking those up from the
// database. This is useful for concurrently verifying a batch of new headers.
func (sb *Backend) verifyCascadingFields(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// The genesis block is the always valid dead-end
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}
	// Ensure that the block's timestamp isn't too close to it's parent
	var parent *types.Header
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		parent = chain.GetHeader(header.ParentHash, number-1)
	}
	if parent == nil || parent.Number.Uint64() != number-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}
	if parent.Time+sb.config.BlockPeriod > header.Time {
		return errInvalidTimestamp
	}

	if err := sb.verifySigner(chain, header, parents); err != nil {
		return err
	}

	return sb.verifyCommittedSeals(chain, header, parents)
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications (the order is that of
// the input slice).
func (sb *Backend) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{}, 1)
	results := make(chan error, len(headers))
	go func() {
		for i, header := range headers {
			err := sb.verifyHeader(chain, header, headers[:i])

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
func (sb *Backend) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	if len(block.Uncles()) > 0 {
		return errInvalidUncleHash
	}
	return nil
}

// verifySigner checks whether the signer is in parent's validator set
func (sb *Backend) verifySigner(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}

	committee, err := sb.committee(header, parents, chain)

	if err != nil {
		return err
	}

	// resolve the authorization key and check against signers
	signer, err := types.Ecrecover(header)
	if err != nil {
		return err
	}

	if header.Coinbase != signer {
		return errInvalidCoinbase
	}

	// Signer should be in the validator set of previous block's extraData.
	for _, committeeMember := range committee.Committee() {
		if committeeMember.Address == signer {
			return nil
		}
	}

	return errUnauthorized
}

// verifyCommittedSeals checks whether every committed seal is signed by one of the parent's committee
func (sb *Backend) verifyCommittedSeals(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	number := header.Number.Uint64()
	// We don't need to verify committed seals in the genesis block
	if number == 0 {
		return nil
	}

	committeeSet, err := sb.committee(header, parents, chain)
	if err != nil {
		return err
	}
	sealed := make(map[common.Address]struct{})

	// The length of Committed seals should be larger than 0
	if len(header.CommittedSeals) == 0 {
		return types.ErrEmptyCommittedSeals
	}

	// Check whether the committed seals are generated by parent's committee
	var power uint64 // Total Voting power for this block
	proposalSeal := tendermintCore.PrepareCommittedSeal(header.Hash(), int64(header.Round), header.Number)
	// 1. Get committed seals from current header
	for _, seal := range header.CommittedSeals {
		// 2. Get the original address by seal and parent block hash
		addr, err := types.GetSignatureAddress(proposalSeal, seal)
		if err != nil {
			sb.logger.Error("not a valid address", "err", err)
			return types.ErrInvalidSignature
		}
		// Every validator can have only one seal. If more than one seals are signed by a
		// validator, the validator cannot be found and errInvalidCommittedSeals is returned.
		if _, ok := sealed[addr]; !ok {
			sealed[addr] = struct{}{}
			_, member, err := committeeSet.GetByAddress(addr)
			if err != nil {
				sb.logger.Error("could not retrieve member")
				return types.ErrInvalidCommittedSeals
			}
			power += member.VotingPower.Uint64()
		} else {
			return types.ErrInvalidCommittedSeals
		}
	}

	// The amount of power should be larger than a Quorum
	if power < committeeSet.Quorum() {
		return types.ErrInvalidCommittedSeals
	}

	return nil
}

// VerifySeal checks whether the crypto seal on a header is valid according to
// the consensus rules of the given engine.
func (sb *Backend) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	// get parent header and ensure the signer is in parent's validator set
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}

	// ensure that the difficulty equals to defaultDifficulty
	if header.Difficulty.Cmp(defaultDifficulty) != 0 {
		return errInvalidDifficulty
	}
	return sb.verifySigner(chain, header, nil)
}

// Prepare initializes the consensus fields of a block header according to the
// rules of a particular engine. The changes are executed inline.
func (sb *Backend) Prepare(chain consensus.ChainReader, header *types.Header) error {
	// unused fields, force to set to empty
	header.Coinbase = sb.Address()
	header.Nonce = emptyNonce
	header.MixDigest = types.BFTDigest

	// copy the parent extra data as the header extra data
	number := header.Number.Uint64()
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	// use the same difficulty for all blocks
	header.Difficulty = defaultDifficulty

	// set header's timestamp
	header.Time = new(big.Int).Add(big.NewInt(int64(parent.Time)), new(big.Int).SetUint64(sb.config.BlockPeriod)).Uint64()
	if int64(header.Time) < time.Now().Unix() {
		header.Time = uint64(time.Now().Unix())
	}
	return nil
}

// Finalize runs any post-transaction state modifications (e.g. block rewards)
// Finaize doesn't modify the passed header.
func (sb *Backend) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (types.Committee, *types.Receipt, error) {
	sb.blockchainInitMu.Lock()
	if sb.blockchain == nil {
		sb.blockchain = chain.(*core.BlockChain) // in the case of Finalize() called before the engine start()
	}
	sb.blockchainInitMu.Unlock()

	committeeSet, receipt, err := sb.AutonityContractFinalize(header, chain, state, txs, receipts)
	if err != nil {
		sb.logger.Error("AutonityContractFinalize", "err", err.Error())
		return nil, nil, err
	}

	return committeeSet, receipt, nil
}

// FinalizeAndAssemble call Finaize to compute post transacation state modifications
// and assembles the final block.
func (sb *Backend) FinalizeAndAssemble(chain consensus.ChainReader, header *types.Header, statedb *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts *[]*types.Receipt) (*types.Block, error) {

	statedb.Prepare(common.ACHash(header.Number), common.Hash{}, len(txs))
	committeeSet, receipt, err := sb.Finalize(chain, header, statedb, txs, uncles, *receipts)
	if err != nil {
		return nil, err
	}
	*receipts = append(*receipts, receipt)
	// No block rewards in BFT, so the state remains as is and uncles are dropped
	header.Root = statedb.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = nilUncleHash

	// add committee to extraData's committee section
	header.Committee = committeeSet
	return types.NewBlock(header, txs, nil, *receipts), nil
}

// AutonityContractFinalize is called to deploy the Autonity Contract at block #1. it returns as well the
// committee field containaining the list of committee members allowed to participate in consensus for the next block.
func (sb *Backend) AutonityContractFinalize(header *types.Header, chain consensus.ChainReader, state *state.StateDB,
	txs []*types.Transaction, receipts []*types.Receipt) (types.Committee, *types.Receipt, error) {
	sb.contractsMu.Lock()
	defer sb.contractsMu.Unlock()

	if header.Number.Int64() == 1 {
		err := sb.blockchain.GetAutonityContract().DeployAutonityContract(chain, header, state)
		if err != nil {
			sb.logger.Error("Deploy autonity contract error", "error", err)
			return nil, nil, err
		}
	}

	committeeSet, receipt, err := sb.blockchain.GetAutonityContract().FinalizeAndGetCommittee(txs, receipts, header, state)
	if err != nil {
		sb.logger.Error("Autonity Contract finalize returns err", "err", err)
		return nil, nil, err
	}
	return committeeSet, receipt, nil
}

// Seal generates a new block for the given input block with the local miner's
// seal place on top.
func (sb *Backend) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	sb.coreMu.RLock()
	isStarted := sb.coreStarted
	sb.coreMu.RUnlock()
	if !isStarted {
		return ErrStoppedEngine
	}

	// update the block header and signature and propose the block to core engine
	header := block.Header()
	number := header.Number.Uint64()

	// Bail out if we're unauthorized to sign a block
	if committeeSet, err := sb.Committee(number); err == nil {
		if _, _, errP := committeeSet.GetByAddress(sb.Address()); errP != nil {
			sb.logger.Error("error validator errUnauthorized", "addr", sb.address)
			return errUnauthorized
		}
	} else {
		sb.logger.Error("couldn't retrieve current block committee", "addr", sb.address)
		return errUnauthorized
	}

	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		sb.logger.Error("Error ancestor")
		return consensus.ErrUnknownAncestor
	}
	block, err := sb.AddSeal(block)
	if err != nil {
		sb.logger.Error("seal error updateBlock", "err", err.Error())
		return err
	}

	// wait for the timestamp of header, use this to adjust the block period
	delay := time.Unix(int64(block.Header().Time), 0).Sub(now())
	select {
	case <-time.After(delay):
		// nothing to do
	case <-sb.stopped:
		return nil
	case <-stop:
		return nil
	}

	sb.setResultChan(results)

	// post block into BFT engine
	sb.postEvent(events.NewUnminedBlockEvent{
		NewUnminedBlock: *block,
	})

	return nil
}

func (sb *Backend) setResultChan(results chan<- *types.Block) {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()

	sb.commitCh = results
}

func (sb *Backend) sendResultChan(block *types.Block) {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()

	sb.commitCh <- block
}

func (sb *Backend) isResultChanNil() bool {
	sb.coreMu.RLock()
	defer sb.coreMu.RUnlock()

	return sb.commitCh == nil
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have based on the previous blocks in the blockchain and the
// current signer.
func (sb *Backend) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	return defaultDifficulty
}

func (sb *Backend) SetProposedBlockHash(hash common.Hash) {
	sb.proposedBlockHash = hash
}

// update timestamp and signature of the block based on its number of transactions
func (sb *Backend) AddSeal(block *types.Block) (*types.Block, error) {
	header := block.Header()
	// sign the hash
	seal, err := sb.Sign(types.SigHash(header).Bytes())
	if err != nil {
		return nil, err
	}

	err = types.WriteSeal(header, seal)
	if err != nil {
		return nil, err
	}

	return block.WithSeal(header), nil
}

// APIs returns the RPC APIs this consensus engine provides.
func (sb *Backend) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{{
		Namespace: "tendermint",
		Version:   "1.0",
		Service:   &API{chain: chain, tendermint: sb},
		Public:    true,
	}}
}

func (sb *Backend) SetBlockchain(bc *core.BlockChain) {
	sb.blockchainInitMu.Lock()
	sb.blockchain = bc
	sb.blockchainInitMu.Unlock()

	sb.currentBlock = bc.CurrentBlock
	sb.hasBadBlock = bc.HasBadBlock
}

// Start implements consensus.Start
func (sb *Backend) Start(ctx context.Context) error {
	// func (sb *Backend) Start(ctx context.Context, chain consensus.ChainReader, currentBlock func() *types.Block, hasBadBlock func(hash common.Hash) bool) error {
	// the mutex along with coreStarted should prevent double start
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if sb.coreStarted {
		return ErrStartedEngine
	}

	sb.stopped = make(chan struct{})

	// clear previous data
	sb.proposedBlockHash = common.Hash{}

	// Start Tendermint
	sb.core.Start(ctx)
	sb.coreStarted = true

	return nil
}

// Stop implements consensus.Stop
func (sb *Backend) Close() error {
	// the mutex along with coreStarted should prevent double stop
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()

	if !sb.coreStarted {
		return ErrStoppedEngine
	}

	// Stop Tendermint
	sb.core.Stop()
	sb.coreStarted = false
	close(sb.stopped)

	return nil
}

// retrieve list of committee for the block at height passed as parameter
func (sb *Backend) savedCommittee(number uint64, chain consensus.ChainReader) (*committee.Set, error) {
	var lastProposer common.Address
	var err error
	if number == 0 {
		number = 1
	}
	parentHeader := chain.GetHeaderByNumber(number - 1)
	if parentHeader == nil {
		return nil, errUnknownBlock
	}
	// For the genesis block, lastProposer is no one (empty).
	if number > 1 {
		lastProposer, err = sb.Author(parentHeader)
		if err != nil {
			return nil, err
		}
	}
	return committee.NewSet(parentHeader.Committee, lastProposer)
}

// retrieve list of committee for the block header passed as parameter
func (sb *Backend) committee(header *types.Header, parents []*types.Header, chain consensus.ChainReader) (*committee.Set, error) {

	// We can't use savedCommittee if parents are being passed :
	// those blocks are not yet saved in the blockchain.
	// autonity will stop processing the received blockchain from the moment an error appears.
	// See insertChain in blockchain.go

	if len(parents) > 0 {
		parent := parents[len(parents)-1]
		lastMiner, err := sb.Author(parent)
		if err != nil {
			return nil, err
		}
		return committee.NewSet(parent.Committee, lastMiner)
	} else {
		return sb.savedCommittee(header.Number.Uint64(), chain)
	}
}

func (sb *Backend) SealHash(header *types.Header) common.Hash {
	return types.SigHash(header)
}
