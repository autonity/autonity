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
  - TestCollusionPVO, it creates a faulty party and nominates a proposer to propose an invalid old proposal which carries
    a valid_round with without having quorum prevotes for that value at that valid_round, and the colluded followers still
    prevote for that old proposal, thus the proposer should be slashed with PO accusation while the followers should be
    slashed by PVO Accusation.
*/
func TestCollusionPVO(t *testing.T) {
	numOfNodes := 8
	users, err := e2e.Validators(t, numOfNodes, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	initCollusionContext(users, autonity.PVO, newCollusionPVOPlaner())

	// creates a network of 8 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(120, 180, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")

	b := getCollusionContext(autonity.PVO)
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
func (p *collusionPVOPlanner) setupRoles(members []*gengen.Validator) (*gengen.Validator, []*gengen.Validator) {
	// To simulate PVO collusion, we ask a member to be leader to propose an invalid old proposal,
	// and the followers should pre-vote for the invalid proposal as valid.
	leader := members[0]
	followers := members[1:]

	// start to setupRoles the message queues by putting messages in it.
	leader.TendermintServices = &interfaces.Services{Broadcaster: newColludedPVOLeader}
	for _, f := range followers {
		f.TendermintServices = &interfaces.Services{Prevoter: newColludedPVOFollower}
	}
	return leader, followers
}

func newColludedPVOFollower(c interfaces.Core) interfaces.Prevoter {
	return &colludedPVOFollower{c.(*core.Core), c.Prevoter()}
}

type colludedPVOFollower struct {
	*core.Core
	interfaces.Prevoter
}

func (c *colludedPVOFollower) SendPrevote(_ context.Context, _ bool) {
	h, r, v := getCollusionContext(autonity.PVO).context()
	// if the leader haven't set up the context, skip.
	if v == nil || h != c.Height().Uint64() {
		return
	}

	// send prevote for the planned invalid proposal.
	vote := message.NewPrevote(r, h, v.Hash(), c.Backend().Sign)
	c.SetSentPrevote(true)
	c.BroadcastAll(vote)
}

func newColludedPVOLeader(c interfaces.Core) interfaces.Broadcaster {
	return &colludedPVOLeader{c.(*core.Core), c.Broadcaster()}
}

type colludedPVOLeader struct {
	*core.Core
	interfaces.Broadcaster
}

func (c *colludedPVOLeader) Broadcast(msg message.Msg) {
	ctx := getCollusionContext(autonity.PVO)
	if !ctx.isReady() {
		// resolve a future height, and setup the context on that height.
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

// setupContext, it resolves a future height and round for the colludedPVOLeader to set up the collusion context.
func (c *colludedPVOLeader) setupContext() {
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

	getCollusionContext(autonity.PVO).setupContext(futureHeight, round, types.NewBlockWithHeader(newBlockHeader(futureHeight)))
}
