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
	"reflect"
	"testing"

	"github.com/clearmatics/autonity/common"
)

func TestHeaderHash(t *testing.T) {
	// 0xcefefd3ade63a5955bca4562ed840b67f39e74df217f7e5f7241a6e9552cca70
	expectedExtra := common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000000f89af8549444add0ec310f115a0e603b2d7db9f067778eaf8a94294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212946beaaed781d2d2ab6350f5c4566a2c6eaac407a6948be76812f765c24641ec63dc2852b378aba2b440b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0")
	expectedHash := common.HexToHash("0xcefefd3ade63a5955bca4562ed840b67f39e74df217f7e5f7241a6e9552cca70")

	// for pos consensus
	header := &Header{MixDigest: BFTDigest, Extra: expectedExtra}
	if !reflect.DeepEqual(header.Hash(), expectedHash) {
		t.Errorf("expected: %v, but got: %v", expectedHash.Hex(), header.Hash().Hex())
	}

	// append useless information to extra-data
	unexpectedExtra := append(expectedExtra, []byte{1, 2, 3}...)
	header.Extra = unexpectedExtra
	if !reflect.DeepEqual(header.Hash(), rlpHash(header)) {
		t.Errorf("expected: %v, but got: %v", rlpHash(header).Hex(), header.Hash().Hex())
	}
}
