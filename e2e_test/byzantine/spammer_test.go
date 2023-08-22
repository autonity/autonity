package byzantine

import (
	"context"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/helpers"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/node"
	"github.com/autonity/autonity/rlp"

	"github.com/stretchr/testify/require"
	"testing"
)

type preVoteSpammer struct {
	*core.Core
	interfaces.Prevoter
}

func (c *preVoteSpammer) SendPrevote(ctx context.Context, isNil bool) {
	logger := c.Logger().New("step", c.Step())

	var prevote = message.Vote{
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

	encodedVote, err := rlp.EncodeToBytes(&prevote)
	if err != nil {
		logger.Error("Failed to encode", "subject", prevote)
		return
	}

	msg := &message.Message{
		Code:          consensus.MsgPrevote,
		Payload:       encodedVote,
		Address:       c.Address(),
		CommittedSeal: []byte{},
	}
	for i := 0; i < 1000; i++ {
		c.Br().SignAndBroadcast(ctx, msg)
	}
	c.SetSentPrevote(true)
}

// TestPrevoteSpammer spams the network by broadcasting 4k preovte messages at once
func TestPrevoteSpammer(t *testing.T) {
	users, err := e2e.Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	//set Malicious users
	users[0].TendermintServices = &node.TendermintServices{Prevoter: &preVoteSpammer{}}
	users[1].TendermintServices = &node.TendermintServices{Prevoter: &preVoteSpammer{}}
	// creates a network of 6 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
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

	// Create committed seal
	seal := helpers.PrepareCommittedSeal(precommit.ProposedBlockHash, c.Round(), c.Height())
	msg.CommittedSeal, err = c.Backend().Sign(seal)
	if err != nil {
		c.Logger().Error("core.sendPrecommit error while signing committed seal", "err", err)
	}

	for i := 0; i < 1000; i++ {
		c.Br().SignAndBroadcast(ctx, msg)
	}
	c.SetSentPrecommit(true)
}

// TestPrecommitSpammer spams the network by broadcasting 4k precommit messages at once
func TestPrecommitSpammer(t *testing.T) {
	users, err := e2e.Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	//set Malicious users
	users[0].TendermintServices = &node.TendermintServices{Precommitter: &precommitSpammer{}}
	// creates a network of 6 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
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
	proposalBlock := message.NewProposal(c.Round(), c.Height(), c.ValidRound(), p, c.Backend().Sign)
	proposal, _ := rlp.EncodeToBytes(proposalBlock)

	c.SetSentProposal(true)
	c.Backend().SetProposedBlockHash(p.Hash())

	for i := 0; i < 1000; i++ {
		c.Br().SignAndBroadcast(ctx, &message.Message{
			Code:          consensus.MsgProposal,
			Payload:       proposal,
			Address:       c.Address(),
			CommittedSeal: []byte{},
		})
	}
}

// TestProposalSpammer spams the network by broadcasting 4k proposal messages at once
func TestProposalSpammer(t *testing.T) {
	users, err := e2e.Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	//set Malicious proposalSender
	users[0].TendermintServices = &node.TendermintServices{Proposer: &proposalSpammer{}}
	// creates a network of 6 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	err = network.WaitForSyncComplete()
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}
