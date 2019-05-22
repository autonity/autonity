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

package tendermint

import (
	"fmt"
	"io"
	"math/big"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/rlp"
)

type Proposal struct {
	Round         *big.Int
	Height        *big.Int
	ValidRound    *big.Int
	ProposalBlock types.Block
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (p *Proposal) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		p.Round,
		p.Height,
		p.ValidRound,
		p.ProposalBlock})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (p *Proposal) DecodeRLP(s *rlp.Stream) error {
	var proposal = new(Proposal)

	if err := s.Decode(&proposal); err != nil {
		return err
	}
	p.Round, p.Height, p.ValidRound, p.ProposalBlock = proposal.Round, proposal.Height, proposal.ValidRound, proposal.ProposalBlock

	return nil
}

type Subject struct {
	Round  *big.Int
	Height *big.Int
	Digest common.Hash
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (sub *Subject) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{sub.Round, sub.Height, sub.Digest})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (sub *Subject) DecodeRLP(s *rlp.Stream) error {
	var subject = new(Subject)

	if err := s.Decode(&subject); err != nil {
		return err
	}
	sub.Round, sub.Height, sub.Digest = subject.Round, subject.Height, subject.Digest
	return nil
}

func (sub *Subject) String() string {
	return fmt.Sprintf("{Round: %v, Height: %v Digest: %v}", sub.Round, sub.Height, sub.Digest.String())
}
