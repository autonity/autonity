package autonity

import (
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/metrics"
	"math/big"
	"sync"
	"testing"
)

const (
	testAddress1 = "70524d664ffe731100208a0154e556f9bb679ae6"
	testAddress2 = "70524d664ffe731100208a0154e556f9bb679ae5"
	testAddress3 = "70524d664ffe731100208a0154e556f9bb679ae4"
)

func TestContract_generateMetricsIDs(t *testing.T) {
	t.Run("test generate block reward metric ID", func(t *testing.T) {
		contract := &Contract{}
		blockNumber := uint64(2)
		blockRewardMetricID := contract.generateBlockRewardMetricsID(blockNumber)
		if blockRewardMetricID != fmt.Sprintf(BlockRewardBlockMetricID, blockNumber) {
			t.Fatal("wrong result expected.")
		}
	})

	t.Run("test generate reward distribution metric ID", func(t *testing.T) {
		contract := &Contract{}
		a := common.Hex2Bytes(testAddress1)
		address := common.BytesToAddress(a)
		metricID := contract.generateRewardDistributionMetricsID(address, Participant, uint64(2))
		expectedID := fmt.Sprintf("contract/block/2/user/%s/participant/reward", address.String())
		if metricID != expectedID {
			t.Fatal("case failed.")
		}
	})

	t.Run("test generate user metrics ID", func(t *testing.T) {
		contract := &Contract{}
		a := common.Hex2Bytes(testAddress1)
		address := common.BytesToAddress(a)
		stakeID, balanceID, commissionRateID, _ := contract.generateUserMetricsID(address, Participant)
		expectedStakeID := fmt.Sprintf(UserMetricIDTemplate, address.String(), "participant", "stake")
		expectedBalanceID := fmt.Sprintf(UserMetricIDTemplate, address.String(), "participant", "balance")
		expectedCommissionRateID := fmt.Sprintf(UserMetricIDTemplate, address.String(), "participant", "commissionrate")
		if stakeID != expectedStakeID || balanceID != expectedBalanceID || commissionRateID != expectedCommissionRateID {
			t.Fatal("test cas failed.")
		}
	})

	t.Run("test generate user metrics ID with wrong role", func(t *testing.T) {
		contract := &Contract{}
		a := common.Hex2Bytes(testAddress1)
		address := common.BytesToAddress(a)
		_, _, _, err := contract.generateUserMetricsID(address, 3)
		if err == nil {
			t.Fatal("test cas failed.")
		}
	})

	t.Run("test resolve user type name", func(t *testing.T) {
		contract := &Contract{}
		name := contract.resolveUserTypeName(Participant)
		if name != "participant" {
			t.Fatal("case failed.")
		}

		name = contract.resolveUserTypeName(Stakeholder)
		if name != "stakeholder" {
			t.Fatal("case failed.")
		}

		name = contract.resolveUserTypeName(Validator)
		if name != "validator" {
			t.Fatal("case failed")
		}

		name = contract.resolveUserTypeName(4)
		if name != "unknown" {
			t.Fatal("case failed.")
		}
	})
}

func TestContract_removeMetricsFromRegistry(t *testing.T) {
	t.Run("remove user metrics from metric registry", func(t *testing.T) {
		// prepare context in metric registry
		contract := &Contract{
			metricDataMutex:  sync.RWMutex{},
			heightLowBounder: 0,
			RWMutex:          sync.RWMutex{},
		}
		blockHeight := uint64(10)
		a1 := common.Hex2Bytes(testAddress1)
		address1 := common.BytesToAddress(a1)
		a2 := common.Hex2Bytes(testAddress2)
		address2 := common.BytesToAddress(a2)
		a3 := common.Hex2Bytes(testAddress3)
		address3 := common.BytesToAddress(a3)

		stakeID1, balanceID1, commisionRateID1, _ := contract.generateUserMetricsID(address1, Participant)
		rewardDistributionMetricID1 := contract.generateRewardDistributionMetricsID(address1, Stakeholder, blockHeight)

		metrics.GetOrRegisterCounter(rewardDistributionMetricID1, nil).Inc(100)
		metrics.GetOrRegisterGauge(stakeID1, nil).Update(100)
		metrics.GetOrRegisterGauge(balanceID1, nil).Update(100)
		metrics.GetOrRegisterGauge(commisionRateID1, nil).Update(100)

		stakeID2, balanceID2, commisionRateID2, _ := contract.generateUserMetricsID(address2, Stakeholder)
		rewardDistributionMetricID2 := contract.generateRewardDistributionMetricsID(address2, Stakeholder, blockHeight)

		metrics.GetOrRegisterCounter(rewardDistributionMetricID2, nil).Inc(200)
		metrics.GetOrRegisterGauge(stakeID2, nil).Update(200)
		metrics.GetOrRegisterGauge(balanceID2, nil).Update(200)
		metrics.GetOrRegisterGauge(commisionRateID2, nil).Update(200)

		stakeID3, balanceID3, commisionRateID3, _ := contract.generateUserMetricsID(address3, Validator)
		rewardDistributionMetricID3 := contract.generateRewardDistributionMetricsID(address3, Stakeholder, blockHeight)

		metrics.GetOrRegisterCounter(rewardDistributionMetricID3, nil).Inc(300)
		metrics.GetOrRegisterGauge(stakeID3, nil).Update(300)
		metrics.GetOrRegisterGauge(balanceID3, nil).Update(300)
		metrics.GetOrRegisterGauge(commisionRateID3, nil).Update(300)

		contract.removeMetricsFromRegistry(address1, blockHeight)
		contract.removeMetricsFromRegistry(address2, blockHeight)
		contract.removeMetricsFromRegistry(address3, blockHeight)

		if metrics.Get(stakeID1) != nil || metrics.Get(balanceID1) != nil || metrics.Get(commisionRateID1) != nil ||
			metrics.Get(rewardDistributionMetricID1) != nil {
			t.Fatal("case failed.")
		}

		if metrics.Get(stakeID2) != nil || metrics.Get(balanceID2) != nil || metrics.Get(commisionRateID2) != nil ||
			metrics.Get(rewardDistributionMetricID2) != nil {
			t.Fatal("case failed.")
		}

		if metrics.Get(stakeID3) != nil || metrics.Get(balanceID3) != nil || metrics.Get(commisionRateID3) != nil ||
			metrics.Get(rewardDistributionMetricID3) != nil {
			t.Fatal("case failed.")
		}

	})
}

func TestContract_CleanUselessMetrics(t *testing.T) {
	t.Run("clean up metrics for removed users, exception case: input address set is empty.", func(t *testing.T) {
		// prepare context in metric registry
		contract := &Contract{
			metricDataMutex:  sync.RWMutex{},
			heightLowBounder: 0,
		}
		blockHeight := uint64(10)
		contract.cleanUselessMetrics(nil, blockHeight)
	})

	t.Run("clean up metrics for removed users, exception case: local address set is empty.", func(t *testing.T) {
		contract := &Contract{
			metricDataMutex:  sync.RWMutex{},
			users:            nil,
			heightLowBounder: 0,
		}

		a1 := common.Hex2Bytes(testAddress1)
		address1 := common.BytesToAddress(a1)
		a2 := common.Hex2Bytes(testAddress2)
		address2 := common.BytesToAddress(a2)
		a3 := common.Hex2Bytes(testAddress3)
		address3 := common.BytesToAddress(a3)

		var users []common.Address
		users = append(users, address1, address2, address3)
		contract.cleanUselessMetrics(users, uint64(10))
		if len(contract.users) != 3 {
			t.Fatal("case failed.")
		}
	})

	t.Run("clean up metrics for removed users, normal case.", func(t *testing.T) {
		a1 := common.Hex2Bytes(testAddress1)
		address1 := common.BytesToAddress(a1)
		a2 := common.Hex2Bytes(testAddress2)
		address2 := common.BytesToAddress(a2)
		a3 := common.Hex2Bytes(testAddress3)
		address3 := common.BytesToAddress(a3)

		var users []common.Address
		users = append(users, address1, address2, address3)

		contract := &Contract{
			metricDataMutex:  sync.RWMutex{},
			users:            users,
			heightLowBounder: 0,
		}

		// user removed, have to clean up and update user set.
		users = users[:len(users)-1]

		contract.cleanUselessMetrics(users, uint64(10))
		if len(contract.users) != 2 {
			t.Fatal("case failed.")
		}
	})
}

func TestContract_MeasureMetricsOfNetworkEconomic(t *testing.T) {
	// to do a mock.
}

func TestContract_removeMetricsOutOfWindow(t *testing.T) {
	t.Run("remove metrics which is out of window, normal case: metrics height in window.", func(t *testing.T) {
		contract := &Contract{
			metricDataMutex:  sync.RWMutex{},
			users:            nil,
			heightLowBounder: 0,
		}

		contract.removeMetricsOutOfWindow(BlockRewardHeightWindow - 1)
		if contract.heightLowBounder != 0 {
			t.Fatal("test case failed.")
		}
	})

	t.Run("remove metrics which is out of window, normal case: metrics height is out of window.", func(t *testing.T) {
		a1 := common.Hex2Bytes(testAddress1)
		address1 := common.BytesToAddress(a1)
		a2 := common.Hex2Bytes(testAddress2)
		address2 := common.BytesToAddress(a2)
		a3 := common.Hex2Bytes(testAddress3)
		address3 := common.BytesToAddress(a3)

		var users []common.Address
		users = append(users, address1, address2, address3)

		contract := &Contract{
			metricDataMutex:  sync.RWMutex{},
			users:            users,
			heightLowBounder: 0,
		}

		contract.removeMetricsOutOfWindow(BlockRewardHeightWindow)
		if contract.heightLowBounder != BlockRewardHeightWindowStepRange {
			t.Fatal("case failed.")
		}
	})
}

func TestContract_measureRewardDistributionMetrics(t *testing.T) {
	t.Run("measure reward distribution metrics, exception case: wrong parameter.", func(t *testing.T) {
		contract := &Contract{}
		a1 := common.Hex2Bytes(testAddress1)
		address1 := common.BytesToAddress(a1)
		a2 := common.Hex2Bytes(testAddress2)
		address2 := common.BytesToAddress(a2)
		a3 := common.Hex2Bytes(testAddress3)
		address3 := common.BytesToAddress(a3)

		var stakeHolders []common.Address
		stakeHolders = append(stakeHolders, address1, address2, address3)
		var rewardFractions []*big.Int
		rewardFractions = append(rewardFractions, common.Big1, common.Big2)
		blockReward := common.Big32
		contract.measureRewardDistributionMetrics(stakeHolders, rewardFractions, blockReward, BlockRewardHeightWindow)
		if contract.heightLowBounder != 0 {
			t.Fatal("case failed.")
		}
	})

	t.Run("measure reward distribution metrics, normal case.", func(t *testing.T) {
		contract := &Contract{}
		a1 := common.Hex2Bytes(testAddress1)
		address1 := common.BytesToAddress(a1)
		a2 := common.Hex2Bytes(testAddress2)
		address2 := common.BytesToAddress(a2)
		a3 := common.Hex2Bytes(testAddress3)
		address3 := common.BytesToAddress(a3)

		var stakeHolders []common.Address
		stakeHolders = append(stakeHolders, address1, address2, address3)
		var rewardFractions []*big.Int
		rewardFractions = append(rewardFractions, common.Big1, common.Big2, common.Big3)
		blockReward := common.Big32
		contract.measureRewardDistributionMetrics(stakeHolders, rewardFractions, blockReward, BlockRewardHeightWindow)
		if contract.heightLowBounder != BlockRewardHeightWindowStepRange {
			t.Fatal("case failed.")
		}
	})
}

/*

type evmMock struct{}

func (evmMock) Call(caller vm.ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	return
}

*/
