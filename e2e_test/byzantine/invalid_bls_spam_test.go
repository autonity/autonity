package byzantine

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core"
	ccore "github.com/autonity/autonity/core"

	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/crypto/blst"
	e2e "github.com/autonity/autonity/e2e_test"
)

// to test invalid bls signature spamming and related disconnection efficacy

func newInvalidSignatureBroadcaster(c interfaces.Core) interfaces.Prevoter {
	return &invalidSignatureBroadcaster{c.(*core.Core), c.Prevoter(), false}
}

type invalidSignatureBroadcaster struct {
	*core.Core
	interfaces.Prevoter
	sent bool
}

// when sending a prevote, do the standard behaviour + send an invalid signature prevote
func (c *invalidSignatureBroadcaster) SendPrevote(ctx context.Context, isNil bool) {
	// send invalid sig
	if !c.sent {
		invalidSigner := func(hash common.Hash) blst.Signature {
			var h common.Hash
			rand.Read(h[:])
			return c.Backend().Sign(h)
		}
		hash := c.CurRoundMessages().ProposalHash()
		self, csize := selfAndCsize(c.Core, c.Height().Uint64())
		prevote := message.NewPrevote(c.Round(), c.Height().Uint64(), hash, invalidSigner, self, csize)
		c.Backend().Gossip(c.CommitteeSet().Committee(), prevote)
		c.sent = true
	}

	// standard behaviour
	c.Prevoter.SendPrevote(ctx, isNil)
}

func TestInvalidBlsSignatureDisconnection(t *testing.T) {
	t.Run("Malicious peer sending an invalid BLS signature should be disconnect for at least 1 epoch", func(t *testing.T) {
		n := 4
		validators, err := e2e.Validators(t, n, "10e36,v,100,0.0.0.0:%s,%s,%s,%s")
		require.NoError(t, err)

		// set malicious handler
		malicious := 0
		validators[malicious].TendermintServices = &interfaces.Services{Prevoter: newInvalidSignatureBroadcaster}

		// creates a network of 4 validators and starts all the nodes in it
		// modify epoch period to ensure that it is > standard p2p suspension period (60 blocks currently)
		// we also modify the PastPerformanceWeight to be 100%, so that validator inactivity always remain 0.
		// We do not want omission jailing to interfere in this test.
		network, err := e2e.NewNetworkFromValidators(t, validators, true, func(genesis *ccore.Genesis) {
			genesis.Config.AutonityContractConfig.EpochPeriod = 100
			genesis.Config.OmissionAccountabilityConfig.PastPerformanceWeight = 10000
			genesis.Config.OmissionAccountabilityConfig.InactivityThreshold = 10000
		})
		require.NoError(t, err)
		defer network.Shutdown(t)

		epochPeriod := network[0].Eth.BlockChain().ProtocolContracts().Cache.EpochPeriod().Uint64()

		// mine some blocks, but do not terminate the epoch
		err = network.WaitToMineNBlocks(epochPeriod-10, int(epochPeriod*2), false)
		require.NoError(t, err)

		for i, node := range network {
			// execution peers should not be affected
			require.Equal(t, n-1, len(node.ExecutionServer().PeersInfo()))

			// consensus peers, malicious peer should be disconnected from everyone
			if i == malicious {
				require.Equal(t, 0, len(node.ConsensusServer().PeersInfo()))
			} else {
				require.Equal(t, n-2, len(node.ConsensusServer().PeersInfo()))
			}
		}

		// close the epoch and wait a bit, the malicious node should be reconnected to the rest of the network
		err = network.WaitToMineNBlocks(epochPeriod*2, int(epochPeriod*4), false)
		require.NoError(t, err)

		for _, node := range network {
			require.Equal(t, n-1, len(node.ExecutionServer().PeersInfo()))
			require.Equal(t, n-1, len(node.ConsensusServer().PeersInfo()))
		}
	})
}
