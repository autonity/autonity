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
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/params"
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
func New(config *tendermintConfig.Config, privateKey *ecdsa.PrivateKey, db ethdb.Database, statedb state.Database, chainConfig *params.ChainConfig, vmConfig *vm.Config, broadcaster *tendermint.Broadcaster, peers consensus.Peers, syncer *tendermint.Syncer, autonityContract *autonity.Contract, verifier *tendermint.Verifier, finalizer *tendermint.DefaultFinalizer) *Backend {
	if chainConfig.Tendermint.BlockPeriod != 0 {
		config.BlockPeriod = chainConfig.Tendermint.BlockPeriod
	}

	pub := crypto.PubkeyToAddress(privateKey.PublicKey).String()
	logger := log.New("addr", pub)

	logger.Warn("new backend with public key")
	latestBlockRetriever := tendermint.NewLatestBlockRetriever(db, statedb)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	backend := &Backend{
		config:               config,
		privateKey:           privateKey,
		address:              address,
		logger:               logger,
		latestBlockRetriever: latestBlockRetriever,
		coreStarted:          false,
		vmConfig:             vmConfig,
		peers:                peers,
		autonityContract:     autonityContract,
		Verifier:             verifier,
		DefaultFinalizer:     finalizer,
	}

	backend.core = tendermint.New(config, backend.privateKey, broadcaster, syncer, address, latestBlockRetriever, statedb, verifier, autonityContract)
	return backend
}

// ----------------------------------------------------------------------------

type Backend struct {
	*tendermint.DefaultFinalizer
	*tendermint.Verifier
	config     *tendermintConfig.Config
	privateKey *ecdsa.PrivateKey
	address    common.Address
	logger     log.Logger
	blockchain *core.BlockChain

	// the channels for tendermint engine notifications
	commitCh          chan<- *types.Block
	proposedBlockHash common.Hash
	coreStarted       bool
	core              tendermint.Tendermint
	stopped           chan struct{}
	coreMu            sync.RWMutex

	// event subscription for ChainHeadEvent event
	broadcaster consensus.Broadcaster

	contractsMu sync.RWMutex
	vmConfig    *vm.Config
	peers       consensus.Peers

	autonityContract     *autonity.Contract
	latestBlockRetriever *tendermint.LatestBlockRetriever
}
