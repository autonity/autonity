package core

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/influxdata/influxdb/pkg/deep"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/log"
)

var (
	testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testAddr   = crypto.PubkeyToAddress(testKey.PublicKey)
)

func TestCheckMessage(t *testing.T) {
	t.Run("valid params given, nil returned", func(t *testing.T) {
		c := &Core{
			round:  1,
			height: big.NewInt(2),
		}

		err := c.checkMessage(1, 2)
		require.NoError(t, err)
	})

	t.Run("given future height, error returned", func(t *testing.T) {
		c := &Core{
			round:  1,
			height: big.NewInt(2),
		}

		err := c.checkMessage(2, 4)
		if !errors.Is(err, constants.ErrFutureHeightMessage) {
			t.Fatalf("have %v, want %v", err, constants.ErrFutureHeightMessage)
		}
	})

	t.Run("given old height, error returned", func(t *testing.T) {
		c := &Core{
			round:  1,
			height: big.NewInt(2),
		}

		err := c.checkMessage(2, 1)
		if !errors.Is(err, constants.ErrOldHeightMessage) {
			t.Fatalf("have %v, want %v", err, constants.ErrOldHeightMessage)
		}
	})

	t.Run("given future round, error returned", func(t *testing.T) {
		c := &Core{
			round:  1,
			height: big.NewInt(3),
		}

		err := c.checkMessage(2, 3)
		if !errors.Is(err, constants.ErrFutureRoundMessage) {
			t.Fatalf("have %v, want %v", err, constants.ErrFutureRoundMessage)
		}
	})

	t.Run("given old round, error returned", func(t *testing.T) {
		c := &Core{
			round:  2,
			height: big.NewInt(2),
		}

		err := c.checkMessage(1, 2)
		if !errors.Is(err, constants.ErrOldRoundMessage) {
			t.Fatalf("have %v, want %v", err, constants.ErrOldRoundMessage)
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

	t.Run("valid vote received", func(t *testing.T) {
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

		c := &Core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			address:          common.HexToAddress("0x1234567890"),
			backlogs:         make(map[common.Address][]message.Msg),
			step:             Propose,
			round:            1,
			height:           big.NewInt(2),
			proposeTimeout:   NewTimeout(Propose, log.New("ProposeTimeout")),
			prevoteTimeout:   NewTimeout(Prevote, log.New("PrevoteTimeout")),
			precommitTimeout: NewTimeout(Precommit, log.New("PrecommitTimeout")),
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

	t.Run("future height message are not processed", func(t *testing.T) {
		msg := message.NewPrevote(2, 4, common.Hash{}, defaultSigner)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(0)

		committeeSet := NewTestCommitteeSet(2)
		val, _ := committeeSet.GetByIndex(0)

		c := &Core{
			logger:   log.New("backend", "test", "id", 0),
			backend:  backendMock,
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]message.Msg),
			round:    2,
			height:   big.NewInt(3),
		}

		c.setLastHeader(&types.Header{Committee: committeeSet.Committee()})

		c.storeBacklog(msg, val.Address)
		c.processBacklog()
	})

	t.Run("future height message are processed when height change", func(t *testing.T) {
		msg := message.NewPrevote(2, 4, common.Hash{}, defaultSigner)
		msg2 := message.NewPrecommit(2, 4, common.Hash{}, defaultSigner)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(0)

		committeeSet := NewTestCommitteeSet(2)
		val, _ := committeeSet.GetByIndex(0)

		c := &Core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			address:          common.HexToAddress("0x1234567890"),
			backlogs:         make(map[common.Address][]message.Msg),
			round:            2,
			height:           big.NewInt(3),
			proposeTimeout:   NewTimeout(Propose, log.New("ProposeTimeout")),
			prevoteTimeout:   NewTimeout(Prevote, log.New("PrevoteTimeout")),
			precommitTimeout: NewTimeout(Precommit, log.New("PrecommitTimeout")),
			committee:        committeeSet,
			messages:         message.NewMap(),
		}
		c.curRoundMessages = c.messages.GetOrCreate(2)

		c.setLastHeader(&types.Header{Committee: committeeSet.Committee()})

		c.storeBacklog(msg, val.Address)
		c.storeBacklog(msg2, val.Address)
		c.SetStep(context.Background(), Prevote)
		c.processBacklog()
		c.setHeight(big.NewInt(4))

		backendMock.EXPECT().Post(gomock.Any()).Times(2)
		c.SetStep(context.Background(), Prevote)
		c.processBacklog()
		timeout := time.NewTimer(2 * time.Second)
		<-timeout.C
	})

	t.Run("untrusted messages are processed when height change", func(t *testing.T) {
		msg := message.NewPrevote(2, 4, common.Hash{}, defaultSigner)
		msg2 := message.NewPrecommit(2, 4, common.Hash{}, defaultSigner)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)

		committeeSet := NewTestCommitteeSet(2)

		c := &Core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			address:          common.HexToAddress("0x1234567890"),
			backlogs:         make(map[common.Address][]message.Msg),
			backlogUntrusted: map[uint64][]message.Msg{},
			round:            2,
			height:           big.NewInt(3),
			proposeTimeout:   NewTimeout(Propose, log.New("ProposeTimeout")),
			prevoteTimeout:   NewTimeout(Prevote, log.New("PrevoteTimeout")),
			precommitTimeout: NewTimeout(Precommit, log.New("PrecommitTimeout")),
			committee:        committeeSet,
			messages:         message.NewMap(),
		}
		c.curRoundMessages = c.messages.GetOrCreate(2)

		c.setLastHeader(&types.Header{Committee: committeeSet.Committee()})

		backendMock.EXPECT().Post(gomock.Any()).Times(0)
		c.storeFutureMessage(msg)
		c.storeFutureMessage(msg2)
		c.SetStep(context.Background(), Prevote)
		c.processBacklog()
		c.setHeight(big.NewInt(4))

		backendMock.EXPECT().Post(gomock.Any()).Times(2)
		c.SetStep(context.Background(), Prevote)

		backendMock.EXPECT().Post(gomock.Any()).Times(0)
		c.processBacklog()
		<-time.NewTimer(2 * time.Second).C
	})

	t.Run("future round message are processed when round change", func(t *testing.T) {
		msg := message.NewPrevote(2, 4, common.Hash{}, defaultSigner)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(0)
		committeeSet := NewTestCommitteeSet(2)
		val, err := committeeSet.GetByIndex(0)
		require.NoError(t, err)

		c := New(backendMock, nil, common.HexToAddress("0x1234567890"), log.Root())
		c.setRound(1)
		c.setHeight(big.NewInt(4))
		c.step = Prevote
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

func TestStoreUncheckedBacklog(t *testing.T) {
	t.Run("save messages in the untrusted backlog", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		backendMock := interfaces.NewMockBackend(ctrl)
		c := &Core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			address:          common.HexToAddress("0x1234567890"),
			backlogs:         make(map[common.Address][]message.Msg),
			backlogUntrusted: make(map[uint64][]message.Msg),
			step:             Prevote,
			round:            1,
			height:           big.NewInt(4),
		}
		var messages []message.Msg
		for i := int64(0); i < MaxSizeBacklogUnchecked; i++ {
			msg := message.NewPrevote(
				i%10,
				uint64(i/(1+i%10)),
				common.Hash{},
				defaultSigner)
			c.storeFutureMessage(msg)
			messages = append(messages, msg)
		}
		found := 0
		for _, msg := range messages {
			for _, umsg := range c.backlogUntrusted[msg.H()] {
				if deep.Equal(msg, umsg) {
					found++
				}
			}
		}
		if found != MaxSizeBacklogUnchecked {
			t.Fatal("unchecked messages lost")
		}
	})

	t.Run("excess messages are removed from the untrusted backlog", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)

		c := &Core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			address:          common.HexToAddress("0x1234567890"),
			backlogs:         make(map[common.Address][]message.Msg),
			backlogUntrusted: make(map[uint64][]message.Msg),
			step:             Prevote,
			round:            1,
			height:           big.NewInt(4),
		}

		var messages []message.Msg
		uncheckedFounds := make(map[uint64]struct{})
		backendMock.EXPECT().RemoveMessageFromLocalCache(gomock.Any()).Times(MaxSizeBacklogUnchecked).Do(func(msg message.Msg) {
			if _, ok := uncheckedFounds[msg.H()]; ok {
				t.Fatal("duplicate message received")
			}
			uncheckedFounds[msg.H()] = struct{}{}
		})

		for i := int64(2 * MaxSizeBacklogUnchecked); i > 0; i-- {
			prevote := message.NewPrevote(i%10, uint64(i), common.Hash{}, defaultSigner)
			c.storeFutureMessage(prevote)
			if i < MaxSizeBacklogUnchecked {
				messages = append(messages, prevote)
			}
		}

		found := 0
		for _, msg := range messages {
			for _, umsg := range c.backlogUntrusted[msg.H()] {
				if deep.Equal(msg, umsg) {
					found++
				}
			}
		}
		if found != MaxSizeBacklogUnchecked-1 {
			t.Fatal("unchecked messages lost")
		}
		if len(uncheckedFounds) != MaxSizeBacklogUnchecked {
			t.Fatal("unchecked messages lost")
		}
	})
}
