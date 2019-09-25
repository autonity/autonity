package core

import (
	"context"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/log"
)

func TestSendPropose(t *testing.T) {
	t.Run("valid block given, proposal is broadcast", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		curRoundState := NewRoundState(big.NewInt(1), big.NewInt(1))
		validRound := big.NewInt(1)

		proposalBlock := NewProposal(curRoundState.round, curRoundState.Height(), validRound, block)
		proposal, err := Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		payloadNoSig, err := expectedMsg.PayloadNoSig()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		payload, err := expectedMsg.Payload()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		valSetMock := validator.NewMockSet(ctrl)
		valSetMock.EXPECT().IsProposer(addr).Return(true).AnyTimes()
		valSetMock.EXPECT().GetProposer()
		valSetMock.EXPECT().Copy()

		valSet := &validatorSet{
			Set: valSetMock,
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().SetProposedBlockHash(block.Hash())
		backendMock.EXPECT().Sign(payloadNoSig).Return([]byte{0x1}, nil)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)

		c := &core{
			address:           addr,
			backend:           backendMock,
			currentRoundState: curRoundState,
			logger:            log.New("backend", "test", "id", 0),
			validRound:        validRound,
			valSet:            valSet,
		}

		c.sendProposal(context.Background(), block)
	})
}

func TestHandleProposal(t *testing.T) {
	t.Run("old proposal given, error returned", func(t *testing.T) {
		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		curRoundState := NewRoundState(big.NewInt(2), big.NewInt(1))
		validRound := big.NewInt(1)

		proposalBlock := NewProposal(big.NewInt(1), curRoundState.Height(), validRound, block)
		proposal, err := Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &Message{
			Code:          msgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		c := &core{
			address:           addr,
			currentRoundState: curRoundState,
			logger:            log.New("backend", "test", "id", 0),
			validRound:        validRound,
		}

		err = c.handleProposal(context.Background(), msg)
		if err != errOldRoundMessage {
			t.Fatalf("Expected %v, got %v", errOldRoundMessage, err)
		}
	})

	t.Run("msg from non-proposer given, error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		curRoundState := NewRoundState(big.NewInt(2), big.NewInt(1))
		validRound := big.NewInt(1)

		proposalBlock := NewProposal(curRoundState.Round(), curRoundState.Height(), validRound, block)
		proposal, err := Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &Message{
			Code:          msgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		valSetMock := validator.NewMockSet(ctrl)
		valSetMock.EXPECT().IsProposer(addr).Return(false)

		valSet := &validatorSet{
			Set: valSetMock,
		}

		c := &core{
			address:           addr,
			currentRoundState: curRoundState,
			logger:            log.New("backend", "test", "id", 0),
			validRound:        validRound,
			valSet:            valSet,
		}

		err = c.handleProposal(context.Background(), msg)
		if err != errNotFromProposer {
			t.Fatalf("Expected %v, got %v", errNotFromProposer, err)
		}
	})

	t.Run("unverified proposal given, error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		curRoundState := NewRoundState(big.NewInt(2), big.NewInt(1))
		validRound := big.NewInt(1)

		proposalBlock := NewProposal(curRoundState.Round(), curRoundState.Height(), validRound, block)
		proposal, err := Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &Message{
			Code:          msgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		sender := validator.NewMockValidator(ctrl)

		valSetMock := validator.NewMockSet(ctrl)
		valSetMock.EXPECT().IsProposer(addr).Return(true).AnyTimes()
		valSetMock.EXPECT().GetProposer()
		valSetMock.EXPECT().Size().AnyTimes()
		valSetMock.EXPECT().Copy()
		valSetMock.EXPECT().GetByAddress(msg.Address).Return(1, sender).AnyTimes()

		valSet := &validatorSet{
			Set: valSetMock,
		}

		var decProposal Proposal
		if err := msg.Decode(&decProposal); err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		var prevote = Vote{
			Round:             big.NewInt(curRoundState.Round().Int64()),
			Height:            big.NewInt(curRoundState.Height().Int64()),
			ProposedBlockHash: common.Hash{},
		}

		encodedVote, err := Encode(&prevote)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		preVoteMsg := &Message{
			Code:          msgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
		}

		payloadNoSig, err := preVoteMsg.PayloadNoSig()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		payload, err := preVoteMsg.Payload()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		event := backlogEvent{
			src: sender,
			msg: msg,
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).Return(time.Nanosecond, consensus.ErrFutureBlock)
		backendMock.EXPECT().Sign(payloadNoSig)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)
		backendMock.EXPECT().Post(event).AnyTimes()

		c := &core{
			address:           addr,
			backend:           backendMock,
			currentRoundState: curRoundState,
			logger:            log.New("backend", "test", "id", 0),
			proposeTimeout:    newTimeout(propose),
			validRound:        validRound,
			valSet:            valSet,
		}

		err = c.handleProposal(context.Background(), msg)
		if err != consensus.ErrFutureBlock {
			t.Fatalf("Expected %v, got %v", consensus.ErrFutureBlock, err)
		}
	})

	t.Run("valid proposal given, no error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		curRoundState := NewRoundState(big.NewInt(2), big.NewInt(1))
		validRound := big.NewInt(1)

		proposalBlock := NewProposal(curRoundState.Round(), curRoundState.Height(), validRound, block)
		proposal, err := Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &Message{
			Code:          msgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		valSetMock := validator.NewMockSet(ctrl)
		valSetMock.EXPECT().IsProposer(addr).Return(true).AnyTimes()
		valSetMock.EXPECT().GetProposer()

		valSet := &validatorSet{
			Set: valSetMock,
		}

		var decProposal Proposal
		if err := msg.Decode(&decProposal); err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(*decProposal.ProposalBlock)

		c := &core{
			address:           addr,
			backend:           backendMock,
			currentRoundState: curRoundState,
			logger:            log.New("backend", "test", "id", 0),
			proposeTimeout:    newTimeout(propose),
			validRound:        validRound,
			valSet:            valSet,
		}

		err = c.handleProposal(context.Background(), msg)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		if !reflect.DeepEqual(curRoundState.proposalMsg, msg) {
			t.Fatalf("%v not equal to  %v", curRoundState.proposalMsg, msg)
		}
	})

	t.Run("valid proposal given, valid round -1, pre-vote is sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		curRoundState := NewRoundState(big.NewInt(2), big.NewInt(1))
		validRound := big.NewInt(-1)

		proposalBlock := NewProposal(curRoundState.Round(), curRoundState.Height(), validRound, block)
		proposal, err := Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &Message{
			Code:          msgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		valSetMock := validator.NewMockSet(ctrl)
		valSetMock.EXPECT().IsProposer(addr).Return(true).AnyTimes()
		valSetMock.EXPECT().GetProposer().AnyTimes()
		valSetMock.EXPECT().Size().AnyTimes()
		valSetMock.EXPECT().Copy()

		valSet := &validatorSet{
			Set: valSetMock,
		}

		var decProposal Proposal
		if err := msg.Decode(&decProposal); err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		var prevote = Vote{
			Round:             big.NewInt(curRoundState.Round().Int64()),
			Height:            big.NewInt(curRoundState.Height().Int64()),
			ProposedBlockHash: common.Hash{},
		}

		encodedVote, err := Encode(&prevote)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		preVoteMsg := &Message{
			Code:          msgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
		}

		payloadNoSig, err := preVoteMsg.PayloadNoSig()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		payload, err := preVoteMsg.Payload()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(*decProposal.ProposalBlock)
		backendMock.EXPECT().Sign(payloadNoSig)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)

		c := &core{
			address:           addr,
			backend:           backendMock,
			currentRoundState: curRoundState,
			lockedValue:       types.NewBlockWithHeader(&types.Header{}),
			lockedRound:       big.NewInt(-1),
			logger:            log.New("backend", "test", "id", 0),
			proposeTimeout:    newTimeout(propose),
			validRound:        validRound,
			valSet:            valSet,
		}

		err = c.handleProposal(context.Background(), msg)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		if !reflect.DeepEqual(curRoundState.proposalMsg, msg) {
			t.Fatalf("%v not equal to  %v", curRoundState.proposalMsg, msg)
		}
	})

	t.Run("valid proposal given, vr < curR, pre-vote is sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		curRoundState := NewRoundState(big.NewInt(2), big.NewInt(1))
		validRound := big.NewInt(0)

		proposalBlock := NewProposal(curRoundState.Round(), curRoundState.Height(), validRound, block)
		proposal, err := Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &Message{
			Code:          msgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		valSetMock := validator.NewMockSet(ctrl)
		valSetMock.EXPECT().IsProposer(addr).Return(true).AnyTimes()
		valSetMock.EXPECT().GetProposer().AnyTimes()
		valSetMock.EXPECT().Size().AnyTimes()
		valSetMock.EXPECT().Copy()

		valSet := &validatorSet{
			Set: valSetMock,
		}

		var decProposal Proposal
		if err := msg.Decode(&decProposal); err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		var prevote = Vote{
			Round:             big.NewInt(curRoundState.Round().Int64()),
			Height:            big.NewInt(curRoundState.Height().Int64()),
			ProposedBlockHash: common.Hash{},
		}

		encodedVote, err := Encode(&prevote)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		preVoteMsg := &Message{
			Code:          msgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
		}

		payloadNoSig, err := preVoteMsg.PayloadNoSig()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		payload, err := preVoteMsg.Payload()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(*decProposal.ProposalBlock)
		backendMock.EXPECT().Sign(payloadNoSig)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)

		c := &core{
			address:           addr,
			backend:           backendMock,
			currentRoundState: curRoundState,
			currentHeightOldRoundsStates: map[int64]roundState{
				0: *curRoundState,
			},
			lockedRound:    big.NewInt(-1),
			lockedValue:    types.NewBlockWithHeader(&types.Header{}),
			logger:         log.New("backend", "test", "id", 0),
			proposeTimeout: newTimeout(propose),
			validRound:     validRound,
			valSet:         valSet,
		}

		err = c.handleProposal(context.Background(), msg)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		if !reflect.DeepEqual(curRoundState.proposalMsg, msg) {
			t.Fatalf("%v not equal to  %v", curRoundState.proposalMsg, msg)
		}
	})
}
