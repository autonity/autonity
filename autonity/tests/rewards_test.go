package tests

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	tenTo18 = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	hundred = new(big.Int).SetUint64(100)
)

type changeTracker struct {
	value          *big.Int
	expectedChange *big.Int
}

func (ct *changeTracker) Validate(newValue *big.Int) error {
	actualChange := new(big.Int).Sub(newValue, ct.value)
	if actualChange.Cmp(ct.expectedChange) != 0 {
		return fmt.Errorf("expected change to be %s, got %s (value %s, newValue %s)",
			ct.expectedChange.String(), actualChange.String(), ct.value.String(), newValue.String())
	}
	return nil
}

func (ct *changeTracker) update(value *big.Int) {
	ct.expectedChange = add(ct.expectedChange, value)
}

func (ct *changeTracker) String() string {
	return fmt.Sprintf(
		"value: %s expectedChange: %s",
		ct.value.String(),
		ct.expectedChange.String(),
	)
}

func newChangeTracker(initialValue *big.Int) *changeTracker {
	return &changeTracker{
		value:          initialValue,
		expectedChange: set(common.Big0),
	}
}

type validatorTracker struct {
	// data
	nodeAddress     common.Address
	treasuryAddress common.Address
	oracleAddress   common.Address
	commission      *big.Int

	// trackers

	// autobonding increases stake
	power          *changeTracker
	selfPower      *changeTracker
	delegatedPower *changeTracker
	// ATN tracking
	nodeAddressAtn     *changeTracker // tips go here
	treasuryAddressAtn *changeTracker // normal atn reward flow
	oracleAddressAtn   *changeTracker // nothing here
	// all NTN is autobonded, so no change should record here
	nodeAddressNtn     *changeTracker
	treasuryAddressNtn *changeTracker
	oracleAddressNtn   *changeTracker
}

func newValidatorTracker(r *Runner, nodeAddress common.Address) *validatorTracker {
	val := validator(r, nodeAddress)
	return &validatorTracker{
		nodeAddress:     val.NodeAddress,
		treasuryAddress: val.Treasury,
		oracleAddress:   val.OracleAddress,
		commission:      val.CommissionRate,

		power:          newChangeTracker(val.BondedStake),
		selfPower:      newChangeTracker(val.SelfBondedStake),
		delegatedPower: newChangeTracker(sub(val.BondedStake, val.SelfBondedStake)),

		nodeAddressAtn:     newChangeTracker(r.GetBalanceOf(val.NodeAddress)),
		treasuryAddressAtn: newChangeTracker(r.GetBalanceOf(val.Treasury)),
		oracleAddressAtn:   newChangeTracker(r.GetBalanceOf(val.OracleAddress)),

		nodeAddressNtn:     newChangeTracker(r.GetNewtonBalanceOf(val.NodeAddress)),
		treasuryAddressNtn: newChangeTracker(r.GetNewtonBalanceOf(val.Treasury)),
		oracleAddressNtn:   newChangeTracker(r.GetNewtonBalanceOf(val.OracleAddress)),
	}
}

func (vt *validatorTracker) Validate(r *Runner) error {
	val := validator(r, vt.nodeAddress) // get up-to-date validator struct

	// power
	powerErr := vt.power.Validate(val.BondedStake)
	if powerErr != nil {
		powerErr = fmt.Errorf("validator power error: %w", powerErr)
	}
	selfPowerErr := vt.selfPower.Validate(val.SelfBondedStake)
	if selfPowerErr != nil {
		selfPowerErr = fmt.Errorf("validator selfPower error: %w", selfPowerErr)
	}
	delegatedPowerErr := vt.delegatedPower.Validate(sub(val.BondedStake, val.SelfBondedStake))
	if delegatedPowerErr != nil {
		delegatedPowerErr = fmt.Errorf("validator delegated power error: %w", delegatedPowerErr)
	}

	// ATN tracking
	nodeAddressAtnErr := vt.nodeAddressAtn.Validate(r.GetBalanceOf(val.NodeAddress))
	if nodeAddressAtnErr != nil {
		nodeAddressAtnErr = fmt.Errorf("validator node address ATN error: %w", nodeAddressAtnErr)
	}
	treasuryAddressAtnErr := vt.treasuryAddressAtn.Validate(r.GetBalanceOf(val.Treasury))
	if treasuryAddressAtnErr != nil {
		treasuryAddressAtnErr = fmt.Errorf("validator treasury address ATN error: %w", treasuryAddressAtnErr)
	}
	oracleAddressAtnErr := vt.oracleAddressAtn.Validate(r.GetBalanceOf(val.OracleAddress))
	if oracleAddressAtnErr != nil {
		oracleAddressAtnErr = fmt.Errorf("validator oracle address ATN error: %w", oracleAddressAtnErr)
	}

	// NTN tracking
	nodeAddressNtnErr := vt.nodeAddressNtn.Validate(r.GetNewtonBalanceOf(val.NodeAddress))
	if nodeAddressNtnErr != nil {
		nodeAddressNtnErr = fmt.Errorf("validator node address NTN error: %w", nodeAddressNtnErr)
	}
	treasuryAddressNtnErr := vt.treasuryAddressNtn.Validate(r.GetNewtonBalanceOf(val.Treasury))
	if treasuryAddressNtnErr != nil {
		treasuryAddressNtnErr = fmt.Errorf("validator treasury address NTN error: %w", treasuryAddressNtnErr)
	}
	oracleAddressNtnErr := vt.oracleAddressNtn.Validate(r.GetNewtonBalanceOf(val.OracleAddress))
	if oracleAddressNtnErr != nil {
		oracleAddressNtnErr = fmt.Errorf("validator oracle address NTN error: %w", oracleAddressNtnErr)
	}

	return errors.Join(
		powerErr, selfPowerErr, delegatedPowerErr,
		nodeAddressAtnErr, treasuryAddressAtnErr, oracleAddressAtnErr,
		nodeAddressNtnErr, treasuryAddressNtnErr, oracleAddressNtnErr,
	)
}

func (vt *validatorTracker) String() string {
	data := [][]string{
		{"validator", vt.nodeAddress.String()},
		{"power", vt.power.String()},
		{"self power", vt.selfPower.String()},
		{"delegated power", vt.delegatedPower.String()},
		{"commission", vt.commission.String()},
		{"node address ATN", vt.nodeAddressAtn.String()},
		{"treasury ATN", vt.treasuryAddressAtn.String()},
		{"oracle ATN", vt.oracleAddressAtn.String()},
		{"node address NTN", vt.nodeAddressNtn.String()},
		{"treasury NTN", vt.treasuryAddressNtn.String()},
		{"oracle NTN", vt.oracleAddressNtn.String()},
	}

	var output = "\n"
	for _, elem := range data {
		elemString := fmt.Sprintf("%-30s %-30s", elem[0], elem[1])
		output = output + elemString + "\n"
	}
	return output
}

type delegationTracker struct {
	nodeAddress common.Address // validator to which stake has been delegated
	delegator   common.Address
	runner      *Runner  // to fetch liquid contracts
	amount      *big.Int // delegation amount in NTN

	// trackers
	atn *changeTracker
	ntn *changeTracker // LNTN value in NTN
}

func newDelegationTracker(r *Runner, nodeAddress common.Address, delegator common.Address, amount *big.Int) *delegationTracker {
	return &delegationTracker{
		nodeAddress: nodeAddress,
		delegator:   delegator,
		runner:      r,
		amount:      set(amount),

		atn: newChangeTracker(unclaimedRewards(r.LiquidStateContract(nodeAddress), delegator)),
		ntn: newChangeTracker(lntnBalanceInNtn(r, nodeAddress, delegator)),
	}
}

func (dt *delegationTracker) String() string {
	return fmt.Sprintf("validator %s, atn: %s, ntn: %s", dt.nodeAddress.String(), dt.atn.String(), dt.ntn.String())
}

func (dt *delegationTracker) Validate() error {
	atnErr := dt.atn.Validate(unclaimedRewards(dt.runner.LiquidStateContract(dt.nodeAddress), dt.delegator))
	if atnErr != nil {
		atnErr = fmt.Errorf("delegator unclaimed rewards error: %w", atnErr)
	}
	ntnErr := dt.ntn.Validate(lntnBalanceInNtn(dt.runner, dt.nodeAddress, dt.delegator))
	if ntnErr != nil {
		ntnErr = fmt.Errorf("delegator ntn balance error: %w", ntnErr)
	}

	return errors.Join(atnErr, ntnErr)
}

type delegatorTracker struct {
	delegator common.Address
	bonds     []*delegationTracker
}

func newDelegatorTracker(delegator common.Address) *delegatorTracker {
	return &delegatorTracker{
		delegator: delegator,
	}
}

func (dt *delegatorTracker) addBond(r *Runner, nodeAddress common.Address, delegator common.Address, amount *big.Int) {
	dt.bonds = append(dt.bonds, newDelegationTracker(
		r,
		nodeAddress,
		delegator,
		amount,
	))
}

func (dt *delegatorTracker) String() string {
	output := fmt.Sprintf("delegator %s\n", dt.delegator)
	for _, bond := range dt.bonds {
		output = output + "\t" + bond.String() + "\n"
	}
	return output
}

func (dt *delegatorTracker) Validate() error {
	var errSum error

	for _, bond := range dt.bonds {
		bondErr := bond.Validate()
		if bondErr != nil {
			errSum = errors.Join(errSum,
				fmt.Errorf("delegation to %s mismatch: %w", bond.nodeAddress, bondErr),
			)
		}
	}
	return errSum
}

// from standard denom to weis
func toWei(standard uint64) *big.Int {
	return mul(new(big.Int).SetUint64(standard), tenTo18)
}

func committee(r *Runner) []IAutonityCommitteeMember {
	c, _, err := r.Autonity.GetCommittee(nil)
	require.NoError(r.T, err)
	return c
}

func treasuryAccount(r *Runner) common.Address {
	autonityTreasuryAccount, _, err := r.Autonity.GetTreasuryAccount(nil)
	require.NoError(r.T, err)
	return autonityTreasuryAccount
}

func withheldRewardsPool(r *Runner) common.Address {
	config, _, err := r.Autonity.GetConfig(nil)
	require.NoError(r.T, err)
	return config.Policy.WithheldRewardsPool
}

func treasuryFee(r *Runner) *big.Int {
	autonityTreasuryFee, _, err := r.Autonity.GetTreasuryFee(nil)
	require.NoError(r.T, err)
	return autonityTreasuryFee
}

func mint(r *Runner, address common.Address, ntn *big.Int) {
	_, err := r.Autonity.Mint(r.Operator, address, ntn)
	require.NoError(r.T, err)
}

func oracleRewardRate(r *Runner) *big.Int {
	config, _, err := r.Autonity.GetConfig(nil)
	require.NoError(r.T, err)
	return config.Policy.OracleRewardRate
}

func proposerRewardRate(r *Runner) *big.Int {
	config, _, err := r.Autonity.GetConfig(nil)
	require.NoError(r.T, err)
	return config.Policy.ProposerRewardRate
}

func bondFrom(r *Runner, delegator common.Address, validator common.Address, amount *big.Int) {
	_, err := r.Autonity.Bond(FromSender(delegator, nil), validator, amount)
	require.NoError(r.T, err)
}

func epochPeriod(r *Runner) *big.Int {
	period, _, err := r.Autonity.GetEpochPeriod(nil)
	require.NoError(r.T, err)
	return period
}

func delegatedStake(r *Runner, nodeAddr common.Address) *big.Int {
	val := validator(r, nodeAddr)
	return sub(val.BondedStake, val.SelfBondedStake)
}

func unclaimedRewards(liquid *ILiquid, delegator common.Address) *big.Int {
	atns, _, err := liquid.UnclaimedRewards(nil, delegator)
	require.NoError(liquid.r.T, err)
	return atns
}

func lntnBalanceInNtn(r *Runner, nodeAddress common.Address, delegator common.Address) *big.Int {
	return lntnToNtn(r, nodeAddress, lntnBalance(r.LiquidStateContract(nodeAddress), delegator))
}

func lntnBalance(liquid *ILiquid, delegator common.Address) *big.Int {
	lntn, _, err := liquid.BalanceOf(nil, delegator)
	require.NoError(liquid.r.T, err)
	return lntn
}

func liquidSupply(r *Runner, nodeAddress common.Address) *big.Int {
	return validator(r, nodeAddress).LiquidSupply
}

// convert LNTN to NTN
func lntnToNtn(r *Runner, nodeAddress common.Address, lntn *big.Int) *big.Int {
	return div(mul(lntn, delegatedStake(r, nodeAddress)), liquidSupply(r, nodeAddress))
}

func commissionRateScaleFactor(r *Runner) *big.Int {
	liquid := r.contractObject(LiquidLogicMetaData, r.Committee.Validators[0].LiquidStateContract)
	scaleFactor, _, err := liquid.call(nil, liquid.abi.Methods["COMMISSION_RATE_SCALE_FACTOR"].Name)
	require.NoError(r.T, err)
	return new(big.Int).SetBytes(scaleFactor)
}

func feeFactorUnitRecip(r *Runner) *big.Int {
	liquid := r.contractObject(LiquidLogicMetaData, r.Committee.Validators[0].LiquidStateContract)
	feeFactor, _, err := liquid.call(nil, liquid.abi.Methods["FEE_FACTOR_UNIT_RECIP"].Name)
	require.NoError(r.T, err)
	return new(big.Int).SetBytes(feeFactor)
}

func supply(r *Runner, liquid *ILiquid) *big.Int {
	totalSupply, _, err := liquid.TotalSupply(nil)
	require.NoError(r.T, err)
	return totalSupply
}

func standardScaleFactor(r *Runner) *big.Int {
	scaleFactor, _, err := r.Autonity.STANDARDSCALEFACTOR(nil)
	require.NoError(r.T, err)
	return scaleFactor
}

func voterPerformance(r *Runner, voter common.Address) *big.Int {
	performance, _, err := r.Oracle.GetRewardPeriodPerformance(nil, voter)
	require.NoError(r.T, err)
	return performance
}

func totalPerformance(r *Runner) *big.Int {
	oracleTest := (*OracleTest)(r.Oracle)
	totPerformance, _, err := oracleTest.GetRewardPeriodAggregatedScore(nil)
	require.NoError(r.T, err)
	return totPerformance
}

func maxCommitteeSize(r *Runner) *big.Int {
	config, _, err := r.Autonity.GetConfig(nil)
	require.NoError(r.T, err)
	return config.Protocol.CommitteeSize
}

func mul(a, b *big.Int) *big.Int {
	return new(big.Int).Mul(a, b)
}

func div(a, b *big.Int) *big.Int {
	return new(big.Int).Div(a, b)
}

func sub(a, b *big.Int) *big.Int {
	return new(big.Int).Sub(a, b)
}

func add(a, b *big.Int) *big.Int {
	return new(big.Int).Add(a, b)
}

func set(a *big.Int) *big.Int {
	return new(big.Int).Set(a)
}

func strictlyPositive(a *big.Int) bool {
	return a.Cmp(common.Big0) > 0
}

func deployOracleTestContract(r *Runner) (common.Address, *OracleTest) {
	config := r.Config
	voters := make([]common.Address, len(config.AutonityContractConfig.Validators))
	treasuries := make([]common.Address, len(config.AutonityContractConfig.Validators))
	validators := make([]common.Address, len(config.AutonityContractConfig.Validators))
	for i, val := range config.AutonityContractConfig.Validators {
		voters[i] = val.OracleAddress
		treasuries[i] = val.Treasury
		validators[i] = *val.NodeAddress
	}

	oracleConfig := OracleTestConfig{
		Autonity:                  params.AutonityContractAddress,
		Operator:                  config.AutonityContractConfig.Operator,
		VotePeriod:                new(big.Int).SetUint64(config.OracleContractConfig.VotePeriod),
		OutlierDetectionThreshold: new(big.Int).SetUint64(config.OracleContractConfig.OutlierDetectionThreshold),
		OutlierSlashingThreshold:  new(big.Int).SetUint64(config.OracleContractConfig.OutlierSlashingThreshold),
		BaseSlashingRate:          new(big.Int).SetUint64(config.OracleContractConfig.BaseSlashingRate),
		NonRevealThreshold:        new(big.Int).SetUint64(config.OracleContractConfig.NonRevealThreshold),
		RevealResetInterval:       new(big.Int).SetUint64(config.OracleContractConfig.RevealResetInterval),
		SlashingRateCap:           new(big.Int).SetUint64(config.OracleContractConfig.SlashingRateCap),
	}

	address, _, oracleTest, err := r.DeployOracleTest(nil,
		voters, validators, treasuries,
		config.OracleContractConfig.Symbols,
		oracleConfig,
	)
	require.NoError(r.T, err)

	return address, oracleTest
}

/*
 * NOTE
 * This test doesn't account for:
 * - provable accountability slashing
 * - omission accountability slashing
 * - oracle slashing
 * - reward loss due to jailing
 *    - if jailed for provable fault, all rewards go to the reporter
 *    - if jailed for omission, all rewards go to a protocol pool
 *
 *  HOWEVER
 *  It does take into account rewards withholding due to omission score
 *  as it might be more common wrt the previous cases
 */
func TestRewardDistribution(t *testing.T) {
	r := Setup(t, func(genesis *params.AutonityContractGenesis) *params.AutonityContractGenesis {
		// decouple autonity treasury and withheld rewards pool for easier accounting
		genesis.WithheldRewardsPool = common.Address{0xca, 0xfe}
		return genesis
	})

	// substitute Oracle with OracleTest
	oracleTestAddress, oracleTest := deployOracleTestContract(r)
	r.NoError(r.Autonity.SetOracleContract(r.Operator, oracleTestAddress))
	r.Oracle = (*Oracle)(oracleTest)

	// substitute OmissionAccountability with OmissionAccountabilityTest
	// need to use evm.Replace as the absenteeComputer precompile only allows the original
	// omission accountability address to make calls
	_, _, _, err := r.Evm.Replace(vm.AccountRef(params.DeployerAddress), generated.OmissionAccountabilityTestBytecode, params.OmissionAccountabilityContractAddress)
	require.NoError(r.T, err)
	omissionTest := &OmissionAccountabilityTest{r.contractObject(OmissionAccountabilityTestMetaData, params.OmissionAccountabilityContractAddress)}

	// fetch committee
	c := committee(r)
	csize := len(c)
	csizeBig := new(big.Int).SetUint64(uint64(csize))

	// bond some delegated stake to the validators
	numDelegators := 4
	delegationAmount := toWei(uint64(500)) // 500 NTN
	delegators := make([]common.Address, 0, numDelegators)
	for i := range numDelegators {
		delegator := common.BytesToAddress([]byte{byte(i + 1)})
		delegators = append(delegators, delegator)

		// mint and bond
		mint(r, delegator, mul(delegationAmount, csizeBig))
		for _, member := range c {
			bondFrom(r, delegator, member.Addr, delegationAmount)
		}
	}
	// apply bonding operations
	_, err = r.Autonity.ApplyStakingOperations(nil)
	require.NoError(t, err)

	// create delegator trackers
	delegatorsTrackers := make([]*delegatorTracker, 0, numDelegators)
	//t.Log("delegator trackers")
	for _, delegator := range delegators {
		tracker := newDelegatorTracker(delegator)
		for _, member := range c {
			tracker.addBond(r, member.Addr, delegator, delegationAmount)
		}
		//t.Log(tracker)
		delegatorsTrackers = append(delegatorsTrackers, tracker)
	}

	period := epochPeriod(r)
	t.Logf("epoch period: %s", period.String())

	// create validator trackers
	t.Logf("committee size: %d", csize)
	validatorTrackers := make([]*validatorTracker, 0, csize)
	totalPower := new(big.Int)
	//t.Log("validator trackers:")
	artificialTotalPerformance := new(big.Int)
	omissionFactor := omissionScaleFactor(r)
	artificialInactivityScore := div(omissionFactor, hundred) // 1%
	for _, member := range c {
		tracker := newValidatorTracker(r, member.Addr)
		validatorTrackers = append(validatorTrackers, tracker)
		totalPower = add(totalPower, tracker.power.value)
		//t.Log(tracker.String())

		// add artificial oracle performance
		r.NoError(oracleTest.SetPerformance(nil, tracker.oracleAddress, hundred))
		artificialTotalPerformance = add(artificialTotalPerformance, hundred)

		// add 1% (100) of inactivity to all validators.
		// will result in 1% of rewards getting withheld
		r.NoError(omissionTest.SetInactivityScore(nil, tracker.nodeAddress, artificialInactivityScore))
	}
	r.NoError(oracleTest.SetRewardPeriodAggregatedScore(nil, artificialTotalPerformance))

	// create autonity treasury tracker
	treasuryAcc := treasuryAccount(r)
	treasuryTrackerAtn := newChangeTracker(r.GetBalanceOf(treasuryAcc))
	treasuryTrackerNtn := newChangeTracker(r.GetNewtonBalanceOf(treasuryAcc))

	// create withheld rewards pool tracker
	withheldPoolAcc := withheldRewardsPool(r)
	withheldAtnTracker := newChangeTracker(r.GetBalanceOf(withheldPoolAcc))
	withheldNtnTracker := newChangeTracker(r.GetNewtonBalanceOf(withheldPoolAcc))

	// position at the last block before epoch end
	r.WaitNBlocks(int(period.Uint64() - 1))

	// total proposer effort for the epoch needs to be > 0 for comprehensive reward distribution
	totProposerEffort := totalProposerEffort(r)
	t.Logf("total proposer effort: %s", totProposerEffort.String())
	require.True(t, strictlyPositive(totProposerEffort))

	// total oracle performance for the epoch needs to be > 0 for comprehensive reward distribution
	totPerformance := totalPerformance(r)
	t.Logf("total oracle performance: %s", totPerformance.String())
	require.True(t, strictlyPositive(totPerformance))

	// compute expected rewards for each committee member and for treasury address

	// artificially set rewards amount
	atnRewards := uint64(100) // in ATN
	ntnRewards := uint64(100) // in NTN
	atn := toWei(atnRewards)
	ntn := toWei(ntnRewards)

	// step 1. compute autonity treasury fee
	autonityTreasuryFee := treasuryFee(r)
	autonityTresuryFeePerc := float64(autonityTreasuryFee.Uint64()) * 100 / float64(tenTo18.Uint64())
	t.Logf("autonity treasury fee %s (%.2f %%)", autonityTreasuryFee.String(), autonityTresuryFeePerc)
	treasuryRewards := div(mul(autonityTreasuryFee, atn), tenTo18)
	t.Logf("expected treasury rewards: %s", treasuryRewards.String())
	treasuryTrackerAtn.update(treasuryRewards)

	// step 2. compute remaining rewards
	remainingAtn := sub(atn, treasuryRewards)
	remainingNtn := set(ntn)

	// step 3. subtract oracle rewards and proposer rewards
	scaleFactor := standardScaleFactor(r)
	oracleRate := oracleRewardRate(r)
	atnOracleRewards := div(mul(oracleRate, remainingAtn), scaleFactor)
	ntnOracleRewards := div(mul(oracleRate, remainingNtn), scaleFactor)

	t.Logf("atn oracle rewards: %s", atnOracleRewards.String())
	t.Logf("ntn oracle rewards: %s", ntnOracleRewards.String())

	maxCsize := maxCommitteeSize(r)
	proposerRate := proposerRewardRate(r)
	atnProposerRewards := div(mul(mul(proposerRate, remainingAtn), csizeBig), mul(scaleFactor, maxCsize))
	ntnProposerRewards := div(mul(mul(proposerRate, remainingNtn), csizeBig), mul(scaleFactor, maxCsize))

	t.Logf("atn proposer rewards: %s", atnProposerRewards.String())
	t.Logf("ntn proposer rewards: %s", ntnProposerRewards.String())

	remainingAtn = sub(remainingAtn, add(atnOracleRewards, atnProposerRewards))
	remainingNtn = sub(remainingNtn, add(ntnOracleRewards, ntnProposerRewards))

	// track how much atn dust is accumulated in total
	totalAtnDust := new(big.Int)
	totalAtnForDelegators := new(big.Int)

	// step 4. compute rewards for each validator
	commissionScaleFactor := commissionRateScaleFactor(r)
	for _, v := range validatorTrackers {
		allocatedAtns := div(mul(v.power.value, remainingAtn), totalPower)
		allocatedNtns := div(mul(v.power.value, remainingNtn), totalPower)

		withheldAtn := div(mul(allocatedAtns, artificialInactivityScore), omissionFactor)
		withheldNtn := div(mul(allocatedNtns, artificialInactivityScore), omissionFactor)
		allocatedAtns = sub(allocatedAtns, withheldAtn)
		allocatedNtns = sub(allocatedNtns, withheldNtn)
		withheldAtnTracker.update(withheldAtn)
		withheldNtnTracker.update(withheldNtn)

		atnDelegatedReward := div(mul(v.delegatedPower.value, allocatedAtns), v.power.value)
		atnSelfReward := sub(allocatedAtns, atnDelegatedReward)
		v.treasuryAddressAtn.update(atnSelfReward)

		ntnDelegatedReward := div(mul(v.delegatedPower.value, allocatedNtns), v.power.value)
		ntnSelfReward := sub(allocatedNtns, ntnDelegatedReward)
		v.selfPower.update(ntnSelfReward)
		v.power.update(ntnSelfReward)

		// apply commission rate on delegated rewards
		atnCommission := div(mul(v.commission, atnDelegatedReward), commissionScaleFactor)
		v.treasuryAddressAtn.update(atnCommission)
		ntnCommission := div(mul(v.commission, ntnDelegatedReward), commissionScaleFactor)
		v.selfPower.update(ntnCommission)
		v.delegatedPower.update(sub(ntnDelegatedReward, ntnCommission))
		v.power.update(ntnDelegatedReward)

		// update expected rewards for delegators
		atnForDelegators := sub(atnDelegatedReward, atnCommission)
		ntnForDelegators := sub(ntnDelegatedReward, ntnCommission)

		// due to how the fees are tracked in the LNTN contract, there is some precision loss.
		// compute how much ATN dust is lost and round off the `atnForDelegators` value
		lntnSupply := supply(r, r.LiquidStateContract(v.nodeAddress))
		atnFeeFactor := div(mul(atnForDelegators, feeFactorUnitRecip(r)), lntnSupply)
		atnMaxClaimable := div(mul(atnFeeFactor, lntnSupply), feeFactorUnitRecip(r))
		atnDust := sub(atnForDelegators, atnMaxClaimable)
		// update stats
		totalAtnDust = add(totalAtnDust, atnDust)
		totalAtnForDelegators = add(totalAtnForDelegators, atnForDelegators)
		// print summary
		t.Logf("validator %s ATN dust  --------------------", v.nodeAddress)
		t.Logf("atn for delegators: %s", atnForDelegators.String())
		t.Logf("atn fee factor: %s", atnFeeFactor.String())
		t.Logf("atn max claimable: %s", atnMaxClaimable.String())
		t.Logf("atn dust: %s (%d)", atnDust.String(), len(atnDust.String()))
		// round off atnForDelegators
		atnForDelegators = sub(atnForDelegators, atnDust)

		for _, tracker := range delegatorsTrackers {
			for _, bond := range tracker.bonds {
				if bond.nodeAddress == v.nodeAddress {
					delegatorAtnReward := div(mul(bond.amount, atnForDelegators), v.delegatedPower.value)
					bond.atn.update(delegatorAtnReward)
					ntnDelegatorReward := div(mul(bond.amount, ntnForDelegators), v.delegatedPower.value)
					bond.ntn.update(ntnDelegatorReward)
				}
			}
		}

		// proposer rewards
		atnProposerReward := div(mul(proposerEffort(r, v.nodeAddress), atnProposerRewards), totProposerEffort)
		v.treasuryAddressAtn.update(atnProposerReward)
		ntnProposerReward := div(mul(proposerEffort(r, v.nodeAddress), ntnProposerRewards), totProposerEffort)
		v.selfPower.update(ntnProposerReward)
		v.power.update(ntnProposerReward)

		// oracle rewards
		atnOracleReward := div(mul(voterPerformance(r, v.oracleAddress), atnOracleRewards), totPerformance)
		v.treasuryAddressAtn.update(atnOracleReward)
		ntnOracleReward := div(mul(voterPerformance(r, v.oracleAddress), ntnOracleRewards), totPerformance)
		v.selfPower.update(ntnOracleReward)
		v.power.update(ntnOracleReward)
	}

	// print summary for lost ATN dust
	t.Logf("ATN dust summary ----------------")
	dustPerc := float64(totalAtnDust.Uint64()) * 100 / float64(totalAtnForDelegators.Uint64())
	t.Logf(
		"\ntotal atn dust %s (%d)\n"+
			"total atn for delegators %s (%d)\n"+
			"dust is %.10f %% of total rewards",
		totalAtnDust.String(), len(totalAtnDust.String()),
		totalAtnForDelegators.String(), len(totalAtnForDelegators.String()),
		dustPerc,
	)

	// update omission scores
	r.NoError(omissionTest.UpdateScores(nil))

	// execute reward redistribution
	t.Logf("executing reward distribution at block %s", r.Evm.Context.BlockNumber.String())
	r.GiveMeSomeMoney(r.Autonity.Address(), atn)
	mint(r, r.Autonity.Address(), ntn)
	_, err = r.Autonity.PerformRedistribution(nil, atn, ntn)
	require.NoError(t, err)

	// check that actual balances match expected ones

	t.Logf("FINAL STATE --------------")

	t.Logf("TREASURY -------")
	t.Logf("atn --> %s", treasuryTrackerAtn.String())
	assert.NoError(t, treasuryTrackerAtn.Validate(r.GetBalanceOf(treasuryAcc)))
	t.Logf("ntn --> %s", treasuryTrackerNtn.String())
	assert.NoError(t, treasuryTrackerNtn.Validate(r.GetNewtonBalanceOf(treasuryAcc)))

	t.Logf("WITHHELD POOL -------")
	t.Logf("atn --> %s", withheldAtnTracker.String())
	assert.NoError(t, withheldAtnTracker.Validate(r.GetBalanceOf(withheldPoolAcc)))
	t.Logf("ntn --> %s", withheldNtnTracker.String())
	assert.NoError(t, withheldNtnTracker.Validate(r.GetNewtonBalanceOf(withheldPoolAcc)))

	t.Logf("VALIDATORS -------")
	for _, tracker := range validatorTrackers {
		t.Log(tracker)
		assert.NoError(t, tracker.Validate(r))
	}

	t.Logf("DELEGATORS -------")
	for _, tracker := range delegatorsTrackers {
		t.Log(tracker)
		assert.NoError(t, tracker.Validate())
	}
}
