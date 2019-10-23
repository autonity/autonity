package core

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
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
			logger:            log.New("backend", "test", "id", 0),
			backend:           backendMock,
			currentRoundState: NewRoundState(big.NewInt(2), big.NewInt(3)),
		}

		c.sendPrecommit(context.Background(), false)
	})

	t.Run("valid proposal given, non nil pre-commit", func(t *testing.T) {
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

		var preCommit = Vote{
			Round:             big.NewInt(curRoundState.Round().Int64()),
			Height:            big.NewInt(curRoundState.Height().Int64()),
			ProposedBlockHash: curRoundState.GetCurrentProposalHash(),
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
			backend:           backendMock,
			address:           addr,
			logger:            log.New("backend", "test", "id", 0),
			valSet:            new(validatorSet),
			currentRoundState: curRoundState,
		}

		c.sendPrecommit(context.Background(), false)
	})

	t.Run("valid proposal given, nil pre-commit", func(t *testing.T) {
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
			backend:           backendMock,
			address:           addr,
			logger:            log.New("backend", "test", "id", 0),
			valSet:            new(validatorSet),
			currentRoundState: curRoundState,
		}

		c.sendPrecommit(context.Background(), true)
	})
}

func TestHandlePrecommit(t *testing.T) {
	t.Run("pre-commit with future height given, error returned", func(t *testing.T) {
		curRoundState := NewRoundState(big.NewInt(1), big.NewInt(2))
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
			address:           addr,
			currentRoundState: curRoundState,
			logger:            log.New("backend", "test", "id", 0),
			valSet:            new(validatorSet),
		}

		err = c.handlePrecommit(context.Background(), expectedMsg)
		if err != errFutureHeightMessage {
			t.Fatalf("Expected %v, got %v", errFutureHeightMessage, err)
		}
	})

	t.Run("pre-commit with invalid signature given, error returned", func(t *testing.T) {
		curRoundState := NewRoundState(big.NewInt(2), big.NewInt(3))
		addr := common.HexToAddress("0x0123456789")

		var preCommit = Vote{
			Round:             big.NewInt(curRoundState.Round().Int64()),
			Height:            big.NewInt(curRoundState.Height().Int64()),
			ProposedBlockHash: curRoundState.GetCurrentProposalHash(),
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
			address:           addr,
			currentRoundState: curRoundState,
			logger:            log.New("backend", "test", "id", 0),
			valSet:            new(validatorSet),
		}

		err = c.handlePrecommit(context.Background(), expectedMsg)
		if err != secp256k1.ErrInvalidSignatureLen {
			t.Fatalf("Expected %v, got %v", secp256k1.ErrInvalidSignatureLen, err)
		}
	})

	t.Run("pre-commit given with no errors, commit called", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		proposal := NewProposal(
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(1),
			types.NewBlockWithHeader(&types.Header{}))
		proposal.ProposalBlock.Hash()

		curRoundState := NewRoundState(big.NewInt(2), big.NewInt(3))
		curRoundState.SetProposal(proposal, nil)

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

		data := PrepareCommittedSeal(addr.Hash())
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

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Commit(*proposal.ProposalBlock, gomock.Any()).Return(nil)

		c := &core{
			address:           addr,
			backend:           backendMock,
			currentRoundState: curRoundState,
			logger:            log.New("backend", "test", "id", 0),
			valSet:            new(validatorSet),
			precommitTimeout:  newTimeout(precommit),
		}

		err = c.handlePrecommit(context.Background(), expectedMsg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})

	t.Run("pre-commit given with no errors, commit cancelled", func(t *testing.T) {
		proposal := NewProposal(
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(1),
			types.NewBlockWithHeader(&types.Header{}))

		curRoundState := NewRoundState(big.NewInt(2), big.NewInt(3))
		curRoundState.SetProposal(proposal, nil)

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

		data := PrepareCommittedSeal(addr.Hash())
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
			address:           addr,
			currentRoundState: curRoundState,
			logger:            log.New("backend", "test", "id", 0),
			valSet:            new(validatorSet),
			precommitTimeout:  newTimeout(precommit),
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err = c.handlePrecommit(ctx, expectedMsg)
		if err != context.Canceled {
			t.Fatalf("Expected %v, got %v", context.Canceled, err)
		}
	})

	t.Run("pre-commit given with no errors, pre-commit timeout triggered", func(t *testing.T) {
		curRoundState := NewRoundState(big.NewInt(2), big.NewInt(3))
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

		data := PrepareCommittedSeal(addr.Hash())
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
			address:           addr,
			currentRoundState: curRoundState,
			logger:            log.New("backend", "test", "id", 0),
			valSet:            new(validatorSet),
			precommitTimeout:  newTimeout(precommit),
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

		err := c.verifyPrecommitCommittedSeal(addrMsg, nil, common.Hash{})
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

		data := PrepareCommittedSeal(common.Hash{})
		hashData := crypto.Keccak256(data)
		sig, err := crypto.Sign(hashData, key)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		err = c.verifyPrecommitCommittedSeal(addrMsg, sig, common.Hash{})
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

		data := PrepareCommittedSeal(addrMsg.Hash())
		hashData := crypto.Keccak256(data)
		sig, err := crypto.Sign(hashData, key)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		err = c.verifyPrecommitCommittedSeal(addrMsg, sig, addrMsg.Hash())
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})
}

func TestHandleCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	block := types.NewBlockWithHeader(&types.Header{})
	addr := common.HexToAddress("0x0123456789")

	backendMock := NewMockBackend(ctrl)
	backendMock.EXPECT().LastCommittedProposal().Return(block, addr)

	valSet := validator.NewMockSet(ctrl)
	valSet.EXPECT().CalcProposer(addr, uint64(0))
	valSet.EXPECT().IsProposer(addr).Return(false)

	backendMock.EXPECT().Validators(uint64(1)).Return(valSet)

	c := &core{
		address:           addr,
		backend:           backendMock,
		currentRoundState: NewRoundState(big.NewInt(2), big.NewInt(3)),
		logger:            log.New("backend", "test", "id", 0),
		proposeTimeout:    newTimeout(propose),
		prevoteTimeout:    newTimeout(prevote),
		precommitTimeout:  newTimeout(precommit),
		valSet:            new(validatorSet),
	}
	c.handleCommit(context.Background())
}
