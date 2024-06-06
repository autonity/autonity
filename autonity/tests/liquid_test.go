package tests

import (
	"math/big"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/params"
)

var (
	staker1 = common.HexToAddress("0x1000000000000000000000000000000000000000")
	staker2 = common.HexToAddress("0x2000000000000000000000000000000000000000")
	staker3 = common.HexToAddress("0x3000000000000000000000000000000000000000")
)

func TestClaimRewards(t *testing.T) {
	// Test 1 validator 1 staker
	r := Setup(t, nil)
	// Mint Newton to some few accounts
	r.Autonity.Mint(Operator, staker1, params.Ntn10000)
	r.Autonity.Mint(Operator, staker2, params.Ntn10000)
	r.Autonity.Mint(Operator, staker3, params.Ntn10000)
	r.Autonity.Bond(&runOptions{origin: staker1}, r.Committee.Validators[0].NodeAddress, params.Ntn10000)
	r.Autonity.Bond(&runOptions{origin: staker2}, r.Committee.Validators[1].NodeAddress, params.Ntn10000)
	r.Autonity.Bond(&runOptions{origin: staker3}, r.Committee.Validators[1].NodeAddress, new(big.Int).Mul(common.Big2, params.Ntn10000))

	// create liquid staking contract per validator
	r.WaitNextEpoch()
	// .. test here claiming rewards, checking if NTN/ATN reward is coherent and accurate.
	// transactions fees can be simulated be sending atns directly to the autonity contract account.
	// todo: Think about in base.go to assign at each epoch the current list of validators / committee
	// in r.validators with the liquid stake contract bindings already prepared so that's easy to manipulate
	// or maybe just create some helpers for it.
}
