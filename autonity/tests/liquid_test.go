package tests

import (
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
	r := setup(t, nil)
	// Mint Newton to some few accounts
	r.autonity.Mint(operator, staker1, params.Ntn10000)
	r.autonity.Mint(operator, staker2, params.Ntn10000)
	r.autonity.Mint(operator, staker3, params.Ntn10000)
	//r.autonity.Bond(&runOptions{origin: staker1}, params.TestAutonityContractConfig.)
}
