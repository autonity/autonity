package params

import (
	"math/big"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/stretchr/testify/require"
)

// can add tests to check hardcoded values or constraints for piccadilly testnet config

func TestPiccadillyGVNTNallocs(t *testing.T) {
	// `PiccadillyGVNTNallocs` is not used anywhere but the values in the list
	// are already present in `PiccadillyGenesisValidators`. Check if it's true
	bondedStakes := make(map[common.Address]*big.Int)
	for _, v := range PiccadillyGenesisValidators {
		bondedStakes[v.ValidatorAddress] = v.BondedStake
	}

	for _, v := range PiccadillyGVNTNallocs {
		stake, ok := bondedStakes[v.Address]
		require.True(t, ok)
		require.Equal(t, stake, v.Value)
	}
}
