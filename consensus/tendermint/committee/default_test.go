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

package committee

import (
	"github.com/clearmatics/autonity/core/types"
	"math/big"
	"strings"
	"testing"

	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/crypto"
)

var (
	testAddress  = "70524d664ffe731100208a0154e556f9bb679ae6"
	testAddress2 = "b37866a925bccd69cfa98d43b510f1d23d78a851"
)

func TestValidatorSet(t *testing.T) {
	testNewValidatorSet(t)
	//testNormalValSet(t)
	//testEmptyValSet(t)
	//testStickyProposer(t)
}

func testNewValidatorSet(t *testing.T) {
	var validators types.Committee
	const ValCnt = 100

	// Create 100 members with random addresses
	for i := 0; i < ValCnt; i++ {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)
		val := types.CommitteeMember{addr, new(big.Int).SetUint64(1)}
		validators = append(validators, val)
	}

	// Create Set
	valSet := NewSet(validators, config.RoundRobin, validators[0].Address)
	if valSet == nil {
		t.Error("the validator byte array cannot be parsed")
		t.FailNow()
	}

	if valSet.Size() != ValCnt {
		t.Errorf("validator set has %d elements instead of %d", valSet.Size(), ValCnt)
	}

	valsMap := make(map[string]struct{})
	for _, val := range validators {
		valsMap[val.String()] = struct{}{}
	}

	// Check members sorting: should be in ascending order
	for i := 0; i < ValCnt-1; i++ {
		val, err := valSet.GetByIndex(i)
		if err != nil {
			t.Error("unexpected error")
		}
		nextVal, err := valSet.GetByIndex(i + 1)
		if err != nil {
			t.Error("unexpected error")
		}
		if strings.Compare(val.String(), nextVal.String()) >= 0 {
			t.Error("validator set is not sorted in ascending order")
		}

		if _, ok := valsMap[val.String()]; !ok {
			t.Errorf("validator set has unexpected element %s. Original members %v, given %v",
				val.String(), validators, valSet.Committee())
		}
	}
}

/*
func testNormalValSet(t *testing.T) {
	b1 := common.Hex2Bytes(testAddress)
	b2 := common.Hex2Bytes(testAddress2)
	addr1 := common.BytesToAddress(b1)
	addr2 := common.BytesToAddress(b2)
	val1 := types.CommitteeMember{addr1, new(big.Int).SetUint64(1)}
	val2 := types.CommitteeMember{addr2, new(big.Int).SetUint64(1)}

	valSet := NewSet(types.Committee{val1, val2}, config.RoundRobin, val1.Address)
	if valSet == nil {
		t.Errorf("the format of validator set is invalid")
		t.FailNow()
	}

	// check size
	if size := valSet.Size(); size != 2 {
		t.Errorf("the size of validator set is wrong: have %v, want 2", size)
	}
	// test get by index
	if val, err := valSet.GetByIndex(0); err != nil || !reflect.DeepEqual(val, val1) {
		t.Errorf("validator mismatch: have %v, want %v", val, val1)
	}
	// test get by invalid index
	if val, err := valSet.GetByIndex(2); err != nil {
		t.Errorf("validator mismatch: have %s, want nil", err)
	}
	// test get by address
	if _, val, err := valSet.GetByAddress(addr2); err != nil ||  !reflect.DeepEqual(val, val2) {
		t.Errorf("validator mismatch: have %v, want %v", val, val2)
	}
	// test get by invalid address
	invalidAddr := common.HexToAddress("0x9535b2e7faaba5288511d89341d94a38063a349b")
	if _, val := valSet.GetByAddress(invalidAddr); val != nil {
		t.Errorf("validator mismatch: have %v, want nil", val)
	}
	// test get proposer
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val1)
	}
	// test calculate proposer
	lastProposer := addr1
	valSet.CalcProposer(lastProposer, uint64(0))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val2) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val2)
	}
	valSet.CalcProposer(lastProposer, uint64(3))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val1)
	}
	// test empty last proposer
	lastProposer = common.Address{}
	valSet.CalcProposer(lastProposer, uint64(3))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val2) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val2)
	}
}

func testEmptyValSet(t *testing.T) {
	valSet := newDefaultSet(Members{}, config.RoundRobin)
	if valSet == nil {
		t.Errorf("validator set should not be nil")
	}
}

func testStickyProposer(t *testing.T) {
	b1 := common.Hex2Bytes(testAddress)
	b2 := common.Hex2Bytes(testAddress2)
	addr1 := common.BytesToAddress(b1)
	addr2 := common.BytesToAddress(b2)
	val1 := New(addr1, new(big.Int).SetUint64(1))
	val2 := New(addr2, new(big.Int).SetUint64(1))

	valSet := newDefaultSet(Members{val1, val2}, config.Sticky)

	// test get proposer
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val1)
	}
	// test calculate proposer
	lastProposer := addr1
	valSet.CalcProposer(lastProposer, uint64(0))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val1)
	}

	valSet.CalcProposer(lastProposer, uint64(1))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val2) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val2)
	}

	// test empty last proposer
	lastProposer = common.Address{}
	valSet.CalcProposer(lastProposer, uint64(3))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val2) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val2)
	}
}
*/
