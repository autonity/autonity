package message

import (
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
)

// represents senders of an aggregated signature

type SendersInfo struct {
	// TODO(lorenzo) refactor into a more efficient data structure
	// e.g. 2 bit for each validator. 00 --> no key, 01 --> 1 key, 10 or 11 --> coefficient of key appended at the end of the initial array
	bits []byte
	// these fields are not serialized, but instead computed at preValidate steps
	// the order matches the != 0 elements in the bitmap
	// it should always be len(addresses) = len(powers) = len(bitmap elements != 0)
	addresses []common.Address `rlp:"-"` //TODO(lorenzo) verify with test
	powers    []*big.Int       `rlp:"-"`
}

//TODO(lorenzo) fix formatting of these comments

// note: bits will always be same length as current committee size
// addresses and power can be shorter (they are elements only for the non-0 ones)
// note that this means that we cannot do random access in addresses and powers array
// e.g. bits [ 0 1 0 2 0 0 ]
//
//	addresses [ 0x1 0x2]
//
// powers [0x1 0x2]
func NewSendersInfo(committeeSize int) SendersInfo {
	return &SendersInfo{bits: make([]byte, committeeSize)}
}

// Len returns how many committee members the sendersInfo contains
// e.g. [0 1 2 0 1] --> len = 5
//
//				[0 9 9 9 2] --> len = 5
//	     [0 1 2]     --> len = 3
func (s SendersInfo) Len() int {
	return len(s.bits)
}

func (s SendersInfo) Bits() []byte {
	return s.bits
}
func (s SendersInfo) Addresses() []common.Address {
	return s.addresses
}

// returns aggregated power of all senders
func (s SendersInfo) Power() *big.Int {
	power := new(big.Int)
	for i, v := range s.bits {
		// regardless of the value, we sum power only once (to prevent counting duplicated messages power twice)
		if v > 0 {
			power.Add(power, s.powers[i])
		}

	}
	return power
}

func (s SendersInfo) Powers() []*big.Int {
	return s.powers
}

// TODO(lorenzo) what if different size
func (s SendersInfo) Merge(other SendersInfo) {
	for i, coefficient := range other.bitmap {
		if coefficient == 0 {
			continue
		}
		if s[i] == 0 {
			s.addresses = append(s.addresses[:i+1], s.addresses[i:]...)
			s.addresses[i] = other.addresses[i]
			s.powers = append(s.powers[:i+1], s.powers[i:]...)
			s.powers[i] = other.powers[i]
		}
		s[i] = s[i] + coefficient
	}
}

func (s SendersInfo) Increment(parent *types.Header, i uint64) {
	s.bits[i]++
	s.addresses = append(s.addresses[:i+1], s.addresses[i:]...)
	s.addresses[i] = parent.Committee[i].Address
	s.powers = append(s.powers[:i+1], s.powers[i:]...)
	s.powers[i] = parent.Committee[i].VotingPower
}

// returns list of indexes of validators that signed
// e.g. for bitmap [0 1 2 1 0] will return [ 1 2 2 3 ]
// the index 2 is repeated because we need to aggregate two times his key
func (s SendersInfo) Flatten() []uint64 {
	var indexes []uint64
	for i, v := range s.bits {
		for j := 0; j < int(v); j++ {
			indexes = append(indexes, uint64(i))
		}
	}
	return indexes
}

// same as before, but repeated indexes are returned only once
func (s SendersInfo) FlattenUniq() []uint64 {
	var indexes []uint64
	for i, v := range s.bits {
		if v > 0 {
			indexes = append(indexes, uint64(i))
		}
	}
	return indexes
}
