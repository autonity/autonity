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
	"github.com/clearmatics/autonity/core/types"
	"math"
	"math/big"
	"reflect"
	"sort"
	"sync"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
)

// ----------------------------------------------------------------------------

type defaultSet struct {
	validators Validators
	policy     config.ProposerPolicy

	proposer    Validator
	validatorMu sync.RWMutex
	selector    ProposalSelector
}

func NewSet(committee types.Committee, policy config.ProposerPolicy) *defaultSet {
	return newDefaultSet(makeValidators(committee), policy)
}

func newDefaultSet(validators Validators, policy config.ProposerPolicy) *defaultSet {
	valSet := &defaultSet{}

	valSet.policy = policy
	valSet.validators = validators

	// sort validator
	sort.Sort(valSet.validators)
	// init proposer
	if valSet.Size() > 0 {
		valSet.proposer = valSet.GetByIndex(0)
	}

	switch policy {
	case config.Sticky:
		valSet.selector = stickyProposer
	case config.RoundRobin:
		valSet.selector = roundRobinProposer
	default:
		valSet.selector = roundRobinProposer
	}

	return valSet
}

func makeValidators(committee types.Committee) []Validator {
	validators := make([]Validator, len(committee))
	for i, c := range committee {
		validators[i] = Validator(c)
	}
	return validators
}

func copyValidators(validators []Validator) []Validator {
	validatorsCopy := make([]Validator, len(validators))
	for i, val := range validators {
		validatorsCopy[i] = New(val.GetAddress(), val.GetVotingPower())
	}

	return validatorsCopy
}

func (valSet *defaultSet) Size() int {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	return len(valSet.validators)
}

func (valSet *defaultSet) List() []Validator {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()

	return copyValidators(valSet.validators)
}

func (valSet *defaultSet) GetByIndex(i uint64) Validator {
	if i < uint64(valSet.Size()) {
		valSet.validatorMu.RLock()
		defer valSet.validatorMu.RUnlock()

		return New(valSet.validators[i].GetAddress(), valSet.validators[i].GetVotingPower())
	}

	return nil
}

func (valSet *defaultSet) GetByAddress(addr common.Address) (int, Validator) {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()

	for i, val := range valSet.validators {
		if addr == val.GetAddress() {
			return i, val
		}
	}
	return -1, nil
}

func (valSet *defaultSet) GetProposer() Validator {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()

	return valSet.getProposer()
}

func (valSet *defaultSet) getProposer() Validator {
	return New(valSet.proposer.GetAddress(), valSet.proposer.GetVotingPower())
}

func (valSet *defaultSet) IsProposer(address common.Address) bool {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()

	_, val := valSet.GetByAddress(address)
	return reflect.DeepEqual(valSet.getProposer(), val)
}

func (valSet *defaultSet) CalcProposer(lastProposer common.Address, round uint64) {
	proposer := valSet.selector(valSet, lastProposer, round)

	valSet.validatorMu.Lock()
	valSet.proposer = proposer
	valSet.validatorMu.Unlock()
}

func (valSet *defaultSet) AddValidator(address common.Address) bool {
	valSet.validatorMu.Lock()
	defer valSet.validatorMu.Unlock()

	for _, v := range valSet.validators {
		if v.GetAddress() == address {
			return false
		}
	}

	valSet.validators = append(valSet.validators, New(address, new(big.Int).SetUint64(1)))
	// TODO: we may not need to re-sort it again
	sort.Sort(valSet.validators)
	return true
}

func (valSet *defaultSet) RemoveValidator(address common.Address) bool {
	valSet.validatorMu.Lock()
	defer valSet.validatorMu.Unlock()

	for i, v := range valSet.validators {
		if v.GetAddress() == address {
			valSet.validators = append(valSet.validators[:i], valSet.validators[i+1:]...)
			return true
		}
	}
	return false
}

func (valSet *defaultSet) Copy() Set {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	return newDefaultSet(copyValidators(valSet.validators), valSet.policy)
}

func (valSet *defaultSet) F() int { return int(math.Ceil(float64(valSet.Size())/3)) - 1 }

func (valSet *defaultSet) Quorum() int { return int(math.Ceil((2 * float64(valSet.Size())) / 3.)) }

func (valSet *defaultSet) Policy() config.ProposerPolicy { return valSet.policy }
