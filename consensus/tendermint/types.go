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
	ProposalBlock *types.Block
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
	var proposal struct {
		Round         *big.Int
		Height        *big.Int
		ValidRound    *big.Int
		ProposalBlock *types.Block
	}

	if err := s.Decode(&proposal); err != nil {
		return err
	}
	p.Round = proposal.Round
	p.Height = proposal.Height
	p.ValidRound = proposal.ValidRound
	p.ProposalBlock = proposal.ProposalBlock

	return nil
}

type Vote struct {
	Round             *big.Int
	Height            *big.Int
	ProposedBlockHash common.Hash
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (sub *Vote) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{sub.Round, sub.Height, sub.ProposedBlockHash})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (sub *Vote) DecodeRLP(s *rlp.Stream) error {
	var vote struct {
		Round             *big.Int
		Height            *big.Int
		ProposedBlockHash common.Hash
	}

	if err := s.Decode(&vote); err != nil {
		return err
	}
	sub.Round = vote.Round
	sub.Height = vote.Height
	sub.ProposedBlockHash = vote.ProposedBlockHash
	return nil
}

func (sub *Vote) String() string {
	return fmt.Sprintf("{Round: %v, Height: %v ProposedBlockHash: %v}", sub.Round, sub.Height, sub.ProposedBlockHash.String())
}
