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

package core

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
)

func testProposal(t *testing.T) {
	pp := &tendermint.Proposal{
		View: &tendermint.View{
			Round:    big.NewInt(1),
			Sequence: big.NewInt(2),
		},
		ProposalBlock: makeBlock(1),
	}
	proposalPayload, _ := Encode(pp)

	m := &message{
		Code:    msgProposal,
		Msg:     proposalPayload,
		Address: common.HexToAddress("0x1234567890"),
	}

	msgPayload, err := m.Payload()
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	decodedMsg := new(message)
	err = decodedMsg.FromPayload(msgPayload, nil)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	var decodedPP *tendermint.Proposal
	err = decodedMsg.Decode(&decodedPP)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	// if block is encoded/decoded by rlp, we cannot to compare interface data type using reflect.DeepEqual. (like tendermint.ProposalBlock)
	// so individual comparison here.
	if !reflect.DeepEqual(pp.ProposalBlock.Hash(), decodedPP.ProposalBlock.Hash()) {
		t.Errorf("proposal hash mismatch: have %v, want %v", decodedPP.ProposalBlock.Hash(), pp.ProposalBlock.Hash())
	}

	if !reflect.DeepEqual(pp.View, decodedPP.View) {
		t.Errorf("view mismatch: have %v, want %v", decodedPP.View, pp.View)
	}

	if !reflect.DeepEqual(pp.ProposalBlock.Number(), decodedPP.ProposalBlock.Number()) {
		t.Errorf("proposal number mismatch: have %v, want %v", decodedPP.ProposalBlock.Number(), pp.ProposalBlock.Number())
	}
}

func testSubject(t *testing.T) {
	s := &tendermint.Subject{
		View: &tendermint.View{
			Round:    big.NewInt(1),
			Sequence: big.NewInt(2),
		},
		Digest: common.BytesToHash([]byte("1234567890")),
	}

	subjectPayload, _ := Encode(s)

	m := &message{
		Code:    msgProposal,
		Msg:     subjectPayload,
		Address: common.HexToAddress("0x1234567890"),
	}

	msgPayload, err := m.Payload()
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	decodedMsg := new(message)
	err = decodedMsg.FromPayload(msgPayload, nil)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	var decodedSub *tendermint.Subject
	err = decodedMsg.Decode(&decodedSub)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	if !reflect.DeepEqual(s, decodedSub) {
		t.Errorf("subject mismatch: have %v, want %v", decodedSub, s)
	}
}

func testSubjectWithSignature(t *testing.T) {
	s := &tendermint.Subject{
		View: &tendermint.View{
			Round:    big.NewInt(1),
			Sequence: big.NewInt(2),
		},
		Digest: common.BytesToHash([]byte("1234567890")),
	}
	expectedSig := []byte{0x01}

	subjectPayload, _ := Encode(s)
	// 1. Encode test
	m := &message{
		Code:          msgProposal,
		Msg:           subjectPayload,
		Address:       common.HexToAddress("0x1234567890"),
		Signature:     expectedSig,
		CommittedSeal: []byte{},
	}

	msgPayload, err := m.Payload()
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	// 2. Decode test
	// 2.1 Test normal validate func
	decodedMsg := new(message)
	err = decodedMsg.FromPayload(msgPayload, func(data []byte, sig []byte) (common.Address, error) {
		return common.Address{}, nil
	})
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	if !reflect.DeepEqual(decodedMsg, m) {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	// 2.2 Test nil validate func
	decodedMsg = new(message)
	err = decodedMsg.FromPayload(msgPayload, nil)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(decodedMsg, m) {
		t.Errorf("message mismatch: have %v, want %v", decodedMsg, m)
	}

	// 2.3 Test failed validate func
	decodedMsg = new(message)
	err = decodedMsg.FromPayload(msgPayload, func(data []byte, sig []byte) (common.Address, error) {
		return common.Address{}, tendermint.ErrUnauthorizedAddress
	})
	if err != tendermint.ErrUnauthorizedAddress {
		t.Errorf("error mismatch: have %v, want %v", err, tendermint.ErrUnauthorizedAddress)
	}
}

func TestMessageEncodeDecode(t *testing.T) {
	testProposal(t)
	testSubject(t)
	testSubjectWithSignature(t)
}
