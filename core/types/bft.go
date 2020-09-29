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

package types

import (
	"errors"
	"strings"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/rlp"
	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/crypto/sha3"
)

var (
	// BFTDigest represents a hash of "Istanbul practical byzantine fault tolerance"
	// to identify whether the block is from BFT consensus engine
	//TODO differentiate the digest between IBFT and Tendermint
	BFTDigest = common.HexToHash("0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365")

	BFTExtraVanity = 32 // Fixed number of extra-data bytes reserved for validator vanity
	BFTExtraSeal   = 65 // Fixed number of extra-data bytes reserved for validator seal

	inmemoryAddresses  = 500 // Number of recent addresses from ecrecover
	recentAddresses, _ = lru.NewARC(inmemoryAddresses)

	// ErrInvalidBFTHeaderExtra is returned if the length of extra-data is less than 32 bytes
	ErrInvalidBFTHeaderExtra = errors.New("invalid pos header extra-data")
	// ErrInvalidSignature is returned when given signature is not signed by given address.
	ErrInvalidSignature = errors.New("invalid signature")
	// ErrInvalidCommittedSeals is returned if the committed seal is not signed by any of parent validators.
	ErrInvalidCommittedSeals = errors.New("invalid committed seals")
	// ErrEmptyCommittedSeals is returned if the field of committed seals is zero.
	ErrEmptyCommittedSeals = errors.New("zero committed seals")
	// ErrNegativeRound is returned if the round field is negative
	ErrNegativeRound = errors.New("negative round")
)

// BFTFilteredHeader returns a filtered header which some information (like seal, committed seals)
// are clean to fulfill the BFT hash rules.
func BFTFilteredHeader(h *Header, keepSeal bool) *Header {
	newHeader := CopyHeader(h)
	if !keepSeal {
		newHeader.ProposerSeal = []byte{}
	}
	newHeader.CommittedSeals = [][]byte{}
	newHeader.Round = 0
	newHeader.Extra = []byte{}
	return newHeader
}

// sigHash returns the hash which is used as input for the BFT
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func SigHash(header *Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()

	// Clean seal is required for calculating proposer seal.
	err := rlp.Encode(hasher, BFTFilteredHeader(header, false))
	if err != nil {
		log.Error("can't hash the header", "err", err, "header", header)
		return common.Hash{}
	}
	hasher.Sum(hash[:0])
	return hash
}

// Ecrecover extracts the Ethereum account address from a signed header.
func Ecrecover(header *Header) (common.Address, error) {
	hash := header.Hash()
	if addr, ok := recentAddresses.Get(hash); ok {
		return addr.(common.Address), nil
	}

	addr, err := GetSignatureAddress(SigHash(header).Bytes(), header.ProposerSeal)
	if err != nil {
		return addr, err
	}
	recentAddresses.Add(hash, addr)
	return addr, nil
}

// WriteSeal writes the extra-data field of the given header with the given seals.
// suggest to rename to writeSeal.
func WriteSeal(h *Header, seal []byte) error {
	if len(seal) != BFTExtraSeal {
		return ErrInvalidSignature
	}
	h.ProposerSeal = make([]byte, len(seal))
	copy(h.ProposerSeal, seal)
	return nil
}

// WriteRound writes the round field of the block header.
func WriteRound(h *Header, round int64) error {
	if round < 0 {
		return ErrNegativeRound
	}
	h.Round = uint64(round)
	return nil
}

// WriteCommittedSeals writes the extra-data field of a block header with given committed seals.
func WriteCommittedSeals(h *Header, committedSeals [][]byte) error {
	if len(committedSeals) == 0 {
		return ErrInvalidCommittedSeals
	}

	for _, seal := range committedSeals {
		if len(seal) != BFTExtraSeal {
			return ErrInvalidCommittedSeals
		}
	}

	h.CommittedSeals = committedSeals
	return nil
}

func RLPHash(v interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, v)
	hw.Sum(h[:0])
	return h
}

// GetSignatureAddress gets the signer address from the signature
func GetSignatureAddress(data []byte, sig []byte) (common.Address, error) {
	// 1. Keccak data
	hashData := crypto.Keccak256(data)
	// 2. Recover public key
	pubkey, err := crypto.SigToPub(hashData, sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pubkey), nil
}

func (c Committee) String() string {
	var ret string
	for _, val := range c {
		ret += "[" + val.Address.String() + " - " + val.VotingPower.String() + "] "
	}
	return ret
}

func (c Committee) Len() int {
	return len(c)
}

func (c Committee) Less(i, j int) bool {
	return strings.Compare(c[i].String(), c[j].String()) < 0
}

func (c Committee) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (m CommitteeMember) String() string {
	return m.Address.String()
}
