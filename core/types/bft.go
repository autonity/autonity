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
	"bytes"
	"errors"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/crypto/sha3"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/rlp"
	lru "github.com/hashicorp/golang-lru"
	"io"
)

var (
	// BFTDigest represents a hash of "Istanbul practical byzantine fault tolerance"
	// to identify whether the block is from BFT consensus engine
	//TODO differentiate the digest between IBFT and Tendermint
	BFTDigest = common.HexToHash("0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365")

	BFTExtraVanity = 32 // Fixed number of extra-data bytes reserved for validator vanity
	BFTExtraSeal   = 65 // Fixed number of extra-data bytes reserved for validator seal

	inmemoryAddresses  = 20 // Number of recent addresses from ecrecover
	recentAddresses, _ = lru.NewARC(inmemoryAddresses)

	// ErrInvalidBFTHeaderExtra is returned if the length of extra-data is less than 32 bytes
	ErrInvalidBFTHeaderExtra = errors.New("invalid pos header extra-data")
	// ErrInvalidSignature is returned when given signature is not signed by given address.
	ErrInvalidSignature = errors.New("invalid signature")
	// ErrInvalidCommittedSeals is returned if the committed seal is not signed by any of parent validators.
	ErrInvalidCommittedSeals = errors.New("invalid committed seals")
	// ErrEmptyCommittedSeals is returned if the field of committed seals is zero.
	ErrEmptyCommittedSeals = errors.New("zero committed seals")
)

type BFTExtra struct {
	Validators    []common.Address
	Seal          []byte
	CommittedSeal [][]byte
}

// EncodeRLP serializes pos into the Ethereum RLP format.
func (pos *BFTExtra) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		pos.Validators,
		pos.Seal,
		pos.CommittedSeal,
	})
}

// DecodeRLP implements rlp.Decoder, and load the pos fields from a RLP stream.
func (pos *BFTExtra) DecodeRLP(s *rlp.Stream) error {
	var bftExtra struct {
		Validators    []common.Address
		Seal          []byte
		CommittedSeal [][]byte
	}
	if err := s.Decode(&bftExtra); err != nil {
		return err
	}
	pos.Validators, pos.Seal, pos.CommittedSeal = bftExtra.Validators, bftExtra.Seal, bftExtra.CommittedSeal
	return nil
}

// ExtractBFTExtra extracts all values of the BFTExtra from the header. It returns an
// error if the length of the given extra-data is less than 32 bytes or the extra-data can not
// be decoded.
func ExtractBFTHeaderExtra(h *Header) (*BFTExtra, error) {
	return ExtractBFTExtra(h.Extra)
}

func ExtractBFTExtra(extra []byte) (*BFTExtra, error) {
	if len(extra) < BFTExtraVanity {
		return nil, ErrInvalidBFTHeaderExtra
	}

	var bftExtra *BFTExtra
	err := rlp.DecodeBytes(extra[BFTExtraVanity:], &bftExtra)
	if err != nil {
		return nil, err
	}
	return bftExtra, nil
}

// BFTFilteredHeader returns a filtered header which some information (like seal, committed seals)
// are clean to fulfill the BFT hash rules. It returns nil if the extra-data cannot be
// decoded/encoded by rlp.
func BFTFilteredHeader(h *Header, keepSeal bool) *Header {
	newHeader := CopyHeader(h)
	bftExtra, err := ExtractBFTHeaderExtra(newHeader)
	if err != nil {
		return nil
	}

	if !keepSeal {
		bftExtra.Seal = []byte{}
	}
	bftExtra.CommittedSeal = [][]byte{}

	payload, err := rlp.EncodeToBytes(&bftExtra)
	if err != nil {
		return nil
	}

	newHeader.Extra = append(newHeader.Extra[:BFTExtraVanity], payload...)

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
	hasher := sha3.NewKeccak256()

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

	// Retrieve the signature from the header extra-data
	bftExtra, err := ExtractBFTHeaderExtra(header)
	if err != nil {
		return common.Address{}, err
	}

	addr, err := GetSignatureAddress(SigHash(header).Bytes(), bftExtra.Seal)
	if err != nil {
		return addr, err
	}
	recentAddresses.Add(hash, addr)
	return addr, nil
}

// PrepareExtra returns a extra-data of the given header and validators
func PrepareExtra(extraData []byte, vals []common.Address) ([]byte, error) {
	extraDataCopy := append([]byte{}, extraData...)

	pos := &BFTExtra{
		Validators:    vals,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}

	payload, err := rlp.EncodeToBytes(&pos)
	if err != nil {
		return extraDataCopy, err
	}

	// compensate the lack bytes if header.Extra is not enough BFTExtraVanity bytes.
	if len(extraDataCopy) < BFTExtraVanity {
		extraDataCopy = append(extraDataCopy, bytes.Repeat([]byte{0x00}, BFTExtraVanity-len(extraDataCopy))...)
	}

	if len(extraDataCopy) > BFTExtraVanity {
		extraDataCopy = extraDataCopy[:BFTExtraVanity]
	}

	return append(extraDataCopy, payload...), nil
}

// WriteSeal writes the extra-data field of the given header with the given seals.
// suggest to rename to writeSeal.
func WriteSeal(h *Header, seal []byte) error {
	if len(seal)%BFTExtraSeal != 0 { // TODO :  len(seal) != BFTExtraSeal
		return ErrInvalidSignature
	}

	bftExtra, err := ExtractBFTHeaderExtra(h)
	if err != nil {
		return err
	}

	bftExtra.Seal = seal
	payload, err := rlp.EncodeToBytes(&bftExtra)
	if err != nil {
		panic(3)
		return err
	}

	h.Extra = append(h.Extra[:BFTExtraVanity], payload...)
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

	bftExtra, err := ExtractBFTHeaderExtra(h)
	if err != nil {
		return err
	}
	bftExtra.CommittedSeal = make([][]byte, len(committedSeals))
	copy(bftExtra.CommittedSeal, committedSeals)

	payload, err := rlp.EncodeToBytes(&bftExtra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:BFTExtraVanity], payload...)
	return nil
}

func RLPHash(v interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
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
