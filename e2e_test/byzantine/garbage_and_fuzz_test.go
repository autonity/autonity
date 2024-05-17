package byzantine

import (
	"context"
	"crypto/sha256"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"math/big"
	"math/rand"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/e2e_test"
)

//todo: (Jason) the fuzz of signature field in the propose, prevote and precommit messages are still missing.

func newRandomBytesBroadcaster(c interfaces.Core) interfaces.Broadcaster {
	return &randomBytesBroadcaster{c.(*core.Core)}
}

type randomBytesBroadcaster struct {
	*core.Core
}

func (s *randomBytesBroadcaster) Broadcast(_ message.Msg) {
	logger := s.Logger().New("step", s.Step())
	logger.Info("Broadcasting random bytes")

	for i := 0; i < 1000; i++ {
		payload, err := e2e.GenerateRandomBytes(2048)
		if err != nil {
			logger.Error("Failed to generate random bytes ", "err", err)
			return
		}
		var hash common.Hash
		copy(hash[:], payload)
		msg := message.Fake{FakeCode: 1, FakePayload: payload, FakeHash: hash}
		s.BroadcastAll(msg)
	}
}

// TestRandomBytesBroadcaster broadcasts random bytes in the network,
// We expect other nodes to detect this misbehaviour and discard these messages
// Receiving nodes should also disconnect misbehaving nodes
func TestRandomBytesBroadcaster(t *testing.T) {
	numOfNodes := 6
	users, err := e2e.Validators(t, numOfNodes, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	f := bft.F(new(big.Int).SetUint64(uint64(numOfNodes)))
	for i := uint64(0); i < f.Uint64(); i++ {
		//set Malicious users
		users[i].TendermintServices = &interfaces.Services{Broadcaster: newRandomBytesBroadcaster}
	}

	// creates a network of 6 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 180, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

func newGarbageMessageBroadcaster(c interfaces.Core) interfaces.Broadcaster {
	return &garbageMessageBroadcaster{c.(*core.Core)}
}

type garbageMessageBroadcaster struct {
	*core.Core
}

func (s *garbageMessageBroadcaster) Broadcast(_ message.Msg) {
	logger := s.Logger().New("step", s.Step())
	var fMsg message.Fake
	f := fuzz.New()
	f.Fuzz(&fMsg)
	logger.Info("Broadcasting random bytes")
	s.BroadcastAll(&fMsg)
}

// TestGarbageMessageBroadcaster broadcasts a garbage Messages in the network,
// We expect other nodes to detect this misbehaviour and discard these messages
// Receiving nodes should also disconnect misbehaving nodes
func TestGarbageMessageBroadcaster(t *testing.T) {
	numOfNodes := 6
	users, err := e2e.Validators(t, numOfNodes, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	f := bft.F(new(big.Int).SetUint64(uint64(numOfNodes)))
	for i := uint64(0); i < f.Uint64(); i++ {
		//set Malicious users
		users[i].TendermintServices = &interfaces.Services{Broadcaster: newGarbageMessageBroadcaster}
	}

	// creates a network of 6 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 180, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

func newFuzzPrecommitSender(c interfaces.Core) interfaces.Precommiter {
	return &fuzzPrecommitSender{c.(*core.Core), c.Precommiter()}
}

type fuzzPrecommitSender struct {
	*core.Core
	interfaces.Precommiter
}

func (c *fuzzPrecommitSender) SendPrecommit(_ context.Context, isNil bool) {
	var precommit *message.Precommit
	r := rand.Int63()
	h := rand.Uint64()
	if isNil {
		precommit = message.NewPrecommit(r, h, common.Hash{}, c.Backend().Sign)
	} else {
		precommit = message.NewPrecommit(r, h, randHash(), c.Backend().Sign)
	}
	c.SetSentPrecommit(true)
	c.BroadcastAll(precommit)
}

// TestFuzzPrecommitter broadcasts a garbage precommit message in the network,
// We expect other nodes to detect this misbehaviour and discard these messages
// Receiving nodes should also disconnect misbehaving nodes
func TestFuzzPrecommitter(t *testing.T) {
	numOfNodes := 6
	users, err := e2e.Validators(t, numOfNodes, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	f := bft.F(new(big.Int).SetUint64(uint64(numOfNodes)))
	for i := uint64(0); i < f.Uint64(); i++ {
		//set Malicious users
		users[i].TendermintServices = &interfaces.Services{Precommiter: newFuzzPrecommitSender}
	}

	// creates a network of 6 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

func newFuzzPrevoter(c interfaces.Core) interfaces.Prevoter {
	return &fuzzPrevoter{c.(*core.Core), c.Prevoter()}
}

type fuzzPrevoter struct {
	*core.Core
	interfaces.Prevoter
}

func (c *fuzzPrevoter) SendPrevote(_ context.Context, isNil bool) {
	var prevote *message.Prevote
	r := rand.Int63()
	h := rand.Uint64()
	if isNil {
		prevote = message.NewPrevote(r, h, common.Hash{}, c.Backend().Sign)
	} else {
		prevote = message.NewPrevote(r, h, randHash(), c.Backend().Sign)
	}
	c.SetSentPrevote(true)
	c.BroadcastAll(prevote)
}

// TestFuzzPrevoter broadcasts a garbage prevote message in the network,
// We expect other nodes to detect this misbehaviour and discard these messages
// Receiving nodes should also disconnect misbehaving nodes
func TestFuzzPrevoter(t *testing.T) {
	numOfNodes := 6
	users, err := e2e.Validators(t, numOfNodes, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	f := bft.F(new(big.Int).SetUint64(uint64(numOfNodes)))
	for i := uint64(0); i < f.Uint64(); i++ {
		//set Malicious users
		users[i].TendermintServices = &interfaces.Services{Prevoter: newFuzzPrevoter}
	}

	// creates a network of 6 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

func newFuzzProposer(c interfaces.Core) interfaces.Proposer {
	return &fuzzProposer{c.(*core.Core), c.Proposer()}
}

type fuzzProposer struct {
	*core.Core
	interfaces.Proposer
}

/*
type structNode struct {
	fName string
	sMap  map[string]*structNode
	fList []string
}

func generateFieldMap(v interface{}) map[string]reflect.Value {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		panic("Need pointer!")
	}
	outMap := make(map[string]reflect.Value)
	e := reflect.ValueOf(v).Elem()
	for i := 0; i < e.NumField(); i++ {
		fmt.Println("handling field => ", e.Type().Field(i).Name)
		if e.Field(i).Type().Kind() == reflect.Ptr {
			fKind := e.Field(i).Type().Elem().Kind()
			if fKind == reflect.Struct {
				fmt.Println("TODO - handle recursively")
			}
		} else if e.Field(i).Type().Kind() == reflect.Struct {
			fmt.Println("TODO - handle recursively")
		}
		outMap[e.Type().Field(i).Name] = e.Field(i)
	}
	return outMap
}
*/
// duplicated with TestInvalidBlockProposal in proposal_test.go
func (c *fuzzProposer) SendProposal(_ context.Context, p *types.Block) {
	f := fuzz.New()
	var num big.Int
	f.Fuzz(&num)
	e2e.FuzBlock(p, &num)
	proposal := message.NewPropose(c.Round(), c.Height().Uint64(), c.ValidRound(), p, c.Backend().Sign)
	c.SetSentProposal(true)
	c.Backend().SetProposedBlockHash(p.Hash())

	c.BroadcastAll(proposal)
}

// TestFuzzProposer broadcasts a garbage proposal message in the network,
// We expect other nodes to detect this misbehaviour and discard these messages
// Receiving nodes should also disconnect misbehaving nodes
func TestFuzzProposer(t *testing.T) {
	numOfNodes := 6
	users, err := e2e.Validators(t, numOfNodes, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	f := bft.F(new(big.Int).SetUint64(uint64(numOfNodes)))
	for i := uint64(0); i < f.Uint64(); i++ {
		//set Malicious users
		users[i].TendermintServices = &interfaces.Services{Proposer: newFuzzProposer}
	}

	// creates a network of 6 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

func randHash() common.Hash {
	randBytes, err := e2e.GenerateRandomBytes(32)
	if err != nil {
		return common.Hash{}
	}
	return sha256.Sum256(randBytes)
}
