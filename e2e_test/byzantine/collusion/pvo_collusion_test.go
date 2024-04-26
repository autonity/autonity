package collusion

import (
	"context"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/cmd/gengen/gengen"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/crypto"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/stretchr/testify/require"
	"testing"
)

/*
*
  - TestCollusionPVO, it creates a faulty party and nominates a proposer to propose an invalid old proposal which carries
    a valid_round with without having quorum prevotes for that value at that valid_round, and the colluded followers
    prevote for that old proposal, thus the proposer should be slashed with PO accusation while the followers should be
    slashed by PVO Accusation.
*/
func TestCollusionPVO(t *testing.T) {
	numOfNodes := 8
	users, err := e2e.Validators(t, numOfNodes, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	initCollusion(users, autonity.PVO, newCollusionPVOPlaner())

	// creates a network of 8 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(120, 180, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")

	b := getCollusion(autonity.PVO)
	// the leader should be slashed by PO accusation, since there is no innocence proof for it.
	leader := crypto.PubkeyToAddress(b.leader.NodeKey.PublicKey)
	detected := e2e.AccountabilityEventDetected(t, leader, autonity.Accusation, autonity.PO, network)
	require.Equal(t, true, detected)

	// while the followers should be slashed by PVO accusation since there are no innocence proof for it.
	for _, f := range b.followers {
		faultyAddress := crypto.PubkeyToAddress(f.NodeKey.PublicKey)
		detected = e2e.AccountabilityEventDetected(t, faultyAddress, autonity.Accusation, autonity.PVO, network)
		require.Equal(t, true, detected)
	}
}

type collusionPVOPlanner struct{}

func newCollusionPVOPlaner() *collusionPVOPlanner {
	return &collusionPVOPlanner{}
}

// setupRoles setup message queues for faulty members, and it also setups implementations of interface: propose, prevote
// and precommit for members.
func (p *collusionPVOPlanner) setupRoles(leader *gengen.Validator, followers []*gengen.Validator) {
	// To simulate PVO collusion, we ask a member to be leader to propose an invalid old proposal,
	// and the followers should pre-vote for the invalid proposal as a valid one.
	leader.TendermintServices = &interfaces.Services{Broadcaster: newColludedPVOLeader}
	for _, f := range followers {
		f.TendermintServices = &interfaces.Services{Prevoter: newColludedPVOFollower}
	}
}

func newColludedPVOFollower(c interfaces.Core) interfaces.Prevoter {
	return &colludedPVOFollower{c.(*core.Core), c.Prevoter()}
}

type colludedPVOFollower struct {
	*core.Core
	interfaces.Prevoter
}

func (c *colludedPVOFollower) SendPrevote(_ context.Context, _ bool) {
	sendPrevote(c.Core, autonity.PVO)
}

func newColludedPVOLeader(c interfaces.Core) interfaces.Broadcaster {
	return &colludedPVOLeader{c.(*core.Core), c.Broadcaster()}
}

type colludedPVOLeader struct {
	*core.Core
	interfaces.Broadcaster
}

func (c *colludedPVOLeader) Broadcast(msg message.Msg) {
	sendProposal(c, autonity.PVO, msg)
}

// setupContext, it resolves a future height and round for the colludedPVOLeader to set up the collusion context.
func (c *colludedPVOLeader) SetupCollusionContext() {
	setupCollusionContext(c, autonity.PVO)
}
