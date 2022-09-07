package malicious

import (
	"context"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	tctypes "github.com/autonity/autonity/consensus/tendermint/core/types"
	"github.com/autonity/autonity/core/types"
	e2etest "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/node"
	"github.com/autonity/autonity/test"
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

	proposalBlock := messageutils.NewProposal(c.Round(), c.Height(), c.ValidRound(), p)
	proposalBlock2 := messageutils.NewProposal(c.Round(), c.Height(), c.ValidRound()-1, p)
	proposal, _ := messageutils.Encode(proposalBlock)
	proposal2, _ := messageutils.Encode(proposalBlock2)

	c.SetSentProposal(true)
	c.Backend().SetProposedBlockHash(p.Hash())

	//send same proposal twice
	c.Br().Broadcast(ctx, &messageutils.Message{
		Code:          messageutils.MsgProposal,
		Msg:           proposal,
		Address:       c.Address(),
		CommittedSeal: []byte{},
	})
	// send 2nd proposal with different validround
	c.Br().Broadcast(ctx, &messageutils.Message{
		Code:          messageutils.MsgProposal,
		Msg:           proposal2,
		Address:       c.Address(),
		CommittedSeal: []byte{},
	})
}

// TestDuplicateProposal broadcasts two proposals with same round and same height but different validround
func TestDuplicateProposal(t *testing.T) {
	users, err := test.Validators(6, "10e18,v,100,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)

	//set Malicious proposalSender
	users[0].CustHandler = &node.CustomHandler{Proposer: &duplicateProposalSender{}}
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

type malProposalSender struct {
	*core.Core
}

// SendProposalFromNonProposer broadcasts a new proposal, only if it is a non-proposer
func SendProposalFromNonProposer(ctx context.Context, c *core.Core, fm []byte) {

	m := new(messageutils.Message)
	if err := m.FromPayload(fm); err != nil {
		c.Logger().Error("can not send proposal, invalid payload", "err", err)
	}

	round, err := m.Round()
	if err != nil {
		c.Logger().Error("failed to decode round value, err:", err)
		return
	}

	height, err := m.Height()
	if err != nil {
		c.Logger().Error("failed to decode height value, err:", err)
		return
	}

	// if we are the proposer for this round, return
	if c.CommitteeSet().GetProposer(round).Address == c.Backend().Address() {
		return
	}
	header := &types.Header{Number: new(big.Int).SetUint64(height.Uint64())}
	block := types.NewBlockWithHeader(header)
	// create a new proposal message
	msgP := e2etest.MsgPropose(c.Backend().Address(), block, height.Uint64(), round, -1)
	fm, err = c.FinalizeMessage(msgP)
	if err != nil {
		return
	}

	if err := c.Backend().Broadcast(ctx, c.CommitteeSet().Committee(), fm); err != nil {
		c.Logger().Error("consensus message broadcast failure, err:", err)
	}
}

//Broadcast overrides the code.Broadcast
func (c *malProposalSender) Broadcast(ctx context.Context, msg *messageutils.Message) {
	logger := c.Logger().New("step", c.Step())

	fm, err := c.FinalizeMessage(msg)
	if err != nil {
		return
	}
	SendProposalFromNonProposer(ctx, c.Core, fm)
	if err := c.Backend().Broadcast(ctx, c.CommitteeSet().Committee(), fm); err != nil {
		logger.Error("consensus message broadcast failure, err:", err)
	}
}

type proposalApprover struct {
	*core.Core
	interfaces.Proposer
}

func (c *proposalApprover) HandleProposal(ctx context.Context, msg *messageutils.Message) error {
	var proposal messageutils.Proposal
	err := msg.Decode(&proposal)
	if err != nil {
		return constants.ErrFailedDecodeProposal
	}
	// Set the proposal for the current round
	c.CurRoundMessages().SetProposal(&proposal, msg, true)

	c.GetPrevoter().SendPrevote(ctx, false)
	c.SetStep(tctypes.Prevote)
	return nil
}

func TestNonProposerWithFaultyApprover(t *testing.T) {
	t.Skip("a malicious proposer will be kicked out from network now, this test is no more valid")
	users, err := test.Validators(6, "10e18,v,100,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)

	//set Malicious proposalSender
	users[0].CustHandler = &node.CustomHandler{Broadcaster: &malProposalSender{}}
	users[1].CustHandler = &node.CustomHandler{Broadcaster: &malProposalSender{}, Proposer: &proposalApprover{}}
	users[2].CustHandler = &node.CustomHandler{Broadcaster: &malProposalSender{}}
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

func TestDuplicateProposalWithFaultyApprover(t *testing.T) {
	users, err := test.Validators(6, "10e18,v,100,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)

	//set Malicious proposalSender
	users[0].CustHandler = &node.CustomHandler{Proposer: &duplicateProposalSender{}}
	users[1].CustHandler = &node.CustomHandler{Proposer: &proposalApprover{}}
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
	proposalBlock := messageutils.NewProposal(c.Round(), c.Height(), c.ValidRound(), p)
	proposal, _ := messageutils.Encode(proposalBlock)

	c.SetSentProposal(true)
	c.Backend().SetProposedBlockHash(p.Hash())

	//send same proposal twice
	c.Br().Broadcast(ctx, &messageutils.Message{
		Code:          messageutils.MsgProposal,
		Msg:           proposal,
		Address:       c.Address(),
		CommittedSeal: []byte{},
	})
}
func TestPartialProposal(t *testing.T) {
	users, err := test.Validators(6, "10e18,v,100,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)

	//set Malicious proposalSender
	users[0].CustHandler = &node.CustomHandler{Proposer: &partialProposalSender{}}
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
	proposalBlock := messageutils.NewProposal(c.Round(), c.Height(), c.ValidRound(), p)
	proposal, _ := messageutils.Encode(proposalBlock)

	c.SetSentProposal(true)
	c.Backend().SetProposedBlockHash(p.Hash())
	var nilAddr common.Address
	fuzz.New().Fuzz(&nilAddr)
	ranBytes, _ := e2etest.GenerateRandomBytes(10000000)

	// junk Address
	junkAddr := common.BytesToAddress(ranBytes)
	//send same proposal twice
	c.Br().Broadcast(ctx, &messageutils.Message{
		Code:          messageutils.MsgProposal,
		Msg:           proposal,
		Address:       junkAddr,
		CommittedSeal: []byte{},
	})
}

func TestInvalidBlockProposal(t *testing.T) {
	//for i := 0; i < 20; i++ {
	users, err := test.Validators(4, "10e18,v,100,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)

	//set Malicious proposalSender
	users[0].CustHandler = &node.CustomHandler{Proposer: &invalidBlockProposer{}}
	// creates a network of 6 users and starts all the nodes in it
	network, err := test.NewNetworkFromValidators(users, true)
	require.NoError(t, err)

	err = network.WaitForSyncComplete()
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(5, 60)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
	network.Shutdown()
	//}
}
