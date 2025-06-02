package tests

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlashingAccessControl(t *testing.T) {
	t.Run("slasher contract itself - anyone can call but no changes to the AC contract", func(t *testing.T) {
		r := Setup(t, nil)
		slasher := r.slasherContract()

		victim := r.Committee.Validators[1]
		blockNumber := new(big.Int).Set(r.Evm.Context.BlockNumber)

		_, _, err := slasher.Jail(r.Operator, victim, blockNumber, new(big.Int).SetUint64(1000), jailed)
		require.NoError(t, err)

		require.Equal(t, uint64(0), getValidator(r, victim.NodeAddress).JailReleaseBlock.Uint64())
		require.Equal(t, active, getValidator(r, victim.NodeAddress).State)

		_, _, err = slasher.Jailbound(r.Operator, victim, jailbound)
		require.NoError(t, err)
		require.Equal(t, active, getValidator(r, victim.NodeAddress).State)

		_, _, _, err = slasher.Slash(r.Operator, victim, new(big.Int).SetUint64(2000)) // 20%
		require.NoError(t, err)
		_, _, _, _, err = slasher.SlashAndJail(r.Operator, victim, new(big.Int).SetUint64(2000), blockNumber, new(big.Int).SetUint64(1000), jailed, jailbound)
		require.NoError(t, err)

		require.Equal(t, active, getValidator(r, victim.NodeAddress).State)
		require.Equal(t, uint64(0), getValidator(r, victim.NodeAddress).TotalSlashed.Uint64())

		// autonity can call, but still no state change if directly calling the slasher
		_, _, err = slasher.Jail(FromAutonity, victim, blockNumber, new(big.Int).SetUint64(1000), jailed)
		require.NoError(t, err)
		_, _, err = slasher.Jailbound(FromAutonity, victim, jailbound)
		require.NoError(t, err)
		_, _, _, err = slasher.Slash(FromAutonity, victim, new(big.Int).SetUint64(2000))
		require.NoError(t, err)
		slashedVal, _, _, _, err := slasher.SlashAndJail(FromAutonity, victim, new(big.Int).SetUint64(2000), blockNumber, new(big.Int).SetUint64(1000), jailed, jailbound)
		require.NoError(t, err)

		// returned validator has the changes
		require.Equal(t, jailed, slashedVal.State)
		require.Less(t, uint64(0), slashedVal.TotalSlashed.Uint64())

		require.Equal(t, active, getValidator(r, victim.NodeAddress).State)
		require.Equal(t, uint64(0), getValidator(r, victim.NodeAddress).TotalSlashed.Uint64())
	})
	t.Run("proxy slashing function in AC", func(t *testing.T) {
		r := Setup(t, nil)

		victim := r.Committee.Validators[1].NodeAddress

		// only accountability contracts can call AC proxy slashing functions
		_, err := r.Autonity.Jail(r.Operator, victim, new(big.Int), jailed)
		t.Log(err)
		require.Error(t, err)
		_, err = r.Autonity.Jailbound(r.Operator, victim, jailbound)
		t.Log(err)
		require.Error(t, err)
		_, err = r.Autonity.Slash(r.Operator, victim, new(big.Int))
		t.Log(err)
		require.Error(t, err)
		_, err = r.Autonity.SlashAndJail(r.Operator, victim, new(big.Int), new(big.Int), jailed, jailbound)
		t.Log(err)
		require.Error(t, err)

		_, err = r.Autonity.Jail(FromSender(r.OmissionAccountability.address, new(big.Int)), victim, new(big.Int), jailed)
		require.NoError(t, err)
		_, err = r.Autonity.Jailbound(FromSender(r.Accountability.address, new(big.Int)), victim, jailbound)
		require.NoError(t, err)
		_, err = r.Autonity.Slash(FromSender(r.OmissionAccountability.address, new(big.Int)), victim, new(big.Int))
		require.NoError(t, err)
		_, err = r.Autonity.SlashAndJail(FromSender(r.Accountability.address, new(big.Int)), victim, new(big.Int), new(big.Int), jailed, jailbound)
		require.NoError(t, err)
	})
}

func TestSlasherFunctions(t *testing.T) {
	t.Run("jailing", func(t *testing.T) {
		r := Setup(t, nil)
		victim := r.Committee.Validators[1].NodeAddress
		fromAccountability := FromSender(r.Accountability.address, new(big.Int))

		jailtime := new(big.Int).SetUint64(500)
		_, err := r.Autonity.Jail(fromAccountability, victim, jailtime, jailed)
		require.NoError(t, err)

		val, _, err := r.Autonity.GetValidator(nil, victim)
		require.NoError(t, err)
		require.Equal(t, jailed, val.State)
		require.Equal(t, jailtime.Uint64()+1, val.JailReleaseBlock.Uint64())
		require.Equal(t, uint64(0), val.TotalSlashed.Uint64())
	})
	t.Run("jailbounding", func(t *testing.T) {
		r := Setup(t, nil)
		victim := r.Committee.Validators[1].NodeAddress
		fromAccountability := FromSender(r.Accountability.address, new(big.Int))

		_, err := r.Autonity.Jailbound(fromAccountability, victim, jailbound)
		require.NoError(t, err)

		val, _, err := r.Autonity.GetValidator(nil, victim)
		require.NoError(t, err)
		require.Equal(t, jailbound, val.State)
		require.Equal(t, uint64(0), val.JailReleaseBlock.Uint64())
		require.Equal(t, uint64(0), val.TotalSlashed.Uint64())
	})
	t.Run("slashing", func(t *testing.T) {
		r := Setup(t, nil)
		victim := r.Committee.Validators[1].NodeAddress
		fromAccountability := FromSender(r.Accountability.address, new(big.Int))

		valBeforeSlash, _, err := r.Autonity.GetValidator(nil, victim)
		require.NoError(t, err)
		availableFunds := valBeforeSlash.BondedStake.Uint64() + valBeforeSlash.UnbondingStake.Uint64() + valBeforeSlash.SelfUnbondingStake.Uint64()
		t.Logf("available funds: %d", availableFunds)

		halfRate := slashingRateScaleFactor(r).Uint64() / 2
		_, err = r.Autonity.Slash(fromAccountability, victim, new(big.Int).SetUint64(halfRate))
		require.NoError(t, err)

		val, _, err := r.Autonity.GetValidator(nil, victim)
		require.NoError(t, err)
		require.Equal(t, active, val.State)
		require.Equal(t, uint64(0), val.JailReleaseBlock.Uint64())
		require.Equal(t, availableFunds/2, val.TotalSlashed.Uint64())
	})
	t.Run("slashing and jailing", func(t *testing.T) {
		r := Setup(t, nil)
		victim := r.Committee.Validators[1].NodeAddress
		fromAccountability := FromSender(r.Accountability.address, new(big.Int))

		valBeforeSlash, _, err := r.Autonity.GetValidator(nil, victim)
		require.NoError(t, err)
		availableFunds := valBeforeSlash.BondedStake.Uint64() + valBeforeSlash.UnbondingStake.Uint64() + valBeforeSlash.SelfUnbondingStake.Uint64()
		t.Logf("available funds: %d", availableFunds)

		halfRate := slashingRateScaleFactor(r).Uint64() / 2
		jailtime := new(big.Int).SetUint64(1000)
		_, err = r.Autonity.SlashAndJail(fromAccountability, victim, new(big.Int).SetUint64(halfRate), jailtime, jailedForInactivity, jailboundForInactivity)
		require.NoError(t, err)

		val, _, err := r.Autonity.GetValidator(nil, victim)
		require.NoError(t, err)
		require.Equal(t, jailedForInactivity, val.State)
		require.Equal(t, jailtime.Uint64()+1, val.JailReleaseBlock.Uint64())
		require.Equal(t, availableFunds/2, val.TotalSlashed.Uint64())

		jailtime = new(big.Int).SetUint64(9999)
		_, err = r.Autonity.SlashAndJail(fromAccountability, victim, slashingRateScaleFactor(r), jailtime, jailedForInactivity, jailboundForInactivity)
		require.NoError(t, err)

		val, _, err = r.Autonity.GetValidator(nil, victim)
		require.NoError(t, err)
		require.Equal(t, jailboundForInactivity, val.State)
		require.Equal(t, uint64(0), val.JailReleaseBlock.Uint64())
		require.Equal(t, availableFunds, val.TotalSlashed.Uint64())
	})
}
