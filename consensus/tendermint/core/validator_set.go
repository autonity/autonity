package core

import (
	"sync"

	"github.com/clearmatics/autonity/consensus/tendermint"
	"github.com/clearmatics/autonity/common"
)

type validatorSet struct {
	sync.RWMutex
	tendermint.ValidatorSet
}

func newValidatorSet(valSet tendermint.ValidatorSet) *validatorSet {
	return &validatorSet{ValidatorSet: valSet}
}

func (v *validatorSet) get() tendermint.ValidatorSet {
	v.RLock()
	valSet := v.ValidatorSet.Copy()
	v.RUnlock()
	return valSet
}

func (v *validatorSet) set(valSet tendermint.ValidatorSet) tendermint.ValidatorSet {
	v.Lock()
	v.ValidatorSet = valSet
	v.Unlock()
	return valSet
}

func (v *validatorSet) Size() int {
	v.RLock()
	size :=  v.ValidatorSet.Size()
	v.RUnlock()
	return size
}

func (v *validatorSet) List() []tendermint.Validator {
	v.RLock()
	list := v.ValidatorSet.List()
	v.RUnlock()
	return list
}

func (v *validatorSet) GetByIndex(i uint64) tendermint.Validator {
	v.RLock()
	val := v.ValidatorSet.GetByIndex(i)
	v.RUnlock()
	return val
}

func (v *validatorSet) GetByAddress(addr common.Address) (int, tendermint.Validator) {
	v.RLock()
	i, val := v.ValidatorSet.GetByAddress(addr)
	v.RUnlock()
	return i, val
}

func (v *validatorSet) GetProposer() tendermint.Validator {
	v.RLock()
	val := v.ValidatorSet.GetProposer()
	v.RUnlock()
	return val
}

func (v *validatorSet) Copy() tendermint.ValidatorSet {
	v.RLock()
	valSet := v.ValidatorSet.Copy()
	v.RUnlock()
	return valSet
}

func (v *validatorSet) Policy() tendermint.ProposerPolicy {
	v.RLock()
	policy := v.ValidatorSet.Policy()
	v.RUnlock()
	return policy
}

