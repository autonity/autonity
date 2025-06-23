// Copyright 2015 The go-ethereum Authors
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

package rlp

import (
	"io"
	"math"
	"reflect"

	"github.com/autonity/autonity/common"
)

// RawValue represents an encoded RLP value and can be used to delay
// RLP decoding or to precompute an encoding. Note that the decoder does
// not verify whether the content of RawValues is valid RLP.
type RawValue []byte

var rawValueType = reflect.TypeOf(RawValue{})

// ListSize returns the encoded size of an RLP list with the given
// content size.
func ListSize(contentSize uint64) uint64 {
	return uint64(headsize(contentSize)) + contentSize
}

// IntSize returns the encoded size of the integer x.
func IntSize(x uint64) int {
	if x < 0x80 {
		return 1
	}
	return 1 + intsize(x)
}

// Split returns the content of first RLP value and any
// bytes after the value as subslices of b.
func Split(b []byte) (k Kind, content, rest []byte, err error) {
	k, ts, cs, err := readKind(b)
	if err != nil {
		return 0, nil, b, err
	}
	return k, b[ts : ts+cs], b[ts+cs:], nil
}

// SplitString splits b into the content of an RLP string
// and any remaining bytes after the string.
func SplitString(b []byte) (content, rest []byte, err error) {
	k, content, rest, err := Split(b)
	if err != nil {
		return nil, b, err
	}
	if k == List {
		return nil, b, ErrExpectedString
	}
	return content, rest, nil
}

// SplitUint64 decodes an integer at the beginning of b.
// It also returns the remaining data after the integer in 'rest'.
func SplitUint64(b []byte) (x uint64, rest []byte, err error) {
	content, rest, err := SplitString(b)
	if err != nil {
		return 0, b, err
	}
	switch {
	case len(content) == 0:
		return 0, rest, nil
	case len(content) == 1:
		if content[0] == 0 {
			return 0, b, ErrCanonInt
		}
		return uint64(content[0]), rest, nil
	case len(content) > 8:
		return 0, b, errUintOverflow
	default:
		x, err = readSize(content, byte(len(content)))
		if err != nil {
			return 0, b, ErrCanonInt
		}
		return x, rest, nil
	}
}

// SplitList splits b into the content of a list and any remaining
// bytes after the list.
func SplitList(b []byte) (content, rest []byte, err error) {
	k, content, rest, err := Split(b)
	if err != nil {
		return nil, b, err
	}
	if k != List {
		return nil, b, ErrExpectedList
	}
	return content, rest, nil
}

// CountValues counts the number of encoded values in b.
func CountValues(b []byte) (int, error) {
	i := 0
	for ; len(b) > 0; i++ {
		_, tagsize, size, err := readKind(b)
		if err != nil {
			return 0, err
		}
		b = b[tagsize+size:]
	}
	return i, nil
}

func readKind(buf []byte) (k Kind, tagsize, contentsize uint64, err error) {
	if len(buf) == 0 {
		return 0, 0, 0, io.ErrUnexpectedEOF
	}
	b := buf[0]
	switch {
	case b < 0x80:
		k = Byte
		tagsize = 0
		contentsize = 1
	case b < 0xB8:
		k = String
		tagsize = 1
		contentsize = uint64(b - 0x80)
		// Reject strings that should've been single bytes.
		if contentsize == 1 && len(buf) > 1 && buf[1] < 128 {
			return 0, 0, 0, ErrCanonSize
		}
	case b < 0xC0:
		k = String
		tagsize = uint64(b-0xB7) + 1
		contentsize, err = readSize(buf[1:], b-0xB7)
	case b < 0xF8:
		k = List
		tagsize = 1
		contentsize = uint64(b - 0xC0)
	default:
		k = List
		tagsize = uint64(b-0xF7) + 1
		contentsize, err = readSize(buf[1:], b-0xF7)
	}
	if err != nil {
		return 0, 0, 0, err
	}
	// Reject values larger than the input slice.
	if contentsize > uint64(len(buf))-tagsize {
		return 0, 0, 0, ErrValueTooLarge
	}
	return k, tagsize, contentsize, err
}

func readSize(b []byte, slen byte) (uint64, error) {
	if int(slen) > len(b) {
		return 0, io.ErrUnexpectedEOF
	}
	var s uint64
	switch slen {
	case 1:
		s = uint64(b[0])
	case 2:
		s = uint64(b[0])<<8 | uint64(b[1])
	case 3:
		s = uint64(b[0])<<16 | uint64(b[1])<<8 | uint64(b[2])
	case 4:
		s = uint64(b[0])<<24 | uint64(b[1])<<16 | uint64(b[2])<<8 | uint64(b[3])
	case 5:
		s = uint64(b[0])<<32 | uint64(b[1])<<24 | uint64(b[2])<<16 | uint64(b[3])<<8 | uint64(b[4])
	case 6:
		s = uint64(b[0])<<40 | uint64(b[1])<<32 | uint64(b[2])<<24 | uint64(b[3])<<16 | uint64(b[4])<<8 | uint64(b[5])
	case 7:
		s = uint64(b[0])<<48 | uint64(b[1])<<40 | uint64(b[2])<<32 | uint64(b[3])<<24 | uint64(b[4])<<16 | uint64(b[5])<<8 | uint64(b[6])
	case 8:
		s = uint64(b[0])<<56 | uint64(b[1])<<48 | uint64(b[2])<<40 | uint64(b[3])<<32 | uint64(b[4])<<24 | uint64(b[5])<<16 | uint64(b[6])<<8 | uint64(b[7])
	}
	// Reject sizes < 56 (shouldn't have separate size) and sizes with
	// leading zero bytes.
	if s < 56 || b[0] == 0 {
		return 0, ErrCanonSize
	}
	return s, nil
}

// AppendUint64 appends the RLP encoding of i to b, and returns the resulting slice.
func AppendUint64(b []byte, i uint64) []byte {
	if i == 0 {
		return append(b, 0x80)
	} else if i < 128 {
		return append(b, byte(i))
	}
	switch {
	case i < (1 << 8):
		return append(b, 0x81, byte(i))
	case i < (1 << 16):
		return append(b, 0x82,
			byte(i>>8),
			byte(i),
		)
	case i < (1 << 24):
		return append(b, 0x83,
			byte(i>>16),
			byte(i>>8),
			byte(i),
		)
	case i < (1 << 32):
		return append(b, 0x84,
			byte(i>>24),
			byte(i>>16),
			byte(i>>8),
			byte(i),
		)
	case i < (1 << 40):
		return append(b, 0x85,
			byte(i>>32),
			byte(i>>24),
			byte(i>>16),
			byte(i>>8),
			byte(i),
		)

	case i < (1 << 48):
		return append(b, 0x86,
			byte(i>>40),
			byte(i>>32),
			byte(i>>24),
			byte(i>>16),
			byte(i>>8),
			byte(i),
		)
	case i < (1 << 56):
		return append(b, 0x87,
			byte(i>>48),
			byte(i>>40),
			byte(i>>32),
			byte(i>>24),
			byte(i>>16),
			byte(i>>8),
			byte(i),
		)

	default:
		return append(b, 0x88,
			byte(i>>56),
			byte(i>>48),
			byte(i>>40),
			byte(i>>32),
			byte(i>>24),
			byte(i>>16),
			byte(i>>8),
			byte(i),
		)
	}
}

// NOTE: l needs to be a list. The list is going to be expanded and the address appended at the end
func AppendAddress(original []byte, addr common.Address) ([]byte, error) {
	// copy list so that original one is not partially modified
	list := make([]byte, len(original), len(original)+common.AddressLength)
	copy(list, original)

	k, tagSize, contentSize, err := readKind(list)
	if err != nil {
		return nil, err
	}
	if k == Byte || k == String {
		return nil, ErrExpectedList
	}
	// rule out extreme edge case
	if contentSize > math.MaxUint64-common.AddressLength {
		return nil, ErrCannotAppend
	}

	newContentSize := contentSize + common.AddressLength
	newTagSize := uint64(headsize(newContentSize))

	switch {
	case tagSize == 1 && newTagSize == 1:
		// remains short list, just need to update tag
		list[0] = byte(0xC0 + newContentSize)
	case tagSize == 1 && newTagSize > 1:
		// becomes long list, next case will update the tag size
		// here update only the first byte of the tag
		list[0] = byte(0xF7)
		fallthrough
	default: // tagSize > 1 && newTagSize > 1 || falling through
		// if the new tag size is bigger than the previous one need to:
		// - update first byte
		// - add bytes after the first one
		if newTagSize > tagSize {
			list = append([]byte{0xF7 + byte(newTagSize-1)}, append(make([]byte, newTagSize-1), list[tagSize:]...)...)
		}
		// update the new size value
		putint(list[1:], newContentSize)
	}
	// append the actual address
	list = append(list, addr[:]...)
	return list, nil
}

func ExtractAddress(original []byte) ([]byte, common.Address, error) {
	k, tagSize, contentSize, err := readKind(original)
	if err != nil {
		return nil, common.Address{}, err
	}
	if k == Byte || k == String {
		return nil, common.Address{}, ErrExpectedList
	}

	if contentSize < common.AddressLength {
		return nil, common.Address{}, ErrCannotExtract
	}

	// copy list so that original one is not partially modified
	list := make([]byte, len(original))
	copy(list, original)

	newContentSize := contentSize - common.AddressLength
	newTagSize := uint64(headsize(newContentSize))

	switch {
	case tagSize == 1 && newTagSize == 1:
		// remains short list
		list[0] = byte(0xC0 + newContentSize)
	case tagSize > 1 && newTagSize == 1:
		// long list becomes short list
		list = append([]byte{0xC0 + byte(newContentSize)}, list[tagSize:]...)
	default: // tagSize > 1 && newTagSize > 1:
		// if the new tag size is smaller than the previous one, need to:
		// - update first byte
		// - remove excess bytes after the first one
		if newTagSize < tagSize {
			list = append([]byte{0xF7 + byte(newTagSize-1)}, list[1+tagSize-newTagSize:]...)
		}

		// update the new size value
		putint(list[1:], newContentSize)
	}
	// extract address and remove it from the list
	list = list[:len(list)-common.AddressLength]
	address := common.BytesToAddress(list[len(list)-common.AddressLength:])
	return list, address, nil
}
