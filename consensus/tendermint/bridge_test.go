package tendermint

import (
	"testing"
	time "time"

	"github.com/clearmatics/autonity/cmd/gengen/gengen"
	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This test just stats and stops a bridge instance, making sure an error is
// returned if the bridge is started or stopped twice, and no error is returned
// otherwise.
func TestStartingAndStoppingBridge(t *testing.T) {
	users, err := Users(1, 1e18, 1, params.UserValidator)
	require.NoError(t, err)
	g, err := gengen.NewGenesis(1, users)
	require.NoError(t, err)
	b, err := createBridge(g, users[0], &syncerMock{}, &broadcasterMock{}, &blockBroadcasterMock{}, &noActionScheduler{})
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

// This test checks that if a freshly started bridge instance, who holds all
// the stake in the system is provided with a block to seal, it will progress
// to commit that block.
func TestBlockGivenToSealIsComitted(t *testing.T) {
	users, err := Users(1, 1e18, 1, params.UserValidator)
	require.NoError(t, err)
	bridges, err := createBridges(users)
	require.NoError(t, err)
	b := bridges.bridges[0] // Only one bridge
	err = b.Start()
	require.NoError(t, err)
	defer bridges.stop() // nolint

	to := time.Millisecond * 100

	proposal, err := b.proposalBlock()
	require.NoError(t, err)
	result := make(chan *types.Block)
	stop := make(chan struct{})
	err = b.Seal(b.blockchain, proposal, result, stop)
	require.NoError(t, err)
	b.pendingMessage(to) // proposal
	b.pendingMessage(to) // prevote
	b.pendingMessage(to) // precommit

	block := b.committedBlock(to, result)

	// Check it is the correct block
	assert.Equal(t, proposal.Hash(), block.Hash())
	// Check it has the right number of committed seals
	assert.Len(t, block.Header().CommittedSeals, 1)
	// Verify the header
	err = b.VerifyHeader(b.blockchain, block.Header(), true)
	assert.NoError(t, err)
}

// This test checks that a call to NewChainHead is required to progress to the
// next block.  It uses a single bridge that has all the stake and lets it
// commit one block, and shows that the bridge instance will then wait for a
// call to NewChainHead before beginning work on the next block.
func TestNewChainHead(t *testing.T) {
	users, err := Users(1, 1e18, 1, params.UserValidator)
	require.NoError(t, err)
	bridges, err := createBridges(users)
	require.NoError(t, err)
	b := bridges.bridges[0] // Only one bridge
	err = b.Start()
	require.NoError(t, err)
	defer bridges.stop() // nolint

	to := time.Millisecond * 50

	proposal, err := b.proposalBlock()
	require.NoError(t, err)
	result := make(chan *types.Block)
	stop := make(chan struct{})
	err = b.Seal(b.blockchain, proposal, result, stop)
	require.NoError(t, err)
	b.pendingMessage(to) // proposal
	b.pendingMessage(to) // prevote
	b.pendingMessage(to) // precommit

	block := b.committedBlock(to, result)

	// Check it is the correct block
	assert.Equal(t, proposal.Hash(), block.Hash())

	// Now we will pass another block to seal, but we don't expect to see a
	// proposal message until NewChainHead is called.
	proposal, err = b.proposalBlock()
	require.NoError(t, err)
	err = b.Seal(b.blockchain, proposal, result, stop)
	require.NoError(t, err)

	// Expect nil message
	proposeMessage := b.pendingMessage(to)
	require.Nil(t, proposeMessage)

	// Now expect the new propose mesage
	err = b.NewChainHead()
	require.NoError(t, err)

	expectedConsensusMessage := &algorithm.ConsensusMessage{
		MsgType:    algorithm.Propose,
		Height:     proposal.NumberU64(),
		Round:      int64(0),
		ValidRound: int64(-1),
		Value:      algorithm.ValueID(proposal.Hash()),
	}
	// Expect new propose message.
	proposeMessage = b.pendingMessage(to)
	require.Equal(t, expectedConsensusMessage, proposeMessage.consensusMessage)

}

// This test constructs a group of 4 bridges, calls Seal with a proposal block
// on the proposer and validates the progression of the bridges through the
// stages of the tendermint algorithm to the point where the block is
// committed. Note that since bridge instances are not connected we are free to
// control when messages are sent in the test.
func TestReachingConsensus(t *testing.T) {
	users, err := Users(4, 1e18, 1, params.UserValidator)
	require.NoError(t, err)
	bridges, err := createBridges(users)
	require.NoError(t, err)
	err = bridges.start()
	require.NoError(t, err)
	defer bridges.stop() // nolint

	proposers, err := bridges.proposer()
	require.NoError(t, err)

	proposer := proposers[0]
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
	proposeMsg := proposer.pendingMessage(to)
	validateProposeMessage(t, proposeMsg, expectedConsensusMessage, proposer, proposal)

	// broadcst the propose message
	err = bridges.broadcast(proposeMsg)
	require.NoError(t, err)

	// check block not yet committed
	block = proposer.committedBlock(to, result)
	require.Nil(t, block)

	// broadcast all prevote messages
	err = bridges.broadcastPendingMessages(to)
	require.NoError(t, err)

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

	// Start brodacsting precommit messages one by one, at this point all the
	// bridges will have handled their own precommit message, so it will only
	// take 2 more to bring the network to agreement on the block.
	expectedConsensusMessage = &algorithm.ConsensusMessage{
		MsgType: algorithm.Precommit,
		Height:  proposal.NumberU64(),
		Round:   int64(0),
		Value:   algorithm.ValueID(proposal.Hash()),
	}

	b := bridges.bridges[0]
	msg := b.pendingMessage(to)
	validateMessage(t, msg, expectedConsensusMessage, b)
	err = bridges.broadcast(msg)
	require.NoError(t, err)

	// check block not yet committed
	block = proposer.committedBlock(to, result)
	require.Nil(t, block)

	b = bridges.bridges[1]
	msg = b.pendingMessage(to)
	validateMessage(t, msg, expectedConsensusMessage, b)
	err = bridges.broadcast(msg)
	require.NoError(t, err)

	// check block not yet committed
	block = proposer.committedBlock(to, result)
	require.Nil(t, block)

	b = bridges.bridges[2]
	msg = b.pendingMessage(to)
	validateMessage(t, msg, expectedConsensusMessage, b)
	err = bridges.broadcast(msg)
	require.NoError(t, err)

	// Now we expect the block to be committed, since 3 of 4 nodes has
	// broadcast their precommit messages and each of those nodes will have
	// processed their own precommit message giving us 3 of 4 commit messages
	// in the nodes that have broadcast their message.

	// We need to check if b is the proposer so that we can correctly pass
	// sealChan to committedBlock.
	var sealChan chan *types.Block
	if b.address == proposer.address {
		sealChan = result
	}
	committedBlock := b.committedBlock(to, sealChan)
	// Check it is the correct block
	assert.Equal(t, proposal.Hash(), committedBlock.Hash())
	// Check it has the right number of committed seals
	assert.Len(t, committedBlock.Header().CommittedSeals, 3)
	// Verify the header
	err = b.VerifyHeader(b.blockchain, committedBlock.Header(), true)
	assert.NoError(t, err)
}
