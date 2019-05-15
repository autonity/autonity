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
	"math"
	"reflect"
	"sort"
	"sync"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
)

type defaultValidator struct {
	address common.Address
}

func (val *defaultValidator) Address() common.Address {
	return val.address
}

func (val *defaultValidator) String() string {
	return val.Address().String()
}

// ----------------------------------------------------------------------------

type defaultSet struct {
	validators tendermint.Validators
	policy     tendermint.ProposerPolicy

	proposer    tendermint.Validator
	validatorMu sync.RWMutex
	selector    tendermint.ProposalSelector
}

func newDefaultSet(addrs []common.Address, policy tendermint.ProposerPolicy) *defaultSet {
	valSet := &defaultSet{}

	valSet.policy = policy
	// init validators
	valSet.validators = make([]tendermint.Validator, len(addrs))
	for i, addr := range addrs {
		valSet.validators[i] = New(addr)
	}
	// sort validator
	sort.Sort(valSet.validators)
	// init proposer
	if valSet.Size() > 0 {
		valSet.proposer = valSet.GetByIndex(0)
	}
	valSet.selector = roundRobinProposer
	if policy == tendermint.Sticky {
		valSet.selector = stickyProposer
	}

	return valSet
}

func (valSet *defaultSet) Size() int {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	return len(valSet.validators)
}

func (valSet *defaultSet) List() []tendermint.Validator {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	return valSet.validators
}

func (valSet *defaultSet) GetByIndex(i uint64) tendermint.Validator {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	if i < uint64(valSet.Size()) {
		return valSet.validators[i]
	}
	return nil
}

func (valSet *defaultSet) GetByAddress(addr common.Address) (int, tendermint.Validator) {
	for i, val := range valSet.List() {
		if addr == val.Address() {
			return i, val
		}
	}
	return -1, nil
}

func (valSet *defaultSet) GetProposer() tendermint.Validator {
	return valSet.proposer
}

func (valSet *defaultSet) IsProposer(address common.Address) bool {
	_, val := valSet.GetByAddress(address)
	return reflect.DeepEqual(valSet.GetProposer(), val)
}

func (valSet *defaultSet) CalcProposer(lastProposer common.Address, round uint64) {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	valSet.proposer = valSet.selector(valSet, lastProposer, round)
}

func calcSeed(valSet tendermint.ValidatorSet, proposer common.Address, round uint64) uint64 {
	offset := 0
	if idx, val := valSet.GetByAddress(proposer); val != nil {
		offset = idx
	}
	return uint64(offset) + round
}

func emptyAddress(addr common.Address) bool {
	return addr == common.Address{}
}

func roundRobinProposer(valSet tendermint.ValidatorSet, proposer common.Address, round uint64) tendermint.Validator {
	if valSet.Size() == 0 {
		return nil
	}
	var seed uint64
	if emptyAddress(proposer) {
		seed = round
	} else {
		seed = calcSeed(valSet, proposer, round) + 1
	}
	pick := seed % uint64(valSet.Size())
	return valSet.GetByIndex(pick)
}

func stickyProposer(valSet tendermint.ValidatorSet, proposer common.Address, round uint64) tendermint.Validator {
	if valSet.Size() == 0 {
		return nil
	}
	var seed uint64
	if emptyAddress(proposer) {
		seed = round
	} else {
		seed = calcSeed(valSet, proposer, round)
	}
	pick := seed % uint64(valSet.Size())
	return valSet.GetByIndex(pick)
}

func (valSet *defaultSet) AddValidator(address common.Address) bool {
	valSet.validatorMu.Lock()
	defer valSet.validatorMu.Unlock()
	for _, v := range valSet.validators {
		if v.Address() == address {
			return false
		}
	}
	valSet.validators = append(valSet.validators, New(address))
	// TODO: we may not need to re-sort it again
	// sort validator
	sort.Sort(valSet.validators)
	return true
}

func (valSet *defaultSet) RemoveValidator(address common.Address) bool {
	valSet.validatorMu.Lock()
	defer valSet.validatorMu.Unlock()

	for i, v := range valSet.validators {
		if v.Address() == address {
			valSet.validators = append(valSet.validators[:i], valSet.validators[i+1:]...)
			return true
		}
	}
	return false
}

func (valSet *defaultSet) Copy() tendermint.ValidatorSet {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()

	addresses := make([]common.Address, 0, len(valSet.validators))
	for _, v := range valSet.validators {
		addresses = append(addresses, v.Address())
	}
	return NewSet(addresses, valSet.policy)
}

func (valSet *defaultSet) F() int { return int(math.Ceil(float64(valSet.Size())/3)) - 1 }

func (valSet *defaultSet) Policy() tendermint.ProposerPolicy { return valSet.policy }
