package asm

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
)

var INT256_MAX, _ = new(big.Int).SetString("57896044618658097711785492504343953926634992332820282019728792003956564819967", 10)
var ORACLE_SCALE_FACTOR = new(big.Int).SetUint64(10000000)
var INVALID_BASKET_ERROR = "0x4ff799c5"

type ACUTestData struct {
	acu *tests.ACU
}

func TestBasicACU(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("Test ACU constructor defaults", setup, func(r *tests.Runner) {
		round, _, err := r.Acu.Round(nil)
		require.NoError(t, err)

		require.Equal(t, uint64(0), round.Uint64())
		scaleFactor, _, err := r.Acu.ScaleFactor(nil)
		require.NoError(t, err)
		require.Equal(t, ORACLE_SCALE_FACTOR, scaleFactor)
	})

	tests.RunWithSetup("Test ACU constructor errors", setup, func(r *tests.Runner) {
		symbols := []string{"FOO"}
		var quantities []*big.Int
		scale := big.NewInt(1)
		autonity := common.Address{}
		operator := common.Address{}
		oracle := common.Address{}

		err := deployACUReverts(r, symbols, quantities, scale, autonity, operator, oracle)
		require.Error(t, err)
	})

}

func deployACUReverts(
	r *tests.Runner,
	symbols_ []string,
	quantities_ []*big.Int,
	scale_ *big.Int,
	autonity common.Address,
	operator common.Address,
	oracle common.Address,
) error {
	args, err := generated.ACUAbi.Pack("", symbols_, quantities_, scale_, autonity, operator, oracle)
	require.NoError(r.T, err)
	data := append(generated.ACUBytecode, args...)
	gas := uint64(math.MaxUint64)
	r.Evm.Origin = params.DeployerAddress
	out, _, _, err := r.Evm.Create(vm.AccountRef(r.Evm.Origin), data, gas, big.NewInt(0))
	if err != nil {
		reason, _ := abi.UnpackRevert(out)
		return fmt.Errorf("%w: %s", err, reason)
	}
	return nil
}
