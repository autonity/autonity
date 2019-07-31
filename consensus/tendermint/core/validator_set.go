package core

import (
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"sync"

	"github.com/clearmatics/autonity/common"
)

type validatorSet struct {
	sync.RWMutex
	validator.Set
}

func (v *validatorSet) set(valSet validator.Set) {
	v.Lock()
	v.Set = valSet
	v.Unlock()
}

func (v *validatorSet) Size() int {
	v.RLock()
	defer v.RUnlock()
	if v.Set == nil {
		return 0
	}
	size := v.Set.Size()
	return size
}

func (v *validatorSet) List() []validator.Validator {
	v.RLock()
	defer v.RUnlock()
	if v.Set == nil {
		return nil
	}

	list := v.Set.List()
	return list
}

func (v *validatorSet) GetByIndex(i uint64) validator.Validator {
	v.RLock()
	defer v.RUnlock()
	if v.Set == nil {
		return nil
	}
	val := v.Set.GetByIndex(i)
	return val
}

func (v *validatorSet) GetByAddress(addr common.Address) (int, validator.Validator) {
	v.RLock()
	defer v.RUnlock()
	if v.Set == nil {
		return -1, nil
	}
	i, val := v.Set.GetByAddress(addr)
	return i, val
}

func (v *validatorSet) GetProposer() validator.Validator {
	v.RLock()
	defer v.RUnlock()
	if v.Set == nil {
		return nil
	}
	val := v.Set.GetProposer()
	return val
}

func (v *validatorSet) Copy() validator.Set {
	v.RLock()
	defer v.RUnlock()
	if v.Set == nil {
		return nil
	}
	valSet := v.Set.Copy()
	return valSet
}

func (v *validatorSet) Policy() config.ProposerPolicy {
	v.RLock()
	defer v.RUnlock()
	if v.Set == nil {
		return 0
	}
	policy := v.Set.Policy()
	return policy
}

func (v *validatorSet) CalcProposer(lastProposer common.Address, round uint64) {
	v.RLock()
	defer v.RUnlock()
	if v.Set == nil {
		return
	}
	v.Set.CalcProposer(lastProposer, round)
}

func (v *validatorSet) IsProposer(address common.Address) bool {
	v.RLock()
	defer v.RUnlock()
	if v.Set == nil {
		return false
	}

	return v.Set.IsProposer(address)
}

func (v *validatorSet) AddValidator(address common.Address) bool {
	v.RLock()
	defer v.RUnlock()
	if v.Set == nil {
		return false
	}

	return v.Set.AddValidator(address)
}

func (v *validatorSet) RemoveValidator(address common.Address) bool {
	v.RLock()
	defer v.RUnlock()
	if v.Set == nil {
		return false
	}
	return v.Set.RemoveValidator(address)
}
