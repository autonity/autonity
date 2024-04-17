package collusion

import (
	"context"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/cmd/gengen/gengen"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/stretchr/testify/require"
	"testing"
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

	initCollusionContext(users, autonity.C1, newCollusionC1Planer())

	// creates a network of 8 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(120, 180, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")

	b := getCollusionContext(autonity.C1)
	// the leader should be slashed by PO accusation, since there is no innocence proof for it.
	leader := crypto.PubkeyToAddress(b.leader.NodeKey.PublicKey)
	detected := e2e.AccountabilityEventDetected(t, leader, autonity.Accusation, autonity.PO, network)
	require.Equal(t, true, detected)

	// while the followers should be slashed by C1 accusation since there are no innocence proof for it.
	for _, f := range b.followers {
		faultyAddress := crypto.PubkeyToAddress(f.NodeKey.PublicKey)
		detected = e2e.AccountabilityEventDetected(t, faultyAddress, autonity.Accusation, autonity.C1, network)
		require.Equal(t, true, detected)
	}
}

type collusionC1Planner struct{}

func newCollusionC1Planer() *collusionC1Planner {
	return &collusionC1Planner{}
}

// setupRoles setup message queues for faulty members, and it also setups implementations of interface: propose, prevote
// and precommit for members.
func (p *collusionC1Planner) setupRoles(members []*gengen.Validator) (*gengen.Validator, []*gengen.Validator) {
	// To simulate C1 collusion, we ask a member to be leader to propose an invalid old proposal,
	// and the followers should pre-commit for the invalid proposal as valid without seeing quorum prevotes of it.
	leader := members[0]
	followers := members[1:]

	// setup collusion context functions
	leader.TendermintServices = &interfaces.Services{Broadcaster: newColludedC1Leader}
	for _, f := range followers {
		f.TendermintServices = &interfaces.Services{Precommiter: newColludedC1Follower}
	}
	return leader, followers
}

func newColludedC1Leader(c interfaces.Core) interfaces.Broadcaster {
	return &colludedC1Leader{c.(*core.Core), c.Broadcaster()}
}

type colludedC1Leader struct {
	*core.Core
	interfaces.Broadcaster
}

func (c *colludedC1Leader) Broadcast(msg message.Msg) {
	ctx := getCollusionContext(autonity.C1)
	if !ctx.isReady() {
		c.setupContext()
		c.BroadcastAll(msg)
		return
	}

	h, r, v := ctx.context()

	if h != c.Height().Uint64() {
		c.BroadcastAll(msg)
		return
	}

	// send invalid proposal with the planed data.
	validRound := r - 1
	p := message.NewPropose(r, h, validRound, v, c.Backend().Sign)
	c.BroadcastAll(p)
}

// setupContext, it resolves a future height and round for the colludedC1Leader to set up the collusion context.
func (c *colludedC1Leader) setupContext() {
	leader := c.Address()
	futureHeight := c.Height().Uint64() + 5
	round := int64(0)

	// make sure the round >= 1, thus we can set the valid round > -1
	for ; ; round++ {
		if round == 0 {
			continue
		}
		if validProposer(leader, futureHeight, round, c.Core) {
			break
		}
	}

	getCollusionContext(autonity.C1).setupContext(futureHeight, round, types.NewBlockWithHeader(newBlockHeader(futureHeight)))
}

func newColludedC1Follower(c interfaces.Core) interfaces.Precommiter {
	return &colludedC1Follower{c.(*core.Core), c.Precommiter()}
}

type colludedC1Follower struct {
	*core.Core
	interfaces.Precommiter
}

func (c *colludedC1Follower) SendPrecommit(_ context.Context, _ bool) {
	h, r, v := getCollusionContext(autonity.C1).context()
	// if the leader haven't set up the context, skip.
	if v == nil || h != c.Height().Uint64() {
		return
	}

	// send precommit for the planned invalid proposal.
	precommit := message.NewPrecommit(r, h, v.Hash(), c.Backend().Sign)
	c.SetSentPrecommit(true)
	c.Broadcaster().Broadcast(precommit)
}
