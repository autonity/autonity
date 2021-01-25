package core

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/rlp"
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
		c := &core{
			round:  1,
			height: big.NewInt(2),
		}

		err := c.checkMessage(1, big.NewInt(2), propose)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}
	})

	t.Run("given nil height, error returned", func(t *testing.T) {
		c := &core{
			round:  1,
			height: big.NewInt(2),
		}

		err := c.checkMessage(1, nil, propose)
		if err != errInvalidMessage {
			t.Fatalf("have %v, want %v", err, errInvalidMessage)
		}
	})

	t.Run("given future height, error returned", func(t *testing.T) {
		c := &core{
			round:  1,
			height: big.NewInt(2),
		}

		err := c.checkMessage(2, big.NewInt(4), propose)
		if err != errFutureHeightMessage {
			t.Fatalf("have %v, want %v", err, errFutureHeightMessage)
		}
	})

	t.Run("given old height, error returned", func(t *testing.T) {
		c := &core{
			round:  1,
			height: big.NewInt(2),
		}

		err := c.checkMessage(2, big.NewInt(1), propose)
		if err != errOldHeightMessage {
			t.Fatalf("have %v, want %v", err, errOldHeightMessage)
		}
	})

	t.Run("given future round, error returned", func(t *testing.T) {
		c := &core{
			round:  1,
			height: big.NewInt(3),
		}

		err := c.checkMessage(2, big.NewInt(3), propose)
		if err != errFutureRoundMessage {
			t.Fatalf("have %v, want %v", err, errFutureRoundMessage)
		}
	})

	t.Run("given old round, error returned", func(t *testing.T) {
		c := &core{
			round:  2,
			height: big.NewInt(2),
		}

		err := c.checkMessage(1, big.NewInt(2), propose)
		if err != errOldRoundMessage {
			t.Fatalf("have %v, want %v", err, errOldRoundMessage)
		}
	})

	t.Run("at propose step, given prevote for same view, error returned", func(t *testing.T) {
		c := &core{
			round:  2,
			height: big.NewInt(2),
			step:   propose,
		}

		err := c.checkMessage(2, big.NewInt(2), prevote)
		if err != errFutureStepMessage {
			t.Fatalf("have %v, want %v", err, errFutureStepMessage)
		}
	})

	t.Run("at propose step, given precommit for same view, error returned", func(t *testing.T) {
		c := &core{
			round:  2,
			height: big.NewInt(2),
			step:   propose,
		}

		err := c.checkMessage(2, big.NewInt(2), precommit)
		if err != errFutureStepMessage {
			t.Fatalf("have %v, want %v", err, errFutureStepMessage)
		}
	})

	t.Run("at prevote step, given precommit for same view, no error returned", func(t *testing.T) {
		c := &core{
			round:  2,
			height: big.NewInt(2),
			step:   prevote,
		}

		err := c.checkMessage(2, big.NewInt(2), precommit)
		if err != nil {
			t.Fatalf("have %v, want %v", err, nil)
		}
	})

	t.Run("at precommit step, given prevote for same view, no error returned", func(t *testing.T) {
		c := &core{
			round:  2,
			height: big.NewInt(2),
			step:   precommit,
		}

		err := c.checkMessage(2, big.NewInt(2), prevote)
		if err != nil {
			t.Fatalf("have %v, want %v", err, nil)
		}
	})
}

func TestStoreBacklog(t *testing.T) {
	t.Run("backlog from self", func(t *testing.T) {
		addr := common.HexToAddress("0x0987654321")
		c := &core{
			logger:  log.New("backend", "test", "id", 0),
			address: addr,
			height:  big.NewInt(1),
			step:    propose,
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
		c := &core{
			logger:   log.New("backend", "test", "id", 0),
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]*types.ConsensusMessage),
		}

		vote := &types.Vote{
			Round:  1,
			Height: big.NewInt(2),
		}

		votePayload, err := types.Encode(vote)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &types.ConsensusMessage{
			Code:       types.msgPrevote,
			Msg:        votePayload,
			decodedMsg: vote,
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
		c := &core{
			logger:   log.New("backend", "test", "id", 0),
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]*types.ConsensusMessage),
		}

		proposal := &types.Proposal{
			Round:         1,
			Height:        big.NewInt(2),
			ValidRound:    1,
			ProposalBlock: types.NewBlockWithHeader(&types.Header{}),
		}

		proposalPayload, err := types.Encode(proposal)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &types.ConsensusMessage{
			Code:       types.msgProposal,
			Msg:        proposalPayload,
			decodedMsg: proposal,
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
		proposal := &types.Proposal{
			Round:         1,
			Height:        big.NewInt(2),
			ValidRound:    1,
			ProposalBlock: types.NewBlockWithHeader(&types.Header{}),
		}

		proposalPayload, err := types.Encode(proposal)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &types.ConsensusMessage{
			Code:       types.msgProposal,
			Msg:        proposalPayload,
			decodedMsg: proposal,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		committeeSet := newTestCommitteeSet(1)
		val, _ := committeeSet.GetByIndex(0)

		expected := backlogEvent{
			msg: msg,
		}

		evChan := make(chan interface{}, 1)

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Post(expected).Do(func(ev interface{}) {
			evChan <- ev
		})

		c := &core{
			logger:   log.New("backend", "test", "id", 0),
			backend:  backendMock,
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]*types.ConsensusMessage),
			step:     propose,
			round:    1,
			height:   big.NewInt(2),
		}

		c.lastHeader = &types.Header{Committee: committeeSet.Committee()}

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
			t.Error("unexpected timeout occurs")
		}
	})

	t.Run("valid vote received, processed at prevote step", func(t *testing.T) {
		vote := &types.Vote{
			Round:  1,
			Height: big.NewInt(2),
		}

		votePayload, err := types.Encode(vote)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &types.ConsensusMessage{
			Code:       types.msgPrevote,
			Msg:        votePayload,
			decodedMsg: vote,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		committeeSet := newTestCommitteeSet(1)
		val, _ := committeeSet.GetByIndex(0)

		expected := backlogEvent{
			msg: msg,
		}

		evChan := make(chan interface{}, 1)

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Post(expected).Do(func(ev interface{}) {
			evChan <- ev
		})

		c := &core{
			logger:   log.New("backend", "test", "id", 0),
			backend:  backendMock,
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]*types.ConsensusMessage),
			step:     propose,
			round:    1,
			height:   big.NewInt(2),
		}

		c.lastHeader = &types.Header{Committee: committeeSet.Committee()}

		c.storeBacklog(msg, val.Address)
		c.processBacklog()

		timeout := time.NewTimer(2 * time.Second)
		//vote should not be processed at propose step
		select {
		case ev := <-evChan:
			t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev))
		case <-timeout.C:
		}
		c.setStep(prevote)
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
			t.Error("unexpected timeout occurs")
		}
	})

	t.Run("same height, but old round", func(t *testing.T) {
		nilRoundVote := &types.Vote{
			Round:  0,
			Height: big.NewInt(1),
		}

		nilRoundVotePayload, err := types.Encode(nilRoundVote)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &types.ConsensusMessage{
			Code:       types.msgPrevote,
			Msg:        nilRoundVotePayload,
			decodedMsg: nilRoundVote,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		committeeSet := newTestCommitteeSet(1)
		val, _ := committeeSet.GetByIndex(0)

		expected := backlogEvent{
			msg: msg,
		}

		evChan := make(chan interface{}, 1)

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Post(expected).Do(func(ev interface{}) {
			evChan <- ev
		})

		c := &core{
			logger:   log.New("backend", "test", "id", 0),
			backend:  backendMock,
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]*types.ConsensusMessage),
			round:    1,
			step:     prevote,
			height:   big.NewInt(1),
		}

		c.lastHeader = &types.Header{Committee: committeeSet.Committee()}

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
			t.Error("unexpected timeout occurs")
		}
	})

	t.Run("future height message are not processed", func(t *testing.T) {
		nilRoundVote := &types.Vote{
			Round:  2,
			Height: big.NewInt(4),
		}

		nilRoundVotePayload, err := types.Encode(nilRoundVote)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &types.ConsensusMessage{
			Code:       types.msgPrevote,
			Msg:        nilRoundVotePayload,
			decodedMsg: nilRoundVote,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(0)

		committeeSet := newTestCommitteeSet(2)
		val, _ := committeeSet.GetByIndex(0)

		c := &core{
			logger:   log.New("backend", "test", "id", 0),
			backend:  backendMock,
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]*types.ConsensusMessage),
			round:    2,
			height:   big.NewInt(3),
		}

		c.lastHeader = &types.Header{Committee: committeeSet.Committee()}

		c.storeBacklog(msg, val.Address)
		c.processBacklog()
	})

	t.Run("future height message are processed when height change", func(t *testing.T) {
		nilRoundVote := &types.Vote{
			Round:  2,
			Height: big.NewInt(4),
		}

		nilRoundVotePayload, err := types.Encode(nilRoundVote)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &types.ConsensusMessage{
			Code:       types.msgPrevote,
			Msg:        nilRoundVotePayload,
			decodedMsg: nilRoundVote,
		}
		msg2 := &types.ConsensusMessage{
			Code:       types.msgPrecommit,
			Msg:        nilRoundVotePayload,
			decodedMsg: nilRoundVote,
		}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(0)

		committeeSet := newTestCommitteeSet(2)
		val, _ := committeeSet.GetByIndex(0)

		c := &core{
			logger:   log.New("backend", "test", "id", 0),
			backend:  backendMock,
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]*types.ConsensusMessage),
			round:    2,
			height:   big.NewInt(3),
		}

		c.lastHeader = &types.Header{Committee: committeeSet.Committee()}

		c.storeBacklog(msg, val.Address)
		c.storeBacklog(msg2, val.Address)
		c.setStep(prevote)
		c.processBacklog()
		c.setHeight(big.NewInt(4))

		backendMock.EXPECT().Post(gomock.Any()).Times(2)
		c.setStep(prevote)
		c.processBacklog()
		timeout := time.NewTimer(2 * time.Second)
		<-timeout.C
	})

	t.Run("untrusted messages are processed when height change", func(t *testing.T) {
		nilRoundVote := &types.Vote{
			Round:  2,
			Height: big.NewInt(4),
		}

		nilRoundVotePayload, err := types.Encode(nilRoundVote)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &types.ConsensusMessage{
			Code:       types.msgPrevote,
			Msg:        nilRoundVotePayload,
			decodedMsg: nilRoundVote,
		}
		msg2 := &types.ConsensusMessage{
			Code:       types.msgPrecommit,
			Msg:        nilRoundVotePayload,
			decodedMsg: nilRoundVote,
		}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)

		committeeSet := newTestCommitteeSet(2)

		c := &core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			address:          common.HexToAddress("0x1234567890"),
			backlogs:         make(map[common.Address][]*types.ConsensusMessage),
			backlogUnchecked: map[uint64][]*types.ConsensusMessage{},
			round:            2,
			height:           big.NewInt(3),
		}

		c.lastHeader = &types.Header{Committee: committeeSet.Committee()}

		backendMock.EXPECT().Post(gomock.Any()).Times(0)
		c.storeUncheckedBacklog(msg)
		c.storeUncheckedBacklog(msg2)
		c.setStep(prevote)
		c.processBacklog()
		c.setHeight(big.NewInt(4))

		backendMock.EXPECT().Post(gomock.Any()).Times(2)
		c.setStep(prevote)

		backendMock.EXPECT().Post(gomock.Any()).Times(0)
		c.processBacklog()
		<-time.NewTimer(2 * time.Second).C
	})

	t.Run("future round message are processed when round change", func(t *testing.T) {
		nilRoundVote := &types.Vote{
			Round:  2,
			Height: big.NewInt(4),
		}

		nilRoundVotePayload, err := types.Encode(nilRoundVote)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		msg := &types.ConsensusMessage{
			Code:       types.msgPrevote,
			Msg:        nilRoundVotePayload,
			decodedMsg: nilRoundVote,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(0)

		committeeSet := newTestCommitteeSet(2)
		val, err := committeeSet.GetByIndex(0)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		c := &core{
			logger:   log.New("backend", "test", "id", 0),
			backend:  backendMock,
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]*types.ConsensusMessage),
			step:     prevote,
			round:    1,
			height:   big.NewInt(4),
		}

		c.lastHeader = &types.Header{Committee: committeeSet.Committee()}

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

		backendMock := NewMockBackend(ctrl)

		c := &core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			address:          common.HexToAddress("0x1234567890"),
			backlogs:         make(map[common.Address][]*types.ConsensusMessage),
			backlogUnchecked: make(map[uint64][]*types.ConsensusMessage),
			step:             prevote,
			round:            1,
			height:           big.NewInt(4),
		}
		var messages []*types.ConsensusMessage

		for i := int64(0); i < MaxSizeBacklogUnchecked; i++ {
			nilRoundVote := &types.Vote{
				Round:  i % 10,
				Height: big.NewInt(i / (1 + i%10)),
			}
			payload, err := rlp.EncodeToBytes(nilRoundVote)
			require.NoError(t, err)
			msg := &types.ConsensusMessage{
				Code:       types.msgPrevote,
				Msg:        payload,
				decodedMsg: nilRoundVote,
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

		backendMock := NewMockBackend(ctrl)

		c := &core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			address:          common.HexToAddress("0x1234567890"),
			backlogs:         make(map[common.Address][]*types.ConsensusMessage),
			backlogUnchecked: make(map[uint64][]*types.ConsensusMessage),
			step:             prevote,
			round:            1,
			height:           big.NewInt(4),
		}

		var messages []*types.ConsensusMessage
		uncheckedFounds := make(map[uint64]struct{})
		backendMock.EXPECT().RemoveMessageFromLocalCache(gomock.Any()).Times(MaxSizeBacklogUnchecked).Do(func(payload []byte) {
			var msg types.ConsensusMessage
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
			nilRoundVote := &types.Vote{
				Round:  i % 10,
				Height: big.NewInt(i),
			}
			payload, err := rlp.EncodeToBytes(nilRoundVote)
			require.NoError(t, err)
			msg := &types.ConsensusMessage{
				Code:       types.msgPrevote,
				Msg:        payload,
				decodedMsg: nilRoundVote,
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
