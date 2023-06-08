package core

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/helpers"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	tctypes "github.com/autonity/autonity/consensus/tendermint/core/types"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/rlp"
	"github.com/golang/mock/gomock"
	"github.com/influxdata/influxdb/pkg/deep"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"testing"
	"time"
)

func TestCheckMessage(t *testing.T) {
	t.Run("valid params given, nil returned", func(t *testing.T) {
		c := &Core{
			round:  1,
			height: big.NewInt(2),
		}

		err := c.CheckMessage(1, big.NewInt(2), tctypes.Propose)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}
	})

	t.Run("given nil height, error returned", func(t *testing.T) {
		c := &Core{
			round:  1,
			height: big.NewInt(2),
		}

		err := c.CheckMessage(1, nil, tctypes.Propose)
		if err != constants.ErrInvalidMessage {
			t.Fatalf("have %v, want %v", err, constants.ErrInvalidMessage)
		}
	})

	t.Run("given future height, error returned", func(t *testing.T) {
		c := &Core{
			round:  1,
			height: big.NewInt(2),
		}

		err := c.CheckMessage(2, big.NewInt(4), tctypes.Propose)
		if err != constants.ErrFutureHeightMessage {
			t.Fatalf("have %v, want %v", err, constants.ErrFutureHeightMessage)
		}
	})

	t.Run("given old height, error returned", func(t *testing.T) {
		c := &Core{
			round:  1,
			height: big.NewInt(2),
		}

		err := c.CheckMessage(2, big.NewInt(1), tctypes.Propose)
		if err != constants.ErrOldHeightMessage {
			t.Fatalf("have %v, want %v", err, constants.ErrOldHeightMessage)
		}
	})

	t.Run("given future round, error returned", func(t *testing.T) {
		c := &Core{
			round:  1,
			height: big.NewInt(3),
		}

		err := c.CheckMessage(2, big.NewInt(3), tctypes.Propose)
		if err != constants.ErrFutureRoundMessage {
			t.Fatalf("have %v, want %v", err, constants.ErrFutureRoundMessage)
		}
	})

	t.Run("given old round, error returned", func(t *testing.T) {
		c := &Core{
			round:  2,
			height: big.NewInt(2),
		}

		err := c.CheckMessage(1, big.NewInt(2), tctypes.Propose)
		if err != constants.ErrOldRoundMessage {
			t.Fatalf("have %v, want %v", err, constants.ErrOldRoundMessage)
		}
	})

	t.Run("at propose step, given prevote for same view, error returned", func(t *testing.T) {
		c := &Core{
			round:  2,
			height: big.NewInt(2),
			step:   tctypes.Propose,
		}

		err := c.CheckMessage(2, big.NewInt(2), tctypes.Prevote)
		if err != constants.ErrFutureStepMessage {
			t.Fatalf("have %v, want %v", err, constants.ErrFutureStepMessage)
		}
	})

	t.Run("at propose step, given precommit for same view, error returned", func(t *testing.T) {
		c := &Core{
			round:  2,
			height: big.NewInt(2),
			step:   tctypes.Propose,
		}

		err := c.CheckMessage(2, big.NewInt(2), tctypes.Precommit)
		if err != constants.ErrFutureStepMessage {
			t.Fatalf("have %v, want %v", err, constants.ErrFutureStepMessage)
		}
	})

	t.Run("at prevote step, given precommit for same view, no error returned", func(t *testing.T) {
		c := &Core{
			round:  2,
			height: big.NewInt(2),
			step:   tctypes.Prevote,
		}

		err := c.CheckMessage(2, big.NewInt(2), tctypes.Precommit)
		if err != nil {
			t.Fatalf("have %v, want %v", err, nil)
		}
	})

	t.Run("at precommit step, given prevote for same view, no error returned", func(t *testing.T) {
		c := &Core{
			round:  2,
			height: big.NewInt(2),
			step:   tctypes.Precommit,
		}

		err := c.CheckMessage(2, big.NewInt(2), tctypes.Prevote)
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
			step:    tctypes.Propose,
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
			backlogs: make(map[common.Address][]*messageutils.Message),
		}

		vote := &messageutils.Vote{
			Round:  1,
			Height: big.NewInt(2),
		}

		votePayload, err := messageutils.Encode(vote)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &messageutils.Message{
			Code:         consensus.MsgPrevote,
			TbftMsgBytes: votePayload,
			DecodedMsg:   vote,
		}

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

	t.Run("proposal message received", func(t *testing.T) {
		c := &Core{
			logger:   log.New("backend", "test", "id", 0),
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]*messageutils.Message),
		}

		proposal := &messageutils.Proposal{
			Round:         1,
			Height:        big.NewInt(2),
			ValidRound:    1,
			ProposalBlock: types.NewBlockWithHeader(&types.Header{}),
		}

		proposalPayload, err := messageutils.Encode(proposal)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &messageutils.Message{
			Code:         consensus.MsgProposal,
			TbftMsgBytes: proposalPayload,
			DecodedMsg:   proposal,
		}

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
		proposal := &messageutils.Proposal{
			Round:         1,
			Height:        big.NewInt(2),
			ValidRound:    1,
			ProposalBlock: types.NewBlockWithHeader(&types.Header{}),
		}

		proposalPayload, err := messageutils.Encode(proposal)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &messageutils.Message{
			Code:         consensus.MsgProposal,
			TbftMsgBytes: proposalPayload,
			DecodedMsg:   proposal,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		committeeSet := helpers.NewTestCommitteeSet(1)
		val, _ := committeeSet.GetByIndex(0)

		expected := backlogEvent{
			msg: msg,
		}

		evChan := make(chan interface{}, 1)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(expected).Do(func(ev interface{}) {
			evChan <- ev
		})

		c := &Core{
			logger:   log.New("backend", "test", "id", 0),
			backend:  backendMock,
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]*messageutils.Message),
			step:     tctypes.Propose,
			round:    1,
			height:   big.NewInt(2),
		}

		c.setLastHeader(&types.Header{Committee: committeeSet.Committee()})

		c.storeBacklog(msg, val.Address)
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
			t.Error("unexpected Timeout occurs")
		}
	})

	t.Run("valid vote received, processed at prevote step", func(t *testing.T) {
		vote := &messageutils.Vote{
			Round:  1,
			Height: big.NewInt(2),
		}

		votePayload, err := messageutils.Encode(vote)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &messageutils.Message{
			Code:         consensus.MsgPrevote,
			TbftMsgBytes: votePayload,
			DecodedMsg:   vote,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		committeeSet := helpers.NewTestCommitteeSet(1)
		val, _ := committeeSet.GetByIndex(0)

		expected := backlogEvent{
			msg: msg,
		}

		evChan := make(chan interface{}, 1)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(expected).Do(func(ev interface{}) {
			evChan <- ev
		})

		c := &Core{
			logger:   log.New("backend", "test", "id", 0),
			backend:  backendMock,
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]*messageutils.Message),
			step:     tctypes.Propose,
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
		c.SetStep(tctypes.Prevote)
		c.processBacklog()

		timeout = time.NewTimer(2 * time.Second)
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
			t.Error("unexpected Timeout occurs")
		}
	})

	t.Run("same height, but old round", func(t *testing.T) {
		nilRoundVote := &messageutils.Vote{
			Round:  0,
			Height: big.NewInt(1),
		}

		nilRoundVotePayload, err := messageutils.Encode(nilRoundVote)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &messageutils.Message{
			Code:         consensus.MsgPrevote,
			TbftMsgBytes: nilRoundVotePayload,
			DecodedMsg:   nilRoundVote,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		committeeSet := helpers.NewTestCommitteeSet(1)
		val, _ := committeeSet.GetByIndex(0)

		expected := backlogEvent{
			msg: msg,
		}

		evChan := make(chan interface{}, 1)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(expected).Do(func(ev interface{}) {
			evChan <- ev
		})

		c := &Core{
			logger:   log.New("backend", "test", "id", 0),
			backend:  backendMock,
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]*messageutils.Message),
			round:    1,
			step:     tctypes.Prevote,
			height:   big.NewInt(1),
		}

		c.setLastHeader(&types.Header{Committee: committeeSet.Committee()})

		c.storeBacklog(msg, val.Address)
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
			t.Error("unexpected Timeout occurs")
		}
	})

	t.Run("future height message are not processed", func(t *testing.T) {
		nilRoundVote := &messageutils.Vote{
			Round:  2,
			Height: big.NewInt(4),
		}

		nilRoundVotePayload, err := messageutils.Encode(nilRoundVote)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &messageutils.Message{
			Code:         consensus.MsgPrevote,
			TbftMsgBytes: nilRoundVotePayload,
			DecodedMsg:   nilRoundVote,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(0)

		committeeSet := helpers.NewTestCommitteeSet(2)
		val, _ := committeeSet.GetByIndex(0)

		c := &Core{
			logger:   log.New("backend", "test", "id", 0),
			backend:  backendMock,
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]*messageutils.Message),
			round:    2,
			height:   big.NewInt(3),
		}

		c.setLastHeader(&types.Header{Committee: committeeSet.Committee()})

		c.storeBacklog(msg, val.Address)
		c.processBacklog()
	})

	t.Run("future height message are processed when height change", func(t *testing.T) {
		nilRoundVote := &messageutils.Vote{
			Round:  2,
			Height: big.NewInt(4),
		}

		nilRoundVotePayload, err := messageutils.Encode(nilRoundVote)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &messageutils.Message{
			Code:         consensus.MsgPrevote,
			TbftMsgBytes: nilRoundVotePayload,
			DecodedMsg:   nilRoundVote,
		}
		msg2 := &messageutils.Message{
			Code:         consensus.MsgPrecommit,
			TbftMsgBytes: nilRoundVotePayload,
			DecodedMsg:   nilRoundVote,
		}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(0)

		committeeSet := helpers.NewTestCommitteeSet(2)
		val, _ := committeeSet.GetByIndex(0)

		c := &Core{
			logger:   log.New("backend", "test", "id", 0),
			backend:  backendMock,
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]*messageutils.Message),
			round:    2,
			height:   big.NewInt(3),
		}

		c.setLastHeader(&types.Header{Committee: committeeSet.Committee()})

		c.storeBacklog(msg, val.Address)
		c.storeBacklog(msg2, val.Address)
		c.SetStep(tctypes.Prevote)
		c.processBacklog()
		c.setHeight(big.NewInt(4))

		backendMock.EXPECT().Post(gomock.Any()).Times(2)
		c.SetStep(tctypes.Prevote)
		c.processBacklog()
		timeout := time.NewTimer(2 * time.Second)
		<-timeout.C
	})

	t.Run("untrusted messages are processed when height change", func(t *testing.T) {
		nilRoundVote := &messageutils.Vote{
			Round:  2,
			Height: big.NewInt(4),
		}

		nilRoundVotePayload, err := messageutils.Encode(nilRoundVote)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &messageutils.Message{
			Code:         consensus.MsgPrevote,
			TbftMsgBytes: nilRoundVotePayload,
			DecodedMsg:   nilRoundVote,
		}
		msg2 := &messageutils.Message{
			Code:         consensus.MsgPrecommit,
			TbftMsgBytes: nilRoundVotePayload,
			DecodedMsg:   nilRoundVote,
		}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)

		committeeSet := helpers.NewTestCommitteeSet(2)

		c := &Core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			address:          common.HexToAddress("0x1234567890"),
			backlogs:         make(map[common.Address][]*messageutils.Message),
			backlogUnchecked: map[uint64][]*messageutils.Message{},
			round:            2,
			height:           big.NewInt(3),
		}

		c.setLastHeader(&types.Header{Committee: committeeSet.Committee()})

		backendMock.EXPECT().Post(gomock.Any()).Times(0)
		c.storeUncheckedBacklog(msg)
		c.storeUncheckedBacklog(msg2)
		c.SetStep(tctypes.Prevote)
		c.processBacklog()
		c.setHeight(big.NewInt(4))

		backendMock.EXPECT().Post(gomock.Any()).Times(2)
		c.SetStep(tctypes.Prevote)

		backendMock.EXPECT().Post(gomock.Any()).Times(0)
		c.processBacklog()
		<-time.NewTimer(2 * time.Second).C
	})

	t.Run("future round message are processed when round change", func(t *testing.T) {
		nilRoundVote := &messageutils.Vote{
			Round:  2,
			Height: big.NewInt(4),
		}

		nilRoundVotePayload, err := messageutils.Encode(nilRoundVote)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &messageutils.Message{
			Code:         consensus.MsgPrevote,
			TbftMsgBytes: nilRoundVotePayload,
			DecodedMsg:   nilRoundVote,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(0)

		committeeSet := helpers.NewTestCommitteeSet(2)
		val, err := committeeSet.GetByIndex(0)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		c := &Core{
			logger:   log.New("backend", "test", "id", 0),
			backend:  backendMock,
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]*messageutils.Message),
			step:     tctypes.Prevote,
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

func TestStoreUncheckedBacklog(t *testing.T) {
	t.Run("save messages in the untrusted backlog", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)

		c := &Core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			address:          common.HexToAddress("0x1234567890"),
			backlogs:         make(map[common.Address][]*messageutils.Message),
			backlogUnchecked: make(map[uint64][]*messageutils.Message),
			step:             tctypes.Prevote,
			round:            1,
			height:           big.NewInt(4),
		}
		var messages []*messageutils.Message

		for i := int64(0); i < MaxSizeBacklogUnchecked; i++ {
			nilRoundVote := &messageutils.Vote{
				Round:  i % 10,
				Height: big.NewInt(i / (1 + i%10)),
			}
			payload, err := rlp.EncodeToBytes(nilRoundVote)
			require.NoError(t, err)
			msg := &messageutils.Message{
				Code:         consensus.MsgPrevote,
				TbftMsgBytes: payload,
				DecodedMsg:   nilRoundVote,
			}
			c.storeUncheckedBacklog(msg)
			messages = append(messages, msg)
		}
		found := 0
		for _, msg := range messages {
			height, err := msg.Height()
			if err != nil {
				t.Fatal("can't retrieve message height")
			}
			for _, umsg := range c.backlogUnchecked[height.Uint64()] {
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
			backlogs:         make(map[common.Address][]*messageutils.Message),
			backlogUnchecked: make(map[uint64][]*messageutils.Message),
			step:             tctypes.Prevote,
			round:            1,
			height:           big.NewInt(4),
		}

		var messages []*messageutils.Message
		uncheckedFounds := make(map[uint64]struct{})
		backendMock.EXPECT().RemoveMessageFromLocalCache(gomock.Any()).Times(MaxSizeBacklogUnchecked).Do(func(payload []byte) {
			var msg messageutils.Message
			err := msg.FromPayload(payload)
			if err != nil {
				t.Fatal("could not decode message payload")
			}
			height, err := msg.Height()
			if err != nil {
				t.Fatal("could not decode message height")
			}
			if _, ok := uncheckedFounds[height.Uint64()]; ok {
				t.Fatal("duplicate message received")
			}
			uncheckedFounds[height.Uint64()] = struct{}{}
		})

		for i := int64(2 * MaxSizeBacklogUnchecked); i > 0; i-- {
			nilRoundVote := &messageutils.Vote{
				Round:  i % 10,
				Height: big.NewInt(i),
			}
			payload, err := rlp.EncodeToBytes(nilRoundVote)
			require.NoError(t, err)
			msg := &messageutils.Message{
				Code:         consensus.MsgPrevote,
				TbftMsgBytes: payload,
				DecodedMsg:   nilRoundVote,
			}
			c.storeUncheckedBacklog(msg)
			if i < MaxSizeBacklogUnchecked {
				messages = append(messages, msg)
			}
		}

		found := 0
		for _, msg := range messages {
			height, err := msg.Height()
			if err != nil {
				t.Error("can't retrieve message height")
			}
			for _, umsg := range c.backlogUnchecked[height.Uint64()] {
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
