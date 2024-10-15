package autonity

import (
	"bytes"
	"errors"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/params/generated"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/ethdb"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
)

// "soft minimum" cache size for the proposer cache.
// "soft minimum" because we can have less entries in it, but once we reach `thresholdSize` we want to keep having at least `thresholdSize`.
// ensure it is >= accountability.accountabilityHeightRange to ensure reduced performance overhead in case of old proposal spamming
const thresholdSize = 256

var (
	AutonityABIKey = []byte("ABISPEC")
)

// Errors
var (
	ErrNoAutonityConfig = errors.New("autonity section missing from genesis")
)

// EVMProvider provides a new evm. This allows us to decouple the contract from *core.blockchain which is required to build a new evm.
type EVMProvider func(header *types.Header, origin common.Address, statedb vm.StateDB) *vm.EVM

// GenesisEVMProvider returns new, preconfigured EVM. Useful for genesis EVMs which depends on fixed set of parameters
// which can be closed in this function
type GenesisEVMProvider func(statedb vm.StateDB) *vm.EVM

type EVMContract struct {
	evmProvider EVMProvider
	contractABI *abi.ABI
	db          ethdb.Database
	chainConfig *params.ChainConfig

	sync.RWMutex
}

func NewEVMContract(evmProvider EVMProvider, contractABI *abi.ABI, db ethdb.Database, chainConfig *params.ChainConfig) *EVMContract {
	return &EVMContract{
		evmProvider: evmProvider,
		contractABI: contractABI,
		db:          db,
		chainConfig: chainConfig,
	}
}

type Cache struct {
	// minimum base fee
	minBaseFee    atomic.Pointer[big.Int]
	minBaseFeeCh  chan *AutonityMinimumBaseFeeUpdated
	subMinBaseFee event.Subscription

	// epoch period
	epochPeriod    atomic.Pointer[big.Int]
	epochPeriodCh  chan *AutonityEpochPeriodUpdated
	subEpochPeriod event.Subscription

	subscriptions *event.SubscriptionScope
	quit          chan struct{}
	done          chan struct{}
}

func newCache(ac *AutonityContract, head *types.Header, state vm.StateDB) (*Cache, error) {
	// initialize minimum base fee through contract call and subscribe to updated event
	minBaseFee, err := ac.MinimumBaseFee(head, state)
	if err != nil {
		return nil, err
	}
	minBaseFeeCh := make(chan *AutonityMinimumBaseFeeUpdated)
	subMinBaseFee, err := ac.WatchMinimumBaseFeeUpdated(nil, minBaseFeeCh)
	if err != nil {
		return nil, err
	}

	// initialize epoch period and subscribe to updated event
	epochPeriod, err := ac.EpochPeriod(head, state)
	if err != nil {
		return nil, err
	}
	epochPeriodCh := make(chan *AutonityEpochPeriodUpdated)
	subEpochPeriod, err := ac.WatchEpochPeriodUpdated(nil, epochPeriodCh)
	if err != nil {
		return nil, err
	}

	scope := new(event.SubscriptionScope)
	subMinBaseFeeWrapped := scope.Track(subMinBaseFee)
	subEpochPeriodWrapped := scope.Track(subEpochPeriod)
	cache := &Cache{
		minBaseFeeCh:   minBaseFeeCh,
		subMinBaseFee:  subMinBaseFeeWrapped,
		epochPeriodCh:  epochPeriodCh,
		subEpochPeriod: subEpochPeriodWrapped,
		subscriptions:  scope,
		quit:           make(chan struct{}),
		done:           make(chan struct{}),
	}

	// store initial values
	cache.minBaseFee.Store(minBaseFee)
	cache.epochPeriod.Store(epochPeriod)

	go cache.Listen()
	return cache, nil
}

//revive:disable:exported - Autonity is one of the contracts, so repetitive naming here is justified
type AutonityContract struct {
	EVMContract
	*AutonityFilterer                                     // allows to watch for Autonity Contract events
	proposers         map[uint64]map[int64]common.Address // map[height][round] --> proposer
}

type ProtocolContracts struct {
	*AutonityContract
	*Cache
	*Accountability
}

func NewProtocolContracts(config *params.ChainConfig, db ethdb.Database, provider EVMProvider, contractBackend bind.ContractBackend, head *types.Header, state vm.StateDB) (*ProtocolContracts, error) {
	if config.AutonityContractConfig == nil {
		return nil, ErrNoAutonityConfig
	}
	ABI := config.AutonityContractConfig.ABI
	// Retrieve a potentially updated ABI from the database
	if raw, err := rawdb.GetKeyValue(db, AutonityABIKey); err == nil || raw != nil {
		jsonABI, err := abi.JSON(bytes.NewReader(raw))
		if err != nil {
			return nil, err
		}
		ABI = &jsonABI
	}

	// create autonity EVM contract
	autonityFilterer, err := NewAutonityFilterer(params.AutonityContractAddress, contractBackend.(bind.ContractFilterer))
	if err != nil {
		return nil, err
	}
	autonityContract := &AutonityContract{
		EVMContract: EVMContract{
			evmProvider: provider,
			contractABI: ABI,
			db:          db,
			chainConfig: config,
		},
		AutonityFilterer: autonityFilterer,
		proposers:        make(map[uint64]map[int64]common.Address),
	}

	// initialize protocol contract cache
	cache, err := newCache(autonityContract, head, state)
	if err != nil {
		return nil, err
	}

	// bind to accountability contract
	accountabilityContract, err := NewAccountability(params.AccountabilityContractAddress, contractBackend)
	if err != nil {
		return nil, err
	}

	contract := ProtocolContracts{
		AutonityContract: autonityContract,
		Cache:            cache,
		Accountability:   accountabilityContract,
	}

	return &contract, nil
}

func (c *Cache) Listen() {
	defer func() {
		c.subscriptions.Close()
		close(c.done)
	}()

	for {
		select {
		case ev := <-c.minBaseFeeCh:
			c.minBaseFee.Store(ev.GasPrice)
		case ev := <-c.epochPeriodCh:
			c.epochPeriod.Store(ev.Period)
		// These should never happen. Errors from subscription can happen only if the subscription is done over an RPC connection.
		// In that case network errors can occur. Since everything is local here, no error should ever occur.
		// we crash the client to avoid using a out-of-date value in case something goes very wrong.
		case <-c.subEpochPeriod.Err():
			log.Crit("protocol contract cache out-of-sync. Please contact the Autonity team.")
		case <-c.subMinBaseFee.Err():
			log.Crit("protocol contract cache out-of-sync. Please contact the Autonity team.")
		case <-c.quit:
			return
		}
	}
}

func (c *Cache) Stop() {
	close(c.quit)
	<-c.done
}

func (c *Cache) MinimumBaseFee() *big.Int {
	return new(big.Int).Set(c.minBaseFee.Load())
}

// TODO: currently this value gets updated when the operator requests a change in the epoch period,
// however the new epoch period gets actually applied at epoch's end.
// for now it is not a big deal because the cached epoch period is only used to disconnect malicious signers
// but it needs to be fixed
func (c *Cache) EpochPeriod() *big.Int {
	return new(big.Int).Set(c.epochPeriod.Load())
}

func (c *AutonityContract) CommitteeEnodes(header *types.Header, db vm.StateDB, asACN bool) (*types.Nodes, error) {
	return c.callGetCommitteeEnodes(db, header, asACN)
}

func (c *AutonityContract) EpochByHeight(header *types.Header, db vm.StateDB, height *big.Int) (*types.EpochInfo, error) {
	return c.callEpochByHeight(db, header, height)
}

// EpochInfo gets the committee, the epoch boundaries and the omission failure delta based on the input header's state.
func (c *AutonityContract) EpochInfo(header *types.Header, db vm.StateDB) (*types.EpochInfo, error) {
	return c.callGetEpochInfo(db, header)
}

func (c *AutonityContract) Committee(header *types.Header, db vm.StateDB) (*types.Committee, error) {
	return c.callGetCommittee(db, header)
}

func (c *AutonityContract) MinimumBaseFee(block *types.Header, db vm.StateDB) (*big.Int, error) {
	if block.Number.Uint64() <= 1 {
		return new(big.Int).SetUint64(c.chainConfig.AutonityContractConfig.MinBaseFee), nil
	}
	return c.callGetMinimumBaseFee(db, block)
}

func (c *AutonityContract) EpochPeriod(block *types.Header, db vm.StateDB) (*big.Int, error) {
	return c.callGetEpochPeriod(db, block)
}

func (c *AutonityContract) EpochID(block *types.Header, db vm.StateDB) (*big.Int, error) {
	return c.callEpochID(db, block)
}

// Proposer election is now computed by committee structure, it is on longer depends on AC contract.
func (c *AutonityContract) Proposer(committee *types.Committee, _ vm.StateDB, height uint64, round int64) (proposer common.Address) {
	c.Lock()
	defer c.Unlock()

	needsTrimming := false

	_, ok := c.proposers[height]
	if !ok {
		c.proposers[height] = make(map[int64]common.Address)
		needsTrimming = true // cannot trim here, we might delete the entry we just created if it is for a very old height
	}

	proposer, ok = c.proposers[height][round]
	if !ok {
		proposer = committee.Proposer(height, round)
		c.proposers[height][round] = proposer
	}

	if needsTrimming {
		// since the chain monotonically grows towards higher heights, we can use the latest buffered height to trim the cache
		// NOTE: this is not 100% accurate because we could be computing Proposer for a very old height. But it works anyways.
		c.trimProposerCache(height)
	}

	return proposer
}

/* the Proposer election function is called from core and from the fault detector.
*
* Core calls it only on current height messages, therefore we don't need to retain old (h,r) proposers for it
* note that it is still good to have a cache for current height, because core will call multiple times the Proposer function for the same (h,r)
*
* The fault detector calls it when checking a proposal and when validating an InvalidProposer fault. For the proposal checking, we will accepts messages that
* are in [currentHeight - accountability.accountabilityHeightRange, currentHeight]. However, I think that it will be rare to process old messages, unless there is
* a malicious peer which is spamming old height messages.
* For proof validation instead, there is no height boundary (could be a proof for a very old height).
 */
func (c *AutonityContract) trimProposerCache(height uint64) {
	// we trim the cache only when the size is double `thresholdSize`
	if height <= thresholdSize*2 || len(c.proposers) <= thresholdSize*2 {
		return
	}

	for h := range c.proposers {
		if h < height-thresholdSize {
			delete(c.proposers, h)
		}
	}
}

func (c *AutonityContract) FinalizeAndGetCommittee(header *types.Header, statedb vm.StateDB) (*types.Receipt, *types.Epoch, error) {
	log.Debug("Finalizing block",
		"balance", statedb.GetBalance(params.AutonityContractAddress),
		"block", header.Number.Uint64())
	upgradeContract, epochInfo, err := c.callFinalize(statedb, header)
	if err != nil {
		return nil, nil, err
	}

	// Create a new receipt for the finalize call
	receipt := types.NewReceipt(nil, false, 0)
	receipt.TxHash = common.ACHash(header.Number)
	receipt.GasUsed = 0
	receipt.Logs = statedb.GetLogs(receipt.TxHash, header.Hash())
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	receipt.BlockHash = header.Hash()
	receipt.BlockNumber = header.Number
	receipt.TransactionIndex = uint(statedb.TxIndex())

	if upgradeContract {
		// warning prints for failure rather than returning error to stuck engine.
		// in any failure, the state will be rollback to snapshot.
		if err = c.upgradeAutonityContract(statedb, header); err != nil {
			log.Warn("Autonity Contracts Upgrade Failed", "err", err)
		}
	}
	return receipt, epochInfo, nil
}

func (c *AutonityContract) upgradeAutonityContract(statedb vm.StateDB, header *types.Header) error {
	log.Info("Initiating Autonity Contracts upgrade", "header", header.Number.Uint64())

	// get contract binary and abi set by system operator before.
	bytecode, newAbi, errContract := c.callRetrieveContract(statedb, header)
	if errContract != nil {
		return errContract
	}

	// take snapshot in case of roll back to former view.
	snapshot := statedb.Snapshot()
	if err := c.replaceAutonityBytecode(header, statedb, bytecode); err != nil {
		statedb.RevertToSnapshot(snapshot)
		return err
	}

	// save new abi in persistent, once node reset, it load from persistent level db.
	if err := rawdb.PutKeyValue(c.db, AutonityABIKey, []byte(newAbi)); err != nil {
		statedb.RevertToSnapshot(snapshot)
		return err
	}

	// upgrade c.ContractStateStore too right after the contract upgrade successfully.
	if err := c.upgradeAbiCache(newAbi); err != nil {
		statedb.RevertToSnapshot(snapshot)
		return err
	}
	log.Info("Autonity Contract upgrade success")
	return nil
}

func (c *AutonityContract) CallContractFunc(statedb vm.StateDB, header *types.Header, packedArgs []byte) ([]byte, uint64, error) {
	return c.EVMContract.CallContractFunc(statedb, header, params.AutonityContractAddress, packedArgs)
}

func (c *AutonityContract) CallContractFuncAs(statedb vm.StateDB, header *types.Header, origin common.Address, packedArgs []byte) ([]byte, error) {
	return c.EVMContract.CallContractFuncAs(statedb, header, params.AutonityContractAddress, origin, packedArgs)
}

func (c *NonStakeableVestingContract) CallContractFuncAs(statedb vm.StateDB, header *types.Header, origin common.Address, packedArgs []byte) ([]byte, error) {
	return c.EVMContract.CallContractFuncAs(statedb, header, params.NonStakeableVestingContractAddress, origin, packedArgs)
}

func (c *StakeableVestingManagerContract) CallContractFuncAs(statedb vm.StateDB, header *types.Header, origin common.Address, packedArgs []byte) ([]byte, error) {
	return c.EVMContract.CallContractFuncAs(statedb, header, params.StakeableVestingManagerContractAddress, origin, packedArgs)
}

func (c *EVMContract) upgradeAbiCache(newAbi string) error {
	c.Lock()
	defer c.Unlock()
	newABI, err := abi.JSON(strings.NewReader(newAbi))
	if err != nil {
		return err
	}
	c.contractABI = &newABI
	return nil
}

// ABI returns the current autonity contract's ABI
func (c *EVMContract) ABI() *abi.ABI {
	return c.contractABI
}

func (c *EVMContract) DeployContractWithValue(header *types.Header, origin common.Address, statedb vm.StateDB, bytecode []byte, value *big.Int, args ...interface{}) error {
	constructorParams, err := c.contractABI.Pack("", args...)
	if err != nil {
		log.Error("contractABI.Pack returns err", "err", err)
		return err
	}

	data := append(bytecode, constructorParams...)
	gas := uint64(0xFFFFFFFF)

	evm := c.evmProvider(header, origin, statedb)

	_, _, _, vmerr := evm.Create(vm.AccountRef(params.DeployerAddress), data, gas, value)
	if vmerr != nil {
		log.Error("evm.Create failed", "err", vmerr)
		return vmerr
	}

	return nil
}

func (c *EVMContract) DeployContract(header *types.Header, origin common.Address, statedb vm.StateDB, bytecode []byte, args ...interface{}) error {
	return c.DeployContractWithValue(header, origin, statedb, bytecode, new(big.Int).SetUint64(0x00), args...)
}

type OracleContract struct {
	EVMContract
}

type AccountabilityContract struct {
	EVMContract
}

type ACUContract struct {
	EVMContract
}

type SupplyControlContract struct {
	EVMContract
}

type StabilizationContract struct {
	EVMContract
}

type UpgradeManagerContract struct {
	EVMContract
}

type InflationControllerContract struct {
	EVMContract
}

type StakeableVestingManagerContract struct {
	EVMContract
}

type NonStakeableVestingContract struct {
	EVMContract
}

type OmissionAccountabilityContract struct {
	EVMContract
}

func NewGenesisEVMContract(genesisEvmProvider GenesisEVMProvider, statedb vm.StateDB, db ethdb.Database, chainConfig *params.ChainConfig) *GenesisEVMContracts {
	evmProvider := func(header *types.Header, origin common.Address, statedb vm.StateDB) *vm.EVM {
		if header != nil {
			panic("genesis evm provider must use nil header")
		}
		return genesisEvmProvider(statedb)
	}
	return &GenesisEVMContracts{
		AutonityContract: AutonityContract{
			EVMContract: EVMContract{
				evmProvider: evmProvider,
				contractABI: &generated.AutonityAbi,
				db:          db,
				chainConfig: chainConfig,
			},
			proposers: make(map[uint64]map[int64]common.Address),
		},
		AccountabilityContract: AccountabilityContract{
			EVMContract: EVMContract{
				evmProvider: evmProvider,
				contractABI: &generated.AccountabilityAbi,
				db:          db,
				chainConfig: chainConfig,
			},
		},
		OmissionAccountabilityContract: OmissionAccountabilityContract{
			EVMContract: EVMContract{
				evmProvider: evmProvider,
				contractABI: &generated.OmissionAccountabilityAbi,
				db:          db,
				chainConfig: chainConfig,
			},
		},
		OracleContract: OracleContract{
			EVMContract: EVMContract{
				evmProvider: evmProvider,
				contractABI: &generated.OracleAbi,
				db:          db,
				chainConfig: chainConfig,
			},
		},
		ACUContract: ACUContract{
			EVMContract: EVMContract{
				evmProvider: evmProvider,
				contractABI: &generated.ACUAbi,
				db:          db,
				chainConfig: chainConfig,
			}},
		SupplyControlContract: SupplyControlContract{
			EVMContract: EVMContract{
				evmProvider: evmProvider,
				contractABI: &generated.SupplyControlAbi,
				db:          db,
				chainConfig: chainConfig,
			},
		},
		StabilizationContract: StabilizationContract{
			EVMContract{
				evmProvider: evmProvider,
				contractABI: &generated.StabilizationAbi,
				db:          db,
				chainConfig: chainConfig,
			},
		},
		UpgradeManagerContract: UpgradeManagerContract{
			EVMContract{
				evmProvider: evmProvider,
				contractABI: &generated.UpgradeManagerAbi,
				db:          db,
				chainConfig: chainConfig,
			},
		},
		InflationControllerContract: InflationControllerContract{
			EVMContract{
				evmProvider: evmProvider,
				contractABI: &generated.InflationControllerAbi,
				db:          db,
				chainConfig: chainConfig,
			},
		},
		StakeableVestingManagerContract: StakeableVestingManagerContract{
			EVMContract{
				evmProvider: evmProvider,
				contractABI: &generated.StakeableVestingManagerAbi,
				db:          db,
				chainConfig: chainConfig,
			},
		},
		NonStakeableVestingContract: NonStakeableVestingContract{
			EVMContract{
				evmProvider: evmProvider,
				contractABI: &generated.NonStakeableVestingAbi,
				db:          db,
				chainConfig: chainConfig,
			},
		},
		statedb: statedb,
	}
}

type GenesisEVMContracts struct {
	AutonityContract
	AccountabilityContract
	OracleContract
	ACUContract
	SupplyControlContract
	StabilizationContract
	UpgradeManagerContract
	InflationControllerContract
	StakeableVestingManagerContract
	NonStakeableVestingContract
	OmissionAccountabilityContract
	statedb vm.StateDB
}

func (c *GenesisEVMContracts) DeployAutonityContract(bytecode []byte, validators []params.Validator, config AutonityConfig) error {
	return c.AutonityContract.DeployContract(nil, params.DeployerAddress, c.statedb, bytecode, validators, config)
}

func (c *GenesisEVMContracts) DeployAccountabilityContract(autonityAddress common.Address, config AccountabilityConfig, bytecode []byte) error {
	return c.AccountabilityContract.DeployContract(nil, params.DeployerAddress, c.statedb, bytecode, autonityAddress, config)
}

func (c *GenesisEVMContracts) DeployOmissionAccountabilityContract(autonityAddress common.Address, operator common.Address, treasuries []common.Address, config OmissionAccountabilityConfig, bytecode []byte) error {
	return c.OmissionAccountabilityContract.DeployContract(nil, params.DeployerAddress, c.statedb, bytecode, autonityAddress, operator, treasuries, config)
}

func (c *GenesisEVMContracts) Mint(address common.Address, amount *big.Int) error {
	return c.AutonityContract.Mint(nil, c.statedb, address, amount)
}

func (c *GenesisEVMContracts) Bond(from common.Address, validatorAddress common.Address, amount *big.Int) error {
	return c.AutonityContract.Bond(nil, c.statedb, from, validatorAddress, amount)
}

func (c *GenesisEVMContracts) FinalizeInitialization(delta uint64) error {
	return c.AutonityContract.FinalizeInitialization(nil, c.statedb, delta)
}

func (c *GenesisEVMContracts) DeployOracleContract(
	voters []common.Address,
	autonityAddress common.Address,
	operator common.Address,
	symbols []string,
	votePeriod *big.Int,
	bytecode []byte,
) error {
	return c.OracleContract.DeployContract(nil, params.DeployerAddress, c.statedb, bytecode, voters, autonityAddress, operator, symbols, votePeriod)
}

func (c *GenesisEVMContracts) DeployACUContract(
	symbols []string,
	quantities []*big.Int,
	scale *big.Int,
	autonity common.Address,
	operator common.Address,
	oracle common.Address,
	bytecode []byte,
) error {
	return c.ACUContract.DeployContract(nil, params.DeployerAddress, c.statedb, bytecode, symbols, quantities, scale, autonity, operator, oracle)
}

func (c *GenesisEVMContracts) DeploySupplyControlContract(autonity common.Address, operator common.Address, stabilizationContract common.Address, bytecode []byte, value *big.Int) error {
	return c.SupplyControlContract.DeployContractWithValue(nil, params.DeployerAddress, c.statedb, bytecode, value, autonity, operator, stabilizationContract)
}

func (c *GenesisEVMContracts) AddBalance(address common.Address, value *big.Int) {
	c.statedb.AddBalance(address, value)
}

func (c *GenesisEVMContracts) DeployStabilizationContract(
	stabilizationConfig StabilizationConfig,
	autonity common.Address,
	operator common.Address,
	oracle common.Address,
	supplyControl common.Address,
	collateral common.Address,
	bytecode []byte,
) error {
	return c.StabilizationContract.DeployContract(nil, params.DeployerAddress, c.statedb, bytecode, stabilizationConfig, autonity, operator, oracle, supplyControl, collateral)
}

func (c *GenesisEVMContracts) DeployUpgradeManagerContract(autonityAddress common.Address, operatorAddress common.Address, bytecode []byte) error {
	return c.UpgradeManagerContract.DeployContract(nil, params.DeployerAddress, c.statedb, bytecode, autonityAddress, operatorAddress)
}

func (c *GenesisEVMContracts) DeployInflationControllerContract(bytecode []byte, param InflationControllerParams) error {
	return c.InflationControllerContract.DeployContract(nil, params.DeployerAddress, c.statedb, bytecode, param)
}

func (c *GenesisEVMContracts) DeployStakeableVestingContract(bytecode []byte, autonityContract common.Address) error {
	return c.StakeableVestingManagerContract.DeployContract(nil, params.DeployerAddress, c.statedb, bytecode, autonityContract)
}

func (c *GenesisEVMContracts) NewStakeableContract(contract params.StakeableVestingData) error {
	return c.StakeableVestingManagerContract.NewContract(nil, c.statedb, contract)
}

func (c *GenesisEVMContracts) DeployNonStakeableVestingContract(bytecode []byte, autonityContract common.Address) error {
	return c.NonStakeableVestingContract.DeployContract(nil, params.DeployerAddress, c.statedb, bytecode, autonityContract)
}

func (c *GenesisEVMContracts) CreateSchedule(schedule params.Schedule) error {
	return c.AutonityContract.CreateSchedule(nil, c.statedb, schedule.VaultAddress, schedule)
}

func (c *GenesisEVMContracts) NewNonStakeableContract(contract params.NonStakeableVestingData) error {
	return c.NonStakeableVestingContract.NewContract(nil, c.statedb, contract)
}
