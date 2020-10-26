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

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint"
	tendermintConfig "github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/contracts/autonity"
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
func New(config *tendermintConfig.Config, privateKey *ecdsa.PrivateKey, db ethdb.Database, statedb state.Database, chainConfig *params.ChainConfig, vmConfig *vm.Config, broadcaster *tendermint.Broadcaster, peers consensus.Peers, syncer *tendermint.Syncer, autonityContract *autonity.Contract, verifier *tendermint.Verifier) *Backend {
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
		config:               config,
		eventMux:             event.NewTypeMuxSilent(logger),
		privateKey:           privateKey,
		address:              address,
		logger:               logger,
		db:                   db,
		recents:              recents,
		coreStarted:          false,
		recentMessages:       recentMessages,
		vmConfig:             vmConfig,
		peers:                peers,
		statedb:              statedb,
		latestBlockRetreiver: tendermint.NewLatestBlockRetriever(db, statedb),
		autonityContract:     autonityContract,
		Verifier:             verifier,
	}

	backend.core = tendermint.New(backend, config, backend.privateKey, broadcaster, syncer, address, tendermint.NewLatestBlockRetriever(db, statedb), statedb, verifier)
	return backend
}

// ----------------------------------------------------------------------------

type Backend struct {
	*tendermint.Verifier
	config     *tendermintConfig.Config
	eventMux   *event.TypeMuxSilent
	privateKey *ecdsa.PrivateKey
	address    common.Address
	logger     log.Logger
	db         ethdb.Database
	blockchain *core.BlockChain

	// the channels for tendermint engine notifications
	commitCh          chan<- *types.Block
	proposedBlockHash common.Hash
	coreStarted       bool
	core              tendermint.Tendermint
	stopped           chan struct{}
	coreMu            sync.RWMutex

	// Snapshots for recent block to speed up reorgs
	recents *lru.ARCCache

	// event subscription for ChainHeadEvent event
	broadcaster consensus.Broadcaster

	//TODO: ARCChace is patented by IBM, so probably need to stop using it
	recentMessages *lru.ARCCache // the cache of peer's messages

	contractsMu sync.RWMutex
	vmConfig    *vm.Config
	peers       consensus.Peers

	autonityContract     *autonity.Contract
	statedb              state.Database
	latestBlockRetreiver *tendermint.LatestBlockRetriever
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

func (sb *Backend) GetContractABI() string {
	// after the contract is upgradable, call it from contract object rather than from conf.
	return sb.autonityContract.GetContractABI()
}

// Whitelist for the current block
func (sb *Backend) WhiteList() []string {
	// TODO this should really return errors
	b, err := sb.latestBlockRetreiver.RetrieveLatestBlock()
	if err != nil {
		panic(err)
	}
	state, err := state.New(b.Root(), sb.statedb, nil)
	if err != nil {
		sb.logger.Error("Failed to get block white list", "err", err)
		return nil
	}

	enodes, err := sb.autonityContract.GetWhitelist(b, state)
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
