package test

import (
	"io/ioutil"
	"os"
	"testing"

	"math/big"

	"time"

	"reflect"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/rlp"
)

func TestWalUpdateHeight(t *testing.T) {
	dir, err := ioutil.TempDir("", "wal_")
	if err != nil {
		t.Fatal(err)
	}

	eventMux := event.NewTypeMuxSilent(log.Root())
	sub := eventMux.Subscribe(events.MessageEvent{}, events.CommitEvent{})

	wal := core.NewWal(log.Root(), dir, sub)
	defer os.RemoveAll(dir) //nolint
	defer wal.Close()

	wal.Start()
	err = wal.UpdateHeight(big.NewInt(1))
	if err != nil {
		t.Fatal(err)
	}
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
	dir := os.TempDir() + "/wal_test"
	eventMux := event.NewTypeMuxSilent(log.Root())
	sub := eventMux.Subscribe(events.MessageEvent{}, events.CommitEvent{})
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlError, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	wal := core.NewWal(log.Root(), dir, sub)

	defer os.RemoveAll(dir)
	defer wal.Close()

	wal.Start()
	err := wal.UpdateHeight(big.NewInt(1))
	if err != nil {
		t.Fatal(err)
	}

	proposal := core.Proposal{
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

	var proposalMsg = core.Message{
		Code:          0,
		Msg:           encodedProposal,
		Address:       common.Address{1},
		CommittedSeal: []byte{},
	}

	encodedProposalMsg, err := rlp.EncodeToBytes(&proposalMsg)
	if err != nil {
		t.Fatal(err)
	}

	eventMux.Post(events.MessageEvent{
		Payload: encodedProposalMsg,
	})

	var preVote = core.Vote{
		Round:  big.NewInt(1),
		Height: big.NewInt(1),
	}

	encodedVote, err := rlp.EncodeToBytes(&preVote)
	if err != nil {
		t.Fatal(err)
	}

	prevotes := []core.Message{}
	for i := 0; i < 5; i++ {
		msg := core.Message{
			Code:          1,
			Msg:           encodedVote,
			Address:       common.Address{uint8(i)},
			CommittedSeal: []byte{},
		}
		prevotes = append(prevotes, msg)
		encMsg, innerErr := rlp.EncodeToBytes(&msg)
		if innerErr != nil {
			t.Fatal(innerErr)
		}
		eventMux.Post(events.MessageEvent{
			Payload: encMsg,
		})

	}

	var precommit = core.Vote{
		Round:  big.NewInt(1),
		Height: big.NewInt(1),
	}

	encodedVote, err = rlp.EncodeToBytes(&precommit)
	if err != nil {
		t.Fatal(err)
	}

	precommits := []core.Message{}
	for i := 0; i < 5; i++ {
		msg := core.Message{
			Code:          2,
			Msg:           encodedVote,
			Address:       common.Address{uint8(i)},
			CommittedSeal: []byte{},
		}
		precommits = append(precommits, msg)
		encMsg, innerErr := rlp.EncodeToBytes(&msg)
		if innerErr != nil {
			t.Fatal(innerErr)
		}
		eventMux.Post(events.MessageEvent{
			Payload: encMsg,
		})

	}

	time.Sleep(time.Second)
	allMsgs := make([]core.Message, 0, 1+len(prevotes)+len(precommits))
	allMsgs = append(allMsgs, proposalMsg)
	allMsgs = append(allMsgs, prevotes...)
	allMsgs = append(allMsgs, precommits...)
	resp, err := wal.Get(big.NewInt(1))
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(resp, allMsgs) {
		t.FailNow()
	}
}
