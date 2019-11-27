package autonity

import (
	"errors"
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/metrics"
	"math/big"
	"sync"
)

const (
	Participant uint8 = iota
	Stakeholder
	Validator
)

const (
	/*
		gauge metrics which tracks stake, balance, and commissionrate of per user, when user is removed, metric
		should be removed from memory too.
		contract/user/0xefqefea...214dafaff/validator/stake
		contract/user/0xefqefea...214dafaff/stakeholder/stake
		contract/user/0xefqefea...214dafaff/participant/stake
		contract/user/0xefqefea...214dafaff/validator/balance
		contract/user/0xefqefea...214dafaff/stakeholder/balance
		contract/user/0xefqefea...214dafaff/participant/balance
		contract/user/0xefqefea...214dafaff/validator/commissionrate
		contract/user/0xefqefea...214dafaff/stakeholder/commissionrate
		contract/user/0xefqefea...214dafaff/participant/commissionrate
		template: contract/user/common.address/[validator|stakeholder|participant]/[stake|balance|commissionrate]
	*/
	UserMetricIDTemplate = "contract/user/%s/%s/%s"

	/*
		counter metrics which counts reward distribution for per block, cannot hold these counters from block0 to
		infinite block number in memory, so we apply a height/time window to keep the counters in reasonable range, for
		those counter which reported to TSDB could be removed for memory recycle.
		template: contract/block/number/user/common.address/[validator|stakeholder|participant]/reward
	*/
	BlockRewardHeightWindow          = 3600 // 1 hour time window to keep those counters in memory.
	BlockRewardHeightWindowStepRange = 600  // each 10 minutes to shrink the window.

	BlockRewardDistributionMetricIDTemplate = "contract/block/%v/user/%s/%s/reward"

	// counter tracks the reward/transactionfee of a specific block
	BlockRewardBlockMetricID = "contract/block/%v/reward"

	// counter counts SUM of the rewards for each block in the history.
	BlockRewardSUMMetricID = "contract/blockreward/sum"

	// gauge metrics which track the global level metrics of economic.
	GlobalMetricIDGasPrice    = "contract/global/mingasprice"
	GlobalMetricIDStakeSupply = "contract/global/stakesupply"

	// gauge metrics which track the network operator balance.
	GlobalOperatorBalanceMetricID = "contract/global/operator/balance"

	RoleUnknown               = "unknown"
	RoleValidator             = "validator"
	RoleStakeHolder           = "stakeholder"
	RoleParticipant           = "participant"
)

// refer to autonity contract abt spec, keep in same meta.
type EconomicMetaData struct {
	Accounts        []common.Address `abi:"accounts"`
	Usertypes       []uint8          `abi:"usertypes"`
	Stakes          []*big.Int       `abi:"stakes"`
	Commissionrates []*big.Int       `abi:"commissionrates"`
	Mingasprice     *big.Int         `abi:"mingasprice"`
	Stakesupply     *big.Int         `abi:"stakesupply"`
}

// refer to autonity contract abi spec, keep in same meta.
type RewardDistributionMetaData struct {
	Result          bool             `abi:"result"`
	Holders         []common.Address `abi:"stakeholders"`
	Rewardfractions []*big.Int       `abi:"rewardfractions"`
	Amount          *big.Int         `abi:"amount"`
}

type EconomicMetrics struct {
	metricDataMutex  sync.RWMutex
	users            []common.Address
	heightLowBounder uint64 // time/height window for keeping reasonable number of metrics in registry.
}

// measure metrics of user's meta data by regarding of network economic.
func (em *EconomicMetrics) SubmitEconomicMetrics(v *EconomicMetaData, stateDB *state.StateDB, height uint64, operator common.Address) {

	if v == nil || stateDB == nil {
		return
	}
	// measure global metrics
	gasPriceGauge := metrics.GetOrRegisterGauge(GlobalMetricIDGasPrice, nil)
	stakeTotalSupplyGauge := metrics.GetOrRegisterGauge(GlobalMetricIDStakeSupply, nil)
	operatorBalanceGauge := metrics.GetOrRegisterGauge(GlobalOperatorBalanceMetricID, nil)
	gasPriceGauge.Update(v.Mingasprice.Int64())
	stakeTotalSupplyGauge.Update(v.Stakesupply.Int64())
	operatorBalanceGauge.Update(stateDB.GetBalance(operator).Int64())

	// measure user metrics
	if len(v.Accounts) != len(v.Usertypes) || len(v.Accounts) != len(v.Stakes) ||
		len(v.Accounts) != len(v.Commissionrates) {
		log.Warn("mismatched data set dumped from autonity contract")
		return
	}

	for i := 0; i < len(v.Accounts); i++ {
		user := v.Accounts[i]
		userType := v.Usertypes[i]
		stake := v.Stakes[i]
		rate := v.Commissionrates[i]
		balance := stateDB.GetBalance(user)

		log.Debug("user: ", user, "userType: ", userType, "stake: ", stake, "rate: ", rate, "balance: ", balance)

		// generate metric ID.
		stakeID, balanceID, commmissionRateID, err := em.generateUserMetricsID(user, userType)
		if err != nil {
			log.Warn("generateUserMetricsID failed.")
			return
		}

		// get or create metrics from default registry.
		stakeGauge := metrics.GetOrRegisterGauge(stakeID, nil)
		balanceGauge := metrics.GetOrRegisterGauge(balanceID, nil)
		commissionRateGauge := metrics.GetOrRegisterGauge(commmissionRateID, nil)

		// submit data to registry.
		stakeGauge.Update(stake.Int64())
		balanceGauge.Update(balance.Int64())
		commissionRateGauge.Update(rate.Int64())
	}

	// clean up useless metrics if there exists.
	em.cleanUselessMetrics(v.Accounts, height)
}

func (em *EconomicMetrics) SubmitRewardDistributionMetrics(v *RewardDistributionMetaData, height uint64) {
	if len(v.Holders) != len(v.Rewardfractions) {
		log.Warn("reward fractions does not distribute to all stake holder.")
		return
	}

	// submit reward distribution metrics to registry.
	for i := 0; i < len(v.Holders); i++ {
		rewardDistributionMetricID := em.generateRewardDistributionMetricsID(v.Holders[i], Stakeholder, height)
		rwdDistributionMetric := metrics.GetOrRegisterCounter(rewardDistributionMetricID, nil)
		rwdDistributionMetric.Inc(v.Rewardfractions[i].Int64())
	}

	// submit block reward metric to registry.
	blockRewardMetricID := em.generateBlockRewardMetricsID(height)
	blockRewardMetric := metrics.GetOrRegisterCounter(blockRewardMetricID, nil)
	blockRewardMetric.Inc(v.Amount.Int64())

	// submit block reward sum metrics to registry.
	sumBlockRewardMetric := metrics.GetOrRegisterCounter(BlockRewardSUMMetricID, nil)
	sumBlockRewardMetric.Inc(v.Amount.Int64())

	// check to remove reward distribution metrics which is out of time/height window.
	em.removeMetricsOutOfWindow(height)
}

func (em *EconomicMetrics) generateBlockRewardMetricsID(blockNumber uint64) string {
	return fmt.Sprintf(BlockRewardBlockMetricID, blockNumber)
}

func (em *EconomicMetrics) generateRewardDistributionMetricsID(address common.Address, role uint8, blockNumber uint64) string {
	userType := em.resolveUserTypeName(role)
	blockMetricsID := fmt.Sprintf(BlockRewardDistributionMetricIDTemplate, blockNumber, address.String(), userType)
	return blockMetricsID
}

func (em *EconomicMetrics) resolveUserTypeName(role uint8) string {
	ret := RoleUnknown
	switch role {
	case Validator:
		ret = RoleValidator
	case Stakeholder:
		ret = RoleStakeHolder
	case Participant:
		ret = RoleParticipant
	}
	return ret
}

func (em *EconomicMetrics) generateUserMetricsID(address common.Address, role uint8) (stakeID string,
	balanceID string, commissionRateID string, err error) {
	if role > Validator {
		return "", "", "", errors.New("invalid parameter")
	}
	userType := em.resolveUserTypeName(role)
	stakeID = fmt.Sprintf(UserMetricIDTemplate, address.String(), userType, "stake")
	balanceID = fmt.Sprintf(UserMetricIDTemplate, address.String(), userType, "balance")
	commissionRateID = fmt.Sprintf(UserMetricIDTemplate, address.String(), userType, "commissionrate")
	return stakeID, balanceID, commissionRateID, nil
}

func (em *EconomicMetrics) removeMetricsFromRegistry(user common.Address, blockNumber uint64) {

	// clean up metrics which counts user's stake, balance and commission rate.
	for role := Participant; role <= Validator; role++ {
		if stakeID, balanceID, commissionRateID, err := em.generateUserMetricsID(user, role); err == nil {
			metrics.DefaultRegistry.Unregister(stakeID)
			metrics.DefaultRegistry.Unregister(balanceID)
			metrics.DefaultRegistry.Unregister(commissionRateID)
		}
	}
	// clean up metrics which counts the removed user's reward.
	for height := em.heightLowBounder; height <= blockNumber; height++ {
		rewardDistributionMetricID := em.generateRewardDistributionMetricsID(user, Stakeholder, blockNumber)
		metrics.DefaultRegistry.Unregister(rewardDistributionMetricID)
	}
}

/*
*  cleanUselessMetrics clean up metric memory from ETH-Metric framework by removed users.
*  Note: when node restart, those metrics registered in the metric registry are auto released.
 */
func (em *EconomicMetrics) cleanUselessMetrics(addresses []common.Address, blockNumber uint64) {
	if len(addresses) == 0 {
		return
	}
	em.metricDataMutex.Lock()
	defer em.metricDataMutex.Unlock()

	if em.users == nil || len(em.users) == 0 {
		em.users = addresses
		return
	}

	for _, user := range em.users {
		found := false
		for _, address := range addresses {
			if user == address {
				found = true
				break
			}
		}

		if !found {
			// to clean up metrics of users who was removed.
			em.removeMetricsFromRegistry(user, blockNumber)
		}
	}
	// load the latest user set from economic contract.
	em.users = addresses
}

/*
* removeMetricsOutOfWindow remove those metrics from memory which is out of window.
 */
func (em *EconomicMetrics) removeMetricsOutOfWindow(blockNumber uint64) {
	em.metricDataMutex.Lock()
	defer em.metricDataMutex.Unlock()
	if blockNumber-em.heightLowBounder < BlockRewardHeightWindow {
		return
	}

	// newLowBounder := blockNumber - ac.heightLowBounder
	newLowBounder := em.heightLowBounder + BlockRewardHeightWindowStepRange
	for height := em.heightLowBounder; height < newLowBounder; height++ {
		for _, user := range em.users {
			blcRwdDistributionID := em.generateRewardDistributionMetricsID(user, Stakeholder, height)
			metrics.DefaultRegistry.Unregister(blcRwdDistributionID)
		}
		blcRwdID := em.generateBlockRewardMetricsID(height)
		metrics.DefaultRegistry.Unregister(blcRwdID)
	}
	// update low bounder with new window edge.
	em.heightLowBounder = newLowBounder
}
