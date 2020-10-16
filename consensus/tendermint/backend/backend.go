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
	"crypto/ecdsa"
	"errors"
	"sync"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint"
	tendermintConfig "github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/ethdb"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/params"
	lru "github.com/hashicorp/golang-lru"
	ring "github.com/zfjagann/golang-ring"
)

const (
	// fetcherID is the ID indicates the block is from BFT engine
	fetcherID = "tendermint"
	// ring buffer to be able to handle at maximum 10 rounds, 20 committee and 3 messages types
	ringCapacity = 10 * 20 * 3
)

var (
	// ErrUnauthorizedAddress is returned when given address cannot be found in
	// current validator set.
	ErrUnauthorizedAddress = errors.New("unauthorized address")
	// ErrStoppedEngine is returned if the engine is stopped
	ErrStoppedEngine = errors.New("stopped engine")
)

// New creates an Ethereum Backend for BFT core engine.
func New(config *tendermintConfig.Config, privateKey *ecdsa.PrivateKey, db ethdb.Database, state state.Database, chainConfig *params.ChainConfig, vmConfig *vm.Config, broadcaster *tendermint.Broadcaster, peers consensus.Peers, syncer *tendermint.Syncer) *Backend {
	if chainConfig.Tendermint.BlockPeriod != 0 {
		config.BlockPeriod = chainConfig.Tendermint.BlockPeriod
	}

	recents, _ := lru.NewARC(inmemorySnapshots)
	recentMessages, _ := lru.NewARC(inmemoryPeers)

	pub := crypto.PubkeyToAddress(privateKey.PublicKey).String()
	logger := log.New("addr", pub)

	logger.Warn("new backend with public key")

	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	backend := &Backend{
		config:         config,
		eventMux:       event.NewTypeMuxSilent(logger),
		privateKey:     privateKey,
		address:        address,
		logger:         logger,
		db:             db,
		recents:        recents,
		coreStarted:    false,
		recentMessages: recentMessages,
		vmConfig:       vmConfig,
		peers:          peers,
	}

	backend.pendingMessages.SetCapacity(ringCapacity)
	backend.core = tendermint.New(backend, config, backend.privateKey, broadcaster, syncer, address, tendermint.NewLatestBlockRetriever(db, state), state)
	return backend
}

// ----------------------------------------------------------------------------

type Backend struct {
	config      *tendermintConfig.Config
	eventMux    *event.TypeMuxSilent
	privateKey  *ecdsa.PrivateKey
	address     common.Address
	logger      log.Logger
	db          ethdb.Database
	blockchain  *core.BlockChain
	hasBadBlock func(hash common.Hash) bool

	// the channels for tendermint engine notifications
	commitCh          chan<- *types.Block
	proposedBlockHash common.Hash
	coreStarted       bool
	core              tendermint.Tendermint
	stopped           chan struct{}
	coreMu            sync.RWMutex

	// Snapshots for recent block to speed up reorgs
	recents *lru.ARCCache

	// we save the last received p2p.messages in the ring buffer
	pendingMessages ring.Ring

	// event subscription for ChainHeadEvent event
	broadcaster consensus.Broadcaster

	//TODO: ARCChace is patented by IBM, so probably need to stop using it
	recentMessages *lru.ARCCache // the cache of peer's messages

	contractsMu sync.RWMutex
	vmConfig    *vm.Config
	peers       consensus.Peers
}

// Commit implements tendermint.Backend.Commit
func (sb *Backend) Commit(block *types.Block, proposer common.Address) {
	sb.logger.Info("Committed", "address", sb.address, "proposer", proposer, "hash", block.Hash(), "number", block.Number().Uint64())
	// - if we are the proposer, send the proposed hash to commit channel,
	//    which is being watched inside the engine.Seal() function.
	// - otherwise, we try to insert the block.
	// -- if success, the ChainHeadEvent event will be broadcasted, try to build
	//    the next block and the previous Seal() will be stopped.
	// -- otherwise, a error will be returned and a round change event will be fired.
	if sb.address == proposer && !sb.isResultChanNil() {
		// feed block hash to Seal() and wait the Seal() result
		sb.sendResultChan(block)
		return
	}

	if sb.broadcaster != nil {
		sb.broadcaster.Enqueue(fetcherID, block)
	}
}

func (sb *Backend) Post(ev interface{}) {
	sb.eventMux.Post(ev)
}

func (sb *Backend) Subscribe(types ...interface{}) *event.TypeMuxSubscription {
	return sb.eventMux.Subscribe(types...)
}

// VerifyProposal implements tendermint.Backend.VerifyProposal
func (sb *Backend) VerifyProposal(proposal types.Block) (time.Duration, error) {
	// Check if the proposal is a valid block
	// TODO: fix always false statement and check for non nil
	// TODO: use interface instead of type
	block := &proposal
	//if block == nil {
	//	sb.logger.Error("Invalid proposal, %v", proposal)
	//	return 0, errInvalidProposal
	//}

	// check bad block
	if sb.blockchain.HasBadBlock(block.Hash()) {
		return 0, core.ErrBlacklistedHash
	}

	// verify the header of proposed block
	err := sb.VerifyHeader(sb.blockchain, block.Header(), false)
	// ignore errEmptyCommittedSeals error because we don't have the committed seals yet
	if err == nil || err == types.ErrEmptyCommittedSeals {
		var (
			receipts types.Receipts

			usedGas        = new(uint64)
			gp             = new(core.GasPool).AddGas(block.GasLimit())
			header         = block.Header()
			proposalNumber = header.Number.Uint64()
			parent         = sb.blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1)
		)

		// We need to process all of the transaction to get the latest state to get the latest committee
		state, stateErr := sb.blockchain.StateAt(parent.Root())
		if stateErr != nil {
			return 0, stateErr
		}

		// Validate the body of the proposal
		if err = sb.blockchain.Validator().ValidateBody(block); err != nil {
			return 0, err
		}

		// sb.blockchain.Processor().Process() was not called because it calls back Finalize() and would have modified the proposal
		// Instead only the transactions are applied to the copied state
		for i, tx := range block.Transactions() {
			state.Prepare(tx.Hash(), block.Hash(), i)
			// Might be vulnerable to DoS Attack depending on gaslimit
			// Todo : Double check
			receipt, receiptErr := core.ApplyTransaction(sb.blockchain.Config(), sb.blockchain, nil, gp, state, header, tx, usedGas, *sb.vmConfig)
			if receiptErr != nil {
				return 0, receiptErr
			}
			receipts = append(receipts, receipt)
		}

		state.Prepare(common.ACHash(block.Number()), block.Hash(), len(block.Transactions()))
		committeeSet, receipt, err := sb.Finalize(sb.blockchain, header, state, block.Transactions(), nil, receipts)
		receipts = append(receipts, receipt)
		//Validate the state of the proposal
		if err = sb.blockchain.Validator().ValidateState(block, state, receipts, *usedGas); err != nil {
			return 0, err
		}

		//Perform the actual comparison
		if len(header.Committee) != len(committeeSet) {
			sb.logger.Error("wrong committee set",
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
				sb.logger.Error("wrong committee member in the set",
					"index", i,
					"currentVerifier", sb.address.String(),
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
		return time.Unix(int64(block.Header().Time), 0).Sub(now()), consensus.ErrFutureBlock
	}
	return 0, err
}

func (sb *Backend) GetContractABI() string {
	// after the contract is upgradable, call it from contract object rather than from conf.
	return sb.blockchain.GetAutonityContract().GetContractABI()
}

// Whitelist for the current block
func (sb *Backend) WhiteList() []string {
	db, err := sb.blockchain.State()
	if err != nil {
		sb.logger.Error("Failed to get block white list", "err", err)
		return nil
	}

	enodes, err := sb.blockchain.GetAutonityContract().GetWhitelist(sb.blockchain.CurrentBlock(), db)
	if err != nil {
		sb.logger.Error("Failed to get block white list", "err", err)
		return nil
	}

	return enodes.StrList
}

func (sb *Backend) ResetPeerCache(address common.Address) {
	ms, ok := sb.recentMessages.Get(address)
	var m *lru.ARCCache
	if ok {
		m, _ = ms.(*lru.ARCCache)
		m.Purge()
	}
}
