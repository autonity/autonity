// Copyright 2014 The go-ethereum Authors
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

package common

import "math/big"

// Common big integers often used
var (
	Big0   = big.NewInt(0)
	Big1   = big.NewInt(1)
	Big2   = big.NewInt(2)
	Big3   = big.NewInt(3)
	Big4   = big.NewInt(4)
	Big5   = big.NewInt(5)
	Big32  = big.NewInt(32)
	Big256 = big.NewInt(256)
)

// inspired by slices.Contains, but compares value instead of pointer
// elements are assumed to be != nil, responsibility to check is of the caller
func Contains(s []*big.Int, v *big.Int) bool {
	return Index(s, v) >= 0
}

// Index returns the index of the first occurrence of v in s,
// or -1 if not present.
func Index(s []*big.Int, v *big.Int) int {
	for i := range s {
		if v.Cmp(s[i]) == 0 {
			return i
		}
	}
	return -1
}
