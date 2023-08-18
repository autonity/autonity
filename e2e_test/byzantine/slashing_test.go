package byzantine

import (
	"context"
	"fmt"
	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/common"
	core2 "github.com/autonity/autonity/core"
	"github.com/autonity/autonity/crypto"
	"math/big"
	"testing"
	"time"

	"github.com/autonity/autonity/autonity"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/ethclient"
	"github.com/autonity/autonity/node"
	"github.com/stretchr/testify/require"
)

// *** e2e tests to write ***:
// Part 1: Slashing tests
// - Staking with penalty absorbing stake
// - Funds moved to treasury
// - Reward redistribution
// - Jail
// - Silence consensus

// Part 2: Accusation flow tests:

//  Validator is accused, submit proof of innocence,
//  Validator is accused, do not submit proof of innocence,
//  Validator is accused and accused again
//  Validators is accused and someone sent direct proof of misbehavior

//  Need to test canAccuse/canSlash for each scenario

// Part 3: Fuzz tests on event handler

// invalid ValidRound message?

func runSlashingTest(ctx context.Context, t *testing.T, nodesCount int, epochPeriod, stake, selfBondedStake uint64, faultyNodes []int, offendersCount, faultsCount uint64, epochs int) (uint64, []autonity.AutonityValidator, []autonity.AutonityValidator) {

	validators, err := e2e.Validators(t, nodesCount, fmt.Sprintf("10e36,v,%d,0.0.0.0:%%s,%%s", selfBondedStake))
	require.NoError(t, err)

	// set Malicious validators
	for _, faultyNodeIndex := range faultyNodes {
		validators[faultyNodeIndex].TendermintServices = &node.TendermintServices{Broadcaster: &InvalidProposal{}}
	}

	validatorsBefore := make([]autonity.AutonityValidator, len(faultyNodes))
	validatorsAfter := make([]autonity.AutonityValidator, len(faultyNodes))

	var baseRate uint64
	var collusionFactor uint64
	var historyFactor uint64
	var slashingPrecision uint64

	// creates a network of validators and starts all the nodes in it
	// and sneak default slashing parameters
	network, err := e2e.NewNetworkFromValidators(t, validators, true, func(genesis *core2.Genesis) {
		baseRate = genesis.Config.AccountabilityConfig.BaseSlashingRateMid
		collusionFactor = genesis.Config.AccountabilityConfig.CollusionFactor
		historyFactor = genesis.Config.AccountabilityConfig.HistoryFactor
		slashingPrecision = genesis.Config.AccountabilityConfig.SlashingRatePrecision
		genesis.Config.AutonityContractConfig.EpochPeriod = epochPeriod

		//if stake is set, it means we need some other account to bond
		if stake > 0 {

			validatorAddress := genesis.Config.AutonityContractConfig.Validators[faultyNodes[0]].NodeAddress

			key, _ := crypto.GenerateKey()
			address := crypto.PubkeyToAddress(key.PublicKey)

			genesis.Alloc[address] = core2.GenesisAccount{
				NewtonBalance: big.NewInt(int64(stake)),
				Balance:       new(big.Int),
				Bonds: map[common.Address]*big.Int{
					*validatorAddress: big.NewInt(int64(stake)),
				},
			}
		}
	})
	require.NoError(t, err)
	defer network.Shutdown()

	dedicatedNode := network[1].WsClient

	autonityContract, err := autonity.NewAutonity(autonity.AutonityContractAddress, dedicatedNode)
	require.NoError(t, err)

	treasuryAccount, err := autonityContract.GetTreasuryAccount(nil)
	require.NoError(t, err)

	balanceBefore, err := autonityContract.BalanceOf(nil, treasuryAccount)
	require.NoError(t, err)

	extraEpochsSlashed := uint64(0)

	// run extra epochs
	for i := 1; i < epochs; i++ {
		// scale timeout with extra 10% of expected time

		timeout, cancel := context.WithTimeout(ctx, time.Duration(float32(epochPeriod)*1.1)*time.Second)
		defer cancel()
		slashingEvents := WaitForSlashingEvents(timeout, t, len(faultyNodes), dedicatedNode)

		for _, slashingEvent := range slashingEvents {
			extraEpochsSlashed += slashingEvent.Amount.Uint64()
		}
	}

	for i, faultNodeIndex := range faultyNodes {
		validatorBefore, err := autonityContract.GetValidator(nil, network[faultNodeIndex].Address)
		require.NoError(t, err)

		validatorsBefore[i] = validatorBefore
	}

	timeout, cancel := context.WithTimeout(ctx, time.Duration(float32(epochPeriod)*1.1)*time.Second)
	defer cancel()
	slashingEvents := WaitForSlashingEvents(timeout, t, len(faultyNodes), dedicatedNode)

	for i, faultNodeIndex := range faultyNodes {
		validatorAfter, err := autonityContract.GetValidator(nil, network[faultNodeIndex].Address)
		require.NoError(t, err)

		validatorsAfter[i] = validatorAfter
	}

	balanceAfter, err := autonityContract.BalanceOf(nil, treasuryAccount)
	require.NoError(t, err)

	// ensure funds were transferred to treasure account
	require.Greater(t, balanceAfter.Uint64(), balanceBefore.Uint64())

	// check if the increase in treasury account matches sum of slashing penalties
	expectedSlashedAmount := extraEpochsSlashed
	for _, slashingEvent := range slashingEvents {
		expectedSlashedAmount += slashingEvent.Amount.Uint64()
	}

	require.Equal(t, balanceAfter.Uint64()-balanceBefore.Uint64(), expectedSlashedAmount)

	// check if slashing amount is calculated properly

	//as per ADR, expected slashing rate is:
	slashingRate := baseRate + (offendersCount * collusionFactor) + (faultsCount * historyFactor)
	if slashingRate >= slashingPrecision {
		slashingRate = slashingPrecision
	}

	for i := range faultyNodes {
		// slashing amounts are based on bonded stake before a penalty
		expectedSlashAmount := new(big.Int).Div(new(big.Int).Mul(big.NewInt(int64(slashingRate)), validatorsBefore[i].BondedStake), new(big.Int).SetUint64(slashingPrecision))

		require.Equal(t, expectedSlashAmount, slashingEvents[i].Amount)
	}

	return expectedSlashedAmount, validatorsBefore, validatorsAfter
}

func TestSimpleSlashing(t *testing.T) {
	runSlashingTest(context.TODO(), t, 4, 40, 0, 100, []int{2}, 1, 0, 1)
}

func TestPenaltyAbsorbingStake(t *testing.T) {

	stake := uint64(200)
	selfBondedStake := uint64(500)

	expectedSlashingAmount, validatorsBefore, validatorsAfter := runSlashingTest(context.TODO(), t, 4, 40, stake, selfBondedStake, []int{2}, 1, 0, 1)

	expectedSelfBondedStake := validatorsBefore[0].SelfBondedStake.Uint64() - expectedSlashingAmount

	if expectedSlashingAmount >= validatorsBefore[0].SelfBondedStake.Uint64() {
		expectedSelfBondedStake = 0
	}

	// make sure selfBondedStake was set properly
	require.Equal(t, selfBondedStake, validatorsBefore[0].SelfBondedStake.Uint64())
	require.Equal(t, stake+selfBondedStake, validatorsBefore[0].BondedStake.Uint64())

	require.Equal(t, expectedSelfBondedStake, validatorsAfter[0].SelfBondedStake.Uint64())
}

func TestMultipleOffender(t *testing.T) {

	// increased number of nodes means there can be more blocks between offences being reported
	// for the test we want them all in one epoch
	runSlashingTest(context.TODO(), t, 6, 100, 0, 100, []int{2, 3}, 2, 0, 1)

}

func TestHistoryFactor(t *testing.T) {

	// Simulating historical slashing is a bit complicated, since there are extra mechanism involved.
	// 1. We do "simple" slashing first, misbehaving node is now put in jail, so it cannot be part of committee
	// and cannot be slashed again until reactivated
	// 2. Once the "sentence", expressed in number of blocks, passes, we re-activate the node
	// 3. We now need to wait for a next epoch
	// 4. Since the node continue to produce invalid proposal, we can expect second slashing soon

	validators, err := e2e.Validators(t, 4, "10e36,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	// set Malicious validators
	faultyNode := 2
	validators[faultyNode].TendermintServices = &node.TendermintServices{Broadcaster: &InvalidProposal{}}

	var chainID *big.Int

	// creates a network of 4 validators and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, validators, true, func(genesis *core2.Genesis) {
		genesis.Config.AutonityContractConfig.EpochPeriod = 50
		genesis.Config.AccountabilityConfig.JailFactor = 1
		chainID = genesis.Config.ChainID
	})
	require.NoError(t, err)
	defer network.Shutdown()

	dedicatedNode := network[1].WsClient

	autonityContract, err := autonity.NewAutonity(autonity.AutonityContractAddress, dedicatedNode)
	require.NoError(t, err)

	accountabilityContract, err := autonity.NewAccountability(autonity.AccountabilityContractAddress, dedicatedNode)
	require.NoError(t, err)

	treasuryAccount, err := autonityContract.GetTreasuryAccount(nil)
	require.NoError(t, err)

	balanceBefore, err := autonityContract.BalanceOf(nil, treasuryAccount)
	require.NoError(t, err)

	// wait for slashing
	timeout, cancel := context.WithTimeout(context.TODO(), 120*time.Second)
	defer cancel()
	slashingEventA := WaitForSlashingEvent(timeout, t, dedicatedNode)

	// wait until we can un-jail (+1 just in case)
	err = network.WaitForHeight(slashingEventA.ReleaseBlock.Uint64()+1, 60)
	require.NoError(t, err)

	//un-jail
	transactOpts, err := bind.NewKeyedTransactorWithChainID(
		validators[faultyNode].TreasuryKey,
		chainID,
	)
	require.NoError(t, err)

	_, err = autonityContract.ActivateValidator(transactOpts, network[faultyNode].Address)
	require.NoError(t, err)

	validatorBefore, err := autonityContract.GetValidator(nil, network[faultyNode].Address)
	require.NoError(t, err)

	// wait for slashing again
	timeout, cancel = context.WithTimeout(context.TODO(), 150*time.Second)
	defer cancel()
	slashingEventB := WaitForSlashingEvent(timeout, t, dedicatedNode)

	balanceAfter, err := autonityContract.BalanceOf(nil, treasuryAccount)
	require.NoError(t, err)

	// ensure funds were transferred to treasure account
	require.Greater(t, balanceAfter.Uint64(), balanceBefore.Uint64())

	// check if the increase in treasury account matches slashing penalty/ies
	require.Equal(t, balanceAfter.Uint64()-balanceBefore.Uint64(), slashingEventA.Amount.Uint64()+slashingEventB.Amount.Uint64())

	// check if slashing amount is calculated properly

	accountabilityConfig, err := accountabilityContract.Config(nil)
	require.NoError(t, err)

	baseRate := accountabilityConfig.BaseSlashingRateMid
	collusionFactor := accountabilityConfig.CollusionFactor
	historyFactor := accountabilityConfig.HistoryFactor
	slashingRatePrecision := accountabilityConfig.SlashingRatePrecision

	// one offender in total
	offendersCount := int64(1)
	// but this should be repeated offence
	faultCount := int64(1)

	//as per ADR, expected slashing rate is:
	// base rate + offendersCount * collusionFactor

	slashingRate := new(big.Int).Add(baseRate,
		new(big.Int).Add(
			new(big.Int).Mul(big.NewInt(offendersCount), collusionFactor),
			new(big.Int).Mul(big.NewInt(faultCount), historyFactor),
		),
	)

	if slashingRate.Cmp(slashingRatePrecision) >= 0 {
		slashingRate = slashingRatePrecision
	}

	// slashing amounts are based on bonded stake before a penalty
	expectedSlashAmount := new(big.Int).Div(new(big.Int).Mul(slashingRate, validatorBefore.BondedStake), slashingRatePrecision)

	require.Equal(t, expectedSlashAmount, slashingEventB.Amount)
}

// Wait for N AccountabilitySlashingEvent to appear on all the nodes in the network
func WaitForSlashingEvents(ctx context.Context, t *testing.T, n int, client *ethclient.Client) []*autonity.AccountabilitySlashingEvent {

	accountabilityContract, err := autonity.NewAccountability(autonity.AccountabilityContractAddress, client)
	require.NoError(t, err)

	// wait for slashing event
	eventsSink := make(chan *autonity.AccountabilitySlashingEvent, n)
	subscription, err := accountabilityContract.WatchSlashingEvent(nil, eventsSink)
	require.NoError(t, err)

	defer subscription.Unsubscribe()

	events := make([]*autonity.AccountabilitySlashingEvent, 0, n)

loop:
	for {
		select {
		case <-ctx.Done():
			t.Error("timeout")
			break loop
		case err := <-subscription.Err():
			t.Errorf("subscription failed: %s", err)
			break loop
		case e := <-eventsSink:
			events = append(events, e)
			if len(events) == n {
				return events
			}
		}
	}

	t.Fatalf("not enough slashing events, wanted %d got %d", n, len(events))
	return events
}

func WaitForSlashingEvent(ctx context.Context, t *testing.T, client *ethclient.Client) *autonity.AccountabilitySlashingEvent {
	e := WaitForSlashingEvents(ctx, t, 1, client)
	return e[0]
}
