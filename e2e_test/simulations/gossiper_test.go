package simulations

import (
	"math"
	"math/rand"
	"testing"

	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/backend"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/log"
)

func newCustomGossiper(b interfaces.Backend) interfaces.Gossiper {
	defaultGossiper := b.Gossiper()
	return &customGossiper{
		Gossiper:       defaultGossiper,
		knownMessages:  defaultGossiper.KnownMessages(),
		recentMessages: defaultGossiper.RecentMessages(),
		address:        defaultGossiper.Address(),
	}
}

type customGossiper struct {
	interfaces.Gossiper
	knownMessages  *lru.ARCCache
	recentMessages *lru.ARCCache
	address        common.Address
}

// this is a test custom gossip function, just to illustrate how to build one
// it gossips only to a random set of ceil(sqrt(N)). It is not optimized.
func (cg *customGossiper) Gossip(committee types.Committee, msg message.Msg) {
	hash := msg.Hash()
	cg.knownMessages.Add(hash, true)

	// determine random subset of committee members to gossip to
	// if by chance we include our own index, we will end up gossiping to
	// ceil(sqrt(N)) - 1
	fullset := rand.Perm(len(committee)) //nolint
	num := uint(math.Ceil(math.Sqrt(float64(len(committee)))))
	subset := fullset[:num] //nolint

	targets := make(map[common.Address]struct{})
	i := 0
	for _, val := range committee {
		if val.Address != cg.address && slices.Contains(subset, i) {
			targets[val.Address] = struct{}{}
		}
		i++
	}

	// simulate network delay before gossiping the packet (between 0 and 200ms)
	// by simulating it at this point, it means we simulate roughly the same delay
	// for all gossips of the same message
	//time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)

	if cg.Broadcaster() == nil || len(targets) == 0 {
		return
	}

	ps := cg.Broadcaster().FindPeers(targets)
	for addr, p := range ps {
		ms, ok := cg.recentMessages.Get(addr)
		var m *lru.ARCCache
		if ok {
			m, _ = ms.(*lru.ARCCache)
			if _, k := m.Get(hash); k {
				// This peer had this event, skip it
				continue
			}
		} else {
			m, _ = lru.NewARC(1024) //   backend.inmemoryMessages  = 1024
		}

		m.Add(hash, true)
		cg.recentMessages.Add(addr, m)

		go p.SendRaw(backend.NetworkCodes[msg.Code()], msg.Payload()) //nolint
	}
}

func (cg *customGossiper) AskSync(_ *types.Header) {
	// I disable the ask sync recovery mechanism, so that I can see if the gossip only is enough to keep the network live
	log.Info("liveness lost, supposed to ask sync (but will not)")
}

// this test just has the purpose of verifying that the customGossiper works as intended
func TestCustomGossiper(t *testing.T) {
	t.Skip("Flacky in CI, remove SKIP only locally.")
	vals, err := e2e.Validators(t, 10, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	for _, val := range vals {
		val.TendermintServices = &interfaces.Services{Gossiper: newCustomGossiper}
	}
	// creates a network of 6 vals and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, vals, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be able to mine blocks respecting the 1 block/s rate
	err = network.WaitToMineNBlocks(10, 60, true)
	require.NoError(t, err)
}
