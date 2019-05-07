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
	"errors"
	"math/big"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/hexutil"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint"
	tendermintCore "github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/rpc"
)

const (
	inmemorySnapshots = 128 // Number of recent vote snapshots to keep in memory
	inmemoryPeers     = 40
	inmemoryMessages  = 1024
)

var (
	// errInvalidProposal is returned when a prposal is malformed.
	errInvalidProposal = errors.New("invalid proposal")
	// errUnknownBlock is returned when the list of validators is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")
	// errUnauthorized is returned if a header is signed by a non authorized entity.
	errUnauthorized = errors.New("unauthorized")
	// errInvalidDifficulty is returned if the difficulty of a block is not 1
	errInvalidDifficulty = errors.New("invalid difficulty")
	// errInvalidExtraDataFormat is returned when the extra data format is incorrect
	errInvalidExtraDataFormat = errors.New("invalid extra data format")
	// errInvalidMixDigest is returned if a block's mix digest is not PoS digest.
	errInvalidMixDigest = errors.New("invalid PoS mix digest")
	// errInvalidNonce is returned if a block's nonce is invalid
	errInvalidNonce = errors.New("invalid nonce")
	// errInvalidUncleHash is returned if a block contains an non-empty uncle list.
	errInvalidUncleHash = errors.New("non empty uncle hash")
	// errInconsistentValidatorSet is returned if the validator set is inconsistent
	errInconsistentValidatorSet = errors.New("inconsistent validator set")
	// errInvalidTimestamp is returned if the timestamp of a block is lower than the previous block's timestamp + the minimum block period.
	errInvalidTimestamp = errors.New("invalid timestamp")
	// errInvalidVotingChain is returned if an authorization list is attempted to
	// be modified via out-of-range or non-contiguous headers.
	errInvalidVotingChain = errors.New("invalid voting blockchain")
	// errInvalidVote is returned if a nonce value is something else that the two
	// allowed constants of 0x00..0 or 0xff..f.
	errInvalidVote = errors.New("vote nonce not 0x00..0 or 0xff..f")
	// errMismatchTxhashes is returned if the TxHash in header is mismatch.
	errMismatchTxhashes = errors.New("mismatch transactions hashes")
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
func (sb *backend) Author(header *types.Header) (common.Address, error) {
	return types.Ecrecover(header)
}

// VerifyHeader checks whether a header conforms to the consensus rules of a
// given engine. Verifying the seal may be done optionally here, or explicitly
// via the VerifySeal method.
func (sb *backend) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	return sb.verifyHeader(chain, header, nil)
}

// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (sb *backend) verifyHeader(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	if header.Number == nil {
		return errUnknownBlock
	}

	// Don't waste time checking blocks from the future
	if header.Time.Cmp(big.NewInt(now().Unix())) > 0 {
		return consensus.ErrFutureBlock
	}

	// Ensure that the extra data format is satisfied
	if _, err := types.ExtractPoSExtra(header); err != nil {
		return errInvalidExtraDataFormat
	}

	// Ensure that the coinbase is valid
	if header.Nonce != (emptyNonce) && !bytes.Equal(header.Nonce[:], nonceAuthVote) && !bytes.Equal(header.Nonce[:], nonceDropVote) {
		return errInvalidNonce
	}
	// Ensure that the mix digest is zero as we don't have fork protection currently
	if header.MixDigest != types.PoSDigest {
		return errInvalidMixDigest
	}
	// Ensure that the block doesn't contain any uncles which are meaningless in PoS
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
func (sb *backend) verifyCascadingFields(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
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
	if parent.Time.Uint64()+sb.config.BlockPeriod > header.Time.Uint64() {
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
func (sb *backend) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
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
func (sb *backend) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	if len(block.Uncles()) > 0 {
		return errInvalidUncleHash
	}
	return nil
}

// verifySigner checks whether the signer is in parent's validator set
func (sb *backend) verifySigner(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}

	validators, err := sb.retrieveValidators(header, parents, chain)

	if err != nil {
		return err
	}

	// resolve the authorization key and check against signers
	signer, err := types.Ecrecover(header)
	if err != nil {
		return err
	}

	// Signer should be in the validator set of previous block's extraData.

	for i := range validators {
		if validators[i] == signer {
			return nil
		}
	}
	return errUnauthorized
}

// verifyCommittedSeals checks whether every committed seal is signed by one of the parent's validators
func (sb *backend) verifyCommittedSeals(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	number := header.Number.Uint64()
	// We don't need to verify committed seals in the genesis block
	if number == 0 {
		return nil
	}

	validatorAddresses, err := sb.retrieveValidators(header, parents, chain)
	if err != nil {
		return err
	}
	validators := validator.NewSet(validatorAddresses, sb.config.GetProposerPolicy())

	extra, err := types.ExtractPoSExtra(header)
	if err != nil {
		return err
	}
	// The length of Committed seals should be larger than 0
	if len(extra.CommittedSeal) == 0 {
		return types.ErrEmptyCommittedSeals
	}

	// Check whether the committed seals are generated by parent's validators
	validSeal := 0
	proposalSeal := tendermintCore.PrepareCommittedSeal(header.Hash())
	// 1. Get committed seals from current header
	for _, seal := range extra.CommittedSeal {
		// 2. Get the original address by seal and parent block hash
		addr, err := types.GetSignatureAddress(proposalSeal, seal)
		if err != nil {
			sb.logger.Error("not a valid address", "err", err)
			return types.ErrInvalidSignature
		}
		// Every validator can have only one seal. If more than one seals are signed by a
		// validator, the validator cannot be found and errInvalidCommittedSeals is returned.
		if validators.RemoveValidator(addr) {
			validSeal += 1
		} else {
			return types.ErrInvalidCommittedSeals
		}
	}

	// The length of validSeal should be larger than number of faulty node + 1
	if validSeal <= 2*validators.F() {
		return types.ErrInvalidCommittedSeals
	}

	return nil
}

// VerifySeal checks whether the crypto seal on a header is valid according to
// the consensus rules of the given engine.
func (sb *backend) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
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
func (sb *backend) Prepare(chain consensus.ChainReader, header *types.Header) error {
	// unused fields, force to set to empty
	header.Coinbase = sb.address
	header.Nonce = emptyNonce
	header.MixDigest = types.PoSDigest

	// copy the parent extra data as the header extra data
	number := header.Number.Uint64()
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	// use the same difficulty for all blocks
	header.Difficulty = defaultDifficulty

	// set header's timestamp
	header.Time = new(big.Int).Add(parent.Time, new(big.Int).SetUint64(sb.config.BlockPeriod))
	if header.Time.Int64() < time.Now().Unix() {
		header.Time = big.NewInt(time.Now().Unix())
	}
	return nil
}

// Finalize runs any post-transaction state modifications (e.g. block rewards)
// and assembles the final block.
//
// Note, the block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (sb *backend) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {

	if sb.blockchain == nil {
		sb.blockchain = chain.(*core.BlockChain) // in the case of Finalize() called before the engine start()
	}

	validators, err := sb.getValidators(header, chain, state)
	// No block rewards in PoS, so the state remains as is and uncles are dropped
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = nilUncleHash

	// add validators to extraData's validators section
	if header.Extra, err = types.PrepareExtra(&header.Extra, validators); err != nil {
		return nil, err
	}

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, nil, receipts), nil
}

func (sb *backend) getValidators(header *types.Header, chain consensus.ChainReader, state *state.StateDB) ([]common.Address, error) {
	var validators []common.Address
	var err error

	if header.Number.Int64() == 1 {
		// Deploy Soma on-blockchain governance contract
		log.Info("Soma Contract Deployer", "Address", sb.config.Deployer)
		contractAddress, err := sb.deployContract(chain, header, state)
		if err != nil {
			return nil, err
		}
		sb.somaContract = contractAddress
		validators, err = sb.retrieveSavedValidators(1, chain)
		if err != nil {
			return nil, err
		}

		// Deploy Glienicke network-permissioning contract
		_, sb.glienickeContract, err = sb.blockchain.DeployGlienickeContract(state, header)
		if err != nil {
			return nil, err
		}

	} else {
		if sb.somaContract == common.HexToAddress("0000000000000000000000000000000000000000") {
			sb.somaContract = crypto.CreateAddress(sb.config.Deployer, 0)
		}

		if sb.glienickeContract == common.HexToAddress("0000000000000000000000000000000000000000") {
			sb.glienickeContract = crypto.CreateAddress(sb.blockchain.Config().GlienickeDeployer, 0)
		}

		validators, err = sb.contractGetValidators(chain, header, state)
		if err != nil {
			return nil, err
		}
	}

	return validators, err
}

// Seal generates a new block for the given input block with the local miner's
// seal place on top.
func (sb *backend) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	// update the block header and signature and propose the block to core engine
	header := block.Header()
	number := header.Number.Uint64()

	// Bail out if we're unauthorized to sign a block
	if _, v := sb.Validators(number).GetByAddress(sb.address); v == nil {
		sb.logger.Error("Error validator")
		return errUnauthorized
	}

	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		sb.logger.Error("Error ancestor")
		return consensus.ErrUnknownAncestor
	}
	block, err := sb.updateBlock(parent, block)
	if err != nil {
		sb.logger.Error("Error updateBlock", err, err.Error())
		return err
	}

	// wait for the timestamp of header, use this to adjust the block period
	delay := time.Unix(block.Header().Time.Int64(), 0).Sub(now())
	select {
	case <-time.After(delay):
	case <-stop:
		return nil
	}

	sb.commitCh = results // results channel stays always the same

	// post block into PoS engine
	go sb.EventMux().Post(tendermint.RequestEvent{
		ProposalBlock: block,
	})

	return nil
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have based on the previous blocks in the blockchain and the
// current signer.
func (sb *backend) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	return defaultDifficulty
}

func (sb *backend) SetProposedBlockHash(hash common.Hash) {
	sb.proposedBlockHash = hash
}

// update timestamp and signature of the block based on its number of transactions
func (sb *backend) updateBlock(parent *types.Header, block *types.Block) (*types.Block, error) {
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
func (sb *backend) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{{
		Namespace: "tendermint",
		Version:   "1.0",
		Service:   &API{chain: chain, tendermint: sb},
		Public:    true,
	}}
}

// Start implements consensus.tendermint.Start
func (sb *backend) Start(chain consensus.ChainReader, currentBlock func() *types.Block, hasBadBlock func(hash common.Hash) bool) error {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if sb.coreStarted {
		return tendermint.ErrStartedEngine
	}

	// clear previous data
	sb.proposedBlockHash = common.Hash{}

	sb.blockchain = chain.(*core.BlockChain)
	sb.currentBlock = currentBlock
	sb.hasBadBlock = hasBadBlock

	if err := sb.core.Start(); err != nil {
		return err
	}

	sb.coreStarted = true
	return nil
}

// Stop implements consensus.tendermint.Stop
func (sb *backend) Close() error {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if !sb.coreStarted {
		return tendermint.ErrStoppedEngine
	}
	if err := sb.core.Stop(); err != nil {
		return err
	}
	sb.coreStarted = false
	return nil
}

// retrieve list of validators for the block at height passed as parameter
func (sb *backend) retrieveSavedValidators(number uint64, chain consensus.ChainReader) ([]common.Address, error) {

	if number == 0 {
		number = 1
	}

	header := chain.GetHeaderByNumber(number - 1)
	if header == nil {
		return nil, errUnknownBlock
	}

	tendermintExtra, err := types.ExtractPoSExtra(header)
	if err != nil {
		return nil, err
	}

	return tendermintExtra.Validators, nil

}

// retrieve list of validators for the block header passed as parameter
func (sb *backend) retrieveValidators(header *types.Header, parents []*types.Header, chain consensus.ChainReader) ([]common.Address, error) {

	var validators []common.Address
	var err error

	/*
		We can't use retrieveSavedValidators if parents are being passed :
		those blocks are not yet saved in the blockchain.
		autonity will stop processing the received blockchain from the moment an error appears.
		See insertChain in blockchain.go
	*/

	if len(parents) > 0 {
		parent := parents[len(parents)-1]
		var tendermintExtra *types.PoSExtra
		tendermintExtra, err = types.ExtractPoSExtra(parent)
		if err == nil {
			validators = tendermintExtra.Validators
		}
	} else {
		validators, err = sb.retrieveSavedValidators(header.Number.Uint64(), chain)
	}
	return validators, err

}

func (sb *backend) SealHash(header *types.Header) common.Hash {
	return types.SigHash(header)
}
