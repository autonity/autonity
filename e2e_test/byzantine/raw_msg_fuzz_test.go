package byzantine

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/backend"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	e2e "github.com/autonity/autonity/e2e_test"
)

/**
 * The raw message fuzzer keeps broadcasting random bytes that replace the rlp encoded bytes stream of the raw msg of
 * ACN protocol, it aims to verify if the consensus network remains live-ness in this context, for those upper layer on
 * top of RLP decoding are covered by each ACN message test set.
 */
func newRawMSGFuzzer(b interfaces.Backend) interfaces.Gossiper {
	defaultGossiper := b.Gossiper()
	return &rawMSGFuzzer{
		Gossiper: defaultGossiper,
		address:  defaultGossiper.Address(),
	}
}

type rawMSGFuzzer struct {
	interfaces.Gossiper
	address common.Address
}

// Faulty node keeps broadcasting fuzz raw message to committee. Every input message of this interface will be fuzzed.
func (fg *rawMSGFuzzer) Gossip(committee types.Committee, msg message.Msg) {
	targets := make([]common.Address, 0)
	i := 0
	for _, val := range committee {
		if val.Address != fg.address {
			targets = append(targets, val.Address)
		}
		i++
	}

	if fg.Broadcaster() == nil || len(targets) == 0 {
		return
	}

	ps := fg.Broadcaster().FindPeers(targets)
	for _, p := range ps {
		randBytes, err := e2e.GenerateRandomBytes(len(msg.Payload()))
		if err != nil {
			panic("Failed to generate random bytes ")
		}
		// send fuzzed raw msg with recognisable msg code.
		go p.SendRaw(backend.NetworkCodes[msg.Code()], randBytes) // nolint
		// send fuzzed raw AskSync message to committee
		go p.SendRaw(backend.SyncNetworkMsg, randBytes) // nolint
		// send fuzzed raw accusation message to committee
		go p.SendRaw(backend.AccountabilityNetworkMsg, randBytes) // nolint
		// send random msg code with fuzzed raw msg.
		go p.SendRaw(rand.Uint64(), randBytes) // nolint
	}
}

func (fg *rawMSGFuzzer) AskSync(_ *types.Header) {
}

func TestRawMessageFuzzer(t *testing.T) {
	numOfNodes := 10
	// create 10 validator nodes with each of them has same voting power.
	vals, err := e2e.Validators(t, numOfNodes, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	f := bft.F(new(big.Int).SetUint64(uint64(numOfNodes)))
	for i := uint64(0); i < f.Uint64(); i++ {
		vals[i].TendermintServices = &interfaces.Services{Gossiper: newRawMSGFuzzer}
	}

	// creates a network of 10 vals and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, vals, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	err = network.WaitForHeight(120, 240)
	require.NoError(t, err)
}
