package core

import (
	"context"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/log"
)

func TestSendPrecommit(t *testing.T) {
	t.Run("proposal is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any()).Times(0)

		messages := message.NewMap()
		c := &Core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: messages.GetOrCreate(0),
			round:            2,
			height:           big.NewInt(3),
		}
		c.SetDefaultHandlers()
		c.precommiter.SendPrecommit(context.Background(), false)
	})

	t.Run("valid proposal given, non nil pre-commit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet, keys := NewTestCommitteeSetWithKeys(7)
		me, _ := committeeSet.GetByIndex(0)
		logger := log.New("backend", "test", "id", 0)
		val, _ := committeeSet.GetByIndex(2)
		addr := val.Address

		proposal := message.NewPropose(
			1,
			2,
			1,
			types.NewBlockWithHeader(&types.Header{}),
			makeSigner(keys[me.Address], me.Address),
		)

		messages := message.NewMap()
		curRoundMessages := messages.GetOrCreate(1)
		curRoundMessages.SetProposal(proposal, false)

		preCommit := message.NewPrecommit(1, 2, curRoundMessages.ProposalHash(), makeSigner(keys[addr], addr))
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Broadcast(gomock.Any(), preCommit)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(makeSigner(keys[addr], addr))

		c := &Core{
			backend:          backendMock,
			address:          addr,
			logger:           logger,
			committee:        committeeSet,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			round:            1,
			height:           big.NewInt(2),
		}
		c.SetDefaultHandlers()
		c.precommiter.SendPrecommit(context.Background(), false)
	})

	t.Run("valid proposal given, nil pre-commit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet, keys := NewTestCommitteeSetWithKeys(7)
		me, _ := committeeSet.GetByIndex(0)
		logger := log.New("backend", "test", "id", 0)
		val, _ := committeeSet.GetByIndex(2)
		addr := val.Address

		proposal := message.NewPropose(
			1,
			2,
			1,
			types.NewBlockWithHeader(&types.Header{}),
			makeSigner(keys[me.Address], me.Address))

		messages := message.NewMap()
		curRoundMessages := messages.GetOrCreate(1)
		curRoundMessages.SetProposal(proposal, true)

		preCommit := message.NewPrecommit(1, 2, common.Hash{}, makeSigner(keys[addr], addr))
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Broadcast(gomock.Any(), preCommit)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(makeSigner(keys[addr], addr))

		c := &Core{
			backend:          backendMock,
			address:          addr,
			logger:           logger,
			curRoundMessages: curRoundMessages,
			messages:         messages,
			committee:        committeeSet,
			height:           big.NewInt(2),
			round:            1,
		}

		c.SetDefaultHandlers()
		c.precommiter.SendPrecommit(context.Background(), true)
	})

}

func TestHandlePrecommit(t *testing.T) {
	t.Run("pre-commit with invalid signature given, panic", func(t *testing.T) {
		committeeSet, keys := NewTestCommitteeSetWithKeys(4)
		member, _ := committeeSet.GetByIndex(1)
		messages := message.NewMap()
		curRoundMessages := messages.GetOrCreate(2)
		preCommit := message.NewPrecommit(2, 3, curRoundMessages.ProposalHash(), makeSigner(keys[member.Address], member.Address))

		c := &Core{
			address:          member.Address,
			round:            2,
			height:           big.NewInt(3),
			curRoundMessages: curRoundMessages,
			logger:           log.New("backend", "test", "id", 0),
			proposeTimeout:   NewTimeout(Propose, log.New("ProposeTimeout")),
			prevoteTimeout:   NewTimeout(Prevote, log.New("PrevoteTimeout")),
			precommitTimeout: NewTimeout(Precommit, log.New("PrecommitTimeout")),
		}
		c.SetDefaultHandlers()
		c.SetStep(context.Background(), Precommit)
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		c.precommiter.HandlePrecommit(context.Background(), preCommit)

	})

	t.Run("pre-commit given with no errors, commit called", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet, keys := NewTestCommitteeSetWithKeys(1)
		member, _ := committeeSet.GetByIndex(0)
		logger := log.New("backend", "test", "id", 0)

		proposal := message.NewPropose(
			2,
			3,
			1,
			types.NewBlockWithHeader(&types.Header{}),
			makeSigner(keys[member.Address], member.Address))

		messages := message.NewMap()
		curRoundMessages := messages.GetOrCreate(2)
		curRoundMessages.SetProposal(proposal, true)

		msg := message.NewPrecommit(2, 3, proposal.Block().Hash(), makeSigner(keys[member.Address], member.Address))
		msg.MustVerify(stubVerifier)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Commit(proposal.Block(), gomock.Any(), gomock.Any()).Return(nil).Do(
			func(proposalBlock *types.Block, round int64, seals [][]byte) {
				if round != 2 {
					t.Fatal("Commit called with round different than precommit seal")
				}
				if !reflect.DeepEqual([][]byte{msg.Signature()}, seals) {
					t.Fatal("Commit called with wrong seal")
				}
			})

		c := &Core{
			address:          member.Address,
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			logger:           logger,
			round:            2,
			height:           big.NewInt(3),
			step:             Precommit,
			committee:        committeeSet,
			proposeTimeout:   NewTimeout(Propose, logger),
			prevoteTimeout:   NewTimeout(Prevote, logger),
			precommitTimeout: NewTimeout(Precommit, logger),
		}

		c.SetDefaultHandlers()
		err := c.precommiter.HandlePrecommit(context.Background(), msg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})

	t.Run("pre-commit given with no errors, pre-commit Timeout triggered", func(t *testing.T) {
		logger := log.New("backend", "test", "id", 0)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		backendMock := interfaces.NewMockBackend(ctrl)
		committeeSet, keys := NewTestCommitteeSetWithKeys(7)
		me, _ := committeeSet.GetByIndex(0)
		proposal := message.NewPropose(
			2,
			3,
			1,
			types.NewBlockWithHeader(&types.Header{}),
			makeSigner(keys[me.Address], me.Address))

		messages := message.NewMap()
		curRoundMessages := messages.GetOrCreate(2)
		curRoundMessages.SetProposal(proposal, true)

		c := &Core{
			address:          me.Address,
			backend:          backendMock,
			curRoundMessages: curRoundMessages,
			messages:         messages,
			logger:           logger,
			round:            2,
			height:           big.NewInt(3),
			step:             Precommit,
			committee:        committeeSet,
			precommitTimeout: NewTimeout(Precommit, logger),
		}
		c.SetDefaultHandlers()
		backendMock.EXPECT().Post(gomock.Any()).Times(1)

		for _, member := range committeeSet.Committee()[1:5] {
			msg := message.NewPrecommit(2, 3, proposal.Block().Hash(), makeSigner(keys[member.Address], member.Address))
			if err := c.precommiter.HandlePrecommit(context.Background(), msg.MustVerify(stubVerifier)); err != nil {
				t.Fatalf("Expected nil, got %v", err)
			}
		}

		msg := message.NewPrecommit(2, 3, common.Hash{}, makeSigner(keys[me.Address], me.Address))
		if err := c.precommiter.HandlePrecommit(context.Background(), msg.MustVerify(stubVerifier)); err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		<-time.NewTimer(5 * time.Second).C

	})
}

func TestHandleCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := log.New("backend", "test", "id", 0)

	addr := common.HexToAddress("0x0123456789")
	testCommittee, keys := GenerateCommittee(3)

	firstKey := keys[testCommittee[0].Address]

	h := &types.Header{Number: big.NewInt(3)}

	// Sign the header so that types.ECRecover works
	seal, err := crypto.Sign(crypto.Keccak256(types.SigHash(h).Bytes()), firstKey)
	require.NoError(t, err)

	err = types.WriteSeal(h, seal)
	require.NoError(t, err)

	h.Committee = testCommittee

	block := types.NewBlockWithHeader(h)
	testCommittee = append(testCommittee, types.CommitteeMember{Address: addr, VotingPower: big.NewInt(1)})
	committeeSet, err := committee.NewRoundRobinSet(testCommittee, testCommittee[0].Address)
	require.NoError(t, err)

	backendMock := interfaces.NewMockBackend(ctrl)
	backendMock.EXPECT().HeadBlock().MinTimes(1).Return(block)

	c := &Core{
		address:          addr,
		backend:          backendMock,
		round:            2,
		height:           big.NewInt(3),
		messages:         message.NewMap(),
		logger:           logger,
		proposeTimeout:   NewTimeout(Propose, logger),
		prevoteTimeout:   NewTimeout(Prevote, logger),
		precommitTimeout: NewTimeout(Precommit, logger),
		committee:        committeeSet,
	}
	c.SetDefaultHandlers()
	c.precommiter.HandleCommit(context.Background())
	if c.round != 0 || c.height.Cmp(big.NewInt(4)) != 0 {
		t.Fatalf("Expected new round")
	}
	// to fix the data race detected by CI workflow.
	err = c.proposeTimeout.StopTimer()
	if err != nil {
		t.Error(err)
	}
}
