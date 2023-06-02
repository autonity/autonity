package byzantine

import (
	"context"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/stretchr/testify/require"
)

func newPreVoteSpammer(c interfaces.Core) interfaces.Prevoter {
	return &preVoteSpammer{c.(*core.Core), c.Prevoter()}
}

type preVoteSpammer struct {
	*core.Core
	interfaces.Prevoter
}

func (c *preVoteSpammer) SendPrevote(_ context.Context, isNil bool) {
	var prevote *message.Prevote
	if isNil {
		prevote = message.NewPrevote(c.Round(), c.Height().Uint64(), common.Hash{}, c.Backend().Sign)
	} else {
		h := c.CurRoundMessages().ProposalHash()
		if h == (common.Hash{}) {
			c.Logger().Error("sendPrecommit Proposal is empty! It should not be empty!")
			return
		}
		prevote = message.NewPrevote(c.Round(), c.Height().Uint64(), common.Hash{}, c.Backend().Sign)
	}

	for i := 0; i < 1000; i++ {
		c.BroadcastAll(prevote)
	}
	c.SetSentPrevote(true)
}

// TestPrevoteSpammer spams the network by broadcasting 4k preovte messages at once
func TestPrevoteSpammer(t *testing.T) {
	users, err := e2e.Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	//set Malicious users
	users[0].TendermintServices = &interfaces.Services{Prevoter: newPreVoteSpammer}
	users[1].TendermintServices = &interfaces.Services{Prevoter: newPreVoteSpammer}
	// creates a network of 6 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

type precommitSpammer struct {
	*core.Core
	interfaces.Precommiter
}

func newPrecommitSpammer(c interfaces.Core) interfaces.Precommiter {
	return &precommitSpammer{c.(*core.Core), c.Precommiter()}
}

func (c *precommitSpammer) SendPrecommit(_ context.Context, isNil bool) {
	var precommit *message.Precommit
	if isNil {
		precommit = message.NewPrecommit(c.Round(), c.Height().Uint64(), common.Hash{}, c.Backend().Sign)
	} else {
		h := c.CurRoundMessages().ProposalHash()
		if h == (common.Hash{}) {
			c.Logger().Error("core.sendPrecommit Proposal is empty! It should not be empty!")
			return
		}
		precommit = message.NewPrecommit(c.Round(), c.Height().Uint64(), h, c.Backend().Sign)
	}
	for i := 0; i < 1000; i++ {
		c.BroadcastAll(precommit)
	}
	c.SetSentPrecommit(true)
}

// TestPrecommitSpammer spams the network by broadcasting 4k precommit messages at once
func TestPrecommitSpammer(t *testing.T) {
	users, err := e2e.Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	//set Malicious users
	users[0].TendermintServices = &interfaces.Services{Precommiter: newPrecommitSpammer}
	// creates a network of 6 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

type proposalSpammer struct {
	*core.Core
	interfaces.Proposer
}

func newProposalSpammer(c interfaces.Core) interfaces.Proposer {
	return &proposalSpammer{c.(*core.Core), c.Proposer()}
}

func (c *proposalSpammer) SendProposal(_ context.Context, p *types.Block) {
	proposal := message.NewPropose(c.Round(), c.Height().Uint64(), c.ValidRound(), p, c.Backend().Sign)
	c.SetSentProposal(true)
	c.Backend().SetProposedBlockHash(p.Hash())
	for i := 0; i < 1000; i++ {
		c.BroadcastAll(proposal)
	}
}

// TestProposalSpammer spams the network by broadcasting 4k proposal messages at once
func TestProposalSpammer(t *testing.T) {
	users, err := e2e.Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	//set Malicious proposalSender
	users[0].TendermintServices = &interfaces.Services{Proposer: newProposalSpammer}
	// creates a network of 6 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	err = network.WaitForSyncComplete()
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}
