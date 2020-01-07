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
	"fmt"
	"github.com/pkg/errors"
	"io"
	"math/big"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/rlp"
)

type ConsensusMsg interface {
	GetRound() *big.Int
	GetHeight() *big.Int
}

type Proposal struct {
	Round         *big.Int
	Height        *big.Int
	ValidRound    *big.Int
	ProposalBlock *types.Block
}

func (p *Proposal) GetRound() *big.Int {
	return p.Round
}

func (p *Proposal) GetHeight() *big.Int {
	return p.Height
}

func NewProposal(r *big.Int, h *big.Int, vr *big.Int, p *types.Block) *Proposal {
	return &Proposal{
		Round:         r,
		Height:        h,
		ValidRound:    vr,
		ProposalBlock: p,
	}
}

// RLP encoding doesn't support negative big.Int, so we have to pass one additionnal field to represents validRound = -1.
// Note that we could have as well indexed rounds starting by 1, but we want to stay close as possible to the spec.
func (p *Proposal) EncodeRLP(w io.Writer) error {
	if p.ProposalBlock == nil {
		// Should never happen
		return errors.New("encoderlp proposal with nil block")
	}

	isValidRoundNil := false
	validRound := new(big.Int).Set(p.ValidRound)
	if p.ValidRound.Int64() == -1 {
		validRound = validRound.SetUint64(255) // to make more obvious bad proposals.
		isValidRoundNil = true
	}

	return rlp.Encode(w, []interface{}{
		p.Round,
		p.Height,
		validRound,
		isValidRoundNil,
		p.ProposalBlock,
	})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (p *Proposal) DecodeRLP(s *rlp.Stream) error {
	var proposal struct {
		Round           *big.Int
		Height          *big.Int
		ValidRound      *big.Int
		IsValidRoundNil bool
		ProposalBlock   *types.Block
	}

	if err := s.Decode(&proposal); err != nil {
		return err
	}

	if proposal.IsValidRoundNil {
		if proposal.ValidRound.Uint64() != 255 {
			return errors.New("bad proposal with isValidround nil")
		}
		proposal.ValidRound = big.NewInt(-1)
	}
	if proposal.ProposalBlock == nil {
		return errors.New("bad proposal with nil decoded block")
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

func (sub *Vote) GetRound() *big.Int {
	return sub.Round
}

func (sub *Vote) GetHeight() *big.Int {
	return sub.Height
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
