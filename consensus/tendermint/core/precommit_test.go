package core

import (
	"context"
	"errors"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"math/big"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/crypto/secp256k1"
	"github.com/clearmatics/autonity/log"
)

func TestSendPrecommit(t *testing.T) {
	t.Run("proposal is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		c := &core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			curRoundMessages: NewRoundMessages(big.NewInt(2), big.NewInt(3)),
		}

		c.sendPrecommit(context.Background(), false)
	})

	t.Run("valid proposal given, non nil pre-commit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := log.New("backend", "test", "id", 0)

		proposal := NewProposal(
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(1),
			types.NewBlockWithHeader(&types.Header{}))

		curRoundState := NewRoundMessages(big.NewInt(2), big.NewInt(3))
		curRoundState.SetProposal(proposal, nil)

		addr := common.HexToAddress("0x0123456789")

		var preCommit = Vote{
			Round:             big.NewInt(curRoundState.Round().Int64()),
			Height:            big.NewInt(curRoundState.Height().Int64()),
			ProposedBlockHash: curRoundState.GetProposalHash(),
		}

		encodedVote, err := Encode(&preCommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrecommit,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{0x1},
			Signature:     []byte{0x1},
		}

		payloadNoSig, err := expectedMsg.PayloadNoSig()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, nil)
		backendMock.EXPECT().Sign(payloadNoSig).Return([]byte{0x1}, nil)

		payload, err := expectedMsg.Payload()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)

		c := &core{
			backend:          backendMock,
			address:          addr,
			logger:           logger,
			committeeSet:     new(validatorSet),
			curRoundMessages: curRoundState,
		}

		c.sendPrecommit(context.Background(), false)
	})

	t.Run("valid proposal given, nil pre-commit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := log.New("backend", "test", "id", 0)

		proposal := NewProposal(
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(1),
			types.NewBlockWithHeader(&types.Header{}))

		curRoundState := NewRoundMessages(big.NewInt(2), big.NewInt(3))
		curRoundState.SetProposal(proposal, nil)

		addr := common.HexToAddress("0x0123456789")

		var preCommit = Vote{
			Round:             big.NewInt(curRoundState.Round().Int64()),
			Height:            big.NewInt(curRoundState.Height().Int64()),
			ProposedBlockHash: common.Hash{},
		}

		encodedVote, err := Encode(&preCommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrecommit,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{0x1},
			Signature:     []byte{0x1},
		}

		payloadNoSig, err := expectedMsg.PayloadNoSig()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, errors.New("seal sign error"))
		backendMock.EXPECT().Sign(payloadNoSig).Return([]byte{0x1}, nil)

		payload, err := expectedMsg.Payload()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)

		c := &core{
			backend:          backendMock,
			address:          addr,
			logger:           logger,
			committeeSet:     new(validatorSet),
			curRoundMessages: curRoundState,
		}

		c.sendPrecommit(context.Background(), true)
	})
}

func TestHandlePrecommit(t *testing.T) {
	t.Run("pre-commit with future height given, error returned", func(t *testing.T) {
		curRoundState := NewRoundMessages(big.NewInt(1), big.NewInt(2))
		addr := common.HexToAddress("0x0123456789")

		var preCommit = Vote{
			Round:  big.NewInt(2),
			Height: big.NewInt(3),
		}

		encodedVote, err := Encode(&preCommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrecommit,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		c := &core{
			address:          addr,
			curRoundMessages: curRoundState,
			logger:           log.New("backend", "test", "id", 0),
			committeeSet:     new(validatorSet),
		}

		err = c.handlePrecommit(context.Background(), expectedMsg)
		if err != errFutureHeightMessage {
			t.Fatalf("Expected %v, got %v", errFutureHeightMessage, err)
		}
	})

	t.Run("pre-commit with invalid signature given, error returned", func(t *testing.T) {
		curRoundState := NewRoundMessages(big.NewInt(2), big.NewInt(3))
		addr := common.HexToAddress("0x0123456789")

		var preCommit = Vote{
			Round:             big.NewInt(curRoundState.Round().Int64()),
			Height:            big.NewInt(curRoundState.Height().Int64()),
			ProposedBlockHash: curRoundState.GetProposalHash(),
		}

		encodedVote, err := Encode(&preCommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrecommit,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		c := &core{
			address:          addr,
			curRoundMessages: curRoundState,
			logger:           log.New("backend", "test", "id", 0),
			committeeSet:     new(validatorSet),
		}
		c.setStep(precommit)
		err = c.handlePrecommit(context.Background(), expectedMsg)
		if err != secp256k1.ErrInvalidSignatureLen {
			t.Fatalf("Expected %v, got %v", secp256k1.ErrInvalidSignatureLen, err)
		}
	})

	t.Run("pre-commit given with no errors, commit called", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := log.New("backend", "test", "id", 0)

		proposal := NewProposal(
			big.NewInt(2),
			big.NewInt(3),
			big.NewInt(1),
			types.NewBlockWithHeader(&types.Header{}))
		proposal.ProposalBlock.Hash()

		curRoundState := NewRoundMessages(big.NewInt(2), big.NewInt(3))
		curRoundState.SetProposal(proposal, nil)
		curRoundState.SetStep(precommit)

		var preCommit = Vote{
			Round:             big.NewInt(curRoundState.Round().Int64()),
			Height:            big.NewInt(curRoundState.Height().Int64()),
			ProposedBlockHash: proposal.ProposalBlock.Hash(),
		}
		// Notice that the current roundstate's round and height can differs from the received precommit seals.
		// This is impossible to happen in the current implementation but actually the specification allows it.
		// see https://github.com/clearmatics/autonity/issues/347
		encodedVote, err := Encode(&preCommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		key, err := generatePrivateKey()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		data := PrepareCommittedSeal(proposal.ProposalBlock.Hash(), preCommit.Round, preCommit.Height)
		hashData := crypto.Keccak256(data)
		sig, err := crypto.Sign(hashData, key)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrecommit,
			Msg:           encodedVote,
			Address:       getAddress(),
			CommittedSeal: sig,
			Signature:     []byte{0x1},
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Commit(proposal.ProposalBlock, gomock.Any(), gomock.Any()).Return(nil).Do(
			func(proposalBlock *types.Block, round *big.Int, seals [][]byte) {
				if round.Cmp(preCommit.Round) != 0 {
					t.Fatal("Commit called with round different than precommit seal")
				}
				if !reflect.DeepEqual([][]byte{expectedMsg.CommittedSeal}, seals) {
					t.Fatal("Commit called with wrong seal")
				}
			})

		c := &core{
			address:          getAddress(),
			backend:          backendMock,
			curRoundMessages: curRoundState,
			logger:           logger,
			committeeSet:     new(validatorSet),
			precommitTimeout: newTimeout(precommit, logger),
		}

		err = c.handlePrecommit(context.Background(), expectedMsg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})

	t.Run("pre-commit given with no errors, commit cancelled", func(t *testing.T) {
		logger := log.New("backend", "test", "id", 0)

		proposal := NewProposal(
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(1),
			types.NewBlockWithHeader(&types.Header{}))

		curRoundState := NewRoundMessages(big.NewInt(2), big.NewInt(3))
		curRoundState.SetProposal(proposal, nil)
		curRoundState.SetStep(prevote)

		addr := getAddress()

		var preCommit = Vote{
			Round:             big.NewInt(curRoundState.Round().Int64()),
			Height:            big.NewInt(curRoundState.Height().Int64()),
			ProposedBlockHash: addr.Hash(),
		}

		encodedVote, err := Encode(&preCommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		key, err := generatePrivateKey()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		data := PrepareCommittedSeal(addr.Hash(), preCommit.Round, preCommit.Height)
		hashData := crypto.Keccak256(data)
		sig, err := crypto.Sign(hashData, key)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrecommit,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: sig,
			Signature:     []byte{0x1},
		}

		c := &core{
			address:          addr,
			curRoundMessages: curRoundState,
			logger:           logger,
			committeeSet:     new(validatorSet),
			precommitTimeout: newTimeout(precommit, logger),
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err = c.handlePrecommit(ctx, expectedMsg)
		if err != context.Canceled {
			t.Fatalf("Expected %v, got %v", context.Canceled, err)
		}
	})

	t.Run("pre-commit given with no errors, pre-commit timeout triggered", func(t *testing.T) {
		curRoundState := NewRoundMessages(big.NewInt(2), big.NewInt(3))
		curRoundState.SetStep(precommit)
		addr := getAddress()

		var preCommit = Vote{
			Round:             big.NewInt(curRoundState.Round().Int64()),
			Height:            big.NewInt(curRoundState.Height().Int64()),
			ProposedBlockHash: addr.Hash(),
		}

		encodedVote, err := Encode(&preCommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		key, err := generatePrivateKey()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		data := PrepareCommittedSeal(addr.Hash(), preCommit.Round, preCommit.Height)
		hashData := crypto.Keccak256(data)
		sig, err := crypto.Sign(hashData, key)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrecommit,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: sig,
			Signature:     []byte{0x1},
		}

		logger := log.New("backend", "test", "id", 0)

		c := &core{
			address:          addr,
			curRoundMessages: curRoundState,
			logger:           logger,
			committeeSet:     new(validatorSet),
			precommitTimeout: newTimeout(precommit, logger),
		}

		err = c.handlePrecommit(context.Background(), expectedMsg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})
}

func TestVerifyPrecommitCommittedSeal(t *testing.T) {
	t.Run("invalid signer address given, error returned", func(t *testing.T) {
		c := &core{
			logger: log.New("backend", "test", "id", 0),
		}

		addrMsg := common.HexToAddress("0x0123456789")

		err := c.verifyCommittedSeal(addrMsg, nil, common.Hash{}, new(big.Int), new(big.Int))
		if err != secp256k1.ErrInvalidSignatureLen {
			t.Fatalf("Expected %v, got %v", secp256k1.ErrInvalidSignatureLen, err)
		}
	})

	t.Run("invalid signer address given, error returned", func(t *testing.T) {
		c := &core{
			logger: log.New("backend", "test", "id", 0),
		}

		addrMsg := common.HexToAddress("0x0123456789")
		key, err := generatePrivateKey()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		data := PrepareCommittedSeal(common.Hash{}, new(big.Int), new(big.Int))
		hashData := crypto.Keccak256(data)
		sig, err := crypto.Sign(hashData, key)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		err = c.verifyCommittedSeal(addrMsg, sig, common.Hash{}, new(big.Int), new(big.Int))
		if err != errInvalidSenderOfCommittedSeal {
			t.Fatalf("Expected %v, got %v", errInvalidSenderOfCommittedSeal, err)
		}
	})

	t.Run("valid params given, no error returned", func(t *testing.T) {
		c := &core{
			logger: log.New("backend", "test", "id", 0),
		}

		addrMsg := getAddress()
		key, err := generatePrivateKey()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		data := PrepareCommittedSeal(addrMsg.Hash(), new(big.Int), new(big.Int))
		hashData := crypto.Keccak256(data)
		sig, err := crypto.Sign(hashData, key)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		err = c.verifyCommittedSeal(addrMsg, sig, addrMsg.Hash(), new(big.Int), new(big.Int))
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})
}

func TestHandleCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := log.New("backend", "test", "id", 0)

	block := types.NewBlockWithHeader(&types.Header{})
	addr := common.HexToAddress("0x0123456789")

	backendMock := NewMockBackend(ctrl)
	backendMock.EXPECT().LastCommittedProposal().MinTimes(1).Return(block, addr)

	valSet := committee.NewMockSet(ctrl)
	valSet.EXPECT().CalcProposer(addr, uint64(0))
	valSet.EXPECT().IsProposer(addr).Return(false)

	backendMock.EXPECT().Validators(uint64(1)).Return(valSet)

	c := &core{
		address:          addr,
		backend:          backendMock,
		curRoundMessages: NewRoundMessages(big.NewInt(2), big.NewInt(3)),
		logger:           logger,
		proposeTimeout:   newTimeout(propose, logger),
		prevoteTimeout:   newTimeout(prevote, logger),
		precommitTimeout: newTimeout(precommit, logger),
		committeeSet:     new(validatorSet),
	}
	c.handleCommit(context.Background())
}
