// Code generated by github.com/fjl/gencodec. DO NOT EDIT.

package eth

import (
	"math/big"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/ethash"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/eth/downloader"
	"github.com/clearmatics/autonity/eth/gasprice"
	"github.com/clearmatics/autonity/miner"
	"github.com/clearmatics/autonity/params"
)

// MarshalTOML marshals as TOML.
func (c Config) MarshalTOML() (interface{}, error) {
	type Config struct {
		Genesis                 *core.Genesis `toml:",omitempty"`
		NetworkId               uint64
		SyncMode                downloader.SyncMode
		DiscoveryURLs           []string
		NoPruning               bool
		NoPrefetch              bool
		TxLookupLimit           uint64                 `toml:",omitempty"`
		Whitelist               map[uint64]common.Hash `toml:"-"`
		LightServ               int                    `toml:",omitempty"`
		LightIngress            int                    `toml:",omitempty"`
		LightEgress             int                    `toml:",omitempty"`
		LightPeers              int                    `toml:",omitempty"`
		UltraLightServers       []string               `toml:",omitempty"`
		UltraLightFraction      int                    `toml:",omitempty"`
		UltraLightOnlyAnnounce  bool                   `toml:",omitempty"`
		SkipBcVersionCheck      bool                   `toml:"-"`
		DatabaseHandles         int                    `toml:"-"`
		DatabaseCache           int
		DatabaseFreezer         string
		TrieCleanCache          int
		TrieDirtyCache          int
		TrieTimeout             time.Duration
		Miner                   miner.Config
		Ethash                  ethash.Config
		Tendermint              config.Config
		TxPool                  core.TxPoolConfig
		GPO                     gasprice.Config
		EnablePreimageRecording bool
		DocRoot                 string `toml:"-"`
		EWASMInterpreter        string
		EVMInterpreter          string
		RPCGasCap               *big.Int                       `toml:",omitempty"`
		Checkpoint              *params.TrustedCheckpoint      `toml:",omitempty"`
		CheckpointOracle        *params.CheckpointOracleConfig `toml:",omitempty"`
		OverrideIstanbul        *big.Int                       `toml:",omitempty"`
		OverrideMuirGlacier     *big.Int                       `toml:",omitempty"`
	}
	var enc Config
	enc.Genesis = c.Genesis
	enc.NetworkId = c.NetworkId
	enc.SyncMode = c.SyncMode
	enc.DiscoveryURLs = c.DiscoveryURLs
	enc.NoPruning = c.NoPruning
	enc.NoPrefetch = c.NoPrefetch
	enc.TxLookupLimit = c.TxLookupLimit
	enc.Whitelist = c.Whitelist
	enc.LightServ = c.LightServ
	enc.LightIngress = c.LightIngress
	enc.LightEgress = c.LightEgress
	enc.LightPeers = c.LightPeers
	enc.UltraLightServers = c.UltraLightServers
	enc.UltraLightFraction = c.UltraLightFraction
	enc.UltraLightOnlyAnnounce = c.UltraLightOnlyAnnounce
	enc.SkipBcVersionCheck = c.SkipBcVersionCheck
	enc.DatabaseHandles = c.DatabaseHandles
	enc.DatabaseCache = c.DatabaseCache
	enc.DatabaseFreezer = c.DatabaseFreezer
	enc.TrieCleanCache = c.TrieCleanCache
	enc.TrieDirtyCache = c.TrieDirtyCache
	enc.TrieTimeout = c.TrieTimeout
	enc.Miner = c.Miner
	enc.Ethash = c.Ethash
	enc.Tendermint = c.Tendermint
	enc.TxPool = c.TxPool
	enc.GPO = c.GPO
	enc.EnablePreimageRecording = c.EnablePreimageRecording
	enc.DocRoot = c.DocRoot
	enc.EWASMInterpreter = c.EWASMInterpreter
	enc.EVMInterpreter = c.EVMInterpreter
	enc.RPCGasCap = c.RPCGasCap
	enc.Checkpoint = c.Checkpoint
	enc.CheckpointOracle = c.CheckpointOracle
	return &enc, nil
}

// UnmarshalTOML unmarshals from TOML.
func (c *Config) UnmarshalTOML(unmarshal func(interface{}) error) error {
	type Config struct {
		Genesis                 *core.Genesis `toml:",omitempty"`
		NetworkId               *uint64
		SyncMode                *downloader.SyncMode
		DiscoveryURLs           []string
		NoPruning               *bool
		NoPrefetch              *bool
		TxLookupLimit           *uint64                `toml:",omitempty"`
		Whitelist               map[uint64]common.Hash `toml:"-"`
		LightServ               *int                   `toml:",omitempty"`
		LightIngress            *int                   `toml:",omitempty"`
		LightEgress             *int                   `toml:",omitempty"`
		LightPeers              *int                   `toml:",omitempty"`
		UltraLightServers       []string               `toml:",omitempty"`
		UltraLightFraction      *int                   `toml:",omitempty"`
		UltraLightOnlyAnnounce  *bool                  `toml:",omitempty"`
		SkipBcVersionCheck      *bool                  `toml:"-"`
		DatabaseHandles         *int                   `toml:"-"`
		DatabaseCache           *int
		DatabaseFreezer         *string
		TrieCleanCache          *int
		TrieDirtyCache          *int
		TrieTimeout             *time.Duration
		Miner                   *miner.Config
		Ethash                  *ethash.Config
		Tendermint              *config.Config
		TxPool                  *core.TxPoolConfig
		GPO                     *gasprice.Config
		EnablePreimageRecording *bool
		DocRoot                 *string `toml:"-"`
		EWASMInterpreter        *string
		EVMInterpreter          *string
		RPCGasCap               *big.Int                       `toml:",omitempty"`
		Checkpoint              *params.TrustedCheckpoint      `toml:",omitempty"`
		CheckpointOracle        *params.CheckpointOracleConfig `toml:",omitempty"`
		OverrideIstanbul        *big.Int                       `toml:",omitempty"`
		OverrideMuirGlacier     *big.Int                       `toml:",omitempty"`
	}
	var dec Config
	if err := unmarshal(&dec); err != nil {
		return err
	}
	if dec.Genesis != nil {
		c.Genesis = dec.Genesis
	}
	if dec.NetworkId != nil {
		c.NetworkId = *dec.NetworkId
	}
	if dec.SyncMode != nil {
		c.SyncMode = *dec.SyncMode
	}
	if dec.DiscoveryURLs != nil {
		c.DiscoveryURLs = dec.DiscoveryURLs
	}
	if dec.NoPruning != nil {
		c.NoPruning = *dec.NoPruning
	}
	if dec.NoPrefetch != nil {
		c.NoPrefetch = *dec.NoPrefetch
	}
	if dec.TxLookupLimit != nil {
		c.TxLookupLimit = *dec.TxLookupLimit
	}
	if dec.Whitelist != nil {
		c.Whitelist = dec.Whitelist
	}
	if dec.LightServ != nil {
		c.LightServ = *dec.LightServ
	}
	if dec.LightIngress != nil {
		c.LightIngress = *dec.LightIngress
	}
	if dec.LightEgress != nil {
		c.LightEgress = *dec.LightEgress
	}
	if dec.LightPeers != nil {
		c.LightPeers = *dec.LightPeers
	}
	if dec.UltraLightServers != nil {
		c.UltraLightServers = dec.UltraLightServers
	}
	if dec.UltraLightFraction != nil {
		c.UltraLightFraction = *dec.UltraLightFraction
	}
	if dec.UltraLightOnlyAnnounce != nil {
		c.UltraLightOnlyAnnounce = *dec.UltraLightOnlyAnnounce
	}
	if dec.SkipBcVersionCheck != nil {
		c.SkipBcVersionCheck = *dec.SkipBcVersionCheck
	}
	if dec.DatabaseHandles != nil {
		c.DatabaseHandles = *dec.DatabaseHandles
	}
	if dec.DatabaseCache != nil {
		c.DatabaseCache = *dec.DatabaseCache
	}
	if dec.DatabaseFreezer != nil {
		c.DatabaseFreezer = *dec.DatabaseFreezer
	}
	if dec.TrieCleanCache != nil {
		c.TrieCleanCache = *dec.TrieCleanCache
	}
	if dec.TrieDirtyCache != nil {
		c.TrieDirtyCache = *dec.TrieDirtyCache
	}
	if dec.TrieTimeout != nil {
		c.TrieTimeout = *dec.TrieTimeout
	}
	if dec.Miner != nil {
		c.Miner = *dec.Miner
	}
	if dec.Ethash != nil {
		c.Ethash = *dec.Ethash
	}
	if dec.Tendermint != nil {
		c.Tendermint = *dec.Tendermint
	}
	if dec.TxPool != nil {
		c.TxPool = *dec.TxPool
	}
	if dec.GPO != nil {
		c.GPO = *dec.GPO
	}
	if dec.EnablePreimageRecording != nil {
		c.EnablePreimageRecording = *dec.EnablePreimageRecording
	}
	if dec.DocRoot != nil {
		c.DocRoot = *dec.DocRoot
	}
	if dec.EWASMInterpreter != nil {
		c.EWASMInterpreter = *dec.EWASMInterpreter
	}
	if dec.EVMInterpreter != nil {
		c.EVMInterpreter = *dec.EVMInterpreter
	}
	if dec.RPCGasCap != nil {
		c.RPCGasCap = dec.RPCGasCap
	}
	if dec.Checkpoint != nil {
		c.Checkpoint = dec.Checkpoint
	}
	if dec.CheckpointOracle != nil {
		c.CheckpointOracle = dec.CheckpointOracle
	}
	return nil
}
