package tendermint

import (
	"crypto/ecdsa"

	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/ethdb"
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
