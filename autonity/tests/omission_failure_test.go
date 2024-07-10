package tests

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint"
)

func TestFinalize(t *testing.T) {
	r := setup(t, nil)
	r.waitNBlocks(tendermint.DeltaBlocks)

	// make validator 0 omission faulty
	//TODO(Lorenzo) base this on epoch period
	start := r.evm.Context.BlockNumber
	for i := 0; i < 80; i++ {
		_, err := r.autonity.Finalize(&runOptions{origin: common.Address{}}, []common.Address{r.committee.validators[0].NodeAddress}, r.committee.validators[1].NodeAddress, common.Big1, new(big.Int).SetInt64(tendermint.DeltaBlocks), false)
		require.NoError(r.t, err, "finalize function error in waitNblocks", i)
		r.evm.Context.BlockNumber = new(big.Int).Add(big.NewInt(int64(i+1)), start)
		r.evm.Context.Time = new(big.Int).Add(r.evm.Context.Time, common.Big1)
	}

	_, err := r.omissionAccountability.Finalize(&runOptions{origin: r.autonity.address}, []common.Address{r.committee.validators[0].NodeAddress}, r.committee.validators[1].NodeAddress, common.Big1, false, common.Big0, new(big.Int).SetInt64(10), true)
	operator = &runOptions{origin: defaultAutonityConfig.Protocol.OperatorAccount}
	require.NoError(r.t, err)
}
