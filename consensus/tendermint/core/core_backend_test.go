package core

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/p2p"
	"github.com/clearmatics/autonity/rpc"
)

func TestCore_VerifySeal(t *testing.T) {
	t.Run("valid params given, no errors returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		header := &types.Header{}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().VerifySeal(nil, header).Return(nil)

		c := &core{
			backend: backendMock,
		}

		err := c.VerifySeal(nil, header)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}
	})

	t.Run("invalid params given, errors returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expected := errors.New("some error")
		header := &types.Header{}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Prepare(nil, header).Return(expected)

		c := &core{
			backend: backendMock,
		}

		err := c.Prepare(nil, header)
		if err != expected {
			t.Fatalf("Expected %v, got %v", expected, err)
		}
	})
}

func TestCore_FinalizeAndAssemble(t *testing.T) {
	t.Run("valid params given, no errors returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		header := &types.Header{}
		block := types.NewBlockWithHeader(header)

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().FinalizeAndAssemble(nil, header, nil, nil, nil, nil).
			Return(block, nil)

		c := &core{
			backend: backendMock,
		}

		b, err := c.FinalizeAndAssemble(nil, header, nil, nil, nil, nil)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		if !reflect.DeepEqual(b, block) {
			t.Fatalf("Expected %v, got %v", block, b)
		}
	})

	t.Run("invalid params given, error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		header := &types.Header{}
		block := types.NewBlockWithHeader(header)
		expected := errors.New("some error")

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().FinalizeAndAssemble(nil, header, nil, nil, nil, nil).
			Return(block, expected)

		c := &core{
			backend: backendMock,
		}

		_, err := c.FinalizeAndAssemble(nil, header, nil, nil, nil, nil)
		if err != expected {
			t.Fatalf("Expected %v, got %v", expected, err)
		}
	})
}

func TestCore_Seal(t *testing.T) {
	t.Run("valid params given, no errors returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		header := &types.Header{}
		block := types.NewBlockWithHeader(header)

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Seal(nil, block, nil, nil).Return(nil)

		c := &core{
			backend: backendMock,
		}

		err := c.Seal(nil, block, nil, nil)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}
	})

	t.Run("invalid params given, error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		header := &types.Header{}
		block := types.NewBlockWithHeader(header)
		expected := errors.New("some error")

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Seal(nil, block, nil, nil).Return(expected)

		c := &core{
			backend: backendMock,
		}

		err := c.Seal(nil, block, nil, nil)
		if err != expected {
			t.Fatalf("Expected %v, got %v", expected, err)
		}
	})
}

func TestCore_SealHash(t *testing.T) {
	t.Run("valid params given, hash returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		header := &types.Header{}
		expected := common.HexToHash("0x0123456789")

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().SealHash(header).Return(expected)

		c := &core{
			backend: backendMock,
		}

		hash := c.SealHash(header)
		if hash != expected {
			t.Fatalf("Expected %v, got %v", expected, hash)
		}
	})
}

func TestCore_CalcDifficulty(t *testing.T) {
	t.Run("valid params given, difficulty returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		parent := &types.Header{}
		expected := big.NewInt(123456789)

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().CalcDifficulty(nil, uint64(123), parent).Return(expected)

		c := &core{
			backend: backendMock,
		}

		difficulty := c.CalcDifficulty(nil, uint64(123), parent)
		if difficulty != expected {
			t.Fatalf("Expected %v, got %v", expected, difficulty)
		}
	})
}

func TestCore_APIs(t *testing.T) {
	t.Run("valid params given, APIs returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var expected []rpc.API

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().APIs(nil).Return(expected)

		c := &core{
			backend: backendMock,
		}

		APIS := c.APIs(nil)
		if !reflect.DeepEqual(APIS, expected) {
			t.Fatalf("Expected %v, got %v", expected, APIS)
		}
	})
}

func TestCore_NewChainHead(t *testing.T) {
	t.Run("backend method called", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().NewChainHead().Return(nil)

		c := &core{
			backend: backendMock,
		}

		err := c.NewChainHead()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}
	})

	t.Run("backend error occurred, error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expected := errors.New("some error")

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().NewChainHead().Return(expected)

		c := &core{
			backend: backendMock,
		}

		err := c.NewChainHead()
		if err != expected {
			t.Fatalf("Expected %v, got %v", expected, err)
		}
	})
}

func TestCore_HandleMsg(t *testing.T) {
	t.Run("backend method called", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addr := common.HexToAddress("0x0123456789")
		data := p2p.Msg{}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().HandleMsg(addr, data).Return(true, nil)

		c := &core{
			backend: backendMock,
		}

		r, err := c.HandleMsg(addr, data)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		if r != true {
			t.Fatalf("Expected <true>, got %v", r)
		}
	})

	t.Run("backend error occurred, error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addr := common.HexToAddress("0x0123456789")
		data := p2p.Msg{}
		expected := errors.New("some error")

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().HandleMsg(addr, data).Return(false, expected)

		c := &core{
			backend: backendMock,
		}

		_, err := c.HandleMsg(addr, data)
		if err != expected {
			t.Fatalf("Expected %v, got %v", expected, err)
		}
	})
}

func TestCore_SetBroadcaster(t *testing.T) {
	t.Run("backend method called", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().SetBroadcaster(nil)

		c := &core{
			backend: backendMock,
		}

		c.SetBroadcaster(nil)
	})
}

func TestCore_Protocol(t *testing.T) {
	t.Run("backend method called", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		name := "test"
		code := uint64(123)

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Protocol().Return(name, code)

		c := &core{
			backend: backendMock,
		}

		n, cd := c.Protocol()
		if n != name {
			t.Fatalf("Expected %v, got %v", name, n)
		}

		if cd != code {
			t.Fatalf("Expected %v, got %v", code, cd)
		}
	})
}

func TestCore_ResetPeerCache(t *testing.T) {
	t.Run("backend method called", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addr := common.HexToAddress("0x0123456789")

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().ResetPeerCache(addr)

		c := &core{
			backend: backendMock,
		}

		c.ResetPeerCache(addr)
	})
}

func TestCore_SyncPeer(t *testing.T) {
	t.Run("backend method called", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addr := common.HexToAddress("0x0123456789")
		curRoundState := NewRoundMessages(big.NewInt(2), big.NewInt(1))

		val := committee.NewMockValidator(ctrl)

		valSetMock := committee.NewMockSet(ctrl)
		valSetMock.EXPECT().GetByAddress(addr).Return(1, val)

		valSet := &validatorSet{
			Set: valSetMock,
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().SyncPeer(addr, gomock.Any())

		c := &core{
			backend:          backendMock,
			curRoundMessages: curRoundState,
			committeeSet:     valSet,
		}

		c.SyncPeer(addr)
	})
}

func TestCore_Close(t *testing.T) {
	t.Run("backend method called", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Close()

		_, cancel := context.WithCancel(context.Background())

		evmux := new(event.TypeMux)

		messageEventSub := evmux.Subscribe(events.MessageEvent{}, backlogEvent{})
		newUnminedBlockEventSub := evmux.Subscribe(events.NewUnminedBlockEvent{})
		committedSub := evmux.Subscribe(events.CommitEvent{})
		timeoutEventSub := evmux.Subscribe(TimeoutEvent{})
		syncEventSub := evmux.Subscribe(events.SyncEvent{})

		stopped := make(chan struct{}, 2)
		stopped <- struct{}{}
		stopped <- struct{}{}

		logger := log.New("backend", "test", "id", 0)

		isStarted := uint32(1)

		c := &core{
			backend:                 backendMock,
			cancel:                  cancel,
			isStarting:              new(uint32),
			isStarted:               &isStarted,
			isStopping:              new(uint32),
			isStopped:               new(uint32),
			committedSub:            committedSub,
			logger:                  logger,
			messageEventSub:         messageEventSub,
			newUnminedBlockEventSub: newUnminedBlockEventSub,
			proposeTimeout:          newTimeout(propose, logger),
			prevoteTimeout:          newTimeout(prevote, logger),
			precommitTimeout:        newTimeout(precommit, logger),
			timeoutEventSub:         timeoutEventSub,
			syncEventSub:            syncEventSub,
			stopped:                 stopped,
		}

		err := c.Close()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}
	})

	t.Run("backend error occurred, error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expected := errors.New("some error")

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Close().Return(expected)

		_, cancel := context.WithCancel(context.Background())

		evmux := new(event.TypeMux)

		messageEventSub := evmux.Subscribe(events.MessageEvent{}, backlogEvent{})
		newUnminedBlockEventSub := evmux.Subscribe(events.NewUnminedBlockEvent{})
		committedSub := evmux.Subscribe(events.CommitEvent{})
		timeoutEventSub := evmux.Subscribe(TimeoutEvent{})
		syncEventSub := evmux.Subscribe(events.SyncEvent{})

		stopped := make(chan struct{}, 2)
		stopped <- struct{}{}
		stopped <- struct{}{}

		isStarted := new(uint32)
		*isStarted = 1

		logger := log.New("backend", "test", "id", 0)
		c := &core{
			backend:                 backendMock,
			cancel:                  cancel,
			isStarting:              new(uint32),
			isStarted:               isStarted,
			isStopping:              new(uint32),
			isStopped:               new(uint32),
			committedSub:            committedSub,
			logger:                  logger,
			messageEventSub:         messageEventSub,
			newUnminedBlockEventSub: newUnminedBlockEventSub,
			proposeTimeout:          newTimeout(propose, logger),
			prevoteTimeout:          newTimeout(prevote, logger),
			precommitTimeout:        newTimeout(precommit, logger),
			timeoutEventSub:         timeoutEventSub,
			syncEventSub:            syncEventSub,
			stopped:                 stopped,
		}

		err := c.Close()
		if err != expected {
			t.Fatalf("Expected %v, got %v", expected, err)
		}
	})

	t.Run("the system is already stopped, nothing done", func(t *testing.T) {
		isStopped := new(uint32)
		isStarted := new(uint32)
		atomic.StoreUint32(isStopped, 1)

		c := &core{
			isStarted: isStarted,
			isStopped: isStopped,
		}

		err := c.Stop()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}
	})

	t.Run("the system is being stopped, nothing done", func(t *testing.T) {
		isStarted := new(uint32)
		atomic.StoreUint32(isStarted, 1)

		isStopping := new(uint32)
		atomic.StoreUint32(isStopping, 1)

		c := &core{
			isStarted:  isStarted,
			isStopped:  new(uint32),
			isStopping: isStopping,
		}

		err := c.Stop()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}
	})
}
