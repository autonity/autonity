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
  - TestCollusionPC, it creates a faulty party and nominates a leader to prepare an invalid old proposal which carries
    a valid_round with without having quorum prevotes for that value at that valid_round, and the colluded followers
    precommit for that old proposal, thus the proposer should be slashed with PO accusation while the followers should be
    slashed by C1 Accusation.
*/
func TestCollusionPC(t *testing.T) {
	numOfNodes := 8
	users, err := e2e.Validators(t, numOfNodes, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	initCollusion(users, autonity.C1, newCollusionC1Planer())

	// creates a network of 8 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(150, 240, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")

	b := getCollusion(autonity.C1)
	// the leader should be slashed by PO accusation, since there is no innocence proof for it.
	leader := crypto.PubkeyToAddress(b.leader.NodeKey.PublicKey)
	err = e2e.AccountabilityEventDetected(t, leader, autonity.Accusation, autonity.PO, network)
	require.NoError(t, err)

	// while the followers should be slashed by C1 accusation since there are no innocence proof for it.
	for _, f := range b.followers {
		faultyAddress := crypto.PubkeyToAddress(f.NodeKey.PublicKey)
		err = e2e.AccountabilityEventDetected(t, faultyAddress, autonity.Accusation, autonity.C1, network)
		require.NoError(t, err)
	}
}

type collusionC1Planner struct{}

func newCollusionC1Planer() *collusionC1Planner {
	return &collusionC1Planner{}
}

// setupRoles setup message queues for faulty members, and it also setups implementations of interface: propose, prevote
// and precommit for members.
func (p *collusionC1Planner) setupRoles(leader *gengen.Validator, followers []*gengen.Validator) {
	// To simulate C1 collusion, we ask a member to be leader to propose an invalid old proposal,
	// and the followers should pre-commit for the invalid proposal as valid without seeing quorum prevotes of it.
	leader.TendermintServices = &interfaces.Services{Broadcaster: newColludedC1Leader}
	for _, f := range followers {
		f.TendermintServices = &interfaces.Services{Precommiter: newColludedC1Follower}
	}
}

func newColludedC1Leader(c interfaces.Core) interfaces.Broadcaster {
	return &colludedC1Leader{c.(*core.Core), c.Broadcaster()}
}

func newColludedC1Follower(c interfaces.Core) interfaces.Precommiter {
	return &colludedC1Follower{c.(*core.Core), c.Precommiter()}
}

type colludedC1Follower struct {
	*core.Core
	interfaces.Precommiter
}

func (c *colludedC1Follower) SendPrecommit(_ context.Context, _ bool) {
	h, r, v := getCollusion(autonity.C1).context()
	// if the leader haven't set up the context, skip.
	if v == nil || h != c.Height().Uint64() {
		return
	}

	// send precommit for the planned invalid proposal.
	committee, err := c.Backend().BlockChain().CommitteeByHeight(h)
	if err != nil {
		panic(err)
	}
	precommit := message.NewPrecommit(r, h, v.Hash(), c.Backend().Sign, committee.MemberByAddress(c.Address()), committee.Len())
	c.SetSentPrecommit(true)
	c.Broadcaster().Broadcast(precommit)
}

type colludedC1Leader struct {
	*core.Core
	interfaces.Broadcaster
}

func (c *colludedC1Leader) Broadcast(msg message.Msg) {
	sendProposal(c, autonity.C1, msg)
}

// setupContext, it resolves a future height and round for the colludedC1Leader to set up the collusion context.
func (c *colludedC1Leader) SetupCollusionContext() {
	setupCollusionContext(c, autonity.C1)
}
