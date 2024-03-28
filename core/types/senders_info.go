package types

import (
	"math/big"

	"github.com/autonity/autonity/common"
)

// represents senders of an aggregated signature

type SendersInfo struct {
	// TODO(lorenzo) refactor into a more efficient data structure
	// e.g. 2 bit for each validator. 00 --> no key, 01 --> 1 key, 10 or 11 --> coefficient of key appended at the end of the initial array
	Bits []byte
	// these fields are not serialized, but instead computed at preValidate steps
	addresses map[int]common.Address `rlp:"-"` //TODO(lorenzo) verify with test
	powers    map[int]*big.Int       `rlp:"-"`
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
func NewSendersInfo(committeeSize int) *SendersInfo {
	return &SendersInfo{Bits: make([]byte, committeeSize), addresses: make(map[int]common.Address), powers: make(map[int]*big.Int)}
}

// Len returns how many committee members the sendersInfo contains
// e.g. [0 1 2 0 1] --> len = 5
//
//				[0 9 9 9 2] --> len = 5
//	     [0 1 2]     --> len = 3
func (s *SendersInfo) Len() int {
	return len(s.Bits)
}

/*
	 //TODO(lorenzo) dekete
		func (s SendersInfo) Bits() []byte {
			return s.bits
		}
*/
func (s *SendersInfo) Addresses() map[int]common.Address {
	return s.addresses
}
func (s *SendersInfo) SetAddresses(addresses map[int]common.Address) {
	s.addresses = addresses
}

// returns aggregated power of all senders
func (s SendersInfo) Power() *big.Int {
	power := new(big.Int)
	for i, v := range s.Bits {
		// regardless of the value, we sum power only once (to prevent counting duplicated messages power twice)
		if v > 0 {
			power.Add(power, s.powers[i])
		}
	}
	return power
}

func (s *SendersInfo) Powers() map[int]*big.Int {
	return s.powers
}

func (s *SendersInfo) SetPowers(powers map[int]*big.Int) {
	s.powers = powers
}

func (s *SendersInfo) Copy() *SendersInfo {
	addresses := make(map[int]common.Address)
	for index, address := range s.addresses {
		addresses[index] = address
	}
	powers := make(map[int]*big.Int)
	for index, power := range s.powers {
		powers[index] = power
	}
	return &SendersInfo{
		Bits:      append(s.Bits[:0:0], s.Bits...),
		addresses: s.addresses,
		powers:    s.powers,
	}
}

func (s *SendersInfo) Merge(other *SendersInfo) {
	if len(s.Bits) != len(other.Bits) {
		// should always merge for the same height --> same committee size --> same legnth
		panic("trying to merge sender information with different length")
	}
	for i, coefficient := range other.Bits {
		if coefficient == 0 {
			continue
		}
		s.addresses[i] = other.addresses[i]
		s.powers[i] = other.powers[i]
		s.Bits[i] = s.Bits[i] + coefficient
	}
}

func (s *SendersInfo) Increment(parent *Header, index int) {
	s.Bits[index]++
	s.addresses[index] = parent.Committee[index].Address
	s.powers[index] = parent.Committee[index].VotingPower
}

// returns list of indexes of validators that signed
// e.g. for bitmap [0 1 2 1 0] will return [ 1 2 2 3 ]
// the index 2 is repeated because we need to aggregate two times his key
func (s *SendersInfo) Flatten() []int {
	var indexes []int
	for i, v := range s.Bits {
		for j := 0; j < int(v); j++ {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

// same as before, but repeated indexes are returned only once
func (s *SendersInfo) FlattenUniq() []int {
	var indexes []int
	for i, v := range s.Bits {
		if v > 0 {
			indexes = append(indexes, i)
		}
	}
	return indexes
}
