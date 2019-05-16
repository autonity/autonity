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
	"testing"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
)

// notice: the normal case have been tested in integration tests.
func TestHandleMsg(t *testing.T) {
	N := uint64(4)
	F := uint64(1)
	sys := newTestSystemWithBackend(N, F)

	closer := sys.Run(true)
	defer closer()

	v0 := sys.backends[0]
	r0 := v0.engine.(*core)

	m, _ := Encode(&tendermint.Subject{
		View: &tendermint.View{
			Sequence: big.NewInt(0),
			Round:    big.NewInt(0),
		},
		Digest: common.BytesToHash([]byte("1234567890")),
	})
	// with a matched payload. msgProposal should match with *tendermint.ProposalBlock in normal case.
	msg := &message{
		Code:          msgProposal,
		Msg:           m,
		Address:       v0.Address(),
		Signature:     []byte{},
		CommittedSeal: []byte{},
	}

	_, val := v0.Validators(0).GetByAddress(v0.Address())
	if err := r0.handleCheckedMsg(msg, val); err != errFailedDecodeProposal {
		t.Errorf("error mismatch: have %v, want %v", err, errFailedDecodeProposal)
	}

	m, _ = Encode(&tendermint.Proposal{
		View: &tendermint.View{
			Sequence: big.NewInt(0),
			Round:    big.NewInt(0),
		},
		ProposalBlock: makeBlock(1),
	})
	// with a unmatched payload. msgPrevote should match with *tendermint.Subject in normal case.
	msg = &message{
		Code:          msgPrevote,
		Msg:           m,
		Address:       v0.Address(),
		Signature:     []byte{},
		CommittedSeal: []byte{},
	}

	_, val = v0.Validators(0).GetByAddress(v0.Address())
	if err := r0.handleCheckedMsg(msg, val); err != errFailedDecodePrevote {
		t.Errorf("error mismatch: have %v, want %v", err, errFailedDecodeProposal)
	}

	m, _ = Encode(&tendermint.Proposal{
		View: &tendermint.View{
			Sequence: big.NewInt(0),
			Round:    big.NewInt(0),
		},
		ProposalBlock: makeBlock(2),
	})
	// with a unmatched payload. tendermint.MsgPrecommit should match with *tendermint.Subject in normal case.
	msg = &message{
		Code:          msgPrecommit,
		Msg:           m,
		Address:       v0.Address(),
		Signature:     []byte{},
		CommittedSeal: []byte{},
	}

	_, val = v0.Validators(0).GetByAddress(v0.Address())
	if err := r0.handleCheckedMsg(msg, val); err != errFailedDecodePrecommit {
		t.Errorf("error mismatch: have %v, want %v", err, errFailedDecodePrecommit)
	}

	m, _ = Encode(&tendermint.Proposal{
		View: &tendermint.View{
			Sequence: big.NewInt(0),
			Round:    big.NewInt(0),
		},
		ProposalBlock: makeBlock(3),
	})
	// invalid message code. message code is not exists in list
	msg = &message{
		Code:          uint64(99),
		Msg:           m,
		Address:       v0.Address(),
		Signature:     []byte{},
		CommittedSeal: []byte{},
	}

	_, val = v0.Validators(0).GetByAddress(v0.Address())
	if err := r0.handleCheckedMsg(msg, val); err == nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	// with malicious payload
	if err := r0.handleMsg([]byte{1}); err == nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
}
