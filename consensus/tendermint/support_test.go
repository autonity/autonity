package tendermint

import (
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"
	"net"
	"testing"
	"time"

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
	"github.com/stretchr/testify/require"
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

func (b *broadcasterMock) Broadcast(message []byte) {
}

// notifyingBlockBroadcaster simply passes broadcast messages to a channel which
// we read from in the test.
type notifyingBroadcaster struct {
	messages chan []byte
	closeCh  chan struct{}
}

func (b *notifyingBroadcaster) Broadcast(message []byte) {
	println("received msghash", common.BytesToHash(crypto.Keccak256(message)).String()[2:6])
	select {
	case b.messages <- message:
	case <-b.closeCh:
	}
}

type blockBroadcasterMock struct{}

func (b *blockBroadcasterMock) Enqueue(id string, block *types.Block) {}

// notifyingBlockBroadcaster simply passes broadcast blocks to a channel which
// we read from within the test.
type notifyingBlockBroadcaster struct {
	blocks  chan *types.Block
	closeCh chan struct{}
}

func (b *notifyingBlockBroadcaster) Enqueue(id string, block *types.Block) {
	select {
	case b.blocks <- block:
	case <-b.closeCh:
	}
}

type noActionScheduler struct {
}

func (s *noActionScheduler) ScheduleTimeout(delay uint, f func()) {
	// do nothing
}

// newTestBridge creates a test bridge instance that wraps a bridge and
// provides methods to wait and intercept broadcast messages and broadcast
// blocks as well as utility methods to generate proposal blocks and determine
// if the test bridge is currently the proposer.
func newTestBridge(
	g *core.Genesis,
	user *gengen.User,
	syncer Syncer) (*testBridge, error) {

	messageChan := make(chan []byte)
	closeChan := make(chan struct{})
	blockChan := make(chan *types.Block)
	b, err := createBridge(
		g,
		user,
		syncer,
		&notifyingBroadcaster{messageChan, closeChan},
		&notifyingBlockBroadcaster{blockChan, closeChan},
		&noActionScheduler{},
	)
	if err != nil {
		return nil, err
	}
	genesisBlock, err := b.latestBlockRetriever.LatestBlock()
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve genesis block: %v", err)
	}
	return &testBridge{
		Bridge:             b,
		messageChan:        messageChan,
		blockChan:          blockChan,
		lastCommittedBlock: genesisBlock,
		closeCh:            closeChan,
	}, nil
}

// createBridge creates a fully running bridge with the exception of the
// syncer, broadcaster and blockBroadcaster provided. For these three
// components we provide test implementations that allow us to intercept the
// sync messages, broadcast messages and broadcast blocks.
func createBridge(
	g *core.Genesis,
	user *gengen.User,
	syncer Syncer,
	broadcaster Broadcaster,
	blockBroadcaster consensus.Broadcaster,
	timeoutScheduler TimeoutScheduler,
) (*Bridge, error) {
	db := rawdb.NewMemoryDatabase()
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

	b := New(
		g.Config.Tendermint,
		user.Key.(*ecdsa.PrivateKey),
		broadcaster,
		syncer,
		verifier,
		finalizer,
		latestBlockRetriever,
		autonityContract,
		timeoutScheduler,
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

func createBridges(users []*gengen.User) (*testBridges, error) {
	g, err := gengen.NewGenesis(1, users)
	if err != nil {
		return nil, err
	}
	bridges := make([]*testBridge, len(users))
	bridgeMap := make(map[common.Address]*testBridge, len(users))
	for i, u := range users {
		b, err := newTestBridge(g, u, &syncerMock{})
		if err != nil {
			return nil, err
		}
		bridges[i] = b
		bridgeMap[b.address] = b
	}
	return &testBridges{bridges, bridgeMap}, nil
}

// testBridges provides a way to manage a group of bridges and provides
// convenience, functions to manipulate the group.
type testBridges struct {
	bridges   []*testBridge
	bridgeMap map[common.Address]*testBridge
}

func (b *testBridges) byAddress(addr common.Address) *testBridge {
	return b.bridgeMap[addr]
}

// proposer gets the proposer for each bridge, they may not be the same if the
// bridges are not synced. The first error to be encountered is returned.
func (b *testBridges) proposer() ([]*testBridge, error) {
	proposers := make([]*testBridge, len(b.bridges))
	for i, bridge := range b.bridges {
		addr, err := bridge.proposer()
		if err != nil {
			return nil, err
		}
		proposers[i] = b.bridgeMap[addr]
	}
	return proposers, nil
}

func (b *testBridges) start() error {
	for _, bridge := range b.bridges {
		err := bridge.Start()
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *testBridges) stop() error {
	for _, bridge := range b.bridges {
		err := bridge.stop()
		if err != nil {
			return err
		}
	}
	return nil
}

// // This assumes that b.sentMessage has not been called on any bridges and that
// // a proposal block has been passed to the Seal function of the proposer for
// // this height and round. No block is returned but the blocks should be
// // available in lastCommittedBlock on each bridge.
// func (b *testBridges) awaitBlock(sealChan chan *types.Block) error {
// 	proposers, err := b.proposer()
// 	if err != nil {
// 		return err
// 	}
// 	p := proposers[0]
// 	to := time.Millisecond * 100
// 	m := p.pendingMessage(to) // get the proposal message
// 	err = b.broadcast(m)      // send it to everyone else
// 	if err != nil {
// 		return err
// 	}
// 	// Now send the prevotes
// 	for _, bridge := range b.bridges {
// 		m = bridge.pendingMessage(to)
// 		err := b.broadcast(m)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	// Now send the precommits
// 	for _, bridge := range b.bridges {
// 		m = bridge.pendingMessage(to)
// 		err := b.broadcast(m)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	// now get the blocks
// 	for _, bridge := range b.bridges {
// 		var sc chan *types.Block
// 		if bridge.address == p.address {
// 			sc = sealChan
// 		}
// 		block := bridge.committedBlock(to, sc)
// 		if err != nil {
// 			return err
// 		}
// 		if block.Hash() != m.value.Hash() {
// 			return fmt.Errorf("unexpected block, expecting: %v, got: %v", m.value.Hash().String(), block.Hash().String())
// 		}
// 	}
// 	return nil
// }

func (b *testBridges) broadcast(m *message) error {
	println("broadcasting", m.consensusMessage.String())
	encoded, err := encodeSignedMessage(m.consensusMessage, b.byAddress(m.address).key, m.value)
	if err != nil {
		return err
	}
	for _, b := range b.bridges {
		if b.address == m.address {
			continue
		}
		size, reader, err := rlp.EncodeToReader(encoded)
		if err != nil {
			return err
		}
		msg := p2p.Msg{
			Code:    tendermintMsg,
			Payload: reader,
			Size:    uint32(size),
		}

		_, err = b.HandleMsg(m.address, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

// broadcastPendingMessages calls pendingMessage for each bridge and forwards
// the returned message to all other bridges. The given timeout is the time to
// wait per bridge for the result from pendingMessage.
func (b *testBridges) broadcastPendingMessages(timeout time.Duration) error {
	// Now send the prevotes
	for _, bridge := range b.bridges {
		m := bridge.pendingMessage(timeout)
		println("broadcasting", m.consensusMessage.String())
		err := b.broadcast(m)
		if err != nil {
			return err
		}
	}
	return nil
}

// testBridge wraps a Bridge instance and provides
type testBridge struct {
	*Bridge
	messageChan        chan []byte
	lastSentMessage    *message
	blockChan          chan *types.Block
	lastCommittedBlock *types.Block
	closeCh            chan struct{}
}

// This closes the test bridge, it is permanent calling this twice will panic.
// Bridges will likely not shut down properly with a direct call to Close since
// they will probably be stuck sending a message on the notifyingBroadcaster or
// notifyingBlockBroadcaster, this method closes the closeCh which releases
// goroutines stuck sending a message on the notifyingBroadcaster or
// notifyingBlockBroadcaster thereby allowing them to close.
func (b *testBridge) stop() error {
	close(b.closeCh)
	return b.Close()
}

// Retrieves the messages from this bridge that have been broadcast by this
// bridge one at a time and in the order they were broadcast. If no message is
// broadcast before the timeout expires then nil is returned.
func (b *testBridge) pendingMessage(timeout time.Duration) *message {
	t := time.NewTimer(timeout)
	for {
		select {
		case m := <-b.messageChan:
			message, err := decodeSignedMessage(m)
			if err != nil {
				panic(fmt.Sprintf("failed to decode signed message: %v", err))
			}
			if message.address != b.address {
				// ignore rebroadcast messages
				continue
			}
			b.lastSentMessage = message
			println("gotmessage", message.consensusMessage.String())
			return message
		case <-t.C:
			println("nomessage")
			return nil
		}
	}
}

// keeps requesting messages from the pending messages untill none are
// returned. This is required to free up the routine from the bridge that might
// be stuck trying to send on the messageChan, so that we can close the bridge.
func (b *testBridge) drainPendingMessages(timeout time.Duration) {
	msg := &message{}
	for msg != nil {
		msg = b.pendingMessage(timeout)
	}
}

func (b *testBridge) proposer() (common.Address, error) {
	var round int64
	if b.lastSentMessage != nil {
		round = b.lastSentMessage.consensusMessage.Round
	}
	return b.proposerAddr(b.lastCommittedBlock.Header(), round)
}

// Waits for the committed block,if sealChan is provided assumes that this
// bridge is the proposer and will wait on that channel, if not it will wait on
// the blockChan.
func (b *testBridge) committedBlock(timeout time.Duration, sealChan chan *types.Block) *types.Block {
	t := time.NewTimer(timeout)
	if sealChan != nil {
		select {
		case block := <-sealChan:
			b.lastCommittedBlock = block
			return block
		case <-t.C:
			return nil
		}
	} else {
		select {
		case block := <-b.blockChan:
			b.lastCommittedBlock = block
			return block
		case <-t.C:
			return nil
		}
	}
}

func (b *testBridge) proposalBlock() (*types.Block, error) {
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

func validateMessage(t *testing.T, msg *message, expectedConsensusMessage *algorithm.ConsensusMessage, b *testBridge) {
	require.Equal(t, expectedConsensusMessage, msg.consensusMessage)
	require.Equal(t, b.address, msg.address)
}

func validateProposeMessage(t *testing.T, proposeMsg *message, expectedConsensusMessage *algorithm.ConsensusMessage, proposer *testBridge, proposal *types.Block) {
	validateMessage(t, proposeMsg, expectedConsensusMessage, proposer)
	// Due to the way that blocks are constructed they can be conceptually
	// equal even if they are not equal in the point of view of the go
	// language. So we just check the hash here.
	require.Equal(t, proposal.Hash(), proposeMsg.value.Hash())
	expectedProposerSeal, err := crypto.Sign(proposal.Hash().Bytes(), proposer.key)
	require.NoError(t, err)
	require.Equal(t, expectedProposerSeal, proposeMsg.proposerSeal)
}
