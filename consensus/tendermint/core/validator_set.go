package core

import (
	"sync"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
)

type validatorSet struct {
	sync.RWMutex
	tendermint.ValidatorSet
}

func (v *validatorSet) get() tendermint.ValidatorSet {
	v.RLock()
	defer v.RUnlock()
	if v.ValidatorSet == nil {
		return nil
	}
	valSet := v.ValidatorSet.Copy()
	return valSet
}

func (v *validatorSet) set(valSet tendermint.ValidatorSet) {
	v.Lock()
	v.ValidatorSet = valSet
	v.Unlock()
}

func (v *validatorSet) Size() int {
	v.RLock()
	defer v.RUnlock()
	if v.ValidatorSet == nil {
		return 0
	}
	size := v.ValidatorSet.Size()
	return size
}

func (v *validatorSet) List() []tendermint.Validator {
	v.RLock()
	defer v.RUnlock()
	if v.ValidatorSet == nil {
		return nil
	}

	list := v.ValidatorSet.List()
	return list
}

func (v *validatorSet) GetByIndex(i uint64) tendermint.Validator {
	v.RLock()
	defer v.RUnlock()
	if v.ValidatorSet == nil {
		return nil
	}
	val := v.ValidatorSet.GetByIndex(i)
	return val
}

func (v *validatorSet) GetByAddress(addr common.Address) (int, tendermint.Validator) {
	v.RLock()
	defer v.RUnlock()
	if v.ValidatorSet == nil {
		return -1, nil
	}
	i, val := v.ValidatorSet.GetByAddress(addr)
	return i, val
}

func (v *validatorSet) GetProposer() tendermint.Validator {
	v.RLock()
	defer v.RUnlock()
	if v.ValidatorSet == nil {
		return nil
	}
	val := v.ValidatorSet.GetProposer()
	return val
}

func (v *validatorSet) Copy() tendermint.ValidatorSet {
	v.RLock()
	defer v.RUnlock()
	if v.ValidatorSet == nil {
		return nil
	}
	valSet := v.ValidatorSet.Copy()
	return valSet
}

func (v *validatorSet) Policy() tendermint.ProposerPolicy {
	v.RLock()
	defer v.RUnlock()
	if v.ValidatorSet == nil {
		return 0
	}
	policy := v.ValidatorSet.Policy()
	return policy
}


func (v *validatorSet) CalcProposer(lastProposer common.Address, round uint64) {
	v.RLock()
	defer v.RUnlock()
	if v.ValidatorSet == nil {
		return
	}
	v.ValidatorSet.CalcProposer(lastProposer, round)
}

func (v *validatorSet) IsProposer(address common.Address) bool {
	v.RLock()
	defer v.RUnlock()
	if v.ValidatorSet == nil {
		return false
	}

	return v.ValidatorSet.IsProposer(address)
}

func (v *validatorSet) AddValidator(address common.Address) bool {
	v.RLock()
	defer v.RUnlock()
	if v.ValidatorSet == nil {
		return false
	}

	return v.ValidatorSet.AddValidator(address)
}

func (v *validatorSet) RemoveValidator(address common.Address) bool {
	v.RLock()
	defer v.RUnlock()
	if v.ValidatorSet == nil {
		return false
	}
	return v.ValidatorSet.RemoveValidator(address)
}