package core

import (
	"errors"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/internal/testlog"
	"github.com/autonity/autonity/log"
	"go.uber.org/mock/gomock"
)

var (
	testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testPower  = big.NewInt(1000)
	testAddr   = crypto.PubkeyToAddress(testKey.PublicKey)
)

func TestCheckMessage(t *testing.T) {
	t.Run("valid params given, nil returned", func(t *testing.T) {
		c := &Core{
			round:  1,
			height: big.NewInt(2),
		}

		err := c.checkMessageStep(1, 2, Propose)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}
	})

	t.Run("given future height, code panics", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("The code did not panic")
			}
		}()
		c := &Core{
			round:  1,
			height: big.NewInt(2),
			logger: testlog.Logger(t, log.LvlDebug),
		}

		_ = c.checkMessageStep(2, 4, Propose)
	})

	t.Run("given old height, error returned", func(t *testing.T) {
		c := &Core{
			round:  1,
			height: big.NewInt(2),
		}

		err := c.checkMessageStep(2, 1, Propose)
		if !errors.Is(err, constants.ErrOldHeightMessage) {
			t.Fatalf("have %v, want %v", err, constants.ErrOldHeightMessage)
		}
	})

	t.Run("given future round, error returned", func(t *testing.T) {
		c := &Core{
			round:  1,
			height: big.NewInt(3),
		}

		err := c.checkMessageStep(2, 3, Propose)
		if !errors.Is(err, constants.ErrFutureRoundMessage) {
			t.Fatalf("have %v, want %v", err, constants.ErrFutureRoundMessage)
		}
	})

	t.Run("given old round, error returned", func(t *testing.T) {
		c := &Core{
			round:  2,
			height: big.NewInt(2),
		}

		err := c.checkMessageStep(1, 2, Propose)
		if !errors.Is(err, constants.ErrOldRoundMessage) {
			t.Fatalf("have %v, want %v", err, constants.ErrOldRoundMessage)
		}
	})

	t.Run("at propose step, given prevote for same view, error returned", func(t *testing.T) {
		c := &Core{
			round:  2,
			height: big.NewInt(2),
			step:   Propose,
		}

		err := c.checkMessageStep(2, 2, Prevote)
		if !errors.Is(err, constants.ErrFutureStepMessage) {
			t.Fatalf("have %v, want %v", err, constants.ErrFutureStepMessage)
		}
	})

	t.Run("at propose step, given precommit for same view, error returned", func(t *testing.T) {
		c := &Core{
			round:  2,
			height: big.NewInt(2),
			step:   Propose,
		}

		err := c.checkMessageStep(2, 2, Precommit)
		if !errors.Is(err, constants.ErrFutureStepMessage) {
			t.Fatalf("have %v, want %v", err, constants.ErrFutureStepMessage)
		}
	})

	t.Run("at prevote step, given precommit for same view, no error returned", func(t *testing.T) {
		c := &Core{
			round:  2,
			height: big.NewInt(2),
			step:   Prevote,
		}

		err := c.checkMessageStep(2, 2, Precommit)
		if err != nil {
			t.Fatalf("have %v, want %v", err, nil)
		}
	})

	t.Run("at precommit step, given prevote for same view, no error returned", func(t *testing.T) {
		c := &Core{
			round:  2,
			height: big.NewInt(2),
			step:   Precommit,
		}

		err := c.checkMessageStep(2, 2, Prevote)
		if err != nil {
			t.Fatalf("have %v, want %v", err, nil)
		}
	})
}

func TestStoreBacklog(t *testing.T) {
	t.Run("backlog from self", func(t *testing.T) {
		addr := common.HexToAddress("0x0987654321")
		c := &Core{
			logger:  log.New("backend", "test", "id", 0),
			address: addr,
			height:  big.NewInt(1),
			step:    Propose,
		}

		val := types.CommitteeMember{
			Address:     addr,
			VotingPower: big.NewInt(1),
		}

		c.storeBacklog(nil, val.Address)

		if c.backlogs[val.Address] != nil {
			t.Fatal("Backlog must be empty!")
		}
	})

	t.Run("vote message received", func(t *testing.T) {
		c := &Core{
			logger:   log.New("backend", "test", "id", 0),
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]message.Msg),
		}

		vote := message.NewPrevote(1, 2, common.Hash{}, defaultSigner)
		val := types.CommitteeMember{
			Address:     common.HexToAddress("0x0987654321"),
			VotingPower: big.NewInt(1),
		}
		c.storeBacklog(vote, val.Address)

		pque := c.backlogs[val.Address]

		savedMsg := pque[0]
		if !reflect.DeepEqual(vote, savedMsg) {
			t.Fatalf("Expected message %+v, but got %+v", vote, savedMsg)
		}
	})

	t.Run("proposal message received", func(t *testing.T) {
		c := &Core{
			logger:   log.New("backend", "test", "id", 0),
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]message.Msg),
		}

		msg := message.NewPropose(1, 2, 1, types.NewBlockWithHeader(&types.Header{}), defaultSigner)
		val := types.CommitteeMember{
			Address:     common.HexToAddress("0x0987654321"),
			VotingPower: big.NewInt(1),
		}
		c.storeBacklog(msg, val.Address)
		pque := c.backlogs[val.Address]

		savedMsg := pque[0]
		if !reflect.DeepEqual(msg, savedMsg) {
			t.Fatalf("Expected message %+v, but got %+v", msg, savedMsg)
		}
	})
}

func TestProcessBacklog(t *testing.T) {
	t.Run("valid proposal received", func(t *testing.T) {

		msg := message.NewPropose(1, 2, 1, types.NewBlockWithHeader(&types.Header{}), defaultSigner)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		committeeSet := NewTestCommitteeSet(1)
		val, _ := committeeSet.GetByIndex(0)

		expected := backlogMessageEvent{
			msg: msg,
		}

		evChan := make(chan any, 1)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(expected).Do(func(ev any) {
			evChan <- ev
		})
		backendMock.EXPECT().ProcessFutureMsgs(uint64(2)).MaxTimes(1)

		c := &Core{
			logger:   log.New("backend", "test", "id", 0),
			backend:  backendMock,
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]message.Msg),
			step:     Propose,
			round:    1,
			height:   big.NewInt(2),
		}

		c.setLastHeader(&types.Header{Committee: committeeSet.Committee()})

		c.storeBacklog(msg, val.Address)
		c.processBacklog()

		timeout := time.NewTimer(2 * time.Second)
		select {
		case ev := <-evChan:
			e, ok := ev.(backlogMessageEvent)
			if !ok {
				t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev))
			}
			if e.msg.Code() != msg.Code() {
				t.Errorf("message code mismatch: have %v, want %v", e.msg.Code(), msg.Code())
			}
		case <-timeout.C:
			t.Error("unexpected Timeout occurs")
		}
	})

	t.Run("valid vote received, processed at prevote step", func(t *testing.T) {

		msg := message.NewPrevote(1, 2, common.Hash{}, defaultSigner)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		committeeSet := NewTestCommitteeSet(1)
		val, _ := committeeSet.GetByIndex(0)

		expected := backlogMessageEvent{
			msg: msg,
		}

		evChan := make(chan any, 1)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(expected).Do(func(ev any) {
			evChan <- ev
		})
		backendMock.EXPECT().ProcessFutureMsgs(uint64(2)).MaxTimes(3)

		c := &Core{
			logger:   log.New("backend", "test", "id", 0),
			backend:  backendMock,
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]message.Msg),
			step:     Propose,
			round:    1,
			height:   big.NewInt(2),
		}

		c.setLastHeader(&types.Header{Committee: committeeSet.Committee()})

		c.storeBacklog(msg, val.Address)
		c.processBacklog()

		timeout := time.NewTimer(2 * time.Second)
		//vote should not be processed at propose step
		select {
		case ev := <-evChan:
			t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev))
		case <-timeout.C:
		}
		c.SetStep(Prevote)
		c.processBacklog()

		timeout = time.NewTimer(2 * time.Second)
		select {
		case ev := <-evChan:
			e, ok := ev.(backlogMessageEvent)
			if !ok {
				t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev))
			}
			if e.msg.Code() != msg.Code() {
				t.Errorf("message code mismatch: have %v, want %v", e.msg.Code(), msg.Code())
			}
		case <-timeout.C:
			t.Error("unexpected Timeout occurs")
		}
	})

	t.Run("same height, but old round", func(t *testing.T) {
		msg := message.NewPrevote(0, 1, common.Hash{}, defaultSigner)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		committeeSet := NewTestCommitteeSet(1)
		val, _ := committeeSet.GetByIndex(0)

		expected := backlogMessageEvent{
			msg: msg,
		}

		evChan := make(chan any, 1)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(expected).Do(func(ev any) {
			evChan <- ev
		})
		backendMock.EXPECT().ProcessFutureMsgs(uint64(1)).MaxTimes(1)

		c := &Core{
			logger:   log.New("backend", "test", "id", 0),
			backend:  backendMock,
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]message.Msg),
			round:    1,
			step:     Prevote,
			height:   big.NewInt(1),
		}

		c.setLastHeader(&types.Header{Committee: committeeSet.Committee()})

		c.storeBacklog(msg, val.Address)
		c.processBacklog()

		timeout := time.NewTimer(2 * time.Second)
		select {
		case ev := <-evChan:
			e, ok := ev.(backlogMessageEvent)
			if !ok {
				t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev))
			}
			if e.msg.Code() != msg.Code() {
				t.Errorf("message code mismatch: have %v, want %v", e.msg.Code(), msg.Code())
			}
		case <-timeout.C:
			t.Error("unexpected Timeout occurs")
		}
	})

	t.Run("future round message are processed when round change", func(t *testing.T) {
		msg := message.NewPrevote(2, 4, common.Hash{}, defaultSigner)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(0)
		backendMock.EXPECT().ProcessFutureMsgs(uint64(4)).MaxTimes(2)

		committeeSet := NewTestCommitteeSet(2)
		val, err := committeeSet.GetByIndex(0)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		c := &Core{
			logger:   log.New("backend", "test", "id", 0),
			backend:  backendMock,
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]message.Msg),
			step:     Prevote,
			round:    1,
			height:   big.NewInt(4),
		}

		c.setLastHeader(&types.Header{Committee: committeeSet.Committee()})

		c.storeBacklog(msg, val.Address)
		c.processBacklog()
		backendMock.EXPECT().Post(gomock.Any()).Times(1)
		c.setRound(2)
		c.processBacklog()
		timeout := time.NewTimer(2 * time.Second)
		<-timeout.C
	})
}
