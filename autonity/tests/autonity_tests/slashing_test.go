package autonitytests

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/stretchr/testify/require"
)

const SLASHING_RATE_PRECISION int64 = 10_000

func TestAfterSlashingEffect(t *testing.T) {
	setup := func() *tests.Runner {
		r := tests.Setup(t, tests.SetInflationReserveZero)
		r.WaitNextEpoch()
		for _, c := range r.Committee.Validators {
			r.NoError(
				r.Autonity.Unbond(
					tests.FromSender(c.Treasury, nil),
					c.NodeAddress,
					c.SelfBondedStake,
				),
			)
			r.NoError(
				r.Autonity.Mint(
					r.Operator,
					c.Treasury,
					big.NewInt(100),
				),
			)
			r.NoError(
				r.Autonity.Bond(
					tests.FromSender(c.Treasury, nil),
					c.NodeAddress,
					big.NewInt(100),
				),
			)
		}
		r.WaitNextEpoch()
		for {
			info := r.CheckErrorAndGetData(
				r.Autonity.GetValidator(nil, r.Committee.Validators[0].NodeAddress),
			).(tests.AutonityValidator)
			if info.SelfUnbondingStake.Cmp(common.Big0) == 0 {
				break
			}
			r.WaitNextEpoch()
		}
		return r
	}

	testAccountability := func(r *tests.Runner, config tests.AccountabilityConfig) *tests.AccountabilityTest {
		_, _, contract, err := r.DeployAccountabilityTest(nil, r.Autonity.Address(), config)
		require.NoError(r.T, err)
		r.NoError(
			r.Autonity.SetAccountabilityContract(
				r.Operator,
				contract.Address(),
			),
		)
		r.NoError(
			contract.FinalizeInitialization(
				tests.FromSender(r.Autonity.Address(), nil),
				r.CheckErrorAndGetData(
					r.Autonity.GetCommittee(nil),
				).([]tests.AutonityCommitteeMember),
			),
		)
		return contract
	}

	tests.RunWithSetup("does not trigger fairness issue (unbondingStake > 0 and delegatedStake > 0)", setup, func(r *tests.Runner) {
		// fairness issue is triggered when delegatedStake or unbondingStake becomes 0 from positive due to slashing
		// it can happen due to slashing rate = 100%
		// it should not happen for slashing amount < totalStake
		config, _, err := r.Accountability.Config(nil)
		require.NoError(r.T, err)

		// modifying config so we get slashingAmount = totalStake - 1, the highest slash possible without triggering fairness issue
		const expectedBondedStake = SLASHING_RATE_PRECISION
		const expectedSlash = expectedBondedStake - 1
		config.Factors.Collusion = new(big.Int).Sub(
			big.NewInt(expectedSlash),
			config.BaseSlashingRates.Mid,
		)

		accountability := testAccountability(r, config)
		unbondFactorDenominator := big.NewInt(10_000_000)
		unbondFactorNumerator := []*big.Int{
			new(big.Int).Div(
				unbondFactorDenominator,
				big.NewInt(10), // 1/10
			),
			new(big.Int).Div(
				new(big.Int).Mul(unbondFactorDenominator, big.NewInt(9)),
				big.NewInt(10), // 9/10
			),
			new(big.Int).Div(
				unbondFactorDenominator,
				big.NewInt(100), // 1/100
			),
			new(big.Int).Div(
				new(big.Int).Mul(unbondFactorDenominator, big.NewInt(99)),
				big.NewInt(100), // 99/100
			),
			new(big.Int).Div(
				unbondFactorDenominator,
				big.NewInt(1000), // 1/1000
			),
			new(big.Int).Div(
				new(big.Int).Mul(unbondFactorDenominator, big.NewInt(999)),
				big.NewInt(1000), // 999/1000
			),
			big.NewInt(1), // 1/10_000_000
			new(big.Int).Sub(
				unbondFactorDenominator,
				big.NewInt(1), // 9_999_999/10_000_000
			),
		}

		delegator := tests.User
		// balance := r.GetBalanceOf(delegator)
		validators := make([]common.Address, 0, len(unbondFactorNumerator))
		for i, c := range r.Committee.Validators {
			if i >= len(unbondFactorNumerator) {
				break
			}
			validators = append(validators, c.NodeAddress)
		}

		validators = append(
			validators,
			generateValidators(r, len(unbondFactorNumerator)-len(validators))...,
		)

		tokenMinted := make([]*big.Int, 0, len(validators))
		for i, v := range validators {
			validatorInfo := r.CheckErrorAndGetData(
				r.Autonity.GetValidator(nil, v),
			).(tests.AutonityValidator)

			fmt.Printf("i %v\n", i)
			fmt.Printf("expectedBondedStake %v\n", expectedBondedStake)
			fmt.Printf("validatorInfo.BondedStake %v\n", validatorInfo.BondedStake)

			mint := new(big.Int).Sub(
				big.NewInt(expectedBondedStake),
				validatorInfo.BondedStake,
			)
			r.NoError(
				r.Autonity.Mint(
					r.Operator,
					delegator,
					mint,
				),
			)
			r.NoError(
				r.Autonity.Bond(
					tests.FromSender(delegator, nil),
					v,
					mint,
				),
			)
			tokenMinted = append(tokenMinted, mint)
		}

		// let bonding apply
		r.WaitNextEpoch()
		for i, v := range validators {
			tokenUnbond := new(big.Int).Div(
				new(big.Int).Mul(
					tokenMinted[i],
					unbondFactorNumerator[i],
				),
				unbondFactorDenominator,
			)
			if tokenUnbond.Cmp(common.Big0) == 0 {
				tokenUnbond = big.NewInt(1)
			}

			r.NoError(
				r.Autonity.Unbond(
					tests.FromSender(delegator, nil),
					v,
					tokenUnbond,
				),
			)
		}

		// let unbonding apply and unbondingStake create
		r.WaitNextEpoch()
		totalStakes := make([]*big.Int, 0, len(validators))
		for _, v := range validators {
			validatorInfo := r.CheckErrorAndGetData(
				r.Autonity.GetValidator(nil, v),
			).(tests.AutonityValidator)

			totalStakes = append(totalStakes, new(big.Int).Add(
				validatorInfo.BondedStake,
				new(big.Int).Add(
					validatorInfo.UnbondingStake,
					validatorInfo.SelfUnbondingStake,
				),
			))
		}
		epochPeriod := r.CheckErrorAndGetData(
			r.Autonity.GetCurrentEpochPeriod(nil),
		).(*big.Int)
		for i, v := range validators {
			slash(r, accountability, 1, v, v, epochPeriod)
			validatorInfo := r.CheckErrorAndGetData(
				r.Autonity.GetValidator(nil, v),
			).(tests.AutonityValidator)

			require.Equal(r.T, uint8(2), validatorInfo.State)
			require.True(r.T, validatorInfo.BondedStake.Cmp(common.Big0) > 0)
			require.True(r.T, validatorInfo.UnbondingStake.Cmp(common.Big0) > 0)
			slashed := new(big.Int).Sub(
				totalStakes[i],
				new(big.Int).Add(
					validatorInfo.BondedStake,
					new(big.Int).Add(
						validatorInfo.UnbondingStake,
						validatorInfo.SelfUnbondingStake,
					),
				),
			)
			require.Equal(r.T, big.NewInt(expectedSlash), slashed)
		}

	})
}

func slash(
	r *tests.Runner,
	accountability *tests.AccountabilityTest,
	epochOffenceCount int64,
	offender, reporter common.Address,
	epochPeriod *big.Int,
) *big.Int {
	config := r.CheckErrorAndGetData(
		accountability.Config(nil),
	).(tests.AccountabilityConfig)
	event := tests.AccountabilityEvent{
		Reporter:       reporter,
		Offender:       offender,
		Block:          big.NewInt(1),
		ReportingBlock: big.NewInt(2),
	}
	r.NoError(
		accountability.Slash(
			nil,
			event,
			big.NewInt(epochOffenceCount),
			epochPeriod,
		),
	)
	return ruleToRate(r, config, event.Rule)
}

func ruleToRate(r *tests.Runner, config tests.AccountabilityConfig, rule uint8) *big.Int {
	if rule == 9 { // equivocation
		return config.BaseSlashingRates.Low
	}
	if rule <= 6 {
		return config.BaseSlashingRates.Mid
	}
	if rule == 7 || rule == 8 { // invalid proposal and invalid proposer
		return config.BaseSlashingRates.High
	}
	require.Fail(r.T, "invalid rule")
	return big.NewInt(0)
}

func generateValidators(r *tests.Runner, count int) []common.Address {
	validatorAddress := make([]common.Address, 0, count)
	for count > 0 {
		validator, signature, _, _, err := tests.RandomValidator()
		require.NoError(r.T, err)
		r.NoError(
			r.Autonity.RegisterValidator(
				tests.FromSender(validator.Treasury, nil),
				validator.Enode,
				validator.OracleAddress,
				validator.ConsensusKey,
				signature,
			),
		)
		validatorAddress = append(validatorAddress, *validator.NodeAddress)
		count--
	}
	return validatorAddress
}
