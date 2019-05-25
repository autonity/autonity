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
	"crypto/ecdsa"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
)

func TestSign(t *testing.T) {
	b := newBackend(t)
	data := []byte("Here is a string....")
	sig, err := b.Sign(data)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	//Check signature recover
	hashData := crypto.Keccak256(data)
	pubkey, _ := crypto.Ecrecover(hashData, sig)
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])
	if signer != getAddress() {
		t.Errorf("address mismatch: have %v, want %s", signer.Hex(), getAddress().Hex())
	}
}

func TestCheckSignature(t *testing.T) {
	key, _ := generatePrivateKey()
	data := []byte("Here is a string....")
	hashData := crypto.Keccak256(data)
	sig, _ := crypto.Sign(hashData, key)
	b := newBackend(t)
	a := getAddress()
	err := b.CheckSignature(data, a, sig)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	a = getInvalidAddress()
	err = b.CheckSignature(data, a, sig)
	if err != types.ErrInvalidSignature {
		t.Errorf("error mismatch: have %v, want %v", err, types.ErrInvalidSignature)
	}
}

func TestCheckValidatorSignature(t *testing.T) {
	vset, keys := newTestValidatorSet(5)

	// 1. Positive test: sign with validator's key should succeed
	data := []byte("dummy data")
	hashData := crypto.Keccak256(data)
	for i, k := range keys {
		// Sign
		sig, err := crypto.Sign(hashData, k)
		if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
		// CheckValidatorSignature should succeed
		addr, err := tendermint.CheckValidatorSignature(vset, data, sig)
		if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
		val := vset.GetByIndex(uint64(i))
		if addr != val.Address() {
			t.Errorf("val address mismatch: have %v, want %v", addr, val.Address())
		}
	}

	// 2. Negative test: sign with any key other than validator's key should return error
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	// Sign
	sig, err := crypto.Sign(hashData, key)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	// CheckValidatorSignature should return ErrUnauthorizedAddress
	addr, err := tendermint.CheckValidatorSignature(vset, data, sig)
	if err != tendermint.ErrUnauthorizedAddress {
		t.Errorf("error mismatch: have %v, want %v", err, tendermint.ErrUnauthorizedAddress)
	}
	emptyAddr := common.Address{}
	if addr != emptyAddr {
		t.Errorf("address mismatch: have %v, want %v", addr, emptyAddr)
	}
}

func TestCommit(t *testing.T) {
	backend := newBackend(t)

	commitCh := make(chan *types.Block, 1)
	backend.commitCh = commitCh

	// Case: it's a proposer, so the Backend.commit will receive channel result from Backend.Commit function
	testCases := []struct {
		expectedErr       error
		expectedSignature [][]byte
		expectedBlock     func() *types.Block
	}{
		{
			// normal case
			nil,
			[][]byte{append([]byte{1}, bytes.Repeat([]byte{0x00}, types.PoSExtraSeal-1)...)},
			func() *types.Block {
				chain, engine := newBlockChain(1, t)
				block, err := makeBlockWithoutSeal(chain, engine, chain.Genesis())
				if err != nil {
					t.Fatal(err)
				}
				expectedBlock, _ := engine.updateBlock(block)
				return expectedBlock
			},
		},
		{
			// invalid signature
			types.ErrInvalidCommittedSeals,
			nil,
			func() *types.Block {
				chain, engine := newBlockChain(1, t)
				block, err := makeBlockWithoutSeal(chain, engine, chain.Genesis())
				if err != nil {
					t.Fatal(err)
				}
				expectedBlock, _ := engine.updateBlock(block)
				return expectedBlock
			},
		},
	}

	for _, test := range testCases {
		expBlock := *test.expectedBlock()

		backend.proposedBlockHash = expBlock.Hash()
		if err := backend.Commit(expBlock, test.expectedSignature); err != nil {
			if err != test.expectedErr {
				t.Errorf("error mismatch: have %v, want %v", err, test.expectedErr)
			}
		}

		if test.expectedErr == nil {
			// to avoid race condition is occurred by goroutine
			select {
			case result := <-commitCh:
				if result.Hash() != expBlock.Hash() {
					t.Errorf("hash mismatch: have %v, want %v", result.Hash(), expBlock.Hash())
				}
			case <-time.After(10 * time.Second):
				t.Fatal("timeout")
			}
		}
	}
}

func TestGetProposer(t *testing.T) {
	chain, engine := newBlockChain(1, t)
	block, err := makeBlock(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}

	insertedCount, err := chain.InsertChain(types.Blocks{block})
	if err != nil {
		t.Errorf("can't insert the block into the chain. Block %v. Chain genesis %v", block, chain.Genesis())
	}
	if insertedCount != 1 {
		t.Errorf("expected to insert only 1 block got %d", insertedCount)
	}

	expected := engine.GetProposer(1)
	actual := engine.Address()
	if actual != expected {
		t.Errorf("proposer mismatch: have %v, want %v", actual.Hex(), expected.Hex())
	}
}

/**
 * SimpleBackend
 * Private key: bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1
 * Public key: 04a2bfb0f7da9e1b9c0c64e14f87e8fb82eb0144e97c25fe3a977a921041a50976984d18257d2495e7bfd3d4b280220217f429287d25ecdf2b0d7c0f7aae9aa624
 * Address: 0x70524d664ffe731100208a0154e556f9bb679ae6
 */
func getAddress() common.Address {
	return common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
}

func getInvalidAddress() common.Address {
	return common.HexToAddress("0x9535b2e7faaba5288511d89341d94a38063a349b")
}

func generatePrivateKey() (*ecdsa.PrivateKey, error) {
	key := "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
	return crypto.HexToECDSA(key)
}

func newTestValidatorSet(n int) (tendermint.ValidatorSet, []*ecdsa.PrivateKey) {
	// generate validators
	keys := make(Keys, n)
	addrs := make([]common.Address, n)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		keys[i] = privateKey
		addrs[i] = crypto.PubkeyToAddress(privateKey.PublicKey)
	}
	vset := validator.NewSet(addrs, tendermint.RoundRobin)
	sort.Sort(keys) //Keys need to be sorted by its public key address
	return vset, keys
}

type Keys []*ecdsa.PrivateKey

func (slice Keys) Len() int {
	return len(slice)
}

func (slice Keys) Less(i, j int) bool {
	return strings.Compare(crypto.PubkeyToAddress(slice[i].PublicKey).String(), crypto.PubkeyToAddress(slice[j].PublicKey).String()) < 0
}

func (slice Keys) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func newBackend(t *testing.T) (b *Backend) {
	_, b = newBlockChain(4, t)
	key, err := generatePrivateKey()
	if err != nil {
		t.Errorf("failed to generate private key: %v", err)
	}

	b.privateKey = key
	return
}
