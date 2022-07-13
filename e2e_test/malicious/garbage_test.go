package malicious

import (
	"context"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"github.com/autonity/autonity/core/types"
	e2etest "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/node"
	"github.com/autonity/autonity/test"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"testing"
)

type randomBytesBroadcaster struct {
	*core.Core
}

func (s *randomBytesBroadcaster) Broadcast(ctx context.Context, msg *messageutils.Message) {
	logger := s.Logger().New("step", s.Step())
	logger.Info("Broadcasting random bytes")

	for i := 0; i < 1000; i++ {
		payload, err := e2etest.GenerateRandomBytes(2048)
		if err != nil {
			logger.Error("Failed to generate random bytes ", "err", err)
			return
		}
		if err = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), payload); err != nil {
			logger.Error("Failed to broadcast message", "msg", msg, "err", err)
			return
		}
	}
}

// TestRandomBytesBroadcaster broadcasts random bytes in the network,
// We expect other nodes to detect this misbehaviour and discard these messages
// Receiving nodes should also disconnect misbehaving nodes
func TestRandomBytesBroadcaster(t *testing.T) {
	users, err := test.Users(6, "10e18,v,100,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)

	//set Malicious users
	users[0].CustHandler = &node.CustomHandler{Broadcaster: &randomBytesBroadcaster{}}
	// creates a network of 6 users and starts all the nodes in it
	network, err := test.NewNetworkFromUsers(users, true)
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 180)
	require.NoError(t, err)
	defer network.Shutdown()
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

type garbageMessageBroadcaster struct {
	*core.Core
}

func (s *garbageMessageBroadcaster) Broadcast(ctx context.Context, msg *messageutils.Message) {
	logger := s.Logger().New("step", s.Step())

	var fMsg messageutils.Message
	err := gofakeit.Struct(&fMsg)
	if err != nil {
		s.Logger().Error("Failed to fake proposal struct, err ", err)
		return
	}
	logger.Info("Broadcasting random bytes")

	payload, err := s.FinalizeMessage(&fMsg)
	if err != nil {
		logger.Error("Failed to finalize message", "msg", fMsg, "err", err)
		return
	}
	if err = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), payload); err != nil {
		logger.Error("Failed to broadcast message", "msg", fMsg, "err", err)
		return
	}
}

// TestGarbageMessageBroadcaster broadcasts a garbage Messages in the network,
// We expect other nodes to detect this misbehaviour and discard these messages
// Receiving nodes should also disconnect misbehaving nodes
func TestGarbageMessageBroadcaster(t *testing.T) {
	users, err := test.Users(6, "10e18,v,100,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)

	//set Malicious users
	users[0].CustHandler = &node.CustomHandler{Broadcaster: &garbageMessageBroadcaster{}}
	// creates a network of 6 users and starts all the nodes in it
	network, err := test.NewNetworkFromUsers(users, true)
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 180)
	require.NoError(t, err)
	defer network.Shutdown()
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

type garbagePrecommitSender struct {
	*core.Core
	interfaces.Precommiter
}

func (c *garbagePrecommitSender) SendPrecommit(ctx context.Context, isNil bool) {
	logger := c.Logger().New("step", c.Step())

	var precommitMsg messageutils.Vote
	err := gofakeit.Struct(&precommitMsg)
	if err != nil {
		logger.Error("Failed to fake precommit struct, err ", err)
		return
	}

	encodedVote, err := messageutils.Encode(&precommitMsg)
	if err != nil {
		logger.Error("Failed to encode", "subject", precommitMsg)
		return
	}
	msg := &messageutils.Message{
		Code:          messageutils.MsgPrecommit,
		Msg:           encodedVote,
		Address:       c.Address(),
		CommittedSeal: []byte{},
	}
	c.SetSentPrecommit(true)
	c.Br().Broadcast(ctx, msg)
}

// TestGarbagePrecommitter broadcasts a garbage precommit message in the network,
// We expect other nodes to detect this misbehaviour and discard these messages
// Receiving nodes should also disconnect misbehaving nodes
func TestGarbagePrecommitter(t *testing.T) {
	users, err := test.Users(6, "10e18,v,100,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)

	//set Malicious users
	users[0].CustHandler = &node.CustomHandler{Precommitter: &garbagePrecommitSender{}}
	// creates a network of 6 users and starts all the nodes in it
	network, err := test.NewNetworkFromUsers(users, true)
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120)
	require.NoError(t, err)
	defer network.Shutdown()
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

type garbagePrevoter struct {
	*core.Core
	interfaces.Prevoter
}

func (c *garbagePrevoter) SendPrevote(ctx context.Context, isNil bool) {
	logger := c.Logger().New("step", c.Step())

	var prevoteMsg messageutils.Vote
	err := gofakeit.Struct(&prevoteMsg)
	if err != nil {
		logger.Error("Failed to fake prevote struct, err ", err)
		return
	}
	encodedVote, err := messageutils.Encode(&prevoteMsg)
	if err != nil {
		logger.Error("Failed to encode", "subject", prevoteMsg)
		return
	}

	msg := &messageutils.Message{
		Code:          messageutils.MsgPrevote,
		Msg:           encodedVote,
		Address:       c.Address(),
		CommittedSeal: []byte{},
	}
	c.Br().Broadcast(ctx, msg)
	c.SetSentPrevote(true)
}

// TestGarbagePrevoter broadcasts a garbage prevote message in the network,
// We expect other nodes to detect this misbehaviour and discard these messages
// Receiving nodes should also disconnect misbehaving nodes
func TestGarbagePrevoter(t *testing.T) {
	users, err := test.Users(6, "10e18,v,100,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)

	//set Malicious users
	users[0].CustHandler = &node.CustomHandler{Prevoter: &garbagePrevoter{}}
	// creates a network of 6 users and starts all the nodes in it
	network, err := test.NewNetworkFromUsers(users, true)
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120)
	require.NoError(t, err)
	defer network.Shutdown()
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

type garbageProposer struct {
	*core.Core
	interfaces.Proposer
}

func (c *garbageProposer) SendProposal(ctx context.Context, p *types.Block) {

	var proposalMsg messageutils.Proposal
	err := gofakeit.Struct(&proposalMsg)
	if err != nil {
		c.Logger().Error("Failed to fake proposal struct, err ", err)
		return
	}
	proposal, _ := messageutils.Encode(proposalMsg)

	c.SetSentProposal(true)
	c.Backend().SetProposedBlockHash(p.Hash())

	//send same proposal twice
	c.Br().Broadcast(ctx, &messageutils.Message{
		Code:          messageutils.MsgProposal,
		Msg:           proposal,
		Address:       c.Address(),
		CommittedSeal: []byte{},
	})
}

// TestGarbagePrevoter broadcasts a garbage proposal message in the network,
// We expect other nodes to detect this misbehaviour and discard these messages
// Receiving nodes should also disconnect misbehaving nodes
func TestGarbageProposer(t *testing.T) {
	users, err := test.Users(6, "10e18,v,100,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)

	//set Malicious users
	users[0].CustHandler = &node.CustomHandler{Proposer: &garbageProposer{}}
	// creates a network of 6 users and starts all the nodes in it
	network, err := test.NewNetworkFromUsers(users, true)
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120)
	require.NoError(t, err)
	defer network.Shutdown()
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}
