package core

import (
	"reflect"
	"testing"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/golang/mock/gomock"
)

func TestValidatorSetSetEmpty(t *testing.T) {
	valSet := validatorSet{}
	if valSet.Set != nil {
		t.Fatalf("nil validator set expected, got %v", &valSet)
	}

	innerSet := Validators([]common.Address{{}})

	valSet.set(innerSet)
	if !reflect.DeepEqual(valSet.Set, innerSet) {
		t.Fatalf("validator set expected %v, got %v", innerSet, valSet.Set)
	}
}

func TestValidatorSetSet(t *testing.T) {
	valSet := validatorSet{}
	if valSet.Set != nil {
		t.Fatalf("nil validator set expected, got %v", &valSet)
	}

	innerSet := Validators([]common.Address{{}})

	valSet.set(innerSet)
	if !reflect.DeepEqual(valSet.Set, innerSet) {
		t.Fatalf("validator set expected %v, got %v", innerSet, valSet.Set)
	}

	innerSet = Validators([]common.Address{{}})
	valSet.set(innerSet)
	if !reflect.DeepEqual(valSet.Set, innerSet) {
		t.Fatalf("updated validator set expected %v, got %v", innerSet, valSet.Set)
	}
}

func TestValidatorSetSizeNil(t *testing.T) {
	valSet := validatorSet{}
	size := valSet.Size()
	if size != 0 {
		t.Fatalf("validator set size expected 0, got %v", size)
	}
}

func TestValidatorSetListNil(t *testing.T) {
	valSet := validatorSet{}
	list := valSet.List()
	if len(list) != 0 {
		t.Fatalf("validator set list expected 0, got %v", list)
	}
}

func TestValidatorSetGetByIndexNil(t *testing.T) {
	valSet := validatorSet{}
	val := valSet.GetByIndex(0)
	if val != nil {
		t.Fatalf("validator expected nil, got %v", val)
	}
}

func TestValidatorSetGetByAddressNil(t *testing.T) {
	valSet := validatorSet{}
	index, val := valSet.GetByAddress(common.Address{})
	if val != nil {
		t.Fatalf("validator expected nil, got %v", val)
	}

	if index != -1 {
		t.Fatalf("validator index expected nil, got %v", index)
	}
}

func TestValidatorSetGetProposerNil(t *testing.T) {
	valSet := validatorSet{}
	val := valSet.GetProposer()
	if val != nil {
		t.Fatalf("validator expected nil, got %v", val)
	}
}

func TestValidatorSetCopyNil(t *testing.T) {
	valSet := validatorSet{}
	valset := valSet.Copy()
	if valset != nil {
		t.Fatalf("validator set expected nil, got %v", valset)
	}
}

func TestValidatorSetPolicyNil(t *testing.T) {
	valSet := validatorSet{}
	policy := valSet.Policy()
	if policy != 0 {
		t.Fatalf("validator set policy expected 0, got %v", policy)
	}
}

func TestValidatorSetCalcProposerNil(t *testing.T) {
	valSet := validatorSet{}
	valSet.CalcProposer(common.Address{}, 0)

	policy := valSet.Policy()
	if policy != 0 {
		t.Fatalf("validator set policy expected 0, got %v", policy)
	}
}

func TestValidatorSetAddValidatorNil(t *testing.T) {
	valSet := validatorSet{}
	res := valSet.AddValidator(common.Address{})
	if res {
		t.Fatalf("validator add result expected false, got %v", res)
	}
}

func TestValidatorSetIsProposerNil(t *testing.T) {
	valSet := validatorSet{}
	res := valSet.IsProposer(common.Address{})
	if res {
		t.Fatalf("validator IsProposer result expected false, got %v", res)
	}
}

func TestValidatorSetRemoveValidatorNil(t *testing.T) {
	valSet := validatorSet{}
	res := valSet.RemoveValidator(common.Address{})
	if res {
		t.Fatalf("validator remove result expected false, got %v", res)
	}
}

func TestValidatorSetSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validatorSetMock := validator.NewMockSet(ctrl)

	expectedSize := 100
	validatorSetMock.EXPECT().
		Size().
		Return(expectedSize)

	valSet := validatorSet{}
	valSet.set(validatorSetMock)

	size := valSet.Size()
	if size != expectedSize {
		t.Fatalf("validator set size expected 0, got %v", size)
	}
}

func TestValidatorSetList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validatorSetMock := validator.NewMockSet(ctrl)

	expectedList := []validator.Validator{validator.New(common.Address{}), validator.New(common.Address{})}
	validatorSetMock.EXPECT().
		List().
		Return(expectedList)

	valSet := validatorSet{}
	valSet.set(validatorSetMock)

	list := valSet.List()
	if !reflect.DeepEqual(list, expectedList) {
		t.Fatalf("validator set list expected %v, got %v", expectedList, list)
	}
}

func TestValidatorSetGetByIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validatorSetMock := validator.NewMockSet(ctrl)

	expectedValidator := validator.New(common.Address{})
	validatorSetMock.EXPECT().
		GetByIndex(uint64(0)).
		Return(expectedValidator)

	valSet := validatorSet{}
	valSet.set(validatorSetMock)

	val := valSet.GetByIndex(0)
	if !reflect.DeepEqual(val, expectedValidator) {
		t.Fatalf("validator expected %v, got %v", expectedValidator, val)
	}
}

func TestValidatorSetGetByAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validatorSetMock := validator.NewMockSet(ctrl)

	expectedAddress := common.Address{}
	expectedAddress[0] = 1

	expectedIndex := 1

	expectedValidator := validator.New(expectedAddress)
	validatorSetMock.EXPECT().
		GetByAddress(expectedAddress).
		Return(expectedIndex, expectedValidator)

	valSet := validatorSet{}
	valSet.set(validatorSetMock)

	index, val := valSet.GetByAddress(expectedAddress)
	if !reflect.DeepEqual(val, expectedValidator) {
		t.Fatalf("validator expected %v, got %v", expectedValidator, val)
	}

	if index != expectedIndex {
		t.Fatalf("validator index expected %v, got %v", expectedIndex, index)
	}
}

func TestValidatorSetGetProposer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validatorSetMock := validator.NewMockSet(ctrl)

	expectedAddress := common.Address{}
	expectedAddress[0] = 1

	expectedValidator := validator.New(expectedAddress)
	validatorSetMock.EXPECT().
		GetProposer().
		Return(expectedValidator)

	valSet := validatorSet{}
	valSet.set(validatorSetMock)

	val := valSet.GetProposer()
	if !reflect.DeepEqual(val, expectedValidator) {
		t.Fatalf("validator expected %v, got %v", expectedValidator, val)
	}
}

func TestValidatorSetCopy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validatorSetMock := validator.NewMockSet(ctrl)

	expectedAddress := common.Address{}
	expectedAddress[0] = 1

	expectedValidatorSet := validator.NewSet([]common.Address{expectedAddress}, 1)

	validatorSetMock.EXPECT().
		Copy().
		Return(expectedValidatorSet)

	valSet := validatorSet{}
	valSet.set(validatorSetMock)

	valSetCopy := valSet.Copy()
	if !reflect.DeepEqual(valSetCopy, expectedValidatorSet) {
		t.Fatalf("validator set expected %v, got %v", expectedValidatorSet, valSetCopy)
	}
}

func TestValidatorSetPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validatorSetMock := validator.NewMockSet(ctrl)

	expectedPolicy := config.ProposerPolicy(1)

	validatorSetMock.EXPECT().
		Policy().
		Return(expectedPolicy)

	valSet := validatorSet{}
	valSet.set(validatorSetMock)

	policy := valSet.Policy()
	if policy != expectedPolicy {
		t.Fatalf("validator set expected policy %v, got %v", expectedPolicy, policy)
	}
}

func TestValidatorSetCalcProposer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validatorSetMock := validator.NewMockSet(ctrl)

	lastProposer := common.Address{}
	lastProposer[0] = 1
	round := uint64(1)

	validatorSetMock.EXPECT().
		CalcProposer(lastProposer, round).
		Return()

	valSet := validatorSet{}
	valSet.set(validatorSetMock)

	valSet.CalcProposer(lastProposer, round)
}

func TestValidatorSetIsProposer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validatorSetMock := validator.NewMockSet(ctrl)

	addr := common.Address{}
	addr[0] = 1

	expectedRes := true

	validatorSetMock.EXPECT().
		IsProposer(addr).
		Return(expectedRes)

	valSet := validatorSet{}
	valSet.set(validatorSetMock)

	res := valSet.IsProposer(addr)
	if res != expectedRes {
		t.Fatalf("validator set proposer result expected %v, got %v", expectedRes, res)
	}
}

func TestValidatorSetAddValidator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validatorSetMock := validator.NewMockSet(ctrl)

	addr := common.Address{}
	addr[0] = 1

	expectedRes := true

	validatorSetMock.EXPECT().
		AddValidator(addr).
		Return(expectedRes)

	valSet := validatorSet{}
	valSet.set(validatorSetMock)

	res := valSet.AddValidator(addr)
	if res != expectedRes {
		t.Fatalf("validator set add validator result expected %v, got %v", expectedRes, res)
	}
}

func TestValidatorSetRemoveValidator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validatorSetMock := validator.NewMockSet(ctrl)

	addr := common.Address{}
	addr[0] = 1

	expectedRes := true

	validatorSetMock.EXPECT().
		RemoveValidator(addr).
		Return(expectedRes)

	valSet := validatorSet{}
	valSet.set(validatorSetMock)

	res := valSet.RemoveValidator(addr)
	if res != expectedRes {
		t.Fatalf("validator set remove validator result expected %v, got %v", expectedRes, res)
	}
}

func Validators(validators []common.Address) validator.Set {
	return validator.NewSet(validators, config.ProposerPolicy(0))
}
