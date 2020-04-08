// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package backend

import (
	"bytes"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/hexutil"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
)

func TestPrepare(t *testing.T) {
	chain, engine := newBlockChain(1)
	header := makeHeader(chain.Genesis(), engine.config)
	err := engine.Prepare(chain, header)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	header.ParentHash = common.BytesToHash([]byte("1234567890"))
	err = engine.Prepare(chain, header)
	if err != consensus.ErrUnknownAncestor {
		t.Errorf("error mismatch: have %v, want %v", err, consensus.ErrUnknownAncestor)
	}
}

func TestSealCommittedOtherHash(t *testing.T) {
	chain, engine := newBlockChain(4)

	block, err := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	otherBlock, err := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	eventSub := engine.Subscribe(events.CommitEvent{})
	eventLoop := func() {
		ev := <-eventSub.Chan()
		_, ok := ev.Data.(events.CommitEvent)
		if !ok {
			t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev.Data))
		}
		err = engine.Commit(otherBlock, 0, [][]byte{})
		if err != nil {
			t.Error("commit should not return error", err.Error())
		}

		eventSub.Unsubscribe()
	}
	go eventLoop()
	seal := func() {
		resultCh := make(chan *types.Block)
		err = engine.Seal(chain, block, resultCh, nil)
		if err != nil {
			t.Error("seal should not return error", err.Error())
		}

		<-resultCh
		t.Error("seal should not be completed")
	}
	go seal()

	const timeoutDura = 2 * time.Second
	timeout := time.NewTimer(timeoutDura)
	<-timeout.C
	// wait 2 seconds to ensure we cannot get any blocks from Istanbul
}

func TestSealCommitted(t *testing.T) {
	chain, engine := newBlockChain(1)
	block, err := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	expectedBlock, _ := engine.AddSeal(block)

	resultCh := make(chan *types.Block)
	err = engine.Seal(chain, block, resultCh, nil)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	finalBlock := <-resultCh
	if finalBlock.Hash() != expectedBlock.Hash() {
		t.Errorf("hash mismatch: have %v, want %v", finalBlock.Hash(), expectedBlock.Hash())
	}
}

func TestVerifyHeader(t *testing.T) {
	chain, engine := newBlockChain(1)

	// errEmptyCommittedSeals case
	block, err := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	block, _ = engine.AddSeal(block)
	err = engine.VerifyHeader(chain, block.Header(), false)
	if err != types.ErrEmptyCommittedSeals {
		t.Errorf("error mismatch: have %v, want %v", err, types.ErrEmptyCommittedSeals)
	}

	header := block.Header()

	// non zero MixDigest
	block, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	header = block.Header()
	header.MixDigest = common.BytesToHash([]byte("123456789"))
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidMixDigest {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidMixDigest)
	}

	// invalid uncles hash
	block, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	header = block.Header()
	header.UncleHash = common.BytesToHash([]byte("123456789"))
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidUncleHash {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidUncleHash)
	}

	// invalid difficulty
	block, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	header = block.Header()
	header.Difficulty = big.NewInt(2)
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidDifficulty {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidDifficulty)
	}

	// invalid timestamp
	block, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	header = block.Header()
	header.Time = new(big.Int).Add(big.NewInt(int64(chain.Genesis().Time())), new(big.Int).SetUint64(engine.config.BlockPeriod-1)).Uint64()
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidTimestamp {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidTimestamp)
	}

	// future block
	block, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	header = block.Header()
	header.Time = new(big.Int).Add(big.NewInt(now().Unix()), new(big.Int).SetUint64(10)).Uint64()
	err = engine.VerifyHeader(chain, header, false)
	if err != consensus.ErrFutureBlock {
		t.Errorf("error mismatch: have %v, want %v", err, consensus.ErrFutureBlock)
	}

	// invalid nonce
	block, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	header = block.Header()
	copy(header.Nonce[:], hexutil.MustDecode("0x111111111111"))
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidNonce {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidNonce)
	}
}

func TestVerifySeal(t *testing.T) {
	chain, engine := newBlockChain(1)
	genesis := chain.Genesis()
	// cannot verify genesis
	err := engine.VerifySeal(chain, genesis.Header())
	if err != errUnknownBlock {
		t.Errorf("error mismatch: have %v, want %v", err, errUnknownBlock)
	}

	block, err := makeBlock(chain, engine, genesis)
	if err != nil {
		t.Fatal(err)
	}

	// change block content
	header := block.Header()
	header.Number = big.NewInt(4)
	block1 := block.WithSeal(header)
	err = engine.VerifySeal(chain, block1.Header())
	if err != errUnknownBlock {
		t.Errorf("error mismatch: have %v, want %v", err, errUnknownBlock)
	}

	// unauthorized users but still can get correct signer address
	privateKey, _ := crypto.GenerateKey()
	engine.SetPrivateKey(privateKey)
	err = engine.VerifySeal(chain, block.Header())
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
}

/* The logic of this needs to change with respect of Autonity contact */
func TestVerifyHeaders(t *testing.T) {
	chain, engine := newBlockChain(1)

	// success case
	headers := []*types.Header{}
	blocks := []*types.Block{}
	size := 100

	var err error
	for i := 0; i < size; i++ {
		var b *types.Block
		if i == 0 {
			b, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
		} else {
			b, err = makeBlockWithoutSeal(chain, engine, blocks[i-1])
		}
		if err != nil {
			t.Fatal(err)
		}

		b, _ = engine.AddSeal(b)

		blocks = append(blocks, b)
		headers = append(headers, blocks[i].Header())
	}

	now = func() time.Time {
		return time.Unix(int64(headers[size-1].Time), 0)
	}

	_, results := engine.VerifyHeaders(chain, headers, nil)

	const timeoutDura = 2 * time.Second
	timeout := time.NewTimer(timeoutDura)
	index := 0
OUT1:
	for {
		select {
		case err := <-results:
			if err != nil {
				/*  The two following errors mean that the processing has gone right */
				if err != types.ErrEmptyCommittedSeals && err != types.ErrInvalidCommittedSeals {
					t.Errorf("error mismatch: have %v, want errEmptyCommittedSeals|errInvalidCommittedSeals", err)
					break OUT1
				}
			}
			index++
			if index == size {
				break OUT1
			}
		case <-timeout.C:
			break OUT1
		}
	}
}

/* The logic of this needs to change with respect of Autonity contact */
func TestVerifyHeadersAbortValidation(t *testing.T) {
	chain, engine := newBlockChain(1)

	// success case
	headers := []*types.Header{}
	blocks := []*types.Block{}
	size := 100

	var err error
	for i := 0; i < size; i++ {
		var b *types.Block
		if i == 0 {
			b, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
		} else {
			b, err = makeBlockWithoutSeal(chain, engine, blocks[i-1])
		}
		if err != nil {
			t.Fatal(err)
		}

		b, _ = engine.AddSeal(b)

		blocks = append(blocks, b)
		headers = append(headers, blocks[i].Header())
	}

	now = func() time.Time {
		return time.Unix(int64(headers[size-1].Time), 0)
	}

	const timeoutDura = 2 * time.Second

	// abort cases
	abort, results := engine.VerifyHeaders(chain, headers, nil)
	timeout := time.NewTimer(timeoutDura)
	index := 0
OUT2:
	for {
		select {
		case err := <-results:
			if err != nil {
				if err != types.ErrEmptyCommittedSeals && err != types.ErrInvalidCommittedSeals {
					t.Errorf("error mismatch: have %v, want errEmptyCommittedSeals|errInvalidCommittedSeals", err)
					break OUT2
				}
			}
			index++
			if index == 5 {
				abort <- struct{}{}
			}
			if index >= size {
				t.Errorf("verifyheaders should be aborted")
				break OUT2
			}
		case <-timeout.C:
			break OUT2
		}
	}
}

/* The logic of this needs to change with respect of Autonity contact */
func TestVerifyErrorHeaders(t *testing.T) {
	chain, engine := newBlockChain(1)

	// success case
	headers := []*types.Header{}
	blocks := []*types.Block{}
	size := 100

	var err error
	for i := 0; i < size; i++ {
		var b *types.Block
		if i == 0 {
			b, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
		} else {
			b, err = makeBlockWithoutSeal(chain, engine, blocks[i-1])
		}
		if err != nil {
			t.Fatal(err)
		}

		b, _ = engine.AddSeal(b)

		blocks = append(blocks, b)
		headers = append(headers, blocks[i].Header())
	}

	now = func() time.Time {
		return time.Unix(int64(headers[size-1].Time), 0)
	}

	const timeoutDura = 2 * time.Second

	// error header cases
	headers[2].Number = big.NewInt(100)
	_, results := engine.VerifyHeaders(chain, headers, nil)
	timeout := time.NewTimer(timeoutDura)
	index := 0
	errors := 0
	expectedErrors := 2

OUT3:
	for {
		select {
		case err := <-results:
			if err != nil {
				if err != types.ErrEmptyCommittedSeals && err != types.ErrInvalidCommittedSeals {
					errors++
				}
			}
			index++
			if index == size {
				if errors != expectedErrors {
					t.Errorf("error mismatch: have %v, want %v", err, expectedErrors)
				}
				break OUT3
			}
		case <-timeout.C:
			break OUT3
		}
	}
}

func TestWriteCommittedSeals(t *testing.T) {

	expectedCommittedSeal := append([]byte{1, 2, 3}, bytes.Repeat([]byte{0x00}, types.BFTExtraSeal-3)...)
	var expectedErr error

	h := &types.Header{}

	// normal case
	err := types.WriteCommittedSeals(h, [][]byte{expectedCommittedSeal})
	if err != expectedErr {
		t.Errorf("error mismatch: have %v, want %v", err, expectedErr)
	}

	if !reflect.DeepEqual(h.CommittedSeals, [][]byte{expectedCommittedSeal}) {
		t.Errorf("extra data mismatch: have %v, want %v", h.CommittedSeals, expectedCommittedSeal)
	}

	// invalid seal
	unexpectedCommittedSeal := append(expectedCommittedSeal, make([]byte, 1)...)
	err = types.WriteCommittedSeals(h, [][]byte{unexpectedCommittedSeal})
	if err != types.ErrInvalidCommittedSeals {
		t.Errorf("error mismatch: have %v, want %v", err, types.ErrInvalidCommittedSeals)
	}
}

func TestAPIs(t *testing.T) {
	b := &Backend{}

	APIS := b.APIs(nil)
	if len(APIS) < 1 {
		t.Fatalf("expected non empty slice")
	}

	if APIS[0].Namespace != "tendermint" {
		t.Fatalf("expected 'tendermint', got %v", APIS[0].Namespace)
	}
}

func TestClose(t *testing.T) {
	t.Run("engine is not running, error returned", func(t *testing.T) {
		b := &Backend{}

		err := b.Close()
		if err != ErrStoppedEngine {
			t.Fatalf("expected %v, got %v", ErrStoppedEngine, err)
		}
	})

	t.Run("engine is running, no errors", func(t *testing.T) {
		b := &Backend{
			coreStarted: true,
			stopped:     make(chan struct{}),
		}

		err := b.Close()
		if err != nil {
			t.Fatalf("expected <nil>, got %v", err)
		}
	})
}

func TestBackendSealHash(t *testing.T) {
	b := &Backend{}

	res := b.SealHash(&types.Header{})
	if res.Hex() == "" {
		t.Fatalf("expected not empty string")
	}
}
