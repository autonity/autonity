package core

import (
	"context"
	"math/big"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
)

func TestSendPrevote(t *testing.T) {
	t.Run("proposal is empty and send prevote nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		committeeSet := NewTestCommitteeSet(4)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any()).Times(1)
		// return random signature just to allow prevote encoding
		backendMock.EXPECT().Sign(gomock.Any()).Times(1).Return(testSignature)
		roundStates := newTendermintState(log.New(), nil, nil)
		c := &Core{
			logger:      log.New("backend", "test", "id", 0),
			backend:     backendMock,
			roundsState: roundStates,
			committee:   committeeSet,
			lastHeader:  &types.Header{Committee: committeeSet.Committee()},
			address:     committeeSet.Committee()[0].Address,
		}
		c.SetHeight(common.Big3)
		c.roundsState.GetOrCreate(2)
		c.SetRound(2)

		c.SetDefaultHandlers()
		c.prevoter.SendPrevote(context.Background(), true)
	})

	t.Run("valid proposal given, non nil prevote", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet, keys := NewTestCommitteeSetWithKeys(4)
		member := committeeSet.Committee()[0]
		signer := makeSigner(keys[member.Address].consensus)
		logger := log.New("backend", "test", "id", 0)
		csize := len(committeeSet.Committee())

		proposal := message.NewPropose(
			1,
			2,
			1,
			types.NewBlockWithHeader(&types.Header{Number: big.NewInt(2)}),
			signer,
			&member)

		roundStates := newTendermintState(log.New(), nil, nil)
		backendMock := interfaces.NewMockBackend(ctrl)

		c := &Core{
			backend:     backendMock,
			address:     member.Address,
			logger:      logger,
			roundsState: roundStates,
			committee:   committeeSet,
			lastHeader:  &types.Header{Committee: committeeSet.Committee()},
		}
		c.SetHeight(common.Big2)
		c.SetRound(1)
		c.UpdateStep(Prevote)
		curMessages := c.roundsState.GetOrCreate(1)
		curMessages.SetProposal(proposal, true)
		expectedMsg := message.NewPrevote(1, 2, curMessages.ProposalHash(), signer, &member, csize)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(signer)
		backendMock.EXPECT().Broadcast(gomock.Any(), expectedMsg)

		c.SetDefaultHandlers()
		c.prevoter.SendPrevote(context.Background(), false)
	})
}

func TestHandlePrevote(t *testing.T) {
	committeeSet, keys := NewTestCommitteeSetWithKeys(4)
	member := committeeSet.Committee()[0]
	signer := makeSigner(keys[member.Address].consensus)
	csize := len(committeeSet.Committee())

	t.Run("pre-vote given with no errors, pre-vote added", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		proposal := message.NewPropose(
			1,
			2,
			1,
			types.NewBlockWithHeader(&types.Header{}),
			signer,
			&member)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any())
		roundStates := newTendermintState(log.New(), nil, nil)
		c := &Core{
			address:        member.Address,
			roundsState:    roundStates,
			logger:         log.Root(),
			committee:      committeeSet,
			prevoteTimeout: NewTimeout(Prevote, log.Root()),
			backend:        backendMock,
		}
		c.SetHeight(common.Big2)
		c.SetRound(1)
		c.UpdateStep(Prevote)
		curRoundMessages := c.roundsState.GetOrCreate(2)
		curRoundMessages.SetProposal(proposal, true)
		prevote := message.NewPrevote(1, 2, curRoundMessages.ProposalHash(), signer, &member, csize)

		c.SetDefaultHandlers()
		err := c.prevoter.HandlePrevote(context.Background(), prevote)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		if s := c.roundsState.GetOrCreate(1).PrevotesPower(curRoundMessages.ProposalHash()); s.Cmp(common.Big1) != 0 {
			t.Fatalf("Expected 1 prevote, but got %d", s)
		}
	})

	t.Run("pre-vote given at pre-vote step, non-nil pre-commit sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		proposal := message.NewPropose(
			2,
			3,
			1,
			types.NewBlockWithHeader(&types.Header{Number: big.NewInt(3)}),
			signer,
			&member)

		backendMock := interfaces.NewMockBackend(ctrl)
		roundStates := newTendermintState(log.New(), nil, nil)
		c := &Core{
			address:          member.Address,
			backend:          backendMock,
			roundsState:      roundStates,
			logger:           log.Root(),
			proposeTimeout:   NewTimeout(Propose, log.Root()),
			prevoteTimeout:   NewTimeout(Prevote, log.Root()),
			precommitTimeout: NewTimeout(Precommit, log.Root()),
			committee:        committeeSet,
			lastHeader:       &types.Header{Committee: committeeSet.Committee()},
		}
		c.SetHeight(common.Big3)
		c.SetRound(2)
		c.UpdateStep(Prevote)
		curRoundMessage := c.roundsState.GetOrCreate(2)
		curRoundMessage.SetProposal(proposal, true)

		// quorum of prevotes for v
		var prevotes []*message.Prevote
		for i := 0; i < 3; i++ {
			val := committeeSet.Committee()[i]
			prevote := message.NewPrevote(2, 3, curRoundMessage.ProposalHash(), makeSigner(keys[val.Address].consensus), &val, csize)
			prevotes = append(prevotes, prevote)
		}
		precommit := message.NewPrecommit(2, 3, curRoundMessage.ProposalHash(), signer, &member, csize)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(signer).AnyTimes()
		backendMock.EXPECT().Broadcast(gomock.Any(), precommit)
		backendMock.EXPECT().Post(gomock.Any()).MaxTimes(3)

		c.SetDefaultHandlers()
		for _, prevote := range prevotes {
			if err := c.prevoter.HandlePrevote(context.Background(), prevote); err != nil {
				t.Fatalf("Expected nil, got %v", err)
			}
		}
		if s := curRoundMessage.PrevotesPower(curRoundMessage.ProposalHash()); s.Cmp(common.Big3) != 0 {
			t.Fatalf("Expected 1 prevote, but got %d", s)
		}

		if !reflect.DeepEqual(c.ValidValue(), curRoundMessage.Proposal().Block()) {
			t.Fatalf("Expected %v, got %v", curRoundMessage.Proposal().Block(), c.ValidValue())
		}
	})

	t.Run("pre-vote given at pre-vote step, nil pre-commit sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		member2 := committeeSet.Committee()[1]

		// quorum of prevotes for nil
		var prevotes []*message.Prevote
		for i := 0; i < 3; i++ {
			val := committeeSet.Committee()[i]
			prevote := message.NewPrevote(2, 3, common.Hash{}, makeSigner(keys[val.Address].consensus), &val, csize)
			prevotes = append(prevotes, prevote)
		}
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(makeSigner(keys[member2.Address].consensus)).AnyTimes()
		backendMock.EXPECT().Post(gomock.Any()).MaxTimes(3)

		precommit := message.NewPrecommit(2, 3, common.Hash{}, makeSigner(keys[member2.Address].consensus), &member2, csize)

		backendMock.EXPECT().Broadcast(gomock.Any(), precommit)

		logger := log.New("backend", "test", "id", 0)
		roundStates := newTendermintState(log.New(), nil, nil)

		c := &Core{
			address:          member2.Address,
			backend:          backendMock,
			roundsState:      roundStates,
			logger:           logger,
			committee:        committeeSet,
			proposeTimeout:   NewTimeout(Propose, logger),
			prevoteTimeout:   NewTimeout(Prevote, logger),
			precommitTimeout: NewTimeout(Precommit, logger),
			lastHeader:       &types.Header{Committee: committeeSet.Committee()},
		}
		c.SetHeight(common.Big3)
		c.SetRound(2)
		c.UpdateStep(Prevote)
		c.roundsState.GetOrCreate(2)
		c.SetDefaultHandlers()
		for _, prevote := range prevotes {
			if err := c.prevoter.HandlePrevote(context.Background(), prevote); err != nil {
				t.Fatalf("Expected nil, got %v", err)
			}
		}
	})

	t.Run("quorum pre-vote given at pre-vote step, pre-vote Timeout triggered", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		member2 := committeeSet.Committee()[1]
		proposal := message.NewPropose(
			2,
			2,
			1,
			types.NewBlockWithHeader(&types.Header{}),
			signer,
			&member)

		backendMock := interfaces.NewMockBackend(ctrl)
		roundState := newTendermintState(log.New(), nil, nil)
		c := &Core{
			address:          member2.Address,
			backend:          backendMock,
			roundsState:      roundState,
			logger:           log.Root(),
			committee:        committeeSet,
			proposeTimeout:   NewTimeout(Propose, log.Root()),
			prevoteTimeout:   NewTimeout(Prevote, log.Root()),
			precommitTimeout: NewTimeout(Precommit, log.Root()),
			lastHeader:       &types.Header{Committee: committeeSet.Committee()},
		}
		c.SetHeight(common.Big2)
		c.SetRound(2)
		c.UpdateStep(Prevote)
		curRoundMessages := c.roundsState.GetOrCreate(2)
		curRoundMessages.SetProposal(proposal, true)

		// 2 prevotes for nil, 1 for v
		var prevotes []*message.Prevote
		val := committeeSet.Committee()[0]
		prevotes = append(prevotes, message.NewPrevote(2, 2, curRoundMessages.ProposalHash(), makeSigner(keys[val.Address].consensus), &val, csize))
		val2 := committeeSet.Committee()[1]
		prevotes = append(prevotes, message.NewPrevote(2, 2, common.Hash{}, makeSigner(keys[val2.Address].consensus), &val2, csize))
		val3 := committeeSet.Committee()[2]
		prevotes = append(prevotes, message.NewPrevote(2, 2, common.Hash{}, makeSigner(keys[val3.Address].consensus), &val3, csize))

		backendMock.EXPECT().Address().AnyTimes().Return(member2.Address)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())
		backendMock.EXPECT().Post(gomock.Any()).MaxTimes(4)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(makeSigner(keys[member2.Address].consensus)).AnyTimes()

		c.SetDefaultHandlers()

		for _, prevote := range prevotes {
			err := c.prevoter.HandlePrevote(context.Background(), prevote)
			if err != nil {
				t.Fatalf("Expected nil, got %v", err)
			}
		}

		require.True(t, c.prevoteTimeout.TimerStarted())
	})
}
