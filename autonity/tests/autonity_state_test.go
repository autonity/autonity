package tests

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/params"
)

func TestInitialState(t *testing.T) {
	r := setup(t, nil)

	r.run("Test token name", func(rr *runner) {
		name, _, err := rr.autonity.Name(nil)
		require.NoError(t, err)
		require.Equal(t, "Newton", name)
	})

	r.run("Test token symbol", func(rr *runner) {
		symbol, _, err := rr.autonity.Symbol(nil)
		require.NoError(t, err)
		require.Equal(t, "NTN", symbol)
	})

	r.run("Test get min base fee", func(rr *runner) {
		mBaseFee, _, err := rr.autonity.GetMinimumBaseFee(r.operator)
		require.NoError(t, err)
		require.Equal(t, rr.genesis.Config.AutonityContractConfig.MinBaseFee, mBaseFee.Uint64())
	})

	r.run("Test get contract version", func(rr *runner) {
		version, _, err := rr.autonity.GetVersion(nil)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(1), version)
	})

	r.run("Test get max committee size", func(rr *runner) {
		cSize, _, err := rr.autonity.GetMaxCommitteeSize(nil)
		require.NoError(t, err)
		require.Equal(t, rr.genesis.Config.AutonityContractConfig.MaxCommitteeSize, cSize.Uint64())
	})

	r.run("Test get operator account", func(rr *runner) {
		operator, _, err := rr.autonity.GetOperator(nil)
		require.NoError(t, err)
		require.Equal(t, r.operator.origin, operator)
	})

	r.run("Test get validators", func(rr *runner) {
		vals, _, err := rr.autonity.GetValidators(nil)
		require.NoError(t, err)
		expectedVals := func() []common.Address {
			var validators []common.Address
			for _, v := range rr.committee.validators {
				validators = append(validators, v.NodeAddress)
			}
			return validators
		}()
		require.True(t, reflect.DeepEqual(vals, expectedVals))
	})

	r.run("Test get committee", func(rr *runner) {
		cmty, _, err := rr.autonity.GetCommittee(nil)
		require.NoError(t, err)
		expectedCmtyAddrs := func() []common.Address {
			var expected []common.Address
			for _, v := range rr.committee.validators {
				expected = append(expected, v.NodeAddress)
			}
			return expected
		}()

		cmtyAddrs := func() []common.Address {
			var addresses []common.Address
			for _, v := range cmty {
				addresses = append(addresses, v.Addr)
			}
			return addresses
		}()

		require.True(t,
			reflect.DeepEqual(
				cmtyAddrs,
				expectedCmtyAddrs,
			),
		)
	})

	r.run("Test get committee enodes", func(rr *runner) {
		committeeEnodes, _, err := rr.autonity.GetCommitteeEnodes(nil)
		require.NoError(t, err)
		expectedEnodes := func() []string {
			var enodes []string
			for _, v := range rr.committee.validators {
				enodes = append(enodes, v.Enode)
			}
			return enodes
		}()
		require.True(t, reflect.DeepEqual(committeeEnodes, expectedEnodes))
	})

	r.run("Test getValidator, balanceOf and totalSupply", func(rr *runner) {
		totalExpectedSupply := big.NewInt(0)
		for _, expectedValidator := range rr.committee.validators {
			totalExpectedSupply = new(big.Int).Add(expectedValidator.BondedStake, totalExpectedSupply)

			balance, _, err := rr.autonity.BalanceOf(nil, expectedValidator.NodeAddress)
			require.NoError(t, err)

			require.Equal(t, balance.Int64(), int64(0), "initial balance of validator is not expected")

			readValidator, _, err := rr.autonity.GetValidator(nil, expectedValidator.NodeAddress)
			require.NoError(t, err)

			require.Equal(t, readValidator.Treasury, expectedValidator.Treasury, "unexpected treasury address")
			require.Equal(t, readValidator.NodeAddress, expectedValidator.NodeAddress, "unexpected node address")
			require.Equal(t, readValidator.Enode, expectedValidator.Enode, "unexpected enode")

			require.Equal(t, rr.genesis.Config.AutonityContractConfig.DelegationRate, readValidator.CommissionRate.Uint64(), "incorrect commission rate")
			require.Equal(t, expectedValidator.BondedStake, readValidator.BondedStake, "incorrect bonded stake")
			require.Equal(t, expectedValidator.TotalSlashed.Int64(), readValidator.TotalSlashed.Int64(), "incorrect total slashed")
			require.Equal(t, expectedValidator.RegistrationBlock.Int64(), readValidator.RegistrationBlock.Int64(), "incorrect registration block")
			require.Equal(t, expectedValidator.State, readValidator.State, "incorrect state")
		}

		// add stakable vesting mint to expected supply
		totalExpectedSupply = new(big.Int).Add(totalExpectedSupply, params.DefaultStakableVestingGenesis.TotalNominal)

		totalSupply, _, err := rr.autonity.TotalSupply(nil)
		require.NoError(t, err)
		require.Equal(t, totalExpectedSupply.String(), totalSupply.String())
	})
}

func TestValidatorCommissionRate(t *testing.T) {
	r := setup(t, func(genesisConfig *core.Genesis) *core.Genesis {
		genesisConfig.Config.AutonityContractConfig.UnbondingPeriod = uint64(0)
		return genesisConfig
	})

	r.run("Test revert with unauthorized caller", func(rr *runner) {
		_, err := rr.autonity.ChangeCommissionRate(&runOptions{origin: rr.committee.validators[2].NodeAddress}, rr.committee.validators[1].NodeAddress, big.NewInt(1337))
		require.Error(t, err, "should revert with incorrect caller")
		require.Contains(t, err.Error(), "execution reverted: require caller to be validator admin account")

		_, err = rr.autonity.ChangeCommissionRate(&runOptions{origin: rr.randomAccount()}, rr.randomAccount(), big.NewInt(1337))
		require.Error(t, err, "should revert with incorrect caller")
		require.Contains(t, err.Error(), "execution reverted: validator must be registered")

		_, err = rr.autonity.ChangeCommissionRate(&runOptions{origin: rr.committee.validators[2].NodeAddress}, rr.committee.validators[2].NodeAddress, big.NewInt(13370))
		require.Error(t, err, "should revert with incorrect rate")
		require.Contains(t, err.Error(), "execution reverted: require correct commission rate")
	})

	r.run("Test change commission rate", func(rr *runner) {
		rr.keepLogs(true)
		_, err := rr.autonity.ChangeCommissionRate(&runOptions{origin: rr.committee.validators[0].NodeAddress}, rr.committee.validators[0].NodeAddress, big.NewInt(1337))
		require.NoError(t, err)

		contains := emitsEvent(rr.logs, rr.autonity.ParseCommissionRateChange, &AutonityTestCommissionRateChange{
			Validator: rr.committee.validators[0].NodeAddress,
			Rate:      big.NewInt(1337),
		})
		require.True(t, contains, "commission rate change log not emitted")

		idx, _, err := rr.autonity.GetCommissionRateChangeQueueFirst(nil)
		require.NoError(t, err)

		commissionRateChangeRequest, _, err := rr.autonity.GetCommissionRateChangeRequest(nil, idx)
		require.NoError(t, err)
		require.Equal(t, rr.committee.validators[0].NodeAddress, commissionRateChangeRequest.Validator)
		require.Equal(t, big.NewInt(1337), commissionRateChangeRequest.Rate)

		_, err = rr.autonity.ApplyNewCommissionRates(&runOptions{origin: rr.origin})
		require.NoError(t, err)

		val, _, err := rr.autonity.GetValidator(nil, rr.committee.validators[0].NodeAddress)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(1337), val.CommissionRate)
	})

	r.run("Test change commission rate only after unbonding period", func(rr *runner) {
		_, err := rr.autonity.SetUnbondingPeriod(rr.operator, big.NewInt(5))
		require.NoError(t, err)

		beforeVal, _, err := rr.autonity.GetValidator(nil, rr.committee.validators[0].NodeAddress)
		require.NoError(t, err)

		initialCommissionRate := beforeVal.CommissionRate

		_, err = rr.autonity.ChangeCommissionRate(&runOptions{origin: rr.committee.validators[0].NodeAddress}, rr.committee.validators[0].NodeAddress, big.NewInt(1337))
		require.NoError(t, err)

		_, err = rr.autonity.ApplyNewCommissionRates(&runOptions{origin: rr.origin})
		require.NoError(t, err)

		intermediateVal, _, err := rr.autonity.GetValidator(nil, rr.committee.validators[0].NodeAddress)
		require.NoError(t, err)

		require.Equal(t, initialCommissionRate, intermediateVal.CommissionRate)

		rr.waitNBlocks(5)

		_, err = rr.autonity.ApplyNewCommissionRates(&runOptions{origin: rr.origin})
		require.NoError(t, err)

		afterVal, _, err := rr.autonity.GetValidator(nil, rr.committee.validators[0].NodeAddress)
		require.NoError(t, err)
		require.Equal(t, int64(1337), afterVal.CommissionRate.Int64())
	})
}

func TestSetProtocolParameters(t *testing.T) {
	r := setup(t, nil)

	r.run("Test set min base fee by operator", func(rr *runner) {
		newBaseFee := big.NewInt(50000)

		rr.keepLogs(true)
		_, err := rr.autonity.SetMinimumBaseFee(rr.operator, newBaseFee)
		require.NoError(t, err)
		require.True(t, emitsEvent(rr.Logs(), rr.autonity.ParseMinimumBaseFeeUpdated, &AutonityTestMinimumBaseFeeUpdated{
			GasPrice: newBaseFee,
		}), "setMinimumBaseFee should emit a MinimumBaseFeeUpdated even with the correct params")

		minBaseFee, _, err := rr.autonity.GetMinimumBaseFee(nil)
		require.NoError(t, err)

		require.Equal(t, int64(50000), minBaseFee.Int64())
	})

	r.run("Test set min base fee fails by non-operator", func(rr *runner) {
		initialMinBaseFee, _, err := rr.autonity.GetMinimumBaseFee(nil)
		require.NoError(t, err)

		_, err = rr.autonity.SetMinimumBaseFee(&runOptions{origin: rr.randomAccount()}, big.NewInt(50000))
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: caller is not the operator")

		minBaseFee, _, err := rr.autonity.GetMinimumBaseFee(nil)
		require.NoError(t, err)
		require.Equal(t, initialMinBaseFee, minBaseFee)
	})

	r.run("Test set max committee size by operator", func(rr *runner) {
		// check that committee size is not 500
		committeeSize, _, err := rr.autonity.GetMaxCommitteeSize(nil)
		require.NoError(t, err)
		require.NotEqual(t, int64(500), committeeSize.Int64())

		_, err = rr.autonity.SetCommitteeSize(rr.operator, big.NewInt(500))
		require.NoError(t, err)

		// verify that committee size is now 500
		maxCommitteeSize, _, err := rr.autonity.GetMaxCommitteeSize(nil)
		require.NoError(t, err)
		require.Equal(t, int64(500), maxCommitteeSize.Int64())
	})

	r.run("Test set committee size fails by non-operator", func(rr *runner) {
		initialCommitteeSize, _, err := rr.autonity.GetMaxCommitteeSize(nil)
		require.NoError(t, err)

		_, err = rr.autonity.SetCommitteeSize(&runOptions{origin: rr.randomAccount()}, big.NewInt(500))
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: caller is not the operator")

		committeeSize, _, err := rr.autonity.GetMaxCommitteeSize(nil)
		require.NoError(t, err)
		require.Equal(t, initialCommitteeSize, committeeSize)
	})

	r.run("Test set unbonding period by operator", func(rr *runner) {
		// check that unbonding period is not 37
		unbondingPeriod, _, err := rr.autonity.GetUnbondingPeriod(nil)
		require.NoError(t, err)
		require.NotEqual(t, int64(37), unbondingPeriod.Int64())

		_, err = rr.autonity.SetUnbondingPeriod(rr.operator, big.NewInt(37))
		require.NoError(t, err)

		// verify that unbonding period is now 37
		unbondingPeriod, _, err = rr.autonity.GetUnbondingPeriod(nil)
		require.NoError(t, err)
		require.Equal(t, int64(37), unbondingPeriod.Int64())
	})

	r.run("Test set unbonding period fails by non-operator", func(rr *runner) {
		unbondingPeriod, _, err := rr.autonity.GetUnbondingPeriod(nil)
		require.NoError(t, err)
		require.NotEqual(t, int64(37), unbondingPeriod.Int64())

		_, err = rr.autonity.SetUnbondingPeriod(&runOptions{origin: rr.randomAccount()}, big.NewInt(37))
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: caller is not the operator")
	})

	r.run("Test extend epoch period by operator", func(rr *runner) {
		newPeriod := big.NewInt(307)

		epochPeriod, _, err := rr.autonity.GetEpochPeriod(nil)
		require.NoError(t, err)
		require.NotEqual(t, int64(307), epochPeriod.Int64())

		rr.keepLogs(true)

		_, err = rr.autonity.SetEpochPeriod(rr.operator, newPeriod)
		nextEpoch, _, err := rr.autonity.GetNextEpochBlock(nil)
		require.NoError(t, err)

		require.NoError(t, err)
		require.True(t, emitsEvent(rr.Logs(), rr.autonity.ParseEpochPeriodUpdated, &AutonityTestEpochPeriodUpdated{
			Period:             newPeriod,
			ToBeAppliedAtBlock: nextEpoch,
		}), "setEpochPeriod should emit EpochPeriodUpdatedEvent with correct param")

		epochPeriod, _, err = rr.autonity.GetEpochPeriod(nil)
		require.NoError(t, err)
		require.Equal(t, int64(307), epochPeriod.Int64())
	})

	r.run("Test extend epoch period fails by non-operator", func(rr *runner) {
		epochPeriod, _, err := rr.autonity.GetEpochPeriod(nil)
		require.NoError(t, err)
		require.NotEqual(t, int64(307), epochPeriod.Int64())

		_, err = rr.autonity.SetEpochPeriod(&runOptions{origin: rr.randomAccount()}, big.NewInt(307))
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: caller is not the operator")

		epochPeriod, _, err = rr.autonity.GetEpochPeriod(nil)
		require.NoError(t, err)
		require.NotEqual(t, int64(307), epochPeriod.Int64())
	})

	r.run("Test set operator account by operator", func(rr *runner) {
		newOperator := rr.randomAccount()

		operator, _, err := rr.autonity.GetOperator(nil)
		require.NoError(t, err)
		require.NotEqual(t, newOperator, operator)

		_, err = rr.autonity.SetOperatorAccount(rr.operator, newOperator)
		require.NoError(t, err)

		operator, _, err = rr.autonity.GetOperator(nil)
		require.NoError(t, err)
		require.Equal(t, newOperator, operator)
	})

	r.run("Test set operator account fails by non-operator", func(rr *runner) {
		newOperator := rr.randomAccount()

		operator, _, err := rr.autonity.GetOperator(nil)
		require.NoError(t, err)
		require.NotEqual(t, newOperator, operator)

		_, err = rr.autonity.SetOperatorAccount(&runOptions{origin: rr.randomAccount()}, newOperator)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: caller is not the operator")

		operator, _, err = rr.autonity.GetOperator(nil)
		require.NoError(t, err)
		require.NotEqual(t, newOperator, operator)
	})

	r.run("Test set treasury account by operator", func(rr *runner) {
		newTreasury := rr.randomAccount()

		treasury, _, err := rr.autonity.GetTreasuryAccount(nil)
		require.NoError(t, err)
		require.NotEqual(t, newTreasury, treasury)

		_, err = rr.autonity.SetTreasuryAccount(rr.operator, newTreasury)
		require.NoError(t, err)

		treasury, _, err = rr.autonity.GetTreasuryAccount(nil)
		require.NoError(t, err)

		require.Equal(t, newTreasury, treasury)
	})

	r.run("Test set treasury account fails by non-operator", func(rr *runner) {
		newTreasury := rr.randomAccount()

		treasury, _, err := rr.autonity.GetTreasuryAccount(nil)
		require.NoError(t, err)
		require.NotEqual(t, newTreasury, treasury)

		_, err = rr.autonity.SetTreasuryAccount(&runOptions{origin: rr.randomAccount()}, newTreasury)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: caller is not the operator")

		treasury, _, err = rr.autonity.GetTreasuryAccount(nil)
		require.NoError(t, err)
		require.NotEqual(t, newTreasury, treasury)
	})

	r.run("Test set treasury fee by operator", func(rr *runner) {
		newTreasuryFee := big.NewInt(54321)

		treasuryFee, _, err := rr.autonity.GetTreasuryFee(nil)
		require.NoError(t, err)
		require.NotEqual(t, newTreasuryFee, treasuryFee)

		_, err = rr.autonity.SetTreasuryFee(rr.operator, newTreasuryFee)
		require.NoError(t, err)

		treasuryFee, _, err = rr.autonity.GetTreasuryFee(nil)
		require.NoError(t, err)
		require.Equal(t, newTreasuryFee, treasuryFee)
	})

	r.run("Test set treasury fee fails by non-operator", func(rr *runner) {
		newTreasuryFee := big.NewInt(54321)

		treasuryFee, _, err := rr.autonity.GetTreasuryFee(nil)
		require.NoError(t, err)
		require.NotEqual(t, newTreasuryFee, treasuryFee)

		_, err = rr.autonity.SetTreasuryFee(&runOptions{origin: rr.randomAccount()}, newTreasuryFee)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: caller is not the operator")

		treasuryFee, _, err = rr.autonity.GetTreasuryFee(nil)
		require.NoError(t, err)
		require.NotEqual(t, newTreasuryFee, treasuryFee)
	})
}

func TestOnlyAccountabilityOnlyProtocol(t *testing.T) {
	r := setup(t, nil)

	r.run("Test updateValidatorAndTransferSlashedFunds can be called by accountability", func(rr *runner) {
		val, _, err := rr.autonity.GetValidator(nil, rr.committee.validators[0].NodeAddress)
		require.NoError(t, err)

		_, err = rr.autonity.UpdateValidatorAndTransferSlashedFunds(
			&runOptions{origin: rr.accountability.address},
			val,
		)
		require.NoError(t, err)
	})

	r.run("Test updateValidatorAndTransferSlashedFunds cannot be called by non-accountability", func(rr *runner) {
		val, _, err := rr.autonity.GetValidator(nil, rr.committee.validators[0].NodeAddress)
		require.NoError(t, err)

		_, err = rr.autonity.UpdateValidatorAndTransferSlashedFunds(
			&runOptions{origin: rr.randomAccount()},
			val,
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: caller is not the slashing contract")
	})

	r.run("Test finalize can be called by protocol", func(rr *runner) {
		_, err := rr.autonity.Finalize(&runOptions{origin: rr.origin})
		require.NoError(t, err)
	})

	r.run("Test finalize cannot be called by non-protocol", func(rr *runner) {
		_, err := rr.autonity.Finalize(&runOptions{origin: rr.randomAccount()})
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: function restricted to the protocol")
	})

	r.run("Test compute committee can be called by the protocol", func(rr *runner) {
		_, err := rr.autonity.ComputeCommittee(&runOptions{origin: rr.origin})
		require.NoError(t, err)
	})

	r.run("Test compute commitee cannot be called by non-protocol", func(rr *runner) {
		_, err := rr.autonity.ComputeCommittee(&runOptions{origin: rr.randomAccount()})
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: function restricted to the protocol")
	})
}

func TestERC20TokenManagement(t *testing.T) {
	r := setup(t, nil)

	r.run("Test mint Newton by operator", func(rr *runner) {
		account := rr.randomAccount()
		amount := big.NewInt(1e18)

		supplyBefore, _, err := rr.autonity.TotalSupply(nil)
		require.NoError(t, err)

		balanceBefore, _, err := rr.autonity.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, int64(0), balanceBefore.Int64())

		rr.keepLogs(true)
		_, err = rr.autonity.Mint(rr.operator, account, amount)
		require.NoError(t, err)
		require.True(t, emitsEvent(rr.Logs(), rr.autonity.ParseMintedStake, &AutonityTestMintedStake{
			Addr:   account,
			Amount: amount,
		}))

		supplyAfter, _, err := rr.autonity.TotalSupply(nil)
		require.NoError(t, err)
		require.Equal(t, new(big.Int).Add(supplyBefore, amount), supplyAfter)

		balanceAfter, _, err := rr.autonity.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, amount, balanceAfter)
	})

	r.run("Test mint Newton fails by non-operator", func(rr *runner) {
		account := rr.randomAccount()
		amount := big.NewInt(1e18)

		supplyBefore, _, err := rr.autonity.TotalSupply(nil)
		require.NoError(t, err)

		balanceBefore, _, err := rr.autonity.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, int64(0), balanceBefore.Int64())

		_, err = rr.autonity.Mint(&runOptions{origin: rr.randomAccount()}, account, amount)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: caller is not the operator")

		supplyAfter, _, err := rr.autonity.TotalSupply(nil)
		require.NoError(t, err)
		require.Equal(t, supplyBefore, supplyAfter)
	})

	r.run("Test burn Newton by operator", func(rr *runner) {
		account := rr.randomAccount()
		amount := big.NewInt(1e18)

		_, err := rr.autonity.Mint(rr.operator, account, amount)
		require.NoError(t, err)

		supplyBefore, _, err := rr.autonity.TotalSupply(nil)
		require.NoError(t, err)

		balanceBefore, _, err := rr.autonity.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, amount, balanceBefore)

		rr.keepLogs(true)
		_, err = rr.autonity.Burn(rr.operator, account, amount)
		require.NoError(t, err)
		require.True(t, emitsEvent(rr.Logs(), rr.autonity.ParseBurnedStake, &AutonityTestBurnedStake{
			Addr:   account,
			Amount: amount,
		}), "burn should emit a burn event with the correct params")

		supplyAfter, _, err := rr.autonity.TotalSupply(nil)
		require.NoError(t, err)
		require.Equal(t, new(big.Int).Sub(supplyBefore, amount), supplyAfter)

		balanceAfter, _, err := rr.autonity.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, int64(0), balanceAfter.Int64())
	})

	r.run("Test burn Newton fails by non-operator", func(rr *runner) {
		account := rr.randomAccount()
		amount := big.NewInt(1e18)

		_, err := rr.autonity.Mint(rr.operator, account, amount)
		require.NoError(t, err)

		supplyBefore, _, err := rr.autonity.TotalSupply(nil)
		require.NoError(t, err)

		balanceBefore, _, err := rr.autonity.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, amount, balanceBefore)

		_, err = rr.autonity.Burn(&runOptions{origin: rr.randomAccount()}, account, amount)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: caller is not the operator")

		supplyAfter, _, err := rr.autonity.TotalSupply(nil)
		require.NoError(t, err)
		require.Equal(t, supplyBefore, supplyAfter)

		balanceAfter, _, err := rr.autonity.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, amount, balanceAfter)
	})

	r.run("Test ERC20 token transfer", func(rr *runner) {
		account1 := rr.randomAccount()
		account2 := rr.randomAccount()

		amount := big.NewInt(1e18)

		_, err := rr.autonity.Mint(rr.operator, account1, amount)
		require.NoError(t, err)

		balanceBefore1, _, err := rr.autonity.BalanceOf(nil, account1)
		require.NoError(t, err)
		require.Equal(t, amount, balanceBefore1)

		balanceBefore2, _, err := rr.autonity.BalanceOf(nil, account2)
		require.NoError(t, err)
		require.Equal(t, int64(0), balanceBefore2.Int64())

		rr.keepLogs(true)
		_, err = rr.autonity.Transfer(&runOptions{origin: account1}, account2, amount)
		require.NoError(t, err)
		require.True(t, emitsEvent(rr.Logs(), rr.autonity.ParseTransfer, &AutonityTestTransfer{
			From:  account1,
			To:    account2,
			Value: amount,
		}), "transfer should emit a Transfer event with the correct params")

		balanceAfter1, _, err := rr.autonity.BalanceOf(nil, account1)
		require.NoError(t, err)
		require.Equal(t, int64(0), balanceAfter1.Int64())

		balanceAfter2, _, err := rr.autonity.BalanceOf(nil, account2)
		require.NoError(t, err)
		require.Equal(t, amount, balanceAfter2)
	})

	r.run("Test ERC20 token transfer fails by insufficient funds", func(rr *runner) {
		account1 := rr.randomAccount()
		account2 := rr.randomAccount()

		amount := big.NewInt(1e18)

		_, err := rr.autonity.Mint(rr.operator, account1, new(big.Int).Div(amount, big.NewInt(2)))
		require.NoError(t, err)

		balanceBefore1, _, err := rr.autonity.BalanceOf(nil, account1)
		require.NoError(t, err)
		require.Equal(t, new(big.Int).Div(amount, big.NewInt(2)), balanceBefore1)

		balanceBefore2, _, err := rr.autonity.BalanceOf(nil, account2)
		require.NoError(t, err)
		require.Equal(t, int64(0), balanceBefore2.Int64())

		_, err = rr.autonity.Transfer(&runOptions{origin: account1}, account2, amount)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: amount exceeds balance")

		balanceAfter1, _, err := rr.autonity.BalanceOf(nil, account1)
		require.NoError(t, err)
		require.Equal(t, new(big.Int).Div(amount, big.NewInt(2)), balanceAfter1)
	})

	r.run("Test ERC20 token approve", func(rr *runner) {
		account1 := rr.randomAccount()
		account2 := rr.randomAccount()

		amount := big.NewInt(1e18)
		_, err := rr.autonity.Mint(rr.operator, account1, amount)
		require.NoError(t, err)

		balanceBefore1, _, err := rr.autonity.BalanceOf(nil, account1)
		require.NoError(t, err)
		require.Equal(t, amount, balanceBefore1)

		balanceBefore2, _, err := rr.autonity.BalanceOf(nil, account2)
		require.NoError(t, err)
		require.Equal(t, int64(0), balanceBefore2.Int64())

		rr.keepLogs(true)
		_, err = rr.autonity.Approve(&runOptions{origin: account1}, account2, amount)
		require.NoError(t, err)
		require.True(t, emitsEvent(rr.Logs(), rr.autonity.ParseApproval, &AutonityTestApproval{
			Owner:   account1,
			Spender: account2,
			Value:   amount,
		}), "approve should emit an approval event with the correct params")

		allowance, _, err := rr.autonity.Allowance(nil, account1, account2)
		require.NoError(t, err)
		require.Equal(t, amount, allowance)

		_, err = rr.autonity.TransferFrom(&runOptions{origin: account2}, account1, account2, amount)
		require.NoError(t, err)

		balanceAfter1, _, err := rr.autonity.BalanceOf(nil, account1)
		require.NoError(t, err)
		require.Equal(t, int64(0), balanceAfter1.Int64())

		balanceAfter2, _, err := rr.autonity.BalanceOf(nil, account2)
		require.NoError(t, err)
		require.Equal(t, amount, balanceAfter2)
	})
}
