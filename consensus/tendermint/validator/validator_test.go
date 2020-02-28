package validator

import (
	"math/big"
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/clearmatics/autonity/common"
)

func TestNewValidator(t *testing.T) {
	addr := common.Address{1}
	value := big.NewInt(10)
	v := New(addr, value)

	if got := v.GetAddress(); got != addr {
		t.Fatalf("got address %v, expected %v", got, addr)
	}
	if gotValue := v.GetVotingPower(); gotValue.Cmp(value) != 0 {
		t.Fatalf("got value %v, expected %v", gotValue, value)
	}
}

func TestSortValidator(t *testing.T) {
	const n = 10
	vals := Validators(make([]Validator, n))
	values := make([]*big.Int, n)
	addrs := make([]common.Address, n)

	sortedVals := Validators(make([]Validator, n))

	for i := 0; i < n; i++ {
		addrs[i] = common.Address{byte(i) + 1}
		values[i] = big.NewInt(int64(i) + 1)

		vals[i] = New(addrs[i], values[i])
		sortedVals[i] = vals[i]
	}

	rand.Shuffle(len(vals), func(i, j int) { vals[i], vals[j] = vals[j], vals[i] })

	sort.Stable(vals)

	if !reflect.DeepEqual(vals, sortedVals) {
		t.Fatalf("got unsorted validators: %v", vals)
	}
}
