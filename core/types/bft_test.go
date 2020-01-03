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
	"math/big"
	"reflect"
	"testing"

	"github.com/clearmatics/autonity/common"
)

func TestHeaderHash(t *testing.T) {
	originalHeader := OriginalHeader{
		ParentHash:  common.HexToHash("0000H45H"),
		UncleHash:   common.HexToHash("0000H45H"),
		Coinbase:    common.HexToAddress("0000H45H"),
		Root:        common.HexToHash("0000H00H"),
		TxHash:      common.HexToHash("0000H45H"),
		ReceiptHash: common.HexToHash("0000H45H"),
		Difficulty:  big.NewInt(1337),
		Number:      big.NewInt(1337),
		GasLimit:    1338,
		GasUsed:     1338,
		Time:        1338,
		Extra:       []byte("Extra data Extra data Extra data  Extra data  Extra data  Extra data  Extra data Extra data"),
		MixDigest:   common.HexToHash("0x0000H45H"),
	}
	PosHeader := originalHeader
	PosHeader.MixDigest = BFTDigest

	originalHeaderHash := common.HexToHash("0x44381ab449d77774874aca34634cb53bc21bd22aef2d3d4cf40e51176cb585ec")
	posHeaderHash := common.HexToHash("0x4ceafbc550a2f60288e7bdfef92a71a65346d184304b526e28cc56a478e12080")

	testCases := []struct {
		header Header
		hash   common.Hash
	}{
		// Non-BFT header tests, PoS fields should not be taken into account.
		{
			Header{},
			common.HexToHash("0xc3bd2d00745c03048a5616146a96f5ff78e54efb9e5b04af208cdaff6f3830ee"),
		},
		{
			Header{OriginalHeader: originalHeader},
			originalHeaderHash,
		},
		{
			Header{OriginalHeader: originalHeader, Committee: Committee{
				{
					Address:     common.HexToAddress("0x1234566"),
					VotingPower: new(big.Int).SetUint64(12),
				},
				{
					Address:     common.HexToAddress("0x13371337"),
					VotingPower: new(big.Int).SetUint64(1337),
				},
			}},
			originalHeaderHash,
		},
		{
			Header{OriginalHeader: originalHeader, ProposerSeal: common.Hex2Bytes("0xbebedead")},
			originalHeaderHash,
		},
		{
			Header{OriginalHeader: originalHeader, Round: new(big.Int).SetUint64(1997)},
			originalHeaderHash,
		},
		{
			Header{OriginalHeader: originalHeader, CommittedSeals: [][]byte{common.Hex2Bytes("0xfacebooc"), common.Hex2Bytes("0xbabababa")}},
			originalHeaderHash,
		},
		{
			Header{OriginalHeader: originalHeader, PastCommittedSeals: [][]byte{common.Hex2Bytes("0xfacebooc"), common.Hex2Bytes("0xbabababa")}},
			originalHeaderHash,
		},
		// BFT header tests
		{
			Header{OriginalHeader: PosHeader}, // test 7
			posHeaderHash,
		},
		{
			Header{OriginalHeader: PosHeader, CommittedSeals: [][]byte{common.Hex2Bytes("0xfacebooc"), common.Hex2Bytes("0xbabababa")}},
			posHeaderHash,
		},
		{
			Header{OriginalHeader: PosHeader, CommittedSeals: [][]byte{common.Hex2Bytes("0x123456"), common.Hex2Bytes("0x777777"), common.Hex2Bytes("0xaaaaaaa")}},
			posHeaderHash,
		},
		{
			Header{OriginalHeader: PosHeader, Committee: Committee{
				{
					Address:     common.HexToAddress("0x1234566"),
					VotingPower: new(big.Int).SetUint64(12),
				},
				{
					Address:     common.HexToAddress("0x13371337"),
					VotingPower: new(big.Int).SetUint64(1337),
				},
			}},
			common.HexToHash("0xf5d460ed44edb6c81ab9ff1979126704e18777986c064d0023aa87bb4a2a7ea5"),
		},
		{
			Header{OriginalHeader: PosHeader, ProposerSeal: common.Hex2Bytes("0xbebedead")},
			common.HexToHash("0x4ceafbc550a2f60288e7bdfef92a71a65346d184304b526e28cc56a478e12080"),
		},
		{
			Header{OriginalHeader: PosHeader, Round: new(big.Int).SetUint64(1997)},
			common.HexToHash("0x79622628fd1cd12bf33623fa6deff06778e65e1f3815bf20cb3dadb026732380"),
		},

		{
			Header{OriginalHeader: PosHeader, PastCommittedSeals: [][]byte{common.Hex2Bytes("0xfacebooc"), common.Hex2Bytes("0xbabababa")}},
			common.HexToHash("0x5d29fd91067324583e8203615ca019679ca5024b8d91cfb3f9710feffd65b6d2"),
		},
	}
	for i := range testCases {
		if !reflect.DeepEqual(testCases[i].hash, testCases[i].header.Hash()) {
			t.Errorf("test %d, expected: %v, but got: %v", i, testCases[i].hash.Hex(), testCases[i].header.Hash().Hex())
		}
	}
}
