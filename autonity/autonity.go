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
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/ethdb"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
)

var (
	DeployerAddress               = common.Address{}
	AutonityContractAddress       = crypto.CreateAddress(DeployerAddress, 0)
	AccountabilityContractAddress = crypto.CreateAddress(DeployerAddress, 1)
	OracleContractAddress         = crypto.CreateAddress(DeployerAddress, 2)
	ACUContractAddress            = crypto.CreateAddress(DeployerAddress, 3)
	SupplyControlContractAddress  = crypto.CreateAddress(DeployerAddress, 4)
	StabilizationContractAddress  = crypto.CreateAddress(DeployerAddress, 5)
	AutonityABIKey                = []byte("ABISPEC")
)

// Errors
var (
	ErrAutonityContract = errors.New("could not call Autonity contract")
	ErrWrongParameter   = errors.New("wrong parameter")
	ErrNoAutonityConfig = errors.New("autonity section missing from genesis")
)

// EVMProvider provides a new evm. This allows us to decouple the contract from *core.blockchain which is required to build a new evm.
type EVMProvider func(header *types.Header, origin common.Address, statedb *state.StateDB) *vm.EVM

// GenesisEVMProvider returns new, preconfigured EVM. Useful for genesis EVMs which depends on fixed set of parameters
// which can be closed in this function
type GenesisEVMProvider func(statedb *state.StateDB) *vm.EVM

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
	minBaseFee    atomic.Pointer[big.Int]
	minBaseFeeCh  chan *AutonityMinimumBaseFeeUpdated
	subMinBaseFee event.Subscription
	subscriptions *event.SubscriptionScope // will be useful when we have multiple subscriptions
	quit          chan struct{}
	done          chan struct{}
}

func newCache(ac *AutonityContract, head *types.Header, state *state.StateDB) (*Cache, error) {
	minBaseFee, err := ac.MinimumBaseFee(head, state)
	if err != nil {
		return nil, err
	}
	minBaseFeeCh := make(chan *AutonityMinimumBaseFeeUpdated)
	subMinBaseFee, err := ac.WatchMinimumBaseFeeUpdated(nil, minBaseFeeCh)
	if err != nil {
		return nil, err
	}
	scope := new(event.SubscriptionScope)
	subMinBaseFeeWrapped := scope.Track(subMinBaseFee)
	cache := &Cache{
		minBaseFeeCh:  minBaseFeeCh,
		subMinBaseFee: subMinBaseFeeWrapped,
		subscriptions: scope,
		quit:          make(chan struct{}),
		done:          make(chan struct{}),
	}
	cache.minBaseFee.Store(minBaseFee)
	go cache.Listen()
	return cache, nil
}

//revive:disable:exported - Autonity is one of the contracts, so repetitive naming here is justified
type AutonityContract struct {
	EVMContract
	*AutonityFilterer // allows to watch for Autonity Contract events

	proposers map[uint64]map[int64]common.Address
}

type ProtocolContracts struct {
	*AutonityContract
	*Cache
	*Accountability
}

func NewProtocolContracts(config *params.ChainConfig, db ethdb.Database, provider EVMProvider, contractBackend bind.ContractBackend, head *types.Header, state *state.StateDB) (*ProtocolContracts, error) {
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
	autonityFilterer, err := NewAutonityFilterer(AutonityContractAddress, contractBackend.(bind.ContractFilterer))
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
	accountabilityContract, _ := NewAccountability(AccountabilityContractAddress, contractBackend)

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
		case <-c.subMinBaseFee.Err():
			// This should never happen. Errors from subscription can happen only if the subscription is done over an RPC connection.
			// In that case network errors can occur. Since everything is local here, no error should ever occur.
			// we crash the client to avoid using a out-of-date value of minBaseFee in case something goes very wrong.
			log.Crit("protocol contract cache out-of-sync. Please contact the Autonity team.")
		case ev := <-c.minBaseFeeCh:
			c.minBaseFee.Store(ev.GasPrice)
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

func (c *AutonityContract) CommitteeEnodes(header *types.Header, db *state.StateDB) (*types.Nodes, error) {
	return c.callGetCommitteeEnodes(db, header)
}

func (c *AutonityContract) MinimumBaseFee(block *types.Header, db *state.StateDB) (*big.Int, error) {
	if block.Number.Uint64() <= 1 {
		return new(big.Int).SetUint64(c.chainConfig.AutonityContractConfig.MinBaseFee), nil
	}
	return c.callGetMinimumBaseFee(db, block)
}

func (c *AutonityContract) Proposer(committee *types.Committee, height uint64, round int64) common.Address {
	c.Lock()
	defer c.Unlock()
	roundMap, ok := c.proposers[height]
	if !ok {
		p := c.electProposer(committee, height, round)
		roundMap = make(map[int64]common.Address)
		roundMap[round] = p
		c.proposers[height] = roundMap
		// since the growing of blockchain height, so most of the case we can use the latest buffered height to gc the buffer.
		c.trimProposerCache(height)
		return p
	}

	proposer, ok := roundMap[round]
	if !ok {
		p := c.electProposer(committee, height, round)
		c.proposers[height][round] = p
		return p
	}

	return proposer
}

// electProposer is a part of consensus, that it elect proposer from epoch head's committee list which was returned
// from autonity contract stable ordered by voting power in evm context.
func (c *AutonityContract) electProposer(committee *types.Committee, height uint64, round int64) common.Address {
	seed := big.NewInt(constants.MaxRound)
	totalVotingPower := committee.TotalVotingPower()

	// for power weighted sampling, we distribute seed into a 256bits key-space, and compute the hit index.
	h := new(big.Int).SetUint64(height)
	r := new(big.Int).SetInt64(round)
	key := r.Add(r, h.Mul(h, seed))
	value := new(big.Int).SetBytes(crypto.Keccak256(key.Bytes()))
	index := value.Mod(value, totalVotingPower)
	// find the index hit which committee member which line up in the committee list.
	// we assume there is no 0 stake/power validators.
	counter := new(big.Int).SetUint64(0)
	for _, member := range committee.Members {
		counter.Add(counter, member.VotingPower)
		if index.Cmp(counter) == -1 {
			return member.Address
		}
	}

	// otherwise, we elect with round-robin.
	return committee.Members[round%int64(committee.Len())].Address
}

func (c *AutonityContract) trimProposerCache(height uint64) {
	if len(c.proposers) > consensus.AccountabilityHeightRange*2 && height > consensus.AccountabilityHeightRange*2 {
		gcFrom := height - consensus.AccountabilityHeightRange*2
		// keep at least consensus.AccountabilityHeightRange heights in the buffer.
		for i := uint64(0); i < consensus.AccountabilityHeightRange; i++ {
			gcHeight := gcFrom + i
			delete(c.proposers, gcHeight)
		}
	}
}

func (c *AutonityContract) FinalizeAndGetCommittee(header *types.Header, statedb *state.StateDB) ([]types.CommitteeMember, *types.Receipt, *big.Int, error) {
	if header.Number.Uint64() == 0 {
		return nil, nil, nil, nil
	}

	log.Debug("Finalizing block",
		"balance", statedb.GetBalance(AutonityContractAddress),
		"block", header.Number.Uint64())

	upgradeContract, committee, lastEpochBlock, err := c.callFinalize(statedb, header)

	if err != nil {
		return nil, nil, nil, err
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
	return committee, receipt, lastEpochBlock, nil
}

func (c *AutonityContract) upgradeAutonityContract(statedb *state.StateDB, header *types.Header) error {
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

func (c *AutonityContract) CallContractFunc(statedb *state.StateDB, header *types.Header, packedArgs []byte) ([]byte, error) {
	return c.EVMContract.CallContractFunc(statedb, header, AutonityContractAddress, packedArgs)
}

func (c *AutonityContract) CallContractFuncAs(statedb *state.StateDB, header *types.Header, origin common.Address, packedArgs []byte) ([]byte, error) {
	return c.EVMContract.CallContractFuncAs(statedb, header, AutonityContractAddress, origin, packedArgs)
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

func (c *EVMContract) DeployContractWithValue(header *types.Header, origin common.Address, statedb *state.StateDB, bytecode []byte, value *big.Int, args ...interface{}) error {
	constructorParams, err := c.contractABI.Pack("", args...)
	if err != nil {
		log.Error("contractABI.Pack returns err", "err", err)
		return err
	}

	data := append(bytecode, constructorParams...)
	gas := uint64(0xFFFFFFFF)

	evm := c.evmProvider(header, origin, statedb)

	_, _, _, vmerr := evm.Create(vm.AccountRef(DeployerAddress), data, gas, value)
	if vmerr != nil {
		log.Error("evm.Create failed", "err", vmerr)
		return vmerr
	}

	return nil
}

func (c *EVMContract) DeployContract(header *types.Header, origin common.Address, statedb *state.StateDB, bytecode []byte, args ...interface{}) error {
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

func NewGenesisEVMContract(genesisEvmProvider GenesisEVMProvider, statedb *state.StateDB, db ethdb.Database, chainConfig *params.ChainConfig) *GenesisEVMContracts {
	evmProvider := func(header *types.Header, origin common.Address, statedb *state.StateDB) *vm.EVM {
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

	statedb *state.StateDB
}

func (c *GenesisEVMContracts) DeployAutonityContract(bytecode []byte, validators []params.Validator, config AutonityConfig) error {
	return c.AutonityContract.DeployContract(nil, DeployerAddress, c.statedb, bytecode, validators, config)
}

func (c *GenesisEVMContracts) DeployAccountabilityContract(autonityAddress common.Address, config AccountabilityConfig, bytecode []byte) error {
	return c.AccountabilityContract.DeployContract(nil, DeployerAddress, c.statedb, bytecode, autonityAddress, config)
}

func (c *GenesisEVMContracts) Mint(address common.Address, amount *big.Int) error {
	return c.AutonityContract.Mint(nil, c.statedb, address, amount)
}

func (c *GenesisEVMContracts) Bond(from common.Address, validatorAddress common.Address, amount *big.Int) error {
	return c.AutonityContract.Bond(nil, c.statedb, from, validatorAddress, amount)
}

func (c *GenesisEVMContracts) FinalizeInitialization() error {
	return c.AutonityContract.FinalizeInitialization(nil, c.statedb)
}

func (c *GenesisEVMContracts) DeployOracleContract(
	voters []common.Address,
	autonityAddress common.Address,
	operator common.Address,
	symbols []string,
	votePeriod *big.Int,
	bytecode []byte,
) error {
	return c.OracleContract.DeployContract(nil, DeployerAddress, c.statedb, bytecode, voters, autonityAddress, operator, symbols, votePeriod)
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
	return c.ACUContract.DeployContract(nil, DeployerAddress, c.statedb, bytecode, symbols, quantities, scale, autonity, operator, oracle)
}

func (c *GenesisEVMContracts) DeploySupplyControlContract(autonity common.Address, operator common.Address, stabilizationContract common.Address, bytecode []byte, value *big.Int) error {
	return c.SupplyControlContract.DeployContractWithValue(nil, DeployerAddress, c.statedb, bytecode, value, autonity, operator, stabilizationContract)
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
	return c.StabilizationContract.DeployContract(nil, DeployerAddress, c.statedb, bytecode, stabilizationConfig, autonity, operator, oracle, supplyControl, collateral)
}
