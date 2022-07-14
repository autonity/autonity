package malicious

import (
	"context"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/helpers"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/node"
	"github.com/autonity/autonity/test"
	"github.com/stretchr/testify/require"
	"testing"
)

type preVoteSpammer struct {
	*core.Core
	interfaces.Prevoter
}

func (c *preVoteSpammer) SendPrevote(ctx context.Context, isNil bool) {
	logger := c.Logger().New("step", c.Step())

	var prevote = messageutils.Vote{
		Round:  c.Round(),
		Height: c.Height(),
	}

	if isNil {
		prevote.ProposedBlockHash = common.Hash{}
	} else {
		if h := c.CurRoundMessages().GetProposalHash(); h == (common.Hash{}) {
			c.Logger().Error("sendPrecommit Proposal is empty! It should not be empty!")
			return
		}
		prevote.ProposedBlockHash = c.CurRoundMessages().GetProposalHash()
	}

	encodedVote, err := messageutils.Encode(&prevote)
	if err != nil {
		logger.Error("Failed to encode", "subject", prevote)
		return
	}

	msg := &messageutils.Message{
		Code:          messageutils.MsgPrevote,
		Msg:           encodedVote,
		Address:       c.Address(),
		CommittedSeal: []byte{},
	}
	for i := 0; i < 4000; i++ {
		c.Br().Broadcast(ctx, msg)
	}
	c.SetSentPrevote(true)
}

// TestPrevoteSpammer spams the network by broadcasting 4k preovte messages at once
func TestPrevoteSpammer(t *testing.T) {
	users, err := test.Validators(6, "10e18,v,100,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)

	//set Malicious users
	users[0].CustHandler = &node.CustomHandler{Prevoter: &preVoteSpammer{}}
	users[1].CustHandler = &node.CustomHandler{Prevoter: &preVoteSpammer{}}
	// creates a network of 6 users and starts all the nodes in it
	network, err := test.NewNetworkFromValidators(users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

type precommitSpammer struct {
	*core.Core
	interfaces.Precommiter
}

func (c *precommitSpammer) SendPrecommit(ctx context.Context, isNil bool) {
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
	seal := helpers.PrepareCommittedSeal(precommit.ProposedBlockHash, c.Round(), c.Height())
	msg.CommittedSeal, err = c.Backend().Sign(seal)
	if err != nil {
		c.Logger().Error("core.sendPrecommit error while signing committed seal", "err", err)
	}

	for i := 0; i < 4000; i++ {
		c.Br().Broadcast(ctx, msg)
	}
	c.SetSentPrecommit(true)
}

// TestPrecommitSpammer spams the network by broadcasting 4k precommit messages at once
func TestPrecommitSpammer(t *testing.T) {
	users, err := test.Validators(6, "10e18,v,100,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)

	//set Malicious users
	users[0].CustHandler = &node.CustomHandler{Precommitter: &precommitSpammer{}}
	// creates a network of 6 users and starts all the nodes in it
	network, err := test.NewNetworkFromValidators(users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

type proposalSpammer struct {
	*core.Core
	interfaces.Proposer
}

func (c *proposalSpammer) SendProposal(ctx context.Context, p *types.Block) {

	proposalBlock := messageutils.NewProposal(c.Round(), c.Height(), c.ValidRound(), p)
	proposal, _ := messageutils.Encode(proposalBlock)

	c.SetSentProposal(true)
	c.Backend().SetProposedBlockHash(p.Hash())

	for i := 0; i < 4000; i++ {
		c.Br().Broadcast(ctx, &messageutils.Message{
			Code:          messageutils.MsgProposal,
			Msg:           proposal,
			Address:       c.Address(),
			CommittedSeal: []byte{},
		})
	}
}

// TestProposalSpammer spams the network by broadcasting 4k proposal messages at once
func TestProposalSpammer(t *testing.T) {
	users, err := test.Validators(6, "10e18,v,100,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)

	//set Malicious proposalSender
	users[0].CustHandler = &node.CustomHandler{Proposer: &proposalSpammer{}}
	// creates a network of 6 users and starts all the nodes in it
	network, err := test.NewNetworkFromValidators(users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	err = network.WaitForSyncComplete()
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}
