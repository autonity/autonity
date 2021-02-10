package tendermint

import (
	"fmt"
	"testing"
	time "time"

	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/params"
	"github.com/davecgh/go-spew/spew"
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
	bridges, err := createBridges(users)
	require.NoError(t, err)
	b := bridges.bridges[0] // Only one bridge
	err = b.Start()
	require.NoError(t, err)
	defer bridges.stop()

	to := time.Millisecond * 100

	proposal, err := b.proposalBlock()
	require.NoError(t, err)
	result := make(chan *types.Block)
	stop := make(chan struct{})
	err = b.Seal(b.blockchain, proposal, result, stop)
	require.NoError(t, err)
	b.pendingMessages(to) // proposal
	b.pendingMessages(to) // prevote
	b.pendingMessages(to) // precommit

	block := b.committedBlock(to, result)

	// Check it is the correct block
	assert.Equal(t, proposal.Hash(), block.Hash())
	// Check it has the right number of committed seals
	assert.Len(t, block.Header().CommittedSeals, 1)
	// Verify the header
	err = b.VerifyHeader(b.blockchain, block.Header(), true)
	assert.NoError(t, err)
}

func TestReachingConsensus(t *testing.T) {
	users, err := Users(4, 1e18, 1, params.UserValidator)
	require.NoError(t, err)
	bridges, err := createBridges(users)
	require.NoError(t, err)
	err = bridges.start()
	require.NoError(t, err)
	defer bridges.stop()

	proposers, err := bridges.proposer()
	require.NoError(t, err)

	proposer := proposers[0] // only one bridge
	proposal, err := proposer.proposalBlock()
	require.NoError(t, err)

	result := make(chan *types.Block)
	stop := make(chan struct{})
	// pass a block to the proposer
	err = proposer.Seal(proposer.blockchain, proposal, result, stop)
	require.NoError(t, err)
	to := time.Millisecond * 100

	// check block not yet committed
	block := proposer.committedBlock(to, result)
	require.Nil(t, block)

	// get the proposal message and validate it
	expectedConsensusMessage := &algorithm.ConsensusMessage{
		MsgType:    algorithm.Propose,
		Height:     proposal.NumberU64(),
		Round:      int64(0),
		ValidRound: int64(-1),
		Value:      algorithm.ValueID(proposal.Hash()),
	}
	proposeMsg := proposer.pendingMessages(to)
	validateProposeMessage(t, proposeMsg, expectedConsensusMessage, proposer, proposal)

	// broadcst the propose message
	err = bridges.broadcast(proposeMsg)
	require.NoError(t, err)
	// time.Sleep(5 * time.Second)

	// check block not yet committed
	block = proposer.committedBlock(to, result)
	require.Nil(t, block)

	// broadcast all prevote messages
	err = bridges.broadcastPendingMessages(to)
	require.NoError(t, err)
	// time.Sleep(5 * time.Second)

	// Validate that the prevotes are as expected
	expectedConsensusMessage = &algorithm.ConsensusMessage{
		MsgType: algorithm.Prevote,
		Height:  proposal.NumberU64(),
		Round:   int64(0),
		Value:   algorithm.ValueID(proposal.Hash()),
	}
	for _, b := range bridges.bridges {
		msg := b.lastSentMessage
		validateMessage(t, msg, expectedConsensusMessage, b)
	}

	// check block not yet committed
	block = proposer.committedBlock(to, result)
	require.Nil(t, block)

	println("precommitsssss")
	// Start brodacsting precommit messages one by one
	expectedConsensusMessage = &algorithm.ConsensusMessage{
		MsgType: algorithm.Precommit,
		Height:  proposal.NumberU64(),
		Round:   int64(0),
		Value:   algorithm.ValueID(proposal.Hash()),
	}
	b := bridges.bridges[0]
	msg := b.pendingMessages(to)
	validateMessage(t, msg, expectedConsensusMessage, b)
	println("precommit 1")
	bridges.broadcast(msg)

	// check block not yet committed
	block = proposer.committedBlock(to, result)
	require.Nil(t, block)

	b = bridges.bridges[1]
	msg = b.pendingMessages(to)
	validateMessage(t, msg, expectedConsensusMessage, b)
	bridges.broadcast(msg)
	println("precommit 2")

	// check block not yet committed
	block = proposer.committedBlock(to, result)
	require.Nil(t, block)

	b = bridges.bridges[2]
	msg = b.pendingMessages(to)
	validateMessage(t, msg, expectedConsensusMessage, b)
	bridges.broadcast(msg)
	println("precommit 3")
	time.Sleep(10 * time.Second)

	// Now we expect the block to be committed, since now a quorum of nodes has
	// broadcast their precommit messages.
	committedBlock := proposer.committedBlock(to, result)
	spew.Dump(committedBlock)
	// Check it is the correct block
	assert.Equal(t, proposal.Hash(), committedBlock.Hash())
	// Check it has the right number of committed seals
	assert.Len(t, committedBlock.Header().CommittedSeals, 3)
	// Verify the header
	err = b.VerifyHeader(b.blockchain, committedBlock.Header(), true)
	assert.NoError(t, err)

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
