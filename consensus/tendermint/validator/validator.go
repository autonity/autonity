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

package validator

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
)

func New(addr common.Address, votingPower int64) *defaultValidator {
	return &defaultValidator{
		address: addr,
		votingPower: votingPower,
	}
}


func NewValidatorsList(votingPower int64, addrs ...common.Address) []tendermint.Validator {
	vals := make([]tendermint.Validator, len(addrs))
	for i := range addrs {
		vals[i] = New(addrs[i], votingPower)
	}

	return vals
}

func NewSet(policy tendermint.ProposerPolicy, vals ...tendermint.Validator) *defaultSet {
	return newDefaultSet(policy, vals...)
}

func ExtractValidators(extraData []byte) []common.Address {
	// get the validator addresses
	addrs := make([]common.Address, len(extraData)/common.AddressLength)
	for i := 0; i < len(addrs); i++ {
		copy(addrs[i][:], extraData[i*common.AddressLength:])
	}

	return addrs
}
