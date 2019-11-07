package core

import (
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/log"
)

func TestCheckMessage(t *testing.T) {
	t.Run("valid params given, nil returned", func(t *testing.T) {
		c := &core{
			currentRoundState: NewRoundState(big.NewInt(1), big.NewInt(2)),
		}

		err := c.checkMessage(big.NewInt(1), big.NewInt(2))
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}
	})

	t.Run("given nil round, error returned", func(t *testing.T) {
		c := &core{}

		err := c.checkMessage(nil, big.NewInt(2))
		if err != errInvalidMessage {
			t.Fatalf("have %v, want %v", err, errInvalidMessage)
		}
	})

	t.Run("given future height, error returned", func(t *testing.T) {
		c := &core{
			currentRoundState: NewRoundState(big.NewInt(2), big.NewInt(3)),
		}

		err := c.checkMessage(big.NewInt(2), big.NewInt(4))
		if err != errFutureHeightMessage {
			t.Fatalf("have %v, want %v", err, errFutureHeightMessage)
		}
	})

	t.Run("given old height, error returned", func(t *testing.T) {
		c := &core{
			currentRoundState: NewRoundState(big.NewInt(2), big.NewInt(3)),
		}

		err := c.checkMessage(big.NewInt(2), big.NewInt(2))
		if err != errOldHeightMessage {
			t.Fatalf("have %v, want %v", err, errOldHeightMessage)
		}
	})

	t.Run("given future round, error returned", func(t *testing.T) {
		c := &core{
			currentRoundState: NewRoundState(big.NewInt(2), big.NewInt(3)),
		}

		err := c.checkMessage(big.NewInt(3), big.NewInt(3))
		if err != errFutureRoundMessage {
			t.Fatalf("have %v, want %v", err, errFutureRoundMessage)
		}
	})

	t.Run("given old round, error returned", func(t *testing.T) {
		c := &core{
			currentRoundState: NewRoundState(big.NewInt(2), big.NewInt(2)),
		}

		err := c.checkMessage(big.NewInt(1), big.NewInt(2))
		if err != errOldRoundMessage {
			t.Fatalf("have %v, want %v", err, errOldRoundMessage)
		}
	})
}

func TestStoreBacklog(t *testing.T) {
	t.Run("backlog from self", func(t *testing.T) {
		addr := common.HexToAddress("0x0987654321")
		c := &core{
			logger:            log.New("backend", "test", "id", 0),
			address:           addr,
			currentRoundState: NewRoundState(big.NewInt(1), big.NewInt(2)),
			wal:               NewWalStub(),
		}

		val := validator.New(addr)
		c.storeBacklog(nil, val)

		if c.backlogs[val] != nil {
			t.Fatal("Backlog must be empty!")
		}
	})

	t.Run("vote message received", func(t *testing.T) {
		c := &core{
			logger:            log.New("backend", "test", "id", 0),
			address:           common.HexToAddress("0x1234567890"),
			currentRoundState: NewRoundState(big.NewInt(1), big.NewInt(2)),
			backlogs:          make(map[validator.Validator]*prque.Prque),
			wal:               NewWalStub(),
		}

		vote := &Vote{
			Round:  big.NewInt(1),
			Height: big.NewInt(2),
		}

		votePayload, err := Encode(vote)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &Message{
			Code: msgPrevote,
			Msg:  votePayload,
		}

		val := validator.New(common.HexToAddress("0x0987654321"))
		c.storeBacklog(msg, val)

		pque := c.backlogs[val]

		savedMsg, _ := pque.Pop()
		if !reflect.DeepEqual(msg, savedMsg) {
			t.Fatalf("Expected message %+v, but got %+v", msg, savedMsg)
		}
	})

	t.Run("proposal message received", func(t *testing.T) {
		c := &core{
			logger:            log.New("backend", "test", "id", 0),
			address:           common.HexToAddress("0x1234567890"),
			backlogs:          make(map[validator.Validator]*prque.Prque),
			currentRoundState: NewRoundState(big.NewInt(1), big.NewInt(2)),
			wal:               NewWalStub(),
		}

		proposal := &Proposal{
			Round:         big.NewInt(1),
			Height:        big.NewInt(2),
			ValidRound:    big.NewInt(1),
			ProposalBlock: types.NewBlockWithHeader(&types.Header{}),
		}

		proposalPayload, err := Encode(proposal)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &Message{
			Code: msgProposal,
			Msg:  proposalPayload,
		}

		val := validator.New(common.HexToAddress("0x0987654321"))

		c.storeBacklog(msg, val)
		pque := c.backlogs[val]

		savedMsg, _ := pque.Pop()
		if !reflect.DeepEqual(msg, savedMsg) {
			t.Fatalf("Expected message %+v, but got %+v", msg, savedMsg)
		}
	})
}

func TestProcessBacklog(t *testing.T) {
	t.Run("valid proposal received", func(t *testing.T) {
		proposal := &Proposal{
			Round:         big.NewInt(1),
			Height:        big.NewInt(2),
			ValidRound:    big.NewInt(1),
			ProposalBlock: types.NewBlockWithHeader(&types.Header{}),
		}

		proposalPayload, err := Encode(proposal)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &Message{
			Code: msgProposal,
			Msg:  proposalPayload,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		valSet := newTestValidatorSet(1)
		val := valSet.GetByIndex(0)

		expected := backlogEvent{
			src: val,
			msg: msg,
		}

		evChan := make(chan interface{}, 1)

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Post(expected).Do(func(ev interface{}) {
			evChan <- ev
		})

		c := &core{
			logger:            log.New("backend", "test", "id", 0),
			backend:           backendMock,
			address:           common.HexToAddress("0x1234567890"),
			backlogs:          make(map[validator.Validator]*prque.Prque),
			currentRoundState: NewRoundState(big.NewInt(1), big.NewInt(2)),
			wal:               NewWalStub(),
		}

		c.storeBacklog(msg, val)
		c.processBacklog()

		timeout := time.NewTimer(2 * time.Second)
		select {
		case ev := <-evChan:
			e, ok := ev.(backlogEvent)
			if !ok {
				t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev))
			}
			if e.msg.Code != msg.Code {
				t.Errorf("message code mismatch: have %v, want %v", e.msg.Code, msg.Code)
			}
		case <-timeout.C:
			t.Error("unexpected timeout occurs")
		}
	})

	t.Run("valid vote received", func(t *testing.T) {
		vote := &Vote{
			Round:  big.NewInt(1),
			Height: big.NewInt(2),
		}

		votePayload, err := Encode(vote)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &Message{
			Code: msgPrevote,
			Msg:  votePayload,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		valSet := newTestValidatorSet(1)
		val := valSet.GetByIndex(0)

		expected := backlogEvent{
			src: val,
			msg: msg,
		}

		evChan := make(chan interface{}, 1)

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Post(expected).Do(func(ev interface{}) {
			evChan <- ev
		})

		c := &core{
			logger:            log.New("backend", "test", "id", 0),
			backend:           backendMock,
			address:           common.HexToAddress("0x1234567890"),
			backlogs:          make(map[validator.Validator]*prque.Prque),
			currentRoundState: NewRoundState(big.NewInt(1), big.NewInt(2)),
			wal:               NewWalStub(),
		}

		c.storeBacklog(msg, val)
		c.processBacklog()

		timeout := time.NewTimer(2 * time.Second)
		select {
		case ev := <-evChan:
			e, ok := ev.(backlogEvent)
			if !ok {
				t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev))
			}
			if e.msg.Code != msg.Code {
				t.Errorf("message code mismatch: have %v, want %v", e.msg.Code, msg.Code)
			}
		case <-timeout.C:
			t.Error("unexpected timeout occurs")
		}
	})

	t.Run("same height, but old round", func(t *testing.T) {
		nilRoundVote := &Vote{
			Round:  big.NewInt(0),
			Height: big.NewInt(0),
		}

		nilRoundVotePayload, err := Encode(nilRoundVote)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &Message{
			Code: msgPrevote,
			Msg:  nilRoundVotePayload,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		valSet := newTestValidatorSet(1)
		val := valSet.GetByIndex(0)

		c := &core{
			logger:            log.New("backend", "test", "id", 0),
			backend:           backendMock,
			address:           common.HexToAddress("0x1234567890"),
			backlogs:          make(map[validator.Validator]*prque.Prque),
			currentRoundState: NewRoundState(big.NewInt(1), big.NewInt(0)),
			wal:               NewWalStub(),
		}

		c.storeBacklog(msg, val)
		c.processBacklog()
	})

	t.Run("future message", func(t *testing.T) {
		nilRoundVote := &Vote{
			Round:  big.NewInt(2),
			Height: big.NewInt(4),
		}

		nilRoundVotePayload, err := Encode(nilRoundVote)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &Message{
			Code: msgPrevote,
			Msg:  nilRoundVotePayload,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		valSet := newTestValidatorSet(2)
		val := valSet.GetByIndex(0)

		c := &core{
			logger:            log.New("backend", "test", "id", 0),
			backend:           backendMock,
			address:           common.HexToAddress("0x1234567890"),
			backlogs:          make(map[validator.Validator]*prque.Prque),
			currentRoundState: NewRoundState(big.NewInt(2), big.NewInt(3)),
			wal:               NewWalStub(),
		}

		c.storeBacklog(msg, val)
		c.processBacklog()
	})
}
