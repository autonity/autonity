package byzantine

import (
	"context"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/node"
	"github.com/autonity/autonity/rlp"
	"github.com/stretchr/testify/require"
	"testing"
)

type precommitService struct {
	*core.Core
	interfaces.Precommiter
}
type malPrecommitService struct {
	*core.Core
	interfaces.Precommiter
}

func (c *malPrecommitService) SendPrecommit(ctx context.Context, isNil bool) {
	logger := c.Logger().New("step", c.Step())

	var precommit = message.Vote{
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

	encodedVote, err := rlp.EncodeToBytes(&precommit)
	if err != nil {
		logger.Error("Failed to encode", "subject", precommit)
		return
	}

	msg := &message.Message{
		Code:          consensus.MsgPrecommit,
		Payload:       encodedVote,
		Address:       c.Address(),
		CommittedSeal: []byte{},
	}

	// nil committed seal
	msg.CommittedSeal = nil
	c.SetSentPrecommit(true)
	c.Br().SignAndBroadcast(ctx, msg)
}

func TestMaliciousPrecommitSender(t *testing.T) {
	users, err := e2e.Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	//set Malicious users
	users[0].TendermintServices = &node.TendermintServices{Precommitter: &malPrecommitService{}}
	// creates a network of 6 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

func TestMaliciousSenderDisc(t *testing.T) {
	users, err := e2e.Validators(t, 4, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)
	valPrecommiter1 := &precommitService{}
	valPrecommiter2 := &precommitService{}

	users[0].TendermintServices = &node.TendermintServices{Precommitter: valPrecommiter1}
	users[1].TendermintServices = &node.TendermintServices{Precommitter: valPrecommiter2}

	// creates a network of users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(1, 60)
	require.NoError(t, err)
	// set a malicious precommitter, this should cause one of the validator to be
	// evicted from the quorum
	valPrecommiter1.Core.SetPrecommitter(&malPrecommitService{})
	valPrecommiter2.Core.SetPrecommitter(&malPrecommitService{})

	// network should stop mining blocks
	err = network.WaitToMineNBlocks(10, 120)
	require.Error(t, err, "Network is not supposed to be mining blocks at this point")
}
