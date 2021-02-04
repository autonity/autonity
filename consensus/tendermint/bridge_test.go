package tendermint

import (
	"crypto/ecdsa"
	"io/ioutil"
	"math"
	"math/big"
	"net"
	"os"
	"testing"
	time "time"

	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/cmd/gengen/gengen"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/rawdb"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/node"
	"github.com/clearmatics/autonity/p2p"
	"github.com/clearmatics/autonity/params"
	"github.com/stretchr/testify/assert"
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

	baseTendermintConfig = config.Config{
		BlockPeriod: 0,
	}
)

func Genesis(key *ecdsa.PrivateKey) (*core.Genesis, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	u := &gengen.User{
		InitialEth: big.NewInt(10 ^ 18), // 1 E
		Key:        key,
		KeyPath:    f.Name(),
		NodeIP:     net.ParseIP("0.0.0.0"),
		NodePort:   0,
		Stake:      1,
		UserType:   params.UserValidator,
	}
	return gengen.NewGenesis(1, []*gengen.User{u})
}

type syncerMock struct{}

func (s *syncerMock) Start()                                              {}
func (s *syncerMock) Stop()                                               {}
func (s *syncerMock) AskSync(lastestHeader *types.Header)                 {}
func (s *syncerMock) SyncPeer(peerAddr common.Address, messages [][]byte) {}

type broadcasterMock struct{}

func (b *broadcasterMock) Broadcast(message []byte) {}

type blockBroadcasterMock struct{}

func (b *blockBroadcasterMock) Enqueue(id string, block *types.Block) {}

// CreateConsensusEngine creates the required type of consensus engine instance for an Ethereum service
func createBridge(
	config *config.Config,
	syncer Syncer,
	broadcaster Broadcaster,
	blockBroadcaster consensus.Broadcaster,
) (*Bridge, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	g, err := Genesis(key)
	if err != nil {
		return nil, err
	}
	db := rawdb.NewMemoryDatabase()
	if err != nil {
		return nil, err
	}
	chainConfig, _, err := core.SetupGenesisBlock(db, g)
	if err != nil {
		return nil, err
	}
	hg, err := core.NewHeaderGetter(db)
	if err != nil {
		return nil, err
	}
	vmConfig := &vm.Config{}
	evmP := core.NewDefaultEVMProvider(hg, *vmConfig, chainConfig)
	autonityContract, err := autonity.NewAutonityContractFromConfig(db, hg, evmP, chainConfig.AutonityContractConfig)
	if err != nil {
		return nil, err
	}
	finalizer := NewFinalizer(autonityContract)
	verifier := NewVerifier(vmConfig, finalizer, config.BlockPeriod)
	// broadcaster := NewBroadcaster(crypto.PubkeyToAddress(key.PublicKey), peers)
	statedb := state.NewDatabase(db)
	latestBlockRetriever := NewBlockReader(db, statedb)
	b := New(
		config,
		key,
		broadcaster,
		syncer,
		verifier,
		finalizer,
		latestBlockRetriever,
		autonityContract,
	)
	isLocalBlock := func(block *types.Block) bool {
		return true
	}
	var txLookupLimit uint64 = 0
	bc, err := core.NewBlockChainWithState(db, statedb, nil, chainConfig, b, *vmConfig, isLocalBlock, core.NewTxSenderCacher(1), &txLookupLimit, hg, autonityContract)
	if err != nil {
		return nil, err
	}
	b.SetExtraComponents(bc, blockBroadcaster)
	return b, nil
}

func TestStartingAndStoppingBridge(t *testing.T) {
	b, err := createBridge(config.DefaultConfig(), &syncerMock{}, &broadcasterMock{}, &blockBroadcasterMock{})
	require.NoError(t, err)
	err = b.Start()
	require.NoError(t, err)
	err = b.Start()
	require.Error(t, err)
	err = b.Close()
	require.NoError(t, err)
	err = b.Close()
	require.Error(t, err)
}

func TestBlockGivenToSealIsComitted(t *testing.T) {
	b, err := createBridge(config.DefaultConfig(), &syncerMock{}, &broadcasterMock{}, &blockBroadcasterMock{})
	require.NoError(t, err)
	err = b.Start()
	require.NoError(t, err)

	block, err := b.latestBlockRetriever.LatestBlock()
	require.NoError(t, err)
	state, err := b.blockchain.State()
	require.NoError(t, err)
	var receipts []*types.Receipt

	header := &types.Header{
		ParentHash: block.Hash(),
		Number:     new(big.Int).Add(block.Number(), common.Big1),
		GasLimit:   math.MaxUint64,
	}
	b.Prepare(b.blockchain, header)
	newBlock, err := b.FinalizeAndAssemble(b.blockchain, header, state, nil, nil, &receipts)

	require.NoError(t, err)
	result := make(chan *types.Block)
	stop := make(chan struct{})
	err = b.Seal(b.blockchain, newBlock, result, stop)
	require.NoError(t, err)
	tm := time.NewTimer(time.Millisecond * 100)
	select {
	case <-tm.C:
		t.Fatalf("Expecting block to have been committed")
	case r := <-result:
		assert.Equal(t, newBlock.Hash(), r.Hash())
	}
}
