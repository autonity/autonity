package byzantine

import (
	"context"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/stretchr/testify/require"
	"testing"
)

type malCommitterService struct {
	*core.Core
	interfaces.Committer
}

func newMalCommitterService(c interfaces.Core) interfaces.Committer {
	return &malCommitterService{c.(*core.Core), c.Committer()}
}

func (c *malCommitterService) Commit(ctx context.Context, _ int64, messages *message.RoundMessages) {
	c.SetStep(ctx, core.PrecommitDone)
	proposal := messages.Proposal()
	if proposal == nil {
		// Should never happen really. Let's panic to catch bugs.
		panic("Core commit called with empty proposal")
	}
	// in this test context, we just set the decision in WAL without commit it to the blockchain.
	// in the test context, we keep reset the node, on the start of the node, it will commit the
	// decision from the WAL to the blockchain to make the blockchain head move forward.
	c.SetDecision(proposal.Block(), proposal.R())
}

func TestCommitDecisionFromWAL(t *testing.T) {
	// create a single node network, let the validator save the decision in WAL without commit it to the blockchain.
	users, err := e2e.Validators(t, 1, "10e18,v,10,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	//set Malicious users
	users[0].TendermintServices = &interfaces.Services{Committer: newMalCommitterService}

	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// keep reset the validator, let the decision be committed from WAL to the blockchain on node restart.
	resets := 4
	for i := 0; i < resets; i++ {
		err = network[0].Close(false)
		require.NoError(t, err)
		network[0].Wait()
		network.WaitToMineNBlocks(1, 10, false)
		err = network[0].Start()
		require.NoError(t, err)
		err = network[0].Eth.StartMining(1)
		require.NoError(t, err)
		network.WaitToMineNBlocks(1, 10, false)
	}

	require.Equal(t, uint64(3), network[0].Eth.BlockChain().CurrentHeader().Number.Uint64())
}
