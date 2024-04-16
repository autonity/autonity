package collusion

import (
	"context"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/cmd/gengen/gengen"
	"github.com/autonity/autonity/common"
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
  - TestCollusionPN, it creates a faulty parity and nominates a proposer to propose an invalid new proposal, while the
    followers of the prevote for that invalid proposal as valid. Thus
*/
func TestCollusionPN(t *testing.T) {
	numOfNodes := 7
	users, err := e2e.Validators(t, numOfNodes, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	createCollusionParty(users, 3, 0, autonity.PN, newCollusionPNPlaner())

	// creates a network of 6 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// Accusation of PVN should rise since followers prevote for the planed invalid value.
	// The faulty party should be slashed, as proposal is not accountable now, thus we cannot
	// slash the malicious proposer.

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(60, 120, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")

	// Accusation of PVN should rise since followers prevote for the planed invalid value.
	// The faulty party should be slashed, as proposal is not accountable now, thus we cannot
	// slash the malicious proposer.
	b := collusionBehaviour(autonity.PN)
	for _, f := range b.followers {
		faultyAddress := crypto.PubkeyToAddress(f.NodeKey.PublicKey)
		detected := e2e.AccountabilityEventDetected(t, faultyAddress, autonity.Accusation, autonity.PVN, network)
		require.Equal(t, true, detected)
	}
}

type collusionPNPlanner struct{}

func newCollusionPNPlaner() *collusionPNPlanner {
	return &collusionPNPlanner{}
}

// plan setup message queues for faulty memebers, and it also setups implementations of interface: propose, prevote
// and precommit for members.
func (p *collusionPNPlanner) plan(_ uint64, _ int64, _ *types.Block, members []*gengen.Validator) (map[common.Address]map[uint64]map[int64]map[core.Step]message.Msg, *gengen.Validator, []*gengen.Validator) {
	// To simulate PN collusion, we ask a member to be leader to propose an invalid proposal,
	// and the followers should pre-vote for the invalid proposal as valid.
	leader := members[0]
	followers := members[1:]

	// start to plan the message queues by putting messages in it.
	leader.TendermintServices = &interfaces.Services{Broadcaster: newColludedProposer}
	for _, f := range followers {
		f.TendermintServices = &interfaces.Services{Prevoter: newColludedPNPrevoter}
	}
	// the collusion of PN does not need message queues.
	return nil, leader, followers
}

func newColludedPNPrevoter(c interfaces.Core) interfaces.Prevoter {
	return &colludedPNPrevoter{c.(*core.Core), c.Prevoter()}
}

type colludedPNPrevoter struct {
	*core.Core
	interfaces.Prevoter
}

func (c *colludedPNPrevoter) SendPrevote(_ context.Context, _ bool) {
	behaviours := collusionBehaviour(autonity.PN)
	if behaviours == nil || behaviours.h != c.Height().Uint64() || behaviours.r != c.Round() {
		return
	}

	// send prevote for the invalid proposal.
	vote := message.NewPrevote(behaviours.r, behaviours.h, behaviours.invalidValue.Hash(), c.Backend().Sign)
	c.SetSentPrevote(true)
	c.BroadcastAll(vote)
}

func newColludedProposer(c interfaces.Core) interfaces.Broadcaster {
	return &colludedProposer{c.(*core.Core), c.Broadcaster()}
}

type colludedProposer struct {
	*core.Core
	interfaces.Broadcaster
}

func (c *colludedProposer) Broadcast(msg message.Msg) {
	behaviours := collusionBehaviour(autonity.PN)
	if behaviours == nil || behaviours.h != c.Height().Uint64() || behaviours.r != c.Round() {
		c.BroadcastAll(msg)
		return
	}

	// send invalid proposal with the planed data.
	p := message.NewPropose(behaviours.r, behaviours.h, -1, behaviours.invalidValue, c.Backend().Sign)
	c.BroadcastAll(p)
}
