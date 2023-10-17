package byzantine

import (
	"context"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/node"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
	"math/big"
	"sync/atomic"
	"testing"
)

type duplicateProposalSender struct {
	*core.Core
	interfaces.Proposer
}

// SendProposal overrides core.sendProposal and send multiple proposals
func (c *duplicateProposalSender) SendProposal(ctx context.Context, p *types.Block) {
	proposal := message.NewPropose(c.Round(), c.Height().Uint64(), c.ValidRound(), p, c.Backend().Sign)
	proposal2 := message.NewPropose(c.Round(), c.Height().Uint64(), c.ValidRound()-1, p, c.Backend().Sign)

	c.SetSentProposal(true)
	c.Backend().SetProposedBlockHash(p.Hash())
	//send same proposal twice
	c.Br().Broadcast(ctx, proposal)
	// send 2nd proposal with different validround
	c.Br().Broadcast(ctx, proposal2)
}

// TestDuplicateProposal broadcasts two proposals with same round and same height but different validround
func TestDuplicateProposal(t *testing.T) {
	users, err := e2e.Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	//set Malicious proposalSender
	users[0].TendermintServices = &node.TendermintServices{Proposer: newDuplicateProposalSender}
	// creates a network of 6 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	err = network.WaitForSyncComplete()
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

func newMalProposalSender(c interfaces.Tendermint) interfaces.Broadcaster {
	return &malProposalSender{c.(*core.Core)}
}

type malProposalSender struct {
	*core.Core
}

func (c *malProposalSender) Broadcast(ctx context.Context, msg message.Message) {
	round := msg.R()
	height := msg.H()
	// if we are the proposer for this round, return
	if c.CommitteeSet().GetProposer(round).Address == c.Backend().Address() {
		return
	}
	header := &types.Header{Number: new(big.Int).SetUint64(height)}
	block := types.NewBlockWithHeader(header)
	// create a new proposal message
	propose := message.NewPropose(round, height, -1, block, c.Backend().Sign)
	c.BroadcastAll(ctx, propose)
}

func newProposalApprover(c interfaces.Tendermint) interfaces.Proposer {
	return &proposalApprover{c.(*core.Core), c.Proposer()}
}

type proposalApprover struct {
	*core.Core
	interfaces.Proposer
}

func (c *proposalApprover) HandleProposal(ctx context.Context, proposal *message.Propose) error {
	// Set the proposal for the current round
	c.CurRoundMessages().SetProposal(proposal, true)
	c.GetPrevoter().SendPrevote(ctx, false)
	c.SetStep(core.Prevote)
	return nil
}

func TestNonProposerWithFaultyApprover(t *testing.T) {
	t.Skip("a malicious proposer will be kicked out from network now, this test is no more valid")
	users, err := e2e.Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	//set Malicious proposalSender
	users[0].TendermintServices = &node.TendermintServices{Broadcaster: newMalProposalSender}
	users[1].TendermintServices = &node.TendermintServices{Broadcaster: newMalProposalSender, Proposer: newProposalApprover}
	users[2].TendermintServices = &node.TendermintServices{Broadcaster: newMalProposalSender}
	// creates a network of 6 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	err = network.WaitForSyncComplete()
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

func TestDuplicateProposalWithFaultyApprover(t *testing.T) {
	users, err := e2e.Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	//set Malicious proposalSender
	users[0].TendermintServices = &node.TendermintServices{Proposer: newDuplicateProposalSender}
	users[1].TendermintServices = &node.TendermintServices{Proposer: newProposalApprover}
	// creates a network of 6 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	err = network.WaitForSyncComplete()
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

func newPartialProposalSender(c interfaces.Tendermint) interfaces.Proposer {
	return &partialProposalSender{c.(*core.Core), c.Proposer()}
}

type partialProposalSender struct {
	*core.Core
	interfaces.Proposer
}

// SendProposal overrides core.sendProposal and send multiple proposals
func (c *partialProposalSender) SendProposal(ctx context.Context, p *types.Block) {
	fakeTransactions := make([]*types.Transaction, 0)
	for i := 0; i < 5; i++ {
		var fakeTransaction types.Transaction
		f := fuzz.New()
		f.Fuzz(&fakeTransaction)
		var tx types.LegacyTx
		f.Fuzz(&tx)
		fakeTransaction.SetInner(&tx)
		fakeTransactions = append(fakeTransactions, &fakeTransaction)
	}
	p.SetTransactions(fakeTransactions)
	proposal := message.NewPropose(c.Round(), c.Height().Uint64(), c.ValidRound(), p, c.Backend().Sign)
	c.SetSentProposal(true)
	c.Backend().SetProposedBlockHash(p.Hash())
	//send same proposal twice
	c.Br().Broadcast(ctx, proposal)
}

func TestPartialProposal(t *testing.T) {
	users, err := e2e.Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	//set Malicious proposalSender
	users[0].TendermintServices = &node.TendermintServices{Proposer: newPartialProposalSender}
	// creates a network of 6 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	err = network.WaitForSyncComplete()
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

func newInvalidBlockProposer(c interfaces.Tendermint) interfaces.Proposer {
	return &invalidBlockProposer{c.(*core.Core), c.Proposer()}
}

type invalidBlockProposer struct {
	*core.Core
	interfaces.Proposer
}

// SendProposal overrides core.sendProposal and send multiple proposals
func (c *invalidBlockProposer) SendProposal(ctx context.Context, p *types.Block) {
	fakeTransactions := make([]*types.Transaction, 0)
	f := fuzz.New()
	for i := 0; i < 5; i++ {
		var fakeTransaction types.Transaction
		f.Fuzz(&fakeTransaction)
		var tx types.LegacyTx
		f.Fuzz(&tx)
		fakeTransaction.SetInner(&tx)

		fakeTransactions = append(fakeTransactions, &fakeTransaction)
	}
	p.SetTransactions(fakeTransactions)
	var hash common.Hash
	f.Fuzz(&hash)
	var atmHash atomic.Value
	atmHash.Store(hash)
	// nil hash
	p.SetHash(atmHash)

	// nil header
	var num big.Int
	f.Fuzz(&num)
	p.SetHeaderNumber(&num)
	proposal := message.NewPropose(c.Round(), c.Height().Uint64(), c.ValidRound(), p, c.Backend().Sign)

	c.SetSentProposal(true)
	c.Backend().SetProposedBlockHash(p.Hash())

	//send same proposal twice
	c.BroadcastAll(ctx, proposal)
}

func TestInvalidBlockProposal(t *testing.T) {
	//for i := 0; i < 20; i++ {
	users, err := e2e.Validators(t, 4, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	//set Malicious proposalSender
	users[0].TendermintServices = &node.TendermintServices{Proposer: newInvalidBlockProposer}
	// creates a network of 6 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)

	err = network.WaitForSyncComplete()
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(5, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
	network.Shutdown()
	//}
}
