// Copyright 2014 The go-ethereum Authors
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

// Package eth implements the Ethereum protocol.
package eth

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"time"

	"github.com/autonity/autonity/accounts"
	"github.com/autonity/autonity/accounts/abi/bind/backends"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/accountability"
	"github.com/autonity/autonity/consensus/tendermint/backend"
	tendermintcore "github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/bloombits"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/state/pruner"
	"github.com/autonity/autonity/core/txpool"
	"github.com/autonity/autonity/core/txpool/legacypool"
	"github.com/autonity/autonity/core/txpool/locals"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/eth/downloader"
	"github.com/autonity/autonity/eth/ethconfig"
	"github.com/autonity/autonity/eth/gasprice"
	"github.com/autonity/autonity/eth/protocols/eth"
	"github.com/autonity/autonity/eth/protocols/snap"
	"github.com/autonity/autonity/eth/tracers"
	"github.com/autonity/autonity/ethdb"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/internal/ethapi"
	"github.com/autonity/autonity/internal/shutdowncheck"
	"github.com/autonity/autonity/internal/version"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/miner"
	"github.com/autonity/autonity/node"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/p2p/dnsdisc"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/rlp"
	"github.com/autonity/autonity/rpc"
)

const (
	maxFullMeshPeers = 20
)

// Config contains the configuration options of the ETH protocol.
// Deprecated: use ethconfig.Config instead.
type Config = ethconfig.Config

// Ethereum implements the Ethereum full node service.
type Ethereum struct {
	// core protocol objects
	config         *ethconfig.Config
	txPool         *txpool.TxPool
	localTxTracker *locals.TxTracker
	blockchain     *core.BlockChain

	handler *handler
	discmix *enode.FairMix

	// DB interfaces
	chainDb ethdb.Database // Block chain database

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	bloomRequests     chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer      *core.ChainIndexer             // Bloom indexer operating during block imports
	closeBloomHandler chan struct{}

	APIBackend *EthAPIBackend

	miner    *miner.Miner
	gasPrice *big.Int

	networkID     uint64
	netRPCService *ethapi.NetAPI

	p2pServer *p2p.Server

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price and address)

	shutdownTracker *shutdowncheck.ShutdownTracker // Tracks if and when the node has shutdown ungracefully

	// Autonity Addons
	faultDetector    *accountability.FaultDetector
	consensusServer  *p2p.Server
	address          common.Address // Local node address
	log              log.Logger
	topologySelector networkTopology
}

// New creates a new Ethereum object (including the
// initialisation of the common Ethereum object)
func New(stack *node.Node, config *ethconfig.Config) (*Ethereum, error) {
	// Ensure configuration values are compatible and sane
	if !config.SyncMode.IsValid() {
		return nil, fmt.Errorf("invalid sync mode %d", config.SyncMode)
	}

	if config.Miner.GasPrice == nil || config.Miner.GasPrice.Cmp(common.Big0) <= 0 {
		stack.Logger().Warn("Sanitizing invalid miner gas price", "provided", config.Miner.GasPrice, "updated", ethconfig.Defaults.Miner.GasPrice)
		config.Miner.GasPrice = new(big.Int).Set(ethconfig.Defaults.Miner.GasPrice)
	}
	if config.NoPruning && config.TrieDirtyCache > 0 {
		if config.SnapshotCache > 0 {
			config.TrieCleanCache += config.TrieDirtyCache * 3 / 5
			config.SnapshotCache += config.TrieDirtyCache * 2 / 5
		} else {
			config.TrieCleanCache += config.TrieDirtyCache
		}
		config.TrieDirtyCache = 0
	}
	stack.Logger().Info("Allocated trie memory caches", "clean", common.StorageSize(config.TrieCleanCache)*1024*1024, "dirty", common.StorageSize(config.TrieDirtyCache)*1024*1024)

	// Assemble the Ethereum object
	chainDb, err := stack.OpenDatabaseWithFreezer("chaindata", config.DatabaseCache, config.DatabaseHandles, config.DatabaseFreezer, "eth/db/chaindata/", false)
	if err != nil {
		return nil, err
	}
	scheme, err := rawdb.ParseStateScheme(config.StateScheme, chainDb)
	if err != nil {
		return nil, err
	}
	// Try to recover offline state pruning only in hash-based.
	if scheme == rawdb.HashScheme {
		if err := pruner.RecoverPruning(stack.ResolvePath(""), chainDb); err != nil {
			log.Error("Failed to recover state", "error", err)
		}
	}
	// Transfer mining-related config to the ethash config.
	chainConfig, err := core.LoadChainConfig(chainDb, config.Genesis)
	if err != nil {
		return nil, err
	}

	var (
		vmConfig = vm.Config{
			EnablePreimageRecording: config.EnablePreimageRecording,
		}
		cacheConfig = &core.CacheConfig{
			TrieCleanLimit:      config.TrieCleanCache,
			TrieCleanNoPrefetch: config.NoPrefetch,
			TrieDirtyLimit:      config.TrieDirtyCache,
			TrieDirtyDisabled:   config.NoPruning,
			TrieTimeLimit:       config.TrieTimeout,
			SnapshotLimit:       config.SnapshotCache,
			Preimages:           config.Preimages,
			StateHistory:        config.StateHistory,
			StateScheme:         scheme,
		}
	)
	if config.VMTrace != "" {
		traceConfig := json.RawMessage("{}")
		if config.VMTraceJsonConfig != "" {
			traceConfig = json.RawMessage(config.VMTraceJsonConfig)
		}
		t, err := tracers.LiveDirectory.New(config.VMTrace, traceConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create tracer %s: %v", config.VMTrace, err)
		}
		vmConfig.Tracer = t
	}
	networkID := config.NetworkId
	if networkID == 0 {
		networkID = chainConfig.ChainID.Uint64()
	}
	// The event mux is shared between the consensus engine and fault detector such that both of them will receive the
	// messages from p2p protocol manager layer.
	evMux := new(event.TypeMux)
	// single instance of msgStore shared by misbehaviour detector and omission fault detector.
	msgStore := tendermintcore.NewMsgStore()
	consensusEngine := ethconfig.CreateConsensusEngine(chainDb, stack, &vmConfig, evMux, msgStore)

	stack.Logger().Info("Initialised chain configuration", "config", chainConfig)

	nodeKey, _ := stack.Config().AutonityKeys()

	eth := &Ethereum{
		config:            config,
		chainDb:           chainDb,
		eventMux:          stack.EventMux(),
		accountManager:    stack.AccountManager(),
		engine:            consensusEngine,
		closeBloomHandler: make(chan struct{}),
		networkID:         networkID,
		gasPrice:          config.Miner.GasPrice,
		bloomRequests:     make(chan chan *bloombits.Retrieval),
		bloomIndexer:      core.NewBloomIndexer(chainDb, params.BloomBitsBlocks, params.BloomConfirms),
		p2pServer:         stack.ExecutionServer(),
		discmix:           enode.NewFairMix(0),
		// Autonity stuff:
		shutdownTracker:  shutdowncheck.NewShutdownTracker(chainDb),
		log:              stack.Logger(),
		address:          crypto.PubkeyToAddress(nodeKey.PublicKey),
		consensusServer:  stack.ConsensusServer(),
		topologySelector: NewGraphTopology(maxFullMeshPeers),
	}

	bcVersion := rawdb.ReadDatabaseVersion(chainDb)
	var dbVer = "<nil>"
	if bcVersion != nil {
		dbVer = fmt.Sprintf("%d", *bcVersion)
	}
	eth.log.Info("Initialising Autonity protocol", "network", config.NetworkId, "dbversion", dbVer)

	if !config.SkipBcVersionCheck {
		if bcVersion != nil && *bcVersion > core.BlockChainVersion {
			return nil, fmt.Errorf("database version is v%d, Geth %s only supports v%d", *bcVersion, version.WithMeta, core.BlockChainVersion)
		} else if bcVersion == nil || *bcVersion < core.BlockChainVersion {
			if bcVersion != nil { // only print warning on upgrade, not on init
				log.Warn("Upgrade blockchain database version", "from", dbVer, "to", core.BlockChainVersion)
			}
			rawdb.WriteDatabaseVersion(chainDb, core.BlockChainVersion)
		}
	}

	txSender := func(tx *types.Transaction) error {
		return eth.txPool.Add([]*types.Transaction{tx}, true)[0]
	}

	eth.APIBackend = &EthAPIBackend{stack.Config().ExtRPCEnabled(), stack.Config().AllowUnprotectedTxs, eth, nil}
	if eth.APIBackend.allowUnprotectedTxs {
		log.Info("Unprotected transactions allowed")
	}
	eth.APIBackend.gpo = gasprice.NewOracle(eth.APIBackend, config.GPO, config.Miner.GasPrice)

	eth.blockchain, err = core.NewBlockChain(
		chainDb,
		cacheConfig,
		config.Genesis,
		eth.engine,
		vmConfig,
		&config.TransactionHistory,
		backends.NewInternalBackend(txSender, eth.APIBackend),
		eth.log)

	if err != nil {
		return nil, err
	}

	// temporary solution
	if be, ok := consensusEngine.(interface {
		SetBlockchain(*core.BlockChain)
	}); ok {
		be.SetBlockchain(eth.blockchain)
	}

	eth.bloomIndexer.Start(eth.blockchain)

	if config.TxPool.Journal != "" {
		config.TxPool.Journal = stack.ResolvePath(config.TxPool.Journal)
	}
	legacyPool := legacypool.New(config.TxPool, eth.blockchain)
	eth.txPool, err = txpool.New(config.TxPool.PriceLimit, eth.blockchain, []txpool.SubPool{legacyPool})
	if err != nil {
		return nil, err
	}
	if !config.TxPool.NoLocals {
		rejournal := config.TxPool.Rejournal
		if rejournal < time.Second {
			log.Warn("Sanitizing invalid txpool journal time", "provided", rejournal, "updated", time.Second)
			rejournal = time.Second
		}
		eth.localTxTracker = locals.New(config.TxPool.Journal, rejournal, eth.blockchain.Config(), eth.txPool)
		stack.RegisterLifecycle(eth.localTxTracker)
	}

	// Permit the downloader to use the trie cache allowance during fast sync
	cacheLimit := cacheConfig.TrieCleanLimit + cacheConfig.TrieDirtyLimit + cacheConfig.SnapshotLimit
	if eth.handler, err = newHandler(&handlerConfig{
		NodeID:         eth.p2pServer.Self().ID(),
		Database:       chainDb,
		Chain:          eth.blockchain,
		TxPool:         eth.txPool,
		Network:        networkID,
		Sync:           config.SyncMode,
		BloomCache:     uint64(cacheLimit),
		EventMux:       eth.eventMux,
		RequiredBlocks: config.RequiredBlocks,
	}); err != nil {
		return nil, err
	}

	eth.miner = miner.New(eth, &config.Miner, chainConfig, eth.EventMux(), eth.engine)
	eth.miner.SetExtra(makeExtraData(config.Miner.ExtraData))
	eth.miner.SetPrioAddresses(config.TxPool.Locals)

	eth.netRPCService = ethapi.NewNetAPI(eth.p2pServer, networkID)

	// Once the chain is initialized, load accountability precompiled contracts in EVM environment before chain sync
	//start to apply accountability TXs if there were any, otherwise it would cause sync failure.
	accountability.LoadPrecompiles() // todo(youssef): doesn't conceptually belong here, should be moved somewhere else
	// Create Fault Detector for each full node for the time being.
	//TODO: I think it would make more sense to move this into the tendermint backend if possible
	eth.faultDetector = accountability.NewFaultDetector(
		eth.blockchain,
		eth.address,
		evMux.Subscribe(events.MessageEvent{}, events.AccountabilityEvent{}, events.OldMessageEvent{}),
		msgStore,
		eth.APIBackend,
		nodeKey,
		eth.blockchain.ProtocolContracts(),
		eth.log)

	// Start the RPC service
	eth.netRPCService = ethapi.NewNetAPI(eth.p2pServer, networkID)

	// Register the backend on the node
	stack.RegisterAPIs(eth.APIs())
	stack.RegisterProtocols(eth.Protocols())
	stack.RegisterLifecycle(eth)

	// Successful startup; push a marker and check previous unclean shutdowns.
	eth.shutdownTracker.MarkStartup()

	return eth, nil
}

func makeExtraData(extra []byte) []byte {
	if len(extra) == 0 {
		// create default extradata
		extra, _ = rlp.EncodeToBytes([]interface{}{
			uint(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch),
			"autonity",
			runtime.Version(),
			runtime.GOOS,
		})
	}
	if uint64(len(extra)) > params.MaximumExtraDataSize {
		log.Warn("Miner extra data exceed limit", "extra", hexutil.Bytes(extra), "limit", params.MaximumExtraDataSize)
		extra = nil
	}
	return extra
}

// APIs return the collection of RPC services the ethereum package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *Ethereum) APIs() []rpc.API {
	apis := ethapi.GetAPIs(s.APIBackend)

	// Append any APIs exposed explicitly by the consensus engine
	apis = append(apis, s.engine.APIs(s.BlockChain())...)

	/* Todo(youssef): revisit aut api
	if _, ok := s.engine.(consensus.BFT); ok {
		apis = append(apis, rpc.API{
			Namespace: "aut",
			Version:   params.Version,
			Service:   NewAutonityContractAPI(s.BlockChain(), s.BlockChain().ProtocolContracts(), s.consensusServer),
			Public:    true,
		})
	}
	*/
	return append(apis, []rpc.API{
		{
			Namespace: "miner",
			Service:   NewMinerAPI(s),
		}, {
			Namespace: "eth",
			Service:   downloader.NewDownloaderAPI(s.handler.downloader, s.blockchain, s.eventMux),
		}, {
			Namespace: "admin",
			Service:   NewAdminAPI(s),
		}, {
			Namespace: "debug",
			Service:   NewDebugAPI(s),
		}, {
			Namespace: "net",
			Service:   s.netRPCService,
		},
	}...)
}

func (s *Ethereum) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *Ethereum) Etherbase() (eb common.Address, err error) {
	s.lock.RLock()
	etherbase := s.address
	s.lock.RUnlock()

	if etherbase != (common.Address{}) {
		return etherbase, nil
	}
	return common.Address{}, fmt.Errorf("address must be explicitly specified")
}

// StopMining terminates the miner, both at the consensus engine level as well as
// at the block creation level.
func (s *Ethereum) StopMining() {
	// Update the thread count within the consensus engine
	type threaded interface {
		SetThreads(threads int)
	}
	if th, ok := s.engine.(threaded); ok {
		th.SetThreads(-1)
	}
	// Stop the block creating itself
	s.miner.Stop()
}

func (s *Ethereum) IsMining() bool      { return s.miner.Mining() }
func (s *Ethereum) Miner() *miner.Miner { return s.miner }

func (s *Ethereum) AccountManager() *accounts.Manager            { return s.accountManager }
func (s *Ethereum) BlockChain() *core.BlockChain                 { return s.blockchain }
func (s *Ethereum) TxPool() *txpool.TxPool                       { return s.txPool }
func (s *Ethereum) EventMux() *event.TypeMux                     { return s.eventMux }
func (s *Ethereum) Engine() consensus.Engine                     { return s.engine }
func (s *Ethereum) FaultDetector() *accountability.FaultDetector { return s.faultDetector }
func (s *Ethereum) ChainDb() ethdb.Database                      { return s.chainDb }
func (s *Ethereum) IsListening() bool                            { return true } // Always listening
func (s *Ethereum) Downloader() *downloader.Downloader           { return s.handler.downloader }
func (s *Ethereum) Synced() bool                                 { return s.handler.synced.Load() }
func (s *Ethereum) SetSynced()                                   { s.handler.enableSyncedFeatures() }
func (s *Ethereum) ArchiveMode() bool                            { return s.config.NoPruning }
func (s *Ethereum) BloomIndexer() *core.ChainIndexer             { return s.bloomIndexer }
func (s *Ethereum) SyncMode() downloader.SyncMode {
	mode, _ := s.handler.chainSync.modeAndLocalHead()
	return mode
}

// Protocols returns all the currently configured
// network protocols to start.
func (s *Ethereum) Protocols() []p2p.Protocol {
	protos := eth.MakeProtocols((*ethHandler)(s.handler), s.networkID, s.discmix)
	if s.config.SnapshotCache > 0 {
		protos = append(protos, snap.MakeProtocols((*snapHandler)(s.handler))...)
	}
	return protos
}

// Start implements node.Lifecycle, starting all internal goroutines needed by the
// Ethereum protocol implementation.
func (s *Ethereum) Start() error {
	// let only tendermint bft engine to start its sub modules
	switch s.engine.(type) {
	case *backend.Backend:
		go s.faultDetector.Start()
		go func() {
			header := s.blockchain.CurrentHeader()
			if header.Number.BitLen() == 0 && header.Time > uint64(time.Now().Unix()) {
				s.genesisCountdown()
			}
			s.validatorController()
		}()
	default:
		s.log.Info("running eth backend without Tendermint BFT engine")
	}

	if err := s.setupDiscovery(); err != nil {
		return err
	}
	// Start the bloom bits servicing goroutines
	s.startBloomHandlers(params.BloomBitsBlocks)

	// Regularly update shutdown marker
	s.shutdownTracker.Start()

	// Start the networking layer
	s.handler.Start(s.p2pServer.MaxPeers)
	return nil
}

func (s *Ethereum) setupDiscovery() error {
	eth.StartENRUpdater(s.blockchain, s.p2pServer.LocalNode())

	// Add eth nodes from DNS.
	dnsclient := dnsdisc.NewClient(dnsdisc.Config{})
	if len(s.config.EthDiscoveryURLs) > 0 {
		iter, err := dnsclient.NewIterator(s.config.EthDiscoveryURLs...)
		if err != nil {
			return err
		}
		s.discmix.AddSource(iter)
	}

	// Add snap nodes from DNS.
	if len(s.config.SnapDiscoveryURLs) > 0 {
		iter, err := dnsclient.NewIterator(s.config.SnapDiscoveryURLs...)
		if err != nil {
			return err
		}
		s.discmix.AddSource(iter)
	}

	// Add DHT nodes from discv5.
	if s.p2pServer.DiscoveryV5() != nil {
		filter := eth.NewNodeFilter(s.blockchain)
		iter := enode.Filter(s.p2pServer.DiscoveryV5().RandomNodes(), filter)
		s.discmix.AddSource(iter)
	}

	return nil
}

// This routine is responsible to communicate to devp2p who are the other consensus members
// if the local node is part of the consensus committee or not. It also control the miner start/stop functions.
func (s *Ethereum) validatorController() {
	var (
		chainHeadCh, epochHeadCh   = make(chan core.ChainHeadEvent), make(chan core.EpochHeadEvent)
		chainHeadSub, epochHeadSub = s.blockchain.SubscribeChainHeadEvent(chainHeadCh), s.blockchain.SubscribeEpochHeadEvent(epochHeadCh)
	)

	updateConsensusEnodes := func(header *types.Header) {
		state, err := s.blockchain.StateAt(header.Root)
		if err != nil {
			s.log.Error("Could not retrieve state at head block", "err", err)
			return
		}
		committee, err := s.blockchain.ProtocolContracts().CallGetCommitteeEnodes(state, header, false)
		if err != nil {
			s.log.Error("Could not retrieve consensus whitelist at head block", "err", err)
			return
		}

		index := s.topologySelector.MyIndex(committee.List, s.p2pServer.LocalNode())
		s.p2pServer.UpdateConsensusEnodes(s.topologySelector.RequestSubset(committee.List, index), committee.List)
	}
	wasValidating := false

	// read the committee base on latest state.
	currentHead := s.blockchain.CurrentBlock()
	currentState, err := s.blockchain.StateAt(currentHead.Root)
	if err != nil {
		panic(err)
	}
	epoch, err := s.blockchain.ProtocolContracts().CallEpochByHeight(currentState, currentHead, currentHead.Number)
	if err != nil {
		panic(err)
	}
	committee := epoch.Committee
	if committee.MemberByAddress(s.address) != nil {
		updateConsensusEnodes(currentHead)
		s.miner.Start()
		s.log.Info("Starting node as validator")
		wasValidating = true
	}

	for {
		select {
		case ev := <-chainHeadCh:
			// Get the block number from the event properly
			blockNum := ev.Header.Number.Uint64()
			s.p2pServer.SetCurrentBlockNumber(blockNum)
		case ev := <-epochHeadCh:
			// epoch head change comes with the committee rotation:
			committee = ev.Header.Epoch.Committee
			// check if the local node belongs to the consensus committee.
			if committee.MemberByAddress(s.address) == nil {
				// if the local node was part of the committee set for the previous block
				// there is no longer the need to retain the full connections and the
				// consensus engine enabled.
				if wasValidating {
					s.log.Info("Local node no longer detected part of the consensus committee, mining stopped")
					s.miner.Stop()
					s.p2pServer.UpdateConsensusEnodes(nil, nil)
					wasValidating = false
				}
				continue
			}
			updateConsensusEnodes(ev.Header)
			// if we were not committee in the past block we need to enable the mining engine.
			if !wasValidating {
				s.log.Info("Local node detected part of the consensus committee, mining started")
				s.miner.Start()
			}
			wasValidating = true
		// Err() channel will be closed when unsubscribing.
		case <-chainHeadSub.Err():
			return
		case <-epochHeadSub.Err():
			return
		}
	}
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Ethereum protocol.
func (s *Ethereum) Stop() error {
	// Stop AFD first,
	s.faultDetector.Stop()
	s.engine.Close()
	// Stop all the peer-related stuff then.
	s.discmix.Close()
	s.handler.Stop()
	// Then stop everything else.
	s.bloomIndexer.Close()
	close(s.closeBloomHandler)
	s.txPool.Close()
	s.miner.Close()
	s.blockchain.Stop()

	// Clean shutdown marker as the last thing before closing db
	s.shutdownTracker.Stop()

	s.chainDb.Close()
	s.eventMux.Stop()

	return nil
}

func (s *Ethereum) genesisCountdown() {
	genesisTime := time.Unix(int64(s.blockchain.Genesis().Time()), 0)
	s.log.Info(fmt.Sprintf("Chain genesis time: %v", genesisTime))
	committee := s.blockchain.Genesis().Header().Epoch.Committee
	if committee.MemberByAddress(s.address) != nil {
		s.log.Warn("**************************************************************")
		s.log.Warn("Local node is detected GENESIS VALIDATOR")
		s.log.Warn("Please remain tuned to our Telegram/Discord channels for announcements")
		s.log.Warn("**************************************************************")
	}
	var (
		lastDays   = 0
		lastHours  = 0
		lastMinute = 0
		lastSecond = 0
	)
	for {
		now := time.Now()
		duration := genesisTime.Sub(now)

		if duration <= 0 {
			s.log.Warn("Launch!")
			go func() {
				time.Sleep(3 * time.Second)
				if s.blockchain.Genesis().Number().Cmp(common.Big0) > 0 {
					s.log.Warn("🚀🚀🚀 LAUNCH SUCCESS 🚀🚀🚀")
				}
			}()
			break
		}
		days := int(duration.Hours() / 24)
		hours := int(duration.Hours()) % 24
		minutes := int(duration.Minutes()) % 60
		seconds := int(duration.Seconds()) % 60

		switch {
		case days > 0:
			if days != lastDays {
				lastDays = days
				s.log.Info(fmt.Sprintf("%d day(s) remaining before genesis", days))
			}
		case (hours == 1 || hours == 2 || hours == 6 || hours == 12) && days == 0:
			if hours != lastHours {
				lastHours = hours
				s.log.Info(fmt.Sprintf("%d hour(s) remaining before genesis", hours))
			}
		case (minutes == 1 || minutes == 5 || minutes == 15 || minutes == 30 || minutes == 45) && hours == 0 && days == 0 && seconds == 0:
			if minutes != lastMinute {
				lastMinute = minutes
				s.log.Info(fmt.Sprintf("%d minute(s) remaining before genesis", minutes))
			}
		case (seconds < 10 || seconds == 30 || seconds == 45) && days == 0 && hours == 0 && minutes == 0:
			if seconds != lastSecond {
				lastSecond = seconds
				s.log.Info(fmt.Sprintf("%d second(s) before genesis", seconds))
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (s *Ethereum) Logger() log.Logger {
	return s.log
}

func (eth *Ethereum) StateAtBlock(header *types.Header, reexec uint64, base *state.StateDB, checkLive bool, preferDisk bool) (statedb *state.StateDB, err error) {
	block := eth.blockchain.GetBlockByHash(header.Hash())
	statedb, _, err = eth.stateAtBlock(context.Background(), block, reexec, base, checkLive, preferDisk)
	return
}
