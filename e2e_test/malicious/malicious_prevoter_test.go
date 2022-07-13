package malicious

import (
	"context"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"github.com/autonity/autonity/consensus/tendermint/core/types"
	"github.com/autonity/autonity/node"
	"github.com/autonity/autonity/test"
	"github.com/stretchr/testify/require"
	"testing"
)

type malPrevoter struct {
	*core.Core
	interfaces.Prevoter
}

//HandlePrevote overrides core.HandlePrevote, It accepts a vote and sends a precommit without checking
// for 2f+1 vote count
func (c *malPrevoter) HandlePrevote(ctx context.Context, msg *messageutils.Message) error {
	var preVote messageutils.Vote
	err := msg.Decode(&preVote)
	if err != nil {
		return constants.ErrFailedDecodePrevote
	}

	prevoteHash := preVote.ProposedBlockHash
	c.AcceptVote(c.CurRoundMessages(), types.Prevote, prevoteHash, *msg)

	// Now we can add the preVote to our current round state
	if err := c.PrevoteTimeout().StopTimer(); err != nil {
		return err
	}
	c.Logger().Debug("Stopped Scheduled Prevote Timeout")

	c.GetPrecommiter().SendPrecommit(ctx, true)
	c.SetStep(types.Precommit)

	return nil
}

func TestMaliciousPrevoter(t *testing.T) {
	users, err := test.Users(6, "10e18,v,100,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)

	//set Malicious users
	users[0].CustHandler = &node.CustomHandler{Prevoter: &malPrevoter{}}
	// creates a network of 6 users and starts all the nodes in it
	network, err := test.NewNetworkFromUsers(users, true)
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120)
	require.NoError(t, err)
	defer network.Shutdown()
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}
