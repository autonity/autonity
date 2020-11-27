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

func TestEconomicMetrics_generateMetricsIDs(t *testing.T) {
	t.Run("test generate block reward metric ID", func(t *testing.T) {
		em := &EconomicMetrics{}
		blockNumber := uint64(2)
		blockRewardMetricID := em.generateBlockRewardMetricsID(blockNumber)
		if blockRewardMetricID != fmt.Sprintf(BlockRewardBlockMetricID, blockNumber) {
			t.Fatal("wrong result expected.")
		}
	})

	t.Run("test generate reward distribution metric ID", func(t *testing.T) {
		em := &EconomicMetrics{}
		a := common.Hex2Bytes(testAddress1)
		address := common.BytesToAddress(a)
		metricID := em.generateRewardDistributionMetricsID(address, Participant, uint64(2))
		expectedID := fmt.Sprintf("contract/block/2/user/%s/participant/reward", address.String())
		if metricID != expectedID {
			t.Fatal("case failed.")
		}
	})

	t.Run("test generate user metrics ID", func(t *testing.T) {
		em := &EconomicMetrics{}
		a := common.Hex2Bytes(testAddress1)
		address := common.BytesToAddress(a)
		stakeID, balanceID, _ := em.generateUserMetricsID(address, Participant)
		expectedStakeID := fmt.Sprintf(UserMetricIDTemplate, address.String(), "participant", "stake")
		expectedBalanceID := fmt.Sprintf(UserMetricIDTemplate, address.String(), "participant", "balance")
		if stakeID != expectedStakeID || balanceID != expectedBalanceID {
			t.Fatal("test case failed.")
		}
	})

	t.Run("test generate user metrics ID with wrong role", func(t *testing.T) {
		em := &EconomicMetrics{}
		a := common.Hex2Bytes(testAddress1)
		address := common.BytesToAddress(a)
		_, _, err := em.generateUserMetricsID(address, 3)
		if err == nil {
			t.Fatal("test case failed.")
		}
	})

	t.Run("test resolve user type name", func(t *testing.T) {
		em := &EconomicMetrics{}
		name := em.resolveUserTypeName(Participant)
		if name != "participant" {
			t.Fatal("case failed.")
		}

		name = em.resolveUserTypeName(Stakeholder)
		if name != "stakeholder" {
			t.Fatal("case failed.")
		}

		name = em.resolveUserTypeName(Validator)
		if name != "validator" {
			t.Fatal("case failed")
		}

		name = em.resolveUserTypeName(4)
		if name != "unknown" {
			t.Fatal("case failed.")
		}
	})
}

func TestEconomicMetrics_removeMetricsFromRegistry(t *testing.T) {
	t.Run("remove user metrics from metric registry", func(t *testing.T) {
		// prepare context in metric registry
		em := &EconomicMetrics{}
		blockHeight := uint64(10)
		a1 := common.Hex2Bytes(testAddress1)
		address1 := common.BytesToAddress(a1)
		a2 := common.Hex2Bytes(testAddress2)
		address2 := common.BytesToAddress(a2)
		a3 := common.Hex2Bytes(testAddress3)
		address3 := common.BytesToAddress(a3)

		stakeID1, balanceID1, _ := em.generateUserMetricsID(address1, Participant)
		rewardDistributionMetricID1 := em.generateRewardDistributionMetricsID(address1, Stakeholder, blockHeight)

		metrics.GetOrRegisterCounter(rewardDistributionMetricID1, nil).Inc(100)
		metrics.GetOrRegisterGauge(stakeID1, nil).Update(100)
		metrics.GetOrRegisterGauge(balanceID1, nil).Update(100)

		stakeID2, balanceID2, _ := em.generateUserMetricsID(address2, Stakeholder)
		rewardDistributionMetricID2 := em.generateRewardDistributionMetricsID(address2, Stakeholder, blockHeight)

		metrics.GetOrRegisterCounter(rewardDistributionMetricID2, nil).Inc(200)
		metrics.GetOrRegisterGauge(stakeID2, nil).Update(200)
		metrics.GetOrRegisterGauge(balanceID2, nil).Update(200)

		stakeID3, balanceID3, _ := em.generateUserMetricsID(address3, Validator)
		rewardDistributionMetricID3 := em.generateRewardDistributionMetricsID(address3, Stakeholder, blockHeight)

		metrics.GetOrRegisterCounter(rewardDistributionMetricID3, nil).Inc(300)
		metrics.GetOrRegisterGauge(stakeID3, nil).Update(300)
		metrics.GetOrRegisterGauge(balanceID3, nil).Update(300)

		em.removeMetricsFromRegistry(address1, blockHeight)
		em.removeMetricsFromRegistry(address2, blockHeight)
		em.removeMetricsFromRegistry(address3, blockHeight)

		if metrics.Get(stakeID1) != nil || metrics.Get(balanceID1) != nil ||
			metrics.Get(rewardDistributionMetricID1) != nil {
			t.Fatal("case failed.")
		}

		if metrics.Get(stakeID2) != nil || metrics.Get(balanceID2) != nil ||
			metrics.Get(rewardDistributionMetricID2) != nil {
			t.Fatal("case failed.")
		}

		if metrics.Get(stakeID3) != nil || metrics.Get(balanceID3) != nil ||
			metrics.Get(rewardDistributionMetricID3) != nil {
			t.Fatal("case failed.")
		}

	})
}

func TestEconomicMetrics_cleanUselessMetrics(t *testing.T) {
	t.Run("clean up metrics for removed users, exception case: input address set is empty.", func(t *testing.T) {
		// prepare context in metric registry
		em := &EconomicMetrics{}
		blockHeight := uint64(10)
		em.cleanUselessMetrics(nil, blockHeight)
	})

	t.Run("clean up metrics for removed users, exception case: local address set is empty.", func(t *testing.T) {
		em := &EconomicMetrics{}

		a1 := common.Hex2Bytes(testAddress1)
		address1 := common.BytesToAddress(a1)
		a2 := common.Hex2Bytes(testAddress2)
		address2 := common.BytesToAddress(a2)
		a3 := common.Hex2Bytes(testAddress3)
		address3 := common.BytesToAddress(a3)

		var users []common.Address
		users = append(users, address1, address2, address3)
		em.cleanUselessMetrics(users, uint64(10))
		if len(em.users) != 3 {
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

		em := &EconomicMetrics{
			metricDataMutex:  sync.RWMutex{},
			users:            users,
			heightLowBounder: 0,
		}

		// user removed, have to clean up and update user set.
		users = users[:len(users)-1]

		em.cleanUselessMetrics(users, uint64(10))
		if len(em.users) != 2 {
			t.Fatal("case failed.")
		}
	})
}

func TestEconomicMetrics_removeMetricsOutOfWindow(t *testing.T) {
	t.Run("remove metrics which is out of window, normal case: metrics height in window.", func(t *testing.T) {
		em := &EconomicMetrics{
			metricDataMutex:  sync.RWMutex{},
			users:            nil,
			heightLowBounder: 0,
		}

		em.removeMetricsOutOfWindow(BlockRewardHeightWindow - 1)
		if em.heightLowBounder != 0 {
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

		em := &EconomicMetrics{
			metricDataMutex:  sync.RWMutex{},
			users:            users,
			heightLowBounder: 0,
		}

		em.removeMetricsOutOfWindow(BlockRewardHeightWindow)
		if em.heightLowBounder != BlockRewardHeightWindowStepRange {
			t.Fatal("case failed.")
		}
	})
}

func TestEconomicMetrics_measureRewardDistributionMetrics(t *testing.T) {
	t.Run("measure reward distribution metrics, exception case: wrong parameter.", func(t *testing.T) {
		em := &EconomicMetrics{}
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
		var distributions RewardDistributionMetaData
		distributions.Amount = blockReward
		distributions.Rewardfractions = rewardFractions
		distributions.Holders = stakeHolders
		distributions.Result = true
		em.SubmitRewardDistributionMetrics(&distributions, BlockRewardHeightWindow)
		if em.heightLowBounder != 0 {
			t.Fatal("case failed.")
		}
	})

	t.Run("measure reward distribution metrics, normal case.", func(t *testing.T) {
		em := &EconomicMetrics{}
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
		var distributions RewardDistributionMetaData
		distributions.Amount = blockReward
		distributions.Rewardfractions = rewardFractions
		distributions.Holders = stakeHolders
		distributions.Result = true

		em.SubmitRewardDistributionMetrics(&distributions, BlockRewardHeightWindow)
		if em.heightLowBounder != BlockRewardHeightWindowStepRange {
			t.Fatal("case failed.")
		}
	})
}

func TestEconomicMetrics_recordMetric(t *testing.T) {
	t.Run("record metrics, normal case 1.", func(t *testing.T) {
		em := &EconomicMetrics{}
		value := big.NewInt(0)
		em.recordMetric("metricID1", value, true)
		metric := metrics.Get("metricID1")
		if metric == nil {
			t.Fatal("case failed.")
		}
	})

	t.Run("record metrics, normal case 2.", func(t *testing.T) {
		em := &EconomicMetrics{}
		value := big.NewInt(0)
		em.recordMetric("metricID2", value, false)
		metric := metrics.Get("metricID2")
		if metric == nil {
			t.Fatal("case failed.")
		}
	})

	t.Run("record metrics, exception case.", func(t *testing.T) {
		em := &EconomicMetrics{}
		em.recordMetric("metricID3", nil, false)
		metric := metrics.Get("metricID3")
		if metric != nil {
			t.Fatal("case failed.")
		}
	})

}
