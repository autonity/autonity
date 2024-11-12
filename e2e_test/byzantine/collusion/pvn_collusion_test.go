package collusion

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/cmd/gengen/gengen"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/crypto"
	e2e "github.com/autonity/autonity/e2e_test"
)

/*
*
  - TestCollusionPVN, it creates a faulty party and nominates a leader to propose an invalid new proposal, while
    the followers prevote for that invalid proposal as a valid one. Since new proposal is not accountable for the time
    being, thus we cannot expect the proposer is slashed, however we can slash those followers by PVN accusation rule.
*/
func TestCollusionPVN(t *testing.T) {
	t.Skip("Flaky test")
	numOfNodes := 8
	users, err := e2e.Validators(t, numOfNodes, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	initCollusion(users, autonity.PVN, newCollusionPVNPlaner())

	// creates a network of 8 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(120, 150, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")

	// Accusation of PVN should rise since followers prevote for the planed invalid value.
	// The followers should be slashed, as a new proposal is not accountable now, thus we cannot
	// slash the malicious proposer in this test.
	b := getCollusion(autonity.PVN)
	for _, f := range b.followers {
		faultyAddress := crypto.PubkeyToAddress(f.NodeKey.PublicKey)
		err = e2e.AccountabilityEventDetected(t, faultyAddress, autonity.Accusation, autonity.PVN, network)
		require.NoError(t, err)
	}
}

type collusionPVNPlanner struct{}

func newCollusionPVNPlaner() *collusionPVNPlanner {
	return &collusionPVNPlanner{}
}

// setupRoles setup message queues for faulty members, and it also setups implementations of interface: propose, prevote
// and precommit for members.
func (p *collusionPVNPlanner) setupRoles(leader *gengen.Validator, followers []*gengen.Validator) {
	// To simulate PVN collusion, we ask a member to be leader to propose an invalid new proposal,
	// and the followers should pre-vote for the invalid proposal as a valid one.
	leader.TendermintServices = &interfaces.Services{Broadcaster: newColludedPVNLeader}
	for _, f := range followers {
		f.TendermintServices = &interfaces.Services{Prevoter: newColludedPVNFollower}
	}
}

func newColludedPVNFollower(c interfaces.Core) interfaces.Prevoter {
	return &colludedPVNFollower{c.(*core.Core), c.Prevoter()}
}

type colludedPVNFollower struct {
	*core.Core
	interfaces.Prevoter
}

func (c *colludedPVNFollower) SendPrevote(_ context.Context, _ bool) {
	sendPrevote(c.Core, autonity.PVN)
}

func newColludedPVNLeader(c interfaces.Core) interfaces.Broadcaster {
	return &colludedPVNLeader{c.(*core.Core), c.Broadcaster()}
}

type colludedPVNLeader struct {
	*core.Core
	interfaces.Broadcaster
}

func (c *colludedPVNLeader) Broadcast(msg message.Msg) {
	sendProposal(c, autonity.PVN, msg)
}

// setupContext, it resolves a future height and round for the colludedPVNLeader to set up the collusion context.
func (c *colludedPVNLeader) SetupCollusionContext() {
	setupCollusionContext(c, autonity.PVN)
}
