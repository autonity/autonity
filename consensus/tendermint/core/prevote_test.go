package core

import (
	"context"
	"math/big"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/log"
)

func TestSendPrevote(t *testing.T) {
	t.Run("proposal is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		c := &core{
			logger:            log.New("backend", "test", "id", 0),
			backend:           backendMock,
			currentRoundState: NewRoundState(big.NewInt(2), big.NewInt(3)),
		}

		c.sendPrevote(context.Background(), false)
	})

	t.Run("valid proposal given, non nil prevote", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		proposal := NewProposal(
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(1),
			types.NewBlockWithHeader(&types.Header{}))

		curRoundState := NewRoundState(big.NewInt(2), big.NewInt(3))
		curRoundState.SetProposal(proposal, nil)

		addr := common.HexToAddress("0x0123456789")

		var preVote = Vote{
			Round:             big.NewInt(curRoundState.Round().Int64()),
			Height:            big.NewInt(curRoundState.Height().Int64()),
			ProposedBlockHash: curRoundState.GetCurrentProposalHash(),
		}

		encodedVote, err := Encode(&preVote)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, nil)

		payload, err := expectedMsg.Payload()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)

		c := &core{
			backend:           backendMock,
			address:           addr,
			logger:            log.New("backend", "test", "id", 0),
			valSet:            new(validatorSet),
			currentRoundState: curRoundState,
		}

		c.sendPrevote(context.Background(), false)
	})
}

func TestHandlePrevote(t *testing.T) {
	t.Run("pre-vote with future height given, error returned", func(t *testing.T) {
		curRoundState := NewRoundState(big.NewInt(1), big.NewInt(2))
		addr := common.HexToAddress("0x0123456789")

		var preVote = Vote{
			Round:  big.NewInt(2),
			Height: big.NewInt(3),
		}

		encodedVote, err := Encode(&preVote)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		c := &core{
			address:           addr,
			currentRoundState: curRoundState,
			logger:            log.New("backend", "test", "id", 0),
			valSet:            new(validatorSet),
		}

		err = c.handlePrevote(context.Background(), expectedMsg)
		if err != errFutureHeightMessage {
			t.Fatalf("Expected %v, got %v", errFutureHeightMessage, err)
		}
	})

	t.Run("pre-vote with old height given, pre-vote added", func(t *testing.T) {
		curRoundState := NewRoundState(big.NewInt(2), big.NewInt(3))
		addr := common.HexToAddress("0x0123456789")

		var preVote = Vote{
			Round:             big.NewInt(1),
			Height:            big.NewInt(3),
			ProposedBlockHash: curRoundState.GetCurrentProposalHash(),
		}

		encodedVote, err := Encode(&preVote)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		c := &core{
			address:                      addr,
			currentRoundState:            curRoundState,
			currentHeightOldRoundsStates: make(map[int64]*roundState),
			logger:                       log.New("backend", "test", "id", 0),
			valSet:                       new(validatorSet),
		}

		err = c.handlePrevote(context.Background(), expectedMsg)
		if err != errOldRoundMessage {
			t.Fatalf("Expected %v, got %v", errOldRoundMessage, err)
		}

		c.currentHeightOldRoundsStatesMu.Lock()
		defer c.currentHeightOldRoundsStatesMu.Unlock()
		oldRoundState := c.currentHeightOldRoundsStates[preVote.Round.Int64()]
		if s := oldRoundState.Prevotes.NilVotesSize(); s != 1 {
			t.Fatalf("Expected 1 nil-prevote, but got %d", s)
		}
	})

	t.Run("pre-vote given with no errors, pre-vote added", func(t *testing.T) {
		proposal := NewProposal(
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(1),
			types.NewBlockWithHeader(&types.Header{}))

		curRoundState := NewRoundState(big.NewInt(2), big.NewInt(3))
		curRoundState.SetProposal(proposal, nil)

		addr := common.HexToAddress("0x0123456789")

		var preVote = Vote{
			Round:             big.NewInt(curRoundState.Round().Int64()),
			Height:            big.NewInt(curRoundState.Height().Int64()),
			ProposedBlockHash: curRoundState.GetCurrentProposalHash(),
		}

		encodedVote, err := Encode(&preVote)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		c := &core{
			address:           addr,
			currentRoundState: curRoundState,
			logger:            log.New("backend", "test", "id", 0),
			valSet:            new(validatorSet),
		}

		err = c.handlePrevote(context.Background(), expectedMsg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		if s := c.currentRoundState.Prevotes.VotesSize(curRoundState.GetCurrentProposalHash()); s != 1 {
			t.Fatalf("Expected 1 prevote, but got %d", s)
		}
	})

	t.Run("pre-vote given at pre-vote step, non-nil pre-commit sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		proposal := NewProposal(
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(1),
			types.NewBlockWithHeader(&types.Header{}))

		curRoundState := NewRoundState(big.NewInt(2), big.NewInt(3))
		curRoundState.SetProposal(proposal, nil)
		curRoundState.SetStep(prevote)

		addr := common.HexToAddress("0x0123456789")

		var preVote = Vote{
			Round:             big.NewInt(curRoundState.Round().Int64()),
			Height:            big.NewInt(curRoundState.Height().Int64()),
			ProposedBlockHash: curRoundState.GetCurrentProposalHash(),
		}

		encodedVote, err := Encode(&preVote)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, nil).AnyTimes()

		var precommit = Vote{
			Round:             big.NewInt(curRoundState.Round().Int64()),
			Height:            big.NewInt(curRoundState.Height().Int64()),
			ProposedBlockHash: curRoundState.GetCurrentProposalHash(),
		}

		encodedVote, err = Encode(&precommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		msg := &Message{
			Code:          msgPrecommit,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{0x1},
			Signature:     []byte{0x1},
		}

		payload, err := msg.Payload()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		backendMock.EXPECT().Broadcast(context.Background(), gomock.Any(), payload)

		c := &core{
			address:           addr,
			backend:           backendMock,
			currentRoundState: curRoundState,
			logger:            log.New("backend", "test", "id", 0),
			prevoteTimeout:    newTimeout(prevote),
			valSet:            new(validatorSet),
		}

		err = c.handlePrevote(context.Background(), expectedMsg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		if s := c.currentRoundState.Prevotes.VotesSize(curRoundState.GetCurrentProposalHash()); s != 1 {
			t.Fatalf("Expected 1 prevote, but got %d", s)
		}

		if !reflect.DeepEqual(c.validValue, c.currentRoundState.Proposal().ProposalBlock) {
			t.Fatalf("Expected %v, got %v", c.currentRoundState.Proposal().ProposalBlock, c.validValue)
		}
	})

	t.Run("pre-vote given at pre-vote step, nil pre-commit sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		curRoundState := NewRoundState(big.NewInt(2), big.NewInt(3))
		curRoundState.SetStep(prevote)

		addr := common.HexToAddress("0x0123456789")

		var preVote = Vote{
			Round:             big.NewInt(curRoundState.Round().Int64()),
			Height:            big.NewInt(curRoundState.Height().Int64()),
			ProposedBlockHash: curRoundState.GetCurrentProposalHash(),
		}

		encodedVote, err := Encode(&preVote)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, nil).AnyTimes()

		var precommit = Vote{
			Round:             big.NewInt(curRoundState.Round().Int64()),
			Height:            big.NewInt(curRoundState.Height().Int64()),
			ProposedBlockHash: common.Hash{},
		}

		encodedVote, err = Encode(&precommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		msg := &Message{
			Code:          msgPrecommit,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{0x1},
			Signature:     []byte{0x1},
		}

		payload, err := msg.Payload()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		backendMock.EXPECT().Broadcast(context.Background(), gomock.Any(), payload)

		c := &core{
			address:           addr,
			backend:           backendMock,
			currentRoundState: curRoundState,
			logger:            log.New("backend", "test", "id", 0),
			prevoteTimeout:    newTimeout(prevote),
			valSet:            new(validatorSet),
		}

		err = c.handlePrevote(context.Background(), expectedMsg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})

	t.Run("pre-vote given at pre-vote step, pre-vote timeout triggered", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		proposal := NewProposal(
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(1),
			types.NewBlockWithHeader(&types.Header{}))

		addr := common.HexToAddress("0x0123456789")

		curRoundState := NewRoundState(big.NewInt(2), big.NewInt(3))
		curRoundState.SetProposal(proposal, nil)
		curRoundState.SetStep(prevote)
		curRoundState.Prevotes.AddVote(addr.Hash(), Message{})

		var preVote = Vote{
			Round:             big.NewInt(curRoundState.Round().Int64()),
			Height:            big.NewInt(curRoundState.Height().Int64()),
			ProposedBlockHash: curRoundState.GetCurrentProposalHash(),
		}

		encodedVote, err := Encode(&preVote)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(addr)

		c := New(backendMock, nil)
		c.currentRoundState = curRoundState
		c.prevoteTimeout = newTimeout(prevote)
		c.valSet = &validatorSet{
			Set: newTestValidatorSet(2),
		}

		err = c.handlePrevote(context.Background(), expectedMsg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})
}
