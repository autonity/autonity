package malicious

import (
	"context"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"github.com/autonity/autonity/node"
	"github.com/autonity/autonity/test"
	"github.com/stretchr/testify/require"
	"testing"
)

type malPrecommitService struct {
	*core.Core
	interfaces.Precommiter
}

func (c *malPrecommitService) SendPrecommit(ctx context.Context, isNil bool) {
	logger := c.Logger().New("step", c.Step())

	var precommit = messageutils.Vote{
		Round:  c.Round(),
		Height: c.Height(),
	}

	if isNil {
		precommit.ProposedBlockHash = common.Hash{}
	} else {
		if h := c.CurRoundMessages().GetProposalHash(); h == (common.Hash{}) {
			c.Logger().Error("core.sendPrecommit Proposal is empty! It should not be empty!")
			return
		}
		precommit.ProposedBlockHash = c.CurRoundMessages().GetProposalHash()
	}

	encodedVote, err := messageutils.Encode(&precommit)
	if err != nil {
		logger.Error("Failed to encode", "subject", precommit)
		return
	}

	msg := &messageutils.Message{
		Code:          messageutils.MsgPrecommit,
		Msg:           encodedVote,
		Address:       c.Address(),
		CommittedSeal: []byte{},
	}

	// Create committed seal
	//	seal := helpers.PrepareCommittedSeal(precommit.ProposedBlockHash, c.Round(), c.Height())
	//	msg.CommittedSeal, err = c.Backend().Sign(seal)
	// nil committed seal
	msg.CommittedSeal = nil
	if err != nil {
		c.Logger().Error("core.sendPrecommit error while signing committed seal", "err", err)
	}

	c.SetSentPrecommit(true)
	c.Br().Broadcast(ctx, msg)
}

func TestMaliciousPrecommitSender(t *testing.T) {
	users, err := test.Users(6, "10e18,v,100,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)

	//set Malicious users
	users[0].CustHandler = &node.CustomHandler{Precommitter: &malPrecommitService{}}
	// creates a network of 6 users and starts all the nodes in it
	network, err := test.NewNetworkFromUsers(users, true)
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120)
	require.NoError(t, err)
	defer network.Shutdown()
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}
