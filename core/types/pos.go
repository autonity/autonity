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
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
	lru "github.com/hashicorp/golang-lru"
	"io"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/crypto/sha3"
	"github.com/clearmatics/autonity/rlp"
)

var (
	// PoSDigest represents a hash of "PoS practical byzantine fault tolerance"
	// to identify whether the block is from PoS consensus engine
	PoSDigest = common.HexToHash("0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365")

	PoSExtraVanity = 32 // Fixed number of extra-data bytes reserved for validator vanity
	PoSExtraSeal   = 65 // Fixed number of extra-data bytes reserved for validator seal

	inmemoryAddresses  = 20 // Number of recent addresses from ecrecover
	recentAddresses, _ = lru.NewARC(inmemoryAddresses)

	// ErrInvalidPoSHeaderExtra is returned if the length of extra-data is less than 32 bytes
	ErrInvalidPoSHeaderExtra = errors.New("invalid pos header extra-data")
	// ErrInvalidSignature is returned when given signature is not signed by given address.
	ErrInvalidSignature = errors.New("invalid signature")
	// ErrInvalidCommittedSeals is returned if the committed seal is not signed by any of parent validators.
	ErrInvalidCommittedSeals = errors.New("invalid committed seals")
	// ErrEmptyCommittedSeals is returned if the field of committed seals is zero.
	ErrEmptyCommittedSeals = errors.New("zero committed seals")
)

type PoSExtra struct {
	Validators    []common.Address
	Seal          []byte
	CommittedSeal [][]byte
}

// EncodeRLP serializes pos into the Ethereum RLP format.
func (pos *PoSExtra) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		pos.Validators,
		pos.Seal,
		pos.CommittedSeal,
	})
}

// DecodeRLP implements rlp.Decoder, and load the pos fields from a RLP stream.
func (pos *PoSExtra) DecodeRLP(s *rlp.Stream) error {
	var posExtra struct {
		Validators    []common.Address
		Seal          []byte
		CommittedSeal [][]byte
	}
	if err := s.Decode(&posExtra); err != nil {
		return err
	}
	pos.Validators, pos.Seal, pos.CommittedSeal = posExtra.Validators, posExtra.Seal, posExtra.CommittedSeal
	return nil
}

// ExtractPoSExtra extracts all values of the PoSExtra from the header. It returns an
// error if the length of the given extra-data is less than 32 bytes or the extra-data can not
// be decoded.
func ExtractPoSExtra(h *Header) (*PoSExtra, error) {
	if len(h.Extra) < PoSExtraVanity {
		return nil, ErrInvalidPoSHeaderExtra
	}

	var posExtra *PoSExtra
	err := rlp.DecodeBytes(h.Extra[PoSExtraVanity:], &posExtra)
	if err != nil {
		return nil, err
	}
	return posExtra, nil
}

// PoSFilteredHeader returns a filtered header which some information (like seal, committed seals)
// are clean to fulfill the PoS hash rules. It returns nil if the extra-data cannot be
// decoded/encoded by rlp.
func PoSFilteredHeader(h *Header, keepSeal bool) *Header {
	newHeader := CopyHeader(h)
	posExtra, err := ExtractPoSExtra(newHeader)
	if err != nil {
		return nil
	}

	if !keepSeal {
		posExtra.Seal = []byte{}
	}
	posExtra.CommittedSeal = [][]byte{}

	payload, err := rlp.EncodeToBytes(&posExtra)
	if err != nil {
		return nil
	}

	newHeader.Extra = append(newHeader.Extra[:PoSExtraVanity], payload...)

	return newHeader
}

// sigHash returns the hash which is used as input for the PoS
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func SigHash(header *Header) (hash common.Hash) {
	hasher := sha3.NewKeccak256()

	// Clean seal is required for calculating proposer seal.
	err := rlp.Encode(hasher, PoSFilteredHeader(header, false))
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
	posExtra, err := ExtractPoSExtra(header)
	if err != nil {
		return common.Address{}, err
	}

	addr, err := GetSignatureAddress(SigHash(header).Bytes(), posExtra.Seal)
	if err != nil {
		return addr, err
	}
	recentAddresses.Add(hash, addr)
	return addr, nil
}

// PrepareExtra returns a extra-data of the given header and validators
func PrepareExtra(extraData *[]byte, vals []common.Address) ([]byte, error) {
	var buf bytes.Buffer

	// compensate the lack bytes if header.Extra is not enough PoSExtraVanity bytes.
	if len(*extraData) < PoSExtraVanity {
		*extraData = append(*extraData, bytes.Repeat([]byte{0x00}, PoSExtraVanity-len(*extraData))...)
	}
	buf.Write((*extraData)[:PoSExtraVanity])

	pos := &PoSExtra{
		Validators:    vals,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}

	payload, err := rlp.EncodeToBytes(&pos)
	if err != nil {
		return nil, err
	}
	return append(buf.Bytes(), payload...), nil
}

// WriteSeal writes the extra-data field of the given header with the given seals.
// suggest to rename to writeSeal.
func WriteSeal(h *Header, seal []byte) error {
	if len(seal)%PoSExtraSeal != 0 { // TODO :  len(seal) != PoSExtraSeal
		return ErrInvalidSignature
	}

	posExtra, err := ExtractPoSExtra(h)
	if err != nil {
		return err
	}

	posExtra.Seal = seal
	payload, err := rlp.EncodeToBytes(&posExtra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:PoSExtraVanity], payload...)
	return nil
}

// WriteCommittedSeals writes the extra-data field of a block header with given committed seals.
func WriteCommittedSeals(h *Header, committedSeals [][]byte) error {
	if len(committedSeals) == 0 {
		return ErrInvalidCommittedSeals
	}

	for _, seal := range committedSeals {
		if len(seal) != PoSExtraSeal {
			return ErrInvalidCommittedSeals
		}
	}

	posExtra, err := ExtractPoSExtra(h)
	if err != nil {
		return err
	}
	posExtra.CommittedSeal = make([][]byte, len(committedSeals))
	copy(posExtra.CommittedSeal, committedSeals)

	payload, err := rlp.EncodeToBytes(&posExtra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:PoSExtraVanity], payload...)
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
