package byzantine

import (
	"context"
	"testing"

	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/core/types"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/node"
	"github.com/stretchr/testify/require"
)

func newMalPrevoter(c interfaces.Tendermint) interfaces.Prevoter {
	return &malPrevoter{c.(*core.Core), c.Prevoter()}
}

type malPrevoter struct {
	*core.Core
	interfaces.Prevoter
}

// HandlePrevote overrides core.HandlePrevote, It accepts a vote and sends a precommit without checking
// for 2f+1 vote count
func (c *malPrevoter) HandlePrevote(ctx context.Context, msg *message.Message) error {
	var preVote message.Vote
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

	c.Precommiter().SendPrecommit(ctx, true)
	c.SetStep(types.Precommit)

	return nil
}

func TestMaliciousPrevoter(t *testing.T) {
	users, err := e2e.Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	//set Malicious users
	users[0].TendermintServices = &node.TendermintServices{Prevoter: newMalPrevoter}
	// creates a network of 6 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}
