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
	"math/big"
	"testing"
)

/*
*
  - TestCollusionPVN, it creates a faulty party and nominates a leader to propose an invalid new proposal, while
    the followers prevote for that invalid proposal as a valid one. Since new proposal is not accountable for the time
    being, thus we cannot expect the proposer is slashed, however we can slash those followers by PVN accusation rule.
*/
func TestCollusionPVN(t *testing.T) {
	numOfNodes := 8
	users, err := e2e.Validators(t, numOfNodes, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	initCollusionContext(users, autonity.PVN, newCollusionPVNPlaner())

	// creates a network of 8 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(60, 120, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")

	// Accusation of PVN should rise since followers prevote for the planed invalid value.
	// The followers should be slashed, as a new proposal is not accountable now, thus we cannot
	// slash the malicious proposer in this test.
	b := getCollusionContext(autonity.PVN)
	for _, f := range b.followers {
		faultyAddress := crypto.PubkeyToAddress(f.NodeKey.PublicKey)
		detected := e2e.AccountabilityEventDetected(t, faultyAddress, autonity.Accusation, autonity.PVN, network)
		require.Equal(t, true, detected)
	}
}

type collusionPVNPlanner struct{}

func newCollusionPVNPlaner() *collusionPVNPlanner {
	return &collusionPVNPlanner{}
}

// setupRoles setup message queues for faulty members, and it also setups implementations of interface: propose, prevote
// and precommit for members.
func (p *collusionPVNPlanner) setupRoles(members []*gengen.Validator) (*gengen.Validator, []*gengen.Validator) {
	// To simulate PVN collusion, we ask a member to be leader to propose an invalid new proposal,
	// and the followers should pre-vote for the invalid proposal as a valid one.
	leader := members[0]
	followers := members[1:]

	// setup collusion context functions
	leader.TendermintServices = &interfaces.Services{Broadcaster: newColludedPVNLeader}
	for _, f := range followers {
		f.TendermintServices = &interfaces.Services{Prevoter: newColludedPVNFollower}
	}
	return leader, followers
}

func newColludedPVNFollower(c interfaces.Core) interfaces.Prevoter {
	return &colludedPVNFollower{c.(*core.Core), c.Prevoter()}
}

type colludedPVNFollower struct {
	*core.Core
	interfaces.Prevoter
}

func (c *colludedPVNFollower) SendPrevote(_ context.Context, _ bool) {
	h, r, invalidValue := getCollusionContext(autonity.PVN).context()
	// if the leader haven't set up the context, skip.
	if invalidValue == nil || h != c.Height().Uint64() || r != c.Round() {
		return
	}

	// send prevote for the planned invalid proposal.
	vote := message.NewPrevote(r, h, invalidValue.Hash(), c.Backend().Sign)
	c.SetSentPrevote(true)
	c.BroadcastAll(vote)
}

func newColludedPVNLeader(c interfaces.Core) interfaces.Broadcaster {
	return &colludedPVNLeader{c.(*core.Core), c.Broadcaster()}
}

type colludedPVNLeader struct {
	*core.Core
	interfaces.Broadcaster
}

func (c *colludedPVNLeader) Broadcast(msg message.Msg) {
	ctx := getCollusionContext(autonity.PVN)
	if !ctx.isReady() {
		c.setupContext()
		c.BroadcastAll(msg)
		return
	}

	h, r, v := ctx.context()

	if h != c.Height().Uint64() || r != c.Round() {
		c.BroadcastAll(msg)
		return
	}

	// send invalid proposal with the planed data.
	p := message.NewPropose(r, h, -1, v, c.Backend().Sign)
	c.BroadcastAll(p)
}

// setupContext, it resolves a future height and round for the colludedPVNLeader to set up the collusion context.
func (c *colludedPVNLeader) setupContext() {
	leader := c.Address()
	futureHeight := c.Height().Uint64() + 5
	round := int64(0)

	for ; ; round++ {
		if !validProposer(leader, futureHeight, round, c.Core) {
			break
		}
	}

	b := types.NewBlockWithHeader(newBlockHeader(futureHeight))
	e2e.FuzBlock(b, new(big.Int).SetUint64(futureHeight))
	getCollusionContext(autonity.PVN).setupContext(futureHeight, round, b)
}
