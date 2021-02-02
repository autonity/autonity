package tendermint

import (
	"crypto/ecdsa"
	"testing"
	time "time"

	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/eth"
	"github.com/clearmatics/autonity/eth/downloader"
	"github.com/clearmatics/autonity/ethdb"
	"github.com/clearmatics/autonity/node"
	"github.com/clearmatics/autonity/p2p"
	"github.com/clearmatics/autonity/params"
	"github.com/stretchr/testify/require"
)

var (
	baseNodeConfig *node.Config = &node.Config{
		Name:    "autonity",
		Version: params.Version,
		P2P: p2p.Config{
			MaxPeers:              100,
			DialHistoryExpiration: time.Millisecond,
		},
		NoUSB:    true,
		HTTPHost: "0.0.0.0",
		WSHost:   "0.0.0.0",
	}

	baseEthConfig = &eth.Config{
		SyncMode:        downloader.FullSync,
		DatabaseCache:   256,
		DatabaseHandles: 256,
		TxPool:          core.DefaultTxPoolConfig,
	}
	baseTendermintConfig = config.Config{
		BlockPeriod: 0,
	}
)

// CreateConsensusEngine creates the required type of consensus engine instance for an Ethereum service
func createBridge(
	key *ecdsa.PrivateKey,
	config *config.Config,
	db ethdb.Database,
	vmConfig *vm.Config,
	peers consensus.Peers,
	state state.Database,
	autonityContract *autonity.Contract,
) consensus.Engine {
	finalizer := NewFinalizer(autonityContract)
	verifier := NewVerifier(vmConfig, finalizer, config.BlockPeriod)
	syncer := NewSyncer(peers)
	broadcaster := NewBroadcaster(crypto.PubkeyToAddress(key.PublicKey), peers)
	latestBlockRetriever := NewBlockReader(db, state)
	return New(
		config,
		key,
		broadcaster,
		syncer,
		verifier,
		finalizer,
		latestBlockRetriever,
		autonityContract,
		state,
	)
}

func TestRun(t *testing.T) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err, )
	b := createBridge(key, &baseTendermintConfig, db ethdb.Database, vmConfig *vm.Config, peers consensus.Peers, state state.Database, autonityContract *autonity.Contract)
}
