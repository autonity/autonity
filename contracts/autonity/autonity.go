package autonity

import (
	"errors"
	"fmt"
	"github.com/clearmatics/autonity/accounts/abi"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/metrics"
	"github.com/clearmatics/autonity/params"
	"math/big"
	"reflect"
	"sort"
	"strings"
	"sync"
)

func NewAutonityContract(
	bc Blockchainer,
	canTransfer func(db vm.StateDB, addr common.Address, amount *big.Int) bool,
	transfer func(db vm.StateDB, sender, recipient common.Address, amount *big.Int),
	GetHashFn func(ref *types.Header, chain ChainContext) func(n uint64) common.Hash,
) *Contract {
	return &Contract{
		bc:          bc,
		canTransfer: canTransfer,
		transfer:    transfer,
		GetHashFn:   GetHashFn,
		//SavedValidatorsRetriever: SavedValidatorsRetriever,
	}
}

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
	GloablMetricIDStakeSupply = "contract/global/stakesupply"
	RoleUnknown               = "unknown"
	RoleValidator             = "validator"
	RoleStakeHolder           = "stakeholder"
	RoleParticipant           = "participant"
)

type ChainContext interface {
	// Engine retrieves the chain's consensus engine.
	Engine() consensus.Engine

	// GetHeader returns the hash corresponding to their hash.
	GetHeader(common.Hash, uint64) *types.Header
}
type Blockchainer interface {
	ChainContext
	GetVMConfig() *vm.Config
	Config() *params.ChainConfig

	UpdateEnodeWhitelist(newWhitelist *types.Nodes)
	ReadEnodeWhitelist(openNetwork bool) *types.Nodes
}

type Contract struct {
	address                  common.Address
	contractABI              *abi.ABI
	bc                       Blockchainer
	SavedValidatorsRetriever func(i uint64) ([]common.Address, error)

	metricDataMutex  sync.RWMutex
	users            []common.Address
	heightLowBounder uint64 // time/height window for keeping reasonable number of metrics in registry.

	canTransfer func(db vm.StateDB, addr common.Address, amount *big.Int) bool
	transfer    func(db vm.StateDB, sender, recipient common.Address, amount *big.Int)
	GetHashFn   func(ref *types.Header, chain ChainContext) func(n uint64) common.Hash
	sync.RWMutex
}

func (ac *Contract) generateBlockRewardMetricsID(blockNumber uint64) string {
	return fmt.Sprintf(BlockRewardBlockMetricID, blockNumber)
}

func (ac *Contract) generateRewardDistributionMetricsID(address common.Address, role uint8, blockNumber uint64) string {
	userType := ac.resolveUserTypeName(role)
	blockMetricsID := fmt.Sprintf(BlockRewardDistributionMetricIDTemplate, blockNumber, address.String(), userType)
	return blockMetricsID
}

func (ac *Contract) resolveUserTypeName(role uint8) string {
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

func (ac *Contract) generateUserMetricsID(address common.Address, role uint8) (stakeID string,
	balanceID string, commissionRateID string, err error) {
	if role > Validator {
		return "", "", "", errors.New("invalid parameter")
	}
	userType := ac.resolveUserTypeName(role)
	stakeID = fmt.Sprintf(UserMetricIDTemplate, address.String(), userType, "stake")
	balanceID = fmt.Sprintf(UserMetricIDTemplate, address.String(), userType, "balance")
	commissionRateID = fmt.Sprintf(UserMetricIDTemplate, address.String(), userType, "commissionrate")
	return stakeID, balanceID, commissionRateID, nil
}

func (ac *Contract) removeMetricsFromRegistry(user common.Address, blockNumber uint64) {

	// clean up metrics which counts user's stake, balance and commission rate.
	for role := Participant; role <= Validator; role++ {
		if stakeID, balanceID, commissionRateID, err := ac.generateUserMetricsID(user, role); err == nil {
			metrics.DefaultRegistry.Unregister(stakeID)
			metrics.DefaultRegistry.Unregister(balanceID)
			metrics.DefaultRegistry.Unregister(commissionRateID)
		}
	}
	// clean up metrics which counts the removed user's reward.
	for height := ac.heightLowBounder; height <= blockNumber; height++ {
		rewardDistributionMetricID := ac.generateRewardDistributionMetricsID(user, Stakeholder, blockNumber)
		metrics.DefaultRegistry.Unregister(rewardDistributionMetricID)
	}
}

/*
*  CleanUselessMetrics clean up metric memory from ETH-Metric framework by removed users.
*  Note: when node restart, those metrics registered in the metric registry are auto released.
 */
func (ac *Contract) CleanUselessMetrics(addresses []common.Address, blockNumber uint64) {
	if len(addresses) == 0 {
		return
	}
	ac.metricDataMutex.Lock()
	defer ac.metricDataMutex.Unlock()

	if ac.users == nil || len(ac.users) == 0 {
		ac.users = addresses
		return
	}

	for _, user := range ac.users {
		found := false
		for _, address := range addresses {
			if user == address {
				found = true
				break
			}
		}

		if !found {
			// to clean up metrics of users who was removed.
			ac.removeMetricsFromRegistry(user, blockNumber)
		}
	}
	// load the latest user set from economic contract.
	ac.users = addresses
}

// measure metrics of user's meta data by regarding of network economic.
func (ac *Contract) MeasureMetricsOfNetworkEconomic(header *types.Header, stateDB *state.StateDB) {
	if header == nil || stateDB == nil || header.Number.Uint64() < 1 {
		return
	}

	// prepare abi and evm context
	deployer := ac.bc.Config().AutonityContractConfig.Deployer
	sender := vm.AccountRef(deployer)
	gas := uint64(0xFFFFFFFF)
	evm := ac.getEVM(header, deployer, stateDB)

	ABI, err := ac.abi()
	if err != nil {
		return
	}

	// pack the function which dump the data from contract.
	input, err := ABI.Pack("dumpNetworkEconomicsData")
	if err != nil {
		log.Warn("cannot pack the method: ", err.Error())
		return
	}

	// call evm.
	value := new(big.Int).SetUint64(0x00)
	ret, _, vmerr := evm.Call(sender, ac.Address(), input, gas, value)
	log.Debug("bytes return from contract: ", ret)
	if vmerr != nil {
		log.Warn("Error Autonity Contract dumpNetworkEconomics")
		return
	}

	// marshal the data from bytes arrays into specified structure.
	v := struct {
		Accounts        []common.Address `abi:"accounts"`
		Usertypes       []uint8          `abi:"usertypes"`
		Stakes          []*big.Int       `abi:"stakes"`
		Commissionrates []*big.Int       `abi:"commissionrates"`
		Mingasprice     *big.Int         `abi:"mingasprice"`
		Stakesupply     *big.Int         `abi:"stakesupply"`
	}{make([]common.Address, 32), make([]uint8, 32), make([]*big.Int, 32),
		make([]*big.Int, 32), new(big.Int), new(big.Int)}

	if err := ABI.Unpack(&v, "dumpNetworkEconomicsData", ret); err != nil { // can't work with aliased types
		log.Warn("Could not unpack dumpNetworkEconomicsData returned value", "err", err, "header.num",
			header.Number.Uint64())
		return
	}

	// measure global metrics
	gasPriceGauge := metrics.GetOrRegisterGauge(GlobalMetricIDGasPrice, nil)
	stakeTotalSupplyGauge := metrics.GetOrRegisterGauge(GloablMetricIDStakeSupply, nil)
	gasPriceGauge.Update(v.Mingasprice.Int64())
	stakeTotalSupplyGauge.Update(v.Stakesupply.Int64())

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
		stakeID, balanceID, commmissionRateID, err := ac.generateUserMetricsID(user, userType)
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
	ac.CleanUselessMetrics(v.Accounts, header.Number.Uint64())
}

//// Instantiates a new EVM object which is required when creating or calling a deployed contract
func (ac *Contract) getEVM(header *types.Header, origin common.Address, statedb *state.StateDB) *vm.EVM {
	coinbase, _ := types.Ecrecover(header)
	evmContext := vm.Context{
		CanTransfer: ac.canTransfer,
		Transfer:    ac.transfer,
		GetHash:     ac.GetHashFn(header, ac.bc),
		Origin:      origin,
		Coinbase:    coinbase,
		BlockNumber: header.Number,
		Time:        new(big.Int).SetUint64(header.Time),
		GasLimit:    header.GasLimit,
		Difficulty:  header.Difficulty,
		GasPrice:    new(big.Int).SetUint64(0x0),
	}
	vmConfig := *ac.bc.GetVMConfig()
	evm := vm.NewEVM(evmContext, statedb, ac.bc.Config(), vmConfig)
	return evm
}

// deployContract deploys the contract contained within the genesis field bytecode
func (ac *Contract) DeployAutonityContract(chain consensus.ChainReader, header *types.Header, statedb *state.StateDB) (common.Address, error) {
	// Convert the contract bytecode from hex into bytes
	contractBytecode := common.Hex2Bytes(chain.Config().AutonityContractConfig.Bytecode)
	evm := ac.getEVM(header, chain.Config().AutonityContractConfig.Deployer, statedb)
	sender := vm.AccountRef(chain.Config().AutonityContractConfig.Deployer)

	//todo do we need it?
	//validators, err = ac.SavedValidatorsRetriever(1)
	//sort.Sort(validators)

	//We need to append to data the constructor's parameters
	//That should always be genesis validators

	contractABI, err := ac.abi()

	if err != nil {
		log.Error("abi.JSON returns err", "err", err)
		return common.Address{}, err
	}

	ln := len(chain.Config().AutonityContractConfig.GetValidatorUsers())
	validators := make(common.Addresses, 0, ln)
	enodes := make([]string, 0, ln)
	accTypes := make([]*big.Int, 0, ln)
	participantStake := make([]*big.Int, 0, ln)
	for _, v := range chain.Config().AutonityContractConfig.Users {
		validators = append(validators, v.Address)
		enodes = append(enodes, v.Enode)
		accTypes = append(accTypes, big.NewInt(int64(v.Type.GetID())))
		participantStake = append(participantStake, big.NewInt(int64(v.Stake)))
	}

	constructorParams, err := contractABI.Pack("",
		validators,
		enodes,
		accTypes,
		participantStake,
		chain.Config().AutonityContractConfig.Operator,
		new(big.Int).SetUint64(chain.Config().AutonityContractConfig.MinGasPrice))
	if err != nil {
		log.Error("contractABI.Pack returns err", "err", err)
		return common.Address{}, err
	}

	data := append(contractBytecode, constructorParams...)
	gas := uint64(0xFFFFFFFF)
	value := new(big.Int).SetUint64(0x00)

	// Deploy the Autonity contract
	_, contractAddress, _, vmerr := evm.Create(sender, data, gas, value)
	if vmerr != nil {
		log.Error("evm.Create returns err", "err", vmerr)
		return contractAddress, vmerr
	}
	ac.Lock()
	ac.address = contractAddress
	ac.Unlock()
	log.Info("Deployed Autonity Contract", "Address", contractAddress.String())

	return contractAddress, nil
}

func (ac *Contract) ContractGetValidators(chain consensus.ChainReader, header *types.Header, statedb *state.StateDB) ([]common.Address, error) {
	if header.Number.Cmp(big.NewInt(1)) == 0 && ac.SavedValidatorsRetriever != nil {
		return ac.SavedValidatorsRetriever(1)
	}
	sender := vm.AccountRef(chain.Config().AutonityContractConfig.Deployer)
	gas := uint64(0xFFFFFFFF)
	evm := ac.getEVM(header, chain.Config().AutonityContractConfig.Deployer, statedb)
	contractABI, err := ac.abi()
	if err != nil {
		return nil, err
	}

	input, err := contractABI.Pack("getValidators")
	if err != nil {
		return nil, err
	}

	value := new(big.Int).SetUint64(0x00)
	//A standard call is issued - we leave the possibility to modify the state
	ret, _, vmerr := evm.Call(sender, ac.Address(), input, gas, value)
	if vmerr != nil {
		return nil, vmerr
	}

	var addresses []common.Address
	if err := contractABI.Unpack(&addresses, "getValidators", ret); err != nil { // can't work with aliased types
		log.Error("Could not unpack getValidators returned value", "err", err)
		return nil, err
	}

	sortableAddresses := common.Addresses(addresses)
	sort.Sort(sortableAddresses)
	return sortableAddresses, nil
}

var ErrAutonityContract = errors.New("could not call Autonity contract")

func (ac *Contract) UpdateEnodesWhitelist(state *state.StateDB, block *types.Block) error {
	newWhitelist, err := ac.GetWhitelist(block, state)
	if err != nil {
		log.Error("could not call contract", "err", err)
		return ErrAutonityContract
	}

	ac.bc.UpdateEnodeWhitelist(newWhitelist)
	return nil
}

func (ac *Contract) GetWhitelist(block *types.Block, db *state.StateDB) (*types.Nodes, error) {
	var (
		newWhitelist *types.Nodes
		err          error
	)

	if block.Number().Uint64() == 1 {
		// use genesis block whitelist
		newWhitelist = ac.bc.ReadEnodeWhitelist(false)
	} else {
		// call retrieveWhitelist contract function
		newWhitelist, err = ac.callGetWhitelist(db, block.Header())
	}

	return newWhitelist, err
}

//blockchain

func (ac *Contract) callGetWhitelist(state *state.StateDB, header *types.Header) (*types.Nodes, error) {
	// Needs to be refactored somehow
	deployer := ac.bc.Config().AutonityContractConfig.Deployer
	sender := vm.AccountRef(deployer)
	gas := uint64(0xFFFFFFFF)
	evm := ac.getEVM(header, deployer, state)

	ABI, err := ac.abi()
	if err != nil {
		return nil, err
	}

	input, err := ABI.Pack("getWhitelist")
	if err != nil {
		return nil, err
	}

	ret, _, vmerr := evm.StaticCall(sender, ac.Address(), input, gas)
	if vmerr != nil {
		log.Error("Error Autonity Contract getWhitelist()")
		return nil, vmerr
	}

	var returnedEnodes []string
	if err := ABI.Unpack(&returnedEnodes, "getWhitelist", ret); err != nil { // can't work with aliased types
		log.Error("Could not unpack getWhitelist returned value")
		return nil, err
	}

	return types.NewNodes(returnedEnodes, false), nil
}

func (ac *Contract) GetMinimumGasPrice(block *types.Block, db *state.StateDB) (uint64, error) {
	if block.Number().Uint64() <= 1 {
		return ac.bc.Config().AutonityContractConfig.MinGasPrice, nil
	}

	return ac.callGetMinimumGasPrice(db, block.Header())
}

func (ac *Contract) SetMinimumGasPrice(block *types.Block, db *state.StateDB, price *big.Int) error {
	if block.Number().Uint64() <= 1 {
		return nil
	}

	return ac.callSetMinimumGasPrice(db, block.Header(), price)
}

func (ac *Contract) callGetMinimumGasPrice(state *state.StateDB, header *types.Header) (uint64, error) {
	// Needs to be refactored somehow
	deployer := ac.bc.Config().AutonityContractConfig.Deployer
	sender := vm.AccountRef(deployer)
	gas := uint64(0xFFFFFFFF)
	evm := ac.getEVM(header, deployer, state)

	ABI, err := ac.abi()
	if err != nil {
		return 0, err
	}

	input, err := ABI.Pack("getMinimumGasPrice")
	if err != nil {
		return 0, err
	}

	value := new(big.Int).SetUint64(0x00)
	ret, _, vmerr := evm.Call(sender, ac.Address(), input, gas, value)
	if vmerr != nil {
		log.Error("Error Autonity Contract getMinimumGasPrice()")
		return 0, vmerr
	}

	minGasPrice := new(big.Int)
	if err := ABI.Unpack(&minGasPrice, "getMinimumGasPrice", ret); err != nil { // can't work with aliased types
		log.Error("Could not unpack minGasPrice returned value", "err", err, "header.num", header.Number.Uint64())
		return 0, err
	}

	return minGasPrice.Uint64(), nil
}

func (ac *Contract) callSetMinimumGasPrice(state *state.StateDB, header *types.Header, price *big.Int) error {
	// Needs to be refactored somehow
	deployer := ac.bc.Config().AutonityContractConfig.Deployer
	sender := vm.AccountRef(deployer)
	gas := uint64(0xFFFFFFFF)
	evm := ac.getEVM(header, deployer, state)

	ABI, err := ac.abi()
	if err != nil {
		return err
	}

	input, err := ABI.Pack("setMinimumGasPrice")
	if err != nil {
		return err
	}

	_, _, vmerr := evm.Call(sender, ac.Address(), input, gas, price)
	if vmerr != nil {
		log.Error("Error Autonity Contract getMinimumGasPrice()")
		return vmerr
	}
	return nil
}

func (ac *Contract) PerformRedistribution(header *types.Header, db *state.StateDB, gasUsed *big.Int) error {
	if header.Number.Uint64() <= 1 {
		return nil
	}
	return ac.callPerformRedistribution(db, header, gasUsed)
}

func (ac *Contract) callPerformRedistribution(state *state.StateDB, header *types.Header, blockGas *big.Int) error {
	// Needs to be refactored somehow
	deployer := ac.bc.Config().AutonityContractConfig.Deployer

	sender := vm.AccountRef(deployer)
	gas := uint64(0xFFFFFFFF)
	evm := ac.getEVM(header, deployer, state)

	ABI, err := ac.abi()
	if err != nil {
		return err
	}

	input, err := ABI.Pack("performRedistribution", blockGas)
	if err != nil {
		log.Error("Error Autonity Contract callPerformRedistribution()", "err", err)
		return err
	}

	value := new(big.Int).SetUint64(0x00)

	ret, _, vmerr := evm.Call(sender, ac.Address(), input, gas, value)
	if vmerr != nil {
		log.Error("Error Autonity Contract callPerformRedistribution()", "err", err)
		return vmerr
	}

	// after reward distribution, update metrics with the return values.
	v := struct {
		Holders         []common.Address `abi:"stakeholders"`
		Rewardfractions []*big.Int       `abi:"rewardfractions"`
		Amount          *big.Int         `abi:"amount"`
	}{make([]common.Address, 32), make([]*big.Int, 32), new(big.Int)}

	if err := ABI.Unpack(&v, "performRedistribution", ret); err != nil { // can't work with aliased types
		log.Error("Could not unpack performRedistribution returned value", "err", err, "header.num", header.Number.Uint64())
		return nil
	}

	ac.measureRewardDistributionMetrics(v.Holders, v.Rewardfractions, v.Amount, header.Number.Uint64())
	return nil
}

/*
* removeMetricsOutOfWindow remove those metrics from memory which is out of window.
 */
func (ac *Contract) removeMetricsOutOfWindow(blockNumber uint64) {
	ac.metricDataMutex.Lock()
	defer ac.metricDataMutex.Unlock()
	if blockNumber-ac.heightLowBounder < BlockRewardHeightWindow {
		return
	}

	// newLowBounder := blockNumber - ac.heightLowBounder
	newLowBounder := ac.heightLowBounder + BlockRewardHeightWindowStepRange
	for height := ac.heightLowBounder; height < newLowBounder; height++ {
		for _, user := range ac.users {
			blcRwdDistributionID := ac.generateRewardDistributionMetricsID(user, Stakeholder, height)
			metrics.DefaultRegistry.Unregister(blcRwdDistributionID)
		}
		blcRwdID := ac.generateBlockRewardMetricsID(height)
		metrics.DefaultRegistry.Unregister(blcRwdID)
	}
	// update low bounder with new window edge.
	ac.heightLowBounder = newLowBounder
}

func (ac *Contract) measureRewardDistributionMetrics(holders []common.Address, rewardFractions []*big.Int,
	blockReward *big.Int, blockNumber uint64) {
	if len(holders) != len(rewardFractions) {
		log.Warn("reward fractions does not distribute to all stake holder.")
		return
	}

	// submit reward distribution metrics to registry.
	for i := 0; i < len(holders); i++ {
		rewardDistributionMetricID := ac.generateRewardDistributionMetricsID(holders[i], Stakeholder, blockNumber)
		rwdDistributionMetric := metrics.GetOrRegisterCounter(rewardDistributionMetricID, nil)
		rwdDistributionMetric.Inc(rewardFractions[i].Int64())
	}

	// submit block reward metric to registry.
	blockRewardMetricID := ac.generateBlockRewardMetricsID(blockNumber)
	blockRewardMetric := metrics.GetOrRegisterCounter(blockRewardMetricID, nil)
	blockRewardMetric.Inc(blockReward.Int64())

	// submit block reward sum metrics to registry.
	sumBlockRewardMetric := metrics.GetOrRegisterCounter(BlockRewardSUMMetricID, nil)
	sumBlockRewardMetric.Inc(blockReward.Int64())

	// check to remove reward distribution metrics which is out of time/height window.
	ac.removeMetricsOutOfWindow(blockNumber)
}

func (ac *Contract) ApplyPerformRedistribution(transactions types.Transactions, receipts types.Receipts, header *types.Header, statedb *state.StateDB) error {
	log.Info("ApplyPerformRedistribution", "header", header.Number.Uint64())
	if header.Number.Cmp(big.NewInt(1)) < 1 {
		return nil
	}
	blockGas := new(big.Int)
	for i, tx := range transactions {
		blockGas.Add(blockGas, new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(receipts[i].GasUsed)))
	}

	ac.MeasureMetricsOfNetworkEconomic(header, statedb)

	log.Info("execution start ApplyPerformRedistribution", "balance", statedb.GetBalance(ac.Address()), "block", header.Number.Uint64(), "gas", blockGas.Uint64())
	if blockGas.Cmp(new(big.Int)) == 0 {
		log.Info("execution start ApplyPerformRedistribution with 0 gas", "balance", statedb.GetBalance(ac.Address()), "block", header.Number.Uint64())
		return nil
	}
	return ac.PerformRedistribution(header, statedb, blockGas)
}

func (ac *Contract) Address() common.Address {
	if reflect.DeepEqual(ac.address, common.Address{}) {
		addr, err := ac.bc.Config().AutonityContractConfig.GetContractAddress()
		if err != nil {
			log.Error("Cant get contract address", "err", err)
		}
		return addr
	}
	return ac.address
}

func (ac *Contract) abi() (*abi.ABI, error) {
	ac.Lock()
	defer ac.Unlock()
	if ac.contractABI != nil {
		return ac.contractABI, nil
	}
	ABI, err := abi.JSON(strings.NewReader(ac.bc.Config().AutonityContractConfig.ABI))
	if err != nil {
		return nil, err
	}
	ac.contractABI = &ABI
	return ac.contractABI, nil

}
