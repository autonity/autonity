package tendermint

import (
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"
	"net"
	"testing"
	time "time"

	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/cmd/gengen/gengen"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/rawdb"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/p2p"
	"github.com/clearmatics/autonity/params"
	"github.com/clearmatics/autonity/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
// baseNodeConfig *node.Config = &node.Config{
// 	Name:    "autonity",
// 	Version: params.Version,
// 	P2P: p2p.Config{
// 		MaxPeers:              100,
// 		DialHistoryExpiration: time.Millisecond,
// 	},
// 	NoUSB:    true,
// 	HTTPHost: "0.0.0.0",
// 	WSHost:   "0.0.0.0",
// }

// baseTendermintConfig = config.Config{
// 	BlockPeriod: 0,
// }
)

func Users(count int, e, stake uint64, usertype params.UserType) ([]*gengen.User, error) {
	users := make([]*gengen.User, count)
	for i := range users {
		key, err := crypto.GenerateKey()
		if err != nil {
			return nil, err
		}
		users[i] = &gengen.User{
			InitialEth: new(big.Int).SetUint64(e),
			Key:        key,
			//We use the empty string here since the key will not be persisted.
			KeyPath: "",
			// We use the zero address here because we won't actualls make or
			// receive any connections.
			NodeIP:   net.ParseIP("0.0.0.0"),
			NodePort: 0,
			Stake:    stake,
			UserType: usertype,
		}
	}
	return users, nil
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

// createBridge creates a fully working bridge, the instance has no missing
// fields or fake fields, except for the syncer, brodcaster and
// blockBroadcaster parameters which are under the caller's control. The
// returned bridge will be the bridge for the user from users who will be the
// proposer for the next block.
func createBridge(
	users []*gengen.User,
	syncer Syncer,
	broadcaster Broadcaster,
	blockBroadcaster consensus.Broadcaster,
) (*Bridge, error) {
	g, err := gengen.NewGenesis(1, users)
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
	config := g.Config.Tendermint
	finalizer := NewFinalizer(autonityContract)
	verifier := NewVerifier(vmConfig, finalizer, config.BlockPeriod)
	statedb := state.NewDatabase(db)
	latestBlockRetriever := NewBlockReader(db, statedb)
	genesisBlock, err := latestBlockRetriever.LatestBlock()
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve genesis block: %v", err)
	}
	state, err := latestBlockRetriever.BlockState(genesisBlock.Root())
	if err != nil {
		return nil, fmt.Errorf("cannot load state from block chain: %v", err)
	}
	// Get initial proposer
	proposer, err := autonityContract.GetProposerFromAC(genesisBlock.Header(), state, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get initial proposer: %v", err)
	}
	var proposerKey *ecdsa.PrivateKey
	for _, u := range users {
		k := u.Key.(*ecdsa.PrivateKey)
		if crypto.PubkeyToAddress(k.PublicKey) == proposer {
			proposerKey = k
		}
	}
	// Construct bridge with initial proposer
	b := New(
		g.Config.Tendermint,
		proposerKey,
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

func getNextProposalBlock(b *Bridge) (*types.Block, error) {
	block, err := b.latestBlockRetriever.LatestBlock()
	if err != nil {
		return nil, err
	}
	state, err := b.blockchain.State()
	if err != nil {
		return nil, err
	}
	var receipts []*types.Receipt
	header := &types.Header{
		ParentHash: block.Hash(),
		Number:     new(big.Int).Add(block.Number(), common.Big1),
		GasLimit:   math.MaxUint64,
	}
	err = b.Prepare(b.blockchain, header)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return b.FinalizeAndAssemble(b.blockchain, header, state, nil, nil, &receipts)
}

func TestStartingAndStoppingBridge(t *testing.T) {
	users, err := Users(1, 1e18, 1, params.UserValidator)
	require.NoError(t, err)
	b, err := createBridge(users, &syncerMock{}, &broadcasterMock{}, &blockBroadcasterMock{})
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
	users, err := Users(1, 1e18, 1, params.UserValidator)
	require.NoError(t, err)
	b, err := createBridge(users, &syncerMock{}, &broadcasterMock{}, &blockBroadcasterMock{})
	require.NoError(t, err)
	err = b.Start()
	require.NoError(t, err)

	proposal, err := getNextProposalBlock(b)
	require.NoError(t, err)
	result := make(chan *types.Block)
	stop := make(chan struct{})
	err = b.Seal(b.blockchain, proposal, result, stop)
	require.NoError(t, err)
	tm := time.NewTimer(time.Millisecond * 100)
	select {
	case <-tm.C:
		t.Fatalf("Expecting block to have been committed")
	case r := <-result:
		// Check it is the correct block
		assert.Equal(t, proposal.Hash(), r.Hash())
		// Check it has the right number of committed seals
		assert.Len(t, r.Header().CommittedSeals, 1)
		// Verify the header
		err := b.VerifyHeader(b.blockchain, r.Header(), true)
		assert.NoError(t, err)
	}

}

func TestReachingConsensus(t *testing.T) {
	users, err := Users(4, 1e18, 1, params.UserValidator)
	require.NoError(t, err)
	b, err := createBridge(users, &syncerMock{}, &broadcasterMock{}, &blockBroadcasterMock{})
	require.NoError(t, err)
	err = b.Start()
	require.NoError(t, err)

	var nonProposers []*gengen.User
	for _, u := range users {
		if u.Key.(*ecdsa.PrivateKey) != b.key {
			nonProposers = append(nonProposers, u)
		}
	}

	proposal, err := getNextProposalBlock(b)
	require.NoError(t, err)
	result := make(chan *types.Block)
	stop := make(chan struct{})

	// pass a block to the proposer
	err = b.Seal(b.blockchain, proposal, result, stop)
	require.NoError(t, err)
	tm := time.NewTimer(time.Millisecond * 100)
	// Do not expect block to be committed, we have not reached quorum
	select {
	case <-tm.C:
	case <-result:
		t.Fatalf("Not expecting block to have been committed")
	}
	prevote := &algorithm.ConsensusMessage{
		Height:  proposal.NumberU64(),
		Round:   0,
		MsgType: algorithm.Prevote,
		Value:   algorithm.ValueID(proposal.Hash()),
	}
	// Send prevotes
	for _, u := range nonProposers {
		handled, err := sendMessage(prevote, u, b)
		require.NoError(t, err)
		require.True(t, handled)
	}

	tm = time.NewTimer(time.Millisecond * 100)
	// Do not expect block to be committed, we have not reached quorum
	select {
	case <-tm.C:
	case <-result:
		t.Fatalf("Not expecting block to have been committed")
	}

	precommit := &algorithm.ConsensusMessage{
		Height:  proposal.NumberU64(),
		Round:   0,
		MsgType: algorithm.Precommit,
		Value:   algorithm.ValueID(proposal.Hash()),
	}
	handled, err := sendMessage(precommit, nonProposers[0], b)
	require.NoError(t, err)
	require.True(t, handled)

	// Do not expect block to be committed, we have not reached quorum
	tm = time.NewTimer(time.Millisecond * 100)
	select {
	case <-tm.C:
	case <-result:
		t.Fatalf("Not expecting block to have been committed")
	}

	handled, err = sendMessage(precommit, nonProposers[1], b)
	require.NoError(t, err)
	require.True(t, handled)

	// Expect block to be committed, we should have reached quorum
	select {
	case <-tm.C:
		t.Fatalf("Expecting block to have been committed")
	case r := <-result:
		// Check it is the correct block
		assert.Equal(t, proposal.Hash(), r.Hash())
		// Check it has the right number of committed seals
		assert.Len(t, r.Header().CommittedSeals, 3)
		// Verify the header
		err := b.VerifyHeader(b.blockchain, r.Header(), true)
		assert.NoError(t, err)
	}

}

func sendMessage(m *algorithm.ConsensusMessage, u *gengen.User, b *Bridge) (bool, error) {
	k := u.Key.(*ecdsa.PrivateKey)
	encoded, err := encodeSignedMessage(m, k, nil)
	if err != nil {
		return false, err
	}
	size, reader, err := rlp.EncodeToReader(encoded)
	if err != nil {
		return false, err
	}
	msg := p2p.Msg{
		Code:    tendermintMsg,
		Payload: reader,
		Size:    uint32(size),
	}
	return b.HandleMsg(crypto.PubkeyToAddress(k.PublicKey), msg)
}

// This test shows that GetSignatureAddressHash does not verify the signature.
func TestSignAndVerify(t *testing.T) {
	t.Skip("Skipped because this sometimes fails")
	h := crypto.Keccak256Hash([]byte{})
	k, err := crypto.GenerateKey()
	require.NoError(t, err)
	sig, err := crypto.Sign(h[:], k)
	require.NoError(t, err)
	sig[0] = 1
	addr, err := types.GetSignatureAddressHash(h[:], sig)
	fmt.Printf("addr: %v error: %v\n", addr.String(), err)
	fmt.Printf("orig: %v\n", crypto.PubkeyToAddress(k.PublicKey).String())
	require.Error(t, err)
}
