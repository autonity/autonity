package core

import (
	"context"
	"math/big"
	"reflect"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
)

func TestSendPrecommit(t *testing.T) {
	t.Run("proposal is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any()).Times(0)

		c := &Core{
			logger:      log.New("backend", "test", "id", 0),
			backend:     backendMock,
			roundsState: newTendermintState(log.New(), nil, nil),
		}
		c.SetHeight(common.Big3)
		c.SetRound(2)

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
			makeSigner(keys[me.Address].consensus),
			&me,
		)

		backendMock := interfaces.NewMockBackend(ctrl)
		roundStates := newTendermintState(log.New(), nil, nil)
		c := &Core{
			backend:     backendMock,
			address:     addr,
			logger:      logger,
			committee:   committeeSet,
			roundsState: roundStates,
			lastHeader:  &types.Header{Committee: committeeSet.Committee()},
		}
		c.SetHeight(common.Big2)
		c.SetRound(1)
		curRoundMessages := roundStates.GetOrCreate(1)
		curRoundMessages.SetProposal(proposal, false)

		preCommit := message.NewPrecommit(1, 2, curRoundMessages.ProposalHash(), makeSigner(keys[addr].consensus), &val, 7)
		backendMock.EXPECT().Broadcast(gomock.Any(), preCommit)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(makeSigner(keys[addr].consensus))

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
			makeSigner(keys[me.Address].consensus),
			&me)

		backendMock := interfaces.NewMockBackend(ctrl)
		roundStates := newTendermintState(log.New(), nil, nil)
		c := &Core{
			backend:     backendMock,
			address:     addr,
			logger:      logger,
			roundsState: roundStates,
			committee:   committeeSet,
			lastHeader:  &types.Header{Committee: committeeSet.Committee()},
		}
		c.SetHeight(common.Big2)
		c.SetRound(1)
		curRoundMessages := roundStates.GetOrCreate(1)
		curRoundMessages.SetProposal(proposal, true)
		preCommit := message.NewPrecommit(1, 2, common.Hash{}, makeSigner(keys[addr].consensus), &val, 7)
		backendMock.EXPECT().Broadcast(gomock.Any(), preCommit)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(makeSigner(keys[addr].consensus))

		c.SetDefaultHandlers()
		c.precommiter.SendPrecommit(context.Background(), true)
	})

}

func TestHandlePrecommit(t *testing.T) {
	t.Run("pre-commit with invalid signature given, panic", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		committeeSet, keys := NewTestCommitteeSetWithKeys(4)
		member, _ := committeeSet.GetByIndex(1)
		roundStates := newTendermintState(log.New(), nil, nil)
		curRoundMessages := roundStates.GetOrCreate(2)
		preCommit := newUnverifiedPrecommit(2, 3, curRoundMessages.ProposalHash(), makeSigner(keys[member.Address].consensus), &member, 4)
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).MaxTimes(1)
		c := &Core{
			backend:          backendMock,
			address:          member.Address,
			roundsState:      roundStates,
			logger:           log.New("backend", "test", "id", 0),
			proposeTimeout:   NewTimeout(Propose, log.New("ProposeTimeout")),
			prevoteTimeout:   NewTimeout(Prevote, log.New("PrevoteTimeout")),
			precommitTimeout: NewTimeout(Precommit, log.New("PrecommitTimeout")),
			committee:        committeeSet,
		}
		c.SetRound(2)
		c.SetHeight(common.Big3)

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

		proposal := generateBlockProposal(2, big.NewInt(3), 1, false, makeSigner(keys[member.Address].consensus), &member)

		roundStates := newTendermintState(log.New(), nil, nil)
		msg := message.NewPrecommit(2, 3, proposal.Block().Hash(), makeSigner(keys[member.Address].consensus), &member, 1)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).MaxTimes(1)
		backendMock.EXPECT().Commit(proposal.Block(), gomock.Any(), gomock.Any()).Return(nil).Do(
			func(proposalBlock *types.Block, round int64, quorumCertificate types.AggregateSignature) {
				if round != 2 {
					t.Fatal("Commit called with round different than precommit seal")
				}

				expectedQuorumCertificate := types.NewAggregateSignature(msg.Signature().Copy(), msg.Signers().Copy())
				if !reflect.DeepEqual(expectedQuorumCertificate, quorumCertificate) {
					t.Fatal("Commit called with wrong seal")
				}
			})

		c := &Core{
			address:          member.Address,
			backend:          backendMock,
			roundsState:      roundStates,
			logger:           logger,
			committee:        committeeSet,
			proposeTimeout:   NewTimeout(Propose, logger),
			prevoteTimeout:   NewTimeout(Prevote, logger),
			precommitTimeout: NewTimeout(Precommit, logger),
		}
		c.SetHeight(common.Big3)
		curRoundMessages := roundStates.GetOrCreate(2)
		curRoundMessages.SetProposal(proposal, true)
		c.UpdateStep(Precommit)
		c.SetRound(2)

		c.SetDefaultHandlers()
		err := c.precommiter.HandlePrecommit(context.Background(), msg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})
	t.Run("quorum pre-commit given with no errors, pre-commit Timeout triggered", func(t *testing.T) {
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
			makeSigner(keys[me.Address].consensus),
			&me)

		roundStates := newTendermintState(log.New(), nil, nil)
		c := &Core{
			address:          me.Address,
			backend:          backendMock,
			logger:           logger,
			roundsState:      roundStates,
			committee:        committeeSet,
			precommitTimeout: NewTimeout(Precommit, logger),
		}
		c.SetHeight(common.Big3)
		c.SetRound(2)
		c.UpdateStep(Precommit)
		curRoundMessages := roundStates.GetOrCreate(2)
		curRoundMessages.SetProposal(proposal, true)

		c.SetDefaultHandlers()
		backendMock.EXPECT().Post(gomock.Any()).Times(6)

		for _, member := range committeeSet.Committee()[1:5] {
			m := member
			msg := message.NewPrecommit(2, 3, proposal.Block().Hash(), makeSigner(keys[member.Address].consensus), &m, 7)
			if err := c.precommiter.HandlePrecommit(context.Background(), msg); err != nil {
				t.Fatalf("Expected nil, got %v", err)
			}
		}

		msg := message.NewPrecommit(2, 3, common.Hash{}, makeSigner(keys[me.Address].consensus), &me, 7)
		if err := c.precommiter.HandlePrecommit(context.Background(), msg); err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		<-time.NewTimer(5 * time.Second).C

	})
}

func TestHandleCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer waitForExpects(ctrl)

	logger := log.New("backend", "test", "id", 0)

	testCommittee, _ := GenerateCommittee(3)

	h := &types.Header{Number: big.NewInt(3)}
	block := types.NewBlockWithHeader(h)
	committeeSet, err := committee.NewRoundRobinSet(testCommittee, testCommittee[0].Address)
	require.NoError(t, err)

	backendMock := interfaces.NewMockBackend(ctrl)
	backendMock.EXPECT().HeadBlock().MinTimes(1).Return(block)
	backendMock.EXPECT().Post(gomock.Any()).MaxTimes(1)
	backendMock.EXPECT().ProcessFutureMsgs(uint64(4)).MaxTimes(1)

	c := &Core{
		address:          testCommittee[0].Address,
		backend:          backendMock,
		roundsState:      newTendermintState(log.New(), nil, nil),
		logger:           logger,
		proposeTimeout:   NewTimeout(Propose, logger),
		prevoteTimeout:   NewTimeout(Prevote, logger),
		precommitTimeout: NewTimeout(Precommit, logger),
		committee:        committeeSet,
	}
	c.SetRound(2)
	c.SetHeight(common.Big3)

	c.SetDefaultHandlers()
	c.precommiter.HandleCommit(context.Background())
	if c.Round() != 0 || c.Height().Cmp(big.NewInt(4)) != 0 {
		t.Fatalf("Expected new round and new height")
	}
	// to fix the data race detected by CI workflow.
	err = c.proposeTimeout.StopTimer()
	if err != nil {
		t.Error(err)
	}
}
