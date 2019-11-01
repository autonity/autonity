package core

import (
	"os"
	"testing"

	"math/big"

	"time"

	"fmt"

	"reflect"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/rlp"
	"github.com/davecgh/go-spew/spew"
)

func TestWalUpdateHeight(t *testing.T) {
	dir := os.TempDir()
	eventMux := event.NewTypeMuxSilent(logger)
	sub := eventMux.Subscribe(events.MessageEvent{}, backlogEvent{}, events.CommitEvent{})

	wal := NewWal(log.Root(), dir, sub)
	defer os.RemoveAll(dir)
	defer wal.Close()

	wal.Start()
	wal.UpdateHeight(big.NewInt(1))
	time.Sleep(time.Millisecond)
	h, err := wal.Height()
	if err != nil {
		t.Fatal(err)
	}

	if check := h.Cmp(big.NewInt(1)); check != 0 {
		t.Fatal("Not equal", check)
	}

}

func TestWalAdd(t *testing.T) {
	dir := os.TempDir()
	eventMux := event.NewTypeMuxSilent(logger)
	sub := eventMux.Subscribe(events.MessageEvent{}, backlogEvent{}, events.CommitEvent{})
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlError, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	wal := NewWal(log.Root(), dir, sub)

	defer os.RemoveAll(dir)
	defer wal.Close()

	wal.Start()
	wal.UpdateHeight(big.NewInt(1))

	proposal := Proposal{
		Round:           big.NewInt(1),
		Height:          big.NewInt(1),
		ValidRound:      big.NewInt(1),
		IsValidRoundNil: big.NewInt(0),
		ProposalBlock:   types.NewBlockWithHeader(&types.Header{}),
	}

	encodedProposal, err := rlp.EncodeToBytes(&proposal)
	if err != nil {
		t.Fatal(err)
	}

	var proposalMsg = Message{
		Code:          msgProposal,
		Msg:           encodedProposal,
		Address:       common.Address{1},
		CommittedSeal: []byte{},
	}

	encodedProposalMsg, err := rlp.EncodeToBytes(&proposalMsg)
	if err != nil {
		t.Fatal(err)
	}

	eventMux.Post(events.MessageEvent{
		encodedProposalMsg,
	})

	var preVote = Vote{
		Round:  big.NewInt(1),
		Height: big.NewInt(1),
	}

	encodedVote, err := rlp.EncodeToBytes(&preVote)
	if err != nil {
		t.Fatal(err)
	}

	prevotes := []Message{}
	for i := 0; i < 5; i++ {
		msg := Message{
			Code:          msgPrevote,
			Msg:           encodedVote,
			Address:       common.Address{uint8(i)},
			CommittedSeal: []byte{},
		}
		prevotes = append(prevotes, msg)
		encMsg, err := rlp.EncodeToBytes(&msg)
		if err != nil {
			t.Fatal(err)
		}
		eventMux.Post(events.MessageEvent{
			encMsg,
		})

	}

	var precommit = Vote{
		Round:  big.NewInt(1),
		Height: big.NewInt(1),
	}

	encodedVote, err = rlp.EncodeToBytes(&precommit)
	if err != nil {
		t.Fatal(err)
	}

	precommits := []Message{}
	for i := 0; i < 5; i++ {
		msg := Message{
			Code:          msgPrecommit,
			Msg:           encodedVote,
			Address:       common.Address{uint8(i)},
			CommittedSeal: []byte{},
		}
		precommits = append(precommits, msg)
		encMsg, err := rlp.EncodeToBytes(&msg)
		if err != nil {
			t.Fatal(err)
		}
		eventMux.Post(events.MessageEvent{
			encMsg,
		})

	}

	time.Sleep(time.Second)
	allMsgs := make([]Message, 0, 1+len(prevotes)+len(precommits))
	allMsgs = append(allMsgs, proposalMsg)
	allMsgs = append(allMsgs, prevotes...)
	allMsgs = append(allMsgs, precommits...)
	if reflect.DeepEqual(wal.Get(big.NewInt(1))) {
		t.FailNow()
	}
}

func TestWalDebugOrig(t *testing.T) {
	dir := os.TempDir()
	eventMux := event.NewTypeMuxSilent(logger)
	sub := eventMux.Subscribe(events.MessageEvent{}, backlogEvent{}, events.CommitEvent{})
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlError, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	wal := NewWal(log.Root(), dir, sub)

	defer os.RemoveAll(dir)
	defer wal.Close()

	wal.Start()
	wal.UpdateHeight(big.NewInt(1))

	proposal := Proposal{
		Round:           big.NewInt(1),
		Height:          big.NewInt(1),
		ValidRound:      big.NewInt(1),
		IsValidRoundNil: big.NewInt(0),
		ProposalBlock:   &types.Block{},
	}

	encodedProposal, err := rlp.EncodeToBytes(&proposal)
	if err != nil {
		t.Fatal(err)
	}

	var proposalMsg = &Message{
		Code:          msgProposal,
		Msg:           encodedProposal,
		Address:       common.Address{1},
		CommittedSeal: []byte{},
	}

	encodedProposalMsg, err := rlp.EncodeToBytes(&proposalMsg)
	if err != nil {
		t.Fatal(err)
	}

	eventMux.Post(events.MessageEvent{
		encodedProposalMsg,
	})
	time.Sleep(time.Second)
	fmt.Println("--------after propose---------")
	spew.Dump(wal.Get(big.NewInt(1)))
	fmt.Println("-----------------")

	var preVote = Vote{
		Round:  big.NewInt(1),
		Height: big.NewInt(1),
	}

	encodedVote, err := rlp.EncodeToBytes(&preVote)
	if err != nil {
		t.Fatal(err)
	}

	prevotes := []Message{}
	for i := 0; i < 5; i++ {
		msg := Message{
			Code:          msgPrecommit,
			Msg:           encodedVote,
			Address:       common.Address{uint8(i)},
			CommittedSeal: []byte{},
		}
		prevotes = append(prevotes, msg)
		encMsg, err := rlp.EncodeToBytes(&msg)
		if err != nil {
			t.Fatal(err)
		}
		eventMux.Post(events.MessageEvent{
			encMsg,
		})

	}

	time.Sleep(time.Second)
	fmt.Println("--------after prevote---------")
	spew.Dump(wal.Get(big.NewInt(1)))
	fmt.Println("-----------------")

	var precommit = Vote{
		Round:  big.NewInt(1),
		Height: big.NewInt(1),
	}

	encodedVote, err = rlp.EncodeToBytes(&precommit)
	if err != nil {
		t.Fatal(err)
	}

	precommits := []Message{}
	for i := 0; i < 5; i++ {
		msg := Message{
			Code:          msgPrecommit,
			Msg:           encodedVote,
			Address:       common.Address{uint8(i)},
			CommittedSeal: []byte{},
		}
		precommits = append(precommits, msg)
		encMsg, err := rlp.EncodeToBytes(&msg)
		if err != nil {
			t.Fatal(err)
		}
		eventMux.Post(events.MessageEvent{
			encMsg,
		})

	}

	time.Sleep(time.Second)

	fmt.Println("--------after precommit---------")
	spew.Dump(wal.Get(big.NewInt(1)))
	fmt.Println("-----------------")
}
